# 4. SQL Fundamentals

## Description

SQL is the language that you use to read relational data. This topic shows `SELECT`, filters, sort, limit, and aliases. You learn projection and selection as two different operations.

Use one term for each concept. A query is a `SELECT` statement. Projection picks columns. Selection picks rows. Complete this topic before you write joins and aggregates.

Write vendor-neutral SQL. Some products use `LIMIT`. The SQL standard also uses `FETCH FIRST`. Both appear in this topic. Do not mix dialects in one statement.

---

## What SQL is (declarative query language)

SQL is a declarative language. You state the result that you want. You do not state the scan algorithm. The DBMS chooses an access path. Topic 11 shows how you read that choice.

SQL covers more than queries. The language has parts:

- data query: `SELECT`
- data change: `INSERT`, `UPDATE`, `DELETE`
- schema: `CREATE`, `ALTER`, `DROP`
- transaction: `BEGIN`, `COMMIT`, `ROLLBACK`

This topic covers query. Later topics cover the other parts.

A declarative statement names tables, columns, and conditions. Example:

```sql
SELECT name, email
FROM customers
WHERE email IS NOT NULL;
```

You do not write a loop. You do not write "open the file." The DBMS finds the rows.

SQL is a standard (ISO/IEC). Products extend the standard. Learn standard forms first. Check the product manual when a statement fails.

Do not build SQL by concatenating user text. Use parameters in programs. This path shows literal values in examples for learning only.

SQL is not a general programming language. It has limited flow control in the core language. Keep application loops in the application. Keep set operations in SQL.

### Questions

#### Theoretical questions

1. What does declarative mean for a `SELECT`?
2. Who chooses the access path for a query?
3. Name four parts of SQL besides a simple query.
4. Why do products differ if SQL is a standard?
5. Why must a program not concatenate user text into SQL?

#### Easy practical tasks

1. Write five sentences that describe SQL. Use only facts from this section.
2. Label each statement as query, change, schema, or transaction: `SELECT`, `INSERT`, `CREATE TABLE`, `COMMIT`.
3. Rewrite a loop in words ("for each customer, print email") as a `SELECT` sketch.
4. Open the SQL page of your DBMS. Write one extension that is not in this handbook.

#### Medium practical tasks

1. Compare SQL with one language that you know. Write six short sentences. Cover loops, types, and who runs the code.
2. Run `SELECT 1`. Then run a `SELECT` from a real table. Write what is declarative in both.
3. Find the ISO SQL Wikipedia overview. Write three sentence fragments that match this section.

#### Advanced practical tasks

1. Write a one-page note: what SQL declares, and what the optimizer still decides. Use one example query.
2. List five vendor extensions in your DBMS (example: `LIMIT`, extra types). Mark each as standard-like or product-only.

---

## `SELECT`, `FROM`, `WHERE`

A basic query has three clauses.

`FROM` names the table (or tables) that you read. The DBMS starts from that table.

`SELECT` names the columns that you return. You can use `*` to return all columns. Use `*` only for exploration. Name columns in saved queries.

`WHERE` keeps the rows that satisfy a condition. Rows that fail the condition do not appear.

```sql
SELECT customer_id, name
FROM customers
WHERE name = 'Ada';
```

Evaluation idea (logical):

1. Start with the rows in `FROM`.
2. Keep rows that match `WHERE`.
3. Return the columns in `SELECT`.

The product can use another physical order. The result must match this meaning.

`WHERE` is optional. Without `WHERE`, the query returns every row. That result can be large. Add a filter when you test on real data.

```sql
SELECT *
FROM customers;
```

Use a semicolon to end a statement in most clients. Some clients send one statement at a time and do not need it. Use the semicolon in scripts.

Do not put a column alias filter in `WHERE` if the product does not allow it. `WHERE` applies before `SELECT` list aliases in standard SQL. Use the column name in `WHERE`.

### Questions

#### Theoretical questions

1. What does `FROM` name?
2. What does `SELECT` name?
3. What does `WHERE` do to rows?
4. What is the logical order of `FROM`, `WHERE`, and `SELECT`?
5. Why is `SELECT *` a poor habit in saved queries?

#### Easy practical tasks

1. Select two columns from a table that you created. Show the statement.
2. Add a `WHERE` that matches one row. Show the result.
3. Run the same query without `WHERE`. Write how many rows you get.
4. Rewrite `SELECT *` as a list of real column names.

#### Medium practical tasks

1. Create `customers` with three rows. Select only `Ada`. Then select everyone except `Ada`.
2. Run a `WHERE` on a column that does not exist. Record the error. Fix the name.
3. Compare `SELECT *` and a named list in the client. Write one risk of `*` when a column is added later.

#### Advanced practical tasks

1. Write three queries against `orders`: all rows, one customer, and a date filter. Use named columns.
2. Document whether your client needs a semicolon. Write a two-statement script that works in that client.

---

## Projection vs selection

Projection and selection are two operations. Do not mix the words.

**Projection** picks attributes (columns). `SELECT name, email` is a projection. The result has fewer columns. The number of rows can stay the same.

**Selection** picks tuples (rows). `WHERE price > 10` is a selection. The result has fewer rows. The number of columns can stay the same.

A typical query does both:

```sql
SELECT name, price
FROM products
WHERE price > 10;
```

This query selects rows with `price > 10`. It projects `name` and `price`.

In mathematics, selection uses a condition on the relation. Projection uses a list of attributes. SQL uses `WHERE` for selection and the `SELECT` list for projection. The keyword `SELECT` is not the same as the word *selection*. That clash confuses beginners. Remember: `SELECT` list = projection. `WHERE` = selection.

```text
products (all columns, all rows)
    | selection: price > 10
    v
fewer rows, all columns
    | projection: name, price
    v
fewer rows, two columns
```

Do not call `WHERE` a projection. Do not call a column list a selection.

### Questions

#### Theoretical questions

1. What does projection pick?
2. What does selection pick?
3. Which SQL clause is projection in a simple query?
4. Which SQL clause is selection in a simple query?
5. Why does the keyword `SELECT` confuse the word selection?

#### Easy practical tasks

1. Write a query that is only projection (no `WHERE`).
2. Write a query that is only selection (`SELECT *` plus `WHERE`).
3. Write a query that is both. Label each part.
4. Draw the three-box flow from this section for `employees` and `dept = 'HR'`.

#### Medium practical tasks

1. On `products`, show row counts: no filter, filter only, filter plus two columns. Write the three counts.
2. Explain to a teammate in six sentences why `SELECT` is a bad name for selection.
3. Take a report request in words. Split it into one selection condition and one projection list.

#### Advanced practical tasks

1. Write a one-page note that maps relational algebra names to SQL clauses for project, select, and (preview) join.
2. Show a case where projection creates duplicate rows (`SELECT city FROM customers`). Explain why selection did not cause those duplicates.

---

## `ORDER BY`, `LIMIT` / `FETCH`

A relation has no order. `ORDER BY` gives an order to the result. You name columns or expressions. Use `ASC` for ascending order (default). Use `DESC` for descending order.

```sql
SELECT name, price
FROM products
ORDER BY price DESC, name ASC;
```

Rows sort by `price` from high to low. Rows with the same price sort by `name`.

Always use `ORDER BY` when the user must see a stable order. Without `ORDER BY`, the DBMS can return any order. That order can change after an upgrade or after an index change.

`LIMIT` keeps the first *n* rows of the sorted result. Many products support `LIMIT`. The SQL standard form is `FETCH FIRST`.

```sql
SELECT name, price
FROM products
ORDER BY price DESC
FETCH FIRST 5 ROWS ONLY;
```

```sql
SELECT name, price
FROM products
ORDER BY price DESC
LIMIT 5;
```

Use `ORDER BY` with `LIMIT` or `FETCH`. A limit without order is not a "top 5" result. It is any five rows.

Some products use `TOP (5)` in the `SELECT` list. That form is not standard. Prefer `FETCH FIRST` or the product `LIMIT` in one dialect.

Do not use `LIMIT` as a substitute for a correct `WHERE`. A limit hides rows. It does not define the business set.

### Questions

#### Theoretical questions

1. Why does a query without `ORDER BY` have no guaranteed order?
2. What does `DESC` do?
3. Why must `LIMIT` come with `ORDER BY` for a "top n" report?
4. What is the standard form that many texts use instead of `LIMIT`?
5. Why is `LIMIT` not a replacement for `WHERE`?

#### Easy practical tasks

1. Select all products ordered by `name`.
2. Select the same rows in reverse name order.
3. Return the first three rows by `price` descending. Use `LIMIT` or `FETCH`.
4. Run the limit query without `ORDER BY`. Write why the result is not "the three cheapest".

#### Medium practical tasks

1. Sort `orders` by date, then by `order_id`. Show a tie on the date.
2. Compare `FETCH FIRST 5 ROWS ONLY` and `LIMIT 5` on your DBMS. Write which form works.
3. Write a query for the five most recent orders for one customer.

#### Advanced practical tasks

1. Write a stable pagination idea with `ORDER BY` key columns (no `OFFSET` yet). Explain why the sort keys must be unique.
2. Document `TOP`, `LIMIT`, and `FETCH` for two products. Write one portable recommendation for this path.

---

## `DISTINCT`

`DISTINCT` removes duplicate rows from the result. The DBMS compares the full selected row. Two rows that match in every selected column become one row.

```sql
SELECT DISTINCT city
FROM customers;
```

This query lists each city once. Many customers can live in the same city.

Without `DISTINCT`, the query returns one city per customer row.

```sql
SELECT city
FROM customers;
```

Use `DISTINCT` when the question is about unique combinations. Do not use `DISTINCT` to hide a bad join that multiplied rows. Fix the join first (topic 6).

`SELECT DISTINCT a, b` treats the pair as the row. `(Paris, 10)` and `(Paris, 11)` stay as two rows.

`COUNT(DISTINCT column)` counts unique non-null values in aggregates. Topic 7 covers that function. `DISTINCT` in the `SELECT` list is the row-level idea.

`DISTINCT` can be expensive on a large result. Prefer a correct `WHERE` and a correct join. Add `DISTINCT` when the unique list is the real question.

Some products allow `DISTINCT ON`. That form is not standard. Do not use it in this path.

### Questions

#### Theoretical questions

1. What does `DISTINCT` remove?
2. What columns does `DISTINCT` compare?
3. Why can `SELECT DISTINCT city` return fewer rows than `SELECT city`?
4. Why must you not use `DISTINCT` to hide a bad join?
5. What does `SELECT DISTINCT a, b` treat as one row?

#### Easy practical tasks

1. Insert two customers in the same city. Select `city` with and without `DISTINCT`.
2. Select `DISTINCT city, country` if you have both columns. Show a pair that is not collapsed.
3. Write a query that lists unique product names.
4. Explain in four sentences when `DISTINCT` is the right tool.

#### Medium practical tasks

1. Create rows that differ only in a column that you do not select. Show that `DISTINCT` on the remaining columns collapses them.
2. Compare row counts: `SELECT customer_id` versus `SELECT DISTINCT customer_id` from `orders`.
3. Write a case where `DISTINCT` is a warning sign (join preview). Describe the better fix in words.

#### Advanced practical tasks

1. Measure or estimate a `DISTINCT` on a large text column versus a `GROUP BY`. Write which one you would use for a unique list.
2. Read whether your DBMS supports `DISTINCT ON`. Write why this path avoids it.

---

## Comparison and `AND` / `OR` / `NOT`

`WHERE` uses predicates. A predicate is true, false, or unknown (`NULL`).

Comparison operators:

- `=` equal
- `<>` not equal (also `!=` in many products)
- `<` less than
- `>` greater than
- `<=` less or equal
- `>=` greater or equal

Combine predicates with `AND`, `OR`, and `NOT`.

```sql
SELECT product_id, name, price
FROM products
WHERE price >= 10
  AND price < 50;
```

`AND` requires both sides to be true. `OR` requires one side to be true. `NOT` inverts a predicate.

Use parentheses when you mix `AND` and `OR`. `AND` binds more tightly than `OR` in SQL. Do not rely on memory. Write the parentheses.

```sql
SELECT *
FROM products
WHERE (category = 'book' OR category = 'map')
  AND price < 20;
```

Without parentheses, the meaning can surprise you.

`NOT` applies to a predicate:

```sql
SELECT *
FROM products
WHERE NOT (price < 10);
```

Prefer a clear comparison (`price >= 10`) when it is easier to read.

A comparison with `NULL` is unknown, not true. `WHERE price = NULL` returns no row. Use `IS NULL` (next section).

Use `<>` in standard SQL for not equal. Check that your product accepts `<>`.

### Questions

#### Theoretical questions

1. What are the three truth values in SQL predicates?
2. What does `AND` require?
3. Why do you write parentheses when you mix `AND` and `OR`?
4. What does `WHERE price = NULL` return? Why?
5. Which standard operator means not equal?

#### Easy practical tasks

1. Select products cheaper than `10`.
2. Select products in a price range with `>=` and `<`.
3. Select rows that match one of two names with `OR`.
4. Rewrite a `NOT (price < 10)` filter as a single comparison.

#### Medium practical tasks

1. Write two queries: mixed `AND`/`OR` with parentheses and the same tokens without parentheses. Compare results.
2. Filter `orders` with a date lower bound and a customer id. Use `AND`.
3. Show that `WHERE qty <> 0` does not include rows where `qty` is `NULL`.

#### Advanced practical tasks

1. Write a truth table for `AND`, `OR`, and `NOT` with true, false, and unknown. Use it on a sample `WHERE`.
2. Translate a business rule with three clauses into SQL with parentheses. Ask a teammate to read the meaning aloud.

---

## `IN`, `BETWEEN`, `LIKE`, `IS NULL`

These predicates appear often in `WHERE`.

**`IN`.** The value equals one item in a list.

```sql
SELECT *
FROM products
WHERE category IN ('book', 'map', 'card');
```

This form matches `OR` of equalities. A subquery can appear in `IN` (topic 6).

**`BETWEEN`.** The value sits in a closed range. `x BETWEEN a AND b` means `x >= a AND x <= b`.

```sql
SELECT *
FROM orders
WHERE ordered_at BETWEEN DATE '2026-01-01' AND DATE '2026-01-31';
```

Both ends are included. For dates and times, confirm the end of the last day. A timestamp `2026-01-31 15:00:00` is inside that range. A timestamp on `2026-02-01` is not. For timestamps, a pair of `>=` and `<` is often safer.

**`LIKE`.** The value matches a text pattern. `%` means any string. `_` means one character.

```sql
SELECT *
FROM customers
WHERE email LIKE '%@example.com';
```

`LIKE` is case-sensitive on some products and not on others. Test your DBMS. Do not use `LIKE '%x%'` on large tables as a first plan. That pattern can force a wide scan.

**`IS NULL` and `IS NOT NULL`.** Test missing values.

```sql
SELECT *
FROM orders
WHERE shipped_at IS NULL;
```

Do not write `= NULL` or `<> NULL`. Those predicates are unknown.

`NOT IN` with a list that contains `NULL` is a trap. The result can be empty. Prefer `NOT EXISTS` later, or keep `NULL` out of the list.

### Questions

#### Theoretical questions

1. What does `IN` test?
2. Is `BETWEEN` inclusive on both ends?
3. What do `%` and `_` mean in `LIKE`?
4. Why must you use `IS NULL` instead of `= NULL`?
5. Why is `LIKE '%text%'` a risky default on a large table?

#### Easy practical tasks

1. Select rows with `category IN` of two values.
2. Select prices `BETWEEN 5 AND 15`.
3. Select names that start with `A` using `LIKE`.
4. Select rows where a date column `IS NULL`.

#### Medium practical tasks

1. Compare `BETWEEN` on a `DATE` column and on a `TIMESTAMP` column at month end. Write the surprise if you find one.
2. Use `LIKE` with `_` to match a four-letter code. Show a match and a miss.
3. Build a list for `IN` from a lookup table in words (subquery preview). Write the idea.

#### Advanced practical tasks

1. Demonstrate `NOT IN` with a `NULL` in the list (if you can). Record the result. Write the safe alternative in words.
2. Write a portable date window for one month using `>=` start and `<` next month. Explain why you prefer it to `BETWEEN` for timestamps.

---

## Aliases

An alias is a temporary name in a query. You alias a column or a table. The alias exists only for that statement.

Column aliases name output columns. Use `AS` for clarity.

```sql
SELECT name AS product_name, price AS unit_price
FROM products;
```

The result headers are `product_name` and `unit_price`. The table columns stay `name` and `price`.

Table aliases shorten names and prepare joins.

```sql
SELECT c.customer_id, c.name
FROM customers AS c
WHERE c.name = 'Ada';
```

`c` is the alias for `customers`. After you set an alias, use that alias for columns in many products. Standard SQL allows the table name or the alias. Be consistent.

Use aliases when names are long or when you join the same table twice (self-join in topic 6). Use short names that still mean something (`c` for customers, `o` for orders). Do not use `t1`, `t2` when you can use a letter that matches the table.

Some products let you omit `AS`. Write `AS` in this path. The extra word is clear.

Do not refer to a column alias in `WHERE` in standard SQL. Repeat the expression or use a subquery. Some products allow the alias in `ORDER BY`. Test `ORDER BY product_name` on your DBMS.

```sql
SELECT price * qty AS line_total
FROM order_lines
ORDER BY line_total;
```

If `ORDER BY line_total` fails, write `ORDER BY price * qty`.

### Questions

#### Theoretical questions

1. What is an alias?
2. Does a column alias rename the stored column?
3. Why do table aliases help in joins?
4. Why does this path keep the word `AS`?
5. Why can a column alias fail in `WHERE`?

#### Easy practical tasks

1. Select `name AS customer_name` from `customers`.
2. Add a table alias `c` and qualify every column.
3. Compute `price * qty AS line_total` from a lines table.
4. Write four sentences on when you alias a table.

#### Medium practical tasks

1. Try `WHERE` with a column alias. Record success or error. Write the portable form.
2. Try `ORDER BY` with a column alias. Record the result on your DBMS.
3. Rewrite a query that uses `t1` into aliases that match table names.

#### Advanced practical tasks

1. Write a query with two table aliases that you will later join (preview). Keep every column qualified.
2. Document alias rules for `WHERE`, `GROUP BY`, and `ORDER BY` on your DBMS in a short table.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the logical path of a query that filters, projects, sorts, and limits.
2. How do projection and selection appear in one `SELECT` that also uses `DISTINCT`?
3. When do you choose `IN`, `BETWEEN`, `LIKE`, or `IS NULL`?
4. Why is SQL declarative even when you add `ORDER BY` and `FETCH`?
5. A teammate writes `WHERE alias_name = 1` and `LIMIT 5` without `ORDER BY`. Which two problems do you name?

#### Easy practical tasks

1. Write a cheat sheet: `SELECT`, `FROM`, `WHERE`, `ORDER BY`, `FETCH`/`LIMIT`, `DISTINCT`, predicates, aliases.
2. Against your practice tables, write one query that uses a filter, two columns, a sort, and a limit.
3. Write one query with `IN` and one query with `LIKE`.
4. Add aliases to both queries.

#### Medium practical tasks

1. Build a "search" query for products: name pattern, price range, and category list. Sort by price.
2. Show unique cities with `DISTINCT`. Then show the same list with a sort.
3. Write two equivalent filters: one with `OR` and one with `IN`. Confirm the same rows.

#### Advanced practical tasks

1. Write a portable top-10 report for latest orders. Document `FETCH` versus `LIMIT` on your DBMS.
2. Create a small script of five `SELECT` statements that a beginner can run after a sample insert. Include one `IS NULL` query.
