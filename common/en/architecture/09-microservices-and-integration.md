# 9. Microservices and Integration

## Description

A microservice is a deployable unit that owns its data and exposes a contract. Integration is the way two systems exchange data or start work. This topic covers the service as a deployable with its own data, network and operations cost, when not to split, sagas, choreography, orchestration, the API gateway, file transfer, shared database, RPC, messaging, and the anti-corruption layer.

Microservices are not a default. A modular monolith is the usual start (Topic 3). You add a second deployable unit only after a real pain. Complete Topics 1 to 8 before this topic. Pair data ownership with Topic 6. Pair remote failure with Topic 8. Pair messages with Topic 7.

Use one term for each concept. A service is a deployable unit with a contract and a data owner. A module is a design boundary inside one deployable unit. Size in lines of code is not the definition. A saga is not a distributed database transaction.

---

## Service as a deployable with its own data

A service in this handbook is a unit that you build, start, stop, and ship on its own clock. The unit has a process (or a small set of processes that you always ship together). The unit owns a data store. Other services do not write the tables of that store.

Ownership is the rule. The owner accepts writes that change the source of truth. Other services read through an API or through events (Topics 5, 6, and 7). A shared table between two deployable units is a hidden contract. Topic 6 names that shape an integration-database anti-pattern.

A service also owns its failures. If the process is down, the owner of that process is the team that ships it. If the schema change is wrong, the same team rolls it back. A service without an owner is not a complete service.

The contract is part of the service. HTTP JSON, gRPC, or events are the public face (Topic 5). You can change internals. You cannot change the contract without a compatibility window.

A library that many processes import is not a service. A package in a monolith is not a service. Those units do not deploy apart and do not own a store.

The word "micro" suggests a small size. Size in lines of code is not the design rule. The design rule is the boundary. A useful boundary matches a bounded context or a clear ownership line (Topic 4).

A bad split follows technical layers: one "API service", one "logic service", and one "database service". Those units cannot change alone. Every feature crosses the network. That shape copies a layered monolith and adds latency (Topic 3).

Do not split a store on the first day "for microservices". One database with module-owned tables is the default in a modular monolith (Topics 3 and 6).

### Questions

#### Theoretical questions

1. What three parts define a service in this handbook?
2. Who may write the source of truth of a service?
3. Why is a shared table between two deployable units a hidden contract?
4. How is a service different from a library?
5. Why is size in lines of code not the design rule?

#### Easy practical tasks

1. Write five sentences that define a service. Use only facts from this section.
2. Make a table: "Unit" and "Service? (yes/no)". Add a package, a library, a worker process with its own store, and a shared reporting database.
3. List four ways a second process can read data that it does not own.
4. Draw one shop with two services. Label the store of each service and the contract between them.

#### Medium practical tasks

1. Take a modular monolith with four modules and one database. Write which module could become a service and which data it must take.
2. Write an ADR title and a one-paragraph context for "orders own the order tables; billing does not write them".
3. Describe a defect that appears when two teams write the same `customers` table from two deployable units.

#### Advanced practical tasks

1. Write a one-page service charter: name, owner, store, contract, and on-call. Fill it for two services in a library loan system.
2. Compare "one database, module-owned tables" with "one database per service" for a team of four. Write operations cost and schema-change cost.

---

## Network and ops cost; when not to split

Each extra deployable unit adds a network hop and an operations bill. The network is not reliable, latency is not zero, and the path is not free (Topic 8). Operations cost includes pipelines, on-call, dashboards, secrets, and schema migrations.

When not to split:

- One team owns the full product and can ship on one clock.
- You cannot draw a stable bounded context (Topic 4).
- You need one transaction for the user action and you have no saga story.
- You have no measure that a module is the bottleneck (Topic 1).
- You cannot staff on-call for another process.
- The split would follow layers (API / logic / database) instead of ownership.

When a split can be valid:

- Two teams must release on independent clocks and a modular monolith already fails that need.
- A measured load needs a different scale path for one part (Topic 12).
- Law or a vendor forces a separate system.
- A hard isolation rule (a fault in one part must not crash the other process).

Write the reason in an ADR with a measure or a constraint. "We might need to scale" is not enough.

A monolith can still run many instances of the same unit. That is scale, not a service split (Topic 12).

If you already split too early, a merge back to a modular monolith can be the correct architecture. Record that in an ADR. Fashion is not a quality attribute.

### Questions

#### Theoretical questions

1. What costs does a second deployable unit add?
2. Name four cases when you must not split.
3. Name three cases when a split can be valid.
4. Why is "we might need to scale" not enough?
5. How is many instances of one unit different from many services?

#### Easy practical tasks

1. Write five sentences about network and operations cost.
2. Make a table: "Situation" and "Split? (yes/no)". Add six rows.
3. List five operations tasks that exist once in a monolith and once per service after a split.
4. Write four sentences to a teammate who wants five services for a student shop.

#### Medium practical tasks

1. Given a five-person team and a 2000-user app, write a one-page recommendation: monolith plus which scale steps.
2. Write an ADR that rejects a three-service split for a classroom project.
3. Compare people hours for one year: one monolith with three instances versus five services.

#### Advanced practical tasks

1. Write a decision tree from constraints to "keep the monolith" or "split this module".
2. Find a public story where a team moved back to a monolith. Write five facts from that story. Do not copy long text.

---

## Saga, choreography, orchestration, API gateway

A saga is a sequence of local transactions that together implement a business action across more than one owner. Each step has a compensating action if a later step fails. A saga is not a distributed ACID transaction. You accept eventual consistency (Topic 4).

Example: reserve stock, then capture payment, then create shipment. If payment fails, a compensate releases stock.

Choreography is a style where each service reacts to events. There is no central director. `OrderPlaced` leads to `PaymentCaptured` leads to `ShipmentCreated`. The flow lives in the subscribers.

Orchestration is a style where a coordinator sends commands and waits for results. The flow lives in the orchestrator.

Pick choreography when the flow is simple and owners must stay decoupled. Pick orchestration when the flow is long, when compensations are many, or when you need one place to see the state of the action.

An API gateway is an edge process that accepts client calls and forwards them to services. It can apply shared policy: TLS termination, authentication check, rate limit, and routing. It must not become the home of domain rules (Topic 3). Price and loan policy stay in the owner.

A gateway is a fault domain. If the gateway is down, public clients fail even if services are up. Give it redundancy (Topic 11).

Do not use a gateway as a hidden integration database. Do not put all joins in the gateway if that makes it a second monolith. A BFF is a relative of a gateway for one client type (Topic 13).

Write who owns the saga state. A missing owner produces lost orders and double charges.

### Questions

#### Theoretical questions

1. What is a saga?
2. How does a saga differ from a distributed ACID transaction?
3. What is choreography?
4. What is orchestration?
5. What may an API gateway do, and what must it not own?

#### Easy practical tasks

1. Write five sentences that define saga, choreography, orchestration, and gateway.
2. Draw reserve-stock, capture-payment, compensate-stock.
3. Make a table: "Duty" and "Gateway or owner". Add TLS, price rule, and rate limit.
4. List four risks of a saga without compensations.

#### Medium practical tasks

1. Design a place-order saga with two compensations. Write the user-visible states.
2. Write an ADR: choreography for mail and search; orchestration for checkout.
3. Explain in eight sentences how a gateway that contains all business rules becomes a monolith at the edge.

#### Advanced practical tasks

1. Write a one-page saga catalog: steps, compensations, owner of state, and timeout.
2. Compare gateway routing with a BFF (preview of Topic 13). Write when you need both.

---

## File, shared DB, RPC, messaging

Hohpe and Woolf describe four classic integration styles. This handbook uses the same four names. See [https://www.enterpriseintegrationpatterns.com/](https://www.enterpriseintegrationpatterns.com/).

File transfer is integration through a file that one system writes and another system reads. The file can sit on a disk, an object store, or an SFTP host. The producer and the consumer do not need to run at the same time. The style fits batch work and partner exchange. It is a poor fit for a user click that needs an answer in 200 ms. You must define format, version, completeness signal, retention, and idempotent load.

A shared database is a store that more than one system uses as a common schema. Both systems issue SQL against the same tables. This handbook treats that style as high coupling (Topic 6). Use it only if you inherit it. Stop new writers. Plan an exit.

RPC is a request and a reply (Topic 5). The caller waits. Use RPC when the user action needs an answer now and both sides can share a contract. Apply timeouts, retries, and breakers (Topic 8).

Messaging is an asynchronous message (Topic 7). The caller does not wait for all subscribers. Use messaging when you need a time split, a buffer, or many independent consumers.

Pick the style from the constraint:

- Same time, need a reply: RPC.
- Different time, batch, partner file: file transfer.
- Different time, many subscribers, jobs: messaging.
- Shared tables: reject for new work.

You can mix styles. A public HTTP API can accept an order. An event can notify search. A nightly file can feed a warehouse.

Security is part of every style. Authenticate the sender. Bound size. Do not put secrets in file names or URLs (Topic 10).

### Questions

#### Theoretical questions

1. What is file transfer as an integration style?
2. Why is a shared database a high-coupling style?
3. When do you pick RPC?
4. When do you pick messaging?
5. Why can you mix styles in one product?

#### Easy practical tasks

1. Write five sentences that name the four styles. Use only facts from this section.
2. Make a table: "Need" and "Style". Add six needs.
3. List four contract items for a nightly CSV drop.
4. Draw HTTP place-order, event to search, and a nightly export file.

#### Medium practical tasks

1. Design a daily catalog file from a vendor to a shop. Include columns, version, and late-file behavior.
2. Write eight sentences that compare a file drop with an HTTP API for the same catalog.
3. Plan an exit from a shared campus database: owner, API, sunset of SQL access.

#### Advanced practical tasks

1. Write a one-page integration map for a campus system: four partners, one style each, and the contract type.
2. Read the EIP site index. Map five pattern names to the four styles in your own words.

---

## Anti-corruption layer

An anti-corruption layer (ACL) is a translation boundary. It protects your model from an external model. Incoming data and outgoing data pass through types that you own. The external names, identifiers, and rules do not leak into your domain (Topic 4).

Use an ACL when you integrate with a vendor, a legacy system, or another bounded context that you do not control.

The ACL can be:

- A package in a modular monolith
- A small adapter service
- A set of mapper functions at the edge of a module

The ACL must not become a second home for your business rules. It translates. The domain still decides what a loan may do.

Without an ACL, a vendor field name appears in your tables, your API, and your UI. A vendor version change then rewrites your system. That is hidden coupling (Topic 2).

Translate identifiers. Their `customerCode` is not your `PatronId` unless you write the map. Store the foreign identifier as a value on your entity if you must round-trip. Do not adopt their whole type.

An ACL is not a shared database view that both sides write. That view is still a shared model.

Test the ACL with fixtures that look like the external payload. Do not hit the vendor from every unit test.

### Questions

#### Theoretical questions

1. What is an anti-corruption layer?
2. When do you use an ACL?
3. What work must stay in the domain, not in the ACL?
4. Why do vendor names in your tables create hidden coupling?
5. How do you handle a foreign identifier?

#### Easy practical tasks

1. Write five sentences that define an ACL.
2. Make a table: "External field" and "Your type". Add six rows for a payment vendor.
3. Draw vendor, ACL, and your domain. Label arrows.
4. List four defects that appear when a vendor JSON struct is your entity.

#### Medium practical tasks

1. Design an ACL for a campus identity system that your loan context consumes. Write the map of identifiers.
2. Write an ADR: ACL package in the monolith, not a new service, for one vendor.
3. Explain in eight sentences how an ACL differs from an integration database.

#### Advanced practical tasks

1. Write a one-page ACL standard: where it lives, what it may import, and how you version the map.
2. Review a public webhook payload. Sketch an ACL that would protect a notes app. Stay defensive. Do not copy secrets.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do service, boundary, and data ownership form one definition?
2. Why do network cost and saga complexity argue for a modular monolith first?
3. When do you pick choreography, and when do you pick orchestration, for the same shop?
4. How do the four integration styles relate to Topics 5, 6, and 7?
5. How does an ACL protect ubiquitous language across a vendor boundary?

#### Easy practical tasks

1. Write a one-page cheat sheet: service definition, when not to split, saga pair, four styles, ACL.
2. For a to-do app, write "one service", four modules, HTTP only, no saga, no ACL.
3. Draw a gateway, two services, and two stores. Mark who may write each store.
4. Bookmark the EIP site and `architecture.topics.md`. Write one sentence on when you open each.

#### Medium practical tasks

1. Write a short integration brief for a campus shop: public HTTP, one vendor file, one event, and a rejected shared DB.
2. Take a teammate design that splits API / logic / database into three services. Rewrite it as a modular monolith or as two ownership services.
3. Write a twelve-week plan: weeks in a monolith, and the pain that would reopen a split.

#### Advanced practical tasks

1. Write a service extract checklist: language, owned tables, contract, saga, ACL, on-call. Fill it for one module.
2. Compare a merge back to a monolith with a further split. Write the quality attributes that each move improves or weakens.
