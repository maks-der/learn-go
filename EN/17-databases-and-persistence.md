# 17. Databases and Persistence

## Description

This topic shows how Go programs store data. The standard package `database/sql` is an interface. A driver implements the low-level calls. You write SQL, you manage transactions, and you pass `context.Context` into queries.

Use one term for each concept. A driver talks to one database product. A pool (`sql.DB`) holds connections. A transaction groups statements so that they commit or roll back together. An ORM generates or maps SQL for you. Do not store secrets in source files. Do not concatenate user input into SQL. Use parameters.

---

## `database/sql` interface

Package `database/sql` does not contain a PostgreSQL or MySQL client. It defines types that drivers implement. Your application imports the driver for its side effect and then uses `sql.DB`.

```go
import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
)

db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
if err != nil {
    return err
}
defer db.Close()

if err := db.PingContext(ctx); err != nil {
    return err
}
```

The blank import registers the driver. `sql.Open` validates the driver name and the configuration. It does not always open a real network connection. `Ping` or `PingContext` proves that a connection works.

`sql.DB` is safe for concurrent use. Do not open a new `sql.DB` per request. Open one pool at process start. Close it at process shutdown.

Common methods:

- `Exec` / `ExecContext` — INSERT, UPDATE, DELETE, DDL
- `Query` / `QueryContext` — many rows
- `QueryRow` / `QueryRowContext` — one row
- `Prepare` / `PrepareContext` — a reusable statement
- `Begin` / `BeginTx` — a transaction

Scan rows into variables:

```go
rows, err := db.QueryContext(ctx, `SELECT id, name FROM items WHERE qty > $1`, min)
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    if err := rows.Scan(&id, &name); err != nil {
        return err
    }
}
if err := rows.Err(); err != nil {
    return err
}
```

Always `Close` `Rows`. Always check `rows.Err()` after the loop. `QueryRow` returns `sql.ErrNoRows` when the row is missing. Test that error with `errors.Is`.

NULL columns need nullable types: `sql.NullString`, `sql.NullInt64`, `sql.NullTime`, or a pointer (`*string`). A normal `string` cannot hold NULL.

Placeholders depend on the driver. PostgreSQL uses `$1`, `$2`. MySQL and SQLite use `?`. Do not switch placeholder style without a look at the driver docs.

`Result` from `Exec` can report `RowsAffected` and `LastInsertId`. `LastInsertId` is not portable. PostgreSQL prefers `INSERT ... RETURNING id` and `QueryRow`.

Do not ignore errors from `Scan`, `Close`, or `Commit`. Do not keep `Rows` open while you run another query on the same connection if you hold the only connection. Close or use a transaction.

### Questions

#### Theoretical questions

1. What does a blank import of a driver package do?
2. Does `sql.Open` always connect to the database? How do you prove the connection?
3. Why do you open one `sql.DB` for the process instead of one per request?
4. What error does `QueryRow` return when no row matches?
5. Why does a NULL column need `sql.NullString` or a pointer?
6. Why is `LastInsertId` a weak portable API?

#### Easy practical tasks

1. Write a program that calls `sql.Open` with a driver name and a DSN from the environment. Call `PingContext`. Print the error if the database is down.
2. Execute `SELECT 1` (or the equivalent for your driver). Scan the value. Print it.
3. Write a `QueryRow` that can return `sql.ErrNoRows`. Handle that error without a panic.
4. List the placeholder style for PostgreSQL and the placeholder style for SQLite.

#### Medium practical tasks

1. Create a table `items(id, name, qty)`. Insert one row with `ExecContext`. Select it with `QueryRowContext` and `Scan`.
2. Query many rows. Forget `rows.Close` in one version. Then add `defer rows.Close()` and `rows.Err()`. Write why both matter.
3. Scan a NULL `name` into `sql.NullString`. Print `Valid` and `String`.

#### Advanced practical tasks

1. Write a small repository type `ItemStore` with `Get`, `List`, and `Create` that use `*sql.DB`. Add tests with a real SQLite file or a test container. Do not mock `sql.DB` unless you already know a pattern.
2. Compare `INSERT ... RETURNING` on PostgreSQL with `LastInsertId` on SQLite. Write six sentences about portability.

---

## Drivers (for example `pgx`, `lib/pq`, SQLite)

A driver translates `database/sql` calls into the protocol of one database. You choose the driver when you choose the database.

**PostgreSQL**

- `github.com/jackc/pgx/v5` is the current client. Use `pgx/v5/stdlib` when you want `database/sql`. Use `pgxpool` and the native `pgx` API when you want PostgreSQL features and fewer layers.
- `github.com/lib/pq` is an older `database/sql` driver. It is in maintenance mode. Do not start a new project on `lib/pq`.

Native `pgx` gives COPY, LISTEN/NOTIFY, and better control of types. `database/sql` gives one API for more than one database. Many teams use `pgx` through `stdlib` so that `sqlx` or migration tools still work.

**SQLite**

- `github.com/mattn/go-sqlite3` uses cgo and the C SQLite library. It is fast and complete. It needs a C compiler.
- `modernc.org/sqlite` is a pure Go port. It avoids cgo. Build and cross-compile are easier. Measure if you need the last bit of speed.

SQLite is one file. It is a good choice for tests, local tools, and small apps. It is a weak choice for many writers on many hosts. Use a server database when many processes must write at once.

**MySQL and others**

Drivers exist for MySQL, SQL Server, and others. The `database/sql` API stays the same. Types, placeholders, and features still differ.

Register names:

- pgx stdlib: `sql.Open("pgx", dsn)`
- lib/pq: `sql.Open("postgres", dsn)`
- mattn sqlite: `sql.Open("sqlite3", path)`

Read the driver README for the DSN format. PostgreSQL URLs look like `postgres://user:pass@host:5432/dbname?sslmode=require`. SQLite DSNs look like a file path plus options.

Keep the driver version in `go.mod`. Run `go test` after a driver major upgrade. Major upgrades can change default types (time zones, bytea, UUID).

Do not import two drivers only to "have options" in production. One product, one driver, one DSN format.

### Questions

#### Theoretical questions

1. What is the role of a database driver in a Go program?
2. Why do new PostgreSQL projects choose `pgx` instead of `lib/pq`?
3. What is the difference between `pgx` native API and `pgx/v5/stdlib`?
4. Why does `mattn/go-sqlite3` need cgo?
5. When is SQLite a poor fit?
6. Where do you find the DSN format for a driver?

#### Easy practical tasks

1. Open the README of `pgx` and of `modernc.org/sqlite`. Write the `sql.Open` driver name for each.
2. Add one driver to a module with `go get`. Show the `require` line in `go.mod`.
3. Write a DSN for local SQLite (`file:test.db` or a path). Open and ping.
4. List three PostgreSQL features that native `pgx` documents and `database/sql` does not expose directly.

#### Medium practical tasks

1. Connect to PostgreSQL with `pgx/v5/stdlib` if you have a server, or skip to SQLite and write why you skipped. Run `SELECT version()` or `sqlite_version()`.
2. Build the same small program with `mattn/go-sqlite3` and with `modernc.org/sqlite` (two modules or build tags). Record whether cgo was required.
3. Read the `lib/pq` repository status. Quote one sentence about maintenance in your own words.

#### Advanced practical tasks

1. Use native `pgxpool` for a query and `database/sql` with `pgx` stdlib for the same query. Compare the code and one benchmark if you have PostgreSQL.
2. Write a short driver-choice report for a team: PostgreSQL service, CI tests, and a CLI tool. Recommend `pgx`, SQLite variant, and whether to stay on `database/sql`.

---

## Prepared statements and transactions

A prepared statement parses SQL once on the server (or in the driver) and then runs with new parameters. Parameters stay out of the SQL text. That is the correct way to pass user values.

```go
stmt, err := db.PrepareContext(ctx, `SELECT name FROM items WHERE id = $1`)
if err != nil {
    return err
}
defer stmt.Close()

var name string
err = stmt.QueryRowContext(ctx, id).Scan(&name)
```

`database/sql` may also prepare statements inside `Query` depending on the driver and the database. An explicit `Prepare` helps when you run the same SQL many times in a loop. Close the statement.

A transaction groups statements. All statements commit, or all statements roll back.

```go
tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
if err != nil {
    return err
}
defer tx.Rollback() // no-op after a successful Commit

if _, err := tx.ExecContext(ctx, `UPDATE accounts SET bal = bal - $1 WHERE id = $2`, amount, fromID); err != nil {
    return err
}
if _, err := tx.ExecContext(ctx, `UPDATE accounts SET bal = bal + $1 WHERE id = $2`, amount, toID); err != nil {
    return err
}
return tx.Commit()
```

`defer tx.Rollback()` is a safe pattern. After `Commit`, `Rollback` returns `sql.ErrTxDone`. You can ignore that error in the defer.

Rules:

- use the `tx` handle inside the transaction, not `db`
- a `Stmt` from `db.Prepare` is not bound to your transaction unless you use `tx.Stmt` or `tx.Prepare`
- check every error before `Commit`
- pick an isolation level on purpose; the default depends on the database
- keep transactions short; do not wait for a user click inside an open transaction

`sql.LevelReadCommitted` and `sql.LevelSerializable` are common choices. Serializable reduces some races and can increase retries. Read the database documentation. Go only passes the isolation request to the driver.

Savepoints exist in some databases. `database/sql` does not provide a portable savepoint API. Use extra SQL if you need them and you accept the vendor lock.

Do not use string format to insert values. A quote in a name then becomes SQL injection.

### Questions

#### Theoretical questions

1. What problem do parameters in a prepared statement solve?
2. Why does `defer tx.Rollback()` stay correct after a successful `Commit`?
3. Why must you call `Exec` on `tx` and not on `db` during a transaction?
4. What does isolation `sql.LevelSerializable` ask the database to do?
5. Why must a transaction stay short?
6. How do you attach a prepared statement to a transaction?

#### Easy practical tasks

1. Prepare `SELECT 1`. Query the statement. Close it.
2. Begin a transaction. Insert a row. Rollback. Show that the row is missing.
3. Begin a transaction. Insert a row. Commit. Show that the row exists.
4. Write a bad string-concatenated query in a comment. Write the parameterized form next to it.

#### Medium practical tasks

1. Transfer a quantity between two rows in one transaction. Add a test that rolls back when the second update fails.
2. Run the same `SELECT` 100 times with `Prepare` and 100 times with `Query`. Compare times on your database. Write one sentence about when prepare helps.
3. Set `TxOptions{ReadOnly: true}` for a read-only transaction. Try an insert. Record the error.

#### Advanced practical tasks

1. Implement a money transfer that retries on a serialization failure. Detect the driver error. Limit the retry count. Use context cancel to stop.
2. Show a bug: query with `db` in the middle of a `tx` that has uncommitted writes. The read may not see the write. Fix the read to use `tx`.

---

## Context-aware queries

Every long database call must accept `context.Context`. The context carries a deadline and a cancel signal. When the client hangs up, you cancel the query. The driver can then stop work on the server if the protocol allows it.

Use the `*Context` methods:

- `PingContext`
- `ExecContext`
- `QueryContext`
- `QueryRowContext`
- `PrepareContext`
- `BeginTx` (the first argument is a context)

```go
ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
defer cancel()

row := db.QueryRowContext(ctx, `SELECT name FROM items WHERE id = $1`, id)
```

In an HTTP handler, start from `r.Context()`. That context cancels when the client goes away. Add a timeout if the request context has no deadline.

`QueryRowContext` still returns `*sql.Row`. `Scan` is the call that waits. If the context is done, `Scan` returns the context error (or a driver wrapper). Use `errors.Is(err, context.DeadlineExceeded)` or `errors.Is(err, context.Canceled)`.

`BeginTx` uses the context for the begin request and can bind cancel to the transaction, depending on the driver and Go version. If the context is cancelled, the transaction must not continue. Do not reuse a `tx` after a cancel.

Do not store a request context in a background worker. The worker outlives the request. Derive a new context with a timeout for jobs.

Do not pass `context.Background()` in a handler when you have `r.Context()`. Background never cancels from the client.

`sql.Conn` from `db.Conn(ctx)` is a reserved connection. Close it. Use it when you need session state (`SET LOCAL`, temporary tables). Prefer normal pool use when you do not need session state.

### Questions

#### Theoretical questions

1. Which methods on `sql.DB` take a context?
2. Why does an HTTP handler start from `r.Context()`?
3. When does `QueryRowContext` wait for the network?
4. Why must a background worker not keep the request context?
5. What do you do with `sql.Conn` when the function returns?
6. What error values do you check after a timeout?

#### Easy practical tasks

1. Call `PingContext` with a context that is already cancelled. Print the error.
2. Use `context.WithTimeout` of 1 nanosecond around a real query. Record the error.
3. In a sample handler, write `QueryRowContext(r.Context(), ...)` (the handler can be a short program that compiles).
4. Write four sentences: request context versus `context.Background()`.

#### Medium practical tasks

1. Start a slow query (`pg_sleep` or a busy SQLite loop). Cancel the context from another goroutine. Record whether the query stops.
2. Pass a deadline that is shorter than the handler timeout. Show which limit wins.
3. Use `db.Conn(ctx)` for two statements that must share a session. Close the conn.

#### Advanced practical tasks

1. Wire an HTTP handler to a store method `Get(ctx, id)`. Cancel the request in an `httptest` client. Show the store error.
2. Document your driver behavior on cancel (does the server query stop?). Use driver docs or a measured test. Write eight sentences.

---

## Connection pooling

`sql.DB` is a pool. When you `Query`, the pool takes a connection, runs the work, and returns the connection. You do not manage sockets in normal code.

Tune the pool with:

```go
db.SetMaxOpenConns(10)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(30 * time.Minute)
db.SetConnMaxIdleTime(5 * time.Minute)
```

Meanings:

- `MaxOpenConns` — maximum connections in use at once (0 means no limit; do not use 0 in production)
- `MaxIdleConns` — connections kept open while idle
- `ConnMaxLifetime` — maximum age; useful when a proxy closes old connections
- `ConnMaxIdleTime` — maximum idle time before the pool closes a connection

If `MaxOpenConns` is too low, goroutines wait for a free connection. If it is too high, the database rejects extra sessions or slows down. Match the limit to the database `max_connections` and to the number of app replicas. Five replicas with 20 open connections each is 100 sessions.

`Ping` at startup fills at least one connection. A burst of traffic then grows the pool up to `MaxOpenConns`.

Health checks: `PingContext` is enough for a simple live check. Do not run a heavy query on every probe. Some platforms probe too often and empty the idle pool.

`Stats()` returns `sql.DBStats`: open connections, in use, idle, wait count. Log those values when you debug pool exhaustion.

Native `pgxpool` is a different pool. Do not wrap `pgxpool` in `database/sql` and also open a second `sql.DB` to the same app without a plan. One pool per process is the simple rule.

Do not copy a `sql.DB` by value in a way that you think duplicates the pool. The `sql.DB` pointer is the handle. Pass `*sql.DB`.

Close the pool at shutdown so that the process drops sockets cleanly. In-flight queries still need a shutdown policy (stop new work, wait, then close).

### Questions

#### Theoretical questions

1. What does `SetMaxOpenConns` limit?
2. Why is `MaxOpenConns = 0` a problem in production?
3. How do you size the pool when many app replicas share one database?
4. What problem does `SetConnMaxLifetime` solve with a proxy?
5. What does `sql.DB.Stats` tell you?
6. Why do you pass `*sql.DB` and not open a pool inside every function?

#### Easy practical tasks

1. Set `MaxOpenConns` to 2 and `MaxIdleConns` to 1. Print `db.Stats()` after ping.
2. Read `go doc sql.DB.SetConnMaxIdleTime`. Write the purpose in one sentence.
3. Open one pool in `main` and pass it to two functions. Do not call `sql.Open` in those functions.
4. List the four pool setters and a typical value for a small service.

#### Medium practical tasks

1. Set `MaxOpenConns(1)`. Run two slow queries in two goroutines. Measure wait time. Then raise the limit to 2 and compare.
2. Log `WaitCount` and `OpenConnections` after a burst of queries.
3. Set a short `ConnMaxIdleTime`. Wait. Show that idle connections drop (via `Stats`).

#### Advanced practical tasks

1. Calculate a pool size for 3 app replicas and a database `max_connections` of 100. Reserve connections for admin and migrations. Write the numbers.
2. Simulate pool exhaustion in a test (low `MaxOpenConns`, many goroutines). Add a timeout context. Record errors. Propose one code change and one config change.

---

## SQL vs ORMs (`sqlx`, `ent`, `gorm`) — trade-offs

You write SQL yourself with `database/sql`. You can add a thin helper. You can use a full ORM. Each step hides more SQL and adds more hidden behavior.

**`database/sql`**

You write SQL. You scan columns. You control every query. The cost is boilerplate and careful `Scan` order.

**`sqlx` (`github.com/jmoiron/sqlx`)**

`sqlx` extends `database/sql`. `Get`, `Select`, and `StructScan` fill structs. Named queries use `:name`. You still write SQL. `sqlx` is a thin layer. Teams that like SQL often stop here.

**`ent`**

`ent` is a code-generation framework. You define a schema in Go. You generate a typed client. Queries are Go code. The tool can emit migrations. The fit is strong when the graph of entities is the product. The cost is a generate step and a learning curve.

**`gorm`**

`gorm` is a full ORM. Models are structs. Associations, hooks, and callbacks are available. You can write less SQL at the start. The cost is hidden queries, N+1 loads, and surprise `UPDATE` shapes. You must still read the SQL that GORM emits.

Trade-offs:

| Approach | Control of SQL | Boilerplate | Surprise queries | Good fit |
| --- | --- | --- | --- | --- |
| `database/sql` | high | high | low | simple stores, hot paths |
| `sqlx` | high | medium | low | most services that use SQL |
| `ent` | medium | low after generate | medium | typed domain, many relations |
| `gorm` | low unless you force SQL | low | high if you are careless | prototypes, CRUD with care |

Rules:

- learn `database/sql` first
- do not pick GORM only to avoid SQL; you still debug SQL in production
- use `sqlx` when struct scan is the only pain
- use `ent` when you accept generation and want compile-time query help
- measure N+1: a loop of `Get` by id is a bug in every layer
- keep raw SQL for reports and updates that the ORM expresses poorly

ORMs do not remove transactions, contexts, or pools. You still set those. You still write migrations (or review generated ones).

### Questions

#### Theoretical questions

1. What problem does `sqlx` solve that `database/sql` does not solve?
2. What extra step does `ent` require that GORM does not always require?
3. What is an N+1 query problem?
4. Why is "we use an ORM so we do not need SQL" a false statement?
5. When is raw `database/sql` the better choice on a path that runs often?
6. Which layer still owns transactions and context when you use GORM or `ent`?

#### Easy practical tasks

1. Map one `items` row to a struct with `database/sql` `Scan` and, if you add `sqlx`, with `Get`. Compare the code length.
2. Read the `sqlx` README. Write the purpose of `Select` in one sentence.
3. Read the `ent` introduction. Write what you generate and from what source.
4. List two GORM features that can hide a query. Give one risk for each.

#### Medium practical tasks

1. Implement `ListItems` with `database/sql` and with `sqlx.Select`. Keep the same SQL. Compare tests.
2. Show an N+1 pattern: list parents, then query children in a loop. Rewrite with one `WHERE parent_id IN (...)` (or a join).
3. Enable SQL logging in GORM or `ent` (or print `sqlx` SQL). Capture one insert. Confirm the columns.

#### Advanced practical tasks

1. Write the same three-table use case (users, orders, items) once with `sqlx` and once with `ent` or GORM. Compare line count, test setup, and one query plan or log.
2. Produce a decision note for your team: default stack (`database/sql` or `sqlx`), when to allow `ent`, when to allow GORM, and a ban on unbounded `Preload` or equivalent.

---

## Migrations

A migration is a versioned change to the schema. Each migration has an identifier and an up script. Many tools also have a down script. The tool records applied versions in a table.

You use a tool such as:

- `github.com/golang-migrate/migrate`
- `github.com/pressly/goose`
- Atlas (`ariga.io/atlas`)

The workflow is the same:

1. Write a new migration file (SQL or generated SQL).
2. Apply it to development.
3. Run tests.
4. Apply the same files to production in order.
5. Do not edit a migration that production already applied. Add a new migration.

Example names:

```text
0001_init.up.sql
0001_init.down.sql
0002_items_qty.up.sql
0002_items_qty.down.sql
```

Up scripts must be safe to run once. Down scripts are for development and emergency rollback. Some teams never run down in production. They add a forward migration that fixes the problem. That policy is valid if you write it down.

Expand and contract is the safe pattern for live systems:

- add a new column (nullable or with a default)
- deploy code that writes both old and new
- backfill
- deploy code that reads the new column
- drop the old column in a later migration

Do not drop a column in the same release that still reads it.

Run migrations from CI or from a controlled job. Do not let every replica run migrations at startup without a lock. Two replicas can race. `migrate` and `goose` can use a database lock. Read the tool docs.

Keep migrations in the repository. Review them like code. Test against a throwaway database. SQLite in tests is not enough if you use PostgreSQL features. Use the same product in integration tests.

Do not mix hand-edited production schema with tool-managed versions. The version table then lies.

### Questions

#### Theoretical questions

1. What does a migration tool store in the database besides your tables?
2. Why must you not edit a migration that production already applied?
3. What is the expand-and-contract pattern?
4. Why is a down migration risky in production?
5. What goes wrong when every replica runs migrations at startup without a lock?
6. Why can SQLite tests miss a PostgreSQL migration bug?

#### Easy practical tasks

1. Install or read the docs of `goose` or `golang-migrate`. Write the command that applies migrations.
2. Write an up script that creates `items(id INTEGER PRIMARY KEY, name TEXT NOT NULL)`.
3. Write a down script that drops `items`.
4. List three file names in order for init, add column, and drop column.

#### Medium practical tasks

1. Apply two migrations to a local SQLite or PostgreSQL database. Show the version table contents.
2. Change an applied migration on purpose. Apply again. Record the tool error or the silent skip. Restore the file. Add a new migration instead.
3. Write an expand step: add `name_new` without dropping `name`.

#### Advanced practical tasks

1. Automate migrations in a small script: up in CI, version print, fail if dirty. Document the dirty-state recovery from the tool docs (do not guess).
2. Plan a live rename of `name` to `title` in five deployments. Write the migration list and the code change for each step. Do not run it against a shared production database.

---

## Redis / caches (optional)

A cache stores data that is expensive to load or to compute. Redis is a remote in-memory store. Common Go clients include `github.com/redis/go-redis/v9`. The standard library does not include a Redis client.

Use a cache when:

- many reads hit the same row or the same JSON
- the data can be slightly old (you set a TTL)
- the database is the source of truth

Do not use Redis as the only copy of important data unless you design for that (and you accept durability limits). A restart or a flush must not lose money records.

Cache-aside pattern:

1. Read the cache.
2. On miss, read the database.
3. Write the cache with a TTL.
4. Return the value.

```go
func (s *Store) GetItem(ctx context.Context, id string) (Item, error) {
    key := "item:" + id
    if b, err := s.rdb.Get(ctx, key).Bytes(); err == nil {
        var it Item
        if err := json.Unmarshal(b, &it); err == nil {
            return it, nil
        }
    }
    it, err := s.dbGet(ctx, id)
    if err != nil {
        return Item{}, err
    }
    if b, err := json.Marshal(it); err == nil {
        _ = s.rdb.Set(ctx, key, b, 5*time.Minute).Err()
    }
    return it, nil
}
```

On write, delete or update the key. A stale cache is a common bug. A short TTL limits the damage.

Other Redis uses: rate limits, locks, session blobs, pub/sub. Each use has failure modes. A lock must expire. A rate limit must not block health checks. Pub/sub is not a durable queue.

`sync.Map` or a guarded `map` is an in-process cache. It does not share across replicas. Redis shares across replicas. Choose the smaller tool that works.

Always pass `context.Context` into Redis calls. Set timeouts. Redis must not hang a handler without a deadline.

This section is optional. Complete SQL first. Add a cache when a measurement shows a hot key.

### Questions

#### Theoretical questions

1. What is the source of truth in a cache-aside design?
2. What does a TTL prevent?
3. Why do you invalidate a key after a write?
4. When is an in-process map enough instead of Redis?
5. Why must a distributed lock expire?
6. Why is this section optional in the learning path?

#### Easy practical tasks

1. Read the `go-redis` getting-started page. Write the `NewClient` options that you need (address, password).
2. Write the four steps of cache-aside in a comment above an empty `GetItem`.
3. Choose a key format `item:{id}`. Write two example keys.
4. List two data types that must not live only in Redis for your project.

#### Medium practical tasks

1. Implement cache-aside for one `Get` with Redis or with an in-memory map if Redis is not available. State which one you used.
2. Update a row and delete the cache key in the same function. Write a test that reads a fresh value.
3. Set a 1-second TTL. Sleep 2 seconds. Show a miss and a new load.

#### Advanced practical tasks

1. Add a singleflight (or `golang.org/x/sync/singleflight`) so that many concurrent misses load the database once. Test with a slow `dbGet`.
2. Write a failure-mode note: Redis down, Redis slow, stale read after write, many goroutines retry at the same time. Give one mitigation for each.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of an HTTP request that loads a row: context, pool, driver, SQL parameters, scan, optional cache. Name each type.
2. When do you use a transaction, when do you use a single `ExecContext`, and when do you use a read-only query?
3. How do `database/sql`, `sqlx`, `ent`, and GORM differ in control of SQL and in surprise queries?
4. What rules keep migrations and connection pools safe when many replicas run?
5. A teammate concatenates SQL, opens a pool per request, skips context, and caches money balances in Redis with no TTL. Which facts from this topic do you use in the review?

#### Easy practical tasks

1. Create a module `example.com/storex`. Open SQLite, apply a create-table statement, insert one row, select it, close the pool. Run `go test`.
2. Write a one-page cheat sheet: `sql.Open`, `PingContext`, placeholders, `Rows.Close`, `BeginTx`, pool setters, `ErrNoRows`, migration order, cache-aside.
3. Export a sample DSN with a password replaced by `$DATABASE_URL`. Show `os.Getenv` in `main`.
4. Draw a diagram: `Handler` → `Store` → `*sql.DB` → driver → database. Add an optional Redis box.

#### Medium practical tasks

1. Build `ItemStore` with `Create` and `Get` on SQLite. Add a transaction that inserts two items and fails the second insert on purpose. Prove rollback.
2. Add `context.WithTimeout` to every store method. Write a test that uses a cancelled context.
3. Write one up migration and apply it with a tool or with a documented `Exec` of the file contents in a test helper. Record the version.

#### Advanced practical tasks

1. Build a small JSON API on top of `ItemStore` with Go 1.22 `ServeMux`. Use context from the request. Add a table test for `ErrNoRows` → HTTP 404.
2. Add an optional Redis (or memory) cache to `Get`. Measure a tight loop of gets. Write when the cache is not worth the invalidation code.
