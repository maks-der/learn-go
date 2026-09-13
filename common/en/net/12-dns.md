# 12. DNS

## Description

DNS maps a name to data. The common case is a name to an IP address. This topic shows resolvers, record types, TTL, `dig` and `nslookup`, split horizon, and DNS over HTTPS or TLS.

Use one term for each concept. A name is a domain name such as `example.com`. A resolver is a program that asks DNS questions. An authoritative server answers from a zone that it hosts. Complete this topic after HTTP. You already use names in URLs.

Practice with `dig` or `nslookup`. Capture UDP port 53 when the query is cleartext. Encrypted DNS hides the query from the path. You still see the resolver address.

---

## Name to address

A program often has a name, not an IP address. The host must find an address before it can send IP packets to that peer. DNS is the usual map from name to address.

The operating system has a stub resolver. The stub sends a query to a recursive resolver. The recursive resolver is often the home router, an ISP resolver, or a public resolver. The answer comes back as one or more IP addresses.

A name can have an IPv4 address (A), an IPv6 address (AAAA), both, or neither. Dual stack uses both when both exist. A later connection uses one address. Happy Eyeballs can try both families.

The map is not only "name to address." Mail uses MX. Services use SRV. Text uses TXT. This section starts with address lookup because browsers and `curl` do that first.

Do not treat a name as a host. One name can point to many addresses. One address can serve many names. Do not cache an address in a long-lived program without a TTL plan. The next section on TTL explains why.

A local file such as `hosts` can override DNS for one name. Use that file only on a machine that you own and only for a test.

### Questions

#### Theoretical questions

1. What does a name-to-address lookup return in the common case?
2. What is a stub resolver?
3. Where does the stub send the query?
4. Can one name have both A and AAAA data?
5. Why is a name not the same thing as one host?

#### Easy practical tasks

1. Write five sentences that explain a browser name lookup before TCP.
2. Resolve `example.com` with `nslookup` or `dig`. Write one IPv4 address if you get one.
3. Make a table: name, A, AAAA. Add `example.com` and your own machine name if DNS has it.
4. Draw: program, stub, recursive resolver, answer.

#### Medium practical tasks

1. Compare `ping example.com` with a ping to the numeric address that DNS returned. Write both targets.
2. Write six sentences: a name with two A records. Who chooses the address?
3. Look at `hosts` on your OS (read only). Write whether any name overrides DNS. Do not publish private names.

#### Advanced practical tasks

1. Capture one cleartext DNS query for a name that you choose. Write the query name and the answer addresses. Use a network that you own.
2. Write a one-page note: name, address, and virtual host. How can many HTTPS names share one address?

---

## Recursive vs authoritative resolver

A recursive resolver accepts a query from a stub. It finds the answer. It can use its cache. If the cache has no valid data, it asks other servers. It starts at a root hint, then a TLD server, then the authoritative server for the zone. Then it returns the answer to the stub.

An authoritative server answers from zone data that an administrator or a registry put there. It does not need to recurse for names in that zone. A server can be authoritative for some zones and recursive for other queries. Many public resolvers recurse and are not authoritative for `example.com`.

The chain of authority uses NS records. The parent zone lists the name servers of the child. You follow NS from the root to the zone. This handbook does not require you to run a root server. You must know that recursion walks that tree.

A recursive resolver is a cache and a walker. An authoritative server is the source for a zone. Do not mix the two words.

A stub does not walk the tree in the usual OS setup. The stub asks one configured resolver. You can walk the tree yourself with `dig +trace` when the tool supports it.

Do not open a recursive resolver to the public Internet on a host that you do not intend as a public service. An open resolver can be abused. Do not change production DNS without a rollback plan.

### Questions

#### Theoretical questions

1. What does a recursive resolver do for a stub?
2. What does an authoritative server use as its source of data?
3. Where does recursion start when the cache is empty?
4. What record type points to the name servers of a zone?
5. Does a typical OS stub walk the root tree itself?

#### Easy practical tasks

1. Write five sentences that compare recursive and authoritative roles.
2. Find the DNS server addresses on your host (`ipconfig /all` or `/etc/resolv.conf`). Write them.
3. Make a table: stub, recursive, authoritative. Add one row for "who asks" and one row for "who stores the zone".
4. Draw the path: stub to recursive to authoritative (idea).

#### Medium practical tasks

1. Use `dig example.com NS` or `nslookup -type=NS example.com`. Write the NS names.
2. Write six sentences: your home router as a recursive forwarder to the ISP.
3. Run `dig +trace example.com` if you have `dig`. Write the first two steps that you see. Skip if the tool is missing.

#### Advanced practical tasks

1. Write a one-page map of a lookup for `www.example.com` from root to the A or AAAA answer. Use public NS data.
2. On a lab VM that you own, install a recursive resolver or use a container image. Query it from that VM. Document the listen address. Do not expose it to the Internet.

---

## Record types: A, AAAA, CNAME, MX, NS, TXT, SRV

DNS stores resource records. Each record has a type, a name, a class (usually IN), a TTL, and data.

A maps a name to an IPv4 address. AAAA maps a name to an IPv6 address. A client that wants IPv4 asks for A. A client that wants IPv6 asks for AAAA.

CNAME maps a name to another name. The resolver then looks up that target. Do not put other record types on the same owner name as a CNAME in the usual zone rules. A CNAME is an alias, not an address.

MX maps a domain to a mail exchanger name and a preference number. The mail system then looks up A or AAAA for that exchanger. A smaller preference value is tried first.

NS maps a zone to a name server name. The parent and the child both publish NS data. Glue A or AAAA records at the parent help when the NS name is inside the child zone.

TXT holds text. Operators use TXT for SPF, domain verification, and other strings. Read the format before you trust a TXT value.

SRV maps a service to a host, a port, and a priority. The name often looks like `_service._proto.example.com`. Some systems use SRV. Many web clients do not.

Do not invent a record type. Use the type that the protocol documents. Do not treat CNAME as a redirect at HTTP. HTTP redirects are a different layer.

### Questions

#### Theoretical questions

1. What data does an A record hold?
2. What data does an AAAA record hold?
3. What does a CNAME point to?
4. Why does MX data include a preference number?
5. Why is a CNAME not an HTTP redirect?

#### Easy practical tasks

1. Write one sentence for each type: A, AAAA, CNAME, MX, NS, TXT, SRV.
2. Query A and AAAA for `example.com`. Write the answers.
3. Query MX for a domain that publishes MX (a public mail domain). Write the exchanger names.
4. Make a table: type, purpose. Add all seven types from this section.

#### Medium practical tasks

1. Query NS and TXT for `example.com`. Write one NS name and whether TXT exists.
2. Write six sentences: a web name that is a CNAME to a CDN name. What lookups follow?
3. Find whether a product you use documents an SRV name. Write the name or write that you found none.

#### Advanced practical tasks

1. Build a glossary of the seven types. Each entry: type, data fields, and one command that shows it.
2. Write a one-page zone sketch on paper for a lab name: SOA idea (one line), NS, A, AAAA, MX. Do not claim it is a full RFC zone.

---

## TTL

TTL is time to live for a DNS record in a cache. The value is in seconds. A recursive resolver must not serve a cached record after the TTL ends. The stub and applications can also cache.

A long TTL reduces queries to authoritative servers. A short TTL lets you change an address faster. After you change a record, some caches still hold the old data until their TTL ends. Plan changes with the old TTL in mind.

TTL 0 means "do not cache" in the usual reading. Some caches still keep the data for a short extra time. Do not depend on instant global change.

Negative answers (no such name, or no records of that type) also have a cache time. That time comes from the SOA of the zone, not from a missing A TTL. A wrong name can "stick" as NXDOMAIN for minutes.

Do not set a one-second TTL on a busy public name without capacity for the query load. Do not ignore TTL in a long-running client. Refresh the name or use a resolver API that honors TTL.

Applications sometimes ignore TTL and cache for the process life. That behavior causes "it works after I restart."

### Questions

#### Theoretical questions

1. What unit does TTL use?
2. What must a recursive resolver do when a TTL ends?
3. Why does a long TTL slow a planned address change?
4. What is a negative cache in one sentence?
5. Why can a process restart fix a stale name?

#### Easy practical tasks

1. Query `example.com` A. Write the TTL that the tool prints.
2. Write five sentences that explain TTL to a teammate.
3. Make a table: TTL 30, TTL 3600, TTL 86400. Add one row for "change speed" and one row for "query load".
4. Convert 300 seconds and 86400 seconds to minutes or hours.

#### Medium practical tasks

1. Query the same name two times in one minute. Write whether the TTL went down (cache) or stayed the same.
2. Write six sentences: you change an A record while the old TTL is 3600. When can a remote user still see the old address?
3. Find the TTL of an MX or NS record for a public domain. Compare it with the A TTL.

#### Advanced practical tasks

1. Write a one-page change plan: lower TTL, wait, change address, then raise TTL. Include a rollback line.
2. Write a small program that resolves a name and prints the addresses. State whether your library exposes TTL. Record what you found.

---

## `dig` / `nslookup`

`dig` is a DNS query tool that is common on Unix and available for Windows. `nslookup` is available on Windows and on many Unix systems. Both send queries and print answers. Prefer `dig` when you have it. The output is clearer for record type and TTL.

Useful `dig` forms:

- `dig example.com A`
- `dig example.com AAAA`
- `dig example.com MX`
- `dig @8.8.8.8 example.com A` (ask a specific resolver)
- `dig +short example.com A` (short output)
- `dig +trace example.com` (walk from the root when policy allows)

Useful `nslookup` forms:

- `nslookup example.com`
- `nslookup -type=MX example.com`
- `nslookup example.com 1.1.1.1` (server as the last argument on many systems)

Read the answer section, not only the first line. Read the status: NOERROR, NXDOMAIN, SERVFAIL, REFUSED. Read the flags: `rd` is recursion desired. `ra` is recursion available.

Do not use these tools to flood a name server. Do not query a private corporate server from a network that has no permission. Use public names and resolvers that you may use.

`Resolve-DnsName` is a Windows PowerShell cmdlet. It can show record types in a table. Use it if you prefer PowerShell.

### Questions

#### Theoretical questions

1. What does `dig example.com A` ask for?
2. What does `@` mean in a `dig` command?
3. What does NXDOMAIN mean?
4. What is the difference between `rd` and `ra` at a high level?
5. Why prefer `dig` when both tools exist?

#### Easy practical tasks

1. Run `nslookup example.com` or `dig example.com`. Save the output. Mark the answer addresses.
2. Write five sentences that explain `+short` versus full `dig` output.
3. Make a table: status name, meaning. Add NOERROR, NXDOMAIN, SERVFAIL.
4. Query a name that does not exist. Write the status.

#### Medium practical tasks

1. Query the same name at your default resolver and at a public resolver (`@` or the `nslookup` server argument). Write whether the answers match.
2. Use `-type` or a type argument for TXT and NS. Write one line from each answer.
3. Write six sentences: `REFUSED` versus `SERVFAIL` as a beginner reading.

#### Advanced practical tasks

1. Write a one-page field guide for `dig` or `nslookup` for this learning path. Include A, AAAA, MX, NS, and a specific server.
2. Script three queries (A, AAAA, NS) and print a small table. Use a public name.

---

## Split horizon

Split horizon (or split view) means one name can return different answers to different clients. The answer depends on the source network, a view on the server, or a different resolver.

A common case: `app.example.com` returns a private address to office hosts and a public address to the Internet. VPN users can see the internal view when they use the office resolver.

Split horizon is not a CNAME. The name is the same. The data differs by viewer. Two people can both be "correct" and still disagree.

Troubleshooting must record which resolver you asked and from which network. `dig @office-resolver` and `dig @public-resolver` can disagree. That disagreement is often the design, not a fault.

Do not publish an internal name on a public zone by accident. Do not assume that your laptop on home Wi-Fi sees the office view. Do not use split horizon as the only access control. A leaked internal address is still a secret to protect.

Geo-DNS is a related idea: different public addresses by client region. The goal is delay or law, not a private LAN.

### Questions

#### Theoretical questions

1. What does split horizon mean for one name?
2. Why can two resolvers give two answers for the same name?
3. What must you record when you debug a split-horizon name?
4. Is split horizon the same as a CNAME?
5. Why is split horizon not enough as the only access control?

#### Easy practical tasks

1. Write five sentences that explain an office view and a public view.
2. Draw two clients, two answers, one name.
3. Make a table: home resolver, office resolver. Add one row for a private app name (invent the addresses).
4. Write four sentences: VPN DNS versus home DNS.

#### Medium practical tasks

1. If you have a workplace or school VPN that you may use, resolve one internal name on VPN and off VPN. Write both results. Skip if you have no such access.
2. Write six sentences: a user who changes DNS to a public resolver while on the office LAN. What can break?
3. Compare geo-DNS and split horizon in a short paragraph. Use STE.

#### Advanced practical tasks

1. Write a one-page runbook: "name works at the office, fails at home." Include resolver, view, and TTL.
2. Design a lab with two views on paper: internal A record and public A record. List who may query each view.

---

## DNS over HTTPS / TLS (awareness)

Classic DNS between stub and recursive resolver is often cleartext UDP (or TCP) on port 53. A path observer can read the query name. A path attacker can try to change the answer if nothing else protects the session.

DNS over TLS (DoT) wraps DNS in TLS, often on port 853. DNS over HTTPS (DoH) sends DNS in HTTPS, often to a resolver URL. Both hide the query and the answer from the path between the stub and that resolver.

You still trust the resolver. The resolver sees the names. DoH and DoT do not hide names from the operator of the resolver. They also do not replace HTTPS for the web site. They only protect the DNS step.

A browser can use DoH and ignore the OS resolver. Then `dig` and the browser can disagree. That disagreement is a common debug surprise.

Do not assume that DoH means "anonymous DNS." Do not disable all local DNS logging in a company without a policy. Enterprises can use their own DoH or DoT endpoints.

This section is awareness. You do not need to implement a DoH client. You need to name the problem (cleartext queries) and the idea (TLS to a resolver).

### Questions

#### Theoretical questions

1. What can a path observer read in classic DNS on port 53?
2. What does DoT wrap?
3. What does DoH use as its transport idea?
4. Who still sees the query name when DoH or DoT is on?
5. Why can a browser and `dig` disagree when DoH is on?

#### Easy practical tasks

1. Write five sentences that compare port 53 cleartext with DoT or DoH.
2. Make a table: classic DNS, DoT, DoH. Add one row for port or URL idea and one row for who you trust.
3. Find whether your OS or browser has a "secure DNS" or DoH setting. Write the name of the setting. Do not turn on a random resolver on a work machine without permission.
4. Draw: stub, TLS, recursive resolver.

#### Medium practical tasks

1. Capture a normal DNS query on a LAN that you own. Write whether you can read the name. Then write what you would expect if DoH were in use.
2. Write six sentences: DoH to a public resolver versus the company resolver. Cover policy and privacy.
3. Read one public vendor page on DoH or DoT. Write four STE sentences. No long quotes.

#### Advanced practical tasks

1. Write a one-page awareness note for a team: what DoH hides, what it does not hide, and how it changes debug with `dig`.
2. If your lab allows it, enable DoH in a browser profile that you own. Compare a name lookup in the browser and in `dig`. Record the difference. Revert the setting if you need the OS resolver.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a full lookup from a URL host to an IP address. Use stub, recursive resolver, record type, and TTL.
2. How do A, AAAA, and CNAME work together on a dual-stack name that is an alias?
3. Why must you name the resolver and the network when two people get different answers?
4. When is packet capture the wrong first DNS tool, and when is it the right tool?
5. A teammate says "DNS is just a hosts file on the Internet." Which facts do you use to correct that sentence?

#### Easy practical tasks

1. Write a one-page cheat sheet: stub, recursive, authoritative, seven record types, TTL, `dig` or `nslookup`, split horizon, DoH or DoT.
2. Run A, AAAA, and NS queries for `example.com`. Save the outputs in one text file. Label each block.
3. Draw one picture: home stub, home resolver, public authoritative path (idea). Label every term from this topic that appears.
4. List which DNS tools you have (`dig`, `nslookup`, `Resolve-DnsName`). Confirm each one with a version or a help line.

#### Medium practical tasks

1. From your machine, identify: configured resolver, one A answer, one TTL, and whether IPv6 answers exist for `example.com`. Write how you found each value.
2. Write a lab plan for the TLS topic that uses one name lookup first. State what you will record. Do not explain the TLS handshake yet.
3. Explain in ten steps how a beginner should query a public name, read the status, and avoid flooding.

#### Advanced practical tasks

1. Build a glossary of 16 DNS terms from this topic. Each entry: term, one-sentence definition, and one command or setting that shows it.
2. Capture one `curl` to a public HTTPS site and, if visible, the DNS query before TCP. Write a timeline. Mark each step as "I can name this now" or "later topic".
