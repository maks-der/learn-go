# 15. Inter-Process Communication

## Description

Processes do not share an address space by default. They still need to exchange data and events. This topic explains the main Unix inter-process communication (IPC) tools: pipes and FIFOs, Unix domain sockets, signals, shared memory with synchronization, and message queues. The last section states when to use each tool.

Complete this topic after I/O and devices. Complete this topic before virtualization (topic 16). You already know file descriptors, `fork`, and mutexes. You now connect two processes on purpose.

Use one term for each concept. A pipe is an anonymous kernel buffer with a read end and a write end. A FIFO is a named pipe that lives in the filesystem namespace. A Unix socket is a socket that uses a filesystem path or an abstract name on one machine. A signal is an asynchronous notification to a process or a thread. Shared memory is a mapping that more than one process can see. Do not mix a pipe with a Unix socket. Do not mix a signal with a return value from `wait`.

---

## Pipes and FIFOs

A pipe is a unidirectional byte stream inside the kernel. `pipe` creates two file descriptors: the read end and the write end. After `fork`, the parent and the child close the ends that they do not need. Data that the writer writes becomes data that the reader reads.

The shell `|` operator is a pipe. The left process writes stdout to the pipe. The right process reads stdin from the pipe.

Pipe rules:

1. Writes of small sizes (up to `PIPE_BUF` bytes) are atomic. Two writers do not interleave those writes.
2. If all write ends close, the reader gets end-of-file (`read` returns 0).
3. If all read ends close, a `write` causes `SIGPIPE` and the error `EPIPE`.
4. A pipe has a limited buffer. A full pipe blocks a blocking `write`.

A FIFO (named pipe) has a path. `mkfifo` creates it. Unrelated processes open the path. One process opens for read. Another opens for write. The kernel still implements a pipe buffer. The directory entry is only a meeting point.

Open of a FIFO blocks until both a reader and a writer exist (usual semantics). `O_NONBLOCK` changes that rule. Read `man 7 fifo`.

Pipes do not preserve message boundaries. Two `write` calls can become one `read`. If you need messages, define a length prefix or use another IPC type.

Do not use a FIFO as a lock file. Use `flock` or another documented lock.

### Questions

#### Theoretical questions

1. What does `pipe` return?
2. What happens when all write ends of a pipe close?
3. What is `PIPE_BUF`?
4. How does a FIFO differ from an anonymous pipe?
5. Why can one `read` return bytes from two `write` calls?

#### Easy practical tasks

1. Open `man 2 pipe` and `man 1 mkfifo`. Write one sentence for each.
2. Run `mkfifo /tmp/os-fifo-practice` and `ls -l` on it. Write the type letter (`p`). Remove the FIFO.
3. Write five sentences about pipes and FIFOs. Use only facts from this section.
4. Draw parent and child after `fork` with one pipe and the unused ends closed.

#### Medium practical tasks

1. Write a C program: `pipe`, `fork`, child writes a string, parent reads and prints it.
2. Run `echo hello | cat` under `strace -e pipe,pipe2,read,write` if your `strace` shows the shell path, or write the same program and trace it.
3. Open a FIFO with two terminals: `cat /tmp/os-fifo-practice` and `echo hi > /tmp/os-fifo-practice`. Write the order that unblocked.

#### Advanced practical tasks

1. Implement a length-prefixed message protocol over a pipe. Send two messages. Show that a naive `read` of 4096 bytes can still be split correctly by your parser.
2. Read `man 7 pipe` on `F_SETPIPE_SZ` and `O_DIRECT` (packet mode). Write a one-page note on when packet mode helps.

---

## Unix sockets

A Unix domain socket (local socket) uses the `AF_UNIX` address family. Both ends live on one machine. There is no TCP stack and no port number in the internet sense.

Forms:

1. Pathname socket: a path in the filesystem (`bind` to `/tmp/app.sock`). Other processes `connect` to that path.
2. Abstract socket (Linux): a name that starts with a NUL byte. No directory entry. The name lives in a kernel table.
3. Unnamed pair: `socketpair` creates two connected sockets. This is like a bidirectional pipe.

Unix sockets are bidirectional. They can be a stream (`SOCK_STREAM`) or datagrams (`SOCK_DGRAM`). Stream sockets preserve a byte stream. Datagram sockets preserve message boundaries.

Unix sockets can pass file descriptors (`SCM_RIGHTS`) and credentials (`SO_PEERCRED`, `SCM_CREDENTIALS`). A privileged helper can open a file and pass the descriptor to a weaker process. That feature is a reason to pick Unix sockets over a TCP port on localhost.

`SOCK_STREAM` Unix sockets look like TCP to the program (`listen`, `accept`, `read`, `write`). The address is a path, not `127.0.0.1`.

Permissions on the socket file and on the directory control who can connect. Abstract sockets use other checks. Topic 17 returns to isolation.

Do not expose a Unix socket on a shared `/tmp` without a safe directory or restrictive mode. Name collisions and stale files are common bugs. Unlink the path before `bind` when you own the service, or use a dedicated directory.

### Questions

#### Theoretical questions

1. What address family do Unix sockets use?
2. What is an abstract socket on Linux?
3. How does `socketpair` differ from `pipe`?
4. What can `SCM_RIGHTS` pass?
5. Why can a Unix socket be safer than TCP on localhost for local-only services?

#### Easy practical tasks

1. Open `man 7 unix` and `man 2 socketpair`. Write one sentence for each.
2. Run `ls -l /var/run` or `/run` and find one `.sock` name. Write the name. Do not connect to it.
3. Write five sentences about Unix sockets. Use only facts from this section.
4. Make a table: pipe, FIFO, `AF_UNIX` stream. Add "bidirectional?" and "has a path?".

#### Medium practical tasks

1. Write a tiny server and client with `AF_UNIX` `SOCK_STREAM` in `/tmp` (practice names). Send one line and close.
2. Use `socat` or `nc -U` if present to talk to your socket. Write the command.
3. Draw `SCM_RIGHTS`: process A opens a file, sends the FD, process B `read`s the same file.

#### Advanced practical tasks

1. Pass a file descriptor over a Unix socket with `sendmsg`/`recvmsg`. Document the cmsg layout from the man pages.
2. Read `man 7 unix` on pathname permissions and abstract names. Write a one-page hardening note for a local service.

---

## Signals

A signal is a limited asynchronous message. The kernel or another process sends a number (`SIGINT`, `SIGTERM`, `SIGCHLD`, `SIGUSR1`). The receiver can ignore some signals, take the default action, or run a handler.

Default actions include terminate, ignore, stop, and continue. `SIGKILL` and `SIGSTOP` cannot be caught. `SIGKILL` is the hard kill.

`kill` sends a signal to a process (the name is historical). `kill -TERM pid` requests a graceful stop if the process handles `SIGTERM`. `kill -9` is `SIGKILL`.

Rules for handlers:

1. A handler can interrupt a blocking system call. The call can fail with `EINTR` or restart, depending on flags.
2. Only async-signal-safe functions are safe in a handler. `printf` and `malloc` are not safe. `write` to a known fd can be safe. Prefer `signalfd` or a self-pipe: the handler only writes a byte, and the main loop reads it.
3. `sigaction` is the modern setup call. Prefer it over `signal`.

`SIGCHLD` tells a parent that a child changed state. The parent still must `wait` to reap the zombie (topic 4).

Signals are not a data channel. You cannot attach a payload of arbitrary size. Real-time signals can queue and carry an integer (`sigqueue`). Even then, a pipe or socket is the usual data path.

Do not use signals as a mutex. A signal can land between any two instructions from the program's point of view (with some kernel rules).

### Questions

#### Theoretical questions

1. What is a signal?
2. Which two signals cannot be caught?
3. Why is `printf` unsafe in a handler?
4. What is the self-pipe (or `signalfd`) pattern?
5. Why must a parent still `wait` after `SIGCHLD`?

#### Easy practical tasks

1. Open `man 7 signal` and `man 2 sigaction`. Write one sentence for each.
2. Run `kill -l`. Write five signal names and numbers.
3. Write five sentences about signals. Use only facts from this section.
4. Make a table: `SIGINT`, `SIGTERM`, `SIGKILL`. Add default action and whether you can catch it.

#### Medium practical tasks

1. Write a program that installs a `SIGTERM` handler that writes one byte to a pipe (or uses `signalfd`). The main loop prints "got term" and exits.
2. Run your program. Send `kill -TERM` from another terminal. Document the result. Then send `kill -KILL` to a copy that only handles `TERM`.
3. Use `strace -e signal=all` or `-e trace=signal` on `sleep 60` and send `SIGINT`. Write the signal name that appears.

#### Advanced practical tasks

1. Read the async-signal-safe list in `man 7 signal-safety`. Write a one-page allowed/forbidden list for a handler.
2. Combine `SIGCHLD` and a non-blocking `waitpid` loop. Spawn three short children. Confirm no zombies with `ps`.

---

## Shared memory + sync

Shared memory maps the same physical frames into two or more processes (topic 10). After setup, a store in one process can be visible in the other process without a `read` system call for each message.

Setup on Linux:

1. POSIX shared memory: `shm_open`, `ftruncate`, `mmap`.
2. `mmap` of a regular file with `MAP_SHARED`.
3. System V: `shmget`, `shmat`.

Shared memory is only a region of bytes. It does not serialize access. Two processes that increment the same counter without a lock still race (topic 7).

Synchronization options:

- A POSIX mutex in the region with `PTHREAD_PROCESS_SHARED`
- A POSIX semaphore (`sem_open` or `sem_init` with pshared)
- A futex (topic 21) behind those APIs
- Atomically updated fields with a documented memory order (harder)

You must agree on the layout: header, mutex, then payload. You must agree on who creates and who unlinks. System V segments can outlive the process. POSIX shm objects stay in `/dev/shm` until `shm_unlink`.

Shared memory is the fastest IPC for large payloads. It is the easiest to get wrong. One wild write corrupts the peer.

Do not put a pointer that is valid only in one address space into the shared region. Use offsets from the start of the mapping.

### Questions

#### Theoretical questions

1. Why is shared memory fast?
2. Why do you still need a lock or an equivalent?
3. What does `PTHREAD_PROCESS_SHARED` change?
4. Why must you store offsets instead of pointers?
5. How can System V shared memory leak?

#### Easy practical tasks

1. Open `man 3 shm_open` and `man 3 pthread_mutexattr_setpshared`. Write one sentence for each.
2. Run `ls /dev/shm`. Write what you see.
3. Write five sentences about shared memory and sync. Use only facts from this section.
4. Draw two address spaces and one physical frame.

#### Medium practical tasks

1. Map POSIX shm in two processes. Write a string in one. Read it in the other. Use a semaphore so the reader waits.
2. Increment a shared counter from two processes without a lock, then with a process-shared mutex. Write the two results.
3. Document unlink and lifetime for your object so a second run still works.

#### Advanced practical tasks

1. Build a small shared ring buffer with a mutex or a pair of semaphores. Move 10,000 messages of 64 bytes. Compare with a pipe of the same payload.
2. Read about `MAP_SHARED` plus `msync`. Write a one-page note: visibility to the peer versus durability on disk.

---

## Message queues

A message queue stores discrete messages in the kernel (or in a library). A reader receives one message at a time. Boundaries are preserved.

POSIX message queues (`mq_open`, `mq_send`, `mq_receive`) have a name (`/myqueue`). Each message has a priority. A queue has a max message count and a max message size. `mq_overview` documents the interface.

System V message queues (`msgget`, `msgsnd`, `msgrcv`) use a numeric key. They also preserve messages. They are older. Many new programs avoid them.

A queue can fill. A blocking send waits. A non-blocking send fails. A reader can block when the queue is empty.

Compared with a pipe, a queue gives messages and optional priority. Compared with a Unix datagram socket, a POSIX queue is a named kernel object with a different permission and persistence story.

Linux also has `netlink` and other special sockets. Those are not POSIX message queues. Do not mix the names.

Message queues are not popular in every codebase. Pipes and Unix sockets cover many needs. Learn queues so that you can read old systems and so that you know a kernel-backed mailbox exists.

### Questions

#### Theoretical questions

1. What does a message queue preserve that a byte-stream pipe does not preserve?
2. What extra field can a POSIX message have?
3. What happens when a POSIX queue is full?
4. How do you open a POSIX queue by name?
5. Why do many new programs still pick sockets over System V queues?

#### Easy practical tasks

1. Open `man 7 mq_overview` and `man 3 mq_open`. Write one sentence for each.
2. Run `ls /dev/mqueue` if the filesystem is mounted. Write what you see.
3. Write five sentences about message queues. Use only facts from this section.
4. Make a table: pipe, POSIX mq, Unix datagram socket. Add "message boundary".

#### Medium practical tasks

1. Write a sender and a receiver with POSIX `mq_send` and `mq_receive`. Send two messages of different lengths.
2. Set a small `mq_maxmsg` and fill the queue. Write the error on the next send.
3. Compare `ipcs -q` (System V) with `ls /dev/mqueue`. Write which namespace you used.

#### Advanced practical tasks

1. Use `mq_notify` to start a handler or a thread when a message arrives. Document the race if you also poll.
2. Read resource limits for message queues (`ulimit`, `/proc/sys/fs/mqueue`). Write a one-page admin note.

---

## When to use each

Pick the weakest tool that matches the need.

Use a pipe when:

- you have a parent and a child
- you need a byte stream
- you do not need a name
- the shell `|` model is enough

Use a FIFO when:

- unrelated processes must meet
- a path is an acceptable rendezvous
- you still want pipe semantics

Use a Unix socket when:

- you need a local server with `listen` and many clients
- you need datagrams or streams
- you need to pass file descriptors or read peer credentials
- you want one API that can later become TCP (with care)

Use signals when:

- you need a small event (stop, reload, child exited)
- you do not need a payload
- you accept asynchrony and handler limits

Use shared memory plus sync when:

- the payload is large or very frequent
- you can define a layout and a lifetime
- you test races

Use a message queue when:

- you want kernel-backed messages and priorities
- an existing system already uses that API

Also consider:

- TCP/UDP sockets (topic 18 and `net.topics.md`) when the peer can be on another host
- `eventfd` and `signalfd` when you want events on a poll set
- files plus `rename` when you need persistence (topic 12)

Do not add shared memory because it sounds fast. Measure. A pipe of 64-byte messages is often enough.

Do not mix five IPC types in one student project without a reason.

### Questions

#### Theoretical questions

1. When is a pipe the first choice?
2. When do you upgrade from a pipe to a Unix socket?
3. When is a signal the wrong tool?
4. When is shared memory worth the extra design?
5. When must you use a network socket instead of Unix IPC?

#### Easy practical tasks

1. Write a six-row table: pipe, FIFO, Unix socket, signal, shm, POSIX mq. Add one typical use.
2. Write five sentences that only restate the decision rules in this section.
3. Open `man 7 pipe`, `man 7 unix`, and `man 7 mq_overview`. Write one line that distinguishes them.
4. List three `eventfd` or `signalfd` uses from `man 2 eventfd` or `man 2 signalfd` if the pages exist.

#### Medium practical tasks

1. Take a two-process homework. Write which IPC you pick and which two you reject. Give one reason each.
2. Draw a decision tree: same parent? need messages? need many clients? need another host?
3. Compare localhost TCP and a Unix socket for a local database. Write six sentences (permissions, ports, FD passing).

#### Advanced practical tasks

1. Implement the same ping-pong of 100,000 small messages over a pipe and over shared memory. Write times and a conclusion.
2. Write a one-page IPC policy for a small Unix service: control plane (signals or a socket) and data plane (pipe or shm).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A parent starts a worker. The parent must send a config blob, then receive logs, then know when the worker dies. Which IPC pieces do you combine, and what do you avoid?
2. How do file descriptors unify pipes, FIFOs, and Unix sockets for `epoll` (topic 14)?
3. Why does shared memory without a memory-order or lock story fail even when the mapping is correct?
4. How do Unix socket credentials relate to user IDs in topic 17?
5. Which IPC objects remain after both processes exit, and which objects do not?

#### Easy practical tasks

1. Write a one-page cheat sheet: `pipe`, `mkfifo`, `AF_UNIX`, `socketpair`, `SIGTERM`/`SIGKILL`, `shm_open`, `mq_open`, when to use each.
2. Run `ls -l /dev/shm /dev/mqueue /tmp | head` and `kill -l | head`. Comment what each namespace is for.
3. Draw one figure: two processes, one pipe, one socket, one shared frame, one signal arrow.
4. Bookmark `man 7 pipe`, `man 7 unix`, `man 7 signal`, and `man 7 mq_overview`.

#### Medium practical tasks

1. Build a mini lab: FIFO for bytes, a Unix socket for a reply, and `SIGUSR1` for "reload". Document the protocol in ten lines. Do not over-engineer.
2. Use `lsof` or `/proc/<pid>/fd` on your lab processes. Write which FDs are pipes or sockets.
3. Trace one program with `strace -e trace=network,signal,ipc,desc`. Write three call names that match this topic.

#### Advanced practical tasks

1. Read the OSTEP chapters on concurrency and on IPC or pipes if present. Write a one-page map to this handbook.
2. Add `SCM_RIGHTS` to your lab so the parent opens a log file and the child only receives the FD. Document why the child never needs the path.
