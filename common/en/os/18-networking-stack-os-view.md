# 18. Networking Stack (OS View)

## Description

The kernel implements sockets, TCP and UDP, routing, and packet filters. This topic is the operating-system view of that stack. You learn the sockets API, how the kernel treats TCP and UDP, network namespaces, and awareness of iptables and nftables. The full networking path lives in `net.topics.md`. This topic does not replace that path.

Complete this topic after security and isolation. Complete this topic before persistence of OS state (topic 19). You already know file descriptors, `epoll`, Unix sockets, and namespaces. You now add internet sockets.

Use one term for each concept. A socket is an endpoint that a process opens with `socket`. TCP is a reliable byte-stream protocol that the kernel implements. UDP is a datagram protocol that the kernel implements. A network namespace is a copy of the network stack view (interfaces, routes, filter tables). A packet filter is kernel code that accepts, drops, or changes packets by rules. Do not mix a Unix socket with an internet socket. Do not mix iptables with the sockets API.

---

## Sockets API

A process talks to the network through the sockets API. The core calls are POSIX and exist on Linux, BSD, and macOS. Windows has a similar model with different names.

Typical stream-server path:

1. `socket(AF_INET, SOCK_STREAM, 0)` or `AF_INET6`
2. `bind` to an address and port
3. `listen`
4. `accept` to get a new connected descriptor
5. `read` / `write` or `recv` / `send`
6. `close`

Typical stream-client path: `socket`, optional `bind`, `connect`, then I/O, then `close`.

UDP uses `SOCK_DGRAM`. There is no `listen` or `accept`. `sendto` and `recvfrom` include addresses.

The socket is a file descriptor on Unix. `epoll` can wait on it (topic 14). `fork` inherits it. `SCM_RIGHTS` can pass it (topic 15).

Address structures (`sockaddr_in`, `sockaddr_in6`) carry IP and port. `getaddrinfo` converts a name and a service string to those structures. Prefer `getaddrinfo` over old `gethostbyname`.

`SO_REUSEADDR` and (on Linux) `SO_REUSEPORT` change bind rules. Read the man page before you copy flags.

Errors are in `errno`. A TCP `read` of 0 means the peer closed. `EAGAIN` means non-blocking and not ready.

This section is the OS interface. Packet formats, routing, and DNS details belong to `net.topics.md` (especially topics 8, 9, and 10).

### Questions

#### Theoretical questions

1. What does `socket` create?
2. What is the server call order for TCP?
3. How does UDP I/O differ from TCP I/O at the call level?
4. Why is a socket a file descriptor on Unix?
5. What does `getaddrinfo` return?

#### Easy practical tasks

1. Open `man 2 socket`, `man 2 bind`, and `man 2 accept`. Write one sentence for each.
2. Open `man 3 getaddrinfo`. Write why it replaces `gethostbyname`.
3. Write five sentences about the sockets API. Use only facts from this section.
4. Draw client `connect` and server `accept` with two file descriptors on the server.

#### Medium practical tasks

1. Write a TCP echo server and client on `127.0.0.1` at a high port. Send one line.
2. Use `ss -lntp` or `netstat -lntp` while the server listens. Write the local address.
3. Trace your client with `strace -e socket,connect,sendto,recvfrom,close`. Write the call list.

#### Advanced practical tasks

1. Add `epoll` (topic 14) to the server so that two clients can stay connected. Document edge cases.
2. Read `man 7 ip` and `man 7 ipv6`. Write a one-page note on `AF_INET` versus `AF_INET6` dual bind (`IPV6_V6ONLY`).

---

## TCP/UDP from the kernel side

The kernel owns protocol state. User space sees descriptors. Inside, TCP has a state machine (LISTEN, SYN-SENT, ESTABLISHED, TIME-WAIT, and the other textbook states). Each connection has a control block: sequence numbers, windows, buffers, and timers.

A TCP `write` copies bytes into a kernel send buffer (or waits). The stack segments the bytes, adds headers, and passes packets to the IP layer and the driver (topic 14). ACKs from the peer free the buffer. `read` copies from the receive buffer.

UDP has no connection state in the TCP sense. The kernel keeps a socket receive queue of datagrams. A datagram can drop if the queue is full. UDP does not retransmit.

The kernel also implements:

- routing lookup (which interface, which next hop)
- neighbor lookup (ARP or NDP)
- fragmentation and reassembly (where the stack still does it)
- checksums (sometimes offloaded to the NIC)

`/proc/net/tcp` and `ss` show connection tuples and states. `netstat` is the older tool.

A process cannot set TCP sequence numbers through the POSIX API. That is kernel state. Socket options (`setsockopt`) expose a small, safe subset: reuse flags, buffer sizes, `TCP_NODELAY`, keepalive.

Kernel TCP is not the only implementation. Some systems use user-space stacks. This path studies the Linux kernel stack.

TIME-WAIT exists so that a delayed packet does not join a new connection that reused the same tuple. Students who bind the same port quickly meet this state.

### Questions

#### Theoretical questions

1. Where does a TCP `write` put bytes before they are on the wire?
2. Why can a UDP `send` succeed and the datagram still be lost?
3. What is a TCP control block at a high level?
4. Why is TIME-WAIT not a user-program bug by itself?
5. What does `TCP_NODELAY` change?

#### Easy practical tasks

1. Open `man 7 tcp` and `man 7 udp`. Write one sentence for each.
2. Run `ss -s` or `ss -ta`. Write how many TCP sockets you see.
3. Write five sentences about the kernel view of TCP and UDP. Use only facts from this section.
4. Make a table: TCP versus UDP. Add "kernel state" and "reliability".

#### Medium practical tasks

1. Run `ss -tan` and find a TIME-WAIT or ESTAB row. Write the columns that you understand.
2. Draw send buffer, NIC, peer, ACK, then buffer free.
3. Read `man 7 socket` for `SO_RCVBUF`. Write six sentences on why a huge user `read` still needs kernel buffers.

#### Advanced practical tasks

1. Read a Linux TCP overview (kernel docs or LWN). Write a one-page map of states that `ss` prints.
2. Compare `TCP_NODELAY` and Nagle in a small experiment: many tiny `write`s with and without the option. Write packet counts from `tcpdump` on lo. Use a VM.

---

## Network namespaces

A network namespace is a Linux namespace (topic 16) that holds a private set of:

- network interfaces
- IPv4 and IPv6 addresses
- routing tables
- packet filter tables
- `/proc/net` view

A process in a new netns does not see the host `eth0` unless you move a device or create a pair.

`ip netns add` creates a named namespace. `unshare --net` creates an anonymous one for a command. Container runtimes give each container a netns.

A common pattern is a veth pair: one end in the host (or a bridge), the other end in the container netns. The host routes or bridges traffic. NAT on the host can hide many containers behind one address (see `net.topics.md` on NAT).

`lo` exists in each netns. A listen on `127.0.0.1` in a container is not the host loopback.

`/proc/self/ns/net` is the namespace handle. `nsenter --net` joins it.

Network namespaces isolate the stack. They do not encrypt traffic. They do not replace firewall policy on the host.

Do not delete host interfaces. Practice with `ip netns` in a VM.

### Questions

#### Theoretical questions

1. What objects does a network namespace isolate?
2. Why is `127.0.0.1` in a container not the host loopback?
3. What is a veth pair for?
4. How do you create a netns with `unshare`?
5. Why is a netns not encryption?

#### Easy practical tasks

1. Open `man 8 ip-netns` and `man 7 network_namespaces`. Write one sentence for each.
2. Run `ls -l /proc/self/ns/net`. Write that the link is a namespace inode.
3. Run `ip link` or `ip addr`. Write two interface names.
4. Write five sentences about network namespaces. Use only facts from this section.

#### Medium practical tasks

1. In a VM, run `sudo ip netns add os-practice` and `sudo ip netns exec os-practice ip link`. Write which interfaces exist. Then delete the netns.
2. Draw host netns, container netns, and a veth pair.
3. Compare `ss` on the host with `ss` inside `unshare --net` (as root in a VM). Write the difference.

#### Advanced practical tasks

1. Create two netns, a veth pair, addresses, and a ping between them. Document every `ip` command. Use a VM.
2. Write a one-page map from Docker `--network bridge` to netns, veth, and a host bridge.

---

## iptables / nftables (awareness)

Linux can filter and transform packets in the kernel. Two user interfaces exist.

iptables is the older tool. It talks to the `xtables` and netfilter hooks. Rules sit in tables (`filter`, `nat`, `mangle`) and chains (`INPUT`, `FORWARD`, `OUTPUT`, and others). Each rule matches a packet and then accepts, drops, or jumps.

nftables is the newer interface. One `nft` tool and a single rule VM replace many `iptables` binaries. Distributions move to nftables. `iptables-nft` can still accept old syntax.

Netfilter is the kernel framework. Both tools program it (with different backends). Hooks sit in the receive, forward, and transmit paths.

Awareness goals:

- Know that a silent drop can be a filter, not an application bug.
- Know `iptables -L` versus `nft list ruleset` (rights needed).
- Know that containers and Kubernetes add more rules. Do not edit them at random.
- Know that a filter is not the sockets API. A process can listen and still receive nothing.

NAT rules rewrite addresses. That topic belongs to `net.topics.md`. Here you only need: the kernel can change a packet before a socket sees it.

Do not flush a production ruleset. Do not practice `iptables -F` on a host that you need.

### Questions

#### Theoretical questions

1. What is netfilter?
2. How does iptables differ from nftables at a high level?
3. What is a silent drop?
4. Why can a process listen and still see no packets?
5. Why must you not flush filter tables on a shared host?

#### Easy practical tasks

1. Open `man 8 iptables` and `man 8 nft`. Write one sentence for each.
2. Run `command -v iptables nft`. Write which tools exist.
3. Write five sentences about packet filters. Use only facts from this section.
4. Make a table: sockets API versus netfilter. Add "who it is for".

#### Medium practical tasks

1. In a VM, run `sudo nft list ruleset` or `sudo iptables -L -n`. Write whether any rules exist. Do not change them unless you own the VM.
2. Draw packet path: NIC, netfilter hook, socket receive queue, `read`.
3. Read a short nftables getting-started page. Write six sentences on table versus chain.

#### Advanced practical tasks

1. In a disposable VM, add one nftables rule that drops a high UDP port, test with `nc`, then delete the rule. Document the rollback first.
2. Write a one-page note: how Docker iptables/nft rules can hide a bind to `0.0.0.0`. No attack recipe.

---

## Full networking path is in `net.topics.md`

This OS topic stops at the kernel interface and the isolation knobs. `net.topics.md` is the networking learning path. Use it for layers, addresses, routing, TCP details, DNS, TLS, and debugging with `tcpdump`.

Pairing:

- OS topic 14 (`epoll`) plus net topic 10 (sockets programming)
- OS topic 16–17 (namespaces, capabilities) plus net topic 21 (hardening)
- OS topic 18 (this file) plus net topics 8 and 9 (UDP and TCP)

Do not skip `net.topics.md` if you write servers. The sockets API is small. Production failures are often routing, MTU, NAT, or TLS.

Do not duplicate that path here. When a later OS lab needs a port or a packet capture, open the matching net topic.

Suggested order after this file: finish OS topics 19–22, or start `net.topics.md` from topic 1 if you want packets first. Both orders work. If you start net first, return here for namespaces and netfilter.

Official starting points: Linux `man 7 ip`, `man 7 tcp`, and the resources list in `net.topics.md`.

### Questions

#### Theoretical questions

1. What does this OS topic cover that `net.topics.md` does not cover in depth?
2. What does `net.topics.md` cover that this topic leaves out?
3. Which net topics pair with `epoll` and with TCP state?
4. Why is the sockets API not enough to debug a production failure?
5. When can you start `net.topics.md` relative to this file?

#### Easy practical tasks

1. Open `net.topics.md` in this repository. Write the titles of topics 8, 9, and 10.
2. Write five sentences about how the two paths split the work. Use only facts from this section.
3. Bookmark `man 7 ip` and `man 7 tcp`.
4. Make a two-column table: "OS handbook" and "net handbook". Add three rows.

#### Medium practical tasks

1. Write a study plan of eight lines: four OS labs and four net labs that share sockets.
2. Run `ip route` and `ip addr`. Write one fact that belongs to net topic 6 (routing) more than to this file.
3. List three tools from `net.topics.md` practice (`ping`, `dig`, `curl`, `tcpdump`) and one OS question each tool can still answer.

#### Advanced practical tasks

1. Read net topic 10 in the topics file and this handbook's sockets section. Write a one-page map of overlapping terms.
2. After you finish net topics 1–6, return here and add a short note (your own file) on what a network namespace changes in those topics.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of one TCP byte from `write` in a process to a driver DMA (topic 14). Name the send buffer and one netfilter hook.
2. How do a Unix socket (topic 15) and an `AF_INET` socket share the API and differ in the kernel stack?
3. A container listens on port 8080. Which namespace, filter, and host NAT facts can hide that port from your laptop?
4. Why does `CAP_NET_ADMIN` (topic 17) matter for `ip` and `nft` but not for a normal `connect`?
5. When do you debug with `ss` and when do you debug with `tcpdump` (net path)?

#### Easy practical tasks

1. Write a one-page cheat sheet: `socket`/`bind`/`listen`/`accept`/`connect`, TCP buffer, UDP queue, netns, veth, iptables, nftables, `net.topics.md`.
2. Run `ss -s`, `ip -br addr`, `ls /proc/self/ns/net`, and `command -v nft iptables`. Comment each.
3. Draw one figure: process, socket, TCP control block, netns, NIC.
4. Bookmark this file, `net.topics.md`, `man 7 tcp`, and `man 8 nft`.

#### Medium practical tasks

1. Write a TCP server in a netns and a client in the host netns connected by veth (VM). Document `ss` in both namespaces.
2. Use `strace` on `curl 127.0.0.1` if a local server exists, or on `curl example.com`. Write socket-related calls only.
3. Map topic 14 readiness I/O onto a TCP listener and ten idle clients.

#### Advanced practical tasks

1. Read OSTEP or kernel networking documentation at a high level. Write a one-page map to this handbook and a pointer to `net.topics.md` for the rest.
2. In a VM, compare `TCP_NODELAY` and a packet filter drop of the peer port. Write how the application errors differ.
