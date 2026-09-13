# 10. Topics and Admin

## Description

Admin work creates and changes topics, sets configuration, and inspects groups. Kafka includes command-line tools in the `bin` folder. Clients also have an Admin API.

This topic covers create, describe, and delete of topics, per-topic config, internal topics, a preview of ACLs, and the common shell tools. Complete this topic after topics 2 and 8. Topic 17 covers security in more detail.

Use one term for each concept. An admin operation changes cluster metadata or topic configuration. An internal topic is a topic that Kafka uses for its own state. An ACL is an access-control rule. Use KRaft clusters. Admin tools use `--bootstrap-server`. Do not use old ZooKeeper flags on new clusters.

---

## Create / describe / delete topics

You create a topic with a name, a partition count, and a replication factor. Example command shape:

```text
kafka-topics --bootstrap-server localhost:9092 --create --topic orders.placed --partitions 3 --replication-factor 1
```

On Windows, use the `.bat` file. On macOS and Linux, use the `.sh` file. Some installs put the tools on `PATH` without a suffix.

`--describe` shows partitions, leaders, replicas, ISR, and configuration when you ask for configs.

`--list` shows topic names.

`--delete` removes a topic. The cluster must allow delete (`delete.topic.enable` or the current equivalent). Delete is not instant. Data leaves disk after the delete process. Do not delete a topic in production without a change process.

Auto-create can exist on a broker (`auto.create.topics.enable`). Do not rely on auto-create in production. Create topics with an explicit partition count and replication factor.

Topic names have rules. Use letters, digits, `.`, `_`, and `-`. Avoid spaces. Avoid names that collide with internal prefixes.

The Admin API in your language does the same operations from a program or a pipeline.

### Questions

#### Theoretical questions

1. What three properties do you set when you create a topic?
2. What does `--describe` show?
3. Why must you not rely on auto-create in production?
4. Why is delete not instant?
5. What address flag do new tools use instead of ZooKeeper?

#### Easy practical tasks

1. Create `admin.demo` with 3 partitions. Describe it. List topics. Write the commands.
2. Write four sentences about topic admin. Use only facts from this section.
3. Find `--help` for `kafka-topics`. Write five flags.
4. Try an invalid topic name. Record the error.

#### Medium practical tasks

1. Delete `admin.demo`. List topics until the name is gone. Write how long it took.
2. Create the same topic twice. Record the error. Use `--if-not-exists` if the tool has it.
3. Write a small Admin API program (or `kcat`/similar) that creates and describes a topic.

#### Advanced practical tasks

1. Create a topic with RF 3 on a three-broker cluster. Describe. Then delete. Confirm that log directories drop the partition folders after delete.
2. Write a one-page change process: who can create topics, required RF, required min ISR, and how you record the change.

---

## Config overrides per topic

Broker defaults apply to all topics. You can override keys on one topic. Examples: `retention.ms`, `cleanup.policy`, `min.insync.replicas`, `max.message.bytes`.

Set overrides at create time (`--config` on `kafka-topics`) or later with `kafka-configs`.

Describe with configs to see the effective values and whether a value is a default or a topic override.

Use overrides for special topics. Do not create a unique snowflake for every topic without a reason. A team standard plus a few overrides is easier to operate.

Some keys are not safe to change on a live topic without a plan. Read the official documentation for the key. Adding partitions is a separate operation from config alter.

Dynamic broker config also exists. That is a cluster-wide change. This section focuses on per-topic overrides.

### Questions

#### Theoretical questions

1. What is a topic config override?
2. Name four keys that teams often override.
3. How do you see whether a value is a default or an override?
4. Why is a unique config for every topic a problem?
5. What is the difference between a topic override and a dynamic broker config?

#### Easy practical tasks

1. Create a topic with `retention.ms=60000`. Describe the configs. Write the value.
2. Make a table: Key, default meaning, when you override.
3. Find `kafka-configs --help`. Write the alter flag names.
4. Write four sentences about overrides. Use only facts from this section.

#### Medium practical tasks

1. Alter `retention.ms` on an existing topic. Describe before and after.
2. Remove the override (delete config) if the tool supports it. Confirm that the default returns.
3. Read official topic config keys. Write three keys that you must not change without a plan.

#### Advanced practical tasks

1. Write a script that applies a standard set of overrides (RF is not a config key—use create for RF; include retention and min ISR). Run it on a test topic.
2. Compare two topics in a table: cleanup policy, retention, min ISR, max message bytes.

---

## Internal topics (`__consumer_offsets`, transaction topics)

Kafka creates internal topics for its own state.

`__consumer_offsets` stores committed offsets for consumer groups. Do not produce application events to this topic. Do not delete it. The group coordinator uses it.

Transaction state uses internal topics (names include transaction state in official docs, for example `__transaction_state`). EOS producers need them. Do not delete them.

Other internal names can appear (cluster metadata in KRaft uses a metadata log on controllers, which is not a normal application topic). Connect and Streams can add more topics when you use those products.

`--list` can show internal topics. Some tools hide them unless you ask. Treat every name that starts with `__` as reserved unless the documentation says otherwise.

Retention and cleanup on `__consumer_offsets` are special. Do not set random overrides on internal topics.

If `__consumer_offsets` is unhealthy, groups cannot commit. That is an operations emergency. Topic 15 covers operations.

### Questions

#### Theoretical questions

1. What does `__consumer_offsets` store?
2. Why must you not delete it?
3. What are transaction internal topics for?
4. How do you recognize many internal names?
5. Why must you not set random retention on `__consumer_offsets`?

#### Easy practical tasks

1. List topics including internal topics if your tool supports that. Write the internal names that you see.
2. Write five sentences about internal topics. Use only facts from this section.
3. Make a table: Internal topic, purpose.
4. Find official names of transaction-related internal topics.

#### Medium practical tasks

1. Describe `__consumer_offsets`. Write partition count and cleanup policy if shown.
2. Commit with a consumer. Watch group describe. Write how that relates to `__consumer_offsets` without dumping the log.
3. Read official documentation on offset storage. Rewrite the idea in STE.

#### Advanced practical tasks

1. Write a one-page "do not touch" list for internal topics: produce, delete, config, consume for curiosity.
2. On a cluster that uses transactions, list transaction-related topics and write their RF. Compare with official recommendations.

---

## ACLs (preview)

An ACL is a rule: a principal may or may not perform an operation on a resource. Example idea: user `alice` can write to topic `orders.placed`.

Kafka authorizes when an authorizer is enabled. Without an authorizer, a client that can reach the listener can do all admin and data operations. A local learning cluster often has no ACLs. A production cluster must have authentication and ACLs (topic 17).

Principals come from the authentication method (SASL user, certificate identity). Resources include topics, groups, and cluster operations. Operations include Read, Write, Create, Delete, Describe, and others.

`kafka-acls` adds and lists rules. Super users bypass ACLs. Do not run production as a super user for applications.

This section is a preview. Do not design a full security system from this page alone. Learn TLS and SASL in topic 17. Then set ACLs.

Wrong ACLs look like "timeout" or "authorization failed". Read the broker log. Do not open the listener to the world to "fix" it.

### Questions

#### Theoretical questions

1. What is an ACL?
2. What happens when no authorizer is enabled?
3. What is a principal?
4. What is a super user?
5. Why is this section only a preview?

#### Easy practical tasks

1. Write four sentences about ACLs. Use only facts from this section.
2. Find `kafka-acls --help`. Write three operations that you see.
3. Make a table: Resource type, example name, example operation.
4. Open the official authorization documentation. Write the URL.

#### Medium practical tasks

1. On a local cluster without ACLs, write why your console producer works without a user name.
2. Read a short official ACL example. Rewrite the example in STE (who, resource, operation).
3. List three production mistakes: no authorizer on a public network, app uses super user, ACL on the wrong group id.

#### Advanced practical tasks

1. If you later enable SASL on a lab (topic 17), add an ACL that allows one user to write one topic and another user to read it. Prove both. This task can wait until topic 17.
2. Write a one-page ACL sketch for a team: producers, consumers, admins, and the internal topics.

---

## `kafka-topics.sh`, `kafka-configs.sh`, `kafka-consumer-groups.sh`

These tools are the daily admin interface.

**`kafka-topics`.** Create, list, describe, delete, and alter partition count. Always pass `--bootstrap-server`.

**`kafka-configs`.** Alter and describe topic or broker configuration. Use the entity type that the help text names (`topics`, `brokers`).

**`kafka-consumer-groups`.** List groups, describe members and lag, delete a group, reset offsets for a stopped group. Reset is dangerous. Reset only when the group is idle and you have a written reason.

Other useful tools exist: `kafka-console-producer`, `kafka-console-consumer`, `kafka-get-offsets`, `kafka-reassign-partitions` (operations), `kafka-acls`. Learn them when you need them.

`--help` is the source of truth for your version. Flag names change. Old blogs show `--zookeeper`. Do not use that flag on a KRaft cluster.

Put the `bin` directory on `PATH`, or write the full path. On Windows, call the `.bat` files from `cmd` or PowerShell.

Prefer scripts in version control for repeated admin. Do not run one-off deletes from memory in production.

### Questions

#### Theoretical questions

1. What is `kafka-topics` for?
2. What is `kafka-configs` for?
3. What is `kafka-consumer-groups` for?
4. Why is offset reset dangerous?
5. Why must you ignore `--zookeeper` in old blogs?

#### Easy practical tasks

1. Run `--help` for all three tools. Write one task that each tool does.
2. Describe a group that you used in topic 5 or 6. Save the output.
3. Write five sentences about the three tools. Use only facts from this section.
4. Make a table: Tool, safe daily use, dangerous use.

#### Medium practical tasks

1. Use `kafka-configs` to add and then describe `retention.ms` on a test topic.
2. Reset a stopped test group to earliest. Consume. Write the full command (including `--execute` if required).
3. Write a small PowerShell or bash wrapper that always sets `BOOTSTRAP=localhost:9092`.

#### Advanced practical tasks

1. Write a checked-in admin script: create topic if missing, apply standard configs, describe. Do not delete.
2. Compare the command-line tools with the Admin API in your language for create and describe. Write when you use each.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the admin path for a new application topic from create to first produce.
2. How do topic overrides, internal topics, and ACLs differ in who may change them?
3. Why do all modern tools take a bootstrap address?
4. What must you check before `--delete` and before an offset reset?
5. How does this topic connect to storage (topic 8) and to groups (topic 6)?

#### Easy practical tasks

1. Create `admin.review` with 3 partitions, set a short retention, describe configs, produce one record, describe a group after consume.
2. Write a one-page cheat sheet for the three tools with the exact flag names from your `--help`.
3. List all topics. Mark application vs internal.
4. Draw who may create topics, who may produce, and who may delete in a future ACL world.

#### Medium practical tasks

1. Write a script that fails if a topic exists without `min.insync.replicas=1` on a laptop (or `=2` on a multi-broker lab).
2. Document ten steps so that a teammate can create a standard topic on your KRaft cluster.
3. Map each subsection to an official operations or admin page.

#### Advanced practical tasks

1. Use only the Admin API (no shell) to create a topic, alter retention, describe, and delete. Record the API types.
2. Write a production admin standard: no auto-create, no ZooKeeper flags, required RF, ACL preview, and a ban on delete without a ticket.
