package usecase

import (
	"context"
	"time"
)

func RetryWithBackoff(ctx context.Context, operation func() error) error {
	const maxAttempts = 5
	const initialDelay = 1 * time.Second
	const maxDelay = 30 * time.Second

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		if ctx.Err() != nil {
			return ctx.Err()
		}

		delay := initialDelay * time.Duration(1<<uint(attempt-1))
		if delay > maxDelay {
			delay = maxDelay
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	return lastErr

}
