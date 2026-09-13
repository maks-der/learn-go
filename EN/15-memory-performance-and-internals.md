# 15. Memory, Performance, and Internals

## Description

This topic explains how Go stores values, how the garbage collector reclaims memory, and how the compiler decides where a value lives. You write faster programs when you understand slices, maps, strings, allocation, and alignment. You also learn the limits of `unsafe` and the limits of compiler optimizations.

Use one term for each concept. Allocation is the act of reserving memory. Escape analysis is the compiler pass that decides stack versus heap. The garbage collector (GC) reclaims heap objects that the program can no longer reach. Do not treat micro-benchmarks as proof without `testing.B` and a clear hypothesis.

---

## Stack vs heap; escape analysis

Each goroutine has a stack. The stack holds local variables that do not escape. Stack allocation is cheap. The runtime grows and shrinks the stack as the goroutine needs more space.

The heap holds objects that live after the function returns, or that the compiler cannot keep on the stack. The garbage collector manages the heap. Heap allocation is more expensive than stack allocation. Heap objects add GC work.

Escape analysis is the compiler pass that decides the location. The compiler does not follow a simple "new means heap" rule. A value that you create with `&T{}` may stay on the stack if nothing retains the pointer. A local variable may move to the heap if you return its address.

Common reasons for a value to escape to the heap:

- you return a pointer to the value
- you store the value in an interface value (the interface may need a heap object)
- you send the value on a channel that outlives the function
- you store the pointer in a struct or slice that already lives on the heap
- you close over the variable in a goroutine that runs after the function returns

Inspect escape analysis with:

```text
go build -gcflags=-m
```

Add `-m` twice for more detail: `-gcflags="-m -m"`. The compiler prints lines such as `moved to heap` and `escapes to heap`. Read those lines as hints. They are not a stable public API. The wording can change between Go versions.

```go
func makeInt() *int {
    x := 42
    return &x // x escapes
}

func sum(a, b int) int {
    return a + b // a and b stay on the stack
}
```

Do not force heap allocation without a reason. Do not hide pointers from the compiler. Write clear code. Profile first. Then change allocation only when a profile shows a cost.

A large local array can also be expensive on the stack. A very large stack frame can cause stack growth. For large buffers, use `make([]byte, n)` and accept a heap slice, or reuse a buffer.

The stack is not shared between goroutines. Do not take the address of a local variable and pass that pointer to another goroutine unless the lifetime is clear. The value may move during stack growth. The runtime updates pointers that the GC and the compiler know about. Do not hide pointers from the runtime (see `unsafe`).

### Questions

#### Theoretical questions

1. What does the stack store for a goroutine?
2. What is escape analysis?
3. Name three reasons for a value to escape to the heap.
4. Why is heap allocation more expensive than stack allocation?
5. Does `new` or `&` always allocate on the heap? Explain.
6. Why is `-gcflags=-m` output a hint and not a stable API?

#### Easy practical tasks

1. Write `makeInt` and `sum` as in this section. Run `go build -gcflags=-m` on the file. Copy the escape lines.
2. Write a function that returns an `int` by value. Confirm that the compiler does not move that result to the heap as a pointer.
3. Write a function that stores an `int` in a variable of type `any`. Read the escape output.
4. Draw two boxes labeled stack and heap. Place a local `int` and a returned `*int` in the diagram.

#### Medium practical tasks

1. Compare a function that returns `T` with a function that returns `*T` for a small struct. Use `-gcflags=-m` and a benchmark. Write which version allocates.
2. Start a goroutine that reads a local variable by pointer. Show the escape line. Then pass the value on a channel instead. Compare the output.
3. Allocate a large `[1 << 20]byte` on the stack in a function. Note stack growth or a better design with a slice. Write your observation.

#### Advanced practical tasks

1. Take a small package. Run `-gcflags=-m` and list every escape. Mark each escape as required by the API or as avoidable. Change one avoidable escape. Show the new output.
2. Explain stack growth: what happens to pointers to stack objects when the runtime copies the stack. Use the official documentation. Write eight short sentences.

---

## Garbage collector (high-level: concurrent, tri-color)

Go uses a concurrent garbage collector. The GC runs at the same time as your goroutines for most of the work. The GC still uses short stop-the-world (STW) pauses for some phases. Those pauses are small on current Go versions. You must not depend on a specific pause length.

The collector uses a tri-color mark algorithm. Each object has a color during a mark cycle:

- white: not yet proven reachable
- gray: reachable, but the GC has not scanned all pointers inside it
- black: reachable, and the GC has scanned its pointers

The GC starts from roots. Roots include global variables, stacks, and registers. The GC marks reachable objects. At the end of the cycle, white objects are garbage. The GC then sweeps unused memory and returns it to the allocator.

Because the program mutates the heap during a concurrent mark, the runtime uses a write barrier. The write barrier records pointer updates so that the GC does not miss an object. You do not call the write barrier yourself. The compiler inserts it.

Go GC is not a generational collector in the classic young/old sense. It is a non-moving, concurrent mark-sweep collector. Object addresses do not change because of GC. You can pass a pointer to C through cgo for the lifetime that you pin or keep reachable. Still treat cgo as a special topic.

The GC pacer decides when to start a cycle. The runtime looks at heap size and allocation rate. You can set a soft goal with `GOGC` (default 100). `GOGC=100` means the heap may grow to about double the live size after a cycle before the next cycle. `GOGC=off` disables GC. Do not disable GC in production.

`GOMEMLIMIT` (Go 1.19 and later) sets a soft memory limit. The runtime tries to keep the total memory near that limit. Use it when the process must stay inside a container limit. Do not use it as a substitute for a correct program.

`runtime.GC()` starts a garbage collection cycle and waits for it to finish. Use it in tests or benchmarks when you need a known state. Do not call it on a request path that runs often.

`runtime.ReadMemStats` and `testing.B` report allocations. Profiles with `pprof` show heap objects. Measure before you tune GC.

### Questions

#### Theoretical questions

1. What does "concurrent garbage collector" mean for your goroutines?
2. What are the three colors in tri-color marking?
3. What is a GC root? Give two examples.
4. Why does a concurrent GC need a write barrier?
5. What does `GOGC=100` mean?
6. Is the Go GC generational? Does it move objects in memory?

#### Easy practical tasks

1. Write a program that allocates many slices in a loop. Print `HeapAlloc` from `runtime.ReadMemStats` before and after the loop.
2. Run the same program with `GOGC=50` and `GOGC=200`. Write the two `GOGC` values and one sentence about heap size.
3. Call `runtime.GC()` in a test after you drop references to a large slice. Write why the test does that.
4. Open the Go GC guide on go.dev. Write the definition of a STW pause in one sentence.

#### Medium practical tasks

1. Write a benchmark that allocates one `[]byte` of size 4096 per iteration. Report `B.ReportAllocs`. Record bytes per operation.
2. Explain write barriers in six sentences. Use the official GC design notes or the Go blog. Do not copy long passages.
3. Set `GOMEMLIMIT` for a small program that allocates in a loop. Record whether the process stays near the limit. Restore the environment.

#### Advanced practical tasks

1. Use `go tool pprof` on a heap profile of a program that leaks slices (keep them in a global). Identify the allocation site.
2. Read about the GC pacer. Write a one-page note: when a cycle starts, what `GOGC` changes, and what `GOMEMLIMIT` changes. Use current Go documentation.

---

## Slice, map, and string internals

A slice is a descriptor with three fields: a pointer to an array, a length, and a capacity. The descriptor is a small value. Copy of a slice copies the descriptor. The copy shares the same backing array.

```go
s := make([]int, 2, 4) // pointer, len=2, cap=4
t := s                 // same array, same len, same cap
t[0] = 9               // s[0] is also 9
```

`append` may write into unused capacity. If capacity is full, `append` allocates a new array, copies the elements, and returns a new descriptor. The old array stays alive if any slice still points to it.

A string is a descriptor with two fields: a pointer to bytes and a length. A string does not have capacity. The bytes are immutable. Two strings may share a backing array. A substring `s[a:b]` may share bytes with `s`.

A map is a hash table. From Go 1.24 the runtime implements maps with Swiss tables. You still use the same language operations: make, insert, lookup, delete. The header that you copy is a pointer-like value. Copy of a map value shares the same table. A nil map allows reads and returns the zero value. A write to a nil map panics.

Map iteration order is random. Do not depend on a sequence of keys. The hash seed and the table layout can change between versions.

```text
slice:  [ ptr | len | cap ] --> array elements
string: [ ptr | len ]       --> immutable bytes
map:    pointer --> hash table (Swiss tables in Go 1.24 and later)
```

`len` and `cap` on a slice are `O(1)`. `len` on a map is `O(1)`. Lookup in a map is average `O(1)`. Worst-case behavior depends on keys and the hash implementation. Do not treat a map as an ordered list.

A slice of structs stores the structs in one array. A slice of pointers stores pointers. The pointer slice can be smaller when structs are large, but it adds heap objects and pointer scanning for the GC.

Clear a slice or a map with the built-in `clear` (Go 1.21 and later). `clear(s)` zeros the elements of slice `s` and keeps length. `clear(m)` deletes all keys of map `m`.

### Questions

#### Theoretical questions

1. What three fields does a slice descriptor hold?
2. What happens to the backing array when `append` grows a slice?
3. How many fields does a string descriptor hold?
4. What map implementation does Go 1.24 and later use?
5. Why does assignment of a map variable not copy the keys and values?
6. What does `clear` do on a slice and on a map?

#### Easy practical tasks

1. Create a slice with `len` 2 and `cap` 4. Append one value without a new array. Print `len` and `cap` before and after.
2. Share a backing array between two slices. Change an index in one slice. Print the other slice.
3. Print `len` of a string that contains a multi-byte UTF-8 character. Write why `len` is not the number of runes.
4. Write a map, copy the map variable, insert a key through the copy, and print the original map.

#### Medium practical tasks

1. Force `append` to allocate a new array. Prove the new allocation with `cap` change or with `&s[0]` before and after (same element index only when the array did not move).
2. Compare a `[]Point` and a `[]*Point` for 10,000 points in a benchmark. Report allocations.
3. Show that range over a map produces a different order across two runs. Print keys from two loops.

#### Advanced practical tasks

1. Write a function that appends to a slice and returns it. Document when the caller still sees the old backing array. Add a test that fails if you ignore the return value of `append`.
2. Read the Go blog or release notes on Swiss tables. Write eight sentences: what changed for the programmer, and what did not change.

---

## String immutability and `[]byte` conversions

A string cannot change its bytes. A statement such as `s[0] = 'A'` does not compile. Immutability lets the runtime share backing arrays between strings. It also makes strings safe to copy as descriptors.

Conversion between `string` and `[]byte` copies the bytes in the safe, supported form:

```go
b := []byte("hello") // copy of the bytes
s := string(b)       // copy of the bytes
```

A copy costs time and memory. In a hot loop, repeated conversions allocate. Prefer `[]byte` for buffers that you mutate. Prefer `string` for keys, text APIs, and values that you do not change.

`strings.Builder` builds a string with a growing buffer. `Builder.String` may share the buffer. Do not use the builder after you call `String` if you still write to the builder. Read the `strings.Builder` documentation for the current rules.

`bytes.Buffer` is a writer that holds bytes. Use it when you need `io.Writer`. Use `strings.Builder` when you only build a string.

The `unsafe` package offers `unsafe.String` and `unsafe.Slice` conversions that do not copy. Those conversions are valid only when you keep the immutability rules. If you convert a `[]byte` to a `string` without a copy, you must not change the slice after that. A later write changes the string in place and breaks the language memory model. Do not use these conversions until you can state the lifetime rules. Prefer a normal conversion.

Range over a string yields runes (Unicode code points) and byte indices. Indexing a string yields a byte. Conversion of a string to `[]rune` copies code points. Choose the form that matches the job.

```go
s := "Gö"
fmt.Println(len(s))      // bytes
fmt.Println(len([]rune(s))) // runes
```

Do not store secrets in strings if you need to erase them. You cannot overwrite string bytes. Use `[]byte` and overwrite the slice when you must clear a secret.

### Questions

#### Theoretical questions

1. Why does the language forbid `s[0] = 'A'` on a string?
2. What is the cost of `[]byte(s)` in the safe conversion?
3. When do you choose `strings.Builder` instead of `bytes.Buffer`?
4. What rule must you follow if you convert a slice to a string without a copy?
5. What does a range over a string produce at each step?
6. Why is a string a poor type for a secret that you must erase?

#### Easy practical tasks

1. Convert a string to `[]byte`, change one index, convert back, and print both values.
2. Build a string with `strings.Builder` from three `WriteString` calls. Print the result.
3. Print `len(s)` and `utf8.RuneCountInString(s)` for a string with `é`.
4. Write a function `HasPrefixBytes(b []byte, prefix string) bool` that does not convert the full slice to a string. Use `bytes.HasPrefix` and `[]byte(prefix)` or `strings.HasPrefix(string(b), prefix)` and then remove the extra copy if you can.

#### Medium practical tasks

1. Benchmark `string(b)` in a loop versus reuse of a string that you create once. Report allocations.
2. Write a function that appends to a `[]byte` buffer and then converts to string once at the end. Compare it to `s = s + part` in a loop.
3. Use `unsafe.String` on a byte slice, then change the slice. Record the unsafe result. Then rewrite the code with a safe copy. State why the unsafe version is invalid.

#### Advanced practical tasks

1. Implement a tokenizer that works on `[]byte` and returns strings only for tokens that the caller stores. Avoid a conversion per byte.
2. Read the documentation of `unsafe.String` and `unsafe.StringData`. Write the preconditions in your own words. Write one valid example and one invalid example.

---

## Allocation-aware coding (`sync.Pool`, pre-size slices)

Allocation-aware coding reduces heap allocations in hot paths. You still write clear code first. You change allocation after a profile or a benchmark shows a cost.

Pre-size slices when you know the length or a good capacity:

```go
out := make([]T, 0, len(in))
for _, v := range in {
    if keep(v) {
        out = append(out, v)
    }
}
```

`make([]T, 0, n)` avoids many growth allocations. `make([]T, n)` is correct when you set every index. Do not guess a huge capacity that you never use. A large unused cap holds a large array alive.

Reuse buffers. A function that encodes JSON or reads a file can accept a `[]byte` and return a sliced prefix:

```go
func appendJSON(buf []byte, v any) ([]byte, error) {
    // encode into buf, grow if needed, return the new slice
    return buf, nil
}
```

`sync.Pool` stores temporary objects for reuse. The pool may drop objects at any time, including during GC. Get may return a new object. You must reset the object before reuse.

```go
var bufPool = sync.Pool{
    New: func() any {
        b := make([]byte, 0, 1024)
        return &b
    },
}

func withBuf(fn func(buf *[]byte)) {
    p := bufPool.Get().(*[]byte)
    *p = (*p)[:0]
    defer bufPool.Put(p)
    fn(p)
}
```

Rules for `sync.Pool`:

- use it for short-lived, high-churn objects
- do not use it as a cache of important data
- reset every object that you get
- put the same type that you get
- do not assume an object stays in the pool

Other allocation tips:

- prefer `strconv.AppendInt` and `bytes` append helpers over format that allocates a new string
- avoid `+` on strings in a loop
- pass a pointer to a large struct when the function does not need a copy
- do not store pointers only to avoid a copy when the struct is small; the heap object and GC scan can cost more
- use `slices.Grow` or `slices.Clip` when those helpers match the job

Measure with:

```text
go test -bench=. -benchmem
```

Do not add `sync.Pool` only because an article mentioned it. Add it when `benchmem` and profiles show allocation of the same object shape on a path that runs often.

### Questions

#### Theoretical questions

1. Why does a pre-sized slice reduce allocations during `append`?
2. What does `sync.Pool` guarantee about stored objects?
3. Why must you reset an object that you get from a pool?
4. When is `make([]T, n)` better than `make([]T, 0, n)`?
5. Why is `sync.Pool` a bad place to store session data?
6. What flag reports allocations in a benchmark?

#### Easy practical tasks

1. Filter a slice of 1000 integers into a new slice with a pre-sized `make`. Print `cap`.
2. Write the same filter with `var out []T` and no cap. Benchmark both versions with `-benchmem`.
3. Write a `sync.Pool` of `[]byte`. Get, append, put. Reset length on get.
4. Replace a loop of `s = s + x` with `strings.Builder`. Compare allocations.

#### Medium practical tasks

1. Write `appendIntCSV(buf []byte, nums []int) []byte` with `strconv.AppendInt`. Benchmark against `fmt.Sprintf`.
2. Add a pool to a small HTTP handler that builds a JSON object in a buffer. Keep the handler correct under concurrent requests.
3. Use `slices.Clip` after you build a slice that you keep for a long time. Explain why clip can help the GC.

#### Advanced practical tasks

1. Profile a program that encodes many small JSON objects. Reduce allocations with a pool or with `json.Encoder` on a reused buffer. Show before and after `benchmem`.
2. Implement a growable buffer type with `Grow` and `Reset`. Add a pool of that type. Write tests for concurrent get and put. Prove that a missed reset would leak data between users (write a failing test, then fix).

---

## `unsafe` package (know the dangers)

The `unsafe` package lets you break the type system. The compiler and the GC assume that you follow the documented rules. If you break those rules, the program can crash, corrupt memory, or miscompile after a Go upgrade.

`unsafe.Pointer` is a pointer type that can convert to and from other pointer types. Conversion is valid only in the cases that the documentation lists. The common legal patterns are:

- convert `*T` to `unsafe.Pointer` and back to `*T`
- convert `unsafe.Pointer` to `uintptr` for arithmetic, then immediately back to `unsafe.Pointer` in the same expression
- use `unsafe.Sizeof`, `unsafe.Alignof`, and `unsafe.Offsetof` for layout facts
- use `unsafe.String`, `unsafe.StringData`, `unsafe.Slice`, and `unsafe.SliceData` with their documented preconditions

```go
n := unsafe.Sizeof(int32(0)) // 4 on current ports
```

Do not store a `uintptr` and use it later as a pointer. The GC does not treat `uintptr` as a pointer. The object can be collected. The integer then does not keep the object alive. Do not hide a pointer in an integer.

Do not use `unsafe` to read another package's unexported fields. Layout can change. The `go vet` tool reports some bad `unsafe` patterns. `vet` does not catch every error.

`unsafe` does not make Go a systems language with manual free. You still have a garbage collector. You still must keep objects reachable while you use pointers to them.

Valid reasons to read `unsafe` code:

- you study the standard library
- you write a package that must interoperate with a C ABI and you accept the cost
- you use a well-reviewed helper such as `sync/atomic` pointer helpers that wrap unsafe internally

Invalid reasons:

- you want a program to run faster without a profile
- you want to avoid a copy that the language already makes cheap
- you want to mutate a string

For this handbook, the required skill is to recognize danger, not to invent new `unsafe` conversions. If a task asks you to use `unsafe`, stay inside the documented conversions.

### Questions

#### Theoretical questions

1. What does the GC assume about pointers that you create with `unsafe`?
2. Why must you not store a pointer in a `uintptr` for later use?
3. What do `Sizeof`, `Alignof`, and `Offsetof` return?
4. Name one legal conversion and one illegal conversion that involve `unsafe.Pointer`.
5. Does `go vet` find every `unsafe` mistake?
6. Why is mutation of a string through `unsafe` a violation of the language rules?

#### Easy practical tasks

1. Print `unsafe.Sizeof` for `bool`, `int32`, `int64`, `string`, and a slice. Record the numbers on your machine.
2. Print `unsafe.Offsetof` for each field of a small struct.
3. Read `go doc unsafe`. Write the rule about `uintptr` in your own words.
4. Run `go vet` on a file that converts a pointer to `uintptr`, stores it, and converts back later. Record whether `vet` warns.

#### Medium practical tasks

1. Use only documented `unsafe.Slice` to view an array as a slice. Do not change bounds beyond the array. Print the slice.
2. Find one use of `unsafe` in the standard library source (for example in `reflect` or `runtime`). Write what the comment says about safety. Do not copy a large block.
3. Write a short report: three bugs that `unsafe` can cause after a Go version upgrade (layout change, GC, compiler).

#### Advanced practical tasks

1. Implement a safe wrapper that uses `unsafe.String` to share bytes with a string, with tests that prove you never mutate the slice after the conversion. Document the lifetime.
2. Read the Go `unsafe` package documentation and the Go compatibility statement. Write when the Go team can change behavior that `unsafe` users depend on.

---

## Alignment and struct padding

The CPU and the ABI require alignment. A type of size 8 often needs an address that is a multiple of 8. The compiler inserts padding bytes so that each field starts at a correct offset.

```go
type Bad struct {
    A bool  // 1 byte
    // 7 bytes padding on a typical 64-bit port
    B int64 // 8 bytes
}

type Good struct {
    B int64
    A bool
    // 7 bytes padding at the end so that the struct size is a multiple of 8
}
```

`Good` still has padding at the end. The size of a struct is a multiple of the alignment of the struct. An array of `Good` then keeps each element aligned.

Field order changes padding. Put fields from largest alignment to smallest when you care about size. Do not reorder exported fields of a struct that you encode with a binary format that depends on layout. Do not reorder fields only for style when a memory profile does not ask for it.

`unsafe.Alignof(v)` is the alignment of `v`. `unsafe.Offsetof(s.Field)` is the offset of the field. `unsafe.Sizeof(v)` is the size including padding.

Atomic operations and some hardware accesses require aligned addresses. The `sync/atomic` types such as `atomic.Int64` handle alignment for you. If you embed `int64` and use `atomic.AddInt64` on a field, the field must be 64-bit aligned. On 32-bit ports this is a real risk. Prefer `atomic.Int64`.

A bool is one byte. A pointer is 4 or 8 bytes by architecture. A string is two pointer-sized words. A slice is three pointer-sized words. These sizes explain why a struct of many small fields plus one pointer can waste space if you order fields poorly.

Empty structs `struct{}` have size zero. A field of type `struct{}` can share an address with the next field. This is useful in a `map[K]struct{}` set. Do not add a zero-size field in the middle of a struct unless you know the layout effect.

### Questions

#### Theoretical questions

1. What is alignment?
2. Why does the compiler insert padding between fields?
3. Why does a struct often have padding at the end?
4. How can field order reduce the size of a struct?
5. Why prefer `atomic.Int64` over `atomic.AddInt64` on a struct field?
6. What is the size of `struct{}` and why does a set use `map[K]struct{}`?

#### Easy practical tasks

1. Define `Bad` and `Good` as in this section. Print `Sizeof` for both.
2. Reorder fields of a struct with `bool`, `int32`, `int64`, and `byte`. Find an order with a smaller size.
3. Print `Alignof` for `byte`, `int32`, `int64`, and `string` on your machine.
4. Draw the layout of a struct with two `bool` fields and one `int64` field. Mark padding.

#### Medium practical tasks

1. Compare `[]Bad` and `[]Good` for one million elements. Estimate memory from `Sizeof` times length. Confirm with a heap profile or `ReadMemStats`.
2. Write a struct that you encode with `encoding/binary` and `binary.Size`. Show that padding is not the same as the encoded form unless you use a packed format. Explain the difference.
3. Put `atomic.Int64` and a `bool` in a struct. Print offsets. State why the atomic type is safer.

#### Advanced practical tasks

1. Design a struct for a cache entry: a `uint64` id, two `bool` flags, a `string` name, and a `[]byte` payload. Minimize size without breaking a stable JSON tag set. Record `Sizeof` before and after.
2. Read the Go specification on size and alignment. Write five rules that a beginner can apply. Do not copy the specification text.

---

## Compiler optimizations you can rely on (and those you cannot)

The Go compiler applies optimizations. Some are stable enough that you can write code in a normal style and expect a good result. Others change between versions. Do not write obscure code to force an optimization.

Optimizations you can usually rely on:

- **inlining of small functions** — a short function often inlines into the caller; keep functions small and simple
- **escape analysis** — values that do not escape stay on the stack
- **bounds-check elimination in simple loops** — a loop of the form `for i := 0; i < len(s); i++` with `s[i]` often loses repeated bounds checks
- **dead-code elimination** — code that is not reachable after constant folding may disappear
- **constant folding** — `2 * 3` becomes `6` at compile time
- **devirtualization in limited cases** — a concrete type stored in an interface in a way the compiler sees can become a direct call

Write normal, clear loops. Do not unroll loops by hand unless a benchmark proves a gain. The compiler and the CPU already handle many simple loops.

Optimizations you must not rely on:

- a specific function always inlines (a comment, an interface call, or a large body can block inlining)
- a specific allocation always stays on the stack
- map iteration order or map growth timing
- GC cycle timing
- a given `append` growth factor (the growth policy can change)
- exact instruction selection or register use
- the same `gcflags=-m` text after an upgrade
- `sync.Pool` keeping an object across a GC
- time.Now monotonic details for security (use `crypto/rand` for secrets)

The ABI and pointer bit patterns are not a user-level contract. Do not parse your own stack. Do not depend on the size of `int` except through `int` itself (`int` is at least 32 bits; on 64-bit ports it is 64 bits).

Build modes change the result. `-race` adds extra memory and slower code. `-gcflags=all=-l` disables inlining. A debug build with extra checks is not the production binary. Benchmark the same flags that you ship.

You can read inlining and escape notes:

```text
go build -gcflags="-m"
```

You can disassemble:

```text
go build -o app .
go tool objdump -s SomeFunc app
```

Use `objdump` to learn, not to lock your source to one instruction sequence.

A practical rule: write the clear version. Measure. Change the algorithm (less work, fewer allocations, better data layout). Do not change the algorithm into an unreadable form to chase a compiler quirk.

### Questions

#### Theoretical questions

1. Name three optimizations that you can usually trust in current Go.
2. Name three behaviors that you must not treat as a contract.
3. Why can an interface method call block a direct inlined call?
4. How does the race detector change what a benchmark measures?
5. Why is hand-unrolling a loop a weak first step?
6. What is bounds-check elimination?

#### Easy practical tasks

1. Write a small function and a large function. Run `-gcflags=-m` and see which one inlines.
2. Build a program with and without `-race`. Compare binary size or a benchmark result.
3. Write a loop `for i := 0; i < len(s); i++` that sums `s[i]`. Keep the code clear. Do not add extra checks.
4. List four items from this section that are not part of the Go compatibility promise for your source code.

#### Medium practical tasks

1. Benchmark a small helper that you expect to inline versus the same logic copied in the caller. Use `-benchmem`. Write whether the results match your expectation.
2. Disable inlining with `-gcflags=all=-l` for a benchmark. Compare times. Restore normal flags.
3. Show that map iteration order is not an optimization that you can rely on. Write a test that sorts keys when the test needs a stable order.

#### Advanced practical tasks

1. Use `go tool objdump` on a function that indexes a slice in a tight loop. Look for bound checks. Write what you see in plain language. Do not claim a guarantee for the next Go version.
2. Take a slow function. Apply one algorithmic change (pre-size, reuse buffer, or fewer interface conversions). Prove the gain with a benchmark. List two compiler-specific changes that you refused to use.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the life of a local variable from allocation to reuse of its memory. Cover stack, heap, escape analysis, and GC colors.
2. How do the descriptors of a slice, a string, and a map differ, and what does assignment copy in each case?
3. When do you convert between `string` and `[]byte`, and when do you keep one type for the whole pipeline?
4. What rules keep `sync.Pool` and `unsafe` from becoming silent production bugs?
5. A teammate reorders struct fields, adds `unsafe`, and sets `GOGC=off` for speed. Which facts do you use to reject each change?

#### Easy practical tasks

1. Create a module `example.com/memx`. Add a slice growth demo, a string conversion demo, and a struct padding demo. Run `go test` and `go build -gcflags=-m`.
2. Write a one-page cheat sheet: escape reasons, tri-color colors, slice header fields, `clear`, `GOGC`, `sync.Pool` reset, and one `unsafe` ban.
3. Run `go test -bench=. -benchmem` on a function that builds a string with `+` and a function that uses `strings.Builder`. Save the output.
4. Draw the tri-color states of three objects during a mark cycle. Label one root.

#### Medium practical tasks

1. Write a program that leaks by appending pointers to a global slice. Then write a version that processes items and drops them. Compare `HeapAlloc` after `runtime.GC()`.
2. Minimize padding in a struct that holds `uint8`, `uint64`, `uint16`, and `uint32`. Print `Sizeof` for two field orders.
3. Document a performance change in ten steps: measure, hypothesize, change one thing, measure again. Apply the steps to one function.

#### Advanced practical tasks

1. Build a small JSON encoder path that reuses a buffer from a pool and pre-sizes the slice. Compare it to `json.Marshal` in a benchmark. State when the extra code is not worth it.
2. Read the current Go GC guide and the `unsafe` package docs. Write a one-page policy for your team: allowed measurements, allowed pools, banned `unsafe` patterns, and required review for struct layout changes in exported binary formats.
