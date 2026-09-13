# 4. Classic Building Blocks

## Description

Classic building blocks are the basic parts of networked systems. This topic covers clients and servers, processes and threads, synchronous and asynchronous calls, stateful and stateless services, idempotency, and timeouts.

These blocks appear in almost every later topic. A monolith still uses processes, timeouts, and idempotency. A distributed system uses the same blocks with more failure.

Complete Topics 1 to 3 before this topic. Read process and thread ideas in `os.topics.md` if those terms are new.

Use one term for each concept. A process is not a thread. A synchronous call is not the same as a blocking operating-system call in every language, but this handbook uses "synchronous" for "the caller waits for the result of this operation".

---

## Client and server

A client is a program that starts a request. A server is a program that waits for a request and then replies. The same process can be a client of one system and a server for another system.

The relation is about who starts the call. A browser is a client of an HTTP server. That HTTP server can be a client of a database server. The database does not call the browser.

A contract binds the client and the server. The contract includes the protocol, the operations, the data shapes, and the error rules. If the server changes the contract, clients break. Topic 8 covers APIs.

A server can serve many clients at the same time. That fact creates concurrency. The server must not mix data from two clients. Isolation is a server duty.

Trust is not equal. A server must not trust a client. The client can send invalid data. The server validates input. The client also must not trust the network. The reply can be late, lost, or altered on an open network.

Deployment can place the client and the server on different machines. The network then sits between them. Topic 11 covers the extra failures of that placement. A client and a server in one process still exist as roles (a caller and a callee).

### Questions

#### Theoretical questions

1. What is a client?
2. What is a server?
3. Can one process be both a client and a server? Explain.
4. Why must a server not trust a client?
5. What belongs in the client-server contract?

#### Easy practical tasks

1. Write five client-server pairs in a web shop (browser, API, database, mail, payments).
2. Make a table: "Program" and "Role (client, server, both)". Add six rows.
3. Draw a sequence: user, browser, API, database. Label who starts each call.
4. List four items in an HTTP contract (method, path, body, status).

#### Medium practical tasks

1. Describe a system where the API is a server to the browser and a client to two backends. Write the failure of each hop.
2. Write ten validation rules that a server applies to a "create user" request.
3. Compare an in-process function call with a client-server call. Write eight sentences.

#### Advanced practical tasks

1. Write a one-page note on reverse roles: webhooks, where your server becomes a client of a callback. Include trust and retry.
2. Design a contract sheet for a file-upload server: size limits, time limits, and error codes. Stay defensive.

---

## Process and thread (from `os.topics.md`)

A process is an instance of a program that the operating system runs. A process has its own address space. Two processes do not share memory unless they use an explicit sharing mechanism.

A thread is a unit of execution inside a process. Threads in one process share the address space. They share heap memory. They can race on the same data. See `os.topics.md` for process control, thread types, and synchronization.

A typical server is one process with many threads, or one process with many lightweight tasks (goroutines, or similar). The architecture view cares about isolation and failure. If one process stops, all threads in that process stop. If one thread deadlocks, other threads in the process can still stall on shared locks.

Multiple processes give a harder boundary. You can restart one process. You can limit memory per process. You pay a cost: you must use the network, pipes, or files to communicate.

A container or a virtual machine runs one or more processes. Those tools do not replace the process model. They package it.

Do not design twenty processes for a student project. Start with one process and clear modules (Topic 3). Add a process when you need a different lifecycle, a different resource limit, or a different failure domain.

### Questions

#### Theoretical questions

1. What is a process?
2. What is a thread?
3. What do threads in one process share?
4. What happens to threads when their process stops?
5. Why does a second process give a harder boundary?

#### Easy practical tasks

1. Write four sentences that compare process and thread. Use `os.topics.md` terms if you already read that path.
2. Make a table: "Property" and "Process vs thread". Add address space, crash scope, and communication.
3. List three ways two processes can communicate.
4. Open a process list on your machine. Write the names of five processes and one guess of their role.

#### Medium practical tasks

1. For a web API plus a background worker, write when they can share one process and when they need two processes.
2. Draw a process with three threads that share a connection pool. Mark a race and a lock.
3. Explain in eight sentences why a memory leak in one thread can kill the whole process.

#### Advanced practical tasks

1. Read the process and thread sections in `os.topics.md`. Write a one-page map from those OS ideas to server design.
2. Compare threads, OS processes, and language lightweight tasks for a 10 000-connection server. Use public facts or a measured small test on your machine.

---

## Synchronous vs asynchronous calls

A synchronous call is a call where the caller waits until the operation finishes and then continues with the result or the error. A function call in one thread is the usual example. An HTTP request where the handler waits for the database is also synchronous from the handler view.

An asynchronous call is a call where the caller does not wait for the full work to finish. The caller can continue. The result arrives later through a callback, a future, a queue, or a second message.

Asynchronous design can improve throughput. The process can start other work while it waits for I/O. Asynchronous design can also improve user time: the API accepts a job and returns `202 Accepted`. The user does not wait for a long report.

Asynchronous design adds complexity. You must store the job. You must show status. You must handle a crash after accept and before finish. Topic 10 covers messages and the outbox.

Synchronous calls are easier to understand. Use them when the work is short and the user needs the result in the same response. Use asynchronous calls when the work is long or when a burst would overload a dependency.

"Asynchronous" is not the same as "fast". A queue can add delay. "Synchronous" is not the same as "slow". A short in-process call can be the fastest path.

### Questions

#### Theoretical questions

1. What is a synchronous call in this handbook?
2. What is an asynchronous call?
3. Why can asynchronous design improve throughput?
4. What extra problems does an accepted job create?
5. When do you keep a call synchronous?

#### Easy practical tasks

1. Classify eight operations as sync or async for a shop (price check, receipt email, monthly report, card charge).
2. Write a sequence for a sync checkout and a sequence for an async report.
3. Make a table: "User sees" and "System does" for `200` with a body versus `202` with a job id.
4. List four ways a later result can return to a user (poll, webhook, email, websocket).

#### Medium practical tasks

1. Design an async PDF export: accept, store, worker, status endpoint. Write failure cases.
2. Explain in eight sentences how a sync handler that waits on a slow email SMTP can exhaust threads.
3. Choose sync or async for image thumbnail creation. Write an ADR paragraph.

#### Advanced practical tasks

1. Write a one-page comparison of async I/O in one process versus a message queue between two processes.
2. Design backpressure for an async job API when the queue length exceeds N. Include the user-visible error.

---

## Stateful vs stateless services

A stateful service stores session or business data in the process memory (or on the local disk of that instance) that the next request needs. A stateless service keeps no such data in the instance. Each request contains the data that the handler needs, or the handler loads state from a shared store.

"Stateless" does not mean "the system has no data". The data lives in a database, a cache, or a file store that all instances share. The instance can die. Another instance can serve the next request.

Stateful instances are hard to scale horizontally. A load balancer must send the same user to the same instance (sticky sessions), or the user loses the session. Sticky sessions reduce the value of extra instances.

Some state is acceptable in memory: caches, connection pools, and local counters. Those items must be rebuildable. If the process stops, a cache miss is allowed. A lost paid order is not allowed in a local-only store.

Affinity to one machine also appears with local uploads or local websocket maps. Prefer a shared store when you have more than one instance.

Choose state location as an ADR. Topic 9 covers data ownership. The source of truth must not be the memory of one web process.

### Questions

#### Theoretical questions

1. What is a stateful service?
2. What is a stateless service?
3. Why is "stateless" compatible with a database?
4. What problem do sticky sessions create?
5. Which in-memory data is acceptable, and which is not?

#### Easy practical tasks

1. Write four sentences that define stateful and stateless. Use only facts from this section.
2. Make a table: "Data" and "Where it must live". Add session login, cart, product price, and a cache of public pages.
3. List three reasons a process can stop (deploy, crash, scale in).
4. Draw one stateful shop (session in process) and one stateless shop (session in a store).

#### Medium practical tasks

1. Redesign a sticky-session chat map to a shared store. Write what you lose and what you gain.
2. Write an ADR: sessions in a database versus a signed cookie. Include security as a quality, not as attack steps.
3. Explain how a local disk upload folder fails when you add a second instance.

#### Advanced practical tasks

1. Write a one-page note on "stateless app, stateful store" and the store as the new bottleneck.
2. Design instance replacement (rolling deploy) for a service that has a 30-second in-memory buffer. How do you drain the buffer?

---

## Idempotency

Idempotency means that a second execution of the same operation has the same effect as the first execution. The client can retry. The system does not apply the business effect twice.

Read operations are often naturally idempotent. `GET /users/1` does not change data. Write operations are not naturally idempotent. Two `POST /transfers` with the same body can move money twice if the server treats them as two transfers.

Networks lose replies. A client that does not see a response does not know if the server committed the work. The client retries. Without idempotency, retries duplicate side effects.

An idempotency key is a client-generated identifier for one intended write. The server stores the key and the result. A retry with the same key returns the stored result and does not perform the write again.

Design operations as "set this state" when you can. `PUT` of a full resource is easier to make idempotent than "add 10 to the balance" without a key.

Idempotency has a time window and a storage cost. You must expire old keys. You must define what "same request" means (same key, same body, or both).

Topic 8 covers idempotency keys on APIs. Topic 10 and Topic 11 cover duplicate messages and idempotent handlers.

### Questions

#### Theoretical questions

1. What is idempotency?
2. Why do retries exist if the client already sent the request?
3. What is an idempotency key?
4. Why is "add 10" hard without a key?
5. What must you decide about the time window of a key?

#### Easy practical tasks

1. Classify eight operations as naturally idempotent or not.
2. Write a sequence: client sends pay, network drops the response, client retries. Show the duplicate without a key.
3. Make a table: "HTTP method" and "Often idempotent?". Add GET, PUT, POST, DELETE.
4. Write four rules for a client that generates idempotency keys.

#### Medium practical tasks

1. Design a server store for idempotency keys: fields, unique rule, and expiry.
2. Write three test cases: first write, retry with same key, second write with a new key.
3. Explain what the server does if the key is the same and the body is different.

#### Advanced practical tasks

1. Write a one-page design for idempotency across two processes (API and worker) for a payment.
2. Compare "natural idempotency" (set state) with "key-based idempotency". Give two operations for each style.

---

## Timeouts

A timeout is a limit on how long a caller waits for an operation. When the limit expires, the caller stops the wait and treats the call as a failure (or as unknown, if the work can still complete).

Every remote call needs a timeout. A call without a timeout can wait forever. Waiting forever holds a thread, a connection, and a user.

Timeouts must nest. If the user-facing budget is 3 seconds, an inner database call cannot have a 30-second timeout. The inner timeout must be smaller than the remaining budget.

A timeout is not a proof that the work did not run. The server can finish after the client gives up. That case is why idempotency matters. The client will retry.

Set timeouts from measures, not from fashion. Start from the P99 of a healthy dependency plus a margin. A timeout that is too short increases errors. A timeout that is too long holds resources.

Apply timeouts to HTTP clients, database clients, and locks. Also apply a timeout to the user request on the server so that abandoned clients do not keep work alive without a bound.

Topic 11 covers retries, backoff, and circuit breakers. Those patterns need timeouts first. A retry without a timeout multiplies load.

### Questions

#### Theoretical questions

1. What is a timeout?
2. Why must every remote call have a timeout?
3. Why must inner timeouts fit the outer budget?
4. Why is a timeout not proof that the work did not run?
5. How do you choose a first timeout value?

#### Easy practical tasks

1. Write a latency budget of 2 seconds with three nested calls. Assign timeouts.
2. Make a table: "Call" and "Timeout". Add DNS, HTTP, and database for one handler.
3. List four resources that a hung call holds.
4. Write five sentences on the link between timeouts and retries.

#### Medium practical tasks

1. Design timeout policy for an API that calls a payment vendor with P99 of 800 ms. Include the user-facing budget.
2. Write what the user sees when the API times out but the payment can still succeed. Link to idempotency.
3. Find default timeouts in one HTTP client library. Write if they are safe for production.

#### Advanced practical tasks

1. Write a one-page timeout standard for a team: required settings, forbidden infinite waits, and how to change a value.
2. Analyze a cascade: service A waits 30 s on B, B waits 30 s on C. Show how short timeouts plus a budget stop the cascade.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do client, server, process, and timeout combine in one HTTP request to a database?
2. Why do stateless instances still need idempotency?
3. When does an asynchronous job still use a synchronous call inside the worker?
4. How does a thread-sharing process change the meaning of "isolate this fault"?
5. Which building blocks belong in almost every ADR for an external vendor call?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic with one diagram.
2. For a login request, label client, server, process, sync call, state location, idempotency need, and timeouts.
3. Classify ten real features from an app that you use as sync or async and as stateful or stateless at the edge.
4. List default timeout values in your language HTTP and database libraries.

#### Medium practical tasks

1. Design a small photo service: upload (async), show album (sync, stateless API), and one worker process. Write timeouts and idempotency keys.
2. Write three ADRs: process split or not, session store, and HTTP client timeouts.
3. Draw a failure story: reply lost, retry, duplicate prevented. Use only terms from this topic.

#### Advanced practical tasks

1. Write a two-page "building block standard" for a team of four: required rules for new outbound calls.
2. Instrument a small local program with a short timeout and a retry. Record measurements. Do not test against systems that you do not own.
