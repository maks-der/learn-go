# 9. Data Architecture

## Description

Data architecture is the set of decisions about where data lives, who owns it, and how other parts read it. This topic covers one database per service versus a shared database, the integration-database anti-pattern, CQRS, caching, search indexes versus the source of truth, and data ownership.

Data decisions are costly to change (Topic 1). A wrong share of tables creates hidden coupling (Topic 3). Pair this topic with `db.topics.md`. For cache details, pair with `redis.topics.md`.

Complete Topics 1 to 8 before this topic. Start with one database in a modular monolith (Topic 5). Split stores only when ownership or a measure requires it.

Use one term for each concept. The source of truth is the store that wins when copies disagree. A cache is not the source of truth. A search index is not the source of truth unless you design it that way (rare).

---

## One database per service vs shared database

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
- Operations cost grows (Topic 2: cost).

Do not copy the "one database per service" slogan into a three-person monolith. You already have one deployable unit. Use table ownership first.

A shared read-only replica of the same owned database is not the same as a shared integration database. Replicas are a scale tool. They still have one owner.

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
2. Estimate people hours for a year of schema changes in a shared-table design versus an owned-store design.

---

## Integration database anti-pattern

An integration database is a database that exists so that many programs can integrate by reading and writing the same tables. The database becomes the contract. This handbook treats that design as an anti-pattern for application integration.

Why it hurts:

- The schema cannot evolve. Every program is a client of every column.
- Invariants have no single owner.
- A "quick" SQL join in a report binds you to internals.
- Security becomes "who has the database password".
- You cannot test one program without the full schema and the full data meaning.

Preferred integration styles (later path, Topic 13): API, messaging, or file transfer with a documented format. Those styles have explicit contracts.

A data warehouse or an analytics replica can be a valid copy for reports. It is still a copy. The owner of the source of truth remains the operational system. Reports must not write back into operational tables without a defined process.

If you inherit an integration database, do not add a fifth writer. Add an API in front of the owner. Move writers one by one. Record the password list and reduce it.

Legal and audit constraints often appear here. Many writers make it hard to say who changed a personal record.

### Questions

#### Theoretical questions

1. What is an integration database?
2. Why is the schema a bad integration contract?
3. Why is a database password a weak security boundary?
4. What integration styles are better for applications?
5. How can a warehouse copy be valid without becoming an integration database?

#### Easy practical tasks

1. Write four sentences that define the anti-pattern.
2. List five programs that might share a campus database in a mud design.
3. Make a table: "Integration method" and "Contract type". Add tables, HTTP API, and events.
4. Write three questions you ask when you find a shared production password.

#### Medium practical tasks

1. Plan a move from a shared report SQL to an API. Write the report fields and the owner.
2. Write an ADR that forbids new writers on a legacy shared schema.
3. Draw the current "all apps to one DB" picture and the target "owner plus copies".

#### Advanced practical tasks

1. Write a two-page recovery plan for a three-app integration database: inventory, owner, API, sunset of SQL access.
2. Compare a managed analytics export (batch files) with live SQL from a BI tool on the production primary. Write the risks.

---

## CQRS (when reads and writes diverge)

CQRS (Command Query Responsibility Segregation) means that the model that you use to change data is not the same model that you use to read data. Commands (writes) go to one side. Queries (reads) go to another side.

A light form is common and enough: write to normalized tables, read from a SQL view, a query tailored table, or a cache. The write model stays strict. The read model is convenient.

A strong form uses a different store for reads (a search index, a document, a cube). A projector updates the read store when writes occur. Consistency becomes eventual (Topic 7).

Use CQRS when:

- Read shapes and write shapes differ a lot.
- Read load is much higher than write load.
- A single model forces painful joins on every page.

Do not use full CQRS when a few SQL queries are enough. A second store is a second failure domain. You must handle lag, rebuilds, and dual bugs.

Commands still need validation and invariants. CQRS is not a reason to skip the write model. Queries must not sneak writes.

Name the lag. "The admin dashboard can lag 2 seconds" is a requirement. "The checkout total can lag 2 seconds" is often not acceptable.

### Questions

#### Theoretical questions

1. What is CQRS?
2. What is a light form of CQRS?
3. What is a strong form of CQRS?
4. When do you use CQRS?
5. Why must you name the read lag?

#### Easy practical tasks

1. Write five sentences that define command side and query side.
2. Make a table: "Screen" and "Same model as write? (yes/no)". Add six screens of a shop.
3. List three extra costs of a second read store.
4. Write a lag requirement for a search page and for a payment receipt.

#### Medium practical tasks

1. Design a light CQRS: write tables plus a read view for an order history page.
2. Write an ADR that rejects a second database for reads on a 100-user app.
3. Draw a projector: write commit → event → update of a read table. Mark failure points.

#### Advanced practical tasks

1. Write a one-page rebuild plan: the read store is wiped. How do you fill it from the source of truth?
2. Compare CQRS with "one model plus better indexes". Write when indexes win.

---

## Caching layers (`redis.topics.md`)

A cache is a store of copies that you can rebuild from the source of truth. The goal is lower latency or lower load on the primary store. See `redis.topics.md` for Redis as a cache.

A cache is optional. The system must still work (perhaps more slowly) if the cache is empty or down, unless you explicitly design a degraded mode.

Common layers:

- Client cache (browser).
- CDN for public files.
- Application memory cache (process-local).
- Shared cache (Redis or similar).
- Database buffer cache (the database already has one).

Invalidation is the hard part. After a write, copies can be stale. Strategies include time-to-live (TTL), delete-on-write, and version keys. Each strategy can fail. TTL is simple and allows staleness. Delete-on-write can miss a key name.

Do not cache data that must be exact at read time unless the TTL is zero or you invalidate correctly. Money balances and authorization decisions need care.

Local process caches do not share deletes. Two instances can disagree until TTL ends. A shared cache reduces that problem and adds a network hop.

Write what you cache, the key format, the TTL, and the stale behavior. Measure hit rate. A cache with a low hit rate still adds complexity.

Do not use the cache as the only store of a write. That design loses data when the cache evicts a key.

### Questions

#### Theoretical questions

1. What is a cache in this handbook?
2. Why must the system survive an empty cache?
3. Name four cache layers.
4. Why is invalidation hard?
5. Why must a cache not be the only store of a write?

#### Easy practical tasks

1. Write four sentences that define cache versus source of truth.
2. Make a table: "Data" and "Cache? (yes/no)". Add public catalog, cart, and payment capture.
3. List three invalidation strategies in one sentence each.
4. Write a key format for `book:{id}` and a TTL choice.

#### Medium practical tasks

1. Design a cache for `GET /books/{id}` with delete-on-write. Write the miss path.
2. Explain in eight sentences why two instances with memory caches can show different titles after an edit.
3. Write an ADR: Redis cache for reads, PostgreSQL as source of truth. Point to `redis.topics.md`.

#### Advanced practical tasks

1. Write a one-page invalidation map for five read endpoints and their write events.
2. Plan a cache failure mode: Redis down. Write timeouts (Topic 4) and user-visible behavior.

---

## Search index vs source of truth

A search index is a structure that is optimized for text search and faceted queries. Elasticsearch and similar products are indexes. A relational table with a full-text feature is still the source of truth if it is the owner. A separate search product is usually a derived copy.

The source of truth accepts writes and holds the authoritative record. The search index receives a projection. Search can lag. Search can drop a document. Search ranking can omit a row that exists in the source.

User rules:

- Create, update, and delete go to the source of truth.
- Search queries can go to the index.
- Open-by-id after a click must be able to load from the source if the index is wrong.

Do not write only to the index. A rebuild would invent data that you never stored.

Reindex is a required operation. You will need to fill the index from the source after a bug or a mapping change. If you cannot rebuild, the index became an accidental source of truth.

Consistency is eventual. Write the lag that search may have. For a shop, a new item that appears in search after 10 seconds can be acceptable. A sold-out item that stays "buyable" in search is a business risk. Combine search hits with a stock check on the source at purchase time.

### Questions

#### Theoretical questions

1. What is a search index in this handbook?
2. What is the source of truth?
3. Why can search omit a row that exists?
4. Why must you be able to rebuild the index?
5. Why must purchase still check the source of truth?

#### Easy practical tasks

1. Write five sentences that compare search and source of truth.
2. Make a table: "Action" and "Which store". Add type-ahead search, open item, and update title.
3. List four fields that a book index might hold and two that must stay only in the source.
4. Write a lag requirement for catalog search.

#### Medium practical tasks

1. Design the write path: save book → project to index. Mark failures and retries.
2. Write the purchase path: search hit → load book from source → check stock.
3. Write an ADR that adds search only when measured query time on SQL fails the target.

#### Advanced practical tasks

1. Write a rebuild runbook: dump from source, index mapping, batch size, and how you detect drift.
2. Compare SQL full-text search with a separate index for a 50 000-document corpus. Write a choice for a small team.

---

## Data ownership

Data ownership means that one module or one service is responsible for a set of data. That owner is the only writer (except controlled migrations). Other parts read through the owner’s API, events, or an approved copy.

Ownership includes:

- Schema change rights
- Invariants
- Retention and deletion (law)
- Access rules
- Backup restore of that data

Without ownership, you have an integration database or a mud (Topic 5). With ownership, you can evolve.

Shared reference data (country codes) can have a small owner or a static package. Do not let every service invent a country table.

Copies must have a status: source, cache, search, warehouse. Each copy has a refresh rule. Each copy has a ban on silent writes.

Ownership is also a stakeholder fact (Topic 1). The owner accepts the quality of that data. A report team that does not own orders must not "fix" order rows in place.

In a monolith, write the owner in the module and in a table comment or a schema file header. In a multi-service system, write the owner in the ADR and in the API.

Deletion and privacy follow the owner. If you copy personal data into three stores, the owner must know all copies. Search and cache are copies.

### Questions

#### Theoretical questions

1. What is data ownership?
2. What duties does an owner have?
3. How do other modules read owned data?
4. Why must copies have a status and a refresh rule?
5. Why do privacy deletions need a list of copies?

#### Easy practical tasks

1. Write an ownership table for eight entities of a clinic (patient, slot, invoice).
2. List four read methods that do not violate ownership.
3. Make a table: "Store" and "Source or copy". Add primary SQL, Redis, and search.
4. Write four sentences on who may `UPDATE` a loan row.

#### Medium practical tasks

1. Assign owners in a modular monolith schema of ten tables. Mark two tables that are currently shared by mistake.
2. Write a deletion sequence for a user: source first, then cache, then search, then warehouse.
3. Write an ADR: module `billing` owns invoices; other modules use events.

#### Advanced practical tasks

1. Write a one-page ownership register: entity, owner, store, copies, retention, and API.
2. Design a quarterly audit that finds tables with more than one writer (logs or privileges). Stay defensive.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do ownership and the integration-database anti-pattern express the same rule?
2. When do CQRS, cache, and search all become "copies with lag"?
3. Why does a modular monolith delay "one database per service"?
4. How do Topics 7 and 9 meet at consistency and aggregates?
5. Which quality attributes suffer first when every app writes the same production database?

#### Easy practical tasks

1. Write a one-page cheat sheet: source of truth, copies, ownership, CQRS, cache, search.
2. For a campus lost-and-found, draw one database, owners, and one optional cache.
3. Write three ADR titles: one store, no integration SQL, search later.
4. Bookmark `db.topics.md` and `redis.topics.md`. Write one task from each that you will do next.

#### Medium practical tasks

1. Write a two-page data architecture for a library: owners, copies, read models, and a search decision.
2. Take a CRUD app. Mark which endpoints are writes to the source and which could use a cache.
3. Role-play a request for a shared production login to a BI tool. Write the safer export path.

#### Advanced practical tasks

1. Write a data architecture standard (three pages) for a team of five: ownership, copies, and forbidden integrations.
2. Design a drift detector: compare counts between source and index daily. Write the alert and the rebuild trigger.
