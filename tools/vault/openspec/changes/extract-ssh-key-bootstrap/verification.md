# Verification

Observed 2026-11-17. The command is owned by tools/vault/cmd/ssh_key; its Go and
BUILD source contain no Forgejo identifiers or infrastructure dependencies.
The dependency query `filter("^//infra/", deps(//tools/vault/cmd/ssh_key:ssh_key))`
returns no targets. The stored schema and cache-validation behavior are unchanged.

The relocated infra/forgejo/test/signing_key acceptance fixture passes using the
new neutral environment inputs and actual Terraform provisioner. Setup wrapper
runfiles were materialized using --script_path without executing Terraform;
its declared relative executable path resolves to the new packaged tool.

All 33 selected acceptance, specification and repository-quality checks pass;
unchanged Forgejo role acceptance is reused from the valid Bazel cache. Semantic
lint and deployment packaging pass. The final aggregate candidate is reviewed
and checked again before updating the existing PR. No live generation, Vault
write, infrastructure plan, or deployment was performed.
