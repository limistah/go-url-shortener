# step-00-manual

Baseline only: plain Go manual wiring in `main.go` style.

- What changed: no Fx container; constructors are called directly in order.
- Why it matters: establishes the pain point (ordering, signature ripple, centralized wiring).
- Verify: `go run ./cmd/server` and exercise `POST /shorten`, `GET /{slug}`.
