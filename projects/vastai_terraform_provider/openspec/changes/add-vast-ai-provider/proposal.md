## Why

Vast.ai GPU rentals need a Terraform lifecycle so operators can plan, import,
update, and destroy instances alongside their other infrastructure.

## What Changes

- Add a Go Terraform Plugin Framework provider named `vastai`.
- Support offer search and on-demand instance lifecycle through Vast.ai REST.
- Add offline CLI acceptance, build targets, and a project landing.

## Capabilities

### New Capabilities

- `provider-lifecycle`: authentication, offer discovery, and instance management.

### Modified Capabilities

None.

## Impact

Adds `projects/vastai_terraform_provider`, root Go dependency declarations,
project registration, and generated inventories. The user approved
`github.com/hashicorp/terraform-plugin-framework v1.18.0` and its required
transitives. No live infrastructure changes or registry publication occur.
