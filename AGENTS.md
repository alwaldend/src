---
title: Agents
---

## Repository guide for agents

Use this file as the repository-wide default. If a more deeply nested
`AGENTS.md` is added, its instructions take precedence for that subtree.

## Repository visibility

This is a public repository. Treat checked-in source, documentation, and test
fixtures as public information. When a task explicitly calls for an external
service, that checked-in public content may be sent to the service without a
separate confidentiality approval. Eval output derived only from checked-in
public fixtures and an isolated public workspace has the same classification.
This does not make credentials, decrypted configuration, other local or
generated artifacts, untracked files, private infrastructure values, or
secret-bearing data safe to disclose; continue to follow the repository's
secret and infrastructure handling rules.

## Scratch files (mandatory)

- Put every task-owned download, generated report, log, extracted archive,
  tool cache, and temporary file under an ignored, task-specific
  `out/<task>/` directory in the applicable workspace.
- Point configurable tool scratch directories at `out/<task>/` as well. Do
  not use `/tmp` for task-owned files merely because they are temporary.
- Use operating-system temporary storage only when a tool cannot be directed
  elsewhere. Remove any task-owned residue there before handoff; never delete
  unrelated files or broad temporary directories.

## Decision-making

Use the `decision-review` skill before committing to a material design,
security, operational, costly, irreversible, or repeatedly failing choice.
Treat both the user's proposal and your current plan as hypotheses: identify
the strongest reason a domain expert would reject them, test that reason
against evidence, and compare credible alternatives. The primary agent owns
the verdict and must not delegate it to a random subagent. Optimize for the
user's actual goal, not agreement or the cheapest satisfaction of stated
constraints. State material trade-offs to the user; routine reversible
implementation choices do not need this review or narration.

Expert guidance, audits, and adversarial review are evidence, not authorization
or new acceptance criteria. Do not turn every hypothetical weakness into
required work: classify findings against the user's stated scope, fix only
in-scope blockers, and report optional hardening separately. If an overbroad
claim would require unrelated machinery, narrow the claim instead of expanding
the implementation. Do not interrupt in-scope work for routine uncertainty or
optional hardening. If completing the requested outcome requires a material
scope expansion that existing context does not authorize, ask before acting;
otherwise choose the smallest reversible in-scope approach and continue.

## Questions

Use the `answer-question` skill whenever the user's message contains a
substantive question, including when it also requests action. When a
substantive question asks for a material decision, use both `answer-question`
and `decision-review`. A question, including “can we,” “should we,” “why,” or
“do we need,” is a request for information and not authorization to modify
files, settings, pull requests, deployments, or other external state.
Read-only investigation is allowed when it supports a truthful answer. If a
message combines a question with an explicit request for action, act only
within that stated scope.

## Making changes

- Use the `project-layout` skill before creating or moving source, or when
  deciding a directory layout anywhere in the repository.
- Read the nearest `README.md`, `BUILD.bazel`, and `include.MODULE.bazel`
  (when present) for the area being changed.
- Name projects using only ASCII letters, digits, and underscores
  (`[a-zA-Z0-9_]+`).
- Prefer a small, target-specific change. This is a large monorepo, so query,
  build, and test the affected Bazel package before considering `//...`.
- Prefer Go for repository automation and scripts. Expose them as Bazel
  `go_binary` targets and invoke them with `bazel_agent run`; use
  another language only when Go or Bazel would materially complicate the task.
- Use the `repo-delivery` skill to finalize implementation work. It owns
  staging, feature-branch commits and pushes, pull request maintenance, review
  comment handling, and the final delivery report.
- Place temporary files in the repository-root `out/` directory. Do not commit
  temporary files.
- Commit binaries only when they are required, are not temporary files, and
  are tracked by Git LFS. Do not commit binaries otherwise.

## Tooling

- Prefer purpose-built tools, MCP capabilities, and repository Bazel targets
  for work they support. Use direct host-shell commands only when no suitable
  tool, MCP, or Bazel target exists.
- Acquire tools used by Bazel hermetically: use pinned, checksummed binary
  archives or Bazel-integrated package managers driven by checked-in manifests
  and lockfiles. Do not silently depend on a host-installed tool or an
  undeclared lifecycle download.
- Hermetic tool acquisition is important; hermetic tool output has lower
  priority. A tool may intentionally access the network or produce
  environment-dependent artifacts when the requested workflow permits it.
  Continue to apply the repository's authorization, secret-handling, and
  external-side-effect rules.

## Documentation

- Document current behavior and supported guarantees first.
- Include future plans only when they are intentional and relevant to users;
  include historical behavior only when it helps operate, migrate, or
  understand the current system.
- Do not add rejected ideas, speculative alternatives, or incidental design
  exploration to ordinary maintained documentation. Keep that material in the
  task's goal or research artifacts, or in an explicitly maintained decision
  record.
- Prefer the smallest durable explanation that helps a reader use or maintain
  the system.

## Searching

Do not use recursive `grep` or `ls`. Use `rg`, `rg --files`, `find` with a
bounded depth, or `bazel_agent query` instead.

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

- Bazel is the primary entry point. Agents must invoke it through
  `bazel_agent` from the applicable workspace root. The runner delegates to the
  Bazelisk-managed `bazel`, enables batch mode, and selects the agent
  configuration. The repository `.bazeliskrc` pins both the Bazel version and
  archive hash. Do not bypass it with a separately installed Bazel binary.
- Use Bazel's filtered stdout and stderr plus targeted test logs for
  diagnostics. Do not write Build Event Protocol output unless the task
  requires it: raw BEP includes the client environment and can contain secrets.
- `MODULE.bazel` is the Bzlmod root. Most dependency families are split into
  `include.MODULE.bazel` files under `tools/`, `third_party/`, and `projects/`.
  Keep a dependency declaration with the owning subsystem rather than adding
  everything to the root module.
- `.bazelrc` imports `tools/bazelrc/root.bazelrc`, which in turn loads the
  generated preset, project flags, and optional ignored `user.bazelrc`.
  Do not put machine-local settings into checked-in rc files.
- Do not hand-edit files that identify themselves as generated. Run the update
  command in their header. Common update targets use a `.update` suffix (for
  example, `bazel_agent run //tools/ansible:requirements.update`).
- When changing BUILD or `.bzl` files, use existing macros and naming patterns
  in the same package. Run the root Buildifier test as well as package tests.
- Run `bazel_agent run //:gazelle` only when a source/dependency change requires
  generated BUILD updates, then review every generated change.

Useful discovery commands:

```sh
bazel_agent query //path/to/package:all
bazel_agent query 'tests(//path/to/package:all)'
bazel_agent query 'rdeps(//..., //path/to/package:target)'
```

`bazel_agent query` can be expensive at repository scope. Substitute the
narrowest reasonable package pattern for `//...` whenever possible.

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
bazel_agent test //path/to/affected/package:all
bazel_agent build //path/to/affected/package:all
bazel_agent test //:buildifier_test           # BUILD/.bzl changes
black --check path/to/changed/python          # Python changes
mypy path/to/changed/python                   # Python changes
bazel_agent test //...                        # only when justified/feasible
```

Not every package exposes all of these targets. Use `bazel_agent query` first
rather than guessing.

The checked-in pre-commit configuration supplies repository hygiene checks.
Install the hook with `bazel_agent run //:write_git_hooks`, and verify it with
`bazel_agent run //:write_git_hooks -- test`. Installation is optional, and
agents should still run the relevant checks explicitly.
