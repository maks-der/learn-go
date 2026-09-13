# 18. Real Project Skills

## Description

A real Go program is more than a single `main` function. You must organize packages, load configuration, write logs, stop the process in a safe way, and ship a binary. This topic shows common project skills that teams use after the language basics.

The Go team does not define one official folder tree for applications. Teams still share common patterns. Use those patterns when they help. Do not copy a large layout for a small tool.

---

## Project layout (`cmd/`, `internal/`, `pkg/`)

The official module layout notes are on [https://go.dev/doc/modules/layout](https://go.dev/doc/modules/layout). That document is a guide. It is not a required standard.

A common application module looks like this:

- `go.mod` at the module root
- `cmd/appname/` for a `package main` that you build and run
- `internal/` for packages that other modules must not import
- optional `pkg/` for packages that other modules may import
- `api/`, `migrations/`, or `docs/` for files that are not Go packages

`cmd/` holds commands. Each subdirectory is one program. Example: `cmd/api` and `cmd/worker`. The directory name is often the binary name.

`internal/` is a compiler rule. A package outside the parent of `internal` cannot import a package under that `internal` tree. Use `internal/` for database access, configuration, and other private code.

`pkg/` is only a convention. The compiler does not treat `pkg/` as special. Some teams use `pkg/` for public library code. Other teams put public packages at the module root. Both are valid.

Keep the first version of a module small. One `main` package at the root is enough for a tiny tool. Add `cmd/` when you have two commands. Add `internal/` when you need a private package. Do not create empty layers such as `models/`, `utils/`, and `helpers/` without a clear job for each package.

Name packages by what they provide. A package name is the last element of the import path. Prefer `store` over `storepkg`. Prefer `httpapi` over `helpers`.

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

1. Add a second command `cmd/worker` that imports the same `internal` package as `cmd/api`. Build both binaries.
2. Create a second module that tries to import the first module `internal` package. Record the compiler error.
3. Move one public type to a `pkg/` package. Import it from `cmd/api`. State one reason to keep that type public.

#### Advanced practical tasks

1. Design a layout for a module with an HTTP API, a CLI admin tool, and a private database package. Write the tree, the import paths, and one sentence for the job of each package.
2. Compare your layout with one public Go repository. List three matches and two differences. Do not copy folders that you cannot explain.

---

## Configuration: environment variables, files, and twelve-factor

Configuration is data that changes between environments. Examples: listen address, database URL, and log level. Configuration is not application code. Do not hard-code environment-specific values in source files.

The Twelve-Factor App method stores configuration in the environment. Read the method at [https://12factor.net/config](https://12factor.net/config). Environment variables work in local shells, containers, and many host platforms.

In Go, read an environment variable with `os.Getenv` or `os.LookupEnv`. `LookupEnv` reports whether the variable is set. Parse numbers with `strconv`. Parse durations with `time.ParseDuration`. Validate values at startup. Stop the process if a required value is missing or invalid.

Configuration files are also common. JSON, TOML, and YAML files can hold structured settings. Load the file at startup. A file is useful for defaults that are not secret. Do not store secrets in a file that you commit to git.

A simple pattern is this order:

1. Set defaults in code.
2. Overlay values from a file if a file path is set.
3. Overlay values from environment variables.
4. Validate the final struct.

Keep one configuration struct. Pass that struct into constructors. Do not read `os.Getenv` from many packages after startup. Late reads make tests and audits hard.

Fail at startup when configuration is wrong. A process that starts with a bad database URL and fails on the first request is harder to operate than a process that exits at once.

### Questions

#### Theoretical questions

1. What is the difference between configuration and code?
2. Why does the Twelve-Factor App method prefer environment variables?
3. What extra information does `os.LookupEnv` give that `os.Getenv` does not give?
4. Why must you validate configuration at startup?
5. Why is it a problem to call `os.Getenv` from many packages during a request?

#### Easy practical tasks

1. Write a program that reads `LISTEN_ADDR` and prints it. If the variable is empty, use `:8080`.
2. Use `os.LookupEnv` for `LOG_LEVEL`. Print `set` or `unset`.
3. Parse `SHUTDOWN_TIMEOUT` with `time.ParseDuration`. Print an error if the value is invalid.
4. Write five sentences that separate secrets, defaults, and environment-specific values.

#### Medium practical tasks

1. Define a `Config` struct with listen address, log level, and database URL. Fill it from defaults, then from environment variables. Reject an empty database URL.
2. Add an optional JSON file overlay. Document the precedence: defaults, file, environment.
3. Write a test that sets environment variables, loads `Config`, and checks the fields. Restore the environment after the test.

#### Advanced practical tasks

1. Support two environments (`dev` and `prod`) with the same binary. Document which values come from the environment and which values stay as defaults.
2. Add a check that refuses to start when `LOG_LEVEL` is not one of `debug`, `info`, `warn`, or `error`. Show the exit code.

---

## Structured logging with `slog`

The `log/slog` package is the structured logger in the standard library (Go 1.21 and later). A structured log record has a message and key-value attributes. A machine can parse those attributes. A text line that mixes free words is harder to search.

Create a handler and a logger:

```text
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("server start", "addr", addr)
```

`slog.NewJSONHandler` writes JSON. Use JSON in production. `slog.NewTextHandler` writes key-value text. Text is easier to read on a laptop.

Levels are `Debug`, `Info`, `Warn`, and `Error`. Set the minimum level in `slog.HandlerOptions`. Do not log at `Info` for events that occur on every tight loop. Do not log secrets. Do not log access tokens or passwords.

`slog.SetDefault` installs a process-wide logger. `slog.Info` then uses that logger. Prefer to pass a `*slog.Logger` into your types. A field on a struct is easy to test.

Add stable keys. Examples: `request_id`, `method`, `path`, `duration_ms`, `err`. Use the same key names in every handler. `logger.With("request_id", id)` returns a logger that includes that attribute on every record.

`slog.LogValuer` lets a type control how it appears in a log. Use it to redact a secret type. `slog.Group` nests attributes under one key.

Write logs to standard output or standard error. The process that runs the binary collects those streams. Do not open a custom log file from the application unless you have a special requirement.

### Questions

#### Theoretical questions

1. What is a structured log record?
2. When do you choose a JSON handler over a text handler?
3. Why must you avoid logging secrets?
4. What does `logger.With` return?
5. Why is a logger field on a struct easier to test than a global logger?

#### Easy practical tasks

1. Write a program that logs `hello` with `slog.Info` and one string attribute.
2. Switch the handler from text to JSON. Show one output line for each handler.
3. Set `HandlerOptions.Level` to `slog.LevelWarn`. Log one `Info` record and one `Warn` record. Confirm that only the warning appears.
4. Add `request_id` with `With`. Log two messages. Confirm that both lines include the id.

#### Medium practical tasks

1. Create an HTTP middleware that logs method, path, status code, and duration. Use `slog`.
2. Define a `token` type that implements `slog.LogValuer` and redacts the value. Log it. Confirm that the secret is not visible.
3. Pass `*slog.Logger` into a service struct. Write a test that uses `slog.NewJSONHandler` with a `bytes.Buffer` and checks the buffer.

#### Advanced practical tasks

1. Add a `ReplaceAttr` function that renames the time key and removes the source path in production. Document the rule.
2. Compare `log.Printf` and `slog` on the same event. Write six sentences about search, levels, and attributes.

---

## Graceful shutdown

Graceful shutdown means the process stops accepting new work and then finishes in-flight work before it exits. The opposite is a hard stop. A hard stop can cut an HTTP response or leave a database transaction open.

Operating systems send signals. Interactive terminals send `SIGINT` on Ctrl+C. Container hosts send `SIGTERM`. Use `signal.NotifyContext` (Go 1.16 and later) or `signal.Notify` to receive those signals.

`http.Server.Shutdown` stops the listener and waits for active connections. Pass a context with a timeout. If the timeout ends, `Shutdown` returns an error. You may then call `Server.Close` for a hard stop.

A typical sequence is:

1. Start `ListenAndServe` in a goroutine.
2. Wait for a signal.
3. Call `Shutdown` with a timeout from configuration.
4. Stop background workers with `context` cancellation.
5. Close the database pool.
6. Exit with a non-zero code if shutdown fails.

`ListenAndServe` returns `http.ErrServerClosed` after a clean `Shutdown`. Treat that value as success. Treat other errors as start or stop failures.

Drain worker pools with `sync.WaitGroup`. Cancel the root context so that outbound HTTP clients and database queries stop. Set timeouts on those operations. Do not wait forever.

Do not call `os.Exit` from a goroutine that skips `defer`. `os.Exit` skips deferred cleanup. Return from `main` after you close resources.

### Questions

#### Theoretical questions

1. What does graceful shutdown wait for?
2. Which signal do many container hosts send to stop a process?
3. What does `http.Server.Shutdown` do to the listener?
4. Why does `ListenAndServe` return `http.ErrServerClosed` after a clean shutdown?
5. Why is `os.Exit` a problem when you still have `defer` cleanup?

#### Easy practical tasks

1. Write a program that listens for `SIGINT` with `signal.NotifyContext` and then prints `bye`.
2. Start `http.Server` on `:8080`. Handle `/health`. Shutdown the server from a signal. Use a 5 second timeout.
3. Log `shutdown start` and `shutdown done` with `slog`.
4. Write the sequence of steps for HTTP server, workers, and database close. Use six short sentences.

#### Medium practical tasks

1. Add an in-flight counter with `sync.WaitGroup` around a slow handler. Confirm that `Shutdown` waits until the handler returns.
2. Make the handler sleep longer than the shutdown timeout. Record the `Shutdown` error. Document what happens to the client.
3. Cancel a background loop from the same context that the signal creates. Show that the loop stops.

#### Advanced practical tasks

1. Combine HTTP shutdown, a worker `WaitGroup`, and `sql.DB.Close`. Measure the time from signal to process exit.
2. Document the stop behavior of your program in a runbook: signals, timeout, exit codes, and what a load balancer must do.

---

## Dependency injection without frameworks

Dependency injection means a type receives the values that it needs. The type does not construct every dependency itself. In Go, the common method is a constructor function.

```text
func NewServer(store Store, log *slog.Logger) *Server
```

`main` (or a small `cmd` package) creates the concrete database, the logger, and the HTTP server. That is the composition root. Lower packages accept interfaces or small structs.

Define interfaces in the package that uses them when the interface is small. A store interface with `Get` and `Save` is enough for an HTTP handler. Do not create an interface for every struct in advance.

Frameworks such as `wire`, `fx`, and `dig` generate or resolve graphs of constructors. They are optional. A team with a small service often does not need them. A long list of constructors in `main` is clear and easy to debug.

Benefits of constructor injection:

- Tests pass a fake `Store`.
- The dependency list is visible in the function signature.
- You can see the graph without a container.

Do not hide required dependencies in package-level variables. Global clients make tests share state. If you need a default, provide `NewServer` plus a small `main` that builds production values.

Pass `context.Context` as the first argument of methods that do I/O. Do not store a request context on a long-lived server struct.

### Questions

#### Theoretical questions

1. What is a composition root in a Go program?
2. Why does a constructor signature help tests?
3. Where should a small interface live when only one package uses it?
4. What problem do package-level client variables create?
5. Why must you not store a request `context.Context` on a long-lived struct?

#### Easy practical tasks

1. Write `NewGreeter(log *slog.Logger)` and a `Greet` method that logs a name. Call it from `main`.
2. Extract a `Clock` interface with `Now() time.Time`. Inject a fake clock in a test.
3. List the constructors in a small module. Draw arrows from `main` to each type.
4. Replace one global `var db` with a field on a struct. Show the new constructor.

#### Medium practical tasks

1. Build an HTTP handler that depends on a `UserStore` interface. Write a fake store and a `httptest` test.
2. Add a third dependency (logger) to the handler constructor. Update all call sites and tests.
3. Compare constructor injection with a `dig` or `fx` example from public docs. Write five sentences on clarity and cost.

#### Advanced practical tasks

1. Split composition into `internal/app.New` that returns a ready `http.Handler` plus a `Close` function. Keep `main` thin.
2. Find a circular import that appears when two packages construct each other. Break the cycle with an interface or a move of types. Record the steps.

---

## Command-line applications (`flag`, optional Cobra)

A command-line application reads arguments and flags, then performs one job. The standard library package is `flag`.

`flag.String`, `flag.Int`, and `flag.Bool` define flags. Call `flag.Parse` before you read values. `flag.Args` returns remaining arguments after flags.

```text
out := flag.String("o", "out.txt", "output file")
flag.Parse()
```

`flag` prints help when the user passes `-h`. Write a clear usage string. Set `flag.Usage` if the default text is not enough.

`os.Args` is the raw argument slice. `os.Args[0]` is the program name. Prefer `flag` over manual parsing when you need typed flags.

One module can hold several commands under `cmd/`. Each command has its own `package main`. That layout avoids a large flag set in one binary.

Cobra (`github.com/spf13/cobra`) is a third-party library for nested subcommands, persistent flags, and generated help. Use Cobra when you have many subcommands and a long help tree. Learn `flag` first. Do not add Cobra to a tool with two flags.

Exit codes matter. Exit `0` on success. Exit a non-zero code on usage errors and runtime errors. `flag` calls `os.Exit(2)` on parse errors by default.

Print user errors to standard error. Print normal program output to standard output. That split lets users pipe output to another program.

### Questions

#### Theoretical questions

1. When must you call `flag.Parse`?
2. What does `flag.Args` return?
3. Why do exit codes matter to scripts?
4. When is Cobra a reasonable choice?
5. Why must user errors go to standard error?

#### Easy practical tasks

1. Write a program with a `-name` string flag and a default. Print `hello` and the name.
2. Add a `-n` integer flag. Print the name that many times.
3. Run the program with `-h`. Save the help text.
4. Create `cmd/greet` in a module. Build the binary with `go build -o greet ./cmd/greet`.

#### Medium practical tasks

1. Implement two subcommands with `flag.NewFlagSet` (`greet` and `list`). Parse `os.Args` yourself. Print usage on a bad command.
2. Add required positional arguments. Exit with a non-zero code when they are missing.
3. Write a test that calls a `Run(args []string) error` function. Do not start a subprocess.

#### Advanced practical tasks

1. Add Cobra to a copy of the tool with three nested commands. Document one feature that `flag` does not give you.
2. Compare binary size and start time of the `flag` version and the Cobra version. Write the two measurements.

---

## REST and JSON services

A REST-style HTTP service exposes resources with URLs and HTTP methods. JSON is a common body format. The standard library packages are `net/http` and `encoding/json`.

Go 1.22 and later extend `http.ServeMux`. You can register a method and a path pattern:

```text
mux.HandleFunc("GET /items/{id}", getItem)
```

`r.PathValue("id")` returns the path segment. Use these patterns before you add a third-party router.

Decode JSON with `json.NewDecoder(r.Body)`. Encode JSON with `json.NewEncoder(w)`. Set `Content-Type` to `application/json`. Set a status code before you write the body. `Encode` writes the body.

Limits protect the server. Wrap `r.Body` with `http.MaxBytesReader`. Use `Decoder.DisallowUnknownFields` when the API must reject unknown keys. Check `Decode` errors. Distinguish invalid JSON from a valid empty object.

Keep handlers thin. Parse the request. Call a service method with `r.Context()`. Write the response. Return after an error write. Do not continue the handler after you write a 4xx or 5xx body.

Use `GET` for reads that do not change state. Use `POST`, `PUT`, `PATCH`, and `DELETE` with clear rules. Document status codes: `200`, `201`, `204`, `400`, `401`, `404`, `409`, and `500` are common.

`httptest.NewRecorder` and `httptest.NewRequest` test handlers without a listen socket. Test the JSON body and the status code.

### Questions

#### Theoretical questions

1. How do you read a path parameter with the Go 1.22 `ServeMux`?
2. Why must you set the status code before you write the body?
3. What does `http.MaxBytesReader` protect against?
4. Why keep HTTP handlers thin?
5. Which HTTP methods must not change server state in a strict REST design?

#### Easy practical tasks

1. Register `GET /health` that writes `{"status":"ok"}` as JSON.
2. Register `GET /items/{id}` that writes the `id` in a JSON object.
3. Register `POST /items` that decodes a JSON name field and returns status `201`.
4. Send a request with `httptest` and check the status code.

#### Medium practical tasks

1. Reject request bodies larger than 1 MiB. Return status `413` or `400` and a JSON error object.
2. Use `DisallowUnknownFields` on `POST`. Write a test that sends an extra field and expects status `400`.
3. Split routing, JSON helpers, and a service method into packages. Test the service without HTTP.

#### Advanced practical tasks

1. Implement list, get, create, and delete for one resource in memory. Document every status code. Add table-driven `httptest` cases.
2. Add a middleware that sets `Content-Type` and recovers from a panic with status `500` and a log line. Test both paths.

---

## Authentication basics (JWT and sessions)

Authentication answers the question: who is this client? Authorization answers the question: what may this client do? Keep those words distinct.

A session is server-side state. The server creates a session record after a successful login. The browser stores a session identifier in a cookie. Later requests send the cookie. The server looks up the session.

Cookie flags matter. Use `HttpOnly` so that script cannot read the cookie. Use `Secure` so that the cookie is sent only on HTTPS. Set `SameSite` to reduce cross-site request abuse. Set a path and an expiry.

A JSON Web Token (JWT) is a compact signed token. A common form has three parts: header, claims, and signature. The parts are Base64-URL strings separated by dots. The signature lets the server detect a change to the claims.

JWT claims often include `sub` (subject), `exp` (expiry), `nbf` (not before), `iss` (issuer), and `aud` (audience). Validate all of those fields. Reject an expired token. A JWT is not encrypted by default. Any holder can read the claims. Do not put a password in a JWT.

Sessions need a store and a revoke path. JWT access tokens are often stateless. Revoke is harder unless you keep a block list or use short expiry plus a refresh token. Short expiry reduces risk when a token leaks.

Transport must be HTTPS in production. Send tokens in the `Authorization` header as `Bearer` for APIs, or use cookies for browsers. Do not put tokens in query strings. Query strings appear in logs.

Password storage is a separate topic. Store a password hash with a slow hash function from a maintained library. Never store a plaintext password. Never log a password.

This section is conceptual. Do not invent your own crypto. Use maintained libraries and a written threat model.

### Questions

#### Theoretical questions

1. What is the difference between authentication and authorization?
2. Where does a session identifier live on the client?
3. What are the three parts of a common JWT?
4. Why is a JWT not private by default?
5. Why must tokens stay out of query strings?

#### Easy practical tasks

1. Write a table with columns `session` and `JWT`. Add five rows (state, revoke, expiry, size, server memory).
2. List the cookie flags `HttpOnly`, `Secure`, and `SameSite`. Write one sentence for each flag.
3. Decode a public sample JWT header and payload on [https://jwt.io](https://jwt.io) in the debugger with a known tutorial token. Write the claim names that you see. Do not use a real production token.
4. Write four sentences that explain `exp` and `aud`.

#### Medium practical tasks

1. Design login and logout for a session cookie API. Write the request flow in eight steps. Include cookie flags.
2. Design an API that uses a short-lived access JWT and a longer-lived refresh token. Write who stores each token and how revoke works.
3. List three mistakes: token in a query, no `exp`, and a secret in the repository. Write the safer alternative for each mistake.

#### Advanced practical tasks

1. Compare session cookies and JWT for a browser app and for a machine client. Write a one-page recommendation with reasons.
2. Read the JWT and cookie documents that your team must follow. Write a checklist of validation steps for incoming tokens.

---

## Docker images for a Go binary

A container image packages a binary and the files that the binary needs. Go can produce a static binary. A static binary does not need `libc` in the image. Set `CGO_ENABLED=0` when you build for a scratch image.

A multi-stage Dockerfile builds in a Go image and copies only the binary into a small runtime image:

```text
FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/app ./cmd/app

FROM scratch
COPY --from=build /out/app /app
ENTRYPOINT ["/app"]
```

`scratch` is empty. There is no shell. There is no CA certificate file unless you copy one. HTTPS clients need CA certificates. Distroless images such as `gcr.io/distroless/static` include certificates and a non-root user. Distroless images still have no package manager and no shell.

Do not run the process as root. Set `USER` to a non-root user in the final image. Distroless static images provide a non-root user for that purpose.

Add a `.dockerignore` file. Exclude `.git`, testdata, and local binaries. A smaller build context makes the build faster.

Pass configuration with environment variables at runtime. Do not bake production secrets into the image. Pin base image versions. Rebuild when you need security updates of the build image.

The binary that you copy must match the container OS and architecture. A Windows binary does not run in a Linux image. See the next section on `GOOS` and `GOARCH`.

### Questions

#### Theoretical questions

1. Why does `CGO_ENABLED=0` help a `scratch` image?
2. What is missing from `scratch` that an HTTPS client may need?
3. Why do teams use a multi-stage build?
4. Why must the process not run as root in the container?
5. Why must you pin base image versions?

#### Easy practical tasks

1. Write a `.dockerignore` that excludes `.git` and local output binaries.
2. Build a static binary with `CGO_ENABLED=0`. Write the command and the binary size.
3. Read the Distroless documentation. Write two differences between `scratch` and Distroless static.
4. List four files that must not appear in a production image.

#### Medium practical tasks

1. Write a multi-stage Dockerfile for `cmd/app`. Build the image. Run a `/health` request against a published port.
2. Switch the final stage from `scratch` to Distroless static. Run as a non-root user. Confirm the process user.
3. Add CA certificates to a `scratch` image. Call an HTTPS URL from the program. Record success or the error.

#### Advanced practical tasks

1. Produce two images for `linux/amd64` and `linux/arm64`. Document the `docker build` platform flags that you use.
2. Scan the image with a vulnerability tool that you have. Write which layer owns each finding. Plan a rebuild.

---

## Cross-compilation with `GOOS` and `GOARCH`

Go cross-compiles by default when `CGO_ENABLED=0`. Set `GOOS` for the target operating system. Set `GOARCH` for the target architecture.

```text
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app-linux-amd64 ./cmd/app
```

Common pairs:

- `windows` / `amd64`
- `linux` / `amd64`
- `linux` / `arm64`
- `darwin` / `arm64`
- `darwin` / `amd64`

List all pairs with `go tool dist list`. Show the current pair with `go env GOOS GOARCH`.

The `go` tool adds `.exe` to the output name when `GOOS=windows` and you do not set `-o`. Set `-o` when you build several targets in one script.

Cgo changes the story. Cgo needs a C compiler for the target. Cross-compilation with Cgo is harder. Prefer pure Go builds for release binaries when you can. Build tags `netgo` and `osusergo` request the pure Go resolver and user lookup on some platforms.

You do not need to run a foreign binary on your machine to produce it. You still must test the binary on the target OS or in a container that matches the target.

`GOARM` selects the ARM version for 32-bit ARM. `GOAMD64` selects the `amd64` microarchitecture level (Go 1.18 and later). Set those variables only when you need a special target.

Do not mix `GOOS` from the environment with a forgotten old value. Scripts must set `GOOS` and `GOARCH` in the same command that runs `go build`.

### Questions

#### Theoretical questions

1. What do `GOOS` and `GOARCH` select?
2. Why does `CGO_ENABLED=0` make cross-compilation easier?
3. What command lists supported `GOOS`/`GOARCH` pairs?
4. Why must you test a Linux binary on Linux or in a Linux container?
5. What is the role of `GOAMD64`?

#### Easy practical tasks

1. Run `go env GOOS GOARCH`. Write the values.
2. Run `go tool dist list`. Count the lines. Write five pairs that you recognize.
3. Cross-compile the same `main` package for `linux/amd64` and `windows/amd64`. Write the two file sizes.
4. Show `go help build` text that mentions `-o`. Write how you name each target binary.

#### Medium practical tasks

1. Write a script that builds four pairs and writes binaries into `dist/`. The script must set `CGO_ENABLED=0`.
2. Enable Cgo on a tiny program that includes `import "C"`. Try a cross-build. Record the error.
3. Build with `GOAMD64=v1` and the default. Compare file sizes if they differ. Document the target machines.

#### Advanced practical tasks

1. Cross-compile for `linux/arm64`. Run the binary in an emulator or on ARM hardware. Record the method and the output.
2. Document a release checklist: `GOOS`, `GOARCH`, Cgo, `go version`, `git` commit, and test host.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `cmd/`, `internal/`, configuration at startup, and `main` as a composition root work together in one module?
2. Why do structured logs, graceful shutdown, and twelve-factor configuration all help an operator?
3. When do you stay with `flag` and `net/http`, and when do you add Cobra or a third-party router?
4. What is the relationship between a static Linux binary, `scratch`, and `GOOS`/`GOARCH`?
5. How do sessions and JWT differ in revoke and server memory, and how does that choice change the API design?

#### Easy practical tasks

1. Create a module with `cmd/api`, `internal/config`, and `internal/logx`. Load `LISTEN_ADDR`, log JSON with `slog`, and serve `GET /health`.
2. Write a one-page cheat sheet: layout, `os.LookupEnv`, `slog` handlers, `Shutdown`, `flag.Parse`, `PathValue`, cookie flags, `CGO_ENABLED=0`, `GOOS`.
3. Build the API for `linux/amd64` with `CGO_ENABLED=0`. Write the binary size and the `go version` line.
4. Draw a sequence diagram from SIGTERM to `Shutdown` to process exit. Include the timeout.

#### Medium practical tasks

1. Add constructor injection for a store interface, `httptest` for `POST /items`, and a fake store. Stop the test server with `Shutdown`.
2. Write a Dockerfile with a build stage and a Distroless final stage. Pass `LISTEN_ADDR` at runtime. Hit `/health`.
3. Add a CLI `cmd/healthcheck` that calls `/health` with `net/http` and exits `1` on failure. Document both commands.

#### Advanced practical tasks

1. Package a small service that uses env configuration, `slog`, graceful shutdown, one JSON resource, and a multi-stage image. Write an operator README with signals, flags, and environment variables.
2. Cross-compile that service for `linux/amd64` and `linux/arm64`. Build two images. Write a table of binary size, image size, and how you tested each target.
