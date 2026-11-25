## Why

The Mullvad load-balancing inbound is enabled, but it has no active balancing
route and the observatory selector matches none of its WireGuard outbounds.

## What Changes

- Route the Mullvad load-balancing inbound through the existing Xray balancer.
- Probe Mullvad WireGuard outbounds using their tag prefix.

## Capabilities

### Modified Capabilities

- `infra-threexui` requires an active Mullvad balancing route and health probes.

## Impact

The shared Terraform node module changes. This rollout targets `njalla1`.
