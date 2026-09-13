# 2. Types and Tables

## Description

This topic shows common data types and how you create tables in PostgreSQL 16 and PostgreSQL 17. You learn integers, exact numeric values, floats, text, boolean, date and time, `uuid`, and `bytea`. You also learn `NULL`, `CREATE TABLE`, identity columns, constraints, defaults, and table inheritance.

Complete topic 1 first. You must connect with `psql` and create a database. Complete this topic before you write joins and application SQL.

Use one term for each concept. A type is a name that PostgreSQL gives to a set of values. A table is a named set of rows. A constraint is a rule that a row must satisfy. An identity column is a column that PostgreSQL fills from a sequence. Inheritance copies column definitions from a parent table. Prefer `numeric` for money. Prefer `timestamptz` for instants. Prefer `text` for strings. Prefer composition (foreign keys) over inheritance for new designs.

---

## Integers, `numeric`, and floating types

PostgreSQL has three binary integer types:

| Type | Aliases | Size | Range |
| --- | --- | --- | --- |
| `smallint` | `int2` | 2 bytes | -32768 to 32767 |
| `integer` | `int`, `int4` | 4 bytes | about -2.1e9 to 2.1e9 |
| `bigint` | `int8` | 8 bytes | about -9e18 to 9e18 |

`integer` is the usual choice for keys and counts that fit in 32 bits. Use `bigint` for high-volume identity keys and large counters. Use `smallint` only when you have a small, closed set and you care about space.

Overflow raises an error:

```sql
SELECT 32767::smallint + 1::smallint;
```

Write integer literals without quotes: `42`. A quoted `'42'` is unknown or text until a cast.

Serial types (`smallserial`, `serial`, `bigserial`) are not true types. They create a sequence and a default. For new tables, prefer identity columns.

Integer division truncates toward zero:

```sql
SELECT 5 / 2;      -- 2 (integer)
SELECT 5::numeric / 2;
```

`numeric` and `decimal` are the same type in PostgreSQL. They store exact base-10 values. You set precision and scale:

```sql
amount numeric(12, 2)
```

Use `numeric` for money and for values that must not lose digits. Do not use `real` or `double precision` for money. Those types are binary floating point. They are approximate.

```sql
SELECT 0.1::numeric + 0.2::numeric;
SELECT 0.1::float8 + 0.2::float8;
```

`real` (`float4`) is 4 bytes. `double precision` (`float8`) is 8 bytes. Use them for scientific measures and for values where a small error is acceptable.

Use `bigint` for `id` when the table can grow past two billion rows. Changing `integer` to `bigint` later is a heavy rewrite.

### Questions

#### Theoretical questions

1. What is the storage size of `smallint`, `integer`, and `bigint`?
2. When do you choose `bigint` for a primary key?
3. What happens when an integer overflows?
4. Why is `numeric` better than `float8` for money?
5. Are `serial` types real types?

#### Easy practical tasks

1. Create a table with one column of each integer type. Insert a valid value in each column.
2. Run `SELECT pg_typeof(42), pg_typeof(42::bigint);`.
3. Insert `32767` into `smallint`. Then try `32768`. Record the error.
4. Compare `0.1::numeric + 0.2::numeric` with `0.1::float8 + 0.2::float8`. Write both results.

#### Medium practical tasks

1. Create a `numeric(12, 2)` column. Insert `19.999`. Write the stored value.
2. Compute `SUM` of `integer` values that overflow `integer`. Use `SUM` and `pg_typeof` on the result.
3. Read the docs for integer types. Write the exact min and max of `integer`.

#### Advanced practical tasks

1. Measure table size with `pg_relation_size` for 100000 `integer` ids versus `bigint` ids. Write the two sizes.
2. Plan a change from `integer` to `bigint` on a primary key. List the objects that you must change (indexes, foreign keys). Do not run a production rewrite.

---

## Text, `boolean`, date/time, `uuid`, `bytea`

Prefer `text` for strings. `varchar(n)` is `text` with a length check. `char(n)` pads with spaces. Do not use `char(n)` for new work.

```sql
SELECT pg_typeof('hello');
SELECT 'hello'::varchar(3);
```

String literals use single quotes. Dollar-quoting is topic 3.

`boolean` stores `true`, `false`, or `NULL`. Input accepts `t`/`f`, `true`/`false`, `yes`/`no`, `1`/`0`. Prefer `true` and `false` in SQL.

```sql
SELECT true IS TRUE;
SELECT NOT false;
```

Date and time types:

| Type | Meaning |
| --- | --- |
| `date` | calendar date |
| `time` | time of day without time zone |
| `timetz` | time of day with time zone (avoid for new work) |
| `timestamp` | date and time without time zone |
| `timestamptz` | date and time with time zone (stored in UTC) |
| `interval` | duration |

Prefer `timestamptz` for instants. Prefer `date` when you store only a calendar day. Prefer `interval` for durations.

```sql
SELECT now();
SELECT DATE '2026-09-13';
SELECT TIMESTAMPTZ '2026-09-13 12:00:00+00';
SELECT INTERVAL '1 day 2 hours';
```

`uuid` stores a 128-bit identifier. Generate values with `gen_random_uuid()` (built in on PostgreSQL 13 and later):

```sql
SELECT gen_random_uuid();
```

`bytea` stores raw bytes. Use it for small binary values. Do not store large files in `bytea` without a size plan. Prefer object storage for large blobs.

```sql
SELECT E'\\xDEADBEEF'::bytea;
```

`json` and `jsonb` are topic 4. Arrays are topic 4. This section is the scalar types that most tables need first.

### Questions

#### Theoretical questions

1. What is the difference between `text` and `varchar(n)`?
2. Why do you prefer `timestamptz` over `timestamp` for instants?
3. What three values can a `boolean` column hold?
4. What function generates a UUID in PostgreSQL 16 and 17 without an extra extension?
5. When do you avoid `bytea` for large files?

#### Easy practical tasks

1. Create a table with `text`, `boolean`, `date`, `timestamptz`, `uuid`, and `bytea`. Insert one valid row.
2. Run `SELECT now(), current_date, gen_random_uuid();`.
3. Insert `'yes'` into a `boolean` column. Select the stored value.
4. Cast `'2026-09-13 12:00:00+00'` to `timestamptz` and to `timestamp`. Write both results.

#### Medium practical tasks

1. Compare `varchar(5)` and `text` when you insert a six-character string. Record the error or the success.
2. Create a `timestamptz` column. Insert the same local clock time in two time zones. Select both rows.
3. Store a short PNG or a hex string in `bytea`. Select `length(data)` and `octet_length(data)`.

#### Advanced practical tasks

1. Read the docs for `timestamp` versus `timestamptz`. Write six sentences on storage and display.
2. Compare `gen_random_uuid()` with `pgcrypto` `gen_random_uuid()` history. Write which one you use on 16 and 17 and why.

---

## `NULL` and three-valued logic

`NULL` means "unknown" or "missing". It is not the string `'null'`. It is not `0`. It is not an empty string.

SQL uses three-valued logic: `true`, `false`, and `unknown`. Any comparison with `NULL` yields `unknown`, not `true`.

```sql
SELECT 1 = NULL;          -- NULL
SELECT NULL = NULL;       -- NULL
SELECT NULL IS NULL;      -- true
SELECT NULL IS NOT NULL;  -- false
```

`WHERE` keeps rows where the condition is `true`. A `WHERE col = NULL` clause keeps no row. Write `WHERE col IS NULL`.

`AND` and `OR` with `NULL`:

```sql
SELECT true AND NULL;     -- NULL
SELECT false AND NULL;    -- false
SELECT true OR NULL;      -- true
SELECT false OR NULL;     -- NULL
```

`IN` lists and `NOT IN` lists that contain `NULL` are a common trap. `x NOT IN (1, NULL)` is never `true`.

`CHECK` constraints treat `unknown` as success. `CHECK (qty > 0)` allows `qty` `NULL` unless you also write `NOT NULL`.

Unique constraints allow several `NULL` values in PostgreSQL for a single-column unique index (each `NULL` is distinct for uniqueness). A unique constraint on several columns treats `NULL` as distinct in that column. Know the rule before you model "optional unique email".

Aggregate functions ignore `NULL` inputs except `COUNT(*)`. `COUNT(col)` skips `NULL`. `SUM` of no non-null values is `NULL`.

Use `COALESCE` when you need a substitute (topic 4). Use `NOT NULL` when the value must exist.

### Questions

#### Theoretical questions

1. What does `1 = NULL` return?
2. Why does `WHERE col = NULL` return no row?
3. How does a `CHECK` constraint treat `unknown`?
4. Does `COUNT(col)` count `NULL` cells?
5. Can a unique column hold more than one `NULL` in PostgreSQL?

#### Easy practical tasks

1. Create a table with a nullable `text` column. Insert `NULL` and `''`. Select both rows with `IS NULL` and with `= ''`.
2. Run the four `AND`/`OR` examples from this section. Save the results.
3. Run `SELECT COUNT(*), COUNT(name) FROM t;` after you insert one `NULL` name.
4. Write five sentences: unknown, `IS NULL`, `WHERE`, `CHECK`, unique.

#### Medium practical tasks

1. Show that `NOT IN (1, NULL)` returns no row for a value `2`. Write why.
2. Add `CHECK (qty > 0)` without `NOT NULL`. Insert `NULL`. Then add `NOT NULL` and try again.
3. Create a unique column. Insert two `NULL` rows. Record the outcome.

#### Advanced practical tasks

1. Read "NULL values" and unique constraints in the 16 or 17 docs. Write the rule for a unique constraint on `(a, b)` when `b` is `NULL`.
2. Design a table where "missing" and "empty string" must stay different. Write the constraints and two sample `INSERT` statements.

---

## `CREATE TABLE` and `GENERATED ALWAYS AS IDENTITY`

`CREATE TABLE` defines a table name, columns, and table constraints.

```sql
CREATE TABLE customers (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email      text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
```

Column order is the order that `SELECT *` uses. Name columns explicitly in application SQL.

Useful clauses:

- `IF NOT EXISTS` — do not fail when the table exists
- `TEMPORARY` / `TEMP` — session table; drop at end of session (or transaction)
- `UNLOGGED` — skip WAL; faster; not crash-safe; not replicated
- `INHERITS` — see the inheritance section
- `PARTITION BY` — topic 10

```sql
CREATE TEMP TABLE scratch (id int);
CREATE UNLOGGED TABLE cache_keys (k text PRIMARY KEY, v text);
```

Use `TEMP` for session scratch data. Use `UNLOGGED` only when you can lose the data.

`CREATE TABLE ... AS SELECT` builds a table from a query. It does not copy constraints.

An identity column uses a sequence. PostgreSQL owns the sequence and ties it to the column.

```sql
CREATE TABLE orders (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    note text
);

INSERT INTO orders (note) VALUES ('first') RETURNING id;
```

Two modes:

- `GENERATED ALWAYS AS IDENTITY` — reject a user-supplied value unless `OVERRIDING SYSTEM VALUE`
- `GENERATED BY DEFAULT AS IDENTITY` — use the sequence when the client omits the column

Prefer `ALWAYS` for new keys. Prefer `bigint` identity for tables that can grow.

`serial` still works. It is a default plus a sequence. The sequence is easier to detach by mistake. Use identity for new tables.

`ALTER TABLE` adds columns and constraints later. Some `ALTER TABLE` forms rewrite the table. Test on a copy.

`DROP TABLE name;` removes the table. `DROP TABLE name CASCADE;` also drops dependent objects. Use `CASCADE` only when you intend that result.

Name tables in the singular or plural form that your team picks. Use one form. Use `snake_case`.

### Questions

#### Theoretical questions

1. What does `CREATE TABLE` define?
2. What is a `TEMPORARY` table?
3. What is an `UNLOGGED` table?
4. What is the difference between `GENERATED ALWAYS` and `GENERATED BY DEFAULT`?
5. Does `CREATE TABLE AS SELECT` copy primary keys?

#### Easy practical tasks

1. Create `customers` as in this section. Run `\d customers`.
2. Insert one row without `id`. Select the generated `id`.
3. Create a `TEMP` table. Insert one row. Reconnect. Confirm that the table is gone.
4. Drop a practice table without `CASCADE`.

#### Medium practical tasks

1. Try `INSERT` of an explicit `id` into a `GENERATED ALWAYS` column. Record the error. Then try `OVERRIDING SYSTEM VALUE`.
2. Create a table with `CREATE TABLE AS SELECT`. Compare `\d` with the source table.
3. Add a column with `ALTER TABLE ... ADD COLUMN`. Insert a row that uses the new column.

#### Advanced practical tasks

1. Compare a `serial` table and an identity table in `\d` and in `pg_get_serial_sequence`. Write three differences.
2. Read `CREATE TABLE` in the official docs. List five options that this section did not show.

---

## Primary key, unique, check, and foreign keys

A **primary key** uniquely identifies a row. It is `NOT NULL` and unique. A table can have one primary key. The primary key can use more than one column (composite).

A **unique** constraint forbids duplicate values in the listed columns. It can allow `NULL` as described earlier.

A **check** constraint is a boolean expression that each row must satisfy (or leave `unknown`).

A **foreign key** says that a value must match a unique key (usually a primary key) in another table, or be `NULL` if the column allows `NULL`.

```sql
CREATE TABLE customers (
    id    bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email text NOT NULL UNIQUE
);

CREATE TABLE orders (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    customer_id bigint NOT NULL REFERENCES customers (id),
    total       numeric(12, 2) NOT NULL CHECK (total >= 0)
);
```

Table-level form is required for multi-column keys:

```sql
PRIMARY KEY (id, created_on),
FOREIGN KEY (customer_id) REFERENCES customers (id)
```

Foreign key actions:

- `NO ACTION` / `RESTRICT` — reject a parent delete or update that still has children
- `CASCADE` — delete or update children
- `SET NULL` — set the child column to `NULL`
- `SET DEFAULT` — set the child column to its default

Write `ON DELETE` and `ON UPDATE` when you create the key. The default is `NO ACTION`.

PostgreSQL creates an index for primary keys and unique constraints. It does not create an index on the foreign-key column automatically. Add that index when you join or delete parents (topic 5).

`NOT NULL` is a constraint. Put it on every column that must have a value.

Do not use a natural key as the only key when the business value can change. Do not omit the foreign key and "check in the app only" if the database can enforce it.

### Questions

#### Theoretical questions

1. How many primary keys can one table have?
2. What does a foreign key require on the parent?
3. What does `ON DELETE CASCADE` do?
4. Does PostgreSQL index a foreign-key column by default?
5. When is a table-level constraint required?

#### Easy practical tasks

1. Create `customers` and `orders` as in this section. Insert one customer and one order.
2. Try an order with a missing `customer_id`. Record the error.
3. Add `CHECK (total >= 0)`. Try `total = -1`. Record the error.
4. Add a unique `email`. Insert the same email twice. Record the error.

#### Medium practical tasks

1. Create a composite primary key. Insert two rows that share one column but not both.
2. Compare `ON DELETE RESTRICT` and `ON DELETE CASCADE` on a copy of the tables. Delete a parent in each case.
3. Add a foreign key with `ALTER TABLE`. Confirm with `\d`.

#### Advanced practical tasks

1. Read deferrable constraints in the docs. Write when `DEFERRABLE INITIALLY DEFERRED` is useful. Run one insert pair that needs it.
2. Design keys for users, posts, and comments. Write `CREATE TABLE` with primary keys and foreign keys. Do not skip `ON DELETE`.

---

## Defaults and `now()`

A default fills a column when `INSERT` omits that column.

```sql
CREATE TABLE events (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    source     text        NOT NULL DEFAULT 'web'
);
```

`now()` is the transaction start time as `timestamptz`. All calls to `now()` in the same transaction return the same value. `clock_timestamp()` returns the wall clock at the call. `statement_timestamp()` returns the start of the current statement.

```sql
BEGIN;
SELECT now();
SELECT pg_sleep(2);
SELECT now(), clock_timestamp();
COMMIT;
```

Use `now()` (or `CURRENT_TIMESTAMP`, the SQL name for the same idea) for "when this row was created in this transaction". Use `clock_timestamp()` only when you need real elapsed time inside one transaction.

The default expression runs at insert time. It is not a trigger. A generated column is different (topic 9).

`DEFAULT` can call a function: `gen_random_uuid()`, `now()`, or a user function. The function must be allowed for the inserting role.

`INSERT ... DEFAULT VALUES` inserts a row that uses only defaults. `UPDATE ... SET col = DEFAULT` restores the default expression.

Do not store `now()` in the application only if the database must own the create time. Two app clocks can disagree. Do not use `timestamp` `DEFAULT now()` when you need a time zone.

### Questions

#### Theoretical questions

1. When does PostgreSQL evaluate a column default?
2. What time does `now()` use inside one transaction?
3. What is the difference between `now()` and `clock_timestamp()`?
4. Is a default the same object as a generated column?
5. What does `SET col = DEFAULT` do?

#### Easy practical tasks

1. Create `events` as in this section. Insert a row that omits `created_at` and `source`. Select the row.
2. Insert a row that overrides `source`. Confirm that `created_at` still uses the default.
3. Run the `BEGIN` / `pg_sleep` example. Write whether the two `now()` values match.
4. Run `UPDATE events SET source = DEFAULT WHERE id = 1;`. Select `source`.

#### Medium practical tasks

1. Compare `now()`, `statement_timestamp()`, and `clock_timestamp()` in one transaction with `pg_sleep(1)`.
2. Use `DEFAULT gen_random_uuid()` on a `uuid` column. Insert two rows. Confirm that the values differ.
3. Show that `DEFAULT now()` on `timestamp` (no time zone) stores a local value. Write the type and one risk.

#### Advanced practical tasks

1. Read "Date/Time Functions" in the 17 docs. Make a table of `now()`, `CURRENT_TIMESTAMP`, `transaction_timestamp()`, and `clock_timestamp()`.
2. Write a policy: which columns the database fills, and which columns the application sends. Include `created_at` and `updated_at` (trigger is topic 9).

---

## Table inheritance (know it exists; prefer composition)

Table inheritance copies columns from a parent table to a child table.

```sql
CREATE TABLE vehicles (
    id    bigint GENERATED ALWAYS AS IDENTITY,
    vin   text NOT NULL
);

CREATE TABLE cars (
    doors int NOT NULL
) INHERITS (vehicles);
```

`cars` has `id`, `vin`, and `doors`. `INSERT` into `cars` also makes the row visible in `SELECT` from `vehicles` unless you use `ONLY`:

```sql
SELECT * FROM vehicles;
SELECT * FROM ONLY vehicles;
```

Inheritance looks useful for "is-a" models. It has sharp limits:

- Unique constraints and primary keys on the parent do not cover children as one global unique set in the way beginners expect.
- Foreign keys to the parent do not see child rows as parent rows.
- Indexes do not cover the whole hierarchy as one index.
- Declarative partitioning (topic 10) is the supported way to split a large table.

Prefer composition. Put shared columns on one table. Point other tables at it with a foreign key. Example: `vehicles` as one table with a `kind` column, or `cars.vehicle_id REFERENCES vehicles (id)`.

Know that `INHERITS` exists so that you can read old schemas and the docs. Do not start a new design on inheritance. Do not mix inheritance with declarative partitioning on the same table.

PostgreSQL 16 and 17 still support inheritance. New features and tools target partitioned tables, not inheritance trees.

### Questions

#### Theoretical questions

1. What columns does a child table receive from `INHERITS`?
2. What does `SELECT FROM ONLY parent` skip?
3. Why can a parent primary key fail as a global unique key for children?
4. Can a foreign key to the parent see a row that lives only in a child?
5. What feature replaces inheritance for large-table splits?

#### Easy practical tasks

1. Create `vehicles` and `cars` as in this section. Insert one car. Select from `vehicles` and from `ONLY vehicles`.
2. Run `\d cars` and `\d vehicles`. List the columns.
3. Write four sentences: inherit, `ONLY`, foreign key limit, prefer composition.
4. Draw two designs: inheritance versus `cars.vehicle_id` foreign key.

#### Medium practical tasks

1. Try a unique `vin` on `vehicles`. Insert the same `vin` into two children (or parent and child). Record what PostgreSQL allows.
2. Try a foreign key from `trips.vehicle_id` to `vehicles(id)` and insert a trip for a car id. Record the result.
3. Read the 17 docs page "Inheritance". Write three warnings in your own words.

#### Advanced practical tasks

1. Compare inheritance and declarative partitioning in a one-page table: unique keys, foreign keys, planner pruning, official recommendation.
2. Find an old schema (or invent one) that uses `INHERITS`. Write a migration sketch to composition. Do not run it on production.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Which types do you pick for money, for an instant, and for a long string? Give one reason each.
2. How do `NULL`, a check constraint, and `NOT NULL` work together on one column?
3. Why is `GENERATED ALWAYS AS IDENTITY` the default choice over `serial` for a new key?
4. When do you add an extra index on a foreign-key column?
5. A teammate wants every vehicle kind as a child table. Which limits do you name first?

#### Easy practical tasks

1. Create `users` with identity `id`, unique `email`, `created_at timestamptz DEFAULT now()`, and `active boolean NOT NULL DEFAULT true`. Insert one row.
2. Add `posts` with a foreign key to `users`. Insert one post. Try a post for a missing user.
3. Run `\d` on both tables. Write the constraints that you see.
4. Insert a row that omits `created_at`. Select `created_at` and `pg_typeof(created_at)`.

#### Medium practical tasks

1. Build users, posts, and comments with foreign keys and `ON DELETE` actions that you choose. Insert a tree of three comments. Delete one post and record what happens to comments.
2. Add a `CHECK` that `rating` is between 1 and 5. Prove that `NULL` is allowed or forbidden based on your `NOT NULL` choice.
3. Compare table size of `numeric(12,2)` versus `float8` for 10000 amounts. Write which type you still pick for money and why.

#### Advanced practical tasks

1. Design a shop schema: customers, products, orders, order_lines. Write full `CREATE TABLE` with types, identity, defaults, and all four constraint kinds.
2. Write a one-page type policy for the team: integers, numeric, text, timestamptz, uuid, boolean, and when JSON waits until topic 4.
