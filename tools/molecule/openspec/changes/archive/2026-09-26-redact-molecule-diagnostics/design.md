## Context

[PR #112 review](https://github.com/alwaldend/src/pull/112#discussion_r4111430016)
found that `collect` retained rendered Ansible headings and the callback stored
`get_name()`. `no_log` does not make these headings safe. Existing specifications
already require sanitized artifacts and credential isolation.

## Goals / Non-Goals

Prevent accidental disclosure through rendered names while preserving phase,
action, outcome, assertion classification, and aggregate counters. This does not
change VM lifecycle behavior or introduce a general secret-value redactor.

## Decisions

Drop rendered play, task, and handler headings altogether. Strictly recognize
recap lines for the runner's two inventory names with numeric counters only.
Callback events use Ansible's generated task UUID instead of the rendered name;
action and failure classification still identify the failed assertion.
Matching or redacting known secret values would miss arbitrary interpolations.

Before changing production code, extend the real guest fixture with templated
play/task/handler names and `no_log`, using a public sentinel. Extend artifact
scanning and assert that failed actions remain identifiable. The focused
`TestFailureCleansResources` run failed on the retained verification log and
callback events, confirming both disclosure paths. Red-run evidence is retained
under ignored `out/molecule-runner-split/diagnostic-red-evidence`.

## Risks / Trade-offs

Human-readable task names are no longer retained; correlate task UUIDs, action
names, phases, and outcomes instead. Private raw logs remain temporary and are
removed with the run. Validate success and failure using real VM execution;
prior KVM/TCG lifecycle acceptance remains relevant because virtualization and
cleanup inputs are unchanged, while diagnostic behavior needs fresh evidence.

## Acceptance evidence

On 2026-09-26 the updated real-VM failure case and public smoke scenario passed
with `bazel_agent bazel test //tools/molecule/test:acceptance_test
//tools/molecule/test:smoke_test --test_filter=TestFailureCleansResources`.
The failure case retained an `ansible.builtin.fail` event with a task ID and
`failure_kind: assertion`; smoke completed all six lifecycle phases. Both
cleanups succeeded. Neither artifact tree contained the templated-name sentinel.
Smoke retained 27 task events and six recap lines; failure retained 23 events
and five recap lines. Source and executable hashes plus result artifacts are
recorded under ignored `out/molecule-runner-split/diagnostic-green-evidence` and
`diagnostic-test-inputs.json`, based on published head `5a19c8a607f` with this
review correction applied. Final candidate checks compare those exact runtime
inputs after formatting and archival.
