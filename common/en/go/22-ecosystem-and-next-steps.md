# 22. Ecosystem and Next Steps

## Description

The Go language is small. The ecosystem is large. This topic shows how you read official style, how you read the standard library, and how you choose extra tools. Learn the standard library first. Add a library when it solves a clear problem.

After this topic, continue with a real program that you use. Read other people source. Contribute a small fix. Keep the habits from topics 1 to 21.

---

## Effective Go and Go Proverbs

Effective Go is the official style and design essay at [https://go.dev/doc/effective_go](https://go.dev/doc/effective_go). It is not the language specification. It teaches how to write clear Go. Read it after you know the basics. Read it again after you ship a module.

The Go Proverbs are short lines that Rob Pike presented in 2015. The list is at [https://go-proverbs.github.io/](https://go-proverbs.github.io/). Proverbs are reminders. They are not laws. Use them when they improve a design. Do not use them to stop a needed change.

Useful proverbs for daily work:

- "Don't communicate by sharing memory; share memory by communicating." Prefer channels when they make ownership clear. Use a mutex when a shared struct is simpler.
- "Concurrency is not parallelism." Many goroutines can run on one core.
- "The bigger the interface, the weaker the abstraction." Keep interfaces small.
- "Make the zero value useful." A zero `struct` must be safe when you can.
- "A little copying is better than a little dependency." Do not import a module for ten lines.
- "Clear is better than clever."
- "Errors are values." Handle them.
- "Gofmt's style is no one's favorite, yet gofmt is everyone's favorite."

Effective Go covers formatting, commentary, names, allocation with `make` and `new`, interfaces, and concurrency. Match those sections to topics that you already completed.

When a proverb conflicts with a measurement or a safety rule, keep the measurement and the safety rule. A proverb is a start, not a proof.

### Questions

#### Theoretical questions

1. How is Effective Go different from the language specification?
2. What does "the bigger the interface, the weaker the abstraction" mean for a Go interface?
3. What does "make the zero value useful" mean for a struct?
4. When may a mutex be clearer than a channel?
5. Why is a proverb not a proof?

#### Easy practical tasks

1. Open Effective Go. Write the section titles in order.
2. Copy five proverbs into a file. Write one sentence under each proverb that uses only your own words.
3. Find the Effective Go advice on names. Change one of your identifiers to match it.
4. Find the proverb about `gofmt`. Run `gofmt` on a file that you format by hand.

#### Medium practical tasks

1. Read Effective Go "Interfaces" and "Concurrency". Write six sentences that link those pages to your own code.
2. Review one of your packages for large interfaces. Split or shrink one interface. Record the before and after.
3. Find a dependency that you can replace with twenty lines of your own code. Measure maintenance cost in sentences, not only line count.

#### Advanced practical tasks

1. Write a one-page style guide for your team that cites Effective Go and three proverbs. Add two local rules that Effective Go does not cover.
2. Watch or read the Go Proverbs talk notes. Add historical context to three proverbs. State one case where a proverb would lead to a poor design.

---

## Reading standard library source

The standard library is the best set of Go examples that every toolchain includes. The source lives in `$GOROOT/src`. You can also browse [https://cs.opensource.google/go/go](https://cs.opensource.google/go/go) or [https://github.com/golang/go](https://github.com/golang/go).

Start with small packages: `strings`, `bytes`, `errors`, `io`. Then read `net/http` in pieces. Do not try to read `runtime` first.

How to read:

1. Open the package documentation on [https://pkg.go.dev](https://pkg.go.dev).
2. Open the same package in `$GOROOT/src`.
3. Read the package comment and the exported types.
4. Follow one function with `go doc -src`.
5. Note tests in `*_test.go`. Tests show intended use.

`go doc -src fmt.Println` prints the source of `Println`. `go env GOROOT` prints the root path.

The standard library prefers simple code over clever code. You will see manual loops, small interfaces, and careful error returns. You will also see low-level code in `runtime` and `sync`. That code is not a model for your HTTP handler.

Build tags and `internal` packages appear in the standard library. Use them as examples of the rules from earlier topics.

When you copy an idea, copy the idea. Do not copy large files. The license is BSD-style. Still, keep your own names and your own comments.

Read the `CONTRIBUTING` and proposal process only when you want to change the standard library. Most learners only read.

### Questions

#### Theoretical questions

1. Where does the standard library source live on your disk?
2. Why start with `strings` instead of `runtime`?
3. What extra information do `*_test.go` files give you?
4. Why is `runtime` a poor model for application code?
5. What command prints the source of a function in the terminal?

#### Easy practical tasks

1. Run `go env GOROOT`. Open `src/strings/strings.go`. Write the package comment in one sentence.
2. Run `go doc -src errors.Join`. Write how `Join` treats nil errors.
3. Open `io.Reader` in the source. List the interface and one type that implements it in that package.
4. Bookmark pkg.go.dev and the Google Git viewer for the Go repo.

#### Medium practical tasks

1. Read `encoding/json` decode of a struct at a high level. Write the steps from `Decode` to field set. Use function names.
2. Compare `log/slog` handler source with how you configured `slog` in topic 18. Write four matches.
3. Read one `net/http` test that uses `httptest`. Write the pattern that you will reuse.

#### Advanced practical tasks

1. Pick a bug that you had in your code. Find a standard library function that solves it. Read the source and the tests. Write what you missed.
2. Trace `http.Server.Serve` down to the accept loop. Write a diagram of five functions. Stop before the full multiplexer internals if time is short.

---

## Popular HTTP libraries (`chi`, `echo`, `gin`)

`net/http` is enough for many APIs, especially with the Go 1.22 `ServeMux` patterns. Learn it first. Add a router or a framework when you need extra features and your team already knows that library.

[github.com/go-chi/chi](https://github.com/go-chi/chi) is a lightweight router. It builds on `net/http`. Handlers stay `http.Handler`. Middleware is `func(http.Handler) http.Handler`. Chi is a small step from the standard library.

[github.com/labstack/echo](https://github.com/labstack/echo) uses its own `echo.Context`. Binding, rendering, and middleware are built in. You leave the plain `http.Handler` signature.

[github.com/gin-gonic/gin](https://github.com/gin-gonic/gin) also uses its own context type. Gin is fast in common benchmarks. It includes binding and a large middleware set. You depend on Gin types in every handler.

Trade-offs:

- Standard library: fewer dependencies, more code you own, easy to hire for.
- Chi: routing and middleware with standard handlers.
- Echo and Gin: faster to start a feature, harder to swap, more automatic binding that hides types.

Do not add all three. Do not wrap Gin in a hexagon with three extra layers in the first version. If you pick a framework, keep handlers thin and keep domain code free of framework types.

Security and updates matter. Pin versions. Read release notes. Replace abandoned modules.

Benchmark your own handlers before you switch for speed. Most APIs are slow because of I/O, not because of the router.

### Questions

#### Theoretical questions

1. Why must you learn `net/http` before Gin or Echo?
2. How does Chi stay close to the standard library?
3. What is the cost of a custom context type on every handler?
4. Why is a router benchmark often the wrong reason to switch?
5. Why must domain packages avoid Gin or Echo types?

#### Easy practical tasks

1. Write the same `GET /health` handler with `ServeMux` and with Chi. Compare the `main` functions.
2. Read the Chi, Echo, and Gin README files. Write one sentence for the target user of each project.
3. List three features that Gin or Echo include that `net/http` does not include.
4. Print the module path and latest version that `go list -m` shows after you add one library.

#### Medium practical tasks

1. Port a three-route JSON API from `ServeMux` to Chi. Keep tests against `http.Handler`.
2. Port the same API to Gin or Echo. Rewrite tests. Count lines that now import the framework.
3. Add the same logging middleware idea in `net/http` and in Chi. Compare signatures.

#### Advanced practical tasks

1. Write a decision record: stay on `net/http`, adopt Chi, or adopt Gin/Echo. Include team skill, test style, and swap cost.
2. Measure request latency for a JSON echo handler on two stacks under the same load tool. Report whether the router is visible in the numbers.

---

## gRPC and Protocol Buffers

gRPC is an RPC framework that uses HTTP/2 and Protocol Buffers. You define a service in a `.proto` file. A compiler generates Go types and client and server stubs.

Protocol Buffers (protobuf) is a binary schema format. `proto3` is the common syntax. Fields have numbers. Those numbers are the stable contract. Do not reuse a field number for a new meaning.

The modern Go toolchain uses `protoc-gen-go` and `protoc-gen-go-grpc`, or the Buf toolchain (`buf.build`). Generated code lives in your module. Commit it if your team requires generated files in git.

gRPC gives:

- typed methods
- deadlines through `context.Context`
- metadata headers
- streaming RPCs (unary, server stream, client stream, bidi)

gRPC needs HTTP/2. Browsers do not call gRPC in the same way that they call JSON REST. Teams often add gRPC-Gateway or a separate JSON API for browsers.

Compatibility rules match topic 21. Add fields. Do not change field numbers. Reserve removed numbers.

Errors use gRPC status codes, not HTTP status codes. Map them with care at a gateway.

Use gRPC when many internal services need a strict schema and streaming. Use JSON HTTP when browsers and simple clients matter more. Both can exist in one company.

TLS and authentication are still required in production. gRPC does not replace topic 18 authentication concepts.

### Questions

#### Theoretical questions

1. What file defines a gRPC service?
2. Why are protobuf field numbers part of the contract?
3. How do deadlines travel in a gRPC call?
4. Why do browsers often need a gateway in front of gRPC?
5. What are the four streaming modes at a high level?

#### Easy practical tasks

1. Write a tiny `.proto` with a `Greeter` service and a `SayHello` method. Do not generate yet if tools are missing. Label field numbers.
2. Read the official Go gRPC quick start. List the plugins that you must install.
3. Compare a JSON field name with a protobuf field number. Write three sentences.
4. List five gRPC status codes and a meaning for each code.

#### Medium practical tasks

1. Generate Go code with `protoc` or Buf for your `Greeter`. Implement the server. Call it from a Go client.
2. Set a 50 ms deadline on the client. Make the server sleep 100 ms. Record the status.
3. Add a field to the message. Show that an old client can talk to a new server if you only add a field.

#### Advanced practical tasks

1. Add TLS to a local gRPC server and client. Document certificate files and the client dial options.
2. Design a migration from JSON HTTP to gRPC for one internal method. Include schema, errors, and a browser path.

---

## GraphQL (optional)

GraphQL is a query language for APIs. The client asks for a graph of fields in one request. The server returns only those fields. GraphQL is optional. Many Go services never need it.

A schema defines types and fields. A resolver fills one field. A Go server often uses [gqlgen](https://github.com/99designs/gqlgen) or a similar library. gqlgen generates types from the schema.

Benefits:

- Clients can fetch nested data in one round trip.
- Mobile apps can avoid unused fields.

Costs:

- The server must bound query depth and cost. A nested query can be expensive.
- The N+1 problem appears when a resolver hits the store once per item. Use a data loader or a batched store call.
- Caching is harder than for simple HTTP GET URLs.
- File upload and auth still need a clear design.

Use GraphQL when many different clients need different shapes of the same graph. Use REST or gRPC when the shapes are few and stable.

Do not expose a GraphQL playground on the public internet in production. Apply the same auth and rate limits as any other API.

Keep business rules out of generated resolver stubs. Call your domain functions.

### Questions

#### Theoretical questions

1. Who chooses the fields in a GraphQL response?
2. What is a resolver?
3. What is the N+1 problem in GraphQL?
4. Why is HTTP GET caching easier for a REST URL than for a GraphQL POST?
5. When is REST a better fit than GraphQL?

#### Easy practical tasks

1. Write a GraphQL query on paper that asks for a user name and two order totals.
2. Read the gqlgen getting-started page. Write the role of `schema.graphqls`.
3. List three denial-of-service risks for a public GraphQL endpoint (depth, breadth, aliases).
4. Compare one REST resource and one GraphQL type for the same "item".

#### Medium practical tasks

1. Build a tiny gqlgen server with one `hello` query. Call it with a POST body.
2. Add a nested field that would N+1. Log each store call. Then batch the store calls and compare the log count.
3. Add a query complexity limit or a depth limit. Prove that a deep query fails.

#### Advanced practical tasks

1. Write a one-page "GraphQL or not" decision for a product with a web app and a public API.
2. Design auth for GraphQL: identity on the request context, per-field checks, and rate limits. Keep secrets out of the schema.

---

## Kubernetes operators and `client-go` (optional, specialized)

Kubernetes is a container orchestrator. An operator is a controller that manages a custom resource. You write a loop: read desired state, read live state, then act to close the gap. That loop is reconcile.

`client-go` is the official Go client for the Kubernetes API. Controller runtimes such as [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) sit on top of `client-go`. [Operator SDK](https://sdk.operatorframework.io/) and [Kubebuilder](https://book.kubebuilder.io/) generate project layouts.

This path is specialized. You need Kubernetes objects, RBAC, and cluster access. You do not need it to write a normal HTTP service.

If you learn it:

- Define a Custom Resource Definition (CRD).
- Write a reconcile function that is idempotent.
- Use `context.Context` and client rate limits.
- Never store cluster-admin credentials in an application that does not need them.
- Test with envtest or a kind cluster.

`client-go` informers cache objects. Direct API calls without a cache can overload the API server. Follow the current controller-runtime examples.

Do not run an operator from a laptop against a production cluster while you learn. Use a local cluster.

Most Go developers never write an operator. Learn `net/http` and modules first.

### Questions

#### Theoretical questions

1. What does a reconcile loop compare?
2. What is a CRD?
3. Why must reconcile be idempotent?
4. Why can naive `client-go` calls overload the API server?
5. Why is this topic optional for a general Go developer?

#### Easy practical tasks

1. Read the Kubebuilder or Operator SDK "what is an operator" page. Write five sentences.
2. List four Kubernetes objects that a normal service uses (Deployment, Service, and two more).
3. Read `go doc` is not enough here; open the `client-go` repository README. Write the purpose of informers in one paragraph.
4. Write three reasons not to start your Go career with operators.

#### Medium practical tasks

1. Install a local cluster (kind or minikube) if your environment allows it. Deploy a public sample operator or a sample CRD. Describe one reconcile.
2. Sketch RBAC for an operator that may only read ConfigMaps and write one CRD. No cluster-admin.
3. Draw the reconcile flow: event, queue, handler, client, status update.

#### Advanced practical tasks

1. Build a tiny controller with controller-runtime that sets a label on a ConfigMap. Test with envtest or kind.
2. Write an operations note: how a bad reconcile can overload the API server, and which client settings reduce that risk.

---

## WebAssembly (`GOOS=js` `GOARCH=wasm`)

Go can compile to WebAssembly (WASM) for a JavaScript host. Set the target:

```text
GOOS=js GOARCH=wasm go build -o app.wasm .
```

The toolchain ships `wasm_exec.js`. On Go 1.23 and earlier, the file is `$GOROOT/misc/wasm/wasm_exec.js`. On Go 1.24 and later, the file is `$GOROOT/lib/wasm/wasm_exec.js`. The browser or Node.js loads that script and the `app.wasm` file. Use the file from the same Go version that built the WASM binary.

The `syscall/js` package talks to JavaScript values. You pass functions to JS. You convert types with care. Memory stays in the WASM linear memory.

Limits:

- The binary is large compared with a hand-written JS module.
- Not every standard library feature fits a browser (no arbitrary TCP listen).
- `syscall/js` code is not reusable on the server without build tags.

Use WASM when you want to reuse a Go library in the browser or in a JS tool. Do not expect a full HTTP server to run in the browser in the same way as on Linux.

TinyGo (next section) can produce smaller WASM for some programs. The standard Go toolchain supports more of the language.

Test in one browser and in Node if you support both. Automate `GOOS=js GOARCH=wasm go test` where tests do not need a real DOM.

Keep JS glue small. Put rules in Go. Put DOM calls behind a thin adapter.

### Questions

#### Theoretical questions

1. Which `GOOS` and `GOARCH` values target JavaScript WASM?
2. What role does `wasm_exec.js` play?
3. Why is a Go WASM file often large?
4. Why can you not listen on an arbitrary TCP port in a browser?
5. What is `syscall/js` for?

#### Easy practical tasks

1. Cross-compile a `fmt.Println("ok")` program to `app.wasm`. Write the file size.
2. Find `wasm_exec.js` in your `GOROOT`. Write the full path.
3. Read the official Go WASM wiki or go.dev blog post. Write four limits.
4. List three program types that fit Go WASM and three that do not.

#### Medium practical tasks

1. Follow the official example to run a Go WASM file in the browser or in Node. Record the commands.
2. Expose a Go function to JS with `syscall/js`. Call it from a small JS file. Return a number.
3. Use a build tag so that `syscall/js` files are not part of the Linux build.

#### Advanced practical tasks

1. Port a pure function (validate an id, hash a string with `crypto`) to WASM. Compare size with TinyGo if you install it.
2. Document a release: Go version, `GOOS`, `GOARCH`, `wasm_exec.js` version match, and browser test notes.

---

## TinyGo (embedded)

TinyGo is a separate compiler for a subset of Go. The site is [https://tinygo.org](https://tinygo.org). TinyGo targets microcontrollers and small WASM binaries. It is not a replacement for the official toolchain on servers.

Differences that you must expect:

- Not every package in the standard library works.
- Reflection and full `net/http` are limited or absent on many targets.
- Goroutines exist but the scheduler and memory model details differ from the GC runtime on `gc`.
- You flash a board with `tinygo flash` and a target name.

Use the official Go toolchain for services, CLIs, and most WASM that needs the full standard library. Use TinyGo when the device memory is small or when you need a small WASM module and your code stays in the supported subset.

Do not assume that a module that builds with `go build` also builds with `tinygo`. Test on the target.

Hardware work needs a data sheet, a probe, and care with voltage. Follow board documentation. Do not power a board from a random source.

TinyGo has its own runtime and its own issue tracker. File bugs there, not on github.com/golang/go, when the problem is TinyGo-only.

### Questions

#### Theoretical questions

1. What problem does TinyGo solve that the official `gc` toolchain does not solve well?
2. Why can a normal module fail to build with TinyGo?
3. Which program class still belongs to official Go?
4. Where do you report a TinyGo-only compiler bug?
5. Why must you test on the real target?

#### Easy practical tasks

1. Open tinygo.org. Write three supported boards or target families.
2. Install TinyGo if your environment allows it. Run `tinygo version`. Save the output.
3. Compile a `Blink` or `hello` example from the TinyGo docs for the `wasm` target or a simulator. Write the command.
4. List five standard library packages that you would not use on a microcontroller.

#### Medium practical tasks

1. Build the same small function with `go` and with `tinygo` for WASM. Compare binary sizes.
2. Read TinyGo language support notes. Write a table of supported and unsupported features that affect your code.
3. Write a program that uses only `machine` GPIO (or a mock) and no `net`. Document the target.

#### Advanced practical tasks

1. Flash a supported board with a sensor read or an LED blink. Record wiring and the `tinygo flash` command. Follow the board safety notes.
2. Write a two-page guide: when your team uses official Go, when it uses TinyGo, and how you share a pure algorithm package between them.

---

## Contributing to open source Go projects

Contribution is a skill. Start small. Read the project `CONTRIBUTING` file and the code of conduct. Follow the license.

For the Go project itself, the process uses Gerrit and a Contributor License Agreement (CLA). The guide is [https://go.dev/doc/contribute](https://go.dev/doc/contribute). You file an issue first for most changes. You do not start with a large rewrite.

For most GitHub Go projects:

1. Search issues. Comment that you want to work on one issue.
2. Fork and clone. Create a branch.
3. Add a test. Change the code. Run `go test ./...`, `go fmt`, and the project linters.
4. Open a pull request with a short why. Link the issue.
5. Answer review comments with new commits or as the project asks.

Good first contributions:

- documentation typos
- extra tests
- small bug fixes with a reproduce case

Bad first contributions:

- a reformat of the whole tree
- a new feature without an issue
- a dependency change that you did not discuss

Be polite. Maintainers have little time. A small patch with a test is easier to accept than a large idea.

You must not commit secrets. You must not add malicious code. You must respect export control and license rules.

If a project is unused or has no license, do not rely on it. Prefer maintained modules.

After you contribute, maintain your own modules with the same respect: tests, `go.mod`, release tags, and a changelog.

### Questions

#### Theoretical questions

1. Why do many Go project guides ask for an issue before a large change?
2. How does the Go project review process differ from a typical GitHub-only project?
3. Why is a full-tree reformat a poor first pull request?
4. What files must you read before you change code?
5. Why does a test help the reviewer?

#### Easy practical tasks

1. Open go.dev/doc/contribute. Write the role of the CLA and of Gerrit in two sentences each.
2. Find a `good first issue` on a Go project that you use. Write the issue number and the asked change. Do not claim it unless you will do the work.
3. Read `CONTRIBUTING.md` of that project. Write the test command that they require.
4. Write a pull request template in your notes: summary, test plan, linked issue.

#### Medium practical tasks

1. Contribute a documentation fix or a test to a project that you use. Follow their process. Record the URL if the work is public.
2. Run the project linters and tests on your machine before you push. Save the commands.
3. Review someone else public pull request (read only). Write three review comments that you would send. Be specific and kind.

#### Advanced practical tasks

1. Propose a small feature through an issue. Wait for agreement. Implement it with tests and docs. Describe the review rounds.
2. Publish a tiny useful module with a license, tests, `golangci-lint`, and a tagged `v0.1.0`. Write how you will respond to issues.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do Effective Go, the standard library source, and the Go Proverbs work together as a learning loop?
2. When do you stay with `net/http`, when do you add Chi, and when do you add gRPC or GraphQL?
3. Why are Kubernetes operators, TinyGo, and WASM specialized compared with a normal Linux service?
4. What compatibility idea is shared by JSON APIs, protobuf field numbers, and generated gRPC code?
5. What makes a first open source contribution easy to accept?

#### Easy practical tasks

1. Write a personal next-steps list of ten items. Map each item to a topic in this handbook or to an official document.
2. Read one Effective Go section and one standard library file on the same day. Write five facts.
3. Create a table: tool, official or third-party, when you allow it (`chi`, gqlgen, `client-go`, TinyGo, `syscall/js`).
4. Run `go doc -src` on a function that you use every week. Write what you learned.

#### Medium practical tasks

1. Rebuild a small JSON API from topic 18. Then write a short note on what would change if you used Chi, gRPC, or GraphQL instead.
2. Compile a WASM binary and a normal Linux binary of the same pure function. Compare size and test method.
3. Clone a popular Go library. Find `CONTRIBUTING`, the module path, and one package that you can read in one hour. Write a summary.

#### Advanced practical tasks

1. Deliver a capstone: a module that follows Effective Go, uses only justified dependencies, has tests, and includes a contribution-style README. Add one optional extra (gRPC or WASM), not both, unless you have time.
2. Write a six-month learning plan: standard library reading, one library deep dive, one community contribution, and one production concern from topic 21.
