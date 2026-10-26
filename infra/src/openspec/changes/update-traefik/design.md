## Context

`third_party/io_traefik_traefik/binary_toolchain.json` owns the binary archive,
integrity, and execution constraints. The shared Ansible role packages it directly.
See proposal.md for the requested extraction.

## Goals / Non-Goals

Keep the independently reviewable dependency change compatible with the existing
Linux AMD64 consumer. Preserve all role and routing changes in PR #202.

## Decisions

Retain the verified 3.7.16 selection and checksum from the original upgrade.
Do not select a different release while extracting that change. Put
`exec_compatible_with` under `toolchain`, as required by the current lock schema.

Create an independent branch from master containing only the lock and this
maintenance record. After publication, remove these changes from PR #202.
Keep the existing live 3.7.16 service running throughout the source split.

## Risks / Trade-offs

Patch releases tighten request handling. The original upgrade exercised actual
routing; independently build and run the binary and package its existing Ansible
consumer before publishing this extraction. Historical deployment evidence stays
with the download change and does not authorize another deployment.

## Migration Plan

Validate the dependency-only candidate and publish it against master. Merge the
version PR before redeploying from PR #202 so a build does not restore the old pin.
The source split itself requires no host changes.
