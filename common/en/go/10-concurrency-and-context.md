# 10. Concurrency and Context

## Description

Go runs many functions at the same time. A goroutine is a function that the Go runtime schedules. A channel moves values between goroutines. The `sync` package protects shared data when a channel is not the simple choice. `context.Context` carries a cancel signal, a timeout, and a deadline across API boundaries.

Complete this topic before you build servers, workers, or other concurrent programs. Use one term for each concept. A goroutine is not an operating-system thread. A channel send is not a function call. A data race is not the same as a deadlock. A context is not a bag of optional parameters.

---

## Goroutines

The `go` keyword starts a goroutine. The new goroutine runs `f` while the current goroutine continues.

```go
go f()
go f(a, b)
go func() {
	doWork()
}()
```

Arguments to `f` are evaluated in the current goroutine. The call of `f` runs in the new goroutine.

The process starts with one goroutine. That goroutine runs `main`. When `main` returns, the process exits. Other goroutines stop at that moment. You must wait if those goroutines must finish. Use a `sync.WaitGroup` or a channel for that wait.

A goroutine has a small stack at start. The stack grows and shrinks as the function needs memory. You can start many goroutines. You must still bound the number when each goroutine holds a resource (a file, a connection, or a large buffer).

From Go 1.22, each iteration of a `for` loop has its own loop variables. This code is safe in Go 1.22 and later:

```go
for _, item := range items {
	go func() {
		process(item)
	}()
}
```

Before Go 1.22, every goroutine could read the last `item`. If you support older compilers, copy the variable inside the loop: `item := item`.

The Go scheduler runs many goroutines on a smaller number of operating-system threads. This model is an M:N scheduler. `GOMAXPROCS` is the maximum number of logical processors that can run Go code at the same time. The default is the number of logical CPUs.

Do not start a goroutine and ignore its error. Return the error on a channel, or collect it with a `WaitGroup` and a mutex.

### Questions

#### Theoretical questions

1. What happens to other goroutines when `main` returns?
2. Where does the compiler evaluate the arguments of `go f(a, b)`?
3. Why is a goroutine not the same as an operating-system thread?
4. What changed for loop variables in Go 1.22, and why does that matter for `go func()` inside a loop?
5. What is `GOMAXPROCS`?

#### Easy practical tasks

1. Write `package main` that starts one goroutine. The goroutine prints `hello`. The `main` function sleeps for 200 milliseconds, then returns. Run the program.
2. Remove the sleep. Run the program ten times. Write what you observe.
3. Start three goroutines that each print a different number. Use a `time.Sleep` in `main` so that the prints appear.
4. Print `runtime.GOMAXPROCS(0)` and `runtime.NumCPU()`. Write both numbers.

#### Medium practical tasks

1. Start a goroutine that sends one `int` on a channel. Receive that value in `main` and print it. Do not use `Sleep`.
2. Write a loop that starts one goroutine per slice element. Each goroutine prints the element. Confirm the program is safe on Go 1.22.
3. Start a goroutine that panics. Observe the process. Write how a panic in one goroutine affects the process.

#### Advanced practical tasks

1. Measure how many goroutines you can start before the program is slow or the machine is short of memory. Record the number and the limit that you hit. Stop the test before the machine becomes unusable.
2. Write a small program that starts work in a goroutine and must observe an error from that work. Use a channel of type `error`. Document the protocol.

---

## Channels, close, and `select`

A channel is a typed conduit. You create a channel with `make`.

```go
ch := make(chan int)    // unbuffered
ch := make(chan int, 4) // buffered, capacity 4
```

Send a value with `ch <- v`. Receive a value with `v := <-ch`. The comma-ok form reports whether the channel is open:

```go
v, ok := <-ch
```

`ok` is `false` when `ch` is closed and empty. `v` is the zero value of the element type in that case.

`close(ch)` marks the channel as closed. Receivers drain remaining buffered values, then they see the closed state.

Channel rules:

- A send to a closed channel panics.
- A close of a closed channel panics.
- A close of a nil channel panics.
- A send on a nil channel blocks forever.
- A receive on a nil channel blocks forever.

The sender closes the channel. The receiver does not close the channel when more than one sender exists.

An unbuffered channel has capacity 0. A send waits until another goroutine receives. A buffered channel stores values while the buffer is not full. Do not use `len(ch)` to decide that a send is safe.

A `select` statement waits for a set of channel operations. `select` runs one ready case. If more than one case is ready, `select` chooses one at random.

```go
select {
case v := <-in:
	handle(v)
case out <- v:
	// sent
case <-time.After(2 * time.Second):
	return errTimeout
}
```

A `default` case runs when no other case is ready. Use `default` for a non-blocking attempt.

`for v := range ch` receives until the channel is closed and the buffer is empty. If nobody closes the channel, the loop waits forever.

### Questions

#### Theoretical questions

1. What does the comma-ok receive form tell you?
2. Who must close a channel, and why?
3. What happens when you send on a closed channel?
4. What does `select` do when two cases are ready?
5. When does `for v := range ch` stop?

#### Easy practical tasks

1. Create an unbuffered `chan string`. Send `"ok"` from a goroutine. Receive in `main` and print the value.
2. Close a channel in the sender. Use comma-ok in the receiver. Print `v` and `ok` for two receives after the close.
3. Write a `select` with two receive cases from two channels. Send one value to one channel from a goroutine. Print which case ran.
4. Receive from a channel with a 100 millisecond timeout. Do not send. Print that the timeout case ran.

#### Medium practical tasks

1. Write a function `func produce(n int) <-chan int` that sends `0` to `n-1` and then closes the channel. Range over the result in `main`.
2. Create `make(chan int, 2)`. Send two values from `main` without a receiver. Print `len` and `cap`. Then receive both values.
3. Implement a non-blocking send helper: return `true` if the send occurred, `false` if the channel was not ready.

#### Advanced practical tasks

1. Design a protocol for one channel that carries both results and a terminal error. Document the message type and the close rule. Implement a small demo.
2. Build a `select` that waits on work, a ticker, and a done channel. Stop the ticker when done. Show that the goroutine returns.

---

## `sync.WaitGroup` and `sync.Mutex`

`sync.WaitGroup` waits for a set of goroutines. Call `Add` before you start each goroutine. Call `Done` when the goroutine finishes. Call `Wait` in the waiter.

```go
var wg sync.WaitGroup
for i := range 3 {
	wg.Add(1)
	go func(n int) {
		defer wg.Done()
		work(n)
	}(i)
}
wg.Wait()
```

`Add` must not run inside the new goroutine when `Wait` can run first. That order is a race. Call `Add` in the starter. Prefer `defer wg.Done()` so that `Done` runs on every path.

Do not copy a `WaitGroup` after you use it. Pass a pointer.

`sync.Mutex` is a mutual exclusion lock. Call `Lock` before you touch shared data. Call `Unlock` after. Prefer `defer mu.Unlock()` after a successful lock.

```go
var mu sync.Mutex
var count int

func inc() {
	mu.Lock()
	defer mu.Unlock()
	count++
}
```

A mutex does not protect data by itself. You must lock every read and every write of that data. Document which fields the mutex protects.

`sync.RWMutex` allows many readers or one writer. Use it when reads are common and writes are rare. Do not hold an `RLock` and then take `Lock` in the same goroutine. That order deadlocks.

`sync.Once` runs a function one time. `sync.Pool` reuses temporary objects. Start with `Mutex` and `WaitGroup`. Add the other types when you measure a need.

Maps are not safe for concurrent write. Protect a shared map with a mutex. Concurrent read with no writes is safe.

### Questions

#### Theoretical questions

1. When must you call `WaitGroup.Add` relative to `go` and `Wait`?
2. Why do you use `defer wg.Done()`?
3. What does a mutex protect?
4. Why must you not copy a `WaitGroup` or a `Mutex` after first use?
5. When do you choose `RWMutex` instead of `Mutex`?

#### Easy practical tasks

1. Start three goroutines that each print a number. Wait with a `WaitGroup`. Do not use `Sleep`.
2. Protect an `int` counter with a `Mutex`. Increment from four goroutines. Print the final value.
3. Forget `Done` on purpose. Observe that `Wait` does not return. Stop the program. Then add `Done`.
4. Write five sentences that state when you use a channel and when you use a mutex.

#### Medium practical tasks

1. Collect errors from several goroutines. Use a mutex and a `[]error`, or use a channel. Wait with a `WaitGroup`.
2. Share a map under a mutex. Insert keys from several goroutines. Read after `Wait`.
3. Show a deadlock: lock a mutex twice in the same goroutine. Then fix the code.

#### Advanced practical tasks

1. Implement a worker pool: a fixed number of goroutines, a job channel, and a `WaitGroup`. Close the job channel to stop.
2. Compare a mutex around a map with one owner goroutine that receives updates on a channel. Write six sentences on complexity and races.

---

## Race detector and goroutine leaks

A data race occurs when two goroutines access the same memory and at least one access is a write, and the accesses are not synchronized.

Build and test with the race detector:

```text
go test -race ./...
go run -race .
```

The detector prints a report when it sees a race. The report shows the two stacks. Fix the race with a mutex, with a channel, or with no shared write. Do not use `time.Sleep` to hide a race.

The race detector uses extra memory and CPU. Use it in tests and in staging. It does not catch every race. A passing `-race` run is not a proof of safety.

A goroutine leak is a goroutine that never returns. Common causes:

- A receive on a channel that nobody closes and nobody sends
- A send on a channel that nobody receives
- A `Wait` that nobody `Done`s
- A `select` that never becomes ready
- Work that ignores context cancel

Leaked goroutines hold stacks and any values they capture. In a server they grow until the process is slow or dies.

Find leaks with `runtime.NumGoroutine` in tests, with `pprof` goroutine profiles, and with careful `select` on `ctx.Done()`. Topic 11 covers `pprof`.

Every goroutine that you start must have a stop rule. Document who closes the channel. Document who cancels the context.

### Questions

#### Theoretical questions

1. What is a data race?
2. What command enables the race detector for tests?
3. Does a passing `-race` run prove that the program has no races?
4. What is a goroutine leak?
5. Name two causes of a leaked goroutine.

#### Easy practical tasks

1. Write a program that increments a shared `int` from two goroutines without a lock. Run `go run -race .`. Record the report.
2. Fix the program with a mutex. Run `-race` again. Confirm that the report is gone.
3. Start a goroutine that receives from a channel that nobody closes. Print `NumGoroutine` after one second. Write why the count stays high.
4. Read `go help test` for `-race`. Write the purpose in one sentence.

#### Medium practical tasks

1. Create a race on a map write from two goroutines. Run `-race`. Then protect the map.
2. Write a function that leaks a goroutine with `time.Tick` in a loop that you call many times. Then fix it with `NewTicker` and `Stop`.
3. Add a test that fails if `NumGoroutine` grows after you start and stop a worker.

#### Advanced practical tasks

1. Capture a goroutine profile with `pprof` after you leak one goroutine on purpose. Identify the stack. Then fix the leak and capture again.
2. Review a small concurrent program. List every goroutine, its stop rule, and whether `-race` covers its shared data.

---

## `context` cancel, timeout, and deadline

`context.Context` carries a deadline, a cancel signal, and optional request values. You pass a context as the first parameter of a function. The callee stops work when the context is done.

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}
```

`Done` closes when the work must stop. `Err` explains why. `Deadline` reports a time if a deadline exists.

`context.Background` returns an empty root context. That context is never cancelled. Use `Background` at the top of `main` and in tests when you start a tree. `context.TODO` also returns an empty context. Use `TODO` only as a marker. Replace it when the call chain is clear.

A context forms a tree. You derive children with `WithCancel`, `WithTimeout`, or `WithDeadline`. When a parent is done, the children are done. Cancel of a child does not cancel the parent.

```go
ctx, cancel := context.WithCancel(parent)
defer cancel()

ctx, cancel := context.WithTimeout(parent, 2*time.Second)
defer cancel()

ctx, cancel := context.WithDeadline(parent, time.Now().Add(2*time.Second))
defer cancel()
```

You must call `cancel` when you create a derived context. `cancel` releases resources. `defer cancel()` is the common form. After timeout, `ctx.Err()` is `context.DeadlineExceeded`. After cancel, `ctx.Err()` is `context.Canceled`.

Check `ctx.Done()` in `select`. Return `ctx.Err()` to the caller. Pass the same context down to `http.NewRequestWithContext` and to database calls.

Do not store a context in a long-lived struct field unless the struct is the work that the context controls. Pass the context down the call chain.

### Questions

#### Theoretical questions

1. What three kinds of data can a `context.Context` carry?
2. What happens to child contexts when the parent is done?
3. Does cancel of a child cancel the parent?
4. Why must you call the `cancel` function from `WithCancel` or `WithTimeout`?
5. What is the difference between `WithTimeout` and `WithDeadline`?

#### Easy practical tasks

1. Read `go doc context.Context`. Write the four methods and one sentence for each method.
2. Create `WithCancel`. Start a goroutine that waits on `ctx.Done()`. Call `cancel`. Print `ctx.Err()`.
3. Create `WithTimeout` of 50 milliseconds. Sleep 100 milliseconds. Print `ctx.Err()`.
4. Draw a tree: one root, two children, one grandchild. Mark the nodes that become done when the first child is cancelled.

#### Medium practical tasks

1. Write a worker that loops with `select` on a job channel and `ctx.Done()`. Cancel from `main`. Show that the worker returns.
2. Compare `WithTimeout` with `select` and `time.After`. Write four sentences on why context composes better.
3. Pass a timed context into `http.NewRequestWithContext`. Use a short timeout. Record the error.

#### Advanced practical tasks

1. Derive a timeout child from a cancelable parent. Cancel the parent first in one run. Let the timeout fire first in a second run. Print `Err` in both runs.
2. Inspect `net/http` documentation for `Request.Context`. Write how the server uses context for a cancelled request.

---

## Why to avoid `WithValue` for optional parameters

`context.WithValue` stores a key and a value in a derived context.

```go
type key struct{}

ctx := context.WithValue(parent, key{}, userID)
id, ok := ctx.Value(key{}).(string)
```

Use a private key type so that other packages cannot collide. The value API is untyped. The caller must assert.

`WithValue` is for request-scoped data that must cross an API that you do not control. Typical cases: a trace id, a request id, and sometimes a user principal that middleware already stored.

Do not use `WithValue` for optional parameters of your own functions. Do not hide a logger, a database handle, or a config struct in context. Those values belong in a struct that you pass as a normal parameter, or in a field of a server type.

Problems with values in context:

- Callers cannot see the required inputs in the function signature.
- Tests must build a context instead of a clear argument list.
- A missing key is a silent zero value after a failed assertion.
- The same context can carry hidden dependencies that create import cycles later.

Pass context for cancel and deadline. Pass data as parameters. If a third-party function only accepts `context.Context`, store the smallest value that that library documents.

Do not store a context inside a value that you then put back into the same context tree. That cycle is hard to cancel.

### Questions

#### Theoretical questions

1. What is a valid use of `WithValue`?
2. Why must the context key be an unexported type in your package?
3. Why is a logger in context a poor design for your own API?
4. How does `WithValue` hide dependencies from the function signature?
5. What do you pass instead of optional parameters in context?

#### Easy practical tasks

1. Store a request id with `WithValue` and a private key. Read it back with `Value` and a type assertion.
2. Read a missing key. Print `ok` from the assertion.
3. Write four sentences that restate the rule: context for cancel, parameters for data.
4. List three values that must not go in context in your own service.

#### Medium practical tasks

1. Write two versions of `func Query`: one takes `ctx` and `userID string`, one takes only `ctx` and reads `userID` from `Value`. Compare call sites in six sentences.
2. Show a key collision: two packages use the string key `"id"`. Then fix one package with a struct key.
3. Refactor a function that pulls a `*sql.DB` from context. Put `*sql.DB` on a struct field instead.

#### Advanced practical tasks

1. Design middleware that stores a trace id in context for a logging library that you do not control. Keep your business functions free of `Value` reads.
2. Read the `context` package docs on `WithValue`. Write five rules in your own words. Do not copy the page.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do goroutines, channels, and `WaitGroup` each wait for work to finish?
2. When do you choose `select` plus context instead of a bare channel receive?
3. How do a mutex and a channel each prevent a data race on a counter?
4. What stop rules prevent a goroutine leak in a worker that reads a channel?
5. Why does `WithValue` not replace a timeout, and why does a timeout not replace a function parameter?

#### Easy practical tasks

1. Start two goroutines. Each sends one number on a channel. `main` receives both, then returns. Use no `Sleep`.
2. Protect a slice append with a mutex while four goroutines add one item each. Wait with a `WaitGroup`. Print `len`.
3. Run `go test -race` on that program as a test file. Record pass or fail.
4. Cancel a `WithTimeout` context of 10 milliseconds in a `select` that also waits on an idle channel. Print which case ran.

#### Medium practical tasks

1. Build a producer and a consumer. Close the channel when the producer ends. Cancel a context if the consumer is too slow.
2. Create a race on purpose. Capture the `-race` report. Fix it. Add a leak test with `NumGoroutine`.
3. Write `func fetch(ctx context.Context, url string)` that uses `NewRequestWithContext` and returns `ctx.Err()` when the context is done.

#### Advanced practical tasks

1. Implement a worker pool that respects `ctx.Done()`, closes jobs only from the owner, and passes `-race`. Document every goroutine.
2. Write a one-page policy: mutex versus channel versus context versus `WithValue`. Apply the policy to a tiny in-memory cache with a stop method.
