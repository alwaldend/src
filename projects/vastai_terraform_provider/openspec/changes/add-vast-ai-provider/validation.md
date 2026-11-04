# Validation

Candidate scope: `projects/vastai_terraform_provider`, its approved Framework
dependency graph, and generated project/OpenSpec integration. Observed
2026-11-04; exact publication identity belongs to the AL delivery receipt.

## Acceptance evidence

- `bazel test //projects/vastai_terraform_provider/test/acceptance:acceptance_test
//infra/src/openspec/validation:projects_vastai_terraform_provider_validate_test
//:repo_quality_test`: all 29 tests passed.
- `bazel build --config=lint //projects/vastai_terraform_provider/...`: passed.
- `bazel test //infra/dns:config_test`: all four constituent tests passed.
- `bazel build //projects/alwaldend.com:site`: passed. Inspected the project
  landing, README documentation, project index, and Go taxonomy HTML for
  project identity and documentation/source links.
- `bazel mod deps --lockfile_mode=update`: passed; no lockfile change required.
- Inventory regeneration ran twice with stable projections. Gazelle generated
  the Go targets and was rerun after final imports.

The CLI acceptance exercises lifecycle, unchanged plans, import reconciliation,
disk replacement, missing remote instances, empty search, invalid configuration,
missing credentials, API failures, malformed replies, and redirect refusal.
The test's `test.outputs/lifecycle.json` records non-sensitive HTTP operations.
It also verifies returned instance credentials are absent from state and
diagnostics. Validation used literal loopback HTTP only; live Vast.ai account
compatibility and Terraform Registry publication remain unverified.

## Session ergonomics

- `root-tidy-generated-imports` (observed): one whole-root Go tidy failed on
  existing Bazel-generated Anki/protobuf imports. Scoped framework dependency
  resolution succeeded. Future provider updates should use scoped resolution
  and Bazel compilation; do not treat whole-root tidy as an acceptance gate.
- `fixture-runfiles-environment` (fixture-tested): Terraform's pinned launcher
  requires runfiles environment variables. The harness now propagates them.
- `fixture-unix-socket-length` (fixture-tested): absolute Bazel temporary paths
  exceeded Unix socket limits. Relative TMPDIR keeps sockets in the same
  private fixture directory with short addresses.
- `formatter-log-bounds` (observed): formatter diagnostics can contain very
  long file lists. Retain logs under ignored out and extract bounded bytes or
  selected diagnostic fields when resuming.

No shared host configuration changed. No additional dependencies were added
beyond the approved framework graph, and no live rental or deployment ran.

## Explicit-null deletion follow-up

The review identified a legitimate absence response that the initial fixture
classified as a failure. The
[official Vast.ai client](https://github.com/vast-ai/vast-python/blob/master/vast.py#L5303-L5313)
treats a successful `instances: null` response as an absent instance.
The added CLI scenario reproduced blocked refresh on the published candidate.
The corrected client preserves raw field presence, maps only explicit null to
the not-found sentinel, and retains errors for missing or mismatched instances.

The CLI suite verifies state removal, recreation planning, recreation, and
destruction for both HTTP 404 and explicit null. Separate read-failure cases
cover HTTP errors, omitted fields, mismatched IDs, and rejected responses with
null fields. Repeatable reports are emitted as
`test.outputs/remote-deletion-http-not-found.json` and
`test.outputs/remote-deletion-null-instance.json`.
