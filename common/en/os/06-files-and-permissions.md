# 6. Files and Permissions

## Description

A filesystem is the OS abstraction for named, persistent data. This topic explains files, directories, paths, and the Unix inode. You learn hard links and symbolic links, permission bits, and file descriptors. You also learn the calls `open`, `read`, `write`, `close`, and `lseek`.

Complete this topic after processes and virtual memory. A file descriptor lives in the process. Complete this topic before filesystem layout on disk (topic 7).

Use one term for each concept. A file is a named inode with data. A directory is a file that maps names to inode numbers. A path is a string of names. A file descriptor is an integer in one process. A hard link is another name for the same inode. A symbolic link is a small file that stores a path. Do not mix a path with an inode. Do not mix a file descriptor with an inode number.

---

## File, directory, path, inode

A file is the object that the OS names and that a program reads or writes. On Unix, a file is not only a document. A directory, a device node, a socket, and a named pipe are files in the inode sense. The type is a field of the inode.

A directory is a file that stores a list of names. Each name points to an inode number. The kernel uses the directory to resolve a path.

A path is a sequence of names that the kernel walks. The names are separated by `/` on Unix. The path `/home/ada/notes.txt` means: start at the root directory, then `home`, then `ada`, then `notes.txt`.

An inode is the kernel and on-disk record of one file. The inode stores the type, the permission bits, the size, the timestamps, the link count, and the pointers to data. The kernel does not store the full path inside the inode. Many paths can reach the same inode (hard links).

An absolute path starts with `/`. The kernel starts the walk at the root of the process. A relative path does not start with `/`. The kernel starts the walk at the current working directory (CWD). `chdir` changes the CWD. A child inherits the CWD at `fork`. `exec` keeps the CWD.

`.` is the current directory in a path. `..` is the parent directory. The root `..` is the root.

`PATH` is not a filesystem path of a file. `PATH` is a list of directories that the shell searches for a command name.

User programs see a tree. The tree starts at `/`. A mount can attach another filesystem at a directory. Topic 7 covers mounts.

Do not think of a file as "bytes on one disk sector". The file is the inode plus the data blocks plus at least one name in a directory (except some unnamed objects).

### Questions

#### Theoretical questions

1. What is a file on Unix at a high level?
2. What does a directory store?
3. What is an inode?
4. Why does the inode not store the full path as its identity?
5. How does an absolute path differ from a relative path?

#### Easy practical tasks

1. Run `ls -l /`. Write five names and whether each looks like a directory (`d` in the first column).
2. Run `stat /` and `stat /tmp`. Write the inode number and the file type for each.
3. Write five sentences: file, directory, path, inode. Use only facts from this section.
4. Draw a tree: `/`, `home`, `ada`, `notes.txt`.

#### Medium practical tasks

1. Run `pwd` and `realpath .` if the tool exists. Write the absolute path of your CWD.
2. Use `cd` in a subshell versus the current shell. Write which CWD change the parent sees.
3. Open `man 7 inode` or `man 7 path_resolution`. Write the steps of a path walk in your own words.

#### Advanced practical tasks

1. Write a C program that prints the CWD with `getcwd` and then `chdir`s to `/tmp` and prints again. This handbook does not contain the source.
2. Write a one-page note: which objects in `/dev` are files, and why Unix calls them files.

---

## Hard link vs symbolic link

A hard link is a directory entry that names an inode. Two hard links to the same inode are two names for the same file. They share size, data, and permissions. `ln` without `-s` creates a hard link. `stat` shows the same inode number. The link count in the inode rises.

You cannot hard-link a directory on most Unix systems (or only `root` can in special cases). You cannot hard-link across filesystems. The inode number is unique per filesystem.

A symbolic link (symlink) is a small file that stores a path. `ln -s` creates it. `ls -l` shows `->`. When a program opens the path with the usual rules, the kernel walks the target. `readlink` reads the stored path. `lstat` gives facts about the symlink itself. `stat` follows the link by default.

A symlink can point to a missing target (a dangling link). A symlink can cross filesystems. A symlink can form a cycle. The kernel has a walk limit (`ELOOP`).

`unlink` or `rm` removes a name. For a hard link, the inode and data stay until the link count is 0 and no process has the file open. For a symlink, `rm` removes the symlink. The target stays.

Do not use a symlink as a lock. Do not follow untrusted symlinks in a privileged directory without care (`O_NOFOLLOW`).

### Questions

#### Theoretical questions

1. What is a hard link?
2. What is a symbolic link?
3. Why can a hard link not cross filesystems?
4. What does a dangling symlink mean?
5. When does `unlink` free the data of a regular file?

#### Easy practical tasks

1. Open `man 1 ln` and `man 2 readlink`. Write one sentence for each.
2. In a practice directory, create a file, a hard link, and a symlink. Run `ls -li`. Write which names share an inode.
3. Write five sentences that contrast hard links and symlinks. Use only facts from this section.
4. Run `readlink -f` on the symlink if the tool exists. Write the result.

#### Medium practical tasks

1. Remove the original name of a hard-linked file. Show that the second name still has the data.
2. Create a dangling symlink. Run `ls -l` and `cat` on it. Write the two results.
3. Draw an inode with link count 2 and two directory names. Draw a symlink inode that stores a path.

#### Advanced practical tasks

1. Read `man 2 open` for `O_NOFOLLOW`. Write when a privileged program must use it.
2. Write a one-page note: `rename` and replacing a file that has two hard links (what the other name sees).

---

## Permissions and file descriptors

Unix permission bits sit on the inode. `ls -l` shows a string such as `-rw-r--r--`. The first character is the type. The next nine characters are owner, group, and other, each with read, write, and execute.

The kernel checks the effective user ID and group IDs of the process (topic 11) against the inode owner and group. Execute on a directory means you may walk through it. Read on a directory means you may list names.

`chmod` changes the bits. `chown` changes owner. `umask` masks bits when a process creates a file. A typical `umask` of `022` makes files `644` and directories `755` if the create mode asked for `666` or `777`.

Extra bits exist: setuid, setgid, and the sticky bit. setuid on an executable can raise the effective UID (topic 11). The sticky bit on `/tmp` limits unlink of other users files.

A file descriptor (FD) is a small integer in one process. `0` is stdin, `1` is stdout, `2` is stderr. `open` returns the next free FD. The process FD table points at a kernel file object. That object has an offset and a status flags word. Two FDs in two processes can point at the same kernel file object after `fork`. Two `open` calls create two offsets.

`ls -l` is not the FD table. `ls -l` is the inode and the name. `lsof` or `/proc/self/fd` shows FDs.

Permissions are checked at `open` (and at some later operations). A process that already has an FD can `read` after `chmod` removes access for new opens. That is a standard Unix rule.

Do not run daily work as UID 0. Do not set files to `0777` to "make it work".

### Questions

#### Theoretical questions

1. What do the nine permission bits mean?
2. What does execute mean on a directory?
3. What is a file descriptor?
4. How does `fork` share file offsets?
5. When does the kernel check permission bits for a normal `read`?

#### Easy practical tasks

1. Run `ls -l` on a file that you own. Write the type and the nine bits.
2. Run `umask` and `ls -l /proc/self/fd`. Write the umask and three FD targets.
3. Open `man 2 chmod` and `man 2 open`. Write one sentence for each.
4. Write five sentences about permissions and FDs. Use only facts from this section.

#### Medium practical tasks

1. Create a file with `echo`. Run `chmod 600` and `chmod 000`. Try `cat` as the owner. Write the result.
2. Draw process, FD table, kernel file object, inode.
3. Run `exec 3>/tmp/os-fd-practice.txt` in bash, then `ls -l /proc/$$/fd`. Write what FD 3 points to. Close it with `exec 3>&-`.

#### Advanced practical tasks

1. Read `man 7 credentials` and `man 2 umask`. Write how umask and the mode argument of `open` combine.
2. Write a one-page note: why `sudo cat` of a secret file can leak the secret to an unprivileged stdout if you redirect without care.

---

## `open`, `read`, `write`, `close`, `lseek`

These five calls are the core Unix file API.

`open` takes a path and flags. `O_RDONLY`, `O_WRONLY`, and `O_RDWR` set the access. `O_CREAT` creates. `O_TRUNC` sets size to 0. `O_APPEND` forces writes to the end. `open` returns an FD or `-1` and `errno`. Always check the return.

`read` copies up to `n` bytes from the FD into a buffer. The return is the number of bytes, `0` on end-of-file, or `-1` on error. A short read is normal. A pipe and a socket often return less than `n`. A disk file can also return less. Loop until you have enough bytes or until EOF.

`write` copies bytes from a buffer to the FD. A short write is also possible. Loop. `write` to a full pipe can block. `write` does not always hit the device. The page cache can hold the data (topic 7). Use `fsync` when you need durability.

`close` releases the FD. The kernel file object lives until the last FD that points at it is closed. After `close`, that integer may be reused. Do not use a closed FD.

`lseek` sets the file offset for FDs that support it. `SEEK_SET`, `SEEK_CUR`, and `SEEK_END` are the whences. Pipes, sockets, and many devices do not support `lseek` (`ESPIPE`). `O_APPEND` writes still go to the end even if you seek.

`pread` and `pwrite` read or write at an offset without changing the FD offset. They help when threads share an FD.

Do not ignore `EINTR`. A signal can interrupt `read` or `write`. Retry or handle the error.

### Questions

#### Theoretical questions

1. What does `open` return on success and on failure?
2. Why must you loop on `read` and `write`?
3. What does `close` release?
4. When does `lseek` fail with `ESPIPE`?
5. How do `pread` and `lseek` plus `read` differ for a shared FD?

#### Easy practical tasks

1. Open `man 2 open`, `man 2 read`, `man 2 write`, `man 2 close`, and `man 2 lseek`. Write one sentence for each.
2. Write the flags `O_RDONLY` and `O_CREAT` in a two-row table with one use each.
3. Write five sentences about these five calls. Use only facts from this section.
4. Run `strace -e openat,read,write,close,lseek cat /etc/hostname` (or a small file). Write the call list.

#### Medium practical tasks

1. Write a C program that copies a file with a 4 KiB buffer: `open`, loop `read`/`write`, `close`. This handbook does not contain the source.
2. Use `lseek` to print the last 16 bytes of a file. Handle a short file.
3. Try `lseek` on a pipe (`pipe` then `lseek`). Write the `errno` name.

#### Advanced practical tasks

1. Add `O_APPEND` and two writers. Show that `lseek` to 0 does not stop append. Write what you saw.
2. Read `man 2 posix_fadvise` or `man 2 fsync`. Write when you call `fsync` after `write` (preview of topic 7).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do path, directory entry, inode, and file descriptor form one open story?
2. When do you use a hard link, and when do you use a symlink?
3. Why can two FDs have two offsets for the same inode?
4. A classmate says "permissions apply on every `read`." Which Unix rule do you use to correct that sentence?
5. Why is a short `write` not always an error?

#### Easy practical tasks

1. Write a cheat sheet: file, directory, path, inode, hard link, symlink, `chmod`, FD 0/1/2, `open`, `read`, `write`, `close`, `lseek`.
2. Run `ls -li` and `stat` on one practice file. Write inode, links, and mode.
3. Draw `open` of a relative path: CWD, directory walk, inode, new FD.
4. Open `man 7 path_resolution` if it exists. Write one fact about symlink walks.

#### Medium practical tasks

1. Write a program that prints the inode of a file with `fstat` after `open`. Compare with `stat` on the path.
2. Use `strace -c` on your copy program. Write the count of `read` and `write`.
3. Document a permission debug checklist: `ls -l`, `id`, directory execute bits, symlink target, `umask`.

#### Advanced practical tasks

1. Read [The Linux Programming Interface](https://man7.org/tlpi/) chapter list on files. Write a one-page map from this topic to those chapters.
2. Design a safe create-and-replace of a config file: write a temp file in the same directory, `fsync`, `rename`. Write why `rename` on the same filesystem is atomic for the name.
