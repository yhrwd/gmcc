package minecraft

import (
	"fmt"
	"strings"
	"time"

	"gmcc/internal/auth/microsoft"
	"gmcc/pkg/httpx"
)

const (
	MINECRAFT_LOGIN_URL        string = "https://api.minecraftservices.com/authentication/login_with_xbox"
	MINECRAFT_ENTITLEMENT_URL  string = "https://api.minecraftservices.com/entitlements/mcstore"
	MINECRAFT_PROFILE_URL      string = "https://api.minecraftservices.com/minecraft/profile"
	MINECRAFT_SESSION_JOIN_URL string = "https://sessionserver.mojang.com/session/minecraft/join"
	MINECRAFT_CERT_URL         string = "https://api.minecraftservices.com/player/certificates"
)

/* -------------------- Minecraft Token -------------------- */

type MinecraftTokenRequest struct {
	IdentityToken string `json:"identityToken"`
}

type MinecraftTokenResponse struct {
	Username    string        `json:"username"`
	Roles       []interface{} `json:"roles"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   int           `json:"expires_in"`
}

type MinecraftEntitlementItem struct {
	Name      string `json:"name"`
	Signature string `json:"signature"`
}

type MinecraftEntitlementResponse struct {
	Items     []MinecraftEntitlementItem `json:"items"`
	Signature string                     `json:"signature"`
	KeyId     string                     `json:"keyId"`
}

type MinecraftProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PlayerCertificatesResponse struct {
	KeyPair struct {
		PrivateKey string `json:"privateKey"`
		PublicKey  string `json:"publicKey"`
	} `json:"keyPair"`
	PublicKeySignature   string    `json:"publicKeySignature"`
	PublicKeySignatureV2 string    `json:"publicKeySignatureV2"`
	ExpiresAt            time.Time `json:"expiresAt"`
	RefreshedAfter       time.Time `json:"refreshedAfter"`
}

// GetMinecraftToken 使用 XSTS 令牌获取 Minecraft 访问令牌
func GetMinecraftToken(xstsResp *microsoft.XSTSResponse) (*MinecraftTokenResponse, error) {
	if xstsResp == nil || len(xstsResp.DisplayClaims.Xui) == 0 {
		return nil, fmt.Errorf("无效的 XSTS 响应")
	}

	identityToken := fmt.Sprintf("XBL3.0 x=%s;%s", xstsResp.DisplayClaims.Xui[0].Uhs, xstsResp.Token)

	// 构造请求
	req := MinecraftTokenRequest{
		IdentityToken: identityToken,
	}

	// 发送请求
	var resp MinecraftTokenResponse
	_, err := httpx.PostJSON(MINECRAFT_LOGIN_URL, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("获取 Minecraft Token 失败: %w", err)
	}

	return &resp, nil
}

func GetMinecraftTokenFromClaims(token, userHash string) (*MinecraftTokenResponse, error) {
	trimmedToken := strings.TrimSpace(token)
	trimmedUserHash := strings.TrimSpace(userHash)
	if trimmedToken == "" || trimmedUserHash == "" {
		return nil, fmt.Errorf("无效的 XSTS 凭据")
	}

	return GetMinecraftToken(&microsoft.XSTSResponse{
		Token: trimmedToken,
		DisplayClaims: struct {
			Xui []struct {
				Uhs string `json:"uhs"`
			} `json:"xui"`
		}{
			Xui: []struct {
				Uhs string `json:"uhs"`
			}{{Uhs: trimmedUserHash}},
		},
	})
}

func VerifyGameOwnership(minecraftToken string) error {
	var resp MinecraftEntitlementResponse
	_, err := httpx.GetWithAuthHeader(
		MINECRAFT_ENTITLEMENT_URL,
		fmt.Sprintf("Bearer %s", minecraftToken),
		&resp,
	)
	if err != nil {
		return fmt.Errorf("验证游戏所有权失败: %w", err)
	}

	if len(resp.Items) == 0 {
		return fmt.Errorf("账户未购买 Minecraft 游戏")
	}
	return nil
}

func GetProfile(minecraftToken string) (*MinecraftProfile, error) {
	var resp MinecraftProfile
	_, err := httpx.GetWithAuthHeader(
		MINECRAFT_PROFILE_URL,
		fmt.Sprintf("Bearer %s", minecraftToken),
		&resp,
	)
	if err != nil {
		return nil, fmt.Errorf("获取 Minecraft Profile 失败: %w", err)
	}
	if strings.TrimSpace(resp.ID) == "" || strings.TrimSpace(resp.Name) == "" {
		return nil, fmt.Errorf("profile 数据不完整")
	}
	return &resp, nil
}

type joinServerRequest struct {
	AccessToken     string `json:"accessToken"`
	SelectedProfile string `json:"selectedProfile"`
	ServerID        string `json:"serverId"`
}

func JoinServer(minecraftAccessToken, selectedProfile, serverID string) error {
	req := joinServerRequest{
		AccessToken:     minecraftAccessToken,
		SelectedProfile: strings.ReplaceAll(selectedProfile, "-", ""),
		ServerID:        serverID,
	}
	_, err := httpx.PostJSON(MINECRAFT_SESSION_JOIN_URL, req, nil)
	if err != nil {
		return fmt.Errorf("join session 失败: %w", err)
	}
	return nil
}

func GetPlayerCertificates(minecraftToken string) (*PlayerCertificatesResponse, error) {
	var resp PlayerCertificatesResponse
	_, err := httpx.PostJSONWithAuthHeader(
		MINECRAFT_CERT_URL,
		fmt.Sprintf("Bearer %s", minecraftToken),
		map[string]interface{}{},
		&resp,
	)
	if err != nil {
		return nil, fmt.Errorf("获取 player certificates 失败: %w", err)
	}
	if strings.TrimSpace(resp.KeyPair.PublicKey) == "" {
		return nil, fmt.Errorf("player certificates 缺少 public key")
	}
	return &resp, nil
}
