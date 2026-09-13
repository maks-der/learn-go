# 12. Generics

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
```

Methods on a generic type may use the type parameters of the type. Methods must not declare new type parameters. If you need a second type parameter on an operation, write a function, not a method.

```go
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

You cannot use a type parameter in a method receiver list as a new name. The receiver uses the type parameters of the type. Two type parameters that use the same constraint can still be different types. `A` and `B` in `Pair` are not the same type.

Do not put a type parameter on a function when the function body never uses `T` as a type. That parameter is noise.

### Questions

#### Theoretical questions

1. Where do you write type parameters on a function?
2. What is the constraint `any` on `T`?
3. Can a method declare its own type parameters?
4. How do you instantiate `Stack` for `int`?
5. What do you return when a generic function has no value of type `T`?

#### Easy practical tasks

1. Write `First[T any]`. Call it on `[]string` and on `[]int`. Print both results.
2. Write `Pair[A, B any]`. Build `Pair[string, int]{"n", 1}`. Print the fields.
3. Add `Pop` to `Stack[T]`. Pop from an empty stack and print the zero value and `false`.
4. Instantiate `First[int]` with an explicit type argument. Then call `First` without one.

#### Medium practical tasks

1. Write `MapSlice[T, U any](s []T, f func(T) U) []U`. Map `[]int` to `[]string`.
2. Write `Convert` as in the example. Convert a `Stack[int]` to `Stack[string]`.
3. Try to add a method `func (s *Stack[T]) Convert[U any](...)`. Record the compiler error. Keep the function form.

#### Advanced practical tasks

1. Design a generic `Result[T any]` with fields `Value T` and `Err error`. Add `Ok` and `Err` helpers. Do not hide panics.
2. Write a generic linked list with `Insert` and `Range`. Use a function with its own type parameter to map nodes to another list.

---

## Constraints and `comparable`

A constraint is an interface. The type argument must satisfy that interface. `any` allows every type. A constraint with methods requires those methods.

```go
type Stringer interface {
	String() string
}

func Join[T Stringer](items []T, sep string) string {
	// call items[i].String()
}
```

An interface can list types. That form is a type set. The type argument must be one of those types or a type whose underlying type is in the set.

```go
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
```

The tilde `~` includes named types whose underlying type is the listed type. `type Count int` satisfies `~int`. `Count` does not satisfy `int` without the tilde.

You can combine methods and type sets in one interface. Package `golang.org/x/exp/constraints` is not in the standard library. The standard library has `cmp.Ordered` from Go 1.21 for ordered types.

```go
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
```

`comparable` is a predeclared constraint. A type argument must support `==` and `!=`. You need `comparable` for map keys and for `==` in the function body.

```go
func Contains[T comparable](s []T, v T) bool {
	for i := range s {
		if s[i] == v {
			return true
		}
	}
	return false
}
```

Slices, maps, and functions are not comparable. `Contains` does not work for `T` equal to `[]int`. Interface values are comparable. A comparison of two interfaces panics when the dynamic types are not comparable.

`comparable` is not a normal interface for variables. You do not write `var x comparable` in ordinary code.

Do not use `comparable` to force `==` on types that need a custom equal method. Write an `Equal` method or pass a function.

### Questions

#### Theoretical questions

1. What does a constraint limit?
2. What does `~int` add that `int` does not add in a type set?
3. What operators does `comparable` require?
4. Why do map keys need a comparable type?
5. Which composite types are not comparable?

#### Easy practical tasks

1. Call a generic `Contains` with `[]string`. Search for one word.
2. Try to instantiate `Contains` with `T` equal to `[]int`. Record the compiler error.
3. Write `Min` with `cmp.Ordered`. Call it with two `int` values and two `string` values.
4. Write `type Count int` and a constraint with `~int`. Pass `Count` to a `Max` function.

#### Medium practical tasks

1. Write `Index[T comparable](s []T, v T) int` that returns `-1` when `v` is missing.
2. Add a `[]byte` field to a struct. Show that the struct is no longer comparable.
3. Declare `map[K]int` in a generic function with `K comparable`. Insert one key.

#### Advanced practical tasks

1. Write a constraint that requires `~string` and a `Width() int` method. Implement the method on a named string type. Call a generic `Pad` function.
2. Write why `comparable` is the wrong constraint for a slice of structs that need field-by-field equal with a custom rule.

---

## Type inference

The compiler infers type arguments from the values that you pass. You can omit the type argument list when inference succeeds.

```go
n, ok := First([]int{1, 2})
```

The compiler sees `[]int` and binds `T` to `int`. You can still write `First[int]([]int{1, 2})`.

Inference uses the arguments, the assignment target, and constraints. Inference can fail when the compiler cannot find one type.

```go
func Empty[T any]() []T {
	return nil
}

s := Empty[int]() // you must write [int]
```

There is no argument that names `T`. The compiler cannot infer `T`.

When two parameters share a type parameter, the arguments must match.

```go
func Both[T any](a, b T) {}
Both(1, 2)    // T is int
// Both(1, "x") // error
```

A function argument can help inference. `MapSlice(nums, strconv.Itoa)` can infer `T` and `U` from `nums` and from `Itoa`.

Unexported type parameters follow the same inference rules. Inference does not care about export.

If inference surprises you, write the type arguments. Explicit arguments are legal and clear.

From Go 1.21, inference is stronger for some function types and unused type parameters. Prefer a readable call. Do not fight the compiler with extra wrappers only to hide `[T]`.

### Questions

#### Theoretical questions

1. When can you omit `[int]` on a generic call?
2. Why does `Empty[T any]() []T` need an explicit type argument?
3. What happens when two arguments disagree on a shared type parameter?
4. Does inference use constraints?
5. When do you write explicit type arguments even if inference works?

#### Easy practical tasks

1. Call `First` without type arguments on a `[]string`. Print the value.
2. Call `Empty` without type arguments. Record the error. Then write `Empty[string]()`.
3. Call `Both(1, 2)`. Then call `Both(1, int64(2))` and record the result.
4. Write `MapSlice` and call it with `[]int{1}` and a function that returns a `string`. Omit type arguments if you can.

#### Medium practical tasks

1. Write `NewPair[A, B any](a A, b B) Pair[A, B]`. Call it without type arguments.
2. Assign `var f func([]int) (int, bool) = First`. Show that the assignment helps or still needs `[int]`.
3. Build a case where inference picks `int` and you wanted `int64`. Fix the call with an explicit argument or a conversion.

#### Advanced practical tasks

1. Read the language specification section on type inference. Write eight short sentences in your own words. Do not copy the text.
2. Design three helpers: one that always infers, one that never infers, and one that infers only from a function argument. Implement them.

---

## When not to use generics

Generics add type parameters to the API. Readers must learn those names. The compiler messages grow. A simple function is often clearer.

Do not use generics when only one type exists. Write `func SumInts(s []int) int`.

Do not use generics to replace an interface that already describes behavior. If every type has `Read`, accept `io.Reader`. An interface is the right tool for mixed values in one slice. `[]io.Reader` holds different concrete types. `[]T` holds one type `T`.

Do not use generics to avoid `any` when you still need a type switch inside. A type switch on `T` with constraint `any` is often worse than a small interface.

Do not write a huge constraint that lists every type in the program. Split the function or use methods.

Do not build a generic "container library" for a small app. A slice and a map are enough. Add `Stack[T]` when you use it in several packages and the operations are stable.

Use generics when:

- the same algorithm works for many types
- you need `==` or `<` on the type parameter
- you need a data structure that stores `T` without `any` and assertions

The standard library uses generics in `slices`, `maps`, and `cmp`. Read those packages. Copy their style: small functions, clear constraints, no extra abstraction.

Measure before you claim that generics are faster. Instantiation can increase binary size. A method on an interface can allocate. Profile if size or speed matters.

### Questions

#### Theoretical questions

1. When is a non-generic function the better API?
2. Why does `io.Reader` beat a generic `Read[T]` for mixed sources?
3. What problem appears when a generic function still uses a type switch on `T`?
4. When is a generic container worth the extra API?
5. Name two standard library packages that use generics.

#### Easy practical tasks

1. Write `SumInts` without generics. Then write `Sum[T ~int | ~int64]`. Call both. Write which one you keep for a single-type package.
2. Store a `bytes.Buffer` and a `os.File` in `[]io.Writer`. Explain why `[]T` cannot hold both.
3. Open `go doc slices`. Write three function names and their constraints.
4. List four rules from this section in a two-column table: "Do" and "Do not".

#### Medium practical tasks

1. Refactor a helper that takes `any` and asserts `int` or `string`. Split it into two functions or one generic function with a constraint. State the better form.
2. Compare `slices.Contains` with your own `Contains`. Delete yours if the standard function is enough.
3. Count instantiations: call a generic function with five types. Write whether the extra API is worth it in that file.

#### Advanced practical tasks

1. Find one generic type in your earlier exercises that you can delete. Replace it with a slice or a small interface. Record the line count before and after.
2. Read the Go blog post on generics. Write a one-page decision list for your team. Do not copy the post.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do type parameters, constraints, and inference work together in one call?
2. When do you choose `any`, `comparable`, `cmp.Ordered`, and a method constraint?
3. Why can a method not add a new type parameter, and what do you write instead?
4. How do generics and interfaces split the job of reuse?
5. What signs tell you to remove a type parameter from an API?

#### Easy practical tasks

1. Write `Last[T any]` and `Contains[T comparable]`. Call both on `[]string`.
2. Write `Stack[T any]` with `Push` and `Len`. Push two `int` values.
3. Call `Min` from `cmp` or your own `Min` with two `float64` values.
4. Draw a table: Feature, Example, When you need it. Add rows for type parameter, constraint, inference, and non-generic function.

#### Medium practical tasks

1. Implement `Filter[T any](s []T, keep func(T) bool) []T` and `Unique[T comparable](s []T) []T`. Test both.
2. Write a generic `Set[T comparable]` with `Add` and `Has`. Do not export the map field.
3. Show one call that needs an explicit type argument and one call that infers. Write why.

#### Advanced practical tasks

1. Build a tiny `slices`-style package with `Map`, `Reduce`, and `Index`. Use the standard `slices` package where it already exists. Document each choice.
2. Design an API that mixes an interface parameter and a type parameter. Write why both appear. Implement one function.
