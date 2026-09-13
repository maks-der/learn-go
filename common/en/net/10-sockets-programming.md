# 10. Sockets Programming

## Description

A socket is the program interface to the transport layer. This topic shows the BSD socket calls, TCP streams versus UDP datagrams, blocking versus non-blocking I/O, `SO_REUSEADDR`, partial TCP reads and writes, timeouts, and the standard library in Go and in Python.

Use one term for each concept. A socket is a file-like object that the kernel binds to a protocol and an address. Complete this topic after TCP and UDP. Pair it with the operating-system topics for the kernel view.

Do not skip errors. A socket call can fail. Close the socket when you finish.

---

## BSD sockets: `socket`, `bind`, `listen`, `accept`, `connect`

The BSD sockets API is the common model. Many languages wrap it. The call names stay the same idea.

`socket` creates a socket. You choose family (`AF_INET`, `AF_INET6`) and type (`SOCK_STREAM` for TCP, `SOCK_DGRAM` for UDP).

`bind` attaches a local address and port. A server binds a known port. A client can bind, or the kernel can pick an ephemeral port on `connect` or on the first send.

`listen` marks a TCP socket as a server. You pass a backlog. The kernel queues completed handshakes.

`accept` returns a new socket for one TCP connection. The original listen socket stays open for more clients.

`connect` on TCP runs the handshake to a remote address. `connect` on UDP sets a default peer. It does not handshake.

A TCP server order is: `socket`, `bind`, `listen`, loop `accept`. A TCP client order is: `socket`, optional `bind`, `connect`. A UDP server is: `socket`, `bind`, `recvfrom`. A UDP client is: `socket`, `sendto` or `connect` then `send`.

Do not `listen` on a UDP socket. Do not `accept` on UDP. Do not forget `bind` on a UDP server if you need a known port.

IPv6 sockets can use `AF_INET6` and a scope id for link-local addresses.

### Questions

#### Theoretical questions

1. What two choices does `socket` need in this section?
2. What does `bind` attach?
3. What does `listen` do?
4. What does `accept` return?
5. How does `connect` differ on TCP and on UDP?

#### Easy practical tasks

1. Write five sentences that list the TCP server call order.
2. Draw a TCP server and a TCP client. Label each call.
3. Make a table: call, TCP server, TCP client, UDP server. Mark used or not.
4. Write the Go or Python names that wrap `socket` and `listen` in your language.

#### Medium practical tasks

1. Write a TCP echo server that binds `127.0.0.1` and a high port. Connect with a client. Use a machine that you own.
2. Write a UDP server and client for one message. Bind the server. Do not `listen`.
3. Write six sentences: why `accept` must happen in a loop for many clients.

#### Advanced practical tasks

1. Bind `::1` and `127.0.0.1` in two programs or one dual-stack design. Document which client address family works.
2. Write a one-page map: BSD name to Go `net` and to Python `socket`. Include `accept` and `connect`.

---

## TCP stream vs UDP datagram

A TCP socket is a byte stream. `read` can return fewer bytes than you asked. Two `write` calls can arrive as one `read`. One `write` can arrive as two `reads`. You must parse a protocol: length prefix, delimiter, or a parser such as HTTP.

A UDP socket is a datagram socket. One `recv` returns one datagram (up to your buffer). If the buffer is too small, the tail can be truncated. Two datagrams stay two datagrams.

TCP has a connection. After `connect` or `accept`, `send` and `recv` use that peer. UDP can use `sendto`/`recvfrom` and can talk to many peers on one socket.

Errors differ. TCP can return 0 bytes on a clean close. UDP does not have that close signal. ICMP errors can surface on a connected UDP socket on some systems.

Choose the API that matches the protocol. Do not read TCP as if each `read` were one message. Do not send a 10 MB UDP datagram as if it were a stream.

Go `net.Conn` is stream-oriented for TCP. `net.PacketConn` is datagram-oriented. Python `socket.SOCK_STREAM` and `SOCK_DGRAM` match the same split.

### Questions

#### Theoretical questions

1. Why can one TCP `read` not equal one application message?
2. What does one UDP `recv` return?
3. What happens if a UDP buffer is too small?
4. How does a TCP `read` of zero bytes often look?
5. Which Go types match stream and datagram?

#### Easy practical tasks

1. Write five sentences that compare stream and datagram APIs.
2. Draw two TCP writes and one TCP read that contains both.
3. Make a table: `read` zero, UDP no packet. Add what the program should do.
4. List three application tricks for TCP message borders (length, newline, parser).

#### Medium practical tasks

1. Write a TCP sender that writes "AB" then "CD" without a pause. Read with a 1-byte buffer. Print chunks.
2. Write a UDP sender that sends two datagrams. Read twice. Show two messages.
3. Write six sentences: why HTTP on TCP needs a parser, not one `read`.

#### Advanced practical tasks

1. Implement a 4-byte length prefix on TCP. Send two messages. Prove that a small `read` loop still works.
2. Write a one-page note: `net.Conn` versus `PacketConn` in Go with one example each.

---

## Blocking vs non-blocking

A blocking socket waits. `accept` waits for a client. `recv` waits for data. `connect` waits for the handshake. The thread or goroutine that calls the function does not continue until the call returns or an error occurs.

A non-blocking socket returns at once when it cannot complete. The error is `EAGAIN` / `EWOULDBLOCK` (or a language wrapper). The program must retry later. The program uses `select`, `poll`, `epoll`, `kqueue`, or `io_uring` on Unix, or I/O completion on Windows. Go uses a poller inside the runtime. You still write blocking-looking `Read` calls. The runtime parks the goroutine.

Non-blocking is not the same as a timeout. A timeout wakes you after a time. Non-blocking wakes you at once with "try later."

A server that uses one blocking thread per connection can hit thread limits. A server that uses non-blocking I/O or goroutines can hold many connections.

Do not spin in a tight loop on a non-blocking socket without a poller. You will waste the CPU. Do not assume Go `SetNonblock` is the usual path. The `net` package already integrates with the poller.

Python can use blocking sockets, `selectors`, or `asyncio`. Pick one model and stay with it in a small program.

### Questions

#### Theoretical questions

1. What does a blocking `recv` wait for?
2. What does a non-blocking call return when no data is ready?
3. How is a timeout different from non-blocking?
4. How does Go hide the poller from a simple `Read`?
5. Why is a busy loop on a non-blocking socket wrong?

#### Easy practical tasks

1. Write five sentences that explain blocking `accept`.
2. Make a table: blocking, non-blocking, timeout. Add one row for "when the call returns".
3. Write four sentences: one thread per client versus many goroutines.
4. Name the poller or selector API on your OS or in your language.

#### Medium practical tasks

1. Write a blocking TCP server. Connect with one client. Observe that a second client waits if you do not accept again yet.
2. In Python, set a socket to non-blocking and call `recv` with no data. Record the exception. In Go, write how `SetReadDeadline` differs.
3. Write six sentences: why a simple beginner server can stay blocking.

#### Advanced practical tasks

1. Write a small `select` or `selectors` server for two sockets, or a Go server with two goroutines. Document the model.
2. Write a one-page note: C non-blocking plus `epoll` versus Go `net`. High-level.

---

## `SO_REUSEADDR`

`SO_REUSEADDR` is a socket option. On many Unix systems it lets a server bind a port that still has connections in `TIME_WAIT` from a previous process. After a restart, `bind` succeeds.

The exact meaning differs on Windows and on Unix. On some systems `SO_REUSEADDR` also allows more than one bind to the same address. `SO_REUSEPORT` (Linux and others) is a different option that lets several sockets share a port for load spread.

Go `Listen` sets `SO_REUSEADDR` by default on some platforms. Python lets you `setsockopt(SOL_SOCKET, SO_REUSEADDR, 1)` before `bind`.

Do not treat `SO_REUSEADDR` as "two servers can serve the same port in a safe way on all OS." Read the platform rules. Do not use `SO_REUSEADDR` to hide a bug that never closes sockets.

A client usually does not need this option. A server that restarts often does.

`TIME_WAIT` on the local port of a client is a different problem (ephemeral ports). `SO_REUSEADDR` on the listener is the common teaching case.

### Questions

#### Theoretical questions

1. What problem does `SO_REUSEADDR` solve for a server restart?
2. Why does `TIME_WAIT` block a new `bind` without the option (Unix teaching case)?
3. How can Windows and Unix differ?
4. What is `SO_REUSEPORT` in one sentence?
5. Does a typical client need `SO_REUSEADDR`?

#### Easy practical tasks

1. Write five sentences that explain a server restart and `TIME_WAIT`.
2. In Python, write the `setsockopt` line. In Go, write where Listen documents reuse.
3. Make a table: `SO_REUSEADDR`, `SO_REUSEPORT`. Add one purpose each.
4. Write four sentences: why two accidental copies of a server can be dangerous if both bind.

#### Medium practical tasks

1. Start a TCP server, connect, stop the server, start it again on the same port. Record whether bind fails. Then set reuse if needed and retry on a lab host.
2. Read the Go `net` source comment or docs on reuse. Write three STE sentences.
3. Write six sentences: ephemeral port exhaustion versus listener `TIME_WAIT`.

#### Advanced practical tasks

1. Compare Windows `SO_REUSEADDR` documentation with Linux. Write a short table of differences. Use official or well-known docs.
2. Write a one-page restart runbook for a small TCP service: close, wait, reuse, health check.

---

## Partial reads/writes on TCP

A TCP `write` can accept only part of your buffer. The return value is the number of bytes that the kernel accepted. You must send the rest in a loop unless your language API documents a full write (Go `net.Conn.Write` can still be short in theory; use `io.Copy` or a loop. Many wrappers retry).

A TCP `read` can return any length from 1 to the buffer size (or 0 on close). You must loop until you have a full message. A fixed `read` of 1024 is not "one HTTP request."

`writev` / `sendmsg` and buffered writers (`bufio` in Go) reduce small writes. They do not change the stream model. You still parse on read.

Signals and interrupts can cut a call. Check the byte count and the error together. Process the bytes first when the API returns bytes plus an error.

Do not ignore a short write. Do not assume `read` fills the buffer. Do not use UDP rules on a TCP socket.

TLS and HTTP libraries handle this for you. When you write a raw protocol, you own the loops.

### Questions

#### Theoretical questions

1. What does a short write mean?
2. Why must you loop on TCP read for a message?
3. Why is "read 1024 bytes" not an HTTP request?
4. What should you do when a call returns bytes and an error?
5. Do HTTP libraries remove the need to know this?

#### Easy practical tasks

1. Write five sentences that explain a short write and a short read.
2. Draw a 100-byte message and three reads of 40, 40, and 20.
3. Make a table: TCP `Write` return, UDP `sendto` return. Add one note each.
4. Write a loop in pseudocode that sends all bytes.

#### Medium practical tasks

1. In Go or Python, send 100000 bytes. Read with a 100-byte buffer. Count how many reads you get.
2. Use `bufio.Reader` in Go or `makefile`/`readuntil` in Python to read a line. Write why a helper helps.
3. Write six sentences: `io.Copy` in Go as a full-copy loop.

#### Advanced practical tasks

1. Implement a read-exact function `readN(conn, n)`. Test with a slow sender.
2. Write a one-page note: how a length-prefix protocol uses partial reads safely.

---

## Timeouts

A timeout sets a deadline. If the call does not finish before the deadline, the API returns an error. Timeouts stop a thread from waiting forever when the peer is gone.

Set a timeout on `connect` so that a dead address fails. Set a timeout on `read` so that a silent peer does not hold a worker. Set a timeout on `write` on a slow or zero-window peer.

Go uses `SetDeadline`, `SetReadDeadline`, and `SetWriteDeadline` on `net.Conn`. Python uses `settimeout` on a socket, or `asyncio` wait. A deadline is often absolute time. A timeout is often a duration. Read your API.

Idle connections on the Internet can die without a FIN. A keepalive or an application ping plus a timeout detects that.

Do not use a 24-hour read timeout on a public server worker if a client can hold the slot. Do not use a 1 ms timeout on a WAN handshake.

A timeout error is not always a bug. It is a signal. Close the socket or retry by policy.

### Questions

#### Theoretical questions

1. What does a read timeout prevent?
2. Why does `connect` need a timeout?
3. What is the difference between a deadline and a duration timeout (idea)?
4. Why can an idle TCP connection die without a FIN?
5. Is a timeout always a program bug?

#### Easy practical tasks

1. Write five sentences that explain a read timeout on a server.
2. In Go, write a `SetReadDeadline` example in comments. In Python, write `settimeout`.
3. Make a table: connect timeout, read timeout, write timeout. Add one risk each if missing.
4. Choose a timeout for a LAN echo and for a public HTTP client. Write the two numbers and why.

#### Medium practical tasks

1. Write a client with a 1 second connect timeout to a black-hole address (a lab IP that does not answer). Record the error.
2. Write a server that sets a short read deadline. Do not send data from the client. Record the timeout.
3. Write six sentences: keepalive versus application ping.

#### Advanced practical tasks

1. Write a Go client that uses `context` with a timeout around `Dial` or HTTP. Document the cancel path.
2. Write a one-page policy: timeouts for an API server (connect, headers, body). Use numbers that you can defend.

---

## Your language's stdlib (`net` in Go, `socket` in Python)

Go package `net` is the standard library for sockets. `net.Listen("tcp", addr)` returns a `Listener`. `Accept` returns a `Conn`. `net.Dial("tcp", addr)` returns a `Conn`. `net.ListenPacket("udp", addr)` is UDP. Use `net.JoinHostPort` and `net.SplitHostPort`. Resolve names with `net.Resolver` or let `Dial` resolve.

Prefer `net.Conn` interfaces in functions. Use `tcpAddr` only when you need the extra fields. Always check errors. Close with `defer conn.Close()` when you own the connection.

Python `socket` is closer to BSD names. `socket.socket(AF_INET, SOCK_STREAM)`, `bind`, `listen`, `accept`, `connect`. `socket.create_server` and `socket.create_connection` are helpers. `socket.timeout` is the timeout exception.

Both languages have HTTP in the standard library or in common packages. This topic is the raw socket. Use HTTP libraries for HTTP. Use raw sockets when you learn or when you write a custom protocol.

Do not mix blocking Python sockets with two threads without a lock on one socket. Do not ignore `Close` in Go. Do not use `Listen("tcp", ":80")` on a privileged port without rights.

IPv6 in both libraries needs brackets in some address strings. Test `::1` and `127.0.0.1`.

### Questions

#### Theoretical questions

1. What does `net.Listen` return in Go?
2. What does `net.Dial` return?
3. Which Python module matches BSD sockets most closely?
4. When should you use an HTTP library instead of raw sockets?
5. Why must you check errors on every socket call?

#### Easy practical tasks

1. Write a Go or Python TCP client that connects to `127.0.0.1` on a port that you open with a server.
2. Write five sentences that map `Listen`/`Dial` to `listen`/`connect`.
3. Read `go doc net.Dial` or Python `help(socket.create_connection)`. Write the purpose in one sentence.
4. Make a table: Go function, Python function. Add four rows.

#### Medium practical tasks

1. Write a Go echo server in `net` and a Python client, or the reverse. Use the same port on localhost.
2. Resolve a name with the stdlib. Print IPv4 and IPv6 addresses.
3. Write six sentences: `defer Close` in Go versus `with` or `close()` in Python.

#### Advanced practical tasks

1. Write a small UDP ping in both languages or in one language with tests. Use timeouts.
2. Write a one-page style guide for your team: when to use `net/http` versus `net.Conn`.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a TCP server from `socket` to a safe `read` loop with a timeout.
2. How do stream semantics, partial I/O, and a length prefix work together?
3. When do you choose blocking plus goroutines, and when do you choose an explicit poller?
4. Why do `SO_REUSEADDR` and `TIME_WAIT` appear in the same operational story?
5. A teammate treats every `recv` as one UDP-like message on TCP. Which facts do you use to correct that design?

#### Easy practical tasks

1. Write a one-page cheat sheet: BSD calls, stream vs datagram, blocking, reuse, partial I/O, timeouts, Go and Python names.
2. Draw a flowchart: create, bind, listen or connect, I/O loop, close.
3. List eight errors you have seen or expect (bind in use, refused, timeout, reset). Write one cause each.
4. Open `go doc net` or Python `socket` docs. Bookmark the `Listen`/`accept` page.

#### Medium practical tasks

1. Write a TCP echo server and client in your language. Add a read timeout. Document how to run both.
2. Write a UDP echo pair. Show that two sends stay two messages.
3. Write a short report: what happens if you omit `SO_REUSEADDR` and restart the server quickly.

#### Advanced practical tasks

1. Write a protocol with a 4-byte length and a payload. Include tests for short reads. Use only localhost.
2. Build a glossary of 12 socket terms. Each entry: term, one sentence, one Go or Python symbol.
