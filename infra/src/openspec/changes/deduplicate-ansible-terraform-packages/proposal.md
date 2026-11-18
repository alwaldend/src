## Why

Ansible packages repeat source mappings, packaged executables and AL command wiring.
Terraform packages repeat command/test inputs, providers and configuration.
Shared package macros will remove this duplication while preserving execution behavior.

## What Changes

- Add an Ansible package macro alongside the existing Ansible rules.
- Add a Terraform package macro alongside the repository AL Terraform integration.
- Migrate existing package graphs, retaining names, selectors, runtime inputs,
  custom operations, visibility and offline test configuration.

## Capabilities

No specification-level behavior changes; this is a refactor (`skip_specs: true`).

## Impact

Shared Bazel integration and its consumers across infrastructure, users, data,
website operations and third-party image configuration. No dependencies or live
infrastructure changes. Acceptance compares expanded targets before/after and
builds representative packages, runs offline tests, lint and repository quality.
