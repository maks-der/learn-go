# 10. Memory: Protection and Sharing

## Description

The MMU does not only translate addresses. It also enforces permissions. This topic explains read, write, and execute bits. You learn shared memory, copy-on-write pages, address-space layout randomization (ASLR), and guard pages.

Complete this topic after address spaces and allocation. You already know PTEs and COW after `fork`. This topic uses those ideas for security and for sharing.

Use one term for each concept. A protection bit is a permission in a page-table entry. Shared memory is a frame that more than one process maps writable or readable on purpose. ASLR is a random shift of map bases. A guard page is an unmapped or inaccessible page that catches overflow. Do not mix shared memory with COW private maps. Do not mix ASLR with a stack canary (another defense).

---

## Protection bits: R/W/X

Each mapping has permissions. Linux `mmap` and `mprotect` use `PROT_READ`, `PROT_WRITE`, and `PROT_EXEC`. `/proc/<pid>/maps` shows `r`, `w`, `x`.

Typical combinations:

- `r-x`: text, some libraries
- `r--`: constants, some file maps
- `rw-`: heap, stack, writable data
- `--x` or `r-x` without write: you cannot store into the page

The CPU and the MMU enforce the bits. A write to a read-only page faults. The kernel may handle the fault as COW, or it may send `SIGSEGV`. An execute of a non-executable page also faults (`SIGSEGV` or `SIGILL` depending on the path).

W^X (write xor execute) is a policy: a page should not be writable and executable at the same time. That policy makes some code-injection attacks harder. JIT compilers request `PROT_WRITE` first, write the code, then `mprotect` to `PROT_READ|PROT_EXEC`.

The NX (no-execute) bit in the PTE is the hardware execute disable. Modern Linux user stacks are not executable by default.

`PROT_NONE` maps a page with no access. Guard pages use this. A touch faults.

User programs cannot bypass these bits. A kernel bug or a wrong `mprotect` in your own process can still weaken protection. Do not `mprotect` a stack to executable without a strong reason.

### Questions

#### Theoretical questions

1. What do R, W, and X mean on a page?
2. What happens when a process writes to a read-only page that is not a COW page?
3. What is W^X?
4. Why is a non-executable stack useful?
5. What is `PROT_NONE`?

#### Easy practical tasks

1. Run `cat /proc/self/maps`. Write one line that is `r-x` and one line that is `rw-`.
2. Open `man 2 mprotect`. Write the three `PROT_` flags.
3. Write five sentences about protection bits. Use only facts from this section.
4. Make a table: "maps letter", "PROT flag", "typical region".

#### Medium practical tasks

1. `mmap` an anonymous page `PROT_READ`. Try to write it. Catch `SIGSEGV` or let it crash. Then `mprotect` to add `PROT_WRITE` and write again.
2. Draw a JIT path: RW, write bytes, then RX.
3. Compare `gcc -z noexecstack` documentation with your binary (`readelf -l a.out` GNU_STACK). Write whether the stack is executable.

#### Advanced practical tasks

1. Read about `PaX`/`NX` history at a high level. Write a one-page timeline of non-executable pages on Linux.
2. Use `mprotect` to remove execute from a page of your own code (dangerous). Document why the next call to that function fails. Do this only in a throwaway program.

---

## Shared memory

Shared memory is a design where two or more processes map the same physical frames and intend to see each other's stores.

Ways on Linux:

1. `mmap` of the same file with `MAP_SHARED`.
2. POSIX shared memory: `shm_open` then `mmap`.
3. System V shared memory: `shmget` and `shmat`.
4. Some anonymous maps with `MAP_SHARED` after `fork` (the child inherits the map).

Shared memory is fast IPC. There is no kernel copy on each store. Processes must synchronize. Use a mutex in the shared region (with `PTHREAD_PROCESS_SHARED`), a semaphore, or atomics. Topic 7 rules still apply. A race in shared memory is still a race.

Shared memory is not the default `fork` data mapping. Default private mappings are COW. Writes do not go to a sibling.

Permissions still apply. You can map shared memory read-only in one process and read-write in another if you set that up.

Cleanup: System V segments can stay after the process exits until `shmctl` removes them. POSIX shm objects live in `/dev/shm` until you `shm_unlink`. Leaks fill that filesystem.

Do not use shared memory for secrets without a plan. Any process with the same access can read the bytes.

### Questions

#### Theoretical questions

1. What is shared memory?
2. How does `MAP_SHARED` on a file differ from `MAP_PRIVATE`?
3. Why do you still need locks in shared memory?
4. How does inherited `MAP_SHARED` anonymous memory differ from COW after `fork`?
5. Why can System V shared memory leak after exit?

#### Easy practical tasks

1. Open `man 3 shm_open` and `man 2 mmap`. Write the two-step POSIX shm path.
2. Run `ls /dev/shm`. Write what you see.
3. Write five sentences about shared memory. Use only facts from this section.
4. Make a table: "COW private" and "MAP_SHARED". Add visibility of writes.

#### Medium practical tasks

1. Write two processes (or parent and child) that `mmap` the same file `MAP_SHARED`. One writes a byte. The other reads it after a `sleep` or a semaphore.
2. Draw two address spaces and one physical frame labeled "shared".
3. Compare POSIX shm and a `MAP_SHARED` file in a short table: namespace, persistence, unlink.

#### Advanced practical tasks

1. Put a `pthread_mutex_t` in shared memory with `PTHREAD_PROCESS_SHARED`. Two processes increment a counter. Document `pthread_mutexattr_setpshared`.
2. Read `man 7 shm_overview` and `man 7 sysv_shm`. Write when you pick each API.

---

## Copy-on-write pages

Copy-on-write pages are a sharing method that becomes private on write. Topic 9 explained COW after `fork`. This section places COW next to true shared memory.

COW uses:

- a shared frame
- read-only PTEs in each process
- a fault-and-copy on write

Uses:

- `fork` without an immediate full copy
- `MAP_PRIVATE` file maps: reads share the page cache; writes get a private frame
- the kernel can share identical file pages across processes that read the same library

COW is sharing of clean pages. It is not a protocol for two processes to cooperate on a buffer. After a write, the processes diverge.

The kernel tracks how many PTEs refer to a frame (a reference count). When the last mapping goes away, the kernel can free the frame. When a COW copy occurs, the writer gets a frame with count 1.

A write to a `MAP_PRIVATE` map does not change the file. A write to a `MAP_SHARED` map can change the file. That is the practical test.

COW faults can fail under OOM. A `fork` of a huge dirty heap is a risk. The child or parent can die on the first write to a COW page if no frame exists.

### Questions

#### Theoretical questions

1. What bits does a COW PTE use before the write?
2. When do two processes stop sharing a COW frame?
3. How does `MAP_PRIVATE` use COW for files?
4. Why is COW the wrong tool for a producer–consumer buffer between two processes?
5. What happens to the frame reference count after a successful COW copy?

#### Easy practical tasks

1. Write a two-column table: "COW" and "shared memory". Add three rows.
2. Draw the PTE change from read-only shared to two writable private frames (only the writer needs a new frame; the other can stay).
3. Open `man 2 mmap` and quote the `MAP_PRIVATE` idea in your own words.
4. Write four sentences that place this section next to topic 9. Use only new facts from this section.

#### Medium practical tasks

1. `mmap` a file `MAP_PRIVATE`, write the first byte in memory, and `cat` the file. Confirm the file is unchanged.
2. Repeat with `MAP_SHARED`. Confirm the file changes after `msync` or close as needed.
3. Explain why shared libraries can be one set of RX frames for many processes.

#### Advanced practical tasks

1. After `fork`, use `/proc/<pid>/smaps` (`Private_Dirty`, `Shared_Clean`) if present. Write how the numbers change after a write in the child.
2. Read kernel documentation on `do_wp_page` (write-protect fault). Write a five-step handler list.

---

## ASLR

Address-space layout randomization (ASLR) is a defense. The kernel places stacks, heaps, libraries, and (when possible) the main executable at random offsets each run.

Goal: an attacker who knows a bug (for example a buffer overflow) should not know a fixed address of a gadget or of the stack. The attacker must leak an address first.

Linux controls ASLR with `kernel.randomize_va_space` (`sysctl`). Common values:

- `0`: off
- `1`: conservative randomization
- `2`: full randomization (typical default)

The main program is randomized when it is a position-independent executable (PIE). Most modern distributions build PIE by default. `cat /proc/self/maps` on two runs of the same program shows different starts when ASLR is on.

ASLR is not complete protection. A leak of one library address can reveal the rest of that library. Partial overwrites and heap sprays still exist. Combine ASLR with NX, stack canaries, and safe code.

Do not turn ASLR off on a production host to "debug" unless you use a controlled environment. For a local debug session, `gdb` and `set disable-randomization` exist. Prefer logging addresses from `/proc/self/maps` while ASLR stays on.

`mmap` without a hint still receives a kernel-chosen address. That choice can include randomization.

### Questions

#### Theoretical questions

1. What is ASLR?
2. What does ASLR make harder for an attacker?
3. Why does the main binary need PIE for full ASLR of text?
4. Why is ASLR not enough by itself?
5. What is a typical default of `kernel.randomize_va_space` on a desktop Linux?

#### Easy practical tasks

1. Run the same program twice and compare `/proc/<pid>/maps` start addresses for `[stack]` or `libc`. Write whether they moved.
2. Run `sysctl kernel.randomize_va_space` if allowed. Write the value.
3. Open `man 8 sysctl` and write how a sysctl name maps to `/proc/sys`.
4. Write five sentences about ASLR. Use only facts from this section.

#### Medium practical tasks

1. Compile with and without PIE if your toolchain allows (`-fno-pie -no-pie` versus default). Compare the text address across two runs each. Write the difference.
2. Draw an address space with arrows that show which regions move under ASLR.
3. Read a short note on ASLR bypass via an information leak. Write why one leaked `libc` pointer is valuable.

#### Advanced practical tasks

1. In a VM, read the kernel documentation for `randomize_va_space`. Write what each value randomizes. Do not set `0` on a machine that you share.
2. Compare Linux ASLR with Windows ASLR at a high level (public docs). Write two similarities and two differences.

---

## Guard pages

A guard page is a page that must not be accessed. It sits next to a buffer or a stack. An overflow that walks off the buffer hits the guard and faults.

Linux thread stacks have a guard region. A overflow of the stack hits `PROT_NONE` or an unmapped gap. The process receives a signal instead of silent corruption of the heap or of another thread stack.

You can add your own guard:

```c
/* allocate 3 pages, make the first PROT_NONE, use the middle page */
```

`mmap` a larger region and `mprotect` the edges to `PROT_NONE`. Your object sits in the middle.

Hardware stack overflow and software heap overflow are different. `malloc` does not place a guard page after every 16-byte object. Heap overflow can still smash the next chunk metadata. Tools such as AddressSanitizer add red zones. Those red zones are a debug allocator feature, not the kernel stack guard.

`MAP_GROWSDOWN` and stack expansion have special rules. A large stack skip can miss the guard and fault anyway, or it can be seen as an attack. Do not allocate huge arrays on the stack.

Guard pages use page granularity. A 1-byte overflow of a small heap object does not hit a guard page unless you aligned the object at the end of a page and placed a guard next to it.

### Questions

#### Theoretical questions

1. What is a guard page?
2. Why do thread stacks use a guard region?
3. Why can `malloc` of a small object not rely on a kernel guard page?
4. How do you build a manual guard with `mmap` and `mprotect`?
5. Why is the guard size at least one page?

#### Easy practical tasks

1. Run `ulimit -s` and `cat /proc/self/maps | grep -i stack`. Write the stack range and the limit.
2. Open `man 3 pthread_attr_setguardsize`. Write the purpose of the call.
3. Write four sentences about guard pages. Use only facts from this section.
4. Draw a stack that grows down into a `PROT_NONE` page.

#### Medium practical tasks

1. `mmap` three pages. `mprotect` page 0 to `PROT_NONE`. Use page 1. Touch page 0 and record the signal.
2. Write a recursive function that eventually overflows the stack (in a VM). Write whether you see a crash. Do not run this on a machine that you cannot reboot.
3. Compare a stack guard with AddressSanitizer red zones in a short table: who inserts them, granularity, production use.

#### Advanced practical tasks

1. Read `man 2 mmap` on `MAP_GROWSDOWN`. Write why a user program should not invent its own stack that way.
2. Measure the default pthread guard size. Change it with `pthread_attr_setguardsize`. Overflow a thread stack in a test program and document the result.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do R/W/X bits, COW, and shared memory all use the same PTE machinery for different goals?
2. Which defenses in this topic help against code injection, and which help against stack smash corruption?
3. When do two processes see the same write, and when do they get a private copy?
4. Why does ASLR move libraries while protection bits stay the same?
5. How does a guard page use `PROT_NONE` differently from a read-only text page?

#### Easy practical tasks

1. Write a one-page cheat sheet: PROT bits, W^X, NX, MAP_SHARED, POSIX shm, COW, ASLR, PIE, guard page.
2. Save two `/proc/self/maps` dumps from two process starts. Highlight one moved address and one permission string.
3. Draw one frame: two processes, first COW, then one process `mprotect` or write, then a third path that is `MAP_SHARED`.
4. Open `man 7 pkeys` if present (memory protection keys). Write one sentence as a preview, or write that the page is missing.

#### Medium practical tasks

1. Build a mini-lab: (1) SIGSEGV on write to RX page, (2) shared file byte, (3) private file byte, (4) guard fault. Write a result table.
2. Document how you would share a ring buffer between two processes: mapping type, mutex location, and why COW is wrong.
3. Use `readelf -l` and `maps` to show PIE and NX stack for one binary.

#### Advanced practical tasks

1. Read Linux `Documentation/admin-guide/sysctl/kernel.rst` on `randomize_va_space` and a pkeys overview. Write a one-page "defense in depth" for process memory.
2. Implement POSIX shm with a process-shared mutex and a guard page at each end of the buffer. Document faults if a neighbor walk goes too far.
