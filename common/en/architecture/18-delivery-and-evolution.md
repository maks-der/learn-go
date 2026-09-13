# 18. Delivery and Evolution

## Description

Delivery and evolution are the ways you ship change without a long freeze. This topic covers continuous integration and continuous delivery (CI/CD), rolling, blue/green, and canary releases, feature flags, expand-contract migrations, compatibility windows, and Architecture Decision Records over time.

A structure that you cannot change is a dead structure (Topic 2, maintainability). Topic 1 introduced ADRs. Topic 8 introduced compatibility. Topic 5 mentioned expand-contract on extract. This topic makes evolution a daily practice. Complete Topics 1 to 17 before this topic.

Use one term for each concept. CI is the automatic build and test of every change. CD is the automatic path to an environment. A rolling update replaces instances in small groups. Blue/green keeps two environments. A canary is a small slice of traffic. A feature flag is a runtime switch. Expand-contract is a multi-step schema change. A compatibility window is a time when old and new clients both work.

---

## CI/CD

Continuous integration (CI) is the habit that every change builds and is tested in a shared pipeline. Developers integrate often. The pipeline compiles, runs tests, and reports fail fast.

Continuous delivery (CD) is the habit that the same pipeline can deploy a proven build to an environment. Some teams stop at a button. Some teams deploy automatically to production after checks. Continuous deployment is the automatic production step. This handbook uses CD for the path. Write if the last step is automatic.

A pipeline is architecture. It is the only supported way to production. A manual copy from a laptop is a supply-chain risk (Topic 14).

Minimum pipeline:

- Fetch the exact revision
- Build
- Unit and contract tests
- Security scan that the team can run
- Store the artifact
- Deploy to a named environment
- Smoke check
- Record who shipped what

Secrets for the pipeline live in the pipeline store, not in the repository (Topic 14).

CI does not replace design review. A green pipeline can still ship a bad boundary. ADRs and code review stay.

If the pipeline is slow, people skip it. Keep the default check set fast. Put long tests in a later stage.

Do not disable a failing test to "unblock" without a ticket and a time limit.

### Questions

#### Theoretical questions

1. What is continuous integration?
2. What is continuous delivery in this handbook?
3. Why is a laptop copy to production a risk?
4. What steps belong in a minimum pipeline?
5. Why must a slow pipeline be treated as a defect?

#### Easy practical tasks

1. Write five sentences that define CI and CD.
2. Make a table: "Step" and "Purpose". Add build, test, store artifact, and smoke check.
3. List four items that must not live in the repository (secrets and similar).
4. Write a smoke check list of five HTTP calls for a catalog API.

#### Medium practical tasks

1. Draw a pipeline from commit to production. Mark a manual approval if you use one.
2. Write a one-page pipeline policy: required checks, who can skip, and how you track a skip.
3. Explain in eight sentences how contract tests protect Topic 8 compatibility during CD.

#### Advanced practical tasks

1. Write an ADR: automatic deploy to production versus a button. Include error-budget policy (Topic 15).
2. Design a pipeline for two services and one shared library. Write how you version the library.

---

## Rolling, blue/green, canary

These are three ways to replace a running version.

A rolling update replaces a subset of instances, then the next subset. Capacity drops during the roll unless you add extra instances first. Old and new code run at the same time. The contract must tolerate that mix (later section).

Blue/green keeps two full environments. Blue serves traffic. You deploy green. You test green. You switch traffic. Rollback is a switch back. Cost is about two environments. Data migrations still need care. A switch does not undo a destructive schema change.

A canary sends a small share of traffic to the new version. You watch SLIs (Topics 15 and 16). If the canary burns the budget, you stop. Then you continue or you roll back.

Pick rolling for simple stateless services with extra capacity. Pick blue/green when you want a fast switch and you can pay for two stacks. Pick canary when the risk is high and you can split traffic.

All three need:

- Health checks
- Drain of in-flight requests
- Compatible data
- A rollback that you practiced

Do not call a Friday night replace of one VM a strategy. Write the name and the abort rule.

Database primaries do not roll like stateless pods. Schema change uses expand-contract (next sections).

### Questions

#### Theoretical questions

1. What is a rolling update?
2. What is blue/green?
3. What is a canary?
4. Why must old and new code tolerate a mix during a roll?
5. Why does a traffic switch not undo a destructive schema change?

#### Easy practical tasks

1. Write five sentences that compare the three styles.
2. Make a table: "Style" and "Extra cost". Add rolling, blue/green, and canary.
3. List six abort signals for a canary (error rate, latency, and four more that you name).
4. Draw blue and green behind a balancer. Label the switch.

#### Medium practical tasks

1. Write a roll plan for six API instances: batch size, drain, and smoke.
2. Explain in eight sentences when a canary lies because the small slice misses a rare tenant.
3. Write a rollback runbook for blue/green that includes "schema already expanded".

#### Advanced practical tasks

1. Write a one-page release-style guide: which service uses which style and why.
2. Compare canary on a gateway weight with canary on a flag (next section). Write operations cost.

---

## Feature flags

A feature flag is a runtime key that selects behavior without a new deploy. The code contains both paths. The flag chooses the path.

Flags help:

- Dark launch (code in production, behavior off)
- Canary of a behavior independent of instance version
- Kill switch for a degrade path (Topic 15)
- Release that is not bound to a marketing date

Flags cost. Each flag is a fork in the code. Old flags become a big ball of mud (Topic 5). Give every flag an owner and a remove date.

Types:

- Release flags: short life, removed after the launch
- Ops flags: degrade and kill switches
- Experiment flags: time-boxed tests

Store flags in a service or a config that you can change without a rebuild. Cache the value. If the flag store is down, fail to a safe default. Safe for a new checkout path is often "off". Safe for a security control is "on" or fail closed (Topic 14).

Do not hide a schema-incompatible change behind a flag if both paths share one destructive column change. Flags do not replace expand-contract.

Log flag evaluations at debug volume only. An evaluation per request at info level fills disks (Topic 16).

Write who may flip production flags. A random flip is a deploy.

### Questions

#### Theoretical questions

1. What is a feature flag?
2. What jobs do flags help?
3. Why must every flag have a remove date?
4. What default do you use if the flag store is down for a new checkout path?
5. Why can a flag not replace expand-contract?

#### Easy practical tasks

1. Write four sentences that define feature flags.
2. Make a table: "Flag type" and "Life". Add release, ops, and experiment.
3. List five flags that you would refuse to keep for a year.
4. Write a safe default for a new recommendation widget and for an authorization check.

#### Medium practical tasks

1. Design a flag for a new search backend. Write evaluation point, default, metrics, and remove date.
2. Explain in eight sentences how flag debt becomes a mud ball.
3. Write a one-page flag policy: owners, audit of flips, and a monthly cleanup.

#### Advanced practical tasks

1. Write an ADR: flags versus extra branches that never merge. Include review cost.
2. Compare a local config flag with a hosted flag service. Write failure mode and cost (Topic 19).

---

## Expand-contract migrations

Expand-contract (also called parallel change) is a sequence that changes a schema or a contract in compatible steps. You expand. You migrate. You contract. You do not flip a breaking change in one deploy.

Example: rename a column.

1. Expand: add the new column. Deploy code that writes both columns and reads the old column.
2. Backfill the new column.
3. Deploy code that reads the new column and still writes both.
4. Contract: stop writes to the old column. Deploy.
5. Remove the old column after a compatibility window.

The same idea applies to APIs: add a field, wait, then remove the old field (Topic 8). The same idea applies to events: add a field, never reuse a name.

Expand-contract needs more deploys. That is the price of online change. A maintenance window with downtime can be honest for a student database. Write the choice.

Do not wrap a destructive `ALTER` in the same release as the only code that understands the new shape if you cannot roll back the `ALTER`.

Test the backfill on a copy of production size. A backfill that locks a hot table is an incident.

Pair with backups (Topic 15). A bad contract step needs a restore plan.

### Questions

#### Theoretical questions

1. What is expand-contract?
2. Why do you add the new column before you remove the old column?
3. What is a backfill?
4. How does the idea apply to API fields?
5. Why is a lock during backfill an incident risk?

#### Easy practical tasks

1. Number five steps for a column rename in your own words.
2. Make a table: "Step" and "Old clients still work? (yes/no)".
3. List four changes that need expand-contract (rename, type change, split table, and one more).
4. Write four sentences on when a downtime window is an honest choice.

#### Medium practical tasks

1. Plan expand-contract for `users.name` split into `given_name` and `family_name`.
2. Explain in eight sentences how a rolling deploy plus a breaking schema change fails.
3. Write a backfill checklist: batch size, locks, metrics, and abort.

#### Advanced practical tasks

1. Write a two-page migration playbook for a hot orders table. Include SLO burn and rollback.
2. Compare expand-contract with a dual-write to a new store (Topic 9). Write when each applies.

---

## Compatibility windows

A compatibility window is a period when more than one version of a contract must work. Clients, messages, and files do not update at the same second.

You need a window when:

- Mobile clients update late (Topic 20)
- Two services roll at different times (Topic 12)
- A partner file format changes (Topic 13)
- Old messages remain in a queue (Topic 10)

Rules:

- Add, do not break, during the window
- Document the start and the end date
- Measure how many clients still use the old shape
- End the window with a contract step, not with hope

Semantic versioning of APIs is a communication tool. The window is the operational fact. A `v2` path can run beside `v1` until `v1` traffic is gone.

Idempotency keys and unknown-field policies belong in the window. Servers must ignore unknown fields that they do not own. Clients must tolerate new fields.

Do not run forever dual support without a date. Forever dual is a hidden product.

Law can force a long window (old export format). Record the law (Topic 1).

The gateway can route by version header. The owner service still must understand the payloads that it accepts.

### Questions

#### Theoretical questions

1. What is a compatibility window?
2. When do you need a window?
3. What does "add, do not break" mean?
4. Why must the window have an end date?
5. Why is forever dual support a hidden product?

#### Easy practical tasks

1. Write five sentences that define a compatibility window.
2. Make a table: "Client type" and "Typical lag". Add browser, mobile, and partner file.
3. List four metrics that tell you that `v1` is still in use.
4. Write a sunset notice in four sentences for a partner CSV version.

#### Medium practical tasks

1. Design `v1` and `v2` of `GET /loans`. Write the window and the reject date for `v1`.
2. Explain in eight sentences how a queue with a one-week retention extends the window.
3. Write a one-page version policy: headers, deprecation, and measurement.

#### Advanced practical tasks

1. Write an ADR: one API with additive fields versus a `/v2` path. Include mobile lag.
2. Plan a window for an event field rename using expand-contract plus consumer order.

---

## ADRs over time

An Architecture Decision Record is a short record of one costly decision (Topic 1). See [https://adr.github.io/](https://adr.github.io/). This section covers the life of ADRs after the first week.

Status changes:

- Proposed
- Accepted
- Deprecated
- Superseded

When a new ADR replaces an old ADR, mark the old one superseded. Point to the new number. Do not delete the old file. History explains strange code.

Review ADRs on a schedule. A quarterly read of accepted ADRs finds decisions that the code already reversed. Either change the code or change the ADR.

Link ADRs to quality attributes and to SLOs. "We use a queue" must name the lag that you accept (Topics 10 and 15).

An ADR index (table) is part of the architecture. New people read the index first.

Do not write an ADR for every rename. Write an ADR when rollback is expensive.

Flags, windows, and pipelines are decisions. They deserve ADRs when they become the team law.

A superseded ADR about a monolith does not vanish after an extract (Topic 12). The extract ADR cites it.

Keep the language factual. A future reader must see the constraint that is now gone (team size, law, vendor).

### Questions

#### Theoretical questions

1. What statuses does this section use for an ADR?
2. Why do you keep a superseded ADR?
3. When do you review ADRs?
4. When do you not write an ADR?
5. How does an extract ADR relate to an older monolith ADR?

#### Easy practical tasks

1. Open [https://adr.github.io/](https://adr.github.io/). Write four sentences on purpose.
2. Make a table: "Status" and "Meaning". Add the four statuses.
3. List five decisions on a student project that need an ADR after month three.
4. Copy an index table: number, title, status, date. Fill five fictional rows.

#### Medium practical tasks

1. Write a superseded pair: ADR 4 accepts one database. ADR 12 extracts search. Cross-link them.
2. Review a public repository that contains ADRs. Summarize one life cycle in six sentences. Do not copy the full text.
3. Write a quarterly ADR-review agenda (one page).

#### Advanced practical tasks

1. Write three ADRs that form a chain: monolith, add queue, add compatibility window. Show context reuse.
2. Design an ADR quality bar: required sections, link to SLO, and a reject list (meeting notes, slogans).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do CI/CD, a release style, and a compatibility window form one safe ship path?
2. How do feature flags and expand-contract solve different evolution problems?
3. Why must schema change and instance replace use different techniques?
4. How do ADRs over time protect evolvability (Topic 2) when people leave?
5. Why does this path place delivery after scalability and reliability?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic.
2. For a to-do API, write a four-stage pipeline, a rolling rule, one flag, and one expand-contract title.
3. Draw a timeline of a column rename across four deploys.
4. Bookmark [https://adr.github.io/](https://adr.github.io/). Write one sentence on when you open it versus when you open the pipeline.

#### Medium practical tasks

1. Write a two-page evolution brief for a campus shop: pipeline, canary, flags, schema steps, windows, and ADR index.
2. Take a teammate plan that is "stop the site on Sunday and migrate". Rewrite it as expand-contract plus a smaller window.
3. Map Topics 8, 12, and 15 onto one release of a billing extract. Write the abort rule.

#### Advanced practical tasks

1. Design a twelve-week delivery standard for a four-person team. Include budget freeze and flag cleanup.
2. Compare this topic with a public twelve-factor or SRE release chapter outline. Map five headings. Do not copy long passages.
