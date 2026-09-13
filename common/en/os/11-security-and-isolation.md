# 11. Security and Isolation

## Description

The kernel isolates users and processes. This topic explains user IDs and least privilege, Linux capabilities, and seccomp. You learn the idea of a sandbox, the difference between a kernel exploit and a user-space exploit, and Secure Boot at a high level.

Complete this topic after virtualization. Complete this topic before the networking stack (topic 12). You already know permission bits, namespaces, and kernel versus user space.

Use one term for each concept. A user ID (UID) is the integer that the kernel stores for a process and for file ownership. Least privilege is the rule that a program must run with the fewest rights that still let it work. A capability is a split of old root power into named rights. Seccomp is a filter on system calls. A sandbox is a combination of isolation tools around a program. Do not mix a capability with a file permission bit. Do not mix Secure Boot with a login password.

---

## UIDs and least privilege

Linux identifies a user by a UID. A process has a real UID, an effective UID, and a saved UID (and filesystem UID). The kernel uses the effective UID for most permission checks. `id` prints the IDs of the shell.

UID 0 is root. Root may pass most permission checks. A program that runs as root can read private files, bind low ports, and load modules (unless other layers stop it).

Least privilege means:

1. Use a normal user account for daily work.
2. Give a service its own UID.
3. Do not run a network daemon as root after bind if you can drop privileges.
4. Use `sudo` for one command, not a root shell for the whole session.

`setuid` binaries run with the file owner effective UID. `passwd` is a classic example. A bug in a `setuid` program is a privilege-escalation risk. Distributions minimize `setuid` root binaries.

A user namespace can map UID 0 inside the namespace to a non-zero host UID. Inside, the process looks like root. On the host, it is not root. This map is a container building block (topic 10). It is not a full security product by itself.

File permission bits (topic 6) still apply. Isolation fails if every file is world-readable or if your service UID can write `/`.

Do not practice as root. Do not disable user checks in a kernel that you need.

### Questions

#### Theoretical questions

1. What is an effective UID?
2. What is UID 0?
3. What does least privilege mean for a service?
4. Why is a `setuid` root binary a risk?
5. What does a user namespace change about UID 0?

#### Easy practical tasks

1. Run `id` and `id -u`. Write UID, GID, and groups.
2. Open `man 7 credentials` and `man 2 seteuid`. Write one sentence for each.
3. Run `find /usr/bin -perm -4000 2>/dev/null | head` (read-only). Write two `setuid` names if any appear.
4. Write five sentences about UIDs and least privilege. Use only facts from this section.

#### Medium practical tasks

1. Write a small C program that prints `getuid` and `geteuid`. Run it normally. This handbook does not contain the source.
2. Draw real UID, effective UID, and a file check on `open`.
3. Read `man 8 sudo`. Write six sentences: how `sudo` changes identity for one command.

#### Advanced practical tasks

1. Read `man 7 user_namespaces`. Write a one-page note: when "root in a user namespace" is not host root.
2. Design a service account: UID, home, and which directories it may write. Do not create it on a shared host.

---

## Capabilities and seccomp

Old Unix treated root as one bit: UID 0 or not. Linux capabilities split that power into named flags. Examples:

- `CAP_NET_BIND_SERVICE`: bind ports below 1024
- `CAP_NET_ADMIN`: network configuration
- `CAP_SYS_ADMIN`: a large set of admin operations (still very wide)
- `CAP_SYS_MODULE`: load kernel modules

A process has capability sets: permitted, effective, inheritable, and bounding (and ambient on modern Linux). File capabilities can grant a subset to a binary without full setuid root. `getcap` and `setcap` manage file caps. `capsh` and `/proc/<pid>/status` (`Cap*`) show process caps.

Drop capabilities that you do not need after startup. A web server may bind port 80, then drop `CAP_NET_BIND_SERVICE`. `CAP_SYS_ADMIN` is still close to root. Do not treat "I used capabilities" as "I am safe".

Seccomp (secure computing) is a Linux filter on system calls. The process installs a Berkeley Packet Filter (BPF) program. Each system call hits the filter. The filter can allow, deny (`ENOSYS` or `EPERM`), kill the thread, or trap.

`seccomp-bpf` is the usual form. libseccomp helps you write filters. Docker and systemd can apply seccomp profiles. A profile that only allows `read`, `write`, `exit`, and `sigreturn` is the historical strict mode. Real programs need a larger allow list.

Seccomp does not check file paths by itself. A process that may `open` can still open many files. Combine seccomp with UIDs, chroot or namespaces, and file permissions.

A bad seccomp profile breaks glibc (it needs unexpected calls). Test. Do not copy a random deny list from an old blog.

Do not confuse capabilities with file mode bits. `r-x` is not `CAP_NET_ADMIN`.

### Questions

#### Theoretical questions

1. What problem do capabilities solve compared with "all or nothing" root?
2. Why is `CAP_SYS_ADMIN` still dangerous?
3. What does seccomp filter?
4. Why is seccomp not enough if `open` is allowed?
5. What is a seccomp allow list?

#### Easy practical tasks

1. Open `man 7 capabilities` and `man 2 seccomp`. Write one sentence for each.
2. Run `grep Cap /proc/self/status`. Write that the fields exist (you do not need to decode hex yet).
3. Write five sentences about capabilities and seccomp. Use only facts from this section.
4. Run `getcap /usr/bin/ping` or another binary if `getcap` exists. Write the result or that none is set.

#### Medium practical tasks

1. Read `man 8 setcap`. Write how a file cap differs from setuid root.
2. Draw a process: drop caps after bind, then install a seccomp profile, then `accept`.
3. Read a Docker default seccomp profile overview (public docs). Write three syscalls that the profile blocks (example names from the docs).

#### Advanced practical tasks

1. Write a tiny program with libseccomp that allows only a few calls, then `open` a file and show the deny. Use a VM. This handbook does not contain the source.
2. Write a one-page note: bounding set versus effective set, and why a child after `exec` may lose caps.

---

## Sandboxing

A sandbox is a combination of isolation tools around a program that must not harm the rest of the system. No single tool is the sandbox. The product is the combination plus a policy.

Typical Linux pieces:

- a dedicated UID and GID
- a mount namespace and a minimal root filesystem (chroot, pivot_root, or an image)
- user, PID, and network namespaces (topic 10)
- cgroup limits (CPU, memory, pids)
- dropped capabilities
- a seccomp profile
- read-only mounts and `noexec` where you can
- Landlock or AppArmor or SELinux (MAC) on some distributions

Browsers and language runtimes use sandboxes for untrusted code (a renderer, a plugin). The trusted process is small. The untrusted process has almost no rights. IPC (topic 9) connects them. Unix sockets and `SCM_RIGHTS` are common.

A sandbox fails if:

- the kernel has a bug (next section)
- the policy allows a dangerous syscall or a writable mount
- you share `/` or `/proc` without care
- you run the helper as host root

`systemd-run` and `bwrap` (bubblewrap) are user tools that apply some of these pieces. Containers are a sandbox style. They are only as strong as the profile.

Do not download untrusted binaries and "sandbox" them on your main laptop without a VM. Do not treat a chroot as enough. chroot is not a security boundary for root.

### Questions

#### Theoretical questions

1. What is a sandbox in this topic?
2. Why is chroot not enough for a root process?
3. Why do browsers split a renderer into a tight sandbox?
4. Name five Linux tools that you can combine in a sandbox.
5. How can IPC keep a privileged helper small?

#### Easy practical tasks

1. Open `man 1 bwrap` if it exists, or a bubblewrap README. Write one sentence about what it wraps.
2. Write five sentences about sandboxes. Use only facts from this section.
3. Make a checklist table: UID, namespaces, cgroups, caps, seccomp. Add "what it limits".
4. Run `man 2 chroot`. Write one warning from the page.

#### Medium practical tasks

1. Draw trusted process, Unix socket, sandboxed child, and a denied `open` of `/home`.
2. Read systemd `ProtectSystem=` and `NoNewPrivileges=` documentation. Write six sentences.
3. Compare a container runtime sandbox with a browser renderer sandbox in a table: who is the attacker, and what is the prize.

#### Advanced practical tasks

1. In a VM, run a process under `bwrap` with a read-only root and no network (follow a current example). Write what still works (`echo`) and what fails.
2. Write a one-page sandbox policy for a homework grader that runs student binaries.

---

## Kernel exploits vs user exploits

A user exploit abuses a bug in a user-space program. Examples: a buffer overflow in a setuid helper, a path traversal in a web server, a command injection in a script. The attacker starts with the rights of that process. If the process is UID 0 or has wide capabilities, the attacker becomes root in user space. The kernel can still be intact.

A kernel exploit abuses a bug in the kernel or in a driver. Examples: a bad `ioctl`, a race in a syscall, a use-after-free in a network protocol. The attacker starts in user mode and reaches kernel mode. Then the attacker can change credentials, disable seccomp, or read any memory. Sandboxes that trust the kernel fail.

Defense in depth:

- fix and update user programs
- least privilege so that a user exploit is smaller
- sandbox so that a user exploit has few syscalls
- update the kernel and minimize loaded modules
- hardware features (NX, SMAP/SMEP, KASLR) raise the cost of some kernel attacks
- VMs add a hypervisor layer (topic 10). A guest kernel exploit is not automatically a host exploit, but hypervisor bugs exist

`/proc/kallsyms` and kernel pointers can help an attacker. Distributions restrict them. Do not treat "I cannot read kallsyms" as enough.

This path does not teach you to write exploits. You learn the split so that you pick the right update and the right isolation layer.

### Questions

#### Theoretical questions

1. What is a user exploit in this section?
2. What is a kernel exploit?
3. Why does a user exploit of a setuid binary matter?
4. Why does a sandbox fail after a kernel exploit?
5. How does a VM change the story of a guest kernel bug?

#### Easy practical tasks

1. Write a two-column table: user exploit and kernel exploit. Add "what you update" and "what isolation remains".
2. Write five sentences about this split. Use only facts from this section.
3. Run `uname -r`. Write that kernel updates matter.
4. Open `man 7 capabilities` and note that a cap is still kernel-enforced. Write one sentence.

#### Medium practical tasks

1. Draw attacker, vulnerable `setuid` program, UID 0, versus attacker, vulnerable driver, kernel mode.
2. Read a public CVE description of a Linux kernel bug (title and impact only). Write six sentences: impact, not a how-to.
3. List three user programs on your machine that listen on the network (`ss -lntp` if allowed). Write why a bug there is a user exploit first.

#### Advanced practical tasks

1. Write a one-page note: SMAP/SMEP and KASLR as high-level kernel self-protection. Use kernel documentation. No exploit steps.
2. Design an update policy: how fast you reboot for a kernel CVE versus restart a service for a user CVE.

---

## Secure boot (high-level)

Secure Boot is a firmware feature (UEFI). The firmware checks a cryptographic signature on the bootloader (or on a shim) before it runs that code. The bootloader can then check the kernel. The goal is: an attacker who writes a disk file must not replace the boot path with unsigned malware.

Keys sit in firmware (db, dbx, platform key). Distributions ship signed shims. You can enroll your own keys on a machine that you own. Enterprise machines can lock keys.

Secure Boot is not:

- a login password
- full disk encryption (that is LUKS or another disk key; often you use both)
- a guarantee against a later kernel exploit
- a guarantee if an attacker has firmware control or a physical DMA device

Measured Boot and TPM (Trusted Platform Module) can record hashes of the boot path. Remote attestation is a later topic. This section only names the idea: the TPM stores measurements. Secure Boot is the policy that refuses an unsigned loader.

If Secure Boot is off, the firmware still boots. The check is skipped. Dual-boot and some lab setups turn it off. Know the trade.

Topic 1 named UEFI. This section adds the signature check. Do not disable Secure Boot on a laptop that you use for banking unless you know why.

### Questions

#### Theoretical questions

1. What does Secure Boot check at a high level?
2. Where do the keys live?
3. Why is Secure Boot not full disk encryption?
4. What can still go wrong after a signed kernel starts?
5. What is the TPM idea in Measured Boot?

#### Easy practical tasks

1. Test whether `/sys/firmware/efi` exists. Write "UEFI" or "not UEFI".
2. Open a distribution Secure Boot page. Write one sentence about shim.
3. Write five sentences about Secure Boot. Use only facts from this section.
4. Make a table: Secure Boot, LUKS, login password. Add "what it protects".

#### Medium practical tasks

1. Run `mokutil --sb-state` if the tool exists. Write the state. Do not change keys.
2. Draw firmware, signature check, GRUB or shim, kernel, `init`.
3. Read a short note on dbx (revoked hashes). Write why a revoked bootloader must not run.

#### Advanced practical tasks

1. Write a one-page high-level map: UEFI Secure Boot, shim, kernel module signatures (`CONFIG_MODULE_SIG`). No key-extraction work.
2. Compare Secure Boot on a laptop with cloud UEFI guests (public cloud docs). Write who owns the keys.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do UID, capabilities, and seccomp form layers on one process?
2. Why is "run as root in a container" still a policy smell even with user namespaces?
3. Which sandbox pieces fail after a kernel exploit, and which VM layer might still hold?
4. How do topic 6 permission bits and topic 11 UIDs meet on one `open`?
5. What does Secure Boot prove, and what does it not prove, about a running system?

#### Easy practical tasks

1. Write a cheat sheet: UID, EUID, setuid, capability, seccomp, sandbox, user exploit, kernel exploit, Secure Boot, TPM (idea).
2. Run `id`, `grep Cap /proc/self/status`, and `ls /proc/self/ns`. Write one line each.
3. Draw a defense stack: firmware, kernel, service UID, seccomp.
4. Bookmark `man 7 capabilities`, `man 7 credentials`, and `man 2 seccomp`.

#### Medium practical tasks

1. Write a six-line hardening plan for a toy network daemon: bind, drop UID, drop caps, seccomp, chroot or namespace.
2. Use `sudo -l` if you have sudo. Write one command that your user may run. Do not escalate for fun.
3. Document a "do not" list: daily root shell, `chmod 777`, disable Secure Boot without a reason, empty seccomp allow of `open` only.

#### Advanced practical tasks

1. Read a systemd sandboxing man page (`systemd.exec`). Write a one-page map from directives to topics 10 and 11.
2. Design (on paper) a two-process helper: unprivileged worker, privileged opener, Unix socket, `SCM_RIGHTS`. State the threat you reduce.
