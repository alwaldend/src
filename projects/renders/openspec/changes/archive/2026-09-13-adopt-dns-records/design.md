## Deployment boundary

Use this project's existing
[Terraform stage](../../../../tf/README.md) and canonical
[DNS declaration](../../../../dnsconfig.json), preserving its scoped access and state.

The shared [cutover procedure](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns adoption ordering, writer coordination, import, verification, and recovery.
Keep ownership enabled after adoption because disabling it would propose deletion.

## Adoption evidence

At `2026-09-13T01:52:27.922263+00:00`, adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled
`tf/dns_variables.tf` SHA256
`5371c2d0418f97a85875084861c7c451f182137a88f09af3a3d50657b5c10ef2`.
The scoped Vault stage created nine AppRole resources. The project imported its
one existing Cloudflare record; the reviewed plan and apply changed no managed
resources. All 80 public-zone records remained equal as JSON objects keyed by
provider ID, with no additions, removals, or changes.

Public CNAME resolution through `1.1.1.1` returned `alwaldend.github.io`; the
preserved provider record has TTL 600 and proxying disabled. Raw plans,
provider inventories, and credentials remain in private ignored task scratch.
