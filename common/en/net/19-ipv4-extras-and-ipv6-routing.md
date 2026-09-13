# 19. IPv4 Extras and IPv6 Routing

## Description

This topic adds IPv4 and IPv6 behaviors that the first addressing and routing topics only named. It shows ICMP redirect, path MTU discovery, IPv6 router advertisements and duplicate address detection, unique local addresses, and transition tools NAT64 and 464XLAT at an awareness level.

Use one term for each concept. A redirect tells a host a better next hop on the same link. Path MTU is the smallest MTU on the path. Complete this topic after IPv4, IPv6, routing, and the MTU preview.

Practice on a LAN that you own. Capture ICMPv6 RAs if the LAN sends them. Do not flood ping. Do not change production PMTU sysctls without a backout.

---

## ICMP redirect

A router can send an ICMP redirect (IPv4) or an ICMPv6 redirect. The message tells a host: for this destination, use a different next hop on the same link. The first hop router is not the best first hop.

Example: two routers on one LAN. The host uses router A as default. The packet for prefix P should go to router B. A sends a redirect. The host can add a host route or a prefix route toward B.

Many hosts ignore redirects or limit them. Attackers on a LAN can send false redirects. Modern OS defaults are cautious. Do not depend on redirects as a clean design. Use a correct default or a more specific route.

Do not enable "accept redirects" on a host that sits on an untrusted LAN unless you know the risk. Do not use redirects as the only way to steer a data center. Use routing protocols or static routes that you document.

A capture shows type 5 on IPv4 ICMP in the common case. Wireshark names the message. You already met echo and time exceeded in topic 6.

### Questions

#### Theoretical questions

1. What does a redirect tell a host?
2. Why would two routers on one LAN cause a redirect?
3. Why do some hosts ignore redirects?
4. Why is a redirect a poor only design for a data center?
5. What ICMP type is the common IPv4 redirect?

#### Easy practical tasks

1. Write five sentences that explain a better next hop on the same link.
2. Make a table: default route, redirect, static host route. Add one row for who writes the route.
3. Draw: host, router A, router B, destination cloud.
4. Write four sentences: untrusted LAN and false redirects.

#### Medium practical tasks

1. Find whether your OS documents "accept ICMP redirects." Write the setting name if you find it. Do not change a work laptop.
2. Write six sentences: a correct static route versus a redirect.
3. In a lab with two routers that you own, create a redirect if you can. Write whether the host route appeared. Skip if you have one router.

#### Advanced practical tasks

1. Write a one-page note: when a redirect is a useful hint and when it is a threat.
2. Capture on a LAN that you own for a few minutes. Write whether you see redirects. Home lab only.

---

## Path MTU discovery

Path MTU discovery (PMTUD) finds the smallest MTU on the path so that the sender can avoid a size that a hop cannot carry.

IPv4 PMTUD: the sender sets the Don't Fragment (DF) bit. A hop that needs to fragment drops the packet and should send ICMP destination unreachable, fragmentation needed, with the next-hop MTU. The sender shrinks the size.

IPv6: routers do not fragment. The hop sends ICMPv6 Packet Too Big. The source shrinks. Topic 3 previewed this.

A firewall that drops those ICMP messages causes black-hole PMTU: small packets work, large TCP segments hang. Topic 6 already warned you. The fix is to allow the needed ICMP, or to clamp MSS, or to set a safe MTU on a tunnel. MSS clamp is a common middlebox trick for TCP.

Do not turn off DF on IPv4 as a first production fix. Fragmentation has its own faults. Do not set a 1280 IPv6 MTU on every LAN without a reason. 1280 is the IPv6 minimum that a path must support.

Tunnels (VPN, GRE) add headers. The effective MTU drops. Users then see PMTU faults after they "only added a VPN."

Ping with a large size and DF (IPv4) is a lab test. Do not flood. Do not test on a network that you do not own.

### Questions

#### Theoretical questions

1. What does PMTUD try to find?
2. What IPv4 ICMP message should a too-small hop send?
3. What ICMPv6 message does a too-small hop send?
4. What is a PMTU black hole?
5. Why do tunnels shrink the effective MTU?

#### Easy practical tasks

1. Write five sentences that explain DF and "fragmentation needed."
2. Find the MTU on one interface. Write the number.
3. Make a table: IPv4 PMTUD, IPv6 PMTUD. Add who fragments and who sends the error.
4. Draw a path: 1500, then a 1400 tunnel, then 1500.

#### Medium practical tasks

1. Ping with a payload near your MTU if your OS allows size and DF. Write the command and the result. One or two probes.
2. Write six sentences: VPN on, large HTTPS hang, small ping works.
3. Explain TCP MSS clamp in four STE sentences.

#### Advanced practical tasks

1. Write a one-page lab: find path MTU to a public host with rising ping sizes. Stop at the first fail. Do not flood.
2. Write a firewall note: which ICMP types a border must pass for PMTUD on dual stack. High-level.

---

## IPv6 RA / DAD

A router advertisement (RA) is an ICMPv6 message. The router sends it periodically and in reply to a router solicitation (RS). The RA can carry prefixes, flags (M, O), MTU, and recursive DNS options. Topic 5 introduced SLAAC. This section adds the control-plane view.

Hosts use the RA to:

- pick a default router (the source of the RA, a link-local address)
- form addresses with SLAAC when the prefix allows it
- learn whether to start DHCPv6
- learn DNS servers when RDNSS is present

Duplicate address detection (DAD) runs before a host uses an IPv6 address. The host sends a neighbor solicitation for its own tentative address. If another host answers, the address is a duplicate. The host must not use it.

Do not disable RAs on a LAN and then expect global IPv6. Do not run two routers that advertise conflicting prefixes without a design. Do not skip DAD in an implementation that you write.

A rogue RA on an untrusted LAN can steal a default route. RA guard on switches is a hardening idea. Topic 21 returns to untrusted networks.

Windows, macOS, and Linux all process RAs by default. A capture filter `icmpv6` shows them.

### Questions

#### Theoretical questions

1. What ICMPv6 message carries prefixes for SLAAC?
2. What does a host send to ask for an RA?
3. What is DAD for?
4. Why is the default router often a link-local address?
5. What can a rogue RA do?

#### Easy practical tasks

1. Write five sentences that list RA jobs.
2. Make a table: RS, RA, NS for DAD. Add one row for who sends.
3. List IPv6 addresses on your host. Write which one looks link-local.
4. Draw: router, RA, host, tentative address, DAD NS.

#### Medium practical tasks

1. Capture an RA on a LAN that you own. Write the source address and any prefix the tool shows.
2. Write six sentences: M and O flags as "use DHCPv6" hints (high-level).
3. Compare a LAN with RA and a LAN with no RA. What IPv6 addresses remain?

#### Advanced practical tasks

1. Write a one-page office IPv6 LAN: RA prefix, DAD, DNS from RDNSS or DHCPv6.
2. On a lab router that you own, change an RA prefix (or use a VM). Document host addresses before and after. No production LAN.

---

## Unique local addresses

Unique local addresses (ULA) are IPv6 addresses in `fc00::/7`. The common practice is `fd00::/8` with a random 40-bit global ID, then a 16-bit subnet, then a 64-bit interface ID. ULA is for local communications. It is not the same as `fe80::/10` link-local. A ULA can route inside one organization across many links. A link-local cannot.

ULA is the IPv6 cousin of RFC 1918 in spirit, not in bits. You must not think that ULA is unroutable on the Internet by magic. Operators should not announce ULA to the public Internet. A default IPv6 route can still try. Policy must stop that.

A host can have ULA and global addresses at the same time. Happy Eyeballs and default address selection choose a source and a destination. A ULA-only destination will not work from the public Internet.

Do not pick `fd00::1` for every lab. Randomize the global ID so that two merged sites do not collide. Do not use ULA as the only address on a server that must be public.

Documentation prefix `2001:db8::/32` is not ULA. Do not mix them.

### Questions

#### Theoretical questions

1. What prefix block is ULA?
2. How is ULA different from link-local?
3. Should operators announce ULA to the public Internet?
4. Can a host have ULA and a global address together?
5. Why randomize the ULA global ID?

#### Easy practical tasks

1. Write five sentences that compare ULA, link-local, and global.
2. Make a table: `10.0.0.0/8`, `fe80::/10`, `fd00::/8`, `2000::/3` idea. Add one row for scope.
3. Write a sample ULA `/48` on paper (`fd` plus hex). Label subnet and host parts.
4. Draw: two LANs, one ULA `/48`, one router.

#### Medium practical tasks

1. Check your host for `fd` or `fc` addresses. Write what you found.
2. Write six sentences: ULA-only app plus a user on the Internet.
3. Compare ULA and IPv4 RFC 1918 in a short paragraph. Use STE.

#### Advanced practical tasks

1. Write a one-page site plan: ULA for internal, global for public, NAT64 later if needed.
2. Generate a random ULA `/48` with a method you document (any fair random hex). Show the bits.

---

## Transition: NAT64, 464XLAT (awareness)

IPv4 and IPv6 hosts must still talk during the long dual-stack era. Transition tools translate or encapsulate. This section is awareness.

NAT64 lets an IPv6-only client reach an IPv4-only server. A translator has an IPv6 prefix (often `64:ff9b::/96` for well-known). The IPv4 address embeds in the IPv6 destination. DNS64 synthesizes AAAA records from A records so that the client opens IPv6 toward the translator.

464XLAT uses NAT64 plus a client CLAT. The CLAT gives IPv4 to local apps that cannot speak IPv6. Mobile networks use this so that an IPv4-only app still works on an IPv6-only radio.

Other names that you may see: 6to4, Teredo, DS-Lite, MAP. Many are legacy or ISP-specific. Do not enable random tunnels on a work PC.

Do not treat NAT64 as a reason to skip IPv6 on servers that you can dual-stack. Do not forget that translation breaks some protocols that carry addresses in payloads.

SIIT is stateless IP/ICMP translation. NAT64 is often stateful for ports. You only need: IPv6-only plus IPv4 world needs a translator or dual stack.

### Questions

#### Theoretical questions

1. What problem does NAT64 solve?
2. What does DNS64 synthesize?
3. What extra job does CLAT do in 464XLAT?
4. Why can translation break some applications?
5. Why is dual stack still the simple server default when you can do it?

#### Easy practical tasks

1. Write five sentences that explain IPv6-only client to IPv4-only server.
2. Make a table: dual stack, NAT64, 464XLAT. Add one row for who needs IPv4 on the host.
3. Write the well-known NAT64 prefix if this section named it.
4. Draw: IPv6 client, NAT64, IPv4 server.

#### Medium practical tasks

1. Read a public NAT64 or 464XLAT explainer from an operator. Write four STE sentences.
2. Write six sentences: an IPv4-only mobile app on an IPv6-only carrier.
3. Check whether your phone or laptop UI mentions IPv6-only or 464XLAT. Write what you see.

#### Advanced practical tasks

1. Write a one-page awareness note: DNS64 plus NAT64 plus a broken literal IPv4 in an app.
2. Design a lab on paper: IPv6-only VLAN, NAT64 gateway, IPv4-only website. List tests (`ping`, `curl`). Do not build it on a network that you do not own.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a dual-stack host that gets an RA, runs DAD, and then sends a large TCP segment through a VPN. Use PMTUD and one extra address type (ULA or global).
2. How do ICMP redirect and a correct routing table overlap, and when should you ignore redirects?
3. Why do NAT64 and ULA both appear in "IPv6 routing extras" even though one is translation?
4. When is a PMTU black hole the first suspect, and when is a missing default router the first suspect?
5. A teammate says "IPv6 is SLAAC only, no routing extras." Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: redirect, PMTUD, RA, DAD, ULA, NAT64, 464XLAT.
2. Save `ipconfig` / `ip addr` IPv6 lines and one interface MTU. Label link-local, ULA if any, global if any.
3. Draw one picture: two LAN routers, RA, host, DF packet, Packet Too Big.
4. List ICMPv6 names from this topic (RA, RS, PTB, redirect, DAD NS).

#### Medium practical tasks

1. From your machine, identify: MTU, IPv6 default (if any), whether an RA is visible in a short capture, and one transition feature you do or do not use.
2. Write a lab plan: safe PMTU ping test plus RA capture. Include "do not flood."
3. Explain in ten steps how a beginner checks IPv6 default and MTU before they disable IPv6.

#### Advanced practical tasks

1. Build a glossary of 14 terms from this topic. Each entry: term, one sentence, one command or ICMP name.
2. Write a short design: IPv6-only office VLAN with DNS64/NAT64 for IPv4-only SaaS, ULA for internal, RA for hosts.
