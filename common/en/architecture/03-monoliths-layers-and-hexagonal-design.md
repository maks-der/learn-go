# 3. Monoliths, Layers, and Hexagonal Design

## Description

A monolith is a system that you deploy as one unit. A style is a set of rules for dependency direction and placement of code. This topic covers the modular monolith as the default, the difference between a deploy unit and logical modules, presentation / domain / infrastructure layers, ports and adapters, and the risk of over-engineering.

A monolith is not a defect. For a small team and a young product, a monolith is often the correct structure. These styles live inside one deploy unit first. They are not a reason to add services.

Complete Topics 1 and 2 before this topic. Keep modules clear inside the monolith. Use one database until data ownership forces a split (Topic 6).

Use one term for each concept. A monolith is a deploy shape. A module is a design boundary. A layer is a band of responsibility. A port is an interface that the domain defines. An adapter is an implementation that talks to the outside world.

---

## Modular monolith as the default

A modular monolith is a monolith with clear module boundaries. Packages (or equivalent units) have names, interfaces, and acyclic dependencies. Modules do not import internals of other modules. The product still builds and deploys as one unit.

The modular monolith gives many benefits of modules without the network. Calls stay in process. Transactions can stay local. Debugging stays in one process. Tests can boot one application.

Rules that keep a modular monolith healthy:

- Each module has one job (high cohesion).
- Modules talk through interfaces, not through shared tables with no owner.
- The dependency graph has no cycles.
- The database schema has owners (a table belongs to one module).
- New features land in an existing module or in a new module. They do not land in a `utils` bag.

A modular monolith is still one fault domain at deploy time. A defect in one module can crash the process. You accept that cost in exchange for simplicity.

A monolith is the right default when the team is small, the domain is young, and the operations budget is small. One deploy pipeline, one runtime, and one database match that budget.

Use a monolith when:

- One team owns the full product.
- You cannot yet draw stable module boundaries that will last a year.
- You need one transaction for a user action.
- You must ship in weeks, not in a year of platform work.
- You do not have on-call capacity for many services.

Do not split for fashion. "We might need to scale one part" is not a reason if that part has no measure that fails. Topic 1 requires measures.

A monolith can still scale. You can run more instances of the same unit behind a load balancer if the app is stateless at the instance (Topic 2). You can add a read replica. You can add a cache. Those steps stay cheaper than a service split.

Write the reason for a split in an ADR. If you cannot write the reason with a measure or a constraint, keep the monolith.

Do not treat "monolith" as "one folder of random files". That shape is a big ball of mud. Symptoms include no owner for a table, cycles, and a change that always touches the whole tree.

### Questions

#### Theoretical questions

1. What is a modular monolith?
2. What benefits do you keep when calls stay in process?
3. When is a monolith the right default?
4. What fault-domain cost does a monolith accept?
5. Why is a possible future scale need not enough for a split?

#### Easy practical tasks

1. Write five sentences that define a modular monolith. Use only facts from this section.
2. Draw a monolith with four modules and one database. Label allowed calls.
3. Make a table: "Rule" and "Defect if you break it". Add four rules.
4. Write a checklist of eight questions. A "yes" majority means "keep the monolith".

#### Medium practical tasks

1. Design a modular monolith for a library loan system: modules, interfaces, and table owners.
2. Write an ADR that rejects a three-service split for a classroom project.
3. Compare a modular monolith with a "layered only" codebase that has no module names. Write eight sentences.

#### Advanced practical tasks

1. Write a one-page standard: how a new teammate adds a module to the monolith.
2. Write a decision tree from constraints (team, time, law, existing systems) to "monolith" or "not yet a split".

---

## Deployment unit vs logical modules

A deployment unit is what you build, ship, and restart together. A logical module is a design boundary inside that unit (or across units).

These two ideas must stay distinct. You can have many logical modules and one deployment unit. That is the modular monolith. You can also have one logical module that you accidentally deploy as many units (copies of the same confusion).

A process, a container, or a virtual machine packages the deployment unit. Those tools do not create a module. If the code has no boundary, the container only ships the mud.

A later extract path exists. You can move a logical module to a new deployment unit when a real pain appears: independent scale, an independent release clock, or a hard isolation rule. Topic 9 covers that split. The extract is cheaper when the module already has an interface and owned tables.

Do not create twenty deployment units to "look like modules". Each extra unit adds a network, a pipeline, and a failure mode (Topic 8).

Write the map: logical module name, package path, table prefix or schema, and deployment unit. In a student project the last column is often "the one app" for every row.

A copy of the same binary behind a load balancer is still one logical system and one deployment shape. It is not many services. Topic 12 covers that scale step.

### Questions

#### Theoretical questions

1. What is a deployment unit?
2. What is a logical module?
3. Why can many logical modules live in one deployment unit?
4. When is a later extract cheaper?
5. Why is a container not a module by itself?

#### Easy practical tasks

1. Write five sentences that separate deploy unit from logical module.
2. Make a table: "Module" and "Deploy unit". Add six rows for a blog that ships as one app.
3. List four costs of a second deployment unit for a two-person team.
4. Draw one box for deploy and four boxes inside it for modules.

#### Medium practical tasks

1. Take a modular monolith with four modules. Write which module could become a second deploy unit and which data it must take.
2. Write an ADR title and a one-paragraph context for "one deploy unit, four logical modules".
3. Describe a defect that appears when two teams ship two processes that still share all tables.

#### Advanced practical tasks

1. Write a one-page extract checklist: interface, owned tables, contract, and on-call. Fill it for one module.
2. Compare "one binary, many instances" with "many binaries, many stores" for a team of four. Write operations cost.

---

## Presentation / domain / infrastructure

A layered style splits the program into bands. A common three-band split is presentation, domain, and infrastructure.

The presentation layer accepts input from users or other programs. It includes HTTP handlers, command-line parsing, and user interface adapters. It translates input into domain calls. It translates domain results into HTTP JSON or other output.

The domain layer holds business rules and domain types. It decides what an order may do. It does not open a SQL connection. It does not format HTML.

The infrastructure layer talks to technical systems: databases, file systems, mail, and clocks. It implements persistence and other details.

Allowed dependency direction in this simple form:

- Presentation can depend on domain.
- Infrastructure can depend on domain.
- Domain must not depend on presentation or on infrastructure.

Some teams add an application layer (use cases) between presentation and domain. That extra band is useful when handlers would contain long scripts. For a small program, three bands are enough.

Layers are not folders with no rules. If every layer imports every other layer, you have names without a style.

A layer is not a module. A module is a job (`pricing`). A layer is a technical band (`http`, `domain`, `sql`). Many teams use both: modules that each contain a thin slice of layers, or layers that contain modules. Pick one map and keep it.

### Questions

#### Theoretical questions

1. What does the presentation layer do?
2. What does the domain layer do?
3. What does the infrastructure layer do?
4. Which layer must not import SQL?
5. How is a layer different from a module?

#### Easy practical tasks

1. Write five sentences that define the three layers. Use only facts from this section.
2. Place ten types or functions of a shop into the three layers (handler, price rule, SQL query, JSON field).
3. Draw allowed arrows between presentation, domain, and infrastructure.
4. List three items that must not live in the domain layer.

#### Medium practical tasks

1. Sketch a "place order" flow across the three layers. Name the types at each boundary.
2. Take a handler that contains SQL. Split it on paper into the three layers.
3. Decide if you need an application layer for a student blog. Write an ADR paragraph.

#### Advanced practical tasks

1. Compare "layers inside modules" with "modules inside layers". Draw both. Write when each map fits a team of four.
2. Write a one-page layer standard: folder names, import rules, and two forbidden examples.

---

## Ports and adapters; the domain does not import SQL

Hexagonal architecture (ports and adapters) places the domain at the center. The outside world talks to the domain through ports. Adapters implement those ports.

A primary port is an operation that an outside actor starts (a use case). A primary adapter is the HTTP handler or the CLI that calls that port.

A secondary port is an interface that the domain needs (store an order, send a mail, read the time). A secondary adapter is the SQL store, the SMTP client, or the system clock.

The domain depends on ports (interfaces). The domain does not depend on adapters. Adapters depend on the domain types or on the port types.

The hexagon picture is a teaching shape. You do not need six sides. You need a center and a rule: details point inward.

The domain must not import SQL. SQL is a detail. A vendor driver type in a domain package couples the rule to one store. You cannot test the rule without that store. You cannot change the store without a rewrite of the rule.

This style makes tests easier. You replace a SQL adapter with a memory adapter. You replace the clock with a fixed clock. Topic 2 described interface versus implementation. Hexagonal design applies that idea to the whole application.

Do not create a port for every function. Create a port when a detail must stay replaceable or when the domain must stay free of a vendor type.

Hexagonal design and layers can live together. Presentation adapters sit on one side. Infrastructure adapters sit on the other side. The domain sits in the center.

Clean architecture is a family of rules with the same dependency rule: source code dependencies point inward. Inner circles hold policy. Outer circles hold details. Keep it small. A large folder tree with no rules is worse than a short layered design.

Map transport types to domain types at the adapter. An HTTP JSON struct is not a domain entity. A SQL row type is not a domain entity. Copy fields at the edge.

### Questions

#### Theoretical questions

1. What is a port?
2. What is an adapter?
3. What is the difference between a primary port and a secondary port?
4. Why must the domain not import SQL?
5. Why does this style help tests?

#### Easy practical tasks

1. Draw a hexagon (or a box) with domain in the center, HTTP on the left, SQL on the right.
2. Write two primary ports and two secondary ports for a library system.
3. Make a table: "Port" and "Adapter examples". Add store, mail, and clock.
4. List three vendor types that must not appear in the domain.

#### Medium practical tasks

1. Design ports for "borrow a book": use case port, book store port, clock port. Write method names.
2. Write a test plan that uses a memory adapter for the store.
3. Convert a layered sketch into ports and adapters. Show what you rename.

#### Advanced practical tasks

1. Write a one-page comparison of three-layer design and hexagonal design. State when the extra names pay off.
2. Implement a paper walkthrough of two adapters for one port (SQL and memory). Show identical domain tests.

---

## Over-engineering warning

Over-engineering is a structure that is larger than the problem. Extra layers, extra ports, extra services, and extra frameworks add cost. They do not add quality unless a measure or a constraint requires them.

Signs of over-engineering:

- A folder tree with more names than types.
- A port for every function, including `now()` wrappers that never change.
- A microservice plan for a one-person course project.
- A "clean architecture" template that the team cannot explain.
- Mapping code that is longer than the rule that it protects.

Signs of under-engineering:

- SQL in the HTTP handler.
- One folder of random files.
- No timeout on outbound calls.
- No owner for a table.

Aim for the smallest style that protects the domain and the quality attributes that you named in Topic 1.

A three-layer student project can be enough. Add hexagon names when you need replaceable adapters or when the handler scripts grow. Add a second deploy unit only after a real boundary pain (Topic 13).

Write an ADR when you add a style. The context must name the problem. "We saw a blog post" is not a problem.

If a new teammate cannot point to the domain package in five minutes, the style failed. Simplify.

### Questions

#### Theoretical questions

1. What is over-engineering in this handbook?
2. Name three signs of over-engineering.
3. Name three signs of under-engineering.
4. When do extra hexagon names pay off?
5. Why must an ADR name a problem before you add a style?

#### Easy practical tasks

1. Write five sentences about the over-engineering warning. Use only facts from this section.
2. Make a table: "Choice" and "Over, under, or fit". Add six rows for a student shop.
3. List four folders you would delete from an empty "enterprise" template.
4. Write four sentences to a teammate who wants six layers for a notes app.

#### Medium practical tasks

1. Review a public "clean architecture" sample. Mark names that protect the domain and names that only add hops.
2. Write an ADR that keeps three layers and rejects a service split for a classroom project.
3. Compare two designs of the same blog: flat handlers with SQL versus three layers. Write eight sentences.

#### Advanced practical tasks

1. Write a one-page "minimum style" standard for a four-person team: folders, import rules, and a reject list.
2. Take a vendor reference architecture. Mark every part that your constraints would remove. Write a one-page residual design.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a modular monolith, layers, and ports work together in one codebase?
2. Why must you keep deploy unit and logical module as two different ideas?
3. What quality attributes from Topic 1 get weaker if you split the monolith too early?
4. How does "the domain does not import SQL" protect tests and change cost?
5. When is a big ball of mud worse than a small over-engineered tree?

#### Easy practical tasks

1. Write a one-page cheat sheet: modular monolith rules, deploy versus module, three layers, port/adapter, over-engineering signs.
2. For a to-do app, name four modules, three layers, one port, and one adapter.
3. Draw one deploy box, three layers, and arrows that do not point at SQL from the domain.
4. Bookmark [https://c4model.com/](https://c4model.com/). Sketch a container view that is still one deployable unit.

#### Medium practical tasks

1. Write a short architecture brief for a campus lost-and-found monolith: modules, layers, table owners, and one rejected service split.
2. Take a handler that mixes HTTP, rules, and SQL. Rewrite it on paper as presentation, domain, and a store port.
3. Write a twelve-week plan: weeks for a modular monolith, and the measure that would reopen a split.

#### Advanced practical tasks

1. Write a mapping guide: HTTP JSON type to domain type to SQL row. Include two forbidden shortcuts.
2. Review a public monolith. Score it 0 to 5 on module boundaries, layer direction, and style size. Explain each score.
