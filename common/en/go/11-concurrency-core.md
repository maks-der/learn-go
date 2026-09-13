# 11. Concurrency (Core)

## Description

Go runs many functions at the same time. A goroutine is a function that the Go runtime schedules. A channel moves values between goroutines. The `sync` package protects shared data when a channel is not the simple choice.

This topic teaches the core tools. Complete this topic before you build servers, workers, or other concurrent programs. Topic 12 covers `context.Context` in more detail. Topic 19 covers advanced concurrency.

Use one term for each concept. A goroutine is not an operating-system thread. A channel send is not a function call. A data race is not the same as a deadlock.

---

## Goroutines: `go f()`

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

Do not start a goroutine and ignore its error. Return the error on a channel, or collect it with a `WaitGroup` and a mutex, or use a helper such as `errgroup` in a later topic.

### Questions

#### Theoretical questions

1. What happens to other goroutines when `main` returns?
2. Where does the compiler evaluate the arguments of `go f(a, b)`?
3. Why is a goroutine not the same as an operating-system thread?
4. What changed for loop variables in Go 1.22, and why does that matter for `go func()` inside a loop?
5. When must you wait for a goroutine to finish?

#### Easy practical tasks

1. Write `package main` that starts one goroutine. The goroutine prints `hello`. The `main` function sleeps for 200 milliseconds, then returns. Run the program.
2. Remove the sleep. Run the program ten times. Write what you observe.
3. Start three goroutines that each print a different number. Use a `time.Sleep` in `main` so that the prints appear.
4. Write five sentences that describe a goroutine. Use only facts from this section.

#### Medium practical tasks

1. Start a goroutine that sends one `int` on a channel. Receive that value in `main` and print it. Do not use `Sleep`.
2. Write a loop that starts one goroutine per slice element. Each goroutine prints the element. Confirm the program is safe on Go 1.22.
3. Start a goroutine that panics. Observe the process. Write how a panic in one goroutine affects the process.

#### Advanced practical tasks

1. Measure how many goroutines you can start before the program is slow or the machine is short of memory. Record the number and the limit that you hit (CPU, memory, or time). Stop the test before the machine becomes unusable.
2. Write a small program that starts work in a goroutine and must observe an error from that work. Use a channel of type `error`. Document the protocol (when you close the channel, who sends).

---

## The Go scheduler (high-level: M:N, `GOMAXPROCS`)

The Go scheduler runs many goroutines on a smaller number of operating-system threads. This model is an M:N scheduler. M goroutines run on N threads.

A simple picture uses three names:

- G is a goroutine.
- M is an operating-system thread.
- P is a logical processor. A P holds the run queue and other scheduler state.

`GOMAXPROCS` is the maximum number of Ps that can run Go code at the same time. From Go 1.5, the default is the number of logical CPUs. Read the value with `runtime.GOMAXPROCS(0)`. Set the value with `runtime.GOMAXPROCS(n)` or with the `GOMAXPROCS` environment variable.

A goroutine that runs a blocking system call can detach from its P. The scheduler can create or reuse another thread so that other goroutines continue. A goroutine that waits on a channel or a mutex does not need an operating-system thread while it waits.

The scheduler is cooperative and preemptive. A long loop without function calls can still be preempted in modern Go. You do not insert manual yield calls for ordinary code.

Do not treat `GOMAXPROCS` as a throughput knob that you turn without a measurement. For CPU-bound work, a value near the number of logical CPUs is the usual start. For I/O-bound work, many goroutines can wait while few threads run.

This section stays at a high level. Topic 20 covers the G, M, and P model in more depth.

### Questions

#### Theoretical questions

1. What does M:N mean for the Go scheduler?
2. What is `GOMAXPROCS`?
3. What is the default value of `GOMAXPROCS` in Go 1.22 and later?
4. What happens when a goroutine makes a blocking system call?
5. Does a waiting goroutine hold an operating-system thread the whole time? Explain.

#### Easy practical tasks

1. Run a program that prints `runtime.GOMAXPROCS(0)` and `runtime.NumCPU()`. Write both numbers.
2. Set the `GOMAXPROCS` environment variable to `1`. Run the same program. Write the new `GOMAXPROCS` value.
3. Print `runtime.NumGoroutine()` at the start of `main`. Write the number and what it represents.
4. Read `go doc runtime.GOMAXPROCS`. Write the purpose of the function in two sentences.

#### Medium practical tasks

1. Run a CPU-bound loop in several goroutines with `GOMAXPROCS=1` and with the default. Compare wall time. Write the two times and one reason for the difference.
2. Draw a diagram with three Gs, two Ms, and two Ps. Label which G is running and which G is in a run queue.
3. Start many goroutines that only sleep. Print `NumGoroutine` and `GOMAXPROCS`. Explain why the first number can be much larger than the second.

#### Advanced practical tasks

1. Read the Go blog post about the scheduler (search for "Go scheduler" on go.dev/blog). Write one page with the meaning of G, M, and P. Do not copy the post.
2. Use `GODEBUG=schedtrace=1000` on a small concurrent program for a few seconds. Write what one trace line shows. Stop the program.

---

## Channels: send, receive, close

A channel is a typed conduit. You create a channel with `make`.

```go
ch := make(chan int)        // unbuffered
ch := make(chan int, 4)     // buffered, capacity 4
```

Send a value with `ch <- v`. Receive a value with `v := <-ch`. The receive operator can stand on the right of a short declaration or an assignment.

The comma-ok form reports whether the channel is open:

```go
v, ok := <-ch
```

`ok` is `false` when `ch` is closed and empty. `v` is the zero value of the element type in that case.

`close(ch)` marks the channel as closed. Receivers drain remaining buffered values, then they see the closed state. After that, every receive returns the zero value.

Channel rules:

- A send to a closed channel panics.
- A close of a closed channel panics.
- A close of a nil channel panics.
- A send on a nil channel blocks forever.
- A receive on a nil channel blocks forever.

The sender closes the channel. The receiver does not close the channel when more than one sender exists. You do not close a channel if no code uses `range` or comma-ok to detect the end.

A channel value is a reference. Copy of the channel variable still refers to the same channel. The zero value of a channel type is `nil`.

### Questions

#### Theoretical questions

1. What does the comma-ok receive form tell you?
2. Who must close a channel, and why?
3. What happens when you send on a closed channel?
4. What happens when you receive from a closed, empty channel?
5. What happens when you send or receive on a nil channel?

#### Easy practical tasks

1. Create an unbuffered `chan string`. Send `"ok"` from a goroutine. Receive in `main` and print the value.
2. Close a channel in the sender. Use comma-ok in the receiver. Print `v` and `ok` for two receives after the close.
3. Write the five panic or block rules from this section in a two-column table: "Operation" and "Result".
4. Declare `var ch chan int` without `make`. Explain why a send on `ch` does not proceed.

#### Medium practical tasks

1. Write a function `func produce(n int) <-chan int` that sends `0` to `n-1` and then closes the channel.
2. Cause a panic with a send on a closed channel. Recover in `main` only to record the panic. Then rewrite the program so that it does not panic.
3. Share one channel among three receivers. Close the channel after ten sends. Confirm that the receivers stop. Write how you stop them.

#### Advanced practical tasks

1. Design a protocol for one channel that carries both results and a terminal error. Document the message type and the close rule. Implement a small demo.
2. Show that two variables of type `chan int` can refer to the same channel. Show that a nil channel is not that shared channel.

---

## Unbuffered vs buffered channels

An unbuffered channel has capacity 0. A send on an unbuffered channel waits until another goroutine receives. A receive waits until another goroutine sends. The two goroutines meet. This meeting is a handshake. The value is transferred, and both sides continue.

A buffered channel has capacity `n`, where `n` is greater than 0. A send stores the value in the buffer and continues while the buffer is not full. A receive takes a value from the buffer and continues while the buffer is not empty. A send blocks when the buffer is full. A receive blocks when the buffer is empty.

```go
unbuf := make(chan int)     // capacity 0
buf := make(chan int, 8)    // capacity 8
```

`cap(ch)` is the capacity. `len(ch)` is the number of values in the buffer at that moment. Do not use `len` to decide that a send or a receive is safe. Another goroutine can change the buffer before your next operation.

Use an unbuffered channel when you want the sender to wait until the receiver takes the value. Use a small buffer when you want to decouple bursts of work. Do not pick a large buffer as a guess for speed. A large buffer can hide a slow consumer. A large buffer can use much memory.

A buffer does not replace a wait for completion. If `main` returns, buffered values that nobody received are lost.

### Questions

#### Theoretical questions

1. What is the capacity of an unbuffered channel?
2. When does a send on a buffered channel block?
3. When does a receive on a buffered channel block?
4. Why must you not use `len(ch)` to decide that a send is safe?
5. When do you choose an unbuffered channel instead of a buffered channel?

#### Easy practical tasks

1. Create `make(chan int, 2)`. Send two values from `main` without a receiver. Print `len` and `cap`. Then receive both values.
2. Repeat the two sends on an unbuffered channel without a receiver. Observe that the program blocks. Stop it. Write why it blocks.
3. Draw two diagrams: handshake on an unbuffered channel, and a buffer with three slots.
4. Write four sentences that compare unbuffered and buffered channels.

#### Medium practical tasks

1. Build a producer that is faster than the consumer. Use capacity `1` and capacity `8`. Measure how long the producer waits. Write the two times.
2. Send three values on a buffered channel of capacity `3`, then close the channel. Receive with a loop until comma-ok is false. Print each value.
3. Explain a case where a large buffer hides a bug. Write that case in six sentences.

#### Advanced practical tasks

1. Implement a semaphore with a buffered channel of empty structs: `make(chan struct{}, n)`. Limit five concurrent workers. Document acquire and release.
2. Profile memory of a program that fills a channel of large structs with a large capacity. Write why a channel of pointers or a smaller capacity can be better.

---

## `select` and timeouts

A `select` statement waits for a set of channel operations. Each `case` is a send or a receive. `select` runs one ready case. If more than one case is ready, `select` chooses one at random. If no case is ready, `select` blocks until one case is ready.

```go
select {
case v := <-in:
    handle(v)
case out <- v:
    // sent
case <-ctx.Done():
    return ctx.Err()
}
```

A `default` case runs when no other case is ready. Use `default` for a non-blocking attempt. Do not add `default` when you mean to wait.

A timeout is a receive from a channel that becomes ready after a duration:

```go
select {
case v := <-ch:
    use(v)
case <-time.After(2 * time.Second):
    return errTimeout
}
```

`time.After` returns a channel. The runtime sends the current time on that channel after the duration. If `select` takes another case first, the timer still runs until it fires. That is acceptable for rare timeouts. For a tight loop, use `time.NewTimer`, receive or stop the timer, and reset it. A timer that you never stop can keep a goroutine and a channel alive until the duration ends.

`time.Tick` in a `select` loop is a common leak. Prefer `time.NewTicker` and call `Stop` when the loop ends.

`select` with one `case` is valid. A `select` with no cases blocks forever.

### Questions

#### Theoretical questions

1. What does `select` do when two cases are ready?
2. What is the role of `default` in a `select`?
3. How do you express a timeout with `select`?
4. Why can `time.After` in a loop waste memory?
5. What happens when `select` has no `case` and no `default`?

#### Easy practical tasks

1. Write a `select` with two receive cases from two channels. Send one value to one channel from a goroutine. Print which case ran.
2. Add a `default` to a `select` on an empty unbuffered channel. Confirm that `default` runs.
3. Receive from a channel with a 100 millisecond timeout. Do not send. Print that the timeout case ran.
4. Read `go doc time.After`. Write the return type and the meaning of the value that arrives.

#### Medium practical tasks

1. Send to two channels in a loop with `select`. Count how many times each case runs in 1000 sends. Write whether the counts are equal.
2. Replace `time.After` in a loop of 10000 iterations with `time.NewTimer` and `Reset`. Compare allocations with `go test -bench` and `-benchmem` if you wrap the loop in a benchmark.
3. Implement a non-blocking send helper: return `true` if the send occurred, `false` if the channel was not ready.

#### Advanced practical tasks

1. Build a `select` that waits on work, a ticker, and `ctx.Done()`. Stop the ticker on cancel. Show that the goroutine returns.
2. Demonstrate a leak with `time.Tick` in a function that you call many times. Then fix the function with `NewTicker` and `Stop`. Use `runtime.NumGoroutine` or `pprof` as evidence.

---

## Range over channels

A `for` range over a channel receives values until the channel is closed and the buffer is empty.

```go
for v := range ch {
    handle(v)
}
```

The loop ends only after `close(ch)` and after all buffered values are received. If nobody closes the channel, the loop waits forever. That wait is a goroutine leak when the rest of the program continues.

The range form does not expose the comma-ok flag. You do not need that flag when you only want values.

You can receive in a plain `for` loop when you need `ok` or when you need `select` for cancel:

```go
for {
    select {
    case v, ok := <-ch:
        if !ok {
            return
        }
        handle(v)
    case <-ctx.Done():
        return
    }
}
```

Do not close the channel in the range loop that receives from it. The sender closes the channel.

Range over a nil channel blocks forever. Range over a closed empty channel does not enter the loop body.

From Go 1.22, `for i := range n` ranges over integers. That form is not a channel range. Use `range ch` only when `ch` is a channel.

### Questions

#### Theoretical questions

1. When does `for v := range ch` stop?
2. What happens if the sender never closes the channel?
3. Why does the receiver not close the channel inside the range loop?
4. What does range over a closed empty channel do?
5. How do you combine range-style receive with cancellation?

#### Easy practical tasks

1. Send three integers on a channel. Close the channel. Range over the channel and print each integer.
2. Omit the close. Run the range in a goroutine. Observe that the program does not finish unless `main` exits. Write the cause.
3. Range over a channel that you close with no sends. Confirm that the loop body does not run.
4. Write a function that receives all values from `<-chan string` with `range` and returns them as a slice.

#### Medium practical tasks

1. Fill a buffered channel, close it, then range. Confirm that you receive the buffered values after `close`.
2. Start two receivers that both range over the same channel. Send ten values and close. Write how the values split.
3. Rewrite a `range` loop as a `for` plus comma-ok. Keep the same behavior.

#### Advanced practical tasks

1. Write a consumer that ranges over a channel and also stops on `ctx.Done()`. You cannot put `range` and `select` in one statement. Show a correct loop.
2. Build a test that fails if a range loop leaks a goroutine. Use a timeout or a `WaitGroup` as the signal that the loop ended.

---

## Directional channels (`chan<-`, `<-chan`)

A channel type can restrict direction.

- `chan T` is bidirectional. You can send and receive.
- `chan<- T` is send-only.
- `<-chan T` is receive-only.

The arrow points the way the data flows. `chan<- T` is a channel that accepts `T`. `<-chan T` is a channel that produces `T`.

You convert a bidirectional channel to a directional channel. You cannot convert a directional channel back to bidirectional.

```go
func produce(out chan<- int) {
    out <- 1
    close(out)
}

func consume(in <-chan int) {
    for v := range in {
        fmt.Println(v)
    }
}

func main() {
    ch := make(chan int)
    go produce(ch)
    consume(ch)
}
```

Directional types document the API. The compiler rejects a send on `<-chan T`. The compiler rejects a receive on `chan<- T`. The compiler rejects `close` on a receive-only channel.

Use directional types on function parameters. Keep the `make` result as `chan T` at the owner.

A return type `<-chan T` tells the caller that the function owns the send side. The caller only receives.

### Questions

#### Theoretical questions

1. What operations are legal on `chan<- T`?
2. What operations are legal on `<-chan T`?
3. Can you convert `<-chan T` to `chan T`? Why?
4. Why do function parameters use directional channel types?
5. Who closes a channel that a function returns as `<-chan T`?

#### Easy practical tasks

1. Write `produce` and `consume` with directional parameters. Connect them with one `chan int`.
2. Try to receive inside a function that takes `chan<- int`. Record the compiler error.
3. Try to close a `<-chan int` parameter. Record the compiler error.
4. Draw the arrows for `chan<- int` and `<-chan int`. Label send and receive.

#### Medium practical tasks

1. Write `func counters() <-chan int` that sends three values and closes the channel. Range over the result in `main`.
2. Pass a bidirectional channel to a send-only parameter and to a receive-only parameter. Show that both calls compile.
3. Split a pipeline into two functions with directional channels. Document which function closes which channel.

#### Advanced practical tasks

1. Design an API with two return channels: `<-chan Result` and `<-chan error`. Document close rules and a safe `select` for the caller.
2. Explain why a struct field of type `chan<- T` is rare. Give one valid use and one misuse.

---

## Common patterns

These four patterns appear in many Go programs. Treat them as one toolkit. Combine them. Do not copy a pattern when a single goroutine is enough.

**Worker pool.** One channel carries jobs. A fixed number of worker goroutines receive jobs. Each worker sends a result on a second channel or writes to a safe collector. A `WaitGroup` waits for the workers. Then you close the result channel. The pool bounds concurrency. The bound is the number of workers.

**Fan-out.** One source sends work to many goroutines. The source can be a channel that many workers share (the worker pool is a fan-out). Fan-out increases throughput when the work is independent.

**Fan-in.** Many producers send to one result channel. You merge the streams. Close the result channel only after every producer is done. A `WaitGroup` in a merge goroutine is the usual close rule: each producer calls `Done`, the merge goroutine waits, then it closes the result channel.

**Pipeline.** Each stage is a goroutine. Stage 1 receives input and sends to stage 2. Stage 2 sends to stage 3. Channels connect the stages. Close the output of a stage when that stage finishes. The next stage ranges until that close.

**Cancellation with `context.Context`.** Pass a `context.Context` into producers, workers, and stages. When the caller cancels, each goroutine returns. Check `ctx.Done()` in `select` next to channel operations. Topic 12 explains `context` in full. In this topic, use `context.WithCancel` or `context.WithTimeout` so that a slow stage does not run forever.

A safe pool or pipeline follows these rules:

1. The owner creates the channels.
2. The sender closes a channel.
3. Every goroutine can exit when the context is cancelled.
4. `WaitGroup.Add` happens before `go`.
5. You do not send on a closed channel.

### Questions

#### Theoretical questions

1. What problem does a worker pool solve?
2. What is the difference between fan-out and fan-in?
3. Who closes the output channel of a pipeline stage?
4. Why does a merge (fan-in) need a `WaitGroup` or an equivalent counter?
5. How does `context.Context` stop a worker that waits on a job channel?

#### Easy practical tasks

1. Write a pool with two workers. Send five jobs. Print five results. Wait with a `WaitGroup`.
2. Draw a pipeline with three stages. Label each channel and the close direction.
3. Write a fan-in of two `<-chan int` into one `<-chan int`. Close the output after both inputs close.
4. Add `context.WithTimeout` to a worker that sleeps. Cancel before the sleep ends. Confirm that the worker returns.

#### Medium practical tasks

1. Implement a worker pool that processes a slice of URLs as fake jobs (`time.Sleep` plus the URL string). Bound the pool to three workers.
2. Build a three-stage pipeline: generate integers, square them, print them. Close each stage in the correct order.
3. Combine fan-out and fan-in: split work to four workers, merge results, return a sorted slice in `main`.

#### Advanced practical tasks

1. Add cancellation to the three-stage pipeline. When the context ends, every stage returns and no goroutine remains. Prove the claim with `runtime.NumGoroutine` or a test.
2. Write a pool that stops when the first job fails. Cancel the context. Drain or avoid sends that can panic. Document the protocol.

---

## `sync.WaitGroup`, `sync.Mutex`, `sync.RWMutex`

`sync.WaitGroup` waits for a set of goroutines.

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    work()
}()
wg.Wait()
```

Call `Add` before `go`. Call `Add` with the number of goroutines, or call `Add(1)` once per goroutine. `Done` subtracts one. `Wait` blocks until the counter is zero. A negative counter panics. Do not copy a `WaitGroup` after you use it. From Go 1.25, `wg.Go(func() { ... })` calls `Add(1)`, runs the function, and calls `Done`. On Go 1.22 through 1.24, use `Add` and `Done`.

`sync.Mutex` is a mutual-exclusion lock. `Lock` acquires the lock. `Unlock` releases the lock. Use `defer mu.Unlock()` immediately after `Lock` when the critical section is the rest of the function. Protect every read and every write of the shared data.

```go
var mu sync.Mutex
var count int

mu.Lock()
count++
mu.Unlock()
```

`sync.RWMutex` allows many readers or one writer. `RLock` and `RUnlock` are for readers. `Lock` and `Unlock` are for writers. A writer waits until all readers unlock. Use `RWMutex` when reads are common and writes are rare. Measure before you replace a `Mutex`. An `RWMutex` is not always faster.

Do not hold a lock across a slow I/O call or a channel operation that can wait. Hold the lock only while you read or write the shared data.

The zero value of `WaitGroup`, `Mutex`, and `RWMutex` is ready to use. Do not copy these values. Share them with pointers when you pass them to functions.

### Questions

#### Theoretical questions

1. Why must `WaitGroup.Add` run before `go`?
2. What happens when `WaitGroup` counter becomes negative?
3. What is the difference between `Mutex` and `RWMutex`?
4. Why do you use `defer mu.Unlock()`?
5. Why must you not copy a `Mutex` after first use?

#### Easy practical tasks

1. Start three goroutines that each print a line. Wait with a `WaitGroup`. Do not use `Sleep`.
2. Protect a shared `int` counter with a `Mutex`. Increment it from four goroutines, 1000 times each. Print the final value. Confirm that it is 4000.
3. Remove the `Mutex` from the counter program. Run `go run -race .`. Write the race detector message.
4. Read `go doc sync.WaitGroup`. List `Add`, `Done`, and `Wait` in your own words.

#### Medium practical tasks

1. Use `RWMutex` on a map. Run eight readers and two writers. Keep all map access inside the correct lock.
2. Call `Add` inside the new goroutine instead of before `go`. Run the program many times. Write the failure that you see or explain why it can skip `Wait`.
3. Hold a `Mutex` during `time.Sleep`. Then move the sleep outside the lock. Compare how long other goroutines wait.

#### Advanced practical tasks

1. If your toolchain is Go 1.25 or later, rewrite a `WaitGroup` example with `wg.Go`. Compare the code with `Add` and `Done`. If your toolchain is older, write the same comparison from the documentation.
2. Find a deadlock: two mutexes locked in opposite order in two goroutines. Record the runtime error. Fix the lock order.

---

## `sync.Once`, `sync.Map`, `sync.Pool`

`sync.Once` runs a function once. Concurrent callers wait until that run finishes.

```go
var once sync.Once

func setup() {
    once.Do(func() {
        loadConfig()
    })
}
```

`Do` runs the function on the first call. Later calls return without a run. If the function panics, `Once` treats the run as done. A later `Do` does not retry. Do not copy a `Once` after first use.

`sync.Map` is a concurrent map. Use it when many goroutines read and write keys and the lock on a normal map is a measured problem. The API uses `any` keys and values: `Load`, `Store`, `LoadOrStore`, `LoadAndDelete`, `Delete`, `Swap`, `CompareAndSwap`, and `Range`. A normal `map` plus a `Mutex` or `RWMutex` is simpler for most programs. Do not copy a `sync.Map` after first use.

`sync.Pool` is a cache of temporary objects. `Get` returns a value of type `any`. Type-assert the result. `Put` returns an object to the pool. The garbage collector can discard pooled objects at any time. A pool is not a reliable store. Use a pool to reduce allocations of short-lived objects that you can reuse (for example, a `bytes.Buffer`). Do not put a dirty object in the pool without a reset.

```go
var bufPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer)
    },
}

b := bufPool.Get().(*bytes.Buffer)
b.Reset()
// use b
bufPool.Put(b)
```

`New` runs when `Get` finds no cached object. `New` can be nil. Then `Get` returns `nil` when the pool is empty.

### Questions

#### Theoretical questions

1. How many times does `Once.Do` run the function in the successful case?
2. What happens when the function in `Once.Do` panics?
3. When do you choose `sync.Map` instead of `map` plus `Mutex`?
4. Why is `sync.Pool` not a reliable store?
5. Why must you `Reset` a `bytes.Buffer` that you get from a pool?

#### Easy practical tasks

1. Call a setup function from three goroutines. Use `Once` so that setup prints `init` only once.
2. Store three keys in a `sync.Map`. Load one key with the comma-ok form (`Load` returns `value, ok`).
3. Put and get a `[]byte` or a `bytes.Buffer` from a `Pool`. Print a value after `Get`.
4. Write a table: type, main methods, one valid use, one misuse.

#### Medium practical tasks

1. Compare a `map` plus `RWMutex` with a `sync.Map` for a read-heavy workload. Write a benchmark for both. Record the results.
2. Show that `Pool` can return a new object after a `runtime.GC()` call. Write what you observe. Do not depend on one GC behavior in production code.
3. Use `LoadOrStore` to register a name only once. Print whether the store happened.

#### Advanced practical tasks

1. Build a small constructor that uses `Once` for lazy initialization and returns an error. `Once` does not retry after panic. Design a pattern that can retry on error (you may need a mutex instead of `Once`).
2. Implement a pool of `[]byte` with a maximum size. Do not `Put` a slice that is too large. Write a test that checks reuse.

---

## Race detector: `go test -race`

A data race occurs when two goroutines access the same memory, at least one access is a write, and the accesses are not synchronized. The race detector finds many data races at run time.

Run tests with the race detector:

```text
go test -race ./...
```

Run a main program with the race detector:

```text
go run -race .
go build -race -o app .
```

The detector instruments the program. The program uses more CPU and more memory. Use `-race` in development and in continuous integration. Do not enable `-race` on every production binary unless you have a special need and you measured the cost.

The detector reports the two stacks that raced. Fix the race with a channel, a mutex, or an atomic operation. Do not add a `time.Sleep` to hide a race. A sleep does not synchronize memory.

The race detector needs CGO on some platforms. On Windows, the official toolchain supports `-race` on `amd64`. If `-race` fails to build, read the error. Install the required C compiler when the message asks for it.

`go test -race` and `go test -covermode=atomic` work together when you want coverage and race detection.

A clean `-race` run does not prove that the program has no races. The detector sees races that the test run executes. Write tests that start the concurrent paths.

### Questions

#### Theoretical questions

1. What three conditions define a data race?
2. Which command runs tests with the race detector?
3. Why is `time.Sleep` not a fix for a race?
4. Why can a green `-race` run still miss a race?
5. What extra cost does the race detector add?

#### Easy practical tasks

1. Write a test that increments a shared `int` from two goroutines without a lock. Run `go test -race`. Save the report.
2. Add a `Mutex`. Run `go test -race` again. Confirm that the report is gone.
3. Run `go test -race -count=1 ./...` in a small module. Write the command output.
4. Read `go help test` and find `-race`. Write the flag description in one sentence.

#### Medium practical tasks

1. Create a race on a map write from two goroutines. Record whether the race detector reports it and whether the runtime panics.
2. Compare test time with and without `-race` on a package that starts many goroutines. Write the two times.
3. Fix a race by changing the design to a channel instead of a mutex. Keep `-race` clean.

#### Advanced practical tasks

1. Add `-race` to a script or a CI configuration for the module. Document when the job runs.
2. Produce a race that appears only with `t.Parallel` tests. Fix the shared state. Show the `-race` report before and after.

---

## Deadlocks, goroutine leaks

A deadlock occurs when every goroutine is waiting and no goroutine can make progress. The Go runtime detects a global deadlock and stops the process with a fatal error. The message lists the goroutines and their wait points.

Typical deadlock causes:

- A send on an unbuffered channel with no receiver, in the only remaining goroutine.
- Two goroutines wait on each other (A waits for B, B waits for A).
- A `WaitGroup.Wait` that never sees enough `Done` calls.
- Two mutexes locked in opposite order.
- A `select` with no ready case and no `default`, when no other goroutine will send.

A goroutine leak is different. The process continues. One or more goroutines stay blocked or stay in a loop. They hold memory and other resources. The runtime does not stop the process.

Typical leak causes:

- A receive on a channel that nobody will send on and nobody will close.
- A `for range ch` when `ch` is never closed.
- A `time.After` or `time.Tick` that you create in a loop and never stop.
- A goroutine that waits on `ctx.Done()` when you never cancel and you never set a deadline.
- A worker that you started, but you dropped the only cancel function.

Detect leaks in tests. Count goroutines before and after the work. Use `runtime.NumGoroutine` with care (other packages start goroutines). Prefer a test that waits for an exit signal with a timeout. Topic 13 covers tests. Topic 19 covers traces.

Fix a deadlock or a leak with a clear lifetime: who starts the goroutine, who closes the channel, who calls cancel, and who waits.

### Questions

#### Theoretical questions

1. What is the difference between a deadlock and a goroutine leak?
2. What does the runtime do when it detects a global deadlock?
3. Give two causes of a deadlock that involve channels.
4. Give two causes of a goroutine leak.
5. Why does a leak not always stop the process?

#### Easy practical tasks

1. Write a `main` that sends on an unbuffered channel and never receives. Run it. Save the deadlock message.
2. Write a goroutine that receives from a channel that `main` never sends on. Let `main` return. Write whether you see a deadlock (you do not, because `main` exits).
3. Keep that receiver goroutine alive with a `WaitGroup` that nobody completes. Observe the deadlock. Write the cause.
4. Make a two-column table: "Deadlock" and "Leak". Add three rows of examples.

#### Medium practical tasks

1. Create a two-mutex deadlock. Save the fatal error. Fix the lock order.
2. Write a function that leaks a goroutine with `for range` on an open channel. Then fix it with `close` or with context cancel.
3. Use `runtime.NumGoroutine` before and after a leaking function and after the fix. Write the three numbers.

#### Advanced practical tasks

1. Write a test that fails when a helper leaks a goroutine. Use a timeout. Then write a second version of the helper that does not leak.
2. Capture a goroutine profile of a leak with `pprof` or `runtime/pprof`. Identify the stuck function in the output.

---

## Prefer communication over shared memory

The Go proverb is: do not communicate by sharing memory; share memory by communicating. The idea is this: pass data on a channel. The receive transfers ownership. After the receive, one goroutine uses the value. You do not need a mutex for that value.

Channels make the flow visible. A pipeline, a worker pool, and a request handler that sends a result on a channel follow this idea.

A mutex is simpler when the data is a single struct in one place. Examples:

- A cache map that many handlers update.
- A counter or a set of counters (or use `sync/atomic` for a single counter).
- A graph of objects that many methods must update in place.

A mutex is the wrong tool when you copy the lock, when you forget to lock a read, or when you hold the lock during a slow call. A channel is the wrong tool when you create many channels only to protect a few fields of one struct.

Choose the tool that matches the ownership:

- One owner at a time: send the value on a channel, or send a message to a single loop (a goroutine that owns the data).
- Many readers and rare writers on one map: `RWMutex` or a copy of the map under a lock.
- Once-only init: `sync.Once`.
- Short critical section on one integer: `sync/atomic` (Topic 19).

Write the ownership in a comment when the code is not obvious. Tests with `-race` are part of the design.

### Questions

#### Theoretical questions

1. What does "share memory by communicating" mean in one sentence?
2. When is a `Mutex` simpler than a channel?
3. What is ownership of a value in a concurrent program?
4. Why is a large set of channels a poor fit for a few fields of one struct?
5. How does a single owner goroutine avoid shared-memory races?

#### Easy practical tasks

1. Rewrite a mutex-protected counter as a goroutine that receives `"inc"` messages on a channel and sends the current count on request.
2. Write six sentences: three cases for channels, three cases for mutexes.
3. Draw ownership of a `[]byte` that a producer sends to a consumer. Mark who may write the slice.
4. Read the Go proverb page or Effective Go on concurrency. Write the proverb and one example from this topic.

#### Medium practical tasks

1. Implement a small in-memory cache in two ways: map plus `RWMutex`, and a single owner goroutine with messages. Compare the code size and a simple benchmark.
2. Find a race in a design that sends a pointer on a channel and then writes the same memory in the sender. Fix the ownership.
3. List every `sync` type from this topic. For each type, write one sentence: "Use this when ...".

#### Advanced practical tasks

1. Design a service struct with a public method that is safe for concurrent use. Choose channel or mutex. Write why. Implement it and run `go test -race`.
2. Review a small public Go project. Find one channel design and one mutex design. Write why each choice fits.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of a value from `make(chan T)` through a send, a receive, and a `close`. Name the blocking rules at each step.
2. How do `GOMAXPROCS`, worker-pool size, and channel buffer capacity differ as limits?
3. A program passes `-race` and still hangs. Which class of bug do you look for first, and why?
4. Compare `select` plus `ctx.Done()` with `for range ch` for a consumer that must stop early.
5. When do you use `sync.Once`, when do you use `sync.Pool`, and when do you use a normal `map` plus `Mutex`?

#### Easy practical tasks

1. Write a module that starts two goroutines, moves one string on an unbuffered channel, waits with a `WaitGroup`, and runs `go test -race`.
2. Make a one-page cheat sheet: `go`, `make(chan T)`, `make(chan T, n)`, `close`, `select`, `range`, `WaitGroup`, `Mutex`, `-race`.
3. Draw one diagram that contains a worker pool, a fan-in merge, and a `context` cancel. Label every channel direction.
4. Run `go doc sync` and `go doc context`. Write three types from each package that this topic uses.

#### Medium practical tasks

1. Build a small pipeline: generate words, count letters in workers, merge counts. Add a timeout with `context.WithTimeout`. Run `go test -race`.
2. Write a test that starts a leaking goroutine on purpose, fails, then a second test that uses the fixed function and passes.
3. Convert a bidirectional channel API into directional parameters. Keep the same behavior. Show the compiler errors that you prevent.

#### Advanced practical tasks

1. Implement a bounded worker pool with cancellation, error return from the first failure, and no leaked goroutines. Document the protocol. Prove the bound with a test that limits concurrency.
2. Use `go test -race -count=20` on the pool. Then capture a goroutine profile after cancel. Confirm that only the expected goroutines remain.
