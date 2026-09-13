# 12. Scale, Delivery, and Cloud

## Description

This topic covers how the system grows, how you ship change, and where the processes run. You learn vertical versus horizontal scale, load balancing, sharding, CDN, CI/CD, rolling, blue/green, and canary releases, feature flags, expand-contract migrations, IaaS, PaaS, FaaS, the twelve-factor app, containers, cost, and vendor lock-in.

Topic 1 defined scale, cost, and maintainability. Topic 2 covered stateful versus stateless services. Topic 11 needs measures for SLIs. Complete Topics 1 to 11 before this topic. Design for the load that you can measure plus a stated growth factor. Do not start with a global shard plan for 30 users.

Use one term for each concept. Vertical scale grows one machine. Horizontal scale grows the number of machines. A shard is a partition of data. A CDN is an edge cache for files. CI is the automatic build and test of every change. CD is the automatic path to an environment. IaaS rents machines. PaaS rents a platform. FaaS rents short-lived functions. Lock-in is the cost to leave a vendor.

---

## Vertical vs horizontal scale, load balancing, sharding, CDN

Vertical scale adds resources to one machine: more CPU, more memory, or a faster disk. The design stays simple. There is a hard ceiling. One machine remains one fault domain (Topic 11).

Horizontal scale adds more machines or more processes that share the work. You need a method to split work: a load balancer, a partition key, or competing consumers. Stateless workers make horizontal scale easier (Topic 2). Sticky sessions make it harder.

Pick vertical first when the load fits one modern host, the team is small, and the state lives in one database that you already operate. Pick horizontal when one host hits a measured ceiling, you need more than one fault domain, or a stateless tier can clone cheaply.

A load balancer spreads incoming connections or requests across more than one healthy instance. Clients see one name. Common methods: round robin, least connections, and hash of a key. Health checks remove a dead instance. A bad check is an outage.

Layer 4 balancing forwards bytes without reading HTTP. Layer 7 balancing can route on path or host. Layer 7 can terminate TLS. Each layer adds policy and cost.

Sharding splits data across stores by a key (user id, tenant, or a hash). Each shard is a smaller database. A query that needs all shards is expensive. A bad key creates a hot shard. Do not shard a campus catalog for fashion. Shard when one primary cannot hold the data or the write rate.

A CDN (content delivery network) is an edge cache for files: images, scripts, and sometimes cacheable HTTP responses. The origin stays the source of truth for the file. A CDN helps "more users" who read the same static bytes. It does not replace a database shard for "more data per user".

Hybrid is common. You scale the web tier horizontally. You scale the database vertically for a long time. Then you add read replicas or shards (Topic 6).

Write the growth axis. "More users" and "more data per user" need different patterns.

### Questions

#### Theoretical questions

1. What is vertical scale?
2. What is horizontal scale?
3. What does a load balancer do?
4. What is a shard?
5. What does a CDN cache?

#### Easy practical tasks

1. Write five sentences that compare vertical and horizontal scale.
2. Make a table: "Tier" and "Usual first move". Add web process and primary database.
3. List four limits of vertical scale.
4. Draw three web processes, one load balancer, and one database.

#### Medium practical tasks

1. Write a one-page scale plan for a campus catalog that grows 10 times in users. State what stays vertical.
2. Explain in eight sentences how a stateful in-memory cart blocks a second web process.
3. Write an ADR: scale the API horizontally, keep one primary database. Include the measure that reopens the database choice.

#### Advanced practical tasks

1. Write a capacity note: current load, 10× load, first bottleneck, and the pattern you apply.
2. Compare a CDN for static files with a search shard plan. Write when each matches a growth axis.

---

## CI/CD; rolling, blue/green, canary

Continuous integration (CI) is the habit that every change builds and is tested in a shared pipeline. Developers integrate often. The pipeline compiles, runs tests, and reports fail fast.

Continuous delivery (CD) is the habit that the same pipeline can deploy a proven build to an environment. Some teams stop at a button. Some teams deploy automatically to production after checks. Continuous deployment is the automatic production step. This handbook uses CD for the path. Write if the last step is automatic.

A pipeline is architecture. It is the only supported way to production. A manual copy from a laptop is a supply-chain risk (Topic 10).

Minimum pipeline:

- Fetch the exact revision
- Build
- Unit and contract tests
- Security scan that the team can run
- Store the artifact
- Deploy to a named environment
- Smoke check
- Record who shipped what

Secrets for the pipeline live in the pipeline store, not in the repository (Topic 10).

A rolling update replaces a subset of instances, then the next subset. Capacity drops during the roll unless you add extra instances first. Old and new code run at the same time. The contract must tolerate that mix (Topic 5).

Blue/green keeps two full environments. Blue serves traffic. You deploy green. You test green. You switch traffic. Rollback is a switch back. Cost is about two environments. A switch does not undo a destructive schema change.

A canary sends a small share of traffic to the new version. You watch SLIs (Topic 11). If the canary burns the budget, you stop.

Pick rolling for simple stateless services with extra capacity. Pick blue/green when you want a fast switch and you can pay for two stacks. Pick canary when the risk is high and you can split traffic.

If the pipeline is slow, people skip it. Keep the default check set fast. Do not disable a failing test to "unblock" without a ticket and a time limit.

### Questions

#### Theoretical questions

1. What is continuous integration?
2. What is continuous delivery in this handbook?
3. What is a rolling update?
4. What is blue/green?
5. What is a canary?

#### Easy practical tasks

1. Write five sentences that define CI, CD, and the three release styles.
2. Make a table: "Step" and "Purpose". Add build, test, store artifact, and smoke check.
3. List four items that must not live in the repository.
4. Write a smoke check list of five HTTP calls for a catalog API.

#### Medium practical tasks

1. Draw a pipeline from commit to production. Mark a manual approval if you use one.
2. Write when you pick rolling versus canary for a campus API.
3. Explain in eight sentences how contract tests protect Topic 5 compatibility during a rolling deploy.

#### Advanced practical tasks

1. Write an ADR: automatic deploy to production versus a button. Include error-budget policy (Topic 11).
2. Design a pipeline for two services and one shared library. Write how you version the library.

---

## Feature flags and expand-contract migrations

A feature flag is a runtime switch that turns a behavior on or off without a new deploy. Flags let you ship dark code, run a canary by user group, or disable a failing feature quickly.

Rules for flags:

- Give each flag an owner and an expiry date.
- A flag that lives forever is a second configuration language. It increases test cost.
- Do not hide a secret behind a flag name in the client.
- Log flag decisions with the correlation identifier (Topic 11).
- Test both states. Untested off is a defect.

Expand-contract (also called parallel change) is a multi-step change that stays compatible. For a column rename:

1. Expand: add the new column. Write both. Read the old column (or both).
2. Migrate: copy data. Change readers to the new column.
3. Contract: stop writes to the old column. Remove the old column.

The same idea applies to APIs (Topic 5) and to message fields (Topic 7). Add, migrate clients, then remove.

A destructive one-step rename during a rolling deploy breaks old instances that still read the old column. Expand-contract avoids that break.

Flags and expand-contract work together. A flag can switch readers after the expand step. Remove the flag in the contract step.

Do not use flags as a substitute for a module boundary. A flag that wraps half the system is over-engineering (Topic 3).

### Questions

#### Theoretical questions

1. What is a feature flag?
2. Why must a flag have an expiry date?
3. What is expand-contract?
4. Why does a one-step rename fail during a rolling deploy?
5. How can a flag and expand-contract work together?

#### Easy practical tasks

1. Write five sentences about flags and expand-contract.
2. Make a table: "Step" and "Readers / writers" for a column rename.
3. List four risks of a forever flag.
4. Write a three-step plan to rename `title` to `book_title`.

#### Medium practical tasks

1. Write an ADR: flags for user-visible features, expand-contract for schema, no forever flags.
2. Design a flag that disables search on the home page when the index is down (Topic 11 degrade).
3. Explain in eight sentences how expand-contract supports blue/green when both stacks share one database.

#### Advanced practical tasks

1. Write a one-page flag standard: naming, owner, expiry, test matrix, and removal.
2. Design an expand-contract for a public JSON field rename. Include two API versions or two fields (Topic 5).

---

## IaaS / PaaS / FaaS, 12-factor, containers

IaaS (infrastructure as a service) gives you virtual machines, disks, and networks. You install the runtime and you operate the OS. You get control. You pay in people time.

PaaS (platform as a service) gives you a place to push an application. The vendor runs the OS and often the balancer. You lose some control. You gain speed.

FaaS (function as a service) runs a function on an event. You do not keep a long-lived process. Cold start, time limits, and local-state limits are part of the model. FaaS fits short jobs and uneven traffic. FaaS is a poor home for a long chatty workflow unless the platform supports that shape.

These are grades of operational work, not grades of quality. A small team can pick PaaS for the API and a managed service for the database.

The twelve-factor app is a set of rules for applications that run on a platform. The document is at [https://12factor.net/](https://12factor.net/). Useful factors for this path: one codebase, explicit dependencies, config in the environment, backing services as attached resources, build/release/run, stateless processes, port binding, process scale, disposable processes, dev/prod parity, logs as event streams, and admin tasks as one-off processes.

These rules support CI/CD and horizontal scale. They do not force microservices (Topic 9). A modular monolith can follow twelve-factor.

A container packages a process with its files and a declared runtime. You build an image. You run the same image in test and in production. Containers are not a cluster. Orchestration (for example Kubernetes) schedules containers, restarts them, and attaches networks. Orchestration is extra operations cost. A student app can run one container on one host.

Do not pick FaaS because the word is fashionable. Measure start time against the latency SLO (Topics 1 and 11). A "serverless" slogan is not a structure. Name the trigger, the time limit, the store, and the failure path.

Law and data location still apply. The grade of service does not remove the constraint (Topic 1).

### Questions

#### Theoretical questions

1. What does IaaS give you?
2. What does PaaS give you?
3. What limits are typical of FaaS?
4. Name four twelve-factor rules that this section lists.
5. How is a container different from an orchestrator?

#### Easy practical tasks

1. Write five sentences that compare IaaS, PaaS, and FaaS.
2. Make a table: "Workload" and "Grade". Add a 24/7 API, a nightly resize, and a database.
3. Open [https://12factor.net/](https://12factor.net/). Write the twelve factor names in order.
4. List four FaaS constraints (time, size, state, and one more).

#### Medium practical tasks

1. Write a one-page pick for a campus catalog: PaaS API plus managed database. Name two rejected grades.
2. Explain in eight sentences how a FaaS cold start can break a 300 ms SLO.
3. Write which twelve-factor rules your student monolith already follows and which it breaks.

#### Advanced practical tasks

1. Write an ADR: PaaS versus virtual machines for the API. Include skill and on-call.
2. Compare FaaS with a worker behind a queue (Topic 7). Write when each wins.

---

## Cost and vendor lock-in

Cost is money and people time (Topic 1). Cloud invoices, licenses, data transfer, and on-call hours are cost. A cheap design that burns the team is still expensive.

Write cost as a constraint. A two-person team cannot operate a large Kubernetes estate and five managed databases. A managed database can be cheaper than people time even if the invoice is higher.

Vendor lock-in is the cost to leave a vendor. Lock-in is not always wrong. A lock that saves a year of platform work can be a good trade. Write the exit cost in the ADR.

Lock-in appears in:

- Proprietary APIs and data formats
- Managed identity and secret stores
- Serverless triggers that have no portable equivalent
- Data egress fees
- Skills that the team has only for one cloud

Reduce lock-in where the cost is low: twelve-factor config, containers, standard SQL, OpenTelemetry (Topic 11), and owned domain types (Topic 4). Accept lock-in where the vendor is the product (a campus identity provider, Topic 10).

Do not copy a vendor reference architecture that assumes a large platform team. Mark every part that your constraints would remove (Topic 3).

Measure. Idle resources, unused logs, and an oversized database class are common leaks. Cost is a quality attribute. Review it like latency.

### Questions

#### Theoretical questions

1. Why is people time part of cost?
2. What is vendor lock-in?
3. Why is lock-in not always wrong?
4. Name four places lock-in appears.
5. How do twelve-factor and owned domain types reduce lock-in?

#### Easy practical tasks

1. Write five sentences about cost and lock-in.
2. Make a table: "Choice" and "Invoice versus people time". Add six rows.
3. List four unused cloud items that can leak money.
4. Write four sentences on why a vendor slide is not a design.

#### Medium practical tasks

1. Write an ADR: managed PostgreSQL versus self-operated IaaS database for a four-person team. Include lock-in.
2. Estimate one year of cost in people hours for Kubernetes plus five services versus PaaS plus a monolith.
3. Explain in eight sentences how data egress can surprise a design that copies a large index every night.

#### Advanced practical tasks

1. Write a one-page cost review: invoice lines, people time, unused resources, and one cut.
2. Write an exit note for one vendor: what you would move, what format you own, and the time estimate.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do stateless instances, a load balancer, and CI/CD work together to ship horizontal scale?
2. Why must expand-contract exist before you trust rolling or blue/green with a shared database?
3. When is a CDN a better first scale step than a shard?
4. How do twelve-factor rules support both PaaS and containers without a service split?
5. How do cost and lock-in change the "when not to split" rule from Topic 9?

#### Easy practical tasks

1. Write a one-page cheat sheet: scale types, balancer, shard, CDN, CI/CD, three releases, flags, expand-contract, IaaS/PaaS/FaaS, twelve-factor, containers, cost, lock-in.
2. For a to-do app, write vertical database, one API instance, a simple CI pipeline, and no shard.
3. Bookmark [https://12factor.net/](https://12factor.net/). Write which three factors you will apply first.
4. Draw blue and green boxes and one database. Mark the expand-contract risk.

#### Medium practical tasks

1. Write a short delivery brief for a campus shop: pipeline, rolling deploy, one flag, and a column rename plan.
2. Take a teammate design that starts with Kubernetes and ten shards. Rewrite it for 2 000 users.
3. Write a twelve-week plan: weeks for CI and one instance, then the measure that adds a second instance.

#### Advanced practical tasks

1. Write a capacity plus cost model: 1×, 10× users, first bottleneck, invoice, and people time.
2. Compare IaaS plus containers with PaaS for the same monolith. Write lock-in and on-call for a five-person team.
