package session

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"gmcc/internal/auth/microsoft"
)

type tokenStore interface {
	Load(accountID string) (*TokenCache, error)
	Save(accountID string, cache *TokenCache) error
}

type authProvider interface {
	MicrosoftProvider
	MinecraftProvider
}

type inflightResult struct {
	session AuthSession
	err     error
	done    chan struct{}
}

type deviceLoginFlow struct {
	info    DeviceLoginInfo
	status  DeviceLoginStatus
	session *AuthSession
	err     error
	cancel  context.CancelFunc
}

type AuthManager struct {
	store    tokenStore
	provider authProvider
	mu       sync.Mutex
	inflight map[string]*inflightResult
	device   map[string]*deviceLoginFlow
}

func NewAuthManager(store tokenStore, provider authProvider) *AuthManager {
	return &AuthManager{
		store:    store,
		provider: provider,
		inflight: map[string]*inflightResult{},
		device:   map[string]*deviceLoginFlow{},
	}
}

func (m *AuthManager) GetSession(ctx context.Context, accountID string) (AuthSession, error) {
	cache, err := m.store.Load(accountID)
	if err != nil {
		return AuthSession{}, err
	}

	now := time.Now().UTC()
	if cache.HasValidMinecraftToken(now) {
		return cache.ToAuthSession(AuthSourceCache), nil
	}
	if !cache.HasMicrosoftRefreshToken() {
		return AuthSession{}, ErrDeviceLoginRequired
	}

	return m.singleFlightRefresh(ctx, strings.TrimSpace(accountID), cache)
}

func (m *AuthManager) Refresh(ctx context.Context, accountID string) (AuthSession, error) {
	cache, err := m.store.Load(accountID)
	if err != nil {
		return AuthSession{}, err
	}
	if !cache.HasMicrosoftRefreshToken() {
		return AuthSession{}, ErrDeviceLoginRequired
	}
	return m.singleFlightRefresh(ctx, strings.TrimSpace(accountID), cache)
}

func (m *AuthManager) singleFlightRefresh(ctx context.Context, accountID string, cache *TokenCache) (AuthSession, error) {
	m.mu.Lock()
	if pending, ok := m.inflight[accountID]; ok {
		m.mu.Unlock()
		select {
		case <-pending.done:
			return pending.session, pending.err
		case <-ctx.Done():
			return AuthSession{}, ctx.Err()
		}
	}

	pending := &inflightResult{done: make(chan struct{})}
	m.inflight[accountID] = pending
	m.mu.Unlock()

	if latest, err := m.store.Load(accountID); err == nil {
		now := time.Now().UTC()
		if latest.HasValidMinecraftToken(now) {
			session := latest.ToAuthSession(AuthSourceCache)
			m.mu.Lock()
			pending.session = session
			pending.err = nil
			close(pending.done)
			delete(m.inflight, accountID)
			m.mu.Unlock()
			return session, nil
		}
	}

	session, err := m.refreshSession(ctx, accountID, cache)

	m.mu.Lock()
	pending.session = session
	pending.err = err
	close(pending.done)
	delete(m.inflight, accountID)
	m.mu.Unlock()

	return session, err
}

func (m *AuthManager) refreshSession(ctx context.Context, accountID string, cache *TokenCache) (AuthSession, error) {
	msCache, err := m.provider.RefreshMicrosoft(ctx, cache.Microsoft.RefreshToken)
	if err != nil {
		return AuthSession{}, normalizeProviderError(err)
	}

	xsts, err := m.provider.GetXSTSFromMicrosoft(ctx, msCache.AccessToken)
	if err != nil {
		return AuthSession{}, normalizeProviderError(err)
	}

	mcCache, err := m.provider.ExchangeMinecraftToken(ctx, xsts)
	if err != nil {
		return AuthSession{}, normalizeProviderError(err)
	}

	if err := m.provider.VerifyOwnership(ctx, mcCache.AccessToken); err != nil {
		return AuthSession{}, normalizeProviderError(err)
	}

	profile, err := m.provider.GetProfile(ctx, mcCache.AccessToken)
	if err != nil {
		return AuthSession{}, normalizeProviderError(err)
	}
	if strings.TrimSpace(profile.ID) == "" || strings.TrimSpace(profile.Name) == "" {
		return AuthSession{}, ErrProfileInvalid
	}
	mcCache.ProfileID = strings.TrimSpace(profile.ID)
	mcCache.ProfileName = strings.TrimSpace(profile.Name)

	next := &TokenCache{
		AccountID: strings.TrimSpace(accountID),
		Microsoft: msCache,
		Minecraft: mcCache,
	}
	if err := m.store.Save(accountID, next); err != nil {
		return AuthSession{}, fmt.Errorf("save refreshed account token cache: %w", err)
	}

	return next.ToAuthSession(AuthSourceRefresh), nil
}

func (m *AuthManager) BeginDeviceLogin(ctx context.Context, accountID string) (DeviceLoginInfo, error) {
	info, err := m.provider.BeginDeviceLogin(ctx, accountID)
	if err != nil {
		return DeviceLoginInfo{}, normalizeProviderError(err)
	}

	flow := &deviceLoginFlow{
		info:   info,
		status: DeviceLoginStatusPending,
	}
	m.mu.Lock()
	m.device[strings.TrimSpace(accountID)] = flow
	m.mu.Unlock()

	return info, nil
}

func (m *AuthManager) GetDeviceLoginStatus(accountID string) (DeviceLoginStatus, *AuthSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	flow, ok := m.device[strings.TrimSpace(accountID)]
	if !ok {
		return DeviceLoginStatusFailed, nil, ErrDeviceLoginRequired
	}

	return flow.status, flow.session, flow.err
}

func (m *AuthManager) CancelDeviceLogin(accountID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	flow, ok := m.device[strings.TrimSpace(accountID)]
	if !ok {
		return nil
	}
	if flow.cancel != nil {
		flow.cancel()
	}
	flow.status = DeviceLoginStatusCancelled
	flow.err = ErrDeviceLoginRequired
	return nil
}

func normalizeProviderError(err error) error {
	if err == nil {
		return nil
	}

	var xstsErr *microsoft.XSTSError
	if errors.As(err, &xstsErr) {
		return ErrXSTSDenied
	}

	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "refresh_token"), strings.Contains(msg, "refresh token"), strings.Contains(msg, "invalid_grant"):
		return ErrRefreshTokenInvalid
	case strings.Contains(msg, "ownership"), strings.Contains(msg, "entitlement"):
		return ErrOwnershipFailed
	case strings.Contains(msg, "profile"):
		return ErrProfileInvalid
	case strings.Contains(msg, "xsts"):
		return ErrXSTSDenied
	default:
		return ErrProviderUnavailable
	}
}
