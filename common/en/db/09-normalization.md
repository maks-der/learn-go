# 9. Normalization

## Description

Normalization is the work of splitting tables so that each fact lives in one place. This topic shows update anomalies, functional dependency, normal forms through BCNF, and when you denormalize on purpose. An orders example ties the ideas together.

Use one term for each concept. A functional dependency means one set of attributes determines another. A normal form is a rule that a table can satisfy. Complete this topic before you add indexes. A messy table wastes indexes.

Normalize to remove anomalies. Do not split tables only to make more objects.

---

## Insert, update, and delete anomalies

An anomaly is a bad side effect when you write a poorly designed table. The same fact sits in many rows. A write can leave the facts in conflict, or you cannot store a fact at all.

**Update anomaly.** You change a fact in one row and you forget the copies. Example: a table `order_rows` stores `customer_name` on every line. Ada changes her name. You update one line. Other lines still show the old name.

**Insert anomaly.** You cannot store a fact until an unrelated fact exists. Example: you cannot store a new customer until that customer has an order, because customer columns live only in the order table.

**Delete anomaly.** You delete a row and you lose a fact that you still need. Example: you delete the last order of a customer and you lose the customer address, because the address lived only on the order row.

```text
Bad table order_rows
order_id  product  customer_name  customer_city
1         pen      Ada            London
1         ink      Ada            London
2         pen      Alan           Paris
```

Update Ada's city in one row only. The table now has two cities for Ada. That split is an update anomaly.

Delete order 2. If Alan has no other row, Alan and Paris disappear. That loss is a delete anomaly.

You cannot insert a customer Grace who has no order. That gap is an insert anomaly.

Normalization moves `customer_name` and `customer_city` to a `customers` table. Each customer fact exists once. Orders point to `customer_id`.

Do not store derived copies of parent attributes on child rows unless you denormalize on purpose (later section) and you control the update path.

### Questions

#### Theoretical questions

1. What is an update anomaly?
2. What is an insert anomaly?
3. What is a delete anomaly?
4. Why does a repeated customer city cause an update anomaly?
5. How does a `customers` table remove those anomalies?

#### Easy practical tasks

1. Copy the `order_rows` grid. Show an update anomaly with two cities for Ada.
2. Show a delete anomaly when you remove Alan's only order.
3. Show an insert anomaly for a customer with no order.
4. Draw the split into `customers` and `order_lines`.

#### Medium practical tasks

1. Build the bad table in SQL. Run the three writes. Record the bad results.
2. Build the split tables. Repeat the three writes. Show that the anomalies are gone.
3. Find a spreadsheet with repeated customer columns. Mark one example of each anomaly.

#### Advanced practical tasks

1. Write a one-page review of a real messy export. Name anomalies and the target tables.
2. Explain in one page why application code that "updates all copies" is a weaker fix than a key to a parent table.

---

## Functional dependency

A functional dependency (FD) says: if two rows have the same values for attributes X, they must have the same values for attributes Y. X determines Y. We write `X → Y`.

Example. In `customers`, `customer_id → name`. One id has one name. Two rows with the same id cannot have two names. The primary key determines every column in the table.

Example. `isbn → title` in a book table if ISBN identifies a book. `title → isbn` is false. Two books can share a title.

Example. In a bad `order_rows` table, `order_id → customer_id` can hold, and `customer_id → customer_city` can hold. Then `order_id → customer_city` also holds by transitivity. The city still belongs with the customer, not with the order line.

A determinant is the left side (X). If X is a candidate key, the dependency is expected. If X is not a key, and Y is not part of X, the table stores a fact about something other than the key. That extra fact is the usual reason to split the table.

```text
order_id, product_id → qty          (line fact)
customer_id → customer_city         (customer fact)
order_id → customer_id              (order fact)
```

Do not store `customer_city` in a table whose key is `(order_id, product_id)`. The city does not depend on the product.

FDs are about meaning, not about the current sample. A sample can look like `city → country` if you have one city name. The rule is still false in the world (Paris exists in more than one country). Use the real rule.

### Questions

#### Theoretical questions

1. What does `X → Y` mean?
2. Why does a primary key determine every column in a well-designed table?
3. Why is `title → isbn` usually false?
4. What is a determinant?
5. Why must you not infer an FD from a small sample alone?

#### Easy practical tasks

1. Write three FDs for a library (`book_id`, `isbn`, `member_id`).
2. Mark a false FD in `order_rows` that you must not keep (`product → customer_city`).
3. Explain `customer_id → customer_city` in four sentences.
4. List columns that `order_id` determines and columns that it does not determine on an order line.

#### Medium practical tasks

1. For a school, write FDs: student, course, enrollment date, instructor office. Group FDs by entity.
2. Show a sample where `city → country` looks true and a second row that breaks it.
3. Take the bad `order_rows` table. List FDs. Mark which determinants are not keys of that table.

#### Advanced practical tasks

1. Write a one-page FD analysis for invoices: invoice, customer, line, product, tax rate.
2. Explain transitivity with `order_id → customer_id` and `customer_id → city`. State where city must live.

---

## 1NF, 2NF, 3NF, and BCNF (high-level)

Normal forms are tests. Each test is stronger than the one before it for the usual path.

**First normal form (1NF).** Each attribute holds a single value, not a list or a nested table. Each row is unique (a key exists). Rows and columns have no hidden repeating groups.

Bad 1NF: `phones = '111,222'` or two columns `phone1`, `phone2` as a repeating group. Good 1NF: a `customer_phones` table with one phone per row.

**Second normal form (2NF).** The table is in 1NF. Every non-key attribute depends on the whole primary key, not on a part of a composite key.

Bad 2NF: `order_lines` with key `(order_id, product_id)` and column `customer_id` that depends only on `order_id`. Move `customer_id` to `orders`.

**Third normal form (3NF).** The table is in 2NF. No non-key attribute depends on another non-key attribute. Non-key attributes depend only on a key.

Bad 3NF: `orders` with `customer_id` and also `customer_city`. City depends on `customer_id`, not on `order_id`. Move city to `customers`.

```text
1NF: atomic values, a key
2NF: no partial dependency on a composite key
3NF: no transitive dependency through a non-key
```

A table with a single-column primary key is automatically in 2NF if it is in 1NF. Partial dependencies need a composite key.

3NF is the usual target for operational systems. It removes the common anomalies.

**Boyce-Codd normal form (BCNF)** is a stricter form of 3NF. The rule in plain words: every determinant is a candidate key.

3NF still allows a special case: a non-key determinant that points into a candidate key. That case is rare. BCNF forbids it.

Example (classic shape). A table `offerings` of a course:

- A course has many teachers.
- A teacher teaches only one course.
- The key is `(course, classroom)` or `(teacher, classroom)` depending on the rules.

If `teacher → course` and `teacher` is not a candidate key of the table, BCNF fails. You split into `teacher_course(teacher, course)` and a table that assigns rooms.

For most beginner schemas, 3NF and BCNF agree. If every determinant is a key, you are in BCNF.

```text
3NF:   non-key attributes depend only on keys
BCNF:  every determinant is a candidate key
```

Do not chase normal forms as a score. Use them to find misplaced facts. Do not spend hours on BCNF puzzles before you can model orders and customers. Learn the slogan: determinants must be keys. If you find a determinant that is not a key, split the table.

BCNF can make some constraints harder to keep in one table. You may need two tables plus an application or a DBMS constraint that spans tables. That trade-off is why some texts stay at 3NF for teaching.

If your table has one candidate key and no extra FDs, BCNF holds.

### Questions

#### Theoretical questions

1. What does 1NF forbid in a cell?
2. What extra rule does 2NF add, and what extra rule does 3NF add?
3. Why is a table with a single-column key already free of partial dependencies?
4. What is the BCNF rule in one sentence, and how does it differ from 3NF at a high level?
5. Why do many beginner tables that are in 3NF also satisfy BCNF?

#### Easy practical tasks

1. Split a `phones` list column into a child table. Write the keys.
2. Move `customer_id` off a line table whose key is `(order_id, product_id)`. Move `customer_city` off `orders`.
3. Write the BCNF slogan on your cheat sheet. Take `customers(customer_id, email, name)` with unique email. List determinants. Is the table in BCNF?
4. Label three defects as 1NF, 2NF, or 3NF.

#### Medium practical tasks

1. Normalize a table `enrollments(student_id, course_id, student_name, course_title, instructor, instructor_office)` step by step. Write 1NF/2NF/3NF tables.
2. Show why `phone1` and `phone2` fail 1NF as a repeating group. Create the bad and good schemas in SQL. Insert the same facts.
3. Read a short public BCNF example (course/teacher). Draw the split tables. Write why unique `email` as a candidate key keeps `email → name` compatible with BCNF.

#### Advanced practical tasks

1. Write a one-page walk-through of 1NF, 2NF, and 3NF on a hospital visit sheet. Then add a one-page comparison of 3NF and BCNF with one example that differs and one that does not.
2. Design `teacher`, `course`, and `assignment` so that `teacher → course` lives in a two-column table whose key is `teacher`.

---

## When to denormalize on purpose

Denormalization is a planned copy of a fact into another table. You do it after you have a normalized model. You accept extra storage and extra update work for a faster or simpler read.

Valid reasons:

- a report that must not join ten tables on every page load
- a cache of a display name on an event row that never changes
- a data warehouse fact table that stores a date and a label for read speed
- a search index that copies fields from several tables

Invalid reasons:

- you did not want to write a join
- you copied columns "just in case" without an owner of updates
- you stored a list in a text column to avoid a child table and you still need to query items

Rules when you denormalize:

- Name the source of truth. One table remains the owner.
- Plan the update: same transaction, trigger, or batch job.
- Document the delay if the copy can be stale.
- Do not denormalize before you measure a real read problem.

Example. `orders` stores `customer_id`. For a shipping label print job, you copy `ship_to_name` onto the order at checkout time. That copy is a snapshot. Later name changes do not rewrite old labels. The snapshot is a business choice, not a mistake.

Example. A `customer_name` column on every `order_lines` row that you update when the customer changes is a poor denormalization. The update cost is high and easy to miss.

Do not confuse denormalization with ignoring 1NF. A JSON document store is a different model. In a relational operational database, start from 3NF.

### Questions

#### Theoretical questions

1. What is denormalization?
2. When is a snapshot column on an order a valid copy?
3. Why is "I do not want to join" a weak reason?
4. What must you document when a copy can be stale?
5. Why do you normalize first?

#### Easy practical tasks

1. Write three valid reasons and three invalid reasons from this section.
2. Design a snapshot `ship_to_name` on `orders`. Write who fills it and when.
3. Mark a copied `customer_city` on lines as a poor copy. Write why.
4. Explain in four sentences the source of truth idea.

#### Medium practical tasks

1. Add a denormalized `item_count` on `orders`. Update it in the same transaction as line changes. Write the risk if you forget.
2. Compare a warehouse-style fact row (copied labels) with an OLTP 3NF order. Write six sentences.
3. Find a cache or search copy in a system you know. Name the owner table.

#### Advanced practical tasks

1. Write a one-page decision record: keep 3NF for checkout, add a reporting table that copies daily totals.
2. Design a stale-read budget: how late a copied name may be, and how you refresh it.

---

## Example: orders and order lines

This example is the standard split.

**Facts.**

- A customer places many orders.
- An order has many lines.
- A line is a quantity of one product.
- Product name and price list live with the product. The line may store the sold price as a snapshot.

**Normalized tables.**

```sql
CREATE TABLE customers (
    customer_id INTEGER PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL
);

CREATE TABLE products (
    product_id INTEGER PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    list_price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE orders (
    order_id INTEGER PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    ordered_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_orders_customers
        FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
);

CREATE TABLE order_lines (
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    qty INTEGER NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (order_id, product_id),
    CONSTRAINT chk_ol_qty CHECK (qty > 0),
    CONSTRAINT fk_ol_orders FOREIGN KEY (order_id) REFERENCES orders (order_id),
    CONSTRAINT fk_ol_products FOREIGN KEY (product_id) REFERENCES products (product_id)
);
```

**Dependencies.**

- `customer_id → name, city`
- `product_id → name, list_price`
- `order_id → customer_id, ordered_at`
- `(order_id, product_id) → qty, unit_price`

`unit_price` on the line is a snapshot of the sale. It can differ from `list_price` later. That copy is intentional.

**Bad combined table (for contrast).**

```text
order_id, product_id, qty, product_name, list_price, customer_name, city
```

This table mixes line, product, and customer facts. It fails 2NF and 3NF. Anomalies appear.

A write of a new order uses one transaction: insert `orders`, insert `order_lines`. The previous topic covers that transaction.

### Questions

#### Theoretical questions

1. Why is `(order_id, product_id)` the line key in this design?
2. Why does `city` not belong on `order_lines`?
3. Why may `unit_price` differ from `list_price`?
4. Which FDs justify table `orders`?
5. What anomalies does the combined table create?

#### Easy practical tasks

1. Draw the four tables and the foreign keys.
2. Insert one customer, one product, one order, two lines (on paper or in SQL).
3. Change the customer city. Show that old lines do not store a stale city.
4. Write the four FD bullets from this section in your notes.

#### Medium practical tasks

1. Implement the four tables. Insert sample rows. Join them to reprint an invoice.
2. Try to insert a line for a missing order. Record the error.
3. Add a second order for the same customer. Show that customer facts stay in one row.

#### Advanced practical tasks

1. Extend the example with `order_status` as a lookup table. Write FDs and keys.
2. Write a one-page invoice reconstruction query plan in words: which tables, which joins, which snapshot columns.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do FDs tell you that a column is in the wrong table?
2. Walk one messy table through 1NF, 2NF, and 3NF to the orders example.
3. When is BCNF extra work beyond 3NF for a beginner shop schema?
4. When is a snapshot column denormalization, and when is it an anomaly?
5. A teammate stores products, customers, and qty in one grid "like Excel." Which three facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: three anomalies, FD, 1NF, 2NF, 3NF, BCNF, denormalize rules, orders split.
2. Normalize a three-column messy list (`student`, `course`, `student_email`) into tables.
3. Draw before-and-after for the `order_rows` example.
4. Label `unit_price` on a line as snapshot or anomaly, with one reason.

#### Medium practical tasks

1. Take a paper receipt. Produce 3NF tables. Insert the receipt as rows.
2. Show an update anomaly on a denormalized copy that you did not plan. Then fix the model.
3. Write FDs for members, books, and loans. Create the tables.

#### Advanced practical tasks

1. Normalize a messy club spreadsheet that includes teams, members, dues, and coaches. Deliver tables, keys, and FDs.
2. Write a decision record for one denormalized reporting table next to the 3NF shop.
