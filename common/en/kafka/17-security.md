# 17. Security

## Description

Security on Kafka has four parts that teams mix: encryption on the network, authentication of clients and brokers, authorization (who may read or write), and limits (quotas). Encryption of data at rest is usually a disk or volume feature.

This topic covers TLS, SASL mechanisms and mTLS, ACLs and super users, encryption at rest at a high level, and quotas. Complete this topic after topics 4, 5, and 10. Apply security on KRaft clusters. Do not add ZooKeeper. Old ZooKeeper ACLs are not the model for new work.

Use one term for each concept. TLS encrypts the TCP connection and can authenticate a certificate. SASL authenticates a user name or token on that connection. An ACL is a rule that allows or denies an operation on a resource. A super user bypasses ACL checks. A quota limits how much a client may produce, fetch, or request.

---

## TLS

TLS (Transport Layer Security) encrypts bytes on the network between clients and brokers, and between brokers. Without TLS, a listener on a plain port sends records in the clear (except the payload is still your serialized bytes; anyone on the network can read them).

You configure a listener with a protocol such as `SSL` or `SASL_SSL`. You give the broker a keystore (the broker certificate and key) and a truststore (the CAs that the broker trusts). Clients get a truststore that trusts the broker certificate. For mTLS, the broker also trusts client certificates (SASL section).

TLS does not replace ACLs. TLS protects the pipe. ACLs decide the API operations.

Certificate expiry stops clients. Monitor not-after dates. Rotate certificates before they expire. Use a process. Do not copy a private key into a git repository.

KRaft controllers need TLS on the controller listener if you encrypt that path. Combined listeners and controller listeners can differ. Read the configuration for your version.

`kafka-topics` and other tools must use the same security protocol and the same bootstrap port. Example shape:

```text
kafka-topics --bootstrap-server localhost:9093 --command-config client.properties --list
```

`client.properties` sets `security.protocol`, truststore location, and passwords through a secure mechanism that you do not put in a handbook of secrets.

### Questions

#### Theoretical questions

1. What does TLS protect on a Kafka listener?
2. Why is TLS not enough without ACLs?
3. What is a keystore versus a truststore at a high level?
4. Why must you monitor certificate expiry?
5. Why can the controller listener need its own TLS settings on KRaft?

#### Easy practical tasks

1. Write five sentences about TLS. Use only facts from this section.
2. Make a table: Listener protocol, encrypted (yes or no), typical port in a lab (example only).
3. Find official TLS listener documentation. Write the property names for keystore location.
4. Draw client → TLS → broker, and broker → TLS → broker.

#### Medium practical tasks

1. Enable TLS on a local KRaft broker (self-signed in a lab). List topics with a command-config file. Save the property keys, not the passwords, in your notes.
2. Connect with the wrong truststore. Record the error. Fix the truststore.
3. Find controller listener TLS notes for KRaft. Write three facts.

#### Advanced practical tasks

1. Enable inter-broker TLS and a client TLS listener. Produce and consume. Write which ports you used.
2. Write a certificate rotation runbook: issue, distribute trust, roll brokers (topic 15), KRaft controller listener, no secrets in git.

---

## SASL: PLAIN, SCRAM, OAUTHBEARER, mTLS

SASL is a family of authentication mechanisms. Kafka uses SASL on a listener such as `SASL_SSL` (SASL on TLS) or `SASL_PLAINTEXT` (SASL without TLS). Do not use `SASL_PLAINTEXT` on a network that you do not trust. PLAIN sends a password. Without TLS, that password is visible on the wire.

**PLAIN** uses a user name and a password. Brokers or a callback check the pair. Use TLS with PLAIN.

**SCRAM** (for example SCRAM-SHA-256 or SCRAM-SHA-512) stores a hashed credential. You create users with `kafka-configs` or the Admin API. SCRAM is common on self-managed clusters.

**OAUTHBEARER** uses an OAuth 2.0 access token. A callback validates the token (issuer, audience, expiry). Cloud products often use this path. You must handle token refresh in the client.

**mTLS** (mutual TLS) authenticates the client with a client certificate. The broker maps the certificate identity to a user. mTLS can be the only authentication, or it can combine with SASL on some designs. The mapping from certificate to principal is configuration.

Inter-broker authentication also uses SASL or mTLS. Brokers must trust each other. KRaft controllers authenticate on the controller listener.

Do not invent a custom password file format. Use the mechanism that your Kafka version documents. Do not reuse production credentials in a handbook or a ticket.

Topic 10 previewed ACLs. Authentication gives a principal. Authorization uses that principal.

### Questions

#### Theoretical questions

1. Why is SASL_PLAINTEXT a poor default on an open network?
2. What does SCRAM store instead of a raw password on the broker side?
3. What does the client send with OAUTHBEARER?
4. How does mTLS prove the client identity?
5. Why must brokers authenticate to each other?

#### Easy practical tasks

1. Make a table: Mechanism, credential type, TLS required as a good practice (yes or no).
2. Write four sentences about SASL. Use only facts from this section.
3. Find official SASL documentation. List the mechanism names that the page gives.
4. Write the difference between authentication and authorization in two sentences.

#### Medium practical tasks

1. Create a SCRAM user on a KRaft lab (if your version supports it). Produce as that user. Write the commands without the password in your saved notes (use a placeholder).
2. Compare PLAIN and SCRAM in the official docs. Write five differences.
3. Read an OAUTHBEARER or cloud auth page. Write how the token is refreshed in one paragraph.

#### Advanced practical tasks

1. Configure `SASL_SSL` with SCRAM (or mTLS) for clients and for inter-broker. Produce and consume. Write the listener names.
2. Write an authentication standard: allowed mechanisms, TLS required, where credentials live, KRaft controller auth, no ZooKeeper.

---

## ACLs and Super users

An **ACL** is a rule: a principal, a host, an operation, a resource (topic, group, cluster, transactional id), and allow or deny. The authorizer on the broker checks ACLs on each request.

Typical operations: READ, WRITE, CREATE, DELETE, ALTER, DESCRIBE, ALTER_CONFIGS. A consumer needs READ on the topic and READ on the group. A producer needs WRITE on the topic. Admin tools need extra cluster or topic operations.

`super.users` is a broker list of principals that bypass ACL checks. Use it for break-glass admin users. Do not put application clients in `super.users`.

If you enable an authorizer and you have no ACLs and no matching super user, clients fail. Add ACLs before you enable enforcement on a live cluster, or use a planned window.

KRaft stores ACL metadata in the cluster metadata (current Kafka). You do not manage Kafka ACLs in ZooKeeper on a new cluster. Use `kafka-acls` with `--bootstrap-server` and a secure command config.

Example shape:

```text
kafka-acls --bootstrap-server localhost:9093 --command-config admin.properties --add --allow-principal User:alice --operation Read --topic orders.placed --group alice-group
```

Deny rules and allow rules interact. Read the authorizer documentation for the default (allow versus deny) on your version.

Topic 10 was a preview. This section is the operations view. Review every ACL change.

### Questions

#### Theoretical questions

1. What five parts does an ACL have at a high level?
2. What two resources does a typical consumer need?
3. What is a super user?
4. Why must you not make an application a super user?
5. Where do you store ACLs on a new KRaft cluster?

#### Easy practical tasks

1. Write five sentences about ACLs. Use only facts from this section.
2. Make a table: Client type, operations, resources.
3. Run `kafka-acls --help`. Write five flags.
4. Find `super.users` in official docs. Write the principal format example that the page uses (rewrite in your own words).

#### Medium practical tasks

1. Enable an authorizer on a KRaft lab. Add ACLs for one producer and one consumer. Confirm a third user fails.
2. List ACLs. Write the command and the output (no secrets).
3. Compare `--bootstrap-server` ACL commands with any old ZooKeeper ACL example you find. Write what you refuse to run.

#### Advanced practical tasks

1. Design ACLs for Connect (topic 12) or Streams (topic 14): internal topics, application topics, and the group or `application.id`. Apply them in a lab.
2. Write an ACL standard: naming of principals, who may be a super user, review process, KRaft only.

---

## Encryption of data at rest (high-level / volume)

**Encryption at rest** means disk bytes are encrypted when they sit on storage. Kafka does not implement a full field-level encryption product in the broker as the default. Teams encrypt the volume or the file system (cloud disk encryption, LUKS, or a vendor volume key).

If someone steals a disk image and the volume key is not with the image, the `.log` files are not readable as plain records. If the broker is up, it has the key or the OS has mounted the volume. Then ACLs and TLS still matter.

Application-level encryption (encrypt the value before produce) is a different control. The broker still stores bytes. Consumers need the key. Schema Registry and converters must match (topic 11). This handbook does not specify a crypto library.

Compaction, retention, and backup copies (snapshots of `log.dirs`) must use the same at-rest policy. A backup on an unencrypted bucket is a leak.

KRaft metadata logs also sit on disk. Encrypt those volumes too on controller nodes.

This section is high-level. Follow your platform encryption standard. Do not invent a custom broker plugin as a first step.

Quotas and ACLs do not encrypt disks. TLS does not encrypt disks.

### Questions

#### Theoretical questions

1. What does encryption at rest protect?
2. Why is volume encryption different from TLS?
3. What happens if a backup of `log.dirs` is not encrypted?
4. Why must KRaft controller disks also use the at-rest policy?
5. How is application-level payload encryption different from volume encryption?

#### Easy practical tasks

1. Write four sentences about encryption at rest. Use only facts from this section.
2. Make a table: Control, protects in transit, protects on disk, protects from a live authorized client.
3. Find your cloud or OS volume encryption name. Write one sentence on how a VM disk is encrypted.
4. Draw: disk → volume key → OS → Kafka `.log` files.

#### Medium practical tasks

1. On a lab VM or Docker volume, write where `log.dirs` lives and whether the volume is encrypted. Write the answer honestly.
2. List three backup targets (snapshot, object copy, another disk). Write the encryption requirement for each.
3. Read official Kafka notes if they mention encryption at rest. Write what Kafka does and does not do.

#### Advanced practical tasks

1. Write a one-page data-protection note: TLS, SASL, ACLs, volume encryption, payload encryption optional, backup rules.
2. Design a key-rotation idea for volume keys versus certificate rotation. Do not implement a custom broker cipher.

---

## Quotas

A **quota** limits a client’s use of the cluster. Kafka supports produce quotas (bytes per second), fetch quotas (bytes per second), and request-rate quotas. You set quotas per user, per client id, or both.

When a client exceeds a quota, the broker delays or throttles that client. Other clients keep capacity. Quotas are a fairness and safety tool. They are not a replacement for ACLs.

Set quotas so that a bad consumer or a test producer cannot fill the disk path or the network for everyone. Connect and Streams clients need quotas that match their job. A too-low produce quota on a source connector creates lag in the external system.

Default quotas can exist in broker configuration. Overrides use `kafka-configs` (topic 10) with `--bootstrap-server`.

Quotas do not hide data. A throttled client still must authenticate and pass ACLs.

Monitor throttle time metrics. A client that is always throttled needs a higher quota or a fix in the client.

KRaft does not replace quotas. Quotas run on brokers for the data plane.

### Questions

#### Theoretical questions

1. What three kinds of quotas does this section name?
2. What does the broker do when a client exceeds a quota?
3. Why are quotas not a replacement for ACLs?
4. Why can a too-low quota break Connect or Streams?
5. How do you apply a quota override?

#### Easy practical tasks

1. Write five sentences about quotas. Use only facts from this section.
2. Make a table: Quota type, unit, typical victim if missing.
3. Find quota configuration in official docs. Write two property or entity names.
4. Draw: two producers, one with a quota, one without. Mark who is delayed.

#### Medium practical tasks

1. Set a very low produce quota on a lab user. Produce in a loop. Write whether you see throttle or errors.
2. List quotas with `kafka-configs`. Write the command shape.
3. Read request-quota documentation. Write how it differs from a byte quota.

#### Advanced practical tasks

1. Design quotas for: interactive app, bulk import, Connect sink. Write numbers as a lab starting point and how you will adjust.
2. Write a quota standard: defaults, how you name client.id, who may raise a quota, metrics you watch.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a produce request that uses TLS, SASL, ACLs, and a produce quota. What check happens at each layer?
2. Why must authentication, authorization, encryption in transit, and encryption at rest all exist together?
3. What is different about securing a KRaft controller listener versus a client listener?
4. How do Connect and Streams change the ACL and quota plan compared with one producer?
5. What from topic 10 do you reuse here, and what is new?

#### Easy practical tasks

1. Write a one-page cheat sheet: TLS, four auth methods, ACL parts, super user, at rest, three quotas.
2. Draw a locked-down lab: ports, protocols, principal, topic ACL, volume lock.
3. Bookmark official security, ACL, and quota pages.
4. List every tool flag in this topic that uses `--bootstrap-server` or `--command-config`.

#### Medium practical tasks

1. Build a KRaft lab: TLS + one SASL mechanism + ACLs + a quota. Produce and consume as the allowed user. Fail as a second user.
2. Map each subsection to one official URL.
3. Write a threat table: stolen disk, stolen password on the wire, extra produce load, missing ACL. Map each to a control in this topic.

#### Advanced practical tasks

1. Add inter-broker security and a rolling restart (topic 15) without downtime in a three-node KRaft lab. Write the order of certificate and ACL steps.
2. Write a production security standard: protocols, mechanisms, ACL review, volume encryption, quotas, secret storage, no ZooKeeper.
