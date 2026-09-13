# 21. Production Hardening

## Description

This topic shows how you run PostgreSQL 16 and PostgreSQL 17 in production. You learn a least-privilege application role, network rules for port `5432`, migration tools, expand-contract schema changes, monitoring, and runbooks for failover and disk full.

Complete topic 20 first. You know slots, replicas, and heavy SQL. This topic is the last operations topic before the ecosystem.

Use one term for each concept. Least privilege means the role has only the rights that the program needs. A migration is a versioned schema change. Expand-contract is a compatible change in small steps. A runbook is a written procedure that an on-call person follows. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## Least-privilege app role

Topic 12 introduced roles and `GRANT`. Production makes that the default, not an extra.

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
-- table grants and default privileges: topic 12
```

Row-level security stays on if you are multi-tenant. The app role must not bypass RLS.

Cloud admin roles are not app roles. A person can use the admin role from a bastion. The connection pool uses `shop_app` only.

Do not debug a production outage by granting `SUPERUSER` to `shop_app`. Do not reuse the same password for migrator and app. Do not leave `PUBLIC` with `CREATE` on `public` if you use that schema (PostgreSQL 15+ already revoked it; do not grant it back).

### Questions

#### Theoretical questions

1. What privileges must the application role lack?
2. Who owns tables in the migrator-and-app split?
3. Why does the app pool never use the cloud admin role?
4. How does RLS change the app role requirements?
5. Why are default privileges part of least privilege?

#### Easy practical tasks

1. `\du` your lab `shop_app`. Confirm `Superuser` is no.
2. Write a grant list for `shop_app` on one schema with three tables.
3. Connect as `shop_app` and try `CREATE TABLE`. Record the error.
4. Write four sentences: app role, migrator, admin, `PUBLIC`.

#### Medium practical tasks

1. Create `shop_migrator` and `shop_app`. Apply defaults (topic 12). Create a table as migrator. Prove the app can `INSERT` and cannot `DROP TABLE`.
2. Revoke `EXECUTE` on a `SECURITY DEFINER` helper from `PUBLIC`. Grant it only to `shop_app`.
3. Draft a secret plan: where the app password lives, where the migrator password lives, who rotates them.

#### Advanced practical tasks

1. Script a production bootstrap: database, schema, roles, defaults, RLS on one tenant table, no superuser grants. Run it on an empty lab database.
2. Read role attributes in the 17 docs. Add `CONNECTION LIMIT` on `shop_app`. Write why you set it.

---

## Network: do not expose `5432` to the world

PostgreSQL listens on TCP port `5432` by default. If that port is reachable on the public internet, attackers scan it. SSL does not make a public `5432` a good idea (topic 12). A strong password is not enough.

Safe pattern:

- `listen_addresses` binds only the interfaces that need it
- a firewall or security group allows `5432` only from the app subnet, PgBouncer, or a bastion
- `pg_hba.conf` matches those sources and requires `scram-sha-256` (and `hostssl` when you use TLS)
- clients use `sslmode=verify-full` when you control certificates
- administrators use a VPN, a bastion, or a vendor private link, not an open `0.0.0.0/0` rule

Topic 2 said not to combine `listen_addresses = '*'` with `trust`. Production also forbids `trust` and `password` (clear text). Use `scram-sha-256`.

Docker `-p 5432:5432` on a laptop is a lab habit. On a cloud VM, a published `5432` to `0.0.0.0` is an incident.

Managed services still need a network rule. "Public IP disabled" or "authorized networks" is the same idea.

Do not put PostgreSQL on the public internet "for a weekend test" on a host that holds real data. Do not use `trust` for any remote line. Do not disable the firewall because a driver failed to connect; fix `pg_hba` and the security group.

Official `pg_hba`: [https://www.postgresql.org/docs/17/auth-pg-hba-conf.html](https://www.postgresql.org/docs/17/auth-pg-hba-conf.html).

### Questions

#### Theoretical questions

1. Why is a public `5432` a problem even with SSL?
2. Who may reach `5432` in the safe pattern?
3. Which `pg_hba` method do you use for remote passwords?
4. Why is Docker publish to `0.0.0.0` dangerous on a cloud VM?
5. What client `sslmode` matches a controlled CA?

#### Easy practical tasks

1. `SHOW listen_addresses;` `SHOW port;`.
2. Read your cloud or local firewall rules for `5432`. Write the source CIDR.
3. Open `pg_hba.conf` or the vendor equivalent. Write the remote lines (mask secrets).
4. Write four sentences: bind, firewall, `pg_hba`, bastion.

#### Medium practical tasks

1. Draw the path: laptop admin → VPN → bastion → PostgreSQL. Label ports.
2. Compare `host`, `hostssl`, and `hostnossl` for the app subnet. Write the line you would deploy.
3. Document how PgBouncer changes the firewall (app to pooler, pooler to database).

#### Advanced practical tasks

1. Write a network checklist for a new environment: bind, SG/NACL, `pg_hba`, TLS, no public IP.
2. Read `listen_addresses` in the 17 docs. Write how a Unix socket avoids TCP for local-only apps.

---

## Migrations (Flyway, golang-migrate, Sqitch)

A migration is a versioned SQL (or code) change that you apply in order. Production schema must not be "whatever someone ran in `psql`".

Common tools:

- **Flyway** — version files such as `V3__add_orders.sql`; tracks `flyway_schema_history`
- **golang-migrate** — `0003_add_orders.up.sql` and a down file; common in Go services
- **Sqitch** — plans with names and dependencies; strong for review

All three need:

- one history table
- a migrator role (previous section)
- a CI step that applies pending files
- a rule that a file that reached production does not change

```text
-- 0004_items_qty.up.sql
ALTER TABLE shop.items ADD COLUMN qty int NOT NULL DEFAULT 0;
```

Down migrations are optional in some teams. If you keep them, test them. A down that drops a column destroys data.

Do not mix Flyway and golang-migrate on the same database. Do not edit `V3__add_orders.sql` after production applied it. Do not run DDL as `shop_app`.

Topic 13 extensions (`CREATE EXTENSION`) belong in a migration that a superuser or a permitted admin runs, not in the app startup on every request.

### Questions

#### Theoretical questions

1. What does a migration tool store in the database?
2. Why must an applied file stay unchanged?
3. Which role runs migrations?
4. What is the risk of a down migration that drops a column?
5. Why is ad-hoc `psql` DDL a problem in production?

#### Easy practical tasks

1. Open the Flyway, golang-migrate, and Sqitch sites. Write the history-table name or equivalent for each (from their docs).
2. Write one `up` file that adds a nullable column.
3. List three files in a fake `migrations/` folder in order.
4. Write four sentences: version, history, migrator, immutable file.

#### Medium practical tasks

1. Apply two migrations with one tool on a lab database (or simulate with a `schema_migrations` table and two SQL files).
2. Change an already-"applied" file on purpose. Write how the tool reacts, or how you would detect drift.
3. Add `CREATE EXTENSION pg_trgm` to a migration. Write who may run it on your host.

#### Advanced practical tasks

1. Design a CI pipeline: lint SQL, apply to a throwaway 16 or 17 database, run app tests, then apply to production with the migrator secret.
2. Compare Flyway and golang-migrate for a Go team. Write three criteria (repeatability, down files, team skill).

---

## Blue/green or expand-contract schema changes

A lock that rewrites a large table can stop the application. PostgreSQL 16 and 17 still rewrite for some `ALTER TABLE` forms. Check the current `ALTER TABLE` page: some changes are a catalog-only update; some scan the table; some rewrite.

**Expand-contract** (expand, migrate, contract):

1. Expand: add a new column or table that the old app ignores (`NULL` or a default that does not rewrite if the docs say so)
2. Dual write or backfill: the app or a job fills the new object
3. Switch reads to the new object
4. Contract: drop the old column or table in a later release

Example: rename `name` to `title` without a long rewrite:

1. Add `title text`
2. Deploy app that writes both
3. Backfill `title` from `name`
4. Deploy app that reads `title`
5. Drop `name`

**Blue/green** for schema can mean two schemas or two databases. You migrate the green copy, switch the pool to green, keep blue for rollback. This is heavier. It fits a major rewrite. You still need data sync.

`CREATE INDEX CONCURRENTLY` (topic 8) is an expand step. `DROP INDEX CONCURRENTLY` is a contract step. `DETACH PARTITION` (topic 15) can drop old data without a huge `DELETE`.

Do not run `ALTER TABLE ... ALTER COLUMN type` on a 200 GB table in business hours without a rewrite plan. Do not add `NOT NULL` without a default or a prior backfill. Do not deploy an app that requires a column that the migration has not applied.

Official `ALTER TABLE`: [https://www.postgresql.org/docs/17/sql-altertable.html](https://www.postgresql.org/docs/17/sql-altertable.html).

### Questions

#### Theoretical questions

1. What are the four expand-contract steps?
2. Why can a rename of a column need more than one release?
3. When is blue/green a better fit than expand-contract?
4. Which index command belongs in the expand step?
5. Why must you read `ALTER TABLE` for each change type?

#### Easy practical tasks

1. Write expand-contract steps for "add `email_verified boolean`".
2. Write expand-contract steps for "rename `name` to `title`".
3. Open the 17 `ALTER TABLE` page. Find one change that notes a table rewrite.
4. Write four sentences: expand, backfill, switch, contract.

#### Medium practical tasks

1. On a lab table with 10000 rows, add a nullable column. Then add `NOT NULL` after a backfill. Time each step.
2. Plan `CREATE INDEX CONCURRENTLY` plus an app deploy that uses the new filter. Write the order.
3. Compare a weekend blue/green cutover with a four-release expand-contract. Write downtime and complexity.

#### Advanced practical tasks

1. Plan a `varchar(20)` to `text` change on a large table. Use the `ALTER TABLE` notes. State whether you rewrite or add a column.
2. Write a playbook that includes rollback at each expand-contract step (what you still can undo).

---

## Monitoring: connections, replication lag, bloat, wraparound

Production needs measurements, not feelings. Topic 18 showed views. This section is the minimum set.

**Connections**

```sql
SELECT count(*) FROM pg_stat_activity;
SHOW max_connections;
SELECT usename, state, count(*)
FROM pg_stat_activity
GROUP BY 1, 2;
```

Alert when used connections stay high, when many sessions are `idle in transaction`, or when the app cannot connect.

**Replication lag** (topic 16 and topic 19)

```sql
SELECT client_addr, state, replay_lsn, sent_lsn
FROM pg_stat_replication;
```

On the replica, use a lag time function if you have replay. Alert when lag exceeds the report SLA.

**Bloat** (topic 11)

```sql
SELECT relname, n_live_tup, n_dead_tup, last_autovacuum
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC
LIMIT 20;
```

Alert on dead-tuple ratio and on size growth that row counts cannot explain.

**Wraparound** (topic 11)

```sql
SELECT datname, age(datfrozenxid)
FROM pg_database
ORDER BY 2 DESC;
```

Alert far before the wraparound limit. Autovacuum must keep up. A failed freeze is an emergency.

Also watch: disk free (topic 18), replication slot retained WAL, `pg_stat_statements` outliers, checkpoint spikes, error logs.

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
2. Join `pg_replication_slots` with a WAL-bytes estimate (topic 18 query). Write a slot alert.
3. Add `application_name` to the connection group-by. Write which apps hold idle transactions.

#### Advanced practical tasks

1. Write a single monitoring SQL script that prints all four areas plus disk-related sizes (`pg_database_size`, slot lag).
2. Map each alert to a first action (pool, replica, vacuum, wraparound runbook). Keep actions high-level.

---

## Runbooks for failover and disk full

A runbook is a procedure that a tired person can follow. Topic 16 covered promote and HA tools. Topic 18 covered a full disk. Production stores the runbook next to the pager.

**Failover runbook** (physical replica, high-level):

1. Confirm the primary is down or unsafe (more than one check).
2. Fence the old primary (stop writes: stop the VM, revoke security group, vendor fence). Do not skip this step.
3. Promote one replica (`pg_promote` or Patroni or vendor failover).
4. Point the pool or DNS at the new primary.
5. Verify writes (`INSERT` in a check table) and replica lag if a new stand-by exists.
6. Open an incident: data-loss window, who promoted, when you restore a new replica.

**Disk-full runbook** (high-level):

1. Confirm which volume is full (`PGDATA`, `pg_wal`, logs).
2. If the server accepts a connection, run slot and archiver queries (topic 18).
3. Free space without deleting `pg_wal` at random.
4. Fix the cause (archive, slot, logs, bloat).
5. Confirm writes resume and WAL recycles.
6. If files were destroyed, restore from backup (topic 16) and stop guessing.

Each runbook needs: owner, last test date, links to dashboards, vendor console steps, and a rollback or "stop and escalate" line.

Do not keep the only copy of the runbook on the primary host. Do not promote two replicas. Do not practice failover only on paper.

### Questions

#### Theoretical questions

1. What is a runbook?
2. Why must you fence the old primary before you promote?
3. What is the last verify step after failover in this section?
4. What must you not delete when the WAL disk is full?
5. Why do you store the runbook off the primary host?

#### Easy practical tasks

1. Write the six failover steps in your own words.
2. Write the six disk-full steps in your own words.
3. Add owner and last test date fields to both.
4. Write four sentences: fence, promote, pool, backup.

#### Medium practical tasks

1. Fill a one-page failover runbook for your lab (Docker Compose or vendor). Include exact commands that you already used in topic 16.
2. Fill a one-page disk-full runbook that references your log path and archive folder.
3. Schedule a drill date. Write what success looks like (time to write, data-loss check).

#### Advanced practical tasks

1. Combine Patroni or one cloud HA page with your runbook. Write which steps the tool does and which steps a person still does.
2. After a lab promote, write the postmortem template: timeline, lag, extra replica, backup still valid.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do least privilege, network rules, and a migrator role form one production boundary?
2. Why must expand-contract and the migration tool agree on order?
3. How do wraparound alerts and a disk-full runbook both protect WAL and freeze work?
4. Why is a replica without a failover runbook not high availability?
5. A teammate wants a public `5432`, app superuser, and `ALTER COLUMN` on a 200 GB table on Friday afternoon. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. Connect as a non-superuser app role. Run `SELECT current_user, rolsuper FROM pg_roles WHERE rolname = current_user;`.
2. Write a cheat sheet: app versus migrator, `5432` rules, three migration tools, expand-contract, four monitors, two runbooks.
3. Run connection, lag (or empty `pg_stat_replication`), bloat, and wraparound queries once.
4. Write one expand-contract sequence for adding a `phone text` column.

#### Medium practical tasks

1. Produce a production checklist with 12 boxes that cover all six sections. Check what your lab already meets.
2. Apply a two-step expand (nullable column, then backfill) with a migration file. Prove `shop_app` never ran DDL.
3. Write warn thresholds for the four monitors and the first link each alert opens (dashboard or runbook).

#### Advanced practical tasks

1. Run a failover drill or a disk-pressure drill on a lab. Fill the runbook with actual times and commands. Fix the runbook where it was wrong.
2. Map this topic to official chapters: roles, `pg_hba`, `ALTER TABLE`, monitoring, high availability. Add one production warning from the docs that this topic did not include.
