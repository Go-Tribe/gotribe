package ip

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// IPInfo 结构体用于解析 JSON 响应
type IPInfo struct {
	Status     string `json:"status"`
	Country    string `json:"country"`
	City       string `json:"city"`
	RegionName string `json:"regionName"`
}

// IsPrivateIP 检查 IP 地址是否为内网地址
func IsPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return true
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		default:
			return false
		}
	}
	return false
}

// GetIPInfo 获取 IP 信息
func GetIPInfo(ctx context.Context, ipAddress string) (string, string, string) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "内网IP", "内网IP", "内网IP"
	}

	if IsPrivateIP(ip) {
		return "内网IP", "内网IP", "内网IP"
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ipAddress)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "内网IP", "内网IP", "内网IP"
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "内网IP", "内网IP", "内网IP"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "内网IP", "内网IP", "内网IP"
	}

	var ipInfo IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
		return "内网IP", "内网IP", "内网IP"
	}

	if ipInfo.Status != "success" {
		return "内网IP", "内网IP", "内网IP"
	}

	return ipInfo.Country, ipInfo.City, ipInfo.RegionName
}
