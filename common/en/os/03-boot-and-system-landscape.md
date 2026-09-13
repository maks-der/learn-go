# 3. Boot and System Landscape

## Description

This topic follows the machine from firmware to a running Linux system. You learn BIOS and UEFI, the GRUB bootloader, kernel initialization, and `init` or `systemd`. You also learn system calls as the border between user space and the kernel, and you learn `/proc` and `/sys`.

Complete this topic after the hardware review. Complete this topic before you study processes in detail. You already know the stage names from topic 1. This topic adds the Linux landscape.

Use one term for each concept. Firmware is the software that the vendor stores on the motherboard. The bootloader loads the kernel. `init` is the first user process. A system call is the request interface into the kernel. `/proc` and `/sys` are virtual filesystems that export kernel state. Do not mix firmware with the bootloader. Do not mix `/proc` with a real disk filesystem.

---

## BIOS and UEFI

Firmware runs when the machine powers on. Two firmware interfaces are common on PCs.

BIOS is the older interface. BIOS uses 16-bit start code and a firmware call table. BIOS reads a boot sector from a disk that uses a Master Boot Record (MBR). The MBR is a small region at the start of the disk. BIOS has size limits and a limited device model.

UEFI is the current interface on most PCs. UEFI uses a richer firmware environment. UEFI reads a boot application from an EFI System Partition (ESP). The ESP is a FAT filesystem. Disks often use a GUID Partition Table (GPT) with UEFI. UEFI can check signatures when Secure Boot is on. Topic 17 covers Secure Boot at a high level.

Firmware jobs:

1. Initialize memory controllers and basic devices.
2. Run self-tests.
3. Find a boot policy: which disk or which UEFI application to start.
4. Load that bootloader or boot application into RAM.
5. Jump to it.

Firmware is not the operating system. Firmware does not schedule your processes. After the kernel runs, Linux can expose firmware tables through `/sys/firmware`. ACPI tables describe power and devices. The kernel parses those tables.

On some machines the firmware is neither PC BIOS nor UEFI. Embedded boards use other ROM monitors. The idea stays the same: a small vendor program loads the next stage.

### Questions

#### Theoretical questions

1. What is firmware?
2. How does BIOS find a bootloader on disk?
3. How does UEFI find a boot application?
4. What is the EFI System Partition?
5. Why is firmware not the operating system?

#### Easy practical tasks

1. Write a two-column table: "BIOS" and "UEFI". Add four comparison rows.
2. On Linux, test whether `/sys/firmware/efi` exists. Write "UEFI boot" or "not UEFI" and the test that you used.
3. Run `ls /sys/firmware` if the directory exists. Write three names that you see.
4. Open `man 8 efibootmgr` if the tool exists. Write one sentence about what the tool manages.

#### Medium practical tasks

1. Run `lsblk -o NAME,PTTYPE,FSTYPE,MOUNTPOINT`. Write whether you see `gpt` or `dos` and whether you see a `vfat` ESP.
2. Draw the UEFI path: firmware, ESP, bootloader file, kernel. Draw the BIOS path: firmware, MBR, bootloader, kernel.
3. Read a distribution page on Secure Boot. Write what the firmware checks before it starts GRUB or shim.

#### Advanced practical tasks

1. Use `efibootmgr -v` when you have rights and UEFI. Write the boot order and the path of one boot entry. Do not change the order.
2. Compare MBR and GPT in a one-page report: size limits, partition count, and which firmware interface usually pairs with each.

---

## Bootloader (GRUB)

The bootloader is the program that loads the kernel. On many Linux distributions the bootloader is GNU GRUB.

GRUB can show a menu. Each menu entry names a kernel image, optional modules, an initial ramdisk (`initramfs` or `initrd`), and a kernel command line. The command line sets the root device, console, and other early options.

Typical GRUB jobs:

1. Read configuration from the boot filesystem.
2. Let the user pick an entry or a rescue entry.
3. Load `vmlinuz` (the kernel image) into memory.
4. Load `initramfs` into memory when the entry requires it.
5. Jump to the kernel entry with the command line.

`initramfs` is a small filesystem in RAM. The kernel uses it to find the real root filesystem. Drivers and `cryptsetup` can live in `initramfs` so that the kernel can open an encrypted disk or a complex logical volume.

GRUB configuration on Debian-style systems lives under `/boot/grub` and `/etc/default/grub`. You regenerate `grub.cfg` with `update-grub` or `grub-mkconfig`. A syntax error can make the machine unbootable. Use a VM for experiments.

Other bootloaders exist: systemd-boot, syslinux, U-Boot on many boards. This path uses GRUB as the example. The job is the same: load a kernel and pass parameters.

### Questions

#### Theoretical questions

1. What is the job of a bootloader?
2. What files does a typical GRUB entry load?
3. What is `initramfs`?
4. What is the kernel command line?
5. Why can a bad GRUB configuration prevent boot?

#### Easy practical tasks

1. Run `ls /boot` on a Linux VM. Write the names that look like a kernel and an initramfs.
2. Run `cat /proc/cmdline`. Write the options that you recognize (`root=`, `ro`, `quiet`).
3. Open `man 7 bootparam`. Write two parameters and their meaning.
4. Write five sentences that describe GRUB. Use only facts from this section.

#### Medium practical tasks

1. Read `/etc/default/grub` if the file exists (use `sudo` only if needed). Write the purpose of `GRUB_CMDLINE_LINUX_DEFAULT` without changing the file.
2. Draw GRUB as a box between firmware and the kernel. Label menu, `vmlinuz`, `initramfs`, and command line.
3. Find `grub.cfg` under `/boot/grub` or `/boot/grub2`. Count the `menuentry` lines. Do not edit the file.

#### Advanced practical tasks

1. In a disposable VM, add `systemd.unit=multi-user.target` or `single` to a GRUB command line for one boot. Document the menu keys that you used and the result. Restore the default.
2. Compare GRUB and systemd-boot in documentation. Write a table: configuration format, UEFI support, and typical distribution.

---

## Kernel init

Kernel initialization starts when the bootloader jumps to the kernel. The kernel is not yet ready to run your shell.

Early kernel init does architecture setup: CPU mode, early console, and a map of physical memory. The kernel then starts the memory allocator, the scheduler structures, and the interrupt system. It probes buses and loads drivers that are built into the kernel.

If `initramfs` is present, the kernel unpacks it and runs a small user-space program from that ramdisk. That program loads extra modules and mounts the real root filesystem. Control then moves to the real root.

The kernel mounts the root filesystem. It sets up `init` as process ID 1. If the kernel cannot find root or cannot start `init`, it panics. A kernel panic is a stop. The machine does not continue to a login.

You can see some init results after boot:

- `dmesg` prints the kernel log.
- `/proc/version` names the kernel.
- `/proc/cmdline` shows the command line that the bootloader passed.

Kernel init is not `systemd` init. Kernel init is the kernel bringing itself up. User-space init starts after the kernel is ready.

### Questions

#### Theoretical questions

1. What is kernel initialization?
2. What is a kernel panic in the boot path?
3. Why does the kernel need `initramfs` on some systems?
4. When does the kernel start process ID 1?
5. How does `dmesg` help you study kernel init?

#### Easy practical tasks

1. Run `uname -r` and `cat /proc/version`. Write the kernel release.
2. Run `dmesg | head` or `journalctl -k -b --no-pager | head` if `dmesg` needs rights. Write two early messages.
3. Open `man 1 dmesg`. Write how to print human-readable timestamps.
4. Write a numbered list of kernel-init steps from this section (five steps).

#### Medium practical tasks

1. Search the kernel log for `Command line`, `memory`, and `root`. Write one line for each match.
2. Draw kernel init as a flowchart: early setup, drivers, optional initramfs, mount root, start `init`, or panic.
3. Compare `/proc/cmdline` with one `dmesg` line that repeats the command line. Confirm that they match.

#### Advanced practical tasks

1. Read the Linux documentation on `initramfs`. Write who runs inside the ramdisk and how `switch_root` or `pivot_root` leads to the real root.
2. Boot a VM with a wrong `root=` parameter on purpose. Record the panic or the initramfs shell. Restore a working entry.

---

## `init` and systemd

The first user process is `init`. Its process ID is 1. If `init` exits, the kernel panics on Linux. `init` must keep running.

Traditional Unix used a program named `init` that read `/etc/inittab` and started a few runlevels. Many Linux systems now use `systemd` as PID 1. `systemd` is an init system and a service manager.

`systemd` jobs at a high level:

1. Mount API filesystems and extra filesystems from unit files.
2. Start services in parallel when their dependencies allow it.
3. Supervise services and restart them when the unit asks for that.
4. Start login services (`getty`, a display manager).
5. Collect service logs through the journal.

A unit is a configuration object. A `.service` unit describes one service. A `.target` unit groups units. `multi-user.target` is a text login system. `graphical.target` adds a graphical login.

You do not need systemd internals for this path. You need the idea: the kernel starts one user process, and that process starts the rest of user space.

Other init systems exist: OpenRC, s6, BusyBox `init`. Embedded images often use a small `init`. The PID 1 contract stays the same.

### Questions

#### Theoretical questions

1. Why is process ID 1 special?
2. What is `systemd` at a high level?
3. What is a systemd unit?
4. What is the difference between kernel init and user-space `init`?
5. What happens on Linux if PID 1 exits?

#### Easy practical tasks

1. Run `ps -p 1 -o pid,ppid,cmd`. Write the command of PID 1.
2. Run `ls /usr/lib/systemd/system | head` if systemd is present. Write five unit file names.
3. Run `systemctl is-system-running` if the command exists. Write the state.
4. Open `man 1 systemd`. Write one sentence that describes PID 1.

#### Medium practical tasks

1. Run `systemctl list-units --type=service --state=running | head`. Write three running services and a guess for each job.
2. Draw a tree: kernel, `systemd`, three services, a login shell. Use `pstree -p 1 | head` to check the idea.
3. Read `man 5 systemd.service`. Write the meaning of `ExecStart` and `WantedBy`.

#### Advanced practical tasks

1. Create a user systemd service that runs `sleep 3600` (user session, not a system-wide change if you can avoid `sudo`). Start it, show it, stop it. Document the commands.
2. Compare `systemd` and a classic `/etc/inittab` `init` in a one-page table: parallelism, dependencies, and logging.

---

## Syscalls as the user–kernel border

A system call is the official border between user space and the kernel. User code cannot jump into arbitrary kernel functions. User code can only enter at the system-call gate, plus traps and exceptions that the CPU already defines.

On Linux, each system call has a number and a name. The CPU puts the number in a register and executes `syscall`. The kernel table jumps to the implementation. Arguments go in registers. A result returns in a register. A negative result in the kernel ABI often becomes `-1` and `errno` in the C library.

Examples:

- `openat` opens a file and returns a file descriptor.
- `read` and `write` move bytes.
- `clone` and `execve` create and replace processes.
- `mmap` maps memory.
- `exit_group` ends the process.

The C library wraps many calls. `fopen` is not a system call. `fopen` calls `open` internally. `malloc` is not a system call on every invocation. `malloc` may call `brk` or `mmap` when it needs more memory.

`strace` prints the system calls of one process. That tool is the practical way to see the border.

A system call is not a function call in the same process only. It is a controlled mode switch. It can block. It can fail. The kernel checks permission on each call.

### Questions

#### Theoretical questions

1. Why is the system-call gate the user–kernel border?
2. How does a Linux user program identify which system call to run?
3. Why is `fopen` not a system call?
4. What does `strace` show?
5. What can the kernel do with the arguments of a system call before it performs the work?

#### Easy practical tasks

1. Open `man 2 syscalls`. Write five system-call names that you recognize.
2. Run `strace -e trace=%file ls /tmp` (or `strace ls /tmp`). Write three file-related calls.
3. Run `man 2 open` and `man 3 fopen`. Write which page is a system call and which page is a library function.
4. Write four sentences that describe a system call. Use only facts from this section.

#### Medium practical tasks

1. Trace `cat` on a small file with `strace -o /tmp/cat.trace cat /etc/hostname`. Count `openat`, `read`, `write`, and `close`.
2. Draw the register path of a system call: user registers, `syscall` instruction, kernel table, return, C library `errno`.
3. Compare `man 2 read` return values with what `strace` prints for `read`. Write how a failure appears.

#### Advanced practical tasks

1. Read `man 2 syscall` and a Linux syscall table for your architecture (x86-64 or arm64). Write the numbers of `read`, `write`, and `exit`.
2. Write a C program that calls `write(1, "ok\n", 3)` with no stdio. Trace it. Confirm that you see `write` and not `fputs`.

---

## `/proc` and `/sys` on Linux

Linux exports kernel state through virtual filesystems. The files are not on a disk. A read calls into the kernel. A write, when allowed, changes kernel settings.

`/proc` is the proc filesystem. It shows processes and some system information.

- `/proc/self` is the process that opens the path.
- `/proc/<pid>/status` and `/proc/<pid>/maps` describe one process.
- `/proc/cpuinfo` and `/proc/meminfo` describe hardware as the kernel sees it.
- `/proc/sys` holds tunables (the `sysctl` tree).

`/sys` is sysfs. It shows devices, buses, drivers, and firmware objects as a tree. `/sys` is the modern place for device attributes. Examples: `/sys/class/net`, `/sys/block`, `/sys/firmware`.

Both filesystems appear in `mount` as `proc` and `sysfs`. You can `cat` many files. Some files are binary. Some files need privilege. Do not write to `/proc/sys` or `/sys` on a shared machine unless you know the effect.

`/proc` is also how many tools work. `ps` reads `/proc`. `free` reads `/proc/meminfo`. Learn the files, then the tools make sense.

Topic 20 returns to tracing. This section only teaches the landscape: the kernel has a directory tree that looks like files.

### Questions

#### Theoretical questions

1. Why is `/proc` called a virtual filesystem?
2. What is `/proc/self`?
3. What kind of objects does `/sys` show?
4. How does `ps` relate to `/proc`?
5. Why is a write to `/proc/sys` dangerous on a shared machine?

#### Easy practical tasks

1. Run `ls /proc | head`. Write which names are process IDs and which names are system files.
2. Run `cat /proc/self/comm` and `readlink /proc/self`. Write what you see.
3. Run `ls /sys/class`. Write five class names.
4. Open `man 5 proc`. Write the purpose of `/proc/<pid>/maps` in one sentence.

#### Medium practical tasks

1. Compare `/proc/self/status` with `ps -p $$ -o pid,rss,stat`. Write two fields that appear in both views.
2. Run `mount | grep -E 'proc|sysfs'`. Write the mount points and the filesystem types.
3. Draw two trees: five useful `/proc` paths and five useful `/sys` paths. Label each path with one use.

#### Advanced practical tasks

1. Read `man 5 sysfs` if present, or kernel sysfs documentation. Write how a device directory in `/sys` relates to a driver.
2. Use `sysctl -a | head` (may need rights). Map one `sysctl` name to a file under `/proc/sys`. Write the mapping rule (dots to slashes).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from UEFI firmware to a `systemd` service start. Name each stage.
2. Where can boot fail, and what do you see at each failure (firmware, GRUB, kernel panic, PID 1)?
3. How do the system-call border and `/proc` differ as ways to talk to the kernel?
4. Why does the kernel need both `initramfs` and a real root filesystem on many desktops?
5. A program uses only libc and never opens `/proc`. Does it still use system calls? Explain.

#### Easy practical tasks

1. Write a one-page cheat sheet: BIOS, UEFI, ESP, GRUB, `vmlinuz`, `initramfs`, `dmesg`, PID 1, systemd unit, system call, `/proc`, `/sys`.
2. Run `cat /proc/cmdline`, `ps -p 1 -o comm=`, and `ls /sys/firmware`. Write one sentence per command.
3. Draw the user–kernel border. Put `strace` on the user side and `dmesg` on the kernel-log side.
4. Bookmark kernel documentation and `man 7 bootparam`. Write when you open each.

#### Medium practical tasks

1. Write a script that records `/proc/version`, `/proc/cmdline`, PID 1 comm, and whether `/sys/firmware/efi` exists. Save a boot report.
2. Trace `true` with `strace -c`. Write how a process that "does nothing" still crosses the system-call border.
3. Walk `/proc/1` with `ls` (some files need rights). Write which files you can read as a normal user and what they tell you about `init`.

#### Advanced practical tasks

1. In a VM, boot once with `init=/bin/sh` on the kernel command line if the image allows it. Write what you get and why systemd did not start. Restore a normal boot.
2. Read about the Linux boot process in the kernel documentation (`start_kernel` overview). Write a one-page map from `start_kernel` to the first user process.
