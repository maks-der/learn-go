# 22. Ecosystem and Next Steps

## Description

This topic shows where you continue after the PostgreSQL 16 and PostgreSQL 17 path. You learn the official tutorial and versioned docs, Postgres Weekly and the wiki, managed cloud products, and how you send a reproducible bug report.

Complete topic 21 first. You can describe a production boundary. After this topic, use a real database and keep reading.

Use one term for each concept. The official docs are the pages on postgresql.org for a major version. The wiki is a community site. A managed service runs PostgreSQL (or a compatible engine) for you. A reproducible bug report lets another person see the same defect. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## Read the official tutorial and current-version docs

The official tutorial is [https://www.postgresql.org/docs/current/tutorial.html](https://www.postgresql.org/docs/current/tutorial.html). The word `current` is the newest stable major version. That version can be newer than 17. Also open the tutorial that matches your server:

- [https://www.postgresql.org/docs/16/tutorial.html](https://www.postgresql.org/docs/16/tutorial.html)
- [https://www.postgresql.org/docs/17/tutorial.html](https://www.postgresql.org/docs/17/tutorial.html)

Complete the tutorial on your cluster even if you already finished this path. The tutorial uses official names and short examples.

The manuals that you will reopen:

| Book | Use |
| --- | --- |
| Tutorial | first path and review |
| SQL Language | types, functions, queries |
| Server Administration | install, config, backup, HA |
| Reference | exact syntax (`CREATE TABLE`, `SET`) |
| Internals | later, after you can operate a server |

Topic 1 listed these books. Now you have context for Administration and the Reference.

Release notes matter at every major upgrade:

- [https://www.postgresql.org/docs/17/release-17.html](https://www.postgresql.org/docs/17/release-17.html)
- [https://www.postgresql.org/docs/16/release-16.html](https://www.postgresql.org/docs/16/release-16.html)

Read the "Migration" section before you move data.

`psql` help stays useful:

```text
\h ALTER TABLE
\?
```

Do not learn production behavior from a random post when the official page exists. Do not bookmark only `docs/current` if production is 16. Do not skip release notes because the new major "looks the same".

### Questions

#### Theoretical questions

1. What does `docs/current` point to?
2. Why open `docs/16` when you run PostgreSQL 16?
3. Which book do you open for exact `CREATE SUBSCRIPTION` syntax?
4. Where do you read compatibility breaks for a major upgrade?
5. What does `\h ALTER TABLE` show?

#### Easy practical tasks

1. Open the 16 tutorial and the 17 tutorial. Write one heading that both share.
2. Complete or re-read one tutorial chapter on your server. Write one command that you ran.
3. Bookmark Administration and Reference for your major version.
4. Open the 17 migration notes. Write two items.

#### Medium practical tasks

1. Map topics 12 to 21 of this path to official chapter titles. Write a two-column table with eight rows.
2. Compare one Reference page in 16 and 17 (example: `COPY` or `CREATE PUBLICATION`). Write one difference or write that the page is the same.
3. Use `\h` for three commands from topic 16 or 21. Write the required arguments.

#### Advanced practical tasks

1. Read "Conventions" in the docs. Write how the docs mark optional syntax. Add two examples from `CREATE INDEX`.
2. Write a quarterly reading plan: tutorial review, release notes, one Administration chapter, one Internals page.

---

## Postgres Weekly / wiki

**Postgres Weekly** is a newsletter at [https://postgresweekly.com/](https://postgresweekly.com/). It links articles, extensions, and release news. Use it to hear about new major versions and tools. It is not a substitute for the official docs.

The **PostgreSQL wiki** is [https://wiki.postgresql.org/](https://wiki.postgresql.org/). The wiki has cookbooks, tuning notes, and conference material. Anyone can find outdated pages. When the wiki and the official docs disagree, the official docs win. Topic 1 stated that rule. Keep it.

Useful wiki habits:

- check the page date and the PostgreSQL version in the text
- follow links to postgresql.org when they exist
- treat bloat queries and `postgresql.conf` snippets as examples, not as a copy-paste policy

Other community sources (high-level): the mailing lists on postgresql.org, Planet PostgreSQL, and conference talks (PGConf). Prefer primary sources for security and upgrade steps.

Do not apply a wiki `shared_buffers` formula to a managed instance without a measurement (topic 18). Do not treat a newsletter link as a vendor support contract.

### Questions

#### Theoretical questions

1. What kind of content does Postgres Weekly collect?
2. What is the wiki URL?
3. Who wins when the wiki and the official docs disagree?
4. Why do you check the version on a wiki page?
5. Is a newsletter a substitute for release notes?

#### Easy practical tasks

1. Open Postgres Weekly. Write the title of the newest issue that the site shows.
2. Open the wiki home. Write three page titles that relate to topics 16 to 21.
3. Open one wiki tuning page. Write the PostgreSQL version that the page names, or write that it names none.
4. Write four sentences: newsletter, wiki, official docs, version check.

#### Medium practical tasks

1. Pick one wiki how-to (backup or vacuum). Compare it with the official chapter. Write two matches and one conflict or gap.
2. From one Weekly issue, list three links. Mark each as docs, blog, or product.
3. Find the wiki page about bug reporting or articles. Write how it points to official process.

#### Advanced practical tasks

1. Subscribe to Postgres Weekly or read three back issues. Write a one-page note: what is news versus what is already in this path.
2. Find a wiki page that still mentions `recovery.conf`. Write the official 16/17 replacement (topic 16).

---

## Cloud: RDS, Cloud SQL, AlloyDB, Aurora PostgreSQL, Crunchy, Supabase

Managed PostgreSQL reduces OS work. You still learn SQL, privileges, vacuum, and plans. The vendor hides some files (`postgresql.conf` may be a parameter group). Some extensions are blocked. Superuser is often replaced by a powerful admin role (topic 12).

Products that this path names (names and features change; read the vendor page):

| Product | Vendor | High-level note |
| --- | --- | --- |
| Amazon RDS for PostgreSQL | AWS | managed instance, parameter groups, Multi-AZ option |
| Amazon Aurora PostgreSQL | AWS | compatible engine, vendor storage and HA |
| Cloud SQL for PostgreSQL | Google Cloud | managed instance, HA option |
| AlloyDB | Google Cloud | PostgreSQL-compatible, vendor HA and features |
| Crunchy | Crunchy Data | PostgreSQL experts; Crunchy Bridge is a managed cloud |
| Supabase | Supabase | PostgreSQL plus an app platform (auth, APIs) |

Aurora and AlloyDB are compatible products. They are not a drop-in promise for every extension or every `postgresql.conf` knobs. Test `pg_dump` from the service and restore to community PostgreSQL 16 or 17 if you must leave.

What you still own on a managed service:

- schema, grants, RLS
- query plans and indexes
- migration files
- application pool size
- what you store in `jsonb`
- a backup policy even when the vendor takes snapshots
- a network rule (topic 21)

What the vendor often owns:

- OS patches
- disk and many HA mechanics
- some PITR buttons
- minor version rollout (you pick a window)

Do not assume `CREATE EXTENSION` works for every name. Do not assume `pg_basebackup` SSH access exists. Do not skip `pg_dump` tests because snapshots exist. Do not use a preview engine version as production without a policy.

### Questions

#### Theoretical questions

1. What work does a managed service remove?
2. What work do you still own?
3. Why are Aurora and AlloyDB not the same object as community PostgreSQL?
4. Why can `CREATE EXTENSION` fail on a cloud instance?
5. Why do you still test `pg_dump` when snapshots exist?

#### Easy practical tasks

1. Open the product page for two names in the table. Write the PostgreSQL major versions that each page lists today.
2. Open one parameter-group or flag list. Write two settings that you recognize from topic 18.
3. Find the backup or PITR page for one vendor. Write the retention that the page states.
4. Write four sentences: managed, compatible engine, extension list, network.

#### Medium practical tasks

1. Compare RDS PostgreSQL and Aurora PostgreSQL on failover and extensions (vendor docs). Write a six-sentence note.
2. Compare Cloud SQL and AlloyDB the same way, or compare Crunchy Bridge and Supabase if you prefer those pages.
3. List five `CREATE EXTENSION` names from this path. Mark each as likely allowed or unknown on one vendor (from their extension list).

#### Advanced practical tasks

1. Write an exit plan: `pg_dump` or logical replication from the service to a community 16 or 17 cluster. Name the gaps you must test (roles, extensions, sequences).
2. Read the shared-responsibility model on one vendor page. Make a two-column table: vendor, you. Add eight rows.

---

## Contribute a reproducible bug report

When PostgreSQL behaves in a way that the docs forbid, you can report it. The official guide is [https://www.postgresql.org/docs/current/bug-report.html](https://www.postgresql.org/docs/current/bug-report.html) (also under `/docs/16/` and `/docs/17/`).

A useful report includes:

- `SELECT version();` and the OS
- the exact SQL or steps
- the result that you got
- the result that you expected, with a docs link
- a small schema that another person can load (`CREATE TABLE`, `INSERT`)
- whether you can reproduce on community PostgreSQL 16 or 17, not only on a fork

```sql
SELECT version();
SHOW server_version_num;
```

Use the official channels that the guide names (form or mailing list). Search existing threads first.

A bug report is not:

- a question about how to write a query
- a vendor outage
- a feature request without a clear defect
- a dump of a production database

Managed forks can add defects that community PostgreSQL does not have. Reproduce on a local `postgres:16` or `postgres:17` image when you can. If you cannot, say that you only have the vendor engine.

Do not attach secrets, personal data, or a full customer dump. Do not demand a fix date. Do not skip the small reproducer.

Topic 1 asked for a five-line template in an advanced task. This section is the full habit.

### Questions

#### Theoretical questions

1. Where is the official bug-reporting guide?
2. Why must the report include `version()`?
3. What is a small schema in a report?
4. Why reproduce on community 16 or 17 when you can?
5. What must you leave out of the attachment?

#### Easy practical tasks

1. Open the 17 bug-report page. Write the channel that the page names.
2. Write a five-line template: version, steps, got, expected, docs URL.
3. Run `SELECT version();` and save the full string.
4. Search the pgsql-bugs list or the form archive for one recent item. Write the subject.

#### Medium practical tasks

1. Write a fake but complete report for a made-up `EXPLAIN` defect. Include `CREATE TABLE` and three `INSERT`s. Do not send it.
2. Compare a good report and a poor report (question only). Write three differences.
3. Document how you would redact a report that started from production SQL.

#### Advanced practical tasks

1. Reproduce a known fixed bug from the release notes on an old lab if you have one, or write why you cannot. Practice the template.
2. Read "Guided Tour" or mailing-list etiquette linked from the docs. Write how you reply if a developer asks for `EXPLAIN (ANALYZE, BUFFERS)`.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the official tutorial, versioned manuals, and release notes work together for a 16-to-17 year?
2. When do you trust Postgres Weekly or the wiki, and when do you stop?
3. What must you still practice on a managed cloud PostgreSQL 16 or 17 service?
4. How does a reproducible report differ from a production incident ticket?
5. A teammate learns only from social media and runs an untested fork in production. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. Bookmark tutorial, docs for 16, docs for 17, `docs/current`, Postgres Weekly, the wiki, and the bug-report page.
2. Write a cheat sheet: five official books, wiki rule, six cloud names, bug-report fields.
3. Run `SELECT version();` and open the matching docs URL.
4. Write one sentence each for RDS, Cloud SQL, AlloyDB, Aurora PostgreSQL, Crunchy, and Supabase.

#### Medium practical tasks

1. Pick the next real project (topics 12 to 21). Write which official chapters and which vendor limits you will read first.
2. Compare one cloud backup button with topic 16 (`pg_dump`, PITR). Write what you still test by hand.
3. Draft a team page: where docs live, which major versions you support, how you file a bug.

#### Advanced practical tasks

1. Complete the official tutorial end to end on PostgreSQL 16 or 17. Add a one-page map from each tutorial chapter to this LearnGO path.
2. Prepare a reproducible report template as a file in your notes (not sent). Include `version()`, a schema, and a redaction checklist. Then write the next feature you will learn from the Internals book or from a conference talk.
