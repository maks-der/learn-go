# 22. Next Steps

## Description

This topic closes the MongoDB path. You complete a **MongoDB University** intro path. You **model** one relational schema as documents and you compare queries. You **enable a replica set** on your machine. You add one **compound index** from a real `explain`. You read the **current manual** for your major version. Complete topics 1 to 21 first.

Use one term for each concept. A **University path** is an official course sequence on [https://learn.mongodb.com/](https://learn.mongodb.com/). A **capstone model** is one domain that you design as tables and as documents. The **current manual** is the documentation that matches the first number of `db.version()`.

After this topic, keep a real collection that you use. Read release notes when you upgrade. Do not stop at tutorials that use only `insertOne` and `find`.

---

## Complete a MongoDB University intro path

MongoDB University is the official training site. It has free intro courses. A typical intro path covers documents, CRUD, the aggregation pipeline, and indexes.

Use University to hear the same ideas with different examples. Use this handbook for STE text and for the question lists. The manual remains the source of truth for options and limits.

How to complete a path:

1. Create an account on [https://learn.mongodb.com/](https://learn.mongodb.com/).
2. Pick an intro course that matches your language if you can (or a generic M001-style intro).
3. Finish the videos and the labs.
4. Map each lesson to a topic number in this path (1 to 21).
5. Write five facts that the course added that this handbook did not repeat.

University labs often use Atlas. That is useful. You still need a local replica set for some later tasks in this topic.

Do not collect certificates instead of reading `explain` output. The certificate is optional. The skill is not.

Course names change. Search for "introduction" and for your major version. A course that targets a very old server is a poor first pick.

If you already finished an intro course, take one next course (aggregation or data modeling) rather than a second intro.

### Questions

#### Theoretical questions

1. What is MongoDB University?
2. Why is the manual still the source of truth after a course?
3. Why do you map lessons to topic numbers?
4. Why is a certificate not a substitute for `explain`?
5. Why avoid a course that targets a very old server?

#### Easy practical tasks

1. Open learn.mongodb.com. Write the title of one intro path that you will take.
2. Create an account or log in. Write the first unit name.
3. Write a two-column map: University unit, handbook topic number. Add five rows (you can fill more as you go).
4. Bookmark University, the manual, and this handbook folder.

#### Medium practical tasks

1. Complete at least one full intro unit. Write ten sentences in your own words.
2. Compare one University lab with the same task in `mongosh` on your machine. Write two differences.
3. List three course items that this handbook already covered. List two that were new.

#### Advanced practical tasks

1. Finish a full intro path. Write a one-page recap that cites topic numbers and official URLs.
2. Take a modeling or aggregation course next. Add a second map to topics 5, 7, and 8.

---

## Model one relational schema as documents and compare queries

Pick a domain that you know. Examples: shop, blog, help-desk, school.

Write a **relational** design first: five to eight tables, keys, and three SQL queries (or written joins if you do not run SQL).

Write a **document** design second: collections, embed vs reference, indexes, and the same three reads as `find` or aggregations.

Compare:

- How many round trips
- Whether a join disappeared
- Where a live shared entity still needs a reference
- Which design stays under the document size limit
- Which design needs a transaction

This exercise proves topic 5 and topic 21. It is not a contest that MongoDB always wins. If the workload is ad-hoc reporting across many entities, say that SQL is a better home for that report.

Keep both diagrams. You will reuse them when you add a replica set and an index.

Do not pick a domain with one table. That does not teach embed vs reference.

Do not embed every table "to win". That recreates a huge document.

If you completed a similar task in topic 1 or 5, reuse the domain and go deeper: add shard-key notes, a validation schema, and one anti-pattern that you rejected.

### Questions

#### Theoretical questions

1. Why do you write the relational design first?
2. What does a disappeared join mean in the document design?
3. When does a reference remain?
4. Why can SQL still win for a report?
5. Why is a one-table domain a poor exercise?

#### Easy practical tasks

1. Name the domain and list five tables.
2. Draw the document tree for one primary read (for example one order).
3. Write three English questions that both designs must answer.
4. Make a table: question, SQL idea, MongoDB idea.

#### Medium practical tasks

1. Implement both sides on a small data set (PostgreSQL or paper SQL, plus MongoDB). Run the three queries. Write counts and a time if you can.
2. Mark each relationship embed, reference, or subset. Give one reason each.
3. Add a JSON Schema for one collection. Write one document that must fail validation.

#### Advanced practical tasks

1. Add a fourth query that is awkward in MongoDB. Write whether you would add `$lookup`, a warehouse, or a relational store.
2. Write a two-page compare: transactions, indexes, and one anti-pattern that you avoided.

---

## Enable a replica set locally

A local **replica set** teaches elections, `hello`, transactions, and change streams. Atlas is a replica set too. A local set lets you stop a member and watch failover.

Minimum for learning: one member. Better: three members on three ports.

Typical ideas (Docker or processes):

- Different `dbPath` and port for each member
- The same `replSetName`
- `rs.initiate()` with a config that lists the hosts
- `rs.status()` until one member is primary

Example initiate shape:

```javascript
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "localhost:27017" },
    { _id: 1, host: "localhost:27018" },
    { _id: 2, host: "localhost:27019" }
  ]
})
```

Host names inside Docker must match what members use to reach each other. `localhost` from a container is the container, not the host. Follow a current Docker Compose example or the official tutorial.

After the set is up:

- Insert with majority write concern
- Run a two-document transaction
- Open a change stream
- Step down the primary and confirm the driver or `mongosh` finds the new primary

Do not expose these ports to the public internet.

Do not skip this task if you only have Atlas. Atlas is valid production practice. Local failover is still worth one afternoon.

If you cannot run three processes, run one member as a replica set. Transactions and change streams work. Elections with a real majority need more members.

### Questions

#### Theoretical questions

1. Why does a one-member replica set still help this path?
2. Why are three members better for learning elections?
3. What does `rs.initiate()` do?
4. Why can `localhost` fail as a host name in Docker?
5. Which three features from this path need a replica set?

#### Easy practical tasks

1. Open the replica-set deploy tutorial for your OS or for Docker. Write the URL.
2. Write the `replSetName` and the three ports that you will use.
3. After initiate, run `db.hello()`. Write `isWritablePrimary` and the hosts.
4. Make a table: task, command. Add initiate, status, stepDown, insert.

#### Medium practical tasks

1. Start the set. Insert one document with `w: "majority"`. Find it.
2. Run `rs.stepDown()`. Time the wait until writes work again. Write the seconds.
3. Open a change stream in a second shell. Insert. Print the event.

#### Advanced practical tasks

1. Build a Docker Compose file for three members and a volume each. Destroy one container. Write what `rs.status()` shows and how you recover.
2. Connect an official driver with a replica-set URI. Retry an insert during step-down. Write if the driver retried.

---

## Add one compound index from a real `explain`

Pick a `find` or an aggregation `$match` that you already use in the capstone model or in a lab collection. Run:

```javascript
db.orders.find({ userId: "u1", status: "open" }).sort({ createdAt: -1 }).explain("executionStats")
```

If the winning stage is `COLLSCAN`, or if `totalDocsExamined` is much larger than `nReturned`, design **one** compound index that matches the filter and the sort (prefix rule from topic 6).

Create the index. Run `explain("executionStats")` again. Write:

- Winning stages before and after
- `totalDocsExamined` before and after
- `executionTimeMillis` before and after (same data size)
- The index key

Use enough documents so that the difference is visible (thousands, not three).

Do not add five indexes. Add one. Prove it.

Do not create the index before the first `explain`. The point is the proof.

If the query is already covered or already `IXSCAN` with a tight examine count, pick a worse query on purpose and fix that one.

Drop the index if it is only for this lesson and it does not belong in the model. Keep it if the capstone needs it.

### Questions

#### Theoretical questions

1. Why do you run `explain` before `createIndex`?
2. What prefix rule do you use for filter plus sort?
3. Why do you need thousands of documents for a fair time compare?
4. Why add one index and not five?
5. What two `executionStats` fields do you record?

#### Easy practical tasks

1. Write the query that you will explain.
2. Run `explain("executionStats")`. Write the winning stage and `totalDocsExamined`.
3. Open the compound-index page. Write the URL.
4. Write the index document `{ field: 1, ... }` that you plan to create.

#### Medium practical tasks

1. Load enough documents. Capture before stats. Create the index. Capture after stats. Write a four-line report.
2. Show that a query that skips the prefix does not use the index well. Write that `explain`.
3. Add a projection and test if the query becomes covered. Write yes or no.

#### Advanced practical tasks

1. Use `allPlansExecution` if two indexes exist. Write why the winner won.
2. Write a team template: query, explain JSON link, index, before/after, who approved.

---

## Read the current manual for your major version

The server major version is the first number of `db.version()`. The manual is [https://www.mongodb.com/docs/manual/](https://www.mongodb.com/docs/manual/). Set the version selector to that number.

Commands, defaults, and limits change. Examples that change across versions:

- Reshard and refine
- Queryable Encryption
- Time series limits
- `explain` field names
- Deprecated shell helpers (`mongo` vs `mongosh`)

Read, do not skim only the landing page. Minimum pages for a close of this path:

- Data modeling introduction
- Replica set elections
- `explain` results
- Transactions limits
- Security checklist
- Release notes for your major version

Atlas has extra docs. Drivers have extra docs. Open the driver page for the language that you use.

Do not bookmark a random blog as the rule for write concern. Link the manual page.

When you upgrade (topic 16), read the new manual before you change FCV.

Print or save a PDF only if you must work offline. The live page is the one that gets fixes.

### Questions

#### Theoretical questions

1. How do you know which manual version to open?
2. Why can a blog be wrong about write concern?
3. Why do Atlas and drivers have separate docs?
4. When do you reread the manual after an upgrade?
5. Why is `mongosh` the current shell in the docs?

#### Easy practical tasks

1. Run `db.version()`. Write the major version.
2. Open the manual. Set the selector. Write the version that the page shows.
3. Open the six page types listed in this section. Write the six URLs.
4. Open your language's driver landing page. Write the URL.

#### Medium practical tasks

1. Pick one command that you use (`updateOne` or `watch`). Read the full page. Write two options that you had not used.
2. Read the release notes for your major version. Write three changes that affect this path.
3. Compare one page on two manual versions (if the site allows). Write one difference or write that they match.

#### Advanced practical tasks

1. Build a personal index of twenty official URLs grouped by topic 1 to 21. Keep it next to your notes.
2. Read the security checklist and the operations checklist. Add any item that this path did not cover to a "read next" list.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do University, the capstone model, a local replica set, one proven index, and the versioned manual form one closing loop?
2. Why does a document redesign without `explain` or a replica set leave gaps from topics 6, 11, and 12?
3. When do you stay on Atlas-only practice, and when must you still run a local replica set?
4. How do topics 21 and 5 show up in the relational-vs-document compare?
5. A teammate finishes videos, never runs `explain`, and reads a blog for elections. Which facts do you use in the review?

#### Easy practical tasks

1. Write a personal next-steps list of ten items. Map each item to a topic number or to an official URL.
2. Create a folder in your notes: `university`, `model`, `replset`, `explain`, `manual`. Put one file in each.
3. Write your `db.version()`, your install type (Community, Docker, Atlas), and the intro course title.
4. Draw a path from topic 1 to 22 with five boxes that you will reuse (CRUD, model, index, replica set, ops).

#### Medium practical tasks

1. Complete the five section tasks of this topic as one weekend project. Write a one-page lab report with URLs and numbers.
2. Revisit suggested practice items 8 to 10 from `mongo.topics.md` (replica set, change stream, shard key on paper). Write the outcome of each.
3. Teach a teammate one section from this topic. Write the three questions that they still had.

#### Advanced practical tasks

1. Deliver a capstone: document model, replica set, one compound index with before/after `explain`, one change stream or transaction, and a short security note (auth on).
2. Write a six-month plan: University next course, one production pattern from topic 19, one internals reread, one contribution to team runbooks.
