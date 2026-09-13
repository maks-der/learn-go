# 7. Functions and Operators

## Description

This topic shows common functions and operators in PostgreSQL 16 and PostgreSQL 17. You learn string, date, and math functions. You learn `COALESCE`, `NULLIF`, `GREATEST`, and `LEAST`. You also learn pattern matching, a short full-text search preview, and JSON operators.

Complete topic 6 first. You can write joins and CTEs. Complete this topic before you study indexes.

Use one term for each concept. A function returns a value from arguments. An operator is a symbol such as `||` or `@>`. Prefer documented function names in application SQL. Prefer operators when they are the usual form (`||` for text, `@>` for `jsonb`).

---

## String, date, and math functions

String functions work on `text` and `varchar`.

```sql
SELECT length('PostgreSQL');
SELECT char_length('PostgreSQL');
SELECT lower('ABC'), upper('abc');
SELECT trim(BOTH ' ' FROM '  nail  ');
SELECT substring('washer' FROM 1 FOR 3);
SELECT left('washer', 3), right('washer', 2);
SELECT replace('red nail', 'red', 'blue');
SELECT split_part('a,b,c', ',', 2);
SELECT position('SQL' IN 'PostgreSQL');
```

`||` concatenates text. `NULL` in a `||` expression makes the result `NULL`. Use `concat` or `concat_ws` when you must skip `NULL`:

```sql
SELECT 'a' || NULL;           -- NULL
SELECT concat('a', NULL, 'b'); -- ab
SELECT concat_ws(',', 'a', NULL, 'c');
```

Date and time functions:

```sql
SELECT now();
SELECT current_date;
SELECT date_trunc('hour', now());
SELECT extract(year FROM DATE '2026-09-13');
SELECT date_part('month', DATE '2026-09-13');
SELECT age(DATE '2026-01-01');
SELECT DATE '2026-09-13' + INTERVAL '1 month';
SELECT timestamptz '2026-09-13 12:00:00+00' AT TIME ZONE 'Europe/Kyiv';
```

`date_trunc` cuts a timestamp to a unit (hour, day, month). `extract` returns a field as `numeric`. `age` returns an `interval`.

Math functions:

```sql
SELECT abs(-3);
SELECT round(12.56, 1);
SELECT ceil(12.1), floor(12.9);
SELECT mod(10, 3);
SELECT power(2, 10);
SELECT sqrt(9);
```

Integer division truncates. Use `numeric` when you need decimal division.

`generate_series` builds a set of values:

```sql
SELECT * FROM generate_series(1, 5);
SELECT * FROM generate_series(
    timestamp '2026-01-01',
    timestamp '2026-01-07',
    interval '1 day'
);
```

Read the official function chapters when you need a function that this list does not show. Do not invent names from other databases (`IFNULL`, `DATEDIFF`). PostgreSQL uses `COALESCE` and date subtraction.

### Questions

#### Theoretical questions

1. What does `||` do for text? What happens with `NULL`?
2. What does `date_trunc('day', ts)` return?
3. What does `extract(year FROM d)` return?
4. Why is `5 / 2` not `2.5` for integers?
5. What does `generate_series(1, 5)` return?

#### Easy practical tasks

1. Run `lower`, `upper`, `trim`, and `length` on one string.
2. Concatenate first name and last name with `||` and with `concat_ws`.
3. Truncate `now()` to day and to month.
4. Select `round(2::numeric / 3, 4)`.

#### Medium practical tasks

1. Split a CSV string with `string_to_array` or `split_part`. Make one row per part with `unnest`.
2. Build seven dates with `generate_series` and format them with `to_char`.
3. Compute the number of days between two dates (`date - date`). Write the type of the result.

#### Advanced practical tasks

1. Read "String Functions and Operators" and list five functions that this section omitted. Write one sentence each.
2. Convert a `timestamptz` to a named zone with `AT TIME ZONE` and back. Explain the result types.

---

## `COALESCE`, `NULLIF`, `GREATEST`, `LEAST`

`COALESCE(a, b, ...)` returns the first argument that is not `NULL`.

```sql
SELECT COALESCE(NULL, NULL, 'nail');
SELECT COALESCE(nickname, email, 'unknown') FROM users;
```

Use `COALESCE` to replace `NULL` in display and in math (`COALESCE(qty, 0)`).

`NULLIF(a, b)` returns `NULL` when `a` equals `b`. Otherwise it returns `a`.

```sql
SELECT NULLIF(0, 0);          -- NULL
SELECT NULLIF(qty, 0);        -- NULL when qty is 0
SELECT 100 / NULLIF(qty, 0);  -- avoids division by zero; result NULL
```

`GREATEST(a, b, ...)` returns the largest argument. `LEAST` returns the smallest. In PostgreSQL, these functions skip `NULL` arguments. They return `NULL` only when every argument is `NULL`.

```sql
SELECT GREATEST(1, 5, 3);
SELECT LEAST(1, 5, 3);
SELECT GREATEST(1, NULL, 3);  -- 3
SELECT GREATEST(NULL, NULL);  -- NULL
```

Some other databases return `NULL` from `GREATEST` when any argument is `NULL`. Do not copy that rule to PostgreSQL.

`GREATEST` and `LEAST` need a common type. Mix types with care. Use an explicit cast when the type is not clear.

These functions are not aggregates. Aggregates are `max` and `min` over rows. `GREATEST` compares values in one row.

```sql
SELECT max(qty) FROM items;           -- aggregate
SELECT GREATEST(qty, min_qty) FROM items; -- per row
```

### Questions

#### Theoretical questions

1. What does `COALESCE` return?
2. What does `NULLIF(a, b)` return when `a = b`?
3. How does PostgreSQL treat `NULL` in `GREATEST`?
4. How is `GREATEST` different from `max`?
5. How can `NULLIF` prevent division by zero?

#### Easy practical tasks

1. Run `COALESCE(NULL, 2, 3)`.
2. Run `NULLIF('a', 'a')` and `NULLIF('a', 'b')`.
3. Run `GREATEST(1, NULL, 4)` and `LEAST(1, NULL, 4)`.
4. Replace `NULL` qty with 0 in a `SELECT` list.

#### Medium practical tasks

1. Use `COALESCE` on three columns to build a display name.
2. Compute `amount / NULLIF(days, 0)` on a small table that contains 0.
3. Compare `GREATEST(a, b)` with a `CASE` that does the same work.

#### Advanced practical tasks

1. Read the official notes for `GREATEST` and `LEAST`. Write how PostgreSQL differs from databases that return `NULL` if any argument is `NULL`.
2. Combine `COALESCE` and `NULLIF` to treat both `NULL` and `''` as missing text.

---

## Pattern matching: `LIKE`, `SIMILAR TO`, POSIX regex

PostgreSQL has three pattern systems.

**`LIKE` / `ILIKE`** (topic 3): `%` and `_`. `ILIKE` is case-insensitive. Fast to learn. Enough for many filters.

**`SIMILAR TO`** uses SQL-standard regular expressions. The pattern must match the whole string. Syntax looks like `LIKE` plus extra tokens (`|`, `*`, `+`, parentheses).

```sql
SELECT name FROM items WHERE name SIMILAR TO '%(nail|screw)%';
```

**POSIX regular expressions** use operators:

| Operator | Meaning |
| --- | --- |
| `~` | match, case-sensitive |
| `~*` | match, case-insensitive |
| `!~` | no match, case-sensitive |
| `!~*` | no match, case-insensitive |

```sql
SELECT name FROM items WHERE name ~ '^[Nn]ail';
SELECT name FROM items WHERE name ~* 'nail';
```

POSIX patterns may match a substring. You do not need `%` at both ends.

Functions:

```sql
SELECT substring('item-42' FROM '[0-9]+');
SELECT regexp_match('item-42', '[0-9]+');
SELECT regexp_replace('a1b2', '[0-9]', 'x', 'g');
SELECT regexp_split_to_table('a,b,c', ',');
```

`regexp_match` returns a `text[]`. `regexp_matches` can return many rows.

Which form to use:

- `LIKE` / `ILIKE` — prefix and simple contains
- `SIMILAR TO` — rare; some teams skip it
- POSIX — validation, extraction, complex search

A leading wildcard or a complex regex often cannot use a normal B-tree index. `pg_trgm` (topic 13) can support `LIKE '%x%'` and some regex with a GIN or GiST index.

Do not run untrusted user regex without a time limit. A bad pattern can use much CPU time. Set `statement_timeout` on application roles.

### Questions

#### Theoretical questions

1. What are the three pattern systems in this section?
2. Does `SIMILAR TO` match the whole string?
3. What is the difference between `~` and `~*`?
4. Does `name ~ 'nail'` need `%`?
5. Why can a user-supplied regex be a risk?

#### Easy practical tasks

1. Filter with `LIKE '%a%'`, `ILIKE '%A%'`, and `~* 'a'`.
2. Use `SIMILAR TO` with two alternatives `|`.
3. Extract digits with `substring(... FROM '[0-9]+')`.
4. Replace all digits with `regexp_replace` and flag `g`.

#### Medium practical tasks

1. Validate an email-like pattern with `~` (keep the pattern simple). Show one match and one miss.
2. Compare `LIKE 'n%'` and `~ '^n'`. Confirm the same rows on your data.
3. Split text to rows with `regexp_split_to_table`.

#### Advanced practical tasks

1. Read "POSIX Regular Expressions" in the docs. Write the meaning of `g`, `i`, and `n` flags for `regexp_replace`.
2. Enable `pg_trgm` on a test database (topic 13). Create a GIN index. Compare `EXPLAIN` for `LIKE '%nail%'` before and after. Skip if you cannot create extensions.

---

## Full-text search preview (`tsvector`, `tsquery`)

Full-text search finds words in documents. It is not `LIKE '%word%'`. PostgreSQL stems words, drops stop words, and ranks results.

`tsvector` is a sorted list of lexemes. `tsquery` is a search query.

```sql
SELECT to_tsvector('english', 'The nails were shining');
SELECT to_tsquery('english', 'nail');
SELECT to_tsvector('english', 'The nails were shining')
       @@ to_tsquery('english', 'nail');
```

`@@` means "the vector matches the query". Stemming maps `nails` to `nail` in the `english` configuration.

Create a generated column or a stored `tsvector` for a table (topic 13 covers generated columns):

```sql
CREATE TABLE posts (
    id      bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title   text NOT NULL,
    body    text NOT NULL,
    tsv     tsvector GENERATED ALWAYS AS (
                to_tsvector('english', title || ' ' || body)
            ) STORED
);
```

Search:

```sql
SELECT id, title
FROM posts
WHERE tsv @@ to_tsquery('english', 'nail & shine')
ORDER BY ts_rank(tsv, to_tsquery('english', 'nail & shine')) DESC;
```

Query forms:

- `to_tsquery` — operators `&` `|` `!` and parentheses; you must write lexemes
- `plainto_tsquery` — plain text to AND of words
- `websearch_to_tsquery` — a web-style string (quotes, `-`)
- `phraseto_tsquery` — phrase

Use a GIN index on `tsvector` for large tables (topic 8).

This section is a preview. Dictionaries, weights (`setweight`), and ranking need more reading in the official "Full Text Search" chapter.

Do not use full-text search for exact SKU lookup. Use `=` or `LIKE` prefix. Do not mix languages in one `english` column without a plan.

### Questions

#### Theoretical questions

1. What does `tsvector` store?
2. What does `tsquery` store?
3. What does `@@` mean?
4. Why can `nails` match `nail` in `english`?
5. What is `plainto_tsquery`?

#### Easy practical tasks

1. Run `to_tsvector('english', 'The nails were shining')`.
2. Test `@@` with `to_tsquery('english', 'nail')`.
3. Test a query that must not match. Write the result.
4. Run `plainto_tsquery('english', 'shining nails')`.

#### Medium practical tasks

1. Create `posts` with a stored `tsvector` column (or compute `to_tsvector` in `WHERE` if you skip generated columns).
2. Insert two posts. Search with `&` and with `|`.
3. Order results with `ts_rank`.

#### Advanced practical tasks

1. Read "Text Search Functions" for `websearch_to_tsquery`. Write one query that uses quotes or `-`.
2. Compare `LIKE '%nail%'` and `@@ to_tsquery('english', 'nail')` on words like `nails` and `email`. Write which rows each method returns.

---

## JSON operators: `->`, `->>`, `@>`, `?`

These operators work on `json` and `jsonb`. Prefer `jsonb` for operators and indexes.

| Operator | Meaning |
| --- | --- |
| `->` | get JSON object field or array element; result is `json` / `jsonb` |
| `->>` | get field or element as `text` |
| `@>` | does the left `jsonb` contain the right `jsonb`? |
| `?` | does the object contain this top-level key? |

```sql
CREATE TABLE events (
    id      bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    payload jsonb NOT NULL
);

INSERT INTO events (payload) VALUES
    ('{"type":"click","n":2,"tags":["a","b"]}');

SELECT payload -> 'type'     AS as_json;
SELECT payload ->> 'type'    AS as_text;
SELECT payload -> 'n';
SELECT (payload ->> 'n')::int;
SELECT payload @> '{"type":"click"}';
SELECT payload ? 'type';
```

`->` keeps JSON type. `->>'type'` is `text`. Cast text to `integer` when you need math.

Array index with `->` is 0-based for JSON arrays:

```sql
SELECT payload -> 'tags' -> 0;
SELECT payload -> 'tags' ->> 0;
```

Related operators (learn the names):

- `?|` — any of these keys exist
- `?&` — all of these keys exist
- `||` — concatenate `jsonb`
- `-` — remove a key

```sql
SELECT payload ?| ARRAY['type','missing'];
SELECT payload - 'n';
```

`@>` is the usual filter for a GIN `jsonb` index (topic 8, topic 14).

Do not use `->>` in a hot `WHERE` if you can write `@>` or a stored column. A cast of `->>` can block a simple index unless you create an expression index.

PostgreSQL 17 adds more SQL/JSON features, such as `JSON_TABLE`. This handbook uses the operators that you need first on 16 and 17.

### Questions

#### Theoretical questions

1. What is the difference between `->` and `->>`?
2. What does `@>` test?
3. What does `?` test?
4. What type is `(payload ->> 'n')` before you cast?
5. Are JSON array indexes 0-based or 1-based?

#### Easy practical tasks

1. Insert one `jsonb` object. Select `-> 'type'` and `->> 'type'`. Use `pg_typeof` on both.
2. Filter with `@>` and a small object.
3. Filter with `?` and a key name.
4. Read the first array element with `->`.

#### Medium practical tasks

1. Cast a numeric JSON field to `int` and add 1.
2. Use `?|` and `?&` with a key array.
3. Update `payload` with `jsonb_set` (look up the function). Show before and after.

#### Advanced practical tasks

1. Create an expression index on `(payload->>'type')` and a GIN index on `payload`. Write two `WHERE` clauses, one for each index (topic 8).
2. Read the `jsonb` operator table in the docs. Add two operators that this section listed only by name. Write an example each.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. When do you use `COALESCE` instead of `||` for display text?
2. When do you use `LIKE`, when POSIX `~`, and when `tsvector`?
3. Why does `GREATEST` skip `NULL` in PostgreSQL, and why must you check other databases?
4. How do `->` and `to_tsvector` both turn a document into something you can search, in different ways?
5. Why is `concat_ws` safer than `||` when a middle column is `NULL`?

#### Easy practical tasks

1. On one `posts` or `events` row, run a string function, a date function, `COALESCE`, a `LIKE`, and a JSON `->>`.
2. Write a cheat sheet: string, date, math, null helpers, three pattern systems, FTS types, four JSON operators.
3. Build a `generate_series` of dates and format them with `to_char`.
4. Search `jsonb` with `@>` and search text with `~*`.

#### Medium practical tasks

1. Write one query that ranks posts with `ts_rank` and another that filters `payload @>`.
2. Replace `IFNULL`-style logic from another product with `COALESCE` and `NULLIF`. Show both.
3. Document team rules: which pattern tool is default, when FTS is required, when `jsonb` operators are required.

#### Advanced practical tasks

1. Combine `regexp_match` and `jsonb` to extract a field and validate it with `~`.
2. Read the function index in the official docs. Map ten functions you used in this topic to their chapter titles.
