# 10. Virtualization and Containers

## Description

Virtualization lets one machine run more than one operating system, or lets one kernel isolate many user environments. This topic explains why virtual machines (VMs) exist, type-1 and type-2 hypervisors, and hardware virtualization. You learn the difference between containers and VMs, and Linux cgroups and namespaces.

Complete this topic after IPC. Complete this topic before security and isolation (topic 11). You already know privilege modes, page tables, and mounts.

Use one term for each concept. A hypervisor is the software that runs virtual machines. A guest is the OS inside a VM. A host is the system that provides the hypervisor or the shared kernel. A container is an isolated user-space environment that shares the host kernel. A namespace is a per-process view of one kernel resource. A cgroup is a kernel limiter and accountant for resource use. Do not mix a container with a VM. Do not mix type 1 and type 2 hypervisors.

---

## Why VMs exist; type-1 vs type-2 hypervisors

A virtual machine is a software illusion of a computer. The guest OS runs as if it owned a CPU, memory, and devices. The hypervisor multiplexes the real hardware among guests.

Reasons to use a VM:

1. Isolation: a crash or a compromised guest must not take down other guests (the goal; bugs exist).
2. Density: many servers share one physical machine.
3. Compatibility: an old OS or a different kernel can run on new hardware.
4. Snapshots and migration: you can copy or move a guest disk and sometimes live memory.
5. Teaching and test: you can break a guest without a second physical PC.

A VM is not the only isolation tool. A process already isolates address spaces. A VM isolates kernels. That extra layer costs CPU, memory, and I/O overhead. You pay that cost when you need a different kernel, a different OS, or a strong trust boundary.

A hypervisor (virtual machine monitor) is the software that creates and runs VMs.

A type-1 hypervisor runs on the hardware. It does not need a full general-purpose host OS under it. Examples in this class: Xen, ESXi, Hyper-V in its bare-metal role, and the idea of Linux KVM as a kernel module that makes the kernel a hypervisor. The line is a bit blurry on Linux because Linux is also a full OS.

A type-2 hypervisor runs as a program on a host OS. VirtualBox and VMware Workstation are typical examples. QEMU can emulate in user space. QEMU plus KVM uses the host kernel as the hypervisor and QEMU as the user-space helper for devices and the GUI.

Comparison:

- Type 1: smaller extra stack, common in data centers, often managed as the main OS of the box.
- Type 2: easy on a laptop, good for labs, one more layer (the host OS) in the trust and performance path.

Cloud vendors sell VMs as a default unit. The same ideas apply to a local QEMU/KVM guest.

A VM does not replace backups. A VM does not make an unsafe guest program safe inside the guest. The guest still needs its own user IDs and updates.

### Questions

#### Theoretical questions

1. What illusion does a VM give to a guest OS?
2. Name three reasons to run a VM.
3. What is a type-1 hypervisor?
4. What is a type-2 hypervisor?
5. Why is the type line blurry for Linux KVM?

#### Easy practical tasks

1. Write five sentences about why VMs exist and about hypervisor types. Use only facts from this section.
2. Run `systemd-detect-virt` or `cat /proc/cpuinfo | grep -i hypervisor` if present. Write whether you look like a guest.
3. Open a QEMU or VirtualBox overview page. Write one sentence about what the product virtualizes.
4. Make a table: type 1 versus type 2. Add "runs on" and "typical example".

#### Medium practical tasks

1. Draw a physical machine, a type-1 hypervisor, and two guests. Then draw a host OS, a type-2 hypervisor, and one guest.
2. Compare a cloud VM and a local lab VM in six sentences: who owns the hypervisor, and what you can snapshot.
3. Read a short note on live migration (idea only). Write what must move (memory, disk, network identity).

#### Advanced practical tasks

1. Read the first pages of a KVM intro (kernel docs or a textbook). Write a one-page list of resources that a guest thinks it owns.
2. Write a policy: when your course work must use a VM instead of WSL or a container (kernel modules, crash tests).

---

## Hardware virtualization

Old virtualization trapped privileged guest instructions and emulated them in software. That path is slow. Modern CPUs add hardware virtualization.

On x86 the extensions are Intel VT-x and AMD-V. The CPU has a guest mode and a host (root) mode. The guest runs many instructions at native speed. A privileged operation or a configured event causes a VM exit. The hypervisor handles the exit and then a VM entry returns to the guest.

Typical exit reasons: a privileged register access, an interrupt injection, an I/O port or MMIO access that the hypervisor emulates, or a hypercall that the guest OS issues on purpose.

Memory virtualization uses an extra translation layer. The guest thinks it has guest physical addresses. The host maps those to real frames. Intel EPT and AMD NPT (nested page tables) let the MMU walk that extra table in hardware. Without EPT/NPT, the hypervisor must shadow page tables in software. Shadows are more exits and more work.

I/O virtualization: the hypervisor can emulate a simple disk and NIC (QEMU devices). Virtio is a paravirtual driver pair: the guest knows it is a VM and uses a fast queue. SR-IOV and device assignment give a guest a real device slice. Those need IOMMU support.

You can check host flags on Linux: `grep -E 'vmx|svm' /proc/cpuinfo`. Nested virtualization is optional and slower.

Hardware virtualization does not remove the need for a hypervisor. It makes the common path faster. A malicious or buggy guest can still stress the host. Firmware and device bugs exist.

### Questions

#### Theoretical questions

1. What is a VM exit?
2. What do Intel VT-x and AMD-V add?
3. What problem do EPT and NPT solve?
4. What is virtio?
5. Why is trap-and-emulate without hardware support slow?

#### Easy practical tasks

1. Run `grep -E 'vmx|svm' /proc/cpuinfo | head`. Write whether the flags exist (or that you are in a guest without the flags).
2. Open a short KVM or VT-x overview. Write one sentence about VM exit.
3. Write five sentences about hardware virtualization. Use only facts from this section.
4. Make a table: guest linear address, guest physical, host physical. Add who maps each step.

#### Medium practical tasks

1. Draw VM entry, guest run, VM exit, hypervisor, VM entry.
2. Compare full emulation (QEMU TCG) with KVM acceleration in six sentences.
3. Read a virtio overview. Write why a paravirtual NIC beats a fully emulated NE2000.

#### Advanced practical tasks

1. Read kernel documentation on KVM. Write a one-page note: vCPU, memslot, and exit handling at a high level.
2. Write when you enable nested KVM and when you must not (lab laptop versus production).

---

## Containers vs VMs

A container is a set of processes that share the host kernel and that have an isolated view of resources. The isolation uses namespaces and cgroups (next section). The container has its own root filesystem view (often a tarball of user space). It does not boot its own kernel.

A VM has a guest kernel. The guest has its own scheduler, page tables, and drivers (virtual or assigned). The guest can be Linux while the host is Linux, or the guest can be another OS.

Compare:

- Kernel: container shares the host kernel. VM has a guest kernel.
- Boot: container starts processes. VM boots firmware and a kernel (topic 1).
- Density: you can pack more containers than VMs on the same RAM for many workloads.
- Isolation: a VM trust boundary is stronger in the usual model (a kernel bug in the host still matters for both). A container escape is a host-kernel bug or a bad config.
- Compatibility: a container cannot run a different kernel ABI. A VM can run another OS.

Docker, Podman, and systemd-nspawn are user tools that set up namespaces, cgroups, and a rootfs. They are not a second kernel. A "container image" is files plus metadata. The running object is processes.

A container on a non-Linux host often sits inside a VM (Docker Desktop). Then you have both layers. For kernel labs, a real Linux VM or a native Linux host is closer to the truth.

Do not treat a container as a VM in a threat model. Do not treat a VM as cheap as a process.

### Questions

#### Theoretical questions

1. What kernel does a container use?
2. What extra isolation does a VM add over a container?
3. Why can you run more containers than VMs on the same machine?
4. Why can a container not run Windows on a Linux host kernel?
5. What is a container image versus a running container?

#### Easy practical tasks

1. Write a two-column table: container versus VM. Add five rows.
2. Run `cat /proc/1/cgroup` and `ls /proc/self/ns` . Write one fact from each.
3. Write five sentences that contrast containers and VMs. Use only facts from this section.
4. Open a Docker or Podman overview. Write one sentence that names namespaces or cgroups.

#### Medium practical tasks

1. Draw host kernel, two containers, and separately host hypervisor plus two guest kernels.
2. Read `man 1 unshare`. Write which namespace flags exist.
3. If you have Podman or Docker, run `podman info` or `docker info` (read only). Write whether the engine uses a VM.

#### Advanced practical tasks

1. Write a one-page threat note: container escape versus guest escape. Name the kernel that an attacker needs to exploit.
2. Compare OCI runtime (runc) documentation with QEMU. Write who creates processes and who emulates devices.

---

## cgroups and namespaces

Linux namespaces split a global kernel resource into per-process views. A process sees only its namespace instance.

Common namespaces:

- PID: process IDs. The container `init` can be PID 1 inside the namespace.
- mount: mount table. The container has its own `/`.
- UTS: hostname.
- network: interfaces, routes, and filter tables (topic 12).
- IPC: System V IPC and POSIX message queues.
- user: UID and GID map. UID 0 inside can map to a non-zero UID on the host.
- time (newer): clock offsets.
- cgroup: the cgroup view.

`clone` and `unshare` create namespaces. `/proc/<pid>/ns` holds namespace inodes. Two processes with the same inode number share that namespace.

Control groups (cgroups) limit and account resources: CPU, memory, pids, I/O. A process sits in a cgroup. The kernel can deny a new `malloc` or kill on memory limit (OOM in that cgroup). CPU shares or max bandwidth limit a noisy neighbor.

cgroup v2 is the current unified hierarchy under `/sys/fs/cgroup`. Older v1 had many controllers as separate mounts. Read the v2 documents for new work.

Namespaces isolate names and views. cgroups isolate amounts. You need both for a container. A namespace without a memory limit still lets a process eat host RAM. A cgroup without a mount namespace still sees the host tree.

User namespaces plus a UID map are a rootless-container building block. They are not a full security product by themselves (topic 11).

Do not run `unshare -m` experiments on a shared host without a VM. You can hide mounts from yourself in confusing ways.

### Questions

#### Theoretical questions

1. What does a PID namespace change?
2. What does a user namespace change about UID 0?
3. What do cgroups limit?
4. Why does a container need both namespaces and cgroups?
5. What is cgroup v2?

#### Easy practical tasks

1. Run `ls /proc/self/ns`. Write the namespace file names.
2. Open `man 7 namespaces` and `man 7 cgroups`. Write one sentence for each.
3. Write five sentences about cgroups and namespaces. Use only facts from this section.
4. Run `cat /proc/self/cgroup`. Write the path that you see.

#### Medium practical tasks

1. In a VM, run `unshare --user --map-root-user --pid --fork --mount-proc /bin/bash` if your policy allows it. Run `id` and `echo $$`. Write UID and PID as seen inside. Exit.
2. Draw a process with a private mount namespace and a cgroup memory max.
3. Read `/sys/fs/cgroup` (v2) if present. Write one controller file name (`memory.max` or similar).

#### Advanced practical tasks

1. Write a one-page map: Docker/Podman flag to namespace and to cgroup controller (from public docs).
2. Read `man 7 user_namespaces`. Write when "root inside" is not host root, and what capabilities still matter (topic 11).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. When do you pay for a hypervisor, and when is a container enough?
2. How do hardware VM exits and guest page tables relate to topic 1 privilege modes and topic 5 translations?
3. A classmate says "a container is a lightweight VM." Which facts do you use to correct that sentence?
4. How do a network namespace and a PID namespace hide different things?
5. Why can a memory cgroup kill a process that `malloc` already returned?

#### Easy practical tasks

1. Write a cheat sheet: VM, guest, host, type 1, type 2, VT-x/SVM, EPT, virtio, container, namespace, cgroup.
2. Run `systemd-detect-virt`, `nproc`, and `ls /proc/self/ns`. Write one line each.
3. Draw two stacks: laptop type-2 lab, and cloud type-1 guests.
4. Bookmark `man 7 namespaces` and `man 7 cgroups`.

#### Medium practical tasks

1. Write a six-line decision: course lab that loads a kernel module versus a course lab that runs a second `nginx`. Pick VM or container.
2. Compare `/proc/1/sched` or PID 1 name on the host with PID 1 inside a container (if you have one). Write the difference.
3. Document a safe experiment rule: use a VM for `unshare` and cgroup writes. Do not test on a shared login host.

#### Advanced practical tasks

1. Read a KVM and a runc overview. Write a one-page table: who owns the page tables that user processes see.
2. Design a tiny "container" on paper: `unshare` flags, a chroot or mount of a rootfs, and a memory.max. Do not claim it is production.
