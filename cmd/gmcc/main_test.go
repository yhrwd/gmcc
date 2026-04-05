package main

import (
	"path/filepath"
	"testing"

	"gmcc/internal/config"
)

func TestBuildRuntime_InitializesVaultAndMetadata(t *testing.T) {
	t.Setenv("GMCC_AUTH_VAULT_KEY", "test-master-key")

	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "config.yaml")
	cfg := config.Default()
	cfg.Auth.Vault.Path = ".authvault"

	runtime, err := buildRuntime(configPath, &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if runtime.VaultRepository == nil {
		t.Fatal("expected vault repository")
	}
	if runtime.AccountRepository == nil {
		t.Fatal("expected account repository")
	}
	if runtime.InstanceRepository == nil {
		t.Fatal("expected instance repository")
	}
	if runtime.AuthManager == nil {
		t.Fatal("expected auth manager")
	}
	if runtime.ResourceManager == nil {
		t.Fatal("expected resource manager")
	}
	if runtime.ClusterManager == nil {
		t.Fatal("expected cluster manager")
	}

	vaultDir := filepath.Join(tmp, ".authvault")
	if got := runtime.VaultRepository; got == nil {
		t.Fatal("expected vault repository to remain available")
	}
	if runtimeBaseDir(configPath) != tmp {
		t.Fatalf("want runtime base dir %q, got %q", tmp, runtimeBaseDir(configPath))
	}
	if resolveRuntimePath(tmp, cfg.Auth.Vault.Path) != vaultDir {
		t.Fatalf("want resolved vault dir %q, got %q", vaultDir, resolveRuntimePath(tmp, cfg.Auth.Vault.Path))
	}
}

func TestBuildRuntime_RequiresVaultKeyEnv(t *testing.T) {
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "config.yaml")
	cfg := config.Default()

	if _, err := buildRuntime(configPath, &cfg); err == nil {
		t.Fatal("expected missing vault key error")
	}
}
