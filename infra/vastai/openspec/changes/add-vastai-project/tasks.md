## 1. Project foundation

- [x] 1.1 Add the project README, AL authentication/backend configuration, and
      internal Bazel targets; build the AL configuration and Vault launcher and
      inspect their packaged inputs without authenticating.
- [x] 1.2 Export the owner OpenSpec workspace and regenerate the validation
      inventory; run the owner's OpenSpec validation test and freshness check.

## 2. Integration validation

- [x] 2.1 Run repository quality checks and affected semantic lint; inspect
      the final diff for secret values and namespace consistency.

## Validation evidence

On 2026-11-04, the AL configuration, Vault launcher, documentation, and Vault
Terraform plan launcher built successfully without execution. Generated AL
JSON selected the expected identity, backend path, mount, and label. The
Terraform runfiles manifest included the new AppRole source. Inventory
regeneration was stable. All 30 selected tests passed, including repository
quality, Terraform formatting, inventory freshness, and both owner OpenSpec
validators. Affected-target semantic lint passed. Source inspection found no
literal credentials. No live Vault or Vast.ai operation was performed.

## 3. API-key injection

- [x] 3.1 Add the labelled KV-to-environment injector and component CLI
      wrapper; build and inspect generated AL configuration and runfiles for
      the exact secret path, field, environment name, label, and plugin binary.
- [x] 3.2 Document the credential contract and CLI entry point; verify
      repository quality, owner OpenSpec validation, existing injector
      behavior checks, and affected semantic lint without reading a live key.

On 2026-11-04, the injector and CLI launcher built successfully. Generated AL
JSON and runfiles verified the secret path, field, environment name, label,
injector binary, and pinned CLI. All 32 selected tests and affected semantic
lint passed. The pinned CLI reads `VAST_API_KEY` directly.

The user then explicitly requested a live read-only check. Running
`bazel run //infra/vastai -- --raw show user` succeeded through the injector
and returned an authenticated account response from `GET /users/current`.
Only a sanitized success receipt was retained; the raw response and logs were
removed without publishing account details. No Vast.ai resources were changed.
