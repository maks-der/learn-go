# 8. Roles and Security

## Description

This topic shows roles, privileges, and connection security in PostgreSQL 16 and PostgreSQL 17. You learn login roles, `GRANT` and `REVOKE`, default privileges, row-level security, `SECURITY DEFINER` functions, SSL, and why the application must not use a superuser.

Complete topic 7 first. You can maintain a cluster. Complete this topic before you write functions and triggers.

Use one term for each concept. A role is a cluster account. A user is a role with `LOGIN`. A privilege is a named right on an object. A policy is a row filter for row-level security. The application role is the login that the program uses. Superuser is a role with `SUPERUSER`. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## Roles vs users

PostgreSQL stores accounts as roles. A role with `LOGIN` can open a session. A role without `LOGIN` is a group. Older text says "user" for a login role. `CREATE USER name` is `CREATE ROLE name LOGIN`. Prefer `CREATE ROLE` in new scripts so that the attributes are visible.

```sql
CREATE ROLE shop_read NOLOGIN;
CREATE ROLE shop_app LOGIN PASSWORD 'secret';
GRANT shop_read TO shop_app;
```

`shop_app` can use privileges of `shop_read` when `shop_app` inherits membership. PostgreSQL 16 and 17 let you set `INHERIT` and `SET` on each membership. `INHERIT true` adds the group privileges to the member. `SET true` lets the member run `SET ROLE` to that group.

Useful role attributes:

- `LOGIN` / `NOLOGIN` — session or group
- `INHERIT` / `NOINHERIT` — default inherit of group privileges (role-level default)
- `CREATEDB` — create databases
- `CREATEROLE` — create roles
- `REPLICATION` — replication connections
- `BYPASSRLS` — skip row-level security
- `SUPERUSER` — almost all checks off

A role belongs to the cluster. You do not create a role inside one database. Privileges on objects still apply inside each database.

`PUBLIC` is a special implicit grant target. Every role is a member of `PUBLIC`. In PostgreSQL 15 and later, `PUBLIC` does not have `CREATE` on schema `public`. PostgreSQL 16 and 17 keep that rule.

Inspect roles:

```sql
SELECT rolname, rolsuper, rolcanlogin, rolbypassrls
FROM pg_roles
ORDER BY rolname;
```

`\du` in `psql` lists roles.

Do not share one login role across unrelated programs. Do not put a production password in a handbook or a git commit.

### Questions

#### Theoretical questions

1. What is the difference between a role and a user in PostgreSQL?
2. What SQL does `CREATE USER shop_app` send?
3. Does a role belong to one database or to the cluster?
4. What does `LOGIN` allow?
5. What is `PUBLIC`?

#### Easy practical tasks

1. Create a `NOLOGIN` group role and a `LOGIN` member. List both with `\du`.
2. Write a two-column table: attribute and one-sentence meaning. Add `LOGIN`, `CREATEDB`, and `SUPERUSER`.
3. Run `SELECT current_user, session_user;`. Write what each name is.
4. Find `rolcanlogin` for `postgres` and for your new group role.

#### Medium practical tasks

1. Grant the group to the login role. Connect as the login role. Show inherited privileges with `\du` and a test `SELECT` after you grant a table privilege to the group.
2. Read the 16 or 17 docs for `GRANT role_name TO`. Write what `INHERIT` and `SET` mean on a membership.
3. Create two login roles. Confirm that one role cannot `SET ROLE` to the other until you grant membership.

#### Advanced practical tasks

1. Design roles for `app`, `migrator`, and `readonly`. Write attributes and group membership. Do not use `SUPERUSER` for those three.
2. Read "Database Roles" in the official docs. Add two attributes that this section did not list. Write one sentence each.

---

## `GRANT` / `REVOKE` and default privileges

A privilege is a right that you grant on an object. The owner of the object has all privileges. Other roles need `GRANT`. `REVOKE` removes a privilege.

Schema privileges:

- `USAGE` — look up objects in the schema
- `CREATE` — create objects in the schema

Table privileges include `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `TRUNCATE`, `REFERENCES`, and `TRIGGER`. Sequence privileges include `USAGE`, `SELECT`, and `UPDATE`. Identity columns use a sequence. An `INSERT` that uses identity needs sequence `USAGE`.

```sql
CREATE SCHEMA shop;
CREATE TABLE shop.items (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL
);

GRANT USAGE ON SCHEMA shop TO shop_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE shop.items TO shop_app;
GRANT USAGE, SELECT ON SEQUENCE shop.items_id_seq TO shop_app;
```

You can grant on every current table in a schema:

```sql
GRANT SELECT ON ALL TABLES IN SCHEMA shop TO shop_read;
GRANT USAGE ON ALL SEQUENCES IN SCHEMA shop TO shop_app;
```

`ALL TABLES` applies to tables that exist now. It does not apply to tables that you create later. Default privileges cover later objects.

`CONNECT` is a database privilege. Without `CONNECT`, the role cannot open a session on that database.

```sql
GRANT CONNECT ON DATABASE shop TO shop_app;
REVOKE CONNECT ON DATABASE shop FROM PUBLIC;
```

`psql` shows privileges with `\dp` (tables) and `\dn+` (schemas).

`ALTER DEFAULT PRIVILEGES` sets grants for objects that a role creates later. It does not change objects that already exist.

```sql
ALTER DEFAULT PRIVILEGES FOR ROLE migrator IN SCHEMA shop
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO shop_app;

ALTER DEFAULT PRIVILEGES FOR ROLE migrator IN SCHEMA shop
    GRANT USAGE, SELECT ON SEQUENCES TO shop_app;
```

`FOR ROLE migrator` means: when `migrator` creates a table in `shop`, apply these grants. Default privileges are per database. `psql` `\ddp` lists them.

A typical split:

- `migrator` owns tables and runs `CREATE TABLE`
- `shop_app` has DML only
- default privileges grant DML from `migrator` to `shop_app`

Do not grant `ALL` on every object to the application role when the program only needs `SELECT` and `INSERT`. Do not grant table privileges without `USAGE` on the schema. Do not assume that `GRANT ON ALL TABLES` is enough after you add tables.

Official: [https://www.postgresql.org/docs/17/sql-grant.html](https://www.postgresql.org/docs/17/sql-grant.html).

### Questions

#### Theoretical questions

1. What does `USAGE` on a schema allow?
2. Why does an identity `INSERT` need a sequence privilege?
3. What does `GRANT SELECT ON ALL TABLES IN SCHEMA shop` miss for new tables?
4. What objects does `ALTER DEFAULT PRIVILEGES` change at the time you run it?
5. Why must you name `FOR ROLE` on default privileges?

#### Easy practical tasks

1. Create a schema and a table. Grant `USAGE` and `SELECT` to a login role. Select one row as that role.
2. Revoke `SELECT`. Try the same `SELECT` as that role. Record the error.
3. List privileges with `\dp` on that table.
4. Set default `SELECT` on tables in one schema for a reader role. Create a new table as the owner. Check `\dp`.

#### Medium practical tasks

1. Grant only `INSERT` (plus schema `USAGE` and sequence `USAGE`). Prove that `SELECT` fails and `INSERT` works.
2. `REVOKE CONNECT ON DATABASE ... FROM PUBLIC`. Connect as a role that has no `CONNECT`. Record the error. Grant `CONNECT` and retry.
3. Create a table before defaults. Create a table after defaults. Compare `\dp` on both.

#### Advanced practical tasks

1. Write a grant script for `readonly` and `app` on one schema. Include `ALL TABLES` and `ALTER DEFAULT PRIVILEGES`. Test each role.
2. Read the privilege list for tables, schemas, and sequences in the 17 docs. Make a three-column table: object, privilege, one example statement that needs it.

---

## Row-level security

Row-level security (RLS) adds a policy on a table. After you enable RLS, a normal role sees only rows that pass `USING`. Inserts and updates must pass `WITH CHECK`.

```sql
ALTER TABLE shop.orders ENABLE ROW LEVEL SECURITY;

CREATE POLICY orders_by_tenant ON shop.orders
    FOR ALL
    TO shop_app
    USING (tenant_id = current_setting('app.tenant_id')::bigint)
    WITH CHECK (tenant_id = current_setting('app.tenant_id')::bigint);
```

The session sets a parameter, then queries:

```sql
SET app.tenant_id = '42';
SELECT * FROM shop.orders;
```

`current_setting('app.tenant_id')` fails if the parameter is missing and you do not use `missing_ok`. Keep the policy expression simple.

The table owner bypasses RLS unless you force it:

```sql
ALTER TABLE shop.orders FORCE ROW LEVEL SECURITY;
```

A role with `BYPASSRLS` or `SUPERUSER` also skips RLS. The application role must not have those attributes.

Policies can be `FOR SELECT`, `FOR INSERT`, `FOR UPDATE`, `FOR DELETE`, or `FOR ALL`. You can attach more than one policy. For the same command, PostgreSQL combines permissive policies with `OR` (default). Restrictive policies use `AND`. Start with one permissive policy per table.

`EXPLAIN` can show a policy filter as a qualifier. Test as the application role, not as the owner.

Do not treat RLS as the only control. You still need `GRANT`. A role without `SELECT` sees no rows. Do not enable RLS and forget `FORCE` if the owner is also the application role.

Official: [https://www.postgresql.org/docs/17/ddl-rowsecurity.html](https://www.postgresql.org/docs/17/ddl-rowsecurity.html).

### Questions

#### Theoretical questions

1. What does `ENABLE ROW LEVEL SECURITY` change?
2. What is the difference between `USING` and `WITH CHECK`?
3. Who bypasses RLS by default?
4. What does `FORCE ROW LEVEL SECURITY` do?
5. Does RLS replace `GRANT`?

#### Easy practical tasks

1. Create a table with `tenant_id`. Insert two tenants. Enable RLS. Create a policy. Set `app.tenant_id` and select.
2. Change `app.tenant_id` to the other tenant. Select again. Write the row count.
3. Query as the table owner without `FORCE`. Write whether the owner sees all rows.
4. List policies with `\d` on the table or `SELECT * FROM pg_policies`.

#### Medium practical tasks

1. Add `FORCE ROW LEVEL SECURITY`. Repeat the owner `SELECT`. Write the difference.
2. Create a `FOR SELECT` policy and a separate `FOR INSERT` policy. Prove that a row that fails `WITH CHECK` cannot insert.
3. Run `EXPLAIN` on a `SELECT` as `shop_app` with a tenant set. Find the policy filter.

#### Advanced practical tasks

1. Build a three-tenant sample. Write policies so that `shop_app` sees one tenant and `shop_read` sees all tenants (or none). Test both roles.
2. Read "Permissive and Restrictive Policies" in the docs. Write one example of each. State why a beginner starts with one permissive policy.

---

## `SECURITY DEFINER` danger

A function runs with the privileges of the caller by default (`SECURITY INVOKER`). `SECURITY DEFINER` runs with the privileges of the owner. Use it only when the caller must not have those privileges, and the function body is a tight, reviewed operation.

Danger: the function can touch objects that the caller cannot touch. If the function uses unqualified names, `search_path` can point to a schema that the caller controls. The function can then call the wrong function or use the wrong table.

Safe pattern for PostgreSQL 16 and 17:

```sql
CREATE FUNCTION shop.touch_item(p_id bigint)
RETURNS void
LANGUAGE sql
SECURITY DEFINER
SET search_path = pg_catalog, shop
AS $$
    UPDATE shop.items SET name = name WHERE id = p_id;
$$;
```

`SET search_path` on the function fixes the path for the call. Include `pg_catalog`. Include only schemas that you trust. Some teams use `SET search_path = ''` and qualify every name.

Also:

- grant `EXECUTE` only to roles that need the function
- revoke `EXECUTE` from `PUBLIC` when you can
- keep the body short
- prefer `SECURITY INVOKER` when `GRANT` is enough
- do not put dynamic SQL together with untrusted text

`SECURITY DEFINER` plus RLS is easy to get wrong. The function owner can bypass RLS. Test as the real caller.

Do not copy a `SECURITY DEFINER` example from the internet without a fixed `search_path`. Do not use `SECURITY DEFINER` to hide a superuser from the application. The function owner must not be a superuser unless the task has no other design.

Official warning: [https://www.postgresql.org/docs/17/sql-createfunction.html](https://www.postgresql.org/docs/17/sql-createfunction.html) (see `SECURITY`).

### Questions

#### Theoretical questions

1. What is the difference between `SECURITY INVOKER` and `SECURITY DEFINER`?
2. Why is `search_path` dangerous on a definer function?
3. How do you fix `search_path` on `CREATE FUNCTION`?
4. Who must own a `SECURITY DEFINER` function in a safe design?
5. When is `SECURITY INVOKER` the better choice?

#### Easy practical tasks

1. Create an invoker function and a definer function that select from a table. Grant `EXECUTE` only. Compare who can read the table directly.
2. Show `prosecdef` in `pg_proc` for both functions.
3. Revoke `EXECUTE` from `PUBLIC` on the definer function. Grant `EXECUTE` to one role.
4. Write four sentences: invoker, definer, `search_path`, `EXECUTE`.

#### Medium practical tasks

1. Create a definer function without `SET search_path`. Read the official warning. Add `SET search_path`. Show `\df+` for the setting.
2. Call a definer function as a role that has no `SELECT` on the table. Confirm that the function still runs if the owner has `SELECT`.
3. Plan (do not implement an attack) a checklist: owner, `search_path`, grants, RLS, dynamic SQL.

#### Advanced practical tasks

1. Write a definer function that writes one allowed column. Refuse other columns. Test a caller that has only `EXECUTE`.
2. Read the 17 docs on `CREATE FUNCTION` security. List three extra settings (`SET` clauses) that you can attach. Write when you would set `row_security`.

---

## SSL and why the app role is not a superuser

SSL encrypts the TCP session between the client and the server. PostgreSQL 16 and 17 use TLS. Older text still says SSL. The setting names keep `ssl`.

Server side (`postgresql.conf`):

- `ssl = on`
- `ssl_cert_file` and `ssl_key_file` — server certificate and key
- optional `ssl_ca_file` — when you verify client certificates

`pg_hba.conf` can require SSL:

```text
hostssl shop shop_app 10.0.0.0/8 scram-sha-256
```

`hostssl` matches only SSL connections. `hostnossl` matches only non-SSL. `host` matches both.

Client `sslmode` (libpq):

| Value | Meaning |
| --- | --- |
| `disable` | no SSL |
| `allow` / `prefer` | try both; `prefer` is the libpq default |
| `require` | SSL required; no certificate check |
| `verify-ca` | SSL plus server cert signed by your CA |
| `verify-full` | `verify-ca` plus host name match |

Production clients use `verify-full` when you control the CA and the host name. `require` stops clear text. It does not stop a wrong server.

```text
postgresql://shop_app@db.example:5432/shop?sslmode=verify-full
```

`psql` shows SSL after connect when you run `\conninfo`.

A superuser bypasses most privilege checks. A superuser can read and write files that the OS user can read, run `COPY` from files, change roles, and skip RLS. The default role `postgres` is a superuser.

The application role must be a normal login:

```sql
CREATE ROLE shop_app LOGIN PASSWORD 'secret';
-- no SUPERUSER, no BYPASSRLS, no REPLICATION
```

Use a superuser (or a tight admin role) for install, some `CREATE EXTENSION` work, full `pg_dumpall` of roles, and crash recovery. Use a migrator role with `CREATE` on the application schema for `CREATE TABLE`. That role is still not a superuser.

```sql
SELECT current_user, session_user;
SELECT rolsuper FROM pg_roles WHERE rolname = current_user;
```

If the program connects as `postgres`, every SQL injection or stolen password is a full cluster compromise. Cloud vendors often hide the true superuser and give a powerful admin role. Treat that admin role as "not for the app".

SSL does not replace passwords or `pg_hba.conf`. Do not expose `5432` on the public internet and then rely on SSL alone. Topic 14 covers network rules. Do not grant `SUPERUSER` to debug a privilege error. Grant the missing privilege.

Official SSL: [https://www.postgresql.org/docs/17/ssl-tcp.html](https://www.postgresql.org/docs/17/ssl-tcp.html).

### Questions

#### Theoretical questions

1. What does `ssl = on` enable on the server?
2. What extra check does `verify-full` add that `require` does not?
3. What checks does a superuser skip?
4. Why must the application role lack `SUPERUSER`?
5. Does SSL replace a password?

#### Easy practical tasks

1. Connect with `psql` and run `\conninfo`. Write whether the session uses SSL.
2. Write a URI that sets `sslmode=require` for a local test.
3. Run `SELECT rolsuper FROM pg_roles WHERE rolname = current_user;` as your app role if you have one.
4. Find `ssl` with `SHOW ssl;` if you have permission.

#### Medium practical tasks

1. Compare a `host` line and a `hostssl` line. Write which clients can match each line.
2. Set `sslmode=disable` against a server that requires SSL (lab). Record the error.
3. Create a non-superuser login. Prove that `CREATE ROLE` or `CREATE EXTENSION` fails. Grant the missing privilege instead of `SUPERUSER`.

#### Advanced practical tasks

1. Enable SSL on a lab server (self-signed is acceptable). Connect with `sslmode=require`. Save `\conninfo`.
2. Write a one-page policy: app role attributes, migrator role, SSL `sslmode`, and `pg_hba.conf` `hostssl`.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a group role, `GRANT` on tables, and default privileges work together when the migrator adds a table?
2. Why is RLS a second filter after `GRANT`, not a replacement?
3. What two mistakes turn `SECURITY DEFINER` into a privilege leak?
4. Why is `sslmode=require` not enough for production if you control DNS?
5. What blast radius does a stolen `shop_app` password have if that role is not a superuser?

#### Easy practical tasks

1. Create `shop_app` and `shop_read`. Grant schema `USAGE` and table `SELECT` to `shop_read`. Connect as each role. Run `SELECT current_user, rolsuper FROM pg_roles WHERE rolname = current_user`.
2. Write a cheat sheet: role, `GRANT`, defaults, RLS, definer, SSL, superuser.
3. List roles with `\du` and table privileges with `\dp`.
4. Write a connection URI with `sslmode=verify-full` and no password in the string.

#### Medium practical tasks

1. Enable RLS on a tenant table. Set default privileges for new tables. Test insert as `shop_app` with a tenant parameter.
2. Revoke `CONNECT` from `PUBLIC`. Prove that only granted roles connect. Restore your lab access first if you work on a shared cluster.
3. Create a definer function with fixed `search_path`. Call it as `shop_app`. Show `\conninfo` SSL state in the same notes.

#### Advanced practical tasks

1. Script a least-privilege shop: migrator owns objects, app has DML, reader has `SELECT`, RLS on `tenant_id`, no superuser, `hostssl` comment in a sample `pg_hba.conf`.
2. Read "Database Roles", "Privileges", "Row Security", and "SSL" in the 17 docs. Add one official warning that this topic did not quote.
