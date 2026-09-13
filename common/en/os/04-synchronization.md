# 4. Synchronization

## Description

Concurrent threads share data. Without rules, they race. This topic explains the critical section and the main locks: mutex, semaphore, and condition variable. You learn deadlock, livelock, starvation, priority inversion at a high level, and the difference between memory visibility and locks.

Complete this topic after threads and scheduling. Complete this topic before you share buffers between threads in later topics.

Use one term for each concept. A mutex is a lock that one thread owns. A semaphore is a counter that threads wait on. A condition variable is a wait queue that a thread uses with a mutex. Deadlock is a permanent wait cycle. Livelock is a state where threads keep changing but make no progress. Starvation is a state where one thread waits without bound while others proceed. Do not mix deadlock with livelock. Do not mix a mutex with a semaphore. Do not treat a sleep as a lock.

---

## Critical section

A critical section is a region of code that must not run on more than one thread at the same time when the threads use the same shared data.

The shared data is the reason for the section. Two threads may execute the same function. They may not both execute the increment of the same counter at the same time.

Properties that a correct lock protocol must give:

1. Mutual exclusion: at most one thread is in the critical section for that data.
2. Progress: if the section is free and threads want to enter, the decision who enters must not wait forever for a thread that is outside the section.
3. Bounded waiting (fairness goal): a thread must not wait forever while others enter again and again.

Keep a critical section small. Do not do file I/O or a long CPU loop inside the section unless the I/O is the shared object that you protect. A large section reduces parallelism.

Entry and exit must pair. If you forget to unlock, other threads block forever. Use a single exit path when you can. In C, check every early `return`.

Not every line of a program is a critical section. Local variables on a private stack do not need a lock. Immutable data that no thread writes does not need a lock. A race exists only when at least one access is a write and the accesses are concurrent and unsynchronized.

A race on `n = n + 1` is not one instruction on the CPU. The CPU loads `n`, adds one, and stores `n`. Two threads can both load the same old value.

### Questions

#### Theoretical questions

1. What is a critical section?
2. What is mutual exclusion?
3. Why must a critical section stay small?
4. Why do private stack variables not need a lock?
5. What happens if a thread never leaves the critical section?

#### Easy practical tasks

1. Mark the critical section in `n = n + 1` on a shared `n`. Write the three hardware-level steps.
2. Write five sentences about critical sections. Use only facts from this section.
3. Open `man 3 pthread_mutex_lock`. Write the pair of functions that wrap a critical section.
4. List three code regions that must not sit inside a lock if you can avoid it (sleep, network wait, large compile).

#### Medium practical tasks

1. Write a program with a shared list. Identify the critical sections for insert and for lookup. You may write comments only.
2. Draw two threads and a box labeled "critical section". Show a third thread that only reads a thread-local counter outside the box.
3. Write a racy counter with two threads. Add comments `/* enter */` and `/* exit */` around the increment. Then put a mutex there in the next section.

#### Advanced practical tasks

1. Read OSTEP on locks and critical sections. Write the three properties (mutual exclusion, progress, bounded waiting) in your own words with one example each.
2. Audit a small open-source C file that uses pthreads. List each critical section and the data that it protects. Write one section that looks too large.

---

## Mutex, semaphore, condition variable

A mutex (mutual exclusion lock) has two operations: lock and unlock. The thread that locks the mutex owns it. Only that thread may unlock it. A mutex is the default tool for a critical section. POSIX: `pthread_mutex_lock`, `pthread_mutex_unlock`. Linux futexes implement the wait (topic 14). You use the POSIX API first.

A semaphore is a counter with wait (P, down) and post (V, up). A binary semaphore can look like a lock, but it does not have an owner in the POSIX `sem_t` model. Any thread may post. A counting semaphore can allow N threads into a pool (for example, N buffer slots). POSIX: `sem_wait`, `sem_post`.

A condition variable lets a thread wait for a boolean condition. The thread must hold a mutex when it waits. `pthread_cond_wait` unlocks the mutex, waits, and locks the mutex again before it returns. Another thread changes the condition, then calls `pthread_cond_signal` or `pthread_cond_broadcast`. Always wait in a loop: `while (!ready) pthread_cond_wait(...)`. A wakeup can be spurious. Another thread can consume the condition.

Typical roles:

- mutex: protect a data structure
- semaphore: count permits or items
- condition variable: wait until a predicate is true (queue not empty, a phase finished)

A mutex is not recursive by default. A thread that locks the same non-recursive mutex twice deadlocks itself. A recursive mutex exists. Prefer a single lock owner path.

Do not use a condition variable without a mutex. Do not use `sleep` to wait for another thread.

### Questions

#### Theoretical questions

1. Who may unlock a mutex?
2. How does a POSIX semaphore differ from a mutex in ownership?
3. Why must `pthread_cond_wait` run while a mutex is held?
4. Why must you wait on a condition variable in a loop?
5. When do you pick a counting semaphore instead of a mutex?

#### Easy practical tasks

1. Open `man 3 pthread_mutex_lock`, `man 3 sem_wait`, and `man 3 pthread_cond_wait`. Write one sentence for each.
2. Write a three-row table: mutex, semaphore, condition variable. Add one typical use per row.
3. Write five sentences about these three tools. Use only facts from this section.
4. Draw a producer and a consumer: mutex, condition variable, queue.

#### Medium practical tasks

1. Fix the racy counter from the critical-section tasks with a `pthread_mutex_t`. Print `N * ITERS` as the expected value. This handbook does not contain the source.
2. Write a bounded buffer of size 1 with a mutex and two condition variables (`not_empty`, `not_full`), or with semaphores. Send ten integers.
3. Show a self-deadlock: one thread locks the same non-recursive mutex twice. Then stop the program. Write what you observed.

#### Advanced practical tasks

1. Implement a reader-writer pattern with a mutex and a condition variable (or `pthread_rwlock`). Document who may enter.
2. Read `man 7 pthreads` on process-shared mutexes. Write when a mutex in shared memory needs `PTHREAD_PROCESS_SHARED`.

---

## Deadlock, livelock, starvation

Deadlock is a permanent wait. Thread A holds lock 1 and wants lock 2. Thread B holds lock 2 and wants lock 1. Neither proceeds.

The Coffman conditions (all four must hold for this kind of deadlock):

1. Mutual exclusion: a resource has at most one owner.
2. Hold and wait: a thread holds one resource and waits for another.
3. No preemption: you cannot force the owner to release the resource.
4. Circular wait: a cycle of waiters exists.

Break one condition and the deadlock cannot form that way. A common user-space fix is lock ordering (topic 14): every thread acquires locks in the same rank order. Another fix is to try a lock and back off (`pthread_mutex_trylock`) with a clear policy. Another fix is one coarse mutex. The coarse mutex reduces deadlock and reduces parallelism.

Livelock is not a sleep cycle. Threads keep running and keep changing state, but they make no useful progress. Example: two threads detect a conflict, both back off, both retry in lockstep, both back off again. The CPU is busy. The work is not done.

Starvation is unfair delay. A thread is ready and is not deadlocked, but it does not get a turn for a long time. A greedy lock, a bad priority scheme, or a waiter that always loses a race can starve. Bounded waiting is the goal that fights starvation.

Do not call every hang a deadlock. First check: is the thread in `S` on a lock, in a livelock loop, or just slow I/O?

Dining philosophers is a teaching story: five threads, five forks, each needs two forks. A cycle of hold-and-wait can deadlock. An order on forks, or a waiter, can prevent the cycle.

### Questions

#### Theoretical questions

1. What is deadlock?
2. Name the four Coffman conditions.
3. What is livelock?
4. What is starvation?
5. How does a global lock order fight circular wait?

#### Easy practical tasks

1. Write a three-column table: deadlock, livelock, starvation. Add one example each.
2. Draw two threads and two mutexes in a deadlock cycle.
3. Open `man 3 pthread_mutex_trylock`. Write how a failed try differs from a blocking lock.
4. Write five sentences about these three failures. Use only facts from this section.

#### Medium practical tasks

1. Write a two-lock deadlock on purpose (thread A locks M1 then M2; thread B locks M2 then M1). Confirm a hang. Then fix it with an order. This handbook does not contain the source.
2. Read a short dining-philosophers description. Write which Coffman condition your fix breaks.
3. Describe a livelock backoff that uses a random delay. Write why a fixed delay can stay in lockstep.

#### Advanced practical tasks

1. Implement dining philosophers with a numbered-fork order. Document the rank. Do not paste a full solution from the internet as your only work.
2. Write a one-page note: how lockdep in the Linux kernel reports a possible deadlock cycle (high-level).

---

## Priority inversion (high-level)

Priority inversion happens when a high-priority thread waits for a lock that a low-priority thread holds, and a medium-priority thread keeps the CPU so that the low-priority thread cannot run and cannot unlock.

Story:

1. Low-priority thread L locks mutex M.
2. High-priority thread H needs M and blocks.
3. Medium-priority thread M runs and preempts L.
4. L does not unlock. H stays blocked. The "high" thread lost in practice.

The inversion is a scheduling plus locking bug. A lock order does not fix it. A smaller critical section helps L finish faster. Priority inheritance is a protocol: while H waits on M, L runs with H priority until L unlocks. Priority ceiling is another protocol: the lock owner gets a high ceiling priority.

Linux realtime mutexes can use priority inheritance (`PTHREAD_PRIO_INHERIT`). The default `SCHED_OTHER` desktop case is less strict. You still see inversion as "the UI froze because a low worker held a lock".

Do not raise a whole process to realtime priority to "fix" inversion. You can freeze the machine. Fix the lock scope. Use inheritance only with a documented realtime design.

This section is high-level. You do not implement a kernel inheritance protocol here.

### Questions

#### Theoretical questions

1. What is priority inversion?
2. What roles do the low, medium, and high threads play?
3. Why does lock ordering not fix inversion?
4. What is priority inheritance in one sentence?
5. Why is a small critical section useful against inversion?

#### Easy practical tasks

1. Draw L, M, H, and mutex M in the inversion story. Number the four steps.
2. Write five sentences about priority inversion. Use only facts from this section.
3. Open `man 3 pthread_mutexattr_setprotocol`. Write the name of the inheritance protocol constant if the page exists.
4. Make a table: "Symptom" and "Cause" for inversion versus deadlock.

#### Medium practical tasks

1. Write a scenario (no need for realtime on your laptop): a GUI thread waits on a lock that a background thread holds while it compresses a file. Write how you shrink the section.
2. Read a short note on the Mars Pathfinder inversion story (public essays). Write six sentences in your own words.
3. Compare inheritance and ceiling in a table: who gets a boost, and when the boost ends.

#### Advanced practical tasks

1. Read `man 7 sched` on realtime policies. Write why a lab on a shared machine must not set `SCHED_FIFO` without a rescue plan.
2. Write a one-page design: where inversion could appear in a thread pool that uses one mutex for the job queue.

---

## Memory visibility vs locks

A lock does two jobs. The first job is mutual exclusion. The second job is visibility. After thread A unlocks, thread B that locks the same mutex must see the writes that A did inside the section.

Visibility is a memory-model idea. CPUs and compilers can reorder loads and stores. A core can keep a write in a store buffer. Without a rule, thread B can see an old value even after A "already wrote" in program order.

A POSIX mutex lock and unlock are synchronization points. They include the barriers that you need for data that the mutex protects. You do not add extra fences around a correct mutex critical section in C.

Atomics (`stdatomic.h`, `__atomic_*`) give visibility for a single object with a stated memory order. They do not replace a mutex for a multi-word structure unless you design a lock-free algorithm. Topic 14 names barriers and lock-free ideas. This section only states the split.

A `volatile` qualifier in C is not a lock. It does not make `n = n + 1` atomic. It does not publish a buffer to another thread on all platforms the way a mutex does. Do not use `volatile` as a substitute for a mutex.

Sleep is not a visibility tool. `sleep(1)` can hide a race in a test. It does not fix the race.

If you share data without a lock or an atomic with a defined order, the program has a data race. The result is undefined in C.

### Questions

#### Theoretical questions

1. What two jobs does a mutex do?
2. Why can a CPU store buffer hide a write from another core?
3. Does `volatile` make an increment atomic? Explain.
4. When is an atomic enough, and when do you still need a mutex?
5. Why does `sleep` fail as a fix for a race?

#### Easy practical tasks

1. Write five sentences that contrast exclusion and visibility. Use only facts from this section.
2. Open `man 3 pthread_mutex_lock` and a C11 atomic overview. Write one sentence for each tool.
3. Make a table: "Tool" and "What it guarantees". Add rows for mutex, atomic, `volatile`, `sleep`.
4. Draw thread A write, unlock, thread B lock, read. Label the visibility point.

#### Medium practical tasks

1. Write a flag `ready` and a payload `data`. Publish with a mutex (set data, then ready). Consume with the same mutex. Then write why a plain `int ready` without a lock is unsafe.
2. Read `man 3 atomic_store` or C11 `memory_order_release` notes. Write six sentences on release and acquire as a preview of topic 14.
3. Use ThreadSanitizer (`gcc -fsanitize=thread`) on a racy counter if your compiler supports it. Write the report in your own words.

#### Advanced practical tasks

1. Write a one-page note: data race versus race condition (a logic race on a locked protocol). Give one example of each.
2. Read the OSTEP chapter on concurrent data structures. Write how a lock around a list gives both exclusion and visibility.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a critical section, a mutex, and visibility form one correct protocol?
2. A program hangs. Which questions distinguish deadlock, livelock, inversion, and a forgotten unlock?
3. Why is a counting semaphore a poor default substitute for every mutex?
4. How can starvation occur even when the Coffman cycle does not exist?
5. Why must two processes that share memory still obey the same lock and visibility rules as two threads?

#### Easy practical tasks

1. Write a cheat sheet: critical section, mutex, semaphore, condition variable, deadlock, livelock, starvation, inversion, visibility, data race.
2. Open `man 7 pthreads`. Write three functions that you will use in labs.
3. Draw a bounded buffer and label which tool protects the array and which tool waits for a slot.
4. List five rules: pair lock and unlock, wait in a loop, keep sections small, order locks, do not use `sleep` as a lock.

#### Medium practical tasks

1. Combine a racy counter lab and a two-lock order lab in one short report: what you saw before and after the fix. Do not paste full source into the report.
2. Use `strace -e futex` on a program that contends on a mutex. Write that the wait enters the kernel. Topic 14 names futex.
3. Design a thread pool queue on paper: mutex, condition variable, job list. Write who signals.

#### Advanced practical tasks

1. Read the OSTEP concurrency chapters. Write a one-page map from textbook lock, CV, and semaphore to POSIX names.
2. Audit a classmate program (or your own) with a checklist: every shared write has a lock or an atomic, every wait uses a loop, every lock has a rank.
