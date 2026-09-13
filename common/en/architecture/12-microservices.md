# 12. Microservices

## Description

A microservice is a deployable unit that owns its data and exposes a contract. This topic covers the service as a deployable unit, the boundary as the real design choice, network and operations cost, choreography versus orchestration, sagas, the API gateway, and cases where you must not use microservices.

Microservices are not a default. A modular monolith is the usual start (Topic 5). You add a second deployable unit only after a real pain: independent scale, an independent release clock, or a hard isolation rule. Complete Topics 1 to 11 before this topic. Pair data ownership with Topic 9. Pair remote failure with Topic 11.

Use one term for each concept. A service is a deployable unit with a contract and a data owner. A module is a design boundary inside one deployable unit. Size in lines of code is not the definition. A saga is not a distributed database transaction.

---

## Service as a deployable with its own data

A service in this handbook is a unit that you build, start, stop, and ship on its own clock. The unit has a process (or a small set of processes that you always ship together). The unit owns a data store. Other services do not write the tables of that store.

Ownership is the rule. The owner accepts writes that change the source of truth. Other services read through an API or through events (Topics 8, 9, and 10). A shared table between two deployable units is a hidden contract. Topic 9 names that shape an integration-database anti-pattern.

A service also owns its failures. If the process is down, the owner of that process is the team that ships it. If the schema change is wrong, the same team rolls it back. A service without an owner is not a complete service.

The contract is part of the service. HTTP JSON, gRPC, or events are the public face (Topic 8). You can change internals. You cannot change the contract without a compatibility window (Topic 18).

A library that many processes import is not a service. A package in a monolith is not a service. Those units do not deploy apart and do not own a store.

Do not split a store on the first day "for microservices". One database with module-owned tables is the default in a modular monolith (Topics 5 and 9).

### Questions

#### Theoretical questions

1. What three parts define a service in this handbook?
2. Who may write the source of truth of a service?
3. Why is a shared table between two deployable units a hidden contract?
4. How is a service different from a library?
5. Why must a service have an owner?

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

## Size is not the point; boundary is

The word "micro" suggests a small size. Size in lines of code is not the design rule. The design rule is the boundary.

A useful boundary matches a bounded context or a clear ownership line (Topic 7). The people who change the rules of billing can ship billing. The people who change the rules of search can ship search. The contract between those units is smaller than the internals.

A bad split follows technical layers: one "API service", one "logic service", and one "database service". Those units cannot change alone. Every feature crosses the network. That shape copies a layered monolith and adds latency (Topic 6).

A bad split follows size targets: "each service must stay under 500 lines". That rule cuts a cohesive job into many remote calls. Cohesion falls. Coupling rises (Topic 3).

Tests for a good boundary:

- One team can explain the language of the unit (ubiquitous language, Topic 7).
- Most changes stay inside one unit.
- The unit can accept a write without a remote write in the same user request.
- You can state the source of truth in one sentence.

If you cannot draw the boundary on a domain diagram, do not draw it on a deploy diagram.

A large service with a clear boundary is better than ten tiny services with a shared table. Prefer fewer, clearer units.

### Questions

#### Theoretical questions

1. Why is size in lines of code not the definition of a service?
2. What makes a boundary useful?
3. Why is a split by technical layer a poor service split?
4. Name four tests for a good boundary.
5. When is one large service better than many small services?

#### Easy practical tasks

1. Write four sentences that replace "micro" with "boundary" in your own words.
2. Make a table: "Split idea" and "Good boundary? (yes/no)". Add "by layer", "by bounded context", and "by file count".
3. List five change requests for a shop. Mark which service each change must touch if boundaries are good.
4. Open Topic 7 notes. Write the name of one bounded context that could be one service later.

#### Medium practical tasks

1. Review a blog diagram that shows 15 services for a student app. Mark three splits that look like layers, not contexts.
2. Write eight sentences that compare a modular monolith boundary with a service boundary.
3. Given a "user service" that every other service calls for every request, write why that boundary is weak.

#### Advanced practical tasks

1. Draw two C4 container views of the same product: five services with weak boundaries and three services with strong boundaries. Write which view you would ship first.
2. Write a one-page review rubric that a team uses to accept or reject a new service proposal. Include language, data owner, and change locality.

---

## Network and ops cost

Each new service adds a network hop. A network hop adds latency, partial failure, timeouts, retries, and idempotency work (Topics 4 and 11). A local function call in a monolith does not need that machinery.

Each new service also adds operations cost. You must build, store, deploy, configure, watch, and page a new unit. You need logs, metrics, traces, and a runbook (Topics 2, 15, and 16). A team of two people cannot operate twenty services with the same quality as a large platform team (Topic 1).

Money cost rises. You pay for more compute, more load balancers, more stores, and more traffic between zones (Topics 2 and 19). People cost rises. On-call load and training load rise.

Debug cost rises. A defect can sit in service A, in the contract, in the broker, or in service B. You need a correlation identifier to join the path (Topic 16).

Test cost rises. Contract tests and environment setup replace a single process boot.

You pay this cost when a measure requires it. Independent scale of a hot read path can pay for one extract. A fashion word does not pay for ten extracts.

Write the cost in the ADR. Name the team hours per week and the new failure modes. A design that hides operations cost is not complete.

### Questions

#### Theoretical questions

1. What failure machinery does a network hop add?
2. What operations work does each new service add?
3. How does a small team size limit the number of services?
4. Why does debug cost rise when you add services?
5. When is the extra cost acceptable?

#### Easy practical tasks

1. Make a two-column table: "Cost type" and "Example". Add latency, money, people, debug, and test.
2. List six new failure modes that appear when one in-process call becomes HTTP.
3. Write five sentences on why twenty services can fail a two-person team.
4. Name four signals that you must have before the first extract (logs, deploy, timeout, owner).

#### Medium practical tasks

1. Estimate weekly operations hours for 3 services versus 12 services for a team of four. State your assumptions.
2. Write a one-page cost note for an extract of a report worker. Include money and people.
3. Take a convert-to-microservices slide. Add the missing operations list. Write which items reject the slide.

#### Advanced practical tasks

1. Build a cost model (table) with columns: service count, deploy pipelines, stores, on-call rotations, and estimated hours. Fill three rows (1, 5, 15).
2. Write an ADR that rejects a split because operations cost exceeds the benefit. Include the measure that would reopen the decision.

---

## Choreography vs orchestration

Choreography and orchestration are two ways to run a multi-step business flow across services.

In choreography, each service reacts to events. Service A publishes `OrderPlaced`. Service B reserves stock and publishes `StockReserved`. Service C captures payment. No central process owns the full sequence. Topic 10 covers events.

In orchestration, one process sends commands and waits for replies. An orchestrator tells billing to capture payment, then tells warehouse to ship. The orchestrator holds the current step.

Choreography benefits:

- Services stay decoupled from a central flow owner.
- New subscribers can react without a change of the producer (if the event is stable).

Choreography costs:

- The full path is hard to see.
- Failure and retry logic spread across many handlers.
- Cycles of events can appear if names are poor.

Orchestration benefits:

- One place shows the sequence and the timeout policy.
- Compensation is easier to list in order.

Orchestration costs:

- The orchestrator can become a hotspot and a coupling hub.
- A chatty orchestrator turns every step into a remote command.

Do not mix the two without a rule. A common rule: choreography for independent reactions (mail, search index). Orchestration for a money path that needs a visible sequence and a timeout.

A workflow engine is one implementation of orchestration. You do not need a workflow product on day one. A state machine in one service can be enough.

### Questions

#### Theoretical questions

1. What is choreography?
2. What is orchestration?
3. What benefit does choreography give?
4. What benefit does orchestration give?
5. When do you pick orchestration for a money path?

#### Easy practical tasks

1. Write four sentences that compare the two styles. Use one shop example.
2. Make a table: "Flow" and "Style". Add "send receipt", "capture payment then ship", and "update search index".
3. List three risks of event cycles in choreography.
4. Draw both styles for "place order → reserve stock → pay".

#### Medium practical tasks

1. Design a three-step loan flow both ways. Write where you see the current step in each design.
2. Write eight sentences on how a new "fraud check" step lands in each style.
3. Identify a coupling hub: an orchestrator that every team must change. Write two ways to shrink it.

#### Advanced practical tasks

1. Write a one-page decision table: flow type, style, timeout owner, and visibility tool (log, trace, or dashboard).
2. Compare a small state machine in the order service with a central workflow product. Write operations cost for a five-person team.

---

## Saga

A saga is a sequence of local transactions. Each local transaction belongs to one service and one store. If a later step fails, the saga runs compensating actions that undo or repair earlier steps.

A saga is not one ACID transaction across services. You do not hold one database lock across a network call (Topic 7). Each step commits. The world can show a temporary disagreement. That is eventual consistency.

Two common saga shapes:

- Orchestrated saga: one coordinator sends commands and records saga state.
- Choreographed saga: each service publishes an event after its local commit. The next service starts. Compensation is also event-driven.

Compensation is a business action, not a magic rollback. `ReleaseStock` after `ReserveStock` is compensation. You cannot always return to the old state. A sent mail cannot unsend. The compensation can be a follow-up mail or a manual case.

Design rules:

- Each step is idempotent (Topics 4 and 11).
- Each step has a timeout.
- You record saga state or you can rebuild it from events.
- You define the user message when the saga is in the middle (reserved but not paid).
- You have a repair list for steps that fail after retries.

Do not start with a saga product. Start with two steps, a state table, and a clear compensation. Add a tool when the number of steps and the audit need grow.

Write the happy path and the fail path in the same ADR. A saga without compensation is an incomplete design.

### Questions

#### Theoretical questions

1. What is a saga?
2. Why is a saga not one ACID transaction across services?
3. What is a compensating action?
4. Why must saga steps be idempotent?
5. What must you define for the user when the saga is in the middle?

#### Easy practical tasks

1. Write five sentences that define a saga. Use only facts from this section.
2. Make a table: "Step" and "Compensation". Add reserve stock, capture payment, and send mail.
3. List four fields of a saga state record (identifier, step, status, time).
4. Draw an orchestrated saga with three steps and one fail path.

#### Medium practical tasks

1. Design a saga for "loan a book": reserve copy, create loan, send mail. Write timeouts and compensations.
2. Explain in eight sentences why a network call inside one SQL transaction is worse than a saga.
3. Write the user-visible states: started, reserved, paid, failed. Map each state to a screen sentence.

#### Advanced practical tasks

1. Write a one-page saga playbook: idempotency keys, outbox (Topic 10), repair list, and metrics.
2. Compare choreographed and orchestrated sagas for a five-step fulfillment flow. Write which style you pick and why.

---

## API gateway

An API gateway is a process at the edge that accepts client calls and forwards them to services. The gateway is a reverse proxy with extra policy. It is not the place for domain rules.

Typical gateway duties:

- TLS termination
- Authentication check (token present and valid)
- Rate limits and size limits
- Routing to a service
- Request identifier injection
- Protocol translation (for example browser HTTP JSON to an internal RPC)

The gateway must not own the source of truth. The gateway must not contain "how to price an order". Those rules belong in the owner service (Topics 6 and 7).

A gateway is a new fault domain. If the gateway is down, clients cannot reach healthy services. You need redundancy and a simple configuration (Topic 15). A huge plugin list makes the gateway a second monolith.

BFF (backend for frontend) is a related idea. A BFF is a gateway-like process that serves one client type. Topic 20 covers BFF.

Do not add a gateway because a slide shows one. One service with a load balancer can be enough. Add a gateway when many services need the same edge policy and you want one public host name.

Clients of public APIs must not need a map of every internal service. The gateway (or an equivalent edge) hides internal topology (Topic 8).

### Questions

#### Theoretical questions

1. What is an API gateway?
2. Which duties belong at the gateway?
3. Which rules must not live in the gateway?
4. Why is the gateway a fault domain?
5. When do you add a gateway?

#### Easy practical tasks

1. Make a table: "Duty" and "Gateway or service". Add TLS, price calculation, rate limit, and stock reserve.
2. Write four sentences that define a gateway. Use only facts from this section.
3. List five risks of a gateway that contains domain rules.
4. Draw clients, one gateway, and three services. Label what the client knows.

#### Medium practical tasks

1. Write a one-page gateway policy: auth header, body size, timeout, and correlation header.
2. Compare "one public service" with "gateway plus four services" for a campus app. Write when the first shape wins.
3. Describe a failure: gateway up, billing down. Write the client status and the log fields.

#### Advanced practical tasks

1. Write an ADR: add or reject a gateway for a three-service system. Include operations cost.
2. Design a rollback: clients call one service directly if the gateway is down. Write the security and routing problems of that fallback.

---

## When not to use microservices

Do not use microservices when the team is small, the product is young, and one database still matches ownership. Topic 5 states the default: a modular monolith.

Do not use microservices to copy a large-company slide. Those slides hide platform teams, SREs, and years of extracts.

Do not use microservices when you cannot name the boundary in domain words. A split that only copies folders to repositories increases cost and does not increase independence.

Do not use microservices when you cannot operate the unit. No on-call, no metrics, no timeout policy, and no backup is a reason to stay in one process.

Do not use microservices to make a resume list. Architecture is structure and costly decisions (Topic 1). A service count is not a quality.

Signals that you must wait:

- You still change the ubiquitous language every week.
- Two modules always change in the same commit.
- You have no request identifier and no SLO.
- The first extract has no rollback.

Signals that one extract can be right:

- One module has a different scale or a different release clock.
- Table ownership is already real.
- The contract is stable.
- The team can operate one more unit.

If the pain is unclear code, fix the modular monolith. If the pain is a real operations or ownership constraint, extract one service. Stop. Measure.

### Questions

#### Theoretical questions

1. What is the default structure for a young product?
2. Why is a large-company slide a poor reason to split?
3. Why must you name the boundary in domain words before you split?
4. Name four signals that you must wait.
5. Name four signals that one extract can be right?

#### Easy practical tasks

1. Write a one-page "do not split" checklist with ten items from this section.
2. Make a table: "Reason to split" and "Valid? (yes/no)". Add resume, scale, fashion word, and isolation law.
3. List five tasks that improve a monolith before any extract.
4. Write four sentences that a reviewer says when a week-two project wants eight services.

#### Medium practical tasks

1. Role-play on paper: a teammate wants a microservice for each database table. Write a one-page reject note that cites Topics 5, 7, and 9.
2. Given a healthy modular monolith with one hot report, write the extract that you would consider and the extracts that you would refuse.
3. Write two ADR summaries: "stay monolith" and "extract search". Include the measure that would change the choice.

#### Advanced practical tasks

1. Write a two-page extraction gate: required evidence before the first network hop. Include SLO, ownership, and rollback.
2. Compare three public post-mortems or blogs that praise a split. Mark which praises name a measure and which praises name a fashion word.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do data ownership, boundary language, and operations cost form one test for a new service?
2. Why does this path place microservices after monoliths, data ownership, messaging, and distributed basics?
3. How do choreography, orchestration, and saga relate without becoming the same idea?
4. What stays in the gateway, and what stays in the owner service, and why does that split protect evolvability?
5. How can an extract that you stop after better modules still be a success?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic.
2. For a to-do product with two developers, write why you stay on one deployable unit. List three later signals that could change that choice.
3. Draw a C4 container diagram with one API, one worker, and one database. Mark what would have to be true before a second database appears.
4. Write three ADR titles: no table split yet, saga for refunds, reject extra gateway.

#### Medium practical tasks

1. Write a short architecture brief (one or two pages) for a campus shop that might extract billing in six months. Include the gate, the saga sketch, and the ops cost.
2. Take a teammate design that is only a list of service names. Rewrite it as boundaries, stores, contracts, and rejected splits.
3. Map Topics 5, 9, 10, and 11 onto one extract of a mail worker. Write what each earlier topic forces you to finish first.

#### Advanced practical tasks

1. Design a twelve-week plan that starts as a modular monolith and allows at most one extract. Write the weekly evidence that keeps or kills the extract.
2. Read a public "we moved to microservices" report. Score it from 0 to 5 on boundary, data ownership, and operations cost. Explain each score. Do not copy long passages.
