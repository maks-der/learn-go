# 11. Testing and Tooling

## Description

Go tests live next to the code. The `testing` package and the `go test` command run tests, benchmarks, examples, and fuzz targets. Coverage, linters, and `pprof` complete the daily toolkit.

This topic teaches each tool as a separate skill. Complete topic 1 before this topic. You must know `go test`, packages, and modules. Topic 10 helps when a test starts goroutines or uses a context.

Use one term for each concept. A test function is not a benchmark. An example is not a fuzz target. A fake is not a mock.

---

## `*_test.go` and table-driven tests

A test file is a Go file whose name ends with `_test.go`. The `go test` command compiles these files with the package under test. The `go build` command ignores them when you build a library or a command.

Put the test file in the same directory as the package. The file can use one of two package names:

- `package foo` — the test is an internal test. The file can use unexported names.
- `package foo_test` — the test is an external test. The file imports `foo` as a client. The file can use only exported names.

Use an external test when you want to test the public API. Use an internal test when you need to check unexported helpers. Put fixtures in a folder named `testdata`. That folder is not a Go package for `go list`.

A test function has this form:

```go
func TestAdd(t *testing.T) {
	got := Add(1, 2)
	if got != 3 {
		t.Fatalf("Add(1, 2) = %d, want 3", got)
	}
}
```

The name must start with `Test`. The next rune must not be a lowercase letter. The only parameter is `*testing.T`.

The test fails when the code calls `t.Fail`, `t.Error`, `t.Errorf`, `t.Fatal`, or `t.Fatalf`. `Fatal` stops the current test function. `Error` records a failure and continues.

A table-driven test stores cases in a slice. One loop runs every case.

```go
func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{name: "positive", in: 2, want: 2},
		{name: "negative", in: -2, want: 2},
		{name: "zero", in: 0, want: 0},
	}
	for _, tc := range tests {
		got := Abs(tc.in)
		if got != tc.want {
			t.Fatalf("%s: Abs(%d) = %d, want %d", tc.name, tc.in, got, tc.want)
		}
	}
}
```

Give each case a `name`. Put the interesting cases in the table: zero, empty, error, and one normal value. Do not copy the same `if` block many times.

### Questions

#### Theoretical questions

1. What suffix marks a file as a test file?
2. What is the difference between `package foo` and `package foo_test` in a test file?
3. Does `go build` compile `_test.go` files into the production binary?
4. What must the name of a test function start with?
5. What is a table-driven test?

#### Easy practical tasks

1. Create a module with `add.go` and `add_test.go` in the same folder. Run `go test`.
2. Rename `add_test.go` to `add_checks.go`. Run `go test`. Write what happens.
3. Write `TestAbs` as a table with three cases. Run `go test -v`.
4. Use `t.Errorf` on a wrong want value. Confirm that the test fails. Then fix the want value.

#### Medium practical tasks

1. Write one internal test that calls an unexported helper. Write one external test that calls only the exported function. Run both with `go test`.
2. Add an error case to a table for a function that returns `(T, error)`. Check `err` in the loop.
3. Put a fixture in `testdata`. Read it from a test with `os.ReadFile`. Confirm `go list ./...` does not treat `testdata` as a package.

#### Advanced practical tasks

1. Split tests for one package across two `_test.go` files. Share a helper in a third `_test.go` file. Run `go test -v`.
2. Add a `//go:build` tag on a test file. Run `go test` with and without `-tags`. Show which tests run.

---

## Subtests and helpers

`t.Run` starts a subtest. Each subtest has a name. `go test` can run one subtest with `-run`.

```go
func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{name: "negative", in: -2, want: 2},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Abs(tc.in)
			if got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}
```

A failure in one subtest does not stop the other subtests. A `Fatal` in a subtest stops only that subtest.

Run one subtest:

```text
go test -run TestAbs/negative
```

A helper function calls `t.Helper()`. Then failure lines point to the test, not to the helper.

```go
func assertEq[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
```

`t.Cleanup` registers a function that runs after the test (or subtest) ends. Use `Cleanup` for temp files and for `cancel` of a context.

`t.Parallel` allows the test to run in parallel with other parallel tests. Call `t.Parallel` at the start of the test or subtest. Do not share writable state across parallel tests. Capture the table row in a new variable when you use older Go versions. Go 1.22 gives each loop iteration its own `tc`.

`t.TempDir` creates a temporary directory. The test framework removes it after the test.

### Questions

#### Theoretical questions

1. What does `t.Run` add that a plain loop does not add?
2. What does `t.Helper` change in a failure report?
3. When does `t.Cleanup` run?
4. What must you avoid when you call `t.Parallel`?
5. How do you run one named subtest from the command line?

#### Easy practical tasks

1. Convert a table loop to `t.Run` with the case name. Run `go test -v`.
2. Run `go test -run TestAbs/negative` for that test.
3. Write `assertEq` with `t.Helper`. Use it in two tests. Fail one test and read the file:line.
4. Create a file in `t.TempDir`. Read it back in the same test.

#### Medium practical tasks

1. Register `t.Cleanup` that removes a side effect you create (for example a package-level flag). Show that later tests see the restored value.
2. Mark two subtests `t.Parallel`. Share a mutex if they touch the same counter. Run `go test -count=20`.
3. Skip a subtest with `t.Skip` when an environment variable is missing. Run the parent test.

#### Advanced practical tasks

1. Build a helper that starts an `httptest.Server` and registers `server.Close` on `t.Cleanup`. Use it in two tests.
2. Compare `t.Fatal` in the parent test with `t.Fatal` in a subtest. Write six sentences on what continues and what stops.

---

## Benchmarks, examples, and fuzz tests

A benchmark function has the form `func BenchmarkXxx(b *testing.B)`. The body runs the work `b.N` times. The tool chooses `N`.

```go
func BenchmarkAbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Abs(-1)
	}
}
```

Run benchmarks with `go test -bench=. -benchmem`. Do not use `t` in a benchmark. Call `b.ResetTimer` after setup if setup is expensive. Stop the timer around setup with `b.StopTimer` and `b.StartTimer` when needed.

An example function has the form `func ExampleXxx()`. The comment after the function can contain `Output:`. `go test` checks that the printed text matches.

```go
func ExampleAbs() {
	fmt.Println(Abs(-2))
	// Output:
	// 2
}
```

Examples appear in `go doc` and on pkg.go.dev. Write examples for the public API.

A fuzz test has the form `func FuzzXxx(f *testing.F)`. You add seed values with `f.Add`. The fuzz function receives `*testing.T` and the typed inputs.

```go
func FuzzAbs(f *testing.F) {
	f.Add(-2)
	f.Fuzz(func(t *testing.T, n int) {
		got := Abs(n)
		if got < 0 {
			t.Fatalf("Abs(%d) = %d", n, got)
		}
	})
}
```

Run fuzzing with `go test -fuzz=FuzzAbs -fuzztime=5s`. The tool writes failing inputs under `testdata/fuzz`. Commit those files. They become regression seeds.

Do not fuzz a function that talks to the network without a bound. Keep the fuzz body fast and deterministic.

### Questions

#### Theoretical questions

1. Who chooses `b.N` in a benchmark?
2. What flag prints allocations in a benchmark?
3. What does the `Output:` comment in an example do?
4. What does `f.Add` provide to a fuzz test?
5. Where does the fuzz engine store a failing input?

#### Easy practical tasks

1. Write `BenchmarkAbs`. Run `go test -bench=Abs -benchmem`.
2. Write `ExampleAbs` with an `Output:` block. Run `go test`.
3. Change the example output on purpose. Record the test failure. Restore the output.
4. Write `FuzzAbs` with one seed. Run `go test -fuzz=FuzzAbs -fuzztime=3s`.

#### Medium practical tasks

1. Add setup before a benchmark. Use `ResetTimer` after you build a large slice.
2. Write an example for a function that prints two lines. Match both lines in `Output:`.
3. Fuzz `strconv.Atoi` on strings. Fail when a successful parse does not round-trip with `Itoa` for small integers. Record any crash or none.

#### Advanced practical tasks

1. Compare two implementations in two benchmarks. Write which is faster and how many allocations each uses.
2. Keep a failing fuzz seed in `testdata/fuzz`. Show that `go test` without `-fuzz` still fails on that seed. Then fix the bug.

---

## `httptest` and coverage

Package `net/http/httptest` builds an HTTP server or a recorder for tests.

`httptest.NewRecorder` implements `http.ResponseWriter` and stores the response.

```go
req := httptest.NewRequest(http.MethodGet, "/health", nil)
rec := httptest.NewRecorder()
Health(rec, req)
if rec.Code != http.StatusOK {
	t.Fatalf("status %d", rec.Code)
}
```

`httptest.NewServer` starts a local server. Close it with `server.Close` or `t.Cleanup(server.Close)`.

```go
srv := httptest.NewServer(http.HandlerFunc(Health))
defer srv.Close()
resp, err := http.Get(srv.URL + "/health")
```

Use `NewRequest` or `NewRequestWithContext` for the incoming request. Set the method, the path, and the body. You do not listen on a real port with a recorder.

Coverage measures which statements `go test` ran.

```text
go test -cover ./...
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out
```

`-cover` prints a percent. `-coverprofile` writes a file. `go tool cover -html` opens a report. Aim to cover the error paths, not only the happy path.

Coverage is not quality. A test can run a line and still not check the result. Write assertions.

Use `-coverpkg` when you need coverage of packages that the test imports but does not sit in. Start with the package under test.

### Questions

#### Theoretical questions

1. What does `httptest.NewRecorder` store?
2. Why do you close `httptest.NewServer`?
3. What does `go test -cover` print?
4. What does `go tool cover -html` show?
5. Why is 100% coverage not the same as a correct package?

#### Easy practical tasks

1. Write a handler that writes `ok`. Test it with `NewRequest` and `NewRecorder`. Check the status and the body.
2. Start `NewServer` with that handler. GET the URL. Check the status.
3. Run `go test -cover` on a package with two functions. Write the percent.
4. Write `cover.out` and open the HTML report. Name one line that is red.

#### Medium practical tasks

1. Test a POST handler that reads JSON. Send a body with `NewRequest` and `strings.NewReader`. Check the status for good JSON and bad JSON.
2. Raise coverage of an error branch that you missed. Show the percent before and after.
3. Use `t.Cleanup` to close the test server. Run the test twice with `-count=2`.

#### Advanced practical tasks

1. Test middleware: wrap a handler that panics or returns 401. Use a recorder. Check status and a header.
2. Produce a coverage profile for `./...` and find a package with 0%. Write whether that package needs tests or is a command with no logic.

---

## Linters and `pprof`

`go vet` is the first static check. Run `go vet ./...` on every change. Topic 1 introduces `vet`.

A linter finds more style and bug patterns. `staticcheck` and `golangci-lint` are common tools. Install them as separate commands. Run them in CI. Do not disable a check without a reason.

```text
go vet ./...
staticcheck ./...
```

`golangci-lint` runs many linters in one process. Pin the version in CI. Start with the default set. Add `errcheck`, `govet`, and `staticcheck` if you configure by hand.

`pprof` profiles CPU, heap, goroutines, and blocking. The `net/http/pprof` package mounts endpoints on a server. Package `runtime/pprof` writes profiles to a file.

```text
go test -cpuprofile=cpu.out -bench=Abs
go tool pprof cpu.out
```

In `pprof`, `top` shows the hottest functions. `list FuncName` shows lines. `web` opens a graph when Graphviz is installed.

Use profiles to explain a slowness that you measured. Do not optimize from a guess. Compare two profiles after a change.

The race detector is also a tool. Run `go test -race ./...` in CI. Topic 10 explains races.

`go generate` runs `//go:generate` comments. It does not run during `go build`. Commit generated files or generate them in CI. Document the command.

### Questions

#### Theoretical questions

1. What class of problems does `go vet` report?
2. Why do teams pin `golangci-lint` in CI?
3. What does `go tool pprof` read?
4. Does `go generate` run when you run `go build`?
5. When do you collect a CPU profile?

#### Easy practical tasks

1. Introduce a `Printf` format bug. Run `go vet`. Fix the verb.
2. Run `go test -bench=Abs -cpuprofile=cpu.out` on a benchmark. Start `go tool pprof cpu.out` and run `top`.
3. Read `go doc net/http/pprof`. Write two endpoint paths.
4. Write a `//go:generate` comment that runs `echo` or a small command. Run `go generate`.

#### Medium practical tasks

1. Install `staticcheck` or `golangci-lint`. Run it on a module. Fix one reported issue.
2. Write a program that allocates in a loop. Capture a heap profile. Name the top allocator.
3. Add a Makefile or a script that runs `fmt`, `vet`, `test -race`, and a linter.

#### Advanced practical tasks

1. Enable `pprof` on a tiny HTTP server. Collect `/debug/pprof/goroutine` while you leak one goroutine on purpose. Then fix the leak.
2. Compare `go test -cover` with a linter report on the same package. Write what each tool can see that the other cannot see.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a table-driven `TestXxx`, a subtest, and a helper work together in one file?
2. When do you write a benchmark, an example, and a fuzz test for the same function?
3. How do `httptest.NewRecorder` and `httptest.NewServer` differ as ways to test HTTP?
4. What do coverage, `vet`, and `pprof` each measure or inspect?
5. Why does CI run `-race`, `vet`, and tests together?

#### Easy practical tasks

1. Write `Abs` with a table test, one example, and `go test -cover`.
2. Add a subtest name for the zero case. Run only that subtest with `-run`.
3. Draw a table: Test, Benchmark, Example, Fuzz. Columns: function prefix, extra flag, typical assertion.
4. Save a coverage profile and a CPU profile for the same package. Write the two file names.

#### Medium practical tasks

1. Test an HTTP JSON handler with a recorder. Add a fuzz test on the decode helper. Run both.
2. Write a script that fails when coverage is below a number that you choose, then run `vet` and `-race`.
3. Add `t.Cleanup` and `t.TempDir` to a test that writes a file and starts a test server.

#### Advanced practical tasks

1. Build a small package with tests, one benchmark, one example, one fuzz target, and an HTTP handler test. Run a full local CI sequence.
2. Profile one slow test or benchmark. Change one line. Show `pprof` top before and after.
