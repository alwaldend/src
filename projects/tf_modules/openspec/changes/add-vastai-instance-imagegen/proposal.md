## Why

Image-generation rentals need a reusable Terraform module that uses Vast.ai
templates while exposing only SSH. The stock web templates can publish several
services, so SSH mode alone is insufficient.

## What Changes

- Add `vastai_instance_imagegen` under `projects/tf_modules`.
- Accept a template hash and select the cheapest qualifying datacenter offer
  within configurable GPU, storage, and price limits; reject templates with public port
  options and require external web tunnels to be disabled.
- Return connection details and a loopback-bound SSH forwarding command.
- Add offline Terraform validation; do not create a deployment root or rent GPU
  capacity. Provider support for template lookup and selection is supporting work.

## Capabilities

### New Capabilities

- `vastai-instance-imagegen`: Template-based image generation with SSH-only access.

### Modified Capabilities

None.

## Impact

Owns module source, packaging, documentation, and offline checks in
`projects/tf_modules`. Supporting provider and SDK changes use their own change
records. No pinned application image or live infrastructure change is included.
