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

// Deprecated: use internal/auth/session.NewTokenCache with account-scoped IDs.
func New(playerID string) *TokenCache {
	return authsession.NewTokenCache(playerID)
}

// Deprecated: use internal/auth/session.TokenStore.Path with account-scoped IDs.
func Path(playerID string) string {
	store := authsession.NewTokenStore(sessionDir)
	return store.Path(playerID)
}

// Deprecated: use internal/auth/session.TokenStore.Load with account-scoped IDs.
func Load(playerID string) (*TokenCache, error) {
	store := authsession.NewTokenStore(sessionDir)
	return store.Load(playerID)
}

// Deprecated: use internal/auth/session.TokenStore.Save with account-scoped IDs.
func Save(playerID string, cache *TokenCache) error {
	store := authsession.NewTokenStore(sessionDir)
	return store.Save(playerID, cache)
}
