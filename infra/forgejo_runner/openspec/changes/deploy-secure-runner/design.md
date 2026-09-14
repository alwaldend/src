## Context

See [proposal](proposal.md). The baseline already owns a Vault AppRole and an
adopted DNS record, but its Ansible runner role is inactive. Reuse those owners
and the shared Fedora XCP-ng cloud-init output. Deployment credentials stay on
the controller; CI has a separate JWT identity.

## Goals / Non-Goals

Provision and configure the requested runner through the owning Terraform and
Ansible targets. Preserve unrelated resources and DNS state. Avoid changes to
the controller host or shared authentication configuration.

## Decisions

**Proceed** with the existing component and shared roles. A runner map provides
future expansion without another project or copied deployment. One root disk
lets cloud-init grow the filesystem used by both tools and workspaces. The
operator confirmed the old VM was absent; retire its stale state without a
destroy request and replace its DNS record with `secure` at a separate address.

**Revise** the initial protected-runner assumption. Forgejo's native runner
model has no protected-branch admission setting. Its security guide explicitly
describes unprotected branches targeting a repository runner. Furthermore,
the pinned 15.0.3 source hardcodes `ref_protected=false` in
`services/actions/context.go` and copies it into Actions JWT claims. Binding
Vault to that claim being true would reject every job. Bind signed repository,
branch, event and workflow claims to branches protected by the owning IaC.
Keep exact refs: Forgejo compiles branch-protection globs with `/` as a separator
in `models/git/protected_branch.go`, whereas Vault globs can match that separator.
A Vault `releases/*` wildcard would admit nested branches outside the existing
Forgejo protection. The workflow reads the candidate's CI configuration and
skips authentication checks on release branches outside the exact allowlist.
This constrains Vault credentials, not runner scheduling. Host execution is
appropriate only for trusted repository writers; workflow filters do not form
a security boundary. The user was offered a dedicated fully protected CI
repository as an alternative to the existing source repository.

Primary evidence inspected on 2026-09-14:

- [Forgejo runner security](https://forgejo.org/docs/latest/admin/actions/security/).
- [Pinned context source](https://codeberg.org/forgejo/forgejo/src/tag/v15.0.3/services/actions/context.go).
- [Pinned JWT source](https://codeberg.org/forgejo/forgejo/src/tag/v15.0.3/services/actions/auth.go).
- [Forgejo 15 OIDC reference](https://forgejo.org/docs/v15.0/user/actions/security-openid-connect/).

## Risks / Trade-offs

Saved plans and raw deployment logs remain private task scratch. Require no
unexpected changes or replacements. Repository writers can execute host jobs
on the runner even when Vault denies their identity.

## Migration Plan

1. Validate and apply the runner's Vault SSH access and CI JWT resources.
2. Bootstrap its XCP-ng OIDC identity, then apply its resource-set allocation.
3. Review and apply VM/DNS provisioning through the existing owner backend.
4. Initialize its repository registration credential in Vault with the packaged
   controller playbook, then deploy the VM with the packaged deployment.
5. Check service health, idempotency, resource capacity and a live CI run.

Rollback starts by stopping the runner through reviewed Ansible configuration;
VM destruction, unrelated changes and shared credential replacement require
their own explicit scope.

## Evidence and next action

Observed on 2026-09-14: the component AppRole authenticated and could update its
own config KV path. Seven initial Terraform formatting, DNS and Buildifier
checks and both Ansible syntax checks passed. The reviewed Vault plan applied
three additions and three updates with no deletions. XCP-ng resource-set access
was updated, and the runner VM was provisioned with eight CPUs, 16 GiB RAM and
a 500 GiB disk. The DNS plan added `secure` and deleted the stale `runner1`
record. A lookup through the owning resolver returned 192.168.10.101 for the
new hostname, and TCP port 22 was reachable.

Registration initially failed because the shared Forgejo login plugin ignored
the selected call's named authentication override. The owning helper now merges
that override with its defaults, with passing regression tests. Registration
then succeeded and stored the credential through the packaged playbook. Next
complete repeat-deployment verification and live CI testing. The host baseline
and development tools deployed successfully. Initial on-host registration
failed because `acl` was absent, preventing Ansible from becoming the runner
account. The package declarations now include it, and the role reports only
safe failure categories. The scoped retry registered the runner, enabled its
service and installed its Bazel tools (26 successful tasks, no failures).

The controller did not trust the new hostname. The initial deployment used a
task-local host-key pin; normal deployments now package the existing public
server CA and generate inventory-scoped trust under ignored operational cache.
The full Terraform reconciliation plan reports no changes and confirms the
requested allocation.

The full repeat deployment used strict CA verification and succeeded: 209 host
tasks, 19 changes in shared baseline configuration, no failures; initial runner
registration was skipped and service health passed. Candidate
`d68b08979c7b889690f72caa26ffa05f15c66e96` passed both delivery validation checks
(root quality and targeted tests, then semantic lint for all affected packages).

Its protected validation branch produced
[CI run 1](https://git.alwaldend.com/alwaldend/src/actions/runs/1), task 1.
The job log identifies `secure` running runner v11.0.0, then reports schema
errors for the `forgejo` expression context before executing any Python.
The workflow now uses the supported `github` compatibility context. Source
inspection also found missing native OIDC request environment support in v11;
the dependency is upgraded to the publisher-checksummed v13.1.0. That binary's
offline workflow validation passes. The upgrade is deployed.

The controller login expired during deployment and the user renewed it.
The scoped Ansible upgrade then succeeded without repeating registration;
an authenticated repository inventory confirms exactly one idle runner named
`secure`, label `secure`, version v13.1.0.

The upgraded source tree `409377b38fc64e27553ab8315f35691fe94e869e` was tested
by validation commit `2024d4e24a06efd618a3ebfe0555181c561f4387`, a fast-forward
child of the first test commit. [CI run 2](https://git.alwaldend.com/alwaldend/src/actions/runs/2)
passed resource checks but returned HTTP 400 during authentication. The
workflow now reports fixed error categories without response text or claims,
requires an audience-specific rejection for its negative check, and rejects
HTTP redirects before credentials can reach another endpoint. The next action
is to validate and test this diagnostic candidate, then fix the observed cause.

Provider v5.8 splits CSV alternatives into arrays by default; the live Vault
role confirms all expected claim arrays. The parsing setting is explicit for
clarity. This is not an authentication-policy change.

Validation retries use a commit with the validated source tree and the observed
prior validation commit as its parent. This keeps the protected test branch
fast-forward while delivery maintains one aggregate feature commit. Inspect the
source difference first, record source/tree/validation OIDs, and use the helper's
explicit prior-commit lease. Archive only after successful live CI verification.

## Operator confirmation and state retirement

The user confirmed that the legacy Proxmox VM is gone. The new secure VM had
already been created at 192.168.10.101, so it retains that address. A declarative
removed block records the old module's retirement with destruction disabled.
The provider still demanded credentials for the retired service while planning
that block. After saving a private state backup, the owning Terraform state
command removed only `module.vm["runner1"].proxmox_vm_qemu.vm`; it reported one
removed instance. No VM destruction was requested. The stale DNS record was
retired through the reviewed DNS plan. Proxmox provider/login dependencies are
removed from the current component.

## Session ergonomics review

- `forgejo-login-call-overrides` (live and fixture-tested): named authentication
  was ignored, producing a registration authorization failure. The owning
  helper fix has a regression test; successful registration confirms routing.
- `runner-unprivileged-acl` (live): suppressed registration output initially
  hid a missing prerequisite. Safe category diagnostics identified permission
  setup, a read-only probe confirmed missing `setfacl`, and the declared package
  fix made registration succeed.
- `runner-controller-host-trust` (live): the controller lacked CA trust for a
  newly provisioned hostname. The deployment now owns its task-local trust
  file, avoiding manual controller configuration.
- `runner-oidc-version-contract` (live and source-verified): a v11 workflow
  failed before execution, and that runner lacks OIDC environment support.
  The pinned upgrade adds an offline workflow-schema validation command before
  another live test.
- One cold build exceeded its initial deadline; the longer warmed retry passed.
  Buildifier checks caught declaration ordering, repaired with the owning
  formatter. No shared skill or host configuration changes are needed.
