# 19. Advanced Concurrency and Systems

## Description

Topic 11 covers goroutines, channels, and mutexes. This topic covers the memory model, atomics, coordination packages, and tools that show the scheduler. Use these tools when a simple mutex or channel is not enough.

A data race is a bug. The memory model defines when one goroutine sees a write from another goroutine. Do not guess. Use the race detector and clear synchronization.

---

## Memory model and happens-before

The Go memory model is the document at [https://go.dev/ref/mem](https://go.dev/ref/mem). It defines happens-before. If event A happens before event B, then B sees the effects of A.

Within one goroutine, the program behaves as if statements run in order. The compiler and CPU may reorder work. That reorder must not change the result for a single goroutine.

Across goroutines, a write is visible only when a happens-before edge exists. A data race occurs when two goroutines access the same variable, at least one access is a write, and there is no happens-before relation. The result of a data race is undefined. Do not treat a race as a benign extra read.

The memory model lists synchronization that creates happens-before edges. Examples:

- The `go` statement happens before the start of the new goroutine.
- A send on a channel happens before the receive of that value completes.
- The close of a channel happens before a receive that returns the zero value because of that close.
- An unlock of a `sync.Mutex` happens before a later lock of the same mutex.
- `sync.Once` runs the function once. That run happens before every `Do` return.
- Package initialization happens before `main` starts.

`sync/atomic` operations on the same variable also synchronize. Use the atomic API. Do not mix plain reads with atomic writes on the same variable.

The race detector (`go test -race` or `go run -race`) reports many races. Use it in development and in CI. The detector does not prove the absence of all races. Tests must run the concurrent paths.

Do not use sleep to "fix" a race. Sleep hides a race. Add a happens-before edge: a channel, a mutex, `WaitGroup`, or an atomic.

### Questions

#### Theoretical questions

1. What does happens-before mean for a read in another goroutine?
2. When is a data race present?
3. Does a send on a channel create a happens-before edge? Explain.
4. Why is a sleep not a valid fix for a race?
5. What does package initialization guarantee relative to `main`?

#### Easy practical tasks

1. Open [https://go.dev/ref/mem](https://go.dev/ref/mem). Write five synchronization points from that page.
2. Write two goroutines that increment a shared `int` without a mutex. Run `go run -race`. Save the race report.
3. Fix the increment with a `sync.Mutex`. Run `-race` again. Confirm that the report is gone.
4. Draw a timeline of a channel send and receive. Mark the happens-before arrow.

#### Medium practical tasks

1. Show that closing a channel unblocks a receiver. Write how that close relates to happens-before.
2. Use `sync.Once` to initialize a value from two goroutines. Assert that setup runs one time and that both callers see the value.
3. Mix a plain read and an `atomic.Store` on the same `int64`. Run `-race`. Record the result. Then use only atomic operations.

#### Advanced practical tasks

1. Read the memory model section on unsynchronized reads. Write a one-page summary with your own diagrams. Do not invent extra guarantees.
2. Find a public Go issue or blog post about a surprising race. Explain the missing happens-before edge in six sentences.

---

## Atomics versus mutexes

`sync.Mutex` protects an invariant that may span several fields. You lock, you read or write the fields, and you unlock. The mutex defines a critical section.

`sync/atomic` updates a single word with well-defined atomic operations. Go 1.19 and later provide types such as `atomic.Int64`, `atomic.Uint64`, `atomic.Bool`, and `atomic.Pointer[T]`. Prefer those types over the older `atomic.AddInt64` function forms when you write new code.

Use an atomic for a counter, a flag, or a single pointer swap. Use a mutex when the update must keep two fields consistent. Example: a struct with `sum` and `count` must change both fields under one lock if readers need a consistent pair.

Atomics are not "mutexes that are faster in every case". A contended mutex can be slower than you want. An atomic hot variable can also be slow because of cache-line sharing. Measure before you replace a mutex.

Do not implement a lock with a spin on an atomic flag unless you have a proven need. A mutex parks the goroutine. A spin loop burns a CPU core and can block the scheduler.

`atomic.CompareAndSwap` (CAS) is the base of many lock-free algorithms. CAS fails when another goroutine changed the value. You retry or you fall back. Retry loops must stay short or must yield.

Read the documentation of each atomic method. A plain assignment to an `int64` is not atomic on all platforms for all alignments in the way you might expect. Use the `atomic` API for shared numeric flags and counters.

### Questions

#### Theoretical questions

1. When must you choose a mutex instead of an atomic?
2. What types did Go 1.19 add in `sync/atomic`?
3. Why is a spin lock on an atomic flag often a bad idea in Go?
4. What does compare-and-swap do when the current value does not match?
5. Why is "atomics are always faster" a false rule?

#### Easy practical tasks

1. Increment `atomic.Int64` from four goroutines. Print the final value. Use a `WaitGroup`.
2. Store and load a flag with `atomic.Bool`. Stop a loop when the flag is true.
3. Protect a `sum` and `count` pair with a mutex. Write one sentence that explains why one atomic is not enough.
4. Read `go doc sync/atomic.Int64`. List five methods.

#### Medium practical tasks

1. Replace a mutex around a single counter with `atomic.Int64`. Keep `-race` clean. Compare benchmark results with `go test -bench`.
2. Use `atomic.Pointer` to publish a read-only config struct. Writers swap a new pointer. Readers load the pointer. Do not mutate the pointed-to struct after publish.
3. Write a CAS loop that inserts a value into an `atomic.Int64` only when the old value is zero. Test success and failure.

#### Advanced practical tasks

1. Benchmark a mutex and an atomic on a contended counter. Change `GOMAXPROCS`. Write a short report with numbers.
2. Find a bug in a design that uses two atomics to update two fields. Demonstrate a torn read. Fix the design with a mutex or an immutable snapshot pointer.

---

## Lock-free patterns (when, and why not)

A lock-free algorithm guarantees that some thread makes progress even when other threads pause. Lock-free is a progress property. It is not a synonym for "uses atomics".

Correct lock-free code is hard. The ABA problem is a common defect. A value changes from A to B and back to A. A CAS can succeed even though the structure is not the same. Garbage collection removes some pointer ABA cases in Go. Logical ABA on counters and reuse of slots still exists.

Most Go programs must not start with a lock-free queue. A `sync.Mutex` plus a slice or a channel is clear. The race detector and code review work well with mutexes. Lock-free code often needs extra proofs and extra tests.

Reasons to avoid lock-free designs:

- The invariant is easy to get wrong.
- Debugging is harder than with a mutex.
- Performance gains are often small after measurement.
- The next reader of the code may introduce a race.

Reasons that can justify extra complexity:

- A profile shows a hot lock.
- The critical section is tiny and well specified.
- A standard type already solves the problem (`atomic.Pointer`, `sync.Map` in special cache cases).

`sync.Map` is not a default map replacement. It fits a narrow set of cases (many keys, write-once or disjoint keys). Read the package comment before you use it.

If you need a concurrent structure, prefer a mutex, a channel, or a well-tested library. Write a lock-free structure only when you can state the progress guarantee and the memory-model edges.

### Questions

#### Theoretical questions

1. What progress property does "lock-free" describe?
2. What is the ABA problem in one paragraph?
3. Why do mutexes fit most Go services?
4. When can `sync.Map` be the wrong choice?
5. Why is "uses atomics" not the same as "lock-free"?

#### Easy practical tasks

1. Write four sentences that distinguish lock-free, blocking, and wait-free at a high level. Use a reference page and cite it.
2. Implement a stack with a mutex and a slice. Write two tests: push and pop from several goroutines.
3. List three package comments (`sync.Map`, `atomic.Pointer`, `chan`) and the job each type already does.
4. Find the word "ABA" in a Go or Wikipedia article. Write the example in your own words.

#### Medium practical tasks

1. Profile a mutex-based counter and an atomic counter. State whether a lock-free design is justified.
2. Read the `sync.Map` documentation. Write a case that fits and a case that does not fit.
3. Review a lock-free snippet from a blog. Mark each atomic and each assumed happens-before edge. List open questions.

#### Advanced practical tasks

1. Implement a simple lock-free stack with `atomic.Pointer` for learning. Write tests with `-race`. Document ABA and why you must not use the code in production without a full review.
2. Compare your stack with a mutex stack under contention. Report throughput and the number of failed CAS retries.

---

## Semaphores and `errgroup`

A semaphore limits how many goroutines enter a section. `golang.org/x/sync/semaphore` provides a weighted semaphore. `NewWeighted(n)` allows a total weight of `n`. `Acquire(ctx, w)` takes weight `w`. `Release(w)` returns it. `TryAcquire` does not wait.

Use a semaphore to cap concurrent outbound calls or file opens. Pass a `context.Context` so that acquire can stop on cancel or timeout.

`golang.org/x/sync/errgroup` runs a set of goroutines and waits for them. `Group.Go` starts a function that returns an `error`. `Wait` returns the first error. `errgroup.WithContext` cancels the shared context when one function returns an error or when `Wait` returns.

`SetLimit(n)` (on current `x/sync`) limits the number of active goroutines in the group. Use `SetLimit` so that `Go` blocks when the group is full. That limit prevents an unbounded spawn.

```text
g, ctx := errgroup.WithContext(parent)
g.SetLimit(8)
g.Go(func() error { return fetch(ctx, url) })
err := g.Wait()
```

`errgroup` is a better default than a raw `WaitGroup` when any error must stop the rest of the work. A raw `WaitGroup` does not cancel siblings.

Do not ignore the error from `Wait`. Do not call `Go` after `Wait` on the same group. Create a new group for a new batch.

Add the module with `go get golang.org/x/sync`. Pin the version in `go.mod`.

### Questions

#### Theoretical questions

1. What does a weighted semaphore limit?
2. What does `errgroup.Wait` return?
3. What extra behavior does `errgroup.WithContext` add?
4. Why does `SetLimit` matter for safety?
5. When is a raw `WaitGroup` enough?

#### Easy practical tasks

1. Add `golang.org/x/sync` to a module. Show the `require` line in `go.mod`.
2. Use `semaphore.NewWeighted(2)` so that only two goroutines print at a time. Start five goroutines.
3. Use `errgroup.Group` to run three functions that return `nil`. Print the result of `Wait`.
4. Make one function return an error. Print the error from `Wait`.

#### Medium practical tasks

1. Fetch several URLs with `errgroup.WithContext` and `SetLimit(4)`. Cancel the others when one fetch fails.
2. Combine a semaphore of weight 3 with an `errgroup` that starts ten tasks. Each task acquires weight 1.
3. Write a test that fails if more than N tasks run at once. Use an atomic counter and a short sleep.

#### Advanced practical tasks

1. Replace an unbounded `for _, item := range items { go work(item) }` with `errgroup` and `SetLimit`. Keep `-race` clean. Compare peak goroutine count with `runtime.NumGoroutine`.
2. Document error policy: first error wins, remaining tasks see `ctx.Done()`, and how you log sibling errors.

---

## Rate limiting and backpressure

Rate limiting controls how often an operation starts. Backpressure slows the producer when the consumer cannot keep up. Both protect memory and downstream systems.

`golang.org/x/time/rate` implements a token bucket. `rate.NewLimiter(r, b)` allows rate `r` events per second with burst `b`. `Wait(ctx)` blocks until a token is available or the context ends. `Allow` is non-blocking. `Reserve` returns a reservation that you can cancel.

Use a limiter on outbound APIs and on expensive inbound handlers. Return HTTP status `429` when you reject a request. Include a `Retry-After` header when you have a clear retry time.

Backpressure in Go is often a blocking channel send. A buffer of size N holds N items. A full channel blocks the sender. That block is useful. An unbounded slice that grows in a goroutine is not backpressure. It is a memory leak risk.

Patterns:

- A worker pool with a bounded job channel.
- A semaphore on in-flight requests.
- `rate.Limiter` before you start heavy work.
- HTTP server timeouts so that stalled clients do not hold slots.

Do not drop work in silence unless the protocol allows loss. Metrics must show drops, waits, and queue depth.

Combine rate limits with context deadlines. A wait that exceeds the deadline must fail. The caller must handle that error.

### Questions

#### Theoretical questions

1. What is the difference between rate limiting and backpressure?
2. What do the two arguments of `rate.NewLimiter` mean?
3. How does a full buffered channel apply backpressure?
4. Which HTTP status code reports a rate-limit reject?
5. Why is an unbounded in-memory queue a risk?

#### Easy practical tasks

1. Add `golang.org/x/time/rate`. Create a limiter of 2 events per second with burst 2. Call `Wait` five times. Print times.
2. Use `Allow` in a loop. Count how many calls succeed in one second.
3. Send 10 values into a channel of capacity 3 from one goroutine. Receive slowly. Describe the pause.
4. Write four sentences on token bucket burst versus sustained rate.

#### Medium practical tasks

1. Wrap an HTTP handler with a limiter. Return `429` when `Allow` is false. Test with `httptest`.
2. Build a worker pool with a job channel of capacity 8 and 4 workers. Measure how a fast producer blocks.
3. Cancel a `Wait` with a timeout context. Assert that you get a context error.

#### Advanced practical tasks

1. Add queue-depth and reject-count logs or metrics to a handler. Generate load. Show how burst and rate change the numbers.
2. Design backpressure across two stages (HTTP accept and outbound HTTP). Write the limits, the errors, and what the client sees.

---

## Scheduler traces and `GODEBUG`

The Go scheduler runs goroutines on operating-system threads. When a program is slow or stuck, you need traces, not guesses.

`GODEBUG` is an environment variable that enables runtime debug output. Examples:

- `GODEBUG=gctrace=1` prints garbage collector traces.
- `GODEBUG=schedtrace=1000` prints scheduler stats every 1000 milliseconds.
- `GODEBUG=schedtrace=1000,scheddetail=1` adds per-P detail.

Go 1.21 and later can also record `GODEBUG` defaults in `go.mod`. See `go help godebug`.

`runtime/trace` records an execution trace. Start a trace, run the workload, stop the trace. `go test -trace=trace.out` writes a test trace. Inspect the file with `go tool trace trace.out`. The viewer shows goroutine timelines, syscalls, and GC.

`pprof` goroutine and CPU profiles are related tools. A goroutine profile shows stacks that are blocked. A mutex profile shows contended locks. Use `net/http/pprof` on a debug port that is not public.

`GOTRACEBACK` controls how much stack output a crash prints. `GOTRACEBACK=all` prints every goroutine. Use that value when you debug a deadlock.

Do not leave verbose `GODEBUG` on in production without a plan. The output is large. Use it for a short debug window.

Read traces with a question: which goroutine waits, which lock is held, and is `GOMAXPROCS` the limit that you expect?

### Questions

#### Theoretical questions

1. What does `GODEBUG=schedtrace=1000` print?
2. What tool opens a `trace.out` file?
3. Why must a `pprof` HTTP endpoint stay off the public internet?
4. What does `GOTRACEBACK=all` add to a crash?
5. Where can a module set `GODEBUG` defaults in modern Go?

#### Easy practical tasks

1. Run a small program with `GODEBUG=gctrace=1`. Save ten lines of output. Label GC start and end if you can.
2. Run the same program with `GODEBUG=schedtrace=500` for two seconds. Write what `gomaxprocs` shows.
3. Run `go help godebug`. Write the purpose of `gctrace` and `schedtrace`.
4. Add `import _ "net/http/pprof"` to a debug build only. List two `pprof` paths.

#### Medium practical tasks

1. Write a test that starts `runtime/trace`, starts many goroutines, and stops the trace. Open `go tool trace` and describe one stall.
2. Capture a goroutine profile from a program that deadlocks on purpose (two mutexes). Use `GOTRACEBACK=all` or `pprof`. Identify the two stacks.
3. Compare a CPU profile and an execution trace for the same benchmark. Write three facts that only the trace shows.

#### Advanced practical tasks

1. Use `schedtrace` while you change `GOMAXPROCS` from 1 to 4. Record runnable queue hints. Explain one line of output.
2. Document a debug playbook: when to use `-race`, `pprof`, `go tool trace`, and `GODEBUG`. Give one example symptom for each tool.

---

## `singleflight` and `semaphore`

`golang.org/x/sync/singleflight` collapses duplicate work for the same key. `Group.Do(key, fn)` runs `fn` once for concurrent callers with that key. All callers wait. All callers receive the same result and error.

`DoChan` returns a channel. `Forget` removes a key so that the next `Do` starts new work.

A typical use is a cache fill. Many incoming requests miss the same cache key. Without `singleflight`, each request hits the database. With `singleflight`, one request fills the cache. The others wait and reuse the value. This pattern reduces a burst of identical work.

`singleflight` is not a cache. It does not store the value after the waiters finish. You still need a cache if later requests must reuse the value.

`semaphore` (the same `x/sync/semaphore` package as in the earlier section) limits concurrency by weight. Use it when keys are different but the resource is shared (connections, memory, file descriptors).

Use `singleflight` when the key is the same and the work is identical. Use a semaphore when the work is different and only the count matters. You can use both: collapse the same key, and cap total in-flight fills.

Handle errors. If `fn` fails, all waiters see the error. Decide whether a failed fill must be retried at once. `Forget` can allow an immediate retry.

Do not call `Do` while you hold a lock that `fn` also needs. That order can deadlock.

### Questions

#### Theoretical questions

1. What problem does `singleflight.Group.Do` solve?
2. Why is `singleflight` not a cache by itself?
3. What does `Forget` change for the next caller?
4. When do you pick a semaphore instead of `singleflight`?
5. How can `Do` deadlock with a mutex?

#### Easy practical tasks

1. Call `Do` with the same key from five goroutines. Let `fn` sleep and increment a counter. Confirm that the counter is 1.
2. Call `Do` with two keys. Confirm that `fn` runs two times.
3. Make `fn` return an error. Confirm that every waiter sees the error.
4. Write five sentences that describe a cache miss storm and how `singleflight` reduces it.

#### Medium practical tasks

1. Put `singleflight` in front of a fake slow store. Measure store calls with and without the group under parallel `httptest` requests for one id.
2. After a failed `Do`, call `Forget` and retry. Write a test for that path.
3. Combine `singleflight` for per-key collapse with a semaphore of 4 for total fills. Document the two limits.

#### Advanced practical tasks

1. Implement a read-through cache: map + mutex + `singleflight`. Write tests for hit, miss, concurrent miss, and error.
2. Compare `Do` and `DoChan` in a handler that must also listen to `ctx.Done()`. Write the cancellation behavior that you observe.

---

## Avoiding unbounded goroutines

A goroutine is cheap compared to an operating-system thread. A goroutine is not free. Each goroutine has a stack and scheduler state. Unbounded `go` statements can exhaust memory or create millions of blocked stacks.

Dangerous patterns:

- `for _, item := range items { go work(item) }` when `items` has no cap.
- `go` per HTTP request that then starts more `go` without a limit.
- A receiver that starts a goroutine for every channel value while the producer is faster than the consumer.
- Forgotten `go` in a package-level function that callers invoke in a loop.

Safer patterns:

- A fixed worker pool and a bounded job channel.
- `errgroup` with `SetLimit`.
- A weighted semaphore around `go` or around the work function.
- `context` cancel so that workers exit.

Every goroutine that you start must have a stop condition. Document who cancels the context. Document who calls `Wait`.

HTTP servers already run a goroutine per connection. Do not add another unbounded fan-out inside the handler. If you need background work after the response, use a bounded queue and reject work when the queue is full.

`runtime.NumGoroutine` helps in tests and in metrics. A leak test can start a function and check that the count returns to a baseline after cancel.

The race detector does not report a leak. A goroutine profile does. Check blocked stacks that point at your code.

### Questions

#### Theoretical questions

1. Why is a goroutine not free?
2. Why is `go` in a range loop over user input dangerous?
3. What is a stop condition for a worker goroutine?
4. Why does an HTTP server already limit how you should spawn extra goroutines?
5. Which profile shows leaked goroutines?

#### Easy practical tasks

1. Start 10000 goroutines that wait on a channel that nobody receives from. Watch memory or `NumGoroutine`. Then stop the program.
2. Rewrite that loop with a worker pool of 8. Confirm that `NumGoroutine` stays near 8 plus the main goroutines.
3. Add a context cancel and a `WaitGroup`. Confirm that workers exit.
4. Write a checklist of four questions to ask before you type `go`.

#### Medium practical tasks

1. Write an HTTP handler that enqueues work on a bounded channel. Return status `503` when the queue is full. Test both paths.
2. Use `errgroup.SetLimit` for a batch API that accepts a list of ids. Cap the list length and the group limit.
3. Write a leak test: run the function, cancel, wait, and compare `NumGoroutine` with a small delta.

#### Advanced practical tasks

1. Find an unbounded `go` in a practice service (or introduce one). Capture a goroutine profile. Fix the bound. Show before and after counts under load.
2. Design a background job system with queue size, worker count, enqueue timeout, and shutdown drain. Write the failure modes.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do happens-before, mutexes, and atomics work together to make a counter and a multi-field struct safe?
2. When do you use `errgroup`, when do you use `singleflight`, and when do you use a semaphore?
3. How do rate limits and bounded channels both apply backpressure, and how do they differ?
4. What evidence do `GODEBUG=schedtrace`, `go tool trace`, and a goroutine profile each give you?
5. Why is a lock-free structure the last choice after a mutex, a channel, and a limit on goroutines?

#### Easy practical tasks

1. Write a one-page map: problem to tool (`-race`, mutex, `atomic.Int64`, `errgroup`, `rate.Limiter`, `singleflight`, `SetLimit`).
2. Create a module that fetches three URLs with `errgroup` and a rate limiter. Print the first error if any.
3. Run one program with `-race` and one with `GODEBUG=gctrace=1`. Save both outputs.
4. Draw a diagram of a worker pool with a bounded queue and a semaphore on outbound HTTP.

#### Medium practical tasks

1. Build a handler that uses `singleflight` per id, a semaphore of 8, and `429` when a limiter rejects. Add `httptest` cases.
2. Introduce a data race, detect it, then fix it with the smallest synchronization that the memory model allows. Write why that edge is enough.
3. Record an execution trace of the handler under parallel load. Write three observations from `go tool trace`.

#### Advanced practical tasks

1. Load-test a small service that used unbounded `go`. Add limits and `errgroup`. Report p95 latency, reject count, and peak `NumGoroutine` before and after.
2. Write an incident playbook: race in CI, deadlock in staging, latency spike, memory growth from goroutines. Name the first command for each case.
