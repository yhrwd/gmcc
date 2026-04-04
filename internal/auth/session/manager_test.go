package session

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestAuthManager_GetSessionSingleFlightRefresh(t *testing.T) {
	provider := newFakeProvider()
	provider.refreshResult = fakeRefreshResult{
		session: AuthSession{
			AccountID:            "acc-main",
			MinecraftAccessToken: "mc-1",
			ProfileID:            "uuid",
			ProfileName:          "Steve",
		},
	}
	store := newFakeTokenStore()
	mgr := NewAuthManager(store, provider)

	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = mgr.GetSession(context.Background(), "acc-main")
		}()
	}
	wg.Wait()

	if provider.refreshCalls != 1 {
		t.Fatalf("want 1 refresh call, got %d", provider.refreshCalls)
	}
}

func TestAuthManager_DeviceLoginStatusLifecycle(t *testing.T) {
	provider := newFakeProvider()
	provider.deviceLoginInfo = DeviceLoginInfo{
		AccountID:       "acc-main",
		UserCode:        "ABCD",
		VerificationURI: "https://microsoft.com/devicelogin",
	}
	store := newFakeTokenStore()
	mgr := NewAuthManager(store, provider)

	if _, err := mgr.BeginDeviceLogin(context.Background(), "acc-main"); err != nil {
		t.Fatalf("begin device login: %v", err)
	}

	status, _, err := mgr.GetDeviceLoginStatus("acc-main")
	if err != nil {
		t.Fatalf("get login status: %v", err)
	}
	if status != DeviceLoginStatusPending {
		t.Fatalf("want pending, got %s", status)
	}
}

func TestAuthManager_ProviderUnavailableNormalizesToAuthFailed(t *testing.T) {
	provider := newFakeProvider()
	provider.refreshErr = errors.New("upstream unavailable")
	mgr := NewAuthManager(newFakeTokenStore(), provider)

	_, err := mgr.Refresh(context.Background(), "acc-main")
	if !errors.Is(err, ErrProviderUnavailable) {
		t.Fatalf("want ErrProviderUnavailable, got %v", err)
	}
}

type fakeTokenStore struct{}

func newFakeTokenStore() *fakeTokenStore {
	return &fakeTokenStore{}
}

type fakeProvider struct {
	refreshCalls    int
	refreshResult   fakeRefreshResult
	refreshErr      error
	deviceLoginInfo DeviceLoginInfo
}

type fakeRefreshResult struct {
	session AuthSession
}

func newFakeProvider() *fakeProvider {
	return &fakeProvider{}
}

func (f *fakeProvider) BeginDeviceLogin(_ context.Context, _ string) (DeviceLoginInfo, error) {
	return f.deviceLoginInfo, nil
}

func (f *fakeProvider) PollDeviceLogin(_ context.Context, _ string, _ DeviceLoginInfo) (MicrosoftTokenCache, error) {
	return MicrosoftTokenCache{}, nil
}

func (f *fakeProvider) RefreshMicrosoft(_ context.Context, _ string) (MicrosoftTokenCache, error) {
	f.refreshCalls++
	if f.refreshErr != nil {
		return MicrosoftTokenCache{}, f.refreshErr
	}
	return MicrosoftTokenCache{ExpiresAt: time.Now().Add(10 * time.Minute)}, nil
}

func (f *fakeProvider) GetXSTSFromMicrosoft(_ context.Context, _ string) (XSTSClaims, error) {
	return XSTSClaims{}, nil
}

func (f *fakeProvider) ExchangeMinecraftToken(_ context.Context, _ XSTSClaims) (MinecraftTokenCache, error) {
	return MinecraftTokenCache{
		AccessToken: f.refreshResult.session.MinecraftAccessToken,
		ProfileID:   f.refreshResult.session.ProfileID,
		ProfileName: f.refreshResult.session.ProfileName,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}, nil
}

func (f *fakeProvider) VerifyOwnership(_ context.Context, _ string) error {
	return nil
}

func (f *fakeProvider) GetProfile(_ context.Context, _ string) (MinecraftProfileData, error) {
	return MinecraftProfileData{
		ID:   f.refreshResult.session.ProfileID,
		Name: f.refreshResult.session.ProfileName,
	}, nil
}
