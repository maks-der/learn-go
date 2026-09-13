# 14. Production and Ecosystem

## Description

This topic shows how you run PostgreSQL 16 and PostgreSQL 17 in production and where you continue after this path. You learn a least-privilege application role, network rules, migrations, expand-contract schema changes, monitoring, a short survey of window functions, foreign data wrappers, and PostGIS, plus official docs, cloud products, and bug reports.

Complete topic 13 first. You can read activity views and memory settings. After this topic, use a real database and keep reading.

Use one term for each concept. Least privilege means the role has only the rights that the program needs. A migration is a versioned schema change. Expand-contract is a compatible change in small steps. A window function computes a value from a related set of rows without collapsing the result. A foreign data wrapper (FDW) reads a remote source as a table. The official docs are the pages on postgresql.org for a major version. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## Least-privilege app role and network exposure

Topic 8 introduced roles and `GRANT`. Production makes that the default, not an extra.

The application role (`shop_app`):

- has `LOGIN`
- has `CONNECT` on the app database
- has `USAGE` on the app schema
- has DML on app tables (`SELECT`, `INSERT`, `UPDATE`, `DELETE` as needed)
- has sequence `USAGE` for identity inserts
- has `EXECUTE` only on functions that the app must call
- does not have `SUPERUSER`, `BYPASSRLS`, `CREATEDB`, `CREATEROLE`, or `REPLICATION`
- does not own tables if a migrator role owns them

The migrator role runs DDL. CI uses the migrator secret. The app never sees that secret.

```sql
CREATE ROLE shop_migrator LOGIN PASSWORD 'migrator-secret';
CREATE ROLE shop_app LOGIN PASSWORD 'app-secret';

GRANT CONNECT ON DATABASE shop TO shop_migrator, shop_app;
GRANT USAGE, CREATE ON SCHEMA shop TO shop_migrator;
GRANT USAGE ON SCHEMA shop TO shop_app;
```

Row-level security stays on if you are multi-tenant. The app role must not bypass RLS.

PostgreSQL listens on TCP port `5432` by default. If that port is reachable on the public internet, attackers scan it. SSL does not make a public `5432` a good idea.

Safe pattern:

- `listen_addresses` binds only the interfaces that need it
- a firewall or security group allows `5432` only from the app subnet, PgBouncer, or a bastion
- `pg_hba.conf` matches those sources and requires `scram-sha-256` (and `hostssl` when you use TLS)
- clients use `sslmode=verify-full` when you control certificates
- administrators use a VPN, a bastion, or a vendor private link, not an open `0.0.0.0/0` rule

Docker `-p 5432:5432` on a laptop is a lab habit. On a cloud VM, a published `5432` to `0.0.0.0` is an incident.

Do not debug a production outage by granting `SUPERUSER` to `shop_app`. Do not put PostgreSQL on the public internet on a host that holds real data. Do not use `trust` for any remote line.

### Questions

#### Theoretical questions

1. What privileges must the application role lack?
2. Who owns tables in the migrator-and-app split?
3. Why is a public `5432` a problem even with SSL?
4. Who may reach `5432` in the safe pattern?
5. Which `pg_hba` method do you use for remote passwords?

#### Easy practical tasks

1. `\du` your lab `shop_app`. Confirm `Superuser` is no.
2. `SHOW listen_addresses;` `SHOW port;`.
3. Connect as `shop_app` and try `CREATE TABLE`. Record the error.
4. Write four sentences: app role, migrator, firewall, `hostssl`.

#### Medium practical tasks

1. Create `shop_migrator` and `shop_app`. Apply defaults (topic 8). Create a table as migrator. Prove the app can `INSERT` and cannot `DROP TABLE`.
2. Draw the path: laptop admin → VPN → bastion → PostgreSQL. Label ports.
3. Document how PgBouncer changes the firewall (app to pooler, pooler to database).

#### Advanced practical tasks

1. Script a production bootstrap: database, schema, roles, defaults, RLS on one tenant table, no superuser grants. Run it on an empty lab database.
2. Write a network checklist for a new environment: bind, security group, `pg_hba`, TLS, no public IP.

---

## Migrations and expand-contract schema changes

A migration is a versioned SQL (or code) change that you apply in order. Production schema must not be "whatever someone ran in `psql`".

Common tools:

- **Flyway** — version files such as `V3__add_orders.sql`; tracks `flyway_schema_history`
- **golang-migrate** — `0003_add_orders.up.sql` and a down file; common in Go services
- **Sqitch** — plans with names and dependencies; strong for review

All three need:

- one history table
- a migrator role
- a CI step that applies pending files
- a rule that a file that reached production does not change

```text
-- 0004_items_qty.up.sql
ALTER TABLE shop.items ADD COLUMN qty int NOT NULL DEFAULT 0;
```

Down migrations are optional in some teams. If you keep them, test them. A down that drops a column destroys data.

A lock that rewrites a large table can stop the application. Check the current `ALTER TABLE` page: some changes are a catalog-only update; some scan the table; some rewrite.

**Expand-contract** (expand, migrate, contract):

1. Expand: add a new column or table that the old app ignores
2. Dual write or backfill: the app or a job fills the new object
3. Switch reads to the new object
4. Contract: drop the old column or table in a later release

Example: rename `name` to `title` without a long rewrite:

1. Add `title text`
2. Deploy app that writes both
3. Backfill `title` from `name`
4. Deploy app that reads `title`
5. Drop `name`

`CREATE INDEX CONCURRENTLY` (topic 5) is an expand step. `DETACH PARTITION` (topic 10) can drop old data without a huge `DELETE`.

Do not mix two migration tools on the same database. Do not edit an applied file after production applied it. Do not run DDL as `shop_app`. Do not add `NOT NULL` without a default or a prior backfill.

Official `ALTER TABLE`: [https://www.postgresql.org/docs/17/sql-altertable.html](https://www.postgresql.org/docs/17/sql-altertable.html).

### Questions

#### Theoretical questions

1. What does a migration tool store in the database?
2. Why must an applied file stay unchanged?
3. What are the four expand-contract steps?
4. Why can a rename of a column need more than one release?
5. Why is ad-hoc `psql` DDL a problem in production?

#### Easy practical tasks

1. Open the Flyway, golang-migrate, and Sqitch sites. Write the history-table name or equivalent for each (from their docs).
2. Write expand-contract steps for "add `email_verified boolean`".
3. Write one `up` file that adds a nullable column.
4. Open the 17 `ALTER TABLE` page. Find one change that notes a table rewrite.

#### Medium practical tasks

1. Apply two migrations with one tool on a lab database (or simulate with a history table and two SQL files).
2. On a lab table with 10000 rows, add a nullable column. Then add `NOT NULL` after a backfill. Time each step.
3. Plan `CREATE INDEX CONCURRENTLY` plus an app deploy that uses the new filter. Write the order.

#### Advanced practical tasks

1. Design a CI pipeline: lint SQL, apply to a throwaway 16 or 17 database, run app tests, then apply to production with the migrator secret.
2. Plan a `varchar(20)` to `text` change on a large table. Use the `ALTER TABLE` notes. State whether you rewrite or add a column. Write rollback at each expand-contract step.

---

## Monitoring: connections, lag, bloat, wraparound

Production needs measurements, not feelings. Topic 13 showed views. This section is the minimum set.

**Connections**

```sql
SELECT count(*) FROM pg_stat_activity;
SHOW max_connections;
SELECT usename, state, count(*)
FROM pg_stat_activity
GROUP BY 1, 2;
```

Alert when used connections stay high, when many sessions are `idle in transaction`, or when the app cannot connect.

**Replication lag** (topic 11)

```sql
SELECT client_addr, state, replay_lsn, sent_lsn
FROM pg_stat_replication;
```

On the replica, use a lag time function if you have replay. Alert when lag exceeds the report SLA.

**Bloat** (topic 7)

```sql
SELECT relname, n_live_tup, n_dead_tup, last_autovacuum
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC
LIMIT 20;
```

Alert on dead-tuple ratio and on size growth that row counts cannot explain.

**Wraparound** (topic 7)

```sql
SELECT datname, age(datfrozenxid)
FROM pg_database
ORDER BY 2 DESC;
```

Alert far before the wraparound limit. Autovacuum must keep up. A failed freeze is an emergency.

Also watch: disk free, replication slot retained WAL, `pg_stat_statements` outliers, checkpoint spikes, error logs.

Do not alert on every `n_dead_tup > 0`. Do not ignore wraparound because "autovacuum is on". Do not collect SQL text that contains secrets.

### Questions

#### Theoretical questions

1. What connection states must you alert on besides "too many backends"?
2. Where do you read replica lag on the primary?
3. Which two columns hint at table bloat in `pg_stat_user_tables`?
4. What does `age(datfrozenxid)` measure?
5. Why is a stale slot a monitoring target?

#### Easy practical tasks

1. Run the four SQL blocks in this section on your lab. Save the outputs.
2. Write a table: signal, view or `SHOW`, warn when.
3. `SHOW autovacuum;`
4. Write four sentences: connections, lag, dead tuples, freeze age.

#### Medium practical tasks

1. Set numeric warn/critical ideas for connections (percent of `max_connections`), lag (seconds), dead-tuple ratio, and `age(datfrozenxid)`.
2. Join `pg_replication_slots` with a WAL-bytes idea. Write a slot alert.
3. Add `application_name` to the connection group-by. Write which apps hold idle transactions.

#### Advanced practical tasks

1. Write a single monitoring SQL script that prints all four areas plus disk-related sizes (`pg_database_size`, slot lag).
2. Map each alert to a first action (pool, replica, vacuum, wraparound runbook). Keep actions high-level.

---

## Window functions, FDW, PostGIS (survey)

This section is a survey. You do not need every feature for the first production app. You must know the names and when to open the docs.

A **window function** uses `OVER (...)`. The result set keeps one row per input row. `GROUP BY` collapses rows.

```sql
SELECT
    id,
    customer_id,
    total,
    sum(total) OVER (PARTITION BY customer_id) AS customer_sum,
    rank() OVER (PARTITION BY customer_id ORDER BY total DESC) AS rnk,
    lag(total) OVER (PARTITION BY customer_id ORDER BY id) AS prev_total
FROM shop.orders;
```

`PARTITION BY` resets the window. `ORDER BY` inside `OVER` is required for `rank` and `lag`. Topic 3 `DISTINCT ON` can replace a simple `row_number()` filter. Official: [https://www.postgresql.org/docs/17/tutorial-window.html](https://www.postgresql.org/docs/17/tutorial-window.html).

A **foreign data wrapper (FDW)** lets you query a remote source as a foreign table. `postgres_fdw` is the common wrapper for another PostgreSQL server.

```sql
CREATE EXTENSION postgres_fdw;

CREATE SERVER remote_shop
    FOREIGN DATA WRAPPER postgres_fdw
    OPTIONS (host 'remote.example', dbname 'shop', port '5432');
```

You then create a user mapping and foreign tables. Pushdown of `WHERE` clauses is limited. FDW is not a substitute for a replica or for logical replication. Official: [https://www.postgresql.org/docs/17/postgres-fdw.html](https://www.postgresql.org/docs/17/postgres-fdw.html).

**PostGIS** is an optional extension for geometry and geography types, GiST spatial indexes, and functions such as `ST_DWithin`. Many installs do not load it.

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
SELECT PostGIS_Version();
```

If `CREATE EXTENSION postgis` fails, write that the package is absent. Spatial search is not a `numeric` pair plus a B-tree. Official: [https://postgis.net/documentation/](https://postgis.net/documentation/).

Do not store every report as a window if a join to a grouped subquery is clearer. Do not point an FDW at a production primary from an untrusted network. Do not treat PostGIS as required for the rest of this path.

### Questions

#### Theoretical questions

1. How does a window differ from `GROUP BY` in the number of result rows?
2. What does `PARTITION BY` reset?
3. What does `postgres_fdw` let you query?
4. Which index class is usual for PostGIS spatial search?
5. Is PostGIS part of every PostgreSQL install?

#### Easy practical tasks

1. Compute `row_number() OVER (ORDER BY id)` on a table.
2. Open the 17 window tutorial. Write one function name that this section listed.
3. Try `CREATE EXTENSION postgres_fdw` or `postgis`. Write success or the error.
4. Write four sentences: window, FDW, PostGIS, optional.

#### Medium practical tasks

1. Rank orders per customer by `total`. Filter `rnk <= 3` in an outer query.
2. Read `CREATE SERVER` in the `postgres_fdw` docs. Write three `OPTIONS` keys.
3. If PostGIS is present, insert one `geography` point. If it is absent, write a fallback `lat`/`lon` schema and three queries you cannot do correctly.

#### Advanced practical tasks

1. Compare `DISTINCT ON` (topic 3) with `row_number()` for "latest order per customer". Write plans if you can.
2. Read FDW limitations (transactions, locks) and PostGIS `geometry` versus `geography`. Write six sentences that a production review can use.

---

## Official docs, cloud products, and bug reports

The official tutorial is [https://www.postgresql.org/docs/current/tutorial.html](https://www.postgresql.org/docs/current/tutorial.html). The word `current` is the newest stable major version. That version can be newer than 17. Also open the tutorial that matches your server:

- [https://www.postgresql.org/docs/16/tutorial.html](https://www.postgresql.org/docs/16/tutorial.html)
- [https://www.postgresql.org/docs/17/tutorial.html](https://www.postgresql.org/docs/17/tutorial.html)

Release notes matter at every major upgrade. Read the "Migration" section before you move data.

The wiki is [https://wiki.postgresql.org/](https://wiki.postgresql.org/). When the wiki and the official docs disagree, the official docs win. **Postgres Weekly** ([https://postgresweekly.com/](https://postgresweekly.com/)) links articles. It is not a substitute for the docs.

Managed PostgreSQL reduces OS work. You still own schema, grants, RLS, plans, migrations, pool size, and a backup test.

| Product | Vendor | High-level note |
| --- | --- | --- |
| Amazon RDS for PostgreSQL | AWS | managed instance, parameter groups, Multi-AZ option |
| Amazon Aurora PostgreSQL | AWS | compatible engine, vendor storage and HA |
| Cloud SQL for PostgreSQL | Google Cloud | managed instance, HA option |
| AlloyDB | Google Cloud | PostgreSQL-compatible, vendor HA and features |
| Crunchy | Crunchy Data | PostgreSQL experts; Crunchy Bridge is a managed cloud |
| Supabase | Supabase | PostgreSQL plus an app platform (auth, APIs) |

Aurora and AlloyDB are compatible products. They are not a drop-in promise for every extension. Test `pg_dump` from the service if you must leave.

When PostgreSQL behaves in a way that the docs forbid, you can report it. The official guide is [https://www.postgresql.org/docs/current/bug-report.html](https://www.postgresql.org/docs/current/bug-report.html).

A useful report includes `SELECT version();`, exact steps, the result that you got, the result that you expected with a docs link, and a small schema. Reproduce on community PostgreSQL 16 or 17 when you can. Do not attach secrets or a customer dump.

Do not learn production behavior from a random post when the official page exists. Do not skip `pg_dump` tests because snapshots exist. Do not demand a fix date on a bug report.

### Questions

#### Theoretical questions

1. What does `docs/current` point to?
2. Who wins when the wiki and the official docs disagree?
3. What work do you still own on a managed service?
4. Why are Aurora and AlloyDB not the same object as community PostgreSQL?
5. Why must a bug report include `version()`?

#### Easy practical tasks

1. Open the 16 tutorial and the 17 tutorial. Write one heading that both share.
2. Open the product page for two names in the table. Write the PostgreSQL major versions that each page lists today.
3. Open the 17 bug-report page. Write the channel that the page names.
4. Run `SELECT version();` and save the full string.

#### Medium practical tasks

1. Open the 17 migration notes. Write two items.
2. Compare RDS PostgreSQL and one other product on failover and extensions (vendor docs). Write a six-sentence note.
3. Write a five-line bug-report template: version, steps, got, expected, docs URL. Do not send it.

#### Advanced practical tasks

1. Write an exit plan: `pg_dump` or logical replication from a managed service to a community 16 or 17 cluster. Name the gaps you must test (roles, extensions, sequences).
2. Write a fake but complete report for a made-up `EXPLAIN` defect. Include `CREATE TABLE` and three `INSERT`s. Redact how you would handle production SQL.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do least privilege, a closed `5432`, and versioned migrations form one production boundary?
2. Why is expand-contract a safer default than a weekend rewrite of a large table?
3. How do connection, lag, bloat, and wraparound alerts map to topics 7, 11, and 13?
4. When do you open window, FDW, or PostGIS docs instead of adding a new table?
5. A teammate wants a public RDS instance, app login as the cloud admin, and wiki tuning snippets only. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. Write a one-page production checklist: roles, network, migrations, four monitors, docs URL for your major version.
2. `\du` the app role. `SHOW listen_addresses;`. Run the wraparound `age(datfrozenxid)` query.
3. Bookmark docs for 16, 17, `current`, and the bug-report guide.
4. Compute `row_number() OVER (ORDER BY id)` on one lab table.

#### Medium practical tasks

1. Apply one expand-contract change (add a column) with a fake history table. Enable RLS on a tenant column. Confirm `shop_app` is not superuser.
2. Run the four monitoring queries. Write warn thresholds that fit your lab.
3. Open one cloud HA page and the official tutorial for your version. Write what the vendor hides and what you still test (`pg_dump`).

#### Advanced practical tasks

1. Draft a runbook: disk full, failover (topic 11), wraparound, stolen app password. Keep steps high-level. Name who may use a superuser.
2. Map this whole path (topics 1–14) to official book titles (Tutorial, SQL, Administration, Reference). Add one Internals page that you will read later.
