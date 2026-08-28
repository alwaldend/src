---
name: host-bot-diagnostics
description: >-
  Audit and troubleshoot Codex, T3 Code, Traefik, firewall, and deployed
  configuration for users/simeonwarren/host_bot using read-only checks. Use
  for listener exposure, mTLS routing, permissions, environment leakage, or
  deployed-versus-Ansible drift; use repo-ansible for implementation.
---

# Host Bot Diagnostics

## Establish the intended state

Read the root `AGENTS.md`, `users/simeonwarren/host_bot/README.md`, the
component's `BUILD.bazel`, `al.lua`, and relevant files under `ansible/`.
Consult `infra/arch/README.md` and `infra/arch/arch.drawio` only when the
topology matters.

Treat the checked-in Ansible configuration as desired state. Do not infer
external reachability from a listening socket alone: account for its address,
active interfaces, routing, and every applicable firewall zone.

## Inspect without changing the host

- Inspect services and listeners with read-only `systemctl` and `ss` commands.
- Query active firewalld zones, services, ports, policies, and rich rules.
- Use `codex doctor --json`, `codex mcp list`, and resolved configuration to
  check Codex health and integration drift.
- Compare permissions with `stat`; inspect environment variable names only.
- Compare observed state with the repository and distinguish intended
  exceptions from unexplained drift.

Do not deploy, restart services, reload firewalld, modify configuration, or
contact a mutating endpoint during diagnosis.

## Protect sensitive data

Never print token values, environment values, authentication files, Vault
responses, certificates, private keys, or rendered secret-bearing
configuration. Report only variable names, file modes, status codes, and
redacted structural facts.

## Report findings

For each finding, state the intended state, observed state, impact, and
narrowest remedy. Include every verification command and its actual result. If
the user also requests implementation, follow `repo-ansible` and
`repo-secrets` as applicable and validate before deploying.
