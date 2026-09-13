# 6. Aggregation

## Description

This topic shows how you compute totals and how you group rows. You learn aggregate functions, `GROUP BY`, and `HAVING`. You also learn how `NULL` behaves in aggregates.

Use one term for each concept. An aggregate reduces many values to one value. A group is a set of rows that share the same grouping key. Complete this topic before you define schemas in detail.

Name the grain of every result. A total without `GROUP BY` is one row. A total with `GROUP BY customer_id` is one row per customer.

---

## `COUNT`, `SUM`, `AVG`, `MIN`, `MAX`

An aggregate function takes a set of values and returns one value.

**`COUNT(*)`.** Count rows in the set. This count includes rows where other columns are `NULL`.

**`COUNT(column)`.** Count rows where that column is not `NULL`.

**`COUNT(DISTINCT column)`.** Count unique non-null values in that column.

**`SUM(column)`.** Add numeric values. Ignore `NULL`. If every value is `NULL`, `SUM` returns `NULL`.

**`AVG(column)`.** Compute the average of non-null numeric values. `AVG` is `SUM / COUNT(column)`, not `SUM / COUNT(*)`, when some values are `NULL`.

**`MIN(column)` and `MAX(column)`.** Return the smallest or largest non-null value. For text, the order follows the collation of the DBMS.

```sql
SELECT
    COUNT(*) AS order_rows,
    COUNT(shipped_at) AS shipped_rows,
    SUM(qty) AS total_qty,
    AVG(qty) AS avg_qty,
    MIN(ordered_at) AS first_order,
    MAX(ordered_at) AS last_order
FROM orders;
```

Without `GROUP BY`, the query returns one row for the whole table (after `WHERE` if you add it).

Do not mix a bare column and an aggregate in the `SELECT` list unless that column is in `GROUP BY`. This rule is the next section.

`COUNT(*)` never returns `NULL`. It returns `0` for an empty set. `SUM` of an empty set is `NULL` in SQL. That difference surprises beginners.

```sql
SELECT COUNT(*) AS n, SUM(qty) AS s
FROM order_lines
WHERE order_id = -1;
```

If no line matches, `n` is `0` and `s` is `NULL`.

Use `COALESCE(SUM(qty), 0)` when you must show `0` for an empty set.

### Questions

#### Theoretical questions

1. What does `COUNT(*)` count?
2. How does `COUNT(column)` differ from `COUNT(*)`?
3. What does `SUM` do with `NULL` values?
4. Why is `AVG` not always `SUM / COUNT(*)`?
5. What does `SUM` return on an empty set?

#### Easy practical tasks

1. Count all rows in `customers`.
2. Count non-null emails if you allow `NULL` email.
3. Compute `MIN(price)` and `MAX(price)` on `products`.
4. Run `SUM` and `COUNT(*)` on a filter that matches no row. Write the two results.

#### Medium practical tasks

1. Compare `COUNT(*)` and `COUNT(shipped_at)` on `orders`.
2. Compute `COUNT(DISTINCT customer_id)` from `orders`.
3. Show `AVG` of a column that contains `NULL`. Write the input values and the result.

#### Advanced practical tasks

1. Write a one-page note: `COUNT(*)`, `COUNT(col)`, `COUNT(DISTINCT col)`, and empty-set `SUM`. Give one example each.
2. Explain when you wrap `SUM` in `COALESCE`. Give a report that must show zero sales.

---

## `GROUP BY`

`GROUP BY` splits rows into groups. Each group shares the same values of the grouping columns. Aggregates run inside each group.

```sql
SELECT customer_id, COUNT(*) AS order_count
FROM orders
GROUP BY customer_id;
```

Each `customer_id` produces one result row. `COUNT(*)` is the number of orders for that customer.

You may group by more than one column:

```sql
SELECT customer_id, status_code, COUNT(*) AS n
FROM orders
GROUP BY customer_id, status_code;
```

The `SELECT` list may contain:

- grouping columns
- aggregates
- expressions that depend only on grouping columns and aggregates

Do not select `name` when you group only by `customer_id`, unless `name` is functionally determined and your DBMS allows it. Standard SQL requires `name` in `GROUP BY` or inside an aggregate. The safe habit is: every non-aggregate column in `SELECT` appears in `GROUP BY`.

```sql
SELECT c.customer_id, c.name, COUNT(*) AS order_count
FROM customers AS c
INNER JOIN orders AS o
    ON o.customer_id = c.customer_id
GROUP BY c.customer_id, c.name;
```

Logical idea:

1. Apply `FROM` and `WHERE`.
2. Form groups with `GROUP BY`.
3. Compute aggregates per group.
4. Apply `HAVING` (next section).
5. Project the `SELECT` list.
6. Apply `ORDER BY`.

`WHERE` filters rows before groups form. Do not put a `COUNT(*)` test in `WHERE`. Use `HAVING`.

A query with aggregates and no `GROUP BY` is one group: the whole filtered table.

### Questions

#### Theoretical questions

1. What does `GROUP BY` create?
2. Which columns may appear in `SELECT` with `GROUP BY`?
3. When does a query without `GROUP BY` still have one group?
4. Why do you add `c.name` to `GROUP BY` when you select `c.name`?
5. Does `WHERE` run before or after groups form?

#### Easy practical tasks

1. Count orders per `customer_id`.
2. Count products per `category`.
3. Group by two columns on a small table. Show the groups by hand.
4. Run a query that selects a non-grouped column with an aggregate. Record the error (or the product extension).

#### Medium practical tasks

1. Join customers and orders. Group by customer id and name. Order by `order_count` descending.
2. Add `WHERE ordered_at >= DATE '2026-01-01'` to a grouped query. Explain that the filter is pre-group.
3. Show that two customers with the same name are two groups when you group by `customer_id`.

#### Advanced practical tasks

1. Write a report: per customer, order count and total qty (join lines). State the grain and the `GROUP BY` list.
2. Document whether your DBMS allows `SELECT customer_id, name` when you `GROUP BY customer_id` and `customer_id` is a primary key.

---

## `HAVING` vs `WHERE`

`WHERE` filters rows before grouping. `HAVING` filters groups after aggregation.

```sql
SELECT customer_id, COUNT(*) AS order_count
FROM orders
WHERE ordered_at >= DATE '2026-01-01'
GROUP BY customer_id
HAVING COUNT(*) >= 3;
```

`WHERE` keeps orders from the year 2026. `GROUP BY` builds per-customer counts from those orders. `HAVING` keeps customers with at least three such orders.

You cannot write `WHERE COUNT(*) >= 3`. `COUNT(*)` does not exist at the row-filter stage.

You can write `HAVING` on a group key, but `WHERE` is clearer for that case.

```sql
SELECT customer_id, COUNT(*) AS order_count
FROM orders
GROUP BY customer_id
HAVING customer_id = 1;
```

Prefer:

```sql
SELECT customer_id, COUNT(*) AS order_count
FROM orders
WHERE customer_id = 1
GROUP BY customer_id;
```

The second form throws away other customers before it groups. That form is cheaper and easier to read.

`HAVING` may use aliases in some products. Standard SQL does not require that. Repeat the aggregate in `HAVING`:

```sql
HAVING COUNT(*) >= 3
```

Do not use `HAVING` as a substitute for `WHERE` on base columns. Use `WHERE` for row rules. Use `HAVING` for group rules.

A left join plus `GROUP BY` plus `HAVING COUNT(o.order_id) = 0` can find customers with no orders. `COUNT(o.order_id)` ignores `NULL`. `COUNT(*)` would be `1` for a customer with a single all-null order side. Prefer `COUNT(o.order_id)` or `NOT EXISTS` for that question.

### Questions

#### Theoretical questions

1. When does `WHERE` filter?
2. When does `HAVING` filter?
3. Why is `WHERE COUNT(*) >= 3` invalid?
4. When do you prefer `WHERE` over `HAVING` for a customer id?
5. Why can `COUNT(*)` after a left join mislead you?

#### Easy practical tasks

1. List customers with at least two orders using `HAVING`.
2. Add a `WHERE` date filter to that query.
3. Rewrite a `HAVING customer_id = 1` into a `WHERE`.
4. Write four sentences that contrast `WHERE` and `HAVING`.

#### Medium practical tasks

1. Find categories with `AVG(price) > 20`. Use `HAVING`.
2. Find categories with at least one product cheaper than `5`. Decide `WHERE` versus `HAVING` and justify.
3. Compare `COUNT(*)` and `COUNT(o.order_id)` for customers left-joined to orders.

#### Advanced practical tasks

1. Write a report: customers with three or more orders in 2026 and a total qty of at least 10. Use `WHERE` and `HAVING`.
2. Explain in one page why `HAVING` is not "a second WHERE" even though the syntax looks similar.

---

## NULL in aggregates

`NULL` means missing. Aggregate functions skip `NULL` inputs, except `COUNT(*)`.

Examples with values `{10, NULL, 30}`:

| Function | Result |
| --- | --- |
| `COUNT(*)` | `3` if these are three rows |
| `COUNT(value)` | `2` |
| `SUM(value)` | `40` |
| `AVG(value)` | `20` |
| `MIN(value)` | `10` |
| `MAX(value)` | `30` |

`AVG` uses two values, not three. The missing value does not count as zero.

If you store `0` instead of `NULL` for "no measurement," `AVG` and `SUM` change. That choice is a model decision. Do not mix `0` and `NULL` for the same meaning.

`GROUP BY` treats `NULL` as one group. All rows with `NULL` in the grouping column go together.

```sql
SELECT coupon_id, COUNT(*) AS n
FROM orders
GROUP BY coupon_id;
```

Orders with no coupon (`coupon_id` `NULL`) form one group. That group is valid.

`COUNT(DISTINCT col)` ignores `NULL`. Several `NULL` values do not add a distinct item.

In `HAVING`, `SUM(col) > 0` is unknown when `SUM` is `NULL`. The group does not pass. Use `COALESCE(SUM(col), 0) > 0` if you must treat empty sums as zero.

Do not use `NULL` as a secret extra category unless you document it. Prefer a real lookup value for "none" when reports must sort and filter that category easily.

### Questions

#### Theoretical questions

1. Which aggregate does not ignore `NULL` in the same way as the others?
2. Why is `AVG` of `{10, NULL, 30}` equal to `20`?
3. How does `GROUP BY` treat `NULL` in the grouping column?
4. Does `COUNT(DISTINCT col)` count `NULL`?
5. Why can `HAVING SUM(col) > 0` drop a group?

#### Easy practical tasks

1. Insert three numeric values including one `NULL`. Compute `SUM`, `AVG`, and `COUNT(col)`.
2. Group by a nullable column. Show the `NULL` group.
3. Compare `COUNT(*)` and `COUNT(col)` on that set.
4. Write four sentences on `0` versus `NULL` in an amount column.

#### Medium practical tasks

1. Compute average shipped delay when `shipped_at` can be `NULL`. State which rows enter `AVG`.
2. Show `COUNT(DISTINCT coupon_id)` versus the number of groups from `GROUP BY coupon_id` (the `NULL` group differs).
3. Use `COALESCE(SUM(qty), 0)` in a left-join report of products with no sales.

#### Advanced practical tasks

1. Write a one-page policy: missing measurements as `NULL`, true zero as `0`, and how each aggregate must be read.
2. Build a report that shows a `NULL` coupon group with the label `none` using `COALESCE` on the label only.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the logical order: `WHERE`, `GROUP BY`, aggregates, `HAVING`, `SELECT`, `ORDER BY`.
2. How do you keep the grain clear in a grouped join of customers and orders?
3. How do `NULL` values change `SUM`, `AVG`, and `GROUP BY`?
4. When do you wrap an aggregate in `COALESCE` for a report?
5. A teammate puts `COUNT(*) > 1` in `WHERE`. Which clause do they need, and why?

#### Easy practical tasks

1. Write a cheat sheet: five aggregates, `GROUP BY`, `HAVING` versus `WHERE`, `NULL` rules.
2. Count rows per group on one practice table. Order by the count.
3. Filter groups with `HAVING`.
4. Compute `MIN` and `MAX` of a date column.

#### Medium practical tasks

1. Write a customer report: order count, first order date, last order date. Sort by count.
2. Join lines and products. Sum qty per product. Keep products with total qty at least 5.
3. Add a grand total row with `UNION ALL` of a grouped query and one total query.

#### Advanced practical tasks

1. Write a monthly sales report (year-month, count, sum) plus a yearly subtotal. Document the method.
2. Compare three ways to list customers with no orders in an aggregate context. Pick one and justify it.
