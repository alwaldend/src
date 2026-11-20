## Why

Traefik deployments need a consistent way to publish metrics to one collector.

## What Changes

Define a shared metrics port default. Each deployment enables Prometheus
in its static Traefik configuration and scopes access in its host firewall.

## Capabilities

### New Capabilities

- `traefik-metrics`: Sets the default port used by deployment-owned metrics listeners.

### Modified Capabilities

None.

## Impact

Defines the shared port in the Traefik role and enables metrics on hosts with checked-in private routing to VictoriaMetrics.
