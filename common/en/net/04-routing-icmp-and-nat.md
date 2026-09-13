# 4. Routing, ICMP, and NAT

## Description

Routing is the choice of the next hop for a packet. ICMP reports network errors and supports `ping` and `traceroute`. NAT translates addresses at an edge. This topic shows the routing table and longest prefix match. You learn ICMP, `ping`, and `traceroute`. You get a survey of static and dynamic routing. You learn source NAT, port forwarding, and peer-to-peer pain. You learn carrier-grade NAT at an awareness level.

Complete this topic after IPv4 and IPv6 addresses. Complete this topic before you study UDP and TCP in depth.

Use one term for each concept. A route is a prefix plus a next hop or an outgoing interface. A routing table is the list of routes that the device uses. ICMP is a control protocol next to IP. Source NAT changes the source address of packets that leave the LAN. Port forwarding maps a public port to a private host. Do not mix a routing table with an ARP table. Do not mix NAT with a full firewall.

You already know the default gateway. This topic puts that gateway in a table with other prefixes.

---

## Routing table and longest prefix match

A routing table maps a destination prefix to a next hop and an interface. The kernel looks up the destination IP of each packet. The kernel chooses one route. The kernel then uses ARP or neighbor discovery for the next hop if the next hop is on a LAN.

A host table is often small: the local subnet, the default route, loopback, and maybe a VPN. A router table can be large. An Internet router can hold hundreds of thousands of prefixes.

Connected routes appear when an interface has an address and a mask. The destination on that subnet is local. The next hop is not a router. The host sends to the destination MAC.

When two routes match the same destination, the device chooses the match with the longest prefix (the most specific mask). A `/32` wins over a `/24`. A `/24` wins over `/0`.

Example: destination `192.168.1.50`. Routes: `0.0.0.0/0` via G1, `192.168.1.0/24` via interface LAN, `192.168.1.50/32` via G2. The host route `/32` wins.

Longest prefix match is not the same as lowest metric. If two routes have the same prefix length, the device uses metric, preference, or equal-cost sharing. Details depend on the OS.

IPv6 uses the same idea. `::/0` is the default. A `/128` is a host route.

On Linux, `ip route` shows IPv4. `ip -6 route` shows IPv6. On Windows, `route print` shows both in one view.

Do not edit a production routing table without a plan. Do not add a default route to a wrong next hop. You can cut your own management path. Do not assume that a smaller metric beats a longer prefix. The prefix length is first in the usual model.

### Questions

#### Theoretical questions

1. What does a routing table map?
2. What is a connected route?
3. What does longest prefix match mean?
4. Which wins for a destination in both prefixes: `/24` or `/16`?
5. What is the IPv6 default prefix?

#### Easy practical tasks

1. Print your IPv4 routing table. Copy the default route line.
2. Write five sentences that explain a connected route on your LAN.
3. Make a table: destination `10.1.2.3`, routes `/8`, `/16`, `/24`, `/0`. Mark the winner.
4. Draw a host, a LAN prefix, and a default next hop.

#### Medium practical tasks

1. Identify loopback, connected, and default routes in your table. Label each line.
2. On paper, pick a winner among `10.0.0.0/8`, `10.1.0.0/16`, and `0.0.0.0/0` for `10.1.9.8`.
3. Write six sentences: what happens when the table has no matching route and no default.

#### Advanced practical tasks

1. Add a host route to a test address on a lab machine that you own, then remove it. Record the commands and `ping` before and after.
2. Write a one-page field guide for `ip route` or `route print`. Name each column that you use.

---

## ICMP, `ping`, `traceroute`

ICMP is the Internet Control Message Protocol. IPv4 uses ICMP. IPv6 uses ICMPv6. Hosts and routers send ICMP to report errors and to support tools. ICMP is not a transport for your application data in the usual design.

`ping` sends an echo request. The target can send an echo reply. A reply proves that this probe reached a stack that answers ICMP. A missing reply is not always "host down." A firewall can drop echo. The host can still serve HTTPS.

`traceroute` (Linux) and `tracert` (Windows) map hops. The tool sends probes with a rising TTL (IPv4) or hop limit (IPv6). Each hop that decrements to zero can send Time Exceeded. The source learns that hop address. Stars mean no ICMP reply. The path can still work for TCP.

Path MTU discovery uses ICMP Packet Too Big (or IPv4 fragmentation needed). A policy that drops all ICMP can break large transfers.

ICMPv6 is required for neighbor discovery. A drop-all-ICMP rule on IPv6 is a serious fault.

Do not treat a successful ping as proof that a TCP port is open. Do not treat traceroute stars as a broken Internet. Do not flood ping at a host that you do not own.

### Questions

#### Theoretical questions

1. What job does ICMP have in this section?
2. What does a ping reply prove?
3. Why can ping fail when a website works?
4. How does traceroute learn a hop address?
5. Why can a drop-all-ICMP policy hurt large packets?

#### Easy practical tasks

1. Ping your gateway. Write the RTT if you get a reply.
2. Write five sentences that explain echo request and echo reply.
3. Run `traceroute` or `tracert` to a public host. Write the first two hops.
4. Make a table: ping success, ping fail, traceroute stars. Add one meaning each.

#### Medium practical tasks

1. Ping a public HTTPS site by name and by address if DNS works. Write whether ICMP is allowed.
2. Write six sentences: TTL expiry and Time Exceeded.
3. Compare IPv4 ping and IPv6 ping to `127.0.0.1` and `::1`.

#### Advanced practical tasks

1. Capture one ping. Write ICMP type and code if the tool shows them.
2. Write a one-page note: when ping is useful and when you must use `curl` or a TCP probe instead.

---

## Static vs dynamic routing (survey)

A static route is a route that an administrator writes. The default gateway on a home host is often a static or DHCP-installed default. A static route does not learn a new path when a link dies, unless a human or a script changes it.

A dynamic routing protocol lets routers exchange prefix information. Inside one organization, OSPF and IS-IS are common link-state protocols. They flood topology and compute paths. RIP is an old distance-vector protocol. You may see it in labs.

Between organizations on the Internet, BGP is the usual protocol. Topic 14 covers BGP at a high level. This section only names it.

Dynamic routing needs extra CPU, extra configuration, and a trust model. A small two-router lab can stay static. A large network needs dynamic routing so that a single link failure does not wait for a human.

DHCP can install a default route on a host. That route is not OSPF. The host still has a static-looking table. The source of the route is DHCP.

Do not run a dynamic protocol on a home LAN without a plan. Do not advertise private prefixes to the public Internet. Do not mix "the route exists" with "the path has capacity."

### Questions

#### Theoretical questions

1. Who writes a static route?
2. What does a dynamic protocol exchange?
3. Name one interior protocol from this survey.
4. Which protocol does this handbook name for the Internet edge?
5. How can DHCP change a host routing table?

#### Easy practical tasks

1. Write five sentences that compare static and dynamic routing.
2. Make a table: home host, office campus, Internet ISP. Add static or dynamic as the usual case.
3. Draw two routers and one static default toward the Internet.
4. List OSPF, RIP, and BGP in one column with one job each.

#### Medium practical tasks

1. Find the source of your default route (DHCP, static, or unknown). Write the evidence.
2. Write six sentences: a static default when the ISP link dies.
3. Read a public one-page OSPF versus BGP comparison. Write four STE sentences. Do not copy long passages.

#### Advanced practical tasks

1. In a lab that you own, add a static route and a more specific static route. Show longest prefix match with ping.
2. Write a one-page survey: when a team stays static and when it adopts OSPF. No vendor CLI dump.

---

## Source NAT, port forwarding, peer-to-peer pain

Source NAT rewrites the source IP of a packet that leaves the inside network. Masquerade is source NAT that uses the current address of the outbound interface. Home routers do this job so that many private hosts share one public IPv4 address.

The box also rewrites the source port when two inside hosts would collide. The box stores a mapping: inside IP, inside port, outside IP, outside port, protocol, and often the remote peer. Return packets match the mapping. The box rewrites the destination back to the inside host.

The mapping is stateful. If the mapping expires, a late reply does not reach the inside host. Idle timeouts differ for UDP and TCP.

Port forwarding (destination NAT) maps a public port to a private host and port. An inbound connection to the public address and that port goes to the inside server. Without a mapping or an outbound state row, a new inbound connection does not reach a private host.

Peer-to-peer applications want two NATed hosts to talk. Neither host has a stable inbound mapping. NAT traversal (STUN, TURN, hole punching) tries to create mappings. It often fails when both sides sit behind strict NAT. Relays then carry the media.

NAT is not a full firewall. A stateful firewall can sit in the same box. Keep the two jobs separate in your notes. NAT is not a proxy. A proxy understands the application.

Do not treat NAT as encryption. Do not expect inbound SSH to a LAN host without a forward or a tunnel. Do not confuse source NAT with dest NAT.

### Questions

#### Theoretical questions

1. What field does source NAT change on the outbound packet?
2. Why can the NAT box also change the source port?
3. What does port forwarding map?
4. Why do two home users have peer-to-peer pain?
5. How is NAT different from a proxy in this section?

#### Easy practical tasks

1. Write your LAN IPv4 and your public IPv4. Label NAT at the edge.
2. Write five sentences that describe an outbound HTTPS request through source NAT.
3. Draw before and after headers: source IP inside, source IP outside.
4. Make a table: source NAT, port forwarding. Add direction of the first packet.

#### Medium practical tasks

1. Write six sentences: two PCs that browse at the same time through one public IP.
2. Find the port-forward UI on a home router that you own, or write "no access." Do not change a work router.
3. Write why a game or a call might need a relay when both peers use NAT.

#### Advanced practical tasks

1. Write a one-page mapping table for three inside hosts that all talk to `8.8.8.8:53` over UDP. Show unique outside ports.
2. On a lab router that you own, add a port forward to a local server, test from another network if you can, then remove it. Document both.

---

## Carrier-grade NAT (awareness)

Carrier-grade NAT (CGNAT or CGN) is NAT in the ISP network. Many customers share a pool of public IPv4 addresses. Your home router may still do NAT. Then you have two NAT layers: home NAT plus CGNAT.

CGNAT makes inbound port forwarding on your home box fail for the public Internet. The public address that a web page shows can belong to a shared pool. A port forward on the home WAN IP never sees the packet.

ISPs use CGNAT because public IPv4 addresses are scarce. IPv6 avoids this sharing when both ends have IPv6. Dual stack can use IPv6 for peer-to-peer while IPv4 stays behind CGNAT.

Awareness facts:

- you may not get a dedicated public IPv4 address
- logs on a server show a shared address
- abuse reports can name many customers at once
- IPv6 or a tunnel can restore inbound reachability

Do not assume that "I have a home router" means you have a unique public IPv4 address. Do not buy a port-forward plan before you test whether the WAN address is public and dedicated.

This section is awareness. You do not configure an ISP CGNAT box.

### Questions

#### Theoretical questions

1. Where does CGNAT sit relative to a home LAN?
2. Why do ISPs use CGNAT?
3. Why can a home port forward fail under CGNAT?
4. How can IPv6 reduce CGNAT pain?
5. Why can one public IPv4 in a log name many customers?

#### Easy practical tasks

1. Write five sentences that explain two NAT layers.
2. Compare your router WAN IPv4 (if you can see it) with a "what is my IP" page. Write whether they match.
3. Make a table: dedicated public IPv4, CGNAT. Add inbound forward result.
4. Draw: LAN, home NAT, ISP CGNAT, Internet.

#### Medium practical tasks

1. Write six sentences: a home camera that needs inbound access on IPv4-only CGNAT.
2. Check whether your ISP documents IPv6 or CGNAT. Write one fact from a public page.
3. Explain why a VPN outbound from the LAN can still work under CGNAT.

#### Advanced practical tasks

1. Write a one-page awareness note for a teammate: how to detect CGNAT without ISP CLI. Use address comparison only.
2. Design on paper an IPv6-first inbound test versus an IPv4 relay. No purchase steps.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do longest prefix match and source NAT both change what a packet does at an edge?
2. Why must a debug story name ICMP filter, routing miss, and NAT mapping as three different causes?
3. When is a static default enough, and when does a campus need a dynamic protocol?
4. What does a public server log show for a host behind source NAT, and what does it hide?
5. How does CGNAT change the meaning of "open a port on my router"?

#### Easy practical tasks

1. Write a cheat sheet: route, LPM, ICMP, ping, traceroute, static, OSPF, BGP name only, SNAT, forward, CGNAT.
2. Print the routing table and ping the gateway. Write one sentence that links the default route to the ping.
3. Draw a packet from `192.168.1.10` to `8.8.8.8` through home NAT. Show addresses before and after.
4. Run traceroute to a public resolver. Mark the hop that leaves your LAN.

#### Medium practical tasks

1. Write a fault tree: cannot reach the Internet. Include default route, gateway ping, NAT, and ICMP filter.
2. Compare IPv4 default route and IPv6 default route on your machine. Write whether both exist.
3. Write six sentences that place ARP (topic 2) after the route lookup for a LAN next hop.

#### Advanced practical tasks

1. Write a one-page lab: static host route versus default, then a NAT mapping story for the same host. Use a lab that you own.
2. Map one HTTPS fetch: route lookup, first hop, NAT rewrite, and whether ICMP was needed. Use tools from this topic.
