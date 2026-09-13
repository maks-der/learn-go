# 11. Application Protocols: HTTP

## Description

HTTP is an application protocol for request and response messages. This topic shows URLs, methods, status codes, headers, bodies, HTTP/1.1 keep-alive, the Host header, cookies, a TLS preview for HTTPS, `curl -v`, and proxies.

Use one term for each concept. A request is the client message. A response is the server message. Complete this topic after TCP and sockets. Later topics cover HTTP/2, HTTP/3, and TLS in more depth.

Practice with `curl -v` and a browser. Capture HTTP on port 80 when the site still allows it. HTTPS hides the HTTP bytes. You still see TCP and TLS.

---

## URL, method, status, headers, body

A URL names a resource. A common form is `scheme://host:port/path?query#fragment`. The scheme is `http` or `https`. The host is a name or an IP address. The default port is 80 for HTTP and 443 for HTTPS. The path is the resource on the server. The query is a string of parameters. The fragment stays in the client. The client does not send the fragment to the server.

A method is the verb of the request. Common methods: `GET` (read), `POST` (submit a body), `PUT` (store), `DELETE` (remove), `HEAD` (headers only), `OPTIONS` (allowed methods). Use `GET` for safe reads. Use `POST` when the request has a side effect that a cache must not repeat as a GET.

A status code is a three-digit number in the response. 1xx is informational. 2xx is success (`200` OK). 3xx is redirect (`301`, `302`, `304`). 4xx is a client error (`400`, `403`, `404`). 5xx is a server error (`500`, `502`, `503`). Read the class first, then the exact code.

Headers are name-value lines. Request headers include `Host`, `User-Agent`, `Accept`, and `Content-Type` when a body exists. Response headers include `Content-Type`, `Content-Length` or `Transfer-Encoding`, `Location` on redirects, and `Set-Cookie`.

The body is optional. `GET` often has no body. `POST` often has a body. The body is bytes. `Content-Type` names the format (`text/html`, `application/json`).

HTTP/1.1 text starts with a request line or a status line, then headers, then a blank line, then the body. Do not forget the blank line.

### Questions

#### Theoretical questions

1. Which part of a URL does the client not send to the server?
2. What is the default port for HTTPS?
3. What does status class 4xx mean?
4. Where do headers sit relative to the body in HTTP/1.1 text?
5. Why is `GET` a poor choice for a bank transfer?

#### Easy practical tasks

1. Split `https://example.com:443/path?x=1#top` into scheme, host, port, path, query, fragment.
2. Write five sentences that explain method, status, header, and body.
3. Make a table: code, class, meaning. Add `200`, `301`, `404`, `500`.
4. Write a tiny HTTP/1.1 request in a text file: `GET / HTTP/1.1`, `Host`, blank line.

#### Medium practical tasks

1. Use `curl -v` on `https://example.com`. Write the method, the status, and two response headers.
2. Send a `POST` with a small body (`curl -d`). Write `Content-Type` and the status.
3. Write six sentences: `Content-Length` versus a body that you did not measure.

#### Advanced practical tasks

1. Write a Go or Python HTTP client that prints status and `Content-Type` for a URL. Do not print a huge body.
2. Write a one-page map of ten status codes that you see in daily work. One sentence each.

---

## HTTP/1.1 keep-alive

HTTP/1.0 often closed the TCP connection after one response. HTTP/1.1 uses persistent connections by default. The client and the server can send more requests and responses on the same TCP connection. That is keep-alive.

Keep-alive avoids a new handshake and a new slow start for every image or script. The client still must parse each response. HTTP/1.1 on one connection is sequential: the client usually waits for a response before the next request (unless it uses pipelining, which is rare and fragile). Many browsers open several TCP connections to one host.

The `Connection: close` header asks to close after this response. A server can also close. The client must handle a close and must retry or open a new connection by policy.

Idle timeouts close a quiet keep-alive connection. A proxy can have a shorter idle time than the client. The next request then needs a new TCP session.

Do not assume that one `curl` without `-v` shows keep-alive. Look at the headers and at the packet capture. Do not pipeline requests unless you know the server supports it.

HTTP/2 multiplexes many streams on one connection. That is a later topic. HTTP/1.1 keep-alive is still one request at a time per connection in the usual case.

### Questions

#### Theoretical questions

1. What does keep-alive reuse?
2. Why does keep-alive help a page with many files?
3. What does `Connection: close` mean?
4. Is HTTP/1.1 pipelining the common case?
5. How is HTTP/2 different at a high level?

#### Easy practical tasks

1. Write five sentences that compare one connection per file and keep-alive.
2. In `curl -v` output, find `Connection` headers. Write them.
3. Draw two HTTP requests on one TCP connection. Label the handshake once.
4. Make a table: HTTP/1.0 old habit, HTTP/1.1 default. Add one row for TCP close.

#### Medium practical tasks

1. Capture two `curl` requests to the same HTTP URL in one connection if you can (`curl` with two URLs). Write whether a second handshake appeared.
2. Write six sentences: idle timeout at a proxy and a client error that looks like a random fail.
3. Compare browser connection count to one host (dev tools) with keep-alive. Write what you see.

#### Advanced practical tasks

1. Write a small HTTP/1.1 client that sends two `GET`s on one TCP socket. Parse two responses. Localhost only.
2. Write a one-page note: keep-alive versus HTTP/2 multiplex. Use HOL from the TCP topic.

---

## Host header and virtual hosts

The Host header is required in HTTP/1.1. It carries the host name (and the port if the port is not the default). Example: `Host: example.com`. Example: `Host: localhost:8080`.

One IP address can serve many sites. The server reads the Host header and selects the site files or the virtual host config. That is a virtual host. Without Host, the server cannot tell `a.example` from `b.example` on the same socket.

TLS has Server Name Indication (SNI) in the handshake. SNI names the host before HTTP starts. The server can pick a certificate. HTTPS virtual hosts need SNI and then still use Host in HTTP. A later TLS topic goes deeper.

A wrong Host can return the default site or `400`. Some servers reject a missing Host. `curl` sets Host from the URL. A raw TCP client that you write must send Host.

Do not put the scheme in Host. Do not use the IP as Host when you need a name-based virtual host, unless that name is the IP on purpose.

DNS maps the name to an address. Host tells the application which name the user asked for. Both are required for name-based hosting.

### Questions

#### Theoretical questions

1. Why is Host required in HTTP/1.1?
2. What is a virtual host?
3. What problem does SNI solve for HTTPS?
4. Who sets Host when you use `curl` with a URL?
5. Why is DNS not enough without Host?

#### Easy practical tasks

1. Write five sentences that explain two names on one IP.
2. From `curl -v`, copy the `Host` request header.
3. Draw: DNS to IP, TCP to port 443, SNI, then HTTP Host.
4. Make a table: Host, SNI, DNS A/AAAA. Add one job each.

#### Medium practical tasks

1. Send a request with a wrong Host to a local server if you can. Write the status.
2. Write six sentences: default virtual host when Host does not match.
3. Open a site that shares an IP with another site (if you know one). Write the Host that `curl` sends.

#### Advanced practical tasks

1. Configure two name-based sites on a local HTTP server (Caddy, nginx, or a tiny Go server). Fetch both names. Document Host.
2. Write a one-page note: why `/etc/hosts` plus Host lets you test a name on localhost.

---

## Cookies

A cookie is a small piece of state that the server asks the client to store. The `Set-Cookie` response header sends a cookie. The client later sends `Cookie` request headers on matching requests.

A cookie has a name, a value, and attributes: `Domain`, `Path`, `Expires` or `Max-Age`, `Secure`, `HttpOnly`, `SameSite`. `Secure` means send only on HTTPS. `HttpOnly` means a page script must not read the cookie. `SameSite` limits cross-site sends.

Cookies are not a secret vault. Anyone who can read the cookie header can impersonate the session if the cookie is a session id. Use HTTPS. Use `Secure` and a strict `SameSite` when you can.

The client stores cookies per browser profile. `curl` does not keep cookies unless you use `-c` and `-b` or `--cookie-jar`.

Do not put passwords in cookies. Do not use a huge cookie as a database. Do not ignore `SameSite` when you debug a missing cookie on a cross-site request.

A cookie is not the only session method. Authorization headers and server-side sessions also exist. This section is the HTTP cookie mechanism.

### Questions

#### Theoretical questions

1. Which header sets a cookie?
2. Which header sends a cookie back?
3. What does the `Secure` attribute mean?
4. What does `HttpOnly` prevent?
5. Why is a session cookie sensitive?

#### Easy practical tasks

1. Write five sentences that describe `Set-Cookie` and `Cookie`.
2. In a browser, open storage tools for a site that you use. Write one cookie name (not a secret value).
3. Make a table: attribute, meaning. Add `Secure`, `HttpOnly`, `Path`.
4. Write a `curl` command that saves cookies to a file (`man curl` or `curl --help`).

#### Medium practical tasks

1. Use `curl -c` and `-b` against a local server that sets a cookie, or against a public demo that you may test. Write the file contents without secrets from a real account.
2. Write six sentences: `SameSite` and a third-party request.
3. Compare a session cookie (no expires) and a persistent cookie. Write a short table.

#### Advanced practical tasks

1. Write a tiny server that sets `Set-Cookie` and checks `Cookie` on the next request. Localhost only.
2. Write a one-page hardening list for session cookies. No exploit steps.

---

## HTTPS preview (TLS)

HTTPS is HTTP over TLS. TLS provides confidentiality, integrity, and authentication of the server (and optionally the client). The browser checks a certificate. The certificate binds a public key to a name. A later topic covers PKI.

The client starts a TLS handshake after TCP. Then the client sends HTTP on the protected channel. A packet capture shows TCP and TLS records. It does not show the HTTP URL path unless you decrypt with a key that you own or you use a local proxy that you control.

The URL scheme `https` means port 443 by default and a TLS requirement. A mixed page that loads HTTP images on an HTTPS site is a security fault in the browser.

Common failures: expired certificate, name mismatch, untrusted issuer, and an interceptor on public Wi-Fi. The client must stop. Do not click through a warning on a real site.

`curl -v` prints the certificate subject and the TLS version. Use that output to learn. `curl -k` skips verify. Do not use `-k` in production.

Do not call HTTPS "encrypted HTTP" and then forget authentication. Encryption without a name check is not enough.

### Questions

#### Theoretical questions

1. What does HTTPS add under HTTP?
2. When does the TLS handshake run relative to TCP?
3. Why does a normal capture hide the HTTP path on HTTPS?
4. Name three certificate failures.
5. Why is `curl -k` dangerous in production?

#### Easy practical tasks

1. Write five sentences that explain HTTPS as HTTP plus TLS.
2. Run `curl -v https://example.com`. Write the TLS version and the certificate subject if shown.
3. Make a table: HTTP port 80, HTTPS port 443. Add one row for capture visibility.
4. Draw: TCP, TLS, HTTP as three layers.

#### Medium practical tasks

1. Visit a site with a browser padlock. Open the certificate view. Write the name and the expiry date.
2. Write six sentences: name mismatch (certificate for A, URL for B).
3. Compare `curl -v http://example.com` and `https://example.com` if HTTP still redirects. Write the status and the `Location` if any.

#### Advanced practical tasks

1. Create a local self-signed or mkcert certificate for a local server. Fetch with `curl` and fix the trust error. Document the commands. Use only your machine.
2. Write a one-page preview of what the TLS topic will add (handshake, PKI, common failures). Do not copy a full TLS spec.

---

## `curl -v`

`curl` is a command-line HTTP client. `-v` (verbose) prints the request headers, the response headers, and connection details. Use `-v` as a learning tool. Use `-sS` when you want a quiet body plus errors.

Useful options:

- `-I` or `--head` for headers only (`HEAD`)
- `-X` to set the method
- `-H` to add a header
- `-d` to send a body
- `-L` to follow redirects
- `-o` to write the body to a file
- `-4` and `-6` to force a family
- `--http1.1` to avoid HTTP/2 in the client

Read the verbose lines. `*` is info about the connection. `>` is the request. `<` is the response. That split trains your eye.

`curl` can speak HTTP/2. The verbose output shows the version. A capture of HTTP/2 is binary. This topic still uses HTTP/1.1 text as the teaching form.

Do not paste verbose output that contains `Authorization` or cookies from a real account. Do not use `curl` to attack a host. Use sites that you own or public examples.

`curl -v` is also useful for TLS errors. The message often names the check that failed.

### Questions

#### Theoretical questions

1. What does `curl -v` print that a normal `curl` hides?
2. What do `>` and `<` mean in verbose mode?
3. What does `-L` do?
4. What does `-I` send as the method?
5. Why must you redact verbose logs?

#### Easy practical tasks

1. Run `curl -v https://example.com`. Save the output. Mark one `>` line and one `<` line.
2. Run `curl -I https://example.com`. Write the status and `content-type`.
3. Write five sentences that explain `-H` and `-d`.
4. Make a table: option, purpose. Add `-v`, `-L`, `-o`, `-4`.

#### Medium practical tasks

1. Follow a redirect with and without `-L` on a URL that redirects (many `http://` URLs). Write both statuses.
2. Send a custom `User-Agent` with `-H`. Confirm it in `-v`.
3. Write six sentences: `--http1.1` versus the default when `curl` negotiates HTTP/2.

#### Advanced practical tasks

1. Write a small script that runs `curl -v` and extracts the status line. Test on `example.com`.
2. Write a one-page `curl` cheat sheet for this learning path. Include TLS and IPv6 flags.

---

## Proxies and reverse proxies

A forward proxy sits near the client. The client sends the request to the proxy. The proxy opens a connection to the origin server. Enterprises use forward proxies for policy and cache. The client must be configured (`HTTP_PROXY`, browser settings). For HTTPS, the client uses `CONNECT` to open a tunnel. The proxy sees the name and port. The TLS payload stays opaque if the proxy does not intercept.

A reverse proxy sits near the servers. The client thinks it talks to the origin. The reverse proxy forwards to an internal app. TLS can terminate at the reverse proxy. Load balancers and CDNs often act as reverse proxies. `Host` and `X-Forwarded-For` (or `Forwarded`) tell the app the original client and name. Trust those headers only from your proxy.

A proxy can be HTTP-aware (layer 7) or only a TCP pass-through (layer 4). This topic uses the HTTP-aware case.

Do not expose an open forward proxy to the Internet. Do not trust `X-Forwarded-For` from the public Internet.

A reverse proxy can add keep-alive to the client and a new connection pool to the app. Timeouts at the proxy are a common source of `502` and `504`.

SOCKS is another proxy family. Browsers and some tools use SOCKS. The idea is still "client talks to a middle hop."

### Questions

#### Theoretical questions

1. Where does a forward proxy sit?
2. Where does a reverse proxy sit?
3. What does `CONNECT` do for HTTPS?
4. Why must an app trust `X-Forwarded-For` only from its proxy?
5. What status class often means a reverse proxy could not get a good answer?

#### Easy practical tasks

1. Write five sentences that compare forward and reverse proxies.
2. Draw a client, a reverse proxy, and two app processes.
3. Make a table: forward proxy, reverse proxy. Add one row for who configures the client.
4. List two products that act as a reverse proxy (nginx, Caddy, a cloud load balancer).

#### Medium practical tasks

1. Read a public nginx reverse-proxy example. Write the idea of `proxy_pass` in four STE sentences. No long config copy.
2. Write six sentences: `502` versus `504` at a reverse proxy (idea).
3. Explain `CONNECT` in a capture if you have a lab proxy. Otherwise write the expected TCP peers.

#### Advanced practical tasks

1. Run a local reverse proxy in front of a local app on a machine that you own. Fetch the proxy URL. Document Host and the backend port.
2. Write a one-page note: which headers a reverse proxy should set and which headers it should strip.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe an HTTPS request from URL parse to TLS to HTTP headers to a possible cookie on the next request.
2. How do keep-alive, Host, and a reverse proxy work together on one public name?
3. Why does `curl -v` help when a browser only shows a blank page?
4. Which facts separate a forward proxy, a reverse proxy, and NAT?
5. A teammate sends a secret in a GET query and in a cookie on HTTP. Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: URL parts, method, status classes, keep-alive, Host, cookies, HTTPS, `curl -v`, proxies.
2. Run `curl -v` on one HTTP and one HTTPS URL. Save both. Redact secrets.
3. Draw a virtual-host box: one IP, two Host values, two document roots.
4. List the default ports and schemes for HTTP and HTTPS.

#### Medium practical tasks

1. Write a lab: local HTTP server, `curl -v`, then add TLS with a local cert. Record each command.
2. Use browser tools and `curl` on the same URL. Compare method, status, and two headers.
3. Write a troubleshooting flow: `404`, `502`, certificate name error, and a missing cookie.

#### Advanced practical tasks

1. Write a small HTTP/1.1 server in Go or Python that handles `GET /` and `Host`. Test with `curl -v`.
2. Build a glossary of 14 HTTP terms. Each entry: term, one sentence, one `curl` flag or header name.
