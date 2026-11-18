## Why

Provide certificate-authenticated access to the OpenCanary web decoy.

## What Changes

Route the public OpenCanary HTTPS hostname to the web decoy over the existing
WireGuard network. Verify the backend certificate and inherit the ingress
client-certificate requirement.

## Capabilities

### New Capabilities

- `opencanary-ingress`: Mutual TLS access to the OpenCanary web decoy.

### Modified Capabilities

None.

## Impact

Only the web route is published through ingress.
