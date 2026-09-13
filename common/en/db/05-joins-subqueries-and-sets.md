# 5. Joins, Subqueries, and Sets

## Description

This topic shows how you combine rows from more than one table. You learn inner joins, outer joins, cross joins, and self-joins. You also learn set operators and subqueries.

Use one term for each concept. A join matches rows from two tables on a condition. A set operator stacks two query results. A subquery is a query inside another statement. Complete this topic before you group rows.

Qualify column names when two tables share a name. Use table aliases.

---

## Inner join

An inner join returns rows that match on the join condition. A row that has no match does not appear.

```sql
SELECT o.order_id, c.name, o.ordered_at
FROM orders AS o
INNER JOIN customers AS c
    ON o.customer_id = c.customer_id;
```

`INNER JOIN` is the same as `JOIN` in standard SQL. Write `INNER JOIN` while you learn. The word is clear.

The `ON` clause is the match rule. Most joins match a foreign key to a primary key. That match is the usual and safe pattern.

If one customer has three orders, the result has three rows. The customer name repeats. That repeat is correct. The grain of the result is the order, not the customer.

If an order has a `customer_id` that does not exist, an inner join drops that order. A foreign key should prevent that case. If you see dropped orders, check the data and the join condition.

```text
customers          orders
id  name           id  customer_id
1   Ada            10  1
2   Alan           11  1
                   12  2

INNER JOIN on customer_id
→ (10, Ada), (11, Ada), (12, Alan)
```

Do not omit the `ON` clause. A join without a condition becomes a cross join in some dialects. That result is a row explosion.

Do not join on the wrong column. `ON o.order_id = c.customer_id` can run and still be meaningless. Check names.

### Questions

#### Theoretical questions

1. What rows does an inner join keep?
2. What happens to a row that has no match?
3. Why can a customer name appear more than once in the result?
4. What is the usual `ON` pattern?
5. Why must you not omit `ON`?

#### Easy practical tasks

1. Join `orders` to `customers`. Select order id and customer name.
2. Draw two small tables and the inner-join result.
3. Write the same join with aliases `o` and `c`.
4. Explain in four sentences the grain of that result.

#### Medium practical tasks

1. Insert a customer with no orders. Run the inner join. Confirm that the customer is absent.
2. Join `order_lines` to `products`. Select product name and quantity.
3. Break the `ON` clause with the wrong pair of columns. Write why the result is wrong even if rows return.

#### Advanced practical tasks

1. Join three tables: customers, orders, and lines. Select one line-level report. State the grain.
2. Compare `INNER JOIN` with a `WHERE` comma-style join if your DBMS still accepts it. Write why `JOIN` plus `ON` is clearer.

---

## Left, right, and full outer join

An outer join keeps unmatched rows from one or both sides. Columns from the side that has no match are `NULL`.

**Left outer join.** Keep every row from the left table. Add matching rows from the right table. When no match exists, right-hand columns are `NULL`.

```sql
SELECT c.customer_id, c.name, o.order_id
FROM customers AS c
LEFT OUTER JOIN orders AS o
    ON o.customer_id = c.customer_id;
```

A customer with no orders still appears. `order_id` is `NULL` for that row.

**Right outer join.** Keep every row from the right table. This join is a left join with the tables swapped. Prefer `LEFT OUTER JOIN` and write the required table on the left. Teams read left joins more easily.

**Full outer join.** Keep unmatched rows from both sides. A customer without orders appears. An order without a customer appears. Many products support `FULL OUTER JOIN`. Some products do not. Check your DBMS.

```sql
SELECT c.customer_id, o.order_id
FROM customers AS c
FULL OUTER JOIN orders AS o
    ON o.customer_id = c.customer_id;
```

`LEFT JOIN` is a common short form of `LEFT OUTER JOIN`. Write `OUTER` while you learn if it helps.

A filter in `WHERE` on a right-hand column can turn a left join into an inner join.

```sql
SELECT c.name, o.order_id
FROM customers AS c
LEFT OUTER JOIN orders AS o
    ON o.customer_id = c.customer_id
WHERE o.ordered_at >= DATE '2026-01-01';
```

Rows with `o.ordered_at` `NULL` fail this `WHERE`. Customers without orders disappear. Put filters on the optional table in the `ON` clause when you must keep unmatched left rows.

```sql
LEFT OUTER JOIN orders AS o
    ON o.customer_id = c.customer_id
   AND o.ordered_at >= DATE '2026-01-01';
```

### Questions

#### Theoretical questions

1. What does a left outer join keep from the left table?
2. What value appears in right-hand columns when there is no match?
3. Why do teams prefer left join over right join?
4. What does a full outer join keep?
5. Why can a `WHERE` on a right-hand column remove unmatched left rows?

#### Easy practical tasks

1. List all customers and their order ids with a left join. Include customers with no orders.
2. Draw the left-join result for the small tables in the inner-join section plus a customer with no orders.
3. Write a query that finds customers with no orders (`order_id IS NULL` after a left join).
4. Explain in four sentences when you need a left join instead of an inner join.

#### Medium practical tasks

1. Run the left join with a date filter in `WHERE` and again in `ON`. Compare the customer list.
2. Try `FULL OUTER JOIN` on your DBMS. Record whether it works.
3. Left-join `products` to `order_lines`. List products that never sold.

#### Advanced practical tasks

1. Write a report: every customer, count of orders (zero when none). You may use a later `GROUP BY` idea in words if you do not know the syntax yet, or use one row per order and count by hand on a tiny set.
2. Document `FULL OUTER JOIN` alternatives with `UNION` of a left join and a right join. Write the idea.

---

## Cross join and self-join

A cross join pairs every row of one table with every row of the other table. The result size is the product of the two sizes.

```sql
SELECT c.name, p.name AS product_name
FROM customers AS c
CROSS JOIN products AS p;
```

If you have 100 customers and 500 products, the result has 50000 rows. Each pair is a row.

A cross join is correct when you need the full combination. Example: every size in a `sizes` table with every color in a `colors` table to build a catalog grid.

A cross join is a mistake when you meant an inner join and you forgot `ON`. Some old syntaxes make this easy:

```sql
SELECT *
FROM orders, customers;
```

This form is a cross join. A later `WHERE orders.customer_id = customers.customer_id` turns it into an inner join. If you forget the `WHERE`, the client returns a huge grid. Use `INNER JOIN ... ON` so that the match is required.

Do not use a cross join to "see if something matches" on large tables. The DBMS may build a huge intermediate result.

Test a cross join only on small lookup tables. Cancel the query if the row count is unexpected.

```text
sizes: S, M
colors: red, blue
CROSS JOIN → (S,red), (S,blue), (M,red), (M,blue)
```

A self-join is a join of a table to itself. You use two aliases. Each alias is one role.

Example. Employees with a manager in the same table:

```sql
SELECT e.name AS employee_name, m.name AS manager_name
FROM employees AS e
INNER JOIN employees AS m
    ON e.manager_id = m.employee_id;
```

`e` is the employee. `m` is the manager. The join matches the foreign key `manager_id` to `employee_id`.

A left self-join keeps employees who have no manager (the top person):

```sql
SELECT e.name AS employee_name, m.name AS manager_name
FROM employees AS e
LEFT OUTER JOIN employees AS m
    ON e.manager_id = m.employee_id;
```

You must alias the table twice. Without aliases, the DBMS cannot tell the two roles apart.

Do not create a loop in your head as a recursive walk unless you use a recursive query. A single self-join walks one step (employee to manager). A full chain needs more joins or a recursive common table expression.

```text
employees
id  name   manager_id
1   Ada    NULL
2   Alan   1
3   Grace  1

self-join → (Alan, Ada), (Grace, Ada)
left self-join also → (Ada, NULL)
```

### Questions

#### Theoretical questions

1. How many rows does a cross join of 10 and 20 rows produce?
2. When is a cross join the correct tool, and when is it a mistake?
3. What is a self-join, and why do you need two aliases?
4. What does `ON e.manager_id = m.employee_id` mean?
5. How many management steps does one self-join walk?

#### Easy practical tasks

1. Cross-join two tiny lookup tables (two rows each). Show all four pairs.
2. Write the row-count formula for a cross join.
3. Create a small `employees` table with a manager column. Insert three rows. Write an inner self-join that lists employee and manager names.
4. Write a left self-join that includes the top manager.

#### Medium practical tasks

1. Build a size-and-color grid with `CROSS JOIN`. Insert the result into a scratch table if you want.
2. Find pairs of customers in the same city with a self-join. Avoid pairing a row with itself (`a.customer_id < b.customer_id`).
3. Forget `ON` on a join of `orders` and `customers` (or use comma syntax). Stop the query if it grows. Write the lesson.

#### Advanced practical tasks

1. Write two self-joins to show employee, manager, and manager of manager. State the limit of this approach.
2. Read a short overview of recursive queries. Write five sentences on when you need them instead of one self-join.

---

## `UNION`, `UNION ALL`, `INTERSECT`, `EXCEPT`

Set operators combine two query results. Each query must have the same number of columns. Types must be compatible. Column names come from the first query.

**`UNION`.** Stack the rows. Remove duplicates.

**`UNION ALL`.** Stack the rows. Keep duplicates.

```sql
SELECT name FROM staff
UNION
SELECT name FROM contractors;
```

```sql
SELECT name FROM staff
UNION ALL
SELECT name FROM contractors;
```

If `Ada` appears in both tables, `UNION` returns one `Ada`. `UNION ALL` returns two.

**`INTERSECT`.** Return rows that appear in both results.

**`EXCEPT`.** Return rows in the first result that do not appear in the second. Some products use `MINUS` for this idea. Prefer `EXCEPT` when the product has it.

```sql
SELECT customer_id FROM orders
INTERSECT
SELECT customer_id FROM newsletter;

SELECT customer_id FROM customers
EXCEPT
SELECT customer_id FROM banned_customers;
```

Set operators compare full rows. `SELECT a, b UNION SELECT a, b` uses both columns.

Use `UNION ALL` when you know the sets do not overlap, or when you must keep duplicates. `UNION` sorts or hashes to remove duplicates. That extra work has a cost.

Both sides should mean the same grain. Do not union a customer list with an order list unless you only want a bag of ids and you label the source.

```sql
SELECT customer_id, 'order' AS source
FROM orders
UNION ALL
SELECT customer_id, 'newsletter' AS source
FROM newsletter;
```

Not every product implements `INTERSECT` and `EXCEPT`. Check the manual. You can express some of these ideas with joins and `EXISTS`.

### Questions

#### Theoretical questions

1. What extra work does `UNION` do that `UNION ALL` does not do?
2. What does `INTERSECT` return?
3. What does `EXCEPT` return?
4. What must be true of the two sides of a set operator?
5. Why do you add a `source` column in a stacked report?

#### Easy practical tasks

1. Union two small name lists. Compare `UNION` and `UNION ALL`.
2. Write an `EXCEPT` (or the product equivalent) that subtracts one id list from another.
3. Write an `INTERSECT` of two id lists.
4. Explain in four sentences when `UNION ALL` is the better default.

#### Medium practical tasks

1. Stack current customers and archived customers with a source label.
2. Try `EXCEPT` on your DBMS. If it fails, write the join/`NOT EXISTS` idea.
3. Union two queries with a different column count. Record the error.

#### Advanced practical tasks

1. Rewrite an `INTERSECT` as an inner join on the full row. Show the same ids on a tiny set.
2. Write a one-page map of set operators versus joins for "in A and B", "in A or B", "in A not B".

---

## Subqueries: scalar, `IN`, `EXISTS`

A subquery is a `SELECT` inside another statement. Parentheses wrap it.

**Scalar subquery.** The subquery returns one column and at most one row. You can use it in a `SELECT` list or in a comparison.

```sql
SELECT name,
       (SELECT COUNT(*) FROM orders AS o WHERE o.customer_id = c.customer_id) AS order_count
FROM customers AS c;
```

If a scalar subquery returns two rows, the statement fails. If it returns no row, the value is `NULL`.

**`IN` subquery.** The outer row matches when the value equals one value from the subquery.

```sql
SELECT name
FROM customers
WHERE customer_id IN (
    SELECT customer_id
    FROM orders
    WHERE ordered_at >= DATE '2026-01-01'
);
```

**`EXISTS` subquery.** The outer row matches when the subquery returns at least one row. You usually correlate: the inner query refers to the outer row.

```sql
SELECT c.name
FROM customers AS c
WHERE EXISTS (
    SELECT 1
    FROM orders AS o
    WHERE o.customer_id = c.customer_id
);
```

`EXISTS` stops when it finds one row. `SELECT 1` is a common dummy list. The columns do not matter.

`NOT EXISTS` finds outer rows with no matching inner row. This form is a clear way to find customers without orders.

```sql
SELECT c.name
FROM customers AS c
WHERE NOT EXISTS (
    SELECT 1
    FROM orders AS o
    WHERE o.customer_id = c.customer_id
);
```

Prefer `EXISTS` over `IN` when the inner query can return `NULL` in a `NOT IN` trap. Prefer a join when you need columns from both tables in the result. Prefer a subquery when you only filter.

A correlated subquery refers to the outer query. An uncorrelated subquery does not. The DBMS can run an uncorrelated subquery once.

### Questions

#### Theoretical questions

1. What is a subquery?
2. What must a scalar subquery return?
3. What does `EXISTS` test?
4. When do you use `NOT EXISTS`?
5. What is a correlated subquery?

#### Easy practical tasks

1. Select customers whose id appears in `orders` with `IN`.
2. Rewrite that query with `EXISTS`.
3. Find customers with no orders with `NOT EXISTS`.
4. Write a scalar subquery that returns one customer name for a known id.

#### Medium practical tasks

1. Compare `IN`, `EXISTS`, and an inner join for "customers who ordered". Write which columns each form can return.
2. Force a scalar subquery to return two rows. Record the error.
3. Show `NOT IN` with a `NULL` in the inner result. Compare `NOT EXISTS` on the same data.

#### Advanced practical tasks

1. Write a query that lists each product and the latest order date with a scalar subquery. State the one-row rule.
2. Read a query plan for `EXISTS` versus `IN` on your DBMS. Write three observations.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. When do you choose inner join, left join, or `NOT EXISTS` for customers and orders?
2. How do you keep a left join from becoming an inner join by accident?
3. When is a cross join correct, and when is a self-join the right tool?
4. When is `UNION ALL` correct and `UNION` wasteful?
5. A teammate writes `FROM a, b` and forgets `WHERE`. Which join did they write, and what is the fix?

#### Easy practical tasks

1. Write a cheat sheet: inner, left, right, full, cross, self-join, set operators, subquery kinds.
2. Run one inner join and one left join on your practice tables.
3. Run `UNION ALL` of two id lists.
4. Run `EXISTS` to test that a customer has at least one order.

#### Medium practical tasks

1. Build a report: each customer name, order id or `NULL`, product names via a second join. State the grain.
2. Find unused products with a left join and with `NOT EXISTS`. Confirm the same ids.
3. Self-join a small hierarchy or a same-city pair query.

#### Advanced practical tasks

1. Write a portable "customers in A not in B" in three ways: `EXCEPT`, `NOT EXISTS`, left join plus `IS NULL`. Run the forms that your DBMS allows.
2. Design a three-table report with one outer join and one inner join. Document why each join type is required.
