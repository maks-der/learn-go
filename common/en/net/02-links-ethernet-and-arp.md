# 2. Links, Ethernet, and ARP

## Description

The physical layer sends bits. The link layer sends frames on one hop. This topic shows bits on wire, fiber, and radio. You learn MAC address and Ethernet frame. You compare switch, hub, and router. You learn ARP and IPv6 neighbor discovery. You learn MTU and VLANs at a high level.

Complete this topic after networks and layers. Complete this topic before you study IPv4 and IPv6 addresses in depth.

Use one term for each concept. A bit is a 0 or a 1 on the medium. A frame is the link-layer unit. A MAC address is the link-layer address of an interface. ARP maps IPv4 to MAC on one link. Neighbor discovery maps IPv6 to MAC on one link. Do not mix a switch with a router. Do not mix a MAC address with an IP address.

You do not design a PHY chip. You name the medium, read a frame, and see why a switch is not a router.

---

## Bits on wire, fiber, and radio

The physical layer turns bits into a signal. A copper wire uses voltage. An optical fiber uses light. A radio link uses electromagnetic waves. The receiver turns the signal back into bits.

Each medium has limits. Copper has a length limit and noise. Fiber has a longer reach and uses light. Radio has shared air, interference, and walls. Capacity and error rate depend on the medium and on the PHY standard.

A network interface has a physical attachment. Examples: an RJ45 port, an SFP cage, a Wi-Fi radio. The driver and the firmware talk to the PHY. The operating system shows the interface as up or down.

Link speed is the bit rate of the PHY. Example: 1 Gbit/s Ethernet. That number is not the application throughput. Headers and gaps use part of the capacity. Radio links change rate when the signal is weak.

Duplex matters on some copper Ethernet links. Full duplex lets both sides send at the same time. Half duplex shares the medium and can collide. Modern switched Ethernet is full duplex. Wi-Fi is a shared radio medium with its own access method.

Do not mix "the cable is bad" with "the IP address is wrong." A down interface is a physical or link fault. A wrong gateway is not a bit-on-wire fault.

### Questions

#### Theoretical questions

1. What does the physical layer do with bits?
2. Name three media and the signal that each one uses.
3. Why is link speed not the same as application throughput?
4. What does an interface down state often mean?
5. What is the difference between full duplex and half duplex?

#### Easy practical tasks

1. Look at one cable or one Wi-Fi icon on your machine. Write the medium: copper, fiber, or radio.
2. Run `ipconfig` or `ip addr`. Write the interface name and whether the interface is up.
3. Write five sentences that compare copper and radio for a home LAN.
4. Make a table: medium, typical home use, one limit. Add three rows.

#### Medium practical tasks

1. Find the negotiated speed of your Ethernet or Wi-Fi interface in the OS. Write the value and the tool that showed it.
2. Unplug a cable or disable Wi-Fi for a short test that you own. Record the interface state. Restore the link.
3. Write six sentences: why a weak Wi-Fi signal can drop throughput even when the IP address is correct.

#### Advanced practical tasks

1. Compare two PHY rates on paper (for example 100 Mbit/s and 1 Gbit/s). Estimate the time to send 10 MB of payload. Ignore headers first, then add 10 percent overhead.
2. Write a one-page lab: how to prove a fault is physical (cable, radio, NIC) before you change IP settings.

---

## MAC address and Ethernet frame

A MAC address is a 48-bit link-layer address on common Ethernet and Wi-Fi interfaces. The usual text form is six octets in hex, for example `aa:bb:cc:dd:ee:ff`. Some systems use a hyphen or no separator.

The first bits mark unicast, multicast, and locally administered addresses. A unicast MAC names one interface. A broadcast MAC is `ff:ff:ff:ff:ff:ff`. A multicast MAC names a group on the link.

An Ethernet frame carries a destination MAC, a source MAC, a type field, a payload, and a frame check sequence. The type field names the next protocol. Example: IPv4 or IPv6. Wi-Fi frames use a different header. The idea stays the same: link addresses plus a payload.

A switch learns unicast MAC addresses per port. The switch forwards a known unicast frame to one port. The switch floods an unknown unicast frame.

A MAC address is not an IP address. The MAC address is valid on one link or one VLAN bridge domain. The IP address can stay the same when the host moves to another link, if the operator assigns it. The MAC address usually stays with the NIC, unless the OS or the user sets a different address.

Privacy features on phones can change the Wi-Fi MAC. That change is a local choice. It does not change the meaning of a MAC address.

Do not type a MAC address into a browser as if it were an IP address. Do not assume that two hosts with the same IP on different links share a MAC.

### Questions

#### Theoretical questions

1. How many bits does a common Ethernet MAC address use?
2. What is the broadcast MAC address?
3. What fields does an Ethernet frame have in this section?
4. Why is a MAC address not an IP address?
5. What does a switch learn from a source MAC?

#### Easy practical tasks

1. Find the MAC address of one interface on your machine. Write it in hex.
2. Write five sentences that explain unicast, multicast, and broadcast MAC addresses.
3. Make a table: address type, example, who receives the frame.
4. Draw a NIC and label the MAC address on the link side and the IP address on the internet side.

#### Medium practical tasks

1. Compare the MAC of your Wi-Fi interface and your Ethernet interface if both exist. Write why they differ.
2. In a capture, find the Ethernet source and destination MAC of one frame. Write both values.
3. Write six sentences: what happens when two hosts use the same MAC on one LAN (the idea of a conflict).

#### Advanced practical tasks

1. Find whether your OS shows a locally administered bit or a randomized Wi-Fi MAC. Write the setting name and the current address.
2. Write a one-page note: how a switch table uses MAC addresses, and why that table is not a routing table.

---

## Switch vs hub vs router

A hub is a simple repeater. A frame that enters one port exits all other ports. Hubs are rare on modern LANs. They create one collision domain on half-duplex copper. You still meet the word in old notes.

A switch forwards Ethernet frames. It learns source MAC addresses. It sends a known unicast to one port. It floods broadcast and unknown unicast. Hosts on a switch do not share one collision domain in the usual full-duplex case. They still share one broadcast domain unless you use VLANs.

A router forwards IP packets between networks. It uses a routing table. It strips the incoming frame and builds a new frame on the outgoing link. The IP source and destination stay the same in the usual case (until NAT). A home "router" is often a switch plus a router plus NAT in one box. Keep the jobs separate in your notes.

A host on one subnet talks to another host on the same subnet through a switch (or a Wi-Fi access point). A host talks to a different subnet through a router. The default gateway is that router.

Do not call a layer-3 switch a mystery. A layer-3 switch is a switch that can also route. The idea is still: frames on a VLAN, packets between VLANs.

Do not expect a switch to understand TCP ports. Do not expect a hub to isolate unicast traffic.

### Questions

#### Theoretical questions

1. What does a hub do with a frame?
2. What does a switch learn?
3. What does a router forward?
4. When does a host use a gateway?
5. Why is a home router more than one job?

#### Easy practical tasks

1. Write five sentences that compare hub, switch, and router.
2. Make a table: device, unit it forwards, address it uses.
3. Draw a LAN: two PCs, one switch, one router to the Internet.
4. Label your home box as switch ports plus a routed WAN port (idea).

#### Medium practical tasks

1. Find the default gateway on your machine. Write whether that device is the next hop off the LAN.
2. Write six sentences: two hosts on one switch with the same subnet versus two different subnets.
3. Trace one `ping` to a LAN neighbor and one `ping` to an Internet address. Write which path needs a router.

#### Advanced practical tasks

1. Write a one-page note: broadcast domain versus collision domain. Place hub, switch, and VLAN in the note.
2. On a lab that you own, compare `arp` or neighbor output before and after a ping to a LAN host. Write what appeared.

---

## ARP and IPv6 neighbor discovery

ARP is the Address Resolution Protocol for IPv4. A host has a next-hop IPv4 address. The host needs the MAC address of that next hop on this link. ARP asks: who has this IPv4 address? The owner replies with its MAC. The host stores a cache entry.

ARP runs only on the local link. A host does not ARP for a remote IPv4 address. The host ARPs for the gateway MAC, then sends the packet to that MAC with the remote IP in the IP header.

An ARP cache can go stale. A wrong static ARP entry can black-hole a host. Duplicate IPv4 addresses cause ARP fights. Two hosts claim one address.

IPv6 does not use ARP. IPv6 uses neighbor discovery (ND). ND uses ICMPv6. Neighbor Solicitation and Neighbor Advertisement map IPv6 addresses to MAC addresses. Router Solicitation and Router Advertisement find routers. Topic 3 covers link-local `fe80::` and SLAAC.

ND and ARP solve the same job: map an internet-layer address to a link-layer address on one hop.

Do not filter all ICMP on IPv6 if you want ND to work. Do not expect `arp -a` to show IPv6 neighbors. Use `ip -6 neigh` or the Windows neighbor table.

### Questions

#### Theoretical questions

1. What two addresses does ARP map?
2. Why does a host ARP for the gateway and not for a remote server IP?
3. What protocol family does IPv6 neighbor discovery use?
4. What happens if two hosts claim the same IPv4 address?
5. Why is ARP limited to one link?

#### Easy practical tasks

1. Run `arp -a` or `ip neigh`. Write one IPv4 address and its MAC if a line exists.
2. Write five sentences that explain an ARP request and an ARP reply.
3. Make a table: IPv4 map, IPv6 map. Add protocol names.
4. Draw: host, gateway IP, ARP question, gateway MAC.

#### Medium practical tasks

1. Ping a LAN neighbor. Then read the neighbor table. Write the new or refreshed line.
2. Write six sentences: packet to 8.8.8.8 on a LAN. Who is the Ethernet destination?
3. In a capture, find an ARP or an ICMPv6 neighbor message. Write the names the tool shows.

#### Advanced practical tasks

1. Write a one-page comparison of ARP and ND. Include cache, ICMP, and router discovery.
2. On a lab that you own, clear one neighbor entry if the OS allows it, ping again, and capture the resolve. Document commands.

---

## MTU and VLANs (high-level)

MTU is the maximum transmission unit. For IPv4 on Ethernet the common MTU is 1500 bytes of IP packet. A larger IP packet must fragment or the sender must use a smaller size. IPv6 hosts do not fragment in the same way. Path MTU discovery finds a safe size. ICMP "packet too big" matters. A filter that drops ICMP can break large transfers.

The Ethernet payload must fit the MTU. TCP often sets MSS from the MTU so that a segment plus headers fits one frame.

A VLAN is a virtual LAN. A switch tags frames with a VLAN ID (802.1Q). Hosts in VLAN 10 do not share a broadcast domain with hosts in VLAN 20 even on the same switch. A router (or a layer-3 switch) forwards between VLANs.

A trunk link carries several VLAN tags. An access port is usually untagged for one VLAN. This topic is high-level. You do not design a campus fabric here.

Do not set a 9000-byte jumbo MTU on one host only. Both ends and the path must agree. Do not treat a VLAN as a security boundary by itself. A mis-tag or a router leak can join the VLANs.

### Questions

#### Theoretical questions

1. What does MTU measure?
2. What is a common Ethernet IPv4 MTU?
3. Why can a drop of ICMP hurt large packets?
4. What does a VLAN ID separate?
5. What is a trunk link in this section?

#### Easy practical tasks

1. Find the MTU of one interface (`ip link` or interface details). Write the number.
2. Write five sentences that explain VLAN as a split broadcast domain.
3. Make a table: access port, trunk port. Add one sentence each.
4. Draw two VLANs and one router between them.

#### Medium practical tasks

1. Write six sentences: MSS and MTU for TCP on Ethernet.
2. Ping with a large size if your OS allows it (`ping -l` or `ping -s`). Write whether the command succeeds. Use a host that you own.
3. Read a home or lab switch UI if you have one. Write whether VLANs exist.

#### Advanced practical tasks

1. Write a one-page note: path MTU discovery on IPv4 and IPv6 at a high level. No exploit steps.
2. In a lab that you own, change MTU on a test interface and ping. Restore the MTU. Record results.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A ping to a LAN IP fails. How do you decide among cable, MAC conflict, and wrong subnet?
2. Why must a router use both a MAC table idea on each LAN and a routing table between LANs?
3. How do ARP and ND fit the encapsulation story from topic 1?
4. Why does a VLAN tag sit at the link layer and not in the TCP header?
5. What stays local to one hop: MAC, IP, or both, and when?

#### Easy practical tasks

1. Write a cheat sheet: medium, MAC, frame, switch, hub, router, ARP, ND, MTU, VLAN.
2. Print your MAC, IPv4, and gateway. Write which values are link-local in meaning.
3. Draw one frame that carries an IPv4 packet to the gateway MAC.
4. Run `arp -a` or `ip neigh` after a web fetch. Write how many neighbors you see.

#### Medium practical tasks

1. Write a lab report: prove your next hop is a router, not a switch-only device. Use addresses and a traceroute hop.
2. Capture one ARP or ND exchange and one IPv4 or IPv6 data frame. Write the Ethernet type field if shown.
3. Compare Wi-Fi and Ethernet on your machine: speed, medium, MAC. Write three differences.

#### Advanced practical tasks

1. Write a one-page campus picture: two access VLANs, one trunk, one router, one MTU policy.
2. Trace a packet from your process to the first hop: IP lookup, ARP or ND, frame on the medium. Use tools from this topic.
