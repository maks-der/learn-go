# 7. Filesystems and Disk

## Description

Topic 6 described the user view of a filesystem: names, inodes, and file descriptors. This topic describes how a Unix filesystem stores that view on a block device, and how the device and the page cache behave. You learn the superblock, the inode table, and data blocks. You also learn journaling, crash consistency, VFS, a short survey of ext4, XFS, NTFS, and APFS, HDD versus SSD, the page cache, `fsync`, and RAID at a high level.

Complete this topic after files and permissions. Complete this topic before devices and I/O (topic 8).

Use one term for each concept. The superblock is the on-disk record of the filesystem as a whole. A journal is a sequential log of intended updates. VFS is the kernel layer that hides the on-disk format from the system-call path. The page cache is the RAM cache of file pages. Durability means the data survives a power loss after the call that promised it. RAID is a layout of many disks that the OS or a controller presents as one volume. Do not mix a journal with user file data. Do not mix `write` with `fsync`.

---

## Superblock, inode table, data blocks

A disk filesystem occupies a range of blocks on a block device. The kernel and `mkfs` divide those blocks into a few roles.

The superblock stores global facts: magic number, block size, inode count, free-block counts, mount state, and pointers to other tables. The kernel reads the superblock at mount. A copy of the superblock often exists in more than one place so that a single bad sector does not destroy the filesystem.

The inode table is an array of inode records. Each inode has a number. The inode stores type, permission bits, owner, size, timestamps, link count, and pointers to data. Topic 6 used the inode as an abstraction. This section names the table that holds those records on disk.

Data blocks hold file contents. A directory is a file. Its data blocks hold names and inode numbers. A regular file stores user bytes in data blocks.

Old Unix layouts used a tree of block pointers inside the inode: direct blocks, then single, double, and triple indirect blocks. Modern ext4 and XFS store extents. An extent is a contiguous run of blocks with a start and a length. Extents shrink metadata for large files.

Free space needs a map. Typical maps are bitmaps (ext4 block groups) or B-trees of free extents (XFS). Allocation policy tries to keep a file in one group so that a later `read` does not jump across the device.

A block group (ext4) packs a bitmap, an inode table slice, and data blocks in one region. Locality is the goal. A small file stays near its inode.

`dumpe2fs` and `tune2fs` show ext4 superblock fields. Use them on a practice image, not on a live root filesystem.

The on-disk layout is not the page cache. A `write` first changes memory. The kernel later writes dirty blocks to the device.

### Questions

#### Theoretical questions

1. What facts does a superblock store?
2. What does one inode record store on disk?
3. What is an extent?
4. Why does a filesystem keep a copy of the superblock?
5. Why is a directory a file with data blocks?

#### Easy practical tasks

1. Open `man 8 mkfs.ext4` and `man 8 dumpe2fs`. Write one sentence for each.
2. Run `df -T` and `findmnt`. Write the type of the root filesystem.
3. Write five sentences: superblock, inode table, data block. Use only facts from this section.
4. Draw three boxes on a disk line: superblock, inode table, data blocks.

#### Medium practical tasks

1. Create a 64 MiB file. Run `mkfs.ext4` on it. Run `dumpe2fs` on the image. Write block size and inode count.
2. Mount the image with `sudo mount -o loop` on a practice directory if you have rights. Create a file. Unmount. Write the inode number from `debugfs` or `ls -i` while mounted.
3. Compare a 1-byte file and a 1 MiB file with `filefrag -v` if the tool exists. Write how many extents you see.

#### Advanced practical tasks

1. Read the ext4 disk layout notes in the kernel documentation. Write a one-page map of block groups, bitmaps, and the inode table.
2. Compare indirect-block trees and extents in a table: metadata size for a 1 GiB sequential file, and cost of a random hole.

---

## Journaling and crash consistency

Crash consistency is the property that after a power loss or a kernel panic, the filesystem metadata (and sometimes data) returns to a defined good state. Without a plan, a crash in the middle of a create can leave a directory entry without an inode, or an inode without a directory name.

A journaling filesystem writes a description of a change to a journal before it writes the main metadata (and sometimes the file data). After a crash, the kernel replays or discards incomplete journal records. The replay brings the filesystem to a consistent metadata state without a full `fsck` walk of every inode.

A journal is a circular log of transactions. A transaction groups related updates: allocate a block, attach it to an inode, and update the directory. The journal commit marks the transaction as complete.

Common journal modes (ext4 names):

1. `data=ordered` (typical default): file data goes to the device before the related metadata hits the journal commit. Metadata is journaled. This mode is a balance of safety and speed.
2. `data=writeback`: metadata is journaled. File data can reach the device in any order relative to that metadata. After a crash, a file can contain old or stale bytes.
3. `data=journal`: data and metadata go through the journal. This mode is the slowest and the most conservative for those updates.

Journaling is not a full backup. It does not keep every old version of user data. It does not replace `fsync` for application durability.

XFS uses a different log design (delayed logging). The idea is the same: a sequential log of metadata intent, then write the real structures.

Copy-on-write filesystems (btrfs, some others) update new blocks and then swing a root pointer. That is another crash-consistency design. This path only names the idea.

A journal has a finite size. A huge rename storm can fill it and stall writers. That stall is a performance fact, not a logic error.

### Questions

#### Theoretical questions

1. What is crash consistency?
2. What does a journal store?
3. What is the difference between `data=ordered` and `data=writeback` on ext4?
4. Why is a journal not a backup?
5. Why can a full journal stall writers?

#### Easy practical tasks

1. Open `man 5 ext4` if it exists, or a distribution ext4 page. Write the default data mode if the page states it.
2. Write five sentences about journaling and crash consistency. Use only facts from this section.
3. Run `findmnt -o TARGET,FSTYPE,OPTIONS /` . Write whether you see a `data=` option.
4. Draw: dirty metadata, journal write, journal commit, write to the real inode.

#### Medium practical tasks

1. Compare `fsck` after a crash on a journaled filesystem versus an old non-journaled filesystem (textbook or man pages). Write six sentences.
2. Read `man 8 tune2fs` for journal options. Write one command that you must not run on a live root without a plan.
3. Make a table: journaled metadata, ordered data, full data journal. Add crash risk for file bytes.

#### Advanced practical tasks

1. Read an OSTEP or kernel note on journaling. Write a one-page example of a create that needs three writes to stay consistent.
2. Write a short comparison of journaling and copy-on-write filesystems. Use public docs. Do not claim you implemented either.

---

## VFS; ext4, XFS, NTFS, APFS (survey)

The virtual filesystem (VFS) is the Linux kernel layer under `open`, `read`, and `stat`. VFS calls filesystem-specific operations. Your program does not name ext4 in the `read` call. A mount attaches a concrete filesystem to a directory.

VFS objects include inodes, dentries (directory cache entries), file objects, and superblocks in memory. The dentry cache speeds path walks. Topic 6 described the user path. VFS implements that walk.

A survey of formats:

- ext4 is a common Linux disk filesystem. It uses extents and a journal. Many distributions use it for the root filesystem.
- XFS is another Linux filesystem. It is strong for large files and parallel allocation. It uses a metadata log. Some distributions use XFS by default.
- NTFS is the usual Windows disk filesystem. Linux can mount it with drivers (`ntfs3` or others). Features and permission models differ from POSIX. Do not treat a USB NTFS disk as a full POSIX tree.
- APFS is the usual Apple filesystem on recent macOS. It uses copy-on-write and snapshots. Linux support is limited. This path does not require an APFS lab.

Other names that you will see: FAT and exFAT on USB sticks, btrfs, ZFS, tmpfs (RAM), NFS and 9p (network). tmpfs is not a disk format. NFS is not a local superblock layout.

Pick a format for the job: USB exchange, Linux root, large data volume, Windows dual boot. This section is a survey. You do not memorize every on-disk field.

`df -T` and `/proc/filesystems` show types that this kernel knows. `mount` and `findmnt` show what is attached.

### Questions

#### Theoretical questions

1. What job does VFS do?
2. Why does a program not name ext4 in `read`?
3. Name one typical use of ext4 and one typical use of XFS.
4. Why is NTFS on Linux not a full POSIX filesystem?
5. What is tmpfs?

#### Easy practical tasks

1. Run `df -T` and `cat /proc/filesystems | head`. Write three types.
2. Open `man 8 mount` and `man 5 fstab`. Write one sentence for each.
3. Write five sentences about VFS and the four named formats. Use only facts from this section.
4. Make a four-row table: ext4, XFS, NTFS, APFS. Add "usual home OS" and "journal or COW".

#### Medium practical tasks

1. Draw user `read`, VFS, ext4, block device.
2. Read a short official note on XFS versus ext4 (distribution or kernel docs). Write three differences.
3. List mounts with `findmnt -D`. Write which ones are disk and which are virtual (`proc`, `sysfs`, `tmpfs`).

#### Advanced practical tasks

1. Read the Linux VFS overview in kernel documentation. Write a one-page summary of inode and dentry.
2. Write a policy: which filesystem you pick for a Linux data disk, a USB stick for Windows, and a VM scratch disk. Give one reason each.

---

## HDD vs SSD, page cache, `fsync`

A hard disk drive (HDD) stores bits on spinning magnetic platters. A seek moves the head. Rotational latency is the wait for the sector. Sequential I/O is much faster than random I/O on an HDD.

A solid-state drive (SSD) stores bits in flash. There is no seek arm. Random reads are much faster than on an HDD. Flash must erase before a rewrite. Wear leveling spreads writes. An SSD still has a larger latency than RAM. The device can have its own write cache.

The page cache is the RAM cache of file pages on Linux. A `read` can hit RAM after the first load. A `write` can return when the data is in RAM. The kernel writes dirty pages later (writeback). `free -h` shows cache. Topic 1 already named this cache.

`fsync` on a file descriptor asks the kernel to write the file data and the metadata that the file needs to a durable device state (as far as the OS can promise). `fdatasync` can skip some metadata. `sync` flushes more of the system. After `write` only, a crash can lose the last bytes.

Durability also depends on the disk write cache and on barriers. `fsync` is necessary. It is not always sufficient if the hardware lies. For this path, treat `fsync` as the application call for "this FD must hit storage".

Do not call `fsync` after every byte in a hot loop unless you measured the need. Do not skip `fsync` on a commit record if you promised durability.

HDD versus SSD still matters for the cache: an SSD recovers faster from cache misses. An HDD random miss is very slow. Sequential `read` after a cold cache still hits the device.

`posix_fadvise` and `madvise` can hint sequential access or cache drop. They are hints.

### Questions

#### Theoretical questions

1. Why is sequential I/O faster than random I/O on an HDD?
2. What flash rule makes small random writes costly on an SSD?
3. What is the page cache?
4. What does `fsync` ask the kernel to do?
5. Why can `write` return before the device has the bytes?

#### Easy practical tasks

1. Run `lsblk -d -o NAME,ROTA,SIZE,MODEL`. Write which devices look rotational.
2. Run `free -h`. Write the cache or buff/cache value.
3. Open `man 2 fsync` and `man 2 fdatasync`. Write one sentence for each.
4. Write five sentences: HDD, SSD, page cache, `fsync`. Use only facts from this section.

#### Medium practical tasks

1. Time a write of a 64 MiB file without `fsync` and with `fsync` (or `dd` with `conv=fdatasync`). Write the two times. Remove the file.
2. Draw user `write`, page cache, writeback, device, `fsync`.
3. Read `/sys/block/<dev>/queue/scheduler` for one disk. Write the active scheduler.

#### Advanced practical tasks

1. Write a tiny logger that `write`s a line and `fsync`s. Compare throughput with a logger that never syncs. This handbook does not contain the source.
2. Read a short note on SSD wear and TRIM (`fstrim`). Write six sentences. Do not run destructive trim experiments on a shared disk.

---

## RAID (high-level)

RAID (redundant array of independent disks) presents many disks as one volume. The goals are extra capacity, extra speed, extra fault tolerance, or a mix. This section is high-level. You do not design a production array here.

Common textbook levels:

- RAID 0 (stripe): split data across disks. More throughput. No redundancy. One disk fail loses the volume.
- RAID 1 (mirror): the same data on two (or more) disks. A one-disk fail can still serve reads. Capacity is that of one disk (for two mirrors).
- RAID 5: stripe plus one parity piece per stripe. One disk fail is recoverable. Rebuild reads the other disks. A second fail during rebuild is a risk.
- RAID 6: two parity pieces. Two disk fails are recoverable. More write cost.
- RAID 10: mirrors that you stripe. A common production mix of speed and redundancy.

Software RAID on Linux is `mdadm`. Hardware RAID is a controller. The kernel can also use device mapper. The filesystem sits on the RAID volume. Journaling and RAID are different layers. RAID does not replace backups. RAID does not fix a file that you deleted. RAID does not fix silent corruption unless you add checksums (other systems).

Rebuild time on large disks is long. During rebuild the array is stressed. Monitor disks. Replace a failed disk with a plan.

Do not build RAID 0 for data that you cannot lose. Do not treat a USB stick pair as a serious array.

### Questions

#### Theoretical questions

1. What illusion does RAID give to the filesystem?
2. What does RAID 0 cost you on a disk fail?
3. How does RAID 1 use capacity?
4. Why is RAID 5 rebuild a risky window?
5. Why does RAID not replace backups?

#### Easy practical tasks

1. Write a five-row table: RAID 0, 1, 5, 6, 10. Add "survives one disk fail?" and "capacity idea".
2. Open `man 8 mdadm` if it exists. Write one sentence about what the tool manages.
3. Write five sentences about RAID. Use only facts from this section.
4. Draw RAID 1: two disks, one volume, a filesystem on top.

#### Medium practical tasks

1. Read a distribution mdadm guide (read only). Write the difference between a degraded array and a clean array.
2. Compare software RAID and hardware RAID in six sentences: who runs the XOR, and who sees `/dev/md0`.
3. Write a policy: which RAID level you would pick for a two-disk home server and why. State that this is not a production design review.

#### Advanced practical tasks

1. Read a short note on RAID and `fsync` (write hole, write cache). Write why a power loss can still hurt an array.
2. Write a one-page note: RAID versus a replicated filesystem or a backup tool. Name what each layer fixes.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do superblock, inode table, and VFS fit one `stat` call?
2. After a crash, what does the journal repair, and what does `fsync` still need to promise for an application?
3. When do you pick ext4, XFS, or a USB FAT/exFAT volume?
4. How do HDD seeks, SSD erase rules, and the page cache change the cost of a random `write`?
5. Which failures does RAID hide, and which failures does it not hide?

#### Easy practical tasks

1. Write a cheat sheet: superblock, inode table, extent, journal, VFS, ext4, XFS, NTFS, APFS, HDD, SSD, page cache, `fsync`, RAID 0/1/5.
2. Run `df -hT`, `free -h`, and `lsblk`. Write one line that explains each command.
3. Draw layers: process, VFS, filesystem, RAID volume, disk.
4. Bookmark the ext4 and XFS man pages or kernel docs that you used.

#### Medium practical tasks

1. Build a loop ext4 image, mount it, copy a file, `umount`, `dumpe2fs`. Write a six-line lab report.
2. Time `dd` with and without `conv=fdatasync` on the image or on `/tmp`. Write the difference.
3. Document a crash-consistency checklist for an app: `write`, `fsync`, `rename`, journal mode of the disk.

#### Advanced practical tasks

1. Read the OSTEP persistence chapters that match journaling and crash consistency. Write a one-page map to ext4 journal modes.
2. Design (on paper) a two-disk Linux data volume: filesystem type, RAID or not, backup, and `fsync` policy for a small database file.
