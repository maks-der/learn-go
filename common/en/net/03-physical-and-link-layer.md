# 3. Physical and Link Layer

## Description

The physical layer sends bits. The link layer sends frames on one hop. This topic shows media, MAC addresses, Ethernet frames, switches, hubs, routers, ARP, neighbor discovery, MTU, and VLANs.

Use one term for each concept. A bit is a 0 or a 1 on the medium. A frame is the link-layer unit. A MAC address is the link-layer address of an interface. Complete this topic before you study IPv4 addressing in depth.

You do not need to design a PHY chip. You need to name the medium, read a frame, and see why a switch is not a router.

---

## Bits on wire / fiber / radio

The physical layer turns bits into a signal. A copper wire uses voltage. An optical fiber uses light. A radio link uses electromagnetic waves. The receiver turns the signal back into bits.

Each medium has limits. Copper has a length limit and noise. Fiber has a longer reach and uses light. Radio has shared air, interference, and walls. Capacity and error rate depend on the medium and on the PHY standard.

A network interface has a physical attachment. Examples: an RJ45 port, an SFP cage, a Wi-Fi radio. The driver and the firmware talk to the PHY. The operating system shows the interface as up or down.

Link speed is the bit rate of the PHY. Example: 1 Gbit/s Ethernet. That number is not the application throughput. Headers and gaps use part of the capacity. Radio links change rate when the signal is weak.

Do not mix "the cable is bad" with "the IP address is wrong." A down interface is a physical or link fault. A wrong gateway is not a bit-on-wire fault.

Duplex matters on some copper Ethernet links. Full duplex lets both sides send at the same time. Half duplex shares the medium and can collide. Modern switched Ethernet is full duplex. Wi-Fi is a shared radio medium with its own access method.

### Questions

#### Theoretical questions

1. What does the physical layer do with bits?
2. Name three media and the signal that each one uses.
3. Why is link speed not the same as application throughput?
4. What does an interface "down" state often mean?
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

## MAC address

A MAC address is a 48-bit link-layer address on common Ethernet and Wi-Fi interfaces. The usual text form is six octets in hex, for example `aa:bb:cc:dd:ee:ff`. Some systems use a hyphen or no separator.

The first bits mark unicast, multicast, and locally administered addresses. A unicast MAC names one interface. A broadcast MAC is `ff:ff:ff:ff:ff:ff`. A multicast MAC names a group on the link.

A switch learns unicast MAC addresses per port. The switch forwards a known unicast frame to one port. The switch floods an unknown unicast frame. Later in this topic you compare switch, hub, and router.

A MAC address is not an IP address. The MAC address is valid on one link (or one VLAN bridge domain). The IP address can stay the same when the host moves to another link, if the operator assigns it. The MAC address usually stays with the NIC, unless the OS or the user sets a different address.

Privacy features on phones can change the Wi-Fi MAC. That change is a local choice. It does not change the meaning of a MAC address.

Do not type a MAC address into a browser as if it were an IP address. Do not assume that two hosts with the same IP on different links share a MAC.

### Questions

#### Theoretical questions

1. How many bits does a common Ethernet MAC address use?
2. What is the broadcast MAC address?
3. What is the difference between a unicast MAC and a multicast MAC?
4. Why is a MAC address not an IP address?
5. What does a switch learn from a source MAC?

#### Easy practical tasks

1. Find the MAC address of one interface on your machine. Write it in hex.
2. Write five sentences that explain unicast, multicast, and broadcast MAC addresses.
3. Make a table: address type, example, and who receives the frame.
4. Draw a NIC and label the MAC address on the link side and the IP address on the internet side.

#### Medium practical tasks

1. Compare the MAC of your Wi-Fi interface and your Ethernet interface if both exist. Write why they differ.
2. In a capture, find the Ethernet source and destination MAC of one frame. Write both values.
3. Write six sentences: what happens when two hosts use the same MAC on one LAN (the idea of a conflict).

#### Advanced practical tasks

1. Find whether your OS shows a locally administered bit or a randomized Wi-Fi MAC. Write the setting name and the current address.
2. Write a one-page note: how a switch table uses MAC addresses, and why that table is not a routing table.

---

## Ethernet frame

An Ethernet frame on the wire has a preamble and start frame delimiter, then the frame that Wireshark often shows: destination MAC, source MAC, optional 802.1Q tag, Ethertype, payload, and FCS.

The Ethertype names the payload. Common values: `0x0800` for IPv4, `0x86DD` for IPv6, `0x0806` for ARP. The payload is the IP packet or the ARP message.

The minimum payload size on classic Ethernet is 46 bytes (without a VLAN tag). The sender pads a short payload. The maximum payload is the MTU, often 1500 bytes for the IP packet.

The FCS is a CRC. The receiver drops a frame with a bad FCS. You often do not see the FCS in a capture on the sending host. A capture on the wire or on some NICs can show errors.

A Wi-Fi link uses 802.11 frames, not the same header as Ethernet. The OS often presents a fake Ethernet header to programs and to many capture setups. This handbook uses Ethernet as the main teaching frame.

Do not call a frame a packet when you talk about the link. Use frame at the link. Use packet at IP.

### Questions

#### Theoretical questions

1. Which two MAC addresses sit at the start of a common Ethernet header?
2. What does the Ethertype field name?
3. What is the usual IPv4 Ethertype value?
4. What does the FCS do?
5. What is the usual maximum IP payload of a standard Ethernet frame?

#### Easy practical tasks

1. Draw an Ethernet header: dest MAC, source MAC, Ethertype, payload, FCS.
2. In Wireshark, open one frame. Write the Ethertype and the protocol name that the tool shows.
3. Write five sentences that explain padding of a short Ethernet payload.
4. Make a table: Ethertype, protocol. Add IPv4, IPv6, and ARP.

#### Medium practical tasks

1. Capture one ARP and one IPv4 frame. Write the Ethertype of each frame.
2. Find the frame length of a small ICMP echo. Write whether the tool shows padding.
3. Write six sentences that compare an Ethernet frame and an IP packet. Use header names.

#### Advanced practical tasks

1. Read the Ethernet header of one captured frame as hex. Mark the destination MAC, the source MAC, and the Ethertype bytes.
2. Write a one-page note: why a capture on the host may hide the preamble and the FCS.

---

## Switch vs hub vs router

A hub is a physical repeater. Every port shares one collision domain. A frame on one port goes to all other ports. Hubs are rare on modern LANs. You still learn the word so that you do not confuse it with a switch.

A switch forwards frames by MAC address. Each port is its own collision domain on modern Ethernet. The switch learns the source MAC on a port. The switch sends a known unicast to one port. The switch floods broadcast, multicast (unless IGMP snooping), and unknown unicast.

A router forwards packets by IP address. The router decapsulates the frame, reads the IP header, chooses the next hop, and encapsulates a new frame. A home "router" box often contains a switch, a Wi-Fi access point, a router, and NAT. The word on the box is not the function.

Use the function, not the plastic box:

- hub: repeat bits to all ports
- switch: forward frames by MAC
- router: forward packets by IP

A layer-3 switch is a switch that can also route. Treat it as a device with both jobs. Do not call a pure layer-2 switch a router.

Broadcast stays on the layer-2 domain. A router does not forward Ethernet broadcast to another IP network. That fact is why ARP stays local.

### Questions

#### Theoretical questions

1. What does a hub do with a frame?
2. How does a switch decide the output port for a known unicast?
3. How does a router decide the next hop?
4. Why does a home gateway box confuse beginners?
5. Does a router forward Ethernet broadcast to another network? Explain.

#### Easy practical tasks

1. Make a three-column table: hub, switch, router. Add one row for address type and one row for broadcast.
2. Write five sentences that explain learning of MAC addresses on a switch.
3. Draw a LAN: two PCs, one switch, one router. Label which device uses MAC and which device uses IP.
4. Look at your home gateway. Write which functions you think it includes (switch, Wi-Fi, router, NAT).

#### Medium practical tasks

1. Find the default gateway MAC in your ARP table (`arp -a` or `ip neigh`). Write why that MAC is the gateway on your LAN, not the remote web server.
2. Write six sentences: what a switch does with a frame to `ff:ff:ff:ff:ff:ff`.
3. Compare a cheap five-port switch and a home router in a short table of features that you can see (WAN port, LAN ports, Wi-Fi).

#### Advanced practical tasks

1. Write a one-page note: why two hosts on the same switch use ARP for each other, but two hosts on different IP networks do not send frames directly.
2. Design a small lab with a switch and a router (real or drawn). List three tests that prove which device is which.

---

## ARP (IPv4) / neighbor discovery (IPv6)

ARP maps an IPv4 address to a MAC address on the local link. The sender broadcasts an ARP request: who has this IPv4 address? The owner unicasts an ARP reply with its MAC. The sender stores the pair in an ARP cache.

You need ARP when the next hop is on the same IPv4 network. The next hop can be the destination host or the default gateway. You do not ARP for a public IPv4 address that is not on your LAN. You ARP for the gateway address.

Neighbor Discovery (ND) is the IPv6 function that replaces ARP. ND uses ICMPv6. A neighbor solicitation asks for the MAC of an IPv6 address. A neighbor advertisement answers. Router advertisements are part of ND and appear again in the IPv6 topic.

The neighbor cache is the IPv6 table of address and MAC pairs. The ARP table is the IPv4 table.

Stale or wrong cache entries cause faults. A host can send frames to the wrong MAC. Clearing the cache (when you own the machine) is a common lab step.

Do not use ARP on IPv6. Do not use ND as if it were IPv4 ARP on the wire. The job is the same. The messages are not the same.

### Questions

#### Theoretical questions

1. What two values does ARP map?
2. Why is an ARP request a broadcast?
3. When do you ARP for the gateway instead of the destination?
4. What IPv6 protocol family carries neighbor discovery?
5. What is a neighbor cache?

#### Easy practical tasks

1. Run `arp -a` or `ip neigh`. Write one IPv4 address and its MAC.
2. Write five sentences that describe an ARP request and an ARP reply.
3. Draw: host A, host B on the same LAN. Show who sends broadcast and who sends unicast.
4. Write one sentence that states when you must not ARP for a remote Internet address.

#### Medium practical tasks

1. Capture an ARP exchange. Write the Ethertype, the requested IPv4 address, and the MAC in the reply.
2. Ping a LAN neighbor that you have not used recently. Watch the ARP table before and after. Write the change.
3. Write six sentences that compare ARP and neighbor solicitation.

#### Advanced practical tasks

1. Clear the ARP entry for your gateway on a machine that you own. Ping the gateway. Capture the new ARP. Record the commands.
2. Write a one-page note: a wrong static ARP entry (lab only). What symptom does a user see? Restore the automatic ARP.

---

## MTU and fragmentation (preview)

The MTU is the maximum IP packet size that a link can carry in one frame, in bytes. Standard Ethernet IPv4 MTU is 1500. The frame can be larger than 1500 because of Ethernet headers.

If an IPv4 packet is larger than the next MTU, a router or the sender can fragment the packet when fragmentation is allowed. Each fragment has an IP header. The destination host reassembles the fragments.

IPv6 routers do not fragment packets. The source host must use path MTU discovery or a safe size. This topic is a preview. Later IPv6 and extra topics go deeper.

A too-small MTU causes odd faults: a `ping` of a small size works, a large TCP session hangs, or a VPN drops. "ICMP fragmentation needed" or IPv6 packet-too-big messages help path MTU discovery.

Do not set a random MTU on a production link without a measurement. Do not confuse Ethernet frame size with IP MTU.

Jumbo frames use a larger MTU on a LAN that all devices support. The Internet path often stays at 1500 or less (tunnels can be smaller).

### Questions

#### Theoretical questions

1. What is an MTU?
2. What is a common Ethernet IPv4 MTU?
3. Who reassembles IPv4 fragments?
4. Do IPv6 routers fragment packets?
5. Why can a small ping work when a large transfer fails?

#### Easy practical tasks

1. Find the MTU of one interface (`ipconfig` / `netsh` or `ip link`). Write the number.
2. Write five sentences that explain why a frame can be larger than the MTU.
3. Draw one 4000-byte IPv4 packet that must split on a 1500 MTU link. Show two or more fragments as boxes.
4. Make a table: IPv4 fragmentation at a router, IPv6 fragmentation at a router. Write allowed or not.

#### Medium practical tasks

1. Ping with a payload size near 1400 and near 1500 if your OS lets you set the size. Write which command you used and what happened.
2. In a capture, find the total length of an IP packet. Compare it with the interface MTU.
3. Write six sentences: what a tunnel does to the effective MTU (the idea only).

#### Advanced practical tasks

1. Read a public short explanation of path MTU discovery. Write five STE sentences. No full RFC copy.
2. Write a lab plan to find the path MTU to a public host with ping size and the DF bit (IPv4). Record results. Do not flood a network that you do not own.

---

## VLANs (high-level)

A VLAN is a virtual LAN. One physical switch can split ports into separate broadcast domains. Frames in VLAN 10 do not flood to VLAN 20 unless a router (or a layer-3 switch) forwards IP packets between them.

An 802.1Q tag sits in the Ethernet header. The tag has a VLAN identifier (VID), often 1 to 4094. A trunk port carries tagged frames for many VLANs. An access port belongs to one VLAN and often sends untagged frames to a PC.

VLANs separate traffic. A guest Wi-Fi VLAN should not see the office printer VLAN. Separation is not a full security system. You still need routing policy and firewalls.

A beginner home LAN often has no visible VLAN. Many home boxes use one LAN. Some boxes add a guest network. That guest network is a VLAN idea.

Do not treat a VLAN as a WAN. Do not treat a VLAN ID as an IP subnet. Operators often map one VLAN to one subnet, but the two numbers are different.

### Questions

#### Theoretical questions

1. What is a VLAN?
2. What does an 802.1Q tag carry?
3. What is the difference between an access port and a trunk port?
4. Why do frames in VLAN 10 not flood into VLAN 20?
5. Is a VLAN ID the same as a subnet mask? Explain.

#### Easy practical tasks

1. Write five sentences that explain why a guest Wi-Fi can use a VLAN.
2. Draw one switch with VLAN 10 and VLAN 20. Show a router that connects the two.
3. Make a table: access port, trunk port. Add one row for tags and one row for typical device.
4. List three reasons an operator uses VLANs.

#### Medium practical tasks

1. Find whether your home or lab switch UI shows VLAN IDs. Write what you see, or write that the box hides VLANs.
2. In a public 802.1Q diagram, write the order: dest MAC, source MAC, tag, Ethertype.
3. Write six sentences: how a VLAN and an IP subnet work together when an operator maps them one to one.

#### Advanced practical tasks

1. Write a one-page design: two VLANs (users, servers) and one router. List which ports are access and which port is a trunk.
2. Capture is optional: if you have a lab trunk, find a VLAN tag in Wireshark. Write the VID. If you have no lab, write how you would filter `vlan` in Wireshark.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of one IPv4 packet from host A to host B on the same Ethernet LAN. Use bits, frame, MAC, ARP, and MTU.
2. How do you tell a physical fault from a wrong VLAN or a wrong MAC table?
3. Why must a host use ARP or ND before it sends the first IP packet to a local next hop?
4. Compare a switch flood and a router forward in one short paragraph.
5. Which facts show that a home gateway is more than a layer-2 switch?

#### Easy practical tasks

1. Write a one-page cheat sheet: medium, MAC, Ethernet fields, switch, hub, router, ARP, ND, MTU, VLAN.
2. Collect from your PC: interface state, MAC, MTU, and one ARP or neighbor entry. Put them in one table.
3. Draw a frame on copper: signal as a line of bits, then the Ethernet header fields.
4. Label your home devices as hub, switch, router, or mixed box. Give one reason for each label.

#### Medium practical tasks

1. Capture ping to the gateway. Write the Ethernet MACs, the ARP that appeared if any, and the ICMP packet size versus MTU.
2. Write a troubleshooting list of eight checks from cable up to ARP. Order the checks from physical to link.
3. Explain in ten sentences how a VLAN-aware switch plus a router replaces two separate physical LANs.

#### Advanced practical tasks

1. Build a glossary of 12 link-layer terms. Each entry: term, one sentence, and one command or Wireshark field.
2. Write a lab that proves a switch is not a router: two subnets, one switch only, then add a router. Predict the ping result in each step. Run it only on equipment that you own.
