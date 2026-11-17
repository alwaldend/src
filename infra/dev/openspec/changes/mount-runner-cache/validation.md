# Validation

The VM root validates with only Vault/XO providers and has no NAS credentials.
NAS namespaces use independent roots/backends and NAS operator authentication.
Their launchers package only their own source; the VM consumes the XO repository.
Existing metadata rejection cases and the isolated storage bootstrap fixture
remain the offline behavior checks. All affected roots are included in final
provider, quality and semantic-lint validation.

The [authorized rollout](deployment.md) records the deployed NAS accounts,
namespace/cache, XCP-ng SR, VM disk, Vault metadata and guest mount. Shared NAS,
development NAS and VM full plans have no changes. Vault's four migrated roots
also have no resource changes, without a Vault restart or reload. Unrelated
XCP-ng OIDC and AppRole changes were left unapplied.

The runner is active with its private mounted cache and 700G GC. A real action
built in two output bases produced identical artifacts and a second-build disk
cache hit. The task-owned probe workspace was removed. The CI diagram's light,
dark and dark-WebP outputs were rendered and inspected; its final directed
runner-to-Forgejo edge has an arrow marker.
