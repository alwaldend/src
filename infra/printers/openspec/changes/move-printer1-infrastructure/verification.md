# Source verification

Source validation was observed on 2026-11-07 before live operations. The
subsequent authorized migration is complete; see the
[deployment receipt](deployment.md) for provisioning, preserved identities,
no-change plans and acceptance limits.

The scoped Bazel invocation passed all 15 tests: formatting for the three
new roots and existing MikroTik and Vault roots, provider schema validation
for the account, printer DNS and printer firewall roots, DNS ownership checks,
and all three affected OpenSpec workspaces. Representative Terraform test
outputs report valid configuration using only pinned provider mirrors.

The repository quality suite passed all 27 tests. Affected semantic lint built
159 targets successfully. Inventory generation ran twice with stable output.
Bzlmod dependency reconciliation succeeded without lockfile changes. The
Random 3.9.1 Linux amd64 archive digest matched the publisher SHA256 list;
only its immutable URL and SRI integrity are checked in.

Inspection confirms the moved firewall keeps its resource address, match,
comment and exact-one-drop precondition, while canonical printer DNS retains
its key, name, address and views. AL stage labels, plugin data, owner backends
and Vault policy paths agree. New password consumers are sensitive and have
no credential outputs. The printers AppRole has no ingress-secret grant.

Offline checks do not exercise live drop lookup failures, bootstrap account
permissions, the new account's live authentication, imports, or Wi-Fi printing.
The deployment procedure and rollback are in
[the migration guide](../../../docs/state-migration.md). The existing live objects now belong to the printers state roots. Final candidate checks, local review and publication
are recorded by task-local delivery evidence under out/printer1-migration.

The provider-environment follow-up removes the redundant RouterOS variables
and reuses the shared injection helper, retaining credential paths and stage
labels. Follow-up results are recorded in the task-local delivery evidence.

The network diagram was rendered through the site Mermaid rule in both light
and dark variants and inspected for readable labels and unobstructed edges.
The printer/router connection is bidirectional; both device DNS labels and
the existing printer resolution are preserved. Wi-Fi forwarding restrictions
are explicitly limited to the LAN so the diagram does not imply WAN denial.
The README consumes the generated documentation assets, with no SVG checked in.
