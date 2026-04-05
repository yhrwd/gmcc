package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func NewTokenCache(accountID string) *TokenCache {
	return &TokenCache{AccountID: strings.TrimSpace(accountID)}
}

type TokenStore struct {
	dir string
	mu  sync.Mutex
}

func NewTokenStore(dir string) *TokenStore {
	return &TokenStore{dir: strings.TrimSpace(dir)}
}

func (s *TokenStore) Path(accountID string) string {
	return filepath.Join(s.dir, sanitizeID(accountID)+".json")
}

func (s *TokenStore) Load(accountID string) (*TokenCache, error) {
	cleanID := strings.TrimSpace(accountID)
	path := s.Path(cleanID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewTokenCache(cleanID), nil
		}
		return nil, fmt.Errorf("read account token cache: %w", err)
	}

	cache := NewTokenCache(cleanID)
	if len(data) == 0 {
		return cache, nil
	}
	var header struct {
		AccountID *string `json:"account_id"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, fmt.Errorf("decode account token cache: %w", err)
	}
	if header.AccountID == nil || strings.TrimSpace(*header.AccountID) == "" {
		return NewTokenCache(cleanID), nil
	}
	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("decode account token cache: %w", err)
	}
	cache.AccountID = strings.TrimSpace(cache.AccountID)
	return cache, nil
}

func (s *TokenStore) Save(accountID string, cache *TokenCache) error {
	if cache == nil {
		return fmt.Errorf("token cache is nil")
	}

	cleanID := strings.TrimSpace(accountID)
	cache.AccountID = cleanID
	cache.UpdatedAt = time.Now().UTC()

	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return fmt.Errorf("create account token cache dir: %w", err)
	}

	payload, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("encode account token cache: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.Path(cleanID)
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o600); err != nil {
		return fmt.Errorf("write account token cache temp: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(path)
		if err2 := os.Rename(tmpPath, path); err2 != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("replace account token cache: %w", err2)
		}
	}

	return nil
}

func (s *TokenStore) Delete(accountID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.Path(accountID)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete account token cache: %w", err)
	}
	return nil
}

func sanitizeID(accountID string) string {
	raw := strings.TrimSpace(accountID)
	if raw == "" {
		return "default"
	}

	var b strings.Builder
	b.Grow(len(raw))
	for _, r := range raw {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}

	sanitized := strings.Trim(b.String(), "_")
	if sanitized == "" {
		return "default"
	}
	return sanitized
}
