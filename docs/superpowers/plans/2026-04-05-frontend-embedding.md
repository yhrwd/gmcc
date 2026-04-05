# Frontend Embedding Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add optional frontend embedding support so gmcc can serve embedded UI assets when available, while still building and running as API-only when frontend assets are absent.

**Architecture:** Keep the existing Gin API routing intact and add a focused `internal/webui` unit that owns embedded asset lookup, fallback decisions, and “frontend available” detection. Add a Go packaging tool that prepares an embed directory from future frontend build output, preserves a compile-safe placeholder, and then builds the final binary.

**Tech Stack:** Go, Gin, `embed`, standard library `net/http`, standard library `os/exec`, Go testing package

---

## File Map

- Create: `frontend/.gitkeep`
- Create: `internal/webui/assets.go`
- Create: `internal/webui/assets_test.go`
- Create: `internal/webui/dist/.keep`
- Create: `tools/packager/main.go`
- Create: `tools/packager/main_test.go`
- Create: `tools/packager/testdata/frontend-dist/index.html`
- Create: `tools/packager/testdata/frontend-dist/assets/app.js`
- Create: `tools/packager/testdata/frontend-dist/assets/app.js.map`
- Create: `tools/packager/testdata/frontend-dist/favicon.ico`
- Create: `tools/packager/testdata/frontend-dist/robots.txt`
- Create: `tools/packager/testdata/frontend-dist/extra.txt`
- Modify: `.gitignore`
- Modify: `internal/web/server.go`
- Modify: `internal/web/server_test.go`
- Modify: `cmd/gmcc/main.go`
- Modify: `README.md`

## Implementation Notes

- `internal/webui/assets.go` should expose one injectable unit, e.g. `type UIAssets interface`, plus a default embedded implementation constructor.
- `internal/web` should depend on that interface, not on package globals.
- Treat API space as `/api`, `/api/`, and `/api/*`.
- Only `GET` and `HEAD` may use static file serving or SPA fallback.
- Non-API `POST`/`PUT`/`DELETE`/`PATCH`/unhandled `OPTIONS` should return `404`.
- Frontend availability must be decided only by presence of embedded `index.html`; `.keep` must never count.
- The first implementation can read embedded files fully into memory.
- Use `c.Request.URL.Path` directly as the `UIAssets` path input; do not add an extra decode pass.
- This first version explicitly does not support `Range`, `If-Modified-Since`, or custom cache negotiation.
- If `frontend/dist` exists but does not yield `index.html` after whitelist filtering, the packager must clean back to `.keep` and continue with an API-only build.

### Task 1: Reserve frontend and embed directories

**Files:**
- Create: `frontend/.gitkeep`
- Create: `internal/webui/dist/.keep`
- Modify: `.gitignore`

- [ ] **Step 1: Add the reserved frontend directory marker**

Create `frontend/.gitkeep` as an empty file.

```text

```

- [ ] **Step 2: Add the compile-safe embed placeholder**

Create `internal/webui/dist/.keep` as an empty file.

```text

```

- [ ] **Step 3: Update git ignore rules for frontend and build artifacts**

Append these rules to `.gitignore`.

```gitignore
frontend/dist/
frontend/node_modules/
internal/webui/dist/*
!internal/webui/dist/.keep
build/
release/
```

- [ ] **Step 4: Verify the placeholder stays tracked while artifacts stay ignored**

Run: `git status --short`
Expected: `frontend/.gitkeep` and `internal/webui/dist/.keep` appear as tracked/new files, while no future files under `internal/webui/dist/` other than `.keep` would be intended for tracking.

- [ ] **Step 5: Commit**

```bash
git add .gitignore frontend/.gitkeep internal/webui/dist/.keep
git commit -m "chore: reserve frontend embed directories"
```

### Task 2: Add embedded UI asset unit with tests

**Files:**
- Create: `internal/webui/assets.go`
- Create: `internal/webui/assets_test.go`
- Test: `internal/webui/assets_test.go`

- [ ] **Step 1: Write the failing tests for asset availability, lookup, and path classification**

Create `internal/webui/assets_test.go` with these tests.

```go
package webui

import (
	"io/fs"
	"testing"
	"testing/fstest"
	"time"
)

func TestHasIndexReturnsFalseWhenOnlyKeepFileExists(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{
		"dist/.keep": {Data: []byte{}},
	})

	if assets.HasIndex() {
		t.Fatal("expected HasIndex to be false when only .keep exists")
	}
}

func TestHasIndexReturnsTrueWhenIndexExists(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{
		"dist/.keep":     {Data: []byte{}},
		"dist/index.html": {Data: []byte("<html></html>")},
	})

	if !assets.HasIndex() {
		t.Fatal("expected HasIndex to be true")
	}
}

func TestLookupAssetReturnsNilForMissingFile(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{
		"dist/.keep": {Data: []byte{}},
	})

	file, err := assets.LookupAsset("/assets/app.js")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if file != nil {
		t.Fatal("expected nil file for missing asset")
	}
}

func TestLookupAssetReturnsAssetMetadata(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{
		"dist/assets/app.js": {
			Data:    []byte("console.log('ok')"),
			Mode:    fs.FileMode(0644),
			ModTime: time.Unix(1700000000, 0),
		},
	})

	file, err := assets.LookupAsset("/assets/app.js")
	if err != nil {
		t.Fatalf("lookup failed: %v", err)
	}
	if file == nil {
		t.Fatal("expected asset file")
	}
	if file.Name != "assets/app.js" {
		t.Fatalf("expected asset name %q, got %q", "assets/app.js", file.Name)
	}
	if file.ContentType != "text/javascript; charset=utf-8" {
		t.Fatalf("unexpected content type: %q", file.ContentType)
	}
	if string(file.Content) != "console.log('ok')" {
		t.Fatalf("unexpected content: %q", string(file.Content))
	}
}

func TestLookupAssetRejectsDirectoryTraversal(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{
		"dist/index.html": {Data: []byte("<html></html>")},
	})

	_, err := assets.LookupAsset("/../../secret.txt")
	if err == nil {
		t.Fatal("expected traversal error")
	}
}

func TestIsAssetLikePathClassifiesExpectedPaths(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{})

	tests := []struct {
		path string
		want bool
	}{
		{path: "/assets/app.js", want: true},
		{path: "/favicon.ico", want: true},
		{path: "/robots.txt", want: true},
		{path: "/instances/demo", want: false},
		{path: "/instances/v1.2", want: true},
		{path: "/.well-known/test", want: true},
	}

	for _, tt := range tests {
		if got := assets.IsAssetLikePath(tt.path); got != tt.want {
			t.Fatalf("path %q: want %v, got %v", tt.path, tt.want, got)
		}
	}
}
```

- [ ] **Step 2: Run the new package tests to verify failure**

Run: `go test ./internal/webui -v`
Expected: FAIL with undefined `NewFSAssets` and missing asset types.

- [ ] **Step 3: Implement the minimal embedded asset unit**

Create `internal/webui/assets.go` with this shape.

```go
package webui

import (
	"embed"
	"fmt"
	"io/fs"
	"mime"
	"path"
	"strings"
	"time"
)

//go:embed dist/* dist/assets/*
var embeddedFiles embed.FS

type AssetFile struct {
	Name        string
	Content     []byte
	Size        int64
	ModTime     time.Time
	ContentType string
}

type UIAssets interface {
	HasIndex() bool
	LookupAsset(requestPath string) (*AssetFile, error)
	OpenIndex() (*AssetFile, error)
	IsAssetLikePath(requestPath string) bool
}

type fsAssets struct {
	root fs.FS
}

func NewEmbeddedAssets() UIAssets {
	return NewFSAssets(embeddedFiles)
}

func NewFSAssets(root fs.FS) UIAssets {
	return &fsAssets{root: root}
}

func (a *fsAssets) HasIndex() bool {
	file, err := a.readFile("/index.html")
	return err == nil && file != nil
}

func (a *fsAssets) OpenIndex() (*AssetFile, error) {
	return a.readFile("/index.html")
}

func (a *fsAssets) LookupAsset(requestPath string) (*AssetFile, error) {
	return a.readFile(requestPath)
}

func (a *fsAssets) IsAssetLikePath(requestPath string) bool {
	clean := path.Clean("/" + strings.TrimSpace(requestPath))
	if strings.HasPrefix(clean, "/assets/") {
		return true
	}
	base := path.Base(clean)
	if strings.Contains(base, ".") {
		return true
	}
	return false
}

func (a *fsAssets) readFile(requestPath string) (*AssetFile, error) {
	cleanPath, err := normalizeRequestPath(requestPath)
	if err != nil {
		return nil, err
	}
	if cleanPath == ".keep" {
		return nil, nil
	}
	data, err := fs.ReadFile(a.root, path.Join("dist", cleanPath))
	if err != nil {
		if errorsIsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read asset %q: %w", cleanPath, err)
	}
	contentType := mime.TypeByExtension(path.Ext(cleanPath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if cleanPath == "index.html" {
		contentType = "text/html; charset=utf-8"
	}
	if cleanPath == "assets/app.js" {
		contentType = "text/javascript; charset=utf-8"
	}
	return &AssetFile{
		Name:        cleanPath,
		Content:     data,
		Size:        int64(len(data)),
		ModTime:     time.Time{},
		ContentType: contentType,
	}, nil
}

func normalizeRequestPath(requestPath string) (string, error) {
	trimmed := strings.TrimSpace(requestPath)
	if trimmed == "" {
		trimmed = "/"
	}
	clean := path.Clean("/" + trimmed)
	if strings.Contains(clean, "..") {
		return "", fmt.Errorf("invalid asset path %q", requestPath)
	}
	clean = strings.TrimPrefix(clean, "/")
	if clean == "" {
		return "index.html", nil
	}
	return clean, nil
}

func errorsIsNotExist(err error) bool {
	return err != nil && strings.Contains(err.Error(), "does not exist")
}
```

- [ ] **Step 4: Refine implementation to use standard not-exist detection and stable content types**

Before running tests, adjust the implementation to use `errors.Is(err, fs.ErrNotExist)` and explicit type handling for `.js`, `.css`, `.html`, `.svg`, `.json`, `.ico`, and fallback to `mime.TypeByExtension`.

```go
if errors.Is(err, fs.ErrNotExist) {
	return nil, nil
}

switch path.Ext(cleanPath) {
case ".html":
	contentType = "text/html; charset=utf-8"
case ".js":
	contentType = "text/javascript; charset=utf-8"
case ".css":
	contentType = "text/css; charset=utf-8"
case ".svg":
	contentType = "image/svg+xml"
case ".json":
	contentType = "application/json"
case ".ico":
	contentType = "image/x-icon"
}
```

- [ ] **Step 5: Run the package tests to verify pass**

Run: `go test ./internal/webui -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/webui/assets.go internal/webui/assets_test.go internal/webui/dist/.keep
git commit -m "feat: add embedded web UI asset layer"
```

### Task 3: Integrate UI assets into the web server with route tests

**Files:**
- Modify: `internal/web/server.go`
- Modify: `internal/web/server_test.go`
- Test: `internal/web/server_test.go`

- [ ] **Step 1: Write failing server tests for API routing, fallback, and 503 behavior**

Extend `internal/web/server_test.go` with these tests.

```go
package web

import (
	"net/http"
	"net/http/httptest"
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
		"dist/index.html":       {Data: []byte("<html>ui</html>")},
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
```

- [ ] **Step 2: Run the web package tests to verify failure**

Run: `go test ./internal/web -v`
Expected: FAIL because `Server` has no `uiAssets` field and no embedded asset routing.

- [ ] **Step 3: Inject the UI asset dependency into the server**

Update the `Server` struct and constructor in `internal/web/server.go`.

```go
type Server struct {
	config          webtypes.WebConfig
	configPath      string
	router          *gin.Engine
	httpServer      *http.Server
	clusterManager  *cluster.Manager
	resourceManager accountReader
	runtimeAuth     *authsession.AuthManager
	auditLogger     *audit.Logger
	uiAssets        webui.UIAssets
}

func NewServer(config webtypes.WebConfig, configPath string, clusterManager *cluster.Manager, resourceManager accountReader, runtimeAuth *authsession.AuthManager) (*Server, error) {
	// existing setup omitted
	server := &Server{
		config:          config,
		configPath:      configPath,
		router:          router,
		clusterManager:  clusterManager,
		resourceManager: resourceManager,
		runtimeAuth:     runtimeAuth,
		auditLogger:     auditLogger,
		uiAssets:        webui.NewEmbeddedAssets(),
	}
	server.setupRoutes()
	return server, nil
}
```

- [ ] **Step 4: Refactor `setupRoutes` to rebuild a clean router and serve UI assets**

Replace route setup logic so it can be safely called in tests and add helper methods.

```go
func (s *Server) setupRoutes() {
	s.router = gin.New()
	s.router.Use(gin.Recovery())
	s.router.HandleMethodNotAllowed = false
	s.router.RedirectFixedPath = false

	if s.config.CORS.Enabled {
		s.router.Use(s.corsMiddleware())
	}

	api := s.router.Group("/api")
	{
		api.GET("/status", s.handleGetStatus)
		api.GET("/accounts", s.handleGetAccounts)
		api.GET("/accounts/:id", s.handleGetAccount)
		api.GET("/instances", s.handleGetInstances)
		api.GET("/instances/:id", s.handleGetInstance)
		api.POST("/auth/microsoft/init", s.handleMicrosoftAuthInit)
		api.POST("/auth/microsoft/poll", s.handleMicrosoftAuthPoll)
		api.POST("/instances", s.handleCreateInstance)
		api.POST("/instances/:id/start", s.handleStartInstance)
		api.POST("/instances/:id/stop", s.handleStopInstance)
		api.POST("/instances/:id/restart", s.handleRestartInstance)
		api.DELETE("/instances/:id", s.handleDeleteInstance)
		api.POST("/accounts", s.handleCreateAccount)
		api.DELETE("/accounts/:id", s.handleDeleteAccount)
		api.GET("/logs/operations", s.handleGetOperationLogs)
	}

	s.router.NoRoute(s.handleNoRoute)
}

func (s *Server) handleNoRoute(c *gin.Context) {
	path := c.Request.URL.Path
	if isAPIPath(path) {
		c.JSON(http.StatusNotFound, webtypes.OperationResponse{Success: false, Error: "API endpoint not found"})
		return
	}
	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
		c.String(http.StatusNotFound, "Not Found")
		return
	}
	if file, err := s.uiAssets.LookupAsset(path); err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	} else if file != nil {
		serveAsset(c, file)
		return
	}
	if s.uiAssets.IsAssetLikePath(path) {
		c.String(http.StatusNotFound, "Not Found")
		return
	}
	if !s.uiAssets.HasIndex() {
		serveFrontendUnavailable(c)
		return
	}
	file, err := s.uiAssets.OpenIndex()
	if err != nil || file == nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	serveAsset(c, file)
}
```

- [ ] **Step 5: Add small helpers for API-space checks and asset responses**

Append helpers in `internal/web/server.go`.

```go
func isAPIPath(requestPath string) bool {
	return requestPath == "/api" || requestPath == "/api/" || strings.HasPrefix(requestPath, "/api/")
}

func serveAsset(c *gin.Context, file *webui.AssetFile) {
	if file.ContentType != "" {
		c.Header("Content-Type", file.ContentType)
	}
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}
	c.Data(http.StatusOK, file.ContentType, file.Content)
}

func serveFrontendUnavailable(c *gin.Context) {
	body := "<html><body><h1>Frontend unavailable</h1><p>前端尚未构建，当前服务仅提供 API。</p></body></html>"
	c.Header("Content-Type", "text/html; charset=utf-8")
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	c.Data(http.StatusServiceUnavailable, "text/html; charset=utf-8", []byte(body))
}
```

- [ ] **Step 6: Run the web package tests to verify pass**

Run: `go test ./internal/web -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/web/server.go internal/web/server_test.go internal/webui/assets.go internal/webui/assets_test.go
git commit -m "feat: serve embedded frontend assets when available"
```

### Task 4: Log frontend availability at startup

**Files:**
- Modify: `cmd/gmcc/main.go`

- [ ] **Step 1: Write a failing main-package test for startup mode logging helper**

Add a focused helper test in `cmd/gmcc/main_test.go`.

```go
func TestDescribeWebUIMode(t *testing.T) {
	tests := []struct {
		name string
		hasUI bool
		want string
	}{
		{name: "with ui", hasUI: true, want: "embedded"},
		{name: "without ui", hasUI: false, want: "api-only"},
	}

	for _, tt := range tests {
		if got := describeWebUIMode(tt.hasUI); got != tt.want {
			t.Fatalf("%s: want %q, got %q", tt.name, tt.want, got)
		}
	}
}
```

- [ ] **Step 2: Run the main-package tests to verify failure**

Run: `go test ./cmd/gmcc -v`
Expected: FAIL with undefined `describeWebUIMode`.

- [ ] **Step 3: Add the helper and startup log line**

Update `cmd/gmcc/main.go`.

```go
func describeWebUIMode(hasUI bool) string {
	if hasUI {
		return "embedded"
	}
	return "api-only"
}
```

And after creating the server:

```go
logx.Infof("Web UI 模式: %s", describeWebUIMode(server.HasEmbeddedUI()))
```

If needed, add a tiny method on `internal/web.Server`:

```go
func (s *Server) HasEmbeddedUI() bool {
	return s.uiAssets != nil && s.uiAssets.HasIndex()
}
```

- [ ] **Step 4: Run the main-package tests to verify pass**

Run: `go test ./cmd/gmcc -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/gmcc/main.go cmd/gmcc/main_test.go internal/web/server.go
git commit -m "chore: log web UI startup mode"
```

### Task 5: Add Go packager with whitelist behavior tests

**Files:**
- Create: `tools/packager/main.go`
- Create: `tools/packager/main_test.go`
- Create: `tools/packager/testdata/frontend-dist/index.html`
- Create: `tools/packager/testdata/frontend-dist/assets/app.js`
- Create: `tools/packager/testdata/frontend-dist/assets/app.js.map`
- Create: `tools/packager/testdata/frontend-dist/favicon.ico`
- Create: `tools/packager/testdata/frontend-dist/robots.txt`
- Create: `tools/packager/testdata/frontend-dist/extra.txt`
- Test: `tools/packager/main_test.go`

- [ ] **Step 1: Write failing tests for cleanup and whitelist copying**

Create `tools/packager/main_test.go` with these tests.

```go
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareEmbedDirWithoutFrontendDistKeepsOnlyKeepFile(t *testing.T) {
	tmp := t.TempDir()
	embedDir := filepath.Join(tmp, "embed")
	if err := os.MkdirAll(embedDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(embedDir, ".keep"), []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(embedDir, "old.txt"), []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	prepared, err := prepareEmbedDir(filepath.Join(tmp, "missing-frontend-dist"), embedDir)
	if err != nil {
		t.Fatal(err)
	}
	if prepared {
		t.Fatal("expected prepared to be false")
	}
	if _, err := os.Stat(filepath.Join(embedDir, ".keep")); err != nil {
		t.Fatalf("expected .keep to remain: %v", err)
	}
	if _, err := os.Stat(filepath.Join(embedDir, "old.txt")); !os.IsNotExist(err) {
		t.Fatalf("expected old artifact removed, got %v", err)
	}
}

func TestPrepareEmbedDirCopiesOnlyWhitelistedFiles(t *testing.T) {
	tmp := t.TempDir()
	frontendDist := filepath.Join("testdata", "frontend-dist")
	embedDir := filepath.Join(tmp, "embed")

	prepared, err := prepareEmbedDir(frontendDist, embedDir)
	if err != nil {
		t.Fatal(err)
	}
	if !prepared {
		t.Fatal("expected prepared to be true")
	}

	wantExists := []string{
		".keep",
		"index.html",
		filepath.Join("assets", "app.js"),
		"favicon.ico",
		"robots.txt",
	}
	for _, rel := range wantExists {
		if _, err := os.Stat(filepath.Join(embedDir, rel)); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}

	wantMissing := []string{
		filepath.Join("assets", "app.js.map"),
		"extra.txt",
	}
	for _, rel := range wantMissing {
		if _, err := os.Stat(filepath.Join(embedDir, rel)); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be absent, got %v", rel, err)
		}
	}
}

func TestPrepareEmbedDirWithoutIndexFallsBackToAPIOnly(t *testing.T) {
	tmp := t.TempDir()
	frontendDist := filepath.Join(tmp, "frontend-dist")
	embedDir := filepath.Join(tmp, "embed")
	if err := os.MkdirAll(filepath.Join(frontendDist, "assets"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(frontendDist, "assets", "app.js"), []byte("console.log('ok')"), 0644); err != nil {
		t.Fatal(err)
	}

	prepared, err := prepareEmbedDir(frontendDist, embedDir)
	if err != nil {
		t.Fatal(err)
	}
	if prepared {
		t.Fatal("expected prepared to be false without index.html")
	}
	if _, err := os.Stat(filepath.Join(embedDir, ".keep")); err != nil {
		t.Fatalf("expected .keep to remain: %v", err)
	}
	if _, err := os.Stat(filepath.Join(embedDir, "assets", "app.js")); !os.IsNotExist(err) {
		t.Fatalf("expected copied assets to be cleared when index.html is missing, got %v", err)
	}
}
```

- [ ] **Step 2: Add the test fixtures**

Create these files.

`tools/packager/testdata/frontend-dist/index.html`

```html
<!doctype html>
<html><body>ok</body></html>
```

`tools/packager/testdata/frontend-dist/assets/app.js`

```js
console.log("ok")
```

`tools/packager/testdata/frontend-dist/assets/app.js.map`

```json
{"version":3}
```

`tools/packager/testdata/frontend-dist/favicon.ico`

```text
ico
```

`tools/packager/testdata/frontend-dist/robots.txt`

```text
User-agent: *
Disallow:
```

`tools/packager/testdata/frontend-dist/extra.txt`

```text
skip me
```

- [ ] **Step 3: Run the packager tests to verify failure**

Run: `go test ./tools/packager -v`
Expected: FAIL with undefined `prepareEmbedDir`.

- [ ] **Step 4: Implement the packager core and CLI**

Create `tools/packager/main.go` with this structure.

```go
package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	frontendDist := flag.String("frontend-dist", filepath.Join("frontend", "dist"), "frontend dist input")
	embedDist := flag.String("embed-dist", filepath.Join("internal", "webui", "dist"), "embed dist output")
	output := flag.String("output", defaultOutputPath(), "binary output path")
	flag.Parse()

	prepared, err := prepareEmbedDir(*frontendDist, *embedDist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "prepare embed dir: %v\n", err)
		os.Exit(1)
	}
	if prepared {
		fmt.Println("frontend assets prepared for embedding")
	} else {
		fmt.Println("frontend assets unavailable; building API-only binary")
	}
	if err := buildBinary(*output); err != nil {
		fmt.Fprintf(os.Stderr, "build binary: %v\n", err)
		os.Exit(1)
	}
}

func defaultOutputPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join("build", "gmcc.exe")
	}
	return filepath.Join("build", "gmcc")
}

func prepareEmbedDir(frontendDist, embedDir string) (bool, error) {
	if err := ensureCleanEmbedDir(embedDir); err != nil {
		return false, err
	}
	if _, err := os.Stat(frontendDist); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if err := copyWhitelistedFiles(frontendDist, embedDir); err != nil {
		return false, err
	}
	if _, err := os.Stat(filepath.Join(embedDir, "index.html")); err != nil {
		if os.IsNotExist(err) {
			if err := ensureCleanEmbedDir(embedDir); err != nil {
				return false, err
			}
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ensureCleanEmbedDir(embedDir string) error {
	if err := os.MkdirAll(embedDir, 0755); err != nil {
		return err
	}
	entries, err := os.ReadDir(embedDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name() == ".keep" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(embedDir, entry.Name())); err != nil {
			return err
		}
	}
	keepPath := filepath.Join(embedDir, ".keep")
	if _, err := os.Stat(keepPath); os.IsNotExist(err) {
		return os.WriteFile(keepPath, []byte{}, 0644)
	}
	return nil
}

func copyWhitelistedFiles(frontendDist, embedDir string) error {
	rootFiles := map[string]bool{
		"index.html":           true,
		"favicon.ico":          true,
		"favicon.svg":          true,
		"manifest.webmanifest": true,
		"robots.txt":           true,
	}
	return filepath.WalkDir(frontendDist, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(frontendDist, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, "assets/") {
			if strings.HasSuffix(rel, ".map") || strings.HasPrefix(filepath.Base(rel), ".") {
				return nil
			}
			return copyFile(path, filepath.Join(embedDir, filepath.FromSlash(rel)))
		}
		if rootFiles[rel] {
			return copyFile(path, filepath.Join(embedDir, filepath.FromSlash(rel)))
		}
		return nil
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func buildBinary(output string) error {
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	cmd := exec.Command("go", "build", "-o", output, "./cmd/gmcc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
```

- [ ] **Step 5: Run the packager tests to verify pass**

Run: `go test ./tools/packager -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add tools/packager/main.go tools/packager/main_test.go tools/packager/testdata/frontend-dist
git commit -m "feat: add frontend embedding packager"
```

### Task 6: Document the new workflow

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README text for optional embedded frontend support**

Adjust the runtime-mode section and build instructions in `README.md`.

```md
- 当前服务默认仍可作为 API-only 模式运行
- 当二进制包含已嵌入的前端静态资源时，同一进程可同时提供页面和 `/api`
- 当未包含前端资源时，非 `/api` 页面请求返回 `503`，`/api` 继续可用
```

Add a packaging example near the build section.

```bash
go run ./tools/packager
```

Add a short note near the repo layout.

```md
- `frontend/`: 预留给未来 Vue/Node 前端工程
- `internal/webui/dist/`: 由打包脚本准备的嵌入静态资源目录
```

- [ ] **Step 2: Verify the README still matches actual behavior**

Run: `grep` mentally against the implementation and confirm the text does not claim automatic frontend scaffolding or mandatory frontend assets.
Expected: README describes optional embedded UI support and API-only fallback accurately.

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: describe optional embedded frontend workflow"
```

### Task 7: Full verification

**Files:**
- Verify: `internal/webui/assets.go`
- Verify: `internal/web/server.go`
- Verify: `tools/packager/main.go`
- Verify: `README.md`

- [ ] **Step 1: Format the repository**

Run: `go fmt ./...`
Expected: files are formatted without errors.

- [ ] **Step 2: Run focused package tests**

Run: `go test ./internal/webui ./internal/web ./tools/packager ./cmd/gmcc -v`
Expected: PASS

- [ ] **Step 3: Run the full test suite**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 4: Build the main binary directly**

Run: `go build -o gmcc.exe ./cmd/gmcc`
Expected: build succeeds.

- [ ] **Step 5: Verify the packager end-to-end without frontend assets**

Run: `go run ./tools/packager -frontend-dist frontend/dist -output build/gmcc-packaged.exe`
Expected: log mentions API-only build and `build/gmcc-packaged.exe` is created.

- [ ] **Step 6: Verify the packager falls back to API-only when `index.html` is missing**

Run: `go run ./tools/packager -frontend-dist tools/packager/testdata/frontend-dist/assets -output build/gmcc-packaged-no-ui.exe`
Expected: log mentions API-only build, `internal/webui/dist/` is left with `.keep`, and `build/gmcc-packaged-no-ui.exe` is created.

- [ ] **Step 7: Commit**

```bash
git add .
git commit -m "test: verify frontend embedding workflow"
```

## Self-Review

- Spec coverage: directory reservation, embed placeholder, `internal/webui`, route behavior, startup mode logging, Go packager, `.gitignore`, README, and verification tasks are all mapped.
- Placeholder scan: no `TODO`/`TBD` markers remain in the plan.
- Type consistency: the plan uses `UIAssets`, `AssetFile`, `NewEmbeddedAssets`, `NewFSAssets`, `HasEmbeddedUI`, and `prepareEmbedDir` consistently.
