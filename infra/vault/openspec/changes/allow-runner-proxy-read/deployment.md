## Authorized application on 2026-11-17

The user authorized the prepared policy change and secure1 deployment. From
source commit `d5b193f67ff25f588969e9ab7a449b98eadc7401`, the owning command was:

```sh
bazel run //infra/vault/tf:tf.apply -- -target=vault_policy.src_infra_dev_runner
```

The interactive apply displayed zero additions, one in-place change and zero
deletions. The reviewed diff added only `read` on
`secrets/data/alwaldend.com/vault1/approles/src_infra_threexui/http_proxies`,
preserving the existing runner-identity read. The agent answered `yes` only
after reviewing that plan. Apply completed successfully at 11:25:54 MSK.
No saved plan, JSON plan/state output or raw plan log was retained.

The same owning plan target with the same resource selector subsequently
reported no changes. This establishes convergence of the selected policy;
unrelated Vault resources were outside this operation's scope.
The development deployment then successfully read the selected existing
credential, and runner-account download and bypass checks passed. See the
[runner rollout](../../../../dev/openspec/changes/configure-runner-proxy/deployment.md).
