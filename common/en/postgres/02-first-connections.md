# 2. First Connections

## Description

This topic shows how a client reaches a PostgreSQL 16 or 17 server. You learn the two main configuration files, the connection fields, and the first `psql` meta-commands. You also learn `search_path` and how you read session settings.

Complete topic 1 first. You need a running server and a working `psql`. Complete this topic before you write application SQL.

Use one term for each concept. `postgresql.conf` sets server parameters. `pg_hba.conf` sets who may connect and how the server authenticates the client. A session is one connection. `search_path` is the list of schemas that unqualified names use.

---

## `postgresql.conf` vs `pg_hba.conf` (high-level)

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

Do not set `listen_addresses = '*'` and `trust` together. Topic 21 covers production network rules.

### Questions

#### Theoretical questions

1. Which settings belong in `postgresql.conf`?
2. Which rules belong in `pg_hba.conf`?
3. What does the server do when two `pg_hba.conf` rules match?
4. What is `scram-sha-256`?
5. Which file does `ALTER SYSTEM` change?

#### Easy practical tasks

1. Find `postgresql.conf` and `pg_hba.conf` on your install or in the Docker data volume. Write both full paths.
2. In `postgresql.conf`, find `port` and `listen_addresses`. Write the values.
3. Count the non-comment lines in `pg_hba.conf`. Write the count.
4. Make a two-column table: "File" and "Question it answers". Add one row per file.

#### Medium practical tasks

1. Change a reloadable parameter (example: `log_min_duration_statement`). Reload. Confirm with `SHOW`. Restore the old value.
2. Add a comment in `pg_hba.conf` that describes your local rule. Reload. Connect again to prove that you did not break login.
3. Read the official docs for `pg_hba.conf`. Write the meaning of `local` versus `host`.

#### Advanced practical tasks

1. Create a role that must use `scram-sha-256` from `127.0.0.1`. Prove that a wrong password fails. Record the client error. Do not write an attack procedure.
2. Find three parameters that need a restart. Write their names and why a reload is not enough (from the docs).

---

## Host, port `5432`, database, user, password

A client needs these fields:

- **host** — name or IP of the server (`localhost`, `127.0.0.1`, or a container name)
- **port** — TCP port; default `5432`
- **database** — the database name; not the cluster name
- **user** — the role with `LOGIN`
- **password** — the role password when the `pg_hba.conf` method needs one

`psql` flags:

```text
psql -h localhost -p 5432 -U postgres -d postgres
```

Libpq environment variables:

- `PGHOST`
- `PGPORT`
- `PGUSER`
- `PGDATABASE`
- `PGPASSWORD` (avoid in shared shells; prefer `.pgpass`)

A Unix-domain socket is a local connection without TCP. On Linux, `psql` without `-h` often uses the socket. `pg_hba.conf` type `local` applies. On Windows, clients usually use TCP.

The host must match `listen_addresses`. The port must match `port`. The database must exist. The role must exist and must have `CONNECT` on that database. The password must satisfy the `pg_hba.conf` method.

Error messages:

- connection refused — nothing listens on that host and port
- no `pg_hba.conf` entry — TCP worked; authentication rules rejected the client
- password authentication failed — the rule matched; the password is wrong
- database does not exist — the name is wrong or you connected to another cluster

Do not confuse the OS user with the PostgreSQL role. They can share a name. They are not the same account.

### Questions

#### Theoretical questions

1. What are the five connection fields in this section?
2. What is the default port?
3. When does `psql` use a Unix-domain socket?
4. What does "connection refused" mean?
5. Why is the OS user not the same object as a PostgreSQL role?

#### Easy practical tasks

1. Connect with all five fields set on the command line. Run `SELECT 1;`.
2. Write a URI that uses the same five fields (use a dummy password).
3. Set `PGHOST`, `PGPORT`, `PGUSER`, and `PGDATABASE`. Connect with `psql` and no flags.
4. Make a table: "Error text" and "What you check first". Add four rows from this section.

#### Medium practical tasks

1. Connect to the same server with `127.0.0.1` and with `localhost`. Write whether both work. If one fails, write the `pg_hba.conf` reason.
2. Try port `5433` on purpose. Record the error. Connect again on the real port.
3. Create a login role with a password. Connect as that role. Then try a wrong password. Record both outcomes.

#### Advanced practical tasks

1. Connect from another machine or from a second container on a user-defined Docker network. Write every field that you changed. If you cannot use a second host, document the limit and use two containers.
2. Read the libpq docs for `sslmode`. Try `sslmode=disable` and `sslmode=prefer` on your local server. Write what each mode does for your server.

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

Do not put application tables in `pg_catalog`. Do not rely on `public` as the only security boundary. Topic 12 covers privileges.

Two tables with the same name can live in two schemas. `search_path` picks the first match. Qualify names in scripts when the path can change.

### Questions

#### Theoretical questions

1. What is an unqualified table name?
2. What does `"$user"` mean in `search_path`?
3. Where does `CREATE TABLE items` put the table?
4. Why is `CREATE` on schema `public` not granted to `PUBLIC` by default in PostgreSQL 15 and later?
5. How does `search_path` choose between `app.items` and `public.items`?

#### Easy practical tasks

1. Run `SHOW search_path;` and `SELECT current_schemas(true);`. Write both results.
2. Create schema `app`. `SET search_path TO app, public`. Create table `items`. Confirm the schema with `\d app.items`.
3. Reset the path with `SET search_path TO DEFAULT`. Show the value again.
4. Write four sentences: qualified name, unqualified name, first schema, `public`.

#### Medium practical tasks

1. Create `app.items` and `public.items`. Change `search_path` so that `SELECT * FROM items` reads `app.items`. Then change the path so that it reads `public.items`.
2. `ALTER ROLE` for your login role to set `search_path`. Reconnect. Confirm `SHOW search_path`.
3. As a non-superuser, try `CREATE TABLE` in `public` without extra grants. Record the error or the success. Write the grant that you need if it fails.

#### Advanced practical tasks

1. Read the docs for `current_schemas`. Explain the boolean argument. Show both `true` and `false`.
2. Design `search_path` for an application role and a report role. Write `ALTER ROLE` statements and one risk if a developer omits the schema name.

---

## `SHOW` and `current_database()`

`SHOW` reads a configuration parameter for the current session:

```sql
SHOW port;
SHOW listen_addresses;
SHOW server_version;
SHOW search_path;
SHOW TimeZone;
SHOW work_mem;
```

`SHOW ALL` lists every parameter. The list is long. Filter in `psql` or query `pg_settings`:

```sql
SELECT name, setting, unit, context
FROM pg_settings
WHERE name IN ('port', 'work_mem', 'search_path');
```

`context` tells you if you can change the value in a session, with a reload, or only at start.

Session functions:

```sql
SELECT current_database();
SELECT current_user;
SELECT session_user;
SELECT current_schema();
SELECT inet_server_addr(), inet_server_port();
SELECT inet_client_addr(), inet_client_port();
```

`current_database()` is the database of this connection. `current_user` is the role for privilege checks. `session_user` is the original login role. `SET ROLE` can change `current_user`.

`current_schema()` is the first writable schema in `search_path`. It is not the same as `current_database()`.

These functions confirm that you connected to the server that you intended. Run them before you run `DELETE`.

`SET` changes a session parameter when the parameter allows it:

```sql
SET TimeZone TO 'UTC';
SHOW TimeZone;
RESET TimeZone;
```

Do not confuse `SHOW` with `SELECT` on a table. `SHOW` is for parameters.

### Questions

#### Theoretical questions

1. What does `SHOW` display?
2. What does `current_database()` return?
3. What is the difference between `current_user` and `session_user`?
4. What does `current_schema()` return?
5. Where do you read the `context` of a parameter?

#### Easy practical tasks

1. Run `SHOW server_version;`, `SHOW port;`, and `SELECT current_database();`. Save the results.
2. Run `SELECT current_user, session_user, current_schema();`.
3. Run `SHOW ALL` and find `shared_buffers`. Write the value.
4. Run `SET TimeZone TO 'UTC';` then `SHOW TimeZone;`. Reset it.

#### Medium practical tasks

1. Query `pg_settings` for five parameters that you use in this topic. Include `context`.
2. Connect to two databases. Run `current_database()` in each session. Write both names.
3. Read the docs for `SET ROLE`. If you have a second role, switch and compare `current_user` and `session_user`.

#### Advanced practical tasks

1. Find three parameters with `context = 'postmaster'` and three with `context = 'user'`. Write what that means for change method.
2. Write a `psql` startup snippet (or `.psqlrc` lines) that prints database, user, server version, and `search_path` on connect.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. A TCP connection succeeds and login fails. Which file do you read first, and why?
2. Describe how host, port, `pg_hba.conf`, and `search_path` each affect a first `SELECT`.
3. Why does PostgreSQL have no `USE database` statement?
4. How do `\conninfo` and `current_database()` give different kinds of information?
5. What is the difference between a session `SET` and a change in `postgresql.conf`?

#### Easy practical tasks

1. Connect, then run `\conninfo`, `\l`, `\dn`, `SHOW search_path;`, and `SELECT current_database();`. Save the session log.
2. Write a cheat sheet with the six meta-commands from this topic and the five connection fields.
3. Create schema `class`, set `search_path`, create one table, describe it with `\d`.
4. Export `SHOW ALL` to a file. Highlight `port`, `listen_addresses`, `search_path`, and `server_version`.

#### Medium practical tasks

1. Break a connection on purpose with a wrong database name, then a wrong port, then a wrong password. Write the three error texts and the first fix for each. Restore a working connection.
2. Document your `pg_hba.conf` local rules in plain language (type, address, method). Do not weaken production rules.
3. Write ten steps that another beginner can follow to connect with `psql` and prove the session with `SHOW` and `current_database()`.

#### Advanced practical tasks

1. Use `psql -E` on `\l`, `\dt`, and `\du`. Save the hidden SQL. Write one sentence per command about which catalog it reads.
2. Configure a dedicated login role that can connect only to one database from `127.0.0.1`. Prove that a connection to a second database fails. Use `pg_hba.conf` and privileges as needed.
