package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_DefaultsAuthVault(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.yaml")
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Auth.Vault.Path != ".authvault" {
		t.Fatalf("want default path, got %q", cfg.Auth.Vault.Path)
	}
	if cfg.Auth.Vault.KeyEnv != "GMCC_AUTH_VAULT_KEY" {
		t.Fatalf("want default key env, got %q", cfg.Auth.Vault.KeyEnv)
	}
	if cfg.Auth.Vault.ScryptN != 1<<20 {
		t.Fatalf("want default scrypt_n, got %d", cfg.Auth.Vault.ScryptN)
	}
	if cfg.Auth.Vault.ScryptR != 8 {
		t.Fatalf("want default scrypt_r, got %d", cfg.Auth.Vault.ScryptR)
	}
	if cfg.Auth.Vault.ScryptP != 1 {
		t.Fatalf("want default scrypt_p, got %d", cfg.Auth.Vault.ScryptP)
	}
	if cfg.Auth.Vault.SaltLen != 32 {
		t.Fatalf("want default salt_len, got %d", cfg.Auth.Vault.SaltLen)
	}
}

func TestGenerateConfigOmitsWebPasswordAuth(t *testing.T) {
	data, err := generateConfigWithComments(Default())
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if strings.Contains(content, "token_expiry") {
		t.Fatalf("config should not contain token_expiry: %s", content)
	}
	if strings.Contains(content, "passwords:") {
		t.Fatalf("config should not contain password list: %s", content)
	}
}
