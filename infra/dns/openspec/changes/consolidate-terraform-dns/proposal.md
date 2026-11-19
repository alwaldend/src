## Why

Standalone DNS owners repeat Terraform root configuration and its Bazel wiring. They already share the record implementation.

## What Changes

Use one shared Bazel entry point for every managed DNS owner. Move DNS
out of mixed service and provisioning roots into dedicated owner roots. Move
canonical declarations into each owner's `dns/` package and remove temporary
transition code.

## Capabilities

- Modified `infra-dns`: DNS commands and Terraform state are independent of service provisioning.

## Impact

Existing DNS-only roots keep their backends and operation names. Extracted
roots use new owner-local backends and require an independently authorized
state transfer. Record declarations, provider IDs, and resource addresses
are preserved.
