package container

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func TestExecute_Failing_Run(t *testing.T) {
	var runCalls, shutdownCalls, terminateCalls atomicCounter

	r := &runnerMock{
		run: func() error {
			runCalls.Increment()
			return errors.New("run error")
		},
		shutdown: func(ctx context.Context) error {
			shutdownCalls.Increment()
			return nil
		},
		terminate: func() error {
			terminateCalls.Increment()
			return nil
		},
	}

	err := Execute(
		context.Background(),
		r,
	)

	if err == nil {
		t.Errorf("expected error but got nil")
	}

	if runCalls.Value() != 1 {
		t.Errorf("expected 1 run call but got %d", runCalls)
	}

	if shutdownCalls.Value() != 0 {
		t.Errorf("expected 0 shutdown call but got %d", shutdownCalls)
	}

	if terminateCalls.Value() != 0 {
		t.Errorf("expected 0 terminate call but got %d", terminateCalls)
	}
}

func TestExecute_Successfully_Shutdown(t *testing.T) {
	var runCalls, shutdownCalls, terminateCalls atomicCounter
	r := &runnerMock{
		run: func() error {
			runCalls.Increment()
			time.Sleep(2 * time.Second)
			return nil
		},
		shutdown: func(ctx context.Context) error {
			shutdownCalls.Increment()
			return nil
		},
		terminate: func() error {
			terminateCalls.Increment()
			return nil
		},
	}

	shutdownChannel := make(chan os.Signal, 1)
	go func() {
		time.AfterFunc(time.Second/2, func() {
			shutdownChannel <- os.Interrupt
		})
	}()

	err := Execute(
		context.Background(),
		r,
		WithShutdownListener(shutdownChannel),
		WithShutdownSignals(os.Interrupt),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if runCalls.Value() != 1 {
		t.Errorf("expected 1 run call but got %d", runCalls)
	}

	if shutdownCalls.Value() != 1 {
		t.Errorf("expected 1 shutdown call but got %d", shutdownCalls)
	}

	if terminateCalls.Value() != 0 {
		t.Errorf("expected 0 terminate call but got %d", terminateCalls)
	}
}

func TestExecute_Failing_Shutdown_Expecting_Terminate(t *testing.T) {
	var runCalls, shutdownCalls, terminateCalls atomicCounter
	r := &runnerMock{
		run: func() error {
			runCalls.Increment()
			time.Sleep(2 * time.Second)
			return nil
		},
		shutdown: func(ctx context.Context) error {
			shutdownCalls.Increment()
			select {
			case <-time.After(2 * time.Second):
				return errors.New("shutdown timeout")
			case <-ctx.Done():
				return ctx.Err()
			}
		},
		terminate: func() error {
			terminateCalls.Increment()
			return nil
		},
	}

	shutdownChannel := make(chan os.Signal, 1)
	go func() {
		time.AfterFunc(time.Second/2, func() {
			shutdownChannel <- os.Interrupt
		})
	}()

	err := Execute(
		context.Background(),
		r,
		WithShutdownListener(shutdownChannel),
		WithShutdownSignals(os.Interrupt),
		WithShutdownTimeout(time.Second/2),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if runCalls.Value() != 1 {
		t.Errorf("expected 1 run call but got %d", runCalls)
	}

	if shutdownCalls.Value() != 1 {
		t.Errorf("expected 1 shutdown call but got %d", shutdownCalls)
	}

	if terminateCalls.Value() != 1 {
		t.Errorf("expected 1 terminate call but got %d", terminateCalls)
	}
}

type atomicCounter struct {
	value int32
}

func (c *atomicCounter) Increment() {
	atomic.AddInt32(&c.value, 1)
}

func (c *atomicCounter) Value() int32 {
	return atomic.LoadInt32(&c.value)
}

type runnerMock struct {
	run       func() error
	shutdown  func(context.Context) error
	terminate func() error
}

func (r *runnerMock) Run() error {
	return r.run()
}

func (r *runnerMock) Shutdown(ctx context.Context) error {
	return r.shutdown(ctx)
}

func (r *runnerMock) Terminate() error {
	return r.terminate()
}
