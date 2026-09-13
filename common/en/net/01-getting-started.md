# 1. Getting Started

## Description

A network lets two programs exchange data. The programs can run on one computer or on two computers. This topic shows the basic words: host, packet, link, LAN, WAN, Internet, client, server, IP address, and port. You also learn why a packet capture helps you.

Use one term for each concept. A host is a computer or a device that has a network address. A packet is a unit of data on the network. A link is a connection between two devices. Complete this topic before you study layered models.

Practice on a real machine. Use `ping`, `ipconfig` or `ip addr`, and a capture tool such as Wireshark or `tcpdump`. Pair this path with the operating-system topics when you study sockets and the kernel.

---

## Why programs talk over a network

A program on one host often needs data or a service on another host. Examples: a browser gets a web page. A game client sends a move. A phone sends a message. The network carries the bytes.

Local memory and local files are not enough when the data lives on a different machine. A process on the same machine can use a local socket. A process on a remote machine must use a network path. The two cases use the same idea: one program sends bytes, and the other program receives bytes.

A network call has a cost. The call is slower than a local function call. The call can fail. The other host can be down. A link can drop packets. You must design for delay and for errors.

Programs talk over a network for share, scale, and location. Share: many clients use one server. Scale: you add hosts. Location: the user is not in the same room as the data.

Do not treat the network as a reliable local variable. The network is a path for packets. Packets can be late, lost, or reordered.

### Questions

#### Theoretical questions

1. Why does a program use a network instead of a local file?
2. Name two reasons a network call can fail.
3. How is a network call different from a local function call?
4. What does "share" mean in this section?
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

## Host, packet, link

A host is an endpoint. A personal computer, a phone, a virtual machine, and a server are hosts when they have a network address. A printer or a camera can also be a host.

A packet is a block of bytes that the network forwards as one unit. The packet has a header and a payload. The header has addresses and control fields. The payload is the data that a higher layer sent.

A link is the connection that carries bits between two devices. A cable, a fiber, or a radio channel is a link. A host attaches to a link through a network interface.

A path from one host to another host can use many links. Each hop is a device that receives a packet on one link and sends the packet on another link. A switch and a router are hop devices.

Use these terms with one meaning:

- host: endpoint with an address
- packet: unit of data on the internet layer and above this handbook, unless a later topic uses "frame"
- link: one hop connection
- interface: the attachment of a host to a link

Do not call a switch a host unless the switch has a management address and you talk about that address. Do not call a full file a packet. The file can split into many packets.

### Questions

#### Theoretical questions

1. What is a host?
2. What two parts does a packet have?
3. What is a link?
4. What is a hop?
5. Why is a file not the same as one packet?

#### Easy practical tasks

1. Run `ipconfig` on Windows or `ip addr` on Linux or macOS. Write the name of one interface and one address.
2. Draw one host, one link, and a second host. Label each part.
3. Write five sentences that use host, packet, link, header, and payload. Use each word one time.
4. List three devices in your room that can be hosts.

#### Medium practical tasks

1. Count the interfaces on your machine. For each interface, write: up or down, IPv4 address if any, IPv6 address if any.
2. Use `ping` to your default gateway. Write the target address and the number of replies.
3. Draw a path with three links and two hops between two hosts. Label each link.

#### Advanced practical tasks

1. Capture one ICMP echo request and one reply with Wireshark or `tcpdump`. Write the length of each packet and the two addresses.
2. Write a short report that explains the difference between a host, a hop, and an interface. Use your home network as the example.

---

## LAN vs WAN vs the Internet

A LAN is a local area network. The LAN covers a home, a room, or a small office. Devices on one LAN share a local broadcast domain or a small set of local segments. Delay is low. You often own the equipment.

A WAN is a wide area network. The WAN covers a city, a country, or many sites of one organization. A leased line, an MPLS service, or a site-to-site VPN can form a WAN. Delay is higher than on a LAN.

The Internet is a public network of networks. Many organizations connect their networks. They use common addresses and common routing. You do not own the full path. You own only your edge.

A home LAN often has one router that attaches to an Internet service provider. Hosts on the LAN use private addresses. The router forwards packets to the WAN link of the provider. Later topics cover NAT and routing.

Rules:

- use LAN for the local network that you control
- use WAN for a long-distance private or provider network
- use Internet for the public internetwork
- do not say "the Internet" when you mean only your LAN

A campus network can have many LANs that one organization owns. Those LANs still are not the Internet.

### Questions

#### Theoretical questions

1. What does LAN mean?
2. What does WAN mean?
3. What is the Internet in one sentence?
4. Who owns the full path of a packet on the Internet?
5. Why is a campus with many LANs still not the Internet?

#### Easy practical tasks

1. Write three examples of a LAN and three examples of a WAN.
2. Draw your home: LAN hosts, one router, and the provider cloud. Label LAN and Internet.
3. Make a table with columns LAN, WAN, and Internet. Add one row for delay, one row for ownership, and one row for example.
4. Write four sentences that explain how a phone on home Wi-Fi reaches a public website.

#### Medium practical tasks

1. Find the public IPv4 address of your home router (search "what is my IP" in a browser, or use a trusted lookup). Compare it with an address from `ipconfig` or `ip addr`. Write which address is on the LAN.
2. Measure `ping` to a LAN host and to a public host. Write both round-trip times.
3. List the networks that a packet can cross from your laptop to `example.com`. Use LAN, provider, and Internet as labels.

#### Advanced practical tasks

1. Traceroute to a public host (`tracert` on Windows, `traceroute` on Unix). Mark which hops look local and which hops look remote. Write how you decided.
2. Write a one-page comparison of a home LAN and a company WAN. Cover ownership, delay, and who can change the path.

---

## Client and server

A client is a program that starts a request. A server is a program that waits for a request and then answers. A browser is a client for HTTP. A web process that listens on a port is a server.

The same host can run a client and a server. The same program can be a client in one talk and a server in another talk. The words describe the role in one exchange, not the hardware.

A server listens. The operating system keeps a socket in the listen state. A client connects or sends a datagram to that socket. Later topics cover TCP, UDP, and sockets.

Many services use one well-known port on the server. The client uses an ephemeral port. The pair of addresses and ports identifies the conversation.

Do not assume that the server is a large machine. A laptop can be a server for a local test. Do not assume that the client is a human. A cron job can be a client.

Peer-to-peer programs can have two roles at the same time. Each side can start a request. This topic still uses client and server for the role in one request.

### Questions

#### Theoretical questions

1. What does a client do?
2. What does a server do?
3. Can one host run both a client and a server? Explain.
4. What does "listen" mean for a server?
5. Why are client and server roles, not machine types?

#### Easy practical tasks

1. Name four client programs and the service that each one uses.
2. Write five sentences that describe a browser as a client and a web process as a server.
3. Draw a client, a server, a request arrow, and a response arrow. Label each part.
4. Find one server process on your machine (`ss -lnt` or `netstat -an`). Write the address and the port if you find one.

#### Medium practical tasks

1. Start a tiny local HTTP server (Python `http.server`, or any simple server). Fetch a file with a browser. Write which program is the client.
2. Use `curl` to request `https://example.com`. Write which program is the client in that command.
3. Write six sentences that compare a human using a browser and a script using `curl`. Cover the client role.

#### Advanced practical tasks

1. Run a local server. Connect from a second device on the LAN. Record the client address that the server log shows.
2. Write a short note on peer-to-peer: when both sides start requests. Give one real example (file share or game).

---

## IP address and port as a mailbox number

An IP address identifies a host interface on an internet. IPv4 uses 32 bits. IPv6 uses 128 bits. Later topics cover the formats.

A port is a 16-bit number on one host. The port selects a process or a socket on that host. The IP address is the building. The port is the mailbox number in that building.

A socket address is the pair of IP address and port, plus the protocol (TCP or UDP). Two conversations can share an IP address when the ports differ. Two conversations can share a port on one host when the remote address or the protocol differs. The full 5-tuple for TCP or UDP is: protocol, local IP, local port, remote IP, remote port.

Well-known ports are the low numbers that common services use. Example: TCP port 80 for HTTP, TCP port 443 for HTTPS, UDP port 53 for DNS. A client often uses a high ephemeral port that the operating system assigns.

Do not confuse an IP address with a MAC address. The MAC address is a link-layer address. Do not treat a port as a physical hole. The port is a number in the transport header.

The mailbox image helps a beginner. Then use the exact terms: IP address, port, and 5-tuple.

### Questions

#### Theoretical questions

1. What does an IP address identify?
2. What does a port select on a host?
3. What is a socket address?
4. What is an ephemeral port?
5. Why is a MAC address not an IP address?

#### Easy practical tasks

1. Write your IPv4 address and one port that a browser uses for HTTPS.
2. Make a table: service name, protocol, and well-known port. Add HTTP, HTTPS, and DNS.
3. Explain the mailbox image in four sentences. Then replace the image with the exact terms.
4. Run `ss -lnt` or `netstat -an`. Copy five listen lines. Mark the port in each line.

#### Medium practical tasks

1. Use `curl -v` on `https://example.com`. From the output, write the remote IP and the remote port.
2. Open two browser tabs to the same site. With a capture or `ss` / `netstat`, show that the client ports can differ.
3. Write the 5-tuple for one TCP connection on your machine. Name each of the five fields.

#### Advanced practical tasks

1. Capture one TCP handshake to a web server. Write the client port, the server port, and the two IP addresses.
2. Write a one-page note: why two processes on one host can both use the network at the same time. Use ports in the explanation.

---

## Packet capture as a learning tool

A packet capture records the bytes that an interface sends and receives. Wireshark is a graphical tool. `tcpdump` is a command-line tool. The capture shows headers and payloads that the host can see.

Capture is a learning tool. You see the real packets, not only a diagram. You see the handshake, the addresses, the ports, and the payload size. You can compare a textbook field with a live field.

Rules for safe capture:

- capture on a network that you own or that you have permission to use
- do not capture passwords or private data on a shared network
- do not publish a capture that contains secrets
- use a filter so that you keep only the packets that you need

A display filter in Wireshark hides packets in the view. A capture filter in `tcpdump` or in Wireshark capture options drops packets before the file grows.

Start with ICMP (`ping`) and with one HTTP or HTTPS session. Name each header that you see. Later topics tell you what each header means.

The capture shows only what this interface sees. Encrypted payloads (TLS) hide application data. You still see addresses, ports, and TLS records.

### Questions

#### Theoretical questions

1. What does a packet capture record?
2. What is the difference between Wireshark and `tcpdump`?
3. What is a capture filter?
4. What is a display filter?
5. Why can a capture miss packets that other hosts send?

#### Easy practical tasks

1. Install Wireshark or confirm that `tcpdump` exists. Write the version.
2. Start a capture. Run `ping` to your gateway. Stop the capture. Count the ICMP packets.
3. Write five rules for safe capture in your own words.
4. Open a sample capture if the tool includes one. Write three column names that you see.

#### Medium practical tasks

1. Capture `ping` to a public host. Write the source IP, the destination IP, and the ICMP type for the request.
2. Use a display filter for ICMP only. Write the filter text that you used.
3. Capture one `curl` to `http://example.com` (HTTP, not HTTPS, if the host still answers). Write the TCP ports.

#### Advanced practical tasks

1. Capture a TLS session with `curl -v https://example.com`. Write what you can read (addresses, ports) and what you cannot read (HTTP body).
2. Write a short lab guide for a classmate: how to capture `ping`, how to save the file, and how to apply one filter. Do not include solutions to later topics.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of one request from a client program to a server program. Use host, packet, link, IP address, and port.
2. Why do LAN, WAN, and the Internet need three different words?
3. How do client, server, IP address, and port work together in one TCP conversation?
4. When is packet capture the wrong first tool, and when is it the right tool?
5. A teammate says "the server is the big computer." Which facts do you use to correct that sentence?

#### Easy practical tasks

1. Write a one-page cheat sheet with these words: host, packet, link, LAN, WAN, Internet, client, server, IP address, port, capture.
2. Run `ipconfig` or `ip addr`, then `ping` to the gateway. Save both outputs in one text file. Label each block.
3. Draw one picture that includes a LAN client, a home router, and a public server. Label every term from this topic that appears.
4. List the tools that this learning path asks you to use (`ping`, `dig`, `curl`, Wireshark or `tcpdump`). Confirm which tools you have installed.

#### Medium practical tasks

1. From your machine, identify: one LAN address, one public address (lookup), one gateway, and one listening port. Write how you found each value.
2. Write a lab plan for the next topic (layered models) that uses one capture. State what you will look for. Do not explain headers yet.
3. Explain in ten steps how a beginner should capture `ping` on a home LAN without saving private payloads from other programs.

#### Advanced practical tasks

1. Build a small glossary of 15 terms from this topic. Each entry: term, one-sentence definition, and one command or tool that shows it.
2. Capture one full `curl` to a public HTTPS site. Write a timeline: DNS if visible, TCP if visible, TLS if visible. Mark each step as "I can name this now" or "later topic".
