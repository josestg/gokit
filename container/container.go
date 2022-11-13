package container

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Runner manages the lifecycle of the application.
type Runner interface {
	// Run starts the application.
	Run() error

	// Shutdown shuts down the application gracefully.
	Shutdown(ctx context.Context) error

	// Terminate is used to force close the application.
	Terminate() error
}

// ShutdownOption is an option to configure the shutdown behavior.
type ShutdownOption struct {
	Signals  []os.Signal
	Timeout  time.Duration
	Listener chan os.Signal
}

// Option is an option to configure the ShutdownOption.
type Option func(*ShutdownOption)

// WithShutdownTimeout configures the shutdown timeout.
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(o *ShutdownOption) {
		o.Timeout = timeout
	}
}

// WithShutdownSignals configures the shutdown signals.
func WithShutdownSignals(signals ...os.Signal) Option {
	return func(o *ShutdownOption) {
		o.Signals = signals
	}
}

// WithShutdownListener configures the shutdown listener.
func WithShutdownListener(listener chan os.Signal) Option {
	return func(o *ShutdownOption) {
		o.Listener = listener
	}
}

// DefaultShutdownOption is the default shutdown option.
func DefaultShutdownOption() ShutdownOption {
	return ShutdownOption{
		Signals:  []os.Signal{os.Interrupt, syscall.SIGTERM},
		Timeout:  5 * time.Second,
		Listener: make(chan os.Signal, 1),
	}
}

// Execute executes the Runner. It blocks until the Runner is gracefully shutdown or terminated.
// The timeout is to give the Runner a chance to shut down gracefully when the listed signals are received.
// The shutdownChannel must be a buffered channel with size 1
func Execute(ctx context.Context, runner Runner, opts ...Option) error {
	// Apply the options.
	shutdownOption := DefaultShutdownOption()
	for _, opt := range opts {
		opt(&shutdownOption)
	}

	// Notify the shutdown channel when the signals are received.
	signal.Notify(shutdownOption.Listener, shutdownOption.Signals...)

	errChannel := make(chan error, 1)
	go func() {
		err := runner.Run()
		if err != nil {
			errChannel <- err
		}
	}()

	// wait unit the runnerMock is done or the signals are received.
	select {
	case <-shutdownOption.Listener:
		ctx, cancel := context.WithTimeout(ctx, shutdownOption.Timeout)
		defer cancel()
		if err := runner.Shutdown(ctx); err != nil {
			if err := runner.Terminate(); err != nil {
				return fmt.Errorf("failed to terminate the runner: %w", err)
			}
		}
	case err := <-errChannel:
		if err != nil {
			return fmt.Errorf("running the runner: %w", err)
		}
	}
	return nil
}
