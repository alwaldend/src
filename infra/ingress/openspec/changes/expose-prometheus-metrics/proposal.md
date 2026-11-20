## Why

Collect Traefik request, latency, and service metrics from shared ingress.

## What Changes

Enable Traefik's native Prometheus endpoint on a dedicated entry point.
Allow access to that port from the VictoriaMetrics guest address.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `infra-ingress`: Provides a private Prometheus endpoint for the metrics collector.

## Impact

Changes the ingress Ansible configuration and exposes the private target to
the VictoriaMetrics scrape configuration.
