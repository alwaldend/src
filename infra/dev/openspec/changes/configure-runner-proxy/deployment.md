## Authorized rollout on 2026-11-17

The user authorized applying the prepared Vault read grant and secure1 runner
configuration. The deployed source was commit
`d5b193f67ff25f588969e9ab7a449b98eadc7401` on `t3/fix-forgejo-ci`.
No merge, other host deployment, proxy-account creation or credential rotation
was performed.

The owning Vault target applied only `vault_policy.src_infra_dev_runner` after
review of its interactive plan: zero additions, one in-place change, zero
deletions. The change adds read access to the existing proxy secret. A subsequent
plan of the same resource reported no changes. The Vault owner's
[receipt](../../../../vault/openspec/changes/allow-runner-proxy-read/deployment.md)
records its scope.

The checked-in Ansible workflow ran with
`--limit secure1.dev.alwaldend.com --tags forgejo_runner`. Its first recap was
37 successful tasks, five changed, zero unreachable, zero failed and nine
skipped. It installed `zip`, wrote the restricted runner configuration and
TLS transport units, enabled the loopback socket and restarted the worker.
Repeating the same selection produced 36 successful tasks, zero changed,
zero unreachable, zero failed and nine skipped.

## Runner-account acceptance

A Python probe read the installed runner environment into memory as
`forgejo-runner`. It emitted only status and non-secret configuration metadata;
no proxy password or credential-bearing URL was logged or retained. It confirmed
configuration mode 0640 and matching uppercase/lowercase proxy settings.

HEAD requests to all six previously blocked upstream endpoints succeeded:

| Endpoint                               | HTTP status | Connection               |
| -------------------------------------- | ----------- | ------------------------ |
| Vault 1.21.4 archive                   | 200         | Local proxy, CONNECT 200 |
| Consul 1.22.7 archive                  | 200         | Local proxy, CONNECT 200 |
| Terraform 1.14.8 archive               | 200         | Local proxy, CONNECT 200 |
| Qt 6.8.3 Qt3D archive                  | 200         | Local proxy, CONNECT 200 |
| Nordic nRF5 SDK 15.3.0 archive         | 200         | Local proxy, CONNECT 200 |
| Arm GNU toolchain 10.3-2021.10 archive | 200         | Local proxy, CONNECT 200 |

The Forgejo HTTPS request returned 200 directly from `192.168.10.40`, with
no CONNECT tunnel. This verifies the local-domain bypass on the actual runner.

CI had already removed its source checkout. A task-owned workspace under the
runner's `home/out/forgejo-ci-proxy-probe` therefore exercised the repository's
exact Bazel 9.3.1 pin and existing Vault archive URL/SRI. Bazel downloaded and
verified the archive successfully; executing the extracted binary reported
Vault v1.21.4. The disposable Bazel server was stopped after verification.

Repeatable probes and sanitized output remain in the task's ignored
`out/forgejo-ci/`: `verify_runner_live.py`, `live-bazel-acceptance.log`,
`remaining-downloads.log`, `ansible-deploy.log` and `ansible-convergence.log`.
The probes demonstrate runner readiness, not a complete new repository CI run.
The checkout, diagram and scan source changes still await the normal PR merge.
