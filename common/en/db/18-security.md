# 18. Security

## Description

This topic shows how you control who connects, what they can do, and how you protect data in transit and at rest. You learn authentication, authorization, least privilege, secrets, TLS, encryption at rest, injection, broad grants, and exposed admin ports.

Use one term for each concept. Authentication proves an identity. Authorization grants rights to that identity. A secret is a password, token, or key that you must not publish. Complete this topic after you can connect with a driver and write parameterized SQL. Security is a property of the whole path, not of one password field.

This path is vendor-neutral. Privilege names differ. The rules do not.

---

## Authentication vs authorization

Authentication answers "who is this session?" The DBMS checks a user name and a credential (password, certificate, token, or an external directory). A failed authentication does not start a session.

Authorization answers "what may this session do?" The DBMS checks privileges: connect, read a table, write a table, create objects, manage roles.

```text
Connect request --> authenticate --> session identity
SQL statement   --> authorize    --> allow or reject
```

A correct password with no `SELECT` privilege still fails `SELECT`. A `GRANT` without a successful login never applies.

Typical identities:

- a human admin
- an application user
- a read-only reporter
- a migration user (schema change)

Do not share one user for all of these jobs. You cannot revoke a leaked application password without also locking the admin. You cannot audit who changed a row if every program uses `postgres` or `sa`.

Authentication methods differ: password, peer/socket, LDAP, certificate. Use a strong method on networks that you do not fully trust. `trust` authentication (any connector is accepted) is only for a locked-down lab.

Do not confuse operating-system login with DBMS login. They can map, but they are not the same check unless you configured that map.

Log authentication failures. Many failures from one host can be an attack or a bad connection string.

### Questions

#### Theoretical questions

1. What question does authentication answer?
2. What question does authorization answer?
3. Why does a valid password not imply `SELECT` success?
4. Why must an application and an admin use different users?
5. When is `trust` authentication acceptable?

#### Easy practical tasks

1. Write the two-step path from this section.
2. List four identities for a class shop. Assign one job each.
3. Find how you create a user (role) in your DBMS.
4. Write one authentication failure and one authorization failure as error stories.

#### Medium practical tasks

1. Create a user that can connect but cannot `SELECT` a table. Prove both: connect works, `SELECT` fails.
2. Read the official authentication-method list. Write which method your local install uses.
3. Enable or locate auth-failure logs. Fail a login on purpose. Find the log line.

#### Advanced practical tasks

1. Write a one-page identity plan: human, app, report, migrate. Include how you rotate each credential.
2. Map one external directory idea (LDAP or SSO) at a high level. Write what the DBMS still authorizes locally.

---

## Least privilege roles

Least privilege means a role has only the rights that it needs. The application user can `SELECT`, `INSERT`, `UPDATE`, and `DELETE` on the tables that the app uses. It cannot `DROP TABLE` on production. It cannot read `pg_authid` or other password catalogs if you can prevent that.

A role is a named set of privileges. Users are roles that can log in (product model differs). You grant a role to a user. You change the role once. All members receive the change.

```sql
-- Shape only; verbs and names differ by product
GRANT SELECT, INSERT, UPDATE, DELETE ON orders TO app_shop;
REVOKE DELETE ON invoices FROM app_shop;
```

Start from zero. Grant what the program needs. Do not grant `SUPERUSER`, `DBA`, or `sysadmin` to an application.

Object owners and `PUBLIC` defaults matter. Some products grant wide rights to `PUBLIC`. Revoke those rights on sensitive objects. Read the default ACL.

Migrations need a stronger user. Run migrations as that user. The application does not use that user at run time.

Read-only roles must not write. A report replica still needs a user that cannot `DELETE`.

Do not grant `SELECT` on all tables if the app only needs five tables. Do not use one write role for an intern tool and for payments.

Review grants when you add a table. A new table can inherit defaults that are too wide. Grant explicitly.

### Questions

#### Theoretical questions

1. What does least privilege require?
2. Why do you grant a role instead of many raw grants on each person?
3. Why must the application not run as a superuser?
4. Why do migrations use a different user than run time?
5. What risk does a `PUBLIC` default grant create?

#### Easy practical tasks

1. Write a privilege list for `app_shop` on `products` and `orders`.
2. Write three rights that `app_shop` must not have in production.
3. Find `GRANT` and `REVOKE` in your DBMS manual.
4. Make a table: role, login yes/no, typical grants.

#### Medium practical tasks

1. Create `app_shop` and `report_shop`. Grant write versus `SELECT` only. Prove with two sessions.
2. List current grants on a practice table (`\dp` or the vendor equivalent). Write what `PUBLIC` has.
3. Add a new table. Check whether `app_shop` can read it. Grant only if you intend that.

#### Advanced practical tasks

1. Write a one-page grant standard: roles, default revoke, new-table checklist, migration user.
2. Review a sample database ACL. List every grant that is wider than the app needs. Write the `REVOKE` list (do not apply it on a shared server).

---

## Secrets for connection strings

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

Certificates and API tokens are secrets too. Treat a client TLS key as you treat a password.

`.env` files help local work. Do not commit `.env`. Add `.env` to the ignore file. A committed `.env.example` contains placeholders only.

Cloud metadata and CI variables leak if you grant a wide role to a pipeline. Use short-lived credentials when the platform offers them.

If a secret leaked, change it. A leak in an old commit stays in Git history. Rotate. Consider the repository as public from that moment.

Do not send production connection strings in email. Do not put them in a ticket body.

### Questions

#### Theoretical questions

1. Why is a password in a Git repository a leak?
2. Where may a process read a production password?
3. Why must development and production use different passwords?
4. Why must you rotate after a leak even if you delete one file?
5. What belongs in `.env.example`?

#### Easy practical tasks

1. Write a connection URI with `PASSWORD` as a placeholder.
2. Add `.env` to a ignore file if you use one. Write an `.env.example` with empty values.
3. List five places a teammate might paste a URI (repo, chat, ticket, image, log). Mark each as forbidden for real secrets.
4. Find whether your driver logs the URI. Write how you redact it.

#### Medium practical tasks

1. Change a small program to read the password from the environment. Prove that the source file has no password.
2. Search your practice repo for `postgres://` or `password=`. Record hits. Remove any real secret.
3. Write a rotate procedure of six steps: create new password, update store, reload app, revoke old, test, record.

#### Advanced practical tasks

1. Write a one-page secret standard: stores, rotation, CI, logs, and response to a leaked commit.
2. Use a local secret manager or the OS key store for one connection. Document the setup without writing the secret.

---

## Encryption in transit (TLS)

TLS encrypts the network bytes between the client and the DBMS. An observer on the network cannot read the SQL text or the rows. TLS also authenticates the server when the client verifies the certificate.

Without TLS on an untrusted network, a packet capture can show passwords and personal data. "The port is on a private VLAN" reduces risk. It does not encrypt. Operators and a mis-routed packet can still see cleartext.

```text
Client --TLS--> DBMS
         ^
   verify server certificate
```

Modes (names differ):

- disable: no TLS
- allow/prefer: may use TLS, may fall back to cleartext
- require: TLS, but the client may skip full certificate checks
- verify-ca / verify-full: TLS plus certificate checks (host name)

Use verify-full (or the product equivalent) in production. A `require` mode without name checks can still accept a different server.

Install a certificate on the DBMS that a client trust store accepts. Self-signed certificates need an explicit CA in the client. Do not set "skip verify" in production to silence an error.

TLS does not replace authentication. TLS does not replace grants. TLS does not encrypt data files on disk.

Internal networks still benefit from TLS. Many incidents start inside the network.

Do not expose the DBMS port on the public internet even with TLS. TLS is not a firewall.

### Questions

#### Theoretical questions

1. What does TLS protect on the wire?
2. What extra check does verify-full add over "TLS required"?
3. Why is a private VLAN not a substitute for TLS?
4. What does TLS not replace?
5. Why is "skip verify" a production risk?

#### Easy practical tasks

1. Find the client TLS mode names for your driver.
2. Draw a packet sniffer versus a TLS session. Write one sentence on what the sniffer sees in each case.
3. Write the production mode that you choose (verify-full or the product name).
4. List two failures: expired certificate, host-name mismatch.

#### Medium practical tasks

1. Connect with TLS disabled and with TLS required on a learning instance if you can. Write the connection option.
2. Read how the server presents a certificate. Write where the cert file lives (docs).
3. Capture the idea of MITM: a client that does not verify the name. Write four sentences. Do not attack a network that you do not own.

#### Advanced practical tasks

1. Write a one-page TLS standard: modes, CA, rotation, and a ban on skip-verify.
2. Enable TLS on a disposable DBMS. Connect with verify-full. Record the failed connect when the name is wrong.

---

## Encryption at rest (high-level)

Encryption at rest protects data files and backups when someone obtains the disk or the backup file. The bytes on the volume are ciphertext without the key.

Layers (you can use more than one):

| Layer | What it encrypts | Who holds the key |
| --- | --- | --- |
| Full-disk or volume | The whole volume | OS or cloud KMS |
| Tablespace or TDE | Database files | DBMS plus KMS |
| Backup encryption | Dump or snapshot files | Backup tool plus KMS |
| Application field encryption | Selected columns | Application |

Volume encryption is the usual first step in a cloud VM. If an attacker gets a running server and the mounted volume, the OS already decrypted the disk. Volume encryption helps when the disk is detached or stolen.

Transparent data encryption (TDE) or equivalent encrypts files that the DBMS writes. The DBMS decrypts pages for a session that is allowed to read. An attacker with only the files still needs the key.

Column encryption in the application hides values from the DBMS and from many admins. You lose ordinary `WHERE` and indexes on the cleartext unless you design searchable encryption. Use it for a few high-sensitivity fields.

Keys are secrets. A key on the same disk as the data is a weak design. Use a key management service. Rotate keys with the product procedure.

Encryption at rest does not hide data from a connected application user. Grants and authentication still apply.

Do not treat "we have disk encryption" as a complete backup policy. Encrypt the backups too. Store the key away from the backup files.

### Questions

#### Theoretical questions

1. What threat does encryption at rest address?
2. Why does volume encryption not stop an attacker on a running mounted server?
3. What does TDE (or equivalent) protect?
4. What query cost can application-level column encryption add?
5. Why must the key not live only beside the data files?

#### Easy practical tasks

1. Fill the four-layer table in your notes with one extra example.
2. Write one threat that TLS covers and one threat that at-rest covers.
3. Find whether your DBMS or cloud vendor documents TDE or volume encryption.
4. Write why backups need their own encryption setting.

#### Medium practical tasks

1. Enable volume encryption on a learning VM or note how a laptop disk is encrypted. Write who unlocks the volume.
2. Compare a cleartext dump and an encrypted backup option in the docs. Write the restore extra step (the key).
3. Pick one column (national id). Write why you would encrypt it in the app or why grants plus TLS are enough for a class project.

#### Advanced practical tasks

1. Write a one-page key policy: where keys live, who rotates, what you do if a backup tape is lost.
2. Design field encryption for one column: key id, rotation, and which queries you lose. Do not implement a weak home-made cipher.

---

## Injection, overly broad grants, exposed admin ports

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

1. How do authentication, authorization, and TLS fail independently on the same connection string?
2. When does a secret leak become a full data breach only because grants were too wide?
3. How do encryption in transit and encryption at rest cover different attackers?
4. Which controls stop injection, and which controls only limit the blast radius?
5. A teammate puts the admin URI in the app, opens port 5432 on `0.0.0.0`, and disables TLS on a "private" cloud subnet. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: authn versus authz, least privilege, secrets, TLS, at rest, three common failures.
2. Create a non-admin user. Connect. Run `SELECT 1`. Show that `DROP DATABASE` fails.
3. Redact a sample URI for a README. Keep host and database name. Remove the password.
4. Write the listen address and TLS mode of your learning instance.

#### Medium practical tasks

1. Build a mini app user: env-based password, grants on two tables only, TLS if available, parameterized SQL.
2. Write a threat table: stolen laptop disk, packet capture, leaked Git password, SQL injection. Map one control each.
3. Run a privilege review on your `learn` database. Produce a grant list and a revoke plan.

#### Advanced practical tasks

1. Write a security review of a small project against this topic. File findings by severity. Fix the ones you own.
2. Design production controls for a payments schema: roles, secrets, TLS verify-full, encrypted backups, no public port. Draw the network.
