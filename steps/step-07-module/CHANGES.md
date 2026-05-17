# step-07-module

Introduces feature-scoped module composition.

- What changed: providers are grouped into `config.Module`, `storage.Module`, `api.Module`, `server.Module`.
- Why it matters: `main` becomes composition-only and teams can own module boundaries.
- Verify: app starts from `cmd/server` with the same endpoint behavior.
