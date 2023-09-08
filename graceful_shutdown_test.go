package gfshutdown_test

import (
	"context"
	"errors"
	"testing"
	"time"

	gfshutdown "github.com/gelmium/graceful-shutdown"
)

func TestGracefulShutdown(t *testing.T) {
	// Define a mock operation that waits for a specified duration before returning
	mockOp := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return nil
		}
	}

	// Define a map of operations to be passed to GracefulShutdown
	ops := map[string]gfshutdown.Operation{
		"op1": mockOp,
		"op2": mockOp,
	}

	// Define a timeout for the GracefulShutdown function
	timeout := 1 * time.Second

	// Call GracefulShutdown with a trigger context that is immediately cancelled
	triggerCtx, cancel := context.WithCancel(context.Background())
	cancel()
	wait := gfshutdown.GracefulShutdown(triggerCtx, timeout, ops)

	// Ensure that the function waits for all operations to complete before returning
	select {
	case r := <-wait:
		// Waited for all operations to complete
		if r != 0 {
			t.Error("GracefulShutdown did not exit gracefully")
		}
	case <-time.After(2 * time.Second):
		t.Error("GracefulShutdown did not wait for all operations to complete")
	}
}

func TestGracefulShutdownWithError(t *testing.T) {
	// Define a mock operation that returns an error
	mockOp := func(ctx context.Context) error {
		return errors.New("operation failed")
	}

	// Define a map of operations to be passed to GracefulShutdown
	ops := map[string]gfshutdown.Operation{
		"op1": mockOp,
	}

	// Define a timeout for the GracefulShutdown function
	timeout := 1 * time.Second

	// Call GracefulShutdown with a trigger context that is immediately cancelled
	triggerCtx, cancel := context.WithCancel(context.Background())
	cancel()
	wait := gfshutdown.GracefulShutdown(triggerCtx, timeout, ops)

	// Ensure that the function waits for all operations to complete before returning
	select {
	case r := <-wait:
		// Waited for all operations to complete
		if r != 0 {
			t.Error("GracefulShutdown did not exit gracefully")
		}
	case <-time.After(2 * time.Second):
		t.Error("GracefulShutdown did not wait for all operations to complete")
	}
}

// Add a test case with timeout
func TestGracefulShutdownWithTimeout(t *testing.T) {
	// Define a mock operation that waits for a specified duration before returning
	mockOp := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
			return nil
		}
	}

	// Define a map of operations to be passed to GracefulShutdown
	ops := map[string]gfshutdown.Operation{
		"op1": mockOp,
	}

	// Define a timeout for the GracefulShutdown function
	timeout := 1 * time.Second

	// Call GracefulShutdown with a trigger context that is immediately cancelled
	triggerCtx, cancel := context.WithCancel(context.Background())
	cancel()
	wait := gfshutdown.GracefulShutdown(triggerCtx, timeout, ops)

	// Ensure that the function waits for all operations to complete before returning
	select {
	case r := <-wait:
		// Gracefully shutdown exit before 2 seconds
		if r != 0 {
			t.Error("GracefulShutdown did not exit gracefully")
		}
		break
	case <-time.After(2 * time.Second):
		// Waited for all operations to complete
		t.Error("GracefulShutdown did not timeout but instead wait for all operations to complete")
	}
}

func TestGracefulShutdownWithTimeoutForceExit(t *testing.T) {
	// Define a mock operation that waits for a specified duration before returning
	mockOp := func(ctx context.Context) error {
		select {
		case <-time.After(2 * time.Second):
			return nil
		}
	}

	// Define a map of operations to be passed to GracefulShutdown
	ops := map[string]gfshutdown.Operation{
		"op1": mockOp,
	}

	// Define a timeout for the GracefulShutdown function
	timeout := 1 * time.Second

	// Call GracefulShutdown with a trigger context that is immediately cancelled
	triggerCtx, cancel := context.WithCancel(context.Background())
	cancel()
	wait := gfshutdown.GracefulShutdown(triggerCtx, timeout, ops)

	// Ensure that the function waits for all operations to complete before returning
	select {
	case r := <-wait:
		// Process did not force exit
		if r != 1 {
			t.Error("GracefulShutdown did not time out")
		}
	case <-time.After(2 * time.Second):
		// Waited for all operations to complete
		t.Error("GracefulShutdown did not timeout but instead wait for all operations to complete")
	}
}
