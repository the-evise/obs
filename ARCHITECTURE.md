# Architecture

This document describes the internal architecture of `obs`.  
It is intended for maintainers and contributors, not end users.

The goal of this architecture is to provide a **stable, minimal, and production-safe observability layer** that scales across many services without hidden behavior.

---

## Design Goals

1. **Explicit lifecycle**
    - No implicit globals without ownership
    - All resources are created, owned, and shut down deliberately

2. **Strong observability contract**
    - Logs, traces, and metrics are always correlated
    - Required fields are guaranteed, not optional

3. **Minimal surface area**
    - Small public API
    - Clear internal boundaries

4. **One-way dependencies**
    - No circular imports
    - Core never depends on middleware

---

## High-Level Overview

`obs` is a **composition layer**, not a framework.

It orchestrates:
- OpenTelemetry (tracing)
- Prometheus (metrics)
- Zap (logging)

It does **not** reimplement these systems.

At runtime, everything is coordinated through:
- a single `Runtime`
- a shared `context.Context` contract

---

## Directory Structure

```text
pkg/obs/
├── config.go
├── runtime.go
├── logger.go
├── tracer.go
├── metrics.go
├── context.go
├── middleware_http.go
├── middleware_grpc.go
├── async.go
└── shutdown.go
```

Each file owns exactly one concern.

---

## Core Modules

### `config.go` — Configuration Contract

**Responsibility**
- Defines the public `Config` struct
- Applies defaults
- Validates user input

**Rules**
- No side effects
- No initialization logic
- No global state

**Used by**
- `runtime.go` only

---

### `runtime.go` — System Orchestrator

**Responsibility**
- Initializes all observability subsystems
- Wires dependencies together
- Owns the full lifecycle

**Public Output**
```go
type Runtime struct {
	Logger   Logger
	Tracer   trace.Tracer
	Meter    metric.Meter
	Shutdown func(ctx context.Context) error
}
```

**Rules**

- The only place allowed to initialize SDKs 
- The only place allowed to set globals 
- Setup must be idempotent

---

### `logger.go` — Logging Adapter

**Responsibility**

- Create the base structured logger
- Attach static fields (service, environment)
- Adapt Zap to the `obs.Logger` interface

**Does NOT**

- Read from context 
- Inject request-scoped fields 

Context enrichment is handled centrally.

---

### `tracer.go` — Tracing Setup

**Responsibility**

- Configure tracer provider 
- Configure sampling 
- Register global tracer provider

**Rules**

- Sampling must be parent-based 
- Defaults must be safe for scale

**Does NOT**

- Start spans 
- Perform propagation

---

### `metrics.go` — Metrics Setup

**Responsibility**

- Create Prometheus registry 
- Expose /metrics endpoint 
- Define base instruments

**Rules**

- No unbounded label sets 
- Route normalization required

**Does NOT**

- Record request metrics directly

---

## Context & Correlation
### `context.go` — Correlation Spine (Critical)

This file is the single source of truth for request-scoped data.

**Responsibility**

- Attach observability data to context 
- Extract data from context

**Owns**

- Context keys 
- Mapping between context → log fields

**Public helpers**

```go
FromContext(ctx) Logger
TraceID(ctx) string
SpanID(ctx) string
RequestID(ctx) string
```


**Rule**

> No other file may directly read from `context.Context`.

---

## Instrumentation Layers
### `middleware_http.go` — HTTP Instrumentation

**Responsibility**

- Create a root span per request
- Generate or propagate request ID 
- Bind logger and span to context 
- Record latency and status 
- Handle panics

**Hard Guarantees**

- Exactly one span per request 
- Exactly two lifecycle logs (start/end)
- Errors are always logged and marked on the span

---

### `middleware_grpc.go` — gRPC Instrumentation

Same guarantees as HTTP, adapted to gRPC semantics.

**No divergence in behavior** between protocols.

---

## Async & Lifecycle
### `async.go` — Safe Goroutines

**Responsibility**

- Preserve context across goroutines 
- Ensure spans and logs remain correlated 
- Recover and log panics

**Public API**

```go
obs.Go(ctx, func(ctx context.Context))
```

### `shutdown.go` — Graceful Termination

**Responsibility**

- Flush exporters 
- Handle OS signals (optional)
- Ensure idempotent shutdown

**Rules**

- Shutdown must be safe to call multiple times 
- Must respect context deadlines 

---

## Request Lifecycle (HTTP)

1. Request enters HTTP middleware 
2. Request ID generated or extracted 
3. Root span started 
4. Logger enriched with:
   - service.name 
   - request_id 
   - trace_id 
   - span_id 
5. Context populated 
6. Handler executes 
7. Metrics recorded 
8. Span ended 
9. End-of-request log emitted

This flow is deterministic and testable.

---

## Why This Architecture

This design ensures:

* predictable behavior at scale
* easy reasoning under incident pressure
* strong defaults with minimal configuration
* long-term API stability

The architecture prioritizes **operational correctness over flexibility.**

---

## Evolution Rules

* New features must not weaken guarantees
* Defaults must remain safe under load
* Any change affecting signals must be documented
* Breaking changes require a major version bump

---

## Summary

`obs` is intentionally boring.

That is the point.

**It provides:**

* one way to do observability
* one lifecycle
* one contract

This is what allows it to scale across teams and services without entropy.