# 2. Layered Models

## Description

A layered model splits network work into levels. Each layer has a job. Each layer uses the layer below it and serves the layer above it. This topic shows the OSI seven-layer map and the TCP/IP four-layer map. You also learn encapsulation, headers, and the end-to-end principle.

Use one term for each concept. A layer is a level with a clear job. A header is the extra bytes that a layer adds. Encapsulation is the act of wrapping data in a header. Complete this topic before you study the physical layer and the link layer.

The models are maps. They help you talk and debug. They are not a law. Real systems can split or merge jobs.

---

## OSI 7 layers (as a map, not a law)

The OSI model has seven layers. The model is a teaching map. Vendors and RFCs do not follow every OSI name in daily work. You still use the map to place a problem.

The seven layers, from the bottom:

1. Physical: bits on a medium (voltage, light, radio)
2. Data link: frames on one link, local addresses
3. Network: packets across many hops, logical addresses
4. Transport: end-to-end byte stream or datagrams, ports
5. Session: dialogs and checkpoints (rare as a separate box in IP networks)
6. Presentation: syntax and encoding (often part of the application)
7. Application: the user protocol (HTTP, DNS, SMTP)

Use the map when you ask: is the cable bad, is the MAC wrong, is the IP wrong, is the port closed, or is the HTTP request wrong? That question list follows the layers from the bottom.

Do not treat OSI as a religion. Session and presentation often live in the application program. Ethernet and IP do not use OSI session packets. The value of OSI is a shared vocabulary.

When a person says "Layer 2" they often mean the link. When a person says "Layer 3" they often mean IP. When a person says "Layer 4" they often mean TCP or UDP. When a person says "Layer 7" they often mean the application protocol.

### Questions

#### Theoretical questions

1. What is the job of the physical layer?
2. What is the job of the network layer in the OSI map?
3. Why does this handbook call OSI a map and not a law?
4. What do people often mean by "Layer 2" and "Layer 3"?
5. Where do session and presentation often live in an IP network?

#### Easy practical tasks

1. Write the seven OSI layer names in order from the bottom.
2. Make a table: layer number, layer name, and one example technology or protocol.
3. Write five debug questions, one for each of layers 1 through 4 and one for layer 7.
4. Label a simple drawing: cable, switch, router, TCP port, and HTTP. Map each label to an OSI layer.

#### Medium practical tasks

1. For a failed web page, write which OSI layer you would check first, second, and third. Give one reason for each choice.
2. Find three "Layer 7" product names (for example a load balancer). Write what "Layer 7" means in that product page in your own words.
3. Compare OSI layer 4 and a TCP connection in six sentences.

#### Advanced practical tasks

1. Read a public OSI overview (ISO or a vendor primer). Write one page: which two layers map poorly to TCP/IP in daily work, and why.
2. Build a troubleshooting flowchart that starts at layer 1 and ends at layer 7. Use only yes or no checks that a beginner can run.

---

## TCP/IP 4 layers: link, internet, transport, application

The TCP/IP model that this handbook uses has four layers:

1. Link: network interface, Ethernet, Wi-Fi, frames, MAC addresses
2. Internet: IP, routing, ICMP
3. Transport: TCP, UDP, ports
4. Application: HTTP, DNS, TLS as used by applications, and other user protocols

Some texts add a physical layer under the link. This handbook keeps bits-on-the-medium with the physical and link topic. The four-layer list still matches how you read a packet: frame, IP, TCP or UDP, then the application bytes.

The internet layer delivers packets to a destination IP address. The transport layer delivers data to a port on that host. The application layer defines the meaning of the bytes.

IP can run on many link types. TCP and UDP can run on IPv4 and IPv6. HTTP can run on TCP (and later on QUIC). That independence is the point of layers.

Use TCP/IP names in lab notes. Use OSI numbers only when a person uses them. Do not mix "Layer 3" and "internet layer" in the same sentence without a map.

The kernel and the network card implement the lower layers. The program implements the application layer. Sockets are the API between the program and the transport layer.

### Questions

#### Theoretical questions

1. Name the four TCP/IP layers in order from the bottom.
2. What does the internet layer deliver?
3. What does the transport layer deliver?
4. Why can IP run on more than one link type?
5. Where do sockets sit in this model?

#### Easy practical tasks

1. Make a two-column table: OSI layer and TCP/IP layer. Map all seven OSI layers.
2. Write four sentences, one per TCP/IP layer, that state the job of that layer.
3. List three application protocols and the transport protocol that each one often uses.
4. Draw a stack: Ethernet, IPv4, TCP, HTTP. Label each TCP/IP layer.

#### Medium practical tasks

1. Open a capture of `ping`. Write which TCP/IP layers you see and which layer is missing.
2. Open a capture of `curl` to a web site. Name the link, internet, transport, and application pieces that you can see.
3. Write six sentences that explain why a program does not send raw Ethernet frames for a normal HTTP request.

#### Advanced practical tasks

1. Compare the four-layer TCP/IP map with a five-layer textbook map that adds physical. Write when each map helps a beginner.
2. Write a lab: capture one UDP DNS query and one TCP HTTP request. For each packet, fill a four-row layer sheet.

---

## Encapsulation and decapsulation

Encapsulation is the wrap step. A sender application gives bytes to TCP or UDP. The transport layer adds a transport header. IP adds an IP header. The link layer adds a frame header and often a trailer (FCS). The physical layer sends bits.

Decapsulation is the unwrap step. The receiver link layer checks the frame. IP checks the IP header and passes the payload to TCP or UDP. The transport layer checks the port and passes the payload to the application.

Each layer treats the bytes from the layer above as payload. That payload can contain more headers. A capture tool shows the nested headers.

On the reverse path, the same process runs in the other host. Intermediate hops often decapsulate only to the internet layer. A router reads the IP header, then encapsulates the packet in a new frame for the next link. The router does not need the HTTP bytes.

Do not add a new header that the next hop cannot parse. Do not assume that every hop sees the application payload. Many hops see only the outer headers.

Tunneling is encapsulation of a packet inside another packet. A VPN can put IP inside UDP inside IP. The idea is the same wrap. Later topics can use that idea. This topic only needs the basic wrap and unwrap.

### Questions

#### Theoretical questions

1. What is encapsulation?
2. What is decapsulation?
3. What does a layer treat as payload?
4. Why does a router often stop at the IP header?
5. What is tunneling in one sentence?

#### Easy practical tasks

1. Draw boxes: HTTP payload, TCP header, IP header, Ethernet header. Show the wrap order at the sender.
2. Write five sentences that describe decapsulation at the receiver.
3. In Wireshark, expand one frame. Write the order of the protocol tree from the top of the packet.
4. Make a list: which headers a home router must read to forward a packet to the provider.

#### Medium practical tasks

1. Capture one ICMP echo. Write which headers you see. State which layer added each header.
2. Capture one TCP segment to port 443. Measure the Ethernet length, the IP total length, and the TCP payload length if the tool shows them.
3. Draw a router hop: incoming frame, IP lookup, outgoing frame. Mark what is removed and what is added.

#### Advanced practical tasks

1. Explain a GRE or WireGuard packet in a public diagram (or a capture if you have a tunnel). Write the nested headers in order. Do not set up an unauthorized tunnel.
2. Write a one-page note: how encapsulation lets HTTP stay the same when the link changes from Ethernet to Wi-Fi.

---

## What each layer adds (header)

A header is a structured set of fields at the start of a protocol data unit. A trailer is extra bytes at the end. Ethernet often has a frame check sequence as a trailer.

Typical additions:

- Link (Ethernet): destination MAC, source MAC, Ethertype, then payload, then FCS
- Internet (IPv4): version, header length, total length, TTL, protocol, header checksum, source IP, destination IP
- Internet (IPv6): version, payload length, next header, hop limit, source IP, destination IP
- Transport (UDP): source port, destination port, length, checksum
- Transport (TCP): source port, destination port, sequence number, acknowledgment number, flags, window, checksum
- Application: method, URL path, headers, body (HTTP) or a DNS query message

Each header answers a question. Ethernet: which NIC on this link. IP: which host on the internet. TCP or UDP: which socket. HTTP: which resource and which method.

The protocol field in IPv4 (or the next-header field in IPv6) names the payload type. The Ethertype names the payload type of the frame. Ports name the application socket, not the payload syntax. The application header names the meaning of the rest.

Header size has a cost. A small payload with a large header wastes capacity. Later topics cover MTU and segmentation.

Do not invent field names. Use the names that the capture tool and the RFC use. Do not skip a header when you debug. A wrong MAC and a wrong IP are different faults.

### Questions

#### Theoretical questions

1. What is a header?
2. What question does the Ethernet header answer?
3. What question does the IP header answer?
4. What question do TCP ports answer?
5. What names the payload type of an IPv4 packet?

#### Easy practical tasks

1. Make a table: layer, header name, and two field names.
2. In Wireshark, open one IPv4 packet. Write the source IP, the destination IP, and the protocol number.
3. Write five sentences that explain why each layer adds its own header.
4. Draw an HTTP request inside TCP inside IP inside Ethernet. Write one field per header.

#### Medium practical tasks

1. Compare an ICMP packet and a UDP packet in a capture. List headers that they share and headers that they do not share.
2. Calculate on paper: Ethernet header 14 bytes, IPv4 header 20 bytes, UDP header 8 bytes, DNS query 40 bytes. What is the frame payload size before FCS?
3. Find the TTL or hop-limit field in a capture. Write the value and the layer that owns the field.

#### Advanced practical tasks

1. For one TCP packet, write every header field that you can name without opening a later topic as a full study. Mark fields that you will study in the TCP topic.
2. Write a one-page field map for a beginner: Ethertype, IP protocol, and TCP port. Explain how a host uses each field in order.

---

## End-to-end principle (high-level)

The end-to-end principle says: put a function at the ends when the function needs the knowledge of the application. The network in the middle should stay simple. Reliability, encryption, and application meaning often belong in the hosts.

Example: IP does not retransmit your file. TCP on the two hosts can retransmit. The routers in the middle forward packets. They do not store the full file.

Example: the network cannot know the full meaning of a bank transfer. The application and TLS on the ends protect the meaning.

The principle is high-level. Real networks add middleboxes: NAT, firewalls, load balancers. Those boxes break the simple picture. Later topics cover NAT and middleboxes. You still use the principle to ask: must this function sit in every router, or only in the two programs?

Functions that every hop needs can sit in the network. Example: a checksum on a single link, or a TTL so that a loop dies. Functions that only the two applications understand should sit at the ends.

Do not expect the Internet to keep your byte order and your retries for you unless your transport or your application does that work. Do not put application policy in every switch.

### Questions

#### Theoretical questions

1. What does the end-to-end principle say in one sentence?
2. Why does IP not retransmit a file?
3. Which functions often belong in the two hosts?
4. What is a middlebox in this section?
5. Why can a TTL sit in the network while a bank transfer cannot?

#### Easy practical tasks

1. Write four sentences that explain the file-and-TCP example.
2. Make a table: function, "ends" or "network", and one reason. Add retransmission, encryption, TTL, and HTTP method.
3. List three middleboxes that you know (NAT, firewall, proxy). Write one sentence each.
4. Draw two hosts and three routers. Mark where TCP reliability lives.

#### Medium practical tasks

1. Write six sentences: what still works if a router does not understand HTTP.
2. Find a public short text on the end-to-end argument (Saltzer, Reed, Clark summary). Write three facts in STE sentences.
3. Explain why NAT is a tension with the end-to-end principle. Use four sentences. Do not describe NAT types yet.

#### Advanced practical tasks

1. Write a one-page note: which functions in HTTPS live at the ends, and which functions live in the IP network.
2. Design a checklist for a new network feature: "must every hop implement this?" Apply the checklist to retransmission and to TTL.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you use OSI numbers and TCP/IP names together without mixing them in one sentence?
2. Describe encapsulation of an HTTP request from the application to the link. Name each header that you add.
3. A packet arrives at a host. Describe decapsulation until the application gets the bytes.
4. How does the end-to-end principle change where you put reliability?
5. Why does a four-layer map help a capture session more than a seven-layer law?

#### Easy practical tasks

1. Write a one-page cheat sheet: OSI seven names, TCP/IP four names, encapsulation, header, end-to-end.
2. Capture any one packet. Fill a form: link header present, IP header present, transport header present, application bytes visible or encrypted.
3. Draw two stacks side by side: OSI and TCP/IP. Draw map lines between them.
4. Explain in five sentences why this topic comes before Ethernet and IP details.

#### Medium practical tasks

1. Take a failed `ping` and a failed `curl`. Write which layer you suspect for each failure. Give one command that tests that layer.
2. Write a lab report template with four sections, one per TCP/IP layer. Use it on one live capture.
3. Compare a switch (link) and a router (internet) in a table of eight cells. Use only facts from this topic and topic 1.

#### Advanced practical tasks

1. Write a teaching script (one page) that uses one Wireshark packet to teach encapsulation. Do not paste a full packet dump of private data.
2. Argue in one page: a firewall that reads HTTP is a middlebox. State which layers it uses and how that sits with the end-to-end principle.
