# 3. SQL in PostgreSQL

## Description

This topic shows core SQL in PostgreSQL 16 and PostgreSQL 17. You write `SELECT`, `INSERT`, `UPDATE`, and `DELETE`. You learn `RETURNING`, `ILIKE`, and dollar-quoting. You also learn joins, `WITH`, recursive queries, `LEFT JOIN LATERAL`, and `DISTINCT ON`.

Complete topics 1 and 2 first. You need a connection and tables with keys. Complete this topic before you study functions, JSON, and indexes.

Use one term for each concept. A statement is one SQL command. `RETURNING` sends result rows from a write. Dollar-quoting is a way to write a string without quote escaping. A join combines rows from two tables. A CTE is a named subquery in a `WITH` clause. A lateral join lets the right side use columns from the left side. `DISTINCT ON` keeps the first row of each group in the current sort order.

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
- `INSERT` must satisfy constraints.
- `SELECT` without `FROM` is valid. Example: `SELECT 1 + 1;`.
- PostgreSQL uses `$1` style parameters in client APIs. In `psql` you write literals. Do not concatenate user text into SQL. Use parameters in programs.

`INSERT ... SELECT` copies from a query:

```sql
INSERT INTO items (id, name, qty)
SELECT id + 100, name, qty FROM items WHERE qty < 10;
```

`UPDATE` can use `FROM` to join other tables. `DELETE` can use `USING`. Learn those forms after you learn joins in this topic.

`LIMIT` and `OFFSET` cap a result. `OFFSET` skips rows. Large `OFFSET` is slow. Topic 13 covers keyset pagination.

```sql
SELECT id, name FROM items ORDER BY id LIMIT 10 OFFSET 20;
```

`FETCH FIRST 10 ROWS ONLY` is the SQL-standard form of `LIMIT`.

PostgreSQL does not require `FROM DUAL`. Use `SELECT` without `FROM` for expressions.

### Questions

#### Theoretical questions

1. What does each of `SELECT`, `INSERT`, `UPDATE`, and `DELETE` do?
2. What happens when `UPDATE` has no `WHERE`?
3. Is `SELECT` without `FROM` valid?
4. Why must programs use parameters instead of string concatenation?
5. What does `INSERT ... SELECT` do?

#### Easy practical tasks

1. Create `items`. Insert three rows. Select all rows.
2. Update one row. Select it again.
3. Delete one row. Count the remaining rows.
4. Run `SELECT 1 + 1;` with no `FROM`.

#### Medium practical tasks

1. Use `INSERT ... SELECT` to copy some rows with new ids. Show the table.
2. `UPDATE items SET qty = qty + 1 FROM ...` after you add a second table. Write the statement.
3. Compare `LIMIT 5` and `FETCH FIRST 5 ROWS ONLY` on the same `ORDER BY`.

#### Advanced practical tasks

1. Write a `DELETE` that uses `USING` to remove items that match a second table.
2. Read the `SELECT` reference. List five clauses that this section did not show (`GROUP BY` may be one).

---

## `RETURNING`, `ILIKE`, dollar-quoting

`RETURNING` adds a result set to `INSERT`, `UPDATE`, or `DELETE`. You get the rows that the statement wrote or removed.

```sql
INSERT INTO items (id, name, qty)
VALUES (4, 'washer', 10)
RETURNING id, name;

UPDATE items
SET qty = qty - 1
WHERE id = 4
RETURNING id, qty;

DELETE FROM items
WHERE id = 4
RETURNING *;
```

Use `RETURNING` when the database fills identity, defaults, or trigger values. The client does not need a second `SELECT`.

`LIKE` matches a pattern. `%` means any string. `_` means one character. `LIKE` is case-sensitive for typical `C` collation text.

`ILIKE` is the case-insensitive form. It is a PostgreSQL extension.

```sql
SELECT name FROM items WHERE name LIKE 'n%';
SELECT name FROM items WHERE name ILIKE 'N%';
```

Escape `%` and `_` with `ESCAPE` when the user text can contain those characters. Topic 4 covers `SIMILAR TO` and regular expressions.

Dollar-quoting writes a string without doubled single quotes:

```sql
SELECT $$it's a nail$$;
SELECT $body$line 1
line 2$body$;
```

The tag between `$` signs must match. Use dollar-quoting for function bodies (topic 9) and for long text in scripts. A dollar-quoted string is still a string literal. It is not a parameter. Do not put user input into dollar quotes in application SQL.

PostgreSQL also casts types in many contexts. Write explicit casts when the type is not obvious: `qty::text` or `CAST(qty AS text)`.

### Questions

#### Theoretical questions

1. What does `RETURNING` return for `DELETE`?
2. How does `ILIKE` differ from `LIKE`?
3. What do `%` and `_` mean in `LIKE`?
4. Why is dollar-quoting useful in a function body?
5. Is a dollar-quoted string a parameter?

#### Easy practical tasks

1. Insert a row with `RETURNING id, name`.
2. Update that row with `RETURNING qty`.
3. Select names with `ILIKE '%a%'`.
4. Run `SELECT $$it's$$;` and the same text with single-quote escaping.

#### Medium practical tasks

1. Delete rows that match a pattern and `RETURNING *` into a look at the removed names.
2. Show `LIKE 'N%'` versus `ILIKE 'N%'` on mixed-case names.
3. Write a function-body style dollar-quoted string that contains a single quote and a newline.

#### Advanced practical tasks

1. Combine `INSERT ... SELECT` with `RETURNING`. Show the new ids.
2. Read `LIKE` escape rules. Write a pattern that matches a literal `%`.

---

## Joins

An **inner join** returns rows that match the join condition.

```sql
SELECT c.email, o.id AS order_id
FROM customers AS c
INNER JOIN orders AS o ON o.customer_id = c.id;
```

`JOIN` without a word is an inner join. Write `INNER JOIN` when you teach.

A **left outer join** returns every row from the left table. When no match exists, the right columns are `NULL`.

```sql
SELECT c.email, o.id AS order_id
FROM customers AS c
LEFT JOIN orders AS o ON o.customer_id = c.id;
```

A **right outer join** returns every row from the right table. A **full outer join** returns every row from both sides.

```sql
SELECT *
FROM customers AS c
FULL JOIN orders AS o ON o.customer_id = c.id;
```

A **cross join** is the Cartesian product. Every left row pairs with every right row. There is no `ON` clause.

```sql
SELECT *
FROM sizes
CROSS JOIN colors;
```

Use a cross join only when you need all pairs (example: size × color).

`USING (customer_id)` is a short form when both columns have the same name. `NATURAL JOIN` matches all same-named columns. Do not use `NATURAL JOIN`. A new column can change the join.

Filter in `ON` versus `WHERE` for outer joins:

- `ON` decides matches
- `WHERE` on a right column after `LEFT JOIN` can remove unmatched left rows and turn the join into an inner join

```sql
-- keeps customers without orders
SELECT c.email, o.id
FROM customers c
LEFT JOIN orders o ON o.customer_id = c.id AND o.id > 10;

-- removes customers with no matching order
SELECT c.email, o.id
FROM customers c
LEFT JOIN orders o ON o.customer_id = c.id
WHERE o.id > 10;
```

Always write aliases. Always write the join condition. Do not join without `ON` unless you mean `CROSS JOIN`.

`UPDATE ... FROM` and `DELETE ... USING` are joins for writes. Put the join in `FROM`/`USING`. Keep `WHERE` for the extra filter.

### Questions

#### Theoretical questions

1. What rows does an inner join return?
2. What rows does a left join return?
3. What is a cross join?
4. Why must you avoid `NATURAL JOIN`?
5. How can `WHERE` after a left join behave like an inner join?

#### Easy practical tasks

1. Create `customers` and `orders`. Insert one customer with orders and one customer without orders. Run an inner join.
2. Run a left join on the same data. Count the rows.
3. Run a cross join of two tiny tables (2 × 3). Count the rows.
4. Rewrite an inner join with `USING` if the column names match.

#### Medium practical tasks

1. Show the `ON` versus `WHERE` pair from this section. Write the two row counts.
2. Write a full outer join. Insert an orphan order if your foreign key allows it, or drop the FK on a copy.
3. Join three tables: customers, orders, order_lines. Use inner joins. Order the result.

#### Advanced practical tasks

1. Draw a box diagram for inner, left, right, and full joins. Add one SQL example each.
2. Write an `UPDATE` that uses `FROM` to copy a value from another table. Use `RETURNING`.

---

## `WITH` and `WITH RECURSIVE`

A CTE names a query. The main query reads that name like a table.

```sql
WITH recent AS (
    SELECT id, customer_id, created_at
    FROM orders
    WHERE created_at >= DATE '2026-01-01'
)
SELECT c.email, r.id
FROM recent AS r
JOIN customers AS c ON c.id = r.customer_id;
```

You can list more than one CTE. Separate them with commas.

```sql
WITH a AS (...),
     b AS (SELECT * FROM a WHERE ...)
SELECT * FROM b;
```

By default a CTE is an optimization fence in older mental models. In PostgreSQL 12 and later the planner can inline many CTEs. You can force materialization with `WITH a AS MATERIALIZED (...)`. You can allow inline with `NOT MATERIALIZED`.

Use a CTE to name a step. Do not nest ten CTEs when a simple join is enough.

`WITH RECURSIVE` builds a working table. The first part (the anchor) runs once. The second part (the recursive term) reads the working table until it adds no row.

```sql
WITH RECURSIVE tree AS (
    SELECT id, parent_id, name, 1 AS depth
    FROM categories
    WHERE parent_id IS NULL
    UNION ALL
    SELECT c.id, c.parent_id, c.name, t.depth + 1
    FROM categories AS c
    JOIN tree AS t ON c.parent_id = t.id
)
SELECT * FROM tree ORDER BY depth, name;
```

Use `UNION ALL` when you do not need to drop duplicate rows. Use `UNION` when you must drop duplicates. Always include a stop condition (depth limit or a finite tree). A cycle without a visited-set check can loop until `statement_timeout` or memory limits.

PostgreSQL 14 and later can use `SEARCH` and `CYCLE` clauses on recursive CTEs. PostgreSQL 16 and 17 keep those clauses. They mark cycles and search order.

A data-modifying CTE can `INSERT`, `UPDATE`, or `DELETE` in `WITH` and then `SELECT` the `RETURNING` rows. The writes run in one statement.

```sql
WITH moved AS (
    DELETE FROM items WHERE qty = 0 RETURNING *
)
INSERT INTO items_archive SELECT * FROM moved;
```

Do not use recursion for a simple list that `generate_series` can build.

### Questions

#### Theoretical questions

1. What is a CTE?
2. What does the recursive term read?
3. When do you write `UNION ALL` in a recursive CTE?
4. What does `AS MATERIALIZED` force?
5. Can a CTE contain `DELETE`?

#### Easy practical tasks

1. Write a non-recursive `WITH` that filters `orders` and then joins `customers`.
2. List two CTEs in one `WITH`. Select from the second.
3. Build a category tree with three rows. Run the recursive query.
4. Write four sentences: CTE, anchor, recursive term, cycle.

#### Medium practical tasks

1. Add a depth limit (`WHERE t.depth < 5`) to the tree query.
2. Compare `EXPLAIN` for a CTE with and without `MATERIALIZED` on a query that you invent.
3. Use a data-modifying CTE to delete zero-qty items and insert them into an archive table.

#### Advanced practical tasks

1. Read `SEARCH` and `CYCLE` in the 17 docs. Add `CYCLE` to a graph that has a loop. Show the cycle mark.
2. Write a recursive query that walks a manager chain (`employees.manager_id`). Stop at the root.

---

## `LEFT JOIN LATERAL`

A lateral join lets the right-hand query use columns from the left-hand row.

```sql
SELECT c.id, c.email, recent.id AS order_id, recent.created_at
FROM customers AS c
LEFT JOIN LATERAL (
    SELECT o.id, o.created_at
    FROM orders AS o
    WHERE o.customer_id = c.id
    ORDER BY o.created_at DESC
    LIMIT 1
) AS recent ON true;
```

Without `LATERAL`, the subquery in `FROM` cannot see `c.id`. `LATERAL` makes that reference legal.

`LEFT JOIN LATERAL ... ON true` keeps every left row. When the subquery returns no row, the right columns are `NULL`. `CROSS JOIN LATERAL` or `INNER JOIN LATERAL` drops left rows that have no match.

A function in `FROM` is implicitly lateral:

```sql
SELECT u.id, g
FROM users AS u
CROSS JOIN LATERAL unnest(u.tags) AS g;
```

`unnest` is topic 4. The pattern is the same: each left row produces a set.

Use `LATERAL` for "top N per group" when a window function (topic 14) is not yet in your toolkit, or when a set-returning function needs left columns.

A correlated subquery in `SELECT` can do similar work. `LATERAL` in `FROM` is often clearer and can be faster when you need several columns from the inner query.

Do not use `LATERAL` for a normal join that `ON` can express. Do not forget `ON true` on a `LEFT JOIN LATERAL` subquery. A missing `ON` is a syntax error for `LEFT JOIN`.

### Questions

#### Theoretical questions

1. What extra right does `LATERAL` give the right-hand query?
2. Why use `LEFT JOIN LATERAL` instead of `CROSS JOIN LATERAL` for "latest order per customer"?
3. What does `ON true` mean in this pattern?
4. Is a function in `FROM` lateral?
5. When must you not use `LATERAL`?

#### Easy practical tasks

1. Create customers and orders. Write the "latest order per customer" query from this section.
2. Change it to `CROSS JOIN LATERAL`. Count rows versus the left join.
3. Write one sentence that explains why `c.id` is visible inside the subquery.
4. Draw left row → inner `LIMIT 1` → output columns.

#### Medium practical tasks

1. Return the latest two orders per customer (`LIMIT 2`). Count rows.
2. Compare a correlated subquery in `SELECT` that returns one id with the `LATERAL` form that returns two columns.
3. Use `LATERAL` with `generate_series(1, n)` where `n` is a column on the left table.

#### Advanced practical tasks

1. `EXPLAIN (ANALYZE, BUFFERS)` the latest-order query. Write the join type that you see.
2. Rewrite the same latest-order problem with `DISTINCT ON` (next section). Compare the two plans.

---

## `DISTINCT ON`

`DISTINCT ON (expr [, ...])` keeps the first row of each distinct value of those expressions. "First" means the first row in the current `ORDER BY`.

```sql
SELECT DISTINCT ON (customer_id)
    customer_id, id, created_at
FROM orders
ORDER BY customer_id, created_at DESC;
```

This query keeps one row per `customer_id`: the latest `created_at`. The `ORDER BY` must start with the `DISTINCT ON` expressions. Then you add the sort that picks the winner (`created_at DESC`).

`DISTINCT ON` is a PostgreSQL extension. It is not standard `DISTINCT`. `DISTINCT` without `ON` drops full-row duplicates.

```sql
SELECT DISTINCT customer_id FROM orders;
```

Rules:

- `ORDER BY` must begin with the same expressions as `DISTINCT ON` (or expressions that match them).
- Without a matching `ORDER BY`, the winner is not defined.
- `DISTINCT ON` runs after `WHERE` and before `LIMIT`.

Use `DISTINCT ON` for "best row per group" when the rule is a sort. A window function (`ROW_NUMBER()`) can do the same job and is standard SQL (topic 14). `DISTINCT ON` is shorter.

Do not use `DISTINCT ON` to hide a bad join that multiplies rows. Fix the join. Do not omit `ORDER BY`.

### Questions

#### Theoretical questions

1. Which row does `DISTINCT ON` keep?
2. What must `ORDER BY` start with?
3. How does `DISTINCT` differ from `DISTINCT ON`?
4. Is `DISTINCT ON` standard SQL?
5. When do you prefer a window function instead?

#### Easy practical tasks

1. Insert two orders per customer. Run the latest-order `DISTINCT ON` query.
2. Change `ORDER BY` to `created_at ASC`. Write which order you keep.
3. Run `SELECT DISTINCT customer_id FROM orders;`.
4. Write four sentences: winner, `ORDER BY`, `DISTINCT`, extension.

#### Medium practical tasks

1. Keep the order with the largest `total` per customer. Show the SQL.
2. Add `LIMIT 10` after `DISTINCT ON`. Write how many customers you can see at most.
3. Compare `DISTINCT ON` with `GROUP BY customer_id` plus `MAX(created_at)`. Write which columns you can keep in each form.

#### Advanced practical tasks

1. `EXPLAIN ANALYZE` `DISTINCT ON` versus `LATERAL LIMIT 1` versus `ROW_NUMBER()` if you already know windows. Write the cheapest plan on your data.
2. Read the `SELECT` reference for `DISTINCT ON`. Copy the official `ORDER BY` rule in your own words.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `RETURNING` and a later `SELECT` differ after an `INSERT` of an identity row?
2. When does a left join plus `WHERE` on the right table lie about "all customers"?
3. What stops a recursive CTE from running forever?
4. Which problem can both `LATERAL` and `DISTINCT ON` solve?
5. Why is `OFFSET 100000` a poor way to page through `orders`?

#### Easy practical tasks

1. Insert an item with dollar-quoted text that contains a quote. `RETURNING` the `id` and `name`.
2. Join customers to orders with `LEFT JOIN`. List customers who have no order (`WHERE o.id IS NULL`).
3. Write a `WITH` named `low` for items with `qty < 5`. Select from `low`.
4. Use `ILIKE` to find names that contain `a` in any case.

#### Medium practical tasks

1. In one script: insert a customer, insert two orders, return the latest order with `DISTINCT ON`, and delete zero-qty items with `RETURNING`.
2. Write a recursive CTE for comments that reply to comments (`parent_id`). Show depth.
3. Update `items.qty` from a `restock` table using `UPDATE ... FROM`. Return new qty values.

#### Advanced practical tasks

1. Build a small shop query file: one inner join, one left join, one CTE, one `LATERAL` latest row, one `DISTINCT ON`. Add `EXPLAIN` for two of them.
2. Read "WITH Queries" and "DISTINCT" in the 17 docs. Write a one-page cheat sheet that a teammate can print.
