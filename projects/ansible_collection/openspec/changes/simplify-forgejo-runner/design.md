## Context

See [proposal.md](proposal.md). No checked-in deployment enables this role.
The `infra/dev` consolidation already retired the old worker identity and token
lookup. Its remaining Vault/OIDC configuration owns SSH, DNS and VM provisioning,
not Actions registration. The shared Forgejo OIDC plugin serves other consumers.

## Goals / Non-Goals

Keep the installed default layout and credential protections while reducing the
role's supported inputs. This is source-only work; it does not enable a worker,
change live authentication or modify shared Forgejo/Vault provisioning.

## Decisions

- Derive internal paths at their uses from `forgejo_runner_dir`. Removing unused
  path overrides is preferable to maintaining compatibility aliases with no callers.
- Keep labels and capacity explicit; serialize a single advanced configuration
  mapping, merging managed fields last. This preserves path agreement without
  maintaining a template expression for every upstream option.
- Put one-time registration in a task entry point shared by installation and
  the acceptance fixture. Keep a short rescue message, stdin input, `no_log`,
  the existing async timeout and the protected registration file. Parsing human
  error strings adds little operational value and retains secret-bearing results.
- Keep the package-management switch for caller-owned OSTree installation.
  Container runtimes remain a caller decision; do not install Podman for host jobs.
- Use local Ansible with a disposable registration executable to exercise real
  tasks and template rendering. Do not add test-only role inputs or invoke a host
  package manager, create users, call systemd or contact Vault/Forgejo.

## Risks / Trade-offs

- Removed inputs break external users → document migration; repository search
  found no active caller, and default installed paths remain unchanged.
- Fixed registration errors give less diagnosis → document checking transport,
  TLS, privilege-switch prerequisites and token validity separately; never
  suggest printing the secret-bearing registration command output.
- Local fixtures cannot establish remote registration, systemd or cross-account
  ACL behavior → explicitly report these coverage limits.

## Failure Cases Before Implementation

Acceptance must cover initial registration, repeated application with undefined
bootstrap inputs, missing initial inputs, a command that prints its token and
fails, a command that exceeds the registration deadline, protected registration
mode, default configuration and an alternate layout, advanced nested settings and attempts to
override managed configuration fields. Preserve the service restart wiring by
source review and render the exact unit. Produce sanitized JSON evidence.

## Migration Plan

Replace removed path overrides with `forgejo_runner_dir`. Replace advanced
timeout/cache scalar inputs with upstream keys under `forgejo_runner_config`.
Deployments that need Podman must supply it through their package owner.
No live migration or re-registration is performed by this change.
