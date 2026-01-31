package auth

import (
	"fmt"
	"gmcc/pkg/network"
)

const (
	MINECRAFT_LOGIN_URL       string = "https://api.minecraftservices.com/authentication/login_with_xbox"
	MINECRAFT_ENTITLEMENT_URL string = "https://api.minecraftservices.com/entitlements/mcstore"
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

// GetMinecraftToken 使用 XSTS 令牌获取 Minecraft 访问令牌
func GetMinecraftToken(xstsResp *XSTSResponse) (*MinecraftTokenResponse, error) {
	identityToken := fmt.Sprintf("XBL3.0 x=%s;%s", xstsResp.DisplayClaims.Xui[0].Uhs, xstsResp.Token)

	// 构造请求
	req := MinecraftTokenRequest{
		IdentityToken: identityToken,
	}

	// 发送请求
	var resp MinecraftTokenResponse
	_, err := network.PostJSON(MINECRAFT_LOGIN_URL, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("获取 Minecraft Token 失败: %w", err)
	}

	return &resp, nil
}

func VerifyGameOwnership(minecraftToken string) error {
	var resp MinecraftEntitlementResponse
	_, err := network.GetWithAuthHeader(
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