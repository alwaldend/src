## Decision

Extend the existing exact-path deployment read policy with the existing
`src_infra_threexui/http_proxies` KV data path. This secret stores a map of
proxy accounts; Vault cannot grant field-level reads. The controller selects
only `t3code` for injection and the runner receives no Vault identity. Keep
the grant read-only, without metadata or sibling access. A dedicated new
secret/account would add credential provisioning beyond reuse requested here.

## Validation

Run Terraform formatting and policy source review. No policy application,
credential rotation or secret writes are authorized by validation.
