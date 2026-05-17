# step-04-annotate-as

Introduces interface-based wiring via annotations.

- What changed: store providers are wrapped with `fx.Annotate(..., fx.As(new(storage.Store)))`.
- Why it matters: handlers depend on `storage.Store` interface, not concrete implementations.
- Verify: `go run ./cmd/server` still serves both endpoints with in-memory store.
