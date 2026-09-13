# 21. Deadlocks, Locks, and SMPs (Advanced)

## Description

Topic 7 introduced mutexes, deadlock, and a first view of visibility. This topic is the next layer on a multiprocessor (SMP) Linux system. You learn lock ordering, the idea of RCU, per-CPU data, memory barriers at a high level, futexes, and lock-free progress conditions.

Complete this topic after internals and tracing. Complete this topic before the practice topic (topic 22). You already know critical sections and Coffman conditions. You now see how the kernel and libc implement waits on many cores.

Use one term for each concept. Lock ordering is a global rank of locks that every thread must follow. RCU is a Linux read-copy-update publish method for readers that do not take a writer mutex. Per-CPU data is a replica of a variable for each CPU. A memory barrier is a constraint on the order of memory operations. A futex is a kernel wait queue keyed by a user address. Wait-free and lock-free are progress guarantees for concurrent algorithms. Do not mix a barrier with a mutex. Do not mix lock-free with "a program that has no locks in the source".

---

## Lock ordering

Deadlock among locks needs a cycle (topic 7). A simple prevention rule is a total order. Every thread that needs more than one lock acquires them in the same rank order. The cycle cannot form.

How to apply the rule:

1. List every lock in the subsystem.
2. Give each lock a rank (a number or a documented name order).
3. Acquire from low rank to high rank (pick one direction and keep it).
4. If you hold a high lock and need a low lock, you must drop and retry, or you must redesign.

The kernel documents lock classes and uses lockdep. Lockdep is a runtime checker. It records the acquire order and reports a possible cycle. You can enable it in a debug kernel. User programs do not have lockdep by default. You write the order in comments and in review.

Ordering is per lock object, not only per type. Two mutexes of the same C type still need a rule (for example, order by address, or order by inode number).

A single global mutex avoids cycles and kills scalability (topic 7 review). Ordering is the middle path: many locks, one rank story.

Inversion (topic 7) is a different bug. Ordering does not fix inversion. Inheritance or a smaller critical section can.

Do not take locks in a signal handler. Do not invent a second order "only when this flag is set" unless you write a proof.

### Questions

#### Theoretical questions

1. How does a total lock order prevent a deadlock cycle?
2. What must you do if you hold a high-rank lock and need a low-rank lock?
3. What does lockdep check?
4. Why is an order by mutex address a valid rule?
5. Why does lock ordering not fix priority inversion?

#### Easy practical tasks

1. Write five sentences about lock ordering. Use only facts from this section.
2. Draw two threads and two locks with a good order and with a cycle.
3. Open `man 3 pthread_mutex_lock`. Write that POSIX does not rank your mutexes for you.
4. Make a table: one big mutex versus ordered fine locks. Add deadlock risk and parallelism.

#### Medium practical tasks

1. Take a two-lock homework (account A and account B). Write the rank rule (order by account id). Write the transfer steps.
2. Read a short lockdep overview in kernel documentation. Write six sentences on what a report means.
3. Audit a small pthread program. List every lock pair that can be held at once. Write the rank.

#### Advanced practical tasks

1. Implement dining philosophers with a global lock order (topic 7 advanced, now with an explicit rank comment in every acquire).
2. Write a one-page lock-order document for a toy kernel-like graph: inode lock, then page lock. Invent ranks. Do not claim it matches Linux.

---

## RCU idea (Linux)

Read-copy-update (RCU) is a Linux synchronization method. Readers do not take a shared mutex. A reader enters an RCU read-side critical section (`rcu_read_lock` in the kernel). The reader then follows pointers that the writer published.

A writer that must change a structure:

1. Copies the structure (or allocates a new one).
2. Updates the copy.
3. Publishes the new pointer with an ordered store.
4. Waits until every CPU (or every reader) has left the old read-side sections (a grace period).
5. Frees the old structure.

The grace period is the heart of RCU. The writer does not free memory that a reader can still see.

User space has userspace RCU libraries. Most application code should still use a mutex. RCU pays when there are many readers and few writers, as in routing tables and some file-descriptor tables.

RCU is not a general replacement for a mutex. Writers still serialize (often with a lock). Readers must not sleep in some RCU flavors. The rules are strict.

This section is the idea, not a license to roll your own list in a homework without a library.

Do not call `kfree` on a node that you just unlinked unless a grace period has ended.

### Questions

#### Theoretical questions

1. What do RCU readers avoid?
2. What is a grace period?
3. Why does the writer copy before publish?
4. When is RCU a good fit?
5. Why is RCU not a drop-in mutex replacement?

#### Easy practical tasks

1. Write five sentences about the RCU idea. Use only facts from this section.
2. Draw reader, old node, new node, publish, grace period, free.
3. Search `man 3 rcu` or write that user-space RCU is a library, not a POSIX mutex.
4. Make a table: reader-writer mutex versus RCU. Add reader cost.

#### Medium practical tasks

1. Read the Linux `Documentation/RCU/whatisRCU.rst` overview (online). Write six sentences in your own words.
2. Map RCU publish to the broken flag from topic 7. Write why a plain store is not enough.
3. List two kernel objects that documentation mentions as RCU-protected (example: dentry). Write why they are read-heavy.

#### Advanced practical tasks

1. Read about `synchronize_rcu` versus `call_rcu`. Write a one-page contrast: blocking wait versus callback.
2. Write why a user-space list with `free` on another thread is not RCU unless you implement a grace period. No need to implement one.

---

## Per-CPU data

A shared counter that every CPU increments needs a cache line that travels between cores. That traffic is expensive. Per-CPU data gives each CPU its own copy. A CPU updates only its copy. A reader that needs the global sum adds the copies (with rules).

Linux uses per-CPU variables for statistics, for some allocators, and for RCU bookkeeping. `this_cpu_inc` is the idea: no lock on the fast path if the CPU cannot migrate mid-update (or the update is atomic on that CPU).

User space can approximate this with thread-local storage (topic 5) or with an array indexed by CPU (`sched_getcpu`). Migration is the hard part. A thread can start an increment on CPU 0 and finish on CPU 1 if you are not careful. The kernel disables preemption or uses atomic ops for some per-CPU updates.

Benefits:

- less lock contention
- less false sharing (topic: one cache line, two variables)

Costs:

- more memory (one replica per CPU)
- a sum that is not instantaneously exact unless you synchronize

Per-CPU is not shared-memory IPC. It is a layout for one process or for the kernel.

Do not build a per-CPU array and then take a global mutex for every increment. You lost the point.

### Questions

#### Theoretical questions

1. What problem does per-CPU data reduce?
2. How do you get a global total from per-CPU counters?
3. Why is thread migration a problem for a naive user-space per-CPU increment?
4. What is false sharing at a high level?
5. Why does per-CPU data use more memory?

#### Easy practical tasks

1. Run `nproc`. Write the CPU count.
2. Write five sentences about per-CPU data. Use only facts from this section.
3. Draw two cores, two counters, and a sum step.
4. Open `man 3 sched_getcpu` if it exists. Write one sentence.

#### Medium practical tasks

1. Write a program with a global atomic counter and one with an array of atomics per `nproc`. Increment from many threads. Compare `perf stat` cache-miss counts if allowed.
2. Read a short Linux per-CPU overview. Write six sentences on `preempt_disable` as an idea (do not write a module).
3. Explain why `top` per-CPU percentages relate to this section but are not per-CPU variables in your process.

#### Advanced practical tasks

1. Implement a thread-local counter plus a stop-the-world sum (join threads, then add). Write when the sum is exact.
2. Write a one-page note on false sharing: two `int` counters in one struct, padding, and `perf` c2c if you have it.

---

## Memory barriers (high-level)

Modern CPUs and compilers reorder memory operations for speed. A write that appears first in your C file can become visible to another core after a later write. Topic 7 showed a broken `data` plus `ready` flag.

A memory barrier (fence) is a hardware or compiler operation that constrains that reorder. Names vary:

- A store-release makes prior stores in that thread visible before the release store.
- A load-acquire sees the release store and then may read the data that the publisher wrote.
- A full barrier orders both directions (stronger, slower).

C11 atomics (`memory_order_release`, `memory_order_acquire`, `memory_order_seq_cst`) express these constraints. A POSIX mutex lock and unlock include sufficient barriers. If you use a mutex, you do not add extra fences for that data.

Linux kernel code uses `smp_wmb`, `smp_rmb`, `smp_mb`, and acquire/release helpers. You see them if you read a kernel path (topic 20).

Rules for this path:

1. Prefer a mutex or C11 atomics with documented orders.
2. Do not sprinkle `asm volatile("" ::: "memory")` and call it a lock.
3. `volatile` is still not a barrier for other cores in the way students hope (topic 7).
4. x86 is stronger than ARM for some orders. A test that passes on a laptop can fail on a phone.

This section does not teach every ARM weak-memory litmus test. It teaches why topic 7 said "locks publish".

### Questions

#### Theoretical questions

1. Why do CPUs reorder memory operations?
2. What does a store-release publish?
3. What does a load-acquire pair with?
4. Why does a mutex remove the need for extra fences on the protected data?
5. Why can an x86 test hide a barrier bug?

#### Easy practical tasks

1. Write five sentences about barriers. Use only facts from this section.
2. Copy the topic 7 flag example and mark where a release and an acquire would sit.
3. Open a C11 `atomic_store_explicit` reference. Write the name of the release order.
4. Make a table: mutex, `memory_order_relaxed`, `memory_order_seq_cst`. Add "exclusion?" and "publish?".

#### Medium practical tasks

1. Write the flag example with C11 release/acquire. Explain in six sentences why `data` is visible.
2. Draw two cores and a store buffer. Show why a fence is needed (idea).
3. Read a short Linux `memory-barriers.txt` introduction (first page). Write three terms in your own words.

#### Advanced practical tasks

1. Compile the broken flag with `-O2` on an ARM VM or a weak-memory emulator if you have one. Write whether you can see a failure. If you only have x86, write why a pass is not a proof.
2. Write a one-page guide: when a student may use `relaxed` (counters) and when they must not (publishing a pointer).

---

## Futex

A futex (fast userspace mutex) is a Linux primitive. The fast path is a user-space atomic compare-and-swap on an integer. If the lock is free, the thread never enters the kernel. If the lock is busy, the thread calls `futex` to sleep on the address. The unlock side writes the integer and calls `futex` to wake waiters if needed.

POSIX `pthread_mutex` on Linux uses futexes. `FUTEX_WAIT` and `FUTEX_WAKE` are the classic operations. Later operations add requeue, pi (priority inheritance), and wait on a value.

A futex word must stay at a stable address. It can live in shared memory for a process-shared mutex (topic 15). The kernel hashes the address (and the mapping) to a wait queue.

Spurious wakes exist. The waiter must re-read the user integer in a loop. This matches the condition-variable `while` loop (topic 7).

You rarely call `futex` yourself. You call `pthread_mutex_lock`. Knowing futexes explains `strace` lines and explains why an uncontended mutex is cheap.

Do not implement a mutex from `futex` until you have read `man 2 futex` and a known-good algorithm. Lost wakes are easy.

### Questions

#### Theoretical questions

1. What is the futex fast path?
2. When does a thread enter the kernel for a mutex?
3. Why must the futex word have a stable address?
4. Why must a waiter loop after `FUTEX_WAIT`?
5. How does a process-shared mutex use a futex?

#### Easy practical tasks

1. Open `man 2 futex`. Write the purpose in one sentence.
2. Write five sentences about futexes. Use only facts from this section.
3. Run `strace -e futex` on a tiny pthread program that locks a mutex in one thread. Write whether a `futex` call appears (it may not if uncontended).
4. Draw user atomic, kernel wait queue, and wake.

#### Medium practical tasks

1. Write two threads that contend on one mutex in a tight loop. Run `strace -c -e futex`. Write the call count.
2. Read the man page on `FUTEX_WAKE`. Write six sentences on who wakes.
3. Compare a pipe wait (topic 15) with a futex wait. Write which one is for bytes and which one is for a word.

#### Advanced practical tasks

1. Read a short note on `FUTEX_LOCK_PI`. Write how it relates to priority inversion (topic 7).
2. Write a one-page trace guide: `strace -e futex,clone,sched_yield` on a contended mutex and on a condvar wait.

---

## Lock-free progress conditions (wait-free / lock-free)

Concurrent algorithms use progress words with a precise meaning.

Wait-free: every thread completes its operation in a bounded number of its own steps, even if other threads stop. No thread depends on another thread to run.

Lock-free: some thread completes an operation in a bounded number of steps. The system as a whole makes progress. One thread can starve if others keep winning a compare-and-swap.

Obstruction-free: a thread completes if it runs by itself for long enough. Others can abort it by running.

A mutex-based algorithm is typically blocking. If the lock owner stops (or sleeps forever), waiters never complete. That is not lock-free.

"Lock-free" does not mean "the source has no mutex". A lock-free queue uses atomics and retry loops. It can still be slower than a mutex. It is harder to prove.

A wait-free structure is even harder. Many "lock-free" student lists are neither correct nor lock-free.

For this path:

1. Use the words only with the definitions above.
2. Prefer a mutex until a measurement says otherwise.
3. C11 `atomic_compare_exchange` is a tool, not a proof of lock-freedom.
4. Progress of the algorithm is not progress of the I/O. A lock-free queue can still block on `write`.

Topic 7 livelock is related: threads run but the algorithm makes no useful progress. A bad CAS loop can livelock.

### Questions

#### Theoretical questions

1. What does wait-free mean?
2. What does lock-free mean?
3. Why is a mutex algorithm usually blocking?
4. Why is "no mutex in the source" not the same as lock-free?
5. How can a CAS loop livelock?

#### Easy practical tasks

1. Write five sentences about the three progress words. Use only facts from this section.
2. Make a three-row table: wait-free, lock-free, blocking. Add one example idea each (mutex, CAS retry, per-thread slot).
3. Open a C11 `atomic_compare_exchange_weak` page. Write that a failure means retry.
4. Draw two threads in a CAS loop on one word.

#### Medium practical tasks

1. Take a mutex queue and a textbook lock-free queue description (no need to implement both). Write six sentences on progress if one thread sleeps in the critical section versus in a CAS retry.
2. Explain obstruction-free in a storyboard of one thread that runs alone after a conflict.
3. Read a short note on the Michael-Scott queue or any famous lock-free queue. Write the progress claim only. Do not paste code from a book.

#### Advanced practical tasks

1. Implement a wait-free increment of a counter with `atomic_fetch_add`. Argue why that one operation is wait-free, and why a stack `pop` is harder.
2. Write a one-page warning: how to read a blog that says "lock-free" and which questions to ask (progress, ABA, reclamation).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do lock ordering, RCU grace periods, and futex waits each prevent a different failure (cycle, use-after-free, busy-spin)?
2. When do you pick a mutex, per-CPU counters, and an RCU-like publish for a read-heavy map?
3. Why do memory barriers appear inside a futex-based mutex even though you never call a fence API?
4. A classmate says a lock-free queue cannot deadlock. Which Coffman condition is gone, and which bugs remain?
5. How does topic 7 visibility preview this topic's release and acquire?

#### Easy practical tasks

1. Write a one-page cheat sheet: lock rank, lockdep, RCU grace period, per-CPU, release/acquire, futex wait/wake, wait-free, lock-free, blocking.
2. Draw one figure: two cores, a mutex word, a futex queue, and a fence on unlock.
3. Run `strace -e futex` and `nproc` and write one line that ties them to SMP.
4. Bookmark `man 2 futex`, C11 atomic orders, and Linux `whatisRCU.rst`.

#### Medium practical tasks

1. Write a design note for a stats counter: per-thread TLS plus a mutex only on scrape. State the progress and the error on the scrape.
2. Use ThreadSanitizer on a program that publishes a pointer without a release. Write the report in your own words.
3. Map each section of this topic to one `strace` or `perf` observation that you could collect on Linux.

#### Advanced practical tasks

1. Read the OSTEP chapters on locks, concurrency bugs, and the optional parallel chapters. Write a one-page map to this handbook.
2. Read a small Linux path that uses RCU (from topic 20). Write the publish and the reader lock names. Do not modify the kernel.
