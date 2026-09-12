## Context

See [proposal](proposal.md) for motivation and scope. The current source
baseline is commit `225d508ce3472498b960cdf5b95594750402e2ed`, inspected on
2026-09-12. This is source evidence, not a live DNS inventory.

The owning [BUILD file](../../../BUILD.bazel) aggregates project JSON inputs,
including external nested modules. The [configuration](../../../dnsconfig.js)
adds common apex/mail records and sends global and dc1 views to Cloudflare,
MikroTik, and BIND. Some JSON entries contain multiple record types. Existing
VM setup consumers also read names and addresses directly from these files.

## Goals / Non-Goals

Enable project-local reconciliation within the existing provider zones while
retaining a single declaration for each record. The contract is in the
[delta specification](specs/infra-dns/spec.md).

This does not introduce public child-zone delegation, change DNS providers,
parse arbitrary DNSControl JavaScript in Terraform, or authorize live changes.
The planning delivery leaves all implementation tasks pending.

## Decisions

### Project state is the ownership boundary

Proceed with a reusable module in the existing
[Terraform module collection](../../../../../projects/tf_modules/README.md).
Each owner calls it from `tf_setup` when that root exists, otherwise `tf`.
Do not instantiate it in both roots or collect all calls into one central
state. If an owner has neither root, add its minimal `tf` root using the
existing backend, Bazel packaging, and AL injection conventions. For nested
modules, establish a supported module dependency at implementation time.

This change owns DNS migration coordination. Before implementing the shared
module or a project batch, create or update linked changes in the affected
owners' OpenSpec workspaces. Those changes own module API and project-specific
bootstrap, state, and deployment requirements; link this coordination contract
instead of copying it. Record the links in the migration ownership map.

Provider instances and credentials come from the calling root. Preserve the
existing Vault/AppRole injection path and select minimum required privileges.
Module placement must not create a dependency cycle: reading declared JSON
for VM setup does not depend on DNS resources being created. Where MikroTik
connectivity depends on bootstrap infrastructure, plan that ordering before
selecting the pilot; do not move DNS into both stages to bypass the issue.

Recovery must also survive loss of a managed DNS name. In particular,
`infra/mikrotik` owns `router1.dc1`, while the current DNSControl and ingress
RouterOS connections use `router1.dc1.alwaldend.com`. Before transferring
that scope, define an independently reachable management endpoint through
the owning IaC, such as an IP endpoint with valid TLS identity or independently
resolved management name. Verify Terraform and rollback access while the
managed DNS record is unavailable in an authorized test. Do not disable TLS
verification to achieve independence. The same check applies to Vault/state
dependencies whose names fall within a migrating scope.

Alternatives: child zones require provider support and delegation changes;
multiple DNSControl writers require complementary ignore rules. Individual
Terraform resources provide ordinary deletion semantics without either.
State separation does not provide subdomain-scoped provider permissions.

### Preserve and normalize the existing declarations

The requested `dnscontrol.json` means the existing `dnsconfig.json` format.
Accept its decoded document as module input so callers retain existing file
paths and VM consumers. Flatten every type member of every logical key,
expand `all` to both views, and reject unknown types and destinations rather
than dropping them silently. Deduplicate repeated destinations after expansion.

Use stable keys including logical record key, type, and destination; mutable
values such as addresses must not form resource identity. Explicitly model
multiple values at the same name/type and provider record versus RRset
granularity. Validate ownership boundaries with exact names and subtrees,
including aliases such as Forgejo's `git` that do not match its directory name.

Compare normalized effective output before migration: names, targets,
relative names versus FQDNs, apex notation, values, type, view, TTL, and
Cloudflare proxy settings. Do not infer uniform TTL from the source: the
current A-record TTL modifier placement differs from AAAA/CNAME placement.
Extend the canonical format only as needed for currently embedded MX/TXT and
other shared records, preserving explicit TTLs and mail priorities.

### Transfer ownership gradually

Keep one small migration ownership map with references to source files,
record scopes, destination views, target Terraform roots, and cutover status.
Do not copy record values into it. Derive the legacy source selection and
ignore rules from this map where practical, so they cannot disagree.

For each transferred subtree, the main operational configuration must stop
declaring its records and include both the exact name and all descendants
(for example `cloud` and `**.cloud`). Generate ignores per provider view;
additional aliases need their own scopes. Apex/shared record transfers need
precise record/type boundaries where a subtree would be too broad. Never
disable the DNSControl ignore safety check or use zone-wide `NO_PURGE` as
the migration mechanism.

An ignore configuration protects only runs using that revision. Coordinate
cutover and prevent old deployment candidates from running; merely merging
an ignore rule does not protect against an older checkout's central push.

### Keep BIND only as a derived documentation view

First investigate the pinned DNSControl's BIND-only generation path. The
full documentation input must include migrated and remaining declarations,
while operational DNSControl input excludes migrated ownership. Reuse one
normalization implementation or its exported normalized data; do not maintain
equivalent JSON-to-record logic independently in Terraform and JavaScript.

One candidate is a provider-free Terraform normalization submodule whose
output can feed both provider resources and the existing BIND renderer.
Select the smallest offline route after checking the pinned tools; no live
Terraform state or provider reads may become a documentation prerequisite.
The renderer must have no live provider configuration or credentials, write
only designated build outputs, and preserve both views with deterministic
serial handling. Label output as declared desired state, not observed DNS.

Retaining DNSControl as a BIND renderer is acceptable after its live writers
are retired. If no practical route meets these constraints, record evidence
and the limitation, remove stale documentation claims, and complete the
core migration without preserving duplicate definitions or translation logic.

## Risks / Trade-offs

- Overlapping states or stale central pushes can revert records -> enforce
  one owner per provider record/RRset, validate scopes, and coordinate cutovers.
- Import identifiers, TTL normalization, or multi-value behavior can differ
  by provider -> confirm pinned Cloudflare and MikroTik resource capabilities
  and compare representative effective output before migrating any project.
- Independent states do not share locks -> disjoint ownership is required;
  coordinate changes to a shared RRset, alias boundary, or zone-level setting.
- Some landing-only projects lack Terraform roots -> include their minimal
  `tf` integration in the rollout, not an exception that leaves a live
  monolithic writer indefinitely.
- The optional export may require substantial tooling -> time-box feasibility
  during implementation and record a reasoned retain/omit decision.

## Migration Plan

1. Inventory declared inputs and roots, including nested modules and shared
   JavaScript records. Establish canonical scope ownership and offline parity
   fixtures without contacting live infrastructure. Link owner-specific
   OpenSpec changes before implementing their respective batches.
2. Implement and validate the reusable module, provider mappings, packaging,
   and credential references. Investigate optional offline BIND generation.
3. Select a pilot covering a representative view and existing Terraform root.
   Prepare its module invocation, imports, legacy exclusions/ignores, and
   reviewed expected plan together; validate offline before operational work.
   Include independent provider, Vault, and backend recovery connectivity
   before any batch that owns those endpoints' DNS records is transferred.
4. Under separately granted exact operational authority, stop competing
   writes, activate the central ignore/exclusion revision, import existing
   records into the designated project state, and verify a no-change adoption
   plan. Only then permit project Terraform writes and resume current central
   runs for the remaining scopes. Preserve existing records during the gap.
5. Repeat per project with independent receipts and ownership-map updates.
   Verify update and deletion behavior, plus preservation of unrelated
   records, using fixtures and an explicitly authorized test environment.
6. Migrate common apex/mail records to the DNS infrastructure's own root,
   audit that every declaration has one owner, and retire all live central
   DNSControl writers and their obsolete credentials/references. Keep only
   the accepted offline documentation path, if feasible.

Rollback for a transferred scope: stop both writers; reconcile the canonical
declarations with any intentional changes already applied; relinquish the
Terraform binding without destroying the live records through an authorized
state operation; restore the central declarations and remove only that
scope's ignores; preview the restoration before allowing central writes.
Never restore central ownership while Terraform can still reconcile it.

## Open Questions

- Which existing project provides the smallest representative pilot without
  bootstrap connectivity dependencies?
- Which pinned MikroTik Terraform resource/import shape preserves all current
  static-record behavior, and which additional provider pin is needed?
- Can the pinned BIND renderer consume a shared provider-free normalized
  export with reproducible serials at reasonable implementation cost?

Next action: perform task 1.1 in [tasks](tasks.md) when implementation is
requested. No provider compatibility tests or live migrations have run.
