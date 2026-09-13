# 12. Middleboxes and LAN Services

## Description

A middlebox is a device or a program on the path that is not the original client or the original server. LAN hosts also need an address, a clock, and a way to find nearby services. This topic shows stateful firewalls. You learn L4 versus L7 load balancers and CDNs. You learn DHCP, NTP, and mDNS. You learn latency versus throughput.

Complete this topic after NAT, TCP, HTTP, and TLS.

Use one term for each concept. A stateful firewall tracks flows. A load balancer spreads connections or requests. A CDN caches and terminates near users. A lease is a timed grant of configuration. NTP syncs clocks. mDNS is DNS-like name lookup on the local link. Latency is delay. Throughput is bytes per unit time. Do not mix a firewall with NAT even when one box does both. Do not mix latency with throughput.

Practice with `curl -v`, traceroute, and timings. Do not change a production firewall without a plan. Do not run a second DHCP server on an office LAN.

---

## Stateful firewalls

A packet filter can allow or drop by address, port, and protocol. A stateful firewall also keeps a table of flows. A reply that matches an allowed outbound flow can pass. A new inbound SYN to a closed or filtered port can drop.

The state key is often the 5-tuple: protocol, local IP, local port, remote IP, remote port. UDP state uses timeouts because there is no FIN. NAT boxes are stateful. Topic 4 covered NAT. A home router is often NAT plus a stateful firewall.

Stateless rules still matter for ICMP and for some anycast paths. A policy that drops all ICMP can break path MTU discovery.

Application-layer firewalls (WAF) inspect HTTP. They are still middleboxes. They can break unusual clients. Log the drop.

`curl` to a closed port fails fast. `curl` to a filtered port can hang until a timeout. That difference is a debug clue.

Do not add allow-any to hide a bug. Do not forget IPv6 rules on a dual-stack host. Do not treat a firewall as a patch for a service that listens on all interfaces with no auth.

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
3. Read your OS firewall UI name (Windows Defender Firewall or `ufw` or `nft`). Write whether it is on. Do not disable a work firewall.

#### Advanced practical tasks

1. Write a one-page host policy: allow SSH from the LAN only, allow outbound HTTPS, drop inbound others. High-level only.
2. On a lab VM that you own, add one deny rule and show that a local client fails. Then remove the rule. Document both.

---

## L4 vs L7 load balancers and CDNs

A load balancer (LB) sits in front of several backends. The client talks to one virtual address or name. The LB picks a backend.

A layer-4 LB forwards TCP or UDP by 5-tuple. It does not read HTTP. It can be fast. It cannot route by URL path or by Host header. TLS can pass through. The backend then terminates TLS.

A layer-7 LB reads HTTP. It can route by path, by header, or by cookie. It often terminates TLS. It can add `X-Forwarded-For`. An L7 LB is a reverse proxy with a pick algorithm.

Algorithms include round robin, least connections, and hash of the client address. Health checks remove a dead backend. Sticky sessions send one client to one backend when the app keeps local state. Prefer stateless apps.

A CDN is a content delivery network. It is a set of caches and TLS endpoints near users. The browser talks to a nearby point of presence. The CDN fetches from the origin on a cache miss. Static files benefit most. Personalized API calls may still go to origin.

Anycast (topic 14) often fronts a CDN. The user still types one name.

Do not put an L4 LB in front of gRPC and expect path-based HTTP routing. Do not hide all backend errors as `502` without logs. Do not use stickiness as a substitute for a shared store. A heavy health check can become the load.

### Questions

#### Theoretical questions

1. What does a load balancer hide from the client?
2. What can an L4 LB not read?
3. What can an L7 LB use to pick a backend?
4. Where does TLS terminate in a common L7 HTTPS design?
5. What does a CDN cache?

#### Easy practical tasks

1. Write five sentences that compare L4, L7, and a CDN.
2. Make a table: L4, L7, CDN. Add one job each.
3. Draw: users, CDN, origin, L7 LB, two backends.
4. Write four sentences: sticky sessions versus a shared store.

#### Medium practical tasks

1. Use `curl -v` on a public CDN-backed site. Write the `server` or `via` or `cf-` style header if any appears. Do not invent fields.
2. Write six sentences: TLS pass-through L4 versus TLS offload L7.
3. Write why `502` can mean the LB is up and the backend is down.

#### Advanced practical tasks

1. Write a one-page design: L4 for a TLS-passthrough database, L7 for an HTTP API, CDN for images.
2. On a lab that you own, run two backends and a reverse proxy that picks by path. Document Host and path rules.

---

## DHCP, NTP, mDNS

DHCP (IPv4) gives a host a temporary configuration: IPv4 address, subnet mask, default gateway, and often DNS resolvers. The grant is a lease. The host must renew the lease before the time ends if it wants to keep the address.

The common IPv4 dance: discover, offer, request, acknowledge. A relay can forward to a DHCP server on another subnet. When the lease ends and renewal fails, the host must stop using that address. Some hosts keep a link-local IPv4 address (`169.254/16`) if DHCP fails. That address does not replace a gateway.

DHCPv6 can assign IPv6 addresses or only options. Topic 3 covered SLAAC versus DHCPv6. Do not expect DHCPv4 to set IPv6.

Reservations bind a MAC or a client identifier to a stable address. The grant is still a lease. A lease is not a security grant. Anyone who can use the LAN can often get a lease.

NTP is the Network Time Protocol. A host asks time servers and adjusts its clock. Correct time matters for TLS certificates, for logs, and for signed tokens. A large clock error looks like an expired certificate. NTP sets UTC. Time zone is a display offset. Do not fix TLS by changing the time zone only.

mDNS is multicast DNS on the local link. Hosts publish names such as `printer.local`. Bonjour is a common implementation. mDNS does not replace global DNS. It does not work across routers unless a gateway exists. Do not publish internal names to a hostile LAN.

Do not run two uncoordinated DHCP servers on one LAN. Do not point a production fleet at a random public NTP server at huge scale. Do not disable time sync to hide a test.

### Questions

#### Theoretical questions

1. What four items does a typical IPv4 DHCP lease include?
2. What must a host do before a lease ends if it wants to keep the address?
3. Why does TLS need a correct clock?
4. What kind of names does mDNS publish?
5. Does mDNS replace global DNS?

#### Easy practical tasks

1. Run `ipconfig /all` or the Unix equivalent. Write IPv4, DHCP server, and lease times if shown.
2. Write five sentences that explain lease, NTP, and mDNS.
3. Make a table: DHCP, NTP, mDNS. Add one protocol job each.
4. Draw: host, DHCP server, NTP server, `.local` name.

#### Medium practical tasks

1. Find the lease expiry on your OS. Write how long the lease still lasts.
2. Write six sentences: two DHCP servers with overlapping pools.
3. Check whether your clock syncs (OS time setting). Write the service name if shown.

#### Advanced practical tasks

1. On a lab LAN that you own, watch a reconnect. Capture DHCP if you can. Write the message names that the tool shows.
2. Write a one-page office plan: DHCP pool, NTP source, and whether mDNS is allowed.

---

## Latency vs throughput

Latency is how long one action waits. Round-trip time (RTT) is a common latency measure. `ping` samples RTT. A TLS handshake needs several RTTs. A short request on a long-delay path is latency-bound.

Throughput is how many bytes move per unit time. A large file on a wide path is throughput-bound when the window and the path allow it. Topic 6 covered the window and congestion control.

Bandwidth is the capacity of the link. Throughput is what you achieve. Throughput cannot exceed the bottleneck capacity. Headers and loss reduce it.

A 1 Gbit/s link with 200 ms RTT can still feel slow for a small API call. The same link can move a large file at high throughput after slow start. Users mix the words. You must not.

Bandwidth-delay product (BDP) is capacity times RTT. The TCP window must be large enough to fill the pipe. Window scaling exists for that reason.

Nagle and delayed ACK can add latency on small writes. Measure before you tune. Do not disable congestion control to chase throughput on the public Internet.

Do not add more parallelism without a measurement. Extra connections help HTTP/1.1 latency for many small files. They also add handshake cost.

### Questions

#### Theoretical questions

1. What is latency in this section?
2. What is throughput?
3. Why can a fast link still feel slow for a small request?
4. What is bandwidth-delay product for?
5. Why must you not mix the two words?

#### Easy practical tasks

1. Write five sentences that define latency and throughput with one example each.
2. Ping your gateway and a public host. Write two RTTs.
3. Make a table: small API, large download. Mark latency-bound or throughput-bound as the usual case.
4. Draw a long fat pipe and a window that does not fill it.

#### Medium practical tasks

1. Time `curl -w` total time and size for a small page and a larger file if you may fetch both. Write which metric dominates.
2. Write six sentences: 20 ms RTT versus 200 ms RTT for a four-way TLS handshake idea.
3. Find your link speed in the OS. Write why that number is not the `curl` throughput.

#### Advanced practical tasks

1. Compute BDP on paper for 100 Mbit/s and 50 ms RTT. Write the byte count.
2. Write a one-page note: when to optimize RTT (CDN, keep-alive) versus when to optimize window and loss.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a stateful firewall, NAT, and an L7 LB all keep state, and how do the keys differ?
2. Why can a CDN improve latency for static files and still leave origin latency for a personalized API?
3. How do DHCP, NTP, and mDNS fail in three different user stories (no address, TLS expire, no printer name)?
4. When is a hang at `curl` a firewall symptom rather than a DNS symptom?
5. How do RTT and window together set throughput for one TCP flow?

#### Easy practical tasks

1. Write a cheat sheet: stateful FW, L4, L7, CDN, DHCP, NTP, mDNS, latency, throughput, BDP.
2. Print your lease info and your clock sync status. Write one sentence each.
3. Draw the path: client, firewall, CDN, LB, backend.
4. Run `ping` and one `curl -w` timing. Write RTT versus total time.

#### Medium practical tasks

1. Write a design review: which middleboxes your daily work path likely has. Guess from headers and traceroute.
2. Write a fault tree: no LAN IP, wrong time, cannot resolve `.local`.
3. Write six sentences: L4 pass-through TLS and what a WAF cannot see.

#### Advanced practical tasks

1. Write a one-page lab network: DHCP, NTP, firewall policy, and an L7 proxy. Use a lab that you own.
2. Measure one download throughput and one ping RTT. Compute whether the flow is near BDP or latency-bound.
