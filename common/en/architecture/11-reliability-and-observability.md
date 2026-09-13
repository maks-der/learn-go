# 11. Reliability and Observability

## Description

Reliability is the ability of the system to do useful work in production. Observability is the ability to explain the current behavior of the system from its outputs. This topic covers SLI, SLO, SLA, error budgets, redundancy, degradation, backpressure, RPO and RTO, logs, metrics, traces, correlation identifiers, RED and USE, OpenTelemetry at a high level, and alerts without fatigue.

Architecture without operations is a wish (Topic 1). Topic 8 covers partial failure. This topic adds targets, recovery, and signals. Complete Topics 1 to 10 before this topic. Pair with [https://sre.google/books/](https://sre.google/books/) for deeper reading. Stay defensive. Do not run experiments on systems that you do not own.

Use one term for each concept. An SLI is a measure. An SLO is a target. An SLA is a promise to a customer. RPO is not RTO. A log is an event record. A metric is a number over time. A trace is a path of spans. A dashboard is for humans who explore. An alert is a page that demands action.

---

## SLI, SLO, SLA, and error budgets

An SLI (service level indicator) is a quantitative measure of a user-visible quality. Examples: the share of HTTP requests that succeed, and the share of requests that finish under 300 ms.

An SLO (service level objective) is a target for an SLI over a window. Example: 99.9 percent of `GET /catalog` requests in 30 days return 2xx or 4xx from the application and finish under 300 ms on the server.

An SLA (service level agreement) is a contractual promise. It includes remedies (credits) when the promise fails. Most internal student systems have SLOs and no SLA. Do not write an SLA that the team cannot measure.

Rules for a useful SLI:

- It matches a user job, not only a CPU graph.
- You can collect it in production.
- You state the event (request, message, batch job) and the success rule.

Availability SLIs often exclude 4xx that the client caused. A user typo is not your downtime. A 500 is. Write the rule. Latency SLIs use a percentile, not only the average (Topic 1).

An error budget is the amount of unreliability that an SLO allows. If the SLO is 99.9 percent success in 30 days, the budget is 0.1 percent of events in that window.

The budget is a control for change speed. When budget remains, you can ship. When budget is gone, you stop risky change and you repair reliability. This rule needs an owner. Product and operations must accept it (Topic 1).

Calculate the budget in events, not only in percent. "0.1 percent of 1 000 000 requests" is 1 000 failed requests. People understand counts.

Write what happens at 50 percent budget consumed and at 100 percent consumed. Example: extra review, feature freeze, or a reliability week.

Too many SLOs hide the few that matter. Start with two: availability of the primary user journey, and latency of that journey.

Google SRE books discuss these terms in depth. See [https://sre.google/books/](https://sre.google/books/). Use the books after you can name one SLI on your own system.

### Questions

#### Theoretical questions

1. What is an SLI?
2. What is an SLO?
3. What is an SLA?
4. What is an error budget?
5. Why must an SLI match a user job?

#### Easy practical tasks

1. Write five sentences that separate SLI, SLO, SLA, and error budget.
2. Make a table: "Term" and "Example". Add one row for each term.
3. Write two SLIs for a library catalog: success rate and latency.
4. Convert 0.1 percent of 1 000 000 requests to a count of failed requests.

#### Medium practical tasks

1. Write SLOs for checkout: availability and P99 latency. Include window and operation name.
2. Explain in eight sentences why an average latency SLO can hide a bad tail.
3. Write what the team does at 50 percent and at 100 percent budget consumed.

#### Advanced practical tasks

1. Read the SLO chapter overview in a public SRE book index. Write ten sentences that map the book terms to this section. Do not copy long passages.
2. Design three SLIs for a worker that sends mail: success, lag, and poison rate. Write how you would measure each.

---

## Redundancy, degradation, backpressure, RPO/RTO

Redundancy is extra capacity that can take over when a part fails. Examples: two API instances, a standby database, two availability zones. Redundancy is not a backup. A copy that you never fail over to is a hope.

Failover is the act of switching to the redundant part. You must test failover. An untested failover is a wish (Topic 1).

Graceful degradation is a planned lesser behavior when a dependency is down (Topic 8). The catalog page can hide search. Checkout must not hide payment and still claim success.

Backpressure is a signal that a consumer cannot accept more work. The producer must slow down, drop with a rule, or fail. Without backpressure, queues grow until memory or disk dies. Timeouts and bounded queues are forms of backpressure. A circuit breaker is a related guard (Topic 8).

RPO (recovery point objective) is the maximum amount of data that you can accept to lose, measured in time. An RPO of 5 minutes means a backup or a replica that is at most 5 minutes behind.

RTO (recovery time objective) is the maximum time to restore useful service after a disaster. An RTO of 1 hour means you practiced the restore in less than 1 hour.

Write RPO and RTO for the source of truth (Topic 6). A cache can be rebuilt. The source cannot. Pair with `db.topics.md` for backup practice.

Do not promise a low RTO if nobody ran a restore. Do not promise a low RPO if you have no replica and no backup schedule.

Chaos experiments (late and careful) test redundancy. Run them only on systems that you own. Start on paper and in a lab.

### Questions

#### Theoretical questions

1. What is redundancy?
2. What is graceful degradation?
3. What is backpressure?
4. What is RPO?
5. What is RTO?

#### Easy practical tasks

1. Write five sentences that define the terms in this section.
2. Make a table: "Part" and "Redundant copy?". Add API instance, database, and cache.
3. Write an RPO and an RTO for a student notes database in one sentence each.
4. List four degrade behaviors for a shop home page.

#### Medium practical tasks

1. Design checkout when payments are down versus when mail is down. They must differ.
2. Write a bounded-queue policy for a mail worker: max depth and what the API returns when the queue is full.
3. Explain in eight sentences why an untested failover is not reliability.

#### Advanced practical tasks

1. Write a one-page disaster-recovery sheet: RPO, RTO, backup owner, last restore test date.
2. Design a paper game day: one instance down, then the primary database down. List expected user results. Do not attack systems that you do not own.

---

## Logs, metrics, traces, correlation IDs

Logs, metrics, and traces are the three common signal types.

A log is a structured event. A useful log has a time, a level, a message, a service name, and identifiers. Prefer structured fields (JSON or key-value) over free text. Bound size. Do not log secrets (Topic 10). Do not log entire request bodies if they can hold personal data.

A metric is a number in a time series: request count, error count, latency histogram, queue depth, and CPU. Metrics are cheap to aggregate. They are poor at telling a unique story. High cardinality (a new series per user identifier) can break the metrics system.

A trace is a tree of spans for one request across processes. Each span has a name, times, and attributes. Traces explain where time went. You sample traces if volume is high. Sampling must still catch errors.

Use all three. Metrics tell you that error rate rose. Logs tell you the error message. Traces tell you which hop was slow.

A correlation identifier is a value that joins records for one unit of work. Common names: request identifier, trace identifier, and causation identifier.

Rules for correlation:

- Create an identifier at the edge if the client did not send one.
- Accept a client identifier only after you validate format and size.
- Never put secrets in the identifier.
- Show the identifier on error pages that operators may ask for.
- Put the identifier in messages and outbound calls (Topic 7).

Without correlation, you cannot prove that a mail worker handled the same order that the API accepted.

Write logs at the boundary: request start, request end, and handled failures. Debug logs in a hot loop fill disks and hide events.

Retention is a design choice. Legal rules can force a minimum or a maximum (Topic 1). Cost rises with volume.

Do not treat stdout screenshots as a production strategy. Ship logs to a store that you can search.

### Questions

#### Theoretical questions

1. What is a log in this handbook?
2. What is a metric?
3. What is a trace?
4. What is a correlation identifier?
5. Why is a new metric series per user identifier a risk?

#### Easy practical tasks

1. Write five sentences that define the three signal types and correlation.
2. Make a table: "Signal" and "Example". Add one log, one metric, and one span attribute.
3. List six fields that every request log must include. Exclude secrets.
4. Write a bad log line and a better structured log line for a failed loan (no personal extras).

#### Medium practical tasks

1. Design a minimum set for an API and a worker: five metrics, five log events, and one trace path.
2. Explain in eight sentences how a histogram metric differs from an average.
3. Write a retention rule for logs and traces. Include a law constraint if your campus has one.

#### Advanced practical tasks

1. Map one user action (place a hold on a book) across logs, metrics, and a trace. Draw the join keys.
2. Write a one-page telemetry standard: structure, cardinality limits, and sampling of traces.

---

## RED / USE and OpenTelemetry (high-level)

RED and USE are two simple methods to pick metrics.

RED is for request-driven services:

- Rate: requests per time
- Errors: failed requests per time (or a ratio)
- Duration: latency distribution

USE is for resources (CPU, disk, pool, broker):

- Utilization: busy share
- Saturation: extra work waiting
- Errors: resource errors

Use RED on the user-facing API. Use USE on the database, the disk, and the queue. Together they support SLIs and capacity work (Topic 12).

OpenTelemetry is a set of standards and tools for traces, metrics, and logs. The project is vendor-neutral. You instrument the application once. You can export to different backends. See the public OpenTelemetry documentation for current APIs.

This path needs a high-level view only:

- A span is a unit of work in a trace.
- Context propagation copies the trace identifier across HTTP and messages.
- An SDK and an exporter send data to a collector or a backend.
- Auto-instrumentation can cover common libraries. You still add business attributes (order id, not a secret).

OpenTelemetry is not a vendor product. A vendor can implement the standard. Do not treat a vendor slide as the standard.

Start small. One trace path through the API and the database beats a full platform that nobody reads.

Do not add high-cardinality attributes on every span. The same cardinality rule as metrics applies.

### Questions

#### Theoretical questions

1. What does RED measure?
2. What does USE measure?
3. When do you apply RED versus USE?
4. What is OpenTelemetry at a high level?
5. What is context propagation?

#### Easy practical tasks

1. Write five sentences about RED, USE, and OpenTelemetry.
2. Make a table: "Metric" and "RED or USE". Add six rows.
3. List four span attributes that are safe and two that must not appear (secrets, raw personal payloads).
4. Draw one trace: browser, API span, database span. Label the trace identifier.

#### Medium practical tasks

1. Design a RED set for `POST /loans` and a USE set for the database connection pool.
2. Write an ADR: use OpenTelemetry SDKs, pick one backend, no second instrumentation API.
3. Explain in eight sentences how context propagation helps a saga debug (Topic 9).

#### Advanced practical tasks

1. Write a one-page instrumentation standard: required spans, attribute deny list, and sampling.
2. Open the public OpenTelemetry documentation index. Write ten sentences that map terms to this section. Do not copy long passages.

---

## Alerts without fatigue

An alert is a notification that demands a human action. A dashboard is a view for a human who already looks. Do not page people from every dashboard chart.

Alert fatigue is the state where people ignore pages because most pages are noise. Fatigue is an architecture defect. It hides a real outage.

Rules for a useful alert:

- It maps to a user-visible SLO or to a fast-moving precursor (disk full, certificate expiry, queue depth).
- A person can do something. "CPU is 61 percent" is often not an action.
- It has an owner and a runbook of five to ten steps.
- It is not a duplicate of three other alerts for the same fault.

Prefer burn-rate or error-budget alerts on the SLI. A single spike of latency can be a blip. A fast burn of the monthly budget is a page.

Tune. If an alert fires and nobody acts, delete it or fix it. Record the change.

On-call is a cost (Topic 1). More services mean more alerts unless you invest in quality (Topic 9). A modular monolith can have a small alert set: API SLO, database up, disk, and the one queue.

Do not alert on every 404. Do not alert on client errors that you already excluded from the SLO.

Show the correlation identifier in the alert body when you can. Operators at 03:00 need one paste value.

### Questions

#### Theoretical questions

1. What is an alert in this handbook?
2. How does an alert differ from a dashboard?
3. What is alert fatigue?
4. What makes an alert useful?
5. Why can more services increase fatigue?

#### Easy practical tasks

1. Write five sentences about alerts without fatigue.
2. Make a table: "Signal" and "Page? (yes/no)". Add SLO burn, CPU 61 percent, disk 95 percent, and one 404.
3. List four items in a minimum runbook (who, what to check, how to degrade, who to call).
4. Write four sentences on why a duplicate alert is a defect.

#### Medium practical tasks

1. Design three alerts for a campus catalog: SLO burn, database down, and queue depth. Write the action for each.
2. Write an ADR: page on SLO burn, not on raw CPU, for the student app.
3. Explain in eight sentences how you retire an alert that nobody acknowledges.

#### Advanced practical tasks

1. Write a one-page alert standard: severity, owner, runbook, and review every 30 days.
2. Map five Topic 8 failure modes to one alert or to a dashboard-only signal.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do SLI, SLO, and error budget use logs, metrics, and traces together?
2. Why must RPO and RTO name the source of truth, not the cache?
3. How do RED and USE support an SLO without a metric explosion?
4. Why does a correlation identifier matter more after you add a broker (Topic 7)?
5. How does alert fatigue weaken the reliability quality attribute from Topic 1?

#### Easy practical tasks

1. Write a one-page cheat sheet: SLI/SLO/SLA/budget, redundancy, degrade, backpressure, RPO/RTO, three signals, RED/USE, OpenTelemetry, alert rules.
2. For a to-do API, write two SLOs, four log fields, three metrics, and two alerts.
3. Draw a request with a correlation identifier through API, database, and mail worker.
4. Bookmark the Google SRE books page. Write one sentence on when you open it.

#### Medium practical tasks

1. Write a short operations brief for a campus shop: SLOs, backup RPO/RTO, telemetry minimum, and a three-alert set.
2. Take a teammate design that pages on every 5xx. Rewrite it as an SLO burn alert plus a dashboard.
3. Write a twelve-week plan: weeks for structured logs and a request id, then SLOs.

#### Advanced practical tasks

1. Write a game-day plus observability plan (paper): dependency down, expected SLI change, expected alert, expected log. Do not run it on systems that you do not own.
2. Map this topic to the suggested practice order in `architecture.topics.md` (structured logs, SLOs, failure modes) in a one-page table.
