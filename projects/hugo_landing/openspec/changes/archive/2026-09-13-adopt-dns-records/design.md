## Deployment boundary

Use this project's existing
[Terraform stage](../../../../tf/README.md) and canonical
[DNS declaration](../../../../dnsconfig.json), preserving its scoped access and state.

The shared [cutover procedure](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns adoption ordering, writer coordination, import, verification, and recovery.
Keep ownership enabled after adoption because disabling it would propose deletion.

## Adoption evidence

At `2026-09-13T01:46:24.253940+00:00`, adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled
`tf/dns_variables.tf` SHA256
`5371c2d0418f97a85875084861c7c451f182137a88f09af3a3d50657b5c10ef2`.
The scoped Vault stage created nine AppRole resources. The project imported its
one existing Cloudflare record, `5518fad3b442824c1a1c71453b958cb2`; the reviewed
plan and apply changed no managed resources. All 80 public-zone records remained
equal as complete JSON objects keyed by provider ID, with no additions, removals,
or changes.

At `2026-09-13T01:50:12Z`, public CNAME resolution for
`hugo-landing.alwaldend.com` through `1.1.1.1` returned `alwaldend.github.io`.
The preserved provider record has TTL 600 and proxying disabled. Raw plans,
provider inventories, and credentials remain in private ignored task scratch.
