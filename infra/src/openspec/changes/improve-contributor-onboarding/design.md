## Context

Deferred; implementation has not started. Authority is planning only. The
current candidate remains `e0d24cb3a`; [proposal evidence](proposal.md) records
the root entry point and unconditional Java truststore startup argument.
No fresh-host build or supported-platform matrix has been validated here.

## Goals / Non-Goals

**Goals:** Make one supported public workflow reproducible and keep custom
trust assumptions with their host owner while preserving certificate checks.

**Non-Goals:** Promise every operating system, change product behavior, add CI,
replace the toolchain, or alter this host's live trust configuration.

## Decisions

- Use the root README as a concise entry point that links the owning setup
  documentation. Extend that setup owner after locating its checked-in source;
  copying the external page into another guide would create two authorities.
- Separate standard public trust from explicit custom-CA setup. Inspect
  existing local Bazel configuration hooks before choosing the mechanism;
  hardcoding another distribution's path would preserve the same assumption.
- Preserve the current host's private-CA behavior through declarative host
  configuration if the shared argument moves. Persistent host changes must
  be mirrored in `users/simeonwarren/host_bot`; implementing this plan does not
  itself authorize deploying that configuration.
- Choose one small public build with a representative inspectable output as
  the smoke workflow. Documentation claims only the host configurations
  actually exercised, rather than inferring portability from configuration.
- Link the Anki onboarding change instead of moving its engine documentation
  into repository setup. Both records can proceed independently.

## Risks / Trade-offs

- Removing the shared truststore override can break private-CA users -> validate
  the explicit custom-trust path before removing the old assumption.
- A smoke target hides important prerequisites -> select a representative
  public target and record exactly what it covers and excludes.
- Host labels become stale -> retain the tested host/environment and candidate
  revision with the repeatable result rather than a broad unsupported claim.

## Migration Plan

Next action after implementation is requested: locate authoritative setup
content, inspect existing host overrides, and record the supported baseline.
Prepare a documentation/configuration candidate; run the documented workflow
in a fresh supported environment without the Fedora-style pathname and verify
custom trust in an isolated configuration. Save the actual command, revision,
and expected output as acceptance evidence. CI and live host deployment remain
outside this change.
