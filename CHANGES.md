# CHANGES

This branch reflects the concrete end-state (`step-09-testing`) of the Uber Fx URL shortener walkthrough.

- Reworked the runnable app to `cmd/server` and concrete package boundaries.
- Added interface-first store wiring with `fx.Annotate` + `fx.As`.
- Added route value groups (`group:"routes"`) and named stores (`name:"hot"`, `name:"cold"`).
- Added module composition and scoped logger decoration in `api` module.
- Added focused `fxtest` examples using `fx.Populate` and `fx.Replace`.
- Added `steps/step-00-manual` through `steps/step-09-testing` folders, each with a `CHANGES.md` narrative.
