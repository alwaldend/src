# Printer1 state migration

Completed on 2026-11-07. The existing DNS records and firewall rule are now
owned by the printers project, with unchanged live identities and matches.
The [deployment receipt](../openspec/changes/move-printer1-infrastructure/deployment.md)
records the operations and verification evidence. No printer record or
firewall rule was destroyed or recreated.

## Current ownership

- `//infra/printers/dns` owns both printer DNS views through the printers
  AppRole's `tf_backend/dns` backend.
- `//infra/printers/tf` owns the Wi-Fi exception through its `tf_backend/tf`
  backend and dedicated RouterOS credential.
- `//infra/mikrotik/accounts` owns the printers account and generated Vault
  credential, using the operator-supplied MikroTik administrator secret.
- `//infra/mikrotik/tf` retains router, switch and bare-metal DNS records.
  It no longer owns printer DNS. The former firewall backend has no bindings.

The MikroTik AppRole no longer has the ingress-secret read policy. The
printers identity cannot read the administrator credential. Do not restore
old deployment jobs or run the former source revision against its old state;
that source would attempt to recreate objects whose bindings have moved.

## Verification

With explicit authority for live review, use the normal text plans:

```sh
bazel run //infra/printers/dns:tf.plan
bazel run //infra/printers/tf:tf.plan
bazel run //infra/mikrotik/accounts:tf.plan
bazel run //infra/mikrotik/tf:tf.plan
```

All should report no resource changes. Check the existing firewall rule still
precedes the enabled terminal forward drop. Its place_before attribute is
creation-only and cannot be reconstructed by import; source ignores updates
to that attribute while retaining guarded insertion for newly created rules.

## Rollback

A rollback requires separate authority and a reviewed plan. Freeze both
writers, restore declarations in the intended owner, remove only the exact
current bindings, and import the verified existing objects into that owner's
backend. Never delete live records or the rule merely to move ownership, and
never keep two active state bindings for one object. If an import fails, keep
writers stopped and restore the binding before allowing an apply.

Use an owning state-only wrapper authenticated by the appropriate owner
AppRole; the old firewall backend is still available for controlled recovery.
Do not revive an ingress-credential dependency to access it. Its temporary
migration target and selector were removed after the empty-state check.
Keep credentials in Vault and retain only sanitized summaries of any recovery.
