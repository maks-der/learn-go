# 1. Getting Started

## Description

MongoDB is a document database. A document database stores data as documents, not as rows in tables. This topic shows what MongoDB is, when a document model fits, how you install MongoDB, and how you use `mongosh`. Complete this topic before you write queries.

Use one term for each concept. A **database** holds collections. A **collection** holds documents. A **document** is one BSON object with fields. Do not call a collection a table. Do not call a document a row.

Learn common database ideas first. Then follow this MongoDB path.

---

## What MongoDB is (document database)

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

### Questions

#### Theoretical questions

1. What is a document database?
2. What process name does a MongoDB server use?
3. What format does MongoDB use to store a document?
4. Must all documents in one collection have the same fields?
5. How is MongoDB different from a key-value store that holds only one blob per key?

#### Easy practical tasks

1. Write five sentences that describe MongoDB. Use only facts from this section.
2. Make a two-column table: "MongoDB term" and "What it means". Add four rows.
3. List three application types that fit a document database. Give one reason for each choice.
4. Open [https://www.mongodb.com/docs/](https://www.mongodb.com/docs/). Write the current major version that the manual shows.

#### Medium practical tasks

1. Compare MongoDB with one relational database that you know. Write six short sentences. Cover data shape, schema, and query language.
2. Draw a simple diagram: client, `mongod` or Atlas, database, collection, document. Label each part.
3. Find the MongoDB release notes for the last two major versions. Write three changes that help a beginner.

#### Advanced practical tasks

1. Read the "Introduction to MongoDB" pages in the official manual. Write a one-page timeline of MongoDB products (Community, Enterprise, Atlas) with one fact each.
2. Explain why a binary document format helps a database more than plain JSON text. Give two reasons. Use the official BSON page.

---

## When documents fit and when a relational model fits better

A document model fits when the application reads and writes a nested object as one unit. Example: an order with line items. The order and the lines are one document. One read returns the full order.

A document model also fits when the shape of data changes often. New fields do not need a table migration for every change. You still plan field names. You still plan types.

A relational model fits better when many entities share data and you must keep that data in one place. Example: a product name that many orders must show as the live name. A relational foreign key and a join keep one product row.

A relational model also fits better when you need many ad-hoc joins, strict multi-row constraints, or heavy reporting across many tables. SQL and a normalized schema are strong for that work.

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

1. When does a nested document replace a parent table and a child table?
2. Why can frequent new fields be easier in a document database?
3. When is a live shared entity better as a relational row?
4. Why can unbounded growth of a document be a problem?
5. Can one system use MongoDB and a relational database together? Give one reason.

#### Easy practical tasks

1. Write three examples of data that fit a document. Write three examples that fit a relational table. Use one sentence each.
2. For a blog post with comments, write one sentence that argues for embed. Write one sentence that argues for a separate collection.
3. List two reporting questions that are easier in SQL. List two application reads that are easier as one document.
4. Make a two-column table: "Fits documents" and "Fits relational". Add four rows.

#### Medium practical tasks

1. Model a shop: customer, order, product. Write which parts you embed and which parts you reference. Give one reason for each choice.
2. Interview one teammate or read one case study. Record one reason they chose MongoDB or they chose a relational database.
3. Draw two diagrams of the same domain: one document tree, one set of tables. Label the difference in a read of one order.

#### Advanced practical tasks

1. Take a small relational schema (five tables). Redesign it as MongoDB collections. Write which joins disappear and which references remain.
2. Find an official MongoDB data-modeling page that compares embed and reference. Quote the rule about data that you access together. Write the URL.

---

## Installing MongoDB Community or using Atlas / Docker

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

## `mongosh` shell

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

## Databases, collections, documents

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

## Official docs: mongodb.com/docs

The primary documentation is [https://www.mongodb.com/docs/](https://www.mongodb.com/docs/). The manual for the server is [https://www.mongodb.com/docs/manual/](https://www.mongodb.com/docs/manual/).

Use the manual that matches your major version. A command can change between versions. Atlas has its own product docs. Drivers have their own docs.

Useful start pages:

- [MongoDB Manual](https://www.mongodb.com/docs/manual/)
- [Data modeling](https://www.mongodb.com/docs/manual/core/data-modeling-introduction/)
- [Aggregation](https://www.mongodb.com/docs/manual/aggregation/)
- [Drivers](https://www.mongodb.com/docs/drivers/)
- [MongoDB University](https://learn.mongodb.com/)

In `mongosh`, `help` and `db.help()` show local help. The web manual is the source of truth for options and limits.

Use the manual when you need exact operator rules. Use University when you learn the first time. Use the driver docs when you write application code.

Do not trust random blog snippets for write-concern or transaction rules. Read the current manual page.

### Questions

#### Theoretical questions

1. Why must the manual version match your server version?
2. Where do you read driver-specific connection options?
3. What is MongoDB University?
4. When do you open the data-modeling pages instead of a CRUD tutorial?
5. Why is a blog post a weak source for transaction rules?

#### Easy practical tasks

1. Open the manual home page. Write the version selector value that you set.
2. Open the CRUD operations page. Write the URL.
3. Open the drivers page. Write the driver name for your language.
4. Bookmark the manual, University, and the data-modeling introduction.

#### Medium practical tasks

1. Find the page for `insertOne`. Write two options that the page lists besides the document.
2. Find the Atlas getting-started page. Write three steps that differ from a local install.
3. Use the manual search for "16 MB". Write the page title that states the document size limit.

#### Advanced practical tasks

1. Compare two manual versions for one command (`find` or `update`). Write one difference, or write that the page is the same.
2. Complete one free MongoDB University intro unit. Map each lesson to a topic number in this path.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from an empty machine to a first inserted document. Name the install choice, the client, and the two `mongosh` commands.
2. A teammate calls a collection a table. Which facts do you use to correct that term?
3. How do Community Server, Atlas, and Docker differ in who starts `mongod`?
4. What is the difference between documentation on mongodb.com/docs and help text inside `mongosh`?
5. Why does this path tell you to learn vendor-neutral database ideas before MongoDB?

#### Easy practical tasks

1. Create database `start`. Insert `{ ok: true }` into `smoke`. Find the document. Write the `_id` that you see.
2. Write a one-page cheat sheet with these items: `mongosh`, `use`, `show dbs`, `show collections`, `insertOne`, `find`, default port.
3. Export or copy your connection method (local host or Atlas host) into a private note. List the database name that you use for practice.
4. Create a folder tree in your notes for later topics: `crud`, `indexes`, `agg`. Put one empty file in each. This is for your study log only.

#### Medium practical tasks

1. Write a small script (PowerShell or bash) that checks `mongosh --version` and then runs a ping. Stop if the shell is missing.
2. Reset your learning database: drop it, create it again, insert two documents. Record each command.
3. Document your install in ten steps so that another beginner can copy it. Include OS and method (Community, Atlas, or Docker).

#### Advanced practical tasks

1. Connect the same Atlas cluster or the same local server from `mongosh` and from one GUI client (Compass or similar). Show one find in each. Compare the displayed types.
2. Read the current "Install MongoDB" page for your OS. Add one production warning from that page to your notes (authentication, bind IP, or similar).
