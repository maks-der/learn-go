# 14. Storage and Internals

## Description

This topic shows how a DBMS stores data on disk and in memory. You learn pages, heap tables, clustered and index-organized tables, the write-ahead log (WAL), the buffer cache, checkpoints, and vacuum (garbage collection).

Use one term for each concept. A page is a fixed-size block of bytes. The buffer cache is the set of pages in memory. The WAL is the durable log of changes. Complete this topic after you understand transactions and indexes. Internals explain why commit, crash recovery, and vacuum exist.

This path is vendor-neutral. Product names for the same idea differ (WAL, redo, journal). The jobs do not.

---

## Pages / blocks

A page (also called a block) is the unit that the DBMS reads from disk and writes to disk. Typical sizes are 4 KiB, 8 KiB, or 16 KiB. The product sets the size. You do not store one row as one file.

A table is a set of pages. Each page holds a header and a number of row versions (or row slots). A large row can span more than one page. A small row shares a page with other rows.

The DBMS addresses a row with a page identifier plus a slot (or an equivalent row pointer). An index leaf stores that pointer, or it stores the row in a clustered structure (next section).

```text
Page 0   : header + rows
Page 1   : header + rows
Page 2   : header + rows
...
```

I/O cost is in pages, not in rows. A query that needs one row still reads at least one page. A sequential scan reads many pages in order. Random reads jump between pages. Sequential reads are often cheaper on the same disk.

When you update a row, the DBMS often writes a new version on a page and leaves the old version until vacuum or until undo no longer needs it. That detail differs by product. The page is still the container.

Do not confuse a page with a filesystem block, though the sizes can match. The DBMS page is the structure that the storage engine understands. The operating system may split or merge I/O around it.

Choose a page size only when the product allows it and when you have a measured reason. For learning, keep the default.

### Questions

#### Theoretical questions

1. What is a page in a DBMS?
2. Why can one I/O still occur when a query needs one row?
3. How does an index usually point to a heap row?
4. Why is a sequential page read often cheaper than a random page read?
5. Is a DBMS page the same object as a filesystem block? Explain.

#### Easy practical tasks

1. Find the default page size of your DBMS in the official docs. Write the size.
2. Draw three pages and six small rows. Show two rows on the same page.
3. Write three sentences: page, row pointer, table as a set of pages.
4. Estimate pages for 10 000 rows of 200 bytes if the page is 8 KiB (ignore headers). Write the estimate.

#### Medium practical tasks

1. Use a vendor tool or query that shows table size in pages or bytes. Record the number for a practice table.
2. Explain why a 2 KiB row and an 8 KiB page still waste space if many pages are half empty.
3. Compare a full scan (many sequential pages) with a primary-key lookup (few pages). Write the I/O difference in words.

#### Advanced practical tasks

1. Write a one-page note: how page size affects row density, TOAST/overflow rows, and I/O. Use your product terms.
2. Fill a table until it occupies many pages. Record size before and after a large delete (before vacuum). Write what you expect from later sections.

---

## Heap vs clustered / index-organized tables (vendor differences)

A heap table stores rows in pages without a required order by primary key. A new row goes into a page with free space. The primary key is a separate unique index. The index points to the heap location.

A clustered table (or index-organized table) stores the row with the primary-key order. The primary key is the table organization. Secondary indexes point to the primary key (or to a clustering key), not to a stable heap slot.

| Organization | Row order | Primary key | Secondary index target |
| --- | --- | --- | --- |
| Heap | No key order | Separate unique index | Page/slot (or equivalent) |
| Clustered / IOT | Primary-key order | The table itself | Usually the primary key |

Product examples (category only): many PostgreSQL user tables are heaps. Some other products cluster by primary key by default. The names "clustered index" and "index-organized table" come from those products. Learn the idea. Then read your product.

Heap benefits: inserts are simple. Updates can leave a new version in a new slot. Vacuum must reclaim old versions.

Clustered benefits: a range scan on the primary key reads rows in order. A point lookup on the primary key can avoid a second heap fetch.

Clustered costs: a secondary index lookup often needs a second step to the primary key. A change of the clustering key can move the row.

Do not assume that `ORDER BY primary_key` is free on a heap. Without an index, the DBMS still sorts. On a clustered table, a primary-key range can match the storage order.

Do not pick a product only because of this difference. Both designs work. Your indexes and queries matter more for a beginner.

### Questions

#### Theoretical questions

1. Where does a new row go in a heap table?
2. What does a clustered or index-organized table use as its organization?
3. What does a secondary index point to in each design?
4. Why can a primary-key range scan be cheaper on a clustered table?
5. Why is `ORDER BY` on a heap primary key not automatically free?

#### Easy practical tasks

1. Label your DBMS as heap-default or cluster-default from the official docs.
2. Draw a heap page plus a separate primary-key index. Draw a clustered leaf that holds the row.
3. Write one insert story for each design in four sentences total.
4. Name one benefit and one cost of clustering.

#### Medium practical tasks

1. Create a table with a primary key. Insert rows in random key order. Read a key range. Write whether you expect storage order to match key order.
2. Add a secondary index. Write the lookup steps for heap versus clustered (two short sequences).
3. Find the product term: heap, clustered index, or index-organized table. Quote the manual heading in your notes.

#### Advanced practical tasks

1. Write a one-page comparison for a team that moves a schema between a heap product and a clustered product. Include secondary indexes.
2. Measure a primary-key range query on a large practice table. Record the plan. Relate the plan to heap or cluster.

---

## WAL / redo log idea

WAL means write-ahead log. The DBMS writes the change to a log on stable storage before it considers the change durable. Some products say redo log or journal. The job is the same: recover after a crash.

A transaction changes pages in memory. Those pages may stay dirty (not yet written to the data files). On `COMMIT`, the DBMS must not lose the committed change. It flushes the WAL records for that transaction. It does not need to flush every dirty data page at commit time.

```text
1. Change pages in the buffer cache
2. Append WAL records
3. Flush WAL to disk on COMMIT
4. Later: write dirty pages to data files
```

Recovery uses the WAL. After a crash, the DBMS replays log records that did not reach the data files. It undoes uncommitted work if the design requires that. You get durability of committed transactions and atomicity of incomplete ones.

WAL also feeds replicas and backup tools in many products. Physical backups and point-in-time recovery depend on a continuous log. A later topic covers that use.

Synchronous_commit (product name differs) controls whether `COMMIT` waits for the log flush. If you turn that wait off, you can lose the last commits after a crash. That setting is a durability trade-off, not a free speed gain.

Do not delete WAL files by hand. The DBMS recycles them after a checkpoint and after replicas or archive jobs no longer need them.

Do not treat a copy of the data directory without the matching log as a safe backup.

### Questions

#### Theoretical questions

1. What does write-ahead mean?
2. Why can `COMMIT` skip a flush of every dirty data page?
3. What does recovery replay after a crash?
4. Why does a replica often need the WAL?
5. What do you lose if `COMMIT` does not wait for a log flush?

#### Easy practical tasks

1. Find the product name for the WAL (WAL, redo, journal). Write the name.
2. Number the four steps in this section. Recite them without the file.
3. Write one sentence that distinguishes WAL files from table data files.
4. Write why you must not delete WAL files in a file manager.

#### Medium practical tasks

1. Locate the WAL or redo directory in a local install (docs). Write the path pattern. Do not delete files.
2. Read the setting that controls whether commit waits for flush. Write the default and the risk of the faster option.
3. Draw crash recovery: dirty page in memory, WAL on disk, data file old. Show the replay.

#### Advanced practical tasks

1. Write a one-page durability note: commit, WAL flush, data-page write, crash. Use your product terms.
2. In a disposable instance, force a shutdown mode that the docs call unsafe (only if the docs describe a lab). Or write a table-top of two shutdown types. Record what WAL must do.

---

## Buffer cache / shared buffers

The buffer cache (shared buffers) is the set of pages that the DBMS keeps in memory. A cache hit means the page is already in memory. A cache miss means the DBMS reads the page from disk (or from the operating-system cache).

All backends of a server DBMS share one cache (product details differ). The cache has a replacement policy. A common idea is to keep hot pages and evict cold pages.

```text
Query needs page P
  P in cache  -->  read memory
  P not in cache  -->  read disk, put P in cache, maybe evict Q
```

Dirty pages sit in the cache until a checkpoint or a background writer writes them. The WAL still protects committed data.

A larger cache can reduce disk reads. A cache that is too large can starve the operating system or other processes. The operating system also caches files. Some products rely on that OS cache more than others.

Do not set the buffer cache to all of RAM. Leave memory for connections, sorts, the OS, and the application.

Do not use the cache as an application cache of query results. The buffer cache holds pages. A result cache is a different layer (a later performance topic).

A full table scan of a table larger than the cache will evict useful pages. That effect is cache pollution. Large analytics scans and small OLTP lookups compete in one cache unless you separate systems.

### Questions

#### Theoretical questions

1. What does the buffer cache store?
2. What is a cache hit?
3. Why can a dirty page stay in memory after `COMMIT`?
4. Why must the buffer cache not consume all RAM?
5. How does a large scan pollute the cache?

#### Easy practical tasks

1. Find the buffer-cache setting name in your DBMS. Write the default.
2. Draw a hit and a miss for page P.
3. Write three memory consumers besides the buffer cache.
4. Write one sentence that distinguishes buffer cache from a query-result cache.

#### Medium practical tasks

1. Read current cache size and, if available, a hit-ratio metric. Write the numbers and what they do not prove.
2. Run a large sequential scan and a small lookup (practice data). Write which workload you expect to evict more pages.
3. Compare DBMS cache and OS file cache in three sentences from the product docs.

#### Advanced practical tasks

1. Write a one-page memory budget for a small server: RAM, buffers, connections, OS. Use rough percents.
2. After you have metrics, propose a cache-size change with a before/after measure plan. Do not change production without that plan.

---

## Checkpoints

A checkpoint is a point when the DBMS writes dirty pages so that recovery can start from a known place. After a checkpoint, older WAL is not required for crash recovery of those flushed pages (replicas and archive may still need that WAL).

The DBMS runs checkpoints on a schedule, on a WAL-volume trigger, or when you request one. A checkpoint is not the same as `COMMIT`. Many commits occur between checkpoints.

```text
WAL  ----c--------c--------c---->
         ^        ^        ^
      checkpoint
```

A long time between checkpoints means a longer replay after a crash. A very frequent checkpoint means more write I/O in normal operation.

During a checkpoint, write traffic to data files rises. Users can see a pause if the product writes in a burst. Many products spread the write (checkpoint completion target).

Do not run a manual checkpoint as a substitute for `COMMIT`. Do not run frequent manual checkpoints to "make it safer" without a measure. The WAL already makes commit durable.

Backup tools often create or wait for a checkpoint so that the data files and the log line up. Follow the backup procedure of the product.

If recovery is slow after a crash, one cause is a large distance from the last checkpoint to the crash. Another cause is a large WAL to replay. Measure. Then change the checkpoint interval with a written reason.

### Questions

#### Theoretical questions

1. What does a checkpoint write?
2. How does a checkpoint differ from `COMMIT`?
3. What is the recovery cost of a long gap between checkpoints?
4. What is the run-time cost of a very frequent checkpoint?
5. Why is a manual checkpoint not a substitute for `COMMIT`?

#### Easy practical tasks

1. Find the checkpoint interval settings in your DBMS docs. Write the names.
2. Draw WAL with three checkpoints and several commits between them.
3. Write two reasons a backup tool waits for a checkpoint.
4. Write one user-visible symptom of a heavy checkpoint.

#### Medium practical tasks

1. Read a server log or a metric for the last checkpoint time if the product exposes it. Record what you found.
2. Relate checkpoint frequency to WAL disk use in five sentences.
3. In a disposable instance, trigger a checkpoint with the official command. Confirm in the log.

#### Advanced practical tasks

1. Write a one-page note: checkpoint, WAL recycle, crash recovery time, and backup. Use your product terms.
2. Propose a checkpoint-tuning experiment with one metric for recovery time and one metric for write latency.

---

## Vacuum / garbage collection idea (vendor-specific later)

Many DBMS designs keep old row versions until no transaction can see them. Update and delete do not always free the space at once. Vacuum (garbage collection, purge, undo cleanup) reclaims that space and maintains related structures.

In a heap with versioning (example category: MVCC heaps), an `UPDATE` makes a new row version. The old version stays until vacuum. Dead versions increase table size (bloat). A scan reads more pages than live data requires.

Vacuum also updates visibility information and can freeze old transaction identifiers (product-specific). Autovacuum or an equivalent background job runs this work. You can also run a manual vacuum in maintenance windows.

```text
UPDATE row  -->  new version + old version
COMMIT
...
VACUUM     -->  reclaim old version if no reader needs it
```

Other products use undo logs. The garbage-collection idea still exists: something must remove data that no transaction needs.

Vacuum is not a backup. Vacuum is not `TRUNCATE`. `VACUUM FULL` or a rebuild (names differ) can compact the table and hold a strong lock. Read the lock rules before you run a full rewrite on production.

Do not turn off autovacuum (or the equivalent) to "go faster" unless you replace it with a tested job. Dead versions accumulate. Queries slow down. Disk grows.

A large delete without vacuum leaves dead space. A later insert can reuse that space after vacuum. Until then, the table can stay large.

Learn the product chapter after this path. The idea is enough here: writes create garbage; a cleaner must run; you watch bloat and cleaner health.

### Questions

#### Theoretical questions

1. Why does an `UPDATE` not always free space immediately?
2. What is table bloat in this section?
3. What job does vacuum or garbage collection do?
4. Why is vacuum not a backup?
5. What happens if the automatic cleaner is off for a long time?

#### Easy practical tasks

1. Find the product name: `VACUUM`, purge, undo cleanup, or another term.
2. Draw live version, dead version, and vacuum.
3. Write three statements: `UPDATE`, `DELETE`, `TRUNCATE`. Mark which ones leave versions that vacuum must handle (in an MVCC heap).
4. Write one risk of `VACUUM FULL` or a table rewrite on a busy system.

#### Medium practical tasks

1. Update many rows on a scratch table. Measure table size. Run the official vacuum. Measure again. Write the two sizes.
2. Read autovacuum (or equivalent) settings. Write how you know the cleaner ran (log or statistic view).
3. Explain in six sentences why a 1 GB table can stay 1 GB after you delete 80 percent of rows, until vacuum or a rewrite.

#### Advanced practical tasks

1. Write a one-page operations note: how you monitor bloat, cleaner lag, and when you schedule a rewrite.
2. Compare MVCC-heap vacuum with undo-log cleanup from two manuals. Write five sentences on what a DBA still watches in both.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do pages, the buffer cache, the WAL, and a checkpoint work together from `UPDATE` to a crash restart?
2. When does table organization (heap versus clustered) change the meaning of a secondary index pointer?
3. Why can a system be durable at `COMMIT` and still need vacuum for space and scan cost?
4. Which I/O is sequential and which I/O is random in: a heap scan, a WAL append, a checkpoint write, an index seek?
5. A teammate copies only table files, deletes WAL, and turns off the cleaner to save CPU. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: page, heap versus cluster, WAL, buffers, checkpoint, vacuum.
2. On a local DBMS, write the official names for WAL, shared buffers, checkpoint, and vacuum.
3. Insert, update, and delete one row. Write which component logged the change and which component may still hold a dirty page.
4. Draw a single timeline: commit, checkpoint, crash, replay, vacuum.

#### Medium practical tasks

1. Create a practice table. Record page size, table size, and buffer setting. Run a scan and a primary-key lookup. Relate each to pages.
2. Write a lab plan (do not damage a shared server): fill, update, checkpoint, vacuum, measure size and a plan.
3. Document the data directory versus the WAL directory for your install. State what a crash-safe copy must include.

#### Advanced practical tasks

1. Write a recovery-and-storage runbook of one page: crash restart, WAL, checkpoint gap, buffer warmup, vacuum after a bulk load.
2. Compare two products on heap versus cluster and on WAL naming. Use official docs. Write a table of terms and one operational difference.
