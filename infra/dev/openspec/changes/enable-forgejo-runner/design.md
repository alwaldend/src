## Context

Forgejo Terraform already created repository runner ID 2 and cached its identity
under the Forgejo AppRole's `runners/secure1` secret. Secure1 is Fedora with a
500-GiB root disk, eight CPUs and 16 GiB RAM. No workflow is currently checked in.

## Goals / Non-Goals

Enable the existing worker identity and install repository build prerequisites.
Do not activate workflows, replace the VM, create registrations or grant jobs
deployment credentials.

## Decisions

- Give the development deployment exact-path read access to the cached worker
  credentials through Vault Terraform. Inject the record only into the Ansible
  controller process. Jobs retain only their own Forgejo worker identity.
- Extend the existing runner registration owner to accept API-issued identity
  data when no local registration exists. Retain interactive registration for
  existing callers; never overwrite an existing worker file.
- Use the existing development package, Bazel and Android task entry points for
  the runner account between the runner's install and service phases. Start a
  fresh worker only after prerequisites are ready. One host-executor slot keeps
  repository jobs serialized.
- Package secure1 SSH inspection with its existing signing flow and canonical
  server CA, retaining strict host verification.

## Risks / Trade-offs

- Host jobs share one service account and persistent filesystem. Restrict this
  repository-scoped worker to trusted code; workflow admission is separate work.
- Invalid/mismatched cached identities must fail before writing worker data.
  Existing local identities must survive redeployment without bootstrap input.
- Initial SDK installation needs outbound downloads and can take several minutes.

## Migration Plan

Validate packaging and syntax, inspect secure1, review and apply only the exact
Vault read policy and development identity binding, then run runner-tagged
Ansible limited to secure1. Verify systemd, identity permissions, tool versions,
Forgejo connectivity and a second idempotent deployment. Preserve sanitized
evidence. Rollback requires an explicitly authorized service stop through IaC;
registration and credentials remain intact.
