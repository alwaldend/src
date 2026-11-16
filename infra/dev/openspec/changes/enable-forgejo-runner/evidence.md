# Runner enablement evidence

## Source and live state

The task starts at `57c3ace6636fda4d2102c4369956a47ddb932140` on the isolated
`t3/enable-forgejo-runner` branch. The verified origin is Forgejo's
`alwaldend/src`, targeting master. Related owner records:

- [Credential import](../../../../../projects/ansible_collection/openspec/changes/import-forgejo-runner-credentials/proposal.md)
- [Vault read grant](../../../../vault/openspec/changes/allow-secure1-runner-credential-read/proposal.md)

Strictly verified `//infra/dev:ssh.secure1` inspection on 2026-11-16 found
`forgejo-runner.service` not found and no runner account. Secure1 runs Fedora 44
and its 500-GiB root filesystem has approximately 499 GiB available. The existing
Forgejo bootstrap receipt records repository-scoped runner ID 2; this task does
not create a registration.

## Validation

The packaged Ansible syntax target and deployment build pass. Terraform
formatting, the actual Ansible credential acceptance fixture and all 29 targets
expanded by the selected test command pass. Semantic lint of the deployment,
SSH, Vault apply and acceptance targets passes. Each of the three owner-local
OpenSpec changes passes strict validation. The credential fixture retains a
sanitized `registration.json` receipt with initial registration, identity import,
repeat deployment, invalid-input rejection and output-suppression outcomes.

## Live operation

The scoped text plan reviewed on 2026-11-17 reports one policy addition, one
in-place development AppRole update and zero destroys. Only the new
`src_infra_dev_runner` policy is added to the existing role; it grants read on
the single secure1 KV data path, without writes, listing or sibling access.
No plan file or raw plan output is retained.

Automatic approval review initially rejected the apply before it executed.
The user then explicitly approved the exact read-only grant. The interactive
apply was reviewed again and completed with one addition, one in-place change
and zero destroys at 00:03 MSK on 2026-11-17. The scoped follow-up plan reports
no changes. Secure1's runner-only Ansible invocation successfully consumed
the record and reached package installation; no credentials appeared in output.
The worker is active/running and enabled, with zero restarts. Its credential
file is owned by the worker at mode 0600; UID 5721 has only its own group.
Sanitized journal inspection reports one successful Forgejo declaration, zero
error/fatal entries and zero TLS errors. An initial inspection summarizer had a
parenthesis typo; correcting that read-only command completed the connection
check. The first runner-only deployment completed at 00:07 MSK with 34 successful
tasks, 22 changes, zero unreachable hosts and zero failures. Android package
receipts are present for API 36, build-tools 36.0.0, platform-tools and NDK
29.0.14206865; the SDK occupies approximately 2.7 GiB. The worker reports version
13.2.1, Node.js 22.23.1 and OpenJDK 27.0.1. Runner-owned Bazel and AL executables
are present and its NDK link resolves to the selected version. Non-secret
identity inspection confirms the existing ID 2 and canonical instance address.

Repeat deployment reported 30 successful tasks and zero changes/failures. Final
review separated role installation from service activation and guarded restart
handlers so failed prerequisites cannot start a fresh worker through forced
handlers. The final candidate was deployed again at 00:13 MSK: 32 successful
tasks, zero changes, zero failures and zero unreachable hosts. This verifies
convergence of the real play; the prerequisite-failure guard was source-reviewed
without injecting a failure into the live host. Dev1 was excluded throughout.

## Current CI

This source revision contains no `.forgejo/workflows` files. Forgejo's task
listing shows eight historical runs, with the latest marked successful; that
history does not establish current workflow activation. The local Node action
launches `bazel run --config=ci //tools/ci`. The command discovers the root and
declared nested workspaces and intends to build then test each, continuing to
collect failures. Its current subprocess launcher still prefixes commands with
`al`, while `al bazel --help` exposes only the affected-test command. This
pre-existing CI launcher defect remains outside runner enablement.

## Session ergonomics

Discovery confirmed that historical runner documentation is not live state.
The initial broad worktree inventory and unfiltered formatter output were
avoidable large reads; later checks use bounded log summaries. No shared skill
or host configuration was changed to address those observations. The owner-local
tasks and this record preserve the exact deployment boundary for continuation.
