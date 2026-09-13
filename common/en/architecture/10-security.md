# 10. Security

## Description

Security architecture is the set of trust boundaries, identity rules, access rules, and secret-handling rules that the structure must support. This topic covers trust boundaries, authentication versus authorization, OAuth 2.0 and OpenID Connect roles, secrets and least privilege, threat modeling, and supply-chain risk from dependencies.

This topic is defensive. Use it to reduce accidents and abuse. Do not use it to plan attacks or to bypass controls. Follow the law and the rules of your school or employer. Complete Topics 1 to 9 before this topic. Topic 1 defines security as a quality attribute. This topic shows structure that supports that quality.

Use one term for each concept. Authentication is not authorization. A client is not a resource server. A secret is not a public identifier. STRIDE is a review checklist. It is not a product.

---

## Trust boundaries

A trust boundary is a line where the level of trust changes. Data and calls that cross the line need extra checks.

Typical zones:

- Public internet and browsers
- Edge (load balancer or API gateway, Topic 9)
- Application processes
- Data stores
- Operator and admin paths
- Third-party APIs

A request from an "internal" network is not automatically trusted. Internal networks get scanned. Misconfiguration happens. Authenticate principals. Authorize actions. Encrypt in transit across hostile or shared networks.

Draw the boundary on a C4 container diagram. Mark each arrow that crosses a line. For each arrow write:

- What identity travels
- What data classification travels (public, internal, personal)
- What happens if the caller is fake or replayed

A monolith still has trust boundaries: browser to process, process to database, process to mail vendor. Microservices add more lines. More lines mean more policy and more chance of a missed check (Topic 9).

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

Authentication answers "who is this principal?". Authorization answers "what may this principal do?". Mix of the two words causes defects.

A principal is a user, a service account, or a machine identity. Authentication binds a request to a principal with a credential or a token. Sessions and tokens are proof of an earlier authentication.

Authorization uses roles, permissions, or policies. "The user logged in" is not enough. A student who is authenticated must not write grades.

Write authorization next to the data owner. The service that owns the grade record must check teacher rights. A gateway can reject a missing token. The gateway must not be the only place that knows "this teacher may change this course" if other paths can reach the service (Topic 9).

Fail closed. If identity is missing or the policy store is down, deny the sensitive action. Record the deny. Do not fail open on admin paths.

Do not put credentials in URLs (Topic 5). Do not log tokens. Use HTTPS on public networks (`net.topics.md`).

Session design is architecture. A server-side session in a store supports logout. A long-lived token in a browser has a theft risk. Write the trade-off (Topic 1).

Service-to-service calls also need identity. A network name is not enough. Use a machine identity that you can revoke.

### Questions

#### Theoretical questions

1. What is authentication?
2. What is authorization?
3. Why is "the user logged in" not enough?
4. Why must the data owner check authorization?
5. What does fail closed mean?

#### Easy practical tasks

1. Write five sentences that separate authentication and authorization.
2. Make a table: "Check" and "Authn or authz". Add six rows.
3. List four principals in a campus system (student, teacher, API worker, admin).
4. Write four sentences on why a gateway-only permission check is weak if the service is reachable.

#### Medium practical tasks

1. Design login and a "change grade" check. Write where each check runs.
2. Write an ADR: server-side session versus token for a campus app. Include logout and theft risk.
3. Explain in eight sentences how a worker that calls the grade API authenticates.

#### Advanced practical tasks

1. Write a one-page access model: roles, resources, and who checks each rule.
2. Design fail-closed behavior when the policy store is down. Write user-visible results for read catalog versus write grade.

---

## OAuth2 / OIDC roles

OAuth 2.0 is a framework for delegated authorization. A user can allow a client application to access a resource without sharing the user password with that client. OpenID Connect (OIDC) is an identity layer on OAuth 2.0. OIDC adds an ID token that states who the user is.

Learn the roles. Do not mix the names.

- Resource owner: the user (or the owner of the data).
- Client: the application that requests access.
- Authorization server: the server that authenticates the user and issues tokens.
- Resource server: the API that accepts a token and serves the data.

OIDC adds the identity provider role (often the same process as the authorization server). The client receives an ID token for authentication and an access token for the API.

Typical web flow at a high level: the user signs in at the authorization server. The client receives tokens. The client calls the resource server with the access token. The resource server validates the token and authorizes the action.

Rules:

- The client is not the authorization server.
- Your API is often the resource server. Validate tokens. Do not invent a private OAuth if a campus provider exists.
- Pick a flow that matches the client type. A public browser client is not a confidential server client. Follow current official guidance for the flow that your provider supports.
- Scopes are coarse. They are not a full authorization model. Still check object-level rights in the owner.

Do not use this section to bypass a login page or to steal tokens. Use it to place roles on a diagram and to avoid a homemade password dance with a vendor.

Store tokens as secrets (next section). Prefer short access-token lifetime and a refresh path that the provider documents.

### Questions

#### Theoretical questions

1. What problem does OAuth 2.0 solve?
2. What does OIDC add?
3. Name the four OAuth roles in this handbook.
4. What is a resource server?
5. Why are scopes not a full authorization model?

#### Easy practical tasks

1. Write five sentences that define OAuth 2.0 and OIDC. Use only facts from this section.
2. Make a table: "Role" and "Example in a campus login". Add four rows.
3. Draw user, client, authorization server, and resource server.
4. List four mistakes: client stores user passwords, API skips token validation, and similar.

#### Medium practical tasks

1. Write an ADR: use the campus OIDC provider; your API is the resource server.
2. Explain in eight sentences the difference between an ID token and an access token at a high level.
3. Map a modular monolith: which process is the client, and which process is the resource server, for a browser app.

#### Advanced practical tasks

1. Write a one-page role diagram for a mobile app plus a web app plus one API. Stay on official role names.
2. Read a public OIDC provider document index. Write ten sentences in your own words about roles. Do not copy long passages.

---

## Secrets and least privilege

A secret is a value that proves identity or that decrypts data: passwords, API keys, private keys, and session signing keys. A public identifier is not a secret. A client id can be public. A client secret must not be public.

Rules for secrets:

- Do not put secrets in source control.
- Do not put secrets in URLs, logs, or error pages.
- Use a secret store or the platform environment that operators control (Topic 12).
- Rotate secrets. Write who can rotate them.
- Bound lifetime. A leaked long-lived key is worse than a leaked short-lived token.

Least privilege is the rule that a principal gets only the rights that the job needs. The API process must not use a database role that can drop all schemas if it only needs to read and write its tables. A student role must not include admin APIs.

Apply least privilege to:

- People (roles in the product)
- Processes (OS user, cloud role, database role)
- Networks (which host can reach which port)
- Data fields (do not return extra personal fields, Topic 5)

A shared production password is a failure of ownership (Topic 6) and of least privilege.

Encrypt in transit on public and shared networks. Encrypt sensitive data at rest when the threat model requires it (next section). Encryption does not replace access control. A process that can read the key can read the data.

### Questions

#### Theoretical questions

1. What is a secret in this handbook?
2. Why must secrets stay out of source control?
3. What is least privilege?
4. Name three places you apply least privilege.
5. Why does encryption not replace access control?

#### Easy practical tasks

1. Write five sentences about secrets and least privilege.
2. Make a table: "Value" and "Secret? (yes/no)". Add six rows.
3. List four places a student project often leaks a key (repo, screenshot, log).
4. Write a database role list: app write, app read replica, admin migrate.

#### Medium practical tasks

1. Write an ADR: secrets in the platform store, never in `git`. Include rotation owner.
2. Design least privilege for a mail worker: what it may read, what it must not write.
3. Explain in eight sentences how a shared production password breaks audit.

#### Advanced practical tasks

1. Write a one-page secret standard: storage, rotation, leak response, and who is paged.
2. Build a privilege matrix: five principals × five actions on grades. Mark allow or deny.

---

## Threat modeling and supply chain

Threat modeling is a structured review of how the system can fail in a security sense. A simple method is STRIDE as a checklist:

- Spoofing: fake identity
- Tampering: change of data
- Repudiation: deny an action because logs are weak
- Information disclosure: leak of data
- Denial of service: make the service unusable
- Elevation of privilege: gain extra rights

Walk each trust-boundary arrow. Ask which STRIDE items apply. Write a mitigation that you already have or that you will add. Stay defensive. Do not write exploit steps.

A supply chain is the set of dependencies that you did not write: libraries, container base images, build tools, and vendors. A malicious or a defective dependency becomes your defect.

Rules for supply chain:

- Pin versions. Record checksums when the toolchain supports it.
- Run a known scan that the team can operate.
- Do not add a library for ten lines (cost and surface).
- The pipeline is the only supported path to production (Topic 12). A laptop copy is a supply-chain risk.
- Review who can publish to your package registry.

Threat modeling without an owner is a poster. Assign an owner. Review again when you add a boundary (a new vendor, a new admin path).

Law and data class belong in the model. Personal data on a new log sink is a new disclosure risk.

### Questions

#### Theoretical questions

1. What is threat modeling in this handbook?
2. Name the six STRIDE items.
3. What is a software supply chain here?
4. Why is a laptop copy to production a supply-chain risk?
5. Why must a new vendor start a new review?

#### Easy practical tasks

1. Write five sentences about threat modeling and supply chain. Use only facts from this section.
2. Make a table: "STRIDE item" and "One mitigation". Add six rows. Stay defensive.
3. List four dependencies in a typical student API (language runtime, web library, database driver, base image).
4. Write four questions you ask before you add a new library.

#### Medium practical tasks

1. Walk five arrows of a campus loan system with STRIDE. Write one mitigation each. Stay defensive.
2. Write an ADR: pin dependencies and forbid unknown `latest` tags on production images.
3. Explain in eight sentences how a log that stores personal data is an information-disclosure risk.

#### Advanced practical tasks

1. Write a one-page threat-model template: assets, boundaries, STRIDE, owners, review date.
2. Design a pipeline check list (names only) for build, test, dependency scan, and signed artifact. Do not write attack tools.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do trust boundaries, authentication, and authorization stack on one request?
2. Why can a correct OAuth role diagram still fail if the resource server skips object-level checks?
3. How do secrets and least privilege support the same quality attribute?
4. Why do more microservices increase both trust lines and supply-chain surface?
5. How does STRIDE use the C4 arrows that Topic 1 introduced?

#### Easy practical tasks

1. Write a one-page cheat sheet: boundary, authn/authz, OAuth roles, secrets, least privilege, STRIDE, supply chain.
2. For a to-do API, draw two boundaries, one login, one owner check, and no secrets in git.
3. Bookmark a campus or public OIDC overview. Write the four OAuth roles in one diagram.
4. Write a deny list: secrets in repo, tokens in logs, `GET` with a password query.

#### Medium practical tasks

1. Write a short security brief for a campus lost-and-found: boundaries, identity, data class, and one vendor ACL (Topic 9).
2. Take a teammate design that trusts the internal network. Add identity and least privilege on paper.
3. Write a twelve-week plan: weeks for login and owner checks, and the review that adds STRIDE.

#### Advanced practical tasks

1. Write a security review checklist that a teammate uses on a pull request (boundaries, authz location, secrets, dependencies).
2. Map this topic to Topic 5 (public API) and Topic 6 (data ownership) in a one-page table: control and owner.
