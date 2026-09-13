# 14. Databases

## Description

This topic shows how Go programs store data. The standard package `database/sql` is an interface. A driver implements the low-level calls. You write SQL, you manage transactions, and you pass `context.Context` into queries.

Use one term for each concept. A driver talks to one database product. A pool (`sql.DB`) holds connections. A transaction groups statements so that they commit or roll back together. An ORM generates or maps SQL for you. Do not store secrets in source files. Do not concatenate user input into SQL. Use parameters.

---

## `database/sql`

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

Always `Close` `rows`. Always check `rows.Err()` after the loop. `QueryRow` returns `sql.ErrNoRows` when no row exists. Detect it with `errors.Is`.

Use `sql.NullString` and the other `Null*` types when a column can be NULL. A plain `string` scan fails on NULL.

### Questions

#### Theoretical questions

1. What does a blank import of a driver package do?
2. Does `sql.Open` always open a network connection?
3. Why must you not open a new `sql.DB` per request?
4. What does `sql.ErrNoRows` mean?
5. Why do you call `rows.Err()` after `Next` returns false?

#### Easy practical tasks

1. Read `go doc database/sql`. Write the purpose of `DB`, `Rows`, and `Tx` in one sentence each.
2. Write a program that calls `sql.Open` with a driver name that is not registered. Record the error.
3. Write the `Scan` loop from this section as a function signature that returns `[]Item` and `error`.
4. List four `database/sql` methods and whether they return rows.

#### Medium practical tasks

1. Use a driver that supports SQLite in-memory (for example `modernc.org/sqlite`). Open, `Ping`, create a table, insert one row, query it.
2. Query a missing id with `QueryRowContext`. Detect `sql.ErrNoRows`.
3. Scan a NULL name into `string` and record the error. Then use `sql.NullString`.

#### Advanced practical tasks

1. Map `rows` into a slice and guarantee `rows.Close` on every return path. Use `defer`.
2. Read the `database/sql` tutorial on go.dev. Write six steps from open to scan in your own words.

---

## Drivers, prepared statements, and transactions

A driver is a separate module. Examples: `pgx` for PostgreSQL, `modernc.org/sqlite` or `github.com/mattn/go-sqlite3` for SQLite, a MySQL driver for MySQL. The first argument of `sql.Open` is the driver name that the driver registers.

Placeholders depend on the database. PostgreSQL uses `$1`, `$2`. SQLite and MySQL often use `?`. Do not copy a `$1` query to a `?` driver without a change.

A prepared statement is SQL that the database parses one time. You then execute it with parameters.

```go
stmt, err := db.PrepareContext(ctx, `SELECT name FROM items WHERE id = $1`)
if err != nil {
	return err
}
defer stmt.Close()
row := stmt.QueryRowContext(ctx, id)
```

`QueryContext` with parameters also prepares or uses a statement inside the pool. An explicit `Prepare` helps when you run the same SQL many times in one function. Always `Close` the statement.

A transaction groups statements.

```go
tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
if err != nil {
	return err
}
defer tx.Rollback()

if _, err := tx.ExecContext(ctx, `UPDATE accounts SET n = n - 1 WHERE id = $1`, from); err != nil {
	return err
}
if _, err := tx.ExecContext(ctx, `UPDATE accounts SET n = n + 1 WHERE id = $1`, to); err != nil {
	return err
}
return tx.Commit()
```

`defer tx.Rollback()` is safe. Rollback after Commit is a no-op. If you return before Commit, Rollback undoes the work.

Do not concatenate user input into SQL. That is SQL injection. Always pass values as parameters.

```go
// wrong
q := "SELECT * FROM items WHERE name = '" + name + "'"
// correct
db.QueryContext(ctx, `SELECT * FROM items WHERE name = $1`, name)
```

`sql.Tx` is not safe for concurrent use. Use one goroutine per transaction.

### Questions

#### Theoretical questions

1. Why do placeholder marks differ between PostgreSQL and SQLite?
2. What does `PrepareContext` return?
3. Why is `defer tx.Rollback()` correct even when you `Commit`?
4. Why must you not build SQL with `+` and user text?
5. Is `sql.Tx` safe for many goroutines?

#### Easy practical tasks

1. Write one SQL string with `?` and one with `$1`. Label the driver that each form fits.
2. Write a `BeginTx` / `Exec` / `Commit` sequence in comments as six steps.
3. Rewrite a concatenated `WHERE` clause to a parameter. Show both strings.
4. Read `go doc sql.TxOptions`. Write the meaning of `ReadOnly`.

#### Medium practical tasks

1. In SQLite, create two tables. Insert in a transaction. Fail the second insert. Show that the first insert is gone after Rollback.
2. Prepare `INSERT` and run it three times with different parameters. Query the count.
3. Start two updates in one transaction that transfer a count between two rows. Print both rows after Commit.

#### Advanced practical tasks

1. Handle a serialization failure: retry the transaction a fixed number of times when the driver error matches. Document the error you treat as retryable.
2. Compare `db.QueryContext` with parameters and an explicit `Prepare` in a loop of 1000 queries. Write when you keep `Prepare`.

---

## Connection pooling and context

`sql.DB` is a pool. The pool opens connections as needed. It reuses idle connections.

Set pool limits:

```go
db.SetMaxOpenConns(10)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(time.Hour)
db.SetConnMaxIdleTime(10 * time.Minute)
```

`MaxOpenConns` bounds open connections to the database. A value that is too high can overload the database. A value that is too low makes requests wait. Start with a small number. Measure.

`MaxIdleConns` keeps unused connections ready. `ConnMaxLifetime` closes old connections so that load balancers and idle timeouts on the server side do not surprise you.

Every query method has a `Context` form. Prefer `QueryContext`, `ExecContext`, `QueryRowContext`, `BeginTx`, and `PingContext`. When the context is done, the driver cancels the query if the database supports cancel.

```go
ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
defer cancel()
row := db.QueryRowContext(ctx, `SELECT name FROM items WHERE id = $1`, id)
```

Pass the request context in a server. Add a timeout if the request context has no deadline. Do not use `context.Background()` in a handler.

`db.Stats()` reports pool counters. Use them when you debug wait time.

Close rows and statements so that the connection returns to the pool. A leaked `Rows` holds a connection. That leak can exhaust `MaxOpenConns`.

Do not store `*sql.Tx` on a long-lived struct across requests. A transaction belongs to one unit of work.

### Questions

#### Theoretical questions

1. What does `SetMaxOpenConns` limit?
2. Why does a leaked `Rows` exhaust the pool?
3. Why do you prefer `QueryContext` over `Query`?
4. What context do you pass in an HTTP handler?
5. Why do you set `ConnMaxLifetime`?

#### Easy practical tasks

1. After `Open`, set `MaxOpenConns` to `2`. Print `db.Stats()` after `Ping`.
2. Write a query function that takes `ctx context.Context` as the first parameter.
3. Read `go doc sql.DB.Stats`. Write three field names and what they count.
4. List four `*Context` methods on `sql.DB`.

#### Medium practical tasks

1. Run more concurrent queries than `MaxOpenConns` in a test. Show that some calls wait. Use a short timeout and record `ctx.Err()`.
2. Forget `rows.Close` in one function. Open many queries. Show `InUse` or wait errors. Then add `Close`.
3. Pass a cancelled context into `PingContext`. Record the error.

#### Advanced practical tasks

1. Tune `MaxOpenConns` for a small local database. Write the numbers and the symptom you used (wait time or DB `max_connections`).
2. Design a `Store` struct that holds `*sql.DB` and methods that take `ctx`. Do not put `ctx` on the struct.

---

## ORM vs raw SQL

Raw SQL means you write the statements. `database/sql` and a driver execute them. You control indexes, joins, and the exact SQL that you run in production.

An ORM (object-relational mapper) maps structs to tables. Some tools generate SQL. Some tools only scan rows into structs. Examples in the Go ecosystem include GORM and sqlc. sqlc generates Go from SQL. That style is not a full ORM.

Benefits of raw SQL:

- The statement is visible.
- You can copy the SQL into a database console.
- You avoid hidden queries.

Costs of raw SQL:

- You write scan code.
- You repeat column lists.

Benefits of an ORM or a mapper:

- Less scan code.
- A common pattern for simple CRUD.

Costs of an ORM:

- Hidden queries and N+1 loads
- Harder tuning
- A large dependency

Prefer raw SQL or sqlc when the team knows SQL. Use an ORM for simple CRUD only when the team already knows that tool. Do not mix two ORMs in one service.

Keep SQL in one package, such as `internal/store`. Do not spread queries across HTTP handlers.

`database/sql` still sits under many ORMs. You still need the pool, context, and transaction rules.

### Questions

#### Theoretical questions

1. What does raw SQL mean in this handbook?
2. What is one benefit of raw SQL?
3. What is one cost of an ORM?
4. How is sqlc different from a full ORM?
5. Where must SQL live in a module layout?

#### Easy practical tasks

1. Make a two-column table: "Raw SQL" and "ORM". Add four comparison rows.
2. Write a struct `Item` and the SQL `SELECT` that fills it. Do not use an ORM.
3. Read the homepage of sqlc or GORM. Write three facts. State which tool you would try first.
4. List three queries that you must understand even if an ORM writes them (insert, select by id, transaction).

#### Medium practical tasks

1. Implement `Get` and `Insert` with `database/sql` only. Keep SQL in `internal/store`.
2. Count the queries that one "list items plus owner name" page needs. Write how an ORM might hide a second query.
3. Replace a string-built filter with parameters. Show the SQL text that stays fixed.

#### Advanced practical tasks

1. Implement the same `Get` with sqlc or with an ORM. Compare generated code size and the SQL that you can read.
2. Write a team rule of one page: when to allow an ORM, when to require raw SQL, and how to review queries.

---

## Migrations

A migration is a change to the database schema. You store migrations as files. You apply them in order. Each file has a version.

Typical jobs:

- create a table
- add a column
- add an index
- backfill data in a separate step

Tools include `golang-migrate`, `goose`, and `atlas`. Some teams apply SQL files in a small Go command. Pick one tool. Do not edit a migration that already ran in production. Add a new version.

```text
migrations/0001_init.up.sql
migrations/0001_init.down.sql
```

The `up` file applies the change. The `down` file reverses it. Many teams stop writing `down` for production and restore from backup instead. If you write `down`, test it.

Run migrations before the new application version needs the new columns. Expand-contract:

1. Add a nullable column.
2. Deploy code that writes both old and new forms.
3. Backfill.
4. Deploy code that reads only the new form.
5. Remove the old column in a later migration.

Do not drop a column in the same release that still reads it.

Apply migrations from a controlled job. Do not race several replicas that all migrate. Use a lock that the tool provides.

Keep migration files in the repository. Embed them with `embed` if the migrate command lives in the same module.

Never put user data dumps in migration files.

### Questions

#### Theoretical questions

1. What is a migration?
2. Why must you not edit a migration that already ran in production?
3. What is expand-contract?
4. Why must two replicas not apply migrations at the same time without a lock?
5. What belongs in a migration file, and what does not?

#### Easy practical tasks

1. Write `0001_init.up.sql` that creates `items(id INTEGER PRIMARY KEY, name TEXT NOT NULL)`.
2. Write a matching `down` that drops `items`.
3. List three migration tools. Write one sentence for each.
4. Draw the five expand-contract steps for a rename of `name` to `title`.

#### Medium practical tasks

1. Apply `0001` to SQLite from a small Go program or a CLI tool. Query `sqlite_master` or list tables.
2. Add `0002` that adds `qty INTEGER NOT NULL DEFAULT 0`. Apply it. Insert a row.
3. Embed the SQL files with `embed.FS`. Print the file names at start.

#### Advanced practical tasks

1. Simulate two processes that migrate at once without a lock. Record the failure. Then use the lock that your tool provides.
2. Write a migration policy of one page: naming, review, expand-contract, and how you handle a failed `up`.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a driver, `sql.DB`, a transaction, and a context timeout share one query?
2. When do you use `QueryContext`, `QueryRowContext`, and `ExecContext`?
3. How do pool limits and `rows.Close` prevent connection exhaustion?
4. Why do parameters, migrations, and expand-contract all protect production data?
5. How do you decide among raw SQL, sqlc, and an ORM for one service?

#### Easy practical tasks

1. Open SQLite in memory. Create `items`. Insert one row. Query it with `QueryRowContext`.
2. Detect `ErrNoRows` for id `999`.
3. Write a cheat sheet: Open, Ping, Exec, Query, Scan, BeginTx, Commit, Rollback, Close.
4. Store the DSN in an environment variable name. Read it with `os.Getenv`. Do not print a password.

#### Medium practical tasks

1. Implement `internal/store` with `Insert` and `Get` in a transaction-free form. Add a transfer that uses `BeginTx`.
2. Set `MaxOpenConns(1)`. Run two timed queries. Show wait or timeout.
3. Add a migration file and apply it before `Ping` in `main`.

#### Advanced practical tasks

1. Build a tiny store: migrations, pool settings, context timeouts, parameterized SQL, and a test with SQLite.
2. Write a review checklist of ten items for a pull request that touches SQL. Include injection, `Close`, context, and schema order.
