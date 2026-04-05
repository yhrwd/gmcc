package web

import (
	"os"
	"path/filepath"
	"testing"

	"gmcc/internal/cluster"
	"gmcc/internal/webtypes"
)

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
