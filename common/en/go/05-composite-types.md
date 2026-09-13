# 5. Composite Types

## Description

A composite type holds other values. This topic covers arrays, slices, maps, and structs. You use these types in almost every Go program.

An array has a fixed length. A slice is a view of an array. A map stores keys and values. A struct groups named fields.

Prefer slices over arrays for most sequences. Prefer maps for lookup by key. Prefer structs for records with a known set of fields.

Length is part of an array type. A slice header stores a pointer, a length, and a capacity. A map value refers to a hash table. A struct value holds its fields.

---

## Arrays

An array type is `[N]T`. `N` is the length. `T` is the element type. The length is part of the type. `[3]int` and `[4]int` are different types.

```go
var a [3]int
b := [3]int{1, 2, 3}
c := [...]int{4, 5, 6}
```

`var a [3]int` sets every element to the zero value. The form `[...]int{4, 5, 6}` lets the compiler count the elements. The type of `c` is `[3]int`.

Assignment copies an array. The copy has its own elements. A change to `b` does not change `a` after `a = b`. Pass of an array to a function also copies the array. A large array as a parameter is expensive. Use a slice when the length can change or when you want to avoid a copy.

```go
a := [2]int{1, 2}
b := a
b[0] = 9
fmt.Println(a[0], b[0]) // 1 9
```

Arrays are comparable when the element type is comparable. You can use `==` on two `[3]int` values. You can use an array of comparable elements as a map key.

Index an array with `a[i]`. The index must be in `0` to `N-1`. An out-of-range index causes a panic.

A multidimensional array is an array of arrays. `[2][3]int` has two rows. Each row has three `int` values. Assignment copies every inner array. `[2][3]int` is not the same as `[][]int`.

Use arrays when the size is part of the meaning. Examples: a hash digest, a size-2 point, a small lookup table. Use slices for lists.

### Questions

#### Theoretical questions

1. Why are `[3]int` and `[4]int` different types?
2. What does assignment of an array copy?
3. What does `[...]int{1, 2}` mean?
4. When is an array type comparable?
5. How is `[2][3]int` different from `[][]int`?

#### Easy practical tasks

1. Declare `var days [7]string`. Set two elements. Print `len(days)`.
2. Create `[3]int` with a literal. Copy it to a second variable. Change one element. Print both arrays.
3. Use `[...]string{"a", "b"}`. Print the type with `fmt.Printf("%T\n", v)`.
4. Declare `[2][2]int`. Set `m[1][0]`. Print all four values.

#### Medium practical tasks

1. Write a function that swaps two elements of `[4]int` and returns the new array. Keep the input unchanged.
2. Compare two `[3]int` values with `==`. Then change one element and compare again.
3. Write `transpose` for `[2][3]int` that returns `[3][2]int`.

#### Advanced practical tasks

1. Use `[16]byte` as a map key for a short digest. Insert two keys. Look up one key.
2. Implement matrix add for `[N][N]int` with `N := 3` as a named constant. Keep arrays, not slices.

---

## Slices: length, capacity, `append`, and the backing array

A slice is a descriptor for a segment of an array. The slice header has three parts: a pointer, a length, and a capacity.

The length is the number of elements that you can index. `len(s)` returns the length. Valid indexes are `0` to `len(s)-1`. The capacity counts elements in the backing array from the slice start to the array end. `cap(s)` returns the capacity. Capacity is always greater than or equal to length.

```go
s := make([]int, 3, 5)
fmt.Println(len(s), cap(s)) // 3 5
```

You can index `s[0]`, `s[1]`, and `s[2]`. You cannot index `s[3]`. That index panics. Extra capacity is unused until `append` or until a slice expression grows the length.

`make` creates a slice. Use `make([]T, len)` or `make([]T, len, cap)`. `append` adds elements at the end. `append` returns a slice. You must use the returned slice.

```go
s := []int{1, 2}
s = append(s, 3)
s = append(s, 4, 5)
s = append(s, []int{6, 7}...)
```

When length is less than capacity, `append` writes into the backing array. When length equals capacity, `append` allocates a new array. The runtime copies the old elements.

`copy` copies elements from a source slice to a destination slice. `copy` returns the number of elements that it copied.

A slicing expression makes a new slice header. The new slice can share the same backing array. The expression is `x[low:high]`. The new length is `high - low`. The three-index form `x[low:high:max]` sets the capacity to `max - low`. Use it when the next `append` must not write into later elements of the original array.

A change to an element of a subslice changes the backing array. Other slices that share that array see the change. When you need an independent slice, `copy` the elements to a new slice.

A slice can be `nil`. A slice can be empty and not `nil`. Both have length `0`. `append` works on a `nil` slice. `encoding/json` encodes a `nil` slice as `null`. It encodes an empty slice as `[]`. Prefer `len(s) == 0` when you only care about elements.

`append` may write into unused capacity that another slice still uses. Fix sharing with a three-index slice or with `copy`. Do not ignore the result of `append`. A small slice of a large backing array keeps the large array alive. Copy the small part out when you keep it for a long time.

### Questions

#### Theoretical questions

1. What three parts does a slice header store?
2. Why must you assign the result of `append`?
3. When does `append` allocate a new backing array?
4. What does the third index in `x[low:high:max]` set?
5. How does `encoding/json` encode a `nil` slice and an empty slice?

#### Easy practical tasks

1. Create `make([]int, 2, 6)`. Print `len` and `cap`.
2. Start from `nil` and `append` three integers. Print the slice.
3. From `[5]int{10, 20, 30, 40, 50}` make a slice of the middle three elements. Print it.
4. Print `s == nil`, `len(s)`, and `cap(s)` for `var s []int` and for `s := []int{}`.

#### Medium practical tasks

1. Start with `make([]int, 0, 4)`. Append one element at a time. Print `len` and `cap` after each append until `len` is `5`.
2. Use a three-index slice so that `append` does not change the original array. Show the array after `append`.
3. Reproduce `base[:2]` plus `append`. Then repeat with `base[:2:2]`. Print `base` in both cases.

#### Advanced practical tasks

1. Implement `clone(s []int) []int` that does not share a backing array. Prove it with a write to the clone.
2. Write a function that takes a slice and returns the first three elements as an isolated slice. Tests must not see later `append` on the result change the input.

---

## Maps: lookup, comma-ok, delete, and iteration

A map type is `map[K]V`. `K` is the key type. `V` is the element type. The key type must be comparable.

`make` creates an empty map that you can write. A composite literal also creates a map.

```go
m := make(map[string]int)
m["a"] = 1
m := map[string]int{"a": 1, "b": 2}
```

Lookup uses `m[key]`. The result is the value. When the key is missing, the result is the zero value of `V`. The comma-ok form reports presence.

```go
v, ok := m["a"]
if !ok {
	fmt.Println("missing")
}
```

`ok` is `true` when the key is in the map. `ok` is `false` when the key is missing. `make(map[K]V, hint)` takes a size hint. The hint is not a length limit. `len(m)` is the number of keys. There is no `cap` for a map.

`delete` removes a key. `delete(m, key)` is a no-op when the key is missing. `delete` on a `nil` map is also a no-op. `clear(m)` removes all keys. After `clear`, `len(m)` is `0`. The map is still not `nil`.

Range over a map visits each key and value. The iteration order is random. Do not write tests that require a fixed print order. If you need a stable order, collect the keys, sort the keys, then look up each key.

```go
for k, v := range m {
	fmt.Println(k, v)
}
```

A map can be `nil`. A read from a `nil` map returns the zero value. A write to a `nil` map panics. Create the map with `make` or with a literal before the first write.

A map value refers to a hash table. Assignment copies the map header. The copy refers to the same table. A function that receives a map can insert, update, and delete keys. The caller sees those changes.

Maps are not safe for concurrent write. Two goroutines must not write the same map at the same time. Concurrent read with no writes is safe. Protect a shared map with `sync.Mutex`. Topic 10 covers goroutines and `sync`.

### Questions

#### Theoretical questions

1. What constraint does the key type have?
2. What does `m[key]` return when the key is missing?
3. What does the comma-ok form tell you?
4. Is map iteration order stable?
5. Which operation on a `nil` map panics?

#### Easy practical tasks

1. `make` a `map[string]int`. Set two keys. Print one value.
2. Look up a missing key without comma-ok. Print the result. Then use comma-ok and print `ok`.
3. Insert two keys. `delete` one key. Print `len` and a comma-ok lookup of the deleted key.
4. Print `m == nil` and `len(m)` for `var m map[int]int` and for `make(map[int]int)`.

#### Medium practical tasks

1. Write `get(m map[string]int, k string) (int, error)` that returns an error when the key is missing.
2. Count word frequencies in a small `[]string` with a map.
3. Print a `map[string]int` in key-sorted order.

#### Advanced practical tasks

1. Write a map with a struct key that has two comparable fields. Insert and look up one key.
2. Write a small type that holds a map and a `sync.Mutex`. Add `Get` and `Set`. Focus on the lock around map access.

---

## Structs, embedding, and struct tags

A struct is a sequence of named fields. Each field has a name and a type.

```go
type Point struct {
	X int
	Y int
}
```

The zero value of a struct has the zero value in each field. Access a field with a dot. `p.X = 3`. Two named struct types are different even when the fields match. A struct is comparable when every field is comparable.

Prefer field names in a composite literal. Names stay correct when someone adds a field. You can take the address of a composite literal: `p := &Point{X: 3, Y: 4}`. Do not mix positional fields and named fields in one struct literal.

A field name that starts with an upper-case letter is exported. Other packages can read and write that field. A field name that starts with a lower-case letter is unexported. Only the same package can use that name. `encoding/json` encodes exported fields by default. A JSON tag does not export a field.

An embedded field lists a type without a new name. The field name becomes the type name. The fields of the inner type are promoted.

```go
type Circle struct {
	Point
	Radius int
}

type Disk struct {
	Center Point
	Radius int
}
```

You can write `c.X` for `Circle`. You must write `d.Center.X` for `Disk`. Embedding is not inheritance. `Circle` is not a `Point`. A function that takes `Point` does not accept `Circle`.

Use embedding when the inner type is a part of the outer type and promotion helps. Use a named field when the inner value is a role such as `Center`. You can embed a pointer type. A `nil` embedded pointer panics when you use a promoted field.

A struct tag is a string literal after a field type. Packages such as `encoding/json` use tags.

```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age,omitempty"`
	Pass string `json:"-"`
}
```

`json:"name"` maps the field to the JSON name `name`. `omitempty` omits the field when the value is the zero value. `json:"-"` skips the field in JSON. The tag does not change the Go field name. Use backticks for the tag string.

### Questions

#### Theoretical questions

1. What does the zero value of a struct contain?
2. What is the difference between `Point{1, 2}` and `Point{X: 1, Y: 2}`?
3. What makes a field exported?
4. Is `Circle` a `Point` when `Circle` embeds `Point`?
5. What does `json:"-"` mean?

#### Easy practical tasks

1. Define `Book` with `Title string` and `Pages int`. Set both fields. Print them.
2. Build `User` with a named-field literal. Print it with `%+v`.
3. Build a `Circle` literal with an inner `Point` literal. Read `c.X` and `c.Radius`.
4. Define `User` with `json:"name"` and `json:"age"`. Marshal one value. Print the JSON.

#### Medium practical tasks

1. Create two packages. Export a type from package `a` with an unexported field. Try to set that field from `package main`. Record the error. Then add a constructor.
2. Embed two structs that both have a field `ID`. Show the compiler error for `v.ID`. Use the long selectors.
3. Unmarshal `{"name":"Ada","age":36}` into `User`. Print the struct. Add `omitempty` on `Age` and marshal a user with `Age` 0.

#### Advanced practical tasks

1. Redesign `Circle` once with embedding and once with `Center Point`. Write six sentences on call sites and clarity.
2. Write `User` with `Pass` hidden from JSON and a constructor that sets `Pass`. Marshal and confirm the secret is absent.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do array value semantics, slice headers, and map headers differ when you assign the value?
2. When does a write through one variable change another variable for arrays, slices, maps, and structs?
3. How do `nil` slices and `nil` maps differ in write behavior?
4. Why is a three-index slice part of a safe `append` plan?
5. How do exported fields, embedding, and tags change what another package and JSON can see?

#### Easy practical tasks

1. Build a `[]Point` from a struct type that you define. Append one point. Print `len` and `cap`.
2. Store those points in a `map[string]Point`. Look up one key with comma-ok.
3. Take a subslice of three integers. `copy` it to a new slice. Change the copy. Print both.
4. Draw a table with rows Array, Slice, Map, Struct. Columns: length in the type, nil write, comparable.

#### Medium practical tasks

1. Write a small in-memory catalog: a struct type, a slice of that type, and a map from id to index. Add one item and look it up.
2. Show one program that hits a slice capacity surprise and the fixed version with isolation. Print both results.
3. Define a struct with an embedded struct, one unexported field, and a json tag. Marshal a value. Explain the JSON.

#### Advanced practical tasks

1. Design types for a grid: choose array or slice of slices. Justify the choice with length, sharing, and growth.
2. Write a package-level constructor that returns a struct with an initialized map field and an empty non-nil slice field. Show safe writes in `main`.
