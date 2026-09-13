# 6. Data Architecture

## Description

Data architecture is the set of decisions about where data lives, who owns it, and how other parts read it. This topic covers one database per service versus a shared database, CQRS when reads and writes diverge, cache and search index versus the source of truth, and data ownership.

Data decisions are costly to change (Topic 1). A wrong share of tables creates hidden coupling (Topic 2). Pair this topic with `db.topics.md`. For cache details, pair with `redis.topics.md`.

Complete Topics 1 to 5 before this topic. Start with one database in a modular monolith (Topic 3). Split stores only when ownership or a measure requires it.

Use one term for each concept. The source of truth is the store that wins when copies disagree. A cache is not the source of truth. A search index is not the source of truth unless you design it that way (rare).

---

## One database per service vs a shared database

A shared database is one database that more than one deployable unit (or more than one unowned module) reads and writes. A "one database per service" rule says that each deployable service owns its data store. Other services use an API or events, not the tables.

In a modular monolith, one database is the default. Modules can still own tables. Ownership is a logical rule. Physical split comes later.

When more than one service writes the same tables:

- You cannot change a column without a multi-service release.
- Invariants spread across codebases.
- Transactions hide a missing API.
- Teams block each other.

When each service has a store:

- You can change the schema inside the owner.
- You lose easy cross-entity transactions.
- You need APIs, events, or a composition layer for joins.
- Operations cost grows (Topic 1: cost).

Do not copy the "one database per service" slogan into a three-person monolith. You already have one deployable unit. Use table ownership first.

A shared read-only replica of the same owned database is not the same as a shared integration database. Replicas are a scale tool. They still have one owner.

An integration database is a database that exists so that many programs can integrate by reading and writing the same tables. The database becomes the contract. This handbook treats that design as an anti-pattern for application integration. The schema cannot evolve. Security becomes "who has the database password". Preferred integration styles are API, messaging, or file transfer (Topic 9).

A data warehouse or an analytics replica can be a valid copy for reports. It is still a copy. The owner of the source of truth remains the operational system. Reports must not write back into operational tables without a defined process.

Write the owner of each store in an ADR. If two services "share for a short time", write the end date. Short times become years.

### Questions

#### Theoretical questions

1. What is a shared database in this handbook?
2. What does "one database per service" mean?
3. Why is one database the default in a modular monolith?
4. What goes wrong when two services write the same tables?
5. How is a read replica different from an integration database?

#### Easy practical tasks

1. Write five sentences that compare shared tables with owned stores.
2. Make a table: "Table" and "Owner module". Add six tables for a shop.
3. List four costs of a second database for a small team.
4. Draw one monolith with one database and two services with two databases.

#### Medium practical tasks

1. Given two services that share `users` and `orders`, write a split plan that starts with APIs, not with a new database product.
2. Write an ADR: one PostgreSQL database, table ownership, no service split.
3. Design a join that a UI needs (user name plus order). Show the monolith query versus two service calls.

#### Advanced practical tasks

1. Write a one-page decision tree: table ownership, then replica, then store split.
2. Write a two-page recovery plan for a three-app integration database: inventory, owner, API, sunset of SQL access.

---

## CQRS when reads and writes diverge

CQRS (command query responsibility segregation) is a split between the write model and the read model. A command changes the source of truth. A query reads a model that you built for that read.

You do not need CQRS for every application. A single model that you read and write is the default. Use CQRS when reads and writes diverge enough that one shape is a bad fit for both.

Signs that reads and writes diverge:

- A write is a small transactional change.
- A read needs a large join, a search ranking, or a dashboard that the write model cannot serve in the latency budget.
- Many readers need a denormalized copy.

A simple form: the same database, extra tables or views that a writer updates in the same transaction or through a projection job. A stronger form: a separate read store that lags (eventual consistency, Topic 4).

CQRS is not event sourcing. Event sourcing stores the source of truth as events (Topic 13). CQRS can exist with normal rows. Do not adopt event sourcing because you adopted CQRS.

Costs:

- Two models to change.
- Lag if the read store is separate.
- Risk that a reader treats a copy as the source of truth.

Do not add CQRS because a blog used the acronym. Measure the read that fails. Then add a read model for that read.

### Questions

#### Theoretical questions

1. What is CQRS?
2. When do you keep one model for read and write?
3. Name three signs that reads and writes diverge.
4. How does CQRS differ from event sourcing?
5. What cost does a separate read store add?

#### Easy practical tasks

1. Write five sentences that define CQRS. Use only facts from this section.
2. Make a table: "Need" and "One model or CQRS?". Add six needs.
3. List four reads in a shop that can stay on the write tables.
4. Draw a command path and a query path that share one database with one extra view.

#### Medium practical tasks

1. Design a catalog list that needs full-text search. Write what stays in the write database and what goes to a read model.
2. Write an ADR that rejects CQRS for a notes app and keeps one schema.
3. Explain in eight sentences how a lagging read model must appear in the user interface.

#### Advanced practical tasks

1. Write a one-page CQRS plan for a campus grade book: write invariants versus a teacher dashboard.
2. Compare a SQL view, a projection table, and a search index as three read models. Write operations cost.

---

## Cache and search index vs source of truth

The source of truth is the store that wins when copies disagree. In a typical business system the source of truth is the operational database (or the owned store of a service).

A cache is a copy that you keep for speed. A cache can be wrong or empty. If the cache and the source disagree, the source wins. You must define how the cache becomes wrong (TTL, explicit delete, or version). You must define what the user sees on a cache miss.

A search index is a copy that you keep for query shapes that the source of truth cannot serve well (full text, ranking, facets). The index is not the source of truth. A rebuild from the source must be possible. A failed index update must not silently invent a new book.

Rules:

- Write the source first (or write source and outbox together, Topic 7).
- Do not write only the cache.
- Do not treat "the document is in the index" as proof that the row exists.
- Bound cache keys. Do not use unbounded personal data as a key if law forbids it.
- Invalidate on write when you can. Accept a short TTL when you cannot.

A cache in front of a database can hide load. It can also hide a stampede when the TTL expires for many keys at once. Topic 12 covers scale. This section only states ownership of truth.

Read replicas are copies of the source for scale. They can lag. They are still the same logical source family. A search product is a different model. Do not mix the words.

If you design a system where the search index is the only store, write that ADR. That design is rare for core business records. It is more common for logs or for documents that you accept as eventual.

### Questions

#### Theoretical questions

1. What is the source of truth?
2. What is a cache in this handbook?
3. Why is a search index not the source of truth in a typical business system?
4. What must you define when the cache and the source disagree?
5. Why must a rebuild from the source be possible?

#### Easy practical tasks

1. Write five sentences that compare cache, search index, and source of truth.
2. Make a table: "Store" and "Wins on conflict?". Add primary database, Redis cache, and search index.
3. List four cache invalidation methods (TTL, delete on write, and similar).
4. Draw a write to the database and a later update of a search index.

#### Medium practical tasks

1. Design a book detail page: source of truth, optional cache, and what happens on cache miss.
2. Write an ADR: PostgreSQL is the source of truth; the search product is a rebuildable index.
3. Explain in eight sentences how a stampede can occur when many TTLs expire together.

#### Advanced practical tasks

1. Write a one-page cache policy: keys, TTL, personal data, and stampede control at a high level.
2. Design a repair job that rebuilds the search index from the source. Write how you detect drift.

---

## Data ownership

Data ownership is the rule that names who may write the source of truth and who may change the schema. An owner is a module or a service, not "everyone with the password".

Ownership includes:

- The tables or collections
- The invariants
- The public read contract (API or events)
- The backup and restore story (`db.topics.md`)
- The classification of the data (public, internal, personal)

A module owns a table when other modules do not write that table. They call the owner or they consume a published event. A foreign key or a shared identifier can still exist. A shared write must not exist.

Without ownership:

- No person can say which invariant is true.
- Schema change needs a meeting of all writers.
- Audit cannot name who changed a personal record.

Write an ownership map. Columns: data, owner, writers, readers, source of truth, copies. Fill it before you add a second database.

Law constraints sit on ownership (Topic 1). The owner of personal data must know retention and delete. A copy in a cache or an index still holds personal data. The owner must include those copies in the delete path.

Ownership is not the same as physical isolation. One database can hold many owners as schemas or as table prefixes. Physical isolation comes when a measure or a constraint requires it (Topic 9).

### Questions

#### Theoretical questions

1. What is data ownership in this handbook?
2. Who may write the source of truth?
3. What happens when a table has no owner?
4. How can one database still have many owners?
5. Why must a delete path include caches and indexes?

#### Easy practical tasks

1. Write five sentences that define data ownership.
2. Make a table: "Table" and "Owner". Add eight rows for a campus shop.
3. List four readers that must not write `orders`.
4. Write four questions you ask when you find a shared production password.

#### Medium practical tasks

1. Fill an ownership map for catalog, checkout, and identity in one monolith.
2. Write an ADR that forbids new writers on a legacy shared schema.
3. Plan a personal-data delete that includes the primary row, a cache key, and a search document.

#### Advanced practical tasks

1. Write a one-page ownership standard for a four-person team: how to add a table, how to add a reader, how to reject a shared write.
2. Estimate people hours for a year of schema changes in a shared-table design versus an owned-table design.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do table ownership in a monolith and "one database per service" relate as two steps?
2. When is CQRS a better answer than a second service with its own database?
3. How do source of truth, cache, and search index fail if you mix the words?
4. Why is an integration database a hidden API in the sense of Topic 5?
5. What Topic 1 constraints (law, team, cost) appear first in a store split?

#### Easy practical tasks

1. Write a one-page cheat sheet: shared versus owned store, CQRS trigger, cache versus source, ownership map columns.
2. For a to-do app, name one database, four owned tables, one optional cache, and no search product.
3. Draw the default: one monolith, one database, module-owned tables.
4. Bookmark `db.topics.md` and `redis.topics.md`. Write one sentence on when you open each file.

#### Medium practical tasks

1. Write a short data brief for a campus lost-and-found: owners, source of truth, one copy, and a rejected shared-SQL report.
2. Take a teammate design that reads production tables from a BI tool. Rewrite the read path as an export or an API.
3. Write a twelve-week plan: weeks for one database, and the measure that would reopen a second store.

#### Advanced practical tasks

1. Write a decision tree that a reviewer uses: ownership first, replica, cache, CQRS, then service store split.
2. Compare a managed analytics export with live SQL from a BI tool on the production primary. Write the risks.
