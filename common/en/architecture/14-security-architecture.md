# 14. Security Architecture

## Description

Security architecture is the set of trust boundaries, identity rules, access rules, and secret-handling rules that the structure must support. This topic covers trust boundaries, authentication versus authorization, OAuth 2.0 and OpenID Connect roles, secrets and keys, least privilege, simple STRIDE threat modeling, and supply-chain risk from dependencies.

This topic is defensive. Use it to reduce accidents and abuse. Do not use it to plan attacks or to bypass controls. Follow the law and the rules of your school or employer. Complete Topics 1 to 13 before this topic. Topic 2 defines security as a quality attribute. This topic shows structure that supports that quality.

Use one term for each concept. Authentication is not authorization. A client is not a resource server. A secret is not a public identifier. STRIDE is a review checklist. It is not a product.

---

## Trust boundaries

A trust boundary is a line where the level of trust changes. Data and calls that cross the line need extra checks. Topic 2 introduced the idea. This section places it on a structure.

Typical zones:

- Public internet and browsers
- Edge (load balancer or API gateway, Topic 12)
- Application processes
- Data stores
- Operator and admin paths
- Third-party APIs

A request from an "internal" network is not automatically trusted. Internal networks get scanned. Misconfiguration happens. Authenticate principals. Authorize actions. Encrypt in transit across hostile or shared networks.

Draw the boundary on a C4 container diagram. Mark each arrow that crosses a line. For each arrow write:

- What identity travels
- What data classification travels (public, internal, personal)
- What happens if the caller is fake or replayed

A monolith still has trust boundaries: browser to process, process to database, process to mail vendor. Microservices add more lines. More lines mean more policy and more chance of a missed check (Topic 12).

Do not invent a "secure zone" that skips authentication because a firewall exists. A firewall is one control. Identity is another control.

Admin interfaces need a stricter boundary. Use a separate host name, a stronger identity, and a shorter session. Log admin actions.

### Questions

#### Theoretical questions

1. What is a trust boundary?
2. Why is an internal network not automatic trust?
3. What three facts do you write on each arrow that crosses a boundary?
4. Why do microservices increase the number of trust boundaries?
5. Why is a firewall not a replacement for authentication?

#### Easy practical tasks

1. Draw two zones for a notes app: browser and server. Mark the boundary and one control on the arrow.
2. Make a table: "Zone" and "Trust (low/medium/high)". Add internet, app process, and admin path.
3. List five arrows in a shop: browser, gateway, API, database, payment vendor.
4. Write four sentences that explain why a monolith still has trust boundaries.

#### Medium practical tasks

1. Draw a C4 container view of a campus loan system. Mark every trust boundary. Classify data on two arrows.
2. Write a one-page rule: no write from the public zone reaches the database without an authenticated principal.
3. Compare admin access through the public API with admin access through a separate host. Write eight sentences.

#### Advanced practical tasks

1. Write a trust-boundary register (table): arrow, data class, identity, control, owner. Fill twelve rows.
2. Take a public architecture diagram. Redraw it with explicit trust lines. List three missing controls. Stay defensive.

---

## Authentication vs authorization

Authentication answers "who is this principal?". Authorization answers "what may this principal do?". Topic 2 uses the same pair. Mix of the two words causes defects.

A principal is a user, a service account, or a machine identity. Authentication binds a request to a principal with a credential or a token. Sessions and tokens are proof of an earlier authentication.

Authorization uses roles, permissions, or policies. "The user logged in" is not enough. A student who is authenticated must not write grades.

Write authorization next to the data owner. The service that owns the grade record must check teacher rights. A gateway can reject a missing token. The gateway must not be the only place that knows "this teacher may change this course" if other paths can reach the service (Topic 12).

Fail closed. If identity is missing or the policy store is down, deny the sensitive action. Record the deny. Do not fail open on admin paths.

Identity is not a display name. Use stable identifiers. Do not put passwords in logs (next sections).

Public pages can skip authentication. They still need other controls: rate limits, size limits, and validation.

### Questions

#### Theoretical questions

1. What is authentication?
2. What is authorization?
3. What is a principal?
4. Why must the data owner enforce authorization?
5. What does fail closed mean when the policy store is down?

#### Easy practical tasks

1. Write five sentences that separate authentication from authorization.
2. Make a table: "Check" and "Authn or authz". Add "password correct", "role is teacher", and "token missing".
3. List four principals in a campus system.
4. Write two requirements: one authentication requirement and one authorization requirement. Include a yes/no check.

#### Medium practical tasks

1. Design roles for a grade book: student, teacher, registrar. Write read and write rights per resource.
2. Write eight sentences on why a gateway-only authorization check fails if an internal caller can reach the service.
3. Describe fail-closed behavior for "edit grade" when the identity service does not reply. Write the user message.

#### Advanced practical tasks

1. Write a one-page access matrix: resource × role × action. Include a service-to-service principal.
2. Compare session cookies and bearer tokens at a high level. Write where authorization still must run. Do not write attack steps.

---

## OAuth2 / OIDC (roles of client, IdP, resource server)

OAuth 2.0 is a framework for delegated authorization. A user can allow a client to access a resource server without sharing a password with that client. OpenID Connect (OIDC) is an identity layer on OAuth 2.0. OIDC adds an ID token that states who authenticated.

Learn the roles. Do not start with a product name.

- Resource owner: the user or organization that owns the data.
- Client: the application that requests access.
- Authorization server (identity provider, IdP): authenticates the user and issues tokens.
- Resource server: the API that holds the data and checks the access token.

The client is often not the resource server. A single-page app is a client. Your API is a resource server. The campus login system can be the IdP.

OIDC is how many teams implement "log in with the campus account". The ID token is for identity. The access token is for API access. Do not treat them as the same token.

Architecture duties:

- Choose one IdP. Do not invent a password store if the organization already has one.
- Register clients. Use the correct client type for the environment (public browser versus confidential server).
- Validate tokens on the resource server. Check issuer, audience, expiry, and signature with the IdP keys.
- Keep scopes small. A scope is not a full authorization model. Your API still checks object-level rights.

Do not put long-lived secrets in a public mobile binary or in a front-end source bundle. Topic 20 notes client limits.

This section is a role map. It is not a grant-type cookbook. Read the IdP documentation of your school or vendor when you implement. Stay on approved flows. Do not design a custom token format without a reason.

### Questions

#### Theoretical questions

1. What problem does OAuth 2.0 address?
2. How does OIDC differ from OAuth 2.0?
3. What are the four roles in this section?
4. What does the resource server check on an access token?
5. Why is a scope not enough for object-level authorization?

#### Easy practical tasks

1. Write four sentences that assign roles for "campus app reads mail through the campus IdP".
2. Make a table: "Role" and "Example". Add client, IdP, resource server, and resource owner.
3. List five token checks (issuer, audience, expiry, signature, and one more that you name).
4. Open your school or a public IdP overview page. Write the product name that acts as IdP.

#### Medium practical tasks

1. Draw the roles for a shop API and a browser app. Label ID token versus access token at a high level.
2. Write an ADR: use the organization IdP. Reject a homemade password table. Include two consequences.
3. Write eight sentences on why a public client must not hold a client secret.

#### Advanced practical tasks

1. Write a one-page role-and-token brief for three apps: browser, server-side web, and internal worker. Stay at role level.
2. Read a public OIDC overview (not a bypass guide). Write ten sentences on ID token versus access token. Do not copy long passages.

---

## Secrets and key management

A secret is a value that grants access: a password, an API token, a private key, or a database credential. A public identifier is not a secret. A key is a secret that a cryptographic function uses. This handbook groups passwords, tokens, and keys as secrets unless a sentence names a key type.

Rules:

- Do not store secrets in source control.
- Do not put secrets in images, tickets, or chat.
- Do not log secrets.
- Inject secrets at runtime from a secret store that operations controls.
- Rotate secrets. Record who can read them.
- Use different secrets per environment.

A secret store is a service or a platform feature that holds values and grants them to identified workloads. Cloud platforms and operators provide these stores. The application reads the value at start or through a sidecar. The application does not write the value into a config file in the repository.

Key management includes generation, storage, rotation, and destruction of keys. Prefer the platform key service for disk encryption and for signing if the team cannot operate a key ceremony.

Certificates expire. An expired TLS certificate is an availability incident (Topic 2). Calendar the expiry. Automate renewal when the platform allows it.

Incident habit: if a secret was in a log or a repository, rotate it. Do not only delete the line.

This material is defensive. Do not use it to extract secrets from systems that you do not own.

### Questions

#### Theoretical questions

1. What is a secret in this handbook?
2. Where must secrets not live?
3. What is a secret store for?
4. Why is certificate expiry an availability incident?
5. What do you do if a secret appeared in a log?

#### Easy practical tasks

1. Make a table: "Item" and "Secret? (yes/no)". Add API tokens, user display names, TLS private keys, and documentation URLs.
2. List four safe injection methods at a high level (platform store, environment from the orchestrator, and two more that you name).
3. Write five sentences on rotation. Include a restart or reload step.
4. Write a bad config snippet in words (password in a file in git) and the replacement in words.

#### Medium practical tasks

1. Plan rotation of a database password in six steps. Include application reload and a check.
2. Write a one-page secret policy for a student project: store, owners, and a rotate-after-exposure rule.
3. Explain in eight sentences why a secret in a container image is still a secret in source history if the Dockerfile copied it.

#### Advanced practical tasks

1. Write an ADR: platform secret store versus encrypted files in the repo. Include team skill as a constraint (Topic 1).
2. Design a certificate expiry dashboard requirement. Link it to an SLO in Topic 15 terms (later). Stay defensive.

---

## Least privilege on networks and data

Least privilege is a rule: give each principal only the rights that the work needs. The rule applies to people, processes, tokens, and network paths. Topic 2 stated the idea. This section applies it to structure.

On data:

- A report job that only reads must use a read-only database role.
- A web process must not use a role that can drop tables.
- Personal data stays in the owner store. Other services receive the minimum fields (Topic 9).
- Admin actions use a separate role and a separate audit trail.

On networks:

- A worker that only talks to a broker does not need a public ingress.
- Databases accept connections from the application network, not from the public internet.
- Management ports stay off the public edge.
- East-west calls use identity, not only a private IP.

Too much privilege hides in "one admin account for all environments" and in "the API uses the database superuser". Those shortcuts fail an audit and enlarge an incident.

Start from deny. Add a path when a feature needs it. Record the path in a diagram.

Least privilege can increase toil. A new feature needs a new grant. That toil is cheaper than a broad credential that can destroy data.

Do not confuse least privilege with a complex mesh that the team cannot operate. A few clear network rules plus strong data roles beat a diagram that no one can change.

### Questions

#### Theoretical questions

1. What is least privilege?
2. Why must a report job not use a role that can drop tables?
3. Why must a database not accept public internet connections?
4. How can least privilege enlarge toil, and why do you still use it?
5. When does a complex network diagram fail least privilege in practice?

#### Easy practical tasks

1. Write four sentences that define least privilege for a process identity.
2. Make a table: "Workload" and "Rights". Add web API, nightly report, and migrator.
3. List five network paths that a small shop must close at the edge.
4. Write two database roles: `app_rw` and `report_ro`. List allowed statements in words.

#### Medium practical tasks

1. Design grants for three services and one database cluster. No service uses a superuser.
2. Write a one-page network rule set for API, worker, database, and broker. Include who can initiate a connection.
3. Review a two-person project that uses one admin credential everywhere. Write an eight-sentence repair plan.

#### Advanced practical tasks

1. Write a privilege matrix: principal × resource × action × environment. Fill fifteen cells for a library system.
2. Compare security groups (or equivalent) with a service mesh at a high level. Write operations cost for a five-person team. Stay defensive.

---

## Threat modeling (STRIDE, simple)

Threat modeling is a structured review of how a system can fail in security terms. You do it to add controls. You do not do it to write attack recipes.

STRIDE is a simple mnemonic:

- Spoofing: a principal pretends to be another principal.
- Tampering: data or code changes without right.
- Repudiation: an actor can deny an action because the log is weak.
- Information disclosure: data reaches a person or a process that must not see it.
- Denial of service: the system cannot do useful work.
- Elevation of privilege: a principal gains rights that the design did not grant.

How to run a short review:

1. Draw the data-flow diagram with trust boundaries.
2. For each process and store, ask the six STRIDE names.
3. Write one realistic failure per name that applies.
4. Write one defensive control per failure.
5. Assign an owner and a test.

Keep the session short. Two hours on one feature is better than a 40-page model that no one reads.

Controls are boring on purpose: authentication, integrity checks, audit logs, encryption in transit and at rest, rate limits, and least privilege. Topic 2 asked for testable security requirements. The model must produce those requirements.

Do not include exploit steps, payloads, or bypass instructions in the notes. Name the failure and the control.

Law and abuse cases belong in the same review: retention, consent, and child-data rules if they apply. Invite the stakeholder who owns that constraint (Topic 1).

### Questions

#### Theoretical questions

1. What is threat modeling in this handbook?
2. What does each STRIDE letter name?
3. What five steps does a short review use?
4. Why must the model produce testable requirements?
5. What must the notes omit?

#### Easy practical tasks

1. Write the STRIDE names and one defensive control each, in your own words.
2. Draw a data-flow of `POST /comments`. Mark one trust boundary.
3. Make a table: "STRIDE name" and "Control". Add six rows. Stay defensive.
4. Write four sentences on why a long unused model fails.

#### Medium practical tasks

1. Run a paper STRIDE review on a public comment form. Write six failures and six controls. Do not write attack steps.
2. Write a one-page template that a team fills in 90 minutes: diagram, six rows, owners.
3. Map repudiation to audit-log fields for "change grade". Include who, what, when.

#### Advanced practical tasks

1. Facilitate a written tabletop: stolen session token as a quality scenario (Topic 2). Describe detection, impact limit, and recovery. Stay defensive.
2. Compare STRIDE with a simple abuse-case list. Write when each method helps a small team.

---

## Supply chain (dependencies)

The software supply chain is the set of source, builds, dependencies, images, and publishers that produce what you run. A defect or a hostile change in a dependency becomes your defect.

Architecture duties:

- Record dependencies (lock files).
- Prefer packages that the team can update.
- Limit install from unknown sources.
- Pin versions in production builds.
- Scan for known published vulnerabilities with a tool that the team can run in CI (Topic 18).
- Separate build from runtime. Do not compile untrusted code on the production host.

Transitive dependencies multiply the surface. One direct library can pull many others. Review the lock file when you add a direct dependency.

A "small helper" that saves 20 lines can add a large tree. Topic 2 treat cost and maintainability as qualities. A dependency is a maintainability and security decision. Write an ADR for unusual or high-risk packages.

Images and base layers are dependencies. Rebuild on a schedule. Do not run an image with a known critical defect in the base layer if a repair exists.

People are part of the chain. A single person who can push to production without review is a supply-chain risk. Use reviews and signed-off pipelines when the organization requires them.

This section is defensive. Use it to reduce accidental exposure. Do not use it to design malicious packages.

### Questions

#### Theoretical questions

1. What is the software supply chain in this handbook?
2. Why does a lock file matter?
3. Why can a small helper be a large risk?
4. How are container images part of the chain?
5. How can a missing review be a supply-chain risk?

#### Easy practical tasks

1. List five parts of a supply chain for a student web app.
2. Open a lock file in a project that you own. Count direct and a sample of transitive names.
3. Write four sentences on pinning versions.
4. Make a table: "Action" and "Risk if skipped". Add lock, pin, scan, and rebuild image.

#### Medium practical tasks

1. Write a one-page dependency policy: who may add a package, what review you require, and how fast you patch.
2. Take one popular library. Write why you accept it or reject it: maintainer activity, license, and tree size. Do not copy long docs.
3. Plan a CI check that fails a build on a stated severity. Write who gets the ticket.

#### Advanced practical tasks

1. Write an ADR: add a dependency scanner and a weekly update window. Include false-positive handling.
2. Draw the chain from developer laptop to production image. Mark five controls. Stay defensive.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do trust boundaries, authentication, and authorization form one path from a browser click to a row in a store?
2. How do OAuth 2.0 roles change what you store in your application versus what the IdP stores?
3. How do secrets, least privilege, and supply-chain controls reduce the size of an incident?
4. Why does this path place security architecture after APIs, data ownership, and integration styles?
5. What makes a STRIDE note useful to an operator at 03:00?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic.
2. For a to-do API, write two trust zones, one IdP choice, three secrets, and three least-privilege rules.
3. Draw a C4 container diagram and overlay trust boundaries and token flow at a high level.
4. Write three ADR titles: organization IdP, secret store, reject superuser in the API.

#### Medium practical tasks

1. Write a short security architecture brief (one or two pages) for a campus lost-and-found system. Stay defensive. Include STRIDE rows and a dependency rule.
2. Take a teammate design that only says "we use OAuth". Rewrite it as roles, token checks, and object-level authorization.
3. Map Topics 2, 8, 9, and 12 onto one "submit assignment" flow. Write the control that each topic forces.

#### Advanced practical tasks

1. Write a two-page threat model for one feature only. Include diagram, STRIDE, owners, and tests. Omit exploit steps.
2. Read a public defensive guide on OIDC or secret stores (official docs). Write ten sentences that you would add to a team standard. Do not copy long passages.
