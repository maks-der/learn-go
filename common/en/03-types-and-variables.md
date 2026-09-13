# 3. Types and Variables

## Description

Go is a statically typed language. Every variable and every value has a type. The compiler checks types before the program runs. A type error stops the build.

This topic explains the basic types and how you declare values. You learn `bool`, `string`, integers, and floats. You learn constants, `var`, short declaration, and zero values. You also learn type conversion, `iota`, and pointers.

Go does not convert types in silence. You write a conversion when types differ. Pointers hold addresses. The zero value of a pointer is `nil`.

Use one term for each concept. A type names a set of values and operations. A variable holds a value of one type. A constant is a named value that does not change. A pointer holds the address of a variable.

---

## Primitive types: `bool`, `string`, integers, floats

A primitive type is a predeclared type for a single value. The main primitive types are `bool`, `string`, integer types, and float types. Go also has complex types. This section focuses on the types that you use first.

`bool` has two values: `true` and `false`. The zero value is `false`. Use `bool` for conditions and flags.

`string` is an immutable sequence of bytes. The zero value is `""`. Go source encodes string literals as UTF-8. `len(s)` counts bytes, not characters. A `for range` over a string yields runes. A rune is a Unicode code point. The type `rune` is an alias of `int32`. The type `byte` is an alias of `uint8`.

Integer types have a fixed size or a platform size. The sized types are `int8`, `int16`, `int32`, `int64`, `uint8`, `uint16`, `uint32`, and `uint64`. `int` and `uint` match the architecture. On a 64-bit system they are 64 bits. On a 32-bit system they are 32 bits. `uintptr` is an unsigned integer that can hold a pointer. Use `int` for most counts and indexes.

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

## Typed vs untyped constants

A constant is a value that the compiler knows. You declare a constant with `const`. You cannot assign a new value to a constant at run time.

```go
const Max = 100
const Title = "demo"
const Ready = true
```

A constant is typed or untyped. A typed constant has a named type. An untyped constant has a kind, such as integer, float, string, or boolean. An untyped constant does not yet have a Go type.

```go
const typedMax int = 100
const untypedMax = 100
```

An untyped numeric constant has high precision in the compiler. You can use it in any numeric context that can hold the value. `untypedMax` can initialize `int`, `int64`, or `float64` when the value fits.

A typed constant must match the target type. `typedMax` has type `int`. You cannot assign it to an `int64` variable without a conversion.

When an untyped constant becomes a variable value, it takes a default type. The default type of an untyped integer is `int`. The default type of an untyped float is `float64`. The default type of an untyped string is `string`. The default type of an untyped boolean is `bool`.

```go
n := 1      // n is int
x := 1.0    // x is float64
s := "hi"   // s is string
```

Use untyped constants for numbers that many types may accept. Use a typed constant when the type is part of the API.

### Questions

#### Theoretical questions

1. What keyword declares a constant?
2. What is an untyped constant?
3. What is the default type of an untyped integer constant?
4. Why can an untyped integer initialize both `int32` and `int64`?
5. Can you change a `const` value after declaration?

#### Easy practical tasks

1. Declare three untyped constants: an integer, a string, and a boolean. Print them.
2. Declare `const Limit int32 = 10`. Assign `Limit` to an `int32` variable. Try to assign it to `int64` and record the result.
3. Write `const Pi = 3.14159` and use it as `float32` and as `float64`. Print both.
4. Try `const N = 1; N = 2` in a function. Record the compiler error.

#### Medium practical tasks

1. Create one typed constant and one untyped constant with the same number. Assign each to `int`, `int64`, and `float64`. Record which assignments compile.
2. Declare a large untyped integer that fits in `int64` but not in `int32`. Assign it to both types. Record the results.
3. Use an untyped constant in a short declaration. Print the type with `%T`. Then give the constant a type and repeat.

#### Advanced practical tasks

1. Write a small report that explains why `const C = 1e20` can initialize `float64` but may fail for `int64`. Show the compiler messages.
2. Build a package API with two typed constants and two untyped constants. Call the API from `main` with several target types. Document each legal use.

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

Package-level variables exist for the life of the process. Local variables exist until no code can use them. The compiler may place them on the stack or on the heap.

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

## Type conversion (no implicit conversion)

Go does not convert types in silence. If `x` has type `T1` and you need type `T2`, you write a conversion. The form is `T2(x)`.

```go
var n int = 5
var m int64 = int64(n)
var f float64 = float64(n)
```

A conversion is legal only when the types can convert. Numeric types can convert to other numeric types. The conversion may truncate or overflow. `int64` to `int32` can lose high bits. `float64` to `int` drops the fraction.

```go
var x float64 = 3.9
var n int = int(x) // n is 3
```

`string` and `[]byte` convert to each other. `string` and `[]rune` convert to each other. `string(r)` for a rune builds a UTF-8 string of that code point. `string(99)` is `"c"` when `99` is a rune or an integer used as a rune conversion.

You cannot convert a `bool` to an integer. You cannot add `int` and `int32` without a conversion. Two named types with the same underlying type convert to each other.

```go
type Celsius float64
type Fahrenheit float64

var c Celsius = 20
var f Fahrenheit = Fahrenheit(c) // legal, but the number is not a real formula
```

A conversion does not change the value in place. The conversion produces a new value. Convert at the edge of an API. Keep one type inside a function when you can.

### Questions

#### Theoretical questions

1. Does Go convert `int` to `int64` without a conversion expression?
2. What is the syntax of a type conversion?
3. What happens when you convert `float64` `3.9` to `int`?
4. Can you convert `bool` to `int`?
5. When can two named types convert to each other?

#### Easy practical tasks

1. Convert an `int` to `int64` and print both values.
2. Convert `3.7` to `int` and print the result.
3. Convert a string to `[]byte`. Change one byte. Convert back to `string` and print both strings.
4. Try `var n int = int64(1)` without a matching conversion. Record the error.

#### Medium practical tasks

1. Write a function that accepts `int64` and call it with an `int` value. Add the conversion at the call site.
2. Convert a large `int64` to `int32` so that the value does not fit. Print both values. Explain the change.
3. Define two named types with underlying type `int`. Convert one to the other. Try to add them without a conversion and record the error.

#### Advanced practical tasks

1. Convert between `string`, `[]byte`, and `[]rune` for a string that contains `€`. Print lengths at each step. Explain each length.
2. Write a small unit converter that uses named types. Allow conversion only through functions that apply a formula. Do not use a raw type conversion in `main`.

---

## `iota` and constant groups

`iota` is a predeclared identifier for constant declarations. In a `const` group, `iota` starts at `0`. `iota` grows by one for each new constant line in that group.

```go
const (
	Sunday = iota // 0
	Monday        // 1
	Tuesday       // 2
)
```

A line that repeats the previous expression still uses the new `iota`. You may write an expression with `iota`:

```go
const (
	Read = 1 << iota // 1
	Write            // 2
	Exec             // 4
)
```

Use `_` to skip a value:

```go
const (
	_ = iota // skip 0
	One
	Two
)
```

A new `const` block resets `iota` to `0`. `iota` does not continue across blocks.

A constant group can mix `iota` and other expressions. Read each line. The compiler copies the expression from the last explicit line when a name has no expression.

```go
const (
	A = iota * 2 // 0
	B            // 2
	C = 100      // 100
	D            // 100
	E = iota     // 4
)
```

Use `iota` for related integer codes. Give the group a type when the codes form an API:

```go
type Day int

const (
	Sunday Day = iota
	Monday
)
```

Do not use `iota` for values that must stay stable in a stored file unless you control every insert. A new line in the middle changes later numbers.

### Questions

#### Theoretical questions

1. What is the first value of `iota` in a `const` group?
2. When does `iota` reset to `0`?
3. How do you skip one `iota` value?
4. What happens on a line that has a name and no expression?
5. Why can a new line in the middle of an `iota` group change later values?

#### Easy practical tasks

1. Declare a `const` group of four weekday names with `iota`. Print each name and value.
2. Use `1 << iota` for three flag constants. Print them in decimal and binary.
3. Skip the zero value with `_`. Print the next two names and values.
4. Write two `const` groups that both use `iota`. Print the first value of each group.

#### Medium practical tasks

1. Give the group a named type `Level`. Declare `Low`, `Medium`, and `High`. Print `%T` and the values.
2. Build a group where one line sets `= 100` and the next lines return to `iota`. Print every value. Explain each number.
3. Insert a new name in the middle of an existing `iota` group. Show how later values change. Then restore the old numbers with explicit values.

#### Advanced practical tasks

1. Model file-permission bits with `iota` and shifts. Write functions that test and combine bits. Print a short demo.
2. Design version codes that must not change after release. Show one `iota` design and one explicit-number design. Write which design you keep and why.

---

## Pointers: `&`, `*`, nil pointers

A pointer holds the address of a variable. The type of a pointer to `T` is `*T`. The operator `&` takes the address of a variable. The operator `*` reads or writes the value at that address.

```go
n := 10
p := &n
fmt.Println(*p) // 10
*p = 20
fmt.Println(n) // 20
```

The zero value of a pointer is `nil`. A `nil` pointer does not point to a variable. A read or write through a `nil` pointer causes a panic.

```go
var p *int
fmt.Println(p == nil) // true
```

Do not dereference a pointer before you know it is not `nil`. Compare the pointer with `nil` first.

Go has no pointer arithmetic. You cannot add `1` to a pointer. You pass a pointer when a function must change the caller variable. You also pass a pointer to avoid a copy of a large struct.

```go
func inc(n *int) {
	*n++
}
```

A pointer value compares equal to another pointer when both point to the same variable, or when both are `nil`.

You may take the address of a struct field. You may store a pointer in a struct field. A pointer to a pointer has type `**T`. That form is rare.

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

1. Dereference a `nil` `*int` inside a function. Recover nothing. Record the panic text. Then add a `nil` check and return.
2. Write two functions that add `1` to an `int`: one by value and one by pointer. Show which change the caller sees.
3. Store a pointer to a struct field. Change the field through the pointer. Print the struct.

#### Advanced practical tasks

1. Return the address of a local variable from a function. Print the value in `main`. Write one sentence about where the variable lives.
2. Build a small linked list with a `node` struct that holds `*node`. Insert three values. Print the list. Handle a `nil` head.

---

## `new` vs taking an address of a value

`new` is a built-in function. `new(T)` allocates a variable of type `T`. The variable holds the zero value. `new` returns `*T`.

```go
p := new(int)
fmt.Println(*p) // 0
*p = 7
```

Taking an address uses `&`. You can write `&x` when `x` is a variable. You can write `&T{...}` for a composite literal. You cannot write `&42` for an integer literal.

```go
x := 7
p1 := &x

type Point struct{ X, Y int }
p2 := &Point{X: 1, Y: 2}
p3 := new(Point)
```

`new(Point)` and `&Point{}` both yield `*Point` to a zero `Point`. Prefer `&Point{}` when you set fields. Prefer a named variable and `&x` when you already have a value. `new` is useful for a zero `*int` or `*bool` without an extra line.

Do not confuse `new` with `make`. `make` creates slices, maps, and channels. `make` returns a ready value, not a pointer to a zero header. This topic does not replace the later slice and map topic.

```go
p := new([]int) // p is *[]int, and *p is nil
s := make([]int, 0) // s is []int, and s is empty
```

Use the form that a reader understands first. Many Go programs never call `new`. Both forms are correct when they produce the pointer that you need.

### Questions

#### Theoretical questions

1. What does `new(T)` return?
2. What value sits in the variable that `new` allocates?
3. Why is `&42` illegal?
4. How do `new(Point)` and `&Point{}` compare for a struct?
5. How does `new` differ from `make`?

#### Easy practical tasks

1. Use `new(int)` and set `*p` to `5`. Print `*p`.
2. Take the address of a local `int` and print the pointer and the value.
3. Use `new` on a struct with two fields. Print the zero fields. Then set the fields.
4. Try `p := &10`. Record the compiler error. Replace it with a legal form.

#### Medium practical tasks

1. Write one function that returns `new(T)` and one function that returns `&T{}` for the same struct. Print both results. State when you prefer each form.
2. Compare `new([]int)` and `make([]int, 0)`. Print types and `== nil` for the slice values.
3. Allocate `*bool` with `new`. Set the value to `true` through the pointer. Pass the pointer to a function that prints it.

#### Advanced practical tasks

1. Benchmark or time a loop that uses `new(Point)` against a loop that uses `&Point{}`. Write the method and the two results. Do not claim a winner without numbers.
2. Design an API that returns `*Config` with defaults. Implement it once with `new` and once with a composite literal. Document which implementation you keep.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do zero values, `nil` pointers, and untyped constants differ when you declare a name with no explicit type?
2. Why must a conversion appear when you pass `int` to a function that wants `int64`, even if the number is small?
3. When do you choose `var`, `:=`, and `const` for a numeric name?
4. How do `iota` and a named integer type work together as a small API?
5. When does a function need a pointer parameter instead of a value parameter?

#### Easy practical tasks

1. Write a program that declares a `bool`, a `string`, an `int`, a `float64`, a `const` count, and a `*int`. Print types and values.
2. Convert the `int` to `float64`, take its address, and change it through the pointer. Print the variable before and after.
3. Declare a typed constant `Max int = 10` and an untyped constant `Base = 2`. Use both in one arithmetic expression that you store in `int64`.
4. Build a `const` group with `iota` for three sizes. Store the selected size in a variable and print it.

#### Medium practical tasks

1. Write a `Scale(p *float64, factor float64)` function. Reject a `nil` pointer. Convert an `int` factor at the call site. Show a zero value, a normal value, and a `nil` call.
2. Create named types `Grams` and `Kilograms`. Convert through a function, not in `main`. Use `new` or `&` to store a result that a helper can update.
3. Mix `var` zero values, `:=`, and a `const` group in one package. Add a test file that checks the zero `bool`, the default type of `1`, and one `iota` value.

#### Advanced practical tasks

1. Implement a tiny thermometer type with `iota` scales, pointer receivers that change a value, and conversions between `float64` and your named type. Do not hide a `nil` dereference.
2. Write a report of ten lines that maps each bug you hit in this topic to one rule: types, constants, zero values, conversion, `iota`, or pointers. Include the compiler or panic text.
