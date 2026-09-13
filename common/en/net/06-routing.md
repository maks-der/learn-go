# 6. Routing

## Description

Routing is the choice of the next hop for a packet. A host and a router both use a routing table. This topic shows longest prefix match, host routes, default routes, ICMP, traceroute, and a short view of static and dynamic routing.

Use one term for each concept. A route is a prefix plus a next hop or an outgoing interface. A routing table is the list of routes that the device uses. Complete this topic before you study NAT.

You already know the default gateway. This topic puts that gateway in a table with other prefixes.

---

## Routing table

A routing table maps a destination prefix to a next hop and an interface. The kernel looks up the destination IP of each packet. The kernel chooses one route. The kernel then ARPs or uses neighbor discovery for the next hop if the next hop is on a LAN.

A host table is often small: the local subnet, the default route, loopback, and maybe a VPN. A router table can be large. An Internet router can hold hundreds of thousands of prefixes.

Each route has a prefix length, a next hop, an interface, and a metric or preference. Some tables also show the source of the route: connected, static, DHCP, or a protocol name.

Connected routes appear when an interface has an address and a mask. The destination on that subnet is local. The next hop is not a router. The host sends to the destination MAC.

Do not edit a production routing table without a plan. Do not add a default route to a wrong next hop. You can cut your own management path.

On Linux, `ip route` shows IPv4. `ip -6 route` shows IPv6. On Windows, `route print` shows both in one view. Read the columns before you change anything.

### Questions

#### Theoretical questions

1. What does a routing table map?
2. What is a connected route?
3. Why is a host table often smaller than a router table?
4. What does the kernel do after it chooses a next hop on a LAN?
5. Name two commands that print a routing table.

#### Easy practical tasks

1. Print your IPv4 routing table. Copy the default route line.
2. Write five sentences that explain a connected route on your LAN.
3. Make a table: destination prefix, next hop, interface. Add two rows from your machine.
4. Draw a host, a LAN prefix, and a default next hop. Link them to table rows.

#### Medium practical tasks

1. Identify loopback, connected, and default routes in your table. Label each line.
2. Compare IPv4 and IPv6 tables. Write one difference that you see.
3. Write six sentences: what happens to an outgoing packet when the table has no matching route and no default.

#### Advanced practical tasks

1. Add a host route to a test address on a lab machine that you own, then remove it. Record the commands and `ping` before and after.
2. Write a one-page field guide for `ip route` or `route print`. Name each column that you use.

---

## Longest prefix match

When two routes match the same destination, the device chooses the match with the longest prefix (the most specific mask). A `/32` wins over a `/24`. A `/24` wins over `/0`.

Example: destination `192.168.1.50`. Routes: `0.0.0.0/0` via G1, `192.168.1.0/24` via interface LAN, `192.168.1.50/32` via G2. The host route `/32` wins.

Longest prefix match is not the same as lowest metric. If two routes have the same prefix length, the device uses metric, preference, or equal-cost sharing. Details depend on the OS and on the routing protocol.

This rule lets you override a default path for one host or one subnet. A VPN can install a more specific prefix so that only that prefix uses the tunnel.

Do not add two different next hops for the same prefix without a plan for metric or ECMP. Do not assume that a smaller metric beats a longer prefix. The prefix length is first in the usual model.

IPv6 uses the same idea. `::/0` is the default. A `/128` is a host route.

### Questions

#### Theoretical questions

1. What does longest prefix match mean?
2. Which wins: `/24` or `/16` for a destination in both prefixes?
3. When does metric matter in this section?
4. Why can a VPN use a more specific prefix?
5. What is the IPv6 default prefix?

#### Easy practical tasks

1. Write five sentences that explain the `192.168.1.50` example.
2. Make a table: destination `10.1.2.3`, routes `/8`, `/16`, `/24`, `/0`. Mark the winner.
3. Draw three prefixes as nested boxes. Show a destination inside the smallest box.
4. In your table, find the most specific route and the least specific route.

#### Medium practical tasks

1. On paper, pick a winner among `10.0.0.0/8`, `10.1.0.0/16`, and `0.0.0.0/0` for `10.1.9.8`.
2. Write six sentences: two routes with the same prefix length and different metrics.
3. Explain why `192.168.1.0/24` and `192.168.1.0/25` can both exist. Give one destination that uses each.

#### Advanced practical tasks

1. In a lab that you own, install a more specific route and a default route. Show that `ping` to one address follows the specific route.
2. Write a one-page note: longest prefix match versus administrative distance in Cisco-like devices (high-level). Use public docs, not a dump of a manual.

---

## Host route vs default route

A host route is a prefix that names one address. IPv4 uses `/32`. IPv6 uses `/128`. The next hop can be a gateway or a local interface.

A default route matches every destination that no other route matched. IPv4 uses `0.0.0.0/0`. IPv6 uses `::/0`. The next hop is the default gateway.

A network route sits between them, for example `/24` or `/64`. Most connected LAN routes are network routes.

Use a host route when one address must use a special path. Use a default route when you cannot list the Internet. Do not use a default route to replace a missing connected route on your LAN. If the LAN prefix is missing, local traffic can take the default path and fail.

DHCP and RA often install the default route. A static administrator can install a host route for a management address.

On a host with two uplinks, two default routes need policy or different metrics. That design is more than this section. First understand one default route.

### Questions

#### Theoretical questions

1. What prefix length is an IPv4 host route?
2. What prefix is the IPv4 default route?
3. When do you use a host route?
4. Why is a connected `/24` not a default route?
5. What installs a default route on many home PCs?

#### Easy practical tasks

1. Write five sentences that compare a host route and a default route.
2. Find a `/32` or `/128` in your table if one exists. Write it, or write that you have none.
3. Make a table: host route, network route, default route. Add one example prefix each.
4. Draw a destination on the Internet and a destination on the LAN. Show which route type each uses on a typical PC.

#### Medium practical tasks

1. Predict the route type for: `127.0.0.1`, your LAN neighbor, `8.8.8.8`. Check the table.
2. Write six sentences: a missing connected route and a present default route. What goes wrong.
3. On Windows or Linux, read the metric of the default route. Write the value and the interface.

#### Advanced practical tasks

1. Add a host route to a public IP via a wrong next hop on a lab host, test, then delete the route. Use only a machine that you own.
2. Write a one-page design: management host route for a server BMC address on a separate path.

---

## ICMP: echo, dest unreachable, time exceeded

ICMP is a control protocol for IPv4. ICMPv6 is the IPv6 counterpart. This section uses IPv4 names. The ideas match ICMPv6 types that you already met (echo, neighbor, packet too big).

Echo request and echo reply are `ping`. They test reachability and round-trip time. A reply proves that a path worked for ICMP. A missing reply does not always prove that the host is down. A firewall can drop ICMP.

Destination unreachable tells the sender that the packet cannot be delivered. Codes include network unreachable, host unreachable, port unreachable, and fragmentation needed. A router or the destination host can send this message.

Time exceeded means the TTL (IPv4) or hop limit (IPv6) reached zero. A router drops the packet and can send time exceeded. Traceroute uses this message.

ICMP messages carry a prefix of the original packet. That copy helps the sender match the error to a flow.

Do not block all ICMP on a network if you want PMTU and traceroute to work. Do not trust ICMP alone as a security scanner. Do not flood ping at a target that you do not own.

### Questions

#### Theoretical questions

1. What do echo request and echo reply test?
2. Why can a missing ping reply lie?
3. What does destination unreachable mean?
4. What does time exceeded mean?
5. Why does an ICMP error include a piece of the original packet?

#### Easy practical tasks

1. Ping a LAN host and a public host. Write the times.
2. Write five sentences that explain TTL and time exceeded.
3. Make a table: ICMP name, who sends it, when. Add echo, dest unreachable, time exceeded.
4. In Wireshark, filter `icmp`. Capture one ping. Write type numbers if the tool shows them.

#### Medium practical tasks

1. Ping a closed TCP port on a host that you own (or use a local unused port). If you see dest unreachable / port unreachable, write it. If a firewall hides it, write that.
2. Set a low TTL on ping if your OS allows it (`ping -i` on Linux for interval is different; use `-t` on Windows for TTL, `-M`/`-t` as documented). Record a time-exceeded if you get one.
3. Write six sentences: fragmentation needed and why a firewall that drops ICMP breaks large transfers.

#### Advanced practical tasks

1. Capture ping and one traceroute hop. Map each ICMP type to a line in this section.
2. Write a one-page policy note: which ICMP types a border firewall should allow for a dual-stack office. High-level only.

---

## `traceroute` / `mtr`

Traceroute sends probes with increasing TTL. Hop 1 expires at the first router. Hop 2 expires at the second router. Each router can send time exceeded. The tool prints the source address of that message and the round-trip time.

Windows uses `tracert`. Many Unix systems use `traceroute`. `mtr` (or WinMTR) repeats the probes and shows loss and latency per hop.

A hop that shows stars (`*`) did not send a useful reply in time. The router can rate-limit ICMP. A firewall can drop the probe. The path can still work for TCP.

The address that you see is the address of the inbound interface that the router uses to send ICMP, in common cases. It is not always the address that your packet used as a next hop. Load balancing can show different hops on different runs.

Do not treat traceroute as a legal map of a foreign network. Do not run aggressive `mtr` against a host that you do not own.

UDP, ICMP, and TCP traceroute variants exist. Firewalls treat them differently. If UDP traceroute fails, a TCP traceroute to port 443 can still show hops.

### Questions

#### Theoretical questions

1. How does traceroute use TTL?
2. What ICMP message do hops send in the common case?
3. What does `*` mean in the output?
4. Why can traceroute fail when a website works?
5. What extra view does `mtr` add?

#### Easy practical tasks

1. Run `tracert` or `traceroute` to `example.com`. Copy the first three hops.
2. Write five sentences that explain increasing TTL.
3. Mark which hops look like your LAN or your ISP. Give one reason.
4. Make a table: hop number, address, time. Add three rows.

#### Medium practical tasks

1. Compare `ping` to the last hop and `ping` to the destination name. Write both results.
2. Run traceroute two times. Write whether any hop address changed.
3. Write six sentences: why a middle hop can show 100 percent loss in `mtr` while the destination still answers.

#### Advanced practical tasks

1. Use `mtr` or WinMTR to a public host for one minute. Write min, avg, max for the destination row. Do not flood.
2. Write a one-page lab: traceroute with ICMP versus TCP if your tool supports it. Compare hop lists.

---

## Static route vs dynamic routing (survey: OSPF, BGP)

A static route is a row that an administrator writes. The row stays until a person or a script changes it. Static routes fit a small network or a special exception.

A dynamic routing protocol updates the table from messages. Neighbors share prefixes. The protocol computes a path. When a link dies, the table can change without a human.

OSPF is a common interior protocol. Routers in one organization flood link state and compute shortest paths. OSPF uses areas in larger designs. This topic only needs: OSPF is for the inside of a network.

BGP is the protocol of the Internet between autonomous systems. BGP is a path-vector protocol. ISPs and large networks use BGP to choose policy, not only shortest hop count. A later topic covers BGP. This topic only needs: BGP connects organizations.

Do not run BGP on a home LAN for learning unless you use a lab. Do not paste a full Internet table into a small router.

Choose static when the network is small and stable. Choose OSPF (or similar) when many internal routers must converge. Choose BGP when you connect to ISPs or to other organizations with policy.

A default route can be static or learned. A home DHCP client learns a static-like default from the box. The box may use a default toward the ISP.

### Questions

#### Theoretical questions

1. Who writes a static route?
2. What does a dynamic protocol do when a link dies?
3. Where do you use OSPF in this survey?
4. Where do you use BGP in this survey?
5. Why is hop count not the only BGP goal?

#### Easy practical tasks

1. Write five sentences that compare static and dynamic routing.
2. Make a table: static, OSPF, BGP. Add one row for "typical place" and one row for "who updates the table".
3. List two cases that fit a static route on a home or lab network.
4. Draw three office routers. Label an OSPF domain and one static default to an ISP.

#### Medium practical tasks

1. On your PC, write which of your routes look static or DHCP-installed. Give a reason.
2. Find a public OSPF one-page overview and a BGP one-page overview. Write three STE sentences each. No long quotes.
3. Write six sentences: a static route to a dead next hop, versus OSPF that can drop that next hop.

#### Advanced practical tasks

1. In a simulator or a lab that you own, add one static route between two subnets. Document the commands. Do not use a production ISP.
2. Write a one-page decision sheet: static versus OSPF versus BGP for a company with two offices and one ISP.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a packet lookup from destination IP to next-hop MAC. Use table, longest prefix, and ICMP if the lookup fails.
2. How do host routes, network routes, and the default route share one table without conflict?
3. Why do ping and traceroute both use ICMP but answer different questions?
4. When is a static default enough, and when do you need OSPF or BGP?
5. A hop shows `*` in traceroute but the website loads. Which facts explain that result?

#### Easy practical tasks

1. Write a one-page cheat sheet: table, LPM, host route, default, ICMP types, traceroute, static, OSPF, BGP.
2. Save `ip route` / `route print` and one `tracert` output in a file. Highlight the default route and hop 1.
3. Draw a router with three interfaces and four table rows that you invent. Show LPM for one destination.
4. Ping the gateway and traceroute to a public host. Write how hop 1 relates to the gateway.

#### Medium practical tasks

1. Build a paper table with five routes. Give three destinations and the winning route for each.
2. Write a lab plan that uses ping, traceroute, and a routing table print to find a wrong default gateway.
3. Compare ICMPv4 time exceeded and ICMPv6 time exceeded in a two-column table. Use a capture if you can.

#### Advanced practical tasks

1. Write a glossary of 12 routing terms. Each entry: term, one sentence, one command or ICMP type.
2. Design a three-router lab (draw or simulate): one static route, one default, one OSPF note. Predict traceroute from a host.
