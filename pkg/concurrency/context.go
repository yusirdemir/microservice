package concurrency

import (
	"context"
)

func RunWithContext[T any](ctx context.Context, fn func() (T, error)) (T, error) {
	var result T
	var err error

	done := make(chan struct{})

	go func() {
		defer close(done)
		result, err = fn()
	}()

	select {
	case <-done:
		return result, err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}
