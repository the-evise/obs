package runtime

import (
	"context"

	"github.com/yourorg/obs/internal/config"
	"github.com/yourorg/obs/internal/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "github.com/yourorg/obs/runtime"

func setupTracerProvider(ctx context.Context, cfg *config.Config) (trace.Tracer, func(context.Context) error, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		tracer, shutdown := noopTracer()
		return tracer, shutdown, err
	}

	res, resErr := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName(cfg)),
			attribute.String("service.environment", serviceEnv(cfg)),
		),
	)
	if resErr != nil {
		res = resource.Default()
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
	)

	otel.SetTracerProvider(tp)

	return tp.Tracer(tracerName), tp.Shutdown, nil
}

func noopTracer() (trace.Tracer, func(context.Context) error) {
	provider := trace.NewNoopTracerProvider()
	return provider.Tracer(tracerName), func(context.Context) error { return nil }
}

func serviceName(cfg *config.Config) string {
	if cfg == nil || cfg.App.Name == "" {
		return "obs"
	}
	return cfg.App.Name
}

func serviceEnv(cfg *config.Config) string {
	if cfg == nil || cfg.App.Env == "" {
		return "development"
	}
	return cfg.App.Env
}

func (r *Runtime) initTracer(ctx context.Context) {
	tracer, shutdown, err := setupTracerProvider(ctx, r.Config)
	if err != nil {
		r.Logger.Error("tracer initialization failed", log.F("error", err))
	}
	r.Tracer = tracer
	r.shutdownTracer = shutdown
}
