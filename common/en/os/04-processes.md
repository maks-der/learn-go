# 4. Processes

## Description

A process is a running instance of a program. This topic explains the difference between a program, a process, and a thread. You learn the process control block, process IDs, and the Unix calls `fork`, `exec`, and `wait`. You also learn zombies, orphans, environment variables, and process states.

Complete this topic after the boot landscape. Complete this topic before you study threads and scheduling. Those topics assume a process.

Use one term for each concept. A program is a file of instructions and data. A process is the OS object that runs a program. A thread is an execution context inside a process. A zombie is a process that has exited and that the parent has not yet waited for. Do not mix `fork` with `exec`. Do not mix a zombie with an orphan.

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

A thread is the unit that the CPU runs. A thread has a stack and a set of registers. A process starts with one thread. The process can create more threads. All threads in one process share the same address space and the same file descriptor table. Topic 5 covers threads.

The program is the stored instructions. The process is the running instance. Each thread is one execution of those instructions inside the same process. The threads share the same address space.

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

1. Run the same binary twice: `/bin/true ; /bin/true` and then `pidof bash` or `pgrep -a bash`. Write how you distinguish the processes.
2. Draw a diagram: one ELF file on disk, two processes, one of those processes with two threads.
3. Use `readlink /proc/$$/exe` and `cat /proc/$$/cmdline`. Write the difference between the program file and the command line.

#### Advanced practical tasks

1. Compile a tiny C program. Run it three times at once (`./a.out &` three times). Show three PIDs and one path from `/proc/<pid>/exe`.
2. Read `man 5 elf` or an ELF overview. Write which parts of the file become the process text and data after `exec`.

---

## Process control block (PCB) idea

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

The PCB is the reason `ps` can print a list. The kernel walks its table of process structures and copies safe fields to user space.

### Questions

#### Theoretical questions

1. What is a process control block?
2. Why does the kernel store registers in the PCB?
3. Name five classes of data that a PCB holds.
4. What is the Linux kernel name that is closest to a per-thread PCB?
5. How does `/proc/<pid>` relate to the PCB idea?

#### Easy practical tasks

1. Make a table: "PCB field" and "Where you see it on Linux". Add five rows (`Pid`, `PPid`, `State`, `Uid`, `VmSize` from `/proc/<pid>/status`).
2. Run `cat /proc/self/status | head -20`. Write five fields that match the PCB list in this section.
3. Open `man 5 proc` for `/proc/<pid>/status`. Write the meaning of `State` and `PPid`.
4. Write four sentences that describe the PCB. Use only facts from this section.

#### Medium practical tasks

1. Compare `/proc/self/status` and `/proc/self/stat`. Write two fields that appear in both forms.
2. Draw a PCB as a box with eight labeled fields. Draw an arrow to the CPU labeled "context load".
3. Read a short description of `task_struct` in kernel documentation or a textbook. Write three fields that the book or page names.

#### Advanced practical tasks

1. Write a program that reads `/proc/self/status` and prints `Pid`, `PPid`, and `State`. Explain which PCB ideas you displayed.
2. Compare the PCB idea with the inode idea (preview of topic 11): both are kernel objects that user space sees only through interfaces. Write one page.

---

## PID, parent, child

The kernel assigns a process ID (PID) to each process. The PID is an integer. On Linux the PID is unique for the life of the process in that PID namespace. After the process exits and is reaped, the kernel can reuse the PID.

The process that creates a new process is the parent. The new process is the child. The child has a parent PID (PPID) that equals the parent's PID at creation. You can see the tree with `ps -ejH` or `pstree`.

Process ID 1 is `init` (often `systemd`). `init` is the ancestor of user processes after boot. If a parent exits before the child, the child becomes an orphan and the kernel reparents the child. On Linux the new parent is often `init` or a subreaper. The orphan section below covers that case.

A parent can wait for a child. The parent then receives the child's exit status. That wait also removes the zombie. The `wait` section covers the calls.

Do not assume that a smaller PID is an older process after PID reuse. Do not assume that two processes with similar command lines share a parent.

Linux also has thread IDs (TIDs). For a single-threaded process, TID and PID match. Topic 5 covers the difference.

### Questions

#### Theoretical questions

1. What is a PID?
2. What is a PPID?
3. What is the parent–child relation?
4. Why can the kernel reuse a PID after a process is gone?
5. What is special about PID 1 in the process tree?

#### Easy practical tasks

1. Run `echo $$` in the shell. Write your shell PID.
2. Run `ps -o pid,ppid,cmd`. Find your shell and its parent. Write both IDs.
3. Run `pstree -p | head`. Write the path from PID 1 to your shell if you can see it.
4. Open `man 2 getpid` and `man 2 getppid`. Write what each call returns.

#### Medium practical tasks

1. Start `sleep 300 &`. Run `ps -o pid,ppid,cmd -p $!`. Write the child PID and the PPID. Kill the sleep with `kill $!`.
2. Draw a process tree with PID 1, a login service, a shell, and two children of the shell.
3. Compare `$$` and `$PPID` in bash. Write a one-line C program that prints `getpid()` and `getppid()` and run it from that shell.

#### Advanced practical tasks

1. Read `man 7 pid_namespaces`. Write how a process can have PID 1 inside a namespace and a different PID on the host.
2. Write a C program that forks once and prints PID and PPID in parent and child. Show that the child's PPID equals the parent's PID.

---

## `fork`, `exec`, `wait` (Unix)

Unix creates processes with a pair of ideas: copy, then replace.

`fork` creates a new process. The child is a copy of the parent. The child gets a new PID. The child inherits file descriptors, signal dispositions (with rules), and the address space as a copy (often copy-on-write; topic 9). `fork` returns 0 in the child. `fork` returns the child PID in the parent. `fork` returns `-1` on failure.

`exec` replaces the current process image with a new program. The family includes `execve`, `execlp`, and others. After a successful `exec`, the same PID runs different code. Open file descriptors stay open unless the program set close-on-exec. `exec` does not return on success.

`wait` and `waitpid` let a parent block until a child changes state. The usual use is to wait until the child exits. The parent then receives the exit status and the kernel frees the zombie.

A typical shell path:

1. The shell reads a command.
2. The shell calls `fork`.
3. The child calls `exec` to run the command.
4. The parent calls `wait` (unless the command is in the background).

Linux also has `clone`. `fork` is a special case of `clone`. `posix_spawn` can combine the steps. Learn `fork`, `exec`, and `wait` first.

```c
#include <unistd.h>
#include <sys/wait.h>
#include <stdio.h>

int main(void)
{
	pid_t pid = fork();
	if (pid == 0) {
		execlp("echo", "echo", "from-child", (char *)0);
		return 127;
	}
	if (pid > 0)
		wait(NULL);
	return 0;
}
```

### Questions

#### Theoretical questions

1. What does `fork` return in the parent and in the child?
2. What does a successful `exec` do to the process image?
3. Why does a parent call `wait`?
4. Why does a shell use `fork` and then `exec` instead of only `exec`?
5. What is `clone` in relation to `fork` on Linux?

#### Easy practical tasks

1. Open `man 2 fork`, `man 3 exec`, and `man 2 wait`. Write one sentence for each page.
2. Type the C program from this section. Compile it. Run it. Write the output.
3. Run `strace -e trace=process echo hi` or `strace -f -e clone,execve,wait4 echo hi`. Write which calls you see.
4. Write the four-step shell path from this section as a numbered list with one extra note about background jobs.

#### Medium practical tasks

1. Change the sample so that the parent prints `wait done` after `wait`. Confirm the order of output.
2. Make the child `exec` a missing binary. Print an error in the child when `exec` fails. Show the parent still waits.
3. Draw a timeline: parent `fork`, child `exec`, child exit, parent `wait`.

#### Advanced practical tasks

1. Write a tiny shell that reads one line, `fork`, `execvp`, and `waitpid`. Handle `exec` failure. Do not implement pipes yet.
2. Compare `vfork` and `posix_spawn` in the manual. Write when a program might avoid a full `fork` of a large process.

---

## Zombies and orphans

A zombie is a process that has exited but that still has an entry in the kernel table. The entry holds the exit status until the parent waits. The zombie does not run. It does not use CPU. It uses a PID slot and a small kernel object.

You see zombies in `ps` with state `Z` and often the mark `<defunct>`. A few short-lived zombies are normal. A large number of zombies means the parent does not call `wait`. The fix is in the parent: wait for children, or set `SIGCHLD` handling so that the parent reaps.

An orphan is a process whose parent has exited. The child still runs. The kernel sets a new parent. On Linux the new parent is PID 1 or a process that marked itself as a child subreaper. The new parent will wait for the orphan when the orphan later exits. That rule prevents a permanent zombie with no waiter.

Do not mix the two words.

- Zombie: child is dead, parent did not wait yet.
- Orphan: parent is dead, child still lives (or the child already died and the new parent must reap it).

A double-fork is a pattern that makes a daemon. The child exits, the grandchild is adopted, and the original parent does not need to wait for the grandchild. You will see that pattern in service programs.

### Questions

#### Theoretical questions

1. What is a zombie process?
2. What resource does a zombie still hold?
3. What is an orphan process?
4. Who waits for an orphan after reparenting?
5. Why does a parent that never calls `wait` create many zombies?

#### Easy practical tasks

1. Write a two-column table: "Zombie" and "Orphan". Add four comparison rows.
2. Open `man 2 wait`. Write what happens if the parent never waits.
3. Run `ps -e -o pid,ppid,stat,cmd | head`. Find the `STAT` column. Write what `Z` would mean.
4. Write five sentences that contrast zombies and orphans. Use only facts from this section.

#### Medium practical tasks

1. Write a C program that forks, the child exits immediately, and the parent sleeps 20 seconds before `wait`. In another terminal, find the zombie with `ps`. Then see it disappear after wait.
2. Write a C program that forks, the parent exits immediately, and the child sleeps 20 seconds. Show that the child's PPID becomes 1 or a subreaper.
3. Draw two timelines: one for a zombie, one for an orphan.

#### Advanced practical tasks

1. Read `man 2 prctl` for `PR_SET_CHILD_SUBREAPER`. Write how a user service manager can reap orphans instead of PID 1.
2. Create many zombies on purpose in a VM (a loop of fork and no wait, with a limit). Show `ps` and then fix the parent with `waitpid(-1, ...)`. Document the danger of an unbounded loop.

---

## Environment variables

An environment variable is a string of the form `NAME=value`. The process holds a list of these strings. Child processes inherit a copy of the list at `fork`. `exec` can keep the same environment or pass a new list (`execve`).

Common variables:

- `PATH`: directories that the shell searches for commands
- `HOME`: the user home directory
- `USER` and `LOGNAME`: the login name
- `LANG` or `LC_*`: locale
- `PWD`: the current directory (the shell often updates it)

The C API uses `getenv`, `setenv`, and the `environ` pointer. The shell uses `export`. A variable that is not exported does not go to child processes.

Environment strings are not secret storage. Any code in the process can read them. Other users can sometimes read them through `/proc/<pid>/environ` when permissions allow it. Do not put passwords in the environment on a shared machine.

The environment is process state. It is not a second filesystem. A change in a child does not change the parent.

### Questions

#### Theoretical questions

1. What is an environment variable?
2. When does a child receive a copy of the parent environment?
3. What does `export` do in the shell?
4. Why is the environment a bad place for a password?
5. Does `setenv` in a child change the parent environment?

#### Easy practical tasks

1. Run `printenv | head`. Write five names that you see.
2. Run `echo $PATH`. Write how many directories the value contains (split on `:`).
3. Open `man 7 environ` and `man 3 getenv`. Write one sentence for each.
4. Run `env -i PATH=/bin /bin/ls /`. Write why `ls` still runs.

#### Medium practical tasks

1. In a shell, set `FOO=1` without `export`, then run `python3 -c 'import os; print(os.environ.get("FOO"))'` or a C `getenv`. Write the result. Repeat with `export FOO=1`.
2. Write a C program that prints `getenv("HOME")` and one argument from `main`. Run it with `HOME=/tmp ./a.out`.
3. Compare `printenv` and `cat /proc/self/environ | tr '\0' '\n' | head`. Write the format difference.

#### Advanced practical tasks

1. Use `execve` to run `/usr/bin/env` with a custom `envp` array that contains only `HELLO=world`. Show the output.
2. Read `man 1 env` and `man 1 printenv`. Write a one-page guide: inherit, clear (`env -i`), and pass one extra variable.

---

## Process states: running, ready, blocked, zombie

A process (or a thread) is not always on the CPU. Textbooks use a small set of states.

- Running: the process executes on a CPU core now.
- Ready: the process could run, but it waits for the scheduler to give it a core.
- Blocked (sleeping): the process waits for an event. Examples: a `read` that has no data yet, a lock, a `wait` for a child, a timer.
- Zombie: the process has exited and the parent has not reaped it.

Linux `ps` state codes do not use the word "ready". A runnable process shows `R`. A process in uninterruptible sleep shows `D`. Interruptible sleep shows `S`. A stopped process (job control or a debugger) shows `T`. A zombie shows `Z`.

The kernel moves the process between states:

1. A running process that uses up its time slice becomes ready.
2. A running process that waits for I/O becomes blocked.
3. A blocked process whose event arrives becomes ready.
4. A process that exits becomes a zombie until `wait`.

Topic 6 explains who picks the next ready process. This section only names the states.

A process can stay blocked for a long time. That is normal. A process that is ready but never runs is a scheduling or priority problem.

### Questions

#### Theoretical questions

1. What is the difference between running and ready?
2. What does blocked mean?
3. Which textbook state maps to `Z` in `ps`?
4. What event moves a process from blocked to ready after a `read`?
5. What does `S` mean in the Linux `ps` STAT column?

#### Easy practical tasks

1. Run `ps -o pid,stat,cmd`. Write the STAT letter of your shell.
2. Run `sleep 30 &` and `ps -o pid,stat,cmd -p $!`. Write the state of `sleep`.
3. Open `man 1 ps` and read the PROCESS STATE CODES. Write the meaning of `R`, `S`, `D`, `T`, and `Z`.
4. Make a state diagram with four textbook boxes and arrows labeled with events from this section.

#### Medium practical tasks

1. Run `cat` with no arguments so that it blocks on stdin. In another terminal, inspect its STAT. Then type a line and see it finish.
2. Reproduce a zombie as in the zombie section. Show `Z` in `ps`. Reap it and show that the PID is gone.
3. Compare `top` or `htop` state letters with `ps`. Write a four-row translation table to textbook names.

#### Advanced practical tasks

1. Read `man 7 signal` for `SIGSTOP` and `SIGCONT`. Stop a `sleep` process and continue it. Show `T` then `S`.
2. Write a one-page map from Linux `task_struct` states (`TASK_RUNNING`, `TASK_INTERRUPTIBLE`, `TASK_UNINTERRUPTIBLE`, `EXIT_ZOMBIE`) to the textbook names.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the life of a process from `fork` to a successful parent `wait`, including one I/O block and the zombie window.
2. Which PCB fields must stay different between parent and child immediately after `fork`?
3. How do environment variables, file descriptors, and the address space behave across `fork` and then `exec`?
4. Why does the OS need both a PID and a parent pointer?
5. A teammate calls a stopped job a zombie. Which facts do you use to correct that?

#### Easy practical tasks

1. Write a one-page cheat sheet: program, process, thread, PCB, PID, PPID, `fork`, `exec`, `wait`, zombie, orphan, `environ`, states `R`/`S`/`Z`.
2. Run `ps -o pid,ppid,stat,cmd --forest`. Mark your shell, one child, and PID 1 on a printout or a copy.
3. Compile and run one `fork`+`exec`+`wait` program and one `getenv` program. Save the source in `os-practice`.
4. Open `man 7 credentials`. Write how a process UID relates to the PCB credentials field.

#### Medium practical tasks

1. Write a parent that creates two children. Each child prints its PID and exits with a different code. The parent `waitpid`s twice and prints both statuses.
2. Use `strace -f` on your tiny shell or on `bash -c 'echo x'`. Label `clone`/`fork`, `execve`, and `wait`.
3. Document a process-tree experiment: start a pipeline `sleep 10 | cat`. Show PIDs, parent, and states, then let it finish.

#### Advanced practical tasks

1. Implement `wait` in a loop that creates N children (N is an argument). Reap all. Prove with `ps` that no zombie remains.
2. Read the OSTEP process chapters. Write a one-page comparison of the textbook PCB with the files you used under `/proc/<pid>`.
