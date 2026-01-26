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

	"gmcc/pkg/crypto"
)

// const clientId = ""
const CLIENT_ID = "feb3836f-0333-4185-8eb9-4cbf0498f947" // Minosoft 2 (microsoft-bixilon2)
const TENANT = "consumers"
const DEVICE_CODE_URL = "https://login.microsoftonline.com/%s/oauth2/v2.0/devicecode?mkt=%s"
const TOKEN_CHECK_URL = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
const XBOX_LIVE_AUTH_URL = "https://user.auth.xboxlive.com/user/authenticate"
const XSTS_URL = "https://xsts.auth.xboxlive.com/xsts/authorize"
const LOGIN_WITH_XBOX_URL = "https://api.minecraftservices.com/authentication/login_with_xbox"
const McTokenExpire = 24 * time.Hour

// StoredToken 统一样式，用于 MS 和 MC token 的存储
type StoredToken struct {
	PlayerID     string
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}

// 设备码响应结构
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
	ExpiresIn       int    `json:"expires_in"`
	Message         string `json:"message"`
}

type MSTokenResponse struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type XBLAuthResponse struct {
	IssueInstant  string `json:"IssueInstant"`
	NotAfter      string `json:"NotAfter"`
	Token         string `json:"Token"`
	DisplayClaims struct {
		XUI []struct {
			UHS string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

type XSTSResponse struct {
	Token         string `json:"Token"`
	DisplayClaims struct {
		XUI []struct {
			UHS string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

type XSTSErrorResponse struct {
	Identity string `json:"Identity"`
	XErr     int64  `json:"XErr"`
	Message  string `json:"Message"`
	Redirect string `json:"Redirect"`
}

type MinecraftLoginResponse struct {
	Username    string   `json:"username"` // 不是 UUID，别被骗
	Roles       []string `json:"roles"`
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int      `json:"expires_in"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

func GetDeviceCode(mkt string) (*DeviceCodeResponse, error) {
	if mkt == "" {
		mkt = "zh-CN"
	}

	data := url.Values{}
	data.Set("client_id", CLIENT_ID)
	data.Set("scope", "XboxLive.signin offline_access openid email")

	deviceCodeURL := fmt.Sprintf(DEVICE_CODE_URL, TENANT, mkt)

	req, err := http.NewRequest("POST", deviceCodeURL, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	var response DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func PollMSToken(deviceCode string, interval int) (*MSTokenResponse, error) {
	tokenURL := fmt.Sprintf(TOKEN_CHECK_URL, TENANT)

	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	data.Set("client_id", CLIENT_ID)
	data.Set("device_code", deviceCode)
	body := data.Encode()

	fmt.Printf("等待用户登录")

	for {
		req, err := http.NewRequest("POST", tokenURL, strings.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			var tokenResp MSTokenResponse
			if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
				return nil, err
			}
			return &tokenResp, nil
		}

		var errorResp TokenErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, err
		}

		switch errorResp.Error {
		case "authorization_pending":
			time.Sleep(time.Duration(interval) * time.Second)

		case "authorization_declined":
			return nil, fmt.Errorf("用户拒绝授权")

		case "bad_verification_code":
			return nil, fmt.Errorf("device_code 无效")

		case "expired_token":
			return nil, fmt.Errorf("device_code 已过期")

		default:
			return nil, fmt.Errorf("未知错误: %s", errorResp.Error)
		}
	}
}

func MSToken(playerID string) (string, error) {
	deviceCode, err := GetDeviceCode("")
	if err != nil {
		return "", err
	}
	fmt.Println(deviceCode.Message)

	msResp, err := PollMSToken(deviceCode.DeviceCode, deviceCode.Interval)
	if err != nil {
		return "", err
	}

	// 存入结构体
	tokenData := StoredToken{
		PlayerID:     playerID,
		AccessToken:  msResp.AccessToken,
		RefreshToken: msResp.RefreshToken,
		ExpiresAt:    time.Now().Unix() + int64(msResp.ExpiresIn),
	}

	// 保存为二进制
	if err := crypto.SaveToken(playerID+"_ms.token", &tokenData); err != nil {
		return "", err
	}

	fmt.Println("Microsoft 登录成功，Token 已保存")
	return msResp.AccessToken, nil
}

func GetXBLAccess(msAccessToken string) (xblToken string, userHash string, err error) {
	payload := fmt.Sprintf(`{
		"Properties": {
			"AuthMethod": "RPS",
			"SiteName": "user.auth.xboxlive.com",
			"RpsTicket": "d=%s"
		},
		"RelyingParty": "http://auth.xboxlive.com",
		"TokenType": "JWT"
	}`, msAccessToken)

	req, err := http.NewRequest("POST", XBOX_LIVE_AUTH_URL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("XBL 验证错误: %s", body)
	}

	var result XBLAuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", err
	}

	if result.Token == "" || len(result.DisplayClaims.XUI) == 0 {
		return "", "", fmt.Errorf("XBL 错误")
	}

	return result.Token, result.DisplayClaims.XUI[0].UHS, nil
}

func parseXSTSError(code int64) error {
	switch code {
	case 2148916227:
		return fmt.Errorf("Xbox 账号已被封禁")
	case 2148916233:
		return fmt.Errorf("账号未注册 Xbox，需要先登录 minecraft.net 初始化")
	case 2148916235:
		return fmt.Errorf("账号所在国家/地区不支持 Xbox Live")
	case 2148916236, 2148916237:
		return fmt.Errorf("账号需要完成 Xbox 成人验证（韩国地区）")
	case 2148916238:
		return fmt.Errorf("未成年账户，必须加入家庭组（Azure ClientId问题）")
	case 2148916262:
		return fmt.Errorf("Xbox 未知错误（2148916262）")
	default:
		return fmt.Errorf("未知 XSTS 错误码: %d", code)
	}
}

func GetXSTSToken(xblToken string) (tokenXSTS string, userHash string, err error) {
	payload := fmt.Sprintf(`{
		"Properties": {
			"SandboxId": "RETAIL",
			"UserTokens": [
				"%s"
			]
		},
		"RelyingParty": "rp://api.minecraftservices.com/",
		"TokenType": "JWT"
	}`, xblToken)

	req, err := http.NewRequest("POST", XSTS_URL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// ❌ XSTS 失败（最关键）
	if resp.StatusCode == 401 {
		var xerr XSTSErrorResponse
		if err := json.Unmarshal(body, &xerr); err != nil {
			return "", "", fmt.Errorf("XSTS 401 且无法解析错误: %s", body)
		}
		return "", "", parseXSTSError(xerr.XErr)
	}

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("XSTS 请求失败: %s", body)
	}

	// 成功
	var result XSTSResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", err
	}

	if result.Token == "" || len(result.DisplayClaims.XUI) == 0 {
		return "", "", fmt.Errorf("XSTS 返回数据不完整")
	}

	return result.Token, result.DisplayClaims.XUI[0].UHS, nil
}

func GetMinecraftToken(playerID string, xstsToken string, userHash string) (string, error) {
	identityToken := fmt.Sprintf(
		"XBL3.0 x=%s;%s",
		userHash,
		xstsToken,
	)

	payload := fmt.Sprintf(`{
		"identityToken": "%s"
	}`, identityToken)

	req, err := http.NewRequest("POST", LOGIN_WITH_XBOX_URL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 修复：检查读取响应体的错误
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取API响应体失败: %w", err)
	}

	if resp.StatusCode == 403 {
		return "", fmt.Errorf("该Azure ID没有调用 Minecraft API 的权限（403）")
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("MC token 获取失败: %s", body)
	}

	var result MinecraftLoginResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("Minecraft access token 为空")
	}

	// 存储 MC token 到 .session
	mcCache := StoredToken{
		PlayerID:    playerID,
		AccessToken: result.AccessToken,
		ExpiresAt:   time.Now().Add(McTokenExpire).Unix(),
	}

	if err := crypto.SaveToken(playerID+"_mc.token", &mcCache); err != nil {
		return "", fmt.Errorf("保存 MC token 失败: %w", err)
	}

	return result.AccessToken, nil
}

func HasMinecraftOwnership(mcAccessToken string) (bool, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.minecraftservices.com/entitlements/mcstore",
		nil,
	)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+mcAccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("entitlements 请求失败: %s", body)
	}

	// 查找game_minecraft
	var result struct {
		Items []struct {
			Name string `json:"name"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	for _, item := range result.Items {
		if item.Name == "game_minecraft" {
			return true, nil
		}
	}

	// 注意：Game Pass 用户这里也是 false（微软设计）
	return false, nil
}

func RefreshToken(playerID string) error {
	var stored StoredToken
	// 1. 加载现有的 MS token
	err := crypto.LoadToken(playerID+"_ms.token", &stored)
	if err != nil {
		return fmt.Errorf("无法加载 MS Token 进行刷新: %w", err)
	}

	// 2. 请求刷新
	v := url.Values{}
	v.Set("client_id", CLIENT_ID)
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", stored.RefreshToken)

	tokenURL := fmt.Sprintf(TOKEN_CHECK_URL, TENANT)
	resp, err := http.PostForm(tokenURL, v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("刷新失败: %s", string(body))
	}

	var result RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// 3. 更新并保存
	stored.AccessToken = result.AccessToken
	stored.RefreshToken = result.RefreshToken
	stored.ExpiresAt = time.Now().Unix() + int64(result.ExpiresIn)

	return crypto.SaveToken(playerID+"_ms.token", &stored)
}
