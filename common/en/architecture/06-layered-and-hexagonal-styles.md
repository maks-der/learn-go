# 6. Layered and Hexagonal Styles

## Description

A style is a set of rules for dependency direction and placement of code. This topic covers layered design, hexagonal design (ports and adapters), a small form of clean architecture, the rule that the domain does not import SQL, mapping between transport types and domain types, and the risk of over-engineering.

These styles live inside one deploy unit first. They are not a reason to add services. Use them to protect the domain from frameworks and stores.

Complete Topics 1 to 5 before this topic. Keep the style small. A three-layer student project can be enough. A large folder tree with no rules is worse than a short layered design.

Use one term for each concept. A layer is a band of responsibility. A port is an interface that the domain defines. An adapter is an implementation that talks to the outside world.

---

## Presentation / domain / infrastructure layers

A layered style splits the program into bands. A common three-band split is presentation, domain, and infrastructure.

The presentation layer accepts input from users or other programs. It includes HTTP handlers, command-line parsing, and user interface adapters. It translates input into domain calls. It translates domain results into HTTP JSON or other output.

The domain layer holds business rules and domain types. It decides what an order may do. It does not open a SQL connection. It does not format HTML.

The infrastructure layer talks to technical systems: databases, file systems, mail, and clocks. It implements persistence and other details.

Allowed dependency direction in this simple form:

- Presentation can depend on domain.
- Infrastructure can depend on domain (or presentation can depend on infrastructure only through interfaces — see the next section).
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

## Hexagonal / ports and adapters

Hexagonal architecture (ports and adapters) places the domain at the center. The outside world talks to the domain through ports. Adapters implement those ports.

A primary port is an operation that an outside actor starts (a use case). A primary adapter is the HTTP handler or the CLI that calls that port.

A secondary port is an interface that the domain needs (store an order, send a mail, read the time). A secondary adapter is the SQL store, the SMTP client, or the system clock.

The domain depends on ports (interfaces). The domain does not depend on adapters. Adapters depend on the domain types or on the port types.

The hexagon picture is a teaching shape. You do not need six sides. You need a center and a rule: details point inward.

This style makes tests easier. You replace a SQL adapter with a memory adapter. You replace the clock with a fixed clock. Topic 3 described interface versus implementation. Hexagonal design applies that idea to the whole application.

Do not create a port for every function. Create a port when a detail must stay replaceable or when the domain must stay free of a vendor type.

Hexagonal design and layers can live together. Presentation adapters sit on one side. Infrastructure adapters sit on the other side. The domain sits in the center.

### Questions

#### Theoretical questions

1. What is a port?
2. What is an adapter?
3. What is the difference between a primary port and a secondary port?
4. Why does the domain depend on ports and not on adapters?
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

## Clean architecture (keep it small)

Clean architecture is a family of rules that Uncle Bob and others popularized. The core rule is the same dependency rule: source code dependencies point inward. Inner circles hold policy. Outer circles hold details.

Typical circles (from the inside):

- Entities: core domain types and rules.
- Use cases: application-specific flows.
- Interface adapters: controllers, presenters, gateways.
- Frameworks and drivers: HTTP framework, database driver.

The names can grow faster than the program. A student project with four circles, twelve folders, and one real feature is over-structured. Keep it small.

A small clean design for a beginner:

- `domain` — entities and rules.
- `app` — use cases (optional if handlers are short).
- `http` — handlers and DTO mapping.
- `store` — SQL adapter.

That set is enough. Add a circle only when a use case script is long or when a second driver appears.

Clean architecture is not a license for many interfaces. It is not a license for a service split. It is a dependency rule plus a warning to keep policy independent of frameworks.

If the team cannot explain the folders in two minutes, the structure is too large.

### Questions

#### Theoretical questions

1. What is the core dependency rule of clean architecture?
2. What do inner circles hold?
3. What do outer circles hold?
4. Why must a beginner keep the structure small?
5. When do you add a use-case folder?

#### Easy practical tasks

1. Write the four classic circles in order from the inside.
2. Map those circles to four folder names for a notes app.
3. Make a table: "Folder" and "May import". Add domain, app, http, and store.
4. Write four sentences that explain "keep it small".

#### Medium practical tasks

1. Take a tutorial that uses eight layers. Reduce it to four folders. Write what you delete.
2. Write an ADR: "we use domain, http, store — we do not use more circles".
3. Review a public "clean architecture" repository. Count folders per feature. Judge if the size fits the features.

#### Advanced practical tasks

1. Write a one-page "minimum clean architecture" for your language with import examples.
2. Compare clean architecture vocabulary with hexagonal vocabulary. Make a term map (entity, use case, port, adapter).

---

## Dependency rule: domain does not import SQL

The dependency rule says that domain code must not import infrastructure details. SQL drivers, HTTP frameworks, and vendor SDKs stay outside the domain.

If the domain imports SQL:

- You cannot test a rule without a database.
- You cannot change the store without a domain rewrite.
- Vendor types leak into business names.
- Transactions and queries become the language of the domain.

The domain may define a port such as `OrderStore` with operations `Save` and `FindByID`. The SQL adapter imports the driver and the domain types. The adapter performs the queries.

The domain may use simple types: numbers, strings, time (through a clock port if tests need a fixed time). The domain must not use a `*sql.DB` value or an HTTP request type.

This rule is the same idea as Topic 3 dependency direction. The style sections only apply the rule with standard names.

Enforcement is practical. Code review can miss an import. A lint rule or a test that fails on forbidden imports is stronger. Write the forbidden list: driver packages, web frameworks, mail SDKs.

An exception is rare. A SQL expression is not a domain rule. If you think you need SQL in the domain, you likely need a better port.

### Questions

#### Theoretical questions

1. What is the dependency rule in one sentence?
2. What goes wrong when the domain imports SQL?
3. How does a store port keep SQL out of the domain?
4. Which types must not appear in the domain?
5. How do you enforce the rule?

#### Easy practical tasks

1. Write a forbidden-import list of six packages or libraries for a typical web app.
2. Draw arrows: domain, `OrderStore` port, SQL adapter, driver.
3. Rewrite a domain function that receives a database handle. Show the new signature.
4. List three tests that you can run without a database after the split.

#### Medium practical tasks

1. Design `OrderStore` operations for place, cancel, and find. Write SQL out of the domain.
2. Find a sample handler online that uses SQL in the handler. Plan a split that also keeps SQL out of the domain.
3. Write a CI idea: fail the build if `domain/` imports a driver.

#### Advanced practical tasks

1. Write a one-page note on transactions: where a transaction starts if the domain cannot import SQL. Compare "unit of work port" versus "application layer".
2. Analyze a leak: a domain type that embeds a vendor JSON tag for a database product. Write the repair.

---

## Mapping DTOs vs domain models

A DTO (data transfer object) is a type that exists for transport. JSON request bodies, JSON response bodies, and row types from a database are DTOs or close cousins. A domain model is a type that holds business meaning and rules.

Do not use one type for all three jobs (HTTP, domain, and SQL). One type forces the domain to accept HTTP names, optional fields for every version, and storage columns.

Map at the edges:

- HTTP adapter: JSON DTO → domain command or domain type.
- SQL adapter: domain type → row, and row → domain type.

Mapping is extra code. That cost is the price of a stable domain. When the JSON field name changes, the domain can stay. When a column name changes, the domain can stay.

Validation split:

- Transport validation: JSON types, required fields, size limits.
- Domain validation: business rules (a discount cannot exceed the total).

Do not skip transport validation. Do not put only transport checks in the domain.

Identity and time need care. An HTTP client may send an id. The domain decides if it accepts that id. A database surrogate key is an infrastructure fact. The domain can hold an identifier type without holding the row type.

Nested JSON is not a nested domain graph by default. Map only what the use case needs. Do not load a full object tree because the DTO has nested objects.

### Questions

#### Theoretical questions

1. What is a DTO?
2. What is a domain model in this section?
3. Why is one type for HTTP, domain, and SQL a problem?
4. Where does mapping run?
5. What is the split between transport validation and domain validation?

#### Easy practical tasks

1. Write a JSON DTO and a domain type for a "create book" command. Show different field names.
2. Make a table: "Change" and "Which type changes". Add JSON rename, column rename, and a new business rule.
3. List four fields that can exist on a DTO but not on the domain type (HTTP-only).
4. Write five sentences on why mapping code is not waste.

#### Medium practical tasks

1. Design mapping for an update of a user profile: PATCH JSON, domain command, SQL row.
2. Write six validation rules and mark each as transport or domain.
3. Show a bad example: a domain type with JSON tags and SQL tags. Write the repair types.

#### Advanced practical tasks

1. Write a one-page policy: when a type may be shared versus when you must map. Include a counter-example.
2. Design versioned DTOs (`v1`, `v2`) that map to one domain command. Topic 8 will reuse this idea.

---

## Over-engineering warning

Over-engineering is a structure that is larger than the problem. Extra layers, extra ports, extra projects, and extra networks appear before a measure or a constraint requires them.

Signals of over-engineering:

- More folders than features.
- An interface for every struct, with one implementation and no test fake.
- A mapper framework for three fields.
- A message broker for a job that a function call can do.
- Copy of a big-company folder tree for a two-person team.

Cost is real. New teammates get lost. Changes touch many empty layers. Performance drops because of extra mapping and extra hops. Topic 2: you cannot maximize all attributes. Complexity spends the maintainability budget.

Rules that keep size honest:

- Start with a modular monolith (Topic 5).
- Add a port when you have two adapters or a test need.
- Add a layer when a handler script is long.
- Add a network hop when a constraint or a measure requires it.
- Delete a folder that has no rule.

A simple design that meets the qualities is better than a "complete" style. You can grow the style. You cannot easily shrink a distributed mud.

Write the reason for each extra band in an ADR. If the reason is "the book had this folder", delete the folder.

### Questions

#### Theoretical questions

1. What is over-engineering in this handbook?
2. Name four signals of over-engineering.
3. How does extra structure harm maintainability?
4. When do you add a port?
5. What reason is not enough for a new folder?

#### Easy practical tasks

1. Write a "too large" folder tree and a "small enough" folder tree for a notes API.
2. Make a table: "Extra item" and "When it is justified". Add port, broker, and fourth layer.
3. List five phrases that hide over-engineering ("for scale", "best practice" with no measure).
4. Write four sentences that you can use in a review to request a smaller design.

#### Medium practical tasks

1. Take a tutorial with many layers. Remove two layers. Write what still works.
2. Write an ADR that rejects a broker for welcome emails in a 50-user app. Name the simpler design.
3. Count interfaces and implementations in a small project. Flag interfaces that have no second implementation and no test fake.

#### Advanced practical tasks

1. Write a one-page "complexity budget" for a semester project: max deploy units, max layers, max ports.
2. Compare two public sample apps (simple and "full clean"). Score them for a beginner team on time-to-first-change.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. What do layered, hexagonal, and clean styles share as one rule?
2. How do these styles support the quality "evolvability" without a service split?
3. When is a DTO leak the same defect as a domain SQL import?
4. How do you choose among three folders and a full clean-circle tree?
5. Why can a style become a big ball of mud if import rules are not enforced?

#### Easy practical tasks

1. Write a one-page cheat sheet: layer names, port names, dependency rule, mapping rule, size rule.
2. Draw one diagram that shows layers and ports for "register user".
3. Write forbidden imports for `domain/` in your language.
4. Bookmark one short explanation of ports and adapters. Write five terms in your words.

#### Medium practical tasks

1. Design a small catalog API with domain, http, and store. Write types and mapping steps. No extra circles.
2. Write three ADRs: style choice, no SQL in domain, and DTO mapping.
3. Refactor a homework handler-on-SQL program on paper into this style. List the file moves.

#### Advanced practical tasks

1. Write a team handbook page (two pages) that teaches the style with one feature end to end.
2. Add a fitness function: forbidden imports plus a test that a use case runs with a memory store. Describe both checks.
