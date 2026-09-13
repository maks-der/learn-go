# 1. Getting Started

## Description

PostgreSQL is an open-source object-relational database system. This topic shows what PostgreSQL is, why the major version matters, and how you install a server. You also learn clusters, databases, schemas, roles, and how a client connects.

Learn common database ideas first (`db.topics.md`). Then complete this topic. Complete this topic before you study types, tables, and SQL.

Use one term for each concept. A cluster is one data directory with one running server. A database lives in a cluster. A schema groups objects inside a database. A role is a login name or a group name. `postgresql.conf` sets server parameters. `pg_hba.conf` sets who may connect. `search_path` is the list of schemas that unqualified names use. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## What PostgreSQL is

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

## Major versions

PostgreSQL uses a major.minor version number. Examples: `16.15` and `17.11`. The first number is the major version. The rest is the minor version.

A minor release fixes defects and security problems. You can apply a minor release without a dump and restore of the data directory. Stop the server. Install the new binaries. Start the server.

A major release can change catalogs, on-disk format, and SQL behavior. You cannot point a PostgreSQL 17 binary at a PostgreSQL 16 data directory. You must upgrade with `pg_upgrade`, dump and restore, or logical replication. Topic 11 covers physical backup. Topic 14 covers upgrade practice.

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

`psql` is the official terminal client. It ships with the PostgreSQL client package. Connect with host, port, user, and database:

```text
psql -h localhost -p 5432 -U postgres -d postgres
```

A newer `psql` can connect to an older server for many tasks. A much older `psql` can miss new SQL. Match client and server major versions when you can.

Do not expose port `5432` to the public internet. Topic 14 covers production network rules.

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
4. Connect with `psql` to the default `postgres` database. Run `SELECT 1;`. Disconnect with `\q`.

#### Medium practical tasks

1. Start PostgreSQL with a named Docker volume. Stop and remove the container. Start a new container on the same volume. Confirm that the data directory still exists.
2. Install or run both 16 and 17 on different ports (`5432` and `5433`). Write the two listen addresses.
3. Compare the official installer with the Docker image. Write a short report: data location, how you stop the server, and how you upgrade the minor version.

#### Advanced practical tasks

1. Build a Docker Compose file with `postgres:17`, a named volume, and a health check. Start it. Show that the health check passes.
2. Read the official Docker image documentation. Write how the image runs scripts from `/docker-entrypoint-initdb.d`. Create one `.sql` file and prove that it runs only on the first init.

---

## Clusters, databases, schemas, and roles

A **cluster** is one data directory initialized by `initdb`. One cluster has one running server. One server can contain many databases.

A **database** is a named catalog inside the cluster. Objects in one database are not visible in another database on the same cluster. You connect to one database at a time.

A **schema** is a namespace inside a database. Tables, views, and functions have a schema name. The default schema for new objects is often `public`.

A **role** is an account in the cluster. A role can log in when it has `LOGIN`. A role without `LOGIN` is a group. Older text says "user" for a login role. PostgreSQL stores both as roles. `CREATE USER` is `CREATE ROLE` with `LOGIN`.

System databases after `initdb`:

- `postgres` — default database for admin work
- `template1` — source for new databases
- `template0` — clean template; do not put objects here

System schemas include `pg_catalog` (system catalogs) and `information_schema` (SQL-standard views). Do not put application tables in `pg_catalog`.

The hierarchy is: cluster → database → schema → table. A role belongs to the cluster. Privileges apply to objects inside a database. Topic 8 covers `GRANT` in detail.

`createdb` and `createuser` are client programs. They send SQL to the server. You can run the same work in `psql`:

```sql
CREATE DATABASE shop;
CREATE ROLE shop_app LOGIN PASSWORD 'secret';
GRANT CONNECT ON DATABASE shop TO shop_app;
```

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
4. Create a database with `createdb` or `CREATE DATABASE`. List databases and find the new name.

#### Medium practical tasks

1. Create a second database. Connect to it. Confirm that a table from the first database is not in `\dt`.
2. Create a schema `app`. Create a table `app.items`. Write the qualified name of the table.
3. Create a login role. Grant `CONNECT`. Connect as that role. Run `SELECT current_user, current_database();`.

#### Advanced practical tasks

1. Read the official docs on `CREATE DATABASE` and templates. Create a database from `template0`. Write why an admin uses `template0`.
2. Design names for one cluster, two databases (app and analytics), and three schemas. Write the rules you used (length, prefix, no reserved words).

---

## Connection URIs, `pg_hba.conf`, and `postgresql.conf` (high-level)

`postgresql.conf` is the main server configuration file. It lives in the data directory by default. It sets listen address, port, memory, logging, and many other parameters. Example lines:

```text
listen_addresses = 'localhost'
port = 5432
max_connections = 100
```

You can also put overrides in `postgresql.auto.conf`. `ALTER SYSTEM` writes that file. Do not edit `postgresql.auto.conf` by hand.

`pg_hba.conf` is the host-based authentication file. Each line is a rule. The server uses the first matching rule. A rule has connection type, database, role, client address, and method.

Example rules (do not copy them to production without a review):

```text
# TYPE  DATABASE  USER  ADDRESS      METHOD
local   all       all                scram-sha-256
host    all       all   127.0.0.1/32 scram-sha-256
host    all       all   ::1/128      scram-sha-256
```

Common methods:

- `scram-sha-256` — password with SCRAM (default password method in current versions)
- `md5` — older password method; avoid for new work
- `trust` — no password; use only on a private local test
- `peer` — local socket; the OS user name must match the role
- `reject` — deny the match

`postgresql.conf` answers "how the server runs". `pg_hba.conf` answers "who may connect". A wrong listen address blocks the network. A wrong `pg_hba.conf` rule rejects the login after the TCP connect.

Reload many settings with `SELECT pg_reload_conf();` or `pg_ctl reload`. Some parameters need a restart. The docs mark each parameter as `SIGHUP` or postmaster.

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

`psql` accepts a URI:

```text
psql "postgresql://shop_app@localhost:5432/shop"
```

Libpq also reads environment variables: `PGHOST`, `PGPORT`, `PGUSER`, `PGDATABASE`. Topic 12 covers client parameters in more detail.

Do not set `listen_addresses = '*'` and `trust` together. Do not commit a URI that contains a password.

### Questions

#### Theoretical questions

1. Which settings belong in `postgresql.conf`?
2. Which rules belong in `pg_hba.conf`?
3. What does the server do when two `pg_hba.conf` rules match?
4. Name the parts of `postgresql://user@host:5432/dbname`.
5. Which file does `ALTER SYSTEM` change?

#### Easy practical tasks

1. Find `postgresql.conf` and `pg_hba.conf` on your install or in the Docker data volume. Write both full paths.
2. In `postgresql.conf`, find `port` and `listen_addresses`. Write the values.
3. Write a connection URI for host `127.0.0.1`, port `5432`, database `shop`, user `shop_app`. Do not include a real production password.
4. Make a two-column table: "File" and "Question it answers". Add one row per file.

#### Medium practical tasks

1. Change a reloadable parameter (example: `log_min_duration_statement`). Reload. Confirm with `SHOW`. Restore the old value.
2. Connect once with a URI and once with `-h -p -U -d`. Write a table of equivalent fields.
3. Read the official docs for `pg_hba.conf`. Write the meaning of `local` versus `host`.

#### Advanced practical tasks

1. Create a role that must use `scram-sha-256` from `127.0.0.1`. Prove that a wrong password fails. Record the client error. Do not write an attack procedure.
2. Find three parameters that need a restart. Write their names and why a reload is not enough (from the docs).

---

## `psql` meta-commands: `\l`, `\c`, `\dt`, `\d`, `\dn`, `\du`

Meta-commands start with a backslash. They are `psql` commands. The server does not receive the backslash line as SQL. `psql` sends SQL for you.

| Command | Purpose |
| --- | --- |
| `\l` | list databases |
| `\c dbname` | connect to another database |
| `\dt` | list tables in the current `search_path` |
| `\d name` | describe a table or other object |
| `\dn` | list schemas |
| `\du` | list roles |

Useful variants:

- `\l+` — databases with size and comment
- `\dt *.*` — tables in all schemas
- `\dt schema.*` — tables in one schema
- `\d+ name` — extra detail (indexes, size comments)
- `\du+` — roles with description
- `\x` — toggle expanded display
- `\q` — quit
- `\conninfo` — show the current connection
- `\h CREATE TABLE` — SQL syntax help
- `\?` — `psql` meta-command help

Examples:

```text
\l
\c learn
\dn
\dt
\d public.items
\du
```

`\c` changes the database. You cannot change database with a SQL `USE` command. PostgreSQL has no `USE`. You reconnect.

`\dt` without a pattern hides tables that are not on `search_path`. A table in schema `app` does not appear until you set the path or you run `\dt app.*`.

`\d` without a name lists more object types than `\dt`. Prefer `\dt` when you want tables only.

These commands help you confirm that you are on the correct database and that your role can see the objects.

### Questions

#### Theoretical questions

1. Is `\l` SQL? Where does it run?
2. How do you change database in PostgreSQL?
3. Why can `\dt` hide a table that exists?
4. What does `\d tablename` show?
5. What does `\du` list?

#### Easy practical tasks

1. Run `\l`, `\dn`, `\dt`, and `\du`. Save each output to a text file.
2. Use `\c` to switch between `postgres` and another database. Run `\conninfo` after each switch.
3. Create a table `demo_meta (id int)`. Run `\dt` and `\d demo_meta`.
4. Run `\?` and find `\l`. Write the help line.

#### Medium practical tasks

1. Create a table in a non-`public` schema. Show that `\dt` misses it. Show that `\dt schema.*` finds it.
2. Compare `\d` and `\d+` on the same table. List three extra lines from `\d+`.
3. Use `\l+` and write the size of each database on your cluster.

#### Advanced practical tasks

1. Run `psql -E` (echo hidden SQL). Run `\dt`. Copy the SQL that `psql` sends. Run that SQL by hand.
2. Write a one-page `psql` cheat sheet with ten meta-commands. Include the six from this section and four more from `\?`.

---

## `search_path` and the `public` schema

`search_path` is a list of schemas. PostgreSQL uses this list to find unqualified names. An unqualified name is `items`. A qualified name is `app.items`.

The default `search_path` is `"$user", public`. `"$user"` means a schema with the same name as the current role. If that schema does not exist, PostgreSQL skips it. Then it uses `public`.

Show the path:

```sql
SHOW search_path;
SELECT current_schemas(true);
```

`current_schemas(true)` includes implicit schemas such as `pg_catalog`. System functions resolve even when you do not write `pg_catalog`.

Set the path for the session:

```sql
SET search_path TO app, public;
```

Set a default for a role:

```sql
ALTER ROLE shop_app SET search_path TO app, public;
```

New objects go to the first schema in `search_path` that exists. If you omit a schema in `CREATE TABLE items`, PostgreSQL creates `public.items` when `public` is the first existing schema.

The `public` schema exists in new databases. In PostgreSQL 15 and later, the `PUBLIC` role does not have `CREATE` on `public` by default. A normal login role needs `GRANT CREATE ON SCHEMA public` or a private schema. PostgreSQL 16 and 17 keep that safer default.

Prefer an application schema (`app` or `shop`) for new work. Qualify names in scripts when two schemas can hold the same object name.

Do not put application tables in `pg_catalog`. Do not rely on `public` `CREATE` for every role. Do not leave two tables with the same name on `search_path` without a plan. The first match wins.

### Questions

#### Theoretical questions

1. What is an unqualified name?
2. What does `"$user"` mean in `search_path`?
3. Where does `CREATE TABLE items` put the table?
4. Does `PUBLIC` have `CREATE` on `public` in PostgreSQL 16 and 17 by default?
5. What does `current_schemas(true)` add that `SHOW search_path` may hide?

#### Easy practical tasks

1. Run `SHOW search_path;` and `SELECT current_schemas(true);`. Save both results.
2. Create schema `app`. `SET search_path TO app, public;`. Create table `items`. Confirm the schema with `\d app.items`.
3. Write a qualified name and an unqualified name for the same table.
4. Grant `USAGE` on `app` to a login role. Connect as that role. Try `SELECT` from `app.items` after you grant table rights.

#### Medium practical tasks

1. Create the same table name in `public` and in `app`. Change `search_path`. Show which table `SELECT * FROM items` reads.
2. `ALTER ROLE` a practice role so that `search_path` is `app, public`. Reconnect as that role. Confirm with `SHOW`.
3. As a non-superuser, try `CREATE TABLE` in `public` without extra grants. Record the error. Then create in `app`.

#### Advanced practical tasks

1. Read the 16 or 17 docs on `search_path` and schema privileges. Write a four-step policy for a new database that does not use `public` for app tables.
2. Find which objects ignore `search_path` (example: some system catalogs). Write two examples from the docs.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from an empty machine to `SELECT 1` in `psql`. Name install, start, and connect.
2. Why is a cluster not the same object as a database?
3. How do a connection URI and `pg_hba.conf` work together when a login fails?
4. What risk do you take when `docs/current` and your server major version are not the same?
5. A teammate wants to store only JSON and skip tables. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. Start PostgreSQL 16 or 17. Connect with `psql`. Run `SELECT version();` and `SELECT current_database();`. Save both results.
2. Write a one-page cheat sheet: cluster, database, schema, role, `psql` meta-commands, URI, `search_path`, major versus minor.
3. Create a database `learn` and a login role `learn_app`. Connect as `learn_app` to `learn` if privileges allow, or document the extra `GRANT` that you need.
4. Open the official download page and the Docker Hub `postgres` page. Write one install command for each method.

#### Medium practical tasks

1. Write a short script (PowerShell or bash) that waits until port `5432` accepts connections, then runs `psql -c "SELECT 1"`.
2. Create two databases and two roles. Draw who can connect to what. Test each pair with `psql`.
3. Document your install in ten steps so that another beginner can copy it. Include version numbers and the paths of `postgresql.conf` and `pg_hba.conf`.

#### Advanced practical tasks

1. Run two major versions at the same time (16 and 17) on different ports. Connect to each. Save `SHOW server_version` for both.
2. Read the `initdb` reference. List five options that change a new cluster. State the default for each option that has a default.
