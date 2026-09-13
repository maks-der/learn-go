# 7. Methods and Interfaces

## Description

A method is a function with a receiver. The receiver binds the method to a type. You call the method on a value or on a pointer of that type.

An interface is a set of method signatures. A type satisfies an interface when the type has those methods. You do not write an implements keyword. Satisfaction is implicit.

This topic shows value receivers, pointer receivers, and the method set of each form. It shows the empty interface `any`. It shows type assertions and type switches. It shows common interfaces in the standard library.

Use methods to give a type behavior. Use interfaces to describe behavior that many types can provide. Accept interfaces in parameters. Return concrete structs from constructors when you can.

---

## Method receivers: value vs pointer

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

You can call a value-receiver method on a value or on a pointer. The compiler copies the value or dereferences the pointer.

You can call a pointer-receiver method on a pointer. You can also call it on an addressable value. The compiler takes the address. A value is addressable when it is a variable. A temporary value is not addressable.

```go
p := Point{X: 1, Y: 2}
p.Scale(2)
(&p).Scale(2)
```

`Point{1, 2}.Scale(2)` does not compile. The literal is not addressable.

The method set of `T` is the set of methods with receiver `T`. The method set of `*T` is the set of methods with receiver `T` or `*T`. Interfaces use method sets. A later section covers satisfaction.

Keep one receiver name for the type. Short names are common. Example: `p` for `Point`.

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

1. Write `Move` with a value receiver that changes `X` and `Y`. Show that the caller `Point` does not change.
2. Rewrite `Move` with a pointer receiver. Show that the caller changes.
3. Store `Point` in a map. Try to call a pointer-receiver method on `m[key]`. Record what happens. Then call it on a local copy.

#### Advanced practical tasks

1. List every legal call form for `Sum` and `Scale` on `p`, `&p`, and a literal. Mark each form valid or invalid.
2. Write one type with both receiver kinds. Document which methods change the value and which methods only read.

---

## When to use pointer receivers

Use a pointer receiver when the method must change the receiver. A value receiver cannot persist a change.

Use a pointer receiver when the struct is large. A value receiver copies every field on each call. A pointer receiver copies one pointer.

Use a pointer receiver when some methods of the type already use a pointer receiver. Keep one style for the type. Mixed styles confuse the method set.

Use a pointer receiver when you declare the interface methods with pointer receivers. Only `*T` has those methods in its method set. `T` does not satisfy that interface.

Use a value receiver when the type is a small immutable value. Examples: `time.Time` methods often use a value receiver. A small struct of numbers can use a value receiver for read-only methods.

Use a value receiver when the type is a small struct that you treat as a value. A named slice type can use a value receiver. A named map type can use a value receiver. Do not replace the header in that method. A value receiver still shares a map header or a slice header. A write to the map through a value receiver is visible to the caller. The header is a copy. The table is not.

Do not pick a pointer receiver only because it looks common. Pick it for mutation, size, consistency, or interfaces.

A constructor that returns `*T` pairs well with pointer methods. A constructor that returns `T` pairs well with value methods. You can still return `T` and take addresses at the call site.

Nil receivers are possible. A pointer-receiver method can receive `nil`. Check for `nil` when that state is valid. Example: a method on a linked list node.

### Questions

#### Theoretical questions

1. Name four reasons to choose a pointer receiver.
2. When is a value receiver a good choice?
3. Why does mixed receiver style cause interface problems?
4. Can a value-receiver method change a map that the struct holds?
5. What happens when a pointer-receiver method gets a `nil` receiver?

#### Easy practical tasks

1. Write a large struct with many fields. Add one pointer-receiver method that sets one field.
2. Write a small `Pair` struct with a value-receiver `Sum` method.
3. Give one type two methods: one pointer and one value. Write why this mix is a problem.
4. Call a pointer method on a `nil` pointer of your type. Handle `nil` inside the method.

#### Medium practical tasks

1. Change a type from all value receivers to all pointer receivers. Update call sites. Note any interface breaks.
2. Time or reason about a method on a struct that holds `[1024]int`. Explain copy cost of a value receiver.
3. Write a constructor `NewT() *T` and two pointer methods. Use only `*T` in `main`.

#### Advanced practical tasks

1. Find three types in the standard library: one that uses value receivers, one that uses pointer receivers, and one that mixes them. Write why each choice fits.
2. Design a type that is safe with a `nil` pointer receiver. Document the `nil` meaning.

---

## Interfaces: implicit satisfaction

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

`Dog` satisfies `Speaker`. `say(Dog{Name: "Rex"})` is valid. No extra keyword is required.

The method names and signatures must match. `Speak() string` is not `Speak() int`. Extra methods on `Dog` are fine. Missing methods mean the type does not satisfy the interface.

If `Speak` has a pointer receiver, `Dog` does not satisfy `Speaker`. `*Dog` does satisfy `Speaker`. Assign `&Dog{}` to `Speaker`.

An interface value has two parts: a dynamic type and a dynamic value. The zero interface value is `nil`. Both parts are unset.

```go
var s Speaker
fmt.Println(s == nil)
s = Dog{Name: "Rex"}
```

A typed-nil pointer in an interface is not a `nil` interface. `var d *Dog; var s Speaker = d` makes `s == nil` false. This bug appears in error returns. Check the concrete pointer before you store it in an interface.

Keep interfaces small. One or two methods is a good size. The standard library uses small interfaces.

Name interfaces with a verb or with `-er` when they have one method. Examples: `Speaker`, `Reader`, `error`.

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
3. Give `Dog` an extra method `Run`. Confirm that `say` still accepts `Dog`.

#### Advanced practical tasks

1. Reproduce a typed-nil `error` return. Show that `err != nil` is true. Then fix the return.
2. Design two small interfaces and one type that satisfies both. Pass the type to two functions.

---

## Empty interface `any` / `interface{}`

The empty interface has no methods. Every type satisfies it. The predeclared name `any` is an alias for `interface{}`. Use `any` in new code. Go 1.18 added `any`. Both names are the same type.

```go
var x any
x = 3
x = "hi"
x = Point{X: 1}
```

`x` can hold any value. You cannot call methods on `x` without a type assertion. You cannot use `x` as an `int` without a conversion step.

`fmt.Println` accepts `...any`. Many APIs use `any` when they store unknown values. Example: `map[string]any` for a loose JSON object.

Prefer a real interface when you know the methods. `any` hides the type. Callers must assert. Tests become harder.

A `nil` `any` value is a `nil` interface. `var x any` has `x == nil` true. After `x = (*int)(nil)`, `x == nil` is false.

Do not use `any` to avoid design. Use a struct, a generic type parameter, or a small interface when you can.

`any` is not a special runtime type beyond `interface{}`. Tools and docs may still show `interface{}`.

### Questions

#### Theoretical questions

1. What methods does the empty interface list?
2. What is the relation between `any` and `interface{}`?
3. Can you call `Speak` on a variable of type `any` without an assertion?
4. When do you prefer a small interface over `any`?
5. When is `var x any` equal to `nil`?

#### Easy practical tasks

1. Assign an `int`, then a `string`, to one `any` variable. Print the variable after each assign.
2. Write a function `printAny(v any)` that calls `fmt.Println(v)`.
3. Store three values of different types in `[]any`. Print the length.
4. Write `var x any` and `x = (*int)(nil)`. Print both `== nil` checks.

#### Medium practical tasks

1. Build `map[string]any` with an `int` and a `string`. Read both keys.
2. Replace an `any` parameter with a small interface in a sample function. Write why the new form is safer.
3. Use `fmt.Sprintf("%T", v)` on an `any` value to print the dynamic type.

#### Advanced practical tasks

1. Decode a small JSON object into `map[string]any`. Assert one number and one string. Handle missing keys.
2. Compare `any`, a one-method interface, and a type parameter for one helper. Write six sentences.

---

## Type assertions and type switches

A type assertion reads the dynamic type of an interface value.

```go
var s Speaker = Dog{Name: "Rex"}
d := s.(Dog)
```

`s.(Dog)` succeeds when the dynamic type is `Dog`. The result is a `Dog` value. The assertion panics when the type does not match.

The comma-ok form does not panic.

```go
d, ok := s.(Dog)
if !ok {
	fmt.Println("not a Dog")
}
```

You can assert to an interface type. The dynamic value must satisfy that interface.

A type switch selects a case by dynamic type.

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

`v` has the type of the matching case. In `default`, `v` has the same type as `s`.

Do not use a type switch as your first design. Prefer methods on the interface. Use a type switch at a boundary. Examples: decode a value, handle a few known types, or adapt old APIs.

An assertion to a pointer type does not match a value type. `*Dog` is not `Dog`. Include the case that you store.

Asserting on a `nil` interface panics in the single-value form. The comma-ok form returns `ok == false`.

### Questions

#### Theoretical questions

1. What does `s.(T)` return when `T` is the dynamic type?
2. What happens when `s.(T)` fails without comma-ok?
3. What does the comma-ok form return on failure?
4. What is the type of `v` in each `case` of a type switch?
5. When do you use a type switch instead of a new interface method?

#### Easy practical tasks

1. Store `Dog` in `Speaker`. Assert to `Dog` with comma-ok. Print `ok`.
2. Assert to `Cat` when the value is a `Dog`. Use comma-ok. Print `ok`.
3. Write a type switch on `any` that handles `int` and `string`.
4. Run the single-value assert on the wrong type. Record the panic.

#### Medium practical tasks

1. Write `asDog(s Speaker) (Dog, error)` that uses comma-ok and returns an error on failure.
2. Add cases for `Dog` and `*Dog` in one type switch. Pass both kinds.
3. Assert an `any` value to `Speaker`. Show success when the dynamic type has `Speak`.

#### Advanced practical tasks

1. Write a function that accepts `error` and uses a type switch for two custom error types. Return a string code for each case.
2. Show that `s.(T)` on a `nil` interface panics, and that the comma-ok form does not. Print both outcomes.

---

## Common interfaces: `error`, `fmt.Stringer`, `io.Reader`, `io.Writer`

The `error` interface has one method.

```go
type error interface {
	Error() string
}
```

`error` is predeclared. Functions return `error` as the last result. A `nil` error means success. Any type with `Error() string` can be an error. Use `errors.New` and `fmt.Errorf` in normal code. The next topic covers error handling in depth.

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

These four interfaces are small. You can satisfy them with one method each. You can test with `bytes.Buffer` and `strings.NewReader`.

Do not write `Read` or `Write` unless you provide the documented meaning. Wrong `io.EOF` rules break callers.

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
3. Read from `strings.NewReader("abc")` into a `[8]byte` or a slice. Print `n` and the bytes.
4. Write a string to `os.Stdout` with `io.Writer` as the variable type.

#### Medium practical tasks

1. Copy from a `strings.NewReader` to a `bytes.Buffer` with `io.Copy`. Print the buffer.
2. Implement `Write` on a type that counts bytes. Use the type as `io.Writer` in `fmt.Fprint`.
3. Implement `Error() string` on a struct. Return that struct as `error`.

#### Advanced practical tasks

1. Write a `Reader` that returns one byte per `Read` from a string. Stop with `io.EOF`.
2. Read the docs for `io.Reader` and `io.Writer`. Write five rules that a correct implementation must follow.

---

## Interface composition

You can embed interfaces in an interface. The new interface includes all methods.

```go
type ReadWriter interface {
	io.Reader
	io.Writer
}
```

`ReadWriter` requires `Read` and `Write`. A type that already has both methods satisfies `ReadWriter`. `*bytes.Buffer` is an example.

The standard library uses this form. `io.ReadCloser` embeds `Reader` and `Closer`. `io.ReadWriteCloser` embeds more.

```go
type Closer interface {
	Close() error
}

type ReadCloser interface {
	io.Reader
	Closer
}
```

Composition keeps interfaces small. You add names without a new method list. Callers that need only `Read` still accept a `ReadCloser`.

A type does not embed the interface unless you write that in the type. Satisfaction stays implicit. `*os.File` satisfies `ReadWriter` because it has the methods.

Do not compose large interfaces. A big interface is hard to satisfy and hard to mock. Compose two or three small interfaces.

You can embed an interface in a struct. The struct then promotes those methods. The embedded field must be non-nil when you call a promoted method. This pattern is a wrapper, not a subclass.

Name the composed interface after the combined job. Do not invent a new method that only forwards without a reason.

### Questions

#### Theoretical questions

1. What methods does `io.ReadWriter` require?
2. Does a type need to name `ReadWriter` to satisfy it?
3. Why does composition help keep interfaces small?
4. Can a function that takes `io.Reader` accept an `io.ReadCloser`?
5. What happens when a struct embeds an interface field that is `nil`?

#### Easy practical tasks

1. Declare your own `ReadWriter` interface by embedding `io.Reader` and `io.Writer`.
2. Assign `&bytes.Buffer{}` to that interface. Write and then read.
3. Write a function that takes `io.Reader`. Pass an `io.ReadCloser` (`io.NopCloser` or a file).
4. List three composed interfaces in package `io`.

#### Medium practical tasks

1. Define `StringCloser` that embeds `fmt.Stringer` and a `Close() error` method.
2. Write a struct that embeds `io.Reader`. Set the field to a `strings.NewReader`. Call `Read` on the struct.
3. Show that a type with only `Read` does not satisfy `ReadWriter`. Record the error.

#### Advanced practical tasks

1. Embed `io.Writer` in a struct and add a `Write` method on the struct. Document which `Write` runs.
2. Design three composed interfaces for a storage API: read, write, and close. Keep each base interface at one method.

---

## The "accept interfaces, return structs" guideline

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
3. Write a test-friendly function that accepts `io.Reader` and count bytes. Use `strings.NewReader` in `main`.

#### Advanced practical tasks

1. Review one of your earlier functions. Change one parameter to an interface and one result to a struct. Record the call-site changes.
2. Find one standard library function that returns an interface and one that returns a struct. Explain both choices.

---

## `comparable` constraint (high-level)

`comparable` is a predeclared interface. You use it as a constraint on a type parameter. A type argument must support `==` and `!=`.

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

`Contains` works for `int`, `string`, and other comparable types. `Contains` does not work for `[]int` as `T`. Slices are not comparable. Maps and functions are not comparable.

Comparable types include booleans, numbers, strings, pointers, channels, arrays of comparable elements, and structs of comparable fields. Interface values are comparable. A comparison of two interfaces panics when the dynamic types are not comparable.

Go 1.20 allows ordinary interface types to satisfy `comparable` as a constraint. You can write `map[K]V` with `K` constrained by `comparable`. That is the rule for map keys.

`comparable` is not a normal interface for variables. You do not write `var x comparable` in ordinary code. You use `comparable` in the type-parameter list.

Do not use `comparable` to force `==` on types that need a custom equal method. Write an `Equal` method or pass a function when equality is not the built-in operator.

This section is high-level. A later topic covers generics in depth. Remember that map keys and `==` need comparable types.

### Questions

#### Theoretical questions

1. What operators does `comparable` require?
2. Which composite types are not comparable?
3. Can you use `comparable` as the type of a normal variable?
4. Why do map keys need a comparable type?
5. What can happen when you compare two interface values?

#### Easy practical tasks

1. Call a generic `Contains` with `[]string`. Search for one word.
2. Try to instantiate `Contains` with `T` equal to `[]int`. Record the compiler error.
3. Show that `Point` with two `int` fields works with `==`.
4. Write `func Same[T comparable](a, b T) bool` and call it with two integers.

#### Medium practical tasks

1. Write `Index[T comparable](s []T, v T) int` that returns `-1` when `v` is missing.
2. Add a `[]byte` field to a struct. Show that the struct is no longer comparable.
3. Declare `map[K]int` in a generic function with `K comparable`. Insert one key.

#### Advanced practical tasks

1. Write why `comparable` is the wrong constraint for a slice of structs that need field-by-field equal with a custom rule.
2. Read the current release notes or spec text on `comparable` and interfaces. Write five sentences on what a beginner must remember.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do method sets decide whether `T` or `*T` may go in an interface variable?
2. How do you choose among a value receiver, a pointer receiver, and no method at all?
3. When do you use `any`, a small interface, and `comparable`?
4. How do type assertions, type switches, and extra interface methods solve the same need in different ways?
5. How do `error`, `fmt.Stringer`, `io.Reader`, and interface composition fit the accept-interfaces guideline?

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
2. Design method receivers and interfaces for a `Cache` type with `Get` and `Set`. Include a note on `comparable` keys. Write signatures only.
