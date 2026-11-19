## Context

See [the proposal](proposal.md) for the adoption scope. The checked-in
[`tf_setup` package](../../../../tf_setup/README.md) retains enabled DNS ownership
and packages the owner's canonical declaration with the shared DNS module.

## Deployment boundary

Use `//infra/forgejo/tf_setup:dns.plan`, `:dns.show`, and `:dns.apply` with
`dns=1`, the existing setup backend, and the `src_infra_dc1_forgejo1` AppRole.
The targets scope planning to `module.dns`; review the exact saved import plan
and apply only that file. They preserve the owning state and avoid starting
unrelated service authentication. Ordinary setup commands retain their existing
credentials and provisioning behavior.

The [shared cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns prerequisites, ordering, writer coordination, credential handling, and
recovery. Keep its procedure canonical; this change records the owner's outcome.

## Decisions and risks

Import the matching existing provider identities into this owner's state so the
transfer preserves the live records. Keep `dns_enabled=true` in the adopted
revision and update the setup documentation to match that supported state;
restoring a disabled default after import would propose record deletion.

Every planned import must match its inspected provider ID and every managed
action must be `no-op`. A subsequent targeted plan verifies no-change
reconciliation. DNS adoption evidence does not establish service or VM health.

## Adoption evidence

At `2026-09-13T02:51:10.878421+00:00`, adoption used source baseline
`9ec1bf2d11863518d7a02c57693ca9c016403f65` with enabled
`tf_setup/dns.tf` SHA256
`e089f896a935fe12cfd6a6e8a0bda4e465e54b67443d417329a03c1ba8ed03e7`.
Authenticated planning, saved-plan apply, and follow-up planning succeeded through
the existing AppRole. This evidence records access, not AppRole bootstrap.

The reviewed saved plan imported six existing records with only `no-op` managed
actions. Apply imported six records, added zero, changed zero, and destroyed zero;
the follow-up plan had no imports or managed changes. Complete inventories of
80 Cloudflare records and 76 RouterOS records remained equal as JSON objects
keyed by provider ID.

| View   | Name                          | Type  | Existing provider ID               |
| ------ | ----------------------------- | ----- | ---------------------------------- |
| dc1    | `git.alwaldend.com`           | A     | `*7A`                              |
| dc1    | `forgejo.alwaldend.com`       | A     | `*3E`                              |
| dc1    | `host1.forgejo.alwaldend.com` | A     | `*3F`                              |
| global | `int.forgejo.alwaldend.com`   | A     | `a9b42028913431c8ef725560554220c2` |
| global | `forgejo.alwaldend.com`       | CNAME | `1c5772c72405e6ebdfe29cec61ab216c` |
| global | `git.alwaldend.com`           | CNAME | `c11ec2a43d9ee2b875cfcd4ac59adb7c` |

At `2026-09-13T02:51:52.760908+00:00`, the saved DNS probe report bound its six
queries to the adoption receipt and canonical declaration hashes. Global queries
through `1.1.1.1` and dc1 queries through `192.168.1.1` matched the corresponding
provider inventory: A answers were `192.168.10.40`, and CNAME answers were
`ingress.alwaldend.com`.

The [shared operational evidence](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/operations.md)
owns writer and authenticated recovery observations and their coverage limits.
These observations do not establish universal disablement of historical writers.
This is a historical DNS snapshot. The enabled source hash and three successful
scoped wrapper builds cover activation; baseline ownership-linter evidence is
reused for the unchanged canonical declaration. Final aggregate source validation
remains with repository delivery. Raw state, plans, inventories, and credentials
remain in task-private ignored scratch.
