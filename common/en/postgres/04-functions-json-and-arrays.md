# 4. Functions, JSON, and Arrays

## Description

This topic shows common functions, JSON types, and arrays in PostgreSQL 16 and PostgreSQL 17. You learn string, date, and math functions. You learn `COALESCE`, `NULLIF`, `GREATEST`, and `LEAST`. You also learn pattern matching, a short full-text search preview, `json` versus `jsonb`, GIN indexes on `jsonb`, and `unnest`.

Complete topic 3 first. You can write joins and `WITH`. Complete this topic before you study indexes and query plans.

Use one term for each concept. A function returns a value from arguments. An operator is a symbol such as `||` or `@>`. `json` is the text JSON type. `jsonb` is the binary JSON type. An array is a typed list (`text[]`, `int[]`). `unnest` turns an array into a set of rows. A GIN index on `jsonb` supports containment and key-exists operators. Prefer `jsonb` for JSON that you search.

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
SELECT random();
```

`round` on `numeric` takes a scale. `round` on `float8` uses a different form. Prefer `numeric` when you need a set scale.

`random()` returns a `float8` in `[0, 1)`. Do not use it for security tokens. Use `gen_random_uuid()` or `pgcrypto` (topic 9) for secrets.

Call functions with argument names when the docs show them. Prefer documented names in application SQL. Prefer operators when they are the usual form (`||` for text).

### Questions

#### Theoretical questions

1. What does `'a' || NULL` return?
2. When do you use `concat` instead of `||`?
3. What does `date_trunc('day', timestamptz)` do?
4. What type does `extract` return?
5. Why is `random()` not a secret generator?

#### Easy practical tasks

1. Run `length`, `lower`, `trim`, and `split_part` on one string.
2. Compute `DATE '2026-09-13' + INTERVAL '7 days'`.
3. Run `round(12.56, 1)` on `numeric` and on `float8`.
4. Concatenate first name and last name with a space. Include a `NULL` last name and compare `||` with `concat_ws`.

#### Medium practical tasks

1. Group events by `date_trunc('month', created_at)`. Count rows per month.
2. Use `substring` and `position` to take the domain from an email.
3. Compare `age(d)` with `now() - d` for one date. Write the two types.

#### Advanced practical tasks

1. Read "String Functions" in the 17 docs. List five functions that this section omitted. Write one example each.
2. Write a `SELECT` that converts `timestamptz` to a labeled local date in two time zones.

---

## `COALESCE`, `NULLIF`, `GREATEST`, `LEAST`

These functions pick a value from a list. They are not aggregates. They work on a fixed list of arguments.

`COALESCE(a, b, ...)` returns the first argument that is not `NULL`.

```sql
SELECT COALESCE(NULL, NULL, 'nail');
SELECT COALESCE(nickname, display_name, 'guest') FROM users;
```

`NULLIF(a, b)` returns `NULL` when `a` equals `b`. Otherwise it returns `a`.

```sql
SELECT NULLIF(qty, 0);
SELECT NULLIF(trim(name), '');
```

A common pair: turn empty string into `NULL`, then fill a default:

```sql
SELECT COALESCE(NULLIF(trim(name), ''), 'unknown');
```

`GREATEST(a, b, ...)` returns the largest non-null value. `LEAST` returns the smallest non-null value. If every argument is `NULL`, the result is `NULL`.

```sql
SELECT GREATEST(1, 5, 3);
SELECT LEAST(created_at, updated_at);
SELECT GREATEST(1, NULL, 2);  -- 2
```

`GREATEST` and `LEAST` compare values of a common type. Mixed types follow the usual type rules.

Do not use `COALESCE` to hide data errors in every query. Fix the schema with `NOT NULL` when the value must exist. Do not write `GREATEST` when you need `MAX` over a group. `MAX` is an aggregate (topic 3 `GROUP BY` / later window topic).

### Questions

#### Theoretical questions

1. What does `COALESCE` return?
2. What does `NULLIF(a, a)` return?
3. How does `GREATEST` treat `NULL` arguments?
4. When do you pair `NULLIF` with `COALESCE`?
5. Why is `GREATEST` not a replacement for `MAX`?

#### Easy practical tasks

1. Select `COALESCE(NULL, 2, 3)`.
2. Select `NULLIF(0, 0)` and `NULLIF(1, 0)`.
3. Select `GREATEST(10, 2, 8)` and `LEAST(10, 2, 8)`.
4. Replace empty `text` with `'n/a'` using `COALESCE` and `NULLIF`.

#### Medium practical tasks

1. Build `users` with nullable `nickname`. Select a display label with `COALESCE`.
2. Use `LEAST(now(), expires_at)` to cap a timestamp.
3. Show `GREATEST` of three `numeric` prices including one `NULL`.

#### Advanced practical tasks

1. Read the docs for `COALESCE` and `GREATEST`. Write the official `NULL` rules in four sentences.
2. Write a query that uses `NULLIF` to avoid division by zero: `amount / NULLIF(qty, 0)`.

---

## Pattern matching and full-text search preview

PostgreSQL has several pattern tools. Pick one on purpose.

`LIKE` / `ILIKE` — SQL patterns (`%`, `_`). Topic 3 covered them. They are enough for prefix and simple contains. `LIKE '%foo%'` cannot use a normal B-tree (topic 5).

`SIMILAR TO` — SQL regular expressions with `LIKE` quoting. Few teams use it. Prefer `LIKE` or POSIX regex.

POSIX regular expressions:

```sql
SELECT name ~ '^[A-Za-z]+$';
SELECT name ~* 'nail';          -- case-insensitive
SELECT substring(name FROM '[0-9]+');
```

`~` is case-sensitive match. `~*` is case-insensitive. `!~` and `!~*` are the negations.

Full-text search (FTS) is a separate system. You store or build a `tsvector`. You query with a `tsquery`. `@@` tests a match.

```sql
SELECT to_tsvector('english', 'The nails are cheap');
SELECT to_tsvector('english', 'The nails are cheap')
       @@ to_tsquery('english', 'nail');
```

`to_tsvector` reduces words (stemming) and drops stop words. `nail` can match `nails`. A GIN index on `tsvector` supports `@@` (topic 5).

Useful FTS functions for a first look:

- `to_tsvector(config, text)`
- `to_tsquery(config, text)` — operators `&`, `|`, `!`
- `plainto_tsquery` — turns plain text into an AND query
- `ts_headline` — highlights matches
- `ts_rank` — ranks matches

```sql
SELECT ts_headline('english', body, query), ts_rank(vec, query)
FROM ...;
```

This section is a preview. FTS has dictionaries, weights, and ranking. Do not treat `ILIKE '%word%'` as search for a large text corpus. Do not mix `~` and `@@` without a plan. Do not skip the text search configuration (`english`, `simple`).

`pg_trgm` (topic 9) adds trigram similarity and can index `ILIKE '%foo%'`. Use it when you need fuzzy text without full FTS.

### Questions

#### Theoretical questions

1. Which operator tests a POSIX regex?
2. What does `to_tsvector` do to words?
3. What does `@@` test?
4. Why is `LIKE '%foo%'` a poor search on a large table?
5. What is a text search configuration?

#### Easy practical tasks

1. Run `name ILIKE '%na%'` on `items`.
2. Run `'PostgreSQL' ~* 'sql'`.
3. Run the `to_tsvector` / `to_tsquery` example. Write true or false.
4. Compare `plainto_tsquery('english', 'cheap nails')` with `to_tsquery('english', 'cheap & nail')`.

#### Medium practical tasks

1. Add a `body text` column. Find rows with `@@ plainto_tsquery('english', '...')`.
2. Show `ts_headline` for one matching row.
3. Write one `~` pattern that checks an email shape (simple, not complete RFC).

#### Advanced practical tasks

1. Read "Full Text Search" introduction in the 17 docs. Write six sentences: vector, query, `@@`, config, rank, index class.
2. Compare `ILIKE`, `~`, and `@@` on the same 1000-row text table with `EXPLAIN`. Write which tool you pick for prefix, regex, and language search.

---

## `json` vs `jsonb` and JSON operators

`json` stores the input text. It keeps key order. It keeps duplicate keys. It keeps insignificant whitespace. Each operation must parse the text.

`jsonb` stores a parsed tree. It drops insignificant whitespace. If a key repeats, the last key wins. Key order is not a data contract. Operators and indexes work on `jsonb`.

```sql
SELECT '{"b": 1, "a": 2}'::json;
SELECT '{"b": 1, "a": 2}'::jsonb;
SELECT pg_typeof('{"a":1}'::jsonb);
```

Prefer `jsonb` for columns that you filter, index, or update in place. Prefer `json` only when you must store the exact text that a client sent.

| Operator | Result |
| --- | --- |
| `->` | field or element; type stays `json` / `jsonb` |
| `->>` | field or element as `text` |
| `#>` / `#>>` | path as JSON or as text |
| `@>` | `jsonb` containment (left contains right) |
| `<@` | left is contained in right |
| `?` | top-level key exists |
| `?\|` / `?&` | any / all of those keys exist |
| `\|\|` | concatenate `jsonb` |

```sql
CREATE TABLE events (
    id      bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    payload jsonb NOT NULL
);

INSERT INTO events (payload)
VALUES ('{"type": "click", "n": 3, "tags": ["a","b"]}');

SELECT payload ->> 'type' AS type,
       payload -> 'n'     AS n_jsonb,
       payload @> '{"type":"click"}' AS is_click,
       payload ? 'tags' AS has_tags
FROM events;
```

`->>` returns `text`. Cast when you need a number: `(payload ->> 'n')::int`.

PostgreSQL 16 adds more SQL/JSON predicates (`IS JSON` and related tests). PostgreSQL 17 adds SQL/JSON query helpers such as `JSON_TABLE` for shaping JSON into rows. You can still do daily work with `->`, `->>`, `@>`, and `jsonb_set`.

Update a key:

```sql
UPDATE events
SET payload = jsonb_set(payload, '{n}', '4', true)
WHERE id = 1;
```

`jsonb_set` takes a path array, a new `jsonb` value, and a create-missing flag.

Do not store every column as JSON. Use tables for fields that you filter, join, and constrain. Use `jsonb` for optional or changing attributes.

### Questions

#### Theoretical questions

1. What does `json` keep that `jsonb` drops?
2. What type does `->>` return?
3. What does `@>` test?
4. When do you prefer `json` over `jsonb`?
5. What does `jsonb_set` change?

#### Easy practical tasks

1. Cast the same object to `json` and `jsonb`. Select both.
2. Insert one `events` row. Select `type` with `->>`.
3. Test `@>` and `?` on that row.
4. Select a nested value with `#>> '{tags,0}'`.

#### Medium practical tasks

1. Update `n` with `jsonb_set`. Select the new payload.
2. Filter `WHERE payload @> '{"type":"click"}'`.
3. Use `jsonb_each_text(payload)` to list keys and values as rows.

#### Advanced practical tasks

1. Read PostgreSQL 17 `JSON_TABLE` (or 16 `IS JSON` if you use 16). Write one example that this section did not show.
2. Design a table with typed columns plus a `jsonb` `attrs`. Write which keys must not live only in JSON.

---

## Indexing `jsonb` with GIN

A GIN index on `jsonb` supports containment and key-exists operators. It does not replace a B-tree on a typed column.

```sql
CREATE INDEX events_payload_gin ON events USING GIN (payload);
```

The default `jsonb` GIN operator class is `jsonb_ops`. It supports `@>`, `?`, `?|`, `?&`, and related searches on keys and values.

`jsonb_path_ops` is a smaller index. It supports `@>` well. It does not support `?` key-exists in the same way. Use it when you only query with containment.

```sql
CREATE INDEX events_payload_path
    ON events USING GIN (payload jsonb_path_ops);
```

Expression indexes extract one field:

```sql
CREATE INDEX events_type_idx ON events ((payload ->> 'type'));
```

That B-tree helps `WHERE payload ->> 'type' = 'click'`. It does not help `@>`.

Rules:

- GIN is large and slow to build and to update. Use it when you query the document.
- A typed column plus a B-tree is cheaper when the key is stable.
- `ANALYZE` the table after a large load (topic 5).
- `EXPLAIN` must show a bitmap index scan or similar before you trust the index.

PostgreSQL 16 and 17 keep these classes. Partial GIN indexes are valid: `CREATE INDEX ... ON events USING GIN (payload) WHERE payload ? 'type';`

Do not create both `jsonb_ops` and `jsonb_path_ops` on the same column without a measured need. Do not index a `json` (text) column with GIN for these operators. Cast or store `jsonb`.

### Questions

#### Theoretical questions

1. Which index method do you use for `jsonb` containment?
2. What does `jsonb_ops` support that `jsonb_path_ops` may not?
3. When is a B-tree on `(payload ->> 'type')` the better index?
4. Why can a GIN index slow `INSERT`?
5. Can you build a partial GIN index?

#### Easy practical tasks

1. Create `events` and a GIN index on `payload`. Run `\d events`.
2. `EXPLAIN` a query with `@> '{"type":"click"}'`.
3. Create a B-tree on `(payload ->> 'type')`. `EXPLAIN` an equality filter on that field.
4. Write four sentences: GIN, `@>`, `?`, expression B-tree.

#### Medium practical tasks

1. Load a few thousand `jsonb` rows. Compare `EXPLAIN ANALYZE` with and without the GIN index.
2. Create `jsonb_path_ops`. Try a `?` query. Write whether the planner uses that index.
3. Build a partial GIN index for one document type.

#### Advanced practical tasks

1. Read "jsonb Indexing" in the 17 docs. Write a table: operator class, operators, size trade-off.
2. Measure index size with `pg_relation_size` for `jsonb_ops` versus `jsonb_path_ops` on the same data.

---

## Arrays and `unnest`

An array is a typed list. Write the type as `int[]`, `text[]`, or `jsonb[]`.

```sql
CREATE TABLE posts (
    id   bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tags text[] NOT NULL DEFAULT '{}'
);

INSERT INTO posts (tags) VALUES (ARRAY['sql', 'postgres']);
INSERT INTO posts (tags) VALUES ('{gin,json}');
```

`{}` is an empty array, not `NULL`. `{a,b}` is the text input form.

Operators and functions:

```sql
SELECT tags[1] FROM posts;                 -- 1-based index
SELECT tags && ARRAY['sql'];               -- overlap
SELECT tags @> ARRAY['sql'];               -- contains
SELECT 'sql' = ANY (tags);
SELECT unnest(tags) AS tag FROM posts;
SELECT cardinality(tags);
SELECT array_append(tags, 'btree');
SELECT array_agg(name) FROM items;
```

`unnest` returns a set of rows. Use it in `FROM` or in the select list. `WITH ORDINALITY` adds a position:

```sql
SELECT p.id, u.tag, u.ord
FROM posts AS p
CROSS JOIN LATERAL unnest(p.tags) WITH ORDINALITY AS u(tag, ord);
```

`array_agg` builds an array from a group. `string_to_array` and `string_agg` move between text and arrays.

A GIN index on `text[]` supports `@>`, `<@`, `&&`, and `=`:

```sql
CREATE INDEX posts_tags_gin ON posts USING GIN (tags);
```

Array limits:

- All elements have the same type.
- Arrays can be multi-dimensional. Prefer one dimension for application data.
- A separate `post_tags` table is easier to constrain and to join. Use an array for a small, simple list.

Do not store a large relational model in arrays. Do not use `tags[0]`. The first element is `tags[1]`.

### Questions

#### Theoretical questions

1. What is the first index of an array in PostgreSQL?
2. What does `unnest` return?
3. What does `&&` mean for arrays?
4. What does `array_agg` build?
5. When do you prefer a child table over `text[]`?

#### Easy practical tasks

1. Create `posts` with `tags text[]`. Insert two rows.
2. Select `tags[1]` and `cardinality(tags)`.
3. Select `unnest(tags)` for one post.
4. Test `@>` and `&&` with `ARRAY['sql']`.

#### Medium practical tasks

1. Use `WITH ORDINALITY` to list tags with positions.
2. `UPDATE` with `array_append`. Select the new array.
3. Create a GIN index on `tags`. `EXPLAIN` a `@>` query.

#### Advanced practical tasks

1. Compare `posts.tags text[]` with a `post_tags(post_id, tag)` table. Write six sentences: query, unique tag, index, FK.
2. Read "Array Functions" in the 17 docs. Write examples for `array_remove` and `unnest` of two arrays in parallel (`unnest(a, b)`).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. When do you keep a field as a typed column instead of a `jsonb` key?
2. How do `COALESCE` and `->>` work together when a key is missing?
3. Which index class do you pick for `jsonb` `@>` versus for `payload ->> 'type' = ...`?
4. Why is full-text `@@` a different tool from `ILIKE`?
5. What goes wrong if you treat array index `0` as the first element?

#### Easy practical tasks

1. Build `events(payload jsonb)` and `posts(tags text[])`. Insert one row each. Select one JSON field and one unnested tag.
2. Use `COALESCE(payload ->> 'type', 'unknown')`.
3. Create a GIN index on `payload` and a GIN index on `tags`.
4. Run `date_trunc('day', now())` and `GREATEST(1, 3, 2)`.

#### Medium practical tasks

1. Write one query that unnests `tags` and filters posts whose `payload` (if you add it) or a second `jsonb` column contains a key.
2. Find rows with FTS `@@` and rows with `ILIKE`. Write two SQL statements on the same `body` column.
3. Update JSON with `jsonb_set` and an array with `array_append` in one transaction.

#### Advanced practical tasks

1. Load sample documents and tag arrays. Prove with `EXPLAIN ANALYZE` that GIN is used for `@>` on both types.
2. Write a one-page policy: when to use `text`, `text[]`, `jsonb`, and `tsvector`. Give one example each.
