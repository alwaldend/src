## 1. Remove the component

- [x] 1.1 Delete `projects/goal` (CLI, API, store, docs, diagrams, landing
      site, and disabled skill sources).

## 2. Update generated and dependent configuration

- [x] 2.1 Remove `goal` from `PROJECTS` in `projects/projects.bzl`.
- [x] 2.2 Remove the committed `goal` record from the global zone snapshot.
- [x] 2.3 Drop `projects_goal` from the repository OpenSpec validation map.
- [x] 2.4 Remove the goal row from `projects/README.md` and retarget the
      `hugo_landing` README examples.

## 3. Update specifications and history

- [x] 3.1 Withdraw the `project-goal` capability specification.
- [x] 3.2 Replace the repository baseline's goal deprecation requirement with
      a removal requirement.
- [x] 3.3 Note the removal in `migration.md` and `README.md`; update the
      workspace `config.yaml` context.
- [x] 3.4 Update `AGENTS.md` so policy states removal instead of deprecation.

## 4. Validate

- [x] 4.1 `bazel_agent bazel test //infra/src/openspec/validation:validate_test`
      passes (49 tests).
- [x] 4.2 `bazel_agent bazel test //infra/dns:config_test` passes.
- [x] 4.3 `bazel_agent bazel test //tools/repo_quality:repo_quality_test`
      passes.
- [x] 4.4 Confirm no surviving reference to `projects/goal` outside git
      history and archived OpenSpec records.
