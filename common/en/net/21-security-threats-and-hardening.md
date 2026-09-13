# 21. Security Threats and Hardening

## Description

A network path can lie, leak, or replay. This topic shows spoofing, sniffing, and replay at a high level, the idea of a SYN flood, TLS interception, least privilege on listen ports, untrusted networks, VPNs and WireGuard at a high level, and the zero-trust idea.

Use one term for each concept. A threat is a class of abuse. Hardening is a control that reduces harm. Complete this topic after TLS, middleboxes, and observability.

This handbook does not give attack steps, payloads, or exploit procedures. You learn what to fear and what to tighten on systems that you own. Practice with configuration reviews, listen-socket lists, and TLS verify. Do not test attacks on a network that you do not own.

---

## Spoofing, sniffing, replay

Spoofing means the sender puts a false identity in a field that a receiver might trust. A false source IP is one form. A false DNS answer is another. A false From header in mail is another. The field looks valid. The peer is wrong.

Sniffing means a party reads frames or packets that were not for its application. A shared radio, a mirrored port, or a compromised host can expose cleartext. TLS and SSH exist so that sniffing the path does not reveal the payload. Addresses and SNI can still leak. Topic 13 stated that.

Replay means a party sends a captured valid message again. A protocol without freshness (time, nonce, or a replay window) can accept the second copy as new. TLS and modern auth tokens include freshness. A naive custom UDP protocol often does not.

Hardening ideas:

- do not trust a source IP as a user identity
- do not send secrets on HTTP, clear FTP, or open MQTT
- use protocols that authenticate and that expire
- treat DNS and mail headers as claims until you check them

Do not capture on a shared LAN to read other people data. Topic 1 already set that rule. Do not build a "demo spoof" against a third party.

Integrity from TLS does not stop a replay at an upper layer if your API accepts the same signed request twice without a nonce. Name the layer.

### Questions

#### Theoretical questions

1. What does spoofing falsify?
2. What does sniffing expose on a cleartext path?
3. What does replay reuse?
4. Why is a source IP a weak user identity?
5. Why can TLS still leave an API open to replay?

#### Easy practical tasks

1. Write five sentences that define the three words.
2. Make a table: spoof, sniff, replay. Add one row for a control that helps.
3. List three cleartext protocols from earlier topics.
4. Draw: observer on a Wi-Fi hop, TLS payload hidden, addresses visible.

#### Medium practical tasks

1. Write six sentences: a custom UDP command without a nonce.
2. Review one of your lab services. Write whether it uses TLS and whether it trusts a source IP.
3. Read a public short note on replay in HTTP APIs (idempotency keys). Write four STE sentences. No exploit steps.

#### Advanced practical tasks

1. Write a one-page threat sheet for a student API: spoof, sniff, replay. Controls only.
2. In a local app that you own, add an expiring token or a nonce idea. Document the check. Do not target a remote system.

---

## SYN flood (idea)

A TCP listener keeps state for half-open connections after a SYN. Topic 9 said that SYN uses kernel memory. A SYN flood is a large number of SYNs that try to fill that state. Honest clients then fail to complete a handshake. The idea is resource exhaustion at the handshake, not a clever payload.

Modern stacks use SYN cookies or large backlogs so that a simple flood is harder. Firewalls and anycast scrubbing absorb volume. You do not need the packet craft. You need the idea: listen state is finite.

Hardening ideas:

- do not expose a debug listener to the Internet
- keep backlog and file-descriptor limits in your capacity plan
- use a provider DDoS control when you serve the public
- fail closed when you cannot accept more connections

Do not run a flood against any host, including a "lab" that shares a path with other people. Do not publish flood commands. Do not confuse a SYN flood with a full application flood (many complete HTTP requests). Both can hurt. The state sits at different layers.

A capture of a real incident can show many SYNs and few ACKs. That pattern is a hint. One SYN from you is a connect. Many SYNs from you are a policy violation.

### Questions

#### Theoretical questions

1. Which TCP state does a SYN flood try to exhaust?
2. Why does a half-open connection cost memory?
3. What is the idea of SYN cookies in one sentence?
4. How is an application-layer flood different at a high level?
5. Why must you not "try a flood" on a shared network?

#### Easy practical tasks

1. Write five sentences that explain the idea without steps.
2. Make a table: one client SYN, many SYNs. Add one row for honest user effect.
3. Write four sentences: listen backlog as a finite resource.
4. Draw a queue of half-open connections.

#### Medium practical tasks

1. On a machine that you own, read how to list TCP listen backlog or `somaxconn` (Linux) or the idea on Windows. Write the name of the setting. Do not tune production.
2. Write six sentences: expose a game server on the Internet versus bind to a VPN address.
3. Find a public vendor page on SYN cookies or flood protection. Write four STE sentences. No tool recipes.

#### Advanced practical tasks

1. Write a one-page capacity note: file descriptors, listen backlog, and a load balancer as a shield. Ideas only.
2. Design a tabletop: users cannot connect, SYN count is high, HTTP logs are empty. What do you check? No attack commands.

---

## TLS interception

TLS interception means a middlebox terminates TLS, then opens a second TLS session to the origin. The client sees a certificate that the interceptor signed. The browser trusts that certificate only if a local CA is in the trust store. Enterprises use this for policy and malware scan. A hostile hotspot can try the same idea. Topic 13 said: do not ignore a warning.

If the client does not trust the interceptor CA, the handshake fails. That failure is a success of verify. If a user clicks through, the interceptor can read HTTP. That is a failure of the human process.

`curl` and language runtimes have their own stores. A laptop can "work" in the browser and fail in Go or Python. Topic 20 named that mismatch.

Hardening ideas:

- keep verify on
- install a company CA only through a managed process
- pin or restrict extra CAs on high-value apps when the threat model needs it
- do not send private data through a random intercept

Do not build an interceptor for a network that you do not own. Do not disable verify to pass a hotel portal. Use a personal network or a known captive-portal flow that does not break TLS to your bank.

HTTPS on the origin is still required after intercept. The second hop must also be TLS if the path is not private.

### Questions

#### Theoretical questions

1. How many TLS sessions does a typical interceptor create?
2. Why must a local CA sit in the client store?
3. Why is a verify failure the correct result without that CA?
4. Why can a browser and `curl` disagree under intercept?
5. Why must the hop to the origin still use TLS on an untrusted path?

#### Easy practical tasks

1. Write five sentences that explain intercept as two TLS sessions.
2. Make a table: no intercept, enterprise intercept, hostile hotspot. Add one row for what the user should do.
3. Draw: client, interceptor, origin. Label two certificates.
4. Write four sentences: click-through as a control failure.

#### Medium practical tasks

1. Check whether your work or school device lists an extra CA (settings name only). Write yes, no, or "I must not open that UI."
2. Write six sentences: Go `x509` error versus a green padlock.
3. Compare `curl -v` issuer on a public site with the browser issuer. Write both. Redact a company CA name if policy requires it.

#### Advanced practical tasks

1. Write a one-page user guide: warning on public Wi-Fi versus expected company intercept. No intercept setup steps.
2. On a machine that you own, document which trust store `curl`, the browser, and one language use. Three short bullets.

---

## Least privilege on listening ports

A process should listen only on the addresses and ports that it needs. A database that binds `0.0.0.0:5432` on a laptop is a larger risk than the same database on `127.0.0.1`. Topic 16 stated that for Redis and Postgres. The idea is general.

Least privilege also means: no extra services, no debug ports on the WAN, no old protocol left "for a test." The host firewall is a second line. The first line is the bind address and the auth on the socket.

Containers and cloud security groups are bind rules at another layer. A "deny inbound" group does not help if you published `0.0.0.0` and then opened the group for the world.

Hardening ideas:

- list listeners with `ss` or `netstat` on a schedule
- bind admin ports to localhost or to a management VLAN
- require auth even on a "private" port
- drop unused packages that open ports

Do not expose SSH, RDP, or a database to `0.0.0.0` on a home PC that uses UPnP. Do not keep a default password on a listener. Do not confuse "I am behind NAT" with "nobody can reach me." Topic 7 covered peer-to-peer and CGNAT. Inbound can still work through a forward or a leak.

IPv6 can publish a global address even when you thought you were "NAT safe." Check both families.

### Questions

#### Theoretical questions

1. What does a bind to `0.0.0.0` mean?
2. Why is `127.0.0.1` a tighter listen for a lab database?
3. Why is a host firewall not enough if you also open a cloud group to the world?
4. Why is NAT not a full listen policy?
5. Why must you list IPv6 listeners too?

#### Easy practical tasks

1. Write five sentences that explain least privilege for sockets.
2. Run `ss -lnt` or `netstat -an`. Mark each listen as localhost, LAN, or all addresses.
3. Make a table: service, desired bind, actual bind. Add three rows from your machine.
4. Write four sentences: default password on a listener.

#### Medium practical tasks

1. Change a local lab service that you own from `0.0.0.0` to `127.0.0.1` if it is safe to do so. Write before and after. Skip production.
2. Write six sentences: UPnP plus a wide bind on a home PC.
3. Check IPv6 listen lines (`[::]`). Write one example if present.

#### Advanced practical tasks

1. Write a one-page listen policy for a laptop lab: app, DB, SSH. Include IPv4 and IPv6.
2. Build a weekly checklist: list sockets, list extra CAs, list scheduled tasks that open ports. No attack items.

---

## Untrusted networks

An untrusted network is a path that you do not control: cafe Wi-Fi, a guest VLAN, a conference LAN, some cellular paths. Hosts on that LAN can spoof, sniff cleartext, and send rogue RAs or rogue DHCP. Topic 15 and topic 19 named rogue discovery.

Assume:

- any host on the LAN can see multicast and some broadcast
- captive portals can intercept HTTP
- TLS warnings are serious
- mDNS names are not identities
- your laptop listeners are reachable on the LAN

Hardening ideas:

- use HTTPS and verified TLS
- use a system firewall on the client
- turn off file share and extra listeners when you travel
- prefer a known VPN to a trusted end (next section)
- do not use the cafe as a place to accept a new SSH host key without an out-of-band check

Do not disable the host firewall "so that AirDrop works" on a hostile LAN. Do not complete a bank login through a certificate warning. Do not plug an unknown USB NIC as a "free adapter" into a work laptop.

A home LAN can also be untrusted if guests and IoT share one flat segment. Split them.

### Questions

#### Theoretical questions

1. What makes a network untrusted in this section?
2. Why is mDNS a weak identity on cafe Wi-Fi?
3. Why is a TLS warning more serious on an untrusted LAN?
4. What can a rogue DHCP server change?
5. Why split guests and IoT from your work laptop?

#### Easy practical tasks

1. Write five sentences that list threats on cafe Wi-Fi.
2. Make a table: trusted home VLAN, cafe Wi-Fi. Add rows for file share and TLS warnings.
3. Write four rules you will follow on the next untrusted LAN.
4. Draw: laptop, cafe AP, other clients, Internet.

#### Medium practical tasks

1. On a trip profile that you own, list shares and listeners you would turn off. Write the list. Do not disable a work control that policy requires.
2. Write six sentences: captive portal HTTP versus TLS to your mail.
3. Compare a guest Wi-Fi design (isolated) with a flat home LAN. Write a short paragraph.

#### Advanced practical tasks

1. Write a one-page travel hardening sheet: firewall, listeners, VPN, certificate warnings. No attack steps.
2. Design a home split: work VLAN, guest VLAN, IoT VLAN. One sentence each for DHCP and discovery.

---

## VPNs and WireGuard (high-level)

A VPN extends a private network over an untrusted path. The client and the server (or two peers) encrypt and authenticate a tunnel. Inner packets use private addresses. Outer packets use the path addresses. Topic 7 mentioned site-to-site VPN as a WAN. This section is the user and the operator view.

WireGuard is a modern VPN protocol. It uses a small set of keys, UDP, and a simple interface (`wg0` on many Unix systems). Other families exist (IPsec, OpenVPN). This topic does not pick a vendor for you. It names the job: authenticated encryption of a tunnel, plus a routing decision (all traffic, or only some prefixes).

A split tunnel sends only some prefixes into the VPN. A full tunnel sends the default route into the VPN. Full tunnel helps on cafe Wi-Fi. Split tunnel can leak DNS or Web traffic to the local LAN if you mis-set the routes.

Hardening ideas:

- verify peer keys out of band
- keep the VPN software updated
- do not treat "I am on the VPN" as a reason to skip TLS to the app
- watch DNS: the resolver must match the trust model

Do not paste private keys. Do not run a random "free VPN" for work secrets. Do not disable a company VPN to bypass policy.

A VPN is not zero trust by itself. It is one secure channel. The next section names that limit.

### Questions

#### Theoretical questions

1. What does a VPN hide from the untrusted path?
2. What still sits in the outer header?
3. What is a split tunnel?
4. Why must you still use TLS to an HTTPS app over a VPN?
5. Why is a random free VPN a poor place for work secrets?

#### Easy practical tasks

1. Write five sentences that explain inner versus outer packets.
2. Make a table: full tunnel, split tunnel. Add one row for cafe use.
3. Write four sentences: WireGuard as "keys plus UDP plus a virtual interface" (idea).
4. Draw: laptop, cafe, VPN peer, office net.

#### Medium practical tasks

1. If you have a VPN that you may use, write whether it is full or split. Give one reason (default route or not). Skip if none.
2. Write six sentences: DNS on the VPN versus DNS on the cafe resolver.
3. Read a public WireGuard conceptual page. Write four STE sentences. No full config copy.

#### Advanced practical tasks

1. Write a one-page choice note: full tunnel on travel, split tunnel on a trusted office LAN.
2. On a lab pair of VMs that you own, bring up a tunnel or write why you used paper only. Never publish keys.

---

## Zero-trust idea (awareness)

Zero trust is a design idea: do not treat a network location as proof of a user or a device. A packet from the "office prefix" is not enough. Each request must authenticate the client and must authorize the action. The path can be the Internet. The control is identity, device health, and policy.

VPN plus a flat trusted LAN is the older model. Once you join, many ports trust you. Lateral movement is easy if one host is weak. Zero trust shrinks that implicit trust. East-west traffic also needs auth.

Awareness facts:

- zero trust is not a product name that you must buy
- you still need TLS, least privilege, and logs
- you still need a device that you can update
- DNS and time still matter

Do not turn off the office firewall because a slide said "zero trust." Do not skip MFA on a VPN and call it zero trust. Do not expose every microservice to the Internet without a policy engine that you understand.

This section is awareness. You can say: location is not identity. Later architecture topics can go deeper.

### Questions

#### Theoretical questions

1. What must you not treat as proof of identity?
2. What does each request still need?
3. How does a flat trusted LAN fail that idea?
4. Is zero trust a reason to skip TLS?
5. Why is MFA on a tunnel not the whole idea?

#### Easy practical tasks

1. Write five sentences that explain "location is not identity."
2. Make a table: castle-and-moat VPN, zero-trust idea. Add one row for east-west.
3. Write four sentences: a printer VLAN that still needs auth for admin.
4. Draw: user, Internet, app, policy check.

#### Medium practical tasks

1. Write six sentences: an intern laptop on the office LAN and a payroll API.
2. Read a public zero-trust explainer from a standards or vendor learning page. Write four STE sentences. No product pitch copy.
3. List three controls from this topic that still apply in a zero-trust story.

#### Advanced practical tasks

1. Write a one-page awareness note: VPN, least privilege, and zero trust as layers, not replacements.
2. Design on paper: one app, identity check, device check, no trust of source IP. High-level only.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a laptop on cafe Wi-Fi that uses a full-tunnel VPN and HTTPS. Use sniffing, TLS verify, and least privilege on local listeners.
2. How do spoofing, a SYN flood idea, and TLS interception each target a different resource (identity, memory, or confidentiality)?
3. Why does zero trust still need the hardening from listen ports and from clocks?
4. When is a VPN the right control, and when is it only one hop in a larger policy?
5. A teammate wants to "test security" with flood tools on the office LAN. Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: three classic threats, SYN idea, intercept, bind addresses, untrusted LAN, VPN, zero trust.
2. Save an `ss`/`netstat` listen list. Mark any bind that is wider than you need.
3. Draw one picture: untrusted LAN, VPN, verified TLS to an app, localhost database.
4. List actions this topic forbids (flood, click-through, paste keys, attack third parties).

#### Medium practical tasks

1. Write a hardening review of your own lab machine: listeners, time sync, TLS verify in `curl`, file shares.
2. Write a travel versus home checklist with ten checks. No exploit items.
3. Explain in ten steps how a beginner reacts to a certificate warning on public Wi-Fi.

#### Advanced practical tasks

1. Build a glossary of 16 terms from this topic. Each entry: term, one sentence, one control. No attack recipes.
2. Write a tabletop incident: intercept warning at a cafe, then a successful VPN. What do you tell the user? Checks only.
