# 1. Getting Started

## Description

An operating system (OS) is the software that controls the computer hardware and that gives programs a stable set of services. This topic explains what an operating system is, how the kernel differs from user space, and how a machine starts. You also learn the main illusions that the OS gives to programs, and you learn why Linux is a good first system.

Complete this topic before you study hardware details, processes, or memory. Later topics assume these terms.

Use one term for each concept. The kernel is the privileged core of the OS. User space is the unprivileged area where user programs run. A system call is the controlled request from a user program to the kernel. Do not mix the names of firmware, bootloader, kernel, and `init`.

---

## What an operating system is

An operating system is system software. It sits between the hardware and the application programs. The hardware includes the CPU, the memory, the disks, and the devices. Application programs include a shell, a browser, a compiler, and a server.

The OS has three main jobs:

1. It manages hardware. It starts devices. It shares the CPU and the memory. It moves data to and from disks and networks.
2. It gives abstractions. A program sees a process, a file, and a socket. The program does not program the disk controller by hand.
3. It isolates programs. One program must not destroy the memory of another program. One user must not read another user's private files without permission.

The OS is not one file. The kernel is the core. Utility programs, libraries, and a package manager sit around the kernel. Together they form the system that you use.

A general-purpose OS such as Linux or Windows supports many programs at the same time. A special-purpose OS can control a device with a fixed set of tasks. This path studies a general-purpose Unix-like OS. The ideas also apply to other systems.

Without an OS, a program would need to know every device, every address, and every interrupt. That work does not scale. The OS hides that work behind a small set of services.

### Questions

#### Theoretical questions

1. What is an operating system?
2. Name the three main jobs of a general-purpose OS.
3. What is the difference between the kernel and the full operating system that a user sees?
4. Why does a program use OS abstractions instead of direct hardware control?
5. What does isolation mean for two programs on one machine?

#### Easy practical tasks

1. Write five sentences that describe an operating system. Use only facts from this section.
2. Make a two-column table: "OS job" and "One example". Add three rows.
3. List three programs that you use. For each program, write one service that the OS must give (CPU time, files, network, or display).
4. Open the manual page `man 1 intro` or `man 7 intro` on a Linux system. Write one sentence from the page that describes the system.

#### Medium practical tasks

1. Compare a general-purpose OS with a special-purpose OS in six short sentences.
2. Draw a three-layer diagram: hardware, operating system, application programs. Label two examples in each layer.
3. Find the official Linux kernel site at [https://www.kernel.org/](https://www.kernel.org/). Write the current stable kernel version that the site shows.

#### Advanced practical tasks

1. Read the first chapter of [Operating Systems: Three Easy Pieces](https://ostep.org/). Write a one-page summary of virtualization, concurrency, and persistence in your own words.
2. Interview a teammate or write a short report: name five OS services that a web server uses during one HTTP request.

---

## Kernel vs user space

The CPU has at least two privilege modes. Kernel mode is the privileged mode. User mode is the unprivileged mode. The kernel runs in kernel mode. User programs run in user space in user mode.

Kernel space is the memory and the code that belong to the kernel. User space is the memory and the code that belong to user programs. A user program must not read or write kernel memory. The CPU and the kernel enforce that rule.

A user program asks the kernel for a service with a system call. Examples: open a file, create a process, send data on a socket. The CPU switches to kernel mode, the kernel does the work, and the CPU returns to user mode.

Code in user space cannot execute privileged instructions. Examples of privileged work: change page tables, mask interrupts, talk to a device register. If a user program tries that work, the CPU stops the instruction and the kernel can kill the process.

The split has a cost. Each system call has a mode switch. The split has a benefit. A crash in one user program does not have to crash the kernel. A bug in the kernel can crash the whole machine.

On Linux, user programs link with a C library such as `glibc` or `musl`. The library wraps system calls. You often call `open` in C. The library issues the real system call for you.

### Questions

#### Theoretical questions

1. What is the difference between kernel mode and user mode?
2. What is user space?
3. How does a user program request a kernel service?
4. Why must a user program not execute privileged instructions?
5. What is one cost and one benefit of the kernel and user-space split?

#### Easy practical tasks

1. Write four sentences that contrast kernel space and user space.
2. List five actions that need a system call (for example, create a file).
3. Run `uname -s -r` on Linux. Write the kernel name and the kernel release.
4. Open `man 2 syscalls`. Write the purpose of section 2 of the manual.

#### Medium practical tasks

1. Run `cat /proc/self/status` and find the lines that start with `Uid` and `VmSize`. Write what those lines tell you about the current process.
2. Use `strace -e trace=openat,read,write,close ls` on a small directory. Count the system calls that you see. Write three call names.
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

Topic 3 covers this path in more detail. This section only names the stages and the order.

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
4. Open `man 8 grub` or `man 7 boot` if the page exists. Write one fact about the bootloader or the boot path.

#### Medium practical tasks

1. On Linux, run `cat /proc/cmdline`. Write what the kernel command line contains in your own words.
2. Draw a timeline from power-on to a login prompt. Label firmware, bootloader, kernel, `init`, and the shell.
3. Find whether your system uses UEFI. On Linux, check if `/sys/firmware/efi` exists. Write the result and what it means.

#### Advanced practical tasks

1. Read a distribution guide on GRUB. Write the path of a typical kernel image and the role of `initramfs`.
2. Compare BIOS boot and UEFI boot in a one-page table: firmware interface, disk partition style, and how the bootloader is stored.

---

## Process, file, device as illusions the OS provides

A program does not see raw hardware. The OS gives illusions. An illusion here is a stable interface that hides a more complex implementation.

A process is the illusion of a private machine. The process has its own address space, its own open files, and its own CPU time. In reality, many processes share one CPU and one physical memory. The kernel switches the CPU and maps memory.

A file is the illusion of a named byte stream. You open a path, you read and write bytes, and you close the file. In reality, the bytes live in blocks on a disk, in a page cache in RAM, or in a device. The same `read` and `write` calls work for many storage types.

A device is the illusion of a file or a small set of calls. On Unix, many devices appear under `/dev`. A program can open `/dev/null` or a terminal device with the same file API. The kernel and the device driver do the real work.

These illusions are not lies that hide all truth. You can still see hardware with tools. The point is that a normal program does not need that view. The program uses the process, the file, and the device interfaces.

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

## Linux, Windows, BSD, macOS — same ideas, different APIs

Linux, Windows, BSD, and macOS are different operating systems. They share the same core ideas: kernel and user space, processes, files, virtual memory, and device drivers. They do not share one application programming interface (API).

Unix-like systems include Linux, BSD, and macOS. They follow many POSIX rules. POSIX is a family of standards for system calls and C library functions. Calls such as `fork`, `exec`, `open`, and `read` exist on those systems with small differences.

Windows uses a different API. Win32 and NT APIs use other names and other structures. Windows still has processes, threads, files, and a kernel. The ideas match. The function names do not match.

Linux is a kernel. A Linux distribution adds `init`, a C library, a package manager, and user programs. Debian, Fedora, and Ubuntu are distributions. They share the Linux kernel and they differ in packages and defaults.

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
2. Run `uname -a` on your Linux or WSL system. Write each field that you understand.
3. Open `man 7 posix` if the page exists, or search the POSIX standard index. Write one sentence about the scope of POSIX.
4. List three POSIX calls that you expect to use later (`open`, `fork`, `read`).

#### Medium practical tasks

1. Compare `man 2 open` on Linux with the Windows `CreateFile` documentation online. Write three differences in parameters or return values.
2. Name two Linux distributions. Write one difference in package manager or release model.
3. Draw a Venn diagram: shared OS ideas in the center, POSIX-only names on one side, Windows-only names on the other.

#### Advanced practical tasks

1. Read the Linux man-pages project page at [https://www.kernel.org/doc/man-pages/](https://www.kernel.org/doc/man-pages/). Write how section 2 and section 3 differ.
2. Write a one-page report: what "Unix-like" means, and why a program that uses only POSIX calls can still fail when you move it between Linux and BSD.

---

## A Linux VM or WSL is enough to start

You do not need a second physical machine. A Linux virtual machine (VM) or Windows Subsystem for Linux (WSL) is enough for this path.

A virtual machine is a guest operating system that runs on a host. VirtualBox, QEMU, and cloud VMs are common. The guest has a virtual CPU, virtual memory, and a virtual disk. The guest kernel is a real Linux kernel.

WSL runs a Linux environment on Windows. WSL 2 uses a real Linux kernel in a lightweight VM. You get a Linux shell, Linux system calls, and `/proc`. Some hardware and some desktop features differ from a full VM.

Install a current long-term distribution, for example Debian or Ubuntu. Use a normal user account. Use `sudo` only when you need a privileged command. Do not practice as `root` for daily work.

Confirm the environment:

```text
uname -s -r
whoami
pwd
ps -p 1 -o comm=
```

You need a shell, a text editor, `man`, a compiler such as `gcc` or `clang`, and basic tools (`ps`, `cat`, `ls`). Add them with the distribution package manager when they are missing.

A container can also give a Linux shell. A container shares the host kernel. A VM has its own kernel. For kernel topics such as `/proc` and system calls, WSL 2 or a VM is closer to a real machine than a minimal container on a non-Linux host.

### Questions

#### Theoretical questions

1. Why is a Linux VM enough to learn operating systems?
2. What does WSL 2 run, at a high level?
3. What is the difference between a VM and a container for kernel study?
4. Why must you use a normal user account for daily practice?
5. Which local commands confirm that you have a usable Linux environment?

#### Easy practical tasks

1. Install or open a Linux VM or WSL. Run `uname -s -r` and save the output in a text file.
2. Run `whoami` and `id`. Write your user name and user ID.
3. Run `command -v gcc cc python3 man`. Write which tools exist.
4. Create a folder `os-practice` in your home directory. Write a one-line `README` that states your distribution name.

#### Medium practical tasks

1. Install `build-essential` or the equivalent compiler package. Compile a one-line C program that prints `ok`. Run the binary.
2. Compare WSL 2 and a full VM in a table: kernel, filesystem from Windows, and networking. Use documentation or your own system.
3. Enable `man` pages if they are missing. Open `man man`. Write how to open a page in section 2.

#### Advanced practical tasks

1. Create a snapshot or a backup of your VM before you experiment. Document the restore steps in five lines.
2. Measure whether `/proc/cpuinfo` and `/proc/sys` exist and look valid. Write a short report: is this environment good enough for topics 2 to 11?

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the full path from power-on to a running shell. Name firmware, bootloader, kernel, `init`, and a user program.
2. A classmate says "the operating system is the kernel." Which facts do you use to correct that sentence?
3. How do the process illusion and the kernel/user-space split work together?
4. Why can you study the same OS ideas on Linux and still fail to compile a Windows-only program?
5. When do you choose a VM, WSL 2, or a container for this learning path?

#### Easy practical tasks

1. Write a one-page cheat sheet with these terms: kernel, user space, system call, firmware, bootloader, `init`, process, file, device, POSIX.
2. Run `uname`, `ps -p 1`, `ls /dev | wc -l`, and `echo $HOME`. Write one line that explains each command.
3. Draw the five start stages and the kernel/user-space border on one page.
4. Bookmark kernel.org, the Linux man-pages, and OSTEP. Write one sentence about when you open each site.

#### Medium practical tasks

1. Write a small script that prints `uname -a`, the name of PID 1, and whether `/sys/firmware/efi` exists. Run the script. Save the output.
2. Use `strace -c true` or `strace -c ls`. Write the three system calls with the highest counts. State that these calls cross the user-kernel border.
3. Document your Linux practice setup in ten steps so that another beginner can copy it.

#### Advanced practical tasks

1. Boot a live Linux ISO in a VM. Compare `/proc/version` and PID 1 with your daily environment. Write five differences.
2. Read `man 7 libc` or the glibc project overview. Write how a user program, the C library, and the kernel share the work of one `open` call.
