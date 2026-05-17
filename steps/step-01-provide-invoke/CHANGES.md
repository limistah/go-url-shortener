# step-01-provide-invoke

Introduces container registration and graph construction.

- What changed: constructor functions are registered with `fx.Provide`; side-effect entrypoint uses `fx.Invoke`.
- Why it matters: removes fragile manual constructor ordering.
- Verify: dependency graph builds with `go run ./cmd/server`.
