# 16. Application Protocols (Survey)

## Description

Many applications speak a protocol above TCP or UDP. This topic is a survey: mail (SMTP and IMAP), FTP versus SFTP, WebSockets, gRPC on HTTP/2, MQTT, and a short view of Redis and Postgres wire protocols.

Use one term for each concept. A survey names the job of each protocol. It does not replace a full mail or database handbook. Complete this topic after HTTP, TLS, and sockets.

Practice with public documentation and with tools that you already have (`curl`, a mail client, a database client). Do not attack a mail server or a database that you do not own.

---

## SMTP / IMAP (mail)

SMTP is the protocol that moves a message toward a mailbox. A mail client or a mail server uses SMTP to submit or to relay. The common submission port on the Internet today is 587 with TLS, or 465 for implicit TLS. Port 25 is still the hop between mail servers. Many ISPs block port 25 from homes.

IMAP is the protocol that a client uses to read and manage messages on a server. The common port is 993 with TLS. POP3 is an older fetch protocol. This survey uses IMAP.

A message has headers and a body. MIME can add parts. SPF, DKIM, and DMARC are DNS and crypto checks on the sender domain. They are not SMTP commands. You already met TXT records in the DNS topic.

Do not send SMTP without TLS on a public path. Do not treat "the mail arrived" as proof that the From header is true. Do not run an open relay.

Mail is store-and-forward. Delay can be minutes. That is not HTTP request-response. Debug with logs and with `dig` on MX, not only with `curl`.

### Questions

#### Theoretical questions

1. What job does SMTP do?
2. What job does IMAP do?
3. Why do many homes fail to send on port 25?
4. Is a From header proof of the sender?
5. How is mail delay different from an HTTP response?

#### Easy practical tasks

1. Write five sentences that separate SMTP and IMAP.
2. Make a table: protocol, typical port with TLS, direction (send or read).
3. Query MX for a public domain. Write the exchanger names.
4. Draw: client, submission, MX hop, IMAP read.

#### Medium practical tasks

1. Open a mail account settings page that you own. Write SMTP host, IMAP host, and TLS on or off. Redact the username if you save the note.
2. Write six sentences: port 25 versus port 587.
3. Find a TXT SPF or DMARC record on a public domain. Write whether it exists.

#### Advanced practical tasks

1. Write a one-page mail path: compose, SMTP, MX, store, IMAP. Name TLS at each hop that should have it.
2. Use `dig` or `nslookup` to collect MX, A, and AAAA for one mail domain. Put the answers in a table.

---

## FTP vs SFTP

FTP is an old file-transfer protocol. Classic FTP uses a control connection and a separate data connection. NAT and firewalls break active FTP. Passive FTP is the usual workaround. FTP often sends passwords and files in cleartext. FTPS adds TLS. The two connections still add pain.

SFTP is file transfer on an SSH channel. SFTP is not FTP and is not FTPS. One TCP connection (SSH) carries commands and data. Firewalls see port 22 in the common case. Topic 13 already introduced SSH.

Choose SFTP or another SSH-based copy (`scp`) when you control both ends and you want a simple secure file move. Choose HTTPS upload or a signed object store for many web apps. Choose FTPS only when a vendor still requires FTP-like commands.

Do not enable cleartext FTP on the Internet. Do not confuse SFTP with FTPS in a runbook. Do not open a wide data-port range for FTP unless you must and you own the firewall.

A packet capture of FTP shows USER lines. A capture of SFTP shows SSH. That difference is a learning tool.

### Questions

#### Theoretical questions

1. Why does classic FTP use two connections?
2. Why is cleartext FTP a problem on a public path?
3. How is SFTP related to SSH?
4. Is SFTP the same as FTPS?
5. Why does NAT hurt active FTP?

#### Easy practical tasks

1. Write five sentences that compare FTP and SFTP.
2. Make a table: FTP, FTPS, SFTP. Add one row for encryption and one row for connections.
3. Write four sentences: passive FTP as a NAT workaround (idea).
4. Draw: FTP control plus data versus one SSH session.

#### Medium practical tasks

1. If you have `sftp` or an SSH client, read `sftp --help` or the man page. Write three command names. Do not connect to a host that you must not use.
2. Write six sentences: why a web product should prefer HTTPS upload over FTP.
3. Find whether a vendor you know still documents FTP. Write the protocol they recommend.

#### Advanced practical tasks

1. On a lab VM that you own, copy one file with SFTP or `scp`. Write the port that you used.
2. Write a one-page decision sheet: FTP, FTPS, SFTP, HTTPS. One sentence each for "when".

---

## WebSockets

A WebSocket is a long-lived, bidirectional message channel that starts as HTTP. The client sends an HTTP request with `Upgrade: websocket`. The server answers `101` if it accepts. Then both sides send frames. The transport is still TCP (or a proxy path). TLS uses `wss://`. Cleartext uses `ws://`.

WebSockets fit chat, live views, and games that need server-push without a new HTTP request each time. They do not replace REST for simple CRUD. They do not replace QUIC by themselves. HTTP/2 and SSE are other push styles.

Proxies must allow the upgrade and must not treat the session as a short HTTP request. Idle timeouts close quiet sockets. Heartbeats keep the path alive.

Do not send secrets on `ws://` on a public LAN. Do not ignore `Origin` checks on a server that you write. Do not assume that a WebSocket is a full mesh network.

A capture after the upgrade is not HTTP/1.1 text. Wireshark can decode WebSocket frames when the session is cleartext or when you have keys.

### Questions

#### Theoretical questions

1. How does a WebSocket start?
2. What status code means the upgrade succeeded?
3. What is the TLS URL scheme for WebSockets?
4. Why do idle timeouts matter?
5. Why is `ws://` a poor choice for secrets on a public path?

#### Easy practical tasks

1. Write five sentences that explain upgrade, then frames.
2. Make a table: HTTP request-response, WebSocket. Add one row for who can send next.
3. Write four sentences: `ws` versus `wss`.
4. Draw: HTTP request, 101, then two arrows.

#### Medium practical tasks

1. Find a public demo or a product that uses WebSockets. Write `ws` or `wss` from the docs. Do not abuse the demo.
2. Write six sentences: a reverse proxy that has a 30-second idle timeout.
3. Compare WebSockets and HTTP/2 multiplexing in a short paragraph (idea).

#### Advanced practical tasks

1. Write a tiny local WebSocket echo in a language that you know, on a machine that you own. Document how you tested it.
2. Write a one-page note: when to use WebSockets versus periodic HTTP poll.

---

## gRPC / HTTP/2

gRPC is an RPC framework. The default mapping uses HTTP/2. Each call can be a stream. Protobuf is the common payload format. The client calls a method on a stub. The server implements the method.

gRPC needs HTTP/2 (or an HTTP/3 mapping in newer stacks). A proxy that only understands HTTP/1.1 text can break gRPC. You need an HTTP/2-aware load balancer or a pass-through of HTTP/2.

TLS and ALPN (`h2`) apply as in topic 14. Many gRPC services use TLS. Some lab setups use cleartext HTTP/2 (`h2c`). Do not use cleartext on a public path.

gRPC is a good fit for service-to-service calls with typed contracts. It is a weaker fit for a simple public browser API unless you generate clients and you control the browsers.

Do not debug gRPC as if it were HTTP/1.1 `curl` text. Use `grpcurl` or a language client. Some `curl` builds can speak HTTP/2 with a binary body. Read the tool docs.

Deadlines and status codes in gRPC are not the same numbers as HTTP status, but a gateway can map them.

### Questions

#### Theoretical questions

1. What HTTP version does common gRPC use?
2. What payload format is common in gRPC?
3. Why can an HTTP/1.1-only proxy break gRPC?
4. What ALPN ID does HTTP/2 gRPC need on TLS?
5. Why is gRPC a poor default for an anonymous public browser form?

#### Easy practical tasks

1. Write five sentences that define gRPC as RPC on HTTP/2.
2. Make a table: REST JSON on HTTP/1.1, gRPC. Add rows for contract and stream support (idea).
3. Write four sentences: protobuf versus JSON as a beginner view.
4. Draw: client stub, HTTP/2, server method.

#### Medium practical tasks

1. Read a public gRPC "hello world" page. Write the service name and one RPC name. No long copy.
2. Write six sentences: L7 load balancer that terminates HTTP/1.1 only.
3. Find whether `grpcurl` or a similar tool exists for your OS. Write yes or no.

#### Advanced practical tasks

1. Run a local gRPC example from an official quickstart on a machine that you own. Write the port and whether TLS was on.
2. Write a one-page note: when gRPC helps and when a simple HTTP JSON API is enough.

---

## MQTT (IoT)

MQTT is a publish-subscribe protocol. A broker sits in the middle. Publishers send messages to topics. Subscribers receive messages for topics that they chose. Clients are often sensors and phones. The transport is usually TCP. TLS is the secure variant. Port 1883 is common for cleartext. Port 8883 is common for TLS.

QoS levels control confirmations. QoS 0 is at most once. QoS 1 is at least once. QoS 2 is exactly once in the MQTT sense. Higher QoS uses more chatter. A constrained radio may stay at QoS 0.

MQTT is a good fit for many-to-many telemetry. It is a poor fit for a large file transfer. It is not a replacement for HTTP in a browser without extra libraries.

Do not expose an MQTT broker to the Internet without auth and TLS. Do not use guessable topic names as the only access control. Do not flood a broker from a script on a shared lab.

A retained message is a last-value store on a topic. A will message is a last word when a client dies. Those features are part of the survey, not a full IoT course.

### Questions

#### Theoretical questions

1. What role does a broker play?
2. What does a topic name select?
3. What is the difference between QoS 0 and QoS 1 at a high level?
4. Which port is common for MQTT over TLS?
5. Why are topic names a poor only access control?

#### Easy practical tasks

1. Write five sentences that explain publish, subscribe, and broker.
2. Make a table: HTTP GET, MQTT subscribe. Add one row for who starts the data.
3. Write four sentences: a temperature sensor and a dashboard.
4. Draw: two publishers, one broker, two subscribers.

#### Medium practical tasks

1. Read a public MQTT primer. Write four STE sentences on QoS. No long quotes.
2. Write six sentences: MQTT on a home LAN versus MQTT on the public Internet.
3. Find one public cloud IoT MQTT endpoint description (host and TLS). Write the TLS requirement only.

#### Advanced practical tasks

1. Run a local broker in a container or VM that you own, or write why you skipped and what you would publish.
2. Write a one-page IoT threat note: open broker, no TLS. Hardening ideas only. No attack steps.

---

## Redis / Postgres wire protocols (awareness)

Redis and PostgreSQL are not HTTP. Each has a wire protocol on TCP. A client library speaks that protocol. You already use sockets. This section is awareness so that a capture does not surprise you.

Redis uses a simple text-oriented protocol (RESP). Commands look like `*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n` on the wire in the classic form. AUTH can send a password. Do not expose Redis to the Internet. Use a password or better auth, bind to localhost or a private net, and use TLS when the path is not private.

PostgreSQL starts with a startup message, then authentication, then query messages. The protocol is binary after startup in modern versions. `sslmode` controls TLS. Do not put a Postgres port on `0.0.0.0` without a firewall and auth.

Awareness facts:

- the port numbers are well known (6379, 5432) but you can change them
- a proxy or pooler can sit in front (PgBouncer, Redis replica)
- application logs plus `ss` tell you more than a guess at the payload
- ORMs still use these protocols under the HTTP API of your app

Do not paste a production capture that contains queries or keys. Do not treat "the port is filtered" as proof that the data is safe if a client on the LAN can connect.

You do not implement these protocols here. You name them and you keep them off public interfaces.

### Questions

#### Theoretical questions

1. Does a typical Redis client speak HTTP?
2. Why is an open Redis port on the Internet a fault?
3. What does Postgres `sslmode` control at a high level?
4. Why can a pooler change what you see in a capture?
5. Why must you redact database captures?

#### Easy practical tasks

1. Write five sentences that separate HTTP APIs from database wire protocols.
2. Make a table: Redis, Postgres. Add default port and "TLS available".
3. Write four sentences: bind to localhost versus bind to all addresses.
4. Draw: app process, private TCP, database.

#### Medium practical tasks

1. If you have Redis or Postgres locally, run `ss -lnt` or `netstat` and write the listen address. If not, write the default ports from this section.
2. Write six sentences: an app HTTP port on the Internet and a database port only on a private net.
3. Read one official page on Redis security or Postgres host-based auth. Write four STE sentences. No long quotes.

#### Advanced practical tasks

1. Write a one-page awareness note: how you would recognize Redis versus HTTP in a cleartext lab capture (idea only).
2. Design a listen and firewall plan for a laptop lab: app, Redis, Postgres. Use only addresses that you own.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe how a user reads mail and how a sensor publishes a temperature. Name the protocols and whether each is request-response or pub-sub.
2. How do TLS and SSH appear across SMTP, SFTP, `wss`, gRPC, and MQTT in this survey?
3. Why do FTP and gRPC both care about middleboxes, but for different reasons?
4. When is HTTP the wrong protocol in this list, even if you can tunnel anything over HTTP?
5. A teammate wants "one protocol for mail, chat, IoT, and SQL." Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: SMTP, IMAP, FTP, SFTP, WebSockets, gRPC, MQTT, Redis, Postgres. One line each.
2. Make a table: protocol, TCP or UDP, TLS usual port if this topic named one.
3. Draw two stacks: browser HTTPS and device MQTT TLS.
4. List which of these programs exist on your machine (mail client, `psql`, `redis-cli`, `ssh`).

#### Medium practical tasks

1. Pick three protocols from this topic. For each, write how you would test a local instance that you own.
2. Write a lab plan: capture SFTP versus a failed FTP idea (do not run cleartext FTP on a public net).
3. Explain in ten steps how you choose a protocol for "live dashboard of sensor data."

#### Advanced practical tasks

1. Build a glossary of 15 protocol names and features from this topic. Each entry: term, one sentence, one port or URL scheme if any.
2. Write a short architecture: mobile app, API (HTTPS or gRPC), broker (MQTT), database (Postgres). Label each hop TLS or private.
