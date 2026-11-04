## Verification

Observed on 2026-11-04 in branch `t3/require-local-pr-review`.

- `bazel run //:format`: passed; no unrelated formatter changes.
- `bazel run //.agents:write_skills`: passed; repeated generation was stable.
  Canonical and generated delivery instructions and references match byte for byte.
- `bazel test //:repo_quality_test //.agents:write_skills_test //:buildifier_test
@rules_al//skills/al-delivery:eval_config_test`: 30 tests passed.
- `bazel build --config=lint @rules_al//skills/al-delivery:skill
@rules_al//openspec:source`: passed, including skill validation. These Markdown
  targets have no language-specific semantic lint actions.
- Owner-scoped `openspec validate require-local-pr-review --strict` through
  `//tools/openspec`: passed.
- `git diff --check`: passed.

Offline evaluation validation checks the configuration and staged skill, not
agent behavior. No live model evaluation was run; the suite's README describes
the external Git and forge operations that it cannot reproduce.

## Local review

`review-agent` review of the complete working-tree diff: **No findings.**
At that review, the comparison base and HEAD were both
`b687079120b12585e872a4d607b6be2f78085c25`; this is an uncommitted candidate.
The review covered the canonical procedure, generated projections, evaluation
scenarios, and change artifacts. It checked initial publication, changed
rebases, review fixes, unchanged recovery, and Forgejo capability limits.
Record final prepared-candidate review under task scratch during delivery.

## Delivery environment

`al delivery provider` identified the configured `github.com/alwaldend/src`
destination and available adapter. `al delivery inspect --base master` then
failed GitHub repository preflight with HTTP 401, `Requires authentication`.
Defect identifier: `delivery-github-authentication-401`. A read-only authenticated
API probe succeeded outside the sandbox; the sandboxed probe failed. This
isolates the observed failure to the execution context, without proving token
expiration or revocation. No credentials or host configuration were changed.

The user selected `--git-transport gh`. The installed AL rejected that flag;
the current source-built task launcher supports it, and its provider and
repository inspection succeeded outside the sandbox. At the user's request,
`bazel run //tools/al:install` updated the installed executable; its provider
check confirmed support for gh transport. Retain this transport selection
throughout delivery. Publication and remote review state belong to delivery receipts and
the forge, rather than a second checked-in status record.
