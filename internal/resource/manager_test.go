package resource

import (
	"errors"
	"testing"

	"gmcc/internal/auth/session"
	"gmcc/internal/state"
)

func TestResourceManager_CreateInstanceValidatesRequiredFields(t *testing.T) {
	rm := NewManager(newFakeAccountsRepo(nil), newFakeInstancesRepo(nil), newFakeAuthStatus(nil, nil))

	tests := []struct {
		name  string
		input CreateInstanceInput
		want  error
	}{
		{name: "missing instance id", input: CreateInstanceInput{AccountID: "acc-main", ServerAddress: "mc.example.com"}, want: ErrInstanceIDRequired},
		{name: "missing account id", input: CreateInstanceInput{InstanceID: "bot-1", ServerAddress: "mc.example.com"}, want: ErrAccountIDRequired},
		{name: "missing server address", input: CreateInstanceInput{InstanceID: "bot-1", AccountID: "acc-main"}, want: ErrServerAddressRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := rm.CreateInstance(tt.input)
			if !errors.Is(err, tt.want) {
				t.Fatalf("want %v, got %v", tt.want, err)
			}
		})
	}
}

func TestResourceManager_CreateInstanceRejectsAccountNotLoggedIn(t *testing.T) {
	accounts := newFakeAccountsRepo([]state.AccountMeta{{AccountID: "acc-main", Enabled: true}})
	instances := newFakeInstancesRepo(nil)
	auth := newFakeAuthStatus(map[string]session.AccountAuthStatus{"acc-main": session.AccountAuthStatusNotLoggedIn}, nil)
	rm := NewManager(accounts, instances, auth)

	_, err := rm.CreateInstance(CreateInstanceInput{InstanceID: "bot-1", AccountID: "acc-main", ServerAddress: "mc.example.com"})
	if !errors.Is(err, ErrAccountNotLoggedIn) {
		t.Fatalf("want ErrAccountNotLoggedIn, got %v", err)
	}
}

func TestResourceManager_CreateInstancePersistsValidatedInstance(t *testing.T) {
	accounts := newFakeAccountsRepo([]state.AccountMeta{{AccountID: "acc-main", Enabled: true}})
	instances := newFakeInstancesRepo([]state.InstanceMeta{{InstanceID: "bot-0", AccountID: "acc-main", ServerAddress: "old.example.com", Enabled: true}})
	auth := newFakeAuthStatus(map[string]session.AccountAuthStatus{"acc-main": session.AccountAuthStatusLoggedIn}, nil)
	rm := NewManager(accounts, instances, auth)

	got, err := rm.CreateInstance(CreateInstanceInput{InstanceID: " bot-1 ", AccountID: " acc-main ", ServerAddress: " mc.example.com ", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	if got.InstanceID != "bot-1" || got.AccountID != "acc-main" || got.ServerAddress != "mc.example.com" {
		t.Fatalf("unexpected saved instance: %+v", got)
	}

	if len(instances.saved) != 2 {
		t.Fatalf("want 2 saved instances, got %d", len(instances.saved))
	}
	if instances.saved[1] != got {
		t.Fatalf("want persisted instance %+v, got %+v", got, instances.saved[1])
	}
}

func TestResourceManager_RestoreResourcesAppliesCanonicalReasons(t *testing.T) {
	accounts := newFakeAccountsRepo([]state.AccountMeta{
		{AccountID: "acc-ok", Enabled: true},
		{AccountID: "acc-disabled", Enabled: false},
		{AccountID: "acc-auth-invalid", Enabled: true},
		{AccountID: "acc-not-logged-in", Enabled: true},
	})
	instances := newFakeInstancesRepo([]state.InstanceMeta{
		{InstanceID: "bot-ok", AccountID: "acc-ok", ServerAddress: "mc-1.example.com", Enabled: true},
		{InstanceID: "bot-disabled", AccountID: "acc-disabled", ServerAddress: "mc-2.example.com", Enabled: true},
		{InstanceID: "bot-auth-invalid", AccountID: "acc-auth-invalid", ServerAddress: "mc-3.example.com", Enabled: true},
		{InstanceID: "bot-not-logged-in", AccountID: "acc-not-logged-in", ServerAddress: "mc-4.example.com", Enabled: true},
		{InstanceID: "bot-missing", AccountID: "acc-missing", ServerAddress: "mc-5.example.com", Enabled: true},
	})
	auth := newFakeAuthStatus(map[string]session.AccountAuthStatus{
		"acc-ok":            session.AccountAuthStatusLoggedIn,
		"acc-auth-invalid":  session.AccountAuthStatusAuthInvalid,
		"acc-not-logged-in": session.AccountAuthStatusNotLoggedIn,
	}, nil)
	rm := NewManager(accounts, instances, auth)

	result, err := rm.RestoreResources()
	if err != nil {
		t.Fatal(err)
	}
	if result.RestoredCount != 1 || result.SkippedCount != 4 {
		t.Fatalf("unexpected counts: %+v", result)
	}
	if len(result.RestoredInstances) != 1 || result.RestoredInstances[0].InstanceID != "bot-ok" {
		t.Fatalf("unexpected restored instances: %+v", result.RestoredInstances)
	}
	if got := result.RestoredInstances[0].AccountID; got != "acc-ok" {
		t.Fatalf("unexpected restored account id: %q", got)
	}

	wantReasons := map[string]ReasonCode{
		"bot-disabled":      ReasonAccountDisabled,
		"bot-auth-invalid":  ReasonAccountAuthInvalid,
		"bot-not-logged-in": ReasonAccountNotLoggedIn,
		"bot-missing":       ReasonAccountNotFound,
	}
	if len(result.SkippedInstances) != len(wantReasons) {
		t.Fatalf("want %d skipped instances, got %d", len(wantReasons), len(result.SkippedInstances))
	}
	for _, skipped := range result.SkippedInstances {
		if wantReasons[skipped.InstanceID] != skipped.Reason {
			t.Fatalf("unexpected skip reason for %s: want %s, got %s", skipped.InstanceID, wantReasons[skipped.InstanceID], skipped.Reason)
		}
	}
}

func TestResourceManager_DeleteAccountRejectsReferencedAccount(t *testing.T) {
	accounts := newFakeAccountsRepo([]state.AccountMeta{{AccountID: "acc-main", Enabled: true}})
	instances := newFakeInstancesRepo([]state.InstanceMeta{{InstanceID: "bot-1", AccountID: "acc-main", ServerAddress: "mc.example.com", Enabled: true}})
	rm := NewManager(accounts, instances, newFakeAuthStatus(nil, nil))

	err := rm.DeleteAccount("acc-main")
	if !errors.Is(err, ErrAccountInUse) {
		t.Fatalf("want ErrAccountInUse, got %v", err)
	}
	if len(accounts.accounts) != 1 {
		t.Fatalf("want accounts unchanged, got %+v", accounts.accounts)
	}
}

func TestResourceManager_DeleteAccountRemovesUnreferencedAccount(t *testing.T) {
	accounts := newFakeAccountsRepo([]state.AccountMeta{
		{AccountID: "acc-main", Enabled: true},
		{AccountID: "acc-other", Enabled: true},
	})
	instances := newFakeInstancesRepo([]state.InstanceMeta{{InstanceID: "bot-1", AccountID: "acc-other", ServerAddress: "mc.example.com", Enabled: true}})
	rm := NewManager(accounts, instances, newFakeAuthStatus(nil, nil))

	if err := rm.DeleteAccount("acc-main"); err != nil {
		t.Fatal(err)
	}
	if len(accounts.accounts) != 1 {
		t.Fatalf("want 1 remaining account, got %d", len(accounts.accounts))
	}
	if accounts.accounts[0].AccountID != "acc-other" {
		t.Fatalf("want remaining account acc-other, got %+v", accounts.accounts[0])
	}
}

func TestResourceManager_ListAccountsReturnsMetadata(t *testing.T) {
	accounts := newFakeAccountsRepo([]state.AccountMeta{{AccountID: "acc-main", Enabled: true, Label: "Main"}})
	rm := NewManager(accounts, newFakeInstancesRepo(nil), newFakeAuthStatus(nil, nil))

	records, err := rm.ListAccounts()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Meta.AccountID != "acc-main" {
		t.Fatalf("unexpected records: %+v", records)
	}
	if records[0].Meta.Label != "Main" {
		t.Fatalf("unexpected label: %+v", records[0].Meta)
	}
}

func TestResourceManager_GetAccountReturnsSingleRecord(t *testing.T) {
	accounts := newFakeAccountsRepo([]state.AccountMeta{{AccountID: "acc-main", Enabled: true}})
	rm := NewManager(accounts, newFakeInstancesRepo(nil), newFakeAuthStatus(nil, nil))

	record, err := rm.GetAccount("acc-main")
	if err != nil {
		t.Fatal(err)
	}
	if record.Meta.AccountID != "acc-main" {
		t.Fatalf("unexpected record: %+v", record)
	}
}

type fakeAccountsRepo struct {
	accounts []state.AccountMeta
	err      error
}

func newFakeAccountsRepo(accounts []state.AccountMeta) *fakeAccountsRepo {
	clone := append([]state.AccountMeta(nil), accounts...)
	return &fakeAccountsRepo{accounts: clone}
}

func (r *fakeAccountsRepo) LoadAll() ([]state.AccountMeta, error) {
	if r.err != nil {
		return nil, r.err
	}
	return append([]state.AccountMeta(nil), r.accounts...), nil
}

func (r *fakeAccountsRepo) SaveAll(accounts []state.AccountMeta) error {
	if r.err != nil {
		return r.err
	}
	r.accounts = append([]state.AccountMeta(nil), accounts...)
	return nil
}

type fakeInstancesRepo struct {
	instances []state.InstanceMeta
	saved     []state.InstanceMeta
	err       error
}

func newFakeInstancesRepo(instances []state.InstanceMeta) *fakeInstancesRepo {
	clone := append([]state.InstanceMeta(nil), instances...)
	return &fakeInstancesRepo{instances: clone}
}

func (r *fakeInstancesRepo) LoadAll() ([]state.InstanceMeta, error) {
	if r.err != nil {
		return nil, r.err
	}
	return append([]state.InstanceMeta(nil), r.instances...), nil
}

func (r *fakeInstancesRepo) SaveAll(instances []state.InstanceMeta) error {
	if r.err != nil {
		return r.err
	}
	r.saved = append([]state.InstanceMeta(nil), instances...)
	r.instances = append([]state.InstanceMeta(nil), instances...)
	return nil
}

type fakeAuthStatus struct {
	statuses map[string]session.AccountAuthStatus
	errs     map[string]error
}

func newFakeAuthStatus(statuses map[string]session.AccountAuthStatus, errs map[string]error) *fakeAuthStatus {
	if statuses == nil {
		statuses = map[string]session.AccountAuthStatus{}
	}
	if errs == nil {
		errs = map[string]error{}
	}
	return &fakeAuthStatus{statuses: statuses, errs: errs}
}

func (f *fakeAuthStatus) GetAccountAuthStatus(accountID string) (session.AccountAuthStatus, error) {
	if err, ok := f.errs[accountID]; ok {
		return "", err
	}
	if status, ok := f.statuses[accountID]; ok {
		return status, nil
	}
	return "", session.ErrAccountNotFound
}
