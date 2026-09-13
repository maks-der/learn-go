# 7. Concurrency and Synchronization

## Description

Concurrent threads share data. Without rules, they race. This topic explains the critical section and the main locks: mutex, semaphore, and condition variable. You learn deadlock and the four Coffman conditions, livelock, starvation, priority inversion, and a first view of memory visibility.

Complete this topic after threads and scheduling. Complete this topic before you share buffers between threads in later projects.

Use one term for each concept. A mutex is a lock that one thread owns. A semaphore is a counter that threads wait on. A condition variable is a wait queue that a thread uses with a mutex. Deadlock is a permanent wait cycle. Do not mix deadlock with livelock. Do not mix a mutex with a semaphore. Do not treat a sleep as a lock.

---

## Critical section

A critical section is a region of code that must not run on more than one thread at the same time when the threads use the same shared data.

The shared data is the reason for the section. Two threads may execute the same function. They may not both execute the increment of the same counter at the same time.

Properties that a correct lock protocol should give:

1. Mutual exclusion: at most one thread is in the critical section for that data.
2. Progress: if the section is free and threads want to enter, the decision who enters must not wait forever for a thread that is outside the section.
3. Bounded waiting (fairness goal): a thread should not wait forever while others enter repeatedly.

The critical section should be small. Do not do file I/O or a long CPU loop inside the section unless the I/O is the shared object that you protect. A large section reduces parallelism.

Entry and exit must pair. If you forget to unlock, other threads block forever. Use a single exit path when you can. In C, check every early `return`.

Not every line of a program is a critical section. Local variables on a private stack do not need a lock. Immutable data that no thread writes does not need a lock. A race exists only when at least one access is a write and the accesses are concurrent and unsynchronized.

### Questions

#### Theoretical questions

1. What is a critical section?
2. What is mutual exclusion?
3. Why should a critical section stay small?
4. Why do private stack variables not need a lock?
5. What happens if a thread never leaves the critical section?

#### Easy practical tasks

1. Mark the critical section in `n = n + 1` on a shared `n`. Write the three hardware-level steps.
2. Write five sentences about critical sections. Use only facts from this section.
3. Open `man 3 pthread_mutex_lock`. Write the pair of functions that wrap a critical section.
4. List three code regions that should not be inside a lock if you can avoid it (sleep, network wait, large compile).

#### Medium practical tasks

1. Write a program with a shared list. Identify the critical sections for insert and for lookup. Do not implement the list yet if you lack time; write the comments.
2. Draw two threads and a box labeled "critical section". Show a third thread that only reads a thread-local counter outside the box.
3. Take the racy counter from topic 5 and add comments `/* enter */` and `/* exit */` around the increment. Then put a mutex there.

#### Advanced practical tasks

1. Read OSTEP on locks and critical sections. Write the three properties (mutex, progress, bounded waiting) in your own words with one example each.
2. Audit a small open-source C file that uses pthreads. List each critical section and the data that it protects. Write one section that looks too large.

---

## Mutex, semaphore, condition variable

A mutex (mutual exclusion lock) has two operations: lock and unlock. The thread that locks the mutex owns it. Only that thread may unlock it. A mutex is the default tool for a critical section. POSIX: `pthread_mutex_lock`, `pthread_mutex_unlock`. Linux futexes implement the wait (topic 21). You use the POSIX API first.

A semaphore is a counter with wait (P, down) and post (V, up). A binary semaphore can look like a lock, but it does not have an owner in the POSIX `sem_t` model. Any thread may post. A counting semaphore can allow N threads into a pool (for example, N buffer slots). POSIX: `sem_wait`, `sem_post`.

A condition variable lets a thread wait for a boolean condition. The thread must hold a mutex when it waits. `pthread_cond_wait` unlocks the mutex, waits, and locks the mutex again before it returns. Another thread changes the condition, then calls `pthread_cond_signal` or `pthread_cond_broadcast` while it holds the mutex (or with documented rules). Always wait in a loop: `while (!ready) pthread_cond_wait(...)`. A wakeup can be spurious, and another thread can consume the condition.

Typical roles:

- mutex: protect a data structure
- semaphore: count permits or items
- condition variable: wait until a predicate is true (queue not empty, a phase finished)

A mutex is not recursive by default. A thread that locks the same non-recursive mutex twice deadlocks itself. A recursive mutex exists. Prefer a single lock owner path.

Do not use a condition variable without a mutex. Do not use `sleep` to wait for another thread.

### Questions

#### Theoretical questions

1. What is a mutex owner?
2. How does a counting semaphore differ from a mutex?
3. Why must `pthread_cond_wait` run while the mutex is locked?
4. Why do you wait on a condition variable in a `while` loop?
5. Why can a thread deadlock itself with a non-recursive mutex?

#### Easy practical tasks

1. Make a three-column table: mutex, semaphore, condition variable. Add "main operations" and "typical use".
2. Open `man 3 pthread_mutex_lock`, `man 3 sem_wait`, and `man 3 pthread_cond_wait`. Write one sentence for each.
3. Write a two-thread program: thread A locks, increments, unlocks; thread B does the same. Join both.
4. Write four sentences that contrast a mutex and a semaphore. Use only facts from this section.

#### Medium practical tasks

1. Implement a bounded buffer (size 4) with a mutex, a `not_empty` condition, and a `not_full` condition. One producer, one consumer, 20 items.
2. Implement the same bound with two semaphores (empty slots and full slots) plus a mutex for the buffer index. Write which tool counts and which tool protects the index.
3. Show a bug: `if (!ready) pthread_cond_wait` without a loop. Explain in comments why a loop is required even if your test is lucky.

#### Advanced practical tasks

1. Add `pthread_cond_broadcast` when you have two consumers. Document why `signal` can stall if two waiters need different predicates.
2. Read `man 7 pthreads` on cancellation and mutexes. Write a one-page note on unlock-on-all-paths and cleanup handlers.

---

## Deadlock: four Coffman conditions

Deadlock is a state where a set of threads wait forever. Each thread holds a resource and waits for a resource that another thread in the set holds.

Coffman, Elphick, and Shoshani named four conditions. All four must hold for this kind of deadlock:

1. Mutual exclusion: at least one resource is not shareable. A mutex is such a resource.
2. Hold and wait: a thread holds a resource and waits for another resource.
3. No preemption: the kernel or a peer cannot take the lock away. The owner must unlock.
4. Circular wait: a cycle exists. A waits for B's lock, B waits for A's lock.

If you break any one condition, that deadlock cannot form.

Common break methods:

- lock ordering: all threads lock mutex A then mutex B, never B then A (breaks circular wait)
- try-lock and back off (can introduce livelock if you are careless)
- do not wait while you hold a lock (breaks hold and wait; not always possible)
- lock-free designs (advanced)

Two mutexes and two threads are enough:

```text
Thread 1: lock(A); lock(B);
Thread 2: lock(B); lock(A);
```

Deadlock is not a race on a counter. The program stops making progress. CPU can be idle.

Detection: a watchdog, a timeout, or a graph of waiters. Linux does not magically break user mutex deadlocks. Your process hangs.

### Questions

#### Theoretical questions

1. What is deadlock?
2. Name the four Coffman conditions.
3. How does global lock ordering break circular wait?
4. Why does "no preemption" apply to a mutex?
5. How does deadlock differ from a lost-update race?

#### Easy practical tasks

1. Write the four conditions as a list. Add one OS example for each.
2. Draw the two-thread, two-lock cycle.
3. Open `man 3 pthread_mutex_trylock`. Write how a try-lock differs from a lock.
4. Write five sentences about deadlock. Use only facts from this section.

#### Medium practical tasks

1. Write a C program that deadlocks on purpose with two mutexes (use a VM). Show that `top` reports the process as sleeping. Kill it.
2. Fix the program with a single lock order. Show that it finishes.
3. Make a table: "Condition" and "How our lock-order fix affects it".

#### Advanced practical tasks

1. Read about banker's algorithm at a high level. Write why most pthread programs use lock order instead of a banker.
2. Build a three-lock cycle (A-B-C). Then apply a total order A < B < C. Document the acquire sequences that remain legal.

---

## Livelock and starvation

Livelock is a state where threads keep changing state but do not make useful progress. They are not stuck on a lock wait in a cycle. They are busy. Example: two threads detect a conflict, both back off, both retry at the same time, forever. Example: two people step aside in the same direction in a corridor.

Starvation is a state where one thread never gets a resource that others keep receiving. The system as a whole can still make progress. Example: SJF always runs short jobs; a long job waits forever. Example: a mutex wake-up policy always prefers the last waiter.

Differences:

- deadlock: no progress, often sleeping in a cycle
- livelock: activity, no useful progress
- starvation: some threads progress, one thread does not
- race: wrong result, the program usually continues

Fixes for livelock: random back-off, a retry limit, a total order so that one side always wins.

Fixes for starvation: fair locks, aging (priority rises with wait time), CFS-like fairness at the scheduler, do not use strict SJF for scheduling.

A spinlock that always fails and retries without a fair queue can starve a thread on a busy lock. A `SCHED_FIFO` thread can starve `SCHED_OTHER` threads (topic 6).

### Questions

#### Theoretical questions

1. What is livelock?
2. What is starvation?
3. How does livelock differ from deadlock?
4. How can random back-off help livelock?
5. How can strict SJF starve a job?

#### Easy practical tasks

1. Make a three-column table: deadlock, livelock, starvation. Add "CPU use" and "progress".
2. Write a corridor analogy in four sentences for livelock, then write why an OS example is better for a report.
3. Open `man 7 sched` and write one sentence about realtime starvation of normal tasks.
4. Write five sentences that contrast the three failure modes. Use only facts from this section.

#### Medium practical tasks

1. Write two threads that use `trylock` on two mutexes and retry immediately on failure with no sleep. Argue whether this can livelock. Add a random `nanosleep` and a retry cap.
2. Draw a timeline where thread C never runs because A and B always get the mutex first (starvation).
3. Read about lock convoy at a high level. Write how it differs from starvation of one thread.

#### Advanced practical tasks

1. Implement a ticket lock or use `PTHREAD_PRIO_INHERIT` documentation to write how fairness or priority can reduce starvation (preview of the next section).
2. Find a public post-mortem of a livelock (database or kernel). Write five sentences: what retried, and what broke the symmetry.

---

## Priority inversion (high-level)

Priority inversion occurs when a high-priority thread waits for a lock that a low-priority thread holds, and a medium-priority thread runs instead of the low-priority thread.

Sequence:

1. Low L locks mutex M.
2. High H needs M and blocks.
3. Medium M-pri becomes ready and preempts L (L is not done with the critical section).
4. H waits as long as M-pri runs. H is inverted: it waits for a lower-priority thread that cannot run.

The scheduler does the correct thing with priorities, but the lock created a hidden dependency.

Priority inheritance is a common fix. While H waits for M, L temporarily runs at H's priority. Medium cannot preempt L. L unlocks. H runs. L returns to its old priority.

Priority ceiling is another fix. The mutex has a ceiling priority. A thread that locks it rises to that ceiling. This method needs static analysis of priorities.

Linux pthread mutexes can use `PTHREAD_PRIO_INHERIT` when the thread uses a realtime policy. The default desktop `SCHED_OTHER` does not expose the same inversion story as a hard realtime system, but inversion can still delay a thread.

Do not give a long critical section to a low-priority thread if a high-priority thread needs the same lock. Keep the section short. That reduces inversion time even without inheritance.

### Questions

#### Theoretical questions

1. What is priority inversion?
2. What roles do L, H, and the medium thread play?
3. What is priority inheritance?
4. What is a priority ceiling at a high level?
5. Why does a short critical section reduce inversion harm?

#### Easy practical tasks

1. Draw the four-step inversion sequence.
2. Write five sentences about inversion. Use only facts from this section.
3. Open `man 3 pthread_mutexattr_setprotocol`. Write the name of the inheritance protocol if the page lists it.
4. Make a table: "Without inheritance" and "With inheritance". Add two rows.

#### Medium practical tasks

1. Write a storyboard for a robot: low-priority logging holds a lock, high-priority control needs the lock, medium-priority UI runs. Propose a fix.
2. Read a short Mars Pathfinder inversion summary (public articles). Write the bug and the inheritance fix in six sentences.
3. Explain why CFS nice values make inversion less famous on a desktop than on an RTOS, but still possible with locks.

#### Advanced practical tasks

1. If you have rights and a VM, try three threads with `SCHED_FIFO` priorities and a mutex without inheritance, then with `PTHREAD_PRIO_INHERIT`. Measure delay of the high-priority thread. Document every command. Do not run this on a shared host.
2. Compare inheritance and ceiling in a one-page table: implementation cost, need for static ceilings, and inversion bound.

---

## Memory visibility vs locks (preview of memory model)

A lock does more than mutual exclusion. A correct mutex also orders memory.

Without a lock or an atomic with the right order, a write in thread A may not be visible to thread B yet. The compiler can keep a value in a register. The CPU can reorder stores and loads. Another core can see an old cache line.

Example of a broken flag:

```c
/* thread A */
data = 42;
ready = 1;

/* thread B */
while (ready == 0) { }
use(data);   /* may see data != 42 */
```

On a real machine this can fail. A mutex around both sides, or an atomic store-release of `ready` and an atomic load-acquire in B, makes the write of `data` visible before `ready` is seen as 1.

Rules for this path:

1. If you share a non-atomic object, use a mutex (or another documented primitive) for all reads and writes of that object.
2. Do not use a plain `int` flag to publish a struct without a memory-order story.
3. `volatile` in C is not a mutex. `volatile` does not make an increment atomic. Do not use `volatile` as a lock.

Topic 21 returns to barriers and lock-free progress. This section only states: locks exclude and they publish.

C11 atomics (`<stdatomic.h>`) give another path. Use them when you measured a need. For a beginner, a mutex is the default.

### Questions

#### Theoretical questions

1. What does visibility mean for a write in thread A and a read in thread B?
2. Why can the broken flag example fail?
3. What two jobs does a mutex do?
4. Why is `volatile` not a substitute for a mutex?
5. What is the beginner rule for shared non-atomic objects?

#### Easy practical tasks

1. Copy the broken flag example and add comments that mark the missing order.
2. Write five sentences about visibility. Use only facts from this section.
3. Open `man 3 pthread_mutex_lock` and search for "synchronize" or read a POSIX memory-sync note. Write one sentence.
4. Make a table: "Tool", "Exclusion?", "Publish writes?". Rows: mutex, `volatile int`, C11 atomic release/acquire.

#### Medium practical tasks

1. Write the flag example with a mutex around `data` and `ready` on both sides. Explain why B can then use `data`.
2. Write the same example with `atomic_store_explicit` and `atomic_load_explicit` if you know C11. If not, write the intended order in comments.
3. Draw two cores, two caches, and a stale line for `data`. Show a lock as a barrier.

#### Advanced practical tasks

1. Compile the broken flag with `-O2` and run it in a loop on a multicore VM. Record whether you see a failure. If you never see it, write why tests are not proofs.
2. Read a short C11 memory-model intro (cppreference or a textbook). Write acquire, release, and sequential consistency in one sentence each.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. You have a bounded queue used by four threads. Which tools do you pick (mutex, semaphore, condition variable), and which Coffman condition do you watch when you add a second lock?
2. How can a program have no deadlock and still fail because of visibility?
3. When is starvation a scheduler story, and when is it a lock-wake story?
4. How does priority inversion combine a mutex with preemptive scheduling?
5. Why does "one big mutex for the whole process" remove many deadlocks but destroy the multicore benefit of threads?

#### Easy practical tasks

1. Write a one-page cheat sheet: critical section, mutex, semaphore, condvar, four Coffman conditions, livelock, starvation, inversion, visibility.
2. Draw one figure that contains a critical section, a deadlock cycle, and a broken flag.
3. Compile a working producer–consumer from this topic. Run it under `time`.
4. Open POSIX `pthread.h` related man pages and list the five function names that you used.

#### Medium practical tasks

1. Break your producer–consumer on purpose (forget a signal, invert lock order, or drop the `while` loop). Record the symptom. Restore the fix.
2. Write a document that maps each bug class in this topic to one tool or one rule that prevents it.
3. Use ThreadSanitizer on a program that forgets a mutex. Write the report in your own words.

#### Advanced practical tasks

1. Implement dining philosophers with five mutexes. Show deadlock, then fix with a global order or a waiter. Document the Coffman condition that you broke.
2. Read the OSTEP chapters on locks, CV, and deadlock. Write a one-page map from each chapter to a section in this handbook.
