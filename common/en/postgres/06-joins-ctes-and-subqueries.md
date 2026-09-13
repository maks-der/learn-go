# 6. Joins, CTEs, and Subqueries

## Description

This topic shows how you combine rows in PostgreSQL 16 and PostgreSQL 17. You learn inner, outer, and cross joins. You learn `WITH` (common table expressions), recursive queries, `LATERAL`, and `DISTINCT ON`.

Complete topic 5 first. You need tables with foreign keys. Complete this topic before you study functions and indexes.

Use one term for each concept. A join combines rows from two tables. A CTE is a named subquery in a `WITH` clause. A recursive CTE can read its own output. A lateral join lets the right side use columns from the left side. `DISTINCT ON` keeps the first row of each group in the current sort order.

---

## Inner / outer / cross joins

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

1. Draw a Venn-style diagram (boxes are fine) for inner, left, right, and full joins. Add one SQL example each.
2. Compare plans with `EXPLAIN` for a join with `ON` only and a join that uses `WHERE` on the right table. Write the join type that you see.

---

## `WITH` (common table expressions)

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
WITH
active_customers AS (
    SELECT id, email FROM customers WHERE email IS NOT NULL
),
order_counts AS (
    SELECT customer_id, count(*) AS n
    FROM orders
    GROUP BY customer_id
)
SELECT a.email, o.n
FROM active_customers AS a
JOIN order_counts AS o ON o.customer_id = a.id;
```

A CTE can be a `INSERT`, `UPDATE`, or `DELETE` with `RETURNING` (data-modifying CTE):

```sql
WITH moved AS (
    DELETE FROM scratch
    WHERE done
    RETURNING *
)
INSERT INTO scratch_archive
SELECT * FROM moved;
```

In PostgreSQL 12 and later, the planner can inline many `SELECT` CTEs. A CTE is not always an optimization fence. Use `WITH cte AS MATERIALIZED (...)` to force materialization. Use `NOT MATERIALIZED` to force inline.

Write a CTE when it makes the query easier to read. Write a subquery when the query is short. Do not nest five CTEs that you never reuse.

A CTE name hides a real table of the same name inside the statement. Use clear names.

### Questions

#### Theoretical questions

1. What is a CTE?
2. How do you write two CTEs in one statement?
3. What is a data-modifying CTE?
4. What does `MATERIALIZED` mean on a CTE?
5. When is a subquery clearer than a CTE?

#### Easy practical tasks

1. Write a CTE that selects three columns. Select from the CTE.
2. Write two CTEs. Join them.
3. Name a CTE `recent` and filter `created_at`.
4. Draw the order: `WITH`, main `SELECT`.

#### Medium practical tasks

1. Use a CTE with `GROUP BY` and join it to a dimension table.
2. Use `DELETE ... RETURNING` in a CTE and `INSERT` the rows into an archive table.
3. Run `EXPLAIN` on a CTE with and without `MATERIALIZED`. Write if the plan changed.

#### Advanced practical tasks

1. Rewrite a nested subquery as two CTEs. Keep the same result. Show both SQL texts.
2. Read "WITH Queries" in the official docs. Write the rules for visibility of data-modifying CTEs (when other CTEs see the change).

---

## `WITH RECURSIVE` for trees

A recursive CTE has an anchor part and a recursive part. The parts sit in a `UNION` (or `UNION ALL`).

```sql
CREATE TABLE org (
    id        int PRIMARY KEY,
    parent_id int REFERENCES org (id),
    name      text NOT NULL
);

INSERT INTO org VALUES
    (1, NULL, 'HQ'),
    (2, 1, 'Engineering'),
    (3, 1, 'Sales'),
    (4, 2, 'Platform');

WITH RECURSIVE tree AS (
    SELECT id, parent_id, name, 1 AS depth
    FROM org
    WHERE parent_id IS NULL
    UNION ALL
    SELECT o.id, o.parent_id, o.name, t.depth + 1
    FROM org AS o
    JOIN tree AS t ON o.parent_id = t.id
)
SELECT * FROM tree
ORDER BY depth, id;
```

The anchor selects the roots. The recursive part joins the base table to the CTE. PostgreSQL repeats the recursive part until it returns no row.

Use `UNION ALL` when you do not need to drop duplicate rows. Use `UNION` when you must drop duplicates. Cycles need extra care. Keep a path array and stop when `id` is already in the path:

```sql
WITH RECURSIVE tree AS (
    SELECT id, parent_id, name, ARRAY[id] AS path
    FROM org
    WHERE id = 1
    UNION ALL
    SELECT o.id, o.parent_id, o.name, t.path || o.id
    FROM org AS o
    JOIN tree AS t ON o.parent_id = t.id
    WHERE o.id <> ALL (t.path)
)
SELECT * FROM tree;
```

`SEARCH` and `CYCLE` clauses exist in current PostgreSQL (SQL-standard cycle detection). You can use them instead of a hand-built path array. See the official `WITH` docs.

Do not recurse without a stop condition. A cycle without a guard can run until the server hits a stack limit, a statement timeout, or a cancel.

Recursive CTEs also generate number series. Prefer `generate_series` for simple ranges.

### Questions

#### Theoretical questions

1. What are the two parts of a recursive CTE?
2. What does the anchor query select in a tree?
3. Why can a cycle be a problem?
4. What does a path array prevent?
5. When do you use `UNION ALL` instead of `UNION`?

#### Easy practical tasks

1. Create `org` and insert the four rows. Run the first recursive query.
2. Add `depth` and order by `depth`.
3. Start the tree at `id = 2` instead of the root. List the rows.
4. Write four sentences: anchor, recursive part, join, stop.

#### Medium practical tasks

1. Add a cycle on a copy of `org`. Query without a guard (use a small `LIMIT` in the outer query if needed). Then add a path guard.
2. Return a text path (`HQ > Engineering > Platform`) with `string_agg` or `array_to_string`.
3. Read `SEARCH` and `CYCLE` in the docs. Rewrite the path example with `CYCLE` if you want the standard form.

#### Advanced practical tasks

1. Compute the count of descendants for each node with a recursive CTE.
2. Compare a recursive CTE with a recursive function in PL/pgSQL (topic 13 preview). Write one benefit of the CTE.

---

## Lateral joins (`LEFT JOIN LATERAL`)

A normal join cannot use a left-row column inside a right-side subquery in the `FROM` list. `LATERAL` allows that.

```sql
SELECT c.id, c.email, last_order.id, last_order.created_at
FROM customers AS c
LEFT JOIN LATERAL (
    SELECT o.id, o.created_at
    FROM orders AS o
    WHERE o.customer_id = c.id
    ORDER BY o.created_at DESC
    LIMIT 1
) AS last_order ON true;
```

For each customer, the subquery reads that `c.id` and returns the latest order. `LEFT JOIN LATERAL` keeps customers that have no order. `last_order` columns are `NULL` for those customers.

`ON true` is common when the lateral subquery already filters. You still need a join condition clause for `LEFT JOIN`.

`CROSS JOIN LATERAL` is the inner form. It drops left rows when the subquery returns no row.

```sql
SELECT c.email, s.n
FROM customers AS c
CROSS JOIN LATERAL (
    SELECT count(*) AS n
    FROM orders o
    WHERE o.customer_id = c.id
) AS s;
```

Set-returning functions in `FROM` are implicit lateral:

```sql
SELECT u.id, u.tag
FROM users AS u
CROSS JOIN LATERAL unnest(u.tags) AS tag;
```

You can write `FROM users u, unnest(u.tags) AS tag`. The comma form is an implicit cross join. Prefer `JOIN LATERAL` for clarity.

Use `LATERAL` for "top N per group" with `LIMIT`. Window functions are another method (topic 20). `DISTINCT ON` is a third method (next section).

A lateral subquery that runs a heavy query per left row can be slow. Check `EXPLAIN`. A join plus window function can be faster on large sets.

### Questions

#### Theoretical questions

1. What does `LATERAL` allow that a normal join does not?
2. Why does the example use `LEFT JOIN LATERAL ... ON true`?
3. How is `CROSS JOIN LATERAL` different from `LEFT JOIN LATERAL`?
4. Why is `unnest` in `FROM` lateral?
5. When can a lateral "top 1" query be slow?

#### Easy practical tasks

1. For each customer, return the latest order with `LEFT JOIN LATERAL` and `LIMIT 1`.
2. Change the join to `CROSS JOIN LATERAL`. Compare the row count.
3. `unnest` a `text[]` column with `CROSS JOIN LATERAL`.
4. Write the meaning of `ON true` in four sentences.

#### Medium practical tasks

1. Return the two latest orders per customer (`LIMIT 2`).
2. Rewrite the latest-order query with a CTE that ranks rows. Compare results with the lateral query.
3. Run `EXPLAIN` on the lateral query. Write the join type.

#### Advanced practical tasks

1. Compare `LATERAL LIMIT 1`, `DISTINCT ON`, and a window `row_number()` for "latest row per group". Write one table: method, SQL size, typical use.
2. Read "LATERAL" in the `SELECT` docs. Write one extra example that uses a function call.

---

## Distinct on (`DISTINCT ON`) — PostgreSQL-specific

`DISTINCT ON (expr [, ...])` keeps the first row of each distinct value of the expressions. "First" follows `ORDER BY`.

```sql
SELECT DISTINCT ON (customer_id)
    customer_id, id, created_at
FROM orders
ORDER BY customer_id, created_at DESC;
```

This returns one row per `customer_id`: the latest `created_at`. The `ORDER BY` must start with the same expressions as `DISTINCT ON`, in the same order. Then you add extra sort keys (`created_at DESC`).

Rules:

- `DISTINCT ON` is a PostgreSQL extension. It is not in the SQL standard.
- `DISTINCT` without `ON` drops full-row duplicates. It is different.
- Without a matching `ORDER BY`, the chosen row is not defined.
- `SELECT DISTINCT ON (a) b` can return any `b` for that `a` if `ORDER BY` does not pin `b`.

Rewrite with a window function when you need standard SQL:

```sql
SELECT customer_id, id, created_at
FROM (
    SELECT
        customer_id, id, created_at,
        row_number() OVER (
            PARTITION BY customer_id
            ORDER BY created_at DESC
        ) AS rn
    FROM orders
) AS s
WHERE rn = 1;
```

Use `DISTINCT ON` for short, clear PostgreSQL queries. Use the window form when you must run the same SQL on another product.

`DISTINCT ON` can use an index that matches the `ORDER BY`. Topic 9 shows how you read the plan.

### Questions

#### Theoretical questions

1. What does `DISTINCT ON (customer_id)` keep?
2. What must `ORDER BY` start with?
3. How is `DISTINCT` different from `DISTINCT ON`?
4. Is `DISTINCT ON` in the SQL standard?
5. What happens when `ORDER BY` does not match `DISTINCT ON`?

#### Easy practical tasks

1. Insert two orders for one customer. Select the latest with `DISTINCT ON`.
2. Change `ORDER BY` to `created_at ASC` and compare which `id` you get.
3. Run `SELECT DISTINCT customer_id FROM orders`.
4. Write the `ORDER BY` rule in one sentence.

#### Medium practical tasks

1. Return the latest order per customer with `DISTINCT ON` and with `row_number()`. Confirm equal keys.
2. Add a second `DISTINCT ON` expression `(customer_id, status)` on a table that has `status`.
3. Run `EXPLAIN` on the `DISTINCT ON` query. Write the first plan node.

#### Advanced practical tasks

1. Read the `SELECT` reference for `DISTINCT ON`. Copy the official `ORDER BY` rule in your own words.
2. Find a case where `DISTINCT ON` and `GROUP BY` are not interchangeable (you need extra columns that are not aggregates). Write both attempts.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Pick a method for "latest child row per parent": left join, CTE, recursive CTE, lateral, or `DISTINCT ON`. Give one reason.
2. When must you use a recursive CTE instead of a normal join?
3. How do data-modifying CTEs and `RETURNING` move rows between tables?
4. Why does a cross join of two large tables cause a problem?
5. What shared rule do `DISTINCT ON` and `LATERAL ... LIMIT 1` both need (a defined order)?

#### Easy practical tasks

1. Write one inner join, one left join, one CTE, and one `DISTINCT ON` on your shop tables.
2. Write a cheat sheet: join types, `WITH`, `WITH RECURSIVE`, `LATERAL`, `DISTINCT ON`.
3. List customers with no orders (`LEFT JOIN ... WHERE o.id IS NULL`).
4. Draw a tree of `org` and write the recursive query that prints it.

#### Medium practical tasks

1. Combine a CTE and a lateral join: CTE filters customers, lateral fetches latest order.
2. Explain three queries with `EXPLAIN` (join, CTE, `DISTINCT ON`). Write the join or scan type for each.
3. Document when your team allows `DISTINCT ON` versus window functions.

#### Advanced practical tasks

1. Write one query that archives old orders with a modifying CTE and prints how many rows moved.
2. Recurse a graph that is not a tree (two parents). Document whether your SQL is a tree walk or a graph walk. Add a cycle guard.
