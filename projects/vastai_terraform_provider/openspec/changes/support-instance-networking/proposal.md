## Why

Terraform consumers need sensitive creation options and observed TCP port mappings.

## What Changes

Add replacement-only sensitive docker_options with an empty compatible default and computed mapped_ports. Preserve imported and legacy instances. Exercise actual Terraform lifecycle fixtures.

## Capabilities

### New Capabilities

- `instance-networking`: Terraform consumers need sensitive creation options and observed TCP port mappings.

## Impact

Owner-local source, packaging, documentation, and offline validation. No new external software or live deployment.
