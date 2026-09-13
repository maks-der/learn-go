# 14. I/O and Devices

## Description

Programs do not talk to device registers. The kernel and a device driver do that work. This topic explains drivers, the interrupt path, and direct memory access (DMA). You learn blocking I/O and non-blocking I/O. You also get a survey of `epoll`, `kqueue`, and IOCP, and an awareness of Linux `io_uring`.

Complete this topic after storage devices. Complete this topic before inter-process communication (topic 15). You already know interrupts from topic 2 and files from topic 11. You now learn how a `read` waits, and how a server waits on many sockets.

Use one term for each concept. A device driver is kernel code that controls one device class. An interrupt handler is the short kernel path that runs when the device signals the CPU. DMA is a device transfer to or from RAM without a CPU copy of each word. Blocking I/O waits in the kernel until the operation can complete. A readiness API tells you when a file descriptor is ready. Do not mix DMA with a system call. Do not mix blocking I/O with a sleep in user code.

---

## Device drivers

A device driver is kernel code that presents a device to the rest of the kernel. The driver programs registers, handles interrupts, and often queues DMA. User programs see a file, a socket, or a block device. They do not see the registers.

Linux groups many drivers as modules. You can load a module with `modprobe` and list modules with `lsmod`. Some drivers are built into the kernel image. Topic 20 covers modules as a tracing topic. This section states the role.

Driver classes:

1. Character drivers: byte streams and `ioctl` (terminals, many sensors).
2. Block drivers: request queues of reads and writes of blocks (disks, NVMe).
3. Network drivers: packets through the network stack (topic 18), not a normal file `read` of frames.

The kernel talks to a driver through a standard set of operations. For a file, those operations include `open`, `read`, `write`, `ioctl`, and `mmap`. VFS or the block layer calls them (topic 12).

A driver must not sleep in a context that forbids sleep (some interrupt code). A driver must not race on shared device state. Those rules are why driver bugs crash the machine (topic 17).

User space talks to many devices through `/dev` or through `sysfs` (`/sys`). `ioctl` is the escape hatch for commands that are not `read` or `write`. Prefer a documented interface. Do not send random `ioctl` codes.

Firmware for a device can load from user space (`firmware` files). The driver still owns the device.

Do not write a kernel driver for this path. Read `lsmod` and `/sys` and treat the driver as a black box with a contract.

### Questions

#### Theoretical questions

1. What is a device driver?
2. What are the three driver classes in this section?
3. How does a user program reach a character driver?
4. Why can a driver bug crash the whole machine?
5. What is `ioctl` for?

#### Easy practical tasks

1. Run `lsmod | head`. Write three module names.
2. Open `man 8 modprobe` and `man 2 ioctl`. Write one sentence for each.
3. Run `ls /sys/class`. Write five class names.
4. Write five sentences about drivers. Use only facts from this section.

#### Medium practical tasks

1. Pick one module from `lsmod`. Run `modinfo <name>` if the tool exists. Write the description field.
2. Draw: user `read`, VFS, driver, device registers.
3. Read `man 5 proc` or sysfs notes for `/sys/class/net`. Write how a network device appears without a `/dev` file.

#### Advanced practical tasks

1. Read the Linux kernel driver model overview (kernel documentation). Write a one-page summary of probe, remove, and a device class.
2. Compare a USB storage device path (block) with a USB serial path (character) in a table: `/dev` name, module, and user API.

---

## Interrupt handling path

A device raises an interrupt when it needs the kernel (data ready, DMA done, error). The CPU stops the current instruction stream and runs an interrupt handler in kernel mode. Topic 2 introduced this idea. This section names the Linux path.

Typical steps:

1. The device asserts an interrupt line or an MSI/MSI-X message.
2. The interrupt controller delivers the event to a CPU.
3. The CPU runs a small top-half handler. That handler acknowledges the device and records work. It must stay short.
4. The kernel schedules a bottom half: a softirq, a tasklet, or a workqueue. That code can do more work and may copy data or wake a process.
5. A blocked process in `read` or a network receive path wakes and later returns to user space.

The split exists because an interrupt can preempt almost anything. Long work in the top half delays other interrupts and the scheduler.

Shared interrupt lines were common on old PIC systems. MSI-X gives many devices their own messages. `/proc/interrupts` counts events per CPU.

A process does not install a C function as a device interrupt handler. The driver does. User programs can wait on file descriptors or signals (`signalfd`, topic 15).

Polling is the other design: the CPU asks the device in a loop. Topic 2 compared polling and interrupts. High-speed NICs sometimes poll in a NAPI loop to reduce interrupt rate. That is an optimization inside the driver.

### Questions

#### Theoretical questions

1. What starts the interrupt handling path?
2. Why must the top-half handler stay short?
3. What is a bottom half for?
4. What does `/proc/interrupts` show?
5. How does NAPI change interrupt rate at a high level?

#### Easy practical tasks

1. Run `cat /proc/interrupts | head`. Write two IRQ names that you recognize (timer, keyboard, or a disk).
2. Open `man 5 proc` for `/proc/interrupts`. Write one sentence.
3. Write five sentences about the interrupt path. Use only facts from this section.
4. Draw the five steps from device assert to a woken `read`.

#### Medium practical tasks

1. Run `watch -n 1 cat /proc/interrupts` for a few seconds while you type or move the mouse. Write which counters increase.
2. Compare interrupt I/O with polling in a six-sentence table story: one key press versus a tight `inb` loop (idea only).
3. Read a short Linux NAPI or softirq note. Write which work moves out of the top half.

#### Advanced practical tasks

1. Read about threaded IRQs in Linux. Write a one-page contrast: hard IRQ handler versus a threaded handler that can sleep.
2. Use `perf` (topic 20) to record `irq` events if you have rights. Write two symbol names. Do this on a VM.

---

## DMA

Direct memory access (DMA) lets a device read or write RAM without a CPU load or store for each word. The CPU sets up a descriptor: address, length, and direction. The device (or a DMA engine) moves the bytes. An interrupt often signals completion.

Without DMA, programmed I/O (PIO) uses the CPU to copy each word from a device register. PIO wastes CPU time on large transfers.

The kernel must give the device addresses that the device can use. Those addresses are I/O-virtual or bus addresses, not a user virtual address. The IOMMU (if present) maps device-visible addresses and can restrict what the device may touch. That restriction is a security feature.

The driver must pin or map pages for the life of the transfer. A page must not move or go to swap while the device writes it. The kernel DMA API handles mapping.

Cache coherence is a hardware and driver problem. The CPU cache and the device must agree on the bytes in RAM. Architectures differ. Drivers call mapping APIs that flush or invalidate caches when needed.

User programs do not start DMA. `read` and `write` may cause the driver to start DMA. `splice` and some `io_uring` paths try to reduce extra copies after DMA.

A buggy device or a missing IOMMU can overwrite RAM. That is one reason firmware and drivers are trusted code.

### Questions

#### Theoretical questions

1. What is DMA?
2. What is PIO?
3. Why is a user virtual address not a DMA address?
4. What extra job does an IOMMU do?
5. Why must pages stay pinned during a DMA transfer?

#### Easy practical tasks

1. Open `man 7 dma` if it exists, or read a kernel DMA-API note online. Write one sentence about mapping.
2. Write five sentences about DMA. Use only facts from this section.
3. Draw: CPU, RAM, device, and an arrow for DMA that does not pass through the CPU.
4. Make a table: PIO versus DMA. Add CPU cost and typical use.

#### Medium practical tasks

1. Read `/sys/class/iommu` or `dmesg` for IOMMU messages if present. Write whether an IOMMU looks active.
2. Trace the idea of one disk `read`: DMA into a page-cache page, then copy or map to the process. Draw it.
3. Read a short IOMMU introduction. Write six sentences on isolation of devices.

#### Advanced practical tasks

1. Read `Documentation/core-api/dma-api.rst` (online). Write a one-page summary of coherent versus streaming DMA mappings.
2. Write why a user-space driver framework (VFIO, UIO) still needs the kernel to map DMA. No exploit steps. Policy and API names only.

---

## Blocking vs non-blocking I/O

A blocking `read` on a pipe, socket, or terminal waits until data exists (or an error or hangup occurs). The process state is blocked (topic 4). The scheduler runs other programs.

A blocking `write` waits until the kernel accepts the bytes (buffer space).

Non-blocking mode (`O_NONBLOCK` on `open` or `fcntl`) changes the rule. If the operation cannot complete now, the call returns `-1` and `errno` is `EAGAIN` or `EWOULDBLOCK`. The process keeps the CPU. It must retry later.

Typical uses:

- Blocking: simple tools, one client, scripts. Easy to write.
- Non-blocking: one thread that must watch many descriptors, or a thread that must not stall on one slow peer.

Non-blocking I/O is not the same as asynchronous I/O. Non-blocking means "return now if you cannot finish". Asynchronous I/O means "start the work and notify me later" (`aio`, `io_uring`). A non-blocking `read` still copies data in the same call when data is ready.

`select`, `poll`, and `epoll` wait until at least one descriptor is ready. Then you call `read` or `write`. That pattern is readiness I/O.

A regular file on a local disk often does not block in the same way as a socket. `O_NONBLOCK` on a regular file does not turn disk I/O into a guaranteed instant return on Linux. Disk waits still happen inside the call. Use `io_uring` or `aio` if you need async file I/O.

Always handle short reads and short writes. Always handle `EINTR`.

### Questions

#### Theoretical questions

1. What does a blocking `read` do when no data is ready?
2. What does `O_NONBLOCK` change?
3. What is the difference between non-blocking I/O and asynchronous I/O?
4. What is readiness I/O?
5. Why is `O_NONBLOCK` on a regular file a weak tool for disk async on Linux?

#### Easy practical tasks

1. Open `man 2 fcntl` and find `O_NONBLOCK`. Write how you set the flag.
2. Open `man 2 read` and find `EAGAIN`. Write the meaning.
3. Write five sentences about blocking and non-blocking I/O. Use only facts from this section.
4. Make a table: blocking `read` on a pipe versus non-blocking `read` on an empty pipe.

#### Medium practical tasks

1. Write two processes and a pipe. Reader uses blocking `read`. Writer sleeps, then writes. Document the reader wait.
2. Set `O_NONBLOCK` on the read end. Read before the write. Write the errno.
3. Draw a server thread that must not block on one client. Mark where `EAGAIN` appears.

#### Advanced practical tasks

1. Read `man 7 pipe` and `man 7 socket` on blocking. Write a one-page note on `SIGPIPE` and `EPIPE` versus `EAGAIN`.
2. Compare POSIX AIO (`aio_read`) with non-blocking sockets in a table: completion, file support, and typical use on Linux.

---

## `epoll` / `kqueue` / IOCP (survey)

A server can have thousands of idle connections. One thread per connection wastes memory. A readiness or completion API lets few threads wait on many descriptors.

`epoll` is the Linux API. You create an epoll instance (`epoll_create1`). You add descriptors with `epoll_ctl`. You wait with `epoll_wait`. The kernel returns a list of events (readable, writable, hangup, error). Edge-triggered and level-triggered modes exist. Edge-triggered mode requires you to drain the descriptor until `EAGAIN`.

`kqueue` is the BSD and macOS API. You create a queue and register filters (read, write, vnode, timer, signal). `kevent` waits. The idea matches `epoll`. The calls differ.

IOCP (I/O completion ports) is the Windows model. You start an overlapped operation. The kernel posts a completion packet when the I/O finishes. This is completion, not only readiness. The thread pool waits on the port.

Survey facts:

- Linux user programs use `epoll` (or `io_uring`).
- BSD and macOS use `kqueue`.
- Windows uses IOCP.
- `select` and `poll` are older and portable. They scale worse with a huge descriptor set.

Libraries (libuv, Go `netpoll`) wrap these APIs. You still must know which model you are on.

This section is a survey. You do not need a production server. You need to name the API and the model (readiness versus completion).

### Questions

#### Theoretical questions

1. What problem do `epoll`, `kqueue`, and IOCP solve?
2. What are the three main `epoll` calls?
3. How does `kqueue` differ from `epoll` at a high level?
4. How does IOCP differ from readiness I/O?
5. Why does `select` scale poorly with many idle sockets?

#### Easy practical tasks

1. Open `man 7 epoll` and `man 2 epoll_wait`. Write one sentence for each.
2. Write a three-row table: Linux, BSD/macOS, Windows. Add the API name.
3. Write five sentences about this survey. Use only facts from this section.
4. Draw level-triggered versus edge-triggered as two timelines of "data available".

#### Medium practical tasks

1. Write a small `epoll` program that waits on `STDIN_FILENO` and prints when it is readable. Type a line.
2. Read `man 2 select` and compare the `fd_set` size idea with `epoll`. Write six sentences.
3. Read a Go or nginx document on the poller. Write which OS API it uses on Linux.

#### Advanced practical tasks

1. Implement a tiny echo server with `epoll` and non-blocking sockets for two clients. Document edge-triggered pitfalls if you use them.
2. Read a Windows IOCP overview (public docs). Write a one-page contrast with `epoll` readiness. Do not write Windows exploit code.

---

## io_uring (Linux, awareness)

`io_uring` is a Linux asynchronous I/O interface. The process and the kernel share two ring buffers in memory: a submission queue (SQ) and a completion queue (CQ). The process posts requests on the SQ. The kernel posts results on the CQ. A system call (`io_uring_enter`) can submit and wait. In some modes the kernel pulls submissions with less need for a call on every I/O.

Goals:

- Async file I/O that works, unlike many POSIX AIO paths on Linux.
- Fewer system calls for a busy server.
- A single API for sockets, files, and other operations (the set grows with kernel version).

You set up a ring with `io_uring_setup`. Libraries such as `liburing` hide the raw mmap of the rings.

This path requires awareness, not mastery. Kernels and distributions differ. Features depend on the kernel version. Some operations still fall back or fail if the filesystem or the operation does not support the path.

Security and complexity are real. The ring is shared memory with the kernel. Keep the kernel updated. Do not treat `io_uring` as a toy on a shared host without reading the distribution notes.

For a beginner server, `epoll` plus non-blocking sockets is enough. Learn `io_uring` when you need async disk I/O or you measure syscall cost.

### Questions

#### Theoretical questions

1. What two rings does `io_uring` use?
2. What problem does `io_uring` solve that `epoll` does not solve well?
3. What does `io_uring_setup` do?
4. Why does this handbook treat `io_uring` as awareness?
5. Why can a shared ring with the kernel be a security concern at a high level?

#### Easy practical tasks

1. Open `man 2 io_uring_setup` if the page exists. Write the purpose in one sentence. If the page is missing, write that the kernel headers or man pages are old.
2. Run `uname -r`. Write the kernel release (features depend on it).
3. Write five sentences about `io_uring`. Use only facts from this section.
4. Make a table: `epoll`, POSIX AIO, `io_uring`. Add "typical I/O type".

#### Medium practical tasks

1. Search whether `liburing` is in your distribution. Write the package name or that it is absent.
2. Draw SQ, kernel, device, CQ, and a user thread that harvests completions.
3. Read a short official `io_uring` overview (kernel docs or man-pages). Write six sentences on submission versus completion.

#### Advanced practical tasks

1. Write a `liburing` program that reads one file asynchronously if your kernel and library allow it. Document the kernel version and every failure.
2. Compare `io_uring` poll mode with `epoll` in a one-page design for a socket server. No need to implement both.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe one disk `read` from a blocking user call to DMA completion and a woken process. Name the driver, the interrupt, and the page cache.
2. When do you choose blocking I/O, `epoll`, and `io_uring` for a student project?
3. How do a character driver and a network driver differ in what a user program opens?
4. Why must the interrupt top half and the DMA completion path stay correct if user programs only call `read`?
5. How do `epoll` and IOCP disagree about when the kernel talks to user space?

#### Easy practical tasks

1. Write a one-page cheat sheet: driver, top half, bottom half, DMA, IOMMU, blocking, `O_NONBLOCK`, `epoll`, `kqueue`, IOCP, `io_uring`.
2. Run `lsmod | wc -l`, `head /proc/interrupts`, `ls /dev | wc -l`, and `man 7 epoll` (open the page). Write one line each.
3. Draw one figure that contains PIO, DMA, and a blocking `read` wait queue.
4. Bookmark `man 7 epoll` and the `io_uring` man page if present.

#### Medium practical tasks

1. Write a program that blocks on a pipe, then a second version with `poll` or `epoll` on the pipe and stdin. Document both.
2. Use `strace -e read,write,poll,epoll_wait` on `cat` and on a small Python or Go HTTP server if you have one. Write which wait API appears.
3. Map topic 2 (interrupts, DMA) and topic 13 (block devices) onto this topic in a short table of terms.

#### Advanced practical tasks

1. Read the OSTEP chapters on I/O devices and on event-based servers. Write a one-page map to each section in this handbook.
2. In a VM, write an `epoll` echo server and load it with many idle connections (a small client loop). Record RSS and CPU. Then write what you would change to try `io_uring` later.
