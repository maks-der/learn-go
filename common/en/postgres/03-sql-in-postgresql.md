# 3. SQL in PostgreSQL

## Description

This topic shows core SQL in PostgreSQL 16 and PostgreSQL 17. You write `SELECT`, `INSERT`, `UPDATE`, and `DELETE`. You learn type output, implicit casts, and `RETURNING`. You also learn `LIMIT`, `OFFSET`, `FETCH`, `LIKE`, `ILIKE`, and dollar-quoting.

Complete topics 1 and 2 first. You need a connection and a schema where you can create tables. Complete this topic before you study types in depth.

Use one term for each concept. A statement is one SQL command. A cast changes a value from one type to another type. An implicit cast is a cast that you do not write. Dollar-quoting is a way to write a string without quote escaping.

---

## `SELECT` / `INSERT` / `UPDATE` / `DELETE`

These four statements are the data commands. `SELECT` reads rows. `INSERT` adds rows. `UPDATE` changes rows. `DELETE` removes rows.

Create a practice table:

```sql
CREATE TABLE items (
    id   integer PRIMARY KEY,
    name text NOT NULL,
    qty  integer NOT NULL DEFAULT 0
);
```

Insert rows:

```sql
INSERT INTO items (id, name, qty) VALUES (1, 'nail', 100);
INSERT INTO items (id, name, qty)
VALUES
    (2, 'screw', 50),
    (3, 'bolt', 25);
```

Select rows:

```sql
SELECT id, name, qty
FROM items
WHERE qty > 10
ORDER BY name;
```

Update rows:

```sql
UPDATE items
SET qty = qty + 5
WHERE name = 'nail';
```

Delete rows:

```sql
DELETE FROM items
WHERE id = 3;
```

Rules:

- `WHERE` limits `UPDATE` and `DELETE`. An `UPDATE` without `WHERE` changes every row. A `DELETE` without `WHERE` removes every row.
- `INSERT` must satisfy constraints. Topic 5 covers constraints.
- `SELECT` without `FROM` is valid. Example: `SELECT 1 + 1;`.
- PostgreSQL uses `$1` style parameters in client APIs. In `psql` you write literals. Do not concatenate user text into SQL. Use parameters in programs.

`INSERT ... SELECT` copies from a query:

```sql
INSERT INTO items (id, name, qty)
SELECT id + 100, name, qty FROM items WHERE qty < 10;
```

`UPDATE` can use `FROM` to join other tables. `DELETE` can use `USING`. Learn those forms after you learn joins in topic 6.

PostgreSQL does not require `FROM DUAL`. Use `SELECT` without `FROM` for expressions.

### Questions

#### Theoretical questions

1. What does each of `SELECT`, `INSERT`, `UPDATE`, and `DELETE` do?
2. What happens when `UPDATE` has no `WHERE` clause?
3. How do you insert more than one row in one statement?
4. Is `SELECT` without `FROM` valid?
5. Why must programs use parameters instead of string concatenation?

#### Easy practical tasks

1. Create `items` as in this section. Insert three rows. Select all columns.
2. Update one row by `id`. Select that row.
3. Delete one row. Select the remaining rows.
4. Run `SELECT 2 * 3 AS product;` with no `FROM`.

#### Medium practical tasks

1. Insert two rows in one `INSERT` with two `VALUES` tuples. Then insert a row with `INSERT ... SELECT`.
2. Write an `UPDATE` that increases `qty` by 1 for every row where `qty < 100`. Show `qty` before and after.
3. Run a `DELETE` with a `WHERE` that matches no row. Write how many rows PostgreSQL reports.

#### Advanced practical tasks

1. Create two tables. Use `UPDATE ... FROM` to copy a column from one table to the other. Use the official syntax.
2. Write a `psql` script that inserts, updates, and deletes, then ends with a `SELECT`. Run it with `psql -f`.

---

## PostgreSQL type output and implicit casts (careful)

Every type has an input function and an output function. `psql` shows the output form. The stored value can differ from the text that you see.

Examples:

```sql
SELECT 1;
SELECT 1.0;
SELECT DATE '2026-09-13';
SELECT TIMESTAMP '2026-09-13 12:00:00';
SELECT 'hello';
```

An untyped string literal has type `unknown` until PostgreSQL assigns a type. The context decides the type:

```sql
SELECT '2026-09-13';              -- text-like unknown
SELECT '2026-09-13'::date;        -- explicit date
SELECT DATE '2026-09-13';         -- typed literal
```

An implicit cast is a cast that you do not write. PostgreSQL adds it when the cast is marked as implicit in the catalogs. Example: `integer` can promote to `bigint` or `numeric` in many expressions.

Careful cases:

- `'1' = 1` can work because PostgreSQL casts the unknown or text value. Do not rely on this in application SQL. Write `1` for an integer column.
- `UNION` picks a common type. Mixed `integer` and `numeric` become `numeric`.
- `IN` lists and `CASE` branches follow type resolution rules. A mismatch raises an error or casts in a way that you did not expect.
- `char(n)` output pads with spaces. Comparison with `text` can surprise you. Prefer `text` or `varchar`.
- `timestamp` and `timestamptz` are different types. An implicit conversion uses the session `TimeZone`. Topic 4 covers this pair.

Write an explicit cast when the type matters:

```sql
SELECT '42'::integer;
SELECT CAST('42' AS integer);
```

See the cast:

```sql
SELECT pg_typeof(1), pg_typeof(1.0), pg_typeof('x');
```

`pg_typeof` shows the resolved type. Use it when a query fails with a type error.

Do not disable type checks. Do not store numbers as `text` to avoid casts.

### Questions

#### Theoretical questions

1. What is a type output function?
2. What is an implicit cast?
3. What type does an untyped string literal start as?
4. Why can `'1' = 1` be a bad pattern in application SQL?
5. What does `pg_typeof` show?

#### Easy practical tasks

1. Run `SELECT pg_typeof(1), pg_typeof(1.0), pg_typeof('1');`. Write the three types.
2. Cast `'2026-09-13'` to `date` in two ways (`::` and `CAST`).
3. Run `SELECT '1' = 1;`. Then run `SELECT '1'::text = 1;`. Record success or error for each.
4. Write four sentences: output text, stored type, implicit cast, explicit cast.

#### Medium practical tasks

1. Build a `UNION` of `integer` and `numeric` literals. Use `pg_typeof` on the result column.
2. Create a `char(5)` column and a `text` column. Insert `'abc'`. Select both. Compare them with `=`. Write what you see.
3. Find three implicit casts in `pg_cast` (`SELECT * FROM pg_cast` with a filter). Write the source type and the target type.

#### Advanced practical tasks

1. Read "Type Conversion" in the official docs. Write the resolution steps for an operator in your own short list.
2. Produce a type error with `CASE` or `UNION`. Fix it with an explicit cast. Save both SQL texts.

---

## `RETURNING`

`RETURNING` is a PostgreSQL extension. It returns rows from `INSERT`, `UPDATE`, or `DELETE` in the same statement. You do not need a second `SELECT`.

```sql
INSERT INTO items (id, name, qty)
VALUES (4, 'washer', 80)
RETURNING id, name;

UPDATE items
SET qty = qty - 1
WHERE id = 4
RETURNING id, qty;

DELETE FROM items
WHERE id = 4
RETURNING *;
```

`RETURNING *` returns every column of the affected row. After `UPDATE`, the values are the new values. After `DELETE`, the values are the old values.

Use `RETURNING` to get a generated key:

```sql
CREATE TABLE notes (
    id   integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    body text NOT NULL
);

INSERT INTO notes (body)
VALUES ('hello')
RETURNING id;
```

Topic 5 covers identity columns. Client drivers read the returned row with `Query` or `QueryRow`, not only `Exec`.

`RETURNING` can compute expressions:

```sql
UPDATE items
SET qty = qty + 10
WHERE id = 1
RETURNING id, qty, qty * 2 AS qty_double;
```

`WITH` can use a `RETURNING` result as a CTE. Topic 6 covers `WITH`.

Do not run a separate `SELECT` for the same key when `RETURNING` already gives the row. Do not assume `RETURNING` works in every other SQL database.

### Questions

#### Theoretical questions

1. Which statements accept `RETURNING`?
2. What does `RETURNING *` mean on `UPDATE`?
3. What does `RETURNING *` mean on `DELETE`?
4. Why do client programs use `RETURNING` for identity columns?
5. Is `RETURNING` in the SQL core that every database implements?

#### Easy practical tasks

1. Insert one `items` row with `RETURNING id, name`.
2. Update that row with `RETURNING qty`.
3. Delete that row with `RETURNING *`.
4. Insert into `notes` (or an identity table) and return only `id`.

#### Medium practical tasks

1. Update three rows in one statement. Return `id` and the new `qty`. Count the returned rows.
2. Use `RETURNING` with an expression column. Show the SQL and the result.
3. Compare `INSERT` plus a later `SELECT` with one `INSERT ... RETURNING`. Write when the single statement is better.

#### Advanced practical tasks

1. Write `WITH deleted AS (DELETE FROM ... RETURNING *) SELECT * FROM deleted;`. Use a copy of a practice table.
2. From a client library that you know (or `psql`), show how you read `RETURNING` columns. Write the API call names.

---

## `LIMIT` / `OFFSET` and `FETCH`

`LIMIT` restricts how many rows the query returns. `OFFSET` skips rows. PostgreSQL applies `OFFSET` first, then `LIMIT`, after `ORDER BY`.

```sql
SELECT id, name
FROM items
ORDER BY id
LIMIT 10;

SELECT id, name
FROM items
ORDER BY id
LIMIT 10 OFFSET 20;
```

Always use `ORDER BY` with `LIMIT` and `OFFSET`. Without `ORDER BY`, the row set is not a stable page.

`FETCH` is the SQL-standard form:

```sql
SELECT id, name
FROM items
ORDER BY id
FETCH FIRST 10 ROWS ONLY;

SELECT id, name
FROM items
ORDER BY id
OFFSET 20
FETCH NEXT 10 ROWS ONLY;
```

`FETCH FIRST 10 ROWS ONLY` matches `LIMIT 10`. You can write `ROW` or `ROWS`. You can write `FIRST` or `NEXT`.

`LIMIT ALL` means no limit. `OFFSET 0` skips nothing.

Problems with large `OFFSET`:

- The server still reads and sorts the skipped rows.
- Page 1000 is slower than page 1.
- Concurrent inserts change pages.

Topic 19 covers keyset pagination (`WHERE id > $1 ORDER BY id LIMIT 10`). Prefer keyset pagination for large lists.

`FETCH` also appears in cursors (`FETCH FORWARD`). That is a different command. This section means the `SELECT` clause.

### Questions

#### Theoretical questions

1. What does `LIMIT` do?
2. What does `OFFSET` do?
3. Why must `ORDER BY` appear with `LIMIT`?
4. What is the `FETCH FIRST` form?
5. Why is a large `OFFSET` slow?

#### Easy practical tasks

1. Insert at least 15 rows. Select the first 5 with `LIMIT` and `ORDER BY`.
2. Select the next 5 with `LIMIT` and `OFFSET`.
3. Rewrite the first query with `FETCH FIRST 5 ROWS ONLY`.
4. Run the same `LIMIT` without `ORDER BY` twice. Write if the order stayed the same.

#### Medium practical tasks

1. Write page 1 and page 3 of size 5 with `OFFSET`. Show both SQL statements.
2. Compare `LIMIT 5` and `FETCH NEXT 5 ROWS ONLY` on the same `ORDER BY`. Confirm equal rows.
3. Time a query with `OFFSET 0` and the same query with a large `OFFSET` on a table of a few thousand rows. Write the two times.

#### Advanced practical tasks

1. Implement the same page with `OFFSET` and with a keyset (`WHERE id > ...`). Compare plans with `EXPLAIN` (topic 9 helps).
2. Read the `SELECT` reference for `LIMIT` and `FETCH`. Write one option that this section did not show (`WITH TIES` if you use it).

---

## `ILIKE` vs `LIKE`

`LIKE` matches a string against a pattern. `%` matches any length. `_` matches one character.

```sql
SELECT name FROM items WHERE name LIKE 'n%';
SELECT name FROM items WHERE name LIKE '_olt';
```

`LIKE` is case-sensitive for the usual collations. `nail` does not match `LIKE 'N%'`.

`ILIKE` is the PostgreSQL case-insensitive form:

```sql
SELECT name FROM items WHERE name ILIKE 'N%';
```

`ILIKE` is not in the SQL standard. Other databases use `UPPER` and `LIKE`, or a separate collation.

Escape a literal `%` or `_` with `ESCAPE`:

```sql
SELECT name FROM items WHERE name LIKE '%\%%' ESCAPE '\';
```

You can write `NOT LIKE` and `NOT ILIKE`.

Locale and collation affect `ILIKE`. Some locales have special case rules. For a stable case-insensitive column, teams often use `citext` (topic 13) or store a normalized copy.

`LIKE 'abc%'` can use a B-tree index. `LIKE '%abc'` cannot use a normal B-tree index. Topic 8 covers indexes. Topic 7 covers `SIMILAR TO` and POSIX regular expressions.

Do not use `LIKE '%' || user_input || '%'` without a plan for performance and for wildcard characters in the input.

`~~` is the operator for `LIKE`. `~~*` is the operator for `ILIKE`. Prefer the keywords in application SQL.

### Questions

#### Theoretical questions

1. What do `%` and `_` mean in `LIKE`?
2. How is `ILIKE` different from `LIKE`?
3. Is `ILIKE` in the SQL standard?
4. How do you match a literal percent character?
5. Which `LIKE` pattern can use a B-tree index?

#### Easy practical tasks

1. Select names that start with a letter of your choice. Use `LIKE`.
2. Repeat the filter with `ILIKE` and a different letter case in the pattern.
3. Select names that contain `o` anywhere. Use `%`.
4. Run `NOT LIKE` and show the rows that fail the pattern.

#### Medium practical tasks

1. Insert a name that contains `%`. Write a `LIKE` query that finds that row and does not treat `%` as a wildcard.
2. Compare `LIKE 'A%'` and `ILIKE 'A%'` on mixed-case data. Write the two row counts.
3. Show `name ~~ 'n%'` and `name LIKE 'n%'`. Confirm that they match the same rows.

#### Advanced practical tasks

1. Read the docs for `LIKE` and collations. Write how locale can change `ILIKE`.
2. Build a case-insensitive search with `LOWER(name) LIKE LOWER($pattern)` and compare it with `ILIKE`. Write one benefit of each form.

---

## Dollar-quoting `$tag$ ... $tag$`

A normal string uses single quotes. A quote inside the string becomes two quotes:

```sql
SELECT 'it''s a nail';
```

Dollar-quoting writes a string between matching tags:

```sql
SELECT $$it's a nail$$;
SELECT $body$it's a nail$body$;
```

The form is `$tag$` + text + `$tag$`. The tag is optional. `$$` is valid. A tag is letters and digits. The tag is case-sensitive.

Use dollar-quoting for:

- function bodies (topic 13)
- long SQL in `DO` blocks
- strings that contain many quotes
- dynamic SQL that contains quotes

Example function body:

```sql
CREATE FUNCTION add_one(i integer)
RETURNS integer
LANGUAGE sql
AS $fn$
    SELECT i + 1;
$fn$;
```

The body is a string. Dollar tags avoid quote stacking.

Rules:

- The close tag must match the open tag.
- The text can contain `$other$` if `other` is not the same tag.
- Dollar-quoting does not expand escape sequences like `E'\n'`. Newlines stay as newlines.
- Do not mix `$$` in nested strings when an inner string also uses `$$`. Use a named tag.

`psql` variable interpolation still applies outside the rules you set. For scripts, prefer explicit tags.

Dollar-quoting is a PostgreSQL feature. Other databases may not accept it.

### Questions

#### Theoretical questions

1. What problem does dollar-quoting solve?
2. What is the shortest dollar-quote pair?
3. Why use a named tag such as `$fn$`?
4. Does dollar-quoting process `E'\n'` escapes?
5. Where do function bodies use dollar-quoting?

#### Easy practical tasks

1. Select a string that contains a single quote. Write it with doubled quotes and with `$$`.
2. Select a string with `$body$ ... $body$`.
3. Create a SQL function with a dollar-quoted body that returns `integer`.
4. Write four sentences: single quote, doubled quote, `$$`, named tag.

#### Medium practical tasks

1. Write a `DO` block with dollar-quoting that `RAISE NOTICE` a message. Run it.
2. Nest a string: outer tag `$outer$`, inner tag `$inner$`. Select the inner text.
3. Fail on purpose with a mismatched close tag. Record the error. Fix the tag.

#### Advanced practical tasks

1. Read `CREATE FUNCTION` in the docs. Rewrite one example from quote-escaping to dollar-quoting.
2. Generate dynamic SQL that contains quotes. Build the string with `format()` and dollar-quoting. Run it with `EXECUTE` in a `DO` block.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe one `INSERT` that uses dollar-quoting in a `text` column and returns the new `id`.
2. Why do implicit casts plus `LIMIT` without `ORDER BY` make a paged API unsafe?
3. When do you choose `ILIKE` instead of `LOWER(col) LIKE LOWER(pattern)`?
4. How do `RETURNING` and a later `SELECT` differ for a concurrent `UPDATE` on the same row?
5. What is the difference between a typed date literal and an unknown string that looks like a date?

#### Easy practical tasks

1. Create a table, insert two rows with one statement, update one row with `RETURNING`, and select with `LIKE`.
2. Write a cheat sheet: four data statements, `RETURNING`, `LIMIT`/`FETCH`, `LIKE`/`ILIKE`, dollar-quoting.
3. Page a sorted list with `FETCH FIRST 3 ROWS ONLY` and with `OFFSET`.
4. Use `pg_typeof` on a `RETURNING` column from an `INSERT` of a string into a `text` column.

#### Medium practical tasks

1. Write a `psql` script that creates a table, loads five rows, demos `ILIKE`, and prints `RETURNING` from a `DELETE`.
2. Find a type mismatch in a `WHERE` clause. Fix it with an explicit cast. Show `EXPLAIN` is not required; show the two SQL texts.
3. Document ten rules for safe SQL in PostgreSQL for a beginner teammate (parameters, `WHERE`, `ORDER BY`, casts).

#### Advanced practical tasks

1. Build a small paging procedure: page size and page number as arguments in a `DO` block or SQL function. Use `LIMIT` and `OFFSET`. Then write a second version that uses a keyset.
2. Read the `INSERT`, `UPDATE`, and `SELECT` reference pages. List three clauses that this topic did not cover. Write one sentence each.
