# 4. Domain Modeling

## Description

Domain modeling is the work of naming the business and of drawing boundaries around meaning. This topic covers ubiquitous language, bounded contexts, aggregates, entities, value objects, domain events, and consistency choices inside a domain.

A domain is the subject area of the product (lending books, taking orders, grading exams). A model is the set of types, rules, and names that the software uses for that subject. A weak model uses only technical names (`table1`, `manager`, `data`). A strong model uses the words of the people who do the work.

Complete Topics 1 to 3 before this topic. Keep models small. This handbook uses Domain-Driven Design (DDD) in a practical form. You do not need a large DDD ceremony.

Use one term for each concept. An entity is not a value object. A bounded context is not a deployable service. A domain event is not a command. Transactional consistency is not eventual consistency.

---

## Ubiquitous language and bounded contexts

A ubiquitous language is a shared set of names that the team and the domain experts use in speech, in documents, and in code. One concept has one name. That name appears in types, table names, logs, and ADRs.

If the expert says "loan" and the code says "contract_x", people translate forever. Translation produces defects. A rule that the expert states does not match the function that the developer writes.

Build the language from conversations. Write a short glossary. Include:

- The name
- A one-sentence meaning
- An example
- A rejected synonym

Do not invent names that sound "more technical". Do not use empty names (`Manager`, `Processor`, `Item2`). Prefer `Loan`, `Hold`, `Fine`.

A bounded context is a boundary inside which a model and a language are consistent. Inside the boundary, each name has one meaning. Outside the boundary, the same word can mean something else.

Example: "Order" in a shopping context is a purchase. "Order" in a warehouse context is a pick list. Those two models must not share one type that tries to be both.

The language is local to a bounded context. The word "customer" can mean two things in two contexts. Do not force one global name when the experts disagree.

A bounded context is a design boundary. It can live in a modular monolith as a module. It can later become a service. The context is not automatically a service. Topic 3: do not split the deploy unit until you must.

Between contexts, you translate. A published event or a DTO carries the data that the other context needs. You do not share the inner entities. An anti-corruption layer (Topic 9) is a translator that protects your model from an external model.

How to find a first context map:

- Listen for the same word with different rules.
- Look for different owners (billing team versus catalog team).
- Look for different consistency needs.
- Look for a document or a screen that is a world of its own.

Keep the first map small. Two or three contexts are enough for a student system (catalog, loan, identity). A map with fifteen contexts for one person is over-engineering (Topic 3).

When the language changes, the model must change. A rename in the business is an architecture signal. Update the glossary and the types together.

Use the language in tests. A test named `cannot_return_loan_if_fine_is_open` teaches the rule. A test named `test1` does not.

### Questions

#### Theoretical questions

1. What is a ubiquitous language?
2. Why does translation between expert words and code produce defects?
3. What is a bounded context?
4. Why is a bounded context not automatically a service?
5. What happens at the boundary between two contexts?

#### Easy practical tasks

1. Write a glossary of eight terms for a library. Include one rejected synonym each.
2. Write two meanings of "account" in two contexts (bank versus website login).
3. Draw a context map for a campus shop: catalog, checkout, identity.
4. List four empty names that hide meaning.

#### Medium practical tasks

1. Split a mixed model (users, invoices, warehouse bins) into contexts. Write what each context owns.
2. Design a translation of "customer id" from identity context to loan context.
3. Write an ADR: three contexts in one monolith, not three services.

#### Advanced practical tasks

1. Write a one-page language guide for a new teammate: how to add a term, how to reject a synonym, how to handle a conflict.
2. Write a one-page context map for an airline (booking, check-in, loyalty). Include two translation points.

---

## Aggregates and consistency boundaries

An aggregate is a cluster of domain objects that you change as one unit. The aggregate has a root. All changes go through the root. The root protects invariants.

The consistency boundary is the edge of the aggregate. Inside the boundary, you can use one transaction and you can keep the invariants true together. Across aggregates, you do not use one large transaction as a default. You use a process, an event, or a later check.

Example: a `Loan` aggregate can include the loan header and the loan lines. The root checks that the due date and the item list stay valid. A `CatalogBook` is a different aggregate. A borrow action may need both. You still decide which aggregate commits first and how the other side learns.

Rules:

- Small aggregates. A giant aggregate locks too much data and is hard to understand.
- One transaction per aggregate as the default (later section).
- Identity of the root is how other aggregates refer to it (an identifier, not a shared mutable object).
- Do not hold a live object graph across contexts.

If every write must lock half the database, the aggregate is too large. If every write needs three remote calls to stay consistent, the boundary is too small or the deploy split is too early.

The aggregate is a design idea. In a modular monolith it can be a package plus owned tables. It is not a microservice by itself.

### Questions

#### Theoretical questions

1. What is an aggregate?
2. What is a consistency boundary in this section?
3. Why must changes go through the root?
4. Why is a giant aggregate a problem?
5. How do other aggregates refer to a root?

#### Easy practical tasks

1. Write five sentences that define an aggregate. Use only facts from this section.
2. Draw a `Loan` aggregate and a `Book` aggregate. Mark the root of each.
3. Make a table: "Object" and "Inside Loan? (yes/no)". Add six rows.
4. List four invariants that a loan root can check.

#### Medium practical tasks

1. Design aggregates for a shop: cart or order, catalog item, and payment. Write what each root protects.
2. Explain in eight sentences why "one transaction for the whole company" is not an aggregate design.
3. Write a borrow flow that touches two aggregates. Name the commit order.

#### Advanced practical tasks

1. Write a one-page aggregate catalog for a library: name, root, invariants, and tables.
2. Review a public schema dump. Propose three aggregates and one aggregate that is too large. Write evidence.

---

## Entities vs value objects

An entity is a domain object with identity that stays the same when attributes change. A `Loan` with identifier `L-42` is still that loan after you change the due date.

A value object is a domain object that you define by its attributes. It has no identity of its own. Two `Money` values with the same amount and the same currency are equal. You replace a value object. You do not change it in place as a named person in the model.

Examples of value objects: money, an email address, a date range, an address, a color. Examples of entities: a customer, a loan, a book copy, an invoice.

Rules:

- Put invariants on the type that owns them. `Money` rejects a negative amount if that is the rule.
- Value objects are often immutable in the model. A change produces a new value.
- Do not give a value object a database surrogate key unless the store requires it. The model still treats it as a value.
- Do not treat every table row as an entity. A join row can be a value or a part of an aggregate.

Identity is not "it has a UUID". Identity is "the business tracks this thing through change". A UUID on an address row does not make Address an entity if the business only cares about the fields.

Use value objects to remove primitive obsession. A `string` named `email` does not check format. An `Email` type can.

### Questions

#### Theoretical questions

1. What is an entity?
2. What is a value object?
3. How do you decide if two value objects are equal?
4. Why is a UUID not enough to make an object an entity?
5. What is primitive obsession in one sentence?

#### Easy practical tasks

1. Classify twelve types as entity or value object for a shop.
2. Write five sentences that compare the two ideas. Use only facts from this section.
3. Make a table: "Type" and "Identity?". Add money, customer, due date, and book copy.
4. List four invariants that belong on a `Money` value object.

#### Medium practical tasks

1. Redesign a homework model that uses `string` for money and email. Write the new types and checks.
2. Explain in eight sentences how an address can be a value object in one context and an entity in another.
3. Write equality rules for `ISBN` and for `BookCopy`.

#### Advanced practical tasks

1. Write a one-page modeling guide: when to add a type, when to keep a primitive, and how to persist a value object.
2. Compare two public APIs: one that exposes raw strings and one that uses typed fields. Write mismatches with this section.

---

## Domain events

A domain event is a fact that something in the domain happened. The name is past tense: `LoanReturned`, `FineOpened`, `OrderPlaced`. The event is part of the model. It is not only a log line.

A command asks the system to do work (`ReturnLoan`). After the aggregate accepts the command, the aggregate can record a domain event. Other modules in the same bounded context can react. Other bounded contexts can receive a published form of the event (Topic 7).

A domain event is not a command. A command can fail. An event states that the change already happened in the model.

A domain event is not a technical "row updated" hook. Name the business fact. `TableLoansUpdated` teaches nothing. `LoanReturned` teaches the rule.

Use events to avoid cycles. The loan module does not import the mail module. The loan module records `LoanReturned`. A mail adapter subscribes.

In a modular monolith you can publish in process. Across processes you need a message tool (Topic 7). The name and the meaning stay the same. The transport changes.

Include a time, an identifier of the aggregate, and a correlation identifier when you publish. Topic 11 uses those fields in logs.

Do not use events to hide a missing module boundary. If two modules must change in one user click and you have no name for the fact, the model is incomplete.

### Questions

#### Theoretical questions

1. What is a domain event?
2. How does a domain event differ from a command?
3. Why is a past-tense name useful?
4. How do events help avoid import cycles?
5. Why is `TableLoansUpdated` a weak event name?

#### Easy practical tasks

1. Write six domain event names for a library. Use past tense.
2. Make a table: "Name" and "Command, event, or query". Add twelve rows.
3. Draw `ReturnLoan` then `LoanReturned` then a mail subscriber.
4. List five fields that a published event must include.

#### Medium practical tasks

1. Design the flow: HTTP return → command → aggregate → event → search index update.
2. Rewrite a mixed function `update_and_email` into a command and an event.
3. Write when a module must send a command to another owner versus publish an event.

#### Advanced practical tasks

1. Write a one-page event catalog: ten events, owner context, and who may subscribe.
2. Compare in-process events in a monolith with broker events. Write what stays the same in the model.

---

## Transactional vs eventual consistency

Transactional consistency means that one transaction makes a set of changes visible together. All of the changes commit, or none of them commit. Inside one aggregate and one database, this is the default. Topic 6 and `db.topics.md` cover transactions.

Eventual consistency means that copies or other aggregates become correct after some time. Readers can see an old value for a period. You use eventual consistency when you cannot or must not lock all data in one transaction.

Use transactional consistency when:

- The invariant is in one aggregate.
- The store is one database that you own.
- The user action must see the new state in the same response.

Use eventual consistency when:

- Two bounded contexts have different owners.
- A second process must do work (mail, search, a report).
- A distributed write would need a long lock or a fragile two-phase commit.

Eventual consistency is not "we hope". You must define:

- How the other side learns (event, poll, or file).
- What the user sees while the copy is old.
- How you repair a missed update.
- How you avoid a double effect (idempotency, Topic 2).

A modular monolith can still use eventual consistency for mail and search. It can use one transaction for the loan row and an outbox row (Topic 7). That pattern keeps the write consistent and lets the subscriber lag.

Do not hide a missing invariant behind "eventual". If the business says the two facts must be true together now, you need one aggregate or an explicit process that the user can see.

### Questions

#### Theoretical questions

1. What is transactional consistency in this handbook?
2. What is eventual consistency?
3. When do you use transactional consistency?
4. When do you use eventual consistency?
5. Why is eventual consistency not the same as "we hope"?

#### Easy practical tasks

1. Write five sentences that compare the two consistency types.
2. Make a table: "Action" and "Transactional or eventual". Add pay, send receipt mail, and update search.
3. List four user messages that tell the truth while a copy lags.
4. Draw one database transaction for a loan and a later mail worker.

#### Medium practical tasks

1. Design borrow-book: what commits in one transaction and what may lag. Write the user-visible states.
2. Explain in eight sentences why two services that share a transaction hide a missing API (preview of Topic 6).
3. Write an ADR: mail is eventually consistent; the loan row is transactional.

#### Advanced practical tasks

1. Write a one-page consistency map for a shop: six actions, consistency type, and repair if a subscriber fails.
2. Compare a two-phase commit slogan with an outbox plus retry. Write operations cost for a four-person team. Do not implement a protocol.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do ubiquitous language, bounded context, and aggregate nest in one model?
2. Why can one word need two models, and how does that affect entities you share?
3. When is a domain event a better tool than a direct import of another module?
4. How do value objects and aggregates share the work of protecting invariants?
5. What quality attributes from Topic 1 change when you pick eventual consistency for a subscriber?

#### Easy practical tasks

1. Write a one-page cheat sheet: glossary, context, aggregate, entity, value object, event, two consistency types.
2. For a to-do app, write four glossary terms, one context, one aggregate, two events, and one eventual subscriber (mail).
3. Draw a context map with two boxes and one translated identifier.
4. Bookmark a short public DDD glossary. Write three terms that match this topic and one term that this topic omitted.

#### Medium practical tasks

1. Write a short domain brief for a campus lost-and-found: language, two contexts, aggregates, and consistency for notify-mail.
2. Take a schema with `tbl_x` names. Rename it to a ubiquitous language on paper. Show before and after.
3. Write a twelve-week modeling plan: glossary first, then aggregates, then events.

#### Advanced practical tasks

1. Write a one-page model review checklist that a teammate uses on a pull request (names, boundaries, consistency).
2. Compare this practical DDD subset with a large DDD book table of contents. Map five book themes to this file and mark ceremony that you skip.
