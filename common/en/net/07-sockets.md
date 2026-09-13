# 7. Sockets

## Description

A socket is the program interface to the transport layer. This topic shows `socket`, `bind`, `listen`, `accept`, and `connect`. You compare a TCP stream and a UDP datagram. You learn partial reads and writes. You learn blocking versus non-blocking, timeouts, and `SO_REUSEADDR`. You learn the standard library in your language.

Complete this topic after TCP and UDP. Pair it with `os.topics.md` for the kernel view of file descriptors.

Use one term for each concept. A socket is a file-like object that the kernel binds to a protocol and an address. A listen socket is not a connection. An accepted socket is one TCP connection. Do not mix `listen` with UDP. Do not skip errors.

Close the socket when you finish.

---

## `socket`, `bind`, `listen`, `accept`, `connect`

The BSD sockets API is the common model. Many languages wrap it. The call names stay the same idea.

`socket` creates a socket. You choose family (`AF_INET`, `AF_INET6`) and type (`SOCK_STREAM` for TCP, `SOCK_DGRAM` for UDP).

`bind` attaches a local address and port. A server binds a known port. A client can bind, or the kernel can pick an ephemeral port on `connect` or on the first send.

`listen` marks a TCP socket as a server. You pass a backlog. The kernel queues completed handshakes.

`accept` returns a new socket for one TCP connection. The original listen socket stays open for more clients.

`connect` on TCP runs the handshake to a remote address. `connect` on UDP sets a default peer. It does not handshake.

A TCP server order is: `socket`, `bind`, `listen`, loop `accept`. A TCP client order is: `socket`, optional `bind`, `connect`. A UDP server is: `socket`, `bind`, `recvfrom`. A UDP client is: `socket`, `sendto` or `connect` then `send`.

IPv6 sockets can use `AF_INET6` and a scope id for link-local addresses.

Do not `listen` on a UDP socket. Do not `accept` on UDP. Do not forget `bind` on a UDP server if you need a known port.

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
4. Write the names in your language that wrap `socket` and `listen`.

#### Medium practical tasks

1. Write a TCP echo server that binds `127.0.0.1` and a high port. Connect with a client. Use a machine that you own.
2. Write a UDP server and client for one message. Bind the server. Do not `listen`.
3. Write six sentences: why `accept` must happen in a loop for many clients.

#### Advanced practical tasks

1. Bind `::1` and `127.0.0.1` in two programs or one dual-stack design. Document which client address family works.
2. Write a one-page map: BSD name to your language API. Include `accept` and `connect`.

---

## TCP stream vs UDP datagram

A TCP socket is a byte stream. `read` can return fewer bytes than you asked. Two `write` calls can arrive as one `read`. One `write` can arrive as two `reads`. You must parse a protocol: length prefix, delimiter, or a parser such as HTTP.

A UDP socket is a datagram socket. One `recv` returns one datagram up to your buffer. If the buffer is too small, the tail can be truncated. Two datagrams stay two datagrams.

TCP has a connection. After `connect` or `accept`, `send` and `recv` use that peer. UDP can use `sendto` and `recvfrom` and can talk to many peers on one socket.

Errors differ. TCP can return 0 bytes on a clean close. UDP does not have that close signal. ICMP errors can surface on a connected UDP socket on some systems.

Choose the API that matches the protocol. Do not read TCP as if each `read` were one message. Do not send a 10 MB UDP datagram as if it were a stream.

In Go, `net.Conn` is stream-oriented for TCP. `net.PacketConn` is datagram-oriented. In Python, `SOCK_STREAM` and `SOCK_DGRAM` match the same split. Other languages use the same idea.

### Questions

#### Theoretical questions

1. Why can one TCP `read` be shorter than the request length?
2. What does one UDP `recv` return?
3. How do you find a message boundary on TCP?
4. What does a TCP read of zero bytes often mean?
5. Can one UDP socket talk to many peers?

#### Easy practical tasks

1. Write five sentences that compare stream and datagram sockets.
2. Make a table: two writes, what TCP read can see, what UDP recv can see.
3. Draw a length-prefix message on a TCP stream.
4. Write four sentences: truncation on a small UDP buffer.

#### Medium practical tasks

1. Write a TCP client that sends "hello" in two writes. Read on the server with a large buffer. Write how many reads you got.
2. Write a UDP server that uses a 4-byte buffer and a 20-byte send. Record what you receive.
3. Write six sentences: HTTP on TCP as a parser on a stream.

#### Advanced practical tasks

1. Implement a tiny length-prefixed echo on TCP. Send two messages back to back. Prove you parse both.
2. Write a one-page note: `net.Conn` versus `PacketConn` or the equivalent in your language.

---

## Partial reads and writes

A partial read means `read` returned n bytes and n is less than you asked. You must call `read` again if you need more. A partial write means `write` accepted n bytes and n is less than you passed. You must call `write` again for the rest.

Partial operations are normal on TCP. Signals, buffer limits, and non-blocking mode make them more common. Your loop must track an offset.

A language helper can hide the loop. Go `io.ReadFull` reads until the buffer fills or an error occurs. Python `socket.sendall` repeats `send`. You still must handle errors and short reads when you use raw calls.

UDP is different. A datagram send is usually all or an error. A short UDP receive is truncation, not "read the rest later" from the same datagram.

Do not ignore the return count. Do not assume that `write` of a 4 KB buffer always sends 4 KB in one call. Do not busy-loop without a timeout or an event wait.

Application protocols on TCP need a framing plan. HTTP/1.1 uses headers plus `Content-Length` or chunked encoding. You cannot wait for a single `read` of the full body unless you know the size.

### Questions

#### Theoretical questions

1. What is a partial read?
2. What is a partial write?
3. Why are partial operations normal on TCP?
4. How is a short UDP receive different from a partial TCP read?
5. What does `sendall` or `ReadFull` hide?

#### Easy practical tasks

1. Write five sentences that explain a read loop with an offset.
2. Make a table: asked 100, got 40. Add the next call.
3. Draw a 3-byte length prefix and a 10-byte body across two reads.
4. Write the helper name in your language that fills a buffer.

#### Medium practical tasks

1. Write a TCP server that reads 1 byte at a time and prints progress. Send 20 bytes from a client.
2. Write six sentences: `Content-Length` and why one `read` is not enough.
3. Force a partial write if you can (small send buffer on a lab). Record the return count.

#### Advanced practical tasks

1. Write a robust TCP reader that reads exactly N bytes or returns an error. Test N across two `write`s.
2. Write a one-page note: framing on TCP versus one datagram per message on UDP.

---

## Blocking vs non-blocking, timeouts, `SO_REUSEADDR`

A blocking call waits until the kernel can finish the work or an error occurs. `accept` blocks until a connection is ready. `read` blocks until data, close, or error. A hang on a dead peer can last a long time if you set no timeout.

A non-blocking socket returns an error such as `EAGAIN` or `EWOULDBLOCK` when the call would wait. The program uses `select`, `poll`, `epoll`, `kqueue`, or an async runtime. Topic 8 in the OS path covers I/O wait. This section needs the idea: do not block a single thread on one slow client if you must serve many clients.

A timeout is a deadline. Set a read timeout and a write timeout on sockets that talk to a network. A connect timeout stops a SYN that never gets an answer. The OS default can be minutes.

`SO_REUSEADDR` lets a process bind a port that still sits in `TIME_WAIT` from a previous instance in the common Unix case. It does not mean two servers share one stream. Details differ on Windows. `SO_REUSEPORT` is a different option on some Unix systems.

Do not wait forever on a public client. Do not set a 10 ms timeout on a large file over a slow path without a plan. Do not treat `SO_REUSEADDR` as a license to run two conflicting servers.

### Questions

#### Theoretical questions

1. What does a blocking `read` wait for?
2. What does a non-blocking call return when it would wait?
3. Why set a connect timeout?
4. What problem does `SO_REUSEADDR` ease after a restart?
5. Is `SO_REUSEPORT` the same option?

#### Easy practical tasks

1. Write five sentences that compare blocking and a timeout.
2. Make a table: `accept`, `read`, `connect`. Add one hang cause each.
3. Find how your language sets a socket deadline or timeout. Write the API name.
4. Write four sentences: restart of a server and `TIME_WAIT`.

#### Medium practical tasks

1. Write a TCP client with a short connect timeout to a closed or filtered port on a host that you own. Write the error.
2. Write six sentences: one thread per client versus non-blocking or a worker pool.
3. Start a server, stop it, start it again on the same port. Write whether you needed `SO_REUSEADDR`.

#### Advanced practical tasks

1. Write a server that sets a read deadline and closes idle clients. Test with a client that sends nothing.
2. Write a one-page note: `SO_REUSEADDR` on Unix versus Windows bind rules. Use public docs.

---

## Your language stdlib

The standard library wraps BSD sockets. You still need the ideas: family, type, bind, listen, accept, connect, deadlines, and close.

In Go, `net.Listen("tcp", "127.0.0.1:0")` binds and listens. `ln.Accept()` accepts. `net.Dial` connects. `net.ListenPacket` is UDP. Always `Close`. Use `context` or `SetDeadline` for timeouts. The course `go.topics.md` covers more Go I/O.

In Python, `socket.socket`, `bind`, `listen`, `accept`, `connect` match the BSD names closely. `create_server` and `create_connection` are helpers. Use `settimeout`.

Other languages have `java.net`, `System.Net.Sockets`, or `tokio`/`async`. Read the docs for:

- how to bind IPv4 and IPv6
- how errors appear (error values versus exceptions)
- how to shut down a write half (`shutdown`)
- whether the helper reads the full buffer

Do not copy a tutorial that ignores errors. Do not listen on `0.0.0.0` on a laptop that shares a LAN unless you want LAN clients. Prefer `127.0.0.1` for local tests.

Resolve names with the stdlib DNS helper, then connect. Topic 9 covers DNS. Happy Eyeballs may try IPv6 and IPv4 for you.

### Questions

#### Theoretical questions

1. What Go call listens on TCP in this section?
2. What Python module matches BSD names closely?
3. Why close a socket?
4. Why prefer `127.0.0.1` for a first echo test?
5. What extra job does a name-based `Dial` or `create_connection` do?

#### Easy practical tasks

1. Write five sentences that map BSD calls to your language.
2. Open the stdlib docs for listen and dial. Write two function names.
3. Make a table: listen, accept, dial, close. Add your language API.
4. Write a four-line program that listens on `127.0.0.1` and a high port, then exits.

#### Medium practical tasks

1. Write a TCP echo in your language. Test with `curl` or a second program.
2. Write six sentences: error handling on `Dial` or `connect` when the port is closed.
3. Bind `127.0.0.1:0` and print the chosen port. Connect to that port.

#### Advanced practical tasks

1. Add deadlines and a clean shutdown to your echo server. Document the API you used.
2. Write a one-page stdlib map for a teammate who knows BSD names only.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Why does a TCP server need two kinds of sockets after `accept`?
2. How do partial reads, blocking, and timeouts interact on one slow client?
3. What breaks if you treat UDP `recv` like a TCP `ReadFull` loop?
4. When is `SO_REUSEADDR` the wrong fix for "address already in use"?
5. How does the stdlib still require you to understand the 5-tuple?

#### Easy practical tasks

1. Write a cheat sheet: socket, bind, listen, accept, connect, stream, datagram, partial I/O, timeout, REUSEADDR.
2. Draw the TCP server call order and the TCP client call order on one page.
3. List the stdlib types you will use for TCP and UDP in your language.
4. Run your empty listen program. Show the port in `ss` or `netstat`.

#### Medium practical tasks

1. Write a lab report: echo server, two clients, logs of accept, and a clean close.
2. Compare a blocking echo with a version that uses a deadline. Write the hang you prevented.
3. Write six sentences: IPv6 `::1` echo and whether your stdlib dual-stack listen worked.

#### Advanced practical tasks

1. Pair this topic with `os.topics.md`: write which file-descriptor calls sit under `read` on Linux (`strace` or docs).
2. Write a small protocol: one UDP discover datagram, then a TCP connection to the printed port. Use a host that you own.
