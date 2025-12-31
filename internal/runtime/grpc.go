package runtime

import (
	"context"
	"time"

	"github.com/yourorg/obs/internal/log"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

/*
   Same responsibilities as HTTP, adapted to gRPC semantics.
   Shares:
       - context helpers
       - metric instruments
       - logging contract
*/

func GRPCUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		rt := getRuntime()
		if rt == nil {
			return handler(ctx, req)
		}

		start := time.Now()

		reqID := ensureRequestID(ctx)
		ctx = contextWithRequestID(ctx, reqID)

		tracer := rt.Tracer
		ctx, span := tracer.Start(
			ctx,
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		logger := rt.Logger.With(
			log.String("request_id", reqID),
			log.String("trace_id", span.SpanContext().TraceID().String()),
			log.String("span_id", span.SpanContext().SpanID().String()),
			log.String("rpc.method", info.FullMethod),
		)

		logger.Debug("rpc started")

		resp, err := handler(ctx, req)

		latency := time.Since(start)

		if err != nil {
			logger.Error("rpc error", log.Error(err))
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}

		logger.Debug(
			"rpc completed",
			log.Int64("latency_ms", latency.Milliseconds()),
		)

		return resp, err
	}
}
