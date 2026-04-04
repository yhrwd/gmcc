package session

import "context"

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
