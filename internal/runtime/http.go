package runtime

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yourorg/obs/internal/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type HTTPServer struct {
	log log.Logger
	srv *http.Server
}

func NewHTTPServer(logger log.Logger, srv *http.Server) *HTTPServer {
	return &HTTPServer{log: logger, srv: srv}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	if s.srv == nil {
		return fmt.Errorf("http server not configured")
	}

	errCh := make(chan error, 1)

	go func() {
		s.log.Info("starting http server")
		errCh <- s.srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		s.log.Info("http server shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return s.srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func HTTPHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rt := getRuntime()
		if rt == nil {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		ctx := r.Context()
		reqID := ensureRequestID(ctx)
		ctx = contextWithRequestID(ctx, reqID)

		tracer := rt.Tracer
		ctx, span := tracer.Start(
			ctx,
			spanNameFromRequest(r),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
			attribute.String("http.scheme", r.URL.Scheme),
		)

		logger := rt.Logger.With(
			log.String("request_id", reqID),
			log.String("trace_id", span.SpanContext().TraceID().String()),
			log.String("span_id", span.SpanContext().SpanID().String()),
			log.String("http.method", r.Method),
			log.String("http.path", r.URL.Path),
		)

		logger.Debug("request started")

		r = r.WithContext(ctx)
		rw := wrapResponseWriter(w)

		defer func() {
			if rec := recover(); rec != nil {
				err := panicToError(rec)

				logger.Error("panic recovered", log.Error(err))
				span.RecordError(err)
				span.SetStatus(codes.Error, "panic")

				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			latency := time.Since(start)

			span.SetAttributes(
				attribute.Int("http.status", rw.status),
			)

			if rw.status >= 500 {
				span.SetStatus(codes.Error, http.StatusText(rw.status))
			} else {
				span.SetStatus(codes.Ok, "")
			}

			logger.Debug(
				"request completed",
				log.Int("http.status", rw.status),
				log.Int64("latency_ms", latency.Milliseconds()),
			)
		}()

		next.ServeHTTP(rw, r)
	})
}

func spanNameFromRequest(r *http.Request) string {
	if r == nil || r.URL == nil {
		return "http.request"
	}
	return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
}

func panicToError(rec any) error {
	switch v := rec.(type) {
	case error:
		return v
	case string:
		return fmt.Errorf("panic: %s", v)
	default:
		return fmt.Errorf("panic: %v", v)
	}
}
