# 9. Standard Library

## Description

The Go standard library is a set of packages that ship with the toolchain. You import these packages without `go get`. This topic shows the packages that most programs need: text, time, files, JSON, HTTP clients, logs, and flags.

Complete topics 1 to 8 before this topic. Prefer the standard library when it solves the problem. This topic uses Go 1.22 or later. Use `log/slog` for structured logs. Do not use `io/ioutil` in new code. Topic 10 covers concurrency and context in depth. Topic 11 covers tests in depth.

Use one term for each package. Read package documentation on pkg.go.dev or with `go doc`.

---

## `fmt`, `strings`, and `strconv`

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
- `%w` wrap an error in `fmt.Errorf` only

Scan functions return the number of items and an `error`. Prefer `bufio.Scanner` or `strconv` when you parse a single line or a single number.

Package `strings` works with text. Common functions: `Contains`, `HasPrefix`, `HasSuffix`, `Split`, `Join`, `ReplaceAll`, `TrimSpace`, `ToLower`, `ToUpper`, `Cut`. Type `strings.Builder` builds a string with few allocations. Type `strings.Reader` implements `io.Reader` over a string.

Package `strconv` converts between basic types and their string forms. Use this package when the input or the output is a single value.

```go
n, err := strconv.Atoi("17")
s := strconv.Itoa(17)
f, err := strconv.ParseFloat("3.14", 64)
```

Every `Parse*` function returns an `error`. Check the error. `Atoi` is a short form of `ParseInt` with base 10 and bit size 0. `FormatFloat` needs a format byte, a precision, and a bit size.

Prefer `strings` when you search for a fixed substring. Prefer `strconv` when you convert one number or one bool. Use `fmt` for humans and for simple logs.

### Questions

#### Theoretical questions

1. What is the difference between `fmt.Println` and `fmt.Printf`?
2. What does the verb `%T` print?
3. What problem does `strings.Builder` solve?
4. What is the difference between `strconv.Atoi` and `strconv.ParseInt`?
5. When must you use `strconv` instead of `fmt` for a conversion?

#### Easy practical tasks

1. Print your name and an integer with `fmt.Printf` and verbs `%s` and `%d`.
2. Use `strings.ToUpper` and `strings.Contains` on one sentence. Print both results.
3. Convert `"17"` to `int` with `Atoi`. Print the value and the error.
4. Convert `17` to a string with `Itoa`. Print the string.

#### Medium practical tasks

1. Scan two integers from a string with `fmt.Sscanf`. Check the count and the error for a good string and for a bad string.
2. Build a line with `strings.Builder`. Write three parts. Print the string.
3. Parse `"ff"` as a hex `int64` with `ParseInt`. Then format that value as hex with `FormatInt`.

#### Advanced practical tasks

1. Write a function that parses a list of integers from a slice of strings and returns one `error` that names the first bad index. Use `strconv`, not `fmt.Sscanf`.
2. Compare `fmt.Sprint(3.14)` with a `strconv` format of the same number. Write which package you pick for a stable machine-readable number.

---

## `time`

Package `time` measures time and clocks. Type `time.Time` is a point in time. Type `time.Duration` is a span of time. A duration is an `int64` number of nanoseconds.

```go
now := time.Now()
later := now.Add(2 * time.Second)
d := later.Sub(now)
```

Common duration constants are `time.Nanosecond`, `time.Microsecond`, `time.Millisecond`, `time.Second`, `time.Minute`, and `time.Hour`.

`time.Sleep` pauses the current goroutine. `time.After` returns a channel that receives the current time after a duration. `time.NewTimer` and `time.NewTicker` give timers that you must stop when you no longer need them. Stop a ticker with `ticker.Stop()`.

Parse and format times with a layout. The reference time is `Mon Jan 2 15:04:05 MST 2006`. That date is the layout key.

```go
t, err := time.Parse(time.RFC3339, "2024-01-02T15:04:05Z")
s := t.Format("2006-01-02")
```

Use `time.RFC3339` for machine exchange. Use `time.UTC` or a named location for storage. `time.Now()` uses the local zone. Convert with `t.UTC()` when you store a timestamp.

Compare times with `t.Before`, `t.After`, and `t.Equal`. Do not compare `Time` values with `==` when location can differ. `Equal` compares the instant.

A zero `Time` is year 1. Test with `t.IsZero()`. Do not treat a zero time as "now".

### Questions

#### Theoretical questions

1. What is the difference between `time.Time` and `time.Duration`?
2. What date is the layout reference for `Format` and `Parse`?
3. Why must you stop a `Ticker`?
4. Why is `t.Equal` safer than `==` for two `Time` values?
5. What does `t.IsZero()` mean?

#### Easy practical tasks

1. Print `time.Now()` with `Format(time.RFC3339)`.
2. Sleep for `200 * time.Millisecond`. Print `start` and `end`.
3. Parse `2024-06-01` with layout `2006-01-02`. Print the year and the month.
4. Add `48 * time.Hour` to a parsed date. Print the new date.

#### Medium practical tasks

1. Time a loop with `time.Since(start)`. Print the duration.
2. Use `time.After` in a `select` with a default or a second case. Show the timeout path.
3. Parse a time in `UTC` and in a named location. Print both. Show `Equal` and `==`.

#### Advanced practical tasks

1. Write a function that parses several layouts until one succeeds. Return an error when all fail.
2. Build a ticker that prints three ticks and then stops. Confirm that the program exits.

---

## `os`, `path/filepath`, and `io`

Package `os` talks to the operating system. Use `os.Open` to read a file. Use `os.Create` or `os.OpenFile` to write a file. Check the error. Close the file with `defer f.Close()` after a successful open.

```go
f, err := os.Open("data.txt")
if err != nil {
	return err
}
defer f.Close()
```

`os.ReadFile` and `os.WriteFile` cover small files. `os.Args` holds command-line arguments. `os.Getenv` and `os.LookupEnv` read environment variables. `os.Exit` stops the process and skips deferred calls.

`os.Stdin`, `os.Stdout`, and `os.Stderr` are open files. They implement `io.Reader` or `io.Writer`.

Package `path/filepath` builds and splits paths for the current operating system. Use `filepath.Join` to join elements. Use `filepath.Base`, `filepath.Dir`, and `filepath.Ext` to split a path. Use `filepath.Abs` for an absolute path. Use `filepath.WalkDir` to visit a tree.

Do not join paths with `+` and `/`. Windows uses `\`. `filepath.Join` picks the correct separator. Package `path` is for slash paths such as URLs. Use `path/filepath` for disk paths.

Package `io` defines `Reader`, `Writer`, `Closer`, and helpers. `io.Copy` copies from a reader to a writer. `io.ReadAll` reads until EOF. `io.NopCloser` wraps a reader with a no-op `Close`. `io.EOF` marks the end of input. `io.MultiWriter` writes to several writers.

Do not use `io/ioutil`. That package is deprecated. The functions moved to `io` and `os`.

### Questions

#### Theoretical questions

1. When do you call `f.Close` after `os.Open`?
2. What is the difference between `path` and `path/filepath`?
3. What does `io.Copy` do?
4. Why must you not join disk paths with string `+` and `/`?
5. Where did `ioutil.ReadFile` move?

#### Easy practical tasks

1. Write a file with `os.WriteFile`. Read it back with `os.ReadFile`. Print the text.
2. Print `os.Args[0]` and `len(os.Args)`.
3. Join two path parts with `filepath.Join`. Print the result.
4. Copy a `strings.NewReader` to `os.Stdout` with `io.Copy`.

#### Medium practical tasks

1. Open a file, `defer` close, and copy the file to a `bytes.Buffer` with `io.Copy`.
2. Walk a small directory with `filepath.WalkDir`. Print each file name.
3. Read `os.LookupEnv` for a name that exists and a name that does not. Print both results.

#### Advanced practical tasks

1. Write a function that copies one file to another with `io.Copy` and correct close of both files on every path.
2. Implement a small `cat` command: read each argument as a path, copy to stdout, wrap errors with the path.

---

## `encoding/json`

Package `encoding/json` encodes and decodes JSON. `json.Marshal` turns a Go value into `[]byte`. `json.Unmarshal` fills a Go value from `[]byte`.

```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

b, err := json.Marshal(User{Name: "Ada", Age: 36})
var u User
err = json.Unmarshal(b, &u)
```

Export the fields that you want in JSON. Use struct tags for names. `omitempty` omits a zero value. `json:"-"` skips a field. Topic 5 covers tags.

`json.NewEncoder` writes to an `io.Writer`. `json.NewDecoder` reads from an `io.Reader`. Use these types for streams and for HTTP bodies.

```go
err := json.NewEncoder(w).Encode(u)
err = json.NewDecoder(r).Decode(&u)
```

Numbers in JSON objects often decode into `float64` when the target is `map[string]any` or `any`. Use a struct with `int` or `json.Number` when you need exact integers.

Unknown object keys are ignored by default when you decode into a struct. `Decoder.DisallowUnknownFields` rejects extra keys.

A `nil` slice or map encodes as `null`. An empty slice encodes as `[]`. An empty map encodes as `{}`.

Always check the error from marshal and unmarshal. Do not ignore a syntax error.

### Questions

#### Theoretical questions

1. What does `json.Marshal` return?
2. Why must you pass a pointer to `Unmarshal`?
3. When do you use `Encoder` and `Decoder` instead of `Marshal`?
4. How does JSON encode a `nil` slice?
5. What type do JSON numbers become in `map[string]any`?

#### Easy practical tasks

1. Marshal a struct with two exported fields. Print the JSON string.
2. Unmarshal that JSON back into a struct. Print the fields.
3. Add `json:"-"` on one field. Confirm that the field is absent in JSON.
4. Marshal a `nil` slice and an empty slice. Print both.

#### Medium practical tasks

1. Decode a JSON array of objects into `[]User`. Print the length.
2. Use `json.NewDecoder` on a `strings.NewReader`. Decode two values in sequence if the input has two objects.
3. Decode into `map[string]any`. Assert one string and one number. Print the types.

#### Advanced practical tasks

1. Enable `DisallowUnknownFields` and feed an extra key. Record the error. Then decode without that option.
2. Write a type with custom `MarshalJSON` and `UnmarshalJSON`. Prove both directions with a test.

---

## `net/http` client

Package `net/http` includes a client and a server. This section covers the client. Topic 13 covers servers.

`http.Get` and `http.Post` are short helpers. They use `http.DefaultClient`. Always close the response body.

```go
resp, err := http.Get("https://example.com")
if err != nil {
	return err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
```

Check `resp.StatusCode`. A 404 response is not an `error` from `Get`. The error is for network and protocol failures.

Use `http.NewRequestWithContext` when you need a method, headers, a body, or a timeout. Pass a context. Topic 10 covers context.

```go
req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
req.Header.Set("Accept", "application/json")
resp, err := http.DefaultClient.Do(req)
```

`http.Client` holds a timeout and a transport. Set `Timeout` on the client for a simple deadline. Reuse one `Client`. Do not create a client per request. The default transport pools connections.

```go
client := &http.Client{Timeout: 5 * time.Second}
```

`http.NewRequest` without context is older. Prefer `NewRequestWithContext`. Do not ignore redirects when you post credentials. Configure `CheckRedirect` when you must.

For tests, use `httptest.Server`. Topic 11 covers `httptest`.

### Questions

#### Theoretical questions

1. Does `http.Get` return an error for HTTP status 404?
2. Why must you close `resp.Body`?
3. Why do you reuse one `http.Client`?
4. What is the difference between `http.Get` and `Client.Do`?
5. What does `Client.Timeout` limit?

#### Easy practical tasks

1. Call `http.Get` on `https://example.com`. Print `StatusCode` and `len(body)`.
2. Print one response header, such as `Content-Type`.
3. Send a request to a URL that does not resolve. Print the error.
4. Set `Accept` on a request with `NewRequestWithContext` and `Background` context. Print the status.

#### Medium practical tasks

1. Create `&http.Client{Timeout: time.Millisecond}`. Call a slow URL or a local delay. Record the timeout error.
2. POST a JSON body with `bytes.NewReader` and `http.NewRequestWithContext`. Print the status.
3. Follow or block redirects with `CheckRedirect`. Record the status sequence for a URL that redirects.

#### Advanced practical tasks

1. Write a helper `func fetch(ctx context.Context, url string) ([]byte, error)` that sets a timeout on the context, closes the body, and wraps errors with the URL.
2. Compare `DefaultClient` with a custom `Transport` that sets `MaxIdleConns`. Write when you change the transport.

---

## `log/slog` and `flag`

Package `log/slog` writes structured logs. A record has a level, a message, and attributes. The default handler writes text. `slog.NewJSONHandler` writes JSON.

```go
slog.Info("started", "port", 8080)
slog.Error("open failed", "err", err)
```

Levels include `Debug`, `Info`, `Warn`, and `Error`. Use `slog.SetDefault` to install a logger. Pass a `*slog.Logger` when a package must not use the global logger.

```go
h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
logger := slog.New(h)
logger.Info("ready", "env", "dev")
```

Put keys in a stable form. Use `err` for an error attribute. Do not log secrets. Do not use `log.Print` for new services. The older `log` package is fine for tiny tools.

Package `flag` parses command-line flags. Define flags before `flag.Parse`.

```go
port := flag.Int("port", 8080, "listen port")
name := flag.String("name", "app", "service name")
flag.Parse()
fmt.Println(*port, *name)
```

`flag.Int` returns `*int`. Dereference the pointer after `Parse`. `flag.Args` returns remaining arguments. `flag.Usage` prints help.

Do not mix `flag` with ad-hoc `os.Args` loops for the same program. Parse once in `main`. Pass values into the rest of the program.

For subcommands, you can use `flag.NewFlagSet`. Many teams later use a small extra module. Start with `flag`.

### Questions

#### Theoretical questions

1. What three parts does an `slog` record have?
2. Why do you prefer `slog` over `log.Print` in a service?
3. When do you call `flag.Parse`?
4. What type does `flag.Int` return?
5. What must you never put in a log attribute?

#### Easy practical tasks

1. Log one `Info` line with two attributes using `slog`.
2. Create a JSON handler on `os.Stdout`. Log one `Error` with an `err` attribute.
3. Define `-n` as an `int` flag with default `1`. Parse and print `*n`.
4. Run the program with `-h` or `-help`. Read the usage text.

#### Medium practical tasks

1. Add `-verbose`. When it is true, set the handler level to `Debug`. Show a debug line that is hidden by default.
2. Parse `-o` as an output path. Open that file and pass it to `slog.NewJSONHandler`.
3. Use `flag.Args` after flags. Require one positional path. Return an error when it is missing.

#### Advanced practical tasks

1. Write `main` that parses flags, builds an `slog.Logger`, and passes both into a `run` function. Keep `main` small.
2. Compare `flag` with a config file plus environment variables. Write six sentences on when each source wins.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `fmt`, `strings`, and `strconv` split the job of text work?
2. When do you choose `os.ReadFile`, `io.Copy`, and `encoding/json` for the same bytes on disk?
3. How do `time.Duration` and `http.Client.Timeout` work together on a client call?
4. Why does JSON need exported fields and tags while `slog` needs stable attribute keys?
5. What does `main` own when you combine `flag`, `slog`, and an HTTP GET?

#### Easy practical tasks

1. Write a program that parses `-name`, logs `hello` with `slog`, and prints the name with `fmt`.
2. Read a JSON file with `os.ReadFile` and `json.Unmarshal`. Print one field.
3. Format `time.Now()` as RFC3339 and write it to a file with `os.WriteFile`.
4. Make a cheat sheet: one function from each package in this topic.

#### Medium practical tasks

1. Build a tiny CLI: `-url` and `-timeout`. GET the URL with a client timeout. Log status and duration with `slog`.
2. Walk a directory, skip non-`.json` files, and count objects that decode into a struct.
3. Copy standard input to a file with `io.Copy` and log errors with `slog.Error`.

#### Advanced practical tasks

1. Write a client that fetches JSON, decodes into a struct, and retries once on a 5xx status. Bound the wait with `time` and context from topic 10 if you already know it.
2. Replace `io/ioutil` usage in an old snippet with `os` and `io`. Document each rename.
