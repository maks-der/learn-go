# 2. Data Modeling Basics

## Description

Data modeling is the work of naming the things you store and the links between them. This topic shows entities, attributes, keys, and relationship types. You also learn `NULL`.

Use one term for each concept. An entity is a thing you store. An attribute is a fact about that thing. A relationship is a link between entities. A primary key identifies one row. A foreign key points to another row. Complete this topic before you study the relational model in formal words.

Draw the model first. Write SQL after the model is clear.

---

## Entities, attributes, and relationships

An entity is a type of thing that the database stores. Examples: customer, product, order, classroom. In a relational design, one entity usually becomes one table.

An attribute is a fact about one entity. Examples: customer email, product price, order date. In a table, one attribute usually becomes one column.

A relationship is a link between two entities. Example: a customer places an order. The relationship has a name and a direction that you can read as a sentence.

Name entities with nouns. Name relationships with verbs. Keep names singular in the model (`customer`, not `customers`) if your team uses that rule. Then choose table names with a team convention (see topic 8).

Example. A library:

- entities: `member`, `book`, `loan`
- attributes of `book`: `isbn`, `title`, `year_published`
- relationships: a member borrows a book through a loan

Do not store two entities in one bag of columns without a reason. A table named `stuff` with mixed people and products is not a model. Split the entities. Add relationships.

A relationship can have attributes. A loan has `loaned_on` and `due_on`. Those facts belong to the loan, not only to the member or the book.

### Questions

#### Theoretical questions

1. What is an entity?
2. What is an attribute?
3. What is a relationship?
4. Why does a relationship sometimes need its own attributes?
5. Why is a table named `stuff` a weak model?

#### Easy practical tasks

1. List five entities for a school. Write two attributes for each entity.
2. Write three relationship sentences for that school (subject, verb, object).
3. Draw boxes for `customer`, `order`, and `product`. Draw lines for the relationships.
4. For a blog, name entities for author, post, and comment. List attributes.

#### Medium practical tasks

1. Model a sports club: member, team, match. Write attributes and relationships. Mark which relationship needs extra attributes.
2. Take a paper form (registration or receipt). Circle entities. Underline attributes. Draw the relationships.
3. Rewrite a messy list of columns (`person_or_company`, `item_or_fee`) into entities. Show the new list.

#### Advanced practical tasks

1. Model a hospital visit: patient, clinician, appointment, diagnosis. Write a one-page model with entities, attributes, and relationships.
2. Find an existing database diagram. List three entities. State one relationship that the diagram misses or that you would change.

---

## Rows/records and columns/fields

A table has a grid shape. A column is a named attribute. All values in one column have the same meaning. A row is one instance of the entity. Example: one row in `customers` is one customer.

The words *record* and *field* mean the same ideas in many tools. A record is a row. A field is a column. Use *row* and *column* in SQL work. Use *record* and *field* when a document or a form uses those words. Do not mix the pairs in one sentence without a map.

```text
customers
+----+--------+------------------+
| id | name   | email            |
+----+--------+------------------+
|  1 | Ada    | ada@example.com  |
|  2 | Alan   | alan@example.com |
+----+--------+------------------+
```

In this grid, `name` is a column. The pair `Ada` plus `ada@example.com` sits in one row. That row is one customer.

Each column has a data type. The type limits the values. Example: `email` is text. `id` is an integer. Topic 8 covers types in more detail.

Do not put two facts in one column when you must query them apart. Example: do not store `"Ada <ada@example.com>"` in one column if you must search by email. Use two columns.

A result set from a query also has rows and columns. That grid is not always a stored table. It is the output of a statement.

### Questions

#### Theoretical questions

1. What is a row in a table?
2. What is a column in a table?
3. How do the words record and field map to row and column?
4. Why must values in one column share a meaning?
5. Why is a query result not always a stored table?

#### Easy practical tasks

1. Draw a table `books` with columns `id`, `title`, `year_published`. Add two rows.
2. Label each part of your drawing as column or row.
3. Split `"Ada <ada@example.com>"` into two columns. Write the column names.
4. Write four sentences that describe the `customers` grid in this section.

#### Medium practical tasks

1. Design `rooms` with four columns. Write one invalid value for each column (wrong meaning, not only wrong type).
2. Convert a spreadsheet sheet with mixed header names into a table sketch with one meaning per column.
3. Show a query result that joins two tables. Mark which columns came from which stored table.

#### Advanced practical tasks

1. Take a CSV with ten columns. Write a target table with better column names and types. Explain three splits or merges.
2. Explain in one page how a document with nested fields maps to rows and columns. Use a single example.

---

## Primary keys and why they must be unique

A primary key is a column, or a set of columns, that identifies one row. No two rows share the same primary key value. A primary key must not be `NULL`.

The DBMS enforces uniqueness. An insert with a duplicate primary key fails. That rule keeps each entity instance distinct.

Choose a primary key that is stable. Do not use a value that users change often. A person can change an email. An email is a weak primary key if that change is common.

Example:

```sql
CREATE TABLE customers (
    customer_id INTEGER PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL
);
```

`customer_id` identifies the customer. Two customers can share a name. They must not share `customer_id`.

A primary key can have more than one column. That key is a composite primary key. Example: `order_id` plus `line_no` in `order_lines`. Together they identify one line.

Every table that you treat as an entity must have a primary key. A table without a primary key allows duplicate rows. You cannot point to one row with a foreign key in a safe way.

Do not reuse a primary key value for a new entity after you delete the old row if other systems still hold the old id. Prefer new values.

### Questions

#### Theoretical questions

1. What does a primary key identify?
2. Why must a primary key be unique?
3. Why must a primary key not be `NULL`?
4. Why is email often a weak primary key?
5. What is a composite primary key?

#### Easy practical tasks

1. Write a `CREATE TABLE` for `products` with a primary key `product_id`.
2. Explain in three sentences why two products can share a name but not a primary key.
3. Mark the primary key in a three-column sketch of `loans`.
4. Write one insert that must fail because the primary key is a duplicate. Use words or SQL.

#### Medium practical tasks

1. Design `order_lines` with a composite primary key. Write the columns and the key.
2. List three candidate attributes for `employees`. Choose a primary key. Reject the others with one reason each.
3. Try to create a table without a primary key (if the DBMS allows it). Insert two identical rows. Write the problem this causes.

#### Advanced practical tasks

1. Write a one-page policy: when the team uses a numeric `id` and when it uses a natural business key. Give two examples each.
2. Model a history table where the same business key appears many times. Choose a primary key that still identifies one row.

---

## Foreign keys and referential integrity

A foreign key is a column, or a set of columns, that points to a primary key (or a unique key) in another table. The pointed row must exist. That rule is referential integrity.

Example. `orders.customer_id` points to `customers.customer_id`. You must not store an order for a customer that does not exist.

```sql
CREATE TABLE orders (
    order_id INTEGER PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    ordered_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_orders_customers
        FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
);
```

The DBMS rejects an insert into `orders` when `customer_id` is not in `customers`. The DBMS also rejects a delete of a customer when orders still point to that customer, unless you define a cascade rule. Topic 8 and a later integrity topic cover cascade options. The default safe action is to reject the delete.

Referential integrity keeps links valid. Without it, reports show orders with missing customers. Those rows are orphans.

A foreign key can be `NULL` when the relationship is optional. Example: `orders.coupon_id` can be `NULL` when the order has no coupon. Do not use `NULL` to mean "unknown customer" if every order must have a customer. Use `NOT NULL` on that foreign key.

The foreign key columns must match the referenced key in number and in type. A single integer key needs a single integer foreign key.

### Questions

#### Theoretical questions

1. What does a foreign key point to?
2. What is referential integrity?
3. What happens when you insert an order with a missing customer id?
4. When may a foreign key be `NULL`?
5. What is an orphan row?

#### Easy practical tasks

1. Write a `CREATE TABLE` for `loans` with foreign keys to `members` and `books`.
2. Draw `customers` and `orders`. Mark the primary key and the foreign key.
3. Write one insert into `orders` that must fail because the customer does not exist.
4. Explain in four sentences why the DBMS must reject that insert.

#### Medium practical tasks

1. Create two tables with a foreign key. Insert a parent row. Insert a child row. Try a child row with a bad id. Record the error.
2. Try to delete the parent while a child exists. Record the result. Write the safe meaning of that result.
3. Design `posts` and `comments` with a required author on posts and an optional editor. Mark which foreign keys are `NOT NULL`.

#### Advanced practical tasks

1. Model three tables in a chain (`countries`, `cities`, `addresses`). Write the foreign keys. State the delete problem if you remove a country.
2. Write a one-page note: enforce foreign keys in the DBMS, not only in application code. Give two failure cases for app-only checks.

---

## One-to-one, one-to-many, many-to-many

A relationship has a cardinality. Cardinality says how many rows on each side can link.

**One-to-one.** One row in A links to at most one row in B. Example: one employee has one locker. Put a unique foreign key in one table. Example: `lockers.employee_id` is unique. Each locker points to one employee. Each employee appears at most once.

**One-to-many.** One row in A links to many rows in B. One row in B links to one row in A. Example: one customer has many orders. Put the foreign key on the "many" side. `orders.customer_id` points to `customers`.

**Many-to-many.** One row in A links to many rows in B. One row in B links to many rows in A. Example: students and courses. A student takes many courses. A course has many students. You cannot put one foreign key in `students` or in `courses`. You add a join table. The next section covers that table.

Read each relationship in both directions. "One customer has many orders" is the same design as "each order has one customer."

Do not put a comma-separated list of ids in a column to fake many-to-many. That design breaks queries and foreign keys. Use a join table.

```text
One-to-many:   customers (1) ----< orders (many)
One-to-one:    employees (1) ----- lockers (1)
Many-to-many:  students (many) >----< courses (many)
```

### Questions

#### Theoretical questions

1. What does cardinality say?
2. Where do you put the foreign key in a one-to-many relationship?
3. How do you enforce one-to-one with a foreign key?
4. Why does many-to-many need more than one foreign key column in a parent table?
5. Why is a comma-separated id list a poor design?

#### Easy practical tasks

1. Label each pair: author and bio page; customer and orders; students and courses.
2. Draw the three cardinality pictures with table names from a library.
3. Write the foreign key column for `orders` in a one-to-many design.
4. Write four sentences that explain one-to-many for invoices and invoice lines.

#### Medium practical tasks

1. Design one-to-one for `users` and `user_settings`. Write the unique foreign key.
2. Change a wrong design that stores `course_ids` in `students` into a many-to-many sketch.
3. Find two one-to-many relationships in a shop. Write both reading directions for each.

#### Advanced practical tasks

1. Model employees, departments, and department managers. Include one-to-many and one-to-one. Write keys.
2. Write a one-page review of a schema that used a text list for tags. Show the many-to-many replacement.

---

## Lookup tables and join tables

A lookup table is a small table that holds a controlled list. Examples: `countries`, `order_statuses`, `currencies`. Other tables point to it with a foreign key. The lookup table keeps the list in one place. You do not repeat the word `shipped` with different spellings in many rows.

```sql
CREATE TABLE order_statuses (
    status_code VARCHAR(20) PRIMARY KEY,
    label VARCHAR(50) NOT NULL
);

CREATE TABLE orders (
    order_id INTEGER PRIMARY KEY,
    status_code VARCHAR(20) NOT NULL,
    CONSTRAINT fk_orders_status
        FOREIGN KEY (status_code) REFERENCES order_statuses (status_code)
);
```

A join table implements many-to-many. It holds foreign keys to both sides. It often has its own attributes, such as `enrolled_on`.

```sql
CREATE TABLE student_courses (
    student_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    enrolled_on DATE NOT NULL,
    PRIMARY KEY (student_id, course_id),
    CONSTRAINT fk_sc_students FOREIGN KEY (student_id) REFERENCES students (student_id),
    CONSTRAINT fk_sc_courses FOREIGN KEY (course_id) REFERENCES courses (course_id)
);
```

Each pair `(student_id, course_id)` appears once. That primary key prevents a duplicate enrollment.

Names for a join table often combine the two entity names: `student_courses`, `order_items`. Some teams use a role name: `enrollments`. Pick one style and keep it.

Do not store the same lookup text in every parent row when the list is shared and must stay consistent. Use a lookup table. Do not add a join table for a true one-to-many. A single foreign key is enough.

### Questions

#### Theoretical questions

1. What does a lookup table store?
2. What does a join table implement?
3. Why does `student_courses` use a composite primary key?
4. When do you not need a join table?
5. Why does a lookup table reduce spelling drift?

#### Easy practical tasks

1. Write a lookup table `countries` with a code and a name. Write a foreign key from `addresses`.
2. Write a join table for `books` and `authors`. Include a `role` attribute.
3. Name three lookup tables for a shop.
4. Explain in four sentences the difference between a lookup table and a join table.

#### Medium practical tasks

1. Create `order_statuses` and `orders`. Insert two statuses. Insert one order. Try an unknown status. Record the error.
2. Create `student_courses`. Insert one enrollment. Try a duplicate pair. Record the error.
3. Redesign a table that has columns `status_text` with mixed values (`Ship`, `shipped`, `SHIPPED`). Move values to a lookup table.

#### Advanced practical tasks

1. Model products, tags, and a join table. Add a unique tag name. Write the keys and one sample query in words.
2. Write a one-page rule set: when a status is a lookup table, when it is a boolean, and when it is a history table.

---

## Surrogate keys vs natural keys

A natural key is a unique key that exists in the business. Examples: ISBN for a book, ISO country code, invoice number that the finance team already uses.

A surrogate key is a unique value that the database generates for identity. Examples: an integer `id`, a UUID. The value has no business meaning. Users do not type it as a business fact.

Both kinds can be a primary key. A natural key is useful when it is unique, stable, and always known at insert time. A surrogate key is useful when the natural key is long, composite, or can change.

Example. `countries.country_code` (`SE`, `UA`) is a stable natural key. `orders.order_id` as an integer is a common surrogate. `customers.email` is unique in many systems but can change, so many teams add `customer_id` and keep email unique as a separate constraint.

Do not use a surrogate key as an excuse to skip uniqueness of the business key. If ISBN must be unique, add a `UNIQUE` constraint on `isbn` even when `book_id` is the primary key.

```sql
CREATE TABLE books (
    book_id INTEGER PRIMARY KEY,
    isbn VARCHAR(13) NOT NULL UNIQUE,
    title VARCHAR(200) NOT NULL
);
```

Do not change a primary key as a routine business edit. If the business value can change, prefer a surrogate primary key and a unique natural key.

### Questions

#### Theoretical questions

1. What is a natural key?
2. What is a surrogate key?
3. Why can email be a poor primary key even when it is unique today?
4. Why do you still put `UNIQUE` on ISBN when `book_id` is the primary key?
5. When is a natural key a good primary key?

#### Easy practical tasks

1. Label each as natural or surrogate: ISBN, `customer_id` integer, ISO currency code, UUID.
2. Write `CREATE TABLE` for `currencies` with a natural primary key.
3. Write `CREATE TABLE` for `customers` with a surrogate primary key and unique email.
4. Write four sentences that compare the two `CREATE TABLE` examples.

#### Medium practical tasks

1. Design `invoices` with a surrogate `invoice_id` and a unique `invoice_number`. Explain who owns each value.
2. List three natural keys in a school. Mark which ones can change. Choose primary keys.
3. Debate `UUID` versus integer surrogate in six short sentences. Cover length, generation, and readability.

#### Advanced practical tasks

1. Write a one-page key policy for a team that integrates two systems. Cover natural keys that cross systems.
2. Model a person who changes national id. Show how a surrogate primary key plus history keeps old references valid.

---

## NULL: missing value, not zero and not empty string

`NULL` means that the value is missing. `NULL` is not the number `0`. `NULL` is not an empty string. `NULL` is not the text `N/A`. Those values are present values.

Use `NULL` when the attribute does not apply or is not known. Example: `employees.middle_name` can be `NULL` when the person has no middle name. Example: `orders.shipped_at` can be `NULL` when the order is not shipped.

Do not store `0` to mean "no price yet" if `0` is also a valid price. Do not store `''` to mean "unknown email" if you also query email. Use `NULL` for missing. Use a real value only when that value is true.

SQL treats `NULL` in a special way. A comparison with `NULL` does not yield true. `price = NULL` is unknown. You test missing values with `IS NULL` and `IS NOT NULL`. Topic 4 covers those predicates.

```sql
SELECT order_id
FROM orders
WHERE shipped_at IS NULL;
```

Aggregate functions ignore `NULL` in general. `SUM` of `{10, NULL, 5}` is `15`. `COUNT(column)` does not count `NULL` in that column. Topic 7 covers this behavior.

A primary key must not be `NULL`. A `NOT NULL` column must have a value on every insert.

Be consistent. If `shipped_at` is `NULL` until ship time, do not also store a status text `not shipped` unless you have a reason. Two sources of the same fact can drift.

### Questions

#### Theoretical questions

1. What does `NULL` mean?
2. Why is `0` not the same as `NULL`?
3. Why is an empty string not the same as `NULL`?
4. How do you test that a column is missing in SQL?
5. Why must a primary key refuse `NULL`?

#### Easy practical tasks

1. Write three columns that should allow `NULL` in `employees`. Write three that must not.
2. Write a `SELECT` that finds orders with `shipped_at IS NULL`.
3. Explain in four sentences why `price = 0` and `price IS NULL` differ for a shop.
4. Mark each stored value as present or missing: `0`, `''`, `NULL`, `N/A`.

#### Medium practical tasks

1. Create a table with a nullable `middle_name`. Insert one row with `NULL` and one with `''`. Select both. Write the difference.
2. Write two queries: unpaid invoices (`paid_at IS NULL`) and zero-amount invoices (`amount = 0`).
3. Find a spreadsheet that uses blank, `0`, and `N/A`. Map each to `NULL` or a real value. Justify three mappings.

#### Advanced practical tasks

1. Write a one-page rule: when your team uses `NULL`, when it uses a default, and when it uses a separate status table.
2. Design `shipments` so that you never need `NULL` for address fields on a shipped row. Show optional draft state versus shipped state.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a shop with entities, keys, and one many-to-many. Name every table and key type.
2. How do you choose between a surrogate primary key and a natural primary key for `products`?
3. Why does referential integrity need a foreign key, not only a matching column type?
4. When is `NULL` the correct value for a date column, and when is a real date required?
5. A teammate stores many tags in one text column. Which modeling facts do you use in the reply?

#### Easy practical tasks

1. Draw a model for a library: members, books, loans. Mark primary keys and foreign keys.
2. Write `CREATE TABLE` statements for two tables in that model with keys. You may omit extra constraints.
3. Fill a cheat sheet: entity, attribute, row, column, primary key, foreign key, `NULL`.
4. Add a lookup table for loan status to the library model.

#### Medium practical tasks

1. Implement customers, orders, and order lines in a local DBMS. Insert sample rows that obey the keys.
2. Try one insert that breaks a primary key and one that breaks a foreign key. Record both errors.
3. Convert a one-page spreadsheet of club members and teams into tables, including a join table if needed.

#### Advanced practical tasks

1. Model a ticketing domain (event, venue, seat, ticket, customer). Write cardinalities and keys. State where `NULL` is allowed.
2. Write a review of an existing schema (class project or public sample). List three modeling defects and a fix for each.
