## Why

The Terraform command maps and repository provider pins were placed under AL.
Their owners should be the reusable Terraform module and the shared external
dependency package, while AL remains a generic command wrapper.

## What Changes

- Move Terraform command maps into rules_terraform with optional caller-supplied wrappers.
- Move shared provider declarations and their offline integration fixture into third_party/terraform.
- Migrate consumers to direct rule-module loads and explicit generic AL wrappers.
- Remove projects/al/rules/terraform; retain root command bindings beside tools/al configuration.

## Capabilities

### Modified Capabilities

- `terraform-execution`: command-map ownership and explicit optional wrapper composition.

## Impact

Terraform BUILD consumers, dependency include discovery, reusable rule tests,
repository documentation, and infrastructure procedures. Provider versions,
operation names, credentials, and live infrastructure remain unchanged.
