# 8. Modules and Packages

## Description

A module is a set of Go packages with a `go.mod` file. A package is one directory of Go files that use the same `package` name. This topic shows how you create a module, how you add and tidy dependencies, and how you version a module for a breaking change. This topic also shows internal packages, build tags, and file embedding.

Complete topic 1 before this topic. Topic 1 introduces `go.mod` and module mode. This topic gives the commands and the design rules. Use one term for each idea. Use module for the versioned unit. Use package for one importable directory. Do not mix GOPATH mode with module mode.

---

## `go mod init`, `go get`, and `go mod tidy`

Run `go mod init` in the empty project root. Pass the module path. The module path is the import prefix for every package in the module:

```text
go mod init example.com/app
```

The command writes `go.mod`. A new file looks like this:

```text
module example.com/app

go 1.22.0
```

The `module` line is the path. Other modules import your packages with this prefix. The `go` line is the language version that the module expects.

Pick a module path that you control. For public code, use a repository URL without a scheme: `github.com/user/app`. Do not use a random short name like `app` when other modules must import the code.

`go get` adds or updates a module in `go.mod`. Pass a module path and a version:

```text
go get example.com/pkg@v1.2.3
go get example.com/pkg@latest
go get example.com/pkg@none
```

From Go 1.16, `go get` changes the module file. `go get` does not install a command binary. Install a command with `go install example.com/cmd/tool@version`.

`go mod tidy` adds missing modules that the source imports. `go mod tidy` removes modules that the source does not import. `go mod tidy` updates `go.sum`. Run `go mod tidy` after you add or delete imports.

`go.sum` stores cryptographic checksums. The go command uses `go.sum` to detect a changed module zip. Do not edit `go.sum` by hand. Restore it with `go mod tidy` when it is wrong.

`go mod vendor` copies the required modules into `vendor/` at the module root. Use `vendor/` when a team requires offline builds. Do not use `vendor/` as a substitute for `go.mod`.

The go command also records indirect modules in `go.mod`. An indirect module is a dependency of a dependency. The comment `// indirect` marks those lines. Do not delete those lines by hand. Run `go mod tidy`.

Keep the exported names of each package few. Give the package one job. Prefer the standard library when it solves the problem. Avoid import cycles.

### Questions

#### Theoretical questions

1. What file does `go mod init` create?
2. What is a module path?
3. What does `go get module@version` change?
4. What two actions does `go mod tidy` perform on requirements?
5. What is the role of `go.sum`?

#### Easy practical tasks

1. Create an empty folder. Run `go mod init example.com/learn/mod`. Show the full contents of `go.mod`.
2. In a new module, `go get rsc.io/quote@v1.5.2` (or another small public module). Show the `require` line in `go.mod`.
3. Import the package, print a value from it, and run `go run .`.
4. Run `go mod tidy`. Show that `go.sum` exists and is not empty.

#### Medium practical tasks

1. Create a module with path `github.com/you/demo`. Add package `greet` in a subfolder. Import it from `main` with the full module prefix. Build the module.
2. Run `go list -m all`. Explain the difference between the main module line and a dependency line.
3. Break one checksum in `go.sum`. Run `go build`. Record the error. Restore the file with `go mod tidy`.

#### Advanced practical tasks

1. Pin one dependency to an older version with `go get`. Then upgrade to `@latest`. Record both `go.mod` lines and a successful build for each version.
2. Run `go mod vendor`. List `vendor/` and build with `-mod=vendor`. Write when vendor mode is the only mode that works.

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

Some old repositories have tags like `v2.0.0` and no `go.mod` from that time. Those modules can appear with a `+incompatible` suffix. New modules must not rely on that pattern.

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

## `internal` packages

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

## Build tags and file suffixes

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

The old form `// +build` still works. New files must use `//go:build`.

Pass a custom tag with `-tags`:

```text
go test -tags=integration ./...
```

```go
//go:build integration

package store_test
```

Use tags for optional integration tests or for optional features. Do not use tags to hide broken files. Fix the files.

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
4. Write a file with `//go:build !windows`. Build on your system. Record whether the file is included.

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

//go:embed version.txt
var version string
```

Import `"embed"` when you use `embed.FS`. Place a `//go:embed` comment on the line above a variable:

```go
import "embed"

//go:embed static/*
var staticFS embed.FS
```

The variable type must be `string`, `[]byte`, or `embed.FS`. The path is relative to the package directory. The path must not contain `.` or `..` as a path element. You can list several patterns on one directive or on several `//go:embed` lines above one `embed.FS` variable.

`embed.FS` implements `fs.FS`. You can open files with `staticFS.Open("static/app.css")`. You can pass `embed.FS` to `http.FileServer` with `http.FS`.

The embed happens at compile time. A change to the file on disk does not change a running binary. Rebuild after you edit an embedded file.

Use embed for HTML templates, SQL files, small static assets, and a version string. Do not embed huge data sets without a reason. The binary grows by the file size.

All files in the pattern must exist at compile time. A missing file is a compile error. Do not embed secrets that you must rotate without a rebuild.

### Questions

#### Theoretical questions

1. When does the compiler read an embedded file?
2. Which variable types accept `//go:embed`?
3. Why do you import `_ "embed"` for a `string` variable?
4. What filesystem interface does `embed.FS` implement?
5. Why must you rebuild after you edit an embedded file?

#### Easy practical tasks

1. Create `version.txt` with one line. Embed it into a `string`. Print the string from `main`.
2. Embed the same file into a `[]byte`. Print `len` of the slice.
3. Create `static/hello.txt`. Embed the folder into `embed.FS`. Read the file with `ReadFile` and print it.
4. Change `version.txt`. Run the old binary if you still have it. Then rebuild and show the new text.

#### Medium practical tasks

1. Serve `embed.FS` with `http.FileServer` and `http.FS`. Open the URL in a browser or with `curl`.
2. Embed two glob patterns into one `embed.FS`. List the names with `fs.WalkDir`.
3. Try `//go:embed ../secret.txt`. Record the compiler error. Move the file into the package directory.

#### Advanced practical tasks

1. Embed HTML templates and parse them with `html/template` from `embed.FS`. Render one page.
2. Compare embed with `os.ReadFile` at start. Write six sentences on deploy, size, and when you must not embed.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a module path, a package import path, and an `internal` fence work together in one repository?
2. When do you run `go mod init`, `go get`, and `go mod tidy` in a new project that later adds a dependency?
3. Why does a breaking API change need `/v2` in the module path?
4. How do file suffixes and `//go:build` select source for one `GOOS`?
5. What belongs in `embed` and what belongs on disk next to the process?

#### Easy practical tasks

1. Create a module `example.com/pack`. Add `internal/store` and `cmd/app`. Embed a `help.txt` string. Build and run.
2. Write a cheat sheet: `go mod init`, `go get`, `go mod tidy`, `go list -m`, `//go:build`, `//go:embed`.
3. Draw a tree with `go.mod`, `internal/`, one `*_windows.go` file, and one embedded file.
4. Run `go list -m` and `go env GOMOD`. Write both outputs.

#### Medium practical tasks

1. Add a public module with `replace`, then `go get` a real version later. Record `go.mod` before and after.
2. Split one command into a public package and `internal` implementation. Add a GOOS-specific file for a path helper.
3. Write an integration test with `//go:build integration` that reads an embedded fixture file.

#### Advanced practical tasks

1. Design a v1 and v2 of the same library on disk. Import both from one command with two `replace` lines. Call both APIs.
2. Cross-compile the command for two operating systems. Show which Go files each build includes and confirm the embedded file is in both binaries.
