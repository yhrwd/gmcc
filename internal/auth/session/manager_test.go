package session

import (
	"context"
	"errors"
	"fmt"
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
		ExpiresAt:       time.Now().UTC().Add(2 * time.Second),
		PollInterval:    10 * time.Millisecond,
	}
	provider.pollResults = []fakePollResult{{err: errAuthorizationPending}}
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

func TestAuthManager_DeviceLoginTransitions(t *testing.T) {
	tests := []struct {
		name           string
		setupProvider  func(*fakeProvider)
		action         func(*testing.T, *AuthManager)
		wantStatus     DeviceLoginStatus
		wantSession    bool
		wantErrIs      error
		wantErrNil     bool
		wantSource     AuthSource
		wantPollAtMost int
	}{
		{
			name: "pending to succeeded",
			setupProvider: func(p *fakeProvider) {
				p.deviceLoginInfo = DeviceLoginInfo{
					AccountID:       "acc-main",
					UserCode:        "ABCD",
					VerificationURI: "https://microsoft.com/devicelogin",
					ExpiresAt:       time.Now().UTC().Add(3 * time.Second),
					PollInterval:    5 * time.Millisecond,
				}
				p.refreshResult = fakeRefreshResult{
					session: AuthSession{
						AccountID:            "acc-main",
						MinecraftAccessToken: "mc-device",
						ProfileID:            "uuid-device",
						ProfileName:          "Steve",
					},
				}
				p.pollResults = []fakePollResult{
					{err: errAuthorizationPending},
					{msCache: MicrosoftTokenCache{AccessToken: "ms-access", RefreshToken: "ms-refresh", ExpiresAt: time.Now().UTC().Add(30 * time.Minute)}},
				}
			},
			wantStatus:     DeviceLoginStatusSucceeded,
			wantSession:    true,
			wantErrNil:     true,
			wantSource:     AuthSourceDeviceLogin,
			wantPollAtMost: -1,
		},
		{
			name: "pending to expired",
			setupProvider: func(p *fakeProvider) {
				p.deviceLoginInfo = DeviceLoginInfo{
					AccountID:       "acc-main",
					UserCode:        "ABCD",
					VerificationURI: "https://microsoft.com/devicelogin",
					ExpiresAt:       time.Now().UTC().Add(-20 * time.Millisecond),
					PollInterval:    5 * time.Millisecond,
				}
			},
			wantStatus:     DeviceLoginStatusExpired,
			wantErrIs:      ErrDeviceLoginRequired,
			wantPollAtMost: 0,
		},
		{
			name: "pending to cancelled",
			setupProvider: func(p *fakeProvider) {
				p.deviceLoginInfo = DeviceLoginInfo{
					AccountID:       "acc-main",
					UserCode:        "ABCD",
					VerificationURI: "https://microsoft.com/devicelogin",
					ExpiresAt:       time.Now().UTC().Add(3 * time.Second),
					PollInterval:    50 * time.Millisecond,
				}
				p.pollResults = []fakePollResult{{err: errAuthorizationPending}}
			},
			action: func(t *testing.T, mgr *AuthManager) {
				t.Helper()
				if err := mgr.CancelDeviceLogin("acc-main"); err != nil {
					t.Fatalf("cancel device login: %v", err)
				}
			},
			wantStatus:     DeviceLoginStatusCancelled,
			wantErrIs:      ErrDeviceLoginRequired,
			wantPollAtMost: -1,
		},
		{
			name: "pending to failed",
			setupProvider: func(p *fakeProvider) {
				p.deviceLoginInfo = DeviceLoginInfo{
					AccountID:       "acc-main",
					UserCode:        "ABCD",
					VerificationURI: "https://microsoft.com/devicelogin",
					ExpiresAt:       time.Now().UTC().Add(3 * time.Second),
					PollInterval:    5 * time.Millisecond,
				}
				p.pollResults = []fakePollResult{{err: errors.New("upstream unavailable")}}
			},
			wantStatus:     DeviceLoginStatusFailed,
			wantErrIs:      ErrProviderUnavailable,
			wantPollAtMost: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := newFakeProvider()
			tt.setupProvider(provider)
			mgr := NewAuthManager(newFakeTokenStore(), provider)

			if _, err := mgr.BeginDeviceLogin(context.Background(), "acc-main"); err != nil {
				t.Fatalf("begin device login: %v", err)
			}

			if tt.action != nil {
				tt.action(t, mgr)
			}

			status, session, err := waitForDeviceLoginFinalState(t, mgr, "acc-main", 2*time.Second)
			if status != tt.wantStatus {
				t.Fatalf("want status %s, got %s", tt.wantStatus, status)
			}

			if tt.wantSession && session == nil {
				t.Fatalf("want session, got nil")
			}
			if !tt.wantSession && session != nil {
				t.Fatalf("want nil session, got %+v", *session)
			}
			if tt.wantSession && session.Source != tt.wantSource {
				t.Fatalf("want source %s, got %s", tt.wantSource, session.Source)
			}

			if tt.wantErrNil && err != nil {
				t.Fatalf("want nil error, got %v", err)
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Fatalf("want error %v, got %v", tt.wantErrIs, err)
			}

			if tt.wantPollAtMost >= 0 {
				if got := provider.getPollCalls(); got > tt.wantPollAtMost {
					t.Fatalf("want poll calls <= %d, got %d", tt.wantPollAtMost, got)
				}
			}
		})
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

func newFakeTokenStore() *fakeTokenStore {
	return &fakeTokenStore{caches: map[string]*TokenCache{}}
}

type fakeTokenStore struct {
	mu     sync.Mutex
	caches map[string]*TokenCache
}

func (f *fakeTokenStore) Load(accountID string) (*TokenCache, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if cache, ok := f.caches[accountID]; ok {
		copy := *cache
		return &copy, nil
	}
	return &TokenCache{
		AccountID: accountID,
		Microsoft: MicrosoftTokenCache{
			RefreshToken: "refresh-token",
		},
	}, nil
}

func (f *fakeTokenStore) Save(accountID string, cache *TokenCache) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	copy := *cache
	f.caches[accountID] = &copy
	return nil
}

type fakeProvider struct {
	mu              sync.Mutex
	refreshCalls    int
	refreshResult   fakeRefreshResult
	refreshErr      error
	deviceLoginInfo DeviceLoginInfo
	pollResults     []fakePollResult
	pollCalls       int
}

type fakeRefreshResult struct {
	session AuthSession
}

type fakePollResult struct {
	msCache MicrosoftTokenCache
	err     error
}

var errAuthorizationPending = errors.New("authorization_pending")

func newFakeProvider() *fakeProvider {
	return &fakeProvider{}
}

func (f *fakeProvider) BeginDeviceLogin(_ context.Context, _ string) (DeviceLoginInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.deviceLoginInfo, nil
}

func (f *fakeProvider) PollDeviceLogin(_ context.Context, _ string, _ DeviceLoginInfo) (MicrosoftTokenCache, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.pollCalls++
	if len(f.pollResults) == 0 {
		return MicrosoftTokenCache{}, errAuthorizationPending
	}

	result := f.pollResults[0]
	f.pollResults = f.pollResults[1:]
	return result.msCache, result.err
}

func (f *fakeProvider) RefreshMicrosoft(_ context.Context, _ string) (MicrosoftTokenCache, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

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
	f.mu.Lock()
	defer f.mu.Unlock()

	return MinecraftProfileData{
		ID:   f.refreshResult.session.ProfileID,
		Name: f.refreshResult.session.ProfileName,
	}, nil
}

func (f *fakeProvider) getPollCalls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.pollCalls
}

func waitForDeviceLoginFinalState(t *testing.T, mgr *AuthManager, accountID string, timeout time.Duration) (DeviceLoginStatus, *AuthSession, error) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		status, session, err := mgr.GetDeviceLoginStatus(accountID)
		if status != DeviceLoginStatusPending {
			return status, session, err
		}
		time.Sleep(10 * time.Millisecond)
	}

	status, session, err := mgr.GetDeviceLoginStatus(accountID)
	t.Fatalf("timed out waiting for final status; current status=%s err=%v", status, err)
	return status, session, fmt.Errorf("unreachable")
}
