# 14. Generics

## Description

Generics let you write functions and types that work with more than one type. You declare type parameters. The compiler substitutes a concrete type when you call the function or when you instantiate the type. The compiler still checks types at compile time.

This topic shows type parameters, constraints, inference, and design choices. Use generics when the same algorithm applies to many types and you need type safety. Do not use generics to hide simple code behind a complex API.

Use one term for each concept. A type parameter is a name in square brackets. A constraint is an interface that limits the type argument. Instantiation is the step that binds a type argument to a type parameter.

---

## Type parameters on functions and types

A generic function declares type parameters after the function name. The constraint follows each type parameter. The parameter list of the function then uses those names as types.

```go
func First[T any](items []T) (T, bool) {
    if len(items) == 0 {
        var zero T
        return zero, false
    }
    return items[0], true
}
```

`T` is a type parameter. `any` is the constraint. The caller may write `First[int](nums)` or rely on inference from `nums`. The result type is `T`. When the slice is empty, the function returns the zero value of `T`.

A generic type declares type parameters after the type name:

```go
type Pair[A any, B any] struct {
    First  A
    Second B
}

type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    i := len(s.items) - 1
    v := s.items[i]
    s.items = s.items[:i]
    return v, true
}
```

Methods on a generic type may use the type parameters of the type. Methods must not declare new type parameters. This rule is still in force in current Go. If you need a second type parameter on an operation, write a function, not a method.

```go
// Allowed: a function with its own type parameters.
func Convert[T any, U any](s Stack[T], f func(T) U) Stack[U] {
    out := Stack[U]{}
    for _, v := range s.items {
        out.Push(f(v))
    }
    return out
}
```

You instantiate a generic type with a type argument list:

```go
var p Pair[string, int]
s := Stack[int]{}
```

Go 1.24 and later allow generic type aliases. An alias may have type parameters, the same as a defined type:

```go
type Handler[T any] = func(context.Context, T) error
```

Use a defined type when you need methods on that name. Use an alias when you only want a shorter name for an existing type.

The standard library uses type parameters in `slices`, `maps`, and `cmp` (Go 1.21 and later). Prefer those packages before you write your own generic helpers for common slice and map work.

### Questions

#### Theoretical questions

1. Where do you write type parameters on a function?
2. What is instantiation of a generic type?
3. May a method declare a new type parameter that the receiver type does not have?
4. What is the zero value of a type parameter `T` when the function does not know the concrete type?
5. What is the difference between a generic defined type and a generic type alias?
6. Why does `Stack[T]` appear on both the type and the method receiver?

#### Easy practical tasks

1. Write `Last[T any](items []T) (T, bool)`. Return the last element or the zero value.
2. Write `Pair[A, B any]` with a function `MakePair` that returns a `Pair`.
3. Add `Len` and `Pop` methods to `Stack[T]`. Call them from `main` with `Stack[string]`.
4. Instantiate `slices.Index` from the `slices` package for a `[]int`. Print the index of one value.

#### Medium practical tasks

1. Write `MapSlice[T, U any](in []T, f func(T) U) []U`. Test it with `int` to `string`.
2. Write a generic `Queue[T]` with `Enqueue`, `Dequeue`, and `Len`. Keep the slice as the backing store.
3. Write `Swap[T any](a, b *T)`. Show that the function works for `int` and for a struct type.

#### Advanced practical tasks

1. Design a generic binary tree node `Node[T any]` with `Left`, `Right`, and `Value`. Write an in-order walk function that calls a function for each value. Do not use `any` as a stored interface value.
2. Try to add a method with its own type parameter. Record the compiler error. Rewrite the operation as a package-level function. Explain the design limit in four sentences.

---

## Constraints and `comparable`

A constraint is an interface. The type argument must satisfy that interface. The compiler uses the constraint to decide which operations are legal on values of the type parameter.

The built-in constraint `comparable` names every type that supports `==` and `!=`. You need `comparable` when you compare values or when you use the type as a map key.

```go
func Contains[T comparable](items []T, want T) bool {
    for _, v := range items {
        if v == want {
            return true
        }
    }
    return false
}

type Set[T comparable] map[T]struct{}
```

A type set constraint lists types. Use `|` between types. Use `~` to include every type whose underlying type is the named type.

```go
type Signed interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

func Abs[T Signed](v T) T {
    if v < 0 {
        return -v
    }
    return v
}
```

`Abs` may take `int` and also a defined type whose underlying type is `int`:

```go
type Count int
```

Without `~`, a defined type with the same underlying type does not satisfy the constraint.

A constraint may include methods and a type set in one interface:

```go
type StringID interface {
    ~string
    ID() string
}
```

The type argument must have underlying type `string` and must have method `ID() string`.

`cmp.Ordered` in package `cmp` is the standard constraint for types that support `<`, `<=`, `>`, and `>=`. Use `cmp.Ordered` for generic min, max, and sort helpers. The `slices` package uses it.

```go
import "cmp"

func Min[T cmp.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

Go 1.21 also added built-in `min` and `max` for ordered values of the same type. Use the built-ins for simple cases. Write a generic function when you need a reusable helper with extra logic.

Current Go checks operations against the type set of the constraint. The type set must support every operation that the function body uses. If the function uses `<`, the constraint must be ordered. If the function uses `==`, the constraint must be `comparable` or a type set of comparable types.

Interfaces that are not constraints in a type parameter list stay ordinary interfaces. You still use ordinary interfaces for runtime polymorphism. Do not confuse a constraint (compile-time type set) with a value of interface type (runtime method set).

### Questions

#### Theoretical questions

1. What operations does `comparable` allow on a type parameter?
2. What does the `~` prefix mean in a type set?
3. Why can you not use `T` as a map key when the constraint is only `any`?
4. What package provides `Ordered` in the standard library?
5. Can a constraint list both a type set and methods? Give a yes or no answer and one reason.
6. What happens at compile time when the function body uses an operation that the constraint does not allow?

#### Easy practical tasks

1. Write `Index[T comparable](items []T, want T) int`. Return `-1` when the value is missing.
2. Write `Set[T comparable]` with `Add`, `Has`, and `Len`.
3. Write a constraint for unsigned integer types with `~`. Write `IsZero[T YourConstraint](v T) bool`.
4. Call `slices.Min` on a `[]string`. Write why `string` satisfies `cmp.Ordered`.

#### Medium practical tasks

1. Write `Max[T cmp.Ordered](items []T) (T, bool)`. Handle the empty slice.
2. Define `type Kelvin float64`. Write `Abs` with a `~float64` constraint and call it with `Kelvin`. Then remove `~` and record the compiler error.
3. Write `Equal[T comparable](a, b []T) bool`. Compare length first, then each index.

#### Advanced practical tasks

1. Write a constraint `Number` that covers signed integers, unsigned integers, and float types. Write `Sum[T Number](items []T) T`. Test three concrete types.
2. Explain why a slice type is comparable only in some cases. Write two type arguments: one that compiles with `Contains` and one that does not. Use the compiler messages as evidence.

---

## `any` vs custom constraints

`any` is an alias for `interface{}`. As a constraint, `any` allows every type. The function may copy values of `T`, put them in slices, and pass them to other functions that accept `T`. The function must not compare values of `T`. The function must not add values of `T`. The function must not call a method on `T` unless that method is in the constraint.

```go
func Clone[T any](items []T) []T {
    out := make([]T, len(items))
    copy(out, items)
    return out
}
```

`Clone` only copies elements. `any` is the correct constraint.

A custom constraint documents the operations that you need. Prefer the smallest constraint that makes the body legal.

```go
type Closer interface {
    Close() error
}

func CloseAll[T Closer](items []T) error {
    for _, item := range items {
        if err := item.Close(); err != nil {
            return err
        }
    }
    return nil
}
```

Here a plain interface parameter `[]Closer` also works. The generic form keeps the concrete type of each element. You do not box each value as an interface if the slice already holds a concrete type.

Use `any` when:

- the function only stores, copies, or moves values
- the caller supplies the behavior as a function argument (`func(T)`)

Use a custom constraint when:

- the function compares, orders, or hashes values
- the function calls methods on `T`
- the function uses operators such as `+` or `<`

Do not use `any` and then convert `T` to `any` (the interface) for type switches as a normal design. That mix loses the benefit of the type parameter. If you need a type switch on an unknown value, an ordinary `any` parameter is clearer.

The empty interface as a value and `any` as a constraint are different roles. A value of type `any` carries a dynamic type at run time. A type parameter constrained by `any` is a single concrete type at compile time after instantiation.

### Questions

#### Theoretical questions

1. What is the difference between `any` as a constraint and `any` as a variable type?
2. Which operations does a type parameter with constraint `any` allow?
3. When is `any` the correct constraint?
4. Why is a smaller custom constraint better than `any` when the body uses methods?
5. How does `CloseAll[T Closer]` differ from `CloseAll(items []Closer)` for the caller?
6. Why is a type switch on a type parameter a sign of a weak generic design?

#### Easy practical tasks

1. Write `Identity[T any](v T) T`. Call it with `int` and with `string`.
2. Write `Apply[T any](v T, f func(T) T) T`. Use it to double an `int`.
3. Write a custom constraint with one method `Name() string`. Write `Names[T YourConstraint](items []T) []string`.
4. Change `Contains` so that the constraint is `any`. Record the compiler error.

#### Medium practical tasks

1. Write two versions of a function that closes many values: one generic, one with `[]io.Closer`. Compare the call sites in comments.
2. Write `Filter[T any](items []T, keep func(T) bool) []T`. Test with even integers.
3. Write a constraint that requires `fmt.Stringer`. Write `Join[T fmt.Stringer](items []T, sep string) string`.

#### Advanced practical tasks

1. Take a function that uses `any` values and type assertions. Rewrite it with type parameters where that rewrite is honest. State one case that must stay as `any`.
2. Measure or reason about allocation: a `[]int` passed to `[]any` versus a generic `[]T`. Write six sentences about boxing and type safety.

---

## Type inference

The compiler may infer type arguments from the types of the ordinary arguments. You often omit the type argument list on a function call.

```go
n := First([]int{1, 2, 3})
// same as First[int]([]int{1, 2, 3})
```

Inference uses the argument types and the result assignment when those types are known. Inference does not guess a type from nothing.

```go
func Empty[T any]() []T {
    return nil
}

// _ = Empty()    // compile error: cannot infer T
_ = Empty[int]()  // explicit type argument
```

A function that returns only `T` and takes no value of type `T` usually needs an explicit type argument.

When more than one type parameter exists, inference uses each argument:

```go
func Transform[T, U any](in []T, f func(T) U) []U {
    out := make([]U, len(in))
    for i, v := range in {
        out[i] = f(v)
    }
    return out
}

strs := Transform([]int{1, 2}, strconv.Itoa)
```

The compiler infers `T` as `int` from the slice. The compiler infers `U` as `string` from `strconv.Itoa`.

Inference for generic types is weaker than inference for generic functions. You must write the type arguments when you name a generic type:

```go
s := Stack[int]{}
// Stack{} is not enough
```

A composite literal of a generic type needs the instantiated type. A helper function can hide that:

```go
func NewStack[T any]() *Stack[T] {
    return &Stack[T]{}
}

s := NewStack[int]() // still explicit, because no argument has type T
```

Unexported type parameters and exported type parameters follow the same inference rules. Inference is a compile-time process. The binary does not contain a separate inference step at run time.

If inference fails, the compiler reports that it cannot infer a type argument. Supply the type argument list. Do not change the algorithm only to force inference.

Go 1.21 and later improved inference for some function types and constraint cases. Write explicit type arguments when the call is hard to read. Clarity is more important than a shorter call.

### Questions

#### Theoretical questions

1. What information does the compiler use to infer a type argument?
2. Why does `Empty[T any]() []T` require an explicit type argument at the call site?
3. Does type inference run at run time? Explain in one sentence.
4. Why must you write `Stack[int]{}` instead of `Stack{}`?
5. What do you do when the compiler says it cannot infer a type argument?
6. How can a function with two type parameters infer both parameters from one call?

#### Easy practical tasks

1. Call `First` with a `[]string` and without a type argument list. Show that inference works.
2. Write `Zero[T any]() T`. Call it as `Zero[int]()` and as `Zero[string]()`.
3. Write a call that fails inference on purpose. Record the compiler message. Fix the call with an explicit type argument.
4. Use `slices.Contains` on a `[]int` without writing `[int]`. Confirm that the program compiles.

#### Medium practical tasks

1. Write `Transform` as in this section. Call it with a method value and with a function literal. Note whether inference succeeds in both cases.
2. Write `NewPair[A, B any](a A, b B) Pair[A, B]`. Show that the helper infers both type arguments from the values.
3. Instantiate a generic type as a map key type or a slice element type. Write the full instantiated name.

#### Advanced practical tasks

1. Build a case with three type parameters where inference succeeds only when you pass a function argument. Document which argument supplies each type parameter.
2. Read the language specification section on type inference. Write five short rules in your own words. Do not copy the specification text.

---

## Generic data structures vs interfaces

A generic data structure stores elements of type `T`. The compiler checks every insert and every read. A `Stack[int]` cannot hold a `string`. There is no type assertion at the use site.

An interface-based structure stores an interface type. A `Stack` of `any` holds every type. The caller must assert the type on the way out. A mistake appears at run time.

```go
// Generic: type-safe at compile time.
type List[T any] struct {
    items []T
}

// Interface: flexible, runtime checks.
type AnyList struct {
    items []any
}
```

Use a generic structure when:

- the collection is homogeneous (one element type)
- you want compile-time checks
- you care about avoiding interface boxing for small types

Use an interface when:

- one collection must hold many different types at once
- the element type is already an interface with methods (`io.Reader`, `error`)
- the set of types is open and you dispatch with a type switch

A method set on an interface describes behavior. A type parameter describes a placeholder type. Do not replace every interface with a type parameter. An `io.Reader` parameter is the correct API for "read bytes". A `Reader[T]` generic type is the wrong model for that job.

Generic structures and interfaces can work together. The structure is generic. The operations on elements use an interface constraint or a function argument:

```go
type Heap[T any] struct {
    items []T
    less  func(a, b T) bool
}
```

This heap does not require `T` to implement an interface. The caller supplies `less`. That design keeps `T` unconstrained and still defines order.

The standard library prefers functions in `slices` and `maps` over many generic container types. Go does not ship a generic list or tree in the standard library. Write a small structure in your module when you need one. Do not import a large container library for a short stack.

Boxing means the compiler stores a value in an interface. The interface value holds a type and a pointer (or a small direct value). A generic instantiation keeps the concrete layout of `T`. For `int` elements, a `[]int` is cheaper than a `[]any`.

### Questions

#### Theoretical questions

1. What compile-time guarantee does `List[int]` give that `[]any` does not give?
2. When is an interface-based collection the better model?
3. What is interface boxing?
4. Why is `io.Reader` still the right parameter type for a copy function?
5. How can a generic heap order elements without a constraint on `T`?
6. Does the standard library provide a generic tree type? What do you use instead?

#### Easy practical tasks

1. Implement `List[T any]` with `Append` and `At`. Use it with `int`.
2. Implement the same list with `[]any`. Show a type assertion in the getter.
3. Put a `string` into the `[]any` list by mistake. Show that the program compiles. Show that a wrong assertion panics or fails.
4. Write four sentences that compare the two lists.

#### Medium practical tasks

1. Write `Heap[T any]` with a `less` function. Insert five integers. Pop them in order.
2. Write a function that copies from `io.Reader` to `io.Writer`. Do not add type parameters. Explain why.
3. Store `fmt.Stringer` values in a slice. Print each value. State why an interface is enough.

#### Advanced practical tasks

1. Implement a generic LRU cache `Cache[K comparable, V any]` with a maximum size. Use a map and a list of keys. Write tests for eviction.
2. Compare assembly or a benchmark of `List[int].At` versus an `[]any` getter with a type assertion. Write which version allocates and why.

---

## When not to use generics

Generics add names, constraints, and instantiation errors. Use them when they remove duplication without hiding the program.

Do not use generics when only one concrete type exists. Write `func Sum(items []int) int`. A later second type is a reason to revisit the design. A possible future type is not a reason.

Do not use generics when an interface already names the behavior. `io.Reader`, `error`, and `fmt.Stringer` are stable APIs. A generic wrapper around them does not help the caller.

Do not use generics to build a framework that your program does not need. A generic dependency-injection container is usually harder to read than a few constructors.

Do not use generics to avoid small, clear functions. Two short functions with clear names are better than one generic function with a weak constraint and a callback that no one understands.

Do not mix generics with reflection as a default style. Reflection (`reflect`) solves open, dynamic cases. Generics solve closed, compile-time cases. If you already use `reflect` to walk unknown structs, a type parameter does not remove that work.

Do not export a generic API that leaks your implementation types. Keep the public surface small. Prefer a concrete type at the package boundary when the package always uses one type.

Good uses of generics include:

- containers (`Set[T]`, `Stack[T]`)
- slice and map helpers that the standard library does not cover
- numeric helpers with `cmp.Ordered` or a number constraint
- small adapters that preserve the concrete type

Check this list before you add a type parameter:

1. Do at least two real types need the same code today?
2. Does a standard function in `slices`, `maps`, or `cmp` already exist?
3. Does an interface already describe the behavior?
4. Can a function argument (`func(T)`) replace a complex constraint?
5. Will a beginner read the signature in one pass?

If the answer is no, skip, interface, callback, or rewrite, do not add generics.

### Questions

#### Theoretical questions

1. Why is one concrete type a reason to avoid generics?
2. Name two standard interfaces that you must not replace with type parameters.
3. What is the problem with a generic API that exists only for a possible future type?
4. When do you choose a function argument instead of a constraint?
5. How do generics and reflection differ in the problems that they solve?
6. What five checks do you run before you add a type parameter?

#### Easy practical tasks

1. Rewrite a generic `Sum[T any]` idea as `SumInts([]int) int`. State why the concrete function is enough.
2. Find one function in `slices` that replaces a helper that you might write. Write the function name and its constraint.
3. List three APIs in the standard library that use interfaces, not type parameters. Give one reason for each.
4. Take a generic function with one call site and one type. Remove the type parameter. Keep the tests green.

#### Medium practical tasks

1. Review a small package (yours or a public one). Mark each generic symbol as "needed" or "not needed". Give one sentence per symbol.
2. Replace a custom constraint with a `func(T) bool` argument. Compare the two signatures in six sentences.
3. Write a bad generic wrapper around `io.Copy`. Then delete it and write why the wrapper failed the checks in this section.

#### Advanced practical tasks

1. Design two APIs for a cache: one generic and one interface-based. Write the exported signatures only. Choose one for a service that stores `User` only. Choose one for a library that many modules import. Justify each choice.
2. Find a public Go module that overuses generics. Quote three signatures (short). Rewrite one signature without generics. Do not paste large files.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a generic function declaration to a compiled call with `int`. Name type parameters, constraints, inference, and instantiation.
2. When do you choose `comparable`, when do you choose `cmp.Ordered`, and when do you choose `any`?
3. Why must methods not declare new type parameters, and what do you write instead?
4. How do generic collections and interface values differ for a slice of integers?
5. A teammate wants to make every package API generic. Which facts from this topic do you use to refuse that plan?

#### Easy practical tasks

1. Create a module `example.com/genericx`. Add `First`, `Contains`, and `Stack[T]`. Add tests for `int` and `string`. Run `go test`.
2. Write a one-page cheat sheet: type parameter syntax, `comparable`, `~`, `cmp.Ordered`, inference failure, and one "do not use" rule.
3. Use `slices.Clone` on a `[]int`. Print the clone. Confirm that a change to the clone does not change the original length.
4. Draw a table with columns "Tool" and "Use when". Rows: type parameter, ordinary interface, `any` value, function argument.

#### Medium practical tasks

1. Implement `Set[T comparable]` and a non-generic `IntSet`. Write the same three tests against both. Record which tests were easier to copy.
2. Write `Merge[T any](a, b []T) []T` and `Unique[T comparable](items []T) []T`. Explain why the constraints differ.
3. Break a program with a missing `~` on a defined integer type. Fix the constraint. Write the error and the fix.

#### Advanced practical tasks

1. Implement a generic directed graph `Graph[T comparable]` with `AddEdge` and `Neighbors`. Write a depth-first search that takes `func(T)`. Add tests for a cycle and for a missing node.
2. Read the `slices` package documentation. List five functions and their constraints. Write one function that the package does not provide and that your project needs. Implement it with tests.
