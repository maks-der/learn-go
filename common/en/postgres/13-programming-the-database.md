# 13. Programming the Database

## Description

This topic shows how you put logic in PostgreSQL 16 and PostgreSQL 17. You learn SQL functions, PL/pgSQL, triggers, generated columns, and extensions. You also learn common extensions: `pgcrypto`, `uuid-ossp`, `citext`, `pg_trgm`, and `pg_stat_statements`.

Complete topic 12 first. You can grant `EXECUTE` and you know `SECURITY DEFINER` risk. Complete this topic before you store JSON documents as a primary model.

Use one term for each concept. A function is a named routine that returns a value or a set. A procedure is a routine that you `CALL` and that can commit. A trigger is a function that the server runs on a table event. A generated column stores a value that the server computes. An extension is a packaged set of objects that you install with `CREATE EXTENSION`. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## SQL functions

A SQL function body is one SQL statement or a list of statements in PostgreSQL 16 and 17. The language name is `sql`. The function runs in the same transaction as the caller.

```sql
CREATE FUNCTION shop.item_label(p_name text, p_qty int)
RETURNS text
LANGUAGE sql
STABLE
AS $$
    SELECT p_name || ' (' || p_qty::text || ')';
$$;

SELECT shop.item_label('nail', 10);
```

Volatility labels help the planner:

- `VOLATILE` — default; can change data or return different values
- `STABLE` — same result inside one scan for the same arguments; can read tables
- `IMMUTABLE` — same result forever for the same arguments; no table reads

Use `IMMUTABLE` only when it is true. A wrong `IMMUTABLE` label can make an index or a generated column store a stale result.

SQL functions can return a table:

```sql
CREATE FUNCTION shop.low_stock(p_min int)
RETURNS TABLE (id bigint, name text, qty int)
LANGUAGE sql
STABLE
AS $$
    SELECT i.id, i.name, i.qty
    FROM shop.items AS i
    WHERE i.qty < p_min;
$$;
```

Prefer SQL functions when the body is a query. The planner can often inline a simple SQL function. Prefer PL/pgSQL when you need variables, loops, or `EXCEPTION`.

`CREATE OR REPLACE FUNCTION` changes the body. Some signature changes need `DROP FUNCTION` first.

Grant `EXECUTE` to the roles that must call the function. Revoke `EXECUTE` from `PUBLIC` when the function is not a public helper.

Do not hide a long business workflow in a nest of SQL functions that nobody can test. Do not mark a function `IMMUTABLE` if it calls `now()` or reads a table.

Official reference: [https://www.postgresql.org/docs/17/xfunc-sql.html](https://www.postgresql.org/docs/17/xfunc-sql.html).

### Questions

#### Theoretical questions

1. What language name do you use for a SQL function?
2. What is the difference between `STABLE` and `IMMUTABLE`?
3. What does `RETURNS TABLE` mean?
4. Why can a simple SQL function be good for the planner?
5. Who may call a function after you revoke `EXECUTE` from `PUBLIC`?

#### Easy practical tasks

1. Create `shop.item_label`. Select it for two argument pairs.
2. Create a SQL function that returns the count of rows in one table.
3. Run `\df shop.*` and `\sf shop.item_label`.
4. Write four sentences: SQL function, volatility, `EXECUTE`, inline.

#### Medium practical tasks

1. Create a `RETURNS TABLE` function. Select from it with `SELECT * FROM shop.low_stock(5);`.
2. Replace the function body with `CREATE OR REPLACE FUNCTION`. Show `\sf` before and after.
3. Compare `EXPLAIN` for a filter in a SQL function versus the same filter written in the outer query.

#### Advanced practical tasks

1. Read "Function Volatility Categories" in the docs. Write one safe example for each of the three labels.
2. Write a SQL function that joins two tables. Grant `EXECUTE` only. Prove that a role without table `SELECT` still fails if the function is `SECURITY INVOKER`.

---

## PL/pgSQL: `$$`, `BEGIN`, exceptions

PL/pgSQL is the procedural language that ships with PostgreSQL. You write blocks with `BEGIN` ... `END`. Dollar-quoting (`$$` or `$tag$`) holds the body so that inner quotes do not break the `CREATE` statement. Topic 3 introduced dollar-quoting.

```sql
CREATE FUNCTION shop.clamp_qty(p_qty int)
RETURNS int
LANGUAGE plpgsql
IMMUTABLE
AS $fn$
BEGIN
    IF p_qty IS NULL THEN
        RETURN 0;
    END IF;
    IF p_qty < 0 THEN
        RETURN 0;
    END IF;
    RETURN p_qty;
END;
$fn$;
```

A block can declare variables:

```sql
CREATE FUNCTION shop.next_label(p_id bigint)
RETURNS text
LANGUAGE plpgsql
STABLE
AS $$
DECLARE
    v_name text;
BEGIN
    SELECT name INTO STRICT v_name
    FROM shop.items
    WHERE id = p_id;
    RETURN v_name;
EXCEPTION
    WHEN NO_DATA_FOUND THEN
        RETURN 'missing';
    WHEN TOO_MANY_ROWS THEN
        RAISE EXCEPTION 'duplicate id %', p_id;
END;
$$;
```

`INTO STRICT` requires exactly one row. `EXCEPTION` handles named conditions. `RAISE` reports an error or a message. `RAISE NOTICE` is for debug. Remove notices from production functions.

PL/pgSQL is `VOLATILE` by default. Set `STABLE` or `IMMUTABLE` when the rules allow it.

Procedures use `CREATE PROCEDURE` and `CALL`. A procedure can `COMMIT` in some cases. Prefer a function that returns a value when you do not need procedure commit behavior.

Do not put an empty `EXCEPTION WHEN OTHERS` that hides every error. Do not loop row by row when one SQL statement can do the work.

Official chapter: [https://www.postgresql.org/docs/17/plpgsql.html](https://www.postgresql.org/docs/17/plpgsql.html).

### Questions

#### Theoretical questions

1. Why does a function body use `$$` or `$tag$`?
2. What does `BEGIN` ... `END` group in PL/pgSQL?
3. What does `INTO STRICT` require?
4. What does `EXCEPTION WHEN NO_DATA_FOUND` handle?
5. How is `CREATE PROCEDURE` different from `CREATE FUNCTION` at a high level?

#### Easy practical tasks

1. Create `shop.clamp_qty`. Select it for `-1`, `NULL`, and `5`.
2. Create a function with `DECLARE` and one variable. Return that variable.
3. Call a function that uses `RAISE NOTICE`. Watch the `psql` message.
4. Write four sentences: dollar-quoting, block, `EXCEPTION`, `RAISE`.

#### Medium practical tasks

1. Write `shop.next_label` with `INTO STRICT` and two handlers. Test missing id and one existing id.
2. Replace `EXCEPTION WHEN OTHERS THEN NULL;` (if you try it) with a named handler. Write why a catch-all is a problem.
3. Convert a three-step PL/pgSQL loop into one SQL statement. Compare `EXPLAIN ANALYZE` if the set is large enough.

#### Advanced practical tasks

1. Write a PL/pgSQL function that updates one row and returns `RETURNING` data through `INTO`. Handle `NO_DATA_FOUND`.
2. Read "Control Structures" and "Error Handling" in the PL/pgSQL docs. Add two condition names that this section did not use. Write one example each.

---

## Triggers: `BEFORE` / `AFTER`, row vs statement

A trigger runs a function when a table event occurs. You create the function first. The function returns `trigger` and uses `TG_` variables. Then you attach it with `CREATE TRIGGER`.

Timing:

- `BEFORE` — runs before the row is written; you can change `NEW`; you can return `NULL` to skip the row
- `AFTER` — runs after the write; use it for side effects that must see the final row
- `INSTEAD OF` — for views

Granularity:

- `FOR EACH ROW` — once per row
- `FOR EACH STATEMENT` — once per statement

```sql
CREATE FUNCTION shop.set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at := clock_timestamp();
    RETURN NEW;
END;
$$;

CREATE TRIGGER items_set_updated_at
    BEFORE UPDATE ON shop.items
    FOR EACH ROW
    EXECUTE FUNCTION shop.set_updated_at();
```

In PostgreSQL 16 and 17 the keyword is `EXECUTE FUNCTION` (older text says `EXECUTE PROCEDURE`). Both forms exist for compatibility. Prefer `EXECUTE FUNCTION`.

`NEW` is the new row version on `INSERT` and `UPDATE`. `OLD` is the old row version on `UPDATE` and `DELETE`. Statement triggers do not have `NEW` and `OLD` row records. They can use transition tables (`REFERENCING NEW TABLE AS ...`) when you need the set of rows.

Use a `BEFORE ROW` trigger to fill columns. Use an `AFTER` trigger to write an audit table. An `AFTER` trigger that updates the same row can loop. Avoid that design.

`updated_at` defaults do not change on `UPDATE`. Topic 5 said to use a trigger or to set the column in the statement. This section is that trigger.

Do not use triggers to replace a foreign key or a check constraint. Do not stack many triggers without a documented order (`WHEN`, names, `pg_trigger.tgenabled`).

Official chapter: [https://www.postgresql.org/docs/17/triggers.html](https://www.postgresql.org/docs/17/triggers.html).

### Questions

#### Theoretical questions

1. What is the difference between `BEFORE` and `AFTER`?
2. What is the difference between a row trigger and a statement trigger?
3. What must a row-level `BEFORE UPDATE` function return to keep the change?
4. What is `NEW` on `INSERT`?
5. Why is a trigger the usual tool for `updated_at`?

#### Easy practical tasks

1. Add `updated_at timestamptz` to a table. Create the `BEFORE UPDATE` trigger. Update a row. Select `updated_at`.
2. List triggers with `\d shop.items`.
3. Insert a row. Confirm that an `UPDATE`-only trigger did not run.
4. Write four sentences: timing, row versus statement, `NEW`, `OLD`.

#### Medium practical tasks

1. Create an `AFTER INSERT` row trigger that writes `id` and `now()` into an audit table. Insert two rows. Select the audit table.
2. In a `BEFORE INSERT` trigger, change `NEW.name`. Insert and show the stored name.
3. Return `NULL` from a `BEFORE INSERT` row trigger. Write what happens to the insert.

#### Advanced practical tasks

1. Write a statement-level `AFTER UPDATE` trigger with `REFERENCING NEW TABLE AS new_rows`. Count the changed rows into a log table.
2. Read trigger firing order in the docs. Create two `BEFORE` triggers on the same event. Document the order that you observe.

---

## Generated columns

A generated column is a column that the server computes from other columns. In PostgreSQL 16 and PostgreSQL 17 the stored form is required. The keyword is `STORED`. The server writes the value on `INSERT` and `UPDATE`. You can index the column. You cannot `UPDATE` the generated column directly.

```sql
CREATE TABLE shop.products (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    price_net numeric(12, 2) NOT NULL,
    tax_rate numeric(4, 3) NOT NULL,
    price_gross numeric(12, 2)
        GENERATED ALWAYS AS (price_net * (1 + tax_rate)) STORED
);
```

PostgreSQL 16 and 17 do not offer virtual generated columns. A later major version can add that feature. Do not write `VIRTUAL` in 16 or 17 scripts.

The expression must be `IMMUTABLE`. It can use other base columns in the same row. It cannot use a subquery. It cannot use a volatile function.

A generated column is not a trigger. The planner sees a real column. `SELECT *` includes it. Topic 8 said that a stored generated column plus a B-tree can replace a long expression index when you also select the value.

```sql
CREATE INDEX products_price_gross_idx ON shop.products (price_gross);
```

Do not copy a generated value into a second writable column. Do not put business rules that need other tables into a generated column. Use a query, a view, or a function.

Official reference: [https://www.postgresql.org/docs/17/ddl-generated-columns.html](https://www.postgresql.org/docs/17/ddl-generated-columns.html).

### Questions

#### Theoretical questions

1. What does `GENERATED ALWAYS AS ... STORED` mean?
2. Can you `UPDATE` only the generated column?
3. Why must the expression be `IMMUTABLE`?
4. Does PostgreSQL 16 or 17 support virtual generated columns?
5. How is a generated column different from a `BEFORE` trigger that fills a column?

#### Easy practical tasks

1. Create `shop.products` with `price_gross`. Insert one row. Select all columns.
2. Update `price_net`. Select `price_gross` again.
3. Try `UPDATE shop.products SET price_gross = 0`. Record the error.
4. Create a B-tree on `price_gross`. Run `\d shop.products`.

#### Medium practical tasks

1. Add a generated `text` column that concatenates two base columns. Filter on it. Compare `EXPLAIN` with an expression index on the same concatenation.
2. Try a generated expression that calls `now()`. Record the error.
3. Write a table: tool, when to use it. Rows: generated column, trigger, view.

#### Advanced practical tasks

1. Migrate a trigger-filled column to a generated column on a copy table. Write the `ALTER TABLE` steps that you used.
2. Read the 16 and 17 generated-column docs. Confirm that both require `STORED`. Write one sentence that you would put in a team style guide.

---

## Extensions: `CREATE EXTENSION`

An extension is a named package. `CREATE EXTENSION` creates its objects and records the name in `pg_extension`. The files live in the PostgreSQL share directory. You do not copy SQL from a blog as a substitute when an extension exists.

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
```

Many extensions need a superuser or a role with `CREATE` on the database plus the package on disk. Cloud vendors often allow a list of extensions and block others.

List what the server can load:

```sql
SELECT name, default_version, installed_version
FROM pg_available_extensions
ORDER BY name;
```

List what this database has:

```sql
SELECT extname, extversion FROM pg_extension ORDER BY 1;
```

`CREATE EXTENSION ... SCHEMA shop` puts objects in that schema when the extension allows it. Some extensions must live in `public` or in a fixed schema.

`DROP EXTENSION` removes the objects. `CASCADE` drops dependent objects. Use `CASCADE` only when you accept that loss.

`ALTER EXTENSION ... UPDATE` moves to a newer extension version that the package provides.

Do not enable every extension "in case you need it". Each extension adds objects, upgrade work, and attack surface. Do not create the same functions by hand when `CREATE EXTENSION` is the supported path.

Official chapter: [https://www.postgresql.org/docs/17/extend-extensions.html](https://www.postgresql.org/docs/17/extend-extensions.html).

### Questions

#### Theoretical questions

1. What does `CREATE EXTENSION` do?
2. Where do you list extensions that you can install?
3. Where do you list extensions that this database already has?
4. What does `DROP EXTENSION ... CASCADE` risk?
5. Why can a cloud service reject `CREATE EXTENSION`?

#### Easy practical tasks

1. Run the `pg_available_extensions` query. Write five names that you recognize.
2. Run `SELECT * FROM pg_extension;`.
3. Read `\dx` in `psql`. Compare it with the query.
4. Open the 17 docs list of contrib modules. Write the URL.

#### Medium practical tasks

1. Create `pg_trgm` or `citext` on a lab database if your role allows it. Show `\dx`.
2. Try `CREATE EXTENSION` of a name that is not available. Record the error.
3. Show the schema of an extension with `\dx+` or `pg_depend` docs. Write where the objects live.

#### Advanced practical tasks

1. Script `CREATE EXTENSION IF NOT EXISTS` for the extensions that your app needs. Run it twice. Confirm that the second run is a no-op.
2. Read how `shared_preload_libraries` relates to `pg_stat_statements` (next section). Write the extra step that `CREATE EXTENSION` alone does not do.

---

## Common extensions: `pgcrypto`, `uuid-ossp` / `pgcrypto` gen, `citext`, `pg_trgm`, `pg_stat_statements`

These contrib extensions appear in many PostgreSQL 16 and 17 apps.

**`gen_random_uuid()`** is built in since PostgreSQL 13. You do not need an extension for version-4 UUIDs in new work:

```sql
SELECT gen_random_uuid();
```

**`uuid-ossp`** still provides other UUID versions (`uuid_generate_v1`, `uuid_generate_v5`, and related functions). Create it only when you need those functions.

**`pgcrypto`** provides `crypt`, `gen_salt`, `digest`, `hmac`, and encryption functions. Older text used `pgcrypto` for `gen_random_uuid()`. Prefer the built-in `gen_random_uuid()` for ids. Use `pgcrypto` for password hashes and digest work. Do not invent your own hash protocol.

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;
SELECT crypt('demo', gen_salt('bf'));
```

**`citext`** is a case-insensitive text type. `WHERE email = 'A@B.com'` matches `a@b.com`. A unique index on `citext` is case-insensitive. Prefer `citext` over `lower(email)` in every query when the column is always case-insensitive. You can also use `LOWER` plus a unique index on `lower(email)` without `citext`.

**`pg_trgm`** adds trigram similarity and can index `LIKE '%x%'` with GIN or GiST. Topic 7 previewed this. Topic 8 covered GIN.

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX items_name_trgm ON shop.items USING gin (name gin_trgm_ops);
SELECT name FROM shop.items WHERE name LIKE '%nail%';
```

**`pg_stat_statements`** records a normalized query and its total time. Topic 9 introduced it. You must load the library in `postgresql.conf`:

```text
shared_preload_libraries = 'pg_stat_statements'
```

Restart the server. Then `CREATE EXTENSION pg_stat_statements` in each database that must query the view. Topic 18 covers operations around this view.

Do not store reversible encryption keys in the same table as the ciphertext without a key plan. Do not enable `pg_stat_statements` and then ignore `pg_stat_statements.track`. Do not use `uuid-ossp` only to generate v4 UUIDs on 16 or 17.

### Questions

#### Theoretical questions

1. Do you need an extension for `gen_random_uuid()` on PostgreSQL 16 or 17?
2. When do you still create `uuid-ossp`?
3. What problem does `citext` solve?
4. What index class does `pg_trgm` often use for `LIKE '%x%'`?
5. Why is `CREATE EXTENSION pg_stat_statements` not enough by itself?

#### Easy practical tasks

1. Select `gen_random_uuid()` twice. Confirm that the values differ.
2. Create `citext` if allowed. Create a `citext` column. Insert `A@B.com` and search with `a@b.com`.
3. Create `pg_trgm` if allowed. Run `SELECT show_trgm('hello');`.
4. Query `pg_extension` for this list and mark which names are installed.

#### Medium practical tasks

1. Create a GIN trigram index. Compare `EXPLAIN` for `LIKE '%nail%'` before and after.
2. Use `pgcrypto` `crypt` and `gen_salt('bf')` on a lab password. Prove that a second `crypt` with the stored hash verifies the same password (read the `crypt` docs).
3. If `pg_stat_statements` is loaded, select the five statements with the highest `total_exec_time`. If it is not loaded, write the config line that you need.

#### Advanced practical tasks

1. Compare `citext` unique constraints with `CREATE UNIQUE INDEX ON t (lower(email))`. Write one benefit and one cost for each.
2. Read the 17 contrib pages for these five extensions. Write a table: extension, one function or view, when your app needs it.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. When do you pick a SQL function, PL/pgSQL, a trigger, or a generated column for the same business rule?
2. How do `GRANT EXECUTE`, `SECURITY INVOKER`, and extensions interact for an application role?
3. Why must `updated_at` and `price_gross` use different tools in this topic?
4. What extra operational step does `pg_stat_statements` need that `citext` does not need?
5. A teammate wants all application logic in triggers. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. In one schema, create a SQL function, a PL/pgSQL function, a generated column, and a `BEFORE UPDATE` trigger. Run one call or write that proves each object works.
2. Write a cheat sheet: SQL function, PL/pgSQL, trigger timings, `STORED`, `CREATE EXTENSION`, the five common extensions.
3. Run `\dx`, `\df`, and `\d` on your practice table. Save the three outputs.
4. Select `gen_random_uuid()` and `version();`. Write the major version.

#### Medium practical tasks

1. Build a small `items` table: identity id, `name`, `qty`, generated label, `updated_at` trigger, `low_stock` function. Grant only what `shop_app` needs.
2. Enable `pg_trgm` and add a GIN index. Show `EXPLAIN` for a substring search.
3. Document who may run `CREATE EXTENSION` on your lab (superuser, cloud console, or neither).

#### Advanced practical tasks

1. Replace a PL/pgSQL row loop with a SQL function plus a generated column. Compare `EXPLAIN ANALYZE` on 10000 rows.
2. Map this topic to official chapters: query language functions, PL/pgSQL, triggers, generated columns, extensions. Add one fact from the docs that this topic did not include.
