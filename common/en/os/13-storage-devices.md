# 13. Storage Devices

## Description

A filesystem sits on a storage device. This topic explains how those devices work at a level that an OS student needs. You learn historical hard-disk geometry, solid-state drives (SSDs) and wear, and the kernel split between block devices and character devices. You also learn the page cache, `fsync` and durability, and RAID at a high level.

Complete this topic after filesystem implementation. Complete this topic before I/O and device drivers (topic 14). You already know that a `write` can stay in RAM. You now learn when the device has the bytes.

Use one term for each concept. A block device transfers data in fixed-size blocks and supports random access. A character device transfers a byte stream and does not look like a disk. The page cache is the RAM cache of file pages. Durability means the data survives a power loss after the call that promised it. RAID is a layout of many disks that the OS or a controller presents as one volume. Do not mix the page cache with the device write cache. Do not mix `write` with `fsync`.

---

## HDD geometry (historical)

A hard disk drive (HDD) stores bits on spinning magnetic platters. A head flies over each surface. The OS and old firmware described the disk with cylinder, head, and sector (CHS).

A cylinder is the set of tracks that the heads can read without a seek, one track per surface. A seek moves the head arm to another cylinder. Rotational latency is the wait for the wanted sector to pass under the head. Transfer time is the time to read or write the sector bytes.

Old BIOS used CHS addresses. Modern disks use logical block addressing (LBA). LBA numbers blocks from 0. The firmware inside the disk maps LBA to physical locations. The OS almost never programs CHS on a current machine.

Geometry still explains old performance advice:

1. Sequential I/O is faster than random I/O because seeks and rotation dominate.
2. Small random reads waste time in motion, not in the bus.
3. The outer tracks of some disks transfer faster (more sectors per track on older zoned disks).

The kernel elevator and I/O schedulers used to merge and sort requests to reduce seeks. On Linux you still see scheduler names (`none`, `mq-deadline`, `bfq`) under `/sys/block/<dev>/queue/scheduler`. SSDs changed the value of sorting. HDDs still care.

Swap and filesystems on HDDs suffer if the working set is random. That fact is why topic 9 mentioned thrashing as a storage problem, not only a RAM problem.

Treat CHS as history and as a model of mechanical cost. Do not compute CHS for a USB disk.

### Questions

#### Theoretical questions

1. What do cylinder, head, and sector mean on an HDD?
2. What is LBA?
3. Why is sequential I/O faster than random I/O on an HDD?
4. What is rotational latency?
5. Why does the OS no longer program CHS on a current machine?

#### Easy practical tasks

1. Open `man 8 hdparm` if it exists, or `man 7 hd`. Write one sentence about disk parameters.
2. Run `lsblk -d -o NAME,TYPE,ROTA,SIZE`. Write which devices have `ROTA` 1 (rotational).
3. Write five sentences about HDD geometry. Use only facts from this section.
4. Draw a platter, a head arm, and one sector under the head.

#### Medium practical tasks

1. Read `/sys/block/<disk>/queue/scheduler` for one disk. Write the active scheduler (brackets in the file).
2. Compare a sequential `dd` read and a random read tool (`fio` if installed) on a practice file. Write which is slower. Do not flood a shared disk.
3. Read a short note on zoned bit recording or CHS-to-LBA. Write six sentences in your own words.

#### Advanced practical tasks

1. Read Linux `Documentation/block/switching-sched.txt` or the current I/O scheduler docs. Write when `bfq` or `mq-deadline` is a typical choice.
2. Write a one-page history: CHS BIOS limits (8 GB and other old caps) and why LBA and GPT replaced that model.

---

## SSD and wear

A solid-state drive (SSD) stores bits in NAND flash (or similar). There is no seek arm. Random reads are much faster than on an HDD. Sequential writes are still often faster than tiny random writes.

Flash has rules:

1. You read a page.
2. You write a page only after erase.
3. You erase a large erase block, not one page.
4. Each erase block has a limited erase count.

The flash translation layer (FTL) inside the SSD maps host LBAs to physical pages. The FTL does wear leveling so that one block does not die first. The FTL also does garbage collection: it moves live pages so that it can erase a block.

A write from the OS can become a read-modify-erase-write inside the device. That amplification is write amplification. Small random writes increase it.

TRIM (the Linux `discard` operation, `fstrim`) tells the SSD which LBAs the filesystem no longer uses. The FTL can then erase those pages without copy. Without TRIM, the SSD can treat stale data as live.

Wear is the consumption of erase cycles. A full disk that is rewritten often dies sooner than a disk with free space for leveling. Over-provisioning is spare flash that the user does not see.

Do not fill an SSD to 100 percent for a write-heavy lab. Do not run `fstrim` in a tight loop as a benchmark.

The kernel still sees an SSD as a block device of LBAs. The FTL is firmware. You cannot see the physical pages from `read`.

### Questions

#### Theoretical questions

1. Why must an SSD erase before it writes a page again?
2. What is the FTL?
3. What is write amplification?
4. What does TRIM tell the device?
5. Why does free space help wear leveling?

#### Easy practical tasks

1. Run `lsblk -d -o NAME,ROTA,DISC-GRAN,DISC-MAX` if the columns exist. Write whether discard looks available.
2. Open `man 8 fstrim`. Write what the command does.
3. Write five sentences about SSD wear. Use only facts from this section.
4. Make a table: HDD seek cost versus SSD erase-block cost.

#### Medium practical tasks

1. Run `lsblk -o NAME,DISC-ALN,MOUNTPOINT` and `findmnt -o TARGET,OPTIONS` for your root. Write whether you see `discard` in options.
2. Draw host LBA, FTL map, and two physical pages (one stale, one live).
3. Read a vendor or kernel note on discard. Write six sentences: when `fstrim` by timer is safer than `discard` as a mount option.

#### Advanced practical tasks

1. Read about NVMe versus SATA SSD at a high level. Write a one-page contrast: queue depth, command set, and why the OS cares.
2. Write a lab policy: how to avoid wear on a student SSD (no huge random `dd` to the raw device, leave free space).

---

## Block devices vs character devices

Unix splits devices into block devices and character devices. The type is a field of the inode of the device node under `/dev`.

A block device stores a linear array of blocks. The kernel may cache those blocks in the page cache. You can `lseek` and `mmap` many block devices. Examples: `/dev/sda`, `/dev/nvme0n1`, `/dev/loop0`, a partition `/dev/sda1`. Filesystems mount on block devices (or on images that look like them).

A character device transfers data as a stream or through special `ioctl` commands. The kernel does not treat it as a cacheable disk. Examples: `/dev/tty`, `/dev/null`, `/dev/zero`, `/dev/urandom`, many GPIO and serial nodes.

`ls -l /dev` shows `b` for block and `c` for character in the first column. Each node has a major and minor number. The major selects a driver. The minor selects an instance.

Some objects blur the line. `/dev/sda` is a block device. If you `open` it and `read`, you still see raw bytes. Network interfaces are not these nodes in the usual model. They use a different kernel path (topic 18).

Raw disks and `dd` of a full device are dangerous. A wrong destination destroys a filesystem. Practice on a loop file or a disposable VM disk.

Topic 14 covers drivers. This section only classifies the node that a program opens.

### Questions

#### Theoretical questions

1. What is a block device?
2. What is a character device?
3. What do major and minor numbers select?
4. Why can the kernel cache block-device pages?
5. How does `ls -l /dev` show the type?

#### Easy practical tasks

1. Run `ls -l /dev/null /dev/zero /dev/urandom`. Write that they are character devices.
2. Run `ls -l /dev/sda /dev/nvme0n1 /dev/loop0` if they exist. Write the type letter.
3. Open `man 4 random` and `man 4 null`. Write one sentence for each.
4. Write five sentences that contrast the two device types. Use only facts from this section.

#### Medium practical tasks

1. Run `stat /dev/null` and `stat` on one block device. Write major and minor (`stat` prints them as device type).
2. Draw VFS, a block driver, and a character driver. Show which path uses the page cache for file data.
3. Read `man 2 mknod`. Write who may create a device node and why a random file is not a device.

#### Advanced practical tasks

1. Read about `/dev/block` and `/sys/dev/block` on Linux. Write how udev builds `/dev` names from sysfs.
2. Write a one-page warning sheet: `dd` to a block device, `losetup`, and how to name a VM disk so that you never hit the host disk.

---

## Buffer cache / page cache

Linux caches file data in RAM as pages. That cache is the page cache. Older texts say buffer cache. On Linux the page cache is the main cache for file bytes. Buffer heads still track some block metadata.

After `read`, the kernel can keep the page. A later `read` of the same offset can copy from RAM. After `write`, the kernel can mark the page dirty and return to the process. Writeback later sends dirty pages to the device.

Benefits: programs see RAM speed for hot files. Many processes that read the same library share frames (topic 10).

Costs: a `write` that returns success is not always on the device. Memory pressure reclaims clean pages first. Dirty pages must write back before reclaim.

`/proc/meminfo` shows `Cached`, `Buffers`, and `Dirty`. `free -h` shows a cache column. The cache is not waste. It is unused RAM at work. The kernel can drop clean cache when a program needs frames.

`posix_fadvise` and `madvise` give hints (sequential, do not reuse). `O_DIRECT` bypasses the page cache for that I/O. Databases sometimes use `O_DIRECT`. Beginners should not use `O_DIRECT` until they can align buffers.

The page cache is per page (typically 4 KiB). The device may use a different block size. The kernel splits and merges I/O.

Topic 9 described demand paging of file maps. `mmap` of a file uses the same page cache.

### Questions

#### Theoretical questions

1. What is the page cache?
2. What is a dirty page?
3. Why is cache in `free` not lost memory?
4. What does `O_DIRECT` skip?
5. How does `mmap` of a file relate to the page cache?

#### Easy practical tasks

1. Run `grep -E '^(Cached|Buffers|Dirty):' /proc/meminfo`. Write the numbers.
2. Run `free -h`. Write the cache or buff/cache column.
3. Open `man 2 posix_fadvise` or `man 2 madvise`. Write one hint name.
4. Write five sentences about the page cache. Use only facts from this section.

#### Medium practical tasks

1. Read a 100 MiB file twice with `dd` or a small program. Use `/usr/bin/time -v` if you have it. Write whether the second read is faster.
2. Draw: process `write`, dirty page, writeback, device.
3. Drop caches only in a disposable VM (`sync; echo 3 > /proc/sys/vm/drop_caches` as root). Compare a read before and after. Do not do this on a shared host.

#### Advanced practical tasks

1. Read `man 2 open` for `O_DIRECT` alignment rules. Write a one-page note on why beginners should wait.
2. Compare writeback knobs `dirty_ratio` and `dirty_bytes` in `Documentation/admin-guide/sysctl/vm.rst`. Write what happens when dirty pages hit the limit.

---

## `fsync`, durability

Durability is the promise that data still exists after a crash or power loss. A successful `write` does not give that promise. The data can sit in the page cache or in a disk write cache.

`fsync(fd)` asks the kernel to write dirty pages for that file to the device and to complete the device write. On success, the data and the needed metadata for that file are on stable storage, as far as the kernel and the device interface can guarantee.

`fdatasync` can skip some metadata (for example, a timestamp) and still flush the data. `sync` flushes a large set of dirty data. `sync` is a blunt tool.

`O_SYNC` and `O_DSYNC` make each `write` wait for durability. They are slow. Prefer explicit `fsync` at a commit point.

Durability also needs the device. Consumer disks can acknowledge a write while bytes sit in a volatile device cache. The kernel sends flush commands (`REQ_PREFLUSH`, FUA). `hdparm` and mount options can change this behavior. A lie from the device breaks the `fsync` promise.

Applications that need a commit:

1. Write the data.
2. `fsync` the file.
3. Write a journal record or `rename` a complete file.
4. `fsync` the directory when the name must survive (the needed case depends on the filesystem).

Databases and mail systems follow documented recipes. Copy them. Do not invent a half flush.

`close` is not `fsync`. Topic 11 already stated this. This section is the reason.

### Questions

#### Theoretical questions

1. What does `fsync` request?
2. Why is a successful `write` not durable?
3. How does `fdatasync` differ from `fsync`?
4. Why can a device cache break durability?
5. Why is `close` not enough for a commit?

#### Easy practical tasks

1. Open `man 2 fsync` and `man 2 fdatasync`. Write one sentence for each.
2. Open `man 2 open` and find `O_SYNC`. Write what it changes.
3. Write five sentences about durability. Use only facts from this section.
4. Make a table: `write`, `fsync`, `sync`, `close`. Add "waits for device?".

#### Medium practical tasks

1. Write a small C program that `write`s one line and calls `fsync`. Run `strace -e write,fsync` on it.
2. Draw the atomic-replace path from topic 12 and mark each `fsync`.
3. Read a PostgreSQL or SQLite durability note (public docs). Write six sentences about `fsync` and a WAL.

#### Advanced practical tasks

1. Measure `write` of 1 MiB with and without `fsync` after each 4 KiB chunk on a practice file. Write the time ratio. Use a VM disk, not a shared NAS.
2. Read about `DIRSYNC` and directory `fsync` on ext4. Write a one-page recipe for "create file and make the name durable".

---

## RAID levels (high-level)

RAID (redundant array of independent disks) combines more than one disk into one volume. The goals are extra size, extra speed, extra fault tolerance, or a mix. Software RAID on Linux is often `mdadm`. Hardware RAID sits in a controller. The OS then sees one block device.

High-level levels:

- RAID 0 (stripe): split data across disks. More throughput. No redundancy. One disk loss destroys the volume.
- RAID 1 (mirror): the same data on two (or more) disks. One disk can fail. Capacity is that of one disk (for two mirrors).
- RAID 5: stripe plus one parity block per stripe. One disk can fail. Rebuild reads the whole array. Risk grows with large disks.
- RAID 6: stripe plus two parity blocks. Two disks can fail. Rebuild is still heavy.
- RAID 10 (1+0): mirrors that are striped. A common choice when you want speed and simple failure rules.

RAID is not a backup. A delete or a silent corruption can hit every mirror. Keep backups on another system.

RAID does not replace `fsync`. A stripe can still lose a write that never reached any disk.

Rebuild (resync) after a disk change loads the array. Performance drops. A second fault during rebuild of RAID 5 is a known risk on large disks.

This path does not require you to build RAID. You must read `lsblk` and `/proc/mdstat` and know what the level means.

### Questions

#### Theoretical questions

1. What does RAID 0 give and what does it not give?
2. What does RAID 1 store on the second disk?
3. How many disk faults can RAID 5 survive?
4. Why is RAID not a backup?
5. What is a rebuild?

#### Easy practical tasks

1. Open `man 8 mdadm` or `man 4 md`. Write one sentence about Linux software RAID.
2. Run `cat /proc/mdstat` if the file exists. Write whether an array is present.
3. Make a five-row table: level, min disks, fault tolerance, extra space idea.
4. Write five sentences about RAID. Use only facts from this section.

#### Medium practical tasks

1. Draw RAID 1 (two disks) and RAID 0 (two disks) with four blocks of user data.
2. Read a distribution RAID guide. Write six sentences: when RAID 10 is chosen over RAID 5.
3. Compare software RAID and a hardware RAID card in a table: who runs the algorithm, and what the OS sees.

#### Advanced practical tasks

1. In a disposable VM with extra virtual disks, create a RAID 1 with `mdadm`, `mkfs`, mount, then fail one disk in software. Document the steps. Do not do this on a host that you need.
2. Write a one-page policy: RAID plus backups plus `fsync` for a small server.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A program calls `write` then `fsync` on an ext4 file on an SSD behind RAID 1. Which layers can still lose the data, and which layers that call is meant to flush?
2. Why did the page cache make more sense when disks were HDDs, and why does it still matter on SSDs?
3. How do block-device LBAs hide both HDD geometry and SSD FTL pages from the filesystem?
4. When does TRIM matter for an SSD that holds a loop filesystem image?
5. Why can RAID 0 plus a daily copy to another disk still be a valid student lab design?

#### Easy practical tasks

1. Write a one-page cheat sheet: CHS, LBA, FTL, TRIM, block vs character, page cache, dirty, `fsync`, RAID 0/1/5/6/10.
2. Run `lsblk -f`, `free -h`, `cat /proc/meminfo | head`, and `ls -l /dev/null`. Comment each output.
3. Draw one stack: program, page cache, block layer, device cache, NAND or platter.
4. Open `man 2 fsync` and `man 8 fstrim`. Write when a beginner uses each.

#### Medium practical tasks

1. Write a small tool that copies a file with `read`/`write` and an optional `fsync` on the destination. Time both modes on a 20 MiB file.
2. Document your practice machine: rotational or not, filesystem type, whether discard exists, whether `mdstat` is empty.
3. Use `strace -e write,fsync,fdatasync,sync` on `dd` and on your copy tool. Write which process issued `fsync`.

#### Advanced practical tasks

1. Read the OSTEP chapters on I/O devices and filesystems that discuss disks and `fsync`. Write a one-page map to this handbook.
2. In a VM, compare write latency with `fio` (or a loop of `pwrite` plus `fsync`) on the virtual disk. Write how a host cache can make the numbers lie.
