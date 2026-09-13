# 13. Application Access

## Description

This topic shows how a program talks to a DBMS. You learn drivers, connection strings, connection pooling, prepared statements, SQL injection, object-relational mappers (ORMs), query builders, raw SQL, the N+1 query problem, migrations, seeds, and fixtures.

Use one term for each concept. A driver is the library that sends protocol messages to the DBMS. A session is one connection. A migration is a versioned schema change. Complete this topic after you can write SQL and use transactions. The application must use those skills without unsafe string concatenation.

This path is vendor-neutral. Driver names and ORM names differ. The risks do not.

---

## Drivers and connection strings

A driver is a library in the application language. The driver implements the client protocol of the DBMS. Your code calls the driver. The driver sends SQL and parameters. The driver returns rows and errors.

You do not open the data files. You connect. A connection string (also called a URI or DSN) is a single text value that names how to connect. Typical fields:

- host
- port
- user
- password
- database name
- extra options (TLS, timeout, schema search path)

```text
# URI shape (product prefixes differ)
postgresql://learn:SECRET@127.0.0.1:5432/learn?sslmode=disable

# Key-value shape
host=127.0.0.1 port=5432 user=learn password=SECRET dbname=learn
```

Use the official driver for your language and product. A generic ODBC or JDBC driver is valid when the platform requires it. Read the driver documentation for types: how a decimal, a timestamp, and a boolean become language values.

A connection is a session. The session has a user, a current database, and transaction state. Close the connection when the work is done, or return it to a pool (next section).

Do not put a production password in source code. Do not commit a connection string that contains a secret. Use an environment variable or a secret store. A later security topic covers secrets in more detail.

Do not use the admin user from the application in production. Create an application user with the rights that the program needs.

If the driver fails to connect, read the error. Typical causes: wrong host, wrong port, server not running, TLS mismatch, or a rejected user.

### Questions

#### Theoretical questions

1. What job does a driver do?
2. Which fields does a connection string usually contain?
3. What state does one connection (session) hold?
4. Why must you not put a production password in source code?
5. Why must the application not use the admin user in production?

#### Easy practical tasks

1. Write a connection URI for your local DBMS with a password placeholder.
2. Name the official driver package for your language and DBMS.
3. List five connection fields in a two-column table: field and example value (no real secret).
4. Connect with a small program or a REPL. Run `SELECT 1`. Record success or the error.

#### Medium practical tasks

1. Connect with a URI and with separate fields (if the driver allows both). Write one difference in setup.
2. Change one field (wrong database name). Record the error text. Fix the field.
3. Read the driver page for type mapping. Write how `NUMERIC` and `BOOLEAN` appear in your language.

#### Advanced practical tasks

1. Write a short program that reads host, port, user, and database from environment variables and connects. Do not print the password.
2. Compare two drivers or two language bindings for the same DBMS. Write start-up, type mapping, and how you close a connection.

---

## Connection pooling

A connection pool is a set of open connections that the application reuses. Open and close of a TCP session and of a DBMS session has a cost. A pool pays that cost once, then lends a connection for one unit of work.

Typical pool settings:

| Setting | Meaning |
| --- | --- |
| Maximum size | Most connections the pool may open |
| Minimum size | Connections to keep open when idle |
| Idle timeout | When the pool closes an unused connection |
| Acquire timeout | How long a request waits for a free connection |

The application borrows a connection, runs statements, commits or rolls back, and returns the connection. "Return" means the pool owns the connection again. The application must not keep a borrowed connection across an HTTP wait or a user think time.

Pool size is not "as large as possible." Each connection uses memory on the DBMS. Too many connections slow the server. A web process with 50 workers does not need 50 idle connections if most workers wait on the network.

```text
Client requests  -->  application pool  -->  few DBMS connections
```

Do not open a new connection per SQL statement in a hot path. Do not use one global connection from many threads unless the driver documents that use. Prefer one connection per concurrent unit of work.

A pool in the application and a pool in a proxy (example: a connection pooler in front of the DBMS) can multiply connections. Count the product: application max times number of application instances.

When the pool is empty, new work waits. That wait is better than an unlimited storm of new connections. A later topic covers connection storms.

### Questions

#### Theoretical questions

1. What cost does a pool avoid?
2. What must the application do after it borrows a connection?
3. Why is a very large pool a risk for the DBMS?
4. Why must you not hold a borrowed connection during a long user wait?
5. How can two layers of pools multiply connections?

#### Easy practical tasks

1. Write four pool settings from this section and one sentence each.
2. Draw the path: request, pool, DBMS, return to pool.
3. Find the pool class or parameter name in your driver or framework.
4. Write one sentence that distinguishes "close the connection" from "return to the pool."

#### Medium practical tasks

1. Configure a small pool (maximum 2). Run three concurrent requests if you can. Record wait or error.
2. Count connections on the DBMS while your app is idle and while it works (`SELECT` from a product activity view).
3. Compute: 4 application instances times max pool 20. Write the worst-case connection count.

#### Advanced practical tasks

1. Write a one-page pool policy: max size, timeout, and who closes leaked connections. Include a leak symptom.
2. Compare an in-app pool with an external pooler at a high level. Write when you would add the second layer.

---

## Prepared statements and SQL injection

A prepared statement is SQL with placeholders. The application sends the statement text once (or the driver does). The application sends parameter values separately. The DBMS does not treat parameter values as SQL syntax.

```sql
-- Placeholder shape differs by product: $1, ?, :name
SELECT *
FROM customers
WHERE email = $1;
```

SQL injection is an attack. The attacker puts SQL syntax into a value. If the application concatenates that value into the statement string, the DBMS runs the attacker's SQL.

```text
-- Unsafe shape (do not use)
"SELECT * FROM users WHERE name = '" + userInput + "'"

-- If userInput is:  ' OR '1' = '1
-- the statement changes meaning
```

Rules:

1. Use bound parameters for every value that comes from outside the statement.
2. Do not build SQL by joining strings with user input.
3. Identifiers (table names, column names) cannot use value parameters in standard SQL. If you must choose an identifier, use a fixed allow-list, not raw input.

Prepared statements also help the DBMS reuse a plan. That benefit is secondary. Safety is the first reason.

An ORM or query builder that uses parameters is safe for values. An ORM that lets you pass a raw string is not automatically safe. Read the escape hatch.

`LIMIT` and `OFFSET` numbers are values. Bind them. Do not concatenate them.

### Questions

#### Theoretical questions

1. How does a prepared statement send values?
2. What is SQL injection?
3. Why does string concatenation of user input create injection?
4. Why can you not bind a table name as a normal value parameter?
5. What is the first reason to use prepared statements?

#### Easy practical tasks

1. Rewrite a concatenated `WHERE email = ...` query as a parameterized query in your driver.
2. Mark three inputs in a login form as "must bind" or "must allow-list."
3. Write the unsafe string and the safe parameter form side by side (in notes, not in production code).
4. Find the placeholder syntax for your DBMS and driver (`$1`, `?`, or a name).

#### Medium practical tasks

1. Build a small query function that accepts an email parameter. Pass a string that contains a quote. Confirm that you get zero or one row, not a syntax change.
2. Read the "raw SQL" page of your ORM or driver. Write one unsafe pattern that the page still allows.
3. Bind `LIMIT` in a paged query. Do not concatenate the page size.

#### Advanced practical tasks

1. Write a one-page rule for the team: parameters, allow-lists for identifiers, and a review checklist for string SQL.
2. Review a sample project (yours or a public demo) for concatenation. List each finding and the fix type (bind or allow-list). Do not run attacks on systems that you do not own.

---

## ORMs vs query builders vs raw SQL

Three common ways to send SQL from an application:

**Raw SQL.** You write the statement. The driver executes it with parameters. You control the text. You own every join and every column list.

**Query builder.** You call functions that build SQL. The builder adds parameters. You still think in tables and columns. The builder reduces string errors. It does not hide the relational model.

**ORM.** An object-relational mapper maps rows to objects (or structs). You load a `Customer`. The ORM writes `SELECT` and `INSERT`. An ORM can hide joins and load graphs of objects.

| Method | Control of SQL | Typical risk |
| --- | --- | --- |
| Raw SQL | Full | Duplication; injection if you concatenate |
| Query builder | High | Complex queries become hard to read |
| ORM | Variable | N+1 loads; hidden statements; leaky mappings |

Use raw SQL or a builder when the query is a report, a join-heavy read, or a batch write. Use an ORM when many simple row operations exist and the team accepts the mapping rules.

An ORM is not a second DBMS. Constraints, types, and transactions still live in the database. If the ORM model and the schema disagree, the schema wins at run time.

Do not treat "we use an ORM" as a security claim. The ORM must still parameterize. Do not treat "we use raw SQL" as a performance claim. A bad raw query is still bad.

Learn SQL first. Then use an ORM as a tool. If you cannot read the SQL that the ORM emits, you cannot tune it.

### Questions

#### Theoretical questions

1. What does a query builder produce?
2. What does an ORM map?
3. When is raw SQL a better fit than an ORM?
4. Why is an ORM not a second DBMS?
5. Why must you still read the SQL that an ORM emits?

#### Easy practical tasks

1. Write the same "customer by id" read as raw SQL and as an ORM or builder call (or as pseudo-code if you have no ORM).
2. Make a three-row comparison table from this section.
3. List two reports in a shop that you would write as raw SQL.
4. Find how your ORM logs SQL. Enable that log in development.

#### Medium practical tasks

1. Insert a row with raw SQL and with an ORM. Compare the generated `INSERT` text.
2. Write a two-table join in SQL. Then try the ORM equivalent. Write which form you can explain line by line.
3. Find one mapping that your ORM cannot express (a partial unique index, a view). Write how you handle it.

#### Advanced practical tasks

1. Write a one-page decision record: ORM for commands, SQL for reports (or a different split). Include risks.
2. Take one ORM query. Capture the SQL. Rewrite that SQL by hand. Compare results and a plan.

---

## N+1 query problem

The N+1 problem is a query shape. You run one query to load N parents. Then you run one query per parent to load children. Total queries: `1 + N`.

Example. You load 50 orders. Then you load the customer for each order in a loop. You run 51 queries. One join or one `WHERE order_id IN (...)` can load the same data in one or two queries.

```text
-- N+1 shape
SELECT * FROM orders;                 -- 1 query
SELECT * FROM customers WHERE id = ?; -- repeated N times

-- Set-based shape
SELECT o.*, c.name
FROM orders AS o
JOIN customers AS c ON c.customer_id = o.customer_id;
```

ORMs cause this shape when you lazy-load a relation inside a loop. The code looks like field access. Each access is a query.

Fixes:

1. Join in SQL.
2. Eager-load the relation in the ORM (one extra query with `IN`, or a join).
3. Load ids, then load children in one query, then attach them in memory.

N+1 is a correctness problem for latency and for pool use. Each extra query uses a round trip. Under load, 1 + N becomes a stall.

Do not "fix" N+1 by opening more connections. Fix the query shape. Do not prefetch the whole database to avoid N+1. Load the columns that the page needs.

Count queries in development. If a list page runs one query per row, you found N+1.

### Questions

#### Theoretical questions

1. What is the N+1 query shape?
2. Why do ORMs hide N+1 in ordinary-looking code?
3. Name two set-based fixes.
4. How does N+1 harm a connection pool?
5. Why is "open more connections" a poor fix?

#### Easy practical tasks

1. Write a loop that would cause N+1 for orders and customers. Then write the join.
2. Draw 1 + 5 queries versus one join for five orders.
3. Find the eager-load API name in your ORM (or write "use a join" if you have no ORM).
4. Count the queries on a list page in a sample app (log lines).

#### Medium practical tasks

1. Build `orders` and `customers`. Write a program that lazy-loads customers in a loop. Count queries. Then eager-load. Count again.
2. Write an `IN` query that loads all customers for a list of order ids.
3. Explain why `SELECT *` on both sides of an eager join can still be expensive.

#### Advanced practical tasks

1. Profile a real list endpoint. Record query count before and after an N+1 fix. Write the two counts and the SQL.
2. Write a team checklist: how you detect N+1 in review (logs, counters, forbidden lazy load in loops).

---

## Migrations as versioned schema change

A migration is a script (or a pair of scripts) that changes the schema. Each migration has a version number or a timestamp. The application or a tool stores which versions already ran.

You do not change production tables by hand in an undocumented session. You add a migration. You review it. You run it in each environment in the same order.

Typical operations in a migration:

- `CREATE TABLE` / `ALTER TABLE` / `DROP TABLE`
- create or drop indexes and constraints
- data backfill that the new schema requires

```text
versions
  001_create_customers.sql
  002_create_orders.sql
  003_add_orders_status.sql
```

An "up" script applies the change. A "down" script reverses it when the tool supports safe reverse. Not every change has a safe reverse. A `DROP COLUMN` that deleted data does not restore those values.

Rules:

1. One logical change per migration when you can.
2. Migrations must be repeatable in order on an empty database and on production.
3. Do not edit a migration that already ran in shared environments. Add a new migration.
4. Pair schema change with application change. Expand, migrate, contract when you need zero downtime (add column, deploy code, then drop old column later).

Seeds are not migrations. A seed loads sample data. A migration changes structure (and sometimes required data).

Do not put environment-specific host names into a migration. Do not run destructive `DROP` on production without a backup and a review.

### Questions

#### Theoretical questions

1. What does a migration version record?
2. Why must you not edit a migration that already ran in a shared environment?
3. What is the difference between an "up" script and a "down" script?
4. Why is a seed not a migration?
5. What does expand-then-contract mean for a column change?

#### Easy practical tasks

1. Name the migration tool that your language or framework uses, or a SQL-file convention.
2. Write a two-file sequence: create `customers`, then add `email`.
3. Write one reason to refuse `DROP TABLE` in a production migration without a backup.
4. List three statements that belong in a migration and two that belong in a seed.

#### Medium practical tasks

1. Run migrations on an empty database. Confirm the version table. Add a new migration. Run again.
2. Attempt to run the same migration twice. Write how the tool prevents a double apply.
3. Write an expand-contract plan for rename of `name` to `display_name` in three deployments.

#### Advanced practical tasks

1. Write a one-page migration standard: naming, review, backup, and forbidden edits. Include one unsafe example.
2. Practice a failing migration (bad SQL) in a disposable database. Record how you repair forward with a new migration.

---

## Seeds and fixtures for development

A seed loads known data into a database for development or demo. A fixture is a fixed set of rows for a test. Both are data, not schema.

Use seeds so that a new machine has customers, products, and orders to query. Use fixtures so that a test starts from a known state.

```text
Development: migrations first, then seeds
Test:       migrations (or a test schema), then fixtures
Production: migrations only, unless you have a controlled reference-data job
```

Rules for seeds:

1. Seeds must run after migrations.
2. Seeds must be safe to run more than once, or the tool must replace the data in a documented way.
3. Do not seed production with fake users unless that is a deliberate, reviewed job.
4. Do not put real personal data from production into a development seed.

Fixtures must be small. A test that loads a full production dump is slow and unclear. Name the facts that the test needs. Example: one customer, two orders, one cancelled order.

Prefer stable keys in fixtures (`customer_id = 1`) so that assertions stay readable. Reset the fixture data at the start of the test, or use a transaction that rolls back.

Do not use a seed as a substitute for a migration of required lookup rows if every environment needs those rows. Required reference data belongs in a migration or in a controlled data job. Optional demo data belongs in a seed.

### Questions

#### Theoretical questions

1. What is the difference between a seed and a fixture?
2. In which order do you run migrations and seeds?
3. Why must you not copy production personal data into a development seed?
4. Why must fixtures stay small?
5. When does reference data belong in a migration instead of a seed?

#### Easy practical tasks

1. Write a seed of three `products` rows for your practice schema.
2. Write a fixture of one customer and two orders for a "list orders" test.
3. Mark four files in a sample project as migration, seed, fixture, or application code.
4. Write two rules from this section that protect production.

#### Medium practical tasks

1. Run migrations and a seed on a disposable database. Drop the database. Repeat. Confirm that the process is repeatable.
2. Write a test that loads fixtures and asserts one join result. Reset data between runs.
3. Separate demo users from required lookup rows (example: order statuses). Put each in the correct place.

#### Advanced practical tasks

1. Design a seed that is idempotent (run twice, same end state). Implement it with `MERGE` / upsert or a clear-and-insert rule.
2. Write a one-page data policy: what developers may seed, what tests may load, and how you anonymize a production sample if the team needs one.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a driver, a pool, and parameterized SQL work together on one web request?
2. When do you choose raw SQL over an ORM if N+1 and a report query both exist on the same page?
3. How do migrations, seeds, and fixtures differ in purpose, environment, and risk?
4. Which application-access failures look like "the DBMS is slow" but start in the client?
5. A teammate concatenates a table name and opens a new connection per row. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: URI fields, pool return rule, bind parameters, ORM versus SQL, N+1 shape, migration order, seed versus fixture.
2. Connect with a driver. Run one parameterized `SELECT`. Return the connection to a pool or close it.
3. Add a migration that creates a table. Add a seed of two rows. Query the rows.
4. Count queries for a two-table list in logs. Write the count.

#### Medium practical tasks

1. Build a small program: pool of two connections, parameterized join, no N+1, schema from migrations, data from a seed.
2. Break the program with an injected quote in a concatenated query (disposable database). Then fix with parameters. Record both outcomes.
3. Write an expand-contract migration plan plus a fixture that tests the new column.

#### Advanced practical tasks

1. Review a small app (yours) against this topic: secrets in URIs, pool size, injection, N+1, migration history, seed safety. Write findings and fixes.
2. Add an ORM and a raw SQL report side by side. Capture SQL, query count, and a one-page recommendation for that app.
