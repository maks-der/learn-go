# 7. Aggregation Framework

## Description

Aggregation processes documents in a **pipeline**. Each **stage** changes the stream. This topic shows `$match`, `$project`, `$group`, `$sort`, `$limit`, `$lookup`, `$unwind`, `$addFields`, `$set`, and `$facet`. You also compare aggregation with many small finds, and you learn memory limits. Complete querying and modeling first.

Use one term for each concept. A **pipeline** is an array of stages. A **stage** is an object with one operator. `$match` filters. `$group` combines documents. `$lookup` reads another collection. Do not use aggregation when a single `find` is enough.

---

## Pipeline idea: `$match`, `$project`, `$group`, `$sort`, `$limit`

You run a pipeline with `aggregate`:

```javascript
db.orders.aggregate([
  { $match: { status: "paid" } },
  { $project: { userId: 1, total: 1, _id: 0 } },
  { $group: { _id: "$userId", spent: { $sum: "$total" } } },
  { $sort: { spent: -1 } },
  { $limit: 10 }
])
```

Order matters. `$match` first reduces the number of documents. Put `$match` as early as you can. `$match` can use indexes when it is at the start.

`$project` shapes fields. Use `1` to keep a field. Use an expression to compute a field. Later stages see only the projected fields.

`$group` requires `_id`. `_id` is the group key. Use `null` to group the whole stream into one document. Accumulators include `$sum`, `$avg`, `$min`, `$max`, `$push`, and `$first`.

`$sort` and `$limit` work as in `find`. A `$sort` plus `$limit` can be optimized to a top-k sort.

The stream is documents. After `$group`, each document is a group result. Field names change. Use the new names in later stages.

Write pipelines as a list. Name each stage in a comment in your notes. Test one stage at a time.

### Questions

#### Theoretical questions

1. What is a pipeline?
2. Why must `$match` sit early when you can do that?
3. What is `_id` in `$group`?
4. What happens to field names after `$project`?
5. When is `find` enough and aggregation not needed?

#### Easy practical tasks

1. Run a pipeline with only `$match`. Compare the count with `countDocuments`.
2. Add `$project` to keep two fields. Write the output keys.
3. `$group` by `status` and `$sum` 1 as `n`. Write the counts.
4. `$sort` and `$limit` the groups. Write the top line.

#### Medium practical tasks

1. Compute average `total` per `userId`. Sort by average descending. Limit 5.
2. Move `$match` after `$group` vs before. Write why the results or cost differ.
3. Use `$group` with `_id: null` to sum all `total`. Write the one number.

#### Advanced practical tasks

1. Read about `$match` index use in aggregation. Run `explain` on a pipeline that starts with `$match`. Write the plan.
2. Build a pipeline of five stages on sample data. Remove one stage at a time to debug. Write the method.

---

## `$lookup` (join-like)

`$lookup` reads documents from another collection. It is the usual join-like stage.

Simple form (equality):

```javascript
{
  $lookup: {
    from: "users",
    localField: "userId",
    foreignField: "_id",
    as: "user"
  }
}
```

The result field `user` is an **array**. A one-to-one match is still an array of one document. An empty array means no match.

Unwind or `$arrayElemAt` if you need one object.

A more general form uses a `pipeline` inside `$lookup` (correlated subquery). You can `$match` and `$project` the foreign documents.

`$lookup` can be expensive. It is not a reason to ignore embed. If every read needs a lookup, embed or snapshot may be better.

Index `foreignField` on the other collection. Without an index, each lookup can scan.

`$lookup` does not replace a relational join engine for ad-hoc analytics. Use it for known access paths.

Collections must be in the same database for the simple form (read the current manual for any exception).

### Questions

#### Theoretical questions

1. What does `$lookup` add to each document?
2. Why is the `as` field an array?
3. Why must `foreignField` have an index?
4. When is embed better than `$lookup`?
5. What is the pipeline form of `$lookup` for?

#### Easy practical tasks

1. Create `users` and `orders` with a matching id. `$lookup` users onto orders. Write the `as` array length.
2. `$project` after lookup to keep `user.name`. Write the shape.
3. Break `_id` types (string vs ObjectId) so lookup fails. Write the empty array. Fix the types.
4. Open the `$lookup` page. Write the URL.

#### Medium practical tasks

1. Use the pipeline form to look up only `{ active: true }` users and project two fields.
2. Compare two finds in the client vs one `$lookup`. Write round trips and when the client finds are clearer.
3. `explain` a `$lookup`. Write whether the foreign collection uses an index.

#### Advanced practical tasks

1. Read about `$lookup` and sharded collections. Write one restriction for your version.
2. Design an API that must not `$lookup` on every request. Write the embed or cache alternative.

---

## `$unwind`

`$unwind` deconstructs an array. Each element becomes a separate document.

```javascript
{ $unwind: "$lines" }
```

A document `{ _id: 1, lines: [ "a", "b" ] }` becomes two documents. Each has `lines: "a"` or `lines: "b"`.

Options:

- `preserveNullAndEmptyArrays: true` — keep documents that have a missing or empty array
- `includeArrayIndex` — store the index in a field

`$unwind` multiplies the number of documents. An array of 1000 elements makes 1000 documents from one parent. Put `$match` before `$unwind` when you can. Unwind only the arrays that you need.

A common pattern: `$unwind` then `$group` to aggregate array elements, or `$unwind` then `$lookup` on each line.

Do not `$unwind` a huge unbounded array. Change the model.

After `$unwind`, field `lines` is no longer an array (unless the element is an array). Later stages must use the new shape.

### Questions

#### Theoretical questions

1. What does `$unwind` do to one document with an array of three elements?
2. What does `preserveNullAndEmptyArrays` do?
3. Why can `$unwind` increase memory use?
4. Why put `$match` before `$unwind`?
5. What is a typical `$unwind` then `$group` use?

#### Easy practical tasks

1. Unwind `tags` on two posts. Count the output documents.
2. Unwind with `preserveNullAndEmptyArrays: true` on a document with no tags. Write if it appears.
3. Add `includeArrayIndex: "i"`. Write the index values.
4. Draw before and after unwind for one sample document.

#### Medium practical tasks

1. Unwind `lines`, then `$group` by `lines.sku` and `$sum` `lines.qty`.
2. Compare `$unwind` of a 2-element array vs a 200-element array. Write output counts.
3. Rebuild an array with `$group` and `$push` after unwind. Write whether order is guaranteed (read the manual).

#### Advanced practical tasks

1. Read about `$unwind` and memory. Write one way to avoid unwind (for example `$reduce` or a different model).
2. Pipeline: match paid orders, unwind lines, lookup products, group by category. Write the stage list.

---

## `$addFields`, `$set`

`$addFields` adds or replaces fields. The other fields stay.

```javascript
{ $addFields: { totalWithTax: { $multiply: [ "$total", 1.25 ] } } }
```

`$set` in aggregation is an alias of `$addFields`. It is not the same object as the update operator `$set` in `updateOne`, but the idea is similar: set fields.

Use `$addFields` when you need a computed field and you still need the old fields. Use `$project` when you want a strict output shape and you drop fields.

You can overwrite a field:

```javascript
{ $set: { name: { $toLower: "$name" } } }
```

Expressions use aggregation operators: `$concat`, `$ifNull`, `$cond`, `$dateToString`, and many more. The first argument of many expressions is a field path with `$`.

`$addFields` does not filter. It does not group. It only computes.

Do not compute in the client if the pipeline already streams the data and the formula is simple. Do not put heavy unique business logic only in a pipeline if the same formula must live in two languages.

### Questions

#### Theoretical questions

1. How is `$addFields` different from `$project`?
2. What is `$set` in an aggregation pipeline?
3. What does a field path like `"$total"` mean?
4. Does `$addFields` remove fields that you do not name?
5. When do you compute in the application instead of the pipeline?

#### Easy practical tasks

1. Add `nTags` as `{ $size: "$tags" }`. Write two output documents.
2. `$set` a lowercase `emailNorm`. Write the new field.
3. Overwrite `qty` with `{ $add: [ "$qty", 1 ] }`. Write the values.
4. Compare one `$project` that keeps all needed fields with one `$addFields`. Write the output keys.

#### Medium practical tasks

1. Use `$cond` to add `level: "low"` or `"high"` from `qty`. Write the expression.
2. Use `$ifNull` to replace a missing `city` with `"unknown"`.
3. Chain `$addFields` then `$match` on the new field. Write why `$match` on a computed field may not use an index.

#### Advanced practical tasks

1. Read the aggregation expression reference. Pick five operators that you did not use. Write one sentence each.
2. Rewrite a small client-side loop (tax, totals) as a pipeline. Compare results on ten documents.

---

## `$facet`

`$facet` runs several pipelines on the **same** input documents. Each subpipeline produces an array in the result.

```javascript
{
  $facet: {
    byStatus: [
      { $group: { _id: "$status", n: { $sum: 1 } } }
    ],
    top: [
      { $sort: { total: -1 } },
      { $limit: 3 }
    ]
  }
}
```

The stage output is one document (or one per incoming document group, but the usual use is after a shared `$match`). Each facet key holds an array of that subpipeline's results.

Use `$facet` for a dashboard: counts, top lists, and a sample, in one round trip.

Costs: the server holds the facet inputs and the facet outputs. Large facets use memory. Each facet is independent. A `$limit` in one facet does not limit another facet.

Do not use `$facet` to replace two cheap finds if the team cannot read the pipeline. Clarity matters.

`$facet` subpipelines cannot use `$out` or `$merge` in older rules. Read the current restriction list.

### Questions

#### Theoretical questions

1. What does `$facet` return?
2. Do facet subpipelines share a `$limit`?
3. Why can `$facet` use a lot of memory?
4. When is a dashboard a good `$facet` use?
5. Why might two finds be clearer than one `$facet`?

#### Easy practical tasks

1. Run `$facet` with two groups: count by `status` and count by `userId`. Write the output keys.
2. Add a facet that `$limit` 1. Write the array length of that key.
3. Put `$match` before `$facet`. Write why that match applies to all facets.
4. Open the `$facet` page. Write one restriction.

#### Medium practical tasks

1. Build a facet: total count, average total, and top 5 orders. Use one `$match`.
2. Compare `$facet` vs two `aggregate` calls. Write time and readability.
3. Cause a large facet (wide `$push`). Record a memory error if you hit the limit, or write the size.

#### Advanced practical tasks

1. Read `$facet` vs `$unionWith`. Write when you use each.
2. Design a report endpoint with three facets. Write the pipeline and the index that the leading `$match` needs.

---

## Aggregation vs many small finds

You can load data with many `find` calls in the application. You can also load data with one pipeline.

Prefer **find** when:

- You need documents as they are stored
- You have a filter, sort, and projection
- You page a list
- The logic is easier to test in the application

Prefer **aggregation** when:

- You group, sum, or reshape
- You must unwind and regroup
- You need a join-like `$lookup` in one trip
- You compute a report

Many small finds are a problem when:

- You run a find per document in a loop (N+1)
- You pull large documents and then filter in memory
- You implement group-by in the client on a huge array

Aggregation is a problem when:

- The pipeline is a second programming language that only one person can change
- You hide business rules in a pipeline that the app also implements
- You `$lookup` graphs that should be a model change

A good default: find for OLTP reads, aggregation for reports and computed views. If a report runs on every user click, store a computed field (next topic pattern).

### Questions

#### Theoretical questions

1. When is `find` the better tool?
2. What is an N+1 find pattern?
3. When does aggregation beat client-side group-by?
4. Why can a pipeline be hard to maintain?
5. When do you store a computed field instead of aggregating on each click?

#### Easy practical tasks

1. Write the same "orders for user" as `find` and as `{ $match }`. Write which you keep.
2. Write an N+1 example in words: list orders, then find each user.
3. Replace that N+1 with one `$lookup` or with an embed. Write the choice.
4. Make a two-column table: "Use find" and "Use aggregate". Add four rows.

#### Medium practical tasks

1. Implement a count-by-status in the client (`toArray` + loop) and in `$group`. Compare times on 10 000 documents.
2. Review a sample app sketch with three endpoints. Mark find vs aggregate for each.
3. Show a pipeline that is harder to read than two finds. Rewrite it as finds or as a simpler pipeline.

#### Advanced practical tasks

1. Read about materialized views or `$merge`. Write when you precompute a report.
2. Profile one aggregation vs one find in Atlas or `explain`. Write `executionStats` for both.

---

## Memory limits and `allowDiskUse`

A pipeline stage has a memory limit. For many versions the limit is **100 megabytes** of RAM per stage unless you allow disk.

```javascript
db.orders.aggregate(pipeline, { allowDiskUse: true })
```

`$group`, `$sort`, and large `$facet` or `$unwind` streams can hit the limit. The error mentions memory or `allowDiskUse`.

`allowDiskUse` lets the stage spill to temporary files. The query can finish. It can be much slower. Disk spill is not a free fix.

Better fixes:

- `$match` earlier
- Project fewer fields
- Avoid `$unwind` of huge arrays
- Index so `$sort` is not a huge memory sort
- Change the model
- Precompute

Atlas and `mongod` also have other limits (time, document size). A 16 MB document limit still applies to documents in the stream in many cases. Very wide `$push` arrays can hit document size.

Do not set `allowDiskUse: true` on every OLTP request. Use it for reports.

`explain` can show disk use in some versions. Test with production-like data size.

### Questions

#### Theoretical questions

1. What is the usual per-stage RAM limit without disk?
2. What does `allowDiskUse` do?
3. Which stages often hit the limit?
4. Why is disk spill not the first fix?
5. Why is `allowDiskUse` a poor default for a user-facing request?

#### Easy practical tasks

1. Find the memory-limit page in the manual. Write the number and the URL.
2. Run a small `$group` with `allowDiskUse: true`. Confirm it works.
3. List three pipeline changes that reduce memory.
4. Write the `aggregate` call form that passes options.

#### Medium practical tasks

1. Build a `$sort` on a large collection without an index. Record whether you need `allowDiskUse`.
2. `$push` many elements in `$group` to approach a large document. Record the error if you hit 16 MB.
3. Compare a pipeline with early `$project` vs late `$project` on a wide document. Write a qualitative memory difference.

#### Advanced practical tasks

1. Read about aggregation memory and 64-bit vs limits in your version. Write any extra Atlas limit that you find.
2. Design a report that cannot use `allowDiskUse` on the request path. Write a nightly `$merge` plan.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of one order document through `$match`, `$unwind`, `$group`, and `$sort`.
2. How do `$lookup` and `$facet` each combine work that you could do with more than one query?
3. When do `$addFields` and `$project` lead to a later `$match` that cannot use an index?
4. How do memory limits change the choice between aggregate, find, and a stored computed field?
5. A teammate writes a 20-stage pipeline for a login page. Which facts do you use to reject it?

#### Easy practical tasks

1. Write one pipeline that matches, groups, sorts, and limits. Run it. Write the result.
2. Write a cheat sheet: stage names in this topic, one line each, plus `allowDiskUse`.
3. `$lookup` then `$unwind` the `as` array. Write the field shape after unwind.
4. Run `explain` on one pipeline. Save the output.

#### Medium practical tasks

1. Build a dashboard pipeline with `$facet` (counts and top N) after a `$match`.
2. Replace an N+1 find loop with one aggregation. Write both. Compare times.
3. Cause or describe a memory error, then fix it with `$match` or `allowDiskUse`. Write which fix you chose and why.

#### Advanced practical tasks

1. Write a five-stage pipeline that uses `$lookup`, `$unwind`, and `$group`. Add indexes that the `$match` and lookup need. Prove with `explain`.
2. Read `$merge` or `$out`. Write a job that stores a daily summary collection. Include the pipeline and the unique key of the summary.
