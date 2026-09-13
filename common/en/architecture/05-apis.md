# 5. APIs

## Description

An API (application programming interface) is a contract that a client can call. This topic covers REST with HTTP and JSON, RPC with gRPC, GraphQL, versioning, backward compatibility, pagination, errors, idempotency keys, and the difference between public and internal APIs.

An API is architecture. A change of the contract is costly for every client. Treat the API as a product. Write the contract. Test the contract.

Complete Topics 1 to 4 before this topic. Pair with `net.topics.md` if HTTP is new. Use defensive design. Do not use this topic to plan abuse of other systems.

Use one term for each concept. REST is not "any JSON over HTTP". RPC is a call style. GraphQL is a query language and a runtime. Those tools solve different problems.

---

## REST / HTTP JSON

REST (Representational State Transfer) is a style for HTTP APIs. A practical REST API uses resources, HTTP methods, and representations (often JSON).

A resource is a noun that the API exposes (`/books/42`, `/loans`). `GET` reads. `POST` creates or starts an action. `PUT` replaces. `PATCH` updates a part. `DELETE` removes. Use methods for their meaning. Do not use `GET` for a change.

JSON is a text format for structured data. Keep field names stable. Use types that JSON can express (object, array, string, number, boolean, null). Dates need a stated format (ISO 8601 is common).

A practical REST checklist:

- Use nouns in paths for resources.
- Use HTTP status codes with a clear meaning (`200`, `201`, `204`, `400`, `401`, `403`, `404`, `409`, `429`, `500`).
- Use `Content-Type: application/json` when the body is JSON.
- Keep authentication out of the URL (use headers).
- Keep the API stateless at the instance (Topic 2).

REST is not a religion. A `POST /loans/42/return` action can be clearer than a generic `PATCH`. Document the rule.

Hypermedia (links in responses) is optional. Many teams ship resource JSON without a full hypermedia design. That choice is acceptable if the contract is documented.

Validate input. Bound sizes. Set timeouts (Topic 2). Log a request identifier.

### Questions

#### Theoretical questions

1. What is a resource in a practical REST API?
2. Why must `GET` not change data?
3. What does a status code communicate?
4. Why must secrets stay out of the URL?
5. When can an action path be better than a generic `PATCH`?

#### Easy practical tasks

1. Design paths and methods for books and loans (list, get, create, return).
2. Make a table: "Status" and "When to use". Add 200, 201, 400, 404, and 409.
3. Write a JSON example for `GET /books/42` (five fields).
4. List four headers that a typical JSON API uses.

#### Medium practical tasks

1. Design `POST /loans` and the error when the book is already on loan. Include status and body shape.
2. Compare `PUT /books/42` with `PATCH /books/42`. Write eight sentences.
3. Write ten validation rules for `POST /books`. Include size limits.

#### Advanced practical tasks

1. Write a one-page REST style guide for a team: paths, methods, dates, and identifiers.
2. Review a public HTTP JSON API. Mark where it follows this section and where it does not. Do not copy large responses.

---

## RPC (gRPC) and GraphQL (when each helps)

RPC (remote procedure call) is a style where the client calls a named procedure on the server. The client looks like a local function. The network is hidden behind a stub.

gRPC is a common RPC system. It uses HTTP/2 and Protocol Buffers (protobuf) for types. You write a `.proto` file. Tools generate client and server code.

gRPC fits internal service-to-service calls when both sides can share the schema and when you want strict types and streaming. gRPC is harder for a simple browser client. Browsers often use gRPC-Web or a JSON HTTP API instead.

Benefits of gRPC:

- Generated types reduce some contract drift.
- Binary encoding can be efficient.
- Streaming is a first-class idea.

Costs of gRPC:

- Browser and cache infrastructure is simpler with JSON HTTP.
- Human debugging of binary payloads needs tools.
- You still need version rules. Generated code does not remove compatibility work.

Treat the `.proto` file as the contract. Review it like an API. Do not change field numbers casually. Protobuf has encoding rules that depend on those numbers.

GraphQL is a query language and a runtime. The client asks for a graph of fields in one request. The server resolves those fields. GraphQL helps when many clients need different shapes of the same graph and when a single HTTP JSON resource would force many round trips.

GraphQL costs:

- A naive resolver can create an N+1 load on the store.
- Authorization must run per field or per type, not only at the URL.
- Caching is harder than caching a `GET` resource.
- A public GraphQL endpoint can expose a large surface.

RPC is not "faster, so always use it". GraphQL is not "modern, so always use it". Measure. A small monolith does not need gRPC inside the process. Use a function call. A single-page student app that reads three resources can use REST.

A typical split: JSON HTTP (REST style) for public or browser clients. gRPC only if an internal second service appears and both teams can operate protobuf. GraphQL when a BFF or a public graph has a proven payload problem (Topic 13).

### Questions

#### Theoretical questions

1. What is RPC?
2. When does gRPC fit well?
3. Why is gRPC harder for a simple browser?
4. When does GraphQL help?
5. Why can a public GraphQL endpoint increase security work?

#### Easy practical tasks

1. Write four sentences that compare REST JSON with gRPC.
2. Write four sentences that compare REST JSON with GraphQL.
3. Make a table: "Client" and "REST, gRPC, or GraphQL?". Add browser, mobile app, and internal worker.
4. List three costs of binary payloads for a beginner team.

#### Medium practical tasks

1. Sketch a `.proto` service with two RPCs (`GetBook`, `BorrowBook`). Describe fields in words.
2. Write an ADR: JSON HTTP for public clients, gRPC only if a second internal service appears.
3. Design one GraphQL query for a loan card. Write two REST calls that the same page would need without GraphQL.

#### Advanced practical tasks

1. Write a one-page comparison of error models: HTTP status plus JSON error versus gRPC status codes.
2. Write a one-page note on GraphQL authorization and N+1. Stay defensive. Do not plan abuse.

---

## Versioning and backward compatibility

Versioning is the practice of changing a contract without a silent break of old clients. Backward compatibility means an old client still works after you deploy a new server (or the reverse, in a stated window).

Rules that keep a JSON HTTP API compatible:

- Add a field. Do not remove a field that clients still read.
- Do not change the meaning of a field.
- Do not change a type (string to object) without a new field name or a new version.
- Treat unknown fields in a request as a validation policy that you write (reject or ignore). Write the rule.

A compatibility window is a time when old and new clients both work. Mobile clients force a long window. Internal clients in one repository can use a short window if you ship them together.

Version in the URL (`/v1/books`) or in a header. URL versions are easy to see. Header versions keep paths short. Pick one rule and keep it.

A breaking change needs a new version or a new field. Document the sunset date of the old version. Do not keep forever versions without an owner.

Protobuf and gRPC have their own rules. New optional fields can be compatible. Reuse of a field number is dangerous. Reserved numbers exist for a reason.

Compatibility is not only types. Error codes, pagination tokens, and default sort order are part of the contract. A change of default sort is a break for a client that did not send a sort.

Topic 12 covers expand-contract for schema and for APIs. The same idea applies: add, migrate clients, then remove.

### Questions

#### Theoretical questions

1. What is backward compatibility?
2. What is a compatibility window?
3. Name three JSON changes that break old clients.
4. Why do mobile clients force a long window?
5. Why is a change of default sort a break?

#### Easy practical tasks

1. Write five sentences that define versioning. Use only facts from this section.
2. Make a table: "Change" and "Compatible? (yes/no)". Add six rows.
3. List four places a version can appear (URL, header, protobuf, event name).
4. Write a sunset notice of six sentences for `/v1` when `/v2` exists.

#### Medium practical tasks

1. Design an add-then-remove plan for a rename of `title` to `bookTitle`. Include two deploys.
2. Write an ADR: URL version versus header version for a campus API.
3. Explain in eight sentences how a rolling deploy requires old and new servers to accept the same requests.

#### Advanced practical tasks

1. Write a one-page compatibility standard: add, deprecate, remove, and who owns the calendar.
2. Design a compatibility window for a `.proto` change that adds a field. Include generated clients.

---

## Pagination, errors, idempotency keys

Pagination splits a large list into pages. A client must not request an unbounded list. Bound the page size. State the maximum.

Two common styles:

- Offset and limit: simple. Unstable if rows insert during a crawl.
- Cursor (a token for the next page): more stable for a live list. The token is opaque to the client.

Write the default page size. Write the sort. A missing sort makes pages meaningless. Return a next token or a next URL. Do not force the client to guess.

Errors are part of the contract. A useful error body includes:

- An application error code that stays stable
- A short message for operators (no secrets)
- A request identifier (Topic 11)
- Optional field errors for validation

Use HTTP status for the class of failure. Use the body for the detail. Do not return `200` with a hidden error flag as the only signal.

Idempotency keys protect unsafe writes (Topic 2). A client sends a key on `POST` that creates a payment or a loan. The server stores the key and the result. A retry with the same key does not create a second side effect.

Rules for keys:

- The client generates the key for one intent.
- The server bounds key size and lifetime.
- A second request with the same key and a different body is a conflict (`409`).
- Read APIs do not need a key.

Filter and sort are also contract. Document allowed fields. Reject unknown sort keys. Bound filter complexity so that a client cannot force a full scan as a denial of service. Stay defensive.

### Questions

#### Theoretical questions

1. Why must a list API bound page size?
2. How does a cursor page differ from offset and limit?
3. What belongs in a useful error body?
4. What does an idempotency key protect?
5. Why is `200` with a hidden error flag a poor contract?

#### Easy practical tasks

1. Design `GET /books` with `limit`, `cursor`, and a default sort.
2. Make a table: "Status" and "Error body fields". Add 400, 404, 409, and 429.
3. Write five sentences about idempotency keys. Use only facts from this section.
4. List four validation errors for `POST /loans` and the status for each.

#### Medium practical tasks

1. Design `POST /payments` with an `Idempotency-Key` header. Write store contents and replay behavior.
2. Compare offset pagination with cursor pagination for a live loan list. Write eight sentences.
3. Write ten rules for filters on `GET /books` (allowed fields, max length, max page).

#### Advanced practical tasks

1. Write a one-page error catalog: codes, status, and whether the client may retry.
2. Design a pagination token that does not leak internal row numbers. Write what happens when the token expires.

---

## Public vs internal APIs

A public API has clients that you do not ship in the same release. Browsers, mobile applications, partners, and other organizations are public clients. The compatibility window is long. Documentation is a product. Authentication, rate limits, and a stable error model are required.

An internal API has clients that one organization owns. The window can be shorter if you ship clients together. You still write a contract. You still version breaking changes. "Internal" is not a license for a silent break.

Differences that matter:

- Public APIs hide internal names and table shapes.
- Public APIs expose less data (least privilege, Topic 10).
- Internal APIs can use gRPC if both sides operate it.
- Public APIs usually stay on HTTP JSON for browsers and partners.
- Rate limits and abuse controls are stricter on public surfaces. Stay defensive.

A modular monolith can expose one public HTTP API and keep internal packages as function calls. That is not two APIs. Do not add a network inside the process for fashion.

When you have two deployable units, the internal contract is still an API (Topic 9). Treat it with timeouts, idempotency, and compatibility.

Do not put admin operations on the same public surface without a stricter identity and a separate audit (Topic 10). A hidden path is not a security control.

Write who may call the API. Write the data class of each field. A public catalog field is not the same as a personal identifier.

### Questions

#### Theoretical questions

1. What is a public API in this handbook?
2. What is an internal API?
3. Why is "internal" not a license for a silent break?
4. Why do public APIs hide table shapes?
5. When is a function call better than an internal HTTP API?

#### Easy practical tasks

1. Write five sentences that compare public and internal APIs.
2. Make a table: "Client" and "Public or internal". Add six rows.
3. List four extra controls that a public API needs (auth, rate limit, and similar).
4. Write four fields that a public book resource must not include (internal keys, secrets).

#### Medium practical tasks

1. Design a public `GET /books` and an internal `GetBookStock` RPC. Write why the payloads differ.
2. Write an ADR: one public JSON API for a campus app, no partner API in year one.
3. Explain in eight sentences how a mobile public client forces a longer compatibility window than a same-repo worker.

#### Advanced practical tasks

1. Write a one-page API product sheet: audience, SLA or SLO, auth, version, and contact.
2. Review a public API documentation site. Mark public versus implied internal details. Do not copy large text.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you choose REST JSON, gRPC, or GraphQL for one client type?
2. Why is an API a costly decision in the sense of Topic 1?
3. How do versioning, pagination tokens, and error codes form one contract?
4. When does an idempotency key matter more on a public API than on an internal function call?
5. What must stay the same when you later extract a module behind the same public paths?

#### Easy practical tasks

1. Write a one-page cheat sheet: REST checklist, gRPC/GraphQL when, compatibility rules, pagination, errors, keys, public versus internal.
2. For a to-do API, write five paths, two error bodies, one page rule, and one idempotent create.
3. Draw a browser as a public client and a worker as an internal client. Label protocols.
4. Bookmark `net.topics.md`. Write three HTTP facts that this topic assumes.

#### Medium practical tasks

1. Write a short API brief for a campus lost-and-found: public resources, version, errors, and a rejected GraphQL plan.
2. Take a teammate design that uses `GET` to create a record. Rewrite the contract.
3. Write a twelve-week plan: weeks for a public JSON API, and the measure that would reopen gRPC.

#### Advanced practical tasks

1. Write a contract test list (names only, no solutions) that protects compatibility during a rolling deploy.
2. Compare two public APIs in a one-page table: versioning, pagination, error shape, and auth placement. Stay defensive.
