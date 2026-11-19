# Writable fixture inventory

Read-only survey of the checked test suites, bounded to deterministic
writable-sandbox or hermetic in-process fixtures covering agent-executed
high-risk trajectory classes.

| Trajectory | Representative evidence | Coverage |
| --- | --- | --- |
| Git (rebase, lease push, remotes) | `tools/repo_delivery/cmd/repo_delivery/git_integration_test.go` and `delivery_integration_test.go` | Lease rejection on concurrent remote change, push rollback on base advance, rebase conflict cleanup, rebase CAS against concurrent head updates, signature preservation, sibling-ref and worktree isolation, SSH-vs-HTTPS endpoint policy, receipt binding to destination/branch/clone |
| Forge / GitHub adapter | `tools/repo_delivery/cmd/repo_delivery/github_test.go`, `review_test.go`, in-memory `integrationForge` | Preflight, PR list/create/edit/view assembly, strict JSON decoding, concurrent metadata-edit refusal, ambiguous create/update reconciliation, expected-head postverification, review-thread reply/resolve, reviewer requests, idempotent retry |
| Bazel execution | `projects/bazel_agent/cmd/bazel_agent/main_test.go` | `--config=agent` injection/ordering, argv assembly, lookup-failure propagation, bounded doctor report, stale-install detection, task-scratch namespace enforcement |
| Runtime (MCP/Cordis) | `projects/mcp_cordis/test/runtime_test.mjs`, `stdio_test.mjs`, `exec_test.mjs` | Full server lifecycle, HMR hot-reload and recovery, plugin invocation, payload limits, process-tree cleanup, concurrent scratch isolation, stdio end-to-end, group-kill timeouts |

## Recorded coverage gaps

- **Terraform:** the existing `//tools/terraform:explicit_operations_test`
  is static only. No offline runner fixture exercises `terraform_binary_map`
  plan plumbing, backend-config injection, or operation policy in a writable
  sandbox.
- **Ansible:** no test fixtures exist in `tools/ansible` or any infra
  Ansible package; only live `al_binary_run` playbook targets are defined.
- **Vault / secret injection:** no test targets in `tools/vault`;
  credential redaction in git errors (`runner_test.go`) is the nearest
  hygiene proxy. Vault injector, tf_backend, login helpers, and environment
  cleaning lack writable fixtures.
- **Bazel agent end-to-end:** command assembly and doctor are well tested;
  no isolated end-to-end workspace execution fixture exists.

These gaps are explicit per the `writable-fixtures` criterion: existing
deterministic fixtures are inventoried and every missing trajectory class is
recorded rather than claimed as covered.
