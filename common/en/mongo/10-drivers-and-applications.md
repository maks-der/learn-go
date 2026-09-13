# 10. Drivers and Applications

## Description

This topic shows how applications talk to MongoDB. You use an official driver. You configure a connection string. You manage `MongoClient` for the process life cycle. You set timeouts. You learn trade-offs of object mapping libraries. You never run an unbounded `find` in an API. Complete CRUD and querying first.

Use one term for each concept. A **driver** is the official client library for a language. A **connection string** (URI) is the configuration string. **`MongoClient`** is the long-lived client object. An **ODM** or mapper maps documents to language types. Do not open a new client for every request.

---

## Official drivers (Node, Python, Go, Java)

MongoDB publishes official drivers. Common drivers:

- Node.js (`mongodb`)
- Python (`pymongo`)
- Go (`go.mongodb.org/mongo-driver`)
- Java (`mongodb-driver-sync` and the reactive driver)

Other official drivers exist (C#, Rust, C). Use the driver for your language from [https://www.mongodb.com/docs/drivers/](https://www.mongodb.com/docs/drivers/).

The official driver:

- Encodes and decodes BSON
- Manages a connection pool
- Implements auth, TLS, and retryable writes
- Matches server commands

A community wrapper can sit on top of the official driver. The official driver must still be in the stack.

Do not use an abandoned third-party client that speaks a private protocol.

Driver versions have a compatibility matrix with server versions. Read the matrix before you upgrade.

The API names differ. The ideas are the same: client, database, collection, filter, update, cursor.

CRUD in a driver uses the same operators as `mongosh`. You pass language objects or BSON documents, not a SQL string.

Install the driver with the language package tool (`npm`, `pip`, `go get`, Maven). Pin a major version.

### Questions

#### Theoretical questions

1. What work does an official driver do besides "open a socket"?
2. Why must a wrapper still use the official driver?
3. Why do you read a compatibility matrix?
4. How does a driver filter compare to a `mongosh` filter?
5. Why do you pin a major version?

#### Easy practical tasks

1. Open the drivers page. Write the install command for your language.
2. Write the import or package name of the official driver.
3. Find the compatibility matrix for your driver. Write the server versions that it lists for the current driver.
4. List four driver objects: client, database, collection, cursor. Use the names in your language.

#### Medium practical tasks

1. Write a 15-line program that pings the server with the official driver. Run it.
2. Compare method names: `insertOne` in `mongosh` vs your driver. Write the two names if they differ.
3. Read the driver changelog for one major version. Write one breaking change.

#### Advanced practical tasks

1. Compare two official drivers (for example Go and Node) on BSON date types. Write how each language type maps.
2. Read about retryable writes in the driver. Write one operation that retries and one that does not.

---

## Connection string and options

A connection string is a URI.

Standalone:

```text
mongodb://localhost:27017
```

With auth and a database name for auth:

```text
mongodb://user:pass@localhost:27017/?authSource=admin
```

Atlas often uses SRV:

```text
mongodb+srv://user:pass@cluster0.example.mongodb.net/
```

`mongodb+srv` uses DNS to find hosts. It also implies TLS in the usual Atlas case.

URI options are query parameters. Examples:

- `retryWrites=true`
- `w=majority`
- `appName=myapp`
- `maxPoolSize=20`

The driver also accepts options in code. Do not set the same option in two places with different values.

Do not put the password in source control. Use an environment variable or a secret store.

Escape special characters in the user name and the password.

`appName` appears in server logs. Set it. It helps operations.

Read the current URI options page. Options change. Some Atlas options live in the Atlas UI (IP list, user roles) and not only in the URI.

### Questions

#### Theoretical questions

1. What is the difference between `mongodb://` and `mongodb+srv://`?
2. What is `authSource`?
3. Why must the password stay out of git?
4. Why set `appName`?
5. What happens if the URI and the code set different pool sizes?

#### Easy practical tasks

1. Connect with a URI from an environment variable. Do not print the password.
2. Add `appName` to the URI or the options. Write the value that you used.
3. Open the connection-string page. Write three option names.
4. Copy the Atlas URI from the UI (if you use Atlas). Mark the user, the host, and the options.

#### Medium practical tasks

1. Set `maxPoolSize` to a small number. Run a few concurrent finds. Write what you observe (wait or error).
2. Compare a URI with `w=majority` to the default. Write the write-concern object that the driver uses (log or debug).
3. Escape a password that contains `@` or `/`. Show a failed parse and a correct escape.

#### Advanced practical tasks

1. Read SRV and TXT records for an Atlas cluster (`nslookup` or `dig`). Write how they relate to the URI (no secrets in notes that you share).
2. Build the client options fully in code with no options in the URI except host. Write why a team does that.

---

## `MongoClient` lifecycle

Create **one** `MongoClient` (or the equivalent name) when the process starts. Reuse it for all requests. Close it when the process stops.

The client owns a **pool** of connections. A new client per request creates many connections. The server will reject or slow down. The application will leak sockets.

Typical sequence:

1. Read the URI from the environment.
2. Create the client.
3. Ping once at startup (fail fast if the URI is wrong).
4. Store the client in application state.
5. Use `client.Database("app").Collection("orders")` per operation (these objects are cheap).
6. On shutdown, call `Disconnect` or `Close` as the driver documents.

Do not close the client after each query.

Database and collection objects do not each open a new pool. They use the client.

In tests, you can create a client per test suite, not per test method, unless you need isolation.

Serverless functions are harder. A global client can be reused across warm invocations. A new client on every cold start is normal. Do not create a client on every single request inside a warm instance.

### Questions

#### Theoretical questions

1. Why do you create one client per process?
2. What does the client pool hold?
3. When do you close the client?
4. Why are `Database` and `Collection` objects cheap?
5. What is the risk in a serverless function that creates a client per request?

#### Easy practical tasks

1. Create a client, ping, list collection names, disconnect. Write the four calls in your language.
2. Find the method name for shutdown in your driver docs.
3. Write where you store the client in a small web app (package variable, struct field, or DI).
4. Draw a diagram: process start → client → pool → `mongod`.

#### Medium practical tasks

1. Intentionally create a client per request in a loop of 100. Watch connection count if you can (serverStatus or Atlas). Then fix it.
2. Add a graceful shutdown that closes the client after the HTTP server stops.
3. In tests, share one client. Write how you isolate data (unique prefix or a test database).

#### Advanced practical tasks

1. Read pool options: min, max, idle time. Write values for a small API (guess, then read a guide).
2. Compare connection count with one client vs a client per request using `serverStatus.connections`. Write the two numbers.

---

## Timeouts

A timeout stops a wait that is too long. Without timeouts, a request can hang.

Common timeouts (names vary by driver):

- **Server selection timeout** — how long the driver waits to find a suitable server
- **Connect timeout** — how long a new TCP/TLS connect may take
- **Socket timeout** or **network timeout** — how long a read or write on a socket may take
- **Operation timeout** — a per-command deadline (`timeoutMS` or `context` with a deadline)

Go and some modern drivers use a **context deadline** for each operation. Set a deadline in the API handler. Pass the context into the driver.

Do not use a 30-minute socket timeout as your only limit on a user request. The user is gone. The pool connection stays busy.

Do not use a 1-millisecond timeout on a query that must read disk. You will see random errors.

Retry: the driver can retry some reads and writes. A short timeout plus many retries can still overload the server. Set one overall deadline.

Atlas and `mongod` also have idle timeouts and cursor timeouts. Those are not a replacement for application deadlines.

Log timeout errors as timeouts. Do not log them as "not found".

### Questions

#### Theoretical questions

1. What problem does a timeout solve?
2. What is server selection timeout?
3. Why does a context deadline belong in an API handler?
4. Why is a very large socket timeout a risk?
5. Why is a timeout error not the same as an empty find?

#### Easy practical tasks

1. Find the timeout option names in your driver. Write three names.
2. Set a short server-selection timeout and point the URI at a closed port. Record the error.
3. In Go, or in a driver with context, pass a 2-second deadline into `Find`. Write the call.
4. Open the driver timeout page. Write the URL.

#### Medium practical tasks

1. Compare a failed connect (bad host) with a slow query. Write which timeout applies to each.
2. Set `timeoutMS` (or equivalent) on one command. Confirm the error when you `sleep` in an aggregation `$function` or a large sort (use a test).
3. Write a policy: default 5s for user reads, longer for a report job. Put the numbers in config.

#### Advanced practical tasks

1. Read about `timeoutMS` in recent MongoDB versions. Write how it relates to driver-side timeouts.
2. Trace one hung request without a deadline (in a test). Then add a deadline. Write the two outcomes.

---

## Object mapping libraries (Mongoose and others) — trade-offs

An **ODM** (object-document mapper) maps language objects to documents. **Mongoose** is a common ODM for Node. Other languages have ODMs and codegen tools.

Benefits:

- A schema in the application
- Casts and defaults
- Hooks (before save)
- A style that looks like a model class

Costs:

- Extra layer on top of the driver
- Magic (implicit casts, implicit collections)
- Version lag behind official driver features
- Hidden N+1 or unbounded queries if you use convenience helpers
- A second schema besides `$jsonSchema`

Use an ODM when the team wants models and the ODM is active. Use the official driver when you want clear BSON and control.

You can use a thin mapper (structs and tags in Go, dataclasses in Python) without a full ODM.

Do not mix two ODMs in one service. Do not ignore the official driver docs when the ODM wraps it.

Mongoose schemas are application validation. They are not a replacement for unique indexes or server validators.

Test the BSON that the ODM writes. Some libraries store dates or numbers in types that you did not expect.

### Questions

#### Theoretical questions

1. What does an ODM map?
2. Name one benefit and one cost of Mongoose-style models.
3. Why is an ODM schema not enough as the only uniqueness control?
4. When is the official driver enough?
5. Why must you inspect the BSON that the ODM writes?

#### Easy practical tasks

1. Read the Mongoose or your ODM home page. Write two features.
2. Write the same insert with the official driver and with an ODM (or a struct mapper). Compare the stored document.
3. List three hooks or middleware that your ODM advertises. Write one risk of hooks.
4. Make a two-column table: "Official driver" and "ODM". Add four rows.

#### Medium practical tasks

1. Show an implicit cast (string to ObjectId or string to number) in an ODM. Write whether you want that.
2. Disable a convenience that loads a full collection. Write the safe API that you use instead.
3. Add a server `$jsonSchema` and an ODM schema. Show a write that one allows and the other rejects.

#### Advanced practical tasks

1. Compare Mongoose, Prisma (or another mapper), and the official Node driver for one endpoint. Write lines of code and control.
2. Read about ODM change streams or transactions support. Write one feature that lagged the official driver in the past.

---

## Avoiding unbounded `find()` in APIs

An API must not return an unbounded collection.

Bad:

```javascript
const all = await col.find({}).toArray()
return all
```

This pattern:

- Grows with the table
- Uses RAM on the app and the client
- Can time out
- Can exceed response size limits

Good:

- Require a filter that is selective, or
- Require `limit` with a maximum (for example 100), and
- Use a sort and a page token or cursor, and
- Project only the fields that the client needs

`find({})` in an admin tool can be acceptable with a forced limit.

Aggregation without `$limit` on a user path has the same problem.

ODM helpers such as `Model.find()` with no limit are the same bug.

Write a shared helper: `findPage(filter, { limit, cursor })`. Reject `limit > 100`.

Test with more documents than the limit. The test fails if the API returns all rows.

Logs: record `limit` and `returned`. Alert if `returned == max` often; the client may need a better filter.

### Questions

#### Theoretical questions

1. Why is `find({}).toArray()` unsafe in an API?
2. What four controls belong on a list endpoint?
3. Why is a maximum limit required even when the client sends `limit`?
4. How is unbounded aggregation the same class of bug?
5. Why test with more documents than the page size?

#### Easy practical tasks

1. Write a bad handler and a good handler in your language (15 lines each). The good handler uses `limit`.
2. Insert 25 documents. Call a `limit(10)` find. Write the array length.
3. Add a projection to the good handler. Write the keys.
4. List three official or team rules that forbid unbounded finds.

#### Medium practical tasks

1. Implement cursor pagination with `_id` or `createdAt`. Write the query for the next page.
2. Reject `limit=100000` in the API. Write the error.
3. Find one ODM call that defaults to no limit. Wrap it.

#### Advanced practical tasks

1. Load-test a bad endpoint vs a good endpoint with 100 000 documents. Write time and memory (even if the bad one fails).
2. Design an export job that must read all documents. Write why it is a job with a cursor, not a public API.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from process start to one `find`: URI, client, timeout, collection, limit.
2. How do official drivers, ODMs, and server validators share responsibility?
3. Which configuration belongs in the URI, which belongs in code, and which belongs in Atlas or `mongod`?
4. Why do pool lifecycle and unbounded finds both cause production incidents?
5. A teammate opens `MongoClient` inside a request and calls `find({})`. Which facts do you use in the review?

#### Easy practical tasks

1. Write a small program: env URI, one client, ping, `find` with `limit(5)`, disconnect.
2. Write a cheat sheet: driver names, URI forms, client lifecycle, timeout types, ODM trade-off, API limit.
3. Set `appName` and a 5-second operation deadline. Run a ping.
4. Document the shutdown hook for your web framework in five steps.

#### Medium practical tasks

1. Build a tiny HTTP list endpoint that requires `limit` ≤ 50 and a projection. Test it with 80 documents.
2. Compare official-driver insert vs ODM insert of the same struct. Dump BSON types of dates and decimals.
3. Break the URI (bad password). Show a fast fail at startup ping, not on the first user request.

#### Advanced practical tasks

1. Measure pool connections under 50 concurrent requests with one client. Write `maxPoolSize` and peak connections.
2. Write a production checklist of ten items from this topic. Apply it to a sample app and mark pass or fail.
