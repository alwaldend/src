# Topology catalog

> Generated deterministic projection. The JSON document at `tools/agents/catalogs/topology.json` is authoritative.

- ID: `agent-system.topology`
- Schema: `agents.alwaldend.com/catalog/v1alpha1/topology-catalog`
- Derivation: `1.0.0`
- Producer: `repository.topology-compiler`
- Source revision: `f5f9e119945ebd65`
- Completeness: `complete`
- JSON digest: `sha256:d288102cbdd5a5a84f6bfb7a4ced8ed89d4f5b69e51b7c0661abca42e930c7af`

## Limitations

None.

## Trees

- `projects` (projects): `projects/README.md` — product
- `infra` (infra): `infra/README.md` — repository_internal
- `tools` (tools): `tools/README.md` — tool
- `data` (data): `data/README.md` — data
- `third_party` (third_party): `third_party/README.md` — third_party
- `users` (users): `users/README.md` — user

## Components

- `agents` — Agents (projects/agents, owned); lifecycle `active`
- `al` — Al (projects/al, owned); lifecycle `active`
- `alwaldend.com` — Alwaldend.com (projects/alwaldend.com, owned); lifecycle `in_progress`
- `android_launcher` — Android launcher (projects/android_launcher, owned); lifecycle `maintenance`
- `ansible_collection` — Ansible collection (projects/ansible_collection, owned); lifecycle `maintenance`
- `autoscroll` — Autoscroll (projects/autoscroll, owned); lifecycle `finished`
- `bazel_agent` — Bazel agent (projects/bazel_agent, owned); lifecycle `active`
- `cgit` — CGit (projects/cgit, owned); lifecycle `maintenance`
- `ci_platform` — Ci platform (projects/ci_platform, owned); lifecycle `abandoned`
- `dotfiles` — Dotfiles (projects/dotfiles, owned); lifecycle `maintenance`
- `goal` — Goal (projects/goal, owned); lifecycle `experimental`
- `infinitime` — Infinitime (projects/infinitime, owned); lifecycle `maintenance`
- `kustomization` — Kustomization (projects/kustomization, owned); lifecycle `in_progress`
- `leetcode_downloader` — Leetcode downloader (projects/leetcode_downloader, owned); lifecycle `in_progress`
- `mcp_cordis` — MCP Cordis (projects/mcp_cordis, owned); lifecycle `active`
- `nexus_security_plugin` — Nexus security plugin (projects/nexus_security_plugin, owned); lifecycle `finished`
- `rules_binary_toolchain` — Rules binary toolchain (projects/rules_binary_toolchain, owned); lifecycle `active`
- `rules_docs` — Rules docs (projects/rules_docs, owned); lifecycle `active`
- `rules_docs_gazelle` — Rules docs Gazelle (projects/rules_docs_gazelle, owned); lifecycle `active`
- `rules_hugo` — Rules Hugo (projects/rules_hugo, owned); lifecycle `active`
- `rules_promptfoo` — Rules Promptfoo (projects/rules_promptfoo, owned); lifecycle `active`
- `rules_promptfoo_gazelle` — Rules Promptfoo Gazelle (projects/rules_promptfoo_gazelle, owned); lifecycle `active`
- `rules_skill` — Rules skill (projects/rules_skill, owned); lifecycle `active`
- `rules_skill_gazelle` — Rules skill Gazelle (projects/rules_skill_gazelle, owned); lifecycle `active`
- `rules_template` — Rules template (projects/rules_template, owned); lifecycle `active`
- `sri` — Sri (projects/sri, owned); lifecycle `finished`
- `tf_modules` — Tf modules (projects/tf_modules, owned); lifecycle `in_progress`
- `useless_qt_gui` — Useless QT GUI (projects/useless_qt_gui, owned); lifecycle `finished`
- `xray_manager` — Xray manager (projects/xray_manager, owned); lifecycle `in_progress`

## Workspaces

- `root` — module `com_alwaldend_src` at `MODULE.bazel`
- `projects.rules_binary_toolchain` — module `rules_binary_toolchain` at `projects/rules_binary_toolchain/MODULE.bazel`
- `projects.rules_docs` — module `rules_docs` at `projects/rules_docs/MODULE.bazel`
- `projects.rules_docs_gazelle` — module `rules_docs_gazelle` at `projects/rules_docs_gazelle/MODULE.bazel`
- `projects.rules_hugo` — module `rules_hugo` at `projects/rules_hugo/MODULE.bazel`
- `projects.rules_promptfoo` — module `rules_promptfoo` at `projects/rules_promptfoo/MODULE.bazel`
- `projects.rules_promptfoo_gazelle` — module `rules_promptfoo_gazelle` at `projects/rules_promptfoo_gazelle/MODULE.bazel`
- `projects.rules_skill` — module `rules_skill` at `projects/rules_skill/MODULE.bazel`
- `projects.rules_skill_gazelle` — module `rules_skill_gazelle` at `projects/rules_skill_gazelle/MODULE.bazel`
- `projects.rules_template` — module `rules_template` at `projects/rules_template/MODULE.bazel`
