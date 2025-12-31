package runtime

import (
	"context"
	"time"

	"github.com/yourorg/obs/internal/log"
)

/*
   Responsibility:
       - Graceful exporter shutdown
       - Signal handling
       - Idempotency
   Owns:
       - WaitForShutdown
       - Runtime.Shutdown
   Does NOT:
       - Create resources
*/

func (r *Runtime) Shutdown() {
	r.shutdownOnce.Do(func() {
		r.Logger.Info("runtime shutting down")
		r.cancel()
		r.wg.Wait()
		r.closeTracer(r.shutdownContext())
	})
}

func (r *Runtime) shutdownContext() context.Context {
	ctx, _ := context.WithTimeout(
		context.Background(),
		r.shutdownTO,
	)
	return ctx
}

func (r *Runtime) ShutdownWithTimeout(d time.Duration) {
	r.shutdownOnce.Do(func() {
		r.Logger.Info("runtime shutting down", log.F("timeout", d))
		r.cancel()
		done := make(chan struct{})
		go func() {
			r.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			r.Logger.Info("runtime complete")
		case <-time.After(d):
			r.Logger.Error("shutdown timedout exceeded")
		}

		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		r.closeTracer(ctx)
	})
}

func (r *Runtime) closeTracer(ctx context.Context) {
	if r.shutdownTracer == nil {
		return
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if err := r.shutdownTracer(ctx); err != nil {
		r.Logger.Error("tracer shutdown failed", log.F("error", err))
	}
}
