# 1. Getting Started

## Description

Go is a compiled programming language. The Go team designed Go for simple syntax, fast compilation, and concurrent programs. This topic shows how you install Go, how you organize a module, and how you use the basic Go commands. Complete this topic before you write Go programs.

Use one term for each concept. Do not mix the old GOPATH mode with Go modules. Use Go modules for all new work.

---

## What Go is and why it exists

Go is a statically typed language. The compiler checks types before the program runs. Go compiles source files to a native binary. The binary runs on the target operating system without a virtual machine.

The Go team made Go to solve problems in large server software. Compilation of large C++ programs was slow. Threads in many languages were complex. Dependency management was difficult. Go gives a small language, a fast compiler, and goroutines for concurrent work.

Go is simple. The language has few keywords. The standard library covers files, HTTP, JSON, and tests. Go is concurrent. The `go` keyword starts a goroutine. Go is compiled. You distribute one binary.

Go is not an object-oriented language in the class-and-inheritance style. Go uses structs and interfaces. Go is not a dynamic scripting language. You declare types. You handle errors with return values.

### Questions

#### Theoretical questions

1. What does "compiled language" mean for a Go program?
2. Why did the Go team add goroutines to the language?
3. How is Go different from a language that uses a virtual machine?
4. Does Go use class inheritance? Explain the Go approach in one sentence.
5. Name three design goals of Go.

#### Easy practical tasks

1. Write five sentences that describe Go. Use only facts from this section.
2. Make a two-column table: "Go property" and "What it means". Add four rows.
3. List three program types that fit Go (for example, a command-line tool). List two program types that do not fit well. Give one reason for each choice.
4. Open [https://go.dev](https://go.dev). Write the current stable Go version that the site shows.

#### Medium practical tasks

1. Compare Go with one language that you know. Write six short sentences. Cover compilation, types, concurrency, and errors.
2. Draw a simple diagram of the path from a `.go` file to a running binary. Label compile and run.
3. Find the Go release notes for the last two major versions. Write three changes that help a beginner.

#### Advanced practical tasks

1. Read the "History of Go" material on go.dev. Write a one-page timeline with years and one fact per year.
2. Explain why a fast compiler helps a large team. Give a numeric example (build time before and after). Use public data or a measured local build.

---

## Installing Go and checking `go version`

Download the official installer from [https://go.dev/dl/](https://go.dev/dl/). Select the package for your operating system. Run the installer. Accept the default install path if you do not have a special requirement.

After the installation, open a new terminal. A new terminal loads the new `PATH`. Run:

```text
go version
```

The output must show `go version go1.x.y` and your operating system and architecture. Example: `go version go1.22.0 windows/amd64`.

If the command is not found, the `PATH` does not include the Go `bin` directory. On Windows, the default is `C:\Program Files\Go\bin`. On macOS and Linux, the default is `/usr/local/go/bin`. Add that directory to `PATH`. Open a new terminal. Run `go version` again.

Also run:

```text
go env GOROOT GOPATH GOMODCACHE
```

`GOROOT` is the install location of the toolchain. `GOPATH` is the workspace for downloaded modules and caches. You do not put your project source in `GOPATH` for module mode.

### Questions

#### Theoretical questions

1. Why must you open a new terminal after you install Go?
2. What information does `go version` show?
3. What is `GOROOT`?
4. What is `GOPATH` in module mode?
5. What is the cause when the terminal shows "command not found" or "not recognized" for `go`?

#### Easy practical tasks

1. Install Go if it is not installed. Run `go version`. Save the full output in a text file.
2. Run `go env`. Find `GOOS` and `GOARCH`. Write their values.
3. Run `go env GOPATH`. Open that folder in the file manager. List the top-level names that you see (`pkg`, `bin`, or empty).
4. Run `go help env`. Write the purpose of `GOROOT` and `GOPATH` in your own words.

#### Medium practical tasks

1. Install two Go versions with a version manager (`goenv` or the official method with extra toolchains). Show `go version` for each. Document the commands.
2. Change `GOPATH` in the user environment. Run `go env GOPATH`. Confirm the new value. Restore the original value.
3. On Windows, show `where.exe go`. On macOS or Linux, show `which go`. Explain how the shell finds `go`.

#### Advanced practical tasks

1. Install Go from source with the bootstrap process in the official documentation. Record each command and the result of `go version`.
2. Compare the official installer with a package manager (`winget`, `apt`, `brew`). Write a short report: version freshness, `GOROOT` path, and update method.

---

## Go modules (do not use GOPATH mode)

Old Go projects used GOPATH mode. All source lived under `$GOPATH/src`. Import paths were directory paths under that tree. That model did not record versions of dependencies in the project.

Go modules are the current standard. A module is a set of packages with a `go.mod` file. The `go.mod` file records the module path and the Go version. The `go.mod` file records required module versions. The `go.sum` file records checksums for those modules.

Always use modules. Run `go mod init` in the project root. Do not set `GO111MODULE=off`. Do not put new project source inside `GOPATH/src` for that reason.

The module path is the import prefix. Example: `github.com/user/app`. Packages in the module use that prefix in import statements from other modules.

```text
go mod init example.com/hello
```

A new `go.mod` looks like this:

```text
module example.com/hello

go 1.22.0
```

### Questions

#### Theoretical questions

1. What file marks a directory as a Go module?
2. What problem did GOPATH mode have with dependency versions?
3. What is a module path?
4. What is the role of `go.sum`?
5. Why must a new project use modules?

#### Easy practical tasks

1. Create an empty folder. Run `go mod init example.com/hello`. Show the contents of `go.mod`.
2. Change the Go version line in `go.mod` with `go mod edit -go=1.22`. Show `go.mod` again.
3. Explain in four sentences the difference between a module and a package.
4. Search `go help modules`. Write the first command that a new project needs.

#### Medium practical tasks

1. Create two modules on disk. Make module B import a package from module A with a `replace` directive in `go.mod`. Build module B.
2. Run `go list -m all` in a module with one dependency. Explain each column of the output.
3. Enable a vendor folder with `go mod vendor`. List what `vendor/` contains. State when a team uses `vendor/`.

#### Advanced practical tasks

1. Reproduce a GOPATH-mode layout in a separate folder (for learning only). Document why module mode is better. Do not use GOPATH mode for real work.
2. Break `go.sum` on purpose (change one checksum). Run `go test`. Record the error. Restore `go.sum` with `go mod tidy`.

---

## Workspace layout: `go.mod`, packages, files

A typical module has this layout:

- `go.mod` at the root
- one or more `.go` files
- optional folders that are extra packages

A package is a directory of Go files with the same `package` name. All files in one directory must use the same package name. Tests use `package name` or `package name_test`.

File names use letters, digits, and `_`. The compiler ignores files with names that start with `.` or `_`. The compiler ignores files with build tags that do not match. Files named `*_windows.go` apply to Windows. Files named `*_test.go` are test files.

Do not put many unrelated packages in one directory. One directory is one package. Put commands in `cmd/toolname` when the module has more than one command. Put private code in `internal/` when other modules must not import it.

```text
example.com/app/
  go.mod
  cmd/app/main.go
  internal/store/store.go
```

### Questions

#### Theoretical questions

1. Where does `go.mod` live in a module?
2. Can two files in one directory use different `package` names? Why?
3. What files does the compiler ignore by name prefix?
4. What is the difference between a package and a module?
5. Why does `cmd/mytool` exist in many repositories?

#### Easy practical tasks

1. Create `hello.go` with `package main` and `func main`. Keep `go.mod` in the same folder. Run the program.
2. Add a subfolder `greet` with `package greet` and an exported function. Call it from `main`.
3. Add a file `_notes.go` with valid Go code. Run `go build`. Confirm that the compiler ignores the file.
4. Draw a tree of a module with `cmd/app`, `internal/store`, and `go.mod`.

#### Medium practical tasks

1. Split one program into three packages: `cmd/app`, `internal/mathx`, and `internal/printx`. Build from the module root.
2. Create two files in one directory with different package names. Record the compiler error. Fix the error.
3. Add `README.md` and a `.go` file. Confirm that `go build` does not compile the markdown file.

#### Advanced practical tasks

1. Design a layout for a module with two commands and one shared library package. Write the folder tree and the import paths.
2. Use a `//go:build` tag on a file. Build with and without the tag. Show which file is included.

---

## `go run`, `go build`, `go test`, `go fmt`, `go vet`

`go run` compiles and runs a package in one step. Use `go run .` in the module directory for `package main`. Use `go run file.go` only for small single-file programs.

`go build` compiles the package. For `package main`, `go build` writes a binary in the current directory. The binary name comes from the directory name. Use `-o name` to set the output name. For a library package, `go build` checks compilation and does not write a file in the current directory.

`go test` compiles test files and runs tests. Test files have the `_test.go` suffix.

`go fmt` formats Go files in the package. The format is the official Go format. Do not debate indentation. Run `go fmt ./...` for the full module.

`go vet` reports suspicious code. Examples: printf format errors, unreachable code, incorrect locks. `go vet` is not a full linter. Run `go vet ./...` with tests.

Typical sequence for a change:

1. Write code.
2. Run `go fmt ./...`.
3. Run `go vet ./...`.
4. Run `go test ./...`.
5. Run `go build -o app .` for a command.

### Questions

#### Theoretical questions

1. What is the difference between `go run` and `go build`?
2. When does `go build` write a binary file?
3. What files does `go test` use?
4. Why does the Go team require one format?
5. What class of errors does `go vet` find that the compiler can miss?

#### Easy practical tasks

1. Write `package main` that prints `ok`. Run it with `go run .`.
2. Build a binary with `go build -o hello.exe .` (or `hello` on Unix). Run the binary.
3. Change the indentation of a file. Run `go fmt`. Show the file after format.
4. Write a `fmt.Printf` call with a wrong format verb. Run `go vet`. Record the message.

#### Medium practical tasks

1. Add `hello_test.go` with a test that always passes. Run `go test -v`.
2. Time `go run` versus `go build` plus run of the binary for the same program. Write the two times.
3. Run `go test ./...` from the module root with two packages. Confirm both packages run.

#### Advanced practical tasks

1. Use `go build -gcflags=-m` on a small program. Read one escape analysis line. Write what it means.
2. Add a `go vet` failure and a test failure in the same package. Fix both. Show the commands that you ran.

---

## Editor setup and official docs

Configure the editor to run format on save. Install the Go extension for Visual Studio Code or the Go plugin for GoLand. Set the format tool to `gofmt` or `goimports`. `gofmt` is the formatter in the Go toolchain. `go fmt` calls `gofmt`. `goimports` formats the file and fixes import lists.

Enable `gopls`. `gopls` is the Go language server. `gopls` gives completion, find-definition, and diagnostics. Point the editor to the same `go` binary that the terminal uses.

Do not mix tabs and spaces by hand. `gofmt` uses tabs for indentation.

The primary documentation is [https://go.dev/doc/](https://go.dev/doc/). It contains the specification, Effective Go, the tutorial, and release notes.

The Tour of Go is [https://go.dev/tour/](https://go.dev/tour/). The tour is an interactive set of lessons. Complete the tour in the browser. The tour covers types, methods, interfaces, and goroutines.

Package documentation is on [https://pkg.go.dev/](https://pkg.go.dev/). For the standard library, open [https://pkg.go.dev/std](https://pkg.go.dev/std). Each package page shows types, functions, and examples.

Use `go doc` in the terminal:

```text
go doc fmt.Println
go doc net/http.Client
```

Use the language specification when you need exact rules. Use Effective Go when you need style. Use the tour when you learn the first time.

### Questions

#### Theoretical questions

1. What extra work does `goimports` do compared with `gofmt`?
2. What is `gopls`?
3. What is the Tour of Go?
4. Where do you read package documentation for `net/http`?
5. When do you open the language specification instead of a tutorial?

#### Easy practical tasks

1. Install the Go extension in your editor. Enable format on save. Save a badly formatted file. Confirm the format change.
2. Run `gofmt -d file.go` on a file with extra spaces. Read the diff.
3. Complete the first two pages of the Tour of Go. Write one fact that you learned.
4. Run `go doc fmt`. Write the package comment in one sentence.

#### Medium practical tasks

1. Complete all Tour of Go pages in the "Basics" section. Write five quiz questions for yourself from that section.
2. Use `go doc -src fmt.Println`. Read the source. Write how `Println` calls `Fprintln`.
3. Add a format check in a script: `gofmt -l .` must print no file names. Fail the script if a name appears.

#### Advanced practical tasks

1. Complete the full Tour of Go. Make a list of topics that this handbook covers later. Map each tour page to a topic number.
2. Configure `golangci-lint` with a `gofmt` check. Run it on a small module. Record the output.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the full path from an empty folder to a running Hello World program. Name each command.
2. Why is a checksum file part of a module?
3. How do `go fmt` and `go vet` work together in a daily workflow?
4. What is the difference between documentation on pkg.go.dev and documentation from `go doc`?
5. A teammate uses GOPATH mode. Which facts do you use to explain the move to modules?

#### Easy practical tasks

1. Create a new module `example.com/start`. Add `main.go` that prints `started`. Format, vet, test (add one test), and build.
2. Write a one-page cheat sheet with these commands: `go version`, `go env`, `go mod init`, `go run`, `go build`, `go test`, `go fmt`, `go vet`, `go doc`.
3. Export `go env` to a file. Highlight `GOOS`, `GOARCH`, `GOROOT`, `GOPATH`, and `GOMODCACHE`.
4. Create a folder tree for a future CLI tool. Add only `go.mod` and empty packages. Run `go list ./...`.

#### Medium practical tasks

1. Write a small script (PowerShell or bash) that runs `go fmt`, `go vet`, and `go test` in sequence. Stop on the first failure.
2. Clone a small public Go repository. Identify `go.mod`, the main package, and one library package. Run `go test ./...`.
3. Document your editor and Go setup in ten steps so that another beginner can copy it.

#### Advanced practical tasks

1. Build the same `main` package for `GOOS=linux` and `GOOS=windows`. Show the two binary names and file sizes. You do not need to run the foreign binary.
2. Create a module that uses one external module from pkg.go.dev. Record `go get`, the `go.mod` require line, and a successful `go run`.
