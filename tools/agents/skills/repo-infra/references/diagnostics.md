# Read-only infrastructure diagnostics

Use this procedure to audit a deployed host or component against its checked-in
desired state. It is read-only and never authorizes a repair.

## Establish the intended state

Read the root `AGENTS.md`, the owning `README.md`, the component's
`BUILD.bazel`, `al.lua`, and the relevant files under `ansible/`. Consult the owning component's README and Mermaid source under `docs/` when
the topology matters; ownerless historical diagrams live in `infra/archive`.

Treat the checked-in Terraform and Ansible as desired state. Do not infer
external reachability from a listening socket alone: account for its bind
address, active interfaces, routing, and every applicable firewall zone.

## Inspect without changing the host

Find the owner's SSH target in its README and BUILD declarations, then confirm
the label with a bounded `bazel_agent bazel query` when uncertain. Run it with
`bazel_agent bazel run //<owner>:ssh -- '<read-only command>'`; names may have
environment suffixes, such as `//infra/download:ssh.local`. XCP-ng provides
`//infra/xcp_ng:ssh`. Inspect the actual target before running it: a target
named `inspect` can include writable probes and is not inherently read-only.

Keep AL configuration, Vault references, injection plugins in `data`, and
selected `run_args` labels aligned. Prefer `//tools/ssh:defs.bzl`'s `ssh_binary`
for direct SSH; password-authenticated targets may use `al_binary_run` with
the established `sshpass` environment mechanism. Never put passwords in command
arguments or print injected credentials. Preserve strict host-key verification.

If a target is missing, adding one is a repository change requiring user
authority and an owner-local OpenSpec record. Reuse the component's existing
authentication rather than creating a second credential path or modifying the
host. Package and document the target before using it. Adding a target does
not authorize the remote command; keep investigation read-only unless the user
explicitly requests a scoped live change.

For XCP-ng hardware checks, use `xl info` for total and free physical-host RAM;
`free` in dom0 reports only the control domain's allocation. Inspect PCI devices
and existing VM assignments before concluding that a GPU is available to a guest.

- Inspect services and listeners with read-only `systemctl` and `ss` commands.
  Map listeners to child processes as well as the primary unit; follow the
  [launcher checks](ansible.md#validate-without-deploying) when interpreting
  invocation flags.
- Query active firewalld zones, services, ports, policies, and rich rules.
- Compare permissions with `stat`; inspect environment variable names only, not
  their values.
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
the user also requests implementation, follow the relevant implementation
procedure and validate before deploying.
