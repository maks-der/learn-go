# 12. Filesystems (Implementation)

## Description

Topic 11 described the user view of a filesystem: names, inodes, and file descriptors. This topic describes how a Unix filesystem stores that view on a block device. You learn the superblock, the inode table, and data blocks. You also learn journaling, crash consistency, the virtual filesystem (VFS) layer, and a short survey of common on-disk formats. Mounts and mount namespaces appear as a preview.

Complete this topic after filesystem abstraction. Complete this topic before storage devices (topic 13). You already know `open` and `read`. You now learn where the bytes live after a crash.

Use one term for each concept. The superblock is the on-disk record of the filesystem as a whole. An inode is the on-disk record of one file. A data block holds file bytes or directory entries. A journal is a sequential log of intended updates. VFS is the kernel layer that hides the on-disk format from the system-call path. Do not mix an inode with a directory entry. Do not mix a journal with the user file data.

---

## Superblock, inode table, data blocks

A disk filesystem occupies a range of blocks on a block device. The kernel and `mkfs` divide those blocks into a few roles.

The superblock stores global facts: magic number, block size, inode count, free-block counts, mount state, and pointers to other tables. The kernel reads the superblock at mount. A copy of the superblock often exists in more than one place so that a single bad sector does not destroy the filesystem.

The inode table is an array of inode records. Each inode has a number. The inode stores type, permission bits, owner, size, timestamps, link count, and pointers to data. Topic 11 used the inode as an abstraction. This section names the table that holds those records on disk.

Data blocks hold file contents. A directory is a file. Its data blocks hold names and inode numbers. A regular file stores user bytes in data blocks.

Old Unix layouts used a tree of block pointers inside the inode: direct blocks, then single, double, and triple indirect blocks. Modern ext4 and XFS store extents. An extent is a contiguous run of blocks with a start and a length. Extents shrink metadata for large files.

Free space needs a map. Typical maps are bitmaps (ext4 block groups) or B-trees of free extents (XFS). Allocation policy tries to keep a file in one group so that a later `read` does not jump across the device.

A block group (ext4) packs a bitmap, an inode table slice, and data blocks in one region. Locality is the goal. A small file stays near its inode.

`dumpe2fs` and `tune2fs` show ext4 superblock fields. `debugfs` can dump an inode. Use them on a practice image, not on a live root filesystem.

The on-disk layout is not the page cache. Topic 13 covers the cache. This section only states that a `write` first changes memory. The kernel later writes dirty blocks to the device.

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
2. Mount the image with `sudo mount -o loop` on a practice directory. Create a file. Unmount. Run `debugfs -R 'stat <path>'` or an equivalent dump. Write the inode number.
3. Compare a 1-byte file and a 1 MiB file with `debugfs` or `filefrag -v`. Write how many extents or blocks you see.

#### Advanced practical tasks

1. Read the ext4 disk layout notes in the kernel documentation. Write a one-page map of block groups, bitmaps, and the inode table.
2. Compare indirect-block trees and extents in a table: metadata size for a 1 GiB sequential file, and cost of a random hole.

---

## Journaling

A journaling filesystem writes a description of a change to a journal before it writes the main metadata (and sometimes the file data). After a crash, the kernel replays or discards incomplete journal records. The replay brings the filesystem to a consistent metadata state without a full `fsck` walk of every inode.

A journal is a circular log of transactions. A transaction groups related updates: allocate a block, attach it to an inode, and update the directory. The journal commit marks the transaction as complete.

Common journal modes (ext4 names):

1. `data=ordered` (typical default): file data goes to the device before the related metadata hits the journal commit. Metadata is journaled. This mode is a balance of safety and speed.
2. `data=writeback`: metadata is journaled. File data can reach the device in any order relative to that metadata. After a crash, a file can contain old or stale bytes.
3. `data=journal`: data and metadata go through the journal. This mode is the slowest and the most conservative for those updates.

Journaling is not a full backup. It does not keep every old version of user data. It does not replace `fsync` for application durability (topic 13).

XFS uses a different log design (delayed logging). The idea is the same: a sequential log of metadata intent, then write the real structures.

A journal has a finite size. A huge rename storm can fill it and stall writers. That stall is a performance fact, not a logic error.

Do not disable the journal on a practice root filesystem. If you experiment, use a loop image.

### Questions

#### Theoretical questions

1. What problem does a journal solve after a crash?
2. What is a journal transaction?
3. What is the difference between `data=ordered` and `data=writeback`?
4. Why is journaling not a backup?
5. What happens when the journal is full?

#### Easy practical tasks

1. Open `man 5 ext4` or `man 8 tune2fs`. Write the three data journaling modes.
2. Run `sudo tune2fs -l` on a practice ext4 image (not a shared host root unless you only read). Write whether a journal exists.
3. Write five sentences about journaling. Use only facts from this section.
4. Make a three-row table: mode, what is journaled, crash risk for file bytes.

#### Medium practical tasks

1. Draw the order of writes for `data=ordered`: file data, journal metadata, journal commit, then checkpoint of metadata.
2. Read a short OSTEP or kernel note on journaling. Write six sentences that separate "consistent metadata" from "the last `write` is durable".
3. Compare ext4 journaling with XFS logging in public documentation. Write two similarities and one difference.

#### Advanced practical tasks

1. On a disposable VM, create an ext4 image, mount with `data=writeback`, and document the mount option from `findmnt -o OPTIONS`. Do not force a crash on a host that you need.
2. Write a one-page note: which application bugs a journal does not fix (missing `fsync`, torn application records).

---

## FAT, ext4, XFS, NTFS, APFS (survey)

Many on-disk formats exist. A survey names the idea of each format. It does not replace a specification.

FAT (FAT16, FAT32, exFAT) stores a file allocation table. The table maps each cluster to the next cluster or to an end mark. FAT has no Unix inode table. Names and size live in directory entries. Permissions and hard links are weak or absent. Firmware and USB media still use FAT. The EFI System Partition is often FAT.

ext4 is a common Linux filesystem. It uses inodes, extents, a journal, and block groups. It supports Unix permissions, hard links, and sparse files. Many distributions use ext4 for the root filesystem.

XFS is a Linux filesystem for large volumes and parallel allocation. It uses B-trees and an extent model. It journals metadata. It grew in data-center and media workloads.

NTFS is the Windows NT filesystem. It uses a Master File Table (MFT). Each file has an MFT record. NTFS supports ACLs, alternate streams, and compression. Linux can mount NTFS through drivers. The on-disk ideas differ from ext4.

APFS is the Apple filesystem. It uses copy-on-write structures and snapshots. Clones and space sharing are design goals. Linux support is limited. You study APFS as a contrast: crash consistency by copy-on-write instead of a classic journal.

Other Linux names that you see: btrfs and Bcachefs (copy-on-write, checksums), F2FS (flash-friendly), tmpfs (RAM), overlayfs (layers). This path does not require you to master them.

Pick a format by the OS and the media. A USB stick that must boot firmware wants FAT. A Linux server disk often wants ext4 or XFS. Do not format a disk that you need.

### Questions

#### Theoretical questions

1. How does FAT find the next cluster of a file?
2. What Unix features does ext4 store that FAT does not store well?
3. What is the MFT in NTFS?
4. How does APFS differ from a classic journaled filesystem at a high level?
5. When is FAT still the correct choice?

#### Easy practical tasks

1. Run `findmnt -t ext4,xfs,vfat,ntfs,btrfs`. Write the types that exist on your system.
2. Open `man 5 ext4` and `man 5 xfs` if present. Write one sentence for each.
3. Make a five-row table: FAT, ext4, XFS, NTFS, APFS. Add "home OS" and "inode or equivalent".
4. Write five sentences that survey the five formats. Use only facts from this section.

#### Medium practical tasks

1. Create a small FAT image with `mkfs.vfat` and an ext4 image. Mount both. Try `ln` (hard link) on each. Write the result.
2. Draw FAT as a linked list of clusters. Draw ext4 as an inode with one extent.
3. Read a public NTFS or APFS overview. Write three features that POSIX `stat` does not show.

#### Advanced practical tasks

1. Read kernel documentation for one extra Linux filesystem (btrfs or F2FS). Write a one-page contrast with ext4: consistency method and target media.
2. Write a decision note for a lab machine: root on ext4, a large data disk on XFS or ext4, and a USB installer on FAT. Give one reason for each.

---

## VFS layer

The virtual filesystem (VFS) is the kernel layer above on-disk filesystems. User programs call `open`, `read`, `stat`, and `mount`. Those calls enter VFS. VFS then calls the filesystem module for the mount that owns the path.

VFS objects on Linux include:

- superblock object (one mounted instance)
- inode object (one file in memory)
- dentry (a cached name-to-inode binding)
- file object (an open file in a process)

A dentry cache (dcache) remembers recent path walks. A page cache remembers file data pages. Both sit behind VFS. Two different on-disk formats still share these caches.

VFS lets you mount ext4, XFS, procfs, and tmpfs in one tree. The path `/proc/self/status` and `/etc/passwd` use the same `open` path. The inode operations differ.

A filesystem module registers a type (`ext4`, `xfs`). `mount` binds a device (or no device) to a directory with that type. Stacked filesystems such as overlayfs sit on other mounts.

Network filesystems (NFS, SMB) also plug into VFS. The "block" behind the inode can be a remote server. The system-call interface stays the same. Failure modes differ (timeouts, stale handles).

Do not think that VFS is a second on-disk format. VFS is in-memory glue. The on-disk format stays in the filesystem driver.

### Questions

#### Theoretical questions

1. What job does VFS do?
2. What is a dentry?
3. Why can `open` work on `/proc` and on ext4 with the same call?
4. What is the difference between a VFS inode and an ext4 on-disk inode?
5. Why is VFS not an on-disk format?

#### Easy practical tasks

1. Open `man 2 mount` and `man 8 mount`. Write who calls each interface.
2. Run `cat /proc/filesystems`. Write five types. Mark which look virtual (nodev).
3. Write five sentences about VFS. Use only facts from this section.
4. Draw: process, VFS, ext4 module, block device.

#### Medium practical tasks

1. Compare `stat /proc/self` and `stat /etc/passwd`. Write which fields still make sense for procfs.
2. Read a Linux VFS overview in kernel documentation. Write the roles of `file`, `dentry`, and `inode` in six sentences.
3. List three filesystem types from `findmnt` and name the VFS object that is common to all three.

#### Advanced practical tasks

1. Read `Documentation/filesystems/vfs.rst` in the kernel tree (online is enough). Write a one-page summary of inode operations versus file operations.
2. Trace `open` of a file with `ftrace` or `bpftrace` if you have rights (topic 20). Write two VFS function names that appear.

---

## Crash consistency

Crash consistency is the property that a sudden power loss or kernel panic leaves the filesystem in a state that the kernel can mount. Metadata must not disagree with itself. Examples of disagreement: a block that is both free and pointed to by an inode, or a directory entry that names a missing inode.

A crash can stop a multi-step update in the middle. Without a plan, `fsck` must walk the volume and repair graphs. That walk is slow on large disks.

Techniques:

1. Journaling: replay a log of metadata intent (this topic).
2. Copy-on-write (COW) filesystems: write new trees, then atomically swing a root pointer (btrfs, APFS).
3. Soft updates and other ordered-write schemes (historical BSD).
4. Careful write order plus `fsck` after crash (old ext2).

Crash consistency is not application consistency. A database that writes a record in two `write` calls can still lose the second call. The filesystem can be consistent while the application file is a torn record. Applications use `fsync`, write-ahead logs, or atomic `rename` of a full file.

Atomic `rename` on the same filesystem is a common idiom: write a temporary file, `fsync` it, then `rename` over the old name. Directory updates are ordered so that the new name appears as one step. This idiom still needs a correct filesystem.

`fsck` (or `fsck.ext4`) still exists. A journal replay can fail. A bug or a torn write in the journal can need a repair. Run `fsck` only on an unmounted filesystem unless the tool documents a read-only check.

Do not pull power on a host to "test crash consistency". Use a VM snapshot or a documented fault-injection lab.

### Questions

#### Theoretical questions

1. What does crash consistency mean for metadata?
2. Give one example of inconsistent metadata.
3. Why is crash consistency not the same as application durability?
4. How does a COW filesystem make a new tree visible?
5. When do you run `fsck` on ext4?

#### Easy practical tasks

1. Open `man 8 fsck.ext4`. Write when the tool expects the filesystem to be unmounted.
2. Write five sentences about crash consistency. Use only facts from this section.
3. Make a table: "filesystem consistent" versus "application record complete". Add one example each.
4. Draw a two-step update that a crash can split.

#### Medium practical tasks

1. Write the atomic-replace steps: create temp file, write, `fsync`, `rename`. Label which step needs the same filesystem.
2. Read OSTEP on crash consistency or journaling. Write six sentences in your own words.
3. Compare journal replay with a full `fsck` in a short table: time, what is checked, when each runs.

#### Advanced practical tasks

1. Read about ext4 `fsync` and `rename` in kernel or LWN notes. Write a one-page guide for a beginner editor that must not lose a file.
2. In a disposable VM, document how the hypervisor can pull power from a guest disk. Do not run the test on a machine that you need. Write the planned observation list.

---

## Mounts and namespaces (preview)

A mount attaches a filesystem instance to a directory in the process namespace. After `mount`, that directory tree shows the new filesystem. The old contents of the mount point are hidden until unmount.

`mount` and `umount` are system calls and also command names. Typical Linux command:

```text
mount -t ext4 /dev/sdb1 /mnt/data
```

`/proc/mounts` and `findmnt` list the mounts that the caller sees. `/etc/fstab` lists mounts that `init` should make at boot.

A bind mount attaches an existing directory tree at a second path. The two paths see the same files. Bind mounts are common in containers.

A mount namespace is a per-process view of the mount table. Two processes can see two different trees. Linux creates this isolation with `clone` flags or `unshare --mount`. Topic 16 covers namespaces in full. This section only states that mounts are not global on modern Linux.

Propagation flags (`shared`, `private`, `slave`) control whether a mount event in one namespace appears in another. Container runtimes set these flags. Beginners can treat them as "mounts can be hidden from other processes".

`chroot` changes the apparent root for a process. It is not a mount namespace. It is a weaker isolation tool. Do not treat `chroot` as a security boundary.

Unmount fails when a process has a working directory or an open file on the mount. `lsof` or `fuser` finds those users (topic 20).

### Questions

#### Theoretical questions

1. What does a mount attach?
2. What does a bind mount do?
3. What is a mount namespace?
4. Why is `chroot` not a mount namespace?
5. Why can `umount` fail?

#### Easy practical tasks

1. Run `findmnt`. Write the source and target of the root mount.
2. Open `man 8 mount` and `man 2 mount`. Write one fact from each.
3. Run `cat /proc/mounts | head`. Write three lines in your own words.
4. Write five sentences about mounts. Use only facts from this section.

#### Medium practical tasks

1. In a user namespace or with sudo in a VM, bind-mount a practice directory on another directory. Write `findmnt` before and after.
2. Draw two processes and two mount namespaces. Show one private mount that the other process does not see.
3. Read `man 7 mount_namespaces`. Write six sentences that preview topic 16.

#### Advanced practical tasks

1. Use `unshare --mount --propagation private` in a VM. Mount a tmpfs. Confirm that the host tree does not show it. Document every command.
2. Write a one-page comparison: `chroot`, mount namespace, and a full container (namespaces plus cgroups). Topic 16 expands this.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A classmate says "VFS stores the inode table on disk." Which facts do you use to correct that sentence?
2. How do the superblock, the journal, and crash consistency work together at mount after a power loss?
3. Why can two processes that call the same `read` still hit two different on-disk formats?
4. When do you pick ext4, XFS, or FAT for a new volume on Linux?
5. How does a mount namespace change what "the filesystem tree" means after topic 11?

#### Easy practical tasks

1. Write a one-page cheat sheet: superblock, inode table, data block, extent, journal, VFS, dentry, crash consistency, mount, bind mount.
2. Run `df -T`, `findmnt`, `cat /proc/filesystems`, and `lsblk -f` if present. Write one line that explains each command.
3. Draw one figure: disk layout, journal, VFS, and a mount point.
4. Bookmark the ext4 and VFS kernel documentation pages. Write one sentence about when you open each page.

#### Medium practical tasks

1. Build a loop ext4 image, mount it, copy a file, unmount, run `fsck.ext4 -n` on the image. Save the commands and the `fsck` summary.
2. Write a short lab script that prints filesystem type, mount options, and whether `/proc` is `nodev`. Run it. Save the output.
3. Use `strace -e mount,umount2` on a harmless `findmnt` or document why those calls do not appear. Write what you learned about who mounts.

#### Advanced practical tasks

1. Read the OSTEP chapters on filesystems and crash consistency. Write a one-page map from those chapters to each section in this handbook.
2. In a disposable VM, add a second virtual disk, partition it, `mkfs.ext4`, add an `/etc/fstab` line, reboot, and document `findmnt`. Do not do this on a host that you need.
