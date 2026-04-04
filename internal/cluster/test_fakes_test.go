package cluster

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	authsession "gmcc/internal/auth/session"
	"gmcc/internal/config"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient"
)

type fakeRunner struct {
	mu       sync.Mutex
	ready    bool
	runErr   error
	runDelay time.Duration
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
		if f.runDelay > 0 {
			time.Sleep(f.runDelay)
		}
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

func (f *fakeRunner) SendCommand(_ string) error {
	return nil
}

func (f *fakeRunner) GetPlayer() *mcclient.Player {
	return nil
}

var errFakeNetworkEOF = errors.New("fake network eof")

type scriptedRunnerFactory struct {
	mu       sync.Mutex
	outcomes []runnerOutcome
	built    []*fakeRunner
}

type runnerOutcome struct {
	err      error
	ready    bool
	runDelay time.Duration
}

func newScriptedRunnerFactory(outcomes ...runnerOutcome) *scriptedRunnerFactory {
	return &scriptedRunnerFactory{outcomes: outcomes}
}

func (f *scriptedRunnerFactory) Build(_ *config.Config, _ *authsession.AuthManager) runner {
	f.mu.Lock()
	defer f.mu.Unlock()

	r := newFakeRunner()
	idx := len(f.built)
	if idx < len(f.outcomes) {
		outcome := f.outcomes[idx]
		r.runErr = outcome.err
		r.ready = outcome.ready
		r.runDelay = outcome.runDelay
	}
	f.built = append(f.built, r)
	return r
}

func (f *scriptedRunnerFactory) BuildCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.built)
}

var testLogInitOnce sync.Once

func initTestLogger(t *testing.T) {
	t.Helper()
	testLogInitOnce.Do(func() {
		if err := logx.Init(t.TempDir(), false, 0, false); err != nil {
			t.Fatalf("init logger: %v", err)
		}
	})
}
