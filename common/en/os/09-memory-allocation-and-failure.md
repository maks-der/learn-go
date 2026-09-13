# 9. Memory: Allocation and Failure

## Description

The kernel and the C library allocate memory in pieces. Those pieces waste space or fail. This topic explains internal and external fragmentation, and it surveys allocators: bump, free list, buddy, and slab. You learn the Linux OOM killer, demand paging, copy-on-write after `fork`, and thrashing.

Complete this topic after address spaces. Complete this topic before protection and sharing. You now know pages. You now learn when pages do not exist yet, and when the system runs out of RAM.

Use one term for each concept. Fragmentation is wasted space. Demand paging loads a page when the process first uses it. Copy-on-write shares a frame until a write. The OOM killer is a Linux kernel path that kills a process to free memory. Thrashing is a state where the machine spends most of its time paging. Do not mix internal fragmentation with external fragmentation. Do not mix a user-space allocator with the kernel page allocator.

---

## Internal vs external fragmentation

Fragmentation is wasted memory that you cannot use for the current request.

Internal fragmentation is waste inside an allocated block. The allocator rounds a request up. A 1-byte `malloc` can consume 16 or 32 bytes of heap metadata plus payload. A 1-byte need that uses a full 4 KiB page wastes almost a page. The unused bytes belong to the allocation. Other objects cannot use them until the block is freed.

External fragmentation is waste outside allocated blocks. Free memory exists, but it is split into holes. No hole is large enough for the next request. The sum of free bytes can still be large. This problem appears with variable-sized segments or with a heap that allocates mixed sizes and frees them in a bad order.

Pages reduce some external fragmentation of physical RAM because frames have one size. The kernel can satisfy a one-page request with any free frame. Multi-page contiguous requests (DMA, huge pages) still suffer external fragmentation of physical memory.

User-space `malloc` fights both types. It uses size classes and caches to limit internal waste. It splits and coalesces free chunks to limit external waste. It never removes all waste.

Do not call every slow program "fragmented". First measure RSS, cache, and page faults.

### Questions

#### Theoretical questions

1. What is internal fragmentation?
2. What is external fragmentation?
3. Why does a 4 KiB page cause internal fragmentation for a 1-byte object if that object has its own page?
4. Why do equal-sized frames help the kernel allocate one page?
5. When does free memory exist but `malloc` still fail (or call `mmap` again) because of holes?

#### Easy practical tasks

1. Write a two-column table: "Internal" and "External". Add three rows of examples.
2. Draw a heap with three used blocks and two holes that cannot satisfy a large request.
3. Open `man 3 malloc`. Write that the size you get can be larger than the size you asked for.
4. Write five sentences about fragmentation. Use only facts from this section.

#### Medium practical tasks

1. Write a program that `malloc`s 1 byte 100000 times and prints RSS from `/proc/self/status` (`VmRSS`). Compare with one `malloc` of 100000 bytes.
2. Explain in six sentences why huge pages can increase internal fragmentation.
3. Search "external fragmentation buddy allocator". Write how free lists per order relate to holes.

#### Advanced practical tasks

1. Read a glibc or jemalloc size-class document. Write three size classes and the internal waste for a request that is one byte over a class.
2. Write a one-page note on physical memory fragmentation and compaction in Linux (high-level).

---

## Allocators (bump, free list, buddy, slab — survey)

An allocator hands out memory and takes it back. Different designs fit different lifetimes.

Bump (arena, linear) allocator: a pointer starts at the base. Each request moves the pointer forward. The allocator does not free single objects, or it frees all at once. It is fast. It has no per-object free list. It does not reuse holes until reset. Compilers and request arenas use this pattern.

Free-list allocator: freed blocks go on a list. The next allocation searches the list (first fit, best fit, or next fit). Coalescing merges neighbor free blocks. This design reuses memory. Search can be slow. External fragmentation depends on the fit rule.

Buddy allocator: blocks have sizes that are powers of two. A split turns a block of size 2^(k+1) into two buddies of size 2^k. A free can merge a block with its buddy if the buddy is free. Linux uses a buddy system for pages. Internal fragmentation occurs when you need 9 pages and you get 16.

Slab allocator: the allocator keeps caches of equal-sized objects (inodes, dentries, `task_struct`). A slab is a group of pages split into those objects. Allocation is fast. Internal waste per object is small. Linux uses slab or a close relative (SLUB) for kernel objects.

User `malloc` is a mix: size classes (slab-like), large `mmap` (bump of a region), and free lists. The kernel page allocator is buddy-based. Do not expect `free` in C to return pages to the OS at once. The library may keep them for the next `malloc`.

This section is a survey. Implement only a toy bump or a toy free list in practice tasks if you write code.

### Questions

#### Theoretical questions

1. How does a bump allocator work?
2. What is a free list?
3. Why are buddy sizes powers of two?
4. What problem does a slab cache solve?
5. Which Linux allocator hands out page frames?

#### Easy practical tasks

1. Make a four-row table: bump, free list, buddy, slab. Add "fast path" and "main waste".
2. Draw a buddy split from 16 pages to two 8-page blocks, then to 4.
3. Open `man 3 malloc` and write one sentence about `free`.
4. Write five sentences that survey the four allocators. Use only facts from this section.

#### Medium practical tasks

1. Implement a bump allocator in C over a fixed 64 KiB buffer. Allocate 10 objects. Show that you cannot free one object.
2. Implement a first-fit free list for a homework struct. Allocate, free, allocate again. Print the list.
3. Read `cat /proc/buddyinfo` if the file exists. Write what the columns mean from `man 5 proc` or kernel docs.

#### Advanced practical tasks

1. Read Linux SLUB documentation. Write how a slab cache differs from the buddy page allocator.
2. Compare glibc `malloc` and a bump arena for a compiler pass in a one-page design: lifetime, threads, and fragmentation.

---

## OOM killer (Linux)

Linux overcommits memory by default. `malloc` and `mmap` can succeed even when RAM plus swap cannot hold every page yet. The kernel hopes that processes will not touch every page.

When the kernel needs a frame and cannot reclaim enough, it can invoke the out-of-memory (OOM) killer. The OOM killer selects a process and sends it a kill (usually `SIGKILL`). The goal is to free memory and keep the system running.

Selection uses a heuristic. `/proc/<pid>/oom_score` and `oom_score_adj` influence the choice. A large process that uses much RAM is a likely victim. The kernel tries not to kill `init`.

You see OOM events in `dmesg` or `journalctl -k`. The log names the victim and the memory state.

OOM is not the same as `malloc` returning `NULL`. With overcommit, `malloc` can return a pointer. The process dies later on a page fault when the kernel cannot supply a frame. Always check `malloc`, but also design for later failure.

You can tune overcommit with `vm.overcommit_memory` (`sysctl`). Do not change it on a shared machine. In a VM you may experiment.

Swap delays OOM. Swap does not remove OOM if the working set stays huge. Thrashing can come first.

### Questions

#### Theoretical questions

1. What is overcommit?
2. What is the OOM killer?
3. Why can `malloc` succeed and the process still die later?
4. What is `oom_score`?
5. Why does the kernel try not to kill PID 1?

#### Easy practical tasks

1. Run `cat /proc/self/oom_score` and `cat /proc/self/oom_score_adj`. Write the numbers.
2. Open `man 5 proc` for `oom_score`. Write the purpose of the file.
3. Run `sysctl vm.overcommit_memory vm.swappiness` if allowed. Write the values.
4. Write five sentences about the OOM killer. Use only facts from this section.

#### Medium practical tasks

1. Search `dmesg` or the journal for `Out of memory`. If you find a past event, write the victim name. If not, write that no event is present.
2. Draw the path: `malloc` success, later write to the page, fault, no frame, OOM kill.
3. Read `Documentation/admin-guide/sysctl/vm.rst` on overcommit. Write the meaning of values 0, 1, and 2 at a high level.

#### Advanced practical tasks

1. In a disposable VM with little RAM, run a program that touches a huge anonymous map. Record the OOM log. Do not run this on a host that you need.
2. Write a one-page policy: when to use `vm.overcommit_memory=2` for a database VM, and what can break.

---

## Demand paging

Demand paging means the kernel does not load every page of a program at `exec` or at `mmap`. It creates a mapping. The page table entries can be invalid or marked not present. The first access faults. The kernel then:

1. Finds the mapping (file, anonymous, stack).
2. Allocates a frame if needed.
3. Fills the frame (read the file, or zero the page).
4. Installs a PTE.
5. Returns to the process. The instruction retries.

Benefits: fast start, less RAM for code that never runs, maps of large files that you only touch in part.

Costs: the first access is slow. A burst of faults at start is normal. A fault on the hot path later is a problem.

Minor page fault: the kernel can satisfy the fault without a storage read (for example, the page is already in the page cache, or a zero page). Major page fault: the kernel must read from a storage device.

`ps` and `time` show major fault counts (`majflt`). `/usr/bin/time -v` prints page faults.

Stack growth on Linux uses faults at the edge of the stack mapping. A large jump beyond the guard can still crash.

Demand paging needs a backing store for file pages and for swapped anonymous pages. Topic 13 covers swap devices in more depth. Here you only need the fault-fill-retry loop.

### Questions

#### Theoretical questions

1. What is demand paging?
2. What are the five handler steps in this section?
3. What is the difference between a minor fault and a major fault?
4. Why does demand paging speed up `exec`?
5. Why is a major fault expensive?

#### Easy practical tasks

1. Run `/usr/bin/time -v true` or `time -v true`. Write the page-fault numbers if they appear.
2. Open `man 2 mmap` and find `MAP_POPULATE`. Write how it changes demand paging.
3. Write a numbered list of the fault path.
4. Write four sentences about demand paging. Use only facts from this section.

#### Medium practical tasks

1. `mmap` a large file (`MAP_PRIVATE`) and touch every Nth page in a loop. Compare time and faults with a version that touches nothing.
2. Compare `MAP_POPULATE` versus default for that file. Write the start-time versus first-touch difference.
3. Draw the retry of the faulting instruction after the kernel installs a PTE.

#### Advanced practical tasks

1. Read `man 2 mincore` or `man 2 posix_madvise`. Write how you can see or influence which pages are resident.
2. Read kernel documentation on page faults. Write how a file-backed fault uses the page cache.

---

## Copy-on-write after `fork`

After `fork`, the child has the same virtual mappings as the parent. If the kernel copied every frame at once, `fork` would be slow and would use much RAM.

Copy-on-write (COW) is the fix. Parent and child share the physical frames. The PTEs are read-only. When either process writes a page:

1. The MMU raises a protection fault.
2. The kernel allocates a new frame.
3. The kernel copies the old page to the new frame.
4. The kernel maps the new frame as writable for the writer.
5. The other process keeps the old frame.

Reads do not copy. A child that only `exec`s soon never dirties most pages. That is why the shell `fork`+`exec` path is cheap.

`MAP_PRIVATE` file maps use the same idea after a write. The file stays unchanged. The process gets a private frame.

COW is not a mutex. Two processes do not see each other's writes after the copy. Before the write they see the same bytes.

`vfork` and `posix_spawn` exist because even COW has a cost on huge processes. A database with a multi-gigabyte heap still pays for page-table copies and for later COW faults if the child writes.

Topic 10 returns to COW as a sharing mechanism. This section stresses `fork` and failure: a COW fault can fail if there is no frame. That fault can kill the process or invoke OOM.

### Questions

#### Theoretical questions

1. Why does `fork` not copy every frame immediately?
2. What permission do COW PTEs use until a write?
3. What does the kernel do on the first write to a COW page?
4. Why is `fork` then `exec` cheap for a shell?
5. How can a COW fault fail?

#### Easy practical tasks

1. Write five sentences about COW after `fork`. Use only facts from this section.
2. Draw parent and child that share a frame, then a write that splits the frame.
3. Open `man 2 fork`. Find the sentence about copy-on-write if present. Write it in your own words.
4. Make a table: "Read after fork" and "Write after fork".

#### Medium practical tasks

1. Write a program that forks. The child writes to a global. The parent sleeps, then prints the global. Show that the parent still has the old value.
2. Trace `bash -c 'true'` with `strace -f -e clone,execve`. Write how little work the child does before `exec`.
3. Compare RSS of parent and child in `/proc` immediately after `fork` in a program that allocated a 64 MiB buffer but did not write it in the child.

#### Advanced practical tasks

1. Touch every page of a 64 MiB buffer in the parent, then `fork`, then touch every page in the child. Measure time and RSS. Explain COW faults.
2. Read `man 2 vfork` and `man 3 posix_spawn`. Write when a large process should avoid a naive `fork`.

---

## Thrashing

Thrashing is a state where the machine spends most of its time moving pages between RAM and a storage device, and little time on useful work.

Cause: the working sets of the runnable processes do not fit in RAM. Each process faults. The kernel evicts a page that another process needs next. The disk or SSD is busy. The CPU waits.

Symptoms:

- high major page fault rate
- high I/O wait
- low useful CPU
- a system that feels frozen
- `free` shows little available memory and heavy swap use

Working set is the set of pages that a process needs in a time window. If the sum of working sets is larger than RAM, the scheduler that runs more processes can make thrashing worse. Some systems swap out a whole process. Linux uses global reclaim and the OOM killer as last steps.

Fixes at a high level:

- add RAM
- reduce the number of concurrent large processes
- reduce the working set (smaller caches, fewer tabs)
- use faster storage (SSD helps but does not remove a huge gap)
- avoid swap for some latency-critical VMs (a policy choice)

Do not confuse thrashing with a single tight CPU loop. A CPU loop has high CPU and few major faults. Thrashing has many major faults and poor progress.

### Questions

#### Theoretical questions

1. What is thrashing?
2. What is a working set?
3. Why can running more processes make thrashing worse?
4. What symptoms show thrashing rather than a CPU-bound loop?
5. Name two high-level fixes.

#### Easy practical tasks

1. Write five sentences about thrashing. Use only facts from this section.
2. Make a table: "CPU-bound busy" and "Thrashing". Add CPU, faults, I/O wait.
3. Run `free -h` and `vmstat 1 3`. Write the `si` and `so` column meanings from `man 8 vmstat`.
4. Open `man 5 proc` for `/proc/meminfo` fields `SwapTotal` and `SwapFree`. Write the values on your system.

#### Medium practical tasks

1. Draw a timeline of two processes that steal each other's pages.
2. Read `swappiness` in kernel vm documentation. Write how it biases reclaim toward file pages versus anonymous pages.
3. Compare `vmstat` at idle and during a large `find` or compile. Write whether you see swap activity.

#### Advanced practical tasks

1. In a VM with small RAM and swap, run two programs that each touch a working set larger than half of RAM. Record `vmstat` and responsiveness. Stop the test. Write the result.
2. Read the OSTEP chapter on replacement and thrashing. Write a one-page summary of working-set and clock algorithms at a high level.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A process `malloc`s 8 GiB, writes 1 MiB, `fork`s, and the child `exec`s. Which ideas (overcommit, demand paging, COW, OOM) apply at each step?
2. How do internal fragmentation in `malloc` and external fragmentation of physical RAM differ in who notices the waste?
3. When does the buddy allocator run, and when does a slab cache run, on a `fork` path?
4. How do thrashing and the OOM killer relate as two outcomes of memory pressure?
5. Why is "check `malloc` for `NULL`" necessary but not sufficient on Linux with overcommit?

#### Easy practical tasks

1. Write a one-page cheat sheet: internal/external fragmentation, bump, free list, buddy, slab, overcommit, OOM, demand paging, minor/major fault, COW, thrashing, working set.
2. Collect `getconf PAGE_SIZE`, `free -h`, `cat /proc/self/oom_score`, and `/usr/bin/time -v true` into a report.
3. Draw one figure: fault, COW split, and OOM as three exits from a missing frame.
4. Open `man 7 mmap` if present, else `man 2 mmap`. Write two failure modes (`MAP_FAILED` versus later SIGKILL).

#### Medium practical tasks

1. Write a lab notebook: one toy bump allocator, one COW `fork` demo, and one `time -v` fault count. Do not include an OOM test on a shared host.
2. Map each allocator type to one Linux component (user malloc, buddy, SLUB) in a table.
3. Use `vmstat` and `free` while you build a large project. Write whether you were CPU-bound, I/O-bound, or near memory pressure.

#### Advanced practical tasks

1. Read Linux `Documentation/admin-guide/mm/`. Write a one-page map from overcommit, reclaim, compaction, OOM, and swap.
2. Implement a tiny user-space allocator (bump + free list) and test it with a program that fragments on purpose. Report internal versus external waste.
