package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gmcc/internal/auth/session"
)

func TestRepository_PutAndGetAccount(t *testing.T) {
	tmp := t.TempDir()
	repo, err := NewRepository(Config{
		Dir:       filepath.Join(tmp, ".authvault"),
		MasterKey: []byte("test-master-key"),
		ScryptN:   1 << 15,
		ScryptR:   8,
		ScryptP:   1,
		SaltLen:   16,
	})
	if err != nil {
		t.Fatal(err)
	}

	record := &session.AccountAuthRecord{
		AccountID: "acc-main",
		UpdatedAt: time.Now().UTC().Round(0),
		Microsoft: session.MicrosoftCredential{RefreshToken: "rtok"},
		Minecraft: session.MinecraftSessionState{ProfileID: "uuid", ProfileName: "Steve"},
	}
	if err := repo.PutAccount(record); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetAccount("acc-main")
	if err != nil {
		t.Fatal(err)
	}
	if got.AccountID != "acc-main" {
		t.Fatalf("want account id acc-main, got %q", got.AccountID)
	}
	if got.Microsoft.RefreshToken != "rtok" {
		t.Fatalf("want refresh token persisted")
	}
	if got.Minecraft.ProfileID != "uuid" || got.Minecraft.ProfileName != "Steve" {
		t.Fatalf("unexpected minecraft state: %+v", got.Minecraft)
	}

	raw, err := os.ReadFile(filepath.Join(tmp, ".authvault", fileNameForAccount("acc-main")))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "rtok") {
		t.Fatalf("vault payload must not contain plaintext refresh token")
	}
}

func TestRepository_PathUsesSanitizedBaseAndHash(t *testing.T) {
	name := fileNameForAccount("acc/main")
	if !strings.HasSuffix(name, ".vault") || !strings.Contains(name, "--") {
		t.Fatalf("unexpected file name %q", name)
	}
	if !strings.HasPrefix(name, "acc-main--") {
		t.Fatalf("want sanitized prefix, got %q", name)
	}
	if fileNameForAccount("acc/main") == fileNameForAccount("acc-main") {
		t.Fatalf("different account ids must not collapse to the same vault file name")
	}
	if strings.Contains(name, "/") {
		t.Fatalf("file name must be filesystem-safe: %q", name)
	}
}

func TestRepository_GetAccountMissingReturnsNotFound(t *testing.T) {
	repo, err := NewRepository(Config{
		Dir:       filepath.Join(t.TempDir(), ".authvault"),
		MasterKey: []byte("test-master-key"),
		ScryptN:   1 << 15,
		ScryptR:   8,
		ScryptP:   1,
		SaltLen:   16,
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := repo.GetAccount("missing"); err != session.ErrAccountNotFound {
		t.Fatalf("want ErrAccountNotFound, got %v", err)
	}
}

func TestRepository_DeleteAccountMissingIsNoop(t *testing.T) {
	repo, err := NewRepository(Config{
		Dir:       filepath.Join(t.TempDir(), ".authvault"),
		MasterKey: []byte("test-master-key"),
		ScryptN:   1 << 15,
		ScryptR:   8,
		ScryptP:   1,
		SaltLen:   16,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.DeleteAccount("missing"); err != nil {
		t.Fatalf("want nil error for missing account delete, got %v", err)
	}
}
