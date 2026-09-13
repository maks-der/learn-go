# 2. Computer Hardware Review

## Description

The operating system controls hardware. You cannot understand the OS if you do not know the main parts of a computer. This topic reviews the CPU, memory, storage, devices, interrupts, timers, and privilege modes.

Complete this topic before you study boot, processes, and virtual memory. Later topics assume these hardware facts.

Use one term for each concept. The CPU executes instructions. A register is storage inside the CPU. RAM is volatile main memory. A disk or an SSD is persistent storage. An interrupt is a hardware signal that stops the current instruction stream so that the kernel can run a handler. Do not mix RAM with disk. Do not mix polling with interrupts.

---

## CPU, registers, ALU

The central processing unit (CPU) executes instructions. The CPU fetches an instruction from memory, decodes it, and executes it. That cycle repeats.

A register is a small, fast storage cell inside the CPU. The instruction set names the registers. Common groups:

- General-purpose registers hold values and addresses. On x86-64, examples are `rax` and `rdi`. On AArch64, examples are `x0` and `x1`.
- The program counter holds the address of the next instruction. On x86-64 this register is `rip`.
- The stack pointer holds the address of the current stack top.
- Status flags record results such as zero or carry.

The arithmetic logic unit (ALU) does integer arithmetic and logic. Examples: add, subtract, and, or, compare. The CPU can also have a floating-point unit and vector units. Those units are not the ALU.

The CPU sees memory as a sequence of addresses. It loads values from memory into registers. It stores registers back to memory. Almost all computation uses registers. Memory access is slower than register access.

A core is an independent execution unit on a chip. A modern chip has many cores. Each core has its own registers. The OS schedules threads on cores. Topic 6 covers scheduling.

The kernel must save and restore registers when it switches processes. That save and restore is part of a context switch. The registers are the live state of the running program.

### Questions

#### Theoretical questions

1. What are the three steps of the instruction cycle?
2. What is a register?
3. What work does the ALU do?
4. Why does a program move values between registers and memory?
5. Why must the kernel save registers during a context switch?

#### Easy practical tasks

1. Write five sentences that describe the CPU. Use only facts from this section.
2. Make a table: "Register role" and "What it holds". Add four rows (program counter, stack pointer, flags, general-purpose).
3. Run `lscpu` or `cat /proc/cpuinfo`. Write the vendor, the model name, and the number of cores or CPUs that the file shows.
4. Open `man 1 lscpu`. Write what `CPU(s)` and `Core(s) per socket` mean.

#### Medium practical tasks

1. Draw a diagram of one core: registers, ALU, and a line to RAM. Label fetch, decode, and execute.
2. Compare x86-64 and AArch64 register names in a short table. Use public ABI documentation or a textbook figure.
3. Run `lscpu -e` if the command exists. Write how many logical CPUs you see and whether the output shows more than one socket.

#### Advanced practical tasks

1. Read the first pages of the Intel or ARM programmer manual on registers. Write one page: name six registers and one sentence for each.
2. Disassemble a tiny C program with `objdump -d` after `gcc -O0`. Find `add` or `mov`. Write how the ALU or a move uses registers.

---

## RAM vs disk vs SSD

Random-access memory (RAM) is the main memory of the machine. The CPU and the kernel treat RAM as the place for running code, stacks, and caches. RAM is volatile. The contents disappear when power is lost.

A hard disk drive (HDD) stores data on magnetic platters. An HDD has moving heads. Access time is milliseconds. Sequential reads are faster than random reads.

A solid-state drive (SSD) stores data in flash memory. An SSD has no moving head. Access time is much smaller than an HDD. An SSD still has a much larger latency than RAM. Flash has a limited number of erase cycles. The SSD controller spreads writes.

Typical latency order, from fast to slow:

1. CPU registers
2. CPU caches
3. RAM
4. SSD
5. HDD
6. Network storage (often slower and more variable)

Capacity order is often the reverse. A disk or an SSD holds more bytes than RAM. RAM holds more bytes than the CPU cache.

The OS uses RAM as a cache for disk and SSD data. That cache is the page cache on Linux. A `write` can hit RAM first. The kernel writes to the device later. Topic 13 covers durability. This section only states the speed and the volatility difference.

Do not say "memory" when you mean "disk". Do not say "disk" when you mean an SSD. When you need a general word, say "storage device".

### Questions

#### Theoretical questions

1. Why is RAM volatile?
2. What is one mechanical difference between an HDD and an SSD?
3. Why is RAM still much faster than an SSD?
4. Why does the OS keep file data in RAM after a read?
5. What is a safe general term for HDD and SSD together?

#### Easy practical tasks

1. Run `free -h`. Write the total RAM and the used RAM.
2. Run `lsblk -d -o NAME,TYPE,SIZE,ROTA,MODEL`. Write which devices look like rotational disks (`ROTA` is 1) and which do not.
3. Make a table: "Technology", "Volatile?", "Typical latency class". Add rows for RAM, SSD, and HDD.
4. Run `df -h`. Write the size of the root filesystem and the mount path.

#### Medium practical tasks

1. Compare `free -h` and `cat /proc/meminfo` (fields `MemTotal` and `Cached`). Write how the page cache appears in the numbers.
2. Time `dd if=/dev/zero of=/tmp/os-hw.bin bs=1M count=256 conv=fdatasync` and then `dd if=/tmp/os-hw.bin of=/dev/null bs=1M`. Write the two rates. Remove the file.
3. Draw the storage hierarchy as a pyramid: registers, cache, RAM, SSD, HDD. Write one capacity note and one speed note.

#### Advanced practical tasks

1. Read `man 1 free` and the `MemAvailable` description in `man 5 proc`. Write how `available` differs from `free`.
2. Write a one-page report: why a process that does not fit in RAM becomes slow. Use the terms RAM, storage device, and later-topic word "paging" only as a preview.

---

## I/O devices and buses

An I/O device is hardware that is not the CPU and not the main RAM. Examples: keyboard, network card, disk controller, graphics adapter, USB controller.

A bus is a communication path that moves data and control signals between components. Modern machines use more than one bus. Examples: PCIe for fast devices, USB for many external devices, SATA or NVMe paths for storage.

The CPU does not control every device bit by bit in a portable OS. A device has registers or a memory-mapped region. The kernel driver reads and writes those locations. Some devices use port I/O on x86. Many devices use memory-mapped I/O (MMIO). MMIO means device registers appear as physical addresses.

A device driver is kernel code that knows one device class. The driver implements a standard kernel interface. The rest of the kernel and the user programs do not need the vendor details.

Data can move with programmed I/O, where the CPU copies each byte or word. Data can also move with direct memory access (DMA). DMA lets the device read or write RAM without a CPU copy for each word. Topic 14 covers DMA in more detail.

User programs rarely talk to the bus. They open a file, a socket, or a device node. The kernel and the driver use the bus.

### Questions

#### Theoretical questions

1. What is an I/O device in this topic?
2. What is a bus?
3. What is memory-mapped I/O?
4. What is the job of a device driver?
5. Why does a user program not program PCIe by hand?

#### Easy practical tasks

1. Run `ls /sys/bus` on Linux. Write five bus names that you see.
2. Run `lspci` if the tool exists. Write three devices from the list.
3. Run `lsusb` if the tool exists. Write two USB devices or write that none are present.
4. Write four sentences: device, bus, driver, user program. Use only facts from this section.

#### Medium practical tasks

1. Pick one line from `lspci -v` or `lspci -nn`. Write the device class and why the kernel needs a driver for it.
2. Draw a diagram: CPU, RAM, PCIe bus, network card. Show a user process that uses a socket, not the bus.
3. Read `man 1 lspci` and `man 8 lsblk`. Write how you map a storage device name to a bus or a controller.

#### Advanced practical tasks

1. Read a kernel documentation page on device drivers at [https://www.kernel.org/doc/html/latest/](https://www.kernel.org/doc/html/latest/). Write the difference between a character device and a block device in two short paragraphs.
2. Trace one USB or network device from `lsusb`/`lspci` to a node under `/sys`. Write the path and two attribute file names.

---

## Interrupts vs polling

The CPU must learn that a device needs service. Two basic methods exist: polling and interrupts.

Polling means the CPU reads a device status register in a loop or on a timer. If the device is ready, the CPU transfers data. If the device is not ready, the CPU wasted work. Polling is simple. Polling can waste CPU time. Polling can also give low latency when the device is almost always ready and the loop is tight.

An interrupt is a signal to the CPU. The device, a timer, or another CPU raises the interrupt. The CPU stops the current user or kernel code (when interrupts are enabled). The CPU runs an interrupt handler in the kernel. The handler services the device. The CPU then returns to the interrupted code.

Interrupts let the CPU run other work while the device is busy. Disk completion, key presses, and network packets often use interrupts. The kernel can also use polling in a fast path, for example busy-poll on a network queue. The two methods can mix.

Interrupts have a cost. Each interrupt causes a mode change and cache disturbance. Too many interrupts hurt performance. The kernel can use interrupt mitigation: the device sends fewer interrupts and the kernel processes a batch.

A software interrupt or a trap is not a device signal. A system call and a page fault use a trap-like mechanism. They still transfer control to the kernel. Topic 3 and topic 8 cover those paths.

### Questions

#### Theoretical questions

1. What is polling?
2. What is an interrupt?
3. When can polling waste CPU time?
4. Why do many devices use interrupts instead of a tight poll loop?
5. What is one cost of a large number of interrupts?

#### Easy practical tasks

1. Write a two-column table: "Polling" and "Interrupts". Add four comparison rows.
2. Run `cat /proc/interrupts | head`. Write three interrupt names or IRQ numbers that you see.
3. Open `man 5 proc` and find the `/proc/interrupts` paragraph. Write the purpose of that file.
4. Give one human example of polling (check a mailbox every minute) and one example of an interrupt (a doorbell). Map each example to the CPU.

#### Medium practical tasks

1. Compare two snapshots of `/proc/interrupts` taken ten seconds apart. Write which IRQ counts increased.
2. Draw the interrupt path: device, interrupt controller, CPU, kernel handler, return to the process.
3. Read `man 7 epoll` as a preview of wait-for-I/O in user space. Write how waiting differs from a user-space poll loop on a status file.

#### Advanced practical tasks

1. Read about `NAPI` or interrupt coalescing in Linux network documentation. Write one page: when the kernel polls a network card after an interrupt.
2. Write a small user-space program that polls a file descriptor with `read` in a tight loop and a second version that blocks. Compare CPU use in `top`.

---

## Clock and timers

A computer needs time for two different jobs.

A wall clock is the time of day. User programs read it for logs and certificates. The wall clock can jump when an administrator sets the clock or when NTP corrects the time.

A monotonic clock only moves forward. The kernel uses it for timeouts and for scheduling. A monotonic clock does not jump backward when the wall clock changes.

Hardware timers raise interrupts at a rate or at a set deadline. Older PCs used the Programmable Interval Timer (PIT). Modern systems use local APIC timers, HPET, or architecture timers such as the ARM generic timer. The CPU also has a cycle counter, for example the Time Stamp Counter (TSC) on x86.

The kernel uses timer interrupts to:

- preempt a process that used its time slice
- wake a process that slept for a timeout
- run delayed kernel work

Linux also supports high-resolution timers. Those timers do not wait for a coarse periodic tick when hardware allows a one-shot timer.

User programs use `clock_gettime`, `sleep`, `nanosleep`, and `select` timeouts. Those calls ask the kernel. The kernel programs the hardware timer.

Do not mix "clock frequency of the CPU" with "timer interrupt rate". The CPU frequency is how fast the core executes. The timer rate is how often the timer interrupts the CPU.

### Questions

#### Theoretical questions

1. What is the difference between a wall clock and a monotonic clock?
2. Why does the kernel need a hardware timer for scheduling?
3. What is a timer interrupt?
4. Why can the wall clock jump?
5. How does a user program wait for a timeout without a busy loop?

#### Easy practical tasks

1. Run `date`. Write the wall-clock time that the command prints.
2. Run `cat /proc/timer_list | head` if the file exists, or `man 2 clock_gettime`. Write one clock name (`CLOCK_REALTIME` or `CLOCK_MONOTONIC`).
3. Sleep with `sleep 2` and write how you know the kernel used a timer.
4. Open `man 7 time`. Write the difference between real time and CPU time in one sentence each.

#### Medium practical tasks

1. Write a small C program that calls `clock_gettime` for `CLOCK_REALTIME` and `CLOCK_MONOTONIC`. Print both values. Run it twice. Write which clock is safe for interval measurement.
2. Run `cat /proc/interrupts` and find a timer-related line if one exists. Write the name.
3. Draw a timeline: a process sleeps for 100 ms, a timer fires, the kernel wakes the process.

#### Advanced practical tasks

1. Read `man 2 clock_gettime` and `man 7 vdso`. Write how Linux can read some clocks without a full system-call trap.
2. Compare tick-based timers and high-resolution timers in kernel documentation. Write when a periodic tick is enough and when a one-shot timer is better.

---

## Privilege rings and CPU modes

The CPU enforces privilege. On x86, privilege rings are numbered. Ring 0 is the most privileged mode. The kernel runs in ring 0. Ring 3 is the least privileged mode. User programs run in ring 3. Rings 1 and 2 are rare in general-purpose OS designs.

Other architectures use other names. AArch64 uses exception levels. EL0 is user mode. EL1 is the kernel. EL2 is a hypervisor. The idea is the same: some instructions and some registers are illegal in the low-privilege mode.

When a user program needs a privileged service, it executes a system-call instruction. On x86-64 the common instruction is `syscall`. The CPU switches to kernel mode, sets a kernel stack, and jumps to a kernel entry. After the service, the kernel returns with `sysret` or an equivalent instruction.

Hardware also raises exceptions: divide by zero, invalid opcode, page fault. Those events enter the kernel. The kernel can send a signal, kill the process, or fix the fault (for example, demand paging in topic 9).

A hypervisor can run the kernel in a less privileged level and trap sensitive instructions. Topic 16 covers virtualization. This section only states that CPU modes exist so that the OS can isolate user code.

If every program ran in ring 0, any bug could overwrite the kernel and all other programs. Privilege modes are the hardware base for the kernel and user-space split in topic 1.

### Questions

#### Theoretical questions

1. What is ring 0 on x86?
2. What is ring 3 on x86?
3. What is an exception level on AArch64 at a high level?
4. What does the `syscall` instruction do?
5. Why is a page fault handled in the kernel and not in the user program?

#### Easy practical tasks

1. Write five sentences that describe CPU privilege. Use only facts from this section.
2. Make a table: "Mode" and "Who runs there". Add rows for kernel and user programs. Add a row for a hypervisor if you include EL2 or VMX as a preview.
3. Open `man 2 syscall`. Write the difference between the `syscall` C wrapper and the CPU instruction, in your own words.
4. List three events that enter the kernel: system call, device interrupt, page fault.

#### Medium practical tasks

1. Draw x86 rings as concentric circles. Place the kernel and a user process. Draw an arrow labeled `syscall`.
2. Run `cat /proc/cpuinfo | grep -E 'flags|Features'` and search for virtualization flags such as `vmx` or `svm`. Write whether the CPU advertises them.
3. Read `man 7 signal` for `SIGSEGV`. Write how a bad user-space access can become a signal after a kernel exception.

#### Advanced practical tasks

1. Read a short overview of x86 protection rings and AArch64 exception levels. Write a one-page comparison table.
2. Use `gdb` or `objdump` on a tiny program and find a `syscall` instruction in the C library stub if you can. Write the function name and why user code does not stay in ring 0.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A process runs on a core, reads a file from an SSD, and sleeps until a timer. Which hardware parts take part, and in what order?
2. Why must the OS treat RAM and an SSD as different resources?
3. How do interrupts, privilege modes, and device drivers work together for one key press?
4. Why does the kernel need both a cycle counter and a programmable timer?
5. What hardware facts from this topic does virtual memory (topic 8) need?

#### Easy practical tasks

1. Write a one-page cheat sheet: CPU, register, ALU, RAM, SSD, HDD, bus, interrupt, poll, timer, ring 0, ring 3.
2. Collect output from `lscpu`, `free -h`, `lsblk`, and `cat /proc/interrupts | head`. Attach a one-line comment to each file.
3. Draw one diagram that shows a core, RAM, a bus, a device, and an interrupt line.
4. Open `man 7 bootparam` or `man 7 cpuset` only if present; otherwise open `man 1 lscpu`. Write two hardware facts that the OS exports to user space.

#### Medium practical tasks

1. Write a script that saves `lscpu`, `free -h`, and `lsblk -f` into a report file with timestamps. Run it on your VM or WSL.
2. Compare two Linux systems (host and VM, or two VMs). Write five hardware differences that `/proc` or `lscpu` shows.
3. Design a table that maps each hardware section in this topic to one later OS topic (boot, processes, scheduling, memory, I/O).

#### Advanced practical tasks

1. Read the OSTEP chapter on limited direct execution or hardware support. Write how the OS uses privilege modes and interrupts to run user code safely.
2. Profile one second of idle system versus one second of `dd` to a file. Use `/proc/interrupts` and `vmstat 1 3`. Write which hardware-related counts move.
