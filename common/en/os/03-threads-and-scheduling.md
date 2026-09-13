# 3. Threads and Scheduling

## Description

A thread is an execution context inside a process. The scheduler chooses which ready thread runs on a CPU core. This topic explains user threads and kernel threads, shared address space, and thread-local storage. You also learn CPU-bound and I/O-bound work, preemption, classic algorithms, the idea of the Linux Completely Fair Scheduler, context switches, affinity, and multicore.

Complete this topic after processes. Complete this topic before synchronization. Locks assume more than one thread.

Use one term for each concept. A kernel thread is a thread that the kernel scheduler sees. A user thread is a thread that a user-space library schedules. A time slice is the interval that a thread may run before the scheduler can preempt it. A context switch is the kernel work that saves one thread and loads another. Affinity is a rule that binds a thread to a set of cores. Do not mix a thread with a process. Do not mix cooperative scheduling with preemptive scheduling.

---

## User threads vs kernel threads

The kernel schedules kernel threads. Each kernel thread has a kernel scheduling identity. On Linux, that identity is a task with its own thread ID (TID). The call `clone` with a shared address space creates a kernel thread in the same process.

A user thread is a context that a library switches in user space. The kernel sees one kernel thread. The library may run many user threads on that one kernel thread. This model is N:1 (many user threads, one kernel thread).

A 1:1 model maps each user-visible thread to one kernel thread. POSIX threads (`pthreads`) on modern Linux use 1:1. `pthread_create` becomes a `clone` that shares memory.

An M:N model maps many user threads onto a pool of kernel threads. Some language runtimes use M:N. The runtime scheduler sits on top of the kernel scheduler.

Properties:

- 1:1: a blocking system call blocks only that kernel thread. Other kernel threads in the process can run. The kernel can place threads on different cores.
- N:1: a blocking system call can block all user threads that share the one kernel thread. The kernel cannot run two user threads of that process on two cores.
- M:N: more complex. The runtime must handle blocking and preemption.

Linux also has kernel-only tasks (workqueues, kthreads). Those are not your program `pthread` threads. They run kernel functions. You see names such as `kworker` in `ps`.

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

1. Write a C program with `pthread_create` that starts two threads. Each thread prints its TID with `gettid()` or `syscall(SYS_gettid)` and the PID with `getpid()`. Show that PIDs match and TIDs differ. This handbook does not contain the source.
2. Draw N:1 and 1:1 diagrams. Label the kernel scheduler.
3. Run `cat /proc/self/status | grep -E 'Threads|Pid|Tgid'`. Write what `Tgid` means (thread group ID, the process PID).

#### Advanced practical tasks

1. Compare Go goroutines or Java virtual threads (documentation only) with pthreads. Write which model each one uses and what happens on a blocking OS call.
2. Read `man 2 clone`. Write which flags `pthread_create` needs so that threads share the address space (`CLONE_VM` and related flags).

---

## Shared address space and thread-local storage

Threads of one process share the same virtual address space. They share:

- the program text
- global and static variables
- the heap
- open file descriptors (the file descriptor table is shared)
- the current working directory and credentials (unless a thread changes them with special calls)

Each thread has a private stack. Local variables in a thread function live on that stack. Each thread has its own registers and its own program counter when it runs.

If thread A writes a global variable, thread B can read the new value. Memory visibility rules apply (topic 4). If thread A writes through a heap pointer that thread B also holds, both threads see the same object.

This share is the reason threads are useful. It is also the reason races exist. Processes do not share an address space by default. Processes use pipes, sockets, or shared memory to talk. Threads already share memory.

Thread-local storage (TLS) is data that belongs to one thread. Other threads do not see that copy. In C11 you use `_Thread_local`. In pthreads you use `pthread_key_create` or compiler TLS. The C library uses TLS for `errno` on many systems. Each thread has its own `errno`.

`fork` in a multithreaded process copies only the calling thread into the child. POSIX lists many restrictions. Avoid `fork` plus threads except for the `fork` then `exec` pattern.

A stack overflow in one thread can hit another mapping. The kernel can place a guard page at the end of a thread stack. Topic 5 covers guard pages.

### Questions

#### Theoretical questions

1. What memory do threads of one process share?
2. What stack does each thread have?
3. What is thread-local storage?
4. Why does `errno` often live in TLS?
5. Why is `fork` in a multithreaded process dangerous if you do not `exec`?

#### Easy practical tasks

1. Write a two-column table: "Shared" and "Private per thread". Add five rows.
2. Open `man 7 pthreads` and find the list of shared attributes. Write four shared items.
3. Draw one address space with text, heap, two stacks, and one TLS block per thread.
4. Write four sentences about shared address space and TLS. Use only facts from this section.

#### Medium practical tasks

1. Write a program with a global counter and a `_Thread_local` counter. Two threads increment both. Print the three values after join. This handbook does not contain the source.
2. Open `man 3 pthread_key_create`. Write one use of a destructor for TLS.
3. Read `man 2 fork` on multithreaded processes. Write three restrictions in your own words.

#### Advanced practical tasks

1. Use `pthread_setspecific` and `pthread_getspecific` for a per-thread name string. Print the name from two threads.
2. Write a one-page note: how a race on a global differs from a safe increment of a TLS counter.

---

## CPU-bound vs I/O-bound

A CPU-bound thread spends most of its time in the running or ready state. It wants the ALU and the caches. Examples: compression, compilation of a large file, a tight numeric loop.

An I/O-bound thread spends most of its time blocked. It waits for a disk, a socket, a user, or a lock. Examples: `read` from a slow pipe, a server that waits for clients, a program that waits for a key.

The scheduler must treat them differently if it wants a responsive system.

- A CPU-bound thread can fill a core for a long time. If the scheduler never preempts it, interactive threads wait.
- An I/O-bound thread often runs for a short burst, then blocks again. If the scheduler runs it soon after it wakes, the user or the remote client sees low latency.

A real process can mix both. A program reads a file (I/O), then parses it (CPU), then writes a result (I/O). The state changes over time.

`top` and `time` help you classify. High `%CPU` with little `wa` (I/O wait) suggests CPU-bound. High I/O wait or a process that is often in `S` suggests I/O-bound.

The scheduler does not need your label. It observes bursts and sleep. You still need the labels to design programs. Do not hold a lock while you do long CPU work if other threads must do I/O.

### Questions

#### Theoretical questions

1. What is a CPU-bound thread?
2. What is an I/O-bound thread?
3. Why does an interactive system preempt long CPU-bound threads?
4. Why must an I/O-bound thread run soon after it wakes?
5. Can one process be CPU-bound in one phase and I/O-bound in another? Explain.

#### Easy practical tasks

1. Write a two-column table: "CPU-bound" and "I/O-bound". Add four example programs.
2. Run `time sleep 2`. Write why user plus sys time is much smaller than real time.
3. Run `time dd if=/dev/zero of=/dev/null bs=1M count=4096`. Write whether this looks CPU-bound or I/O-bound on your machine.
4. Open `man 1 time`. Write the meaning of real, user, and sys.

#### Medium practical tasks

1. Write a C loop that adds numbers for a few seconds. Run `time` and `top` while it runs. Classify it.
2. Write a program that `read`s from a pipe with no writer for five seconds (or use `sleep`). Classify it with `ps` STAT.
3. Draw a timeline for a web worker: wait for socket, parse request (CPU), read a file (I/O), write response.

#### Advanced practical tasks

1. Use `pidstat` or `top -H` on a compiler if you have one. Write which threads look CPU-bound.
2. Read the OSTEP scheduling chapter on workload. Write a one-page summary of turnaround time versus response time for the two workloads.

---

## Preemptive vs cooperative scheduling

A preemptive scheduler can stop a running thread without the thread permission. The timer interrupt is the usual trigger. The kernel also preempts when a higher-priority thread becomes ready, on systems that use priority preemption.

A cooperative scheduler switches only when the thread yields. The thread calls a yield function, or it blocks on I/O. If the thread never yields and never blocks, it keeps the core. A bug in one thread can freeze the system.

General-purpose operating systems use preemption for user threads. Linux preempts user tasks. Linux also has kernel preemption options so that a thread in a kernel path can be preempted in more places.

Cooperative switching still exists:

- user-space fibers and some language runtimes yield in the library
- an old embedded loop can run without a timer
- a thread that disables preemption in the kernel (a spinlock section) is briefly cooperative on that core

Preemption gives fairness and latency. It has a cost: more context switches and the need for synchronization. Cooperative scheduling is simpler for some small systems. It is not enough for a multiuser Linux desktop.

A program cannot disable the Linux user-preemption timer for other processes. A process can waste CPU in a loop. The scheduler will still share the core with others, unless the process uses a realtime policy with care.

### Questions

#### Theoretical questions

1. What is preemptive scheduling?
2. What is cooperative scheduling?
3. What hardware event often causes preemption?
4. Why can a cooperative system freeze?
5. Where does cooperative switching still appear on a modern Linux machine?

#### Easy practical tasks

1. Write five sentences that contrast preemption and cooperation. Use only facts from this section.
2. Make a table: "Who decides to switch" for preemptive and cooperative.
3. Open `man 2 sched_yield`. Write when a thread would call it.
4. List two events that can preempt a Linux user thread. Then list one event that is a voluntary block, not preemption. Label each event.

#### Medium practical tasks

1. Write a CPU loop and a second process that prints the time every 200 ms. Confirm that the printer still runs. Explain how preemption makes that possible.
2. Draw two timelines: a cooperative hog that never yields, and a preemptive time slice that interrupts the hog.
3. Read `man 7 sched`. Write one sentence about the default Linux policy `SCHED_OTHER` and preemption.

#### Advanced practical tasks

1. Read a short note on kernel preemption in Linux. Write how a spinlock section is briefly non-preemptible on that core.
2. Write a one-page report: why a language runtime with cooperative user threads still needs kernel preemption of its worker threads.

---

## FIFO, Round Robin, priority, CFS (idea)

A scheduling algorithm chooses the next ready thread.

First-in first-out (FIFO), also called First-Come First-Served: the thread that became ready first runs until it blocks or exits. There is no time slice in the simple form. A long CPU-bound thread delays everyone behind it. Linux realtime policy `SCHED_FIFO` is a related idea for privileged tasks. Do not use realtime policies in daily labs without a plan to recover the machine.

Round Robin: each thread gets a time slice. When the slice ends, the thread goes to the tail of the ready queue. Short slices improve response time. Short slices increase context-switch cost.

Priority: each thread has a rank. The scheduler picks a higher-priority ready thread first. Static priority is a fixed number. Dynamic priority can change with sleep and CPU use. Linux `nice` is a user hint for the default policy. A lower nice value means a higher share (more favorable). Only privileged users can set a strongly negative nice.

The Linux Completely Fair Scheduler (CFS) is the idea behind the default policy for normal tasks. CFS does not use a simple global Round Robin queue only. It tracks virtual runtime. A thread that ran less gets a turn. Sleeping threads get a boost when they wake so that interactive work feels fast. The exact data structure (a tree of tasks) is an implementation detail. You need the idea: fairness over time, not a fixed FIFO stamp.

No single algorithm is best for all goals. Turnaround time, response time, fairness, and throughput conflict. A desktop wants response time. A batch machine wants throughput.

### Questions

#### Theoretical questions

1. How does FIFO treat a long CPU-bound thread that is first in line?
2. What does Round Robin do when a time slice ends?
3. What does a lower `nice` value mean on Linux?
4. What is the CFS idea of virtual runtime?
5. Why do response time and throughput conflict?

#### Easy practical tasks

1. Write a four-row table: FIFO, Round Robin, priority, CFS. Add one sentence per row.
2. Open `man 1 nice` and `man 2 nice`. Write the nice range.
3. Run `ps -o pid,ni,pri,cmd`. Write the nice value of your shell.
4. Write five sentences about these algorithms. Use only facts from this section.

#### Medium practical tasks

1. Run `nice -n 19 yes > /dev/null` and an un-niced `yes > /dev/null` for a few seconds. Use `top` to compare `%CPU`. Then stop both (`killall yes` only on your machine).
2. Draw a Round Robin timeline for three threads and a slice of one unit.
3. Read `man 7 sched` on `SCHED_OTHER`, `SCHED_FIFO`, and `SCHED_RR`. Write one sentence for each.

#### Advanced practical tasks

1. Read the OSTEP chapter on scheduling. Write a one-page comparison of SJF (preview) and Round Robin for response time.
2. Read a CFS overview (kernel documentation or a textbook). Write how vruntime relates to sleep. Do not claim you read the full kernel source.

---

## Context switch, affinity, and multicore

A context switch is the kernel path that stops one thread and starts another on the same core. The kernel saves registers to the PCB (topic 2). The kernel loads the next thread registers. If the next thread is in another process, the kernel also switches the address space (page tables and TLB effects). Topic 5 covers the TLB.

A context switch has a cost: save and restore, cache and TLB warm-up, and kernel time. Too many switches waste CPU. Too few switches hurt latency. The time slice is a trade.

Affinity is a rule that binds a thread to a set of cores. Linux has `sched_setaffinity` and the `taskset` command. Soft affinity is the scheduler habit to keep a thread on the last core so that the cache stays warm. Hard affinity is a mask that you set.

Multicore means more than one core can run threads at the same time. True parallelism needs more than one kernel thread (or more than one process). Two user threads on one kernel thread do not use two cores at once.

Cache coherence makes memory look consistent across cores, with rules. Two threads that share a counter still need a lock or an atomic (topic 4). Multicore does not remove races.

A thread that moves across cores can lose cache warmth. A thread that is pinned to a busy core can wait while another core is idle. Pin only when you measure a need.

`/proc/cpuinfo` and `nproc` show how many logical CPUs you have. Hyper-threading can show more logical CPUs than physical cores.

### Questions

#### Theoretical questions

1. What does the kernel save in a context switch?
2. Why is a switch between processes often more expensive than a switch between threads of one process?
3. What is CPU affinity?
4. Why do two user threads on one kernel thread not run on two cores at once?
5. Does multicore remove the need for locks on shared data? Explain.

#### Easy practical tasks

1. Run `nproc` and `lscpu`. Write how many logical CPUs you see.
2. Open `man 1 taskset`. Write one example command that prints a PID affinity.
3. Write five sentences about context switch, affinity, and multicore. Use only facts from this section.
4. Run `taskset -p $$` if the tool exists. Write the affinity mask of your shell.

#### Medium practical tasks

1. Start two CPU loops. Use `top` and watch whether both get CPU on a multicore machine.
2. Draw one core: thread A, switch, thread B. Label saved registers.
3. Pin a CPU loop with `taskset -c 0` and leave a second loop free. Write what `top` shows for each.

#### Advanced practical tasks

1. Read `man 2 sched_setaffinity`. Write a small program that pins itself to core 0 and prints `sched_getcpu` if available.
2. Write a one-page note: cache-line bouncing when two cores increment one shared counter. Topic 4 will add the lock.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the 1:1 thread model and multicore scheduling work together?
2. Why is a race possible on a global and not possible on a TLS integer that only one thread writes?
3. A desktop mixes a compiler and a text editor. Which scheduling ideas keep the editor responsive?
4. Why does CFS still need preemption even if it aims at fairness?
5. When do you pin a thread, and when do you let the scheduler pick a core?

#### Easy practical tasks

1. Write a cheat sheet: kernel thread, user thread, TLS, CPU-bound, I/O-bound, preemption, FIFO, Round Robin, nice, CFS, context switch, affinity.
2. Run `ps -eL | wc -l` and `nproc`. Write one sentence that relates thread count to core count.
3. Draw one process with two stacks and two cores. Show only one thread running per core.
4. Open `man 7 pthreads` and `man 7 sched`. Write one fact from each page.

#### Medium practical tasks

1. Write a two-thread program that burns CPU. Run `top -H`. Write two TIDs. This handbook does not contain the source.
2. Use `time` on a CPU loop and on `sleep 1`. Classify both. Then write which scheduler behavior each one needs.
3. Measure `sched_yield` in a tight loop with `strace -c` or `perf stat` if you have it. Write that yield is a system call.

#### Advanced practical tasks

1. Read the OSTEP concurrency and scheduling chapters that match this topic. Write a one-page map from textbook terms to Linux `ps` and `nice`.
2. Design (no full code) an M:N runtime: when a worker blocks on `read`, what must the runtime do so that other user threads continue?
