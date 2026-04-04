package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTokenStore_SaveAndLoadByAccountID(t *testing.T) {
	dir := t.TempDir()
	store := NewTokenStore(filepath.Join(dir, ".session"))
	cache := &TokenCache{
		AccountID: "acc-main",
		UpdatedAt: time.Now().UTC(),
		Minecraft: MinecraftTokenCache{
			AccessToken: "mc-token",
			ProfileID:   "player-uuid",
			ProfileName: "Steve",
			ExpiresAt:   time.Now().Add(10 * time.Minute),
		},
	}

	if err := store.Save("acc-main", cache); err != nil {
		t.Fatalf("save cache: %v", err)
	}

	loaded, err := store.Load("acc-main")
	if err != nil {
		t.Fatalf("load cache: %v", err)
	}
	if loaded.AccountID != "acc-main" {
		t.Fatalf("want account id acc-main, got %q", loaded.AccountID)
	}
}

func TestTokenStore_DoesNotFallbackToLegacyPlayerCache(t *testing.T) {
	dir := t.TempDir()
	store := NewTokenStore(filepath.Join(dir, ".session"))
	legacyPath := filepath.Join(dir, ".session", "legacy-player.json")
	writeLegacyPlayerCache(t, legacyPath)

	loaded, err := store.Load("legacy-player")
	if err != nil {
		t.Fatalf("load cache: %v", err)
	}
	if loaded.HasValidMinecraftToken(time.Now()) {
		t.Fatalf("legacy player cache must not be treated as account cache")
	}
}

func writeLegacyPlayerCache(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir legacy dir: %v", err)
	}

	legacy := map[string]any{
		"player_id": "legacy-player",
		"microsoft": map[string]any{
			"access_token":  "legacy-ms-token",
			"refresh_token": "legacy-refresh-token",
			"expires_at":    time.Now().UTC().Add(30 * time.Minute),
		},
		"minecraft": map[string]any{
			"access_token": "legacy-mc-token",
			"profile_id":   "legacy-uuid",
			"profile_name": "LegacySteve",
			"expires_at":   time.Now().UTC().Add(30 * time.Minute),
		},
	}

	buf, err := json.Marshal(legacy)
	if err != nil {
		t.Fatalf("marshal legacy cache: %v", err)
	}
	if err := os.WriteFile(path, buf, 0o644); err != nil {
		t.Fatalf("write legacy cache: %v", err)
	}
}
