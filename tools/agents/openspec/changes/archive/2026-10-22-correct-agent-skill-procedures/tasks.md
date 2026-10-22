## 1. Correct the procedures

- [x] 1.1 Correct Android and Codex safety boundaries; verify the instructions preserve data and require declarative host deployment.
- [x] 1.2 Correct runner examples; compare them with the runner's accepted command grammar.
- [x] 1.3 Preserve authorized implementation through OpenSpec planning; add an eval scenario for an explicit fix request.
- [x] 1.4 Correct blog output and Hugo draft guidance; compare with the owning BUILD declaration.

## 2. Validate and deliver

- [x] 2.1 Validate skill packages, offline eval configurations, discovery, formatting, and repository quality against the candidate.
- [x] 2.2 Review the aggregate diff and select the delivery scope, aggregate message, and validation targets; use the delivery receipt for subsequent lint, publication, and pull-request verification state.

## Evidence

- The preflight ran the eight affected `eval_config_test` targets,
  `//.agents:write_skills_test`, and `//:repo_quality_test`: all 35 tests passed.
- Runner examples were checked against `bazelArguments` in
  `projects/bazel_agent/cmd/bazel_agent/main.go`; site guidance was checked
  against the draft flag selection and `out_dir` in
  `projects/alwaldend.com/BUILD.bazel`.
- Offline eval validation proves packaging and configuration, not model behavior.
- Receipt-bound final checks and publication continue through `repo-delivery`;
  ignored logs and receipts remain under `out/agent-skill-corrections/`.
