# 8. Schema Design Patterns

## Description

This topic shows common MongoDB schema patterns. A pattern is a known way to shape documents. You will learn the attribute, bucket, outlier, computed, subset, and extended-reference patterns, and a schema version field. Complete data modeling first.

Use one term for each concept. A **pattern** is a design that you apply on purpose. Do not apply every pattern to every collection. Name the pattern in your notes when you use it.

Official pattern names come from MongoDB documentation and talks. Use those names so that a team can talk with one vocabulary.

---

## Attribute pattern

The **attribute pattern** stores many optional predicates as an array of key-value documents.

Instead of many sparse fields:

```javascript
{ color: "red", size: "M", material: null, season: null }
```

you store:

```javascript
{
  attrs: [
    { k: "color", v: "red" },
    { k: "size", v: "M" }
  ]
}
```

You index `{ "attrs.k": 1, "attrs.v": 1 }`. A query for color red uses `$elemMatch` on `k` and `v`.

Use this pattern when:

- Many documents have different optional fields
- You must index those fields
- A wide set of field names would create many sparse indexes

Do not use this pattern for a small, stable set of fields. Normal fields are easier to read and to validate.

The application still owns the allowed keys. The array is not a dump of random JSON without a list of keys.

Reads that need a document shape with named fields can `$project` or reshape in the application.

### Questions

#### Theoretical questions

1. What problem does the attribute pattern solve?
2. Why is a pair `k` and `v` easier to index than 50 sparse fields?
3. When must you use `$elemMatch` on `attrs`?
4. Why is this pattern a poor fit for five stable fields?
5. Who lists the allowed keys?

#### Easy practical tasks

1. Store two products with different `attrs`. Query color red with `$elemMatch`.
2. Create the compound index on `attrs.k` and `attrs.v`. Run `explain`.
3. Draw the sparse-field shape and the attribute-array shape for the same product.
4. Open the official attribute-pattern page if it exists. Write the URL or the search title that you used.

#### Medium practical tasks

1. Compare three single-field indexes vs one attribute index for ten optional fields. Write index count.
2. Write a validator idea: `attrs.k` must be in an enum. List five keys.
3. Reshape `attrs` to a map in the application. Write the function steps.

#### Advanced practical tasks

1. Read about the attribute pattern and wildcard indexes. Write when a wildcard index is an alternative.
2. Design a catalog with 200 possible predicates. Write why attributes beat 200 fields.

---

## Bucket pattern

The **bucket pattern** groups many small events into one document. The document is a **bucket** for a time range or a size limit.

Example: sensor readings for one device and one hour:

```javascript
{
  deviceId: "d1",
  start: ISODate("2026-01-15T10:00:00Z"),
  n: 3,
  readings: [
    { t: ISODate("2026-01-15T10:01:00Z"), c: 21.1 },
    { t: ISODate("2026-01-15T10:02:00Z"), c: 21.2 },
    { t: ISODate("2026-01-15T10:03:00Z"), c: 21.0 }
  ]
}
```

Benefits:

- Fewer documents
- Smaller `_id` index
- One read for a window

You must bound the bucket. Close the bucket when `n` reaches a maximum or when the time window ends. Start a new document.

MongoDB time series collections use a related idea inside the server. You can still use an application-level bucket for other data (logs, metrics).

Do not put an unbounded month of events in one bucket. That is an unbounded array.

Writes use `$push` and `$inc` on `n`. Concurrent writes to the same bucket need care (retry or a unique bucket key).

### Questions

#### Theoretical questions

1. What is a bucket document?
2. Why does bucketing reduce the number of documents?
3. When do you close a bucket?
4. How is this different from one document per event?
5. Why must a bucket stay bounded?

#### Easy practical tasks

1. Write two bucket documents for two hours of one device. Use a small `readings` array.
2. Query buckets where `start` is in one day. Write the filter.
3. `$push` one reading and `$inc` `n`. Write the update.
4. Make a table: "One doc per event" vs "Bucket". Add three rows.

#### Medium practical tasks

1. Choose a bucket key: `deviceId` + hour. Write the unique index.
2. Compare count of documents for 10 000 events as one-per-event vs 100 events per bucket.
3. Read time series collections in the manual (high-level). Write one difference from a manual bucket.

#### Advanced practical tasks

1. Handle two writers that `$push` to the same bucket. Write a retry when the bucket is full.
2. Design a rollup: hourly buckets into a daily summary. Write the job keys.

---

## Outlier pattern

The **outlier pattern** keeps the common case in the main document. It moves a rare huge case to another collection.

Example: most books have fewer than 20 coauthors. A few books have 500 contributors. You embed authors up to 20. You set `hasExtraAuthors: true` and you store the rest in `book_authors`.

The application:

1. Reads the book.
2. If the flag is true, reads the extra collection.

Most reads do not need the extra query.

Use this pattern when:

- The distribution has a long tail
- The common case must stay small and fast
- You still must store the tail

Do not use this pattern if every document is large. Then the main model is wrong.

The bound must be real. Measure the percentiles. If the 99th percentile is still huge, embed is the wrong default.

Writes that add an author must move data when the bound is crossed.

### Questions

#### Theoretical questions

1. What is an outlier in this pattern?
2. Why does a flag exist on the main document?
3. When is the extra read acceptable?
4. Why is this pattern useless if every document is huge?
5. What must a write do when the bound is crossed?

#### Easy practical tasks

1. Write a book with 3 embedded authors and a book with a flag and extra authors.
2. Write the two-step read in numbered steps.
3. Pick a bound (for example 20). Write why you picked that number (even if it is a guess).
4. Draw the common path and the outlier path.

#### Medium practical tasks

1. Measure a list length in sample data (or invent a histogram). Mark a bound at p95.
2. Write the update that moves the 21st author to the extra collection.
3. Compare outlier vs always-reference. Write who pays the extra read.

#### Advanced practical tasks

1. Find the official outlier-pattern description. Write the example domain that they use.
2. Design outliers for a chat: last 50 messages embedded, older in `messages`. Write flags and indexes.

---

## Computed pattern

The **computed pattern** stores a value that you can calculate from other fields. You update the stored value when the source changes.

Example: store `commentCount` on a post. Increment it when you insert a comment. Do not `$count` all comments on every page view.

Example: store `lineTotal` and `orderTotal` when you change a line.

Benefits:

- Fast reads
- Simple `find`

Costs:

- Writes must keep the value correct
- A bug can make the stored value stale
- You may need a repair job

Use this pattern when reads of the computed value are frequent and the formula is stable.

Do not store a computed value that changes on every read for another reason (current time) unless you also store `computedAt` and you accept staleness.

You can recompute with an aggregation job. The job is the repair path, not the hot path.

Single-document computations can use an aggregation update pipeline. Multi-document counts need a careful increment or a transaction.

### Questions

#### Theoretical questions

1. What is a computed field in this pattern?
2. Why increment `commentCount` instead of counting on read?
3. What is the main cost of this pattern?
4. When is a nightly recompute acceptable?
5. Why is "minutes since created" a poor stored field without extra rules?

#### Easy practical tasks

1. Add `commentCount: 0` to a post. `$inc` it when you insert a comment.
2. Show a wrong count if you insert a comment and forget `$inc`.
3. List three computed fields for a shop order.
4. Make a table: "Compute on read" vs "Store computed". Add three rows.

#### Medium practical tasks

1. Write a repair aggregation that sets `commentCount` from a `comments` collection.
2. Update `orderTotal` in the same document as the lines (single document). Write the `$set` or pipeline update.
3. Compare a computed field with `$facet` on each request. Write latency in qualitative terms.

#### Advanced practical tasks

1. Design a concurrent increment of `likeCount`. Write why `$inc` is required.
2. Read the official computed pattern. Write one example that is not a count.

---

## Subset pattern

The **subset pattern** stores a **small copy** of related data for the common read. The full document lives in another collection.

Example: a product list stores `{ productId, name, thumbUrl, price }`. The `products` collection stores the long description, all images, and warehouse fields.

The list page does not load 50 KB of description per row.

This pattern is denormalization with a size goal. You copy a subset, not the whole product.

Updates to a source field that is in the subset must update the copies, or you accept staleness. Price on a list is often live enough to update in the background, or you read price live and you still subset the heavy fields.

Do not subset fields that change every second unless you have a refresh plan.

Do not duplicate the entire document. Then you have two sources of truth for everything.

### Questions

#### Theoretical questions

1. What does a subset store?
2. Why does a list page need a subset?
3. How is a subset different from a full duplicate?
4. What must you do when a subset field changes in the source?
5. When is staleness acceptable?

#### Easy practical tasks

1. Write a full product and a subset used on a card.
2. Write a `find` projection that builds a subset from the full product (when you do not store a copy).
3. List fields that must stay in the full document only.
4. Draw list read vs detail read.

#### Medium practical tasks

1. Compare stored subset vs projection from one collection. Write when a second collection is still useful (for example a different team or a different life cycle).
2. Write an update path when `name` changes: update `products` and all `cart_items` subsets.
3. Estimate payload for 50 list items with and without a 5 KB description.

#### Advanced practical tasks

1. Find the official subset pattern. Write how it relates to the extended reference.
2. Design a subset that is stale by at most 5 minutes. Write a job or a cache rule.

---

## Extended reference

An **extended reference** stores an id **and** a few fields from the referenced document.

```javascript
{
  productId: ObjectId("..."),
  name: "Nail",
  price: NumberDecimal("1.50")
}
```

The id is the live link. The extra fields are a snapshot or a cache for the common read.

This pattern sits between a bare reference and a full embed.

Use it when:

- You must join or navigate by id
- You must display a label without a lookup
- The extra fields are few and have a known freshness rule

Decide if the extra fields are a **snapshot** (order history) or a **cache** (can refresh).

Do not extend the reference with twenty fields. Then you have a subset that is too large, or a hidden embed.

Keep field names clear. `name` on an order line is the snapshot name. Do not pretend it is always the catalog name.

### Questions

#### Theoretical questions

1. What two parts does an extended reference have?
2. How is it different from a bare ObjectId?
3. How is it different from a full embed of the product?
4. Why must you decide snapshot vs cache?
5. Why is a 20-field extension a problem?

#### Easy practical tasks

1. Write an order line as a bare id, as an extended reference, and as a full embed. Three snippets.
2. Query orders by `productId`. Write the filter.
3. Display an order without `$lookup`. Write which fields you use.
4. Make a table: "Bare ref", "Extended ref", "Embed". Add three rows.

#### Medium practical tasks

1. For invoice lines, mark each extra field as snapshot. For a cart, mark each extra field as cache. Write the difference.
2. Write a refresh job that updates cached `name` on open carts when the product name changes.
3. Compare `$lookup` vs extended reference on a list of 50 orders. Write the read count.

#### Advanced practical tasks

1. Read the official extended-reference pattern. Write the example that they use.
2. Design an extended reference that includes `schemaVersion` of the product. Write when you refresh.

---

## Schema versioning field

A **schema version** field records the shape of the document. A common name is `schemaVersion` or `v`.

```javascript
{ schemaVersion: 2, name: "Ann", displayName: "Ann" }
```

The application reads the version. It runs a **reader** that understands version 1 and version 2. A **writer** writes only the current version, or it writes both fields during a transition.

This is the expand-contract idea:

1. Deploy code that reads old and new.
2. Write new fields (expand).
3. Backfill old documents.
4. Remove old-field reads (contract).

The version field makes the shape explicit. You can query `{ schemaVersion: 1 }` for a backfill.

Do not increment the version for every small add if the reader can treat a missing field as a default. Increment when the meaning of a field changes or when the shape is hard to detect.

Do not store two full incompatible meanings in one field without a version or a new name.

Validation can allow more than one version with `oneOf` (next topic).

### Questions

#### Theoretical questions

1. What does `schemaVersion` record?
2. Why must a reader understand more than one version during a change?
3. What is expand-contract?
4. When do you increment the version?
5. How does a version field help a backfill job?

#### Easy practical tasks

1. Write a v1 user `{ name }` and a v2 user `{ name, displayName, schemaVersion: 2 }`.
2. Write a reader rule: if version is missing, treat as 1.
3. Query `{ schemaVersion: { $ne: 2 } }` for a backfill list.
4. List three changes that need a new version and three that do not.

#### Medium practical tasks

1. Write a four-step expand-contract plan to rename `name` to `displayName`.
2. Add `schemaVersion` to an existing test collection. Update old documents with `$set`.
3. Sketch `oneOf` for v1 and v2 (fields only). No need to run the validator yet.

#### Advanced practical tasks

1. Read expand-contract for databases (MongoDB or general). Write a one-page plan for a field type change (string to object).
2. Design a dual-write period: writers write v2, a job converts v1. Write the stop condition for the job.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Which patterns reduce document size on the hot path (outlier, subset, bucket), and how do they differ?
2. Which patterns add duplication on purpose (computed, subset, extended reference), and what extra write do they need?
3. When do you pick the attribute pattern instead of a polymorphic set of named fields?
4. How does `schemaVersion` support a change that uses the computed or subset pattern?
5. A teammate applies all seven patterns to one collection. Which facts do you use to stop that?

#### Easy practical tasks

1. Name the pattern for: last 100 metrics in one hour document; `likeCount` on a post; product id plus name on a line.
2. Write a cheat sheet: seven pattern names, one sentence each.
3. Add `schemaVersion: 1` to one sample document for each of two collections that you use.
4. Draw a shop that uses extended reference and computed `orderTotal` only. No other patterns.

#### Medium practical tasks

1. Pick a real or invented product catalog. Apply attribute or not. Apply subset for lists. Write why.
2. Design buckets for click events and outliers for a viral post's comments. Write bounds and flags.
3. Write a change that needs `schemaVersion` 1 to 2 and a computed field repair. List the deploy order.

#### Advanced practical tasks

1. Read the official building-with-patterns series. Map each official pattern in this topic to a URL and one extra fact.
2. Review a model from topic 5. Add at most three patterns. Write which problems they solve and the write paths you added.
