# 13. TLS and Security Basics

## Description

TLS protects a byte stream between two programs. This topic shows confidentiality, integrity, authenticity, certificates, PKI, a high-level handshake, HTTPS as HTTP over TLS, common certificate failures, and SSH as a separate secure channel.

Use one term for each concept. Confidentiality means an observer cannot read the payload. Integrity means a change is detectable. Authenticity means you talk to the intended peer. Complete this topic after DNS and the HTTP preview of HTTPS.

Practice with `curl -v` and a browser certificate view. Capture TLS records. You see versions and handshake messages. You do not see HTTP paths on a normal HTTPS session.

---

## Confidentiality, integrity, authenticity

Confidentiality hides the payload from a path observer. TLS encrypts application bytes after the handshake. Addresses and ports stay visible. The Server Name Indication (SNI) name is often still visible in older TLS. TLS 1.3 can encrypt more of the handshake. This section needs the idea: encryption is not invisibility of the connection.

Integrity lets the receiver detect a change to the ciphertext or the records. A middlebox that flips bits should cause a failure, not a silent wrong page. TLS uses MAC or AEAD for this job. You do not need the math. You need the property.

Authenticity binds the session to a peer identity. For HTTPS, the identity is a DNS name on a certificate. Encryption without a name check is not enough. An attacker on the path can terminate TLS with a different certificate if the client does not verify.

The three words work together. Confidentiality without authenticity can hide data from a third party and still give the data to the wrong peer. Authenticity without integrity is not a complete channel. Integrity without confidentiality can still leak secrets.

Do not say "secure" without naming which of the three you mean. Do not treat a VPN as a substitute for TLS to a public API unless you designed that trust model.

SSH uses the same three ideas with a different protocol and a different identity model (host keys and user keys).

### Questions

#### Theoretical questions

1. What does confidentiality hide?
2. What does integrity detect?
3. What does authenticity bind?
4. Why is encryption without a name check not enough?
5. Which parts of a TLS session stay visible on the path in the common case?

#### Easy practical tasks

1. Write five sentences that define the three words with one example each.
2. Make a table: property, what a path attacker must not do. Add three rows.
3. Draw a client, a path observer, and a server. Mark what TLS hides.
4. Write four sentences: a padlock icon versus the three properties.

#### Medium practical tasks

1. Use `curl -v` on `https://example.com`. Write the TLS version if the tool shows it.
2. Write six sentences: HTTP on port 80 versus HTTPS for a password form.
3. Compare SNI visibility with HTTP Host visibility on HTTPS. Write a short paragraph.

#### Advanced practical tasks

1. Write a one-page note: three properties, what TLS does not hide (addresses, timing, SNI idea).
2. Capture a TLS session. List fields you can read and fields you cannot read. Use a site that you may fetch.

---

## Certificates and PKI

A certificate binds a public key to a name (and sometimes to other names). A certificate authority (CA) signs that bind. The client has a store of trusted root CA certificates. The operating system and the browser ship those roots.

PKI is public key infrastructure. PKI is the set of CAs, certificates, revocation data, and policies that make the bind usable. You do not run a public CA for this topic. You must know the trust chain: server certificate, optional intermediates, and a root that the client already trusts.

The client checks:

- the name in the URL matches a name on the certificate
- the current time is inside the validity interval
- the chain leads to a trusted root
- the signature on each link is valid
- revocation status when the client implements that check

A self-signed certificate is signed by its own key. A browser does not trust it unless you add an exception or you install a local CA such as mkcert. Use self-signed or mkcert only on machines that you own.

Do not send a private key in chat or in a repository. Do not disable verify to "make it work" on a real user path. Do not treat a lock icon on a phishing name as safety. The lock only says TLS to that name succeeded.

A client certificate is optional. Some APIs use mutual TLS. This topic only needs: the common web case authenticates the server.

### Questions

#### Theoretical questions

1. What two things does a server certificate bind?
2. What is a root CA in the client store?
3. Name four checks that a client can apply to a certificate.
4. What is a self-signed certificate?
5. Why is a lock icon on a wrong name not safety?

#### Easy practical tasks

1. Write five sentences that explain a trust chain: server, intermediate, root.
2. Open a browser certificate view for `https://example.com`. Write the subject name and the issuer.
3. Make a table: check, what fails if it is wrong. Add name, dates, and untrusted root.
4. Write four sentences: mkcert or a local CA versus a public CA.

#### Medium practical tasks

1. Run `curl -v https://example.com`. Write the subject and the expiry if shown.
2. Write six sentences: why an intermediate CA exists (the idea).
3. Find where your OS or browser lists trusted roots (settings name only). Write the path or the settings label.

#### Advanced practical tasks

1. Create a local certificate for a local server on a machine that you own. Fetch with `curl` and then trust the cert or use mkcert. Document the commands. Do not use this cert on the public Internet.
2. Write a one-page PKI map for a teammate: leaf, intermediate, root, private key location.

---

## Handshake (high-level)

The TLS handshake runs after the TCP handshake for TLS on TCP. The two sides agree on a TLS version and on algorithms. The server sends a certificate (usual web case). The client verifies the certificate. The two sides derive keys. Then they send encrypted application data.

TLS 1.3 is the current common version on the public web. The wire messages differ from TLS 1.2. The teaching story is the same: authenticate the server, then protect the bytes.

ALPN can run inside the handshake. ALPN selects the application protocol (HTTP/1.1, HTTP/2, or others). The next topic covers ALPN with HTTP/2 and HTTP/3.

A handshake can fail before any HTTP request. `curl -v` then prints a TLS error, not an HTTP status. A capture shows an alert or a close. You will debug that pattern again in the observability topic.

Do not invent a custom "encryption" of HTTP headers on port 80 and call it TLS. Do not skip verify. Do not paste a full handshake dump that contains secrets from a lab that used session tickets you must keep private.

Resumption can skip a full handshake on a later connection. The idea is speed. The security still depends on the first trust decision.

QUIC includes TLS 1.3 in its own UDP handshake. That is the HTTP/3 topic. This section is TLS on TCP unless you read ahead.

### Questions

#### Theoretical questions

1. When does the TLS handshake run relative to TCP in HTTPS?
2. What does the client verify during the handshake?
3. What does the handshake produce for later records?
4. How do you see a handshake failure in `curl -v` versus an HTTP 404?
5. What is resumption in one sentence?

#### Easy practical tasks

1. Write five sentences that list handshake steps at a high level.
2. Draw: TCP handshake, TLS handshake, HTTP request.
3. Make a table: TLS 1.2 idea, TLS 1.3 idea. Add one row: "still verify the name".
4. Run `curl -v https://example.com`. Mark the lines that look like TLS, not HTTP.

#### Medium practical tasks

1. Capture a ClientHello if your tool names it. Write the destination port and whether you see a server name field.
2. Write six sentences: a failed handshake versus a failed HTTP login (401).
3. Compare `curl --tlsv1.2` and the default on a public site if your `curl` supports the flag. Write what the tool reports. Do not force old TLS on a host that you do not own as a test of weaknesses.

#### Advanced practical tasks

1. Write a one-page handshake map: messages you expect in a capture (ClientHello, ServerHello, certificate, application data). No byte-level spec.
2. On a local server that you own, capture a full handshake. Label three packets. Use only your machine.

---

## HTTPS = HTTP over TLS

HTTPS is HTTP on a TLS session. The default port is 443. The URL scheme is `https`. The client parses the URL, resolves the name, connects with TCP, completes TLS, then sends the HTTP request on the protected stream.

The HTTP methods, status codes, and headers are the same ideas as topic 11. TLS does not replace `Host`. The client still sends `Host`. The path observer does not see that header when TLS works.

A reverse proxy often terminates TLS. The proxy speaks HTTPS to the client and HTTP or HTTPS to the app. The app must know whether the original request was HTTPS (`X-Forwarded-Proto` or an equivalent). Trust that header only from your proxy.

HSTS tells a browser to use HTTPS for a name for a time. Mixed content is HTTP resources on an HTTPS page. Browsers block or warn. Fix the resource URLs.

Do not put secrets in an HTTP URL and then "add TLS later" without a review. Old links and logs can still leak. Do not call a site HTTPS if TLS verify is off.

`curl` uses HTTPS when the URL says `https`. `curl -k` is still HTTP over TLS with verify off. That is not the same as safe HTTPS.

### Questions

#### Theoretical questions

1. What three steps sit between URL parse and the HTTP request on HTTPS?
2. What is the default HTTPS port?
3. Does TLS remove the need for the Host header?
4. What does TLS termination at a reverse proxy mean?
5. How is `curl -k` different from a verified HTTPS fetch?

#### Easy practical tasks

1. Write five sentences that expand "HTTPS = HTTP over TLS".
2. Run `curl -v https://example.com`. Write the scheme, the port, and the HTTP status.
3. Make a table: HTTP, HTTPS. Add rows for port, TLS, and capture of the path.
4. Draw layers: IP, TCP, TLS, HTTP.

#### Medium practical tasks

1. Fetch the same name with `http://` and `https://` if the host answers both. Write the status and any `Location`.
2. Write six sentences: TLS terminate at a proxy, then HTTP on a private LAN to the app.
3. Find HSTS in a browser for a site that uses it, or write that you did not find the UI. Record what you saw.

#### Advanced practical tasks

1. Run a local HTTP server and a local TLS frontend (Caddy, nginx, or `caddy`/`openssl s_server` style) on a machine that you own. Fetch HTTPS locally. Document ports.
2. Write a one-page note: which HTTP fields a path observer sees on HTTP versus HTTPS.

---

## Common failures: expired cert, wrong host, MITM on public Wi-Fi

Expired certificate: the validity end time is in the past (or the start time is in the future). The client must fail. Cause: the operator did not renew. Your clock can also be wrong. A wrong clock looks like an expiry fault.

Wrong host: the certificate names do not include the URL host. Cause: wrong vhost, a default cert, or a raw IP in the URL. The client must fail.

Untrusted issuer: the chain does not reach a root that this client trusts. Cause: corporate intercept CA not installed, a self-signed server, or a missing intermediate.

MITM on public Wi-Fi: a captive portal or an attacker presents a different certificate. A correct client shows a warning. Do not ignore the warning. A portal can also intercept HTTP and redirect you. Use HTTPS and a name that you typed or bookmarked.

Other common failures: revoked certificate, weak old TLS that the client refuses, and a name that does not resolve. Read the exact error. `curl` and browsers use different words for the same check.

Do not click through a certificate warning on a bank, mail, or work site. Do not use public Wi-Fi for secrets without TLS that verifies. Do not set the system clock back to "fix" expiry.

A company TLS interceptor is a later hardening topic. This section needs: the client sees an unexpected issuer.

### Questions

#### Theoretical questions

1. What does an expired certificate mean?
2. What does a wrong-host failure mean?
3. Why can a wrong system clock look like expiry?
4. What should a client do on a certificate warning on public Wi-Fi?
5. Name one failure that is not expiry and not a name mismatch.

#### Easy practical tasks

1. Write five sentences that list expired, wrong host, and untrusted issuer.
2. Make a table: failure, what you check first. Add three rows.
3. Write four sentences you would tell a user who wants to "just click continue".
4. Run `curl -v` to a public HTTPS site that works. Write that the checks passed.

#### Medium practical tasks

1. Point `curl` at a name with HTTPS and a wrong `Host` via `--resolve` or a local hosts trick on a machine that you own, or explain the expected error in six sentences if you cannot lab it safely.
2. Turn off Wi-Fi time sync only if you may, or write why a wrong clock breaks TLS. Do not leave a work laptop on a wrong date.
3. Write six sentences: captive portal versus a TLS MITM warning.

#### Advanced practical tasks

1. On a local server that you own, serve a certificate for name A and fetch name B. Record the client error. Then fix the name. Document both steps.
2. Write a one-page triage: expiry, name, issuer, clock. Each row: how you confirm, how you do not "fix" it unsafely.

---

## SSH as a secure channel

SSH is a protocol for a secure channel, usually on TCP port 22. Common uses: a remote shell, file copy (`scp`, `sftp`), and port forwarding. SSH is not HTTP. SSH is not TLS. The goals match: confidentiality, integrity, authenticity.

The server has a host key. The client stores a fingerprint in `known_hosts` after the first success (or you pin the fingerprint). A changed host key is a warning. It can mean a new server or an attacker. Do not ignore a changed-key warning on a host that should be stable.

The user authenticates with a password or, better, with a public key pair. Keep the private key private. Use a passphrase on the key file.

SSH can tunnel other TCP ports. That tunnel is a tool, not a full network design. Topic 21 covers VPNs at a high level.

Do not disable host-key checks to save time. Do not reuse one user key on a public paste. Do not expose SSH to the Internet without a plan (keys, updates, and limited users).

`ssh -v` prints handshake steps. Use it to learn. Use it on hosts that you may access.

### Questions

#### Theoretical questions

1. What does SSH protect at a high level?
2. What is a host key for?
3. Why is a changed host key a warning?
4. Why is a key pair better than a password for many servers?
5. How is SSH different from HTTPS in role?

#### Easy practical tasks

1. Write five sentences that compare SSH and TLS as secure channels.
2. Make a table: SSH, HTTPS. Add rows for default port, typical user, and identity (host key versus certificate name).
3. Write four sentences: `known_hosts` and the first connection.
4. Find whether `ssh` exists on your machine. Write the version (`ssh -V`).

#### Medium practical tasks

1. Read `man ssh` or a vendor SSH page for `-L` or `-D` at a high level. Write four STE sentences. Do not open a tunnel to a network that you must not use.
2. Write six sentences: password auth versus public-key auth.
3. If you have a lab VM that you own, connect with SSH once. Write how the host key prompt looked. Skip if you have no VM.

#### Advanced practical tasks

1. Write a one-page SSH hardening list for a lab VM: keys, no unused accounts, updates. No attack steps.
2. Compare `sftp` and HTTPS file upload in a short note: channel, identity, and typical use.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe an HTTPS session from name lookup to HTTP. Use the three security properties and one certificate check.
2. How do PKI, the handshake, and a reverse proxy that terminates TLS fit on one public name?
3. Why can SSH and TLS both be "secure" and still not replace each other?
4. When is `curl -k` the wrong tool, and when is a local trust store the right fix?
5. A teammate says "the Wi-Fi password encrypts my HTTPS." Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: three properties, certificate chain, handshake order, HTTPS layers, three failures, SSH host key.
2. Run `curl -v https://example.com` and open the same site in a browser certificate view. Write two facts that match.
3. Draw one picture: client, public Wi-Fi, server. Label where TLS must verify.
4. List default ports: HTTPS and SSH.

#### Medium practical tasks

1. Write a lab: local TLS server, a name mismatch, then a fix. Record commands. Use only your machine.
2. Build a troubleshooting flow: handshake error, HTTP 403, SSH changed host key. One box each.
3. Compare `curl -v` TLS lines with a Wireshark TLS summary for one fetch. Write three fields that both show.

#### Advanced practical tasks

1. Write a small HTTPS client in Go or Python that verifies TLS and prints the peer certificate names. Test on `example.com`.
2. Build a glossary of 14 terms from this topic. Each entry: term, one sentence, one tool that shows it.
