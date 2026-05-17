# step-06-named-values

Introduces named dependencies for multiple stores.

- What changed: hot and cold `storage.Store` providers are tagged with `name:"hot"` and `name:"cold"`.
- Why it matters: same interface type can be injected unambiguously into one consumer.
- Verify: `go run ./cmd/server` and confirm shorten/lookup flow uses dual-store wiring.
