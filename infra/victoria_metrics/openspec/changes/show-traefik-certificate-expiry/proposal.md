## Why

Operators need one place to see which certificates reported by Traefik are
approaching expiration across the monitored instances.

## What Changes

- Add a certificate expiration panel to the Monitoring overview dashboard.

## Capabilities

### New Capabilities

- `traefik-certificate-expiry` shows days remaining for certificates reported by scraped Traefik instances.

### Modified Capabilities

None.

## Impact

The VictoriaMetrics VMUI dashboard uses the existing Traefik metrics.
