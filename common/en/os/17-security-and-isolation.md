# 17. Security and Isolation

## Description

The kernel isolates users and processes. This topic explains user IDs and least privilege, Linux capabilities, and seccomp. You learn the idea of a sandbox, the difference between a kernel exploit and a user-space exploit, and Secure Boot at a high level.

Complete this topic after virtualization. Complete this topic before the networking stack (topic 18). You already know permission bits, namespaces, and kernel versus user space. You now learn how Linux limits what a process may do.

Use one term for each concept. A user ID (UID) is the integer that the kernel stores for a process and for file ownership. Least privilege is the rule that a program must run with the fewest rights that still let it work. A capability is a split of old root power into named rights. Seccomp is a filter on system calls. A sandbox is a combination of isolation tools around a program. Do not mix a capability with a file permission bit. Do not mix Secure Boot with a login password.

---

## User IDs and least privilege

Linux identifies a user by a UID. A process has a real UID, an effective UID, and a saved UID (and filesystem UID). The kernel uses the effective UID for most permission checks. `id` prints the IDs of the shell.

UID 0 is root. Root may pass most permission checks. A program that runs as root can read private files, bind low ports, and load modules (unless other layers stop it).

Least privilege means:

1. Use a normal user account for daily work (topic 1).
2. Give a service its own UID.
3. Do not run a network daemon as root after bind if you can drop privileges.
4. Use `sudo` for one command, not a root shell for the whole session.

`setuid` binaries run with the file owner's effective UID. `passwd` is a classic example. A bug in a `setuid` program is a privilege-escalation risk. Distributions minimize `setuid` root binaries.

A user namespace can map UID 0 inside the namespace to a non-zero host UID. Inside, the process looks like root. On the host, it is not root. This map is a container building block (topic 16). It is not a full security product by itself.

File permission bits (topic 11) still apply. Isolation fails if every file is world-readable or if your service UID can write `/`.

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

1. Write a small C program that prints `getuid` and `geteuid`. Run it normally.
2. Draw real UID, effective UID, and a file check on `open`.
3. Read `man 8 sudo`. Write six sentences: how `sudo` changes identity for one command.

#### Advanced practical tasks

1. Read `man 7 user_namespaces`. Write a one-page note: when "root in a user namespace" is not host root.
2. Design a service account: UID, home, and which directories it may write. No need to create it on a shared host.

---

## Capabilities (Linux)

Old Unix treated root as one bit: UID 0 or not. Linux capabilities split that power into named flags. Examples:

- `CAP_NET_BIND_SERVICE`: bind ports below 1024
- `CAP_NET_ADMIN`: network configuration
- `CAP_SYS_ADMIN`: a large set of admin operations (still very wide)
- `CAP_SYS_MODULE`: load kernel modules
- `CAP_DAC_OVERRIDE`: bypass file permission checks

A process has a capability set. File capabilities can grant bits when you `exec` a binary (`setcap`). This is a finer tool than `setuid` root, but a wide capability is still almost root.

`capsh`, `getpcaps`, and `/proc/self/status` (Cap* fields) show sets. `man 7 capabilities` is the reference.

Ambient and inheritable sets control what children keep. Container runtimes drop most capabilities. systemd unit files can use `CapabilityBoundingSet`.

`CAP_SYS_ADMIN` is a warning sign. It covers many operations. Dropping "everything except `CAP_SYS_ADMIN`" is not least privilege.

Capabilities do not replace UIDs. A process can be UID 1000 and still have `CAP_NET_BIND_SERVICE`. A process can be UID 0 and have an empty bounding set in a tight sandbox.

Do not `setcap` random binaries on a machine that you need.

### Questions

#### Theoretical questions

1. What problem do capabilities solve compared with a single root bit?
2. What does `CAP_NET_BIND_SERVICE` allow?
3. Why is `CAP_SYS_ADMIN` a warning?
4. How can a file grant a capability at `exec`?
5. Why can a non-zero UID still have a dangerous capability?

#### Easy practical tasks

1. Open `man 7 capabilities`. Write three capability names and one sentence each.
2. Run `grep Cap /proc/self/status`. Write the field names.
3. Run `command -v capsh` and `command -v getpcaps`. Write which tools exist.
4. Write five sentences about capabilities. Use only facts from this section.

#### Medium practical tasks

1. Decode `CapEff` with `capsh --decode=` if `capsh` exists. Write two bits that are on or off.
2. Read a systemd unit that sets `CapabilityBoundingSet` (example on the web or on your system). Write which bits it keeps.
3. Make a table: `setuid` root versus `setcap cap_net_bind_service+ep`. Add residual risk.

#### Advanced practical tasks

1. In a VM, start a process with `capsh --drop=...` or a systemd unit and show that `ping` or bind to port 80 fails without the needed cap. Document commands.
2. Read the capability handbook section on `CAP_SYS_ADMIN`. Write a one-page list of operation classes it still covers.

---

## Seccomp

Seccomp (secure computing) is a Linux filter on system calls. After a process installs a filter, the kernel checks each system call against the filter. The filter can allow, deny, return an error, or kill the thread.

The original mode allowed only a tiny set of calls (`read`, `write`, `exit`, and a few others). That mode is too strict for a normal libc program.

Seccomp-BPF lets you write a Berkeley Packet Filter program on the syscall number and argument values (with limits). libseccomp is the usual helper library. Container runtimes install a default profile. Chrome and other sandboxes use seccomp.

A filter is a blacklist or a whitelist. A whitelist of needed calls is tighter. A blacklist of famous bad calls is weaker. New syscalls can appear. A deny-by-default list ages better if you maintain it.

Seccomp does not filter file paths by itself. A process that may call `open` can still open many files. Combine seccomp with chroot, namespaces, and file permissions.

`prctl(PR_SET_SECCOMP, ...)` and `seccomp()` install filters. Once installed, a typical process cannot relax the filter. A privileged helper can be a second process.

Do not copy a random BPF filter from the internet onto a production host. Test in a VM. A wrong filter breaks `exec` or DNS and looks like a mysterious `EPERM`.

### Questions

#### Theoretical questions

1. What does seccomp filter?
2. What is the difference between classic seccomp and seccomp-BPF?
3. Why is a syscall whitelist tighter than a blacklist?
4. Why is seccomp not enough if `open` is allowed?
5. What happens when a filter denies a call?

#### Easy practical tasks

1. Open `man 2 seccomp` and `man 2 prctl`. Write one sentence for each.
2. Search `man 7 seccomp` if it exists. Write the purpose in one sentence.
3. Write five sentences about seccomp. Use only facts from this section.
4. Make a table: allow, error, kill. Add one use each.

#### Medium practical tasks

1. Read a Docker or systemd `SystemCallFilter` example. Write six sentences: which calls they block by name.
2. Draw libc `open`, syscall number, seccomp, then the real `open` path.
3. Run `grep Seccomp /proc/self/status`. Write the number and look up the meaning in `man 5 proc`.

#### Advanced practical tasks

1. Write a tiny C program that installs a libseccomp filter that allows `write`, `exit`, and `exit_group` only, then prints one line. Document every failure with libc.
2. Read about seccomp and `io_uring` or `ptrace`. Write a one-page note: new syscalls can bypass an old blacklist. No exploit steps.

---

## Sandboxing idea

A sandbox is a bundle of restrictions around a program. The program still runs. It sees fewer objects and fewer calls.

Typical Linux ingredients:

1. UID and GID (least privilege)
2. File permissions and a dedicated root or overlay
3. Namespaces (topic 16)
4. Cgroups (topic 16)
5. Capabilities
6. Seccomp
7. Landlock or LSM policies (SELinux, AppArmor) on some systems
8. A VM when the kernel itself is untrusted (topic 16)

The idea is defense in depth. One layer has bugs. A second layer must still stop a class of harm.

A sandbox is not a proof of safety. A kernel bug can escape a container. A VM escape is rarer and still exists. An application bug can still leak data that the sandbox allows the process to read.

Browsers sandbox renderer processes. The renderer parses hostile HTML. The browser assumes that the renderer is a target of attack. The OS tools limit what a compromised renderer can open.

For this path, a sandbox is a design: list what the process needs, then deny the rest. Start with a normal UID. Add layers when the threat grows.

Do not build a sandbox only from `chroot`. Topic 12 already stated that `chroot` is weak.

### Questions

#### Theoretical questions

1. What is a sandbox in this handbook?
2. Why does a sandbox use more than one tool?
3. Why is a sandbox not a proof of safety?
4. Why do browsers sandbox renderer processes?
5. Why is `chroot` alone a weak sandbox?

#### Easy practical tasks

1. Write five sentences about the sandbox idea. Use only facts from this section.
2. Make an eight-row list of ingredients from this section. Tick which ones you already used in earlier topics.
3. Open `man 7 landlock` if it exists, or write that the page is absent.
4. Draw a process inside three rings: UID, seccomp, namespace.

#### Medium practical tasks

1. Pick `cat /etc/hostname`. Write a sandbox recipe that still allows that job and that blocks `socket`.
2. Read an AppArmor or systemd sandbox example (`ProtectSystem`, `PrivateTmp`). Write six sentences.
3. Compare a container and a VM as sandboxes in a table (topic 16 terms).

#### Advanced practical tasks

1. Write a one-page sandbox plan for a student HTTP server: UID, bind port, filesystem view, seccomp outline, cgroup memory max.
2. Read the Landlock or seccomp documentation overview. Write how a path-based LSM differs from a syscall filter.

---

## Kernel exploits vs user exploits

An exploit is a use of a bug to break a security rule. This section names classes. It does not teach how to write an exploit.

A user-space exploit attacks a user program. Examples: a buffer overflow in a setuid helper, a path traversal in a web server, a malicious document that crashes a parser. The attacker gains the rights of that process (and maybe more if the process is privileged). ASLR, NX, and W^X (topic 10) raise the cost. They do not remove bugs.

A kernel exploit attacks the kernel. A successful kernel exploit can give ring 0 on the host. Then user IDs, seccomp, and containers on that kernel are not a boundary. This is why you update the kernel and why you do not load random modules (topic 20).

A guest kernel exploit in a VM takes the guest. A hypervisor or host kernel exploit is worse. Isolation goals differ (topic 16).

Privilege escalation is any jump from fewer rights to more rights. It can stay in user space (UID 1000 to UID 0) or enter the kernel.

Defense in this course:

1. Run less code as root.
2. Keep the kernel and libc updated.
3. Do not disable mitigations for a small benchmark on a machine that you need.
4. Treat driver and module sources as trusted code.

This handbook never asks you to write an exploit, a proof of concept attack, or a bypass.

### Questions

#### Theoretical questions

1. What is a user-space exploit in this section?
2. What can a kernel exploit break that a user-space exploit does not break?
3. How do ASLR and NX help, and what do they not do?
4. Why does a kernel exploit collapse container isolation on that host?
5. What is privilege escalation?

#### Easy practical tasks

1. Write five sentences that contrast the two exploit classes. Use only facts from this section.
2. Run `uname -r`. Write that kernel updates matter for kernel bugs.
3. Make a table: user exploit versus kernel exploit. Add "typical gain".
4. Open `man 1 uname` and write where you read the release.

#### Medium practical tasks

1. Draw a compromised renderer (user) versus a compromised kernel. Show which processes die or leak.
2. Read a public CVE description for a fixed Linux bug (high level). Write the component (user program or kernel) and the fixed version. No exploit details.
3. List three mitigations from topic 10 and one isolation tool from this topic. Write how they stack.

#### Advanced practical tasks

1. Write a one-page incident note template: what to collect (`uname`, package versions, whether the process was root) without attack steps.
2. Read the kernel self-protection documentation overview. Write six sentences on why the kernel must not trust user pointers.

---

## Secure boot (high-level)

Secure Boot is a firmware feature. At power-on, firmware checks a cryptographic signature on the bootloader (and the chain can continue to the kernel). If the signature is not trusted, firmware can refuse to boot.

The trust anchor is keys in firmware (platform key, key exchange key, and allowed or forbidden lists in the UEFI model). A vendor or the user enrolls keys. Topic 3 named UEFI. This section names the check.

Goals:

- Stop an unsigned bootloader from running before the OS.
- Make a class of persistent bootkits harder.

Limits:

- Secure Boot does not encrypt your files. Disk encryption is a separate feature (LUKS on Linux).
- Secure Boot does not stop a signed but buggy bootloader or kernel.
- If an attacker has firmware control or your enrolled key, the chain does not help.
- Some hardware and some drivers need a signed shim. Distributions document this.

Linux often uses a signed shim and GRUB or a signed UKI (unified kernel image). You can see EFI variables under `/sys/firmware/efi` when the system booted in UEFI mode.

Secure Boot is not a login. It is not a sandbox. It is a boot-time authenticity check.

Do not turn Secure Boot off on a machine that you need unless you know why (a lab of unsigned kernels). Do not enroll random keys.

### Questions

#### Theoretical questions

1. What does Secure Boot check at a high level?
2. Where does the trust anchor live?
3. What does Secure Boot not encrypt?
4. Why can a signed buggy kernel still be a problem?
5. How does Secure Boot differ from a user password?

#### Easy practical tasks

1. Check whether `/sys/firmware/efi` exists. Write UEFI or not (topic 3).
2. Open a distribution Secure Boot page (Debian or Ubuntu). Write one sentence about shim.
3. Write five sentences about Secure Boot. Use only facts from this section.
4. Make a table: Secure Boot, LUKS, login password. Add the question each answers.

#### Medium practical tasks

1. Run `mokutil --sb-state` if the tool exists. Write the state. If it is missing, write that.
2. Draw firmware, signed bootloader, kernel, `init`. Mark the signature check.
3. Read a high-level UEFI Secure Boot overview. Write six sentences on keys. No key-extraction work.

#### Advanced practical tasks

1. Write a one-page lab policy: when a student VM may disable Secure Boot (unsigned test kernels) and what must stay on a laptop.
2. Compare Secure Boot with measured boot or TPM-backed attestation at a high level (names and goals only).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do UID, capabilities, seccomp, and a namespace each answer a different question: who, which admin right, which call, which objects?
2. A service binds port 80 and then handles hostile HTTP. Which privileges do you drop after bind, and what do you still need?
3. Why does Secure Boot not replace least privilege after the kernel starts?
4. How does a kernel exploit change the value of a perfect seccomp filter?
5. When do you add a VM to a sandbox instead of another Linux LSM?

#### Easy practical tasks

1. Write a one-page cheat sheet: UID/EUID, least privilege, `setuid`, capabilities, `CAP_SYS_ADMIN`, seccomp, sandbox layers, kernel vs user exploit, Secure Boot.
2. Run `id`, `grep Cap /proc/self/status`, `grep Seccomp /proc/self/status`, and `ls /sys/firmware/efi`. Comment each.
3. Draw one figure: boot signature check, then a user process with dropped caps and a seccomp box.
4. Bookmark `man 7 capabilities`, `man 7 credentials`, and `man 2 seccomp`.

#### Medium practical tasks

1. Write a systemd-style sandbox sketch for a tiny static file server (even if you only write the keys on paper): user, `PrivateTmp`, syscall filter idea, capability set.
2. Use `ps -o pid,user,group,cmd` on PID 1 and on your shell. Write why PID 1 is a different identity story.
3. Map topics 10, 11, 16, and 17 in a table: protection bit, file mode, namespace, capability.

#### Advanced practical tasks

1. Read the Linux man-pages for capabilities and seccomp plus a distribution hardening guide. Write a one-page map to this handbook.
2. In a disposable VM, create a user, run a server as that user, and list what `/proc/self/status` shows for Cap and Seccomp. Do not weaken the host that you need.
