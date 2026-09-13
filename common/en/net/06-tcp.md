# 6. TCP

## Description

TCP is a transport protocol that provides a reliable byte stream. TCP builds a connection, numbers bytes, retransmits lost data, and controls the send rate. This topic shows the three-way handshake. You learn sequence numbers, retransmission, and the window. You learn congestion control: slow start and AIMD. You learn close and `TIME_WAIT`. You learn head-of-line blocking.

Complete this topic after UDP. Complete this topic before you write socket programs that assume a stream.

Use one term for each concept. A segment is the TCP unit on the wire. A connection is the state that the two hosts keep. The receive window protects the receiver. The congestion window protects the network. Do not mix flow control with congestion control. Do not treat TCP as a message protocol.

TCP does not keep message boundaries. The application must parse the stream.

---

## Three-way handshake

A TCP connection starts with a three-way handshake. The client sends a segment with the SYN flag and an initial sequence number. The server replies with SYN and ACK, and with its own initial sequence number. The client sends ACK. After that, both sides can send data.

The handshake agrees on initial sequence numbers and on options such as the maximum segment size (MSS). The handshake also proves that both sides can receive packets in the common case.

A server socket in the listen state accepts new SYNs. The kernel keeps a queue of connections that finished the handshake. `accept` in the next topic takes one connection from that queue.

A handshake can fail. A firewall can drop SYN. The server can send RST if no process listens. The client can time out. A SYN flood is a later security topic. This section only needs: SYN uses kernel memory.

The first data byte can ride with the third ACK in some stacks. The three flags SYN, SYN-ACK, ACK are the teaching sequence.

Do not send application data as if the connection exists before the handshake completes (unless you study TCP Fast Open later). Do not call the handshake UDP-like. UDP has no handshake.

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

1. Capture `curl` to a website. Write client ISN and server ISN if the tool shows them, or write relative numbers.
2. Connect to a closed port on a host that you own. Capture RST if it appears. Write the flags.
3. Write six sentences: listen queue and why a busy server can refuse new SYNs.

#### Advanced practical tasks

1. Capture a failed handshake (wrong port or dropped SYN in a lab). Write the last flag that you saw.
2. Write a one-page note: MSS option in the SYN and how it relates to MTU from topic 2.

---

## Sequence numbers, retransmission, and the window

TCP numbers the byte stream. The sequence number in a segment is the number of the first byte in that segment (after the handshake consumes one sequence number for SYN). The acknowledgment number is the next byte that the sender of the ACK expects.

A cumulative ACK says: I have all bytes before this number. Selective ACK (SACK) is an option that names extra blocks. This topic needs cumulative ACK first.

Sequence numbers are 32-bit and wrap. Relative numbers in Wireshark start at 0 for display. The wire has the absolute value. Each side has its own sequence space. Do not add them.

The sender keeps a copy of unacknowledged bytes. If an ACK does not arrive in time, the sender retransmits. Loss detection uses a retransmission timeout (RTO) and fast retransmit (often after several duplicate ACKs).

Reliability is not a bound on delay. A terrible path can retransmit for a long time. Applications still need timeouts. A successful `write` means the local stack accepted the bytes. It does not prove that the peer read them. A RST ends the stream.

Flow control protects the receiver. The receiver advertises a window: how many more bytes it can accept. The sender must not send more unacknowledged data than the window allows. A window of zero means the sender must stop except for window probes.

Flow control is not congestion control. Flow control is about the receive buffer. Congestion control is about the network path.

Do not treat a sequence number as a packet counter. A large segment jumps the sequence by many bytes. Do not set a tiny socket receive buffer without a reason.

### Questions

#### Theoretical questions

1. What does the sequence number in a data segment name?
2. What does the acknowledgment number mean?
3. What does the sender keep until ACK?
4. What does the advertised window mean?
5. How is flow control different from congestion control?

#### Easy practical tasks

1. Write five sentences that explain sequence and ACK with a 100-byte send.
2. Draw a lost segment and a later retransmit with the same sequence.
3. In a capture, write a sequence number and the next ACK that you see.
4. Make a table: flow control, congestion control. Add one sentence each.

#### Medium practical tasks

1. Capture a short HTTP request. Compute how the sequence advances with the payload length.
2. Write six sentences: a slow reader and a shrinking window.
3. Find `tcp.analysis.retransmission` in Wireshark on any capture. If none, write how you searched.

#### Advanced practical tasks

1. Write a program that reads slowly from TCP on a lab that you own. Capture the window.
2. Write a one-page note: SACK versus cumulative ACK. High-level only.

---

## Congestion control (slow start, AIMD)

Congestion control protects the network. Too many senders can fill a queue and drop packets. TCP treats loss (and in modern stacks, delay signals) as a sign to slow down.

Slow start: the congestion window (`cwnd`) starts small and grows quickly (often exponentially) while ACKs arrive. The sender can send more data in flight.

AIMD: additive increase, multiplicative decrease. After slow start, the window grows slowly. After loss, the window drops by a large factor (classic Reno halves). Details differ by algorithm (Cubic, BBR). This topic needs the idea, not the formulas of every variant.

The sender uses the minimum of `cwnd` and the receive window. The path and the receiver both can limit the rate.

Congestion control is why a single TCP flow does not stay at full line rate after a drop. It is also why many flows can share a link in a rough way.

A new connection pays slow start. Keep-alive HTTP (topic 8) and HTTP/2 (topic 11) reuse a connection so that later requests skip a new slow start.

Do not disable congestion control to go faster on the public Internet. Do not treat every loss as congestion on Wi-Fi. Loss can be noise. Modern stacks try to tell the difference. The idea still starts with "loss can mean slow down."

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
2. Write six sentences: Cubic versus the AIMD idea at a high level.
3. Time a large download. Write that you cannot see `cwnd` without extra tools, and name one extra tool (`ss -i` on Linux).

#### Advanced practical tasks

1. On a Linux lab, run `ss -i` during a transfer. Write one `cwnd` value. Skip if you have no Linux.
2. Write a one-page comparison of loss-based and delay-based ideas (Reno or Cubic versus BBR at a high level).

---

## Close and `TIME_WAIT`

A clean TCP close uses FIN segments. Each side closes its send direction. The common teaching sequence is a four-segment close: FIN, ACK, FIN, ACK. Some stacks combine FIN and ACK. The idea is half-close: one side can still receive after it stopped sending.

A RST aborts the connection. Data in flight can be lost from the application view. Prefer a clean close when you can.

`TIME_WAIT` is a state on the side that performed the active close. That host keeps the 5-tuple busy for a short time (often two maximum segment lifetimes). The wait drops delayed segments from the old connection so that a new connection with the same tuple does not receive old data.

Many short client connections can fill `TIME_WAIT` on the client. That is a reason to reuse connections. Servers that active-close many connections can also see `TIME_WAIT`. `SO_REUSEADDR` (topic 7) helps bind during some waits. It does not delete the need for the state.

Do not kill `TIME_WAIT` with a kernel tweak on a public server without a reason. Do not treat `TIME_WAIT` as a leak. It is a safety timer.

The application sees a read of zero bytes on a clean close. That is not UDP.

### Questions

#### Theoretical questions

1. What flag starts a clean close?
2. What is a half-close?
3. How does RST differ from FIN?
4. Why does `TIME_WAIT` exist?
5. Which side usually enters `TIME_WAIT`?

#### Easy practical tasks

1. Write five sentences that describe FIN and `TIME_WAIT`.
2. Draw a four-segment close.
3. Make a table: FIN close, RST abort. Add what the peer can lose.
4. Run `ss -tan` or `netstat -an`. Write whether you see `TIME_WAIT`.

#### Medium practical tasks

1. Capture a short `curl` to HTTP or HTTPS. Find FIN or RST at the end. Write which you see.
2. Write six sentences: many short connections and `TIME_WAIT` on the client.
3. Compare a clean echo-server close with a kill of the process. Write the difference you expect in a capture.

#### Advanced practical tasks

1. Write a TCP client that closes and a server that prints a zero-length read. Use a machine that you own.
2. Write a one-page note: two MSL and delayed segments. High-level only.

---

## Head-of-line blocking

Head-of-line (HOL) blocking means a loss or a stall at the front of a queue delays everything behind it.

In TCP, one lost segment delays delivery of later bytes to the application until the gap fills. TCP must present an ordered stream. HTTP/2 multiplexes many requests on one TCP connection. A lost TCP segment still stalls all streams at the TCP layer. That is a reason HTTP/3 uses QUIC on UDP (topic 11).

HOL also appears at other layers. HTTP/1.1 on one connection waits for one response before the next request in the usual client. A switch queue can delay a large flow and a small flow that share a port. Name the layer when you say HOL.

TCP HOL is not a bug in the spec. It is the cost of an ordered reliable stream.

Do not open hundreds of TCP connections to avoid HOL without a measurement. Each connection pays a handshake and slow start. Do not blame HTTP headers for a TCP loss stall on HTTP/2.

### Questions

#### Theoretical questions

1. What does head-of-line blocking mean?
2. Why does a lost TCP segment delay later bytes to the application?
3. Why does HTTP/2 still have TCP HOL?
4. Why does HTTP/3 move to QUIC in this story?
5. Why is TCP HOL not a spec bug?

#### Easy practical tasks

1. Write five sentences that explain TCP HOL with one lost segment and two later segments.
2. Draw a gap in the sequence and data that waits in the receiver.
3. Make a table: HTTP/1.1 wait, TCP HOL, HTTP/3 idea. Add one line each.
4. Write four sentences: ordered stream versus independent messages.

#### Medium practical tasks

1. Write six sentences: one image and one API call on one HTTP/2 TCP connection when a segment is lost.
2. In a capture, find a retransmission. Write that later sequence numbers wait in the idea even if the tool does not show the app stall.
3. Read a public short note on HTTP/2 HOL. Write four STE sentences.

#### Advanced practical tasks

1. Write a one-page comparison: TCP stream HOL versus UDP independent datagrams. Use topics 5 and 6.
2. Design on paper two requests that must not share one TCP connection if loss on one must not delay the other. Then note HTTP/3 as another option.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the handshake, sequence numbers, and the window work together for the first 1000 bytes?
2. Why does a new TCP connection feel slower than a reused connection on a lossy path?
3. What should an application still time out if TCP is reliable?
4. How do `TIME_WAIT` and HOL both come from TCP's care for an ordered stream?
5. Where does RST sit relative to retransmission: does retry continue?

#### Easy practical tasks

1. Write a cheat sheet: SYN, ISN, ACK, window, cwnd, slow start, AIMD, FIN, TIME_WAIT, HOL.
2. Capture one handshake and write the three flags in order.
3. Draw send window as the min of receive window and cwnd.
4. Run `ss -tan` or `netstat -an` after a web browse. Write two TCP states that you see.

#### Medium practical tasks

1. Write a lab report: handshake, one data segment, ACK, close. Use `curl` and Wireshark.
2. Compare UDP from topic 5 and TCP in a completed table: handshake, order, retry, close, HOL.
3. Write six sentences: a zero window and a congestion window drop as two different slows.

#### Advanced practical tasks

1. Skim the TCP state machine names in RFC 9293. Write LISTEN, ESTABLISHED, TIME_WAIT, and CLOSE-WAIT in your own words.
2. Write a one-page teaching talk: why TCP is a stream, with one capture screenshot description (no secrets).
