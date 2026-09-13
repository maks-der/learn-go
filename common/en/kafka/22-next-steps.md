# 22. Next Steps

## Description

This topic is a practice list. You apply topics 1–21 on a KRaft cluster. You write commands and small programs. You observe the cluster. You do not get solutions in this handbook. You write the results in your own notes.

This topic covers produce and consume in your language, a partition-assignment drawing for three consumers, compaction on a changelog topic, Avro and schema compatibility checks, a small Kafka Streams job, and a slow poll that causes a rebalance. Complete the earlier topics first. Use KRaft. Do not start ZooKeeper.

Use one term for each concept. A practice is a task that you run. An observation is a fact that you record from output or UI. A drawing is a diagram that you make before or after the run. Do not copy a finished program from this file. This file states the goal and the checks.

---

## Produce and consume in your language

Write a producer and a consumer in a language that you already use (Java, Go, Python, or another language with a maintained Kafka client).

Goals:

- Connect with `bootstrap.servers` to a local KRaft cluster (topic 1).
- Produce at least three records to a topic that you create (topic 10). Use a key on at least two records (topic 3).
- Consume with a `group.id` and print key, value, partition, and offset (topics 2 and 5).
- Close the client on shutdown.

Checks (you write the answers):

- The consumer sees all records if it starts from the beginning or if it starts before you produce.
- A second start of the same group does not reprint committed records if you committed (topic 5).
- The client does not mention ZooKeeper.

Do not use only the console tools for this practice. Console tools are allowed as a second check.

Serialize as UTF-8 strings or JSON first. Avro comes in a later section.

If the client needs a framework, keep the program small. One file is enough.

### Questions

#### Theoretical questions

1. Which configuration key names the cluster bootstrap?
2. What four fields must you print for each consumed record in this practice?
3. Why must this practice use a language client, not only the console?
4. What proves that the group committed offsets?
5. Why must the client not use ZooKeeper?

#### Easy practical tasks

1. Create topic `next.lang` with three partitions on KRaft. Write the exact command.
2. Write a producer outline (function names only, no full program) that sends three records.
3. Write a consumer outline that polls and prints. No full program.
4. Make a table: Language, client library name, document URL.

#### Medium practical tasks

1. Run your producer and consumer. Save commands and output in your notes (no secrets).
2. Restart the consumer. Record whether old records appear. Write why.
3. Produce with two keys. Write the partitions that you saw. Relate to topic 3.

#### Advanced practical tasks

1. Add a second consumer in a second group. Confirm both groups can read the same records from the beginning.
2. Write a short report: library version, Kafka version, KRaft proof (`process.roles` or vendor mode), and one error that you fixed.

---

## Draw partition assignment for a group of 3 consumers

Take one topic with a known partition count. Use three consumers with the same `group.id`.

Before you start the consumers, **draw** a possible assignment: which consumer owns which partition. Use the rule from topic 6: one partition belongs to at most one active member.

Then start the three consumers. Describe the group (`kafka-consumer-groups --describe` or the Admin API). **Draw the real assignment** next to your guess.

Change the partition count or the member count (stop one consumer). Draw again. Record the rebalance.

Use three partitions first (one each). Then try six partitions (some members get two). Then try two partitions (one member is idle). Write what "idle" means.

This practice is a drawing practice. A correct drawing is the deliverable. Do not skip the paper or file diagram.

KRaft does not change the assignment rule. The assignor name can appear in the describe output (topic 6).

### Questions

#### Theoretical questions

1. What rule limits how many active members can share one partition?
2. What happens when there are more members than partitions?
3. What happens when there are more partitions than members?
4. What output command shows the live assignment?
5. Why do you draw before you run?

#### Easy practical tasks

1. Draw three consumers and three partitions. Guess an assignment.
2. Draw three consumers and two partitions. Mark the idle member.
3. Draw three consumers and six partitions. Guess counts per member.
4. Write the `kafka-consumer-groups` describe command shape for your group.

#### Medium practical tasks

1. Run the three-consumer lab. Compare guess and actual. Write differences.
2. Stop one consumer. Draw before and after. Save describe output.
3. Write the assignor name if the output shows it. Relate to topic 6.

#### Advanced practical tasks

1. Repeat with cooperative sticky (if you set it). Draw two rebalances. Write which partitions moved.
2. Write a one-page note: static membership (topic 6) and how your drawing would change if you used it.

---

## Enable compaction on a changelog topic

Create a topic that you treat as a changelog (topics 8 and 13). Set `cleanup.policy=compact`. Use a key on every record.

Goals:

- Produce a sequence: key `A` value `1`, key `B` value `1`, key `A` value `2`.
- Wait for compaction or force a condition that lets compaction run (segment roll, `min.cleanable.dirty.ratio`, enough closed segments). Read topic 8. Do not invent a delete of `.log` files.
- Consume from the beginning after compaction had a chance to run. Write which values remain for `A` and `B`.
- Produce a tombstone (null value) for `A`. After compaction, write whether `A` is gone.

Checks:

- You can describe the topic and show `cleanup.policy=compact`.
- You can explain why the active segment might still hold old values (topic 8).
- The cluster is KRaft.

Optional: set `compact,delete` and a retention (topic 8) if you want both behaviors. Write why you chose that.

Do not compact `__consumer_offsets` by hand. That topic is already internal.

### Questions

#### Theoretical questions

1. What does `cleanup.policy=compact` keep per key?
2. What is a tombstone in this practice?
3. Why can old values still appear if you consume immediately?
4. Which topic type in Streams also uses compaction (topic 13)?
5. Why must you not delete segment files to force compaction?

#### Easy practical tasks

1. Create `next.compact` with compaction. Write the create or `kafka-configs` command.
2. Describe the topic. Write the cleanup policy from the output.
3. Write the three produce records (A/B/A) as a table: key, value, expected latest.
4. Draw a compacted log before and after a tombstone for `A`.

#### Medium practical tasks

1. Run produce, wait or roll segments, consume from beginning. Write the records that remain.
2. Find the configuration that affects when compaction runs. Write the key names from the docs.
3. Compare this topic with a delete-only topic that has the same produces. Write the difference after retention.

#### Advanced practical tasks

1. Measure how long your lab took until `A=1` disappeared. Write segment settings that you used.
2. Write a changelog standard: when to compact, tombstone rules, and why Streams internal topics must not be edited by hand.

---

## Add Avro + schema compatibility checks

Run a Schema Registry (Confluent or Karapace) next to KRaft Kafka (topic 11).

Goals:

- Define an Avro schema for a small record (for example, order id and amount).
- Register it with a subject. Use TopicNameStrategy unless you write a reason.
- Produce and consume with the Avro serializer and deserializer.
- Set a compatibility mode (start with BACKWARD).
- Register a second schema version that is compatible. Confirm the register succeeds.
- Attempt a change that the mode must reject. Confirm the register fails. Export the error text.

Checks:

- The wire value is not raw JSON (magic byte and schema id, topic 11).
- The consumer of version 1 can read version 2 if you chose BACKWARD and a compatible add.
- You did not start ZooKeeper.

Do not paste a full schema from a vendor tutorial as your only work. Write your own field names.

JSON Schema or Protobuf is acceptable if you cannot run Avro. The compatibility check is required.

### Questions

#### Theoretical questions

1. What three parts does a Confluent-style Avro value have on the wire?
2. What does BACKWARD allow a consumer to do?
3. What is a subject?
4. Why do you try a rejected schema on purpose?
5. Why is the registry not Kafka?

#### Easy practical tasks

1. Write an Avro schema with two fields in your notes (your own names).
2. Write the subject name that TopicNameStrategy would use for your topic and a value.
3. Find the registry HTTP path to list subjects. Write it.
4. Make a table: Change, BACKWARD ok (yes or no) — add optional field, remove field, rename field.

#### Medium practical tasks

1. Register v1, produce, consume. Write the schema id.
2. Register a compatible v2. Produce v2. Consume with a v1 reader if you can. Write who succeeds.
3. Register an incompatible schema. Save the error. Then stop.

#### Advanced practical tasks

1. Add a CI-style check: a script that registers a schema file and fails the run on incompatibility. Point it at the lab registry.
2. Write a schema standard: format, compatibility, strategy, KRaft Kafka, registry product name.

---

## Build a small Kafka Streams job

Write one Kafka Streams topology (topic 14). Java is the usual language.

Minimum job (pick one):

- Word count, or
- Filter orders by amount, or
- Stream-table join of orders and customers

Goals:

- Set `application.id` and `bootstrap.servers` to KRaft.
- Read at least one source topic. Write one sink topic.
- If the job is stateful, list the changelog topics after the first run.
- Unit-test the topology with `TopologyTestDriver` (topic 14) **or** write why you ran only an integration test.
- Set `processing.guarantee` to `at_least_once` or `exactly_once_v2`. Write the choice.

Checks:

- Two instances with the same `application.id` split partitions (topic 6 and 14).
- You do not write to changelog topics with a console producer.
- No ZooKeeper property exists in your config.

Do not paste the official word-count sample unchanged. Change topic names and at least one rule.

If you cannot run JVM Streams, write the topology on paper and run the same idea as a small consumer-producer program. State that it is not Kafka Streams.

### Questions

#### Theoretical questions

1. What does `application.id` name?
2. When does Streams create a changelog topic?
3. Why must you not produce into a changelog topic?
4. What is the minimum test this practice asks for?
5. Which processing guarantee did this handbook recommend that you name in the config?

#### Easy practical tasks

1. Write the topology as a list of nodes (source → process → sink). No code.
2. Name `application.id`, source topic, and sink topic.
3. Make a table: Stateless or stateful, changelog expected (yes or no).
4. Find `StreamsBuilder` in your project or docs. Write one method that you will call.

#### Medium practical tasks

1. Run the job on KRaft. Produce input. Consume the sink. Save output.
2. List topics. Mark internal names.
3. Run `TopologyTestDriver` for one case (one in, one out) if you use Java.

#### Advanced practical tasks

1. Start a second instance. Draw assignment (reuse the drawing skill). Kill one instance. Write restore or rebalance behavior.
2. Enable `exactly_once_v2`. Repeat a crash test. Write observations without claiming a formal proof.

---

## Break a consumer (slow poll) and watch a rebalance

Write or change a consumer so that it **violates** the poll contract on purpose (topics 5, 6, and 19).

Goals:

- Run two consumers in one group on a topic with at least two partitions.
- In one consumer, sleep (or block) longer than `max.poll.interval.ms` without calling poll. Lower `max.poll.interval.ms` in the lab so that you do not wait 5 minutes.
- Watch logs: leave group, revoke, assign.
- Describe the group during the event. Draw the assignment before and after.

Then **fix** the consumer: remove the sleep, or pause partitions and keep polling (topic 19). Confirm the group stays stable under a slow sink that uses pause.

Checks:

- You can quote a log line that shows the rebalance (in your notes, not in this handbook).
- You can explain session timeout versus poll interval (topic 5) in your own words.
- You did not blame KRaft for your sleep.

Do not run this against a production group.

This is the last practice. It combines consumers, groups, and backpressure.

### Questions

#### Theoretical questions

1. Which configuration do you lower in the lab so that the break happens soon?
2. What does the coordinator do when a member misses the poll interval?
3. What must you see in the second consumer when the first is expelled?
4. What is the difference between this break and a pause-based slow sink?
5. Why is this unsafe in production?

#### Easy practical tasks

1. Write the two keys: `max.poll.interval.ms` and `session.timeout.ms`. Write which one you will break.
2. Draw two consumers, two partitions, then one consumer left alone after the break.
3. Write a sleep time that is greater than your poll interval (numbers only).
4. List three log words you will search for (for example: revoke, assign, leave).

#### Medium practical tasks

1. Run the break. Save describe output and timestamps.
2. Restore a healthy poll. Confirm URP is unrelated (topic 15). Write that the incident was the client.
3. Repeat with pause/resume instead of sleep. Write that a rebalance did not occur (or did, and why).

#### Advanced practical tasks

1. Measure time from start of sleep to completed rebalance. Write the number and the configs.
2. Write a consumer runbook: how you detect a poll-interval rebalance, how you fix it, KRaft cluster unchanged.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the six practices map to topics 3, 5, 6, 8, 11, 14, and 19?
2. What one KRaft fact must be true in every practice on this list?
3. Why does this topic refuse to include solution programs?
4. Which practice proves compaction, and which practice proves group assignment?
5. How would you add one extra practice from the suggested list in `kafka.topics.md` (TLS + SASL) without changing the six sections?

#### Easy practical tasks

1. Write a checklist of the six practices with empty boxes. Tick them in your notes when done.
2. Draw one diagram that includes all six: client language, three members, compacted keys, registry, Streams, rebalance.
3. Bookmark official pages you used in the six practices.
4. List every topic name that you created in this topic (`next.*` or your names).

#### Medium practical tasks

1. Run all six practices on one KRaft Compose or local cluster. Write a single timeline of commands (no solution code).
2. Map each subsection to the earlier handbook file that you needed most.
3. Write a "what I still cannot do" list of three items and which topic you will reread.

#### Advanced practical tasks

1. Add TLS + SASL from topic 17 to the language client practice. Keep all six practices working. Write the extra property keys (no secrets).
2. Write your own next-ten-days plan: one Streams job in production shape, one Connect source, one DR drawing from topic 18. Still no ZooKeeper.
