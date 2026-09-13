# 13. Evolution and Practice

## Description

This topic turns the path into habit. You survey backend-for-frontend (BFF) and server-side rendering versus single-page applications. You learn why event sourcing and CRDTs are rarely the default. You write ADRs and a C4 context diagram. You split a monolith only after a real boundary pain. You read *Designing Data-Intensive Applications* (DDIA) and you write one post-mortem a week.

The file `architecture.topics.md` also lists a suggested practice order of ten steps. Use that list after you finish these sections. Complete Topics 1 to 12 before you treat this topic as done. Pair this path with `db.topics.md`, `net.topics.md`, and a language path. Architecture remains trade-offs (Topic 1).

Use one term for each concept. A BFF is a backend that serves one client type. SSR renders HTML on the server. An SPA renders in the browser after a script load. Event sourcing is not "we use events". A CRDT is not a general merge for every form. A post-mortem is a learning document after an incident.

---

## BFF, SSR vs SPA (survey)

This section is a survey. Clients are part of the system (Topic 1). They hold state, they cache, and they fail on poor networks (`net.topics.md`).

A BFF (backend for frontend) is a server-side process that exists to serve one client type. A web BFF shapes data for the browser. A mobile BFF shapes data for the mobile application. Each BFF talks to owner services or to a modular monolith (Topics 3 and 9).

The BFF can aggregate two or three reads into one client call, hide internal service names, apply client-specific authentication details, and trim fields that the client does not need (Topic 10).

The BFF must not become the home of domain rules. Price and loan policy stay in the owner (Topics 3 and 4). If every new business rule lands in the BFF, you built a second monolith at the edge.

A BFF is a relative of an API gateway (Topic 9). A gateway applies shared edge policy. A BFF applies one client shape. Do not add a BFF for a single-page student app that already talks to one API.

Server-side rendering (SSR) produces HTML on the server for a request. The browser shows content without a large client render first. SSR can help first paint, simple caches, and clients with little JavaScript.

A single-page application (SPA) loads a script and then renders in the browser. Later navigation can stay on the client. An SPA can fit rich interactive tools. It adds a public JavaScript bundle, a harder first paint, and a longer compatibility window with the API (Topic 5).

Many products mix the two: SSR for the first page, then client render for the next actions. Write the mix. UX consistency is what the person sees. Service consistency is what the stores agree (Topic 4). A spinner that hides a lagging subscriber is a UX lie if the write already failed.

Do not treat the browser as a thin decoration in front of microservices. A chatty SPA against ten services needs a BFF or a gateway aggregation. A modular monolith with one JSON API is enough for many campus apps.

### Questions

#### Theoretical questions

1. What is a BFF?
2. Which rules must not live in a BFF?
3. What is SSR?
4. What is an SPA?
5. How does UX consistency differ from service consistency?

#### Easy practical tasks

1. Write five sentences that define BFF, SSR, and SPA. Use only facts from this section.
2. Make a table: "Duty" and "BFF or owner service". Add field trim, loan policy, and aggregate reads.
3. Draw a browser, a web BFF, and two owner services.
4. List four risks of a BFF that contains domain rules.

#### Medium practical tasks

1. Design a mobile BFF response for a loan card versus a web catalog page. Write why the payloads differ.
2. Write an ADR: SSR pages plus a small script, not a large SPA, for a campus catalog.
3. Explain in eight sentences how a BFF reduces chattiness without becoming the source of truth.

#### Advanced practical tasks

1. Write a one-page BFF charter: client, owner services, auth, SLO, and a reject list of domain rules.
2. Compare BFF aggregation with GraphQL at a high level (Topic 5). Write operations cost for a four-person team.

---

## Event sourcing and CRDTs (rarely the default)

Event sourcing stores the source of truth as a sequence of events. You rebuild state by replay. A current balance is a projection. The event log is the source of truth (Topic 6).

This is not the same as "publish events after you write a row". That pattern is event-driven integration (Topic 7). Event sourcing changes the write model.

Costs of event sourcing:

- You must version events forever.
- You must design upcast rules.
- Personal-data delete and correction become hard (law, Topic 1).
- Projections lag. UX must tell the truth.
- Replay and storage cost rise.
- People need training.

When event sourcing can fit: a hard audit of every change is the product, you already operate an event log, and you can state delete and correct rules. When it does not fit: a normal CRUD campus app, a young ubiquitous language (Topic 4), or a two-person team.

CQRS often appears with event sourcing (Topic 6). CQRS without event sourcing is enough for many read/write splits. Do not start a side project on event sourcing to look advanced. Start with a row and an outbox (Topic 7).

A CRDT (conflict-free replicated data type) is a data structure that replicas can update without a lock. Merge is defined so that replicas converge. You must use types that have that merge math.

CRDTs help some collaborative editors and some edge replicas. They do not merge arbitrary business objects. A bank balance is not a casual CRDT problem. Money usually needs a single writer or a carefully designed type and an audit.

Most campus apps use server wins plus a user prompt when a mobile client was offline. That is enough. Do not add a CRDT library because a talk used the acronym.

If you study these ideas, write an ADR that names the law constraint and the default that you keep. Do not copy vendor slogans.

### Questions

#### Theoretical questions

1. What is event sourcing?
2. How does it differ from publish-after-write?
3. Name four costs of event sourcing.
4. What is a CRDT at a high level?
5. Why is a bank balance not a casual CRDT problem?

#### Easy practical tasks

1. Write five sentences that define event sourcing and CRDTs. Use only facts from this section.
2. Make a table: "Style" and "Source of truth". Add row plus outbox, and event sourcing.
3. List five law or UX problems that an immortal event log creates.
4. Write four sentences on why CQRS can exist without event sourcing.

#### Medium practical tasks

1. Write a one-page reject ADR for event sourcing on a notes app. Include the default that you keep.
2. Explain in eight sentences how a "right to erase" request collides with an immortal log.
3. Write when a collaborative notes field might study CRDTs and when server wins is enough.

#### Advanced practical tasks

1. Write a version-and-upcast note for two event shapes of `GradeChanged`. Do not implement a framework.
2. Compare event sourcing with a normal table plus an audit table. Write operations cost for a five-person team.

---

## Write ADRs and a C4 context diagram

Pick a system that you own or that you may change. A to-do tool, a campus club site, or a library of notes is enough. The project must run. A slide deck is not a side project.

Write ADRs in the repository (Topic 1). See [https://adr.github.io/](https://adr.github.io/). Start with three:

1. Source of truth (the primary store)
2. Identity (how a principal authenticates)
3. Deploy unit (one process or more)

Add ADRs when a decision is costly to reverse: public API style, a queue, a region, a secret store (Topics 5, 7, 10, and 12).

Rules from this path:

- One decision per ADR
- Context includes constraints (team, time, law, existing systems)
- Consequences include the quality that you weaken
- Status can change. Do not delete a superseded ADR (Topic 12)

Do not write twenty ADRs on day one. Three good ADRs beat a folder of slogans. Review the ADRs after you ship. If the code disagrees, fix the code or supersede the ADR.

The C4 model is a small set of diagram levels. The site is [https://c4model.com/](https://c4model.com/). This path uses context first. Add a container view when you have more than one process or store.

A context diagram shows people, your system as one box, and external systems. Each relation has a verb. A container diagram shows the deployable units and stores inside your system: web application, API, worker, database, and broker.

Rules for diagrams:

- One job per diagram
- Names match the code and the ADRs
- Trust boundaries can overlay later (Topic 10)
- Do not draw twenty boxes that the code does not have

A context diagram plus three ADRs is a complete first architecture pack for a student system.

### Questions

#### Theoretical questions

1. Why must the side project run?
2. Which three ADRs does this section name first?
3. What does a C4 context diagram show?
4. What does a C4 container diagram show?
5. What do you do when code and an accepted ADR disagree?

#### Easy practical tasks

1. Open [https://adr.github.io/](https://adr.github.io/). Write the purpose of an ADR in four sentences.
2. Open [https://c4model.com/](https://c4model.com/). Write the names of the four C4 levels in order.
3. Write an ADR title and a one-paragraph context for the primary store.
4. Draw a context diagram for a to-do app: user, your system, and one external mail provider.

#### Medium practical tasks

1. Write the three starter ADRs. Cross-link a consequence of the store ADR to the deploy ADR.
2. Draw a container diagram that is still one API process and one database. Label the relation.
3. Review a public repository that contains ADRs. Summarize one ADR in six sentences. Do not copy the full text.

#### Advanced practical tasks

1. Write an ADR index for eight decisions. Include status, date, and related quality attributes (Topic 1).
2. Draw context and container for a library loan system. Overlay two trust boundaries (Topic 10).

---

## Split a monolith only after a real boundary pain

A modular monolith is the default (Topic 3). A split is a costly decision (Topic 1). You split only after a real boundary pain.

Real pain looks like this:

- Two teams block each other on one release clock, and package boundaries already exist.
- A measured load needs a different scale path for one owned module (Topic 12).
- A fault in one module takes down a process that law or a contract requires to stay up.
- A vendor or a law forces a separate system.

Not pain:

- A blog post about microservices
- A line-count target
- A possible future scale need with no measure
- A desire to use a broker on day one

Before a split, the module must already have a language, owned tables, a contract, and tests (Topics 4, 5, and 6). The extract is a move of an existing boundary onto a new deploy unit (Topic 9). If you cannot draw the boundary on a domain diagram, do not draw it on a deploy diagram.

After a split you pay network cost, saga cost, and operations cost (Topics 8 and 9). Write those costs in the ADR. Name the quality that you weaken.

If the pain disappears after you improve the modular monolith (clearer packages, a replica, a queue for one job), keep the monolith. That result is a success.

The suggested practice order in `architecture.topics.md` puts "split a service only if the boundary is real" last. Follow that order.

### Questions

#### Theoretical questions

1. What is a real boundary pain in this handbook?
2. Name three items that are not pain.
3. What must exist before an extract?
4. What costs appear after a split?
5. Why is an improved monolith a success?

#### Easy practical tasks

1. Write five sentences about the split rule. Use only facts from this section.
2. Make a table: "Signal" and "Pain? (yes/no)". Add six rows.
3. Copy the ten-item suggested practice order from `architecture.topics.md` into your notes. Tick what you already did.
4. Write four sentences to a teammate who wants to split before the first user.

#### Medium practical tasks

1. Write an ADR that rejects a split for your side project and names the pain that would reopen it.
2. Take a modular monolith map. Write which module could extract and which data it must take. Then write why you wait.
3. Explain in eight sentences how a queue for mail (Topic 7) can remove a false reason to split.

#### Advanced practical tasks

1. Write a one-page extract checklist and fill it with "not ready" for each row of your project.
2. Compare a merge back to a monolith with a further split for a public case study. Write five facts. Do not copy long text.

---

## DDIA and one post-mortem a week

*Designing Data-Intensive Applications* (Kleppmann) is a book that covers data models, storage, replication, partitions, and derived data. The site is [https://dataintensive.net/](https://dataintensive.net/). Read it after you ship a modular monolith with one database. Pair chapters with Topics 6 to 8 and with `db.topics.md`. Do not use the book as a reason to start with a global shard plan.

A post-mortem (also called a learning review) is a short document after an incident or after a lab failure. It states what happened, what you expected, what you learned, and what you will change. It does not blame a person.

Write one post-mortem a week while you practice:

- A test that failed in a surprising way
- A timeout that you forgot
- A deploy that broke a contract
- A game-day note from Topic 11

Rules:

- Facts first. Time, symptom, impact, and the correlation identifier if you have one.
- One or two actions with an owner and a date.
- Share it with the team. Keep secrets out (Topic 10).
- A week without a production incident still has a lab post-mortem.

Official and standard resources for this path also include the C4 model, ADRs, twelve-factor, Enterprise Integration Patterns, and Google SRE books (`architecture.topics.md`). Pick one next book after DDIA. Finish it. Do not start four books in one week.

Practice in the order that `architecture.topics.md` lists if you can: modular monolith, structured logs, timeouts, three ADRs, cache or replica, one async job, SLOs, expand-contract, failure modes, then a split only if the boundary is real.

Do not put solutions to course labs into this handbook folder. Keep your practice in your own repository.

### Questions

#### Theoretical questions

1. When do you open DDIA relative to this path?
2. What is a post-mortem in this handbook?
3. Why does a post-mortem not blame a person?
4. What belongs in a weekly post-mortem if production had no incident?
5. Why must you not start four books in one week?

#### Easy practical tasks

1. Open [https://dataintensive.net/](https://dataintensive.net/). Write the part or chapter titles that match Topics 6 to 8.
2. Write a five-line post-mortem template: what happened, expected, learned, action, owner.
3. Copy the suggested practice order from `architecture.topics.md`. Schedule the first four items.
4. Bookmark C4, ADR, twelve-factor, EIP, DDIA, and SRE books. Write one sentence on when you open each.

#### Medium practical tasks

1. Read one DDIA chapter overview or table of contents section. Write a one-page map to this repository's matching handbook.
2. Write a real post-mortem for a defect in your side project this week. Keep secrets out.
3. Run the first four items of the suggested practice order if you have not run them. Write one sentence each.

#### Advanced practical tasks

1. Build a reading calendar: two handbook topics and one DDIA chapter per week for two months.
2. After four weekly post-mortems, write a two-page self-review: which topics 1 to 12 you still cannot explain, and which practice-order item is next.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do BFF and SSR/SPA choices depend on whether you have one API or many services?
2. Why does this path treat event sourcing and CRDTs as recognition topics, not as defaults?
3. How do three ADRs and a C4 context diagram prove that you used Topics 1 to 3?
4. What evidence of "real boundary pain" would convince you to extract after you already have packages and owned tables?
5. How do a weekly post-mortem and DDIA together beat either habit alone?

#### Easy practical tasks

1. Write a one-page cheat sheet: BFF, SSR, SPA, event sourcing warning, CRDT warning, three ADRs, C4 context, split rule, DDIA, weekly post-mortem, practice order.
2. Create a folder `architecture-practice` with three empty notes: `adrs.md`, `c4.md`, `postmortems.md`. Write one pass test in each note.
3. Draw a roadmap from topic 1 to topic 13 with four practice boxes on it.
4. Bookmark DDIA, C4, ADR, and `architecture.topics.md`.

#### Medium practical tasks

1. Fill the three notes as you finish the work. Each note must contain your diagrams and observations, not a pasted internet solution.
2. Write a short architecture pack for a campus lost-and-found: context diagram, three ADRs, and a reject of event sourcing.
3. Use the suggested practice order as a twelve-week calendar. Assign each item to a week.

#### Advanced practical tasks

1. After the four labs in the practice order (monolith, logs, timeouts, ADRs), write a two-page self-review against Topics 1 to 12.
2. Start exactly one next-step track (DDIA part 1, `db.topics.md` topic 1, or `net.topics.md` topic 1). Write the first week's log. Do not put official lab solutions into `common/en/architecture`.
