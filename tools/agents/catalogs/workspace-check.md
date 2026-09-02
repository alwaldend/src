# Workspace-check catalog

> Generated deterministic projection. The JSON document at `tools/agents/catalogs/workspace-check.json` is authoritative.

- ID: `agent-system.workspace-check`
- Schema: `agents.alwaldend.com/catalog/v1alpha1/workspace-check-catalog`
- Derivation: `1.0.0`
- Producer: `repository.workspace-check-compiler`
- Source revision: `87a51cd7fadb87ac72737f1ef54eef09`
- Completeness: `complete`
- JSON digest: `sha256:40d46b3b913dc607a7c55a9b0b499cd4d500fdc99a779a5b29843ef4bcb696c0`

## Limitations

None.

## Workspaces

- `projects.rules_binary_toolchain` — module `rules_binary_toolchain` at `projects/rules_binary_toolchain/MODULE.bazel`
  - bazelIgnore: true, rootOverride: true, docsAggregation: false, fullCheck: true
  - phase `projects.rules_binary_toolchain.check` via `repository.bazel-operations`: bazel_agent test //...
- `projects.rules_docs` — module `rules_docs` at `projects/rules_docs/MODULE.bazel`
  - bazelIgnore: true, rootOverride: true, docsAggregation: true, fullCheck: true
  - phase `projects.rules_docs.check` via `repository.bazel-operations`: bazel_agent test //...
- `projects.rules_docs_gazelle` — module `rules_docs_gazelle` at `projects/rules_docs_gazelle/MODULE.bazel`
  - bazelIgnore: true, rootOverride: true, docsAggregation: false, fullCheck: true
  - phase `projects.rules_docs_gazelle.check` via `repository.bazel-operations`: bazel_agent test //...
- `projects.rules_hugo` — module `rules_hugo` at `projects/rules_hugo/MODULE.bazel`
  - bazelIgnore: true, rootOverride: true, docsAggregation: false, fullCheck: true
  - phase `projects.rules_hugo.check` via `repository.bazel-operations`: bazel_agent test //...
- `projects.rules_promptfoo` — module `rules_promptfoo` at `projects/rules_promptfoo/MODULE.bazel`
  - bazelIgnore: true, rootOverride: true, docsAggregation: false, fullCheck: true
  - phase `projects.rules_promptfoo.check` via `repository.bazel-operations`: bazel_agent test //...
- `projects.rules_promptfoo_gazelle` — module `rules_promptfoo_gazelle` at `projects/rules_promptfoo_gazelle/MODULE.bazel`
  - bazelIgnore: true, rootOverride: true, docsAggregation: false, fullCheck: true
  - phase `projects.rules_promptfoo_gazelle.check` via `repository.bazel-operations`: bazel_agent test //...
- `projects.rules_skill` — module `rules_skill` at `projects/rules_skill/MODULE.bazel`
  - bazelIgnore: true, rootOverride: true, docsAggregation: false, fullCheck: true
  - phase `projects.rules_skill.check` via `repository.bazel-operations`: bazel_agent test //...
- `projects.rules_skill_gazelle` — module `rules_skill_gazelle` at `projects/rules_skill_gazelle/MODULE.bazel`
  - bazelIgnore: true, rootOverride: true, docsAggregation: false, fullCheck: true
  - phase `projects.rules_skill_gazelle.check` via `repository.bazel-operations`: bazel_agent test //...
- `projects.rules_template` — module `rules_template` at `projects/rules_template/MODULE.bazel`
  - bazelIgnore: true, rootOverride: true, docsAggregation: false, fullCheck: true
  - phase `projects.rules_template.check` via `repository.bazel-operations`: bazel_agent test //...
- `root` — module `com_alwaldend_src` at `MODULE.bazel`
  - bazelIgnore: false, rootOverride: false, docsAggregation: false, fullCheck: true
  - phase `root.check` via `repository.bazel-operations`: bazel_agent test //...
