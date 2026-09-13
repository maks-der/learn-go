# 21. Advanced Topics

## Description

Advanced topics are ideas that you must recognize and that you must not adopt as a default. This topic covers event sourcing, CRDTs at awareness level, multi-region active-active, data mesh and lakehouse at awareness level, platform engineering and internal developer platforms, and formal methods with TLA+ as an optional tool for hard protocols.

These ideas are costly. They solve rare problems. Complete Topics 1 to 20 before this topic. You need a modular monolith, clear data ownership, and honest operations before you add any item here. Pair event ideas with Topic 10. Pair region ideas with Topics 15 and 19.

Use one term for each concept. Event sourcing is not "we use events". A CRDT is not a general merge for every form. Active-active is not multi-AZ. A data mesh is not a lake. A platform is not a cluster logo. TLA+ is a specification language. You rarely implement Raft yourself (Topic 11).

---

## Event sourcing (costly; rare default)

Event sourcing stores the source of truth as a sequence of events. You rebuild state by replay. A current balance is a projection. The event log is the source of truth (Topic 9).

This is not the same as "publish events after you write a row". That pattern is event-driven integration (Topic 10). Event sourcing changes the write model.

Costs:

- You must version events forever
- You must design upcast rules
- PII delete and correction become hard (law, Topic 1)
- Projections lag. UX must tell the truth (Topic 20)
- Replay and storage cost rise
- People need training

When it can fit:

- A hard audit of every change is the product
- You already have a team that operates an event log
- You can state delete and correct rules

When it does not fit:

- A normal CRUD campus app
- A young product with a changing language (Topic 7)
- A two-person team

CQRS often appears with event sourcing (Topic 9). CQRS without event sourcing is enough for many read/write splits.

Do not start a side project on event sourcing to look advanced. Start with a row and an outbox.

If you study it, read a careful source and write an ADR that names the law constraint. Do not copy vendor slogans.

### Questions

#### Theoretical questions

1. What is event sourcing?
2. How does it differ from publish-after-write?
3. Name four costs.
4. When can event sourcing fit?
5. Why is a young ubiquitous language a reason to wait?

#### Easy practical tasks

1. Write five sentences that define event sourcing. Use only facts from this section.
2. Make a table: "Style" and "Source of truth". Add row plus outbox, and event sourcing.
3. List five law or UX problems that replay creates.
4. Write four sentences on why CQRS can exist without event sourcing.

#### Medium practical tasks

1. Write a one-page reject ADR for event sourcing on a notes app. Include the default that you keep.
2. Explain in eight sentences how a "right to erase" request collides with an immortal log.
3. Sketch one audit domain (grade changes) where a full event history might be the product. Write the projection that the teacher sees.

#### Advanced practical tasks

1. Write a version-and-upcast note for two event shapes of `GradeChanged`. Do not implement a framework.
2. Compare event sourcing with a normal table plus an audit table. Write operations cost for a five-person team.

---

## CRDTs (awareness)

A CRDT (conflict-free replicated data type) is a data structure that replicas can update without a lock. Merge is defined so that replicas converge. You must use types that have that merge math.

CRDTs help some collaborative editors and some edge replicas. They do not merge arbitrary business objects. A bank balance is not a casual CRDT problem. Money usually needs a single writer or a carefully designed type and an audit.

Awareness rules:

- You can name CRDT when a paper or a product claims automatic merge
- You ask: which type? which operations? what is the user-visible conflict?
- You do not add a CRDT library to "solve offline" for a shop cart without a proof

Topic 20 said most campus apps use server wins plus a prompt. That rule stays. CRDTs are an alternative for a small set of types (sets, counters, some text).

Cost includes memory, complexity, and a team that can explain merge to support staff.

Do not invent a CRDT. Use a reviewed library if a real product need appears. Write an ADR.

This section is awareness. You do not need to implement a type.

### Questions

#### Theoretical questions

1. What is a CRDT at a high level?
2. Why can replicas merge without a lock?
3. Why is a bank balance a poor casual CRDT example?
4. What questions do you ask when a product claims CRDTs?
5. What remains the default for campus offline conflicts?

#### Easy practical tasks

1. Write four sentences that define CRDTs for a beginner.
2. Make a table: "Data" and "CRDT likely? (yes/no)". Add a shared grocery set and a loan due date with law rules.
3. List three costs of a CRDT library.
4. Write five sentences that compare server-wins with CRDT merge.

#### Medium practical tasks

1. Read a public CRDT overview page. Write six sentences in your own words. Do not copy long passages.
2. Explain in eight sentences how a counter CRDT still needs a product meaning for "undo".
3. Write an ADR title and context that rejects CRDTs for a grade book.

#### Advanced practical tasks

1. Write a one-page awareness note: three CRDT kinds (set, counter, text) and one product each.
2. Compare CRDTs with operational transform at a one-paragraph awareness level. State that you would not implement either in this path.

---

## Multi-region active-active

Active-active multi-region means more than one region accepts writes at the same time. Users can work in two geographies without a failover wait.

This is not multi-AZ in one region (Topic 19). This is not active-passive DR (Topic 15).

Hard problems:

- Latency between regions (Topic 11)
- Conflict if two regions write the same key
- Session and identity stickiness
- Legal location (some data must not leave a region)
- Test cost and people cost
- CAP trade-offs in practical words (Topic 11)

Most products that need a second region use active-passive: one region writes, the other is standby. RPO and RTO are enough (Topic 15).

Active-active can fit a product with a natural partition (tenant or region as the write owner) and a rare global key. A global unique ISBN catalog with writes in two continents is a research project unless you already have a platform team.

Do not sell "active-active" on a slide if you have one primary database in one region. That claim is false.

If you need it, you will also need Topic 16 at a high level of discipline and a conflict rule that product accepts.

Consensus across regions is slow. You rarely implement it (Topic 11). You configure a product that already solved a subset.

### Questions

#### Theoretical questions

1. What is multi-region active-active?
2. How does it differ from multi-AZ?
3. How does it differ from active-passive DR?
4. Name four hard problems.
5. When is a partition by tenant a way to avoid dual writers on one key?

#### Easy practical tasks

1. Write five sentences that define active-active.
2. Make a table: "Claim" and "True?". Add two AZs one region, hot standby region, and two writing regions.
3. List five extra costs versus a single region.
4. Draw two regions that both accept writes. Mark a conflicting key.

#### Medium practical tasks

1. Write a one-page reject ADR for active-active on a campus shop. Accept multi-AZ plus backups.
2. Explain in eight sentences how a legal data-location rule can forbid a write in the second region.
3. Design a tenant-partitioned write: campus A in region A, campus B in region B. Write the global report lag.

#### Advanced practical tasks

1. Write a conflict table: last-write-wins, merge, and reject. Give one user sentence each.
2. Compare a vendor global-table product with two independent systems plus files (Topic 13). Write lock-in (Topic 19).

---

## Data mesh / lakehouse (awareness)

A data mesh is an organization and architecture idea: domain teams own analytical data products, with platform support and shared rules. A lakehouse is a technical style that stores large analytical data in files plus table formats, and serves warehouse-like queries.

These are awareness terms. They appear on slides. They are not a default for an operational catalog API.

Awareness rules:

- Operational source of truth stays with the owner (Topic 9)
- Analytical copies can lag. That is eventual consistency (Topic 7)
- A mesh without owners is a data swamp
- A lakehouse without access rules is a leak (Topic 14)
- You do not need a mesh to export a nightly file (Topic 13)

A small team uses a replica or a file export for reports. That is enough.

If your organization adopts a mesh, your job as an application architect is to publish a clear data product: schema, owner, freshness SLO, and privacy class. You do not rebuild the shop as a mesh.

Do not rename a shared reporting database as a mesh. Topic 9 still calls that an integration risk if writers multiply.

Read one public overview later. Write the terms in your own words. Do not copy long passages.

### Questions

#### Theoretical questions

1. What is a data mesh at awareness level?
2. What is a lakehouse at awareness level?
3. Why is a mesh without owners a swamp?
4. What must an operational system still own?
5. What report path is enough for a small team?

#### Easy practical tasks

1. Write four sentences that separate mesh (organization) from lakehouse (technology).
2. Make a table: "Need" and "Tool". Add nightly CSV, live checkout, and a company-wide analytic platform.
3. List four fields of a data-product charter (schema, owner, freshness, privacy).
4. Write five sentences that reject a mesh for a two-person shop.

#### Medium practical tasks

1. Design a nightly analytical copy of loans. Write freshness and access rules.
2. Explain in eight sentences how a lakehouse table is not the source of truth for a return at the desk.
3. Review a vendor "mesh" slide. Mark what is organization change versus a new store.

#### Advanced practical tasks

1. Write a one-page awareness glossary: warehouse, lake, lakehouse, mesh, operational DB. One sentence each in your words.
2. Compare a mesh data product with an anti-corruption layer on ingest (Topic 13). Write who owns quality.

---

## Platform engineering and internal developer platforms

Platform engineering is the work of building an internal product that other developers use to ship. An internal developer platform (IDP) is that product: templates, pipelines, golden paths, and paved infrastructure.

The platform is not Kubernetes by itself (Topic 19). The platform is the opinionated path: "here is how this company deploys a service with logs, identity, and a database".

Goals:

- Reduce every-team reinvention
- Encode Topics 14 to 18 as defaults
- Keep the golden path smaller than the full cloud catalog

Costs:

- A platform team needs people
- A bad platform becomes a gate that nobody can change
- One-size templates can fight a real constraint

A two-person product team is not a platform team. They can still keep a small golden path: one pipeline, one PaaS, one managed database. That is platform work at a tiny scale.

Do not wait for an IDP to start a modular monolith. Do not build an IDP to justify microservices (Topic 12). The extract still needs a boundary.

If you join a company with a platform, learn the golden path first. An exception needs an ADR and a time limit.

Measure the platform as a product: time to first deploy, and time to a new environment. Those are SLIs for the platform (Topic 15).

### Questions

#### Theoretical questions

1. What is platform engineering?
2. What is an internal developer platform?
3. Why is a cluster not a platform by itself?
4. What risk does a platform-as-gate create?
5. What is a golden path for a two-person team?

#### Easy practical tasks

1. Write five sentences that define an IDP.
2. Make a table: "Item" and "Platform or product app". Add pipeline template, loan policy, and log standard.
3. List four defaults that a golden path must include (identity, secrets, pipeline, and one more).
4. Write four sentences on why a two-person team must not start with a platform rewrite.

#### Medium practical tasks

1. Write a one-page golden path for a campus team: PaaS, managed DB, CI, and structured logs.
2. Explain in eight sentences how a platform can encode least privilege without becoming a ticket hell.
3. Write two platform SLIs: time to first deploy and time to add a preview environment.

#### Advanced practical tasks

1. Write an ADR: adopt the company golden path. List two allowed exceptions.
2. Compare a central platform team with each squad running raw cloud accounts. Write cost and security.

---

## Formal methods (TLA+) — optional, high value for hard protocols

Formal methods use mathematics to specify and check designs. TLA+ is a language and a set of tools for this work. You write a specification of states and allowed steps. A model checker looks for a broken invariant.

High value appears when the protocol is small and the bugs are costly: consensus, a lock service, a subtle failover. Topic 11 said you rarely implement Raft. If you design a new protocol, TLA+ can pay.

This work is optional. Most application features do not need it. Tests, types, and reviews are enough.

How to use it if you choose:

1. Specify the protocol, not the whole shop
2. Name invariants (no split brain, no double spend)
3. Keep the model tiny
4. Let the checker find a counterexample
5. Fix the spec. Then write code that matches the spec

Do not specify a 200-page application. The model will not finish. Do not treat a green model as a complete implementation.

You can read a public TLA+ example later. You do not need TLA+ to finish this path.

If a vendor claims a protocol is "proven", ask what was specified. The claim can be marketing.

Pair with chaos and failover tests (Topic 15). A spec does not replace a restore drill.

### Questions

#### Theoretical questions

1. What are formal methods in this handbook?
2. What is TLA+ used for?
3. When is the value high?
4. Why must the model stay tiny?
5. Why does a green model not replace a restore drill?

#### Easy practical tasks

1. Write four sentences that define TLA+ for a beginner.
2. Make a table: "Problem" and "TLA+ worth it? (yes/no)". Add a new consensus idea and a catalog CSS change.
3. List three invariants that a failover spec might name.
4. Write five sentences on why this path marks TLA+ optional.

#### Medium practical tasks

1. Write a one-page spec in plain sentences (not TLA+ syntax) for "one primary, one standby, never two writers".
2. Explain in eight sentences how a model-checker counterexample helps a review.
3. Find a public TLA+ overview. Write six sentences in your own words. Do not copy long passages.

#### Advanced practical tasks

1. Try a tiny public TLA+ tutorial if you want extra work. Write what invariant you checked. This task is optional.
2. Compare TLA+ with a tabletop game day (Topic 15). Write when you use each.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. What shared warning applies to event sourcing, CRDTs, and active-active?
2. How do data mesh and lakehouse differ from the operational source of truth?
3. How can a small golden path give platform value without an IDP department?
4. When does TLA+ pay, and when does it waste a semester?
5. Why does this path place advanced topics after you can operate a simple system?

#### Easy practical tasks

1. Write a one-page cheat sheet: each idea, default (use / avoid), and one cheaper substitute.
2. For a to-do app, mark all six ideas as avoid. Write one substitute for each.
3. Write six ADR titles that reject the six ideas for a four-person campus team.
4. Draw a poster: "rare default" with the six names and one cost each.

#### Medium practical tasks

1. Write a two-page "not yet" brief that a teammate can attach to a fashion-word slide.
2. Take a public conference talk title that uses three of these words. Rewrite the talk as problems and substitutes.
3. Map Topics 9, 10, 15, and 19 onto a reject of active-active event sourcing in two regions.

#### Advanced practical tasks

1. Write a decision tree: if audit-is-the-product, if offline-merge-is-the-product, if region-write-is-the-product. End most branches at "do not".
2. Read one careful public article on one idea only. Score the article 0 to 5 on costs named. Explain the score. Do not copy long passages.
