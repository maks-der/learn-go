# 15. Reliability and Operations

## Description

Reliability and operations are the practices that keep a system doing useful work in production. This topic covers service level indicators, objectives, and agreements (SLI, SLO, SLA), error budgets, redundancy and failover, graceful degradation, backpressure, disaster recovery with RPO and RTO, and chaos experiments.

Architecture without operations is a wish (Topic 1). Topic 2 defines reliability and availability. Topic 11 covers partial failure. This topic adds targets, budgets, and recovery. Complete Topics 1 to 14 before this topic. Pair with [https://sre.google/books/](https://sre.google/books/) for deeper reading. Stay defensive. Do not run experiments on systems that you do not own.

Use one term for each concept. An SLI is a measure. An SLO is a target. An SLA is a promise to a customer. RPO is not RTO. Degradation is not failover. Chaos work is late and careful.

---

## SLI, SLO, SLA

An SLI (service level indicator) is a quantitative measure of a user-visible quality. Examples: the share of HTTP requests that succeed, and the share of requests that finish under 300 ms.

An SLO (service level objective) is a target for an SLI over a window. Example: 99.9 percent of `GET /catalog` requests in 30 days return 2xx or 4xx from the application and finish under 300 ms on the server.

An SLA (service level agreement) is a contractual promise. It includes remedies (credits) when the promise fails. Most internal student systems have SLOs and no SLA. Do not write an SLA that the team cannot measure.

Rules for a useful SLI:

- It matches a user job, not only a CPU graph.
- You can collect it in production (Topic 16).
- You state the event (request, message, batch job) and the success rule.

Availability SLIs often exclude 4xx that the client caused. A user typo is not your downtime. A 500 is. Write the rule.

Latency SLIs use a percentile, not only the average (Topic 2). State the operation and the hop.

Too many SLOs hide the few that matter. Start with two: availability of the primary user journey, and latency of that journey. Add a freshness SLI if a report can lag.

Google SRE books discuss these terms in depth. See [https://sre.google/books/](https://sre.google/books/). Use the books after you can name one SLI on your own system.

### Questions

#### Theoretical questions

1. What is an SLI?
2. What is an SLO?
3. What is an SLA?
4. Why must an SLI match a user job?
5. Why do many student systems have an SLO and no SLA?

#### Easy practical tasks

1. Write five sentences that separate SLI, SLO, and SLA.
2. Make a table: "Term" and "Example". Add one row for each term.
3. Write two SLIs for a library catalog: success rate and latency.
4. List four events that must not count as your downtime (client 4xx and similar).

#### Medium practical tasks

1. Write SLOs for checkout: availability and P99 latency. Include window and operation name.
2. Explain in eight sentences why an average latency SLO can hide a bad tail.
3. Draft a one-page SLO sheet for a campus login plus a downstream API. State what you exclude.

#### Advanced practical tasks

1. Read the SLO chapter overview in a public SRE book index. Write ten sentences that map the book terms to this section. Do not copy long passages.
2. Design three SLIs for a worker that sends mail: success, lag, and poison rate. Write how you would measure each (Topic 16).

---

## Error budgets

An error budget is the amount of unreliability that an SLO allows. If the SLO is 99.9 percent success in 30 days, the budget is 0.1 percent of events in that window.

The budget is a control for change speed. When budget remains, you can ship. When budget is gone, you stop risky change and you repair reliability. This rule needs an owner. Product and operations must accept it (Topic 1).

A budget is not a license to cause failures. It is a shared number that makes trade-offs visible. A feature that burns the budget in two days is a bad trade unless a stakeholder accepts the outage.

Calculate the budget in events, not only in percent. "0.1 percent of 1 000 000 requests" is 1 000 failed requests. People understand counts.

You need a clean SLI before you need a budget. A noisy SLI makes a noisy budget. Exclude client errors that you already excluded from the SLO.

Write what happens at 50 percent budget consumed and at 100 percent consumed. Example: extra review, feature freeze, or a reliability week.

Error budgets fail when leadership ignores them. An SLO without a budget policy is a poster.

### Questions

#### Theoretical questions

1. What is an error budget?
2. How does a budget control change speed?
3. Why is a budget not a license to cause failures?
4. Why must you express the budget in event counts as well as percent?
5. What happens to the policy if leadership ignores the budget?

#### Easy practical tasks

1. Compute the failed-request count for 99.0 percent and 99.9 percent on 500 000 requests.
2. Write four sentences that define an error budget.
3. Make a table: "Budget left" and "Policy". Add 100 percent, 50 percent, and 0 percent.
4. List three changes that burn budget (bad deploy, dependency timeout, bad schema).

#### Medium practical tasks

1. Write a one-page error-budget policy for a four-person team. Include who can stop a release.
2. Given 2000 allowed failures and 1800 already used, write the decision for a risky schema change.
3. Explain in eight sentences how a noisy 4xx classification can fake a healthy budget.

#### Advanced practical tasks

1. Design a monthly budget report: SLI, consumed budget, top three burners, and next actions.
2. Compare "feature freeze at zero budget" with "only freeze the hot path". Write stakeholder conflicts.

---

## Redundancy and failover

Redundancy is extra capacity that can continue work when one part fails. Failover is the act of switching from a failed part to a healthy part.

Examples:

- Two application processes behind a load balancer
- A database replica that can become primary
- Two availability zones (Topic 19)

Redundancy is not a number of boxes. If both processes sit on one host, one host failure still stops the service. Independent failure domains matter.

Failover can be automatic or manual. Automatic failover is faster and can be wrong (split brain, Topic 11). Manual failover is slower and needs a trained person. Write the choice in an ADR.

Test failover. An untested replica is a wish (Topic 2). A test includes:

- Detection time
- Switch time
- Data loss, if any
- Client behavior (retries, errors)

Health checks drive many failovers. A bad health check removes all healthy nodes or keeps a dead node in the pool. Keep the check simple. Check a real dependency if the process is useless without it. Do not make the check so heavy that it becomes the outage.

Active-passive means a standby waits. Active-active means more than one unit serves traffic. Active-active needs careful data design (Topics 9 and 21).

Cost rises with redundancy. Stakeholders must accept the bill (Topic 2).

### Questions

#### Theoretical questions

1. What is redundancy?
2. What is failover?
3. Why is a second process on the same host weak redundancy?
4. Why must you test failover?
5. What risk does a bad health check create?

#### Easy practical tasks

1. Write five sentences that define redundancy and failover.
2. Make a table: "Part" and "Independent domain?". Add two processes on one VM and two processes in two zones.
3. List four items that a failover test must record.
4. Draw active-passive for one database. Label the client.

#### Medium practical tasks

1. Write a failover runbook for a web process: detect, switch, verify, and communicate.
2. Explain in eight sentences how an automatic failover can cause split brain. Stay at idea level (Topic 11).
3. Design a health check for an API that needs the database. Write what you check and what you do not check.

#### Advanced practical tasks

1. Write an ADR: single-zone app with a multi-AZ database versus multi-AZ app. Include cost and RTO.
2. Compare active-passive and active-active for a session store. Write the data problem that active-active adds.

---

## Graceful degradation

Graceful degradation is a planned weaker behavior when a dependency fails. The system continues a subset of useful work. Topic 11 introduced the idea. This section makes it an operations contract.

Examples:

- Search is down: show a static popular list
- Recommendations are down: show the catalog
- Mail is down: accept the order and queue the mail
- Payments are down: refuse checkout with a clear message (do not pretend success)

Degradation is a product choice. A hidden empty page is not graceful. A clear message is. Write the user text in the same document as the failure matrix.

Each dependency needs a row:

- Fail the request
- Degrade
- Queue
- Serve a subset

Timeouts turn "unknown" into a local decision (Topic 4). That decision can be wrong. Degradation must be safe if the work later completes (idempotency, Topic 11).

Feature flags can turn off a non-critical widget (Topic 18). Cache can serve stale reads if the stakeholder accepts staleness (Topic 9).

Do not degrade security controls. A down policy store on an admin write must fail closed (Topic 14).

Measure degradation. A metric "search_degraded" tells you that users see the fallback. An SLO can allow fallback time or can treat fallback as a partial failure. Write the rule.

### Questions

#### Theoretical questions

1. What is graceful degradation?
2. Why is an empty page not graceful?
3. What four behaviors can a dependency row list?
4. Why must degradation stay safe if the remote work later completes?
5. Why must you not degrade an admin authorization check?

#### Easy practical tasks

1. Write four sentences that define graceful degradation.
2. Make a table: "Dependency" and "Behavior". Add search, mail, payments, and recommendations.
3. Write three user messages for three degrade paths. Use plain words.
4. List three features that you must not degrade.

#### Medium practical tasks

1. Write a failure matrix for a shop with four dependencies. Include user text and a metric name.
2. Explain in eight sentences how a stale cache is a degrade choice, not a silent lie, if you label freshness.
3. Design a flag that hides recommendations. Write who can flip it and how you record the flip.

#### Advanced practical tasks

1. Write a one-page degrade playbook: timeout, fallback, metric, and return-to-normal.
2. Compare "fail the whole page" with "render the page without one widget" for a campus portal. Write accessibility and support cost.

---

## Backpressure

Backpressure is a signal from a busy consumer to a producer to slow down. Without backpressure, queues grow until memory or disk dies. Latency explodes (Topic 2).

Forms of backpressure:

- HTTP 429 or 503 with a retry hint when the server is saturated
- A full queue that rejects new publish
- A bounded worker pool that does not accept more jobs
- TCP windowing at a lower layer (`net.topics.md`)

The application must still set bounds. An unbounded in-memory list of tasks is a defect. Topic 11 bulkheads isolate pools. Backpressure tells the caller that the pool is full.

Timeouts without backpressure cause retry storms. Callers retry, the server gets busier, and the outage grows. Combine timeout, retry budget, and a reject path (Topic 11).

User-facing APIs must fail fast when saturated. A 30-second wait that then fails is worse than a fast 503. Background producers can pause.

Write the limit: max queue depth, max concurrent handlers, max accept rate. Test the limit in a lab that you own. Do not flood a system that you do not own.

Backpressure is an architecture choice. If every service waits forever, one slow store knocks down the graph. If every service sheds load, users see errors, and the core can survive.

### Questions

#### Theoretical questions

1. What is backpressure?
2. What happens when queues have no bound?
3. Name three application-level forms of backpressure.
4. How do retries without a budget fight backpressure?
5. Why is a fast 503 often better than a long wait?

#### Easy practical tasks

1. Write five sentences that define backpressure.
2. Make a table: "Limit" and "Unit". Add queue depth, concurrent handlers, and accept rate.
3. List four symptoms of missing backpressure (memory growth, disk full, P99 rise, retry storm).
4. Write a user message for a saturated checkout.

#### Medium practical tasks

1. Design bounds for an API and a mail worker. Write what each unit does when the bound hits.
2. Explain in eight sentences how a bulkhead (Topic 11) and backpressure work together.
3. Write a lab plan on a system that you own: raise load until 503 appears. Do not test on foreign hosts.

#### Advanced practical tasks

1. Write a one-page saturation policy: shed order, metrics, and when operators add capacity.
2. Compare load shedding at the gateway with shedding in each service. Write the fairness problem.

---

## Disaster recovery (RPO/RTO)

Disaster recovery (DR) is the plan to restore a service after a large loss: a region down, a deleted store, or a corrupted disk. Daily failover of one process is not DR. DR assumes a bigger blast.

RPO (recovery point objective) is the maximum amount of data loss that stakeholders accept, measured in time. An RPO of 5 minutes means that a restore can lose up to 5 minutes of writes.

RTO (recovery time objective) is the maximum time from the disaster to useful service. An RTO of 4 hours means that users can work again within 4 hours.

Backups implement RPO. Restore drills implement RTO. A backup that no one restores is a wish.

Write:

- What you back up (database, object store, config)
- How often you back up
- Where copies live (another zone or region, Topic 19)
- How you restore
- Who decides to declare a disaster
- How you communicate to users

Ransomware and accidental `DROP` are DR cases. Least privilege and backups both matter (Topic 14).

Do not promise RPO zero if you cannot pay for synchronous copies across sites. Synchronous copies increase latency and cost (Topics 2 and 11).

Test restores on a schedule. Record the clock. Change the RTO if the drill misses the target.

### Questions

#### Theoretical questions

1. What is disaster recovery in this handbook?
2. What is RPO?
3. What is RTO?
4. Why is an unrestored backup a wish?
5. Why can RPO zero be the wrong promise?

#### Easy practical tasks

1. Write five sentences that separate RPO and RTO.
2. Make a table: "System" and "Example RPO/RTO". Add student notes and a bank-like transfer tool (fictional).
3. List six items that a DR document must name.
4. Write four sentences on why a region-level event needs a different plan than a single process crash.

#### Medium practical tasks

1. Write a one-page DR card for a campus catalog: backup, restore steps, owners, and user text.
2. Given RPO 15 minutes and a backup every 60 minutes, write why the design fails the RPO.
3. Plan a restore drill: environment, clock start, success check, and rollback of the drill.

#### Advanced practical tasks

1. Write an ADR: backups in a second region versus a second vendor. Include cost and lock-in (Topic 19).
2. Design RPO/RTO for two data stores (database and object files) that must stay consistent enough for the business. Write the gap that remains.

---

## Chaos experiments (careful, late)

A chaos experiment is a planned fault that you inject to learn if the system matches its reliability design. You do this late. You do this with consent. You do this on systems that you own or that you have written permission to test.

Chaos work is not random destruction. It is a hypothesis test. Example: "If we stop the search process, the catalog page still renders in 400 ms and the degrade metric rises."

Prerequisites:

- SLOs and dashboards exist (Topic 16)
- Timeouts and degrade paths exist
- You can stop the experiment
- Stakeholders know the window
- Production experiments start small, in working hours, with a rollback

Start in a lab. Kill one process. Break one dependency name in a test environment. Game days on paper also teach. A tabletop is a chaos experiment with no injector.

Do not run chaos against third-party systems. Do not run chaos as a joke. Do not hide the experiment from on-call.

If the experiment fails the hypothesis, you gained a defect list. That is success. If you only produce an outage without a write-up, you produced harm.

This path places chaos last in the topic. You must earn it with SLIs, backups, and degrade paths.

### Questions

#### Theoretical questions

1. What is a chaos experiment in this handbook?
2. Why must chaos wait until late?
3. What is a hypothesis in this context?
4. Why is a tabletop still useful?
5. What makes a failed hypothesis a success?

#### Easy practical tasks

1. Write four sentences that define a careful chaos experiment.
2. List eight prerequisites from this section.
3. Write one hypothesis for "mail worker down".
4. Make a table: "Allowed?" and "Action". Add lab process kill, paper game day, and unannounced production kill.

#### Medium practical tasks

1. Write a 90-minute game-day script: inject, watch SLIs, decide stop, write notes.
2. Explain in eight sentences why chaos without a degrade path only produces an outage.
3. Draft a permission note that you would send to a product owner before a lab experiment.

#### Advanced practical tasks

1. Design a one-page chaos policy for a school project: environments, approval, and forbidden targets.
2. Read a public write-up of a game day (not an attack guide). Write ten sentences on hypothesis and learning. Do not copy long passages.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do SLI, SLO, error budget, and release speed form one operations loop?
2. How do redundancy, degradation, and backpressure protect different kinds of failure?
3. Why must RPO and RTO appear in the same document as backup design?
4. Why does this path place chaos after SLOs and DR?
5. How do Topics 2 and 11 differ from this topic without replacing it?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic.
2. For a to-do API, write one SLI, one SLO, one budget policy sentence, and one degrade row.
3. Draw a poster: process down, database down, broker down, disk full. Write the first operator action for each.
4. Bookmark [https://sre.google/books/](https://sre.google/books/). Write one sentence on when you open it.

#### Medium practical tasks

1. Write a two-page operations brief for a campus shop: SLOs, budget, failover, degrade, RPO/RTO, and a rejected early chaos plan.
2. Take a teammate design that only says "we will be highly available". Rewrite it with numbers and tests.
3. Map a broker-down drill from Topic 10 onto SLIs, backpressure, and user text from this topic.

#### Advanced practical tasks

1. Write a tabletop (90 minutes) that walks a zone loss. Include failover, RPO, communication, and budget burn.
2. Compare this topic to one public SRE workbook chapter outline. Map five of their headings to your system. Do not copy long passages.
