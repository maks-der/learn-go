# 14. HTTP/2, HTTP/3, QUIC

## Description

HTTP/2 and HTTP/3 are newer HTTP mappings. HTTP/2 uses TLS on TCP in the common web case. HTTP/3 uses QUIC on UDP. This topic shows multiplexing, header compression at an awareness level, QUIC on UDP, why HTTP/3 exists, and ALPN.

Use one term for each concept. A stream is one request-response conversation inside one connection. Multiplexing means many streams share one connection. Complete this topic after TLS and HTTP/1.1.

Practice with `curl -v` and a browser developer view that shows the HTTP version. A capture of HTTP/2 or HTTP/3 is binary. Name the version first. Do not expect HTTP/1.1 text on the wire.

---

## Multiplexing

HTTP/1.1 keep-alive reuses one TCP connection. The usual client still waits for one response before it sends the next request on that connection. Browsers open several TCP connections to reduce the wait.

HTTP/2 multiplexes many streams on one connection. Each stream has an identifier. Frames carry stream data. One lost TCP segment still stalls the whole connection at the TCP layer. Multiplexing does not remove TCP head-of-line blocking. The TCP topic already named that limit.

HTTP/3 multiplexes streams on QUIC. A loss on one stream does not block other streams at the transport layer in the QUIC design. That is a reason HTTP/3 exists. A later section repeats this with loss.

The server can push in HTTP/2 in the old design. Many browsers limit or ignore push. Do not build a new service on server push. Use a normal request.

Do not open hundreds of HTTP/2 connections to "go faster" without a measurement. One connection per origin is the teaching default. Do not confuse a stream with a TCP connection.

Priority and flow control exist per stream. This topic only needs: many requests can be in flight on one connection.

### Questions

#### Theoretical questions

1. What does multiplexing mean for HTTP/2?
2. How does HTTP/1.1 keep-alive differ from HTTP/2 multiplexing in the usual client?
3. Why does HTTP/2 still have TCP head-of-line blocking?
4. What identifies one HTTP/2 stream?
5. Why is server push a poor base for a new design?

#### Easy practical tasks

1. Write five sentences that compare one-at-a-time HTTP/1.1 and multiplexed HTTP/2.
2. Draw one connection with three streams.
3. Make a table: HTTP/1.1, HTTP/2, HTTP/3. Add one row for transport and one row for multiplexing.
4. Write four sentences: stream versus TCP connection.

#### Medium practical tasks

1. Run `curl -v https://example.com`. Write the HTTP version that `curl` reports.
2. Write six sentences: a page with twenty small files on HTTP/1.1 versus HTTP/2 (idea).
3. In a browser network panel, load a public site. Write whether the panel shows h2 or h3.

#### Advanced practical tasks

1. Force `curl --http1.1` and default `curl` on the same HTTPS URL. Write both versions and both total times for a small page.
2. Write a one-page note: what multiplexing helps and what TCP loss still hurts on HTTP/2.

---

## HPACK / QPACK (awareness)

HTTP/1.1 repeats headers as text on every request. `Cookie`, `User-Agent`, and `Authorization` can be large. HTTP/2 uses HPACK to compress headers. HTTP/3 uses QPACK.

Both schemes use tables of header fields. The two sides agree on names and values that they already sent. The wire then carries a compact form. You do not need the Huffman bit layout for this path.

Awareness facts:

- compression reduces repeat header bytes
- HPACK assumes ordered delivery (TCP)
- QPACK is built for QUIC streams that can arrive out of order
- compressed headers are still inside TLS or QUIC crypto on the public web

Header compression had security stories when it mixed secrets with attacker-controlled fields. Modern stacks follow current RFC rules. Do not invent your own header compression.

Do not look at a binary capture and expect to read `Host:` as ASCII. The decoder in Wireshark can show headers after it parses HPACK. Encrypted payloads still need keys that you own.

This section is awareness. You must know the names HPACK and QPACK and why they exist. You do not implement them here.

### Questions

#### Theoretical questions

1. Why does HTTP/1.1 waste bytes on headers?
2. What does HPACK compress?
3. Why does HTTP/3 use QPACK instead of HPACK?
4. Are compressed headers cleartext on the public Internet path?
5. Why must you not invent your own header compression?

#### Easy practical tasks

1. Write five sentences that explain header compression as "do not send the same cookie text every time."
2. Make a table: HPACK, QPACK. Add one row for HTTP version and one row for transport assumption.
3. List four headers that repeat on every request to one site.
4. Write four sentences: Wireshark decoded headers versus ASCII on the wire.

#### Medium practical tasks

1. Compare `curl -v` header list size in lines for one request. Write that HPACK is not visible as text in `-v`.
2. Write six sentences: a large cookie on HTTP/1.1 versus HTTP/2 (idea).
3. Read a public one-page HPACK or QPACK overview. Write four STE sentences. No long quotes.

#### Advanced practical tasks

1. Write a one-page awareness note: HPACK, QPACK, and why QUIC needed a new scheme.
2. In Wireshark, open an HTTP/2 sample (tool sample or your own decrypt lab). Write whether the tool shows decoded headers. Use only captures that you may share.

---

## QUIC over UDP

QUIC is a transport protocol that runs over UDP. HTTP/3 uses QUIC. QUIC includes TLS 1.3 for keys. QUIC is not "raw UDP with no reliability." QUIC can retransmit stream data. QUIC can also support unreliable datagrams in later extensions. This topic uses reliable HTTP/3 streams.

Why UDP: user-space evolution, one round trip for crypto and transport in the common 1-RTT case, and no TCP head-of-line across streams. NAT and firewalls already pass UDP for many users. Some networks still block UDP. Then the client falls back to HTTP/2 or HTTP/1.1 on TCP.

A QUIC connection uses connection IDs. A client IP or port change (NAT rebinding, a phone on a new network) can keep the connection when the IDs match the design. TCP ties the connection to the 4-tuple more tightly.

Do not write a custom "QUIC" with raw UDP in production. Use a library. Do not assume that UDP port 443 is always open. Measure.

`curl` can use HTTP/3 when the build supports it. The verbose output then says `HTTP/3`. A capture shows UDP, not TCP, for that session.

QUIC still has congestion control. A loss still slows the sender. QUIC does not make a thin link fat.

### Questions

#### Theoretical questions

1. Which transport port protocol does QUIC use?
2. Does QUIC include TLS?
3. Why can a UDP block force a fallback to HTTP/2?
4. What extra identity can a QUIC connection use besides the 4-tuple?
5. Does QUIC remove congestion control?

#### Easy practical tasks

1. Write five sentences that correct "QUIC is unreliable UDP."
2. Make a table: TCP plus TLS, QUIC. Add rows for header layer and HOL across streams.
3. Draw: IP, UDP, QUIC, HTTP/3.
4. Write four sentences: NAT rebinding and connection IDs (idea).

#### Medium practical tasks

1. Run `curl -V` and write whether the build lists HTTP3. Then try `curl -v --http3` if the flag exists. Write the result.
2. Write six sentences: a hotel network that allows TCP 443 and drops UDP.
3. Capture a browser session to a site that uses HTTP/3 if you can. Write TCP or UDP. Skip if you cannot tell.

#### Advanced practical tasks

1. Write a one-page comparison: TCP 4-tuple versus QUIC connection ID for a mobile client.
2. Design a fallback: try HTTP/3, then HTTP/2. Write timeouts as ideas, not as a vendor config dump.

---

## Why HTTP/3 exists (loss + HOL)

TCP delivers an ordered byte stream. A lost segment delays all later bytes for the application. HTTP/2 puts many streams on that one stream. One loss delays all HTTP/2 streams. That is TCP head-of-line blocking. Topic 9 stated this fact.

HTTP/3 exists to keep HTTP semantics (methods, status, headers, body) and to avoid that TCP HOL cost. QUIC recovers each stream with its own packet numbers. A lost packet for stream 1 does not hold stream 2 in the transport receive queue.

HTTP/3 also aims for a faster handshake: TLS and transport combine. The first request can need fewer round trips than TCP plus TLS 1.2 on a cold connection. Real gains depend on the path and on resumption.

HTTP/3 does not remove loss. A loss still delays that stream and still signals congestion. HTTP/3 does not fix a bad application. A single huge response on one stream still takes time.

Do not enable HTTP/3 only because the version number is larger. Measure the page or the API. Do not blame HTTP/3 when DNS or TLS verify failed first.

Some CDNs and browsers already use HTTP/3. Your `curl` build might not. That mismatch is normal.

### Questions

#### Theoretical questions

1. What happens to HTTP/2 streams when TCP loses one segment?
2. How does HTTP/3 change that loss story at a high level?
3. Does HTTP/3 remove packet loss?
4. What handshake benefit does HTTP/3 aim for?
5. Why is a larger version number not a sufficient reason to switch?

#### Easy practical tasks

1. Write five sentences that explain loss plus HOL as the reason for HTTP/3.
2. Draw HTTP/2 on TCP with one hole in the byte stream and three HTTP streams waiting.
3. Make a table: what HTTP/3 changes, what HTTP/3 does not change. Add two rows each.
4. Write four sentences: HTTP semantics stay, transport changes.

#### Medium practical tasks

1. Write six sentences: a lossy Wi-Fi hop and a page with many small API calls on h2 versus h3 (idea).
2. Find a public HTTP/3 or QUIC explainer. Write four STE sentences. No long quotes.
3. Time one page load in a browser with HTTP/3 if the panel shows h3, and compare a `curl --http1.1` fetch. Write that the test is not fair and why.

#### Advanced practical tasks

1. Write a one-page note that links topic 9 HOL, HTTP/2 multiplexing, and HTTP/3. Use your own words.
2. Design a lab that you own: a local HTTP/3 server if a tool lets you, or write why you cannot run one and what you would measure.

---

## ALPN

ALPN is Application-Layer Protocol Negotiation. The client and the server use ALPN during the TLS (or QUIC) handshake. They choose an application protocol. Common IDs: `http/1.1`, `h2`, `h3`.

Without ALPN, a server on port 443 would not know whether the client wants HTTP/1.1 or HTTP/2 after TLS. The old upgrade dance on cleartext HTTP is not the HTTPS default.

The client offers a list. The server picks one ID that it supports. If the lists do not meet, the handshake or the next step fails. `curl -v` can print the agreed ALPN protocol.

ALPN is not DNS. ALPN is not a URL scheme. `https` is still the scheme for HTTP/2 and HTTP/3 on the web. The version is a handshake result.

Do not invent an ALPN ID for a private protocol without a registry plan if you need interop. Do not assume that a successful TLS handshake means HTTP/2. Read the agreed protocol.

QUIC uses ALPN for `h3`. TCP plus TLS uses ALPN for `h2` or `http/1.1`.

### Questions

#### Theoretical questions

1. What does ALPN select?
2. When does ALPN run?
3. Name two ALPN IDs for HTTP.
4. Does the `https` scheme tell you HTTP/2 versus HTTP/3?
5. Who offers the list and who picks?

#### Easy practical tasks

1. Write five sentences that explain ALPN as "agree the application protocol inside the handshake."
2. Make a table: ID, meaning. Add `http/1.1`, `h2`, `h3`.
3. Run `curl -v https://example.com`. Write any ALPN line that you see.
4. Draw: ClientHello with ALPN list, server choice.

#### Medium practical tasks

1. Compare `curl --http1.1 -v` and default `-v` ALPN offers if visible. Write the difference.
2. Write six sentences: a server that only speaks HTTP/1.1 and a client that offers only `h2`.
3. Find ALPN in a Wireshark TLS handshake if the name is visible. Write the strings. Use a fetch that you may capture.

#### Advanced practical tasks

1. Write a one-page map: URL scheme, DNS, TCP or UDP, ALPN, HTTP version.
2. In Go or Python, open a TLS connection with ALPN offers and print the negotiated protocol to a public host or to a local server that you own.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a browser fetch of `https://example.com` that can become HTTP/2 or HTTP/3. Use ALPN and one transport fact.
2. How do multiplexing, HPACK or QPACK, and HOL fit in one comparison of HTTP/2 and HTTP/3?
3. Why does QUIC on UDP still need TLS ideas from topic 13?
4. When is HTTP/1.1 still the right teaching or debug choice?
5. A teammate says "HTTP/3 is faster because UDP has no congestion control." Which facts do you use to correct that sentence?

#### Easy practical tasks

1. Write a one-page cheat sheet: multiplexing, HPACK, QPACK, QUIC, HOL, ALPN.
2. Run `curl -v` on one HTTPS URL. Save the output. Mark version and ALPN if present.
3. Draw one picture: two stacks, TCP+TLS+h2 and UDP+QUIC+h3.
4. List which tools on your machine can show the HTTP version (`curl`, browser panel).

#### Medium practical tasks

1. From your machine, identify: `curl` HTTP version to `example.com`, browser version to the same name if you can, and whether UDP appears in a short capture.
2. Write a lab plan: force HTTP/1.1, allow HTTP/2, try HTTP/3. State what you will record.
3. Explain in ten steps how a beginner reads `curl -v` to name the HTTP version without reading binary frames.

#### Advanced practical tasks

1. Build a glossary of 12 terms from this topic. Each entry: term, one sentence, one command or UI that shows it.
2. Write a short design: an API that must work on networks that block UDP. State the fallback chain.
