## Why

The Grafana service deployment needs a reusable Ansible role.

## What Changes

Package the Fedora Grafana installation, configuration, and service checks
as `alwaldend.main.grafana`.

## Capabilities

### New Capabilities

- `grafana-role`: Install and configure Grafana with OIDC and protected credentials.

### Modified Capabilities

None.

## Impact

The Grafana infrastructure playbook consumes the role and supplies deployment inputs.
