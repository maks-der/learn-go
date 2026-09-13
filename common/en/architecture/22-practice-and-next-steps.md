# 22. Practice and Next Steps

## Description

Practice turns the path into skill. This topic covers Architecture Decision Records on a real side project, C4 context and container diagrams, the rule that you split a modular monolith only after a real pain, a weekly post-mortem habit, and three books that you read after you ship a monolith.

The file `architecture.topics.md` also lists a suggested practice order of ten steps. Use that list after you finish these sections. Complete Topics 1 to 21 before you treat this topic as done. Pair this path with `db.topics.md`, `net.topics.md`, and a language path. Architecture remains trade-offs (Topic 1).

Use one term for each concept. An ADR records one costly decision. A C4 context diagram shows the system and its neighbors. A C4 container diagram shows deployable units and stores. A post-mortem is a learning document after an incident. A book is a next step, not a substitute for a running system.

---

## Write ADRs for a real side project

Pick a system that you own or that you may change. A to-do tool, a campus club site, or a library of notes is enough. The project must run. A slide deck is not a side project.

Write ADRs in the repository (Topic 1). See [https://adr.github.io/](https://adr.github.io/). Start with three:

1. Source of truth (the primary store)
2. Identity (how a principal authenticates)
3. Deploy unit (one process or more)

Add ADRs when a decision is costly to reverse: public API style, a queue, a region, a secret store (Topics 8, 10, 14, and 19).

Rules from this path:

- One decision per ADR
- Context includes constraints (team, time, law, existing systems)
- Consequences include the quality that you weaken
- Status can change. Do not delete a superseded ADR (Topic 18)

Do not write twenty ADRs on day one. Three good ADRs beat a folder of slogans.

Review the ADRs after you ship. If the code disagrees, fix the code or supersede the ADR.

Use the same language as the code. Name rejected options. A future you is the first reader.

If the project is a class assignment, still write the three ADRs. They train the habit that Topic 1 started.

### Questions

#### Theoretical questions

1. Why must the side project run?
2. Which three ADRs does this section name first?
3. When do you add a fourth ADR?
4. What must context include?
5. What do you do when code and an accepted ADR disagree?

#### Easy practical tasks

1. Open [https://adr.github.io/](https://adr.github.io/). Write the purpose of an ADR in four sentences.
2. List five decisions on your project that need an ADR. List five that do not.
3. Write an ADR title and a one-paragraph context for the primary store.
4. Copy an ADR template into the project. Fill headings only, then fill one full ADR.

#### Medium practical tasks

1. Write the three starter ADRs. Cross-link a consequence of the store ADR to the deploy ADR.
2. Write two ADRs that would conflict if both were accepted. Record which one wins and why.
3. Review a public repository that contains ADRs. Summarize one ADR in six sentences. Do not copy the full text.

#### Advanced practical tasks

1. Write an ADR index for eight decisions. Include status, date, and related quality attributes (Topic 2).
2. After two weeks of coding, supersede one ADR with a new ADR. Keep the old file.

---

## Draw a C4 context and container diagram

The C4 model is a small set of diagram levels. The site is [https://c4model.com/](https://c4model.com/). This path uses context and container first.

A context diagram shows people, your system as one box, and external systems. Each relation has a verb. A container diagram shows the deployable units and stores inside your system: web application, API, worker, database, and broker.

Rules:

- One job per diagram
- Names match the code and the ADRs
- Trust boundaries can overlay later (Topic 14)
- Do not fill the page with framework logos

Draw by hand or with a simple tool. The value is the labels, not the tool.

Update the diagrams when a container appears or dies. A stale diagram is a lie (Topic 5, mud symptoms).

Do not start with a code (class) diagram. C4 code level is optional and late.

If you extracted a service (Topic 12), the container view must show two units and two owners of data. If you did not extract, do not draw fake services.

Show one external identity provider if you use one (Topic 14). Show the file drop if you integrate by file (Topic 13).

Keep a short text next to the picture: three costly decisions and the SLO names (Topics 1 and 15).

### Questions

#### Theoretical questions

1. What does a C4 context diagram show?
2. What does a C4 container diagram show?
3. Why must names match the code?
4. Why do you avoid a class diagram at the start?
5. What must the container view show after a real extract?

#### Easy practical tasks

1. Open [https://c4model.com/](https://c4model.com/). Write the four C4 levels in order.
2. Draw a context diagram for your side project. Use three people or systems at least.
3. Draw a container diagram for the same project. Label each store.
4. Write one verb on every arrow.

#### Medium practical tasks

1. Add a note that lists three costly decisions next to the container view.
2. Compare your diagram with the repository folders. Write every mismatch.
3. Overlay two trust boundaries on the container view (Topic 14). Stay defensive.

#### Advanced practical tasks

1. Draw context and container for a library loan system that is not your project. Write five sentences on what each view hides.
2. Review a public architecture picture. Redraw it as C4 context and container. Mark empty fashion words that you drop.

---

## Split a modular monolith only after a real pain

The default remains a modular monolith (Topic 5). A split is an extract of one module to a new deployable unit with its own data when ownership requires it (Topics 9 and 12).

Real pain means a measure or a hard constraint:

- Independent scale that you measured (Topics 2 and 17)
- An independent release clock that a window cannot solve (Topic 18)
- Isolation that law or a stakeholder requires (Topics 1 and 14)

Not real pain:

- A fashion word
- A resume list
- A slide from a large company
- A folder count

Finish the extract path in Topic 5 before the network hop: clear interface, table owners, observability, timeouts, rollback. If the pain dies at step 3, stop. That stop is a success.

One extract is a semester of operations learning. Do not extract five units in one burst.

Write the success measure in the extract ADR. If the measure does not move, roll back.

The suggested practice order in `architecture.topics.md` places a split last. Follow that order.

If you never split, you can still finish this path. A healthy monolith is a complete architecture.

### Questions

#### Theoretical questions

1. What is a real pain in this section?
2. Name three items that are not real pain.
3. Why can a stop after better modules be a success?
4. Why is one extract enough for a long time?
5. Why can you finish this path with no split?

#### Easy practical tasks

1. Write a one-page "do not split" checklist with ten items.
2. Make a table: "Reason" and "Pain? (yes/no)". Add resume, measured CPU ceiling, and fashion word.
3. Number the Topic 5 extract steps in one sentence each.
4. Write four sentences that a reviewer says to a week-two eight-service plan.

#### Medium practical tasks

1. On your side project, name the module that would extract first if pain appeared. List the missing rules today.
2. Write an extract ADR that stays in `proposed` until a measure exists.
3. Role-play on paper: reject a split of every table into a service. Cite Topics 5, 7, 9, and 12.

#### Advanced practical tasks

1. Write a two-page gate: evidence before the first network hop. Include SLO, ownership, and rollback.
2. If you already extracted too early, write a merge-back plan. A merge can be the correct evolution (Topic 18).

---

## Read one post-mortem a week

A post-mortem (or incident review) is a document that states what happened, what the system did, and what you will change. Good reviews avoid blame of a person. They name missing tests, missing SLOs, and missing boundaries.

Read one public post-mortem each week. Sources include company engineering blogs and public status-page write-ups. Prefer texts that include a timeline and a repair.

How to read:

1. Write the user impact in one sentence
2. Write the first failed assumption (often a fallacy from Topic 11)
3. Name one quality attribute that moved (Topic 2)
4. Name one architecture decision that made the incident larger or smaller
5. Write one habit that you will add to your project

Ten minutes is enough. A notebook of twelve reviews is a course.

Do not treat a marketing "we learned and we are stronger" page as a post-mortem. Look for facts: times, SLIs, and changes.

Stay defensive. Use incidents to add controls. Do not use write-ups to learn attack steps.

Share one review with a teammate each month. Speaking the story fixes the lesson.

If you cause an incident in a lab that you own, write a post-mortem. That document is an ADR cousin.

### Questions

#### Theoretical questions

1. What is a post-mortem in this handbook?
2. Why do good reviews avoid blame of a person?
3. What five notes do you write after a read?
4. What makes a marketing page a poor source?
5. How does a lab post-mortem relate to an ADR?

#### Easy practical tasks

1. Find one public incident write-up. Write the user impact in one sentence.
2. Make a table: "Week" and "URL". Reserve four empty rows for the next month.
3. List four facts that you want in every review (time, SLI, decision, repair).
4. Write four sentences on why blame of a person hides a missing test.

#### Medium practical tasks

1. Complete the five-note method for one public review. Do not copy long passages.
2. Map the incident to one topic in this path (timeout, cache, deploy, or secret).
3. Write a one-page template that you will fill after your own lab outage.

#### Advanced practical tasks

1. Compare two reviews of similar failures (for example, a bad deploy). Write which review teaches more and why.
2. Build a twelve-week reading list. Tag each item with a topic number from `architecture.topics.md`.

---

## Books: DDIA, Fundamentals of Software Architecture, Building Microservices (after you built a monolith)

Books deepen the path. They do not replace a running monolith. Read them after you have modules, one database, logs, and a few ADRs.

*Designing Data-Intensive Applications* (Kleppmann) is a data-systems book. The site is [https://dataintensive.net/](https://dataintensive.net/). Use it to deepen Topics 9 to 11 and 17. Read after you can explain source of truth, replication at a high level, and why a network call is not a local call.

*Fundamentals of Software Architecture* (Richards and Ford) is an architecture-practice book. Use it to deepen quality attributes, styles, and team topology ideas. Read after Topics 1 to 6.

*Building Microservices* (Newman) is a service-split book. Read it only after you built and operated a monolith. The book will otherwise look like a shopping list. Pair it with Topic 12 and with the "when not to" section.

Also keep:

- [https://c4model.com/](https://c4model.com/)
- [https://adr.github.io/](https://adr.github.io/)
- [https://12factor.net/](https://12factor.net/)
- [https://www.enterpriseintegrationpatterns.com/](https://www.enterpriseintegrationpatterns.com/)
- [https://sre.google/books/](https://sre.google/books/)

How to read a book in this path:

- Map each chapter to a topic number
- Write five sentences per chapter in your own words
- Reject one idea that your constraints forbid
- Do not copy long passages

If you only have time for one book, pick DDIA after a small system, or *Fundamentals* if data systems feel too heavy. Do not start with microservices.

Return to `architecture.topics.md` and run the ten practice steps. Then pick a language path and `db.topics.md` for depth.

### Questions

#### Theoretical questions

1. Why must a monolith come before *Building Microservices*?
2. Which topics does DDIA deepen?
3. Which topics does *Fundamentals of Software Architecture* deepen?
4. What four notes do you write per chapter?
5. Which book do you pick if you can read only one and you already have a small system?

#### Easy practical tasks

1. Open [https://dataintensive.net/](https://dataintensive.net/). Write the book title and the author name.
2. Make a table: "Book" and "Read after". Add the three books.
3. Bookmark the five official sites in this section. Write one sentence on when you open each.
4. Write a twelve-week reading plan with one book and weekly post-mortems.

#### Medium practical tasks

1. Read one chapter overview or table of contents of DDIA. Map five headings to topic numbers. Do not copy long text.
2. Write a one-page "what I will not adopt yet" list after you skim *Building Microservices* contents.
3. Compare this handbook path with the *Fundamentals* part titles. Write six matches and two gaps.

#### Advanced practical tasks

1. Write a reading seminar plan (six sessions) for a campus club. Include ADRs and C4 homework.
2. After you finish one book, write a one-page errata of the book against your constraints (team, law, cost). Stay factual.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do ADRs, C4 views, and a running project form one evidence pack of your architecture?
2. Why does the path refuse a split until pain, and how does a weekly post-mortem teach that refusal?
3. How do the three books divide data, practice, and services without replacing Topics 1 to 21?
4. What does the suggested practice order in `architecture.topics.md` add that this topic does not repeat section by section?
5. How will you know that you are ready to leave the beginner path?

#### Easy practical tasks

1. Write a one-page cheat sheet: three ADRs, two C4 views, split rule, post-mortem method, and book order.
2. For your side project, write the next seven days: one ADR, one diagram update, and one post-mortem note.
3. Copy the ten-step practice order from `architecture.topics.md` into your notes. Mark what you already did.
4. Write three sentences that define "ready for Topic 12 extract" in measures, not in wishes.

#### Medium practical tasks

1. Write a two-page personal architecture brief for the next semester: project, ADRs, diagrams, SLO, and a rejected fashion word.
2. Take a teammate plan that is only a technology list. Rewrite it as structure, decisions, and a practice calendar.
3. Map your brief to `db.topics.md` and `net.topics.md`. Write two skills that you still lack.

#### Advanced practical tasks

1. Run a 90-minute self-review: diagrams, ADRs, SLO, and one failure drill on a system that you own. Write the gaps.
2. After twelve post-mortems and one book, write a one-page "what I changed in my project". Point to commits or ADRs, not to slogans.
