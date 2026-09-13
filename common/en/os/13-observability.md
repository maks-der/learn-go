# 13. Observability

## Description

You can watch a live Unix system. This topic explains the daily tools `ps`, `top`, `lsof`, and `strace`. You learn `/proc/self/maps`, a first view of `perf` and eBPF, and how the OS records logs, crash dumps, and time. You learn the difference between a wall clock and a monotonic clock.

Complete this topic after networking from the OS. Complete this topic before SMP locks and next steps (topic 14). You already know processes, maps, and system calls.

Use one term for each concept. `ps` lists processes. `top` shows a live CPU and memory view. `lsof` lists open files. `strace` traces system calls. `/proc/self/maps` is the address-space listing of the current process. `perf` samples hardware and kernel events. eBPF is a safe in-kernel program model for tracing and filtering. Wall-clock time is civil time (seconds since an epoch, adjustable). Monotonic time only moves forward for measuring intervals. Do not mix `strace` with `perf`. Do not mix wall-clock time with a monotonic timer.

---

## `ps`, `top`, `lsof`, `strace`

`ps` prints a snapshot of processes. Useful columns: `pid`, `ppid`, `stat`, `pcpu`, `pmem`, `cmd`. `ps -eL` or `ps -T` adds threads. Topic 2 and topic 3 used `ps`. This section treats `ps` as an observability tool. Read `STAT` letters. Read `TIME` as CPU time, not wall time.

`top` (or `htop`) refreshes the view. You see load average, memory, and the busiest processes. `top -H` shows threads. `P` in `top` can sort by CPU. Do not leave a mistaken `k` (kill) on a shared host. `uptime` prints load average only.

`lsof` lists open files: regular files, pipes, sockets, cwd. `lsof -p <pid>` is a process view. `lsof /path` shows who has a file open. `ss` is often better for sockets. `lsof` is better when you ask "which process holds this file". `/proc/<pid>/fd` is the raw list.

`strace` traces system calls of one process. `strace -p <pid>` attaches (needs rights). `strace -f` follows children. `strace -e openat,read,write` filters. `strace -c` counts calls. `strace` is slow. It changes timing. It can hide races or create them. Use it to see the user-kernel border (topic 1). Do not leave `strace` on a production process without a plan.

Typical debug order:

1. `ps` or `top`: who is busy or stuck
2. `lsof` or `/proc/pid/fd`: what is open
3. `strace`: which syscall blocks or fails
4. `ss` if the FD is a socket (topic 12)

`/usr/bin/time -v` (GNU time) prints faults and RSS. It is not the shell keyword `time` on all shells.

### Questions

#### Theoretical questions

1. What does `ps` show that `top` also shows, and what is different?
2. What does `lsof` list besides regular files?
3. Why does `strace` change program timing?
4. What does `strace -f` do?
5. When do you pick `ss` instead of `lsof`?

#### Easy practical tasks

1. Open `man 1 ps`, `man 1 top`, `man 8 lsof`, and `man 1 strace`. Write one sentence for each.
2. Run `ps -o pid,stat,cmd` and `lsof -p $$ | head` if `lsof` exists. Write two open files of your shell.
3. Write five sentences about these four tools. Use only facts from this section.
4. Run `strace -c true`. Write the three calls with the highest counts.

#### Medium practical tasks

1. Run `top` for five seconds (or `top -b -n 1`). Write the process with the highest `%CPU`.
2. Start `sleep 60` and run `strace -p <pid>`. Write the syscall that sleeps. Detach. Stop `sleep`.
3. Use `lsof -p` on a process that has a file open in another terminal (`less /etc/hostname`). Write the FD.

#### Advanced practical tasks

1. Trace `cat` of a small file with `strace -e openat,read,write,close`. Write how many `read`s you see and the buffer size if visible.
2. Write a one-page note: when `strace` is the wrong tool (short CPU loops, need for `perf`).

---

## `/proc/self/maps`

`/proc/<pid>/maps` lists virtual memory regions of a process. `/proc/self/maps` is the same file for the caller.

Each line has:

- address range
- permissions `rwxp` (topic 5)
- offset
- device and inode (for file maps)
- pathname (`[heap]`, `[stack]`, `[vdso]`, a library path)

This is the live view of topic 5. You can see text, heap, stack, mmap regions, and the vDSO.

`/proc/self/smaps` adds sizes (`Rss`, `Pss`, swap). It is heavier. `pmap` is a pretty printer.

`/proc/self/status` has `VmSize`, `VmRSS`, and UIDs. Topic 1 already used it.

Deleted maps show `(deleted)` when a binary or a library was replaced but still mapped. That line is a common "I upgraded but the process still runs old code" fact.

A map line is not a page table dump. It is the VMA (virtual memory area) list. Many pages inside a VMA can still be unfaulted (topic 5).

`/proc/self/maps` is observability for memory. Combine it with `ps` RSS and with `pmap -x`. Do not parse `maps` with weak string code in production if you can use a library. The format is stable enough for labs.

The fourth permission letter `p` or `s` is private versus shared. Shared maps match topic 5 shared memory.

### Questions

#### Theoretical questions

1. What does one `maps` line describe?
2. What do `[heap]` and `[stack]` mean?
3. How does `smaps` differ from `maps`?
4. What does `(deleted)` on a map line mean?
5. Why is a VMA not the same as one page?

#### Easy practical tasks

1. Run `cat /proc/self/maps | head -n 15`. Write one `r-x` line and one `rw-` line.
2. Open `man 5 proc` for `maps`. Write the meaning of the permission field.
3. Write five sentences about `/proc/self/maps`. Use only facts from this section.
4. Run `pmap $$ | head` if `pmap` exists. Write how it relates to `maps`.

#### Medium practical tasks

1. Write a C program that `malloc`s 16 MiB and touches it. Print `maps` before and after. Mark the heap or an `mmap` region.
2. Compare `VmRSS` in `status` with `Rss` in `smaps` (first heap block). Write that they are related, not always one number.
3. Draw three VMAs: text, heap, stack. Label a hole.

#### Advanced practical tasks

1. Find a `(deleted)` map on a long-running process after a package upgrade (or describe how you would). Write the risk.
2. Write a one-page note: `Pss` versus `Rss` for shared libraries.

---

## `perf` and eBPF (intro)

`perf` is the Linux performance tool. It samples events: CPU cycles, cache misses, scheduler events, and some tracepoints. `perf stat` counts events for one command. `perf record` writes a sample file. `perf report` shows a profile. You need rights or `perf_event_paranoid` settings. Some VMs hide counters.

`perf` answers "where does the CPU time go" better than `strace`. It does not print every syscall by default. You can use `perf trace` as a lighter `strace`-like view on some systems.

eBPF (extended Berkeley Packet Filter) is a virtual machine in the kernel. You load a verified program. The program runs at a hook: a tracepoint, a kprobe, a socket, a cgroup, or XDP. The verifier rejects unsafe programs (unbounded loops, bad pointer use). User space reads maps (hash tables, rings) that the program fills.

Typical eBPF uses: `bpftrace` one-liners, BCC tools (`execsnoop`, `opensnoop`), systemd and Cilium policies, and custom tracers. You do not write C kernel modules for those jobs.

Intro rules:

- start with `perf stat -- sleep 1` or `perf stat ls`
- read a `bpftrace` example (`tracepoint:syscalls:sys_enter_openat`)
- do not load unknown BPF objects as root
- eBPF is not a substitute for `ps` on day one

`strace` is ptrace. `perf` is sampling or counting. eBPF is in-kernel hooks with maps. Pick the tool that matches the question.

Kernel modules (`lsmod`) are another internals path. A module is not eBPF. Topic 8 named modules as drivers. Do not `insmod` random files.

### Questions

#### Theoretical questions

1. What does `perf stat` show?
2. What does `perf record` produce?
3. What is eBPF?
4. What does the BPF verifier reject?
5. How does eBPF differ from `strace`?

#### Easy practical tasks

1. Open `man 1 perf` if it exists. Write one sentence about `perf stat`.
2. Run `perf stat -- true` if you have rights. Write one counter name that you see, or write the error.
3. Write five sentences about `perf` and eBPF. Use only facts from this section.
4. Search a `bpftrace` one-liner gallery (read only). Write one example name.

#### Medium practical tasks

1. Run `perf stat -- dd if=/dev/zero of=/dev/null bs=1M count=1024`. Write whether the task looks CPU-bound from the counters.
2. Draw user program, `perf` sample, kernel event, report.
3. Read a BCC `opensnoop` overview. Write six sentences: what hook it uses and what it prints.

#### Advanced practical tasks

1. If policy allows, run `bpftrace -e 'tracepoint:syscalls:sys_enter_execve { printf("%s\n", comm); }'` for a few seconds in a VM. Write three command names. Stop the probe.
2. Write a one-page intro: `perf_event_paranoid` levels and why a container may not see hardware events.

---

## Logs, crash dumps, wall clock vs monotonic time

Linux systems that use systemd store many logs in the journal. The daemon is `journald`. The tool is `journalctl`. Units and the kernel send structured records.

Useful views:

- `journalctl -b` : this boot
- `journalctl -k` : kernel messages
- `journalctl -u ssh` : one unit
- `journalctl -p err` : priority

`dmesg` reads the kernel ring buffer. `journalctl -k` and `dmesg` overlap. Some programs still write text files under `/var/log`. Logs need rotation and size limits. Logs can contain secrets. Do not publish a raw dump.

A crash dump is a copy of memory and CPU state after a panic (kernel) or after a process abort (core dump). `ulimit -c` and `coredumpctl` (systemd) manage user core files. `kdump` captures a kernel crash. A dump is large and sensitive. Use it to debug, then restrict access.

Wall-clock time is `CLOCK_REALTIME`. It is the civil clock. NTP or a user can step it. A step can jump backward. Do not use wall-clock time to measure a duration in a program.

Monotonic time is `CLOCK_MONOTONIC` (or `CLOCK_MONOTONIC_RAW`). It only moves forward for interval measure. A suspend can pause some monotonic clocks. `CLOCK_BOOTTIME` includes suspend. Use monotonic time for timeouts and benchmarks.

`date` prints the wall clock. `cat /proc/uptime` is seconds since boot (not a wall clock). `timedatectl` shows NTP state on systemd machines.

Hibernation writes RAM to disk and powers off. Suspend keeps RAM powered. Those are power states, not logs. They still change clocks and can lose devices. This section only names them so that you do not mix them with a crash dump.

### Questions

#### Theoretical questions

1. What does `journalctl -k` show?
2. What is a crash dump?
3. Why must you not measure a duration with `CLOCK_REALTIME`?
4. What is `CLOCK_MONOTONIC` for?
5. How does a core dump differ from a kernel crash dump?

#### Easy practical tasks

1. Open `man 1 journalctl` and `man 2 clock_gettime`. Write one sentence for each.
2. Run `date` and `cat /proc/uptime`. Write both values and which clock class each is.
3. Write five sentences about logs, dumps, and clocks. Use only facts from this section.
4. Run `ulimit -c`. Write the core-file size limit.

#### Medium practical tasks

1. Run `journalctl -b -n 10 --no-pager` if you have rights. Write two unit or kernel lines in your own words.
2. Write a C snippet plan: start `CLOCK_MONOTONIC`, `sleep(1)`, stop, print elapsed. Then run it. This handbook does not contain the source.
3. Read `man 1 coredumpctl` if it exists. Write how you list dumps without sending them anywhere.

#### Advanced practical tasks

1. Read `man 7 time` and `man 2 clock_nanosleep`. Write a one-page note: which clock to use for a 50 ms poll timeout.
2. Write a policy: who may read `journalctl -a` and core files on a multiuser machine.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. In which order do you use `top`, `lsof`, `strace`, and `perf` for a stuck server?
2. How do `/proc/self/maps` and `ps` RSS tell different memory stories?
3. When is eBPF a better hook than `strace -p` on a hot process?
4. Why can a log timestamp and a monotonic timeout disagree after an NTP step?
5. What must you hide when you share a `journalctl` snippet or a core dump?

#### Easy practical tasks

1. Write a cheat sheet: `ps`, `top`, `lsof`, `strace`, `maps`, `smaps`, `perf stat`, eBPF, `journalctl`, `dmesg`, `CLOCK_REALTIME`, `CLOCK_MONOTONIC`.
2. Run `ps`, `cat /proc/self/maps | wc -l`, and `date`. Write one line each.
3. Draw a debug flow: symptom, tool, next tool.
4. Bookmark `man 5 proc`, `man 1 strace`, and `man 1 journalctl`.

#### Medium practical tasks

1. Pick a command (`ls` or your echo server). Observe it with `strace -c` and `/proc/<pid>/maps` while it runs. Write a six-line report.
2. Compare `time ls` (shell) with `/usr/bin/time -v ls` if present. Write extra fields that GNU time adds.
3. Document a data-handling rule: no secrets in labs that you paste into chat, redact `journalctl` user names if needed.

#### Advanced practical tasks

1. Write a one-page observability plan for the tiny shell in topic 14: which tool shows zombies, which tool shows `execve`, which clock you use for a timeout.
2. Read an OSTEP or Linux chapter on tracing. Write a map from textbook "profiling" to `perf` and from "ptrace" to `strace`.
