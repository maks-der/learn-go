# 8. Transport: UDP

## Description

UDP is a transport protocol that sends datagrams. UDP does not build a connection. UDP does not retransmit. This topic shows ports, datagrams, checksums, the lack of reliability and order, and when UDP is the right choice.

Use one term for each concept. A datagram is one UDP message. A port selects a socket on a host. Complete this topic before you study TCP.

UDP is simple. The application must handle loss, delay, and duplicate datagrams when those events matter.

---

## Port numbers

A UDP port is a 16-bit number. The UDP header has a source port and a destination port. The destination port selects the receiving socket. The source port tells the receiver where to send a reply.

The pair of IP addresses, UDP, and the two ports is the 4-tuple that identifies a UDP flow in many firewalls and NAT boxes. (TCP uses a 5-tuple with the same idea plus the protocol field already set to TCP.)

Well-known ports are the low numbers. DNS often uses UDP port 53. DHCP uses 67 and 68. Many games and voice tools use high ports that the vendor documents.

The client often uses an ephemeral source port. The server binds a known destination port. A server can also bind port 0 and let the OS pick a port. Then the client must learn that port by some other method.

Ports are per protocol. TCP port 53 and UDP port 53 are different sockets. A process can listen on both. Do not say "port 53" without TCP or UDP when the difference matters.

A firewall rule must name UDP when the service is UDP. A TCP allow does not pass UDP.

### Questions

#### Theoretical questions

1. How many bits does a UDP port use?
2. What does the destination port select?
3. Why does a client need a source port?
4. Are TCP port 53 and UDP port 53 the same socket?
5. What is an ephemeral port in UDP?

#### Easy practical tasks

1. Write five sentences that explain source port and destination port for a DNS query.
2. Make a table: service, UDP port. Add DNS and DHCP.
3. Run `ss -lnu` or `netstat -an` and find a UDP listen line if one exists. Write it.
4. Draw a client ephemeral port and server port 53.

#### Medium practical tasks

1. Use `dig` or `nslookup` and a capture. Write the client UDP port and the server port.
2. Write six sentences: a firewall that allows TCP 53 only. What happens to normal DNS.
3. Two processes cannot bind the same UDP address and port (usual case). Try on a lab host if you can. Record the error.

#### Advanced practical tasks

1. Write a one-page note: NAT mappings for UDP ports and why the source port on the WAN can differ.
2. Capture two DNS queries. Show that the client port can change. Write both 4-tuples.

---

## Connectionless datagram

UDP is connectionless. There is no handshake before the first datagram. The sender writes a datagram. The network tries to deliver that datagram. The receiver may get it, get a duplicate, or get nothing.

Each `send` is one datagram. Each `recv` returns at most one datagram (when you use a datagram API). The OS does not glue two UDP messages into one stream.

The UDP length field covers the UDP header and the payload. The maximum payload is limited by the IP MTU minus IP and UDP headers, unless the sender allows IP fragmentation. Large UDP datagrams can fragment. Fragmentation is fragile. Many applications keep UDP payloads small.

There is no UDP "listen" state like TCP listen. A bound UDP socket can receive datagrams from any peer, unless the program connects the UDP socket (a filter of the default peer). `connect` on UDP does not create a network handshake. It only sets a default remote address in the kernel.

Do not expect a "connection reset" as a required UDP feature. An ICMP port unreachable can arrive if the remote port is closed. Many hosts do not send it. Many firewalls drop it.

Do not treat a successful `send` as delivery. The call only queues the datagram.

### Questions

#### Theoretical questions

1. What does connectionless mean for UDP?
2. Does UDP glue two messages into a stream?
3. What does `connect` on a UDP socket do?
4. Why are large UDP datagrams risky?
5. Does a successful `send` prove delivery?

#### Easy practical tasks

1. Write five sentences that compare a UDP datagram and a TCP byte stream (preview).
2. Draw two datagrams from A to B. Show that they are separate messages.
3. Make a table: TCP handshake, UDP first packet. Add one row for "bytes on the wire before data".
4. Write four sentences about ICMP port unreachable as optional feedback.

#### Medium practical tasks

1. Write a tiny UDP sender and receiver in any language (or use `nc -u`). Send one short message. Confirm the receiver prints one message.
2. Send two messages. Confirm the receiver can tell them apart as two reads.
3. Write six sentences: UDP `connect` as a filter, not a handshake.

#### Advanced practical tasks

1. Send a UDP payload larger than the path MTU on a lab that you own. Capture fragments or a failure. Write what you see.
2. Write a one-page API note: `sendto` versus `send` after `connect` for UDP.

---

## Checksum

The UDP header has a checksum field. The checksum covers the UDP header, the payload, and a pseudo-header with IP addresses. The receiver drops a datagram with a bad checksum.

IPv4 allows a UDP checksum of zero to mean "no checksum" in old rules. Do not rely on that. IPv6 requires a UDP checksum. Modern stacks compute the checksum.

The checksum is not a cryptographic hash. A middlebox or a bit flip can still produce a collision in theory. The checksum only catches common errors. Applications that need authenticity use TLS, DTLS, or an application MAC.

Hardware can offload the checksum. A capture on the sending host can show a zero or a wrong checksum before the NIC fills the field. Do not debug a "bad checksum" on the sender capture without care.

UDP-Lite exists for partial checksums on media. This handbook does not require UDP-Lite.

Do not disable UDP checksums for a "speed" trick on a general network. Do not treat a good checksum as proof that the sender is trusted.

### Questions

#### Theoretical questions

1. What does the UDP checksum cover (the idea)?
2. What does the receiver do with a bad checksum?
3. Why does IPv6 require a UDP checksum?
4. Why is the checksum not a cryptographic check?
5. Why can a sender-side capture show a strange checksum?

#### Easy practical tasks

1. Write five sentences that explain why UDP checks the payload.
2. In Wireshark, open one UDP packet. Write whether the tool says the checksum is valid.
3. Make a table: checksum, TLS. Add one row for "what attack each stops".
4. Draw the UDP header: ports, length, checksum.

#### Medium practical tasks

1. Compare an IPv4 UDP DNS packet and an IPv6 UDP packet in a capture. Write the checksum field presence.
2. Write six sentences: NIC offload and why a local capture can confuse you.
3. Find the UDP length field. Compare it with the payload size that the tool shows.

#### Advanced practical tasks

1. Read a short public note on the UDP pseudo-header. Write six STE sentences. No RFC dump.
2. Write a one-page note: when an application adds its own CRC or MAC on top of UDP.

---

## No reliability, no order

UDP does not retransmit. If a datagram is lost, UDP does not send it again. The application or a layer above UDP must retry if a retry is required.

UDP does not number datagrams. The network can reorder packets. The receiver can get message 2 before message 1. The application must tolerate reorder or must add sequence numbers.

UDP does not pace the sender. A fast sender can overflow a slow receiver or a thin link. That is not "flow control" in the TCP sense. Some applications add their own windows.

UDP does not avoid congestion by itself. A bad UDP sender can fill a link. QUIC and many media stacks add congestion control in the application or in a library. Raw UDP has no such state.

Duplicates can occur. A retry at the application can meet a late original datagram. The application must ignore duplicates if they break the logic.

Do not use raw UDP for a file copy unless you write a protocol. Do not blame UDP when you forgot to handle loss.

### Questions

#### Theoretical questions

1. What does UDP do when a datagram is lost?
2. Can UDP datagrams arrive out of order?
3. Does UDP provide flow control?
4. Who must add congestion control if you need it on UDP?
5. Why can duplicates appear?

#### Easy practical tasks

1. Write five sentences that list what UDP does not provide.
2. Make a table: loss, reorder, duplicate. Add one application reaction each.
3. Draw three datagrams with numbers 1, 2, 3. Show an arrival order 1, 3, 2.
4. Write four sentences: why a game can accept loss but a bank file cannot.

#### Medium practical tasks

1. Design on paper a tiny protocol: sequence number and ack. State what UDP still does not do.
2. Write six sentences: a sender that loops `send` as fast as possible on UDP. What can go wrong.
3. Compare DNS retry (application retry) with TCP retransmission (later topic). Write a short paragraph.

#### Advanced practical tasks

1. Write a UDP program that adds a sequence number. Send many messages on a lab. Detect reorder if you can create it, or simulate reorder in the receiver test.
2. Write a one-page note: how QUIC on UDP adds reliability as an option. High-level only.

---

## When UDP is the right choice (DNS, games, QUIC, media)

Choose UDP when one independent message is the natural unit. DNS queries are short request-response datagrams. A second query is a new datagram. Timeout and retry sit in the resolver.

Choose UDP when late data is useless. Games and media often prefer a new sample over a late retransmission. A late packet can be worse than a loss.

Choose UDP when you want to build your own transport. QUIC uses UDP so that it can run in user space and can change without kernel TCP. HTTP/3 uses QUIC. You will study that later. This topic only needs: QUIC is not "unreliable HTTP." QUIC adds its own reliability and TLS.

Choose UDP for some tunnels and for DHCP. The protocols are simple and local or short.

Choose TCP when you need a reliable byte stream and you do not want to write that logic. File transfer, mail, and many APIs use TCP (or QUIC).

Do not choose UDP only because it is "faster." A lost UDP datagram can make the user wait for an application retry. Measure.

Do not put a large file on raw UDP in production.

### Questions

#### Theoretical questions

1. Why does DNS fit UDP (classic case)?
2. Why do games often accept loss?
3. Why does QUIC run on UDP?
4. When is TCP the better default?
5. Why is "UDP is faster" an incomplete reason?

#### Easy practical tasks

1. Write five sentences that match DNS, a game, QUIC, and a file copy to UDP or TCP.
2. Make a table: application, UDP or TCP, one reason. Add four rows.
3. List two media examples (voice, video) and why a late packet is a problem.
4. Draw a DNS query and reply as two datagrams.

#### Medium practical tasks

1. Capture a DNS query. Write payload size and whether one datagram was enough.
2. Write six sentences: HTTP/1.1 on TCP versus HTTP/3 on QUIC (idea only).
3. Find whether a game or voice app you use documents UDP ports. Write the ports if public.

#### Advanced practical tasks

1. Write a one-page decision sheet: five questions that choose UDP or TCP for a new service.
2. Design a tiny media sender: 20 ms packets, drop late packets. Write the rules. You do not need to implement the full codec.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a DNS query from ports to checksum to "no retry in UDP."
2. How do NAT mappings and connectionless datagrams work together?
3. Why must an application add sequence numbers if order matters?
4. Which facts tell you that UDP is the right tool for a given service?
5. A teammate wants raw UDP for a backup file. Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: ports, datagram, checksum, no reliability, when to use UDP.
2. Capture one UDP flow. Write 4-tuple, length, and checksum valid or not.
3. Draw UDP versus TCP as two boxes of features. Use only this topic for UDP.
4. Run `dig` and write the UDP port pair from a capture or from verbose output.

#### Medium practical tasks

1. Write a lab: UDP echo server and client. Send three messages. Document loss if you see it.
2. Compare UDP DNS and a later TCP HTTP capture in a table of eight cells (you can take HTTP in the next topics).
3. Write a troubleshooting flow: "UDP service does not work" that checks bind, firewall protocol, and NAT timeout.

#### Advanced practical tasks

1. Write a glossary of 10 UDP terms. Each entry: term, one sentence, one tool field.
2. Implement a UDP ping: timestamp in the payload, print RTT. Handle loss with a timeout. No solutions in this file; do the work in your editor.
