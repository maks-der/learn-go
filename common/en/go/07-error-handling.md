# 7. Error Handling

## Description

Go treats an error as a normal value. A function returns an `error` when the operation can fail. The caller checks the value and decides the next step. This topic shows the `error` interface, how you create and wrap errors, and how you inspect an error chain. This topic also shows when `panic` is correct.

Complete topic 4 before this topic. Topic 4 introduces `panic` and `recover`. This topic gives the rules for production code. Use one term for each idea. Use `error` for a value that implements the `error` interface. Use `panic` only for a failure that the program must not continue.

---

## The `error` interface

The predeclared `error` interface has one method:

```go
type error interface {
	Error() string
}
```

Any type that has a method `Error() string` is an `error`. The method returns a message for humans and for logs. Do not parse that string in logic. Inspect the typed value or the chain.

The zero value of the `error` interface is `nil`. A `nil` error means success. A non-nil error means failure. Check the error as soon as the function returns:

```go
n, err := w.Write(p)
if err != nil {
	return err
}
```

By convention the `error` is the last result. A typical signature is `func Open(name string) (*File, error)`. When `err != nil`, do not use the other results unless the package documents a special case.

An interface value is `nil` only when the type and the data are both unset. A function that returns `error` must return a `nil` interface, not a nil pointer stored in an interface. This bug is common:

```go
var p *MyError = nil
return p // the error is not nil
```

Return `nil` directly. Do not return a nil pointer as an `error`.

### Questions

#### Theoretical questions

1. What method does the `error` interface require?
2. What does a `nil` error mean?
3. Why is `error` the last result in a function signature?
4. Why must you not parse `Error()` text to decide program logic?
5. When is an interface value that holds a nil pointer not equal to `nil`?

#### Easy practical tasks

1. Write a type `msgError` with one `string` field and a method `Error() string`. Create a value and print it with `fmt.Println`.
2. Write a function `func half(n int) (int, error)` that returns an error when `n` is odd. Call it with `2` and with `3`. Check `err` in both cases.
3. Write `if err != nil` after a call to `os.Open` on a name that does not exist. Print `err`.
4. List four standard library functions that return `error`. Write the package name for each.

#### Medium practical tasks

1. Return a nil pointer of a named error type as `error`. Compare the result with `nil`. Record the result. Then return a true `nil` interface and compare again.
2. Write a function with three results where `error` is not last. Run `go vet`. Record the message if the tool reports a problem. Move `error` to the last place.
3. Write two types that both implement `error`. Store each value in a variable of type `error`. Print `err.Error()` for both.

#### Advanced practical tasks

1. Read the source of `os.PathError` in the standard library. Write how the type implements `error` and which fields it stores.
2. Design a function API that returns `(T, error)` and a second API that returns only `error`. Write when each shape is correct. Give one example for each.

---

## `errors.New`, `fmt.Errorf`, and `%w`

The package `errors` creates a simple error:

```go
err := errors.New("not found")
```

`errors.New` stores the string. The string does not change later. Use `errors.New` for a fixed message with no extra data.

`fmt.Errorf` formats a message like `fmt.Sprintf` and returns an `error`:

```go
err := fmt.Errorf("user %d not found", id)
```

The verb `%w` wraps another error. The new error keeps the original error in a chain:

```go
err := fmt.Errorf("load config: %w", fileErr)
```

Use `%w` when the caller must still detect the original error with `errors.Is` or `errors.As`. Use `%v` or `%s` when you only need text and you do not keep the chain. A format string must use `%w` only with an `error` operand.

From Go 1.20 you can wrap more than one error in one call:

```go
err := fmt.Errorf("copy %s to %s: %w; %w", src, dst, errRead, errWrite)
```

`errors.Join` also combines errors. `errors.Join` returns `nil` when every input is `nil`. Use `errors.Join` when two or more independent steps fail and you must report all of them.

Do not wrap the same error many times without a reason. Each wrap adds a frame of context. Keep the message short. Do not start the message with a capital letter. Do not end the message with a period. That is the style of the standard library.

Add context at the layer that has the data. Wrap at the point where you call another function. Do not put secrets in the message. Do not wrap at every line inside one function when no new fact exists.

### Questions

#### Theoretical questions

1. What is the difference between `errors.New` and `fmt.Errorf` without `%w`?
2. What does the verb `%w` add that `%v` does not add?
3. What operand type must you pass to `%w`?
4. When do you use `errors.Join` instead of a single `fmt.Errorf`?
5. What style rule does the standard library use for the first letter of an error message?

#### Easy practical tasks

1. Create an error with `errors.New("denied")`. Print the result of `err.Error()`.
2. Create an error with `fmt.Errorf("port %d in use", 8080)`. Print the error.
3. Wrap `io.EOF` with `fmt.Errorf("read body: %w", io.EOF)`. Print the new error.
4. Open a missing file. Return `fmt.Errorf("open %s: %w", path, err)` from a helper. Print the error in `main`.

#### Medium practical tasks

1. Wrap an error two times with `%w`. Print the outer message. Then use `errors.Unwrap` twice and print each step.
2. Use `errors.Join` on one `nil` error and one non-nil error. Print the join result. Then join two non-nil errors and print that result.
3. Wrap `io.EOF` without `%w` (use `%v`). Show that `errors.Is` does not find `io.EOF`.

#### Advanced practical tasks

1. Use one `fmt.Errorf` with two `%w` verbs (Go 1.20 or later). Confirm `errors.Is` finds each inner error.
2. Compare `errors.Join(errA, errB)` with `fmt.Errorf("%w: %w", errA, errB)`. Write how the printed text differs and how `errors.Is` behaves for both.

---

## `errors.Is` and `errors.As`

After you wrap errors, `==` does not find the inner value. Type assertion on the outer error does not find the inner type. Use the functions in package `errors`.

`errors.Is(err, target)` walks the chain. It returns true when a value in the chain matches `target`. A match uses `==` or a method `Is(error) bool` on a type in the chain. Use `errors.Is` for sentinel errors:

```go
if errors.Is(err, os.ErrNotExist) {
	// the file is missing
}
```

`errors.As(err, &target)` walks the chain. It finds the first value that matches the type of `target`. `target` must be a non-nil pointer to a type that implements `error`, or a pointer to an interface. On success, `errors.As` writes the value into `target` and returns true:

```go
var pe *os.PathError
if errors.As(err, &pe) {
	fmt.Println(pe.Path)
}
```

`errors.Unwrap(err)` returns the next single error in the chain. It returns `nil` when there is no next error. `errors.Unwrap` does not walk a join of many errors. `errors.Is` and `errors.As` do walk a join.

Do not write `err == os.ErrNotExist` after a wrap. Do not write `err.(*os.PathError)` after a wrap. Use `errors.Is` and `errors.As`.

When `err` is `nil`, `errors.Is(err, target)` is true only when `target` is also `nil`.

### Questions

#### Theoretical questions

1. Why does `err == sentinel` fail after `fmt.Errorf` with `%w`?
2. What does `errors.Is` compare at each step of the chain?
3. What argument type must you pass as the second parameter of `errors.As`?
4. What does `errors.Unwrap` return when the error has no inner error?
5. Does `errors.Unwrap` walk every error from `errors.Join`? Explain.

#### Easy practical tasks

1. Wrap `os.ErrNotExist` with `%w`. Show that `==` is false and `errors.Is` is true.
2. Wrap an `*os.PathError` from `os.Open` on a missing file. Use `errors.As` and print `Path`.
3. Call `errors.Unwrap` on an error from `fmt.Errorf("x: %w", io.EOF)`. Print the result.
4. Call `errors.Is(nil, nil)` and `errors.Is(nil, io.EOF)`. Write the two results.

#### Medium practical tasks

1. Build a chain of three wraps around `io.EOF`. Write a loop that calls `errors.Unwrap` until the result is `nil`. Print each step.
2. Use a type assertion on a wrapped `*os.PathError`. Record the panic or the failed assertion. Then use `errors.As` and record success.
3. Join two different sentinels with `errors.Join`. Show that `errors.Is` finds each sentinel.

#### Advanced practical tasks

1. Add a method `Is(target error) bool` on a custom error type so that `errors.Is(err, os.ErrNotExist)` is true for that type. Prove it with a small program.
2. Show what `errors.Unwrap` returns for an `errors.Join` result. Compare that result with `errors.Is` on each joined error.

---

## Sentinel errors vs custom error types

A sentinel error is a package-level variable that callers compare with `errors.Is`:

```go
var ErrNotFound = errors.New("not found")
```

The standard library uses sentinels. Examples: `io.EOF`, `os.ErrNotExist`, `sql.ErrNoRows`. A sentinel is correct when the failure is one known case and the caller needs no extra fields.

A custom error type is a struct (or another named type) with an `Error() string` method. Add fields for data that the caller must read:

```go
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string {
	return e.Query + ": " + e.Err.Error()
}

func (e *QueryError) Unwrap() error {
	return e.Err
}
```

Add `Unwrap` when the type wraps another error. Then `errors.Is` and `errors.As` walk through the type.

Use a sentinel when the case is stable and has no extra data. Use a custom type when the caller must read a path, a status code, or another field. Do not export a large set of sentinels that never appear in caller code. Do not use a custom type when a sentinel is enough.

Compare sentinels with `errors.Is`, not with `==`, when any wrap is possible. Detect custom types with `errors.As`, not with a type assertion, when any wrap is possible.

### Questions

#### Theoretical questions

1. What is a sentinel error?
2. Name two sentinel errors from the standard library.
3. When must you use a custom error type instead of a sentinel?
4. Why does a wrapping custom type need an `Unwrap` method?
5. Why is `errors.Is` safer than `==` for an exported sentinel?

#### Easy practical tasks

1. Declare `var ErrDenied = errors.New("denied")` in a small package. Return it from a function. Detect it with `errors.Is` in `main`.
2. Write a struct `CodeError` with fields `Code int` and `Msg string`. Implement `Error()`. Print a value.
3. Make a two-column table: "Sentinel" and "Custom type". Add four rows of differences.
4. Find `io.EOF` in `go doc`. Write one sentence about when a reader returns it.

#### Medium practical tasks

1. Write `QueryError` with `Unwrap`. Wrap `sql.ErrNoRows` (or `errors.New("no rows")`). Show that `errors.Is` and `errors.As` both succeed.
2. For three APIs (file missing, HTTP 404, checksum mismatch), choose sentinel or custom type. Write one reason for each choice.
3. Export a sentinel and wrap it in another package with `%w`. Detect the sentinel from `main` with `errors.Is`.

#### Advanced practical tasks

1. Implement `Is` on a custom type so that several HTTP status codes map to one sentinel. Prove `errors.Is` for two codes.
2. Read `os.PathError` and `fs.PathError`. Write how a sentinel (`os.ErrNotExist`) and a custom type work together in one failure.

---

## When panic is appropriate

An `error` is for a failure that the caller can handle. A `panic` stops the current goroutine. If nobody calls `recover`, the process prints a stack and exits.

Use `panic` only when the program cannot continue in a safe way. Typical cases:

- A programmer invariant is false. Example: a switch misses a case that the type system cannot exclude, and you treat that as a bug.
- Initialization cannot finish and no useful program remains. Example: a required packed asset is missing from the binary.
- You reach a state that means memory corruption or an impossible contract.

Do not use `panic` for expected failures. Missing files, bad user input, network timeouts, and HTTP 4xx cases are errors. The caller must receive an `error`.

Do not use `panic` as a second return path in a library. Library code that panics on normal input forces every caller to use `recover`. That design is hard to test and hard to compose.

`recover` belongs at a defined boundary. The HTTP server in the standard library recovers from a panic in a handler so that one request does not stop the process. Do not scatter `recover` in every function. Recover at the boundary, and log the panic.

A nil pointer dereference, an index out of range, and a send on a closed channel also panic. Those panics are bugs. Fix the code. Do not hide them with a wide `recover` unless you are at a process boundary.

A `Must*` helper calls a function that returns `(T, error)` and panics when the error is not `nil`. Use `Must*` when the input is a fixed literal or an embedded file and a failure means the binary is wrong. Examples: `regexp.MustCompile` and `template.Must`. Do not use `Must*` on user input.

### Questions

#### Theoretical questions

1. What happens to a goroutine when it panics and nobody calls `recover`?
2. Name two situations where `panic` is appropriate.
3. Why must a library not panic on bad user input?
4. Where is `recover` appropriate in a server?
5. When is a `Must*` helper correct?

#### Easy practical tasks

1. Make a two-column table: "Situation" and "error or panic". Add five rows (missing file, bad JSON from a client, missing embed file at init, division by zero that you can check, impossible enum case).
2. Write a function that panics when a length is negative. Call it with `-1`. Read the stack in the output.
3. Write four sentences that explain when you return an `error` instead of a panic.
4. Use `regexp.MustCompile` with a valid constant pattern. Then try a bad pattern and record the panic.

#### Medium practical tasks

1. In `init`, panic if a compiled regex is wrong. In `main`, return an `error` if a user flag is wrong. Show both programs.
2. Convert a function that panics on a missing key into a function that returns `(T, error)`. Keep a panic only for a broken internal invariant.
3. Start a tiny HTTP handler that panics. Use the default server. Record whether the process stays up after one request.

#### Advanced practical tasks

1. Write a team rule of one page: when to panic, when to return `error`, and where to `recover`. Apply the rule to a small CLI and a small HTTP handler.
2. Trigger a panic with a nil map write. Explain the stack. Then prevent the panic with a `make` call. Do not use `recover` for this fix.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a low-level `os` error to a caller that uses `errors.Is` after two wraps.
2. How do you choose among `errors.New`, `fmt.Errorf` with `%w`, a sentinel, and a custom type?
3. Why do `errors.Is` and `errors.As` exist when `==` and type assertions already exist?
4. When must a public API hide the inner error instead of wrapping with `%w`?
5. How do you decide between a returned error, a `Must*` panic at init, and a recover at a server boundary?

#### Easy practical tasks

1. Write `func load(path string) ([]byte, error)` that wraps `os.ReadFile` with `%w` and the path. Call it for a missing file and print the error.
2. Detect `os.ErrNotExist` with `errors.Is` on that wrapped error. Print `true` or `false`.
3. Declare `ErrEmpty` as a sentinel. Return it when the path is `""`. Detect it with `errors.Is`.
4. Draw a table: Create error, Wrap error, Test sentinel, Test type, Panic. Add one function name or form in each cell.

#### Medium practical tasks

1. Build three packages: `store`, `service`, `cmd`. Each layer wraps with `%w` and one new word. Print the final error. Show `errors.Is` for the root sentinel.
2. Write a custom `*PathError` with `Unwrap`. Use `errors.As` in `main` to print the path after a wrap.
3. Review five `return err` lines in a small program. Keep the wrap only at package boundaries. Record which lines you change.

#### Advanced practical tasks

1. Write a helper `func wrap(op, name string, err error) error` that returns `nil` when `err` is `nil`. Use it at three call sites. Prove that a nil input stays nil.
2. Design an error policy for a library: which errors stay in the chain, which errors become sentinels, and which errors lose the cause. Write the policy in one page and implement one function that follows it.
