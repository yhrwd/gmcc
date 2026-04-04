package cluster

import (
	"context"
	"errors"
	"sync"
	"testing"

	"gmcc/internal/logx"
)

type fakeRunner struct {
	mu       sync.Mutex
	ready    bool
	runErr   error
	blockRun bool
	runCh    chan struct{}
	stopCh   chan struct{}
}

func newFakeRunner() *fakeRunner {
	return &fakeRunner{
		runCh:  make(chan struct{}, 1),
		stopCh: make(chan struct{}, 1),
	}
}

func (f *fakeRunner) Run(ctx context.Context) error {
	f.mu.Lock()
	runErr := f.runErr
	block := f.blockRun
	f.mu.Unlock()

	select {
	case f.runCh <- struct{}{}:
	default:
	}

	if runErr != nil {
		return runErr
	}

	if block {
		<-ctx.Done()
		select {
		case f.stopCh <- struct{}{}:
		default:
		}
		return ctx.Err()
	}

	<-ctx.Done()
	select {
	case f.stopCh <- struct{}{}:
	default:
	}
	return ctx.Err()
}

func (f *fakeRunner) IsReady() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.ready
}

func (f *fakeRunner) SetReady(v bool) {
	f.mu.Lock()
	f.ready = v
	f.mu.Unlock()
}

var errFakeNetworkEOF = errors.New("fake network eof")

var testLogInitOnce sync.Once

func initTestLogger(t *testing.T) {
	t.Helper()
	testLogInitOnce.Do(func() {
		if err := logx.Init(t.TempDir(), false, 0, false); err != nil {
			t.Fatalf("init logger: %v", err)
		}
	})
}
