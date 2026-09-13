# 20. Frontend / Client Architecture (Survey)

## Description

Client architecture is the structure of the programs that people use: browsers, mobile applications, and similar clients. This topic is a survey. It covers the backend-for-frontend (BFF) pattern, mobile offline behavior, server-side rendering versus single-page applications on the web, and the difference between user-experience consistency and service consistency.

Clients are part of the system (Topic 1). They hold state, they cache, and they fail on poor networks (`net.topics.md`). Complete Topics 1 to 19 before this topic. Do not treat the browser as a thin decoration in front of microservices. Pair APIs with Topic 8. Pair identity with Topic 14.

Use one term for each concept. A BFF is a backend that serves one client type. Offline is a first-class mode, not a defect. SSR renders HTML on the server. An SPA renders in the browser after a script load. UX consistency is what the person sees. Service consistency is what the stores agree (Topic 7).

---

## BFF (backend for frontend)

A BFF (backend for frontend) is a server-side process that exists to serve one client type. A web BFF shapes data for the browser. A mobile BFF shapes data for the mobile application. Each BFF talks to owner services or to a modular monolith (Topics 5 and 12).

The BFF can:

- Aggregate two or three reads into one client call
- Hide internal service names
- Apply client-specific authentication details
- Trim fields that the client does not need (Topic 14, least privilege on data)

The BFF must not become the home of domain rules. Price and loan policy stay in the owner (Topics 6 and 7). If every new business rule lands in the BFF, you built a second monolith at the edge.

A BFF is a relative of an API gateway (Topic 12). A gateway applies shared edge policy. A BFF applies one client shape. Some teams combine them. Write the rule. A huge plugin list is a risk.

Do not add a BFF for a single-page student app that already talks to one API. Add a BFF when two clients need different payloads or different release clocks (Topic 18).

The BFF is a fault domain. If the mobile BFF is down, the mobile app fails even if owner services are up. Give it SLOs and redundancy (Topic 15).

Version the BFF contract with the same care as any public API (Topic 8). Mobile clients force long compatibility windows.

### Questions

#### Theoretical questions

1. What is a BFF?
2. What jobs can a BFF do?
3. Which rules must not live in a BFF?
4. How does a BFF differ from an API gateway?
5. When do you not add a BFF?

#### Easy practical tasks

1. Write five sentences that define a BFF. Use only facts from this section.
2. Make a table: "Duty" and "BFF or owner service". Add field trim, loan policy, and aggregate reads.
3. Draw a browser, a web BFF, and two owner services.
4. List four risks of a BFF that contains domain rules.

#### Medium practical tasks

1. Design a mobile BFF response for a loan card versus a web catalog page. Write why the payloads differ.
2. Explain in eight sentences how a BFF reduces chattiness without becoming the source of truth.
3. Write an ADR: one public API versus a web BFF plus a mobile BFF for a campus shop.

#### Advanced practical tasks

1. Write a one-page BFF charter: client, owner services, auth, SLO, and a reject list of domain rules.
2. Compare BFF aggregation with GraphQL at a high level (Topic 8). Write operations cost for a four-person team.

---

## Mobile offline

Mobile clients lose network. Offline is a normal mode. Architecture must define what the person can do with no network and what happens when the network returns.

Typical offline work:

- Read a local cache of recent data
- Queue writes (a loan request, a note)
- Show a clear offline mark
- Replay the queue with idempotency keys (Topics 4 and 8)

Conflict is the hard part. The server can change while the client is offline. You need a rule: server wins, client wins, or merge. CRDTs are a later awareness topic (Topic 21). Most campus apps use server wins plus a user prompt.

Security does not vanish offline. Tokens live on the device. Use the platform store for tokens. Do not log tokens (Topic 14). Encrypted local stores still have device-theft risk. Write the data class that may sit on disk.

Sync is an integration style (Topic 13). It can use RPC when online and a file or message queue of local operations. Dual write on the device is still a dual-write problem. Use one local outbox.

Compatibility windows grow. Users skip updates (Topic 18). The server must accept old clients for the window.

Do not promise "full offline checkout" if payments need a live vendor. Promise the subset that you can test.

Measure sync lag and conflict rate (Topic 16). Those numbers are quality attributes (Topic 2).

### Questions

#### Theoretical questions

1. Why is offline a normal mode on mobile?
2. What write path do you use when the network is down?
3. What three conflict rules can you name?
4. Why do tokens on a device need extra care?
5. Why must you not promise a full offline payment if the vendor must be live?

#### Easy practical tasks

1. Write four sentences that define mobile offline for a beginner.
2. Make a table: "Action" and "Allowed offline? (yes/no)". Add read last catalog, place a paid order, and save a draft note.
3. List five user-visible states: online, offline, syncing, conflict, and failed sync.
4. Write a user message for a conflict on a note title.

#### Medium practical tasks

1. Design a local outbox for "request a book hold". Include idempotency and replay.
2. Explain in eight sentences how a long compatibility window interacts with an old offline client.
3. Write a one-page data-class policy: what may live on disk and for how long.

#### Advanced practical tasks

1. Write a sync playbook: backoff, poison local jobs, and operator support steps.
2. Compare server-wins with a user prompt. Write support cost for a campus help desk.

---

## Web: SSR vs SPA

SSR (server-side rendering) means the first HTML for a page is produced on the server. The browser shows content before or without a large client render. SPA (single-page application) means the server sends a shell and scripts. The browser builds the page and talks to APIs.

SSR benefits:

- Faster first content for many documents
- Simpler share of links and some search indexing
- Less need for a large client runtime on the first paint

SSR costs:

- Server CPU per request
- Cache rules for HTML (Topic 17)
- A mix of server templates and client scripts if you add interactivity

SPA benefits:

- Rich in-page interaction after load
- A clear split between API and UI teams

SPA costs:

- Slow or empty first paint if scripts are large
- Extra care for history, accessibility, and search
- Token storage and XSS surface (Topic 14). Stay defensive.

Many products mix: SSR for the first document, then client render for widgets. That mix is a structure. Name who renders each part.

The choice is architectural when you change the latency budget, the cache design, and the BFF shape. The choice is not architectural when it is only a framework slogan (Topic 1).

Do not run an SPA that chatters twenty internal services from the browser. That leaks topology and trust (Topics 12 and 14). Use one API or a BFF.

Measure Largest Contentful Paint or a similar first-content SLI if the web is the product. Pair with Topic 15.

### Questions

#### Theoretical questions

1. What is SSR?
2. What is an SPA?
3. What first-paint cost does an SPA add?
4. Why is a mix of SSR and client widgets a structure that you must name?
5. Why must a browser not call twenty internal services?

#### Easy practical tasks

1. Write five sentences that compare SSR and SPA.
2. Make a table: "Page" and "SSR, SPA, or mix". Add a public catalog, an admin grid, and a login page.
3. List four quality attributes that the choice changes (latency, cost, security, and one more).
4. Draw SSR versus SPA for the first request. Label who produces HTML.

#### Medium practical tasks

1. Write a one-page choice for a campus library site. Include search and a logged-in account page.
2. Explain in eight sentences how HTML cache at a CDN differs for SSR public pages and for personalized SPA APIs (Topic 17).
3. Write defensive token-storage notes for an SPA at a high level. Do not write attack steps.

#### Advanced practical tasks

1. Write an ADR: SSR catalog plus a small account SPA. Include the latency budget.
2. Compare a BFF that returns HTML (SSR) with a BFF that returns JSON (SPA). Write operations and cache cost.

---

## Consistency of UX vs service consistency

Service consistency is the agreement of stores and processes on facts (Topic 7). Transactional consistency and eventual consistency are service terms.

UX consistency is what the person believes is true on the screen. A UI can look consistent while services still catch up. A UI can look broken while services are correct.

Examples:

- The client shows "hold placed" from a local optimistic update. The server later rejects the hold. UX and service disagree until you repair the screen.
- The catalog page is stale for 30 seconds because of a cache. Service truth is fine. UX is old.
- Two tabs show two prices. The service has one price. The UX is inconsistent.

Design the UX for the consistency that you actually have. If the mail is async (Topic 17), do not show "mail sent" before the worker succeeds. Show "mail queued".

Optimistic UI is a choice. It needs a visible fail path and idempotent writes. Silent revert is a trust defect.

Across devices, the same user can see two states during sync. Write the rule. A spinner or a "last updated" time is part of the architecture, not decoration.

Do not force a distributed transaction to make a spinner disappear (Topics 11 and 12). Change the sentence on the button instead.

Accessibility and language belong here. An error that only appears as a red box without text is an inconsistent UX for many users.

SLIs can include UX facts: "share of hold requests that show a final success or fail in 10 seconds". That is stricter than "API 200".

### Questions

#### Theoretical questions

1. What is service consistency in this section?
2. What is UX consistency?
3. How can an optimistic update create a disagreement?
4. Why must async work not show "done" too early?
5. How can an SLI include a UX fact?

#### Easy practical tasks

1. Write four sentences that separate UX consistency from service consistency.
2. Make a table: "Screen text" and "Honest? (yes/no)". Add "mail sent", "mail queued", and "hold placed (pending)".
3. List five cases where two devices can disagree for a short time.
4. Write a better button sentence for an async export.

#### Medium practical tasks

1. Design an optimistic "add to list" flow. Write success, fail, and retry text.
2. Explain in eight sentences how a 30-second catalog cache is acceptable if the page shows freshness.
3. Write two SLIs: one for API success and one for "user sees a final state".

#### Advanced practical tasks

1. Write a one-page UX consistency standard for a campus shop: optimistic rules, stale labels, and multi-device sync.
2. Compare a saga in the middle (Topic 12) with the screens that a person must see. Map each saga state to one sentence.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do BFF, render style, and trust boundaries form one client-to-owner path?
2. How does mobile offline change APIs, windows, and idempotency at the same time?
3. Why can a correct service still produce a poor UX, and why can a smooth UX hide a service lie?
4. Why is this topic a survey after services and delivery, not a UI-framework course?
5. What makes a client choice architectural in the sense of Topic 1?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic.
2. For a notes product, write one BFF sentence, one offline rule, one render choice, and one honest status text.
3. Draw a C4 container view: browser, optional BFF, API, worker, and a mobile app with a local store.
4. Write three ADR titles: no BFF yet, SSR catalog, server-wins on note conflicts.

#### Medium practical tasks

1. Write a two-page client brief for a campus library: web mix, mobile offline hold queue, and UX sentences for a saga.
2. Take a teammate plan that is only "React plus microservices". Rewrite it as contracts, trust, and consistency.
3. Map Topics 8, 14, 17, and 18 onto one mobile release. Write the compatibility window.

#### Advanced practical tasks

1. Design a semester client architecture for web plus Android. Include BFF decision, offline subset, and SLIs.
2. Read a public offline-first or SSR overview. Write ten sentences that you would keep as team rules. Do not copy long passages.
