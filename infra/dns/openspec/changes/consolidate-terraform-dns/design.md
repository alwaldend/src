## Context

`dns_records` owns the schema, normalization, and provider resources. Standalone roots also need provider configuration, an HTTP backend, and Bazel operation targets.

## Decisions

- Package the shared Terraform files at the owner's existing working directory. The module name, source, provider alias, and relative declaration path retain their existing meanings.
- Select either the Cloudflare-only or Cloudflare-and-RouterOS entry point when declaring the Bazel targets. Local Inference requires only Cloudflare and keeps DNS independent of the private router.
- Keep the owner responsible for its declaration and AL configuration, including its AppRole and Vault state path.
- Give extracted roots a separate AL configuration so service wrappers using broad labels cannot select their backend or credentials.
- Store declarations alongside each DNS root. Complete the live transfer using the previous revision for non-destructive removal and the new revision for imports before ordinary plans or applies.

## Verification

Build every owner root and its affected service root. Compare provider
requirements and resource addresses with the previous definitions. Run
offline Terraform and declaration checks against the relocated inputs.
