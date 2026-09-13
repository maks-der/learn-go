# 1. Networks and Layers

## Description

A network lets two programs exchange data. The programs can run on one host or on two hosts. This topic shows why programs talk over a network. It also shows host, packet, LAN, WAN, and Internet. You learn IP address and port. You learn OSI as a map and the TCP/IP four layers. You learn encapsulation. You learn packet capture as a learning tool.

Complete this topic first. Complete this topic before you study links and Ethernet. Pair this path with `os.topics.md` when you study sockets.

Use one term for each concept. A host is a device that has a network address. A packet is a unit of data that the network forwards. A layer is a level with one job. Encapsulation is the act of wrapping data in a header. Do not mix a host with a switch that has no address in your notes. Do not mix a packet with a full file.

Practice on a machine that you own. Use `ping`, `ipconfig` or `ip addr`, and a capture tool such as Wireshark or `tcpdump`.

---

## Why programs talk over a network

A program on one host often needs data or a service on another host. A browser gets a web page. A game client sends a move. A phone sends a message. The network carries the bytes.

Local memory and local files are not enough when the data lives on a different machine. Two processes on the same host can use a local socket. Two processes on two hosts must use a network path. The two cases use the same idea. One program sends bytes. The other program receives bytes.

A network call has a cost. The call is slower than a local function call. The call can fail. The other host can be down. A link can drop packets. You must design for delay and for errors.

Programs talk over a network for share, scale, and location. Share: many clients use one server. Scale: you add hosts. Location: the user is not in the same room as the data.

Do not treat the network as a reliable local variable. The network is a path for packets. Packets can be late, lost, or reordered.

### Questions

#### Theoretical questions

1. Why does a program use a network instead of a local file?
2. Name two reasons a network call can fail.
3. How is a network call different from a local function call?
4. What does share mean in this section?
5. Why must you design for delay?

#### Easy practical tasks

1. List five programs on your computer that send data to another host. Write one sentence for each program.
2. Write four sentences that explain why a browser uses a network.
3. Make a two-column table: "Local action" and "Network action". Add four rows.
4. Time a local file open and a fetch of a small public web page. Write the two times.

#### Medium practical tasks

1. Draw a diagram of one program that talks to a remote service. Label the two hosts and the path.
2. Write six short sentences that compare a local function call and a network call. Cover speed, errors, and location.
3. Find one service that you use every day. Write which host starts the talk and which host answers.

#### Advanced practical tasks

1. Measure the time of 20 `ping` echoes to a local gateway and to a public host. Write the minimum, the maximum, and the average for each target.
2. Write a one-page note: what happens to a user when the remote host is down. Use one real product as the example.

---

## Host, packet, LAN, WAN, Internet

A host is an endpoint. A personal computer, a phone, a virtual machine, and a server are hosts when they have a network address. A printer or a camera can also be a host.

A packet is a block of bytes that the network forwards as one unit. The packet has a header and a payload. The header has addresses and control fields. The payload is the data that a higher layer sent.

A LAN is a local area network. Hosts on a LAN sit on a short path. A home Wi-Fi network is a LAN. An office Ethernet is a LAN.

A WAN is a wide area network. A WAN spans a long distance. A company link between two cities is a WAN. The path uses many hops.

The Internet is a network of networks. Many operators connect their networks. A packet can cross several networks from your host to a public server.

A path from one host to another host can use many links. Each hop is a device that receives a packet on one link and sends the packet on another link. A switch and a router are hop devices.

Use these terms with one meaning:

- host: endpoint with an address
- packet: unit of data on the internet layer in this handbook, unless a later topic uses frame
- LAN: local network
- WAN: long-distance network
- Internet: the public network of networks

Do not call a switch a host unless the switch has a management address and you talk about that address. Do not call a full file a packet. The file can split into many packets.

### Questions

#### Theoretical questions

1. What is a host?
2. What two parts does a packet have?
3. What is a LAN?
4. What is a WAN?
5. What is the Internet in this section?

#### Easy practical tasks

1. Write five sentences that define host, packet, LAN, WAN, and Internet.
2. Label your home network as LAN or WAN. Label the path to a public website as LAN, WAN, or both.
3. Make a table: term, one example. Add host, packet, LAN, WAN, Internet.
4. Draw two LANs and one Internet cloud between them.

#### Medium practical tasks

1. Run `ipconfig` or `ip addr`. Write one host address that you see and whether it looks local.
2. Write six sentences: a file of 10 MB and many packets.
3. Count the hops in a `traceroute` or `tracert` to a public host. Write how many hops you see.

#### Advanced practical tasks

1. Write a one-page map of one daily path: your host, your LAN, your ISP, the Internet, the server LAN.
2. Capture one packet to a public host. Write the source address and the destination address. Use a network that you own.

---

## IP address and port

An IP address names an interface on the internet layer. IPv4 uses 32 bits. IPv6 uses 128 bits. Topic 3 covers the text form. This section needs the idea: the IP address selects a host (or a host interface) on the path.

A port is a 16-bit number on the transport layer. TCP and UDP each have their own port space. The port selects a socket on that host. A browser often uses TCP port 443 for HTTPS. DNS often uses UDP port 53.

A program that waits for clients binds a known port. A client often uses an ephemeral source port. The kernel picks that port.

The 5-tuple identifies a TCP or UDP flow in many tools: protocol, source IP, source port, destination IP, destination port. Firewalls and NAT boxes use that tuple.

An IP address without a port is not enough for two programs on the same host. Two servers can share one IP and use two ports. Two clients can share one IP and use two source ports.

Do not type a port into a field that wants only an IP address. Do not say port 80 when you mean the HTTP protocol on that port. The port is a number. The protocol is the language.

Loopback addresses stay on the same host. IPv4 loopback is `127.0.0.1`. IPv6 loopback is `::1`. A packet to loopback does not leave the machine.

### Questions

#### Theoretical questions

1. What does an IP address select?
2. What does a port select?
3. Are TCP port 80 and UDP port 80 the same socket?
4. What is an ephemeral port?
5. What is the IPv4 loopback address?

#### Easy practical tasks

1. Write five sentences that explain IP address and port for a web browser.
2. Make a table: service, common port, TCP or UDP. Add HTTPS and DNS.
3. Draw a host with one IP and two listening ports.
4. Write the 5-tuple fields in a list.

#### Medium practical tasks

1. Run `ss -lnt` or `netstat -an`. Copy two listen lines. Write the address and the port.
2. Use `curl` to a public HTTPS site. Write the destination port that you expect.
3. Write six sentences: two programs on `127.0.0.1` with two ports.

#### Advanced practical tasks

1. Capture one TCP connection. Write the 5-tuple from the first data packet.
2. Write a one-page note: why a server needs a known port and why a client can use an ephemeral port.

---

## OSI as a map; TCP/IP four layers

A layered model splits network work into levels. Each layer has a job. Each layer uses the layer below it and serves the layer above it.

The OSI model has seven layers. The model is a teaching map. Vendors and RFCs do not follow every OSI name in daily work. You still use the map to place a problem.

The seven OSI layers, from the bottom:

1. Physical: bits on a medium (voltage, light, radio)
2. Data link: frames on one link, local addresses
3. Network: packets across many hops, logical addresses
4. Transport: end-to-end byte stream or datagrams, ports
5. Session: dialogs and checkpoints (rare as a separate box in IP networks)
6. Presentation: syntax and encoding (often part of the application)
7. Application: the user protocol (HTTP, DNS, SMTP)

The TCP/IP model that this path uses has four layers:

1. Link: physical plus data link (Ethernet, Wi-Fi)
2. Internet: IP
3. Transport: TCP, UDP
4. Application: HTTP, DNS, TLS as the program sees them

Use the map when you ask: is the cable bad, is the MAC wrong, is the IP wrong, is the port closed, or is the HTTP request wrong? That question list follows the layers from the bottom.

When a person says Layer 2 they often mean the link. When a person says Layer 3 they often mean IP. When a person says Layer 4 they often mean TCP or UDP. When a person says Layer 7 they often mean the application protocol.

Do not treat OSI as a law. Session and presentation often live in the application program. Ethernet and IP do not use OSI session packets. The value of OSI is a shared vocabulary.

### Questions

#### Theoretical questions

1. What is the job of the physical layer?
2. What is the job of the network layer in the OSI map?
3. Name the four TCP/IP layers in this handbook.
4. What do people often mean by Layer 2 and Layer 3?
5. Why does this handbook call OSI a map and not a law?

#### Easy practical tasks

1. Write the seven OSI layer names in order from the bottom.
2. Make a table: OSI layer, TCP/IP layer, one example protocol or medium.
3. Write five debug questions, one for each of layers 1 through 4 and one for layer 7.
4. Label a simple drawing: cable, switch, router, TCP port, and HTTP. Map each label to an OSI layer.

#### Medium practical tasks

1. Place `ping`, Ethernet, TCP, and HTTPS on both maps. Write one line each.
2. Write six sentences: session and presentation in a real browser.
3. Find a fault story from your own use (no Wi-Fi, wrong URL, timeout). Name the layer that you suspect first.

#### Advanced practical tasks

1. Write a one-page comparison of OSI seven layers and TCP/IP four layers. Use only facts from this section.
2. Read a public short OSI overview. Write four STE sentences about what the map is for. Do not copy long passages.

---

## Encapsulation

Encapsulation is the wrap of higher-layer data in a lower-layer header. The application writes a message. TCP or UDP adds a transport header. IP adds an internet header. Ethernet adds a frame header and a trailer.

On the receive path the stack removes headers in the reverse order. That act is decapsulation.

A header is the extra bytes that a layer adds. The payload of one layer is often the full PDU of the layer above. PDU means protocol data unit. This handbook uses frame on the link, packet on IP, and segment or datagram on transport.

Example path for an HTTP GET on Ethernet:

1. HTTP bytes
2. TCP header plus HTTP bytes
3. IP header plus that TCP segment
4. Ethernet header plus that IP packet plus a frame check

Each hop on a router path keeps the IP packet in the usual case. The hop builds a new link frame for the next link. The MAC addresses change. The IP addresses stay the same until NAT (topic 4).

Do not call the Ethernet unit a packet in lab notes when you mean a frame. Do not forget that encryption (TLS) sits above TCP in the common HTTPS case. The capture then shows TLS records, not HTTP text.

The end-to-end idea says: keep as much function as you can in the endpoints. The network forwards packets. Middleboxes exist. Topic 12 covers them. Encapsulation still starts at the host stack.

### Questions

#### Theoretical questions

1. What is encapsulation?
2. What is decapsulation?
3. What header sits around TCP in a normal IPv4 send?
4. What changes at a router hop: the IP addresses or the Ethernet addresses, in the usual case?
5. What name does this handbook use for the link unit?

#### Easy practical tasks

1. Write five sentences that walk HTTP into Ethernet for one GET.
2. Draw four boxes: HTTP, TCP, IP, Ethernet. Show wrap order.
3. Make a table: layer, unit name, address type if any.
4. List headers that you expect in a cleartext HTTP capture.

#### Medium practical tasks

1. Open a capture of any TCP session. Write which headers Wireshark shows for one packet.
2. Write six sentences: why a router rewrites the Ethernet header and keeps the IP header.
3. Draw send and receive. Mark encapsulation on the left and decapsulation on the right.

#### Advanced practical tasks

1. In a capture, measure header sizes for Ethernet, IP, and TCP on one packet. Write the three sizes and the payload size.
2. Write a one-page note: encapsulation versus NAT rewrite. How is a header change different from a wrap?

---

## Packet capture as a learning tool

A packet capture records frames or packets on an interface. Wireshark is a graphical tool. `tcpdump` is a command-line tool. You see headers that programs hide.

Capture is a learning tool. You confirm a handshake. You see a DNS query. You see a TLS Client Hello. You see a reset.

A capture filter selects packets before they enter the file. A display filter hides packets in the view after the file exists. Use a tight capture filter on a busy link. Example capture idea: ICMP only, or one host address.

Safety rules:

- capture on a network that you own or that you have permission to use
- do not publish files that contain cookies, tokens, or passwords
- stop the capture when you have the event
- prefer your own lab or a VM

A capture does not replace `ping` or `ss`. Start with a simple tool. Open a capture when you need bytes.

Do not capture on a shared office Wi-Fi to read other people traffic. Do not leave a capture running for hours on a laptop that visits many sites.

Topic 13 returns to Wireshark filters and to debug order. This section only teaches: capture shows the layers.

### Questions

#### Theoretical questions

1. What does a packet capture record?
2. What is the difference between a capture filter and a display filter?
3. Name two capture tools.
4. Why can a capture file be a secret leak?
5. When do you open a capture instead of `ping` only?

#### Easy practical tasks

1. Write five sentences that explain capture as a learning tool.
2. Start Wireshark or `tcpdump`, ping a host that you own, stop the capture. Write whether you see ICMP.
3. Make a table: safe capture, unsafe capture. Add two rows each.
4. Write one capture-filter idea for DNS and one for ICMP.

#### Medium practical tasks

1. Capture `ping` to your gateway. Write the source IP and the destination IP of one echo request.
2. Write six sentences: why you stop a capture after the event.
3. Open a saved capture. Count packets. Write the first protocol that you recognize.

#### Advanced practical tasks

1. Write a one-page lab: capture a TCP handshake to a host that you own. Name the three flags. Redact addresses if you share the note.
2. Compare Wireshark and `tcpdump` in a short table: UI, filter language, typical use.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A browser cannot load a page. How do you use the layer map to choose the first check?
2. How do host, IP address, and port work together for one TCP connection?
3. Why can two programs on one LAN still fail if encapsulation is correct but the port is closed?
4. What stays the same and what changes when a packet moves from your LAN to the Internet?
5. Why is a packet capture not enough if you never name the layer that you look at?

#### Easy practical tasks

1. Write a one-page cheat sheet: host, packet, LAN, WAN, Internet, IP, port, OSI, TCP/IP four layers, encapsulation, capture.
2. Run `ping` to `127.0.0.1` and to your gateway. Write one sentence for each result.
3. Draw one HTTPS request as four layers. Leave TLS as "application protection" for topic 10.
4. Bookmark a Wireshark filter page. Write when you open it.

#### Medium practical tasks

1. Write a script or a note that records your IPv4 address, your default gateway, and one `ping` RTT. Save a lab report.
2. Map one real fault from your week onto the four TCP/IP layers. Write the symptom and the layer that you checked first.
3. Capture loopback traffic if your OS allows it, or capture LAN ICMP. Write which layers you can see.

#### Advanced practical tasks

1. Write a one-page teaching talk: why programs talk over a network, then show one captured packet and name every header.
2. Compare a local socket talk on one host with a talk across the Internet. Use layers, delay, and failure modes.
