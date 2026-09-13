# 17. Middleboxes and Performance

## Description

A middlebox is a device or a program on the path that is not the original client or the original server. This topic shows stateful firewalls, load balancers, CDNs, Nagle and delayed ACK, bandwidth-delay product, latency versus throughput, and head-of-line delay at many layers.

Use one term for each concept. A stateful firewall tracks flows. A load balancer spreads connections or requests. A CDN caches and terminates near users. Complete this topic after NAT, TCP, HTTP, and TLS.

Practice with `curl -v`, traceroute, and timings. Do not change a production firewall without a plan. Measure before you tune.

---

## Firewalls (stateful)

A packet filter can allow or drop by address, port, and protocol. A stateful firewall also keeps a table of flows. A reply that matches an allowed outbound flow can pass. A new inbound SYN to a closed port can drop.

The state key is often the 5-tuple (protocol, local IP, local port, remote IP, remote port). UDP state uses timeouts because there is no FIN. NAT boxes are stateful. Topic 7 covered NAT. A home router is often NAT plus a stateful firewall.

Stateless rules still matter for ICMP and for some anycast paths. A policy that drops all ICMP can break path MTU discovery. Topic 6 and topic 19 cover that cost.

Do not add "allow any any" to hide a bug. Do not forget IPv6 rules on a dual-stack host. Do not treat a firewall as a patch for a service that listens on all interfaces with no auth.

Application-layer firewalls (WAF) inspect HTTP. They are still middleboxes. They can break unusual clients. Log the drop.

`curl` to a closed port fails fast. `curl` to a filtered port can hang until a timeout. That difference is a debug clue.

### Questions

#### Theoretical questions

1. What extra table does a stateful firewall keep?
2. What is a common state key for TCP?
3. Why does UDP state need a timeout?
4. Why can a drop-all-ICMP policy hurt large transfers?
5. How can a hang versus an instant fail hint at a filter?

#### Easy practical tasks

1. Write five sentences that compare stateless and stateful filters.
2. Make a table: new inbound SYN, inbound reply to your outbound. Add allow or drop in a typical home policy.
3. Draw: LAN host, firewall, Internet server, state row.
4. List IPv4 and IPv6 as two rule families.

#### Medium practical tasks

1. Time `curl` to a closed local port and to a filtered public port if you have a known filter that you may test. Write both behaviors. Do not scan a range.
2. Write six sentences: IPv4 works, IPv6 is filtered, dual-stack browser.
3. Read your OS firewall UI name (Windows Defender Firewall or `ufw`/`nft`). Write whether it is on. Do not disable a work firewall.

#### Advanced practical tasks

1. Write a one-page host policy: allow SSH from the LAN only, allow outbound HTTPS, drop inbound others. High-level only.
2. On a lab VM that you own, add one deny rule and show that a local client fails. Then remove the rule. Document both.

---

## Load balancers (L4 vs L7)

A load balancer (LB) sits in front of several backends. The client talks to one virtual address or name. The LB picks a backend.

A layer-4 LB forwards TCP or UDP by 5-tuple. It does not read HTTP. It can be fast. It cannot route by URL path or by `Host` header. TLS can pass through. The backend then terminates TLS.

A layer-7 LB reads HTTP. It can route by path, by header, or by cookie. It often terminates TLS. It can add `X-Forwarded-For`. Topic 11 covered reverse proxies. An L7 LB is a reverse proxy with a pick algorithm.

Algorithms include round robin, least connections, and hash of the client address. Health checks remove a dead backend. Sticky sessions send one client to one backend when the app keeps local state. Prefer stateless apps.

Do not put an L4 LB in front of gRPC and expect path-based HTTP routing. Do not hide all backend errors as `502` without logs. Do not use stickiness as a substitute for a shared store.

Health checks must use a cheap path. A heavy health check can become the load.

### Questions

#### Theoretical questions

1. What does a load balancer hide from the client?
2. What can an L4 LB not read?
3. What can an L7 LB use to pick a backend?
4. Where does TLS terminate in a common L7 HTTPS design?
5. Why are sticky sessions a design smell for a new app?

#### Easy practical tasks

1. Write five sentences that compare L4 and L7.
2. Make a table: L4, L7. Add rows for HTTP path routing and TLS terminate.
3. Draw: clients, one VIP, three backends.
4. Write four sentences: health check versus user traffic.

#### Medium practical tasks

1. Read a public cloud or nginx load-balancer page. Write four STE sentences. No long config copy.
2. Write six sentences: `502` from an LB when all backends are down.
3. Explain `X-Forwarded-For` trust again in four sentences for an L7 LB.

#### Advanced practical tasks

1. Run two local backends and one reverse proxy on a machine that you own. Fetch the proxy URL. Write which backend answered if you can tell.
2. Write a one-page choice sheet: L4 pass-through versus L7 terminate for an HTTPS API.

---

## CDN

A CDN is a content delivery network. Many edge nodes store copies of static objects (and sometimes more). The user name often points to the CDN. DNS or anycast sends the client to a nearby edge. Topic 18 covers anycast.

The origin is the source of truth. The CDN fetches from the origin on a cache miss. Cache-Control and CDN rules decide how long an object stays. A purge drops an object early.

TLS can terminate at the edge. The edge-to-origin path can use a new TLS session. You must still protect the origin (allow-lists, authenticated pulls).

A CDN helps throughput for large files and helps delay for small static files when the edge is close. A CDN does not fix a slow origin API by itself if every request is dynamic and uncacheable.

Do not cache personalized HTML without a plan. Do not forget that a purge is not instant everywhere. Do not point a DNS name at a CDN and leave the origin open to the world if the design assumed "only the CDN may fetch."

`curl -v` can show the CDN in headers (`cf-ray`, `x-cache`, `age`, or vendor names). Those headers are hints, not a standard.

### Questions

#### Theoretical questions

1. What does an edge cache store?
2. What happens on a cache miss?
3. Who can terminate TLS in a common CDN design?
4. Why can a fully dynamic API see little CDN gain?
5. Why must the origin still have access control?

#### Easy practical tasks

1. Write five sentences that explain origin, edge, and miss.
2. Run `curl -v` on a public static site. Write any cache or CDN header that you see.
3. Make a table: cached image, uncached API POST. Add one row for CDN help.
4. Draw: user, nearest edge, origin.

#### Medium practical tasks

1. Compare `Age` or `x-cache` on two `curl` fetches of the same URL. Write what changed.
2. Write six sentences: a purge that is slow in one region.
3. Find the DNS name of a CDN target for a public site (`dig` CNAME). Write the CNAME if present.

#### Advanced practical tasks

1. Write a one-page cache policy: images 1 day, HTML 60 seconds, POST never. High-level.
2. Map `curl -v` times (`time_namelookup`, `time_connect`, `time_total` if you use `-w`) for a CDN URL and for a raw origin if you have one that you may fetch.

---

## Nagle, delayed ACK (classic pitfalls)

Nagle's algorithm waits to send a small TCP segment if unacknowledged data is in flight. The goal is to avoid many tiny packets. Delayed ACK waits a short time (often up to 200 ms) before it sends an ACK, in the hope that a reply will piggyback.

Together they can stall a request-response protocol that sends small writes and waits. Each side waits for the other. The user sees a 200 ms gap on every message. Games and some RPC stacks turn Nagle off (`TCP_NODELAY`) for that reason.

Do not disable Nagle on a bulk file transfer without a measurement. Many small packets can waste header bytes. Do not blame the disk when a capture shows 200 ms ACK delays.

HTTP/2 and QUIC change the framing. The classic pitfall is still common on HTTP/1.1 and on custom TCP protocols that write a small header, then wait.

Wireshark can show the time between a small request and the next ACK. Look at timestamps.

This section is a classic pitfall. Modern stacks still implement these knobs. Know the names.

### Questions

#### Theoretical questions

1. What does Nagle delay?
2. What does delayed ACK delay?
3. Why can the two together stall a small RPC?
4. What socket option turns Nagle off on many OSes?
5. When can many small packets be worse than a short wait?

#### Easy practical tasks

1. Write five sentences that explain the 200 ms stall story.
2. Make a table: Nagle on, `TCP_NODELAY`. Add one row for a chat message, one row for a large download.
3. Draw a timeline: small write, wait, delayed ACK, next write.
4. Write four sentences: piggyback ACK as the good case.

#### Medium practical tasks

1. Find `TCP_NODELAY` in your language stdlib docs. Write the function or flag name.
2. Write six sentences: a protocol that writes 1 byte, then waits for 1 byte.
3. In a local echo that you own, time many tiny messages. Write whether you can see a floor near 200 ms. If not, write that modern stacks may ACK sooner.

#### Advanced practical tasks

1. Write a one-page note: when to set `TCP_NODELAY` and when to keep Nagle.
2. Capture a tiny request-response on a lab TCP app. Mark ACK times. Use only your machine.

---

## Bandwidth-delay product

Bandwidth-delay product (BDP) is bandwidth times round-trip time. The result is a number of bytes. It is the amount of data that can be in flight if you want to fill the pipe.

Example idea: 100 Mbit/s and 40 ms RTT. Convert units with care. The window must be large enough or the sender stays idle while it waits for ACKs.

TCP flow control and congestion control cap in-flight data. Topic 9 covered windows. A long fat network (high bandwidth and high delay) needs a large window. Autotuning on modern OSes often helps. A tiny application buffer can still cap the rate.

Do not raise a window to "max" on a lossy path without congestion control. Do not confuse bits per second with bytes in the BDP formula. Convert first.

A satellite path has a large delay. A LAN has a small delay. The same bandwidth needs different in-flight bytes.

`ping` gives RTT. A speed test gives bandwidth. BDP needs both.

### Questions

#### Theoretical questions

1. What two quantities does BDP multiply?
2. What unit should you use for the product when you size a TCP window?
3. Why does a long RTT need more data in flight to fill a link?
4. Which TCP limits cap in-flight data?
5. Why is a LAN BDP often small?

#### Easy practical tasks

1. Write five sentences that explain BDP as "bytes in the pipe."
2. Measure `ping` RTT to a LAN host and to a public host. Write both.
3. Make a table: low RTT LAN, high RTT WAN. Add one row for in-flight bytes at the same bandwidth (idea).
4. Write four sentences: idle sender that waits for ACK.

#### Medium practical tasks

1. Compute a BDP on paper: pick a bandwidth and an RTT from your pings. Show the unit conversion.
2. Write six sentences: a 1 Gbit LAN versus a 100 Mbit path with 100 ms RTT.
3. Find whether your OS documents TCP window autotune. Write one sentence.

#### Advanced practical tasks

1. Write a one-page worksheet: RTT, bandwidth, BDP, and a window size check.
2. Time a large download and record RTT. Write whether the throughput is far below bandwidth and list two possible caps (window, loss, server).

---

## Latency vs throughput

Latency is how long one action waits. Round-trip time is a latency. Time to first byte is a latency. Throughput is how many bits or transactions you finish per unit time.

A link can have high throughput and high latency (a full airplane of disks). A link can have low latency and low throughput (a serial console). Users mix the two words. You must not.

Many small serial RPCs suffer latency. One large stream suffers throughput and window limits. Parallel requests hide latency if the work is independent. HTTP/2 multiplexing is one form of that idea.

Do not optimize throughput when the user waits on the first token. Do not add more bandwidth when the app does twenty serial DNS plus TLS plus tiny calls. Measure the critical path.

`curl -w` can print `time_namelookup`, `time_connect`, `time_appconnect`, `time_starttransfer`, and `time_total`. Those numbers split DNS, TCP, TLS, and wait for the first body byte.

Queues add latency. A full bufferbloat queue can raise RTT when a large download shares a home link. That is delay, not always loss.

### Questions

#### Theoretical questions

1. What is latency in one sentence?
2. What is throughput in one sentence?
3. Why can a high-bandwidth path still feel slow?
4. How do serial RPCs interact with RTT?
5. Which `curl -w` time is closest to "first byte"?

#### Easy practical tasks

1. Write five sentences that keep the two words apart.
2. Run `curl -w` with a format that prints the times above on `https://example.com`. Write the numbers.
3. Make a table: action, latency or throughput. Add page first paint and a 1 GB file copy.
4. Draw a critical path: DNS, TCP, TLS, request, first byte.

#### Medium practical tasks

1. Compare `time_total` for HTTP/1.1 and default `curl` on one URL. Write both.
2. Write six sentences: twenty serial API calls on a 200 ms RTT path.
3. Start a large download and ping the gateway if you may. Write whether RTT rose (bufferbloat idea). Stop the download.

#### Advanced practical tasks

1. Write a one-page method: how you decide to attack latency first or throughput first for a given bug report.
2. Build a small timing table for one page: DNS, TCP, TLS, TTFB, total. Use `curl -w` or a browser panel.

---

## Head-of-line at every layer

Head-of-line (HOL) blocking means a first item holds later items. Topic 9 covered TCP HOL. Topic 14 covered HTTP/2 on TCP versus HTTP/3. HOL also appears in other places.

A switch or a router queue: one large flow can fill a buffer. Other flows wait. Fair queues and AQM try to reduce that wait. You do not configure a core router here. You must know that "the link is not full" can still hide a queue.

A lock in an application: one mutex holds all handlers. The network is fine. The process is not.

HTTP/1.1 on one connection: one slow response holds the next request. That is application HOL on a sequential protocol.

A disk or a thread pool can be a HOL site. Debug must name the layer. A packet capture that looks healthy plus a slow handler is a clue.

Do not buy more bandwidth when the HOL site is a single-thread lock. Do not enable HTTP/3 when the HOL site is your database.

Ask: which queue, which lock, which stream?

### Questions

#### Theoretical questions

1. What is HOL in one general sentence?
2. Name three layers where HOL can appear.
3. How can a capture look fine when the user still waits?
4. Why does more bandwidth not fix a mutex HOL?
5. How does HTTP/1.1 sequential use create HOL?

#### Easy practical tasks

1. Write five sentences that list TCP HOL, HTTP/1.1 HOL, and app-lock HOL.
2. Make a table: layer, HOL example. Add four rows.
3. Draw a queue of packets with the first one stuck.
4. Write four questions you ask to find the layer.

#### Medium practical tasks

1. Write six sentences: a CDN and HTTP/3 that still hit a single-thread origin.
2. From a real slow page that you may test, guess one HOL site and one test that would confirm it.
3. Compare a browser waterfall (many files) with one `curl` of the HTML. Write what the waterfall adds.

#### Advanced practical tasks

1. Write a one-page HOL map for a web request: DNS, TCP, TLS, HTTP, app, DB. One HOL risk each.
2. In a local server that you own, add an artificial global lock or a sleep on one handler. Show that other requests wait. Then remove it. Document the test.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a mobile fetch of a static image through a CDN and an L7 edge. Use firewall state, TLS terminate, and cache hit or miss.
2. How do BDP, latency, and throughput work together when you size a long path?
3. Why can Nagle plus delayed ACK look like "the server is slow" in an APM tool?
4. When is an L4 load balancer enough, and when must you use L7?
5. A teammate says "add a CDN and open the firewall any any." Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: stateful firewall, L4, L7, CDN, Nagle, BDP, latency, throughput, HOL.
2. Run `curl -w` and `curl -v` on one HTTPS URL. Save both. Mark CDN or cache headers if any.
3. Draw one picture: client, firewall, LB, two backends, CDN edge as an optional box.
4. List which knobs this topic said not to turn first (`TCP_NODELAY`, allow any, more bandwidth).

#### Medium practical tasks

1. Write a troubleshooting flow: hang, `502`, 200 ms stalls, low throughput on a fat long path.
2. Measure RTT and a download rate to one public file. Compute a rough BDP and compare it with the rate in words.
3. Explain in ten steps how a beginner uses `curl` times before they change a firewall.

#### Advanced practical tasks

1. Build a glossary of 16 terms from this topic. Each entry: term, one sentence, one tool or header.
2. Write a lab on a machine that you own: two backends, one proxy, one slow handler. Record timings. No production changes.
