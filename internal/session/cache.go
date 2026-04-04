package session

import (
	authsession "gmcc/internal/auth/session"
)

const (
	sessionDir = ".session"
)

type MicrosoftTokenCache = authsession.MicrosoftTokenCache

type MinecraftTokenCache = authsession.MinecraftTokenCache

type TokenCache = authsession.TokenCache

// Deprecated: official-auth cache ownership moved to internal/auth/session.
// Use internal/auth/session.NewTokenCache with account-scoped accountID.
func New(accountID string) *TokenCache {
	return authsession.NewTokenCache(accountID)
}

// Deprecated: use internal/auth/session.TokenStore.Path with account-scoped accountID.
func Path(accountID string) string {
	store := authsession.NewTokenStore(sessionDir)
	return store.Path(accountID)
}

// Deprecated: use internal/auth/session.TokenStore.Load with account-scoped accountID.
func Load(accountID string) (*TokenCache, error) {
	store := authsession.NewTokenStore(sessionDir)
	return store.Load(accountID)
}

// Deprecated: use internal/auth/session.TokenStore.Save with account-scoped accountID.
func Save(accountID string, cache *TokenCache) error {
	store := authsession.NewTokenStore(sessionDir)
	return store.Save(accountID, cache)
}
