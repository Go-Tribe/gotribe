// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
	"gotribe/pkg/api/v1"
	"gotribe/pkg/llm"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	contentTypeJSON = "application/json"
	contentTypeSSE  = "text/event-stream"
	openai          = "openai"

	pointLogTypeChat   = "chat"
	pointLogReasonChat = "对话消耗"
)

// usage 用于解析 OpenAI 风格返回的 usage，支持缓存 token（cached_tokens / prompt_tokens_details.cached_tokens）.
type usage struct {
	PromptTokens        int `json:"prompt_tokens"`
	CompletionTokens    int `json:"completion_tokens"`
	TotalTokens         int `json:"total_tokens"`
	CachedTokens        int `json:"cached_tokens"`
	PromptTokensDetails *struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"prompt_tokens_details"`
}

// getCachedTokens 先取 CachedTokens，没有则取 PromptTokensDetails.CachedTokens，都没有则为 0.
func (u *usage) getCachedTokens() int {
	if u == nil {
		return 0
	}
	if u.CachedTokens > 0 {
		return u.CachedTokens
	}
	if u.PromptTokensDetails != nil && u.PromptTokensDetails.CachedTokens > 0 {
		return u.PromptTokensDetails.CachedTokens
	}
	return 0
}

// Completions 转发 POST /v1/chat/completions，支持 stream=true（SSE）与 stream=false（JSON），并做对话扣费.
func (ctrl *ChatController) Completions(c *gin.Context) {
	log.C(c).Infow("Chat completions proxy called")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	var reqMap map[string]interface{}
	if err := json.Unmarshal(body, &reqMap); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	modelVal, _ := reqMap["model"].(string)
	if modelVal == "" {
		core.WriteResponse(c, errno.ErrInvalidParameter, nil)
		return
	}

	username := c.GetString(known.XUsernameKey)
	projectID := c.GetString(known.XProjectIDKey)
	user, err := ctrl.ds.Users().Get(c.Request.Context(), v1.UserWhere{Username: username})
	if err != nil {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	// 创建 conversation_log（接口一进来就创建）
	convLog, err := ctrl.createConversationLog(c, user, projectID, modelVal)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	// BeforeCompletions：点数不足则拒绝
	if err := ctrl.BeforeCompletions(c, user, projectID); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	// 根据配置 model 选 name，再按 name 从 llm_model 表取第一条
	backendName, err := llm.SelectBackendName(modelVal)
	if err != nil {
		if err == llm.ErrLLMConfigNotFound || err == llm.ErrNoBackend {
			core.WriteResponse(c, errno.ErrPageNotFound, nil)
			return
		}
		core.WriteResponse(c, errno.InternalServerError, nil)
		return
	}

	llmModel, err := ctrl.ds.LLMModels().GetByName(c.Request.Context(), backendName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			core.WriteResponse(c, errno.ErrPageNotFound, nil)
		} else {
			core.WriteResponse(c, errno.InternalServerError, nil)
		}
		return
	}

	// 使用库中的 model 请求上游
	reqMap["model"] = llmModel.ModelName
	upstreamBody, err := json.Marshal(reqMap)
	if err != nil {
		core.WriteResponse(c, errno.InternalServerError, nil)
		return
	}

	upstreamURL := strings.TrimSuffix(llmModel.BaseURL, "/")
	if llmModel.Type == openai {
		upstreamURL = upstreamURL + "/chat/completions"
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, upstreamURL, bytes.NewReader(upstreamBody))
	if err != nil {
		core.WriteResponse(c, errno.InternalServerError, nil)
		return
	}

	apikey := llmModel.ApiKey
	if llmModel.Type == openai {
		apikey = "Bearer " + llmModel.ApiKey
	}
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Authorization", apikey)
	ctrl.setHeadersFromLLMModel(req, llmModel)
	ctrl.setClientHeaders(c, req, llmModel)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.C(c).Errorw("Upstream LLM request failed", "err", err)
		core.WriteResponse(c, errno.InternalServerError, nil)
		return
	}
	defer resp.Body.Close()

	stream, _ := reqMap["stream"].(bool)
	if stream {
		usageOut := ctrl.writeStreamResponse(c, resp)
		ctrl.AfterCompletions(c, convLog, llmModel, usageOut)
	} else {
		respBody, _ := io.ReadAll(resp.Body)
		usageOut := ctrl.parseUsageFromJSON(respBody)
		ctrl.writeJSONResponse(c, resp, respBody)
		ctrl.AfterCompletions(c, convLog, llmModel, usageOut)
	}
}

func (ctrl *ChatController) setHeadersFromLLMModel(req *http.Request, m *model.LLMModelM) {
	if m.Header == "" {
		return
	}
	var h map[string]string
	if _ = json.Unmarshal([]byte(m.Header), &h); h != nil {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}
}

// setClientHeaders 根据 llm_model.ClientHeader（JSON 数组字符串，配置的 header key），从当前请求头取出对应 key 的值塞入上游请求头.
func (ctrl *ChatController) setClientHeaders(c *gin.Context, req *http.Request, m *model.LLMModelM) {
	if m.ClientHeader == "" {
		return
	}
	var keys []string
	if err := json.Unmarshal([]byte(m.ClientHeader), &keys); err != nil || len(keys) == 0 {
		return
	}
	for _, k := range keys {
		v := c.Request.Header.Get(k)
		if v != "" {
			req.Header.Set(k, v)
		}
	}
}

// createConversationLog 创建一条 conversation_log，并返回带 SID 的 log（用于后续扣费 eventID）.
func (ctrl *ChatController) createConversationLog(c *gin.Context, user *model.UserM, projectID, modelName string) (*model.ConversationLogM, error) {
	requestID := c.GetHeader(known.XRequestIDKey)
	appVersionId := c.GetHeader(known.XAppVersionId)
	clientName := c.GetHeader(known.XClientName)
	clientVersion := c.GetHeader(known.XClientVersion)
	deviceModel := c.GetHeader(known.XDeviceModel)
	os := c.GetHeader(known.XOS)
	osVersion := c.GetHeader(known.XOSVersion)
	osArch := c.GetHeader(known.XOSArch)
	cpuModel := c.GetHeader(known.XCPUModel)

	convLog := &model.ConversationLogM{
		UserID:          user.UserID,
		Username:        user.Username,
		ModelName:       modelName,
		StartTime:       time.Now(),
		RequestID:       requestID,
		Status:          1,
		DeductionStatus: known.ClDeductionStatusPending,
		DeductionType:   known.ClDeductionTypeWhite,
		IPAddress:       c.ClientIP(),
		AppVersionId:    appVersionId,
		ClientName:      clientName,
		ClientVersion:   clientVersion,
		DeviceModel:     deviceModel,
		OS:              os,
		OSVersion:       osVersion,
		OSArch:          osArch,
		CPUModel:        cpuModel,
		Ext:             "{}",
	}
	if err := ctrl.ds.ConversationLogs().Create(c.Request.Context(), convLog); err != nil {
		return nil, errno.InternalServerError
	}
	return convLog, nil
}

// BeforeCompletions 检查用户点数是否不低于配置的最低可对话点数，不足则返回错误.
func (ctrl *ChatController) BeforeCompletions(c *gin.Context, user *model.UserM, projectID string) error {
	username := user.Username
	user, err := ctrl.ds.Users().Get(c.Request.Context(), v1.UserWhere{Username: username})
	if err != nil {
		return errno.InternalServerError
	}
	points, err := ctrl.b.Point().GetAvailablePoints(c.Request.Context(), user.UserID, projectID)
	if err != nil {
		return errno.InternalServerError
	}
	minPoints := llm.GetMinPointsToChat()
	avail := float64(0)
	if points != nil {
		avail = *points
	}
	if avail < minPoints {
		return errno.ErrInsufficientPoints
	}
	return nil
}

// AfterCompletions 根据 usage 与 llm_model 计费，扣减用户点数并更新 conversation_log.
func (ctrl *ChatController) AfterCompletions(c *gin.Context, convLog *model.ConversationLogM, llmModel *model.LLMModelM, usageOut *usage) {
	ctx := c.Request.Context()
	now := time.Now()
	endTime := now
	startTime := convLog.StartTime
	totalMs := int(endTime.Sub(startTime).Milliseconds())

	convLog.EndTime = &endTime
	convLog.TotalTimeMs = totalMs
	if usageOut != nil {
		convLog.PromptTokens = usageOut.PromptTokens
		convLog.CompletionTokens = usageOut.CompletionTokens
		convLog.CacheTokens = usageOut.getCachedTokens()
		convLog.TotalTokens = usageOut.TotalTokens
	}

	cost := ctrl.calcPointsCost(llmModel, usageOut)
	convLog.PointsCost = cost

	// 白名单用户不扣积分，状态设为无需扣费，类型记为白名单用户
	if ctrl.isWhiteUsername(convLog.Username) {
		convLog.DeductionStatus = known.ClDeductionStatusNoDeduct
		convLog.DeductionType = known.ClDeductionTypeWhite
		_ = ctrl.ds.ConversationLogs().Update(ctx, convLog)
		return
	}

	if cost <= 0 {
		convLog.DeductionStatus = known.ClDeductionStatusNoDeduct
		_ = ctrl.ds.ConversationLogs().Update(ctx, convLog)
		return
	}

	if convLog.Username == "" {
		convLog.DeductionStatus = known.ClDeductionStatusNoDeduct
		_ = ctrl.ds.ConversationLogs().Update(ctx, convLog)
		return
	}

	err := ctrl.b.Point().SubPoints(ctx, convLog.Username, c.GetString(known.XProjectIDKey), pointLogTypeChat, pointLogReasonChat, convLog.SID, cost)
	if err != nil {
		convLog.DeductionStatus = known.ClDeductionStatusFailed
		convLog.ErrorMessage = err.Error()
	} else {
		convLog.DeductionStatus = known.ClDeductionStatusSuccess
		convLog.DeductionTime = &now
	}
	_ = ctrl.ds.ConversationLogs().Update(ctx, convLog)
}

func (ctrl *ChatController) isWhiteUsername(username string) bool {
	for _, u := range llm.GetWhiteUsernames() {
		if u == username {
			return true
		}
	}
	return false
}

func (ctrl *ChatController) calcPointsCost(m *model.LLMModelM, u *usage) float64 {
	if u == nil {
		return 0
	}
	cached := u.getCachedTokens()
	nonCachedPrompt := u.PromptTokens - cached
	if nonCachedPrompt < 0 {
		nonCachedPrompt = 0
	}
	in := float64(nonCachedPrompt) / 1000 * m.InputPrice
	cache := float64(cached) / 1000 * m.CachePrice
	out := float64(u.CompletionTokens) / 1000 * m.OutputPrice
	return in + cache + out
}

func (ctrl *ChatController) parseUsageFromJSON(body []byte) *usage {
	var v struct {
		Usage *usage `json:"usage"`
	}
	if err := json.Unmarshal(body, &v); err != nil || v.Usage == nil {
		return nil
	}
	return v.Usage
}

func (ctrl *ChatController) writeStreamResponse(c *gin.Context, resp *http.Response) (usageOut *usage) {
	c.Header("Content-Type", contentTypeSSE)
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(resp.StatusCode)
	for k, v := range resp.Header {
		if strings.EqualFold(k, "Content-Type") || strings.EqualFold(k, "Transfer-Encoding") {
			continue
		}
		for _, vv := range v {
			c.Writer.Header().Add(k, vv)
		}
	}
	c.Writer.Flush()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		c.Writer.Write(line)
		c.Writer.Write([]byte("\n"))
		c.Writer.Flush()
		if bytes.HasPrefix(line, []byte("data:")) {
			data := bytes.TrimPrefix(line, []byte("data:"))
			if bytes.Equal(data, []byte("[DONE]")) {
				continue
			}
			var chunk struct {
				Usage *usage `json:"usage"`
			}
			if json.Unmarshal(data, &chunk) == nil && chunk.Usage != nil {
				usageOut = chunk.Usage
			}
		}
	}
	return usageOut
}

func (ctrl *ChatController) writeJSONResponse(c *gin.Context, resp *http.Response, body []byte) {
	c.Status(resp.StatusCode)
	for k, v := range resp.Header {
		for _, vv := range v {
			c.Writer.Header().Add(k, vv)
		}
	}
	_, _ = c.Writer.Write(body)
}
