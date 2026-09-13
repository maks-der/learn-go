# 20. Internals (Advanced)

## Description

Internals explain how Kafka stores positions and how clients talk to brokers. You do not need this topic to produce and consume. You need it when you debug lag, replication, or a commit that does not move.

This topic covers log end offset versus high watermark, the produce and fetch protocol, the KRaft controller quorum, index files, and the consumer offset commit protocol. Complete this topic after topics 2, 5, 6, 8, and 9. Use KRaft. This topic mentions ZooKeeper only as old history that you do not install.

Use one term for each concept. The log end offset (LEO) is the next offset that the replica will write. The high watermark (HW) is the offset that is fully replicated to the ISR. The controller quorum is the KRaft Raft group that stores cluster metadata. An index file maps offsets or times to positions in a segment. A commit is a write of a group offset to `__consumer_offsets`.

---

## Log end offset vs high watermark

Each replica of a partition has a **log end offset (LEO)**. LEO is the offset of the next record that this replica will append. If the last stored record is offset 41, LEO is 42. Different replicas can have different LEO values when a follower is behind.

The **high watermark (HW)** is the largest offset that the leader has written and that all replicas in the current ISR have also written (the exact definition follows the current Kafka replication design: HW advances when the ISR has the record). Consumers that use the normal read path can fetch records **below the HW**. They cannot read records that are only on the leader and not yet on the ISR (those offsets are above HW).

`acks=all` waits for the ISR (and `min.insync.replicas`, topic 9). After that produce succeeds, those records can become visible as HW moves.

A follower fetches from the leader and raises its LEO. When the ISR is caught up, the leader raises HW.

If you compare LEO on the leader with HW, the difference is data that is not yet fully replicated (or not yet marked). A large gap means lag (topic 9).

Do not use LEO as the consumer "latest" in operations talk without saying HW. "End of log" in a consumer API is the log end that the client is allowed to see (HW for a normal consumer).

KRaft does not change LEO and HW. Those numbers belong to the partition log, not to the metadata quorum.

### Questions

#### Theoretical questions

1. What is the log end offset?
2. What is the high watermark?
3. Why can a consumer not read above the HW?
4. How does a follower raise its LEO?
5. Why must you not mix LEO and HW in an incident report?

#### Easy practical tasks

1. Write five sentences about LEO and HW. Use only facts from this section.
2. Draw a leader LEO=100 and a follower LEO=90. Mark a possible HW.
3. Make a table: Term, who holds it, who may read up to it.
4. Find official or design-doc wording for high watermark. Rewrite it in STE.

#### Medium practical tasks

1. On a two-broker lab, stop a follower. Produce. Describe or use a tool to reason about HW versus leader LEO. Write what consumers still see.
2. Start the follower. Watch replication. Write when consumers can see the new records.
3. Relate HW to `acks=all` and ISR in six sentences (topic 9).

#### Advanced practical tasks

1. Measure offsets: produce a known count, read `LogEndOffset` and `HighWatermark` metrics if available. Write both.
2. Write a one-page debug note: consumer at end, producer success, HW not moved — what you check (ISR, follower fetch).

---

## Fetch and produce protocol

Clients and brokers speak the **Kafka protocol** on TCP. The protocol is a set of request types with versions. Common data-plane types:

- **Produce**: the producer sends record batches for partitions. The broker answers with errors and base offsets (and more fields in current versions).
- **Fetch**: the consumer (or follower) asks for records from partitions, with a fetch offset, min bytes, and max wait (topic 16). The broker answers with record batches up to the HW for consumers.

Followers also use Fetch to replicate. The follower fetch offset is that replica’s LEO. This is the replication path, not a consumer group.

Other request types exist: Metadata, FindCoordinator, JoinGroup, OffsetCommit, OffsetFetch, Txn requests. You do not memorize every field. You must know that a client is not "just HTTP".

The client uses `bootstrap.servers` to get **metadata** (leaders, epochs). Then produce and fetch go to the **leader** of each partition (fetch from follower is an optional path).

Protocol versions are negotiated. An old client on a new broker uses an older version. A too-old client can fail.

KRaft changes how metadata is stored. The Metadata request still exists. Clients still bootstrap to brokers. There is no ZooKeeper protocol for clients.

Do not implement a raw protocol parser as a first exercise. Use a maintained client. Read the protocol guide when you debug with a packet tool or broker logs.

### Questions

#### Theoretical questions

1. What does a Produce request carry?
2. What does a Fetch request ask for?
3. Who else uses Fetch besides a consumer?
4. Why does the client call Metadata before it produces to a new partition?
5. Do clients speak to ZooKeeper in current Kafka?

#### Easy practical tasks

1. Write five sentences about produce and fetch. Use only facts from this section.
2. Make a table: Request, typical sender, typical receiver.
3. Find the official protocol / API page. Write three request names.
4. Draw: bootstrap → metadata → produce to leader → fetch from leader.

#### Medium practical tasks

1. Enable client debug logs for one produce and one consume. Write the request names that you see (or the client API calls that map to them).
2. Read the Fetch request fields in the protocol docs (offset, max wait, min bytes). Map them to consumer config keys.
3. Compare follower fetch and consumer fetch in five sentences.

#### Advanced practical tasks

1. Use a current protocol listing to write the versioning idea (flexible versions) in six STE sentences. Do not implement it.
2. Write a debug checklist: metadata stale, not leader, fetch offset out of range — which request and which error.

---

## Controller quorum (KRaft)

**KRaft** is Kafka Raft. A **controller quorum** is a set of Kafka processes with the controller role. They replicate a **metadata log** with Raft. One node is the active controller. The others are voters (and observers in some setups).

The metadata log stores topics, partitions, replicas, ACLs (current versions), and other cluster state. It does not store application records. Application records stay in partition logs in `log.dirs` (topic 8).

`process.roles` can be `broker`, `controller`, or both. Combined mode is common in small labs. Production can split controllers and brokers.

A Raft **majority** must be up to accept metadata changes (create topic, move leader). If a majority of controllers stop, the cluster cannot elect leaders or change topics. Existing partition leaders can still serve produce and fetch for a time, but operations degrade. Do not stop a majority (topic 15).

Clients do not connect to ZooKeeper. Tools use `--bootstrap-server`. `kafka-storage format` writes the cluster id for KRaft (topic 1).

This handbook mentions ZooKeeper only to say: do not install it for new clusters. Old Kafka stored metadata in ZooKeeper. You do not learn that path here.

Controller epoch and partition leader epoch protect stale leaders. A stale broker must not append as leader after a new epoch.

### Questions

#### Theoretical questions

1. What does the KRaft metadata log store?
2. What does the metadata log not store?
3. What is a Raft majority in this section?
4. What does `process.roles` select?
5. Why do admin tools use `--bootstrap-server` and not ZooKeeper?

#### Easy practical tasks

1. Write five sentences about the controller quorum. Use only facts from this section.
2. Make a table: Combined node, dedicated controller. Add when you use each.
3. Find KRaft configuration keys (`process.roles`, `controller.quorum.voters` or the current name). Write them from the docs.
4. Draw three voters. Mark majority. Mark a forbidden "two down" state.

#### Medium practical tasks

1. On a three-controller KRaft lab, stop one controller. Create a topic. Write whether it works. Stop a second. Try again. Write the result.
2. Read official KRaft migration or concepts pages. Write ten facts. Mention ZooKeeper in at most one sentence that says you do not use it.
3. Describe how a leader epoch stops a stale leader. Use official wording rewritten in STE.

#### Advanced practical tasks

1. Inspect controller logs during a topic create. Write the events you see (without pasting secrets).
2. Write a KRaft operations page: voter count, disk for metadata, TLS on controller listener (topic 17), rolling restart rule.

---

## Index files (.index, .timeindex)

A partition segment is a `.log` file plus index files (topic 8). Two common indexes:

- **`.index`** (offset index): maps a record offset to a byte position in the `.log` file. The broker uses it to start a fetch at an offset without a full scan of the segment.
- **`.timeindex`** (time index): maps a timestamp to an offset. The broker uses it for time-based retention and for fetches by timestamp (`offsetsForTimes` in the client).

Indexes are sparse. They do not store every record. Kafka finds the nearest entry and then scans a small range in the `.log` file.

There can be more files (for example a transaction index in some versions). Do not edit any of them by hand. A bad index can make the partition unreadable until recovery.

Recovery rebuilds indexes from the `.log` if they are missing or corrupt. That takes time on large segments. That is one reason huge segments hurt restart (topic 8 and 15).

KRaft metadata segments are a different log. Do not mix those files with application partition indexes.

Do not delete indexes to "save disk" while the broker runs.

### Questions

#### Theoretical questions

1. What does `.index` map?
2. What does `.timeindex` map?
3. Why are indexes sparse?
4. Why must you not edit index files?
5. Why do large segments make index rebuild slow?

#### Easy practical tasks

1. Write four sentences about index files. Use only facts from this section.
2. In a lab `log.dirs`, list files for one partition. Write the extensions that you see.
3. Make a table: File, purpose, safe to edit (no).
4. Find official log format or storage notes. Write one sentence about the offset index.

#### Medium practical tasks

1. Produce enough data to close a segment. List `.log`, `.index`, and `.timeindex` for that base offset.
2. Use a time-based seek in a consumer (beginning of a timestamp). Write how the time index helps (from docs and this section).
3. Read about `log.index.interval.bytes` or the current sparse-interval key. Write the meaning in your own words.

#### Advanced practical tasks

1. In a stopped lab broker (never in production), note what official recovery does if an index is missing. Do not corrupt a production disk. Write the recovery idea.
2. Write a storage note: segment size, index interval, why you never delete indexes by hand, KRaft metadata separate.

---

## Consumer offset commit protocol

A consumer group stores committed offsets in the internal topic `__consumer_offsets` (topics 6 and 10). The **offset commit protocol** is the request path:

1. The client finds the **group coordinator** (FindCoordinator).
2. The client sends **OffsetCommit** with group id, member (or generation), partitions, and offsets (plus metadata).
3. The coordinator appends a commit record to `__consumer_offsets` and answers success or error.

**OffsetFetch** reads the last committed offsets when a member starts or rebalances.

Auto-commit and manual commit (topic 5) both use this protocol. The difference is when the client calls it.

Transactional commit (topic 7) can send offsets inside a transaction (`sendOffsetsToTransaction`). Then the commit of the transaction and the offsets is one atomic unit with the produce. Isolation and `__consumer_offsets` still apply.

A commit is not a consume. A commit is a produce to an internal compacted topic. That topic has partitions. The coordinator owns the group.

If the coordinator moves, the client finds it again. Commits can fail with a retriable error. The client must retry. A lost commit after a process crash causes a replay (at-least-once).

KRaft does not replace `__consumer_offsets`. Group data is still in that topic on brokers. You do not commit offsets to ZooKeeper in current Kafka.

### Questions

#### Theoretical questions

1. Which topic stores group offsets?
2. What is the role of the group coordinator in a commit?
3. What is the difference between OffsetCommit and OffsetFetch?
4. How can a transaction include offsets?
5. Do current clients commit offsets to ZooKeeper?

#### Easy practical tasks

1. Write five sentences about the commit protocol. Use only facts from this section.
2. Make a table: Request, when the client sends it.
3. Find OffsetCommit in the official protocol docs. Write three fields that the page names.
4. Draw: consumer → coordinator → `__consumer_offsets`.

#### Medium practical tasks

1. Commit manually in a lab. Use `kafka-consumer-groups --describe`. Write the offset. Consume one record. Commit. Describe again.
2. Stop a consumer after process but before commit (or kill). Restart. Write whether a record is repeated.
3. Read `__consumer_offsets` notes (compaction, how many partitions). Write three facts.

#### Advanced practical tasks

1. Use a transactional producer/consumer path (topic 7) and confirm offsets commit with the transaction. Write the isolation setting on a second consumer.
2. Write a commit-protocol note: coordinator, internal topic, retry, KRaft, no ZooKeeper, link to topic 5 pitfalls.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do LEO, HW, Fetch, and a consumer offset work together when a follower is behind?
2. What two logs exist on a KRaft node (metadata versus partition), and which indexes belong to which?
3. How does the produce protocol relate to HW before a consumer Fetch can see a record?
4. Why is `__consumer_offsets` an application of the same log idea as an application topic?
5. What do you refuse from old internals articles that start ZooKeeper?

#### Easy practical tasks

1. Write a one-page cheat sheet: LEO, HW, Produce, Fetch, KRaft quorum, two indexes, commit requests.
2. Draw one partition: `.log`, `.index`, `.timeindex`, leader LEO, HW, consumer offset.
3. Bookmark official protocol, KRaft, and storage pages.
4. List every offset kind named in this topic (LEO, HW, fetch offset, committed offset) in one table.

#### Medium practical tasks

1. On a KRaft lab, produce, consume, describe group, and list segment files. Write one paragraph that uses all five subsection ideas.
2. Map each subsection to one official URL.
3. Write a glossary of ten internals terms from this topic in STE.

#### Advanced practical tasks

1. Debug a staged incident: follower down, produce `acks=all` behavior, consumer at HW, then restore follower. Write LEO/HW/ISR notes.
2. Write an internals standard for on-call: which metrics (HW, LEO, URP, controller), KRaft quorum, never edit indexes, never use ZooKeeper.
