# Project DNS cutover

This source change removes central DNSControl write entrypoints and prepares
owner Terraform integrations. It does not establish deployed ownership.
`dns_enabled` remains false before adoption. Provider initialization can still
require credentials; it is not an offline check.

## Prepare the owner

Run `bazel_agent bazel run //infra/dns:lint` to discover the current declarations
and their owners. Select each owner's `tf_setup` root if present, otherwise
`tf`; nested owners use their documented adapters in `//infra/dns/projects`.
Their local OpenSpec changes retain the source and bootstrap requirements.
No checked-in ownership registry or export step is required.

Before operational commands, review and separately authorize the owning Vault
Terraform changes. New role modules live under
`infra/vault/tf/approles/<name>` and do not join the broad infrastructure group.
The existing `cloudflare.com/dns_token` entry supplies the Cloudflare token.
If its optional `cloudflare_zone_id` is absent or empty, the shared Terraform
module resolves the configured zone name; verify the token permits listing
and reading that zone. The existing DNS MikroTik entry supplies `username`
and `password`. DNS reuses the configured ingress RouterOS endpoint through
the shared infra AL configuration. Credential entries need no additional
metadata fields, and secret values remain in Vault.

PVE service roots also select the endpoint-only AL injector. The provider
validates `PM_API_URL` before target pruning when its Terraform block omits
`pm_api_url`. AL reuses the existing configured PVE URL for that validation;
the DNS selector does not start PVE login or supply PVE credentials. The
[offline regression](../../../../../../projects/tf_modules/dns_records/test/providers/README.md)
checks targeted import and saved-plan apply with existing service state and
zero PVE API requests. Required root variables still use their existing Vault
readers; the DNS commands retain the same owning root and backend.

Endpoint-owning scopes require an independent authenticated recovery path.
The configured MikroTik hostname does not itself establish recovery access.
An independent management path must remain reachable without the managed router
name and have a valid matching TLS identity. Do not assume the existing
certificate covers an IP address or disable TLS verification. Verify Vault
and backend access independently before adopting their DNS names too.
Record the authenticated recovery tests in the coordination change's
[operational evidence](operations.md). Existing certificates and networking
need no changes when those tests succeed.

A global-only landing owner such as `projects/sri` is the smallest pilot after
its AppRole/backend is deployed. Continue with other global-only owners,
non-endpoint infrastructure, common apex/mail, and endpoint-owning batches
after their independent recovery prerequisites pass.

## Adopt one owner

1. Audit active and scheduled DNS writers before any import, stop identified
   competing writers, and coordinate one owner at a time. Record the audit's
   coverage and unavailable observations. Removing commands from this
   candidate does not disable an older checkout, external scheduler or manual
   writer. Keep identified competing writers stopped throughout adoption;
   remaining records stay live during the transfer.
2. Obtain an explicitly authorized live inventory. Match declarations by view,
   name, type, value, TTL and priority. Select enabled, non-dynamic RouterOS
   rows. Retain sanitized record IDs and comparisons; keep raw state/plans in
   private ignored task scratch.
3. Set the owner's `dns_enabled` default permanently to true in its reviewed
   transfer revision. Run source tests and the runtime ownership linter. Keep
   resource changes stopped while the existing records are adopted.
4. For a service root, populate its optional `dns_cloudflare_import_ids` and
   `dns_routeros_import_ids` maps in a private JSON variable file. Each map keys
   the normalized logical/type/view record key to the actual provider ID.
   Roots declare only their applicable maps, each defaulting to `{}`. The
   checked-in import blocks select the owning module addresses. Cloudflare
   uses `<zone_id>/<record_id>`; RouterOS uses its exact `*HEX` row identifier.
   Run `:dns.plan -- -var-file=<absolute-private-file> -out=<absolute-plan>`
   and inspect that file through `:dns.show -- -json <absolute-plan>`.
   These targets retain the same owning root and backend and already target
   `module.dns`. They select `dns=1` alone; adding the ordinary service labels
   would also start unrelated authentication.
5. Require every planned import to match its inspected provider ID and every
   resource action to be `no-op`: no additions, removals, replacements or
   attribute changes. Resolve discrepancies before proceeding. Apply only the
   reviewed file through `:dns.apply -- <absolute-plan>`, then require another
   targeted plan to show no changes. The apply target rejects bare apply,
   additional arguments, and nonempty `TF_CLI_ARGS` or `TF_CLI_ARGS_*` values.
   Compare unrelated state and the complete provider inventories before and
   after adoption. Retain a sanitized receipt with imported IDs, owning root
   and reviewed source revision in the owner change; remove private artifacts.
6. Permit only the adopted revision. Test scoped update/deletion and unrelated
   record preservation in an explicitly authorized environment. Never restore
   `dns_enabled=false` after import: it would propose deleting the bindings'
   live records. Repeat with separate receipts for each owner.

If fresh complete provider inventories prove a declared record is absent,
record that evidence separately from imports. Review an additions-only saved
plan for exactly those missing declarations, preserving every existing record.
Apply the inspected file, then verify a no-change plan and complete provider
inventories. An incomplete inventory or failed lookup does not prove absence.

Roots containing only DNS can continue using their explicit `:tf.import` or
documented adapter with actual IDs, followed by a reviewed no-change plan.
The shared module's `import_addresses` output owns the relative addresses;
prefix them with the caller module address. Global-only roots call
`module.dns.cloudflare_dns_record.records`; combined roots call
`module.dns.module.global.cloudflare_dns_record.records` and
`module.dns.module.dc1[0].routeros_ip_dns_record.records`. Never guess IDs or
use ambiguous name-based RouterOS imports for multivalue names.

Targeting is limited to the DNS operation and its dependencies. It does not
validate the health of the services elsewhere in the root. Their ordinary
Terraform wrappers retain their original labels and authentication behavior.

## Rollback

Stop both writers and reconcile declarations with intentional applied changes.
Use a separately authorized Terraform state operation to relinquish bindings
without deleting live records. Only then restore the central writer from
pre-migration revision `cb4b5bd3a37d29dc778ed140ab72d05efef7816a`, reconcile its
declarations, and review its preview before resuming it. Never restore a
central writer while any project can still reconcile overlapping records.

Archive this coordination change only after every owner's adoption and
recovery evidence is complete.
