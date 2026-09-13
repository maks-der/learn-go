# 13. Integration Styles

## Description

Integration is the way two systems exchange data or start work. This topic covers four classic styles: file transfer, shared database, remote procedure call (RPC), and messaging. It also surveys Hohpe and Woolf patterns and the anti-corruption layer.

Integration is a costly decision (Topic 1). A hidden share of tables is still integration (Topic 9). Complete Topics 1 to 12 before this topic. Pair contracts with Topic 8. Pair messages with Topic 10. Pair remote failure with Topic 11.

Use one term for each concept. File transfer moves a document. A shared database is one store that more than one system writes or reads as if it were a private schema. RPC is a request and a reply. Messaging is an asynchronous message. An anti-corruption layer is a translation boundary. A pattern name is not a design if you omit the constraint.

---

## File transfer

File transfer is integration through a file that one system writes and another system reads. The file can sit on a disk, an object store, or an SFTP host. The format is part of the contract: CSV, JSON lines, XML, or a fixed-width text.

The producer and the consumer do not need to run at the same time. That time split is the main benefit. A nightly export can feed a warehouse. A partner can drop a catalog file once a day.

You must define:

- Location and name rule
- Format and version
- Encoding and time zone
- Completeness signal (a done file, a checksum, or an object event)
- Retention and delete rule
- What happens when a file is late, empty, or duplicate

File transfer is a poor fit for a user click that needs an answer in 200 ms. The style fits batch work, partner exchange, and bulk load.

Idempotent load is required. A retry can send the same file again (Topic 4). Poison files need a quarantine folder, not an infinite retry (Topic 10).

Do not treat an undocumented spreadsheet as a contract. Write the columns. Version the header. Reject unknown required fields with a report.

Security is part of the style. Use authenticated transport. Limit who can write the drop location. Do not put secrets in the file name.

### Questions

#### Theoretical questions

1. What is file transfer as an integration style?
2. What benefit does the time split give?
3. Why is file transfer a poor fit for a user click with a tight latency budget?
4. What signal tells the consumer that a file is complete?
5. Why must a file load be idempotent?

#### Easy practical tasks

1. Write five sentences that define file transfer. Use only facts from this section.
2. Make a table: "Contract item" and "Example". Add format, name rule, completeness signal, and retention.
3. List four failure modes of a nightly CSV drop.
4. Write a done-file rule and a checksum rule in four sentences.

#### Medium practical tasks

1. Design a daily catalog file from a vendor to a shop. Include columns, version, and late-file behavior.
2. Write eight sentences that compare a file drop with an HTTP API for the same catalog.
3. Plan a quarantine path for a file that fails validation. Include operator notice and retry policy.

#### Advanced practical tasks

1. Write a one-page file-exchange contract for student grades: privacy, retention, and access list.
2. Compare object-store events with a polling job for new files. Write operations cost and duplicate risk.

---

## Shared database

A shared database is a store that more than one system uses as a common schema. Both systems issue SQL (or an equivalent query) against the same tables. Hohpe and Woolf list this style as a classic option. This handbook treats it as a high-coupling option.

Topic 9 names the usual defect: an integration database. A column change needs a multi-system release. Invariants spread. Transactions hide a missing API. Teams block each other.

A shared database can appear to be fast. A join across two owners looks cheap. The cost arrives later as a freeze on schema change and as unclear ownership.

When a shared database is still in place:

- Name one owner per table if you can.
- Stop new cross-owner writes.
- Put a view or a read replica in front of reporting if the report is the only extra reader.
- Plan an exit: API, events, or file export.

A modular monolith with one database is not this style. One deployable unit and module-owned tables is an internal design (Topics 5 and 9). Shared database as integration means more than one system, team, or deployable unit treats the schema as a public contract.

Do not add a second application user that writes the same tables "for speed". That choice is an architecture decision. Write an ADR if you accept it. Write the date that you will remove it.

### Questions

#### Theoretical questions

1. What is a shared database as an integration style?
2. How does this style differ from one database inside a modular monolith?
3. Why does a column change become expensive?
4. What hidden contract does a cross-owner join create?
5. What exit styles can replace a shared database?

#### Easy practical tasks

1. Write four sentences that define shared-database integration.
2. Make a table: "Access" and "Shared-database risk (high/low)". Add two apps write orders; one app reads a reporting replica; one module owns tables in one process.
3. List five symptoms from Topic 9 that appear when two services share tables.
4. Write an ADR title that rejects a new writer on the order tables.

#### Medium practical tasks

1. Given two teams and one `customers` table, write a six-month exit plan to an API or events.
2. Write eight sentences on why a reporting join can stay on a replica while writes must not share tables.
3. Review a vendor "single data platform" slide. Mark where it becomes an integration database.

#### Advanced practical tasks

1. Write a one-page coupling report: tables, writers, readers, and change freeze risk. Use a fictional campus system.
2. Design expand-contract steps that move one writer off a shared table without a flag-day cut (Topic 18).

---

## RPC

RPC (remote procedure call) is integration through a request and a reply. The caller waits. The callee does work and returns a result or an error. HTTP JSON APIs and gRPC are common forms (Topic 8).

RPC fits a user-facing question that needs an answer now: "Is this ISBN on the shelf?" RPC also fits a command that must accept or reject before the caller continues.

RPC costs:

- Both sides must be up at the same time (unless you hide a queue behind the callee).
- Latency includes the network (Topic 11).
- Timeouts and retries are mandatory.
- A chain of RPC calls multiplies failure (A calls B calls C).

Idempotency keys belong on writes (Topics 4 and 8). A retry after a timeout can double a side effect if the handler is not idempotent.

Do not use RPC for work that the user does not need in the same click. Send a message and return "accepted" (Topic 10). That choice is an integration style change, not only a performance trick.

Version the contract. Backward compatibility is part of RPC (Topic 8). An internal RPC is still a contract if another team owns the other side.

A local function is not RPC. When you extract a service, a former function becomes RPC. Add the failure policy at that moment (Topic 12).

### Questions

#### Theoretical questions

1. What is RPC as an integration style?
2. When is a request-reply call the right style?
3. Why does a chain of RPC calls multiply failure?
4. Why do writes need idempotency keys?
5. When must you replace RPC with messaging?

#### Easy practical tasks

1. Write five sentences that define RPC integration.
2. Make a table: "Need" and "RPC or message". Add "show stock now", "send receipt", and "rebuild search index".
3. List four headers or metadata fields that an RPC call must carry (timeout, identity, request identifier, idempotency key).
4. Draw A → B → C. Mark three timeouts.

#### Medium practical tasks

1. Design `GET /copies/{isbn}/availability` and `POST /loans`. Write timeout and idempotency rules.
2. Write eight sentences on a retry that doubles a payment if the handler is not idempotent.
3. Compare REST JSON and gRPC for one internal read. Write when each contract fits (Topic 8).

#### Advanced practical tasks

1. Write a one-page RPC policy: deadline propagation, retries, and when to open a circuit (Topic 11).
2. Map a three-call checkout onto a latency budget. Show that the chain cannot meet 300 ms if each hop is 150 ms.

---

## Messaging

Messaging is integration through messages that a broker or an event log stores and delivers (Topic 10). The sender does not wait for all receivers to finish.

Messaging fits a time split, a burst buffer, and a boundary between owners. The producer commits a fact or a job. Consumers work later.

Messaging costs:

- You must design duplicates, order, and poison messages (Topic 10).
- You must solve dual write (outbox or an equivalent).
- Operators must run a broker or a log.
- Users can see lag. Eventual consistency becomes visible (Topic 7).

Commands, events, and queries are different message kinds. Prefer request-reply APIs for queries. Prefer events for facts. Prefer queue messages for jobs.

Do not use a broker to hide a missing module boundary (Topic 10). A message that requires an immediate reply from three owners is often an RPC flow in disguise. That disguise adds a broker without removing temporal coupling.

Write the delivery promise. At-least-once delivery is the usual honest promise. Exactly-once is a property of a careful handler plus a store, not a slogan on a vendor page.

Pair this section with `kafka.topics.md` if you need an event log. Pair it with Topic 12 if the flow is a saga.

### Questions

#### Theoretical questions

1. What is messaging as an integration style?
2. What three needs make messaging a fit?
3. What costs does a broker add?
4. Why are queries a poor default on a message bus?
5. What is the usual honest delivery promise?

#### Easy practical tasks

1. Write four sentences that define messaging integration.
2. Make a table: "Kind" and "Example". Add command, event, and query.
3. List five operational items you need before the first production queue.
4. Write five sentences that compare RPC and messaging for "send a receipt".

#### Medium practical tasks

1. Design an `OrderPlaced` event and a mail consumer. Include duplicate handling and a poison path.
2. Write when a saga (Topic 12) must use messages instead of a chain of RPC writes.
3. Explain dual write in eight sentences. Point to the outbox as the usual repair (Topic 10).

#### Advanced practical tasks

1. Write a one-page messaging standard: headers, version, idempotency, and retention.
2. Compare a broker queue and an event log for partner integration. Write replay and operations cost.

---

## Hohpe / Woolf patterns (survey)

Gregor Hohpe and Bobby Woolf documented a pattern language for messaging systems. The catalog is public at [https://www.enterpriseintegrationpatterns.com/](https://www.enterpriseintegrationpatterns.com/). This section is a survey. It is not a full catalog.

A pattern is a named solution to a repeated problem in a context. Use a pattern when you can state the problem. Do not collect pattern names as decoration (Topic 1).

Useful groups for this path:

- Channel: how messages travel (point-to-point, publish-subscribe).
- Message: command, document, or event.
- Endpoint: sender and receiver adapters.
- Routing: content-based router, splitter, aggregator, filter.
- Transformation: translator, content enricher, claim check.
- System management: wire tap, control bus, detour.
- Failure: dead-letter channel, invalid-message channel.

Many product features map to these names. A dead-letter queue is a dead-letter channel. A routing key is often a content-based router. An outbox is an application-side reliability pattern that sits next to the catalog.

How to use the catalog:

1. Name the problem in one sentence.
2. Open the pattern page.
3. Write the context and the forces in your words.
4. Accept or reject the pattern in an ADR.
5. Record the product feature that implements it.

Do not invent new names for a pattern that the catalog already names. Shared language helps reviews.

The four styles in the earlier sections (file, shared database, RPC, messaging) are the top-level choices. Most Hohpe and Woolf patterns refine messaging. File and RPC still need versioning and error reports.

### Questions

#### Theoretical questions

1. What is a pattern in this handbook?
2. Where is the Hohpe and Woolf catalog?
3. Name four pattern groups from this survey.
4. How does a dead-letter queue map to the catalog?
5. What five steps do you use before you adopt a pattern?

#### Easy practical tasks

1. Open [https://www.enterpriseintegrationpatterns.com/](https://www.enterpriseintegrationpatterns.com/). Write the names of six patterns that you opened.
2. Make a table: "Product feature" and "Pattern name". Add pub/sub, DLQ, and routing key.
3. Write four sentences that explain why a pattern name without a problem is empty.
4. List three patterns that help operators see traffic (wire tap or equivalent).

#### Medium practical tasks

1. Pick Content-Based Router. Write the problem, the context, and one campus example in one page.
2. Map your Topic 10 design (outbox, inbox, poison) to catalog names. Note what the catalog does not name.
3. Write a review checklist of eight questions that cite pattern forces, not product logos.

#### Advanced practical tasks

1. Write a one-page glossary: ten pattern names and one sentence each in your own words. Do not copy long passages.
2. Take a public integration diagram. Label three Hohpe and Woolf patterns. Mark one box that is only a product name.

---

## Anti-corruption layer

An anti-corruption layer (ACL) is a boundary that translates an external model into your model. The external system keeps its names and its rules. Your domain keeps its names and its rules (Topic 7). The ACL sits between them.

Without an ACL, foreign names leak. Your code starts to use the vendor account identifier as if it were your `PatronId`. A vendor field change then breaks your domain.

An ACL can be:

- An adapter module in a hexagonal design (Topic 6)
- A small dedicated service that only translates
- A translator in a message pipeline (Hohpe and Woolf)

The ACL maps types, identifiers, and errors. It rejects data that your domain cannot accept. It does not become a second home for the vendor rules. If you copy all vendor behavior into the ACL, you only moved the corruption.

Use an ACL when:

- You integrate with a system that you do not own.
- The ubiquitous language differs.
- You must protect your model from frequent vendor change.

Do not use an ACL for two modules that you own and that must share a language. That case needs a better boundary or one context, not a translation tax.

Test the ACL with contract tests on both sides. Version the map. Record rejected fields.

A shared database has no ACL if both sides query the same columns. That is one reason shared database is hard to evolve.

### Questions

#### Theoretical questions

1. What is an anti-corruption layer?
2. What leaks when you omit an ACL?
3. Where can an ACL live in a structure?
4. When must you add an ACL?
5. When must you not add an ACL?

#### Easy practical tasks

1. Write five sentences that define an ACL. Use only facts from this section.
2. Make a table: "External name" and "Your name". Add five fields from a payment vendor to a shop.
3. List four jobs of an ACL (type map, identifier map, error map, reject).
4. Draw your domain, the ACL, and a vendor API.

#### Medium practical tasks

1. Design an ACL for a campus directory. Map vendor user records to `Patron`. Write two rejected cases.
2. Write eight sentences that compare an ACL with a shared table of vendor columns inside your domain.
3. Place an ACL in a hexagonal diagram (Topic 6). Label the port and the adapter.

#### Advanced practical tasks

1. Write a one-page ACL test plan: fixtures, version, and what you do when the vendor adds a required field.
2. Compare a translating module with a translating service. Write when the extra deployable unit is worth the cost (Topic 12).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do time split, coupling, and failure mode differ across file transfer, shared database, RPC, and messaging?
2. Why does this path treat shared database as a known style and as a usual anti-pattern at the same time?
3. How do Hohpe and Woolf names help a review without becoming decoration?
4. How does an anti-corruption layer protect a bounded context (Topic 7) during integration?
5. What makes an integration choice architectural in the sense of Topic 1?

#### Easy practical tasks

1. Write a one-page cheat sheet of the four styles, ACL, and three catalog pattern names.
2. For a library that must receive a vendor MARC dump and also show live shelf status, pick two styles and write one sentence each.
3. Make a table: "Style" and "Typical latency". Add the four styles in qualitative words (seconds, hours, milliseconds).
4. Bookmark [https://www.enterpriseintegrationpatterns.com/](https://www.enterpriseintegrationpatterns.com/). Write one sentence on when you open it.

#### Medium practical tasks

1. Write a short integration brief for a campus shop: payments (RPC), receipts (messaging), catalog (file), and a rejected shared database.
2. Take a teammate design that only lists product names (SFTP, Kafka, Postgres). Rewrite it as styles, contracts, and failure behavior.
3. Map Topics 8, 9, 10, and 11 onto one partner integration. Write what each topic forces you to write down.

#### Advanced practical tasks

1. Design a two-year exit from a shared reporting database to files plus an API. Include ACL, compatibility windows (Topic 18), and owners.
2. Score four public integration blogs from 0 to 5 on contract clarity and failure design. Explain each score. Do not copy long passages.
