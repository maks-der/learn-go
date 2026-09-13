# 13. Performance and Security

## Description

This topic shows how you find a slow database and how you protect who connects and what they can do. You learn measurement, cardinality, hot rows, pagination, caching, authentication, authorization, secrets, TLS, encryption at rest, injection, and exposed admin ports.

Use one term for each concept. A metric is a number that you collect over time. Authentication proves an identity. Authorization grants rights to that identity. Complete this topic after you can read a query plan and connect with a driver.

Do not change a setting before you have a baseline. Measure, change one thing, measure again. Security is a property of the whole path, not of one password field.

---

## Measure first: slow query log and metrics

A slow query log records statements that exceed a time threshold. You set the threshold. The log gives the SQL text, the duration, and sometimes a plan. Start here when users say "it is slow" and you do not know which statement.

Metrics are time series: connections, CPU, disk I/O, cache hit ratio, lock waits, replication lag, transaction rate. Metrics show the shape of the problem (CPU bound, I/O bound, lock bound, connection bound). They do not always name the SQL.

```text
User report --> metrics (which resource) --> slow log or traces (which SQL) --> plan --> change
```

A baseline is a measurement before the change. Record: query text, duration, plan, row counts, and the time of day. Without a baseline you cannot know if the change helped.

Typical first questions:

1. Is the host saturated (CPU, disk, memory)?
2. Is one statement slow, or are many statements waiting?
3. Is the wait on I/O, on a lock, or on a sequential scan?

Do not enable a very low slow-log threshold on a busy production system without a plan. The log can fill the disk. Use a sample, a higher threshold, or a product-specific statement store.

Do not tune by folklore ("increase shared buffers and it will be fast"). Correlate a metric with a statement.

Application metrics matter. A 200 ms API time with 5 ms in the DBMS is not a database problem. Measure both sides.

### Questions

#### Theoretical questions

1. What does a slow query log record?
2. What do metrics show that a single log line may not show?
3. What is a baseline?
4. Why can a very low slow-log threshold harm a busy system?
5. Why must you measure the application and the DBMS?

#### Easy practical tasks

1. Find the slow-query log setting and the threshold name in your DBMS.
2. Write the four-step path from this section in your notes.
3. List five metrics that you would graph for a small server.
4. Write three baseline fields that you store before an index change.

#### Medium practical tasks

1. Turn on a slow log in a learning instance. Run a sequential scan on a large scratch table. Find the line in the log.
2. Record CPU and disk busy (OS or DBMS) during that scan. Write which resource rose.
3. Compare one API timer with the DBMS duration for the same request (log both). Write the two numbers.

#### Advanced practical tasks

1. Write a one-page measure playbook: logs, metrics, traces, and when you stop at the application.
2. Build a weekly review: top five statements by total time. Keep the list. Do not change production in this task unless you own it.

---

## Cardinality, hot rows, and pagination

Cardinality is the number of distinct values in a column (or in a group of columns). High cardinality means many distinct values. Example: a unique email column has cardinality equal to the number of rows. A boolean column has cardinality 2.

The optimizer uses statistics to estimate how many rows a filter returns. If statistics are stale, the estimate is wrong. A wrong estimate can pick a nested loop that should be a hash join, or a scan that should be a seek. After a large load, refresh statistics (product command: `ANALYZE` or equivalent).

Do not create an index on a boolean that is true for 99 percent of rows and expect a seek to help those queries. Low selectivity favors a scan.

Do not confuse cardinality with table size. A table of 10 million rows can have a column of cardinality 10.

A hot row is a row that many transactions try to change at the same time. Example: a single counter row for "tickets remaining," or one parent order that every line update locks.

Lock contention is wait time on locks. Sessions queue. Throughput drops. CPU can look idle while sessions wait.

```text
Session A: UPDATE counters SET n = n - 1 WHERE id = 1;  -- holds row lock
Session B: same row  -->  wait
```

Isolation and row locks are correct. The design is wrong if one row is a global bottleneck.

Mitigations (choose after you measure):

1. Split the counter into many rows and sum them.
2. Use an atomic increment that the product documents, if it reduces transaction time.
3. Move the hot fact out of the transactional table (queue, cache with a later reconcile) when the business allows a short delay.
4. Shorten the transaction. Do not hold the hot row lock while you call a network service.

Do not raise isolation to serializable as a first fix for a hot row. Do not add more application servers that all update the same row. You add more waiters.

Pagination returns a page of rows. Two common methods exist.

**OFFSET pagination.** You request `LIMIT n OFFSET k`. The DBMS still walks or sorts k + n rows, then discards k. Large k is expensive.

```sql
SELECT order_id, ordered_at
FROM orders
ORDER BY ordered_at DESC, order_id DESC
LIMIT 20 OFFSET 40;
```

**Keyset pagination** (seek method). You ask for rows after the last key of the previous page. The query uses a range on the sort key. Cost stays near one page.

```sql
SELECT order_id, ordered_at
FROM orders
WHERE (ordered_at, order_id) < (:last_at, :last_id)
ORDER BY ordered_at DESC, order_id DESC
LIMIT 20;
```

Keyset needs a unique sort key (add the primary key to break ties). OFFSET can jump to page 500 (at a cost). OFFSET also drifts: a new insert at the top can shift rows between pages.

Do not use `OFFSET` for deep pages on large tables. Do not use `SELECT *` on a wide table for a list page. Total count (`SELECT COUNT(*)`) is a separate cost. Do not run a full count on every page if the UI can show "more" without a total.

An index that matches `ORDER BY` (and the keyset predicate) keeps the page cheap. Read the plan.

### Questions

#### Theoretical questions

1. What is cardinality in this section, and what can stale statistics do to a plan?
2. What is a hot row, and why do more application servers not fix one-row contention?
3. Why is a large `OFFSET` expensive, and what does a keyset page use instead?
4. Why must the keyset sort key be unique, and how can `OFFSET` pages drift when inserts occur?
5. Why must you not hold a row lock during an external HTTP call?

#### Easy practical tasks

1. Estimate cardinality for: primary key, country code, boolean `is_deleted` on a live-only table. Find the command that refreshes statistics in your DBMS.
2. Name two hot-row examples in a shop or a ticket system. Draw three sessions and one locked row.
3. Write an `OFFSET` query for page 3 of size 10. Write a keyset query that continues after a given `(ordered_at, order_id)`.
4. Make a two-column table: OFFSET versus keyset. Add four rows. Write one UI that needs jump-to-page and one that needs only Next.

#### Medium practical tasks

1. Load thousands of rows. Explain a filter. Refresh statistics. Explain again. Write what changed.
2. Open two sessions. Update the same row in an uncommitted transaction in session A. Update in session B. Record the wait. Commit A. Time 100 sequential updates to one counter versus 100 updates spread across 20 counters.
3. Fill a table with thousands of rows. Time `OFFSET 0`, `OFFSET 1000`, and `OFFSET 10000` with `LIMIT 20`. Time the keyset query for the same late page. Write the times and the plan access method.

#### Advanced practical tasks

1. Write a one-page note: when you run `ANALYZE` after loads, how you detect estimate errors, and a ticket-inventory design without one global row.
2. Implement Next/Previous with keyset in a small program. Prove that a new insert does not duplicate a row on Next. Write a pagination standard: default method, when OFFSET is allowed, and how you avoid `COUNT(*)`.

---

## Caching in front of the database

A cache stores query results or objects in a fast store (memory, key-value product). The application reads the cache first. On a miss, it reads the DBMS and fills the cache.

Caching helps when:

- the same read is frequent
- the result can be a few seconds or minutes old
- the DBMS is the bottleneck

```text
Read:  cache hit --> return
       cache miss --> DBMS --> fill cache --> return
Write: DBMS commit --> invalidate or update cache keys
```

Invalidation is the hard part. If you update a row and leave the old cache key, users see stale data. If you invalidate too many keys, the cache is useless and the DBMS sees a storm of misses (thundering herd).

Cache keys must include every input that changes the result (user, page, locale). A key that omits the user can leak data across users. That leak is a security bug.

Do not cache a write path that must stay consistent with constraints unless you still commit in the DBMS. The cache is not the system of record.

Do not add a cache to hide an N+1 query or a missing index. Fix the query first. Measure. Then cache if the read is still hot.

A short TTL (time to live) limits staleness without perfect invalidation. TTL is a trade-off, not a full design.

The DBMS buffer cache is not this cache. Buffer cache holds pages. Application cache holds results.

### Questions

#### Theoretical questions

1. When does an application cache help?
2. What must a write do to cached data?
3. What is a thundering herd after a large invalidation?
4. Why must a cache key include the user when the result is per user?
5. Why must you not use a cache as the system of record?

#### Easy practical tasks

1. Draw hit, miss, and invalidate for a product-price cache.
2. Write a cache key for "orders for customer 42, page after id 100."
3. List three reads that may use a 30-second TTL and two that must not.
4. Write one sentence that distinguishes application cache from buffer cache.

#### Medium practical tasks

1. Implement a tiny in-memory cache for `GET product` in a learning app. Invalidate on update. Show stale data if you skip invalidate.
2. Measure DBMS query count with and without the cache for 50 identical reads.
3. Write a stampede rule: single fill on miss (lock or request coalesce). Describe it in six sentences.

#### Advanced practical tasks

1. Write a one-page cache policy: keys, TTL, invalidation, personal data, and what you never cache.
2. Design cache keys for a join-heavy report. List every invalidation event. Mark keys that you cannot maintain (then do not cache).

---

## Authentication vs authorization and least privilege

Authentication answers "who is this session?" The DBMS checks a user name and a credential (password, certificate, token, or an external directory). A failed authentication does not start a session.

Authorization answers "what may this session do?" The DBMS checks privileges: connect, read a table, write a table, create objects, manage roles.

```text
Connect request --> authenticate --> session identity
SQL statement   --> authorize    --> allow or reject
```

A correct password with no `SELECT` privilege still fails `SELECT`. A `GRANT` without a successful login never applies.

Typical identities: a human admin, an application user, a read-only reporter, a migration user (schema change).

Do not share one user for all of these jobs. You cannot revoke a leaked application password without also locking the admin. You cannot audit who changed a row if every program uses `postgres` or `sa`.

Least privilege means a role has only the rights that it needs. The application user can `SELECT`, `INSERT`, `UPDATE`, and `DELETE` on the tables that the app uses. It cannot `DROP TABLE` on production.

```sql
-- Shape only; verbs and names differ by product
GRANT SELECT, INSERT, UPDATE, DELETE ON orders TO app_shop;
REVOKE DELETE ON invoices FROM app_shop;
```

Start from zero. Grant what the program needs. Do not grant `SUPERUSER`, `DBA`, or `sysadmin` to an application.

Migrations need a stronger user. Run migrations as that user. The application does not use that user at run time.

Read-only roles must not write. A report replica still needs a user that cannot `DELETE`.

Object owners and `PUBLIC` defaults matter. Some products grant wide rights to `PUBLIC`. Revoke those rights on sensitive objects.

Do not grant `SELECT` on all tables if the app only needs five tables. Review grants when you add a table. A new table can inherit defaults that are too wide. Grant explicitly.

`trust` authentication (any connector is accepted) is only for a locked-down lab. Log authentication failures. Many failures from one host can be an attack or a bad connection string.

### Questions

#### Theoretical questions

1. What question does authentication answer, and what question does authorization answer?
2. Why does a valid password not imply `SELECT` success?
3. What does least privilege require, and why must the application not run as a superuser?
4. Why must an application and an admin use different users, and why do migrations use a different user than run time?
5. What risk does a `PUBLIC` default grant create?

#### Easy practical tasks

1. Write the two-step path from this section. List four identities for a class shop. Assign one job each.
2. Write a privilege list for `app_shop` on `products` and `orders`. Write three rights that `app_shop` must not have in production.
3. Find how you create a user (role) in your DBMS. Find `GRANT` and `REVOKE` in your DBMS manual.
4. Write one authentication failure and one authorization failure as error stories.

#### Medium practical tasks

1. Create a user that can connect but cannot `SELECT` a table. Prove both: connect works, `SELECT` fails.
2. Create `app_shop` and `report_shop`. Grant write versus `SELECT` only. Prove with two sessions. List current grants on a practice table. Write what `PUBLIC` has.
3. Add a new table. Check whether `app_shop` can read it. Grant only if you intend that. Enable or locate auth-failure logs. Fail a login on purpose. Find the log line.

#### Advanced practical tasks

1. Write a one-page identity and grant standard: human, app, report, migrate, default revoke, new-table checklist.
2. Review a sample database ACL. List every grant that is wider than the app needs. Write the `REVOKE` list (do not apply it on a shared server).

---

## Secrets, TLS, and encryption at rest

A connection string often contains a password. That password is a secret. A secret in a Git repository, a chat log, or a container image is a leak.

Rules:

1. Store secrets in an environment variable, a secret manager, or a file with tight permissions. Do not store them in source.
2. Use a different password in development, test, and production.
3. Rotate a leaked password. Rotate on a schedule if the policy requires it.
4. Do not print the connection string in logs. Drivers can log the URI. Turn that off or redact.
5. Prefer a URI without a password plus a separate password field that the process reads at start.

```text
# Environment (example names)
DATABASE_HOST=127.0.0.1
DATABASE_USER=app_shop
DATABASE_PASSWORD=...   # from a secret store, not from the repo
```

`.env` files help local work. Do not commit `.env`. A committed `.env.example` contains placeholders only.

If a secret leaked, change it. A leak in an old commit stays in Git history. Rotate. Consider the repository as public from that moment.

TLS encrypts the network bytes between the client and the DBMS. An observer on the network cannot read the SQL text or the rows. TLS also authenticates the server when the client verifies the certificate.

Without TLS on an untrusted network, a packet capture can show passwords and personal data. "The port is on a private VLAN" reduces risk. It does not encrypt.

Modes (names differ): disable, allow/prefer, require, verify-ca / verify-full. Use verify-full (or the product equivalent) in production. A `require` mode without name checks can still accept a different server.

Do not set "skip verify" in production to silence an error. TLS does not replace authentication. TLS does not replace grants. TLS does not encrypt data files on disk. Do not expose the DBMS port on the public internet even with TLS. TLS is not a firewall.

Encryption at rest protects data files and backups when someone obtains the disk or the backup file. The bytes on the volume are ciphertext without the key.

Layers (you can use more than one):

| Layer | What it encrypts | Who holds the key |
| --- | --- | --- |
| Full-disk or volume | The whole volume | OS or cloud KMS |
| Tablespace or TDE | Database files | DBMS plus KMS |
| Backup encryption | Dump or snapshot files | Backup tool plus KMS |
| Application field encryption | Selected columns | Application |

Volume encryption helps when the disk is detached or stolen. If an attacker gets a running server and the mounted volume, the OS already decrypted the disk.

Transparent data encryption (TDE) or equivalent encrypts files that the DBMS writes. An attacker with only the files still needs the key.

Column encryption in the application hides values from the DBMS and from many admins. You lose ordinary `WHERE` and indexes on the cleartext unless you design searchable encryption.

Keys are secrets. A key on the same disk as the data is a weak design. Use a key management service.

Encryption at rest does not hide data from a connected application user. Grants and authentication still apply. Encrypt the backups too. Store the key away from the backup files.

### Questions

#### Theoretical questions

1. Why is a password in a Git repository a leak, and where may a process read a production password?
2. What does TLS protect on the wire, and what extra check does verify-full add over "TLS required"?
3. Why is a private VLAN not a substitute for TLS, and what does TLS not replace?
4. What threat does encryption at rest address, and why does volume encryption not stop an attacker on a running mounted server?
5. Why must the key not live only beside the data files, and why must you rotate after a leak even if you delete one file?

#### Easy practical tasks

1. Write a connection URI with `PASSWORD` as a placeholder. Add `.env` to a ignore file if you use one. Write an `.env.example` with empty values.
2. Find the client TLS mode names for your driver. Write the production mode that you choose (verify-full or the product name).
3. Fill the four-layer at-rest table in your notes. Write one threat that TLS covers and one threat that at-rest covers.
4. List five places a teammate might paste a URI (repo, chat, ticket, image, log). Mark each as forbidden for real secrets.

#### Medium practical tasks

1. Change a small program to read the password from the environment. Prove that the source file has no password. Search your practice repo for `password=`. Record hits.
2. Connect with TLS disabled and with TLS required on a learning instance if you can. Write the connection option. Write a rotate procedure of six steps.
3. Compare a cleartext dump and an encrypted backup option in the docs. Write the restore extra step (the key). Pick one column (national id). Write why you would encrypt it in the app or why grants plus TLS are enough for a class project.

#### Advanced practical tasks

1. Write a one-page secret and TLS standard: stores, rotation, CI, logs, modes, CA, and a ban on skip-verify.
2. Write a one-page key policy: where keys live, who rotates, what you do if a backup tape is lost. Enable TLS on a disposable DBMS. Connect with verify-full. Record the failed connect when the name is wrong.

---

## Injection and exposed admin ports

Three common failures appear together in incidents.

**Injection.** The application builds SQL from untrusted text. The attacker runs SQL as the application user. Parameterized statements prevent value injection. Least privilege limits the damage if injection still occurs. An app user without `DROP` cannot drop tables through injection.

**Overly broad grants.** `GRANT ALL` or a superuser application turns any bug into full control. A leaked password becomes a full dump. Review grants. Revoke unused rights.

**Exposed admin ports.** The DBMS listen address is `0.0.0.0` on a public IP. Scanners find the port. They try passwords. They try old bugs. Bind to localhost or to a private network. Use a firewall. Use a VPN or a bastion for admin tools. Do not publish `5432`, `3306`, or vendor admin UI ports to the internet.

```text
Public internet --x--> DBMS port
Admin --> VPN / bastion --> DBMS
App   --> private net   --> DBMS
```

Also close unused extensions and unused languages if the product allows it. Remove sample users and default passwords.

Defense in layers:

1. No public DBMS port
2. TLS with verification
3. Strong unique passwords or certificates
4. Least privilege
5. Parameterized SQL
6. Encryption at rest for disks and backups

Do not rely on one layer. Do not test injection on systems that you do not own.

### Questions

#### Theoretical questions

1. How does least privilege reduce the effect of SQL injection?
2. Why is `GRANT ALL` to the application a high-severity finding?
3. Why must the DBMS port not listen on a public address?
4. What path should an admin tool use instead of a public port?
5. Why do you need more than one of the six layers?

#### Easy practical tasks

1. Write the six layers as a checklist.
2. Find the listen address setting in your DBMS. Write the learning value (`localhost`) and a forbidden value for a laptop on café Wi-Fi.
3. List default ports for two DBMS products. Mark them as "do not publish."
4. Write one `REVOKE` that you would apply after a `GRANT ALL` mistake (on a disposable database).

#### Medium practical tasks

1. Confirm that your learning DBMS is not reachable from another machine. If it is, bind to localhost or add a firewall rule that you understand.
2. Combine a parameterized query with a user that cannot `DROP`. Write how an injection attempt would fail to drop a table.
3. Search a public advisory about an exposed database port. Write five sentences in your own words. Do not copy the advisory.

#### Advanced practical tasks

1. Write a one-page hardening guide for a small production shop: bind, firewall, TLS, roles, secrets, backup encryption.
2. Review a compose file or a cloud security group. List every open port. Close or justify each. Record the change in a disposable environment.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a slow log, statistics, and a plan work together before you add an index or a cache?
2. When is a hot row the cause, and when is a deep `OFFSET` the cause, of a slow page?
3. How do authentication, authorization, and TLS fail independently on the same connection string?
4. How do encryption in transit and encryption at rest cover different attackers, and which controls stop injection versus only limit the blast radius?
5. A teammate raises buffers, adds a cache, puts the admin URI in the app, and opens port 5432 on `0.0.0.0` without a baseline. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: measure, cardinality, hot rows, keyset, cache, authn versus authz, secrets, TLS, at rest, six layers.
2. Enable or locate a slow log. Copy one line into your notes (sanitize secrets). Create a non-admin user. Connect. Run `SELECT 1`. Show that `DROP DATABASE` fails.
3. Write one keyset query and one OFFSET query for the same list. Redact a sample URI for a README. Keep host and database name. Remove the password.
4. Write the listen address and TLS mode of your learning instance.

#### Medium practical tasks

1. Build a scratch workload: large table, stale stats, then `ANALYZE`, then a keyset page. Record plans and times. Produce a lock wait with two sessions. Write how the metrics differ.
2. Build a mini app user: env-based password, grants on two tables only, TLS if available, parameterized SQL. Write a threat table: stolen laptop disk, packet capture, leaked Git password, SQL injection. Map one control each.
3. Write a tuning ticket template: baseline, hypothesis, one change, result. Run a privilege review on your `learn` database. Produce a grant list and a revoke plan.

#### Advanced practical tasks

1. Tune one real slow statement (or a made-slow one): measure, stats, index or rewrite, pagination fix. Write before/after. Write a security review of a small project against this topic. File findings by severity. Fix the ones you own.
2. Design production controls for a payments schema: roles, secrets, TLS verify-full, encrypted backups, no public port, cache policy, connection budget. Draw the network.
