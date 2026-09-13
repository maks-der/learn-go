# 6. Scheduling

## Description

The scheduler chooses which ready thread runs on a CPU core. This topic explains CPU-bound and I/O-bound work, preemption, and classic algorithms: FIFO, SJF, and Round Robin. You learn priority and `nice`, the idea of the Linux Completely Fair Scheduler, the cost of a context switch, and CPU affinity.

Complete this topic after processes and threads. Complete this topic before you study locks in depth. Lock convoys and priority inversion use scheduler facts.

Use one term for each concept. A time slice is the interval that a thread may run before the scheduler can preempt it. A context switch is the kernel work that saves one thread and loads another. Affinity is a rule that binds a thread to a set of cores. Do not mix priority with `nice` without the Linux mapping. Do not mix cooperative scheduling with preemptive scheduling.

---

## CPU-bound vs I/O-bound

A CPU-bound thread spends most of its time in the running or ready state. It wants the ALU and the caches. Examples: compression, compilation of a large file, a tight numeric loop.

An I/O-bound thread spends most of its time blocked. It waits for a disk, a socket, a user, or a lock. Examples: `read` from a slow pipe, a server that waits for clients, a program that waits for a key.

The scheduler must treat them differently if it wants a responsive system.

- A CPU-bound thread can fill a core for a long time. If the scheduler never preempts it, interactive threads wait.
- An I/O-bound thread often runs for a short burst, then blocks again. If the scheduler runs it soon after it wakes, the user or the remote client sees low latency.

A real process can mix both. A program reads a file (I/O), then parses it (CPU), then writes a result (I/O). The state changes over time.

`top` and `time` help you classify. High `%CPU` with little `wa` (I/O wait) suggests CPU-bound. High I/O wait or a process that is often in `S` suggests I/O-bound.

The scheduler does not need your label. It observes bursts and sleep. You still need the labels to design programs: do not hold a lock while you do long CPU work if other threads must do I/O.

### Questions

#### Theoretical questions

1. What is a CPU-bound thread?
2. What is an I/O-bound thread?
3. Why does an interactive system preempt long CPU-bound threads?
4. Why should an I/O-bound thread run soon after it wakes?
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

1. Use `pidstat` or `top -H` on a browser or a compiler if you have one. Write which threads look CPU-bound.
2. Read the OSTEP scheduling chapter on workload. Write a one-page summary of turnaround time versus response time for the two workloads.

---

## Preemptive vs cooperative

A preemptive scheduler can stop a running thread without the thread's permission. The timer interrupt (topic 2) is the usual trigger. The kernel also preempts when a higher-priority thread becomes ready, on systems that use priority preemption.

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

1. Read about `CONFIG_PREEMPT` in Linux kernel documentation. Write how kernel preemption differs from user preemption.
2. Compare a Go scheduler (documentation) with the Linux scheduler: which layer preempts goroutines, and when does the OS still preempt the OS thread?

---

## FIFO, SJF, Round Robin

Textbook algorithms explain the trade-offs. Linux does not run these in their pure form for normal tasks. Learn them so that you can read papers and later understand CFS.

FIFO (First In, First Out), also called FCFS (First Come, First Served): the scheduler runs threads in arrival order. A thread runs until it blocks or finishes. There is no time slice. A long CPU-bound thread that arrives first delays everyone. FIFO is simple. FIFO has poor response time when a long job is first.

SJF (Shortest Job First): the scheduler runs the job with the shortest next CPU burst. SJF reduces average turnaround time. The scheduler must know or guess the burst length. A stream of short jobs can starve a long job. Preemptive SJF is SRTF (Shortest Remaining Time First): a new shorter job can preempt the current job.

Round Robin (RR): the scheduler gives each ready thread a time slice. At the end of the slice, the thread goes to the tail of the ready queue. RR improves response time for interactive work. If the slice is too large, RR looks like FIFO. If the slice is too small, the machine spends too much time in context switches.

Metrics that you use when you compare algorithms:

- turnaround time: finish time minus arrival time
- response time: first run minus arrival time
- waiting time: time in the ready queue
- fairness: whether one thread waits forever

Work through small Gantt charts on paper. Use three jobs with known burst lengths. Do not start with a real Linux trace.

### Questions

#### Theoretical questions

1. How does FIFO choose the next thread?
2. Why can FIFO have a long response time?
3. What information does SJF need?
4. What is the starvation risk of SJF?
5. How does Round Robin use a time slice?

#### Easy practical tasks

1. Draw a Gantt chart for FIFO with jobs A=8, B=4, C=2 that all arrive at time 0 in order A, B, C. Compute turnaround times.
2. Draw SJF for the same jobs. Compute turnaround times.
3. Draw RR with slice 2 for the same jobs. Compute response times.
4. Open any OSTEP scheduling chapter figure and copy one chart by hand. Label the algorithm.

#### Medium practical tasks

1. Add arrival times (A at 0, B at 1, C at 2) and repeat FIFO and SRTF on paper. Write the difference.
2. Explain in six sentences what happens to RR when the slice is 1 versus 100 for three equal jobs of length 30.
3. Make a table: algorithm, needs burst estimate?, can starve?, good for interactive?

#### Advanced practical tasks

1. Write a small simulator in Python or C that reads a list of jobs and prints Gantt charts for FIFO, SJF, and RR. Use only integer times.
2. Read about convoy effect in FIFO. Write a one-page example with one long job and many short jobs.

---

## Priority and nice

A priority is a number that the scheduler uses to prefer one thread over another. Higher priority usually means "run sooner" or "get more CPU". The exact number scale depends on the OS.

On Linux, a normal process has a nice value from `-20` to `19`. A smaller nice value means a higher priority for the default scheduler. Nice `0` is the default. Nice `19` is the least favored. Only privileged users can move a process to a negative nice value.

```text
nice -n 10 ./cpu-loop
renice 15 -p PID
```

Linux also has realtime policies `SCHED_FIFO` and `SCHED_RR` with priorities 1 to 99. A realtime thread can starve normal threads if it does not block. Do not use realtime policies on a shared machine without a plan.

Priority is not a promise of a fixed share. The default Linux scheduler still tries to be fair among threads at similar nice values. Nice biases the share.

Priority inversion is a later topic: a low-priority thread holds a lock that a high-priority thread needs. Topic 7 covers that problem.

User tools: `top` shows `NI` (nice) and `PR` (priority as `top` computes it). `ps` can show `ni`.

### Questions

#### Theoretical questions

1. What is a scheduling priority?
2. What is the Linux nice range for normal processes?
3. Does a larger nice value mean more CPU on Linux? Explain.
4. Why are realtime policies dangerous on a desktop?
5. Who may set a negative nice value?

#### Easy practical tasks

1. Run `ps -o pid,ni,pri,cmd`. Write the nice value of your shell.
2. Run `nice -n 10 sleep 30 &` and `ps -o pid,ni,cmd -p $!`. Write the nice value.
3. Open `man 1 nice` and `man 1 renice`. Write the difference between the two commands.
4. Write four sentences about nice. Use only facts from this section.

#### Medium practical tasks

1. Start two CPU loops. `renice` one of them to 19. Use `top` to see whether the nicer process gets less CPU. Write the observation.
2. Draw a scale from nice `-20` to `19` and mark "more favored" and "less favored".
3. Read `man 7 sched` on `SCHED_FIFO`. Write one difference from `SCHED_OTHER`.

#### Advanced practical tasks

1. Read `man 2 setpriority` and `man 2 sched_setscheduler`. Write which call changes nice and which call changes policy.
2. Write a one-page warning sheet for realtime scheduling: starvation, priority inversion preview, and when audio or robotics stacks use it.

---

## Completely Fair Scheduler idea (Linux)

The Completely Fair Scheduler (CFS) is the default Linux scheduler for `SCHED_OTHER` (also called `SCHED_NORMAL`) for a long period of Linux history. The idea is to give each runnable task a fair share of CPU time.

CFS does not use a simple Round Robin list as the only structure. It tracks virtual runtime (`vruntime`). A task that runs accumulates `vruntime`. A task that waits does not accumulate `vruntime` in the same way. The scheduler picks the runnable task with the smallest `vruntime`. That task is the one that has had the least fair time so far.

Nice values weight the accumulation. A nicer (less favored) task's `vruntime` grows faster, so it is picked less often.

CFS uses a red-black tree of runnable tasks ordered by `vruntime`. The pick is the leftmost node. You do not implement that tree in this topic. You only need the idea: fair share, not a fixed FIFO queue.

Linux continues to change the scheduler. Newer kernels add or replace parts (for example EEVDF in recent kernels). The learning point stays: Linux aims for fairness and good latency for desktop and server loads, not for a textbook SJF.

You can see scheduler statistics in `/proc/<pid>/sched` and in `perf sched` (topic 20). Beginners can stay with `top` and `nice`.

### Questions

#### Theoretical questions

1. What problem does CFS try to solve?
2. What is `vruntime` in one sentence?
3. How does CFS choose the next task at a high level?
4. How does nice affect CFS?
5. Why is CFS not the same as textbook Round Robin?

#### Easy practical tasks

1. Write five sentences that describe CFS. Use only facts from this section.
2. Open `man 7 sched` and find `SCHED_OTHER`. Write that it is the default policy.
3. Run `cat /proc/self/sched | head` if the file exists. Write two field names that you see.
4. Make a table: "Textbook RR" vs "CFS idea". Add three rows.

#### Medium practical tasks

1. Start two equal CPU loops at nice 0. Confirm that `top` shows about half the CPU each on a single core (`taskset -c 0` if you have `taskset`).
2. Draw a number line of `vruntime` for two tasks. Show why the smaller `vruntime` is picked.
3. Read a kernel documentation page on CFS or EEVDF. Write one sentence about what your kernel version uses if the page says.

#### Advanced practical tasks

1. Read `Documentation/scheduler/sched-design-CFS.rst` from the kernel tree or the HTML docs. Write a one-page summary of `vruntime` and the red-black tree.
2. Compare CFS with EEVDF from public kernel notes. Write what "earliest eligible" adds to fairness.

---

## Context switch cost

A context switch is the kernel path that stops one thread and starts another on the same core.

The kernel must:

1. Save general-purpose registers and the program counter.
2. Save the stack pointer.
3. Switch the kernel bookkeeping to the next thread.
4. Switch address spaces if the next thread belongs to another process (change page tables, topic 8).
5. Restore the next thread's registers.
6. Return to user mode if the next thread is a user thread.

Costs:

- direct CPU time in the kernel
- cache and TLB disturbance (the next thread needs different memory)
- pipeline disruption

A switch between two threads of the same process is often cheaper than a switch between two processes. The page tables can stay. The stacks still differ.

A high switch rate appears when:

- the time slice is too small
- many threads wake each other (ping-pong)
- too many threads fight for few cores

Do not optimize for zero switches. Blocking I/O must switch. Optimize for unnecessary switches: too many threads, too much lock contention, and busy-wait that yields in a tight loop.

`vmstat` shows `cs` (context switches) per second. `pidstat -w` shows voluntary and involuntary switches. Voluntary means the thread blocked. Involuntary means the scheduler preempted it.

### Questions

#### Theoretical questions

1. What is a context switch?
2. Why is a process switch often more expensive than a thread switch in the same process?
3. What is TLB disturbance at a high level?
4. What is the difference between a voluntary and an involuntary switch?
5. Why does a small time slice increase context-switch cost?

#### Easy practical tasks

1. Run `vmstat 1 5`. Write the `cs` column meaning from `man 8 vmstat`.
2. Write a numbered list of context-switch steps from this section.
3. Open `man 1 pidstat` if installed. Write what `-w` shows.
4. Write four sentences about context-switch cost. Use only facts from this section.

#### Medium practical tasks

1. Compare `vmstat` idle versus during two fighting CPU loops. Write how `cs` changes.
2. Write two threads that increment a locked counter in a tight loop. Compare `pidstat -w` with a version that increments an atomic without a long hold. Write a hypothesis.
3. Draw the cache/TLB cost as "cold caches after switch" next to the register-save cost.

#### Advanced practical tasks

1. Use `perf stat -e context-switches,cpu-migrations` on a short program if `perf` works. Write the counts.
2. Read OSTEP on context switch. Write a one-page estimate: what is hardware cost versus kernel software cost.

---

## Affinity and multicore

A machine has more than one core. The scheduler may run a thread on any core, or it may keep the thread on one core.

Cache affinity (soft affinity) is the idea that a thread should return to the core that already has its cache lines. Linux tries to keep a thread on the same core when that is fair.

CPU affinity (hard affinity) is a mask that limits the cores that a thread may use. Linux exposes this with `sched_setaffinity` and the `taskset` command.

```text
taskset -c 0 ./cpu-loop
```

Reasons to set affinity:

- measurement: pin a benchmark to one core
- realtime: keep a thread away from noisy cores
- NUMA: keep a thread near its memory (advanced)

Reasons not to set affinity in a normal app:

- you can leave a core idle while another core is busy
- you fight the scheduler
- a VM already virtualizes CPUs

A thread migration is a move from one core to another. Migration has a cost (caches, TLB). The scheduler migrates to balance load.

`/proc/<pid>/status` has `Cpus_allowed`. `lscpu` shows the topology: sockets, cores, and hardware threads (SMT). Two hardware threads on one core share some execution units. They are not two full cores.

### Questions

#### Theoretical questions

1. What is CPU affinity?
2. What is cache affinity?
3. Why can a hard affinity mask hurt throughput?
4. What is a thread migration?
5. Why are two SMT threads on one core not two independent cores?

#### Easy practical tasks

1. Run `nproc` and `lscpu | head`. Write the number of logical CPUs.
2. Run `taskset -p $$`. Write the affinity mask of your shell.
3. Open `man 1 taskset`. Write how to start a command on core 0.
4. Write four sentences about affinity. Use only facts from this section.

#### Medium practical tasks

1. Start two CPU loops with `taskset -c 0` for both. Compare `top` with two loops that have no mask. Write the difference.
2. Read `/proc/self/status` for `Cpus_allowed_list`. Write the list.
3. Draw a two-core chip. Show a thread that stays on core 0 for cache affinity, then a migration to core 1 for load balance.

#### Advanced practical tasks

1. Read `man 2 sched_setaffinity` and `man 7 numa` if present. Write how NUMA affinity differs from a simple CPU mask.
2. Measure a memory-heavy loop pinned to one core versus bouncing across cores (`taskset` vs default) with `time`. Write a careful conclusion (noise is common).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A laptop runs one compiler (CPU-bound) and one editor that waits for keys (I/O-bound). How do preemption, CFS fairness, and nice work together?
2. When do you use a textbook Gantt chart, and when do you use `top` on Linux?
3. How do context-switch costs change your choice of time slice and thread count?
4. Why is `taskset -c 0` useful in a benchmark and harmful as a default in a server?
5. How does `SCHED_FIFO` break the CFS fairness idea?

#### Easy practical tasks

1. Write a one-page cheat sheet: CPU-bound, I/O-bound, preemptive, FIFO, SJF, RR, nice, CFS, `vruntime`, context switch, affinity.
2. Collect `nproc`, `nice` of your shell, `vmstat 1 3`, and `taskset -p $$` into a report.
3. Draw one Gantt chart (RR) and one Linux picture (two tasks, `vruntime`, two cores).
4. Open `man 7 sched`. Write five policy or concept names that the page lists.

#### Medium practical tasks

1. Run the same CPU loop at nice 0 and nice 19 against a competitor at nice 0, both pinned to one core. Write the CPU share that you observe.
2. Write a short script that prints context switches from `vmstat` before and during `find /usr >/dev/null`. Explain the increase.
3. Design a table that maps each algorithm in this topic to a metric it optimizes (turnaround, response, fairness).

#### Advanced practical tasks

1. Build the job simulator from the algorithm section and add nice weights as a simple multiplier on slice or on `vruntime`. Document the model limits.
2. Read the current Linux scheduler documentation for your kernel series. Write what default scheduler your kernel uses and one feature that textbook RR does not have.
