package logx

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSummaryDoesNotWriteDebugOnlyDetails(t *testing.T) {
	resetLogxStateForTest(t)

	var console bytes.Buffer
	consoleLogger = log.New(&console, "", 0)

	dir := t.TempDir()
	writer, err := newBoundedRotatingWriter(filepath.Join(dir, "gmcc.log"), 1024, 2)
	if err != nil {
		t.Fatalf("new bounded writer: %v", err)
	}
	defer writer.Close()

	fileWriter = writer
	fileLogger = log.New(fileWriter, "", 0)
	debugEnabled = false

	Summaryf("info", "instance ready")
	Debugf("provider detail: %s", "device code pending")

	consoleOutput := console.String()
	if !strings.Contains(consoleOutput, "[INFO] instance ready") {
		t.Fatalf("expected summary output in console, got %q", consoleOutput)
	}
	if strings.Contains(consoleOutput, "provider detail") {
		t.Fatalf("expected debug detail to stay hidden, got %q", consoleOutput)
	}

	data, err := os.ReadFile(filepath.Join(dir, "gmcc.log"))
	if err != nil {
		t.Fatalf("read summary log: %v", err)
	}
	fileOutput := string(data)
	if !strings.Contains(fileOutput, "[INFO] instance ready") {
		t.Fatalf("expected summary output in file, got %q", fileOutput)
	}
	if strings.Contains(fileOutput, "provider detail") {
		t.Fatalf("expected debug detail to stay out of summary file, got %q", fileOutput)
	}
}

func TestInfofRemainsBackwardCompatible(t *testing.T) {
	resetLogxStateForTest(t)

	var console bytes.Buffer
	consoleLogger = log.New(&console, "", 0)

	Infof("legacy info message: %s", "ok")

	if got := console.String(); !strings.Contains(got, "[INFO] legacy info message: ok") {
		t.Fatalf("expected Infof to log summary-compatible output, got %q", got)
	}
}

func resetLogxStateForTest(t *testing.T) {
	t.Helper()

	mu.Lock()
	defer mu.Unlock()

	consoleLogger = nil
	fileLogger = nil
	debugEnabled = false

	if eventWriter != nil {
		_ = eventWriter.Close()
		eventWriter = nil
	}
	if fileWriter != nil {
		_ = fileWriter.Close()
		fileWriter = nil
	}
}
