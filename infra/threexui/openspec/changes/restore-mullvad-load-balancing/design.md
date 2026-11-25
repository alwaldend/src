## Context

Both nodes use the same Terraform module. Xray uses prefix matching for
observatory selectors, and the current selector includes a literal asterisk.

## Goals

- Select a healthy Mullvad WireGuard outbound for the dedicated inbound.
- Apply only the required changes to `njalla1`.

## Decisions

Restore the dedicated routing rule after the existing advertised-domain and
private-IP blocking rules. Select the `out-mullvad-min-` tag prefix for probes
and make routing depend on observatory configuration. Keep the existing
`leastLoad` strategy and its `blocked` fallback.

Review a targeted interactive Terraform apply for the Njalla routing and
observatory resources. Verify both settings through a subsequent targeted plan.

The first rollout saved the correct selector under a named object that Xray
ignored. Build the patched 3x-ui provider with isolated Go dependencies. Use
its `default` observatory tag to write the selector directly under
`observatory` and read back the same Terraform tag after migration.
