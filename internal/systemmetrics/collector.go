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
