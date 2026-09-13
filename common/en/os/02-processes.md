# 2. Processes

## Description

A process is a running instance of a program. This topic explains the difference between a program, a process, and a thread. You learn process IDs, parent and child, and the process control block idea. You also learn the Unix calls `fork`, `exec`, and `wait`, plus zombies, orphans, process states, and environment variables.

Complete this topic after topic 1. Complete this topic before you study threads and scheduling. Those topics assume a process.

Use one term for each concept. A program is a file of instructions and data. A process is the OS object that runs a program. A thread is an execution context inside a process. A zombie is a process that has exited and that the parent has not yet waited for. An orphan is a process whose parent has exited. Do not mix `fork` with `exec`. Do not mix a zombie with an orphan.

---

## Program vs process vs thread

A program is a passive file. On Linux the file is often an ELF binary or a script with an interpreter. The program does not run until the kernel loads it.

A process is the active OS object. The process has:

- a process ID
- an address space (virtual memory)
- a list of open file descriptors
- credentials (user ID, group ID)
- at least one thread

One program can have many processes. Two users can run `/bin/bash` at the same time. Each shell is a different process. One process can also load a new program over itself with `exec`. The process ID can stay the same after `exec`.

A thread is the unit that the CPU runs. A thread has a stack and a set of registers. A process starts with one thread. The process can create more threads. All threads in one process share the same address space and the same file descriptor table. Topic 3 covers threads.

The program is the stored instructions. The process is the running instance. Each thread is one execution of those instructions inside the same process.

The kernel tracks processes. User tools such as `ps` list them. A process can create child processes. A process can be stopped, continued, or killed.

### Questions

#### Theoretical questions

1. What is the difference between a program and a process?
2. Can two processes run the same program file? Explain.
3. What is a thread in relation to a process?
4. What state does a process have that a program file does not have?
5. What happens to the process ID when the process calls `exec`?

#### Easy practical tasks

1. Run `ps -o pid,cmd`. Find your shell. Write the PID and the program path.
2. Start a second `sleep 60` in another terminal. Run `ps -C sleep -o pid,cmd`. Write how many processes share the same program name.
3. Write six sentences: two about programs, two about processes, two about threads. Use only facts from this section.
4. Open `man 1 ps`. Write the meaning of the `CMD` column.

#### Medium practical tasks

1. Run `/bin/true` twice. Then run `pgrep -a bash` or `pidof bash`. Write how you distinguish the processes.
2. Draw a diagram: one ELF file on disk, two processes, one of those processes with two threads.
3. Use `readlink /proc/$$/exe` and `cat /proc/$$/cmdline`. Write the difference between the program file and the command line.

#### Advanced practical tasks

1. Compile a tiny C program. Run it three times at once (`./a.out &` three times). Show three PIDs and one path from `/proc/<pid>/exe`.
2. Read `man 5 elf` or an ELF overview. Write which parts of the file become the process text and data after `exec`.

---

## PID, parent, child, and the PCB idea

The kernel gives each process a process ID (PID). The PID is an integer. Tools such as `kill` and `ps` use it. On Linux you also see a thread ID (TID). For a single-thread process, PID and TID often match. Topic 3 covers TIDs.

A process that creates another process is the parent. The new process is the child. The child has a parent PID (PPID). On Linux, `ps -o pid,ppid,cmd` shows both. Process ID 1 (`init` or `systemd`) is the ancestor of user processes after boot.

The kernel must store the state of each process. The textbook name for that data is the process control block (PCB). Linux does not use the name PCB in the source. Linux uses `task_struct` for a thread and related structures for the address space and files. The idea is the same.

A PCB holds at least:

- identifiers: PID, parent PID
- state: running, ready, blocked, zombie
- CPU context: registers when the process is not on the CPU
- memory description: page tables or a pointer to them
- file table: open file descriptors
- credentials: user, group, capabilities
- accounting: CPU time, nice value
- signal state

When the scheduler picks a process, the kernel loads the CPU context from the PCB. When the process blocks or when its time slice ends, the kernel stores the context back.

You do not allocate a PCB in your C program. The kernel allocates it when it creates the process. You see pieces of the PCB through `/proc/<pid>`.

A PID can be reused after the process is fully reaped. Do not store an old PID and assume it still names the same process after a long time.

### Questions

#### Theoretical questions

1. What is a PID?
2. What is a parent process and what is a child process?
3. What is the PCB idea?
4. Name five kinds of data that a PCB holds.
5. Why can a PID be reused later?

#### Easy practical tasks

1. Run `ps -o pid,ppid,cmd`. Write your shell PID and its parent PID.
2. Open `man 5 proc` and find `/proc/<pid>`. Write three file names under that directory.
3. Write five sentences about PID, parent, child, and the PCB. Use only facts from this section.
4. Run `cat /proc/self/stat | awk '{print $1, $4}'`. Write that the first field is PID and the fourth is PPID (see `man 5 proc`).

#### Medium practical tasks

1. Draw a tree: PID 1, your login session, your shell, and one `sleep` child.
2. Compare `ps -p 1 -o pid,ppid,cmd` with `ps -p $$ -o pid,ppid,cmd`. Write who is the ancestor.
3. List six `/proc/<pid>` files and one PCB field that each file exposes.

#### Advanced practical tasks

1. Read a short note on Linux `task_struct` (kernel documentation or a textbook figure). Write how the textbook PCB maps to Linux names.
2. Write a one-page report: why a supervisor that stores PIDs must also check start time or a pidfd on modern Linux.

---

## `fork`, `exec`, `wait`

Unix creates a process with `fork` (or `clone` on Linux). `fork` copies the calling process. The parent gets the child PID as the return value. The child gets 0. Both continue from the same instruction after `fork`.

`exec` replaces the address space of the calling process with a new program. The PID stays the same. Open file descriptors stay open unless the close-on-exec flag is set. The stack, heap, and mappings of the old program are gone.

A typical child path is `fork`, then `exec` in the child. A typical parent path is `fork`, then `wait` or `waitpid` for that child.

`wait` and `waitpid` let the parent collect the child exit status. The kernel keeps a zombie until the parent waits. Topic 14 asks you to write a tiny shell with this pattern.

`fork` is expensive if you copy the full address space at once. Linux uses copy-on-write pages (topic 5). The child shares physical frames until a write.

`vfork` and `posix_spawn` exist. For this path, learn `fork` plus `exec` first.

Errors: `fork` can fail (`EAGAIN`, `ENOMEM`). `exec` can fail if the file is missing or not executable. After a failed `exec`, the child still runs the old program. The child must `_exit` after a failed `exec`. If the child returns to the parent code, you get two shells or two servers.

### Questions

#### Theoretical questions

1. What does `fork` return in the parent and in the child?
2. What does `exec` change, and what does it keep?
3. Why does the parent call `wait` or `waitpid`?
4. What must the child do if `exec` fails?
5. Why is `fork` plus `exec` the usual create-and-run pattern?

#### Easy practical tasks

1. Open `man 2 fork`, `man 3 execvp`, and `man 2 waitpid`. Write one sentence for each.
2. Write the parent and child return values of `fork` as a two-row table.
3. Run `strace -e fork,clone,execve,wait4 /bin/true` or `strace -f -e execve ls`. Write one call that you see.
4. Write five sentences about `fork`, `exec`, and `wait`. Use only facts from this section.

#### Medium practical tasks

1. Write a C program: `fork`, child prints `child` and exits, parent `waitpid`s and prints `parent`. This handbook does not contain the source.
2. Change the program so that the child calls `execlp("date", "date", (char *)0)`. Confirm that the parent still waits.
3. Draw the address space before `fork`, after `fork`, and after `exec` in the child.

#### Advanced practical tasks

1. Trace your program with `strace -f -e fork,clone,execve,wait4,exit_group`. Write the child PID story.
2. Read `man 2 clone`. Write how `fork` relates to `clone` on Linux. Do not write a new launcher from that page yet.

---

## Zombies and orphans

A zombie is a process that has exited. The kernel still keeps a small PCB so that the parent can read the exit status. `ps` shows state `Z`. A zombie does not run. A zombie holds little memory. A large number of zombies still wastes PID slots.

The parent must call `wait` or `waitpid` to reap the zombie. If the parent ignores `SIGCHLD` with `SA_NOCLDWAIT` or sets a handler that reaps, the kernel can collect children without a zombie list. The default is: you wait, or you get zombies.

An orphan is a live process whose parent has exited. The kernel reparents the orphan. On Linux the new parent is usually `init` (PID 1) or a designated subreaper. `init` waits for orphans so that they do not stay zombies forever.

A zombie is dead and not yet reaped. An orphan is alive and has a new parent. Do not mix the words.

A long-lived server that `fork`s workers must wait or use `SIGCHLD`. If it never waits, the zombie table grows.

`kill` cannot "run" a zombie. The parent must wait. Killing the parent can let `init` reap the zombie.

### Questions

#### Theoretical questions

1. What is a zombie process?
2. What is an orphan process?
3. Who reaps an orphan after the original parent exits?
4. Why does a zombie still need a PID?
5. How does a server avoid a pile of zombies?

#### Easy practical tasks

1. Open `man 2 wait` and find the zombie description. Write one sentence.
2. Write a two-column table: "Zombie" and "Orphan". Add four comparison rows.
3. Run `ps -el | head` or `ps -o pid,stat,cmd`. Write what the `STAT` letter `Z` means if you see it (or write that you see none).
4. Write five sentences about zombies and orphans. Use only facts from this section.

#### Medium practical tasks

1. Write a program that `fork`s and the parent `sleep`s without `wait`. In another terminal run `ps -o pid,ppid,stat,cmd -p <child>`. Write the `STAT` letter. Then stop the parent. This handbook does not contain the source.
2. After the parent exits, run `ps` again on the child if it still sleeps. Write the new PPID.
3. Draw a timeline: child exits, zombie exists, parent waits, PCB is gone.

#### Advanced practical tasks

1. Read `man 2 prctl` for `PR_SET_CHILD_SUBREAPER`. Write when a user process, not PID 1, reaps orphans.
2. Write a one-page note: `SIGCHLD`, `waitpid(-1, ...)`, and why a signal handler must be careful with async-signal-safe calls.

---

## Process states and environment variables

A process moves through states. Common textbook names:

- running: on a CPU
- ready (runnable): wants a CPU, waits in the scheduler queue
- blocked (sleeping): waits for I/O, a child, or a lock
- stopped: received `SIGSTOP` or a tty stop
- zombie: exited, not reaped

Linux `ps` letters include `R` (running or runnable), `S` (interruptible sleep), `D` (uninterruptible sleep, often disk), `T` (stopped), `Z` (zombie). One letter is a summary. `/proc/<pid>/status` has a `State:` line.

The scheduler picks a ready process. Topic 3 covers that choice. A blocked process does not consume a time slice until it wakes.

Environment variables are a list of `NAME=value` strings. The kernel stores a copy for the process. `exec` can pass a new list (`execve`) or keep the current list (`execvp` uses `environ`). A child after `fork` inherits a copy of the parent environment.

Common variables: `PATH`, `HOME`, `USER`, `LANG`. `PATH` is a list of directories that the shell searches for a command. `PATH` is not a filesystem path of one file.

A process can change its own environment with `setenv` or `putenv`. That change does not change the parent. The shell exports variables so that children see them.

Do not put secrets in the environment on a shared host if other users can read `/proc/<pid>/environ`.

### Questions

#### Theoretical questions

1. Name five process states in this section.
2. What does Linux `ps` letter `S` mean?
3. When does a process leave the blocked state?
4. What is an environment variable?
5. Does a child share the same environment memory as the parent after `fork`? Explain.

#### Easy practical tasks

1. Run `ps -o pid,stat,cmd`. Write the `STAT` letter of your shell.
2. Run `echo $PATH` and `printenv HOME`. Write both values.
3. Open `man 7 environ`. Write one sentence about inheritance.
4. Write five sentences about process states and the environment. Use only facts from this section.

#### Medium practical tasks

1. Run `sleep 30 &` then `ps -o pid,stat,cmd -p <pid>`. Write why the state is `S`.
2. Run `cat /proc/self/environ | tr '\0' '\n' | head`. Write three names that you see.
3. Draw a state diagram: ready, running, blocked, stopped, zombie. Label `fork`, I/O wait, `wait`, and `SIGSTOP`.

#### Advanced practical tasks

1. Compare `R`, `S`, and `D` in `man 1 ps`. Write why `D` is hard to kill.
2. Write a C program that prints one variable with `getenv` and then `execve`s `/usr/bin/env` with a custom list. Write how the child list differs. This handbook does not contain the source.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A classmate says "a process is a program on disk." Which facts do you use to correct that sentence?
2. How do the PCB idea and `/proc/<pid>` work together?
3. Why can the same PID value name a new process after the old one is reaped?
4. How do `fork`, `exec`, and `wait` form one create-run-reap story?
5. When is a sleeping process not a zombie, and when is a dead process not an orphan?

#### Easy practical tasks

1. Write a one-page cheat sheet: program, process, thread, PID, PPID, PCB, `fork`, `exec`, `wait`, zombie, orphan, `PATH`.
2. Run `ps -o pid,ppid,stat,cmd` and `echo $$`. Mark your shell on the list.
3. Draw parent and child after `fork` only, then after the child `exec`s `date`.
4. Open `man 7 signal` and find `SIGCHLD`. Write one sentence that links it to this topic.

#### Medium practical tasks

1. Write a program that prints PID and PPID in parent and child, then `waitpid`s. Save the output. This handbook does not contain the source.
2. Use `strace -f` on `/bin/sh -c date`. Write which process calls `execve`.
3. Document a five-step checklist to debug a zombie: `ps`, PPID, whether the parent waits, `SIGCHLD`, restart.

#### Advanced practical tasks

1. Read the `fork` chapter in [OSTEP](https://ostep.org/). Write a one-page summary of copy-on-write as a preview of topic 5.
2. Compare `posix_spawn` documentation with `fork` plus `exec`. Write three reasons a project picks `posix_spawn`.
