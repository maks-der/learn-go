# 1. Getting Started

## Description

MongoDB is a document database. A document database stores data as documents, not as rows in tables. This topic shows what MongoDB is, when a document model fits, how you install MongoDB, and how you use `mongosh`. You also learn databases, collections, documents, BSON types, and the document size limit.

Learn common database ideas first (`db.topics.md`). Then complete this topic. Complete this topic before you write queries.

Use one term for each concept. A **database** holds collections. A **collection** holds documents. A **document** is one BSON object with fields. Do not call a collection a table. Do not call a document a row.

---

## What MongoDB is and when documents fit

MongoDB is a document database management system. The server process is `mongod`. Clients send commands to `mongod`. Atlas is the MongoDB cloud service. The same document model applies to a local server and to Atlas.

A document is a set of field-and-value pairs. Values can be strings, numbers, dates, arrays, and nested documents. Documents in one collection do not need the same fields. The application still has a schema. The database does not force one schema unless you add validation.

MongoDB stores documents in BSON. BSON is a binary format. You write queries in a JSON-like syntax. The server stores and returns BSON.

MongoDB is not a relational database. It does not use SQL as the primary language. It does not require joins for every relationship. You can still model relationships. You embed related data, or you store a reference.

MongoDB is not a key-value store only. A document can have many fields. You can query those fields. You can index those fields.

Typical uses:

- Operational data for web and mobile applications
- Content that has a nested shape
- Event data and catalogs
- Services that already exchange JSON

A document model fits when the application reads and writes a nested object as one unit. Example: an order with line items. The order and the lines are one document. One read returns the full order.

A document model also fits when the shape of data changes often. New fields do not need a table migration for every change. You still plan field names. You still plan types.

A relational model fits better when many entities share data and you must keep that data in one place. Example: a product name that many orders must show as the live name. A relational foreign key and a join keep one product row.

A relational model also fits better when you need many ad-hoc joins, strict multi-row constraints, or heavy reporting across many tables.

Use documents when:

- Data that you access together has a clear nest
- The document stays below the size limit
- Arrays stay bounded
- The application owns the schema

Use a relational model when:

- Many-to-many links are the center of the model
- You need declarative referential integrity as the main control
- Analysts run many new join queries
- Transactions across many independent rows are the default write

You can use both. Many systems store operational documents in MongoDB and store reports in a relational or warehouse system.

### Questions

#### Theoretical questions

1. What is a document database?
2. What process name does a MongoDB server use?
3. When does a nested document replace a parent table and a child table?
4. How is MongoDB different from a key-value store that holds only one blob per key?
5. When is a live shared entity better as a relational row?

#### Easy practical tasks

1. Write five sentences that describe MongoDB. Use only facts from this section.
2. Make a two-column table: "Fits documents" and "Fits relational". Add four rows.
3. List three application types that fit a document database. Give one reason for each choice.
4. Open [https://www.mongodb.com/docs/](https://www.mongodb.com/docs/). Write the current major version that the manual shows.

#### Medium practical tasks

1. Compare MongoDB with one relational database that you know. Write six short sentences. Cover data shape, schema, and query language.
2. Model a shop: customer, order, product. Write which parts you embed and which parts you reference. Give one reason for each choice.
3. Draw two diagrams of the same domain: one document tree, one set of tables. Label the difference in a read of one order.

#### Advanced practical tasks

1. Read the "Introduction to MongoDB" pages in the official manual. Write a one-page timeline of MongoDB products (Community, Enterprise, Atlas) with one fact each.
2. Take a small relational schema (five tables). Redesign it as MongoDB collections. Write which joins disappear and which references remain.

---

## Installing Community, Atlas, or Docker

You can run MongoDB in three common ways.

**MongoDB Community Server.** Download the installer from the MongoDB download page. Select your operating system. Install the server. Start the `mongod` service. Community Server is free. You manage the host, the disk, and the upgrades.

**Atlas.** Atlas is the official cloud service. Create a project and a cluster. The free tier is enough for learning. Atlas manages the server. You connect with a connection string. You must set a database user and an IP access list.

**Docker.** Docker runs `mongod` in a container. A typical learning command uses the official `mongo` image and publishes port `27017`. Data in the container is lost if you do not mount a volume. Use a volume when you want data to stay.

For this path, pick one method and keep it. Do not mix three environments without a note of which one you use.

After the server runs, install `mongosh` if the installer did not add it. `mongosh` is the official shell.

Default port is `27017`. Do not expose that port to the public internet on a learning machine without authentication.

### Questions

#### Theoretical questions

1. What is the difference between Community Server and Atlas?
2. What does Docker give you that a local installer also gives you?
3. Why must you set a user and an IP list on Atlas?
4. What happens to data in a Docker container if you do not mount a volume?
5. What is the default MongoDB port?

#### Easy practical tasks

1. Install Community Server, or create an Atlas cluster, or start a Docker container. Write which method you used.
2. Confirm that port `27017` listens (local or Docker) or that the Atlas cluster shows "available".
3. Write the install path or the Atlas connection host in a text file. Do not write passwords in a shared file.
4. Open the official install page for your operating system. Write the package name that you used.

#### Medium practical tasks

1. Start MongoDB with Docker and a named volume. Stop the container. Start it again. Confirm that a test database is still there.
2. Create an Atlas free cluster. Add your IP. Create a database user with only the rights that you need for learning.
3. Compare Community Server and Atlas in a short table: who manages backups, who manages TLS, who manages upgrades.

#### Advanced practical tasks

1. Run Community Server as a replica set with one member (required later for transactions). Record the config file or the command flags.
2. Compare Docker Desktop, a Linux VM, and Atlas for a beginner. Write a short report: start time, cost, and reset method.

---

## `mongosh`

`mongosh` is the MongoDB Shell. It is a command-line client. You use it to run commands, to inspect data, and to test queries. The old `mongo` shell is not the current tool. Use `mongosh`.

Connect to a local server:

```text
mongosh
```

Connect with a URI:

```text
mongosh "mongodb://localhost:27017"
```

Connect to Atlas with the URI that Atlas shows. The URI includes the user name. The shell asks for the password, or you put credentials in the URI. Do not commit a URI with a password.

After you connect, you are in a JavaScript environment. You select a database with `use`:

```text
use learn
```

`db` is the current database. `db.help()` shows common helpers. `show dbs` lists databases. `show collections` lists collections in the current database.

A simple insert and find:

```javascript
db.notes.insertOne({ title: "hello", n: 1 })
db.notes.find()
```

`mongosh` prints documents as extended JSON. You can assign results to variables. You can write small scripts. For daily learning, type commands one by one.

Exit with `exit` or `Ctrl+C` as the shell documents.

The web manual is the source of truth for options and limits. Use the manual that matches your major version. Do not trust random blog snippets for write-concern or transaction rules.

### Questions

#### Theoretical questions

1. What is `mongosh`?
2. Why must you use `mongosh` and not the old `mongo` shell for new work?
3. What does the `use` command do?
4. What does `db` refer to in the shell?
5. Why is a connection URI with a password a risk in a shared file?

#### Easy practical tasks

1. Connect with `mongosh`. Run `db.hello()` or `db.runCommand({ ping: 1 })`. Save the output in a text file.
2. Run `show dbs`. Write the database names that you see.
3. Run `help` or `db.help()`. Write three helper names.
4. Insert one document into `db.notes`. Find it. Show the printed document.

#### Medium practical tasks

1. Connect to two targets (local and Atlas, or two local ports). Write the two URI forms. Do not include passwords in your notes if you share them.
2. Write a small `.mongodb.js` script that inserts two documents and counts them. Run it with `mongosh` as the manual shows.
3. Use `mongosh --eval` to run one find command from the system shell. Record the full command.

#### Advanced practical tasks

1. Enable authentication on a local `mongod`. Connect with a user in `mongosh`. Record the extra URI options.
2. Compare `mongosh` snippets with the same operations in one official driver. Write which types look different (dates, ObjectId).

---

## Databases, collections, and documents

A MongoDB **deployment** holds one or more databases. A **database** holds collections and views. A **collection** holds documents. A **document** is the unit that you insert, update, and find.

Names:

- Database names are strings. Do not use reserved names such as `admin`, `local`, and `config` for application data.
- Collection names are strings. Use a clear plural noun, for example `orders`.
- Field names are strings. Do not start application field names with `$`. Do not use `.` in a field name.

MongoDB creates a database when you first store data in it. MongoDB creates a collection when you first insert a document, unless you create the collection on purpose.

Each document has an `_id` field. If you omit `_id`, MongoDB adds an ObjectId. `_id` is unique in the collection.

A collection is not a table. There is no fixed column list. Two documents can look different. That flexibility is not a reason to store random shapes without a plan.

Commands in `mongosh` use `db.collectionName.method()`. Example: `db.orders.insertOne({ ... })`.

Drop a collection with `db.orders.drop()`. Drop a database with `db.dropDatabase()`. Those commands remove data. Do not run them on a shared cluster without care.

### Questions

#### Theoretical questions

1. What is the containment order: deployment, database, collection, document?
2. When does MongoDB create a database?
3. What field must each document have?
4. Why must you not use `admin` for application data?
5. How is a collection different from a relational table?

#### Easy practical tasks

1. Run `use learn`. Insert one document into `items`. Run `show collections`. Confirm `items` exists.
2. Insert two documents with different fields into the same collection. Find both.
3. Run `db.getName()`. Write the current database name.
4. Draw a tree: deployment → two databases → collections → one sample document each.

#### Medium practical tasks

1. Create a collection with `db.createCollection("typed")`. Insert one document. Compare this with an insert that creates the collection.
2. List databases with `db.adminCommand({ listDatabases: 1 })`. Write the size fields that you see.
3. Drop a test collection. Confirm that `show collections` no longer lists it.

#### Advanced practical tasks

1. Read the naming-restriction page in the manual. Write five illegal or reserved name rules.
2. Compare `local`, `admin`, and `config` in the manual. Write one sentence for the role of each system database.

---

## JSON vs BSON, field types, `_id`, and the 16 MB limit

JSON is a text format for objects, arrays, strings, numbers, booleans, and null. Humans read JSON. Many APIs send JSON.

BSON is a binary encoding. MongoDB uses BSON on disk and on the wire. BSON includes a type byte and a length for each value. The server does not parse JSON text for each stored document.

BSON has types that JSON does not have. Examples: ObjectId, Date, Binary, Decimal128, Timestamp, Regular Expression. `mongosh` prints those types as extended JSON. Extended JSON uses wrappers such as `{ "$oid": "..." }` and `{ "$date": "..." }`.

A JSON number is not a full match for BSON numbers. BSON has 32-bit integers, 64-bit integers, and 64-bit floating-point values. Drivers map language numbers to BSON types. A decimal money value needs `Decimal128`, not a binary float.

Do not edit raw BSON by hand. Use the shell, a driver, or a tool.

Common BSON field types:

- **String.** UTF-8 text. Use strings for names, codes, and tokens.
- **Number.** Integer or double. Use `NumberInt`, `NumberLong`, or a normal number in `mongosh`. Use `NumberDecimal` for exact decimals.
- **Bool.** `true` or `false`.
- **Date.** A UTC datetime. Use `ISODate("2026-01-15T12:00:00Z")` in `mongosh`. Do not store dates as local strings if you must sort them as time.
- **ObjectId.** A 12-byte identifier.
- **Array.** An ordered list of values. Prefer one type in an array when you can.
- **Embedded document.** A nested object. Example: `address: { city: "Oslo", zip: "0001" }`.
- **Null.** A missing value. Null is not the same as a missing field.
- **Binary.** Raw bytes. Use binary for small blobs. Do not use binary for large files.

Type matters for queries and indexes. The number `1` and the string `"1"` are not equal. Pick a type for each field and keep it.

Every document has `_id`. `_id` is unique in the collection. MongoDB creates an index on `_id`. You cannot drop that index.

If you omit `_id` on insert, the driver or the server adds an ObjectId. You can set `_id` to another type. The value must be unique. The value must be stable.

ObjectId is 12 bytes: a timestamp in seconds, a random value, and a counter. The timestamp part lets you sort by insert time in a rough way. Do not treat ObjectId as a strict clock for all nodes.

`_id` is immutable in practice. You do not change `_id` with `$set`. To use a new `_id`, you insert a new document and you delete the old document. Keep `_id` small.

One MongoDB document has a maximum size of 16 megabytes. The limit includes BSON overhead. The limit protects the server and the client. A document that grows without a bound will hit the limit.

Do not store large files as one document. Use GridFS or object storage for large files. `Object.bsonsize(doc)` in `mongosh` returns the BSON size in bytes for an in-memory document.

The 16 MB limit is not the limit for a collection. A collection can hold many documents. A cursor returns batches.

### Questions

#### Theoretical questions

1. What is the difference between JSON and BSON?
2. Why is `NumberDecimal` better than a double for money?
3. What happens if you insert a document without `_id`?
4. What is the maximum size of one document?
5. Does the 16 MB limit apply to a collection?

#### Easy practical tasks

1. In `mongosh`, insert `{ n: 1, d: ISODate() }`. Find the document. Write how the shell prints the date.
2. Insert `{ n: 1 }` and `{ n: "1" }`. Query `{ n: 1 }`. Write which document matches.
3. Insert a document without `_id`. Print `_id` and its type. Then run `Object.bsonsize` on that document.
4. Open the official BSON page and the document size-limit page. Write both URLs.

#### Medium practical tasks

1. Insert a date with `ISODate` and the same instant as a string. Sort the collection by that field. Write what happens.
2. Query for `{ color: null }` on a collection that has a missing `color` and a `color: null`. Write how many documents match. Then write a query that uses `$type`.
3. Compare ObjectId and a UUID string as `_id`. Write two advantages for each.

#### Advanced practical tasks

1. Explain why money must not use a binary float. Insert `0.1 + 0.2` as a double and as `NumberDecimal`. Compare the stored values.
2. Try to insert a document that is larger than 16 MB (generate it in a script). Record the error. Do not keep the huge payload.

---

## Schema-less does not mean schema-free

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

1. Describe the path from an empty machine to a first inserted document. Name the install choice, the client, and the two `mongosh` commands.
2. A teammate says "MongoDB stores JSON" and calls a collection a table. Which facts do you use to correct both terms?
3. How do Community Server, Atlas, and Docker differ in who starts `mongod`?
4. How do the 16 MB limit, field types, and `_id` rules work together when you design one document?
5. What is the difference between a flexible schema and no schema?

#### Easy practical tasks

1. Create database `start`. Insert `{ ok: true }` into `smoke`. Find the document. Write the `_id` that you see.
2. Write a one-page cheat sheet: `mongosh`, `use`, `show dbs`, `show collections`, `insertOne`, `find`, default port, JSON vs BSON, 16 MB.
3. Export or copy your connection method (local host or Atlas host) into a private note. List the database name that you use for practice.
4. Insert one complete sample document for a `books` collection. Use string, number, bool, date, array, and an embedded document. Write `Object.bsonsize`.

#### Medium practical tasks

1. Write a small script (PowerShell or bash) that checks `mongosh --version` and then runs a ping. Stop if the shell is missing.
2. Reset your learning database: drop it, create it again, insert two documents. Record each command.
3. Export one collection to JSON (`mongoexport` or a GUI). Open the file. Mark extended JSON keys. Import it back into a new collection.

#### Advanced practical tasks

1. Connect the same Atlas cluster or the same local server from `mongosh` and from one GUI client (Compass or similar). Show one find in each. Compare the displayed types.
2. Compare ObjectId, UUID, and a natural key for `_id` in a one-page table: size, sort order, uniqueness source, and driver support.
