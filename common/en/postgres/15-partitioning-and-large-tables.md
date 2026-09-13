# 15. Partitioning and Large Tables

## Description

This topic shows declarative partitioning in PostgreSQL 16 and PostgreSQL 17. You learn range, list, and hash partitions, partition pruning, indexes on partitioned tables, why inheritance is the old path, and how you archive old partitions.

Complete topic 14 first. You can model wide rows. Complete this topic before you study backup and replicas.

Use one term for each concept. A partitioned table is the parent. A partition is a child table that holds a slice of the rows. The partition key is the column list that chooses the slice. Partition pruning skips partitions that cannot match the query. Declarative partitioning is `PARTITION BY` on `CREATE TABLE`. Inheritance is a different, older feature. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## Declarative partitioning (range, list, hash)

Declarative partitioning splits one table into partitions. You `INSERT` into the parent. The server routes the row to one partition. A `SELECT` on the parent can read many partitions.

Three methods:

- **RANGE** — contiguous bounds (`FOR VALUES FROM ... TO ...`). Usual for `date` and `timestamptz`.
- **LIST** — discrete values (`FOR VALUES IN (...)`). Usual for a tenant id or a status that you split on purpose.
- **HASH** — modulus buckets (`FOR VALUES WITH (MODULUS n, REMAINDER r)`). Usual when you need even size and you always filter by the hash key.

```sql
CREATE TABLE shop.orders (
    id bigint GENERATED ALWAYS AS IDENTITY,
    created_on date NOT NULL,
    tenant_id bigint NOT NULL,
    total numeric(12, 2) NOT NULL,
    PRIMARY KEY (id, created_on)
) PARTITION BY RANGE (created_on);

CREATE TABLE shop.orders_2026_01
    PARTITION OF shop.orders
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

CREATE TABLE shop.orders_2026_02
    PARTITION OF shop.orders
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
```

The upper bound is exclusive. A row with `created_on = '2026-02-01'` goes to `orders_2026_02`.

A default partition catches values that no other partition claims:

```sql
CREATE TABLE shop.orders_default
    PARTITION OF shop.orders DEFAULT;
```

Use a default partition with care. It can grow without limit if you forget new months.

Hash example:

```sql
CREATE TABLE shop.sessions (
    sid uuid NOT NULL,
    user_id bigint NOT NULL,
    PRIMARY KEY (sid)
) PARTITION BY HASH (sid);

CREATE TABLE shop.sessions_p0 PARTITION OF shop.sessions
    FOR VALUES WITH (MODULUS 4, REMAINDER 0);
-- p1, p2, p3 in the same way
```

Unique keys and primary keys on a partitioned table must include the partition key. That is why `PRIMARY KEY (id, created_on)` appears above. A global unique `id` alone is not valid on the parent in PostgreSQL 16 and 17.

Do not partition a small table. Partitioning adds planning cost. Do not mix range months with unrelated list keys on the same parent without a documented design. Do not insert a row that matches no partition and has no default; the statement fails.

Official chapter: [https://www.postgresql.org/docs/17/ddl-partitioning.html](https://www.postgresql.org/docs/17/ddl-partitioning.html).

### Questions

#### Theoretical questions

1. What is the difference between RANGE, LIST, and HASH?
2. Is the RANGE upper bound inclusive or exclusive?
3. Why must a primary key on a partitioned table include the partition key?
4. What does a default partition hold?
5. What happens when no partition matches and there is no default?

#### Easy practical tasks

1. Create a range-partitioned table with two monthly partitions. Insert one row in each month.
2. Insert a date that has no partition and no default. Record the error.
3. Add a default partition. Insert that date again. Select from the default table.
4. Run `\d+ shop.orders` and list partition names.

#### Medium practical tasks

1. Create a list-partitioned table on `tenant_id` with two tenants. Insert and select from the parent.
2. Create a hash-partitioned table with modulus 4. Insert 20 rows. Count rows per partition.
3. Try `PRIMARY KEY (id)` on a table that is `PARTITION BY RANGE (created_on)`. Record the error. Fix the key.

#### Advanced practical tasks

1. Design keys for "orders by month" and "sessions by hash of sid". Write `CREATE TABLE` for both. Include primary keys that PostgreSQL accepts.
2. Read "Partitioning" in the 17 docs. Write the rules for `NULL` in RANGE and LIST. Test one `NULL` insert.

---

## Partition pruning

Partition pruning skips partitions that cannot contain matching rows. The planner needs a constraint on the partition key (or a value that it can prove).

```sql
EXPLAIN (ANALYZE, COSTS OFF)
SELECT id, total
FROM shop.orders
WHERE created_on >= DATE '2026-02-01'
  AND created_on <  DATE '2026-03-01';
```

The plan must list `orders_2026_02` and must not scan `orders_2026_01` when the bounds are constants.

Pruning can occur at plan time or at run time. Run-time pruning applies when the value is a parameter that the executor learns later (`PREPARE` / bind, or a join). `EXPLAIN (ANALYZE)` shows removed partitions more clearly than `EXPLAIN` without analyze.

A filter on a column that is not the partition key does not prune. `WHERE total > 100` still can scan every month.

`SET enable_partition_pruning = off` is a debug switch. Keep it on in production.

Functions on the partition key can block pruning. Prefer `created_on >= $1` over `date_trunc('month', created_on) = $1` when `created_on` is the key.

Do not hide the partition key inside an expression that the planner cannot prove. Do not expect pruning from a `jsonb` field unless that field is the partition key.

### Questions

#### Theoretical questions

1. What does partition pruning skip?
2. Which `WHERE` clause prunes a range on `created_on`?
3. Why can a parameter still prune at run time?
4. Does `WHERE total > 100` prune month partitions?
5. What debug setting disables pruning?

#### Easy practical tasks

1. `EXPLAIN` a query that filters one month. Write which partitions appear.
2. `EXPLAIN` the same table with no date filter. Write how many partitions appear.
3. `SHOW enable_partition_pruning;`.
4. Write four sentences: prune, partition key, constant, run time.

#### Medium practical tasks

1. `EXPLAIN (ANALYZE)` a prepared statement with a date parameter. Find run-time pruning in the text if it appears.
2. Compare `WHERE created_on >= $d AND created_on < $d2` with `WHERE date_trunc('month', created_on) = $d`. Write which plan prunes.
3. Turn `enable_partition_pruning` off in a session. Repeat `EXPLAIN`. Then set it on. Do not leave it off.

#### Advanced practical tasks

1. Join `orders` to `tenants` with a filter only on `tenants`. Write whether orders partitions prune. Add a date filter and compare.
2. Read "Partition Pruning" in the 16 or 17 docs. Write one sentence about pruning during execution versus planning.

---

## Indexes on partitioned tables

An index on the parent is a partitioned index. PostgreSQL creates a matching index on each partition. New partitions get the index when you attach them if you created the index on the parent.

```sql
CREATE INDEX orders_tenant_idx ON shop.orders (tenant_id);
```

Each partition has its own B-tree. There is no single cluster-wide B-tree for all months. A unique index on the parent must include the partition key, the same rule as the primary key.

You can create an index on one partition only:

```sql
CREATE INDEX orders_2026_01_total_idx ON shop.orders_2026_01 (total);
```

Use a local index when one slice needs a different shape. Prefer parent indexes so that you do not forget new months.

`CREATE INDEX CONCURRENTLY` on a partitioned parent has extra rules in PostgreSQL 16 and 17. You often create the index on the parent without `CONCURRENTLY`, or you build per partition concurrently and attach. Read the current `CREATE INDEX` page before you run this on production.

`REINDEX` can target the parent or one partition. Topic 8 covered `REINDEX` and `CONCURRENTLY`.

Primary keys, unique constraints, and foreign keys have limits across partitions. A foreign key from a normal table to a partitioned table is supported when the referenced unique key exists. A foreign key from a partitioned table to another table is also common. Cross-partition unique keys that omit the partition key are not available.

Do not expect one unique `email` across all list partitions if `email` is not in the partition key. Do not create fifty different index definitions by hand without a parent index.

### Questions

#### Theoretical questions

1. What does `CREATE INDEX` on the parent create on each partition?
2. Must a unique parent index include the partition key?
3. When do you index one partition only?
4. Is there one B-tree for all range partitions?
5. Why can a unique `email` fail as a parent constraint?

#### Easy practical tasks

1. Create an index on the parent. Run `\d` on the parent and on one partition.
2. Create an extra index on one partition only. List indexes on that partition.
3. Query `pg_indexes` where `tablename` like `orders%`.
4. Write four sentences: parent index, local index, unique key, partition key.

#### Medium practical tasks

1. Attach a new month partition. Confirm that the parent index appears on the new table.
2. Try a unique index on `(id)` only. Record the error.
3. `EXPLAIN` a query that filters `tenant_id` and one month. Write which index each node uses if the planner shows it.

#### Advanced practical tasks

1. Read `CREATE INDEX` on partitioned tables in the 17 docs. Write the concurrent-build options that the page allows.
2. Plan foreign keys: `order_lines` to `orders`. Write a key that works with `PRIMARY KEY (id, created_on)`. Include `created_on` on the child if required.

---

## Inheritance vs declarative partitioning

Table inheritance lets a child table receive columns from a parent. Topic 5 said to prefer composition. Inheritance also served as the old manual partition method: children with `CHECK` constraints, `INSERT` triggers, and `constraint_exclusion`.

Declarative partitioning is the supported method for new large tables in PostgreSQL 16 and 17.

| Task | Inheritance | Declarative |
| --- | --- | --- |
| Route `INSERT` on the parent | trigger or app | built in |
| Prune unused children | `constraint_exclusion` | partition pruning |
| Unique key across slices | easy to get wrong | key must include partition key |
| Attach / detach | `INHERIT` / `NO INHERIT` | `ATTACH` / `DETACH` |
| New work | do not use for shards | use this |

```sql
-- old pattern (do not use for new shards)
CREATE TABLE logs_all (t timestamptz, msg text);
CREATE TABLE logs_2026 () INHERITS (logs_all);
```

`SELECT` from the parent can include children. `SELECT ... ONLY` skips children. That behavior is inheritance, not declarative partitioning.

A declarative parent is not a normal heap for your rows. Rows live in partitions. Inheritance parents can hold rows and children can hold rows. That mix surprises beginners.

Know inheritance so that you can read old schemas. Migrate old inherit-and-trigger shards to `PARTITION BY` when you can plan the move.

Do not create a new inherit tree to split a fact table. Do not mix `INHERITS` and `PARTITION BY` on the same table.

### Questions

#### Theoretical questions

1. What feature do you use for new partitioned tables?
2. How did old designs route inserts into inherit children?
3. What does `ONLY` do on a `SELECT` from an inherit parent?
4. Can a declarative parent hold rows in its own heap the same way?
5. Why does topic 5 tell you to prefer composition over inheritance for columns?

#### Easy practical tasks

1. Draw the table from this section. Add one row of your own (task, inherit, declarative).
2. Create a tiny inherit parent and child (lab). Insert into the child. Select from the parent and from `ONLY` the parent.
3. Compare `\d+` of an inherit parent with `\d+` of a partitioned parent.
4. Write four sentences: inherit, `ONLY`, declarative, prune.

#### Medium practical tasks

1. Read "Table Inheritance" and "Declarative Partitioning" in the docs. Write three incompatibilities.
2. Try `PARTITION BY` and `INHERITS` together if you are unsure. Record the error. If the server accepts a form, write why you still avoid it.
3. List production risks of trigger-based routing (missed trigger, wrong `CHECK`, double insert).

#### Advanced practical tasks

1. Write a migration sketch from an inherit-and-check layout to `PARTITION BY RANGE`. Do not run it on real data. Name `ATTACH PARTITION`.
2. Find `constraint_exclusion` in the docs. Write why declarative pruning replaced it for new work.

---

## Archiving old partitions

A partition is a table. You can detach it, dump it, move it, or drop it. That is the main operations win of range partitioning.

Detach without a long wait on PostgreSQL 14 and later (16 and 17 included):

```sql
ALTER TABLE shop.orders DETACH PARTITION shop.orders_2025_01 CONCURRENTLY;
```

`CONCURRENTLY` cannot run in a transaction block. After a successful detach, the table is a normal table. You can `pg_dump` it, move it to an archive schema, or `DROP TABLE` when the backup is safe.

Attach an old table that already has the right rows and bounds:

```sql
ALTER TABLE shop.orders ATTACH PARTITION shop.orders_2025_01
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
```

Attach checks the bounds. The check can take a lock and scan. Prepare the table with a matching `CHECK` constraint to speed attach when the docs say that helps.

Archive flow that many teams use:

1. Create next month partitions before the month starts.
2. Detach months that the app must not query online.
3. Dump or copy the detached table to cold storage.
4. Drop the detached table after you verify the dump.

A default partition that collected "forgotten" months is harder to split. Keep a calendar job that creates partitions.

Do not `DELETE` millions of old rows from a non-partitioned table if you can partition by date and drop a slice. Do not drop a partition before a backup that you tested. Topic 16 covers `pg_dump` of one table.

### Questions

#### Theoretical questions

1. What does `DETACH PARTITION` do to the child table?
2. Why use `DETACH ... CONCURRENTLY`?
3. Can `DETACH ... CONCURRENTLY` run inside `BEGIN`?
4. What must be true before you `DROP` a detached partition?
5. Why is a monthly partition easier to archive than a huge `DELETE`?

#### Easy practical tasks

1. Detach one lab partition (without `CONCURRENTLY` is acceptable). Select from the parent and from the detached table.
2. Attach it again with the same bounds. Select from the parent.
3. Write the four-step archive flow in your own words.
4. List partitions with `\d+ shop.orders`.

#### Medium practical tasks

1. `DETACH ... CONCURRENTLY` on a quiet lab table. Record whether you used a transaction block and what happened.
2. `pg_dump` a detached table (topic 16 preview: `-t shop.orders_2025_01`). Restore it under a new name.
3. Create next month's partition with a script that uses `generate_series` of dates. Show the `CREATE TABLE` output.

#### Advanced practical tasks

1. Write a runbook: create future partitions, detach old ones, dump, verify row counts, drop. Include a check that the app no longer needs that month.
2. Read `ATTACH PARTITION` locking and constraint notes in the 17 docs. Write how a matching `CHECK` speeds attach.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do RANGE bounds, pruning, and detach form one lifecycle for a time-series table?
2. Why do unique indexes and foreign keys become harder after you partition?
3. When is HASH a better method than RANGE?
4. How do you explain inheritance to a teammate who copies an old blog about "table partitioning"?
5. A teammate wants 200 list partitions on a low-cardinality status column that every query ignores. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. Build a two-month `orders` parent. Insert rows. `EXPLAIN` a one-month query. Detach the older month.
2. Write a cheat sheet: RANGE, LIST, HASH, prune, parent index, inherit, detach.
3. Query `pg_partition_tree('shop.orders');` if the function exists in your version. Save the result.
4. Insert one row and select it from the parent and from the partition table name.

#### Medium practical tasks

1. Add a parent index on `tenant_id`. Add a third month. Prove the index exists on the new partition. Run a pruned query.
2. Compare `pg_relation_size` of each partition after an uneven load.
3. Document a calendar: when you create next partitions and when you detach.

#### Advanced practical tasks

1. Load 100000 rows across four months. Compare `EXPLAIN ANALYZE` for a month filter versus a filter that cannot prune. Write buffer counts if you use `BUFFERS`.
2. Map this topic to the official partitioning chapter. Add one limit (foreign keys, identity, or unique) that this topic did not detail.
