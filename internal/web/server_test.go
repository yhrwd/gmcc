package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"gmcc/internal/cluster"
	"gmcc/internal/webtypes"
	"gmcc/internal/webui"
)

func newTestServerWithAssets(t *testing.T, fsys fstest.MapFS) *Server {
	t.Helper()

	server, err := NewServer(webtypes.WebConfig{}, "", cluster.NewManager(cluster.ClusterConfig{}, nil), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	server.uiAssets = webui.NewFSAssets(fsys)
	server.setupRoutes()
	return server
}

func TestAuditLogDirUsesConfigPathBaseDir(t *testing.T) {
	configPath := filepath.Join(`C:\runtime`, "config.yaml")
	got := auditLogDir(configPath)
	want := filepath.Join(`C:\runtime`, "logs", "audit")
	if got != want {
		t.Fatalf("want audit log dir %q, got %q", want, got)
	}
}

func TestAuditLogDirFallsBackToRelativeLogsDir(t *testing.T) {
	got := auditLogDir("")
	want := filepath.Join("logs", "audit")
	if got != want {
		t.Fatalf("want fallback audit log dir %q, got %q", want, got)
	}
}

func TestNewServerInitializesConfigRelativeAuditDir(t *testing.T) {
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "config.yaml")

	server, err := NewServer(webtypes.WebConfig{}, configPath, cluster.NewManager(cluster.ClusterConfig{}, nil), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	auditDir := filepath.Join(tmp, "logs", "audit")
	info, err := os.Stat(auditDir)
	if err != nil {
		t.Fatalf("expected audit dir to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", auditDir)
	}
	if server.auditLogger == nil {
		t.Fatal("expected audit logger to be initialized")
	}
}

func TestNoRouteReturns503WhenFrontendUnavailable(t *testing.T) {
	server := newTestServerWithAssets(t, fstest.MapFS{
		"dist/.keep": {Data: []byte{}},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	if resp.Code != http.StatusServiceUnavailable {
		t.Fatalf("want %d, got %d", http.StatusServiceUnavailable, resp.Code)
	}
	if got := resp.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Fatalf("unexpected content type: %q", got)
	}
}

func TestNoRouteServesIndexForSPAPath(t *testing.T) {
	server := newTestServerWithAssets(t, fstest.MapFS{
		"dist/index.html": {Data: []byte("<html>ui</html>")},
	})

	req := httptest.NewRequest(http.MethodGet, "/instances/demo", nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("want %d, got %d", http.StatusOK, resp.Code)
	}
	if body := resp.Body.String(); body != "<html>ui</html>" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestNoRouteServesStaticAsset(t *testing.T) {
	server := newTestServerWithAssets(t, fstest.MapFS{
		"dist/index.html":    {Data: []byte("<html>ui</html>")},
		"dist/assets/app.js": {Data: []byte("console.log('ok')")},
	})

	req := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("want %d, got %d", http.StatusOK, resp.Code)
	}
	if got := resp.Header().Get("Content-Type"); got != "text/javascript; charset=utf-8" {
		t.Fatalf("unexpected content type: %q", got)
	}
	if body := resp.Body.String(); body != "console.log('ok')" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestNoRouteReturns404ForNonGetPageMethod(t *testing.T) {
	server := newTestServerWithAssets(t, fstest.MapFS{
		"dist/index.html": {Data: []byte("<html>ui</html>")},
	})

	req := httptest.NewRequest(http.MethodPost, "/instances/demo", nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("want %d, got %d", http.StatusNotFound, resp.Code)
	}
}

func TestNoRouteKeepsAPI404InAPISpace(t *testing.T) {
	server := newTestServerWithAssets(t, fstest.MapFS{
		"dist/index.html": {Data: []byte("<html>ui</html>")},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/not-found", nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("want %d, got %d", http.StatusNotFound, resp.Code)
	}
}

func TestHasEmbeddedUIReflectsAssetAvailability(t *testing.T) {
	withUI := newTestServerWithAssets(t, fstest.MapFS{
		"dist/index.html": {Data: []byte("<html>ui</html>")},
	})
	if !withUI.HasEmbeddedUI() {
		t.Fatal("expected embedded UI to be available")
	}

	withoutUI := newTestServerWithAssets(t, fstest.MapFS{
		"dist/.keep": {Data: []byte{}},
	})
	if withoutUI.HasEmbeddedUI() {
		t.Fatal("expected embedded UI to be unavailable")
	}
}
