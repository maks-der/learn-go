# 16. I/O, Files, and Networking

## Description

This topic shows how Go programs read bytes, write bytes, and talk over the network. The `io.Reader` and `io.Writer` interfaces are the center of that design. Files, buffers, TCP, HTTP, and TLS all implement or consume those interfaces.

Use one term for each concept. A reader produces bytes. A writer consumes bytes. A handler serves an HTTP request. A pattern in `ServeMux` (Go 1.22 and later) matches a method and a path. Do not ignore errors from I/O calls. Do not hold a file or a connection open after you finish the work. Close the resource, usually with `defer`.

---

## `io.Reader` / `io.Writer` composition

`io.Reader` has one method:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

`Read` fills `p` with up to `len(p)` bytes. It returns the number of bytes and an error. `io.EOF` means there are no more bytes. A `Read` can return `n > 0` and `io.EOF` in the same call. Process the bytes first. Then handle `EOF`.

`io.Writer` has one method:

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

`Write` writes `len(p)` bytes or returns an error. Check that `n` equals `len(p)` when you need a full write. `io.ErrShortWrite` reports a short write without a more specific error.

You compose readers and writers. You wrap one value with another value. The wrapper also implements `io.Reader` or `io.Writer`.

Useful helpers in package `io`:

- `io.Copy(dst, src)` copies from a reader to a writer until EOF
- `io.CopyN(dst, src, n)` copies at most `n` bytes
- `io.LimitReader(r, n)` stops after `n` bytes
- `io.MultiReader(r1, r2, ...)` reads each reader in order
- `io.MultiWriter(w1, w2, ...)` writes the same bytes to each writer
- `io.TeeReader(r, w)` reads from `r` and writes a copy to `w`
- `io.NopCloser(r)` adds a no-op `Close` to a reader
- `io.ReadAll(r)` reads all remaining bytes (watch memory)

```go
limited := io.LimitReader(file, 1024)
n, err := io.Copy(os.Stdout, limited)
```

`io.ReadCloser` and `io.WriteCloser` add `Close()`. HTTP bodies and files are closers. Close them.

Prefer `io.Reader` and `io.Writer` in function parameters. Accept the interface. Return a concrete type when you own the implementation. This matches the guideline "accept interfaces, return structs".

Do not ignore the byte count. Do not treat every error as fatal before you use `n`. Do not use `io.ReadAll` on an unbounded network stream. Use `io.LimitReader` or a decoder that stops.

`io.Pipe` creates a connected pair of `PipeReader` and `PipeWriter`. A write blocks until a read consumes the data (with internal buffering details). Use a pipe when one goroutine produces bytes and another consumes them.

### Questions

#### Theoretical questions

1. What does `Read` return when the stream ends?
2. Why must you process `n` bytes before you handle `io.EOF`?
3. What does `io.LimitReader` prevent?
4. What is the difference between `io.MultiWriter` and two separate `Write` calls in your code?
5. Why do HTTP bodies implement `io.ReadCloser` and not only `io.Reader`?
6. When is `io.ReadAll` the wrong tool?

#### Easy practical tasks

1. Use `io.Copy` to copy from a `strings.NewReader` to `os.Stdout`.
2. Wrap a reader with `io.LimitReader` of 5 bytes. Copy to stdout. Use a longer string.
3. Use `io.MultiWriter` to write one line to a `bytes.Buffer` and to stdout.
4. Call `Read` in a loop on a small reader. Print `n` and `err` each time until `EOF`.

#### Medium practical tasks

1. Write a function `CopyLimit(dst io.Writer, src io.Reader, n int64) (int64, error)` that uses `io.Copy` and `io.LimitReader`. Test it.
2. Use `io.TeeReader` to copy a file to stdout and to a hash (`crypto/sha256`). Print the hex checksum.
3. Build a `io.MultiReader` from three `strings.NewReader` values. Read all bytes with a small buffer (4 bytes). Print each chunk.

#### Advanced practical tasks

1. Implement `io.Reader` on a struct that yields Fibonacci numbers as decimal lines. Compose it with `io.LimitReader`. Write a test that reads a fixed byte count.
2. Use `io.Pipe` with two goroutines: one writes JSON lines, one decodes them. Close the pipe. Handle errors from both sides.

---

## `bufio` for efficient I/O

A raw `Read` or `Write` on a file or a socket can do one system call per call. Many small calls are slow. Package `bufio` adds a memory buffer in the process.

`bufio.Reader` wraps an `io.Reader`. It fills an internal buffer from the inner reader. `Read`, `ReadByte`, `ReadString`, and `Peek` use that buffer.

`bufio.Writer` wraps an `io.Writer`. `Write` stores bytes in a buffer. `Flush` writes the buffer to the inner writer. You must call `Flush` before you need the data on the other side. `defer writer.Flush()` is common. Still check the error from `Flush`.

`bufio.Scanner` reads tokens. The default split is by lines. Scanner is the simple tool for line-oriented text.

```go
f, err := os.Open("data.txt")
if err != nil {
    return err
}
defer f.Close()

sc := bufio.NewScanner(f)
for sc.Scan() {
    line := sc.Text()
    _ = line
}
if err := sc.Err(); err != nil {
    return err
}
```

Scanner has a default maximum token size (64 KiB). A longer line fails the scan. Raise the limit with `sc.Buffer(buf, max)` when you control the input. Do not raise the limit for untrusted input without a bound that you accept.

`Peek` returns the next bytes without removing them from the buffer. Use `Peek` when you need to look at a prefix (for example, a magic number).

`bufio.NewReadWriter` holds a reader and a writer for a bidirectional stream such as a TCP connection.

Rules:

- wrap a file or a network connection when you do many small reads or writes
- do not wrap a reader that is already fully in memory unless you need `Scanner` or `Peek`
- flush writers
- handle `Scanner.Err`
- do not share one `bufio.Reader` or `bufio.Writer` across goroutines without a mutex

### Questions

#### Theoretical questions

1. Why are many small `Write` calls on a file slow?
2. What does `Flush` do on a `bufio.Writer`?
3. What is the default split function of `bufio.Scanner`?
4. What happens when a line is longer than the scanner token limit?
5. What does `Peek` return?
6. Why must you not share a `bufio.Reader` between goroutines without a lock?

#### Easy practical tasks

1. Write five lines to a file with `bufio.NewWriter` and `Flush`. Read them back with `bufio.Scanner`.
2. Forget `Flush` on purpose. Read the file. Then add `Flush`. Compare the file contents.
3. Use `ReadString('\n')` on a `bufio.Reader` that wraps a `strings.NewReader`.
4. Print `Buffered` on a writer before and after `Flush`.

#### Medium practical tasks

1. Benchmark 10,000 one-byte writes to a file with and without `bufio.Writer`. Report times.
2. Set a custom scanner buffer for a line that is 100,000 bytes. Scan that line from a reader.
3. Use `Peek(2)` to detect a UTF-8 BOM (`0xEF, 0xBB`) and then read the rest. Test with and without a BOM.

#### Advanced practical tasks

1. Write a line reader that uses `bufio.Reader.ReadSlice` or `ReadBytes` and that handles lines longer than one buffer. Add tests for empty lines and for a missing final newline.
2. Wrap a TCP echo connection (or a `net.Pipe`) with `bufio.NewReadWriter`. Send many small frames. Measure the effect of flush frequency.

---

## File read/write, directories, permissions

Package `os` opens files and directories. Package `path/filepath` builds paths for the operating system. Package `io/fs` describes a read-only file tree.

Open a file for read:

```go
f, err := os.Open(name)
if err != nil {
    return err
}
defer f.Close()
```

Create or truncate a file for write:

```go
f, err := os.Create(name) // mode 0666 before umask
```

Open with full control:

```go
f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o600)
```

Flags control create, append, truncate, and exclusive create. The permission bits (`0o600`) apply when the function creates the file. The process umask can clear bits. On Windows, permission bits have a smaller effect than on Unix.

Read a small file with `os.ReadFile`. Write a small file with `os.WriteFile`. Both close the file for you. Do not use them for huge files.

```go
data, err := os.ReadFile("config.json")
err = os.WriteFile("out.txt", data, 0o644)
```

Directories:

- `os.Mkdir` creates one directory
- `os.MkdirAll` creates the full path
- `os.ReadDir` lists a directory
- `os.Remove` removes a file or an empty directory
- `os.RemoveAll` removes a tree
- `os.Rename` renames or moves a file (not always atomic across disks)

`filepath.Join` joins elements with the correct separator. `filepath.Clean` removes extra separators. `filepath.EvalSymlinks` resolves links. For untrusted path parts, use `filepath.Join` with a fixed root and check that the result stays under that root. Go 1.24 adds `os.Root` for operations that stay inside a directory. Prefer that API when you need a jailed root.

File modes use `fs.FileMode`. `info.IsDir()` reports a directory. `info.Mode().Perm()` reports permission bits.

Always close files. A missed close leaks descriptors. `defer f.Close()` is the normal form. If you need the close error on write, assign `defer` to a named return or close explicitly after `Flush`.

Do not build paths with string `+` and `/` only. Windows uses `\`. `filepath` handles both.

### Questions

#### Theoretical questions

1. What is the difference between `os.Open` and `os.Create`?
2. When do the permission bits in `os.OpenFile` apply?
3. Why is `os.ReadFile` a bad choice for an unbounded file from a user?
4. What does `os.MkdirAll` create that `os.Mkdir` does not create?
5. Why do you use `filepath.Join` instead of string concatenation?
6. What problem does a fixed root (or `os.Root`) solve for user-supplied paths?

#### Easy practical tasks

1. Write `hello` to `out.txt` with `os.WriteFile`. Read it back with `os.ReadFile`. Print the string.
2. Create a directory `data` with `os.MkdirAll`. List it with `os.ReadDir`.
3. Open a file with `os.O_APPEND`. Write two lines in two open-close cycles. Show both lines in the file.
4. Print `filepath.Join("a", "b", "c.txt")` on your operating system.

#### Medium practical tasks

1. Copy a file with `os.Open`, `os.Create`, and `io.Copy`. Sync the destination with `f.Sync` if you need durability. Close both files.
2. Write a function that rejects a user path if it escapes a root directory. Test `../` input.
3. Set file permission `0o600` on Unix. Print `info.Mode()`. On Windows, write what the operating system allows.

#### Advanced practical tasks

1. Walk a directory tree with `filepath.WalkDir`. Count files and total size. Skip directories named `.git`.
2. Implement atomic replace: write to `name.tmp`, `Sync`, then `Rename` to `name`. Explain when this is atomic on your operating system.

---

## TCP/UDP with `net`

Package `net` provides TCP and UDP. `net.Listen` accepts TCP connections. `net.Dial` opens a client connection. Both return `net.Conn`.

`net.Conn` implements `io.Reader` and `io.Writer`. It also has `Close`, `SetDeadline`, `SetReadDeadline`, and `SetWriteDeadline`.

```go
ln, err := net.Listen("tcp", "127.0.0.1:0") // :0 selects a free port
if err != nil {
    return err
}
defer ln.Close()

go func() {
    conn, err := ln.Accept()
    if err != nil {
        return
    }
    defer conn.Close()
    io.Copy(conn, conn) // echo
}()

addr := ln.Addr().String()
c, err := net.Dial("tcp", addr)
```

For each `Accept`, start a goroutine if you handle connections in parallel. Close the connection when the handler finishes. Set deadlines so that a stuck peer does not block a goroutine forever.

```go
conn.SetDeadline(time.Now().Add(10 * time.Second))
```

A deadline applies to future I/O until you change it. `SetReadDeadline` and `SetWriteDeadline` split the two directions.

UDP is datagram-based. `net.ListenPacket("udp", addr)` returns a `net.PacketConn`. `ReadFrom` and `WriteTo` include the peer address. There is no accept loop. One socket can talk to many peers. UDP does not guarantee order or delivery.

```go
pc, err := net.ListenPacket("udp", "127.0.0.1:0")
n, addr, err := pc.ReadFrom(buf)
_, err = pc.WriteTo(buf[:n], addr)
```

Resolve addresses with `net.ResolveTCPAddr` or pass a string to `Dial` and `Listen`. Use `context` with `net.Dialer.DialContext` so that cancel stops a slow dial.

`net.Pipe` creates an in-memory pair of connections. Use it in tests instead of a real port when you can.

Do not use package `io/ioutil`. That package is deprecated. Use `io` and `os`. Do not listen on `:80` in examples that you run without rights. Use `127.0.0.1:0` or a high port.

### Questions

#### Theoretical questions

1. What interface does a TCP connection implement for bytes?
2. Why must you set a deadline on a connection?
3. What is the difference between `Listen` plus `Accept` and `ListenPacket`?
4. Does UDP guarantee that a datagram arrives?
5. How do you cancel a dial that takes too long?
6. Why is `127.0.0.1:0` useful in tests?

#### Easy practical tasks

1. Write a TCP echo server and client in one `main` test or program. Send `ping`. Read `ping` back.
2. Print the address of a listener that uses port `0`.
3. Send a UDP packet to yourself and print the payload.
4. Call `net.Dial` to a closed port. Print the error type or message.

#### Medium practical tasks

1. Accept many TCP connections. Handle each in a goroutine. Add a `WaitGroup` in a test that opens three clients.
2. Set a short read deadline. Show that a slow client causes a timeout error.
3. Use `net.Dialer` and `DialContext` with a cancelled context. Record the error.

#### Advanced practical tasks

1. Write a line-oriented TCP protocol: the server reads lines with `bufio.Scanner` and replies `OK <line>`. Add a test with `net.Pipe` or a local listen.
2. Write a UDP server that replies with the length of the datagram. Document what happens when the datagram is larger than your buffer.

---

## HTTP servers: `http.Handler`, `ServeMux` (Go 1.22+ patterns)

`http.Handler` is the server interface:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

`http.HandlerFunc` converts a function to a handler. `http.ResponseWriter` writes headers and the body. `http.Request` holds the method, URL, headers, body, and context.

Go 1.22 enhanced `http.ServeMux`. A pattern may include a method, a host, and path wildcards.

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /health", health)
mux.HandleFunc("GET /items/{id}", getItem)
mux.HandleFunc("POST /items", createItem)
mux.HandleFunc("GET /files/{path...}", getFile)
```

Rules for patterns:

- `GET /path` matches GET (and HEAD for GET patterns)
- `{name}` matches one path segment
- `{name...}` matches the rest of the path (including slashes)
- `{$}` matches the end of the path (`/{$}` is only `/`)
- a more specific pattern wins
- the mux reports a conflict when two patterns overlap in an ambiguous way at registration

Read a wildcard with `Request.PathValue`:

```go
func getItem(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintln(w, id)
}
```

`Request.SetPathValue` exists for other routers that want the same API.

Start a server:

```go
s := &http.Server{
    Addr:              "127.0.0.1:8080",
    Handler:           mux,
    ReadHeaderTimeout: 5 * time.Second,
}
err := s.ListenAndServe()
```

Set timeouts on the server. A missing `ReadHeaderTimeout` is a common review comment. Use `Server.Shutdown` with a context for graceful stop (later topics cover the full pattern).

`http.NotFound`, `http.Error`, and `http.Redirect` write common responses. Write the status with `WriteHeader` before the body. You can write the header only once in the useful sense: later `WriteHeader` calls do not change the status.

Do not use the global `http.DefaultServeMux` in a library. Build your own mux and pass it to `http.Server`. The global mux is easy to misuse with unexpected registrations.

Go 1.22 still accepts old patterns such as `/static/`. Those patterns keep trailing-slash redirect behavior. Prefer explicit method patterns for new routes.

### Questions

#### Theoretical questions

1. What method must an `http.Handler` implement?
2. How do you read `{id}` from a Go 1.22 pattern?
3. What does `{path...}` match that `{path}` does not match?
4. Why set `ReadHeaderTimeout` on `http.Server`?
5. What happens when two `ServeMux` patterns conflict?
6. Why must a library avoid `http.DefaultServeMux`?

#### Easy practical tasks

1. Register `GET /health` that writes `ok`. Call it with `httptest.NewServer` and `http.Get`.
2. Register `GET /items/{id}`. Request `/items/42`. Print `PathValue("id")` in the handler and in the test.
3. Register `POST /items`. Send a POST and a GET. Show that GET does not run the POST handler.
4. Use `http.Error` with status 400 in a handler. Check the status in a test.

#### Medium practical tasks

1. Add `GET /files/{path...}`. Request `/files/a/b.txt`. Print the wildcard value.
2. Register `/{$}` and `/`. Explain the difference in a comment. Test both paths.
3. Set `ReadHeaderTimeout` and one other timeout on `http.Server`. Start and shut down the server in a test.

#### Advanced practical tasks

1. Build a mux with `GET /users/{id}` and `GET /users/{id}/posts/{postID}`. Write table tests for five paths, including a miss.
2. Compare `ServeMux` patterns with a third-party router in a short note (six to ten sentences). Keep the example on the standard mux working without the third-party package.

---

## Middleware pattern

Middleware is a function that takes an `http.Handler` and returns an `http.Handler`. The returned handler runs code before and after `next.ServeHTTP`.

```go
func withLog(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}
```

You stack middleware from the outside:

```go
h := withLog(withAuth(mux))
```

The outer wrapper runs first on the way in. It runs last on the way out. Order matters. Put panic recovery on the outside. Put auth after logging if you want to log failed auth. Put request-id generation early so that later wrappers can read the id.

Common middleware jobs:

- log method, path, status, and duration
- recover from panics and write status 500
- require a header or a cookie
- set a timeout with `http.TimeoutHandler` or a context
- add security headers
- limit request body size with `http.MaxBytesReader`

To record the status code, wrap `http.ResponseWriter`:

```go
type statusWriter struct {
    http.ResponseWriter
    code int
}

func (w *statusWriter) WriteHeader(code int) {
    w.code = code
    w.ResponseWriter.WriteHeader(code)
}
```

If the handler never calls `WriteHeader`, `Write` sends status 200. Your wrapper must treat `code == 0` as 200 after the handler returns.

`http.TimeoutHandler` is standard middleware for a time limit. `http.AllowQuerySemicolons` and other adapters exist in `net/http`. Prefer small wrappers that you can test with `httptest`.

Do not put business rules that belong in one route into global middleware. Do not swallow errors without a log or a metric. Do not hold a mutex across `next.ServeHTTP` unless you intend to serialize all requests.

Pass values with `context.WithValue` only for request-scoped data such as a request id. Do not use context values for optional function parameters in general code. Topic 12 covers that rule.

### Questions

#### Theoretical questions

1. What is the type of a typical middleware function?
2. Which wrapper runs first when you write `withLog(withAuth(mux))`?
3. Why do you wrap `ResponseWriter` to log a status code?
4. What status do you assume when the handler never calls `WriteHeader`?
5. Name three jobs that fit middleware and one job that does not fit middleware.
6. Why is a mutex around `next.ServeHTTP` usually a mistake?

#### Easy practical tasks

1. Write `withLog` that logs method and path. Wrap a mux. Make one request.
2. Write `withHeader` that sets `X-App: 1` before `next`. Check the header in a test.
3. Wrap a handler with `http.TimeoutHandler` and a short limit. Sleep in the handler. Record the status.
4. Draw the call order of three wrappers around a mux.

#### Medium practical tasks

1. Write panic-recovery middleware. Panic in a handler. Assert status 500 and that the test process still runs.
2. Write `statusWriter` and log status plus duration. Test a handler that calls `WriteHeader(201)` and a handler that only writes a body.
3. Limit the body with `http.MaxBytesReader` in middleware. POST a large body. Record the error path.

#### Advanced practical tasks

1. Write request-id middleware that reads `X-Request-ID` or creates an id, stores it in the context, and writes it on the response. Read it from a later handler.
2. Implement a small middleware chain type `func Use(h http.Handler, mw ...func(http.Handler) http.Handler) http.Handler`. Test order with a buffer that each wrapper writes to.

---

## JSON APIs: decode, encode, validation

Use `encoding/json` for JSON APIs. Decode the request body with `json.NewDecoder(r.Body)`. Encode the response with `json.NewEncoder(w)`.

```go
type CreateItem struct {
    Name string `json:"name"`
    Qty  int    `json:"qty"`
}

func createItem(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

    var in CreateItem
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    if err := dec.Decode(&in); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }
    if err := dec.Decode(&struct{}{}); err != io.EOF {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }
    if in.Name == "" || in.Qty < 1 {
        http.Error(w, "invalid fields", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(map[string]string{"name": in.Name})
}
```

`Decoder` reads from a stream. `json.Unmarshal` reads a full `[]byte`. Prefer `Decoder` for HTTP bodies. Prefer `Encoder` for HTTP responses. `Encoder` adds a trailing newline. Clients accept that.

`DisallowUnknownFields` rejects extra keys. A second `Decode` that must return `io.EOF` rejects a body with two JSON values. Both checks help you reject sloppy input.

Validation is not a JSON feature. After decode, check fields in Go. Use extra types (`uuid`, `time.Time` with a custom format) when the domain needs them. Struct tags name JSON keys. Exported fields take part in encode and decode. Unexported fields do not.

`json.Number` and `Decoder.UseNumber` keep large numbers as strings. Default decode of a number into `any` uses `float64` and can lose integer precision.

Go 1.24 adds `omitzero` for some omit cases on zero values. `omitempty` still omits empty values as defined in the package documentation. Read the current `encoding/json` docs for the tag that you use.

Do not decode into `map[string]any` when you have a stable schema. A struct documents the API. Do not trust client data. Set `Content-Type` on success. Write a problem body with a stable shape for errors when you build a real API.

### Questions

#### Theoretical questions

1. Why prefer `json.Decoder` over `json.Unmarshal` for an HTTP body?
2. What does `DisallowUnknownFields` reject?
3. Why do you validate fields after a successful decode?
4. Which fields of a struct take part in JSON encode?
5. What is the risk of decoding a JSON number into `any`?
6. Why wrap the body with `http.MaxBytesReader` before decode?

#### Easy practical tasks

1. Decode `{"name":"a","qty":2}` into a struct. Print the fields.
2. Encode a struct to stdout with `json.NewEncoder`.
3. Send unknown field `{"name":"a","extra":1}` with `DisallowUnknownFields`. Print the error.
4. Write a handler that returns `400` when `name` is empty.

#### Medium practical tasks

1. Write `POST /items` as in this section. Use `httptest` to cover valid JSON, bad JSON, unknown field, and empty name.
2. Decode a large integer into `any` and into `int64` or `json.Number`. Compare the values.
3. Add `Content-Type: application/json` and a version field on every success response. Test the header.

#### Advanced practical tasks

1. Write a `decodeJSON` helper that limits bytes, disallows unknown fields, and rejects trailing data. Use it in two handlers. Add tests.
2. Design an error JSON object with `code` and `message`. Return it for decode errors and for validation errors. Keep status codes correct (400 versus 422 if you choose 422). Document the choice.

---

## WebSockets (high-level; `gorilla/websocket` or `nhooyr.io/websocket`)

A WebSocket is a persistent, bidirectional channel that starts as an HTTP request. The client sends an Upgrade request. The server switches the protocol if it accepts. After the handshake, both sides send messages: text, binary, and control frames (ping, pong, close).

The Go standard library does not provide a full WebSocket implementation for applications. You use a library.

`github.com/gorilla/websocket` is a widely used library. It upgrades with `Upgrader.Upgrade`. It reads and writes messages with `ReadMessage` and `WriteMessage`. You must run one reader loop. You must not write from two goroutines without a mutex. The project documents that concurrency model.

```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // do not use this in production
}

func echo(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer c.Close()
    for {
        mt, msg, err := c.ReadMessage()
        if err != nil {
            break
        }
        if err := c.WriteMessage(mt, msg); err != nil {
            break
        }
    }
}
```

`nhooyr.io/websocket` is the old module path. That module is deprecated. New code must use `github.com/coder/websocket`. The Coder module is the maintained continuation. It has first-class `context.Context` support and helpers in `websocket/wsjson`.

```go
c, err := websocket.Accept(w, r, nil)
defer c.CloseNow()
_, data, err := c.Read(r.Context())
err = c.Write(r.Context(), websocket.MessageText, data)
```

High-level rules for both libraries:

- treat the handshake as HTTP; check origin, auth, and method before upgrade
- do not set `CheckOrigin` to always true on a public server
- use TLS (`wss:`) on the public network
- close with a close code when you can
- apply backpressure; do not allow unbounded queues of outgoing messages
- cancel the context when the HTTP request ends

Use WebSockets for live updates, collaborative editors, and games. Do not use WebSockets when a single HTTP request or a short poll is enough. Load balancers and proxies must allow the Upgrade path.

This handbook stays high-level. You must read the chosen library documentation for ping periods and close codes before you ship.

### Questions

#### Theoretical questions

1. How does a WebSocket connection start?
2. Why does the standard library not replace a WebSocket library for most applications?
3. What is the current module path for the library that was `nhooyr.io/websocket`?
4. Why must you not accept every browser origin on a public server?
5. What concurrency rule does `gorilla/websocket` document for writes?
6. When do you keep HTTP instead of a WebSocket?

#### Easy practical tasks

1. Open the documentation for `github.com/gorilla/websocket` and `github.com/coder/websocket`. Write three differences in your own words.
2. Write an echo handler with one library. Use `httptest` only if the library supports it; otherwise write a small `main` and a note on how you tested.
3. List the message types that your library uses (text, binary, close).
4. Change `CheckOrigin` or the accept options to reject an empty origin. Record the API that you used.

#### Medium practical tasks

1. Add a ping loop or enable the library ping option. Document how a dead peer is detected.
2. Broadcast one message to many connections. Protect the connection set with a mutex. Do not write to a connection from two goroutines at once.
3. Close a connection with a normal close code. Show the client error or close status.

#### Advanced practical tasks

1. Build a small chat: HTTP page optional, WebSocket endpoint required. Limit the payload size. Drop the slowest client when the outbound queue is full.
2. Compare the two libraries on context cancel, JSON helpers, and default origin checks. Write a one-page recommendation for a new service.

---

## TLS basics with `crypto/tls`

TLS encrypts a TCP connection and authenticates the server (and optionally the client). HTTPS is HTTP over TLS. `wss` is WebSocket over TLS.

Package `crypto/tls` implements TLS. `http.Server.ListenAndServeTLS` uses it. You can also wrap a `net.Listener`:

```go
cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
cfg := &tls.Config{
    Certificates: []tls.Certificate{cert},
    MinVersion:   tls.VersionTLS12,
}
ln, err := tls.Listen("tcp", "127.0.0.1:8443", cfg)
```

`tls.Config` is the control type. Important fields:

- `Certificates` — server certificates
- `GetCertificate` — choose a certificate by Server Name Indication (SNI)
- `ClientCAs` and `ClientAuth` — client certificate authentication
- `RootCAs` — extra trust roots for a client
- `MinVersion` — refuse old protocol versions
- `NextProtos` — ALPN, for example `h2` and `http/1.1`

Do not set `InsecureSkipVerify` to true except in a local experiment. That flag skips server certificate checks. A public client must verify the server name and the chain.

For HTTP:

```go
err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", mux)
```

Or set `Server.TLSConfig` and call `ListenAndServeTLS`.

Local development uses a self-signed certificate. Browsers and `net/http` clients reject it unless you add the certificate to a trust pool or (only in a test) use a custom `RootCAs`. Tools such as `mkcert` install a local CA. Do not commit a private key to a public repository.

`crypto/tls` also provides `tls.Dial` and `tls.Dialer` for clients. `http.Client` uses TLS for `https` URLs. You can set `http.Transport.TLSClientConfig`.

Prefer TLS 1.2 as a minimum. Prefer TLS 1.3 when both sides support it (the handshake selects it). Do not invent your own crypto on top of raw TCP when TLS already solves the problem.

Certificates expire. Automate renewal (for example ACME). This topic only requires the Go API and the idea of a certificate plus a private key.

### Questions

#### Theoretical questions

1. What two services does TLS provide for a connection?
2. What does `MinVersion` prevent?
3. Why is `InsecureSkipVerify` dangerous?
4. What is SNI and which `tls.Config` field uses it to pick a certificate?
5. How does `http.ListenAndServeTLS` differ from `http.ListenAndServe`?
6. Why must you not commit `key.pem` to a public repository?

#### Easy practical tasks

1. Read `go doc crypto/tls.Config`. Write the purpose of `MinVersion` and `Certificates`.
2. Generate a self-signed certificate with `crypto/tls` documentation or `openssl` / `mkcert`. List the two file names that `LoadX509KeyPair` needs.
3. Write a client that sets `InsecureSkipVerify` in a test-only config. Mark the file as test-only in a comment.
4. Open a public `https` URL with `http.Get`. Print the status. You do not set a custom TLS config.

#### Medium practical tasks

1. Start `ListenAndServeTLS` in a test or a small program with a self-signed pair. Dial with a client that uses `RootCAs` loaded from the same certificate.
2. Set `MinVersion: tls.VersionTLS12`. Document how you would test a client that offers only an older version (describe the setup; you do not need to run old TLS if the toolchain blocks it).
3. Add HTTPS redirects or HSTS only as a written plan of three steps. Do not enable HSTS on a host that you do not control.

#### Advanced practical tasks

1. Configure optional client certificates (`tls.RequestClientCert` or `RequireAndVerifyClientCert`). Document the trust store (`ClientCAs`).
2. Wrap a raw TCP echo server with `tls.Listen` and a `tls.Dial` client. Reuse your echo protocol from the TCP section. Record the extra files and config fields that you needed.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of bytes from a file through `bufio` and `io.Copy` to a TCP connection. Name each interface.
2. How do Go 1.22 `ServeMux` patterns and middleware work together on one server?
3. What checks do you apply to a JSON request body before you trust the fields?
4. How do TLS, HTTP, and WebSockets relate on one listening port?
5. A teammate wants to ignore I/O errors, skip `Flush`, accept every WebSocket origin, and disable TLS verify. Which facts from this topic do you use in the review?

#### Easy practical tasks

1. Create a module `example.com/iox`. Add a file copy command, a `GET /health` server, and a JSON echo handler. Run tests with `httptest`.
2. Write a one-page cheat sheet: `io.Copy`, `LimitReader`, `bufio.Scanner`, `os.OpenFile` flags, `PathValue`, middleware signature, `DisallowUnknownFields`, WebSocket upgrade, `MinVersion`.
3. Use `curl` or `http.Get` against `httptest.NewServer` for `GET /health`. Save the status and body.
4. Draw a sequence diagram: client, TLS, HTTP handler, JSON decode.

#### Medium practical tasks

1. Build `GET /files/{name}` that reads a file from a fixed directory. Reject names that escape the directory. Return JSON errors.
2. Add logging middleware and a 1 MiB body limit to a `POST /items` JSON handler. Write three `httptest` cases.
3. Write a TCP client that sends one JSON object and a server that decodes it with `json.Decoder` on the connection.

#### Advanced practical tasks

1. Build a small HTTPS JSON API with a Go 1.22 mux, middleware, and a WebSocket route at `GET /ws`. Use a self-signed certificate in tests. Document how a browser client would trust it.
2. Load-test a handler that uses `bufio` versus a handler that writes small pieces without a buffer. Report times and a conclusion. Keep TLS off for the local test if you need a simpler setup. State that production still uses TLS.
