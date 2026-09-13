# 5. Data Modeling

## Description

This topic shows how you design documents and collections. You choose embed or reference. You model one-to-one, one-to-many, and many-to-many. You keep arrays bounded. You choose a snapshot or a live reference. Complete CRUD and querying first.

Use one term for each concept. **Embed** means you store related data inside the same document. **Reference** means you store an id of a document in another collection. **Bounded** means the array has a maximum size that you control. A **snapshot** copies values at write time. A **live reference** reads the current related document later.

The main rule: store together the data that you access together.

---

## Embed vs reference

**Embed** when the related data belongs to the parent and you read it with the parent. Example: an order and its line items. One `find` returns the order.

**Reference** when the related data has its own life cycle, many parents share it, or the related data is large. Example: a product catalog. Many orders point to the same product id. You update the product in one place.

Embed benefits:

- One read
- One atomic write to the parent document
- A clear nest for the application

Embed costs:

- Duplication if you copy shared data
- Document growth
- The 16 MB limit
- Harder queries that start from the child

Reference benefits:

- One copy of shared data
- Smaller parent documents
- Independent queries on each collection

Reference costs:

- Extra reads or a `$lookup`
- No single-document atomicity across collections (unless you use a transaction)

You can combine both. Store a product id and a snapshot of the price on the order line. That pattern is an extended reference (later topic).

Do not embed a full unbounded child set. Do not reference when every screen needs the child list and the list is small and private to the parent.

### Questions

#### Theoretical questions

1. What does embed mean?
2. What does reference mean?
3. Why is one read a benefit of embed?
4. Why is a shared product a poor full embed in every order?
5. How does the 16 MB limit affect embed?

#### Easy practical tasks

1. Write one embed example and one reference example in four sentences total.
2. Draw two diagrams of a blog post and comments: embed vs reference.
3. List three questions that are easier if comments are a separate collection.
4. Open the official data-modeling introduction. Write the URL and the access-together rule in your own words.

#### Medium practical tasks

1. Model `orders` and `products` for a shop. Write which you embed and which you reference. Give one reason each.
2. Count reads: "show order with product names" for embed of names vs reference only. Write the read count.
3. Take a one-to-many pair from a relational schema. Write the MongoDB embed form and the reference form.

#### Advanced practical tasks

1. Read the official "Embedded Data vs References" page. Write three decision questions from that page.
2. Design a fallback: embed the last 20 comments, reference the rest. Write the two collections and the read path.

---

## One-to-one, one-to-many, many-to-many in documents

**One-to-one.** Embed if the side is small and private. Example: a user and a profile settings object. Use a second collection if the side is large or secret and you rarely read it.

**One-to-many.** Embed the many side when it is bounded and you always load it with the parent. Example: an order and lines. Reference the many side when it grows without a bound. Example: a user and all orders. Store `userId` on the order.

**Many-to-many.** Use references. Example: students and courses. Store an array of course ids on the student if the array is bounded. Store a mapping collection if the link has attributes (grade, term) or if both sides have large arrays.

Patterns:

- Parent with embedded array of children
- Child with a parent id (the many side owns the reference)
- Two-way arrays of ids (only if both arrays stay small)
- A link collection `{ studentId, courseId, grade }`

Do not store two huge id arrays that grow forever. That is two unbounded arrays.

Cardinality is an estimate. Write the expected maximum. If you cannot write a maximum, use a reference or a link collection.

### Questions

#### Theoretical questions

1. When do you embed a one-to-one side?
2. Where do you put the reference in a one-to-many design when the many side is unbounded?
3. Why is a link collection useful for many-to-many?
4. What is the risk of two-way id arrays?
5. Why must you write a maximum size for an embedded array?

#### Easy practical tasks

1. Model user + settings as one-to-one embed. Write one sample document.
2. Model user + orders as one-to-many with `userId` on orders. Write two sample documents.
3. Model students and courses as a link collection. Write three sample documents.
4. Make a table: relationship type, embed or reference, one reason.

#### Medium practical tasks

1. For comments on a post, write the embed design and the `postId` design. State a number of comments that changes your choice.
2. Add a grade to a student-course link. Show why the link collection holds that field better than two id arrays.
3. Draw cardinality notes: "1 user to N orders, N ~ 1000 per year". Pick a design.

#### Advanced practical tasks

1. Redesign a relational many-to-many (three tables) as MongoDB collections. Write the queries for "courses for student" and "students for course".
2. Read about the subset pattern (preview). Apply it to a user with 10 000 followers. Write what you embed.

---

## The "data that is accessed together is stored together" rule

This rule is the center of MongoDB modeling. If a screen or an API always reads A with B, store A and B in one document when size allows.

The rule is about **access**, not about real-world ownership only. A customer "owns" orders in the business. The application may still store orders in a separate collection because the list is unbounded.

Work from queries:

1. List the reads and writes.
2. For each read, list the fields.
3. Group fields that always appear together.
4. Check size and growth.
5. Split when growth or sharing requires a split.

Do not copy a relational schema table-for-table without this step. A third normal form schema is a starting map of entities. It is not the final collection list.

Do not embed data that you access together if that embed duplicates a huge shared object on every write. Then store a reference and accept an extra read, or store a small snapshot.

The rule also applies to writes. If two fields must change together, one document gives you single-document atomicity.

### Questions

#### Theoretical questions

1. What does "accessed together" mean in this rule?
2. Why can business ownership differ from the collection design?
3. Why is a list of queries the first modeling input?
4. How does single-document atomicity relate to this rule?
5. Why is a table-for-table copy a weak MongoDB model?

#### Easy practical tasks

1. Write three API reads for a shop. Mark which entities each read needs.
2. For "get order by id", write one document that satisfies the read.
3. For "list orders for user", write why orders are not all embedded in the user.
4. Find the official sentence of this rule. Write it in your own words. Include the URL.

#### Medium practical tasks

1. Take five queries on a blog. Design collections so that four queries need one find. Write the query that still needs two finds.
2. Show a case where access-together conflicts with sharing. Write your compromise (snapshot or extra read).
3. Rewrite a 3NF schema of four tables into two collections. Justify each merge with a query.

#### Advanced practical tasks

1. Instrument a small app (or imagine one): log which collections each endpoint reads. Propose one embed that removes a round trip.
2. Write a one-page modeling worksheet: queries, size, growth, share, atomicity. Fill it for a chat room + messages.

---

## Bounded arrays vs unbounded arrays

A **bounded array** has a maximum number of elements that you enforce. Example: last 50 notifications, up to 10 addresses, up to 100 order lines for a small shop.

An **unbounded array** grows without a plan. Example: all events for a device in one document, all comments for a popular post, all orders for a user.

Unbounded arrays cause:

- Document growth toward 16 MB
- Slow updates (the server rewrites a larger document)
- Large reads when you need only a few elements
- Painful indexes (multikey)

Use a bounded array when the bound is real. Enforce the bound in the application. Use `$push` with `$slice` to keep the last N items.

Use a separate collection when the list has no bound. Store a parent id on each child. Index that id.

A "large but finite" list can still be too large to embed. A bound of 100 000 is not a practical embed. Write the byte size, not only the count.

Do not use an array as a dump for every future child.

### Questions

#### Theoretical questions

1. What is a bounded array?
2. Why is an unbounded array a problem for the 16 MB limit?
3. How does `$push` with `$slice` keep a bound?
4. Where do you store children when the array is unbounded?
5. Why is a large finite bound still a bad embed?

#### Easy practical tasks

1. Write three bounded array examples and three unbounded examples.
2. Insert a document. `$push` a tag. Repeat. Write the array.
3. Use `$push` and `$slice` to keep the last 3 items. Show the array after 5 pushes.
4. Estimate BSON size of 1000 elements of `{ n: 1 }`. Use `Object.bsonsize`.

#### Medium practical tasks

1. Model notifications: last 20 in the user document, full history in `notifications`. Write both shapes.
2. Compare update time for `$push` on a 10-element array vs a 5000-element array (measure on your machine). Write the two times.
3. Write an application check that rejects a sixth address if the max is 5.

#### Advanced practical tasks

1. Read the anti-pattern page about unbounded arrays. Write two official reasons. Include the URL.
2. Design a migration from an embedded comment array to a `comments` collection. Write the steps and the read compatibility plan.

---

## Snapshot of related data vs live reference

A **snapshot** copies fields from a related entity at write time. Example: store `priceAtOrder` and `productName` on the order line. Later product changes do not change the order.

A **live reference** stores only an id. Each read loads the current related document. Example: store `productId` and read the current name for a catalog page.

Use a snapshot when history must not change. Money, legal text, and invoices need snapshots.

Use a live reference when the current value is the truth. A product page must show the current name and stock.

You can store both: `productId` plus snapshot fields. The id lets you join. The snapshot lets you display the past.

Do not snapshot huge objects. Copy the fields that the parent read needs. That is a subset (later pattern).

Do not use a live reference for a completed order price. The customer must see the price that they paid.

### Questions

#### Theoretical questions

1. What is a snapshot in this topic?
2. What is a live reference?
3. Why must an invoice line snapshot the price?
4. Why must a catalog page use a live name?
5. Why store both an id and snapshot fields?

#### Easy practical tasks

1. Write an order line with `productId`, `nameSnapshot`, and `priceSnapshot`.
2. Change the product name in a `products` collection. Write which document still has the old name.
3. List three fields that you would snapshot and three that you would always read live.
4. Make a two-column table: "Snapshot" and "Live". Add four rows.

#### Medium practical tasks

1. Design a user profile shown on comments. Choose snapshot of `displayName` vs live read. Write a privacy or rename consequence for each.
2. Write the update path when a live field changes and 1 million parents reference it. Write why you do not update all parents.
3. Model a "quote" that becomes an "order". Write which fields freeze at convert time.

#### Advanced practical tasks

1. Read the extended-reference pattern (preview). Write how it combines snapshot and reference.
2. Design a correction flow: the snapshot was wrong. Write whether you edit the order, add an adjustment document, or both.

---

## Polymorphic collections (when and when not)

A **polymorphic collection** stores documents that share a topic but not the same fields. Example: an `events` collection with `{ type: "login", ip: "..." }` and `{ type: "purchase", sku: "..." }`.

Use one collection when:

- You query the documents together (all events for a user)
- The documents share a few indexed fields (`userId`, `createdAt`, `type`)
- The variant fields stay small

Do not use one collection when:

- The shapes have no shared query
- Indexes become many and sparse
- Validators cannot describe the mix without a weak schema
- Teams own the variants as separate products

An alternative is one collection per type. That split helps when each type has its own indexes and life cycle.

A `type` or `_t` field names the variant. The application reads `type` before it reads variant fields. JSON Schema can use `oneOf` (later topic).

Polymorphic does not mean "any JSON". You still list the allowed types.

### Questions

#### Theoretical questions

1. What is a polymorphic collection?
2. Why do shared query fields matter?
3. When do separate collections beat one polymorphic collection?
4. What is the role of a `type` field?
5. Why is "any JSON" not a schema?

#### Easy practical tasks

1. Insert two event types into `events`. Find all events for one `userId`.
2. Find only `type: "login"`. Write the filter.
3. Write three allowed event types and their required fields.
4. List two indexes that all events can share.

#### Medium practical tasks

1. Compare one `events` collection vs `logins` and `purchases`. Write two queries that are easier in each design.
2. Sketch a validator idea: `type` plus required fields per type. Use a short `oneOf` outline (no need to run it yet).
3. Count indexes: five types with two extra fields each, all in one collection. Write the index pressure in words.

#### Advanced practical tasks

1. Read official guidance on polymorphic patterns. Write when MongoDB documents recommend it. Include the URL.
2. Design a change stream on `events` that routes by `type`. Write why one collection helps that stream.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do embed, reference, bounds, and snapshot work together for an order with products and a user?
2. Which relationship types usually embed, and which types usually reference? Give one exception each.
3. Why does "accessed together" beat "looks like a table" as a design rule?
4. How do unbounded arrays and polymorphic dumps fail the same size-and-index test?
5. A teammate embeds all orders in `users`. Which facts do you use to reject that design?

#### Easy practical tasks

1. Write sample documents for `users`, `orders`, and `products` that follow this topic. Keep arrays bounded.
2. Write a cheat sheet: embed vs ref, 1-1 / 1-N / N-N, access-together, bounded array, snapshot vs live, polymorphic.
3. Mark each array in your sample as bounded or unbounded. Write the max size.
4. Draw one diagram of your shop model with arrows for references and boxes for embeds.

#### Medium practical tasks

1. Take a relational schema of six tables. Produce a MongoDB model. Write a query list and which find covers each query.
2. Pick one embed in your model and write the condition that would force a split to a collection.
3. Write two event types and one order type. Decide one or two collections. Justify with queries.

#### Advanced practical tasks

1. Write a one-page design review of your model: growth, 16 MB, atomicity, duplication, and the hardest query.
2. Model a social "follow" graph. Compare embed of id arrays vs a `follows` collection. Estimate size at 1 million follows.
