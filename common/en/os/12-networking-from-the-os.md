# 12. Networking from the OS

## Description

The kernel implements sockets, TCP and UDP, routing, and packet filters. This topic is the operating-system view of that stack. You learn the sockets API, how the kernel treats TCP and UDP, network namespaces, and awareness of iptables and nftables. The full protocol path lives in `net.topics.md`. This topic does not replace that path.

Complete this topic after security and isolation. Complete this topic before observability (topic 13). You already know file descriptors, `epoll`, Unix sockets, and namespaces.

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

The socket is a file descriptor on Unix. `epoll` can wait on it (topic 8). `fork` inherits it. `SCM_RIGHTS` can pass it (topic 9).

Address structures (`sockaddr_in`, `sockaddr_in6`) carry IP and port. `getaddrinfo` converts a name and a service string to those structures. Prefer `getaddrinfo` over old `gethostbyname`.

`SO_REUSEADDR` and (on Linux) `SO_REUSEPORT` change bind rules. Read the man page before you copy flags.

Errors are in `errno`. A TCP `read` of 0 means the peer closed. `EAGAIN` means non-blocking and not ready.

This section is the OS interface. Packet formats, routing, and DNS details belong to `net.topics.md`.

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

1. Write a TCP echo server and client on `127.0.0.1` at a high port. Send one line. This handbook does not contain the source.
2. Use `ss -lntp` or `netstat -lntp` while the server listens. Write the local address.
3. Trace your client with `strace -e socket,connect,sendto,recvfrom,close`. Write the call list.

#### Advanced practical tasks

1. Add `epoll` (topic 8) to the server so that two clients can stay connected. Document edge cases.
2. Read `man 7 ip` and `man 7 ipv6`. Write a one-page note on `AF_INET` versus `AF_INET6` dual bind (`IPV6_V6ONLY`).

---

## TCP/UDP in the kernel

The kernel owns the TCP and UDP state machines. Your process does not send raw segments if it uses the sockets API. The kernel builds headers, computes checksums (or offloads them), and queues packets on a device.

TCP in the kernel:

- connection table keyed by local and remote address and port
- sequence numbers, acknowledgements, retransmission timers
- a send buffer and a receive buffer (`SO_SNDBUF`, `SO_RCVBUF`)
- congestion control (the algorithm is a kernel module or a sysctl choice)
- `TIME_WAIT` after an active close

A `write` on a TCP socket copies bytes into the send buffer (or blocks). The kernel segments later. A `read` copies from the receive buffer. Topic 8 blocking rules apply.

UDP in the kernel:

- no connection state in the TCP sense (connected UDP sockets still only cache a peer)
- datagrams; a `recv` is one datagram (truncation is possible)
- no retransmission
- checksums; the kernel or the NIC can fail a bad packet

The kernel routing table picks the output interface. The neighbor table (ARP or NDP) finds the next MAC. Drivers and DMA (topic 8) send the frame. Incoming packets walk the reverse path: driver, IP, TCP/UDP, socket queue, wake the process.

`netstat` and `ss` show sockets and some TCP states. `/proc/net/tcp` and `/proc/net/udp` are older text tables.

You do not implement TCP in this topic. You learn that the state lives in the kernel and that `ss` shows it. `net.topics.md` covers handshake, window, and congestion in protocol detail.

### Questions

#### Theoretical questions

1. Where does the TCP state machine run when you use sockets?
2. What happens to bytes in `write` on a TCP socket before they hit the wire?
3. What TCP state can remain after `close`?
4. How does a UDP `recv` differ from a TCP `read` for message boundaries?
5. What kernel tables sit under an outgoing packet after TCP or UDP?

#### Easy practical tasks

1. Open `man 7 tcp` and `man 7 udp`. Write one sentence for each.
2. Run `ss -tuln`. Write one listening socket or write that you see none besides defaults.
3. Write five sentences about TCP and UDP in the kernel. Use only facts from this section.
4. Make a table: send buffer versus receive buffer. Add which call fills each.

#### Medium practical tasks

1. Start your echo server. Run `ss -ti` on the connection while you send data. Write one TCP info field if it appears.
2. Draw process `write`, TCP send buffer, IP, NIC, peer.
3. Read `man 7 socket` on `SO_RCVBUF`. Write six sentences on why a full receive buffer can stall the sender.

#### Advanced practical tasks

1. Compare `cat /proc/net/tcp` with `ss -t`. Write how a hex address maps (high-level). You may use a converter; write the steps.
2. Write a one-page note: kernel congestion control as a sysctl (`tcp_congestion_control`). Point to `net.topics.md` for AIMD.

---

## Network namespaces

A network namespace is a Linux namespace (topic 10) for the network stack. Each namespace has its own:

- network interfaces (except that physical devices sit in one namespace at a time)
- IPv4 and IPv6 addresses
- routing table
- iptables or nftables tables
- sockets (a listen port in namespace A is not the port in namespace B)

The initial namespace is the host namespace. `ip netns` creates a named namespace. A container runtime creates a network namespace per container by default.

A `veth` pair is a common link: one end in the container namespace, one end in the host or in a bridge namespace. The host NATs or routes to the outside. Topic 10 named the namespace. This section names the stack copy.

`ip netns exec <name> ss -l` shows sockets inside that namespace. `nsenter -n` is another way in.

Loopback `lo` exists in each namespace. You must bring it up (`ip link set lo up`) or `127.0.0.1` does not work inside.

Network namespaces isolate views. They do not replace a firewall policy by themselves. A container with `--network host` shares the host namespace. That mode is weaker isolation.

Do not delete host interfaces in a lab on a shared machine. Use a VM.

### Questions

#### Theoretical questions

1. What objects belong to a network namespace?
2. Why can two namespaces both listen on port 80?
3. What is a `veth` pair for?
4. Why must you bring up `lo` in a new namespace?
5. What does host networking do to namespace isolation?

#### Easy practical tasks

1. Open `man 8 ip-netns` and `man 7 network_namespaces` if they exist. Write one sentence for each.
2. Run `ls /var/run/netns` or `ip netns list`. Write whether any named namespaces exist.
3. Write five sentences about network namespaces. Use only facts from this section.
4. Run `ip addr` and write how many interfaces you see in the current namespace.

#### Medium practical tasks

1. In a VM, create a namespace with `ip netns add os-net-lab`, run `ip netns exec os-net-lab ip addr`, then delete it. Write what you saw on `lo`.
2. Draw host namespace, `veth`, container namespace, `lo` in each.
3. Read a Docker bridge overview. Write six sentences: how a container IP reaches the host.

#### Advanced practical tasks

1. Create two namespaces and a `veth` pair. Assign addresses and `ping` across. Document every command. Use a VM.
2. Write a one-page note: network namespace plus user namespace (rootless networking limits).

---

## iptables / nftables (awareness)

Linux packet filters sit in the kernel netfilter framework. Administrators write rules that match packets and then accept, drop, reject, or change them (NAT).

`iptables` is the older user tool and the older table names (`filter`, `nat`, `mangle`). Many scripts still use it. `iptables` can be a frontend over nftables on current distributions.

`nftables` is the current framework. The user tool is `nft`. One ruleset can replace several iptables tables. Syntax is different. Concepts stay: hooks (input, output, forward), matches, and verdicts.

Awareness rules for this path:

- Filters run in the kernel. A process `connect` can fail because a rule drops the packet. `errno` can be a timeout, not "iptables said no" in the socket API.
- Container and cloud systems add many rules. Do not flush tables on a shared host (`iptables -F`).
- `ss` shows sockets. `iptables -L` or `nft list ruleset` shows filters (needs rights).
- NAT (source NAT, port forward) is a filter-table job, not a sockets-API job. `net.topics.md` covers NAT in the protocol path.

You do not write a production firewall here. You learn that the OS has a packet filter and that it is not the sockets API.

`nft` and `iptables` need `CAP_NET_ADMIN` or root. Least privilege (topic 11) applies.

### Questions

#### Theoretical questions

1. What is netfilter in one sentence?
2. How does iptables relate to nftables on a current Linux?
3. Why can a `connect` fail without a sockets-API error that names the firewall?
4. What is one job of the `nat` table or an nft NAT chain?
5. Why must you not flush filter tables on a shared host?

#### Easy practical tasks

1. Open `man 8 iptables` and `man 8 nft` if they exist. Write one sentence for each.
2. Write five sentences about iptables and nftables. Use only facts from this section.
3. Make a table: sockets API versus packet filter. Add "who writes it" and "when it runs".
4. Run `ls /proc/net` | `head`. Write that filter state lives in the kernel (you may see nf files).

#### Medium practical tasks

1. If you have rights in a VM, run `nft list ruleset` or `iptables -L -n`. Write one chain name. Do not change rules.
2. Draw packet: NIC, filter hook, routing, TCP, socket.
3. Read a short nftables wiki "first ruleset" page. Write six sentences. Do not apply it on a shared host.

#### Advanced practical tasks

1. Write a one-page awareness note: conntrack and why NAT needs state. Point to `net.topics.md`.
2. Compare a distribution `firewalld` or `ufw` overview with raw `nft`. Write who generates the ruleset.

---

## Full protocol path is in `net.topics.md`

This OS topic stops at the kernel interface. The networking learning path in `net.topics.md` covers the protocol story in order:

- layers and encapsulation
- Ethernet, MAC, ARP
- IPv4 and IPv6 addresses
- routing, ICMP, NAT
- UDP and TCP behavior (handshake, window, congestion)
- HTTP and TLS
- DNS and more

Use this OS topic when you write `socket` and when you debug `ss`, namespaces, and filters. Use `net.topics.md` when you ask why a checksum exists, how a three-way handshake looks on the wire, or how a prefix is routed.

Pair the paths:

- OS sockets API + net TCP chapter
- OS network namespaces + net addressing and routing
- OS iptables awareness + net NAT and middleboxes
- OS `strace` on `connect` + net packet capture

Do not expect this file to teach CIDR or TLS. Do not expect `net.topics.md` to teach `task_struct` or page tables.

A lab that needs both: write a TCP client (this topic), capture the handshake with `tcpdump` (`net.topics.md`).

### Questions

#### Theoretical questions

1. What does this OS topic cover that `net.topics.md` does not replace?
2. What does `net.topics.md` cover that this file does not teach?
3. Which OS objects match the net chapters on TCP and UDP?
4. Which OS objects match the net chapters on NAT?
5. Why must you pair `connect` with a packet capture to see the handshake?

#### Easy practical tasks

1. Open `net.topics.md` in the repo. Write the titles of the first four topics.
2. Write five sentences about the split between OS and net paths. Use only facts from this section.
3. Make a two-column table: "Learn in OS" and "Learn in net". Add five rows.
4. Bookmark `man 7 ip` and the `net.topics.md` TCP section when you reach it.

#### Medium practical tasks

1. Write a study plan: three OS labs and three net labs that you will run in the same week.
2. Draw two books: kernel socket and on-the-wire segment. Label which file teaches each.
3. Run `man 7 netdevice` if it exists. Write one fact that is OS-side, not protocol-side.

#### Advanced practical tasks

1. Write a one-page map from this topic's five sections to specific `net.topics.md` headings.
2. Design a combined exam question: one sockets-API part and one packet-format part. Do not write the full answer key.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the sockets API, kernel TCP buffers, and a NIC driver form one `write`?
2. When do two processes need a Unix socket (topic 9) instead of `AF_INET`?
3. How does a network namespace change what `ss` and iptables show?
4. Why is a packet filter not a substitute for TLS or for application auth?
5. Which questions send you to `net.topics.md` instead of another `man 2` page?

#### Easy practical tasks

1. Write a cheat sheet: `socket`, `bind`, `listen`, `accept`, `connect`, TCP buffer, UDP datagram, netns, veth, iptables, nftables, `net.topics.md`.
2. Run `ss -tuln` and `ip addr`. Write one listening port and one address.
3. Draw client, kernel TCP, IP, filter, NIC.
4. Open `man 7 socket` and `net.topics.md`. Write one sentence about when you open each.

#### Medium practical tasks

1. Write a UDP ping-pong on localhost (two programs). Compare `ss -ul` with your TCP echo. This handbook does not contain the source.
2. Use `strace -e socket,bind,listen,accept,connect` on the TCP pair. Write the order.
3. Document a debug checklist: `ss`, `ip addr`, namespace, filter rules, then `tcpdump` (net path).

#### Advanced practical tasks

1. Combine `epoll` (topic 8) and TCP (this topic) in a small multi-client echo. Write a six-line design of the event loop.
2. In a VM, put a server in a network namespace and connect from the host over veth. Write the addresses. Point leftover protocol questions to `net.topics.md`.
