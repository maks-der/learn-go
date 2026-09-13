# 10. Standard Library Essentials

## Description

The Go standard library is a set of packages that ship with the toolchain. You import these packages without `go get`. This topic shows the packages that most programs need: text, numbers, time, files, encoding, HTTP clients, logs, context, command-line flags, tests, and locks.

Complete topics 1 to 9 before this topic. Use one term for each package. Prefer the standard library when it solves the problem. This topic uses Go 1.22 or later. Use `log/slog` for structured logs. Use `math/rand/v2` for new random numbers. Do not use `io/ioutil` in new code. Topic 11 covers concurrency in depth. Topic 12 covers `context` in depth. Topic 13 covers tests in depth. This topic gives the first working view of those packages.

---

## `fmt` — printing and scanning

Package `fmt` formats values as text and reads text into values. The common print functions are:

- `Print`, `Println`, `Printf` write to standard output
- `Fprint`, `Fprintln`, `Fprintf` write to an `io.Writer`
- `Sprint`, `Sprintln`, `Sprintf` return a `string`

`Println` adds spaces between operands and adds a newline. `Printf` uses a format string. Common verbs:

- `%v` default format
- `%+v` default format with field names for structs
- `%#v` Go syntax
- `%T` type
- `%s` string
- `%q` quoted string
- `%d` integer
- `%f` floating-point
- `%p` pointer
- `%w` wrap an error in `fmt.Errorf` only

`fmt.Errorf` returns an `error`. See topic 8 for `%w`.

The scan functions read text:

- `Scan`, `Scanln`, `Scanf` read from standard input
- `Fscan`, `Fscanln`, `Fscanf` read from an `io.Reader`
- `Sscan`, `Sscanln`, `Sscanf` read from a `string`

Scan functions return the number of items and an `error`. Check the error. A scan from a user often fails. Prefer `bufio.Scanner` or `strconv` when you parse a single line or a single number.

Use `fmt` for humans and for simple logs. Use `strconv` when you convert one number or one bool. Use `encoding/json` when you exchange structured data.

### Questions

#### Theoretical questions

1. What is the difference between `fmt.Println` and `fmt.Printf`?
2. Which function writes formatted text to an `io.Writer`?
3. What does the verb `%T` print?
4. What two results does `fmt.Sscanf` return?
5. When must you use `strconv` instead of `fmt` for a conversion?

#### Easy practical tasks

1. Print your name and an integer with `fmt.Printf` and verbs `%s` and `%d`.
2. Use `fmt.Sprintf` to build the string `id=42` without printing.
3. Print a struct with `%v`, `%+v`, and `%#v`. Write how the three lines differ.
4. Use `fmt.Fprintln` to write one line into a `strings.Builder`.

#### Medium practical tasks

1. Scan two integers from a string with `fmt.Sscanf`. Check the count and the error for a good string and for a bad string.
2. Write a `Printf` call with a wrong verb for the operand type. Run `go vet`. Record the message.
3. Compare `fmt.Sprint(3.14)` with a `strconv` format of the same number. Write which package you pick for a stable machine-readable number.

#### Advanced practical tasks

1. Implement a type with method `Format(f fmt.State, verb rune)` so that `%s` and `%v` print different text. Prove both verbs.
2. Read the source of `fmt.Println` with `go doc -src`. Write how it reaches `Fprintln` and `os.Stdout`.

---

## `strconv` — string conversions

Package `strconv` converts between basic types and their string forms. Use this package when the input or the output is a single value.

Integer conversions:

- `Atoi` parses an `int` from a decimal string
- `Itoa` formats an `int` as a decimal string
- `ParseInt` and `ParseUint` take a base and a bit size
- `FormatInt` and `FormatUint` take a base

Other conversions:

- `ParseFloat`, `FormatFloat`
- `ParseBool`, `FormatBool`
- `Quote`, `Unquote` for Go string literals
- `AppendInt` and the other `Append*` functions write into a `[]byte`

Every `Parse*` function returns an `error`. Check the error. Do not ignore a parse error. A bad string is normal input, not a panic.

`Atoi` is a short form of `ParseInt` with base 10 and bit size 0. Bit size 0 means the size of `int` on the current machine.

`FormatFloat` needs a format byte (`'f'`, `'e'`, `'g'`), a precision, and a bit size. Pick a precision on purpose. Do not use `fmt.Sprintf("%f", x)` when you need a fixed machine format. Use `strconv`.

`Quote` adds quotes and escapes. Use it when you print a string that can contain spaces or quotes.

### Questions

#### Theoretical questions

1. What is the difference between `strconv.Atoi` and `strconv.ParseInt`?
2. Why must you check the error from `ParseFloat`?
3. What does bit size `0` mean in `ParseInt`?
4. When do you use `AppendInt` instead of `Itoa`?
5. What does `strconv.Quote` add to a string?

#### Easy practical tasks

1. Convert `"17"` to `int` with `Atoi`. Print the value and the error.
2. Convert `17` to a string with `Itoa`. Print the string.
3. Parse `"true"` and `"yes"` with `ParseBool`. Record both results.
4. Quote the string `say "hi"` with `strconv.Quote`. Print the result.

#### Medium practical tasks

1. Parse `"ff"` as a hex `int64` with `ParseInt`. Then format that value as hex with `FormatInt`.
2. Use `FormatFloat` with `'f'`, precision `2`, and bit size `64`. Compare the text with `fmt.Sprintf("%.2f", x)`.
3. Feed `Atoi` an empty string, a word, and a number that does not fit in `int`. Record the three errors.

#### Advanced practical tasks

1. Write a function that parses a list of integers from a slice of strings and returns one `error` that names the first bad index. Use `strconv`, not `fmt.Sscanf`.
2. Build a decimal string with `AppendInt` in a loop for the numbers `1` to `5`. Show the final `[]byte`. Explain why append can be cheaper than `Itoa` plus `+`.

---

## `strings`, `bytes`, `unicode`, `regexp`

These packages work with text. Use `string` when you do not need to change bytes in place. Use `[]byte` when you need a mutable buffer. Package `strings` and package `bytes` have many of the same function names.

Common `strings` functions: `Contains`, `HasPrefix`, `HasSuffix`, `Split`, `Join`, `ReplaceAll`, `TrimSpace`, `ToLower`, `ToUpper`, `Cut`, `CutPrefix`, `CutSuffix`. Type `strings.Builder` builds a string with few allocations. Type `strings.Reader` implements `io.Reader` over a string.

Package `bytes` adds `Buffer`. `bytes.Buffer` is a writable byte buffer. Use it when you collect bytes for an `io.Writer`.

Package `unicode` classifies runes. Examples: `unicode.IsLetter`, `unicode.IsDigit`, `unicode.IsSpace`, `unicode.ToUpper`. A `rune` is a Unicode code point. Range a string to see runes, not bytes. Package `unicode/utf8` encodes and decodes UTF-8.

Package `regexp` implements regular expressions. The engine is RE2. The engine does not support backreferences. `regexp.Compile` returns an error. `regexp.MustCompile` panics. Use `MustCompile` only for a pattern that is a constant in the source (topic 8).

Prefer `strings` when you search for a fixed substring. A regular expression is slower. A wrong pattern is a common error. Compile a pattern once. Reuse the `*regexp.Regexp` value.

### Questions

#### Theoretical questions

1. When do you use `[]byte` instead of `string`?
2. What problem does `strings.Builder` solve?
3. What does `unicode.IsLetter` accept as input?
4. Why does RE2 omit backreferences?
5. When must you use `strings.Contains` instead of `regexp`?

#### Easy practical tasks

1. Use `strings.Cut` on `name=value`. Print the two parts and the found flag.
2. Join `[]string{"a", "b", "c"}` with `strings.Join` and a comma.
3. Use `bytes.Contains` on a `[]byte` slice. Print the bool.
4. Call `unicode.IsDigit` on `'7'` and on `'a'`. Print both results.

#### Medium practical tasks

1. Build a 3-line text with `strings.Builder`. Then split the result on newlines.
2. Compile a pattern that matches an email-like string with `regexp.Compile`. Return the compile error to the caller for a bad pattern.
3. Range over a string that contains a non-ASCII letter. Print each `rune` and its `unicode.IsLetter` result. Also print each byte index.

#### Advanced practical tasks

1. Replace a `regexp` search for a fixed word with `strings.EqualFold` or `strings.Contains`. Measure both on a large string with `testing.Benchmark` (see topic 13 for full benchmark rules). Write the two times.
2. Read `go doc strings.Builder`. Write a function that writes 10 000 integers into a `Builder` and a function that uses `+` in a loop. Compare `len` of the results. Do not claim a winner without a measurement.

---

## `time` — time, duration, timers, tickers

Package `time` represents a moment (`time.Time`) and a length (`time.Duration`). A `Duration` is an `int64` count of nanoseconds. The package defines `time.Second`, `time.Millisecond`, and the other units.

`time.Now()` returns the current local time. Use `t.UTC()` for UTC. Use `time.Date` to build a time from parts. Use `time.Parse` and `t.Format` with a layout. The layout is the reference time, not a `YYYY` pattern:

```text
Mon Jan 2 15:04:05 MST 2006
```

That date is the only layout language. Example: `"2006-01-02"` formats a year-month-day. `"15:04:05"` formats a 24-hour clock.

`time.Since(t)` is `time.Now().Sub(t)`. `time.Until(t)` is `t.Sub(time.Now())`.

`time.Sleep` blocks the current goroutine. `time.After` returns a channel that receives once. `time.NewTimer` returns a `*time.Timer` that you can stop. `time.NewTicker` returns a `*time.Ticker` that receives on a period.

Stop a ticker when you no longer need it:

```go
ticker := time.NewTicker(time.Second)
defer ticker.Stop()
```

If you ignore `Stop`, the ticker goroutine can leak. A timer that you do not stop can also leak until it fires. For cancellation of work, prefer `context.WithTimeout` (see the `context` section and topic 12).

Load a time zone with `time.LoadLocation`. The name `"Local"` is the system zone. The name `"UTC"` is UTC. A missing zone database is an error on some systems.

Do not compare times with `==` when the locations differ. Use `t.Equal`.

### Questions

#### Theoretical questions

1. What unit does `time.Duration` store?
2. Why is the layout string a real date in the year 2006?
3. What is the difference between `time.After` and `time.NewTicker`?
4. Why must you call `Stop` on a ticker that you create?
5. Why is `t.Equal` safer than `==` for `time.Time` values?

#### Easy practical tasks

1. Print `time.Now()` with `Format` and layout `"2006-01-02 15:04:05"`.
2. Sleep for `200 * time.Millisecond`. Print `time.Since` a start time.
3. Parse `"2024-01-15"` with layout `"2006-01-02"`. Print the `Time` in UTC.
4. Create a `Duration` of 90 seconds. Print it with `String` and with `Seconds`.

#### Medium practical tasks

1. Start a ticker of 100 milliseconds. Receive three times. Call `Stop`. Print the three times.
2. Parse a time that does not match the layout. Record the error. Fix the layout.
3. Load location `"America/New_York"` (or `"UTC"` if the zone file is missing). Convert `time.Now()` into that location.

#### Advanced practical tasks

1. Write a function that waits for a `time.Timer` or a done channel, whichever is first. Use `select`. Explain why this pattern belongs with `context` in a real API.
2. Compare `time.Now().Sub(t)` and `time.Since(t)` in a loop of 1000 calls. Then write why monotonic time matters for durations (read `go doc time`). Write four sentences.

---

## `math` / `math/rand/v2`

Package `math` provides floating-point functions and constants. Examples: `math.Sqrt`, `math.Pow`, `math.Abs`, `math.Hypot`, `math.Floor`, `math.Ceil`, `math.Mod`, `math.NaN`, `math.Inf`, `math.IsNaN`, `math.Pi`. `math.Max` and `math.Min` take `float64`. For integers, use the built-in `min` and `max` (Go 1.21 and later).

Do not use `math` for money that must be exact. Floating-point values have rounding error.

Package `math/rand/v2` is the random number package for Go 1.22 and later. Import `math/rand/v2`. The top-level functions are safe for concurrent use and use a hidden seed from the runtime:

```go
n := rand.IntN(100) // 0 <= n < 100
```

`rand.N` works for other integer types. `rand.Float64` returns a value in `[0.0, 1.0)`.

Create a local generator when you need a repeatable sequence:

```go
src := rand.NewPCG(1, 2)
r := rand.New(src)
```

`rand.NewChaCha8` is another source. Pass a seed that you control in tests. Do not use a fixed seed when each run must produce a different sequence.

The old package `math/rand` still exists. New code must use `math/rand/v2`. The old global source needed `Seed`. The v2 top-level API does not use that `Seed` function.

Do not use `math/rand/v2` for keys, tokens, or passwords. Use `crypto/rand` for those values.

### Questions

#### Theoretical questions

1. Why must you not use `float64` arithmetic for exact currency?
2. What range does `rand.IntN(100)` return?
3. Why does new code import `math/rand/v2` instead of `math/rand`?
4. When do you create `rand.New(rand.NewPCG(...))` instead of the top-level functions?
5. Which package must you use for a session token?

#### Easy practical tasks

1. Print `math.Sqrt(2)` and `math.Pi`.
2. Print `rand.IntN(6)` three times from `math/rand/v2`.
3. Use built-in `min` and `max` on two integers. Use `math.Min` on two `float64` values.
4. Run `go doc math/rand/v2`. Write the package comment in one sentence.

#### Medium practical tasks

1. Create two `rand.Rand` values with the same `PCG` seeds. Show that `IntN(10)` matches for the first five calls.
2. Detect `math.NaN()` with `math.IsNaN`. Then compute `math.Sqrt(-1)` and test that result.
3. Write a shuffle of a small slice with `rand.Shuffle` from `math/rand/v2`. Print before and after.

#### Advanced practical tasks

1. Compare `math/rand` (v1) and `math/rand/v2` in a short table: import path, seed, `Intn` versus `IntN`, concurrency. Write one migration of a v1 call to v2.
2. Read 16 bytes from `crypto/rand`. Encode them as hex. Write why this source is the correct source for a secret and why `math/rand/v2` is not.

---

## `os`, `os/exec`, `path/filepath`, `io/fs`

Package `os` talks to the operating system. Common file functions: `os.Open`, `os.Create`, `os.OpenFile`, `os.ReadFile`, `os.WriteFile`, `os.ReadDir`, `os.MkdirAll`, `os.Remove`, `os.CreateTemp`, `os.MkdirTemp`. Always check the error. Close a file that you open with `Open` or `Create`. `ReadFile` closes the file for you.

`os.Args` is the command-line argument slice. `os.Args[0]` is the program name. `os.Getenv` and `os.LookupEnv` read environment variables. `os.Exit` stops the process at once. `os.Exit` does not run deferred calls. Prefer a `return` from `main` when you can.

Package `os/exec` starts other programs. `exec.Command("go", "version")` builds a `*exec.Cmd`. `Output` runs the command and returns stdout. `CombinedOutput` returns stdout and stderr together. `Run` waits and returns an error. Use `exec.CommandContext` when you have a `context.Context` (topic 12).

Package `path/filepath` builds and splits operating-system paths. Use `filepath.Join` so that separators match the system. Use `filepath.Clean`, `filepath.Abs`, and `filepath.WalkDir`. Package `path` uses a slash and is for URL-like paths. Do not use `path` for Windows file paths.

Package `io/fs` defines file system interfaces. `fs.FS` is the core interface. `os.DirFS` turns a directory into an `fs.FS`. `embed.FS` also implements `fs.FS` (topic 9). `fs.WalkDir` walks any `fs.FS`. Prefer `fs.WalkDir` and `filepath.WalkDir` over the older `filepath.Walk`.

Sentinel errors such as `os.ErrNotExist` need `errors.Is` (topic 8).

### Questions

#### Theoretical questions

1. What is the difference between `os.ReadFile` and `os.Open` plus `io.ReadAll`?
2. Why must you not call `os.Exit` from a function that uses `defer` to unlock or to close?
3. When do you use `filepath.Join` instead of `path.Join`?
4. What does `os.DirFS` return?
5. How do you pass a deadline to a child process with `os/exec`?

#### Easy practical tasks

1. Write a file with `os.WriteFile` and read it back with `os.ReadFile`. Print the text.
2. Print `os.Args` and `os.Getenv("PATH")` (or `Path` on Windows).
3. Join `"tmp"` and `"out.txt"` with `filepath.Join`. Print the result.
4. Run `exec.Command("go", "env", "GOOS")`. Print `Output()`.

#### Medium practical tasks

1. Create a temp directory with `os.MkdirTemp`. Write a file in it. List the directory with `os.ReadDir`. Remove the tree.
2. Use `filepath.WalkDir` on a small folder. Print each path and whether it is a directory.
3. Start a command that does not exist. Record the error. Then use `errors.Is` or `exec.ErrNotFound` as the documentation describes.

#### Advanced practical tasks

1. Open a file with `os.OpenFile` and flags for append. Write two lines from two runs. Show the final contents.
2. Implement a function that reads a file from either `os.DirFS(".")` or an `embed.FS` through one `fs.FS` parameter. Prove both sources.

---

## `io` and `io/ioutil` (prefer `io` and `os`)

Package `io` defines the stream interfaces. The core interfaces are `io.Reader`, `io.Writer`, `io.Closer`, and `io.Seeker`. Many types implement them: files, buffers, HTTP bodies, and compressors.

Common helpers:

- `io.Copy` copies from a `Reader` to a `Writer`
- `io.ReadAll` reads until EOF
- `io.LimitReader` limits the number of bytes
- `io.MultiReader` and `io.MultiWriter` compose streams
- `io.TeeReader` writes a copy while you read
- `io.NopCloser` wraps a `Reader` with a close that does nothing
- `io.Discard` is a `Writer` that drops bytes
- `io.EOF` is the sentinel for end of input. Treat `EOF` as success when a read finishes. Do not treat every `EOF` as a hard failure.

Package `io/ioutil` is deprecated since Go 1.16. New code must not import `io/ioutil`. Use this map:

- `ioutil.ReadAll` → `io.ReadAll`
- `ioutil.ReadFile` → `os.ReadFile`
- `ioutil.WriteFile` → `os.WriteFile`
- `ioutil.ReadDir` → `os.ReadDir`
- `ioutil.NopCloser` → `io.NopCloser`
- `ioutil.Discard` → `io.Discard`
- `ioutil.TempFile` → `os.CreateTemp`
- `ioutil.TempDir` → `os.MkdirTemp`

`os.ReadDir` returns `[]os.DirEntry`. The old `ioutil.ReadDir` returned `[]fs.FileInfo`. The return types differ. Adjust the call site.

Prefer small functions that take `io.Reader` or `io.Writer`. Then tests can pass a `bytes.Buffer` instead of a real file.

### Questions

#### Theoretical questions

1. What method does `io.Reader` require?
2. When is `io.EOF` a normal result?
3. Why is `io/ioutil` wrong in new Go 1.22 code?
4. What is the replacement for `ioutil.ReadFile`?
5. Why do functions that accept `io.Reader` test more easily?

#### Easy practical tasks

1. Copy a string into a `bytes.Buffer` with `io.Copy` and `strings.NewReader`.
2. Read all bytes from that buffer with `io.ReadAll`. Print the string.
3. Write three lines to `io.Discard`. Confirm that the program prints nothing from that write.
4. Open `go doc io/ioutil`. Copy the deprecation sentence.

#### Medium practical tasks

1. Limit a reader to 5 bytes with `io.LimitReader`. Read all. Print the length.
2. Replace every `ioutil` call in a short sample that you write with `ioutil` on purpose. Switch the sample to `io` and `os`. Show both import blocks.
3. Use `io.MultiWriter` to write the same line to a file and to a `bytes.Buffer`.

#### Advanced practical tasks

1. Implement a function `func count(r io.Reader) (int, error)` that counts bytes with a 4 KiB buffer and `Read`. Do not use `ReadAll` for the count. Test it with a file and with a `strings.Reader`.
2. Wrap a reader with `io.TeeReader` so that a hash and a file write happen in one `io.Copy`. Print the hash and the file size.

---

## `encoding/json`, `encoding/xml`, `encoding/csv`

These packages encode structured data.

Package `encoding/json` converts Go values to JSON and back. `json.Marshal` returns `[]byte`. `json.Unmarshal` fills a pointer. Use exported fields. Set names with a struct tag:

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age,omitempty"`
}
```

`omitempty` drops an empty field from the output. Use `json.NewEncoder` and `json.NewDecoder` for streams. A decoder reads one value after another from an `io.Reader`. `Decoder.DisallowUnknownFields` rejects extra keys. When you unmarshal into `any`, numbers become `float64`. Use a struct when you know the shape.

Package `encoding/xml` is the same idea for XML. `xml.Marshal` and `xml.Unmarshal` use tags such as `xml:"name"` and `xml:",attr"`. XML is common for older APIs. Prefer JSON for new HTTP APIs unless the other side requires XML.

Package `encoding/csv` reads and writes comma-separated records. `csv.NewReader(r).Read()` returns one record (`[]string`) or an error. `ReadAll` reads every record. `csv.NewWriter(w)` writes records. Call `Flush` and then check `w.Error()`. Set `Comma` when the file uses a semicolon or a tab.

Check every encode and decode error. Bad input is normal. Do not panic on a decode error.

### Questions

#### Theoretical questions

1. Why must a field be exported for `json.Marshal` to see it?
2. What does `omitempty` do?
3. When do you use `json.NewDecoder` instead of `json.Unmarshal`?
4. What Go type does JSON use for a number when the target is `any`?
5. Why must you call `Flush` on a `csv.Writer`?

#### Easy practical tasks

1. Marshal a struct to JSON and print the string.
2. Unmarshal `{"name":"Ada"}` into a struct. Print the name.
3. Write two CSV rows with `csv.NewWriter` and `Flush`. Print the buffer.
4. Marshal a small struct to XML and print the string.

#### Medium practical tasks

1. Decode a JSON array from a `strings.Reader` with `json.NewDecoder`. Loop until `io.EOF`.
2. Use `DisallowUnknownFields` on an object that has an extra key. Record the error.
3. Read a CSV line with a wrong number of fields. Set `FieldsPerRecord` and record the error.

#### Advanced practical tasks

1. Write a struct with `json` and `xml` tags. Produce both encodings. Then decode both back and compare the values.
2. Stream a large JSON array without `ReadAll` of the full body. Use `Decoder.Token` or a decode loop. Write why this design uses less memory.

---

## `net/http` client basics (`http.Get`, `http.Client`)

Package `net/http` contains a client and a server. This section is the client only. Topic 16 covers servers.

`http.Get(url)` sends a GET request with `http.DefaultClient`. Check the error. Then check `resp.StatusCode`. Then read `resp.Body`. Then close the body:

```go
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
```

If `err` is nil, you must close `Body`. A non-2xx status is not an `error` from `Get`. You must test the status code.

`http.DefaultClient` has no timeout. A stuck server can block the goroutine for a long time. Create a client with a timeout:

```go
client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Get(url)
```

Use `http.NewRequestWithContext` and `client.Do` when you have a `context.Context`. Topic 12 shows cancellation. For this topic, set `Timeout` on the client.

`http.Post` and `http.PostForm` send a body. The same rules apply: check the error, close the body, check the status.

Do not leak connections. Close the body even when you do not read it. If you need only the status, read a small amount or use `io.Copy(io.Discard, resp.Body)` before close so that the transport can reuse the connection.

Use `https://example.com` or a local `httptest` server for practice. Do not send tests to a server that you do not control.

### Questions

#### Theoretical questions

1. Does `http.Get` return an error when the status is `404`?
2. Why must you close `resp.Body`?
3. What is the timeout of `http.DefaultClient`?
4. How do you send a request with a client timeout of 10 seconds?
5. Why is a context useful on a request? Give one sentence. Topic 12 gives the full answer.

#### Easy practical tasks

1. Call `http.Get("https://example.com")`. Print `StatusCode`. Close the body.
2. Read the body with `io.ReadAll`. Print `len(body)`.
3. Create `http.Client{Timeout: 5 * time.Second}`. Call `Get` on the same URL.
4. Run `go doc net/http.Client`. Write the purpose of the `Timeout` field.

#### Medium practical tasks

1. Handle a URL that does not resolve. Print the error. Then handle a URL that returns a non-2xx status (use `https://example.com/status/404` only if you have such a host, or use `httptest`).
2. Use `http.NewRequest` with method `HEAD` and `client.Do`. Print `StatusCode` and `ContentLength`.
3. Compare `http.Get` and a client with `Timeout: 1 * time.Millisecond` on a slow URL or a blocked address. Record the timeout error.

#### Advanced practical tasks

1. Build a small `httptest.Server` that sleeps 2 seconds. Call it with a client timeout of 200 milliseconds. Record the error. Then call it with a 5 second timeout and read the body.
2. Write a helper `func get(ctx context.Context, url string) ([]byte, error)` that uses `NewRequestWithContext`. Cancel the context in a test and record the error. Keep the deep context rules for topic 12.

---

## `log` and `log/slog`

Package `log` writes simple lines to a writer. The default logger writes to standard error. Common calls: `log.Print`, `log.Printf`, `log.Println`. `log.Fatal` logs and then calls `os.Exit(1)`. `log.Panic` logs and then panics. `log.SetFlags` changes the prefix (date, time, file). `log.SetOutput` changes the writer.

`log.Fatal` skips deferred cleanup because it calls `os.Exit`. Do not use `Fatal` inside a library. Return an `error`.

Package `log/slog` is the structured logger in Go 1.21 and later. A log record has a level, a message, and attributes. Common levels: `Debug`, `Info`, `Warn`, `Error`.

```go
slog.Info("listen", "port", 8080)
slog.Error("open", "path", path, "err", err)
```

You can use typed attributes: `slog.Int("port", 8080)`, `slog.String("path", path)`.

Create a logger with a handler:

```go
h := slog.NewJSONHandler(os.Stdout, nil)
logger := slog.New(h)
logger.Info("ready")
```

`slog.NewTextHandler` writes key=value text. `slog.SetDefault` changes the default logger that `slog.Info` uses. `logger.With("req_id", id)` returns a logger that always includes that attribute.

Use `log` for a tiny tool that needs one line. Use `slog` for services and for any program that another system must parse. Do not build JSON by hand with `fmt` when `slog` can do it.

`slog` does not replace `error` return values. Log at the boundary. Return the error to the caller when the caller must act.

### Questions

#### Theoretical questions

1. Where does the default `log` logger write?
2. Why is `log.Fatal` dangerous together with `defer`?
3. What extra data does `slog` store that `log.Println` does not store in a structured way?
4. What is the difference between `slog.NewJSONHandler` and `slog.NewTextHandler`?
5. Why must a library function return an `error` instead of calling `log.Fatal`?

#### Easy practical tasks

1. Print a line with `log.Println`. Then set a prefix with `log.SetPrefix`.
2. Call `slog.Info` with one message and two key-value pairs.
3. Create a JSON handler on a `bytes.Buffer`. Log one `Info` event. Print the buffer.
4. Run `go doc log/slog`. Write one sentence about the package.

#### Medium practical tasks

1. Use `slog.SetDefault` with a text handler that includes source (`HandlerOptions{AddSource: true}`). Log once. Show the source key in the output.
2. Create a logger with `With("component", "store")`. Log from two functions. Show that the attribute appears on both lines.
3. Replace three `log.Printf` calls with `slog` attributes. Keep the same facts in the log.

#### Advanced practical tasks

1. Write an `slog.Handler` wrapper that drops `Debug` events (or use `HandlerOptions{Level: slog.LevelInfo}`). Prove that `Info` appears and `Debug` does not.
2. Log an `error` with `slog.Error` and then still return the error to the caller. Write a four-sentence rule for when a layer logs and when a layer only returns.

---

## `context` package overview

Package `context` carries a deadline, a cancel signal, and optional request values across API boundaries. This section is an introduction. Topic 12 is the deep dive.

The root contexts are `context.Background()` and `context.TODO()`. Use `Background` at the start of `main` or in a test when no other context exists. Use `TODO` only when you have not yet passed a real context and you will fix the call.

Child contexts add control:

- `WithCancel` returns a context and a `cancel` function
- `WithTimeout` cancels after a `time.Duration`
- `WithDeadline` cancels at a `time.Time`
- `WithValue` stores a key-value pair

Always call `cancel` when you create a cancelable context. Use `defer cancel()`. The call is safe after the context is already done.

The first parameter of a function that can block must be `ctx context.Context`. Do not store a context in a struct for a long lifetime unless the struct represents that one operation.

Listen for the end of work:

```go
select {
case <-ctx.Done():
    return ctx.Err()
case res := <-out:
    return res, nil
}
```

`ctx.Err()` is `context.Canceled` or `context.DeadlineExceeded` after the context is done.

Do not use `WithValue` for optional parameters that belong in the function signature. Topic 12 explains the risks. For this topic, know that the value API exists and that you must keep keys private.

`http.NewRequestWithContext` attaches a context to a client request. `exec.CommandContext` stops a child process when the context is done.

### Questions

#### Theoretical questions

1. What three kinds of data can a `context.Context` carry?
2. What is the difference between `Background` and `TODO`?
3. Why must you call the `cancel` function from `WithTimeout`?
4. Where does `ctx` appear in a function signature?
5. Which two errors can `ctx.Err()` return after `Done` is closed?

#### Easy practical tasks

1. Create a context with `WithTimeout` of 50 milliseconds. Sleep 80 milliseconds. Print `ctx.Err()`.
2. Create a context with `WithCancel`. Call `cancel`. Print `ctx.Err()`.
3. Write a function signature that takes `ctx context.Context` as the first parameter and a `url string` as the second.
4. Run `go doc context`. List the four `With*` functions.

#### Medium practical tasks

1. Start a `select` that waits on `ctx.Done()` and on `time.After(time.Second)`. Use a 10 millisecond timeout. Print which case runs.
2. Pass `context.Background()` into `http.NewRequestWithContext` for `https://example.com`. Call `Do` on a client with a timeout. Close the body.
3. Store a request id with `WithValue` and a private key type. Read it back. Write one sentence about why a string key is a bad idea (topic 12 gives more).

#### Advanced practical tasks

1. Write a loop that reads from a channel and also stops on `ctx.Done()`. Show a clean return when you cancel.
2. List three APIs in the standard library that take `context.Context`. Write this list for topic 12. Do not implement a full cancel tree in this task.

---

## `flag` / `os.Args` for CLI

A command-line program receives arguments in `os.Args`. `os.Args[0]` is the program name or the path of the binary. `os.Args[1:]` are the arguments that the user types.

Package `flag` parses flags such as `-n` and `--name`. Define flags before you parse:

```go
n := flag.Int("n", 1, "count")
name := flag.String("name", "", "user name")
flag.Parse()
```

`flag.Int` returns a pointer. Read `*n` after `Parse`. `flag.Parse` exits the process with a help message when the user passes `-h` or a bad flag. The default `flag.CommandLine` set writes that message to standard error.

`flag.Args()` returns the remaining non-flag arguments. `flag.NArg()` is the count of those arguments. `flag.NFlag()` is the count of flags that the user set.

Use `flag` when you need typed flags and help text. Use `os.Args` when you have only positional arguments and no flags. Do not parse `os.Args` by hand when `flag` already covers the syntax.

Define flags in `main` or in an `init` that belongs to the command package. Do not register global flags inside a library package. A library must accept values as parameters.

`flag.Duration` parses values such as `250ms`. That type matches `time.Duration`.

### Questions

#### Theoretical questions

1. What is `os.Args[0]`?
2. Why does `flag.Int` return a pointer?
3. What does `flag.Args` return after `Parse`?
4. Why must a library not call `flag.String` at package level?
5. When is `os.Args` enough and `flag` not required?

#### Easy practical tasks

1. Print every element of `os.Args` with its index.
2. Add a `-name` string flag with a default of `world`. Print `hello` and the name.
3. Add a `-n` int flag. Loop `n` times and print the index.
4. Run the program with `-h`. Copy the help text that `flag` prints.

#### Medium practical tasks

1. Require one positional file name after flags. Use `flag.Args()`. Return an error from `main` logic when the name is missing (print and `os.Exit(2)` or return a code if you use a `run` function).
2. Add `flag.Duration("timeout", time.Second, "wait")`. Print the duration.
3. Use `flag.Var` or a custom `flag.Value` for a comma-separated list. Parse `-tags=a,b`. Print the slice.

#### Advanced practical tasks

1. Split `main` into `func run(args []string, out io.Writer) error`. Parse flags with a new `flag.FlagSet` so that tests can call `run` without `os.Exit`. Write two tests: good flags and bad flags.
2. Compare `flag` with a third-party CLI package in six short sentences. State when you stay with `flag`.

---

## `testing` package (overview; see topic 13)

Package `testing` is the test framework that `go test` runs. This section is an overview. Topic 13 covers table tests, benchmarks, examples, fuzzing, and coverage.

A test file name ends with `_test.go`. A test function has this shape:

```go
func TestSum(t *testing.T) {
    if sum(2, 3) != 5 {
        t.Fatalf("got %d", sum(2, 3))
    }
}
```

The name must start with `Test` and then a capital letter or a digit. The parameter is `*testing.T`.

`t.Error` and `t.Errorf` fail the test and continue. `t.Fatal` and `t.Fatalf` fail the test and stop that function. Use `Fatal` when the next lines cannot run. Use `t.Helper()` in a helper so that the reported line is the call site.

`t.Run` starts a subtest. Topic 13 shows table-driven tests that use `t.Run`.

Run tests:

```text
go test ./...
go test -v
```

Other entry points exist: `Benchmark` functions, `Example` functions, and `Fuzz` functions. You will write those in topic 13. For this topic, write one `Test` function per package that you care about.

Tests live in the same package (`package store`) or in an external test package (`package store_test`). The external form can import only the exported API.

Do not put `t.Fatal` in production code. The `testing` package is for tests.

### Questions

#### Theoretical questions

1. What file name suffix does `go test` use?
2. What is the difference between `t.Error` and `t.Fatal`?
3. What is the signature of a test function?
4. Why does topic 13 exist if this section already runs `go test`?
5. What extra power does `package foo_test` give you?

#### Easy practical tasks

1. Write `func Add(a, b int) int` and `TestAdd` that checks `Add(2, 2) == 4`. Run `go test`.
2. Make the test fail on purpose. Read the output. Fix the test.
3. Run `go test -v`. Write what `-v` adds.
4. Run `go doc testing.T`. List `Error`, `Fatal`, and `Run`.

#### Medium practical tasks

1. Add a helper `assertEq` that calls `t.Helper` and `t.Fatalf`. Use it twice in one test.
2. Add `t.Run` for two cases: `2+3` and `0+0`. Run `go test -v`.
3. Move the test to `package xxx_test` and import your module path. Keep the test passing.

#### Advanced practical tasks

1. Write a failing test, a passing test, and a test that calls `t.Skip`. Run `go test -v` and explain the three results in four sentences.
2. Open topic 13 in `TOPICS.md`. Make a list of test features that this section did not cover. Keep that list for later study.

---

## `sync` and `sync/atomic` (overview; see topic 11)

Package `sync` provides locks and other primitives for goroutines that share memory. This section is an overview. Topic 11 covers goroutines, channels, races, and the full set of patterns.

`sync.Mutex` has `Lock` and `Unlock`. Unlock in a `defer` after a successful lock. Do not copy a `Mutex` after you use it. Put the mutex above the data that it protects in the same struct.

`sync.RWMutex` allows many readers or one writer. Use it only when you measure a need. A `Mutex` is simpler.

`sync.WaitGroup` waits for a set of goroutines. Call `Add` before you start the goroutine. Call `Done` at the end of the goroutine. Call `Wait` in the parent. Prefer `Add(1)` next to each `go` statement. Do not use a negative `Add` by mistake.

Other types that topic 11 explains: `sync.Once`, `sync.Map`, `sync.Pool`, `sync.Cond`. Know the names. Do not use `sync.Map` as the default map. Use a normal map plus a mutex unless the documentation of `sync.Map` matches your case.

Package `sync/atomic` provides atomic operations on integers and pointers. Go 1.19 added types such as `atomic.Int64` and `atomic.Bool`:

```go
var n atomic.Int64
n.Add(1)
v := n.Load()
```

Atomics are for simple counters and flags. They do not make a full algorithm safe by themselves. If you need two fields to change together, use a `Mutex`.

The race detector is `go test -race`. Topic 11 shows how to read the report. A data race is a bug.

Prefer a channel when goroutines must hand off a value. Use a mutex when goroutines must update a struct in place. Topic 11 develops that choice.

### Questions

#### Theoretical questions

1. What two methods does `sync.Mutex` provide?
2. Why do you call `WaitGroup.Add` before `go`?
3. When is `atomic.Int64` enough and a `Mutex` not required?
4. Why must you not copy a struct that contains a `Mutex` after first use?
5. Which topic explains worker pools and `select`?

#### Easy practical tasks

1. Protect an `int` with a `Mutex`. Increment it once. Print the value.
2. Start two goroutines that each call `Done` on a `WaitGroup`. `Wait` in `main`.
3. Increment an `atomic.Int64` two times. Print `Load()`.
4. Run `go doc sync.WaitGroup`. Write the order of `Add`, `Done`, and `Wait`.

#### Medium practical tasks

1. Increment a counter from 10 goroutines, 100 times each, with a mutex. Print `1000`. Then remove the mutex and run `go test -race` (or `go run -race`). Record the race report if it appears.
2. Use `sync.Once` to initialize a value from two goroutines. Print the value once.
3. Compare a mutex around a `map[string]int` with an `atomic.Int64` counter. Write which tool fits each job.

#### Advanced practical tasks

1. Write a tiny map store with a `RWMutex`. Do one write and two reads from different goroutines. Then list what topic 11 must still teach (channels, leak, deadlock).
2. Read `go doc sync/atomic.Int64`. Use `CompareAndSwap` to set a flag from `0` to `1` only once. Prove that a second swap fails.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A program reads a flag, reads a file, decodes JSON, and logs a failure. Which standard library packages do you use for each step, and which deprecated package do you avoid?
2. How do `io.Reader`, `os.Open`, `embed.FS`, and `http.Response.Body` share one design idea?
3. When do you pick `fmt`, `strconv`, `strings`, and `encoding/json` for text?
4. What is the difference between a client `http.Client.Timeout`, a `time.Ticker`, and `context.WithTimeout`?
5. Which parts of this topic are only overviews, and which later topics complete them?

#### Easy practical tasks

1. Write a CLI that parses `-n` with `flag`, converts a string argument with `strconv.Atoi`, and prints the sum with `fmt`. Add one test.
2. Make a one-page cheat sheet with one line per package group in this topic. Include `math/rand/v2` and `log/slog`. Mark `io/ioutil` as deprecated.
3. Draw a flow: `os.Open` → `io.Copy` → `os.Stdout`. Label `error` checks and `Close`.
4. Run `go doc` on `fmt`, `strconv`, `time`, and `slog`. Write one sentence per package from the package comment.

#### Medium practical tasks

1. Write a program that loads a JSON file from disk (or from `embed`), prints one field, and logs with `slog` if the file is missing. Use `errors.Is` with `os.ErrNotExist`.
2. Call an `httptest` server with a client timeout. Decode a JSON body. Write the result as one CSV row.
3. Start a ticker and a `WaitGroup` with two goroutines. Stop the ticker. Wait. Run `go test -race` on a test that starts this program logic.

#### Advanced practical tasks

1. Build a small tool that uses `flag`, `os`, `encoding/json`, `log/slog`, `context.WithTimeout`, and `net/http` together. The tool must fetch one URL, decode JSON, and exit with a non-zero code on failure. Use `httptest` in a test. Do not use `ioutil`. Do not use `math/rand` v1.
2. Write a migration note of one page for a Go 1.15 module that still imports `io/ioutil` and `math/rand`. List every replacement that this topic requires for Go 1.22 (`os`/`io`, `rand/v2`, `slog`). Apply the note to a tiny sample module that you create.
