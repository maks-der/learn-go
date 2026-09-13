# 13. Testing and Tooling

## Description

Go tests live next to the code. The `testing` package and the `go test` command run tests, benchmarks, examples, and fuzz targets. Coverage, linters, and `pprof` complete the daily toolkit.

This topic teaches each tool as a separate skill. Complete Topic 1 before this topic. You must know `go test`, packages, and modules. Topic 11 and Topic 12 help when a test starts goroutines or uses a context.

Use one term for each concept. A test function is not a benchmark. An example is not a fuzz target. A fake is not a mock. `go generate` does not run when you run `go build`.

---

## `*_test.go` files

A test file is a Go file whose name ends with `_test.go`. The `go test` command compiles these files with the package under test. The `go build` command ignores them when you build a library or a command.

Put the test file in the same directory as the package. The file can use one of two package names:

- `package foo` — the test is an internal test. The file can use unexported names.
- `package foo_test` — the test is an external test. The file imports `foo` as a client. The file can use only exported names.

Use an external test when you want to test the public API. Use an external test when an internal test would create an import cycle. Use an internal test when you need to check unexported helpers.

Rules for test files:

- The file name must end with `_test.go`.
- The compiler ignores a test file during `go build` of the non-test package.
- A folder named `testdata` is not a Go package for `go list`. Put fixtures in `testdata`.
- Build tags on test files work the same as on other files.

Do not put production code in a `_test.go` file. Do not put tests in a file that lacks the `_test.go` suffix. `go test` will not run those tests.

`go test` compiles the package plus the test files into a test binary. Then it runs that binary. `go test ./...` runs tests in every package under the current directory.

### Questions

#### Theoretical questions

1. What suffix marks a file as a test file?
2. What is the difference between `package foo` and `package foo_test` in a test file?
3. Does `go build` compile `_test.go` files into the production binary?
4. What is the `testdata` directory for?
5. When do you choose an external test package?

#### Easy practical tasks

1. Create a module with `add.go` and `add_test.go` in the same folder. Use `package` names that match. Run `go test`.
2. Rename `add_test.go` to `add_checks.go`. Run `go test`. Write what happens.
3. Add a `testdata/sample.txt` file. Confirm that `go list ./...` does not treat `testdata` as a package.
4. Write four sentences that compare an internal test and an external test.

#### Medium practical tasks

1. Write one internal test that calls an unexported helper. Write one external test that calls only the exported function. Run both with `go test`.
2. Create an import cycle on purpose with an internal test that imports a client package that imports the package under test. Record the error. Switch the test to `package foo_test` if that removes the cycle.
3. Run `go test -c` and find the test binary. Write the flag purpose from `go help test`.

#### Advanced practical tasks

1. Split tests for one package across two `_test.go` files. Share a helper in a third `_test.go` file. Run `go test -v`.
2. Add a `//go:build` tag on a test file. Run `go test` with and without `-tags`. Show which tests run.

---

## `TestXxx(t *testing.T)`

A test function has this form:

```go
func TestAdd(t *testing.T) {
    got := Add(1, 2)
    if got != 3 {
        t.Fatalf("Add(1, 2) = %d, want 3", got)
    }
}
```

The name must start with `Test`. The next rune must not be a lowercase letter. `Testadd` is not a test. `TestAdd` is a test. The only parameter is `*testing.T`.

`go test` runs each `TestXxx` function. The test fails when the code calls `t.Fail`, `t.Error`, `t.Errorf`, `t.Fatal`, `t.Fatalf`, or `t.FailNow`. The test is skipped when the code calls `t.Skip`, `t.Skipf`, or `t.SkipNow`.

`t.Fatal` and `t.Fatalf` stop that test function. `t.Error` and `t.Errorf` record a failure and continue. Use `Fatal` when later lines are not safe. Use `Error` when you want more checks in the same function.

Other useful methods:

- `t.Log` and `t.Logf` write logs. `go test` shows them with `-v` or when the test fails.
- `t.Cleanup(fn)` runs `fn` after the test (and after subtests) in last-in first-out order.
- `t.TempDir()` creates a temporary directory. The test removes it at the end.
- `t.Setenv` sets an environment variable for the test.
- `t.Parallel()` marks the test as parallel with other parallel tests.
- From Go 1.24, `t.Context()` returns a context that is cancelled when the test finishes. From Go 1.24, `t.Chdir` changes the working directory for the test.

Do not call `t.Parallel` if the test shares unsynchronized package state with other tests. Run `go test -race` when tests run in parallel.

`go test` flags that you use often: `-v`, `-count=1` (disable result cache), `-run TestAdd`, `-timeout`, `-race`.

### Questions

#### Theoretical questions

1. Which names does `go test` treat as test functions?
2. What is the difference between `t.Fatalf` and `t.Errorf`?
3. When does `go test` print `t.Log` lines?
4. What does `t.Cleanup` do?
5. Why can `t.Parallel` hide or reveal a data race?

#### Easy practical tasks

1. Write `TestAdd` that checks one sum. Run `go test -v`.
2. Make the test fail with `t.Errorf`. Run `go test`. Save the output.
3. Add `t.Log` of the inputs. Run with and without `-v`. Write the difference.
4. Run `go test -run TestDoesNotExist`. Write the result.

#### Medium practical tasks

1. Write a test that creates a file in `t.TempDir()`, writes bytes, and reads them back.
2. Use `t.Cleanup` to close a file. Confirm that the cleanup runs when the test fails.
3. Mark two tests with `t.Parallel`. Run `go test -v` and observe the order. Write what you see.

#### Advanced practical tasks

1. If your toolchain is Go 1.24 or later, use `t.Context()` with a helper that waits on `ctx.Done()`. Show that the context is cancelled at the end of the test. If your toolchain is older, use `t.Cleanup` plus `WithCancel` and write the same idea.
2. Write a flaky test that depends on another test in the same package through a package-level variable. Fix the tests so that `go test -count=20` is stable.

---

## Table-driven tests

A table-driven test stores cases in a slice. The test loops over the slice and runs the same checks.

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {name: "zeros", a: 0, b: 0, want: 0},
        {name: "positive", a: 2, b: 3, want: 5},
        {name: "negative", a: -1, b: 1, want: 0},
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := Add(tc.a, tc.b)
            if got != tc.want {
                t.Fatalf("Add(%d, %d) = %d, want %d", tc.a, tc.b, got, tc.want)
            }
        })
    }
}
```

Each case has a `name`. The name appears in the `go test` output. Use `t.Run` so that you can run one case with `-run TestAdd/positive`.

From Go 1.22, each loop iteration has its own `tc`. The closure can use `tc` without a manual copy. Before Go 1.22, you wrote `tc := tc` inside the loop.

Put invalid inputs and error cases in the same table when the API returns an error. Add a `wantErr` field. Do not hide a case in a comment. Add a row.

Keep the case struct small. Share setup that is the same for every row outside the loop. Put setup that is unique in the row.

Do not generate random rows as the only tests. Random rows belong in a fuzz target. The table is the list of cases that you understand.

### Questions

#### Theoretical questions

1. What is a table-driven test?
2. Why does each row need a `name`?
3. How do you run one row with `go test`?
4. What did Go 1.22 change for loop variables in these tests?
5. When do you add a `wantErr` field to the case struct?

#### Easy practical tasks

1. Write a table-driven test for a function that returns the larger of two integers. Add four rows.
2. Run `go test -v`. Confirm that each row name appears.
3. Run `go test -run 'TestMax/zeros' -v` (use a name that exists). Confirm that only that row runs.
4. Add a row that fails. Read the output. Write which row failed.

#### Medium practical tasks

1. Test a parser that returns `(value, error)`. Use `want` and `wantErr` in the table. Cover success and failure.
2. Move common setup out of the loop. Keep one field that changes per row. Run the test.
3. Convert three separate `TestXxx` functions into one table. Keep the same checks.

#### Advanced practical tasks

1. Build a table for a function with a slice input. Include nil, empty, and long slices. Use `cmp` or a careful `reflect.DeepEqual` only if you must. Prefer a clear comparison.
2. Add a `t.Parallel` call inside each `t.Run`. Run `go test -race -count=10`. Confirm that the table is safe.

---

## Subtests: `t.Run`

`t.Run` starts a subtest. The parent test waits for the subtest. The first argument is the name. The second argument is a function.

```go
func TestSplit(t *testing.T) {
    t.Run("empty", func(t *testing.T) {
        // ...
    })
    t.Run("comma", func(t *testing.T) {
        // ...
    })
}
```

The full name is `TestSplit/empty`. Spaces in the name become underscores in the run filter. Keep names short and unique.

`t.Run` returns `false` when the subtest failed or the parent must stop. You can use the return value. Most tests ignore it.

A failure in a subtest fails the parent. `t.Fatalf` in a subtest stops that subtest. It does not stop other subtests that the parent still starts. If you need to stop the parent, check the return value of `t.Run` or call `t.FailNow` in the parent.

You can nest `t.Run`. Group cases by input class. Do not nest without a reason.

`t.Parallel` in a subtest runs that subtest in parallel with other parallel subtests of the same parent. The parent does not finish until those subtests finish. From Go 1.22, loop variables are safe in parallel subtests that close over `tc`.

Use subtests for table rows. Use subtests for two implementations of the same interface. Use subtests when setup is expensive and you want to share it in the parent (careful with mutation).

### Questions

#### Theoretical questions

1. What is the full name of a subtest?
2. Does `t.Fatalf` in a subtest stop the parent from starting the next subtest?
3. When does a parent test finish if a subtest called `t.Parallel`?
4. Why do you nest `t.Run`?
5. How does `t.Run` work with a table-driven test?

#### Easy practical tasks

1. Write `TestMath` with two subtests: `add` and `sub`. Run `go test -v`.
2. Run only one subtest with `-run`. Write the command and the output.
3. Fail one subtest. Confirm that the other subtest still runs.
4. Write the full names of three nested subtests in a tree: parent, child, grandchild.

#### Medium practical tasks

1. Share a setup value in the parent. Use it in two subtests. Do not mutate it in a way that breaks the second subtest.
2. Mark subtests as parallel. Give each subtest its own copy of data. Run `go test -race`.
3. Use the `bool` return of `t.Run` to skip later subtests when setup fails.

#### Advanced practical tasks

1. Write a suite that tests two implementations of one interface with the same subtest table. Use an outer `t.Run` per implementation.
2. Measure `go test` time with sequential subtests and with parallel subtests for a slow fake I/O case. Write the two times.

---

## `t.Helper()`, `t.Fatalf`, `t.Error`

`t.Helper` marks the calling function as a helper. The next failure line in the `go test` output is the caller, not the helper.

```go
func assertEq[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Fatalf("got %v, want %v", got, want)
    }
}
```

Call `t.Helper` at the start of every helper that can fail. If you omit `t.Helper`, the failure points at the helper file. You then waste time in the wrong place.

`t.Fatalf(format, args...)` logs a message and stops the test function (or the subtest). It is `Logf` plus `FailNow`.

`t.Error(args...)` logs a message and marks the test as failed. The test function continues. `t.Errorf` is the format form.

Choose this way:

- Use `t.Fatalf` when the test cannot continue (nil pointer, setup failed, wrong type).
- Use `t.Error` or `t.Errorf` when later checks still have meaning.
- Use a helper with `t.Helper` when more than one test needs the same check.

Do not panic in a test to signal failure. Use `t.Fatal`. A panic fails the test, but the output is harder to read.

`t.Fail` marks failure and continues without a message. `t.FailNow` stops without a new message. Prefer `Error` and `Fatal` so that the output has a reason.

`t.Fatal` does not run the rest of the function. `defer` and `t.Cleanup` still run.

### Questions

#### Theoretical questions

1. What does `t.Helper` change in the test output?
2. What is the difference between `t.Error` and `t.Fatalf`?
3. When must a test stop at the first failure?
4. Does `t.Fatal` still run `defer` and `t.Cleanup`?
5. Why is a panic a poor way to fail a test?

#### Easy practical tasks

1. Write a helper without `t.Helper`. Fail a test. Write the file and line in the output.
2. Add `t.Helper` to the helper. Fail again. Write the new file and line.
3. Use `t.Errorf` twice in one test. Confirm that both messages appear.
4. Use `t.Fatalf` for the first check. Confirm that the second check does not run.

#### Medium practical tasks

1. Write `assertEq` as a helper. Use it in two tests. Run `go test -v` on a failure.
2. Compare output of `t.Fail`, `t.Error("reason")`, and `t.Fatal("reason")`. Write a three-row table.
3. Call `t.Fatal` after you register `t.Cleanup`. Confirm that cleanup runs. Print a log from cleanup with `-v`.

#### Advanced practical tasks

1. Write a helper that wraps `t.Run` and a check. Call `t.Helper` in the right functions so that a failure points at the test row.
2. Build a small assert package in `_test.go` (same module) with `Helper`. Document when the helper uses `Fatal` versus `Error`.

---

## Benchmarks: `BenchmarkXxx`

A benchmark function has this form:

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(1, 2)
    }
}
```

The name must start with `Benchmark`. The next rune must not be a lowercase letter. The parameter is `*testing.B`.

`b.N` is the number of iterations. The `testing` package chooses `N` so that the benchmark runs for a stable time. You must put the work under test inside the loop.

From Go 1.24, you can write `for b.Loop() { ... }`. `Loop` reports the timer in a safer way. On Go 1.22 and 1.23, use the `b.N` loop.

Run benchmarks:

```text
go test -bench=BenchmarkAdd -benchmem
go test -bench=. -count=5
```

`-bench` is a regular expression. `-benchmem` prints allocations. `-count` repeats the benchmark.

Timer control:

- `b.ResetTimer()` after setup.
- `b.StopTimer()` and `b.StartTimer()` around code that is not under test.
- `b.ReportAllocs()` to always report allocations.

`b.Run` starts a sub-benchmark. `b.RunParallel` runs the body on several goroutines. `b.SetBytes` reports throughput.

Do not use `testing.T` methods on `b`. `B` has `Error`, `Fatal`, `Helper`, and `Cleanup` that match the test versions.

Keep the compiler from removing the work. Assign the result to a package-level variable if the compiler can see that the result is unused and the function is simple. Read the benchmark results with care. A micro-benchmark is not a full program.

### Questions

#### Theoretical questions

1. Who chooses the value of `b.N`?
2. Why must the work under test run inside the loop?
3. What does `-benchmem` add to the output?
4. When do you call `b.ResetTimer`?
5. What is `b.Loop` and from which Go version is it available?

#### Easy practical tasks

1. Write `BenchmarkAdd` for a small function. Run `go test -bench=BenchmarkAdd`.
2. Add `-benchmem`. Write the columns that you see (`ns/op`, `B/op`, `allocs/op`).
3. Run `go test -bench=. -count=3`. Write why more than one count helps.
4. Read `go doc testing.B`. List five methods in your own words.

#### Medium practical tasks

1. Add setup that allocates a large slice. Call `b.ResetTimer` after the setup. Compare with a version that does not reset the timer.
2. Write two benchmarks: one that appends to a slice without a capacity, one that pre-allocates. Compare `B/op`.
3. Use `b.Run` to benchmark three input sizes. Run `go test -bench=BenchmarkJoin -benchmem`.

#### Advanced practical tasks

1. Use `b.RunParallel` on a function that is safe for concurrent use. Compare with a sequential benchmark. Write when parallel results help.
2. Show compiler dead-code removal: a benchmark that the compiler can empty, and a fix that keeps the work. Use `go test -bench=.` and, if you need, `go test -bench=. -gcflags=-m`.

---

## Examples: `ExampleXxx` (and `Output:` comments)

An example is a function whose name starts with `Example`. Examples appear in the documentation on pkg.go.dev and in `go doc`. `go test` compiles examples. `go test` runs an example when the function has an output comment.

```go
func ExampleAdd() {
    fmt.Println(Add(1, 2))
    // Output:
    // 3
}
```

The comment must use the exact form. The text after `Output:` is compared with the standard output of the example. Leading spaces in the comment follow the `go test` rules. Keep the output simple.

`// Unordered output:` compares lines without order.

Name rules:

- `Example()` is the package example.
- `ExampleAdd` documents function `Add`.
- `ExampleT` documents type `T`.
- `ExampleT_M` documents method `M` on `T`.
- `ExampleAdd_second` is a second example for `Add`. The suffix starts with a lowercase letter.

If you omit `Output:` and `Unordered output:`, `go test` compiles the example and does not execute it as a test.

Do not put assertions with `testing.T` in an example. The output is the check. Print with `fmt`. Keep the example short. The reader copies it.

Examples live in `_test.go` files. You can place them in the same file as tests.

### Questions

#### Theoretical questions

1. When does `go test` execute an example?
2. What does the `Output:` comment compare?
3. How do you name a second example for the same function?
4. What is `Unordered output:` for?
5. Where do examples appear besides `go test`?

#### Easy practical tasks

1. Write `ExampleAdd` with an `Output:` comment. Run `go test -v`.
2. Change the print so that the output does not match. Run `go test`. Save the failure.
3. Add a package-level `Example` that prints the package purpose in one line.
4. Run `go test -run Example`. Write which functions ran.

#### Medium practical tasks

1. Write `ExampleT_Method` for a type with a method. Open the documentation with `go doc` or pkg.go.dev layout and confirm the name rule.
2. Write an example with `Unordered output:` that prints two lines in one order. Confirm that a swap of the two print lines still passes.
3. Add an example without `Output:`. Confirm that `go test` still compiles the package.

#### Advanced practical tasks

1. Write three examples for one function: success path, error path, and a suffix name. Keep each example under ten lines.
2. Compare a unit test and an example for the same function. Write when you keep both and when the example is enough for documentation.

---

## Fuzzing: `FuzzXxx`

A fuzz target has this form:

```go
func FuzzParse(f *testing.F) {
    f.Add("hello")
    f.Add("")
    f.Fuzz(func(t *testing.T, s string) {
        _, err := Parse(s)
        if err != nil {
            return
        }
    })
}
```

The name must start with `Fuzz`. The next rune must not be a lowercase letter. `f.Add` loads seed values. `f.Fuzz` registers the function that the fuzzer calls. The fuzz function takes `*testing.T` and arguments. Allowed argument types include `string`, `[]byte`, fixed integers, `bool`, `rune`, `float32`, and `float64`. Check the current `testing` documentation for the full list.

`go test` without `-fuzz` runs the seed corpus only. That run is like a table of seeds.

Run the fuzzer:

```text
go test -fuzz=FuzzParse -fuzztime=10s
```

The fuzzer writes interesting inputs under `testdata/fuzz`. A failing input is saved. `go test` without `-fuzz` replays that input. Commit useful failing cases. Fix the code. Keep the case if it is a valid regression test.

Write a fuzz function that finds a property. Example: parse then format returns the same bytes. Do not only ignore every error. An error can be correct. A panic or a wrong round trip is a bug.

Do not fuzz a function that deletes files or that calls the network. Keep the target in memory. Use `t.Skip` for inputs that are out of scope when that is the right rule.

The `-fuzz` flag accepts one package. Do not pass `./...` with `-fuzz`.

### Questions

#### Theoretical questions

1. What is the difference between `f.Add` and the function that `f.Fuzz` registers?
2. What does `go test` do with a fuzz target when you omit `-fuzz`?
3. Where does the fuzzer store failing inputs?
4. Which types can the fuzz function accept as arguments?
5. What is a good property to check in a fuzz function?

#### Easy practical tasks

1. Write `FuzzReverse` for a function that reverses a `[]byte`. Add two seeds. Run `go test` without `-fuzz`.
2. Run `go test -fuzz=FuzzReverse -fuzztime=5s`. Write the last lines of the output.
3. Read `go doc testing.F`. List `Add` and `Fuzz` in your own words.
4. Add a seed that panics the function. Run `go test`. Save the failure.

#### Medium practical tasks

1. Fuzz a parse function. Check that a successful parse does not panic on a second parse of the same input.
2. After a failure, find the file under `testdata/fuzz`. Run `go test` with no `-fuzz` and confirm that the case still fails. Then fix the code.
3. Limit fuzzing with `-fuzztime` and `-fuzzminimizetime` (if your version has it). Write the flags that you used.

#### Advanced practical tasks

1. Write a round-trip property: `Parse(Format(v))` equals `v` for valid `v`. Use two arguments if you need. Document invalid inputs that you skip.
2. Compare a table-driven test and a fuzz target for the same parser. Write which bugs each method found in a short experiment.

---

## Test coverage: `go test -cover`

Coverage measures which statements the tests executed.

```text
go test -cover
go test -coverprofile=cover.out
go tool cover -func=cover.out
go tool cover -html=cover.out
```

`-cover` prints a percent for the package. `-coverprofile` writes a file. `go tool cover -func` prints functions. `go tool cover -html` opens a page that colors executed lines.

Cover modes:

- `set` — a statement was executed or not.
- `count` — how many times a statement ran.
- `atomic` — like `count`, safe for concurrent tests. Use `atomic` with `-race`.

```text
go test -covermode=atomic -coverprofile=cover.out
```

`-coverpkg` sets which packages count. Use it when a test in one package exercises another package. Default coverage is the package under test.

A high percent does not mean the tests are good. Coverage does not check that an assertion exists. A test can execute a line and not check the result.

Aim coverage at the packages that you change. Do not chase 100 percent when the last lines are `panic("unreachable")` or generated code. Exclude generated files in the tool that you use for reports when the team agrees.

`go test -cover` is enough for a local check. Store a profile in continuous integration when the team wants a report. Do not fail the build on a small drop without a team rule.

### Questions

#### Theoretical questions

1. What does a coverage percent measure?
2. What is the difference between `-cover` and `-coverprofile`?
3. When do you use `-covermode=atomic`?
4. Why is 100 percent coverage not the same as correct tests?
5. What does `-coverpkg` change?

#### Easy practical tasks

1. Run `go test -cover` on a package with one function and one test. Write the percent.
2. Add `-coverprofile=cover.out`. Run `go tool cover -func=cover.out`. Write the function list.
3. Open `go tool cover -html=cover.out`. Write which lines the report marks as not executed.
4. Add a test for an uncovered branch. Run coverage again. Write the new percent.

#### Medium practical tasks

1. Compare `covermode=set` and `covermode=count` on a loop. Write what the profile shows.
2. Use `-coverpkg` so that a test in `package foo_test` counts coverage of `package foo` and one helper package. Write the command.
3. Find a line that ran but had no assertion. Add an assertion. Write why coverage did not catch the gap.

#### Advanced practical tasks

1. Write a script that fails when `go tool cover -func` shows a total below a number that you choose. Run it on a small module.
2. Generate a coverage report for `./...` and list packages that have no tests. Write a plan for the first package that you will test.

---

## Mocks vs fakes vs stubs; interfaces for testability

Tests need to replace a dependency. The three common replacements are different tools.

A **stub** returns a fixed result. The stub does not implement real behavior. Example: `Get` always returns the same user.

A **fake** is a working small implementation. Example: a map that stores users in memory. The fake follows the same rules as the real store for the cases that you care about.

A **mock** records calls and checks expectations. Example: the test fails if `Save` was not called with this ID. A mock is often generated. The test can become a list of calls instead of a list of behaviors.

Prefer a fake when the dependency is a store, a clock, or a queue. Prefer a stub when the result is a single canned value. Use a mock when the call itself is the behavior (a mail sender must be called once). Do not mock every interface. The test then copies the production call graph.

Design the production code with interfaces at the boundary:

```go
type UserStore interface {
    Find(ctx context.Context, id string) (User, error)
}

type Service struct {
    Store UserStore
}
```

The production `main` function injects a real store. The test injects a fake store. The `Service` type does not import the database package.

Keep interfaces small. Define the interface in the package that uses it, not in the package that implements it, when that split helps. This is the "accept interfaces, return structs" guideline from Topic 7.

Do not add an interface only for a mock when a fake of a concrete type is enough. Do not use the empty interface for a fake.

### Questions

#### Theoretical questions

1. What is the difference between a stub and a fake?
2. What extra job does a mock do that a fake does not do?
3. Why do tests inject a store through an interface?
4. When is a mock the wrong tool?
5. Where do you define a small interface for a dependency?

#### Easy practical tasks

1. Write a `Store` interface with one method. Write a fake that uses a map. Test a `Service` that calls the fake.
2. Replace the fake with a stub that always returns an error. Test the error path of the service.
3. Write a four-column table: term, definition, one use, one misuse.
4. List three dependencies that you would fake and one that you might mock. Give a reason for each.

#### Medium practical tasks

1. Refactor a function that calls `http.Get` directly. Pass an `HTTPDoer` interface (or `*http.Client` with a custom `Transport`). Test with a fake transport or `httptest`.
2. Write a mock by hand that records the last ID. Assert the call after `Service.Delete`. Keep the mock under twenty lines.
3. Compare a test that uses a fake store with a test that uses a mock store for the same service. Write which test is easier to read.

#### Advanced practical tasks

1. Design a `Clock` interface (`Now() time.Time`) and a fake clock. Test a timeout rule without `time.Sleep`.
2. Review a test file that uses a generated mock. Rewrite one test with a fake. Write what you gained and what you lost.

---

## `httptest` for HTTP handlers

The `net/http/httptest` package tests HTTP servers and clients without a real network port in the simple case, and with a local server when you need one.

Unit-test a handler with a request and a recorder:

```go
func TestHealth(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)
    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
    }
}
```

`httptest.NewRequest` builds an incoming request. Prefer `httptest.NewRequestWithContext` when the handler uses `r.Context()`. `httptest.ResponseRecorder` implements `http.ResponseWriter`. Read `rec.Code`, `rec.Body`, and `rec.Header()`.

Integration-test with a server:

```go
srv := httptest.NewServer(handler)
defer srv.Close()
resp, err := srv.Client().Get(srv.URL + "/health")
```

`NewServer` listens on the local machine. `srv.URL` is the base URL. You must call `Close`. `NewTLSServer` starts TLS. `NewUnstartedServer` lets you set fields before `Start`.

From Go 1.22, you can register `GET /health` on `http.NewServeMux` and test that pattern. Use `req := httptest.NewRequest(http.MethodGet, "/health", nil)` and `mux.ServeHTTP`.

Cancel tests: create a context, pass it with `NewRequestWithContext`, cancel it, and assert that the handler returns. For `NewServer`, use a client request with `NewRequestWithContext` and `client.Do`.

Do not listen on a fixed port in tests when `httptest` can choose a port. Do not leak a server. Use `t.Cleanup(srv.Close)`.

### Questions

#### Theoretical questions

1. What is the difference between `ResponseRecorder` and `NewServer`?
2. Why do you call `srv.Close`?
3. Which helper do you use when the handler reads `r.Context()`?
4. What fields of the recorder do you check after `ServeHTTP`?
5. When do you choose `NewTLSServer`?

#### Easy practical tasks

1. Write a handler that writes `ok` and status 200. Test it with `NewRequest` and `NewRecorder`.
2. Assert the `Content-Type` header in that test.
3. Start `NewServer` for the same handler. Call `Get` on `srv.URL`. Close the server.
4. Read `go doc net/http/httptest`. List `NewRequest`, `NewRecorder`, and `NewServer` in your own words.

#### Medium practical tasks

1. Test a `POST` handler that decodes JSON. Send a `bytes.Buffer` body. Assert the status and the response JSON.
2. Test a Go 1.22 mux pattern with a path wildcard. Use `r.PathValue` in the handler. Assert the value.
3. Cancel the request context in a handler that waits on `Done`. Assert the status or the early return.

#### Advanced practical tasks

1. Write a client function that calls a URL. Point it at `httptest.NewServer`. Test success and a 500 status.
2. Test middleware: wrap a handler that records that the middleware ran. Use `NewRecorder` and one `NewServer` test.

---

## `go generate`

`go generate` runs commands that you write in special comments. The compiler does not run those commands. You run `go generate` on purpose.

```go
//go:generate stringer -type=Status
```

The comment must start at the beginning of the line. The form is `//go:generate` plus a command. The command runs in the package directory.

```text
go generate .
go generate ./...
```

Common uses:

- `stringer` for `String` methods on integer types.
- Code for protocol buffers.
- A mock generator, when the team accepts generated mocks.
- Embed-related helpers, when `embed` is not enough.

The generated file is normal Go source. Commit the file if the team wants builds without the generator. Or generate in continuous integration if the team wants that rule. Pick one rule.

Do not use `go generate` for the build itself. `go build` must work after the files exist. Do not hide required steps only in a generate comment without a README or a script.

`go generate` does not run tests. It does not format files unless the command does. Run `go fmt` on generated files if the generator does not format them.

Variables such as `$GOFILE` and `$GOPACKAGE` are available in the command. Read `go help generate` for the list.

### Questions

#### Theoretical questions

1. Does `go build` run `//go:generate` directives?
2. Where does the generate command run?
3. What is a valid reason to use `go generate`?
4. Why must a teammate know whether generated files are committed?
5. What is the exact comment prefix for a generate directive?

#### Easy practical tasks

1. Add `//go:generate echo generate-ok` to a package file. Run `go generate -v .`. Write the output.
2. Run `go help generate`. List three `$GO*` variables.
3. Write four sentences: what generate is, what it is not, when you run it, where the comment lives.
4. Put the directive in the wrong place (after a space). Run `go generate`. Write what happens.

#### Medium practical tasks

1. Install or use a small generator (for example `stringer` if you can, or a tiny local program). Generate a file. Show the command and the file.
2. Add `go generate ./...` to a script that also runs `go test`. Document the order.
3. Create two packages with generate directives. Run `go generate ./...` from the module root. Confirm both run.

#### Advanced practical tasks

1. Write a tiny generator program in `internal/cmd/gen` that writes a Go file. Call it from `//go:generate`. Make the output `gofmt` clean.
2. Compare `go generate` with a `//go:embed` file. Write when you generate code and when you embed data.

---

## Linters: `staticcheck`, `golangci-lint`

A linter reports problems that the compiler accepts. `go vet` is the first linter. It ships with Go. Run `go vet ./...` on every change.

`staticcheck` is a widely used linter. It finds dead code, wrong `errors.Is` use, unused values, and many other issues. The checks have names such as `SA4006`. Install the official `staticcheck` binary. Run:

```text
staticcheck ./...
```

`golangci-lint` runs many linters in one command. It can run `govet`, `staticcheck`, `gofmt`, and others. Install the official release. Run:

```text
golangci-lint run
```

A configuration file (`.golangci.yml` or the current name in the `golangci-lint` documentation) selects linters. Start with a small set. Enable `staticcheck` and `govet`. Add more linters when the team agrees.

Rules for a beginner:

- Fix `go vet` first.
- Then run `staticcheck` or `golangci-lint`.
- Do not disable a check without a reason.
- Do not format the configuration debate. The official Go format stays `gofmt`.

`golangci-lint` version and the Go version must work together. Read the install page for your Go version.

Linters are not tests. A clean lint run does not prove behavior. Tests prove behavior. Linters reduce mistakes.

Continuous integration must run `go test` and at least one linter. Topic 1 already uses `go vet`. This topic adds `staticcheck` and `golangci-lint`.

### Questions

#### Theoretical questions

1. What class of problems does a linter find that `go test` can miss?
2. What is `staticcheck`?
3. What is `golangci-lint`?
4. Why must you still run tests when the linter is clean?
5. Which linter ships with the Go toolchain?

#### Easy practical tasks

1. Run `go vet ./...` on your module. Write the result.
2. Install `staticcheck`. Run `staticcheck ./...`. Write the result.
3. Write a small dead-assignment bug (`x := 1` then `x = 2` without a read, or an unused result). Run `staticcheck`. Record the check name.
4. Read the `staticcheck` documentation for one `SA` check. Write the check name and the idea.

#### Medium practical tasks

1. Install `golangci-lint`. Run `golangci-lint run` with the default or a minimal config. Write which linters ran.
2. Add a configuration file that enables `govet` and `staticcheck`. Run the linter again.
3. Fix every finding in a small package. Show a clean run.

#### Advanced practical tasks

1. Add `golangci-lint` and `go test -race` to a script. Fail the script on the first error. Document the tools and versions.
2. Compare `staticcheck` alone with `golangci-lint` on the same module. Write two findings that only one setup reported, or write that they matched.

---

## Profiling: `pprof` (CPU, heap, goroutines)

`pprof` records where a program spends time and memory. The `runtime/pprof` package writes profiles. The `net/http/pprof` package serves profiles over HTTP. The `go tool pprof` command reads a profile.

CPU profile. The sampler notes the call stack at a rate. Use it when the program is slow.

Heap profile. The profile shows allocations. Use it when memory grows or when allocation rate is high.

Goroutine profile. The profile lists goroutines and their stacks. Use it when you suspect a leak or a deadlock that is not global.

For tests and benchmarks:

```text
go test -cpuprofile=cpu.out -bench=BenchmarkAdd
go test -memprofile=mem.out -bench=BenchmarkAdd
go tool pprof cpu.out
```

For a server, import the pprof handlers:

```go
import _ "net/http/pprof"
```

Register them on a mux that you control, or use `http.DefaultServeMux` on a debug port. The paths look like `/debug/pprof/profile` (CPU), `/debug/pprof/heap`, and `/debug/pprof/goroutine`. Do not expose that port on a public network without protection.

```text
go tool pprof http://127.0.0.1:6060/debug/pprof/profile?seconds=30
go tool pprof http://127.0.0.1:6060/debug/pprof/heap
go tool pprof http://127.0.0.1:6060/debug/pprof/goroutine
```

`go tool pprof` can start a local page with `-http`. Use `top`, `list`, and `web` in the pprof prompt.

You can also write a profile from code with `pprof.StartCPUProfile` and `pprof.WriteHeapProfile`. Stop the CPU profile with `pprof.StopCPUProfile`.

Read CPU time as a hint, not as a proof. Measure before you change the code. Measure after. Topic 15 covers allocation and escape analysis.

### Questions

#### Theoretical questions

1. What does a CPU profile record?
2. What does a heap profile help you find?
3. When do you take a goroutine profile?
4. Which import serves pprof over HTTP?
5. Why must you protect a `/debug/pprof` endpoint?

#### Easy practical tasks

1. Run a benchmark with `-cpuprofile=cpu.out`. Run `go tool pprof -top cpu.out` (or the equivalent). Write the top function.
2. Run the same benchmark with `-memprofile=mem.out`. Write the top allocator.
3. Read `go doc runtime/pprof`. Write the purpose of `StartCPUProfile`.
4. Make a three-row table: profile type, command or path, one question the profile answers.

#### Medium practical tasks

1. Start a small HTTP server with `net/http/pprof` on `localhost`. Fetch `/debug/pprof/goroutine?debug=1`. Write how many goroutines you see.
2. Create a goroutine leak on purpose. Compare a goroutine profile before and after. Then fix the leak.
3. Use `-http` with `go tool pprof` if your environment allows a browser. Write the name of one box or one line that you inspect. If you cannot open a browser, use `list FunctionName` in the pprof prompt.

#### Advanced practical tasks

1. Profile a CPU-bound function and a memory-heavy function. Change one line (pre-allocate a slice, or reduce work). Show `top` before and after.
2. Write a test or a script that captures a heap profile of a function under test with `runtime/pprof`. Document the steps so that a teammate can repeat them.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a `_test.go` file to a coverage HTML report. Name each command and the role of `TestXxx`.
2. How do table-driven tests, subtests, and fuzz seeds work together for one parser?
3. A teammate wants a mock for every interface and 100 percent coverage. Which facts do you use to choose fakes, tables, and a coverage target?
4. When do you use `httptest.NewRecorder`, when do you use `httptest.NewServer`, and when do you still write a fake `RoundTripper`?
5. How do `go generate`, `staticcheck`, and `pprof` differ from `go test` in a daily workflow?

#### Easy practical tasks

1. Create a module with `Sum`, `TestSum` (table with three rows), `ExampleSum`, and `BenchmarkSum`. Run `go test -cover`, `go test -bench=BenchmarkSum`, and `go test -run Example`.
2. Make a one-page cheat sheet: `go test` flags (`-v`, `-run`, `-cover`, `-race`, `-bench`, `-fuzz`), `t.Run`, `t.Helper`, `t.Fatalf`, `httptest`, `go generate`, `staticcheck`.
3. Draw a diagram of a `Service` with a `Store` interface, a fake store in tests, and `httptest` for the HTTP handler.
4. Run `go doc testing` and list `T`, `B`, `F`, and `Example` naming rules in four short sentences.

#### Medium practical tasks

1. Write a handler plus `Sum` helper. Test the helper with a table. Test the handler with `httptest`. Add `-coverprofile` and open the HTML report. Add one missing assertion that coverage did not require.
2. Add a fuzz target for `Sum` or for a parse helper. Run it for five seconds. Keep at least two seeds in `f.Add`.
3. Run `go vet` and `staticcheck` (or `golangci-lint`) on the module. Fix findings. Show a clean test run with `-race`.

#### Advanced practical tasks

1. Build a small package: public API, table tests, one fake dependency, one `httptest` server test, one benchmark, one example, and a fuzz target. Write a README section that lists the commands to run all of them.
2. Capture a CPU or heap profile of the benchmark. Change the code to reduce allocations. Show `-benchmem` and a `pprof` `top` line before and after. Keep `go test -race` clean.
