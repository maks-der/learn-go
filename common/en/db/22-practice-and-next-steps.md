# 22. Practice and Next Steps

## Description

This topic is a practice close of the common database path. You design a small schema, you write SQL by hand, you read a query plan, you take a backup and you restore it, then you choose a product to learn next (PostgreSQL, MongoDB, or Redis).

Use one term for each concept. Practice means you run commands on a local DBMS. A specialization is a vendor path after this vendor-neutral path. Complete this topic after you can use constraints, transactions, indexes, and backups. This topic does not teach those ideas again. It asks you to combine them.

Do not skip the SQL-by-hand step. An ORM later will emit SQL that you must read.

---

## Design a schema for a small domain (library, shop, tickets)

Pick one domain. Design tables before you write application code. Write the design on paper or in a short document. Then implement it in SQL.

**Library.** Members, copies of works, loans. A work has a title. A copy has a barcode. A loan links a member and a copy with `loaned_at` and `due_at`. A copy cannot be on two open loans.

**Shop.** Customers, products, orders, order lines. A line has quantity and unit price as exact money. An order has a customer and a status. A product has a unique sku.

**Tickets.** Events, seats or ticket types, customers, reservations. A seat cannot be sold twice for the same event. Use a unique constraint and a transaction when you reserve.

Minimum design contents:

1. Entities and relationships (one-to-many, many-to-many).
2. Primary keys and foreign keys with `ON DELETE` actions.
3. Types: money, dates, UTC instants, booleans.
4. Unique constraints that match the business rules.
5. Audit columns if you need them.
6. Which deletes are hard and which are soft.

```text
Example shop grain
customers 1--* orders 1--* order_lines *--1 products
```

Normalize to 3NF for this exercise. Do not denormalize for speed before you have a query and a plan.

Write two or three queries that the system must answer. Design keys and indexes for those queries after the tables exist. Do not invent ten indexes on day one.

Do not model the whole world. A library does not need a full HR system. Stop at the domain that you named.

### Questions

#### Theoretical questions

1. What six items belong in the minimum design?
2. Why does a copy have only one open loan?
3. Why do you write queries before you add many indexes?
4. Why is 3NF the target for this exercise?
5. What do you omit when you "do not model the whole world"?

#### Easy practical tasks

1. Choose library, shop, or tickets. Write five sentences that bound the domain.
2. Draw the boxes and crow-foot (or 1..*) marks.
3. List primary keys and two unique business rules.
4. Write two questions the database must answer (example: "open loans for a member").

#### Medium practical tasks

1. Write `CREATE TABLE` for all tables with types, PK, FK, and unique constraints. Run them in `learn`.
2. Write the `ON DELETE` action for each foreign key. Justify each in one sentence.
3. Insert a valid sample (at least two parents and three children). Insert one row that must fail. Save the error.

#### Advanced practical tasks

1. Add a rule that needs a transaction (transfer a loan, checkout with stock, reserve a seat). Write the steps and a rollback case.
2. Review the schema against topics 2, 8, 10, and 12. Write a one-page review with pass/fail items.

---

## Write SQL by hand before you use an ORM

Write the statements in a CLI. Do not start with an ORM for this exercise. You must see the text that the DBMS runs.

Cover this set on your domain schema:

1. `INSERT` of a parent and children in one transaction.
2. `UPDATE` with a `WHERE` that uses a key.
3. `DELETE` or a soft-delete `UPDATE` that matches your policy.
4. `SELECT` with a join and a `WHERE`.
5. `GROUP BY` report (counts or sums).
6. A subquery or `EXISTS` if the report needs it.
7. `LIMIT` with a keyset page, not only `OFFSET 0`.

```sql
BEGIN;
INSERT INTO orders (order_id, customer_id, placed_at, status)
VALUES (1, 1, CURRENT_TIMESTAMP, 'placed');
INSERT INTO order_lines (order_id, product_id, quantity, unit_price)
VALUES (1, 10, 2, 19.99);
COMMIT;
```

Save the scripts in files. Run them more than once only if they are idempotent, or reset the data first.

Read the errors. A foreign-key error means your insert order is wrong or the parent is missing. A unique error means the business rule worked.

After the scripts work, you may try an ORM. Compare the SQL that the ORM emits with your scripts. If you cannot explain the ORM SQL, do not ship it.

Do not concatenate user input in these scripts. Use parameters in the application later. In the CLI, you may use literal values that you wrote.

### Questions

#### Theoretical questions

1. Why do you write SQL in a CLI before an ORM?
2. Which seven statement kinds does this section ask for?
3. What does a foreign-key error tell you?
4. When may you run the same script twice?
5. What must you compare after you introduce an ORM?

#### Easy practical tasks

1. Write and run the seven kinds of statements on your schema (one of each).
2. Save the scripts as `.sql` files in a notes folder.
3. Cause one FK error and one unique error on purpose. Write the two messages.
4. Write a keyset `SELECT` for the next page of orders or loans.

#### Medium practical tasks

1. Wrap a multi-table write in `BEGIN`/`COMMIT`. Force a `ROLLBACK` on the second insert. Show that the parent did not stay (or that you used a savepoint on purpose).
2. Write a `GROUP BY` report that would be an N+1 loop in an application. Keep it as one statement.
3. Parameterize one `SELECT` in a small program or in the CLI prepare feature if it exists.

#### Advanced practical tasks

1. Write a one-page SQL script pack: schema, seed, report, page. Another beginner must run it from zero.
2. Implement the same report in an ORM. Paste the emitted SQL. Mark extra queries. Fix N+1 if it appears.

---

## Read a query plan

A query plan is the access path that the optimizer chose. You read it before you add indexes and after you add them.

Steps:

1. Write the SQL.
2. Run `EXPLAIN` (or the vendor name) without a write side effect.
3. Find the scan or seek on each table.
4. Find the join type.
5. Compare estimated rows with actual rows if the product shows both.
6. Change one thing (an index or a rewrite). Explain again.

```text
Tokens to search for (names differ)
Seq Scan / Table Scan
Index Scan / Index Seek
Nested Loop / Hash Join / Merge Join
Sort / Aggregate
```

Use a table that is large enough to matter. A plan on 10 rows teaches little. Load thousands of rows if you need to.

Do not run `EXPLAIN ANALYZE` on a `DELETE` in production. The form that executes the statement will delete.

Check that the join key has an index on the large side. Check that the filter you expect is not hidden in a scan of the whole table.

Write the plan in your notes in four lines: access on A, access on B, join, cost or time. You do not need to memorize every node type.

If the estimate is far from actual, refresh statistics and explain again (topic 17).

### Questions

#### Theoretical questions

1. What does a query plan show?
2. What six steps does this section name?
3. Why is a plan on 10 rows a weak lesson?
4. Why must you not run an executing explain on a production `DELETE`?
5. What do you check on the join key of a large table?

#### Easy practical tasks

1. Find the explain command for your DBMS.
2. Explain a `SELECT` by primary key. Write the access token.
3. Explain a join report from your domain. Write the join token.
4. List five plan tokens that you will search for.

#### Medium practical tasks

1. Load enough rows to change a plan. Explain a selective `WHERE` before and after an index.
2. Compare estimate and actual rows on one node. Write the two numbers.
3. Rewrite a query that used a function on an indexed column (if you have one). Explain both.

#### Advanced practical tasks

1. Write a one-page annotated plan for your slowest practice query. Label each node in your words.
2. Tune one join: baseline plan and time, one index, new plan and time. Write the delta.

---

## Take a backup and restore it

A backup that you restore is the only backup that you proved. Use your practice database.

Steps:

1. Note a proof value (`COUNT(*)` of a main table, or a specific row).
2. Take a logical dump or a physical backup with the official tool.
3. Store the file off the live data directory if you can.
4. Change the live data (insert a row or update a name).
5. Restore to a new database name or a new instance. Do not overwrite the only copy until you are sure.
6. Run the proof query. Confirm that the extra change is absent if you restored the older copy.
7. Write the duration. That duration is a lower bound on RTO for this method.

```text
learn  -->  dump file  -->  learn_restored
proof: COUNT(*) = 17
```

If you only have a dump tool, use it. If you can do PITR on a disposable instance, do that as a second lab.

Do not call the replica a backup in this exercise. Take a file. Restore the file.

After the restore works, write a five-line runbook: command, file path, restore command, proof query, who to call (you).

Encrypting a class dump is optional. Do not put a dump that contains real personal data in a public folder.

### Questions

#### Theoretical questions

1. What is a proof value?
2. Why do you restore to a new name first?
3. What do you expect if you change data after the backup and you restore the backup?
4. Why is a replica not enough for this exercise?
5. What does the restore duration bound?

#### Easy practical tasks

1. Run `COUNT(*)` on your main table. Write the number.
2. Dump the database. Write the command and the file size.
3. Restore to `learn_restored` (or equivalent). Run `COUNT(*)` again.
4. Write the five-line runbook.

#### Medium practical tasks

1. Insert a row after the dump. Restore the dump. Show that the new row is not in the restored copy and still is in live `learn`.
2. Time the dump and the restore. Write RTO as "at least this long" for this size.
3. Fail the restore (wrong path). Record the error. Restore correctly.

#### Advanced practical tasks

1. If the product allows it, take a base backup and restore to a time before one write (PITR lab). Write T0, T1, and the proof.
2. Write a restore test that you could run every month. Include isolation of the copy and how you destroy it.

---

## Then specialize: PostgreSQL, MongoDB, Redis

This path was vendor-neutral. Next you pick one product and you learn its names, tools, and limits.

**PostgreSQL.** A relational DBMS. Choose it to go deeper on SQL, types, WAL, `EXPLAIN`, roles, and logical replication. It matches most of this path one-to-one.

**MongoDB.** A document DBMS. Choose it to learn collections, documents, indexes on JSON, replica sets, and when a document aggregate is the unit of work. Keep a relational system of record in mind for money if you still need one.

**Redis.** A key-value (and data-structure) store. Choose it to learn keys, TTL, persistence options, and in-memory access. Use it as a cache, session store, or lock, not as the only order ledger.

```text
After this path
  --> PostgreSQL  if you need SQL and a general record store
  --> MongoDB     if you need document aggregates
  --> Redis       if you need fast key access and TTL
```

You can learn more than one. Learn one well enough to install, back up, and explain a plan or an equivalent.

Do not start three vendor paths in the same week. Finish a first install, a first backup, and a first real schema or key design.

Official docs are the source of truth for commands. This handbook stays vendor-neutral on purpose.

When you specialize, take the ideas with you: types, transactions, indexes, backups, least privilege, measure first. Only the syntax and the internals names change.

### Questions

#### Theoretical questions

1. Which product maps most directly to this path?
2. When do you choose MongoDB after this path?
3. When do you choose Redis after this path?
4. Why must Redis not be the only order ledger?
5. What three jobs do you finish on one product before you start the next?

#### Easy practical tasks

1. Write your choice and three reasons that use this path (not marketing).
2. Open the official install page of that product. Write the version that you will use.
3. Map five terms from this path to product terms (example: WAL, dump, role).
4. Write one thing the product will not replace from this path (example: Redis will not replace 3NF invoices).

#### Medium practical tasks

1. Install the chosen product (or start a container). Run the vendor equivalent of `SELECT 1` or `PING`.
2. Recreate a small piece of your domain in that product (tables, documents, or keys). Write what you could not express.
3. Find the official backup page. Write the first command that you will learn next.

#### Advanced practical tasks

1. Write a one-page specialization plan: 10 sessions, each with one outcome (install, schema, backup, plan, roles).
2. Compare your domain in the relational design and in the new product. Write a table of features: easier, harder, missing.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do schema design, hand-written SQL, a plan, and a restore prove that you finished the common path?
2. When do you add an ORM and a second product, and what must already work?
3. Which practice from the suggested path (install, join, transaction, dump) does each section of this topic complete?
4. How do you choose PostgreSQL, MongoDB, or Redis from the access path of your domain, not from popularity?
5. A teammate wants to skip SQL and backups and "just use an ORM plus a cloud replica." Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: design list, seven SQL kinds, plan tokens, restore proof, specialization choice.
2. Confirm that `learn` has your domain tables, at least one join, and a dump file on disk.
3. Write the proof `COUNT(*)` and the dump file name in one note.
4. Write the next product name and the official doc URL (install page).

#### Medium practical tasks

1. Run a single sitting: one new constraint, one join report, one explain, one dump/restore to a new name. Write the four outcomes.
2. Give your schema and scripts to a peer (or to your future self in a README). List the steps they need. Remove one hidden manual step.
3. Write RPO/RTO for your class project from the restore time that you measured.

#### Advanced practical tasks

1. Package the domain: migrations or SQL files, seed, report, backup command, restore command, plan screenshot or text. Another machine must reproduce it.
2. Write a close-out of the common path: what you can teach, what you will learn in the vendor path, and one risk you still do not cover (HA, CDC, or sharding).
