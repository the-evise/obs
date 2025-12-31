# Repository Guidelines

## Project Structure & Module Organization
Runtime logic lives in `internal/`: `internal/runtime` wires logging, tracing, metrics, and shutdown; `internal/log`, `internal/metrics`, and `internal/tracer` provide the concrete signal adapters; `internal/config` owns validation/defaults; `internal/service` hosts sample services plus the readiness checker used by `cmd/checker`. Executables belong under `cmd/` (for example, `go run ./cmd/checker` exercises the observability stack end-to-end). Keep shared or exported APIs under `pkg/` so consumers never depend on `internal/`. Place assets or fixtures beside the code that consumes them to avoid stale cross-package references.

## Build, Test, and Development Commands
`go mod tidy` keeps dependencies aligned; run it after adding imports. `go build ./...` must stay clean before any push. Execute `go run ./cmd/checker` locally to verify bootstrap flows and `/healthz` endpoints. Always run `go test ./...` (add `-race` when touching concurrency) and `go vet ./...` before opening a pull request. The CI workflow mirrors this and also runs `golangci-lint run`, so running the same command locally avoids red builds.

## Coding Style & Naming Conventions
Use Go 1.22, enforced gofmt/goimports formatting, and keep files ASCII-only. Favor small, domain-scoped packages (`internal/runtime/http.go`, `internal/service/example`) and keep exported identifiers prefixed with the subsystem (`Runtime`, `CheckerService`). Interfaces should be named after the behavior (`Logger`, `HealthReporter`) and defined where they are consumed. Log/metric names should stay lowercase with dots (`runtime.start.duration`) and match the naming seen in `internal/metrics`. Document any exported type or function with a full-sentence comment.

## Testing Guidelines
Adopt table-driven tests and store them next to the code (e.g., `internal/runtime/runtime_test.go`). Name tests `Test<Thing>` and helpers `<thing>Helper`. Where spans or metrics are involved, rely on the noop implementations in `internal/metrics/noop.go` or add fakes under the same package. Target at least the coverage that keeps `go test ./... -cover` from regressing (document deltas in PRs). For services, wire `internal/service/example` against the checker binary in integration tests to keep HTTP/gRPC middleware behavior deterministic.

## Commit & Pull Request Guidelines
History currently uses concise, imperative subjects (`first commit`). Keep following that pattern: one short (<72 char) subject, optional wrapped body describing motivation, and reference the issue ID (`Fixes #123`) when applicable. Pull requests must: describe the change (“what” and “why”), list validation commands (`go test ./...`, `golangci-lint run`), link any incidents or design docs, and attach screenshots/log excerpts when altering emitted signals. Update `ARCHITECTURE.md` or `README.md` if the public contract or guarantees move.

## Observability & Configuration Tips
`internal/config` already sets safe defaults (OTLP exporter, Prometheus scrape port). Extend configuration by adding fields to the `Config` struct, validating inside the same package, and threading values through `internal/runtime/runtime.go` so one call to `runtime.Run(ctx, cfg)` still owns every subsystem. Never mutate global loggers or tracers outside the runtime module; add helper methods or context accessors instead. Keep secrets out of the repo—inject endpoints/keys via environment variables and document them in this file or `README.md` when new knobs appear.
