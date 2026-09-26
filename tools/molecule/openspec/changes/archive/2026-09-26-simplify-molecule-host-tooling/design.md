## Context

The user requested fixes for PR #112 discussions
[4111302969](https://github.com/alwaldend/src/pull/112#discussion_r4111302969)
and [4111329349](https://github.com/alwaldend/src/pull/112#discussion_r4111329349).
The first rejects the third-party QEMU binary distributor and requests a
reputation warning in the dependency skill; the second questions setup complexity.
Host inspection found `/usr/bin/qemu-system-x86_64` and `/usr/bin/qemu-img`,
both supplied by Fedora packages at 10.2.2-1.fc44. No host mutation is needed.

## Goals / Non-Goals

Use reputable host QEMU, remove redundant packaging, and explain the minimum
workflow. Preserve the accepted disposable-VM, package mapping, credential
isolation, interruption, cleanup, and diagnostic guarantees. This review does
not add role changes, production deployments, host package installation, or a
new virtualization framework.

## Decisions

**Verdict: revise.** A content hash establishes artifact integrity, not the
reputation of an unrelated binary distributor. The existing distribution
packages remove that trust decision and the extra QEMU/firmware toolchain.
[QEMU's download guidance](https://www.qemu.org/download/) documents both
upstream source releases and distribution packages. Building official QEMU
inside Bazel is credible but adds its native dependency closure and firmware
build maintenance without improving the requested local test workflow.

Host QEMU is a runtime prerequisite, like OpenSSH. Generated test configuration
identifies its two paths; preflight checks executable availability. The host
installation supplies matching firmware and modules. Results retain executable
versions and hashes; host libraries and firmware are explicitly non-hermetic.
Bazel build actions remain sandboxed and offline.

Molecule's [Ansible-native approach](https://docs.ansible.com/projects/molecule/ansible-native/)
uses ordinary inventory and lifecycle playbooks. The existing wrapper supplies
Bazel package mappings and a QEMU lifecycle with bounded cleanup. Retaining
that single tested adapter avoids shifting VM creation/recovery into each
consumer. The initial documentation will lead with one test command and a short
prerequisite list; the detailed rule and recovery contract stays available.

Update the existing dependency skill with publisher-selection guidance,
distinguishing official upstream or established distribution/vendor releases
from obscure repackagers. Do not use star counts as a trust threshold. Existing
offline Promptfoo coverage gains cases for both distinctions; it validates the
harness only, not model behavior.

## Risks / Trade-offs

Host QEMU versions, libraries, firmware, and modules can differ across hosts.
Run full VM acceptance, including KVM and TCG, with the installed QEMU and retain
its identity. An unavailable executable must fail before creating a guest.
Changing the runtime input invalidates prior VM acceptance for this revision.

## Migration Plan

Extend acceptance before implementation, remove the QEMU package and module
include through the owning generator, and keep the public scenario rule intact.
Run full runner acceptance and the public smoke target, skill/configuration
validation, OpenSpec, lint, and repository quality gates. Publish the reviewed
candidate and reply to both threads through repo-delivery. PR #111 remains the
consumer follow-up after the tooling prerequisite merges.

## Acceptance evidence

Candidate `51b7ff54d8b54e36e4a45ad50c591b5e1cd993aa` passed the full
`//tools/molecule/test:acceptance_test` and public
`//tools/molecule/test:smoke_test` on 2026-09-26. Validation ran from
12:29:20 to 12:42:35 UTC with Fedora QEMU and qemu-img 10.2.2-1.fc44.
Generated configuration uses the two `/usr/bin` executables and carries no
bundled firmware input. Retained results match the host binary versions and
SHA-256 digests.

KVM smoke, simultaneous guests, repeated execution, and explicit TCG completed
all six lifecycle phases, including idempotence and reboot persistence. Injected
verification failure, boot timeout, and SIGTERM returned failure with successful
cleanup. Invalid accelerator, missing image, and missing host QEMU failed before
guest creation. SIGKILL recovery removed only the owned processes and runtime;
two destroy attempts succeeded while an unrelated process survived.

Affected semantic lint, repository quality, both affected OpenSpec workspace
checks, dependency-skill staging and offline Promptfoo configuration, and exact
skill discovery passed. Offline Promptfoo validation does not establish model
behavior. Module generation removed the QEMU include without lockfile churn.

The task-local validation receipt and copied result artifacts are under
`out/molecule-runner-split/`; these are disposable evidence, not source inputs.
Subsequent edits only archive these artifacts and synchronize the stable
specification. Final lint, quality, and specification checks must cover that
published candidate; runtime evidence remains applicable while runner, test,
configuration, dependencies, and host executable inputs are unchanged.
