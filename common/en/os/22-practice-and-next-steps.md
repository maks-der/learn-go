# 22. Practice and Next Steps

## Description

This topic turns the path into work that you write and run. You plan a tiny shell, a race that you can see and then fix, a toy allocator, and an `strace` of a real command. You also plan how to read [OSTEP](https://ostep.org/) next to each earlier topic. The last section names xv6, OS course labs, and internships.

Complete this topic after the advanced lock topic. You already have topics 1 to 21. This topic does not add a new kernel subsystem. It asks you to connect the subsystems.

Use one term for each concept. A tiny shell is a loop that reads a line, splits it, and runs a program with `fork`, `exec`, and `wait`. A race is an unsynchronized concurrent access that you can demonstrate. A toy allocator is a user-space `malloc`/`free` over a fixed buffer. A lab is a graded or self-directed exercise with a clear pass test. Do not treat this file as a solution key. Do not paste a full shell or a full allocator from the internet into your tree and call it practice.

---

## Write a tiny shell (`fork`/`exec`)

A tiny shell is the best single program for topics 4 and 11. You reuse one process as the parent. Each command runs in a child.

Minimum behavior:

1. Print a prompt.
2. Read one line from stdin.
3. Split the line into a command and arguments (spaces are enough).
4. `fork`.
5. The child calls `execvp` (or `execve` with a path search that you write later).
6. The parent calls `waitpid` and then prints the prompt again.
7. Exit the loop on end-of-file or on a `exit` built-in.

In scope for a first week: external programs (`ls`, `date`), a working directory that the child inherits, and a non-zero exit status that you print.

Out of scope until the minimum works: pipes, redirects, job control, quotes, and glob. Those features are extra labs. Pipes need topic 15. Redirects need `dup2` and `open`.

Checks that you write yourself:

- `waitpid` reaps the child. `ps` shows no zombie.
- `exec` failure in the child prints an error and `_exit`s. The parent must not `exec` by mistake.
- A built-in `cd` must run in the parent. A `cd` in the child does not change the shell.

Topic 4 already defined `fork`, `exec`, and zombies. This section only states the product.

Do not start from a 2000-line GitHub shell. Write the loop first. Add one feature at a time.

### Questions

#### Theoretical questions

1. Why must `cd` run in the parent?
2. What does the child do if `execvp` returns?
3. Why does the parent call `waitpid`?
4. Which topic 4 bug appears if the parent never waits?
5. Why are pipes an extra lab, not the first lab?

#### Easy practical tasks

1. Open `man 2 fork`, `man 3 execvp`, and `man 2 waitpid`. Write one sentence for each.
2. Write a ten-line design: prompt, parse, fork, exec, wait. No code yet.
3. Run `ps -o pid,ppid,stat,cmd` while a real shell runs `sleep 30`. Write the parent and child.
4. List five commands that your tiny shell must be able to run when the minimum is done.

#### Medium practical tasks

1. Write the minimum shell. Test `ls`, `true`, `false`, and a bad name. Record exit statuses. This handbook does not contain the source.
2. Add a built-in `cd` and `pwd`. Document one failed `cd` (bad path).
3. Use `strace -f -e fork,clone,execve,wait4,waitpid` on your shell for one `ls`. Write the child PID story.

#### Advanced practical tasks

1. Add one extra feature only after the minimum is stable: `|` between two commands, or `>` redirect. Write a design of FDs before you code.
2. Write a one-page test list: zombies, `exec` fail, `cd`, EOF, and a long argument line.

---

## Write a multithreaded program and hit a race; then fix it

A race that you cannot see does not teach you. Build a program that is wrong on purpose, show a bad result, then add a lock.

A standard lab:

1. Start `N` threads (`N` at least 4).
2. Each thread adds 1 to a shared `int` (or `long`) `ITERS` times (`ITERS` at least 100000).
3. Join all threads.
4. Print the counter. The correct value is `N * ITERS`.

Without a lock, the value is often smaller. That is the lost update from topic 5 and topic 7. If your machine never shows a loss, raise `ITERS`, disable the optimizer on that file, or run more threads. A pass without a loss is not a proof of safety.

Then fix the program with `pthread_mutex_lock` around the increment, or with `atomic_fetch_add`. Print the correct total. Measure time with and without the lock (`clock_gettime` monotonic, topic 19).

Optional second race: a linked-list push from two threads without a lock. The list can lose a node or corrupt a pointer. Repair it with one mutex around push.

Tools: ThreadSanitizer (`-fsanitize=thread`) on the broken binary. Write the report in your own words (topic 7).

Do not "fix" the race with `sleep`. Do not use `volatile` as the fix (topic 7).

### Questions

#### Theoretical questions

1. What is the correct counter value in the lab?
2. Why can a fast machine hide the race?
3. Why is `sleep` not a fix?
4. What two correct fixes does this section allow?
5. What does ThreadSanitizer add that a wrong total already showed?

#### Easy practical tasks

1. Open `man 3 pthread_create` and `man 3 pthread_mutex_lock`. Write one sentence for each.
2. Write the expected total for 8 threads and 1 000 000 iterations.
3. Write five sentences about the lab. Use only facts from this section.
4. Make a table: broken increment, mutex increment, atomic increment. Add "exclusion".

#### Medium practical tasks

1. Implement the broken counter. Save a wrong total. Then implement the mutex fix. Save a right total. Do not put the solution in this repository handbook.
2. Run the broken program under ThreadSanitizer. Write three lines of the report in your own words.
3. Time the mutex version and the atomic version with a monotonic clock. Write the two durations.

#### Advanced practical tasks

1. Add a broken list push. Capture a crash or a lost node. Then fix it. Document the critical section.
2. Write a one-page note that maps this lab to topic 7 (critical section, visibility) and topic 21 (futex under the mutex).

---

## Implement a user-space allocator toy

A toy allocator sits in one process. It is not the kernel buddy allocator (topic 9). You request a large arena with `mmap` or with a static array. You then implement `toy_malloc` and `toy_free`.

Goals of a first toy:

1. Split the arena into blocks. Each block has a header (size, free flag, maybe a next pointer).
2. `toy_malloc(n)` finds a free block that is large enough (first fit is enough).
3. You may split a large block.
4. `toy_free(p)` marks the block free. You may coalesce with a neighbor.
5. Alignment: return pointers that are aligned for an 8-byte or 16-byte object.

Pass tests that you write:

- Allocate 10 objects, write a pattern, free them, allocate again.
- A double-free detector (a flag or a simple abort) is a plus.
- External fragmentation: many small frees, then one large alloc that fails. That failure is a lesson, not a crash of the process if you check NULL.

Out of scope: multiple threads (add a mutex only if you want a second week), `mmap` of each object, and glibc compatibility of every edge.

Topic 9 named bump, free list, buddy, and slab. A bump allocator that never frees is a valid day-one toy. A free list is the week-one toy. Buddy is extra.

Do not replace `malloc` in a whole process with `LD_PRELOAD` until the toy passes tests on its own API.

### Questions

#### Theoretical questions

1. Where does the arena come from?
2. What does a block header store?
3. What is first fit?
4. How do you demonstrate external fragmentation?
5. Why is a bump allocator a valid first toy?

#### Easy practical tasks

1. Open `man 2 mmap` and `man 3 malloc`. Write how the toy differs from libc `malloc`.
2. Draw an arena with three blocks: two used, one free.
3. Write five sentences about the toy. Use only facts from this section.
4. Make a table: bump, free list, buddy. Add "can free?" from topic 9.

#### Medium practical tasks

1. Implement a bump allocator over 64 KiB. Show that `toy_free` is a no-op or is unsupported. Then implement a free list in a second file.
2. Write tests that fill the arena and check a NULL return.
3. Use `/proc/self/maps` to find your `mmap` arena. Write the range.

#### Advanced practical tasks

1. Add coalesce on `toy_free`. Write a test: free two neighbors, then alloc their sum.
2. Write a one-page comparison of your free list with glibc `malloc` from public docs: threads, and mmap of large blocks.

---

## Use `strace` on a real command

Topic 20 introduced `strace`. This lab is a full story of one real command. Pick `cat FILE`, `ls DIR`, or `uname -a`. Do not start with a browser.

Method:

1. Write the command that you run.
2. Run `strace -o /tmp/os-strace.txt -f -tt command ...`
3. Read the file with a pager.
4. List the first `execve`.
5. List `openat` (or `open`) of the target file or directory.
6. Count `read` and `write` if the command copies bytes.
7. Note `mmap`, `brk`, or both (loader and libc).
8. Note the final `exit_group`.

Write a one-page report: which calls match topics 4, 8, and 11. Which calls you do not know yet. Look those names up in `man 2`.

Flags that help: `-c` for counts, `-e trace=file` for path calls, `-s 64` for buffer previews. Do not publish buffers that contain secrets.

`strace -f` follows children. A shell pipeline needs `-f`. Your tiny shell is a good second target.

A command that talks to the network shows `socket` and `connect` (topic 18). That is a bonus report, not the first one.

### Questions

#### Theoretical questions

1. Why is `execve` first in the trace of a new command?
2. Why can `cat` show many `mmap` calls before the file `read`?
3. What does `-f` change?
4. Why must you not publish a raw trace from a login tool?
5. Which handbook topics should appear in a `cat` report?

#### Easy practical tasks

1. Open `man 1 strace`. Write the meaning of `-o`, `-c`, and `-e`.
2. Run `strace -c cat /etc/hostname`. Write the three calls with the highest counts.
3. Write five sentences about this lab. Use only facts from this section.
4. Make a table: `execve`, `openat`, `read`, `write`, `exit_group`. Add one topic number each.

#### Medium practical tasks

1. Produce the one-page `cat` report that this section describes. Save the trace in your practice folder.
2. Compare `strace cat FILE` with `strace -e trace=file cat FILE`. Write what the filter hides.
3. Trace your tiny shell that runs `true`. Write how `-f` shows the child.

#### Advanced practical tasks

1. Trace `ls --color=never` on a small directory. Map `getdents64` (or `readdir` in libc) to topic 11 directories.
2. Write a second report on `curl` of a URL only if you already started `net.topics.md`. Otherwise stop at `socket` and `connect` names.

---

## Read [OSTEP](https://ostep.org/) chapters that match each section

[Operating Systems: Three Easy Pieces](https://ostep.org/) is a free textbook. It groups the field into virtualization, concurrency, and persistence. This learning path used the same ideas in a Unix handbook order.

How to read:

1. Finish a handbook topic (or a cluster of topics).
2. Open the OSTEP chapter list.
3. Read the chapter that matches. Examples: processes with topic 4, paging with topic 8, locks with topic 7, files and crash consistency with topics 11–13, I/O with topics 13–14.
4. Write a one-page map: OSTEP section titles to handbook headings. Use your own words.
5. Do the book's homework only when you want extra work. This repository does not ship those answers.

OSTEP is not Linux-specific. The book uses abstract machines and some xv6 or C examples. When the book name and the Linux name differ, keep both in your map (topic 1 already said APIs differ).

Read again after topic 21. Concurrency chapters then make more sense.

Do not try to read the whole book in one weekend before topic 4. Pair it.

Official site: [https://ostep.org/](https://ostep.org/).

### Questions

#### Theoretical questions

1. What three pieces does OSTEP name?
2. When do you open an OSTEP chapter relative to a handbook topic?
3. Why can an OSTEP name differ from a Linux name?
4. Why pair the book instead of reading it all first?
5. What belongs in a one-page map?

#### Easy practical tasks

1. Open [https://ostep.org/](https://ostep.org/). Write five chapter titles that match topics 4, 6, 7, 8, and 11.
2. Write five sentences about how this path uses OSTEP. Use only facts from this section.
3. Bookmark the PDF or the chapter HTML that you use.
4. Make a table: handbook topic number, OSTEP chapter name (your match).

#### Medium practical tasks

1. Read one OSTEP process or address-space chapter. Write a one-page map to this repository's matching handbook.
2. Read one OSTEP concurrency chapter after you finish the race lab. Add three sentences that you could not have written before the lab.
3. Compare OSTEP persistence (files, disks) with topics 12 and 13 in a six-row table.

#### Advanced practical tasks

1. Build a full index: topics 1–21 versus OSTEP chapters. Mark gaps (namespaces, eBPF) that the book does not cover.
2. Write a study calendar: two handbook topics and one OSTEP chapter per week.

---

## Then: Linux kernel internships / xv6 / operating systems course labs

After this path you can choose a harder environment. Three common next steps exist.

xv6 is a small teaching OS from MIT. The book and the source are public: [https://pdos.csail.mit.edu/6.1810/2024/xv6/book-riscv-rev4.pdf](https://pdos.csail.mit.edu/6.1810/2024/xv6/book-riscv-rev4.pdf). xv6 is small enough to read. Labs ask you to add system calls, paging features, or locks. Use xv6 when you want a full kernel that you can change.

University OS labs (MIT 6.1810, Stanford, or a local course) give a schedule, a test suite, and a mentor. Use them when you want deadlines. The labs often use xv6 or a similar kernel.

Linux kernel internships and outreach programs expect more. You need C, git, mailing-list habits, and one area (driver, filesystem, scheduler, or net). Start with `Documentation/`, a tiny cleanup that a maintainer asked for, and a VM. Topic 20 taught you how to read a path. Topic 17 taught you that a bad patch can panic a machine. Do not send a huge unsolicited rewrite.

Also useful:

- The Linux Programming Interface (Kerrisk) for POSIX depth
- kernel.org documentation
- `net.topics.md` if you build servers
- The suggested practice order in `os.topics.md` (commands, `mmap`, ring buffer, namespaces, FUSE toy)

Pick one next step. Finish it. Do not start xv6, a university lab, and a kernel internship in the same week.

### Questions

#### Theoretical questions

1. Why is xv6 easier to change than Linux?
2. What does a university OS lab add that this handbook does not add?
3. What skills does a Linux kernel internship expect beyond this path?
4. Why is a tiny documented patch better than a huge rewrite?
5. When do you open `net.topics.md` instead of xv6?

#### Easy practical tasks

1. Open the xv6 book URL. Write the table of contents section names that match processes and paging.
2. Open `os.topics.md` and copy the ten-item suggested practice order into your notes. Tick what you already did.
3. Write five sentences about next steps. Use only facts from this section.
4. Bookmark xv6, OSTEP, kernel.org docs, and the Linux man-pages.

#### Medium practical tasks

1. Write a one-page plan: xv6 or a local course or a Linux starter area. Give three reasons.
2. Clone xv6 or open the source online. Find `fork` or the trap entry. Write the file name. Do not submit a lab solution here.
3. List the remaining items from the suggested practice order in `os.topics.md`. Schedule three of them.

#### Advanced practical tasks

1. Complete one official xv6 or course lab on your own machine. Keep the solution in your private repo. This handbook stays without solutions.
2. Write a six-month plan: finish `net.topics.md` topics 1–10, read three OSTEP chapters that you skipped, and either xv6 lab 1 or a documented kernel newbie task.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the tiny shell, the race lab, and the toy allocator each prove a different cluster of topics (processes, concurrency, memory)?
2. Why does an `strace` of `cat` plus an OSTEP files chapter beat either activity alone?
3. Which next step do you pick if you want to change a kernel, and which next step do you pick if you want POSIX depth without a kernel tree?
4. How do you keep this file's "no solutions" rule when you look at a classmate's shell?
5. When is the suggested practice order in `os.topics.md` a better checklist than starting xv6?

#### Easy practical tasks

1. Write a one-page cheat sheet: tiny shell steps, race lab numbers, allocator header fields, `strace` report outline, OSTEP pairing rule, xv6 versus Linux internship.
2. Create a folder `os-practice` with four empty notes: `shell.md`, `race.md`, `alloc.md`, `strace.md`. Write one pass test in each note.
3. Draw a roadmap from topic 1 to topic 22 with four practice boxes on it.
4. Bookmark OSTEP, xv6, `os.topics.md`, and `net.topics.md`.

#### Medium practical tasks

1. Fill the four notes with results as you finish each lab. Each note must contain commands and observations, not a pasted internet solution.
2. Run the first four items of the suggested practice order in `os.topics.md` if you have not run them. Write one sentence each.
3. Use `strace` on your shell, your race program, and a libc `malloc` test. Write one new fact per binary.

#### Advanced practical tasks

1. After the four labs, write a two-page self-review: which topics 1–21 you still cannot explain, and which OSTEP chapters you reread.
2. Start exactly one next-step track (xv6 lab 1, a course lab, or `net.topics.md` topic 1). Write the first week's log. Do not put official lab solutions into `common/en/os`.
