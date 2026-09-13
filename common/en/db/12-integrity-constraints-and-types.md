# 12. Integrity, Constraints, and Types

## Description

This topic shows how you keep stored values correct. You learn domain types, money and floating-point errors, time zones, cascading foreign keys, soft delete versus hard delete, and audit columns.

Use one term for each concept. A domain is the set of allowed values for a column. A constraint is a rule that the DBMS checks on write. A type is the DBMS name for a domain. Complete this topic after you can create tables, keys, and `CHECK` rules. Types and constraints are the rules that those tables must obey.

This path is vendor-neutral. Type names differ by product. The ideas do not.

---

## Domain types: integers, decimals, text, dates, booleans

A domain type limits what a column can store. The DBMS rejects a value that does not match the type. Types are the first integrity rule.

**Integer.** An integer is a whole number. Use an integer for counts, identifiers that are numbers, and quantities that cannot have a fraction. Example: `quantity` on an order line. Do not use an integer for money.

**Decimal.** A decimal stores a number with a fixed scale. Scale is the number of digits after the decimal point. Precision is the total number of digits. Example: `NUMERIC(12, 2)` can store values such as `1999.50` with two fraction digits. Use a decimal when the fraction must be exact.

**Text.** Text stores characters. Use a bounded type when the product gives one (`VARCHAR(n)`). Use an unbounded type only when the length has no useful limit. Store codes in a stable alphabet. Do not store a date or a number as text if you must sort or calculate on that value.

**Date and time.** A date type stores a calendar day. A timestamp type stores a day and a clock time. A time type stores a clock time without a day. Use these types for time values. Do not store a date as text if you must compare dates.

**Boolean.** A boolean stores true or false. Some products also allow unknown (`NULL`). Use a boolean for a yes/no fact. Example: `is_active`. Do not use the integers `0` and `1` unless the product has no boolean type.

Choose the narrowest type that fits the domain. A narrow type uses less space. A narrow type also rejects more bad writes.

```sql
CREATE TABLE products (
    product_id   INTEGER PRIMARY KEY,
    sku          VARCHAR(32) NOT NULL,
    name         VARCHAR(200) NOT NULL,
    unit_price   NUMERIC(12, 2) NOT NULL,
    in_stock     BOOLEAN NOT NULL,
    released_on  DATE
);
```

`NULL` is not a type. `NULL` means the value is missing. A type plus `NOT NULL` is stronger than a type alone.

Do not mix units in one column. Example: do not store kilograms and grams in the same numeric column without a unit column. The type cannot see the unit.

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

1. Change `unit_price` from text to `NUMERIC(12, 2)` on a scratch table. Record the error if bad text exists. Clean the data. Change the type.
2. Compare `VARCHAR(32)` and an unbounded text type for `sku`. Write three reasons to bound the length.
3. Find the official integer and decimal type names in your DBMS. Write one sentence for each name.

#### Advanced practical tasks

1. Design types for a library loan: copy, member, due date, returned flag, fine amount. Justify each type in one sentence.
2. Write a one-page note: what the DBMS checks at the type layer versus what a `CHECK` constraint must add.

---

## Money and floating-point pitfalls

Money is a count of a currency unit. The count must be exact. A floating-point type (`REAL`, `FLOAT`, `DOUBLE PRECISION`) is a binary approximation. Many decimal fractions have no exact binary form. Example: `0.1` plus `0.2` can fail an equality test in floating point.

Do not store money in a floating-point column. Use a decimal type with a fixed scale, or store an integer count of the smallest unit (cents). Both methods keep exact sums.

```sql
-- Prefer a decimal with a known scale
unit_price NUMERIC(12, 2) NOT NULL

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
2. Design `order_lines` with either `NUMERIC(12, 2)` or cents. Do not mix the two.
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

## Cascading foreign keys: `ON DELETE` / `ON UPDATE`

A foreign key says that a child value must match a parent key. The DBMS can also define what happens when the parent row changes.

`ON DELETE` sets the action when you delete the parent row.

| Action | Effect on child rows |
| --- | --- |
| `NO ACTION` / `RESTRICT` | Delete of the parent fails while children exist |
| `CASCADE` | The DBMS deletes the matching child rows |
| `SET NULL` | The DBMS sets the foreign key to `NULL` (column must allow `NULL`) |
| `SET DEFAULT` | The DBMS sets the foreign key to the default value |

`ON UPDATE` sets the action when the parent key value changes. The same action names apply. Surrogate keys rarely change. Natural keys change more often. Prefer stable keys so that `ON UPDATE` stays rare.

```sql
CREATE TABLE orders (
    order_id     INTEGER PRIMARY KEY,
    customer_id  INTEGER NOT NULL,
    CONSTRAINT fk_orders_customer
        FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);
```

`CASCADE` on delete is convenient and dangerous. A delete of one customer can remove all orders and, if those tables also cascade, all order lines. Use `RESTRICT` (or `NO ACTION`) when history must stay. Use `CASCADE` for owned parts that have no meaning without the parent. Example: order lines belong to an order.

`SET NULL` means "this child stays, but it has no parent." That action fits optional links. It does not fit a required customer on an order.

Do not use cascade as a substitute for a planned delete procedure. Know the full child graph before you enable `CASCADE`.

Default actions differ by product. Read the product manual. Write the action in the `CREATE TABLE` statement so that the schema is explicit.

### Questions

#### Theoretical questions

1. What does `ON DELETE RESTRICT` do when children exist?
2. What does `ON DELETE CASCADE` do to child rows?
3. When is `SET NULL` a valid delete action?
4. Why do surrogate keys make `ON UPDATE` rare?
5. Why is `CASCADE` dangerous on a customer delete?

#### Easy practical tasks

1. Write `ON DELETE` choices for: order lines under an order; orders under a customer; a nullable tag on a product.
2. Add a foreign key on a practice child table with `ON DELETE RESTRICT`. Try to delete the parent. Record the error.
3. Make a three-row table: parent, child, recommended `ON DELETE` action.
4. Write one sentence that distinguishes `CASCADE` from a `DELETE` that you write in SQL by hand.

#### Medium practical tasks

1. Build `customers`, `orders`, and `order_lines`. Put `CASCADE` on lines, `RESTRICT` on orders. Delete an order. Delete a customer. Record both results.
2. Change a natural key that children reference, with `ON UPDATE CASCADE`. Confirm that child keys change.
3. Draw the delete graph for four tables in a shop. Mark each edge `CASCADE` or `RESTRICT`.

#### Advanced practical tasks

1. Write a delete procedure for a customer that must keep invoices. Use `RESTRICT` plus a status column. Document the steps.
2. Read the default `ON DELETE` action in your DBMS. Write whether an omitted clause is safe for your schema.

---

## Soft delete vs hard delete

A hard delete removes the row. After `COMMIT`, the row is gone from the table. Related rows follow the foreign-key actions.

A soft delete keeps the row and marks it as deleted. Typical marks:

- a boolean `is_deleted`
- a timestamp `deleted_at` that is `NULL` for live rows

```sql
ALTER TABLE customers
    ADD COLUMN deleted_at TIMESTAMP;

-- Soft delete
UPDATE customers
SET deleted_at = CURRENT_TIMESTAMP  -- store UTC; see the time-zone section
WHERE customer_id = 42;

-- Live rows
SELECT *
FROM customers
WHERE deleted_at IS NULL;
```

Use a hard delete when the row has no legal or business value after removal, and when no child history must stay.

Use a soft delete when you must hide the row from normal queries but keep it for audit, undo, or unique-key history.

Soft delete has costs. Every query that means "live rows" must filter the mark. Unique keys still see the deleted row unless you design a partial unique index or you change the key. Foreign keys still see the row. A "deleted" customer can still block a new customer with the same email.

Do not mix the two methods on the same table without a written rule. Do not use soft delete as a backup. A backup is a copy of the database. A soft-deleted row is still in the live database.

Restore after a hard delete needs a backup or an audit log. Restore after a soft delete is an `UPDATE` that clears the mark.

### Questions

#### Theoretical questions

1. What does a hard delete do to the row?
2. What two marks does this section name for a soft delete?
3. Why must live queries filter `deleted_at`?
4. How can a soft-deleted email block a new customer?
5. Why is a soft delete not a backup?

#### Easy practical tasks

1. Add `deleted_at` to a practice table. Soft-delete one row. Select live rows only.
2. Write one hard `DELETE` and one soft `UPDATE` for the same primary key (on a scratch row).
3. List three tables in a shop. Choose hard or soft delete for each. Give one reason.
4. Write the live-row `WHERE` clause for `is_deleted` and for `deleted_at`.

#### Medium practical tasks

1. Put a `UNIQUE` constraint on email. Soft-delete a customer. Insert the same email again. Record what happens. Propose a fix in three sentences.
2. Write a view `v_live_customers` that hides soft-deleted rows. Select from the view.
3. Compare undo: restore from backup versus clear `deleted_at`. Write four differences.

#### Advanced practical tasks

1. Design unique emails that allow reuse after soft delete (partial unique index or composite key). Implement the design if your DBMS allows it.
2. Write a one-page policy: which tables hard-delete, which soft-delete, who can purge, and how purge interacts with foreign keys.

---

## Audit columns (`created_at`, `updated_at`)

Audit columns record when a row was created and when a row last changed. They are facts about the row, not business event times.

Typical columns:

| Column | Meaning | Who writes it |
| --- | --- | --- |
| `created_at` | Instant of first insert | DBMS default or the application, once |
| `updated_at` | Instant of last update | DBMS trigger or the application, on every update |
| `created_by` | User or service that inserted the row | Application (optional) |
| `updated_by` | User or service that last updated the row | Application (optional) |

```sql
CREATE TABLE products (
    product_id  INTEGER PRIMARY KEY,
    name        VARCHAR(200) NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);
```

Store these instants in UTC. Set `created_at` on insert. Do not change `created_at` on update. Set `updated_at` on insert (same as `created_at`) and on every later update.

A default of "now" on the column helps `created_at`. A trigger or application code must maintain `updated_at`. A default alone does not change on `UPDATE`.

Audit columns are not a full audit log. A full audit log stores old values and the statement or the row image. `updated_at` only stores the last write time.

Do not use `updated_at` as a business "status changed at" unless that is the only status clock you need. An order has `placed_at`, `shipped_at`, and `cancelled_at` as business times. Those columns have meaning. `updated_at` only says that some column changed.

Do not trust the client clock as the only source if users can set their clock. Prefer `CURRENT_TIMESTAMP` (or the vendor equivalent) in the DBMS.

### Questions

#### Theoretical questions

1. What fact does `created_at` record?
2. Why must you not change `created_at` on update?
3. Why is a column default not enough for `updated_at`?
4. How does an audit column differ from a full audit log?
5. Why is `updated_at` a weak substitute for `shipped_at`?

#### Easy practical tasks

1. Add `created_at` and `updated_at` to a practice table. Insert one row. Write the two values.
2. Update the row. Show that `created_at` stayed and `updated_at` changed (after you add a trigger or you set the column in `UPDATE`).
3. Make a glossary: `created_at`, `updated_at`, `placed_at`. One sentence each.
4. Write an `INSERT` that omits audit columns and relies on defaults, if your DBMS allows that.

#### Medium practical tasks

1. Create a trigger or application helper that sets `updated_at` on `UPDATE`. Prove it with two updates.
2. Add `created_by` as text. Insert as user `learn`. Do not take the name from an untrusted client without a note in your comments.
3. Compare DBMS `CURRENT_TIMESTAMP` with a value that you pass from the client. Write which source you choose and why.

#### Advanced practical tasks

1. Design audit columns plus one business timestamp table for orders. Draw which clock each report must use.
2. Write a one-page standard for a team: UTC, defaults, triggers, and what a row-level audit log would add later.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a domain type, a money rule, and a time-zone rule work together on one `orders` row?
2. When do you combine `ON DELETE RESTRICT`, a soft-delete mark, and audit columns on the same parent table?
3. Which store do you choose for a price, a birthday, and an order instant, and why must those three choices differ?
4. What integrity failures remain after types are correct if foreign-key actions and delete policy are wrong?
5. A teammate stores prices as `FLOAT` and deletes customers with `ON DELETE CASCADE`. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: types, money, UTC, cascade actions, soft versus hard delete, audit columns.
2. Create `products` with a decimal price, a boolean flag, a date, and audit columns. Insert two rows.
3. Add a child table with `ON DELETE RESTRICT`. Soft-delete a parent. Show that the child still references the parent.
4. List every column on your practice `orders` table. Mark type, zone rule, and delete role.

#### Medium practical tasks

1. Build a three-table shop. Apply exact money, UTC instants, restrict on customer delete, cascade on line delete, and audit columns. Run one insert and one failed delete.
2. Write a live-customer view and an order report that uses `placed_at`, not `updated_at`.
3. Document a delete policy for the three tables in eight yes/no questions.

#### Advanced practical tasks

1. Implement a customer archive: soft delete, unique email after delete, restrict on invoices, cascade on a preferences child. Record each statement and result.
2. Write a schema review checklist of twelve items from this topic. Apply it to a public sample schema or to your shop. Mark pass or fail for each item.
