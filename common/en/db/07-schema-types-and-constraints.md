# 7. Schema, Types, and Constraints

## Description

This topic shows how you create tables, choose types, and enforce rules. You learn `CREATE TABLE`, `ALTER TABLE`, `DROP TABLE`, common types, money and time-zone pitfalls, constraints, foreign-key actions, and naming.

Use one term for each concept. A schema definition statement changes objects, not row values. `CREATE`, `ALTER`, and `DROP` are data definition language (DDL). A type limits the values in a column. A constraint is a rule that the DBMS checks on write.

Practice on a throwaway database. Do not drop a table that other people use.

---

## `CREATE TABLE`, `ALTER TABLE`, `DROP TABLE`

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

`CREATE TABLE` fails when the name already exists. Use a new name, or drop the old table only if that is safe.

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

Add a constraint after the table exists:

```sql
ALTER TABLE products
    ADD CONSTRAINT chk_products_price CHECK (price >= 0);
```

The add fails if a current row breaks the check. Clean the data first, or fix the rule.

Change a type when the existing values convert cleanly. A change from `VARCHAR(200)` to `VARCHAR(50)` fails if a value is longer than 50.

Do not use `ALTER TABLE` as a daily data-fix tool. Use `UPDATE` for values. Use `ALTER` for structure.

DDL can lock a table. A long alter can block writers. Run large alters in a maintenance window in production. This path only needs small practice tables.

`DROP TABLE` removes the table object and its data. The name is gone. Indexes on that table are gone. The action is not an empty table. To keep the table and remove rows, use `DELETE` or `TRUNCATE`.

```sql
DROP TABLE scratch_items;
```

If other tables reference the table with a foreign key, the DBMS rejects the drop by default. That rejection is a safety rule.

Some products allow `DROP TABLE customers CASCADE`. `CASCADE` drops dependent objects or drops foreign keys that point to the table. The exact list differs by vendor. `CASCADE` can remove more than you expect.

Do not confuse `DROP TABLE ... CASCADE` (remove objects) with `ON DELETE CASCADE` (remove child rows when a parent row is deleted). The next sections cover foreign-key actions.

Some products have `DROP TABLE IF EXISTS`. That form avoids an error when the table is already gone. Use it in scripts that you rerun.

### Questions

#### Theoretical questions

1. What does `CREATE TABLE` define?
2. Why does a new `NOT NULL` column need a default when rows exist?
3. What does `DROP TABLE` remove that `TRUNCATE` does not remove?
4. Why does a foreign key block a drop by default?
5. Why is `ALTER` the wrong tool to fix one email value?

#### Easy practical tasks

1. Create `products` with an integer key, a name, and a decimal price.
2. Add a nullable `phone` column to `customers`. Then drop that column after you test it.
3. Create a scratch table. Drop it. Confirm that `SELECT` fails.
4. Write four sentences on the difference between `DROP TABLE` and `TRUNCATE`.

#### Medium practical tasks

1. Try `CREATE TABLE` twice with the same name. Record the error.
2. Add `NOT NULL` to a column that still has `NULL`. Record the error. Update the rows. Retry.
3. Try to drop a parent that a child references. Record the error.

#### Advanced practical tasks

1. Write a four-step plan to replace a column type on a large table without one dangerous step. Stay high-level.
2. Write a rerunnable DDL script: `DROP TABLE IF EXISTS` in a safe order, then `CREATE TABLE`. Run it twice.

---

## Integers, decimals, text, dates, booleans

A domain type limits what a column can store. The DBMS rejects a value that does not match the type. Types are the first integrity rule.

**Integer.** An integer is a whole number. Use an integer for counts, identifiers that are numbers, and quantities that cannot have a fraction. Example: `quantity` on an order line. Do not use an integer for money unless you store cents.

**Decimal.** A decimal stores a number with a fixed scale. Scale is the number of digits after the decimal point. Precision is the total number of digits. Example: `DECIMAL(12, 2)` can store values such as `1999.50` with two fraction digits. Use a decimal when the fraction must be exact.

**Text.** Text stores characters. Use a bounded type when the product gives one (`VARCHAR(n)`). Use an unbounded type only when the length has no useful limit. Store codes in a stable alphabet. Do not store a date or a number as text if you must sort or calculate on that value.

**Date and time.** A date type stores a calendar day. A timestamp type stores a day and a clock time. A time type stores a clock time without a day. Use these types for time values. Do not store a date as text if you must compare dates.

**Boolean.** A boolean stores true or false. Some products also allow unknown (`NULL`). Use a boolean for a yes/no fact. Example: `is_active`. Do not use the integers `0` and `1` unless the product has no boolean type.

Common type families (names differ by product):

| Family | Typical standard-like name | Use |
| --- | --- | --- |
| Integer | `INTEGER` | Counts, ids, quantities |
| Decimal | `DECIMAL(p, s)` | Money and exact fractions |
| Floating | `DOUBLE PRECISION` | Scientific measures, not money |
| Text | `VARCHAR(n)` | Strings; prefer `VARCHAR` |
| Date | `DATE` | Calendar day |
| Time | `TIMESTAMP` | Day plus time |
| Boolean | `BOOLEAN` | True or false |

Choose the narrowest type that fits the domain. A narrow type uses less space. A narrow type also rejects more bad writes.

`NULL` is not a type. `NULL` means the value is missing. A type plus `NOT NULL` is stronger than a type alone.

Do not mix units in one column. Example: do not store kilograms and grams in the same numeric column without a unit column. The type cannot see the unit.

Do not invent a type per product in this path when a standard-like type exists. Write `INTEGER` and `DECIMAL`. Adjust to the product when you implement.

### Questions

#### Theoretical questions

1. What does a domain type limit?
2. When do you use an integer instead of a decimal?
3. What is the difference between precision and scale on a decimal type?
4. Why must you not store a calendar date as text if you compare dates?
5. Is `NULL` a domain type? Explain.

#### Easy practical tasks

1. For each column, choose integer, decimal, text, date, or boolean: age in years, email, unit price, birth date, newsletter flag.
2. Write a `CREATE TABLE` for `customers` with five columns and a correct type for each column.
3. Make a two-column table: "Bad type" and "Better type". Add four rows from this section.
4. Write two `INSERT` statements: one that the type must accept, one that the type must reject.

#### Medium practical tasks

1. Insert a string that is too long for `VARCHAR(5)`. Record the error.
2. Compare `VARCHAR(32)` and an unbounded text type for `sku`. Write three reasons to bound the length.
3. Find the official integer and decimal type names in your DBMS. Write one sentence for each name.

#### Advanced practical tasks

1. Write a type-choice sheet for a shop: money, qty, flags, dates, codes. Map each to a type on your DBMS.
2. Design types for a library loan: copy, member, due date, returned flag, fine amount. Justify each type in one sentence.

---

## Money and floating-point pitfalls

Money is a count of a currency unit. The count must be exact. A floating-point type (`REAL`, `FLOAT`, `DOUBLE PRECISION`) is a binary approximation. Many decimal fractions have no exact binary form. Example: `0.1` plus `0.2` can fail an equality test in floating point.

Do not store money in a floating-point column. Use a decimal type with a fixed scale, or store an integer count of the smallest unit (cents). Both methods keep exact sums.

```sql
-- Prefer a decimal with a known scale
unit_price DECIMAL(12, 2) NOT NULL

-- Or an integer count of cents
unit_price_cents INTEGER NOT NULL
```

Pick one method for the whole schema. Do not mix dollars as float and cents as integer in the same database without a written rule.

Rounding is a business rule. Two-decimal money often uses banker's rounding or round-half-up. The DBMS decimal type stores the digits. Your application or a documented SQL function must apply the rounding rule at the same step every time.

Currency is not only a number. A price of `10.00` is incomplete if you do not know the currency. Use a currency code column (`USD`, `EUR`) when more than one currency exists. Do not add amounts in different currencies.

Taxes and discounts create more rounding points. Compute from a single source of truth (unit price and quantity). Store the line total that you charged. Do not recompute a historical invoice from live tax tables unless the product requires that.

Floating-point types remain useful for measurements and scientific values. They are the wrong type for ledgers.

### Questions

#### Theoretical questions

1. Why is a floating-point type a poor store for money?
2. What two exact stores for money does this section name?
3. Why is a number without a currency code incomplete when more than one currency exists?
4. Who applies the rounding rule if the type only stores digits?
5. When is a floating-point type still a correct choice?

#### Easy practical tasks

1. Write one sentence that explains why `0.1 + 0.2` is a risk for money in floating point.
2. Design `order_lines` with either `DECIMAL(12, 2)` or cents. Do not mix the two.
3. Add a `currency_code CHAR(3)` column to a price table. Write two valid example rows.
4. Label three columns in a shop as money or not money: quantity, tax amount, product name.

#### Medium practical tasks

1. In your DBMS, add `0.1` and `0.2` in a float type and in a decimal type. Record the two results.
2. Write a line-total formula from unit price and quantity. State where you round: per line, per order, or both.
3. Convert three sample prices between a decimal column and a cents column on paper. Show that the sum stays exact.

#### Advanced practical tasks

1. Write a one-page rule for a team: money type, scale, rounding, currency column, and how you store historical invoices.
2. Find one public incident or article about money in floating point. Write five sentences in your own words. Do not copy the source.

---

## Time zones and UTC

A timestamp without a time zone is a local clock reading. The DBMS does not know the zone. A timestamp with a time zone (or an instant stored in UTC) names one point on the timeline.

UTC is a standard time scale. Many systems store instants in UTC. The application converts to a local zone only when it shows the value to a person.

Use one rule for all instants:

1. Receive a time from a client with an explicit zone or as UTC.
2. Convert the value to UTC.
3. Store UTC.
4. Convert to the user zone only in the presentation layer.

A calendar date (`DATE`) has no clock time and no zone. A birthday is a date. A meeting start is an instant. Do not store a meeting start as a date only.

A "local civil time" is a different need. Example: a shop that opens at 09:00 in each city. That value is a time plus a zone name, not a single UTC instant, until you pick a day.

```sql
-- Instant: prefer UTC in the column or a timestamptz type
created_at TIMESTAMP NOT NULL  -- store UTC; document the rule

-- Civil date: no zone
birth_date DATE NOT NULL
```

Do not append a zone abbreviation such as `EST` as free text. Abbreviations collide. Use an IANA zone name (`America/New_York`) when you must store a zone.

Daylight saving changes the offset. A stored UTC instant does not move. A stored local time without a zone can become ambiguous in the overlap hour.

Do not compare timestamps that mix UTC and an unknown local zone. Convert first.

### Questions

#### Theoretical questions

1. What is the difference between a timestamp without a zone and an instant in UTC?
2. When do you convert UTC to a local zone?
3. Why is a birthday a date and not a timestamp?
4. Why is a zone abbreviation a weak store?
5. What problem does daylight saving create for a local time without a zone?

#### Easy practical tasks

1. Label each value as date, local time, or instant: birthday, "store opens at 09:00", "order placed".
2. Write the four-step UTC rule from this section in your notes.
3. Add `created_at` to a practice table. Write in a comment that the column stores UTC.
4. Name one IANA zone for a city that you know.

#### Medium practical tasks

1. Insert one row with the current UTC time. Read it. Write how your client shows the zone.
2. Find the timestamp-with-time-zone type name in your DBMS. Write how it differs from timestamp without a zone in that product.
3. Describe a daylight-saving overlap in four sentences. State what you store for an order in that hour.

#### Advanced practical tasks

1. Design columns for an event product: event date in a city, door time as civil time, and sale instant in UTC. Justify each column.
2. Write a one-page team rule: storage zone, API format, and display zone. Include one bad example that you reject.

---

## `NOT NULL`, `UNIQUE`, `CHECK`, `DEFAULT`

Constraints are rules on a table. The DBMS rejects a change that breaks a constraint.

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

## Foreign-key actions, soft delete, and audit columns

A foreign key says that a child value must match a parent key. The DBMS can also define what happens when the parent row changes.

`ON DELETE` sets the action when you delete the parent row.

| Action | Effect on child rows |
| --- | --- |
| `NO ACTION` / `RESTRICT` | Delete of the parent fails while children exist |
| `CASCADE` | The DBMS deletes the matching child rows |
| `SET NULL` | The DBMS sets the foreign key to `NULL` (column must allow `NULL`) |
| `SET DEFAULT` | The DBMS sets the foreign key to the default value |

`ON UPDATE` sets the action when the parent key value changes. Surrogate keys rarely change. Prefer stable keys so that `ON UPDATE` stays rare.

```sql
CREATE TABLE orders (
    order_id INTEGER PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    CONSTRAINT fk_orders_customer
        FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);
```

`CASCADE` on delete is convenient and dangerous. A delete of one customer can remove all orders and, if those tables also cascade, all order lines. Use `RESTRICT` (or `NO ACTION`) when history must stay. Use `CASCADE` for owned parts that have no meaning without the parent. Example: order lines belong to an order.

A hard delete removes the row. After `COMMIT`, the row is gone from the table. Related rows follow the foreign-key actions.

A soft delete keeps the row and marks it as deleted. Typical marks: a boolean `is_deleted`, or a timestamp `deleted_at` that is `NULL` for live rows.

```sql
UPDATE customers
SET deleted_at = CURRENT_TIMESTAMP
WHERE customer_id = 42;

SELECT *
FROM customers
WHERE deleted_at IS NULL;
```

Use a hard delete when the row has no legal or business value after removal. Use a soft delete when you must hide the row from normal queries but keep it for audit, undo, or unique-key history.

Soft delete has costs. Every query that means "live rows" must filter the mark. Unique keys still see the deleted row unless you design a partial unique index. A "deleted" customer can still block a new customer with the same email.

Do not use soft delete as a backup. A backup is a copy of the database. A soft-deleted row is still in the live database.

Audit columns record when a row was created and when a row last changed. They are facts about the row, not business event times.

Typical columns: `created_at` (instant of first insert), `updated_at` (instant of last update), optional `created_by` and `updated_by`.

Store these instants in UTC. Set `created_at` on insert. Do not change `created_at` on update. Set `updated_at` on insert and on every later update. A column default of "now" helps `created_at`. A trigger or application code must maintain `updated_at`. A default alone does not change on `UPDATE`.

Audit columns are not a full audit log. `updated_at` only stores the last write time. An order has `placed_at`, `shipped_at`, and `cancelled_at` as business times. Those columns have meaning. `updated_at` only says that some column changed.

Do not trust the client clock as the only source if users can set their clock. Prefer `CURRENT_TIMESTAMP` (or the vendor equivalent) in the DBMS.

### Questions

#### Theoretical questions

1. What does `ON DELETE RESTRICT` do when children exist?
2. What does `ON DELETE CASCADE` do to child rows, and why is it dangerous on a customer delete?
3. What two marks does this section name for a soft delete?
4. Why is a soft delete not a backup?
5. Why must you not change `created_at` on update, and why is a default not enough for `updated_at`?

#### Easy practical tasks

1. Write `ON DELETE` choices for: order lines under an order; orders under a customer; a nullable tag on a product.
2. Add `deleted_at` to a practice table. Soft-delete one row. Select live rows only.
3. Add `created_at` and `updated_at` to a practice table. Insert one row. Write the two values.
4. Write one sentence that distinguishes `CASCADE` from a `DELETE` that you write in SQL by hand.

#### Medium practical tasks

1. Build `customers`, `orders`, and `order_lines`. Put `CASCADE` on lines, `RESTRICT` on orders. Delete an order. Delete a customer. Record both results.
2. Put a `UNIQUE` constraint on email. Soft-delete a customer. Insert the same email again. Record what happens. Propose a fix in three sentences.
3. Update a row. Show that `created_at` stayed and `updated_at` changed (after you add a trigger or you set the column in `UPDATE`).

#### Advanced practical tasks

1. Write a delete procedure for a customer that must keep invoices. Use `RESTRICT` plus a status column. Document the steps.
2. Design audit columns plus one business timestamp table for orders. Draw which clock each report must use.

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
2. How do a domain type, a money rule, and a time-zone rule work together on one `orders` row?
3. How do `NOT NULL`, `UNIQUE`, `CHECK`, and `DEFAULT` work together on one column?
4. When do you combine `ON DELETE RESTRICT`, a soft-delete mark, and audit columns on the same parent table?
5. A teammate stores prices as `FLOAT`, names tables `tblCustomers`, and deletes customers with `ON DELETE CASCADE`. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: `CREATE`/`ALTER`/`DROP`, types, money, UTC, constraints, cascade, soft delete, audit, naming.
2. Create `products` with a decimal price, a boolean flag, a date, and audit columns. Insert two rows.
3. Add a child table with `ON DELETE RESTRICT`. Soft-delete a parent. Show that the child still references the parent.
4. Drop the child table, then drop the parent.

#### Medium practical tasks

1. Build a small shop schema from scratch with named constraints, exact money, and UTC instants. Violate each constraint once.
2. Alter a table to add a unique column. Backfill values. Then set `NOT NULL`.
3. Document every object name in that schema against your naming list.

#### Advanced practical tasks

1. Implement a customer archive: soft delete, unique email after delete, restrict on invoices, cascade on a preferences child. Record each statement and result.
2. Compare two DBMS type names for integer, decimal, boolean, and timestamp. Write a portable DDL and a product-adjusted DDL.
