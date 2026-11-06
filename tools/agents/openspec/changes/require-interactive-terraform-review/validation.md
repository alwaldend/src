## Verification

- `git diff --check`: passed.
- OpenSpec strict validation: passed.
- `//tools/agents/skills/repo-infra:eval_config_test`: passed offline.
- `//.agents:write_skills_test`: passed.
- `bazel build --config=lint` for the skill and owning OpenSpec source: passed;
  the skill validation aspect also passed.
- `//:repo_quality_test`: all 27 tests passed before final preparation; the
  prepared candidate is checked again through the required delivery gates.

Offline Promptfoo validation checks configuration and staging, not whether a
model follows the policy. No live Terraform command was run. The representative
output is the packaged skill, inspected against text-mode review, interactive
approval, partial-failure recovery, and incompatible-target scenarios.

## Artifact audit

The read-only audit inventories registered repository workspaces, their
existing linked Bazel outputs, `/tmp`, and `/var/tmp`. It records only paths,
sizes, permissions, owners, timestamps, and classifications. Binary detection
uses ZIP member names without extracting contents. Likely JSON and text plan
files are classified locally without emitting values. No unrelated files are
removed, and no plan payload is retained as evidence.

The task-private report is `out/terraform-plan-policy/plan-audit.json` in the
task worktree. It includes scan bounds, inaccessible paths, and continuation
state. Its findings are observations at the recorded time, not evidence that
arbitrarily named exports, excluded dependencies, symlinks, or inaccessible
directories are clear.

## Review and publication

The entire task diff is reviewed with the pinned `review-agent` procedure.
Exact comparison refs, final candidate identity, review findings, and delivery
receipts remain under `out/terraform-plan-policy/`. Component entry points
that require saved plans remain incompatible with the new agent policy;
their existing guards and reusable Terraform APIs are preserved.
