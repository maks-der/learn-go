# 12. Administration and Security

## Description

This topic shows how you run and protect a MongoDB deployment. You configure `mongod`. You learn the WiredTiger storage engine. You take backups, you practice restore, and you plan upgrades. You authenticate users and you assign roles. You limit network access, you enable TLS, and you use typed driver APIs. You learn field-level encryption at a high level.

Complete drivers and replica-set ideas first. Do not run administration commands on a shared cluster without a role and a change window.

Use one term for each concept. **`mongod`** is the database process. **WiredTiger** is the default storage engine. **Authentication** proves who the client is. **Authorization** (RBAC) decides what that user can do. **TLS** encrypts the network.

A learning database on localhost still needs a password before you expose a port. Atlas needs a database user and an IP access list (or a private network) before you open the cluster to an application.

---

## `mongod` and WiredTiger

`mongod` reads a YAML configuration file. The default path depends on the operating system. You can also pass flags. Prefer a file for anything that must survive a restart.

Common settings (names are examples; read your version):

- `storage.dbPath` — data files
- `systemLog.path` — log file
- `net.port` — listen port (`27017` is the default)
- `net.bindIp` — interfaces
- `replication.replSetName` — replica-set name
- `security.authorization` — user access on or off

Example shape:

```yaml
storage:
  dbPath: /var/lib/mongo
systemLog:
  destination: file
  path: /var/log/mongodb/mongod.log
net:
  port: 27017
  bindIp: 127.0.0.1
```

Start:

```text
mongod --config /etc/mongod.conf
```

`mongosh` is the client. `mongod` is the server. Do not confuse the two processes.

On Windows, MongoDB often runs as a service. On Docker, flags or an env file set the same ideas. Atlas hides `mongod.conf`. You set equivalent ideas in the Atlas UI (cluster tier, IP list, backup).

Validate a change on a test host. A wrong `dbPath` or `bindIp` can start an empty instance or expose the port.

Do not enable a public `bindIp` without authentication and TLS.

**WiredTiger** is the default storage engine. It stores documents and indexes on disk. It uses a cache in RAM. It supports document-level concurrency. Two updates to two documents can run at the same time.

WiredTiger compresses data. Compression reduces disk. Compression uses CPU. The default is fine for most learning systems.

WiredTiger is not MMAPv1. MMAPv1 is gone from current MongoDB. Old blog posts about record padding and collection-level locks do not apply to WiredTiger.

Checkpoints write a consistent snapshot to the data files. Between checkpoints, the **journal** holds recent writes. Journaling is on by default. After a crash, WiredTiger replays the journal and returns to a consistent state. Do not turn the journal off on a system that you care about.

Write concern `j: true` waits until the write is in the journal. The journal is not a user backup.

Do not delete files under `dbPath` to "free space" while `mongod` runs.

`db.serverStatus().wiredTiger` shows cache and other counters. Atlas shows storage charts.

### Questions

#### Theoretical questions

1. What process is `mongod`?
2. What does `storage.dbPath` select?
3. What is the default storage engine?
4. What does document-level concurrency mean?
5. What problem does the journal solve after a crash?

#### Easy practical tasks

1. Find the config file or the Docker command that starts your `mongod`. Write the path or the command.
2. Run `db.serverStatus().storageEngine` or the equivalent. Write the name.
3. Open the configuration-file page and the WiredTiger page. Write both URLs.
4. Write five settings from this section and one sentence each.

#### Medium practical tasks

1. Change the log path on a lab instance. Restart. Confirm new log lines.
2. Read `wiredTiger.cache` in `serverStatus`. Write current and max bytes.
3. Insert with `{ writeConcern: { w: 1, j: true } }` and with `{ w: "majority" }`. Write if both succeed on your set.

#### Advanced practical tasks

1. Document a production config checklist: bind IP, auth, TLS, replica set, paths, log rotation, journal always on.
2. Explain how WiredTiger cache and the OS file cache both hold data. Write six sentences.

---

## Backup, restore, and upgrades

A **backup** is a copy that you can restore. A secondary is not a backup. A delayed member is not a full backup.

Common methods:

**`mongodump` / `mongorestore`.** Logical export of BSON. Good for small data and for selected collections. Slow and heavy on large clusters. You can dump from a hidden secondary.

**Filesystem or volume snapshots.** Snapshot the disk (or the cloud volume) with a consistent method. For self-managed replica sets, the manual describes snapshot plus the journal, or snapshot on a secondary with flush procedures.

**Atlas backup.** Atlas offers cloud snapshots and, on some tiers, more frequent backups and point-in-time restore. Use Atlas backup for Atlas clusters unless you have a written exception.

Pick a method that matches size and RPO (how much data you can lose). A nightly `mongodump` of 2 TB is often the wrong tool.

Encrypt backups. Control who can download them. A backup has all the data.

Do not run `mongodump` against the primary at peak without a plan. Prefer a hidden secondary.

Do not copy `dbPath` with `cp` or zip while `mongod` writes, unless the official procedure allows that exact method.

A **restore drill** is a practice restore. You restore to a safe target. You measure time. You check data. You write what failed.

A drill answers:

- How long to restore (RTO)
- Who has the passwords and the roles
- Whether the backup is complete
- Whether applications can point to the restored data

Steps in spirit:

1. Pick a backup (dump, snapshot, or Atlas snapshot).
2. Restore to a new cluster or a new database name. Do not overwrite production on the first drill.
3. Run checks: document counts, a sample find, an application smoke test.
4. Record the clock and the problems.
5. Repeat on a schedule (for example every 90 days).

`mongorestore` can drop collections if you pass drop flags. Read the flags. A wrong restore can delete data.

WiredTiger can leave free space inside files after deletes. The disk file does not always shrink. **Compact** rewrites data to release space (with version-specific behavior). Compact can take a long time. Do not run compact on the primary at peak. Compact does not fix a bad model.

A **major version** is the first number (for example 7 to 8). MongoDB supports an upgrade path. You do not jump over unsupported gaps. Read **Upgrade** for your from-version and to-version.

Typical ideas:

- Upgrade binaries on secondaries first, then the primary (rolling)
- Set **feature compatibility version** (`setFeatureCompatibilityVersion`) only when the docs say so
- Read compatibility of drivers and of Atlas
- Test on a copy of the data
- Read deprecated commands and removed options

Atlas: you pick a version in the UI. Atlas still needs a change window and an application test.

A downgrade can be hard or impossible after you enable new FCV features. Do not set FCV in production on the same day as the binary upgrade unless the procedure says that you must.

Do not upgrade production on a Friday without a team and a rollback plan that the manual still allows.

### Questions

#### Theoretical questions

1. Why is a secondary not a backup?
2. When is a disk snapshot better than `mongodump`?
3. What is a restore drill?
4. Why do you upgrade secondaries before the primary?
5. What is feature compatibility version?

#### Easy practical tasks

1. Open the backup methods page and the Atlas backup page. Write both URLs.
2. Run `mongodump --help`. Write two options (`--db`, `--out`).
3. Write the major version of your server (`db.version()`).
4. Make a table: method, good for, poor for. Add dump, snapshot, Atlas.

#### Medium practical tasks

1. Dump one small database with `mongodump`. Restore it to a new database name with `mongorestore`. Confirm a document. Time the restore.
2. Read removed features between your version and the next. Write two items that could break an app.
3. Write a test plan: CRUD, index build, transaction, change stream after upgrade.

#### Advanced practical tasks

1. Read consistent snapshot steps for a replica set on your OS or cloud. Write the steps in order. Do not skip journal notes.
2. Perform a major upgrade on a local replica set in Docker. Record commands, FCV, and one failed write if any.

---

## Authentication and role-based access

**Authentication** checks credentials. MongoDB supports several mechanisms. The common ones for applications are **SCRAM** and **x.509**.

**SCRAM** (Salted Challenge Response Authentication Mechanism) uses a user name and a password. Current servers prefer **SCRAM-SHA-256**. The client and the server prove the password without sending the password in the clear (you still need TLS on the network).

Create a user on a self-managed server after you enable authorization:

```javascript
use admin
db.createUser({
  user: "app",
  pwd: "use-a-secret-store",
  roles: [ { role: "readWrite", db: "shop" } ]
})
```

Put the user and the password in the URI, or let the driver prompt. Do not commit the password.

**x.509** uses certificates. The client presents a certificate. The server maps the certificate subject to a user. Replica-set members can also use x.509 to authenticate to each other.

Atlas: you create a database user in the UI. Atlas uses SCRAM for those users by default. You can add other methods that Atlas documents.

Enable authorization on Community Server (`security.authorization: enabled`). Without it, a client that can reach the port can do anything.

`localhost exception` lets you create the first user on a fresh self-managed server. Read that page before you lock yourself out.

Do not create a user with the `root` role for the application.

**Role-based access control (RBAC)** grants privileges through **roles**. A privilege is an action on a resource (a database, a collection, or the cluster).

Built-in roles (examples):

- `read` — find on a database
- `readWrite` — find and write
- `dbAdmin` — indexes and some admin on a database
- `userAdmin` / `userAdminAnyDatabase` — manage users
- `clusterMonitor` — many read-only cluster stats
- `backup` / `restore` — dump and restore related actions
- `root` — all actions (operators only, and only when needed)

Assign the smallest role that works. The application needs `readWrite` on its database, not `root`.

You can create a custom role when built-in roles are too wide:

```javascript
db.createRole({
  role: "orderWriter",
  privileges: [
    { resource: { db: "shop", collection: "orders" }, actions: [ "find", "insert", "update" ] }
  ],
  roles: []
})
```

Atlas has built-in roles in the UI. Custom roles exist on Atlas too. Use them.

Do not share one user across many applications. A leak then has a large blast radius.

Do not grant `dropCollection` to the web application if the application never drops collections.

### Questions

#### Theoretical questions

1. What does authentication prove?
2. What does SCRAM use as credentials?
3. Why must you enable authorization on a self-managed server?
4. Why is `readWrite` on one database better than `root` for an API?
5. When do you create a custom role?

#### Easy practical tasks

1. Open the authentication page and the built-in roles page. Write both URLs.
2. On Atlas, view Database Access. Write the user name (not the password) and the auth method.
3. Write the minimum role for a reporting job that only finds in `shop`.
4. Make a table: user, role, database. Add app, analyst, operator.

#### Medium practical tasks

1. Enable authorization on a local `mongod`. Create a user. Connect with that user. Confirm a bad password fails.
2. Create two users: `app` with `readWrite` on `learn`, and `reader` with `read`. Prove the reader cannot insert.
3. Create a custom role that can `find` and `insert` on `shop.orders` only. Assign it. Prove `shop.items` insert fails.

#### Advanced practical tasks

1. Design roles for app, migrator, backup, and human admin. Write actions that each must not have.
2. Compare Atlas database users with self-managed `admin` users. Write how rotation of a password differs.

---

## Bind IP, TLS, and typed APIs (no string-built queries)

A MongoDB port on the public internet without controls is a common incident. Limit who can open a TCP connection.

**Self-managed `bindIp`.** Bind to `127.0.0.1` for a local-only lab. Bind to a private interface for a VPC. Do not bind to `0.0.0.0` without authentication, TLS, and a firewall.

**Firewall.** Allow port `27017` (or your port) only from application subnets. Deny the rest.

**Atlas IP access list.** Atlas refuses clients that are not on the list (unless you use a private endpoint and the matching mode). Add the application egress IPs. Do not add `0.0.0.0/0` except on a short-lived learning cluster that has no real data, and remove it after.

**VPC and private networking.** Place `mongod` in a private subnet. On Atlas, use VPC peering or a private endpoint so that traffic does not use the public internet.

Network control is not authentication. You need both. A stolen password from a laptop on the allowlist still works. A public port with a strong password is still a scan target.

Do not keep a personal home IP on the Atlas list after you leave the project.

**TLS** (Transport Layer Security) encrypts bytes on the network. Clients and servers authenticate the certificate. Without TLS, a network observer can read queries and documents.

Atlas requires TLS for client connections. The official URI includes `tls=true` (or the older `ssl=true`). The driver verifies the Atlas certificate.

Self-managed:

- Set `net.tls.mode` (for example `requireTLS`)
- Point to a PEM certificate and key
- Distribute a CA file to clients
- Use TLS for replica-set member traffic too

Do not use a self-signed certificate in production without a real CA plan. Do not disable certificate validation (`tlsInsecure`, `sslAllowInvalidCertificates`) except on a short lab, and never with real data.

Certificates expire. Calendar the renewal. An expired certificate stops the application.

`mongodb+srv` on Atlas already implies TLS. A `mongodb://` URI to a self-managed host may not. Be explicit.

TLS is not field-level encryption. The server still sees plaintext documents after TLS ends.

MongoDB query injection happens when you build a query from **strings** or from raw user JSON and you pass it to the driver without control.

A classic mistake is to take HTTP JSON and pass it as the filter:

```javascript
// unsafe if body is the whole filter
collection.find(req.body)
```

The client can send `{ "role": "admin" }` or operators such as `$gt` and `$ne`. The filter is then the attacker's filter.

Another mistake is to concatenate strings into `$where` or into a JavaScript snippet. Do not use `$where` for application filters.

Safe habits:

- Build a filter object in code. Set each field from a **typed** value (string, number, ObjectId).
- Allow-list the fields that the client may filter.
- Convert ids with the driver `ObjectId` type. Do not accept an operator object where you expect an id.
- Use the official driver API, not a string of JavaScript that you `eval`.

This is the same class of bug as SQL injection. The language is BSON, not SQL. The rule is the same: user input is data, not query structure.

Do not log full queries that contain passwords or tokens.

### Questions

#### Theoretical questions

1. What does `bindIp` control?
2. Why is `0.0.0.0/0` on Atlas a risk?
3. What does TLS protect on the wire?
4. Why is `find(req.body)` unsafe?
5. Why is `$where` a poor application tool?

#### Easy practical tasks

1. Write the `bindIp` value or Atlas network mode that you use now.
2. Inspect your connection string. Write if it includes `tls` or `mongodb+srv`.
3. Write a safe filter object for `GET /orders?status=open` with a hardcoded field name.
4. Make a table: unsafe pattern, safe pattern. Add three rows.

#### Medium practical tasks

1. On a lab VM or Docker, bind to localhost only. Confirm a connection from another host fails.
2. In your language, write a function that accepts `status` as a string and returns `{ status }` only if the value is in `{ "open", "closed" }`.
3. Show a unit test: attacker body `{ "status": { "$ne": null } }` must not become the filter.

#### Advanced practical tasks

1. Design network controls for three environments (dev, staging, prod). Include who may add an IP.
2. Enable `requireTLS` on a local replica set with a lab CA. Connect with `mongosh` and the CA file. Record the flags.

---

## Field-level encryption (high-level)

**Client-side field-level encryption (CSFLE)** encrypts selected fields in the driver before the document reaches the server. The server stores ciphertext. Operators who read the collection without keys do not see the plaintext.

**Queryable Encryption** is a later product that lets you encrypt fields and still run some equality (and, in current versions, more) queries on those fields. The server uses encrypted indexes. Read the current manual. The exact query types grow with versions.

High-level pieces:

- A **key management** service (KMIP, cloud KMS, or a local key for labs only)
- A **data encryption key** per field or per store, wrapped by a master key
- A driver that is new enough
- A schema or automatic encryption map that names the fields

Use these tools for high-sensitivity fields (national id, some health data, some payment tokens) when policy requires that the database administrator must not see plaintext.

Limits:

- Not every operator works on ciphertext
- Random encryption (CSFLE) blocks equality queries unless you use the queryable product correctly
- Key loss means data loss
- CPU and size increase

This topic stays high-level. Do not invent a home-grown XOR "encryption" in the application.

TLS plus RBAC plus CSFLE solve different problems. You can need all three.

Atlas has product pages and helpers. Enterprise features and licensing apply in some setups. Read what your edition includes.

### Questions

#### Theoretical questions

1. Where does CSFLE encrypt a field?
2. What does the server store for an encrypted field?
3. What extra problem does Queryable Encryption try to solve?
4. Why is key loss permanent data loss?
5. Why is TLS not enough when an operator can `find` the collection?

#### Easy practical tasks

1. Open the CSFLE and Queryable Encryption pages. Write both URLs.
2. Write three fields that might need encryption and three that might not in a shop app.
3. Make a table: TLS, RBAC, CSFLE. Add what an attacker without keys can still see.
4. Write one sentence: local lab master key vs cloud KMS.

#### Medium practical tasks

1. Read which query types Queryable Encryption supports in your manual version. Write the list.
2. Draw key hierarchy: master key, data key, field ciphertext.
3. Compare hashing a search token vs Queryable Encryption for email lookup. Write two differences.

#### Advanced practical tasks

1. Read automatic vs explicit encryption in the driver. Write when you would use explicit encrypt calls.
2. Write a one-page policy: which fields are encrypted, who holds KMS IAM, how you rotate keys, how you test a restore.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `mongod` config, WiredTiger, and the journal work together on one acknowledged write?
2. How do authentication, RBAC, and network allowlists stop different attackers?
3. When do TLS and field-level encryption both apply on one connection?
4. Why can a least-privilege user still leak data if the handler passes `req.body` to `find`?
5. A teammate disables the journal, binds `0.0.0.0`, uses `root` in the API, and never runs a restore drill. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: config keys, WiredTiger, dump vs snapshot vs Atlas, SCRAM, roles, bindIp, TLS, typed filters, CSFLE vs TLS.
2. Audit your learning cluster: auth on?, TLS on?, IP list?, app role name?, backup method, last drill date or "never".
3. Draw defense layers: network, TLS, auth, RBAC, field encryption, safe API.
4. List five secrets that must not appear in git (URI password, PEM key, KMS key, dump file, trigger secret).

#### Medium practical tasks

1. Run one dump/restore drill on a lab database. Complete a short report: date, duration, errors.
2. Create an app user with least privilege. Point a small script at it. Confirm a privilege error on `dropDatabase`.
3. Fix or rewrite one unsafe find into an allow-listed filter. Add a test.

#### Advanced practical tasks

1. Produce an operations and security standard: config baseline, backup RPO/RTO, drill calendar, role matrix, network, TLS, encryption fields, injection tests.
2. Run a review of a sample repo (yours or a demo). File five findings with severity and a fix type (no exploit steps).
