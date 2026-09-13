# 16. Observability

## Description

Observability is the ability to explain the current behavior of the system from its outputs. This topic covers logs, metrics, and traces; correlation identifiers; RED and USE methods; OpenTelemetry at a high level; dashboards versus alerts; and alert fatigue.

Topic 2 defined observability as a quality attribute. Topic 15 needs measures for SLIs. This topic shows the signals and the habits. Complete Topics 1 to 15 before this topic. Design signals with the feature. Do not add them only after the first incident.

Use one term for each concept. A log is an event record. A metric is a number over time. A trace is a path of spans. A correlation identifier joins records. A dashboard is for humans who explore. An alert is a page that demands action. OpenTelemetry is a standard for telemetry, not a vendor product.

---

## Logs, metrics, traces

Logs, metrics, and traces are the three common signal types.

A log is a structured event. A useful log has a time, a level, a message, a service name, and identifiers. Prefer structured fields (JSON or key-value) over free text. Bound size. Do not log secrets (Topic 14). Do not log entire request bodies if they can hold personal data.

A metric is a number in a time series: request count, error count, latency histogram, queue depth, and CPU. Metrics are cheap to aggregate. They are poor at telling a unique story. High cardinality (a new series per user identifier) can break the metrics system (Topic 2).

A trace is a tree of spans for one request across processes. Each span has a name, times, and attributes. Traces explain where time went. You sample traces if volume is high. Sampling must still catch errors.

Use all three. Metrics tell you that error rate rose. Logs tell you the error message. Traces tell you which hop was slow.

Write logs at the boundary: request start, request end, and handled failures. Debug logs in a hot loop fill disks and hide events.

Retention is a design choice. Legal rules can force a minimum or a maximum (Topic 1). Cost rises with volume (Topic 2).

Do not treat stdout screenshots as a production strategy. Ship logs to a store that you can search.

### Questions

#### Theoretical questions

1. What is a log in this handbook?
2. What is a metric?
3. What is a trace?
4. Why do you need all three signal types?
5. Why is a new metric series per user identifier a risk?

#### Easy practical tasks

1. Write five sentences that define the three signal types.
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

## Correlation IDs

A correlation identifier is a value that joins records for one unit of work. Common names: request identifier, trace identifier, and causation identifier.

A request identifier labels one incoming HTTP call. Every log line for that call repeats the value. If the API starts a worker job, the job log repeats the same value or a child value that points to the parent.

A trace identifier labels a distributed trace. Span identifiers label each hop. Context propagation copies these values across HTTP headers or message headers (Topic 10).

Rules:

- Create an identifier at the edge if the client did not send one.
- Accept a client identifier only after you validate format and size.
- Never put secrets in the identifier.
- Show the identifier on error pages that operators may ask for.
- Put the identifier in messages and outbound calls.

Without correlation, you cannot prove that a mail worker handled the same order that the API accepted. Dual-write and saga debug become guesswork (Topics 10 and 12).

Operators at 03:00 need one paste value. Train the team to ask for that value first.

Do not invent a new header name if a standard that you already use defines one. OpenTelemetry and many gateways already define headers (next section). Pick one convention and write it down.

### Questions

#### Theoretical questions

1. What is a correlation identifier?
2. How does a request identifier differ from a trace identifier?
3. When do you create an identifier at the edge?
4. Why must a worker log the identifier from the API?
5. Why do operators want one paste value?

#### Easy practical tasks

1. Write four sentences that define correlation identifiers.
2. Make a table: "Hop" and "Header or field". Add browser, API, message, and worker.
3. List five places that must print the identifier.
4. Write a user-facing error sentence that includes a request identifier.

#### Medium practical tasks

1. Design identifier rules for HTTP and for a queue message. Include size and charset.
2. Explain in eight sentences how a saga (Topic 12) uses a correlation identifier and a causation identifier.
3. Review a homework log with no identifier. Write the first three changes.

#### Advanced practical tasks

1. Write a one-page propagation standard: incoming header, log field, outbound header, and message field.
2. Draw a sequence for one order. Show parent and child identifiers. Write how you search the log store.

---

## RED / USE methods

RED and USE are two simple methods to choose metrics.

RED fits request-driven services:

- Rate: requests per second
- Errors: failed requests per second or error ratio
- Duration: latency distribution

USE fits resources (CPU, disk, queue, thread pool):

- Utilization: the share of time that the resource is busy
- Saturation: extra work that waits (queue depth)
- Errors: device or pool errors

Use RED on the SLI path (Topic 15). Use USE when you debug a tight resource. A high duration with low utilization can mean a remote wait. A high utilization with rising saturation means you are out of capacity.

Do not collect every kernel metric on day one. Start with RED for the user journey and USE for the database and the broker.

Name metrics so that a new operator can guess the meaning. `http_requests_total` with labels `route` and `status` is clearer than `x1`.

Label with care. Route, method, and status are common. User identifier as a label is usually too much cardinality.

The methods do not replace traces. They tell you where to open a trace.

### Questions

#### Theoretical questions

1. What does RED measure?
2. What does USE measure?
3. When do you apply RED?
4. When do you apply USE?
5. Why must user identifier stay out of metric labels in most systems?

#### Easy practical tasks

1. Write five sentences that compare RED and USE.
2. Make a table: "Method" and "Three numbers". Fill RED and USE.
3. List four resources that need USE in a small shop.
4. Write three metric names for `POST /loans` in RED style.

#### Medium practical tasks

1. Design a RED dashboard for one API and a USE dashboard for its database.
2. Explain in eight sentences how high duration plus low CPU utilization points to a dependency.
3. Write label rules: allowed labels and forbidden labels.

#### Advanced practical tasks

1. Write a one-page metric catalog for API, worker, broker, and database. Mark RED or USE for each line.
2. Compare RED with a custom business metric (orders placed). Write how both appear in one SLO story.

---

## OpenTelemetry (high-level)

OpenTelemetry is a set of standards and libraries for traces, metrics, and logs. The project lives at [https://opentelemetry.io/](https://opentelemetry.io/). This section is high level. It is not an SDK tutorial.

The useful ideas:

- A common data model for spans, metrics, and logs
- Context propagation across process boundaries
- A collector that can receive telemetry and export it to a backend
- Instrumentation libraries that wrap HTTP, gRPC, and common stores

Architecture choices:

- What you instrument (your code and your edge)
- What you sample
- Where the collector runs
- Which backend stores traces and metrics
- How much you spend (Topic 2)

OpenTelemetry does not replace the need for questions. You still name the 03:00 questions (Topic 2). The standard reduces vendor lock-in of the instrumentation API. The storage vendor can still lock data (Topic 19).

Do not enable every auto-instrumentation on a huge monolith on day one. You will drown in spans. Start with the inbound HTTP span and the outbound HTTP or SQL span.

Keep attributes bounded. Do not attach full payloads.

A correlation identifier that matches the trace identifier simplifies search. Align your log field with the trace context when you can.

If the team cannot run a collector, a simpler stack (structured logs plus one metrics agent) is still observability. Adopt OpenTelemetry when the extra hops justify the common API.

### Questions

#### Theoretical questions

1. What is OpenTelemetry at a high level?
2. What is context propagation?
3. What job does a collector do?
4. Why can storage still create lock-in if the API is standard?
5. Why must you not enable every auto-instrumentation at once?

#### Easy practical tasks

1. Open [https://opentelemetry.io/](https://opentelemetry.io/). Write the three signal types that the home page names.
2. Write four sentences that define OpenTelemetry for a beginner.
3. Make a table: "Part" and "Job". Add SDK, collector, and backend.
4. List three first spans that you would enable on an API.

#### Medium practical tasks

1. Write a one-page adoption plan: inbound HTTP, outbound SQL, collector, and one backend.
2. Explain in eight sentences how sampling changes what you can debug.
3. Write an ADR: OpenTelemetry SDK versus only vendor agents. Include team skill.

#### Advanced practical tasks

1. Read the public OpenTelemetry overview on traces. Write ten sentences on context and spans. Do not copy long passages.
2. Design attribute limits and a deny list (payloads, secrets). Stay defensive.

---

## Dashboards vs alerts

A dashboard is a screen of charts and tables that a human reads. An alert is a condition that notifies a person when it stays true.

Dashboards answer "what is going on?" during a review or an incident. Alerts answer "must a human act now?"

Good dashboard habits:

- Start from RED for the user journey
- Show error ratio, latency, and saturation
- Link to logs and traces with the same time range
- Keep one screen for the primary journey

Good alert habits:

- Alert on symptoms that users feel, or on a bound that predicts user pain
- Require a duration (not a single spike)
- Name the next step in the alert text
- Page a human only when that human can act

A chart without an owner is decoration. An alert without a runbook is noise.

Do not alert on every dashboard line. That path causes alert fatigue (next section).

SLO burn alerts (Topic 15) are often better than raw CPU alerts. CPU can be high while users are fine. Users can be in pain while CPU is low.

Static thresholds need review when load changes. A fixed "100 requests per second" alert can fire every exam week and stay silent in summer.

### Questions

#### Theoretical questions

1. What is a dashboard for?
2. What is an alert for?
3. Why must an alert include a next step?
4. Why can a CPU alert miss user pain?
5. Why do static thresholds fail when load is seasonal?

#### Easy practical tasks

1. Write five sentences that separate dashboards from alerts.
2. Make a table: "Signal" and "Dashboard, alert, or both". Add P99 latency, disk 95 percent full, and a debug counter.
3. List four panels for a primary-journey dashboard.
4. Write one alert in words: condition, duration, and next step.

#### Medium practical tasks

1. Design a dashboard plus three alerts for a catalog API. Include an SLO burn idea.
2. Explain in eight sentences why paging on a single 5-second spike is a defect.
3. Write a runbook header that an alert must link to: check, decide, act, communicate.

#### Advanced practical tasks

1. Write a one-page dashboard standard: one journey screen, one resource screen, and naming rules.
2. Compare symptom alerts with cause alerts. Write when each is valid.

---

## Alert fatigue

Alert fatigue is a state where people ignore alerts because too many alerts are useless. The system still pages. Humans stop trusting the page.

Causes:

- Alerts on noisy metrics
- No duration
- Duplicate alerts for one fault
- Alerts that clear and fire in a loop
- Alerts that no one can act on
- Missing ownership

Effects:

- Slow response to a real incident
- Silent disable of alerts
- Loss of the error-budget loop (Topic 15)

Repairs:

- Delete or downgrade alerts that did not need a human in the last N weeks
- Group alerts that share one cause
- Use warning channels for slow trends and pages for user-visible breaks
- Tune thresholds after a season change
- Write a weekly alert-review habit

A small team can support a small set of pages. Write a maximum. Example: five paging alerts for a student shop. Everything else is a dashboard or a ticket.

Night pages need a higher bar than weekday tickets. Protect sleep. Tired operators cause defects.

Do not add an alert as a souvenir of an incident without a stop condition. Many teams add alerts after pain and never remove them.

Measure fatigue: pages per week, percent that led to an action, and time to acknowledge. Those numbers are operations quality.

### Questions

#### Theoretical questions

1. What is alert fatigue?
2. Name four causes.
3. How does fatigue break incident response?
4. Why must a small team cap paging alerts?
5. Why do night pages need a higher bar?

#### Easy practical tasks

1. Write four sentences that define alert fatigue.
2. Make a table: "Alert" and "Page or ticket". Add disk 70 percent, SLO burn, and a missing optional widget.
3. List five repair actions from this section.
4. Write a weekly review checklist with four questions.

#### Medium practical tasks

1. Take a fictional list of 20 alerts. Mark keep, downgrade, or delete. Write one reason each for six of them.
2. Explain in eight sentences how duplicate alerts for one database fault appear.
3. Write a policy: maximum paging alerts and the approval to add one.

#### Advanced practical tasks

1. Design an alert-quality report: pages, action rate, and top noisy alerts. Include who reads it.
2. Compare a two-person on-call with a twenty-person SRE group. Write how the paging set must differ.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do logs, metrics, traces, and correlation identifiers answer one 03:00 question together?
2. How do RED and USE choose metrics without a random shopping list?
3. What does OpenTelemetry standardize, and what must you still decide?
4. How do dashboards, alerts, and fatigue control the same telemetry in different ways?
5. Why does this path place observability after SLOs and not only after quality attributes?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic.
2. For a to-do API, write five log fields, three RED metrics, one USE metric, and one alert in words.
3. Draw a request path with a correlation identifier across API and worker.
4. Bookmark [https://opentelemetry.io/](https://opentelemetry.io/). Write one sentence on when you open it.

#### Medium practical tasks

1. Write a two-page observability brief for a campus shop: signals, identifiers, dashboards, alerts, and a fatigue cap.
2. Take a teammate design that only says "we will add ELK and Grafana". Rewrite it as questions, signals, and owners.
3. Map Topics 2, 12, and 15 onto one checkout path. Write the SLI and the three signals that prove it.

#### Advanced practical tasks

1. Write a tabletop (60 minutes) that uses only telemetry to find a slow payment hop. Include identifiers and RED charts.
2. Read a public OpenTelemetry or SRE observability overview. Write ten sentences that you would add to a team standard. Do not copy long passages.
