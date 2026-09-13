# 8. Devices and I/O

## Description

Programs do not talk to device registers. The kernel and a device driver do that work. This topic explains drivers, the interrupt path, and direct memory access (DMA). You learn blocking I/O and non-blocking I/O. You also get a survey of `epoll`, `kqueue`, and IOCP, and an awareness of Linux `io_uring`.

Complete this topic after filesystems and disk. Complete this topic before inter-process communication (topic 9). You already know interrupts from topic 1 and files from topic 6.

Use one term for each concept. A device driver is kernel code that controls one device class. An interrupt handler is the short kernel path that runs when the device signals the CPU. DMA is a device transfer to or from RAM without a CPU copy of each word. Blocking I/O waits in the kernel until the operation can complete. A readiness API tells you when a file descriptor is ready. Do not mix DMA with a system call. Do not mix blocking I/O with a sleep in user code.

---

## Drivers, interrupts, DMA

A device driver is kernel code that presents a device to the rest of the kernel. The driver programs registers, handles interrupts, and often queues DMA. User programs see a file, a socket, or a block device. They do not see the registers.

Linux groups many drivers as modules. You can load a module with `modprobe` and list modules with `lsmod`. Some drivers are built into the kernel image.

Driver classes:

1. Character drivers: byte streams and `ioctl` (terminals, many sensors).
2. Block drivers: request queues of reads and writes of blocks (disks, NVMe).
3. Network drivers: packets through the network stack (topic 12), not a normal file `read` of frames.

The kernel talks to a driver through a standard set of operations. For a file, those operations include `open`, `read`, `write`, `ioctl`, and `mmap`. VFS or the block layer calls them (topic 7).

A device raises an interrupt when it needs the kernel (data ready, DMA done, error). The CPU stops the current instruction stream and runs an interrupt handler in kernel mode.

Typical steps:

1. The device asserts an interrupt line or writes an MSI-X message.
2. The CPU enters the kernel. The kernel finds the handler.
3. The top half (hard IRQ) does little work: acknowledge the device, wake a queue, or schedule a bottom half.
4. The bottom half (softirq, tasklet, or workqueue) does more work and may sleep in some contexts.
5. The kernel returns to the interrupted thread, or it schedules another thread.

Keep the hard IRQ short. Sleep is forbidden in many interrupt contexts. A driver bug can crash the whole machine (topic 11).

Direct memory access (DMA) lets the device read or write RAM without a CPU load and store of each word. The driver sets up a DMA descriptor. The device transfers a block. The device interrupts when the transfer is done. The CPU is free for other work during the transfer.

DMA needs coherent addresses and constraints (alignment, bounce buffers). IOMMU hardware can isolate device DMA. That is a later security topic. This section only states that DMA is a device-to-RAM path.

User space talks to many devices through `/dev` or through `sysfs` (`/sys`). `ioctl` is the escape hatch for commands that are not `read` or `write`. Prefer a documented interface. Do not send random `ioctl` codes.

Do not write a kernel driver for this path. Read `lsmod` and `/sys` and treat the driver as a black box with a contract.

### Questions

#### Theoretical questions

1. What is a device driver?
2. What are the three driver classes in this section?
3. What is the difference between the interrupt top half and the bottom half?
4. What is DMA?
5. Why can a driver bug crash the whole machine?

#### Easy practical tasks

1. Run `lsmod | head`. Write three module names.
2. Open `man 8 modprobe` and `man 2 ioctl`. Write one sentence for each.
3. Run `ls /sys/class`. Write five class names.
4. Write five sentences about drivers, interrupts, and DMA. Use only facts from this section.

#### Medium practical tasks

1. Pick one module from `lsmod`. Run `modinfo <name>` if the tool exists. Write the description field.
2. Draw: user `read`, VFS, driver, interrupt, DMA, RAM.
3. Read `man 5 proc` or sysfs notes for `/sys/class/net`. Write how a network device appears without a `/dev` file.

#### Advanced practical tasks

1. Read the Linux kernel driver model overview (kernel documentation). Write a one-page summary of probe, remove, and a device class.
2. Compare a USB storage device path (block) with a USB serial path (character) in a table: `/dev` name, module, and user API.

---

## Blocking vs non-blocking I/O

Blocking I/O is the default for many FDs. A `read` that needs data that is not ready puts the thread to sleep in the kernel. The thread is not on the CPU. When the data arrives (interrupt, DMA done, or peer write), the kernel wakes the thread. `write` to a full pipe or socket can also block.

Non-blocking I/O means the call returns at once if it cannot complete. You set `O_NONBLOCK` on the FD (`fcntl` or `open`). `read` returns `-1` with `EAGAIN` or `EWOULDBLOCK` when no data is ready. `write` can return `EAGAIN` when the buffer is full. The thread stays runnable. You must retry later.

Blocking I/O is simple. One thread per connection is easy to write and can waste memory and context switches at large scale. Non-blocking I/O needs a readiness loop. You do not sleep in `read`. You wait in `epoll` or a similar call (next section).

`O_NONBLOCK` is not the same as asynchronous I/O. Classic non-blocking I/O still copies in the system call when the FD is ready. The thread still spends time in `read`. True async completion (io_uring, some `aio` paths) can start I/O and collect results later.

A timeout is a third style. `setsockopt` `SO_RCVTIMEO` or `alarm` can limit a blocking `read`. `poll` has a timeout argument.

Do not spin on `EAGAIN` without a wait. A busy loop burns a core. Do not set `O_NONBLOCK` and then ignore `EAGAIN`.

### Questions

#### Theoretical questions

1. What does a blocking `read` do when no data is ready?
2. What does `O_NONBLOCK` change about `read`?
3. What errno values mean "try again"?
4. Why is a busy loop on `EAGAIN` a bad design?
5. How does non-blocking I/O differ from asynchronous completion?

#### Easy practical tasks

1. Open `man 2 fcntl` and `man 2 read`. Write how you set `O_NONBLOCK`.
2. Write a two-column table: blocking and non-blocking. Add three rows.
3. Write five sentences about blocking and non-blocking I/O. Use only facts from this section.
4. Run `cat` with no arguments (stdin). Write that it blocks. Stop with Control-D or Control-C.

#### Medium practical tasks

1. Write a program that sets a pipe read end to non-blocking and `read`s with no writer. Write the errno name. This handbook does not contain the source.
2. Draw a thread in `S` on a blocking `read`, then an interrupt, then a wake.
3. Compare `sleep` in user code with a blocking `read`. Write which one the kernel I/O path uses.

#### Advanced practical tasks

1. Add a timeout to a blocking socket `read` with `SO_RCVTIMEO` or `poll`. Document the return on timeout.
2. Write a one-page note: one thread per client versus one thread plus non-blocking FDs. Give a scale where each design is enough.

---

## `epoll` / `kqueue` / IOCP (survey)

A readiness API waits until one or more file descriptors can `read` or `write` without blocking (or until an error, hangup, or timeout).

`select` and `poll` are old POSIX calls. `select` has an FD number limit on many systems. Both copy the FD list into the kernel on each call. They still work for small sets.

`epoll` is the Linux readiness API for many FDs. You create an epoll instance (`epoll_create1`). You add FDs with `epoll_ctl`. You wait with `epoll_wait`. The kernel tracks the set. Edge-triggered and level-triggered modes exist. Level-triggered is easier: if the FD is still readable, `epoll_wait` reports it again. Edge-triggered reports a change. With edge-triggered mode you must drain the FD or you can stall.

`kqueue` is the BSD and macOS readiness API. You register filters (read, write, vnode, timer). The idea matches `epoll`: one wait set, many FDs, events out.

IOCP (I/O completion ports) is the usual Windows model for high-scale I/O. It is a completion API more than a readiness API. You start an operation. The port reports when the operation is done. The thread that waits is not the same design as `epoll_wait` plus `read`.

Survey rules:

- Linux servers: learn `epoll`.
- BSD and macOS: learn `kqueue`.
- Windows: learn IOCP if you write native servers.
- Portable libraries (libuv, Go netpoller) hide the choice.

`epoll` does not work on all FD types in the same way. Regular files on Linux are often "always ready". Use `epoll` for sockets, pipes, and many devices. Topic 12 uses `epoll` on sockets.

This section is a survey. You implement a small `epoll` loop in practice tasks. You do not implement IOCP on Linux.

### Questions

#### Theoretical questions

1. What does a readiness API tell you?
2. Why is `epoll` cheaper than `select` for a large FD set?
3. What is the difference between level-triggered and edge-triggered `epoll`?
4. What system uses `kqueue`?
5. How does IOCP differ from `epoll` at a high level?

#### Easy practical tasks

1. Open `man 7 epoll` and `man 2 poll`. Write one sentence for each.
2. Write a three-row table: `epoll`, `kqueue`, IOCP. Add the usual OS.
3. Write five sentences about these APIs. Use only facts from this section.
4. Run `man 2 epoll_wait` and write the return value on timeout (0).

#### Medium practical tasks

1. Write a program with a pipe: the parent `epoll_wait`s on the read end, the child writes a line after one second. Print the line. This handbook does not contain the source.
2. Draw `epoll_create`, `epoll_ctl` add, `epoll_wait`, `read`.
3. Read `man 7 epoll` on `EPOLLET`. Write one rule that you must follow in edge-triggered mode.

#### Advanced practical tasks

1. Add a second client FD (two pipes or two sockets) to the same epoll set. Document how you map the event to the FD.
2. Read a Go or libuv netpoller overview. Write which OS API each platform uses.

---

## io_uring (awareness)

`io_uring` is a Linux I/O interface. The process maps two ring buffers that it shares with the kernel: a submission queue (SQ) and a completion queue (CQ). The process posts I/O requests on the SQ. The kernel posts results on the CQ. A system call (`io_uring_enter`) can submit and wait. In some modes the kernel polls the SQ and you avoid a call per I/O.

Goals:

1. Fewer system calls than one `read` per buffer.
2. Completion, not only readiness: you start a `read` and later you see that it finished.
3. One API for files, sockets, and other operations (the set grows over kernel versions).

`io_uring` is more complex than `epoll`. You must manage rings, memory, and completion order. Security bugs in early versions led to restricted defaults on some systems. Use a current kernel and a maintained library (`liburing`) if you write real code.

This path only asks for awareness. You do not need to ship `io_uring` in a first server. Learn blocking I/O, then `epoll`, then read an `io_uring` overview.

`io_uring` does not replace the driver or DMA. It is a user-kernel submission path. The device still interrupts. DMA still moves bytes.

Do not copy an old blog that uses unsafe flags. Read the current `man 2 io_uring_setup` notes and the kernel documentation.

### Questions

#### Theoretical questions

1. What two rings does `io_uring` use?
2. How does completion differ from `epoll` readiness?
3. Why can `io_uring` reduce system-call count?
4. Why must a beginner still learn `epoll` first?
5. Does `io_uring` replace DMA? Explain.

#### Easy practical tasks

1. Open `man 2 io_uring_setup` if the page exists, or the kernel `io_uring` documentation online. Write one sentence about the setup call.
2. Write a two-column table: `epoll` and `io_uring`. Add "readiness or completion" and "typical first use".
3. Write five sentences about `io_uring`. Use only facts from this section.
4. Run `uname -r`. Write whether your kernel is 5.1 or newer (io_uring arrived around then). Do not enable experimental flags.

#### Medium practical tasks

1. Read the `liburing` README overview. Write the names of submit and wait helpers.
2. Draw process SQ, kernel, CQ, and one `read` request.
3. Compare a blog from 2019 with current man pages. Write one warning that changed (security or API).

#### Advanced practical tasks

1. Write a one-page awareness note: three operations that `io_uring` can submit, and three reasons you would still use `read` plus `epoll`.
2. If your kernel and policy allow it, build the `liburing` echo example from upstream docs. Do not invent flags. Document the kernel version.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a driver, an interrupt, and DMA complete one blocking `read` from a device?
2. When do you pick blocking I/O, `O_NONBLOCK` plus `epoll`, or `io_uring`?
3. Why is `epoll` on Linux not the same design as IOCP on Windows?
4. What goes wrong if the hard IRQ handler sleeps?
5. Why are regular files a poor fit for an `epoll`-only server loop?

#### Easy practical tasks

1. Write a cheat sheet: driver, module, IRQ top/bottom, DMA, blocking, `O_NONBLOCK`, `EAGAIN`, `epoll`, `kqueue`, IOCP, `io_uring`.
2. Run `lsmod | wc -l`, `ls /dev | wc -l`, and `man 7 epoll` (open the page). Write one line each.
3. Draw a server: many sockets, one `epoll_wait`, no busy loop.
4. Bookmark `man 7 epoll` and the io_uring kernel docs.

#### Medium practical tasks

1. Combine a non-blocking pipe and `epoll_wait` in one small program. Write a six-line report of the event order.
2. Use `strace -e epoll_wait,read,write` on your program. Write the call list.
3. Document a readiness checklist: set non-block, add to epoll, wait, handle `EAGAIN`, never spin.

#### Advanced practical tasks

1. Read Linux `Documentation/core-api/irq/` or an interrupt overview. Write a one-page map from device IRQ to a waking `wait_queue`.
2. Design (on paper) a Linux TCP echo server: `epoll` first. Mark where `io_uring` could replace `read` later. Topic 12 adds the socket calls.
