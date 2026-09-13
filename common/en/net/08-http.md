# 8. HTTP

## Description

HTTP is an application protocol for request and response messages. This topic shows URL, method, status, headers, and body. You learn HTTP/1.1 keep-alive and the Host header. You learn cookies, proxies, and `curl -v`. You get an HTTPS preview.

Complete this topic after TCP and sockets. Later topics cover HTTP/2, HTTP/3, and TLS in more depth.

Use one term for each concept. A request is the client message. A response is the server message. The Host header names the virtual host. Keep-alive reuses a TCP connection. Do not mix a URL fragment with a query. Do not treat HTTPS as a different application verb set. HTTPS is HTTP over TLS.

Practice with `curl -v` and a browser. Capture HTTP on port 80 when the site still allows it. HTTPS hides the HTTP bytes. You still see TCP and TLS.

---

## URL, method, status, headers, body

A URL names a resource. A common form is `scheme://host:port/path?query#fragment`. The scheme is `http` or `https`. The host is a name or an IP address. The default port is 80 for HTTP and 443 for HTTPS. The path is the resource on the server. The query is a string of parameters. The fragment stays in the client. The client does not send the fragment to the server.

A method is the verb of the request. Common methods: `GET` (read), `POST` (submit a body), `PUT` (store), `DELETE` (remove), `HEAD` (headers only), `OPTIONS` (allowed methods). Use `GET` for safe reads. Use `POST` when the request has a side effect that a cache must not repeat as a GET.

A status code is a three-digit number in the response. 1xx is informational. 2xx is success (`200` OK). 3xx is redirect (`301`, `302`, `304`). 4xx is a client error (`400`, `403`, `404`). 5xx is a server error (`500`, `502`, `503`). Read the class first, then the exact code.

Headers are name-value lines. Request headers include `Host`, `User-Agent`, `Accept`, and `Content-Type` when a body exists. Response headers include `Content-Type`, `Content-Length` or `Transfer-Encoding`, `Location` on redirects, and `Set-Cookie`.

The body is optional. `GET` often has no body. `POST` often has a body. The body is bytes. `Content-Type` names the format (`text/html`, `application/json`).

HTTP/1.1 text starts with a request line or a status line, then headers, then a blank line, then the body. Do not forget the blank line.

Do not use `GET` for a bank transfer. Do not send a secret in a URL query if you can put it in a header or a POST body. Logs often store the URL.

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

1. Write an HTTP client in your language that prints status and `Content-Type` for a URL. Do not print a huge body.
2. Write a one-page map of ten status codes that you see in daily work. One sentence each.

---

## HTTP/1.1 keep-alive and the Host header

HTTP/1.0 often closed the TCP connection after one response. HTTP/1.1 uses persistent connections by default. The client and the server can send more requests and responses on the same TCP connection. That is keep-alive.

Keep-alive avoids a new handshake and a new slow start for every image or script. The client still must parse each response. HTTP/1.1 on one connection is sequential: the client usually waits for a response before the next request. Pipelining is rare and fragile. Many browsers open several TCP connections to one host.

The `Connection: close` header asks to close after this response. Idle timeouts close a quiet keep-alive connection. A proxy can have a shorter idle time than the client.

The Host header is required in HTTP/1.1. It names the site. Many servers host several names on one IP address. The Host header selects the site. A request without Host is invalid in HTTP/1.1.

The default port can be omitted in the Host header (`example.com` means port 80 or 443 by scheme). A non-default port appears as `example.com:8080`.

HTTP/2 multiplexes many streams on one connection. That is topic 11. HTTP/1.1 keep-alive is still one request at a time per connection in the usual case.

Do not assume that one `curl` without `-v` shows keep-alive. Look at the headers and at the packet capture. Do not send a request to an IP and forget Host when the server is name-based.

### Questions

#### Theoretical questions

1. What does keep-alive reuse?
2. Why does keep-alive help a page with many files?
3. What does `Connection: close` mean?
4. Why is the Host header required in HTTP/1.1?
5. Is HTTP/1.1 pipelining the common case?

#### Easy practical tasks

1. Write five sentences that explain keep-alive and Host.
2. Make a table: HTTP/1.0 close, HTTP/1.1 keep-alive. Add one cost each.
3. Write a request line plus a Host header for `https://example.com/`.
4. Draw one TCP connection and two HTTP request-response pairs.

#### Medium practical tasks

1. Use `curl -v` twice in one process if you can (`curl` connection reuse) or compare two separate `curl` runs. Write what `-v` says about reuse if it appears.
2. Send `curl` to a name and to the numeric IP with `-H "Host: ..."`. Write both statuses on a host that you may test.
3. Write six sentences: virtual host and one IPv4 address.

#### Advanced practical tasks

1. Write a tiny HTTP/1.1 server that reads Host and returns two different bodies for two names. Bind `127.0.0.1`.
2. Capture keep-alive if a site still speaks clear HTTP, or use your local server. Write the TCP stream index for two requests.

---

## Cookies, proxies, `curl -v`

A cookie is a small piece of state that the server sets with `Set-Cookie`. The client stores it and sends it back on later requests with the `Cookie` header. Cookies implement sessions, preferences, and tracking. Flags matter: `HttpOnly`, `Secure`, `SameSite`, `Path`, and `Domain`.

Do not treat a cookie as a password that you print in a ticket. Do not store secrets in cookies without `Secure` and a short life on a real site. A capture of HTTP shows cookie values in cleartext.

A proxy sits between the client and the origin. A forward proxy is the client choice (office filter, `HTTP_PROXY`). A reverse proxy sits in front of servers (load balancer, ingress). The client of a reverse proxy thinks it talks to the origin.

Proxies can add `X-Forwarded-For` or `Forwarded`. Those headers are claims. The last hop that you trust can set them. Topic 12 covers load balancers.

`curl -v` prints the request, the response headers, TLS facts on HTTPS, and errors. Use it before you write a client. Useful flags: `-I` for HEAD, `-d` for a body, `-H` for a header, `-x` for a proxy, `-o` for the body file, `-w` for timings.

Do not commit a `curl` line that contains a bearer token. Do not trust a proxy that you do not control with HTTPS unless you also understand TLS interception (topic 10).

### Questions

#### Theoretical questions

1. Which header sets a cookie?
2. Which header sends a cookie back?
3. What is a forward proxy versus a reverse proxy?
4. Why is `X-Forwarded-For` a claim?
5. What does `curl -v` show that a silent `curl` hides?

#### Easy practical tasks

1. Write five sentences that explain cookies and `curl -v`.
2. Make a table: `Set-Cookie` flag, one purpose. Add `Secure` and `HttpOnly`.
3. Run `curl -v https://example.com`. Write two request headers and two response headers.
4. Draw: browser, forward proxy, origin.

#### Medium practical tasks

1. Run `curl -I` on a public URL. Write status and `server` or `content-type` if present.
2. Write six sentences: a session cookie versus a token in a URL.
3. Read your browser cookie UI for one site that you use. Write two cookie names. Do not publish values.

#### Advanced practical tasks

1. Send `curl -v -x` through a local proxy that you own, or skip and write the flags you would use. Document the CONNECT line if HTTPS.
2. Write a one-page note: reverse proxy and the Host header. What the origin must trust.

---

## HTTPS preview

HTTPS is HTTP over a TLS session on TCP in the common web case. The URL scheme is `https`. The default port is 443. The client first builds TCP, then completes a TLS handshake, then sends HTTP.

TLS provides confidentiality, integrity, and authenticity when the client verifies the certificate. Topic 10 covers TLS. This section needs the preview: the padlock is not a different HTTP.

`curl -v` on HTTPS shows the certificate subject, the issuer, and the TLS version when the tool prints them. A browser shows a certificate dialog.

Cleartext HTTP on port 80 can redirect to HTTPS with `301` or `308` and a `Location` header. HSTS tells the browser to use HTTPS next time. You still need TLS configured on 443.

On the wire, a path observer sees IP addresses, ports, and often SNI (the name in the TLS handshake). The observer does not see the path or the cookie on a normal HTTPS session.

Do not send passwords on `http://` on a public path. Do not disable certificate verify to hide a lab error on a real user path. Do not think HTTPS removes the need for the Host header. TLS SNI and HTTP Host both name the site.

HTTP/3 uses QUIC on UDP. The scheme is still `https`. Topic 11 covers that mapping.

### Questions

#### Theoretical questions

1. What sits under HTTP on common HTTPS?
2. What default port does HTTPS use?
3. What three properties does TLS aim to provide (preview)?
4. What can a path observer still see on HTTPS?
5. Does HTTPS change HTTP methods?

#### Easy practical tasks

1. Write five sentences that define HTTPS as HTTP plus TLS.
2. Run `curl -v https://example.com`. Write the TLS version if shown.
3. Make a table: HTTP port 80, HTTPS port 443. Add what the capture shows.
4. Draw: TCP, TLS, HTTP as three stacked boxes.

#### Medium practical tasks

1. Compare `curl -v http://example.com` and `curl -v https://example.com` if the HTTP URL redirects. Write the status and Location.
2. Write six sentences: SNI versus Host.
3. Open the certificate view in a browser for one site. Write the name and the issuer. Do not export the private key (you do not have it).

#### Advanced practical tasks

1. Write a one-page preview: what topic 10 must add (PKI, handshake, failures).
2. Capture an HTTPS session. List fields you can read and fields you cannot read. Use a site that you may fetch.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do URL, Host, and SNI work together for one HTTPS request to a name-based server?
2. Why can keep-alive and cookies both need a long-lived client state?
3. When does a 3xx response plus Location change the next request line?
4. How does a reverse proxy change what `curl -v` shows versus what the origin process sees?
5. Why is a 200 on HTTP not the same safety as a 200 on HTTPS?

#### Easy practical tasks

1. Write a cheat sheet: URL parts, methods, status classes, Host, keep-alive, cookie, proxy, curl -v, HTTPS.
2. Run `curl -v https://example.com` and write method, status, Host, and TLS yes or no.
3. Draw a request with headers and a blank line, then a response with a short body.
4. Bookmark MDN HTTP. Write when you open it.

#### Medium practical tasks

1. Write a local HTTP/1.1 server that logs method, path, and Host. Hit it with `curl -v` and a browser.
2. Trace one page load in a browser network panel. Write two status codes and whether connections were reused.
3. Write six sentences: cookie `Secure` plus HTTPS preview.

#### Advanced practical tasks

1. Implement GET and POST on your local server. Show `Content-Length` on POST. Use only a host that you own.
2. Write a one-page map from browser address bar to TCP write: DNS name (topic 9), TLS (topic 10), then HTTP bytes.
