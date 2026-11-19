# project-dns Specification

## Purpose

Record the retirement of `rules_binary_toolchain` landing infrastructure while preserving
its reusable Bazel module.

## Requirements

### Requirement: Keep landing infrastructure retired

The module SHALL follow the [tools boundary](../../../../README.md) and SHALL
have no dedicated landing DNS declarations, Terraform root, or operational
source exports.

#### Scenario: Inspect the module after landing retirement

- **WHEN** the module is consumed from `tools/rules_binary_toolchain`
- **THEN** its reusable Bazel rules remain available without landing infrastructure.
