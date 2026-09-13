# 2. CRUD and Queries

## Description

CRUD means Create, Read, Update, and Delete. This topic shows how you insert documents, find documents, change fields, replace documents, and delete documents. You also shape results with projection, sort, skip, limit, and cursors. You query nested fields, arrays, and strings.

Complete the getting-started topic first. Complete this topic before you study data modeling.

Use one term for each concept. **insert** adds a new document. **find** reads documents. **update** changes fields in a document. **replace** writes a full new document body. **delete** removes a document. **upsert** updates a match or inserts a document when no match exists. **Projection** selects fields in the result. A **cursor** is the iterator that the server uses for `find`.

All examples use `mongosh` and collection helpers such as `db.items.insertOne`.

---

## `insertOne` / `insertMany`

`insertOne` writes one document. The server returns an acknowledgment and the `_id` of the new document.

```javascript
db.items.insertOne({ name: "nail", qty: 100 })
```

If you omit `_id`, MongoDB adds an ObjectId. If you set `_id`, that value is the key.

`insertMany` writes many documents in one call. You pass an array.

```javascript
db.items.insertMany([
  { name: "pin", qty: 50 },
  { name: "bolt", qty: 20 }
])
```

By default, `insertMany` stops at the first error (`ordered: true`). Set `ordered: false` to continue after an error. Duplicate `_id` is a common error.

Insert is a write. A failed insert does not change the collection for that document. Do not retry an insert with the same `_id` unless you intend to fail on a duplicate.

Do not call `insertOne` in a loop when you have a batch. Use `insertMany` or a bulk write. A loop of single inserts is slower.

`insertOne` and `insertMany` do not replace an existing document. Use update, replace, or upsert for that case.

### Questions

#### Theoretical questions

1. What does `insertOne` return to the client?
2. What happens when you omit `_id` on insert?
3. What is the default `ordered` behavior of `insertMany`?
4. Why is a loop of `insertOne` a poor batch method?
5. Does `insertOne` overwrite a document that already has the same `_id`?

#### Easy practical tasks

1. Insert one document with `insertOne`. Write the returned `_id`.
2. Insert three documents with `insertMany`. Run `find`. Count the documents.
3. Insert a document with a chosen `_id`. Find it by that `_id`.
4. Run `insertMany` with two documents that share `_id`. Record the error.

#### Medium practical tasks

1. Run `insertMany` with `ordered: false` and one duplicate `_id` in the middle. Write how many documents exist after the call.
2. Time 100 `insertOne` calls versus one `insertMany` of 100 documents. Write the two times.
3. Use `db.items.bulkWrite` with two insert operations. Confirm both documents exist.

#### Advanced practical tasks

1. Read the write-concern page. Insert with `w: "majority"` if you have a replica set. Write what extra wait that option adds.
2. Write a script that inserts 1000 documents in batches of 100. Record total time and any duplicate-key errors.

---

## `find` / `findOne` and query operators

`find` returns a cursor. A cursor is an iterator over matching documents. `mongosh` prints the first batch when you type `find()`.

```javascript
db.items.find()
db.items.find({ qty: 100 })
```

An empty filter `{}` matches all documents. A filter `{ qty: 100 }` matches documents where `qty` equals 100. Equality is the default. You do not need `$eq` for a simple equal match.

`findOne` returns one document or `null`. It does not return a cursor.

```javascript
db.items.findOne({ name: "nail" })
```

`findOne` is not a special server command with a different match rule. It is a find that limits the result to one document.

Do not use `find()` on a large collection in a program without a filter, a limit, or both. APIs must not return an unbounded result.

`countDocuments` counts matches. `estimatedDocumentCount` uses metadata. Use `countDocuments` when you need a correct filter count.

A filter can use operators. An operator is a key that starts with `$`.

- `$eq` — equal. `{ qty: { $eq: 100 } }` is the same as `{ qty: 100 }`.
- `$gt` — greater than. `$gte`, `$lt`, and `$lte` are the related operators.
- `$in` — the field value is in a list. `{ name: { $in: ["nail", "pin"] } }`.
- `$and` — all conditions must match. A filter with several fields is an implicit AND.
- `$or` — at least one condition must match. You pass an array of filters.
- `$exists` — the field is present (`true`) or absent (`false`). A field with value `null` still exists.

Example:

```javascript
db.items.find({ qty: { $gt: 10 } })
db.items.find({ $or: [ { name: "nail" }, { qty: { $lt: 5 } } ] })
db.items.find({ color: { $exists: false } })
```

`$in` is not `$or` of full documents. `$in` tests one field against a list of values.

`$and` is useful when you need two operators on the same field. Prefer a simple implicit AND when it is clear.

Do not build a filter from a string that includes user text. Use a structured object.

Comparison uses BSON types. A string does not compare as a number.

### Questions

#### Theoretical questions

1. What does `find` return?
2. What does `findOne` return when no document matches?
3. What does `$in` test?
4. Does `{ color: { $exists: true } }` match `{ color: null }`?
5. Why can `{ qty: { $gt: 10 } }` miss a document where `qty` is the string `"100"`?

#### Easy practical tasks

1. Insert three items. Run `find()`. Write how many documents the shell prints.
2. Run `findOne` with a filter that matches. Write the `name` field.
3. Find documents with `qty` greater than 10, and documents where `name` is in a list of two values.
4. Find documents that miss the field `color`. Then find documents that match `$or` of two equality filters.

#### Medium practical tasks

1. Compare `find({}).limit(1)` and `findOne({})` in `mongosh`. Write how the printed result differs.
2. Write one query that needs `$and` because the same field has two operators. Run it.
3. Query `$gt` on a field that has both numbers and strings. Write which documents match.

#### Advanced practical tasks

1. Read the BSON comparison order page. Write where null, numbers, and strings sit in that order.
2. Open a cursor in a driver or with `forEach` on a large find. Write one observation about batches or `getMore`.

---

## `updateOne` / `updateMany`, `replaceOne`, `deleteOne` / `deleteMany`, upsert

`updateOne` changes the first matching document. `updateMany` changes all matching documents.

An update uses an **update document** with operators. You do not pass a full replacement document to these methods.

- `$set` — set a field to a value. Creates the field if it is missing.
- `$unset` — remove a field.
- `$inc` — add a number to a numeric field. The field must be a number or missing (then MongoDB treats it as 0).

```javascript
db.items.updateOne({ name: "nail" }, { $set: { qty: 90 } })
db.items.updateMany({ qty: { $gt: 0 } }, { $inc: { qty: -1 } })
db.items.updateOne({ name: "nail" }, { $unset: { obsolete: "" } })
```

The filter is the first argument. The update is the second argument. A missing filter matches the first document in `updateOne`. That is dangerous. Always pass a filter that you intend.

`$set` can set nested fields with dot notation: `{ $set: { "addr.city": "Oslo" } }`.

Updates are atomic on one document. Two `$inc` operations on the same document do not lose a count if they use `$inc`, not read-then-write in the client.

`replaceOne` writes a new document body over a matching document. The second argument is a **replacement document**, not an update operator document.

```javascript
db.items.replaceOne(
  { name: "nail" },
  { name: "nail", qty: 0, unit: "box" }
)
```

The replacement document must not contain update operators such as `$set`. The replacement becomes the full content, except `_id`. You cannot change `_id` in a replace. A replace can drop fields that you omit. That is the point. It is also the risk.

`replaceOne` updates at most one document. There is no `replaceMany`.

`deleteOne` removes the first matching document. `deleteMany` removes all matching documents.

```javascript
db.items.deleteOne({ name: "pin" })
db.items.deleteMany({ qty: 0 })
```

A missing or empty filter on `deleteMany` deletes all documents in the collection. Always pass a filter that you intend. Prefer `drop()` when you intend to remove the whole collection.

`deleteOne` with a filter that matches many documents still deletes only one. Prefer `_id` when you must delete one specific document.

An **upsert** is an update or a replace that inserts a document when no document matches the filter.

```javascript
db.items.updateOne(
  { name: "wire" },
  { $set: { qty: 10 } },
  { upsert: true }
)
```

If a document with `name: "wire"` exists, MongoDB applies `$set`. If it does not exist, MongoDB creates a document. `replaceOne` also accepts `upsert: true`.

`$inc` with upsert creates the document and then increments.

Two upserts with the same filter can cause a duplicate-key error if they both insert. Handle that error and retry, or use a unique index on the filter fields.

Do not upsert with a loose filter. You can insert documents that you did not intend.

### Questions

#### Theoretical questions

1. What is the difference between `updateOne` and `replaceOne`?
2. Why is `$inc` safer than read, add, and `$set` for a counter?
3. What happens if you run `deleteMany({})`?
4. What does `upsert: true` do when the filter matches nothing?
5. Why can two concurrent upserts fail with a duplicate key?

#### Easy practical tasks

1. `$set` the `qty` of one item. Then `$inc` `qty` by 5. Find the document. Write the new value.
2. Replace one document with a new body that has one extra field. Confirm old extra fields are gone if you omitted them.
3. Insert two documents. `deleteOne` by `_id`. Confirm one remains.
4. Upsert a document that does not exist. Find it. Write the fields that MongoDB stored.

#### Medium practical tasks

1. Use `$set` with dot notation to add `addr.city` without replacing the rest of `addr`.
2. Compare `replaceOne` and `updateOne` with `$set` on the same starting document. Show which fields survive.
3. Compare `deleteMany({})` and `db.items.drop()` on a test collection. Write what each returns and what `show collections` shows.

#### Advanced practical tasks

1. Compare `updateOne` with a pipeline update (`[ { $set: ... } ]`) in the manual. Write one thing a pipeline update can do that `$set` alone cannot.
2. Design an idempotent "put item by sku" API with upsert and a unique index. Write the filter, the update, and the error retry.

---

## Projection, sort, skip, limit, and cursors

A projection is the second argument to `find` (or a `.project()` helper). It selects which fields the server returns.

Include mode: set fields to `1`.

```javascript
db.items.find({ qty: { $gt: 0 } }, { name: 1, qty: 1 })
```

`_id` is included unless you exclude it:

```javascript
db.items.find({}, { name: 1, _id: 0 })
```

Exclude mode: set fields to `0`. You cannot mix include and exclude of normal fields in one projection. You can exclude `_id` while you include other fields.

Projection reduces network payload. It does not change the stored document. Do not project a password hash to a public API.

`sort` orders the cursor. You pass fields and direction: `1` is ascending, `-1` is descending.

```javascript
db.items.find().sort({ qty: -1 }).limit(10)
```

`limit` sets the maximum number of documents that the cursor returns. `skip` omits that many documents in the sort order, then returns the rest (or the limit).

```javascript
db.items.find().sort({ name: 1 }).skip(20).limit(10)
```

`skip` on a large number is expensive. The server still walks the skipped documents. For deep pages, use a range filter on the sort key (`qty < lastQty`) instead of a large skip.

Always set `sort` when you use `skip` and `limit`. Without a sort, page order is not a stable contract.

`find` returns a **cursor**. The cursor is not the full result set in client memory. The server sends **batches**. When the client needs more documents, it sends `getMore`.

```javascript
const c = db.items.find({ qty: { $gt: 0 } })
c.forEach((doc) => { print(doc.name) })
```

`toArray` loads all remaining documents into memory. Do not call `toArray` on an unbounded result.

A cursor can time out on the server if the client does not consume it. The default idle timeout is 10 minutes for many deployments.

You can set a batch size. A small batch size increases round trips. A huge batch size increases memory.

Close a cursor when you stop early in a driver. `explain` on a cursor shows the plan.

### Questions

#### Theoretical questions

1. What does a projection change: the stored document or the result?
2. Can you include `name` and exclude `notes` in the same projection?
3. Why is a large `skip` expensive?
4. Why does `find` return a cursor instead of a full array?
5. Why is `toArray` dangerous on a large match?

#### Easy practical tasks

1. Find all items. Project only `name`. Write whether `_id` appears. Then project `name` and set `_id` to `0`.
2. Insert five documents with different `qty`. Sort by `qty` descending. Apply `limit(2)`. Write the two names.
3. Apply `skip(1).limit(2)` on a sort by `name`. Write the names.
4. Assign a `find` cursor to a variable. Call `hasNext` and `next` once. Write the `_id`.

#### Medium practical tasks

1. Try a mixed include/exclude projection. Record the error.
2. Page a collection of 25 documents into pages of 5 with `skip` and `limit`. Then implement the next page with a range filter instead of `skip`.
3. Set `batchSize(2)` and iterate. Count how many times you call `next` for six documents.

#### Advanced practical tasks

1. Compare `skip` paging and cursor-key paging on 10 000 documents. Time page 1 and page 200. Write the four times.
2. Read the cursor timeout and `noCursorTimeout` page. Write why `noCursorTimeout` is risky.

---

## Dot notation, array queries, and regular expressions

An embedded document is a field whose value is a document.

```javascript
{
  name: "Ann",
  addr: { city: "Oslo", zip: "0001" }
}
```

**Dot notation** names a path: `"addr.city"`.

```javascript
db.people.find({ "addr.city": "Oslo" })
```

A query `{ addr: { city: "Oslo" } }` can fail to match if `addr` also has `zip`. The filter means "addr equals this exact document". Prefer dot notation for one field.

You can `$set` a nested field with the same dot path. Arrays of embedded documents also use dot notation: `"items.sku"`. That query matches if **any** array element has `sku` with that value.

Do not use a dot in an application field name. The dot is the path separator.

Arrays need special operators when the match is not a single equality.

**Equality to the full array** uses `{ tags: ["a", "b"] }`. Order and exact contents matter.

**Equality to one element** uses `{ tags: "a" }`. The document matches if any element equals `"a"`.

`$all` matches if the array contains all listed values (any order):

```javascript
db.posts.find({ tags: { $all: ["db", "mongo"] } })
```

`$size` matches an exact length. `$size` cannot take a range. For a range of lengths, store `tagCount` or use aggregation.

`$elemMatch` requires **one** array element to match all conditions:

```javascript
db.orders.find({
  lines: { $elemMatch: { sku: "nail", qty: { $gte: 10 } } }
})
```

Without `$elemMatch`, `{ "lines.sku": "nail", "lines.qty": { $gte: 10 } }` can match **two different** elements.

You can match strings with a regular expression.

```javascript
db.items.find({ name: /nail/i })
db.items.find({ name: { $regex: "nail", $options: "i" } })
```

`i` makes the match case-insensitive.

A regex that is not a prefix match often cannot use an index well. A leading `.*` forces a scan of many values. A prefix regex such as `/^nail/` can use an index on `name` more often.

Do not pass raw user text into a regex without escape of special characters. The user can write a pattern that is slow or that matches too much.

Prefer exact match or `$in` when you know the full string. Prefer a normalized lowercase field for simple case-insensitive equality. Test regex queries with `explain`.

### Questions

#### Theoretical questions

1. What does `"addr.city"` mean in a filter?
2. Why can `{ addr: { city: "Oslo" } }` miss a document that has `addr.zip`?
3. When do you need `$elemMatch` instead of two dotted predicates?
4. Why is `/.*nail.*/` often bad for indexes?
5. Why must you escape user text in a regex?

#### Easy practical tasks

1. Insert two people with different `addr.city`. Find by `"addr.city"`. Project only `"addr.zip"`.
2. Insert posts with `tags` arrays. Find `{ tags: "mongo" }`. Then find tags that contain both `"db"` and `"mongo"` with `$all`.
3. Find arrays with `$size: 2`. Write an `$elemMatch` query on an order with two lines.
4. Find names that match `/nail/i`. Write the same query with `$regex` and `$options`.

#### Medium practical tasks

1. Reproduce the false match: two dotted predicates without `$elemMatch`. Then fix it with `$elemMatch`.
2. Store `tagCount`. Query `tagCount: { $gte: 3 }`. Write why this replaces `$size` for a range.
3. Compare `/^na/` and `/na/` with `explain` on an indexed `name` field. Write which plan looks better.

#### Advanced practical tasks

1. Read about multikey indexes and `$elemMatch`. Write whether one compound index on two array fields is valid.
2. Design a product search: exact sku, prefix name, and tag all-of. Write which operator you use for each and which index you would add (names only).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Map Create, Read, Update, Delete to MongoDB method names. Include replace and upsert.
2. How do projection, sort, skip, and limit work together on one cursor?
3. Why does MongoDB separate `updateOne` (operators) from `replaceOne` (full body)?
4. Why does dot notation need `$elemMatch` for some array documents but not for a single embedded document?
5. A teammate pages with `skip: 50000` and calls `find().toArray()` in an API. Which facts do you use in the review?

#### Easy practical tasks

1. In one collection, insert two documents, find one, update one field, replace the other, delete one. Write the final `find` output.
2. Write a cheat sheet: insert, find, update, replace, delete, upsert, projection, sort, skip cost, `$elemMatch`, regex index warning.
3. Find items with `qty > 0`, project `name` only, sort by `name`, limit 3. Write the commands as one chain.
4. Query an embedded field and an array field in the same collection. Write both filters.

#### Medium practical tasks

1. Write a small script that upserts a counter document 20 times with `$inc`. Write the final count.
2. Use `bulkWrite` with one insert, one update, and one delete. Confirm the collection state.
3. Add an index on `name`. Compare `explain` for equality, prefix regex, and leading-wildcard regex.

#### Advanced practical tasks

1. Read `bulkWrite` ordered vs unordered. Cause one duplicate-key error in a batch. Write how many operations applied.
2. In a driver, stream a cursor to a file without `toArray`. Write the loop and how you handle errors and close.
