---
title: Agents
---

## Repository guide for agents

Use this file as the repository-wide default. If a more deeply nested
`AGENTS.md` is added, its instructions take precedence for that subtree.

## Decision-making

For material design, security, and operational decisions, ask what the best
expert in that field would do and why they would reject your current choice;
if you can name that reason, don't make the choice. Optimize for what that
expert would judge correct, never for what satisfies the stated constraints
most cheaply. State material trade-offs to the user; routine implementation
choices need not be narrated.

## Making changes

- Read the nearest `README.md`, `BUILD.bazel`, and `include.MODULE.bazel`
  (when present) for the area being changed.
- Name projects using only ASCII letters, digits, and underscores
  (`[a-zA-Z0-9_]+`).
- Prefer a small, target-specific change. This is a large monorepo, so query,
  build, and test the affected Bazel package before considering `//...`.
- Use the `repo-delivery` skill to finalize implementation work. It owns
  staging, feature-branch commits and pushes, pull request maintenance, review
  comment handling, and the final delivery report.
- Do not commit binaries. Add binary paths to `.gitignore` instead.

## Searching

Do not use recursive `grep` or `ls`. Use `rg`, `rg --files`, `find` with a
bounded depth, or `bazel query` instead.

## Repository map and boundaries

- `projects/`: product and reusable project code. It may be public, published,
  and used in build actions.
- `infra/`: private infrastructure definitions (Terraform, Ansible, Flux,
  DNS, and host/service configuration). It must not be published or consumed
  by build actions.
- `tools/`: repo-wide build rules and toolchains. Except for toolchain types,
  these are private and not for production build actions; tools may be used by
  tests. Tool targets intended across the repo normally use
  `visibility = ["//:__subpackages__"]`.
- `data/`: private checked-in data and documentation assets. It may be used in
  builds, but must not be public or published.
- `third_party/`: private vendored or externally sourced code. It may be used
  in builds, but must not be public or published.
- `users/`: private user-specific code and infrastructure. It must not be
  published or consumed by builds.

The authoritative policy for each tree is its top-level `README.md`. Preserve
these visibility and publication boundaries when adding dependencies.

## Bazel and dependency management

- Bazel is the primary entry point. `.bazeliskrc` pins the supported version;
  use `bazel`, not a separately installed unpinned binary.
- `MODULE.bazel` is the Bzlmod root. Most dependency families are split into
  `include.MODULE.bazel` files under `tools/`, `third_party/`, and `projects/`.
  Keep a dependency declaration with the owning subsystem rather than adding
  everything to the root module.
- `.bazelrc` imports `tools/bazelrc/root.bazelrc`, which in turn loads the
  generated preset, project flags, and optional ignored `user.bazelrc`.
  Do not put machine-local settings into checked-in rc files.
- Do not hand-edit files that identify themselves as generated. Run the update
  command in their header. Common update targets use a `.update` suffix (for
  example, `bazel run //tools/ansible:requirements.update`).
- When changing BUILD or `.bzl` files, use existing macros and naming patterns
  in the same package. Run the root Buildifier test as well as package tests.
- Run `bazel run //:gazelle` only when a source/dependency change requires
  generated BUILD updates, then review every generated change.

Useful discovery commands:

```sh
bazel query //path/to/package:all
bazel query 'tests(//path/to/package:all)'
bazel query 'rdeps(//..., //path/to/package:target)'
```

`bazel query` can be expensive at repository scope. Substitute the narrowest
reasonable package pattern for `//...` whenever possible.

## Infrastructure safety

- Treat all files under `infra/`, `users/`, and secret-bearing `data/`
  subtrees as sensitive. Never paste credentials, private keys, state, plan
  output, inventories, or decrypted configuration into logs or commits.
- `al.lua` files use the repository's custom configuration DSL and often wire
  Vault AppRole authentication into generated commands. Follow a nearby
  service's `al.lua` and Bazel target rather than invoking tools directly.
- Terraform is wrapped by `terraform_binary_map`. Typical targets are
  `:tf.fmt_check`, `:tf.plan`, and `:tf.apply`; the first two are validation,
  while `apply`, `deploy`, `destroy`, `import`, `migrate`, `state`, and
  `force_unlock` can change remote state.
- Do not run any state-changing infrastructure target unless the user
  explicitly requests that exact operation and scope. For ordinary code
  changes, limit verification to formatting, validation, queries, builds, and
  tests that do not contact or mutate live systems.
- Do not commit `.terraform/`, state files, plans, environment files, or
  machine-local credentials. Existing `.gitignore` rules are a backstop, not
  permission to create or inspect secret material unnecessarily.

## Style

Follow `.editorconfig` and the closest existing files:

- UTF-8, LF endings, a final newline, no trailing whitespace, and a preferred
  maximum line length of 79.
- Four spaces by default (including JSON); tabs only for Go and Makefiles; two
  spaces for YAML, Markdown, HTML, Proto, and TOML.
- Python formatting and analysis settings live in `pyproject.toml`. Note that
  the project metadata supports Python 3.10+, while some tool configurations
  deliberately target Python 3.9 compatibility; do not casually synchronize
  those values.
- Never add broad formatting churn to a focused change.

## Verification

Choose checks that match the files changed, in this order:

```sh
git diff --check
bazel test //path/to/affected/package:all
bazel build //path/to/affected/package:all
bazel test //:buildifier_test                 # BUILD/.bzl changes
black --check path/to/changed/python          # Python changes
mypy path/to/changed/python                   # Python changes
bazel test //...                              # only when justified/feasible
```

Not every package exposes all of these targets. Use `bazel query` first rather
than guessing.

The checked-in pre-commit configuration supplies repository hygiene checks.
The hook itself is generated with `bazel run //:write_git_hooks`; installation
is optional, and agents should still run the relevant checks explicitly.
