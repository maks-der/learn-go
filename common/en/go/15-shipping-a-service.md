# 15. Shipping a Service

## Description

A real Go program is more than a single `main` function. You must organize packages, load configuration, write logs, stop the process in a safe way, and ship a binary. This topic shows how you ship a service.

The Go team does not define one official folder tree for applications. Teams still share common patterns. Use those patterns when they help. Do not copy a large layout for a small tool.

Use one term for each concept. Configuration is data that you load at start. Graceful shutdown is a stop that finishes in-flight work. A secret is a password or a key. Do not store secrets in source files.

---

## `cmd/` and `internal/` layout

The official module layout notes are on [https://go.dev/doc/modules/layout](https://go.dev/doc/modules/layout). That document is a guide. It is not a required standard.

A common application module looks like this:

- `go.mod` at the module root
- `cmd/appname/` for a `package main` that you build and run
- `internal/` for packages that other modules must not import
- optional `pkg/` for packages that other modules may import
- `api/`, `migrations/`, or `docs/` for files that are not Go packages

`cmd/` holds commands. Each subdirectory is one program. Example: `cmd/api` and `cmd/worker`. The directory name is often the binary name. Keep `main` small. Call `run()` from `internal` or from the same package.

`internal/` is a compiler rule. A package outside the parent of `internal` cannot import a package under that `internal` tree. Use `internal/` for database access, configuration, and other private code. Topic 8 covers the import fence.

`pkg/` is only a convention. The compiler does not treat `pkg/` as special.

Keep the first version of a module small. One `main` package at the root is enough for a tiny tool. Add `cmd/` when you have two commands. Add `internal/` when you need a private package. Do not create empty layers such as `models/`, `utils/`, and `helpers/` without a clear job for each package.

Name packages by what they provide. Prefer `store` over `storepkg`. Prefer `httpapi` over `helpers`.

```text
example.com/shop/
  go.mod
  cmd/api/main.go
  cmd/worker/main.go
  internal/store/store.go
  internal/httpapi/server.go
```

### Questions

#### Theoretical questions

1. Does the Go toolchain require a `cmd/` folder? Explain.
2. What compiler rule does `internal/` enforce?
3. Why is `pkg/` different from `internal/`?
4. When is a single `main` package at the module root enough?
5. What is a good reason to add a second directory under `cmd/`?

#### Easy practical tasks

1. Draw a folder tree for a module with `cmd/api`, `internal/store`, and `go.mod`. Write the import path of `store`.
2. Create that tree in a new module. Put a function in `internal/store`. Call it from `cmd/api`. Build the command.
3. Read [https://go.dev/doc/modules/layout](https://go.dev/doc/modules/layout). Write four facts from that page.
4. List three package names that describe a job. List three package names that do not describe a job.

#### Medium practical tasks

1. Split `cmd/api` and `cmd/worker`. Share `internal/store`. Build both binaries.
2. Try to import `internal/store` from a second module. Record the error.
3. Move logic out of `main` into `func run() error`. Make `main` call `os.Exit` only after `run` returns an error.

#### Advanced practical tasks

1. Design a module with an API command, a migrate command, and `internal` packages for config, store, and HTTP. Write the tree and import paths.
2. Compare a tiny tool at the module root with the `cmd/` layout. Write six sentences on when you grow into `cmd/`.

---

## Configuration and structured logs

Configuration is data that the process reads at start. Sources include flags, environment variables, and files.

Use `flag` for a small CLI. Use environment variables for a service in a container. Example: `PORT`, `DATABASE_URL`, `LOG_LEVEL`. Read them in one function. Return a `Config` struct.

```go
type Config struct {
	Addr string
	DSN  string
}

func Load() (Config, error) {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	return Config{Addr: ":" + addr, DSN: dsn}, nil
}
```

Do not read `os.Getenv` in every package. Load once in `main` or in `internal/config`. Pass `Config` down.

Validate required fields at start. Fail fast. Do not start the HTTP server when the DSN is empty.

Structured logs use `log/slog`. Topic 9 introduces `slog`. Install a JSON handler in production. Use text in local development if you prefer.

```go
h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
slog.SetDefault(slog.New(h))
slog.Info("listen", "addr", cfg.Addr)
```

Log a request id if you have one. Log errors with `"err", err`. Do not log secrets. Do not log full request bodies that can contain passwords.

Set the log level from configuration. `DEBUG` is for development. `INFO` is the default for production.

### Questions

#### Theoretical questions

1. Why do you load configuration in one function?
2. What must happen when a required environment variable is missing?
3. Why is `slog` better than `fmt.Println` in a service?
4. What must you never put in a log record?
5. Why do containers often use environment variables?

#### Easy practical tasks

1. Write `Load` that reads `PORT` with a default and requires `DATABASE_URL`. Print the config without the password if you mask it.
2. Log `started` with `slog.Info` and the listen address.
3. Parse `LOG_LEVEL` as `info` or `debug`. Set `HandlerOptions.Level`.
4. List five environment variable names for a typical API.

#### Medium practical tasks

1. Combine `flag` and environment: a flag overrides an environment value. Document the order.
2. Load a JSON config file when `CONFIG` is set. Still allow `PORT` to override the file.
3. Pass `*slog.Logger` into `internal/httpapi` instead of the global logger. Log one request.

#### Advanced practical tasks

1. Redact `DATABASE_URL` in logs and in `fmt.Sprintf("%+v", cfg)`. Show the redacted string.
2. Write a config policy of one page: sources, precedence, required keys, and log rules.

---

## Graceful shutdown

Graceful shutdown stops accept of new work. It waits for in-flight requests. Then it exits.

`http.Server.Shutdown` is the standard method. It closes the listener. It waits until handlers return or the context is done.

```go
s := &http.Server{Addr: cfg.Addr, Handler: mux}

go func() {
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("listen", "err", err)
	}
}()

ch := make(chan os.Signal, 1)
signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
<-ch

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := s.Shutdown(ctx); err != nil {
	return err
}
```

`signal.Notify` must use a buffered channel. `SIGTERM` is common in containers. `SIGINT` is Ctrl+C.

After HTTP shutdown, close the database pool. Cancel worker contexts. Wait for worker `WaitGroup`. Order the stop: stop accept first, then finish work, then close dependencies.

`os.Exit` skips deferred closes. Prefer `return` from `run` so that `defer db.Close()` runs.

`Shutdown` does not stop a handler that ignores the request context. Handlers must watch `r.Context()`. Topic 10 covers context.

Set `http.Server` timeouts so that a stuck client does not block shutdown forever. The shutdown context is a last bound.

### Questions

#### Theoretical questions

1. What does `http.Server.Shutdown` wait for?
2. Why must the signal channel have capacity at least `1`?
3. Which signals do you listen for in a container?
4. Why is `os.Exit` a problem during shutdown?
5. Why must a handler use `r.Context()`?

#### Easy practical tasks

1. Write a server that logs `bye` after `Shutdown`. Send SIGINT and confirm the log.
2. Print `http.ErrServerClosed` handling in your listen goroutine.
3. Add a 2 second shutdown timeout. Write why you chose that number.
4. Draw the stop order: signal, Shutdown, db.Close, return.

#### Medium practical tasks

1. Add a handler that sleeps 3 seconds. Shutdown with a 1 second timeout. Record the error. Then increase the timeout.
2. Cancel a worker context when the signal arrives. Wait with a `WaitGroup`.
3. Show that `os.Exit(0)` in the signal handler skips `defer` in `run`. Then fix the design.

#### Advanced practical tasks

1. Implement `run` that starts HTTP and a worker. Shutdown both. Pass `-race`.
2. Read `go doc http.Server.Shutdown`. Write five rules for handlers and idle connections.

---

## Docker and cross-compilation

Go compiles to a static binary for many operating systems. You set `GOOS` and `GOARCH`.

```text
set GOOS=linux
set GOARCH=amd64
go build -o api ./cmd/api
```

On Unix:

```text
GOOS=linux GOARCH=amd64 go build -o api ./cmd/api
```

You do not need CGO for most services. Set `CGO_ENABLED=0` for a static binary that copies into a small image.

A typical Dockerfile uses a build stage and a runtime stage:

```text
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

FROM gcr.io/distroless/static:nonroot
COPY --from=build /api /api
USER nonroot
ENTRYPOINT ["/api"]
```

Copy `go.mod` first so that dependency download caches. Do not run as root in the final image. Distroless and scratch images have no shell. Pass configuration with environment variables.

Pin base image versions. Scan images in CI when you can.

Cross-compile on your machine for a remote OS. You do not need Docker only to produce a Linux binary. Docker is for a repeatable runtime.

`go build -trimpath` removes local paths from the binary. `-ldflags="-s -w"` reduces size. Use them when you know the trade-off (harder debug).

### Questions

#### Theoretical questions

1. What do `GOOS` and `GOARCH` select?
2. Why do you set `CGO_ENABLED=0` for many service images?
3. Why does a Dockerfile copy `go.mod` before the full source?
4. Why must the process not run as root in the final image?
5. Can you produce a Linux binary without Docker?

#### Easy practical tasks

1. Cross-compile `cmd/api` for `linux/amd64`. Show the file size. You do not need to run it on Windows.
2. Run `go env GOOS GOARCH`. Write the values.
3. Write a two-stage Dockerfile as text in your notes. Label build and runtime.
4. List three environment variables that the container must receive.

#### Medium practical tasks

1. Build a local image and run it with `PORT`. Curl `/health` if you have Docker.
2. Compare image size of `golang` as runtime versus `distroless` or `scratch`. Write the two sizes.
3. Add `-trimpath` to the build. Confirm the binary still runs.

#### Advanced practical tasks

1. Multi-arch build notes: `amd64` and `arm64`. Write the `GOARCH` values and one CI command.
2. Produce a binary with `CGO_ENABLED=1` and a C dependency. Record why the scratch image fails. Then use a glibc or distroless-cc base.

---

## Timeouts, retries, and secrets

A service must bound time. Set HTTP server timeouts. Set HTTP client timeouts. Set context timeouts on outbound calls and on database queries.

```go
client := &http.Client{Timeout: 5 * time.Second}
ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
defer cancel()
```

The client timeout includes connect, redirects, and body read. The context timeout cancels the request. Use both when you need cancel of the whole tree.

A retry repeats a call after a failure. Retry only when the operation is idempotent or when you know a repeat is safe. GET is often safe. POST that creates a resource is not safe without an idempotency key.

Use exponential backoff and a maximum attempt count. Add jitter so that many clients do not retry at the same time. Do not retry on `400` or `401`. Retry on timeout and on `503` when the API documents it.

```go
for attempt := 0; attempt < 3; attempt++ {
	err := call(ctx)
	if err == nil {
		return nil
	}
	if !retryable(err) {
		return err
	}
	time.Sleep(backoff(attempt))
}
```

Honor `ctx.Done()` in the sleep. A retry must not ignore cancel.

A secret is a password, a token, or a private key. Load secrets from the environment or from a secret store. Do not commit `.env` files that contain real secrets. Do not put secrets in logs, in error messages, or in client-facing JSON.

Rotate secrets without a code change when you can. Use a new environment value and a restart.

TLS files are secrets. File permissions must not be world-readable.

### Questions

#### Theoretical questions

1. What does `http.Client.Timeout` include?
2. When is a retry unsafe?
3. What is jitter on a backoff?
4. Where do you load secrets?
5. Why must you honor `ctx.Done()` during a retry sleep?

#### Easy practical tasks

1. Set a client timeout of 1 millisecond against a slow URL. Print the error.
2. Make a two-column table: "Retry" and "Do not retry". Add four status or error rows.
3. Write `Load` that requires `API_TOKEN`. Fail if it is empty. Do not print the token.
4. List three places that must not receive a secret (logs, Git, error JSON).

#### Medium practical tasks

1. Write `backoff(attempt)` that doubles and adds a small random jitter. Print three sleeps.
2. Retry a GET three times on `503`. Stop on `200` or on cancel.
3. Show a wrapped error that includes a DSN. Then redact it.

#### Advanced practical tasks

1. Implement retries that share a parent context of 2 seconds. Show that attempts stop when the context is done.
2. Write a secret checklist for a pull request: env names, log review, Docker, and TLS files.

---

## Input validation

Validate input at the edge. The HTTP handler or the CLI is the edge. Do not trust query parameters, JSON fields, or headers.

Check:

- required fields
- length limits
- format (email, UUID, enum)
- ranges (page size, ids)
- unexpected fields when the API is strict

```go
if in.Name == "" || len(in.Name) > 100 {
	http.Error(w, "bad name", http.StatusBadRequest)
	return
}
```

Return `400` for input that the client can fix. Return `404` when an id is well formed but missing. Do not use `500` for validation.

Validate before you talk to the database. A failed validation must not open a transaction.

Limit lists. Reject `page_size=1000000`. Set a maximum.

Use `http.MaxBytesReader` for JSON bodies. Topic 13 covers that limit. Use `DisallowUnknownFields` when extra keys are a mistake.

For HTML forms, encode output with `html/template` so that scripts do not run. For JSON APIs, still bound strings.

Do not validate only in the database. The database constraint is a last stop. The handler still returns a clear error.

Keep validation functions pure when you can. They are easy to test. Topic 11 covers table-driven tests.

### Questions

#### Theoretical questions

1. Where do you validate input?
2. What status code do you use for a bad field?
3. Why do you validate before a transaction?
4. Why do you cap `page_size`?
5. Why is a database constraint not enough?

#### Easy practical tasks

1. Reject an empty `name` with `400`. Test with `httptest`.
2. Reject `name` longer than 20 runes. Test both sides of the limit.
3. Reject `page_size` below `1` or above `100`.
4. Write five validation rules for a `User` JSON object.

#### Medium practical tasks

1. Table-test a `func validateItem(Item) error` with empty, too long, and ok cases.
2. Combine JSON decode, `DisallowUnknownFields`, and `validateItem`. Return JSON errors.
3. Show a request that skips validation and causes a driver error. Then add validation and show `400`.

#### Advanced practical tasks

1. Validate a nested JSON list. Fail on the first bad index. Name the index in the error.
2. Write a validation policy: field limits, UTF-8, HTML vs JSON, and how tests cover the edge.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `cmd/`, `internal/`, and a small `main` work together in a shippable module?
2. What must `run` do from load config to listen to shutdown to close the database?
3. How do structured logs, redaction, and environment configuration protect operations?
4. Why do timeouts, retries, and validation all belong at the process edge?
5. What does a container image add that cross-compilation alone does not add?

#### Easy practical tasks

1. Create `cmd/api` with `/health`, `slog`, and `PORT`. Build a Linux binary.
2. Write a cheat sheet: layout, config, slog, Shutdown, GOOS, timeout, secret, validate.
3. Draw the lifecycle: start, load, migrate, listen, signal, shutdown, exit.
4. Mask a DSN in a log line. Show the printed text.

#### Medium practical tasks

1. Implement `run() error` with config, mux, graceful shutdown, and `db.Close` (use SQLite if you have topic 14).
2. Add one retrying outbound GET with a client timeout and a context bound.
3. Write a Dockerfile and a `.dockerignore`. List what you exclude.

#### Advanced practical tasks

1. Ship a tiny service: layout, JSON API with validation, slog, shutdown, and a multi-stage Docker build. Document env vars.
2. Run `-race` on a shutdown test that starts the server and sends SIGINT in-process or calls `Shutdown` directly.
