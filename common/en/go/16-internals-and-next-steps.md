# 16. Internals and Next Steps

## Description

This topic explains how Go stores values, how the garbage collector reclaims memory, and how slices and maps work inside. You also learn when to avoid reflection and cgo. The last sections point to Effective Go, the standard library, routers, gRPC, and further reading.

Use one term for each concept. Allocation is the act of reserving memory. Escape analysis is the compiler pass that decides stack versus heap. The garbage collector (GC) reclaims heap objects that the program can no longer reach. Do not treat micro-benchmarks as proof without `testing.B` and a clear hypothesis.

Complete topics 1 to 15 before this topic. This topic is a map, not a full runtime manual.

---

## Stack, heap, and escape analysis

Each goroutine has a stack. The stack holds local variables that do not escape. Stack allocation is cheap. The runtime grows and shrinks the stack as the goroutine needs more space.

The heap holds objects that live after the function returns, or that the compiler cannot keep on the stack. The garbage collector manages the heap. Heap allocation is more expensive than stack allocation. Heap objects add GC work.

Escape analysis is the compiler pass that decides the location. The compiler does not follow a simple "new means heap" rule. A value that you create with `&T{}` may stay on the stack if nothing retains the pointer. A local variable may move to the heap if you return its address.

Common reasons for a value to escape to the heap:

- you return a pointer to the value
- you store the value in an interface value (the interface may need a heap object)
- you send the value on a channel that outlives the function
- you store the pointer in a struct or slice that already lives on the heap
- you close over the variable in a goroutine that runs after the function returns

Inspect escape analysis with:

```text
go build -gcflags=-m
```

Add `-m` twice for more detail: `-gcflags="-m -m"`. The compiler prints lines such as `moved to heap` and `escapes to heap`. Read those lines as hints. They are not a stable public API. The wording can change between Go versions.

```go
func makeInt() *int {
	x := 42
	return &x // x escapes
}

func sum(a, b int) int {
	return a + b // a and b stay on the stack
}
```

Do not force heap allocation without a reason. Write clear code. Profile first. Then change allocation only when a profile shows a cost.

A large local array can be expensive on the stack. For large buffers, use `make([]byte, n)` and accept a heap slice, or reuse a buffer.

The stack is not shared between goroutines. Do not hide pointers from the runtime with `unsafe` unless you know the rules.

### Questions

#### Theoretical questions

1. What lives on a goroutine stack?
2. What lives on the heap?
3. What is escape analysis?
4. Why can a returned `&x` move `x` to the heap?
5. How do you print escape decisions from the compiler?

#### Easy practical tasks

1. Build a file that returns `&x` and a file that returns `x + 1`. Run `go build -gcflags=-m`. Write the escape lines.
2. Write five sentences that compare stack and heap.
3. Store an `int` in an `any` variable. Check whether `-m` says the value escapes.
4. Read `go doc cmd/compile`. Write where `-m` is documented.

#### Medium practical tasks

1. Return a pointer to a local struct. Then rewrite the function to return the struct by value. Compare `-m` output.
2. Send a pointer on a channel from a function that then returns. Explain the lifetime.
3. Benchmark a function that allocates every call versus a function that fills a stack struct. Use `testing.B`.

#### Advanced practical tasks

1. Find one escape in a small HTTP handler (interface conversion or `fmt`). Write whether you must fix it.
2. Read a Go blog post on escape analysis. Write one page in your own words. Do not copy the post.

---

## Garbage collector (high-level)

The garbage collector finds heap objects that the program cannot reach. It reclaims that memory. You do not call `free` in Go.

Go uses a concurrent, tri-color, mark-and-sweep collector in current versions. The collector runs with your program. It stops all goroutines only for short times. Those stops are stop-the-world (STW) pauses. Modern Go keeps STW short.

A simple picture:

1. Mark: the collector traces pointers from roots (stacks, globals, registers).
2. Objects that it does not mark are garbage.
3. Sweep: the collector reclaims unmarked objects.

The collector is not a real-time system. A huge heap or a huge number of pointers can still cost CPU. Allocation rate matters. Many short-lived objects create work.

`GOGC` sets a target heap growth percent. The default is `100`. That value means the collector aims to run after the heap grows by about 100% from the live size. `GOGC=off` disables GC. Do not disable GC in production.

`GOMEMLIMIT` (Go 1.19 and later) sets a soft memory limit. The collector works harder near that limit. Use it when a container has a memory cap.

`runtime.GC()` starts a collection. Use it in tests or in tools. Do not call it in a hot request path.

`runtime.ReadMemStats` and `pprof` heap profiles show allocation. Topic 11 covers `pprof`.

You do not tune GC first. You fix leaks (goroutines, unbounded slices, global caches). Then you profile. Then you change `GOGC` or `GOMEMLIMIT` with a measurement.

### Questions

#### Theoretical questions

1. What does the garbage collector reclaim?
2. What is a stop-the-world pause?
3. What does `GOGC=100` mean at a high level?
4. What does `GOMEMLIMIT` bound?
5. Why do you not call `runtime.GC()` on every request?

#### Easy practical tasks

1. Print `GOGC` with `os.Getenv` or document the default `100`.
2. Run a loop that appends to a slice of large structs. Print `runtime.MemStats.HeapAlloc` before and after (read `go doc runtime.MemStats`).
3. Write four sentences on mark and sweep in your own words.
4. Read `go doc runtime/debug.SetGCPercent`. Write what `SetGCPercent(-1)` does.

#### Medium practical tasks

1. Allocate in a loop. Call `runtime.GC()` once. Compare `HeapAlloc` before and after.
2. Capture a heap profile of a program that keeps a global slice growing. Name the function that allocates.
3. Explain a goroutine leak as a GC problem: the objects stay reachable.

#### Advanced practical tasks

1. Run the same allocation loop with `GOGC=50` and `GOGC=200`. Compare pause or CPU with `pprof` or a simple timer. Write the limits of your test.
2. Read the Go GC guide on go.dev. Write eight short sentences. Do not copy the page.

---

## Slice and map internals

A slice header is three words: a pointer to an array, a length, and a capacity. Topic 5 covers the programmer view. This section adds the internal view.

`append` grows the backing array when `len` equals `cap`. The runtime picks a new capacity. The growth factor is not a promise in your code. Old arrays become garbage when no slice still points to them. A small slice that still points to a huge array keeps the huge array alive. Copy the small part out.

A slice of pointers or a slice of interfaces holds pointers. The GC must scan those pointers. A slice of numbers is cheaper to scan.

A string is a header: a pointer and a length. The bytes are immutable. A slice of a string shares the backing bytes. `string` and `[]byte` conversions can copy. Go 1.22 and later optimize some conversions. Do not rely on a copy or a share without a need.

A map is a hash table. The runtime stores buckets. A map value is a pointer to that table. Assignment shares the table. Topic 5 covers that rule.

Growth of a map allocates new buckets and evicts overflow. Iteration order is random on purpose. The runtime randomizes the start so that programs do not depend on order.

An empty map still has a header. A `nil` map has no table. A write to a `nil` map panics.

Concurrent write to a map is fatal. The runtime detects many races and crashes with `concurrent map writes`. The race detector also reports races.

Do not use `unsafe` to walk map buckets. Use `range` and APIs.

`clear(s)` for a slice sets elements to zero and keeps length in recent Go versions. `clear(m)` deletes all keys. Read the current spec for the slice `clear` rule in your Go version.

### Questions

#### Theoretical questions

1. What three parts does a slice header store?
2. Why can a small subslice keep a large array alive?
3. What does a map assignment copy?
4. Why is map iteration order random?
5. What happens on concurrent map write?

#### Easy practical tasks

1. Print `unsafe.Sizeof` of a slice variable (the header). Print `len` and `cap` of a slice with 3 elements.
2. Show a subslice that keeps a 1_000_000-byte array alive. Then `copy` the one byte you need.
3. Assign a map to a second variable. Insert one key. Print the first map.
4. Range a map three times. Record whether the print order changed.

#### Medium practical tasks

1. Append until `cap` jumps. Print `len` and `cap` in a loop. Write that the jump sizes are an implementation detail.
2. Convert a large `[]byte` to `string` and back. Benchmark if you want. Write whether you must copy for safety.
3. Trigger `concurrent map writes` with two goroutines (short test). Then add a mutex.

#### Advanced practical tasks

1. Read the runtime map documentation or a Go blog post on maps. Write how buckets and growth work in six sentences of your own.
2. Design an API that returns a slice that must not alias caller memory. Prove isolation with a test.

---

## Reflection and cgo (when to avoid them)

Package `reflect` inspects types and values at run time. `encoding/json` uses reflection to read struct fields and tags. You can write `reflect.ValueOf(x).Kind()`.

Reflection is slow compared with generated or typed code. Reflection can panic on a wrong kind. Reflection breaks some compiler checks. You lose autocomplete and type safety.

Use reflection when you write a generic encoder, a test helper that prints any value, or a plugin-style registry. Do not use reflection to read a field that you can name in source. Do not use reflection in a hot request path without a measurement.

```go
v := reflect.ValueOf(p)
fmt.Println(v.Kind(), v.Type())
```

`unsafe` lets you break the type system. `unsafe.Pointer` can point at any address. Misuse causes memory corruption. The garbage collector may not see a hidden pointer. Do not use `unsafe` in application services. The standard library uses `unsafe` in a few places with strict review.

cgo calls C from Go. You write `import "C"` and C source in comments. cgo is useful for a vendor library that has no Go API. Costs:

- slower builds
- more complex cross-compilation
- extra threads and calling convention cost
- `CGO_ENABLED=0` builds fail
- harder Docker scratch images

Do not use cgo to call a tiny C function that you can write in Go. Prefer a pure Go driver (for example `modernc.org/sqlite` or `pgx`) when you can.

If you must use cgo, isolate it in one package. Set `CGO_ENABLED=1` only for that build. Document the C toolchain.

### Questions

#### Theoretical questions

1. What does `reflect` inspect?
2. Why is reflection a poor default in a hot path?
3. What is the main risk of `unsafe`?
4. Name three costs of cgo.
5. When is cgo a reasonable choice?

#### Easy practical tasks

1. Print `reflect.TypeOf` of an `int`, a `string`, and a struct.
2. Read `go doc reflect.Kind`. List five kinds.
3. Write four sentences on when you must not use `unsafe`.
4. Run `go env CGO_ENABLED`. Write the value.

#### Medium practical tasks

1. Use `reflect` to read an exported struct field by name. Then delete that code and use `p.Name`.
2. Build a module with `CGO_ENABLED=0`. Confirm that a pure Go program builds.
3. Find one standard library package that uses `reflect` (`encoding/json`). Write why it needs it.

#### Advanced practical tasks

1. Compare a json marshal of a struct with a hand-written encode of the same fields. Benchmark both. Write when reflection is acceptable.
2. Read the cgo documentation. Write a one-page "avoid unless" list for your team.

---

## Effective Go and reading the standard library

[Effective Go](https://go.dev/doc/effective_go) is the official style and design guide. It covers formatting, commentary, names, control structures, functions, data, initialization, methods, interfaces, and concurrency. Read it after you know the syntax.

The language specification is [https://go.dev/ref/spec](https://go.dev/ref/spec). Use it when you need an exact rule. Do not start there as a beginner.

The standard library is a teacher. Read small packages first:

- `errors`
- `io`
- `bytes`
- `strings`
- `sort`
- `context`
- `net/http` (selected files)

Open source with `go doc -src` or in `$GOROOT/src`. Read exported functions and the tests next to them. Tests in the standard library show table-driven style.

Copy style, not large internals. The runtime is not a model for application structure.

Other official pages:

- [A Tour of Go](https://go.dev/tour/)
- [Go Blog](https://go.dev/blog/)
- [Go by Example](https://gobyexample.com/)
- [pkg.go.dev/std](https://pkg.go.dev/std)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)

Read release notes when you upgrade. A new `slices` function can replace your helper.

Do not treat random blog posts as stronger than Effective Go and the spec. Check the Go version that a post uses.

### Questions

#### Theoretical questions

1. What is Effective Go for?
2. When do you open the language specification?
3. Why is the standard library a good reading list?
4. What command prints source for `fmt.Println`?
5. Why must you check the Go version of a tutorial?

#### Easy practical tasks

1. Read Effective Go "Formatting" and "Commentary". Change one of your files to match those rules. Record the changes.
2. Run `go doc -src errors.New`. Write how the function stores the string.
3. Open `io/io.go` in `GOROOT`. Write the `Reader` interface as you see it.
4. Bookmark Effective Go, the spec, and pkg.go.dev.

#### Medium practical tasks

1. Read `context/context.go` documentation comments. Write how `WithCancel` relates to the tree model from topic 10.
2. Find a table-driven test in the standard library. Write the file path and what the table checks.
3. Complete any Tour of Go pages that you skipped. Map each page to a topic number in this path.

#### Advanced practical tasks

1. Read Effective Go "Interfaces" and "Concurrency". Review one of your packages against five rules. List matches and misses.
2. Pick `net/http/server.go` and read `Serve` or `ListenAndServe` at a high level. Write a one-page map of types. Do not copy source.

---

## Routers, gRPC, and further reading

Go 1.22 `ServeMux` covers many HTTP APIs. You may still choose a router for extra features: middleware stacks, automatic OpenAPI, or older Go versions. Popular routers include `chi` and `echo`. Evaluate a module before you add it. Prefer the standard mux until you feel a real limit.

gRPC is an RPC system over HTTP/2. It uses Protocol Buffers. The Go implementation is `google.golang.org/grpc`. Use gRPC when many internal services need typed contracts and streaming. Use JSON HTTP when browsers and simple clients matter. Do not add gRPC to a single public CRUD API without a reason.

Further skills after this path:

- worker pools and `errgroup` (`golang.org/x/sync/errgroup`)
- `sqlc` or careful ORM use (topic 14)
- OpenTelemetry for traces
- more `pprof` and benchmarks (topic 11)
- security review: TLS, secrets, validation (topic 15)
- reading Go proposals and release notes

Suggested practice order from the learning path:

1. Tour of Go plus tiny CLI programs
2. Structs, slices, maps: an in-memory list
3. HTTP JSON API with `net/http`
4. Tests, table-driven cases, and `-race`
5. Context timeouts and graceful shutdown
6. `database/sql`
7. A worker pool
8. Dockerize and cross-compile
9. Profile one slow path with `pprof`
10. Rebuild a small tool that you actually use

Official resources stay the source of truth: go.dev, the spec, Effective Go, and pkg.go.dev.

Do not learn every extra module at once. Ship a small service with the standard library. Add one dependency when the standard library is not enough.

### Questions

#### Theoretical questions

1. When is the standard `ServeMux` enough?
2. What problem does gRPC solve that JSON HTTP may not solve?
3. Why must you not add a router only because a tutorial uses one?
4. What is `errgroup` for at a high level?
5. What is the last practice item in the suggested order, and why does it matter?

#### Easy practical tasks

1. Write five sentences that compare `ServeMux` and a third-party router.
2. Open the gRPC Go quick start page. Write three facts. Do not implement a full service yet.
3. Copy the ten practice items into your notes. Mark which items you already finished.
4. List four official URLs from this section.

#### Medium practical tasks

1. Build or revisit a JSON API with only `net/http`. List features that would force a router.
2. Read `go doc golang.org/x/sync/errgroup` (after `go get` if needed). Write how `Wait` returns the first error.
3. Profile one function from your API with `pprof`. Write one finding.

#### Advanced practical tasks

1. Sketch a system with one JSON edge API and one internal gRPC service. Write why each protocol fits.
2. Write a six-month reading plan: one official doc per week and one small project per month. Keep the standard library first.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do escape analysis, the heap, and the garbage collector share the story of a pointer that you return?
2. Why do slice capacity and map buckets matter when you debug memory growth?
3. When do reflection, `unsafe`, and cgo each become the wrong tool?
4. How do Effective Go and the standard library source replace random tutorials?
5. What do you learn next after a standard-library JSON service: router, gRPC, or practice item 10?

#### Easy practical tasks

1. Run `-gcflags=-m` on `makeInt` and `sum`. Paste one line per function into your notes and explain it.
2. Print `MemStats.HeapAlloc` after you append 100_000 integers. Then `nil` the slice and call `runtime.GC()`. Print again.
3. Write a one-page glossary: stack, heap, GC, slice header, map bucket, reflect, cgo.
4. Open Effective Go and the Tour. Write which document you open for style and which you open for a first lesson.

#### Medium practical tasks

1. Find a slice alias bug in an old exercise. Fix it with `copy`. Show `-race` or a unit test.
2. Replace a `reflect` field read with typed code. Benchmark if the path is hot.
3. Read `errors` and `io` source comments. List three style habits that you will copy.

#### Advanced practical tasks

1. Write a report of one page: memory of one service (allocs, maps, goroutines), tools (`-m`, `pprof`, `GOGC`), and one change you will not make without data.
2. Rebuild a tool that you use. Use modules, tests, `-race`, and a small `cmd/` layout. Write what this path still does not cover for that tool.
