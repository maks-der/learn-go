# 14. Internet Routing and Hardening

## Description

The public Internet is a set of networks that exchange routes. A path can also lie, leak, or flood. This topic shows AS numbers, the BGP idea, and peering versus transit. You learn anycast and RPKI at an awareness level. You learn spoofing, SYN flood, and untrusted networks. You learn VPNs and WireGuard at a high level. You draw one HTTPS request from the browser to the origin.

Complete this topic after routing, TLS, middleboxes, and debugging.

Use one term for each concept. An autonomous system (AS) is a network with one routing policy. BGP is the routing protocol between autonomous systems. Peering is a settlement-free or policy interconnect. Transit is paid carriage to the rest of the Internet. Hardening is a control that reduces harm. Do not announce prefixes that you do not hold. Do not test floods or spoofing on a network that you do not own.

This handbook does not give attack steps, payloads, or exploit procedures. You learn what to fear and what to tighten on systems that you own.

---

## AS numbers, BGP idea, peering vs transit

An AS number (ASN) is an identifier for an autonomous system. IANA and the RIRs assign ASNs. A public ASN appears in the global BGP table. A private ASN exists for local use and must not leak to the public Internet in the usual policy.

A small site often has no ASN. The site buys transit from an ISP. The ISP AS originates or aggregates the site prefixes. A large company, a cloud, or an ISP has one or more ASNs.

BGP sends updates about prefixes. Each update has an AS_PATH: the list of ASNs that the announcement already crossed. A router prefers a path by policy, then by path length, then by other tie breaks. Hop count of IP routers inside an AS is not the same as AS_PATH length.

A loop is visible when your ASN already sits in AS_PATH. You drop that update. BGP is not OSPF. OSPF floods link state inside one organization. BGP hides internal topology and lets operators apply policy.

Peering connects two ASes so that they exchange traffic for their own prefixes (and their customers, by policy). Peering often happens at an Internet exchange or a private interconnect. Transit is a customer-provider link. The provider carries traffic to and from the rest of the Internet for a fee.

A home user has neither peering nor an ASN. The home user is a customer of an ISP that has both.

iBGP runs inside an AS. eBGP runs between ASNs. This topic only needs the names.

Do not advertise all connected routes to the Internet. Do not accept a full table on a small CPU without a plan. Do not assume the shortest AS_PATH is the best user path. Policy and congestion matter. Convergence can take time.

### Questions

#### Theoretical questions

1. What does an ASN identify?
2. What does AS_PATH list?
3. How is BGP different from OSPF in this survey?
4. What is peering in this section?
5. What is transit in this section?

#### Easy practical tasks

1. Write five sentences that define AS, BGP, peering, and transit.
2. Look up the origin AS for a public prefix (looking glass or whois). Write the ASN.
3. Make a table: home LAN, ISP, large cloud. Add one row for has a public ASN.
4. Draw: your ISP AS, a peer AS, and a transit provider (idea).

#### Medium practical tasks

1. Use `whois` or a web whois on a public IP that you may query. Write the ASN and the organization name if shown.
2. Write six sentences: two offices, no ASN, two ISP transits.
3. Find a public list of well-known ASNs (a cloud or a CDN). Write two names and numbers.

#### Advanced practical tasks

1. Write a one-page note: when a company should get an ASN (high-level, no registrar procedure dump).
2. Map one public website: name, A or AAAA, prefix, origin AS. Use public tools.

---

## Anycast and RPKI (awareness)

Anycast means many hosts share one destination IP address. Routers pick the nearest or the best host by routing metric. DNS resolvers and CDNs use anycast so that users hit a nearby site. A single address is not a single machine.

Anycast failover is a route change. Some flows break. Some flows move. You do not debug anycast with "the" server log on one rack only.

RPKI is Resource Public Key Infrastructure. It lets an address holder sign a route origin authorization (ROA). A ROA says which ASN may originate a prefix. Routers can validate BGP updates (Route Origin Validation). Invalid updates can be dropped by policy.

RPKI does not encrypt user traffic. RPKI does not stop every leak. It reduces accidental or hostile origin of a prefix by the wrong ASN. Awareness is enough: you must know the names ROA and ROV.

A bad leak or hijack can send a prefix to the wrong AS. Users see wrong sites, black holes, or interception. Operators filter, use RPKI, and talk to peers. You do not announce a test prefix on the public Internet.

Do not treat anycast as a load balancer algorithm inside one rack. Do not treat RPKI as TLS. Do not configure BGP on a device that you do not own.

### Questions

#### Theoretical questions

1. What does anycast share among many hosts?
2. Why can one resolver IP be many machines?
3. What does a ROA authorize?
4. Does RPKI encrypt HTTPS?
5. What user symptom can a route leak cause?

#### Easy practical tasks

1. Write five sentences that explain anycast and RPKI.
2. Make a table: unicast web server, anycast DNS. Add how many places the IP can live.
3. Draw: two cities, one anycast address, two hosts.
4. Write four sentences: RPKI versus TLS.

#### Medium practical tasks

1. Query a public resolver address (for example a well-known DoH or DNS IP from a public page). Write that anycast may hide the city.
2. Write six sentences: why traceroute to an anycast address can look different from two homes.
3. Read a public RPKI overview. Write four STE sentences. Do not copy long passages.

#### Advanced practical tasks

1. Write a one-page awareness note: leak, ROA, and what a looking glass can show. No attack steps.
2. Compare two looking-glass views of one prefix if public tools allow it. Write AS_PATH differences.

---

## Spoofing, SYN flood, untrusted networks

Spoofing means the sender puts a false identity in a field that a receiver might trust. A false source IP is one form. A false DNS answer is another. The field looks valid. The peer is wrong.

A SYN flood is a large number of TCP SYNs that try to fill listen state. Topic 6 said SYN uses kernel memory. Honest clients then fail to complete a handshake. The idea is resource exhaustion at the handshake. Modern stacks use SYN cookies or large backlogs. Providers absorb volume. You need the idea: listen state is finite.

An untrusted network is a path that you do not control: cafe Wi-Fi, a hotel LAN, a guest VLAN, the public Internet. Treat the path as hostile. Use TLS with verify on. Use SSH with host-key checks. Avoid cleartext protocols. Prefer a VPN when the policy requires a trusted egress.

Hardening ideas:

- do not trust a source IP as a user identity
- do not expose a debug listener to the Internet
- keep backlog and file-descriptor limits in your capacity plan
- fail closed when you cannot accept more connections
- do not send secrets on HTTP or open MQTT

This handbook does not give flood commands or spoof recipes. Do not run a flood against any host, including a lab that shares a path with other people. Do not capture on a shared LAN to read other people data.

Integrity from TLS does not stop an API replay if your handler accepts the same request twice without a nonce. Name the layer.

### Questions

#### Theoretical questions

1. What does spoofing falsify?
2. What resource does a SYN flood try to exhaust?
3. Why is a source IP a weak user identity?
4. What is an untrusted network in this section?
5. Why does TLS verify matter on cafe Wi-Fi?

#### Easy practical tasks

1. Write five sentences that define spoof, SYN flood idea, and untrusted network.
2. Make a table: threat, one hardening control. Add three rows.
3. List listen sockets on your machine (`ss` or `netstat`). Write which bind to all addresses.
4. Draw: cafe Wi-Fi, TLS to a public API, addresses visible.

#### Medium practical tasks

1. Write six sentences: a custom UDP command without a nonce (replay idea). No exploit steps.
2. Review one lab service that you own. Write whether it uses TLS and whether it trusts a source IP.
3. Write a host listen policy: localhost only versus LAN versus Internet.

#### Advanced practical tasks

1. Write a one-page threat sheet for a student API: spoof, flood idea, untrusted path. Controls only.
2. On a lab that you own, bind a service to `127.0.0.1` and show that a LAN client cannot connect. Document bind address.

---

## VPNs and WireGuard (high-level)

A VPN is a virtual private network. It encrypts and tunnels packets between a client and a gateway, or between two sites. The inner packets get a new outer header. The path observer sees the tunnel endpoints, not the inner destinations, when the tunnel works.

People use a VPN to reach a private LAN, to meet a policy egress, or to reduce exposure on an untrusted network. A VPN is not a full security program. The inner services still need auth. Split tunnel versus full tunnel changes which prefixes use the VPN.

WireGuard is a modern VPN protocol. It uses a small set of cryptographic primitives and a simple key model. Each peer has a public key. Allowed IPs define which prefixes travel in the tunnel. The common transport is UDP. This topic is high-level. You do not harden a production mesh here.

Other VPN families exist: IPsec, OpenVPN, TLS-based vendor clients. The idea stays: a tunnel plus keys plus a routing decision.

A VPN does not replace HTTPS to a public API if the API is outside the tunnel. A VPN plus clear HTTP inside a hostile inner network is still a fault. Do not disable TLS because a VPN exists unless the design says the tunnel is the only path and you accept that model.

Do not paste private keys into chat. Do not run a public VPN server without auth and updates. Do not treat a commercial VPN as anonymity. Addresses and timing still leak.

SSH tunnels and HTTPS CONNECT proxies are related ideas. They are not WireGuard.

### Questions

#### Theoretical questions

1. What does a VPN tunnel hide from a path observer?
2. What does a path observer still see?
3. What does WireGuard use to identify a peer in this section?
4. Why is a VPN not a full security program?
5. Does a VPN replace TLS to a public website?

#### Easy practical tasks

1. Write five sentences that explain VPN and WireGuard at a high level.
2. Make a table: no VPN on cafe Wi-Fi, VPN to office, HTTPS to a bank. Add what each protects.
3. Draw: client, UDP tunnel, gateway, inner LAN.
4. Write four sentences: split tunnel versus full tunnel.

#### Medium practical tasks

1. If you have a VPN that you own, write the transport (UDP or other) and whether all prefixes go into the tunnel. Skip if you have none.
2. Write six sentences: WireGuard Allowed IPs as a routing policy.
3. Compare SSH local forward and a VPN in four sentences.

#### Advanced practical tasks

1. Write a one-page design: remote admin uses WireGuard to a bastion, then SSH. Name what TLS still covers.
2. On a lab that you own, read a WireGuard quick start from official docs. Write the peer fields in STE. Do not publish keys.

---

## Draw one HTTPS request from browser to origin

This section is a single picture in words. You draw it on paper. An HTTPS request from a browser to an origin crosses almost every topic in this path.

Teaching path (typical IPv4 home user, HTTP/2 on TCP):

1. The user types a URL or clicks a link. The fragment stays in the browser.
2. The stub resolver asks a recursive resolver. DNS returns A or AAAA (topic 9). TTL starts.
3. The host looks up the routing table (topic 4). The destination is not local. The next hop is the default gateway.
4. ARP or ND finds the gateway MAC (topic 2). The Ethernet frame goes to the home switch or Wi-Fi AP.
5. The home router applies NAT and a stateful firewall (topics 4 and 12). The source IP becomes the public or CGNAT address.
6. ISP routing and BGP carry the packet across ASes (this topic). A CDN anycast address may attract the packet to a nearby city (topic 12).
7. TCP three-way handshake runs (topic 6). `ss` would show SYN-SENT then ESTABLISHED.
8. TLS handshake runs with SNI and a certificate check (topic 10). ALPN may select `h2`.
9. HTTP/2 carries the request on a stream: method, path, Host, headers, optional body (topics 8 and 11). Cookies may go out.
10. An L7 load balancer or the origin answers. The response status and body return on the same stream. The browser renders.

Variants: HTTP/3 uses QUIC on UDP and skips the TCP handshake. IPv6 may skip NAT. DoH skips clear UDP 53. A VPN wraps steps 5 to 6 in an outer tunnel.

Draw boxes. Label each box with one topic number. The picture is the exam for this path.

Do not skip DNS when the user typed a name. Do not skip TLS when the scheme is `https`. Do not label the home box as only a switch.

### Questions

#### Theoretical questions

1. Where does DNS sit relative to TCP in this path?
2. What changes at the home NAT box?
3. Which two handshakes occur before HTTP on TCP HTTPS?
4. What can a CDN change in the path?
5. How does HTTP/3 change steps 7 and 8?

#### Easy practical tasks

1. Draw the ten-step path on one page. Label topics.
2. Write five sentences that walk one step from each of DNS, NAT, TCP, TLS, and HTTP.
3. Make a table: step, tool that can show it (`dig`, `arp`, `curl -v`, capture).
4. Circle the steps that a path observer can still see on HTTPS.

#### Medium practical tasks

1. Run `dig`, `ping`, `curl -v`, and a short capture for `https://example.com`. Write one fact per tool that matches a step.
2. Write six sentences: IPv6 dual-stack variant of the same picture.
3. Add a WireGuard tunnel to the drawing. Mark which steps move inside the tunnel.

#### Advanced practical tasks

1. Write a one-page annotated diagram: browser to origin, including CGNAT or IPv6 and a CDN. No secrets.
2. Present the diagram to a peer. Time yourself. Fill gaps from topics 1 to 13.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do AS_PATH policy, anycast, and a CDN together explain two users who see two different RTTs to one name?
2. Why do RPKI, TLS, and a VPN protect three different claims (route origin, peer name, inner prefixes)?
3. What hardening still matters if BGP and TLS both work?
4. How does an untrusted LAN change the HTTPS picture without changing the origin certificate?
5. Where do SYN-flood controls sit in the ten-step picture: DNS, TCP, or HTTP?

#### Easy practical tasks

1. Write a cheat sheet: ASN, BGP, peering, transit, anycast, RPKI, spoof, SYN flood idea, untrusted net, VPN, WireGuard, HTTPS path.
2. Draw the HTTPS path from memory. Then compare with the section list.
3. Look up one origin AS for a site that you use. Write the number.
4. List three controls you will keep on cafe Wi-Fi.

#### Medium practical tasks

1. Write a table that maps each of the ten HTTPS steps to a debug tool from topic 13.
2. Write a one-page policy: listen binds, TLS verify, no cleartext, VPN when required.
3. Write six sentences: peering versus transit for a small SaaS that buys one ISP.

#### Advanced practical tasks

1. Build a full-path lab report for one `curl` to a public HTTPS URL: DNS, first hop, TLS, HTTP version, and origin AS.
2. Write a hardening checklist for a student VM that you will put on a LAN: bind, firewall, SSH keys, updates. No attack steps.
