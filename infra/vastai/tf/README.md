---
title: Vast.ai account configuration
description: Account settings and protected team identities
---

This root uses the personal account credential to manage account environment
settings and the team names in `../config.json`. The configuration currently
retains only alwaldend. Its identity is protected against deletion and its
funded account must not be recreated or replaced.

The account root creates no project API keys or rental contracts. The separate
[alwaldend root](../alwaldend/README.md) uses its operator-provided team credential
to manage team keys and settings. Both roots read the same `environment_variables`
section: ENABLE_HTTPS, WEB_ENABLE_HTTPS and WEB_ENABLE_AUTH are true. Teams do
not inherit the account settings, so the team root applies the same map.

```sh
bazel run //infra/vastai/tf:tf.plan
bazel run //infra/vastai/tf:tf.apply
```

Review the interactive apply plan. Reject changes to alwaldend ownership,
account identity or credit. Never save plans or use automatic approval.

The prior live test root's three objects, the test_vastai_provider team, its
generated credentials, and its Vault entry were deleted through their original
states. The two imported teams projects_vastai_terraform_prov and
test_vastai_terraform_provider remain tracked pending deletion: the provider
refuses destruction because their management credentials were not recovered.
They are absent from desired configuration; they have not been detached from
state or reported as deleted.

The existing alwaldend identity and three account variable addresses were
reconciled within the same state. Completed migration blocks were removed.
No resource was transferred between states.

`//infra/vastai/test/catalog:catalog_test` exercises both actual roots with mock
providers, verifies shared settings and full-access key declarations, and rejects
invalid key names. It also verifies public-key metadata rotation and refusal of
missing or invalid SSH metadata. It uses no live credentials or infrastructure.

## Host SSH access

The account root manages `vastai_ssh_key.host` for each `ssh_key_entities` entry
in the shared configuration. It reads the entity's `sshpubkey` from Vault at
`identity/entity/name/<entity>`. The component AppRole has an exact-path metadata
read for the configured host user; the root contains no public-key literal or
private material. The same Vault identity is injected alongside the account
credential through the existing AL configuration.

This is account-wide registration, explicitly requested by the operator.
Vast owns propagation to instances; verify access before running deployment.
An already registered key must be imported into its matching resource before
apply rather than duplicated. Updating this source does not apply policy,
register a key or modify an existing instance. Review the owning interactive
Vault and account applies separately when deployment is authorized.
