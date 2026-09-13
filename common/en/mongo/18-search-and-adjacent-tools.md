# 18. Search and Adjacent Tools

## Description

This topic shows MongoDB tools that sit next to ordinary collections. **Atlas Search** uses Lucene for full-text search. **Time series collections** store measurements that have a time. **GridFS** stores files that exceed the document size limit. MongoDB can act as a **cache**, but **Redis** is usually the better cache. Complete querying, indexes, and aggregation first.

Use one term for each concept. A **search index** is a Lucene index that Atlas maintains. A **time series collection** is a collection type that buckets measurements. **GridFS** is a driver convention that uses `files` and `chunks` collections. A **cache** is a fast store of data that another system already owns.

Pick the tool that matches the problem. Do not use GridFS for 1 KB JSON. Do not use Atlas Search for a single equality on `_id`.

---

## Atlas Search (Lucene)

**Atlas Search** is a search engine that Atlas builds on **Apache Lucene**. You create a search index on a collection. You query with the aggregation stage `$search` (and related stages such as `$searchMeta`).

Atlas Search is not the same as a MongoDB text index (`$text`). The text index is a server feature with limited scoring and language support. Atlas Search is stronger for:

- Full-text relevance
- Fuzzy match
- Autocomplete
- Facets (with the Search facet features)
- Analyzers (how text splits into tokens)

You define mappings. Dynamic mappings index many fields. Static mappings name the fields and the analyzer. Static mappings are easier to control.

A typical pipeline:

```javascript
db.articles.aggregate([
  {
    $search: {
      index: "default",
      text: { query: "replica set", path: "body" }
    }
  },
  { $limit: 20 },
  { $project: { title: 1, score: { $meta: "searchScore" } } }
])
```

Search is eventually consistent with the collection. A write is not always in the index in the same millisecond. Design the UI for a short delay, or refresh after a known wait in admin tools.

Atlas Search runs on Atlas. Community Server does not include Atlas Search. Alternatives on self-managed stacks are an external search engine plus CDC (topic 15).

Do not create a search index on every field "just in case". Indexes cost memory and build time.

Do not use `$search` as the first idea for `{ _id: id }`. Use `find`.

Read the Atlas Search documentation for your Atlas version. Operators change.

### Questions

#### Theoretical questions

1. What library does Atlas Search use?
2. How is Atlas Search different from a `$text` index?
3. What aggregation stage runs an Atlas Search query?
4. Why can a just-written document miss a search result?
5. Can a self-managed Community Server run Atlas Search?

#### Easy practical tasks

1. Open the Atlas Search page. Write the URL.
2. Write three query types that fit Lucene search and three that fit `find`.
3. Make a table: `$text`, Atlas Search. Add three rows.
4. If you have Atlas, list whether Search is available on your tier.

#### Medium practical tasks

1. Create a search index on a test collection (or write the JSON mapping). Run one `$search`. Write the first `_id`.
2. Compare a `find` regex with `$search` fuzzy for a misspelled word. Write which result is useful.
3. Read analyzers (`lucene.standard` vs language). Write one reason to pick a language analyzer.

#### Advanced practical tasks

1. Design mappings for a product catalog: title, description, sku, facets for brand. Write static mappings in outline form.
2. Compare Atlas Search with a change-stream plus external Lucene/Elasticsearch. Write cost, consistency, and who operates the cluster.

---

## Time series collections

A **time series collection** stores documents that have a **time field** and, usually, a **meta field**. MongoDB groups documents into **buckets** on disk. The bucket layout saves space and can speed some range scans on time.

Create:

```javascript
db.createCollection("readings", {
  timeseries: {
    timeField: "ts",
    metaField: "sensorId",
    granularity: "seconds"
  }
})
```

Each measurement is still a document in the API. You insert one reading per event. You do not manage buckets by hand.

Good fit:

- Sensors
- Metrics
- Click or event streams with a clear timestamp
- Data that you often query by time range and by a device or host

Poor fit:

- Documents that you update often in place
- Highly relational operational data
- Data without a real time field

Secondary indexes are allowed with limits that depend on version. You often index the meta field and the time field.

TTL on time series collections can drop old buckets. Read the current TTL page for time series.

You cannot treat a time series collection as a normal collection for every command. Some update and delete patterns are restricted. Read the limitations page before you migrate a hot collection.

Do not store a huge unbounded array of readings inside one sensor document. That is the anti-pattern from modeling. Use a time series collection or a bucket pattern that you control.

### Questions

#### Theoretical questions

1. What two fields do you set when you create a time series collection?
2. What is a bucket in this context?
3. Why do sensors fit this collection type?
4. Why are frequent in-place updates a poor fit?
5. Why is an unbounded array of readings on one document a poor alternative?

#### Easy practical tasks

1. Open the time series collection page. Write the URL.
2. Write a sample reading document with `ts` and `sensorId`.
3. Make a table: time series, normal collection. Add three differences.
4. List three `granularity` values from the manual.

#### Medium practical tasks

1. Create a time series collection. Insert 100 readings. Run a time-range find. Write `explain` stage names if you can.
2. Read update and delete limitations. Write two forbidden or costly operations.
3. Compare the bucket pattern from topic 8 with a time series collection. Write two similarities and two differences.

#### Advanced practical tasks

1. Design retention: 7 days hot, TTL, and a monthly archive collection. Write the jobs.
2. Read compression and `granularity` advice. Pick a granularity for 1 Hz sensors and for 1-minute metrics. Explain.

---

## GridFS (large files)

The document size limit is 16 MB. **GridFS** stores a larger file as many **chunks** plus one **file** metadata document.

Default collections:

- `fs.files` — file name, length, content type, upload date
- `fs.chunks` — `{ files_id, n, data }` where `data` is a binary chunk

The default chunk size is 255 KB on many versions. Drivers implement upload and download. You do not assemble chunks by hand in ordinary application code.

```javascript
// driver APIs differ; mongosh may use a GridFSBucket helper
```

Use GridFS when:

- You must keep the file in MongoDB (simple ops, one backup domain)
- The file is larger than 16 MB or you want a stream API
- You accept database backup and cost for binary data

Do not use GridFS when:

- Object storage (S3 or similar) is available and the file is large or numerous
- You only need a URL to a CDN
- The file is a few kilobytes (store a BinData field or a URL)

GridFS files are not ordinary documents for `find` of the whole file. You query `fs.files` for metadata. You stream chunks through the driver.

Indexes: GridFS needs the indexes that the driver creates (`files_id` + `n` unique on chunks). Do not drop them.

Backup size grows with every file. A dump of GridFS is heavy. Snapshots include the binaries.

Atlas has no special "GridFS product". It is still collections. The 16 MB limit still applies to each chunk document, which is fine at 255 KB.

### Questions

#### Theoretical questions

1. What size limit makes GridFS necessary for one file?
2. What are the two GridFS collections?
3. What does one chunk document hold?
4. When is object storage a better place for the file?
5. Why must you keep the GridFS chunk index?

#### Easy practical tasks

1. Open the GridFS page. Write the URL and the default chunk size.
2. Draw `fs.files` (one row) and three `fs.chunks` for one file.
3. Make a table: BinData field, GridFS, S3. Add max size and typical use.
4. Write why a 2 KB avatar does not need GridFS.

#### Medium practical tasks

1. Upload a file larger than 1 MB with an official driver GridFS API. Download it. Compare checksums.
2. Query `fs.files` for `filename`. Write the `length` that you see.
3. Estimate document count for a 100 MB file at 255 KB chunks. Write the number.

#### Advanced practical tasks

1. Read how to use a custom bucket name and chunk size. Write when a larger chunk helps.
2. Design file storage for user uploads: which files go to S3, which (if any) go to GridFS, how you back up each.

---

## MongoDB as a cache (usually Redis is better)

A **cache** holds copies of data for fast reads. The source of truth is often another store (or the same MongoDB collection).

MongoDB can serve hot documents quickly when they fit in RAM. That is a database with a warm cache, not a purpose-built cache product.

**Redis** (or a similar in-memory store) is usually better as a cache because:

- TTL and eviction policies are first-class
- Latency is typically lower for small keys
- You can evict without a WiredTiger checkpoint story
- Many cache patterns (locks, counters, session blobs) are already documented for Redis

When MongoDB is enough without a cache:

- The working set fits in RAM
- A find by `_id` or a covered query already meets the SLA
- You do not want a second system

When a cache in front of MongoDB helps:

- A hot key that would otherwise hammer one document
- A computed aggregation that is expensive and can be seconds stale
- Session data that you do not want in the primary working set

If you still use MongoDB as a cache collection:

- Set a TTL index on `expiresAt`
- Keep documents small
- Accept that TTL is not instant to the second
- Do not treat it as the only copy of data that you cannot rebuild

Do not add Redis "because everyone does". Measure. Do not use MongoDB as a cache "to avoid Redis" if you then implement eviction poorly.

Do not store sessions in a huge MongoDB document that grows without TTL.

### Questions

#### Theoretical questions

1. What is a cache in this section?
2. Why is Redis often a better cache than MongoDB?
3. When is MongoDB without Redis enough?
4. What MongoDB feature expires cache documents?
5. Why is a cache a poor only copy of data that you cannot rebuild?

#### Easy practical tasks

1. Write five sentences: MongoDB hot working set vs Redis cache.
2. Open a Redis vs MongoDB comparison that you trust, or the Redis handbook in this project. Write one fact.
3. Make a table: session store, page fragment, product by `_id`. Pick MongoDB, Redis, or either.
4. Write a TTL index key pattern for `{ expiresAt: 1 }`.

#### Medium practical tasks

1. Implement a tiny cache collection with TTL. Read through cache then collection. Write the two paths.
2. Measure a heavy aggregation vs a cached document of its result. Write the times and the stale window.
3. List eviction problems if you only delete cache rows from the application and you never set TTL.

#### Advanced practical tasks

1. Design a cache key policy for a multi-instance API (what you include in the key, how you invalidate on write).
2. Write a team rule: when you may add Redis, when you must not, and how you prove the need with metrics.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do Atlas Search, time series collections, and GridFS each change the default "one ordinary collection" model?
2. When do `$search`, a time-range find on a time series collection, and a GridFS download solve three different user actions?
3. Why is Redis the usual cache while MongoDB still holds the source documents?
4. How does CDC from topic 15 relate to Atlas Search and to an external search engine?
5. A teammate stores 50 MB images as one document, puts click logs in an unbounded array, and adds Redis without a measurement. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: `$search` vs `$text`, timeField/metaField, `fs.files`/`fs.chunks`, TTL cache, Redis.
2. Draw a system: API, MongoDB, Atlas Search, optional Redis, object storage for big files.
3. Map four features of an app (search box, sensor chart, video, session) to a tool from this topic.
4. Open one official page from each section. Write the four URLs.

#### Medium practical tasks

1. For a blog product, write which fields use Atlas Search, which stay in `find`, and where images live.
2. Write a one-page "adjacent tools" decision table with a yes/no test for each tool.
3. If Atlas is available, create a search index or a time series collection and run one query. Write the result count.

#### Advanced practical tasks

1. Design an IoT backend: time series for readings, Search for device docs, S3 for firmware files, Redis for last-value cache. Write data flow and TTLs.
2. Write a production standard: who may create Search indexes, max GridFS usage, cache policy, and when to refuse MongoDB-as-cache.
