## Why

The user requested a third_party download namespace equivalent to projects, with a versioned X OpenAPI snapshot published to both hosts.

## What Changes

- Support projects and third_party versioned publication roots.

## Capabilities

### New Capabilities

- `publication-namespaces`: Support projects and third_party versioned publication roots.

### Modified Capabilities

- `ssh-release-deployment`: namespace-aware public paths and project-only website steps.

## Impact

tools/rules_release; related changes are owned by infra/download, tools/rules_release, projects/ansible_collection and internal/x_sdk.
