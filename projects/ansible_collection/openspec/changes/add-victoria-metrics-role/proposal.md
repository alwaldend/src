## Why

Make the VictoriaMetrics service configuration reusable across Ansible inventories.

## What Changes

Add the `alwaldend.main.victoria_metrics` role. The caller supplies the
controller-side binary while the role manages the guest service.

## Capabilities

### New Capabilities

- `victoria-metrics-role`: Installs and configures the VictoriaMetrics single-node service.

### Modified Capabilities

None.

## Impact

The VictoriaMetrics infrastructure playbook imports the role from the shared
collection.
