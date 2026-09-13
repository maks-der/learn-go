# 13. Debugging

## Description

Observability means you can see enough of the path and the endpoints to explain a fault. This topic shows `ping`, `traceroute`, `ss`, and `tcpdump`. You learn Wireshark filters. You learn loss versus delay versus reorder. You learn TLS handshake failures and DNS TTL surprises.

Complete this topic after transport, DNS, TLS, and middleboxes.

Use one term for each concept. A symptom is what the user sees. A cause sits on one layer. Loss is a missing packet. Delay is late arrival. Reorder is a change of packet order. Do not mix a capture filter with a display filter. Do not stop at ping when the app uses TCP 443.

Practice on hosts and networks that you own or that you have permission to test. Do not scan. Do not flood. Do not publish captures that contain secrets.

---

## `ping`, `traceroute`, `ss`, `tcpdump`

These four tools answer different questions.

`ping` sends ICMP echo or an OS variant. A reply is reachability for that probe plus an RTT sample. A missing reply is not always host down. Topic 4 stated that.

`traceroute`, `tracert`, or `mtr` maps hops with TTL. Stars are not always a broken path. Topic 4 stated that.

`ss` on Linux and `netstat` on Windows show sockets. Listen sockets and established tuples tell you whether the process bound the port that you think. Teaching commands: `ss -lnt` and `netstat -an`. `ss` is the Linux name in the topic list. Use `netstat` when `ss` is absent.

`tcpdump` writes packets on one interface with a capture filter. It is the command-line partner of Wireshark. Topic 1 introduced capture safety.

Order of use for a cannot-connect report:

1. Is the name resolved?
2. Does `ping` or a TCP probe reach the address?
3. Is anyone listening (`ss` or `netstat`)?
4. What do the packets show (`tcpdump` or Wireshark)?

`curl -v` sits between these tools and the application. Use it for HTTP and TLS.

Do not start with a full packet dump when `ss` already shows no listener. Do not ping only and then stop when the app uses TCP 443 and ICMP is filtered.

### Questions

#### Theoretical questions

1. What does a ping reply prove?
2. What does `ss -lnt` show?
3. Why can traceroute show stars when HTTPS works?
4. What is a capture filter for?
5. Why is ping failed not the end of an HTTP debug?

#### Easy practical tasks

1. Write five sentences that assign one question to each of the four tools.
2. Ping your gateway. Write RTT.
3. Run `ss -lnt` or `netstat -an`. Copy two listen lines.
4. Make a table: tool, question it answers. Add four rows.

#### Medium practical tasks

1. Traceroute to a public host. Mark local versus remote hops.
2. Write six sentences: no listener on the port versus a filter that drops SYN.
3. Run `tcpdump` or Wireshark with a capture filter for ICMP while you ping. Write the filter text.

#### Advanced practical tasks

1. Write a one-page runbook: four tools in order for connection refused versus timed out.
2. On a local server that you own, show listen, then connect, then capture the handshake. Document commands.

---

## Wireshark filters

A display filter hides packets in the view. It does not drop them from the file unless you export the displayed set. A capture filter drops packets before they enter the file. Topic 1 named both. This section builds display-filter skill.

Useful display filters for this path:

- `icmp` or `icmpv6`
- `dns`
- `tcp.port == 443`
- `udp.port == 53`
- `tls`
- `http` (cleartext HTTP)
- `ip.addr == 192.0.2.1` (example address)
- `tcp.flags.syn == 1`
- `tcp.analysis.retransmission`
- `quic` when the tool supports it

Capture filters use a different language (`tcp port 443`, `host 192.0.2.1`). Do not paste a display filter into `tcpdump` and expect it to work.

Color and the packet detail tree show layers. Follow TCP stream reassembles a stream for cleartext. HTTPS streams stay encrypted.

Do not filter so hard that you hide the ICMP error that explains the TCP stall. Do not share a pcap with cookies. Use a display filter to find the event, then export a small range.

Name the layer in your notes: Ethernet, IP, TCP, TLS, HTTP.

### Questions

#### Theoretical questions

1. What does a display filter change?
2. What does a capture filter change?
3. Why can `http` show nothing on an HTTPS site?
4. What filter finds SYN segments in this section?
5. Why must you not paste a display filter into `tcpdump`?

#### Easy practical tasks

1. Write five sentences that compare the two filter kinds.
2. Open a capture. Apply `dns`. Write how many packets remain.
3. Make a table: display filter, what you hope to see. Add four rows from this section.
4. Write a capture-filter form for TCP port 443.

#### Medium practical tasks

1. Capture a `curl` to HTTPS. Apply `tls` and `tcp.flags.syn == 1`. Write what each view shows.
2. Write six sentences: follow TCP stream on HTTP versus HTTPS.
3. Find a retransmission filter hit or write that the capture has none.

#### Advanced practical tasks

1. Write a one-page filter cheat sheet: ten display filters and three capture filters that you will reuse.
2. Export a tiny pcap of one handshake with no secrets. Write the filters that you used.

---

## Loss vs delay vs reorder

Loss means a packet does not arrive. TCP retransmits. UDP does not. The application sees a stall, a retry, or a gap.

Delay means packets arrive late. RTT grows. Queues on a busy link add delay. A satellite or a distant CDN origin adds delay. Delay is not loss. `ping` min versus max shows jitter when the values spread.

Reorder means packets arrive in a different order than the sender sent. TCP can reorder in the receive buffer. Duplicate ACKs can appear. UDP applications must tolerate reorder or they must drop late data.

The three faults need different controls. More bandwidth does not fix a long RTT. A bigger window does not fix 20 percent loss. A CDN can cut delay. FEC or a lower rate can help loss on radio. Do not guess. Measure.

Tools: `ping` for RTT samples, `mtr` or traceroute for hop delay, Wireshark for retransmission and out-of-order flags, application timers for user delay.

Do not call every slow page packet loss. Do not call every retransmission a congested core. Wi-Fi noise is loss too.

### Questions

#### Theoretical questions

1. What is loss?
2. What is delay?
3. What is reorder?
4. Why does more bandwidth not fix a long RTT?
5. How does TCP show loss in a capture?

#### Easy practical tasks

1. Write five sentences that define the three words.
2. Make a table: loss, delay, reorder. Add one symptom each.
3. Ping a host 20 times. Write min, max, and whether you saw loss.
4. Draw a late packet and a missing packet as two pictures.

#### Medium practical tasks

1. In Wireshark, find `tcp.analysis.out_of_order` or retransmission. Write which you found.
2. Write six sentences: voice on UDP and a late packet.
3. Compare gateway RTT and a distant host RTT. Write the extra delay.

#### Advanced practical tasks

1. Write a one-page decision sheet: which tool first for loss, delay, and reorder.
2. On a lab that you own, add delay with a netem-like tool if you have it, or use a distant host. Record ping and a small `curl` time.

---

## TLS handshake failures and DNS TTL surprises

TLS handshake failures abort before HTTP. Topic 10 listed classes: name mismatch, clock and expiry, untrusted issuer, protocol mismatch, interception, verify off.

Debug order for TLS:

1. Does DNS return the address that you think?
2. Does TCP connect (`curl -v` shows connect, `ss` shows ESTABLISHED)?
3. Does Client Hello leave (capture `tls`)?
4. Does the server send a certificate?
5. What does `curl` or the browser print (verify, name, clock)?

A capture that shows only SYN without SYN-ACK is not a TLS problem. A capture that shows a TLS alert is a TLS problem.

DNS TTL surprises happen after a record change. Old caches still serve the old A or AAAA. Some users reach the new host. Some users reach the old host. Some users get NXDOMAIN if you removed a name too soon.

Debug order for TTL:

1. `dig` at two resolvers
2. compare TTL and answers
3. flush only a cache that you own
4. keep the old host up until TTLs end

Encrypted DNS in the browser can disagree with `dig` (topic 9). Always name which resolver you queried.

Do not lower TTL after the outage only. Do not blame TLS when `dig` already points to a dead address. Do not publish certificate errors that contain internal host names if policy forbids it.

### Questions

#### Theoretical questions

1. Why is SYN without SYN-ACK not a TLS failure?
2. What is the first TLS question after TCP connects?
3. Why can two users see two different IPs after a DNS change?
4. Why can `dig` and a browser disagree?
5. Why keep the old host up until TTL ends?

#### Easy practical tasks

1. Write five sentences that split TLS fail from DNS fail.
2. Make a table: symptom, first command. Add cert error and wrong A.
3. Run `dig example.com A` and `curl -v https://example.com`. Write address and TLS result.
4. Draw: stale cache user versus fresh cache user.

#### Medium practical tasks

1. Write six sentences: expired cert versus NXDOMAIN.
2. Query one name at two resolvers. Write both answers and TTLs.
3. Read a `curl` TLS error on a lab that you own (self-signed). Write the class from topic 10.

#### Advanced practical tasks

1. Write a one-page combined runbook: DNS, TCP, TLS, HTTP status. Four stages.
2. Simulate a TTL surprise on a lab domain that you own, or write the time math for TTL 3600 and a change at 10:00.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you choose among ping, ss, and a TLS error string for one user report?
2. When must you open Wireshark instead of adding another `curl -v` flag?
3. How can loss, delay, and a stale DNS answer produce the same user sentence: it is slow or broken?
4. Why does a display filter that shows only `http` hide the real HTTPS fault?
5. What is the safest way to share debug evidence without secrets?

#### Easy practical tasks

1. Write a cheat sheet: ping, traceroute, ss, tcpdump, display versus capture, loss, delay, reorder, TLS fail, TTL.
2. Run the four tools once each on your machine. Write one line of output each.
3. Apply `dns` and `tls` to a capture. Write packet counts.
4. Bookmark the Wireshark filter reference. Write when you open it.

#### Medium practical tasks

1. Debug a local closed port and a local listening port. Write the tool that distinguished them.
2. Write a table: user symptom, layer, first tool, second tool. Add five rows.
3. Write six sentences: HTTP/3 UDP miss in a TCP-only capture.

#### Advanced practical tasks

1. Write a full lab report: start a server, break DNS hosts file or a cert, then recover. Use only a host that you own.
2. Build a personal filter and command card. Use it on one real fault this week and attach notes (no secrets).
