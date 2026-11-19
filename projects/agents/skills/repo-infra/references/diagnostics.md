# Read-only infrastructure diagnostics

Use this procedure to audit a deployed host or component against its checked-in
desired state. It is read-only and never authorizes a repair.

## Establish the intended state

Read the root `AGENTS.md`, the owning `README.md`, the component's
`BUILD.bazel`, `al.lua`, and the relevant files under `ansible/`. Consult
`infra/arch/README.md` and `infra/arch/arch.drawio` only when the topology
matters.

Treat the checked-in Terraform and Ansible as desired state. Do not infer
external reachability from a listening socket alone: account for its bind
address, active interfaces, routing, and every applicable firewall zone.

## Inspect without changing the host

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
