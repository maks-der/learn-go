# 14. JSON, Arrays, and Documents

## Description

This topic shows JSON and arrays in PostgreSQL 16 and PostgreSQL 17. You learn `json` versus `jsonb`, GIN indexes on `jsonb`, arrays and `unnest`, when JSON belongs in PostgreSQL, and limited check constraints on JSON shape.

Complete topic 13 first. You can create extensions and functions. Complete this topic before you partition large tables.

Use one term for each concept. `json` is the text JSON type. `jsonb` is the binary JSON type. A document is one JSON value in a column. An array is a typed list (`text[]`, `int[]`, `jsonb[]`). `unnest` turns an array into a set of rows. A GIN index on `jsonb` supports containment and key-exists operators. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## `json` vs `jsonb`

Topic 4 introduced the two types. This section is the working rule.

`json` stores the input text. It keeps key order. It keeps duplicate keys. It keeps insignificant whitespace. Each operation must parse the text.

`jsonb` stores a parsed tree. It drops insignificant whitespace. If a key repeats, the last key wins. Key order is not a data contract. Operators and indexes work on `jsonb`.

```sql
SELECT '{"b": 1, "a": 2}'::json;
SELECT '{"b": 1, "a": 2}'::jsonb;
SELECT pg_typeof('{"a":1}'::jsonb);
```

Prefer `jsonb` for columns that you filter, index, or update in place. Prefer `json` only when you must store the exact text that a client sent.

Topic 7 listed operators. Review:

| Operator | Result |
| --- | --- |
| `->` | field or element; type stays `json` / `jsonb` |
| `->>` | field or element as `text` |
| `@>` | `jsonb` containment |
| `?` | `jsonb` top-level key exists |

```sql
CREATE TABLE shop.events (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    payload jsonb NOT NULL
);

INSERT INTO shop.events (payload)
VALUES ('{"type": "click", "n": 3}');

SELECT payload ->> 'type' AS type,
       payload -> 'n'     AS n_jsonb,
       payload @> '{"type":"click"}' AS is_click
FROM shop.events;
```

PostgreSQL 16 adds more SQL/JSON predicates (`IS JSON` and related tests). PostgreSQL 17 adds SQL/JSON query helpers such as `JSON_TABLE` for shaping JSON into rows. You can still do daily work with `->`, `->>`, `@>`, and `jsonb_set`.

Update a key:

```sql
UPDATE shop.events
SET payload = jsonb_set(payload, '{n}', '4', true)
WHERE id = 1;
```

Do not store JSON as `text` if you need the parser to reject invalid JSON. Do not use `json` because the name is shorter. Do not compare `jsonb` with `=` when you meant containment (`@>`).

Official chapter: [https://www.postgresql.org/docs/17/datatype-json.html](https://www.postgresql.org/docs/17/datatype-json.html).

### Questions

#### Theoretical questions

1. What does `json` keep that `jsonb` drops?
2. Which type do you index for search?
3. What is the difference between `->` and `->>`?
4. What does `@>` test?
5. When is `json` the better column type?

#### Easy practical tasks

1. Cast the same object to `json` and `jsonb`. Select both.
2. Insert one `jsonb` row. Select `-> 'type'` and `->> 'type'`. Run `pg_typeof` on both.
3. Run `SELECT '{"a":1}'::jsonb @> '{"a":1}';`.
4. Open the 17 JSON type page. Write one 16-or-17 feature name that this section mentioned.

#### Medium practical tasks

1. Use `jsonb_set` to add a key and to change a key. Select before and after.
2. Compare `'{"b":1,"a":2}'::json` and the same literal as `jsonb`. Write what changes.
3. Run `EXPLAIN` on `WHERE payload @> '{"type":"click"}'` without an index (next section adds GIN).

#### Advanced practical tasks

1. On PostgreSQL 17, try `JSON_TABLE` from the official docs on a small `jsonb` value. On 16, try `IS JSON`. Save the query that your major version accepts.
2. Read duplicate-key rules for `json` and `jsonb` in the docs. Write two examples that show the difference.

---

## Indexing `jsonb` with GIN

A B-tree on `(payload->>'type')` helps equality on one extracted key. A GIN index on the `jsonb` column helps containment and key-exists on many keys.

```sql
CREATE INDEX events_payload_gin ON shop.events USING gin (payload);
```

The default operator class is `jsonb_ops`. It supports `@>`, `?`, `?|`, `?&`, and related searches. It indexes keys and values. It is larger.

`jsonb_path_ops` supports `@>` only. It is smaller and can be faster for containment. It does not support `?`.

```sql
CREATE INDEX events_payload_path_gin
    ON shop.events USING gin (payload jsonb_path_ops);
```

Check the plan:

```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT id
FROM shop.events
WHERE payload @> '{"type":"click"}';
```

You want a bitmap index scan or an index scan that names the GIN index. A sequential scan on a large table means the index does not match the predicate. `@>` matches GIN. `payload->>'type' = 'click'` matches a B-tree on that expression, not the default GIN containment path.

Expression GIN is also valid:

```sql
CREATE INDEX events_tags_gin ON shop.events USING gin ((payload -> 'tags'));
```

Topic 8 covered GIN as a class. This section is the `jsonb` pattern. `CREATE INDEX CONCURRENTLY` still applies on a busy table.

Do not create both `jsonb_ops` and `jsonb_path_ops` on the same column without a measurement. Do not expect GIN to speed `ORDER BY payload`. Do not index huge documents if you only filter one small key; a B-tree on the expression can be enough.

Official operator classes: [https://www.postgresql.org/docs/17/datatype-json.html#JSON-INDEXING](https://www.postgresql.org/docs/17/datatype-json.html#JSON-INDEXING).

### Questions

#### Theoretical questions

1. Which operator class is the default for `USING gin (payload)`?
2. What operators does `jsonb_path_ops` support?
3. Why is `jsonb_path_ops` often smaller?
4. Which index type fits `payload->>'type' = 'click'`?
5. Does GIN help `ORDER BY payload`?

#### Easy practical tasks

1. Create a GIN index on `payload`. Run `\d shop.events`.
2. Run `EXPLAIN` for `@>` on a matching document.
3. Create a B-tree on `(payload->>'type')`. Run `EXPLAIN` for `payload->>'type' = 'click'`.
4. Write four sentences: `jsonb_ops`, `jsonb_path_ops`, B-tree extract, `@>`.

#### Medium practical tasks

1. Load enough rows that a sequential scan is costly. Compare `EXPLAIN ANALYZE` with and without GIN for `@>`.
2. Create `jsonb_path_ops`. Run `EXPLAIN` for `payload ? 'type'`. Write whether that index can serve `?`.
3. Compare `pg_relation_size` of `jsonb_ops` and `jsonb_path_ops` on the same data.

#### Advanced practical tasks

1. Build two queries: one containment, one extracted equality. Give each its best index. Prove with `EXPLAIN`.
2. Read the JSON indexing docs for 16 and 17. Write one extra operator or path query that GIN can support. Cite the page.

---

## Arrays and `unnest`

An array is a single column that holds a list of values of one type. Syntax uses `{` in literals and `[]` in constructors.

```sql
CREATE TABLE shop.posts (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tags text[] NOT NULL DEFAULT '{}'
);

INSERT INTO shop.posts (tags) VALUES (ARRAY['sql', 'gin']);
INSERT INTO shop.posts (tags) VALUES ('{sql,json}');
```

Operators and functions:

- `tags @> ARRAY['sql']` — contains all listed elements
- `tags && ARRAY['gin', 'fts']` — overlaps
- `tags[1]` — first element (arrays are 1-based)
- `array_length(tags, 1)` — length of the first dimension
- `unnest(tags)` — one row per element

```sql
SELECT p.id, u.tag
FROM shop.posts AS p
CROSS JOIN LATERAL unnest(p.tags) AS u(tag);
```

Topic 6 covered `LATERAL`. `unnest` is a common lateral partner.

A GIN index on `text[]` supports `@>` and `&&`:

```sql
CREATE INDEX posts_tags_gin ON shop.posts USING gin (tags);
```

You can store `jsonb[]`, but a `jsonb` array inside one document (`{"tags":[...]}`) is often clearer. Do not mix both styles without a rule.

An array is not a child table. You cannot put a foreign key on one element. You cannot attach a unique constraint to one element without extra work. Topic 4 said to use arrays for small sets. This section adds `unnest` and GIN.

Do not replace `orders` and `order_lines` with `line_ids bigint[]` when you need line totals, constraints, and joins. Do not use `0`-based indexing in SQL; the first element is `[1]`.

Official chapter: [https://www.postgresql.org/docs/17/arrays.html](https://www.postgresql.org/docs/17/arrays.html).

### Questions

#### Theoretical questions

1. What is the index of the first array element in PostgreSQL?
2. What does `unnest` return?
3. What does `tags @> ARRAY['sql']` mean?
4. What does `tags && ARRAY['a','b']` mean?
5. Why is an array a poor substitute for a child table?

#### Easy practical tasks

1. Create `shop.posts` and insert two tag arrays.
2. Select `tags[1]` and `array_length(tags, 1)`.
3. `SELECT unnest(tags) FROM shop.posts;`
4. Write a `LATERAL unnest` query that returns `id` and `tag`.

#### Medium practical tasks

1. Create a GIN index on `tags`. `EXPLAIN` a `@>` filter.
2. Find posts that overlap `ARRAY['json','gin']`.
3. Compare `unnest` in the `SELECT` list with `CROSS JOIN LATERAL unnest`. Write the row counts.

#### Advanced practical tasks

1. Normalize `tags` into `post_tags(post_id, tag)` with `INSERT ... SELECT id, unnest(tags)`. Compare a join query with the array query.
2. Read array functions in the 17 docs. Add `array_agg` and `cardinality`. Write one example each on your `posts` table.

---

## When JSON belongs in PostgreSQL vs a document DB

PostgreSQL can store documents. A document-only database can store documents. The choice is a model choice, not a fashion choice.

JSON belongs in PostgreSQL when:

- most data already has tables, keys, and transactions
- you need joins from a document to a relational row in one transaction
- the document is a payload, settings bag, or sparse attribute set
- you still want PostgreSQL backup, roles, and constraints on the typed columns

JSON does not belong as the only model when:

- every query is a deep path into a changing document and you never join
- you need a document API and a managed document product already exists
- you store large opaque blobs that you never filter (use `bytea` or object storage and keep a key)

Rules that keep a hybrid design safe:

- typed columns for values that you filter, sort, join, or constrain
- `jsonb` for values that vary by row or change shape
- a GIN or expression index only for paths that you measure
- a check constraint when the shape is small and stable (next section)

PostgreSQL is not "MongoDB with SQL" when you skip keys. Topic 1 said tables stay the primary model. A document database is not "wrong". It is a different product with a different operations model.

Do not put the whole application state in one `jsonb` column named `data`. Do not duplicate the same attribute as a column and as JSON without a write rule. Do not pick JSON to avoid a migration. Topic 21 covers migrations.

### Questions

#### Theoretical questions

1. When does `jsonb` belong next to typed columns?
2. When is a document-only store a better fit?
3. Why keep join keys as typed columns?
4. What is the risk of one `data jsonb` column for the whole row?
5. Does a GIN index make a `jsonb` column a complete substitute for a table?

#### Easy practical tasks

1. Write five sentences: three cases for `jsonb` in PostgreSQL, two cases against a JSON-only schema.
2. Draw a table `orders` with typed `id`, `customer_id`, `total`, and `payload jsonb`. Label why each column is typed or JSON.
3. List two document-database products and two PostgreSQL features they do not share (from this path: roles, `GRANT`, `VACUUM`, or SQL joins).
4. Open the official JSON docs intro. Write one sentence that the docs use for JSON in PostgreSQL.

#### Medium practical tasks

1. Take a sample document with `customer_id`, `sku`, and `qty`. Split it into tables plus a small `jsonb` extras column. Write the `CREATE TABLE` statements.
2. Write a query that joins `events.payload->>'sku'` to `products.sku`. Then rewrite it with a typed `sku` column. Compare `EXPLAIN`.
3. Interview yourself: list ten attributes of an order. Mark each as column, JSON, or child table. Give one reason each.

#### Advanced practical tasks

1. Write a one-page decision record: when your team allows a new `jsonb` column. Include index, backup, and constraint checks.
2. Compare PostgreSQL `jsonb` with one document database using only official docs. Cover transactions, joins, and indexes in six short sentences.

---

## Check constraints on JSON shape (limited)

A `jsonb` column accepts any valid JSON unless you add a constraint. PostgreSQL 16 and 17 do not ship a full JSON Schema validator as a core constraint. You can test types, keys, and small shapes with SQL.

```sql
ALTER TABLE shop.events
    ADD CONSTRAINT events_payload_is_object
    CHECK (jsonb_typeof(payload) = 'object');

ALTER TABLE shop.events
    ADD CONSTRAINT events_payload_has_type
    CHECK (payload ? 'type');

ALTER TABLE shop.events
    ADD CONSTRAINT events_type_is_text
    CHECK (jsonb_typeof(payload -> 'type') = 'string');
```

`IS JSON` (PostgreSQL 16 and 17) tests text or JSON syntax in expressions. A `jsonb` column already rejected invalid JSON on input. The useful checks are shape and types of fields.

Limits:

- each `CHECK` is a boolean SQL expression
- deep optional trees are hard to express
- error messages are not field-level schema errors
- a constraint that calls a `STABLE` function must stay cheap on every write

You can wrap checks in an `IMMUTABLE` SQL function and call it from `CHECK`. The function must stay `IMMUTABLE` for the constraint to be valid.

For a large schema, validate in the application and keep a few database checks for invariants that must never fail (type is present, id is a number).

Do not copy a full JSON Schema into fifty `CHECK` clauses. Do not use a `VOLATILE` function in `CHECK`. Do not assume that `@>` proves a complete schema; it only proves containment of the sample that you wrote.

### Questions

#### Theoretical questions

1. Does a `jsonb` column reject unknown keys by default?
2. What does `jsonb_typeof` return for an object?
3. Why is a full JSON Schema a poor fit for many `CHECK` clauses?
4. What volatility must a function have if you call it from `CHECK`?
5. What does `payload ? 'type'` not prove about the type of the value?

#### Easy practical tasks

1. Add `CHECK (jsonb_typeof(payload) = 'object')`. Insert an array. Record the error.
2. Add a `? 'type'` check. Insert `{}`. Record the error.
3. Insert a valid object that passes both checks.
4. Write four sentences: valid JSON, shape, `CHECK`, application validation.

#### Medium practical tasks

1. Add a check that `payload ->> 'n'` is a digit string or that `jsonb_typeof(payload -> 'n') = 'number'`. Test a bad value.
2. Write an `IMMUTABLE` SQL function that returns boolean for a small shape. Use it in `CHECK`.
3. Compare `payload @> '{"type":"click"}'` as a check versus as a query filter. Write which role each form has.

#### Advanced practical tasks

1. Document five invariants for an event document. Implement three as `CHECK`. Leave two in the application. Explain the split.
2. Read JSON functions in the 17 docs (`jsonb_typeof`, `jsonb_exists`, `IS JSON`). Write one constraint that this section did not show.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `jsonb`, GIN, arrays, and `CHECK` work together on one events table?
2. Why can a B-tree on `(payload->>'type')` and a GIN on `payload` both be correct?
3. When do you `unnest` an array instead of storing a JSON array in `jsonb`?
4. How does the hybrid-model rule change your index plan?
5. A teammate wants to drop all typed columns and keep one `jsonb` document. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. Create `events(payload jsonb)` and `posts(tags text[])`. Insert two rows each. Run one `@>` query on each table.
2. Write a cheat sheet: `json`, `jsonb`, GIN classes, `unnest`, hybrid rule, `CHECK` limits.
3. Create one GIN index and one array GIN index. Show `\d` for both tables.
4. Run `SELECT jsonb_typeof('[]'::jsonb), jsonb_typeof('{}'::jsonb);`.

#### Medium practical tasks

1. Load 10000 event documents with two `type` values. Add GIN. Compare `EXPLAIN ANALYZE` for `@>` before and after.
2. Build a report: unnest tags, join to a small `tag_dim` table. Write the SQL.
3. Add two shape checks to `events`. Write a short test list: two inserts that fail, two that work.

#### Advanced practical tasks

1. Design a schema for "order plus flexible extras". Use typed money, a child `lines` table, and `extras jsonb` with a GIN index and two checks. Load sample data and three queries.
2. Map this topic to official chapters: JSON types, JSON functions, arrays, `CREATE INDEX`. Add one PostgreSQL 17 SQL/JSON fact that this topic did not include.
