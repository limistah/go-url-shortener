# step-09-testing

Introduces Fx-native test composition.

- What changed: tests use `fxtest.New`, `fx.Populate`, and `fx.Replace` for handler-level verification.
- Why it matters: real dependency graph tests with lightweight unit-test ergonomics.
- Verify: run `go test ./...` and confirm handler tests pass.
