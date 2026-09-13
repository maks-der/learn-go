# 17. Security

## Description

This topic shows how you protect a MongoDB deployment. You authenticate users. You assign **roles**. You limit **network** access. You enable **TLS**. You learn field-level encryption at a high level. You use typed driver APIs so that you do not build queries from raw strings. Complete drivers and administration ideas first.

Use one term for each concept. **Authentication** proves who the client is. **Authorization** (RBAC) decides what that user can do. **TLS** encrypts the network. **Queryable Encryption** and **client-side field-level encryption** protect field values so that the server can store ciphertext.

A learning database on localhost still needs a password before you expose a port. Atlas needs a database user and an IP access list (or a private network) before you open the cluster to an application.

---

## Authentication (SCRAM, x.509)

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

Atlas: you create a database user in the UI. Atlas uses SCRAM for those users by default. You can add other methods that Atlas documents (LDAP, AWS IAM, OIDC on some tiers). You do not run `mongod` flags by hand.

Enable authorization on Community Server (`security.authorization: enabled`). Without it, a client that can reach the port can do anything.

`localhost exception` lets you create the first user on a fresh self-managed server. Read that page before you lock yourself out.

Do not use the same password on Atlas and in a public gist.

Do not create a user with the `root` role for the application.

### Questions

#### Theoretical questions

1. What does authentication prove?
2. What does SCRAM use as credentials?
3. What does an x.509 client present?
4. Why must you enable authorization on a self-managed server?
5. Why is `root` a poor role for an application user?

#### Easy practical tasks

1. Open the authentication page. Write the URL and the SCRAM mechanism name that the page prefers.
2. On Atlas, view Database Access. Write the user name (not the password) and the auth method.
3. Write a URI shape with a user name and a placeholder password.
4. Make a table: SCRAM, x.509. Add credential type and one typical use.

#### Medium practical tasks

1. Enable authorization on a local `mongod`. Create a user. Connect with that user. Confirm a bad password fails.
2. Create two users: `app` with `readWrite` on `learn`, and `reader` with `read`. Prove the reader cannot insert.
3. Read the localhost exception page. Write when it ends.

#### Advanced practical tasks

1. Read member x.509 vs client x.509. Write two sentences each.
2. Compare Atlas database users with self-managed `admin` users. Write how rotation of a password differs.

---

## Role-based access control

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

Auditing (Enterprise / Atlas) can record auth events. Read if your edition includes it.

List privileges with `rolesInfo` and the Atlas UI.

### Questions

#### Theoretical questions

1. What does a role contain?
2. Why is `readWrite` on one database better than `root` for an API?
3. When do you create a custom role?
4. Why must two applications use two users?
5. What is a privilege resource?

#### Easy practical tasks

1. Open the built-in roles page. Write the URL and five role names.
2. Write the minimum role for a reporting job that only finds in `shop`.
3. Make a table: user, role, database. Add app, analyst, operator.
4. In Atlas or `rolesInfo`, inspect one user. Write the roles that you see.

#### Medium practical tasks

1. Create a custom role that can `find` and `insert` on `shop.orders` only. Assign it. Prove `shop.items` insert fails.
2. Grant `clusterMonitor` to a metrics user. Run one status command. Write which command you used.
3. Review an existing user. Remove one extra role. Write what broke, or that nothing broke.

#### Advanced practical tasks

1. Design roles for app, migrator, backup, and human admin. Write actions that each must not have.
2. Read privilege actions for `changeStream` and `compact`. Write which role you would use for each job.

---

## Network exposure: bind IP, VPC, IP allowlists

A MongoDB port on the public internet without controls is a common incident. Limit who can open a TCP connection.

**Self-managed `bindIp`.** Bind to `127.0.0.1` for a local-only lab. Bind to a private interface for a VPC. Do not bind to `0.0.0.0` without authentication, TLS, and a firewall.

**Firewall.** Allow port `27017` (or your port) only from application subnets. Deny the rest.

**Atlas IP access list.** Atlas refuses clients that are not on the list (unless you use a private endpoint and the matching mode). Add the application egress IPs. Do not add `0.0.0.0/0` except on a short-lived learning cluster that has no real data, and remove it after.

**VPC and private networking.** Place `mongod` in a private subnet. On Atlas, use VPC peering or a private endpoint so that traffic does not use the public internet.

**SRV and DNS.** Clients resolve Atlas hosts. Still restrict IP or private connectivity. DNS is not access control.

Network control is not authentication. You need both. A stolen password from a laptop on the allowlist still works. A public port with a strong password is still a scan target.

Do not publish `mongod` on a cloud VM with the vendor "allow all" security group.

Do not keep a personal home IP on the Atlas list after you leave the project.

### Questions

#### Theoretical questions

1. What does `bindIp` control?
2. Why is `0.0.0.0/0` on Atlas a risk?
3. How does a VPC peering path differ from a public IP list?
4. Why is a firewall not a substitute for authentication?
5. Why must you remove old laptop IPs from an allowlist?

#### Easy practical tasks

1. Write the `bindIp` value or Atlas network mode that you use now.
2. Open Atlas Network Access. Write how many IP entries you see (no need to publish them).
3. Open the bind-IP and Atlas network pages. Write both URLs.
4. Draw: internet, firewall, `mongod`. Mark the allowed arrow.

#### Medium practical tasks

1. On a lab VM or Docker, bind to localhost only. Confirm a connection from another host fails.
2. Add and remove a test IP on Atlas (or write the steps if you cannot change the project).
3. Compare public TLS + allowlist vs private endpoint. Write two benefits each.

#### Advanced practical tasks

1. Design network controls for three environments (dev, staging, prod). Include who may add an IP.
2. Read Atlas private endpoint limits. Write one routing pitfall (wrong region or missing DNS).

---

## TLS

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

Do not mix "TLS to the database" with "HTTPS to the user" and think you are done. The path from app server to `mongod` needs TLS too.

### Questions

#### Theoretical questions

1. What does TLS protect on the wire?
2. Does Atlas require TLS for clients?
3. Why is `tlsInsecure` dangerous?
4. Why do certificates need a renewal date?
5. Does TLS hide document fields from `mongod`?

#### Easy practical tasks

1. Inspect your connection string. Write if it includes `tls` or `mongodb+srv`.
2. Open the TLS configuration page. Write the URL.
3. Write four sentences: TLS, certificate, CA, expire.
4. Make a table: Atlas, local lab without TLS, production self-managed. Add "TLS required?".

#### Medium practical tasks

1. Connect with an official driver and TLS options visible in code. Write the option names.
2. Read replica-set TLS (member certificates). Write why members need TLS too.
3. Plan a certificate rotation without downtime at a high level (new cert, reload or rolling restart).

#### Advanced practical tasks

1. Enable `requireTLS` on a local replica set with a lab CA. Connect with `mongosh` and the CA file. Record the flags.
2. Write a TLS standard: min version, forbidden insecure flags, who owns renewal, Atlas vs self-managed.

---

## Field-level encryption / Queryable Encryption (high-level)

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

## No injection via string-built queries (use typed APIs)

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
- For SQL-style libraries that sit on MongoDB, still bind parameters. Prefer the official driver.

ODM libraries can help if they treat user input as values, not as full query documents. They can also hide a `$where`. Read the docs.

This is the same class of bug as SQL injection. The language is BSON, not SQL. The rule is the same: user input is data, not query structure.

Do not log full queries that contain passwords or tokens.

### Questions

#### Theoretical questions

1. Why is `find(req.body)` unsafe?
2. What is an allow-list of filter fields?
3. Why is `$where` a poor application tool?
4. How is this similar to SQL injection?
5. Why must an id parameter become an `ObjectId` (or the type you use) in code?

#### Easy practical tasks

1. Write a safe filter object for `GET /orders?status=open` with a hardcoded field name.
2. Open a MongoDB injection or security checklist page. Write the URL.
3. Make a table: unsafe pattern, safe pattern. Add three rows.
4. Write why `$gt` in user JSON is a problem if you pass the JSON as the filter.

#### Medium practical tasks

1. In your language, write a function that accepts `status` as a string and returns `{ status }` only if the value is in `{ "open", "closed" }`.
2. Show a unit test: attacker body `{ "status": { "$ne": null } }` must not become the filter.
3. Review one handler in a sample app. Mark any raw body-to-find path.

#### Advanced practical tasks

1. Read operator injection (`$regex`, `$expr`) risks. Write two extra allow-list rules.
2. Write a team standard: typed filters, no `$where`, no `eval`, logging redaction, code review checks.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do authentication, RBAC, and network allowlists stop different attackers?
2. When do TLS and field-level encryption both apply on one connection?
3. Why can a least-privilege user still leak data if the handler passes `req.body` to `find`?
4. How do SCRAM users, x.509 members, and Atlas UI users fit one production picture?
5. A teammate binds `0.0.0.0`, disables TLS validation, and uses `root` in the API. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: SCRAM, x.509, built-in roles, bindIp, Atlas list, TLS flags, CSFLE vs TLS, typed filters.
2. Audit your learning cluster: auth on?, TLS on?, IP list?, app role name?.
3. Draw defense layers: network, TLS, auth, RBAC, field encryption, safe API.
4. List five secrets that must not appear in git (URI password, PEM key, KMS key, trigger secret, dump file).

#### Medium practical tasks

1. Create an app user with least privilege. Point a small script at it. Confirm a privilege error on `dropDatabase`.
2. Write a secure connection checklist of ten items for Atlas and ten for Community Server.
3. Fix or rewrite one unsafe find into an allow-listed filter. Add a test.

#### Advanced practical tasks

1. Produce a security standard: auth mechanisms, role matrix, network, TLS, encryption fields, injection tests, secret rotation.
2. Run a review of a sample repo (yours or a demo). File five findings with severity and a fix type (no exploit steps).
