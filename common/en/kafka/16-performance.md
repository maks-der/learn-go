# 16. Performance

## Description

Performance in Kafka is the rate of records and the delay of each record. You can raise throughput. You can lower latency. One change often hurts the other. You measure both. You change one setting at a time.

This topic covers throughput versus latency, producer batching, consumer fetch waits, partition count versus broker count, zero-copy transfer, and client versus broker bottlenecks. Complete this topic after topics 3, 4, 5, and 15. Test on a KRaft cluster. Do not use ZooKeeper for new work.

Use one term for each concept. Throughput is records or bytes per second. Latency is the time from produce send to a defined end (ack, or consume). A bottleneck is the part that limits the rate. Zero-copy is a transfer that avoids an extra copy in user space. Batching groups records before a send.

---

## Throughput vs latency

**Throughput** is how much data the system moves in a period. Teams measure records per second or bytes per second. You must say if the number is produce, consume, or both.

**Latency** is how long one record waits. Produce latency often means time to the ack. End-to-end latency means time to a consumer process. Those two numbers are not the same.

A large batch raises throughput. The first record in the batch waits for `linger.ms` or for the batch to fill (topic 4). That wait raises latency.

A small batch or `linger.ms=0` can lower latency and lower throughput. More requests hit the broker. CPU and request queues grow.

Many partitions can raise parallelism and throughput. They can also raise metadata, memory, and end-to-end delay if each partition has little data (later section).

Do not optimize without a number. Write the target: for example, 50 000 records per second produce, p99 produce ack under 20 ms. Then measure.

KRaft metadata is not the usual produce bottleneck on a stable cluster. A saturated controller can still slow admin operations and new leadership. Measure brokers and clients first for data-plane tests.

### Questions

#### Theoretical questions

1. What is throughput in this handbook?
2. What two latency ends must you not mix?
3. How can a larger batch raise latency?
4. Why must you write a numeric target before you tune?
5. Is KRaft the first place you look for a produce bottleneck?

#### Easy practical tasks

1. Write five sentences about throughput and latency. Use only facts from this section.
2. Make a table: Change, throughput effect, latency effect.
3. Write one numeric target for a lab topic (rate and p99).
4. Find official performance or configuration notes that mention batching and latency. Write one sentence.

#### Medium practical tasks

1. Produce with a console or a small client at two linger values. Estimate records per second. Write both numbers and the method.
2. Draw a time line: records wait in a batch, send, ack, consume.
3. List three metrics you will collect (producer, broker, consumer) from official monitoring names.

#### Advanced practical tasks

1. Measure produce p50 and p99 and records per second for two batch settings on KRaft. Write a short report with limits of the test.
2. Write a one-page SLO draft: throughput, produce latency, end-to-end latency, and how you measure each.

---

## Batch size and linger

The producer groups records into batches per partition. `batch.size` is the maximum batch size in bytes (topic 4). `linger.ms` is the extra wait to add more records before send.

When the batch is full, the producer sends even if linger time is not over. When linger time is over, the producer sends even if the batch is not full.

Compression (topic 4) runs on the batch. A larger batch often compresses better. Compression uses CPU on the client and on the broker when it decompresses for some paths.

`buffer.memory` limits how much the producer can wait to send. If the buffer is full, `send` blocks or fails per configuration.

Do not set `batch.size` larger than what the broker accepts. Broker `message.max.bytes` and topic `max.message.bytes` limit the record batch.

For low latency, keep linger small. For high throughput, raise linger and batch size until the broker or the network is the limit.

Measure after each change. A huge batch can raise produce latency and can cause larger fetch delays on the consumer.

KRaft does not change these producer keys.

### Questions

#### Theoretical questions

1. What does `batch.size` limit?
2. What does `linger.ms` wait for?
3. When does the producer send before linger ends?
4. How does compression relate to batch size?
5. What broker or topic byte limit can reject a large batch?

#### Easy practical tasks

1. Write four sentences about batch and linger. Use only facts from this section.
2. Make a table: `linger.ms` 0 versus 20, expected latency, expected throughput.
3. Find `batch.size` and `linger.ms` in official producer docs. Write the defaults.
4. List three producer keys from this section and one sentence each.

#### Medium practical tasks

1. Run a producer with `linger.ms=0` and with `linger.ms=20`. Write rate and a latency sample.
2. Change `batch.size`. Keep linger fixed. Write what happens to request count if you can see client metrics.
3. Read `message.max.bytes` documentation. Write how it relates to `batch.size`.

#### Advanced practical tasks

1. Enable compression for one codec. Compare CPU and throughput at two batch sizes. Write the result.
2. Write a producer performance profile for two apps: a UI event path (low linger) and a bulk import (higher linger).

---

## Fetch min bytes / max wait

The consumer fetch request asks the broker for records. `fetch.min.bytes` is the minimum data the broker should gather before it answers (if it can). `fetch.max.wait.ms` is the maximum time the broker waits to reach that minimum.

If the partition has enough data, the broker answers at once. If the partition is idle, the broker can wait up to `fetch.max.wait.ms`. That wait raises consume latency and can raise throughput per request (fewer empty or tiny fetches).

`max.partition.fetch.bytes` and `fetch.max.bytes` limit how much one fetch returns. A too-small limit causes more round trips. A too-large limit can increase memory and pause.

These keys are consumer configuration (topic 5). They do not change produce batching.

For interactive consumers, keep wait small. For bulk sink jobs, a larger min bytes and wait can be better.

Do not set `fetch.min.bytes` so high that a low-traffic partition always waits the full `fetch.max.wait.ms` if your SLA cannot accept that delay.

KRaft does not change fetch keys. Fetch still goes to the partition leader (or a replica if you configure replica fetch).

### Questions

#### Theoretical questions

1. What does `fetch.min.bytes` ask the broker to do?
2. What does `fetch.max.wait.ms` limit?
3. When does the broker ignore the wait and answer at once?
4. What do the max fetch byte keys limit?
5. Why can a high min bytes hurt a low-traffic topic?

#### Easy practical tasks

1. Write five sentences about fetch wait. Use only facts from this section.
2. Make a table: Setting pair, better for latency or throughput.
3. Find the three fetch keys in official consumer docs. Write the defaults.
4. Draw: consumer fetch request → broker wait → response.

#### Medium practical tasks

1. Consume a quiet topic with two `fetch.max.wait.ms` values. Write the delay you feel or measure.
2. Consume a busy topic. Write whether min bytes still adds delay.
3. Compare fetch wait with producer linger in a six-sentence note: who waits, where.

#### Advanced practical tasks

1. Tune a sink consumer for throughput: raise min bytes and wait, measure records per poll and end-to-end delay.
2. Write a consumer profile table for two group types: interactive API versus Connect sink.

---

## Number of partitions vs broker count

Partitions give parallelism (topic 3). Produce, consume, and Streams tasks scale with partitions (and with group members).

Each partition has a leader, replicas, files, and memory. Too many partitions on one broker raise:

- heap and metadata size
- open files and disk seeks
- controller and election work
- time for URP to clear after a restart (topic 15)

A small cluster with tens of thousands of partitions is a common failure mode. Start with a partition count that matches expected consumers and throughput. Add partitions later if you must. Shrinking is hard (topic 3).

Broker count must support the replica factor and the leadership load. Three brokers with RF=3 can hold many partitions, but each broker leads about one third. If you add brokers, you must reassign (topic 15) or new partitions will land on new brokers only for new topics.

More partitions do not always raise throughput. If each partition has a tiny produce rate, you pay overhead and you get no extra disk sequential benefit.

KRaft can handle more partitions than old ZooKeeper-based Kafka in many versions. That is not a reason to create unbounded partitions. Read the current Apache guidance for partition limits on your version.

### Questions

#### Theoretical questions

1. What cost does each partition add on a broker?
2. Why can too many partitions slow a rolling restart?
3. Why does adding brokers not move old partitions by itself?
4. When do extra partitions not raise throughput?
5. Why is KRaft not a license for unbounded partition counts?

#### Easy practical tasks

1. Write four sentences about partitions versus brokers. Use only facts from this section.
2. Make a table: Partition count, expected consumers, broker count in a lab plan.
3. Describe a topic. Write partition count and who leads each partition.
4. Find official notes on partition limits or "too many partitions". Write three facts.

#### Medium practical tasks

1. Create two topics with 3 and 30 partitions on the same KRaft lab. Produce the same total records. Compare produce time and broker metrics if you can.
2. Add a broker (or pretend with a diagram if you have one node). Write how you would reassign leadership.
3. Calculate a rough leadership count per broker: partitions × topics / brokers. Write the number for your lab.

#### Advanced practical tasks

1. Read a current Apache or Confluent article on partition count. Write a max partitions per broker that you will use as a lab rule, with the source.
2. Write a capacity note: target consumers, RF, broker count, partition count, and a reassignment step when you scale brokers.

---

## Zero-copy transfer (awareness)

**Zero-copy** means the broker can send data from the page cache (or disk) to the network socket without an extra copy into a user-space buffer. On Linux, this often uses `sendfile` or an equivalent path.

Kafka can use this path for common fetch responses of already-compressed batches. The consumer receives the same bytes that the producer sent (for that batch format).

This section is **awareness**. You do not enable a "zero-copy switch" in daily client code. You must not fight the path: avoid unnecessary decompression on the broker, keep page cache healthy (topic 15), and keep fetch sizes reasonable.

Some paths cannot use zero-copy. Examples: the broker must transcode, decrypt, or filter bytes. Security processing (topic 17) can add copies. Treat that as a cost.

Zero-copy does not remove network limits. It reduces CPU copies on the broker.

Do not claim zero-copy in a design review unless you know the fetch path that you use. Read current Kafka documentation if you need the exact conditions.

KRaft is not part of the fetch zero-copy path. Zero-copy is a data-plane broker feature.

### Questions

#### Theoretical questions

1. What does zero-copy mean in this handbook?
2. Why is this section "awareness"?
3. What OS feature often implements this path?
4. What can prevent zero-copy?
5. Does KRaft provide zero-copy?

#### Easy practical tasks

1. Write five sentences about zero-copy. Use only facts from this section.
2. Make a table: Path, extra copy likely (yes, no, or maybe).
3. Find a mention of sendfile or zero-copy in Kafka docs or a design note. Write the URL.
4. Draw: page cache → socket → consumer. Label "no extra user copy" on the broker.

#### Medium practical tasks

1. Compare a fetch of compressed batches with a path that needs broker-side filter (from docs). Write which is closer to zero-copy.
2. Explain in six sentences how page cache (topic 15) and zero-copy work together.
3. List two security features (topic 17) that can add CPU on the broker.

#### Advanced practical tasks

1. Read a current Kafka performance write-up. Rewrite the zero-copy paragraph in STE. Add one sentence on when it does not apply.
2. Write an awareness note for your team: what you do not tune, and what you still measure (CPU, network).

---

## Client vs broker bottlenecks

A **client bottleneck** is in the producer or consumer process: CPU (serialize, compress), wait on `send`, slow poll loop (topic 5), small thread pool, or a single-thread client that cannot fill the network.

A **broker bottleneck** is in the Kafka process or its machine: disk I/O, network card, CPU (request handlers, TLS), heap or GC (topic 15), too many partitions, or a saturated request queue.

A **network bottleneck** sits between them: bandwidth, cross-zone cost, or a load balancer that you should not put in the data path without a design.

How you tell:

- High producer I/O wait and low broker disk: client or network.
- High broker disk utilization and rising produce latency: broker storage.
- Consumer lag with idle broker and idle CPU on consumer: often the application after poll (topic 19 backpressure).
- All brokers busy, one client idle: scale or fix that client.

Do not raise partitions or linger as the first step. Measure client metrics and broker metrics. Change one layer.

Managed Kafka (MSK, Confluent Cloud) hides some broker metrics. You still measure the client. You still watch the vendor dashboard for broker saturation.

KRaft controller saturation shows up in admin slowness and failover time, not in every produce on a healthy leader.

### Questions

#### Theoretical questions

1. What is a client bottleneck?
2. What is a broker bottleneck?
3. How can a slow poll loop look like a Kafka performance problem?
4. Why is "add partitions" a poor first step?
5. Where do you look first on a managed Kafka service?

#### Easy practical tasks

1. Write five sentences about client versus broker limits. Use only facts from this section.
2. Make a table: Symptom, likely layer, first metric.
3. List four producer metrics and four consumer metrics from official client docs.
4. Draw three boxes: client, network, broker. Write one failure in each.

#### Medium practical tasks

1. Create a lab: a consumer that sleeps in the poll loop. Watch lag. Write why the broker is not the cause.
2. Create a lab: a tiny `batch.size` and many small requests. Watch broker request rate. Write what you see.
3. Read official monitoring pages. Map five metrics to client or broker.

#### Advanced practical tasks

1. Load-test produce until something breaks. Write the first resource that hit 80 percent (CPU, disk, network, or client). Use KRaft.
2. Write a performance debug checklist: seven steps from client metrics to broker disk, including "do not add partitions first".

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do linger, fetch wait, and batch size work together on one record path from produce to consume?
2. When do you raise partitions, and when do you fix a client bottleneck instead?
3. How do page cache, zero-copy, and disk I/O form one broker story?
4. What performance work is different on KRaft versus old ZooKeeper clusters, and what is the same?
5. How do topics 4, 5, and 15 feed this topic?

#### Easy practical tasks

1. Write a one-page cheat sheet: throughput versus latency, linger, fetch wait, partitions, zero-copy, bottleneck layers.
2. Draw the full path with the wait points (linger, broker handle, fetch wait, poll).
3. Bookmark official producer, consumer, and operations monitoring pages.
4. Write one sentence each for every configuration key named in this topic.

#### Medium practical tasks

1. Run one produce/consume lab on KRaft with two profiles (low latency, high throughput). Fill a table of settings and results.
2. Map each subsection to one official URL.
3. Take a fictional incident: high lag, low broker CPU. Write a five-step debug from this topic.

#### Advanced practical tasks

1. Produce a short performance report for your lab cluster: method, numbers, bottleneck, next change. No ZooKeeper.
2. Write a team standard: SLOs, default linger and fetch values, partition policy, and a ban on tuning without metrics.
