# 15. DHCP, Time, and LAN Services

## Description

LAN hosts need an address, a clock, and a way to find nearby services. This topic shows DHCP leases, NTP, mDNS and Bonjour, and a short view of other discovery protocols.

Use one term for each concept. A lease is a timed grant of configuration. NTP is a protocol that syncs clocks. mDNS is DNS-like name lookup on the local link. Complete this topic after IPv4 addressing and the IPv6 SLAAC preview.

Practice on a LAN that you own. Read `ipconfig /all` or `nmcli` / `ip addr`. Do not run a second DHCP server on an office LAN.

---

## DHCP lease

DHCP (IPv4) gives a host a temporary configuration: IPv4 address, subnet mask, default gateway, and often DNS resolvers. The grant is a lease. The lease has a duration. The host must renew the lease before the time ends if it wants to keep the address.

The common IPv4 dance: discover, offer, request, acknowledge. The new host can broadcast. A relay can forward to a DHCP server on another subnet. This topic needs the idea: the server is the source of the lease. The host does not pick a random LAN address in the usual office setup.

When the lease ends and renewal fails, the host must stop using that address. Some hosts keep a link-local IPv4 address (`169.254/16`) if DHCP fails. That address does not replace a gateway.

DHCPv6 can assign IPv6 addresses or only options. Topic 5 covered SLAAC versus DHCPv6. A home box can do both. Do not expect DHCPv4 to set IPv6.

Do not run two uncoordinated DHCP servers on one LAN. Do not shrink a pool so that guests get no address. Do not treat a lease as a security grant. Anyone who can use the LAN can often get a lease.

Reservations bind a MAC or a client identifier to a stable address. The grant is still a lease. The address can change if you change the reservation.

### Questions

#### Theoretical questions

1. What four items does a typical IPv4 DHCP lease include?
2. What must a host do before a lease ends if it wants to keep the address?
3. What IPv4 range do some hosts use when DHCP fails?
4. Does DHCPv4 configure IPv6 addresses?
5. Why is a lease not a security grant?

#### Easy practical tasks

1. Write five sentences that explain a lease as a timed address plus options.
2. Run `ipconfig /all` or the Unix equivalent. Write the IPv4 address, DHCP server, and lease times if shown.
3. Make a table: discover, offer, request, ack. Add one row for who sends.
4. Draw: host, LAN broadcast, DHCP server, lease.

#### Medium practical tasks

1. Find the lease expiry on your OS. Write how long the lease still lasts.
2. Write six sentences: two DHCP servers with overlapping pools.
3. Compare your DNS resolver from DHCP with a manual resolver that you do not apply. Write the DHCP value only.

#### Advanced practical tasks

1. On a lab LAN that you own, watch a renew or a reconnect. Capture DHCP if you can. Write the message names that the tool shows.
2. Write a one-page office DHCP plan: pool size, reservation for a printer, and what happens when the server is down.

---

## NTP

NTP is the Network Time Protocol. A host asks time servers and adjusts its clock. Correct time matters for TLS certificates, for logs, for Kerberos, and for signed tokens.

A large clock error looks like an expired certificate. Topic 13 named that fault. NTP (or SNTP, or an OS time service) is the usual fix. The OS often uses `time.windows.com`, `timed.apple.com`, or a pool such as `pool.ntp.org`. Enterprises use internal NTP.

NTP is not a human time zone. Time zone is a display offset. NTP sets UTC. Do not "fix TLS" by changing the time zone only.

A host with no NTP and a dead CMOS battery can boot with a year in the past. Then every TLS session fails. Check the clock before you debug PKI.

Do not point a production fleet at a random public NTP server that you may not use at that scale. Do not disable time sync to hide a test. Do not trust logs that have no synced clock when you compare two hosts.

Precision beyond one second is out of scope. You need "the year and the day are correct."

### Questions

#### Theoretical questions

1. What does NTP adjust?
2. Why does TLS need a correct clock?
3. Is a time zone the same as NTP time?
4. What can a dead battery plus no NTP do to a boot clock?
5. Why must you sync clocks before you compare logs on two hosts?

#### Easy practical tasks

1. Write five sentences that link NTP and certificate validity.
2. Check your OS time sync setting. Write whether it is on.
3. Make a table: UTC, time zone, NTP. Add one row for purpose.
4. Write your current clock and whether it matches a trusted phone clock.

#### Medium practical tasks

1. Find which time server your OS uses if the UI shows it. Write the name.
2. Write six sentences: a laptop that slept for a week and then failed TLS for one minute.
3. Compare `timedatectl` or Windows "Internet time" with a browser certificate date. Write both.

#### Advanced practical tasks

1. Write a one-page note: how you would confirm clock error before you revoke a certificate.
2. On a lab VM that you own, see how the OS lists NTP. Do not set a wrong date on a work machine.

---

## mDNS / Bonjour

mDNS is multicast DNS. Hosts send DNS-like queries to a multicast address on the local link. They do not need a unicast DNS server for those names. The common suffix is `.local`.

Bonjour is an Apple name for mDNS plus DNS-SD (service discovery). Printers, speakers, and laptops advertise names and services. Windows and Linux also implement mDNS in many builds.

A name such as `printer.local` is not a public Internet name. It must not leave the link in the usual design. A router should not forward that multicast to the WAN.

mDNS uses UDP. It is chatty on a quiet capture. Topic 5 already said that modern OSes show IPv6 and LAN services. mDNS is one of those services.

Do not use `.local` as a public corporate zone if you also want mDNS. The mix confuses resolvers. Do not trust an mDNS name as proof of identity on an untrusted LAN. Anyone on the link can answer.

Split horizon and mDNS are different. mDNS is link-scoped. Split horizon is a unicast DNS view.

### Questions

#### Theoretical questions

1. Where do mDNS queries go?
2. What suffix is common for mDNS names?
3. What is Bonjour in one sentence?
4. Should a home router forward mDNS to the Internet?
5. Why is an mDNS name a weak identity on public Wi-Fi?

#### Easy practical tasks

1. Write five sentences that compare mDNS and unicast DNS.
2. Make a table: `example.com`, `printer.local`. Add rows for resolver and scope.
3. List three device types that often use mDNS.
4. Draw: two hosts, multicast query, one answer.

#### Medium practical tasks

1. On a LAN that you own, look for `.local` names in a printer dialog or `dns-sd` / OS browse UI. Write one name if you find one.
2. Write six sentences: a `.local` corporate zone plus mDNS.
3. Capture a few packets with a multicast DNS filter if your tool has one (`mdns` or UDP 5353). Write the count. Home lab only.

#### Advanced practical tasks

1. Write a one-page note: when to use mDNS and when to use DHCP plus unicast DNS.
2. Design a guest VLAN: mDNS isolation yes or no. Give one reason each way.

---

## Discovery protocols (awareness)

Hosts also use other discovery protocols. This section names them. You do not configure all of them here.

SSDP / UPnP lets devices find media boxes and some routers. UPnP can open port mappings. That feature is a risk on a hostile LAN. Turn it off on a router if you do not need it.

LLMNR is a Windows name query on the local link. It is similar in spirit to mDNS. Some guides disable LLMNR for hardening. That is a policy choice.

WS-Discovery finds printers and some Windows services. ARP and IPv6 neighbor discovery find MAC addresses. They are not application service discovery, but they appear in the same captures.

DHCP options can point to a TFTP server or a VoIP controller. That is discovery through the lease, not multicast chat.

Do not enable every discovery feature on a guest network. Do not treat a discovered URL as trusted. Do not expose UPnP to the WAN.

Awareness means: you can name the packet in Wireshark and you know the scope is often the LAN.

### Questions

#### Theoretical questions

1. What risk does UPnP port mapping add?
2. What is LLMNR in one sentence?
3. How can DHCP act as discovery without multicast?
4. Why must a guest network limit discovery?
5. Is ARP a service-discovery protocol in the mDNS sense?

#### Easy practical tasks

1. Write five sentences that list SSDP, LLMNR, and mDNS as LAN discovery.
2. Make a table: protocol, typical use, scope. Add three rows.
3. Find a UPnP or "user-friendly" setting on a home router UI if you may open it. Write the name of the setting. Do not change a router that you do not own.
4. Write four sentences: discovered printer versus a typed IP.

#### Medium practical tasks

1. In a one-minute LAN capture that you own, write which discovery names you see (mDNS, SSDP, LLMNR, ND).
2. Write six sentences: UPnP on a home router facing the Internet.
3. Read a public short note on LLMNR hardening. Write four STE sentences. No attack steps.

#### Advanced practical tasks

1. Write a one-page awareness sheet: DHCP, NTP, mDNS, SSDP, LLMNR. One sentence each plus one "do not".
2. Design a home lab filter list in Wireshark for this topic. Write the filter texts. Do not include attack recipes.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a PC that joins a home LAN. Use DHCP, a clock source, and one discovery protocol.
2. How do DHCP DNS options and mDNS both answer "what is the address of this name" in different scopes?
3. Why can a correct DHCP lease still fail TLS if NTP is wrong?
4. When is a second DHCP server a lab tool, and when is it an outage?
5. A teammate says "the LAN is trusted, so discovery can stay open." Which facts do you use in a reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: lease fields, NTP versus time zone, `.local`, one other discovery name.
2. Save `ipconfig /all` or `ip addr` plus a time-sync screenshot name. Label DHCP and clock.
3. Draw one picture: router as DHCP and NTP client of the ISP, hosts as DHCP clients.
4. List ports or multicast names that this topic mentioned (67/68 idea, 5353, NTP).

#### Medium practical tasks

1. From your machine, identify: lease times, DNS from DHCP, time sync on or off, one LAN discovery sign.
2. Write a lab plan for a guest Wi-Fi: DHCP on, mDNS off or isolated, UPnP off. State why.
3. Explain in ten steps how a beginner checks "I got an address" before they debug a website.

#### Advanced practical tasks

1. Build a glossary of 12 terms from this topic. Each entry: term, one sentence, one command or UI.
2. Write a short incident note: "printer disappeared from the laptop." Include DHCP, mDNS, and VLAN as possible causes. No solutions beyond checks to run.
