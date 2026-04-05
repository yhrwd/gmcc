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
