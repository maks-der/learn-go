# 6. Composite Types

## Description

A composite type holds other values. This topic covers arrays, slices, maps, and structs. You use these types in almost every Go program.

An array has a fixed length. A slice is a view of an array. A map stores keys and values. A struct groups named fields.

Prefer slices over arrays for most sequences. Prefer maps for lookup by key. Prefer structs for records with a known set of fields.

Length is part of an array type. A slice header stores a pointer, a length, and a capacity. A map value refers to a hash table. A struct value holds its fields.

Each section below is one idea. Read the array sections first. Then read slices. Then read maps. Then read structs.

---

## Arrays: fixed length and value semantics

An array type is `[N]T`. `N` is the length. `T` is the element type. The length is part of the type. `[3]int` and `[4]int` are different types.

```go
var a [3]int
b := [3]int{1, 2, 3}
c := [...]int{4, 5, 6}
```

`var a [3]int` sets every element to the zero value. The composite literal `[3]int{1, 2, 3}` sets each index. The form `[...]int{4, 5, 6}` lets the compiler count the elements. The type of `c` is `[3]int`.

Assignment copies an array. The copy has its own elements. A change to `b` does not change `a` after `a = b`.

```go
a := [2]int{1, 2}
b := a
b[0] = 9
fmt.Println(a[0], b[0])
```

This prints `1` and `9`. The arrays do not share storage.

Pass of an array to a function also copies the array. A large array as a parameter is expensive. Use a slice when the length can change or when you want to avoid a copy.

Arrays are comparable when the element type is comparable. You can use `==` on two `[3]int` values. You can use an array of comparable elements as a map key.

Index an array with `a[i]`. The index must be in `0` to `N-1`. An out-of-range index causes a panic.

Use arrays when the size is part of the meaning. Examples: a hash digest, a size-2 point, a small lookup table. Use slices for lists.

### Questions

#### Theoretical questions

1. Why are `[3]int` and `[4]int` different types?
2. What does assignment of an array copy?
3. What does `[...]int{1, 2}` mean?
4. When is an array type comparable?
5. Why is a large array a bad function parameter?

#### Easy practical tasks

1. Declare `var days [7]string`. Set two elements. Print `len(days)`.
2. Create `[3]int` with a literal. Copy it to a second variable. Change one element. Print both arrays.
3. Write a function that takes `[2]int` and returns the sum. Show that the caller array does not change.
4. Use `[...]string{"a", "b"}`. Print the type with `fmt.Printf("%T\n", v)`.

#### Medium practical tasks

1. Write a function that swaps two elements of `[4]int` and returns the new array. Keep the input unchanged.
2. Compare two `[3]int` values with `==`. Then change one element and compare again.
3. Pass an array and a pointer to an array to two functions. Show which call can change the caller storage.

#### Advanced practical tasks

1. Use `[16]byte` as a map key for a short digest. Insert two keys. Look up one key.
2. Measure or reason about copy cost of `[1000000]int` versus a slice header. Write six sentences.

---

## Multidimensional arrays

A multidimensional array is an array of arrays. `[2][3]int` has two rows. Each row has three `int` values.

```go
var m [2][3]int
m[0][1] = 7
n := [2][3]int{
	{1, 2, 3},
	{4, 5, 6},
}
```

The inner type is `[3]int`. The outer type is `[2][3]int`. Assignment copies every inner array. The copy is deep for the values. There is no shared backing array.

Index from the outer array first. `n[1][2]` is `6` in the example. An inner index out of range panics.

`[2][3]int` is not the same as `[][]int`. `[][]int` is a slice of slices. Each inner slice can have a different length. A slice of slices can share backing arrays. An array of arrays cannot change length.

You can range over the outer array. Each element is an inner array. Range over that inner array to visit each value.

```go
for i := range n {
	for j := range n[i] {
		fmt.Println(i, j, n[i][j])
	}
}
```

Use a multidimensional array when every row has the same fixed size. Use a slice of slices when rows grow or when rows have different lengths.

A function parameter of type `[2][3]int` copies all six integers. Pass a pointer when you must change the caller array.

### Questions

#### Theoretical questions

1. What is the element type of `[2][3]int`?
2. Does assignment of `[2][3]int` share storage?
3. How is `[2][3]int` different from `[][]int`?
4. Which index is the outer index in `m[i][j]`?
5. When do you choose an array of arrays over a slice of slices?

#### Easy practical tasks

1. Declare `[2][2]int`. Set `m[1][0]`. Print all four values.
2. Write a composite literal for `[2][3]string` with six words.
3. Copy a `[2][2]int` value. Change one cell in the copy. Print both.
4. Range over a `[3][2]int` and print each pair `(i, j, value)`.

#### Medium practical tasks

1. Write `transpose` for `[2][3]int` that returns `[3][2]int`.
2. Write `rowSum` that returns `[2]int` sums for each row of `[2][3]int`.
3. Compare a `[2][2]int` with a `[][]int` that you build with the same numbers. Change one inner slice and explain the difference.

#### Advanced practical tasks

1. Implement matrix add for `[N][N]int` with `N := 3` as a named constant. Keep arrays, not slices.
2. Build a 3-D array `[2][2][2]int`. Write a function that returns the sum of all elements.

---

## Slice length vs capacity

A slice is a descriptor for a segment of an array. The slice header has three parts: a pointer, a length, and a capacity.

The length is the number of elements that you can index. `len(s)` returns the length. Valid indexes are `0` to `len(s)-1`.

The capacity counts elements in the backing array. The count starts at the slice start. The count ends at the array end. `cap(s)` returns the capacity. Capacity is always greater than or equal to length.

```go
s := make([]int, 3, 5)
fmt.Println(len(s), cap(s))
```

This prints `3 5`. You can index `s[0]`, `s[1]`, and `s[2]`. You cannot index `s[3]`. That index panics. `append` can use the extra capacity.

An index must use the length, not the capacity. Extra capacity is unused until `append` or until a slice expression grows the length.

`len` and `cap` are built-in functions. They work on slices, arrays, and some other types. For a slice, both are fast.

A zero slice value has length `0` and capacity `0`. The next sections cover `nil` slices and `append`.

Do not store capacity in your own field unless you wrap a slice on purpose. Use `len` and `cap` on the slice value.

### Questions

#### Theoretical questions

1. What three parts does a slice header store?
2. What does `len(s)` count?
3. What does `cap(s)` count?
4. Can you index an element in the unused capacity?
5. What is the relation between `len(s)` and `cap(s)`?

#### Easy practical tasks

1. Create `make([]int, 2, 6)`. Print `len` and `cap`.
2. Print `len` and `cap` of `[]int{1, 2, 3}`.
3. Index the last valid element of a slice. Then write the index that would panic.
4. Print `len` and `cap` of `var s []int`.

#### Medium practical tasks

1. Start with `make([]int, 0, 4)`. Append one element at a time. Print `len` and `cap` after each append until `len` is `5`.
2. Take a slice of an array. Print `len` and `cap` of the slice. Relate `cap` to the array length.
3. Write a function that returns `len(s)` and `cap(s)` for a `[]int`. Call it before and after one `append`.

#### Advanced practical tasks

1. Build a slice with `len` 2 and `cap` 8. Use a slice expression to make a new slice with a different `len` and the same `cap` start rule. Print both headers.
2. Explain a case where `len` stays small and `cap` is large. State one memory risk.

---

## `make`, `append`, and `copy`

`make` creates a slice, a map, or a channel. For a slice, use `make([]T, len)` or `make([]T, len, cap)`.

```go
a := make([]int, 3)
b := make([]int, 0, 8)
```

`a` has length `3` and capacity `3`. Each element is the zero value. `b` has length `0` and capacity `8`. `b` can grow with `append` without a new array at first.

`append` adds elements at the end. `append` returns a slice. You must use the returned slice.

```go
s := []int{1, 2}
s = append(s, 3)
s = append(s, 4, 5)
s = append(s, []int{6, 7}...)
```

When length is less than capacity, `append` writes into the backing array. When length equals capacity, `append` allocates a new array. The runtime copies the old elements. The new capacity is larger.

`copy` copies elements from a source slice to a destination slice. `copy` returns the number of elements that it copied. That number is the smaller of `len(dst)` and `len(src)`.

```go
dst := make([]int, 2)
src := []int{9, 8, 7}
n := copy(dst, src)
```

`n` is `2`. `dst` is `{9, 8}`. `copy` handles overlap in a correct way.

`make` with only a length sets capacity equal to length. Use an extra capacity when you know the final size. That choice reduces allocations.

Do not ignore the result of `append`. The variable on the left must receive the new header.

### Questions

#### Theoretical questions

1. What is the difference between `make([]T, n)` and `make([]T, n, c)`?
2. Why must you assign the result of `append`?
3. When does `append` allocate a new backing array?
4. How many elements does `copy` copy?
5. How do you append all elements of one slice to another slice?

#### Easy practical tasks

1. `make` a `[]string` of length `3`. Set index `1`. Print the slice.
2. Start from `nil` and `append` three integers. Print the slice.
3. `copy` two elements from `[]int{1, 2, 3}` into a new slice of length `2`.
4. Append a slice to another slice with `...`. Print the result.

#### Medium practical tasks

1. Preallocate with `make([]int, 0, 10)`. Append ten values in a loop. Confirm that `cap` stays `10`.
2. Use `copy` to clone a slice so that later `append` on the clone does not change the original.
3. Show `copy` with overlap: copy `s[0:2]` to `s[1:3]` on a four-element slice. Print before and after.

#### Advanced practical tasks

1. Write `appendAll` that appends many `[]int` values into one slice. Use a capacity hint from the total length.
2. Compare `make([]T, n)` plus index fills with `make([]T, 0, n)` plus `append`. Write when each form is correct.

---

## Slicing expressions and backing arrays

A slicing expression makes a new slice header. The expression uses an existing array or slice. The new slice can share the same backing array.

```go
a := [5]int{0, 1, 2, 3, 4}
s := a[1:4]
```

`s` has length `3`. The elements are `1`, `2`, and `3`. The capacity is `4`. Four elements remain from index `1` to the end of the array.

The expression is `x[low:high]`. The new length is `high - low`. The new capacity is `cap(x) - low` when `x` is a slice. Indexes follow the half-open rule. `low` is in. `high` is out.

You can omit `low` or `high`. `s[:]` keeps the same length. `s[1:]` drops the first element. `s[:2]` keeps the first two elements.

The three-index form is `x[low:high:max]`. That form sets the capacity to `max - low`. Use it when the next `append` must not write into later elements of the original array.

```go
s := a[1:3:3]
```

`s` has length `2` and capacity `2`. `append` on `s` allocates a new array. The original array does not receive the new element.

A change to an element of `s` changes the backing array. Other slices that share that array see the change.

```go
t := a[0:2]
s[0] = 9
fmt.Println(a, t)
```

Plan sharing. When you need an independent slice, `copy` the elements to a new slice.

A slice of a slice does not copy elements. It only builds a new header.

### Questions

#### Theoretical questions

1. What does `x[low:high]` share with `x`?
2. How do you compute the new length?
3. What does the third index in `x[low:high:max]` set?
4. Why is the range half-open?
5. When do you `copy` instead of slice?

#### Easy practical tasks

1. From `[5]int{10, 20, 30, 40, 50}` make a slice of the middle three elements. Print it.
2. Print `s[:]`, `s[1:]`, and `s[:2]` for `s := []int{1, 2, 3, 4}`.
3. Change one element of a subslice. Print the original slice.
4. Print `len` and `cap` of `a[2:]` when `a` is `[6]int`.

#### Medium practical tasks

1. Use a three-index slice so that `append` does not change the original array. Show the array after `append`.
2. Make two overlapping slices of one array. Change one index in the overlap. Print both slices.
3. Write a function that returns `s[1:len(s)-1]` and documents the panic when `len(s) < 2`.

#### Advanced practical tasks

1. Implement `clone(s []int) []int` that does not share a backing array. Prove it with a write to the clone.
2. Show a case where `s[low:high]` has a large capacity and a later `append` overwrites data that another slice still uses.

---

## Nil vs empty slice

A slice can be `nil`. A slice can be empty and not `nil`. Both have length `0`. They are not the same in every API.

```go
var n []int
e := []int{}
m := make([]int, 0)
```

`n` is `nil`. `len(n)` is `0`. `cap(n)` is `0`. `e` and `m` are empty. They are not `nil`. `len` is `0` for both.

`append` works on a `nil` slice. The first `append` allocates a backing array. Range over a `nil` slice does zero iterations. `copy` to or from a `nil` slice copies zero elements.

`n == nil` is `true` for `var n []int`. `e == nil` is `false`. Prefer `len(s) == 0` when you only care about elements.

`encoding/json` encodes a `nil` slice as `null`. It encodes an empty slice as `[]`. That difference matters in APIs.

A function that returns a slice with no elements can return `nil`. That return is fine. Document the JSON choice when you encode the result.

Do not write special cases for `nil` before `append` or `range`. Those operations already accept `nil`.

### Questions

#### Theoretical questions

1. How do you create a `nil` slice?
2. How do you create an empty slice that is not `nil`?
3. Does `append` work on a `nil` slice?
4. How does `encoding/json` encode a `nil` slice and an empty slice?
5. When do you test `s == nil` instead of `len(s) == 0`?

#### Easy practical tasks

1. Print `s == nil`, `len(s)`, and `cap(s)` for `var s []int` and for `s := []int{}`.
2. Range over a `nil` slice. Confirm that the loop body does not run.
3. `append` one value to a `nil` `[]string`. Print the result and `s == nil` after the append.
4. Encode a `nil` slice and an empty slice with `json.Marshal`. Print both outputs.

#### Medium practical tasks

1. Write `nonEmpty(s []int) []int` that returns `nil` when `len(s) == 0` and returns `s` otherwise.
2. Show `copy` from a `nil` source into a destination of length `3`. Print the destination.
3. Write a function that always returns an empty non-nil slice when there are no items. Use `make` or `[]T{}`.

#### Advanced practical tasks

1. Build a JSON object with a field of type `[]int`. Show `null` versus `[]` for callers. Write which form you pick and why.
2. Compare `var s []T`, `s := []T{}`, and `s := make([]T, 0)` in a table: nil, len, cap, JSON.

---

## Common slice pitfalls

Two common pitfalls are append sharing and capacity surprises.

`append` may write into unused capacity. That capacity can belong to an array that another slice still uses. The other slice then sees new values. The other slice can also lose old values when `append` overwrites them.

```go
base := []int{1, 2, 3, 4}
a := base[:2]
b := append(a, 9)
fmt.Println(base, a, b)
```

`a` has capacity `4`. `append` writes `9` into `base[2]`. `base` becomes `{1, 2, 9, 4}`. A reader can miss this result.

When `append` needs more capacity, it allocates a new array. Later writes do not change the old array. The same code can share or not share. The result depends on `cap`.

```go
a := []int{1, 2}
b := append(a, 3)
```

If `cap(a)` is `2`, `b` has a new array. If `cap(a)` is larger, `b` may share with `a`.

Fix sharing with a three-index slice: `a := base[:2:2]`. Then `append` cannot write into `base`. Fix sharing with `copy` to a new slice before `append`.

Another pitfall is to ignore the `append` result. The old header can stay at the old length. The new header has the new length. Always assign `s = append(s, x)`.

Another pitfall is to keep a small slice of a large backing array. The garbage collector cannot free the large array while the small slice lives. Copy the small part out when you keep it for a long time.

Do not assume that `append` always allocates. Check `cap` or cut the capacity when you need isolation.

### Questions

#### Theoretical questions

1. When can `append` change an element that another slice still shows?
2. Why does the same `append` code sometimes allocate and sometimes not allocate?
3. How does a three-index slice stop an overwrite?
4. What goes wrong when you ignore the result of `append`?
5. Why can a small slice keep a large array alive?

#### Easy practical tasks

1. Reproduce the `base[:2]` plus `append` example. Print `base` after the append.
2. Repeat the example with `base[:2:2]`. Print `base` again.
3. Write `s = append(s, x)` and a second version that does not assign. Compare prints.
4. Print `cap` of a subslice and explain whether the next `append` can share.

#### Medium practical tasks

1. Write `grow(s []int, x int) []int` that never writes into the caller backing array.
2. Show two calls of the same helper: one with extra capacity and one without. Record both results.
3. Keep `s[:1]` from a large `make([]byte, 1_000_000)` slice. Then copy the one byte out. Explain the memory difference.

#### Advanced practical tasks

1. Write a function that takes a slice and returns the first three elements as an isolated slice. Tests must not see later `append` on the result change the input.
2. Find one extra pitfall in your own code or in a review of a classmate file. Write the cause and the fix. Do not reuse the examples above as the only case.

---

## Map creation, lookup, and comma-ok

A map type is `map[K]V`. `K` is the key type. `V` is the element type. The key type must be comparable.

`make` creates an empty map that you can write.

```go
m := make(map[string]int)
m["a"] = 1
```

A composite literal also creates a map.

```go
m := map[string]int{"a": 1, "b": 2}
```

Lookup uses `m[key]`. The result is the value. When the key is missing, the result is the zero value of `V`. You cannot see a missing key with a single lookup when zero is a valid stored value.

The comma-ok form reports presence.

```go
v, ok := m["a"]
if !ok {
	fmt.Println("missing")
}
```

`ok` is `true` when the key is in the map. `ok` is `false` when the key is missing.

`make(map[K]V, hint)` takes a size hint. The hint is not a length limit. The map can grow past the hint. The hint can reduce early growth.

`len(m)` is the number of keys. There is no `cap` for a map.

Do not use a slice, a map, or a function as a key type. Those types are not comparable.

### Questions

#### Theoretical questions

1. What constraint does the key type have?
2. What does `m[key]` return when the key is missing?
3. What does the comma-ok form tell you?
4. What does the hint in `make(map[K]V, hint)` mean?
5. Why is a missing key a problem when the zero value is valid data?

#### Easy practical tasks

1. `make` a `map[string]int`. Set two keys. Print one value.
2. Look up a missing key without comma-ok. Print the result.
3. Look up the same missing key with comma-ok. Print `ok`.
4. Create a map with a literal that has three pairs. Print `len(m)`.

#### Medium practical tasks

1. Write `get(m map[string]int, k string) (int, error)` that returns an error when the key is missing.
2. Count word frequencies in a small `[]string` with a map.
3. Use `make(map[string]int, 100)` and insert ten keys. Print `len`.

#### Advanced practical tasks

1. Write a map with a struct key that has two comparable fields. Insert and look up one key.
2. Compare lookup with comma-ok versus a stored pointer value. Write when each form shows absence.

---

## Delete and iteration order

`delete` removes a key. `delete(m, key)` is a no-op when the key is missing. `delete` on a `nil` map is also a no-op.

```go
m := map[string]int{"a": 1, "b": 2}
delete(m, "a")
delete(m, "missing")
```

After the first `delete`, `"a"` is gone. The second `delete` does nothing.

Range over a map visits each key and value.

```go
for k, v := range m {
	fmt.Println(k, v)
}
```

The iteration order is random. The runtime changes the start position. Do not write tests that require a fixed print order. Do not treat the first pair from `range` as a minimum or as insertion order.

If you need a stable order, collect the keys. Sort the keys. Then look up each key.

```go
keys := make([]string, 0, len(m))
for k := range m {
	keys = append(keys, k)
}
sort.Strings(keys)
```

You can range over keys only: `for k := range m`. You can ignore the key: `for _, v := range m`.

`clear(m)` removes all keys. `clear` exists in Go 1.21 and later. After `clear`, `len(m)` is `0`. The map is still not `nil`.

Do not delete keys from a map while you assume a fixed visit order. You may delete the current key during `range`. Do not depend on visits of keys that you add during the same `range`.

### Questions

#### Theoretical questions

1. What does `delete` do when the key is missing?
2. Is map iteration order stable?
3. How do you print map keys in a fixed order?
4. What does `clear(m)` do to `len(m)` and to `m == nil`?
5. Why must a test not compare a `range` print to one fixed string?

#### Easy practical tasks

1. Insert two keys. `delete` one key. Print `len` and a comma-ok lookup of the deleted key.
2. Range over a map of three pairs. Run the program two times. Record whether the print order changed.
3. Range over keys only and collect them in a slice. Print the slice length.
4. Call `clear` on a map. Print `len` and `m == nil`.

#### Medium practical tasks

1. Print a `map[string]int` in key-sorted order.
2. Delete a key inside `range` when the value is `0`. Print the map after the loop.
3. Write `keys(m map[string]int) []string` that returns a sorted key list.

#### Advanced practical tasks

1. Write a test that checks map contents without a fixed iteration order. Use lookups or a sort.
2. Add keys during a `range` and document which new keys you saw. Run the program more than one time.

---

## Nil map vs empty map

A map can be `nil`. A map can be empty and not `nil`. Reads differ from writes.

```go
var n map[string]int
e := map[string]int{}
m := make(map[string]int)
```

`n` is `nil`. `e` and `m` are empty maps. `len` is `0` for all three. `n == nil` is `true`. The others are not `nil`.

A read from a `nil` map returns the zero value. Comma-ok returns `ok == false`. Range over a `nil` map does zero iterations. `delete` on a `nil` map does nothing. `len` of a `nil` map is `0`.

A write to a `nil` map panics. The message is `assignment to entry in nil map`.

```go
var n map[string]int
n["a"] = 1 // panic
```

Create the map with `make` or with a literal before the first write.

`encoding/json` encodes a `nil` map as `null`. It encodes an empty map as `{}`.

Return a `nil` map when you have no pairs and you do not write more pairs. Return `make(map[K]V)` when the caller will insert keys.

Do not leave a struct field as a `nil` map if later code writes to it. Initialize that field in a constructor.

### Questions

#### Theoretical questions

1. What operations are safe on a `nil` map?
2. Which operation on a `nil` map panics?
3. How do you create an empty map that accepts writes?
4. How does `encoding/json` encode a `nil` map and an empty map?
5. Why does a struct field of map type need `make` before writes?

#### Easy practical tasks

1. Print `m == nil` and `len(m)` for `var m map[int]int` and for `make(map[int]int)`.
2. Look up a key in a `nil` map with comma-ok. Print `v` and `ok`.
3. Write to a `make` map. Then write the `nil` write in a small program and recover or let it panic once on purpose.
4. Marshal a `nil` map and an empty map to JSON. Print both.

#### Medium practical tasks

1. Write `ensure(m map[string]int) map[string]int` that returns `m` when `m != nil` and returns a new map otherwise.
2. Range over a `nil` map and over an empty map. Confirm both loops skip the body.
3. Add a struct with a map field. Show a panic when you write without `make`. Then fix the constructor.

#### Advanced practical tasks

1. Design a function that returns a map. State when the function returns `nil` and when it returns an empty map. Give one caller for each case.
2. Compare `var m map[K]V`, `m := map[K]V{}`, and `m := make(map[K]V)` in a table: nil, write safe, JSON.

---

## Maps are reference-like

A map value refers to a hash table. Assignment copies the map header. The copy refers to the same table. A write through one variable is visible through the other variable.

```go
a := map[string]int{"x": 1}
b := a
b["x"] = 2
fmt.Println(a["x"])
```

This prints `2`. `a` and `b` share the table.

A function that receives a map can insert, update, and delete keys. The caller sees those changes. The function cannot replace the caller map with `nil` by assignment to the parameter. The parameter is a copy of the header.

Maps are not safe for concurrent write. Two goroutines must not write the same map at the same time. A write and a read at the same time are also not safe. The runtime can stop the program with `fatal error: concurrent map writes`. `recover` does not handle that fault.

Concurrent read with no writes is safe. Many goroutines can read a map that no goroutine changes.

Protect a shared map with `sync.Mutex`. Give the map to one goroutine. Or use `sync.Map` when that type matches the need. Later topics cover goroutines and `sync`.

A map is not a pointer in the source. You write `map[K]V`, not `*map[K]V`, for normal use. The reference-like behavior is in the header.

Do not pass a map to another goroutine and write from both sides without a lock.

### Questions

#### Theoretical questions

1. What does assignment of a map copy?
2. Can a function insert a key that the caller then sees?
3. Are concurrent writes to one map safe?
4. Are concurrent reads safe when there is no write?
5. Does `recover` stop `fatal error: concurrent map writes`?

#### Easy practical tasks

1. Assign a map to a second variable. Change one key through the second variable. Print the first map.
2. Write a function that sets `m[k] = v`. Show that `main` sees the new pair.
3. Write a function that does `m = nil` on the parameter. Show that `main` still has a non-nil map.
4. Print two sentences that state the concurrency rule for maps.

#### Medium practical tasks

1. Write `merge(dst, src map[string]int)` that writes all `src` keys into `dst`.
2. Clone a map by ranging and inserting into a new `make` map. Change the clone. Show that the original is stable.
3. Document a design that avoids shared writes: one function owns the map and returns results.

#### Advanced practical tasks

1. Write a small type that holds a map and a `sync.Mutex`. Add `Get` and `Set`. Do not start extra goroutines if you are not ready. Focus on the lock around map access.
2. Read the Go FAQ or blog text on map concurrency. Write five sentences in your own words. Do not copy the page.

---

## Defining structs

A struct is a sequence of named fields. Each field has a name and a type.

```go
type Point struct {
	X int
	Y int
}

type User struct {
	Name string
	Age  int
}
```

`type` binds a name to the struct type. Use a named type when you pass values, write methods, or document meaning. You can also write an anonymous struct type. Named types are the common form.

The zero value of a struct has the zero value in each field. `var p Point` has `X == 0` and `Y == 0`.

Access a field with a dot. `p.X = 3`. Read with `p.X`.

Two named struct types are different even when the fields match. `Point` is not `User`. You can convert between two struct types when the fields match. The names, types, order, and tags must be the same.

A struct is comparable when every field is comparable. You can use `==` on two `Point` values. You cannot use `==` when a field is a slice or a map.

Put one struct definition in a package that owns that concept. Keep field names clear. Use comments on the type when the meaning is not obvious.

An empty struct `struct{}` has no fields. `struct{}{}` is a common token. Its size is zero.

### Questions

#### Theoretical questions

1. What does the zero value of a struct contain?
2. How do you read and write a field?
3. When are two named struct types different?
4. When is a struct type comparable?
5. What is `struct{}` used for?

#### Easy practical tasks

1. Define `Book` with `Title string` and `Pages int`. Set both fields. Print them.
2. Print the zero value of `Book` with `%+v`.
3. Define two structs with the same field list and different type names. Try to assign one to the other. Record the compiler message.
4. Compare two `Point` values with `==`.

#### Medium practical tasks

1. Write `func Mid(a, b Point) Point` that returns the midpoint with integer division.
2. Add a slice field to a struct. Show that `==` no longer compiles.
3. Convert one struct type to another struct type that has the same fields. Print the result.

#### Advanced practical tasks

1. Design three structs for a shop: `Money`, `Item`, and `Cart`. Use named types. Do not add methods yet.
2. Use `struct{}` as the value type of a set-like `map[string]struct{}`. Insert two names. Test membership.

---

## Composite literals

A composite literal builds a value of an array, a slice, a map, or a struct. For a struct, list fields in braces.

```go
p := Point{1, 2}
q := Point{X: 1, Y: 2}
r := Point{X: 1}
```

The first form is a positional literal. The order matches the field order in the type. The second form uses field names. The third form sets `X` and leaves `Y` at the zero value.

Prefer field names. Names stay correct when someone adds a field. A positional literal breaks or changes meaning when the field list changes.

You can take the address of a composite literal.

```go
p := &Point{X: 3, Y: 4}
```

The type of `p` is `*Point`. This form is common in constructors.

You must name the type. `Point{X: 1}` is valid. `{X: 1}` is not valid as a stand-alone value.

Nested literals can omit the inner type in some positions. `[]Point{{1, 2}, {3, 4}}` is valid. Keep the code easy to read. Use names when the nest is deep.

Do not mix positional fields and named fields in one struct literal. The compiler rejects that mix.

Use a pointer literal when the value must be shared or when a method needs a pointer. The next topic covers methods.

### Questions

#### Theoretical questions

1. What is the difference between `Point{1, 2}` and `Point{X: 1, Y: 2}`?
2. What value do omitted named fields get?
3. Why do field names survive a new field better than positions?
4. What is the type of `&Point{X: 1}`?
5. Can one struct literal mix positional fields and named fields?

#### Easy practical tasks

1. Build `User` with a named-field literal. Print it with `%+v`.
2. Build `User` with only `Name` set. Print `Age`.
3. Build `*Point` with `&Point{...}`. Print `p.X`.
4. Build `[]Point` with two inner literals.

#### Medium practical tasks

1. Add a new field to a struct. Fix every positional literal. Then rewrite those literals with names.
2. Write `NewUser(name string) User` that returns a literal with `Age` left at zero.
3. Write a nested struct literal for a type that contains another struct field.

#### Advanced practical tasks

1. Compare `var p Point` plus field assigns with one literal. Write when you choose each style.
2. Build a map `map[string]Point` with literals for two keys. Print both points.

---

## Exported fields

A field name that starts with an upper-case letter is exported. Other packages can read and write that field. A field name that starts with a lower-case letter is unexported. Only the same package can use that name.

```go
type Account struct {
	ID    string
	email string
}
```

`ID` is exported. `email` is unexported. Code in another package can set `a.ID`. That code cannot set `a.email`.

Export a field when callers must set or read it. Keep a field unexported when the type must protect it. Use methods in the next topic to expose controlled access.

`encoding/json` encodes exported fields by default. It ignores unexported fields. A JSON tag does not export a field. The name rule still applies.

The struct type name also follows the export rule. An exported type with unexported fields is a common pattern. Callers use constructors and methods.

Two structs in different packages do not share unexported fields for conversion. Field names must match for conversion. Unexported names include the package identity for this check.

Do not export a field only to fill it in a test in another package. Put the test in the same package or add a constructor.

### Questions

#### Theoretical questions

1. What makes a field exported?
2. Can another package assign to an unexported field?
3. Does `encoding/json` encode unexported fields by default?
4. Can a JSON tag export a field?
5. Why does a type keep some fields unexported?

#### Easy practical tasks

1. Define a struct with one exported field and one unexported field. Set both fields in `package main`.
2. Print the struct with `%+v` and with `json.Marshal`. Compare the outputs.
3. Write four sentences on when you export a field.
4. Rename a field from `name` to `Name`. List what callers can do after the rename.

#### Medium practical tasks

1. Create two packages. Export a type from package `a` with an unexported field. Try to set that field from `package main`. Record the error.
2. Add a constructor `New` in package `a` that sets the unexported field. Call `New` from `main`.
3. Marshal a struct that has exported and unexported fields. Show which fields appear in JSON.

#### Advanced practical tasks

1. Design `Account` so that `email` stays unexported. List the constructor and the exported fields. Do not write methods that the next topic must introduce unless you need them.
2. Explain struct conversion across packages when unexported fields exist. Write five sentences from the specification idea.

---

## Embedding vs named fields

A named field has a name and a type. An embedded field lists a type without a new name. The field name becomes the type name.

```go
type Point struct {
	X int
	Y int
}

type Circle struct {
	Point
	Radius int
}

type Disk struct {
	Center Point
	Radius int
}
```

`Circle` embeds `Point`. The fields `X` and `Y` are promoted. You can write `c.X`. You can still write `c.Point.X`. `Disk` has a named field `Center`. You must write `d.Center.X`. There is no `d.X`.

Embedding is not inheritance. There is no subclass. `Circle` is not a `Point`. A function that takes `Point` does not accept `Circle`.

Promoted fields can conflict. If two embedded types share a field name, you must use the longer selector. The compiler rejects an ambiguous `c.X` when two `X` fields are promoted.

Embedding also promotes methods. The next topic covers method sets. Promotion of methods is a reason to embed a type.

Use embedding when the inner type is a part of the outer type and promotion helps. Use a named field when the inner value is a role such as `Center` or `Owner`.

You can embed a pointer type. `*Point` as an embedded field can be `nil`. A use of a promoted field then panics. Initialize the pointer.

### Questions

#### Theoretical questions

1. How do you write an embedded field?
2. What is field promotion?
3. Is `Circle` a `Point` when `Circle` embeds `Point`?
4. What do you write when promotion is ambiguous?
5. When do you choose a named field instead of embedding?

#### Easy practical tasks

1. Build a `Circle` literal with an inner `Point` literal. Read `c.X` and `c.Radius`.
2. Build a `Disk` with `Center`. Show that `d.X` does not compile.
3. Print `c.Point` and `c.X` for the same value.
4. Embed `Point` and add a field `Name string`. Set all fields in one literal.

#### Medium practical tasks

1. Embed two structs that both have a field `ID`. Show the compiler error for `v.ID`. Use the long selectors.
2. Write a function that takes `Point`. Try to pass a `Circle`. Record the error.
3. Embed `*Point`. Leave it `nil`. Show the panic on `c.X`. Then set the pointer.

#### Advanced practical tasks

1. Redesign `Circle` once with embedding and once with `Center Point`. Write six sentences on call sites and clarity.
2. List three standard library types that embed another type. Write the outer type and the embedded type for each.

---

## Struct tags

A struct tag is a string literal after a field type. The compiler stores the tag. `reflect` reads the tag. Packages such as `encoding/json` use tags.

```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age,omitempty"`
	Pass string `json:"-"`
}
```

The usual form is `` `key:"value"` ``. More than one key can appear in one tag. Separate keys with spaces.

`json:"name"` maps the field to the JSON name `name`. `omitempty` omits the field when the value is the zero value. `json:"-"` skips the field in JSON.

The tag does not change the Go field name. `u.Name` is still the name in Go. The tag is metadata.

Use backticks for the tag string. A raw string can hold double quotes around the value.

`encoding/json` only sees exported fields. A tag on an unexported field does not encode that field.

Keep tags short. Use the names that the wire format needs. Do not invent a new tag key unless you write the code that reads it.

Other packages use tags too. Examples: `xml`, `yaml` in extra modules, and database mappers. This topic uses `json` as the main example.

Invalid tag syntax is still a string. The json package may ignore a bad value. Check the official json tag rules when you add options.

### Questions

#### Theoretical questions

1. Where do you write a struct tag?
2. What package reads tags at run time?
3. What does `json:"name"` change?
4. What does `json:"-"` mean?
5. Does a tag export an unexported field?

#### Easy practical tasks

1. Define `User` with `json:"name"` and `json:"age"`. Marshal one value. Print the JSON.
2. Add `omitempty` on `Age`. Marshal a user with `Age` 0. Print the JSON.
3. Add `json:"-"` on a field. Confirm that the field is absent in JSON.
4. Print a field tag with `reflect` for one field. Use the course docs or `go doc reflect.StructTag`.

#### Medium practical tasks

1. Unmarshal `{"name":"Ada","age":36}` into `User`. Print the struct.
2. Use a different JSON name than the Go name. Show both names in a short table.
3. Add a second key in one tag, for example a comment in your notes plus `json`. Show that json still works.

#### Advanced practical tasks

1. Write `User` with `Pass` hidden from JSON and a constructor that sets `Pass`. Marshal and confirm the secret is absent.
2. Read the `encoding/json` docs for tag options. List four options. Write one sentence for each option.

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
