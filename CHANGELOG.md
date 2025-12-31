# Changelog

All notable changes to this project will be documented in this file.

This project follows **Semantic Versioning (SemVer)**:
- **MAJOR**: incompatible API changes
- **MINOR**: backward-compatible functionality
- **PATCH**: backward-compatible bug fixes

The format is inspired by *Keep a Changelog*, but tailored for infrastructure libraries.

---

## [Unreleased]

### Added
- Explicit runtime lifecycle (`Setup`, `Shutdown`)
- Defined observability contract (guaranteed fields)
- HTTP and gRPC middleware with strict behavioral guarantees

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

---

## [0.1.0] — Initial Public Release

### Added
- Global observability bootstrap via `MustSetup`
- Structured JSON logging (Zap-based)
- OpenTelemetry tracing with OTLP exporter
- Prometheus metrics with `/metrics` endpoint
- HTTP middleware:
    - request ID injection
    - span creation
    - latency and status recording
- gRPC unary server interceptor
- Context-based logger access via `FromContext(ctx)`
- Example service under `cmd/example`

### Notes
- This release establishes the **baseline API**.
- No API stability guarantees beyond `0.x`.
- Defaults favor simplicity over fine-grained control.

---

## Versioning Policy

- **0.x.y**  
  Rapid iteration. APIs may change without notice.

- **1.0.0**  
  API is considered stable.  
  Breaking changes require a major version bump.

- **1.x.y**
    - `x` increases for new features
    - `y` increases for bug fixes and internal improvements

---

## Upgrade Guidance

Breaking changes will always include:
- clear migration notes
- rationale for the change
- before/after examples where applicable

If a change is operationally risky (metrics, logging, tracing behavior), it will be highlighted explicitly.

---

## Commit Discipline (Recommended)

To keep this changelog accurate:
- use conventional commits where possible
- link PRs or issues in entries
- update `[Unreleased]` continuously

---

This changelog is part of the library’s **trust contract**.
