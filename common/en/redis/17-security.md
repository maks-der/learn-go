# 17. Security

## Description

Redis does not start as a public internet service. This topic covers passwords, ACLs, bind addresses, protected mode, TLS, and commands that can destroy data or stall the instance.

Complete this topic before you put Redis on a shared network. Use a lab instance for every practical task that changes auth or command names. Do not practice on a shared production instance.

Use one term for each concept. Authentication is proof of identity to Redis. An ACL user is a named account with a command list and a key pattern. Protected mode is a default safety when Redis has no bind and no password. TLS encrypts the connection. A dangerous command is a command that can wipe data, stall the server, or change security settings.

This topic is hardening. It does not describe attacks.

---

## `requirepass` / ACLs (Redis 6+)

Older Redis used one password for the whole instance:

```text
CONFIG SET requirepass labpass
AUTH labpass
```

`requirepass` is a single shared secret. Every client that knows it has full power. Redis 6 adds Access Control Lists (ACLs). Prefer ACLs on Redis 6 and newer.

ACL ideas:

- A default user exists. You can set a password on it or disable it.
- `ACL SETUSER` creates or changes a user.
- A user has rules: commands (`+GET`, `-FLUSHALL`), categories (`+@read`), and key patterns (`~cache:*`).
- `AUTH username password` selects the user.

```text
ACL LIST
ACL WHOAMI
ACL SETUSER app on >apppass ~cache:* +@read +SET +DEL -@dangerous
```

Syntax details are on the `ACL` command pages. Check your version. The `>` prefix in `SETUSER` sets a password. `on` enables the user.

Application clients use a least-privilege user. A human operator uses a stronger user for `CONFIG` and `ACL` work. Do not give the application `FLUSHALL`, `DEBUG`, or `KEYS`.

Store passwords outside the repository. Use a secret manager or the environment. Do not put production passwords in shell history if you can use a prompt or a file with tight permissions.

`ACL SAVE` writes the ACL file when you use an ACL file. `CONFIG REWRITE` can persist `redis.conf` changes. A `CONFIG SET` password that you do not persist can disappear after restart. Test restart.

Hosted Redis often requires a user and TLS. Follow the vendor docs.

### Questions

#### Theoretical questions

1. What is the limit of a single `requirepass`?
2. Which Redis version introduced ACLs?
3. What three kinds of ACL rules does this section name?
4. Why must the application user lack `FLUSHALL`?
5. Why must you test auth after a restart?

#### Easy practical tasks

1. Run `ACL LIST` on your instance. Write the user names.
2. Run `ACL WHOAMI`. Save the reply.
3. Write four sentences: `requirepass` vs ACL users.
4. Open the ACL tutorial on redis.io. Write the `AUTH` form with a user name.

#### Medium practical tasks

1. On a lab instance, create a user that can only `PING` and `GET` keys that start with `lab:`. Prove `SET` fails. Then delete the user.
2. Set a password on the default user in a lab. Connect with `redis-cli`. Then restore open access only if you are on localhost lab.
3. Read how your client sends username and password. Write the option names.

#### Advanced practical tasks

1. Design three users: `app`, `readonly`, `ops`. Write command categories and key patterns for each.
2. Persist ACLs (`ACL SAVE` or ACL file). Restart the lab instance. Confirm users still exist. Document the files.

---

## Bind address and protected mode

`bind` sets the addresses that Redis listens on.

```text
CONFIG GET bind
```

`bind 127.0.0.1` (and `::1` for IPv6 localhost) accepts only local clients. That is the correct default for a laptop lab.

`bind 0.0.0.0` listens on all IPv4 addresses. Any host that can reach the port can attempt a connection. Combine that with no auth and you have an open data store on the network.

Protected mode: when Redis has no bind configuration that you intended and no password, Redis may refuse remote commands. The client sees an error that tells you to bind or to set a password.

```text
CONFIG GET protected-mode
```

Protected mode is a safety net. It is not a substitute for a firewall and for ACLs.

In Docker, publishing `-p 6379:6379` can expose Redis on the host. Bind inside the container and publish only to `127.0.0.1:6379` on the host when you learn:

```text
docker run --name redis-learn -p 127.0.0.1:6379:6379 redis
```

A cloud security group or firewall must allow only the application subnets. Do not rely on "nobody knows the IP."

`protected-mode no` plus `bind 0.0.0.0` plus no ACL is a misconfiguration. Do not run that on a routable network.

### Questions

#### Theoretical questions

1. What does `bind 127.0.0.1` restrict?
2. What does `bind 0.0.0.0` mean?
3. When does protected mode refuse remote work?
4. Why is protected mode not enough alone?
5. What extra risk does `docker run -p 6379:6379` add on a laptop?

#### Easy practical tasks

1. `CONFIG GET bind`. Save the value.
2. `CONFIG GET protected-mode`. Save the value.
3. Write four sentences: bind vs protected mode.
4. Write a Docker publish flag that listens only on the host loopback.

#### Medium practical tasks

1. From the same machine, connect to `127.0.0.1`. Write success or fail.
2. Read your `redis.conf` comments for `bind` and `protected-mode`. Copy the intent in your own words (six sentences).
3. Draw host, firewall, application subnet, Redis bind address.

#### Advanced practical tasks

1. On a disposable lab VM or container network, show that a remote client fails in protected mode, then fix the lab with bind to localhost only. Do not expose the instance to the public internet.
2. Write a network checklist: bind, firewall, Docker publish, security group, no public `6379`.

---

## TLS

TLS encrypts bytes on the network. Without TLS, a packet capture on the path can show commands and values. Auth passwords also travel in clear text without TLS (RESP `AUTH`).

Redis can listen with TLS (`tls-port`, certificates, and related `tls-*` settings). Exact file names and flags depend on the version. Official docs: search "Redis TLS" on redis.io.

Typical pieces:

- Server certificate and private key
- CA certificate that clients trust
- Optional client certificates (mutual TLS)
- `redis-cli --tls` and client library TLS options

Hosted Redis often requires TLS. The vendor gives a CA bundle or a "rediss://" URL.

TLS does not replace ACLs. TLS protects the path. ACLs decide what the client may do after it connects.

Certificates expire. A restart or a client fail after expiry is an operations event (topic 19). Calendar the expiry.

Do not disable certificate verification in production to "make it work." Fix the CA and the host name.

A local Docker lab on `127.0.0.1` may skip TLS. A Redis host on a shared network must not skip TLS if the network is not already a trusted encrypted fabric that your security team named as equivalent.

### Questions

#### Theoretical questions

1. What does TLS protect on the Redis path?
2. Why is `AUTH` without TLS a risk on a shared network?
3. Does TLS replace ACL users?
4. What extra flag does `redis-cli` use for TLS?
5. What happens when a certificate expires?

#### Easy practical tasks

1. Open the official Redis TLS page. Write three `tls-*` setting names.
2. Write four sentences: TLS vs ACL.
3. Find your client TLS options (CA file, skip verify). Write the names. Mark which option you must not use in production.
4. If you have a hosted Redis that requires TLS, write the connection URL scheme (`rediss` or flags). If not, write "no hosted TLS in this lab".

#### Medium practical tasks

1. Start Redis with TLS from the official example (or Redis Stack / vendor template) in a lab. Connect with `redis-cli --tls`. Record the flags. Use only lab certificates.
2. Connect once with verification on. Document the error if the CA is wrong (do not leave verification off).
3. Compare port numbers: cleartext `6379` vs `tls-port`. Write both.

#### Advanced practical tasks

1. Enable mutual TLS in a lab. Document server and client certificate files and the client flags.
2. Write an expiry runbook: who owns the cert, alert lead time, and how you reload Redis or the proxy.

---

## Dangerous commands: `FLUSHALL`, `KEYS`, `DEBUG`, `CONFIG`

Some commands are correct in a lab and harmful on a shared instance.

`FLUSHALL` deletes all keys on the instance (all logical DBs). `FLUSHDB` deletes the current logical DB. There is no undo.

`KEYS` walks the keyspace on the main thread. On a large instance it can stall Redis (topic 12, topic 16).

`DEBUG` includes subcommands that can crash the process or block it (`DEBUG SEGFAULT`, `DEBUG SLEEP`). Those are developer tools. They are not application commands.

`CONFIG GET` can reveal settings. `CONFIG SET` can change bind, persistence, and directories. A client that can `CONFIG SET` can weaken security or break persistence.

Other commands in the same class: `SHUTDOWN`, `REPLICAOF` / `SLAVEOF`, `MODULE LOAD`, `SCRIPT DEBUG`. The `@dangerous` ACL category groups many of them. Check `ACL CAT dangerous` on your version.

The application user must not have these commands. An operator uses them in a controlled window with a change ticket.

Do not run `FLUSHALL` as a deploy step on a durable Redis. Topic 13 and topic 21 explain cache warming instead.

### Questions

#### Theoretical questions

1. What does `FLUSHALL` delete?
2. Why is `KEYS` dangerous on a large instance?
3. Why is `DEBUG` not an application command?
4. What can `CONFIG SET` change that affects security?
5. What is `@dangerous` in ACLs?

#### Easy practical tasks

1. Run `ACL CAT dangerous` if the command exists. Write five command names.
2. Write four sentences: `FLUSHALL` vs `DEL` of one prefix with `SCAN`.
3. Open the `FLUSHALL` command page. Write the time complexity note.
4. Write one sentence each: `FLUSHALL`, `KEYS`, `DEBUG`, `CONFIG SET`.

#### Medium practical tasks

1. On an empty lab instance only, run `FLUSHDB` after you `SET` one key. Confirm `DBSIZE` is 0. Do not do this on shared data.
2. Compare `CONFIG GET dir` and `CONFIG GET dbfilename`. Write why those paths matter for backup and for access.
3. In your client, try to call `FLUSHALL` as a non-ops user (or write the ACL that would block it). Record the error.

#### Advanced practical tasks

1. Write an ops policy: which roles may run `CONFIG`, `FLUSHALL`, `DEBUG`, and `KEYS`. Include a ban on application use.
2. Build an ACL for `app` that includes `@read`, `@write`, and `@hash` (as needed) and excludes `@dangerous`, `KEYS`, and `FLUSHALL`. Test on lab keys.

---

## Rename or disable commands

Older setups renamed dangerous commands in `redis.conf`:

```text
rename-command FLUSHALL ""
rename-command KEYS ""
rename-command CONFIG ""
```

An empty new name disables the command. A secret name hides the command from casual use. A secret name is a weak control. The name can leak. Prefer ACLs.

Redis 7+ deprecates `rename-command` in favor of ACLs. Read the comments in your `redis.conf`. If `rename-command` is ignored or warned, use `ACL SETUSER` to remove commands.

Disable means the server rejects the command for everyone, including you, until you change config and restart (or you keep one admin user with the command). Plan how ops will still run a controlled `CONFIG` if you disable it globally.

`COMMAND` lists commands. After a disable, the name disappears or fails.

Do not rename `AUTH` or `ACL` in a way that locks you out of a remote instance. Keep a local console path (localhost, container exec) that you tested.

Application code must use `SCAN`, not a renamed `KEYS`. If you disable `KEYS`, tutorials that still show `KEYS *` will fail. That is what you want in production.

### Questions

#### Theoretical questions

1. What does `rename-command FLUSHALL ""` do?
2. Why is a secret rename weaker than an ACL?
3. What does Redis 7 prefer instead of `rename-command`?
4. What lockout risk exists if you disable `CONFIG` and `ACL` for all users?
5. Which iteration command remains after you disable `KEYS`?

#### Easy practical tasks

1. Open your `redis.conf`. Find `rename-command` comments. Write one example line.
2. Write four sentences: rename vs ACL deny.
3. Run `COMMAND INFO FLUSHALL`. Write whether the command exists on your instance.
4. List three commands you would deny for an application user.

#### Medium practical tasks

1. In a disposable lab config, disable `FLUSHALL` with the method your version documents (ACL or rename). Prove `FLUSHALL` fails. Restore the lab.
2. Document how you would still flush a lab dataset (`FLUSHDB` on a dedicated lab instance, or delete by prefix).
3. Read deprecation notes for `rename-command` on your version. Write the note in your own words.

#### Advanced practical tasks

1. Write a `redis.conf` / ACL snippet for a cache instance: app user, ops user, disabled `DEBUG`, no `KEYS` for app.
2. Test a lockout recovery plan on a lab container: how you attach to localhost and fix ACL without the network. Write the steps.

---

## Do not expose Redis to the public internet

Redis on a public address with no ACL and no TLS is a data leak and a wipe risk. Bots scan port `6379`. This is a common incident class.

Required controls (use all that apply):

- Bind to private addresses only
- Firewall or security group: application subnets only
- ACL users with least privilege
- TLS on any path that leaves a trusted host
- Protected mode left on unless you have bind and auth on purpose
- Dangerous commands denied for application users
- No public inbound `6379` / `tls-port`

A "temporary test" on `0.0.0.0` in a cloud VM is still public if the security group is open.

Managed Redis in a VPC is still not public if you follow the vendor private-endpoint model. Public endpoints need extra review (TLS, strong auth, IP allow lists).

Kubernetes: a `Service` of type `LoadBalancer` on Redis can publish it. Prefer ClusterIP and network policy.

Do not put Redis on the public internet to "save time." Use SSH tunnel, VPN, or a bastion to reach a lab in the cloud.

If you already exposed a lab, treat all keys as compromised. Rotate passwords. Take the instance off the public address. Start a new instance for learning if you are not sure.

### Questions

#### Theoretical questions

1. Why is port `6379` on a public IP a problem?
2. Name five controls from this section.
3. Why is a "short test" on `0.0.0.0` still dangerous in a cloud VM?
4. What Kubernetes Service type can publish Redis by mistake?
5. What do you do if a lab was public for an unknown time?

#### Easy practical tasks

1. Write a five-line "never public" checklist.
2. Write four sentences: private bind vs public `LoadBalancer`.
3. Check your lab: is Redis reachable only on `127.0.0.1`? Write how you tested (from the same host).
4. Open your cloud or Docker UI (if any). Write whether port `6379` is listed as public.

#### Medium practical tasks

1. Draw a recommended production path: app subnet → Redis private IP → no internet route to Redis.
2. Compare Redis Cloud "public endpoint" vs "private endpoint" in vendor docs. Write three differences.
3. Write a tunnel command you would use to reach a remote lab (SSH local forward) without opening `6379` to the world. Do not scan third-party hosts.

#### Advanced practical tasks

1. Write a team standard: bind, ACL, TLS, firewall, command denies, and a quarterly review.
2. Review an existing compose file or Helm chart for Redis. Mark public ports, missing passwords, and missing resource limits. Write findings.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do bind, firewall, TLS, and ACLs each stop a different failure (reach, sniff, act)?
2. Why is `requirepass` alone a weak production story on Redis 6+?
3. Which dangerous commands hurt availability (`KEYS`, `DEBUG SLEEP`) versus durability (`FLUSHALL`)?
4. When do you still use `rename-command`, and when do you refuse it?
5. What is the shortest correct reply to "just open 6379 so I can debug from home"?

#### Easy practical tasks

1. Write a cheat sheet: `ACL LIST`, `AUTH`, `bind`, `protected-mode`, `--tls`, `FLUSHALL`, `rename-command`.
2. Export `CONFIG GET bind`, `protected-mode`, and `aclfile` (or related). Highlight them.
3. Draw an attacker on the internet and three closed doors (firewall, TLS+auth, ACL).
4. Write lab rules: local bind, throwaway passwords, no production `FLUSHALL`.

#### Medium practical tasks

1. Harden a disposable lab: ACL app user, bind localhost, `FLUSHALL` denied. Connect as app and as ops. Record both sessions.
2. Write an incident checklist for "Redis was bound to 0.0.0.0 for a day": rotate, audit, rebuild, review firewall.
3. Map each topic 12 client setting (timeouts, password, TLS) to a security control in this topic.

#### Advanced practical tasks

1. Produce a one-page Redis security baseline for your team. Include versions (ACL, TLS), network, commands, and secrets.
2. Compare your baseline with the official Redis security page. Add any control you missed. Write the URL.
