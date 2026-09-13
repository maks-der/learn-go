# 4. Querying Documents

## Description

This topic shows how you shape find results. You project fields, you sort, you skip, you limit, and you use cursors. You also query embedded documents and arrays. Complete the CRUD topic first.

Use one term for each concept. **Projection** selects fields in the result. **Sort** sets the order. **Skip** and **limit** page the result. A **cursor** is the iterator that the server uses for `find`. **Dot notation** names a field inside an embedded document or an array.

Do not run unbounded `find` in an application. Always think about limit, filter, and indexes.

---

## Projection (include / exclude fields)

A projection is the second argument to `find` (or a `.project()` helper). It selects which fields the server returns.

Include mode: set fields to `1`.

```javascript
db.items.find({ qty: { $gt: 0 } }, { name: 1, qty: 1 })
```

`_id` is included unless you exclude it:

```javascript
db.items.find({}, { name: 1, _id: 0 })
```

Exclude mode: set fields to `0`.

```javascript
db.items.find({}, { notes: 0 })
```

You cannot mix include and exclude of normal fields in one projection. You can exclude `_id` while you include other fields. That is the usual exception.

Projection reduces network payload. It does not change the stored document. A smaller projection can make a query **covered** if an index holds all returned fields. That is a later performance topic.

Do not project a password hash to a public API. Projection is part of security for reads.

`findOne` accepts a projection in the same way.

### Questions

#### Theoretical questions

1. What does a projection change: the stored document or the result?
2. How do you exclude `_id` when you include `name`?
3. Can you include `name` and exclude `notes` in the same projection?
4. Why does projection help a network-heavy API?
5. What is a covered query at a high level?

#### Easy practical tasks

1. Find all items. Project only `name`. Write whether `_id` appears.
2. Project `name` and set `_id` to `0`. Confirm `_id` is absent.
3. Exclude one large field. Compare the printed size of the result.
4. Use `findOne` with a projection. Write the returned keys.

#### Medium practical tasks

1. Try a mixed include/exclude projection. Record the error.
2. Project a nested field with dot notation: `{ "addr.city": 1 }`. Write the shape of the result.
3. Compare payload size with and without a projection on a document that has a large array. Use `Object.bsonsize` on one result.

#### Advanced practical tasks

1. Read about aggregation `$project` vs find projection. Write one feature that `$project` has that find projection lacks.
2. Design an API response that must never include `hash` or `internal`. Write the projection and a test document.

---

## Sort, skip, limit

`sort` orders the cursor. You pass fields and direction: `1` is ascending, `-1` is descending.

```javascript
db.items.find().sort({ qty: -1 })
```

`limit` sets the maximum number of documents that the cursor returns.

```javascript
db.items.find().sort({ qty: -1 }).limit(10)
```

`skip` omits that many documents in the sort order, then returns the rest (or the limit).

```javascript
db.items.find().sort({ name: 1 }).skip(20).limit(10)
```

Sort uses the field values and BSON type order. Missing fields and mixed types can surprise you. Keep types stable.

`skip` on a large number is expensive. The server still walks the skipped documents. For deep pages, use a range filter on the sort key (`qty < lastQty`) instead of a large skip.

A sort can use an index. A sort that does not match an index can fail or use a lot of memory. The server has a sort memory limit. Later topics cover indexes and `allowDiskUse` for aggregation.

Always set `sort` when you use `skip` and `limit`. Without a sort, page order is not a stable contract.

### Questions

#### Theoretical questions

1. What do `1` and `-1` mean in `sort`?
2. Why is a large `skip` expensive?
3. Why must you sort when you page with `skip` and `limit`?
4. What is a range-key alternative to `skip`?
5. What happens when you sort a field that has mixed types?

#### Easy practical tasks

1. Insert five documents with different `qty`. Sort by `qty` descending. Write the order of names.
2. Apply `limit(2)` to that sort. Write the two names.
3. Apply `skip(1).limit(2)` on a sort by `name`. Write the names.
4. Sort by two fields: `qty` then `name`. Write the order.

#### Medium practical tasks

1. Page a collection of 25 documents into pages of 5 with `skip` and `limit`. Write the five commands.
2. Implement the next page with a range filter on `_id` or `qty` instead of `skip`. Write the two queries (page 1 and page 2).
3. Sort on a missing field for some documents. Write where those documents appear.

#### Advanced practical tasks

1. Read the sort-memory limit in the manual. Cause a large in-memory sort in aggregation or find (use a test collection). Record the error or the disk plan.
2. Compare `skip` paging and cursor-key paging on 10 000 documents. Time page 1 and page 200. Write the four times.

---

## Cursors

`find` returns a **cursor**. The cursor is not the full result set in client memory. The server sends **batches**. When the client needs more documents, it sends `getMore`.

In `mongosh`, a `find()` print shows the first batch. You can iterate:

```javascript
const c = db.items.find({ qty: { $gt: 0 } })
c.forEach((doc) => { print(doc.name) })
```

Cursor helpers include `hasNext`, `next`, `toArray`, `limit`, `sort`, and `project`. `toArray` loads all remaining documents into memory. Do not call `toArray` on an unbounded result.

A cursor can **time out** on the server if the client does not consume it. The default idle timeout is 10 minutes for many deployments. Tailable cursors are a special case (capped collections, change streams).

You can set a batch size. A small batch size increases round trips. A huge batch size increases memory.

Close a cursor when you stop early in a driver. `mongosh` often closes it for you at the end of a script.

`explain` on a cursor shows the plan. Use it when you study indexes.

### Questions

#### Theoretical questions

1. Why does `find` return a cursor instead of a full array?
2. What is `getMore`?
3. Why is `toArray` dangerous on a large match?
4. What happens if a client leaves a cursor idle for a long time?
5. What is a batch size?

#### Easy practical tasks

1. Assign a `find` cursor to a variable. Call `hasNext` and `next` once. Write the `_id`.
2. Use `forEach` to print one field from all matches. Keep the collection small.
3. Call `toArray` on a `limit(3)` cursor. Write the array length.
4. Run `find().explain()`. Write the `stage` name that you see.

#### Medium practical tasks

1. Set `batchSize(2)` and iterate. Count how many times you call `next` for six documents.
2. In a driver, open a cursor, read one document, then close the cursor. Write the API names.
3. Compare `find().limit(100).toArray()` with a `forEach` that counts. Write memory use in qualitative terms.

#### Advanced practical tasks

1. Read the cursor timeout and `noCursorTimeout` page. Write why `noCursorTimeout` is risky.
2. Capture `getMore` in Atlas profiler, logs, or a driver event. Write one batch size that you observed.

---

## Embedded documents and dot notation

An embedded document is a field whose value is a document.

```javascript
{
  name: "Ann",
  addr: { city: "Oslo", zip: "0001" }
}
```

**Dot notation** names a path: `"addr.city"`.

Equality on the whole embedded document requires an exact match of fields and order in some cases. Prefer dot notation for one field:

```javascript
db.people.find({ "addr.city": "Oslo" })
```

A query `{ addr: { city: "Oslo" } }` can fail to match if `addr` also has `zip`. The filter means "addr equals this exact document".

You can `$set` a nested field with the same dot path. You can project a nested field with the same path.

Arrays of embedded documents also use dot notation: `"items.sku"`. That query matches if **any** array element has `sku` with that value. The next section covers `$elemMatch` when two fields must match on the **same** element.

Do not use a dot in an application field name. The dot is the path separator.

### Questions

#### Theoretical questions

1. What does `"addr.city"` mean in a filter?
2. Why can `{ addr: { city: "Oslo" } }` miss a document that has `addr.zip`?
3. How does `"items.sku"` match an array of objects?
4. Why must application field names not contain a dot?
5. How do you set only `addr.zip` without replacing `addr.city`?

#### Easy practical tasks

1. Insert two people with different `addr.city`. Find by `"addr.city"`.
2. Project only `"addr.zip"`. Write the result shape.
3. `$set` `"addr.zip"` on one person. Confirm `city` is still there.
4. Query the whole `addr` object with an exact document. Show a miss when one extra nested field exists.

#### Medium practical tasks

1. Insert an array `items: [ { sku: "A", n: 1 }, { sku: "B", n: 2 } ]`. Query `"items.sku": "A"`. Write the full parent document that you get.
2. Query `"items.sku": "A"` and `"items.n": 2` without `$elemMatch`. Write whether the document matches. Explain.
3. Build a three-level nest (`a.b.c`). Query and project the leaf.

#### Advanced practical tasks

1. Read about `$getField` or fields with dots/dollars in modern MongoDB. Write when you would avoid such names anyway.
2. Model an address as embed vs a separate collection. Write two queries: one with dot notation, one with a follow-up find.

---

## Array queries: `$elemMatch`, `$size`, `$all`

Arrays need special operators when the match is not a single equality.

**Equality to the full array** uses `{ tags: ["a", "b"] }`. Order and exact contents matter.

**Equality to one element** uses `{ tags: "a" }`. The document matches if any element equals `"a"`.

`$all` matches if the array contains all listed values (any order):

```javascript
db.posts.find({ tags: { $all: ["db", "mongo"] } })
```

`$size` matches an exact length:

```javascript
db.posts.find({ tags: { $size: 3 } })
```

`$size` cannot take a range. For a range of lengths, store `tagCount` or use aggregation.

`$elemMatch` requires **one** array element to match all conditions:

```javascript
db.orders.find({
  lines: { $elemMatch: { sku: "nail", qty: { $gte: 10 } } }
})
```

Without `$elemMatch`, `{ "lines.sku": "nail", "lines.qty": { $gte: 10 } }` can match **two different** elements.

You can combine `$elemMatch` with projection to return only the matched element (`$` positional projection or `$elemMatch` in projection). Read the current manual for the form that you use.

### Questions

#### Theoretical questions

1. What is the difference between `{ tags: "a" }` and `{ tags: ["a"] }`?
2. What does `$all` require?
3. Why can `$size` not express "at least 3 elements"?
4. When do you need `$elemMatch` instead of two dotted predicates?
5. What does `$elemMatch` mean for "the same array element"?

#### Easy practical tasks

1. Insert posts with `tags` arrays. Find `{ tags: "mongo" }`.
2. Find tags that contain both `"db"` and `"mongo"` with `$all`.
3. Find arrays with `$size: 2`.
4. Insert one order with two lines. Write an `$elemMatch` query that matches only one intended line pair.

#### Medium practical tasks

1. Reproduce the false match: two dotted predicates without `$elemMatch`. Then fix it with `$elemMatch`.
2. Project the matching array element with the positional `$` or `$elemMatch` projection. Write the result.
3. Store `tagCount`. Query `tagCount: { $gte: 3 }`. Write why this replaces `$size` for a range.

#### Advanced practical tasks

1. Read about multikey indexes and `$elemMatch`. Write whether one compound index on two array fields is valid (the restriction).
2. Design queries for "all lines in stock" vs "at least one line in stock". Write both filters.

---

## Regular expressions (careful with indexes)

You can match strings with a regular expression.

```javascript
db.items.find({ name: /nail/i })
db.items.find({ name: { $regex: "nail", $options: "i" } })
```

`i` makes the match case-insensitive.

A regex that is not a prefix match often **cannot use an index** well. A leading `.*` forces a scan of many values. A prefix regex such as `/^nail/` can use an index on `name` more often.

Case-insensitive search without a collation or a normalized field can also skip a simple index. Atlas Search is a better tool for rich text search. That is a later topic.

Do not pass raw user text into a regex without escape of special characters. The user can write a pattern that is slow (a complexity attack) or that matches too much.

Prefer exact match or `$in` when you know the full string. Prefer a normalized lowercase field for simple case-insensitive equality.

Test regex queries with `explain`. If you see `COLLSCAN` on a large collection, change the query or add a better field.

### Questions

#### Theoretical questions

1. What does the `i` option do?
2. Why is `/.*nail.*/` often bad for indexes?
3. When can a regex use an index?
4. Why must you escape user text in a regex?
5. When is a lowercase copy of a field better than a regex?

#### Easy practical tasks

1. Find names that match `/nail/i`. Write the matches.
2. Write the same query with `$regex` and `$options`.
3. Run `explain` on the regex query. Write the stage name.
4. Replace the regex with an equality query on the same data. Run `explain` again. Write the stage name.

#### Medium practical tasks

1. Compare `/^na/` and `/na/` with `explain` on an indexed `name` field. Write which plan looks better.
2. Insert a user string that contains `.` or `*`. Show a wrong match if you do not escape. Then escape and show the fix.
3. Add a `nameLower` field. Query equality on `nameLower`. Write why this can use a normal index.

#### Advanced practical tasks

1. Read about collation and case-insensitive indexes. Write one way to get case-insensitive equality without a regex.
2. Read Atlas Search vs `$regex` (high-level). Write two features that Search has that `$regex` does not.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do projection, sort, skip, and limit work together on one cursor?
2. Why does dot notation need `$elemMatch` for some array documents but not for a single embedded document?
3. Which query features in this topic can ignore an index if you write them badly?
4. What is the difference between a cursor batch and a `limit`?
5. A teammate pages with `skip: 50000`. Which facts do you use to propose a better method?

#### Easy practical tasks

1. Find items with `qty > 0`, project `name` only, sort by `name`, limit 3. Write the commands as one chain.
2. Write a cheat sheet: projection include/exclude, sort direction, skip cost, cursor, dot path, `$elemMatch`, regex index warning.
3. Query an embedded field and an array field in the same collection. Write both filters.
4. Run `explain` on one find that you wrote in this topic. Save the output in a file.

#### Medium practical tasks

1. Build a paged API sketch: filter, sort, limit, projection. Write the `mongosh` form for page 2 with skip and the form with a range key.
2. Write three array queries (`$all`, `$size`, `$elemMatch`) on the same collection. Draw which documents match each.
3. Add an index on `name`. Compare `explain` for equality, prefix regex, and leading-wildcard regex.

#### Advanced practical tasks

1. In a driver, stream a cursor to a file without `toArray`. Write the loop and how you handle errors and close.
2. Design a product search: exact sku, prefix name, and tag all-of. Write which operator you use for each and which index you would add (names only).
