package session

import (
	"errors"
	"strings"
	"time"
)

var ErrAccountNotFound = errors.New("account auth record not found")

type AccountAuthStatus string

const (
	AccountAuthStatusLoggedIn    AccountAuthStatus = "logged_in"
	AccountAuthStatusNotLoggedIn AccountAuthStatus = "not_logged_in"
	AccountAuthStatusAuthInvalid AccountAuthStatus = "auth_invalid"
)

type MicrosoftCredential = MicrosoftTokenCache

type MinecraftSessionState = MinecraftTokenCache

type AccountAuthRecord struct {
	AccountID     string                `json:"account_id"`
	UpdatedAt     time.Time             `json:"updated_at"`
	LastAuthError string                `json:"last_auth_error,omitempty"`
	Microsoft     MicrosoftCredential   `json:"microsoft"`
	Minecraft     MinecraftSessionState `json:"minecraft"`
}

type AccountRecordRepository interface {
	GetAccount(accountID string) (*AccountAuthRecord, error)
	PutAccount(record *AccountAuthRecord) error
	DeleteAccount(accountID string) error
}

func NewAccountAuthRecord(accountID string) *AccountAuthRecord {
	return &AccountAuthRecord{AccountID: strings.TrimSpace(accountID)}
}

func (r *AccountAuthRecord) HasValidMicrosoftAccess(now time.Time) bool {
	if r == nil {
		return false
	}
	return tokenUsable(r.Microsoft.AccessToken, r.Microsoft.ExpiresAt, now)
}

func (r *AccountAuthRecord) HasMicrosoftRefreshToken() bool {
	if r == nil {
		return false
	}
	return strings.TrimSpace(r.Microsoft.RefreshToken) != ""
}

func (r *AccountAuthRecord) HasValidMinecraftSession(now time.Time) bool {
	if r == nil {
		return false
	}
	if strings.TrimSpace(r.Minecraft.ProfileID) == "" || strings.TrimSpace(r.Minecraft.ProfileName) == "" {
		return false
	}
	return tokenUsable(r.Minecraft.AccessToken, r.Minecraft.ExpiresAt, now)
}

func (r *AccountAuthRecord) ToAuthSession(source AuthSource) AuthSession {
	if r == nil {
		return AuthSession{}
	}
	return AuthSession{
		AccountID:            strings.TrimSpace(r.AccountID),
		MinecraftAccessToken: strings.TrimSpace(r.Minecraft.AccessToken),
		ProfileID:            strings.TrimSpace(r.Minecraft.ProfileID),
		ProfileName:          strings.TrimSpace(r.Minecraft.ProfileName),
		MicrosoftExpiresAt:   r.Microsoft.ExpiresAt,
		MinecraftExpiresAt:   r.Minecraft.ExpiresAt,
		Source:               source,
	}
}

func classifyAuthRecord(record *AccountAuthRecord) AccountAuthStatus {
	if record == nil {
		return AccountAuthStatusNotLoggedIn
	}
	if record.HasMicrosoftRefreshToken() {
		return AccountAuthStatusLoggedIn
	}
	if strings.TrimSpace(record.LastAuthError) != "" {
		return AccountAuthStatusAuthInvalid
	}
	return AccountAuthStatusNotLoggedIn
}
