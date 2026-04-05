package logx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEventWriterRotatesAtConfiguredSize(t *testing.T) {
	dir := t.TempDir()

	w, err := newStructuredWriter(filepath.Join(dir, "gmcc-events.jsonl"), 128, 2)
	if err != nil {
		t.Fatalf("new structured writer: %v", err)
	}
	defer w.Close()

	for i := 0; i < 20; i++ {
		if _, err := w.Write([]byte("{\"event_type\":\"instance.lifecycle\",\"action\":\"ready\"}\n")); err != nil {
			t.Fatalf("write event: %v", err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}

	if len(entries) < 2 {
		t.Fatalf("expected rotated files, got %d", len(entries))
	}
	if len(entries) > 3 {
		t.Fatalf("expected writer to retain bounded archives, got %d files", len(entries))
	}
}
