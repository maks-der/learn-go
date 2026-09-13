# 2. Language Basics

## Description

A Go source file has a fixed order. The file starts with a package clause. Import declarations come next. Other declarations come after the imports.

This topic explains that file skeleton. You learn packages, import paths, and `func main`. You also learn comments, names, export rules, and semicolons.

The compiler checks these rules first. A program does not run when the file structure is wrong. Complete this topic before you study types.

Use one term for each concept. A package is a directory of Go files with the same package name. A command is a `package main` that contains `func main`. An import path is the string in an `import` declaration.

---

## Packages and `package main`

A package is the unit that the compiler builds. All `.go` files in one directory form one package. Those files must use the same package name. A file starts with a package clause:

```go
package greet
```

The package name is an identifier. The package name is often the same as the folder name. The package name is not the import path. Other files import the package by import path. Code inside the package uses the package name as a prefix only from outside.

`package main` is special. A `package main` defines a command. A command builds to a binary. A library package uses a different name, such as `greet` or `http`. A library package does not build to a binary by itself.

The program entry sits in `package main`. A module can contain many packages. Only a `package main` with `func main` starts a process. Put commands in a folder such as `cmd/app` when the module has more than one command.

Two files in one directory must not use different package names. The compiler reports an error. Test files may use `package name` or `package name_test`. The `name_test` form is an external test package. External tests import the package like a client.

```go
package main

import "fmt"

func main() {
	fmt.Println("command")
}
```

### Questions

#### Theoretical questions

1. What is a package in Go?
2. Why must all files in one directory use the same package name?
3. What does `package main` mark?
4. What is the difference between a package name and an import path?
5. When does `go build` write a binary for a package?

#### Easy practical tasks

1. Create a module. Add `package main` and print `hello`. Run the program with `go run .`.
2. Add a folder `msg` with `package msg`. Export one function. Call it from `main`.
3. Put two files in one folder with different package names. Record the compiler error. Fix the names.
4. Draw a folder tree with `package main` in `cmd/app` and `package store` in `store`. Write each import path.

#### Medium practical tasks

1. Split one command into three files in the same `package main`. Keep one `func main`. Build the command.
2. Add `package msg_test` in folder `msg`. Import `msg` and call one exported function from a test.
3. Create two commands under `cmd/alpha` and `cmd/beta`. Build each command from the module root.

#### Advanced practical tasks

1. Design a module with one library package and two commands that import it. Write the tree, package names, and import paths. Build both commands.
2. Move a library package to a new folder name that does not match the package name. Update imports. Document what still works and what confuses readers.

---

## `import` and import paths

An import declaration names a package that this file uses. The import path is a string. For the standard library, the import path is the path after `src`, such as `"fmt"` or `"net/http"`. For a module, the import path is the module path plus the package folder.

```go
package main

import (
	"fmt"
	"os"

	"example.com/hello/msg"
)
```

Group imports in one `import` block. `gofmt` puts standard library paths in the first group. `gofmt` puts other paths in the next group. A blank line separates the groups.

The compiler rejects an unused import. Remove the import or use a name from that package. `goimports` can add and remove imports for you.

An import binds the package name in this file. The default name is the package clause of the imported package. Change the name with an alias when two packages share a name:

```go
import (
	fmt1 "fmt"
	myhttp "net/http"
)
```

A blank import uses `_` as the name. The file does not use names from that package. The compiler still runs the imported package `init` functions. Use a blank import only for a required side effect.

A dot import uses `.` as the name. Names from that package enter the file block without a prefix. Do not use a dot import in normal code. The prefix shows the source of each name.

An import cycle is an error. Package A must not import B if B imports A, directly or through other packages.

### Questions

#### Theoretical questions

1. What is an import path?
2. Why does the compiler reject an unused import?
3. What does a blank import do?
4. When do you use an import alias?
5. What is an import cycle?

#### Easy practical tasks

1. Write a program that imports `fmt` and `strings`. Print the result of `strings.ToUpper("go")`.
2. Add an unused import. Run `go build`. Record the error. Remove the import.
3. Import `fmt` with the alias `out`. Call `out.Println`. Run the program.
4. Run `go list -f "{{.ImportPath}}" ./...` in a module with two packages. Write each import path.

#### Medium practical tasks

1. Create packages `a` and `b` that import each other. Record the cycle error. Break the cycle with a third package.
2. Use a blank import of a small package that prints a line in `init`. Show the print when `main` runs.
3. Add one standard library import and one local import. Run `go fmt`. Show how `gofmt` groups the paths.

#### Advanced practical tasks

1. Import the same package twice with two aliases in one file. Record the compiler result. Explain the rule in three sentences.
2. Build a package that other code must not import from another module. Place it under `internal/`. Import it from `main` in the same module. Try an import from a second module and record the error.

---

## `func main()` as the program entry point

The program entry point is `func main` in `package main`. The signature has no parameters and no results:

```go
func main() {
}
```

The runtime starts the process. The runtime initializes the `main` package. The runtime runs `init` functions first. The runtime then calls `func main`. When `main` returns, the process exits with status `0`.

Do not give `main` parameters. Do not give `main` results. Read command-line arguments from `os.Args`. Call `os.Exit` when you need a non-zero status. `os.Exit` stops the process at once. Deferred calls do not run after `os.Exit`.

A `package main` must contain exactly one `func main`. Two `main` functions in the same package are an error. A library package must not define the program entry. You may name a function `main` in a library package. That function is not the process entry.

`func main` can call other functions in the same package or in imported packages. Keep `main` small. Put logic in other functions and packages.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: app name")
		os.Exit(1)
	}
	fmt.Println(os.Args[1])
}
```

### Questions

#### Theoretical questions

1. What is the exact signature of the program entry function?
2. In which package must the program entry live?
3. What happens when `func main` returns?
4. How does `os.Exit` differ from a return from `main`?
5. Can a library package start a process? Explain.

#### Easy practical tasks

1. Write `func main` that prints `ready`. Run it with `go run .`.
2. Print `len(os.Args)` and each argument. Run the binary with two extra words.
3. Add a second `func main` in another file of the same package. Record the error. Remove the extra function.
4. Move `func main` into a library package. Run `go run .` from that folder. Record what happens.

#### Medium practical tasks

1. Write `main` that returns status `1` with `os.Exit` when no argument exists. Return status `0` when an argument exists. Show both runs.
2. Add an `init` function that prints `init`. Let `main` print `main`. Show the order of the two lines.
3. Split setup and work into two functions. Keep `main` as three or fewer calls. Run the program.

#### Advanced practical tasks

1. Compare `os.Exit(1)` with `panic("fail")` from `main`. Record process status and whether a `defer` in `main` runs.
2. Create two `package main` folders in one module. Give each its own `func main`. Build each path. Explain why both can exist.

---

## Comments (`//`, `/* */`)

A comment is text that the compiler ignores. Go has two comment forms. A line comment starts with `//` and ends at the end of the line. A block comment starts with `/*` and ends with `*/`.

```go
// Package greet formats short messages.
package greet

// Hello returns a greeting for name.
func Hello(name string) string {
	/* This block comment
	   spans two lines. */
	return "hi " + name
}
```

Use `//` for almost all comments. Use `/* */` for a long comment or to disable a block of code for a short test. Block comments do not nest. A `/*` inside a block comment does not start a new comment.

A package comment sits above the package clause with no blank line. Package comments document the package on pkg.go.dev. A function comment sits above the function. Start a function comment with the function name. The same rule applies to types, variables, and constants.

The compiler also reads some comments as directives. Example: `//go:build` selects files at build time. Place a directive at the start of the file. Follow the official format. A normal comment must not look like a directive by accident.

Do not write comments that only repeat the code. Write comments that state intent or a constraint. Keep comments in the same language as the project.

### Questions

#### Theoretical questions

1. Where does a `//` comment end?
2. Can you nest `/* */` comments?
3. Where do you place a package comment?
4. What is a `//go:build` comment?
5. What must the first words of a function comment include?

#### Easy practical tasks

1. Add a package comment and a function comment to a small command. Run `go doc` in that folder.
2. Write one `//` comment and one `/* */` comment in the same file. Run `go run .`.
3. Place a comment between `package` and `import` that is not a directive. Confirm that the program still builds.
4. Disable one statement with `/* */`. Run the program. Restore the statement.

#### Medium practical tasks

1. Write comments for two exported functions so that `go doc -all` shows both comments.
2. Add a `//go:build ignore` file next to `main.go`. Confirm that `go build` skips that file.
3. Write a nested `/* /* */ */` comment and record the compiler error. Replace it with two `//` lines.

#### Advanced practical tasks

1. Document a package with a package comment that is several sentences. Publish nothing. Show `go doc` output and explain each sentence.
2. Compare `//go:build` with a file suffix such as `_windows.go`. Write a short note on when you use each form.

---

## Identifiers, keywords, and naming conventions

An identifier names a package, function, type, variable, constant, or label. An identifier starts with a letter or `_`. The next characters are letters, digits, or `_`. A letter includes `_` and Unicode letters.

These names are valid: `total`, `Total`, `HTTPClient`, `_tmp`, `x2`. These names are not valid: `2x`, `user-name`, `package`.

A keyword is a reserved word. You must not use a keyword as an identifier. The keywords include `package`, `import`, `func`, `var`, `const`, `type`, `if`, `else`, `for`, `range`, `switch`, `case`, `break`, `continue`, `goto`, `defer`, `return`, `go`, `select`, `chan`, `map`, `struct`, `interface`, `fallthrough`, and `default`.

The blank identifier `_` discards a value. Use `_` when a function returns a value that you do not need.

Name packages with short lowercase words. Do not use mixed case in a package name. Do not use `_` in a package name unless you have a strong reason.

Name other identifiers with MixedCaps. Start an exported name with an uppercase letter. Start an unexported name with a lowercase letter. Keep acronyms in full caps in exported names, such as `HTTP` and `URL`.

Use short names in a small scope. `i` is fine in a short loop. Use a longer name in a wide scope. Do not repeat the package name in every function name.

```go
package counter

func Add(total int, n int) int {
	return total + n
}
```

### Questions

#### Theoretical questions

1. What characters may start an identifier?
2. Why is `package` not a valid identifier?
3. What is the blank identifier?
4. What is the MixedCaps rule?
5. How do you write an exported name that contains `HTTP`?

#### Easy practical tasks

1. Declare five valid identifiers and use them in one program. Print one value.
2. Try to name a variable `func`. Record the compiler error. Rename the variable.
3. Name a package `userinfo` and a folder `userinfo`. Export `Count`. Call `userinfo.Count` from `main`.
4. List ten keywords from a file that you write. Next to each keyword, write one word that describes its role.

#### Medium practical tasks

1. Rename a function from `get_user_url` to a MixedCaps name. Update all calls. Run `go test` or `go run .`.
2. Use `_` to ignore the index from a `for range` loop over a slice of names. Print only the names.
3. Create two identifiers that differ only by case, such as `total` and `Total`. Print both. Write when this is a bad idea.

#### Advanced practical tasks

1. Read Effective Go on names. Change one of your packages to match those rules. Record each old name and each new name.
2. Write a small scanner that reads a `.go` file as text and counts identifier-like words. This is a learning tool. Compare the count with the number of names you expect.

---

## Exported vs unexported names (capital letter)

Export is a visibility rule. A name is exported when the first letter is an uppercase Unicode letter. Other packages can use an exported name. A name is unexported when the first letter is a lowercase letter or `_`. Other packages cannot use an unexported name.

The rule applies to functions, types, variables, constants, fields, and methods. The rule does not use a keyword such as `public`. The case of the first letter is the only switch.

```go
package msg

var Default = "hi" // exported
var suffix = "!"   // unexported

func Hello(name string) string { // exported
	return wrap(name)
}

func wrap(name string) string { // unexported
	return Default + " " + name + suffix
}
```

Code in the same package can use exported names and unexported names. Code in another package can use only exported names. An external test package `name_test` is another package. External tests see only exported names.

Export the names that form the package API. Keep helper names unexported. Do not export a name only because a test in another package wants it. Put such a test in the same package, or export a small test helper.

`go doc` lists exported names. Unexported names stay hidden in the public document.

### Questions

#### Theoretical questions

1. What makes a name exported?
2. Can another package call an unexported function?
3. Does the export rule apply to struct fields?
4. Why is there no `public` keyword in Go?
5. What names does an external test package see?

#### Easy practical tasks

1. Export `Hello` from package `greet`. Call `greet.Hello` from `main`. Run the program.
2. Change `Hello` to `hello`. Build `main`. Record the error. Restore the export.
3. Add an unexported helper in `greet`. Call it from `Hello` only. Confirm that `main` cannot call the helper.
4. Export a package-level variable and print it from `main`. Then unexport it and record the new error.

#### Medium practical tasks

1. Create a struct with one exported field and one unexported field. Set both fields inside the package. Try to set the unexported field from `main` and record the error.
2. Write tests in `package greet` that call an unexported function. Write tests in `package greet_test` that cannot call it. Show both results.
3. Export only two names from a package that has five functions. Use `go doc` to show the public API.

#### Advanced practical tasks

1. Design a package API with three exported names and four unexported names. Write a one-page note that explains each export choice.
2. Break a program by exporting a name that you later need to change. Then hide the name and add a new exported function. Document the compatibility lesson.

---

## Semicolons and how the compiler inserts them

The Go grammar uses semicolons between statements. You almost never type a semicolon. The lexer inserts a semicolon when a line ends after a token that can end a statement.

The lexer inserts a semicolon after a newline when the last token is an identifier, a number, a string, `break`, `continue`, `fallthrough`, `return`, `++`, `--`, `)`, `]`, or `}`.

This rule is why the opening `{` of `func`, `if`, and `for` stays on the same line:

```go
func main() {
	x := 1
	if x > 0 {
		return
	}
}
```

If you move `{` to the next line after `if x > 0`, the lexer inserts a semicolon after `0`. The compiler then reads `if x > 0;` and a block. That form is not valid for `if`.

A return value must stay on the same line as `return`, or you must use parentheses:

```go
return x

return (
	x
)
```

Do not write `return` and then `x` on the next line. The lexer inserts a semicolon after `return`. The function returns without the value.

`gofmt` removes visible semicolons that you do not need. You may write `a := 1; b := 2` on one line. Prefer one statement per line.

### Questions

#### Theoretical questions

1. Who inserts most semicolons in Go source?
2. After which tokens does the lexer insert a semicolon at a newline?
3. Why must `{` stay on the same line as `func` or `if`?
4. What goes wrong when `return` is the last token on a line and the value is on the next line?
5. Does `gofmt` require you to type semicolons?

#### Easy practical tasks

1. Write a valid `func main` with `{` on the same line. Run it.
2. Move the `{` of `if` to the next line. Record the compiler error. Restore the brace.
3. Write two statements on one line with a semicolon. Run `go fmt`. Show the file after format.
4. Write `return` and a value on one line. Then split them across two lines and record the error or the wrong result.

#### Medium practical tasks

1. Build a list of five tokens that trigger semicolon insertion. For each token, write one tiny file that shows the effect.
2. Write a function with a multi-result `return` that uses parentheses across lines. Run the function.
3. Put `++` at the end of a line. Show that the next line is a new statement. Print the value after the increment.

#### Advanced practical tasks

1. Read the semicolon rule in the language specification. Rewrite the rule in ten short sentences of your own. Do not copy the specification text.
2. Find a real compiler error in your code that came from semicolon insertion. Fix the layout. Explain the token that caused the insert.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the required order of a Go source file from the first line to the first function.
2. How do package names, import paths, and exported names work together when `main` calls another package?
3. Why does a blank import still compile when the file uses no name from that package?
4. What visibility do comments have for `go doc` compared with unexported functions?
5. How does semicolon insertion interact with MixedCaps names at the end of a line?

#### Easy practical tasks

1. Write a three-file command: `main.go`, `help.go` in `package main`, and `text/format.go` as a library. Print one formatted line.
2. Save `go doc` output for the library package in a text file. Add one line that names each unexported helper that `go doc` hides.
3. Add a second library package that imports the first. Call only exported names across that chain from `main`.
4. Run `go fmt` on the module. Confirm that imports are in groups and that no extra semicolons remain.

#### Medium practical tasks

1. Create a small CLI that uses `os.Args`, one imported library package, and `os.Exit`. Cover a missing-argument path and a success path.
2. Add an unused import, a bad brace layout, and a lowercase call to an intended API. Fix the three compiler errors. Record each message.
3. Write an external test that uses only exported names and an internal test that uses one unexported helper. Run `go test ./...`.

#### Advanced practical tasks

1. Build a module with `cmd/tool`, `internal/parse`, and `internal/parse/token`. Document every package name and import path. Show that `main` cannot import a name that you keep unexported.
2. Write a one-page map of one real file from the standard library. Mark the package clause, import groups, exported names, comments, and any line where a semicolon is inserted.
