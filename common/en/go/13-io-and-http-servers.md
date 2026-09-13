# 13. I/O and HTTP Servers

## Description

This topic shows how Go programs read bytes, write bytes, and serve HTTP. The `io.Reader` and `io.Writer` interfaces are the center of that design. Files, buffers, TCP, HTTP, and TLS all implement or consume those interfaces.

Use one term for each concept. A reader produces bytes. A writer consumes bytes. A handler serves an HTTP request. A pattern in `ServeMux` (Go 1.22 and later) matches a method and a path. Do not ignore errors from I/O calls. Close each file and each connection when you finish the work.

---

## `io.Reader` / `io.Writer` and `bufio`

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

`Write` writes `len(p)` bytes or returns an error. Check that `n` equals `len(p)` when you need a full write.

Useful helpers in package `io`:

- `io.Copy(dst, src)` copies from a reader to a writer until EOF
- `io.LimitReader(r, n)` stops after `n` bytes
- `io.MultiReader` and `io.MultiWriter`
- `io.TeeReader(r, w)` reads from `r` and writes a copy to `w`
- `io.NopCloser(r)` adds a no-op `Close` to a reader
- `io.ReadAll(r)` reads all remaining bytes (watch memory)

```go
limited := io.LimitReader(file, 1024)
n, err := io.Copy(os.Stdout, limited)
```

Do not use `io.ReadAll` on an unbounded network stream. Use `io.LimitReader` or a decoder that stops.

Package `bufio` adds buffering. `bufio.NewReader` reduces system calls on small reads. `bufio.NewWriter` collects writes. Call `Flush` on a writer before you need the data on the other side. `bufio.Scanner` reads lines or tokens.

```go
sc := bufio.NewScanner(f)
for sc.Scan() {
	line := sc.Text()
	_ = line
}
if err := sc.Err(); err != nil {
	return err
}
```

Set `sc.Buffer` when a line can be larger than the default token size. Check `sc.Err()` after the loop.

Prefer `io.Reader` and `io.Writer` in function parameters. Accept the interface. Return a concrete type when you own the implementation.

### Questions

#### Theoretical questions

1. What two results does `Read` return?
2. How must you handle `n > 0` and `io.EOF` in the same `Read`?
3. What does `io.Copy` do?
4. Why do you call `Flush` on a `bufio.Writer`?
5. Why is `io.ReadAll` dangerous on a network body?

#### Easy practical tasks

1. Copy a `strings.NewReader` to a `bytes.Buffer` with `io.Copy`. Print the buffer.
2. Wrap a reader with `io.LimitReader` of 3 bytes. Copy to stdout.
3. Scan three lines from a `strings.NewReader` with `bufio.Scanner`. Print each line.
4. Write to a `bufio.Writer` that wraps a `bytes.Buffer`. Print the buffer before and after `Flush`.

#### Medium practical tasks

1. Write a function `func slurp(r io.Reader, max int64) ([]byte, error)` that uses `LimitReader` and `ReadAll`.
2. Use `io.TeeReader` to copy input to stdout while you also decode the same bytes into a buffer.
3. Set a small scanner buffer. Feed a long line. Record the error. Increase the buffer and retry.

#### Advanced practical tasks

1. Implement `io.Reader` that returns one byte per `Read` from a string. Stop with `io.EOF`. Copy it with `io.Copy`.
2. Compose `MultiWriter` of a file and a `bytes.Buffer`. Prove both receive the same bytes.

---

## Files and directories

Package `os` opens files. `os.Open` opens for read. `os.Create` truncates and writes. `os.OpenFile` takes flags such as `os.O_APPEND` and `os.O_RDWR`. Check the error. Close the file with `defer f.Close()` after a successful open.

```go
f, err := os.Open(path)
if err != nil {
	return err
}
defer f.Close()
```

`os.ReadFile` and `os.WriteFile` cover small files. `os.ReadFile` loads the whole file into memory. Do not use it for huge files.

`os.Mkdir` creates one directory. `os.MkdirAll` creates parents. `os.Remove` removes a file. `os.RemoveAll` removes a tree. `os.Rename` renames or moves a file. `os.Stat` returns `fs.FileInfo`. Use `errors.Is(err, os.ErrNotExist)` for a missing path.

Package `path/filepath` joins and splits disk paths. Use `filepath.Join`, `filepath.WalkDir`, and `filepath.Abs`. Do not join disk paths with `+` and `/`.

```go
err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
		return nil
	}
	fmt.Println(path)
	return nil
})
```

`os.DirFS` and `embed.FS` implement `fs.FS`. You can open files through that interface. Topic 8 covers `embed`.

Always check I/O errors. A successful `Write` can still fail on `Close` when the system flushes data. Check the error from `Close` when a write must persist.

### Questions

#### Theoretical questions

1. When do you call `f.Close` after `os.Open`?
2. When do you choose `os.ReadFile` instead of `os.Open`?
3. How do you test that a path is missing?
4. Why must you not join disk paths with string `+` and `/`?
5. Why can `Close` return an error after a successful `Write`?

#### Easy practical tasks

1. Write a file with `os.WriteFile`. Read it back with `os.ReadFile`. Print the text.
2. Create a directory with `os.MkdirAll`. Write a file inside it.
3. Call `os.Stat` on a missing path. Detect `os.ErrNotExist` with `errors.Is`.
4. Join two path parts with `filepath.Join`. Print the result.

#### Medium practical tasks

1. Open a file, `defer` close, and copy to a `bytes.Buffer` with `io.Copy`.
2. Walk a small directory with `filepath.WalkDir`. Print each file name. Skip directories in the print.
3. Append one line with `os.OpenFile` and `os.O_APPEND`. Read the file and print it.

#### Advanced practical tasks

1. Write a function that copies one file to another with `io.Copy` and checks `Close` on both files.
2. Implement a small `ls` that prints mode, size, and name for each entry in a directory. Use `os.ReadDir`.

---

## TCP with `net`

Package `net` talks TCP and UDP. `net.Listen("tcp", addr)` listens. `ln.Accept()` returns a `net.Conn` for each client. `net.Dial("tcp", addr)` connects as a client.

```go
ln, err := net.Listen("tcp", "127.0.0.1:0")
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
	io.Copy(conn, conn)
}()

c, err := net.Dial("tcp", ln.Addr().String())
```

`net.Conn` implements `io.Reader` and `io.Writer`. You can wrap it with `bufio`. Set deadlines with `conn.SetDeadline`, `SetReadDeadline`, and `SetWriteDeadline`. A deadline is the main way to stop a stuck read.

Close the listener when the process stops. Close each connection when the exchange ends. Accept runs in a loop in a real server. Bound the work per connection with a context or a deadline.

`net.JoinHostPort` builds `host:port`. `net.SplitHostPort` splits it. Use these helpers instead of string format for IPv6.

Resolve names with `net.Resolver` or with `net.Dial`. Prefer `Dialer.DialContext` so that cancel works.

Do not use TCP as your first HTTP API. Use `net/http` for HTTP. Use `net` when you build a custom protocol or when you learn sockets.

### Questions

#### Theoretical questions

1. What does `net.Listen` return?
2. What interfaces does `net.Conn` implement?
3. Why do you set a deadline on a connection?
4. What helper builds a host and port string?
5. When do you use `net` instead of `net/http`?

#### Easy practical tasks

1. Listen on `127.0.0.1:0`. Print `ln.Addr()`. Close the listener.
2. Dial that address from the same program after `Accept` in a goroutine. Write one byte. Read it on the server side.
3. Set a read deadline in the past. Read. Print the error.
4. Split `"127.0.0.1:8080"` with `net.SplitHostPort`. Print host and port.

#### Medium practical tasks

1. Write an echo server: read a line, write the same line. Connect with `net.Dial` and `bufio`.
2. Use `Dialer{Timeout: 50 * time.Millisecond}` against a closed port. Record the error.
3. Accept two clients. Handle each in a goroutine. Close the listener from `main` after a delay.

#### Advanced practical tasks

1. Implement a tiny line protocol: `PING` returns `PONG`, `QUIT` closes. Test with a client in the same test file.
2. Cancel a `DialContext` while the connect is in progress (use a reserved or blackhole address with care). Record `ctx.Err()`.

---

## HTTP servers, mux, and middleware

Package `net/http` serves HTTP. `http.Handler` has one method: `ServeHTTP(ResponseWriter, *Request)`. `http.HandlerFunc` converts a function to a handler.

```go
http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
})
err := http.ListenAndServe(":8080", nil)
```

From Go 1.22, `ServeMux` matches method and path patterns. `"GET /items/{id}"` binds `id`. Read the value with `r.PathValue("id")`. Older Go versions need a third-party router or manual path parse. Prefer the standard mux in new code on Go 1.22 or later.

`http.ListenAndServe` uses `http.DefaultServeMux` when the handler is `nil`. Prefer an explicit `*http.ServeMux` in production. You then do not share global routes with tests.

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /health", health)
s := &http.Server{Addr: ":8080", Handler: mux}
err := s.ListenAndServe()
```

Middleware is a function that wraps a handler. The wrapper runs before and after the next handler.

```go
func withLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("req", "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
```

Set timeouts on `http.Server`: `ReadTimeout`, `WriteTimeout`, and `IdleTimeout`. Topic 15 covers graceful shutdown with `Shutdown`.

The server recovers from a panic in a handler so that one request does not stop the process. Log panics in your own middleware if you need extra data.

Always write a status. `Write` without `WriteHeader` sends `200`. Do not write the header twice.

### Questions

#### Theoretical questions

1. What method does `http.Handler` require?
2. What does a Go 1.22 pattern such as `GET /items/{id}` match?
3. Why do you prefer an explicit `ServeMux` over `nil` in `ListenAndServe`?
4. What is middleware?
5. Which `http.Server` fields set timeouts?

#### Easy practical tasks

1. Write `GET /health` that returns `ok`. Start the server. Call it with `http.Get` or `curl`.
2. Add `GET /items/{id}` and print `r.PathValue("id")`.
3. Wrap the mux with `withLog`. Show one log line per request.
4. Use `httptest.NewRecorder` to test `/health` without a listen port.

#### Medium practical tasks

1. Register `POST /items` and `GET /items/{id}`. Return 405 or 404 for a wrong method or path. Print the status.
2. Add middleware that rejects requests without a header `X-Token`. Return 401.
3. Set `ReadTimeout` and `WriteTimeout` on `http.Server`. Document the values that you chose.

#### Advanced practical tasks

1. Write middleware that limits the body with `http.MaxBytesReader`. POST a large body. Record the error.
2. Chain three middleware wrappers: log, auth, recover. Test each layer with `httptest`.

---

## JSON APIs

A JSON API reads JSON from the request body and writes JSON to the response. Use `json.NewDecoder(r.Body)` and `json.NewEncoder(w)`.

```go
type Item struct {
	Name string `json:"name"`
}

func create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var in Item
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(in)
}
```

Set `Content-Type` to `application/json`. Write the header before the body. Use `http.Error` for a simple text error. For a JSON error object, encode a struct and set the status first.

Limit the body. Wrap `r.Body` with `http.MaxBytesReader`. Do not `Decode` an unbounded body.

`Decode` reads one JSON value. Extra trailing data can be an error if you call `More` or decode again. Unknown fields are ignored by default. Use `DisallowUnknownFields` when the API must be strict.

Check the request method and path in the mux. Do not run a JSON decode on `GET` when there is no body.

Return useful status codes: `200` or `201` on success, `400` for bad input, `404` when a resource is missing, `409` for a conflict, `500` only for unexpected faults. Do not leak internal errors in the body.

### Questions

#### Theoretical questions

1. Why do you use `Decoder` on `r.Body` instead of `ReadAll` plus `Unmarshal` for large bodies?
2. When do you set `Content-Type`?
3. What status do you use for invalid JSON?
4. Why do you call `MaxBytesReader`?
5. When is `DisallowUnknownFields` useful?

#### Easy practical tasks

1. Write `POST /echo` that decodes `Item` and encodes it back. Test with `httptest`.
2. Send invalid JSON. Confirm status `400`.
3. Set `Content-Type` and check it on the recorder.
4. Add `json:"name"` on the struct. Send `{"name":"x"}`. Print the field.

#### Medium practical tasks

1. Add `GET /items/{id}` that returns JSON or `404`. Store items in a map with a mutex.
2. Enable `DisallowUnknownFields`. Send an extra key. Record the status.
3. Limit the body to 64 bytes. Send a large JSON object. Record the result.

#### Advanced practical tasks

1. Write a JSON error type with `code` and `message`. Use it for 400 and 404. Keep 500 as a generic message.
2. Decode an array of items from one POST. Validate each name. Return 400 if any item fails.

---

## TLS basics

TLS encrypts HTTP. HTTPS is HTTP over TLS. Package `crypto/tls` and `net/http` support TLS.

`http.ListenAndServeTLS(addr, certFile, keyFile, handler)` loads a certificate and a key from disk. The files are PEM.

```go
err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", mux)
```

`http.Server` also has `ListenAndServeTLS`. You can set `TLSConfig` on the server. `MinVersion: tls.VersionTLS12` is a common minimum.

For a client, `https://` URLs use TLS. The default client verifies the server certificate against system roots. `InsecureSkipVerify` disables that check. Do not use `InsecureSkipVerify` in production. Use it only in a local test with a known risk.

```go
tr := &http.Transport{
	TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
}
client := &http.Client{Transport: tr}
```

Create a local certificate for development with `crypto/tls` test helpers or with a tool such as `mkcert`. Do not commit a production private key. Topic 15 covers secrets.

`httptest.NewTLSServer` starts a TLS server in tests. The client from `srv.Client()` trusts that certificate.

TLS does not replace authentication. TLS protects the channel. You still check tokens or sessions.

### Questions

#### Theoretical questions

1. What two files does `ListenAndServeTLS` need?
2. What does the default HTTP client verify on HTTPS?
3. Why must you not set `InsecureSkipVerify` in production?
4. What is a common `MinVersion` for TLS?
5. Does TLS replace an application token check?

#### Easy practical tasks

1. Read `go doc http.ListenAndServeTLS`. Write the parameter list in your own words.
2. Start `httptest.NewTLSServer` with a handler. GET with `srv.Client()`. Print the status.
3. Print `srv.URL` and confirm that the scheme is `https`.
4. Write four sentences that separate TLS, HTTP, and application auth.

#### Medium practical tasks

1. Set `TLSClientConfig.MinVersion` on a transport. Call a public HTTPS URL. Print the status.
2. Point a client at the test TLS server with `http.DefaultClient`. Record the certificate error. Then use `srv.Client()`.
3. Document where you store `cert.pem` and `key.pem` on a development machine. Do not put keys in the module.

#### Advanced practical tasks

1. Generate a self-signed certificate with the standard library or `openssl`. Serve HTTPS on localhost. Curl with `-k` or with the cert. Record the steps.
2. Read `tls.Config` fields `Certificates`, `ClientAuth`, and `MinVersion`. Write one sentence for each field.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `io.Reader`, files, TCP connections, and HTTP bodies share one I/O model?
2. When do you choose `bufio.Scanner`, `io.Copy`, and `json.Decoder`?
3. How do a mux pattern, middleware, and a JSON handler split one request?
4. Why do deadlines exist on `net.Conn` and timeouts exist on `http.Server`?
5. What does TLS change and what does it not change for a JSON API?

#### Easy practical tasks

1. Write a file, scan its lines with `bufio`, and print the count.
2. Add `GET /health` and `POST /echo` JSON on one mux. Test both with `httptest`.
3. Copy a file to stdout with `io.Copy`. Check `Close`.
4. Draw a stack: TLS, HTTP, mux, middleware, JSON handler, `io.Reader` body.

#### Medium practical tasks

1. Build a small item API in memory: create and get by id. Limit body size. Return JSON errors.
2. Echo bytes over TCP in a test. Then expose the same echo as HTTP POST.
3. Wrap the API with log middleware and a `MaxBytesReader` limit. Prove both with `httptest`.

#### Advanced practical tasks

1. Serve the API with `NewTLSServer`. Call it from a client that decodes JSON. Shutdown is not required yet.
2. Write a handler that streams JSON lines with `json.Encoder` and `Flush` on `http.Flusher` when available. Test with a recorder.
