## Current design

The component-owned `dns/` Terraform project follows `infra/download/dns`.
It selects only Cloudflare, DNS credentials and the component AppRole, with
its own `/tf_backend/dns` state. Canonical records remain in dnsconfig.json;
GPU provisioning uses a separate root and backend.

The delivered project contains only the active DNS configuration. It exposes
no legacy backend, migration command, import map or backward-compatibility path.
The user explicitly rejected retaining that scaffolding.

## Completed state transfer

The three existing record bindings were transferred without provider mutations
before this simplification. The execution evidence records the completed
operation; no migration configuration is retained in the project.
