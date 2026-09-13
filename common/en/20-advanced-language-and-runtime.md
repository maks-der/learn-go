# 20. Advanced Language and Runtime

## Description

This topic covers tools that sit under normal application code. Reflection, code generation, Cgo, assembly, and the scheduler model are powerful. They also add cost and risk. Use them when a simpler Go design is not enough.

Prefer the type system, interfaces, and generators that you can read. Do not start a feature with `reflect` or Cgo. Measure and isolate those tools when you need them.

---

## Reflection (`reflect`) — use it sparingly

The `reflect` package inspects types and values at run time. `reflect.TypeOf` returns a `reflect.Type`. `reflect.ValueOf` returns a `reflect.Value`. `Kind` reports the underlying kind (`struct`, `slice`, `ptr`, and others).

`encoding/json` and similar packages use reflection. That use is a reason that JSON encode and decode are slower than hand-written field copies. Application code that uses `reflect` for business rules is often slower and harder to type-check.

Common operations:

- `Value.NumField` and `Field` walk exported struct fields.
- `Interface` returns the value as `any`.
- `Set` writes a value when `CanSet` is true.
- `Type.Implements` checks an interface.

`CanSet` is false for unexported fields and for values that are not addressable. `ValueOf(x)` on a non-pointer often is not settable. Pass a pointer when you must set fields.

Reflection skips compile-time checks. A typo in a field name becomes a run-time error. Tests must cover those paths.

Use reflection for:

- generic codecs and formatters
- tools that must accept unknown structs
- thin adapters around tags such as `json`

Do not use reflection for:

- everyday field access
- a type switch that you can write with interfaces
- performance-critical inner loops

Go 1.22 and later still treat `reflect` as a last-choice API. Generics remove some old uses of `any` plus reflect. Prefer generics or interfaces when they express the design.

`unsafe` plus `reflect` is a sharp pair. Do not mix them unless you follow the current `unsafe` rules and you have a review.

### Questions

#### Theoretical questions

1. What is the difference between `reflect.Type` and `reflect.Value`?
2. Why is `CanSet` often false for `ValueOf` on a non-pointer?
3. Why does reflection move errors from compile time to run time?
4. Which standard library packages justify the cost of reflection?
5. When do generics replace a use of `reflect`?

#### Easy practical tasks

1. Print `TypeOf` and `Kind` for an `int`, a `string`, and a struct.
2. Walk the exported fields of a struct with `NumField` and print each name.
3. Try to `Set` a field on a non-pointer value. Record the panic or the `CanSet` result.
4. Read `go doc reflect.Value.Set`. Write the preconditions in your own words.

#### Medium practical tasks

1. Write a function that copies matching exported fields from one struct to another with `reflect`. Test two structs that share three field names.
2. Compare a JSON marshal of a struct with a hand-written encode of the same fields. Measure with a benchmark.
3. Replace one `reflect` helper with a generic function. Show that the compiler catches a type error that `reflect` missed.

#### Advanced practical tasks

1. Implement a tiny required-tag checker: walk fields and fail when a tagged field is the zero value. Use `IsZero`. Document the limits.
2. Read `encoding/json` encode or decode code in the standard library. Write a one-page map of how it uses `reflect`. Do not copy long source into your notes.

---

## `go:generate` and AST tools

`go generate` runs commands that you write in special comments. The comment form is:

```text
//go:generate stringer -type=Color
```

The comment must start at the beginning of the line. Run `go generate ./...` from the module root. `go generate` is not part of `go build`. You run it when you change the input of the generator.

Common generators:

- `stringer` for `String` methods on integer types
- mock generators for interfaces
- Protocol Buffer compilers (see topic 22)

Commit the generated Go files when your team agrees to that rule. The CI must fail when generated files are stale. A common check is to run `go generate` and then `git diff --exit-code`.

The `go/parser`, `go/ast`, and `go/token` packages parse Go source into an abstract syntax tree (AST). `parser.ParseFile` returns an `*ast.File`. You walk nodes with `ast.Inspect` or `ast.Walk`.

`go/types` type-checks packages. `golang.org/x/tools/go/packages` loads packages with full module awareness. Use `go/packages` for tools that must follow imports. Use `go/parser` alone for a single file experiment.

AST tools are the right level for linters and generators. Do not parse Go with regular expressions. Regular expressions miss comments, build tags, and nested syntax.

Keep generators deterministic. The output must not include the local time or a random id. Stable output keeps reviews small.

`//go:generate` can call `go run` on a generator in your module. That pattern keeps the generator version in `go.mod`.

### Questions

#### Theoretical questions

1. Does `go build` run `go generate`? Explain.
2. Where must a `//go:generate` comment start?
3. Why is an AST better than a regular expression for Go source?
4. When do you need `go/packages` instead of `parser.ParseFile`?
5. Why must generator output be deterministic?

#### Easy practical tasks

1. Add a `//go:generate echo hello` comment. Run `go generate`. Show the printed line.
2. Install or run `stringer` on a small integer type. Show the generated `String` method.
3. Parse a file with `parser.ParseFile`. Print the package name from the AST.
4. Use `ast.Inspect` to print every function name in a file.

#### Medium practical tasks

1. Write a generator that reads a Go file and writes a list of exported names. Hook it with `//go:generate go run`.
2. Add a CI-style script that runs `go generate ./...` and fails when `git status` shows a change.
3. Load a package with `golang.org/x/tools/go/packages`. Print the import path and the number of syntax files.

#### Advanced practical tasks

1. Write a small analyzer that reports exported functions without a doc comment. Use `go/ast` and `go/parser`.
2. Compare `go/parser` on one file with `go/packages` on a module that has build tags. Document files that appear in only one approach.

---

## Cgo (costs and when to avoid)

Cgo lets a Go file import the virtual package `C` and call C code. You write C in comments above `import "C"`. The `cgo` tool generates glue.

Cgo has costs:

- Build time increases.
- Cross-compilation needs a C cross-compiler.
- Each call from Go to C has overhead. A tight loop of small C calls is often slower than Go.
- The Go garbage collector must track pointers that you pass. The cgo pointer rules are strict. You must not store a Go pointer in C memory in a way that hides it from the collector.
- `C.CString` allocates. You must free it with `C.free` when the C API does not take ownership.
- Some platforms and `scratch` images become harder.

Set `CGO_ENABLED=0` to force a pure Go build. The build fails if a package needs Cgo. That failure is useful. Many teams keep `CGO_ENABLED=0` in release pipelines.

Reasons to use Cgo:

- You must call a vendor C library that has no Go API.
- You must use an operating-system API that the standard library does not expose.

Reasons to avoid Cgo:

- You want simple `GOOS`/`GOARCH` builds.
- You want a static binary in `scratch`.
- You can use a pure Go library of similar quality.

The standard library uses Cgo for some network and user lookups on some platforms. Build tags `netgo` and `osusergo` request the pure Go implementations where they exist.

Do not use Cgo for speed without a benchmark that includes call overhead. Do not wrap every C function if you can batch work on the C side.

`CGO_CFLAGS` and `CGO_LDFLAGS` pass flags to the C compiler and linker. Pin toolchains. Document the C library version.

### Questions

#### Theoretical questions

1. What import path enables Cgo in a file?
2. Why does Cgo make cross-compilation harder?
3. What must you do after `C.CString` in many programs?
4. Why can a small C function in a loop be slower than Go?
5. What does `CGO_ENABLED=0` do in a release build?

#### Easy practical tasks

1. Run `go env CGO_ENABLED`. Write the value.
2. Build a normal program with `CGO_ENABLED=0`. Confirm that it builds.
3. Read [https://pkg.go.dev/cmd/cgo](https://pkg.go.dev/cmd/cgo). Write four pointer rules in short sentences.
4. List three libraries that you would call only through Cgo and three that you would replace with a Go module.

#### Medium practical tasks

1. Write a tiny Cgo program that calls `C.puts` or a small C add function. Measure ten thousand calls versus a Go add.
2. Cross-compile that Cgo program to another `GOOS`. Record the error or the extra compiler that you needed.
3. Build with `netgo` and `osusergo` tags on a program that uses `net` and `os/user`. Document the difference in `go build -x` output if you see one.

#### Advanced practical tasks

1. Wrap a small C library (or a single `.c` file in the repo). Document memory ownership and a `Close` function on the Go side.
2. Write a decision record: Cgo versus rewrite versus a separate C process. Include build, security, and operations.

---

## Assembly awareness (`go tool objdump`, `//go:nosplit`)

Go programs compile to machine code. You can inspect that code. You almost never write assembly for an application.

`go build -gcflags=-S` prints assembly during compile. `go tool objdump` disassembles a binary or a package archive. `go tool objdump -s FuncName binary` limits the output to one function.

The Go assembler uses a form that looks like Plan 9 assembly. It is not the same syntax as GNU `as` for every platform. Architecture files live under `runtime` and `internal/bytealg` in the standard library.

`//go:nosplit` is a compiler directive. It means the function must not need a stack split (a stack growth check) at entry. The runtime uses this directive in low-level code. If you add `//go:nosplit` to a function that needs more stack than the limit, the program can crash. Do not add this directive to application code.

Other directives that you may see:

- `//go:noescape` tells the compiler that pointer arguments do not escape through the call.
- `//go:linkname` links a local name to another symbol. This is not a public API. The Go team can break it.

Use object dumps to learn. Example questions: did the compiler inline a function, and did a bounds check remain?

Do not copy assembly from a blog into a service to "make it faster". Write clear Go. Use `go test -bench` and the compiler report `-gcflags=-m` first.

This section is awareness only. You do not need to write assembly to be an effective Go programmer.

### Questions

#### Theoretical questions

1. What does `go tool objdump` show?
2. What does `//go:nosplit` claim about a function?
3. Why must application code avoid `//go:nosplit`?
4. What is `//go:linkname` a sign of?
5. Which tool do you use before you even look at assembly when you care about speed?

#### Easy practical tasks

1. Build a small binary. Run `go tool objdump` and find `main.main`. Write five instruction lines.
2. Compile with `go build -gcflags=-S` and save the output for one function.
3. Search the standard library for `//go:nosplit`. Open one file. Write why that file is runtime code.
4. Run `go help build` and `go tool objdump -h`. Write the flags that you used.

#### Medium practical tasks

1. Compare `objdump` for a function before and after you mark another function to force no inline (`//go:noinline`). Describe the call.
2. Use `-gcflags=-m` and match one "inlining" line with the assembly.
3. Disassemble `bytes.IndexByte` or a similar function. Write whether you see a call into an architecture-specific body.

#### Advanced practical tasks

1. Read the Go assembler document on go.dev. Write a half-page glossary: `TEXT`, `NOSPLIT`, `ABIInternal` at a high level.
2. Explain a crash risk of `//go:nosplit` on a function that calls a large stack frame. Use documentation, not an experiment on a production host.

---

## Runtime internals: the G/M/P model

The Go scheduler is an M:N scheduler. Many goroutines (G) run on a smaller number of operating-system threads (M). A logical processor (P) holds the resources that a thread needs to run Go code.

- G is a goroutine. It has a stack and a state (running, runnable, waiting).
- M is a machine. It is an operating-system thread.
- P is a processor. The number of P values is `GOMAXPROCS`. A thread must hold a P to run Go code.

A P has a local run queue of goroutines. There is also a global run queue. Idle processors steal work from other processors. This is work stealing.

When a goroutine makes a blocking system call, the M can detach from its P. The P can move to another M so that other goroutines continue. When the system call returns, the M tries to obtain a P again.

The network poller waits for socket readiness. Goroutines that block on network I/O park. They do not need a P while they wait.

`GOMAXPROCS` defaults to the number of CPU cores. A higher value allows more Go code to run in parallel. It does not increase speed when the work is not parallel. A value that is too high can add overhead.

`runtime.Gosched` yields the CPU. Most programs do not need it. Channels, mutexes, and I/O already park.

You do not set G, M, or P in application code. You observe them with traces and `GODEBUG=schedtrace`. Use this model to understand why a blocked syscall does not freeze the whole process, and why a tight CPU loop without function calls can delay preemption on older Go versions. Current Go versions have asynchronous preemption (Go 1.14 and later).

Do not tune `GOMAXPROCS` from rumor. Measure.

### Questions

#### Theoretical questions

1. What do G, M, and P stand for?
2. What does `GOMAXPROCS` set?
3. What happens to P when an M blocks in a syscall?
4. What is work stealing in this model?
5. Why does a parked network goroutine not need a P while it waits?

#### Easy practical tasks

1. Print `runtime.GOMAXPROCS(0)` and `runtime.NumCPU()`. Write both numbers.
2. Run a CPU-bound loop with `GOMAXPROCS=1` and `GOMAXPROCS=4`. Compare wall time for four parallel loops.
3. Read a short official or Go blog description of G/M/P. Write five sentences in your own words.
4. Draw three boxes labeled G, M, and P. Draw who holds whom while a goroutine runs Go code.

#### Medium practical tasks

1. Start many goroutines that block on `time.Sleep`. Watch `NumGoroutine` and a `schedtrace` line. Relate waiting Gs to Ms.
2. Open a socket and block on read with a deadline. Explain where the goroutine parks.
3. Compare a program that only does CPU math with a program that only does HTTP client calls. How do M counts differ in `schedtrace`?

#### Advanced practical tasks

1. Read the scheduler file comments in `$GOROOT/src/runtime/proc.go` (header comments only). Write a one-page summary of run queues and handoff.
2. Use `go tool trace` to show a syscall that releases a P. Screenshot or describe the viewer rows. Relate the view to G/M/P.

---

## Finalizers, and why they are rarely the right tool

`runtime.SetFinalizer(obj, func)` registers a function that the garbage collector may call when `obj` becomes unreachable. The finalizer is not a destructor from C++.

Limits:

- The runtime does not guarantee that a finalizer runs at a given time.
- The runtime does not guarantee that a finalizer runs before process exit.
- A finalizer can "resurrect" an object when it stores the pointer again. That behavior is surprising.
- Finalizers run on a dedicated goroutine. They can deadlock if they take locks in a bad order.
- Cycles of objects with finalizers can delay or prevent collection.

Use explicit cleanup. Offer `Close() error`. Call `Close` with `defer` when you open a file, a connection, or a handle. Track ownership in the type that created the resource.

`runtime.KeepAlive(x)` keeps `x` reachable until that line. Use it when a finalizer or a cgo call must not run while you still use memory that `x` owns. This is a rare need.

Go 1.24 adds `runtime.AddCleanup` as a cleanup registration that avoids some finalizer resurrection problems. Prefer `AddCleanup` over `SetFinalizer` when you are on Go 1.24 or later and you still need GC-tied cleanup. Prefer `Close` over both.

Do not use a finalizer to flush a log or to commit a transaction. Those actions must run at a known time. Do not use a finalizer to unlock a mutex.

A finalizer is a last-chance net for a leak that you cannot yet find. A test that checks `Close` is a better net.

### Questions

#### Theoretical questions

1. Does `SetFinalizer` guarantee that the function runs before exit?
2. What does resurrection mean for a finalizer?
3. Why is `Close` plus `defer` the normal design?
4. What does `runtime.KeepAlive` prevent?
5. What did Go 1.24 add as a cleanup alternative to `SetFinalizer`?

#### Easy practical tasks

1. Read `go doc runtime.SetFinalizer`. Write five warnings from that documentation.
2. Write a type with `Close` that sets a flag. Test that `defer` calls `Close`.
3. Write one sentence that states when a finalizer might never run.
4. Search `SetFinalizer` in the standard library. Name two packages that use it and why they might.

#### Medium practical tasks

1. Register a finalizer that sets an atomic flag. Allocate, drop the pointer, call `runtime.GC` twice. Record whether the flag is set. Explain why a pass or fail is not a stable test.
2. Show `KeepAlive` in a tiny example with a comment that states what must stay alive.
3. Compare `SetFinalizer` and `AddCleanup` in the Go 1.24+ documentation (or release notes). Write a table of differences.

#### Advanced practical tasks

1. Design a file wrapper that must not leak descriptors. Use `Close` and a test that opens many files. Do not rely on a finalizer for the test to pass.
2. Find a public post-mortem where a finalizer caused a deadlock or a leak. Write six sentences on the failure and the fix.

---

## Plugin package limitations

The `plugin` package loads a shared object at run time on some Unix platforms. `plugin.Open` opens a `.so` file. `Lookup` finds an exported symbol.

Limits are severe:

- Windows is not supported.
- The plugin and the host must be built with the same Go toolchain and compatible flags.
- The plugin model does not match Go modules well. Dependency versions must line up.
- You cannot unload a plugin in a safe, general way.
- Crashes and `init` side effects occur in the same process. A bad plugin can take down the host.
- Cross-compilation of plugins is painful.

Because of those limits, most teams do not use `plugin` in production. Prefer these designs:

- Compile features into the binary with build tags.
- Run extensions as separate processes and use RPC or HTTP.
- Use a scripting host only when you accept that extra runtime.

If you still experiment with `plugin`, keep it off the critical path. Load only signed artifacts from a controlled path. Treat the plugin as part of the same trust domain as the host.

Read `go doc plugin` for the current platform list. Do not assume macOS, Linux, and FreeBSD all behave the same in every Go version.

### Questions

#### Theoretical questions

1. What does `plugin.Open` load?
2. Why must the plugin and the host use the same toolchain?
3. Why is unload a problem?
4. What happens to the host process when a plugin panics in `init`?
5. What process-isolation design replaces `plugin` for many teams?

#### Easy practical tasks

1. Read `go doc plugin`. Write the supported operating systems that the page lists.
2. Write four reasons not to use `plugin` in a service that you operate.
3. Sketch a side process that talks JSON over stdin and stdout. List three benefits versus `plugin`.
4. Find a build tag example in topic 9 or in go.dev. Write how a tag replaces a plugin for two features.

#### Medium practical tasks

1. On Linux, follow the official `plugin` example if your OS supports it. Record every build command. If your OS does not support it, record the error.
2. Change the Go patch version of the plugin or the host on purpose (if you can). Record the load error.
3. Compare `plugin` with `net/rpc` or gRPC for an extension. Write a table: isolation, versioning, deploy, latency.

#### Advanced practical tasks

1. Write an architecture note that forbids `plugin` in your team except for a named experiment. Include the isolation alternative.
2. Research one production outage or issue tracker thread about Go plugins. Summarize the version or platform problem in six sentences.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `reflect`, `go generate`, and generics each move work in time (compile time versus run time)?
2. Why do Cgo, `plugin`, and assembly all make builds and operations harder?
3. How does the G/M/P model explain a blocked syscall and a CPU-bound loop?
4. Why are finalizers a poor substitute for `Close`, and when is `AddCleanup` still not enough?
5. When is an AST-based tool the right way to change or check code?

#### Easy practical tasks

1. Write a cheat sheet: `reflect.TypeOf`, `//go:generate`, `CGO_ENABLED`, `go tool objdump`, G/M/P, `SetFinalizer`, `plugin.Open`.
2. Parse one of your files with `go/parser` and print function names. Then format the file with `gofmt`.
3. Build a binary with `CGO_ENABLED=0` and disassemble `main.main`. Write the first few instructions.
4. Draw G, M, and P during a blocking DNS lookup and during a tight math loop.

#### Medium practical tasks

1. Add a `stringer` generate step and a small `reflect` debug dump behind a flag. Document why the dump must not run on every request.
2. Write a decision table: problem, first choice, last choice. Include JSON tags, C library, extra feature flag, resource cleanup, and hot function.
3. Read header comments in `runtime/proc.go` and `cmd/cgo` documentation. Write ten facts that affect application authors.

#### Advanced practical tasks

1. Build a tiny linter with `go/packages` that forbids `runtime.SetFinalizer` and `plugin.Open` in `internal/`. Run it on a sample that fails and a sample that passes.
2. Write a two-page report: "Tools we allow in production code". Cover `reflect`, Cgo, assembly directives, finalizers, and plugins. Give a yes or no and one reason for each.
