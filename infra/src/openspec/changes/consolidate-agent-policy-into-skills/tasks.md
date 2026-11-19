## 1. Consolidate the root policy

- [x] 1.1 Reduce `AGENTS.md` to repository-wide constraints and one
      "When to load a skill" routing table ordered by task phase.
- [x] 1.2 Make initial reads task-dependent; require BUILD and MODULE
      inspection only when implementation or validation depends on it.
- [x] 1.3 Replace the absolute one-copy rule with one authoritative source,
      distinguishing duplicates from generated projections and fixtures.
- [x] 1.4 Simplify the retry rule to its principle and remove the fixed
      attempt counts and cache specifics.
- [x] 1.5 Fix the generated-file ambiguity: exclude disposable build outputs
      and include required generator-maintained files.

## 2. Add the workspace skill

- [x] 2.1 Add `projects/agents/skills/repo-workspace` with its `SKILL.md`,
      `BUILD.bazel`, `agents/openai.yaml`, and offline eval.
- [x] 2.2 Register its `:skill` label in `.agents/BUILD.bazel` and regenerate
      the discovery links.

## 3. Transfer extracted facts

- [x] 3.1 Move commit-subject and trailer conventions into
      `tools/repo_delivery/skills/repo-delivery/references/commits.md` and
      link it from the delivery skill.
- [x] 3.2 Add the Go and Bazel-native automation rules to `repo-bazel`.
- [x] 3.3 Confirm the nested `.bazelrc` rule, infrastructure packaging,
      secret-file handling, and verification sequence each have an owner.

## 4. Reconcile drift

- [x] 4.1 Align `git-rebase-remote` and `repo-delivery` on preserving a
      replaced commit's reachable progress while permitting rewritten
      ancestry.

## 5. Validate

- [x] 5.1 `bazel_agent bazel test //:repo_quality_test` passes.
- [x] 5.2 `bazel_agent bazel test //:buildifier_test` passes.
- [x] 5.3 `bazel_agent bazel test //.agents:write_skills_test` passes.
- [x] 5.4 `bazel_agent bazel build` of the changed skill packages succeeds and
      their `eval_config_test` targets pass.
- [x] 5.5 `OPENSPEC_PROJECT=infra/src bazel_agent bazel run //tools/openspec --
validate consolidate-agent-policy-into-skills --type change --strict
--no-interactive` passes.
- [x] 5.6 Run the ergonomics fixtures against the revision and record routing,
      authority, and read-cost verdicts.
