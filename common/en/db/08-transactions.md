# 8. Transactions

## Description

This topic shows how the DBMS keeps writes safe when more than one session works at the same time. You learn ACID, `BEGIN` / `COMMIT` / `ROLLBACK`, autocommit, isolation anomalies, isolation levels, locks, deadlocks, and retry.

Use one term for each concept. A session is one client connection. A transaction is a unit of work in a session. Isolation is the rule that limits what concurrent transactions see. Complete this topic before you normalize a schema. You need transactions when you split tables and still write one business event.

Practice with two CLI sessions when you can. One session is not enough to see interleaving.

---

## ACID

ACID is a set of four properties for transactions.

**Atomicity.** All writes in the transaction happen, or none happen. A failure in the middle does not leave a half-finished order.

**Consistency.** A committed transaction leaves the database in a state that satisfies all constraints. A transaction that would break a foreign key does not commit. Consistency here means "rules of the schema," not "the business is always correct."

**Isolation.** Concurrent transactions do not see each other's uncommitted writes (at typical levels). The DBMS makes it appear that transactions run in some serial order, to a degree that the isolation level defines. Weaker levels allow more interleaving.

**Durability.** After `COMMIT` succeeds, the committed data survives a crash. The DBMS writes a log to stable storage. Commit means the change is persistent.

Example. Transfer 50 from account A to account B:

1. Subtract 50 from A.
2. Add 50 to B.
3. Commit.

Atomicity keeps the pair together. Consistency keeps balances inside `CHECK` rules if you have them. Isolation stops a second session from reading A after step 1 and B before step 2 as if that were a committed state (at a sufficient isolation level). Durability keeps the result after a power loss.

Do not treat ACID as a slogan. Name the property that you need. "We need ACID" is vague. "This transfer must be atomic and durable" is clear.

Do not assume that every DBMS and every setting gives the same isolation. Read the default isolation level of your product.

### Questions

#### Theoretical questions

1. What does atomicity guarantee?
2. What does consistency mean in ACID?
3. What does isolation limit?
4. What does durability guarantee after `COMMIT`?
5. Why is "we need ACID" a weak requirement statement?

#### Easy practical tasks

1. Write one sentence per ACID letter with the money-transfer example.
2. Map a failed second insert of an order line to atomicity.
3. Map a foreign-key rejection to consistency.
4. Make a four-row table: property, meaning, failure if missing.

#### Medium practical tasks

1. Run a two-statement transfer in a transaction and roll back. Show atomicity.
2. Read the default isolation level in your DBMS docs. Write the name.
3. Unplug the idea of durability: write what must exist on disk after commit (high-level: log).

#### Advanced practical tasks

1. Write a one-page note that separates schema consistency from business correctness. Give one example of a committed but "wrong" sale.
2. Compare ACID claims of a relational DBMS and a typical key-value cache. Write six sentences.

---

## `BEGIN` / `COMMIT` / `ROLLBACK`

You start an explicit transaction with `BEGIN` (or `START TRANSACTION` on some products). You end it with `COMMIT` or `ROLLBACK`.

```sql
BEGIN;
UPDATE accounts SET balance = balance - 50 WHERE account_id = 1;
UPDATE accounts SET balance = balance + 50 WHERE account_id = 2;
COMMIT;
```

```sql
BEGIN;
UPDATE accounts SET balance = balance - 50 WHERE account_id = 1;
ROLLBACK;
```

After `ROLLBACK`, the subtract does not remain.

`COMMIT` makes the changes visible to other sessions (subject to isolation) and durable. `ROLLBACK` cancels all changes of this transaction. The table definitions that you created inside the transaction may also roll back, depending on the product. Test DDL. Some products auto-commit DDL.

If the client disconnects before `COMMIT`, the DBMS rolls back the open transaction in the usual case.

Savepoints let you roll back part of a transaction. Syntax: `SAVEPOINT name`, then `ROLLBACK TO SAVEPOINT name`. Use a savepoint when one optional step can fail and you still want to commit the rest. Keep this feature rare. Prefer small transactions.

Do not nest transactions as if SQL had true nested commit. Some products offer savepoints or ignore a second `BEGIN`. Check the manual.

Write `BEGIN` immediately before the writes. Write `COMMIT` immediately after the last check. Do not wait for a user click inside an open transaction.

### Questions

#### Theoretical questions

1. What statement starts an explicit transaction in this path?
2. What does `COMMIT` do?
3. What does `ROLLBACK` do on a disconnect before commit?
4. What is a savepoint?
5. Why must you not wait for a user inside an open transaction?

#### Easy practical tasks

1. Begin, insert a scratch row, commit. Select the row from a second session if you can.
2. Begin, insert a scratch row, rollback. Confirm that the row is gone.
3. Write the product aliases for `BEGIN` if they exist (`START TRANSACTION`).
4. Explain in four sentences when you choose rollback.

#### Medium practical tasks

1. Begin a transfer. Fail the second update on purpose. Roll back. Show both balances.
2. Try `CREATE TABLE` inside a transaction and roll back. Write whether the table remains.
3. Use a savepoint: insert two rows, roll back to the savepoint, commit. Show which row remains.

#### Advanced practical tasks

1. Document the full transaction keyword set of your DBMS: begin, commit, rollback, savepoint.
2. Write a script that aborts when a check fails (`balance < 0`) and rolls back. Show the messages.

---

## Autocommit vs explicit transactions

Autocommit means each statement is its own transaction. The client or the server commits after every successful statement. A later failed statement does not undo the earlier ones.

Explicit transactions mean you send `BEGIN` and later `COMMIT` or `ROLLBACK`. Several statements share one fate.

Many GUI clients default to autocommit. Many CLI clients also default to autocommit. Turn it off when you practice multi-step writes, or always wrap those writes in `BEGIN`.

```text
Autocommit on:
  INSERT order header;     -- committed
  INSERT order line;       -- fails
  -- header remains

Explicit:
  BEGIN;
  INSERT order header;
  INSERT order line;       -- fails
  ROLLBACK;
  -- header does not remain
```

Autocommit is convenient for a single `SELECT` or a single `INSERT`. It is dangerous for a business event that needs two writes.

Some products treat DDL as autocommit even inside a transaction. A `CREATE INDEX` can commit the work so far. Do not mix long DML and DDL in one transaction until you know the product.

A connection pool can return a connection that still has an open transaction if a program forgets to commit. That bug leaks locks. Always finish the transaction in application code (`COMMIT` or `ROLLBACK` in `finally` / `defer`). This path stays in SQL, but you must know the risk.

Do not disable autocommit globally on a shared GUI and then walk away. Your session can hold locks.

### Questions

#### Theoretical questions

1. What does autocommit do after each statement?
2. Why can autocommit leave an order header without lines?
3. When is autocommit acceptable?
4. Why can DDL break an explicit transaction?
5. What happens if a pooled connection stays in an open transaction?

#### Easy practical tasks

1. Find the autocommit setting in your CLI or GUI. Write where it is.
2. With autocommit on, run a good insert and a bad insert. Show the leftover row.
3. Repeat with an explicit transaction and rollback. Show that both are gone.
4. Write four sentences on when you switch autocommit off.

#### Medium practical tasks

1. Toggle autocommit off. Run two updates. Commit. Confirm from a second session.
2. Leave a transaction open. From a second session, try to update the same row. Write what you see (block or wait). Then commit or roll back.
3. Read whether DDL autocommits in your DBMS. Write the fact.

#### Advanced practical tasks

1. Write an application-level checklist: begin, work, commit or rollback in all exit paths. Keep it vendor-neutral.
2. Compare two clients (CLI and GUI) for autocommit defaults. Write a safe practice rule for each.

---

## Dirty read, non-repeatable read, phantom, lost update

Concurrent transactions can produce anomalies. Learn the names. The isolation level controls which anomalies can appear.

**Lost update.** Two sessions read the same value. Each adds one. Each writes. The last write wins. One increment disappears.

```text
Stock is 10.
T1 reads 10. T2 reads 10.
T1 writes 11. T2 writes 11.
Result is 11, not 12.
```

Use one `UPDATE stock = stock + 1` in a single statement, or lock the row, or use a higher isolation plus retries.

**Dirty read.** A session reads a value that another session has not committed. The writer then rolls back. The reader used a value that never existed as committed data.

**Non-repeatable read.** A session reads a row. Another session commits a change to that row. The first session reads the same row again and sees a new value.

**Phantom read.** A session reads a set of rows that match a `WHERE`. Another session commits an insert that also matches. The first session reads again and sees a new row. The extra row is a phantom.

```text
T1: SELECT COUNT(*) FROM orders WHERE customer_id = 1;  -- 3
T2: INSERT a new order for customer 1; COMMIT;
T1: SELECT COUNT(*) again;  -- 4  (phantom)
```

Do not treat all anomalies as equal. A dirty read is usually unacceptable for money. A phantom may be acceptable for a loose report.

The next section maps anomalies to isolation levels.

### Questions

#### Theoretical questions

1. What is a lost update?
2. What is a dirty read?
3. What is a non-repeatable read?
4. What is a phantom read?
5. Why is a dirty read worse than a phantom for a bank transfer?

#### Easy practical tasks

1. Draw the lost-update timeline with stock `10`.
2. Write four definitions in your own words, one sentence each.
3. Label a story: "I saw an order that then vanished" as dirty read or something else.
4. Explain a phantom with a `COUNT` query in four sentences.

#### Medium practical tasks

1. Reproduce a lost update with two sessions and read-plus-write (not `SET col = col + 1`). Show the wrong final value.
2. Try to read an uncommitted row from a second session. Write whether your default isolation allows it.
3. In one long transaction, select a row twice. Change it from another session between the reads. Record whether the value changed.

#### Advanced practical tasks

1. Write a lab script with two sessions for each anomaly that your isolation level can show. Record the results.
2. Explain why `UPDATE ... SET stock = stock + 1` avoids the lost update that the read-then-write pattern causes.

---

## Isolation levels

SQL defines isolation levels. Each level allows or forbids anomalies. The names from weakest to strongest:

1. `READ UNCOMMITTED`
2. `READ COMMITTED`
3. `REPEATABLE READ`
4. `SERIALIZABLE`

Typical anomaly map (standard idea; products differ):

| Level | Dirty read | Non-repeatable | Phantom |
| --- | --- | --- | --- |
| Read uncommitted | Possible | Possible | Possible |
| Read committed | No | Possible | Possible |
| Repeatable read | No | No | Possible (in the standard) |
| Serializable | No | No | No |

Many products use `READ COMMITTED` as the default. Some products implement `REPEATABLE READ` so that phantoms do not appear. Do not memorize only the standard table. Read the product page.

Set the level for a transaction:

```sql
BEGIN;
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
-- statements
COMMIT;
```

The exact `SET` syntax varies. Some products set the level on the session.

**Read uncommitted.** Almost never use this level for a relational business database. It allows dirty reads.

**Read committed.** Each statement sees only committed data. Two statements in one transaction can see different committed versions. This level is a common default.

**Repeatable read.** Rows that you already read stay stable in this transaction. Another session's update to those rows waits or does not appear.

**Serializable.** The DBMS must make the result equal to some serial order of transactions. If it cannot, it aborts one transaction. You retry.

Stronger isolation increases waits and aborts. Do not set serializable on every report. Use it for operations that must not interleave, such as a unique booking of the last seat.

Lost update is not always in the simple table. Row locks and `UPDATE` expressions still matter at read committed.

### Questions

#### Theoretical questions

1. What is the weakest isolation level in the SQL list?
2. Which typical default forbids dirty reads but allows non-repeatable reads?
3. What extra guarantee does serializable claim?
4. Why can a product disagree with the standard phantom column?
5. Why is serializable not the default for every query?

#### Easy practical tasks

1. Read the default isolation level of your DBMS. Write it.
2. Write the four level names in order from weak to strong.
3. Fill a copy of the anomaly table from this section.
4. Explain in four sentences when you would choose serializable.

#### Medium practical tasks

1. Set isolation to `READ COMMITTED`. Try the two-read experiment on one row.
2. Set isolation to `REPEATABLE READ` or `SERIALIZABLE` if the product allows it. Repeat the experiment. Write the difference.
3. Force a serializable abort if you can (two conflicting transactions). Record the error.

#### Advanced practical tasks

1. Write a one-page isolation guide for a shop: browse catalog, place order, run daily report.
2. Compare your product isolation page with the standard table. List two differences.

---

## Locks, deadlocks, and retry

A lock is a reservation that one transaction holds so that another transaction cannot do a conflicting action.

**Row lock.** The transaction locks one row (or the index entry that leads to it). Writers of that row wait. Readers may still read, depending on the isolation level and the product (some products use versions instead of blocking readers).

**Table lock.** The transaction locks the whole table. Other sessions cannot write, or cannot even read, depending on the lock mode. Table locks appear for some DDL and for some `LOCK TABLE` statements.

Lock modes (high-level):

- shared: many holders can read
- exclusive: one holder can write

An `UPDATE` or `DELETE` takes an exclusive row lock. A `SELECT` may take no lock or a shared lock. Some products never lock rows for a plain `SELECT` because they read a version.

```sql
BEGIN;
UPDATE products SET price = 9.99 WHERE product_id = 10;
-- row 10 is locked until COMMIT or ROLLBACK
COMMIT;
```

A second session that updates the same row waits. If it waits too long, it can time out.

Do not lock a table for a long report if a row-level plan exists. Do not lock rows that you do not need. Keep transactions short so that locks release soon.

`SELECT ... FOR UPDATE` (where the product supports it) locks selected rows as if you will write them. Use it when you read a seat and then write a booking in the same transaction.

A deadlock happens when two transactions wait for each other. Neither can proceed.

```text
T1 locks row A. T2 locks row B.
T1 wants row B. T2 wants row A.
Both wait forever unless the DBMS breaks the cycle.
```

The DBMS detects the cycle. It aborts one transaction. That session receives an error. The other session continues.

Your program must retry the aborted transaction from the start. Do not retry only the last statement. The aborted transaction lost all of its writes.

Reduce deadlocks:

- Lock rows in a stable order (always account `1` then account `2`).
- Keep transactions short.
- Do not hold a lock while you call a slow network service.
- Avoid mixing many random row updates in one transaction.

```sql
-- Always update the lower id first
UPDATE accounts SET balance = balance - 50 WHERE account_id = 1;
UPDATE accounts SET balance = balance + 50 WHERE account_id = 2;
```

If another session uses the same id order, the deadlock pair is less likely.

Do not treat a deadlock as a rare crash that you ignore. Under load, deadlocks appear. Retry is part of the design.

The abort is a rollback. Constraints still hold. No half write remains from the aborted transaction.

Some products report a deadlock with a specific error code. Catch that code. Back off a little. Retry a limited number of times. Then fail.

You do not pick page locks as a beginner. The DBMS chooses the granularity. Think in rows and tables.

### Questions

#### Theoretical questions

1. What does a lock reserve, and what does a row lock cover?
2. What happens when two sessions update the same row?
3. What is a deadlock, and how does the DBMS usually break it?
4. Why must you retry the whole transaction after a deadlock abort?
5. Why does a stable lock order help, and why are deadlocks part of normal design under load?

#### Easy practical tasks

1. In session A, begin and update one row. Do not commit. In session B, update the same row. Write what you observe. Then commit session A.
2. Draw the two-row deadlock timeline.
3. Write the retry rule in four sentences. List four ways to reduce deadlocks.
4. Explain in four sentences why short transactions release locks sooner.

#### Medium practical tasks

1. Try `SELECT ... FOR UPDATE` if the product has it. Show that a second writer waits.
2. Reproduce a deadlock with two sessions and two rows if you can. Record the error.
3. Rewrite a transfer so that both sessions lock ids in the same order. Find the deadlock error code or message in your DBMS docs.

#### Advanced practical tasks

1. Write a booking pattern: select seat for update, check free, update, commit. Stay conceptual if the syntax differs. Include a retry policy: max attempts, backoff, which errors retry.
2. Read whether your DBMS uses versions (MVCC) so that readers do not wait for writers. Write six sentences. Then write a one-page incident note: how you would explain a deadlock abort to a beginner developer.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do ACID, explicit transactions, and isolation levels work together in a money transfer?
2. Which anomaly remains possible at read committed, and how do you design around it?
3. When do you use row locks versus a stronger isolation level?
4. What must an application do after a deadlock abort?
5. A teammate leaves autocommit on and runs two updates for one business event. Which two risks do you name?

#### Easy practical tasks

1. Write a cheat sheet: ACID, begin/commit/rollback, autocommit, four anomalies, four levels, locks, deadlock retry.
2. Run one commit demo and one rollback demo on a scratch account table.
3. Write the default isolation level of your DBMS from memory, then check the docs.
4. Draw a deadlock and mark which session the DBMS may abort.

#### Medium practical tasks

1. With two sessions, show a wait on the same row and then a commit.
2. Set a non-default isolation level and run a two-read test.
3. Write a transfer script that rolls back when a `CHECK` on balance fails.

#### Advanced practical tasks

1. Design a last-seat booking: isolation or `FOR UPDATE`, conflict behavior, and retry. Implement a tiny practice version.
2. Write a lab report with two sessions that demonstrates one anomaly and one isolation level that removes it.
