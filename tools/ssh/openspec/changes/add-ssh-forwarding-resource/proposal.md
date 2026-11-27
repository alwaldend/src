## Why

AL commands that access a private service need a verified SSH forward for
the duration of each invocation.

## What Changes

Add a reusable AL resource provider that opens a loopback listener on an
available port after authenticating the SSH host and reaching the destination.

## Capabilities

### New Capabilities

- `ssh-forwarding-resource` manages invocation-scoped local SSH forwards.
