# 2. Documents and BSON

## Description

MongoDB stores each document as BSON. BSON is Binary JSON. This topic shows how JSON and BSON differ, how large a document can be, which field types you use, and how `_id` works. Complete this topic before you write CRUD commands.

Use one term for each concept. **JSON** is a text format. **BSON** is the storage format. A **field** is a name and a value in a document. **ObjectId** is the default type for `_id`. Schema-less storage does not mean that the application has no schema.

---

## JSON vs BSON

JSON is a text format for objects, arrays, strings, numbers, booleans, and null. Humans read JSON. Many APIs send JSON.

BSON is a binary encoding. MongoDB uses BSON on disk and on the wire. BSON includes a type byte and a length for each value. The server does not parse JSON text for each stored document.

BSON has types that JSON does not have. Examples: ObjectId, Date, Binary, Decimal128, Timestamp, Regular Expression. `mongosh` prints those types as extended JSON. Extended JSON uses wrappers such as `{ "$oid": "..." }` and `{ "$date": "..." }`.

A JSON number is not a full match for BSON numbers. BSON has 32-bit integers, 64-bit integers, and 64-bit floating-point values. JSON has one number syntax. Drivers map language numbers to BSON types. A decimal money value needs `Decimal128`, not a binary float.

BSON is not a replacement for JSON in every API. Your application can still send JSON to clients. The driver converts between language objects, JSON, and BSON.

Do not edit raw BSON by hand. Use the shell, a driver, or a tool.

### Questions

#### Theoretical questions

1. What is the difference between JSON and BSON?
2. Where does MongoDB use BSON?
3. Name three BSON types that JSON does not have.
4. What is extended JSON?
5. Why is a JSON number not enough for all BSON number types?

#### Easy practical tasks

1. Write four sentences that compare JSON and BSON. Use only facts from this section.
2. Make a two-column table: "JSON type" and "Related BSON type". Add five rows.
3. In `mongosh`, insert `{ n: 1, d: ISODate() }`. Find the document. Write how the shell prints the date.
4. Open the official BSON page. Write the URL and one sentence from the page.

#### Medium practical tasks

1. Insert one integer and one float. Use `typeof` in `mongosh` or inspect the type in Compass. Write the two BSON types.
2. Convert a document to extended JSON with `EJSON` in `mongosh` if the shell provides it, or export with `mongoexport`. Show the `$oid` or `$date` keys.
3. Compare how one official driver shows ObjectId and Date in language types. Write two lines of sample types.

#### Advanced practical tasks

1. Read the BSON specification summary in the manual. Write how a document encodes field length and type. Use your own words.
2. Explain why money must not use a binary float. Insert `0.1 + 0.2` as a double and as `NumberDecimal`. Compare the stored values.

---

## Document size limit (16 MB)

One MongoDB document has a maximum size of 16 megabytes. This limit is a hard limit for a single document. The limit includes BSON overhead, not only the text that you type.

The limit protects the server and the client. A very large document is expensive to read, write, and cache. A document that grows without a bound will hit the limit.

Do not store large files as one document. Use GridFS or object storage for large files. Do not push unbounded events into one array in one document.

If a document is near the limit, change the model. Split the data. Use references. Use a bucket pattern for time series. Those designs are later topics.

`Object.bsonsize(doc)` in `mongosh` returns the BSON size in bytes for an in-memory document. Use it when you test growth.

The 16 MB limit is not the limit for a collection. A collection can hold many documents. The 16 MB limit is not the limit for a result set. A cursor returns batches.

### Questions

#### Theoretical questions

1. What is the maximum size of one document?
2. Why does MongoDB set a maximum document size?
3. Does the limit apply to a collection?
4. Where do you store a file that is larger than 16 MB?
5. What happens if an update would make a document larger than 16 MB?

#### Easy practical tasks

1. In `mongosh`, insert a small document. Run `Object.bsonsize` on that document. Write the number.
2. Write three data types that must not live in one growing document (for example, a full chat history). Give one reason each.
3. Find the manual page that states the 16 MB limit. Write the URL.
4. Make a table: "Item" and "Limited to 16 MB?". Add rows for document, collection, and database.

#### Medium practical tasks

1. Build a document with a large string (for example 1 MB). Measure `Object.bsonsize`. Estimate how many such fields fit before 16 MB.
2. Write a short plan to split a "user with all orders" document that will grow past 16 MB.
3. Compare GridFS and an external object store in three sentences. Use the official GridFS page for facts.

#### Advanced practical tasks

1. Try to insert a document that is larger than 16 MB (generate it in a script). Record the error. Do not keep the huge payload.
2. Read how WiredTiger stores documents at a high level. Write why a huge document is bad for cache even when it is under 16 MB.

---

## Field types: string, number, bool, date, objectId, array, embedded document, null, binary

Common BSON field types:

- **String.** UTF-8 text. Use strings for names, codes, and tokens.
- **Number.** Integer or double. Use `NumberInt`, `NumberLong`, or a normal number in `mongosh`. Use `NumberDecimal` for exact decimals.
- **Bool.** `true` or `false`.
- **Date.** A UTC datetime. Use `ISODate("2026-01-15T12:00:00Z")` in `mongosh`. Do not store dates as local strings if you must sort them as time.
- **ObjectId.** A 12-byte identifier. See the next section.
- **Array.** An ordered list of values. Values can have mixed types. Prefer one type in an array when you can.
- **Embedded document.** A nested object. Example: `address: { city: "Oslo", zip: "0001" }`.
- **Null.** A missing value. Null is not the same as a missing field. A query for null can match both, unless you write a stricter query.
- **Binary.** Raw bytes. Use binary for small blobs. Do not use binary for large files.

Other types exist: Regular Expression, JavaScript, Timestamp, MinKey, MaxKey. You will see them less often in application documents.

Type matters for queries and indexes. The number `1` and the string `"1"` are not equal. A Date and an ISO date string are not the same type.

Pick a type for each field and keep it. Later topics cover validation.

### Questions

#### Theoretical questions

1. Why must you not store a datetime as a local display string if you sort by time?
2. What is the difference between a missing field and a field with value `null`?
3. Why is `NumberDecimal` better than a double for money?
4. Can one array hold values of different types? Should you do that?
5. How does an embedded document differ from a reference to another collection?

#### Easy practical tasks

1. Insert one document that uses string, number, bool, date, array, embedded document, and null. Find it.
2. Insert `{ n: 1 }` and `{ n: "1" }`. Query `{ n: 1 }`. Write which document matches.
3. Make a table: "Type" and "Example value in mongosh". Add eight rows.
4. Open the BSON types page in the manual. List two types that this section did not describe in detail.

#### Medium practical tasks

1. Insert a date with `ISODate` and the same instant as a string. Sort the collection by that field. Write what happens.
2. Query for `{ color: null }` on a collection that has a missing `color` and a `color: null`. Write how many documents match. Then write a query that uses `$type`.
3. Store a small binary with `BinData` or `Binary.createFromBase64`. Find the document. Write the type that the shell shows.

#### Advanced practical tasks

1. Write a one-page type-mapping table for your language: language type → BSON type for string, int, long, decimal, date, bool, array, object, null, binary.
2. Read about `Decimal128` precision. Write the precision and why it still is not a full money library.

---

## `_id` and ObjectId

Every document has `_id`. `_id` is unique in the collection. MongoDB creates an index on `_id`. You cannot drop that index.

If you omit `_id` on insert, the driver or the server adds an ObjectId. You can set `_id` to another type. Teams often use ObjectId, a UUID string, or a natural key. The value must be unique. The value must be stable.

ObjectId is 12 bytes:

- A timestamp in seconds
- A random value
- A counter

The timestamp part lets you sort by insert time in a rough way. Do not treat ObjectId as a strict clock for all nodes. Clocks can differ.

`mongosh` shows ObjectId as `ObjectId("...")`. The hex string is 24 characters.

`_id` is immutable in practice. You do not change `_id` with `$set`. To use a new `_id`, you insert a new document and you delete the old document.

Do not use a large embedded document as `_id`. Keep `_id` small. The `_id` index stores the key.

### Questions

#### Theoretical questions

1. What happens if you insert a document without `_id`?
2. Must `_id` be an ObjectId?
3. What are the parts of an ObjectId?
4. Why is `_id` a poor place for a large embedded document?
5. How do you "change" `_id` if the field is not updated in place?

#### Easy practical tasks

1. Insert a document without `_id`. Print `_id` and its type.
2. Insert a document with `_id: "sku-1"`. Find it with `{ _id: "sku-1" }`.
3. Extract the timestamp from an ObjectId in `mongosh` with `getTimestamp()`. Write the date.
4. Try to insert two documents with the same `_id`. Record the error.

#### Medium practical tasks

1. Compare ObjectId and a UUID string as `_id`. Write two advantages for each.
2. Insert documents with `_id` of different types (`1`, `"1"`). Find each. Explain why they are different keys.
3. Read the ObjectId page in the manual. Write whether ObjectIds from two drivers can collide in normal use. Use the page, not a guess.

#### Advanced practical tasks

1. Generate 10 000 ObjectIds in a script. Check uniqueness. Write the method and the result.
2. Design `_id` for an events collection that must be unique across shards later. Write why a monotonic field as the only `_id` can hurt distribution (high-level).

---

## Schema-less does not mean schema-free (the application still has a schema)

MongoDB does not require a schema when you create a collection. You can insert `{ a: 1 }` and `{ b: "x" }` into the same collection. That behavior is **schema-less storage**.

The application still has a schema. The application expects field names. The application expects types. The application expects arrays or embedded documents in a known shape. If the shape changes without a plan, queries break. Indexes waste space. Validation fails later.

**Schema-free** would mean that no one owns the shape. That is not how a working system operates. Someone owns the schema: the application, a schema file, or a validator on the collection.

Good practice:

- Write a document shape for each collection.
- Use stable field names.
- Use one type per field.
- Add a `schemaVersion` field when the shape must change (later topic).
- Add JSON Schema validation when the collection is stable enough (later topic).

Flexible schema is a tool. It helps during change. It is not a reason to skip design.

### Questions

#### Theoretical questions

1. What does schema-less storage mean in MongoDB?
2. Why does the application still have a schema?
3. What breaks when two documents use different types for the same field name?
4. Who can "own" the schema if the database does not force one?
5. How is a flexible schema useful during a change?

#### Easy practical tasks

1. Insert two different shapes into `db.mixed`. Write a `find` that returns only documents that have field `a`.
2. Write a one-page schema note for a `users` collection: field name, type, required or not.
3. List three bugs that a mixed-type field can cause. Use one sentence each.
4. Find a manual page that discusses schema validation or data modeling. Write how it treats "flexible schema".

#### Medium practical tasks

1. Write two versions of a `product` document (v1 and v2). List the fields that a query must handle during a migration.
2. Compare "schema in the application only" with "schema in the database validator". Write two benefits of each.
3. Review a sample collection that you created. Mark fields that already drift (name or type). Write a fix plan.

#### Advanced practical tasks

1. Read about expand-contract schema changes (later topic preview). Write a three-step change that adds `displayName` without breaking old readers.
2. Argue for or against a single `events` collection with many shapes. Use the polymorphic-collection idea. Write a one-page decision.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A teammate says "MongoDB stores JSON". Which facts do you use to correct that sentence?
2. How do the 16 MB limit, field types, and `_id` rules work together when you design one document?
3. Why can extended JSON appear in a file export but not as the on-disk format?
4. What is the difference between a flexible schema and no schema?
5. Which BSON types must you decide before you write the first insert for a money-and-time document?

#### Easy practical tasks

1. Insert one complete sample document for a `books` collection. Use string, number, bool, date, array, and an embedded document. Write `Object.bsonsize`.
2. Write a cheat sheet: JSON vs BSON, 16 MB, `_id`, ObjectId parts, null vs missing.
3. Print one ObjectId as a hex string and as `getTimestamp()`. Save the output.
4. Create `schema-notes.md` in your study folder. Describe `books` in ten lines. No code from this handbook is required.

#### Medium practical tasks

1. Export one collection to JSON (`mongoexport` or a GUI). Open the file. Mark extended JSON keys. Import it back into a new collection.
2. Write a small script that rejects a document if a required field is missing or has the wrong type. Run it on two sample documents.
3. Measure BSON size before and after you add a 100-element array of short strings. Write the two sizes.

#### Advanced practical tasks

1. Build a document that approaches 16 MB with repeated fields. Stop before the limit. Graph size versus element count. Write when you would split the model.
2. Compare ObjectId, UUID, and a natural key for `_id` in a one-page table: size, sort order, uniqueness source, and driver support.
