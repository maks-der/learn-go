# 2. Program Structure and Types

## Description

A Go source file has a fixed order. The file starts with a package clause. Import declarations come next. Other declarations come after the imports. Every variable and every value has a type. The compiler checks types before the program runs.

This topic explains the file skeleton, names, primitive types, variables, constants, and pointers. Complete this topic before you study control flow.

Use one term for each concept. A package is a directory of Go files with the same package name. A command is a `package main` that contains `func main`. A type names a set of values and operations. A pointer holds the address of a variable.

---

## `package main`, `import`, and `func main`

A package is the unit that the compiler builds. All `.go` files in one directory form one package. Those files must use the same package name. A file starts with a package clause:

```go
package greet
```

The package name is an identifier. The package name is often the same as the folder name. The package name is not the import path. Other files import the package by import path.

`package main` is special. A `package main` defines a command. A command builds to a binary. A library package uses a different name, such as `greet`. A library package does not build to a binary by itself.

An import declaration names a package that this file uses. The import path is a string. For the standard library, the import path is `"fmt"` or `"net/http"`. For a module, the import path is the module path plus the package folder.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("command")
}
```

The compiler rejects an unused import. Remove the import or use a name from that package. `gofmt` puts standard library paths in the first group. `gofmt` puts other paths in the next group.

The program entry point is `func main` in `package main`. The signature has no parameters and no results. The runtime initializes the `main` package. The runtime then calls `func main`. When `main` returns, the process exits with status `0`.

Do not give `main` parameters. Do not give `main` results. Read command-line arguments from `os.Args`. Call `os.Exit` when you need a non-zero status. `os.Exit` stops the process at once. Deferred calls do not run after `os.Exit`.

A `package main` must contain exactly one `func main`. Keep `main` small. Put logic in other functions and packages.

### Questions

#### Theoretical questions

1. What is a package in Go?
2. What does `package main` mark?
3. What is the difference between a package name and an import path?
4. What is the exact signature of the program entry function?
5. How does `os.Exit` differ from a return from `main`?

#### Easy practical tasks

1. Create a module. Add `package main` and print `hello`. Run the program with `go run .`.
2. Add a folder `msg` with `package msg`. Export one function. Call it from `main`.
3. Add an unused import. Run `go build`. Record the error. Remove the import.
4. Print `len(os.Args)` and each argument. Run the binary with two extra words.

#### Medium practical tasks

1. Split one command into three files in the same `package main`. Keep one `func main`. Build the command.
2. Create packages `a` and `b` that import each other. Record the cycle error. Break the cycle with a third package.
3. Write `main` that returns status `1` with `os.Exit` when no argument exists. Return status `0` when an argument exists. Show both runs.

#### Advanced practical tasks

1. Design a module with one library package and two commands that import it. Write the tree, package names, and import paths. Build both commands.
2. Compare `os.Exit(1)` with `panic("fail")` from `main`. Record process status and whether a `defer` in `main` runs.

---

## Comments, identifiers, keywords, and exported names

A comment is text that the compiler ignores. A line comment starts with `//` and ends at the end of the line. A block comment starts with `/*` and ends with `*/`. Use `//` for almost all comments. Block comments do not nest.

```go
// Package greet formats short messages.
package greet

// Hello returns a greeting for name.
func Hello(name string) string {
	return "hi " + name
}
```

A package comment sits above the package clause with no blank line. A function comment sits above the function. Start a function comment with the function name.

An identifier names a package, function, type, variable, constant, or label. An identifier starts with a letter or `_`. The next characters are letters, digits, or `_`. These names are valid: `total`, `Total`, `HTTPClient`. These names are not valid: `2x`, `user-name`.

A keyword is a reserved word. You must not use a keyword as an identifier. The keywords include `package`, `import`, `func`, `var`, `const`, `type`, `if`, `for`, `return`, `go`, `defer`, and `struct`.

The blank identifier `_` discards a value. Use `_` when a function returns a value that you do not need.

Name packages with short lowercase words. Name other identifiers with MixedCaps. Start an exported name with an uppercase letter. Start an unexported name with a lowercase letter. Keep acronyms in full caps in exported names, such as `HTTP`.

A name is exported when the first letter is an uppercase Unicode letter. Other packages can use an exported name. A name is unexported when the first letter is a lowercase letter or `_`. Other packages cannot use an unexported name. The rule applies to functions, types, variables, constants, fields, and methods.

```go
package msg

var Default = "hi" // exported

func Hello(name string) string {
	return wrap(name)
}

func wrap(name string) string { // unexported
	return Default + " " + name
}
```

The compiler also inserts semicolons. The opening `{` of `func`, `if`, and `for` must stay on the same line. Do not write comments that only repeat the code.

### Questions

#### Theoretical questions

1. Where does a `//` comment end?
2. What makes a name exported?
3. Why is `package` not a valid identifier?
4. What is the blank identifier?
5. Can another package call an unexported function?

#### Easy practical tasks

1. Add a package comment and a function comment to a small command. Run `go doc` in that folder.
2. Try to name a variable `func`. Record the compiler error. Rename the variable.
3. Export `Hello` from package `greet`. Call `greet.Hello` from `main`. Run the program.
4. Change `Hello` to `hello`. Build `main`. Record the error. Restore the export.

#### Medium practical tasks

1. Write comments for two exported functions so that `go doc -all` shows both comments.
2. Create a struct with one exported field and one unexported field. Set both fields inside the package. Try to set the unexported field from `main` and record the error.
3. Write tests in `package greet` that call an unexported function. Write tests in `package greet_test` that cannot call it. Show both results.

#### Advanced practical tasks

1. Design a package API with three exported names and four unexported names. Write a one-page note that explains each export choice.
2. Read the semicolon rule in the language specification. Rewrite the rule in ten short sentences of your own. Do not copy the specification text.

---

## Primitive types: `bool`, `string`, integers, floats

A primitive type is a predeclared type for a single value. The main primitive types are `bool`, `string`, integer types, and float types.

`bool` has two values: `true` and `false`. The zero value is `false`. Use `bool` for conditions and flags.

`string` is an immutable sequence of bytes. The zero value is `""`. Go source encodes string literals as UTF-8. `len(s)` counts bytes, not characters. A `for range` over a string yields runes. A rune is a Unicode code point. The type `rune` is an alias of `int32`. The type `byte` is an alias of `uint8`.

Integer types have a fixed size or a platform size. The sized types are `int8`, `int16`, `int32`, `int64`, `uint8`, `uint16`, `uint32`, and `uint64`. `int` and `uint` match the architecture. On a 64-bit system they are 64 bits. Use `int` for most counts and indexes.

Float types are `float32` and `float64`. They follow IEEE-754. Prefer `float64` for new code. The zero value of every numeric type is `0`.

```go
var ok bool = true
var name string = "Go"
var n int = 22
var x float64 = 1.5
```

Do not mix types in an operation. `int` and `int64` are different types. Convert one value first. Integer overflow wraps at run time. A constant overflow is a compile error.

### Questions

#### Theoretical questions

1. What values can a `bool` hold?
2. What does `len` return for a `string`?
3. What is the difference between `int` and `int64`?
4. Which float type do you use for most new code?
5. What is the zero value of `string`?

#### Easy practical tasks

1. Declare one `bool`, one `string`, one `int`, and one `float64`. Print each value with `fmt.Printf` and `%T`.
2. Print `len("Gö")` and then range over the string. Write the byte length and each rune.
3. Assign `int8` to the maximum value `127`. Add `1` and print the result.
4. Print `true`, `false`, and the zero `bool` from a `var` with no value.

#### Medium practical tasks

1. Write a function that counts runes in a string without `utf8.RuneCountInString`. Compare the result with `len`.
2. Store the same number in `int32` and `int64`. Try to add them without a conversion. Record the error. Fix the add.
3. Compare `float32` and `float64` for `0.1 + 0.2 == 0.3`. Print both results. Write one sentence about equality of floats.

#### Advanced practical tasks

1. Write a table that lists each integer type, its size in bits, and its minimum and maximum values. Print the table from a program that uses `math` limits.
2. Decode a string that contains mixed ASCII and non-ASCII letters. Print each byte index and each rune start index. Explain the difference.

---

## `var`, short declaration `:=`, and zero values

A variable holds a value that can change. The zero value is the value of a variable that you do not initialize. Every type has a zero value.

`var` declares a variable. You may give a type, a value, or both:

```go
var a int
var b int = 3
var c = 3
```

When you omit the value, Go stores the zero value. The zero value of `bool` is `false`. The zero value of `string` is `""`. The zero value of every numeric type is `0`. The zero value of a pointer is `nil`.

You may group `var` declarations:

```go
var (
	name string
	age  int
)
```

Short declaration uses `:=`. Short declaration works only inside a function. The compiler infers the type from the value. At least one name on the left must be new in that block.

```go
func demo() {
	count := 0
	count, err := next() // err is new; count is redeclared
	_ = count
	_ = err
}
```

Redeclaration with `:=` requires the same type. Redeclaration does not create a second variable. It assigns to the existing variable.

Use `var` at package level. Use `var` when you want the zero value and no right-hand side. Use `:=` for most local variables that have a value at once. Do not use `:=` when the inferred type is wrong. Write `var x int64 = 0` in that case.

### Questions

#### Theoretical questions

1. What is a zero value?
2. Where is `:=` legal?
3. What does `var n int` store before any assignment?
4. What is the redeclare rule for `:=`?
5. Why do package-level variables use `var` and not `:=`?

#### Easy practical tasks

1. Declare `var` variables of type `bool`, `string`, `int`, and `float64` with no values. Print them.
2. Use `:=` to declare a local `name` and print it.
3. Try `:=` at package level. Record the compiler error.
4. Write `var n = 5` and `k := 5` in a function. Print both types with `%T`.

#### Medium practical tasks

1. Call a function that returns two values. Use `:=` so that one name is new and one name already exists. Print both values.
2. Try to redeclare a name with `:=` and a different type in the same block. Record the error. Fix the types.
3. Compare `var buf []byte` and `buf := []byte{}`. Print whether each value equals `nil`.

#### Advanced practical tasks

1. Write a function that needs an `int64` zero. Show one version with `var` and one version with `:=` and a conversion. Explain which version is clearer.
2. Trace a package-level `var` and a local `var` with prints in `init` and in `main`. Document the order of initialization.

---

## Constants, `iota`, and type conversion

A constant is a value that the compiler knows. You declare a constant with `const`. You cannot assign a new value to a constant at run time.

A constant is typed or untyped. A typed constant has a named type. An untyped constant has a kind, such as integer, float, string, or boolean. An untyped numeric constant has high precision in the compiler. You can use it in any numeric context that can hold the value.

```go
const Max = 100
const typedMax int = 100
```

When an untyped constant becomes a variable value, it takes a default type. The default type of an untyped integer is `int`. The default type of an untyped float is `float64`.

`iota` is a predeclared identifier for constant declarations. In a `const` group, `iota` starts at `0`. `iota` grows by one for each new constant line in that group. A new `const` block resets `iota` to `0`.

```go
const (
	Sunday = iota // 0
	Monday        // 1
	Tuesday       // 2
)

const (
	Read = 1 << iota // 1
	Write            // 2
	Exec             // 4
)
```

Use `_` to skip a value. Give the group a type when the codes form an API. Do not use `iota` for values that must stay stable in a stored file unless you control every insert.

Go does not convert types in silence. If `x` has type `T1` and you need type `T2`, you write `T2(x)`. Numeric types can convert to other numeric types. The conversion may truncate or overflow. You cannot convert a `bool` to an integer. You cannot add `int` and `int32` without a conversion.

```go
var n int = 5
var m int64 = int64(n)
var x float64 = 3.9
var k int = int(x) // k is 3
```

`string` and `[]byte` convert to each other. Convert at the edge of an API. Keep one type inside a function when you can.

### Questions

#### Theoretical questions

1. What is an untyped constant?
2. What is the first value of `iota` in a `const` group?
3. Does Go convert `int` to `int64` without a conversion expression?
4. What happens when you convert `float64` `3.9` to `int`?
5. When does `iota` reset to `0`?

#### Easy practical tasks

1. Declare three untyped constants: an integer, a string, and a boolean. Print them.
2. Declare a `const` group of four weekday names with `iota`. Print each name and value.
3. Convert an `int` to `int64` and print both values.
4. Convert `3.7` to `int` and print the result.

#### Medium practical tasks

1. Create one typed constant and one untyped constant with the same number. Assign each to `int`, `int64`, and `float64`. Record which assignments compile.
2. Use `1 << iota` for three flag constants. Print them in decimal and binary.
3. Define two named types with underlying type `int`. Convert one to the other. Try to add them without a conversion and record the error.

#### Advanced practical tasks

1. Write a small report that explains why `const C = 1e20` can initialize `float64` but may fail for `int64`. Show the compiler messages.
2. Convert between `string`, `[]byte`, and `[]rune` for a string that contains `€`. Print lengths at each step. Explain each length.

---

## Pointers: `&`, `*`, and `nil`

A pointer holds the address of a variable. The type of a pointer to `T` is `*T`. The operator `&` takes the address of a variable. The operator `*` reads or writes the value at that address.

```go
n := 10
p := &n
fmt.Println(*p) // 10
*p = 20
fmt.Println(n) // 20
```

The zero value of a pointer is `nil`. A `nil` pointer does not point to a variable. A read or write through a `nil` pointer causes a panic. Compare the pointer with `nil` first.

Go has no pointer arithmetic. You cannot add `1` to a pointer. You pass a pointer when a function must change the caller variable. You also pass a pointer to avoid a copy of a large struct.

```go
func inc(n *int) {
	*n++
}
```

`new(T)` allocates a variable of type `T` with the zero value and returns `*T`. You can also write `&T{...}` for a composite literal. You cannot write `&42` for an integer literal. Do not confuse `new` with `make`. `make` creates slices, maps, and channels.

The compiler may allocate a local variable on the heap when you return its address. This is escape analysis. You still write `&` and `*`. You do not manage the heap by hand.

### Questions

#### Theoretical questions

1. What does `&` produce?
2. What does `*` do when you apply it to a pointer?
3. What is the zero value of a pointer type?
4. What happens when you dereference a `nil` pointer?
5. Does Go allow pointer arithmetic?

#### Easy practical tasks

1. Take the address of an `int`. Print the pointer with `%p` and the value with `*`.
2. Change the `int` through the pointer. Print the variable after the change.
3. Declare `var p *string` and print `p == nil`.
4. Write a function `inc(n *int)` that adds `1`. Call it from `main`.

#### Medium practical tasks

1. Dereference a `nil` `*int` inside a function. Record the panic text. Then add a `nil` check and return.
2. Write two functions that add `1` to an `int`: one by value and one by pointer. Show which change the caller sees.
3. Use `new(int)` and set `*p` to `5`. Print `*p`. Compare this form with a named variable and `&`.

#### Advanced practical tasks

1. Return the address of a local variable from a function. Print the value in `main`. Write one sentence about where the variable lives.
2. Build a small linked list with a `node` struct that holds `*node`. Insert three values. Print the list. Handle a `nil` head.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the required order of a Go source file from the first line to the first function.
2. How do package names, import paths, and exported names work together when `main` calls another package?
3. When do you choose `var`, `:=`, and `const` for a numeric name?
4. Why must a conversion appear when you pass `int` to a function that wants `int64`?
5. When does a function need a pointer parameter instead of a value parameter?

#### Easy practical tasks

1. Write a three-file command: `main.go`, `help.go` in `package main`, and `text/format.go` as a library. Print one formatted line.
2. Write a program that declares a `bool`, a `string`, an `int`, a `float64`, a `const` count, and a `*int`. Print types and values.
3. Convert the `int` to `float64`, take its address, and change it through the pointer. Print the variable before and after.
4. Build a `const` group with `iota` for three sizes. Store the selected size in a variable and print it.

#### Medium practical tasks

1. Create a small CLI that uses `os.Args`, one imported library package, and `os.Exit`. Cover a missing-argument path and a success path.
2. Write a `Scale(p *float64, factor float64)` function. Reject a `nil` pointer. Convert an `int` factor at the call site.
3. Mix `var` zero values, `:=`, and a `const` group in one package. Add a test file that checks the zero `bool` and one `iota` value.

#### Advanced practical tasks

1. Build a module with `cmd/tool`, `internal/parse`, and exported names only where a client needs them. Document every package name and import path.
2. Implement a tiny thermometer type with `iota` scales, pointer updates, and conversions between `float64` and your named type. Do not hide a `nil` dereference.
