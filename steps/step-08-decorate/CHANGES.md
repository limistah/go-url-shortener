# step-08-decorate

Introduces scoped decoration for cross-cutting behavior.

- What changed: `api.Module` decorates `*zap.Logger` to include `module=api`.
- Why it matters: module-local behavior is injected without touching each constructor.
- Verify: API logs include the module field while non-API logs remain undecorated.
