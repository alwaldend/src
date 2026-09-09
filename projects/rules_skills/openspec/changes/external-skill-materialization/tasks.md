## 1. Archive packaging

- [x] 1.1 Add `skill_archive` and the `skill_archives` macro for archive-owned skill directories
- [x] 1.2 Fail analysis when a declared root has no SKILL.md
- [x] 1.3 Cover archive packaging with an analysis fixture
- [x] 1.4 Derive archive skill targets from a wildcard so new upstream skills need no edit

## 2. Discovery writer and check

- [x] 2.1 Add `skills_write` with `symlinks` and `archives` inputs
- [x] 2.2 Generate the updater with lock-based reconciliation
- [x] 2.3 Generate the exact-state check, including payload comparison
- [x] 2.4 Resolve the workspace from an explicit marker
- [x] 2.5 Check in real scripts and pass them data instead of generating bash in Starlark
- [x] 2.6 Rename the module and directory to `rules_skills`
- [x] 2.7 Replace `skill_discovery_links` and its tests

## 3. Consumers

- [x] 3.1 Declare the pinned OpenSpec archive under `third_party/`
- [x] 3.2 Write `.agents/skills` from `//.agents:write_skills`
- [x] 3.3 Remove the checked-in OpenSpec skill copies
- [x] 3.4 Update the registry entry, READMEs, and the openspec skill
- [x] 3.5 Accept `compatibility` in SKILL.md frontmatter

## 4. Verification

- [x] 4.1 `bazel run //.agents:write_skills`
- [x] 4.2 `bazel test //.agents:write_skills_test`
- [x] 4.3 Confirm the check fails on tampered content
- [x] 4.4 `bazel test @rules_skills//main/go:go_test`
- [x] 4.5 `openspec validate --all --strict`
- [x] 4.6 Preserve the published `rules-skill` landing identity through Terraform `moved` blocks
