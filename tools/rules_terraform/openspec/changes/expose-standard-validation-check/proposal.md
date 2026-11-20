## Why

Offline Terraform roots need the same validation command and test naming.

## What Changes

Allow test maps to opt into a standard Terraform validation target alongside
the default format check.

## Capabilities

### New Capabilities

- `terraform-validation-check`: Adds the standard validation target for backend-free inputs.

## Impact

Terraform package owners can reuse the validation command through the test map.
