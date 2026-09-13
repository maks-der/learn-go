# 1. What an OS Does

## Description

An operating system (OS) is the software that controls the computer hardware and that gives programs a stable set of services. This topic explains the kernel and user space, the start path of a machine, and the main illusions that the OS gives to programs. You also learn the hardware that the OS manages, the system-call border, and `/proc` and `/sys` on Linux.

Complete this topic before you study processes, memory, or files. Later topics assume these terms.

Use one term for each concept. The kernel is the privileged core of the OS. User space is the unprivileged area where user programs run. A system call is the controlled request from a user program to the kernel. Firmware is the vendor software that runs first after power-on. Do not mix firmware, bootloader, kernel, and `init`. Do not mix RAM with a storage device.

This path uses Linux and other Unix-like systems. Practice on a Linux machine, a Linux virtual machine, or WSL 2.

---

## Kernel vs user space

The CPU has at least two privilege modes. Kernel mode is the privileged mode. User mode is the unprivileged mode. The kernel runs in kernel mode. User programs run in user space in user mode.

Kernel space is the memory and the code that belong to the kernel. User space is the memory and the code that belong to user programs. A user program must not read or write kernel memory. The CPU and the kernel enforce that rule.

A user program asks the kernel for a service with a system call. Examples: open a file, create a process, send data on a socket. The CPU switches to kernel mode. The kernel does the work. The CPU returns to user mode.

Code in user space cannot execute privileged instructions. Examples of privileged work: change page tables, mask interrupts, talk to a device register. If a user program tries that work, the CPU stops the instruction. The kernel can kill the process.

The split has a cost. Each system call has a mode switch. The split has a benefit. A crash in one user program does not have to crash the kernel. A bug in the kernel can crash the whole machine.

On Linux, user programs link with a C library such as `glibc` or `musl`. The library wraps system calls. You often call `open` in C. The library issues the real system call for you.

The OS is not one file. The kernel is the core. Utility programs, libraries, and a package manager sit around the kernel. Together they form the system that you use.

### Questions

#### Theoretical questions

1. What is the difference between kernel mode and user mode?
2. What is user space?
3. How does a user program request a kernel service?
4. Why must a user program not execute privileged instructions?
5. What is one cost and one benefit of the kernel and user-space split?

#### Easy practical tasks

1. Write four sentences that contrast kernel space and user space.
2. List five actions that need a system call. Example: create a file.
3. Run `uname -s -r` on Linux. Write the kernel name and the kernel release.
4. Open `man 2 syscalls`. Write the purpose of section 2 of the manual.

#### Medium practical tasks

1. Run `cat /proc/self/status` and find the lines that start with `Uid` and `VmSize`. Write what those lines tell you about the current process.
2. Use `strace -e trace=openat,read,write,close ls` on a small directory. Write three call names that you see.
3. Draw the path of one `read` from a C program: user code, C library, system call, kernel, return.

#### Advanced practical tasks

1. Read `man 2 intro`. Write a short report: what a system call is, how an error returns, and what `errno` means.
2. Compare two C libraries (`glibc` and `musl`) in public documentation. Write how each one wraps system calls. Give one reason a project picks `musl`.

---

## Firmware, bootloader, kernel, init, user programs

A machine does not start the kernel as the first instruction after power-on. The start path has fixed stages.

1. Firmware runs first. On many PCs the firmware is UEFI. Older PCs use BIOS. Firmware tests hardware. Firmware finds a boot device. Firmware loads a bootloader.
2. The bootloader runs next. On many Linux systems the bootloader is GRUB. The bootloader loads the kernel image into memory. The bootloader can also load an initial ramdisk. The bootloader passes a command line to the kernel.
3. The kernel starts. The kernel sets up memory, interrupts, and drivers. The kernel mounts a root filesystem. The kernel starts the first user process.
4. The first user process is `init`. On many Linux systems `init` is `systemd`. `init` has process ID 1. `init` starts services and login programs.
5. User programs run after that. A shell, a display manager, and your applications are user programs.

Each stage has a small job. Firmware does not start your shell. The kernel does not show a graphical login by itself. `init` does not implement file read in place of the kernel.

If one stage fails, the next stage does not run. A bad bootloader configuration can stop the kernel from loading. A kernel panic can stop `init` from starting. A failed `init` can leave the system without services.

`initramfs` is a small filesystem in RAM. The kernel uses it to find the real root filesystem. Drivers can live in `initramfs` so that the kernel can open a complex disk.

On Linux, run `ps -p 1 -o comm=` to see the name of process ID 1. Run `cat /proc/cmdline` to see the kernel command line.

### Questions

#### Theoretical questions

1. What is the order of firmware, bootloader, kernel, `init`, and user programs?
2. What is the job of firmware at power-on?
3. What is the job of the bootloader?
4. Why does the kernel start `init` as the first user process?
5. What happens when one stage in the start path fails?

#### Easy practical tasks

1. Write the five stages of start as a numbered list. Add one sentence to each stage.
2. Run `ps -p 1 -o pid,comm`. Write the name of process ID 1.
3. Run `echo $SHELL`. Write the path of your login shell. State that the shell is a user program.
4. Open `man 7 boot` if the page exists. Write one fact about the boot path.

#### Medium practical tasks

1. On Linux, run `cat /proc/cmdline`. Write what the kernel command line contains in your own words.
2. Draw a timeline from power-on to a login prompt. Label firmware, bootloader, kernel, `init`, and the shell.
3. Find whether your system uses UEFI. On Linux, check if `/sys/firmware/efi` exists. Write the result and what it means.

#### Advanced practical tasks

1. Read a distribution guide on GRUB. Write the path of a typical kernel image and the role of `initramfs`.
2. Compare BIOS boot and UEFI boot in a one-page table: firmware interface, disk partition style, and how the bootloader is stored.

---

## Process, file, and device as illusions

A program does not see raw hardware. The OS gives illusions. An illusion here is a stable interface that hides a more complex implementation.

A process is the illusion of a private machine. The process has its own address space, its own open files, and its own CPU time. In reality, many processes share one CPU and one physical memory. The kernel switches the CPU and maps memory.

A file is the illusion of a named byte stream. You open a path, you read and write bytes, and you close the file. In reality, the bytes live in blocks on a storage device, in a page cache in RAM, or in a device. The same `read` and `write` calls work for many storage types.

A device is the illusion of a file or a small set of calls. On Unix, many devices appear under `/dev`. A program can open `/dev/null` or a terminal device with the same file API. The kernel and the device driver do the real work.

These illusions are not lies that hide all truth. You can still see hardware with tools. A normal program does not need that view. The program uses the process, the file, and the device interfaces.

Later topics add more illusions: virtual memory, threads, and sockets. The same idea applies. The OS virtualizes a resource and isolates users of that resource.

### Questions

#### Theoretical questions

1. What does "illusion" mean for an OS abstraction?
2. What private things does a process appear to own?
3. Why can the same `read` call work for a disk file and for some devices?
4. What is `/dev` on a Unix system?
5. How does the process illusion relate to one CPU that many programs share?

#### Easy practical tasks

1. Run `ls /dev | head`. Write five device names and one guess for each purpose.
2. Run `echo hello > /tmp/os-illusion.txt` and `cat /tmp/os-illusion.txt`. State which illusion you used.
3. Run `ps -o pid,cmd`. Write how many processes you see and why that number is not one.
4. Write six sentences: two about processes, two about files, two about devices. Use only facts from this section.

#### Medium practical tasks

1. Compare `/dev/zero`, `/dev/null`, and `/dev/urandom` with `man` pages. Write one use for each.
2. Copy a small file with `cat` and with `dd`. Write which file abstraction both tools use.
3. Draw a diagram: one physical disk, the kernel, and three processes that each have an open file. Label the illusion and the shared hardware.

#### Advanced practical tasks

1. Read `man 4 null` and `man 4 zero`. Write how those devices implement `read` and `write`.
2. Write a one-page essay: name one illusion that is leaky (a program can still see the hardware). Give a Linux command that shows the leak.

---

## CPU, memory, disk, interrupts, and privilege modes

The OS controls hardware. You cannot understand the OS if you do not know the main parts of a computer.

The central processing unit (CPU) executes instructions. The CPU fetches an instruction from memory, decodes it, and executes it. That cycle repeats. A register is a small, fast storage cell inside the CPU. The program counter holds the address of the next instruction. The stack pointer holds the address of the current stack top. The kernel must save and restore registers when it switches processes.

A core is an independent execution unit on a chip. A modern chip has many cores. Each core has its own registers. The OS schedules threads on cores.

Random-access memory (RAM) is the main memory of the machine. RAM is volatile. The contents disappear when power is lost. A storage device holds data after power loss. A hard disk drive (HDD) has moving heads. A solid-state drive (SSD) stores data in flash. RAM is much faster than an SSD. An SSD is much faster than an HDD for random reads.

Typical latency order, from fast to slow: CPU registers, CPU caches, RAM, SSD, HDD. The OS uses RAM as a cache for storage data. That cache is the page cache on Linux.

An I/O device is hardware that is not the CPU and not the main RAM. Examples: keyboard, network card, disk controller. A bus moves data between components. Examples: PCIe, USB, NVMe.

An interrupt is a hardware signal that stops the current instruction stream so that the kernel can run a handler. A timer interrupt lets the kernel preempt a running program. A device interrupt says that data is ready or that a transfer is done. Polling is the opposite idea: the CPU asks the device again and again. Interrupts free the CPU for other work.

Privilege modes are the hardware base of kernel space and user space. User mode cannot change page tables or mask interrupts. Kernel mode can. The CPU raises the privilege level only through a controlled gate: a system call, an exception, or an interrupt.

Do not say "memory" when you mean a storage device. When you need a general word for HDD and SSD, say "storage device".

### Questions

#### Theoretical questions

1. What are the three steps of the instruction cycle?
2. Why is RAM volatile, and why does a storage device still matter?
3. What is an interrupt?
4. Why does the kernel need a timer interrupt?
5. How do privilege modes support the kernel and user-space split?

#### Easy practical tasks

1. Run `lscpu` or `cat /proc/cpuinfo`. Write the model name and the number of CPUs that the file shows.
2. Run `free -h`. Write the total RAM.
3. Run `lsblk -d -o NAME,TYPE,SIZE,ROTA,MODEL`. Write which devices look rotational (`ROTA` is 1).
4. Write five sentences: CPU, RAM, storage, interrupt, privilege mode. Use only facts from this section.

#### Medium practical tasks

1. Draw a pyramid: registers, cache, RAM, SSD, HDD. Write one speed note and one capacity note.
2. Compare `free -h` and `cat /proc/meminfo` (fields `MemTotal` and `Cached`). Write how the page cache appears in the numbers.
3. Open `man 7 time` or a short note on the timer interrupt. Write how a time slice can start from a timer.

#### Advanced practical tasks

1. Disassemble a tiny C program with `objdump -d` after `gcc -O0`. Find `mov` or `add`. Write how the CPU uses a register.
2. Write a one-page report: why a process that does not fit in RAM becomes slow. Use the terms RAM, storage device, and the later word "paging" only as a preview.

---

## Syscalls, `/proc`, and `/sys`

A system call is the official border between user space and the kernel. On Linux, section 2 of the manual documents system calls. The C library is in section 3. You write `open` in C. The library issues `openat` or another kernel call.

A system call has a name and a number. The number is the ABI. The number can differ across CPU architectures. A program must call the C library. Do not guess a raw syscall number.

`/proc` is a virtual filesystem. It is not a normal disk tree. The kernel builds the files when you read them. `/proc` exports process state and some kernel counters. Examples:

- `/proc/self` is the calling process
- `/proc/<pid>/status` holds UIDs and memory sizes
- `/proc/cpuinfo` describes the CPU
- `/proc/cmdline` is empty at the process view; `/proc/cmdline` of the kernel is `/proc/cmdline` at the root of `/proc`

`/proc/self/maps` lists virtual memory regions. Topic 5 and topic 13 use that file.

`/sys` is another virtual filesystem. It exports devices, drivers, and firmware tables. Examples: `/sys/class`, `/sys/block`, `/sys/firmware`. You can read many files as a normal user. A write to `/sys` can change device state. Do not write to `/sys` on a shared machine without a reason.

`/proc` and `/sys` look like files. They are kernel interfaces. A crash of a disk filesystem does not mean `/proc` is gone. A kernel panic does.

`strace` prints system calls of a process. Use `strace` to see the border. Topic 13 covers `strace` in more detail.

### Questions

#### Theoretical questions

1. What is a system call?
2. Why is `/proc` not a normal disk filesystem?
3. What does `/proc/self` mean?
4. What kind of data does `/sys` export?
5. Why must a program use the C library instead of a guessed syscall number?

#### Easy practical tasks

1. Open `man 5 proc` and `man 5 sysfs` if the pages exist. Write one sentence for each.
2. Run `ls /proc | head` and `ls /sys | head`. Write five names from each list.
3. Run `cat /proc/self/status | head`. Write the `Name` and `Pid` lines.
4. Write five sentences about system calls, `/proc`, and `/sys`. Use only facts from this section.

#### Medium practical tasks

1. Run `strace -c true`. Write the three system calls with the highest counts.
2. Compare `/proc/1/comm` and `ps -p 1 -o comm=`. Write whether they match.
3. Draw user program, C library, system call, kernel, `/proc` file. Label which path is a call and which path is a read.

#### Advanced practical tasks

1. Read `man 2 syscall`. Write how a raw `syscall()` differs from a libc wrapper.
2. Write a one-page note: three `/proc` files and three `/sys` files that a beginner can read safely. State what each file shows.

---

## Linux, Windows, BSD, macOS — same ideas, different APIs

Linux, Windows, BSD, and macOS are different operating systems. They share the same core ideas: kernel and user space, processes, files, virtual memory, and device drivers. They do not share one application programming interface (API).

Unix-like systems include Linux, BSD, and macOS. They follow many POSIX rules. POSIX is a family of standards for system calls and C library functions. Calls such as `fork`, `exec`, `open`, and `read` exist on those systems with small differences.

Windows uses a different API. Win32 and NT APIs use other names and other structures. Windows still has processes, threads, files, and a kernel. The ideas match. The function names do not match.

Linux is a kernel. A Linux distribution adds `init`, a C library, a package manager, and user programs. Debian, Fedora, and Ubuntu are distributions. They share the Linux kernel. They differ in packages and defaults.

This path uses Linux and POSIX names. When a later topic says `fork` or `/proc`, the text means a Unix-like system. The same idea on Windows has another name.

Do not treat "Linux" and "Unix" as the same word. Unix is a family and a history. Linux is a kernel that behaves like Unix for user programs. BSD systems are also Unix-like. macOS is a certified Unix with a Darwin kernel.

### Questions

#### Theoretical questions

1. Which ideas do Linux, Windows, BSD, and macOS share?
2. What is POSIX?
3. Why is Linux a kernel and not a full distribution?
4. How does the Windows API differ from the POSIX API at a high level?
5. Why does this learning path use Linux names?

#### Easy practical tasks

1. Make a table with four OS names and one API family for each (POSIX or Windows).
2. Run `uname -a` on your Linux system. Write each field that you understand.
3. Open `man 7 posix` if the page exists, or search the POSIX standard index. Write one sentence about the scope of POSIX.
4. List three POSIX calls that you expect to use later (`open`, `fork`, `read`).

#### Medium practical tasks

1. Compare `man 2 open` on Linux with public Windows `CreateFile` documentation. Write three differences in parameters or return values.
2. Name two Linux distributions. Write one difference in package manager or release model.
3. Draw a Venn diagram: shared OS ideas in the center, POSIX-only names on one side, Windows-only names on the other.

#### Advanced practical tasks

1. Read the Linux man-pages project page at [https://www.kernel.org/doc/man-pages/](https://www.kernel.org/doc/man-pages/). Write how section 2 and section 3 differ.
2. Write a one-page report: what "Unix-like" means, and why a program that uses only POSIX calls can still fail when you move it between Linux and BSD.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the full path from power-on to a running shell. Name firmware, bootloader, kernel, `init`, and a user program.
2. A classmate says "the operating system is the kernel." Which facts do you use to correct that sentence?
3. How do the process illusion and the kernel and user-space split work together?
4. Why can you study the same OS ideas on Linux and still fail to compile a Windows-only program?
5. How do `/proc` and a timer interrupt both help you see that many programs share one machine?

#### Easy practical tasks

1. Write a one-page cheat sheet with these terms: kernel, user space, system call, firmware, bootloader, `init`, process, file, device, interrupt, POSIX, `/proc`, `/sys`.
2. Run `uname`, `ps -p 1`, `ls /dev | wc -l`, and `echo $HOME`. Write one line that explains each command.
3. Draw the five start stages and the kernel and user-space border on one page.
4. Bookmark [kernel.org](https://www.kernel.org/), the Linux man-pages, and [OSTEP](https://ostep.org/). Write one sentence about when you open each site.

#### Medium practical tasks

1. Write a small script that prints `uname -a`, the name of PID 1, and whether `/sys/firmware/efi` exists. Run the script. Save the output.
2. Use `strace -c ls` on a small directory. Write the three system calls with the highest counts. State that these calls cross the user-kernel border.
3. Document your Linux practice setup in ten steps so that another beginner can copy it.

#### Advanced practical tasks

1. Boot a live Linux ISO in a VM. Compare `/proc/version` and PID 1 with your daily environment. Write five differences.
2. Read `man 7 libc` or the glibc project overview. Write how a user program, the C library, and the kernel share the work of one `open` call.
