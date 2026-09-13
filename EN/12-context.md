# 12. Context

## Description

`context.Context` carries a deadline, a cancel signal, and optional request values across API boundaries. You pass a context as the first parameter of a function. The callee stops work when the context is done.

This topic teaches why context exists and how you use the standard constructors. Complete Topic 11 before this topic. You must know goroutines, channels, and `select`.

Use one term for each concept. A context is not a bag of optional parameters. A cancel function is not optional. You must call cancel when you create a derived context.

---

## Why `context.Context` exists

A server handles many requests at the same time. Each request starts work: a database query, an HTTP call, a worker. The client can disconnect. A deadline can expire. The process can shut down. The work must stop. If the work does not stop, you leak goroutines, you hold connections, and you waste CPU.

`context.Context` is the standard type for that signal. The interface is in the `context` package:

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

`Done` closes when the work must stop. `Err` explains why. `Deadline` reports a time if a deadline exists. `Value` returns a request-scoped value.

Before `context`, packages invented their own cancel channels and timeouts. Those types did not compose. A function that calls two libraries had no common cancel type. The standard `Context` type solves that. From Go 1.7, `net/http` uses context on the request. The rest of the standard library followed.

A context forms a tree. You start from a root. You derive children with `WithCancel`, `WithTimeout`, `WithDeadline`, or `WithValue`. When a parent is done, the children are done. Cancel of a child does not cancel the parent.

Do not use context as a general configuration object. Do not store a context in a long-lived struct field unless the struct is the work that the context controls (a rare case). Pass the context down the call chain.

### Questions

#### Theoretical questions

1. What three kinds of data can a `context.Context` carry?
2. Why did Go add a standard context type instead of a cancel channel in every package?
3. What happens to child contexts when the parent is done?
4. Does cancel of a child cancel the parent?
5. Why must in-flight work stop when a client disconnects?

#### Easy practical tasks

1. Read `go doc context.Context`. Write the four methods and one sentence for each method.
2. Write five sentences that explain why context exists. Use only facts from this section.
3. Draw a tree: one root, two children, one grandchild. Mark the nodes that become done when the first child is cancelled.
4. Open [https://pkg.go.dev/context](https://pkg.go.dev/context). Write the package comment in your own words.

#### Medium practical tasks

1. Compare a custom `done chan struct{}` parameter with a `context.Context` parameter. Write six sentences. Cover composition, deadlines, and standard library support.
2. Find three functions in the standard library that take `context.Context` as the first parameter. Write their import paths and names.
3. Explain in a short note why a leaked HTTP handler goroutine is a context problem as well as a concurrency problem.

#### Advanced practical tasks

1. Read the Go blog post "Go Concurrent Patterns: Context" or the `context` package design notes. Write one page: problem, tree model, and one rule of use. Do not copy the text.
2. Inspect `net/http` documentation for `Request.Context`. Write how the server uses context for a cancelled request.

---

## `context.Background`, `TODO`

`context.Background` returns an empty root context. That context is never cancelled. It has no deadline. It has no values. Use `Background` at the top of a `main` function, in `init` when you must, and in tests when you start a tree.

```go
ctx := context.Background()
```

`context.TODO` also returns an empty context. Use `TODO` when you do not know which context to pass. `TODO` is a marker for later work. Replace `TODO` with a real context when the call chain is clear.

```go
ctx := context.TODO()
```

`Background` and `TODO` are not cancelled when a request ends. If you pass them into a server handler instead of `r.Context()`, the handler does not stop when the client disconnects.

Do not create a new `Background` in every helper. Derive from the context that the caller passed. Do not use `TODO` in a public API that already has a request or a parent context.

In tests, `Background` is normal for a unit test that does not need a deadline. Prefer `WithTimeout` in tests that can hang.

The two functions return a context that implements the `Context` interface. You do not close them. You do not call a cancel function on them. You derive children when you need cancel or a deadline.

### Questions

#### Theoretical questions

1. What is the difference between `Background` and `TODO`?
2. Does `Background` ever become done on its own?
3. When must you use `Background`?
4. When must you replace a `TODO` context?
5. Why is `Background` the wrong context inside an HTTP handler that already has a request?

#### Easy practical tasks

1. Write `package main` that prints `Background().Err()` and the `ok` value from `Deadline()`. Write the results.
2. Search a small module (or write three files) and replace one `TODO` with a parent context. Record the change.
3. Write four sentences: two valid uses of `Background`, two invalid uses.
4. Run `go doc context.Background` and `go doc context.TODO`. Write the intent of each function in your own words.

#### Medium practical tasks

1. Start a goroutine that waits on `ctx.Done()` with `ctx := context.Background()`. Explain why that wait never ends.
2. Write a public function that takes `ctx context.Context`. In `main`, pass `Background`. In a second call, pass `TODO`. Document which call you will change later.
3. Write a test helper that accepts a `context.Context`. Use `Background` in the test. Write why a timeout context can be better.

#### Advanced practical tasks

1. Audit a small public Go repository for `context.TODO()`. Classify each use as "valid placeholder" or "must pass a parent". Write the list.
2. Explain why `Background` and `TODO` are not the same type in documentation even when the current implementation looks similar. Use the package comments as evidence.

---

## `WithCancel`, `WithTimeout`, `WithDeadline`

These functions derive a child context from a parent. Each function returns a new context and a cancel function.

```go
ctx, cancel := context.WithCancel(parent)
defer cancel()

ctx, cancel := context.WithTimeout(parent, 2*time.Second)
defer cancel()

ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
defer cancel()
```

You must call `cancel`. The cancel function releases resources that the child holds (a timer, a node in the tree). `defer cancel()` is the normal form. A second call of `cancel` is safe. It does nothing.

`WithCancel` is done when you call `cancel` or when the parent is done.

`WithTimeout` is done when the duration ends, when you call `cancel`, or when the parent is done. `WithTimeout` is `WithDeadline` with `time.Now().Add(timeout)` as the deadline.

`WithDeadline` is done when the clock reaches the deadline, when you call `cancel`, or when the parent is done.

After the child is done, `ctx.Err()` is `context.Canceled` or `context.DeadlineExceeded`. If the parent is done first, the child reports the parent error.

From Go 1.21:

- `WithCancelCause` records a cause. `context.Cause(ctx)` returns that error.
- `WithTimeoutCause` and `WithDeadlineCause` exist for the same purpose.
- `WithoutCancel(parent)` returns a child that is not cancelled when the parent is cancelled. The child has no deadline and a nil `Done` channel. Values from the parent remain. Read `go doc context.WithoutCancel` before you use this function.
- `AfterFunc(ctx, f)` runs `f` in a new goroutine after `ctx` is done.

Do not ignore the cancel function. Do not wait only on a timeout and skip `cancel` after success. The timer can remain until the duration ends if you never cancel.

### Questions

#### Theoretical questions

1. Why must you call the cancel function that `WithTimeout` returns?
2. When is a `WithTimeout` context done?
3. What is the difference between `WithTimeout` and `WithDeadline`?
4. What values can `ctx.Err()` return after a derived context is done?
5. What extra information does `WithCancelCause` store?

#### Easy practical tasks

1. Create a `WithCancel` context. Start a goroutine that waits on `ctx.Done()`. Call `cancel` from `main`. Print `ctx.Err()`.
2. Create a `WithTimeout` context of 50 milliseconds. Wait on `ctx.Done()`. Print `ctx.Err()`. Confirm that the error is `DeadlineExceeded`.
3. Create a `WithDeadline` context with a time in the past. Check `ctx.Err()` at once. Write the result.
4. Write a function that takes a parent context and a duration, derives a timeout child, and always calls `cancel` with `defer`.

#### Medium practical tasks

1. Derive a timeout child from a parent. Cancel the parent first. Write which error the child reports.
2. Derive a timeout child. Finish the work before the timeout. Call `cancel`. Explain what resource you release.
3. Use `WithCancelCause` (Go 1.21 or later). Cancel with a custom error. Print `context.Cause(ctx)` and `ctx.Err()`.

#### Advanced practical tasks

1. Write a benchmark or a loop that creates many `WithTimeout` contexts and does not call `cancel`. Then call `cancel` in a second version. Compare memory or goroutine count.
2. Use `context.AfterFunc` to close a channel or to log when a context is done. Show that the function runs after cancel and after timeout.

---

## `WithValue` (and why to avoid it for optional params)

`WithValue` returns a child context that stores one key and one value.

```go
type traceIDKey struct{}

ctx := context.WithValue(parent, traceIDKey{}, "abc123")
id, ok := ctx.Value(traceIDKey{}).(string)
```

The lookup walks from the child to the parent. A child can shadow a key. `Value` returns `nil` when the key is absent.

Use `WithValue` only for request-scoped data that is not part of the function signature. Typical values are a trace identifier, a request identifier, or a user identity for logging and diagnostics.

Do not use `WithValue` for optional parameters. Do not put a database handle, a logger that the function must have, or a configuration struct in the context. Those values belong in parameters or in a struct that you pass in a clear way. Hidden parameters make the function hard to test and hard to read.

The key type must be comparable. Use an unexported struct type as the key. A string key can collide with another package that uses the same string.

```go
// good key
type ctxKey struct{}

// collision risk
ctx = context.WithValue(ctx, "user", u)
```

`WithValue` does not change cancel or deadline. The child is done when the parent is done.

Do not store a large object in every derived context if you can avoid it. Each `WithValue` adds a node to the tree.

### Questions

#### Theoretical questions

1. What is a valid use of `WithValue`?
2. Why must you not put optional function parameters in the context?
3. Why must the key be an unexported type in your package?
4. What does `Value` return when the key is missing?
5. Does `WithValue` change the deadline of the parent?

#### Easy practical tasks

1. Store a string request identifier with an unexported key. Read it in a helper function. Print the value.
2. Use a string key `"id"` in two packages in the same program (two files can simulate two packages). Show that the keys can collide. Then switch to unexported types.
3. Write four sentences: two valid context values, two invalid context values.
4. Read `go doc context.WithValue`. Write the warning from the documentation in your own words.

#### Medium practical tasks

1. Refactor a function that reads a logger from context into a function that takes the logger as a parameter. Keep a request identifier in the context. Write why the split is clearer.
2. Walk a chain of three `WithValue` nodes. Shadow one key. Print `Value` for each key at the leaf.
3. Write a test that fails when a handler requires a value that the test forgot to store. Then fix the test.

#### Advanced practical tasks

1. Design a small middleware that stores a request identifier in the context and a handler that logs it. Do not store the HTTP response writer in the context.
2. Measure or reason about a deep `WithValue` chain on a path that runs many times. Write when you would pack several diagnostic fields into one struct value instead of many `WithValue` calls.

---

## Propagating context through APIs

The first parameter of a function that can block, I/O, or start work is `ctx context.Context`. The name is `ctx`. The type is the interface, not a pointer.

```go
func FetchUser(ctx context.Context, id string) (User, error) {
    return db.QueryUser(ctx, id)
}
```

Pass the same `ctx` to every callee that does work for this request. Derive a child when the callee needs a shorter timeout or a local cancel.

```go
func FetchUser(ctx context.Context, id string) (User, error) {
    ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
    defer cancel()
    return db.QueryUser(ctx, id)
}
```

Do not store `ctx` in a struct that outlives the request. A typical server stores no context on the service struct. Each method receives `ctx` from the handler.

An exception is a worker that you start with a context and that exits when the context is done. The worker holds the context for its lifetime, not for the lifetime of the process.

Do not start a goroutine with a request context and then return from the handler while that goroutine still uses the context. The server cancels the request context when the handler returns. The goroutine sees `Canceled`. If the work must continue after the response, derive a new context from `Background` or from the server shutdown context, copy the values that you need, and set a new timeout. `context.WithoutCancel` (Go 1.21) can detach from parent cancel. You still must set a deadline and you still must avoid leaks.

Return `ctx.Err()` when you stop because the context is done. Do not hide that error without a reason. The caller needs to know that the work was cancelled.

Keep the context parameter even when the current body does not use it. The next I/O call will need it. A linter can require the first parameter to be a context.

### Questions

#### Theoretical questions

1. Where does the `context.Context` parameter sit in a function signature?
2. When do you derive a child instead of passing the parent unchanged?
3. Why must you not store a request context on a long-lived service struct?
4. What happens if a goroutine uses `r.Context()` after the HTTP handler returns?
5. What error do you return when you stop because the context is done?

#### Easy practical tasks

1. Write `Fetch` with `ctx` as the first parameter. Call `time.Sleep` in a loop and return `ctx.Err()` when `ctx.Done()` is ready.
2. From `main`, call `Fetch` with `WithTimeout` of 50 milliseconds and a sleep of 200 milliseconds. Print the error.
3. Draw a call chain of four functions. Mark where `ctx` is passed and where one function derives a shorter timeout.
4. Write five rules for context in APIs. Use imperative sentences.

#### Medium practical tasks

1. Write a service struct with a method `Get(ctx context.Context, id string)`. Do not put `ctx` in the struct. Call the method from a fake handler.
2. Start a goroutine from a function that receives `ctx`. Wait for the goroutine in that function. Pass `ctx` into the goroutine. Cancel from the test.
3. Show a bug: the handler starts a goroutine with `r.Context()` and returns. The goroutine prints `ctx.Err()`. Write the error that you see.

#### Advanced practical tasks

1. Implement "work continues after the response" in a safe way: copy the request identifier, use `Background` or `WithoutCancel`, set `WithTimeout`, and wait on shutdown. Document the lifetime.
2. Review three public APIs (standard library or a small module). Write whether they take `ctx` first. Note any exception and a reason.

---

## Listening for `ctx.Done()`

`ctx.Done()` returns a channel. The channel is closed when the context is done. A receive from `Done` succeeds at that moment. Do not send on `Done`. You only receive.

The usual wait is a `select`:

```go
select {
case <-ctx.Done():
    return ctx.Err()
case v := <-results:
    return handle(v)
}
```

Check `ctx.Err()` after `Done` is closed. `Err` is `nil` while the context is not done. After the context is done, `Err` is not `nil`.

A polling loop can check `Err` without a `select`:

```go
for {
    if err := ctx.Err(); err != nil {
        return err
    }
    if done := step(); done {
        return nil
    }
}
```

Polling without a wait burns CPU. Prefer `select` on `Done` and on the work channel.

`Done` can return `nil` for a context that is never cancelled (`Background`). A receive on a nil channel blocks forever. That is correct for `Background`: the work is not cancelled.

Do not close `Done` yourself. The `context` package owns that channel.

Combine `Done` with I/O that does not take a context. Start the I/O. Wait in `select` on `Done` and on a result channel. If `Done` wins, abandon or cancel the I/O in the way that the library allows. If the library has no cancel, you can leak a goroutine. Prefer APIs that accept a context.

From Go 1.21, `context.AfterFunc` can run cleanup when `Done` closes. You can still use `select` in the main flow.

### Questions

#### Theoretical questions

1. What does a receive from `ctx.Done()` mean?
2. When is `ctx.Err()` equal to `nil`?
3. What does `Done()` return for `context.Background()`, and what does a receive do?
4. Why is a tight loop that only checks `ctx.Err()` a poor wait?
5. Who closes the `Done` channel?

#### Easy practical tasks

1. Wait on `ctx.Done()` in `main` with a `WithCancel` context. Call `cancel` from a goroutine after 20 milliseconds. Print `ctx.Err()`.
2. Print `ctx.Err()` before and after timeout. Write both values.
3. Write a `select` with `ctx.Done()` and a `time.After` case. Use a long timeout on the context and a short `time.After`. Print which case ran.
4. Read `go doc context.Canceled` and `go doc context.DeadlineExceeded`. Write when each error appears.

#### Medium practical tasks

1. Write a worker that reads from a job channel and from `ctx.Done()` in one `select`. Cancel while the channel is empty. Confirm that the worker returns.
2. Wrap a `time.Sleep` with a context: use a timer and `select` on `ctx.Done()`. Return early on cancel.
3. Show that `Background().Done()` is nil. Write a comment that explains the block on receive.

#### Advanced practical tasks

1. Call an API that has no context parameter (for example a blocking send on a channel). Add a helper that abandons the wait on `ctx.Done()`. Document whether the inner goroutine can leak.
2. Use `errors.Is(err, context.Canceled)` and `errors.Is(err, context.DeadlineExceeded)` in a test. Cover both cases.

---

## HTTP and context cancellation

Incoming HTTP requests carry a context. `r.Context()` is done when the client disconnects, when the handler returns, or when the server hits a configured timeout.

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    row, err := store.Load(ctx, r.PathValue("id"))
    if err != nil {
        if errors.Is(err, context.Canceled) {
            return
        }
        http.Error(w, "unavailable", http.StatusServiceUnavailable)
        return
    }
    // write row
}
```

From Go 1.22, `ServeMux` patterns can include method and path wildcards. `r.PathValue` reads a wildcard. That change is routing. The context rules stay the same.

Outgoing HTTP calls must use a context:

```go
req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
if err != nil {
    return err
}
resp, err := client.Do(req)
```

Do not use `http.NewRequest` and then ignore cancel. `NewRequestWithContext` is the current API. `Request.WithContext` exists for a copy of a request.

`http.Client` also has `Timeout`. That timeout covers the whole call. The context deadline and the client timeout can both apply. Prefer one clear deadline. A context from the incoming request is the usual parent for outgoing calls. Derive a shorter timeout when the downstream must fail first.

Server timeouts:

- `ReadTimeout` and `WriteTimeout` on `http.Server`
- `http.TimeoutHandler` for a handler deadline
- `Server.Shutdown(ctx)` for graceful shutdown

`Shutdown` stops new connections. It waits for handlers. The `ctx` on `Shutdown` sets a limit for that wait. Handlers still see their request context.

Register a signal in `main`. Cancel a root context on `SIGINT` or `SIGTERM`. Pass that context to `Shutdown`. Topic 18 covers graceful shutdown in a full service.

Test cancel with `httptest` (Topic 13). Close the client connection or cancel the request context and assert that the handler returns.

### Questions

#### Theoretical questions

1. When does `r.Context()` become done on the server?
2. Which function creates an outgoing request that respects cancel?
3. What is the difference between `http.Client.Timeout` and a context deadline on the request?
4. What does `Server.Shutdown(ctx)` wait for?
5. Why must a handler pass `r.Context()` into database and HTTP calls?

#### Easy practical tasks

1. Write a handler that waits on `r.Context().Done()` and writes nothing after cancel. Start the server. Stop the client. Write how you observe the cancel.
2. Write a client call with `NewRequestWithContext` and a 50 millisecond timeout. Point the client at a handler that sleeps 500 milliseconds. Print the error.
3. Read `go doc net/http.Request.Context`. Write the cancel conditions in your own words.
4. List `ReadTimeout`, `WriteTimeout`, and `Shutdown` in a table with one purpose per row.

#### Medium practical tasks

1. Build a handler that calls a second `httptest` server with the incoming `r.Context()`. Cancel the outer request. Confirm that the inner client call stops.
2. Configure `http.TimeoutHandler` with a short limit. Hit a slow handler. Write the status code and the body.
3. Call `Server.Shutdown` with a timeout context from a goroutine after one request starts. Write whether the handler completes or the shutdown context ends first.

#### Advanced practical tasks

1. Implement `main` with signal cancel, an `http.Server`, and `Shutdown`. Prove that an in-flight handler can finish within the shutdown timeout.
2. Write an integration test that cancels the client context in the middle of a streamed response. Record the server-side `ctx.Err()` and the client error.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the full tree from `Background` to a timeout child that also stores a request identifier. Name every constructor and the cancel rule.
2. A teammate puts a `*sql.DB` in the context so that helpers stay short. Which facts do you use to reject that design?
3. How do `ctx.Done()`, `ctx.Err()`, and `errors.Is` work together in a library function?
4. Compare cancel on an incoming HTTP request with cancel on `Server.Shutdown`. What does each signal stop?
5. When do you use `TODO`, when do you use `Background`, and when do you derive `WithTimeout`?

#### Easy practical tasks

1. Write a small module: `Fetch(ctx, name)` sleeps, listens for `Done`, and returns `ctx.Err()`. Call it from `main` with `WithTimeout`. Print the error.
2. Make a one-page cheat sheet: `Background`, `TODO`, `WithCancel`, `WithTimeout`, `WithDeadline`, `WithValue`, `Done`, `Err`, `NewRequestWithContext`.
3. Draw an HTTP request: client disconnect, handler, database call, outgoing HTTP call. Label each context.
4. Run `go doc context` and list every `With*` function that your Go version prints. Mark the functions that this topic requires you to know.

#### Medium practical tasks

1. Write two tests for `Fetch`: one test cancels with `WithCancel`, one test uses `WithTimeout`. Assert the error with `errors.Is`.
2. Add a request identifier with `WithValue`. Log the identifier in `Fetch`. Do not pass the identifier as a second optional channel or map.
3. Write a handler plus `httptest` client. Cancel the client context. Assert that the handler observes `Canceled`.

#### Advanced practical tasks

1. Build a tiny API client: every method takes `ctx`, uses `NewRequestWithContext`, and maps `context` errors to a typed timeout error. Document the mapping. Test cancel and deadline.
2. Implement graceful shutdown: signal, root cancel, `Shutdown`, and one background worker that exits on the same root context. Show that no handler goroutine and no worker remains after shutdown.
