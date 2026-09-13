# 6. Indexes

## Description

An index is a data structure that helps the server find documents without a full collection scan. This topic shows single-field indexes, compound indexes, unique indexes, multikey indexes, TTL indexes, and a high-level view of text, geospatial, and hashed indexes. You also hide indexes and you read `explain`. Complete querying first.

Use one term for each concept. An **index** is a sorted structure on one or more fields. A **collection scan** (`COLLSCAN`) reads every document. An **index scan** (`IXSCAN`) reads the index, then fetches documents if needed. A **prefix** of a compound index is the leftmost field list.

Indexes speed reads. Indexes slow writes. Each insert, update, and delete must update the indexes. Do not add an index for every field.

---

## Single-field indexes

A single-field index uses one field. MongoDB already has a unique index on `_id`.

Create an index in `mongosh`:

```javascript
db.orders.createIndex({ userId: 1 })
```

`1` is ascending. `-1` is descending. For a single field, direction rarely changes equality lookups. Direction matters more for sort.

Use a single-field index when many queries filter or sort on that field alone.

List indexes:

```javascript
db.orders.getIndexes()
```

Drop an index by name:

```javascript
db.orders.dropIndex("userId_1")
```

An index has a name. MongoDB can generate the name. You can set `name` in `createIndex` options.

A query uses an index when the filter or sort matches the index. `explain` shows the plan. If you see `COLLSCAN` on a large collection, you likely need an index or a better filter.

Do not create a duplicate index on the same key pattern. The create command fails or is a no-op depending on options.

### Questions

#### Theoretical questions

1. What problem does an index solve?
2. Which index exists on every collection by default?
3. What do `1` and `-1` mean in `createIndex`?
4. Why do extra indexes slow writes?
5. What does `COLLSCAN` mean?

#### Easy practical tasks

1. Create an index on `{ userId: 1 }`. Run `getIndexes()`. Write the name.
2. Insert five orders with `userId`. Find `{ userId: "u1" }`. Run `explain`. Write the stage.
3. Drop the index. Run `explain` again. Write the stage.
4. Open the official index introduction. Write the URL.

#### Medium practical tasks

1. Time a find on a field before and after you add an index. Use at least 10 000 documents. Write the two times.
2. Create an index with a custom `name`. Drop it by that name.
3. Compare `explain("queryPlanner")` and `explain("executionStats")`. Write one extra field that execution stats show.

#### Advanced practical tasks

1. Read about `background` index builds in older versions vs current index builds. Write the current default behavior for your version.
2. Create an index that is not used by any query. Write how you detect that with `$indexStats` or Atlas.

---

## Compound indexes and prefix rule

A **compound index** uses two or more fields in a fixed order.

```javascript
db.orders.createIndex({ userId: 1, createdAt: -1 })
```

The **ESR** guide (Equality, Sort, Range) helps you order keys: fields you match for equality, then fields you sort, then fields you filter as a range.

The **prefix rule**: a compound index on `{ a: 1, b: 1, c: 1 }` can support queries on:

- `{ a }`
- `{ a, b }`
- `{ a, b, c }`

It does not support a query on `{ b }` or `{ c }` alone as an index prefix. The server may still use a later key in some plans, but you must not rely on that. Put the first filter field first when you can.

Sort: an index can satisfy a sort if the sort keys match the index order (or the reverse, with rules). A mismatch can force an in-memory sort.

One compound index can replace two single-field indexes when queries always include the left prefix. Do not create both `{ a: 1, b: 1 }` and `{ a: 1 }` without a reason. The prefix covers `{ a }`.

Order is part of the index. `{ a: 1, b: 1 }` is not `{ b: 1, a: 1 }`.

### Questions

#### Theoretical questions

1. What is a compound index?
2. What is the prefix rule?
3. Why does `{ a: 1, b: 1 }` not replace an index for queries on `b` only?
4. What does ESR mean?
5. When can one compound index replace a single-field index on the first field?

#### Easy practical tasks

1. Create `{ userId: 1, createdAt: -1 }`. List indexes.
2. Find `{ userId: "u1" }` and sort `{ createdAt: -1 }`. Run `explain`. Write if the plan uses the index.
3. Find `{ createdAt: { $gt: ISODate("2020-01-01") } }` only. Run `explain`. Write the stage.
4. Draw the three prefixes of `{ a: 1, b: 1, c: 1 }`.

#### Medium practical tasks

1. Compare `{ userId: 1, createdAt: -1 }` and `{ createdAt: -1, userId: 1 }` for a user-history query. Write which prefix matches.
2. Create two indexes that overlap (`{ a: 1 }` and `{ a: 1, b: 1 }`). Write whether you need both.
3. Force a sort that does not match the index. Show `SORT` or a memory sort in `explain`.

#### Advanced practical tasks

1. Read the official compound-index and sort pages. Write one rule about sort direction and index direction.
2. Design indexes for three queries on `orders`. Use as few compound indexes as you can. Write the key patterns and which query each serves.

---

## Unique indexes

A **unique index** rejects a second document with the same key.

```javascript
db.users.createIndex({ email: 1 }, { unique: true })
```

The `_id` index is unique.

Unique indexes enforce uniqueness that the application also wants. The database is the last guard. Two concurrent inserts cannot both succeed with the same email.

A unique index treats missing fields as `null` (with rules). Two documents that omit `email` can conflict on a unique `email` index. Use a **partial unique index** when you want uniqueness only when the field exists:

```javascript
db.users.createIndex(
  { email: 1 },
  { unique: true, partialFilterExpression: { email: { $exists: true, $type: "string" } } }
)
```

A unique compound index is unique on the combination of fields.

Do not use a unique index as your only validation of format. Unique is about equality of the key, not about email syntax.

If a duplicate exists, `createIndex` fails. Clean the data first.

### Questions

#### Theoretical questions

1. What does a unique index prevent?
2. Why is a unique index useful under concurrent inserts?
3. How can two documents without `email` conflict on a unique `email` index?
4. What is a partial unique index?
5. What is unique in a unique compound index?

#### Easy practical tasks

1. Create a unique index on `email`. Insert one user. Insert the same email again. Record the error.
2. Run `getIndexes()`. Write the `unique` flag.
3. Create a unique compound index on `{ tenantId: 1, sku: 1 }`. Show that the same sku can exist in two tenants.
4. Find the unique-index page in the manual. Write the URL.

#### Medium practical tasks

1. Insert two documents that omit `email` after a unique index on `email`. Record the result. Then try a partial unique index.
2. Build a unique index on a collection that already has a duplicate. Record the error. Remove the duplicate. Create the index.
3. Compare application check vs unique index. Write a race that only the index catches.

#### Advanced practical tasks

1. Read about unique indexes and `null` in your server version. Write the exact rule that you tested.
2. Design uniqueness for a sparse optional username. Write the index options and two test inserts.

---

## Multikey indexes (arrays)

When you index a field that holds an array, MongoDB creates a **multikey** index. Each array element becomes a key in the index.

```javascript
db.posts.createIndex({ tags: 1 })
```

A query `{ tags: "mongo" }` can use that index.

Multikey indexes have limits. You cannot create a compound index that has **more than one** array field in the general case. An index `{ tags: 1, scores: 1 }` fails if both fields are arrays.

A compound index can include one array field and other scalar fields: `{ authorId: 1, tags: 1 }`.

Embedded documents in an array: you can index `"lines.sku"`. The index is multikey.

Multikey indexes can be large. Each element is a key. A document with 10 000 tags creates many index entries. Keep arrays bounded.

`explain` can show `isMultiKey: true`.

### Questions

#### Theoretical questions

1. What is a multikey index?
2. Why can `{ tags: "mongo" }` use an index on `tags`?
3. Why can you not compound-index two array fields in the usual case?
4. Why do large arrays make a multikey index expensive?
5. What does `isMultiKey` mean in `explain`?

#### Easy practical tasks

1. Index `tags`. Insert a post with three tags. Find one tag. Run `explain`.
2. Run `getIndexes()`. Write whether the index is marked multikey after the insert.
3. Index `"lines.sku"` on orders. Query `{ "lines.sku": "nail" }`.
4. Try to create `{ tags: 1, aliases: 1 }` on documents where both are arrays. Record the result.

#### Medium practical tasks

1. Measure index size with `collStats` before and after you add 1000 tags across many documents. Write the sizes.
2. Compare a query `{ tags: { $all: ["a", "b"] } }` with `explain`. Write how the index is used.
3. Show a compound `{ userId: 1, tags: 1 }`. Query both fields. Write the plan.

#### Advanced practical tasks

1. Read the multikey restrictions page. Write two restrictions in your own words.
2. Compare a tags array plus multikey index vs a `post_tags` collection. Write when the extra collection wins.

---

## TTL indexes

A **TTL index** (Time To Live) deletes documents after a delay. You create it on a **date** field.

```javascript
db.sessions.createIndex({ expiresAt: 1 }, { expireAfterSeconds: 0 })
```

If `expireAfterSeconds` is `0`, MongoDB deletes the document when `expiresAt` is in the past. If you set a positive number, MongoDB adds that many seconds to the field value.

The TTL monitor runs in the background. Deletes are not instant at the exact second. Do not use TTL for a lock that must expire at a precise instant.

The field must be a Date (or an array of Dates). A string date does not expire.

TTL is one index per special use. You cannot use a compound TTL index to expire on two fields in the old model; read the current manual for allowed forms. The common pattern is one date field.

Use TTL for sessions, logs, and temporary tokens. Do not use TTL as the only legal-hold policy if you must keep data.

### Questions

#### Theoretical questions

1. What does a TTL index do?
2. What type must the indexed field have?
3. Why is deletion not exact to the second?
4. What does `expireAfterSeconds: 0` mean when the field is an absolute expire date?
5. When must you not use TTL as the only retention control?

#### Easy practical tasks

1. Create a TTL index on `expiresAt` with `expireAfterSeconds: 0`. Insert a document with `expiresAt` one minute in the past. Wait and confirm the delete (or read how to check the monitor).
2. Insert a document with a string in `expiresAt`. Confirm it does not expire.
3. Run `getIndexes()`. Write the `expireAfterSeconds` option.
4. Open the TTL page in the manual. Write the URL.

#### Medium practical tasks

1. Use `expireAfterSeconds: 60` on `createdAt`. Insert now. Write when you expect the delete.
2. Explain why a TTL index is also an index that queries can use. Run a find on `expiresAt` with `explain`.
3. Compare TTL delete with an application job that runs `deleteMany`. Write two differences.

#### Advanced practical tasks

1. Read about the TTL monitor interval and replica set behavior. Write who runs the delete.
2. Design a collection that keeps 7 days of events. Write the date field, the index, and one query that still needs a non-TTL index.

---

## Text, geospatial, hashed (high-level)

These index types solve special queries. Learn the idea. Read the full pages before production use.

**Text index.** A text index supports `$text` search on string fields. It tokenizes words. It is not a full search engine. Atlas Search (Lucene) is stronger for product search. One collection can have at most one text index (legacy rule; confirm for your version).

**Geospatial index.** `2dsphere` indexes GeoJSON points and shapes. You query `$near` and `$geoWithin`. Use it for "find places near this point". Store coordinates in the GeoJSON form that the manual shows. A bad coordinate format will not use the index.

**Hashed index.** A hashed index stores a hash of the field. MongoDB uses hashed indexes for hashed sharding. Equality queries can use a hashed index. Range queries cannot use a hashed index well.

Do not create a text index to replace a simple equality index on `email`. Do not create a geospatial index on a normal number field. Do not create a hashed index for a range sort.

### Questions

#### Theoretical questions

1. What query operator uses a text index?
2. Why is Atlas Search often better than a text index for product search?
3. What is a `2dsphere` index for?
4. What is a hashed index for in a cluster?
5. Why can a hashed index not support a range query well?

#### Easy practical tasks

1. Create a text index on `title`. Insert two titles. Run a `$text` query. Write the matches.
2. Read the GeoJSON point example in the manual. Write one valid point document.
3. Read the hashed index page. Write one sentence about sharding.
4. Make a table: index type, typical query, one limit.

#### Medium practical tasks

1. Compare `$text` and `$regex` on the same field. Write which uses the text index.
2. Create a `2dsphere` index. Run a `$near` query from the manual example (use a tiny data set).
3. Explain why `{ email: 1 }` is better than a text index for login lookup.

#### Advanced practical tasks

1. Read Atlas Search vs text index. Write three features that Search adds.
2. Read hashed vs ranged shard keys (preview of sharding). Write one workload that fits hashed.

---

## Hidden indexes and `explain`

A **hidden index** exists and the server updates it, but the query planner does not use it.

```javascript
db.orders.createIndex({ status: 1 }, { hidden: true })
db.orders.hideIndex("status_1")
db.orders.unhideIndex("status_1")
```

Hide an index before you drop it. If the load is fine, drop it. If a query becomes slow, unhide it. This is safer than an immediate drop.

`explain` shows how the server runs a query.

```javascript
db.orders.find({ userId: "u1" }).explain("executionStats")
```

Verbosity:

- `queryPlanner` — the selected plan
- `executionStats` — documents examined and time
- `allPlansExecution` — more than one plan

Read these fields:

- `stage` — `COLLSCAN`, `IXSCAN`, `FETCH`, `SORT`
- `nReturned` — documents returned
- `totalDocsExamined` — documents read
- `totalKeysExamined` — index keys read

A good plan examines about as many documents as it returns, or fewer with a covered query. A bad plan examines far more documents than it returns.

Do not guess. Run `explain` on a realistic data size.

### Questions

#### Theoretical questions

1. What does a hidden index do to writes and to the planner?
2. Why do you hide an index before you drop it?
3. What is the difference between `queryPlanner` and `executionStats`?
4. What does `totalDocsExamined` much larger than `nReturned` suggest?
5. What is a `FETCH` stage?

#### Easy practical tasks

1. Create an index. Hide it. Run `explain` on a query that used it. Write the new stage.
2. Unhide the index. Run `explain` again. Write the stage.
3. Run `explain("executionStats")` on a `COLLSCAN` query. Write `totalDocsExamined`.
4. Find the hidden-index page. Write the URL.

#### Medium practical tasks

1. Compare `IXSCAN` + `FETCH` with a projection that might cover the query. Write whether `FETCH` remains.
2. Use `allPlansExecution` on a query with two possible indexes. Write which plan won.
3. Hide an index in a test. Measure one slow query. Unhide. Measure again.

#### Advanced practical tasks

1. Read `$indexStats`. Write how you see operations that used each index.
2. Write a drop checklist: hide, watch metrics, explain slow queries, then drop. Run it on a unused test index.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do single-field, compound, unique, and multikey indexes differ in what they store?
2. Which index types in this topic change data (TTL) or change planner behavior (hidden) without a drop?
3. Why does the prefix rule matter when you add a second field to an index?
4. How do you prove that a query uses an index?
5. A teammate adds an index on every field. Which facts do you use to stop that?

#### Easy practical tasks

1. On `orders`, create `{ userId: 1, createdAt: -1 }` and a unique index on `{ orderNo: 1 }`. Write `getIndexes()`.
2. Write a cheat sheet: createIndex, prefix rule, unique, multikey, TTL, text/geo/hashed, hide, explain stages.
3. Run `explain("executionStats")` on one equality query. Save `stage` and `totalDocsExamined`.
4. Drop a test index after you hide it first. Write the two commands.

#### Medium practical tasks

1. Design indexes for: find orders by user sorted by date, unique email, tags search, session expiry. Write four `createIndex` calls.
2. Load 20 000 documents. Compare `COLLSCAN` and `IXSCAN` times for the same filter. Write the two `explain` summaries.
3. Show one query that a compound index serves and one query on the right field that it does not serve.

#### Advanced practical tasks

1. Use `$indexStats` and `explain` to find one unused index in a test database. Hide it, then drop it.
2. Read about covered queries. Build one find that is covered. Write the index, the projection, and the `explain` stage list.
