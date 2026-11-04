## Authentication observations

On 2026-11-04, the same `gh api --hostname github.com user` query returned HTTP
401 inside the task sandbox and succeeded outside it through execution approval.
Only an authenticated boolean was projected; no account details or tokens were
displayed. Neither GH_TOKEN nor GITHUB_TOKEN was set. A credential-availability
probe captured CLI token output privately in process memory and reported only
an availability boolean: false inside the sandbox.

The documented sanitized status command returned:

- Sandbox: active true, state error, tokenSource default, exit zero.
- Approved context: active true, state success, tokenSource keyring, exit zero.

This establishes an execution-context credential-access problem. It does not
establish an expired token or justify changing global configuration. The changed
AL launcher reproduced the new diagnostic inside the sandbox; the identical
`delivery inspect --git-transport gh --base master` invocation succeeded through
execution approval outside it. The sanitized provider report identified
github.com/alwaldend/src, the authorized destination. No credentials, remotes,
global Git settings, or host permissions were changed.

## Validation and review

- The CLI fixture was written before the implementation. Both HTTP 401 cases
  failed on the missing advice; HTTP 403 cases passed.
- The full delivery suite and repository quality, buildifier, and skill
  discovery checks passed: 30 tests. Initial invocations used an unavailable
  root-package label for AL's skill; the corrected
  `@rules_al//skills/al-delivery:eval_config_test` passed separately.
- Semantic lint passed for `@rules_al//internal/delivery:go_test` and
  `@rules_al//skills/al-delivery:{skill,eval_config_test}`.
- OpenSpec strict validation passed. Configured formatters ran. Skill discovery
  regeneration ran twice and the authentication reference matches its projection.
- Local review using review-agent found no actionable defects. Inspection covered
  the complete tracked diff and new artifacts, API-error cause preservation,
  credential redaction, both Git modes, JSON status semantics, approval refusal,
  and the generated skill projection. Offline evaluation validation establishes
  configuration validity, not model behavior or live credential-store access.

After preparation, rerun ordinary required tests and lint against the returned
candidate, inspect the base-to-head diff and representative output, then publish
and verify through the same local AL launcher and gh transport. Logs, launcher,
commit message, inspection, and receipt belong under
`out/delivery-authentication/`; Git and AL receipts own final publication state.

## Agent ergonomics

Stable defect: `DELIVERY-SANDBOX-CREDENTIAL-ACCESS` (live evidence). A sandbox
401 previously led to repeated transport or login suggestions despite a valid
keyring login. The narrow remedy is conditional diagnostic advice and a linked
recovery reference with three evaluation cases. Measure future recovery by one
sanitized context comparison followed by the same approved AL invocation,
without login changes or token copying.

Discovery overhead included oversized initial source reads and one incorrect
nested-module label selection. Subsequent inspection used bounded source regions
and the owning `@rules_al` labels. No host migration or additional agent was
needed; the maintained change and ignored receipt directory provide resume state.
