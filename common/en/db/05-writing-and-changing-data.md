# 5. Writing and Changing Data

## Description

This topic shows how you add rows, change rows, and remove rows. You learn `INSERT`, `UPDATE`, `DELETE`, and `TRUNCATE`. You also learn the upsert idea and a first look at transactions around writes.

Use one term for each concept. A write is a statement that changes stored data. A transaction groups writes so that they all succeed or all fail. Complete this topic before you study isolation in topic 9.

Practice on a copy of data. Do not run unchecked `UPDATE` or `DELETE` on a shared database.

---

## `INSERT` (single and multi-row)

`INSERT` adds new rows to a table. You name the table and the columns. You give values that match those columns.

```sql
INSERT INTO customers (customer_id, name, email)
VALUES (1, 'Ada', 'ada@example.com');
```

Name the columns. Do not rely on the physical column order. A later `ALTER TABLE` can change that order.

Insert more than one row in one statement:

```sql
INSERT INTO customers (customer_id, name, email)
VALUES
    (2, 'Alan', 'alan@example.com'),
    (3, 'Grace', 'grace@example.com');
```

The multi-row form is one statement. In a transaction it succeeds as a whole or it fails as a whole.

You can insert from a query:

```sql
INSERT INTO customers_archive (customer_id, name, email)
SELECT customer_id, name, email
FROM customers
WHERE customer_id < 100;
```

The `SELECT` list must match the insert column list in count and in types.

If a column has a `DEFAULT` or an identity generator, you can omit that column. The DBMS fills it. If a column is `NOT NULL` and has no default, you must supply a value.

An insert fails when it breaks a constraint: duplicate primary key, bad foreign key, `NOT NULL`, or `CHECK`. The failed statement does not add the bad row. In a multi-row insert, many products reject the whole statement.

Do not insert user-controlled text without parameters in a program. Use placeholders. This handbook uses literals only for learning.

### Questions

#### Theoretical questions

1. Why must you name columns in `INSERT`?
2. What does a multi-row `VALUES` list do?
3. When do you use `INSERT ... SELECT`?
4. What happens when an insert breaks a primary key?
5. When may you omit a column from the insert list?

#### Easy practical tasks

1. Insert one customer with a named column list.
2. Insert two products in one `INSERT`.
3. Omit a column that has a default. Select the row and show the default.
4. Write an `INSERT` that must fail on a duplicate key. Run it. Record the error.

#### Medium practical tasks

1. Create `customers_archive`. Copy two rows with `INSERT ... SELECT`.
2. Insert a child row before the parent row. Record the foreign-key error. Fix the order.
3. Compare one statement with three rows versus three statements with one row. Write one difference for errors.

#### Advanced practical tasks

1. Write a script that inserts a customer, an order, and two lines in a safe order. State the key dependencies.
2. Document how your DBMS reports `INSERT` row counts. Show a successful multi-row insert and a failed one.

---

## `UPDATE` and the danger of a missing `WHERE`

`UPDATE` changes existing rows. You set columns. You filter with `WHERE`.

```sql
UPDATE customers
SET email = 'ada.lovelace@example.com'
WHERE customer_id = 1;
```

The `SET` clause assigns new values. You can assign more than one column:

```sql
UPDATE products
SET price = 12.50, name = 'Notebook A5'
WHERE product_id = 10;
```

**Danger.** An `UPDATE` without `WHERE` changes every row in the table.

```sql
UPDATE products
SET price = 0;
```

This statement sets every product price to `0`. That change is a common accident. Always write `WHERE` first in your mind. Many teams write the `SELECT` with the same `WHERE` and check the row count. Then they change `SELECT` to `UPDATE`.

```sql
SELECT product_id, price
FROM products
WHERE product_id = 10;
```

Confirm one row. Then run the `UPDATE` with the same `WHERE`.

You can update from an expression:

```sql
UPDATE products
SET price = price * 1.10
WHERE category = 'book';
```

This statement raises book prices by ten percent.

An update fails when the new values break a constraint. Example: a new email that is not unique.

Some products support `UPDATE ... FROM` or a join in `UPDATE`. Syntax differs. The portable idea is: change rows that you can name with a key. Prefer `WHERE customer_id = ...` or `WHERE customer_id IN (...)`.

Do not run a bulk `UPDATE` on production without a transaction and a tested `WHERE`. Topic 9 shows how you roll back.

### Questions

#### Theoretical questions

1. What does `UPDATE` change?
2. What happens when `UPDATE` has no `WHERE`?
3. Why do teams run a `SELECT` with the same `WHERE` first?
4. What can stop an `UPDATE` after you write new values?
5. Why is a key in `WHERE` safer than a loose text match?

#### Easy practical tasks

1. Update one customer email by `customer_id`.
2. Write the `SELECT` that you run before that update.
3. Update two columns of one product in one statement.
4. Write (but do not run on shared data) an `UPDATE` that would zero all prices. Explain the danger in three sentences.

#### Medium practical tasks

1. Raise prices for one category with an expression. Select before and after.
2. Try an update that breaks a unique email. Record the error.
3. Count rows that match a `WHERE`. Confirm the count before you update.

#### Advanced practical tasks

1. Write a safe checklist for a bulk price change: select, count, transaction, update, verify, commit.
2. Compare vendor `UPDATE ... FROM` docs with a portable `WHERE id IN (SELECT ...)`. Write one portable statement.

---

## `DELETE` vs `TRUNCATE`

`DELETE` removes rows that match `WHERE`. The table stays. The schema stays.

```sql
DELETE FROM orders
WHERE order_id = 500;
```

A `DELETE` without `WHERE` removes every row. That statement is as dangerous as an `UPDATE` without `WHERE`. Confirm the filter with a `SELECT` first.

`DELETE` fires per-row actions that the DBMS defines: foreign-key checks, triggers, and write-ahead logging of each row. A large `DELETE` can take a long time and can write a large log.

`TRUNCATE` removes all rows from a table in one storage operation. It is a schema-level reset of the data. Syntax:

```sql
TRUNCATE TABLE products;
```

`TRUNCATE` does not take a `WHERE` in standard use. It is "empty this table." Many products refuse `TRUNCATE` when other tables reference the table, unless you add extra clauses. Those clauses differ by vendor. Do not truncate a parent table without a plan for the children.

Compare:

| Action | Typical use | Filter | Rollback | Notes |
| --- | --- | --- | --- | --- |
| `DELETE` | Remove some or all rows | `WHERE` allowed | Yes, in a transaction | Checks each row |
| `TRUNCATE` | Empty a table fast | No `WHERE` | Product-dependent | Resets storage; watch foreign keys |

Some products treat `TRUNCATE` as a transaction-safe statement. Others have limits. Test your DBMS in a practice database.

Do not use `TRUNCATE` when you must keep some rows. Use `DELETE` with `WHERE`. Do not use `DROP TABLE` when you only need to remove rows. `DROP TABLE` removes the table object (topic 8).

Child rows can block `DELETE` of a parent. The DBMS protects referential integrity. Delete children first, or use a planned cascade rule.

### Questions

#### Theoretical questions

1. What does `DELETE` remove?
2. What does `TRUNCATE` remove?
3. Why is `DELETE` without `WHERE` dangerous?
4. Why can `TRUNCATE` fail on a parent table?
5. What is the difference between `TRUNCATE` and `DROP TABLE`?

#### Easy practical tasks

1. Insert a throwaway row. Delete it by primary key.
2. Write the `SELECT` that you run before that delete.
3. Write a `TRUNCATE` statement for a scratch table. Run it only on that scratch table.
4. Make a three-row comparison table: `DELETE`, `TRUNCATE`, `DROP TABLE`.

#### Medium practical tasks

1. Fill a scratch table with ten rows. Delete a subset. Then truncate the rest.
2. Try `DELETE` on a parent that still has children. Record the error. Delete in a safe order.
3. Try `TRUNCATE` on a referenced parent. Record the product message.

#### Advanced practical tasks

1. Time a large `DELETE` versus `TRUNCATE` on a scratch table (thousands of rows). Write the two times and one reason.
2. Write a one-page rule: when the team allows `TRUNCATE` in production. Include foreign keys and backup.

---

## `MERGE` / upsert idea (vendor syntax differs)

An upsert is a write that inserts a row when the key is new and updates the row when the key exists. The business idea is "put this row; do not fail on a duplicate key."

The SQL standard includes `MERGE`. Many products implement `MERGE`. Other products use a different form:

- `INSERT ... ON CONFLICT` (PostgreSQL)
- `INSERT ... ON DUPLICATE KEY UPDATE` (MySQL)
- `MERGE` (several products)

Learn the idea first. Then read the syntax for your DBMS. Do not copy a PostgreSQL `ON CONFLICT` clause into a product that does not have it.

Standard-shaped `MERGE` (shape only; test your product):

```sql
MERGE INTO products AS t
USING (
    SELECT 10 AS product_id, 'Notebook' AS name, 4.50 AS price
) AS s
ON t.product_id = s.product_id
WHEN MATCHED THEN
    UPDATE SET name = s.name, price = s.price
WHEN NOT MATCHED THEN
    INSERT (product_id, name, price)
    VALUES (s.product_id, s.name, s.price);
```

The `ON` clause is the match rule. When a row matches, the statement updates. When no row matches, the statement inserts.

Upsert needs a unique key. The match must be unambiguous. If two target rows match, `MERGE` must fail.

A naive program that `SELECT`s, then `INSERT`s or `UPDATE`s in two steps, can lose to a concurrent session. Two sessions can both see "missing" and both insert. Use one upsert statement or a transaction with a proper lock (topic 9).

Do not use upsert to hide a modeling error. If you never expect a duplicate, a plain `INSERT` plus an error is clearer.

### Questions

#### Theoretical questions

1. What does an upsert do?
2. Why does upsert need a unique key?
3. Why do vendor syntaxes differ?
4. What is the role of the `ON` clause in `MERGE`?
5. Why can `SELECT` then `INSERT` fail under two sessions?

#### Easy practical tasks

1. Write in six sentences the upsert idea without product syntax.
2. Find the official upsert or `MERGE` page for your DBMS. Write the statement name.
3. Draw a flow: match on `product_id` → update; no match → insert.
4. List two cases: one that needs upsert, one that must fail on duplicate.

#### Medium practical tasks

1. Run the product upsert twice with the same key and two prices. Show that the second run updates.
2. Try a `MERGE` or upsert that can match two rows. Record the error if the product reports it.
3. Write a portable fallback: transaction, select by key, then insert or update. State the concurrency limit.

#### Advanced practical tasks

1. Compare `MERGE` and the product-specific upsert on the official docs. Write a one-page syntax map.
2. Write a load script idea for a reference table (`currencies`) that you refresh every day with upsert.

---

## Transactions around writes (preview of section 9)

A transaction is a group of statements that the DBMS treats as one change. All statements commit, or all statements roll back. No partial group remains.

You start a transaction, you write, you commit or you roll back.

```sql
BEGIN;

INSERT INTO orders (order_id, customer_id, ordered_at)
VALUES (100, 1, CURRENT_TIMESTAMP);

INSERT INTO order_lines (order_id, product_id, qty)
VALUES (100, 10, 2);

COMMIT;
```

If the second insert fails, you `ROLLBACK`. The order header must not stay without lines if your rule requires lines.

```sql
BEGIN;

UPDATE accounts SET balance = balance - 50 WHERE account_id = 1;
UPDATE accounts SET balance = balance + 50 WHERE account_id = 2;

ROLLBACK;
```

After `ROLLBACK`, both balances return to the start of the transaction. This example is the preview. Topic 9 covers isolation, locks, and dirty reads.

Many clients run in autocommit mode. Each statement commits alone. A failed second statement does not undo the first. For a multi-step write, turn off autocommit or start an explicit transaction.

Do not leave a transaction open while you wait for a human. Locks can block other sessions. Commit or roll back as soon as the writes are complete.

Do not commit before you verify a bulk change if you still might undo it. Verify inside the transaction, then commit. If the client disconnects, many products roll back an open transaction.

`SELECT` can run inside a transaction. You can read the rows that you just wrote before you commit. Other sessions may not see those rows yet. Topic 9 explains why.

### Questions

#### Theoretical questions

1. What does a transaction guarantee for a group of writes?
2. What does `COMMIT` do?
3. What does `ROLLBACK` do?
4. What is the risk of autocommit for a two-step insert?
5. Why must you not leave a transaction open for a long time?

#### Easy practical tasks

1. Start a transaction. Insert one scratch row. Roll back. Confirm that the row is gone.
2. Start a transaction. Insert one scratch row. Commit. Confirm that the row stays.
3. Write the three statements `BEGIN`, `COMMIT`, `ROLLBACK` in your product dialect if the words differ.
4. Explain in four sentences why an order and its lines belong in one transaction.

#### Medium practical tasks

1. Insert a parent and a child in one transaction. Fail the child on purpose. Roll back. Confirm that the parent is gone.
2. With autocommit on, run a good insert and a bad insert. Show that the first row remains.
3. Update two related rows in one transaction. Select them before commit. Then roll back.

#### Advanced practical tasks

1. Write a money-transfer script: begin, two updates, a check, commit or rollback. Do not use this as a banking product.
2. Document autocommit in your GUI and your CLI. Write how you start an explicit transaction in each.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a safe path to change one row: read, write, verify.
2. When do you choose `DELETE`, `TRUNCATE`, or `DROP TABLE`?
3. How does a transaction change the meaning of a multi-statement insert of an order?
4. What must exist before an upsert is a safe single statement?
5. A teammate runs `UPDATE products SET price = 9.99` and then disconnects. Which two habits would have reduced the damage?

#### Easy practical tasks

1. Write a cheat sheet: `INSERT`, multi-row insert, `UPDATE`, `DELETE`, `TRUNCATE`, `MERGE`/upsert, `BEGIN`/`COMMIT`/`ROLLBACK`.
2. Insert two rows, update one, delete one. Show `SELECT` after each step.
3. Run one rollback demo on a scratch table.
4. Write a `WHERE` checklist for `UPDATE` and `DELETE`.

#### Medium practical tasks

1. Load a lookup table with a multi-row insert. Change one label. Remove one unused code (if no child points to it).
2. Build an order header and lines in one transaction. Commit. Then delete in a safe order in a second transaction.
3. Try an upsert twice on the same key. Show insert then update.

#### Advanced practical tasks

1. Write a practice script that transfers stock from one warehouse row to another with a transaction and a rollback path.
2. Compare row-by-row `DELETE` of all rows with `TRUNCATE` on a scratch table. Write when each is correct.
