package example

import (
	"context"
	"time"
)

func (s *ExampleService) doWork(ctx context.Context) error {
	select {
	case <-time.After(2 * time.Second):
		return nil // simulate success

	case <-ctx.Done():
		return ctx.Err()
	}
}
