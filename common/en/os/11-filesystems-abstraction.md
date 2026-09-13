# 11. Filesystems (Abstraction)

## Description

A filesystem is the OS abstraction for named, persistent data. This topic explains files, directories, and paths. You learn absolute and relative paths, the Unix inode, hard links and symbolic links, and permission bits. You also learn file descriptors and the calls `open`, `read`, `write`, `close`, and `lseek`.

Complete this topic after processes and memory maps. A file descriptor lives in the process. Complete this topic before filesystem implementation (topic 12). This topic stays at the user-visible interface.

Use one term for each concept. A file is a named inode with data. A directory is a file that maps names to inode numbers. A path is a string of names. A file descriptor is an integer in one process. A hard link is another name for the same inode. A symbolic link is a small file that stores a path. Do not mix a path with an inode. Do not mix a file descriptor with an inode number.

---

## File, directory, path

A file is the object that the OS names and that a program reads or writes. On Unix, a file is not only a document. A directory, a device node, a socket, and a named pipe are files in the inode sense. The type is a field of the inode.

A directory is a file that stores a list of names. Each name points to an inode number. The kernel uses the directory to resolve a path.

A path is a sequence of names that the kernel walks. The names are separated by `/` on Unix. The path `/home/ada/notes.txt` means: start at the root directory, then `home`, then `ada`, then `notes.txt`.

The kernel does not store the full path inside the inode. Many paths can reach the same inode (hard links). The inode stores the type, the permission bits, the size, the timestamps, and the pointers to data.

User programs see a tree. The tree starts at `/`. A mount can attach another filesystem at a directory. Topic 12 covers mounts. This topic treats one tree.

`.` is the current directory in a path. `..` is the parent directory. The root's `..` is the root.

Do not think of a file as "bytes on one disk sector". The file is the inode plus the data blocks plus at least one name in a directory (except some unnamed objects).

### Questions

#### Theoretical questions

1. What is a file on Unix at a high level?
2. What does a directory store?
3. What is a path?
4. Why does the inode not store the full path as its identity?
5. What do `.` and `..` mean?

#### Easy practical tasks

1. Run `ls -l /` . Write five names and whether each looks like a directory (`d` in the first column).
2. Run `man 7 path_resolution` if the page exists, else `man 1 ls`. Write one fact about path walk.
3. Write five sentences: file, directory, path. Use only facts from this section.
4. Draw a tree: `/`, `home`, `ada`, `notes.txt`.

#### Medium practical tasks

1. Run `stat /` and `stat /tmp`. Write the inode number and the file type for each.
2. Use `find /tmp -maxdepth 1 | head` and explain that `find` walks names, not raw inodes.
3. Compare `ls` and `ls -d /tmp`. Write what `-d` changes.

#### Advanced practical tasks

1. Read `man 7 inode` or `man 7 path_resolution`. Write the steps of a path walk in your own words.
2. Write a one-page note: which objects in `/dev` are files, and why Unix calls them files.

---

## Absolute vs relative paths

An absolute path starts with `/`. The kernel starts the walk at the root of the process (the root directory, which `chroot` can change). Example: `/etc/passwd`.

A relative path does not start with `/`. The kernel starts the walk at the current working directory (CWD) of the process. Example: `notes.txt` or `../bin/app`.

The CWD is process state. `chdir` changes it. A child inherits the CWD at `fork`. `exec` keeps the CWD.

The same relative path can name two different files if two processes have two CWDs. The same absolute path names the same file (if the root is the same and nothing raced).

Scripts should use absolute paths for system files (`/etc`, `/usr`). Relative paths are right for files next to the user project. A program that opens `config.ini` without a directory depends on how the user started it. That dependence is a common bug.

`realpath` and `readlink -f` print a canonical absolute path when the path exists. Symbolic links can make the canonical path look different from the path that the user typed.

`PATH` is not a filesystem path of a file. `PATH` is a list of directories that the shell searches for a command name. Do not mix `PATH` with the path of an open file.

### Questions

#### Theoretical questions

1. How does an absolute path start on Unix?
2. Where does a relative path walk begin?
3. What system call changes the CWD?
4. Why can a relative path be a bug in a program that you start from any directory?
5. How does `PATH` differ from a file path?

#### Easy practical tasks

1. Run `pwd` and `ls notes-does-not-exist` in your home directory. Then run `ls /etc/hostname` (or `/etc/os-release`). Write which path was absolute.
2. Open `man 2 chdir` and `man 3 getcwd`. Write one sentence for each.
3. Write four sentences about absolute and relative paths. Use only facts from this section.
4. Make a table: "Path string", "Absolute?", "Start of walk". Add three rows.

#### Medium practical tasks

1. Write a C program that prints `getcwd` and then `chdir("/tmp")` and prints again. Run it. Confirm the shell CWD did not change.
2. From two different directories, run `cat os-practice/README` versus `cat "$HOME/os-practice/README"`. Write which form is stable.
3. Draw the walk of `../etc/passwd` from `/usr/bin`.

#### Advanced practical tasks

1. Read `man 7 path_resolution` on `AT_FDCWD` and `openat`. Write why `openat` is safer than `chdir` plus `open` in a multithreaded process.
2. Use `realpath` on a path with `..` and a symlink. Write the input and the output.

---

## Inode (Unix)

An inode is the kernel and filesystem object that stores the metadata of a file. Each inode has a number in that filesystem. `ls -i` and `stat` print the number.

The inode holds:

- file type (regular, directory, symlink, device, ...)
- permission bits and owner
- link count (how many hard links)
- size
- timestamps (access, change, modification — exact fields depend on the filesystem)
- pointers to data (or to extents; topic 12)

The inode does not hold the name. The directory holds the name and the inode number.

Commands:

```text
stat file
ls -i file
```

Two names with the same device and the same inode number are the same file (hard links). Two files on two filesystems can show the same inode number. The pair (filesystem, inode number) is the identity.

When the last hard link is removed and no process has the file open, the filesystem can free the inode and the data. If a process still has the file open, the data can stay until `close`. That is why a deleted log can still grow.

`stat` versus `lstat`: `lstat` does not follow a symbolic link. `stat` follows it and shows the target inode.

Topic 12 describes how ext4 stores the inode table. This section only needs the abstraction.

### Questions

#### Theoretical questions

1. What metadata does an inode store?
2. Where is the file name stored?
3. Why is the identity (filesystem, inode number) and not the inode number alone?
4. When can the filesystem free an inode after `unlink`?
5. What is the difference between `stat` and `lstat`?

#### Easy practical tasks

1. Run `stat /etc/hostname` or `stat /etc/os-release`. Write inode, size, and links.
2. Run `ls -i /etc/hostname /etc/os-release`. Write the two numbers.
3. Open `man 2 stat` and `man 1 stat`. Write which one is a system call.
4. Write five sentences about inodes. Use only facts from this section.

#### Medium practical tasks

1. Create a file. Run `stat`. Remove the file while `sleep 60` has it open (`sleep 60 < file` in another test). Compare `ls` and `ls -l /proc/<pid>/fd`. Write what you see.
2. Draw an inode box and two directory entries that point to it.
3. Compare `stat` and `stat -L` on a symlink (or `stat` versus `lstat` in C). Write the inode numbers.

#### Advanced practical tasks

1. Read `man 7 inode`. Write a table of file types (`S_IFREG`, `S_IFDIR`, `S_IFLNK`).
2. Write a C program that calls `stat` and prints `st_ino`, `st_nlink`, `st_size`, and `st_mode` in octal.

---

## Hard link vs symbolic link

A hard link is a directory entry that points to an existing inode on the same filesystem. `ln file other` creates a second name. Both names are equal. There is no "original" in the inode. The link count increases. `unlink` of one name decreases the count.

Rules:

- you cannot hard-link a directory on Linux (except `.` and `..` that the filesystem manages)
- you cannot hard-link across filesystems
- a hard link works after you rename the other name

A symbolic link (symlink) is a file of type symlink. The data of that file is a path string. `ln -s target name` creates it. The kernel follows the string when you `open` the name (unless you use `lstat` or `O_NOFOLLOW`).

Rules:

- a symlink can point across filesystems
- a symlink can be dangling (the target does not exist)
- a symlink can form a cycle; the kernel stops after a hop limit (`ELOOP`)
- `rm` of a symlink removes the link, not the target

`ls -l` shows `->` for a symlink. `readlink` prints the stored path.

Use a hard link when you want two names for the same data on one filesystem. Use a symlink when you want a pointer that can cross mounts or that you can retarget.

Do not use a symlink when a confused program might follow a link that an attacker controls in `/tmp`. That is a security topic (symlink races). Prefer `openat` and `O_NOFOLLOW` in careful programs.

### Questions

#### Theoretical questions

1. What is a hard link?
2. What is a symbolic link?
3. Why can a hard link not cross filesystems?
4. What is a dangling symlink?
5. What error can a symlink cycle produce?

#### Easy practical tasks

1. Create a file `a`. Run `ln a b` and `ln -s a c`. Run `ls -li a b c`. Write inode numbers and the `->` line.
2. Open `man 1 ln` and `man 2 readlink`. Write the difference between `ln` and `ln -s`.
3. Write a two-column table: hard link vs symlink. Add four rows.
4. Run `readlink c` from the task above. Write the string.

#### Medium practical tasks

1. `unlink` `a`. Show that `b` still has the data and that `c` is dangling if it pointed to `a` by name.
2. Try `ln a /tmp/a2` if `a` is not on the same filesystem as `/tmp`. Write the error.
3. Draw both link types: two arrows to one inode, versus an arrow to a path string.

#### Advanced practical tasks

1. Create a symlink loop (`x -> y`, `y -> x`). Run `cat x`. Write the error. Remove the links.
2. Read `man 2 open` for `O_NOFOLLOW` and `man 2 unlink`. Write a safe delete-or-open note for `/tmp`.

---

## Permissions: user, group, other; `chmod`, `chown`

Unix permission bits sit in the inode. Three classes exist:

- user (owner)
- group
- other (everyone else)

Each class has three bits: read `r`, write `w`, execute `x`. `ls -l` shows them as `rwxr-xr-x`. The first character is the type (`-` file, `d` directory, `l` symlink).

For a regular file:

- `r` allows `open` for read
- `w` allows change of data
- `x` allows `exec` of the file as a program

For a directory:

- `r` allows listing of names
- `w` allows create and delete of names (with extra rules)
- `x` allows walk through the directory (search). Without `x` you cannot use a file inside even if you know the name, in the usual model

`chmod` changes the bits. You can use octal (`chmod 644 file`) or symbols (`chmod u+x file`).

`chown` changes the owner and optionally the group. Only root can change the owner on typical Linux. `chgrp` changes the group when you have the right.

The kernel also has extra bits: setuid, setgid, and the sticky bit. On a directory, the sticky bit (`/tmp`) limits who can delete names. This section requires you to know that they exist. Do not set setuid on your practice files.

Linux also has ACLs and capabilities. Topic 17 covers more isolation. The nine `rwx` bits are the base.

A process has a real user ID and an effective user ID. The kernel uses the effective IDs for the permission check. `sudo` changes the effective user.

### Questions

#### Theoretical questions

1. What are the three permission classes?
2. What does execute mean on a directory?
3. What does `chmod 644` mean for a regular file?
4. Who may `chown` a file to another user on typical Linux?
5. What does the sticky bit do on `/tmp`?

#### Easy practical tasks

1. Run `ls -l` on a file that you own. Write the ten-character mode string and explain each part.
2. Run `chmod 600` on a practice file. Run `ls -l`. Restore `644` or `664`.
3. Open `man 1 chmod` and `man 2 chmod`. Write one sentence for each.
4. Run `id`. Write your UID, GID, and group names.

#### Medium practical tasks

1. Create a directory `d` with `chmod 700`. Put a file inside. Use a second user if you have one, or explain what that user would see.
2. Compare `chmod u+x` on a script with `chmod +x`. Run the script as `./script`. Write why `x` is required.
3. Draw the nine bits as a table: user/group/other × r/w/x.

#### Advanced practical tasks

1. Read `man 7 inode` on setuid and `man 1 chmod` on the sticky bit. Write a one-page warning sheet.
2. Write a C program that `stat`s a file and prints whether the current process could read it (compare `st_mode` with `geteuid` and `getegid`). Handle the owner/group/other cases.

---

## File descriptors

A file descriptor (FD) is a small integer that refers to an open file description in the kernel. The process holds a table of FDs. The first three standard FDs are:

- `0` stdin
- `1` stdout
- `2` stderr

`open` returns the next free FD (the lowest number). `close` frees that slot. After `close`, that number can refer to a new open.

The FD is not the inode number. Two FDs can refer to the same inode. One FD can be duplicated with `dup` so that two numbers refer to the same open file description. Those two share the file offset.

The open file description is a kernel object. It holds the offset, the access mode, and a pointer to the inode (or more precisely to a file object). `fork` copies the FD table. Parent and child share the same open file descriptions for the FDs that were open at `fork`. A `read` in the child moves the offset that the parent also sees.

`/proc/self/fd` is a directory of symlinks from FD numbers to paths (when the kernel can show a path). `ls -l /proc/self/fd` is the practical view.

There is a per-process limit (`ulimit -n`) and a system limit. A leak of FDs (`open` without `close`) hits the limit.

Do not confuse FD `1` with inode `1`. They are different spaces.

### Questions

#### Theoretical questions

1. What is a file descriptor?
2. What are FDs 0, 1, and 2?
3. What does an open file description store that the FD number does not store?
4. How does `fork` affect FDs and offsets?
5. What happens if you never `close` FDs?

#### Easy practical tasks

1. Run `ls -l /proc/self/fd`. Write what `0`, `1`, and `2` point to.
2. Run `ulimit -n`. Write the soft limit.
3. Open `man 2 close` and `man 2 dup`. Write one sentence for each.
4. Write five sentences about FDs. Use only facts from this section.

#### Medium practical tasks

1. Write a program that `open`s a file and prints the FD number. Run `ls -l /proc/<pid>/fd` while the program sleeps.
2. Use `dup2` to redirect stdout to a file. Print `hello` and show the file contents.
3. Draw: process FD table, open file description (offset), inode.

#### Advanced practical tasks

1. After `fork`, let the child `read` one byte. Let the parent `read` the next byte from the same FD without a seek. Write why the parent does not see the first byte again.
2. Read `man 2 fcntl` for `F_DUPFD_CLOEXEC` and `FD_CLOEXEC`. Write why close-on-exec matters for `exec`.

---

## `open`, `read`, `write`, `close`, `lseek`

These five calls are the core POSIX file I/O interface.

`open` or `openat` turns a path into an FD. Flags include `O_RDONLY`, `O_WRONLY`, `O_RDWR`, `O_CREAT`, `O_TRUNC`, `O_APPEND`, and `O_CLOEXEC`. Always check for `-1` and read `errno`. When you pass `O_CREAT`, you also pass a mode (permission bits, modified by `umask`).

`read` copies up to `n` bytes from the open file into a user buffer. The return value is the number of bytes, `0` on end of file, or `-1` on error. A short read is allowed. A loop is required if you need exactly `n` bytes (except for some devices).

`write` copies bytes from a user buffer to the file. A short write is allowed. A loop is required for exact counts. `O_APPEND` makes each write go to the current end (atomic with respect to other `O_APPEND` writers on normal files, with documented limits).

`close` releases the FD. It does not always flush storage to the device. `fsync` is a later topic. `close` can still return an error (for example, a delayed write error). Check it in careful programs.

`lseek` changes the offset of the open file description. `SEEK_SET`, `SEEK_CUR`, and `SEEK_END` set the new position. `lseek` does not work on all FD types (pipes, some sockets). It returns the new offset or `-1`.

```c
int fd = open("notes.txt", O_RDONLY);
if (fd < 0)
	return 1;
char buf[64];
ssize_t n = read(fd, buf, sizeof buf);
close(fd);
```

`stdio` (`fopen`, `fread`) uses these calls inside the C library. Learn the system calls first. Then you understand buffers and `fflush`.

### Questions

#### Theoretical questions

1. What does `open` return on success and on failure?
2. Why must you loop on `read` if you need an exact byte count?
3. What does `O_APPEND` change?
4. Does `close` guarantee that data is on the disk?
5. When does `lseek` fail?

#### Easy practical tasks

1. Open `man 2 open`, `man 2 read`, `man 2 write`, `man 2 close`, and `man 2 lseek`. Write one sentence for each.
2. Type the sample. Print `n` and the bytes (add a NUL with care).
3. Run `strace -e openat,read,write,close,lseek cat /etc/hostname`. Write the call list.
4. Write a program that writes `hello\n` with `write(1, ...)`. Run it.

#### Medium practical tasks

1. Write a file copy: `open` source and dest, loop `read`/`write`, `close` both. Copy a small text file. Diff the result.
2. Use `lseek` to write a byte at offset 1_000_000 in a new file. Run `ls -l` and `du`. Write why the size and the disk use can differ (hole / sparse file).
3. Open with `O_APPEND` and write from two processes. Read the file. Write what you observe.

#### Advanced practical tasks

1. Implement a loop that handles `EINTR` and short writes. Document the retry rules from the man pages.
2. Compare `pread`/`pwrite` with `lseek` plus `read` in a multithreaded process. Write why positioned I/O avoids a shared offset race.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from the string `/home/ada/a.txt` to a `read` of one byte. Name directory walk, inode, FD, and offset.
2. How do a hard link, a symlink, and a second FD to the same file differ as "two ways to see data"?
3. Why do directory execute bits and file read bits both matter when you `open` a nested path?
4. What process state from topic 4 (CWD, FD table, umask) changes how these calls behave?
5. Which objects in this topic are per process, and which objects live in the filesystem until unlink?

#### Easy practical tasks

1. Write a one-page cheat sheet: file, directory, absolute/relative, inode, hard link, symlink, `rwx`, `chmod`, FD 0/1/2, `open`/`read`/`write`/`close`/`lseek`.
2. Run `stat`, `ls -li`, `id`, `ulimit -n`, and `ls /proc/self/fd` and comment each output.
3. Draw one figure: path, inode, two names, one FD table, one offset.
4. Open `man 7 path_resolution` and `man 2 open`. Write three flags that change path walk or create.

#### Medium practical tasks

1. Write a small tool that prints inode, mode, and link count for each argument, and that copies one file with the five system calls.
2. Create a directory tree with a symlink and a hard link. Document `stat` versus `lstat` and `readlink`.
3. Use `strace` on your copy tool. Count `read` and `write` calls for a 100 KiB file with a 4 KiB buffer.

#### Advanced practical tasks

1. Implement `openat`-based walk of a directory that you `open` once. Refuse symlinks with `O_NOFOLLOW` on each step. Document the security goal.
2. Read the OSTEP files and directories chapter. Write a one-page map from that chapter to each section in this handbook.
