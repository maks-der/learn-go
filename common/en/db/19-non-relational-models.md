# 19. Non-Relational Models

## Description

This topic is a survey of data models that are not the relational table model. You learn key-value stores, document stores, wide-column stores, graph databases, search engines as stores, polyglot persistence, and the CAP trade-off in practical words.

Use one term for each concept. A data model is the way a product organizes values and queries. A system of record is the store that wins when copies disagree. Complete this topic after you can model tables and keys. You need the relational model as a comparison.

This path is vendor-neutral. Product names are examples of a category. Do not pick a product only by popularity.

---

## Key-value stores

A key-value store maps a key to a value. The primary operation is get and set by key. The store does not parse the value as a first-class query model. The value can be a string, a blob, or a small structure that the application encodes.

```text
SET session:42 --> { user_id: 7, expiry: ... }
GET session:42
DEL session:42
```

Use a key-value store when:

- you always know the key
- you need low latency
- the value is a cache, a session, a lock, or a flag
- you can accept eviction or a short TTL

Do not use a key-value store as the only store for orders with many query shapes ("all orders last month for this city"). You would scan keys or you would build secondary indexes by hand.

Keys need a naming scheme. Example: `user:7:session`. A collision overwrites a value. Treat the key space as a schema.

Many key-value products add extras: lists, expiry, pub/sub. Those extras do not make the product a relational DBMS. Transactions, if they exist, often cover one key or a small set of keys.

Durability differs. Some configurations keep data only in memory. A restart loses the data. Read the persistence setting before you store something that you cannot rebuild.

The relational comparison: a table with a primary key and no useful secondary access is close to key-value. The key-value product is then a simpler, faster specialist.

### Questions

#### Theoretical questions

1. What is the primary access path in a key-value store?
2. When is a key-value store a good fit?
3. Why is "all orders in a city" a poor key-value query?
4. What happens if two features use the same key string?
5. Why must you read the persistence setting before you store facts that you cannot rebuild?

#### Easy practical tasks

1. Design keys for: session, feature flag, page cache. Write three key strings.
2. Write get/set/delete for a shopping cart by `cart_id`.
3. Make a two-column table: fit versus poor fit. Add four rows.
4. Name one product in this category from public knowledge. Mark it as an example.

#### Medium practical tasks

1. Implement a tiny dictionary in your language that mimics get/set/TTL. Write what a real product adds (network, eviction).
2. Store a JSON blob under a key. Show that you cannot query a field without reading the blob.
3. Read a persistence page of one product. Write memory-only versus on-disk in three sentences.

#### Advanced practical tasks

1. Write a one-page design: sessions in a key-value store, orders in a relational database. Include expiry and rebuild.
2. Compare two key-value products at a high level: persistence, cluster, and value types. Use official docs.

---

## Document stores

A document store keeps a document, often JSON, as the unit of data. One document can nest arrays and objects. Documents in one collection can differ in fields (flexible schema).

```text
{
  "order_id": 1001,
  "customer": { "name": "Ada", "city": "Kyiv" },
  "lines": [
    { "sku": "A-10", "qty": 2 }
  ]
}
```

Use a document store when:

- the unit of read and write is a whole document
- the shape changes often
- you load one aggregate by id
- you accept that relations across documents are weaker than foreign keys

Query languages can filter on nested fields and can index those fields. That ability is stronger than a pure key-value store. It is not the same as arbitrary relational joins with constraints.

Integrity is often in the application. Unique indexes exist in many products. Multi-document transactions exist in some products and are newer or limited. Do not assume ACID across many documents unless you read the manual.

Duplication is common. You copy a customer name into the order document. An update of the customer name does not change old orders unless you write that job.

The relational comparison: a document is close to a row plus nested child rows in one JSON column. The document store makes that shape the default. Joins are not the default.

Do not store unbounded growing arrays in one document (all events of a user forever). Documents have size limits. Split when the aggregate grows without bound.

### Questions

#### Theoretical questions

1. What is the unit of data in a document store?
2. When is a document store a good fit?
3. Where do many integrity rules live if the store has no foreign keys?
4. Why does a copied customer name go stale?
5. Why must you not grow one document without a bound?

#### Easy practical tasks

1. Write a JSON document for a blog post with tags and an author object.
2. Write the same blog data as two relational tables. List two differences.
3. List three fields that you would index on a `orders` collection.
4. Name one product in this category. Mark it as an example.

#### Medium practical tasks

1. Design an order document versus 3NF tables. Write which queries are easier in each model.
2. Read a size-limit page for one document product. Write the limit and a split rule.
3. Describe a unique index on `email` in a document collection. Write one race that you still must handle.

#### Advanced practical tasks

1. Write a one-page model for a product catalog in documents. Include variants, price, and a change that would force a rewrite of all documents.
2. Compare multi-document transactions in one product with a relational transaction. Write five sentences from the manual.

---

## Wide-column stores

A wide-column store organizes data by a row key and a flexible set of columns, often grouped in column families. You do not load a fixed wide row of all columns on every read. You read the columns that you ask for.

Writes often append a new cell version. The row key determines the partition (the node that owns the row). Queries that match one row key are efficient. Queries that scan many keys can be expensive.

```text
Row key: user#42
  family profile: name=Ada, city=Kyiv
  family events:  2026-09-13T10:00=login, 2026-09-13T11:00=logout
```

Use a wide-column store when:

- you write a large volume of events or metrics
- you read by a known row key or a key range
- you can design the key for the query (query-first design)
- you accept eventual consistency options on some products

The primary-key design is the schema. If you need "all events for a user on a day," put user and day in the row key. If you need a different access path later, you add another table (another key design) and you write twice.

Secondary indexes exist in some products and have limits. Do not expect ad-hoc SQL joins.

The relational comparison: a table clustered by a composite primary key, with sparse columns, is the closest idea. The wide-column product targets huge scale and partition tolerance.

Do not pick this model for a small shop with many join-heavy reports. The operational cost is high. The query flexibility is low.

### Questions

#### Theoretical questions

1. What two parts identify a cell in a wide-column store?
2. Which queries are efficient?
3. What does query-first key design mean?
4. How do you add a new access path if the row key does not match?
5. Why is this model a poor default for a small join-heavy shop?

#### Easy practical tasks

1. Design a row key for "metrics for host H on day D."
2. Draw one row with two column families.
3. Write one query that fits the key and one query that does not.
4. Name one product in this category. Mark it as an example.

#### Medium practical tasks

1. Design two tables (two key layouts) for inbox-by-user and message-by-id. Write the double-write rule.
2. Read a partition or row-size warning in a manual. Write three sentences on hot partitions.
3. Compare a time-series row key with a relational `events` table plus indexes. Write four differences.

#### Advanced practical tasks

1. Write a one-page key design for a chat inbox. Include hot-user risk and a split rule.
2. Compare consistency settings of one wide-column product (one replica versus quorum). Write when a read can miss a write.

---

## Graph databases

A graph database stores nodes and edges (relationships). A node has a type and properties. An edge has a type, a direction, and properties. A query walks paths: neighbors, shortest path, or a pattern of types.

```text
(Ada:Person)-[:WORKS_AT]->(Acme:Company)
(Ada:Person)-[:KNOWS]->(Grace:Person)
```

Use a graph database when the question is the relationship:

- friends of friends
- who can reach this resource
- recommendation along a path
- fraud rings (shared devices and accounts)

A relational database can store edges as a join table. Deep or variable-length paths become recursive SQL or many joins. The graph product makes those paths the default query.

Integrity: some products constrain edge types. They do not replace a payroll ledger. Store money in a system of record. Use the graph for the link questions.

Do not model a simple one-to-many invoice as a graph only to follow a trend. A table is enough.

Indexes still exist on properties (find a node by email). After you have the start node, the walk uses edges.

Writes that add many edges need a plan. A supernode (a node with millions of edges) is a hot spot. Split or constrain that pattern.

### Questions

#### Theoretical questions

1. What are the two basic units in a graph database?
2. When is a graph model a good fit?
3. How do you store the same edges in a relational database?
4. Why is a graph a poor system of record for money?
5. What is a supernode problem?

#### Easy practical tasks

1. Draw five nodes and six edges for a classroom (people, courses).
2. Write one path question and one property lookup (find by email).
3. Write the same "enrolled in" fact as a join table.
4. Name one product in this category. Mark it as an example.

#### Medium practical tasks

1. Express "friends of friends who are not already friends" as a graph pattern and as SQL. Write which text is shorter.
2. Read a warning about dense nodes in a product manual. Write a mitigation.
3. Design a graph for file-access: user, group, file. Write one authorization path query in words.

#### Advanced practical tasks

1. Write a one-page design: relational users and payments, graph for recommendations. Include the sync path.
2. Compare a recursive SQL query with a graph walk for depth 4. Write cost risks from docs or from a small experiment.

---

## Search engines as data stores

A search engine indexes tokens and fields so that you can rank text. Users type words. The engine returns documents by relevance, filters, and facets.

```text
Index: title, body, tags, updated_at
Query: "integrity constraints" AND tag:db
Result: ranked hits, highlights
```

Use a search engine when users search text, when you need fuzzy match, or when you need relevance. Do not use a search index as the only copy of an order. Search indexes are often eventually consistent with the system of record. A reindex can rebuild from the record store. If you lose only the search node, you rebuild. If you lose only the record store, you lose the truth.

Writes go to the system of record first. Then you update the index (sync or async). A failed index update must be retryable. Search products often do not give the same transaction story as a relational DBMS.

Mapping (schema of the index) controls analysis: lowercase, stemming, language. A wrong analyzer makes queries miss. Treat the mapping as a schema.

Do not run `SELECT`-style joins across many indexes as a substitute for a warehouse. Use the engine for find-and-filter. Use a database for the ledger.

Access control must apply to hits. A search that ignores row grants can leak documents. Filter by tenant or permission in the query.

### Questions

#### Theoretical questions

1. What does a search engine optimize for?
2. Why is a search index a poor only copy of an order?
3. In which order do you write the record store and the index?
4. What does an analyzer change?
5. Why must a search query include a permission filter?

#### Easy practical tasks

1. Write three user queries that need search, not `LIKE '%x%'`.
2. Draw write path: app, DBMS, index.
3. List four fields on a blog post index.
4. Name one search product. Mark it as an example.

#### Medium practical tasks

1. Compare `LIKE` and a token index for a 1 million row article table. Write four differences.
2. Read a refresh or refresh-interval setting. Write how soon a new document can appear in results.
3. Design a tenant filter on every query. Write the field and the risk if you omit it.

#### Advanced practical tasks

1. Write a one-page dual-write design: relational post, search document, retry, rebuild from SQL.
2. Plan a mapping change that requires reindex. Write downtime or alias-switch steps from a product guide.

---

## Polyglot persistence: one app, more than one store

Polyglot persistence means one application uses more than one data model on purpose. Example: relational orders, key-value sessions, search for the catalog, a graph for recommendations.

Each store has a job. One store is the system of record for each fact. Other stores are derived.

```text
Order placed --> relational DBMS (record)
             --> search index (catalog availability text)
             --> cache invalidate (session cart)
```

Benefits: each query uses a model that fits. Costs: more operations, more failure modes, more consistency work. You must define the write order and the repair job.

Rules:

1. Name the system of record for each fact.
2. Do not write a price in three stores as three equals without a sync rule.
3. Start with one store. Add a second store when a measured need exists (search latency, session volume).
4. Keep transactions in the record store. Use outbox or a retry job for derived stores (later topic).

Do not add a store because a blog post praised it. Do not run five products for a class project.

A single relational DBMS can cover a large part of a small system. Polyglot is an architecture choice under load or under a special query, not a badge.

### Questions

#### Theoretical questions

1. What does polyglot persistence mean?
2. What is a system of record in this section?
3. Name two costs of extra stores.
4. When do you add a second store?
5. Where do transactions live in a polyglot design?

#### Easy practical tasks

1. Assign stores to: invoice, session, full-text help pages.
2. Draw one user action that touches two stores.
3. Write one fact that must have a single record store.
4. List three operational tasks that double when you add a DBMS product.

#### Medium practical tasks

1. Design a small shop with two stores only. Justify why a third store is not needed yet.
2. Write a sync failure: record committed, search not updated. Write the user-visible symptom and a repair.
3. Compare "JSON in PostgreSQL" versus a document DBMS for one collection. Write when one product is enough.

#### Advanced practical tasks

1. Write a one-page polyglot map: facts, record store, derived stores, write order, repair.
2. Review a public architecture that uses three stores. Mark each as record or derived. Note any unclear owner.

---

## CAP theorem in practical words (trade-offs, not a slogan)

CAP is a way to talk about a partition (a network split) in a distributed store.

- **C (consistency)** here means all clients see the same latest write (linearizability, in informal words: one up-to-date copy).
- **A (availability)** means every request to a non-failing node receives a response (not an error that says "I cannot decide").
- **P (partition tolerance)** means the system continues while the network drops messages between nodes.

On a single node, you do not choose CAP. When the network between nodes fails, a product must choose: refuse some writes or reads (keep one answer), or keep serving and accept that copies diverge.

```text
Partition: site A cannot talk to site B
  Option 1: A stops writes so B cannot miss them  (prefer one copy)
  Option 2: A and B both accept writes            (prefer availability, diverge)
```

Practical words:

- A relational primary that refuses writes when it cannot reach a sync replica prefers one copy.
- A multi-primary that accepts writes on both sides of a split prefers availability. You merge later.
- "We are AP" or "we are CP" as a slogan hides the real setting (quorum size, retry, client failover).

CAP does not say that you drop durability or that SQL is obsolete. It does not apply to a single SQLite file on one disk in the usual way.

When you read a vendor page, ask: what happens when two sites cannot talk? Who accepts writes? How do clients find a node? What do you repair after the split?

Do not use CAP to end a design talk. Use it to start the partition story.

### Questions

#### Theoretical questions

1. What does partition mean in this section?
2. What do you give up if both sides accept writes during a split?
3. What do you give up if a primary refuses writes during a split?
4. Why is "we are AP" a weak requirement?
5. Why does CAP not describe a single-node embedded file in the usual way?

#### Easy practical tasks

1. Write C, A, and P in one short sentence each (practical words).
2. Label option 1 and option 2 from this section on a two-site drawing.
3. Write one product setting that is closer to "one copy" and one that is closer to "keep serving."
4. Write one sentence that you will not use as a slogan.

#### Medium practical tasks

1. Take a primary plus async replica. Write what clients see if the primary is in the minority partition.
2. Take a quorum of three nodes. Write how a majority still accepts a write.
3. Read a CAP paragraph in a vendor doc. Rewrite it as a partition story in six sentences.

#### Advanced practical tasks

1. Write a one-page partition table-top for a shop: two sites, who writes, what users see, how you merge.
2. Compare two stores (relational HA versus a distributed key-value) on one partition. Use manuals. No slogans.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you choose key-value, document, wide-column, graph, or search from the access path, not from the product name?
2. When does polyglot persistence need a system of record, and how does CAP describe a split between two of those stores?
3. Which models treat secondary access as a new table or index that you design up front?
4. What integrity jobs stay in the application when you leave the relational model?
5. A teammate wants one document store for invoices, search, sessions, and friend graphs. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: five models, polyglot, CAP partition story.
2. Fill a table: model, unit of data, typical query, poor query.
3. Assign one example product per model. Mark all as examples.
4. Pick a system of record for a library loan. Name one derived store if you need search.

#### Medium practical tasks

1. Split a social shop across two models. Write write order and one repair.
2. Write a partition story for your split. State who accepts checkout writes.
3. Rebuild one aggregate (order) as tables and as a document. List queries that break in each form.

#### Advanced practical tasks

1. Write an architecture note: three stores, record versus derived, operations cost, and a CAP table-top.
2. Take a public system design. Relabel each database as a model from this topic. Mark slogans and replace them with partition behavior.
