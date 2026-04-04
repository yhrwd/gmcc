package session

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"gmcc/internal/auth/microsoft"
	"gmcc/internal/auth/minecraft"
)

type MicrosoftProvider interface {
	BeginDeviceLogin(ctx context.Context, accountID string) (DeviceLoginInfo, error)
	PollDeviceLogin(ctx context.Context, accountID string, info DeviceLoginInfo) (MicrosoftTokenCache, error)
	RefreshMicrosoft(ctx context.Context, refreshToken string) (MicrosoftTokenCache, error)
	GetXSTSFromMicrosoft(ctx context.Context, accessToken string) (XSTSClaims, error)
}

type MinecraftProvider interface {
	ExchangeMinecraftToken(ctx context.Context, xsts XSTSClaims) (MinecraftTokenCache, error)
	VerifyOwnership(ctx context.Context, accessToken string) error
	GetProfile(ctx context.Context, accessToken string) (MinecraftProfileData, error)
}

type ProviderSet struct {
	Microsoft MicrosoftProvider
	Minecraft MinecraftProvider
}

func NewLiveProviderSet() ProviderSet {
	return ProviderSet{
		Microsoft: NewLiveMicrosoftProvider(),
		Minecraft: LiveMinecraftProvider{},
	}
}

func (p ProviderSet) BeginDeviceLogin(ctx context.Context, accountID string) (DeviceLoginInfo, error) {
	if p.Microsoft == nil {
		return DeviceLoginInfo{}, fmt.Errorf("microsoft provider is nil: %w", ErrProviderUnavailable)
	}
	return p.Microsoft.BeginDeviceLogin(ctx, accountID)
}

func (p ProviderSet) PollDeviceLogin(ctx context.Context, accountID string, info DeviceLoginInfo) (MicrosoftTokenCache, error) {
	if p.Microsoft == nil {
		return MicrosoftTokenCache{}, fmt.Errorf("microsoft provider is nil: %w", ErrProviderUnavailable)
	}
	return p.Microsoft.PollDeviceLogin(ctx, accountID, info)
}

func (p ProviderSet) RefreshMicrosoft(ctx context.Context, refreshToken string) (MicrosoftTokenCache, error) {
	if p.Microsoft == nil {
		return MicrosoftTokenCache{}, fmt.Errorf("microsoft provider is nil: %w", ErrProviderUnavailable)
	}
	return p.Microsoft.RefreshMicrosoft(ctx, refreshToken)
}

func (p ProviderSet) GetXSTSFromMicrosoft(ctx context.Context, accessToken string) (XSTSClaims, error) {
	if p.Microsoft == nil {
		return XSTSClaims{}, fmt.Errorf("microsoft provider is nil: %w", ErrProviderUnavailable)
	}
	return p.Microsoft.GetXSTSFromMicrosoft(ctx, accessToken)
}

func (p ProviderSet) ExchangeMinecraftToken(ctx context.Context, xsts XSTSClaims) (MinecraftTokenCache, error) {
	if p.Minecraft == nil {
		return MinecraftTokenCache{}, fmt.Errorf("minecraft provider is nil: %w", ErrProviderUnavailable)
	}
	return p.Minecraft.ExchangeMinecraftToken(ctx, xsts)
}

func (p ProviderSet) VerifyOwnership(ctx context.Context, accessToken string) error {
	if p.Minecraft == nil {
		return fmt.Errorf("minecraft provider is nil: %w", ErrProviderUnavailable)
	}
	return p.Minecraft.VerifyOwnership(ctx, accessToken)
}

func (p ProviderSet) GetProfile(ctx context.Context, accessToken string) (MinecraftProfileData, error) {
	if p.Minecraft == nil {
		return MinecraftProfileData{}, fmt.Errorf("minecraft provider is nil: %w", ErrProviderUnavailable)
	}
	return p.Minecraft.GetProfile(ctx, accessToken)
}

type LiveMicrosoftProvider struct {
	mu          sync.Mutex
	deviceCodes map[string]*microsoft.DeviceCodeResponse
}

func NewLiveMicrosoftProvider() *LiveMicrosoftProvider {
	return &LiveMicrosoftProvider{deviceCodes: map[string]*microsoft.DeviceCodeResponse{}}
}

func (p *LiveMicrosoftProvider) BeginDeviceLogin(ctx context.Context, accountID string) (DeviceLoginInfo, error) {
	if err := ctx.Err(); err != nil {
		return DeviceLoginInfo{}, err
	}

	resp, err := microsoft.BeginDeviceCodeLogin()
	if err != nil {
		return DeviceLoginInfo{}, err
	}

	info := DeviceLoginInfoFromMicrosoftResponse(accountID, resp)

	p.mu.Lock()
	p.deviceCodes[strings.TrimSpace(accountID)] = resp
	p.mu.Unlock()

	return info, nil
}

func (p *LiveMicrosoftProvider) PollDeviceLogin(ctx context.Context, accountID string, _ DeviceLoginInfo) (MicrosoftTokenCache, error) {
	if err := ctx.Err(); err != nil {
		return MicrosoftTokenCache{}, err
	}

	trimmedAccountID := strings.TrimSpace(accountID)

	p.mu.Lock()
	deviceCode := p.deviceCodes[trimmedAccountID]
	p.mu.Unlock()

	if deviceCode == nil {
		return MicrosoftTokenCache{}, ErrDeviceLoginRequired
	}

	resp, err := microsoft.PollDeviceCodeLogin(deviceCode)
	if err != nil {
		return MicrosoftTokenCache{}, err
	}
	return MicrosoftCacheFromToken(resp), nil
}

func (p *LiveMicrosoftProvider) RefreshMicrosoft(ctx context.Context, refreshToken string) (MicrosoftTokenCache, error) {
	if err := ctx.Err(); err != nil {
		return MicrosoftTokenCache{}, err
	}

	resp, err := microsoft.RefreshMicrosoftToken(refreshToken)
	if err != nil {
		return MicrosoftTokenCache{}, err
	}
	return MicrosoftCacheFromToken(resp), nil
}

func (p *LiveMicrosoftProvider) GetXSTSFromMicrosoft(ctx context.Context, accessToken string) (XSTSClaims, error) {
	if err := ctx.Err(); err != nil {
		return XSTSClaims{}, err
	}

	resp, err := microsoft.GetXSTSTokenFromAccessToken(accessToken)
	if err != nil {
		return XSTSClaims{}, err
	}
	return XSTSClaimsFromMicrosoftResponse(resp)
}

type LiveMinecraftProvider struct{}

func (LiveMinecraftProvider) ExchangeMinecraftToken(ctx context.Context, xsts XSTSClaims) (MinecraftTokenCache, error) {
	if err := ctx.Err(); err != nil {
		return MinecraftTokenCache{}, err
	}

	resp, err := minecraft.GetMinecraftTokenFromClaims(xsts.Token, xsts.UserHash)
	if err != nil {
		return MinecraftTokenCache{}, err
	}

	return MinecraftCacheFromToken(resp), nil
}

func (LiveMinecraftProvider) VerifyOwnership(ctx context.Context, accessToken string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return minecraft.VerifyGameOwnership(accessToken)
}

func (LiveMinecraftProvider) GetProfile(ctx context.Context, accessToken string) (MinecraftProfileData, error) {
	if err := ctx.Err(); err != nil {
		return MinecraftProfileData{}, err
	}

	profile, err := minecraft.GetProfile(accessToken)
	if err != nil {
		return MinecraftProfileData{}, err
	}
	return MinecraftProfileDataFromResponse(profile)
}

func DeviceLoginInfoFromMicrosoftResponse(accountID string, resp *microsoft.DeviceCodeResponse) DeviceLoginInfo {
	if resp == nil {
		return DeviceLoginInfo{AccountID: strings.TrimSpace(accountID)}
	}

	pollInterval := time.Duration(resp.Interval) * time.Second
	if pollInterval <= 0 {
		pollInterval = 5 * time.Second
	}

	return DeviceLoginInfo{
		AccountID:       strings.TrimSpace(accountID),
		VerificationURI: strings.TrimSpace(resp.VerificationURI),
		UserCode:        strings.TrimSpace(resp.UserCode),
		ExpiresAt:       time.Now().UTC().Add(time.Duration(resp.ExpiresIn) * time.Second),
		PollInterval:    pollInterval,
	}
}

func MicrosoftCacheFromToken(resp *microsoft.TokenResponse) MicrosoftTokenCache {
	if resp == nil {
		return MicrosoftTokenCache{}
	}

	return MicrosoftTokenCache{
		AccessToken:  strings.TrimSpace(resp.AccessToken),
		RefreshToken: strings.TrimSpace(resp.RefreshToken),
		ExpiresAt:    time.Now().UTC().Add(time.Duration(resp.ExpiresIn) * time.Second),
	}
}

func XSTSClaimsFromMicrosoftResponse(resp *microsoft.XSTSResponse) (XSTSClaims, error) {
	if resp == nil {
		return XSTSClaims{}, fmt.Errorf("xsts response is nil")
	}
	if strings.TrimSpace(resp.Token) == "" {
		return XSTSClaims{}, fmt.Errorf("xsts token is empty")
	}
	if len(resp.DisplayClaims.Xui) == 0 || strings.TrimSpace(resp.DisplayClaims.Xui[0].Uhs) == "" {
		return XSTSClaims{}, fmt.Errorf("xsts user hash is empty")
	}

	return XSTSClaims{Token: strings.TrimSpace(resp.Token), UserHash: strings.TrimSpace(resp.DisplayClaims.Xui[0].Uhs)}, nil
}

func MinecraftCacheFromToken(resp *minecraft.MinecraftTokenResponse) MinecraftTokenCache {
	if resp == nil {
		return MinecraftTokenCache{}
	}
	return MinecraftTokenCache{
		AccessToken: strings.TrimSpace(resp.AccessToken),
		ExpiresAt:   time.Now().UTC().Add(time.Duration(resp.ExpiresIn) * time.Second),
	}
}

func MinecraftProfileDataFromResponse(profile *minecraft.MinecraftProfile) (MinecraftProfileData, error) {
	if profile == nil {
		return MinecraftProfileData{}, fmt.Errorf("minecraft profile is nil")
	}

	id := strings.TrimSpace(profile.ID)
	name := strings.TrimSpace(profile.Name)
	if id == "" || name == "" {
		return MinecraftProfileData{}, ErrProfileInvalid
	}

	return MinecraftProfileData{ID: id, Name: name}, nil
}
