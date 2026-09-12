## 1. Inventory and compatibility

- [ ] 1.1 Inventory declared DNS sources, aliases, shared JavaScript records,
      and owning Terraform roots, including nested modules; deliver an ownership
      map referencing sources and covering every record without copying values.
- [ ] 1.2 Confirm pinned Cloudflare and MikroTik Terraform resource/import
      support; record provider mappings and offline fixtures covering current
      types, multi-value records, TTLs, proxy settings, and normalization.
- [ ] 1.3 Select the pilot and migration batches; verify each owner uses
      `tf_setup` when present, otherwise `tf`, and document bootstrap dependency
      ordering and any new root/backend packaging required.

## 2. Shared Terraform module

- [ ] 2.0 Create or update the module owner's linked OpenSpec change; verify
      its API requirements and dependency contract are owned there and referenced
      by the migration map before implementing the module.
- [ ] 2.1 Add the reusable module and canonical JSON normalization; verify
      multi-type entries, all/global/dc1 routing, destination deduplication,
      invalid-input rejection, and stable identities with offline tests.
- [ ] 2.2 Add provider-specific record resources and caller-supplied provider
      wiring; verify effective output parity and that changes/deletions target
      only owned resources using provider mocks or equivalent offline fixtures.
- [ ] 2.3 Package module inputs, provider dependencies, backend references,
      and AL/Vault injection through existing repository conventions; pass
      package validation and confirm JSON-based VM consumers still work.

## 3. Coexistence and pilot

- [ ] 3.0 Prepare independent management/recovery access for scopes owning
      provider, Vault, or state endpoint names; deliver owning IaC and verify
      authenticated access without the managed DNS record in an authorized
      test before transferring those scopes, especially `infra/mikrotik`.
- [ ] 3.1 Add central source exclusions and exact-name/descendant ignore
      rules from ownership metadata; verify migrated scopes and aliases are
      protected in their views and remaining scopes still reconcile normally.
- [ ] 3.2 Add the pilot invocation to its selected root and prepare imports,
      cutover instructions, old-writer controls, and rollback; deliver reviewed
      source plus offline checks before requesting any live operation.
- [ ] 3.3 Under explicit operational authority, perform pilot adoption;
      retain sanitized evidence of preserved record identities, no-change
      adoption plan, and a single active owner for the migrated scope.

## 4. Read-only BIND feasibility

- [ ] 4.1 Evaluate a BIND-only renderer consuming canonical normalized data;
      record a retain/omit verdict with evidence for offline execution and shared
      translation, without adding duplicate record definitions.
- [ ] 4.2 If feasible, implement the documentation target; verify complete
      global/dc1 output before and after migration, deterministic generation,
      source provenance, and absence of provider access or credentials. Otherwise
      document the limitation and identify stale output/claims for removal.

## 5. Gradual rollout and retirement

- [ ] 5.1 Expand the inventory into one reviewable checklist item per project
      batch before migrating it; prepare each selected root, imports, ignores,
      and rollback, link its owning OpenSpec change with local requirements,
      and pass its offline checks independently.
- [ ] 5.2 Under explicit authority per operational scope, adopt each prepared
      batch; verify unchanged effective records and unrelated-scope preservation,
      and retain sanitized ownership and validation receipts.
- [ ] 5.3 Move common apex/mail declarations to their canonical owner and
      migrate them through the module; verify MX priorities, TXT values, TTLs,
      multi-value apex records, and both views against the baseline.
- [ ] 5.4 Audit complete ownership and retire central live DNSControl
      entry points and obsolete credential references; verify every declared
      record has one Terraform owner and retained BIND output is documentation only.
- [ ] 5.5 Update operating documentation and affected owner specifications;
      pass package checks, strict OpenSpec validation, and delivery gates, then
      archive this change only when the migration and its evidence are complete.
