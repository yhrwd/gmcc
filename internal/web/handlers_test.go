package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"gmcc/internal/cluster"
	"gmcc/internal/resource"
	"gmcc/internal/state"
	"gmcc/internal/web/audit"
)

type fakeAccountReader struct {
	list        []resource.AccountRecord
	get         map[string]resource.AccountRecord
	err         error
	createInput resource.CreateAccountInput
	deleteID    string
	deleteErr   error
}

func (f *fakeAccountReader) ListAccounts() ([]resource.AccountRecord, error) {
	if f.err != nil {
		return nil, f.err
	}
	return append([]resource.AccountRecord(nil), f.list...), nil
}

func (f *fakeAccountReader) GetAccount(accountID string) (resource.AccountRecord, error) {
	if f.err != nil {
		return resource.AccountRecord{}, f.err
	}
	record, ok := f.get[accountID]
	if !ok {
		return resource.AccountRecord{}, resource.ErrAccountNotFound
	}
	return record, nil
}

func (f *fakeAccountReader) CreateAccount(in resource.CreateAccountInput) (state.AccountMeta, error) {
	if f.err != nil {
		return state.AccountMeta{}, f.err
	}
	f.createInput = in
	meta := state.AccountMeta{AccountID: in.AccountID, Enabled: in.Enabled, Label: in.Label, Note: in.Note}
	f.list = append(f.list, resource.AccountRecord{Meta: meta})
	if f.get == nil {
		f.get = map[string]resource.AccountRecord{}
	}
	f.get[in.AccountID] = resource.AccountRecord{Meta: meta}
	return meta, nil
}

func (f *fakeAccountReader) DeleteAccount(accountID string) error {
	f.deleteID = accountID
	if f.deleteErr != nil {
		return f.deleteErr
	}
	delete(f.get, accountID)
	return nil
}

func TestHandleGetAccountsReturnsAccountMetadataView(t *testing.T) {
	gin.SetMode(gin.TestMode)

	clusterManager := cluster.NewManager(cluster.ClusterConfig{}, nil)

	server := &Server{
		clusterManager: clusterManager,
		resourceManager: &fakeAccountReader{
			list: []resource.AccountRecord{{
				Meta: state.AccountMeta{AccountID: "acc-main", Enabled: true, Label: "Main", Note: "Primary"},
			}},
		},
		auditLogger: newTestAuditLogger(t),
	}

	ctx, recorder := newJSONContext(http.MethodGet, "/api/accounts")
	server.handleGetAccounts(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", recorder.Code)
	}

	var payload struct {
		Accounts []map[string]any `json:"accounts"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Accounts) != 1 {
		t.Fatalf("want 1 account, got %d", len(payload.Accounts))
	}
	account := payload.Accounts[0]
	if account["id"] != "acc-main" {
		t.Fatalf("unexpected account id: %+v", account)
	}
	if _, exists := account["player_id"]; exists {
		t.Fatalf("account payload should not expose player_id without auth session: %+v", account)
	}
	if account["auth_status"] != "not_logged_in" {
		t.Fatalf("unexpected auth status: %+v", account)
	}
	if _, exists := account["status"]; exists {
		t.Fatalf("account payload should not expose runtime status: %+v", account)
	}
}

func TestHandleGetInstancesReturnsRuntimeInstanceView(t *testing.T) {
	gin.SetMode(gin.TestMode)

	clusterManager := cluster.NewManager(cluster.DefaultClusterConfig(), nil)
	if err := clusterManager.CreateInstance("bot-1", cluster.AccountEntry{ID: "acc-main", ServerAddress: "mc.example.com", Enabled: true}); err != nil {
		t.Fatalf("create instance: %v", err)
	}

	server := &Server{clusterManager: clusterManager, auditLogger: newTestAuditLogger(t)}

	ctx, recorder := newJSONContext(http.MethodGet, "/api/instances")
	server.handleGetInstances(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", recorder.Code)
	}

	var payload struct {
		Instances []map[string]any `json:"instances"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Instances) != 1 {
		t.Fatalf("want 1 instance, got %d", len(payload.Instances))
	}
	instance := payload.Instances[0]
	if instance["id"] != "bot-1" {
		t.Fatalf("unexpected instance id: %+v", instance)
	}
	if instance["account_id"] != "acc-main" {
		t.Fatalf("unexpected account id: %+v", instance)
	}
	if _, exists := instance["auth_status"]; exists {
		t.Fatalf("instance payload should not expose account auth status: %+v", instance)
	}
}

func TestHandleCreateAccountCreatesOnlyAccountMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reader := &fakeAccountReader{get: map[string]resource.AccountRecord{}}
	clusterManager := cluster.NewManager(cluster.DefaultClusterConfig(), nil)
	server := &Server{resourceManager: reader, clusterManager: clusterManager, auditLogger: newTestAuditLogger(t)}

	body := bytes.NewBufferString(`{"id":"acc-main","label":"Main","note":"Primary"}`)
	ctx, recorder := newBodyContext(http.MethodPost, "/api/accounts", body)
	server.handleCreateAccount(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", recorder.Code)
	}
	if reader.createInput.AccountID != "acc-main" {
		t.Fatalf("unexpected create input: %+v", reader.createInput)
	}
	if got := len(clusterManager.ListInstances()); got != 0 {
		t.Fatalf("account create should not create runtime instances, got %d", got)
	}
}

func TestHandleCreateInstanceAllowsDistinctInstanceAndAccountIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	clusterManager := cluster.NewManager(cluster.DefaultClusterConfig(), nil)
	server := &Server{clusterManager: clusterManager, auditLogger: newTestAuditLogger(t)}

	body := bytes.NewBufferString(`{"id":"bot-1","account_id":"acc-main","server_address":"mc.example.com","enabled":true}`)
	ctx, recorder := newBodyContext(http.MethodPost, "/api/instances", body)
	server.handleCreateInstance(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("want 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	info, err := clusterManager.GetInstanceInfo("bot-1")
	if err != nil {
		t.Fatal(err)
	}
	if info.AccountID != "acc-main" {
		t.Fatalf("unexpected account id: %+v", info)
	}
}

func TestHandleDeleteAccountReturnsErrorWhenAccountInUse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reader := &fakeAccountReader{deleteErr: resource.ErrAccountInUse}
	server := &Server{resourceManager: reader, auditLogger: newTestAuditLogger(t)}

	ctx, recorder := newBodyContext(http.MethodDelete, "/api/accounts/acc-main", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "acc-main"}}
	server.handleDeleteAccount(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", recorder.Code)
	}
	if reader.deleteID != "acc-main" {
		t.Fatalf("unexpected delete id: %q", reader.deleteID)
	}
}

func newJSONContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(method, path, nil)
	return ctx, recorder
}

func newBodyContext(method, path string, body *bytes.Buffer) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	var requestBody *bytes.Buffer
	if body == nil {
		requestBody = bytes.NewBuffer(nil)
	} else {
		requestBody = body
	}
	ctx.Request = httptest.NewRequest(method, path, requestBody)
	ctx.Request.Header.Set("Content-Type", "application/json")
	return ctx, recorder
}

func newTestAuditLogger(t *testing.T) *audit.Logger {
	t.Helper()
	logger, err := audit.NewLogger(t.TempDir(), 1)
	if err != nil {
		t.Fatal(err)
	}
	return logger
}
