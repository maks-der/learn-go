# 9. Inter-Process Communication

## Description

Processes do not share an address space by default. They still need to exchange data and events. This topic explains the main Unix inter-process communication (IPC) tools: pipes and FIFOs, Unix domain sockets, signals, and shared memory with synchronization. The last section states when to use each tool.

Complete this topic after devices and I/O. Complete this topic before virtualization (topic 10). You already know file descriptors, `fork`, and mutexes.

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

1. Write a C program: `pipe`, `fork`, child writes a string, parent reads and prints it. This handbook does not contain the source.
2. Run `echo hello | cat` under `strace -e pipe,pipe2,read,write` if your `strace` shows the path, or write the same program and trace it.
3. Open a FIFO with two terminals: `cat /tmp/os-fifo-practice` and `echo hi > /tmp/os-fifo-practice`. Write the order that unblocked.

#### Advanced practical tasks

1. Implement a length-prefixed message protocol over a pipe. Send two messages. Show that a naive `read` of 4096 bytes can still be split correctly by your parser.
2. Read `man 7 pipe` on `F_SETPIPE_SZ` and `O_DIRECT` (packet mode). Write a one-page note on when packet mode helps.

---

## Unix sockets

A Unix domain socket (local socket) uses the `AF_UNIX` address family. Both ends live on one machine. There is no TCP stack and no port number in the internet sense.

Forms:

1. Pathname socket: a path in the filesystem (`bind` to `/tmp/app.sock`). Other processes `connect` to that path.
2. Abstract socket (Linux): a name that starts with a null byte. It does not use a directory entry. It disappears when no process holds it.

Stream Unix sockets (`SOCK_STREAM`) behave like a local TCP pipe: `listen`, `accept`, a bidirectional byte stream. Datagram Unix sockets (`SOCK_DGRAM`) send messages. `SOCK_SEQPACKET` preserves message boundaries on a connected socket.

Unix sockets can pass file descriptors with `SCM_RIGHTS` and credentials with `SCM_CREDENTIALS` (Linux). A privileged helper can open a file and pass the FD to an unprivileged process. Topic 11 uses this idea in sandboxes.

`connect` to a pathname socket needs permission on the path. Remove a stale socket file before `bind` if your server owns that path. `unlink` the path on a clean exit.

Unix sockets are often faster than TCP localhost because they skip IP and routing. They do not work across machines. Topic 12 covers internet sockets.

Do not mix `AF_UNIX` with `AF_INET`. The `sun_path` field is not an IP address.

### Questions

#### Theoretical questions

1. What address family do Unix sockets use?
2. How does a pathname socket differ from a Linux abstract socket?
3. What extra data can a Unix socket pass with `SCM_RIGHTS`?
4. Why can a Unix socket be faster than TCP on `127.0.0.1`?
5. What must a server do about a stale `.sock` file?

#### Easy practical tasks

1. Open `man 7 unix` and `man 2 socket`. Write one sentence for each.
2. Run `ss -xl` or `ls /tmp/*.sock` if any exist. Write one path that you see or write that you see none.
3. Write five sentences about Unix sockets. Use only facts from this section.
4. Draw server `bind`/`listen`/`accept` and client `connect` on a path.

#### Medium practical tasks

1. Write a stream Unix-socket echo: server binds `/tmp/os-unix-practice.sock`, client sends one line. Remove the path after. This handbook does not contain the source.
2. Use `strace -e socket,bind,connect,sendmsg,recvmsg` on the client. Write the call list.
3. Read `man 7 unix` on `SCM_RIGHTS`. Write six sentences on passing an FD.

#### Advanced practical tasks

1. Pass a file descriptor of an opened file from parent to child over a Unix socket pair (`socketpair`). The child reads the file through the received FD.
2. Write a one-page note: when a service uses a Unix socket in `/run` and how systemd socket activation relates (high-level).

---

## Signals

A signal is an asynchronous notification. The kernel or another process sends it. Examples: `SIGINT` (Control-C), `SIGTERM` (polite stop), `SIGKILL` (cannot catch), `SIGCHLD` (child state change), `SIGPIPE` (write to a closed pipe), `SIGSEGV` (bad memory access).

A process can:

- take the default action (exit, ignore, core dump, stop)
- ignore the signal (not `SIGKILL` or `SIGSTOP`)
- install a handler with `sigaction`

A handler runs in the process (on a thread). It must only call async-signal-safe functions. `printf` and `malloc` are not safe in a handler. A common pattern is: the handler sets a `volatile sig_atomic_t` flag. The main loop reads the flag.

`kill` sends a signal to a process. `kill -9` is `SIGKILL`. Prefer `SIGTERM` first. `raise` sends a signal to the caller.

Signals are not a data pipe. You do not send a file through a signal. You can use `sigqueue` and `siginfo` for a small integer. For bytes, use a pipe or a socket.

`fork` inherits handlers. `exec` resets handlers to default except ignored signals that stay ignored (with rules). Read `man 7 signal`.

Do not use signals as the only protocol between two busy servers. Races and handler limits make that design fragile. A self-pipe or `signalfd` on Linux can turn a signal into an FD for `epoll` (topic 8).

### Questions

#### Theoretical questions

1. What is a signal?
2. Which two signals can you not catch?
3. Why must a handler avoid `printf`?
4. What does `SIGCHLD` tell the parent?
5. Why are signals a poor byte-stream IPC?

#### Easy practical tasks

1. Open `man 7 signal` and `man 2 sigaction`. Write one sentence for each.
2. Run `sleep 30` and send `kill -TERM` from another terminal. Write the result.
3. Write five sentences about signals. Use only facts from this section.
4. Make a table: `SIGINT`, `SIGTERM`, `SIGKILL`, `SIGPIPE`. Add the default action.

#### Medium practical tasks

1. Install a `SIGINT` handler that sets a flag. The main loop prints and exits when the flag is set. Send Control-C. This handbook does not contain the source.
2. Draw sender `kill`, kernel, handler, flag, main loop.
3. Read `man 2 signalfd` if the page exists. Write how a signal becomes readable on an FD.

#### Advanced practical tasks

1. Combine `SIGCHLD` and `waitpid` in a loop that starts two children. Reap both. Avoid a busy loop.
2. Write a one-page note: async-signal-safe functions (`man 7 signal-safety`) and the self-pipe trick.

---

## Shared memory plus sync

Shared memory is a region that two or more processes map. A store in one process can be a load in another. There is no kernel copy on each store. Topic 5 named `MAP_SHARED`, `shm_open`, and `shmget`. This section adds the rule: shared memory without synchronization is a data race.

You must pair shared memory with a lock or with atomics that have a defined memory order (topic 4). Options:

1. A POSIX mutex in the shared region with `PTHREAD_PROCESS_SHARED`.
2. A POSIX semaphore (`sem_open` or a shared `sem_t`).
3. A futex on a shared word (topic 14). Harder.
4. A pipe or eventfd used only as a wakeup, while the payload sits in the map.

The protocol must define who writes which field and when the data is valid. A sequence number plus a mutex is a simple pattern. A ring buffer needs head and tail indexes and memory visibility.

Shared memory is fast for large payloads (a frame, a queue of records). Setup is harder than a pipe. Lifetime and cleanup matter: `shm_unlink` removes the name. The map can live until the last process unmaps.

Do not put a pointer that is valid only in one address space into the shared region unless both processes map at the same base (rare and fragile). Use offsets from the start of the map.

Security: any process that can map the object can read it. Use a strict mode on `shm_open` and a private directory. Topic 11 covers credentials.

### Questions

#### Theoretical questions

1. Why is shared memory fast?
2. Why must you add synchronization?
3. What mutex attribute do you need for a mutex in shared memory?
4. Why are raw pointers inside the shared region dangerous?
5. What does `shm_unlink` remove?

#### Easy practical tasks

1. Open `man 3 shm_open` and `man 3 pthread_mutexattr_setpshared`. Write one sentence for each.
2. Write five sentences about shared memory plus sync. Use only facts from this section.
3. Draw two processes, one shared frame, one mutex beside the payload.
4. Make a table: "Payload in the map" and "Who waits". Add mutex and semaphore rows.

#### Medium practical tasks

1. Map a shared integer with `shm_open` in two programs. Use a process-shared mutex or a named semaphore. Increment 100000 times in each. Print the total. This handbook does not contain the source.
2. List cleanup steps: `munmap`, `close`, `shm_unlink`. Write who runs `shm_unlink`.
3. Design a ring buffer on paper: head, tail, slots, and which lock protects which index.

#### Advanced practical tasks

1. Implement a small shared ring of 8 messages with a mutex and a condition variable (process-shared) or two semaphores.
2. Write a one-page note: when a pipe is enough and when a shared ring is worth the complexity.

---

## When to use each

Pick the smallest tool that matches the job.

Use a pipe or a FIFO when:

- you have a byte stream
- the processes are a parent and a child, or they can meet at a path
- you do not need random access or a second reader of the same bytes
- the shell `|` model fits

Use a Unix socket when:

- you need a server that accepts many clients on one machine
- you need datagrams or `SOCK_SEQPACKET` messages
- you need to pass file descriptors or credentials
- you want a path in `/run` like many system services

Use a signal when:

- you need a small event: stop, reload, child died
- you do not need a payload (or you only need `siginfo`)
- you accept handler rules

Use shared memory plus sync when:

- the payload is large or very frequent
- you can write a correct protocol
- you accept cleanup and permission work

Use internet sockets (topic 12) when the other end can be another machine. Use a file plus `flock` when the job is a lock on a resource, not a stream.

You can combine tools. A common pattern is: shared memory for data, eventfd or a socket for wakeup. Another pattern is: a pipe for bytes, `SIGCHLD` for death.

Do not start with shared memory for a first lab. Start with a pipe. Do not send large data through signals.

Message queues (`mq_open`, System V queues) exist. This path does not require them. If you need messages, a Unix datagram socket or a length-prefixed pipe is enough for labs.

### Questions

#### Theoretical questions

1. When is a pipe the right first choice?
2. When do you need a Unix socket instead of a pipe?
3. When is a signal enough?
4. When does shared memory win on performance?
5. Why do you still use internet sockets for some local designs?

#### Easy practical tasks

1. Write a five-row table: pipe, FIFO, Unix socket, signal, shared memory. Add one typical use each.
2. Write five sentences about choosing IPC. Use only facts from this section.
3. List three combinations (data tool plus event tool).
4. Open `man 7 ipc` if it exists. Write one sentence about the System V set (awareness).

#### Medium practical tasks

1. For each scenario, pick one tool: shell pipeline, local privileged helper that opens a file, reload a daemon, 4K video frames between two processes. Write one reason each.
2. Draw a decision flowchart: same machine? byte stream? need FD pass? large payload?
3. Compare a FIFO and a Unix stream socket in six sentences: open rules, bidirection, `accept`.

#### Advanced practical tasks

1. Write a one-page policy for a class project: default IPC, when to upgrade, and how you test cleanup.
2. Read `man 7 mq_overview` (optional). Write why this path still prefers sockets and pipes for beginners.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `fork` and a pipe form one parent-child channel?
2. Which IPC types are file descriptors, and why does that matter for `epoll`?
3. How do signals and `waitpid` work together without treating `SIGCHLD` as a byte pipe?
4. Why is shared memory plus a forgotten lock worse than a slow pipe?
5. A teammate wants one tool for all IPC. Which facts do you use to reject that plan?

#### Easy practical tasks

1. Write a cheat sheet: `pipe`, `PIPE_BUF`, FIFO, `AF_UNIX`, `SCM_RIGHTS`, `sigaction`, `SIGKILL`, `shm_open`, process-shared mutex.
2. Run `echo a | wc -c` and `ss -x | head`. Write one fact from each command.
3. Draw four boxes: pipe, Unix socket, signal, shared memory. Add one arrow of data or event.
4. Bookmark `man 7 pipe`, `man 7 unix`, and `man 7 signal`.

#### Medium practical tasks

1. Build a pipeline of two children from one parent (`ls` then `wc` style) with two `pipe`s or one pipe. Write a six-line report.
2. Use `strace -f -e pipe,socket,connect,kill` on `sh -c 'echo hi | cat'`. Write the IPC that you see.
3. Document an IPC debug checklist: unused pipe ends, stale `.sock`, `EPIPE`, missing `shm_unlink`, handler safety.

#### Advanced practical tasks

1. Combine a Unix socket for control and a shared-memory region for a payload. Write the protocol on one page. Implement a tiny demo if you have time.
2. Read the OSTEP IPC or concurrency chapters that match pipes and shared memory. Write a one-page map to POSIX names.
