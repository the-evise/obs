package internal

import (
	"context"
	"fmt"

	"github.com/yourorg/obs/internal/log"
	"github.com/yourorg/obs/internal/runtime"
)

// Go runs fn in a separate goroutine while keeping request context
// and recovering from panics so signals stay correlated.
func Go(ctx context.Context, fn func(ctx context.Context)) {
	if fn == nil {
		return
	}

	if ctx == nil {
		ctx = context.Background()
	}

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				logger := runtime.LoggerFromContext(ctx)
				if logger == nil {
					logger = log.NewStdLogger()
				}
				logger.Error("goroutine panic", log.Error(fmt.Errorf("panic: %v", rec)))
			}
		}()

		fn(ctx)
	}()
}
