# 5. Virtual Memory

## Description

Virtual memory gives each process its own address space. This topic explains virtual addresses, pages, page tables, the MMU, and the TLB. You learn the layout of text, BSS, heap, stack, and `mmap`. You also learn demand paging, copy-on-write after `fork`, thrashing, allocators, the Linux OOM killer, protection bits, shared memory, ASLR, and guard pages.

Complete this topic after processes and privilege modes. Complete this topic before files that use `mmap` in depth.

Use one term for each concept. A virtual address is the address that a user program uses. A physical address is the address on the RAM chips. A page is a fixed-size block of virtual memory. The MMU translates virtual addresses to physical addresses. Demand paging loads a page when the process first uses it. Copy-on-write shares a frame until a write. Do not mix a virtual address with a physical address. Do not mix the TLB with the page table.

---

## Virtual addresses, pages, page tables, MMU, TLB

Virtual memory is the illusion that each process has a large private array of bytes. The program uses virtual addresses. The hardware and the kernel map those addresses to physical RAM, to a storage device, or to nothing (an unmapped hole).

Goals:

1. Isolation: process A cannot read process B memory by default.
2. Relocation: the same binary can run at the same virtual addresses in every process. The kernel places the pages in any free physical frames.
3. Size illusion: a process can use more virtual memory than the machine has RAM. The kernel can store some pages on a storage device.

A virtual address space has mapped regions and holes. Access to a hole causes a fault. The kernel may send `SIGSEGV` to the process.

The kernel does not map each byte. It maps pages. On many Linux x86-64 systems the base page is 4096 bytes. Huge pages are larger (2 MiB or 1 GiB) and are optional. A physical page-sized chunk of RAM is a frame. A mapping says: virtual page number V points to frame P, with permission bits.

A page table is the data structure that stores those mappings. Linux uses multi-level page tables. On x86-64 a common walk has four or five levels. Each virtual address splits into indexes and a page offset. The offset is the low 12 bits for a 4 KiB page. Those bits do not go through the table. They select the byte inside the page.

The memory management unit (MMU) is the hardware that walks the page table (or uses the TLB) and enforces permissions. User code cannot write page tables. A privileged instruction and the MMU use the table.

The translation lookaside buffer (TLB) is a cache of recent translations. A hit avoids a full page-table walk. A context switch to another address space often flushes or tags the TLB. That cost is why a process switch is heavier than a thread switch in the same process (topic 3).

If a mapping is missing or the permission fails, the MMU raises a page fault. The kernel handler decides: grow the stack, load a demand-paged file page, copy-on-write, or kill the process.

`/proc/<pid>/maps` shows regions, not every page.

### Questions

#### Theoretical questions

1. What is a virtual address, and what is a physical address?
2. Name three goals of virtual memory.
3. What is a page table?
4. What work does the MMU do?
5. What is the TLB, and why does a process switch affect it?

#### Easy practical tasks

1. Run `cat /proc/self/maps | head`. Write two address ranges that you see.
2. Open `man 5 proc` for `/proc/<pid>/maps`. Write the meaning of the address column.
3. Write five sentences: virtual address, page, page table, MMU, TLB. Use only facts from this section.
4. Make a table: "Virtual" and "Physical". Add three comparison rows.

#### Medium practical tasks

1. Write a C program that prints the address of `main` and of a local variable. State that those are virtual addresses.
2. Draw one physical RAM bar and two process address spaces that both contain `0x400000`.
3. Draw a 4 KiB page split: virtual page number and 12-bit offset.

#### Advanced practical tasks

1. Read the OSTEP virtual memory intro. Write a one-page summary of isolation, sharing, and the address-space API.
2. Compare 32-bit user space (a historical 3G/1G split) with 64-bit user space. Write why 32-bit processes hit a size limit sooner.

---

## Stack, heap, BSS, text, `mmap`

A Linux user address space has named regions. `/proc/self/maps` shows them.

Text is the machine code. It is usually read-execute and not writable. It comes from the ELF file.

Data (sometimes called `.data`) holds initialized globals. BSS holds uninitialized globals. The kernel (or the loader) zeros BSS. BSS is writable.

The heap is the region that grows with `brk` or with `malloc` (the C library often uses `mmap` for large blocks). You allocate objects on the heap. You free them with `free`. The heap is shared among threads of the process (topic 3).

The stack holds frames for function calls: return addresses, saved registers, and local variables. The main thread stack grows down on many architectures. Each thread has its own stack.

`mmap` maps a file or anonymous memory into the address space. The call returns a pointer. `MAP_PRIVATE` plus a file is often copy-on-write. `MAP_SHARED` publishes writes to the file or to other processes (see the last section). `MAP_ANONYMOUS` is RAM that is not backed by a user file path.

The loader uses `mmap` for the program and for shared libraries. Your program can `mmap` a data file and change bytes in memory. Topic 6 treats the file as a path. This section treats the map as a region.

Typical layout (addresses are examples, ASLR shifts them): text at a low region, then data and BSS, then heap growing up, `mmap` regions in the middle, stack at a high region growing down. Holes sit between regions.

Do not return the address of a local stack variable from a function that has returned. That address is not valid. Do not `free` a stack address.

### Questions

#### Theoretical questions

1. What does the text region hold?
2. What is BSS?
3. How does the heap differ from the stack?
4. What does `mmap` do?
5. Why is a pointer to a returned local variable invalid?

#### Easy practical tasks

1. Run `cat /proc/self/maps`. Write lines that look like `[heap]`, `[stack]`, and a library path.
2. Open `man 2 mmap` and `man 3 malloc`. Write one sentence for each.
3. Write five sentences: text, BSS, heap, stack, `mmap`. Use only facts from this section.
4. Draw the regions from low to high with arrows for heap grow and stack grow.

#### Medium practical tasks

1. Write a C program that prints addresses of a function, a global, a `malloc` block, and a local. Match each to a region.
2. `mmap` a small file `MAP_SHARED`, change one byte, `msync` or close, and `cat` the file. Write what you saw.
3. Compare `malloc(16)` and `mmap` of 16 bytes (or 4 KiB). Write why the library does not `mmap` every tiny object.

#### Advanced practical tasks

1. Read `man 2 brk`. Write how `malloc` can use `brk` and `mmap` together (glibc overview).
2. Write a one-page note: what `[vdso]` and `[vvar]` mean in `maps` (high-level).

---

## Demand paging, COW after `fork`, thrashing

Demand paging means the kernel does not load every page when `exec` starts. It maps the file. When the process first touches a page, a fault occurs. The kernel reads that page from the file or zeros an anonymous page. Startup is faster. RSS (resident set size) grows with use.

Copy-on-write (COW) after `fork` shares physical frames between parent and child. Both page tables point at the same frames. The PTEs are read-only. When either process writes, a fault occurs. The kernel copies the frame, maps the copy writable for the writer, and leaves the other process on the old frame (or copies too if needed). `exec` in the child soon drops the old mappings. COW makes `fork` cheap when the child execs quickly.

Thrashing is a state where the machine spends most of its time paging. The working set of the running processes does not fit in RAM. The kernel writes pages to swap and reads them back again and again. CPU use looks busy or wait-heavy. Useful work is slow. Fixes: add RAM, reduce the working set, kill a hog, or move work to another machine. Swap is not free RAM.

`vmstat` and `si`/`so` columns show swap in and swap out. `free -h` shows swap use. Some swap is normal. Constant swap I/O is a problem.

A minor fault can be a COW or an anonymous page that was not yet backed. A major fault reads a storage device. `ps` and `/proc/<pid>/stat` expose fault counters.

### Questions

#### Theoretical questions

1. What is demand paging?
2. What does COW do after `fork`?
3. Why is `fork` plus a quick `exec` cheap with COW?
4. What is thrashing?
5. What is the difference between a minor fault and a major fault?

#### Easy practical tasks

1. Run `free -h`. Write total RAM and swap size.
2. Open `man 1 vmstat`. Write what `si` and `so` mean.
3. Write five sentences: demand paging, COW, thrashing. Use only facts from this section.
4. Run `ps -o pid,min_flt,maj_flt,cmd -p $$`. Write the two fault columns.

#### Medium practical tasks

1. Write a program that `malloc`s a large array but does not touch it. Compare `VmSize` and `VmRSS` in `/proc/self/status`. Then touch every page and compare again.
2. Draw parent and child after `fork`: shared frames, then a write in the child and a new frame.
3. Run `vmstat 1` for five seconds while the machine is idle. Write whether swap I/O is zero.

#### Advanced practical tasks

1. Read the OSTEP chapter on paging and TLBs. Write a one-page summary of demand paging and the working set.
2. Write a note: why a container memory limit can cause reclaim and kills even when the host has free RAM (preview of topic 10).

---

## Allocators and the OOM killer

An allocator hands out memory and takes it back. User `malloc` is a user-space allocator. The kernel page allocator hands out frames.

A bump (arena) allocator moves a pointer forward. It is fast. It does not free single objects, or it frees all at once.

A free-list allocator puts freed blocks on a list. The next allocation searches the list. Coalescing merges neighbor free blocks.

A buddy allocator uses sizes that are powers of two. A split turns a block into two buddies. Linux uses a buddy system for pages. Internal fragmentation occurs when you need 9 pages and you get 16.

A slab allocator keeps caches of equal-sized objects. Linux uses slab or SLUB for kernel objects such as inodes.

Internal fragmentation is waste inside an allocated block. External fragmentation is waste in holes that are too small for the next request. Pages reduce some external fragmentation of RAM because frames have one size.

The Linux out-of-memory (OOM) killer runs when the kernel cannot free enough memory. It picks a process (a score) and sends `SIGKILL`. The goal is to recover the system. The victim can be a large server, not the process that just called `malloc`. `dmesg` or `journalctl -k` can show an OOM dump.

`malloc` can return `NULL`. On Linux, overcommit can make `malloc` succeed and a later touch can kill the process. Check return values. Design for failure. Do not assume infinite RAM.

Do not expect `free` to return pages to the OS at once. The C library may keep them for the next `malloc`.

### Questions

#### Theoretical questions

1. How does a bump allocator work?
2. What is the difference between internal and external fragmentation?
3. What problem does a slab cache solve?
4. What does the OOM killer do?
5. Why can `malloc` succeed and a later write still kill the process on Linux?

#### Easy practical tasks

1. Open `man 3 malloc`. Write what happens when allocation fails (return value).
2. Write a four-row table: bump, free list, buddy, slab. Add one sentence each.
3. Write five sentences about allocators and the OOM killer. Use only facts from this section.
4. Run `cat /proc/meminfo | grep -E 'MemTotal|SwapTotal|Commit'`. Write the three lines.

#### Medium practical tasks

1. Write a program that `malloc`s 1 byte many times and prints `VmRSS`. Compare with one large `malloc`. This handbook does not contain the source.
2. Read `man 5 proc` for `/proc/sys/vm/overcommit_memory` (or the sysctl doc). Write the three mode numbers in your own words.
3. Search a past OOM line in `dmesg` if you have rights. If none exists, read a sample OOM dump online. Write which field is the victim name.

#### Advanced practical tasks

1. Implement a toy bump allocator over a fixed buffer (`alloc` only, then a reset). Do not replace libc in a real program.
2. Write a one-page note: how cgroup memory limits interact with the OOM killer (preview of topic 10).

---

## Protection bits, shared memory, ASLR, guard pages

Each mapping has permissions. Linux `mmap` and `mprotect` use `PROT_READ`, `PROT_WRITE`, and `PROT_EXEC`. `/proc/<pid>/maps` shows `r`, `w`, `x`.

The CPU and the MMU enforce the bits. A write to a read-only page faults. The kernel may handle the fault as COW, or it may send `SIGSEGV`. W^X (write xor execute) is a policy: a page must not be writable and executable at the same time. Modern Linux user stacks are not executable by default.

Shared memory is a design where two or more processes map the same physical frames and intend to see each other stores. Ways on Linux: `mmap` of the same file with `MAP_SHARED`, POSIX `shm_open`, System V `shmget`. Shared memory is fast IPC. Processes must synchronize (topic 4 and topic 9). Default `fork` private mappings are COW, not this design.

Address-space layout randomization (ASLR) shifts the base of stacks, heaps, and libraries. The goal is to make some exploits harder. You still see the real addresses in `/proc/self/maps` of your own process. Do not treat ASLR as a full security product. Topic 11 covers more isolation.

A guard page is an unmapped or `PROT_NONE` page next to a stack or a buffer. A touch faults. The kernel can grow a stack or send a signal. Guard pages catch some overflows. They do not catch all bugs.

`PROT_NONE` maps a page with no access. A touch faults.

User programs cannot bypass these bits. A kernel bug or a wrong `mprotect` in your own process can still weaken protection.

### Questions

#### Theoretical questions

1. What do R, W, and X mean on a page?
2. What is W^X?
3. How does shared memory differ from COW after `fork`?
4. What does ASLR randomize?
5. What is a guard page?

#### Easy practical tasks

1. Run `cat /proc/self/maps`. Write one line that is `r-x` and one line that is `rw-`.
2. Open `man 2 mprotect` and `man 3 shm_open`. Write one sentence for each.
3. Write five sentences: protection bits, shared memory, ASLR, guard pages. Use only facts from this section.
4. Make a table: "maps letter", "PROT flag", "typical region".

#### Medium practical tasks

1. `mmap` an anonymous page `PROT_READ`. Try to write it. Then `mprotect` to add `PROT_WRITE` and write again.
2. Run `cat /proc/self/maps` twice in two shells. Write whether stack or library bases differ (ASLR).
3. Draw two processes and one shared frame. Label a mutex that they must use.

#### Advanced practical tasks

1. Map POSIX shared memory in two processes. Write a number in one. Read it in the other. Add a semaphore or a process-shared mutex.
2. Write a one-page note: NX stack, ASLR, and W^X as a set. State what they do not stop.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the MMU, the page table, and protection bits form one isolation mechanism?
2. A process has a huge `VmSize` and a small `VmRSS`. Which ideas in this topic explain that pair?
3. How do demand paging, COW, and `exec` after `fork` fit one create-process story?
4. Why is the OOM killer a policy on top of virtual memory, not a replacement for `malloc` checks?
5. When do you use `mmap` for a file, and when do you use `read` into a heap buffer?

#### Easy practical tasks

1. Write a cheat sheet: virtual address, frame, PTE, MMU, TLB, text, BSS, heap, stack, `mmap`, demand paging, COW, thrashing, OOM, ASLR, guard page.
2. Run `cat /proc/self/status | grep -E 'VmSize|VmRSS|VmStk'`. Write the three values.
3. Draw fault paths: hole, COW write, demand file page, permission fail.
4. Open `man 5 proc` and bookmark the `maps` and `status` sections.

#### Medium practical tasks

1. Write a small program that `mmap`s a file, changes it, and also prints `/proc/self/maps` for that range. This handbook does not contain the source.
2. Use `time` and `/usr/bin/time -v` (if present) on a compiler. Write major page faults if the tool shows them.
3. Document a five-step response to thrashing: measure swap, find the hog, reduce load, check OOM logs, add RAM or limits.

#### Advanced practical tasks

1. Read the OSTEP paging and TLB chapters. Write a one-page lab plan: what you will measure on Linux (`maps`, faults, `vmstat`).
2. Compare POSIX shared memory and a `MAP_SHARED` file in a table: lifetime, name, and who cleans up.
