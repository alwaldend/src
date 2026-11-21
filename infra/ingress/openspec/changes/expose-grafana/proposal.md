## Why

Remote operators need to reach Grafana through the shared ingress.

## What Changes

Route the Grafana service hostname through shared ingress with its standard
client certificate policy and verified HTTPS upstream.

## Capabilities

### New Capabilities

- `grafana-ingress`: Client-authenticated ingress for Grafana.

### Modified Capabilities

None.

## Impact

Adds one shared ingress route to the Grafana-owned backend.
