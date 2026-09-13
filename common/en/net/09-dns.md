# 9. DNS

## Description

DNS maps a name to data. The common case is a name to an IP address. This topic shows recursive versus authoritative resolver. You learn A, AAAA, CNAME, MX, NS, and TXT. You learn TTL. You learn `dig`. You learn DNS over HTTPS and DNS over TLS at an awareness level.

Complete this topic after HTTP. You already use names in URLs.

Use one term for each concept. A name is a domain name such as `example.com`. A resolver is a program that asks DNS questions. An authoritative server answers from a zone that it hosts. TTL is the time that a cache may keep an answer. Do not mix a recursive resolver with an authoritative server. Do not treat a name as one host.

Practice with `dig` or `nslookup`. Capture UDP port 53 when the query is cleartext. Encrypted DNS hides the query from the path. You still see the resolver address.

---

## Recursive vs authoritative resolver

The operating system has a stub resolver. The stub sends a query to a recursive resolver. The recursive resolver is often the home router, an ISP resolver, or a public resolver.

A recursive resolver accepts a query from a stub. It finds the answer. It can use its cache. If the cache has no valid data, it asks other servers. It starts at a root hint, then a TLD server, then the authoritative server for the zone. Then it returns the answer to the stub.

An authoritative server answers from zone data that an administrator or a registry put there. It does not need to recurse for names in that zone. A server can be authoritative for some zones and recursive for other queries. Many public resolvers recurse and are not authoritative for `example.com`.

The chain of authority uses NS records. The parent zone lists the name servers of the child. You follow NS from the root to the zone.

A stub does not walk the tree in the usual OS setup. The stub asks one configured resolver. You can walk the tree yourself with `dig +trace` when the tool supports it.

Do not open a recursive resolver to the public Internet on a host that you do not intend as a public service. An open resolver can be abused. Do not change production DNS without a rollback plan.

### Questions

#### Theoretical questions

1. What is a stub resolver?
2. What does a recursive resolver do for a stub?
3. What does an authoritative server use as its source of data?
4. Where does recursion start when the cache is empty?
5. What record type points to the name servers of a zone?

#### Easy practical tasks

1. Write five sentences that separate stub, recursive, and authoritative.
2. Draw: program, stub, recursive resolver, authoritative server.
3. Make a table: resolver type, who queries it, what it stores.
4. Find your configured resolver in `ipconfig /all` or `/etc/resolv.conf`. Write the address.

#### Medium practical tasks

1. Run `dig +trace example.com` if the tool supports it, or read a public trace example. Write the first two steps you understand.
2. Write six sentences: why a home router is often recursive and not authoritative for `example.com`.
3. Compare `nslookup` server line with your OS resolver. Write whether they match.

#### Advanced practical tasks

1. Write a one-page walk of the tree from root to `example.com` using public `dig` output.
2. Write a one-page note: why an open recursive resolver is a risk. No attack steps.

---

## A, AAAA, CNAME, MX, NS, TXT

A record types name the data that DNS returns.

- A: IPv4 address
- AAAA: IPv6 address
- CNAME: another name (an alias). The client then looks up that name
- MX: mail exchanger host and a preference number
- NS: name server of a zone
- TXT: text. People use it for SPF, domain verification, and other claims

A name can have A, AAAA, both, or neither. Dual stack uses both when both exist. Happy Eyeballs can try both families.

CNAME has rules. A name that has a CNAME must not have other data in the usual zone rules. Do not put a CNAME on the zone apex if your DNS host forbids it. Some hosts use ANAME or ALIAS as a vendor workaround.

MX values are host names, not IP addresses. You still need A or AAAA for the exchanger.

NS at the parent and NS at the child must agree, or you get a lame delegation.

TXT can be long. A verifier looks for a specific string. TXT is not encryption.

Do not cache an address in a long-lived program without a TTL plan. Do not treat a From header as proven because MX exists. MX is a delivery path, not a proof of the sender.

### Questions

#### Theoretical questions

1. What does an A record return?
2. What does an AAAA record return?
3. What does a CNAME point to?
4. What does an MX record name?
5. What is a common use of TXT?

#### Easy practical tasks

1. Query A and AAAA for `example.com`. Write the answers.
2. Write five sentences that define the six types in this section.
3. Make a table: type, one-line meaning.
4. Query MX for a public domain that you may look up. Write the exchanger names.

#### Medium practical tasks

1. Query NS for `example.com`. Write the name server names.
2. Write six sentences: a CNAME to a CDN name, then A on the CDN name.
3. Find a TXT record on a public domain. Write whether SPF or another token appears. Do not copy secrets.

#### Advanced practical tasks

1. Use `dig` to collect A, AAAA, MX, NS, and TXT for one domain. Put the answers in a table.
2. Write a one-page note: CNAME at the apex and why vendors add ALIAS. Use public docs.

---

## TTL

TTL is time to live for a DNS answer in a cache. The value is in seconds. A recursive resolver must not serve a cached answer after the TTL ends. The stub and the OS also cache.

A long TTL reduces queries and can hide a change. A short TTL lets you move an address faster. It increases query load. Some resolvers clamp very short or very long TTLs.

When you change an A record, old caches still have the old address until the old TTL ends. Plan the change before you drop the old server. That surprise is a common outage.

Negative answers (NXDOMAIN) also have a TTL from SOA parameters. A wrong "does not exist" can cache.

Applications that store an IP forever ignore DNS. Mobile clients and long-lived servers must refresh.

Do not set TTL to 0 as a habit. Some caches will still hold a minimum. Do not lower TTL only after the outage started. Lower it before the change.

Topic 13 returns to DNS TTL surprises in debug.

### Questions

#### Theoretical questions

1. What does TTL measure?
2. Who must honor TTL in this section?
3. Why does a long TTL hide a DNS change?
4. Why lower TTL before a planned move?
5. What is a negative cache in this section?

#### Easy practical tasks

1. Query `example.com` A with `dig`. Write the TTL that the tool shows.
2. Write five sentences that explain cache and TTL.
3. Make a table: TTL 30, TTL 86400. Add one benefit and one cost each.
4. Query the same name twice. Write whether the TTL dropped.

#### Medium practical tasks

1. Write six sentences: a one-hour TTL and a server that you already turned off.
2. Compare TTL on A and on NS if both appear. Write both numbers.
3. Find your OS or browser DNS cache flush command. Write the command. Do not flush a work machine if policy forbids it.

#### Advanced practical tasks

1. Write a one-page cutover plan: lower TTL, wait, change A, keep old host, then raise TTL.
2. Write a small program that resolves a name twice and prints the addresses. Discuss whether the stdlib exposed TTL.

---

## `dig`

`dig` is a DNS query tool from BIND. It prints the question, the answer, the authority, and the extra sections. It prints the resolver that it used and the time.

Useful forms:

- `dig example.com A`
- `dig example.com AAAA`
- `dig example.com MX`
- `dig example.com NS`
- `dig example.com TXT`
- `dig @1.1.1.1 example.com A` (ask a specific resolver)
- `dig +short example.com A` (short output)
- `dig +trace example.com` (walk the tree)

`nslookup` is common on Windows. It can set type and server. `Resolve-DnsName` is a PowerShell cmdlet. Learn one tool well.

`dig` talks to a resolver. It does not prove that your browser uses the same resolver or the same cache. Encrypted DNS in the browser can bypass the OS resolver.

Do not flood a public resolver with a script. Do not query names that you have no reason to query in a loop. Do not treat `dig` SERVFAIL as "the site is down" without a second resolver.

Read the status: NOERROR, NXDOMAIN, SERVFAIL, REFUSED.

### Questions

#### Theoretical questions

1. What sections does `dig` print in a full answer?
2. What does `@` do in a `dig` command?
3. What does `+short` change?
4. Why can `dig` and a browser disagree?
5. What does NXDOMAIN mean?

#### Easy practical tasks

1. Run `dig example.com A` or `nslookup -type=A example.com`. Write the address and the status.
2. Write five sentences that explain `dig` as a query tool.
3. Make a table: flag or form, purpose. Add `+short` and `@`.
4. Run A and AAAA queries. Write both command lines.

#### Medium practical tasks

1. Query the same name at two resolvers (`@` your router and a public resolver if policy allows). Write any difference.
2. Write six sentences: SERVFAIL versus NXDOMAIN.
3. Use `dig +trace` or a documented trace. Write how many NS steps you count.

#### Advanced practical tasks

1. Write a small cheat sheet of ten `dig` forms that you will reuse.
2. Compare `dig` UDP and `dig +tcp` on a name. Write whether both succeed.

---

## DoH / DoT (awareness)

DNS over TLS (DoT) sends DNS on a TLS session, often on port 853. DNS over HTTPS (DoH) sends DNS in HTTPS, often to a resolver URL.

Both hide the query name from a path observer that is not the resolver. The resolver still sees the query. The IP addresses of the resolver still show on the path.

Awareness facts:

- the OS, the browser, or an app can choose DoH or DoT
- a home router capture on UDP 53 can go quiet
- enterprise filters that rely on seeing UDP 53 need a new plan
- DoH uses HTTP and TLS ideas from topics 8 and 10

DoH and DoT do not replace authoritative servers. They protect the stub-to-recursive hop. They do not encrypt the later walk if the recursive resolver still uses clear DNS to the authority (some resolvers encrypt more hops; you do not need that detail).

This section is awareness. You must know the names and what they hide. You do not need to run a DoH server.

Do not assume that encrypted DNS makes the destination of HTTPS invisible. SNI and addresses still leak unless other controls exist. Do not disable a company DoH policy on a work machine.

### Questions

#### Theoretical questions

1. What does DoT wrap?
2. What does DoH wrap?
3. What does encrypted DNS hide from a path observer?
4. Who still sees the query name?
5. Does DoH replace the authoritative server?

#### Easy practical tasks

1. Write five sentences that define DoH and DoT.
2. Make a table: UDP 53, DoT, DoH. Add port or URL idea.
3. Find whether your OS or browser has a "secure DNS" setting. Write the name of the setting. Do not change a work policy.
4. Draw: stub, TLS or HTTPS, recursive resolver.

#### Medium practical tasks

1. Write six sentences: why a UDP 53 capture can be empty when the browser still resolves names.
2. Read a public DoH resolver URL from a vendor page. Write the scheme and that you will not flood it.
3. Compare privacy of DoH with HTTPS SNI visibility in four sentences.

#### Advanced practical tasks

1. Write a one-page awareness note for a team: what encrypted DNS changes in debug (`dig` versus browser).
2. If you own a lab, capture UDP 53 during `dig` and then resolve in a browser with secure DNS on. Write the difference. Skip if you cannot control the setting.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do recursive cache, TTL, and an A record change work together during a cutover?
2. Why must you name both A and AAAA when you debug a dual-stack site?
3. How do NS and authoritative data differ from a CNAME alias?
4. When is `dig @resolver` the wrong tool to explain a browser failure?
5. What stays visible on the path with DoH: the query name, the resolver IP, or both, and which?

#### Easy practical tasks

1. Write a cheat sheet: stub, recursive, authoritative, A, AAAA, CNAME, MX, NS, TXT, TTL, dig, DoH, DoT.
2. Run `dig example.com A` and `dig example.com AAAA`. Write answers and TTLs.
3. Draw the walk from stub to authority with a cache hit as a shortcut.
4. Bookmark an RFC or MDN DNS glossary. Write when you open it.

#### Medium practical tasks

1. Build a table for one public domain: A, AAAA, MX, NS, TXT present or not.
2. Write a fault story: NXDOMAIN versus wrong A versus expired TTL after a move.
3. Write six sentences: mail MX plus A of the exchanger.

#### Advanced practical tasks

1. Use `dig +trace` and a second resolver query. Write a one-page map of who answered what.
2. Write a client that looks up A and AAAA and tries `curl` to one address of each family when both exist.
