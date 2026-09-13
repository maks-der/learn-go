# 19. Cloud and Infrastructure

## Description

Cloud and infrastructure are the platforms that run the processes, stores, and networks that your structure needs. This topic covers IaaS, PaaS, and FaaS; the twelve-factor app; containers and a Kubernetes survey; multi-AZ and multi-region; managed databases and brokers; cost as a constraint; and vendor lock-in trade-offs.

Infrastructure is an architecture decision (Topic 1). Topic 15 needs zones for failover. Topic 17 needs units that you can clone. Complete Topics 1 to 18 before this topic. Pair application networking with `net.topics.md`. Do not treat a vendor slide as a design.

Use one term for each concept. IaaS rents machines. PaaS rents a platform for your code. FaaS rents short-lived functions. A container is a packaged process. Orchestration schedules containers. An availability zone is a failure domain inside a region. A region is a geographic group. Lock-in is the cost to leave a vendor.

---

## IaaS / PaaS / FaaS

IaaS (infrastructure as a service) gives you virtual machines, disks, and networks. You install the runtime and you operate the OS. You get control. You pay in people time.

PaaS (platform as a service) gives you a place to push an application. The vendor runs the OS and often the balancer. You lose some control. You gain speed.

FaaS (function as a service) runs a function on an event. You do not keep a long-lived process. Cold start, time limits, and local-state limits are part of the model. FaaS fits short jobs and uneven traffic. FaaS is a poor home for a long chatty workflow unless the platform supports that shape.

These are grades of operational work, not grades of quality. A small team can pick PaaS for the API and IaaS or a managed service for the database.

Mix is normal. A function that resizes images can sit next to a long-running API. Write the constraints: timeout, payload size, and identity (Topic 14).

Do not pick FaaS because the word is fashionable. Measure start time against the latency SLO (Topics 2 and 15).

A "serverless" slogan is not a structure. Name the trigger, the time limit, the store, and the failure path.

Law and data location still apply. The grade of service does not remove the constraint (Topic 1).

### Questions

#### Theoretical questions

1. What does IaaS give you?
2. What does PaaS give you?
3. What limits are typical of FaaS?
4. Why is a mix of grades normal?
5. Why is "serverless" not a structure by itself?

#### Easy practical tasks

1. Write five sentences that compare IaaS, PaaS, and FaaS.
2. Make a table: "Workload" and "Grade". Add a 24/7 API, a nightly resize, and a database.
3. List four FaaS constraints (time, size, state, and one more).
4. Write four sentences on people-time cost of IaaS for a two-person team.

#### Medium practical tasks

1. Write a one-page pick for a campus catalog: PaaS API plus managed database. Name two rejected grades.
2. Explain in eight sentences how a FaaS cold start can break a 300 ms SLO.
3. Draw triggers and stores for one function and one API. Label identity.

#### Advanced practical tasks

1. Write an ADR: PaaS versus virtual machines for the API. Include skill and on-call.
2. Compare FaaS with a worker behind a queue (Topic 17). Write when each wins.

---

## 12-factor app

The twelve-factor app is a set of rules for applications that run on a platform. The document is at [https://12factor.net/](https://12factor.net/). This section is a survey. It is not a replacement for the site.

Useful factors for this path:

- One codebase, many deploys
- Explicit dependencies
- Config in the environment, not in the code
- Backing services as attached resources
- Build, release, run as separate stages
- Stateless processes (Topic 4)
- Port binding
- Concurrency by process scale (Topic 17)
- Disposable processes (fast start, graceful stop)
- Dev/prod parity
- Logs as event streams (Topic 16)
- Admin tasks as one-off processes

These rules support CI/CD and horizontal scale (Topics 17 and 18). They do not force microservices (Topic 12). A modular monolith can follow twelve-factor.

Config in the environment must not become secrets in a ticket (Topic 14). Use a secret store and inject at run.

Dev/prod parity does not mean identical hardware. It means the same backing-service types and a short gap between versions.

Do not treat the list as law when a factor fights a constraint. A desktop UI tool is not a twelve-factor target. A batch that must pin to one host for a hardware key is a stated exception. Write the exception in an ADR.

Read the site. Write the factors in your own words. Do not copy long passages.

### Questions

#### Theoretical questions

1. What is the twelve-factor app document?
2. Why does config belong outside the code?
3. How do disposable processes help rolling deploys?
4. Why can a monolith follow twelve-factor?
5. When do you record an exception?

#### Easy practical tasks

1. Open [https://12factor.net/](https://12factor.net/). Write the twelve factor names in order.
2. Write four sentences on logs as streams.
3. Make a table: "Factor" and "Your project (yes/no)". Fill six rows.
4. List three secrets that must not sit in the codebase.

#### Medium practical tasks

1. Score a homework app against twelve factors. Write one repair per miss.
2. Explain in eight sentences how build/release/run maps to your pipeline (Topic 18).
3. Write a one-page twelve-factor subset that a four-person team will enforce.

#### Advanced practical tasks

1. Write an ADR that accepts ten factors and rejects two for a stated constraint.
2. Compare twelve-factor with a stateful service that stores sessions in memory. Write the scale cost (Topic 17).

---

## Containers and orchestration (Kubernetes survey)

A container is a packaged runtime for a process: files, libraries, and a start command. The host shares a kernel. The container is not a full virtual machine. Images are supply-chain artifacts (Topic 14).

Orchestration schedules containers across machines. It restarts failed processes. It can run a desired replica count. Kubernetes is a common orchestrator. This section is a survey. It is not a cluster tutorial.

Kubernetes ideas that architects must know:

- Pod: one or more containers that share a network identity
- Deployment: a desired replica set and a roll strategy
- Service: a stable name in the cluster
- Job: a finite batch
- ConfigMap and Secret: configuration injection
- Namespace: a scope for names and some policy

You also inherit new failure modes: image pull, scheduler pressure, bad resource limits, and a control plane outage.

A small team can use a managed container service or a PaaS and never run a raw cluster. Topic 1 team constraints apply. "We use Kubernetes" is not a design (Topic 1).

Resource requests and limits are architecture. If you omit them, one process can starve others (bulkhead idea, Topic 11).

Do not run a database in a cluster if the team cannot operate storage, backup, and failover. Prefer a managed database (next sections) unless you have a reason.

Learn enough to read a diagram. Defer operator skill until you have a platform owner (Topic 21).

### Questions

#### Theoretical questions

1. What is a container?
2. What does orchestration add?
3. What is a Pod in this survey?
4. Why is "we use Kubernetes" not a design?
5. Why can a database in a cluster be the wrong default?

#### Easy practical tasks

1. Write five sentences that define containers and orchestration.
2. Make a table: "Object" and "Job". Add Pod, Deployment, Service, and Job.
3. List four new failure modes that a cluster adds.
4. Write four sentences on image as a supply-chain artifact.

#### Medium practical tasks

1. Draw a Deployment of three API Pods and one Service. Label the balancer outside if you have one.
2. Explain in eight sentences how a missing memory limit becomes a noisy-neighbor incident.
3. Write an ADR: managed container service versus a self-operated cluster. Include people cost.

#### Advanced practical tasks

1. Write a one-page Kubernetes reading list for an architect (official concepts only). Do not copy long passages.
2. Compare a PaaS push with a cluster Deployment for a two-person team. Write a reject rule for raw Kubernetes.

---

## Multi-AZ / multi-region

An availability zone (AZ) is a failure domain inside a cloud region. Zones have independent power and network in the vendor model. A region is a geographic grouping of zones.

Multi-AZ means you run copies in more than one zone in one region. A zone loss must not stop the service if the design is honest. You need:

- Instances in more than one zone
- A store that can fail over across those zones
- A balancer that does not live in one zone only
- Tested failover (Topic 15)

Multi-region means you run in more than one region. Distance adds latency (Topic 11). Data copy adds RPO questions (Topic 15). Active-passive region failover is simpler than active-active (Topic 21).

Pick multi-AZ before multi-region. Many products need zone resilience. Few student systems need a second continent.

Law can force a region (data location). That is a constraint, not a scale pattern.

Cost doubles easily. Cross-zone and cross-region traffic has a price. Write the bill in the ADR.

Do not claim multi-AZ if the database is a single node in one zone. The API copies do not save you.

Name the RTO for a region loss. If the answer is "we accept a day", a backup in a second region can be enough (Topic 15).

### Questions

#### Theoretical questions

1. What is an availability zone in this handbook?
2. What is a region?
3. What must be true for an honest multi-AZ claim?
4. Why do you pick multi-AZ before multi-region?
5. Why can API copies in two zones still fail if the database has one zone?

#### Easy practical tasks

1. Write five sentences that compare multi-AZ and multi-region.
2. Make a table: "Event" and "Survives if multi-AZ is honest?". Add one zone loss and one region loss.
3. List four cost items that appear when you add a region.
4. Draw two zones, two API instances, and a multi-AZ database.

#### Medium practical tasks

1. Write a one-page multi-AZ design for a campus shop. Include the store.
2. Explain in eight sentences how a region failover with RPO 15 minutes looks to a user who just paid.
3. Write an ADR that rejects multi-region for a campus-only app. Include the law if data must stay local.

#### Advanced practical tasks

1. Design RTO/RPO for zone loss versus region loss. Fill a table with tests that you would run.
2. Compare a second-region backup with a hot standby. Write people cost and money cost.

---

## Managed databases and brokers

A managed database or broker is a vendor-operated store. You choose size, network rules, and backup windows. The vendor applies many OS patches and runs the control plane. You still design schema, queries, and users.

Why architects pick managed stores:

- Consensus and failover are already built (Topic 11)
- Backups and restore APIs exist (Topic 15)
- The team can stay small

You still own:

- Schema and expand-contract (Topic 18)
- Privileges (Topic 14)
- Capacity and slow queries (Topic 2)
- The bill

A managed broker does not remove poison messages or dual write (Topic 10). A managed database does not remove ownership (Topic 9).

Limits exist: maximum connections, maximum size, forbidden features, and region list. Read the limits before you promise a design.

Do not run a homemade cluster to "learn Raft" in production (Topic 11). Learn on a lab. Ship a managed primary.

Exit cost is lock-in (next section). Export and backup formats matter.

Network placement matters. A public database endpoint is a trust-boundary defect (Topic 14). Use private networking when the platform offers it.

### Questions

#### Theoretical questions

1. What is a managed store in this handbook?
2. What do you still own?
3. What does a managed broker not remove?
4. Why is a homemade production cluster a risk for a small team?
5. Why must you read platform limits before you promise a design?

#### Easy practical tasks

1. Write four sentences that define a managed database.
2. Make a table: "Task" and "Vendor or you". Add OS patch, schema, backup schedule, and query plan.
3. List five limits that you will look up on a vendor page.
4. Write four sentences on a public database endpoint as a trust defect.

#### Medium practical tasks

1. Write a one-page managed PostgreSQL choice: size, multi-AZ, backup, and private network.
2. Explain in eight sentences how a connection limit becomes an application architecture problem.
3. Compare a managed broker with a containerized broker on one VM. Write operations cost.

#### Advanced practical tasks

1. Write an ADR: managed store versus self-operated store. Include lock-in and skill.
2. Design an export drill: take a backup out of the vendor and restore it somewhere that you control (lab). Write the clock.

---

## Cost as an architectural constraint

Cost is a quality attribute (Topic 2). In the cloud, cost is also a monthly constraint that can kill a design. Idle resources still bill. Egress bills. Extra zones bill.

Write cost with a unit and a period. Include people hours. A cheap machine that needs a full-time operator is not cheap.

Architecture choices that move cost:

- Many services (Topic 12)
- Chatty cross-AZ traffic
- Unbounded logs and traces (Topic 16)
- Oversize managed stores
- Always-on FaaS replacements that would be cheaper as one process
- Forgotten environments (old preview stacks)

Set a budget owner. Alert on forecast, not only on surprise invoices.

Premature distribution increases cost. Start with one process and one managed database if that meets SLOs.

Reserved capacity and committed use can lower money cost and raise lock-in. Write the trade.

Students and small teams must also count time-to-learn. A platform with 200 products can waste a semester.

Do not optimize a 20-dollar bill with a 200-hour rewrite. Optimize the item that dominates.

### Questions

#### Theoretical questions

1. Why is cloud cost an architecture constraint?
2. Why must people hours appear in the cost number?
3. Name four design choices that raise a cloud bill.
4. Why is a forgotten preview environment a cost defect?
5. When is a rewrite a worse cost than the invoice?

#### Easy practical tasks

1. Write five sentences that define cost as a constraint.
2. Make a table: "Item" and "Unit". Add compute, egress, log ingest, and on-call hours.
3. List six idle resources that a student project can forget to stop.
4. Write a monthly budget sentence for a campus shop lab.

#### Medium practical tasks

1. Estimate a one-page bill for API, database, balancer, and logs at a stated request rate. State assumptions.
2. Explain in eight sentences how chatty microservices raise cross-AZ cost.
3. Write a cost-alert rule in words: threshold, window, and owner.

#### Advanced practical tasks

1. Write a year cost model: users × 10, and which component dominates.
2. Compare always-on instances with FaaS for a spiky campus workload. Show the break-even idea.

---

## Vendor lock-in trade-offs

Vendor lock-in is the extra cost to leave a vendor. Lock-in is not automatically a defect. A managed database that saves two people-years can be the correct lock-in.

Sources of lock-in:

- Proprietary APIs and data types
- Images and functions that run only on one platform
- IAM and network models
- Billing discounts that punish exit
- Skills that do not transfer

Reduce lock-in when exit is likely: open formats, twelve-factor backing-service URLs, standard SQL where it is enough, OpenTelemetry instead of a vendor-only SDK (Topic 16).

Accept lock-in when:

- The team is small
- The feature is not your differentiator
- Export exists and you tested it
- The SLO is easier to meet with the product

Write the exit sketch in the ADR. A one-page "how we would leave" is enough. If you cannot write the sketch, you do not understand the lock.

Multi-cloud as a default is often a fashion word. Two clouds double operations cost. Use two clouds when a constraint forces it (law, a customer, or a true DR vendor split).

Do not rewrite portable code for a semester to avoid a lock that you will never exercise.

### Questions

#### Theoretical questions

1. What is vendor lock-in?
2. Why is lock-in not automatically a defect?
3. Name four sources of lock-in.
4. When do you accept lock-in?
5. Why is multi-cloud a costly default?

#### Easy practical tasks

1. Write four sentences that define lock-in.
2. Make a table: "Choice" and "Exit cost (low/medium/high)". Add standard SQL, a proprietary function trigger, and object files in a common format.
3. List five questions that an ADR must answer about a vendor.
4. Write an exit sketch in five sentences for a managed Postgres.

#### Medium practical tasks

1. Write an ADR that accepts a vendor queue. Include the exit sketch and the reason you accept lock-in.
2. Explain in eight sentences how a proprietary identity API locks the application.
3. Compare "portable to two clouds" with "portable to one other PaaS". Write the real goal (avoid downtime versus avoid price).

#### Advanced practical tasks

1. Write a one-page lock-in register: service, lock type, export test date, and owner.
2. Read a public vendor migration guide outline. Write ten sentences on what makes exit feasible. Do not copy long passages.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do IaaS, PaaS, and FaaS change people cost without changing the need for SLOs?
2. How do twelve-factor rules and containers support the same deploy habit?
3. Why is multi-AZ a reliability design and multi-region often a different, later design?
4. How do managed stores, cost, and lock-in form one vendor decision?
5. Why does this path place cloud after delivery and reliability?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic.
2. For a to-do API, pick PaaS, one region, two AZs if the platform allows it, and a managed database. Write one lock-in sentence.
3. Draw a C4 deployment view: edge, API, database, and zone boxes.
4. Bookmark [https://12factor.net/](https://12factor.net/). Write one sentence on when you open it.

#### Medium practical tasks

1. Write a two-page infrastructure brief for a campus shop: grades of service, factors you enforce, AZ plan, managed store, budget, and lock-in.
2. Take a teammate slide that only lists product names. Rewrite it as failure domains, bills, and owners.
3. Map Topics 14, 15, and 18 onto one "first production on a cloud" plan.

#### Advanced practical tasks

1. Design a twelve-month platform path: PaaS now, container service later, reject raw Kubernetes until a named owner exists.
2. Compare two public reference architectures from vendors. Mark parts that your constraints would remove (Topic 1). Do not copy diagrams as if they were your design.
