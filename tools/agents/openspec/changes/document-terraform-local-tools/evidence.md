## Verification

The canonical `repo-infra` skill builds and validates. Its required offline
Promptfoo configuration check passes, and managed skill discovery was refreshed.
The two new cases cover stable creation-only execution and interrupted creation
without cached credentials. This checks harness loading, not model behavior.

The linked Forgejo implementation supplies isolated executable and production
Terraform provisioner evidence: ordinary applies do not rerun after tool edits,
and explicit replacement reuses Vault without creating another runner. No live
evaluation, plan, apply, registration or Vault write was performed.

The Terraform reference remains the sole procedure owner; SKILL.md routes to it
and the Forgejo owner retains command-specific facts and recovery details.
