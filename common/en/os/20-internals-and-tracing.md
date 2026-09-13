# 20. Internals and Tracing

## Description

You can watch a live Unix system. This topic explains syscall tables, `/proc/self/maps`, and the daily tools `strace`, `lsof`, `ps`, and `top`. You get an introduction to `perf` and eBPF, awareness of kernel modules, and an optional path through a small piece of kernel source.

Complete this topic after persistence of OS state. Complete this topic before the advanced lock topic (topic 21). You already know processes, maps, and system calls. You now observe them on a running machine.

Use one term for each concept. A syscall table maps a number to a kernel function. `/proc/self/maps` is the address-space listing of the current process. `strace` traces system calls. `lsof` lists open files. `perf` samples hardware and kernel events. eBPF is a safe in-kernel program model for tracing and filtering. A kernel module is code that you can load into the running kernel. Do not mix `strace` with `perf`. Do not mix a module with a user-space library.

---

## Syscall tables

A system call has a name (`read`) and a number (`__NR_read`). The number is the ABI. User space puts the number in a register (for example `rax` on x86-64) and executes a trap instruction (`syscall`). The kernel indexes a table of function pointers and runs the handler.

The table is per architecture. x86-64 and aarch64 numbers differ. 32-bit compatibility tables exist on 64-bit kernels. A binary from another ABI fails in surprising ways if you assume one number.

Linux headers define the numbers (`unistd.h`). `ausyscall` or `cat /usr/include/asm/unistd_64.h` can show them. `strace` prints names. The raw trace still uses numbers inside the kernel.

Adding a system call is a kernel-development event. User programs should call libc. libc wraps the ABI and handles `restart_syscall` and 64-bit offsets.

`/proc/kallsyms` lists kernel symbols (root may see all). The syscall handlers have names such as `__x64_sys_read`. You do not need to memorize the list. You need to know that the table exists and that `strace` decodes it.

A wrong number is not a C compiler error. It is an `ENOSYS` or a wild call. This is why you do not invoke `syscall()` with a guessed number.

### Questions

#### Theoretical questions

1. What does a syscall table store?
2. Why do syscall numbers differ across architectures?
3. What instruction starts a system call on x86-64?
4. Why should a program call libc instead of a raw number?
5. What does `ENOSYS` mean?

#### Easy practical tasks

1. Open `man 2 syscall` and `man 2 syscalls`. Write one sentence for each.
2. Run `grep __NR_read /usr/include/asm/unistd_64.h` or an equivalent header. Write the number if you find it.
3. Write five sentences about syscall tables. Use only facts from this section.
4. Draw user register, `syscall` instruction, table, `sys_read`.

#### Medium practical tasks

1. Use `strace -e raw=read` (or a similar flag) on `true` or `cat` and compare with a normal `strace`. Write what "raw" changes.
2. Compare one syscall number on x86-64 with the man-page note for arm64 (online). Write that the names match and the numbers may not.
3. Read `man 2 intro`. Write six sentences on ABI versus API.

#### Advanced practical tasks

1. Read a Linux `syscall_64.tbl` excerpt online. Write a one-page note: how a new syscall gets a number.
2. Write a tiny program that calls `syscall(SYS_getpid)` and `getpid()`. Confirm they match.

---

## `/proc/self/maps`

`/proc/<pid>/maps` lists virtual memory regions of a process. `/proc/self/maps` is the same file for the caller.

Each line has:

- address range
- permissions `rwxp` (topic 10)
- offset
- device and inode (for file maps)
- pathname (`[heap]`, `[stack]`, `[vdso]`, a library path)

This is the live view of topic 8. You can see text, heap, stack, mmap regions, and the vDSO.

`/proc/self/smaps` adds sizes (`Rss`, `Pss`, swap). It is heavier. `pmap` is a pretty printer.

`/proc/self/status` has `VmSize`, `VmRSS`, and UIDs. Topic 1 already used it.

Deleted maps show `(deleted)` when a binary or a library was replaced but still mapped. That line is a common "I upgraded but the process still runs old code" fact.

Do not parse `maps` with weak string code in production if you can use a library. The format is stable enough for labs.

A map line is not a page table dump. It is the VMA (virtual memory area) list. Many pages inside a VMA can still be unfaulted (topic 9).

### Questions

#### Theoretical questions

1. What does one `maps` line describe?
2. What does `[vdso]` mean at a high level?
3. What does `(deleted)` mean on a path?
4. How does `maps` differ from a full page-table dump?
5. What extra data does `smaps` add?

#### Easy practical tasks

1. Run `cat /proc/self/maps`. Write one `r-x` line and one `rw-` line.
2. Open `man 5 proc` for `/proc/[pid]/maps`. Write one sentence.
3. Run `head /proc/self/status`. Write `VmSize` and `VmRSS`.
4. Write five sentences about `maps`. Use only facts from this section.

#### Medium practical tasks

1. Run `cat /proc/self/maps` from a C program and from the shell. Write whether the stack address range looks different.
2. Draw three VMAs: text, heap, stack. Copy the permission letters from your `maps`.
3. Use `pmap $PID` on your shell if `pmap` exists. Write the total.

#### Advanced practical tasks

1. `mmap` a file and find its line in `/proc/self/maps` from the same process (read the file after `mmap`). Write the path and the range.
2. Read about `[vvar]` and `[vdso]`. Write a one-page note on why `gettimeofday` can avoid a real syscall.

---

## `strace`, `lsof`, `ps`, `top`

These tools are the daily OS lab kit.

`strace` attaches to a process and prints system calls and signals. `strace ls` starts `ls`. `strace -p PID` attaches to a running process (rights needed). `-e` filters. `-c` counts. `-o` writes a file. `strace` is slow. It changes timing. It can reveal secrets (file names, buffer contents with `-s`).

`lsof` lists open files. "File" includes sockets, pipes, and cwd. `lsof -p PID` and `lsof /path` are common. `/proc/<pid>/fd` is the raw view. `ls -l /proc/self/fd` is enough for many labs.

`ps` lists processes. `ps -o pid,ppid,stat,cmd` is a useful format. `ps` reads `/proc`. Topic 4 used `ps`.

`top` and `htop` sample CPU and memory over time. They show the scheduler's view (topic 6), not a proof of a bug. `pidstat` and `vmstat` add more counters.

Use the tools in this order for a hang:

1. `ps` and `top`: is the process runnable, sleeping, or zombie?
2. `lsof` or `/proc/fd`: what is it waiting on?
3. `strace -p`: which syscall is blocked?

Do not `strace` a production database without a plan. Do not share an `strace` log that contains tokens.

### Questions

#### Theoretical questions

1. What does `strace` print?
2. What does `lsof` list that `ls` of a directory does not list?
3. Where does `ps` get its data on Linux?
4. Why is `top` not a proof of a race?
5. What is a safe first step when a process hangs?

#### Easy practical tasks

1. Open `man 1 strace`, `man 8 lsof`, `man 1 ps`, and `man 1 top`. Write one sentence for each.
2. Run `strace -c true`. Write the three calls with the highest counts if any appear.
3. Run `ls -l /proc/self/fd`. Write what 0, 1, and 2 point to.
4. Run `ps -o pid,stat,cmd -p $$`. Write the STAT letter.

#### Medium practical tasks

1. Run `strace -e openat,execve ls /tmp` and write three paths.
2. Start `sleep 300` in another terminal. Use `ps` and `ls -l /proc/<pid>/fd`. Then `strace -p` and interrupt with Ctrl-C on `strace` only.
3. Compare `top` batch mode `top -b -n 1` with `ps`. Write one field that `top` adds.

#### Advanced practical tasks

1. Trace `cat file` with `strace -tt -T -o trace.txt`. Write the `openat`/`read`/`write`/`close` story in ten lines.
2. Use `lsof -i` or `ss` plus `lsof -p` on a small server from topic 18. Write which FDs are sockets.

---

## `perf`, eBPF (intro)

`perf` is the Linux profiler and event tool. It can count CPU cycles, cache misses, and kernel tracepoints. `perf stat` counts a command. `perf record` plus `perf report` samples a call graph. You need rights (`perf_event_paranoid`) and matching debug symbols for a clear report.

`perf` is statistical. A short run can lie. A production system can refuse access.

eBPF (extended Berkeley Packet Filter) is a virtual machine in the kernel. You load a verified program. The program attaches to a tracepoint, a kprobe, a perf event, or a network hook. The verifier rejects unsafe programs (unbounded loops, invalid memory). User space reads maps or ring buffers.

Typical tools: `bpftrace` (one-liners), BCC, `libbpf`. Examples: count `read` sizes, trace `vfs_read`, or build a histogram of syscall latency.

eBPF is not a replacement for `strace` in every lab. `strace` is simpler for one process. eBPF is better when you cannot attach to every process or when the probe must stay cheap.

Capabilities and unprivileged eBPF policy changed across kernels. Many systems need root for probes. Do not load unknown BPF objects.

This section is an introduction. You do not need to write a production tracer.

### Questions

#### Theoretical questions

1. What does `perf stat` measure?
2. Why is `perf record` statistical?
3. What does the eBPF verifier reject at a high level?
4. When do you pick `strace` instead of eBPF?
5. Why can unprivileged BPF be disabled?

#### Easy practical tasks

1. Open `man 1 perf` if present. Write two subcommands.
2. Run `perf stat true` if allowed. Write one counter name. If it fails, write the error.
3. Run `command -v bpftrace`. Write whether it exists.
4. Write five sentences about `perf` and eBPF. Use only facts from this section.

#### Medium practical tasks

1. Run `perf list | head` if allowed. Write three event names.
2. Draw user program, syscall, kprobe or tracepoint, eBPF program, user map read.
3. Read a `bpftrace` one-liner example from the man page or docs. Write what it counts in six sentences. Do not run hostile scripts.

#### Advanced practical tasks

1. If you have rights in a VM, `perf record -g` a CPU-bound loop and open `perf report`. Write the hottest symbol.
2. Write a `bpftrace` script that counts `tracepoint:syscalls:sys_enter_openat` for ten seconds. Document the need for root. Use a VM.

---

## Kernel modules (awareness)

A kernel module is an object file that the kernel can load (`insmod`, `modprobe`) and unload (`rmmod`). Drivers are often modules (topic 14). Filesystems can be modules. `lsmod` lists loaded modules. `/lib/modules/$(uname -r)/` holds the files.

A module runs in kernel mode. A bug can panic the machine. A hostile module is a kernel exploit. Only load modules from your distribution or from a source that you trust and that you built for this kernel version.

`modprobe` resolves dependencies. `modinfo` prints license, description, and parameters.

Some features are built-in (`=y`) and never appear as a loadable file. `/lib/modules` does not list them as `.ko` in use, but `/proc/kallsyms` still has the symbols.

Secure Boot can require a signed module (topic 17). Out-of-tree modules (NVIDIA, ZFS, DKMS) break on some kernel upgrades. That is an operations fact.

This path does not ask you to write a module. Awareness means: you can read `lsmod`, you do not `rmmod` at random, and you treat `/dev` plus a new module as a trust change.

Do not download a random `.ko` and load it.

### Questions

#### Theoretical questions

1. What is a kernel module?
2. Why can a module panic the machine?
3. What does `modprobe` add over `insmod`?
4. Why can Secure Boot reject a module?
5. Why must the module match the running kernel version?

#### Easy practical tasks

1. Run `lsmod | head`. Write three names.
2. Open `man 8 modprobe` and `man 8 lsmod`. Write one sentence for each.
3. Run `modinfo` on one name from `lsmod`. Write the description.
4. Write five sentences about modules. Use only facts from this section.

#### Medium practical tasks

1. Run `uname -r` and `ls /lib/modules/$(uname -r) | head`. Write that the directory matches the release.
2. Draw user `modprobe`, kernel, and a new `/dev` node.
3. Read `dmesg` after boot for "Loading" or driver lines. Write two module-related lines.

#### Advanced practical tasks

1. Read the kernel module documentation overview. Write a one-page note on `init` and `exit` of a module (idea only).
2. Write a policy: who may `insmod` on a lab VM versus on a laptop.

---

## Reading a small kernel path (optional)

You can read Linux source without writing a driver. The goal is to match a name from `strace` or `/proc` to a function.

How to start:

1. Pick one call that you already use (`openat`, `read`, `clone`, `mmap`).
2. Open the kernel tree on [https://elixir.bootlin.com/linux](https://elixir.bootlin.com/linux) or a local clone.
3. Search the syscall entry (`sys_read` or `ksys_read`).
4. Read only the top function and one helper. Ignore `#ifdef` forests on the first pass.
5. Write the path in five boxes: libc, syscall, VFS, filesystem, block or page cache.

Rules:

- Match your `uname -r` series when you can (or read a current stable).
- Do not try to understand every lock on day one (topic 21).
- xv6 is a smaller OS for a full read (topic 22).
- Licensed source is public. You still do not copy huge files into homework. Summarize.

Optional is optional. If you skip this section, you can still finish the path. If you read one path, you remember VFS (topic 12) better.

A good first path is `read`: `ksys_read` to `vfs_read` to `ext4_file_read_iter` (or the current names). A good second path is `clone` for `fork`.

### Questions

#### Theoretical questions

1. Why start from a syscall that you already traced?
2. What is Elixir (or a local clone) for?
3. Why is xv6 easier to read as a whole?
4. What five boxes does this section suggest?
5. Why must you not paste a huge kernel file into a report?

#### Easy practical tasks

1. Open Elixir or a kernel GitHub mirror. Search `vfs_read`. Write the file path that you see.
2. Write five sentences about how to read a small path. Use only facts from this section.
3. Bookmark kernel.org documentation and Elixir.
4. Make a table: xv6 versus Linux. Add "size" and "use in this course".

#### Medium practical tasks

1. Write a five-box diagram for `read` with the real function names that you found (names can differ by version).
2. Compare `man 2 read` with the first thirty lines of the handler. Write three checks that the man page already told you.
3. Find `sys_clone` or `kernel_clone`. Write one sentence on what `fork` shares.

#### Advanced practical tasks

1. Read one OSTEP chapter that matches your path (files or processes). Write a one-page map from the book to the Linux functions.
2. Walk `mmap` from the syscall to `do_mmap` (or current name) and stop. Write the flags that you recognize from topic 8.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. You have a process that is "stuck". How do you combine `ps`, `/proc/fd`, `maps`, and `strace` before you touch `perf`?
2. Why can `strace` and eBPF both show `openat` and still be the wrong first tool for a CPU-bound loop?
3. How does a syscall number connect `strace` output to a kernel function name in Elixir?
4. When is a kernel module the explanation for a new `/dev` node that `lsof` shows?
5. Why is `/proc/self/maps` not enough to know which pages are in RAM?

#### Easy practical tasks

1. Write a one-page cheat sheet: syscall number, `maps`, `strace`, `lsof`, `ps`, `top`, `perf`, eBPF, `lsmod`, Elixir.
2. Run `strace -c true`, `cat /proc/self/maps | wc -l`, `ls /proc/self/fd | wc -l`, and `lsmod | wc -l`. Comment each number.
3. Draw one figure: user, libc, syscall table, VFS, and a `perf` sample arrow.
4. Bookmark `man 1 strace`, `man 5 proc`, and Elixir.

#### Medium practical tasks

1. Write a script that prints PID, `stat`, FD count, and the first `maps` line for the script itself. Run it.
2. Trace `uname` and then find the handler name online. Write the map from the `strace` line to the function.
3. Use `perf stat` and `strace -c` on the same `dd` command if allowed. Write what each tool measures.

#### Advanced practical tasks

1. Read a small Linux path and the matching OSTEP chapter. Write a one-page dual map.
2. In a VM, record `bpftrace` or `perf` while you run the hang-debug trio (`ps`, `lsof`, `strace`) on a blocked `sleep`. Write which tool saw the `nanosleep` wait most clearly.
