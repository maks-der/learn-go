# 10. TLS

## Description

TLS protects a byte stream between two programs. This topic shows confidentiality, integrity, and authenticity. You learn certificates and PKI. You learn the handshake at a high level. You learn common failures. You learn SSH as a separate secure channel.

Complete this topic after DNS and the HTTP preview of HTTPS.

Use one term for each concept. Confidentiality means an observer cannot read the payload. Integrity means a change is detectable. Authenticity means you talk to the intended peer. A certificate binds a public key to a name. PKI is the set of CAs and policies that make that bind usable. Do not say "secure" without naming which of the three properties you mean. Do not treat SSH as TLS.

Practice with `curl -v` and a browser certificate view. Capture TLS records. You see versions and handshake messages. You do not see HTTP paths on a normal HTTPS session.

---

## Confidentiality, integrity, authenticity

Confidentiality hides the payload from a path observer. TLS encrypts application bytes after the handshake. Addresses and ports stay visible. The Server Name Indication (SNI) name is often still visible in older TLS. TLS 1.3 can encrypt more of the handshake. Encryption is not invisibility of the connection.

Integrity lets the receiver detect a change to the records. A middlebox that flips bits must cause a failure, not a silent wrong page. TLS uses MAC or AEAD for this job. You do not need the math. You need the property.

Authenticity binds the session to a peer identity. For HTTPS, the identity is a DNS name on a certificate. Encryption without a name check is not enough. An attacker on the path can terminate TLS with a different certificate if the client does not verify.

The three words work together. Confidentiality without authenticity can hide data from a third party and still give the data to the wrong peer. Authenticity without integrity is not a complete channel. Integrity without confidentiality can still leak secrets.

Do not treat a VPN as a substitute for TLS to a public API unless you designed that trust model. SSH uses the same three ideas with a different protocol and a different identity model.

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

A certificate binds a public key to a name and sometimes to other names. A certificate authority (CA) signs that bind. The client has a store of trusted root CA certificates. The operating system and the browser ship those roots.

PKI is public key infrastructure. PKI is the set of CAs, certificates, revocation data, and policies that make the bind usable. You do not run a public CA for this topic. You must know the trust chain: server certificate, optional intermediates, and a root that the client already trusts.

The client checks:

- the name in the URL matches a name on the certificate
- the current time is inside the validity interval
- the chain leads to a trusted root
- the signature on each link is valid
- revocation status when the client implements that check

A self-signed certificate is signed by its own key. A browser does not trust it unless you add an exception or you install a local CA such as mkcert. Use self-signed or mkcert only on machines that you own.

A client certificate is optional. Some APIs use mutual TLS. This topic only needs: the common web case authenticates the server.

Do not send a private key in chat or in a repository. Do not disable verify to make it work on a real user path. Do not treat a lock icon on a phishing name as safety. The lock only says TLS to that name succeeded.

### Questions

#### Theoretical questions

1. What does a certificate bind?
2. What is a CA in this section?
3. What is a trust chain?
4. Why does a browser reject a typical self-signed certificate?
5. What name must match the certificate?

#### Easy practical tasks

1. Write five sentences that explain certificate, CA, and root store.
2. Open a browser certificate view for `example.com`. Write the subject name and the issuer.
3. Make a table: check, what fails if it is wrong. Add name, time, and root.
4. Draw: leaf certificate, intermediate, root.

#### Medium practical tasks

1. Use `curl -v` and write the certificate subject if printed.
2. Write six sentences: why the private key must stay on the server.
3. Find where your OS stores trusted roots (name of the store). Do not export a pile of certificates.

#### Advanced practical tasks

1. Create a self-signed certificate for a local server that you own. Show that `curl` fails verify and that `-k` or a custom CA changes the result. Use only your lab.
2. Write a one-page PKI note: leaf, intermediate, root, and revocation as names only.

---

## Handshake (high-level)

The TLS handshake agrees on a TLS version, a cipher suite, and keys. It authenticates the server with the certificate in the common web case. Then the two sides encrypt application data.

A teaching sequence for TLS 1.2-style talk:

1. Client Hello: version, random, offered ciphers, SNI
2. Server Hello: chosen version and cipher
3. Server certificate
4. Key exchange messages
5. Change to encrypted records
6. Application data (HTTP)

TLS 1.3 combines steps and encrypts more of the handshake. You still need Client Hello, a server choice, a certificate, and keys. You do not need every extension.

ALPN (Application-Layer Protocol Negotiation) lets the client offer `http/1.1` and `h2`. The server picks. Topic 11 uses ALPN.

The handshake needs extra round trips on a new connection. Session resumption can shorten a later handshake. TCP still needs its own handshake first on classic HTTPS.

Do not confuse the TCP three-way handshake with the TLS handshake. Both happen for HTTPS on TCP. Do not expect to read HTTP in a capture before the handshake finishes.

A failure in the handshake aborts before HTTP. `curl` prints a TLS error. Topic 13 debugs common failures.

### Questions

#### Theoretical questions

1. What does the handshake agree on?
2. What is Client Hello for?
3. What is SNI?
4. What does ALPN choose?
5. How is the TLS handshake different from the TCP handshake?

#### Easy practical tasks

1. Write five sentences that describe a high-level handshake.
2. Draw TCP handshake, then TLS handshake, then HTTP.
3. Make a table: TLS 1.2 idea, TLS 1.3 idea. Add "more encryption of handshake" as one row.
4. Run `curl -v https://example.com`. Write ALPN or HTTP version if shown.

#### Medium practical tasks

1. In a capture, find Client Hello. Write whether SNI appears in the tool.
2. Write six sentences: why a new HTTPS connection pays two handshakes on TCP.
3. Compare two `curl -v` runs to the same host. Write whether the tool mentions reuse or session.

#### Advanced practical tasks

1. Write a one-page handshake map for TLS 1.3 from a public short overview. Use STE. Do not copy long passages.
2. Force `curl --tlsv1.2` and default `curl` if the site allows both. Write both versions.

---

## Common failures

Common TLS failures have names. Learn the symptom and the usual cause.

Name mismatch: the URL host is not on the certificate. A wrong Host or a raw IP against a name-only cert fails.

Expired or not-yet-valid certificate: the clock is wrong or the operator forgot to renew. NTP (topic 12) matters. A laptop with a dead clock fails everywhere.

Untrusted issuer: self-signed cert, missing intermediate, or a corporate inspect CA that the client does not have.

Revoked certificate: the client checked and refused. Not all clients check the same way.

Protocol or cipher mismatch: old client, old server, or a middlebox that blocks TLS 1.3.

Interception: a proxy presents its own certificate. The client must trust that corporate CA. On an unknown network, that is a warning.

Hostname verify disabled: the program used `-k` or `InsecureSkipVerify`. The session can encrypt to the wrong peer.

Do not fix a real user path with verify off. Do not ignore an expired cert because "it works in the browser" if the browser showed a click-through. Do not blame DNS for a name mismatch until you read the certificate names.

`curl -v` and the browser error text are the first tools. A capture shows whether Client Hello leaves and whether an alert returns.

### Questions

#### Theoretical questions

1. What fails when the URL name is not on the certificate?
2. Why can a wrong clock look like an expired certificate?
3. What does an untrusted issuer mean?
4. What does verify off remove?
5. How can a proxy cause a new issuer?

#### Easy practical tasks

1. Write five sentences that list five failure classes from this section.
2. Make a table: error idea, first check. Add clock, name, and issuer.
3. Read a `curl` help line for insecure. Write why you will not use it on production.
4. Draw: correct name versus wrong name on the certificate.

#### Medium practical tasks

1. Point `curl` at a known bad name if you have a lab cert, or read a public expired-cert screenshot description. Write the error class.
2. Write six sentences: missing intermediate versus self-signed.
3. Check your machine clock. Write the source of time (NTP or manual).

#### Advanced practical tasks

1. On a lab that you own, serve a cert with the wrong name and one expired cert (short validity). Record `curl` errors. Then restore a valid local CA.
2. Write a one-page runbook: TLS fail versus TCP fail versus DNS fail. No attack steps.

---

## SSH as a secure channel

SSH is a protocol that provides a secure channel. People use it for a remote shell, for file copy (`scp`, SFTP), and for tunnels. SSH is not TLS and is not HTTP.

SSH authenticates the server with a host key. The client stores the host key fingerprint in `known_hosts`. A changed key is a warning. The client authenticates the user with a password or, better, a public key pair.

SSH provides confidentiality, integrity, and authenticity for the session. The identity model is keys and user accounts, not a public CA store in the usual OpenSSH setup. Some shops use certificates for SSH. That is optional here.

The common port is 22. The wire is not HTTP. A capture shows SSH, not TLS Client Hello.

Use SSH on a path that you may use. Use a key with a passphrase. Do not reuse a password from a website. Do not disable host-key checks to hide a lab warning on a real jump host.

TLS remains the right choice for public HTTPS APIs. SSH remains the right choice for admin shells. A VPN (topic 14) is a third design.

Do not expose SSH to the whole Internet without extra controls (keys only, allow lists). This handbook does not give attack steps.

### Questions

#### Theoretical questions

1. What job does SSH do in this section?
2. How does a client authenticate an SSH server in the usual setup?
3. How does a user authenticate to SSH in the preferred setup?
4. Is SSH the same protocol as TLS?
5. What file stores host keys on a typical client?

#### Easy practical tasks

1. Write five sentences that compare SSH and TLS identity.
2. Make a table: HTTPS, SSH. Add port and identity store.
3. Open `man ssh` or a Windows OpenSSH help page. Write two client options you recognize.
4. Draw: client, SSH channel, remote shell.

#### Medium practical tasks

1. If you have an SSH server that you own, connect and write whether you used a key. Skip if you have none.
2. Write six sentences: `known_hosts` warning when a host key changes.
3. Compare SFTP (topic 11 survey) with HTTPS upload in four sentences.

#### Advanced practical tasks

1. Write a one-page note: when to use TLS, when to use SSH, when to use a VPN. High-level only.
2. Generate a key pair on a lab that you own. Show the public key format. Do not publish the private key.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the three properties, the certificate name check, and SNI work together for HTTPS?
2. Why can a correct TLS handshake still be the wrong peer if you skip verify?
3. What must you check when HTTPS fails but `ping` and `dig` succeed?
4. How is PKI trust different from SSH `known_hosts` trust?
5. Why do TCP handshake failures and TLS alerts need different tools in the story?

#### Easy practical tasks

1. Write a cheat sheet: confidentiality, integrity, authenticity, cert, CA, handshake, SNI, ALPN, common errors, SSH.
2. Run `curl -v https://example.com`. Write version, issuer if shown, and HTTP status.
3. Draw the trust chain and the SSH host-key pin as two pictures.
4. Bookmark a TLS error glossary. Write when you open it.

#### Medium practical tasks

1. Write a lab: self-signed local HTTPS, `curl` fail, install trust or use a lab CA, `curl` success. Use a host that you own.
2. Write six sentences: clock, NTP, and expired certs.
3. Capture Client Hello and an application data record. Write which one you can read.

#### Advanced practical tasks

1. Write a one-page incident sheet: four TLS failures and the first command for each.
2. Compare mutual TLS (client cert) with SSH user keys in a short design note. Do not implement a public CA.
