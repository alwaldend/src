## Why

The Flux AppRole needs the standard DC1 guest resources to provision its K3s VM
through Xen Orchestra.

## What Changes

Register the existing Flux AppRole in the resource-set inventory.

## Impact

The existing per-AppRole authorization and resource-set contract applies to the
Flux VM.
