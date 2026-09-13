# 11. HTTP/2, HTTP/3, and Other Protocols

## Description

HTTP/2 and HTTP/3 are newer HTTP mappings. HTTP/2 uses TLS on TCP in the common web case. HTTP/3 uses QUIC on UDP. This topic shows multiplexing and why HTTP/3 exists. You learn QUIC over UDP. You get a survey of mail, WebSockets, gRPC, and MQTT.

Complete this topic after HTTP and TLS.

Use one term for each concept. A stream is one request-response conversation inside one connection. Multiplexing means many streams share one connection. QUIC is a transport on UDP that HTTP/3 uses. A survey names the job of each protocol. It does not replace a full mail or messaging handbook. Do not mix a stream with a TCP connection. Do not mix SFTP with FTP.

Practice with `curl -v` and a browser developer view that shows the HTTP version. A capture of HTTP/2 or HTTP/3 is binary. Name the version first.

---

## Multiplexing and why HTTP/3 exists

HTTP/1.1 keep-alive reuses one TCP connection. The usual client still waits for one response before it sends the next request on that connection. Browsers open several TCP connections to reduce the wait.

HTTP/2 multiplexes many streams on one connection. Each stream has an identifier. Frames carry stream data. Many requests can be in flight. Priority and flow control exist per stream.

One lost TCP segment still stalls the whole connection at the TCP layer. Multiplexing does not remove TCP head-of-line blocking. Topic 6 named that limit.

HTTP/3 multiplexes streams on QUIC. A loss on one stream does not block other streams at the transport layer in the QUIC design. That is a reason HTTP/3 exists. Lossy radio links and many small objects make the difference visible.

The server can push in HTTP/2 in the old design. Many browsers limit or ignore push. Do not build a new service on server push. Use a normal request.

ALPN selects `h2` or `http/1.1` on TLS. HTTP/3 uses a different setup (HTTPS records and Alt-Svc or prior knowledge). `curl -v` often prints the HTTP version.

Do not open hundreds of HTTP/2 connections to go faster without a measurement. One connection per origin is the teaching default. Do not confuse a stream with a TCP connection.

Header compression (HPACK on HTTP/2, QPACK on HTTP/3) reduces repeat headers. You do not implement those codecs here. Encrypted payloads still need keys that you own.

### Questions

#### Theoretical questions

1. What does multiplexing mean for HTTP/2?
2. How does HTTP/1.1 keep-alive differ from HTTP/2 multiplexing in the usual client?
3. Why does HTTP/2 still have TCP head-of-line blocking?
4. Why does HTTP/3 exist in this section?
5. Why is server push a poor base for a new design?

#### Easy practical tasks

1. Write five sentences that compare one-at-a-time HTTP/1.1 and multiplexed HTTP/2.
2. Draw one connection with three streams.
3. Make a table: HTTP/1.1, HTTP/2, HTTP/3. Add one row for transport and one row for multiplexing.
4. Write four sentences: stream versus TCP connection.

#### Medium practical tasks

1. Run `curl -v https://example.com`. Write the HTTP version that `curl` reports.
2. Write six sentences: a page with twenty small files on HTTP/1.1 versus HTTP/2.
3. In a browser network panel, load a public site. Write whether the panel shows h2 or h3.

#### Advanced practical tasks

1. Force `curl --http1.1` and default `curl` on the same HTTPS URL. Write both versions and both total times for a small page.
2. Write a one-page note: what multiplexing helps and what TCP loss still hurts on HTTP/2.

---

## QUIC over UDP

QUIC is a transport protocol that runs over UDP. HTTP/3 maps HTTP semantics onto QUIC streams. QUIC includes encryption in the transport. You do not run cleartext HTTP/3 on the public web.

QUIC combines handshake ideas so that a new session can take fewer round trips than TCP plus TLS on a cold start. Details depend on version and on resumption. This topic needs the idea: QUIC sits where TCP plus TLS sat.

QUIC ports are UDP. Firewalls and NAT that allow only TCP 443 can block HTTP/3. The client then falls back to HTTP/2 or HTTP/1.1 on TCP. A UDP idle timeout can kill a quiet QUIC path.

Because QUIC is UDP, a capture filter for `tcp.port == 443` misses HTTP/3. Use UDP and the port that the client uses (often 443). The payload is not HTTP/1.1 text.

QUIC has its own congestion control and retransmission per stream. You do not tune it here.

Do not disable UDP 443 on a lab and then conclude HTTP is dead. Check the fallback. Do not send your own large custom UDP protocol and call it QUIC. QUIC is a standard with a handshake and crypto.

Topic 5 said UDP has no reliability. QUIC adds reliability on top of UDP. The UDP header is still there.

### Questions

#### Theoretical questions

1. What transport does QUIC use?
2. What application protocol uses QUIC in this topic?
3. Why can a TCP-only firewall block HTTP/3?
4. Why does a TCP port-443 display filter miss HTTP/3?
5. How does QUIC change the "UDP has no reliability" fact?

#### Easy practical tasks

1. Write five sentences that place QUIC under HTTP/3.
2. Make a table: HTTPS HTTP/2, HTTPS HTTP/3. Add TCP or UDP and TLS or built-in crypto.
3. Draw: UDP, QUIC, HTTP/3.
4. Write four sentences: fallback to TCP when UDP is blocked.

#### Medium practical tasks

1. In a browser, see if a site uses h3. Then block UDP in a lab that you own if you can, or write the expected fallback.
2. Write six sentences: NAT UDP timeout and a long-lived HTTP/3 tab.
3. Run `curl --http3` if your build supports it. Write success or "not built."

#### Advanced practical tasks

1. Capture UDP 443 while you load a site that uses HTTP/3. Write whether you see UDP. Use a network that you own.
2. Write a one-page note: QUIC handshake versus TCP plus TLS. High-level only.

---

## Mail, WebSockets, gRPC, MQTT (survey)

Many applications speak a protocol above TCP or UDP. This survey names the job of each protocol.

Mail: SMTP moves a message toward a mailbox. The common submission port is 587 with TLS, or 465 for implicit TLS. Port 25 is the hop between mail servers. Many ISPs block port 25 from homes. IMAP is the protocol that a client uses to read and manage messages. The common port is 993 with TLS. A message has headers and a body. SPF, DKIM, and DMARC are DNS and crypto checks on the sender domain. They are not SMTP commands. Mail is store-and-forward. Delay can be minutes. Debug with logs and with `dig` on MX, not only with `curl`.

WebSockets: a client starts with HTTP, then upgrades to a long-lived bidirectional message channel on the same TCP (or QUIC) session. Browsers use WebSockets for live updates. The protocol is not request-response after the upgrade. You still need TLS (`wss://`) on a public path.

gRPC: remote procedure calls that use HTTP/2 (and HTTP/3 in newer stacks). Messages are often Protobuf, not JSON. You need HTTP/2 features such as streams. An L7 proxy must understand gRPC or you use an L4 pass-through. Deadlines and status codes have gRPC names plus HTTP mapping.

MQTT: a publish and subscribe protocol for devices and brokers. It uses TCP in the common case. MQTT over TLS exists. Topics are broker strings, not HTTP paths. QoS levels control retry. A retained message is not HTTP cache. Do not open an MQTT broker to the Internet without auth.

Other names that you will meet: SFTP is file transfer on SSH, not FTP. FTP uses two connections and often cleartext. Prefer SFTP or HTTPS upload. Redis and Postgres have their own TCP protocols. Treat each as a framed stream.

Do not send SMTP without TLS on a public path. Do not treat the From header as proof of the sender. Do not run an open mail relay or an open MQTT broker.

### Questions

#### Theoretical questions

1. What job does SMTP do?
2. What job does IMAP do?
3. What does a WebSocket upgrade change?
4. Why does gRPC usually need HTTP/2?
5. What pattern does MQTT use?

#### Easy practical tasks

1. Write five sentences that assign one job to mail, WebSockets, gRPC, and MQTT.
2. Make a table: protocol, typical transport, one port or URL scheme.
3. Query MX for a public domain. Write the exchanger names.
4. Draw: HTTP upgrade to WebSocket as two phases.

#### Medium practical tasks

1. Open mail settings on an account that you own. Write SMTP host, IMAP host, and TLS on or off. Redact the username if you save the note.
2. Write six sentences: port 25 versus port 587.
3. Find whether a public API you use offers gRPC or only JSON HTTP. Write the evidence.

#### Advanced practical tasks

1. Write a one-page mail path: compose, SMTP, MX, store, IMAP. Name TLS at each hop that must have it.
2. Write a survey table: WebSocket versus gRPC stream versus MQTT topic. Add one fit each.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do HTTP/2 multiplexing, TCP HOL, and HTTP/3 on QUIC form one story?
2. Why must a firewall story name UDP when the user says HTTPS works in some browsers only?
3. When do you pick WebSockets, gRPC, or MQTT instead of REST HTTP/1.1?
4. How does mail store-and-forward differ from an HTTP response time?
5. What does ALPN have to do with seeing `h2` in `curl -v`?

#### Easy practical tasks

1. Write a cheat sheet: multiplex, HOL, QUIC, UDP 443, SMTP, IMAP, WebSocket, gRPC, MQTT.
2. Run `curl -v https://example.com` and write the HTTP version.
3. Draw HTTP/1.1, HTTP/2, and HTTP/3 as three stacks.
4. List four protocols from this topic that need TLS or SSH on a public path.

#### Medium practical tasks

1. In a browser network panel, collect HTTP versions for five requests on one page.
2. Write a fault tree: site works on TCP 443, fails on HTTP/3 only. Include UDP and NAT.
3. Write six sentences: gRPC through an L7 load balancer that only knows HTTP/1.1.

#### Advanced practical tasks

1. Write a one-page protocol pick list for a new product: public website, live UI, device telemetry, mail notify.
2. Capture or log ALPN on a lab server that you own. Show `h2` versus `http/1.1`.
