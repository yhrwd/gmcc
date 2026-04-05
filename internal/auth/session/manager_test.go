package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"gmcc/internal/logx"
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
	store := newFakeRecordStore(&AccountAuthRecord{
		AccountID: "acc-main",
		Microsoft: MicrosoftCredential{RefreshToken: "refresh-token"},
	})
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

func TestAuthManager_GetAccountAuthStatus(t *testing.T) {
	tests := []struct {
		name   string
		record *AccountAuthRecord
		want   AccountAuthStatus
	}{
		{name: "missing", record: nil, want: AccountAuthStatusNotLoggedIn},
		{
			name: "valid refresh token",
			record: &AccountAuthRecord{
				AccountID: "acc",
				Microsoft: MicrosoftCredential{RefreshToken: "rtok"},
			},
			want: AccountAuthStatusLoggedIn,
		},
		{
			name: "history only",
			record: &AccountAuthRecord{
				AccountID: "acc",
				Minecraft: MinecraftSessionState{ProfileID: "uuid", ProfileName: "Steve"},
			},
			want: AccountAuthStatusNotLoggedIn,
		},
		{
			name: "last auth error without refresh token",
			record: &AccountAuthRecord{
				AccountID:     "acc",
				LastAuthError: ErrRefreshTokenInvalid.Error(),
			},
			want: AccountAuthStatusAuthInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newFakeRecordStore(tt.record)
			mgr := NewAuthManager(store, newFakeProvider())

			got, err := mgr.GetAccountAuthStatus("acc")
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("want %s, got %s", tt.want, got)
			}
		})
	}
}

func TestAuthManager_GetAccountProfile(t *testing.T) {
	store := newFakeRecordStore(&AccountAuthRecord{
		AccountID: "acc-main",
		Minecraft: MinecraftSessionState{ProfileID: "uuid-main", ProfileName: "Steve"},
	})
	mgr := NewAuthManager(store, newFakeProvider())

	profile, err := mgr.GetAccountProfile("acc-main")
	if err != nil {
		t.Fatal(err)
	}
	if profile.ProfileID != "uuid-main" || profile.ProfileName != "Steve" {
		t.Fatalf("unexpected profile: %+v", profile)
	}
}

func TestClassifyAuthRecord(t *testing.T) {
	tests := []struct {
		name   string
		record *AccountAuthRecord
		want   AccountAuthStatus
	}{
		{name: "missing", record: nil, want: AccountAuthStatusNotLoggedIn},
		{
			name: "valid refresh token",
			record: &AccountAuthRecord{
				AccountID: "acc",
				Microsoft: MicrosoftCredential{RefreshToken: "rtok"},
			},
			want: AccountAuthStatusLoggedIn,
		},
		{
			name: "history only",
			record: &AccountAuthRecord{
				AccountID: "acc",
				Minecraft: MinecraftSessionState{ProfileID: "uuid", ProfileName: "Steve"},
			},
			want: AccountAuthStatusNotLoggedIn,
		},
		{
			name: "last auth error without refresh token",
			record: &AccountAuthRecord{
				AccountID:     "acc",
				LastAuthError: ErrRefreshTokenInvalid.Error(),
			},
			want: AccountAuthStatusAuthInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := classifyAuthRecord(tt.record); got != tt.want {
				t.Fatalf("want %s, got %s", tt.want, got)
			}
		})
	}
}

func TestAccountAuthRecord_ToAuthSession(t *testing.T) {
	record := &AccountAuthRecord{
		AccountID: "acc-main",
		Microsoft: MicrosoftCredential{
			ExpiresAt: time.Now().UTC().Add(30 * time.Minute),
		},
		Minecraft: MinecraftSessionState{
			AccessToken: "mc-token",
			ExpiresAt:   time.Now().UTC().Add(10 * time.Minute),
			ProfileID:   "player-uuid",
			ProfileName: "Steve",
		},
	}

	session := record.ToAuthSession(AuthSourceCache)
	if session.AccountID != "acc-main" {
		t.Fatalf("want account id acc-main, got %q", session.AccountID)
	}
	if session.MinecraftAccessToken != "mc-token" {
		t.Fatalf("want minecraft token mc-token, got %q", session.MinecraftAccessToken)
	}
	if session.ProfileID != "player-uuid" || session.ProfileName != "Steve" {
		t.Fatalf("unexpected profile data: %+v", session)
	}
	if session.Source != AuthSourceCache {
		t.Fatalf("want source %s, got %s", AuthSourceCache, session.Source)
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
	store := newFakeRecordStore()
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
			wantErrIs:      ErrRefreshUpstream,
			wantPollAtMost: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := newFakeProvider()
			tt.setupProvider(provider)
			mgr := NewAuthManager(newFakeRecordStore(&AccountAuthRecord{
				AccountID: "acc-main",
				Microsoft: MicrosoftCredential{RefreshToken: "refresh-token"},
			}), provider)

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
	mgr := NewAuthManager(newFakeRecordStore(&AccountAuthRecord{
		AccountID: "acc-main",
		Microsoft: MicrosoftCredential{RefreshToken: "refresh-token"},
	}), provider)

	_, err := mgr.Refresh(context.Background(), "acc-main")
	if !errors.Is(err, ErrRefreshUpstream) {
		t.Fatalf("want ErrRefreshUpstream, got %v", err)
	}
}

func TestAuthManager_GetSessionUsesValidMicrosoftAccessWithoutRefreshToken(t *testing.T) {
	provider := newFakeProvider()
	provider.refreshResult = fakeRefreshResult{
		session: AuthSession{
			AccountID:            "acc-main",
			MinecraftAccessToken: "mc-from-access",
			ProfileID:            "uuid",
			ProfileName:          "Steve",
		},
	}
	store := newFakeRecordStore()
	store.records["acc-main"] = &AccountAuthRecord{
		AccountID: "acc-main",
		Microsoft: MicrosoftCredential{
			AccessToken: "ms-access",
			ExpiresAt:   time.Now().Add(10 * time.Minute),
		},
	}
	mgr := NewAuthManager(store, provider)

	session, err := mgr.GetSession(context.Background(), "acc-main")
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if session.MinecraftAccessToken != "mc-from-access" {
		t.Fatalf("unexpected session token: %q", session.MinecraftAccessToken)
	}
	if provider.refreshCalls != 0 {
		t.Fatalf("refresh should not be called, got %d", provider.refreshCalls)
	}
}

func TestAuthManager_PersistsLastAuthError(t *testing.T) {
	provider := newFakeProvider()
	provider.refreshErr = errors.New("invalid_grant")
	store := newFakeRecordStore(&AccountAuthRecord{
		AccountID: "acc-main",
		Microsoft: MicrosoftCredential{RefreshToken: "refresh-token"},
	})
	mgr := NewAuthManager(store, provider)

	_, err := mgr.Refresh(context.Background(), "acc-main")
	if !errors.Is(err, ErrRefreshTokenInvalid) {
		t.Fatalf("want ErrRefreshTokenInvalid, got %v", err)
	}
	record, loadErr := store.GetAccount("acc-main")
	if loadErr != nil {
		t.Fatalf("load record: %v", loadErr)
	}
	if record.LastAuthError == "" {
		t.Fatalf("expected LastAuthError to be persisted")
	}
}

func TestAuthManager_EmitsAuthEvents(t *testing.T) {
	tests := []struct {
		name         string
		action       func(t *testing.T, mgr *AuthManager)
		wantAction   string
		wantResult   string
		wantAuthErr  string
		wantContains string
	}{
		{
			name: "cache hit",
			action: func(t *testing.T, mgr *AuthManager) {
				t.Helper()
				store := mgr.store.(*fakeRecordStore)
				store.records["acc-main"] = &AccountAuthRecord{
					AccountID: "acc-main",
					Minecraft: MinecraftSessionState{
						AccessToken: "mc-cache",
						ProfileID:   "uuid-cache",
						ProfileName: "Steve",
						ExpiresAt:   time.Now().Add(10 * time.Minute),
					},
				}
				if _, err := mgr.GetSession(context.Background(), "acc-main"); err != nil {
					t.Fatalf("get session: %v", err)
				}
			},
			wantAction:   "cache_hit",
			wantResult:   "success",
			wantContains: "account session cache hit",
		},
		{
			name: "refresh",
			action: func(t *testing.T, mgr *AuthManager) {
				t.Helper()
				if _, err := mgr.Refresh(context.Background(), "acc-main"); err != nil {
					t.Fatalf("refresh: %v", err)
				}
			},
			wantAction:   "refresh",
			wantResult:   "success",
			wantContains: "account session refreshed",
		},
		{
			name: "device login start and clear",
			action: func(t *testing.T, mgr *AuthManager) {
				t.Helper()
				provider := mgr.provider.(*fakeProvider)
				provider.deviceLoginInfo = DeviceLoginInfo{
					AccountID:       "acc-main",
					UserCode:        "ABCD",
					VerificationURI: "https://microsoft.com/devicelogin",
					ExpiresAt:       time.Now().UTC().Add(2 * time.Second),
					PollInterval:    50 * time.Millisecond,
				}
				if _, err := mgr.BeginDeviceLogin(context.Background(), "acc-main"); err != nil {
					t.Fatalf("begin device login: %v", err)
				}
				if err := mgr.Clear("acc-main"); err != nil {
					t.Fatalf("clear: %v", err)
				}
			},
			wantAction:   "clear",
			wantResult:   "success",
			wantContains: "account session cleared",
		},
		{
			name: "auth failure",
			action: func(t *testing.T, mgr *AuthManager) {
				t.Helper()
				mgr.provider.(*fakeProvider).refreshErr = errors.New("invalid_grant")
				_, err := mgr.Refresh(context.Background(), "acc-main")
				if !errors.Is(err, ErrRefreshTokenInvalid) {
					t.Fatalf("want ErrRefreshTokenInvalid, got %v", err)
				}
			},
			wantAction:   "auth_failed",
			wantResult:   "failed",
			wantAuthErr:  "refresh_token_invalid",
			wantContains: "account session authentication failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logDir := t.TempDir()
			if err := logx.Init(logDir, true, 1024*1024, false); err != nil {
				t.Fatalf("init logx: %v", err)
			}
			defer logx.Close()

			provider := newFakeProvider()
			provider.refreshResult = fakeRefreshResult{
				session: AuthSession{
					AccountID:            "acc-main",
					MinecraftAccessToken: "mc-1",
					ProfileID:            "uuid",
					ProfileName:          "Steve",
				},
			}
			mgr := NewAuthManager(newFakeRecordStore(&AccountAuthRecord{
				AccountID: "acc-main",
				Microsoft: MicrosoftCredential{RefreshToken: "refresh-token"},
			}), provider)

			tt.action(t, mgr)

			events := readAuthEvents(t, filepath.Join(logDir, "gmcc-events.jsonl"))
			event := events[len(events)-1]
			if event["event_type"] != "auth.session" {
				t.Fatalf("want auth.session event type, got %#v", event["event_type"])
			}
			if event["action"] != tt.wantAction {
				t.Fatalf("want action %q, got %#v", tt.wantAction, event["action"])
			}
			if event["account_id"] != "acc-main" {
				t.Fatalf("want account_id acc-main, got %#v", event["account_id"])
			}
			if tt.wantResult != "" && event["result"] != tt.wantResult {
				t.Fatalf("want result %q, got %#v", tt.wantResult, event["result"])
			}
			if tt.wantAuthErr == "" {
				if _, ok := event["auth_error"]; ok {
					t.Fatalf("did not expect auth_error, got %#v", event["auth_error"])
				}
			} else if event["auth_error"] != tt.wantAuthErr {
				t.Fatalf("want auth_error %q, got %#v", tt.wantAuthErr, event["auth_error"])
			}
			if !strings.Contains(event["message"].(string), tt.wantContains) {
				t.Fatalf("want message containing %q, got %#v", tt.wantContains, event["message"])
			}
		})
	}
}

func TestAuthManager_ClearRemovesCacheAndFlowState(t *testing.T) {
	provider := newFakeProvider()
	provider.deviceLoginInfo = DeviceLoginInfo{
		AccountID:       "acc-main",
		UserCode:        "ABCD",
		VerificationURI: "https://microsoft.com/devicelogin",
		ExpiresAt:       time.Now().Add(2 * time.Second),
		PollInterval:    20 * time.Millisecond,
	}
	store := newFakeRecordStore()
	store.records["acc-main"] = &AccountAuthRecord{AccountID: "acc-main", LastAuthError: "old"}
	mgr := NewAuthManager(store, provider)

	if _, err := mgr.BeginDeviceLogin(context.Background(), "acc-main"); err != nil {
		t.Fatalf("begin device login: %v", err)
	}
	if err := mgr.Clear("acc-main"); err != nil {
		t.Fatalf("clear auth state: %v", err)
	}
	_, err := store.GetAccount("acc-main")
	if !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("expected account to be deleted, got %v", err)
	}
	status, _, statusErr := mgr.GetDeviceLoginStatus("acc-main")
	if status != DeviceLoginStatusFailed || !errors.Is(statusErr, ErrDeviceLoginRequired) {
		t.Fatalf("expected cleared device flow state, got status=%s err=%v", status, statusErr)
	}
}

func TestAuthManager_GetSessionSingleFlightsValidMicrosoftAccessPath(t *testing.T) {
	provider := newFakeProvider()
	provider.refreshResult = fakeRefreshResult{
		session: AuthSession{
			AccountID:            "acc-main",
			MinecraftAccessToken: "mc-single-flight",
			ProfileID:            "uuid",
			ProfileName:          "Steve",
		},
	}
	store := newFakeRecordStore()
	store.records["acc-main"] = &AccountAuthRecord{
		AccountID: "acc-main",
		Microsoft: MicrosoftCredential{
			AccessToken: "ms-access",
			ExpiresAt:   time.Now().Add(10 * time.Minute),
		},
	}
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

	if provider.getXSTSCalls() != 1 {
		t.Fatalf("expected 1 XSTS call, got %d", provider.getXSTSCalls())
	}
}

func newFakeRecordStore(records ...*AccountAuthRecord) *fakeRecordStore {
	store := &fakeRecordStore{records: map[string]*AccountAuthRecord{}}
	for _, record := range records {
		if record == nil {
			continue
		}
		copy := *record
		store.records[record.AccountID] = &copy
	}
	return store
}

type fakeRecordStore struct {
	mu      sync.Mutex
	records map[string]*AccountAuthRecord
}

func (f *fakeRecordStore) GetAccount(accountID string) (*AccountAuthRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if record, ok := f.records[accountID]; ok {
		copy := *record
		return &copy, nil
	}
	return nil, ErrAccountNotFound
}

func (f *fakeRecordStore) PutAccount(record *AccountAuthRecord) error {
	if record == nil {
		return errors.New("record is nil")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	copy := *record
	f.records[record.AccountID] = &copy
	return nil
}

func (f *fakeRecordStore) DeleteAccount(accountID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.records, accountID)
	return nil
}

type fakeProvider struct {
	mu              sync.Mutex
	refreshCalls    int
	xstsCalls       int
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
	f.mu.Lock()
	f.xstsCalls++
	f.mu.Unlock()
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

func (f *fakeProvider) getXSTSCalls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.xstsCalls
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

func readAuthEvents(t *testing.T, path string) []map[string]any {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read auth event log: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	events := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var event map[string]any
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("unmarshal auth event %q: %v", line, err)
		}
		events = append(events, event)
	}
	if len(events) == 0 {
		t.Fatal("expected at least one auth event")
	}
	return events
}
