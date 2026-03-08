package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	sessionDir      = ".session"
	tokenExpirySkew = 30 * time.Second
)

type MicrosoftTokenCache struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type MinecraftTokenCache struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	ProfileID   string    `json:"profile_id"`
	ProfileName string    `json:"profile_name"`
}

type TokenCache struct {
	PlayerID  string              `json:"player_id"`
	UpdatedAt time.Time           `json:"updated_at"`
	Microsoft MicrosoftTokenCache `json:"microsoft"`
	Minecraft MinecraftTokenCache `json:"minecraft"`
}

func New(playerID string) *TokenCache {
	return &TokenCache{PlayerID: strings.TrimSpace(playerID)}
}

func Path(playerID string) string {
	return filepath.Join(sessionDir, sanitizePlayerID(playerID)+".json")
}

func Load(playerID string) (*TokenCache, error) {
	path := Path(playerID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return New(playerID), nil
		}
		return nil, fmt.Errorf("读取 token 缓存失败: %w", err)
	}

	cache := New(playerID)
	if len(data) == 0 {
		return cache, nil
	}
	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("解析 token 缓存失败: %w", err)
	}
	if strings.TrimSpace(cache.PlayerID) == "" {
		cache.PlayerID = strings.TrimSpace(playerID)
	}
	return cache, nil
}

func Save(playerID string, cache *TokenCache) error {
	if cache == nil {
		return fmt.Errorf("token 缓存为空")
	}
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		return fmt.Errorf("创建 token 缓存目录失败: %w", err)
	}

	cache.PlayerID = strings.TrimSpace(playerID)
	cache.UpdatedAt = time.Now().UTC()

	payload, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 token 缓存失败: %w", err)
	}

	path := Path(playerID)
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o600); err != nil {
		return fmt.Errorf("写入临时 token 缓存失败: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(path)
		if err2 := os.Rename(tmpPath, path); err2 != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("替换 token 缓存失败: %w", err2)
		}
	}

	return nil
}

func (c *TokenCache) HasValidMicrosoftAccess(now time.Time) bool {
	if c == nil {
		return false
	}
	return tokenUsable(c.Microsoft.AccessToken, c.Microsoft.ExpiresAt, now)
}

func (c *TokenCache) HasMicrosoftRefreshToken() bool {
	if c == nil {
		return false
	}
	return strings.TrimSpace(c.Microsoft.RefreshToken) != ""
}

func (c *TokenCache) HasValidMinecraftToken(now time.Time) bool {
	if c == nil {
		return false
	}
	if strings.TrimSpace(c.Minecraft.ProfileID) == "" || strings.TrimSpace(c.Minecraft.ProfileName) == "" {
		return false
	}
	return tokenUsable(c.Minecraft.AccessToken, c.Minecraft.ExpiresAt, now)
}

func tokenUsable(token string, expiresAt time.Time, now time.Time) bool {
	if strings.TrimSpace(token) == "" || expiresAt.IsZero() {
		return false
	}
	return now.Add(tokenExpirySkew).Before(expiresAt)
}

func sanitizePlayerID(playerID string) string {
	raw := strings.TrimSpace(playerID)
	if raw == "" {
		return "default"
	}

	var b strings.Builder
	b.Grow(len(raw))
	for _, r := range raw {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}

	sanitized := strings.Trim(b.String(), "_")
	if sanitized == "" {
		return "default"
	}
	return sanitized
}
