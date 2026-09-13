# 3. IPv4 and IPv6 Addresses

## Description

IPv4 uses a 32-bit address. IPv6 uses a 128-bit address. This topic shows IPv4 dotted decimal, mask, and CIDR. You learn private ranges, loopback, and the default gateway. You learn why IPv6 exists, address form, and compression. You learn link-local `fe80::`, SLAAC, and dual stack.

Complete this topic after links and Ethernet. Complete this topic before you study routing, ICMP, and NAT.

Use one term for each concept. An IPv4 address is a 32-bit identifier for an interface. An IPv6 address is a 128-bit identifier. A subnet is a group of addresses that share a prefix. A mask or a prefix length marks the network part. Dual stack means IPv4 and IPv6 run together. Do not mix a private IPv4 address with a public IPv4 address. Do not mix link-local IPv6 with a global IPv6 address.

Work on paper and on your machine. Wrong masks cause "it pings the LAN but not the Internet" faults.

---

## IPv4 dotted decimal, mask, CIDR

An IPv4 address is 32 bits. Dotted decimal writes four octets. Each octet is 8 bits. Each octet is a number from 0 to 255. Example: `192.168.1.10`.

You must convert between binary and decimal for one octet. Example: `11000000` is 192. The conversion uses powers of two: 128, 64, 32, 16, 8, 4, 2, 1. The host and the router use the 32-bit value. Dots are only for humans.

An IPv4 address has a network part and a host part. The network part is the same for all interfaces on one subnet. The host part is unique on that subnet. The subnet mask or the prefix length sets the split.

CIDR writes a prefix length after a slash. `192.168.1.0/24` means 24 network bits and 8 host bits. A `/16` uses 16 network bits. The old class A, B, and C names are history. CIDR is the current form.

Two interfaces can talk on the LAN without a router when they share the same network part and the mask is correct. If the network parts differ, a router must forward the packet.

The all-zero host part is the network identifier in common practice (`192.168.1.0/24`). The all-one host part is the directed broadcast (`192.168.1.255/24`). Do not assign those two addresses to a host on a typical Ethernet subnet.

Host bits give the count of addresses: 2 to the power of the host-bit count. Usable host addresses are often that count minus 2, when you reserve network and broadcast.

Leading zeros in an octet can confuse some tools. Some parsers read a leading zero as octal. Write octets without leading zeros.

Do not mix two different masks on one LAN without a plan. Hosts will disagree about who is local.

### Questions

#### Theoretical questions

1. How many bits does an IPv4 address use?
2. What is the range of one octet in decimal?
3. What does `/24` mean in CIDR?
4. Who forwards a packet when the network parts differ?
5. What is the all-one host part on a typical subnet?

#### Easy practical tasks

1. Convert `10` and `255` to 8-bit binary.
2. Convert `11000000` and `00000001` to decimal.
3. Write `192.168.0.1` as 32 bits grouped in four octets.
4. Make a table of powers of two from 1 to 128.

#### Medium practical tasks

1. Convert these addresses to binary: `10.0.0.1`, `172.16.5.4`, `8.8.8.8`.
2. On paper, write the network address and broadcast of `192.168.10.50/24`.
3. Write six sentences that explain why `256` cannot be an octet.

#### Advanced practical tasks

1. Subnet a `/24` into four `/26` networks on paper. Write each prefix, mask, and host range.
2. Write a small program or a spreadsheet that converts one octet both ways. Test 0, 1, 127, 128, and 255.

---

## Private ranges, loopback, gateway

RFC 1918 defines private IPv4 ranges. They are not unique on the public Internet. Routers on the public Internet must not route them as global destinations.

The private ranges:

- `10.0.0.0/8`
- `172.16.0.0/12`
- `192.168.0.0/16`

A home LAN often uses `192.168.0.0/24` or `192.168.1.0/24`. An enterprise often uses `10.0.0.0/8` pieces. NAT (topic 4) lets many private hosts share one public IPv4 address.

Loopback stays on the same host. `127.0.0.0/8` is the IPv4 loopback range. `127.0.0.1` is the usual address. A packet to loopback does not leave the machine.

The default gateway is the next hop for destinations that are not on a connected subnet. The host puts that router IP in the routing table as `0.0.0.0/0`. The gateway must sit on a connected LAN of the host. You cannot use a gateway that is not neighbor-reachable.

A "what is my IP" page shows a public address. That address is often the NAT address, not your LAN address. Both values are real. They name different points.

Other special IPv4 uses exist. `169.254.0.0/16` is link-local (APIPA) when DHCP fails. `224.0.0.0/4` is multicast. You do not need every special range here. You must know private, loopback, and gateway.

Do not assign a private address to a public server that must be reachable from the Internet without NAT. Do not ping `8.8.8.8` and think a reply proves your LAN address is public.

### Questions

#### Theoretical questions

1. Name the three RFC 1918 private ranges.
2. Why are private addresses not unique on the Internet?
3. What is `127.0.0.1` for?
4. What prefix is the IPv4 default route?
5. Why must the gateway be on a connected subnet?

#### Easy practical tasks

1. Write your LAN IPv4 and your public IPv4 (trusted lookup). Label which is private.
2. Write five sentences that explain private range, loopback, and gateway.
3. Make a table: address, private or public or loopback. Add `10.1.1.1`, `8.8.8.8`, `127.0.0.1`.
4. Draw a LAN host, a gateway, and the Internet.

#### Medium practical tasks

1. Print your routing table. Copy the default route. Write the gateway IP.
2. Ping `127.0.0.1` and ping your gateway. Write both results.
3. Write six sentences: why two homes can both use `192.168.1.10`.

#### Advanced practical tasks

1. Write a one-page home map: private range, mask, gateway, public IPv4, loopback.
2. Decide on paper whether each destination is local or needs the gateway: `192.168.1.20` and `1.1.1.1` for a host `192.168.1.10/24` with gateway `192.168.1.1`.

---

## Why IPv6 exists; address form and compression

IPv4 has about 4.3 billion addresses. NAT and private ranges delayed the end of free public IPv4. They did not create more public addresses. Many networks still need unique addresses for hosts, mobiles, and servers.

IPv6 gives a much larger address space. 128 bits can number a huge number of interfaces. Operators can give a home a `/56` or a `/64` and still have room.

IPv6 is not IPv4 with longer dotted decimal. The address form is hex. Broadcast is gone. Neighbor discovery replaces ARP. ICMP is required. Some middleboxes treat IPv6 poorly. You still need IPv4 on many paths. That is dual stack.

An IPv6 address is eight groups of 16 bits. Each group is up to four hex digits. Groups use colons. Example: `2001:0db8:0000:0000:0000:0000:0000:0001`.

Compression rules:

- you may omit leading zeros in a group: `2001:db8:0:0:0:0:0:1`
- you may replace one run of groups that are all zero with `::`
- you may use `::` only one time in an address

The compressed form of the example is `2001:db8::1`. Documentation uses `2001:db8::/32`. Do not use documentation addresses as real global addresses.

A LAN often uses `/64`. A single address is `/128`. Text in URLs puts IPv6 in brackets: `http://[2001:db8::1]:80/`.

Do not wait for the day IPv4 dies to learn IPv6. Your LAN already has IPv6 link-local addresses. Do not write two `::` markers. Do not treat NAT as a design that removes the need for IPv6.

### Questions

#### Theoretical questions

1. Why is the IPv4 address space a problem?
2. How many bits does an IPv6 address use?
3. What does `::` mean?
4. Why may you use `::` only one time?
5. How do you write an IPv6 address with a port in a URL?

#### Easy practical tasks

1. Write five sentences that explain IPv4 scarcity and IPv6 size.
2. Expand `2001:db8::1` to eight groups.
3. Compress `fe80:0000:0000:0000:0202:b3ff:fe1e:8329`.
4. Make a table: IPv4, IPv6. Add rows for bit length, text form, and broadcast.

#### Medium practical tasks

1. Decide which of these are valid: `2001:db8::1::2`, `2001:db8:0:0:0:0:0:1`, `2001:db8::`. Explain.
2. Write the first 64 bits of `2001:db8:1:2:3:4:5:6` as a `/64` prefix.
3. Find one IPv6 address on your machine. Write the full form and the compressed form.

#### Advanced practical tasks

1. Write a small program that compresses and expands IPv6 text, or use a trusted library. Test `::1`, `2001:db8::`, and a link-local address.
2. Write a one-page note for a manager: why the team must test IPv6 even if the product uses IPv4 today.

---

## Link-local `fe80::`, SLAAC, dual stack

A link-local IPv6 address starts with `fe80::/10`. In practice you see `fe80:` and then zeros, then an interface identifier. The address is valid only on one link. Routers do not forward link-local packets to another link.

Every IPv6-capable interface has a link-local address when IPv6 is on. Neighbor discovery uses link-local addresses. The same link-local prefix appears on every link. You must add a zone or interface when you ping: `ping fe80::1%eth0` or `ping fe80::1%12` on Windows. The `%` part is the zone index. It is not part of the 128-bit address on the wire.

SLAAC is Stateless Address Autoconfiguration. A router sends Router Advertisements. A host builds a global or unique-local address from a prefix plus an interface identifier. Privacy addresses can change the interface identifier. DHCPv6 can still give options or addresses. SLAAC does not need DHCPv6 for a basic address.

Dual stack means IPv4 and IPv6 run at the same time. A name can have an A record and an AAAA record. A client can try both. Happy Eyeballs can pick a family that works. A fault on one family does not always show on the other.

Link-local is not a private global prefix. Unique local addresses (`fc00::/7`, commonly `fd00::/8`) are a different idea. Global addresses often start in `2000::/3` for current global unicast.

Do not filter out all `fe80` traffic on a LAN if you want IPv6 to work. Do not use a link-local address as a public server address. Do not disable IPv6 on a test and then forget that production still has it.

### Questions

#### Theoretical questions

1. What prefix marks a typical link-local IPv6 address?
2. Why do you add `%` and an interface name when you ping `fe80::`?
3. What does SLAAC use from the router?
4. What does dual stack mean?
5. Do routers forward link-local packets to another link?

#### Easy practical tasks

1. List IPv6 addresses on your machine. Mark link-local versus other types if you can.
2. Write five sentences that explain `fe80::`, SLAAC, and dual stack.
3. Make a table: link-local, unique local, global. Add one prefix hint each.
4. Draw a host with IPv4, `fe80::`, and one global IPv6 if present.

#### Medium practical tasks

1. Ping `::1`. Then try a link-local ping with a zone if you have a peer on the LAN.
2. Write six sentences: a site with IPv4 NAT and IPv6 SLAAC on the same LAN.
3. Resolve `example.com` A and AAAA. Write whether dual stack is possible for that name.

#### Advanced practical tasks

1. Read a short public SLAAC note. Write the steps: RA, prefix, interface identifier, Duplicate Address Detection (name only).
2. Write a one-page dual-stack test plan: ping, DNS A/AAAA, and one HTTPS fetch on each family if the path allows it.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A host can ping a LAN neighbor but not the Internet. Which IPv4 fields do you check first?
2. How do CIDR prefix length and a gateway work together on one host?
3. Why can IPv6 link-local exist when no global IPv6 prefix is present?
4. What breaks if you treat `192.168.1.1` as a unique Internet identity?
5. How does dual stack change a "the site is down" report?

#### Easy practical tasks

1. Write a cheat sheet: dotted decimal, /24, RFC 1918, loopback, gateway, IPv6 groups, `::`, `fe80::`, SLAAC, dual stack.
2. Run `ipconfig` or `ip addr`. Write IPv4, mask or prefix, gateway, and any IPv6.
3. Compress and expand one IPv6 address from your machine or `2001:db8:0:0:0:0:0:1`.
4. Draw IPv4 private LAN plus IPv6 link-local on the same NIC.

#### Medium practical tasks

1. On paper, pick a host address in `10.0.0.0/22`. Write mask, network, broadcast, and a valid gateway.
2. Compare `ping -4` and `ping -6` to a name that has both families if your tool supports it. Write both outcomes.
3. Write a short lab: two hosts, same `/24`, different `/25`. Explain who can talk without a router.

#### Advanced practical tasks

1. Subnet a `/24` into `/26`s as the practice list asks. Then write one IPv6 `/64` plan for the same LAN idea.
2. Write a one-page addressing standard for a home lab: IPv4 private range, DHCP or static, IPv6 SLAAC, and when you use loopback.
