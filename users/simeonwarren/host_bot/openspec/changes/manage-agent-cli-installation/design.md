## Context

See proposal.md for motivation. Host Bot already defines the managed user's
home through `dev_vm_bazel_home` and owns CLI deployment in its playbook.
Codex 0.164.0 is installed on the live host; Grok is absent.

## Goals / Non-Goals

Provision missing commands through the user-supplied upstream installers.
Keep provider selection, login, and existing Codex configuration unchanged.

## Decisions

- Add two user-scoped shell tasks to the existing Host Bot playbook, beside
  its other CLI installations. A shared role would add an unnecessary owner
  for this host-only integration.
- Use the exact official installer pipelines with Bash `pipefail`, ensuring
  a curl failure cannot be hidden by a successful interpreter exit.
- Set each upstream install-directory variable to the managed user's
  `~/.local/bin`, and set Codex's non-interactive mode. Guard each task with
  its installed command path and expose `codex_install` and `grok_install`
  tags. Existing installations are adopted without an upgrade.
- Retain the upstream installers' current-release selection as explicitly
  requested. Pinning downloaded release archives would depart from the
  supplied installation method.

## Risks / Trade-offs

- Mutable upstream installers and release selection limit reproducibility;
  use the explicitly approved publisher URLs and record observed versions.
- Installer-owned files, completions, and shell integration follow upstream
  behavior; their provisioning is declared by these Ansible tasks.
- A missing curl executable or failed download must fail visibly. The
  installers require the host's existing curl and Bash executables.

## Migration Plan

Build and lint the packaged playbook, inspect the selected task list, then
run only the two installer tags on Host Bot as authorized by the user.
Record command versions and verify the installed-command guards through
Ansible's shell module check mode. Keep the logs in ignored task scratch.

The normal deployment wrapper was attempted and stopped before task execution
because its Vault AppRole secret-ID request returned HTTP 403. The user also
explicitly authorized running the supplied installer command directly. That
exact exception was used only to install Grok for the current Host Bot user;
the installation task above records the reconciliation into checked-in IaC.
No Vault configuration or permissions were changed. Codex was already present.
