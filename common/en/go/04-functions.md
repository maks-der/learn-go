# 4. Functions

## Description

A function is a named block of statements. You declare a function one time. You call the function many times. This topic shows how you write functions, how you return values, and how you use function values.

Go functions can return more than one value. The common pattern is a result and an error. Later topics use that pattern in every program.

A function can capture variables from an outer block. That function is a closure. A function can call itself. That call is recursion. The `init` function runs before `main`. Panic stops normal control flow.

Use functions to split work. Give each function one clear job. Return errors for expected failures. Use panic only for faults that the program cannot handle.

---

## Declaration, multiple return values, and named returns

A function declaration starts with `func`. Next comes the name. Next come the parameters in parentheses. Next comes the result type when the function returns a value. The body is a block.

```go
func add(a int, b int) int {
	return a + b
}
```

Adjacent parameters of the same type can share one type name. `func add(a, b int) int` is the same as `func add(a int, b int) int`.

You cannot declare a named function inside another function. You can declare an anonymous function inside a function. Two functions in one package cannot have the same name. Call a function with the correct number of arguments. The argument types must match the parameter types. Go does not convert types in a call.

A function can return more than one result. Write the result types in parentheses. The common pattern is `(value, error)`. The last result is often `error`. A `nil` error means success. A non-nil error means failure. Check the error before you use the value.

```go
func div(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("divide by zero")
	}
	return a / b, nil
}

q, err := div(10, 2)
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(q)
```

Use the blank identifier `_` when you must ignore one result. Do not ignore an `error` result.

You can name the result parameters. The names become variables in the function body. Go sets them to the zero value at the start. A `return` with no expressions is a naked return. Use a naked return only in a short function.

```go
func rect(w, h int) (area int, perim int) {
	area = w * h
	perim = 2 * (w + h)
	return
}
```

A deferred function can change a named result. Prefer an explicit `return` when the function has more than a few lines. Callers do not need to use the same names.

### Questions

#### Theoretical questions

1. What keyword starts a function declaration?
2. How do you write the result types for two results?
3. What does a `nil` error mean in the `(value, error)` pattern?
4. What is a named return value?
5. What is a naked return?

#### Easy practical tasks

1. Write `double` that takes one `int` and returns that value times two. Call it from `main`.
2. Write `swap` that takes two `int` values and returns them in reverse order. Print both results.
3. Call a function that returns `(int, error)`. Handle the error with `if err != nil`.
4. Rewrite `rect` with named results `area` and `perim`. Use a naked return.

#### Medium practical tasks

1. Write `parsePair` that reads two integers from one string. Return both integers and an error.
2. Chain two calls that each return `(int, error)`. Stop at the first error.
3. Add a `defer` that changes a named `error` result. Show the value in `main`.

#### Advanced practical tasks

1. Write `div` and a caller that tests success, divide by zero, and a negative input rule that you define.
2. Write a function that opens a resource, uses `defer`, and sets a named `error` when cleanup fails.

---

## Variadic functions

A variadic parameter accepts zero or more arguments of one type. Write `...T` as the last parameter. Inside the function that parameter is a slice of type `[]T`.

```go
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	fmt.Println(sum())
	fmt.Println(sum(1, 2, 3))
}
```

Only the last parameter can be variadic. `func bad(x ...int, y int)` is not valid.

You can pass an existing slice with `...` after the slice. `sum(values...)` passes the elements. The function still sees a slice.

```go
values := []int{4, 5, 6}
fmt.Println(sum(values...))
```

`append` is a built-in variadic function. `fmt.Println` is also variadic.

An empty call `sum()` gives a slice of length zero. The slice can be `nil`. Range over that slice is safe.

Do not mix a slice argument and extra elements in one call. `sum(values..., 9)` is not valid.

### Questions

#### Theoretical questions

1. Where must a variadic parameter sit in the list?
2. What is the type of a variadic parameter inside the function?
3. How do you pass a slice to a variadic parameter?
4. How many extra arguments can a caller omit?
5. Why is `func f(a ...int, b int)` invalid?

#### Easy practical tasks

1. Write `join` that takes `...string` and prints each string on one line.
2. Call your `sum` with no arguments, with three literals, and with a slice and `...`.
3. Write `max` that takes `...int` and returns the largest value. Define the empty-list behavior.
4. Call `fmt.Println` with four arguments of different types.

#### Medium practical tasks

1. Write `avg` that takes `...float64` and returns the average and an error when the list is empty.
2. Write `appendInt` that wraps `append` for `[]int` and returns the new slice. Show a call that grows the slice.
3. Pass one slice to two variadic functions. Show that each function can read the same elements.

#### Advanced practical tasks

1. Write `concat` that takes a separator and `...string`. Return one string. Do not add a separator at the end.
2. Measure or reason about `f(s...)` versus a parameter of type `[]T`. Write when you choose each form.

---

## Functions as values and closures

A function is a value. You can assign a function to a variable. You can pass a function as an argument. You can return a function.

The type of a function is its signature. `func(int) int` is a type. Two functions have the same type when the parameter types and the result types match. The names in the signature are not part of the type.

```go
func square(n int) int { return n * n }

func apply(n int, fn func(int) int) int {
	return fn(n)
}

func main() {
	f := square
	fmt.Println(apply(4, f))
}
```

The zero value of a function type is `nil`. A call of a `nil` function value causes a panic. Check for `nil` when a function parameter is optional.

An anonymous function has no name. You write `func(...) ... { ... }`. You can call it at once.

A closure is a function value that refers to variables from an outer block. The function can read those variables. The function can assign those variables.

```go
func counter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}
```

The inner function keeps `n` alive after `counter` returns. Each call of `counter` creates a new `n`. A closure captures the variable, not a copy of the current value.

Go 1.22 gives each iteration of a `for` loop its own loop variable. Closures that capture the loop variable in Go 1.22 see the value for that iteration.

Use function values for callbacks and for small strategies. Use closures for small state. Prefer a struct when the state has several fields.

### Questions

#### Theoretical questions

1. What is a function type?
2. When do two functions have the same type?
3. What is the zero value of a function variable?
4. What does a closure capture?
5. What did Go 1.22 change about loop variables and closures?

#### Easy practical tasks

1. Write `apply` as in the example. Pass `square` and a second function that adds one.
2. Write an anonymous function that returns `10`. Call it in the same statement.
3. Write `counter` as in the example. Call the result three times. Print the three numbers.
4. Create two counters. Show that they do not share `n`.

#### Medium practical tasks

1. Write `filter` that takes `[]int` and `func(int) bool`. Return the values that pass the test.
2. Write `adder(base int) func(int) int` that returns a function that adds `base`.
3. Store three functions in a map from `string` to `func(int) int`. Call each by name.

#### Advanced practical tasks

1. Write `compose` that takes two `func(int) int` values and returns `g(f(x))`.
2. Implement a limited-use token with a closure: the first `n` calls succeed, then the function returns an error.

---

## Recursion

Recursion is a function call to the same function. Each call must move toward a base case. The base case returns without a new recursive call.

```go
func fact(n int) int {
	if n <= 1 {
		return 1
	}
	return n * fact(n-1)
}
```

Use a clear base case. A missing base case causes infinite recursion. Infinite recursion grows the stack. The program then panics with a stack overflow.

Go does not guarantee tail-call optimization. A deep recursive call can overflow the stack even when the last action is a call. Use a loop when the depth can be large.

Recursion fits trees, nested structures, and divide-and-conquer algorithms. A slice or a number with a small bound can use recursion for learning.

A recursive function can have more than one recursive call. Example: a tree walk that visits left and right children.

Keep the recursive function pure when you can. Pass the needed data as parameters. Do not hide the progress in a global variable.

### Questions

#### Theoretical questions

1. What is a base case?
2. What happens when a recursive function has no base case?
3. Does Go promise tail-call optimization?
4. When do you replace recursion with a loop?
5. Name one data shape that fits recursion.

#### Easy practical tasks

1. Write `fact` for `n >= 0`. Print `fact(5)`.
2. Write `sumN` that returns `1 + 2 + ... + n` with recursion.
3. Write `countdown` that prints `n`, then `n-1`, down to `1`.
4. Write the base case for a function that walks a string by index.

#### Medium practical tasks

1. Write recursive `fib`. Then write a loop version. Compare the two for `n := 10`.
2. Write a recursive function that reverses a string.
3. Write a recursive search in a slice that returns the index or `-1`.

#### Advanced practical tasks

1. Walk a small binary tree with a recursive function. Print values in order.
2. Find a recursion depth that panics on your machine. Then rewrite the same job with a loop. Record the two outcomes.

---

## `init` (when it runs; why not to overuse it)

An `init` function has this form:

```go
func init() {
	fmt.Println("init runs")
}
```

You do not call `init`. The runtime calls `init` during package initialization. `init` has no parameters and no results. You cannot refer to `init` by name.

Package initialization follows this order:

1. The runtime initializes imported packages first.
2. The runtime initializes package-level variables in the current package.
3. The runtime runs `init` functions in the current package.
4. After all packages initialize, the runtime calls `main` in package `main`.

One file can have more than one `init` function. Those functions run in source order. Do not depend on file order across a package if you can avoid it.

A blank import `_ "pkg"` imports a package only for its `init` side effects. Database drivers use that pattern. The import is easy to miss in a review.

Avoid many `init` functions. They hide work. Tests cannot skip them. The start order across packages is hard to see. Prefer an exported `Setup` function that `main` calls.

Use `init` only for registration that must run at start. Example: register a driver. Keep `init` short. Do not read files or start servers in `init` without a strong reason.

### Questions

#### Theoretical questions

1. Who calls `init`?
2. When does `init` run relative to `main`?
3. Can you call `init` from your code?
4. What is a blank import `_ "pkg"` for?
5. Why must you avoid many `init` functions?

#### Easy practical tasks

1. Add one `init` that prints `init`. Add `main` that prints `main`. Run the program. Record the order.
2. Add two `init` functions in one file. Record the print order.
3. Write a package-level variable with an initializer that prints `var`. Add `init` and `main` prints. Record the order.
4. Write three sentences that state when `init` is acceptable.

#### Medium practical tasks

1. Create package `a` that imports package `b`. Put a print in each `init`. Record the order from `main`.
2. Use a blank import of a small local package. Show that `init` in that package runs.
3. Replace an `init` that sets a map with an explicit `Setup` function. Call `Setup` from `main`.

#### Advanced practical tasks

1. Document the start sequence of a module with three packages. Draw the import graph and the `init` order.
2. Find one standard library package that uses `init`. Write what it registers and why a user might still avoid extra `init` in app code.

---

## Panic and recover

A panic stops normal execution of the current goroutine. The runtime runs deferred functions. Then the program crashes unless a deferred function calls `recover`.

```go
func mustPositive(n int) {
	if n < 0 {
		panic("n must be >= 0")
	}
}
```

An error is a value. You return an error for expected failures. Examples: a missing file, bad user input, a network timeout. The caller checks the error and chooses the next step.

A panic is for a fault that the program cannot handle in a clean way. Examples: a broken invariant, a nil dereference, an out-of-range index. Treat a panic as a bug or as a last stop.

`recover` stops a panic. `recover` works only inside a deferred function. `recover` returns the value that was passed to `panic`. `recover` returns `nil` when there is no panic.

```go
func safe() (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("panic: %v", v)
		}
	}()
	panic("boom")
}
```

Do not use `panic` for ordinary control flow. Do not hide panics in a library without a documented boundary. Some servers recover at the request boundary. That recover protects the process. It does not fix the bug.

Prefer `error` results in your functions. Use `panic` when the program is wrong and must stop. Use `recover` only at a clear boundary.

### Questions

#### Theoretical questions

1. What is the difference between an error and a panic?
2. When do deferred functions run during a panic?
3. Where must you call `recover` for it to work?
4. What does `recover` return when there is no panic?
5. Why is panic a bad tool for user input errors?

#### Easy practical tasks

1. Write a function that panics on an empty string. Call it from `main` with a non-empty string and then with an empty string.
2. Write `safe` as in the example. Print the returned error.
3. Call `recover` in a normal function that is not deferred. Show that it does not stop a later panic in the same function.
4. List three failures that must be errors and two failures that may be panics.

#### Medium practical tasks

1. Write `mustAtoi` that panics when `strconv.Atoi` fails. Then write a wrapper that recovers and returns an `error`.
2. Add three `defer` prints around a panic. Record the print order.
3. Show that a panic in function `A` runs defers in `A`, then defers in the caller, unless `recover` stops it.

#### Advanced practical tasks

1. Write a small worker function that recovers from a panic and returns an error to the caller. The process must keep running.
2. Read how `net/http` documents panic recovery in handlers. Write five sentences on why a server recovers at that boundary.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of a call that returns `(T, error)` from the callee to a correct check in the caller.
2. How do named results, `defer`, and `recover` work together in one function?
3. When do you choose a variadic parameter, a slice parameter, or a function value parameter?
4. Why does the runtime initialize imported packages before package-level variables and `init` in the importer?
5. How do you decide between a returned error, a panic, and a closure that stores an error?

#### Easy practical tasks

1. Write `repeat(s string, n int) (string, error)` that rejects `n < 0`. Use a named result for the error.
2. Write a variadic `printAll(...any)` that prints each argument on its own line.
3. Write `applyTwice(fn func(int) int, x int) int`. Pass an anonymous function that adds one.
4. Draw a table with columns Function form, When you use it, and One risk. Add rows for named returns, closures, `init`, and panic.

#### Medium practical tasks

1. Write a recursive function and a function-value helper that solve the same small job. State which form you keep.
2. Build a tiny plugin map: name to `func([]int) int`. Add `sum` and `max`. Reject unknown names with an error.
3. Write a package with `init` that you then remove. Move the same setup into a function that tests can skip.

#### Advanced practical tasks

1. Write a function that uses a named `error`, a `defer` with `recover`, and a helper that may panic. Return a single `error` to `main`.
2. Design five functions for a line-based calculator. List signatures only. Mark which signatures use multiple results, variadic parameters, or function values.
