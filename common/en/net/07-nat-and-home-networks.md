# 7. NAT and Home Networks

## Description

NAT translates addresses (and often ports) at a network edge. Home networks use NAT because many hosts share one public IPv4 address. This topic shows why NAT is common, source NAT, port forwarding, peer-to-peer problems, carrier-grade NAT, and hairpin NAT.

Use one term for each concept. Source NAT changes the source address of packets that leave the LAN. Port forwarding maps a public port to a private host. Complete this topic before you study transport details that assume a clear 5-tuple.

NAT is not a full firewall. A stateful firewall can sit in the same box. Keep the two jobs separate in your notes.

---

## Why NAT is common

Public IPv4 addresses are scarce. A home or a small office gets one public IPv4 address, or a small set. Many phones, PCs, and TVs need addresses. RFC 1918 private addresses solve the local count. NAT at the edge lets those private hosts share the public address.

NAT also hides the internal layout from the public Internet. That hide is a side effect. It is not a complete security design. Hosts still need updates and a firewall.

IPv6 can give each host a global address. Many homes still use IPv4 NAT because the ISP path, the TV, or the old printer is IPv4-only. Dual stack can use NAT on IPv4 and no NAT on IPv6.

You will see NAT in almost every home lab. Your LAN address is private. Your "what is my IP" page shows the NAT public address. Those two addresses are the reason this topic exists.

Do not treat NAT as the same as a proxy. A proxy understands the application. Basic NAT rewrites IP and transport ports.

Do not assign public IPv4 meaning to `192.168.1.10`. That address is local.

### Questions

#### Theoretical questions

1. Why do homes use private IPv4 plus NAT?
2. Why is NAT not a complete security design?
3. How can IPv6 reduce the need for NAT?
4. What two IPv4 addresses does a typical home user see?
5. How is NAT different from a proxy in this section?

#### Easy practical tasks

1. Write your LAN IPv4 and your public IPv4 (trusted lookup). Label which is private.
2. Write five sentences that explain IPv4 scarcity and NAT.
3. Draw a LAN with three hosts and one public address. Label NAT at the edge.
4. Make a table: device, private IP. Add three home devices.

#### Medium practical tasks

1. Count devices on your LAN that get a DHCP lease (router UI if you have access). Write why one public IPv4 is not enough without NAT.
2. Write six sentences: IPv4 NAT at home and IPv6 global addresses if your ISP offers IPv6.
3. Compare a "what is my IP" page with `ipconfig`. Write which programs see which address.

#### Advanced practical tasks

1. Write a one-page note for a home user: what NAT does, what a firewall does, and what NAT does not encrypt.
2. Inventory your home: public IPv4 count (often one), private range, and whether IPv6 is on. Use only access that you own.

---

## Source NAT / masquerade

Source NAT rewrites the source IP of a packet that leaves the inside network. Masquerade is source NAT that uses the current address of the outbound interface. Linux `iptables` and `nftables` call that masquerade. Home routers do the same job in the UI.

The box also rewrites the source port when two inside hosts would collide. The box stores a mapping: inside IP, inside port, outside IP, outside port, protocol, and the remote peer. Return packets match the mapping. The box rewrites the destination back to the inside host.

The mapping is stateful. If the mapping expires, a late reply does not reach the inside host. Idle timeouts differ for UDP and TCP.

From the Internet, all inside hosts can look like one public IP. Logs on a public server show the NAT address, not `192.168.1.10`.

Do not confuse source NAT with dest NAT. Source NAT changes the sender as seen outside. Dest NAT (port forwarding) changes the destination as seen inside.

Do not expect a new inbound connection from the Internet to reach a private host without a mapping. Source NAT creates mappings for outbound traffic first.

### Questions

#### Theoretical questions

1. What field does source NAT change on the outbound packet?
2. What is masquerade?
3. Why can the NAT box also change the source port?
4. What does a NAT mapping store (the idea)?
5. Why does a public server log show the public IP, not the LAN IP?

#### Easy practical tasks

1. Write five sentences that describe an outbound HTTP request through source NAT.
2. Draw before and after headers: source IP inside, source IP outside.
3. Make a table: inside 5-tuple, outside 5-tuple. Invent one TCP example.
4. Explain in four sentences why two PCs can browse at the same time through one public IP.

#### Medium practical tasks

1. Capture on your PC (inside). Write the source IP of an outbound packet. Then look at a public echo service or a server log if you have one. Compare sources.
2. Write six sentences: UDP mapping timeout versus a long TCP session (idea).
3. On a Linux lab that you own, find a masquerade or `snat` rule (`nft list ruleset` or similar). Write the line. Skip this if you have no Linux NAT box.

#### Advanced practical tasks

1. Write a one-page mapping table for three inside hosts that all talk to `8.8.8.8:53` over UDP. Show unique outside ports.
2. Capture WAN and LAN of a lab router that you own (or a VM NAT). Show the source rewrite. Do not capture on a network that you do not own.

---

## Port forwarding

Port forwarding is destination NAT for inbound traffic. The router listens on a public IP and a port. The router rewrites the destination to a private IP and port. Example: public `203.0.113.5:8080` goes to `192.168.1.20:80`.

You use port forwarding when a server on the LAN must accept a new connection from the Internet. Games, cameras, and self-hosted web servers often need this. Many vendors also call it "virtual server" or "inbound rule."

The forward is not only a rewrite. The firewall must allow the traffic. Some boxes combine the two settings. Write both in your notes: NAT map and allow rule.

A port can map to a different internal port. The public port and the private port need not match. The protocol (TCP or UDP) must match the service.

Do not forward a port to a host that you do not manage. Do not forward admin ports to the Internet without extra protection. Prefer a VPN when you can.

UPnP and PCP can install forwards automatically. That is convenient and risky. Know whether your box enables UPnP.

### Questions

#### Theoretical questions

1. What does port forwarding change on an inbound packet?
2. Why do you need a firewall allow as well as a NAT map?
3. Can the public port differ from the private port?
4. What problem does port forwarding solve that source NAT does not?
5. What is the risk of UPnP in one sentence?

#### Easy practical tasks

1. Write five sentences that explain a web server on `192.168.1.20` and a public port 8080.
2. Draw an inbound packet: public dest, then private dest after the rewrite.
3. Open your home router UI if you own it. Write whether a port-forward page exists. Do not expose a real service.
4. Make a table: public IP, public port, private IP, private port, protocol. Invent one row.

#### Medium practical tasks

1. On a lab that you own, set a forward to a local HTTP server. Test from a phone on mobile data if you can, not only from the LAN. Write the result.
2. Write six sentences: why a test from the LAN to the public IP can fail without hairpin NAT (preview the later section).
3. List three services that people forward and one reason each is risky.

#### Advanced practical tasks

1. Write a one-page hardening checklist for one forwarded TCP port (updates, bind address, firewall, logs).
2. Compare UPnP off versus a manual forward. Write a short policy for a home lab.

---

## NAT and peer-to-peer pain

Peer-to-peer programs want two hosts behind NAT to talk. Neither host has a public listen mapping at the start. Outbound source NAT does not create a mapping that the other peer can guess in a safe way for all NAT types.

NAT traversal techniques include STUN (learn the public mapping), TURN (relay through a public server), and ICE (choose a candidate). Many video and game systems use these. This topic needs the pain and the idea, not a full WebRTC course.

Symmetric NAT is hard: the mapping depends on the remote address. A hole that works for a STUN server may not work for the peer. Full-cone NAT is easier. Home boxes differ.

A relay always works if both peers can reach the relay. The cost is bandwidth on the relay. Direct paths are cheaper when NAT permits them.

Do not assume that "disable firewall" fixes P2P. The NAT mapping is a different problem. Do not open the whole box to the Internet to make a game work.

IPv6 end-to-end can remove this pain when both sides have global addresses and a firewall allow. Many peers still sit on IPv4 NAT.

### Questions

#### Theoretical questions

1. Why is it hard for two NATed hosts to start a session?
2. What does a STUN server tell a client (the idea)?
3. What does a TURN relay do?
4. Why is a relay a reliable fallback?
5. How can IPv6 reduce this pain?

#### Easy practical tasks

1. Write five sentences that explain two homes, two NATs, and one game.
2. Draw two NATs and a public relay. Show a path that always works.
3. Make a table: STUN, TURN, ICE. Add one sentence each.
4. List two apps that you use that might need NAT traversal (voice, game). Do not reverse engineer them.

#### Medium practical tasks

1. Write six sentences: outbound chat works, inbound unsolicited UDP fails. Use mappings.
2. Find a public short ICE or WebRTC NAT overview. Write four STE sentences. No long quotes.
3. Compare "call a friend" on mobile data versus on home Wi-Fi. Write one NAT-related hypothesis.

#### Advanced practical tasks

1. Write a one-page design: a small P2P file tool. State when you need a relay.
2. In a lab with two NAT VMs that you own, try a direct UDP echo and then a relay. Record which path works. No attack on third-party NATs.

---

## Carrier-grade NAT (awareness)

Carrier-grade NAT (CGNAT) is NAT in the ISP network. The customer may get a private or shared address on the WAN port of the home router, not a unique public IPv4. The range `100.64.0.0/10` (RFC 6598) is reserved for shared address space on service-provider NAT.

You still have a home NAT. The path can be NAT then NAT (double NAT). Port forwarding on the home box does not create a public mapping on the ISP NAT. Inbound services fail. Some ISPs sell a real public IPv4 as an extra product.

CGNAT also hurts some peer-to-peer paths and some VPNs that expect a public address. Logging at a public server shows an ISP address that many customers share. Abuse handling needs extra logs at the ISP.

IPv6 from the ISP is a way out. 464XLAT and similar transition tools appear in later topics. This section only needs awareness: your WAN address is not always public.

Do not blame only your home router when inbound ports fail. Check the WAN address. If it is `100.64` or another private range, you are behind CGNAT or another extra NAT.

### Questions

#### Theoretical questions

1. Where does CGNAT sit?
2. What is `100.64.0.0/10`?
3. Why does home port forwarding fail behind CGNAT?
4. Why does a public server log become harder to read under CGNAT?
5. What WAN address pattern should make you suspect extra NAT?

#### Easy practical tasks

1. Look at your router WAN IPv4 if you own the box. Write whether it looks public, RFC 1918, or `100.64/10`.
2. Write five sentences that explain double NAT.
3. Make a table: unique public IPv4, CGNAT. Add rows for port forward and for "what is my IP".
4. Draw ISP CGNAT, home NAT, and a PC.

#### Medium practical tasks

1. Compare router WAN IP and a "what is my IP" page. Write if they match. If they differ, write a possible reason.
2. Write six sentences: a home camera app that uses a vendor relay because of CGNAT.
3. Find your ISP IPv6 status. Write whether IPv6 could avoid IPv4 inbound problems.

#### Advanced practical tasks

1. Write a one-page customer note: how to detect CGNAT and what to ask the ISP (public IPv4 or IPv6).
2. Research 464XLAT at a high level from a public source. Write five STE sentences. No implementation of a transition hack on a foreign network.

---

## Hairpin NAT

Hairpin NAT (NAT loopback, NAT reflection) is a rewrite for packets that stay in the site but use the public address as the destination. Example: a LAN phone opens `http://203.0.113.5:8080`, which is a port forward to a LAN web server. The packet must go to the router and come back to the server. The router must translate source or destination so that the return path works.

Without hairpin NAT, the client on the LAN cannot use the public name or the public IP that the rest of the Internet uses. The client must use the private IP or a split DNS name (later DNS topic).

Hairpin is a local feature. It is not required by IP. Cheap routers differ. Some implement it. Some do not. Some implement it only for some ports.

A clean design uses an internal name that resolves to the private IP on the LAN, and the public name on the Internet. Hairpin is a convenience when you cannot do that.

Do not test a port forward only from the LAN to the public IP and then claim that the Internet works. Test from a real outside path. Do not assume hairpin exists.

Source and dest NAT can both apply on a hairpin packet. The exact rewrite depends on the vendor. You only need the symptom and the idea.

### Questions

#### Theoretical questions

1. What is hairpin NAT?
2. Why does a LAN client use the public IP in this problem?
3. Why can the return path fail without a special rewrite?
4. What is a cleaner alternative to hairpin?
5. Why is a LAN test of the public IP a weak proof of an inbound service?

#### Easy practical tasks

1. Write five sentences that explain a LAN browser and a public URL of a home server.
2. Draw the U-turn: client, router, server, all on one LAN.
3. Make a table: split DNS, hairpin NAT, use private IP. Add one pro each.
4. Check your router manual or UI for "NAT loopback" or "hairpin". Write what you find.

#### Medium practical tasks

1. If you have a lab forward, try the public IP from the LAN and from mobile data. Write both results.
2. Write six sentences: why the server might see the client as the router address after hairpin.
3. Design a split-horizon idea in four sentences (names only). Full DNS is a later topic.

#### Advanced practical tasks

1. Write a one-page test plan for a home HTTPS service: outside path, LAN private name, LAN public name.
2. On a Linux NAT lab that you own, read public notes on `nft` hairpin or `net.ipv4.ip_forward` plus NAT. Write the idea, not a copy of a long wiki.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe an outbound web request and its reply through source NAT. Then describe an inbound request through port forwarding.
2. How do CGNAT and hairpin NAT create two different "the public IP does not work" stories?
3. Why do peer-to-peer programs need extra servers when both sides use NAT?
4. Which facts separate NAT, firewall, and proxy?
5. A teammate says NAT is enough security. Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: why NAT, masquerade, port forward, P2P pain, CGNAT, hairpin.
2. Record LAN IP, WAN IP of the router, and public lookup IP. Label each.
3. Draw a home network with NAT and one optional forwarded port. Do not publish real ports.
4. List three tests: outbound browse, inbound service, LAN-to-public-IP. Write what each test proves.

#### Medium practical tasks

1. Write a troubleshooting flow for "friends cannot connect to my game." Include NAT type, UPnP, CGNAT, and relay.
2. Build a table of five applications you use. Mark likely NAT behavior: outbound only, needs forward, needs traversal.
3. Compare IPv4 NAT home and IPv6 end-to-end in eight cells.

#### Advanced practical tasks

1. Write a glossary of 12 NAT terms. Each entry: term, one sentence, one symptom or command.
2. Design a lab with two virtual NATs and one public-like router. List pings and one TCP service test. Run only on equipment that you own.
