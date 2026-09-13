# 8. APIs

## Description

An API (application programming interface) is a contract that a client can call. This topic covers REST with HTTP and JSON, RPC with gRPC, GraphQL, versioning, backward compatibility, pagination, filtering, errors, idempotency keys, and the difference between public and internal APIs.

An API is architecture. A change of the contract is costly for every client. Treat the API as a product. Write the contract. Test the contract.

Complete Topics 1 to 7 before this topic. Pair with `net.topics.md` if HTTP is new. Use defensive design. Do not use this topic to plan abuse of other systems.

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
- Keep the API stateless at the instance (Topic 4).

REST is not a religion. A `POST /loans/42/return` action can be clearer than a generic `PATCH`. Document the rule.

Hypermedia (links in responses) is optional. Many teams ship resource JSON without a full hypermedia design. That choice is acceptable if the contract is documented.

Validate input. Bound sizes. Set timeouts (Topic 4). Log a request identifier.

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

## RPC (gRPC)

RPC (remote procedure call) is a style where the client calls a named procedure on the server. The client looks like a local function. The network is hidden behind a stub.

gRPC is a common RPC system. It uses HTTP/2 and Protocol Buffers (protobuf) for types. You write a `.proto` file. Tools generate client and server code.

gRPC fits internal service-to-service calls when both sides can share the schema and when you want strict types and streaming. gRPC is harder for a simple browser client. Browsers often use gRPC-Web or a JSON HTTP API instead.

Benefits:

- Generated types reduce some contract drift.
- Binary encoding can be efficient.
- Streaming is a first-class idea.

Costs:

- Browser and cache infrastructure is simpler with JSON HTTP.
- Human debugging of binary payloads needs tools.
- You still need version rules. Generated code does not remove compatibility work.

RPC is not "faster, so always use it". Measure. A small monolith does not need gRPC inside the process. Use a function call.

Treat the `.proto` file as the contract. Review it like an API. Do not change field numbers casually. Protobuf has encoding rules that depend on those numbers.

### Questions

#### Theoretical questions

1. What is RPC?
2. What is gRPC in one sentence?
3. When does gRPC fit well?
4. Why is gRPC harder for a simple browser?
5. Why must you still version a protobuf contract?

#### Easy practical tasks

1. Write four sentences that compare REST JSON with gRPC.
2. List five procedure names for a loan service in RPC style.
3. Make a table: "Client" and "JSON HTTP or gRPC?". Add browser, mobile app, and internal worker.
4. Write three costs of binary payloads for a beginner team.

#### Medium practical tasks

1. Sketch a `.proto` service with two RPCs (`GetBook`, `BorrowBook`). Describe fields in words.
2. Write an ADR: JSON HTTP for public clients, gRPC only if a second internal service appears.
3. Explain field numbers in protobuf at a high level and why reuse of a number is dangerous.

#### Advanced practical tasks

1. Write a one-page comparison of error models: HTTP status plus JSON error versus gRPC status codes.
2. Design a compatibility window for a `.proto` change that adds a field. Include generated clients.

---

## GraphQL (when it helps, when it hurts)

GraphQL is a query language and a runtime for APIs. The client sends a query that names the fields it needs. The server returns JSON that matches the query shape. One endpoint is common (`POST /graphql`).

GraphQL helps when many clients need different slices of a large graph (mobile wants few fields, an admin console wants many fields). It can reduce the count of round trips that a fixed REST resource set would need.

GraphQL hurts when:

- The team is small and a few REST resources are enough.
- Queries become expensive because clients ask for deep graphs (N+1 reads).
- Caching at HTTP is simpler with `GET` resources than with one `POST` query.
- Authorization is field-sensitive and easy to miss.
- File upload and simple cache CDNs are awkward.

A GraphQL schema is still a contract. A rename of a field is a breaking change. You still version or deprecate fields.

Do not use GraphQL as a fashion front for one table. If the client always needs the same five fields, REST is simpler.

Protect the server. Bound query depth and cost. Authenticate. Do not expose internal columns that a client must not see.

### Questions

#### Theoretical questions

1. What is GraphQL?
2. When can GraphQL reduce round trips?
3. Name four cases where GraphQL hurts?
4. Why is a GraphQL schema still a contract?
5. Why must you bound query depth?

#### Easy practical tasks

1. Write a small query in words: book title and author name only.
2. Make a table: "Client need" and "REST or GraphQL?". Add five rows.
3. List three security checks for a GraphQL endpoint (defensive).
4. Write four sentences on N+1 reads in a nested query.

#### Medium practical tasks

1. Design a REST alternative with two endpoints that covers a mobile screen. Compare with one GraphQL query.
2. Write an ADR that rejects GraphQL for a two-resource student API.
3. Sketch a schema with types `Book` and `Author`. Mark one field that you would not expose.

#### Advanced practical tasks

1. Write a one-page cost model: persisted queries, depth limit, and who may run expensive reports.
2. Compare public GraphQL and REST docs of one product. Write when each style fits that product.

---

## Versioning

Versioning is the method that you use to publish more than one contract at the same time, or to evolve a contract in a controlled way.

Common methods:

- Path version: `/v1/books` and `/v2/books`.
- Header version: `Accept: application/vnd.myapp.v1+json`.
- Additive evolution without a number for a small internal API.

Path versions are easy to see in logs and in proxies. Header versions keep URLs stable. Pick one method. Document it. Do not mix both without a reason.

A version is not a folder that you copy forever. `v2` must have a reason: a breaking change that you cannot avoid. Additive fields do not always need `v2`.

Keep `v1` until clients move. Set a sunset date. Measure which clients still call `v1`. Topic 2: observability.

Internal APIs can evolve faster if one team owns all clients. Public APIs cannot. The next sections cover compatibility and public versus internal APIs.

Write compatibility tests. Store example requests and responses. Replay them in continuous integration.

### Questions

#### Theoretical questions

1. What is API versioning?
2. What is a path version?
3. What is a header version?
4. When do you need `v2`?
5. Why do you keep `v1` after you ship `v2`?

#### Easy practical tasks

1. Write four sentences that compare path version and header version.
2. Design `/v1/users` and `/v2/users` with one field that moves.
3. Make a table: "Change" and "Needs a new version? (yes/no)". Add six changes.
4. Write a sunset notice of eight lines for `v1`.

#### Medium practical tasks

1. Document a version policy in one page: method, additive rules, sunset, and metrics.
2. Plan dual support for six months. Include logs that record the version.
3. Find a public version policy. Write three rules from that policy in your words.

#### Advanced practical tasks

1. Write a six-month migration plan from `v1` to `v2` with numbers (share of traffic) and an end date.
2. Compare path versioning for a browser app and header versioning for a machine client. Write a choice.

---

## Backward compatibility

A compatible change lets old clients continue to work. A breaking change forces clients to change.

Compatible examples:

- Add an optional field.
- Add an endpoint.
- Add an enum value only if old clients ignore unknown values (state the rule).
- Relax a constraint (accept a longer name if you still store it).

Breaking examples:

- Remove or rename a field.
- Change the type or the meaning of a field.
- Turn a success into an error for the same request.
- Reuse a field name for a new meaning.
- Change identity format (`42` versus a UUID) without a mapping period.

Rules that keep a version stable:

- Do not remove or rename fields.
- Do not change field types.
- Add new fields as optional.
- Add new endpoints instead of changing the meaning of an old endpoint.

Unknown-field rule: the server ignores unknown fields from clients if you want an older server to accept a newer client. The client ignores unknown fields from servers if you want an old client to accept a new server. Write the rule. JSON libraries differ.

Semantic change is a silent break. The field `status` still exists but "open" now means a different business state. That change is breaking even if the JSON parses.

### Questions

#### Theoretical questions

1. What is a compatible change?
2. What is a breaking change?
3. Why is a rename a breaking change?
4. Why is a new meaning of an old field a silent break?
5. What is the unknown-field rule?

#### Easy practical tasks

1. Make a table of eight changes. Mark each change as compatible or breaking.
2. Write five rules that keep `v1` stable.
3. Give three examples of a semantic break with the same JSON shape.
4. Write four sentences on why tests with stored responses help.

#### Medium practical tasks

1. Design an additive change: add `middleName`. Show old and new JSON. Write client behavior.
2. Plan a field rename with a compatibility window (both names). Then remove the old name in `v2`.
3. Write consumer-driven contract tests in words for two clients.

#### Advanced practical tasks

1. Write a one-page compatibility checklist that a reviewer uses on every API pull request.
2. Analyze a public changelog. Classify ten changes as compatible or breaking.

---

## Pagination, filtering, errors

Pagination splits a large list into pages. Without pagination, a list endpoint can return too much data. Latency and memory grow. Availability can drop.

Common styles:

- Offset and limit: `?offset=40&limit=20`. Simple. Unstable if rows insert during walk.
- Cursor: `?cursor=abc&limit=20`. More stable for infinite scroll. Harder to jump to page 10.

Always bound `limit`. A client must not request one million rows. Return a next cursor or a next offset. Document empty pages.

Filtering restricts the set (`?status=open`). Sorting must use an allow-list of fields. Do not pass raw SQL. That is a defect.

Errors need a stable shape. Example fields: `code`, `message`, `request_id`. Do not leak internals (stack traces, SQL, or paths on disk) to a public client.

Use the correct status class:

- 4xx: the client can fix the request (or cannot have access).
- 5xx: the server or a dependency failed.

Do not use `200` with a hidden error field as the only signal if you can use a status code. Some GraphQL designs use `200` with an `errors` array. Document that exception.

Rate limits use `429`. Validation errors use `400` or `422` if you adopt that convention. Be consistent.

### Questions

#### Theoretical questions

1. Why do list endpoints need pagination?
2. What is the difference between offset pagination and cursor pagination?
3. Why must `limit` have a maximum?
4. Why must sort fields use an allow-list?
5. What must a public error body not contain?

#### Easy practical tasks

1. Design `GET /books` with `limit`, `offset`, and `status` filter.
2. Write a JSON error object with `code`, `message`, and `request_id`.
3. Make a table: "Status" and "Error code name". Add 400, 401, 404, 409, and 429.
4. List four filters that a loan list needs and two filters that you reject.

#### Medium practical tasks

1. Compare offset and cursor for a busy feed where new rows insert often. Write eight sentences.
2. Design a 400 body that lists field errors for `POST /books`.
3. Write the maximum page size and the default page size in a style guide paragraph.

#### Advanced practical tasks

1. Write a one-page pagination standard: defaults, max, stable sort, and how a client walks all rows.
2. Design defensive filtering for search text (length, charset) without writing an attack.

---

## Idempotency keys for writes

An idempotency key is a client-generated identifier for one intended write. The client sends the key in a header (a common name is `Idempotency-Key`). The server stores the key and the outcome. A retry with the same key does not apply the write again.

Use keys on payments, transfers, and other writes that must not duplicate. Topic 4 defined idempotency. This section applies it to HTTP APIs.

Client rules:

- Generate a new key for a new intent.
- Reuse the key only for a retry of the same intent.
- Keep the key for the timeout and retry window.

Server rules:

- Store the key, a fingerprint of the request (or the full body hash), the status, and the response.
- If the key exists and the body matches, return the stored response.
- If the key exists and the body differs, reject the request (`409` or `422`).
- Expire old keys.

`GET` does not need a key. `PUT` of a full resource can be naturally idempotent. `POST` that creates a new payment usually needs a key.

Document the header, the time-to-live of the key, and the error when the body does not match.

Keys do not replace authentication. A key is not a secret that proves identity. Treat keys as unique identifiers. Still authenticate.

### Questions

#### Theoretical questions

1. What is an idempotency key on an HTTP write?
2. When does the client reuse a key?
3. What does the server store?
4. What happens if the key matches and the body differs?
5. Why does a key not replace authentication?

#### Easy practical tasks

1. Write a sequence of three HTTP calls: first pay, lost response, retry with the same key.
2. Make a table: "Method and path" and "Key required?". Add five rows.
3. List four fields in the server key table.
4. Write a client algorithm in six steps to generate and reuse a key.

#### Medium practical tasks

1. Design the header, TTL, and `409` body for a mismatched replay.
2. Write three tests: first success, retry, and conflicting body.
3. Explain how two tabs in a browser can send two keys for two intents (two payments).

#### Advanced practical tasks

1. Write a one-page design for keys in a load-balanced API with two instances (shared store).
2. Compare `PUT` natural idempotency with `POST` plus a key for "create payment".

---

## Public vs internal APIs

A public API has clients that you do not control: third-party developers, unknown apps, or customers. An internal API has clients that your organization owns.

Public APIs need a slower change rate, a published policy, stronger compatibility, and clearer errors. You cannot update all clients on the same day. Sunset dates must be long. Documentation is part of the product.

Internal APIs can move faster if one pipeline deploys all clients, or if you accept a coordinated release. You still need a contract. You can skip some ceremony (no public developer portal) but you must not skip timeouts, auth, and logs.

Security differs. A public API is on a low-trust boundary. Rate limits, stronger authentication, and less detailed errors are normal. An internal API still authenticates. "Inside the network" is not enough (Topic 2).

Do not expose an internal admin API on the public internet. Use a separate entry, a separate identity, and least privilege.

Some teams use the same handlers with different gateways. That design is acceptable if authorization is not only "the gateway said so" without checks in the service.

Write the audience on the API document: public, partner, or internal. Version policy follows the audience.

### Questions

#### Theoretical questions

1. What is a public API?
2. What is an internal API?
3. Why do public APIs change more slowly?
4. Why is "inside the network" not enough for an internal API?
5. Why must an admin API stay off the public internet?

#### Easy practical tasks

1. Make a table: "Rule" and "Public / internal / both". Add versioning, rate limit, and detailed stack traces.
2. Write four sentences that explain a partner API (between public and internal).
3. List five documents that a public API needs and two that an internal API can skip.
4. Classify eight endpoints of a shop as public or internal.

#### Medium practical tasks

1. Write a one-page public API policy: versions, support window, and contact.
2. Design a split: public catalog API and internal admin API. Write identity and network placement at a high level.
3. Review a public API and guess which rules would be lighter if it were internal only.

#### Advanced practical tasks

1. Write an ADR: one codebase, two gateways, two auth policies. Include the risk if a handler misses a check.
2. Design a deprecation program for a public field used by unknown clients. Include metrics and a long sunset.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you choose among REST JSON, gRPC, and GraphQL for one product with a browser and one worker?
2. How do versioning and backward compatibility work together?
3. Why do pagination and error shape belong in the same contract as resources?
4. How do idempotency keys connect Topic 4 to HTTP writes?
5. What quality attributes (Topic 2) change when an internal API becomes public?

#### Easy practical tasks

1. Write a one-page cheat sheet of all API terms in this topic.
2. Design a public `v1` books API: three resources, errors, pagination, and one write with a key.
3. Write three ADR titles: style, version method, and public versus internal.
4. Bookmark one HTTP status reference and one protobuf style guide. Write when you open each.

#### Medium practical tasks

1. Write a two-page API design for campus locker rental. Include compatibility rules and a sunset plan.
2. Take a homework JSON API. Add pagination, a stable error shape, and an idempotency header on create.
3. Role-play a breaking change request from product. Write the compatible alternative.

#### Advanced practical tasks

1. Write a full API standard (three pages) that a team of four can apply in review.
2. Compare OpenAPI and protobuf as contract sources. Write a pipeline idea that fails CI on a breaking change.
