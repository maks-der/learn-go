# 3. CRUD Basics

## Description

CRUD means Create, Read, Update, and Delete. This topic shows the basic write and read methods in MongoDB. You insert documents, you find documents, you update fields, you replace documents, and you delete documents. Complete this topic before you study projections, sort, and array queries.

Use one term for each concept. **insert** adds a new document. **find** reads documents. **update** changes fields in a document. **replace** writes a full new document body. **delete** removes a document. **upsert** updates a match or inserts a document when no match exists.

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

## `find` / `findOne`

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

Do not use `find()` on a large collection in a program without a filter, a limit, or both. The next topic covers `limit`. APIs must not return an unbounded result.

`countDocuments` counts matches. `estimatedDocumentCount` uses metadata. Use `countDocuments` when you need a correct filter count.

### Questions

#### Theoretical questions

1. What does `find` return?
2. What does `findOne` return when no document matches?
3. What does an empty filter match?
4. Why is `find()` without a limit risky in an application?
5. When do you use `countDocuments` instead of `estimatedDocumentCount`?

#### Easy practical tasks

1. Insert three items. Run `find()`. Write how many documents the shell prints.
2. Run `findOne` with a filter that matches. Write the `name` field.
3. Run `findOne` with a filter that does not match. Write the result.
4. Run `db.items.countDocuments({})`. Write the number.

#### Medium practical tasks

1. Compare `find({}).limit(1)` and `findOne({})` in `mongosh`. Write how the printed result differs.
2. Use `find` with two equality fields (`name` and `qty`). Confirm only the intended document matches.
3. Run `estimatedDocumentCount` and `countDocuments`. Write when the two numbers can differ (read the manual).

#### Advanced practical tasks

1. Open a cursor in a driver or with `forEach` on a large find. Print batch behavior from the logs or from `explain`. Write one observation.
2. Read about `find` vs `getMore` in the wire protocol (high-level). Write how the client gets the next batch.

---

## Query operators: `$eq`, `$gt`, `$in`, `$and`, `$or`, `$exists`

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

`$and` is useful when you need two operators on the same field, or when you want an explicit list. Prefer a simple implicit AND when it is clear.

Do not build a filter from a string that includes user text. Use a structured object. Later topics cover injection.

Comparison uses BSON types. Types have a comparison order. A string does not compare as a number.

### Questions

#### Theoretical questions

1. When do you need `$eq` instead of a plain equality field?
2. What does `$in` test?
3. When is `$and` required instead of two fields in one object?
4. Does `{ color: { $exists: true } }` match `{ color: null }`?
5. Why can `{ qty: { $gt: 10 } }` miss a document where `qty` is the string `"100"`?

#### Easy practical tasks

1. Find documents with `qty` greater than 10.
2. Find documents where `name` is in a list of two values.
3. Find documents that miss the field `color`.
4. Find documents that match `$or` of two equality filters.

#### Medium practical tasks

1. Write one query that needs `$and` because the same field has two operators. Run it.
2. Compare `{ a: 1, b: 2 }` with `{ $and: [ { a: 1 }, { b: 2 } ] }`. Show that they match the same documents on your data.
3. Query `$gt` on a field that has both numbers and strings. Write which documents match.

#### Advanced practical tasks

1. Read the BSON comparison order page. Write where null, numbers, and strings sit in that order.
2. Build a filter from a list of user-selected names with `$in`. Write why this is safer than string concatenation.

---

## `updateOne` / `updateMany` with `$set`, `$unset`, `$inc`

`updateOne` changes the first matching document. `updateMany` changes all matching documents.

An update uses an **update document** with operators. You do not pass a full replacement document to these methods (that is `replaceOne`).

- `$set` — set a field to a value. Creates the field if it is missing.
- `$unset` — remove a field.
- `$inc` — add a number to a numeric field. The field must be a number or missing (then MongoDB treats it as 0).

```javascript
db.items.updateOne({ name: "nail" }, { $set: { qty: 90 } })
db.items.updateMany({ qty: { $gt: 0 } }, { $inc: { qty: -1 } })
db.items.updateOne({ name: "nail" }, { $unset: { obsolete: "" } })
```

The filter is the first argument. The update is the second argument. A missing filter matches the first document in `updateOne`. That is dangerous. Always pass a filter that you intend.

`updateMany` with a filter that matches too much is also dangerous. The same rule as SQL `UPDATE` without a good `WHERE` applies.

`$set` can set nested fields with dot notation: `{ $set: { "addr.city": "Oslo" } }`.

Updates are atomic on one document. Two `$inc` operations on the same document do not lose a count if they use `$inc`, not read-then-write in the client.

### Questions

#### Theoretical questions

1. What is the difference between `updateOne` and `updateMany`?
2. What does `$set` do when the field does not exist?
3. Why is `$inc` safer than read, add, and `$set` for a counter?
4. What happens if you call `updateOne` with an empty filter?
5. How do you remove a field in an update?

#### Easy practical tasks

1. `$set` the `qty` of one item. Find the document. Confirm the new value.
2. `$inc` `qty` by 5. Find the document. Write the new value.
3. `$unset` a field that you added. Confirm the field is gone.
4. Run `updateMany` with a filter that matches two documents. Confirm both changed.

#### Medium practical tasks

1. Run two `$inc` updates from two `mongosh` sessions on the same counter document. Write the final value and why it is correct.
2. Use `$set` with dot notation to add `addr.city` without replacing the rest of `addr`.
3. Run `updateOne` with a filter that matches zero documents. Write `matchedCount` and `modifiedCount`.

#### Advanced practical tasks

1. Read about `arrayFilters` or `$min` / `$max` (pick one). Run one example. Write when you would use it.
2. Compare `updateOne` with a pipeline update (`[ { $set: ... } ]`) in the manual. Write one thing a pipeline update can do that `$set` alone cannot.

---

## `replaceOne`

`replaceOne` writes a new document body over a matching document. The filter selects the document. The second argument is a **replacement document**, not an update operator document.

```javascript
db.items.replaceOne(
  { name: "nail" },
  { name: "nail", qty: 0, unit: "box" }
)
```

The replacement document must not contain update operators such as `$set`. The replacement becomes the full content, except `_id`. `_id` stays the same unless you include the same `_id`. You cannot change `_id` in a replace.

Use `replaceOne` when the client has a full new document. Use `updateOne` when you change a few fields.

A replace can drop fields that you omit. That is the point. It is also the risk. If the client sends a partial object by mistake, you lose fields.

`replaceOne` updates at most one document. There is no `replaceMany`.

### Questions

#### Theoretical questions

1. How is `replaceOne` different from `updateOne` with `$set`?
2. What happens to fields that you omit in the replacement document?
3. Can `replaceOne` change `_id`?
4. Why is there no `replaceMany` in the usual API?
5. When do you choose replace instead of update?

#### Easy practical tasks

1. Replace one document with a new body that has one extra field. Find it. Confirm old extra fields are gone if you omitted them.
2. Try `replaceOne` with `{ $set: { qty: 1 } }` as the replacement. Record the error.
3. Replace a document and include the same `_id`. Confirm `_id` is unchanged.
4. Run `replaceOne` with a filter that matches nothing. Write `matchedCount`.

#### Medium practical tasks

1. Load a document with `findOne`, change two fields in your editor, `replaceOne` the full body. Write the risk of this pattern.
2. Compare `replaceOne` and `updateOne` with `$set` on the same starting document. Show which fields survive.
3. Use `replaceOne` with `upsert: true` (preview). Write whether a new document is inserted when no match exists.

#### Advanced practical tasks

1. Read replace and update in the driver for your language. Write the method names and which one takes operators.
2. Design a "save full aggregate" API. Write when the API must use replace and when it must use `$set` to avoid lost fields.

---

## `deleteOne` / `deleteMany`

`deleteOne` removes the first matching document. `deleteMany` removes all matching documents.

```javascript
db.items.deleteOne({ name: "pin" })
db.items.deleteMany({ qty: 0 })
```

A missing or empty filter on `deleteMany` deletes all documents in the collection. That is a common accident. Always pass a filter that you intend. Prefer `drop()` when you intend to remove the whole collection.

Deletes do not run a SQL-style `TRUNCATE` log in the same way. They remove documents. Indexes update. The `_id` values are gone. A later insert can reuse a custom `_id` if you set it.

`deleteOne` with a filter that matches many documents still deletes only one. Do not assume which one unless the filter is unique (for example `_id`).

Deleted data is not in the collection. Recovery needs a backup or a delay member. Do not use delete to test recovery on production data.

### Questions

#### Theoretical questions

1. What is the difference between `deleteOne` and `deleteMany`?
2. What happens if you run `deleteMany({})`?
3. Why is `_id` the safest filter for `deleteOne`?
4. When do you use `drop()` instead of `deleteMany`?
5. Does a delete free the `_id` for a later insert of the same custom `_id`?

#### Easy practical tasks

1. Insert two documents. `deleteOne` by `_id`. Confirm one remains.
2. `deleteMany` with a filter on `qty`. Confirm the remaining count.
3. Run `deleteOne` with a filter that matches nothing. Write `deletedCount`.
4. Insert a document, delete it, insert the same `_id` again. Confirm the insert works.

#### Medium practical tasks

1. Compare `deleteMany({})` and `db.items.drop()` on a test collection. Write what each returns and what `show collections` shows.
2. Delete in a transaction on a replica set (preview). Abort the transaction. Confirm the document is still there.
3. Write a checklist of three checks before you run `deleteMany` on a shared database.

#### Advanced practical tasks

1. Read about change streams or the oplog at a high level. Write how a delete appears as an event.
2. Time `deleteMany` of 10 000 documents versus `drop()` of a test collection. Write the two times and when `drop` is the right tool.

---

## `upsert`

An **upsert** is an update or a replace that inserts a document when no document matches the filter.

```javascript
db.items.updateOne(
  { name: "wire" },
  { $set: { qty: 10 } },
  { upsert: true }
)
```

If a document with `name: "wire"` exists, MongoDB applies `$set`. If it does not exist, MongoDB creates a document. The new document includes the filter fields and the update fields, with rules that the manual describes.

`replaceOne` also accepts `upsert: true`. The replacement document becomes the new document.

Upsert is useful for "set this key to this value" and for counters that must start at zero. `$inc` with upsert creates the document and then increments.

Race: two upserts with the same filter can cause a duplicate-key error if they both insert. Handle that error and retry, or use a unique index on the filter fields.

Do not upsert with a loose filter. You can insert documents that you did not intend.

The insert path of an upsert still requires a unique `_id`. MongoDB generates `_id` if you do not set it.

### Questions

#### Theoretical questions

1. What does `upsert: true` do when the filter matches?
2. What does `upsert: true` do when the filter matches nothing?
3. Why can two concurrent upserts fail with a duplicate key?
4. Why must the filter for an upsert be specific?
5. How does `$inc` behave with an upsert on a new document?

#### Easy practical tasks

1. Upsert a document that does not exist. Find it. Write the fields that MongoDB stored.
2. Upsert the same filter again with a different `$set`. Confirm there is still one document.
3. Run `updateOne` without upsert on a missing filter. Confirm no insert.
4. Use `replaceOne` with upsert and a full body. Find the new document.

#### Medium practical tasks

1. Upsert with `$setOnInsert` and `$set`. Write which fields change on insert only. Use the manual.
2. Create a unique index on `name`. Run two upserts in two shells at the same time. Record whether one fails.
3. Compare upsert with `findOne` plus `insertOne` in the client. Write one race that upsert avoids or still has.

#### Advanced practical tasks

1. Read the upsert field-merge rules when the filter has operators. Write one surprise (for example `$or` in the filter).
2. Design an idempotent "put item by sku" API with upsert and a unique index. Write the filter, the update, and the error retry.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Map Create, Read, Update, Delete to MongoDB method names. Include replace and upsert.
2. Why does MongoDB separate `updateOne` (operators) from `replaceOne` (full body)?
3. Which write methods can create a document, and which methods only change or remove existing documents?
4. How do filters work the same way for find, update, and delete?
5. A teammate wants to "edit a document" in the client and send it back. Which method do they need, and what can they lose?

#### Easy practical tasks

1. In one collection, insert two documents, find one, update one field, replace the other, delete one. Write the final `find` output.
2. Write a cheat sheet with method names and one line each: insert, find, update, replace, delete, upsert.
3. Run `insertOne`, then `updateOne` with `$inc`, then `findOne`. Write the three commands and the final `qty`.
4. Create `crud-notes.md` with your six method names and one warning for each (empty filter, ordered insert, and similar).

#### Medium practical tasks

1. Write a small script that upserts a counter document 20 times with `$inc`. Write the final count.
2. Copy one document, change it, and write both an `$set` update and a `replaceOne` version. Run both on two clones. Compare remaining fields.
3. Use `bulkWrite` with one insert, one update, and one delete. Confirm the collection state.

#### Advanced practical tasks

1. Read `bulkWrite` ordered vs unordered. Cause one duplicate-key error in a batch. Write how many operations applied.
2. On a replica set, run a multi-document transaction that inserts and then deletes. Commit. Then run one that aborts. Show the collection after each.
