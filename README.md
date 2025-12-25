# obs — Opinionated Observability Bootstrap for Go

`obs` is a small, opinionated library that provides a **production-ready observability baseline** for Go services.

It standardizes how logging, tracing, metrics, and request correlation are wired so teams stop re-implementing the same setup in every repository.

One setup call. Consistent signals. No framework lock-in.

---

## Why This Exists

Observability in Go is usually:
- repetitive
- inconsistent across services
- added late (often during incidents)

`obs` exists to make **correct observability the default**, from the first line of code.

---

## What You Get (Guaranteed)

Every service using `obs` gets:

- **Structured JSON logs**
- **OpenTelemetry tracing** (OTLP exporter)
- **Prometheus metrics** exposed at `/metrics`
- **Request / trace / span correlation**
- **HTTP and gRPC auto-instrumentation**

These guarantees apply uniformly across services.

---

## Non-Goals

`obs` intentionally does **not**:

- replace OpenTelemetry
- hide or abstract vendors
- provide dashboards or alerts
- act as a web or application framework

It is infrastructure glue. Nothing more.