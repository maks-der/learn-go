# 16. Virtualization

## Description

Virtualization lets one machine run more than one operating system, or lets one kernel isolate many user environments. This topic explains why virtual machines (VMs) exist, hypervisor types, traps and emulation, and hardware virtualization on x86. You learn memory virtualization (EPT/NPT), the difference between containers and VMs, and Linux cgroups and namespaces.

Complete this topic after IPC. Complete this topic before security and isolation (topic 17). You already know privilege modes, page tables, and mount namespaces as a preview. You now see a full isolation stack.

Use one term for each concept. A hypervisor is the software that runs virtual machines. A guest is the OS inside a VM. A host is the system that provides the hypervisor or the shared kernel. A container is an isolated user-space environment that shares the host kernel. A namespace is a per-process view of one kernel resource. A cgroup is a kernel limiter and accountant for resource use. Do not mix a container with a VM. Do not mix type 1 and type 2 hypervisors.

---

## Why VMs exist

A virtual machine is a software illusion of a computer. The guest OS runs as if it owned a CPU, memory, and devices. The hypervisor multiplexes the real hardware among guests.

Reasons to use a VM:

1. Isolation: a crash or a compromised guest should not take down other guests (the goal; bugs exist).
2. Density: many servers share one physical machine.
3. Compatibility: an old OS or a different kernel can run on new hardware.
4. Snapshots and migration: you can copy or move a guest disk and sometimes live memory.
5. Teaching and test: you can break a guest without a second physical PC (topic 1).

A VM is not the only isolation tool. A process already isolates address spaces. A VM isolates kernels. That extra layer costs CPU, memory, and I/O overhead. You pay that cost when you need a different kernel, a different OS, or a strong trust boundary.

Cloud vendors sell VMs as the default unit. The same ideas apply to a local QEMU/KVM guest.

A VM does not replace backups. A VM does not make an unsafe guest program safe inside the guest. The guest still needs its own user IDs and updates.

### Questions

#### Theoretical questions

1. What illusion does a VM give to a guest OS?
2. Name three reasons to run a VM.
3. What extra isolation does a VM add over a process?
4. What cost do you pay for that isolation?
5. Why is a VM not a backup?

#### Easy practical tasks

1. Write five sentences about why VMs exist. Use only facts from this section.
2. Run `systemd-detect-virt` or `cat /proc/cpuinfo | grep -i hypervisor` if present. Write whether you look like a guest.
3. Open the QEMU or VirtualBox overview page. Write one sentence about what the product virtualizes.
4. Make a table: process isolation versus VM isolation. Add "own kernel?".

#### Medium practical tasks

1. Draw a physical machine, a hypervisor, and two guests. Label CPU, RAM, and disk.
2. Compare a cloud VM and a local lab VM in six sentences: who owns the hypervisor, and what you can snapshot.
3. Read a short note on live migration (idea only). Write what must move (memory, disk, network identity).

#### Advanced practical tasks

1. Read the first chapter of a KVM or hypervisor intro (kernel docs or a textbook). Write a one-page list of resources that a guest thinks it owns.
2. Write a policy: when your course work must use a VM instead of WSL or a container (kernel modules, crash tests).

---

## Hypervisor type 1 vs 2

A hypervisor (virtual machine monitor) is the software that creates and runs VMs.

A type 1 hypervisor runs on the hardware. It does not need a full general-purpose host OS under it. Examples in this class: Xen, ESXi, Hyper-V in its bare-metal role, and the idea of Linux KVM as a kernel module that makes the kernel a hypervisor. The line is a bit blurry on Linux because Linux is also a full OS.

A type 2 hypervisor runs as a program on a host OS. VirtualBox and VMware Workstation are typical examples. QEMU can emulate in user space. QEMU plus KVM uses the host kernel as the hypervisor and QEMU as the user-space helper for devices and the GUI.

Comparison:

- Type 1: less extra host software in the path, common in data centers.
- Type 2: easy to install on a laptop, good for labs, extra overhead and extra attack surface in the host OS.

KVM (Kernel-based Virtual Machine) is a Linux kernel module. It exposes `/dev/kvm`. A user-space process (often QEMU) configures the guest and handles I/O. This model is the usual Linux server hypervisor.

Do not argue about marketing labels. Ask: what runs in the most privileged CPU mode, and what handles device emulation?

### Questions

#### Theoretical questions

1. What is a hypervisor?
2. What is a type 1 hypervisor?
3. What is a type 2 hypervisor?
4. What is KVM on Linux?
5. Why is the type 1 versus type 2 line blurry on Linux?

#### Easy practical tasks

1. Open `man 8 kvm` or `man 1 qemu-system-x86_64` if present. Write one sentence.
2. Run `ls -l /dev/kvm` if the node exists. Write the permissions.
3. Write a two-column table: type 1 and type 2. Add one example each.
4. Write five sentences about hypervisor types. Use only facts from this section.

#### Medium practical tasks

1. Draw VirtualBox on Windows as type 2. Draw KVM plus QEMU on Linux. Label host kernel and guest.
2. Read a KVM API overview. Write six sentences: who calls `ioctl` on `/dev/kvm`.
3. List three device types that user-space QEMU typically emulates (disk, NIC, firmware).

#### Advanced practical tasks

1. Read `Documentation/virt/kvm` (online). Write a one-page map of VM create, vCPU run, and exit to user space.
2. Compare Xen PV versus HVM at a high level in a table. No exploit content.

---

## Traps and emulation

A guest kernel tries to run privileged instructions. On a machine without hardware virtualization, those instructions cannot run in user mode. The CPU traps. The hypervisor emulates the instruction and then resumes the guest.

This model is trap-and-emulate. Classic papers describe it. Some x86 instructions were not classically virtualizable: they behaved differently in unprivileged mode without a trap. Hypervisors then used binary translation or paravirtualization.

Paravirtualization (PV) means the guest OS knows it is a guest. It calls hypercalls instead of touching fake hardware. Xen popularized this idea. PV reduces traps. It needs a modified guest.

Emulation of devices is similar. A guest `out` to an I/O port or a store to an MMIO page exits to the hypervisor. QEMU or the kernel then updates a virtual disk or a virtual NIC.

Exits are expensive. A high exit rate kills performance. The goal of hardware virt and of virtio (paravirtual devices) is fewer exits and faster I/O.

You can think of a VM exit as a controlled trip from guest mode to hypervisor mode. It is not a Linux system call from a user process, but the idea "trap to a more privileged layer" is the same family as topic 1.

### Questions

#### Theoretical questions

1. What is trap-and-emulate?
2. Why did some x86 instructions break classic virtualization?
3. What is a hypercall?
4. What is MMIO emulation?
5. Why are VM exits expensive?

#### Easy practical tasks

1. Write five sentences about traps and emulation. Use only facts from this section.
2. Draw: guest privileged instruction, trap, hypervisor emulate, resume.
3. Look up `virtio` in `man 7 virtio` or a kernel virtio note if present. Write one sentence.
4. Make a table: full emulation versus paravirtual device. Add "guest knows?".

#### Medium practical tasks

1. Read a short note on virtio-net or virtio-blk. Write six sentences: why a virtio disk is faster than a fully emulated IDE disk.
2. Compare a system call (user to Linux) with a VM exit (guest to hypervisor) in a table: who traps, who resumes.
3. List three operations that likely cause exits (timer, I/O port, `hlt` wait).

#### Advanced practical tasks

1. Read a classic virtualization paper abstract (Popek and Goldberg, or a KVM overview). Write the three properties of a virtualizable architecture in your own words.
2. In a KVM guest, read `/proc/cpuinfo` flags and write whether hypervisor and virtio appear. Document the host only if you own it.

---

## Hardware virt (VT-x, AMD-V)

Intel VT-x and AMD-V (SVM) add CPU modes for guests. The processor can run guest kernel code in a guest mode. Many privileged operations cause a VM exit with a reason code. The hypervisor does not emulate every instruction in software.

The hypervisor sets a control structure (VMCS on Intel, VMCB on AMD). That structure holds guest state, host state, and which events cause exits.

Benefits:

- The guest kernel can use its normal privileged instructions when the control bits allow it.
- Exit reasons are explicit (CPUID, I/O, EPT fault, interrupt window).
- Nested paging (next section) reduces memory virtualization cost.

`CPUID` in a guest often still exits so that the hypervisor can hide or expose features. A guest that sees `hypervisor` in CPU flags knows it is virtualized.

You enable these features in firmware (BIOS/UEFI) on a lab PC. `/proc/cpuinfo` on Linux shows `vmx` (Intel) or `svm` (AMD) when the host CPU and firmware allow KVM.

Hardware virt does not remove the need for a hypervisor. It makes the common path faster. Device emulation and scheduling remain software.

Do not confuse VT-x with VT-d. VT-d is Intel's IOMMU brand for device assignment. Topic 14 mentioned the IOMMU.

### Questions

#### Theoretical questions

1. What do VT-x and AMD-V add?
2. What is a VMCS (or VMCB) for?
3. What is a VM exit reason?
4. Why can `CPUID` still exit?
5. How does `vmx` or `svm` in `/proc/cpuinfo` help a KVM host?

#### Easy practical tasks

1. Run `grep -E 'vmx|svm' /proc/cpuinfo | head`. Write whether the flag exists.
2. Write five sentences about hardware virt. Use only facts from this section.
3. Make a table: VT-x/AMD-V versus VT-d/IOMMU. Add one job each.
4. Open `man 1 lscpu` and write whether a hypervisor flag is listed.

#### Medium practical tasks

1. Draw guest mode, VM exit, host hypervisor, VM resume.
2. Read a short Intel VT-x or AMD-V overview (public). Write six sentences on guest mode versus host mode.
3. If `/dev/kvm` is missing, write three possible causes (firmware, CPU, module) from documentation. Do not force a host change.

#### Advanced practical tasks

1. Read KVM `Documentation/virt/kvm/api.rst` on `KVM_RUN` and exit reasons. Write a one-page list of five exit types in your own words.
2. Compare TCG (pure QEMU emulation) with KVM acceleration in a table: need for VT-x, speed, and use in class.

---

## Memory virtualization (EPT/NPT)

A guest OS has its own page tables. Those tables map guest virtual addresses to guest physical addresses (GPA). GPA is not a host physical address (HPA). The hypervisor must map GPA to HPA.

Without hardware support, the hypervisor used shadow page tables: it built host page tables that the CPU used, and it trapped guest page-table updates. That method caused many exits.

Extended Page Tables (EPT, Intel) and Nested Page Tables (NPT, AMD) add a second translation in hardware:

1. Guest virtual address to GPA (guest page tables).
2. GPA to HPA (EPT/NPT tables that the hypervisor owns).

The TLB and the MMU walk both levels. A miss can be expensive, but the guest can update its own page tables without a trap in the common case.

An EPT fault occurs when the GPA is not mapped or the permission is wrong. The hypervisor then allocates a frame, maps a device MMIO page, or injects a fault into the guest.

Huge pages and EPT huge pages reduce walk cost. Memory overcommit and ballooning are extra hypervisor tricks. This section only needs the two-level walk.

Containers do not use EPT. They use the host page tables of ordinary processes.

### Questions

#### Theoretical questions

1. What is a guest physical address?
2. What problem do EPT and NPT solve?
3. What is a shadow page table?
4. What happens on an EPT fault?
5. Why do containers not need EPT?

#### Easy practical tasks

1. Write five sentences about EPT/NPT. Use only facts from this section.
2. Draw a two-level walk: GVA to GPA to HPA.
3. Make a table: guest PTE versus EPT entry. Add who writes each table.
4. Look up `ept` in `lscpu` flags or kernel docs. Write one sentence.

#### Medium practical tasks

1. Compare a host process page fault with an EPT fault in six sentences.
2. Read a short KVM MMU note. Write how a guest `mmap` becomes GPA mappings and then EPT mappings (idea).
3. Explain why DMA from an assigned device needs an IOMMU when GPA is not HPA (topic 14).

#### Advanced practical tasks

1. Read about two-dimensional page walks and TLB cost. Write a one-page performance note for a guest with a large working set.
2. Write a contrast: Linux `mmap` in a container versus `mmap` in a KVM guest. Name which translation layers exist.

---

## Containers vs VMs

A container is a set of isolated processes that share the host kernel. Linux implements this isolation with namespaces and cgroups (next section). A container image is a root filesystem plus metadata. The runtime (`runc`, Docker, Podman) starts the process tree.

A VM has a guest kernel. The hypervisor isolates that kernel from the host and from other guests.

Comparison:

- Kernel: a container uses the host kernel. A VM has its own kernel.
- Isolation: a VM boundary is stronger in the usual threat model. A kernel bug can break all containers on that host.
- Size and start time: containers are smaller and start faster.
- Devices and modules: a container cannot load an arbitrary kernel module. A VM can (inside the guest).
- OS mix: a container does not run a Windows kernel on a Linux host. A VM can.

Use a container when you want a repeatable user-space environment and you trust the host kernel. Use a VM when you need another kernel, a crash lab, or a stronger boundary.

WSL 2 is a lightweight VM with a Linux kernel. A Linux Docker container on a Linux host is not a VM.

Do not say "the container has its own OS" unless you mean its own user-space root filesystem. The kernel is shared.

### Questions

#### Theoretical questions

1. What kernel does a Linux container use?
2. What kernel does a VM use?
3. Why can one kernel bug affect all containers on a host?
4. When do you pick a VM for an OS course lab?
5. How does WSL 2 differ from a Docker container on Linux?

#### Easy practical tasks

1. Run `cat /proc/1/cgroup` and `ls /proc/self/ns`. Write that you see namespace links.
2. Write five sentences that contrast containers and VMs. Use only facts from this section.
3. Make a table: start time, own kernel, isolation goal.
4. Open `man 1 docker` or `man 1 podman` if present. Write the first sentence of the description.

#### Medium practical tasks

1. Draw host kernel, two containers, and one VM with a guest kernel on the same machine.
2. Read a short runtime overview (`runc`). Write six sentences: namespaces plus cgroups plus a root filesystem.
3. Compare `uname -r` in the host and in a container (if you have one). Write whether the release matches.

#### Advanced practical tasks

1. Write a one-page decision guide for your lab: VM, container, WSL 2, or dual boot.
2. Read about seccomp and user namespaces as extra container layers (topic 17). Write how they change the threat model without becoming a VM.

---

## cgroups and namespaces (Linux)

Linux namespaces split the view of kernel objects. A process has a set of namespace memberships. `ls /proc/self/ns` shows inode links for each type.

Main namespace types:

- mount (`mnt`): mount table (topic 12)
- process ID (`pid`): PID 1 inside the namespace is not the host PID 1
- network (`net`): interfaces, routes, and firewall (topic 18)
- user (`user`): user IDs map to different host UIDs
- UTS: hostname
- IPC: System V IPC and POSIX mq isolation
- cgroup namespace: the view of the cgroup tree
- time namespace: some clocks (newer kernels)

`unshare` and `clone` create namespaces. `nsenter` joins them. A container runtime does this for you.

Control groups (cgroups) limit and account CPU, memory, I/O, and other resources. The kernel can kill a cgroup that exceeds a memory max (OOM in that cgroup). `systemd` and container runtimes create cgroup trees under `/sys/fs/cgroup`.

Namespaces hide names. Cgroups limit amounts. You need both for a container. You also need a filesystem view (pivot_root or overlay) and often seccomp (topic 17).

`unshare` as a normal user may fail for some types. User namespaces allow more without root on many systems. The exact policy depends on the distribution.

Do not treat namespaces as a complete security product. They are building blocks.

### Questions

#### Theoretical questions

1. What does a namespace isolate?
2. What does a cgroup limit?
3. What does a PID namespace change about PID 1?
4. Why does a container need both namespaces and cgroups?
5. What does a user namespace map?

#### Easy practical tasks

1. Run `ls -l /proc/self/ns`. Write the namespace file names.
2. Open `man 7 namespaces` and `man 7 cgroups`. Write one sentence for each.
3. Run `cat /proc/self/cgroup`. Write the first line.
4. Write five sentences about cgroups and namespaces. Use only facts from this section.

#### Medium practical tasks

1. Run `unshare --user --map-root-user --pid --fork --mount-proc ps -p 1` if your system allows it. Write what PID 1 looks like inside. If the command fails, write the error.
2. Draw a process with arrows to mnt, pid, and net namespaces.
3. Read `/sys/fs/cgroup` (cgroup v2). Write two controller names that you see.

#### Advanced practical tasks

1. Create a user and mount namespace, mount a tmpfs, and show that the host does not see it (VM preferred). Document every command.
2. Write a one-page map: Docker or Podman flags to namespaces and to cgroup limits (`--memory`, `--cpus`).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A classmate says "Docker is a type 2 hypervisor." Which facts do you use to correct that sentence?
2. How do VT-x, EPT, and virtio work together to keep a guest fast?
3. When does trap-and-emulate still happen on a modern KVM guest?
4. How do mount namespaces (topic 12) and this topic's container story fit one process start path?
5. Why can a VM on KVM still need QEMU in user space?

#### Easy practical tasks

1. Write a one-page cheat sheet: VM, guest, host, type 1/2, KVM, trap, VT-x, EPT, container, namespace, cgroup.
2. Run `ls /dev/kvm`, `ls /proc/self/ns`, `cat /proc/self/cgroup`, and `systemd-detect-virt` if present. Comment each.
3. Draw one figure: hardware, host kernel (KVM), QEMU, guest kernel, and a container next to it that skips the guest kernel.
4. Bookmark `man 7 namespaces`, `man 7 cgroups`, and KVM documentation.

#### Medium practical tasks

1. Write a short script that prints `uname -r`, namespace inode numbers, and whether `/dev/kvm` exists. Run it on the host and inside a container if you have one. Save both outputs.
2. Start a disposable VM (or use an existing lab VM). From the guest, record `lscpu` hypervisor bits. From the host, record `lsmod | grep kvm`.
3. Map topic 1 (privilege modes) and topic 8 (page tables) onto VM exits and EPT in a table.

#### Advanced practical tasks

1. Read the OSTEP or kernel documentation material on virtualization if you have it. Write a one-page map to each section in this handbook.
2. Build the smallest container-like environment you can with `unshare` and a directory as root (no Docker). Document what is still shared with the host (kernel, some devices).
