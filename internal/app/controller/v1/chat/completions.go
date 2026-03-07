// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package chat

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/log"
	"gotribe/pkg/llm"

	"github.com/gin-gonic/gin"
)

const (
	contentTypeJSON = "application/json"
	contentTypeSSE  = "text/event-stream"
	openai          = "openai"
)

// Completions 转发 POST /v1/chat/completions，支持 stream=true（SSE）与 stream=false（JSON）.
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

	backend, err := llm.SelectBackend(modelVal)
	if err != nil {
		if err == llm.ErrLLMConfigNotFound || err == llm.ErrNoBackend {
			core.WriteResponse(c, errno.ErrPageNotFound, nil)
			return
		}
		core.WriteResponse(c, errno.InternalServerError, nil)
		return
	}

	// 使用后端配置的 model 请求上游
	reqMap["model"] = backend.Model
	upstreamBody, err := json.Marshal(reqMap)
	if err != nil {
		core.WriteResponse(c, errno.InternalServerError, nil)
		return
	}

	upstreamURL := backend.BaseURL
	if backend.Type == openai {
		upstreamURL = backend.BaseURL + "/chat/completions"
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, upstreamURL, bytes.NewReader(upstreamBody))
	if err != nil {
		core.WriteResponse(c, errno.InternalServerError, nil)
		return
	}

	apikey := backend.APIKey
	if backend.Type == openai {
		apikey = "Bearer " + backend.APIKey
	}
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Authorization", apikey)
	for k, v := range backend.Header {
		req.Header.Set(k, v)
	}
	//for k, v := range c.Request.Header {
	//	if strings.EqualFold(k, "Content-Length") || strings.EqualFold(k, "Authorization") {
	//		continue
	//	}
	//	for _, vv := range v {
	//		req.Header.Add(k, vv)
	//	}
	//}

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
		ctrl.writeStreamResponse(c, resp)
	} else {
		ctrl.writeJSONResponse(c, resp)
	}
}

func (ctrl *ChatController) writeStreamResponse(c *gin.Context, resp *http.Response) {
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
	_, _ = io.Copy(c.Writer, resp.Body)
}

func (ctrl *ChatController) writeJSONResponse(c *gin.Context, resp *http.Response) {
	c.Status(resp.StatusCode)
	for k, v := range resp.Header {
		for _, vv := range v {
			c.Writer.Header().Add(k, vv)
		}
	}
	_, _ = io.Copy(c.Writer, resp.Body)
}
