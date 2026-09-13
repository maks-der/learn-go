# 20. Observability and Debugging

## Description

Observability means you can see enough of the path and the endpoints to explain a fault. This topic shows `ping`, `traceroute`, `ss` / `netstat`, `tcpdump`, Wireshark display filters, loss versus delay versus reordering, application logs plus the path, TLS handshake failures, and DNS TTL surprises.

Use one term for each concept. A symptom is what the user sees. A cause sits on one layer. Complete this topic after transport, DNS, TLS, and middleboxes.

Practice on hosts and networks that you own or that you have permission to test. Do not scan. Do not flood. Do not publish captures that contain secrets.

---

## `ping`, `traceroute`, `ss` / `netstat`, `tcpdump`

These four tools answer different questions.

`ping` sends ICMP echo (or an OS variant). A reply is reachability for that probe plus an RTT sample. A missing reply is not always "host down." Topic 6 stated that.

`traceroute` / `tracert` / `mtr` maps hops with TTL. Stars are not always a broken path. Topic 6 stated that.

`ss` (Linux) and `netstat` (Windows and others) show sockets. Listen sockets and established tuples tell you whether the process bound the port that you think. `ss -lnt` and `netstat -an` are the teaching commands.

`tcpdump` writes packets on one interface with a capture filter. It is the command-line partner of Wireshark. Topic 1 introduced capture safety.

Order of use for a "cannot connect" report:

1. is the name resolved?
2. does `ping` or a TCP probe reach the address?
3. is anyone listening (`ss` / `netstat`)?
4. what do the packets show (`tcpdump` or Wireshark)?

Do not start with a full packet dump when `ss` already shows no listener. Do not ping only and then stop when the app uses TCP 443 and ICMP is filtered.

`curl -v` sits between these tools and the application. Use it for HTTP and TLS. This section is the lower tools.

### Questions

#### Theoretical questions

1. What does a ping reply prove?
2. What does `ss -lnt` show?
3. Why can traceroute show stars when HTTPS works?
4. What is a capture filter for?
5. Why is "ping failed" not the end of an HTTP debug?

#### Easy practical tasks

1. Write five sentences that assign one question to each of the four tools.
2. Ping your gateway. Write RTT.
3. Run `ss -lnt` or `netstat -an`. Copy two listen lines.
4. Make a table: tool, question it answers. Add four rows.

#### Medium practical tasks

1. Traceroute to a public host. Mark local versus remote hops as you did in topic 1.
2. Write six sentences: no listener on the port versus a filter that drops SYN.
3. Run `tcpdump` or Wireshark with a capture filter for ICMP while you ping. Write the filter text.

#### Advanced practical tasks

1. Write a one-page runbook: four tools in order for "connection refused" versus "timed out."
2. On a local server that you own, show listen, then connect, then capture the handshake. Document commands.

---

## Wireshark display filters

A display filter hides packets in the view. It does not drop them from the file (unless you export the displayed set). A capture filter drops packets before they enter the file. Topic 1 named both. This section builds display-filter skill.

Useful display filters for this path:

- `icmp` or `icmpv6`
- `dns`
- `tcp.port == 443`
- `udp.port == 53`
- `tls`
- `http` (cleartext HTTP)
- `ip.addr == 192.0.2.1` (example address)
- `tcp.flags.syn == 1`

Combine with `and`, `or`, and `not`. Wireshark colors TCP problems. The expert pane hints at retransmits. Do not trust color alone.

Name resolution in Wireshark can send extra DNS. Turn it off when you debug DNS itself.

Do not apply a display filter and then claim that "the packet never happened" if you only hid it. Do not share a file that still contains other conversations.

Save a small pcap with a capture filter when you can. Use a display filter when you already captured too much.

### Questions

#### Theoretical questions

1. What does a display filter change?
2. What does a capture filter change?
3. Why turn off name resolution when you debug DNS?
4. Why can a filter hide a packet that exists in the file?
5. Name two filters that isolate HTTPS TCP.

#### Easy practical tasks

1. Write five sentences that compare display and capture filters.
2. Open a capture. Apply `icmp`. Write how many packets remain visible.
3. Make a table: filter text, purpose. Add four rows from this section.
4. Write a filter for DNS only.

#### Medium practical tasks

1. Capture a `curl` to HTTPS. Apply `tls` and then `tcp.port == 443`. Write whether the counts match.
2. Write six sentences: a teammate who used a display filter and said "DNS never ran."
3. Find the expert or analyze menu for TCP retransmits. Write how many if any.

#### Advanced practical tasks

1. Write a one-page filter cheat sheet for this learning path (ICMP, DNS, TCP handshake, TLS).
2. Export only displayed packets to a new file after a filter. Write the two file sizes.

---

## Packet loss vs delay vs reordering

Loss means a packet does not arrive. Delay means it arrives late. Reordering means a later sequence arrives before an earlier one. The user word "slow" can mean any of the three.

TCP treats loss as a signal to retransmit and often to slow down. Topic 9 covered that. Delay without loss raises RTT and can stall request-response apps. Reordering can look like loss if the stack is eager. UDP apps must handle all three.

How you see them:

- ping: missing replies (loss or filter), RTT variance (delay)
- mtr: percent loss per hop (careful: ICMP loss is not TCP loss)
- Wireshark: retransmission, duplicate ACK, out-of-order
- application: timeout versus long wait with no error

Do not fix "loss" with a larger window when the path is dropping 10 percent. Do not blame the server when the Wi-Fi hop shows delay spikes. Do not treat one `mtr` hop with ICMP loss as a broken user path.

Bufferbloat is extra delay in a full queue. Throughput can stay high. Interactive sessions feel bad. Topic 17 named queues.

A good report states: loss rate, RTT min/avg/max, and whether order broke. Guess after you measure.

### Questions

#### Theoretical questions

1. What is loss?
2. What is delay?
3. What is reordering?
4. Why can mtr hop loss lie about TCP?
5. How can high throughput still feel slow?

#### Easy practical tasks

1. Write five sentences that separate the three words.
2. Ping a public host 20 times. Write min, max, and whether any probe failed.
3. Make a table: symptom, possible cause. Add timeout, high RTT, out-of-order in Wireshark.
4. Draw three timelines: loss, delay, reorder.

#### Medium practical tasks

1. Use `mtr` or WinMTR for one minute to a public host. Write destination loss and last-hop ICMP loss.
2. Write six sentences: Wi-Fi delay spike versus a dead origin.
3. In a capture of a download, search for retransmission. Write a count.

#### Advanced practical tasks

1. Write a one-page measurement sheet: how you record loss, delay, and reorder for a ticket.
2. On a lab that you own, induce delay with a documented tool if you have one, or write a paper plan. Do not attack a shared network.

---

## Application logs + network path

A complete debug uses both ends and the path. The app log shows the error string, the URL, and the time. The network tools show whether the packet left and what came back. Either side alone lies.

Examples:

- app says `connection refused`: look for a listener and a RST
- app says `timed out`: look for dropped SYN or no RST
- app says `certificate verify failed`: look at TLS, not at AAAA
- app says `no such host`: look at DNS, not at TCP
- app is fine, user is not: look at a different resolver, DoH, or a proxy

Align clocks. Topic 15 covered NTP. A 5-minute skew makes two log files useless as a timeline.

Redact secrets. Keep request IDs. One ID in the app and in a proxy log beats a 2 GB pcap.

Do not change three middleboxes at once. Do not ignore the app error code because a ping succeeded. Do not collect only client logs when the server never saw the SYN.

A good ticket has: time (UTC), client IP, dest name and IP, dest port, app error, one ping, one `ss` or listen proof, and an optional small pcap.

### Questions

#### Theoretical questions

1. Why can ping success and app failure both be true?
2. Which tool set matches `no such host`?
3. Why must clocks match on two hosts?
4. What should a good ticket include?
5. Why is a request ID useful across a proxy?

#### Easy practical tasks

1. Write five sentences that pair one app error with one network check.
2. Make a table: app string, first tool. Add four rows from this section.
3. Write a ticket template with the fields this section listed.
4. Check whether your app logs timestamps with a time zone.

#### Medium practical tasks

1. Trigger a refused connection to a closed local port with `curl`. Write the app error and what `ss` showed.
2. Write six sentences: server logs empty, client timeout.
3. Compare a browser error page with `curl -v` for the same URL. Write both texts (redact).

#### Advanced practical tasks

1. Write a one-page correlation guide: proxy `502`, app log, and a pcap of the proxy-to-app hop.
2. Add a request ID to a tiny local server that you own. Show it in the client and in the log.

---

## TLS handshake failures

Topic 13 listed expired cert, wrong host, and untrusted issuer. This section is how you observe them.

`curl -v` prints the check that failed. Browsers use a different sentence for the same check. Read both. `openssl s_client -connect host:443 -servername host` is a common extra tool when you have it. It prints the chain. Do not use it as an attack tool. Use it to view a chain that you may fetch.

A capture shows ClientHello, ServerHello or an alert, and a close. You still cannot read the HTTP path. You can see the failure time relative to TCP. If TCP never completes, the fault is not the certificate.

Split the cases:

- TCP fails: address, port, firewall
- TLS fails: chain, name, clock, protocol version
- HTTP fails: status, app

Corporate TLS intercept changes the issuer. The OS store may trust the intercept CA. `curl` on the same machine may use a different store. That mismatch is a common "works in the browser" bug.

Do not add `-k` to a production script to hide a failure. Do not send a user to "click through." Do not forget IPv6: one family can fail TLS on a different vhost.

### Questions

#### Theoretical questions

1. How do you tell a TCP failure from a TLS failure with `curl -v`?
2. Why can a browser trust a site when `curl` does not?
3. What does a TLS alert in a capture add if `curl` already failed?
4. Why check the clock again here?
5. Why can IPv4 and IPv6 disagree on TLS?

#### Easy practical tasks

1. Write five sentences that split TCP, TLS, and HTTP failures.
2. Run `curl -v https://example.com`. Write that TLS succeeded and the HTTP status.
3. Make a table: failure class, first command. Add three rows.
4. Write the `curl` flag that skips verify and why you must not use it in production.

#### Medium practical tasks

1. Fetch a local self-signed site that you own with `curl` and with a browser. Write both errors. Then fix trust on one client.
2. Write six sentences: intercept CA in the OS store, missing CA in the language runtime.
3. If `openssl s_client` exists, connect to `example.com:443` with SNI. Write the subject. If not, skip.

#### Advanced practical tasks

1. Write a one-page triage tree: handshake error messages to clock, name, issuer, or version.
2. Capture a failed local handshake (wrong name). Label ClientHello and the alert. Use only your machine.

---

## DNS TTL surprises

Topic 12 defined TTL. This section is the debug surprise: you changed a record, some users still see the old address, and some users see the new address. Both can be correct.

Sources of split views:

- recursive caches with remaining TTL
- stub or browser cache
- application process cache that ignores TTL
- split horizon (different answers by design)
- DoH in the browser versus `dig` on the OS resolver
- anycast resolvers that are not in sync for a short time

What you record: query time, resolver address, answer bytes, TTL remaining, and the client network. `dig` two times. Watch TTL decrease. That decrease proves a cache.

Do not lower TTL to 1 second on a huge name in a panic without capacity. Do not flush a public resolver (you cannot). Flush your stub (`ipconfig /flushdns` on Windows, or the OS equivalent) when you test your own host.

Negative cache: a typo can stick as NXDOMAIN. Wait or query a resolver that did not cache the miss.

A CDN cutover plan starts with a TTL drop, a wait of the old TTL, then the change. Topic 12 asked for that plan. Here you use it in an incident.

### Questions

#### Theoretical questions

1. Why can two users see two addresses after a change?
2. What does a decreasing TTL in repeated `dig` prove?
3. Why can a browser and `dig` disagree?
4. What is a negative cache surprise?
5. What must a cutover wait for after you lower TTL?

#### Easy practical tasks

1. Write five sentences that list cache places for a name.
2. Dig a name two times. Write both TTLs.
3. Flush your local DNS cache if you may. Write the command.
4. Make a table: cache place, how you see it.

#### Medium practical tasks

1. Query the same name at two resolvers. Write answers and TTLs.
2. Write six sentences: app still uses an old IP after you flush the OS cache.
3. Document your OS stub flush and browser DoH toggle locations (names only).

#### Advanced practical tasks

1. Write a one-page incident: "half the users on the old IP." Include TTL, resolver, and split horizon.
2. Write a small client that resolves a name in a loop for one minute and prints the address. State whether it honors TTL.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a debug of "the site is slow" that uses ping, a socket list, a display filter, and one app log line.
2. How do loss, a TLS alert, and a DNS TTL split produce three different user stories that all sound like "it is down"?
3. Why must you name the resolver, the clock, and the listen socket before you blame the code?
4. When is a pcap the wrong first artifact, and when is it the right artifact?
5. A teammate pastes a 500 MB capture and no UTC time. Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: four CLI tools, five display filters, loss versus delay, log plus path, TLS split, TTL surprises.
2. Run ping, `ss`/`netstat`, and `curl -v` to one URL. Save three labeled outputs.
3. Draw a flowchart: DNS fail, TCP fail, TLS fail, HTTP fail.
4. List redaction rules for logs and pcaps from this topic.

#### Medium practical tasks

1. From your machine, produce a ticket-quality pack for a successful `https://example.com` fetch: times, IPs, ports, TLS OK, HTTP status.
2. Write a lab plan that creates one refused and one TLS name error on localhost. State what each tool should show.
3. Explain in ten steps how a beginner applies a display filter without deleting the rest of the file.

#### Advanced practical tasks

1. Build a glossary of 16 observability terms from this topic. Each entry: term, one sentence, one command or filter.
2. Write a short playbook: correlate browser DoH, OS `dig`, and a proxy log for one name. No solutions beyond the checks.
