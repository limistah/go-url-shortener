# step-02-lifecycle

Introduces lifecycle-managed server startup and shutdown.

- What changed: `server.NewServer` appends `OnStart` and `OnStop` hooks via `fx.Lifecycle`.
- Why it matters: handles run/stop semantics cleanly without hand-written signal plumbing.
- Verify: start with `go run ./cmd/server`, stop with Ctrl-C and observe graceful shutdown.
