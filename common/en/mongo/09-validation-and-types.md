# 9. Validation and Types

## Description

This topic shows how you validate documents in MongoDB. You attach a JSON Schema to a collection. You set `validationLevel` and `validationAction`. You compare application validation with database validation. You keep field names and types stable. Complete documents, CRUD, and modeling first.

Use one term for each concept. A **validator** is a rule on a collection. **JSON Schema** is the usual validator language. **validationLevel** selects which documents the server checks. **validationAction** selects error or warn. The application still validates user input.

---

## JSON Schema validation on collections

You can create or collMod a collection with a validator:

```javascript
db.createCollection("users", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["email", "createdAt"],
      properties: {
        email: { bsonType: "string", pattern: "^.+@.+$" },
        createdAt: { bsonType: "date" },
        age: { bsonType: "int", minimum: 0 }
      },
      additionalProperties: true
    }
  }
})
```

`$jsonSchema` describes BSON types with `bsonType`. Use `string`, `int`, `long`, `double`, `decimal`, `bool`, `date`, `objectId`, `array`, `object`, `null`, `binData`.

`required` lists fields that must exist. `properties` describes named fields. `additionalProperties: false` rejects unknown fields. That setting is strict. Use it when the schema is stable.

You can use `oneOf` for versions or variants. You can describe arrays with `items`.

Existing collections: `db.runCommand({ collMod: "users", validator: { $jsonSchema: { ... } } })`.

A validator that is too strict will block writes. Test on a copy. Insert valid and invalid documents.

The schema is not a full application model. It does not replace unique indexes. It does not replace authorization.

### Questions

#### Theoretical questions

1. What does `$jsonSchema` describe?
2. What does `required` mean?
3. What does `additionalProperties: false` do?
4. How do you add a validator to a collection that already exists?
5. Why is a unique index still required when you have a schema?

#### Easy practical tasks

1. Create `users` with a schema that requires `email` as a string. Insert a valid document.
2. Insert a document without `email`. Record the error.
3. Insert `age` as a string when the schema wants `int`. Record the error.
4. Open the JSON Schema validation page. Write the URL.

#### Medium practical tasks

1. Set `additionalProperties: false`. Insert an extra field. Record the result. Then allow extra fields.
2. Describe an array of objects (`addresses`) with `items`. Insert one valid and one invalid address.
3. Use `collMod` to add `displayName` as an optional string. Insert both old and new shapes if your schema allows it.

#### Advanced practical tasks

1. Write `oneOf` for two event types (`login`, `purchase`) with different required fields. Test two inserts.
2. Read about `$jsonSchema` vs query expressions as a validator. Write one rule that is easier as a query expression.

---

## `validationLevel` / `validationAction`

**validationLevel** controls which documents the server validates.

- `strict` — validate all inserts and all updates. This is the default.
- `moderate` — validate inserts and validate updates only if the existing document already matches the schema. Documents that were invalid before the validator can still update if they stay invalid in some cases. Read the current manual for the exact moderate rule.

Use `strict` for new collections. Use `moderate` when you add a validator to a collection that has old invalid documents and you must still update them.

**validationAction** controls what the server does on a failed check.

- `error` — reject the write. This is the default.
- `warn` — accept the write and log a warning.

Use `warn` when you introduce a schema and you want to find violations without blocking production writes. Watch the logs. Fix the application. Then set `error`.

```javascript
db.runCommand({
  collMod: "users",
  validator: { $jsonSchema: { /* ... */ } },
  validationLevel: "moderate",
  validationAction: "warn"
})
```

Do not leave `warn` forever. Invalid data will grow.

`strict` plus `error` is the usual end state for a stable collection.

### Questions

#### Theoretical questions

1. What is the default `validationLevel`?
2. When do you use `moderate`?
3. What is the difference between `error` and `warn`?
4. Why is `warn` only a temporary step?
5. Does `moderate` validate every insert?

#### Easy practical tasks

1. Set `validationAction: "error"`. Prove that a bad insert fails.
2. On a test collection, set `warn`. Insert a bad document. Confirm that it exists. Read how you would see the warning (log or Atlas).
3. Write the `collMod` command that sets both options. You can keep a simple schema.
4. Find the options page in the manual. Write the URL.

#### Medium practical tasks

1. Create invalid old documents. Add a strict validator. Try to update one old document. Record the result.
2. Change to `moderate`. Update an old invalid document (for example add a field). Record the result. Read the manual if the update fails.
3. Plan a rollout: warn → fix → error. Write the three deploy steps.

#### Advanced practical tasks

1. Compare `moderate` rules in two manual versions if they differ. Write the rule in your own words with a test.
2. Use Atlas or logs to count validation warnings for one day of test writes. Write how you found the count.

---

## Application validation vs database validation

**Application validation** runs in the client or the API. Examples: a JSON Schema library, a type decoder, Mongoose, Bean Validation. It checks input before the write. It returns a clear error to the user.

**Database validation** runs in `mongod`. It checks every writer: the API, a script, `mongosh`, a migration. It is the last guard.

Use both.

The application:

- Validates user-facing rules (password strength, email confirm)
- Gives field-level errors
- Can use richer types than the collection schema

The database:

- Blocks bad writes from tools and forgotten scripts
- Enforces type and required fields
- Does not depend on one language

Do not rely on the application only. A future job will write without that code.

Do not rely on the database only. The user will see a driver error, not a form error. The schema cannot express every business rule.

Keep the two schemas close. If `age` is `int` in MongoDB, the API must not send a float that becomes a different BSON type.

Drivers can send types that pass the application but fail `bsonType`. Test the real BSON, not only JSON.

### Questions

#### Theoretical questions

1. What does application validation do that the database does not do well?
2. What does database validation catch that the application can miss?
3. Why use both?
4. Why can a JSON API value fail `bsonType` after the driver encodes it?
5. Why is a driver error a poor user-facing message?

#### Easy practical tasks

1. Write three rules that belong in the API and three rules that belong in `$jsonSchema`.
2. Insert from `mongosh` a document that your app would reject. Show that the database still needs a validator.
3. Map one field: TypeScript or Go type → BSON type → `bsonType` string.
4. Make a two-column table: "Application" and "Database". Add four rows.

#### Medium practical tasks

1. Implement a small validator in your language for `email` and `createdAt`. Then add the same rules as `$jsonSchema`. Test both paths.
2. Find one rule that you cannot express in `$jsonSchema` (for example "email is unique"). Write which tool enforces it (unique index).
3. Show a number that JSON accepts (`1.0`) and write the BSON type your driver sends. Check if the schema allows it.

#### Advanced practical tasks

1. Generate `$jsonSchema` from an application type (or the reverse) with a tool or a short script. Write the gaps.
2. Design a policy: which team can change the collection validator, and how the API schema stays in review. Write five rules.

---

## Stable field names and types

Stable names and types keep queries, indexes, and validators working.

Rules:

- Pick a field name and keep it. Do not mix `user_id`, `userId`, and `userID`.
- Pick a BSON type and keep it. Do not store `qty` as int and as string.
- Use one date type. Prefer BSON Date, not a date string, if you sort and index by time.
- Use `Decimal128` for money. Do not mix decimal and double.
- Use one id type for a reference. Do not mix ObjectId and string for the same target.

If you must change a name or a type, use expand-contract and `schemaVersion`. Query both forms during the change. Backfill. Then remove the old form.

Indexes do not fix mixed types. A query `{ qty: { $gt: 10 } }` misses `"11"`.

Validation helps after you set it. Old mixed data still exists until you clean it.

Name style: teams often use camelCase in documents. That is a convention. Consistency matters more than the style.

Do not reuse a field name for a new meaning. Use a new name or a new version.

### Questions

#### Theoretical questions

1. Why do mixed types break range queries?
2. Why must a reference field use one id type?
3. Why is a new field name safer than a new meaning on an old name?
4. How does `schemaVersion` help a type change?
5. Why does an index not make `"11"` match `$gt: 10`?

#### Easy practical tasks

1. Insert `qty: 5` and `qty: "5"`. Query `{ qty: { $gt: 1 } }`. Write which documents match.
2. List field names in one of your collections. Mark any duplicates that differ only by case or underscore.
3. Write a one-page field dictionary: name, bsonType, required, notes.
4. Find two documents with different types for the same field (create them). Write a `$type` query that finds the strings.

#### Medium practical tasks

1. Write an expand-contract plan to change `price` from double to decimal. Include reads, writes, and backfill.
2. Add a validator that locks `qty` to `int`. Show that a string update fails.
3. Normalize `userId` that was stored as string and as ObjectId. Write the migration update.

#### Advanced practical tasks

1. Scan a collection with an aggregation that groups by `$type` of a field. Write the type histogram.
2. Design a lint check in CI that compares application structs to the collection validator. Write what the check can and cannot prove.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `$jsonSchema`, `validationLevel`, and `validationAction` work together from first rollout to a stable collection?
2. Which rules belong in the application, which belong in the validator, and which belong in a unique index?
3. Why do stable types matter for indexes, aggregation, and validation at the same time?
4. What is the risk if you set `additionalProperties: false` on day one of an unstable product?
5. A teammate says "MongoDB has no schema". Which facts do you use to correct that sentence?

#### Easy practical tasks

1. Create `products` with a small `$jsonSchema`, `strict`, and `error`. Insert one good document and one bad document.
2. Write a cheat sheet: bsonType names, required, additionalProperties, level, action, app vs db, stable names.
3. Write a field dictionary for `products` with four fields and types.
4. Run `collMod` to set `warn`, then set `error` again on a test collection.

#### Medium practical tasks

1. Add a validator to a collection that already has mixed types. Use `moderate` and `warn`. Write how you find and fix the bad documents.
2. Align one application struct with the collection schema. Fix one mismatch that you find.
3. Write `oneOf` for `schemaVersion` 1 and 2 of the same collection. Test one insert per version.

#### Advanced practical tasks

1. Produce a validation rollout report: counts of warn-log violations, top failing fields, date you switch to error.
2. Compare JSON Schema validation with client-side decoding in your language. Write a table of rules each layer enforces for a real document.
