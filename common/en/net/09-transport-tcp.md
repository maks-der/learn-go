# 9. Transport: TCP

## Description

TCP is a transport protocol that provides a reliable byte stream. TCP builds a connection, numbers bytes, retransmits lost data, and controls the send rate. This topic shows the three-way handshake, sequence and acknowledgment numbers, retransmission, flow control, congestion control, the four-way close, `TIME_WAIT`, and head-of-line blocking.

Use one term for each concept. A segment is the TCP unit on the wire. A connection is the state that the two hosts keep. Complete this topic before you write socket programs that assume a stream.

TCP does not keep message boundaries. The application must parse the stream.

---

## Connection: three-way handshake

A TCP connection starts with a three-way handshake. The client sends a segment with the SYN flag and an initial sequence number. The server replies with SYN and ACK, and with its own initial sequence number. The client sends ACK. After that, both sides can send data.

The handshake agrees on initial sequence numbers and on options such as the maximum segment size (MSS). The handshake also proves that both sides can receive packets (in the common case).

A server socket in the listen state accepts new SYNs. The kernel keeps a queue of connections that finished the handshake. `accept` in the next topic takes one connection from that queue.

A handshake can fail. A firewall can drop SYN. The server can send RST if no process listens. The client can time out. A SYN flood is a later security topic. This section only needs: SYN uses kernel memory.

Do not send application data as if the connection exists before the handshake completes (unless you study TCP Fast Open later). Do not call the handshake "UDP-like." UDP has no handshake.

The first data byte can ride with the third ACK in some stacks (or later). The three flags SYN, SYN-ACK, ACK are the teaching sequence.

### Questions

#### Theoretical questions

1. What three segments form the common handshake?
2. What does the SYN flag ask for?
3. What does the handshake agree on besides "we are open"?
4. What can the server send if no process listens?
5. Does UDP use this handshake?

#### Easy practical tasks

1. Write five sentences that describe SYN, SYN-ACK, and ACK.
2. Draw the three arrows between client and server. Label flags.
3. Make a table: segment, flags, who sends it.
4. In Wireshark, filter `tcp.flags.syn == 1`. Find a handshake. Write the three times.

#### Medium practical tasks

1. Capture `curl` to a website. Write client ISN and server ISN if the tool shows them (or write relative numbers).
2. Connect to a closed port on a host that you own. Capture RST if it appears. Write the flags.
3. Write six sentences: listen queue and why a busy server can refuse new SYNs.

#### Advanced practical tasks

1. Capture a failed handshake (wrong port or dropped SYN in a lab). Write the last flag that you saw.
2. Write a one-page note: MSS option in the SYN. How it relates to MTU from topic 3.

---

## Sequence and ack numbers

TCP numbers the byte stream. The sequence number in a segment is the number of the first byte in that segment (after the handshake consumes one sequence number for SYN). The acknowledgment number is the next byte that the sender of the ACK expects.

A cumulative ACK says: I have all bytes before this number. Selective ACK (SACK) is an option that names extra blocks. This topic needs cumulative ACK first.

Sequence numbers are 32-bit and wrap. Relative numbers in Wireshark start at 0 for display. The wire has the absolute value. Use relative numbers when you teach. Use absolute numbers when you debug a wrap or a middlebox.

Each side has its own sequence space. The client sequence and the server sequence are independent. Do not add them.

Retransmissions reuse the sequence numbers of the lost bytes. Duplicate ACKs repeat an acknowledgment number when a gap appears.

Do not treat a sequence number as a packet counter. A large segment jumps the sequence by many bytes. Do not treat the ACK number as "how many packets I got."

### Questions

#### Theoretical questions

1. What does the sequence number in a data segment name?
2. What does the acknowledgment number mean?
3. What is a cumulative ACK?
4. Why does Wireshark show relative sequence numbers?
5. Do the two sides share one sequence space?

#### Easy practical tasks

1. Write five sentences that explain sequence and ACK with a 100-byte send.
2. Draw a send of 10 bytes and an ACK of the next expected byte.
3. In a capture, write a sequence number and the next ACK that you see.
4. Make a table: SYN, data, FIN. Add how many sequence numbers each consumes (teaching rule).

#### Medium practical tasks

1. Capture a short HTTP request. Compute how the sequence advances with the payload length.
2. Write six sentences: a gap in the stream and a duplicate ACK (idea).
3. Turn off relative sequence numbers in Wireshark for one packet. Write the absolute value.

#### Advanced practical tasks

1. Find a retransmission in a capture (`tcp.analysis.retransmission`). Write the sequence number that was sent again.
2. Write a one-page note: SACK versus cumulative ACK. High-level only.

---

## Reliable delivery and retransmission

TCP delivers a reliable ordered stream. The sender keeps a copy of unacknowledged bytes. If an ACK does not arrive in time, the sender retransmits. The receiver can reorder segments and can hold out-of-order data until the gap fills.

Loss detection uses a retransmission timeout (RTO) and fast retransmit (often after several duplicate ACKs). Exact timers are complex. This topic needs: timeout and fast retransmit exist.

The receiver sends ACKs. Delayed ACK can wait for a second segment or a short timer. That delay interacts with Nagle (later performance topic). For now: ACKs are not always immediate.

Reliability is not a bound on delay. A terrible path can retransmit for a long time. Applications still need timeouts.

TCP checksums the header and payload (with a pseudo-header). A bad segment is dropped. Reliability then looks like loss.

Do not assume that "TCP never loses data" if the connection resets. A RST ends the stream. Do not assume that a write() return means the peer read the bytes. It means the local stack accepted the bytes.

### Questions

#### Theoretical questions

1. What does the sender keep until ACK?
2. What two ideas detect loss in this section?
3. Is reliability the same as a short delay?
4. What happens to a segment with a bad TCP checksum?
5. Does a successful `write` prove that the peer read the data?

#### Easy practical tasks

1. Write five sentences that describe retransmission after loss.
2. Draw a lost segment and a later retransmit with the same sequence.
3. Make a table: ACK received, timeout, fast retransmit. Add one action each.
4. In Wireshark, find `tcp.analysis.retransmission` on any capture. If none, write how you searched.

#### Medium practical tasks

1. Use a lab or a poor Wi-Fi moment to capture a retransmission. Write the time between the two sends.
2. Write six sentences: delayed ACK as a reason that ACKs are not one-per-segment.
3. Compare UDP "application retry" and TCP retransmission in a short table.

#### Advanced practical tasks

1. Read a short public note on RTO. Write six STE sentences. No RFC dump.
2. Write a one-page note: connection reset versus a clean close. What the application sees.

---

## Flow control (window)

Flow control protects the receiver. The receiver advertises a window: how many more bytes it can accept. The sender must not send more unacknowledged data than the window allows.

The window is in the TCP header. Window scaling is an option that extends the range on fast long paths. The handshake must agree on scaling.

A window of zero means the sender must stop (except window probes). The receiver later sends a window update.

Flow control is not congestion control. Flow control is about the receive buffer. Congestion control is about the network path.

A small receive buffer or a slow application `read` shrinks the window. The sender slows down even if the network is empty.

Do not set a tiny socket receive buffer without a reason. Do not confuse the window with the congestion window (`cwnd`).

### Questions

#### Theoretical questions

1. What does the advertised window mean?
2. Who does flow control protect?
3. What does a zero window mean?
4. How is flow control different from congestion control?
5. What is window scaling for?

#### Easy practical tasks

1. Write five sentences that explain a receive window.
2. In a capture, find the window field on an ACK. Write the value.
3. Draw a receiver buffer and a window advertisement.
4. Make a table: flow control, congestion control. Add one sentence each.

#### Medium practical tasks

1. Write six sentences: a slow reader and a shrinking window.
2. Find whether a handshake in your capture has a window scale option.
3. Compare two ACKs in one flow. Write if the window grew or shrank.

#### Advanced practical tasks

1. Write a program that reads slowly from TCP on a lab. Capture the window. Use a machine that you own.
2. Write a one-page note: bandwidth-delay product and why a large window is needed on a long fat pipe (preview of later topics).

---

## Congestion control (idea: slow start, AIMD)

Congestion control protects the network. Too many senders can fill a queue and drop packets. TCP treats loss (and in modern stacks, delay signals) as a sign to slow down.

Slow start: the congestion window (`cwnd`) starts small and grows quickly (often exponentially) while ACKs arrive. The sender can send more data in flight.

AIMD: additive increase, multiplicative decrease. After slow start, the window grows slowly. After loss, the window drops by a large factor (classic Reno halves). Details differ by algorithm (Cubic, BBR). This topic needs the idea, not the formulas of every variant.

The sender uses the minimum of `cwnd` and the receive window. The path and the receiver both can limit the rate.

Congestion control is why a single TCP flow does not stay at full line rate after a drop. It is also why many flows can share a link in a rough way.

Do not disable congestion control to "go faster" on the public Internet. Do not treat every loss as congestion on Wi-Fi (loss can be noise). Modern stacks try to tell the difference. The idea still starts with "loss can mean slow down."

### Questions

#### Theoretical questions

1. What does congestion control protect?
2. What is slow start in one sentence?
3. What does AIMD mean?
4. What two limits cap the data in flight?
5. Why must a public sender not turn congestion control off?

#### Easy practical tasks

1. Write five sentences that explain slow start as "start small, grow fast."
2. Draw `cwnd` as a rising line, then a drop after loss.
3. Make a table: slow start, AIMD increase, loss decrease.
4. Write four sentences: receive window versus `cwnd`.

#### Medium practical tasks

1. Find the congestion-control name on Linux (`sysctl net.ipv4.tcp_congestion_control`) or write the Windows default from a public doc. Record it.
2. Write six sentences: Cubic versus the AIMD idea (high-level).
3. Time a large download. Write that you cannot see `cwnd` without extra tools, and name one extra tool (ss -i on Linux).

#### Advanced practical tasks

1. On Linux lab, run `ss -i` during a transfer. Write one `cwnd` value. Skip if you have no Linux.
2. Write a one-page comparison of loss-based and delay-based ideas (Reno/Cubic versus BBR at a high level).

---

## Four-way close and `TIME_WAIT`

A clean TCP close uses FIN. Each side must close its send direction. Common teaching sequence: A sends FIN, B ACKs, B sends FIN, A ACKs. That is four segments. Some stacks combine FIN and ACK. You still need both directions to close.

After the active closer sends the last ACK, that host can enter `TIME_WAIT`. The host waits so that a late duplicate cannot open a new connection with the same 4-tuple and so that the last ACK can be resent if it was lost.

`TIME_WAIT` lasts a few minutes (often 2*MSL). Many short connections can fill the port space. `SO_REUSEADDR` in the next topic helps the listener. Do not kill `TIME_WAIT` as a first production trick.

A RST aborts the connection. Data in flight can be lost. Use RST only when the state is wrong.

Half-close: one side sends FIN and can still receive. The other side can still send. Some protocols use this. Many applications close both directions together.

Do not forget to close sockets. A process exit closes them. An explicit close is clearer.

### Questions

#### Theoretical questions

1. Why can a clean close need four segments?
2. What does FIN mean?
3. Why does `TIME_WAIT` exist?
4. How is RST different from FIN?
5. What is a half-close?

#### Easy practical tasks

1. Write five sentences that describe FIN from both sides.
2. Draw the four-way close. Label `TIME_WAIT` on one host.
3. In Wireshark, filter `tcp.flags.fin == 1`. Find a close. Write who sent the first FIN.
4. Make a table: FIN, ACK, RST. Add one meaning each.

#### Medium practical tasks

1. Capture a short `curl`. Write whether you see four close segments or a combined segment.
2. Run `ss -tan` or `netstat -an` after many short connections. Look for `TIME_WAIT`. Write a count if you see any.
3. Write six sentences: why a server that binds a port can fail just after a restart (preview `SO_REUSEADDR`).

#### Advanced practical tasks

1. Write a one-page note: 2*MSL and late duplicates. High-level.
2. In a lab, write a client that closes and a server that still sends. Observe half-close if you can.

---

## Head-of-line blocking

TCP delivers an ordered stream. If byte 1 is lost, the receiver does not give bytes 2 to 100 to the application until byte 1 arrives (or the connection fails). Those later bytes wait. That wait is head-of-line (HOL) blocking.

HOL blocking hurts when one stream carries many independent messages. HTTP/1.1 on one TCP connection can stall all requests when one packet is lost. HTTP/2 multiplexes many streams on one TCP connection. A loss still blocks all HTTP/2 streams at the TCP layer. HTTP/3 on QUIC avoids that TCP HOL problem. A later topic covers HTTP/2 and HTTP/3.

UDP does not have this TCP HOL behavior. A lost datagram does not hold the next datagram in the kernel stream. The application can read a later message.

HOL can also appear at other layers (a switch queue, an application lock). This section means TCP HOL.

Do not put unrelated urgent messages behind a large file on the same TCP connection if a stall is bad. Use another connection or another transport.

Do not confuse HOL with a zero window. A zero window is flow control. HOL here is ordered delivery after loss.

### Questions

#### Theoretical questions

1. What is TCP head-of-line blocking?
2. Why does a lost first byte block later bytes?
3. Why does HTTP/2 still have TCP HOL?
4. How does UDP differ for independent messages?
5. How is HOL different from a zero window?

#### Easy practical tasks

1. Write five sentences that explain HOL with a lost segment.
2. Draw a queue of bytes 1 to 5 with byte 1 missing.
3. Make a table: HTTP/1.1 one connection, HTTP/2 on TCP, HTTP/3. Add one HOL note each (idea).
4. Write four sentences: two TCP connections as a way to avoid one stall.

#### Medium practical tasks

1. Write six sentences: a chat message stuck behind a large download on the same TCP connection.
2. Find a public HOL diagram for HTTP/2 versus HTTP/3. Write four STE sentences. No long quotes.
3. Capture a loss if you can. Write whether later segments arrived before the retransmit.

#### Advanced practical tasks

1. Write a one-page note: why QUIC uses UDP to reduce TCP HOL for HTTP/3.
2. Design an app with one reliable stream and one unreliable stream. State why one TCP socket cannot give both.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a TCP session from handshake to data to close. Use sequence numbers and one retransmission.
2. How do flow control and congestion control limit the sender at the same time?
3. Why does `TIME_WAIT` appear after a short client connection?
4. How does HOL blocking change the way you map many messages onto one connection?
5. A teammate says TCP guarantees that the peer has processed the data. Which facts do you use to correct that sentence?

#### Easy practical tasks

1. Write a one-page cheat sheet: handshake, seq/ack, retry, window, `cwnd`, close, HOL.
2. Capture one full `curl` to HTTP or HTTPS. Label handshake, data, close in the packet list.
3. Draw the TCP state names you know (at least listen, established, time-wait) as boxes.
4. Compare one UDP DNS flow and one TCP web flow in a table of six cells.

#### Medium practical tasks

1. Write a lab report template with sections: handshake, first data seq, window, close flags.
2. Use `ss -tan` during a download. Write states that you see.
3. Write a troubleshooting flow: connect timeout versus RST versus a stall after connect.

#### Advanced practical tasks

1. Write a glossary of 14 TCP terms. Each entry: term, one sentence, one Wireshark field or flag.
2. Read RFC 9293 state-machine names (skim). Map five names to this topic in a table. Do not copy pages of the RFC.
