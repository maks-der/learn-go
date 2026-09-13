# 18. BGP and the Internet

## Description

BGP is the routing protocol between autonomous systems on the public Internet. This topic is intermediate. It shows AS numbers, the path-vector idea, peering versus transit, anycast, why a bad leak can break many paths, and RPKI at an awareness level.

Use one term for each concept. An autonomous system (AS) is a network with one routing policy. A path-vector protocol carries a list of AS hops. Complete this topic after the routing survey in topic 6.

You do not need a full Internet table on a home router. Practice with public looking glasses, `whois` or `dig` for AS data, and news of past leaks. Do not announce prefixes that you do not hold.

---

## AS numbers

An AS number (ASN) is an identifier for an autonomous system. IANA and the RIRs assign ASNs. A public ASN appears in the global BGP table. A private ASN exists for local use and must not leak to the public Internet in the usual policy.

A small site often has no ASN. The site buys transit from an ISP. The ISP AS originates or aggregates the site prefixes. A large company, a cloud, or an ISP has one or more ASNs.

You can look up an ASN for a prefix with a looking glass or with `whois`. Many tools print `origin AS`. That value is the AS that the tool thinks originates the prefix. It is not a proof of ownership by itself. RPKI adds a signed check. A later section covers RPKI.

Do not treat an ASN as a company legal name. One company can have many ASNs. One ASN can serve many brands. Do not configure BGP on a device that you do not own.

IPv4 and IPv6 prefixes both travel in BGP. Dual stack means two families, not two Internets with no policy link.

### Questions

#### Theoretical questions

1. What does an ASN identify?
2. Who assigns public ASNs at a high level?
3. Why does a home LAN usually have no ASN?
4. What is a private ASN for?
5. Why is "origin AS" not full proof of ownership?

#### Easy practical tasks

1. Write five sentences that define AS and ASN.
2. Look up the origin AS for a public prefix (search a looking glass for `example.com` or a well-known resolver). Write the ASN.
3. Make a table: home LAN, ISP, large cloud. Add one row for "has a public ASN?"
4. Draw: your ISP AS and one neighbor AS (idea).

#### Medium practical tasks

1. Use `whois` or a web whois on a public IP that you may query. Write the ASN and the organization name if shown.
2. Write six sentences: two offices, no ASN, two ISP transits.
3. Find a public list of well-known ASNs (a cloud or a CDN). Write two names and numbers.

#### Advanced practical tasks

1. Write a one-page note: when a company should get an ASN and PI or PA address space (high-level, no registrar procedure dump).
2. Map one public website: name, A or AAAA, prefix, origin AS. Use public tools.

---

## BGP as a path-vector protocol (idea)

BGP sends updates about prefixes. Each update has an AS_PATH: the list of ASNs that the announcement already crossed. A router prefers a path by policy, then by path length, then by other tie breaks. Hop count of IP routers inside an AS is not the same as AS_PATH length.

Path-vector means: you see the AS sequence, not a full map of every router. A loop is visible when your ASN already sits in AS_PATH. You drop that update.

BGP is not OSPF. OSPF floods link state inside one organization. Topic 6 said that. BGP scales because it hides the internal topology and because it lets operators apply policy (prefer this peer, do not use that path for this prefix).

iBGP runs inside an AS. eBGP runs between ASNs. This topic only needs the names. You do not design a route reflector here.

Do not advertise all connected routes to the Internet. Do not accept a full table on a small CPU without a plan. Do not assume the shortest AS_PATH is the best user path. Policy and congestion matter.

Convergence can take time. A withdraw must spread. Users see a partial outage during that time.

### Questions

#### Theoretical questions

1. What does AS_PATH list?
2. How does a router detect a simple AS loop?
3. How is BGP different from OSPF in this survey?
4. What is eBGP versus iBGP in one sentence each?
5. Why is the shortest AS_PATH not always the best user path?

#### Easy practical tasks

1. Write five sentences that explain path-vector as "a list of AS hops."
2. Make a table: OSPF, BGP. Add rows for place and what is advertised.
3. Draw a path of three ASNs and write an AS_PATH.
4. Write four sentences: prefix plus next hop plus AS_PATH (idea).

#### Medium practical tasks

1. Open a looking glass traceroute or BGP view for a public prefix. Write the AS_PATH if shown.
2. Write six sentences: a prefer-customer policy versus shortest AS_PATH.
3. Read a public BGP "best path" summary. Write four STE sentences. No full RFC copy.

#### Advanced practical tasks

1. Write a one-page idea map: update, withdraw, AS_PATH, policy.
2. In a simulator that you own (or a paper lab), show a looped AS_PATH and say why it is dropped. No production peers.

---

## Peering vs transit

Transit is a paid (or contract) service. A transit provider forwards your traffic to the rest of the Internet. You send the provider your prefixes. The provider sends you a default or a full table. Money and an SLA often exist.

Peering is an exchange of traffic between two ASNs, often "settlement-free" at an Internet exchange (IX) or on a private link. Each side sends only its own prefixes and its customer prefixes, not a full Internet table. Peering reduces transit cost and can reduce delay.

A default-free AS has enough BGP data to reach all prefixes without a default route. Large ISPs aim for that. A stub AS uses a default toward transit.

Do not call every neighbor a peer. Peering has a policy meaning. Do not announce transit prefixes to a peer (that is a leak). Do not expect a free peer to carry your traffic to the whole Internet.

Relationships: customer, provider, peer. The valley-free idea says a path should not go customer-to-provider after it already went to a peer or a customer in ways that break the money graph. You only need: policy follows business.

### Questions

#### Theoretical questions

1. What does a transit provider give you that a typical peer does not?
2. What prefixes does a peer usually advertise?
3. What is a stub AS?
4. What is a leak in the peering sense?
5. Why can peering reduce delay?

#### Easy practical tasks

1. Write five sentences that compare peering and transit.
2. Make a table: transit, peering. Add rows for table size and money (typical).
3. Draw: stub AS, transit AS, rest of Internet.
4. Write four sentences: Internet exchange as a meeting point (idea).

#### Medium practical tasks

1. Find a public IX traffic or member page. Write the IX name and one fact.
2. Write six sentences: a company with one transit and one peer at an IX.
3. Explain "settlement-free" in four STE sentences.

#### Advanced practical tasks

1. Write a one-page business-plus-routing note: when to buy a second transit versus when to peer.
2. Sketch valley-free allowed paths on paper for provider, customer, and peer. High-level only.

---

## Anycast

Anycast means many hosts in different places use the same prefix (and often the same IP). BGP announces that prefix from each site. A client hits the site that its local BGP path prefers. DNS resolvers and CDNs use anycast.

Anycast is not multicast. The packet still has one destination address. Only one site handles that packet. The site can change if BGP changes. A long TCP session can break if the path moves to a different host that has no state. Designs keep state or use short exchanges (DNS UDP) or anycast-aware load balancers.

Health: a site must withdraw the prefix when it is down. If it keeps the announcement, BGP still sends users into a black hole.

Do not anycast a stateful database listener without a design. Do not treat anycast as a guarantee of the geographically nearest site. Policy can prefer a farther AS.

`ping` to an anycast address measures one site, not all sites. Two clients can reach two sites and still be "correct."

### Questions

#### Theoretical questions

1. What does anycast share across sites?
2. How does BGP make a client reach one site?
3. How is anycast different from multicast?
4. Why can a TCP session break when the anycast path moves?
5. What must a failed site do with its announcement?

#### Easy practical tasks

1. Write five sentences that explain anycast with a public DNS IP.
2. Make a table: unicast, anycast, multicast. Add one row for who receives.
3. Draw two cities, one prefix, two clients.
4. Write four sentences: DNS UDP as a fit for anycast.

#### Medium practical tasks

1. `dig` or `ping` a well-known anycast resolver from two networks if you can (home and phone). Write whether RTT differs a lot.
2. Write six sentences: anycast CDN versus a single-origin IP.
3. Traceroute to an anycast IP from your home. Write the last few hops. Do not claim you know the city.

#### Advanced practical tasks

1. Write a one-page note: anycast plus health withdraw. Include one TCP risk.
2. Compare two public looking-glass views of the same anycast prefix. Write whether AS_PATH differs.

---

## Why a bad BGP leak takes down the Internet

A leak is an announcement that should not leave a network, or that goes to the wrong neighbors. Example: a customer announces a full table or a more-specific prefix that it learned from one provider to another provider. Traffic then takes a path that cannot carry it. Paths fail. Users see outages far from the leak source.

A hijack is a different fault: someone announces a prefix that they should not originate. A more-specific prefix (`/24` versus `/16`) can win by longest prefix match. Topic 6 covered longest prefix. The Internet then sends traffic to the wrong AS.

Both faults spread because BGP trusts neighbors by default. Filters, max-prefix limits, and RPKI reduce the blast. They do not make BGP "safe by magic."

Do not announce more-specifics that punch a hole in someone else aggregate unless you have the right. Do not accept a huge prefix count from a small peer. Do not treat a social media outage as "the whole Internet" without a looking glass.

History has public postmortems. Read one as a lesson. You do not replay it.

### Questions

#### Theoretical questions

1. What is a leak in one sentence?
2. What is a hijack in one sentence?
3. Why can a more-specific prefix steal traffic?
4. Why do distant users suffer?
5. What knob limits how many prefixes you accept from a neighbor?

#### Easy practical tasks

1. Write five sentences that separate leak and hijack.
2. Make a table: leak, hijack. Add one row for intent (mistake versus theft idea).
3. Draw a customer that accidentally sends a full table to a provider.
4. Write four sentences: longest prefix match as the weapon of a more-specific.

#### Medium practical tasks

1. Read one public postmortem of a BGP leak or hijack. Write six STE sentences. No sensational quotes.
2. Write six sentences: max-prefix on a peering session.
3. Explain why a `/24` can beat a `/16` for one address.

#### Advanced practical tasks

1. Write a one-page incident sketch: how an operator confirms a leak with a looking glass. Checks only. No attack steps.
2. Design a filter idea: accept only customer prefixes that you listed. High-level.

---

## RPKI (awareness)

RPKI is Resource Public Key Infrastructure. Address holders publish a signed ROA (Route Origin Authorization). The ROA says which ASN may originate a prefix (and a max length). A BGP router can validate: valid, invalid, or unknown.

Invalid means the announcement does not match a ROA. Many networks drop invalids. Unknown means no ROA. The announcement can still be right. Valid means a matching ROA exists.

RPKI does not sign the full AS_PATH. A leak that keeps a valid origin can still happen. Path security is a later idea (ASPA and similar). This section is awareness of origin checks.

Do not treat RPKI valid as "this path is the best." Do not create a ROA that is tighter than your real announcements if you still need longer prefixes. A wrong ROA can make your own prefix invalid.

You can view ROA state in public validators and in some looking glasses. You do not need to run a certificate authority for this course.

### Questions

#### Theoretical questions

1. What does a ROA authorize?
2. What are the three validation states?
3. Does RPKI sign the full AS_PATH?
4. What can a wrong ROA do to your own prefix?
5. Why can a leak still occur when the origin is valid?

#### Easy practical tasks

1. Write five sentences that explain ROA as "who may originate this prefix."
2. Make a table: valid, invalid, unknown. Add one row for a typical action (drop or accept idea).
3. Look up RPKI state for a well-known prefix in a public validator. Write the state.
4. Draw: holder, ROA, validator, BGP router.

#### Medium practical tasks

1. Write six sentences: drop invalids versus "RPKI is optional."
2. Read a public RPKI one-pager from a RIR or an operator. Write four STE sentences.
3. Explain max-length in a ROA in four sentences (idea).

#### Advanced practical tasks

1. Write a one-page awareness note: what RPKI stops and what it does not stop.
2. For the prefix of a public site that you mapped earlier, record origin AS and RPKI state from public tools.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe how a home HTTP request leaves an ISP AS and reaches a CDN anycast prefix. Use ASN, transit or peer, and BGP as a path-vector idea.
2. How do longest prefix match, a leak of a more-specific, and RPKI origin checks interact?
3. Why do peering and transit need two words if both are BGP neighbors?
4. When is a looking glass the right tool, and when is traceroute enough?
5. A teammate says "BGP is just OSPF for the world." Which facts do you use to correct that sentence?

#### Easy practical tasks

1. Write a one-page cheat sheet: ASN, AS_PATH, eBGP, peering, transit, anycast, leak, hijack, ROA.
2. Look up one public IP: ASN, AS_PATH if shown, RPKI state if shown. Save the output.
3. Draw one picture: stub, transit, peer, anycast CDN.
4. List public tools you used (whois, looking glass, validator).

#### Medium practical tasks

1. Write a lab plan on paper: you never announce in production; you only read public views. State three queries.
2. Compare two paths to the same name from a looking glass in two cities if the tool allows. Write AS_PATH differences.
3. Explain in ten steps how a beginner reads a leak postmortem without copying rumors.

#### Advanced practical tasks

1. Build a glossary of 14 BGP terms from this topic. Each entry: term, one sentence, one public tool.
2. Write a short tabletop: a more-specific appears, RPKI invalid. What do you check first? No commands that change the Internet.
