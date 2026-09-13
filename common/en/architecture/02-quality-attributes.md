# 2. Quality Attributes

## Description

A quality attribute is a non-functional property of a system. Quality attributes include reliability, performance, scalability, security, maintainability, observability, and cost. This topic defines each attribute and shows how you measure it.

You cannot maximize all quality attributes at the same time. Each strong choice has a cost in another attribute. Architecture work is the selection of the attributes that matter and the attributes that you accept as weaker.

Complete Topic 1 before this topic. Use measurable statements. Do not use empty words such as "fast" or "secure" without a check.

Use one term for each concept. Availability is not the same as reliability. Latency is not the same as throughput.

---

## Reliability / availability

Reliability is the ability of the system to do the correct work for a period of time. A reliable system completes the intended function without a defect that stops the user. Mean time between failures (MTBF) is one measure of reliability.

Availability is the share of time that the system can accept and complete useful work. A common expression is:

`availability = uptime / (uptime + downtime)`

Teams also use "nines". "Three nines" means 99.9 percent availability. Three nines is about 8.8 hours of downtime in one year. "Four nines" is about 53 minutes in one year. More nines need more redundancy and more operational work.

Availability does not mean that every answer is correct. A system can be available and still return a wrong result. Reliability cares about correct work. Availability cares about the system being ready.

Planned maintenance is still downtime if users cannot work. Some teams exclude planned windows from the number. Write the rule. Users do not care about the internal rule. Users care about the time that they cannot work.

Faults occur. Disks fail. Processes stop. Dependencies time out. Reliability and availability improve when you detect a fault, isolate it, and recover. Later topics cover timeouts, redundancy, and graceful degradation.

Do not promise five nines if the team cannot test failover. A promise without a test is a wish.

### Questions

#### Theoretical questions

1. What is reliability in this handbook?
2. What is availability?
3. How do reliability and availability differ?
4. What does "three nines" mean for downtime in one year?
5. Why is a high availability number without a failover test a wish?

#### Easy practical tasks

1. Convert 99.0 percent, 99.9 percent, and 99.99 percent to approximate downtime hours or minutes in one year. Show the arithmetic.
2. Write four sentences that separate "the site answers" from "the site answers correctly".
3. List five faults that can decrease availability in a web application.
4. Make a two-column table: "Event" and "Counts as downtime? (yes/no)". Add six rows. Include planned maintenance.

#### Medium practical tasks

1. Write two measurable availability requirements for a library catalog: one for students in term time, one for a nightly batch.
2. Draw a simple diagram of one process and one database. Mark three faults. Write the user-visible effect of each fault.
3. Calculate the availability of two components in series if each component has 99.0 percent availability. State the assumption.

#### Advanced practical tasks

1. Compare MTBF, MTTR, and availability. Write a one-page note with one numeric example that uses all three terms.
2. Read a public status-page history for a month. Estimate availability from published incidents. List what the page does not tell you.

---

## Performance (latency, throughput)

Performance is how well the system uses time and capacity for a given workload. Two primary measures are latency and throughput.

Latency is the time from the start of a request to the completion of the response. Teams report percentiles. P50 is the median. P99 is the latency that 99 percent of requests stay below. Averages hide slow requests. Prefer percentiles.

Throughput is the count of completed operations in a unit of time. Examples: requests per second, messages per second, and rows written per second. Throughput has a limit when queues grow or when a dependency saturates.

Latency and throughput are related. If you push throughput near the limit, queues grow, and latency rises. A design that only quotes "requests per second" without a latency target is incomplete.

Write the operation that you measure. "The system is fast" is not a measure. "P99 of `POST /checkout` is less than 300 ms on the server, without the client network" is a measure. State if the number includes the database.

Capacity tests and production measures differ. A test on an empty database is not the production shape. Record data size, concurrency, and the mix of reads and writes.

Later topics cover caching, async offload, and load leveling. Those patterns exist to move latency or throughput. They do not remove the need for a measure.

### Questions

#### Theoretical questions

1. What is latency?
2. What is throughput?
3. Why do teams use P99 instead of only the average?
4. What happens to latency when throughput is near the limit?
5. Why must a performance requirement name the operation?

#### Easy practical tasks

1. Define P50, P95, and P99 in four sentences.
2. Write three performance requirements for a search box. Include latency and throughput.
3. Make a table: "Metric" and "Unit". Add latency, throughput, error rate, and queue depth.
4. List four causes of high latency in a request that reads a database.

#### Medium practical tasks

1. Sketch a latency budget for a page that calls two HTTP APIs and one database. Assign milliseconds to each step. The total must stay under 400 ms.
2. Explain in eight sentences how a long tail (high P99) can exist when P50 looks good.
3. Design a small load test plan: operations, concurrency, duration, and success criteria. Do not run an attack on a system that you do not own.

#### Advanced practical tasks

1. Collect public latency numbers from one API provider. Write how they define the measure (which hop, which percentile).
2. Write a one-page analysis of queueing: arrival rate, service time, and when latency explodes. Use a simple numeric example.

---

## Scalability

Scalability is the ability of the system to keep required latency and throughput when the load grows. Load can be more users, more data, or more events.

Vertical scale adds resources to one machine: more CPU, more memory, or a faster disk. Vertical scale is simple. It has a hard ceiling. One machine also remains a single fault domain.

Horizontal scale adds more machines or more processes. Horizontal scale needs a method to split work: load balancing, partitions, or shards. Horizontal scale also needs stateless workers or careful session design. Topic 4 covers stateful and stateless services.

A system that is not scalable shows a sharp rise in latency or a drop in success rate when load grows. A system that is scalable keeps the quality targets in the stated load range.

Scalability is not the same as a large user count today. A classroom tool for 30 users does not need a global shard plan. Design for the load that you can measure plus a stated growth factor.

Cost and complexity grow with horizontal scale. More processes mean more failure modes. Topic 11 covers distributed-system basics.

Write the growth axis. "Scale the system" is empty. "Keep P99 of reads under 200 ms when the item table grows from 1 million to 20 million rows" is a scalability requirement.

### Questions

#### Theoretical questions

1. What is scalability?
2. What is vertical scale?
3. What is horizontal scale?
4. Why is a large user count today not the same as a scalability design?
5. What extra failure modes appear when you add more processes?

#### Easy practical tasks

1. Write four sentences that compare vertical scale and horizontal scale.
2. List three growth axes (users, data, events). Give one example system for each axis.
3. Make a table: "Technique" and "What it splits". Add load balancer, replica, and shard.
4. Write one scalability requirement for a photo album with 100 users and one for 1 million users.

#### Medium practical tasks

1. For a blog with one database, write when you add a read replica and when you do not. Give three signals.
2. Draw a before-and-after diagram: one server versus three web processes and one database. Label new failure points.
3. Estimate a capacity plan for 10 times more daily orders. List which component hits the limit first and why.

#### Advanced practical tasks

1. Write a one-page note on the difference between scalability and elasticity (add and remove capacity automatically).
2. Take a public architecture of a known product. Identify the growth axis that the design serves. Name one axis that the design does not serve well.

---

## Security

Security is the protection of data and operations against unauthorized use, disclosure, change, or destruction. Architecture sets trust boundaries, identity, access rules, and the handling of secrets.

Authentication answers "who is this principal?". Authorization answers "what may this principal do?". Those two ideas are different. Topic 14 in the later path covers them in more depth.

A trust boundary is a line where the level of trust changes. The public internet is a low-trust zone. An internal admin network is a higher-trust zone. Do not treat a request from the internal network as safe without authentication.

Least privilege is a rule: give each process, user, and token only the rights that the work needs. A report job that only reads must not use an account that can drop tables.

Secrets are passwords, tokens, and keys. Do not store secrets in source control. Do not log secrets. Use a secret store that the operations team controls.

This handbook is defensive. Use security design to reduce accidents and abuse. Do not use this material to plan attacks. Follow the law and the rules of your school or employer.

A security requirement must be testable. "Encrypt data" is incomplete. "TLS 1.2 or newer on all public HTTP" and "personal data at rest uses the platform disk encryption" are closer to requirements.

### Questions

#### Theoretical questions

1. What is security as a quality attribute?
2. What is the difference between authentication and authorization?
3. What is a trust boundary?
4. What is least privilege?
5. Why must a security requirement be testable?

#### Easy practical tasks

1. Write five defensive security requirements for a notes API. Include identity, access, and secrets.
2. Make a table: "Item" and "Secret? (yes/no)". Add API keys, user display names, database passwords, and public documentation URLs.
3. Draw two trust zones for a web app: browser and server. Mark the boundary.
4. List four places where a beginner accidentally stores a password. Write the safe alternative for each place.

#### Medium practical tasks

1. Design roles for a school grade book: student, teacher, and registrar. Write what each role can read and write.
2. Write a short threat list (six items) for a public comment form. For each item, write one defensive control. Do not write attack steps.
3. Plan secret rotation for a database password in six steps. Include application restart or reload.

#### Advanced practical tasks

1. Write a one-page security quality scenario: a stolen session cookie. Describe detection, impact limit, and recovery. Stay defensive.
2. Map STRIDE names (spoofing, tampering, repudiation, information disclosure, denial of service, elevation of privilege) to six controls in a small web app. One sentence per name.

---

## Maintainability / evolvability

Maintainability is the ease of change of the system by the people who own it. Evolvability is the ability of the system to accept new requirements without a full rewrite. The two terms are close. This handbook uses maintainability for daily change and evolvability for larger change.

Measures include time to understand a module, time to add a field, time to fix a defect, and the size of a change that touches many modules. High coupling decreases maintainability. Topic 3 covers cohesion and coupling.

A system that is hard to test is hard to change. Hidden dependencies, global mutable state, and a shared database with unclear owners all decrease evolvability.

Documentation and ADRs improve maintainability when they stay short and current. An 80-page design that no one updates does not help. A module with a clear interface and three ADRs does help.

Team skill is part of this quality. A clever structure that only one person understands is not maintainable. Prefer a simple structure that the team can operate.

You often trade performance tricks for maintainability. A hand-tuned shortcut that saves 2 ms and blocks a schema change for six months is a bad trade if latency is already inside the target.

### Questions

#### Theoretical questions

1. What is maintainability?
2. What is evolvability?
3. How does coupling decrease maintainability?
4. Why does a system that is hard to test resist change?
5. How can a performance shortcut harm evolvability?

#### Easy practical tasks

1. Write four measures of maintainability for a student project.
2. List five code or data smells that make change slow.
3. Make a table: "Change" and "Modules you must touch". Add four changes for a shop.
4. Write five sentences on why a simple module boundary helps a new teammate.

#### Medium practical tasks

1. Take a small repository that you own. Time how long a new reader needs to find the write path for one feature. Write the path and the obstacles.
2. Propose three changes that would reduce the number of modules touched for a common change. Do not implement a large rewrite.
3. Write an ADR that rejects a "clever" library because the team cannot maintain it.

#### Advanced practical tasks

1. Design a one-year evolution plan for a monolith that must add a second client (mobile). List what must stay stable.
2. Compare two public codebases of similar size. Write which one is easier to change and why. Use coupling and tests as evidence.

---

## Observability

Observability is the ability to explain the current behavior of the system from its outputs. The common outputs are logs, metrics, and traces.

Logs are event records. A useful log has a time, a level, a message, and a request identifier. Metrics are numbers over time: request count, error count, latency, and queue depth. Traces follow one request across processes. A trace uses a trace identifier.

Without observability, you cannot prove a quality target in production. Availability, latency, and error rate need measures. Observability is the means to collect those measures.

Observability has a cost. High-volume debug logs increase storage cost and can hide important events. Cardinality explosions in metrics (a new time series per user identifier) can break the metrics system.

Correlation is required. A request identifier that appears in logs and traces lets you join signals. Topic 16 in the later path covers RED and USE methods and OpenTelemetry at a high level.

Design observability with the feature. Do not add logs only after the first incident. Name the questions that you must answer at 03:00: Is the dependency down? Is the queue growing? Is one tenant the cause?

### Questions

#### Theoretical questions

1. What is observability?
2. What are logs, metrics, and traces?
3. Why do quality targets need observability?
4. How can observability increase cost?
5. What is a request identifier used for?

#### Easy practical tasks

1. Write four questions that an operator must answer during an outage.
2. Make a table: "Signal" and "Example". Add one log, one metric, and one trace field.
3. List five fields that every request log must include.
4. Write a bad log line and a better log line for a failed payment (no secrets).

#### Medium practical tasks

1. Design a minimum observability set for a two-process system: API and worker. List logs, metrics, and one trace path.
2. Explain cardinality in metrics with a numeric example (one series per user versus one series per route).
3. Write an alert rule in words: condition, duration, and who receives it. Explain how you avoid a noisy alert.

#### Advanced practical tasks

1. Map one user action (place order) across logs, metrics, and a trace. Draw the identifiers that join the signals.
2. Read the OpenTelemetry overview documentation. Write ten sentences on traces and context propagation. Do not copy long passages.

---

## Cost

Cost is the money, time, and people that the system consumes. Architecture choices change cost. A managed database can increase money cost and decrease people cost. Many small services can increase money cost and people cost.

Money cost includes compute, storage, network egress, licenses, and vendor support. People cost includes development time, on-call time, and training. Time-to-delivery is also a cost. A design that is "cheaper per request" but ships six months late can be the expensive design.

Cost is a quality attribute because stakeholders must accept it. A design that meets latency and fails the budget is not acceptable.

Write cost with a unit and a period. "The system is cheap" is empty. "Compute plus database is less than X per month at Y requests per day" is a target. Include the people hours for operations.

Premature distribution increases cost. Extra networks, extra stores, and extra pipelines each add invoices and failure modes. Start small. Measure. Then pay for scale when a quality target fails.

Cost also includes waste: idle capacity, unused indexes, and debug systems that stay on. Observability and capacity reviews find waste.

### Questions

#### Theoretical questions

1. What kinds of cost does architecture include?
2. How can a managed service increase money cost and decrease people cost?
3. Why is time-to-delivery a cost?
4. Why is a design that misses the budget a failed design?
5. How does premature distribution increase cost?

#### Easy practical tasks

1. List eight cost items for a small web application on a public cloud.
2. Write a monthly cost target with units for a student project (even if the number is low).
3. Make a table: "Choice" and "Cost that rises". Add managed database, extra service, and extra region.
4. Write four sentences on people cost of on-call for a two-person team.

#### Medium practical tasks

1. Compare "one virtual machine plus one database" with "three services plus a broker plus two databases" for a 200-user tool. Write a cost argument, not a fashion argument.
2. Estimate storage cost growth if logs stay for 90 days at 5 GB per day. Show the arithmetic.
3. Write an ADR that selects a managed database because people cost is higher than money cost for the team.

#### Advanced practical tasks

1. Build a one-year cost model with three load levels. Include compute, storage, and 4 hours of operations per week.
2. Read a public cloud pricing page for one database product. Write which meter (storage, I/O, hours) dominates for a read-heavy app.

---

## You cannot maximize all of them

Quality attributes compete. A choice that improves one attribute often decreases another attribute. That fact is the core of architecture.

Examples:

- More redundancy can increase availability and increase cost.
- Stronger security checks can increase safety and increase latency.
- More modules can increase evolvability and decrease performance if each call is a network call.
- More logs can increase observability and increase cost.
- A simple monolith can increase maintainability for a small team and decrease independent scalability of one part.

A quality attribute scenario states the stimulus, the environment, and the measurable response. Use scenarios to make the trade-off explicit. "We accept P99 of 400 ms to keep one database and one deploy unit" is a decision.

Stakeholders must accept the weaker attributes. If no one accepts the weaker attribute, the project has no architecture. It has a conflict.

Do not add every pattern from a book. Each pattern spends a budget of complexity. Spend the budget on the attributes that the stakeholders named.

Record the rejected quality in the ADR. Future readers must see that you did not "forget" security or cost. You chose.

### Questions

#### Theoretical questions

1. Why can you not maximize all quality attributes?
2. Give two pairs of attributes that often compete.
3. What is a quality attribute scenario?
4. Who must accept the weaker attributes?
5. Why do you record a rejected quality in an ADR?

#### Easy practical tasks

1. Write six trade-off sentences in the form "If we improve X, Y can get worse because ...".
2. Make a table with seven attributes from this topic. Mark high, medium, or low priority for a classroom quiz app.
3. For a payments app, mark the same table again. Write three differences.
4. Write a one-paragraph stakeholder statement that accepts higher cost for higher availability.

#### Medium practical tasks

1. Write three quality attribute scenarios for one system: latency, availability, and cost. Show where they conflict.
2. Take a design that uses many services for a small team. Rewrite it to maximize maintainability and cost. List what you lose.
3. Facilitate a paper workshop: five stakeholders, five attributes, limited budget of 10 points to assign. Record the result.

#### Advanced practical tasks

1. Write an ADR that explicitly lowers one attribute to save another. Include numbers and the owner of the risk.
2. Find a public engineering blog that describes a trade-off (consistency versus latency, or cost versus isolation). Summarize the trade-off in one page. Do not copy long quotations.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you turn an empty quality word into a requirement that a test can check?
2. Which attributes need production signals, and why does that link observability to the others?
3. How do team size and cost change the availability target that you can promise?
4. Why is security a quality attribute and not only a feature list?
5. What sequence do you use: measure a quality, then add a pattern, or add a pattern, then hope?

#### Easy practical tasks

1. Write a one-page quality sheet for a campus bike-rental app. Include one measurable target for each attribute in this topic.
2. Draw a radar-style table (not a pretty chart) with scores 1 to 5 for your last project. Justify each score in one sentence.
3. List ten empty phrases and rewrite each phrase as a measurable quality statement.
4. Bookmark one public SLO or status-page vocabulary guide. Write five terms and your definitions.

#### Medium practical tasks

1. Write a quality workshop agenda (60 minutes) that ends with three accepted trade-offs and two ADRs to draft.
2. For an online exam system, pick the three attributes that must win. Write the attributes that you will allow to be weaker. Give reasons.
3. Design a monthly quality review: which graphs, which costs, which security checks. Keep the review to one hour.

#### Advanced practical tasks

1. Write a full quality attribute tree for one product: attributes, sub-attributes, measures, and current gaps.
2. Compare two reference architectures (monolith on one host versus split services). Score them on all attributes in this topic for a five-person team.
