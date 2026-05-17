# step-03-in-out

Introduces structured dependency params/results.

- What changed: handler constructors use `fx.In` parameter structs.
- Why it matters: constructor signatures stay readable as dependencies grow.
- Verify: `go build ./...` and run server to confirm unchanged endpoint behavior.
