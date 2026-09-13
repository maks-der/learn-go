# 4. IPv4 Addressing

## Description

IPv4 uses a 32-bit address. Humans write the address as four decimal octets. This topic shows binary, the network part, the host part, subnet masks, CIDR, private ranges, loopback, broadcast, multicast, and the default gateway.

Use one term for each concept. An IPv4 address is a 32-bit identifier for an interface. A subnet is a group of addresses that share a prefix. A mask or a prefix length marks the network part. Complete this topic before you study IPv6 and routing.

Work on paper and on your machine. Wrong masks cause "it pings the LAN but not the Internet" faults.

---

## Binary and dotted decimal

An IPv4 address is 32 bits. Dotted decimal writes four octets. Each octet is 8 bits. Each octet is a number from 0 to 255. Example: `192.168.1.10` is four octets.

You must convert between binary and decimal for one octet. Example: `11000000` is 192. Example: `10101000` is 168. The conversion uses powers of two: 128, 64, 32, 16, 8, 4, 2, 1.

The full address in binary is 32 bits with no dots. Dots are only for humans. The host and the router use the 32-bit value.

Leading zeros in an octet do not change the value. `192.168.1.010` can confuse some tools because some parsers read a leading zero as octal. Write octets without leading zeros.

Do not write an IPv4 address with more than four parts. Do not use hex for IPv4 in normal lab notes unless a tool shows hex.

A packet has a source IPv4 address and a destination IPv4 address. Each address is 32 bits in the header.

### Questions

#### Theoretical questions

1. How many bits does an IPv4 address use?
2. What is the range of one octet in decimal?
3. Why do humans use dotted decimal?
4. What is the binary weight of the left bit in an octet?
5. Why can a leading zero in an octet cause a tool error?

#### Easy practical tasks

1. Convert `10` and `255` to 8-bit binary.
2. Convert `11000000` and `00000001` to decimal.
3. Write `192.168.0.1` as 32 bits grouped in four octets.
4. Make a table of powers of two from 1 to 128.

#### Medium practical tasks

1. Convert these addresses to binary: `10.0.0.1`, `172.16.5.4`, `8.8.8.8`.
2. A tool shows `c0 a8 01 01` as hex octets. Write the dotted decimal address.
3. Write six sentences that explain why `256` cannot be an octet.

#### Advanced practical tasks

1. Write a small program or a spreadsheet that converts one octet both ways. Test 0, 1, 127, 128, and 255.
2. Read 8 bytes of an IPv4 header in a capture hex view. Mark the source and destination address bytes. Write them in dotted decimal.

---

## Network vs host part

An IPv4 address has a network part and a host part. The network part is the same for all interfaces on one subnet. The host part is unique on that subnet.

The split is not fixed at the old class A, B, or C boundaries. The subnet mask or the prefix length sets the split. A `/24` uses 24 bits for the network and 8 bits for the host. A `/16` uses 16 bits for the network.

Two interfaces can talk on the LAN without a router when they share the same network part and the mask is correct. If the network parts differ, a router must forward the packet.

The all-zero host part is the network identifier in common practice (for example `192.168.1.0/24`). The all-one host part is the directed broadcast address (for example `192.168.1.255/24`). Do not assign those two addresses to a host on a typical Ethernet subnet.

Host bits give the count of addresses in the subnet: 2 to the power of the host-bit count. Usable host addresses are often that count minus 2, when you reserve network and broadcast.

Do not mix two different masks on one LAN without a plan. Hosts will disagree about who is local.

### Questions

#### Theoretical questions

1. What is the network part of an address?
2. What is the host part of an address?
3. Who forwards a packet when the network parts differ?
4. What is the all-one host part on a typical subnet?
5. How do you compute the number of addresses in a subnet?

#### Easy practical tasks

1. For `192.168.1.10/24`, write the network part and the host part in decimal and in bits.
2. For `10.0.5.9/16`, write how many bits are network bits.
3. Draw a box for the network bits and a box for the host bits for `/8`, `/16`, and `/24`.
4. Write five sentences that explain why two hosts with `192.168.1.10` and `192.168.2.10` need a router if both use `/24`.

#### Medium practical tasks

1. List all usable host addresses for `192.168.10.0/30`. Show the network and broadcast addresses.
2. Two hosts: `172.16.4.1/22` and `172.16.5.2/22`. Decide if they are on the same subnet. Show the network identity.
3. Write six sentences about a mask mismatch: host A uses `/24`, host B uses `/16`, same first three octets.

#### Advanced practical tasks

1. On paper, split `192.168.50.0/24` into two equal subnets. Write each network, mask, and usable range.
2. Write a one-page note: old classful A/B/C versus a mask that you choose. Why this handbook uses the mask, not the class.

---

## Subnet mask and CIDR (`/24`)

A subnet mask is 32 bits. Ones mark the network part. Zeros mark the host part. Example: `255.255.255.0` is 24 ones and 8 zeros. That mask is `/24` in CIDR.

CIDR is Classless Inter-Domain Routing. The prefix length is the number of network bits. You write `192.168.1.0/24`. You can also write the mask `255.255.255.0`.

Common masks:

- `/32` one address (a host route)
- `/31` two addresses (point-to-point, special rules)
- `/30` four addresses, two usable hosts on a typical LAN
- `/24` 256 addresses, 254 usable on a typical LAN
- `/16` 65536 addresses
- `/8` about 16 million addresses

A host uses the mask to decide: is the destination local, or must I send the packet to a gateway? The host does a bitwise AND of the destination with the mask. The host compares the result with its own network identity.

Do not write `/24` on an address without a network identity when you mean a single host. `192.168.1.10/24` on an interface means the host address is `.10` and the subnet is `192.168.1.0/24`.

Wildcard bits are the inverse of the mask. Some ACL tools use wildcards. This topic uses prefix length and dotted masks.

### Questions

#### Theoretical questions

1. What do the one-bits in a subnet mask mean?
2. What does `/24` mean?
3. What is CIDR?
4. How does a host decide that a destination is local?
5. What does `192.168.1.10/24` on an interface mean?

#### Easy practical tasks

1. Write the dotted mask for `/8`, `/16`, `/24`, and `/32`.
2. Write the prefix length for `255.255.255.128` and `255.255.0.0`.
3. Make a table: prefix, total addresses, usable hosts on a typical LAN. Add `/24`, `/25`, `/30`.
4. Write five sentences that explain the AND of `192.168.1.50` with `255.255.255.0`.

#### Medium practical tasks

1. Convert `/20` to a dotted mask. Write the number of host bits.
2. Host `10.4.9.10` mask `255.255.252.0`. Write the network identity and the broadcast address.
3. On paper, subnet a `/24` into four `/26` networks. Write each network address.

#### Advanced practical tasks

1. Write a small program that prints network and broadcast for an address and a prefix length. Test three pairs.
2. Find the prefix on your interface. Compute the network identity by hand. Confirm with the OS output.

---

## Private ranges (`10/8`, `172.16/12`, `192.168/16`)

RFC 1918 sets private IPv4 ranges. Those addresses are not unique on the public Internet. Routers on the public Internet must not route them as global destinations.

The ranges:

- `10.0.0.0/8` (`10.0.0.0` to `10.255.255.255`)
- `172.16.0.0/12` (`172.16.0.0` to `172.31.255.255`)
- `192.168.0.0/16` (`192.168.0.0` to `192.168.255.255`)

A home LAN often uses `192.168.0.0/24` or `192.168.1.0/24`. A lab often uses `10.0.0.0/8` pieces. The choice is local.

NAT (next topics) translates a private source to a public address at the edge. Private addresses can repeat in different homes. That is why two homes can both use `192.168.1.1`.

Other special ranges exist (documentation, link-local `169.254.0.0/16`, CGNAT `100.64.0.0/10`). This section focuses on RFC 1918. You will see `169.254` when IPv4 autoconfig fails to get a DHCP lease.

Do not announce private prefixes to the public Internet. Do not assume that a `10.` address is always "internal" if you are on a large campus that also uses public space.

### Questions

#### Theoretical questions

1. Why do private IPv4 ranges exist?
2. What is the prefix of the `10` private range?
3. What is the full span of `172.16/12`?
4. Why can two homes use `192.168.1.1`?
5. What does `169.254.0.0/16` often mean on a host?

#### Easy practical tasks

1. Write the three RFC 1918 prefixes and one example address in each.
2. Look at your LAN IPv4 address. Write whether it is in a private range.
3. Make a table: range, prefix, typical use. Add the three private ranges.
4. Write five sentences that explain why a public website is not at `192.168.1.1`.

#### Medium practical tasks

1. Decide which of these are private: `10.1.2.3`, `172.15.0.1`, `172.16.0.1`, `192.168.100.50`, `8.8.8.8`. Show the rule.
2. Find your public IPv4 (trusted lookup). Compare it with your LAN IPv4. Write which range each address belongs to.
3. Write six sentences: a company uses `10.0.0.0/8` internally and one public address at the NAT edge.

#### Advanced practical tasks

1. Read a short public summary of RFC 1918 (not a full RFC dump). Write eight STE sentences.
2. Write a one-page plan: pick a private range for a lab with three subnets. Avoid overlap with a common home `192.168.1.0/24` if you will VPN later.

---

## Loopback `127.0.0.1`

The loopback range is `127.0.0.0/8`. Packets to this range stay in the host. They do not go out on a physical interface. `127.0.0.1` is the common address. The name `localhost` often maps to `127.0.0.1` and to `::1`.

A server can bind to `127.0.0.1` so that only programs on the same host can connect. A server that binds to `0.0.0.0` listens on all IPv4 interfaces.

`ping 127.0.0.1` tests the IPv4 stack in the host. It does not test the cable or the LAN.

Do not use `127.0.0.1` as the address of another machine. Do not publish a service on loopback when you want LAN clients to connect.

IPv6 loopback is `::1`. It is not `127.0.0.1`. Dual-stack hosts have both.

### Questions

#### Theoretical questions

1. What is the IPv4 loopback prefix?
2. Why does a loopback packet not leave the host?
3. What does a bind to `127.0.0.1` mean for remote clients?
4. What does `ping 127.0.0.1` test?
5. What is the IPv6 loopback address?

#### Easy practical tasks

1. Run `ping 127.0.0.1`. Write the result.
2. Resolve `localhost` with `ping localhost` or `getent` / `nslookup`. Write the address or addresses.
3. Write five sentences that compare `127.0.0.1` and your LAN address.
4. Draw a process on a host that talks to `127.0.0.1`. Show that the path stays inside the host.

#### Medium practical tasks

1. Start a local server bound to `127.0.0.1`. Try to fetch it from the same host and, if you can, from another device. Write both results.
2. Compare `ping 127.0.0.1` and `ping` to your LAN address. Write what each test proves.
3. Write six sentences about `0.0.0.0` as a listen address versus `127.0.0.1`.

#### Advanced practical tasks

1. Find a process that listens on `127.0.0.1` (`ss` or `netstat`). Write the port and why loopback is a good choice for that service if you can tell.
2. Write a one-page note: when a developer binds to loopback by mistake and a teammate cannot connect on the LAN.

---

## Broadcast and multicast (idea)

A unicast address names one interface. A broadcast address names all hosts on a subnet (or all hosts on the local link for limited broadcast `255.255.255.255`). A multicast address names a group of interested hosts.

IPv4 limited broadcast is `255.255.255.255`. A directed broadcast uses the subnet broadcast address, for example `192.168.1.255` for `192.168.1.0/24`. Many networks block directed broadcast.

IPv4 multicast uses `224.0.0.0/4`. Example: `224.0.0.1` is all hosts on the local subnet. Routing of multicast across networks needs extra protocols. This topic only needs the idea.

ARP uses link broadcast. DHCP uses broadcast or a mix of broadcast and unicast. Video and some discovery protocols use multicast.

Do not flood multicast on a WAN without a design. Do not treat broadcast as a way to reach the Internet.

IPv6 has no broadcast. IPv6 uses multicast for the jobs that IPv4 broadcast did.

### Questions

#### Theoretical questions

1. What is unicast?
2. What is broadcast?
3. What is multicast?
4. What is the IPv4 limited broadcast address?
5. Does IPv6 use broadcast? Explain.

#### Easy practical tasks

1. For `10.0.0.0/24`, write the directed broadcast address.
2. Write five sentences that compare unicast, broadcast, and multicast.
3. Make a table: address, type. Add `8.8.8.8`, `255.255.255.255`, `224.0.0.1`, `192.168.1.255`.
4. List two protocols that use broadcast or multicast on a LAN.

#### Medium practical tasks

1. In a capture, find one frame to `ff:ff:ff:ff:ff:ff`. Write the IPv4 destination if the payload is IP or ARP.
2. Write six sentences: why a switch floods broadcast, and why a router does not forward `255.255.255.255`.
3. Find `224.0.0.251` or similar in a capture if present (mDNS). Write the idea: group address, not one host.

#### Advanced practical tasks

1. Write a one-page note: DHCP discover as broadcast (idea). Map it to UDP ports only at a high level.
2. Explain why IPv6 neighbor discovery uses multicast, not broadcast. Use facts from this topic and topic 3.

---

## Default gateway

The default gateway is the next-hop IP address that a host uses when the destination is not local. The host ARPs for the gateway MAC. The host sends the frame to that MAC. The IP destination in the packet stays the remote address.

On a home LAN the gateway is often `.1` in the subnet, for example `192.168.1.1`. That is a habit, not a rule. The gateway must be an address on the same subnet as the host.

A host can have more than one route. The default route is `0.0.0.0/0` on IPv4. The next topic covers the routing table. This section only needs the idea: one last-resort next hop.

If the gateway is wrong, LAN ping can work and Internet ping fails. If the gateway is down, the same symptom appears. If the host mask is wrong, the host can treat remote addresses as local and skip the gateway.

Do not set a gateway that is not on the local subnet. Do not confuse the gateway IP with the public IP of the NAT box. They can be different addresses on different interfaces of the same box.

### Questions

#### Theoretical questions

1. When does a host use the default gateway?
2. Does the IP destination change when the host sends to the gateway? Explain.
3. Why must the gateway address sit on the local subnet?
4. What IPv4 prefix is the default route?
5. Why can LAN ping work when the gateway is wrong?

#### Easy practical tasks

1. Find your IPv4 default gateway (`ipconfig` or `ip route`). Write the address.
2. Ping the gateway. Write the result.
3. Write five sentences that explain ARP for the gateway, not for the remote web server.
4. Draw: host, gateway, Internet cloud. Label local destination versus remote destination.

#### Medium practical tasks

1. Compare gateway IP, LAN host IP, and public IP (lookup). Write the role of each address.
2. Write six sentences: a host with no default gateway. What still works, and what fails.
3. On Windows `route print` or on Unix `ip route`, find the default route line. Copy it and name each field that you know.

#### Advanced practical tasks

1. In a capture of a ping to `8.8.8.8`, write the Ethernet destination MAC and the IP destination. Explain why they do not name the same device.
2. Write a one-page troubleshooting guide: LAN works, Internet fails. Include mask, gateway, and NAT as later checks. Do not solve NAT yet.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe how a host uses mask, network part, and gateway together when it sends to a LAN neighbor and when it sends to a public address.
2. Why do private ranges, loopback, and broadcast need different rules?
3. How do dotted decimal and CIDR describe the same 32-bit value in two ways?
4. A teammate writes `192.168.1.0` as a host address on `/24`. Which facts do you use to correct that choice?
5. How does multicast differ from "send to the gateway"?

#### Easy practical tasks

1. Write a one-page cheat sheet: binary, mask, CIDR, RFC 1918, `127.0.0.1`, broadcast, multicast, gateway.
2. From your PC, record: IPv4 address, mask or prefix, gateway, and whether the address is private.
3. Convert your LAN address to binary. Mark network bits and host bits.
4. Draw `192.168.1.0/24` as a line of addresses. Mark network, one host, gateway, and broadcast.

#### Medium practical tasks

1. On paper, subnet a `/24` into four `/26`s. Assign a gateway in each subnet. Write the usable ranges.
2. Write a lab: two hosts on the same `/24`, then change one host to a different `/24`. Predict ping. Run only on machines that you own.
3. Build a table of ten IPv4 addresses that you see in a day. Label each as public, private, loopback, or special.

#### Advanced practical tasks

1. Write a small tool that takes address plus prefix and prints network, broadcast, first host, last host, and gateway candidate `.1` if it is usable.
2. Capture one LAN ping and one Internet ping. For each, write IP source, IP destination, Ethernet destination, and whether a gateway MAC was used.
