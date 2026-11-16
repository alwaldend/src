## Context

See proposal.md for motivation. Forgejo generates a UUID and secret on creation
and returns the secret only once. The installed fj CLI has no runner registration
command. The service Terraform already receives Forgejo and Vault environment
credentials through AL.

## Goals / Non-Goals

Bootstrap credentials with a packaged local tool during an authorized apply.
Keep secret values in Vault. Runner host deployment, credential rotation,
remote adoption, deletion and drift reconciliation remain explicit operations.

## Decisions

Proceed with a component-owned Go tool using already pinned Cobra, Vault API
and the standard HTTP client. A stable terraform_data resource invokes its
ensure command through an explicit executable interpreter. There are no
triggers_replace, timestamps, file hashes, or secret-valued resource attributes.
An ordinary apply does not rerun the tool, even when tool source changes.

The tool reuses complete matching cached credentials from Vault. When no record
exists, it checks all runner-list pages for a same-name runner before one
creation POST and a create-only Vault write. A failed save or lost creation
response is reconciled explicitly: the same-name check prevents a blind duplicate
on retry. Permission and transport failures never mean the cache is absent.
The owning Terraform backend lock serializes execution. Independent concurrent
invocations are unsupported.

The source declaration selects secure1 for alwaldend/src, matching infra/dev's
documented runner host. The secret stays below the Forgejo AppRole's existing
KV namespace, requiring no new writer policy or provider.

The user explicitly selected creation-only execution and caching. A pending
reservation protocol and tool-content replacement triggers were rejected as
unnecessary machinery. A managed REST provider would provide a fuller lifecycle
but adds a dependency; this command is a bounded bootstrap instead.

## Risks / Trade-offs

- Lost creation response or failed Vault save leaves a remote registration
  without cached credentials → stop on a same-name collision and require
  separately authorized reconciliation.
- Credential drift is not polled → retain cached credentials until an operator
  explicitly reconciles them; source changes do not invalidate them.
- Removing the Terraform resource retains the runner and Vault record → cleanup
  is an explicit operation and no destructive provisioner is installed.
- Upstream errors can contain credentials → suppress their text while retaining
  error causes, bound API responses, and never print response bodies.

## Migration Plan

No state migration is required. Review an owning interactive apply before the
initial registration. Existing same-name runners require verified adoption or
cleanup outside this command. Replacing the bootstrap resource reuses cached
credentials. Changing an existing declaration's name or repository requires
explicit reconciliation; provisioner edits alone do not reconfigure the runner.
