# 7. Domain Modeling

## Description

Domain modeling is the work of naming the business and of drawing boundaries around meaning. This topic covers ubiquitous language, bounded contexts, aggregates, entities, value objects, domain events, and consistency choices inside a domain.

A domain is the subject area of the product (lending books, taking orders, grading exams). A model is the set of types, rules, and names that the software uses for that subject. A weak model uses only technical names (`table1`, `manager`, `data`). A strong model uses the words of the people who do the work.

Complete Topics 1 to 6 before this topic. Keep models small. This handbook uses Domain-Driven Design (DDD) in a practical form. You do not need a large DDD ceremony.

Use one term for each concept. An entity is not a value object. A bounded context is not a deployable service. A domain event is not a command.

---

## Ubiquitous language

A ubiquitous language is a shared set of names that the team and the domain experts use in speech, in documents, and in code. One concept has one name. That name appears in types, table names, logs, and ADRs.

If the expert says "loan" and the code says "contract_x", people translate forever. Translation produces defects. A rule that the expert states does not match the function that the developer writes.

Build the language from conversations. Write a short glossary. Include:

- The name
- A one-sentence meaning
- An example
- A rejected synonym

Do not invent names that sound "more technical". Do not use empty names (`Manager`, `Processor`, `Item2`). Prefer `Loan`, `Hold`, `Fine`.

The language is local to a bounded context (next section). The word "customer" can mean two things in two contexts. Do not force one global name when the experts disagree.

When the language changes, the model must change. A rename in the business is an architecture signal. Update the glossary and the types together.

Use the language in tests. A test named `cannot_return_loan_if_fine_is_open` teaches the rule. A test named `test1` does not.

### Questions

#### Theoretical questions

1. What is a ubiquitous language?
2. Why does translation between expert words and code produce defects?
3. What does a glossary entry contain?
4. Why can one word have two meanings in two contexts?
5. Why do tests use the language?

#### Easy practical tasks

1. Write a glossary of eight terms for a library. Include one rejected synonym each.
2. Rewrite five technical names (`UserMgr`, `TblOrd`) into domain names.
3. Make a table: "Expert sentence" and "Type or function name". Add five rows.
4. List four empty names that hide meaning.

#### Medium practical tasks

1. Listen to or imagine a 10-minute talk with a librarian. Extract twelve terms and three rules.
2. Rename a homework model on paper so that names match a real domain. Show before and after.
3. Find two documents (UI and database) that use different names for the same idea. Write a unification plan.

#### Advanced practical tasks

1. Write a one-page language guide for a new teammate: how to add a term, how to reject a synonym, how to handle a conflict.
2. Compare the public vocabulary of a product (help pages) with a typical schema dump from a similar app. Write mismatches.

---

## Bounded context (DDD, practical)

A bounded context is a boundary inside which a model and a language are consistent. Inside the boundary, each name has one meaning. Outside the boundary, the same word can mean something else.

Example: "Order" in a shopping context is a purchase. "Order" in a warehouse context is a pick list. Those two models must not share one type that tries to be both.

A bounded context is a design boundary. It can live in a modular monolith as a module. It can later become a service. The context is not automatically a service. Topic 5: do not split the deploy unit until you must.

Between contexts, you translate. A published event or a DTO carries the data that the other context needs. You do not share the inner entities. An anti-corruption layer (later path) is a translator that protects your model from an external model.

How to find a first context map:

- Listen for the same word with different rules.
- Look for different owners (billing team versus catalog team).
- Look for different consistency needs.
- Look for a document or a screen that is a world of its own.

Keep the first map small. Two or three contexts are enough for a student system (catalog, loan, identity). A map with fifteen contexts for one person is over-engineering (Topic 6).

Draw a context map: names, relations (upstream/downstream or "this context publishes X"), and the owner of each context.

### Questions

#### Theoretical questions

1. What is a bounded context?
2. Why can one word need two models?
3. Why is a bounded context not automatically a service?
4. What happens at the boundary between two contexts?
5. Name three signals that help you find a context?

#### Easy practical tasks

1. Write two meanings of "account" in two contexts (bank versus website login).
2. Draw a context map for a campus shop: catalog, checkout, identity.
3. Make a table: "Term" and "Meaning in context A / B". Add four terms.
4. List four types that must not be shared across your two "Order" models.

#### Medium practical tasks

1. Split a mixed model (users, invoices, warehouse bins) into contexts. Write what each context owns.
2. Design a translation of "customer id" from identity context to loan context.
3. Write an ADR: three contexts in one monolith, not three services.

#### Advanced practical tasks

1. Write a one-page context map for an airline (booking, check-in, loyalty). Include two translation points.
2. Take a public API that mixes two languages (commerce and content). Propose a context split and the published types.

---

## Aggregates and consistency boundaries

An aggregate is a cluster of domain objects that you change as one unit. The aggregate has a root. All changes go through the root. The root protects invariants.

A consistency boundary is the border inside which you require an immediate, consistent update (often one transaction). The aggregate is the usual consistency boundary in this handbook.

Example: an `Order` root plus `OrderLine` objects. The invariant "sum of lines equals total" is checked on the root. You do not update a line from outside the root.

Do not make a large aggregate. A large aggregate locks more data and creates contention. If two users can change two parts independently, those parts are likely two aggregates.

References across aggregates use identifiers, not object graphs that you save in one transaction. `Order` can store `CustomerID`. `Order` does not store the full `Customer` graph if customer is another aggregate.

One user action can touch more than one aggregate. Then you cannot have one simple invariant across both unless you use one transaction across both (tight coupling) or you accept eventual consistency (later section). Prefer small aggregates and explicit processes.

Persistence often stores one aggregate per transaction. That rule keeps the boundary real. If every use case updates ten aggregates in one transaction, you do not have ten aggregates. You have one hidden blob.

### Questions

#### Theoretical questions

1. What is an aggregate?
2. What is an aggregate root?
3. What is a consistency boundary in this handbook?
4. Why must aggregates stay small?
5. How do aggregates refer to each other?

#### Easy practical tasks

1. Draw an `Order` aggregate with a root and lines. Write one invariant.
2. Make a table: "Object" and "Root or part". Add order, line, customer, and product.
3. Write four operations that must go through the `Loan` root.
4. List three invariants for a `ShoppingCart` aggregate.

#### Medium practical tasks

1. Split a "god" `School` object (students, rooms, grades) into aggregates. Justify each boundary.
2. Design identifiers between `Invoice` and `Customer`. Write what a transaction may update.
3. Write two use cases: one that stays in one aggregate, one that needs two. Explain consistency.

#### Advanced practical tasks

1. Write a one-page rule: how to choose an aggregate from a contention problem (two users edit the same root).
2. Compare a fine-grained aggregate design with a coarse design for a calendar. Use locking and user conflict as criteria.

---

## Entities vs value objects

An entity is an object with a stable identity. Two entities can have the same attributes and still be different. Example: two users named Ada with different ids. You track an entity across time.

A value object is an object that you define by its attributes. Two value objects with the same attributes are the same. Example: money amount `10 EUR`, a date range, an address (if you do not give the address its own identity). You can replace a value object. You do not update it in place as a tracked identity.

Use entities for things that the business tracks: `Loan`, `Customer`, `Invoice`. Use value objects for measures, codes, and small compounds: `Money`, `ISBN`, `EmailAddress`.

Value objects are often immutable. A change of an address on a person is a replacement of the value, not a silent mutation of a shared object.

Do not give a database row id to every value. An id does not make an entity. The business does. A `Money` row id is usually a leak of storage into the model.

Equality differs:

- Entities: equal if the identity is equal.
- Value objects: equal if the attributes are equal.

Validation belongs on both. An `EmailAddress` value can reject a bad string. An entity can reject a forbidden state change.

### Questions

#### Theoretical questions

1. What is an entity?
2. What is a value object?
3. How does equality differ?
4. Why is a database id not enough to make an entity?
5. Why are value objects often immutable?

#### Easy practical tasks

1. Classify twelve concepts of a shop as entity or value object.
2. Write equality rules for `User` and for `Money`.
3. Make a table: "Concept" and "Identity? (yes/no)". Add ISBN, loan, color, and payment.
4. Write five sentences that explain replacement of an address value.

#### Medium practical tasks

1. Design `Money` (amount plus currency) as a value object. Write invalid combinations.
2. Redesign a model where `Address` has a table id but no business identity. Show the value-object form.
3. Write two invariants: one on a value object, one on an entity state change.

#### Advanced practical tasks

1. Write a one-page note on when an address becomes an entity (shipping locations that you track over years).
2. Compare language tools (records, structs, classes) for value objects. Write how you implement equality.

---

## Domain events

A domain event is a record that something meaningful happened in the domain. The name is in the past tense: `LoanReturned`, `PaymentCaptured`, `SeatReserved`. The event is a fact. You do not undo a fact. You can record a later compensating fact.

Domain events help communication inside a model and across contexts. Other modules can react. A `LoanReturned` event can start a fine calculation in another module without a direct import cycle.

An event is not a command. A command is a request (`ReturnLoan`). The command can fail. An event is the result after a successful change (`LoanReturned`). Topic 10 covers commands versus events versus queries in messaging.

Keep the payload small. Include identifiers, time, and the data that a subscriber needs. Do not attach the full aggregate graph if a subscriber only needs the id and the return date.

You can use events in process (function calls to listeners) first. You do not need a broker to have a domain event. Add a broker when a process or a team boundary requires it (Topic 10).

Write events in the ubiquitous language. `RowUpdated` is not a domain event. `HoldExpired` is a domain event.

Store events when other systems must not miss them. The outbox pattern (Topic 10) is one method. In a small monolith, a table of events or an in-process dispatcher can be enough.

### Questions

#### Theoretical questions

1. What is a domain event?
2. How does an event differ from a command?
3. Why does the name use the past tense?
4. Why can events reduce import cycles?
5. Why do you not need a broker to have a domain event?

#### Easy practical tasks

1. Write eight domain event names for a library. Reject two technical names.
2. Make a table: "Command" and "Event if success". Add four pairs.
3. List fields for `BookHeld`: ids, time, and two extra facts.
4. Write four sentences on why `RowUpdated` is a weak event.

#### Medium practical tasks

1. Design in-process listeners for `LoanReturned` (fines, notifications). Write the module arrows.
2. Choose the payload for `OrderPlaced` so that a mail module does not need the order table.
3. Write three events that a billing context would consume from a shop context.

#### Advanced practical tasks

1. Write a one-page event catalog for one context: name, payload, producer, consumers, and consistency note.
2. Design a version rule for event payloads (add a field, never reuse a name). Link to Topic 8 ideas.

---

## Transactional vs eventual consistency inside a domain

Transactional consistency means that a set of changes commits as one unit. After the commit, all those changes are visible. If the commit fails, none of the changes are visible. A database transaction is the usual tool.

Eventual consistency means that the system can show a temporary disagreement. A later process, message, or retry makes the views agree. The user can see a short delay.

Inside one aggregate, prefer transactional consistency. The root invariant must hold after the transaction.

Across aggregates or across contexts, eventual consistency is often the honest model. A payment captured and a loyalty point added are two facts. They can use two transactions and a `PaymentCaptured` event.

Do not pretend that a distributed system is transactionally consistent if it is not. A network call inside a database transaction is a common defect. The remote work can succeed and the local commit can fail, or the reverse.

Choose consistency from the business. "The shelf count and the loan record must match in the same second" can force one transaction or one aggregate. "The weekly report can lag 5 minutes" must not force a distributed transaction.

Measures help. Write the maximum lag that a stakeholder accepts. That lag is a quality attribute (Topic 2).

Compensating actions handle failure in eventual designs: a `PaymentFailed` fact after a reserved seat, or a timeout that releases the seat. Topic 12 later covers sagas. You can start with a simple retry and a manual repair list.

### Questions

#### Theoretical questions

1. What is transactional consistency?
2. What is eventual consistency?
3. Where do you prefer transactional consistency?
4. Why is a network call inside a database transaction a defect?
5. How does the business decide the consistency style?

#### Easy practical tasks

1. Write four sentences that compare the two consistency styles.
2. Make a table: "Rule" and "Transactional or eventual". Add six business rules for a shop.
3. List three user-visible delays that are acceptable and one that is not for a payment.
4. Write a compensating action for a reserved seat when payment fails.

#### Medium practical tasks

1. Design borrow-book with one aggregate and one transaction. Then write what happens if fines live in another aggregate.
2. Draw a timeline of eventual consistency for "order placed" and "inventory reserved". Show a window of disagreement.
3. Write an ADR that rejects a distributed transaction for email after order.

#### Advanced practical tasks

1. Write a one-page decision sheet: invariant, max lag, owner, and style (transactional or eventual). Fill five rows.
2. Research a public explanation of sagas at a high level. Map it to compensating events in this section. Do not copy long text.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do ubiquitous language and bounded context stop a single "User" type from absorbing the whole product?
2. How do aggregates and consistency choices protect invariants without a large lock?
3. When is a domain event the right tool to keep modules acyclic (Topic 3)?
4. How do entities and value objects change the way you design equality and persistence?
5. Why does this topic stay useful inside a modular monolith?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms with one library example each.
2. Build a mini model: glossary (six terms), two contexts, two aggregates, four events.
3. Draw a context map and one aggregate diagram for a bike rental.
4. List ten names from your last project and mark language, entity/value, and context.

#### Medium practical tasks

1. Write a two-page domain design for a cafeteria card: language, contexts, aggregates, events, and consistency table.
2. Write two ADRs: aggregate boundaries, and in-process events (no broker).
3. Take a CRUD schema (five tables). Propose a model that is not only rows. Show invariants.

#### Advanced practical tasks

1. Write an event-and-aggregate paper design that you could implement in one process with one database.
2. Compare a DDD overview chapter (public) with this practical subset. Write what you skip on purpose for a small team.
