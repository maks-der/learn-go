# 9. Packages and Modules

## Description

A module is a set of Go packages with a `go.mod` file. A package is one directory of Go files that use the same `package` name. This topic shows how you create a module, how you add and tidy dependencies, and how you version a module for a breaking change. This topic also shows internal packages, small package APIs, build tags, and file embedding.

Complete topic 1 before this topic. Topic 1 introduces `go.mod` and module mode. This topic gives the commands and the design rules. Use one term for each idea. Use *module* for the versioned unit. Use *package* for one importable directory. Do not mix GOPATH mode with module mode.

---

## Creating a module with `go mod init`

Run `go mod init` in the empty project root. Pass the module path. The module path is the import prefix for every package in the module:

```text
go mod init example.com/app
```

The command writes `go.mod`. A new file looks like this:

```text
module example.com/app

go 1.22.0
```

The `module` line is the path. Other modules import your packages with this prefix. Example: package `internal/store` has import path `example.com/app/internal/store`. The `go` line is the language version that the module expects.

Pick a module path that you control. For public code, use a repository URL without a scheme: `github.com/user/app`. For private experiments, a domain that you own is enough: `example.com/learn/mod`. Do not use a random short name like `app` when other modules must import the code.

`go mod init` does not create Go source files. Add packages after the command. Place `go.mod` at the root. Do not place a second `go.mod` in a subfolder unless you intend a separate module.

From Go 1.21 the toolchain can add a `toolchain` line. That line selects a Go toolchain version. You can set the language version later:

```text
go mod edit -go=1.22.0
```

Run `go list -m` to print the module path. Run `go env GOMOD` to print the path of the active `go.mod` file.

### Questions

#### Theoretical questions

1. What file does `go mod init` create?
2. What is a module path?
3. Why must a public module path match a repository location that you control?
4. Where must `go.mod` live in a normal module?
5. What does the `go` line in `go.mod` record?

#### Easy practical tasks

1. Create an empty folder. Run `go mod init example.com/learn/mod`. Show the full contents of `go.mod`.
2. Run `go list -m` and `go env GOMOD`. Write both outputs.
3. Change the language version with `go mod edit -go=1.22.0`. Show `go.mod` again.
4. Add `main.go` with `package main` in the same folder. Run `go run .`.

#### Medium practical tasks

1. Create a module with path `github.com/you/demo`. Add package `greet` in a subfolder. Import it from `main` with the full module prefix. Build the module.
2. Run `go mod init` in a folder that already has `go.mod`. Record the error. Use a new folder instead.
3. Create two folders, each with its own `go.mod`. Explain in four sentences when two modules are correct and when one module is correct.

#### Advanced practical tasks

1. Initialize a module, add a `toolchain` line with `go mod edit`, and write what `go version` and `go env GOTOOLCHAIN` show after the edit.
2. Map a local module path to a future public path. Write the `go.mod` `module` line, the import paths, and the git remote URL that you plan to use.

---

## `go get`, `go mod tidy`, `go mod vendor`

`go get` adds or updates a module in `go.mod`. Pass a module path and a version:

```text
go get example.com/pkg@v1.2.3
go get example.com/pkg@latest
```

Use `@none` to remove a requirement:

```text
go get example.com/pkg@none
```

From Go 1.16, `go get` changes the module file. `go get` does not install a command binary. Install a command with `go install example.com/cmd/tool@version`.

`go get ./...` looks at the imports in the current module and updates `go.mod` to match. Prefer `go mod tidy` for day-to-day cleanup.

`go mod tidy` adds missing modules that the source imports. `go mod tidy` removes modules that the source does not import. `go mod tidy` updates `go.sum`. Run `go mod tidy` after you add or delete imports.

`go.sum` stores cryptographic checksums. The go command uses `go.sum` to detect a changed module zip. Do not edit `go.sum` by hand. Restore it with `go mod tidy` when it is wrong.

`go mod vendor` copies the required modules into `vendor/` at the module root. The folder `vendor/modules.txt` lists the copied modules. Build with vendor mode when you need a build without the module proxy:

```text
go mod vendor
go build -mod=vendor ./...
```

Use `vendor/` when a team requires offline builds or a fixed copy of source in the repository. Do not use `vendor/` as a substitute for `go.mod`. Keep `go.mod` and `go.sum` in every module.

The go command also records indirect modules in `go.mod`. An indirect module is a dependency of a dependency. The comment `// indirect` marks those lines. Do not delete those lines by hand. Run `go mod tidy`.

### Questions

#### Theoretical questions

1. What does `go get module@version` change?
2. How do you install a command binary from a module in Go 1.16 and later?
3. What two actions does `go mod tidy` perform on requirements?
4. What is the role of `go.sum`?
5. When does a team copy modules into `vendor/`?

#### Easy practical tasks

1. In a new module, `go get rsc.io/quote@v1.5.2` (or another small public module). Show the `require` line in `go.mod`.
2. Import the package, print a value from it, and run `go run .`.
3. Run `go mod tidy`. Show that `go.sum` exists and is not empty.
4. Run `go get rsc.io/quote@none` after you remove the import. Show that the `require` line is gone.

#### Medium practical tasks

1. Run `go mod vendor`. List `vendor/` and read the first lines of `vendor/modules.txt`. Build with `-mod=vendor`.
2. Run `go list -m all`. Explain the difference between the main module line and a dependency line.
3. Break one checksum in `go.sum`. Run `go test ./...` or `go build`. Record the error. Restore the file with `go mod tidy`.

#### Advanced practical tasks

1. Pin one dependency to an older version with `go get`. Then upgrade to `@latest`. Record both `go.mod` lines and a successful build for each version.
2. Compare a build with the default module mode and a build with `-mod=vendor` after you delete the module cache path for that module (or use a machine without network). Write when vendor mode is the only mode that works.

---

## Semantic import versioning (`/v2`)

Go modules use semantic versions: `vMAJOR.MINOR.PATCH`. A major version `v2` or higher is a different module from `v1`. The import path must include the major version suffix:

```text
module example.com/note/v2
```

Importers write:

```go
import "example.com/note/v2"
```

Major version `v0` and `v1` have no suffix. These module paths are valid:

- `example.com/note` for `v0.x` and `v1.x`
- `example.com/note/v2` for `v2.x`
- `example.com/note/v3` for `v3.x`

The rule keeps one import path mapped to one major API. A program can import `example.com/note` and `example.com/note/v2` in the same build when both are needed during a migration.

A breaking change in a tagged `v1` or higher module requires a new major version. A breaking change is a change that can fail an existing importer at compile time or that changes documented behavior in an incompatible way. Adding a new exported function is not a breaking change. Removing an exported function is a breaking change.

Do not put `v2` code in module path `example.com/note` without the `/v2` suffix. The go command rejects that layout for modules with a `go.mod` file.

Versions `v0.x` make no compatibility promise. Importers must treat `v0` as unstable.

Some old repositories have tags like `v2.0.0` and no `go.mod` from that time. Those modules can appear with a `+incompatible` suffix. New modules must not rely on that pattern. Add a correct `go.mod` with `/v2` for new major versions.

### Questions

#### Theoretical questions

1. When must a module path end with `/v2`?
2. Why can one program import `example.com/note` and `example.com/note/v2` together?
3. What change requires a new major version?
4. Does `v1` use a `/v1` suffix in the module path?
5. What does a `v0` major version mean for compatibility?

#### Easy practical tasks

1. Write a `go.mod` that declares `module example.com/learn/note/v2` and `go 1.22.0`. Add `package note` and an exported function.
2. In a second module, import `example.com/learn/note/v2` with a `replace` directive that points at the first folder. Call the function.
3. Make a two-column table: "Version" and "Module path". Add rows for `v0.3.0`, `v1.2.0`, `v2.0.0`, and `v3.1.4`.
4. Read `go help modules` or the module reference. Write one sentence about the import compatibility rule.

#### Medium practical tasks

1. Start a `v1` module. Copy it to a `v2` module with a changed function signature. Import both from a third module with `replace`. Call both APIs.
2. Tag a mental release plan: list three changes that stay in `v1` and two changes that need `v2`. Write one reason for each.
3. Try `module example.com/learn/note` with import path `example.com/learn/note/v2` in the same module. Record the go command error. Fix the `module` line.

#### Advanced practical tasks

1. Read the official blog post on v2 modules. Write a six-step checklist that a maintainer follows to publish `v2` next to `v1`.
2. Explain `+incompatible` in four short sentences. State why a new module must not use that pattern.

---

## Internal packages (`internal/`)

A directory named `internal` limits who can import the packages under it. Only the parent of `internal` and the tree under that parent can import those packages.

Example layout:

```text
example.com/app/go.mod
example.com/app/cmd/tool/main.go
example.com/app/internal/store/store.go
example.com/app/pkg/api/api.go
```

Package `example.com/app/cmd/tool` can import `example.com/app/internal/store`. Package `example.com/app/pkg/api` can import it. A different module cannot import `example.com/app/internal/store`. A sibling module path such as `example.com/other` cannot import it.

The last `internal` element in the path is the fence. In `example.com/app/internal/store/sql` the package `sql` is still internal to `example.com/app`.

Use `internal/` for code that is not part of the public API. Put shared implementation there. Put the small public surface in a package that is not under `internal/`.

Do not use `internal/` as the only design tool. A clear public package still matters. `internal/` only blocks imports from outside the parent tree.

The go command enforces the rule at compile time. The error mentions `internal`.

### Questions

#### Theoretical questions

1. Which packages may import `example.com/app/internal/store`?
2. Can another module import your `internal` package?
3. What is the "parent of `internal`" in the import path?
4. Why do teams put implementation code under `internal/`?
5. Does `internal/` replace the need for a small public API?

#### Easy practical tasks

1. Create `internal/store` with an exported function. Call it from `package main` at the module root. Build.
2. Draw a tree with `cmd/app`, `internal/store`, and `go.mod`. Write the import path of `store`.
3. Write four sentences that explain who can import `internal` and who cannot.
4. Run `go list ./...` in a module that has `internal/`. Confirm that the internal package appears in the list.

#### Medium practical tasks

1. Create two modules. Try to import `internal` from module A in module B. Record the compiler error. Then move the needed function to a public package and import that.
2. Place `internal` under `pkg/api/internal/db`. Show that `pkg/api` can import `db` and that `cmd/tool` cannot. Record both results.
3. Split a module into `internal/impl` and a public package that calls `impl`. Keep the public function list short.

#### Advanced practical tasks

1. Design a module with two commands and one internal library. Write the folder tree, every import path, and which imports the go command allows.
2. Find an `internal` package in the standard library tree on your machine (`GOROOT`). Write the import path and one package that is allowed to import it.

---

## Package design: small APIs, few dependencies

A package is a unit of API and a unit of compilation. Keep the exported names few. Export a type or a function only when an importer outside the package needs it. Hide the rest with a lowercase name or with `internal/`.

Give the package one job. Package `fmt` formats text. Package `errors` handles errors. Do not build a package named `util` that mixes files, HTTP, and random helpers. Split that code.

Keep the dependency list short. Each `require` in `go.mod` is a contract that you must update and trust. Prefer the standard library when it solves the problem. Add a module when the standard library is not enough and the extra module is stable.

Avoid import cycles. Package A must not import B if B imports A. Cycles are compile errors. Fix a cycle with a smaller third package or with an interface that breaks the cycle.

Accept interfaces and return concrete types when that guideline fits (topic 7). Do not export a large interface that only one type implements unless you need it for tests.

Document every exported name. The first sentence is the summary. Run `go doc` on your package and read it as a new user.

Name the package after what it provides, not after what it contains. Prefer `package store` over `package storehelpers`. The import path already has the module prefix.

### Questions

#### Theoretical questions

1. When must a name stay unexported?
2. Why is a package named `util` often a design problem?
3. Why must a module keep few dependencies?
4. What happens when two packages import each other?
5. What does a reader see first in `go doc` for an exported function?

#### Easy practical tasks

1. Write a package with two exported functions and two unexported helpers. Call only the exported functions from `main`.
2. List the exported names of `io` with `go doc io`. Count them. Write why the set is small.
3. Open `go.mod` of a small public module. Count the `require` lines. Write one sentence about the size.
4. Rename a package from `helpers` to a name that states the job. Update the import.

#### Medium practical tasks

1. Split one file that formats text, writes a file, and computes a hash into three packages. Keep dependencies one-way.
2. Create an import cycle on purpose. Record the error. Break the cycle with an interface in a fourth small package.
3. Review five exported names in one of your packages. Unexport every name that `main` does not need.

#### Advanced practical tasks

1. Write a one-page design for a `config` package: exported types, unexported types, standard library imports, and zero extra modules. Implement the public API only.
2. Compare `net/http` and a large third-party router. Write six short sentences about API size, dependencies, and when you stay with the standard library.

---

## Build tags / file suffixes (`_windows.go`, `//go:build`)

The go command includes a file only when the file matches the current build context. The context includes `GOOS`, `GOARCH`, and custom tags.

A file suffix sets `GOOS` or `GOARCH` without a comment:

- `foo_windows.go` builds only on Windows
- `foo_linux.go` builds only on Linux
- `foo_amd64.go` builds only on `amd64`
- `foo_windows_amd64.go` builds only on Windows `amd64`

Do not use a suffix that is not a known `GOOS` or `GOARCH`. The suffix `_test.go` marks a test file. That suffix is not a GOOS name.

A build constraint is a comment at the top of the file. Use the `//go:build` form. Put a blank line after the constraint. Then write the `package` line:

```go
//go:build windows

package pathx
```

The expression can use `&&`, `||`, and `!`:

```go
//go:build linux && amd64
```

```go
//go:build !windows
```

The old form `// +build` still works. New files must use `//go:build`. The go command can write both forms when you run `gofmt` on a file that has only one form.

Pass a custom tag with `-tags`:

```text
go test -tags=integration ./...
```

```go
//go:build integration

package store_test
```

Use tags for optional integration tests or for optional features. Do not use tags to hide broken files. Fix the files.

`cgo` and `race` are special tags. You do not set them as file suffixes.

Two files in one package must not declare the same function for the same build context. Put the shared signature in a file without a tag. Put the operating system body in `foo_windows.go` and `foo_unix.go`.

### Questions

#### Theoretical questions

1. Which operating system includes `net_windows.go`?
2. What must follow a `//go:build` line before the `package` line?
3. How do you enable a custom tag named `integration`?
4. Why must new files use `//go:build` instead of only `// +build`?
5. Why can two files in one package both declare `func open()` if they have different GOOS suffixes?

#### Easy practical tasks

1. Add `hello_windows.go` and `hello_linux.go` (or `hello_darwin.go`) that each set a string. Print the string from `main`. Build on your operating system.
2. Add `//go:build ignore` to a file. Confirm that `go build` skips the file.
3. Run `go env GOOS GOARCH`. Write the file suffix that matches your machine.
4. Write a file with `//go:build !windows`. Build on Windows and on a Unix-like system if you can. Record which build includes the file.

#### Medium practical tasks

1. Implement `func UserCacheDir() (string, error)` in `cache_windows.go` and `cache_unix.go` with `//go:build` on the Unix file. Call it from `main`.
2. Add an integration test file with `//go:build integration`. Show that `go test` skips it and `go test -tags=integration` runs it.
3. Put `//go:build` on the first line without a blank line before `package`. Record the result. Insert the blank line and fix the file.

#### Advanced practical tasks

1. Cross-compile with `GOOS=linux GOARCH=amd64` and with `GOOS=windows GOARCH=amd64`. Show which tagged files each build compiles (`go list -f '{{.GoFiles}}'`).
2. Design a package that uses a file suffix for GOOS and a custom tag for an optional driver. Write the file names, the `//go:build` lines, and the build commands.

---

## Embedding files with `embed`

Package `embed` includes files in the binary at compile time. The compiler reads files from the package directory. You do not open those files from disk at run time.

Import the package. A blank import is enough when you embed into a `string` or a `[]byte`:

```go
import _ "embed"
```

Import `"embed"` when you use `embed.FS`.

Place a `//go:embed` comment on the line above a variable:

```go
//go:embed version.txt
var version string
```

```go
//go:embed logo.png
var logo []byte
```

```go
//go:embed static/*
var static embed.FS
```

The variable must be a `string`, a `[]byte`, or an `embed.FS`. The pattern is relative to the directory of the Go file. The pattern must not contain `..` and must not start with `/`.

The compiler skips files whose names start with `.` or `_` unless you use the `all:` prefix (Go 1.18 and later):

```go
//go:embed all:templates
var templates embed.FS
```

Embed is correct for HTML templates, SQL files, and small static assets. The data lives in the binary. A large asset increases the binary size.

The file must exist at compile time. A missing file is a compile error. Use `embed` for files that ship with the program. Use `os.ReadFile` for files that the user supplies.

You cannot embed a file that sits outside the package directory. Move the file or place the embed in the package that owns the file.

### Questions

#### Theoretical questions

1. When does the compiler read an embedded file?
2. Which types may a `//go:embed` variable have?
3. Why must the embed path stay inside the package directory?
4. What does the `all:` prefix change?
5. When must you use `os.ReadFile` instead of `embed`?

#### Easy practical tasks

1. Create `hello.txt` next to a Go file. Embed it into a `string`. Print the string.
2. Embed the same file into a `[]byte`. Print `len` of the slice.
3. Embed a folder of two files into `embed.FS`. Read one file with `fs.ReadFile`.
4. Run `go doc embed`. Write the purpose of `embed.FS` in one sentence.

#### Medium practical tasks

1. Embed a `text/template` file. Parse it with `template.ParseFS` and execute it.
2. Omit a file that `//go:embed` names. Record the compile error. Restore the file.
3. Embed a directory without `all:`. Add a file named `.hidden`. Show whether the FS contains it. Repeat with `all:`.

#### Advanced practical tasks

1. Build a tiny static file server that uses `http.FileServer` with `http.FS` and an `embed.FS`. Serve two files. You do not need a public port if you use `httptest`.
2. Compare binary size before and after you embed a 1 MB file. Write the two sizes and one sentence about when embed is the wrong tool.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the full path from an empty folder to a module that imports a v2 dependency, hides implementation in `internal/`, and embeds a config template.
2. How do `go get`, `go mod tidy`, and `go mod vendor` work together in one release workflow?
3. Why does a breaking API change force both a new module path and a new import path?
4. How do file suffixes and `//go:build` solve the same problem in different ways?
5. A teammate adds a second `go.mod` in `internal/` and exports a large `util` package. Which rules from this topic do you use in a review?

#### Easy practical tasks

1. Create `example.com/learn/box`. Add `internal/store`, a public `greet` package, `//go:embed hello.txt`, and `main`. Format, tidy, and build.
2. Write a one-page cheat sheet: `go mod init`, `go get`, `go mod tidy`, `go mod vendor`, `/v2`, `internal/`, `//go:build`, `//go:embed`.
3. Draw a module tree with `cmd/app`, `internal/engine`, `go.mod`, `vendor/`, and one `*_windows.go` file. Label what each part is for.
4. Run `go list -m all` and `go list ./...` in that module. Write what each command lists.

#### Medium practical tasks

1. Write a script that runs `go mod tidy`, `go test ./...`, and `go build -mod=readonly .` (or `-mod=vendor` after `go mod vendor`). Stop on the first failure.
2. Publish a local `v1` and `v2` of the same library with `replace` in a consumer module. Call one function from each major version in `main`.
3. Add a Windows file and a Unix file for a path helper, plus an `internal` package that both files may use. Build for your `GOOS` and list `GoFiles`.

#### Advanced practical tasks

1. Create a module that vendors one external dependency, embeds a `go:embed` file system, and uses `-tags=slow` for a long test. Document the commands that a CI system must run.
2. Read the module reference on go.dev. Write a report of one page: module path rules, `/v2`, `replace`, `exclude`, and when you use `vendor/`. Implement only the parts that this topic already uses (`replace`, `/v2`, `vendor`).
