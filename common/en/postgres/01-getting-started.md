# 1. Getting Started

## Description

PostgreSQL is an open-source object-relational database system. This topic shows what PostgreSQL is, why the major version matters, and how you install a server. You also learn clusters, databases, schemas, roles, and how you connect.

Learn common database ideas first (`db.topics.md`). Then complete this topic. Complete this topic before you study connections, SQL, and types.

Use one term for each concept. A cluster is one data directory with one running server. A database lives in a cluster. A schema groups objects inside a database. A role is a login name or a group name. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## What PostgreSQL is (object-relational, open source)

PostgreSQL is a relational database. You store data in tables. You query data with SQL. PostgreSQL is also object-relational. You can define types, operators, and functions. Tables can inherit columns from other tables. Most new work uses normal tables, types, and constraints. Do not treat inheritance as the default design.

The PostgreSQL Global Development Group maintains the project. The license is the PostgreSQL License. You may use, change, and distribute the software. The project is not a single vendor product.

PostgreSQL runs on Linux, macOS, and Windows. Client programs connect over TCP. The default port is `5432`. The server process is `postgres`. Older text sometimes says `postmaster` for the main process.

PostgreSQL follows SQL standards for many features. It also adds PostgreSQL-specific features. Examples: `RETURNING`, `ILIKE`, `DISTINCT ON`, and `jsonb`. Later topics in this path cover those features.

PostgreSQL is not a document-only store. You can store JSON. Tables, types, and constraints stay the primary model. PostgreSQL is not a small embedded file like SQLite. You run a server. Clients connect to that server.

### Questions

#### Theoretical questions

1. What does "object-relational" mean for PostgreSQL?
2. Who maintains PostgreSQL?
3. What license does PostgreSQL use?
4. What is the default TCP port?
5. How is PostgreSQL different from an embedded file database?

#### Easy practical tasks

1. Write five sentences that describe PostgreSQL. Use only facts from this section.
2. Make a two-column table: "PostgreSQL property" and "What it means". Add four rows.
3. List three application types that fit PostgreSQL. List two types that do not fit well. Give one reason for each choice.
4. Open [https://www.postgresql.org/](https://www.postgresql.org/). Write the current stable major versions that the site shows.

#### Medium practical tasks

1. Compare PostgreSQL with one other database that you know. Write six short sentences. Cover license, SQL, server process, and types.
2. Draw a simple diagram: client, TCP port `5432`, `postgres` process, data directory. Label each box.
3. Open the official "About" page. Write three facts that this section does not already list.

#### Advanced practical tasks

1. Read the PostgreSQL License text. Write a one-page summary with permissions and conditions. Use your own words.
2. Find the current versioning policy on postgresql.org. Write how long a major version receives fixes. Use the official table.

---

## Versions and why the major version matters

PostgreSQL uses a major.minor version number. Examples: `16.15` and `17.11`. The first number is the major version. The rest is the minor version.

A minor release fixes defects and security problems. You can apply a minor release without a dump and restore of the data directory. Stop the server. Install the new binaries. Start the server.

A major release can change catalogs, on-disk format, and SQL behavior. You cannot point a PostgreSQL 17 binary at a PostgreSQL 16 data directory. You must upgrade with `pg_upgrade`, dump and restore, or logical replication. Topic 16 covers upgrade methods.

The major version controls which features you have. Documentation URLs include the major version. Example: [https://www.postgresql.org/docs/17/](https://www.postgresql.org/docs/17/). The word `current` in the docs URL points to the newest stable major version. That major version can be newer than 17.

Use the same major version in development, test, and production. Do not develop on 17 and deploy to 16 without a feature check. Do not skip minor updates for a long time.

Show the server version after you connect:

```sql
SELECT version();
SHOW server_version;
SHOW server_version_num;
```

`server_version_num` is an integer. Example: PostgreSQL 17.11 reports `170011`. Compare versions with that number in scripts.

### Questions

#### Theoretical questions

1. What is the difference between a major release and a minor release?
2. Can a PostgreSQL 17 server open a PostgreSQL 16 data directory? Why?
3. Why must development and production use the same major version?
4. What does `SELECT version();` return?
5. What is `server_version_num`?

#### Easy practical tasks

1. Open [https://www.postgresql.org/support/versioning/](https://www.postgresql.org/support/versioning/). Write the final support date for PostgreSQL 16 and PostgreSQL 17.
2. Write one sentence that explains when you need `pg_upgrade`.
3. Make a table with columns "Change type" and "Typical upgrade method". Add rows for minor and major.
4. Find the release notes for PostgreSQL 16 and PostgreSQL 17. Write one feature from each major version.

#### Medium practical tasks

1. Compare `docs/16`, `docs/17`, and `docs/current` in the browser. Write which major version `current` uses today.
2. Plan a move from 16 to 17 in five steps. Do not run the upgrade yet. Name the official tool for each step that needs a tool.
3. Read the "Migration" section of the 17 release notes. List three compatibility items that a beginner must check.

#### Advanced practical tasks

1. Write a short policy for a team: how you pick a major version, how you apply minor updates, and how you test them.
2. Compare dump/restore and `pg_upgrade` from the official upgrade docs. Write when each method is the better choice.

---

## Installing PostgreSQL or using the official Docker image

You can install PostgreSQL with the official installer, with a package manager, or with Docker.

The official Windows and macOS installers are on [https://www.postgresql.org/download/](https://www.postgresql.org/download/). Linux distributions provide packages. The PostgreSQL Global Development Group also publishes the PGDG apt and yum repositories. Prefer a repository that names the major version. Example: `postgresql-17`.

The official Docker image is `postgres` on Docker Hub. Pin a major tag. Examples: `postgres:16` and `postgres:17`. Do not use `postgres:latest` for real work. The latest tag can change major version.

A minimal Docker run:

```text
docker run --name pg17 -e POSTGRES_PASSWORD=secret -p 5432:5432 postgres:17
```

Useful environment variables:

- `POSTGRES_PASSWORD` — password for the superuser (required)
- `POSTGRES_USER` — superuser name (default `postgres`)
- `POSTGRES_DB` — first database (default is the user name)

The image runs `initdb` on the first start when the data volume is empty. Data lives in `/var/lib/postgresql/data` inside the container. Mount a volume if you must keep the data after you remove the container.

After install, confirm that the server listens. On a host install, the service name is often `postgresql` or `postgresql-17`. In Docker, use `docker logs pg17` and look for "ready to accept connections".

Do not expose port `5432` to the public internet. Topic 21 covers production network rules.

### Questions

#### Theoretical questions

1. Why must you pin a Docker major tag such as `postgres:17`?
2. What does `POSTGRES_PASSWORD` set?
3. When does the official image run `initdb`?
4. Where does the official image store the data directory inside the container?
5. Why is a package that names the major version useful?

#### Easy practical tasks

1. Install PostgreSQL 16 or 17, or start `postgres:17` in Docker. Record the method and the major version.
2. Run `docker pull postgres:17` or open the installer download page. Write the exact image tag or installer file name.
3. List the three environment variables from this section and the default of each variable when the image sets a default.
4. Draw the path from your terminal to the server: client host, port, container or service name.

#### Medium practical tasks

1. Start PostgreSQL with a named Docker volume. Stop and remove the container. Start a new container on the same volume. Confirm that the data directory still exists.
2. Install or run both 16 and 17 on different ports (`5432` and `5433`). Write the two listen addresses.
3. Compare the official installer with the Docker image. Write a short report: data location, how you stop the server, and how you upgrade the minor version.

#### Advanced practical tasks

1. Build a Docker Compose file with `postgres:17`, a named volume, and a health check. Start it. Show that the health check passes.
2. Read the official Docker image documentation. Write how the image runs scripts from `/docker-entrypoint-initdb.d`. Create one `.sql` file and prove that it runs only on the first init.

---

## `psql` and a GUI (`pgAdmin`, DBeaver, DataGrip)

`psql` is the official terminal client. It ships with the PostgreSQL client package. Use `psql` for scripts, `\copy`, and meta-commands. Topic 2 covers meta-commands in detail.

Connect with host, port, user, and database:

```text
psql -h localhost -p 5432 -U postgres -d postgres
```

Set `PGPASSWORD` only for short tests. Prefer a `.pgpass` file or a secret manager. Do not put the password in a shared script.

A GUI helps you browse trees of databases and tables. Common tools:

- **pgAdmin** — official web GUI
- **DBeaver** — free multi-database GUI
- **DataGrip** — commercial IDE

The GUI still sends SQL to the same server. Learn `psql` even if you use a GUI. Error text, `EXPLAIN`, and scripts are easier in `psql`.

`psql` reads `~/.psqlrc` on start (or `%APPDATA%\postgresql\psqlrc.conf` on Windows). You can set prompts and defaults there. Do not hide errors in that file.

Check the client version:

```text
psql --version
```

A newer `psql` can connect to an older server for many tasks. A much older `psql` can miss new SQL. Match client and server major versions when you can.

### Questions

#### Theoretical questions

1. What is `psql`?
2. Why must you still learn `psql` when you use a GUI?
3. Name three GUI tools from this section.
4. Why is a password in a shared script a problem?
5. Why match the `psql` major version with the server major version?

#### Easy practical tasks

1. Run `psql --version`. Save the full output.
2. Connect with `psql` to the default `postgres` database. Run `SELECT 1;`. Disconnect with `\q`.
3. Open pgAdmin, DBeaver, or DataGrip. Create a connection to the same server. Run `SELECT 1`.
4. Write the `psql` command line that sets host, port, user, and database.

#### Medium practical tasks

1. Use `psql -c "SELECT current_user;"` without an interactive session. Save the output.
2. Compare the object tree in a GUI with `\l` and `\dt` in `psql`. Write three objects that both views show.
3. Create a `.pgpass` file (or the Windows `pgpass.conf`). Connect without a password prompt. Document the file location and the file mode.

#### Advanced practical tasks

1. Write a `psql` script file that runs two `SELECT` statements. Run it with `psql -f`. Show the output.
2. Configure SSL or a non-default port in both `psql` and one GUI. Record every connection field that you set.

---

## Clusters, databases, schemas, roles

A **cluster** is one data directory initialized by `initdb`. One cluster has one running server. One server can contain many databases.

A **database** is a named catalog inside the cluster. Objects in one database are not visible in another database on the same cluster. You connect to one database at a time.

A **schema** is a namespace inside a database. Tables, views, and functions have a schema name. The default schema for new objects is often `public`. Topic 2 covers `search_path`.

A **role** is an account in the cluster. A role can log in when it has `LOGIN`. A role without `LOGIN` is a group. Older text says "user" for a login role. PostgreSQL stores both as roles. `CREATE USER` is `CREATE ROLE` with `LOGIN`.

System databases after `initdb`:

- `postgres` — default database for admin work
- `template1` — source for new databases
- `template0` — clean template; do not put objects here

System schemas include `pg_catalog` (system catalogs) and `information_schema` (SQL-standard views). Do not put application tables in `pg_catalog`.

The hierarchy is: cluster → database → schema → table. A role belongs to the cluster. Privileges apply to objects inside a database. Topic 12 covers `GRANT` in detail.

### Questions

#### Theoretical questions

1. What is a cluster in PostgreSQL?
2. Can one `SELECT` read tables from two databases in the same cluster without a foreign server? Why?
3. What is a schema?
4. What is the difference between a login role and a group role?
5. What is the purpose of `template1`?

#### Easy practical tasks

1. Draw the hierarchy cluster → database → schema → table. Add `public` and one application table.
2. Write one sentence each for `postgres`, `template0`, and `template1`.
3. List three object types that live in a schema.
4. Explain in four sentences why a role is a cluster object and a table is a database object.

#### Medium practical tasks

1. After you can connect, run queries that list databases, schemas, and roles (topic 2 shows the commands). Save the three lists.
2. Create a second database. Connect to it. Confirm that a table from the first database is not in `\dt`.
3. Create a schema `app`. Create a table `app.items`. Write the qualified name of the table.

#### Advanced practical tasks

1. Read the official docs on `CREATE DATABASE` and templates. Create a database from `template0`. Write why an admin uses `template0`.
2. Design names for one cluster, two databases (app and analytics), and three schemas. Write the rules you used (length, prefix, no reserved words).

---

## `createdb`, `createuser`, connection URIs

`createdb` and `createuser` are client programs. They send SQL to the server. You can run the same work in `psql`:

```sql
CREATE DATABASE shop;
CREATE ROLE shop_app LOGIN PASSWORD 'secret';
GRANT CONNECT ON DATABASE shop TO shop_app;
```

`createdb shop` is the same idea as `CREATE DATABASE shop`. `createuser shop_app` is the same idea as `CREATE ROLE shop_app LOGIN`. Use SQL when you need full options. Use the programs for short shell scripts.

A connection URI (also called a connection URL) is one string:

```text
postgresql://shop_app:secret@localhost:5432/shop
```

URI parts:

- scheme: `postgresql://` or `postgres://`
- user and optional password
- host and optional port (default `5432`)
- database name
- optional query parameters (`sslmode`, `connect_timeout`)

Key-value form is also valid:

```text
postgresql://?host=localhost&port=5432&dbname=shop&user=shop_app
```

`psql` accepts a URI:

```text
psql "postgresql://shop_app@localhost:5432/shop"
```

Libpq also reads environment variables: `PGHOST`, `PGPORT`, `PGUSER`, `PGDATABASE`. Topic 17 covers client parameters in more detail.

Do not commit a URI that contains a password. Use a secret store or `.pgpass`.

### Questions

#### Theoretical questions

1. What SQL statement does `createdb` send?
2. What SQL statement does `createuser` send?
3. Name the parts of `postgresql://user@host:5432/dbname`.
4. What environment variables set host, port, user, and database?
5. Why must a URI with a password stay out of git?

#### Easy practical tasks

1. Create a database with `createdb` or `CREATE DATABASE`. List databases and find the new name.
2. Create a login role with `createuser` or `CREATE ROLE`. List roles and find the new name.
3. Write a connection URI for host `127.0.0.1`, port `5432`, database `shop`, user `shop_app`. Do not include a real production password.
4. Connect with `psql` and that URI (password prompt is acceptable).

#### Medium practical tasks

1. Create a database and a role. Grant `CONNECT`. Connect as that role. Run `SELECT current_user, current_database();`.
2. Connect once with a URI and once with `-h -p -U -d`. Write a table of equivalent fields.
3. Set `PGHOST`, `PGPORT`, `PGUSER`, and `PGDATABASE` in your shell. Run `psql` with no extra flags. Then unset the variables.

#### Advanced practical tasks

1. Read the libpq URI docs. Add `sslmode=prefer` and `connect_timeout=5` to a URI. Show a successful connection.
2. Write a small shell script that creates a database and a role from arguments. The script must not echo the password.

---

## Official docs: postgresql.org/docs

The primary documentation is [https://www.postgresql.org/docs/](https://www.postgresql.org/docs/). Open the major version that matches your server. This path uses 16 and 17:

- [https://www.postgresql.org/docs/16/](https://www.postgresql.org/docs/16/)
- [https://www.postgresql.org/docs/17/](https://www.postgresql.org/docs/17/)
- [https://www.postgresql.org/docs/current/](https://www.postgresql.org/docs/current/) (newest stable major)

The official tutorial is [https://www.postgresql.org/docs/current/tutorial.html](https://www.postgresql.org/docs/current/tutorial.html). Complete the tutorial on your 16 or 17 server.

Useful books inside the docs:

- Tutorial — first SQL and `psql`
- SQL Language — statements and types
- Server Administration — install, config, backup
- Reference — exact syntax for each command

The wiki is [https://wiki.postgresql.org/](https://wiki.postgresql.org/). The wiki has tips. The official docs win when the two sources disagree.

`psql` can open help for a SQL command:

```text
\h CREATE TABLE
\?
```

`\h` shows SQL syntax. `\?` shows `psql` meta-commands.

Use the reference page when you need exact syntax. Use the tutorial when you learn the first time. Use release notes when a major upgrade changes behavior.

### Questions

#### Theoretical questions

1. Why must the docs URL include the major version?
2. What is the official tutorial URL pattern?
3. When do you trust the official docs over the wiki?
4. What does `\h CREATE TABLE` show?
5. What does `\?` show?

#### Easy practical tasks

1. Open the 17 tutorial. Complete the first page. Write one fact that you learned.
2. Open the reference page for `SELECT`. Write the purpose of the `FROM` clause in one sentence.
3. Run `\h CREATE DATABASE` in `psql`. Write the required argument.
4. Bookmark docs for 16, 17, and `current`.

#### Medium practical tasks

1. Find the "Conventions" page in the docs. Write how the docs mark optional syntax.
2. Open `version()` function reference. Write the return type.
3. Compare one reference page in 16 and 17 (example: `CREATE INDEX`). Write one difference or write that the page is the same.

#### Advanced practical tasks

1. Map this handbook topic to official doc chapters. List chapter titles for install, `psql`, and `CREATE DATABASE`.
2. Read "Bug Reporting Guidelines" in the docs. Write a five-line template for a reproducible report.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from an empty machine to `SELECT 1` in `psql`. Name install, start, and connect.
2. Why is a cluster not the same object as a database?
3. How do `createdb` and a connection URI work together in a new project?
4. What risk do you take when `docs/current` and your server major version are not the same?
5. A teammate wants to store only JSON and skip tables. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. Start PostgreSQL 16 or 17. Connect with `psql`. Run `SELECT version();` and `SELECT current_database();`. Save both results.
2. Write a one-page cheat sheet: cluster, database, schema, role, `psql`, URI, major versus minor.
3. Create a database `learn` and a login role `learn_app`. Connect as `learn_app` to `learn` if privileges allow, or document the extra `GRANT` that you need.
4. Open the official download page and the Docker Hub `postgres` page. Write one install command for each method.

#### Medium practical tasks

1. Write a short script (PowerShell or bash) that waits until port `5432` accepts connections, then runs `psql -c "SELECT 1"`.
2. Create two databases and two roles. Draw who can connect to what. Test each pair with `psql`.
3. Document your install in ten steps so that another beginner can copy it. Include version numbers.

#### Advanced practical tasks

1. Run two major versions at the same time (16 and 17) on different ports. Connect to each. Save `SHOW server_version` for both.
2. Read the `initdb` reference. List five options that change a new cluster. State the default for each option that has a default.
