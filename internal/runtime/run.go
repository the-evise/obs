package runtime

import (
	"context"
	"errors"
	"sync"

	"github.com/yourorg/obs/internal/log"
)

// Run starts all registered services and blocks until the runtime shuts down.
func (r *Runtime) Run() error {
	r.listenForSignals()

	var wg sync.WaitGroup
	errCh := make(chan error, len(r.services))

	for _, svc := range r.services {
		svc := svc
		if svc == nil {
			continue
		}

		wg.Add(1)
		go func(s Service) {
			defer wg.Done()

			if err := s.Run(r.ctx); err != nil && !errors.Is(err, context.Canceled) {
				errCh <- err
			}
		}(svc)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if isFatal(err) {
			r.Logger.Error("fatal service failure", log.Error(err))
			r.Shutdown()
			break
		}

		r.Logger.Error("service error", log.Error(err))
	}

	<-r.ctx.Done()

	return r.ctx.Err()
}
