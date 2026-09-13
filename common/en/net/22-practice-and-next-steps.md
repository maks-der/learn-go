# 22. Practice and Next Steps

## Description

This topic turns the path into practice. You draw one HTTPS request with every header layer, you write a TCP echo pair, you write a tiny UDP ping, and you capture HTTP plus a TLS ClientHello. Then you continue with Beej's Guide, a standard textbook, and the RFCs for HTTP and TLS.

Use one term for each concept. A header is the control prefix that a layer adds. Echo means the server returns the bytes that it read. Complete this topic after topics 1 to 21. Do not skip the safety rules: own machine, no floods, no secrets in files that you share.

Practice in this order if you can: addresses and `ping`, `dig`, `curl -v`, echo, capture, then reading. The list at the end of `net.topics.md` matches that order.

---

## Draw a packet from browser to `https://example.com` (every header)

Pick one moment: the first TCP segment that carries a TLS ClientHello, or the first HTTP/2 or HTTP/3 request after the handshake. Name every header that the host adds. Use the TCP/IP map from topic 2.

A complete drawing includes:

- application: HTTP method, path, `Host`, other headers (after TLS decrypt, or as a logical layer)
- TLS record (or QUIC plus TLS) when the scheme is `https`
- TCP (ports, flags, sequence) or UDP plus QUIC
- IPv4 or IPv6 (source, destination, TTL or hop limit)
- Ethernet or Wi-Fi frame (MAC addresses) on the LAN
- the DNS step before the first packet if you draw the full story

Write the 5-tuple. Write whether the path uses NAT (topic 7). Write ALPN if you know the HTTP version (topic 14). Write that a path observer does not see the HTTP path when TLS works.

Do not copy a textbook figure without labels in your own words. Do not skip the link layer because "it is just Ethernet." The switch uses MAC. The router uses IP.

If the browser uses HTTP/3, draw UDP and QUIC. If it uses HTTP/2, draw TCP and TLS. `curl -v` and a capture tell you which stack you used.

### Questions

#### Theoretical questions

1. Which layer adds MAC addresses on a LAN?
2. Which header holds the TCP or UDP ports?
3. Why is HTTP not visible in a normal HTTPS capture?
4. Where does NAT change a header in a home path?
5. What extra box do you add if ALPN selected `h3`?

#### Easy practical tasks

1. Draw IPv4, TCP, TLS, and HTTP as four boxes. Label one field each.
2. Write the 5-tuple for a browser to `example.com` on port 443 (use a real IP from `dig`).
3. Make a table: layer, header example. Add five rows.
4. Write five sentences that include DNS before TCP.

#### Medium practical tasks

1. Run `dig` and `curl -v` to `https://example.com`. Fill the drawing with real addresses, ports, and HTTP version.
2. Add NAT to the drawing: private source, public source. Write which device changes it.
3. Write six sentences: Ethernet hop to the home router, then a WAN header that you do not see.

#### Advanced practical tasks

1. Produce one full-page figure: every header from browser to `example.com` for the stack that you measured. Label "visible on the path" versus "inside TLS."
2. Draw a second figure for HTTP/3 if your browser used `h3`, or write why you still used TCP.

---

## Write a TCP echo server and client

An echo server reads bytes from a TCP connection and writes the same bytes back. An echo client writes a line, then reads the reply. Topics 9 and 10 covered the stream, the handshake, and the socket calls.

Requirements for the exercise:

- bind the server to `127.0.0.1` or to another address that you own
- pick a high port that is free
- handle a partial `read` and a partial `write`
- close the sockets
- print errors

Use Go `net` or Python `socket`. Do not use a framework that hides the stream. The goal is to see the 5-tuple and the stream.

Test with your client and with `curl` or `ncat` if you want a second client. Capture the handshake on loopback if your OS allows it.

Do not bind `0.0.0.0` on a shared network for this exercise. Do not echo on a privileged port. Do not send secrets into echo. Do not treat echo as a production protocol.

When it works, add a timeout (topic 10). Then write one sentence about `TIME_WAIT` after you stop the client.

### Questions

#### Theoretical questions

1. Why can one `read` return fewer bytes than you wrote?
2. Why bind to localhost for this lab?
3. What does the server write back?
4. Which topic listed `socket`, `bind`, `listen`, `accept`, `connect`?
5. Why is echo a poor production API?

#### Easy practical tasks

1. Write five sentences that describe echo as a stream, not as messages.
2. Pick a port and write the server bind address.
3. Make a table: call, role (server or client). Add `listen`, `accept`, `connect`.
4. Draw: client, loopback, server, bytes there and back.

#### Medium practical tasks

1. Implement the pair in one language. Send one line. Save the output.
2. Send a payload larger than one system `write`. Show that you looped until all bytes returned.
3. Write six sentences: blocking `read` versus a timeout.

#### Advanced practical tasks

1. Add a second concurrent client if your server accepts a loop of `accept`. Write what you saw.
2. Capture the handshake on loopback. Write client port, server port, and the three handshake flags.

---

## Write a tiny UDP ping

A tiny UDP ping sends a datagram with a sequence number and a time (or a token). The peer returns the same payload. The sender measures RTT when the matching token returns. Topic 8 covered loss, order, and checksums.

Requirements:

- one socket on each side
- a timeout on the sender so that a loss does not hang forever
- a printed result: reply, timeout, or unexpected payload
- localhost first, then a second device on a LAN that you own if you want

This is not ICMP `ping`. This is your protocol on UDP. Firewalls treat it as UDP, not as echo request.

Do not send a tight loop of datagrams to a public host. Do not claim that you measured path MTU with this tool unless you sized the payload on purpose (topic 19). Do not skip the timeout.

If you never get a reply on a real LAN, check bind address, port, and the host firewall (topic 21).

Compare the exercise with DNS: short request, short reply, retry at the application.

### Questions

#### Theoretical questions

1. Why must the sender use a timeout?
2. How do you match a reply to a request?
3. Why is this tool not ICMP ping?
4. What does UDP not fix for you?
5. Why can a LAN firewall drop your datagram when `ping` still works?

#### Easy practical tasks

1. Write five sentences that describe the payload and the reply.
2. Choose a port and a payload layout (seq plus token). Write it.
3. Make a table: ICMP ping, your UDP ping. Add one row for protocol.
4. Draw two datagrams: request and reply.

#### Medium practical tasks

1. Implement sender and responder on localhost. Record one RTT and one timeout (stop the responder to force a timeout).
2. Send three datagrams. Write whether any arrived out of order (localhost may never reorder).
3. Write six sentences: application retry versus TCP retransmission.

#### Advanced practical tasks

1. Run the responder on a second device on a LAN that you own. Record RTT. Stop if a policy forbids it.
2. Write a one-page comparison: your UDP ping, ICMP ping, and `curl` time to first byte.

---

## Capture your own HTTP and TLS ClientHello

Use a network that you own. Start Wireshark or `tcpdump`. Fetch `http://example.com` if the host still answers HTTP (you may see a redirect). Fetch `https://example.com`. Stop the capture. Save a small file.

Find:

- DNS query and answer if they are cleartext
- TCP handshake to port 80 or 443
- HTTP request line if the session is cleartext
- TLS ClientHello: destination port, SNI if visible, advertised versions
- ALPN list if the tool shows it

Write what you cannot read: HTTP path on HTTPS, cookies on HTTPS, the certificate names inside encrypted handshake parts that your TLS version hides.

Do not capture on a work VPN if policy forbids it. Do not upload the pcap to a public site. Do not include other conversations. Use a capture filter such as `host` of the resolved address.

Compare the file with `curl -v` from the same fetch. The two views should agree on ports and on success or failure.

If HTTP/3 appears, find UDP and QUIC. Name it. That is a valid result.

### Questions

#### Theoretical questions

1. What is a ClientHello?
2. Why can SNI appear when HTTP stays hidden?
3. What capture filter keeps the file small?
4. Why must you not publish the pcap?
5. How does `curl -v` complement the capture?

#### Easy practical tasks

1. Write five rules for a safe capture (topic 1 plus this section).
2. Start a capture, ping the gateway, stop. Confirm that you can save a file.
3. Make a table: visible field, hidden on HTTPS. Add four rows.
4. Write a display filter for `tls` and one for `http`.

#### Medium practical tasks

1. Capture `curl` to `https://example.com`. Write ClientHello time relative to the SYN.
2. Capture an HTTP fetch or a redirect. Write the first request line if visible.
3. Write six sentences: what you would tell a classmate who expects to read the HTTPS path.

#### Advanced practical tasks

1. Label a screenshot or a note: DNS, SYN, ClientHello, application data. Mark each as named in topics 12 to 14.
2. Write a short lab guide for a classmate: filters, what to click, what not to share. No solutions to the drawing exercise.

---

## Then: Beej's Guide, Computer Networking: A Top-Down Approach, RFCs for HTTP and TLS

After the labs, read. Do not read everything at once.

[Beej's Guide to Network Programming](https://beej.us/guide/bgnet/) teaches sockets in C with a clear voice. Use it to compare with Go `net` and Python `socket`. The calls match topic 10.

*Computer Networking: A Top-Down Approach* (Kurose and Ross) is a common textbook. It follows application to link, which matches how you used HTTP first. Use it for depth on congestion control, DNS, and HTTP. Do not copy pages into this handbook.

RFCs are the protocol standards:

- HTTP: start with the current HTTP semantics and HTTP/1.1 message RFCs (the IETF HTTP working group documents; numbers change when they revise). [MDN HTTP](https://developer.mozilla.org/en-US/docs/Web/HTTP) is easier for daily headers.
- TLS: start with TLS 1.3 (RFC 8446) after you know the handshake story. Skim. Do not memorize cipher suites first.
- TCP: RFC 9293 is the current TCP spec. Topic list item 10 in the suggested practice order asked you to skim the state machine.

Other official links from the path: [IETF RFC editor](https://www.rfc-editor.org/), [Cloudflare Learning Center](https://www.cloudflare.com/learning/), [Wireshark User's Guide](https://www.wireshark.org/docs/).

Read with a purpose. One RFC section plus one capture beats a full RFC in a weekend.

Do not treat a blog as a substitute for the RFC when you implement a parser. Do not paste large RFC text into notes that you publish. Summarize.

Pair the next OS topics for kernel sockets and the file descriptor model.

### Questions

#### Theoretical questions

1. What does Beej's Guide add that this path only surveyed?
2. Why is a top-down textbook a fit after HTTP-first practice?
3. Why skim RFC 8446 after topic 13, not before topic 1?
4. What is a better first HTTP read than a full RFC for a beginner?
5. Why summarize an RFC instead of copying it?

#### Easy practical tasks

1. Open Beej's Guide. Write the chapter titles for sockets and for `select` or equivalent.
2. Open RFC 8446. Write the title and the year. Do not copy the abstract in full.
3. Make a table: resource, when you will read it. Add four rows.
4. Bookmark MDN HTTP and the RFC editor. Write the two URLs.

#### Medium practical tasks

1. Read Beej on `send` and `recv`. Write six sentences that match your echo lab.
2. Skim the TCP state machine in RFC 9293. Write the state names that you already used (LISTEN, ESTABLISHED, TIME-WAIT).
3. Read MDN on one header (`Host` or `Cache-Control`). Write four STE sentences.

#### Advanced practical tasks

1. Write a one-page reading plan for the next month: Beej, one textbook chapter, one RFC section, one capture per week.
2. Compare one textbook figure of HTTPS with your drawing from the first section. Write three differences.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the practice path from a drawing of `https://example.com` to echo, UDP ping, a capture, and a reading list. Name one topic that each lab reuses.
2. How do TCP echo and UDP ping together prove that you understood streams versus datagrams?
3. Why does a ClientHello capture belong with the drawing exercise, not instead of it?
4. When should you read an RFC, and when should you stay with MDN or Beej?
5. A teammate wants to skip the labs and only read Kurose. Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page checklist of the four labs and the three reading sources.
2. Confirm the tools from the path: `ping`, `dig` or `nslookup`, `curl`, Wireshark or `tcpdump`, a compiler or interpreter for echo.
3. Draw a week plan: two labs, one capture review, one reading block.
4. List files you will keep private (pcaps, keys, verbose logs).

#### Medium practical tasks

1. Complete the four labs on a machine that you own. Write a status line for each: done or blocked, and why.
2. Write a lab report template: command, output path, figure name, safety note. Leave the answers blank.
3. Explain in ten steps how a beginner goes from `ipconfig`/`ip addr` to a labeled ClientHello without publishing secrets.

#### Advanced practical tasks

1. Build a personal glossary of 20 terms from topics 12 to 22. Each entry: term, one sentence, one lab or tool.
2. Write a next-quarter plan: one small public RFC you will skim, one socket program you will extend, and one OS topic you will pair. No solutions to the labs in this file.
