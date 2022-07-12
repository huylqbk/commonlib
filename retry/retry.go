package retry

import (
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// implement retry by fibonacci algorithm
type state [2]time.Duration

type BackoffFunc func() (time.Duration, bool)
type RetryFunc func(ctx context.Context) error

type backoff struct {
	state   unsafe.Pointer
	next    BackoffFunc
	l       sync.Mutex
	attempt uint64
}

// RetryHandle implement retry mechanism by fibonacci algorithm
// ctx for pass context
// base is init backoff value
// maxRetry > 0 is limit of retry
// f is retry function handler
func RetryHandle(ctx context.Context, base time.Duration, maxRetry uint64, f RetryFunc) error {
	if base < 0 {
		return errors.New("base is invalid")
	}

	b := &backoff{
		state: unsafe.Pointer(&state{0, base}),
	}
	b.defaultNext()
	if maxRetry > 0 {
		b.maxRetry(maxRetry)
	}

	return do(ctx, b, f)
}

func (b *backoff) defaultNext() {
	b.next = BackoffFunc(func() (time.Duration, bool) {
		for {
			curr := atomic.LoadPointer(&b.state)
			currState := (*state)(curr)
			next := currState[0] + currState[1]

			if next <= 0 {
				return math.MaxInt64, false
			}

			if atomic.CompareAndSwapPointer(&b.state, curr, unsafe.Pointer(&state{currState[1], next})) {
				return next, false
			}
		}
	})
}

func (b *backoff) maxRetry(max uint64) {
	b.next = BackoffFunc(func() (time.Duration, bool) {
		b.l.Lock()
		defer b.l.Unlock()

		if b.attempt >= max {
			return 0, true
		}
		b.attempt++

		val, stop := b.next()
		if stop {
			return 0, true
		}

		return val, false
	})
}

type retryableError struct {
	err error
}

// RetryableError marks an error as retryable.
func RetryableError(err error) error {
	if err == nil {
		return nil
	}
	return &retryableError{err}
}

// Unwrap implements error wrapping.
func (e *retryableError) Unwrap() error {
	return e.err
}

// Error returns the error string.
func (e *retryableError) Error() string {
	if e.err == nil {
		return "retryable: <nil>"
	}
	return "retryable: " + e.err.Error()
}

// context passed to the RetryFunc.
func do(ctx context.Context, b *backoff, f RetryFunc) error {
	for {
		// Return immediately if ctx is canceled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := f(ctx)
		if err == nil {
			return nil
		}

		// Not retryable
		var rerr *retryableError
		if !errors.As(err, &rerr) {
			return err
		}

		next, stop := b.next()
		if stop {
			return rerr.Unwrap()
		}

		// ctx.Done() has priority, so we test it alone first
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		t := time.NewTimer(next)
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
			continue
		}
	}
}
