# URL Shortener — Uber Fx Build Spec (Concrete Demo)

This repository now contains a concrete minimal URL shortener aligned to the talk narrative.

## App endpoints

- `POST /shorten`
- `GET /{slug}`

## Entry point

```bash
go run ./cmd/server
```

## One-command test

```bash
go test ./...
```

## Quick verify

```bash
curl -s -X POST localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com"}'
```

Use the returned slug:

```bash
curl -i -L localhost:8080/<slug>
```

## Structure used in the talk

- `cmd/server/main.go`
- `internal/config/config.go`
- `internal/storage/store.go`
- `internal/handler/shorten.go`
- `internal/handler/redirect.go`
- `internal/handler/route.go`
- `internal/api/module.go`
- `internal/server/server.go`

## Progressive steps

Because this environment works on one active branch, branch checkpoints are mirrored as folders:

- `steps/step-00-manual`
- `steps/step-01-provide-invoke`
- `steps/step-02-lifecycle`
- `steps/step-03-in-out`
- `steps/step-04-annotate-as`
- `steps/step-05-value-groups`
- `steps/step-06-named-values`
- `steps/step-07-module`
- `steps/step-08-decorate`
- `steps/step-09-testing`

Each step folder includes `CHANGES.md` narrative.

## Walkthrough

See [`docs/walkthrough.md`](docs/walkthrough.md) for the visual dependency-flow explanation, step mapping, troubleshooting, and pattern cheat sheet.
