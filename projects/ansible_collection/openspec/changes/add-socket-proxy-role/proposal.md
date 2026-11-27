## Why

TCP socket activation is useful across hosts and can be configured with the
same systemd units.

## What Changes

- Package a role accepting the unit name, listener, and backend endpoint.

## Capabilities

### New Capabilities

- `socket-proxy` manages a socket-activated TCP proxy.

## Impact

Flux consumes the reusable role for HTTPS.
