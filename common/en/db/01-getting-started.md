# 1. Getting Started

## Description

A database stores data so that programs can find it, change it, and keep it after the program stops. This topic shows what a database is, how it differs from a file, and how you start a local database management system (DBMS).

Use one term for each concept. Data is raw values. Information is data in a useful context. A database is an organized collection of data. A DBMS is the software that manages the database. Complete this topic before you model tables or write SQL.

This path is vendor-neutral. The examples use standard SQL where SQL applies. Product names appear only as examples of a category.

---

## What a database is and why programs use one

A database is a collection of related data that a program stores and retrieves. The data stays on disk after the program exits. Many programs can use the same data. The DBMS controls who reads and who writes.

Programs use a database when they must keep facts across sessions. A shop keeps products, customers, and orders. A ticket system keeps seats and bookings. A log service keeps events. In each case the program writes data and later reads it again.

A database gives structure. You define tables or collections. You define keys. You define rules. The DBMS checks those rules. A database also gives concurrent access. Two users can work at the same time. The DBMS isolates their changes.

A program that holds all facts only in memory loses those facts when the process stops. A program that appends lines to a text file can keep facts, but it does not give queries, transactions, or access control. A database is the usual tool when many records, many users, or strong rules exist.

Example of a small shop database:

- table `products` with product name and price
- table `customers` with name and email
- table `orders` with customer, product, quantity, and order date

The shop program inserts an order. Later the same program, or a report program, reads the orders. The data lives in the database, not only in the shop process.

### Questions

#### Theoretical questions

1. What does a database keep after a program exits?
2. Why do two programs use one database instead of two copies of the same file?
3. What does the DBMS check when you define rules?
4. Why is memory alone not enough for a shop order list?
5. Name three jobs that a DBMS does for a program.

#### Easy practical tasks

1. Write five sentences that describe a database. Use only facts from this section.
2. Make a two-column table: "Program need" and "How a database helps". Add four rows.
3. List three program types that need a database. List two that do not. Give one reason for each choice.
4. Draw three boxes: program, DBMS, stored data. Draw arrows for write and read.

#### Medium practical tasks

1. Describe a library system in six short sentences. Name the facts that must stay after the program stops.
2. Compare a shopping list in a phone note with a shop inventory in a database. Write four differences.
3. Find a public description of one DBMS. Write three services that the product claims. Mark each as structure, access, or durability.

#### Advanced practical tasks

1. Interview a teammate or read a post-mortem of a system that lost data. Write one page: what was stored, what failed, and which database job was missing.
2. Design a one-page diagram of a ticket system with two programs (booking and report) and one database. Label each data flow.

---

## Data vs information vs a database vs a DBMS

Data is a raw value. Example: `42`, `2026-09-13`, or the text `Ada`. Data has no meaning until you name the context.

Information is data plus meaning. Example: "Ada placed 42 orders on 2026-09-13." The same number `42` is data. The sentence is information.

A database is the organized collection. The collection has a defined structure. Example: tables `customers` and `orders` with named columns. The database holds many rows of data. Queries turn those rows into information.

A DBMS is the software process that manages databases. You install a DBMS. You create a database in that DBMS. You connect to the DBMS. The DBMS reads and writes the files on disk. You do not edit those files by hand.

Keep these terms separate:

| Term | Meaning | Example |
| --- | --- | --- |
| Data | A raw value | `19.99` |
| Information | Data in context | Price of SKU `A-10` is `19.99` |
| Database | Organized collection | The `shop` database |
| DBMS | Managing software | The server process that hosts `shop` |

Do not call the DBMS "the database" when you mean the product. Say "the DBMS" for the software. Say "the database" for the collection. One DBMS can host many databases.

### Questions

#### Theoretical questions

1. What is the difference between data and information?
2. What does a database contain that a single value does not?
3. What is a DBMS?
4. Can one DBMS host more than one database?
5. Why must you not edit DBMS data files by hand?

#### Easy practical tasks

1. Label each item as data, information, database, or DBMS: `PostgreSQL`, `orders` collection, `7`, "7 seats remain".
2. Write one sentence that turns the data `London` and `14:05` into information.
3. Make a four-row glossary with the four terms from this section.
4. Name two databases that one DBMS can host in a learning setup (for example `shop` and `class`).

#### Medium practical tasks

1. Take a paper receipt. List five data values. Write one information sentence for the receipt as a whole.
2. Draw a stack: applications, DBMS, database files, disk. Write one sentence for each layer.
3. Read the start page of one DBMS manual. Quote the product name and write whether the page describes the DBMS or a sample database.

#### Advanced practical tasks

1. Write a one-page note that explains the four terms to a person who only knows spreadsheets. Use a single running example.
2. Compare two DBMS products at a high level. For each, state: process model, typical database object, and one use case. Do not copy marketing text.

---

## When a file is enough and when a database is better

A file is enough when one program writes the data, one user reads it, and the structure is simple. Examples: a configuration file, a static list, a small cache, or an export that another tool will import.

A file is a poor store when many writers exist, when you must query by many fields, when you must keep rules, or when you must recover after a crash. A database is better in those cases.

Use a file when all of these are true:

- one writer, or writers that never overlap
- you read the whole file or you scan it in a simple way
- loss of the last write is acceptable, or you copy the file as a backup
- no other program must query the same data at the same time

Use a database when any of these are true:

- many users or many processes write at the same time
- you need queries such as "all orders for one customer in one month"
- you need rules such as "each order must point to a real customer"
- you need a transaction: all writes succeed, or none succeed
- you need users and permissions

Example. A personal note app can store notes in one JSON file. A classroom booking system must not give the same seat to two students. The booking system needs a database.

Do not put a production multi-user store in a shared spreadsheet or a single JSON file. Those tools do not give safe concurrent writes.

### Questions

#### Theoretical questions

1. When is a file a sufficient store?
2. Why is a file a poor store for two writers at the same time?
3. Which need forces a transaction?
4. Why is a shared spreadsheet a weak production store?
5. Name four signals that you must move from a file to a database.

#### Easy practical tasks

1. For each case, choose file or database: app settings, bank transfers, a poem draft, a hotel booking.
2. Write four short sentences that justify one of those choices.
3. List three file formats that beginners use as a store (`CSV`, `JSON`, `SQLite` file is a database). Mark which item is not "only a file".
4. Describe one failure that happens when two users edit the same CSV at the same time.

#### Medium practical tasks

1. Take a real folder of CSV exports. Write when those files are enough and when a database is better.
2. Design a checklist with eight yes/no questions that decide file versus database. Apply it to a blog and to a payment ledger.
3. Measure or estimate: how you find one row in a 100000-line CSV versus a table with an index. Write the difference in words.

#### Advanced practical tasks

1. Write a one-page decision record for a team: keep `config.yaml` as a file, move user accounts to a database. Include risks.
2. Reproduce a lost-update story with two editors and one shared file. Document the steps and the lost change. Do not use this method in production.

---

## Client, server, and embedded databases

A client is a program that sends requests to a database. The client can be your application, a command-line tool, or a graphical tool.

A server DBMS runs as a separate process. The client connects over a network or a local socket. Many clients connect at the same time. Examples of this shape: PostgreSQL, MySQL, Microsoft SQL Server. The server owns the data files. The server applies permissions and transactions.

An embedded DBMS runs inside your process. The library opens a file. There is no separate server process. Example of this shape: SQLite. One file holds the database. Embedded databases are simple to start. They fit single-user tools, tests, and small apps.

Compare the two shapes:

| Shape | Process | Typical use | Concurrent writers |
| --- | --- | --- | --- |
| Server | Separate DBMS process | Shared apps, many users | Designed for this |
| Embedded | Inside the application | Local apps, tests, tools | Limited; often one writer |

A local server on your machine is still a server. "Local" means the network address is your computer. "Embedded" means the DBMS is a library in the same process.

Do not treat SQLite and a server DBMS as the same operations problem. Backup, users, and network access differ. Learn the ideas on either shape. Then learn the operations of the product that you will run.

### Questions

#### Theoretical questions

1. What is a database client?
2. What process owns the data files in a server DBMS?
3. What is an embedded DBMS?
4. Why does a multi-user web app usually use a server DBMS?
5. Is a DBMS on `localhost` embedded? Explain.

#### Easy practical tasks

1. Label each product as server or embedded: SQLite, PostgreSQL, MySQL. Use public docs if you are not sure.
2. Draw two diagrams: app plus SQLite file; app plus network plus server DBMS.
3. Write three reasons to use an embedded DBMS in a test suite.
4. Write three reasons to use a server DBMS for a class project that two people share.

#### Medium practical tasks

1. Install or locate one embedded and one server DBMS. Write the start command or library name for each.
2. Compare connection setup: file path versus host, port, user, and password. Make a two-column table.
3. Read the concurrency limits of SQLite in its official FAQ. Write three sentences in your own words.

#### Advanced practical tasks

1. Run the same three SQL statements against an embedded DBMS and a server DBMS. Record differences in connection and in types.
2. Write a short architecture note: when a desktop app must move from embedded to server. Include backup and second-user access.

---

## Common categories: relational, document, key-value, wide-column, graph, search

A data model is the way a DBMS organizes values. Learn the common categories. Pick a category after you know the access pattern. Do not pick a product only by popularity.

**Relational.** Data lives in tables. Rows share a fixed set of columns. You relate tables with keys. SQL is the usual language. Use this model for structured records and strong rules. Example domains: orders, payroll, bookings.

**Document.** Data lives in documents, often JSON. One document can nest related fields. Documents in one collection can differ in shape. Use this model when each item is a self-contained document and the shape changes. Example domains: profiles, content pages.

**Key-value.** You store a value under a key. You read by key. You do not query inside the value as a first-class feature. Use this model for caches, sessions, and simple lookups. Example domains: session store, feature flags.

**Wide-column.** You store rows with a flexible set of columns, grouped by column family. Writes and reads often use a row key. Use this model for large write volumes and time-series style access. Example domains: metrics, large event tables.

**Graph.** Data is nodes and edges. You query paths and neighborhoods. Use this model when the relationships are the main question. Example domains: social links, access paths, recommendations.

**Search.** A search engine indexes text and structured fields. You query by relevance, tokens, and filters. Use this model when users search text. A search store is not a full replacement for a system of record.

Many systems use more than one category. That approach is polyglot persistence. Learn the relational model first in this path. The other categories appear again in a later survey topic.

```text
Relational   : tables, rows, SQL
Document     : JSON documents
Key-value    : get / set by key
Wide-column  : row key plus columns
Graph        : nodes and edges
Search       : inverted index, ranked hits
```

### Questions

#### Theoretical questions

1. What is a data model in this section?
2. Which category uses tables and SQL as the usual interface?
3. When do you choose a key-value store?
4. Why is a search engine not a full system of record?
5. What does polyglot persistence mean?

#### Easy practical tasks

1. Match each need to one category: "find friends of friends", "store a session by id", "invoice lines with tax rules", "full-text blog search".
2. Write one sentence per category in the list above.
3. Make a table with columns "Category", "Unit of data", "Typical query".
4. Name one product in each category from public knowledge. Mark them as examples only.

#### Medium practical tasks

1. Take a social app. Split its data across two categories. Justify each split in three sentences.
2. Explain why a shopping cart can live in a key-value store while orders live in a relational database.
3. Read a short overview of one non-relational product. Write three things it does not replace in a relational DBMS.

#### Advanced practical tasks

1. Write a one-page comparison of relational versus document for a product catalog. Include queries and rules.
2. Design a system that uses relational, key-value, and search. Draw the write path and the read path for one user action.

---

## Installing a local DBMS (or using Docker)

Install a local DBMS so that you can practice. A local install keeps data on your machine. You do not need a cloud account for this path.

Two common methods exist:

1. Install the official package for your operating system.
2. Run a container with Docker.

Pick one method. Use the same method for the rest of this topic.

**Official package.** Download the installer from the vendor site. Run the installer. Set a password for the admin user if the installer asks. Accept the default port if no other service uses it. Typical ports: PostgreSQL `5432`, MySQL `3306`.

**Docker.** You need Docker Engine. Pull an official image. Map the port. Set the password in an environment variable. Persist data in a named volume. Example shape (PostgreSQL as an illustration):

```text
docker run --name learn-db -e POSTGRES_PASSWORD=devpass -p 5432:5432 -v learn-db-data:/var/lib/postgresql/data -d postgres:16
```

Change the image and the variables when you use a different product. Read the official image page. Do not copy passwords from examples into production.

After the install, verify the process. On the host, confirm that the port listens. In Docker, run `docker ps` and confirm that the container is up.

Create a practice database. Use a name such as `learn`. Do not use the system database for your tables. System databases are for the DBMS itself.

If the service does not start, read the server log. Common causes: the port is in use, the data directory has the wrong owner, or the password environment variable is missing.

### Questions

#### Theoretical questions

1. Why do you install a local DBMS for this path?
2. What are the two common install methods in this section?
3. Why do you persist a Docker volume for database data?
4. Why must you not practice in a system database?
5. What is a common cause when the DBMS does not start?

#### Easy practical tasks

1. Install a local relational DBMS or start a container. Write the product name and version.
2. Confirm that the server process or container is running. Save the command and the output.
3. Create a database named `learn`. Record the statement or the GUI action.
4. Write the host, port, user, and database name that you will use. Do not publish a real production password.

#### Medium practical tasks

1. Stop the DBMS and start it again. Confirm that the `learn` database still exists.
2. Change the host port mapping in Docker or the listen port in a config file. Connect on the new port. Then restore the default.
3. Read the official Docker image page for your DBMS. List the required environment variables.

#### Advanced practical tasks

1. Run two versions of the same DBMS on different ports. Write how you keep the data directories separate.
2. Write a compose file with one DBMS service, a volume, and a health check. Start it and connect.

---

## Connecting with a GUI client and a CLI

A client connects with a connection string or with separate fields. Typical fields:

- host (example: `127.0.0.1` or `localhost`)
- port
- user
- password
- database name

A CLI is a terminal client. You type SQL and you see text results. Use the CLI when you learn. The CLI shows the exact statement. Examples of CLI tools: `psql`, `mysql`, `sqlite3`, `sqlcmd`. The name depends on the product.

A GUI client is a graphical program. You browse tables, you edit rows, and you run SQL in a window. Use a GUI to see structure. Do not use only the GUI. You must still type SQL.

Connect with the CLI first. Example shape (the program name changes by product):

```text
psql -h 127.0.0.1 -p 5432 -U learn -d learn
```

After you connect, run a probe statement:

```sql
SELECT 1;
```

The result must show `1`. This test proves that authentication works and that the session can run SQL.

In the GUI, create a connection with the same host, port, user, and database. Run `SELECT 1` again. Browse the list of databases and schemas. Confirm that `learn` appears.

Do not commit passwords into a repository. Use the client password prompt or a local config file that you do not share. Do not expose the DBMS port on a public network for a learning install.

If the client fails, read the error. Typical causes: wrong port, wrong user, wrong database name, server not running, or a firewall block.

### Questions

#### Theoretical questions

1. Which fields does a client need to connect to a server DBMS?
2. Why must you learn the CLI and not only a GUI?
3. What does `SELECT 1` prove?
4. Why must you not put a password in a shared repository?
5. Name three causes of a failed connection.

#### Easy practical tasks

1. Connect with the CLI. Run `SELECT 1`. Save the full session output.
2. Connect with a GUI client. Run `SELECT 1`. Capture a screenshot for your notes.
3. List the databases that the client shows. Mark the practice database.
4. Write your connection fields in a private note: host, port, user, database. Omit the password or use a placeholder.

#### Medium practical tasks

1. Fail a connection on purpose (wrong port). Record the error. Fix the port and connect.
2. Run `SELECT current_user` or the vendor equivalent. Write who you are in the session.
3. Compare the CLI result and the GUI result for `SELECT 1`. Write one advantage of each client.

#### Advanced practical tasks

1. Connect through a connection URI if your client supports it. Show the URI with a password placeholder. Connect successfully.
2. Create a read-only user if the DBMS allows it. Connect as that user. Try `SELECT` and `CREATE TABLE`. Record which statement fails.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from an empty machine to a successful `SELECT 1`. Name install, start, create database, and connect.
2. How do data, information, a database, and a DBMS relate in one running example?
3. When do you keep a file, and when do you install a DBMS, for a two-person class project?
4. What is the difference between an embedded DBMS and a server DBMS on `localhost`?
5. A teammate wants only a document store for invoices with tax rules. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: terms, categories, connection fields, and `SELECT 1`.
2. Create the `learn` database if it does not exist. Connect with CLI and GUI. Run `SELECT 1` in both.
3. Fill a table with six categories and one learning use case each.
4. Export `SELECT 1` results from the CLI to a text file. Keep the file in your notes.

#### Medium practical tasks

1. Write a short script (PowerShell or bash) that checks whether the DBMS port is open and then runs a CLI `SELECT 1`.
2. Document your install in ten steps so that another beginner can copy it. Include the verify step.
3. Start the DBMS, stop it, and show which client error appears when the server is down. Then start it again.

#### Advanced practical tasks

1. Run one relational DBMS and one embedded DBMS. Execute `SELECT 1` on both. Write a comparison of start-up and connection.
2. Create two practice databases on one server. Connect to each. Show that objects in one database do not appear in the other.
