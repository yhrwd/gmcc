package session

import (
	"errors"
	"strings"
	"time"

	"gmcc/internal/constants"
)

type AuthSource string

const (
	AuthSourceCache       AuthSource = "cache"
	AuthSourceRefresh     AuthSource = "refresh"
	AuthSourceDeviceLogin AuthSource = "device_login"
)

type DeviceLoginStatus string

const (
	DeviceLoginStatusPending   DeviceLoginStatus = "pending"
	DeviceLoginStatusSucceeded DeviceLoginStatus = "succeeded"
	DeviceLoginStatusExpired   DeviceLoginStatus = "expired"
	DeviceLoginStatusCancelled DeviceLoginStatus = "cancelled"
	DeviceLoginStatusFailed    DeviceLoginStatus = "failed"
)

type AuthSession struct {
	AccountID            string
	MinecraftAccessToken string
	ProfileID            string
	ProfileName          string
	MicrosoftExpiresAt   time.Time
	MinecraftExpiresAt   time.Time
	Source               AuthSource
}

type DeviceLoginInfo struct {
	AccountID       string
	VerificationURI string
	UserCode        string
	ExpiresAt       time.Time
	PollInterval    time.Duration
}

type MicrosoftTokenCache struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type XSTSClaims struct {
	Token    string
	UserHash string
}

type MinecraftTokenCache struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	ProfileID   string    `json:"profile_id"`
	ProfileName string    `json:"profile_name"`
}

type MinecraftProfileData struct {
	ID   string
	Name string
}

type AccountProfile struct {
	ProfileID   string
	ProfileName string
}

type TokenCache struct {
	AccountID     string              `json:"account_id"`
	UpdatedAt     time.Time           `json:"updated_at"`
	LastAuthError string              `json:"last_auth_error,omitempty"`
	Microsoft     MicrosoftTokenCache `json:"microsoft"`
	Minecraft     MinecraftTokenCache `json:"minecraft"`
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

func (c *TokenCache) ToAuthSession(source AuthSource) AuthSession {
	if c == nil {
		return AuthSession{}
	}
	return AuthSession{
		AccountID:            strings.TrimSpace(c.AccountID),
		MinecraftAccessToken: strings.TrimSpace(c.Minecraft.AccessToken),
		ProfileID:            strings.TrimSpace(c.Minecraft.ProfileID),
		ProfileName:          strings.TrimSpace(c.Minecraft.ProfileName),
		MicrosoftExpiresAt:   c.Microsoft.ExpiresAt,
		MinecraftExpiresAt:   c.Minecraft.ExpiresAt,
		Source:               source,
	}
}

func tokenUsable(token string, expiresAt time.Time, now time.Time) bool {
	if strings.TrimSpace(token) == "" || expiresAt.IsZero() {
		return false
	}
	return now.Add(constants.TokenExpirySkew).Before(expiresAt)
}

var (
	ErrDeviceLoginRequired = errors.New("device login required")
	ErrRefreshTokenInvalid = errors.New("refresh token invalid")
	ErrRefreshUpstream     = errors.New("refresh upstream failed")
	ErrProviderUnavailable = errors.New("provider unavailable")
	ErrOwnershipFailed     = errors.New("ownership failed")
	ErrProfileInvalid      = errors.New("profile invalid")
	ErrXSTSDenied          = errors.New("xsts denied")
)
