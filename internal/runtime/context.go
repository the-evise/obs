package runtime

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/yourorg/obs/internal/log"
)

/*
   Responsibility:
       - Attach observability data to context
       - Extract data from context
   Owns:
       - Context keys
       - Field extraction logic
   >LoggerFromContext(ctx) log.Logger
   >TraceID(ctx) string
   >RequestID(ctx) string
   All request-scoped correlation flows through this file.
   No other file is allowed to read context directly.
*/

type ctxKey string

const (
	ctxLoggerKey   ctxKey = "obs.logger"
	ctxRequestID   ctxKey = "obs.request_id"
	ctxTraceID     ctxKey = "obs.trace_id"
	ctxSpanID      ctxKey = "obs.span_id"
	requestIDBytes        = 12 // 24 hex chars
)

var activeRuntime atomic.Pointer[Runtime]

func setRuntime(r *Runtime) {
	if r == nil {
		activeRuntime.Store(nil)
		return
	}
	activeRuntime.Store(r)
}

func getRuntime() *Runtime {
	return activeRuntime.Load()
}

// ContextWithLogger stores the request-scoped logger in ctx.
func ContextWithLogger(ctx context.Context, logger log.Logger) context.Context {
	if logger == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

// LoggerFromContext returns the logger previously attached via ContextWithLogger.
func LoggerFromContext(ctx context.Context) log.Logger {
	if ctx == nil {
		return nil
	}
	logger, _ := ctx.Value(ctxLoggerKey).(log.Logger)
	return logger
}

func contextWithRequestID(ctx context.Context, id string) context.Context {
	if id == "" {
		id = newRequestID()
	}
	return context.WithValue(ctx, ctxRequestID, id)
}

// RequestID fetches the correlation identifier for the current request.
func RequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	value, _ := ctx.Value(ctxRequestID).(string)
	return value
}

func ensureRequestID(ctx context.Context) string {
	if id := RequestID(ctx); id != "" {
		return id
	}
	return newRequestID()
}

func contextWithTrace(ctx context.Context, traceID, spanID string) context.Context {
	if traceID != "" {
		ctx = context.WithValue(ctx, ctxTraceID, traceID)
	}
	if spanID != "" {
		ctx = context.WithValue(ctx, ctxSpanID, spanID)
	}
	return ctx
}

// TraceID returns the current trace identifier if one has been recorded.
func TraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	value, _ := ctx.Value(ctxTraceID).(string)
	return value
}

// SpanID returns the span identifier connected to ctx.
func SpanID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	value, _ := ctx.Value(ctxSpanID).(string)
	return value
}

func newRequestID() string {
	buf := make([]byte, requestIDBytes)
	if _, err := rand.Read(buf); err != nil {
		// rand.Read should not fail; fall back to timestamp-based value
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(buf)
}
