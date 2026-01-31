package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	network "gmcc/pkg/network"
)

const (
	MICROSOFT_TENANT    string        = "consumers"
	MICROSOFT_CLIENT_ID string        = "feb3836f-0333-4185-8eb9-4cbf0498f947"
	MAX_CHECK_TIME      time.Duration = 15 * time.Minute
)

const (
	DEVICE_CODE_URL     string = "https://login.microsoftonline.com/%s/oauth2/v2.0/devicecode"
	TOKEN_CHECK_URL     string = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
	XBOX_LIVE_AUTH_URL  string = "https://user.auth.xboxlive.com/user/authenticate"
	XSTS_URL            string = "https://xsts.auth.xboxlive.com/xsts/authorize"
)

func init() {}

/* -------------------- Mircosoft Token -------------------- */
type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// TokenResponse 处理来自 /token 终结点的成功或错误响应
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`

	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GetDeviceCode 获取设备码
func getDeviceCode(clientID string) (*DeviceCodeResponse, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("scope", "XboxLive.signin offline_access")

	var resp DeviceCodeResponse
	_, err := network.PostForm(
		fmt.Sprintf(DEVICE_CODE_URL, MICROSOFT_TENANT),
		form,
		&resp,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll DeviceCode Token 使用从 DeviceCode 请求得到的响应轮询
func pollDeviceCodeToken(clientID string, devResp *DeviceCodeResponse) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	form.Set("device_code", devResp.DeviceCode)

	interval := time.Duration(devResp.Interval) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}

	deadline := time.Now().Add(time.Duration(devResp.ExpiresIn) * time.Second)

	for time.Now().Before(deadline) {
		var tokenResp TokenResponse
		_, err := network.PostForm(fmt.Sprintf(TOKEN_CHECK_URL, MICROSOFT_TENANT), form, &tokenResp)

		if err != nil && tokenResp.Error == "" {
			return nil, err
		}

		if tokenResp.Error != "" {
			switch tokenResp.Error {
			case "authorization_pending":
			case "authorization_declined", "bad_verification_code", "expired_token":
				return nil, fmt.Errorf("验证失败: %s: %s", tokenResp.Error, tokenResp.ErrorDescription)
			default:
				return nil, fmt.Errorf("验证失败: %s: %s", tokenResp.Error, tokenResp.ErrorDescription)
			}
		}

		if tokenResp.AccessToken != "" {
			return &tokenResp, nil
		}

		time.Sleep(interval)
	}

	return nil, fmt.Errorf("设备代码过期或轮询超时")
}

// Get Microsoft Token 获取 Microsoft 设备码并轮询获取访问令牌
func getMicrosoftToken() (*TokenResponse, error) {
	devResp, err := getDeviceCode(MICROSOFT_CLIENT_ID)
	if err != nil {
		return nil, err
	}

	fmt.Println("请在浏览器打开下面的链接并输入代码完成登录：")
	if devResp.VerificationURIComplete != "" {
		fmt.Printf("自动链接: %s\n", devResp.VerificationURIComplete)
	} else {
		fmt.Printf("链接: %s\n", devResp.VerificationURI)
		fmt.Printf("代码: %s\n", devResp.UserCode)
	}

	return pollDeviceCodeToken(MICROSOFT_CLIENT_ID, devResp)
}

func refreshMSToken(refreshToken string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("client_id", MICROSOFT_CLIENT_ID)
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	form.Set("scope", "XboxLive.signin offline_access")
	var tokenResp TokenResponse
	_, err := network.PostForm(fmt.Sprintf(TOKEN_CHECK_URL, MICROSOFT_TENANT), form, &tokenResp)
	if err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

/* -------------------- Xbox Live Token -------------------- */

type XBLRequest struct {
	Properties struct {
		AuthMethod string `json:"AuthMethod"`
		SiteName   string `json:"SiteName"`
		RpsTicket  string `json:"RpsTicket"`
	} `json:"Properties"`
	RelyingParty string `json:"RelyingParty"`
	TokenType    string `json:"TokenType"`
}

type XBLResponse struct {
	IssueInstant  time.Time `json:"IssueInstant"`
	NotAfter      time.Time `json:"NotAfter"`
	Token         string    `json:"Token"`
	DisplayClaims struct {
		Xui []struct {
			Uhs string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

func getXBLToken(microsoftAccessToken string) (*XBLResponse, error) {
	xblReq := XBLRequest{}
	xblReq.Properties.AuthMethod = "RPS"
	xblReq.Properties.SiteName = "user.auth.xboxlive.com"
	xblReq.Properties.RpsTicket = "d=" + microsoftAccessToken
	xblReq.RelyingParty = "http://auth.xboxlive.com"
	xblReq.TokenType = "JWT"
	var xblResp XBLResponse
	_, err := network.PostJSON(XBOX_LIVE_AUTH_URL, xblReq, &xblResp)
	if err != nil {
		return nil, err
	}
	return &xblResp, nil
}

/* -------------------- XSTS Token -------------------- */

type XSTSRequest struct {
	Properties struct {
		SandboxID  string   `json:"SandboxId"`
		UserTokens []string `json:"UserTokens"`
	} `json:"Properties"`
	RelyingParty string `json:"RelyingParty"`
	TokenType    string `json:"TokenType"`
}

type XSTSResponse struct {
	IssueInstant  time.Time `json:"IssueInstant"`
	NotAfter      time.Time `json:"NotAfter"`
	Token         string    `json:"Token"`
	DisplayClaims struct {
		Xui []struct {
			Uhs string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

// XSTSErrorResponse 表示XSTS授权错误的响应
type XSTSErrorResponse struct {
	Identity string `json:"Identity"`
	XErr     int    `json:"XErr"`
	Message  string `json:"Message"`
	Redirect string `json:"Redirect"`
}

// XSTSError 自定义错误类型，用于返回更友好的错误信息
type XSTSError struct {
	XErr     int
	Message  string
	Redirect string
	Identity string
}

func (e *XSTSError) Error() string {
	msg := fmt.Sprintf("XSTS错误 (XErr: %d)", e.XErr)
	if e.Message != "" {
		msg += ": " + e.Message
	}
	switch e.XErr {
	case 2148916227:
		msg += " - 账号已被Xbox封禁"
	case 2148916233:
		msg += " - 账号未注册Xbox。需通过minecraft.net登录创建"
	case 2148916235:
		msg += " - 账号所属国家/地区不支持Xbox Live服务"
	case 2148916236, 2148916237:
		msg += " - 需在Xbox页面完成成人验证（韩国）"
	case 2148916238:
		msg += " - 未成年账户（未满18岁），必须由成人账户添加到家庭组才能继续"
	case 2148916262:
		msg += " - 未知错误（很少）"
	}
	if e.Redirect != "" {
		msg += fmt.Sprintf(" (Redirect: %s)", e.Redirect)
	}
	return msg
}

func getXSTSToken(xblToken string) (*XSTSResponse, error) {
	xstsReq := XSTSRequest{}
	xstsReq.Properties.SandboxID = "RETAIL"
	xstsReq.Properties.UserTokens = []string{xblToken}
	xstsReq.RelyingParty = "rp://api.minecraftservices.com/"
	xstsReq.TokenType = "JWT"

	// 直接发送请求以便在错误时能解析响应体
	payload, err := json.Marshal(xstsReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, XSTS_URL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		// 尝试解析为 XSTS 错误响应
		var errResp XSTSErrorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.XErr != 0 {
			return nil, &XSTSError{
				XErr:     errResp.XErr,
				Message:  errResp.Message,
				Redirect: errResp.Redirect,
				Identity: errResp.Identity,
			}
		}
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var xstsResp XSTSResponse
	if err := json.Unmarshal(body, &xstsResp); err != nil {
		return nil, err
	}
	return &xstsResp, nil
}

func tokenWorkflow(microsoftToken *TokenResponse) (*XSTSResponse, error) {
	// 获取 Xbox Live 令牌
	xblResp, err := getXBLToken(microsoftToken.AccessToken)
	if err != nil {
		return nil, err
	}

	// 获取 XSTS 令牌
	xstsResp, err := getXSTSToken(xblResp.Token)
	if err != nil {
		return nil, err
	}

	return xstsResp, nil
}

/* -------------------- Final Token -------------------- */
func GetToken() (*XSTSResponse, error) {
	microsoftToken, err := getMicrosoftToken()
	if err != nil {
		return nil, err
	}

	return tokenWorkflow(microsoftToken)
}

func RefreshToken(refresh_token string) (*XSTSResponse, error) {
	microsoftToken, err := refreshMSToken(refresh_token)
	if err != nil {
		return nil, err
	}

	return tokenWorkflow(microsoftToken)
}
