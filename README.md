# go-url-shortener

Minimal URL shortener for teaching Uber Fx dependency injection with a visual walkthrough.

## Start here

- Walkthrough: [`/docs/walkthrough.md`](/docs/walkthrough.md)
- Final runnable app: [`/cmd/shortener/main.go`](/cmd/shortener/main.go)
- Progressive checkpoints: [`/checkpoints`](/checkpoints)

## One-command run

```bash
go run ./cmd/shortener
```

Optional named-secondary storage demo:

```bash
USE_SECONDARY_STORE=1 go run ./cmd/shortener
```

## One-command test

```bash
go test ./...
```

## API

### Create short URL

```bash
curl -s -X POST http://localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://go.dev"}'
```

Response example:

```json
{"slug":"u1","short_url":"http://localhost:8080/u1"}
```

### Resolve slug

```bash
curl -i http://localhost:8080/u1
```

Expected: `307 Temporary Redirect` to original URL.

## Learning map

See `/docs/walkthrough.md` for:
- before/after wiring narrative
- step-by-step visuals (mermaid graphs)
- troubleshooting common Fx errors
- cheat sheet and slide-to-code index
