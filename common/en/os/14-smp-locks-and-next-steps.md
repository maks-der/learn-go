# 14. SMP Locks and Next Steps

## Description

Topic 4 introduced mutexes, deadlock, and a first view of visibility. This topic is the next layer on a multiprocessor (SMP) Linux system, plus the practice that closes the path. You learn lock ordering, per-CPU data, futexes, the idea of RCU, and memory barriers at a high level. You then plan a tiny shell, a data race that you hit and fix, and a reading path through OSTEP and xv6.

Complete this topic after observability. You already know critical sections, Coffman conditions, and `fork` plus `exec`.

Use one term for each concept. Lock ordering is a global rank of locks that every thread must follow. Per-CPU data is a replica of a variable for each CPU. A futex is a kernel wait queue keyed by a user address. RCU is a Linux read-copy-update publish method for readers that do not take a writer mutex. A memory barrier is a constraint on the order of memory operations. A tiny shell is a loop that reads a line, splits it, and runs a program with `fork`, `exec`, and `wait`. Do not mix a barrier with a mutex. Do not treat this file as a solution key.

---

## Lock ordering and per-CPU data

Deadlock among locks needs a cycle (topic 4). A simple prevention rule is a total order. Every thread that needs more than one lock acquires them in the same rank order. The cycle cannot form.

How to apply the rule:

1. List every lock in the subsystem.
2. Give each lock a rank (a number or a documented name order).
3. Acquire from low rank to high rank (pick one direction and keep it).
4. If you hold a high lock and need a low lock, you must drop and retry, or you must redesign.

The kernel documents lock classes and uses lockdep. Lockdep is a runtime checker. It records the acquire order and reports a possible cycle. You can enable it in a debug kernel. User programs do not have lockdep by default. You write the order in comments and in review.

Ordering is per lock object, not only per type. Two mutexes of the same C type still need a rule (for example, order by address, or order by inode number).

A single global mutex avoids cycles and kills scalability. Ordering is the middle path: many locks, one rank story.

Priority inversion (topic 4) is a different bug. Ordering does not fix inversion.

Per-CPU data is a copy of a statistic or a cache for each CPU. A thread that runs on CPU 2 updates only the CPU 2 slot. No lock is needed for that update if no other CPU writes that slot. The kernel uses `DEFINE_PER_CPU` and similar macros. User space can approximate this with an array indexed by `sched_getcpu`, with care: the thread can migrate between the load of the index and the store. A seqcount or a migrate-disable region (kernel) fixes that. For a lab, disable affinity change or use atomics.

Per-CPU counters still need a sum that is consistent if you read all slots. Readers can see a slightly stale total. That is often enough for stats.

Do not take locks in a signal handler. Do not invent a second order "only when this flag is set" unless you write a proof.

### Questions

#### Theoretical questions

1. How does a total lock order prevent a deadlock cycle?
2. What must you do if you hold a high-rank lock and need a low-rank lock?
3. What does lockdep check?
4. What is per-CPU data?
5. Why can a user-space `sched_getcpu` index race with migration?

#### Easy practical tasks

1. Write five sentences about lock ordering and per-CPU data. Use only facts from this section.
2. Draw two threads and two locks with a good order and with a cycle.
3. Open `man 3 pthread_mutex_lock`. Write that POSIX does not rank your mutexes for you.
4. Make a table: one big mutex versus ordered fine locks versus per-CPU counters. Add deadlock risk and parallelism.

#### Medium practical tasks

1. Take a two-lock homework (account A and account B). Write the rank rule (order by account id). Write the transfer steps.
2. Read a short lockdep overview in kernel documentation. Write six sentences on what a report means.
3. Write a per-CPU style counter array of size `nproc`. Two threads increment "their" slot. Print the sum. Document the migration caveat.

#### Advanced practical tasks

1. Implement dining philosophers with a global lock order (topic 4) and an explicit rank comment on every acquire. This handbook does not contain the source.
2. Write a one-page lock-order document for a toy graph: inode lock, then page lock. Invent ranks. Do not claim it matches Linux.

---

## Futex, RCU (idea), memory barriers (high-level)

A futex (fast userspace mutex) is a kernel wait queue that is keyed by a user-space address. The fast path of `pthread_mutex_lock` is an atomic in user space. Only a contended lock enters the kernel with `futex`. `strace -e futex` on a contended mutex shows that path (topic 4 review). You do not call `futex` in a first program. You call pthreads. Know that the wait is not a spin in the kernel for every lock.

Read-copy-update (RCU) is a Linux synchronization method. Readers do not take a shared mutex. A reader enters an RCU read-side critical section (`rcu_read_lock` in the kernel). The reader then follows pointers that the writer published.

A writer that must change a structure:

1. Copies the structure (or allocates a new one).
2. Updates the copy.
3. Publishes the new pointer with an ordered store.
4. Waits until every CPU (or every reader) has left the old read-side sections (a grace period).
5. Frees the old structure.

The grace period is the heart of RCU. The writer does not free memory that a reader can still see. User space has userspace RCU libraries. Most application code must still use a mutex. RCU pays when there are many readers and few writers.

A memory barrier (fence) is a constraint on the order of memory operations. The CPU and the compiler must not reorder some loads and stores across the barrier. Topic 4 said that a mutex already includes the barriers that you need for data that the mutex protects. Explicit barriers (`atomic_thread_fence`, `smp_mb` in the kernel) appear in lock-free code and in device drivers.

High-level pairs:

- acquire: later loads and stores stay after the acquire
- release: earlier loads and stores stay before the release
- full barrier: both directions

Do not sprinkle barriers "to be safe" around a mutex. Do not call a program lock-free only because you see no `pthread_mutex_t` in the source. Lock-free is a progress guarantee. It is hard to get right.

This section is high-level. You do not implement RCU or a lock-free list here.

### Questions

#### Theoretical questions

1. What is a futex?
2. Why does an uncontended mutex not need a system call?
3. What is an RCU grace period?
4. What does a release store publish to an acquire load?
5. Why is "no mutex in the source" not the same as lock-free?

#### Easy practical tasks

1. Open `man 2 futex` and a short RCU overview (kernel docs). Write one sentence for each.
2. Write five sentences about futex, RCU, and barriers. Use only facts from this section.
3. Draw uncontended lock (atomic only) versus contended lock (futex wait).
4. Make a table: mutex, RCU reader, explicit barrier. Add "typical user".

#### Medium practical tasks

1. Run a two-thread mutex hammer under `strace -e futex`. Write whether you see `FUTEX_WAIT`. This handbook does not contain the source.
2. Draw RCU: readers, old object, new object, grace period, free.
3. Read C11 `memory_order_acquire` and `memory_order_release` notes. Write six sentences that link them to topic 4 visibility.

#### Advanced practical tasks

1. Write a one-page note: when the kernel uses RCU (example: a list of tasks or a routing table) and why a mutex would hurt readers.
2. Read an OSTEP or herdtools memory-model intro at a high level. Write three reorderings that a barrier would stop. No exploit work.

---

## Write a tiny shell (`fork`/`exec`)

A tiny shell is the best single program for topic 2 and topic 6. You reuse one process as the parent. Each command runs in a child.

Minimum behavior:

1. Print a prompt.
2. Read one line from stdin.
3. Split the line into a command and arguments (spaces are enough).
4. `fork`.
5. The child calls `execvp` (or `execve` with a path search that you write later).
6. The parent calls `waitpid` and then prints the prompt again.
7. Exit the loop on end-of-file or on a `exit` built-in.

In scope for a first week: external programs (`ls`, `date`), a working directory that the child inherits, and a non-zero exit status that you print.

Out of scope until the minimum works: pipes, redirects, job control, quotes, and glob. Those features are extra labs. Pipes need topic 9. Redirects need `dup2` and `open`.

Checks that you write yourself:

- `waitpid` reaps the child. `ps` shows no zombie.
- `exec` failure in the child prints an error and `_exit`s. The parent must not `exec` by mistake.
- A built-in `cd` must run in the parent. A `cd` in the child does not change the shell.

Topic 2 already defined `fork`, `exec`, and zombies. This section only states the product.

Do not start from a 2000-line GitHub shell. Write the loop first. Add one feature at a time. This handbook does not contain the source.

### Questions

#### Theoretical questions

1. Why must `cd` run in the parent?
2. What does the child do if `execvp` returns?
3. Why does the parent call `waitpid`?
4. Which topic 2 bug appears if the parent never waits?
5. Why are pipes an extra lab, not the first lab?

#### Easy practical tasks

1. Open `man 2 fork`, `man 3 execvp`, and `man 2 waitpid`. Write one sentence for each.
2. Write a ten-line design: prompt, parse, fork, exec, wait. No code yet.
3. Run `ps -o pid,ppid,stat,cmd` while a real shell runs `sleep 30`. Write the parent and child.
4. List five commands that your tiny shell must be able to run when the minimum is done.

#### Medium practical tasks

1. Write the minimum shell. Test `ls`, `true`, `false`, and a bad name. Record exit statuses.
2. Add a built-in `cd` and `pwd`. Document one failed `cd` (bad path).
3. Use `strace -f -e fork,clone,execve,wait4,waitpid` on your shell for one `ls`. Write the child PID story.

#### Advanced practical tasks

1. Add one extra feature only after the minimum is stable: `|` between two commands, or `>` redirect. Write a design of FDs before you code.
2. Write a one-page test list: zombies, `exec` fail, `cd`, EOF, and a long argument line.

---

## Hit a data race, then fix it

A race that you cannot see does not teach you. Build a program that is wrong on purpose, show a bad result, then add a lock.

A standard lab:

1. Start `N` threads (`N` at least 4).
2. Each thread adds 1 to a shared `int` (or `long`) `ITERS` times (`ITERS` at least 100000).
3. Join all threads.
4. Print the counter. The correct value is `N * ITERS`.
5. Without a lock, the printed value is often smaller. If it is correct, raise `ITERS` or `N`, or run again. A race is not guaranteed on every machine in one run.
6. Add a `pthread_mutex_t` around the increment. Confirm the correct total.
7. Optional: use ThreadSanitizer (`-fsanitize=thread`) on the racy build. Read the report.

The increment is the critical section (topic 4). The mutex gives exclusion and visibility. Do not use `sleep` to "fix" the race. Do not use `volatile` as the fix.

A second step: two locks in opposite order. Hang, then apply lock ordering from the first section.

Observe with topic 13: `time` the racy and the locked versions. The lock is slower. Correctness is the first goal. Then you may split counters (per-CPU idea) if you measure a need.

This handbook does not contain the source. Do not paste a full solution from the internet as your only work.

### Questions

#### Theoretical questions

1. Why can `N * ITERS` fail without a lock?
2. Why might one run still print the correct total?
3. What two jobs does the mutex do in the fix?
4. Why is `volatile` not the fix?
5. How do you use TSan in this lab?

#### Easy practical tasks

1. Write the expected formula `N * ITERS` for `N=4` and `ITERS=100000`.
2. Open `man 3 pthread_create` and `man 3 pthread_mutex_lock`. Write one sentence for each.
3. Write five sentences about the race lab. Use only facts from this section.
4. List three ways to make the race more visible (more iterations, more threads, TSan).

#### Medium practical tasks

1. Implement the racy counter. Save a wrong output. Then add the mutex. Save a correct output.
2. Run `time` on both versions. Write the two elapsed times.
3. Draw load, add, store for two threads without a lock.

#### Advanced practical tasks

1. Replace the mutex with an atomic add (`atomic_fetch_add`). Compare totals and times. Write when a mutex is still required (multi-field updates).
2. Add a second shared structure that needs two locks. Show a deadlock, then fix the order. Document ranks.

---

## Read OSTEP and xv6

[Operating Systems: Three Easy Pieces](https://ostep.org/) (OSTEP) is a free textbook. It matches this path: virtualization (CPU and memory), concurrency, and persistence. Read the chapters next to the topics, not only at the end.

A useful pairing:

- OSTEP CPU virtualization and process chapters: topics 2 and 3
- OSTEP lock and CV chapters: topics 4 and 14
- OSTEP paging and TLB chapters: topic 5
- OSTEP files and journaling chapters: topics 6 and 7
- OSTEP concurrency bugs: this topic's race lab

Do the OSTEP homework questions in your own words. The book has figures. This handbook stays in STE and does not copy those figures.

[xv6](https://pdos.csail.mit.edu/6.1810/2024/xv6/book-riscv-rev4.pdf) is a small teaching Unix. The book and the source fit in a course. You see `fork`, page tables, a scheduler, and a simple filesystem in one tree. Use the RISC-V edition that your course picks. You do not need to replace Linux labs with xv6. You use xv6 to see a whole kernel that you can read.

Also keep:

- [Linux man-pages](https://www.kernel.org/doc/man-pages/)
- [The Linux Programming Interface](https://man7.org/tlpi/) for POSIX depth

Suggested extra practice (no solutions here): `mmap` a file and change it; a ring buffer between two threads; `strace` on `cat`; a pipe between two processes; a process in a new namespace; a toy filesystem in a file (FUSE optional).

Do not treat xv6 as Linux. Line numbers and APIs differ. Do not submit xv6 source as your tiny shell.

### Questions

#### Theoretical questions

1. What three big themes does OSTEP use?
2. How do you pair OSTEP paging chapters with topic 5?
3. What is xv6?
4. Why is xv6 not a substitute for Linux `man` pages?
5. Which resource do you open for a POSIX call signature?

#### Easy practical tasks

1. Open [ostep.org](https://ostep.org/). Write the titles of the first three chapters that you will read.
2. Open the xv6 book PDF. Write the chapter title that matches processes.
3. Write five sentences about OSTEP and xv6. Use only facts from this section.
4. Bookmark man-pages and TLPI. Write one sentence about when you open each.

#### Medium practical tasks

1. Read one OSTEP chapter that matches a topic you already finished. Write a one-page summary in STE.
2. Skim one xv6 source file (`proc.c` or `vm.c` in the edition that you picked). Write six sentences: what the file owns.
3. Make a table: this handbook topic number, OSTEP chapter, xv6 file (guess, then check).

#### Advanced practical tasks

1. Write a semester plan: which OSTEP chapters you read each week next to topics 2 to 14.
2. Run xv6 in an emulator if your course provides instructions. Run `ls` or the equivalent. Write how the environment differs from Linux. Follow only official lab docs.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do lock ordering, futex waits, and per-CPU stats fit one scalable counter service?
2. When do you stay with a mutex, and when do you only read about RCU?
3. How do the tiny shell and the race lab prove topics 2 and 4 on a real machine?
4. Which topic 13 tools do you apply to the shell and to the race program?
5. What do you take from OSTEP and xv6 that this Linux handbook cannot replace?

#### Easy practical tasks

1. Write a cheat sheet: lock rank, lockdep, per-CPU, futex, RCU grace period, acquire/release, tiny shell, racy counter, OSTEP, xv6.
2. Draw a path: topics 1 to 13 as boxes, this topic as labs plus reading.
3. List the five official resources from `os.topics.md` (OSTEP, xv6, man-pages, TLPI, and this path).
4. Open `man 7 pthreads` and the OSTEP concurrency table of contents. Write one fact from each.

#### Medium practical tasks

1. Write a personal lab report template: goal, commands, `strace` or `ps` evidence, result. Use it for the shell or the race. Do not paste full source into the template file.
2. Combine lock ordering with the two-lock account transfer on paper, then implement it after the race lab is correct.
3. Document a reading log: date, OSTEP chapter, three sentences, one Linux command that matches the chapter.

#### Advanced practical tasks

1. After the minimum shell works, add a pipe (topic 9) and trace it. Write a one-page design plus evidence. No copied shell.
2. Write a one-page "what next" plan: more OSTEP projects, an OS course with xv6 labs, or a deeper Linux chapter (cgroups, io_uring). Pick one and give a reason.
