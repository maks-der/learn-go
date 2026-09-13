# 9. Replica Sets

## Description

This topic shows how MongoDB copies data across members. A **replica set** is a group of `mongod` processes that hold the same data set. One member is the **primary**. The other data members are **secondaries**. You learn elections, the oplog, stale reads, and hidden or delayed members.

Complete CRUD, drivers, and transactions first. Complete this topic before you study sharding.

Use one term for each concept. The **primary** accepts client writes. A **secondary** applies the same writes from the **oplog**. An **election** selects a new primary when the current primary is not available. A **stale read** is a read that does not yet include a recent write.

Replication gives high availability. Replication is not a backup by itself. You still need a backup that you store off the replica set.

---

## Primary, secondaries, and elections

A replica set has an odd number of voting members in a typical production layout. Three voting members are a common minimum. The set elects one primary. The other data-bearing members are secondaries.

Writes go to the primary. The primary records each write in the oplog. Each secondary pulls those oplog entries and applies them. After apply, the secondary has the same documents as the primary, with a possible delay.

A replica set can include:

- Data-bearing members (`mongod` with a full copy)
- An **arbiter** (votes, holds no application data)

An arbiter is a special case. Atlas and current guidance prefer data-bearing members. Do not add an arbiter unless you understand the write-concern and storage trade-off.

Members talk on the replica-set ports. Clients use a connection string that lists seeds and `replicaSet=Name`, or an Atlas SRV URI.

Useful commands:

```javascript
rs.status()
rs.conf()
db.hello()
```

`rs.status()` shows health, state, and optimes. `db.hello()` shows if the current member is primary and lists hosts.

A standalone `mongod` is not a replica set. Multi-document transactions and change streams need a replica set (or a sharded cluster). For learning, you can start a one-member replica set.

The primary can change. Your driver must discover the new primary. Official drivers do this. Do not hard-code one host as "always primary".

Do not write to a secondary. Secondaries reject normal application writes.

An **election** selects a new primary. Voting members start an election when they cannot reach the primary, or when the primary steps down.

A member can become primary only if it can win a **majority of the votes**. In a three-member set, two votes are a majority. If two members are down, the last member cannot elect itself. The set has no primary. Writes fail.

The primary **steps down** when:

- You run `rs.stepDown()`
- A more eligible member appears
- The primary cannot see a majority (it becomes secondary to avoid a split-brain write)

**Priority** affects who wins. A member with priority `0` cannot become primary. Use that for a hidden or delayed member.

Elections take time. During the election, the set has no primary. Writes wait or fail. Official drivers retry some writes. Design the application for a short write outage.

Atlas creates a replica set when you create a cluster. You do not start `mongod` by hand. You still must know that a failover can occur and that a write with `w: 1` can roll back.

### Questions

#### Theoretical questions

1. What is a replica set?
2. Which member accepts application writes?
3. Why is an odd number of voting members common?
4. Why can one surviving member in a three-member set not accept writes?
5. What does priority `0` prevent?

#### Easy practical tasks

1. Run `db.hello()` on your deployment. Write if you have a replica set and which host is primary.
2. Run `rs.status()` or the Atlas metrics view. Write the member count and the states.
3. Draw three members. Mark votes. Show why two down members stop writes.
4. Open the replica-set introduction and the elections page. Write both URLs.

#### Medium practical tasks

1. Start a three-member replica set in Docker or on localhost. Insert one document. Find it on each member with a direct connection.
2. On a test replica set, run `rs.stepDown()`. Time the gap until `hello` shows a new primary. Write the seconds.
3. Stop two of three members. Attempt an insert. Write the error. Start the members again.

#### Advanced practical tasks

1. Read about voting and priority in `rs.conf()`. Set one member to priority `0`. Confirm it does not become primary.
2. Design a five-member set for one region and one analytics member. Assign priorities. Explain the majority if the analytics site is down.

---

## Oplog

The **oplog** is a special capped collection. The default name is `local.oplog.rs`. The primary writes an entry for each operation that secondaries must apply.

Secondaries tail the oplog. They apply entries in order. The apply is how replication works.

The oplog is **capped**. Old entries disappear when the collection is full. The time span of the remaining entries is the **oplog window**.

If a secondary is down longer than the oplog window, it cannot catch up from the oplog. It needs an initial sync (a full copy). That is expensive.

Oplog size is configurable. Atlas lets you grow the oplog. On Community Server you set oplog size in configuration or with a command that the manual lists for your version.

The oplog is not a general event bus for your application. Change streams use the oplog internally. Your application must use change streams or another approved API, not a raw tail of `local.oplog.rs`, unless you are an operator who knows the risks.

Do not store application data in the `local` database.

`db.printReplicationInfo()` and `rs.printReplicationInfo()` show oplog window estimates on many versions. Atlas shows oplog metrics in the UI.

### Questions

#### Theoretical questions

1. What collection holds the oplog?
2. Who writes oplog entries for normal application writes?
3. What happens when the oplog is full?
4. What is the oplog window?
5. Why is a long secondary outage a risk for the oplog window?

#### Easy practical tasks

1. On a replica set, run `use local` then `show collections`. Confirm `oplog.rs` exists (if your role allows it).
2. Run `rs.printReplicationInfo()` or the Atlas oplog chart. Write the window that you see.
3. Open the oplog page in the manual. Write the URL.
4. Make a two-column table: "Oplog" and "Application collection". Add three differences.

#### Medium practical tasks

1. Insert 1000 documents. Compare oplog timestamps or Atlas oplog rate before and after. Write the change.
2. Read how to resize the oplog for your version. Write the steps. Do not run a resize on a shared cluster.
3. Explain initial sync vs oplog catch-up in six sentences.

#### Advanced practical tasks

1. Read which operations create oplog entries (including no-ops). Write three entry types at a high level.
2. Estimate an oplog window for a write-heavy lab. Propose a larger oplog or a faster recovery plan. Write the numbers.

---

## Stale reads from secondaries

A **stale read** is a read that misses a write that already succeeded on the primary. Replication is asynchronous unless you wait with write concern and read rules.

Default read preference is `primary`. Those reads are not secondary stale reads.

If you set read preference `secondary`, `secondaryPreferred`, or `nearest`, the driver can read a secondary. That secondary can lag.

Lag causes:

- A user does not see the document that the user just created
- A unique check on a secondary misses a new key
- Reports that do not match the primary

Use secondary reads when the application accepts lag. Examples: some dashboards, some search previews, some full-collection reports.

Do not use a secondary read to confirm a write that you just sent. Use the primary, or use a session with causal consistency and the concerns from the transactions topic.

Write concern `majority` does not remove stale reads on a secondary. Majority means the write is durable on a majority. A given secondary can still be behind.

Measure lag. Atlas shows replication lag. `rs.status()` shows optime dates. Alert when lag grows.

### Questions

#### Theoretical questions

1. What is a stale read?
2. Which read preference modes can hit a secondary?
3. Why can a user miss a document that the user just inserted?
4. Does `w: "majority"` make every secondary current?
5. When is a secondary read acceptable?

#### Easy practical tasks

1. Set read preference `secondary` in `mongosh` or a driver. Run a find. Write the host if the client shows it.
2. Open the stale-reads or read-preference page. Write one warning from that page.
3. Draw: client write, primary, two secondaries, a find that misses the write. Label lag.
4. Make a table: endpoint type, read preference, accept stale? Add three rows.

#### Medium practical tasks

1. Insert on the primary. Immediately find with `secondary`. Repeat until you see a miss or write that you did not see a miss. Write the result.
2. Compare `rs.status()` optimes for primary and one secondary. Write the time gap.
3. Write an API rule: checkout reads vs catalog reads. Give one reason for each.

#### Advanced practical tasks

1. Read about `maxStalenessSeconds`. Configure it in a driver. Write what happens when all secondaries are too stale.
2. Combine this section with write concern and causal consistency from topic 8. Write a one-page "read your writes" policy.

---

## Hidden and delayed members

A **hidden** member is a secondary that client drivers do not use for reads. Set `hidden: true`. A hidden member must have priority `0`. It can still vote. It still replicates.

Use a hidden member for backups or for jobs that must not take application read traffic.

A **delayed** member applies oplog entries after a fixed delay. The delay is `secondaryDelaySecs` (older docs used `slaveDelay`). A delayed member must be hidden and must have priority `0`.

Use a delayed member as a short "oops" window. Example: a delay of 3600 seconds keeps a copy from one hour ago. A bad drop can be repaired from that member if you act inside the delay. This is not a full backup strategy.

Do not point the application at a delayed member. The data is intentionally old.

Do not hide the only secondaries that you need for majority if you do not understand votes. Hidden members often still vote. Plan the majority.

Atlas offers related node types (for example analytics nodes) on some tiers. Read the Atlas page for your cluster. The idea is the same: keep special members off the default read path.

### Questions

#### Theoretical questions

1. What does `hidden: true` change for drivers?
2. Why must a hidden member have priority `0`?
3. What does a delayed member wait for?
4. Why is a delayed member not a complete backup plan?
5. Can a hidden member vote?

#### Easy practical tasks

1. Open the hidden-member and delayed-member pages. Write both URLs.
2. Write the three settings that a delayed member must have.
3. Draw a four-member set: primary, two normal secondaries, one delayed hidden member. Mark who serves reads.
4. List two jobs that fit a hidden member. List two jobs that do not.

#### Medium practical tasks

1. On a local test set, add a hidden member or change `rs.conf()` to hide one member. Confirm `hello` on a client does not list it for reads.
2. Set a short delay on a lab member (for example 60 seconds). Insert a document. Find it on the delayed member before and after the delay. Write the times.
3. Compare snapshot backup vs delayed member for a mistaken `drop()`. Write three differences.

#### Advanced practical tasks

1. Read vote configuration with a hidden member. Write how majority changes if you add a hidden voter.
2. Design a restore drill that uses a delayed member. Write the steps and the time limit.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the primary, the oplog, and secondaries work together on one insert?
2. When do elections, hidden members, and delayed members each change who can be primary?
3. Why is replication necessary for transactions and change streams but not sufficient as the only backup?
4. How do stale secondary reads and automatic failover create different application risks?
5. A teammate adds an arbiter and reads from `nearest` for checkout. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: primary, secondary, election majority, oplog window, stale read, hidden, delayed, Atlas SRV.
2. On one deployment (local or Atlas), record: member count, current primary, oplog window or "unknown".
3. Draw a timeline: write, oplog, secondary apply, election, new primary.
4. Mark three of your application reads as primary-only or secondary-acceptable. Give one reason each.

#### Medium practical tasks

1. Build a local three-member set. Step down the primary. Confirm the driver or `mongosh` finds the new primary and that a test document remains.
2. Write a one-page comparison: hidden member backup job vs Atlas snapshot vs delayed member. Include one limitation each.
3. Measure lag during a bulk insert. Write the peak lag and whether a secondary read would have been safe for your API.

#### Advanced practical tasks

1. Produce a replica-set standard for your team: member count, priorities, hidden/delayed policy, write concern, read preference, oplog sizing, Atlas vs self-managed.
2. Simulate a secondary that falls off the oplog window (or write the procedure from the manual). Record how you detect the problem and what initial sync means for downtime.
