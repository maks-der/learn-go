# 21. Production, Security, and Architecture

## Description

A service that works on a laptop still needs a production design. You must version the API, observe the process, bound time, validate input, and protect secrets. You must also keep the package layout small enough that a new teammate can read it.

This topic is defensive. Use it to reduce accidents and abuse. Do not use it to plan attacks. Follow the law and your employer rules.

---

## API versioning and compatibility

An API is a contract. Clients depend on URLs, fields, status codes, and error shapes. A breaking change forces every client to change at the same time. That cost is high.

A compatible change adds a field or an endpoint. An old client ignores the new field. A breaking change removes a field, renames a field, changes a meaning, or turns a success into an error for the same request.

Common version methods:

- Path version: `/v1/items` and `/v2/items`.
- Header version: `Accept: application/vnd.myapp.v1+json`.
- A single URL with additive evolution and no version number for a small internal API.

Path versions are easy to see in logs. Header versions keep URLs stable. Pick one method and document it.

Rules that keep v1 stable:

- Do not remove or rename JSON fields.
- Do not change the type of a field.
- Do not reuse a field name for a new meaning.
- Add new fields as optional.
- Add new endpoints instead of overloading an old one with a new meaning.

When you must break the contract, publish v2. Keep v1 until clients move. Set a sunset date. Log the version that each client uses.

Write compatibility tests. Store example request and response files. Replay them in CI. Consumer-driven contracts help when many teams share the API.

Go struct tags must stay stable. A `json:"name,omitempty"` change can hide a field. Treat tag changes as API changes.

### Questions

#### Theoretical questions

1. What is a breaking change in a JSON HTTP API?
2. What is an additive change?
3. Why do teams keep `/v1` after they ship `/v2`?
4. How can a `json` tag change break a client?
5. What is the benefit of replay tests on stored example responses?

#### Easy practical tasks

1. Write a table of five changes. Mark each change as compatible or breaking.
2. Design `/v1/users` and `/v2/users` for a user resource. List fields that move between versions.
3. Document your version method (path or header) in ten lines.
4. Find the version policy of one public API. Write three rules from that policy.

#### Medium practical tasks

1. Implement two handlers for `/v1/items` and `/v2/items` that share one store. v2 adds one field. Test both JSON shapes.
2. Add a sunset header on v1. Log a warning when a client calls v1.
3. Write CI tests that fail when a required v1 field disappears from the encoder.

#### Advanced practical tasks

1. Plan a six-month migration from v1 to v2. Include metrics, dual-write if needed, and an end date for v1.
2. Compare path versioning with header versioning for a browser app and for a machine client. Write a one-page choice.

---

## Observability: metrics, tracing, and OpenTelemetry

Observability is the ability to explain the current behavior of a system from its outputs. The three common signals are logs, metrics, and traces.

Logs are events. Topic 18 covers `slog`. Metrics are numbers over time: request count, error count, latency, queue depth. Traces follow one request across services. A trace has a trace id. Each step is a span.

OpenTelemetry is a vendor-neutral set of APIs and exporters. The Go modules live under `go.opentelemetry.io/otel`. You create a `Tracer` and a `Meter`. You start a span in a handler. You pass `context.Context` so that child spans stay in the same trace.

Context propagation sends the trace id to the next service. HTTP uses headers such as W3C `traceparent`. Instrument the client and the server. If only one side is instrumented, the trace breaks.

Prometheus is a common metrics system. Many Go services expose `/metrics` for Prometheus to scrape. The `expvar` package and `runtime/metrics` expose process facts. Use `runtime/metrics` for GC and memory. Do not invent your own scrape format when Prometheus is the team standard.

Cardinality is the number of unique label combinations. Do not put a user id or a raw URL path with ids into a metric label. High cardinality can overload the metrics system.

Sample traces in production when traffic is high. Keep all error traces if you can. Align log fields with `trace_id` so that you can jump from a log line to a trace.

Do not expose metrics or pprof on a public listener. Bind debug endpoints to an internal port or a private network.

### Questions

#### Theoretical questions

1. What is the difference between a metric and a trace?
2. What does a span represent?
3. Why must you pass `context.Context` into traced functions?
4. What is metric cardinality, and why is a user id a bad label?
5. Why must `/metrics` stay off the public internet?

#### Easy practical tasks

1. List five metrics for an HTTP service (include at least one latency metric).
2. Read the OpenTelemetry Go getting-started page. Write the purpose of a `TracerProvider`.
3. Add a `request_id` or `trace_id` field to one `slog` line. Show the output.
4. Open `go doc runtime/metrics`. Write two metric names that relate to memory or GC.

#### Medium practical tasks

1. Add a counter and a histogram to a handler with a metrics library that your team allows. Scrape or print the values after ten `httptest` calls.
2. Create a parent span and a child span around a fake store call. Print the trace id from both.
3. Propagate a trace header through an HTTP client call to a second handler. Confirm that both spans share a trace id.

#### Advanced practical tasks

1. Export traces and metrics to a local OpenTelemetry Collector or a backend that you run. Capture a screenshot or a file of one trace. Document the ports.
2. Write a cardinality review of your labels. Remove or hash any high-cardinality label. Record the before and after series count if you can.

---

## Timeouts, retries, and circuit breakers

A timeout bounds how long an operation may run. A retry repeats a failed operation. A circuit breaker stops calls to a dependency that is already failing. Together they protect threads, goroutines, and user time.

Use `context.WithTimeout` or `WithDeadline` for outbound calls and for handler budgets. Set `http.Client.Timeout` or, better, timeouts on the transport and a context per request. Set server `ReadTimeout`, `WriteTimeout`, and `IdleTimeout` on `http.Server`.

Retries are safe when the operation is idempotent. A GET that reads data is often safe to retry. A POST that creates an order may create two orders if you retry. Use idempotency keys when the client may retry a write.

Backoff must grow. Wait a short time, then a longer time. Add jitter (random delay) so that many clients do not retry on the same tick. Cap the number of retries. Honor `Retry-After` when the server sends it.

A circuit breaker keeps three ideas: closed (calls pass), open (calls fail at once), half-open (a probe call tests the dependency). Open the circuit after a threshold of errors. Close it after probes succeed. Libraries such as `sony/gobreaker` implement this pattern. You can also implement a small version with an atomic state and a clock.

Do not retry on a 400-class error that will fail again. Do retry on a timeout or a 503 when the protocol says the server is overloaded. Log the attempt number.

Every retry multiplies load. A retry storm can take a service down. Circuit breakers and rate limits reduce that risk.

### Questions

#### Theoretical questions

1. Why does a timeout need a `context.Context` that you pass down the stack?
2. When is a retry unsafe?
3. What is jitter in backoff?
4. What are the closed, open, and half-open states?
5. How can retries make an outage worse?

#### Easy practical tasks

1. Wrap a `Sleep` that is longer than 50 ms in a 20 ms timeout context. Print the error.
2. Set `http.Server` read and write timeouts. Write the values and why you chose them.
3. Write a backoff sequence of four delays with jitter. Print the values.
4. List five HTTP status codes. Mark retry or do not retry for each code.

#### Medium practical tasks

1. Write a client helper that retries an idempotent GET three times with backoff. Test a fake transport that fails twice and then succeeds.
2. Refuse to retry a POST without an `Idempotency-Key` header. Test both paths.
3. Implement a tiny circuit breaker that opens after three errors and fails fast for one second. Test the fast-fail path.

#### Advanced practical tasks

1. Combine timeout, retry, and breaker on an outbound HTTP call. Add metrics for attempts, timeouts, and open-circuit rejects.
2. Review a handler budget: 200 ms total, 50 ms for auth, 100 ms for the store. Show how child timeouts derive from the parent deadline.

---

## Input validation and injection risks

Treat all input as untrusted. Input includes URL paths, query strings, headers, JSON bodies, cookies, and file names. Validate first. Then use the values.

Injection happens when you build a command, a query, or a path by joining strings with user data. The user data can change the meaning of the command.

SQL: use parameterized queries with `database/sql`. Do not join a user string into a SQL statement. Topic 17 covers this rule. The same rule applies to NoSQL query builders that accept raw strings.

Operating-system commands: do not pass user input to a shell. If you must run a process, use `os/exec` with a fixed executable path and an argument list. Do not use `bash -c` with a joined string.

File paths: do not open `filepath.Join(root, userPath)` without a check. A path that contains `..` can leave the root directory. Clean the path. Check that the result stays under the root.

HTML: if you generate HTML, encode user text. Use `html/template`. Do not write raw user strings into HTML.

Headers and logs: reject new lines in values that you write into headers or into log lines that other systems parse.

Return a clear 400 error for invalid input. Do not return stack traces or SQL text to the client. Log details on the server.

Validation libraries help. You can also write small checks: length limits, allowed character sets, and enumerated values. Prefer allow lists over block lists.

This section teaches defense. Do not practice injection against systems that you do not own.

### Questions

#### Theoretical questions

1. Why is a JSON body still untrusted after the JSON parses?
2. What is a parameterized SQL query?
3. Why is `bash -c` plus user text dangerous?
4. What path element lets a join leave a directory root?
5. Why must error responses omit stack traces in production?

#### Easy practical tasks

1. Write a validator that accepts an id of 1 to 32 letters or digits and rejects every other string. Add table tests.
2. Rewrite a fake `Query("... WHERE id = " + id)` into a parameterized form. Do not run it against a real database if you do not have one.
3. Use `html/template` to print a user name. Pass a name that contains `<`. Show the encoded output.
4. List six input sources on an HTTP request.

#### Medium practical tasks

1. Reject a file download name that leaves a `data/` root after `filepath.Clean`. Test `..` segments.
2. Start a subprocess with `os/exec` and a fixed argument list. Show that you do not call a shell.
3. Add `http.MaxBytesReader` and JSON field length checks. Test an oversize body.

#### Advanced practical tasks

1. Review one of your handlers for SQL, path, header, and log injection. Write a table of inputs, checks, and remaining risks.
2. Add a shared validation package for ids, emails, and page limits. Use it from two handlers. Prove a reject path with `httptest`.

---

## Secrets management

A secret is a value that grants access: a password, an API token, a private key, or a database URL that contains a password. Do not commit secrets to git. Do not put secrets in a Docker image layer. Do not log secrets.

Store secrets outside the source tree. Common places:

- environment variables that the host injects
- a cloud secret manager
- a file that the host mounts with strict permissions

Environment variables are easy and easy to leak in debug dumps. Secret managers add rotation, access logs, and short-lived credentials. Prefer a secret manager in production.

Use a `.env` file only on a laptop if your team allows it. Add `.env` to `.gitignore`. Commit a `.env.example` file with empty values or fake values.

Scan the repository with a secret scanner in CI. If a secret lands in git, rotate the secret. Assume that every clone and every fork has a copy. History rewrite does not reach every copy.

Separate secrets by environment. A development key must not open production data. Give each service its own credential. Rotate on a schedule and after a person leaves the team.

In Go, load secrets at startup into a struct. Do not print the struct with `%+v` if it contains secrets. Implement `String` or `slog.LogValuer` to redact.

TLS private keys and JWT signing keys are secrets. Store them in the same way as passwords.

### Questions

#### Theoretical questions

1. What must you do when a secret is committed to git?
2. Why is a `.env.example` file useful?
3. Why do secret managers beat environment variables for production?
4. Why must each service have its own credential?
5. How can `%+v` leak a secret?

#### Easy practical tasks

1. Create `.gitignore` with `.env`. Create `.env.example` with `DATABASE_URL=`.
2. Write a `Secret` type that prints `redacted` in `fmt` and `slog`.
3. List five places a token can leak (logs, traces, images, git, error pages).
4. Read your host documentation for one secret manager. Write the steps to fetch a secret at process start.

#### Medium practical tasks

1. Load a token from the environment. Fail startup when it is empty. Prove that `slog` does not print the token.
2. Add a CI secret scan (a tool that you can run locally). Commit a fake `AKIA` pattern on purpose in a branch, watch the scan fail, then remove it. Use only a fake value.
3. Split development and production config so that production requires a secret manager path. Document both.

#### Advanced practical tasks

1. Design rotation for a database password: two valid passwords during a window, then revoke the old one. Write the steps and the Go config change.
2. Audit a Dockerfile and a compose file for copied key files. Write a report of findings and fixes.

---

## Least privilege for OS and network

Least privilege means the process has only the rights that it needs. A stolen process then does less harm.

Operating system:

- Run the process as a non-root user in the container and on the host.
- Give the working directory and the listen port only the needed file permissions.
- Do not mount the Docker socket into an application container.
- Drop Linux capabilities that you do not need.
- Use read-only root file systems when the app does not write to disk.

Network:

- Listen on the interface that you need. An internal admin port must not bind `0.0.0.0` on a public host without a firewall.
- Allow egress only to the databases and APIs that the service must call.
- Use TLS for data in transit.
- Separate a public listener from a debug listener.

In Go, `http.Server` can serve two muxes on two addresses: public and internal. `net.Listen` can use a Unix socket for local-only access.

Cloud IAM roles must not use a wildcard when a single bucket or topic is enough. Database users must not be administrators. Grant `SELECT` and `INSERT` only on the tables that the service uses.

Review privilege when you add a feature. A new outbound URL is a new trust decision.

### Questions

#### Theoretical questions

1. What does least privilege mean for a process user?
2. Why is the Docker socket a high-privilege mount?
3. Why split public HTTP and pprof into two listeners?
4. Why must a database user not be an administrator?
5. What is egress control?

#### Easy practical tasks

1. Run a binary as a non-root user in a container (topic 18). Write the `USER` line.
2. Bind a debug server to `127.0.0.1:6060` and the public server to `:8080`. Document who can reach each port.
3. List the outbound hosts that one of your programs needs.
4. Write five file paths that the process must not write.

#### Medium practical tasks

1. Add a second `http.Server` for `/metrics` on an internal address. Prove that the public mux does not include `/metrics`.
2. Create a database role sketch (or a real role) with only the tables that the app uses.
3. Set a read-only root file system in a container and a writable `/tmp`. Confirm that a write to `/etc` fails.

#### Advanced practical tasks

1. Write a network policy (Kubernetes or firewall rules) that allows ingress 8080 and egress to one database port. Document deny-by-default.
2. Review a cloud role for the service. Remove one extra permission. Record the test that still passes.

---

## Feature flags and config hot-reload

A feature flag is a switch that changes behavior without a new deploy. Flags can be booleans or small enumerations. Use flags to hide an unfinished path or to roll out a change to a fraction of traffic.

Load flags from configuration or from a flag service. Pass a snapshot into the request path. Do not read a mutable global without a defined concurrency policy.

Hot-reload means the process reads a new configuration file or a new flag set while it runs. Hot-reload is useful. It is also a source of races and partial updates.

Safe hot-reload rules:

- Validate the new config first. Keep the old snapshot if validation fails.
- Swap a pointer to an immutable snapshot with `atomic.Pointer`.
- Do not change half of the fields while handlers run.
- Decide what must not reload (listen address, database URL). Restart for those values.

Dangerous hot-reload ideas:

- Reload TLS keys without a clear error path.
- Reload only some replicas so that the fleet disagrees.
- Parse a file on every request.

Feature flags need an owner and a removal date. An old flag is extra code. Test both values of every flag that you keep.

Do not use flags to hide security checks. A flag that disables authentication is a high-risk control. Protect it like a secret change.

### Questions

#### Theoretical questions

1. What problem does a feature flag solve?
2. Why must a hot-reload validate before it swaps?
3. Why is an immutable snapshot plus `atomic.Pointer` a useful pattern?
4. Which config values must require a restart?
5. Why must old flags be removed?

#### Easy practical tasks

1. Add a boolean `FEATURE_NEW_LIST` from the environment. Branch a handler on the value.
2. Write tests for both flag values.
3. List three settings that may reload and three that must restart.
4. Write a flag inventory table: name, owner, default, removal date.

#### Medium practical tasks

1. Watch a config file with `fsnotify` or a poll loop. Swap an `atomic.Pointer[Config]` after validation. Prove that a bad file does not replace a good snapshot.
2. Show a race: two fields updated without a snapshot. Then fix it with a pointer swap.
3. Serve a `/debug/flags` internal endpoint that prints non-secret flag values.

#### Advanced practical tasks

1. Design a 10 percent rollout of a new JSON field. Include metrics, a disable flag, and a rollback that does not need a rebuild.
2. Write a policy: who may change flags in production, how you audit the change, and how you test both paths in CI.

---

## Domain-driven design in Go (packages as boundaries)

Domain-driven design (DDD) is a way to structure software around the language of the business. A bounded context is a boundary where a word has one meaning. In Go, a package is the natural boundary.

Do not copy a Java folder tree (`entities`, `repositories`, `services`, `usecases`) into every module. That tree often creates circular imports and empty types. Name packages after the domain: `checkout`, `catalog`, `account`.

Put the core types and rules in a package that does not import `net/http` or a database driver. HTTP adapters and SQL adapters import the domain package. The domain package does not import the adapters.

Use the same words in code that the product team uses. If the team says "invoice", do not name the type `BillDoc`.

Keep aggregates small. An aggregate is a cluster of objects that you change as one. In Go, that is often one struct plus methods and a store interface.

Do not introduce a domain event bus in the first version. Start with functions and types. Add events when multiple contexts must react.

Package `internal/` keeps the domain private to the module when no other module must import it. Export only the types that adapters need.

Circular imports are a signal. If `checkout` and `catalog` import each other, split a small `money` or `sku` package, or move a function.

### Questions

#### Theoretical questions

1. What is a bounded context in one sentence?
2. Why is a Go package a better boundary than a `models` folder?
3. Why must a domain package avoid `net/http` imports?
4. What does ubiquitous language mean for type names?
5. What does a circular import often tell you about package design?

#### Easy practical tasks

1. Pick a small domain (library loans or shop cart). Write ten domain words and ten type or function names.
2. Draw two packages: `loan` and `httpapi`. Draw import arrows. The domain must not import HTTP.
3. List three names that leak HTTP or SQL into the domain (`UserRow`, `JSONUser`).
4. Read a short DDD glossary. Write definitions of entity and value object in your own words.

#### Medium practical tasks

1. Implement `cart` with `Add` and `Total` and no imports of `net/http`. Add `internal/httpapi` that decodes JSON and calls `cart`.
2. Split a circular import that you create on purpose between two domain packages. Record the split.
3. Write package comments that state the bounded context and the words that belong in the package.

#### Advanced practical tasks

1. Map a real product workflow to three packages. Write which types may cross a boundary and which types must not.
2. Compare a layered `service`/`repo` tree with a domain-package tree for the same app. Write six sentences on imports and test speed.

---

## Hexagonal and clean architecture in Go (keep it small)

Hexagonal architecture (ports and adapters) puts the application in the center. Ports are interfaces. Adapters talk to HTTP, databases, and queues. Clean architecture is a related idea with inner and outer rings.

In Go, a port is often a small interface:

```text
type OrderStore interface {
    Save(ctx context.Context, o Order) error
}
```

The HTTP handler and the SQL store both depend on the center. The center does not depend on the SQL driver.

Keep it small. A module with two interfaces and three packages is enough for many services. A module with twenty layers is not more professional. Extra layers hide the flow.

Rules that stay useful:

- `main` wires adapters to the application.
- Interfaces are small and defined next to the user when that is natural.
- Tests replace adapters with fakes.
- You can swap PostgreSQL for a fake store without changing the domain functions.

Rules that grow too large:

- A `UseCase` type for every function.
- A DTO for every field copy between identical structs.
- Interfaces with one implementation and one test fake that still live in a separate `ports` package too early.

Start with one `internal/app` package and one `internal/httpapi` package. Split when a file is hard to name or when imports cycle.

Topic 18 constructor injection is the usual wiring method. You do not need a framework to be hexagonal.

### Questions

#### Theoretical questions

1. What is a port in hexagonal architecture?
2. What is an adapter?
3. Why can too many layers hide the flow?
4. Where does `main` sit in this design?
5. When do you split `internal/app`?

#### Easy practical tasks

1. Draw a hexagon. Place domain, HTTP, and SQL on the diagram.
2. Write an `OrderStore` interface with two methods. Write a fake and a comment for a SQL adapter.
3. List four extra types that you will not add to a three-endpoint service.
4. Point to the composition root in a small module.

#### Medium practical tasks

1. Build `CreateOrder` in a package that imports only the standard library and your domain types. Call it from an HTTP adapter.
2. Swap the store from a map fake to another fake that fails. Keep the handler tests. Change only wiring in `main` or `NewApp`.
3. Count packages and interfaces in your module. Write a one-paragraph justification for each extra package.

#### Advanced practical tasks

1. Refactor a layered copy (`handlers`, `services`, `repos`) into fewer packages. Keep behavior. Show the import graph before and after.
2. Write a team guide of two pages: "Hexagonal Go, minimum size". Include examples of too little structure and too much structure.

---

## Event-driven systems and message queues

An event-driven system reacts to events. A message queue or a log carries those events. Common systems are Apache Kafka, NATS, and RabbitMQ.

Kafka is a durable distributed log. Consumers read by offset. Many teams use Kafka for high-volume event streams and for replay.

NATS is a lightweight messaging system. NATS Core is for pub/sub. JetStream adds persistence. Many teams use NATS for internal fan-out and for simpler operations than Kafka.

RabbitMQ is a broker with queues, exchanges, and routing keys. Many teams use RabbitMQ for work queues and for routing by key.

Delivery terms:

- At-most-once: a message may disappear. No duplicate.
- At-least-once: a message may arrive two times. You must make handlers idempotent.
- Exactly-once: hard across systems. Often you emulate it with idempotency keys and a store.

Producers must define a schema. Version the payload. Consumers must ignore unknown fields when you can. The same compatibility rules as HTTP APIs apply.

Consumers must handle poison messages (messages that always fail). Use a retry queue and a dead-letter queue. Log the payload id, not a secret.

Do not process a message in an unbounded goroutine storm. Use a worker limit (topic 19). Pass a context. Commit or ack only after the side effect succeeds, or use an outbox pattern when you must write a database and a message together.

Choose a broker with the operations team. The programming part is often the smaller part. Disk, lag, and consumer groups are the larger part.

### Questions

#### Theoretical questions

1. What is the difference between a durable log and a transient pub/sub channel?
2. Why does at-least-once delivery require idempotent consumers?
3. What is a dead-letter queue for?
4. Why must you version event payloads?
5. What problem does an outbox pattern address?

#### Easy practical tasks

1. Write a table: Kafka, NATS, RabbitMQ. Add rows for persistence, typical use, and one operations concern.
2. Design an `OrderPlaced` event JSON with a version field and an event id.
3. Write four sentences on at-least-once versus at-most-once for an email sender.
4. List three fields that must not appear in an event that you log (passwords, tokens, card numbers).

#### Medium practical tasks

1. Implement an in-process pub/sub with a bounded channel per subscriber. Drop or block on overflow. Document the choice.
2. Write a consumer loop that retries twice and then writes the message id to a dead-letter slice. Test a poison handler.
3. Make the handler idempotent with a processed-id set. Deliver the same event id twice. Assert one side effect.

#### Advanced practical tasks

1. Design an outbox table and a publisher worker for `OrderPlaced`. Write the failure cases (crash after DB write, crash after publish).
2. Compare Kafka consumer groups with RabbitMQ competing consumers for one work queue. Write a one-page operations and code impact note.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do API compatibility, event payload versioning, and feature flags work together during a long rollout?
2. Why do timeouts, least privilege, and secret redaction all limit the harm from a fault?
3. How do traces, metrics, and structured logs answer different questions about the same slow request?
4. When do DDD package boundaries and hexagonal ports lead to the same interfaces?
5. Why is an idempotent consumer required when you also retry HTTP calls and use a queue?

#### Easy practical tasks

1. Write a one-page production checklist: version, `/metrics` bind address, timeouts, validation, `.gitignore` for secrets, non-root user, flag owner, package names, event id.
2. Draw a request path: public listener, handler timeout, store call, optional queue publish. Mark where you validate and where you trace.
3. Redact a config struct and add a v1 JSON example file for CI replay.
4. List the ports that one service listens on and the hosts that it may call.

#### Medium practical tasks

1. Build a small module: v1 JSON API, `slog` plus a request duration metric, context timeout, input validation, constructor-injected store. Add `httptest` cases.
2. Add a feature flag and an in-process event that a second package consumes in a bounded worker. Show a poison-message path.
3. Run the service as non-root in a container with a secret from the environment. Confirm that logs omit the secret.

#### Advanced practical tasks

1. Write an architecture review of your module: compatibility policy, observability gaps, timeout map, threat notes (injection and secrets), package graph, and queue choice.
2. Add OpenTelemetry spans and a circuit breaker on one outbound call. Prove a fast-fail when the breaker is open and a full trace when the call works.
