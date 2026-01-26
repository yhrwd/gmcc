package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 定义玩家信息的结构体（与返回 JSON 字段一一对应）
// 嵌套结构体对应 skins 和 capes 数组中的对象
type MinecraftProfile struct {
	ID    string `json:"id"`    // 玩家 UUID
	Name  string `json:"name"`  // Minecraft 用户名
	Skins []Skin `json:"skins"` // 皮肤列表
	Capes []Cape `json:"capes"` // 披风列表
}

// 皮肤信息结构体
type Skin struct {
	ID      string `json:"id"`      // 皮肤 ID
	State   string `json:"state"`   // 状态（ACTIVE 表示活跃）
	URL     string `json:"url"`     // 皮肤纹理 URL
	Variant string `json:"variant"` // 皮肤类型（CLASSIC 经典款）
	Alias   string `json:"alias"`   // 别名（如 STEVE）
}

// 披风信息结构体
type Cape struct {
	ID          string `json:"id"`          // 披风 ID
	State       string `json:"state"`       // 状态（ACTIVE 表示活跃）
	URL         string `json:"url"`         // 披风纹理 URL
	Alias       string `json:"alias"`       // 别名（如 Migrator）
}

// 错误响应结构体（请求失败时解析返回的错误信息）
type ErrorResponse struct {
	Path         string `json:"path"`         // 请求路径
	Error        string `json:"error"`        // 错误码（如 NOT_FOUND）
	ErrorMessage string `json:"errorMessage"` // 错误描述
}

// GetMinecraftProfile 发送请求获取玩家信息
// 参数：minecraftToken - 之前获取的 Minecraft 访问令牌（Bearer Token）
// 返回：*MinecraftProfile - 玩家信息结构体指针；error - 错误信息（请求失败/解析失败时返回）
func GetPlayerInfo(mcToken string) (*MinecraftProfile, error) {
	reqURL := "https://api.minecraftservices.com/minecraft/profile"
	client := &http.Client{}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败：%v", err)
	}

	// Authorization: Bearer <令牌>
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mcToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败：%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("请求失败，状态码：%d，解析错误响应失败：%v", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("请求失败：错误码=%s，描述=%s", errResp.Error, errResp.ErrorMessage)
	}

	var profile MinecraftProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("解析玩家信息失败：%v", err)
	}

	return &profile, nil
}