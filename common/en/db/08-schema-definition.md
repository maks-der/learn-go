# 8. Schema Definition

## Description

This topic shows how you create tables, change tables, and drop tables. You learn common data types at a high level, table constraints, indexes as schema objects, and naming conventions.

Use one term for each concept. A schema definition statement changes objects, not row values. `CREATE`, `ALTER`, and `DROP` are data definition language (DDL). Complete this topic before you study transactions in depth. DDL and transactions interact in product-specific ways.

Practice on a throwaway database. Do not drop a table that other people use.

---

## `CREATE TABLE`, data types (high-level)

`CREATE TABLE` defines a base table. You name the table. You name each column. You give each column a type. You add constraints.

```sql
CREATE TABLE products (
    product_id INTEGER PRIMARY KEY,
    sku VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL
);
```

Common type families (names differ by product):

| Family | Typical standard-like name | Use |
| --- | --- | --- |
| Integer | `INTEGER` | Counts, ids, quantities |
| Decimal | `DECIMAL(p, s)` | Money and exact fractions |
| Floating | `DOUBLE PRECISION` | Scientific measures, not money |
| Text | `VARCHAR(n)`, `CHAR(n)` | Strings; prefer `VARCHAR` |
| Date | `DATE` | Calendar day |
| Time | `TIMESTAMP` | Day plus time |
| Boolean | `BOOLEAN` | True or false |

Use `DECIMAL` (or the product exact numeric type) for money. Do not use a binary floating type for currency. Rounding errors appear.

`VARCHAR(n)` has a maximum length. Pick a length that fits the domain. `VARCHAR(2)` is enough for many country codes. `VARCHAR(200)` fits many names. Some products have unlimited text types. Use them when the length is truly unbounded.

`TIMESTAMP` may or may not include a time zone, depending on the type name. A later topic covers time zones. For this path, store timestamps in one agreed convention and document it.

`BOOLEAN` is not present in every product. Some products use `0`/`1` or `'t'`/`'f'`. Check the manual.

Do not invent a type per product in this path when a standard-like type exists. Write `INTEGER` and `DECIMAL`. Adjust to the product when you implement.

`CREATE TABLE` fails when the name already exists. Use a new name, or drop the old table only if that is safe.

### Questions

#### Theoretical questions

1. What does `CREATE TABLE` define?
2. Why do you use `DECIMAL` for money?
3. What does `VARCHAR(n)` limit?
4. Why do type names differ across products?
5. What happens when you create a table that already exists?

#### Easy practical tasks

1. Create `products` with an integer key, a name, and a decimal price.
2. Create `customers` with a text email column and a timestamp.
3. List the types that your DBMS used when you asked for `INTEGER` and `BOOLEAN`.
4. Write four sentences on why `DOUBLE PRECISION` is a poor money type.

#### Medium practical tasks

1. Create three related tables for a library. Use types that match each attribute.
2. Try `CREATE TABLE` twice with the same name. Record the error.
3. Insert a string that is too long for `VARCHAR(5)`. Record the error.

#### Advanced practical tasks

1. Write a type-choice sheet for a shop: money, qty, flags, dates, codes. Map each to a type on your DBMS.
2. Compare `CHAR` and `VARCHAR` in the product docs. Write six sentences on padding and when `CHAR` is acceptable.

---

## `ALTER TABLE`

`ALTER TABLE` changes an existing table. You add a column, drop a column, change a type, or add a constraint. The exact clauses differ by vendor. Learn the ideas and then read the product syntax.

Add a column:

```sql
ALTER TABLE customers
    ADD email_verified BOOLEAN NOT NULL DEFAULT FALSE;
```

A new `NOT NULL` column needs a default if rows already exist. The DBMS must fill the new column for every old row.

Drop a column:

```sql
ALTER TABLE customers
    DROP COLUMN email_verified;
```

Dropping a column destroys that data. Other objects (views, indexes) can block the drop.

Rename a column or a table (syntax varies):

```sql
ALTER TABLE customers
    RENAME COLUMN name TO full_name;
```

Add a constraint after the table exists:

```sql
ALTER TABLE products
    ADD CONSTRAINT chk_products_price CHECK (price >= 0);
```

The add fails if a current row breaks the check. Clean the data first, or fix the rule.

Change a type when the existing values convert cleanly. A change from `VARCHAR(200)` to `VARCHAR(50)` fails if a value is longer than 50.

Do not use `ALTER TABLE` as a daily data-fix tool. Use `UPDATE` for values. Use `ALTER` for structure.

Some products cannot alter a column in one step. They add a new column, copy data, drop the old column, and rename. Plan that path for large tables.

DDL can lock a table. A long alter can block writers. Run large alters in a maintenance window in production. This path only needs small practice tables.

### Questions

#### Theoretical questions

1. What kinds of change does `ALTER TABLE` make?
2. Why does a new `NOT NULL` column need a default when rows exist?
3. What can block `DROP COLUMN`?
4. Why can a type shrink fail?
5. Why is `ALTER` the wrong tool to fix one email value?

#### Easy practical tasks

1. Add a nullable `phone` column to `customers`.
2. Add a `CHECK` on `products.price`.
3. Drop the `phone` column after you test it.
4. Write four sentences on the default problem for `NOT NULL` adds.

#### Medium practical tasks

1. Add `NOT NULL` to a column that still has `NULL`. Record the error. Update the rows. Retry.
2. Shrink a `VARCHAR` and show a failure when a value is too long.
3. Rename a column. Update one query that used the old name.

#### Advanced practical tasks

1. Write a four-step plan to replace a column type on a large table without one dangerous step. Stay high-level.
2. Read whether your DBMS locks the table for `ADD COLUMN`. Write five sentences from the docs.

---

## `DROP TABLE` and cascading deletes (careful)

`DROP TABLE` removes the table object and its data. The name is gone. Indexes on that table are gone. The action is not an empty table. To keep the table and remove rows, use `DELETE` or `TRUNCATE` (topic 5).

```sql
DROP TABLE scratch_items;
```

If other tables reference the table with a foreign key, the DBMS rejects the drop by default. That rejection is a safety rule.

Some products allow:

```sql
DROP TABLE customers CASCADE;
```

`CASCADE` drops dependent objects or drops foreign keys that point to the table. The exact list differs by vendor. `CASCADE` can remove more than you expect. Views and child constraints can disappear.

**Cascading deletes** are a different idea. They are a foreign-key action on `DELETE` of a row, not a `DROP TABLE`:

```sql
CONSTRAINT fk_orders_customers
    FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
    ON DELETE CASCADE
```

When you delete a customer, the DBMS deletes the orders that point to that customer. This rule is powerful and dangerous. A delete of one parent can remove a large child tree.

Prefer `ON DELETE RESTRICT` or no cascade (the usual default: reject the delete) while you learn. Add `ON DELETE CASCADE` only when the child cannot exist without the parent and you intend that wipe.

Do not confuse:

- `DROP TABLE ... CASCADE` — remove objects
- `ON DELETE CASCADE` — remove child rows when a parent row is deleted

Do not run `DROP TABLE` in a shared database without a list of dependents. Query metadata or use the GUI to list foreign keys and views first.

Some products have `DROP TABLE IF EXISTS`. That form avoids an error when the table is already gone. Use it in scripts that you rerun.

### Questions

#### Theoretical questions

1. What does `DROP TABLE` remove?
2. Why does a foreign key block a drop by default?
3. What is the risk of `DROP TABLE ... CASCADE`?
4. What does `ON DELETE CASCADE` do?
5. How do you empty a table without dropping it?

#### Easy practical tasks

1. Create a scratch table. Drop it. Confirm that `SELECT` fails.
2. Write the difference between `DROP TABLE` and `TRUNCATE` in four sentences.
3. Draw a parent and a child. Mark what `ON DELETE CASCADE` removes.
4. List objects that you must check before you drop `customers`.

#### Medium practical tasks

1. Try to drop a parent that a child references. Record the error.
2. Add `ON DELETE CASCADE` on a practice child. Delete the parent. Show that the child rows are gone.
3. Try `DROP TABLE ... CASCADE` on your DBMS (practice schema only). Write what the product removed.

#### Advanced practical tasks

1. Write a one-page policy: when `ON DELETE CASCADE` is allowed (example: order lines under an order) and when it is forbidden (example: customers).
2. Document `DROP ... CASCADE` versus `ON DELETE CASCADE` for your DBMS with one official quote each, in your own words.

---

## Constraints: `NOT NULL`, `UNIQUE`, `CHECK`, `DEFAULT`

Constraints are rules on a table. Topic 3 placed them in integrity kinds. This section shows the SQL forms.

**`NOT NULL`.** The column must have a value on every row.

**`UNIQUE`.** No two rows share the same value (for `NULL`, products differ: many allow several `NULL` in a unique column). A unique constraint can cover more than one column.

**`CHECK`.** An expression must be true or unknown for each row. `CHECK (price >= 0)` rejects negative prices. `NULL` price makes the check unknown, so a nullable price can still be `NULL` unless you also add `NOT NULL`.

**`DEFAULT`.** When an insert omits the column, the DBMS stores this value. A default is not a constraint in the same sense, but it belongs in the column definition.

```sql
CREATE TABLE products (
    product_id INTEGER PRIMARY KEY,
    sku VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT chk_products_price CHECK (price >= 0)
);
```

`PRIMARY KEY` implies `UNIQUE` and `NOT NULL`. You do not need to write those two again on the same columns.

Name your constraints (`chk_products_price`). A name appears in error messages. Unnamed constraints get a generated name that is hard to read.

You can attach `UNIQUE`, `CHECK`, and foreign keys as table constraints. You can attach `NOT NULL` and `DEFAULT` as column attributes.

Do not use `CHECK` for rules that need another table. Use a foreign key. Some products allow subqueries in `CHECK`. That feature is not portable.

Do not rely on `DEFAULT` to hide a missing business value. If the application must send `active`, require it in the application even if the column has a default.

### Questions

#### Theoretical questions

1. What does `NOT NULL` forbid?
2. What does `UNIQUE` forbid?
3. Why can `CHECK (price >= 0)` still allow `NULL`?
4. When does `DEFAULT` apply?
5. Why do you name a `CHECK` constraint?

#### Easy practical tasks

1. Create a table with `NOT NULL`, `UNIQUE`, `CHECK`, and `DEFAULT`.
2. Insert a row that omits the defaulted column. Select the default.
3. Insert a row that breaks the `CHECK`. Record the error name.
4. Insert a duplicate `sku`. Record the error.

#### Medium practical tasks

1. Add a composite `UNIQUE (order_id, product_id)` on lines. Show a duplicate pair fail.
2. Show two `NULL` values in a unique column if the product allows it. Write the rule you observed.
3. Alter a table to add a named `CHECK`. Violate it. Drop the constraint if the syntax exists.

#### Advanced practical tasks

1. Write a constraint map for `orders`: which rules are `NOT NULL`, `UNIQUE`, `CHECK`, foreign key, or application-only.
2. Compare unique-null behavior on two products or on your product plus the SQL standard note. Write one page.

---

## Indexes as a schema object (details in section 11)

An index is a schema object that speeds up lookup. It is not a second table of business data. The DBMS maintains the index when you write rows.

You create an index with `CREATE INDEX`. A unique constraint or a primary key also creates an index in most products.

```sql
CREATE INDEX idx_orders_customer_id
    ON orders (customer_id);
```

```sql
CREATE UNIQUE INDEX idx_customers_email
    ON customers (email);
```

The unique index enforces uniqueness and supports fast lookup by email. A `UNIQUE` constraint is the logical rule. The index is the physical helper. Topic 3 made that split. Topic 11 explains B-trees, composite indexes, and when an index hurts writes.

Name indexes with a clear prefix (`idx_`). Include the table and the columns when the name length allows.

Do not add an index for every column. Each index costs space and write time. Add an index when a query filter or a join key needs it, or when a unique rule needs it.

`DROP INDEX` removes the index. It does not remove table rows. The name and syntax differ (`DROP INDEX` versus an `ALTER TABLE` drop constraint) when the index backs a constraint.

List indexes in the client or in metadata views. Confirm that you did not create a duplicate index on the same columns.

Treat indexes as part of the schema. Check them into the same change process as tables. Do not create random indexes on production without a query reason.

### Questions

#### Theoretical questions

1. What does an index speed up?
2. Who maintains the index when you insert a row?
3. How does a unique index relate to a unique constraint?
4. Why must you not index every column?
5. What does `DROP INDEX` not delete?

#### Easy practical tasks

1. Create an index on `orders(customer_id)`.
2. Create a unique index on a text column that must be unique.
3. List indexes on that table in your client.
4. Write four sentences on the cost of an extra index.

#### Medium practical tasks

1. Drop a non-constraint index. Confirm that the table and rows remain.
2. Try to create a unique index on a column that has duplicates. Record the error.
3. Name three columns that you would index in a shop and two that you would not. Give reasons.

#### Advanced practical tasks

1. Write a one-page note that points to topic 11: what you created now, and what you will measure later.
2. Find the metadata view for indexes. Write the query and the meaning of two columns in the result.

---

## Naming conventions

A naming convention is a team rule for object names. The DBMS accepts many names. Readers need one style.

Recommended habits for this path:

- Use lowercase letters, digits, and `_`.
- Use a singular or a plural table style. Pick one. This path often uses plural table names (`customers`) and singular entity talk (`customer`).
- Name primary keys `table_id` or `id`. Pick one. `customer_id` is clear in joins.
- Name foreign keys after the referenced key (`customer_id` in `orders`).
- Name indexes `idx_table_columns`.
- Name unique constraints `uq_table_columns`.
- Name checks `chk_table_rule`.
- Name foreign keys `fk_table_ref`.

Examples:

```text
customers.customer_id
orders.customer_id
idx_orders_customer_id
fk_orders_customers
chk_products_price
```

Do not use spaces in names. Do not use reserved words (`order`, `user`, `table`) as raw names. If you must, quote the name. Quoted names are case-sensitive on many products. Avoid the need to quote.

Do not mix `camelCase` and `snake_case` in one schema. SQL readers expect `snake_case` in many ecosystems.

Keep names in English if the team language for code is English. Keep names stable. A rename breaks programs.

Do not encode the type in the name (`name_varchar`). The type lives in the catalog. Do not prefix every table with the schema name (`sales_sales_orders`). The schema already namespaces the table.

Document the convention in a short page. Follow it in every `CREATE` statement.

### Questions

#### Theoretical questions

1. Why does a team need a naming convention?
2. Why is `customer_id` a useful foreign-key name?
3. Why do you avoid reserved words as table names?
4. Why must you not put spaces in names?
5. Why is a type suffix on a column a weak habit?

#### Easy practical tasks

1. Rename a messy list (`CustomerTbl`, `custName`) to a conventional list.
2. Write names for indexes and foreign keys on `orders.customer_id`.
3. List three reserved words that you must not use as raw table names.
4. Write four sentences that choose singular versus plural tables for your notes.

#### Medium practical tasks

1. Apply one convention to a three-table library schema. Write every object name.
2. Create a table with a quoted reserved name if you can. Query it. Write why you will not do this again.
3. Review an existing schema (sample or project). List five names that break the convention.

#### Advanced practical tasks

1. Write a one-page naming standard for a team: tables, columns, keys, indexes, views.
2. Map names from a vendor sample database to your convention. Do not rename a vendor sample in place; write a translation table.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from `CREATE TABLE` to a later `ALTER` and a safe `DROP`.
2. How do `NOT NULL`, `UNIQUE`, `CHECK`, and `DEFAULT` work together on one column?
3. When is `ON DELETE CASCADE` acceptable, and when is `DROP TABLE CASCADE` too wide?
4. Why is an index a schema object if it does not change the meaning of a row?
5. A teammate names tables `tblCustomers` and columns `Name`. Which convention facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: `CREATE`, types, `ALTER`, `DROP`, constraints, `CREATE INDEX`, naming.
2. Create two tables with types, keys, one check, and one index. Insert one valid row each.
3. Add a column with a default. Select it.
4. Drop the index, then drop the child table, then drop the parent.

#### Medium practical tasks

1. Build a small shop schema from scratch with named constraints. Violate each constraint once.
2. Alter a table to add a unique column. Backfill values. Then set `NOT NULL`.
3. Document every object name in that schema against your naming list.

#### Advanced practical tasks

1. Write a rerunnable DDL script: `DROP TABLE IF EXISTS` in a safe order, then `CREATE TABLE`. Run it twice.
2. Compare two DBMS type names for integer, decimal, boolean, and timestamp. Write a portable DDL and a product-adjusted DDL.
