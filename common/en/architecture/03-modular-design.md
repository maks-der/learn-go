# 3. Modular Design

## Description

A module is a part of a system with a boundary and a responsibility. Modular design is the work of choosing those boundaries. This topic covers cohesion, coupling, encapsulation, interfaces, dependency direction, cycles, and packages as boundaries.

Good modules let you change one part without a rewrite of the other parts. Bad modules spread one change across the whole system. Topic 2 called that effect maintainability.

Complete Topics 1 and 2 before this topic. Write small modules first. Do not start with many network services. A module can live in one process.

Use one term for each concept. A module is a design boundary. A package is a language boundary. A library is a distributable module. Those terms are close. This handbook keeps them distinct.

---

## Cohesion and coupling

Cohesion is how strongly the elements inside a module belong together. A module with high cohesion does one job. Examples: "price an order", "hash a password", "write an audit line". A module with low cohesion mixes unrelated jobs. Example: "prices, emails, and PDF layout in one file".

Coupling is how strongly one module depends on another module. A module with high coupling cannot change without changes in many other modules. A module with low coupling depends on a small, stable interface.

Aim for high cohesion and low coupling. Those two goals work together. When a module does one job, other modules need less knowledge of its internals.

Some coupling is necessary. A checkout module must call a price module. That coupling is acceptable when the price interface is small and stable.

Hidden coupling is worse than visible coupling. Hidden coupling includes a shared global variable, a shared database table with no owner, and a hidden assumption about file format. Visible coupling is an import or a documented interface.

A change request is a test of the design. If a small business change touches many modules, cohesion is low or coupling is high. Record that signal. Topic 5 uses the same signal for a "big ball of mud".

### Questions

#### Theoretical questions

1. What is cohesion?
2. What is coupling?
3. Why do high cohesion and low coupling work together?
4. When is coupling acceptable?
5. What is hidden coupling?

#### Easy practical tasks

1. Write five sentences that define cohesion and coupling. Use only facts from this section.
2. Classify six folders in a project that you know as high or low cohesion. Give one reason each.
3. Make a table: "Dependency" and "Visible or hidden". Add four examples.
4. List three change requests for a shop. For each request, guess how many modules you must touch.

#### Medium practical tasks

1. Split a low-cohesion module (auth + email + reports) into three modules. Write the new names and the remaining calls.
2. Draw two module diagrams of the same system: one with high coupling, one with lower coupling. Label the interfaces.
3. Take a public small repository. Find one file that mixes two jobs. Write a split plan of ten lines.

#### Advanced practical tasks

1. Write a one-page review of coupling types (content, common, control, stamp, data) with one software example each. Keep the language simple.
2. Measure a codebase with a simple count: incoming and outgoing dependencies per package. Interpret the two worst packages.

---

## Encapsulation

Encapsulation is the rule that a module hides its internals and shows only a small interface. Internals include data structures, SQL, and helper functions. Callers must not depend on those internals.

Encapsulation protects invariants. An invariant is a rule that must stay true. Example: "an order total equals the sum of line totals." If every caller writes to the same tables, the invariant breaks.

Access through functions or methods lets the module check the rule. Direct access to fields or tables bypasses the check.

Encapsulation is not the same as "make every field private and add getters". A getter that exposes the full internal structure is a leak. Hide the structure. Expose operations that match the job: `placeOrder`, not `getOrderMap`.

Language features help: private names, packages, and modules. They do not replace design. A public package that exports twenty types has a large surface. A large surface is weak encapsulation.

When you change internals and you do not change the interface, callers continue to work. That is the practical test of encapsulation.

### Questions

#### Theoretical questions

1. What is encapsulation?
2. What is an invariant?
3. Why does direct table access weaken encapsulation?
4. Why can a getter leak internals?
5. What is the practical test of encapsulation?

#### Easy practical tasks

1. Write four operations that a "wallet" module can expose. Write three internals that it must hide.
2. Make a table: "Exposed item" and "Leak? (yes/no)". Add six items for a user profile module.
3. Explain in five sentences how encapsulation protects an invariant for a bank account balance.
4. List language tools that hide names in a language that you know.

#### Medium practical tasks

1. Redesign a module that exposes a raw SQL row type to callers. Write the new operations and the hidden type.
2. Find a public class or package with many public fields. Write an encapsulation plan. Do not paste large source.
3. Write two tests that fail if an invariant breaks. Write the invariant in one sentence.

#### Advanced practical tasks

1. Compare encapsulation at the type level and at the process level (a service that does not share its database). Write eight sentences.
2. Design a module API for a calendar booking rule set so that you can change storage from files to SQL without a caller change.

---

## Interface vs implementation

An interface is the contract that a caller can use. The contract includes operations, types, errors, and meaning. An implementation is the code and data that satisfy the contract.

Callers depend on the interface. Callers must not depend on the implementation. If they do, you cannot replace the implementation.

A stable interface changes slowly. An implementation can change often. That split is how you evolve a system. Topic 8 applies the same idea to public APIs.

An interface that mirrors one database table is a weak interface. It copies the implementation into the contract. A better interface uses operations of the domain: `reserveSeat`, not `updateSeatsColumn`.

More than one implementation can satisfy one interface. Examples: a fake in-memory store for tests and a SQL store for production. Tests become easier. Topic 6 uses this idea in ports and adapters.

Do not create an interface for every type. An interface with one implementation and no test fake can be extra noise. Add an interface when you have two implementations or a clear stability need.

Document the meaning. "Returns the user" is incomplete if "user" can be missing. State empty cases and errors.

### Questions

#### Theoretical questions

1. What is an interface in this handbook?
2. What is an implementation?
3. Why must callers depend on the interface only?
4. When is an interface that copies a table a weak interface?
5. When do you add an interface, and when do you wait?

#### Easy practical tasks

1. Write an interface of five operations for a key-value store. Do not mention files or SQL.
2. List three implementations that could satisfy that interface.
3. Make a table: "Change" and "Breaks the interface? (yes/no)". Add six changes.
4. Write four sentences that explain a test fake.

#### Medium practical tasks

1. Design a `Clock` interface so that tests can fix the time. Write the production implementation idea and the test implementation idea.
2. Take a function that opens SQL inside the caller. Split interface and implementation on paper.
3. Document error cases for a `GetUser(id)` operation. Include not found and store down.

#### Advanced practical tasks

1. Write a one-page rule for your team: when a new interface is required. Include a counter-example of interface noise.
2. Compare a language `interface` type with an HTTP API as an interface. What is the same? What is different?

---

## Dependency direction

A dependency points from the module that needs a service to the module that provides the service. In code, an import is a dependency. In processes, a client depends on a server.

Dependency direction must follow stability and responsibility. A high-level policy module must not depend on a low-level detail such as a SQL dialect. The detail depends on an interface that the policy defines. Topic 6 calls this the dependency rule.

If the domain module imports the HTTP framework, the domain is stuck to that framework. If the HTTP adapter imports the domain, you can change the framework.

Draw arrows from the dependent to the provider. Cycles in those arrows are a problem. The next section covers cycles.

Unstable modules must not sit at the center of many arrows. A report layout that changes every week must not be the dependency of the order core.

Use dependency inversion when two modules must talk and you must keep the core free of details. The core defines the interface. The detail implements the interface.

Package rules, lint rules, and code review can enforce direction. A rule that is only in a slide will decay.

### Questions

#### Theoretical questions

1. What is a dependency in code?
2. Why must a domain module not import a SQL dialect?
3. What is dependency inversion in one sentence?
4. Why is an unstable module a bad center of dependencies?
5. How can a team enforce dependency direction?

#### Easy practical tasks

1. Draw three boxes: HTTP, domain, SQL. Draw allowed arrows. Draw a forbidden arrow.
2. Write five sentences that explain why the HTTP layer can depend on the domain.
3. List four "details" (SQL, SMTP, file system, payment vendor) that the core must not import.
4. Make a table: "Import" and "Allowed? (yes/no)" for six import pairs.

#### Medium practical tasks

1. Redesign a sketch where `order` imports `smtp`. Introduce an interface. Show the new arrows.
2. Walk a repository that you own. Find one import that points the wrong way. Write a fix plan.
3. Write a package lint wish-list (three rules) that would catch wrong-direction imports.

#### Advanced practical tasks

1. Write a one-page comparison of "core depends on libraries" versus "libraries adapt to core". Use one feature as the example.
2. Design a dependency graph for a billing system with six modules. Mark each arrow with "policy" or "detail".

---

## Acyclic dependencies

A dependency cycle exists when module A depends on B and B depends on A, directly or through other modules. Cycles make build order hard. Cycles make change hard. A change in A can require a change in B that requires a change in A.

Cycles also hide ownership. If two modules import each other, neither module is the owner of the shared idea. The shared idea needs a third module or a clearer split.

Break a cycle with one of these methods:

- Move the shared type or interface to a third module.
- Invert one dependency with an interface.
- Merge the two modules if they are one job.
- Pass data as values so that one side does not need the other type.

A cycle of three or more modules is still a cycle. Draw the graph. Do not only look at pairs.

Some languages allow import cycles and fail at compile time. Some languages allow them and fail at runtime. Treat a cycle as a design defect even if the compiler is silent.

Acyclic dependencies let you test a module with fakes at the edges. A cycle often forces you to boot half of the system for one test.

### Questions

#### Theoretical questions

1. What is a dependency cycle?
2. Why do cycles make change hard?
3. Why do cycles hide ownership?
4. Name three methods to break a cycle.
5. Why do cycles harm tests?

#### Easy practical tasks

1. Draw a cycle of two modules and a cycle of three modules. Label the arrows.
2. Write four sentences on why a shared type can live in a third module.
3. List three signals that a cycle exists (build error, odd import, test setup).
4. Make a table: "Fix" and "When to use it". Add merge, invert, and extract.

#### Medium practical tasks

1. Given A = orders, B = users, and each needs a type from the other, break the cycle on paper.
2. Find a cycle in a project or invent a realistic one from a blog. Write the break plan in ten steps.
3. Write two tests that become possible after you break a cycle.

#### Advanced practical tasks

1. Write a one-page note on layering versus a cycle: how a three-layer rule prevents some cycles and misses others.
2. Design a weekly review that detects new cycles (tool or manual graph). Include a fail rule for the pipeline.

---

## Packages / modules / libraries as boundaries

A package is a language unit: a directory of source files with one name. A module (in design) is a responsibility boundary. A library is a package or a set of packages that you distribute and version.

These units are tools for boundaries. A folder is not a boundary if every type is public and every package imports every other package.

Put one cohesive job in one package. Name the package after the job, not after a technical layer only. `pricing` is clearer than `utils`. `utils` and `common` grow into low-cohesion bags.

A library boundary is stronger than a folder. A library has a version. A change of the public types is a release. Use a library when more than one program must share the code and you can accept release cost.

An `internal` package (in some languages) blocks other modules from import. Use that tool for code that must stay private to one program.

Do not split a small program into many libraries. The release cost exceeds the benefit. Start with packages in one program. Extract a library when a second program needs the same code and the interface is stable.

Deployment boundaries are stronger again. Topic 5 covers the monolith as one deployable unit with many packages. Topic 12 covers services that deploy apart. Do not confuse a package split with a service split.

### Questions

#### Theoretical questions

1. What is a package in this handbook?
2. What is the difference between a design module and a library?
3. Why is `utils` a risky package name?
4. When do you extract a library?
5. Why is a package split not the same as a service split?

#### Easy practical tasks

1. Draw a folder tree for a shop with packages `cmd/app`, `pricing`, `cart`, and `internal/store`.
2. Write five package names that show a job. Write three package names that hide a job.
3. Make a table: "Boundary" and "Cost of change". Add package, library, and deployable service.
4. List four rules for what may be public in a package.

#### Medium practical tasks

1. Take a single-folder program. Propose a package split of four packages. Write import rules.
2. Write an ADR: "we keep one module (one program) and we do not extract a library yet". Give three reasons.
3. Design an `internal` area and a public library area for a toolkit that a second program will use next year.

#### Advanced practical tasks

1. Compare language module systems (Go modules, Java modules, or another pair). Write how each enforces a boundary.
2. Plan an extract of a library from a monolith: version, public types, changelog, and the first consumer.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do cohesion, coupling, and encapsulation work together in one sentence each, then in one combined paragraph?
2. How does dependency direction protect evolvability from Topic 2?
3. When does a small interface hide a large implementation change?
4. Why do cycles and hidden coupling produce the same symptom for a change request?
5. How do you decide among a package, a library, and a later service split?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic with one example each.
2. For a forum application, name six modules, their interfaces (two operations each), and forbidden imports.
3. Draw an acyclic dependency graph for those six modules.
4. Review your last homework program. Write three modular defects and one fix each.

#### Medium practical tasks

1. Write a modular design for a library loan system in two pages: modules, interfaces, dependency arrows, and two invariants.
2. Convert a "utils" bag of twelve functions into cohesive packages. Show the mapping table.
3. Write three ADRs: package layout, no library extract, and a rule against import cycles.

#### Advanced practical tasks

1. Design a fitness function (automated check) for modular rules: no forbidden imports, max package size, no cycles. Specify the tool idea and the fail behavior.
2. Compare the module map of a well-known open-source application with this topic. Write what you would change for a five-person team.
