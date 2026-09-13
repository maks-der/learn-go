# 6. Validation and Schema Patterns

## Description

This topic shows how you validate documents in MongoDB and how you apply common schema patterns. You attach a JSON Schema to a collection. You set `validationLevel` and `validationAction`. You compare application validation with database validation. You learn the attribute, bucket, outlier, computed, subset, and extended-reference patterns. You add a schema version field.

Complete data modeling first. Complete this topic before you study drivers and production schema changes.

Use one term for each concept. A **validator** is a rule on a collection. **JSON Schema** is the usual validator language. A **pattern** is a design that you apply on purpose. Do not apply every pattern to every collection.

---

## JSON Schema validation, `validationLevel`, `validationAction`

You can create or collMod a collection with a validator:

```javascript
db.createCollection("users", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["email", "createdAt"],
      properties: {
        email: { bsonType: "string", pattern: "^.+@.+$" },
        createdAt: { bsonType: "date" },
        age: { bsonType: "int", minimum: 0 }
      },
      additionalProperties: true
    }
  }
})
```

`$jsonSchema` describes BSON types with `bsonType`. Use `string`, `int`, `long`, `double`, `decimal`, `bool`, `date`, `objectId`, `array`, `object`, `null`, `binData`.

`required` lists fields that must exist. `properties` describes named fields. `additionalProperties: false` rejects unknown fields. That setting is strict. Use it when the schema is stable.

You can use `oneOf` for versions or variants. You can describe arrays with `items`.

Existing collections: `db.runCommand({ collMod: "users", validator: { $jsonSchema: { ... } } })`.

A validator that is too strict will block writes. Test on a copy. Insert valid and invalid documents.

The schema is not a full application model. It does not replace unique indexes. It does not replace authorization.

**validationLevel** controls which documents the server validates.

- `strict` — validate all inserts and all updates. This is the default.
- `moderate` — validate inserts and validate updates only if the existing document already matches the schema. Documents that were invalid before the validator can still update in some cases. Read the current manual for the exact moderate rule.

Use `strict` for new collections. Use `moderate` when you add a validator to a collection that has old invalid documents and you must still update them.

**validationAction** controls what the server does on a failed check.

- `error` — reject the write. This is the default.
- `warn` — accept the write and log a warning.

Use `warn` when you introduce a schema and you want to find violations without blocking production writes. Watch the logs. Fix the application. Then set `error`.

```javascript
db.runCommand({
  collMod: "users",
  validator: { $jsonSchema: { /* ... */ } },
  validationLevel: "moderate",
  validationAction: "warn"
})
```

Do not leave `warn` forever. Invalid data will grow.

`strict` plus `error` is the usual end state for a stable collection.

### Questions

#### Theoretical questions

1. What does `$jsonSchema` describe?
2. What does `additionalProperties: false` do?
3. What is the default `validationLevel`?
4. When do you use `moderate`?
5. What is the difference between `error` and `warn`?

#### Easy practical tasks

1. Create `users` with a schema that requires `email` as a string. Insert a valid document. Then insert a document without `email`. Record the error.
2. Insert `age` as a string when the schema wants `int`. Record the error.
3. Set `validationAction: "error"`. Prove that a bad insert fails.
4. Open the JSON Schema validation page. Write the URL.

#### Medium practical tasks

1. Set `additionalProperties: false`. Insert an extra field. Record the result. Then allow extra fields.
2. Create invalid old documents. Add a strict validator. Try to update one old document. Then change to `moderate` and retry.
3. Plan a rollout: warn → fix → error. Write the three deploy steps.

#### Advanced practical tasks

1. Write `oneOf` for two event types (`login`, `purchase`) with different required fields. Test two inserts.
2. Read about `$jsonSchema` vs query expressions as a validator. Write one rule that is easier as a query expression.

---

## Application validation vs database validation

**Application validation** runs in the client or the API. Examples: a JSON Schema library, a type decoder, Mongoose, Bean Validation. It checks input before the write. It returns a clear error to the user.

**Database validation** runs in `mongod`. It checks every writer: the API, a script, `mongosh`, a migration. It is the last guard.

Use both.

The application:

- Validates user-facing rules (password strength, email confirm)
- Gives field-level errors
- Can use richer types than the collection schema

The database:

- Blocks bad writes from tools and forgotten scripts
- Enforces type and required fields
- Does not depend on one language

Do not rely on the application only. A future job will write without that code.

Do not rely on the database only. The user will see a driver error, not a form error. The schema cannot express every business rule.

Keep the two schemas close. If `age` is `int` in MongoDB, the API must not send a float that becomes a different BSON type.

Drivers can send types that pass the application but fail `bsonType`. Test the real BSON, not only JSON.

Stable names and types keep queries, indexes, and validators working. Pick a field name and keep it. Pick a BSON type and keep it. If you must change a name or a type, use expand-contract and `schemaVersion`.

### Questions

#### Theoretical questions

1. What does application validation do that the database does not do well?
2. What does database validation catch that the application can miss?
3. Why use both?
4. Why can a JSON API value fail `bsonType` after the driver encodes it?
5. Why is a driver error a poor user-facing message?

#### Easy practical tasks

1. Write three rules that belong in the API and three rules that belong in `$jsonSchema`.
2. Insert from `mongosh` a document that your app would reject. Show that the database still needs a validator.
3. Map one field: TypeScript or Go type → BSON type → `bsonType` string.
4. Make a two-column table: "Application" and "Database". Add four rows.

#### Medium practical tasks

1. Implement a small validator in your language for `email` and `createdAt`. Then add the same rules as `$jsonSchema`. Test both paths.
2. Find one rule that you cannot express in `$jsonSchema` (for example "email is unique"). Write which tool enforces it (unique index).
3. Show a number that JSON accepts (`1.0`) and write the BSON type your driver sends. Check if the schema allows it.

#### Advanced practical tasks

1. Generate `$jsonSchema` from an application type (or the reverse) with a tool or a short script. Write the gaps.
2. Design a policy: which team can change the collection validator, and how the API schema stays in review. Write five rules.

---

## Attribute, bucket, outlier, computed, subset, and extended-reference patterns

A **pattern** is a known way to shape documents. Official pattern names come from MongoDB documentation and talks. Use those names so that a team can talk with one vocabulary. Name the pattern in your notes when you use it.

**Attribute pattern.** Store many optional predicates as an array of key-value documents.

```javascript
{
  attrs: [
    { k: "color", v: "red" },
    { k: "size", v: "M" }
  ]
}
```

You index `{ "attrs.k": 1, "attrs.v": 1 }`. A query for color red uses `$elemMatch` on `k` and `v`. Use this pattern when many documents have different optional fields and you must index those fields. Do not use this pattern for a small, stable set of fields.

**Bucket pattern.** Group many small events into one document. The document is a **bucket** for a time range or a size limit.

```javascript
{
  deviceId: "d1",
  start: ISODate("2026-01-15T10:00:00Z"),
  n: 3,
  readings: [
    { t: ISODate("2026-01-15T10:01:00Z"), c: 21.1 },
    { t: ISODate("2026-01-15T10:02:00Z"), c: 21.2 }
  ]
}
```

Benefits: fewer documents, a smaller `_id` index, one read for a window. You must bound the bucket. Close the bucket when `n` reaches a maximum or when the time window ends. MongoDB time series collections use a related idea inside the server.

**Outlier pattern.** Keep the common case in the main document. Move a rare huge case to another collection. Example: most books have fewer than 20 coauthors. A few books have 500 contributors. You embed authors up to 20. You set `hasExtraAuthors: true` and you store the rest in `book_authors`. Most reads do not need the extra query. Measure the percentiles. If the 99th percentile is still huge, embed is the wrong default.

**Computed pattern.** Store a value that you can calculate from other fields. You update the stored value when the source changes. Example: store `commentCount` on a post. Increment it when you insert a comment. Do not `$count` all comments on every page view. The cost is that writes must keep the value correct. A bug can make the stored value stale. You may need a repair job.

**Subset pattern.** Store a **small copy** of related data for the common read. The full document lives in another collection. Example: a product list stores `{ productId, name, thumbUrl, price }`. The `products` collection stores the long description. Updates to a source field that is in the subset must update the copies, or you accept staleness. Do not duplicate the entire document.

**Extended reference.** Store an id **and** a few fields from the referenced document.

```javascript
{
  productId: ObjectId("..."),
  name: "Nail",
  price: NumberDecimal("1.50")
}
```

The id is the live link. The extra fields are a snapshot or a cache for the common read. This pattern sits between a bare reference and a full embed. Decide if the extra fields are a **snapshot** (order history) or a **cache** (can refresh). Do not extend the reference with twenty fields.

Do not apply all six patterns to one collection.

### Questions

#### Theoretical questions

1. What problem does the attribute pattern solve?
2. When do you close a bucket?
3. Why does a flag exist on the main document in the outlier pattern?
4. Why increment `commentCount` instead of counting on read?
5. What two parts does an extended reference have?

#### Easy practical tasks

1. Store two products with different `attrs`. Query color red with `$elemMatch`.
2. Write two bucket documents for two hours of one device. `$push` one reading and `$inc` `n`.
3. Write a book with 3 embedded authors and a book with a flag and extra authors.
4. Write an order line as a bare id, as an extended reference, and as a full embed. Three snippets.

#### Medium practical tasks

1. Compare three single-field indexes vs one attribute index for ten optional fields. Write index count.
2. Choose a bucket key: `deviceId` + hour. Write the unique index. Compare document count for 10 000 events as one-per-event vs 100 events per bucket.
3. Write a repair aggregation that sets `commentCount` from a `comments` collection. Then write an update path when `name` changes in a subset.

#### Advanced practical tasks

1. Read the official building-with-patterns series. Map each official pattern in this section to a URL and one extra fact.
2. Design a catalog with 200 possible predicates, buckets for click events, and outliers for a viral post's comments. Write bounds, flags, and write paths.

---

## Schema version field

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

Validation can allow more than one version with `oneOf`.

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
3. Sketch `oneOf` for v1 and v2 (fields only). Run the validator if you can.

#### Advanced practical tasks

1. Read expand-contract for databases (MongoDB or general). Write a one-page plan for a field type change (string to object).
2. Design a dual-write period: writers write v2, a job converts v1. Write the stop condition for the job.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `$jsonSchema`, `validationLevel`, and `validationAction` work together from first rollout to a stable collection?
2. Which rules belong in the application, which belong in the validator, and which belong in a unique index?
3. Which patterns reduce document size on the hot path (outlier, subset, bucket), and how do they differ?
4. How does `schemaVersion` support a change that uses the computed or subset pattern?
5. A teammate says "MongoDB has no schema" and applies all six patterns to one collection. Which facts do you use to correct both ideas?

#### Easy practical tasks

1. Create `products` with a small `$jsonSchema`, `strict`, and `error`. Insert one good document and one bad document.
2. Write a cheat sheet: bsonType names, level, action, app vs db, six pattern names, `schemaVersion`.
3. Name the pattern for: last 100 metrics in one hour document; `likeCount` on a post; product id plus name on a line.
4. Add `schemaVersion: 1` to one sample document for each of two collections that you use.

#### Medium practical tasks

1. Add a validator to a collection that already has mixed types. Use `moderate` and `warn`. Write how you find and fix the bad documents.
2. Write `oneOf` for `schemaVersion` 1 and 2 of the same collection. Test one insert per version.
3. Pick a real or invented product catalog. Apply attribute or not. Apply subset for lists. Write why.

#### Advanced practical tasks

1. Produce a validation rollout report: counts of warn-log violations, top failing fields, date you switch to error.
2. Review a model from topic 3. Add at most three patterns. Write which problems they solve and the write paths you added.
