# 6. Methods and Interfaces

## Description

A method is a function with a receiver. The receiver binds the method to a type. You call the method on a value or on a pointer of that type.

An interface is a set of method signatures. A type satisfies an interface when the type has those methods. You do not write an implements keyword. Satisfaction is implicit.

This topic shows value receivers, pointer receivers, and implicit satisfaction. It shows `any`, type assertions, and type switches. It shows common interfaces in the standard library.

Use methods to give a type behavior. Use interfaces to describe behavior that many types can provide. Accept interfaces in parameters. Return concrete structs from constructors when you can.

---

## Value receivers vs pointer receivers

A method declaration has a receiver in parentheses before the name.

```go
type Point struct {
	X int
	Y int
}

func (p Point) Sum() int {
	return p.X + p.Y
}

func (p *Point) Scale(k int) {
	p.X *= k
	p.Y *= k
}
```

The receiver `(p Point)` is a value receiver. The method gets a copy of the struct. A change to `p` does not change the caller value.

The receiver `(p *Point)` is a pointer receiver. The method gets a pointer. A change to `p.X` changes the caller storage.

You can call a value-receiver method on a value or on a pointer. You can call a pointer-receiver method on a pointer. You can also call it on an addressable value. The compiler takes the address. A value is addressable when it is a variable. A temporary value is not addressable. `Point{1, 2}.Scale(2)` does not compile.

The method set of `T` is the set of methods with receiver `T`. The method set of `*T` is the set of methods with receiver `T` or `*T`. Interfaces use method sets.

Use a pointer receiver when the method must change the receiver. Use a pointer receiver when the struct is large. Use a pointer receiver when some methods of the type already use a pointer receiver. Keep one style for the type.

Use a value receiver when the type is a small immutable value. A named map type can use a value receiver. A write to the map through a value receiver is visible to the caller. The header is a copy. The table is not.

A pointer-receiver method can receive `nil`. Check for `nil` when that state is valid. A constructor that returns `*T` pairs well with pointer methods.

### Questions

#### Theoretical questions

1. Where do you write the receiver in a method declaration?
2. What does a value receiver copy?
3. What does a pointer receiver let the method change?
4. When can you call a pointer-receiver method on a value?
5. What is the method set of `*T` compared with the method set of `T`?

#### Easy practical tasks

1. Add `Sum` with a value receiver on `Point`. Call it from `main`.
2. Add `Scale` with a pointer receiver. Call it on a variable. Print the fields.
3. Try `Point{1, 2}.Scale(2)`. Record the compiler error.
4. Call `Sum` on a `*Point` value. Print the result.

#### Medium practical tasks

1. Write `Move` with a value receiver that changes `X` and `Y`. Show that the caller `Point` does not change. Then rewrite `Move` with a pointer receiver.
2. Store `Point` in a map. Try to call a pointer-receiver method on `m[key]`. Record what happens. Then call it on a local copy.
3. Call a pointer method on a `nil` pointer of your type. Handle `nil` inside the method.

#### Advanced practical tasks

1. List every legal call form for `Sum` and `Scale` on `p`, `&p`, and a literal. Mark each form valid or invalid.
2. Find three types in the standard library: one that uses value receivers, one that uses pointer receivers, and one that mixes them. Write why each choice fits.

---

## Implicit interface satisfaction

An interface type lists methods. A type satisfies the interface when the method set of that type contains every listed method. You do not declare the link.

```go
type Speaker interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return d.Name
}

func say(s Speaker) {
	fmt.Println(s.Speak())
}
```

`Dog` satisfies `Speaker`. `say(Dog{Name: "Rex"})` is valid. No extra keyword is required. The method names and signatures must match. Extra methods on `Dog` are fine.

If `Speak` has a pointer receiver, `Dog` does not satisfy `Speaker`. `*Dog` does satisfy `Speaker`. Assign `&Dog{}` to `Speaker`.

An interface value has two parts: a dynamic type and a dynamic value. The zero interface value is `nil`. Both parts are unset.

A typed-nil pointer in an interface is not a `nil` interface. `var d *Dog; var s Speaker = d` makes `s == nil` false. This bug appears in error returns. Check the concrete pointer before you store it in an interface.

Keep interfaces small. One or two methods is a good size. Name interfaces with a verb or with `-er` when they have one method. Examples: `Speaker`, `Reader`, `error`.

You can embed interfaces in an interface. The new interface includes all methods. `io.ReadWriter` requires `Read` and `Write`. A type that already has both methods satisfies `ReadWriter`. Do not compose large interfaces.

### Questions

#### Theoretical questions

1. How does a type satisfy an interface?
2. Do you write an implements keyword?
3. When does `T` fail to satisfy an interface that `*T` satisfies?
4. What two parts does an interface value store?
5. Why is a typed-nil pointer inside an interface not equal to `nil`?

#### Easy practical tasks

1. Define `Speaker` and `Dog` as in the example. Call `say` with a `Dog`.
2. Add `Cat` with `Speak`. Call `say` with a `Cat`.
3. Remove `Speak` from `Cat`. Record the compiler error at the `say` call.
4. Assign `var d *Dog` to `Speaker`. Print `s == nil` and `d == nil`.

#### Medium practical tasks

1. Change `Speak` to a pointer receiver. Fix the call sites so the program compiles.
2. Write `say` in one package and the type in another package. Show that satisfaction still works.
3. Declare your own `ReadWriter` interface by embedding `io.Reader` and `io.Writer`. Assign `&bytes.Buffer{}` to that interface.

#### Advanced practical tasks

1. Reproduce a typed-nil `error` return. Show that `err != nil` is true. Then fix the return.
2. Design two small interfaces and one type that satisfies both. Pass the type to two functions.

---

## `any`, type assertions, and type switches

The empty interface has no methods. Every type satisfies it. The predeclared name `any` is an alias for `interface{}`. Use `any` in new code.

```go
var x any
x = 3
x = "hi"
x = Point{X: 1}
```

You cannot call methods on `x` without a type assertion. Prefer a real interface when you know the methods. `any` hides the type. Do not use `any` to avoid design.

A `nil` `any` value is a `nil` interface. After `x = (*int)(nil)`, `x == nil` is false.

A type assertion reads the dynamic type of an interface value.

```go
var s Speaker = Dog{Name: "Rex"}
d := s.(Dog)
```

`s.(Dog)` succeeds when the dynamic type is `Dog`. The assertion panics when the type does not match. The comma-ok form does not panic.

```go
d, ok := s.(Dog)
if !ok {
	fmt.Println("not a Dog")
}
```

A type switch selects a case by dynamic type. `v` has the type of the matching case.

```go
switch v := s.(type) {
case Dog:
	fmt.Println(v.Name)
case *Dog:
	fmt.Println(v.Name)
default:
	fmt.Println("other")
}
```

Do not use a type switch as your first design. Prefer methods on the interface. Use a type switch at a boundary. An assertion to a pointer type does not match a value type. `*Dog` is not `Dog`. Asserting on a `nil` interface panics in the single-value form. The comma-ok form returns `ok == false`.

### Questions

#### Theoretical questions

1. What is the relation between `any` and `interface{}`?
2. What does `s.(T)` return when `T` is the dynamic type?
3. What happens when `s.(T)` fails without comma-ok?
4. What is the type of `v` in each `case` of a type switch?
5. When do you use a type switch instead of a new interface method?

#### Easy practical tasks

1. Assign an `int`, then a `string`, to one `any` variable. Print the variable after each assign.
2. Store `Dog` in `Speaker`. Assert to `Dog` with comma-ok. Print `ok`.
3. Write a type switch on `any` that handles `int` and `string`.
4. Write `var x any` and `x = (*int)(nil)`. Print both `== nil` checks.

#### Medium practical tasks

1. Write `asDog(s Speaker) (Dog, error)` that uses comma-ok and returns an error on failure.
2. Add cases for `Dog` and `*Dog` in one type switch. Pass both kinds.
3. Build `map[string]any` with an `int` and a `string`. Read both keys with assertions.

#### Advanced practical tasks

1. Write a function that accepts `error` and uses a type switch for two custom error types. Return a string code for each case.
2. Decode a small JSON object into `map[string]any`. Assert one number and one string. Handle missing keys.

---

## Common interfaces: `error`, `fmt.Stringer`, `io.Reader`, `io.Writer`

The `error` interface has one method.

```go
type error interface {
	Error() string
}
```

`error` is predeclared. Functions return `error` as the last result. A `nil` error means success. Any type with `Error() string` can be an error. Use `errors.New` and `fmt.Errorf` in normal code. Topic 7 covers error handling in depth.

`fmt.Stringer` has one method.

```go
type Stringer interface {
	String() string
}
```

`fmt` calls `String` for verbs such as `%s` and `%v`. Implement `String` for a clear print form. Do not panic in `String`.

`io.Reader` has one method.

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

`Read` fills `p` with bytes. It returns the count and an error. `io.EOF` marks the end. Many types are readers: files, buffers, HTTP bodies.

`io.Writer` has one method.

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

`Write` writes bytes from `p`. `fmt.Fprint` writes to any `io.Writer`. `os.Stdout` is a writer. `bytes.Buffer` is a reader and a writer.

```go
var buf bytes.Buffer
fmt.Fprint(&buf, "hello")
```

These four interfaces are small. You can satisfy them with one method each. You can test with `bytes.Buffer` and `strings.NewReader`. Do not write `Read` or `Write` unless you provide the documented meaning. Wrong `io.EOF` rules break callers.

### Questions

#### Theoretical questions

1. What method does `error` require?
2. When does `fmt` call `String`?
3. What does `Read` return at the end of the input?
4. Why can `fmt.Fprint` write to a file and to a buffer?
5. Why are these interfaces small on purpose?

#### Easy practical tasks

1. Add `String() string` to `Point`. Print the point with `fmt.Println`.
2. Return `errors.New("bad")` from a function. Print `err.Error()`.
3. Read from `strings.NewReader("abc")` into a slice. Print `n` and the bytes.
4. Write a string to `os.Stdout` with `io.Writer` as the variable type.

#### Medium practical tasks

1. Copy from a `strings.NewReader` to a `bytes.Buffer` with `io.Copy`. Print the buffer.
2. Implement `Write` on a type that counts bytes. Use the type as `io.Writer` in `fmt.Fprint`.
3. Implement `Error() string` on a struct. Return that struct as `error`.

#### Advanced practical tasks

1. Write a `Reader` that returns one byte per `Read` from a string. Stop with `io.EOF`.
2. Read the docs for `io.Reader` and `io.Writer`. Write five rules that a correct implementation must follow.

---

## Accept interfaces, return structs

This guideline is a design rule. Function parameters accept interface types. Constructors and helpers return concrete struct types.

```go
type Store struct {
	path string
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func Save(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}
```

`NewStore` returns `*Store`. Callers see all methods of `Store`. You can add methods later without a new interface.

`Save` accepts `io.Writer`. Callers pass a file, a buffer, or a test writer. `Save` does not depend on `*os.File`.

If `NewStore` returned an interface, callers could not call extra methods. Tests would need a wider interface. New methods would break the interface or stay hidden.

Return an interface when you must hide the type. Examples: you return one of several types, or you must keep the concrete type unexported and small. Those cases are the exception.

Accept a concrete type when the function needs that exact type and no other type makes sense. Do not invent an interface for one implementation in the same package with no test need.

A parameter of type `any` is not a good substitute for a small interface. State the methods that you need.

This guideline works with implicit satisfaction. You do not change the concrete type to add `io.Writer` to a parameter list. You only need the methods.

### Questions

#### Theoretical questions

1. What does a function accept under this guideline?
2. What does a constructor return under this guideline?
3. Why is a returned interface often too narrow or too wide?
4. When is a returned interface acceptable?
5. Why is `any` a weak parameter type for this guideline?

#### Easy practical tasks

1. Write `NewPoint(x, y int) Point` that returns a struct, not an interface.
2. Write `printStringer(s fmt.Stringer)` that accepts an interface.
3. Rewrite a function that takes `*bytes.Buffer` so that it takes `io.Writer`.
4. Write three sentences that state the guideline in your own words.

#### Medium practical tasks

1. Create `*FileStore` with `Save` and `Load`. Return `*FileStore` from `New`. Accept `io.Reader` in `Load`.
2. Try to return an interface from `New` with only `Save`. Show that `main` cannot call `Load` without an assert.
3. Write a test-friendly function that accepts `io.Reader` and counts bytes. Use `strings.NewReader` in `main`.

#### Advanced practical tasks

1. Review one of your earlier functions. Change one parameter to an interface and one result to a struct. Record the call-site changes.
2. Find one standard library function that returns an interface and one that returns a struct. Explain both choices.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do method sets decide whether `T` or `*T` may go in an interface variable?
2. How do you choose among a value receiver, a pointer receiver, and no method at all?
3. When do you use `any`, a small interface, and a type assertion?
4. How do type assertions, type switches, and extra interface methods solve the same need in different ways?
5. How do `error`, `fmt.Stringer`, `io.Reader`, and the accept-interfaces guideline fit together?

#### Easy practical tasks

1. Write a struct with one pointer method and one value method. Store the value in an interface that lists only the value method. Print a successful call.
2. Implement `fmt.Stringer` and `error` on two different types. Print both.
3. Write `func Dump(w io.Writer, v fmt.Stringer) error` and call it with a `bytes.Buffer`.
4. Draw a table: Receiver kind, Changes caller, In method set of `T`, In method set of `*T`.

#### Medium practical tasks

1. Build a tiny pipeline: a type that satisfies `io.Reader`, a function that accepts `io.Reader`, and a constructor that returns the concrete type.
2. Use a type switch on `any` to route `int`, `string`, and `error` values to three print lines.
3. Refactor a function that returns an interface so that it returns a struct. Keep one parameter as an interface.

#### Advanced practical tasks

1. Write a package with an unexported struct, an exported constructor that returns the struct or a pointer, and a function that accepts `io.Writer`. Show a test-style caller that uses `bytes.Buffer`.
2. Design method receivers and interfaces for a `Cache` type with `Get` and `Set`. Write signatures only.
