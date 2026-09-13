# 19. Persistence of OS State

## Description

The kernel and the user-space services keep state that must survive a reboot or that must explain a crash. This topic explains system logs (`journald` and the Windows Event Log), crash dumps, hibernation and suspend, and time. You learn the difference between a wall clock and a monotonic clock, and you learn NTP at a high level.

Complete this topic after the networking stack. Complete this topic before internals and tracing (topic 20). You already know filesystems, `fsync`, and `init`. You now learn where the OS writes its own history.

Use one term for each concept. A system log is a durable or semi-durable record of events from the kernel and from services. A crash dump is a copy of memory (and CPU state) after a panic. Suspend keeps RAM powered and stops the machine. Hibernation writes RAM to disk and powers off. Wall-clock time is civil time (seconds since an epoch, adjustable). Monotonic time only moves forward for measuring intervals. Do not mix `journald` with a user application log file. Do not mix wall-clock time with a monotonic timer.

---

## Logging (`journald`, Event Log)

Linux systems that use systemd store many logs in the journal. The daemon is `journald`. The tool is `journalctl`. Units (services) and the kernel send structured records. The journal can live in memory only, or on disk under `/var/log/journal`.

Useful views:

- `journalctl -b` : this boot
- `journalctl -k` : kernel messages
- `journalctl -u ssh` : one unit
- `journalctl -p err` : priority

The journal is not the only log. Some programs still write text files under `/var/log`. `rsyslog` or `syslog-ng` can read the journal or `/dev/log` and write files. The kernel ring buffer is also in `dmesg`. `journalctl -k` and `dmesg` overlap.

Windows uses the Event Log. Applications and the system write events to named logs (Application, System). Event Viewer reads them. The idea matches: a central store, severity, a timestamp, and a source. The API and the files differ.

Logs need rotation and disk limits. A full disk can block services. `journald` has size caps. `logrotate` manages text files.

Logs can contain secrets (passwords pasted into a command line, tokens). Do not publish a raw `journalctl` dump.

Least privilege applies. A normal user sees a subset. `systemd-journal` group or root sees more.

This section is persistence of OS events. Application logs (a web server `access.log`) are a separate design. Still write them with rotation.

### Questions

#### Theoretical questions

1. What does `journald` store?
2. What is the difference between `journalctl -k` and a user unit log?
3. How does the Windows Event Log compare at a high level?
4. Why do logs need a size limit?
5. Why can a journal dump contain secrets?

#### Easy practical tasks

1. Open `man 1 journalctl` and `man 8 systemd-journald`. Write one sentence for each.
2. Run `journalctl -b -n 20 --no-pager` if you have rights. Write two facility or unit names that you see.
3. Run `dmesg | tail` if allowed. Write one kernel line in your own words.
4. Write five sentences about system logging. Use only facts from this section.

#### Medium practical tasks

1. Compare `ls /var/log` with `journalctl --disk-usage` if the command exists. Write where bytes live.
2. Draw kernel, `journald`, a service, and a text file under `/var/log`.
3. Read `man 5 journald.conf` for `Storage=` and `SystemMaxUse=`. Write six sentences.

#### Advanced practical tasks

1. In a VM, set a small journal size, generate messages, and show that old boots disappear. Document the config. Do not do this on a shared host.
2. Write a one-page policy: what a student may paste from `journalctl` into a homework report.

---

## Crash dumps

A crash dump is a snapshot that you collect when the kernel or a process dies in a way that you must debug later.

A kernel panic can trigger `kdump` on Linux. A reserved crash kernel boots and writes the old memory to disk or to the network. The file is a `vmcore`. Tools such as `crash` read it. This path is an awareness topic. You need extra packages and a matching debug kernel to analyze a dump.

A user process can leave a core dump. The kernel writes a file (`core`) when a process receives some signals (`SIGSEGV`, `SIGABRT`) if `ulimit -c` is not zero. `systemd` may store cores through `coredumpctl`. The dump contains memory. Treat it as secret.

`sysctl kernel.core_pattern` chooses the name or a pipe to a helper.

A crash dump is not a backup of user documents. It is a forensic and debugging artifact. It is large. It can contain keys and passwords that were in RAM.

Windows has minidumps and full memory dumps after a bugcheck (stop error). The idea matches: save state, then reboot.

Do not enable `kdump` on a laptop with almost no disk unless you plan the space. Do not upload a core file to a public site.

For this course, `dmesg` after a VM panic and a user `coredumpctl list` are enough practice.

### Questions

#### Theoretical questions

1. What is a kernel crash dump for?
2. What is a user core dump?
3. Why is a core file sensitive?
4. What does `kernel.core_pattern` select?
5. How does a Windows minidump compare at a high level?

#### Easy practical tasks

1. Run `ulimit -c`. Write the value.
2. Open `man 5 core` and `man 1 coredumpctl` if present. Write one sentence for each.
3. Run `sysctl kernel.core_pattern` if allowed. Write the pattern.
4. Write five sentences about crash dumps. Use only facts from this section.

#### Medium practical tasks

1. Run `coredumpctl list` if systemd is present. Write whether dumps exist. Do not share them.
2. Draw panic, crash kernel, `vmcore` on disk, reboot into the normal kernel.
3. Read `man 8 kdump` or a distribution kdump page. Write six sentences on reserved memory.

#### Advanced practical tasks

1. In a disposable VM, allow core dumps, run a program that raises `SIGSEGV`, and inspect `coredumpctl info`. Then delete the dump.
2. Write a one-page note: when to collect `kdump` versus when `journalctl -k` after a reboot is enough.

---

## Hibernation / suspend

Suspend to RAM (S3, "sleep") keeps memory powered. The CPU stops. Devices enter low power. A wake event (lid, keyboard, WoL) resumes. Resume is fast. A power loss during suspend loses RAM and is equal to a crash.

Hibernation (suspend to disk) writes the contents of RAM to a swap device or a file, then powers off. On the next boot, the bootloader and kernel restore that image. Resume is slower. The machine can lose AC power.

Hybrid sleep writes a hibernation image and then suspends to RAM. If power remains, resume is fast. If power is lost, the disk image is used.

Linux uses `/sys/power/state` and systemd `systemctl suspend` or `hibernate`. Swap must be large enough for hibernation. Secure Boot and disk encryption add constraints. A bad resume can corrupt mounts if devices appear in a different order. Distributions document supported setups.

Windows has Sleep, Hibernate, and Fast Startup (a logoff plus a partial hibernation of the kernel session). The names differ. The RAM-versus-disk idea matches.

Do not hibernate a VM guest as your first lab unless the hypervisor documents it. Host sleep while a USB disk is mounted can also surprise you.

Suspend is not a substitute for `fsync`. Applications must still commit data. The kernel syncs more on hibernation, but you must not rely on a surprise power loss during suspend.

### Questions

#### Theoretical questions

1. What does suspend to RAM keep powered?
2. Where does hibernation write RAM?
3. What does a power loss during suspend do?
4. Why must swap be large for hibernation?
5. How does Windows Fast Startup relate to hibernation at a high level?

#### Easy practical tasks

1. Open `man 8 systemd-sleep` or `man 1 systemctl` for `suspend`. Write one sentence.
2. Run `cat /sys/power/state` if the file exists. Write the words that you see.
3. Write five sentences about suspend and hibernation. Use only facts from this section.
4. Make a table: suspend, hibernate, power off. Add "RAM contents" and "power use".

#### Medium practical tasks

1. Run `swapon --show` and `free -h`. Write whether swap could hold RAM for hibernation (rough compare).
2. Draw the hibernation path: RAM, swap, power off, boot, restore.
3. Read a distribution hibernate page. Write six sentences on encryption and resume.

#### Advanced practical tasks

1. In a laptop or VM that you own, run one suspend and resume. Write five devices that still worked. Do not force this on a shared lab PC.
2. Write a one-page warning: hibernate plus a dual-boot filesystem that is still mounted is unsafe.

---

## Time: wall clock vs monotonic

Programs ask the OS for time. Two families matter.

Wall-clock time (`CLOCK_REALTIME` on Linux) is civil time. It is the clock that `date` prints. NTP and an administrator can step it. A leap second or a bad sync can move it backward or jump it forward. Use it for logs, certificates, and "what time was this file written".

Monotonic time (`CLOCK_MONOTONIC`) measures intervals. It does not go backward when the wall clock jumps. It is the right clock for timeouts, benchmarks, and "sleep 50 ms". It can stop or jump across suspend depending on the clock flavor. `CLOCK_BOOTTIME` includes suspend time on Linux. `CLOCK_MONOTONIC` typically does not include the sleep period.

`CLOCK_MONOTONIC_RAW` is a hardware-backed monotonic clock that NTP does not slew. Use it when you must avoid NTP adjustment on the interval clock.

APIs:

- `clock_gettime`
- `gettimeofday` (older, wall clock)
- `time` (seconds)
- Go `time.Now` versus `time.Since` (the library maps to these ideas)

File timestamps (`mtime`) are wall-clock values. A wrong wall clock makes `make` and logs lie.

Do not measure elapsed time with `CLOCK_REALTIME`. A step of one hour breaks your timeout.

Topic 2 mentioned hardware timers. This section is the user-visible clocks that those timers feed.

### Questions

#### Theoretical questions

1. What is wall-clock time for?
2. What is a monotonic clock for?
3. Why can `CLOCK_REALTIME` move backward?
4. How does `CLOCK_BOOTTIME` differ from `CLOCK_MONOTONIC` on Linux?
5. Why is `mtime` a wall-clock value?

#### Easy practical tasks

1. Open `man 2 clock_gettime` and `man 7 time`. Write one sentence for each.
2. Run `date -u` and `cat /proc/uptime`. Write that the two numbers are different kinds of time.
3. Write five sentences about the two clock families. Use only facts from this section.
4. Make a table: `CLOCK_REALTIME`, `CLOCK_MONOTONIC`, `CLOCK_BOOTTIME`. Add one use each.

#### Medium practical tasks

1. Write a C program that prints both `CLOCK_REALTIME` and `CLOCK_MONOTONIC` twice with a `sleep(1)`. Write which difference is about one second.
2. Draw a wall-clock step during NTP and a monotonic line that does not step.
3. Read `man 2 clock_nanosleep`. Write six sentences on a relative monotonic sleep.

#### Advanced practical tasks

1. Compare `CLOCK_MONOTONIC` across a suspend on a machine that you own. Write whether the monotonic delta included the sleep.
2. Write a one-page API guide for a server: log timestamps, request deadlines, and certificate checks.

---

## NTP (high-level)

Network Time Protocol (NTP) and related protocols (chrony, `systemd-timesyncd`, SNTP) sync the wall clock to time servers. The kernel clock drifts. Without sync, logs on two machines disagree, and TLS certificates look not yet valid or expired.

A client queries servers, estimates delay, and slews or steps `CLOCK_REALTIME`. A slew changes the frequency slowly. A step jumps the clock. A large step can break running programs that used the wall clock as a timer.

Linux tools:

- `timedatectl` (systemd)
- `chronyc` (chrony)
- `ntpq` (ntpd)

`/etc/adjtime` and RTC (hardware clock) store time across power-off. UTC versus local time in the RTC is a firmware and OS agreement. Wrong RTC mode causes a one-hour or several-hour error at boot.

NTP is not authentication of users. Authenticated NTP and NTS exist for a tighter trust story. This section only needs: the wall clock is a distributed agreement, not a pure hardware fact.

Do not run two NTP daemons that fight. Do not step a production clock by hours without a window.

Full NTP packet detail belongs to `net.topics.md` (time services). Here you need the OS effect on `CLOCK_REALTIME`.

### Questions

#### Theoretical questions

1. What problem does NTP solve for the OS?
2. What is the difference between a slew and a step?
3. Why can a step break a program?
4. What is the RTC for?
5. Why is NTP not a user-login protocol?

#### Easy practical tasks

1. Run `timedatectl` if present. Write whether NTP is active.
2. Open `man 1 timedatectl` or `man 1 chronyc`. Write one sentence.
3. Write five sentences about NTP. Use only facts from this section.
4. Make a table: kernel drift, NTP slew, NTP step.

#### Medium practical tasks

1. Run `timedatectl show` or `chronyc tracking` if present. Write the offset field in your own words.
2. Draw RTC, kernel clock, NTP server, and `CLOCK_REALTIME`.
3. Read a short chrony or timesyncd page. Write six sentences on why laptops use a client, not a public stratum-1 server farm that you host.

#### Advanced practical tasks

1. Read about NTS or NTP authentication at a high level. Write a one-page trust note: who you believe for time.
2. Write a lab policy: how to set time in an offline VM (`timedatectl set-time`) and why that VM must not run a second NTP daemon.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A machine panics, reboots, and then has the wrong hour in the logs. Which stores do you inspect (journal, RTC, NTP, crash dump), and in what order?
2. Why is a monotonic timeout still useful during an NTP step that rewrites log timestamps?
3. How do hibernation and a crash dump both copy RAM, and how do their goals differ?
4. When is `journalctl -b -1` more useful than a core dump?
5. Why must you treat both a `vmcore` and a journal export as sensitive?

#### Easy practical tasks

1. Write a one-page cheat sheet: `journalctl`, Event Log, `dmesg`, `kdump`, core, suspend, hibernate, `CLOCK_REALTIME`, `CLOCK_MONOTONIC`, NTP.
2. Run `timedatectl` or `date`, `journalctl -b -n 5 --no-pager`, `ulimit -c`, and `cat /sys/power/state`. Comment each.
3. Draw one timeline: boot, NTP slew, a service error in the journal, a suspend, a resume.
4. Bookmark `man 1 journalctl`, `man 2 clock_gettime`, and `man 5 core`.

#### Medium practical tasks

1. Write a small program that logs a wall-clock timestamp and a monotonic start offset for each line. Run it. Save ten lines.
2. Document your system: journal on disk or not, NTP client name, swap size, core pattern.
3. Use `journalctl -o short-precise` and explain why precise stamps still need a correct wall clock.

#### Advanced practical tasks

1. Read OSTEP or kernel documentation on timekeeping if you have it. Write a one-page map to this handbook.
2. In a VM, disconnect NTP, set a wrong time, generate a log line, then fix NTP. Write how you would explain the skew in a report. Do not do this on a host that other people use for Kerberos or TLS.
