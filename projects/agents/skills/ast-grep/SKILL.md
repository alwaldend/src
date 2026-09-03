---
name: ast-grep
description: >-
  Perform structure-aware code search and rewrite in this repository with
  ast-grep. Use when a task needs to find or transform code by syntax shape
  (pattern, AST node, semantic construct) rather than by plain text or regex;
  use repo-bazel and bazel-agent for build and test orchestration.
---

# Search and rewrite code with ast-grep

Use the repository-pinned `ast-grep` binary for searches that plain-text tools
cannot express, such as matching call sites by argument shape, function
signatures, imports, or control-flow structure, and for precise rewrites
across a language's syntax tree.

## Use the repository toolchain

Read the nearest `AGENTS.md` and the tool's `tools/ast-grep/README.md`. Read
the pinned version and archive contract from `tools/ast-grep/binary_toolchain
.json`; do not hard-code a remembered version.

Follow `bazel-agent` for every Bazel invocation. Run ast-grep through the
repository alias and pass pattern and options after Bazel's separator:

```sh
bazel_agent bazel run //tools/ast-grep:ast-grep -- --pattern 'console.log($A)' --lang javascript .
```

Do not substitute an ast-grep binary from `PATH`, a separately downloaded
binary, or an extracted ELF for repository work. The pinned toolchain is the
only supported entry point. Verify the pattern against the actual source:
structural matching can silently miss code when the pattern's implied AST
shape does not match the real syntax.

## Match by syntax shape

ast-grep patterns describe AST structure with metavariables like `$A`,
`$FUNC`, and wildcards `$$$` for zero or more nodes. Match only the semantic
shape you intend; an over-broad pattern can match unrelated code.

- Use `--lang` to select the language and `--json` or `--interactive` when
  inspecting matches.
- Use `--rewrite` and `--edit` for transforms; review a diff (`--diff`) before
  applying changes to tracked sources.
- Scope searches to the owning project or package when the task is local.

## Keep rewrites safe

Prefer read-only search until the user authorizes edits. For rewrites, stage
the change on a dedicated feature branch in its own worktree, review the
generated diff, and run the owning package's tests after the transform. Never
rewrite generated files, vendored code, or lockfiles with ast-grep.

## Validate patterns

When a pattern is ambiguous, confirm it on a small checked-in fixture first
with `--pattern` before running it at repository scope, and use `--json` to
inspect the matched nodes. Some languages are not supported by the pinned
ast-grep build; prefer the native language tooling or a focused search when
ast-grep cannot parse the file.
