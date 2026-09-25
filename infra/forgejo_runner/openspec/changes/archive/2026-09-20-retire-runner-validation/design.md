## Context

Runner acceptance succeeded in live run 8 on validation commit
`896d36f944104b0586fbe02d5f01a5e4ff382487`. Review then explicitly requested
removal of the harness, rather than retaining a branch-only test.

## Goals / Non-Goals

Remove runtime acceptance code, helper targets, configuration and current
references. Preserve runner provisioning, default-branch Vault identity and
historical evidence. Do not delete the historical remote branch or claim an
undeployed policy change has taken effect.

## Decisions

Delete `ci.json` and its Bazel filegroup because no validation consumer remains.
The Vault role derives its sole allowed ref from the existing repository catalog.
Keep Android prerequisites: Bazel resolves NDK toolchains before the installer
can be run. Remove Node's extra CA selection after the deployed Node succeeded
against both internal services using system trust alone.

## Risks / Trade-offs

The deployed Vault role can retain its previous validation ref until a scoped
owning Terraform deployment. Source validation cannot establish live policy
convergence. Historical archived records remain unchanged.

## Verification

Runner Ansible packaging and affected lint pass; the Vault Terraform format
test passes. Strict OpenSpec validation passes after explicitly replacing the
retired smoke requirement. Current code and documentation no longer reference
the removed harness; archived records retain historical provenance. Repository
command behavior is verified by the related repository entry-point change.
