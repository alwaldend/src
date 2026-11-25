## 1. Restore Mullvad load balancing

- [x] 1.1 Restore the inbound rule and correct the observatory selector.
- [x] 1.2 Validate the Terraform and OpenSpec configuration.
- [x] 1.3 Deploy the reviewed change to Njalla and verify the applied settings.

On 2026-11-25, the targeted apply was initially blocked by retired
`module.node_work1` entries in the shared service state. After the approved
[Work1 state reconciliation](../../../../dns/openspec/changes/consolidate-terraform-dns/tasks.md),
Terraform updated Njalla's observatory and routing in place. The first
verification plan showed the balancer missing. A second targeted apply
restored the balancer, and the final plan found no changes across the three
Njalla resources.

## 2. Correct the observatory serialization

- [x] Build the patched provider from the pinned upstream source and verify that
      Terraform can install it from the repository mirror.

An isolated Terraform run installed the Bazel-built provider, imported the
legacy tagged configuration from a loopback panel fixture, and migrated it
with one flat observatory update. The existing routing configuration survived
the update, and the next plan found no changes.
