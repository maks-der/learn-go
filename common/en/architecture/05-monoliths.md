# 5. Monoliths

## Description

A monolith is a system that you deploy as one unit. All modules of the product run in that unit (one process or a small set of processes that you always ship together). This topic covers the modular monolith, when a monolith is the default, the difference between a deploy unit and logical modules, a later extract path, and symptoms of a big ball of mud.

A monolith is not a defect. For a small team and a young product, a monolith is often the correct structure. Later topics cover other styles. Those styles have a higher operations cost.

Complete Topics 1 to 4 before this topic. Keep modules clear inside the monolith (Topic 3). Use one database until data ownership forces a split (Topic 9).

Use one term for each concept. A monolith is a deploy shape. A module is a design boundary. Those two ideas are not the same.

---

## Modular monolith

A modular monolith is a monolith with clear module boundaries. Packages (or equivalent units) have names, interfaces, and acyclic dependencies. Modules do not import internals of other modules. The product still builds and deploys as one unit.

The modular monolith gives many benefits of modules without the network. Calls stay in process. Transactions can stay local. Debugging stays in one process. Tests can boot one application.

Rules that keep a modular monolith healthy:

- Each module has one job (high cohesion).
- Modules talk through interfaces, not through shared tables with no owner.
- The dependency graph has no cycles.
- The database schema has owners (a table belongs to one module).
- New features land in an existing module or in a new module. They do not land in a `utils` bag.

A modular monolith is still one fault domain at deploy time. A defect in one module can crash the process. You accept that cost in exchange for simplicity.

Do not treat "monolith" as "one folder of random files". That shape is a big ball of mud. The last section of this topic describes the symptoms.

Start here for a new product unless a hard constraint forces more than one deployable unit on day one.

### Questions

#### Theoretical questions

1. What is a modular monolith?
2. What benefits do you keep when calls stay in process?
3. Name four rules that keep a modular monolith healthy.
4. What fault-domain cost does a monolith accept?
5. How is a modular monolith different from one folder of random files?

#### Easy practical tasks

1. Write five sentences that define a modular monolith. Use only facts from this section.
2. Draw a monolith with four modules and one database. Label allowed calls.
3. Make a table: "Rule" and "Defect if you break it". Add four rules.
4. List three features of a blog and the module that owns each feature.

#### Medium practical tasks

1. Design a modular monolith for a library loan system: modules, interfaces, and table owners.
2. Write a package lint rule list (five rules) that protects module boundaries.
3. Compare a modular monolith with a "layered only" codebase that has no module names. Write eight sentences.

#### Advanced practical tasks

1. Write a one-page standard: how a new teammate adds a module to the monolith.
2. Review a public monolith repository. Mark modules that exist and modules that are only folders. Write evidence.

---

## When a monolith is the right default

A monolith is the right default when the team is small, the domain is young, and the operations budget is small. One deploy pipeline, one runtime, and one database match that budget.

Use a monolith when:

- One team owns the full product.
- You cannot yet draw stable module boundaries that will last a year.
- You need one transaction for a user action.
- You must ship in weeks, not in a year of platform work.
- You do not have on-call capacity for many services.

Do not split for fashion. "We might need to scale one part" is not a reason if that part has no measure that fails. Topic 2 requires measures.

A monolith can still scale. You can run more instances of the same unit behind a load balancer if the app is stateless at the instance (Topic 4). You can add a read replica. You can add a cache. Those steps stay cheaper than a service split.

A monolith is the wrong default when law or a vendor forces a separate system, when two products have no shared domain, or when two teams must release on independent clocks and a modular monolith already fails that need.

Write the reason for a split in an ADR. If you cannot write the reason with a measure or a constraint, keep the monolith.

### Questions

#### Theoretical questions

1. When is a monolith the right default?
2. Why is a possible future scale need not enough for a split?
3. How can a monolith scale without a service split?
4. When is a monolith the wrong default?
5. What must an ADR for a split contain?

#### Easy practical tasks

1. Write a checklist of eight questions. A "yes" majority means "keep the monolith".
2. Make a table: "Situation" and "Monolith default? (yes/no)". Add six rows.
3. List five operations tasks that exist once and stay once in a monolith.
4. Write four sentences to a teammate who wants "microservices" for a student shop.

#### Medium practical tasks

1. Given a five-person team and a 2000-user app, write a one-page recommendation: monolith plus which scale steps.
2. Write an ADR that rejects a three-service split for a classroom project.
3. Find a public story where a team moved back to a monolith. Write five facts from that story. Do not copy long text.

#### Advanced practical tasks

1. Write a decision tree from constraints (team, time, law, existing systems) to "monolith" or "not yet a split".
2. Compare cost for one year: one monolith with three instances versus five services. Use people hours as the main cost.

---

## Deployment unit vs logical modules

A deployment unit is what you build, ship, and restart together. A logical module is a design boundary inside that unit (or across units).

These two ideas must stay distinct. You can have many logical modules and one deployment unit. That is the modular monolith. You can also have one logical module that you accidentally deploy as many units (copies of the same confusion).

A change of the deployment unit is costly. You need a pipeline, a runtime, health checks, and a failure plan for each unit. A change of a logical module is cheaper: a package, an interface, a test.

Do not add a deployment unit to "make the module real". The module is already real if the boundary is real in code and in data ownership.

Maps that help:

- C4 container view shows deployment units and data stores.
- A package diagram shows logical modules.

A container is not a module. A container is a packaging and isolation tool. You can put a modular monolith in one container.

When you later extract a module, you change its deployment unit. The logical module must already exist. Extraction then moves a boundary that you already understand. The next section covers that path.

### Questions

#### Theoretical questions

1. What is a deployment unit?
2. What is a logical module?
3. Why is a new deployment unit costly?
4. Why is a container not a module?
5. Which diagram shows units, and which diagram shows modules?

#### Easy practical tasks

1. Write five sentences that separate deploy unit and logical module.
2. Draw one deploy unit that contains five logical modules.
3. Make a table: "Item" and "Unit or module". Add container image, package, process, and interface.
4. List four operational items that each new deploy unit needs.

#### Medium practical tasks

1. Sketch C4 container and a package diagram for the same library system. Write what each view hides.
2. A teammate puts each package in its own container "for architecture". Write a review that rejects the plan.
3. Count deploy units and modules in a system that you use. Guess. Then write what you would need to confirm the count.

#### Advanced practical tasks

1. Write a one-page glossary for your team: module, package, process, container, service, deploy unit.
2. Design a migration that changes modules but not deploy units for two quarters. List the benefits.

---

## Modular monolith extraction path

An extraction path is a sequence that moves one module from the monolith to a separate deployment unit. You use this path only after a real pain: independent scale, an independent release clock, or a hard isolation constraint.

Typical sequence:

1. Make the module real inside the monolith. Clear interface. Clear table owners. No cycles.
2. Stop other modules from reading those tables. Use the interface only.
3. Add observability on the module boundary (Topic 2).
4. Extract a library or keep the package, but treat the interface as a future network contract.
5. Run the module logic in-process behind that contract until the contract is stable.
6. Move the data if the new unit must own its store (Topic 9). Use expand-contract steps. Do not flip a switch on a shared table.
7. Introduce a network call. Add timeouts, retries, and idempotency (Topics 4 and 11).
8. Deploy the new unit. Keep a compatibility window. Keep a rollback.

Extract one module. Do not extract five modules in one project. Learn the operations cost.

If the pain disappears after step 1 to 3, stop. You already have a better monolith. That result is a success.

Write an ADR at the start and an ADR at the cut-over. Include the measure that defines success (latency, error rate, or deploy frequency).

### Questions

#### Theoretical questions

1. When do you start an extraction?
2. Why must the module be real before the network appears?
3. Why do you add observability before the split?
4. What building blocks appear when the call becomes a network call?
5. Why can a stop after better boundaries be a success?

#### Easy practical tasks

1. Number the eight steps in your own words (one sentence each).
2. Make a table: "Step" and "Still one deploy unit? (yes/no)".
3. List three measures that could justify an extract.
4. Write four risks of a first network hop.

#### Medium practical tasks

1. Plan an extract of a "search" module from a shop monolith. Write data ownership and the first network API.
2. Write a rollback plan for step 7 and step 8.
3. Identify a module in a project that is not ready to extract. List the missing rules.

#### Advanced practical tasks

1. Write a two-page extraction playbook for one module: schema, dual write or expand-contract, contract tests, and SLO.
2. Compare "extract a read-only report" with "extract the checkout write path". Explain which extract is safer and why.

---

## Big-ball-of-mud symptoms

A big ball of mud is a system with no visible architecture. Code and data grow by local patches. Boundaries do not hold. The term is a known name for this failure. This handbook uses it as a technical name.

Symptoms include:

- Every change touches many folders.
- There is no owner for a table or a type.
- Cycles exist everywhere.
- `utils`, `common`, and `helpers` hold business rules.
- You cannot test one job without a full database and a full boot.
- Names do not match the domain.
- Diagrams and code tell different stories.
- People fear a change of one field.

A monolith can be clean. A set of services can also be a big ball of mud (a distributed mud). The deploy shape does not save you. Modules and ownership save you.

Recovery starts with boundaries, not with a rewrite. Pick one job. Give it an interface. Move rules into that module. Stop new writes to leaked tables. Write tests around the invariant. Repeat.

A full rewrite is rarely the first step. A rewrite without new rules produces a new mud. Topic 1: record decisions. Topic 3: cohesion and coupling.

Treat fear of change as a quality signal (maintainability). Measure time-to-change for a standard request. Use that number in reviews.

### Questions

#### Theoretical questions

1. What is a big ball of mud?
2. Name five symptoms.
3. Why can a set of services still be a mud?
4. Why is a rewrite without new rules a risk?
5. What is the first recovery step?

#### Easy practical tasks

1. Make a checklist of ten symptoms. Score a project that you know (0 or 1 per item).
2. Write five sentences that explain why a monolith is not automatically a mud.
3. List four "bag" package names that often hide a mud.
4. Write a fear-of-change story in six sentences from a project that you know.

#### Medium practical tasks

1. Plan a four-week recovery for one leaked table: owner, interface, tests, and a stop rule for new leaks.
2. Draw the current dependency spaghetti of a small project (honest). Then draw a target acyclic graph.
3. Write a review comment template that rejects a change that adds a new cycle or a new `utils` rule.

#### Advanced practical tasks

1. Write a one-page "mud index" with weights. Apply it to two codebases. Interpret the scores.
2. Design a six-month recovery roadmap that does not stop feature work. Include capacity (for example 20 percent of time).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do Topics 3 and 5 work together: modules inside one deploy unit?
2. Why does this path teach monoliths before microservices?
3. What quality attributes (Topic 2) usually improve when you keep a monolith on a small team?
4. What quality attributes can get worse in a large, dirty monolith?
5. How do you know that a pain is "extract a module" and not "clean the module"?

#### Easy practical tasks

1. Write a one-page cheat sheet: definitions, default rules, extract steps, mud symptoms.
2. For a campus event app, name six logical modules and argue for one deploy unit.
3. Draw C4 context and container for that app with one container for the app and one for the database.
4. Write three ADR titles that a monolith project needs in month one.

#### Medium practical tasks

1. Write a two-page architecture brief for a modular monolith: modules, data owners, deploy, and the extract trigger.
2. Role-play a review: one person wants eight services. Write both sides and the decision.
3. Take a homework project. List mud symptoms and a two-week cleanup that does not add services.

#### Advanced practical tasks

1. Write a playbook that a team follows before any extract: measures, module fitness, and a go/no-go checklist.
2. Compare a public modular-monolith talk or article with this topic. Write what you accept and what you reject for a three-person team.
