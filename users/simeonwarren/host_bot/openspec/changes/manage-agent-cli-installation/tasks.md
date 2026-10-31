## 1. Declare installation

- [x] 1.1 Add the user-scoped Codex and Grok installer tasks and documentation; verify the packaged task list, installer paths, and failure propagation.
- [x] 1.2 Validate the OpenSpec change, Ansible package, repository quality, and affected-target lint.

## 2. Deploy and deliver

- [x] 2.1 Attempt the scoped deployment, then use the explicitly authorized direct installer if the wrapper is unavailable; verify both CLI versions and preserve the logs.
- [x] 2.2 Verify the installed-command guards through the packaged playbook in shell-module check mode; both installation tasks must report no changes.

## Validation evidence

Observed on 2026-10-24:

- `al bazel build //users/simeonwarren/host_bot/ansible:ansible //users/simeonwarren/host_bot/ansible:ansible_bin` passed.
- The packaged `ansible_bin` passed `--syntax-check`; its task list with
  `--tags codex_install,grok_install` contained exactly the two installer tasks.
- Both pipelines use Bash `pipefail` so curl failures propagate to Ansible;
  each task's `creates` path matches the upstream installation-directory variable.
- OpenSpec strict validation, all 26 repository quality tests, and
  affected-package lint passed.
- The `:ansible` deployment wrapper failed before executing tasks because its
  Vault AppRole secret-ID request returned HTTP 403. The exact user-authorized
  Grok pipeline then installed successfully with `GROK_BIN_DIR` selecting
  `~/.local/bin`; it also provisioned upstream-owned completions and shell integration.
- `codex --version` reported `codex-cli 0.164.0` and `grok --version` reported
  `grok 1.0.70 (f50904eaee6c)`.
- The packaged playbook with the same tags, local connection,
  `ansible_become=false`, and `--check` reported `ok=3 changed=0 failed=0`.
  This verifies the installed-command guards, not a successful Vault-backed
  deployment. The deployment and check logs are under ignored
  `out/manage-agent-cli-installation/`.

Publication state belongs to the task-local repository delivery receipt and Git.

## Session ergonomics

- `HOST-BOT-CLI-VAULT-BOOTSTRAP` (live): one scoped deployment was blocked by
  a Vault HTTP 403 before tasks ran. Its current authentication flow also selects
  host-wide injection for a user-only installation. Consider a separately
  authorized local installer entry point if this recurs; retain the normal
  Vault-backed deployment for tasks that need it. A fresh agent can resume from
  the sanitized failure and the existing installer tasks without retrying Vault.
- `CLI-INSPECTION-OUTPUT-SIZE` (routed): one earlier Cordis result was printed
  with both text and structured representations and exceeded the output bound.
  Subsequent reads selected only structured content and narrowed paths and
  line ranges. Keep that bounded projection when resuming.
