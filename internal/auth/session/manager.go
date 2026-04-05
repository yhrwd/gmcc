package session

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"gmcc/internal/auth/microsoft"
	"gmcc/internal/logx"
)

type tokenStore interface {
	Load(accountID string) (*TokenCache, error)
	Save(accountID string, cache *TokenCache) error
	Delete(accountID string) error
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
		m.emitAuthEvent("cache_hit", "account session cache hit", accountID, nil, "success")
		return cache.ToAuthSession(AuthSourceCache), nil
	}
	if cache.HasValidMicrosoftAccess(now) {
		return m.completeSessionFromMicrosoft(ctx, strings.TrimSpace(accountID), cache.Microsoft, AuthSourceRefresh)
	}
	if !cache.HasMicrosoftRefreshToken() {
		return AuthSession{}, m.persistAuthFailure(accountID, cache, ErrDeviceLoginRequired)
	}

	return m.singleFlightRefresh(ctx, strings.TrimSpace(accountID), cache)
}

func (m *AuthManager) Refresh(ctx context.Context, accountID string) (AuthSession, error) {
	cache, err := m.store.Load(accountID)
	if err != nil {
		return AuthSession{}, err
	}
	if !cache.HasMicrosoftRefreshToken() {
		return AuthSession{}, m.persistAuthFailure(accountID, cache, ErrDeviceLoginRequired)
	}
	return m.singleFlightRefresh(ctx, strings.TrimSpace(accountID), cache)
}

func (m *AuthManager) Clear(accountID string) error {
	trimmedAccountID := strings.TrimSpace(accountID)

	m.mu.Lock()
	if flow, ok := m.device[trimmedAccountID]; ok {
		if flow.cancel != nil {
			flow.cancel()
		}
		delete(m.device, trimmedAccountID)
	}
	if pending, ok := m.inflight[trimmedAccountID]; ok {
		pending.err = ErrDeviceLoginRequired
		close(pending.done)
		delete(m.inflight, trimmedAccountID)
	}
	m.mu.Unlock()

	if err := m.store.Delete(trimmedAccountID); err != nil {
		return err
	}

	m.emitAuthEvent("clear", "account session cleared", trimmedAccountID, nil, "success")
	return nil
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
			m.emitAuthEvent("cache_hit", "account session cache hit", accountID, nil, "success")
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
		normalized := normalizeProviderError(err)
		return AuthSession{}, m.persistAuthFailure(accountID, cache, normalized)
	}

	return m.completeSessionFromMicrosoft(ctx, accountID, msCache, AuthSourceRefresh)
}

func (m *AuthManager) completeSessionFromMicrosoft(ctx context.Context, accountID string, msCache MicrosoftTokenCache, source AuthSource) (AuthSession, error) {

	xsts, err := m.provider.GetXSTSFromMicrosoft(ctx, msCache.AccessToken)
	if err != nil {
		normalized := normalizeProviderError(err)
		return AuthSession{}, m.persistAuthFailure(accountID, nil, normalized)
	}

	mcCache, err := m.provider.ExchangeMinecraftToken(ctx, xsts)
	if err != nil {
		normalized := normalizeProviderError(err)
		return AuthSession{}, m.persistAuthFailure(accountID, nil, normalized)
	}

	if err := m.provider.VerifyOwnership(ctx, mcCache.AccessToken); err != nil {
		normalized := normalizeProviderError(err)
		return AuthSession{}, m.persistAuthFailure(accountID, nil, normalized)
	}

	profile, err := m.provider.GetProfile(ctx, mcCache.AccessToken)
	if err != nil {
		normalized := normalizeProviderError(err)
		return AuthSession{}, m.persistAuthFailure(accountID, nil, normalized)
	}
	if strings.TrimSpace(profile.ID) == "" || strings.TrimSpace(profile.Name) == "" {
		return AuthSession{}, m.persistAuthFailure(accountID, nil, ErrProfileInvalid)
	}
	mcCache.ProfileID = strings.TrimSpace(profile.ID)
	mcCache.ProfileName = strings.TrimSpace(profile.Name)

	next := &TokenCache{
		AccountID: strings.TrimSpace(accountID),
		UpdatedAt: time.Now().UTC(),
		Microsoft: msCache,
		Minecraft: mcCache,
	}
	if err := m.store.Save(accountID, next); err != nil {
		return AuthSession{}, fmt.Errorf("save refreshed account token cache: %w", err)
	}

	switch source {
	case AuthSourceRefresh:
		m.emitAuthEvent("refresh", "account session refreshed", accountID, nil, "success")
	case AuthSourceDeviceLogin:
		m.emitAuthEvent("device_login_succeeded", "device login succeeded", accountID, nil, "success")
	}

	return next.ToAuthSession(source), nil
}

func (m *AuthManager) BeginDeviceLogin(ctx context.Context, accountID string) (DeviceLoginInfo, error) {
	trimmedAccountID := strings.TrimSpace(accountID)
	info, err := m.provider.BeginDeviceLogin(ctx, trimmedAccountID)
	if err != nil {
		return DeviceLoginInfo{}, normalizeProviderError(err)
	}
	if strings.TrimSpace(info.AccountID) == "" {
		info.AccountID = trimmedAccountID
	}
	if info.PollInterval <= 0 {
		info.PollInterval = 2 * time.Second
	}

	m.emitAuthEvent("device_login_started", "device login started", trimmedAccountID, nil, string(DeviceLoginStatusPending))

	pollCtx, cancel := context.WithCancel(context.Background())

	flow := &deviceLoginFlow{
		info:   info,
		status: DeviceLoginStatusPending,
		cancel: cancel,
	}

	m.mu.Lock()
	previous := m.device[trimmedAccountID]
	m.device[trimmedAccountID] = flow
	m.mu.Unlock()

	if previous != nil && previous.cancel != nil {
		previous.cancel()
	}

	go m.pollDeviceLogin(pollCtx, trimmedAccountID, flow)

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
	trimmedAccountID := strings.TrimSpace(accountID)

	m.mu.Lock()
	flow, ok := m.device[trimmedAccountID]
	if !ok {
		m.mu.Unlock()
		return nil
	}
	flow.status = DeviceLoginStatusCancelled
	flow.err = ErrDeviceLoginRequired
	cancel := flow.cancel
	m.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	m.emitAuthEvent("device_login_cancelled", "device login cancelled", trimmedAccountID, ErrDeviceLoginRequired, string(DeviceLoginStatusCancelled))

	return nil
}

func (m *AuthManager) pollDeviceLogin(ctx context.Context, accountID string, flow *deviceLoginFlow) {
	ticker := time.NewTicker(flow.info.PollInterval)
	defer ticker.Stop()

	for {
		if deviceLoginExpired(flow.info.ExpiresAt) {
			m.finishDeviceFlowExpired(accountID, flow)
			return
		}

		msCache, err := m.provider.PollDeviceLogin(ctx, accountID, flow.info)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			if isDeviceLoginPending(err) {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					continue
				}
			}
			if isDeviceLoginExpiredError(err) {
				m.finishDeviceFlowExpired(accountID, flow)
				return
			}
			m.finishDeviceFlowFailed(accountID, flow, normalizeProviderError(err))
			return
		}

		session, completeErr := m.completeSessionFromMicrosoft(ctx, accountID, msCache, AuthSourceDeviceLogin)
		if completeErr != nil {
			if errors.Is(completeErr, ErrDeviceLoginRequired) {
				m.finishDeviceFlowExpired(accountID, flow)
				return
			}
			m.finishDeviceFlowFailed(accountID, flow, completeErr)
			return
		}

		m.finishDeviceFlowSucceeded(accountID, flow, session)
		return
	}
}

func (m *AuthManager) finishDeviceFlowSucceeded(accountID string, flow *deviceLoginFlow, session AuthSession) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current, ok := m.device[accountID]
	if !ok || current != flow {
		return
	}

	snapshot := session
	current.status = DeviceLoginStatusSucceeded
	current.session = &snapshot
	current.err = nil
}

func (m *AuthManager) finishDeviceFlowExpired(accountID string, flow *deviceLoginFlow) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current, ok := m.device[accountID]
	if !ok || current != flow {
		return
	}

	current.status = DeviceLoginStatusExpired
	current.session = nil
	current.err = ErrDeviceLoginRequired
	m.emitAuthEvent("device_login_expired", "device login expired", accountID, ErrDeviceLoginRequired, string(DeviceLoginStatusExpired))
}

func (m *AuthManager) finishDeviceFlowFailed(accountID string, flow *deviceLoginFlow, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current, ok := m.device[accountID]
	if !ok || current != flow {
		return
	}

	current.status = DeviceLoginStatusFailed
	current.session = nil
	current.err = err
	m.emitAuthEvent("device_login_failed", "device login failed", accountID, err, string(DeviceLoginStatusFailed))
}

func (m *AuthManager) emitAuthEvent(action, message, accountID string, authErr error, result string) {
	logx.Emit(newAuthEvent(action, message, accountID, authErr, result))
}

func newAuthEvent(action, message, accountID string, authErr error, result string) logx.Event {
	level := "info"
	if authErr != nil {
		level = "warn"
	}

	return logx.NewAuthEvent(level, action, message, "", strings.TrimSpace(accountID), authErrorCode(authErr), result)
}

func authErrorCode(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, ErrDeviceLoginRequired):
		return "device_login_required"
	case errors.Is(err, ErrRefreshTokenInvalid):
		return "refresh_token_invalid"
	case errors.Is(err, ErrRefreshUpstream):
		return "refresh_upstream_failed"
	case errors.Is(err, ErrProviderUnavailable):
		return "provider_unavailable"
	case errors.Is(err, ErrOwnershipFailed):
		return "ownership_failed"
	case errors.Is(err, ErrProfileInvalid):
		return "profile_invalid"
	case errors.Is(err, ErrXSTSDenied):
		return "xsts_denied"
	default:
		return "unknown"
	}
}

func deviceLoginExpired(expiresAt time.Time) bool {
	if expiresAt.IsZero() {
		return false
	}
	return !time.Now().UTC().Before(expiresAt)
}

func isDeviceLoginPending(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "authorization_pending") || strings.Contains(msg, "authorization pending") || strings.Contains(msg, "slow_down") || strings.Contains(msg, "slow down")
}

func isDeviceLoginExpiredError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "expired_token") || strings.Contains(msg, "expired token") || strings.Contains(msg, "expired")
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
	case strings.Contains(msg, "temporarily unavailable"), strings.Contains(msg, "upstream unavailable"), strings.Contains(msg, "timeout"), strings.Contains(msg, "context deadline"):
		return ErrRefreshUpstream
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

func (m *AuthManager) persistAuthFailure(accountID string, base *TokenCache, authErr error) error {
	if authErr == nil {
		return nil
	}
	m.emitAuthEvent("auth_failed", "account session authentication failed", accountID, authErr, "failed")
	cache := NewTokenCache(accountID)
	if base != nil {
		copy := *base
		cache = &copy
	}
	cache.AccountID = strings.TrimSpace(accountID)
	cache.LastAuthError = authErr.Error()
	cache.UpdatedAt = time.Now().UTC()
	if saveErr := m.store.Save(accountID, cache); saveErr != nil {
		return fmt.Errorf("persist auth failure: %w", saveErr)
	}
	return authErr
}
