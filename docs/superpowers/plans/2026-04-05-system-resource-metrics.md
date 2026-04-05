# System Resource Metrics Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a new `GET /api/resources` endpoint that returns cross-platform host CPU and memory metrics for the gmcc dashboard.

**Architecture:** Add a focused `internal/systemmetrics` package that wraps a cross-platform library behind a tiny collector interface and returns a stable `Snapshot` model. Inject that collector into the existing web server, expose a new read-only API handler, and document the endpoint in the frontend-facing API docs.

**Tech Stack:** Go, Gin, `github.com/shirou/gopsutil/v4`, Go testing package, standard library time/json/http

---

## File Map

- Create: `internal/systemmetrics/collector.go`
- Create: `internal/systemmetrics/collector_test.go`
- Modify: `internal/webtypes/types.go`
- Modify: `internal/web/server.go`
- Modify: `internal/web/handlers.go`
- Modify: `internal/web/handlers_test.go`
- Modify: `cmd/gmcc/main.go`
- Modify: `docs/api.md`
- Modify: `go.mod`
- Modify: `go.sum`

## Implementation Notes

- Keep the endpoint at `GET /api/resources`.
- Use a tiny internal collector interface; Web code must not depend directly on `gopsutil` types.
- `cpu_percent` must mean normalized whole-machine usage in the `0-100` range.
- Return complete snapshots only; if CPU or memory collection fails, return `500` instead of partial data.
- `collected_at` must be emitted as a UTC `RFC3339` timestamp.
- Keep this endpoint outside `logOperation`; it is a read-only status route like the other list/read APIs.

### Task 1: Add the system metrics collector package

**Files:**
- Create: `internal/systemmetrics/collector.go`
- Create: `internal/systemmetrics/collector_test.go`
- Modify: `go.mod`
- Modify: `go.sum`

- [ ] **Step 1: Add failing tests for snapshot shape and failure propagation**

Create `internal/systemmetrics/collector_test.go` with these tests.

```go
package systemmetrics

import (
	"errors"
	"testing"
	"time"
)

type fakeCPUReader struct {
	percent float64
	err     error
}

func (f fakeCPUReader) ReadCPUPercent() (float64, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.percent, nil
}

type fakeMemoryReader struct {
	snapshot MemorySnapshot
	err      error
}

func (f fakeMemoryReader) ReadMemory() (MemorySnapshot, error) {
	if f.err != nil {
		return MemorySnapshot{}, f.err
	}
	return f.snapshot, nil
}

func TestCollectorCollectReturnsSnapshot(t *testing.T) {
	collector := NewCollector(fakeCPUReader{percent: 12.5}, fakeMemoryReader{snapshot: MemorySnapshot{
		TotalBytes:     16,
		UsedBytes:      8,
		AvailableBytes: 7,
		UsedPercent:    50,
	}})

	snapshot, err := collector.Collect()
	if err != nil {
		t.Fatalf("collect failed: %v", err)
	}
	if snapshot.CPUPercent != 12.5 {
		t.Fatalf("want cpu percent 12.5, got %v", snapshot.CPUPercent)
	}
	if snapshot.Memory.TotalBytes != 16 {
		t.Fatalf("unexpected memory snapshot: %+v", snapshot.Memory)
	}
	if snapshot.CollectedAt.IsZero() {
		t.Fatal("expected collected_at to be set")
	}
	if snapshot.CollectedAt.Location() != time.UTC {
		t.Fatalf("expected UTC collected_at, got %v", snapshot.CollectedAt.Location())
	}
}

func TestCollectorCollectReturnsErrorWhenCPUFails(t *testing.T) {
	collector := NewCollector(fakeCPUReader{err: errors.New("cpu failed")}, fakeMemoryReader{})

	_, err := collector.Collect()
	if err == nil {
		t.Fatal("expected collect error")
	}
}

func TestCollectorCollectReturnsErrorWhenMemoryFails(t *testing.T) {
	collector := NewCollector(fakeCPUReader{percent: 12.5}, fakeMemoryReader{err: errors.New("memory failed")})

	_, err := collector.Collect()
	if err == nil {
		t.Fatal("expected collect error")
	}
}
```

- [ ] **Step 2: Run the new package tests to verify failure**

Run: `go test ./internal/systemmetrics -v`
Expected: FAIL with undefined `NewCollector`, `MemorySnapshot`, and collector types.

- [ ] **Step 3: Add the gopsutil dependency**

Update `go.mod` by adding the dependency below.

```go
require github.com/shirou/gopsutil/v4 v4.25.3
```

Then resolve `go.sum` with:

```bash
go mod tidy
```

- [ ] **Step 4: Implement the collector package**

Create `internal/systemmetrics/collector.go` with this structure.

```go
package systemmetrics

import (
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type MemorySnapshot struct {
	TotalBytes     uint64
	UsedBytes      uint64
	AvailableBytes uint64
	UsedPercent    float64
}

type Snapshot struct {
	CPUPercent  float64
	Memory      MemorySnapshot
	CollectedAt time.Time
}

type Collector interface {
	Collect() (Snapshot, error)
}

type CPUReader interface {
	ReadCPUPercent() (float64, error)
}

type MemoryReader interface {
	ReadMemory() (MemorySnapshot, error)
}

type collector struct {
	cpuReader    CPUReader
	memoryReader MemoryReader
}

func NewCollector(cpuReader CPUReader, memoryReader MemoryReader) Collector {
	return &collector{cpuReader: cpuReader, memoryReader: memoryReader}
}

func NewDefaultCollector() Collector {
	return NewCollector(gopsutilCPUReader{}, gopsutilMemoryReader{})
}

func (c *collector) Collect() (Snapshot, error) {
	cpuPercent, err := c.cpuReader.ReadCPUPercent()
	if err != nil {
		return Snapshot{}, fmt.Errorf("read cpu percent: %w", err)
	}
	memorySnapshot, err := c.memoryReader.ReadMemory()
	if err != nil {
		return Snapshot{}, fmt.Errorf("read memory snapshot: %w", err)
	}
	return Snapshot{
		CPUPercent:  normalizePercent(cpuPercent),
		Memory:      memorySnapshot,
		CollectedAt: time.Now().UTC(),
	}, nil
	}

type gopsutilCPUReader struct{}

func (g gopsutilCPUReader) ReadCPUPercent() (float64, error) {
	values, err := cpu.Percent(200*time.Millisecond, false)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("empty cpu percent result")
	}
	return values[0], nil
}

type gopsutilMemoryReader struct{}

func (g gopsutilMemoryReader) ReadMemory() (MemorySnapshot, error) {
	stats, err := mem.VirtualMemory()
	if err != nil {
		return MemorySnapshot{}, err
	}
	return MemorySnapshot{
		TotalBytes:     stats.Total,
		UsedBytes:      stats.Used,
		AvailableBytes: stats.Available,
		UsedPercent:    normalizePercent(stats.UsedPercent),
	}, nil
}

func normalizePercent(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return math.Round(value*100) / 100
}
```

- [ ] **Step 5: Run the package tests to verify pass**

Run: `go test ./internal/systemmetrics -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum internal/systemmetrics/collector.go internal/systemmetrics/collector_test.go
git commit -m "feat: add system resource metrics collector"
```

### Task 2: Add the API response models

**Files:**
- Modify: `internal/webtypes/types.go`

- [ ] **Step 1: Write the failing JSON shape test in the web handler test file first**

Before editing models, add this assertion block to the future success test in `internal/web/handlers_test.go` so the JSON shape is locked in.

```go
if _, ok := payload["cpu_percent"].(float64); !ok {
	t.Fatalf("expected cpu_percent number, got %+v", payload)
}
memory, ok := payload["memory"].(map[string]any)
if !ok {
	t.Fatalf("expected memory object, got %+v", payload)
}
if _, ok := memory["total_bytes"].(float64); !ok {
	t.Fatalf("expected total_bytes number, got %+v", memory)
}
if _, ok := payload["collected_at"].(string); !ok {
	t.Fatalf("expected collected_at string, got %+v", payload)
}
```

- [ ] **Step 2: Add the response structs to `internal/webtypes/types.go`**

Append these types.

```go
type ResourceMemoryView struct {
	TotalBytes     uint64  `json:"total_bytes"`
	UsedBytes      uint64  `json:"used_bytes"`
	AvailableBytes uint64  `json:"available_bytes"`
	UsedPercent    float64 `json:"used_percent"`
}

type ResourceSnapshotView struct {
	CPUPercent  float64            `json:"cpu_percent"`
	Memory      ResourceMemoryView `json:"memory"`
	CollectedAt time.Time          `json:"collected_at"`
}
```

- [ ] **Step 3: Run the web package tests to verify the repository still compiles**

Run: `go test ./internal/web -run TestHandleGetResources -v`
Expected: FAIL because the handler does not exist yet, but the model types compile.

- [ ] **Step 4: Commit**

```bash
git add internal/webtypes/types.go internal/web/handlers_test.go
git commit -m "feat: add resource metrics response models"
```

### Task 3: Wire the collector into the web server and add the handler

**Files:**
- Modify: `internal/web/server.go`
- Modify: `internal/web/handlers.go`
- Modify: `internal/web/handlers_test.go`
- Modify: `cmd/gmcc/main.go`

- [ ] **Step 1: Add failing handler tests for success, missing collector, and collector failure**

Append these test helpers and tests to `internal/web/handlers_test.go`.

```go
type fakeResourceCollector struct {
	snapshot systemmetrics.Snapshot
	err      error
}

func (f fakeResourceCollector) Collect() (systemmetrics.Snapshot, error) {
	if f.err != nil {
		return systemmetrics.Snapshot{}, f.err
	}
	return f.snapshot, nil
}

func TestHandleGetResourcesReturnsSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := &Server{
		resourceCollector: fakeResourceCollector{snapshot: systemmetrics.Snapshot{
			CPUPercent: 12.5,
			Memory: systemmetrics.MemorySnapshot{
				TotalBytes:     16,
				UsedBytes:      8,
				AvailableBytes: 7,
				UsedPercent:    50,
			},
			CollectedAt: time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC),
		}},
		auditLogger: newTestAuditLogger(t),
	}

	ctx, recorder := newJSONContext(http.MethodGet, "/api/resources")
	server.handleGetResources(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("want 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if got := payload["cpu_percent"]; got != 12.5 {
		t.Fatalf("unexpected cpu_percent: %+v", payload)
	}
	if _, ok := payload["cpu_percent"].(float64); !ok {
		t.Fatalf("expected cpu_percent number, got %+v", payload)
	}
	memory, ok := payload["memory"].(map[string]any)
	if !ok {
		t.Fatalf("expected memory object, got %+v", payload)
	}
	if _, ok := memory["total_bytes"].(float64); !ok {
		t.Fatalf("expected total_bytes number, got %+v", memory)
	}
	if collectedAt, ok := payload["collected_at"].(string); !ok || collectedAt != "2026-04-05T12:00:00Z" {
		t.Fatalf("unexpected collected_at: %+v", payload)
	}
}

func TestHandleGetResourcesReturns503WithoutCollector(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := &Server{auditLogger: newTestAuditLogger(t)}

	ctx, recorder := newJSONContext(http.MethodGet, "/api/resources")
	server.handleGetResources(ctx)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("want 503, got %d", recorder.Code)
	}
}

func TestHandleGetResourcesReturns500OnCollectFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := &Server{
		resourceCollector: fakeResourceCollector{err: errors.New("boom")},
		auditLogger:       newTestAuditLogger(t),
	}

	ctx, recorder := newJSONContext(http.MethodGet, "/api/resources")
	server.handleGetResources(ctx)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("want 500, got %d", recorder.Code)
	}
}
```

- [ ] **Step 2: Run the targeted web tests to verify failure**

Run: `go test ./internal/web -run TestHandleGetResources -v`
Expected: FAIL with missing `resourceCollector` field and undefined handler.

- [ ] **Step 3: Add the collector dependency to `Server` and register the route**

Update `internal/web/server.go`.

```go
type Server struct {
	config            webtypes.WebConfig
	configPath        string
	router            *gin.Engine
	httpServer        *http.Server
	clusterManager    *cluster.Manager
	resourceManager   accountReader
	runtimeAuth       *authsession.AuthManager
	auditLogger       *audit.Logger
	uiAssets          webui.UIAssets
	resourceCollector systemmetrics.Collector
}

func NewServer(config webtypes.WebConfig, configPath string, clusterManager *cluster.Manager, resourceManager accountReader, runtimeAuth *authsession.AuthManager, resourceCollector systemmetrics.Collector) (*Server, error) {
	// existing setup omitted
	server := &Server{
		config:            config,
		configPath:        configPath,
		clusterManager:    clusterManager,
		resourceManager:   resourceManager,
		runtimeAuth:       runtimeAuth,
		auditLogger:       auditLogger,
		uiAssets:          webui.NewEmbeddedAssets(),
		resourceCollector: resourceCollector,
	}
	server.setupRoutes()
	return server, nil
}
```

And register the route inside `setupRoutes()`.

```go
api.GET("/resources", s.handleGetResources)
```

- [ ] **Step 4: Implement `handleGetResources` in `internal/web/handlers.go`**

Add the handler and mapping helper.

```go
func (s *Server) handleGetResources(c *gin.Context) {
	if s.resourceCollector == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "resource metrics collector not initialized"})
		return
	}

	snapshot, err := s.resourceCollector.Collect()
	if err != nil {
		logx.Warnf("采集系统资源失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect system resources"})
		return
	}

	c.JSON(http.StatusOK, webtypes.ResourceSnapshotView{
		CPUPercent: snapshot.CPUPercent,
		Memory: webtypes.ResourceMemoryView{
			TotalBytes:     snapshot.Memory.TotalBytes,
			UsedBytes:      snapshot.Memory.UsedBytes,
			AvailableBytes: snapshot.Memory.AvailableBytes,
			UsedPercent:    snapshot.Memory.UsedPercent,
		},
		CollectedAt: snapshot.CollectedAt,
	})
}
```

- [ ] **Step 5: Inject the default collector from `cmd/gmcc/main.go`**

Update the server construction call.

```go
server, err := web.NewServer(
	runtime.WebConfig,
	configPath,
	runtime.ClusterManager,
	runtime.ResourceManager,
	runtime.AuthManager,
	systemmetrics.NewDefaultCollector(),
)
```

Add the import.

```go
"gmcc/internal/systemmetrics"
```

- [ ] **Step 6: Update any existing `NewServer(...)` call sites in tests**

For current tests that call `NewServer(...)`, pass `nil` for the new collector dependency where resource metrics are not under test.

```go
server, err := NewServer(webtypes.WebConfig{}, configPath, cluster.NewManager(cluster.ClusterConfig{}, nil), nil, nil, nil)
```

- [ ] **Step 7: Run the focused web and main tests to verify pass**

Run: `go test ./internal/web ./cmd/gmcc -v`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add internal/web/server.go internal/web/handlers.go internal/web/handlers_test.go cmd/gmcc/main.go internal/webtypes/types.go
git commit -m "feat: expose system resource metrics API"
```

### Task 4: Document the new endpoint

**Files:**
- Modify: `docs/api.md`

- [ ] **Step 1: Add the new API documentation section**

Insert a new subsection under status/query APIs.

```md
### 3.6 获取宿主机系统资源

`GET /api/resources`

响应示例：

```json
{
  "cpu_percent": 12.5,
  "memory": {
    "total_bytes": 17179869184,
    "used_bytes": 8589934592,
    "available_bytes": 7516192768,
    "used_percent": 50.0
  },
  "collected_at": "2026-04-05T12:00:00Z"
}
```

字段：

- `cpu_percent`: 宿主机整体 CPU 使用率（0-100）
- `memory.total_bytes`: 总内存字节数
- `memory.used_bytes`: 已使用内存字节数
- `memory.available_bytes`: 可用内存字节数
- `memory.used_percent`: 内存使用率（0-100）
- `collected_at`: 采样时间，UTC RFC3339 时间戳

错误：

- `503 {"error":"resource metrics collector not initialized"}`
- `500 {"error":"failed to collect system resources"}`
```

- [ ] **Step 2: Verify the docs match the implemented route and payload**

Run: compare the docs section against `internal/web/handlers.go` and `internal/webtypes/types.go`.
Expected: route, field names, and error messages match exactly.

- [ ] **Step 3: Commit**

```bash
git add docs/api.md
git commit -m "docs: add system resource metrics API"
```

### Task 5: Full verification

**Files:**
- Verify: `internal/systemmetrics/collector.go`
- Verify: `internal/web/server.go`
- Verify: `internal/web/handlers.go`
- Verify: `docs/api.md`

- [ ] **Step 1: Format all changed packages**

Run: `go fmt ./internal/systemmetrics ./internal/web ./internal/webtypes ./cmd/gmcc`
Expected: formatting succeeds.

- [ ] **Step 2: Run focused tests**

Run: `go test ./internal/systemmetrics ./internal/web ./cmd/gmcc -v`
Expected: PASS

- [ ] **Step 3: Run the full test suite**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 4: Build the main binary**

Run: `go build -o gmcc.exe ./cmd/gmcc`
Expected: build succeeds.

- [ ] **Step 5: Commit**

```bash
git add .
git commit -m "test: verify system resource metrics endpoint"
```

## Self-Review

- Spec coverage: route, cross-platform collector package, gopsutil dependency, UTC/RFC3339 timestamp behavior, complete-failure semantics, web injection, docs, and tests are all covered.
- Placeholder scan: no `TODO`/`TBD` markers remain in the plan.
- Type consistency: the plan uses `Collector`, `Snapshot`, `MemorySnapshot`, `ResourceSnapshotView`, and `handleGetResources` consistently.
