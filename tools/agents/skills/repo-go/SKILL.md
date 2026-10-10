---
name: repo-go
description: >-
  Implement, refactor, and review Go code in this Bazel monorepo, including
  contextual error wrapping, subprocess boundaries, and Go dependency wiring.
  Use alongside repo-bazel and bazel-agent for builds and validation.
---

# Work with Go in this repository

Read the owning README and BUILD.bazel before editing. Keep the package's
existing style and use project-layout for new source locations; a focused
change does not authorize moving unrelated legacy packages.

## Organize code and initialization

Avoid `init()` functions. Initialize components explicitly from their owning
entry point so ordering is visible and failures can be returned with context.
Keep methods in the same source file as the declaration of their receiver type.

## Document package boundaries

Every first-party Go package must have one package documentation comment in
`doc.go`, immediately before its package declaration. Start with `Package <name>` and describe the package purpose, responsibilities, and important
boundaries or guarantees. Include command packages and handwritten wrappers
for generated code; do not edit generated files or upstream dependencies.
Keep the comment accurate when behavior changes and include `doc.go` in the
owning Bazel target. Module-wide architecture belongs in the owning README.

Document every top-level declaration, including exported and unexported types,
functions, methods, constants, and package variables. Document every struct
field, including embedded and private fields. Start each comment with the
declared name and explain its purpose or behavior; record meaningful contracts
without restating the identifier. Preserve accurate existing comments and
update them alongside changes to the declaration. Apply the same rules to
named declarations and struct fields in tests. Generated declarations rely on
their source schema or generator documentation; do not edit generated files.

## Preserve error causes

Always wrap propagated non-nil errors with useful operation and resource
context using `fmt.Errorf("operation %q: %w", resource, err)`. Do not return
`err`, `cmd.Run()`, or another fallible call directly. Check the result, wrap
its failure, and return nil on success. This applies to subprocesses,
filesystem operations, parsing, and errors from other project functions.

Use `%w` for the underlying error so callers retain `errors.Is` and
`errors.As` behavior; `%v` and string concatenation lose that cause. Create a
plain error only when the function detects a condition without an underlying
error. Keep messages precise without repeating the same context at every
layer. Never include credentials or secret-bearing command arguments.

Handle meaningful close and flush errors. A cleanup error must not hide the
primary failure; preserve both when returning both is appropriate. Log an
error at its handling boundary instead of logging and returning it at every
layer.

## Commands and dependencies

Use Cobra for Go command-line interfaces. Define subcommands with their own
flags and required-flag validation, pass the command context to operations,
and keep output and error handling at the command boundary. Reuse the existing
repository-pinned Cobra module and declare it in the owning Bazel deps.

Pass local subprocess arguments as an argv slice through `exec.CommandContext`.
Keep shell use limited to operations that need shell syntax; preserve literal
shell arguments with the existing quoting mechanism. Propagate cancellation
and connect standard input, output, and error deliberately.

Use repo-external-dependency for dependency selection, approval, and updates.
It covers inspecting existing options and obtaining approval before introducing
any new external dependency. The root
`go.mod` and `go.sum` feed the Go extension in `tools/go/include.MODULE.bazel`;
nested modules keep their own manifests.
Wire imports into the package's Bazel deps and expose newly used root modules
through the owning `use_repo` declaration. Use the owning Go command to update
module metadata and the repository workflow for generated locks.

## Build and verify

Use bazel-agent and repo-bazel for every build, test, formatter, or tool run.
Use the repository-pinned Go entry point `//tools/go:go` when a Go module
command is needed; do not substitute host Go. Inspect generator ownership
before changing generated source or BUILD declarations.

Write behavioral failure cases before implementation when isolation is
necessary. Prefer an isolated end-to-end command fixture with a repeatable
result artifact. Keep ordinary checks package-scoped, include affected
consumers, and complete repo-delivery's required gates. Offline skill eval
validation checks packaging and configuration, not agent behavior.
