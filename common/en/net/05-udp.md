# 5. UDP

## Description

UDP is a transport protocol that sends datagrams. UDP does not build a connection. UDP does not retransmit. This topic shows ports and connectionless datagrams. You learn the checksum, the lack of reliability, and the lack of order. You learn when UDP is the right choice.

Complete this topic after routing and NAT. Complete this topic before you study TCP.

Use one term for each concept. A datagram is one UDP message. A port selects a socket on a host. The checksum detects some corruption. Do not mix a UDP datagram with a TCP byte stream. Do not expect UDP to retry for you.

UDP is simple. The application must handle loss, delay, and duplicate datagrams when those events matter.

---

## Ports and connectionless datagrams

A UDP port is a 16-bit number. The UDP header has a source port and a destination port. The destination port selects the receiving socket. The source port tells the receiver where to send a reply.

The pair of IP addresses, UDP, and the two ports is the 4-tuple that identifies a UDP flow in many firewalls and NAT boxes.

Well-known ports are the low numbers. DNS often uses UDP port 53. DHCP uses 67 and 68. Many games and voice tools use high ports that the vendor documents.

The client often uses an ephemeral source port. The server binds a known destination port. A server can also bind port 0 and let the OS pick a port. Then the client must learn that port by some other method.

Ports are per protocol. TCP port 53 and UDP port 53 are different sockets. A process can listen on both. Do not say "port 53" without TCP or UDP when the difference matters.

UDP is connectionless. There is no handshake before the first datagram. The sender writes a datagram. The network tries to deliver that datagram. The receiver may get it, get a duplicate, or get nothing.

One `send` is one datagram. One `recv` returns one datagram (up to your buffer). If the buffer is too small, the tail can be truncated. Two datagrams stay two datagrams. UDP keeps message boundaries.

A UDP socket can talk to many peers with `sendto` and `recvfrom`. `connect` on UDP sets a default peer. It does not handshake.

NAT mappings for UDP use timeouts because there is no FIN. A quiet flow can lose its mapping.

Do not `listen` on a UDP socket. Do not `accept` on UDP. A firewall rule must name UDP when the service is UDP. A TCP allow does not pass UDP.

### Questions

#### Theoretical questions

1. How many bits does a UDP port use?
2. What does the destination port select?
3. Are TCP port 53 and UDP port 53 the same socket?
4. What does connectionless mean for the first datagram?
5. Does UDP keep message boundaries?

#### Easy practical tasks

1. Write five sentences that explain source port and destination port for a DNS query.
2. Make a table: service, UDP port. Add DNS and DHCP.
3. Run `ss -lnu` or `netstat -an` and find a UDP listen line if one exists. Write it.
4. Draw a client ephemeral port and server port 53.

#### Medium practical tasks

1. Use `dig` or `nslookup` and a capture. Write the client UDP port and the server port.
2. Write six sentences: a firewall that allows TCP 53 only. What happens to normal DNS.
3. Two processes cannot bind the same UDP address and port in the usual case. Try on a lab host if you can. Record the error.

#### Advanced practical tasks

1. Write a one-page note: NAT mappings for UDP ports and why the source port on the WAN can differ.
2. Capture two DNS queries. Show that the client port can change. Write both 4-tuples.

---

## Checksum; no reliability and no order

The UDP header includes a checksum. The checksum covers the UDP header, the payload, and a pseudo-header with IP addresses. A bad checksum causes the receiver to drop the datagram. IPv4 allows a zero UDP checksum in old rules. IPv6 requires a UDP checksum. Do not rely on a zero checksum.

A checksum is not a cryptographic integrity check. A middlebox or an attacker can recompute it. TLS and DTLS sit above when you need authenticity.

UDP does not retransmit. If a datagram is lost, UDP does not send it again. The application can retry, or the user can retry, or the data can be skippable (a voice frame).

UDP does not order datagrams. The second send can arrive before the first. The application must use sequence numbers if order matters.

UDP does not prevent duplicates. A NAT or a retry can deliver two copies. The application must ignore a duplicate if that matters.

UDP does not have flow control or congestion control in the base protocol. A careless sender can flood a path. Modern uses (QUIC, some game protocols) add their own control. Topic 11 covers QUIC.

A successful `send` means the local stack accepted the datagram. It does not mean the peer received it.

Do not treat a UDP checksum as proof that the sender is the real peer. Do not build a file transfer on bare UDP without a reliability layer. Do not ignore truncation when the receive buffer is small.

### Questions

#### Theoretical questions

1. What does a bad UDP checksum cause at the receiver?
2. Why is a UDP checksum not authenticity?
3. What does UDP do when a datagram is lost?
4. Can UDP datagrams arrive out of order?
5. Does a successful `send` prove that the peer got the datagram?

#### Easy practical tasks

1. Write five sentences that list what UDP does not provide.
2. Make a table: loss, reorder, duplicate. Add one application response each.
3. Draw two datagrams that cross and arrive in reverse order.
4. Write four sentences: checksum versus TLS.

#### Medium practical tasks

1. Write six sentences: DNS over UDP and a lost query. What does the stub resolver do in the usual design?
2. Compare a 100-byte UDP message and a 100-byte TCP send in a short table: boundaries, retry, order.
3. In a capture of DNS, write whether you see one query datagram and one answer datagram.

#### Advanced practical tasks

1. Write a small UDP client and server that print a sequence number that you put in the payload. Send several datagrams on a lab that you own.
2. Write a one-page note: IPv6 UDP checksum requirement and why a zero checksum is a problem on IPv6.

---

## When UDP is the right choice

UDP is the right choice when you want a small, fast, message-shaped send and you accept loss or you add your own reliability.

Good fits:

- DNS queries that are short and that retry at the application
- DHCP discover and offer on a LAN
- Real-time voice or video where a late packet is useless
- Games that send frequent state and skip old state
- QUIC and HTTP/3, which add reliability on UDP
- Discovery protocols such as mDNS on the local link

Poor fits:

- a byte stream that must arrive complete and in order with no extra code
- a large file with no application retry
- a protocol that needs a clear connection close signal

TCP is the usual choice for HTTP/1.1, mail submission, and SSH. Topic 6 covers TCP. Topic 8 covers HTTP. You can still use UDP for HTTP/3.

Latency and simplicity matter. UDP has no handshake delay. The first datagram can be the request. NAT and firewalls still need a mapping. "No handshake" does not mean "always passes the edge."

Choose UDP when the application owns loss policy. Choose TCP when you want the kernel to own retransmission and order.

Do not pick UDP only because it feels faster. Measure. A lost request plus an application retry can be slower than TCP. Do not send a 10 MB UDP datagram. Path MTU and middleboxes will hurt you.

### Questions

#### Theoretical questions

1. Why is DNS a common UDP use?
2. Why can real-time voice prefer UDP?
3. When is TCP a better default than UDP?
4. Does "no handshake" mean the datagram always passes NAT?
5. Why is a huge UDP datagram a poor idea?

#### Easy practical tasks

1. Write five sentences that list three good UDP uses and two poor UDP uses.
2. Make a table: DNS, file download, live call. Add UDP or TCP as the usual choice.
3. Draw a voice frame that arrives late. Label it as skippable.
4. List DHCP, DNS, and HTTP/3 as UDP users in one column.

#### Medium practical tasks

1. Write six sentences: a game state tick on UDP versus a bank transfer on TCP.
2. Time `dig` for a name and `curl` for a small HTTPS page. Write that the protocols differ. Do not treat the times as a UDP-versus-TCP proof.
3. Read a public note on why HTTP/3 uses UDP. Write four STE sentences.

#### Advanced practical tasks

1. Write a one-page decision sheet: pick UDP or TCP for chat, file sync, live telemetry, and DNS. Give one reason each.
2. Design on paper a tiny UDP protocol: sequence number, timeout retry, max payload. Do not implement a flood.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the 4-tuple, NAT timeout, and connectionless send interact for a home DNS query?
2. Why must an application that needs order add fields when it uses UDP?
3. What does the UDP checksum protect, and what does it not protect?
4. How do ports stay independent from TCP even when the number is the same?
5. When does QUIC change the "UDP has no reliability" story without changing this topic's base UDP facts?

#### Easy practical tasks

1. Write a cheat sheet: port, datagram, connectionless, checksum, no retry, no order, when to use UDP.
2. Resolve a name with `dig` or `nslookup`. Write UDP port 53 as the usual server port.
3. Draw one UDP header: source port, destination port, length, checksum.
4. Run `ss -lnu` or the Windows equivalent. Write "no listen" if the list is empty.

#### Medium practical tasks

1. Capture a DNS query and answer. Write sizes. Write whether the answer was one datagram.
2. Write a short comparison table: UDP this topic versus TCP next topic. Leave TCP cells as questions you will fill after topic 6.
3. Explain in six sentences why a stateful firewall needs a timer for UDP.

#### Advanced practical tasks

1. Write a UDP echo server and client in your language. Show one truncated receive if you shrink the buffer. Use a host that you own.
2. Write a one-page lab: UDP through home NAT for DNS, and why an inbound UDP server on the LAN needs a forward.
