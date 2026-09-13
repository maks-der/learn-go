# 5. Threads

## Description

A thread is an execution context inside a process. This topic explains user threads and kernel threads, the shared address space, and thread-local storage. You learn why programs create threads, and you get a first view of race conditions.

Complete this topic after processes. Complete this topic before scheduling and synchronization. Those topics assume more than one thread.

Use one term for each concept. A kernel thread is a thread that the kernel scheduler sees. A user thread is a thread that a user-space library schedules. Thread-local storage is data that belongs to one thread. A race condition is a bug that depends on the order of unsynchronized accesses. Do not mix a thread with a process. Do not mix a race with a deadlock (topic 7).

---

## User threads vs kernel threads

The kernel schedules kernel threads (also called light-weight processes on some systems). Each kernel thread has a kernel scheduling identity. On Linux, that identity is a task with its own thread ID (TID). The call `clone` with shared address space creates a kernel thread in the same process.

A user thread is a coroutine-like context that a library switches in user space. The kernel sees one kernel thread. The library may run many user threads on that one kernel thread. This model is N:1 (many user threads, one kernel thread).

A 1:1 model maps each user-visible thread to one kernel thread. POSIX threads (`pthreads`) on modern Linux use 1:1. `pthread_create` becomes a `clone` that shares memory.

An M:N model maps many user threads onto a pool of kernel threads. Some language runtimes use M:N. The runtime scheduler sits on top of the kernel scheduler.

Properties:

- 1:1: a blocking system call blocks only that kernel thread. Other kernel threads in the process can run. The kernel can place threads on different cores.
- N:1: a blocking system call can block all user threads that share the one kernel thread. The kernel cannot run two user threads of that process on two cores.
- M:N: more complex. The runtime must handle blocking and preemption.

Linux also has kernel-only tasks (workqueues, kthreads). Those are not your program's `pthread` threads. They run kernel functions. You see names such as `kworker` in `ps`.

This path uses POSIX threads as the default: one `pthread` is one kernel thread.

### Questions

#### Theoretical questions

1. What is a kernel thread?
2. What is a user thread in the N:1 model?
3. Why does a blocking `read` hurt all N:1 user threads on one kernel thread?
4. What model do modern Linux pthreads use?
5. What is an M:N thread model?

#### Easy practical tasks

1. Write a three-row table: N:1, 1:1, M:N. Add columns for "kernel sees" and "multicore".
2. Open `man 7 pthreads`. Write one sentence about Linux NPTL.
3. Run `ps -eL | head` or `ps -T -p 1`. Write the difference between PID and LWP or SPID if the columns appear.
4. Write five sentences that contrast user threads and kernel threads. Use only facts from this section.

#### Medium practical tasks

1. Write a C program with `pthread_create` that starts two threads. Each thread prints its TID with `gettid()` or `syscall(SYS_gettid)` and the PID with `getpid()`. Show that PIDs match and TIDs differ.
2. Draw N:1 and 1:1 diagrams. Label the kernel scheduler.
3. Run `cat /proc/self/status | grep -E 'Threads|Pid|Tgid'`. Write what `Tgid` means (thread group ID, the process PID).

#### Advanced practical tasks

1. Compare Go goroutines or Java virtual threads (documentation only) with pthreads. Write which model each one uses and what happens on a blocking OS call.
2. Read `man 2 clone`. Write which flags `pthread_create` needs so that threads share the address space (`CLONE_VM` and related flags).

---

## Shared address space

Threads of one process share the same virtual address space. They share:

- the program text
- global and static variables
- the heap
- open file descriptors (the file descriptor table is shared)
- the current working directory and credentials (unless a thread changes them with special calls)

Each thread has a private stack. Local variables in a thread function live on that stack. Each thread has its own registers and its own program counter when it runs.

If thread A writes a global variable, thread B can read the new value (subject to memory visibility; topic 7). If thread A writes through a heap pointer that thread B also holds, both threads see the same object.

This share is the reason threads are useful. It is also the reason races exist. Processes do not share an address space by default. Processes use pipes, sockets, or shared memory to talk. Threads already share memory.

A stack overflow in one thread can, in theory, hit another mapping. The kernel can place a guard page at the end of a thread stack. Topic 10 covers guard pages.

`fork` in a multithreaded process copies only the calling thread into the child. POSIX lists many restrictions. Avoid `fork` plus threads except for the `fork` then `exec` pattern.

### Questions

#### Theoretical questions

1. What memory do threads of one process share?
2. What stack does each thread have?
3. Why can two threads increment the same global counter without extra OS objects?
4. How does shared address space differ from two processes after `fork`?
5. Why is `fork` in a multithreaded process dangerous if you do not `exec`?

#### Easy practical tasks

1. Write a two-column table: "Shared" and "Private per thread". Add five rows.
2. Open `man 7 pthreads` and find the list of shared attributes. Write four shared items.
3. Draw one address space with text, heap, and two stacks.
4. Write four sentences about shared address space. Use only facts from this section.

#### Medium practical tasks

1. Write two threads that both print the address of a global variable and the address of a local variable. Show that the global address matches and the local addresses differ.
2. Let thread A `malloc` an integer, pass the pointer to thread B (or store it in a global), and let thread B write the integer. Print it from A after a join.
3. Compare `/proc/self/maps` before and after `pthread_create` of a thread with a large stack. Write what new mapping appeared if you can see it.

#### Advanced practical tasks

1. Read `man 2 fork` about multithreaded fork. Write the async-signal-safe rule and the `pthread_atfork` idea.
2. Write a program that opens a file in `main` and reads from it in two threads (take care with offset). Document whether the file offset is shared.

---

## Thread-local storage

Thread-local storage (TLS) is memory that the runtime and the kernel give to each thread. A TLS variable has one instance per thread. A write in thread A does not change the instance in thread B.

In C11, `_Thread_local` or `thread_local` marks a TLS variable. In pthreads, `pthread_key_create`, `pthread_setspecific`, and `pthread_getspecific` implement keys. Compilers also use TLS for the C library (`errno` is typically thread-local).

TLS is useful for:

- `errno` and similar per-thread error values
- a per-thread buffer that must not need a lock
- a thread name or a thread-local random state

TLS is not a lock. TLS does not make a shared heap object safe. If two threads hold a pointer to the same heap object, TLS does not protect that object.

The compiler and the dynamic loader set up TLS at thread start. The cost is small for a few variables. A huge TLS block in a program that creates thousands of threads wastes memory.

Do not store a pointer to a thread's stack object in another thread without a lifetime rule. TLS and stacks go away when the thread exits, unless you detach and leak, or you join and the stack is reused.

### Questions

#### Theoretical questions

1. What is thread-local storage?
2. How does a TLS variable differ from a global variable?
3. Why is `errno` typically thread-local?
4. Does TLS make a shared heap object safe? Explain.
5. What happens to TLS when the thread exits?

#### Easy practical tasks

1. Open `man 3 pthread_setspecific` or `man 3 errno`. Write one sentence about per-thread data.
2. Write a C program with `thread_local int n` (or `__thread int n`). Let two threads each set `n` to a different value and print `n`.
3. Make a table: "Global", "TLS", "Stack local". Add a row for who sees the value.
4. Write four sentences about TLS. Use only facts from this section.

#### Medium practical tasks

1. Implement a pthread key that stores a small struct. Set it in each thread. Get it and print a field. Delete the key in `main` after joins.
2. Show that `errno` after a failed `open` in thread A does not overwrite a later `errno` in thread B (force two different errors).
3. Draw TLS as a slot next to each thread stack, plus one shared heap.

#### Advanced practical tasks

1. Read how ELF TLS works (initial exec vs dynamic). Write a one-page overview for a beginner.
2. Measure RSS of a process that creates 500 threads with a large `thread_local` array versus a small TLS variable. Write the memory difference and the risk.

---

## Why threads exist (I/O wait, multicore)

A single thread that waits for I/O cannot use the CPU for other work in that thread. The process can still run other threads. That is the I/O reason for threads.

Example: a server waits for a client on a socket. If one thread blocks in `read`, another thread can handle a second client. The alternative is one thread and non-blocking I/O with `epoll` (topic 14). Threads are often simpler to write. They use more memory per connection.

The second reason is multicore. Two CPU-bound threads can run on two cores at the same time. Two processes can also use two cores. Threads share memory, so they avoid extra IPC for shared data. That share is both the benefit and the risk.

Threads are not free:

- each thread has a stack
- the kernel must schedule more tasks
- locks can serialize the work so that only one thread runs the critical part
- bugs become races

Use processes when you need a strong isolation boundary. Use threads when the tasks share much data and you accept the sync cost. Use a thread pool when you have many short tasks. Do not create one thread per tiny task without a bound.

A single-threaded program is easier to reason about. Add threads when you have a measured wait or a measured need for more cores.

### Questions

#### Theoretical questions

1. How do threads help when one call blocks on I/O?
2. How do threads help on a multicore CPU?
3. What is one alternative to threads for many sockets?
4. Why can locks cancel the multicore benefit?
5. When do you choose a new process instead of a new thread?

#### Easy practical tasks

1. Write five sentences: three reasons to use threads, two costs. Use only facts from this section.
2. Make a table: "I/O-bound work" and "CPU-bound work". Add two thread-related notes to each column.
3. Run `nproc` or `lscpu`. Write how many logical CPUs you have and the maximum useful CPU-bound threads before oversubscription (a first guess).
4. Open `man 3 pthread_create`. Write the meaning of the start function argument.

#### Medium practical tasks

1. Write two threads that each sleep for two seconds. Measure wall time of the program with `time`. Write why the wall time is about two seconds, not four.
2. Write two threads that each do a large CPU loop (busy work). Measure wall time on a machine with at least two cores. Compare with one thread that does both loops in sequence.
3. Draw a server: four client connections, a thread per client, and one shared log file. Mark where a lock will be needed (preview).

#### Advanced practical tasks

1. Compare a two-thread blocking server with a single-threaded `poll` loop in a short design document. Do not implement the full server unless you want the extra work.
2. Read about the C10k problem at a high level. Write why "one thread per connection" fails at large scale.

---

## Race conditions (preview of sync)

A race condition occurs when the result of a program depends on the uncontrolled order of operations on shared data.

Classic example: two threads increment a shared counter `n`. The increment is not one CPU operation in the abstract. It is a read, an add, and a write. Both threads can read the same value. Both threads can write the same new value. One increment is lost.

```c
/* unsafe */
n = n + 1;
```

The program can fail only sometimes. That makes races hard to find. A test can pass on one CPU and fail on two CPUs.

A race is not the same as a data race in the C memory model, but they overlap. A data race is concurrent conflicting accesses to a non-atomic object without synchronization. Topic 7 covers mutexes and visibility. This section only shows the problem.

Fixes that you will study next:

- a mutex around the critical section
- an atomic increment
- do not share the data (use TLS or a separate object per thread, then combine)

Do not "add a sleep" to hide a race. A sleep changes timing. It does not make the access safe.

Tools: ThreadSanitizer (`-fsanitize=thread`) on Clang or GCC can report many C/C++ data races. Use it in development.

### Questions

#### Theoretical questions

1. What is a race condition?
2. Why is `n = n + 1` unsafe on a shared `n` without a lock or an atomic?
3. Why can a racy program pass tests?
4. Why is `sleep` not a fix?
5. What is one safe design that avoids sharing the counter during the loop?

#### Easy practical tasks

1. Write four sentences that describe a race. Use only facts from this section.
2. Draw the lost-update timeline: two threads, three steps each, one lost increment.
3. Open `man 3 pthread_mutex_lock` as a preview. Write the two calls that pair around a critical section.
4. List three shared objects that often race: a counter, a linked list, a `FILE *` log.

#### Medium practical tasks

1. Write two threads that each increment a global `int` one million times without a lock. Print the result. Run the program ten times. Write the values that you see.
2. Add a `pthread_mutex_t` around the increment. Print the result. Confirm that the value is two million.
3. Build the unsafe program with `-fsanitize=thread` if your compiler supports it. Write the tool report in your own words.

#### Advanced practical tasks

1. Replace the mutex with C11 `atomic_int` and `atomic_fetch_add`. Compare the result and a rough `time` measurement with the mutex version.
2. Read the OSTEP concurrency intro. Write a one-page explanation of too-much-milk or lost update in terms of this section.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A process has three pthreads. What do they share, what do they not share, and who schedules them on Linux?
2. How do TLS, a stack local, and a heap object differ when two threads run the same function?
3. When does an extra thread help I/O, and when does it only add stack memory?
4. Why does the 1:1 model fit Linux servers better than N:1 for blocking sockets?
5. How does a race on a shared counter relate to the shared address space?

#### Easy practical tasks

1. Write a one-page cheat sheet: user thread, kernel thread, 1:1, TID, shared address space, TLS, I/O wait, multicore, race, mutex preview.
2. Run `ps -T -p $$` on your shell. Write how many threads you see.
3. Compile a two-thread "hello" program. Run `strace -f -e clone,exit ./a.out` and write the `clone` flags if they appear.
4. Draw one process box with two threads, one TLS slot each, and one racy global.

#### Medium practical tasks

1. Write a program with three threads: two increment a locked counter, one prints the counter every 100 ms until a join. Document the output shape.
2. Compare two processes that increment a counter in a file with two threads that increment a global. Write which pair needs IPC and which pair needs a mutex.
3. Use `/proc/<pid>/task` on your threaded program. List the TIDs and match them to your printouts.

#### Advanced practical tasks

1. Build a small thread pool (fixed N worker threads, a queue of function pointers). Run CPU tasks and I/O tasks. Write when the pool helps and when the queue needs a condition variable (topic 7).
2. Read Linux `man 2 gettid` and `man 7 pthreads`. Write a one-page map from POSIX thread to Linux thread group and TID.
