# 4. Operators and Control Flow

## Description

Operators combine values. Control flow chooses the next statement. Go keeps both parts small. You write expressions with operators. You write `if`, `switch`, and `for` for decisions and loops.

This topic explains arithmetic, comparison, and logical operators. You learn `if` with and without a short statement. You learn value switch and type switch. You learn the three `for` forms. You also learn `break`, `continue`, labels, `goto`, and `defer`.

Go has no `while` keyword. `for` is the only loop. Go has `goto`. Do not use `goto` in new code. Use `defer` to schedule a call that runs when the function returns.

Use one term for each concept. An operator produces a value from operands. A statement changes control or assigns a value. A label names a statement for `break`, `continue`, or `goto`. A deferred call runs after the surrounding function ends.

---

## Arithmetic, comparison, and logical operators

Arithmetic operators work on numbers. The operators are `+`, `-`, `*`, `/`, and `%`. `%` is the remainder. `%` works on integers only. Integer `/` truncates toward zero. `+` also concatenates strings.

```go
sum := 3 + 4
mod := 10 % 3
name := "Go" + "1.22"
```

`++` and `--` add or subtract one. They are statements. They are not expressions. You cannot write `x = i++`. You write `i++` on its own line.

Bitwise operators work on integers. The operators are `&`, `|`, `^`, `&^`, `<<`, and `>>`. `^` is xor when it has two operands. `^` is bitwise complement when it has one operand. `&^` is bit clear.

Comparison operators are `==`, `!=`, `<`, `<=`, `>`, and `>=`. The result has type `bool`. You can compare numbers, strings, and other comparable types. You cannot compare slices, maps, or functions with each other. You can compare those types with `nil`.

Logical operators are `&&`, `||`, and `!`. `&&` and `||` use short-circuit evaluation. If the left side of `&&` is `false`, Go does not evaluate the right side. If the left side of `||` is `true`, Go does not evaluate the right side.

```go
ok := n > 0 && name != ""
```

Assignment operators include `=`, `+=`, `-=`, `*=`, and `/=`. A comparison is not an assignment. Use `==` to compare. Use `=` to assign.

Operator precedence controls the order of work. `*` and `/` bind more tightly than `+` and `-`. `&&` binds more tightly than `||`. Use parentheses when the order is not obvious.

Go 1.21 and later provide built-in `min` and `max` for ordered values. Use them when you need the smaller or larger of two or more values.

### Questions

#### Theoretical questions

1. What is the result type of `3 > 1`?
2. Why is `i++` not legal inside a larger expression?
3. What does short-circuit evaluation mean for `&&`?
4. Can you compare two slices with `==`?
5. What does integer division `7 / 2` produce?

#### Easy practical tasks

1. Print `7+2`, `7-2`, `7*2`, `7/2`, and `7%2` for `int` values.
2. Concatenate two strings with `+` and print the result.
3. Write `i++` and `i--` in a function. Print `i` after each statement.
4. Print the results of `true && false`, `true || false`, and `!true`.

#### Medium practical tasks

1. Show short-circuit behavior. Call a function that prints `right` only on the right side of `&&` and `||`. Use one true left side and one false left side.
2. Compare two structs with `==` when all fields are comparable. Then add a slice field and record the compiler error.
3. Use bitwise `&`, `|`, and `<<` to set and test a flag bit. Print the value before and after.

#### Advanced practical tasks

1. Write a function that reports overflow for `uint8` addition without using a wider type in the public result. Print cases that wrap and cases that do not.
2. Build a small expression evaluator for `+`, `-`, `*`, and `/` on `int` tokens. Reject `%` on a non-integer path. Print four sample expressions.

---

## `if` / `else` (including `if` with a short statement)

`if` runs a block when a condition is `true`. The condition must have type `bool`. Parentheses around the condition are optional. Braces around the block are required.

```go
if n > 0 {
	fmt.Println("positive")
}
```

`else` runs when the condition is `false`. `else if` chains more tests. The `else` keyword must stay on the same line as the closing `}` of the previous block. Semicolon insertion makes a new line before `else` invalid.

```go
if n > 0 {
	fmt.Println("positive")
} else if n == 0 {
	fmt.Println("zero")
} else {
	fmt.Println("negative")
}
```

`if` can start with a short statement. The statement runs first. The condition runs next. The names from the short statement exist in the `if` block and in the `else` blocks.

```go
if v := len(name); v > 0 {
	fmt.Println(v)
} else {
	fmt.Println("empty")
}
```

Use the short statement for a local value that you need only in the `if`. A common form is `if err := work(); err != nil`. Do not use the short statement when the value must stay in the outer function.

There is no ternary operator. Do not write `x = cond ? a : b`. Use `if` and assign in each branch.

### Questions

#### Theoretical questions

1. What type must an `if` condition have?
2. Are braces required for a one-line `if` body?
3. Where can you use names from an `if` short statement?
4. Why must `else` stay on the same line as `}`?
5. Does Go have a ternary `? :` operator?

#### Easy practical tasks

1. Write `if` that prints `ok` when an `int` is greater than `10`.
2. Add `else` that prints `small`.
3. Write `if v := 3; v%2 == 1` and print `odd`.
4. Write an `else if` chain for negative, zero, and positive numbers. Test three values.

#### Medium practical tasks

1. Use `if n, err := strconv.Atoi(s); err != nil` to parse a string. Print the error path and the success path.
2. Move `else` to the next line after `}`. Record the compiler error. Restore the layout.
3. Show that a name from the short statement is not visible after the `if`/`else`. Try to print it below the `if` and record the error.

#### Advanced practical tasks

1. Rewrite a nested `if` chain into a flat `else if` chain. Keep the same results for five inputs. Print a before and after table.
2. Write a function that uses a short statement to open a scope for two locals. Return a value from both branches. Add tests for both branches.

---

## `switch` (value and type switch)

A value `switch` compares an expression to a list of cases. Go runs the first matching case. Go does not fall through to the next case unless you write `fallthrough`.

```go
switch day {
case 1:
	fmt.Println("Monday")
case 2, 3:
	fmt.Println("mid")
default:
	fmt.Println("other")
}
```

A case can list several values. `default` runs when no case matches. `default` is optional. Cases evaluate from top to bottom.

A `switch` may omit the expression. That form is `switch true`. Each case is a `bool` expression. This form replaces a long `if` / `else if` chain.

A `switch` may start with a short statement, like `if`:

```go
switch n := len(name); n {
case 0:
	fmt.Println("empty")
default:
	fmt.Println(n)
}
```

A type switch compares the dynamic type of an interface value. The form is `switch v := x.(type)`. Each case names a type. Inside a case, `v` has that type.

```go
func describe(x any) {
	switch v := x.(type) {
	case int:
		fmt.Println("int", v)
	case string:
		fmt.Println("string", v)
	default:
		fmt.Println("other")
	}
}
```

`any` is the alias for `interface{}` from Go 1.18. Use `any` in new code. A type switch cannot use `fallthrough`. Use a type switch when one function accepts several types through an interface.

### Questions

#### Theoretical questions

1. Does a matching `case` continue into the next `case` by default?
2. What does `fallthrough` do?
3. What is a `switch` with no expression?
4. What does `switch v := x.(type)` examine?
5. Can a type switch use `fallthrough`?

#### Easy practical tasks

1. Write a value `switch` on an `int` with three cases and `default`. Print each path.
2. Put two values in one `case` list. Show that both values print the same line.
3. Write a `switch` with no expression that tests `n < 0`, `n == 0`, and `n > 0`.
4. Write a type switch on `any` for `int` and `string`. Call it twice.

#### Medium practical tasks

1. Use `fallthrough` in a value `switch` so that one input prints two lines. Show the order of the lines.
2. Add a short statement to a `switch` that parses a string length. Cover empty and non-empty cases.
3. Type-switch on `any` with `int`, `float64`, and `[]byte`. Print a distinct line for each type.

#### Advanced practical tasks

1. Replace a type switch with a value switch and explain why that rewrite fails or becomes hard to read. Then keep the type switch and add a test per case.
2. Write a dispatcher that uses a value `switch` on a command string and a type switch on a payload of type `any`. Handle unknown command and unknown type.

---

## `for` as the only loop (`for`, `for range`, infinite `for`)

`for` is the only loop statement in Go. There is no `while` and no `do-while`. Three forms cover all loops.

The full form has init, condition, and post:

```go
for i := 0; i < 3; i++ {
	fmt.Println(i)
}
```

The condition-only form repeats while the condition is `true`:

```go
n := 3
for n > 0 {
	n--
}
```

The infinite form has no condition. The loop runs until `break`, `return`, or `panic`:

```go
for {
	break
}
```

`for range` walks an array, a slice, a string, a map, or a channel. For a slice, the first value is the index and the second value is a copy of the element.

```go
for i, v := range []int{10, 20} {
	fmt.Println(i, v)
}
```

You may omit values. `for i := range s` binds the index. `for _, v := range s` binds the element. `for range s` binds nothing.

From Go 1.22, `for i := range n` walks integers from `0` to `n-1` when `n` is an integer. From Go 1.22, each iteration has its own loop variable. A closure that captures `i` sees that iteration value.

```go
for i := range 3 {
	fmt.Println(i) // 0, 1, 2
}
```

Range over a string yields a byte index and a rune. Range over a map yields key and value in a random order. Range over a channel receives until the channel is closed.

Do not change a slice or map in a way that hides the range rules until you learn those types in a later topic. Keep first loops simple.

### Questions

#### Theoretical questions

1. Why does Go have no `while` keyword?
2. What are the three parts of the full `for` header?
3. What does `for range` over a slice yield?
4. What integers does `for i := range 3` produce in Go 1.22?
5. How did loop-variable capture change in Go 1.22?

#### Easy practical tasks

1. Print numbers `0` to `4` with a full `for`.
2. Print the same numbers with `for i := range 5`.
3. Range over a slice of three names and print index and name.
4. Write an infinite `for` that `break`s after it prints one line.

#### Medium practical tasks

1. Range over a string that contains `ü`. Print each index and rune.
2. Use a condition-only `for` to remove digits from an integer by dividing by `10` until the value is `0`. Print the count of digits.
3. Collect three functions in a slice during a `for`. Each function prints the loop index. Call them after the loop on Go 1.22 or later. Record the three numbers.

#### Advanced practical tasks

1. Compare a full `for` and `for range` on the same slice while you append during the loop. Print indexes and values. Write the rule that you observe.
2. Walk a map with `for range` three times. Show that the print order can change. Then collect keys, sort them, and print a stable order.

---

## `break`, `continue`, labeled breaks

`break` leaves the inner `for`, `switch`, or `select`. The statement after that block runs next.

`continue` ends the current `for` iteration. The loop starts the next iteration. `continue` does not apply to `switch`.

```go
for i := 0; i < 5; i++ {
	if i%2 == 0 {
		continue
	}
	if i == 3 {
		break
	}
	fmt.Println(i)
}
```

A `break` inside a `switch` that sits in a loop leaves the `switch` only. The loop continues. Use a label when you must leave the loop from inside the `switch`.

A label is an identifier and a colon before a statement. `break Label` leaves the labeled `for`, `switch`, or `select`. `continue Label` starts the next iteration of the labeled `for`.

```go
Outer:
	for i := 0; i < 3; i++ {
		switch i {
		case 1:
			break Outer
		}
		fmt.Println(i)
	}
```

Place the label on the loop that you want to control. Use a clear name such as `Outer` or `Rows`. Do not use a label when a simple `return` or a boolean flag is enough.

`fallthrough` is not `break`. `fallthrough` continues a value `switch` into the next case. `break` stops the `switch`.

### Questions

#### Theoretical questions

1. What does `break` leave when it has no label?
2. What does `continue` skip?
3. Why does `break` inside a `switch` not stop the outer `for`?
4. What is a labeled `break`?
5. How does `continue` with a label differ from `break` with a label?

#### Easy practical tasks

1. Print numbers `0` to `9` but `continue` on even numbers.
2. Break a loop when a value equals `7`. Print the values before `7`.
3. Put `break` in a `switch` inside a `for`. Show that the loop continues.
4. Add a label and `break` the loop from inside that `switch`. Show that the loop stops.

#### Medium practical tasks

1. Write a nested loop over rows and columns. Use `continue` on the inner loop to skip one column. Print the pairs.
2. Use a labeled `continue` to skip a full row when the row index is `2`. Print the remaining pairs.
3. Replace a labeled `break` with a `return` from a helper function. Keep the same printed result.

#### Advanced practical tasks

1. Search a 2D slice for the first negative number. Leave both loops with a labeled `break`. Return the row, column, and a found flag.
2. Rewrite the same search with no label. Use a named function and `return`. Compare the two versions in six short sentences.

---

## `goto` (know it exists; almost never use it)

`goto` jumps to a labeled statement in the same function. The language includes `goto`. The Go team does not use `goto` as a normal control tool.

```go
func demo(n int) {
	if n < 0 {
		goto Done
	}
	fmt.Println(n)
Done:
	fmt.Println("end")
}
```

`goto` must not jump over the declaration of a variable into the scope of that variable. The compiler rejects that jump. `goto` must not jump into another block in a way that skips initialization.

Do not use `goto` to build loops. Use `for`. Do not use `goto` to exit nested blocks in new code. Use `return`, `break`, or a labeled `break`. Do not use `goto` to handle errors. Use extra `if` lines or a helper function.

You learn `goto` so that you can read old code and compiler messages. If you see `goto` in a new change, replace it with structured control flow.

A label that `goto` uses is the same form of name as a `break` label. Keep labels rare. Two uses of the same label name in one function are an error.

### Questions

#### Theoretical questions

1. What does `goto` do?
2. Why must you not use `goto` to write a loop?
3. What jump over a variable declaration does the compiler reject?
4. What structured tools replace `goto` for nested loops?
5. Is `goto` legal across two functions?

#### Easy practical tasks

1. Write a function with one `goto` and one label. Print a line before the jump and a line at the label.
2. Replace that `goto` with an `if`. Keep the same output.
3. Try `goto` into a block that declares a new variable. Record the compiler error.
4. List three statements from this topic that you must use instead of `goto`.

#### Medium practical tasks

1. Find a `goto` example in older public Go code or in a tutorial. Rewrite the function without `goto`. Show both versions.
2. Use `goto` to skip one print, then rewrite the skip with `if`. Run both functions and compare the output.
3. Try `goto` to a label in another function. Record the result. Explain the scope of a label.

#### Advanced practical tasks

1. Write a nested-loop search with `goto` for exit. Then delete `goto` and use a labeled `break`. Then use a helper and `return`. Measure only clarity. Write which version a beginner must keep.
2. Read the `goto` rules in the language specification. Write eight short sentences that restate the rules. Do not copy the specification text.

---

## `defer` basics

`defer` schedules a function call. The call runs when the surrounding function returns. The call also runs when the function panics, before the panic continues.

```go
func demo() {
	defer fmt.Println("after")
	fmt.Println("before")
}
```

This program prints `before` and then `after`.

Arguments of a deferred call evaluate at the `defer` line. The call waits. The argument values do not wait.

```go
n := 1
defer fmt.Println(n) // prints 1
n = 2
```

Several `defer` calls run in last-in, first-out order. The last `defer` runs first.

```go
defer fmt.Println("first")
defer fmt.Println("second")
```

This program prints `second` and then `first`.

Use `defer` to close a file, unlock a mutex, or restore a setting. Place `defer` soon after a successful open or lock. The close then runs on every return path.

`os.Exit` stops the process and does not run deferred calls. A `return` from `main` does run deferred calls.

A `defer` in a loop schedules one call per iteration. Those calls wait until the function returns. Do not defer `Close` in a long loop when you can close in the same iteration.

Named result parameters can change in a deferred function. A later topic covers that pattern. In this topic, use `defer` for simple cleanup and for print order.

### Questions

#### Theoretical questions

1. When does a deferred call run?
2. When does Go evaluate the arguments of a deferred call?
3. In what order do several `defer` calls run?
4. Do deferred calls run after `os.Exit`?
5. What problem appears when you `defer` `Close` inside a long loop?

#### Easy practical tasks

1. Write a function that prints `start`, then `defer` prints `end`. Run it.
2. Defer two prints. Show the last-in, first-out order.
3. Pass a variable to `fmt.Println` in `defer`. Change the variable after `defer`. Print what the deferred call shows.
4. Defer a call that prints `done` and then `return` in the middle of `main`. Confirm that `done` still prints.

#### Medium practical tasks

1. Compare `return` from `main` with `os.Exit(0)` while a `defer` print exists. Record which run shows the print.
2. Open a file, `defer` `Close`, write one line, and read the file back in a second function.
3. Put `defer fmt.Println(i)` inside a `for` that runs three times. Print the order of the numbers after the loop body.

#### Advanced practical tasks

1. Write a function that defers a closure that prints a variable by reference and a `defer` that prints the value as an argument. Change the variable before return. Explain the two lines.
2. Simulate cleanup of three resources with `defer`. Fail on the second resource. Show that earlier `defer` calls still run. Do not leak a fake resource.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do operators, `if`, and `switch` share the type `bool`?
2. When do you choose `for range` instead of a full `for`, and when do you choose an infinite `for`?
3. How do labeled `break` and `defer` each change the path out of nested work?
4. Why is `goto` a reading skill and not a writing tool in this handbook?
5. When a function returns from inside a `switch` in a loop, which deferred calls run, and what does the loop do?

#### Easy practical tasks

1. Write a program that uses `%`, `&&`, `if`, and a full `for` to print odd numbers below `10`.
2. Add a value `switch` that labels each number as `small`, `mid`, or `big`. Keep the loop.
3. Range over a slice of words. `continue` on an empty word. `defer` one summary print that runs at the end of `main`.
4. Count spaces in a sentence with `for range` over a string and a comparison. Print the count.

#### Medium practical tasks

1. Parse flags from `os.Args` with `if` short statements and a `switch` on the command word. Include `help` and an unknown-command path.
2. Walk a slice with `for range`. Use a type switch on `any` values in the slice. `break` a labeled loop when you see a negative `int`.
3. Combine `defer` cleanup with an `if err != nil` chain in one function that writes a temp file. Show the success path and the error path.

#### Advanced practical tasks

1. Build a tiny REPL loop: infinite `for`, `switch` on text commands, `continue` on empty input, `break` on `quit`, and `defer` a goodbye print. Reject `goto`.
2. Write a function that processes a matrix with nested loops, a type switch on each cell of type `any`, and deferred row counters. Add tests for empty input, a `nil` pointer cell, and an early labeled `break`.
