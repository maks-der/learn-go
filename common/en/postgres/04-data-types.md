# 4. Data Types

## Description

This topic shows the common data types in PostgreSQL 16 and PostgreSQL 17. You learn integers, exact numeric values, approximate floats, text, boolean, date and time, `uuid`, JSON, `bytea`, arrays, and `NULL`.

Complete topic 3 first. You must be able to create a table and run `SELECT`. Complete this topic before you design constraints.

Use one term for each concept. A type is a name that PostgreSQL gives to a set of values. Storage size is the size on disk for a value of that type, not including tuple headers. Prefer the smallest integer that holds your range. Prefer `numeric` for money. Prefer `timestamptz` for instants. Prefer `text` for strings. Prefer `jsonb` for JSON that you search.

---

## Integers: `smallint`, `int`, `bigint`

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

Serial types (`smallserial`, `serial`, `bigserial`) are not true types. They create a sequence and a default. For new tables, prefer identity columns (topic 5).

Integer division truncates toward zero:

```sql
SELECT 5 / 2;      -- 2 (integer)
SELECT 5::numeric / 2;
```

Use `bigint` for `id` when the table can grow past two billion rows. Changing `integer` to `bigint` later is a heavy rewrite.

### Questions

#### Theoretical questions

1. What is the storage size of `smallint`, `integer`, and `bigint`?
2. When do you choose `bigint` for a primary key?
3. What happens when an integer overflows?
4. What is the result type of `5 / 2`?
5. Are `serial` types real types?

#### Easy practical tasks

1. Create a table with one column of each integer type. Insert a valid value in each column.
2. Run `SELECT pg_typeof(42), pg_typeof(42::bigint);`.
3. Insert `32767` into `smallint`. Then try `32768`. Record the error.
4. Make a table: type, size, one valid value, one value that does not fit.

#### Medium practical tasks

1. Compare `id integer` and `id bigint` in `\d` after you create two tables. Write the reported types.
2. Compute `SUM` of `integer` values that overflow `integer`. Use `SUM` and `pg_typeof` on the result.
3. Read the docs for integer types. Write the exact min and max of `integer`.

#### Advanced practical tasks

1. Measure table size with `pg_relation_size` for 100000 `integer` ids versus `bigint` ids. Write the two sizes.
2. Plan a change from `integer` to `bigint` on a primary key. List the objects that you must change (indexes, foreign keys). Do not run a production rewrite.

---

## `numeric` / `decimal` for money and exact values

`numeric` and `decimal` are the same type in PostgreSQL. They store exact base-10 values. You set precision and scale:

```sql
amount numeric(12, 2)
```

Precision `12` is the total count of digits. Scale `2` is the count of digits after the decimal point. `numeric` without precision stores values of variable size. Unconstrained `numeric` is slower than a constrained one for some work.

Use `numeric` for money, legal quantities, and any value that must not lose digits. Do not use `real` or `double precision` for money.

```sql
SELECT 0.1::numeric + 0.2::numeric;
SELECT 0.1::numeric + 0.2::numeric = 0.3::numeric;
```

`numeric` arithmetic is exact for values that fit. Division can produce many digits. Use `round`:

```sql
SELECT round(1::numeric / 3, 4);
```

Money type `money` exists. It depends on locale. Prefer `numeric` plus a currency code column.

Do not compare `numeric` and floats without an explicit plan. Cast the float to `numeric` or the reverse only when you accept rounding.

Storage grows with the number of digits. Very large precision is not free.

### Questions

#### Theoretical questions

1. Are `numeric` and `decimal` different types in PostgreSQL?
2. What is precision? What is scale?
3. Why is `numeric` the type for money?
4. Why avoid the `money` type for new work?
5. What does unconstrained `numeric` mean?

#### Easy practical tasks

1. Create a column `price numeric(10, 2)`. Insert `19.9` and `19.90`. Select both.
2. Run `0.1::numeric + 0.2::numeric` and compare with `0.3::numeric`.
3. Run `round(10::numeric / 3, 2)`.
4. Write four sentences: exact value, precision, scale, money.

#### Medium practical tasks

1. Insert a value that does not fit `numeric(4, 2)`. Record the error.
2. Compare `numeric(10, 2)` and unconstrained `numeric` in `\d`. Write the displayed types.
3. Sum a list of prices as `numeric`. Show that the sum matches a hand calculation.

#### Advanced practical tasks

1. Compare `0.1 + 0.2` as `numeric` and as `double precision`. Write both results and which type you use for money.
2. Read the `numeric` docs for `NaN`. Insert `NaN` if your version allows it. Write how comparisons with `NaN` behave.

---

## `real` / `double precision` (approximate)

`real` is 4-byte IEEE floating point (`float4`). `double precision` is 8-byte IEEE floating point (`float8`). These types are approximate. Many decimal fractions have no exact binary form.

```sql
SELECT 0.1::double precision + 0.2::double precision;
SELECT 0.1::float8 + 0.2::float8 = 0.3::float8;
```

The second query is often false. Do not use floats for money or for exact equality of decimals.

Use floats for scientific measurements, some statistics, and values where speed matters more than exact decimal digits.

`real` has about 6 decimal digits of precision. `double precision` has about 15. Prefer `double precision` when you need a float.

Special values: `Infinity`, `-Infinity`, and `NaN`.

```sql
SELECT 'Infinity'::float8;
SELECT 'NaN'::float8;
```

`NaN` comparisons are special. `NaN = NaN` is false. Use `x = x` is false as a NaN test, or use `x IS DISTINCT FROM` patterns with care. In PostgreSQL, `NaN` sorts as greater than other values in some contexts. Read the docs when you sort floats.

Do not mix `numeric` and floats in a money column. Do not store a float and expect a printed decimal to match the input text.

### Questions

#### Theoretical questions

1. What does "approximate" mean for `real` and `double precision`?
2. How many bytes does each type use?
3. Why can `0.1 + 0.2 = 0.3` fail for floats?
4. When is a float type acceptable?
5. Name three special float values.

#### Easy practical tasks

1. Run `0.1::float8 + 0.2::float8` and the equality test with `0.3::float8`.
2. Run the same two expressions with `numeric`.
3. Select `'Infinity'::float8` and `'NaN'::float8`.
4. Make a table: use case, choose `numeric` or `float8`. Add four rows.

#### Medium practical tasks

1. Insert the same decimal into `real` and `double precision`. Select with many digits (`extra_float_digits`). Compare the text.
2. Sort a column that contains `NaN` and normal numbers. Write the order.
3. Find `pg_typeof(1.23e4)`. Write the type of a scientific literal.

#### Advanced practical tasks

1. Read the official float docs. Write how PostgreSQL sorts `NaN`.
2. Write a short report: three bugs that appear when a team stores money in `double precision`.

---

## Text: `text`, `varchar`, `char`

PostgreSQL text types:

- `text` — string of any length (limited by row size)
- `varchar(n)` / `character varying(n)` — string up to `n` characters
- `varchar` without `n` — same behavior as `text`
- `char(n)` / `character(n)` — blank-padded to `n` characters

Prefer `text` or `varchar(n)` when you need a length limit as a constraint. Prefer `text` plus a `CHECK (char_length(col) <= n)` when you want a clear rule. `varchar(n)` rejects a longer string on input.

`char(n)` pads with spaces. That padding surprises comparisons and concatenation. Do not use `char(n)` for new application columns.

```sql
CREATE TABLE labels (
    a text,
    b varchar(10),
    c char(10)
);

INSERT INTO labels VALUES ('hi', 'hi', 'hi');
SELECT a, b, c, length(a), length(b), length(c);
```

`length` for `char(n)` counts the padded value. `text` and `varchar` do not pad.

PostgreSQL stores short strings inline. Long strings can use TOAST. You do not pick TOAST by hand.

There is no performance win for `varchar(10)` versus `text` for short strings. The old myth from other databases does not apply here.

Use `text` for names, comments, and tokens. Use `varchar(n)` when a standard or a form has a hard length.

### Questions

#### Theoretical questions

1. What is the difference between `text` and `varchar(n)`?
2. What does `char(n)` add on input?
3. Why does this handbook tell you to avoid `char(n)`?
4. Is `varchar` without `n` different from `text` in behavior?
5. Does `varchar(10)` store faster than `text` for short strings?

#### Easy practical tasks

1. Create the `labels` table. Insert `'hi'`. Select `length` for each column.
2. Insert a string of 11 characters into `varchar(10)`. Record the error.
3. Compare `'hi'::char(10) = 'hi'::text`. Write the result.
4. Write four sentences: `text`, `varchar(n)`, `char(n)`, TOAST.

#### Medium practical tasks

1. Add `CHECK (char_length(a) <= 10)` on a `text` column. Try a longer value.
2. Concatenate `char(10)` with `text`. Show the spaces in the result (`quote_literal`).
3. Read `\d+` on a table with a long `text` column after you insert a large value. Note TOAST if the docs or `\d+` mention it.

#### Advanced practical tasks

1. Read "Character Types" in the official docs. Quote the project advice about `char(n)` in your own words.
2. Compare `octet_length` and `char_length` for a string with non-ASCII characters. Write both numbers.

---

## `boolean`

`boolean` stores `true`, `false`, or `NULL`. Input accepts several texts:

```sql
SELECT true, false;
SELECT 't'::boolean, 'f'::boolean;
SELECT 'yes'::boolean, 'no'::boolean;
SELECT 'on'::boolean, 'off'::boolean;
SELECT '1'::boolean, '0'::boolean;
```

Prefer the tokens `true` and `false` in SQL. `psql` prints `t` and `f` by default.

Do not use `integer` 0 and 1 as a boolean column. Use `boolean`.

`WHERE flag` means `WHERE flag = true`. A `NULL` flag is unknown. The row does not pass. Topic section on `NULL` covers three-valued logic.

```sql
CREATE TABLE flags (
    id integer PRIMARY KEY,
    active boolean NOT NULL DEFAULT true
);
```

`NOT` inverts `true` and `false`. `NOT NULL` is still `NULL`.

Some drivers map `boolean` to a native bool. Check the driver when you scan values.

### Questions

#### Theoretical questions

1. What values can a `boolean` column hold?
2. Why prefer `boolean` over `integer` 0 and 1?
3. What does `WHERE flag` mean when `flag` is `NULL`?
4. What text does `psql` print for `true`?
5. What does `NOT NULL` return?

#### Easy practical tasks

1. Create `flags`. Insert `true` and `false`. Select the rows.
2. Cast `'yes'`, `'no'`, `'on'`, and `'off'` to `boolean`.
3. Run `SELECT NOT true, NOT false, NOT NULL;`.
4. Filter `WHERE active` and `WHERE NOT active`.

#### Medium practical tasks

1. Insert `NULL` into a nullable `boolean`. Test `WHERE flag`, `WHERE flag IS NULL`, and `WHERE flag IS NOT TRUE`.
2. Compare `\pset null` display with a `NULL` boolean in `psql`.
3. Read the `boolean` docs. List all accepted input strings.

#### Advanced practical tasks

1. From a client language, insert and read `boolean` and `NULL`. Write the host types.
2. Write three `CHECK` examples that use `boolean` columns (example: one of two flags must be true).

---

## Date/time: `date`, `time`, `timestamp`, `timestamptz`

Common types:

- `date` — calendar date, no time of day
- `time` / `time without time zone` — time of day
- `timetz` / `time with time zone` — avoid for new work
- `timestamp` / `timestamp without time zone` — date and time, no zone
- `timestamptz` / `timestamp with time zone` — an instant; stored in UTC

Prefer `date` when you need only a day. Prefer `timestamptz` when you need an instant (created_at, event time). Prefer `timestamp` when you need a local date-time without a zone (a wall clock in one building). Do not mix the two timestamp types without a rule.

```sql
SELECT DATE '2026-09-13';
SELECT TIMESTAMP '2026-09-13 15:30:00';
SELECT TIMESTAMPTZ '2026-09-13 15:30:00+02';
SHOW TimeZone;
```

PostgreSQL stores `timestamptz` as UTC. Display uses the session `TimeZone`.

```sql
SET TimeZone TO 'UTC';
SELECT TIMESTAMPTZ '2026-09-13 15:30:00+02';
SET TimeZone TO 'Europe/Kyiv';
SELECT TIMESTAMPTZ '2026-09-13 15:30:00+02';
```

`now()` is the start time of the current transaction. `CURRENT_TIMESTAMP` is the same as `now()`. `clock_timestamp()` changes during the transaction. Topic 5 covers defaults.

`interval` is a duration:

```sql
SELECT DATE '2026-09-13' + INTERVAL '7 days';
SELECT now() - INTERVAL '1 hour';
```

Do not store Unix epoch integers when `timestamptz` works. Do not use `timetz`.

### Questions

#### Theoretical questions

1. What does `date` store?
2. What is the difference between `timestamp` and `timestamptz`?
3. How does PostgreSQL store `timestamptz`?
4. What session setting changes `timestamptz` display?
5. Why avoid `timetz`?

#### Easy practical tasks

1. Select `DATE '2026-09-13'` and `pg_typeof` of that value.
2. Insert `now()` into a `timestamptz` column. Select it.
3. `SET TimeZone` to `UTC` and to a named zone. Select the same `timestamptz` twice.
4. Add `INTERVAL '1 day'` to a `date`.

#### Medium practical tasks

1. Insert the same literal into `timestamp` and `timestamptz`. Change `TimeZone`. Compare the two columns.
2. Compare `now()` and `clock_timestamp()` inside a transaction that includes `pg_sleep(1)`.
3. Compute age with `age(date_col)` or `now() - timestamptz_col`.

#### Advanced practical tasks

1. Read "Date/Time Types" in the docs. Write how DST affects `timestamp` versus `timestamptz`.
2. Convert a Unix epoch (seconds) to `timestamptz` with `to_timestamp`. Convert back with `extract(epoch FROM ...)`.

---

## `uuid`, `json` / `jsonb`, `bytea`

`uuid` stores a 128-bit identifier. The text form is `a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11`. Generate values with `gen_random_uuid()` (in core since PostgreSQL 13):

```sql
SELECT gen_random_uuid();
CREATE TABLE events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    payload jsonb NOT NULL
);
```

`json` stores JSON text. It keeps whitespace and key order. It accepts duplicate keys. `jsonb` stores a parsed binary form. It drops insignificant whitespace. It does not keep duplicate keys (the last key wins). `jsonb` supports indexes and containment operators. Prefer `jsonb` for columns that you search. Prefer `json` only when you must keep the exact input text.

```sql
SELECT '{"a":1}'::json;
SELECT '{"a":1}'::jsonb;
SELECT '{"a":1}'::jsonb @> '{"a":1}';
```

Topic 7 covers JSON operators. Topic 14 covers JSON in depth.

`bytea` stores raw bytes. Input can be hex format:

```sql
SELECT '\xDEADBEEF'::bytea;
```

Use `bytea` for small binary values. Large files often belong in object storage. A very large `bytea` uses TOAST and memory in the client.

Do not store JSON as unconstrained `text` if you need validation. `json` and `jsonb` reject invalid JSON on input.

### Questions

#### Theoretical questions

1. What does `uuid` store?
2. What function generates a random UUID in current PostgreSQL?
3. How is `json` different from `jsonb`?
4. Which JSON type do you use when you search keys?
5. What does `bytea` store?

#### Easy practical tasks

1. Run `SELECT gen_random_uuid();` twice. Confirm the values differ.
2. Cast a small object to `json` and to `jsonb`. Select both.
3. Insert invalid JSON into a `jsonb` column. Record the error.
4. Select `'\x00FF'::bytea`.

#### Medium practical tasks

1. Create `events` as in this section. Insert two rows. Select `id` and `payload`.
2. Compare `'{"b":1,"a":2}'::json` and the same text as `jsonb`. Write what changes.
3. Use `octet_length` on a `bytea` value. Write the length.

#### Advanced practical tasks

1. Read the `jsonb` docs on duplicate keys and key order. Write two examples.
2. Compare `uuid` and `bigint` identity for a primary key. Write two reasons for each choice.

---

## Arrays

An array is an ordered list of values of one element type. Write the type as `integer[]` or `text[]`.

```sql
CREATE TABLE packs (
    id integer PRIMARY KEY,
    tags text[],
    nums integer[]
);

INSERT INTO packs VALUES (1, ARRAY['red', 'blue'], ARRAY[1, 2, 3]);
INSERT INTO packs VALUES (2, '{green,blue}', '{4,5}');
```

Subscripts start at 1:

```sql
SELECT nums[1], nums[2] FROM packs WHERE id = 1;
```

Useful operators and functions:

- `= ` — equal arrays
- `@>` — left contains right
- `<@` — left is contained by right
- `&&` — overlap
- `unnest(arr)` — one row per element
- `array_length(arr, 1)` — length of the first dimension
- `cardinality(arr)` — total elements

```sql
SELECT * FROM packs WHERE tags @> ARRAY['blue'];
SELECT unnest(tags) FROM packs WHERE id = 1;
```

Arrays can be multidimensional. Most applications use one dimension.

`NULL` elements are allowed. An array value can also be `NULL`.

Do not use arrays as a replacement for a child table when you need keys, constraints, and joins. Use arrays for small sets (tags, days of week). Topic 14 covers `unnest` in more depth.

A GIN index can speed containment on arrays (topic 8).

### Questions

#### Theoretical questions

1. What is the element type of `integer[]`?
2. What is the first subscript in PostgreSQL arrays?
3. What does `unnest` do?
4. What does `@>` mean for arrays?
5. When must you use a child table instead of an array?

#### Easy practical tasks

1. Create `packs`. Insert one row with `ARRAY[...]`.
2. Select the first element of `nums`.
3. Filter rows with `@>` and one tag.
4. `SELECT unnest(ARRAY['a','b']);`.

#### Medium practical tasks

1. Insert an array that contains `NULL`. Select `array_length` and `cardinality`.
2. Update `tags` with `array_append` or concatenation `||`.
3. Compare `' {1,2,3}'::integer[]` and `ARRAY[1,2,3]`. Confirm equality.

#### Advanced practical tasks

1. Create a GIN index on a `text[]` column (topic 8). Write the `CREATE INDEX` statement and a matching `WHERE` clause.
2. Read "Arrays" in the docs. Write how a slice `nums[1:2]` works. Show one query.

---

## `NULL` and three-valued logic

`NULL` means unknown or missing. It is not the number 0. It is not an empty string. Three-valued logic uses `true`, `false`, and `unknown`.

```sql
SELECT 1 = 1;          -- true
SELECT 1 = 2;          -- false
SELECT 1 = NULL;       -- NULL (unknown)
SELECT NULL = NULL;    -- NULL (unknown)
```

`WHERE` keeps rows where the condition is `true`. `WHERE col = NULL` never keeps a row. Use `IS NULL` and `IS NOT NULL`:

```sql
SELECT * FROM items WHERE qty IS NULL;
SELECT * FROM items WHERE qty IS NOT NULL;
```

`IS DISTINCT FROM` treats two `NULL` values as not distinct:

```sql
SELECT NULL IS DISTINCT FROM NULL;     -- false
SELECT 1 IS DISTINCT FROM NULL;        -- true
```

`AND` and `OR` with `NULL`:

- `true AND NULL` is `NULL`
- `false AND NULL` is `false`
- `true OR NULL` is `true`
- `false OR NULL` is `NULL`

`NOT NULL` is `NULL`.

Aggregates such as `SUM` and `AVG` skip `NULL` inputs. `COUNT(col)` skips `NULL`. `COUNT(*)` counts rows.

A unique constraint allows many `NULL` values in PostgreSQL (NULL is distinct from NULL in that rule), unless you use `UNIQUE NULLS NOT DISTINCT` (PostgreSQL 15 and later).

Always decide `NULL` or `NOT NULL` on each column. Prefer `NOT NULL` when a value must exist.

### Questions

#### Theoretical questions

1. What does `NULL` mean?
2. Why is `WHERE col = NULL` wrong?
3. What does `IS DISTINCT FROM` do with two `NULL` values?
4. Does `COUNT(col)` count `NULL` cells?
5. What is `UNIQUE NULLS NOT DISTINCT`?

#### Easy practical tasks

1. Run `SELECT 1 = NULL, NULL = NULL, NULL IS NULL;`.
2. Insert a row with a `NULL` quantity. Select it with `IS NULL`.
3. Run `SELECT COUNT(*), COUNT(qty) FROM items;` when one `qty` is `NULL`.
4. Write the four `AND`/`OR` results from this section.

#### Medium practical tasks

1. Build a `UNIQUE` column that allows two `NULL` values. Then try `UNIQUE NULLS NOT DISTINCT` if you use PostgreSQL 15 or later.
2. Use `COALESCE(qty, 0)` in a `SUM`. Compare with `SUM(qty)`.
3. Test `WHERE flag IS NOT TRUE` on `true`, `false`, and `NULL` boolean values.

#### Advanced practical tasks

1. Read "Comparison Functions" for `IS NOT DISTINCT FROM`. Rewrite a join condition that must treat `NULL` keys as equal.
2. Write a one-page guide: when a column is `NOT NULL`, when it is nullable, and how `CHECK` interacts with `NULL` (unknown passes `CHECK` unless you write `IS NOT NULL`).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Pick a type for money, a type for an instant, a type for a tag list, and a type for a search document in JSON. Give one reason each.
2. How do three-valued logic and a `boolean` column change a `WHERE` clause?
3. Why do `char(n)` and `real` both cause equality surprises, for different reasons?
4. When do you store `uuid` instead of `bigint` for a key?
5. What is the difference between a `NULL` array and an array of `NULL` elements?

#### Easy practical tasks

1. Create one table that uses `bigint`, `numeric(12,2)`, `text`, `boolean`, `timestamptz`, `uuid`, `jsonb`, and `text[]`. Insert one valid row.
2. Write a cheat sheet with every type in this topic and one "use this when" sentence.
3. Show `pg_typeof` for an integer literal, a numeric literal, a date literal, and `gen_random_uuid()`.
4. Select a row with `IS NULL` and a row with `IS NOT NULL` on the same column.

#### Medium practical tasks

1. Load five rows with mixed types. Write three `WHERE` clauses that fail because of types or `NULL`. Fix each clause.
2. Change session `TimeZone` and compare a `timestamp` column with a `timestamptz` column on the same insert literal.
3. Document your type rules for a new project in ten lines (keys, money, time, text, JSON, arrays, `NULL`).

#### Advanced practical tasks

1. Design a `payments` table and an `events` table with types only (no extra features from later topics except `NOT NULL`). Write `CREATE TABLE` and five sample `INSERT` statements.
2. Read the "Data Types" chapter contents page. List four types that this topic did not cover. Write one sentence each for when you would open that page.
