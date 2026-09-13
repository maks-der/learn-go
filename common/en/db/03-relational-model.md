# 3. Relational Model

## Description

The relational model is the theory behind tables, keys, and constraints. This topic maps the formal words to the tables that you already know. You also learn catalogs, views, and the difference between the logical schema and physical storage.

Use one term for each concept. A relation is a table with unique rows. A tuple is a row. An attribute is a column. A constraint is a rule that the DBMS enforces. Complete this topic before you write SQL queries.

Keep the model vendor-neutral. Storage details differ by product. The logical ideas stay the same.

---

## Relations, tuples, and attributes (in plain words)

A relation is a set of tuples that share the same attributes. In practice you call the relation a table. The set idea matters: a true relation does not contain two identical tuples. A SQL table can allow duplicate rows if you do not define a key. Always define a key so that the table behaves like a relation.

A tuple is one row. It holds one value for each attribute. The order of tuples has no meaning in the model. You add order only when you query with `ORDER BY`.

An attribute is a named column with a domain. The domain is the set of allowed values. Example: an attribute `qty` has domain "positive integer". The name `qty` is unique inside the relation.

Example. Relation `products`:

- attributes: `product_id`, `name`, `price`
- one tuple: `(10, 'Notebook', 4.50)`

The relational model does not use a nested table inside a cell as a first-class idea. Each attribute value is atomic in the basic model. Topic 10 covers first normal form.

SQL names differ from theory names. Learn both:

| Theory | SQL practice |
| --- | --- |
| Relation | Table |
| Tuple | Row |
| Attribute | Column |
| Domain | Data type plus extra rules |
| Degree | Number of columns |
| Cardinality | Number of rows |

Do not depend on row order on disk. A query without `ORDER BY` can return any order. The model treats the relation as a set.

### Questions

#### Theoretical questions

1. What is a relation in plain words?
2. What is a tuple?
3. Why does a relation need unique tuples?
4. Why does tuple order have no meaning in the model?
5. How does a domain limit an attribute?

#### Easy practical tasks

1. Map the words relation, tuple, and attribute onto a sketch of `customers`.
2. Write one tuple for `products` as a list of named values.
3. Fill the theory-versus-SQL table with one extra example row in your notes.
4. Explain in four sentences why `ORDER BY` is not part of the stored relation.

#### Medium practical tasks

1. Create a table without a key. Insert two identical rows. Write how this table fails the set idea.
2. Count the degree and the cardinality of a table that you create with three columns and five rows.
3. Rewrite a nested contact list (name plus many phones in one cell) as relations with atomic attributes.

#### Advanced practical tasks

1. Read a short public summary of Codd's relational model. Write one page that maps five terms to SQL.
2. Explain why a result of `UNION` is closer to a relation than a result of a query that allows duplicates. Use an example.

---

## Schemas and catalogs

A schema is a named group of objects. Objects include tables, views, and constraints. A schema keeps names organized. Two schemas can each have a table named `users`. The full name is `schema_name.users`.

A catalog is a database or a higher container, depending on the product. In the SQL standard, a catalog contains schemas. In many products you connect to a database, and that database contains schemas. Learn the three-level name when your product uses it: `catalog.schema.table`.

Example names:

```text
learn.public.customers
shop.sales.orders
```

`learn` is the database. `public` or `sales` is the schema. `customers` and `orders` are tables.

The word *schema* also means the design: the list of tables and columns. That meaning is the logical schema. The named schema object is a namespace. Use "schema object" when you mean the namespace. Use "logical schema" when you mean the design.

Create a schema for a subject area when the database grows. Example: `sales`, `hr`, `inventory`. Do not put every experiment in the system schema if the product gives you a default schema for users.

The catalog also stores metadata: table names, column types, keys. You query metadata through information views. The view names differ by vendor. The idea is the same. The DBMS describes the database in tables.

Do not guess object names. List schemas and tables from the client or from the metadata views.

### Questions

#### Theoretical questions

1. What does a schema object group?
2. What is a catalog in the SQL standard idea?
3. Why can two tables share the name `users`?
4. What are the two meanings of the word schema?
5. What does metadata describe?

#### Easy practical tasks

1. Write the three-part name for table `orders` in database `shop` and schema `sales`.
2. List the schemas that your local DBMS shows in the practice database.
3. Create a schema named `practice` if the product allows it. Create one table in it.
4. Explain in four sentences why a default schema exists.

#### Medium practical tasks

1. Create the same table name in two schemas. Insert different rows. Select from each full name.
2. Find the metadata view or command that lists columns for one table. Write the command and two columns that it shows.
3. Draw a catalog with two schemas and three tables. Label each level.

#### Advanced practical tasks

1. Compare naming in two products (for example PostgreSQL versus SQLite). Write how each product uses database, schema, and table.
2. Write a one-page naming plan for a company database with `sales`, `hr`, and `audit` schemas.

---

## Integrity constraints (domain, entity, referential)

An integrity constraint is a rule that every stored change must satisfy. The DBMS rejects a change that breaks a constraint. Constraints keep bad data out.

**Domain constraint.** A value must belong to the attribute domain. Types enforce part of this rule: an integer column rejects `'abc'`. Extra rules tighten the domain: `price >= 0`, `email` matches a pattern, `status` is one of a list. `NOT NULL` is a domain rule: the value must be present.

**Entity constraint.** Each tuple must be identifiable. A primary key implements this rule. No two rows share the key. The key is not `NULL`. Entity integrity is the formal name for this rule.

**Referential constraint.** A foreign key value must match an existing key or must be `NULL` when the design allows `NULL`. Referential integrity is the formal name. Topic 2 showed the same idea with tables.

Example:

```sql
CREATE TABLE products (
    product_id INTEGER PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    CONSTRAINT chk_products_price CHECK (price >= 0)
);

CREATE TABLE order_lines (
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    qty INTEGER NOT NULL,
    PRIMARY KEY (order_id, product_id),
    CONSTRAINT chk_ol_qty CHECK (qty > 0),
    CONSTRAINT fk_ol_products FOREIGN KEY (product_id) REFERENCES products (product_id)
);
```

`CHECK (price >= 0)` is a domain constraint. `PRIMARY KEY` is an entity constraint. `FOREIGN KEY` is a referential constraint.

Do not rely only on the application for these rules. A second client can skip application checks. Put the rules in the database.

### Questions

#### Theoretical questions

1. What does an integrity constraint do on a write?
2. What is a domain constraint?
3. What is entity integrity?
4. What is a referential constraint?
5. Why must the DBMS enforce constraints, not only the application?

#### Easy practical tasks

1. Label each rule: `price >= 0`, primary key on `product_id`, foreign key to `products`.
2. Write a `CHECK` that keeps `qty` greater than `0`.
3. Write three domain values that a `DECIMAL` price must reject.
4. Explain entity integrity in four sentences with `customers` as the example.

#### Medium practical tasks

1. Create `products` with a `CHECK` on `price`. Try a negative price. Record the error.
2. Create `order_lines` with a foreign key. Try a missing product. Record the error.
3. List every constraint on a three-table shop sketch. Mark each as domain, entity, or referential.

#### Advanced practical tasks

1. Write a one-page constraint map for enrollments: types, keys, and checks. State one rule that you cannot express with a simple `CHECK`.
2. Compare application validation and DBMS constraints. Give two attacks or bugs that only the DBMS stops.

---

## Candidate keys, superkeys, composite keys

A superkey is a set of attributes that uniquely identifies a tuple. The set can contain extra attributes that you do not need for uniqueness.

A candidate key is a minimal superkey. You cannot remove an attribute and still keep uniqueness. A table can have more than one candidate key.

The primary key is the candidate key that you choose as the main identifier. Other candidate keys become unique constraints.

A composite key is a key with more than one attribute. A composite key can be a candidate key or a primary key.

Example. `employees`:

- `employee_id` is unique (surrogate)
- `national_id` is unique in that country
- `email` is unique in the company

Each of these can be a candidate key. `employee_id` plus `email` is a superkey, but it is not minimal. The pair is not a candidate key.

Example. `order_lines`:

- `(order_id, line_no)` is a composite candidate key
- `(order_id, product_id)` can also be a candidate key if one product appears once per order

```sql
CREATE TABLE employees (
    employee_id INTEGER PRIMARY KEY,
    national_id VARCHAR(32) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL
);
```

Do not add columns to a key without a reason. A wider key is harder to use as a foreign key. Prefer a minimal candidate key as the primary key.

### Questions

#### Theoretical questions

1. What is a superkey?
2. What is a candidate key?
3. How does a primary key relate to candidate keys?
4. What is a composite key?
5. Why is `employee_id` plus `email` not a candidate key when each column is already unique?

#### Easy practical tasks

1. List two candidate keys for `books` if `book_id` and `isbn` are unique.
2. Mark one superkey that is not a candidate key in `employees` from this section.
3. Write a composite primary key for `student_courses`.
4. Explain in four sentences why a candidate key is minimal.

#### Medium practical tasks

1. Design `flights` with `flight_id` and a unique triple (`carrier`, `flight_no`, `depart_date`). Name all candidate keys.
2. Create a table with two unique columns. Insert a clash on the second unique column. Record the error.
3. Draw a table of keys: superkey, candidate, primary, composite. Add one example each from a school.

#### Advanced practical tasks

1. Write a one-page key analysis for a medical record: patient id, national id, and visit number. Choose the primary key.
2. Show how a composite primary key becomes a composite foreign key in a child table. Write both `CREATE TABLE` statements.

---

## Views vs base tables

A base table is a stored table. Rows live in the database files. `CREATE TABLE` defines a base table. Inserts write to that table.

A view is a named query. `CREATE VIEW` stores the query text, not a second copy of the rows (in the basic case). When you select from the view, the DBMS runs the query.

```sql
CREATE VIEW open_orders AS
SELECT order_id, customer_id, ordered_at
FROM orders
WHERE shipped_at IS NULL;
```

```sql
SELECT order_id
FROM open_orders
WHERE customer_id = 15;
```

Views give a stable name to a frequent query. Views hide columns that a role must not see. Views can join tables so that a report user reads one object.

A view is not a cache by default. Each select can read the base tables again. Some products offer materialized views that store a snapshot. That feature is extra. Treat a normal view as a saved `SELECT`.

You can update some simple views. Many views are read-only because the query is a join or an aggregate. Do not assume that `INSERT` into a view works. Test the product. Prefer writes to base tables.

Do not use a view as a substitute for a missing table in the model. If an entity needs storage, create a base table. Use a view to show or restrict data.

### Questions

#### Theoretical questions

1. What does a base table store?
2. What does `CREATE VIEW` store in the basic case?
3. When does the DBMS run the view query?
4. Why is a normal view not a cache?
5. Why do many views reject `INSERT`?

#### Easy practical tasks

1. Write a view that lists customers with only `customer_id` and `name`.
2. Write a `SELECT` from that view.
3. Explain in four sentences the difference between `open_orders` and table `orders`.
4. Name two reasons to give a report user a view instead of the base table.

#### Medium practical tasks

1. Create `orders` and view `open_orders`. Insert a shipped row and an open row. Select from the view. Confirm the shipped row is absent.
2. Change a base table column that the view uses. Select from the view. Write what happened.
3. Try `INSERT` into a view that is a simple filter. Record whether the product allows it.

#### Advanced practical tasks

1. Write a view that joins `customers` and `orders` and returns one row per order with customer name. State why writes to this view are a poor idea.
2. Read the product docs on materialized views. Write five sentences: what is stored, when it refreshes, and when you use it.

---

## Logical schema vs physical storage (high-level)

The logical schema is the design that users and programs see: tables, columns, keys, views, and constraints. You change the logical schema with `CREATE`, `ALTER`, and `DROP`.

Physical storage is how the DBMS places bytes on disk: files, pages, indexes, and heap or clustered layout. The product chooses many of these details. You tune some of them later (topic 11 and later operations topics).

The relational model lets you work at the logical level. You write `SELECT` against tables. You do not read pages by hand. The DBMS maps the logical row to physical pages.

Example. Two products store the same `customers` table in different file layouts. Your SQL stays the same. The query plan can differ. The logical schema is portable. The physical plan is not fully portable.

Indexes are physical helpers that you declare. They do not change the meaning of the table. A unique index can support a logical unique constraint. The constraint is logical. The index is physical.

Do not design tables around a rumor about disk order. Design keys and constraints for meaning. Add indexes for access paths after you know the queries.

Do not confuse a tablespace or filegroup with a schema namespace. A schema is a name group. A tablespace is a storage location. Products use different words.

### Questions

#### Theoretical questions

1. What objects belong to the logical schema?
2. What is physical storage in this section?
3. Why can the same `SELECT` have different plans on two products?
4. How does an index relate to a unique constraint?
5. Why must you not design tables from a rumor about disk order?

#### Easy practical tasks

1. Make a two-column table: "Logical" and "Physical". Add five rows from this section.
2. Write three logical objects in your practice database.
3. Explain in four sentences why SQL hides pages from the beginner.
4. Label each item: `CREATE TABLE`, index on `email`, data file, `PRIMARY KEY` meaning.

#### Medium practical tasks

1. Create a table and an index. Use the client to show that the table meaning did not change. Write the objects you see.
2. Find the product term for the data file or tablespace. Write one sentence that separates it from schema.
3. Draw a stack: SQL, logical schema, indexes, pages, files. Write one sentence per layer.

#### Advanced practical tasks

1. Read a high-level storage page for your DBMS. Write one page: heap versus clustered table if the product has both. Stay high-level.
2. Write a change plan: add a column (logical) versus move a tablespace (physical). State who must approve each change.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Map a full shop (`customers`, `orders`, `order_lines`) to relation, tuple, attribute, and the three integrity kinds.
2. How do candidate keys and a primary key appear in `CREATE TABLE` for `employees`?
3. When do you add a view, and when do you add a base table?
4. What stays the same across vendors at the logical level, and what can change in storage?
5. A teammate says "the schema is the folder on disk." Which two meanings do you separate in the reply?

#### Easy practical tasks

1. Write a cheat sheet: relation, tuple, schema object, catalog, view, constraint kinds, superkey, candidate key.
2. Create two tables with primary keys, one foreign key, and one `CHECK`. Insert one valid set of rows.
3. Create a view that hides one column of a table. Select from the view.
4. Write the full object name of one table in your DBMS (database and schema if they exist).

#### Medium practical tasks

1. Implement `products` and `order_lines` from this topic. Force one domain failure, one entity failure, and one referential failure.
2. Add a second unique column as a candidate key. Show that a duplicate on that column fails.
3. Document the metadata command that lists constraints on one table. Paste the output into your notes.

#### Advanced practical tasks

1. Design a logical schema for a library in one page: schemas, tables, keys, two views. Omit physical file names.
2. Compare views and base tables for a security role that may read orders but must not read payment columns. Write the view and the grant idea in words.
