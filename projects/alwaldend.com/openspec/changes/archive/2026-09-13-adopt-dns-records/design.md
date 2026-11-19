## Deployment boundary

Use this project's existing
[Terraform stage](../../../../tf/README.md) and canonical
[DNS declaration](../../../../dnsconfig.json), preserving its scoped access and state.
The existing stage also contains service infrastructure. The owning
`//projects/alwaldend.com/tf:dns.plan`, `dns.show`, and `dns.apply` wrappers
target `module.dns` in the same root and backend, using `dns=1` without starting
Proxmox login. Apply consumes the reviewed saved plan. This workflow verifies
DNS and its dependencies; it does not establish overall PVE service health.

The shared [cutover procedure](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns adoption ordering, writer coordination, import, verification, and recovery.
Keep ownership enabled after adoption because disabling it would propose deletion.

The [shared writer audit](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md#writer-audit)
owns writer-control evidence and its observation limits for external hosts,
external CI, and historical checkouts.

## Earlier blocked attempt

At `2026-09-13T01:37Z`, deployment was paused before any DNS import or source
activation. The reviewed Vault prerequisite apply attached only the project's
DNS policy to its existing role, with no policy removals. It also created the
11 passive DNS policies selected by the shared policy-module dependency.

The ordinary Terraform wrapper required Proxmox authentication. Its configured
host `192.168.1.216:8006` had a direct route but failed neighbor discovery;
the separately declared IPv6 address failed the same way. This prevented TLS
authentication before Terraform could plan. Checked-in inventory provided no
alternative node endpoint. That attempt left DNS ownership disabled and the
records unchanged. The scoped workflow subsequently resolved the DNS adoption
blocker while retaining the existing service root and VM configuration.

The earlier Vault changes above are retained as previously recorded history.
Current existing-role access is verified by the successful authenticated scoped
DNS operations; the earlier prerequisite plan was not re-audited for finalization.

## Completed scoped adoption

At `2026-09-13T02:45:35.476602+00:00`, adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled
`tf/dns_variables.tf` SHA256
`5371c2d0418f97a85875084861c7c451f182137a88f09af3a3d50657b5c10ef2`.
The reviewed scoped plan contained two exact Cloudflare imports and two managed
`no-op` actions. Applying that saved plan completed both declarative imports
with zero additions, changes, or deletions. The follow-up scoped plan contained
two DNS `no-op` actions and no pending imports. All 80 Cloudflare records
remained equal as complete JSON objects keyed by provider ID.

At `2026-09-13T02:47:22.059808+00:00`, post-adoption A and AAAA queries for
`pages.alwaldend.com` through `1.1.1.1` matched the declared addresses, without
errors or truncation. The DNS evidence is bound to the adoption receipt and
canonical declaration by their SHA256 hashes.

The unchanged declaration reuses the baseline ownership-linter evidence.
The source default/hash and scoped plan checks cover activation; final aggregate
source quality gates remain with repository delivery. Raw plans, import maps,
inventories, and credentials remain in private ignored task scratch.
