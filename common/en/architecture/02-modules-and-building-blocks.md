# 2. Modules and Building Blocks

## Description

A module is a part of a system with a boundary and a responsibility. Classic building blocks are the basic parts of networked systems: clients, servers, sync and async calls, state, idempotency, and timeouts. This topic covers cohesion, coupling, encapsulation, dependency direction, and packages as boundaries.

Good modules let you change one part without a rewrite of the other parts. Bad modules spread one change across the whole system. Topic 1 called that effect maintainability.

Complete Topic 1 before this topic. Write small modules first. Do not start with many network services. A module can live in one process. Read process and thread ideas in `os.topics.md` if those terms are new.

Use one term for each concept. A module is a design boundary. A package is a language boundary. A process is not a thread. A synchronous call is not the same as a blocking operating-system call in every language. This handbook uses "synchronous" for "the caller waits for the result of this operation".

---

## Cohesion, coupling, encapsulation

Cohesion is how strongly the elements inside a module belong together. A module with high cohesion does one job. Examples: "price an order", "hash a password", "write an audit line". A module with low cohesion mixes unrelated jobs. Example: "prices, emails, and PDF layout in one file".

Coupling is how strongly one module depends on another module. A module with high coupling cannot change without changes in many other modules. A module with low coupling depends on a small, stable interface.

Aim for high cohesion and low coupling. Those two goals work together. When a module does one job, other modules need less knowledge of its internals.

Some coupling is necessary. A checkout module must call a price module. That coupling is acceptable when the price interface is small and stable.

Hidden coupling is worse than visible coupling. Hidden coupling includes a shared global variable, a shared database table with no owner, and a hidden assumption about file format. Visible coupling is an import or a documented interface.

Encapsulation is the rule that a module hides its internals and shows only a small interface. Internals include data structures, SQL, and helper functions. Callers must not depend on those internals.

Encapsulation protects invariants. An invariant is a rule that must stay true. Example: "an order total equals the sum of line totals." If every caller writes to the same tables, the invariant breaks.

Encapsulation is not the same as "make every field private and add getters". A getter that exposes the full internal structure is a leak. Hide the structure. Expose operations that match the job: `placeOrder`, not `getOrderMap`.

A change request is a test of the design. If a small business change touches many modules, cohesion is low or coupling is high.

### Questions

#### Theoretical questions

1. What is cohesion?
2. What is coupling?
3. Why do high cohesion and low coupling work together?
4. What is encapsulation?
5. What is hidden coupling?

#### Easy practical tasks

1. Write five sentences that define cohesion, coupling, and encapsulation. Use only facts from this section.
2. Classify six folders in a project that you know as high or low cohesion. Give one reason each.
3. Make a table: "Dependency" and "Visible or hidden". Add four examples.
4. List three change requests for a shop. For each request, guess how many modules you must touch.

#### Medium practical tasks

1. Split a low-cohesion module (auth plus email plus reports) into three modules. Write the new names and the remaining calls.
2. Draw two module diagrams of the same system: one with high coupling, one with lower coupling. Label the interfaces.
3. Take a public small repository. Find one file that mixes two jobs. Write a split plan of ten lines.

#### Advanced practical tasks

1. Write a one-page review of coupling types (content, common, control, stamp, data) with one software example each. Keep the language simple.
2. Measure a codebase with a simple count: incoming and outgoing dependencies per package. Interpret the two worst packages.

---

## Dependency direction

Dependency direction is the rule that says which module may import or call which other module. A healthy graph has a clear direction. A cycle is a path that returns to the start. Cycles make change hard. You cannot understand or test one module alone when it depends on a cycle.

Point dependencies toward stable modules. A domain rule is more stable than an HTTP handler. A store adapter is less stable than the domain interface that it implements. Topic 3 uses this rule in layers and hexagon design.

Allowed examples:

- A presentation module may call a domain module.
- A domain module must not import a SQL driver.
- A utility that many modules import must stay small and free of business rules.

Forbidden examples:

- Two modules that import each other.
- A domain module that imports a web framework type.
- A "utils" package that imports every other package.

A cycle can hide in data. Two modules that share a table with no owner form a data cycle. The import graph can look clean. The change cost is still high.

Draw the graph. Tools can list imports. You can also draw four boxes and arrows on paper. If you cannot draw the graph, the direction is not clear.

When you need a callback from a stable module to an unstable module, use an interface that the stable module owns. The unstable module implements the interface. That is inversion of the dependency.

### Questions

#### Theoretical questions

1. What is dependency direction?
2. Why is a cycle a problem?
3. Which way must a domain module not depend?
4. How can a shared table hide a cycle?
5. What is inversion of a dependency in one sentence?

#### Easy practical tasks

1. Draw four modules with allowed arrows only. Write why each arrow is allowed.
2. Make a table: "Import" and "Allowed? (yes/no)". Add six rows for a shop.
3. List three signs of a cycle in daily work (slow change, hard tests, or similar).
4. Write five sentences that define a healthy dependency graph. Use only facts from this section.

#### Medium practical tasks

1. Take a small repository. Draw the package import graph. Mark one cycle or write that none exists.
2. Break a paper cycle between `orders` and `mail` with an interface that `orders` owns.
3. Write a lint rule list of five import rules for a student monolith.

#### Advanced practical tasks

1. Write a one-page standard: how a new teammate adds a dependency. Include a reject example.
2. Compare "modules inside layers" with "layers inside modules" as two maps of direction. Draw both.

---

## Client/server, sync vs async, stateful vs stateless

A client is a program that starts a request. A server is a program that waits for a request and then replies. The same process can be a client of one system and a server for another system.

The relation is about who starts the call. A browser is a client of an HTTP server. That HTTP server can be a client of a database server. The database does not call the browser.

A contract binds the client and the server. The contract includes the protocol, the operations, the data shapes, and the error rules. A server must not trust a client. The client can send invalid data. The server validates input.

A synchronous call is a call where the caller waits for the result of this operation. An asynchronous call is a call where the caller does not wait for the full result in the same wait. The caller can continue and can receive a result later (a message, a poll, or a callback).

Synchronous calls are simple to read. They hold a resource while they wait. Asynchronous calls can absorb a burst. They need a place to store the work and a rule for failure.

A stateful service stores client-specific data in the process memory (or on the local disk of that instance). The next request must reach the same instance, or the data must be shared. A stateless service stores client-specific data outside the instance (a database or a cache). Any healthy instance can handle the next request.

Stateless instances are easier to clone behind a load balancer. Sticky sessions make scale harder. Some state is necessary (the database). The goal is to keep instance-local session state small.

A process is an instance of a program that the operating system runs. A thread is a unit of execution inside a process. Threads in one process share memory. If one process stops, all threads in that process stop. See `os.topics.md`.

Do not design twenty processes for a student project. Start with one process and clear modules. Add a process when you need a different lifecycle, a different resource limit, or a different failure domain.

### Questions

#### Theoretical questions

1. What is a client? What is a server?
2. What is a synchronous call in this handbook?
3. What is an asynchronous call in this handbook?
4. How does a stateless service differ from a stateful service?
5. Why do sticky sessions make scale harder?

#### Easy practical tasks

1. Write five client-server pairs in a web shop (browser, API, database, mail, payments).
2. Make a table: "Call" and "Sync or async". Add six rows.
3. Draw a sequence: user, browser, API, database. Label who starts each call.
4. List four items that belong in a client-server contract.

#### Medium practical tasks

1. Describe a system where the API is a server to the browser and a client to two backends. Write the failure of each hop.
2. Compare an in-process function call with a client-server call. Write eight sentences.
3. Design a cart: in-memory on one instance versus rows in a database. Write how you add a second API instance in each case.

#### Advanced practical tasks

1. Write a one-page note on reverse roles: webhooks, where your server becomes a client of a callback. Include trust and retry.
2. Write a state map for a login session: what lives in the browser, what lives in the API process, and what lives in a store.

---

## Idempotency and timeouts

Idempotency is the property that a repeated request has the same effect as one request. A second `PUT` of the same resource must not create a second resource. A second "charge this payment key" must not charge twice.

Networks retry. Users double-click. Workers replay messages. If the handler is not idempotent, a retry creates a defect that is worse than a timeout.

An idempotency key is a client identifier for one intended action. The server stores the key and the result. A retry with the same key returns the stored result and does not run the side effect again.

`GET` is usually idempotent. `PUT` and `DELETE` can be idempotent if you design them that way. A bare `POST` that always creates a new row is not idempotent unless you add a key or a natural unique rule.

A timeout bounds the wait. Every remote call needs a timeout. A call without a timeout can wait forever. Forever waits fill threads and hide failure.

Rules for timeouts:

- Set a timeout on each outbound call (HTTP, database, broker).
- Keep the total wait inside the user budget.
- A timeout is a local decision. The remote work can still complete.
- After a timeout, do not assume that the work did not run. Use idempotency for a retry.

A timeout is not a retry. A retry repeats the call. Retry only idempotent operations, or operations that you protect with a key. Topic 8 covers retries in a distributed system.

Do not use `sleep` as a timeout design. Do not wait without a bound "until the server answers".

### Questions

#### Theoretical questions

1. What is idempotency?
2. What is an idempotency key?
3. Why must a retry assume that the first call may have completed?
4. What does a timeout bound?
5. Why is a timeout not the same as a retry?

#### Easy practical tasks

1. Classify eight HTTP operations as idempotent, not idempotent, or "can be if designed".
2. Write five sentences that define idempotency and timeouts. Use only facts from this section.
3. List four places a student API must set a timeout (database, HTTP client, and similar).
4. Make a table: "Event" and "Need an idempotency key? (yes/no)". Add payment charge, catalog read, and mail send.

#### Medium practical tasks

1. Design `POST /payments` with an `Idempotency-Key` header. Write what the server stores and what it returns on a replay.
2. Write a timeout budget for a page that calls three backends. Show that the sum can exceed the user budget.
3. Explain in eight sentences why a double-click on "Place order" needs a unique rule.

#### Advanced practical tasks

1. Write a one-page idempotency standard: key format, storage time, and conflict when two different bodies share one key.
2. Design a failure matrix: timeout on payment versus timeout on email. Write user-visible behavior for each.

---

## Packages as boundaries

A package is a language boundary. In many languages a package (or a module, or a namespace) groups types and functions. The compiler or the runtime can hide unexported names. That hide is a tool for encapsulation. It does not replace design.

Use packages as the first module boundary in a monolith. One package does one job. Other packages import the public names only. They must not import internals. They must not share a bag named `utils` that grows without an owner.

Rules that keep package boundaries useful:

- One job per package.
- A small public surface.
- No import cycles.
- Tables and files that the package owns have a clear owner name.
- Tests can import the package without booting the whole world when you can.

A folder is not a boundary if every file is public and every package imports every other package. Names on folders without rules are decoration.

A library is a distributable package (or set of packages). A library does not own a production data store. A service does. Topic 9 defines a service.

Language features help: private names, packages, and modules. A public package that exports twenty types has a large surface. A large surface is weak encapsulation.

When you change internals and you do not change the public names, callers continue to work. That is the practical test of a package boundary.

Start with packages. Add processes and networks later. Topic 3 keeps those packages inside one deploy unit.

### Questions

#### Theoretical questions

1. What is a package in this handbook?
2. Why is a folder not automatically a boundary?
3. How is a library different from a service?
4. What is the practical test of a package boundary?
5. Why is a large public surface weak encapsulation?

#### Easy practical tasks

1. Draw a monolith with four packages and allowed imports only.
2. Make a table: "Name" and "Job". Add six package names for a blog.
3. List five names that must stay unexported in a `pricing` package.
4. Write four sentences that explain why `utils` is a risk.

#### Medium practical tasks

1. Split a single-folder homework into three packages on paper. Write the public functions.
2. Write a package lint list (five rules) that protects module boundaries.
3. Compare two public repositories: one with clear packages and one with a flat folder. Write eight sentences.

#### Advanced practical tasks

1. Write a one-page standard: how a new teammate adds a package to the monolith.
2. Review a public monolith repository. Mark packages that exist as real boundaries and folders that are only names. Write evidence.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do cohesion, coupling, and encapsulation work together in one module?
2. Why can a clean import graph still have high coupling through a shared table?
3. When do you keep a call synchronous, and when do you make it asynchronous?
4. How do timeouts and idempotency protect a client that retries?
5. Why are packages the first boundary, before processes and services?

#### Easy practical tasks

1. Write a one-page cheat sheet: cohesion, coupling, encapsulation, direction, client/server, sync/async, state, idempotency, timeout, package rules.
2. For a to-do API, name four packages, two synchronous calls, one timeout, and one idempotent write.
3. Draw a client, a server, and a database. Mark who waits and who stores session state.
4. Bookmark `os.topics.md`. Write three terms that this topic uses from that path.

#### Medium practical tasks

1. Write a short design for a library loan API in one process: packages, one stateful store, and timeouts on the store.
2. Take a teammate design that uses global variables for the cart. Rewrite it as a package with an interface and a store.
3. Make a change-request table: five small business changes and the packages that each change must touch.

#### Advanced practical tasks

1. Write a one-page review of a small public API client: coupling, timeouts, and idempotency of retries. Stay defensive.
2. Design a module map for a campus shop that can later extract one package. Write what you must not share across the future cut.
