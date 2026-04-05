package state

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestAccountRepository_SaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	repo := NewAccountRepository(filepath.Join(tmp, ".state", "accounts.yaml"))
	accounts := []AccountMeta{{AccountID: "acc-main", Enabled: true, Label: "Main"}}

	if err := repo.SaveAll(accounts); err != nil {
		t.Fatal(err)
	}

	got, err := repo.LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 account, got %d", len(got))
	}
	if got[0].AccountID != "acc-main" {
		t.Fatalf("want account id acc-main, got %q", got[0].AccountID)
	}
	if !got[0].Enabled {
		t.Fatal("want account enabled to persist")
	}
	if got[0].Label != "Main" {
		t.Fatalf("want label Main, got %q", got[0].Label)
	}
}

func TestAccountRepository_RejectsDuplicateIDs(t *testing.T) {
	tmp := t.TempDir()
	repo := NewAccountRepository(filepath.Join(tmp, ".state", "accounts.yaml"))

	err := repo.SaveAll([]AccountMeta{{AccountID: "acc-main"}, {AccountID: "acc-main"}})
	if err == nil {
		t.Fatal("want duplicate error")
	}
	if !strings.Contains(err.Error(), "duplicate account_id") {
		t.Fatalf("want duplicate account_id error, got %v", err)
	}
}

func TestAccountRepository_RejectsEmptyAccountID(t *testing.T) {
	tmp := t.TempDir()
	repo := NewAccountRepository(filepath.Join(tmp, ".state", "accounts.yaml"))

	err := repo.SaveAll([]AccountMeta{{Enabled: true}})
	if err == nil {
		t.Fatal("want validation error")
	}
	if !strings.Contains(err.Error(), "account_id is required") {
		t.Fatalf("want account_id required error, got %v", err)
	}
}

func TestInstanceRepository_SaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	repo := NewInstanceRepository(filepath.Join(tmp, ".state", "instances.yaml"))
	instances := []InstanceMeta{{
		InstanceID:    "bot-1",
		AccountID:     "acc-main",
		ServerAddress: "mc.example.com",
		Enabled:       true,
	}}

	if err := repo.SaveAll(instances); err != nil {
		t.Fatal(err)
	}

	got, err := repo.LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 instance, got %d", len(got))
	}
	if got[0].InstanceID != "bot-1" {
		t.Fatalf("want instance id bot-1, got %q", got[0].InstanceID)
	}
	if got[0].AccountID != "acc-main" {
		t.Fatalf("want account id acc-main, got %q", got[0].AccountID)
	}
	if got[0].ServerAddress != "mc.example.com" {
		t.Fatalf("want server address mc.example.com, got %q", got[0].ServerAddress)
	}
	if !got[0].Enabled {
		t.Fatal("want instance enabled to persist")
	}
}

func TestInstanceRepository_RejectsDuplicateIDs(t *testing.T) {
	tmp := t.TempDir()
	repo := NewInstanceRepository(filepath.Join(tmp, ".state", "instances.yaml"))

	err := repo.SaveAll([]InstanceMeta{{InstanceID: "bot-1", AccountID: "acc-main", ServerAddress: "mc.example.com"}, {InstanceID: "bot-1", AccountID: "acc-alt", ServerAddress: "mc2.example.com"}})
	if err == nil {
		t.Fatal("want duplicate error")
	}
	if !strings.Contains(err.Error(), "duplicate instance_id") {
		t.Fatalf("want duplicate instance_id error, got %v", err)
	}
}

func TestInstanceRepository_RejectsMissingRequiredFields(t *testing.T) {
	tmp := t.TempDir()
	repo := NewInstanceRepository(filepath.Join(tmp, ".state", "instances.yaml"))

	err := repo.SaveAll([]InstanceMeta{{InstanceID: "bot-1"}})
	if err == nil {
		t.Fatal("want validation error")
	}
	if !strings.Contains(err.Error(), "account_id is required") {
		t.Fatalf("want account_id required error, got %v", err)
	}
}
