# 12. Clients and Pooling

## Description

This topic shows how programs connect to PostgreSQL 16 and PostgreSQL 17. You learn libpq connection strings, PgBouncer pool modes, drivers (`psycopg`, `pgx`, JDBC), prepared statements, `LISTEN` / `NOTIFY`, and `COPY`.

Complete topic 11 first. You know that a replica is read-only. Complete this topic before you change memory settings.

Use one term for each concept. libpq is the official C client library. A connection string is a URI or a key-value list that libpq understands. A pool reuses server sessions. PgBouncer is a common pooler. A driver is a language binding. A prepared statement is a parsed statement that you execute with parameters. `COPY` is the bulk transfer command. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## `libpq` connection parameters

Topic 1 showed a URI. libpq also accepts keyword strings:

```text
host=localhost port=5432 dbname=shop user=shop_app sslmode=verify-full
```

URI form:

```text
postgresql://shop_app@localhost:5432/shop?sslmode=verify-full&connect_timeout=5
```

Common parameters:

| Name | Role |
| --- | --- |
| `host` | host name or path to a Unix socket |
| `port` | default `5432` |
| `dbname` | database |
| `user` | login role |
| `password` | prefer `.pgpass` or a secret store |
| `sslmode` | topic 8 |
| `connect_timeout` | seconds to wait for the TCP connect |
| `application_name` | shows in `pg_stat_activity` |
| `options` | server command-line style options |

Environment variables: `PGHOST`, `PGPORT`, `PGUSER`, `PGDATABASE`, `PGSSLMODE`, `PGAPPNAME`. A parameter in the string overrides the environment.

`application_name` helps operations:

```text
postgresql://shop_app@localhost:5432/shop?application_name=shop-api
```

```sql
SELECT application_name, state, query
FROM pg_stat_activity
WHERE datname = current_database();
```

libpq is the base for `psql`, many drivers, and `pg_dump`. If a driver documents "libpq parameters", the table above applies.

Do not put a password in process lists or CI logs. Do not use `sslmode=disable` on a network that you do not trust. Do not share one `application_name` for every microservice if you must debug connections.

Official: [https://www.postgresql.org/docs/17/libpq-connect.html](https://www.postgresql.org/docs/17/libpq-connect.html).

### Questions

#### Theoretical questions

1. Name two forms of a libpq connection string.
2. Which environment variable sets the database name?
3. What does `application_name` show?
4. What does `connect_timeout` limit?
5. Does a URI parameter override `PGHOST`?

#### Easy practical tasks

1. Connect with `psql` and a URI that sets `application_name=learn17`. Query `pg_stat_activity` for that name.
2. Connect with a keyword string (`psql "host=... dbname=..."`).
3. Write a table of five parameters and their meaning.
4. `psql --help` and find the option that passes a connection URI.

#### Medium practical tasks

1. Set `PGHOST`, `PGPORT`, `PGUSER`, and `PGDATABASE`. Run `psql` with no flags. Then override `dbname` in a URI.
2. Add `sslmode` and `connect_timeout` to a URI. Show `\conninfo`.
3. Compare the 16 and 17 libpq parameter lists. Write one parameter that you did not know.

#### Advanced practical tasks

1. Read `.pgpass` format in the libpq docs. Create a file. Connect without a password prompt. Write the file mode that the docs require on Unix.
2. Write a small program or `psql` wrapper that refuses to start if the URI contains a password. Document how you detect it.

---

## PgBouncer: transaction vs session pooling

Each PostgreSQL session uses memory and a process (or a backend slot). A web app that opens one session per HTTP request will exhaust `max_connections`. A pool keeps a smaller set of server sessions and lends them to clients.

**PgBouncer** is a popular pooler. It sits between the app and PostgreSQL. Three pool modes:

- **session** — one server connection for the whole client connection. All features work. Pooling gain is small if clients stay connected.
- **transaction** — the server connection returns to the pool at `COMMIT` or `ROLLBACK`. High reuse. Many session features break.
- **statement** — return after each statement. Rare. Breaks multi-statement transactions.

Features that need **session** mode (or a direct connection):

- `LISTEN` / `NOTIFY`
- temporary tables that survive across transactions
- `SET` that must last for the client session
- advisory locks held across transactions
- prepared statements on the server session, unless PgBouncer is configured to track them

Transaction mode is the usual choice for a stateless HTTP API that uses short transactions. Session mode is the usual choice for a worker that `LISTEN`s or holds state.

PgBouncer authenticates clients and opens server connections as configured. You still use a least-privilege role (topic 8). You still limit `max_connections` on PostgreSQL (topic 13).

Do not set `max_connections` to 1000 and skip a pool. Do not use transaction pooling and then depend on `LISTEN`. Do not run PgBouncer as a public internet proxy without TLS and auth.

Official: [https://www.pgbouncer.org/](https://www.pgbouncer.org/).

### Questions

#### Theoretical questions

1. Why does a web app need a pool?
2. When does transaction mode return a server connection to the pool?
3. Which mode do you use for `LISTEN`?
4. What happens to a temporary table in transaction mode after `COMMIT`?
5. Is PgBouncer a substitute for `GRANT`?

#### Easy practical tasks

1. Write a three-row table: mode, when the connection returns, one feature that breaks or works.
2. Open the PgBouncer config docs. Find `pool_mode`. Write the default that the page shows.
3. `SHOW max_connections;` on your server.
4. Draw: app, PgBouncer, PostgreSQL. Label two client connections and one server connection in transaction mode.

#### Medium practical tasks

1. Install PgBouncer in a lab or read a sample `pgbouncer.ini`. Write the values you would set for `pool_mode`, `listen_port`, and a `databases` line.
2. List five session features from this section. Mark each as session-only or safe in transaction mode.
3. Compare PgBouncer with a driver-side pool (one paragraph). Write who limits `max_connections`.

#### Advanced practical tasks

1. Run PgBouncer in transaction mode against a lab. Run a short `BEGIN` ... `COMMIT` workload. Watch `pg_stat_activity` count versus client count.
2. Read PgBouncer docs on prepared statements (`max_prepared_statements` in current PgBouncer). Write whether your version can track prepares in transaction mode.

---

## Drivers: `psycopg`, `pgx`, JDBC

A driver sends startup, queries, and `COPY` through the PostgreSQL protocol. Many drivers use libpq. Some speak the protocol in the language.

**`psycopg`** is the usual Python driver (psycopg 3). It supports connection strings, parameters, and `COPY`. Use bind parameters. Do not build SQL with f-strings.

**`pgx`** is a common Go driver. It can use the protocol without libpq. Connection strings still look like libpq URIs. Prefer `pgx` with a pool (`pgxpool`) in the application, plus PgBouncer when many app hosts share one cluster.

**JDBC** (`org.postgresql.Driver`) is the usual Java driver. The URL form is `jdbc:postgresql://host:5432/shop`. JDBC has its own pool products (HikariCP). The same session-versus-transaction rules apply if a second pooler sits in front.

All three must:

- use parameters (`$1` or `?` as the driver documents)
- set `application_name`
- close or return connections
- request `sslmode` that matches topic 8
- handle errors and retries without a double `INSERT` if the commit is uncertain

```python
# Python idea (psycopg 3)
# conn.execute("SELECT id FROM shop.items WHERE name = %s", ("nail",))
```

```go
// Go idea (pgx)
// conn.Query(ctx, "SELECT id FROM shop.items WHERE name = $1", "nail")
```

Do not log full connection strings. Do not disable TLS in production to "make the driver work". Do not mix two pools on the same process without a size plan (driver pool times PgBouncer times `max_connections`).

### Questions

#### Theoretical questions

1. What does a driver do?
2. Which driver does this section name for Python, Go, and Java?
3. Why must you use bind parameters?
4. What extra URL prefix does JDBC use?
5. Why can two pools in series exhaust `max_connections`?

#### Easy practical tasks

1. Open the official pages for psycopg 3, pgx, and the PostgreSQL JDBC driver. Write one connection-string example from each page.
2. Write the same `SELECT` with parameters in Python style and Go style.
3. Add `application_name` to a driver URL on paper.
4. List three errors that a driver must not hide (connect fail, unique violation, serialization failure).

#### Medium practical tasks

1. Write a 20-line program in one language that connects, `SELECT 1`, and closes. Use a parameter query for a real table if you have one.
2. Read how your driver reports `sslmode`. Set `require` or `verify-full` in a lab.
3. Compare driver pool defaults (min/max). Write a size that stays under your `max_connections`.

#### Advanced practical tasks

1. Implement a query that uses `RETURNING` (topic 3) in your driver. Print the returned id.
2. Read the driver docs on pipelining or batch (pgx or psycopg). Write one case where a batch helps and one case where it hides errors.

---

## Prepared statements

A prepared statement is a parse and analyze that the session keeps. You then `EXECUTE` with parameters. The protocol does the same work without SQL `PREPARE` when the driver sends a parse message.

```sql
PREPARE item_by_name (text) AS
    SELECT id, qty FROM shop.items WHERE name = $1;

EXECUTE item_by_name('nail');
DEALLOCATE item_by_name;
```

Benefits:

- parameters stay out of the SQL text (safer than concatenation)
- the server can skip repeat parse work
- `pg_stat_statements` sees a normalized shape

Costs:

- the plan can stay generic; a one-off plan might be better for a rare value
- the prepare lives in the session; transaction pooling can break unnamed or session-bound prepares
- `PREPARE` is per session, not cluster-wide

Drivers often prepare implicitly after a few repeats. You can turn that off when plans become bad (driver settings such as "prepare threshold").

`EXPLAIN EXECUTE item_by_name('nail');` shows the plan of the execute.

Do not concatenate user text into `PREPARE`. Do not assume a prepared statement is a cache across PgBouncer transaction mode. Do not keep thousands of named prepares in a long session without `DEALLOCATE`.

Official: [https://www.postgresql.org/docs/17/sql-prepare.html](https://www.postgresql.org/docs/17/sql-prepare.html).

### Questions

#### Theoretical questions

1. What work does `PREPARE` do before `EXECUTE`?
2. How do you remove a named prepare?
3. Why do parameters belong in `EXECUTE` and not in a string join?
4. Why can transaction pooling break prepares?
5. What does `EXPLAIN EXECUTE` show?

#### Easy practical tasks

1. `PREPARE` a `SELECT`. `EXECUTE` it twice with two values.
2. `DEALLOCATE` the name. `EXECUTE` again. Record the error.
3. `EXPLAIN EXECUTE` your prepare.
4. Write four sentences: parse, parameter, session, pool.

#### Medium practical tasks

1. Prepare an `INSERT ... RETURNING`. Execute it three times. Select the new ids.
2. Compare `pg_prepared_statements` before and after `PREPARE`.
3. In a driver, run the same parameterized query many times. Read the driver docs on automatic prepare. Write the setting name.

#### Advanced practical tasks

1. Build a case where a generic plan is worse than a custom plan (uneven values). Read "generic plan" in the 16 or 17 docs. Write what you observed.
2. Test a named `PREPARE` through PgBouncer transaction mode if you have it. Document success or failure.

---

## `LISTEN` / `NOTIFY` and `COPY`

`NOTIFY` sends a payload to all sessions that `LISTEN` on the same channel in the same database. The notify is visible after commit. A rollback does not notify.

```sql
-- session A
LISTEN shop_jobs;

-- session B
NOTIFY shop_jobs, 'order:42';
```

In `psql`, session A prints an asynchronous notice after session B commits. In a driver, you must poll or use the driver API for notifications.

Rules:

- channel names are identifiers
- the payload is `text` (keep it small; put a key, not a document)
- listeners must use a session that stays connected
- PgBouncer transaction mode is the wrong pool for listeners
- `NOTIFY` is not a job queue with retry and persistence

Use `LISTEN` for "wake up and then `SELECT` the work table". Store the work in a table. The notify is only a signal.

`COPY` moves rows between a table and a file or a stream. It is the fast path for bulk load and unload.

```sql
COPY shop.items (name, qty) FROM STDIN WITH (FORMAT csv, HEADER true);
```

In `psql`, `\copy` is a client-side variant. `\copy` reads a file on the client host. SQL `COPY` to a file path runs on the server and needs a superuser (or a granted role on a configured program). Prefer `\copy` or driver `COPY` APIs for application loads.

```text
\copy shop.items (name, qty) FROM 'items.csv' WITH (FORMAT csv, HEADER true)
\copy shop.items TO 'items_out.csv' WITH (FORMAT csv, HEADER true)
```

Binary format is faster for large numeric tables. CSV is easier to debug. Row triggers still run. Constraints still apply.

Do not use `NOTIFY` as the only store of an order. Do not `LISTEN` on a new connection per HTTP request. Do not build 100000 single-row `INSERT` statements if `COPY` is available.

Official: [https://www.postgresql.org/docs/17/sql-notify.html](https://www.postgresql.org/docs/17/sql-notify.html), [https://www.postgresql.org/docs/17/sql-copy.html](https://www.postgresql.org/docs/17/sql-copy.html).

### Questions

#### Theoretical questions

1. When does a listener see `NOTIFY`?
2. What happens to `NOTIFY` after `ROLLBACK`?
3. What is the difference between `COPY` and `\copy`?
4. Which pool mode breaks `LISTEN`?
5. Why do you still need a table for the work after `NOTIFY`?

#### Easy practical tasks

1. Open two `psql` sessions. `LISTEN` in A. `NOTIFY` in B. Write what A prints.
2. `\copy` two rows from a small CSV into a table. Select them.
3. Query `pg_listening_channels();` in the listener session.
4. Run `\h COPY` and write two `WITH` options.

#### Medium practical tasks

1. `BEGIN; NOTIFY ...; ROLLBACK;` in B. Write whether A received the payload.
2. Compare wall time of 1000 single `INSERT`s versus one `COPY` of 1000 rows (lab). Write the two times.
3. `NOTIFY` with a payload that is a row id. In A, `SELECT` that row from a table after the notice.

#### Advanced practical tasks

1. Write a small worker loop: `LISTEN`, wait, `SELECT` unprocessed rows, update them. Use two sessions. Use session pooling or a direct connection.
2. Read `COPY` permissions and notify payload limits in the 17 docs. Write who may `COPY` to a server file and a team rule for notify payload content.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `application_name`, a pool, and `pg_stat_activity` help you debug a stuck API?
2. Which features in this topic force session pooling or a direct connection?
3. How do bind parameters, `PREPARE`, and `COPY` each move data without string-joined SQL?
4. Why is `NOTIFY` a signal and `COPY` a bulk path, not a queue product?
5. A teammate opens a new PostgreSQL session per HTTP request and concatenates SQL. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. Connect with a URI that sets `application_name` and `sslmode`. Run `SELECT 1`. Show `\conninfo`.
2. Write a cheat sheet: libpq parameters, three pool modes, three drivers, `PREPARE`, `LISTEN`, `\copy`.
3. `PREPARE` a `SELECT`. `EXECUTE` it. `DEALLOCATE` it.
4. Draw app → PgBouncer (transaction) → PostgreSQL and a second path for a `LISTEN` worker.

#### Medium practical tasks

1. Write one parameterized query in `psql` (`PREPARE`) and the same query in one driver. Add `application_name`.
2. Load a CSV with `\copy`. Then `NOTIFY` a channel that a second session listens on.
3. Size a driver pool and a PgBouncer pool so that the product stays under `max_connections`. Show the arithmetic.

#### Advanced practical tasks

1. Build a small program that uses parameters, `RETURNING`, and `COPY` or a batch insert. Document pool mode and why `LISTEN` is absent.
2. Map this topic to official libpq, `PREPARE`, `NOTIFY`, and `COPY` pages. Add one driver-specific setting that this topic did not name.
