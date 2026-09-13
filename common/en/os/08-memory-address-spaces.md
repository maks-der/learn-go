# 8. Memory: Address Spaces

## Description

Virtual memory gives each process its own address space. This topic explains pages, page tables, the MMU, and the TLB. You learn the layout of text, BSS, heap, and stack. You learn the idea of `mmap` and the historical role of segmentation on x86.

Complete this topic after processes and hardware privilege. Complete this topic before allocation failure and protection bits. Those topics assume virtual addresses.

Use one term for each concept. A virtual address is the address that a user program uses. A physical address is the address on the RAM chips. A page is a fixed-size block of virtual memory. The MMU translates virtual addresses to physical addresses. Do not mix a page with a segment. Do not mix the TLB with the page table.

---

## Virtual memory

Virtual memory is the illusion that each process has a large private array of bytes. The program uses virtual addresses. The hardware and the kernel map those addresses to physical RAM, to disk (later topics), or to nothing (an unmapped hole).

Goals:

1. Isolation: process A cannot read process B memory by default.
2. Relocation: the same binary can run at the same virtual addresses in every process. The kernel places the pages in any free physical frames.
3. Size illusion: a process can use more virtual memory than the machine has RAM. The kernel can store some pages on a storage device (demand paging, topic 9).

Without virtual memory, programs would use physical addresses. Two programs could not both use address `0x400000`. The OS would have to relocate every binary at load time, and isolation would be weaker.

A virtual address space has mapped regions and holes. Access to a hole causes a fault. The kernel may send `SIGSEGV` to the process.

On 64-bit Linux, user space uses a large canonical address range. The kernel uses a separate range. A user program cannot map kernel addresses. Topic 1 already stated that split.

Virtual memory is not only "more memory than RAM". Even a small process uses virtual memory for isolation and for a simple pointer model.

### Questions

#### Theoretical questions

1. What is a virtual address?
2. What is a physical address?
3. Name three goals of virtual memory.
4. What happens when a process touches an unmapped address?
5. Why can two processes use the same virtual address for different data?

#### Easy practical tasks

1. Write five sentences about virtual memory. Use only facts from this section.
2. Run `cat /proc/self/maps | head`. Write two address ranges that you see.
3. Open `man 5 proc` for `/proc/<pid>/maps`. Write the meaning of the address column.
4. Make a table: "Virtual" and "Physical". Add three comparison rows.

#### Medium practical tasks

1. Run two `cat /proc/self/maps` in two shells. Compare the first text mapping. Write whether the virtual start looks the same (ASLR may change it; topic 10).
2. Draw one physical RAM bar and two process address spaces that both contain `0x400000`.
3. Write a C program that prints the address of `main` and of a local variable. State that those are virtual addresses.

#### Advanced practical tasks

1. Read the OSTEP virtual memory intro. Write a one-page summary of isolation, sharing, and the "address space" API.
2. Compare 32-bit user space (3G/1G split as a historical example) with 64-bit user space. Write why 32-bit processes hit a size limit sooner.

---

## Pages and page tables

The kernel does not map each byte. It maps pages. A page is a fixed size. On many Linux/x86-64 systems the base page is 4096 bytes. Huge pages are larger (2 MiB or 1 GiB) and are optional.

A physical page-sized chunk of RAM is a frame (or page frame). A mapping says: virtual page number V points to frame P, with permission bits (topic 10).

A page table is the data structure that stores those mappings. The CPU walks the page table on a translation (or uses the TLB). Linux uses multi-level page tables. On x86-64 a common walk has four or five levels. Each virtual address splits into indexes and a page offset.

The offset is the low 12 bits for a 4 KiB page. Those bits do not go through the table. They select the byte inside the page.

The kernel builds and updates page tables. User code cannot write them. A privileged instruction and the MMU use the table.

A process has its own page tables. After `fork`, the child gets a copy of the mappings. The kernel often uses copy-on-write (topic 9) so that the physical frames stay shared until a write.

`/proc/<pid>/maps` shows regions, not every page. `/proc/<pid>/pagemap` can show frame numbers (advanced, some rights needed).

If a mapping is missing, the MMU raises a page fault. The kernel handler decides: grow the stack, load a demand-paged file page, copy-on-write, or kill the process.

### Questions

#### Theoretical questions

1. What is a page?
2. What is a frame?
3. What is a page table?
4. Why does a 4 KiB page use a 12-bit offset?
5. Who is allowed to write page tables?

#### Easy practical tasks

1. Run `getconf PAGE_SIZE`. Write the page size.
2. Write a virtual address as "page number + offset" for page size 4096. Use an example address from `/proc/self/maps`.
3. Open `man 2 mmap` and find `MAP_ANONYMOUS`. Write that mappings are page-aligned.
4. Draw a two-level page table for a tiny address (indexes and offset).

#### Medium practical tasks

1. Read `cat /proc/self/maps`. Compute how many pages a 0x1000-length mapping uses.
2. Write a program that `malloc`s 16 KiB and prints the pointer. Align the size to the page size in a comment.
3. Search kernel documentation for "page table" levels on x86-64. Write the number of levels for a standard 4-level setup.

#### Advanced practical tasks

1. Read about 5-level paging (`la57`) in Linux docs. Write why more levels exist.
2. Use `page-types` or a documented pagemap example if you can. Write one virtual page to frame mapping. Stop if you lack rights.

---

## MMU

The memory management unit (MMU) is the hardware that translates virtual addresses to physical addresses. The MMU sits on the CPU (or next to it). Every user load and store goes through translation when paging is on.

The MMU reads page-table entries (PTEs). It checks present bits and permission bits. If the translation succeeds, the access goes to the physical address. If the translation fails, the CPU raises a page fault or a protection fault and enters the kernel.

The kernel programs the MMU:

- it sets the page-table base register (for example `CR3` on x86)
- it changes PTEs
- it tells the CPU to drop stale TLB entries after a change

Without an MMU, a general-purpose OS cannot give cheap isolation per process. Some microcontrollers have no MMU. They use a memory protection unit (MPU) or they run a single address space. This path assumes an MMU.

The MMU does not allocate RAM. The kernel allocator and the buddy system pick frames. The MMU only enforces the current tables.

Context switch to another process includes an MMU switch: a new page-table base. That is why a process switch is heavier than a thread switch in the same process (topic 6).

### Questions

#### Theoretical questions

1. What is the MMU?
2. What does the MMU do with a page-table entry?
3. What happens on a translation failure?
4. What is the page-table base register used for?
5. Why does a process context switch change MMU state?

#### Easy practical tasks

1. Write five sentences about the MMU. Use only facts from this section.
2. Make a table: "MMU" and "Kernel". Add rows for translate, allocate frame, handle fault.
3. Draw CPU, MMU, page tables in RAM, and physical RAM.
4. Open `man 2 mprotect`. Write how a program asks the kernel to change permission bits that the MMU will later enforce.

#### Medium practical tasks

1. Write a program that writes to a pointer `NULL`. Run it. Write that the MMU and kernel produced `SIGSEGV`.
2. Explain in six sentences why two threads of one process do not switch `CR3` on every thread switch.
3. Read a short x86 paging overview. Write the role of the present bit in a PTE.

#### Advanced practical tasks

1. Read kernel documentation on page-fault handling. Write the path from MMU fault to `do_page_fault` (or the current name) in five steps.
2. Compare MMU and MPU in a one-page table for embedded versus Linux servers.

---

## TLB

The translation lookaside buffer (TLB) is a cache of recent virtual-to-physical translations. A page-table walk is expensive (several memory reads). The TLB stores the result of a walk.

On a TLB hit, the MMU does not walk the full table. On a TLB miss, the hardware (or the software on some architectures) walks the table and fills the TLB.

The TLB is small. It cannot hold every mapping. A program with a large random working set can miss often. That cost is part of memory performance (topic 15 in other paths; here you only need the idea).

When the kernel changes a PTE, the TLB can hold a stale translation. The kernel must invalidate TLB entries. On x86, a write to `CR3` flushes the TLB for that CPU (with some exceptions for global pages). The kernel can also invalidate one page with a special instruction.

A process switch that changes page tables often flushes or tags the TLB. Some CPUs use PCID or ASID tags so that they do not flush every entry on every switch.

Do not mix the TLB with the CPU data cache. The data cache stores bytes of memory. The TLB stores translations. Both can miss after a context switch.

### Questions

#### Theoretical questions

1. What is the TLB?
2. What is a TLB hit and a TLB miss?
3. Why must the kernel invalidate TLB entries after a PTE change?
4. How does a process switch affect the TLB?
5. How does the TLB differ from the data cache?

#### Easy practical tasks

1. Write four sentences about the TLB. Use only facts from this section.
2. Make a table: "Hit" and "Miss". Add what the MMU does in each case.
3. Draw a TLB as a small table next to a large page table.
4. Search `man 1 perf` or `perf list` for a `dTLB` or `iTLB` event if `perf` exists. Write the event name.

#### Medium practical tasks

1. Explain why a tight loop on one array page is friendly to the TLB.
2. Explain why a walk over a huge linked list with nodes on many pages can miss the TLB.
3. Read about ASID or PCID in a CPU manual summary. Write one sentence on tagged TLB entries.

#### Advanced practical tasks

1. Use `perf stat` on a small program and a large random-access program if events work. Compare TLB-miss counts. Write a careful conclusion.
2. Read Linux documentation on TLB flush. Write why a shootdown on other CPUs is needed on an SMP machine.

---

## Stack, heap, BSS, text

A Unix process address space has standard regions. Names come from the ELF object model and from C.

- Text (code): executable instructions. Usually read-only and executable. Shared between processes that run the same file when the kernel can share the frames.
- Rodata: read-only data (string literals). Often next to text.
- Data: initialized global and static variables. Writable.
- BSS: global and static variables that start as zero. The file does not store zeros. The kernel maps anonymous zero pages.
- Heap: dynamic allocation (`malloc`). The heap grows toward higher addresses on many layouts. `brk`/`sbrk` used to grow it. Modern `malloc` also uses `mmap`.
- Stack: local variables, return addresses, and call frames. The stack grows toward lower addresses on most Linux ABIs. Each thread has a stack (topic 5).
- Memory mappings: extra regions from `mmap`, shared libraries (`.so`), and the vDSO.

`/proc/<pid>/maps` shows these as ranges with permissions `r-x`, `rw-`, `r--`. The pathname column shows the file or `[heap]`, `[stack]`, `[vdso]`.

The compiler and the linker assign symbols to sections. The kernel loads the ELF program headers, not the section names. You still use the names text, BSS, heap, and stack in conversation.

Do not store large objects on the stack. The stack has a limit (`ulimit -s`). A deep recursion or a large local array can fault. The heap or `mmap` is the place for large buffers.

### Questions

#### Theoretical questions

1. What does the text region contain?
2. What is BSS, and why is it not stored as zeros in the file?
3. How does the heap differ from the stack?
4. Which way does the stack grow on typical Linux?
5. What does `[vdso]` mean in `/proc/self/maps` at a high level?

#### Easy practical tasks

1. Run `cat /proc/self/maps`. Find `[stack]`, `[heap]` if present, and an `r-xp` mapping for the executable or `libc`.
2. Write a C program with a global `int g = 1;`, a global `int z;`, a `malloc(8)`, and a local `int l`. Print all four addresses.
3. Open `man 3 malloc` and `man 2 brk`. Write how the heap can grow.
4. Run `ulimit -s`. Write the stack size limit in kilobytes.

#### Medium practical tasks

1. Draw a map from low to high (or the order that your `/proc/self/maps` uses): text, data, BSS, heap, libraries, stack.
2. Compare `size a.out` after you compile a tiny program. Write the text, data, and BSS column values.
3. Create a large local array (several megabytes) and run it. Record a crash or a success. Then move the array to `malloc`. Write the difference.

#### Advanced practical tasks

1. Read `man 5 elf` program headers (`PT_LOAD`). Write how two `PT_LOAD` segments become maps.
2. Use `readelf -l a.out` and match virtual addresses to `/proc/<pid>/maps` while the program sleeps. Write the match.

---

## `mmap` idea

`mmap` is a system call that maps a region of virtual address space.

Two common uses:

1. File-backed map: the pages show the contents of a file. A read of the memory reads the file through the page cache. `MAP_SHARED` can write back to the file. `MAP_PRIVATE` uses copy-on-write for writes.
2. Anonymous map: no file. The pages are zero. `malloc` of large blocks often uses this. Thread stacks use anonymous maps.

```c
void *p = mmap(NULL, 4096, PROT_READ | PROT_WRITE,
    MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
```

`munmap` removes the mapping. `mprotect` changes permissions. `madvise` gives hints.

`mmap` is page-granular. The address and the length should follow page rules. The kernel rounds and aligns.

Compared with `read`: a file map can avoid an extra user buffer. Compared with `brk`: `mmap` can place regions anywhere in the holes of the address space and can return them independently.

Failures: `mmap` returns `MAP_FAILED`. A later access can still fault if you use `PROT_NONE` or if the file is truncated. Check the return value.

Topic 10 covers sharing and protection bits. This section only introduces the map as a first-class region in the address space.

### Questions

#### Theoretical questions

1. What does `mmap` do?
2. What is an anonymous mapping?
3. What is the difference between `MAP_SHARED` and `MAP_PRIVATE` for a file map?
4. Why must lengths follow page size?
5. How does a large `malloc` relate to `mmap`?

#### Easy practical tasks

1. Open `man 2 mmap`. Write the meaning of `MAP_ANONYMOUS` and `PROT_WRITE`.
2. Type the anonymous example. Print the pointer. `munmap` it.
3. Run `strace -e mmap,munmap` on that program. Write the arguments that you see.
4. Find `libc` mappings in `/proc/self/maps`. Write that they came from `mmap` of a file.

#### Medium practical tasks

1. `mmap` a small text file with `PROT_READ` and `MAP_PRIVATE`. Print the first bytes as a string if the file is text.
2. Draw three maps: text file, anonymous heap-like region, stack.
3. Compare `read` of a file into a buffer with `mmap` of the same file in a short table: copies, lifetime, errors.

#### Advanced practical tasks

1. Write a file with `MAP_SHARED`, change a byte in memory, `msync`, and confirm the file changed.
2. Read `man 2 mmap` on `MAP_FIXED`. Write why `MAP_FIXED` is dangerous for a beginner.

---

## Segmentation (historical / x86)

Segmentation divides memory into variable-sized segments. Each segment has a base, a limit, and permissions. A logical address is a segment selector plus an offset. The hardware adds the base and checks the limit.

Historical x86 (16-bit and early 32-bit) used segments for almost all accesses (`CS`, `DS`, `SS`, `ES`). Operating systems used segments to isolate processes before paging was common, or together with paging.

Modern 64-bit Linux on x86-64 does not use segmentation for user address spaces in the old way. Paging does isolation. Segment bases for `CS`, `DS`, and `SS` are unused in 64-bit mode, with small exceptions (`FS` and `GS` bases still hold TLS pointers).

Why this topic still names segmentation:

- textbooks contrast segments (variable size) with pages (fixed size)
- old x86 material will confuse you if you think Linux still gives each process a `DS` limit
- protection rings (topic 2) used to pair with segment descriptors

Other architectures never used x86-style segments. They used paging or an MPU.

Do not implement a segment table in a Linux user program. Use pages, `mmap`, and TLS via the compiler.

### Questions

#### Theoretical questions

1. What is a segment in the historical x86 model?
2. What is a logical address in that model?
3. Why did paging replace segmentation for isolation on modern Linux?
4. What leftover use do `FS` or `GS` have on x86-64 Linux?
5. How does a segment differ from a page?

#### Easy practical tasks

1. Write a two-column table: "Segment" and "Page". Add size, isolation method, and Linux x86-64 use.
2. Write four sentences about historical segmentation. Use only facts from this section.
3. Open a short x86-64 ABI note on TLS and `FS`. Write one sentence.
4. Draw a segment with base and limit, then a page table next to it.

#### Medium practical tasks

1. Read an OSTEP or textbook page that still teaches 8086 segments. Write three facts that do not apply to Linux x86-64 user space.
2. Explain why a variable-sized segment can cause external fragmentation of the address space more than fixed pages.
3. Find `fs` or TLS in `/proc/self/maps` or in `gdb` `info registers` if you use gdb. Write what you saw.

#### Advanced practical tasks

1. Read the Intel SDM overview of long mode segmentation. Write which segment checks still occur and which bases are ignored.
2. Compare x86 segments with PowerPC or ARM address-space IDs in a one-page survey (paging only).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of one user load: virtual address, TLB or page-table walk, MMU check, physical RAM.
2. How do text, heap, stack, and `mmap` regions appear in one address space without overlapping?
3. Why does the kernel, not the MMU, allocate frames?
4. When does a page fault mean "grow or fill a page", and when does it mean "kill the process"?
5. Why is segmentation in this handbook marked historical for Linux x86-64?

#### Easy practical tasks

1. Write a one-page cheat sheet: VA, PA, page, frame, PTE, MMU, TLB, text, BSS, heap, stack, `mmap`, segment.
2. Save `/proc/self/maps` and label five lines by hand.
3. Compile a program that prints `PAGE_SIZE`, a `mmap` pointer, and a stack address.
4. Open `man 7 address_spaces` if present, else `man 2 mmap`. Write two facts about Linux maps.

#### Medium practical tasks

1. Write a script that prints `getconf PAGE_SIZE`, `ulimit -s`, and the first and last lines of `/proc/self/maps`.
2. Draw one figure that includes the MMU, the TLB, a two-level table, and the ELF regions.
3. Trace a tiny program with `strace -e mmap,brk,mprotect`. Match new maps to `/proc/<pid>/maps` while it sleeps.

#### Advanced practical tasks

1. Read the kernel documentation on x86 paging. Write a one-page walk of a 4-level translation with example bit fields.
2. Compare `MAP_PRIVATE` file maps with anonymous maps and with `brk` heap in a report: who backs the pages, and what `fork` will share (preview of topic 9).
