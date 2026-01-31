package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HTTPResponse 包含完整的HTTP响应信息
type HTTPResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// HTTPError 表示HTTP请求错误
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http %d: %s", e.StatusCode, e.Body)
}

// PostFormWithResponse 发送表单 POST 并解析 JSON 到 ptr，返回完整的响应信息
func PostForm(rawURL string, form url.Values, ptr interface{}) (*HTTPResponse, error) {
	req, err := http.NewRequest(http.MethodPost, rawURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败：%w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return doRequest(req, ptr)
}

// PostJSONWithResponse 发送JSON格式的POST请求，并解析JSON响应到ptr，返回完整的响应信息
func PostJSON(rawURL string, reqBody interface{}, ptr interface{}) (*HTTPResponse, error) {
	// 序列化请求体
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化JSON请求体失败：%w", err)
	}

	// 创建请求
	req, err := http.NewRequest(http.MethodPost, rawURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败：%w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return doRequest(req, ptr)
}

// doRequest 执行HTTP请求的核心逻辑
func doRequest(req *http.Request, ptr interface{}) (*HTTPResponse, error) {
	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败：%w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败：%w", err)
	}

	// 构造响应对象
	httpResp := &HTTPResponse{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       body,
	}

	// 尝试解析JSON响应
	if ptr != nil {
		if err := json.Unmarshal(body, ptr); err != nil && resp.StatusCode < 300 {
			return httpResp, fmt.Errorf("解析JSON响应失败：%w", err)
		}
	}

	// 检查状态码
	if resp.StatusCode >= 300 {
		return httpResp, &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       strings.TrimSpace(string(body)),
		}
	}

	return httpResp, nil
}

// GetWithAuthHeader 发送带有Authorization头的GET请求，并解析JSON响应到ptr，返回完整的响应信息
func GetWithAuthHeader(rawURL string, authToken string, ptr interface{}) (*HTTPResponse, error) {
	// 创建请求
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败：%w", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Accept", "application/json")

	return doRequest(req, ptr)
}
