package runtime

/*
   Responsibility:
       - Creates and wires all subsystems
       - Owns the lifecycle
       - Exposes the public Runtime
   Owns:
       - Logger
       - Tracer provider
       - Meter provider
       - Shutdown coordination
   Dependency Direction:
   config -> runtime -> (logger, tracer, metrics)
   This is the only place where:
   - OTel SDKs are initialized
   - Exporters are created
   - Globals (if any) are set
*/

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/yourorg/obs/internal/config"
	"github.com/yourorg/obs/internal/health"
	"github.com/yourorg/obs/internal/log"
	"github.com/yourorg/obs/internal/metrics"
	"go.opentelemetry.io/otel/trace"
)

type Runtime struct {
	// Process-level context
	ctx      context.Context
	cancel   context.CancelFunc
	services []Service

	// Infrastructure
	Config         *config.Config
	Logger         log.Logger
	Metrics        metrics.Metrics
	Tracer         trace.Tracer
	shutdownTracer func(context.Context) error

	// Lifecycle control
	wg sync.WaitGroup

	shutdownOnce sync.Once
	shutdownTO   time.Duration
}

type FailureType int

const (
	FailureTransient FailureType = iota
	FailureFatal
)

type Failure struct {
	Type FailureType
	Err  error
}

func (f *Failure) Error() string {
	return f.Err.Error()
}

func (f *Failure) Unwrap() error {
	return f.Err
}

func Fatal(err error) error {
	if err == nil {
		return nil
	}
	return &Failure{Type: FailureFatal, Err: err}
}

func Transient(err error) error {
	if err == nil {
		return nil
	}
	return &Failure{Type: FailureTransient, Err: err}
}

func NewRuntime(
	ctx context.Context,
	logger log.Logger,
	cfg *config.Config,
	shutdownTimeOut time.Duration,
) *Runtime {
	ctx, cancel := context.WithCancel(ctx)

	runtime := &Runtime{
		ctx:        ctx,
		cancel:     cancel,
		services:   []Service{},
		Config:     cfg,
		Logger:     logger,
		shutdownTO: shutdownTimeOut,
	}

	if runtime.Logger == nil {
		runtime.Logger = log.NewStdLogger()
	}

	if runtime.Metrics == nil {
		runtime.Metrics = metrics.NewNoop()
	}

	runtime.initTracer(ctx)

	setRuntime(runtime)

	return runtime
}

func (r *Runtime) Ready() bool {
	for _, svc := range r.services {
		report := svc.Health()
		if report.Status == health.StatusUnhealthy {
			return false
		}
	}
	return true
}

func (r *Runtime) Context() context.Context {
	return r.ctx
}

func (r *Runtime) RegisterServices(svc ...Service) {
	for _, s := range svc {
		if s == nil {
			continue
		}
		r.services = append(r.services, s)
	}
}

func (r *Runtime) Go(fn func(ctx context.Context) error) {
	r.wg.Add(1)

	go func() {
		defer r.wg.Done()

		err := fn(r.ctx)
		if err != nil {
			r.Logger.Error(
				"goroutine failed",
				log.F("error", err),
			)
		}
		if err == nil || errors.Is(err, context.Canceled) {
			return
		}

		r.handleFailure(err)
	}()
}

func (r *Runtime) handleFailure(err error) {
	r.Logger.Error("service failed", log.F("error", err))

	if isFatal(err) {
		r.Shutdown()
	}
}

func isFatal(err error) bool {
	var f *Failure
	if errors.As(err, &f) {
		return f.Type == FailureFatal
	}
	return true
}
