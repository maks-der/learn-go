# 1. Architecture and Quality Attributes

## Description

Software architecture is the structure of a system and the decisions that are costly to change. This topic defines architecture. It also defines quality attributes. You learn functional requirements, quality requirements, stakeholders, and constraints. You also learn how to record a decision in an Architecture Decision Record (ADR).

Complete this topic before you design modules or styles. Pair this path with a language path and with `db.topics.md` and `net.topics.md`.

Use one term for each concept. Do not use "architecture" as a synonym for "framework" or "diagram". Architecture is structure, costly decisions, and trade-offs. You cannot maximize all quality attributes at the same time.

---

## What architecture is (structure plus decisions that are costly to change)

Architecture is the structure of a software system. Structure includes the parts of the system and the relations between those parts. Parts can be modules, processes, data stores, and interfaces.

Architecture is also the set of decisions that are costly to change later. A costly decision is a choice that forces a large rewrite if you reverse it. Examples: one database or many databases, one deployable unit or many deployable units, a synchronous call or a message.

A design choice is architectural when a change of that choice has a large cost. The color of a button is not architectural. The choice of the data store that holds the source of truth is architectural.

Architecture exists even when you do not write a document. Every system has a structure. If you do not choose the structure, the structure appears from daily code changes. That accidental structure is still architecture. It is often hard to change.

Architecture work is about trade-offs. A trade-off is a choice that improves one quality and can decrease another quality. Record the structure and the costly decisions. Use simple diagrams and short text. The C4 model is one method to show context and containers. See [https://c4model.com/](https://c4model.com/).

Architecture is not a framework name. Spring, Django, and similar tools are not architecture. You can use the same framework for many different structures. Architecture is not a fashion word. Words such as "cloud native" do not specify structure, constraints, or trade-offs. A box-and-line picture that ignores time, team size, law, and existing systems is not a design.

### Questions

#### Theoretical questions

1. What two parts make architecture in this handbook?
2. When is a design choice architectural?
3. Give one example of a costly decision and one example of a cheap decision.
4. What is a trade-off?
5. Does a system have architecture when no diagram exists? Explain.

#### Easy practical tasks

1. Write five sentences that define architecture. Use only facts from this section.
2. Make a two-column table: "Decision" and "Costly to change? (yes/no)". Add six rows from a system that you know.
3. List four parts that can appear in a system structure. Give one relation between two of those parts.
4. Open [https://c4model.com/](https://c4model.com/). Write the names of the four C4 levels in order.

#### Medium practical tasks

1. Pick a small program that you wrote. Name three decisions that are costly to change. Give one reason for each name.
2. Draw a simple context diagram: users, your system, and two external systems. Label each relation with one verb.
3. Rewrite these three phrases into useful statements or mark them as empty: "we are cloud native", "we use Kubernetes", "we have a modern API".

#### Advanced practical tasks

1. Read the C4 context and container pages. Draw both views for a library loan system. Write five sentences that explain what each view hides.
2. Find a public post-mortem of a production incident. Name one architectural decision that made the incident larger or smaller. Give evidence from the text.

---

## Functional vs quality requirements

A requirement is a condition that the system must satisfy. Teams split requirements into two groups.

A functional requirement says what the system must do. Examples: "A user can place an order." "The system sends a receipt." Functional requirements describe behavior and features.

A quality requirement says how well the system must do the work. Quality requirements are quality attributes. Examples: latency, availability, security, and cost. Later sections in this topic name the common attributes.

A quality attribute needs a measure. "The system must be fast" is not a requirement. "P99 latency of the checkout HTTP call is less than 300 ms in region A" is a requirement. "The service is secure" is not a requirement. "Only authenticated users can read order records" is closer to a requirement.

Functional requirements and quality attributes interact. A new feature can increase latency. A security control can decrease throughput. Record both types. Do not treat quality as an afterthought.

Write requirements so that a test or a review can check them. If no person can say pass or fail, the text is a wish, not a requirement.

This handbook uses "constraint" for team, time, law, and existing systems. Use "quality attribute" for reliability, performance, and similar properties. Do not mix the two words.

### Questions

#### Theoretical questions

1. What is a functional requirement?
2. What is a quality requirement?
3. Why must a quality attribute have a measure?
4. How can a new feature change a quality attribute?
5. When is a requirement only a wish?

#### Easy practical tasks

1. Classify ten statements as functional, quality, constraint, or wish. Use a four-column table.
2. Rewrite "the API must be fast" into two measurable quality requirements.
3. Write three functional requirements for a note-taking application.
4. Write three quality requirements for the same application. Include a number or a clear yes/no check.

#### Medium practical tasks

1. Take a public product page. Extract five functional items and five quality items. Mark items that have no measure.
2. Write a one-page requirement list for a classroom booking system. Include both types. Add one conflict between a feature and a quality.
3. For one quality requirement, write how you would measure it in a test or in production.

#### Advanced practical tasks

1. Write a quality attribute scenario (source, stimulus, environment, response, measure) for availability and one for security.
2. Turn six vague wishes into testable requirements. Show the before text and the after text.

---

## Stakeholders and constraints

A stakeholder is a person or a group that has an interest in the system. Stakeholders include users, the product owner, developers, operators, security staff, legal staff, and the paying organization.

Different stakeholders care about different qualities. A user cares about latency and clear errors. An operator cares about logs and deploy time. A finance stakeholder cares about cost. A legal stakeholder cares about retention and consent.

Architecture work includes the identification of stakeholders. If you miss a stakeholder, you miss a constraint or a quality. The missed item appears later as a surprise.

A stakeholder can have a veto. A security team can block a design that stores secrets in source control. A legal team can block a design that sends personal data to another region.

You cannot satisfy every stakeholder fully. Record who accepts each risk. A decision without an owner is not a complete decision.

A constraint is a limit that you cannot ignore. Constraints reduce the set of possible designs. A design that ignores a constraint fails in operation or in audit.

Team constraints include skill, size, and working hours. A team of two people cannot operate twenty services with the same quality as a large platform team.

Time constraints include a launch date and a seasonal peak. A design that needs six months of migration is not valid when the launch is in four weeks.

Law constraints include privacy, retention, accessibility, and licenses. Law can force data location, consent records, and audit logs. Law is not optional.

Existing systems are constraints. You must integrate with a payment provider, a directory, or an old database. A green-field diagram that replaces all existing systems is a wish if those systems stay.

Write constraints in the same document as decisions. A future reader must see why you rejected a "better" structure.

### Questions

#### Theoretical questions

1. What is a stakeholder?
2. What happens when a design misses a stakeholder?
3. What is a constraint in this handbook?
4. How does team size change the number of deployable units that you can operate?
5. Why is an existing system a constraint?

#### Easy practical tasks

1. List eight stakeholder types for an online shop. Write one quality that each type cares about.
2. Make a two-column table: "Stakeholder" and "Fear if the system fails". Add six rows.
3. Classify eight limits as team, time, law, or existing system.
4. Write three sentences that explain why "rewrite everything" can fail as a plan.

#### Medium practical tasks

1. For a hospital appointment system, map five stakeholders to quality attributes. Mark two conflicts.
2. Given: three developers, twelve weeks, one legacy SQL database, and a privacy rule for student data. Propose a structure in one page. Name two designs that you reject.
3. Role-play on paper: write a short statement from a user, an operator, and a finance person about the same design. Then write one compromise.

#### Advanced practical tasks

1. Build a RACI-style table (responsible, accountable, consulted, informed) for five architecture decisions in a small company system.
2. Write a constraint register with columns: constraint, owner, source, impact on structure, date. Fill twelve rows for a fictional bank transfer tool.

---

## Architecture Decision Records

An Architecture Decision Record (ADR) is a short document that records one architecture decision. The ADR states the context, the decision, and the consequences. Teams keep ADRs in the repository. See [https://adr.github.io/](https://adr.github.io/).

Write an ADR when the decision is costly to change. Examples: the primary database, the authentication method, the deploy shape (one process or many), and the public API style.

A typical ADR includes:

- Title and a number
- Status: proposed, accepted, superseded
- Context: facts, constraints, and the problem
- Decision: the choice in one short paragraph
- Consequences: good results and bad results
- Date and authors

An ADR is not a full design document. An ADR is not a meeting note. One ADR covers one decision. If you have three independent decisions, write three ADRs.

Status changes over time. When a new ADR replaces an old ADR, mark the old ADR as superseded. Do not delete the old ADR. History explains the structure that people see in the code.

Write in the same language as the team. Use facts. Name rejected options in one list. A future reader must see why you did not pick the popular option.

Start ADRs early. Three ADRs on a small project already help: data store, identity, and deploy.

### Questions

#### Theoretical questions

1. What is an Architecture Decision Record?
2. When do you write an ADR?
3. What sections does a typical ADR include?
4. Why do you keep a superseded ADR?
5. Why is one ADR limited to one decision?

#### Easy practical tasks

1. Open [https://adr.github.io/](https://adr.github.io/). Write the purpose of an ADR in four sentences.
2. Write an ADR title and a one-paragraph context for "use PostgreSQL as the source of truth".
3. List five decisions on a student project that need an ADR. List five decisions that do not.
4. Copy an ADR template into your notes. Fill the headings only.

#### Medium practical tasks

1. Write a full ADR for authentication: session cookie versus token. Include two rejected options and three consequences.
2. Write two ADRs that conflict if both are accepted. Then write which ADR wins and why.
3. Review a public repository that contains ADRs. Summarize one ADR in six sentences. Do not copy the full text.

#### Advanced practical tasks

1. Write three ADRs for one side project: database, public API style, and deploy unit. Cross-link the ADRs where a consequence of one ADR is context for another.
2. Design an ADR index with number, title, status, date, and related quality attributes. Fill it for a system with eight decisions.

---

## Reliability, performance, scale, security, maintainability, cost

A quality attribute is a non-functional property of a system. This path uses six attributes as a first set: reliability, performance, scale, security, maintainability, and cost. Later topics add observability and operations targets.

Reliability is the ability of the system to do the correct work for a period of time. Availability is the share of time that the system can accept and complete useful work. A common expression is `availability = uptime / (uptime + downtime)`. "Three nines" means 99.9 percent availability. That is about 8.8 hours of downtime in one year. A system can be available and still return a wrong result. Reliability cares about correct work. Availability cares about the system being ready.

Performance is how fast and how much work the system does. Latency is the time of one operation. Throughput is the number of operations in a period. Use a percentile for latency, not only the average. P99 is a common target. Average latency can hide a bad tail.

Scale (scalability) is the ability to keep the required latency and throughput when load grows. Vertical scale adds resources to one machine. Horizontal scale adds more machines. Measure the load before you split the system.

Security is the protection of data and operations from unauthorized use. Authentication answers who the principal is. Authorization answers what the principal may do. A trust boundary is a line where the level of trust changes. Topic 10 covers security structure.

Maintainability is the ability to change the system at a safe cost. Clear modules, tests, and ADRs improve maintainability. A fashion stack that the team cannot operate decreases maintainability.

Cost is money and people time. Cloud invoices, licenses, and on-call hours are cost. A design that saves money and burns the team is still expensive.

Each attribute needs a measure. Do not promise five nines if the team cannot test failover. Do not say "secure" without a check.

### Questions

#### Theoretical questions

1. How do reliability and availability differ?
2. What is the difference between latency and throughput?
3. What is the difference between vertical scale and horizontal scale?
4. What is the difference between authentication and authorization?
5. Why is people time part of cost?

#### Easy practical tasks

1. Write one measurable sentence for each of the six attributes in this section.
2. Convert 99.0 percent, 99.9 percent, and 99.99 percent to approximate downtime in one year. Show the arithmetic.
3. Make a table: "Attribute" and "Bad empty word". Add six rows (`fast`, `secure`, and similar).
4. List five faults that can decrease availability in a web application.

#### Medium practical tasks

1. For a campus catalog, write a latency SLO and an availability SLO. Include the operation name.
2. Compare two designs for the same shop: one cheap host and one redundant pair. Write eight sentences on reliability versus cost.
3. Write how you would measure maintainability for a twelve-week student project (change time, defect count, or review time).

#### Advanced practical tasks

1. Write a quality attribute scenario for performance and one for maintainability. Use source, stimulus, environment, response, and measure.
2. Find a public status page. Extract one reliability number and one performance number. Mark missing measures.

---

## You cannot maximize all of them

You cannot maximize all quality attributes at the same time. Each strong choice has a cost in another attribute. Architecture work is the selection of the attributes that matter and the attributes that you accept as weaker.

Examples of the conflict:

- More redundancy can improve reliability and can increase cost.
- A strict security control can decrease latency performance.
- Horizontal scale can improve throughput and can decrease maintainability (more operations).
- A cheap single host can decrease cost and can decrease reliability.
- Fast delivery of a feature can decrease maintainability if you skip module boundaries.

Stakeholders disagree on which attribute wins. Record the winner and the owner. An ADR is the right place for that record.

Do not collect patterns as decoration. Use a pattern only when it solves a stated quality problem. "We use microservices" is not a design if you do not state latency, data ownership, and failure behavior.

A useful architecture review asks: What is costly to change? Which quality is at risk? Who must accept the risk?

Write the rejected quality in plain words. "We accept weaker isolation in exchange for one deploy unit" is a complete trade-off. "We want everything" is not a design.

### Questions

#### Theoretical questions

1. Why can you not maximize all quality attributes at the same time?
2. Give one pair of attributes that often conflict.
3. Who must accept a quality that you weaken?
4. Why is "we use microservices" not a complete design?
5. What three questions does a useful architecture review ask?

#### Easy practical tasks

1. Write five sentences about trade-offs. Use only facts from this section.
2. Make a table: "Choice" and "Attribute that can get weaker". Add six rows.
3. List four fashion words. For each word, write one concrete quality or constraint.
4. Write one trade-off sentence for a two-person team that ships in eight weeks.

#### Medium practical tasks

1. Take a teammate design that is only a technology list. Rewrite it as structure plus decisions plus one rejected quality.
2. Write a one-page checklist that a reviewer uses to reject empty architecture slides.
3. For a library loan system, pick two attributes to protect and two attributes to accept as weaker. Give reasons.

#### Advanced practical tasks

1. Collect five public "architecture" slides or blog posts. Score each from 0 to 5 on structure, constraints, and trade-offs. Explain each score.
2. Write a one-page trade-off matrix for a fictional payments tool: six attributes, one winner, one loser, one owner.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do structure, costly decisions, stakeholders, and constraints fit in one architecture definition?
2. Why does this path start with requirements and ADRs before styles such as layers or hexagon?
3. What is the difference between a quality attribute and a constraint in this handbook?
4. How does an empty fashion word hide a missing stakeholder?
5. What makes an accidental structure expensive after one year?

#### Easy practical tasks

1. Write a one-page cheat sheet: architecture definition, requirement types, stakeholder list, constraint types, six attributes, ADR sections.
2. For a to-do application, write two functional requirements, two quality requirements, four stakeholders, four constraints, and one ADR title.
3. Draw a C4 context diagram for that to-do application. Add a note that lists three costly decisions.
4. Bookmark [https://c4model.com/](https://c4model.com/) and [https://adr.github.io/](https://adr.github.io/). Write one sentence on when you open each site.

#### Medium practical tasks

1. Write a short architecture brief (one or two pages) for a campus lost-and-found system. Include requirements, stakeholders, constraints, and two ADR summaries.
2. Make a timeline of five decisions for a twelve-week project. Mark which decisions you record in week 1.
3. Interview one teammate or write a simulated interview script. Turn vague wishes into six testable requirements and one named trade-off.

#### Advanced practical tasks

1. Compare C4 context/container with an informal box diagram for the same system. Write when each view helps a stakeholder.
2. Read one chapter overview from a public architecture book index. Map five book themes to this topic and to later topics in `architecture.topics.md`.
