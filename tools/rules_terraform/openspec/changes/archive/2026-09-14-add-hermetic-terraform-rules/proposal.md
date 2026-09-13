## Why

Terraform currently downloads providers during command execution, and 43 source
lockfiles duplicate release selections. Reusable Bazel rules should acquire
pinned providers and make the runtime dependency graph explicit.

## What Changes

- Add standalone `tools/rules_terraform` with a Bzlmod provider extension,
  provider metadata, executable rules, and offline checks.
- Package selected archives in Terraform's filesystem mirror layout and make
  initialization use only that mirror.
- Use one version per provider source, with a resolution boundary that can
  support explicit override policy later. Standardize Yandex on 0.203.0.
- Preserve AL/Vault integration in `projects/al/rules/terraform`.
- **BREAKING**: remove `tools/terraform`, migrate its consumers, and remove source
  `.terraform.lock.hcl` files.

## Capabilities

### New Capabilities

- `terraform-execution`: pinned acquisition, provider runfiles, and controlled
  Terraform execution for reusable Bazel consumers.

### Modified Capabilities

None.

## Impact

All repository Terraform command/test targets, DNS fixture tests, formatters,
root aliases, dependency declarations, documentation, and the standalone
workspace validation matrix. This migration authorizes source implementation
and offline checks; it does not apply infrastructure changes.
