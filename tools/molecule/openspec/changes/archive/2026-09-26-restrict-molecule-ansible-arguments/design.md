## Context

[PR #112 review](https://github.com/alwaldend/src/pull/112#discussion_r4111472002) identified unrestricted Ansible arguments.
For example, an extra variable can change `ansible_connection` to `local`.
Current consumer declarations use only `--skip-tags=ssh_sign_key` and `--diff`.

## Goals / Non-Goals

Prevent command-line arguments from replacing runner-owned connection and
inventory settings while preserving current consumer behavior. Scenario source
remains declared test code; this change does not sandbox arbitrary playbooks.

## Decisions

Use an allowlist instead of enumerating dangerous Ansible flags. Accept literal
`--diff` and a single `--skip-tags=` argument with a comma-separated list of
ordinary tag identifiers. All other arguments fail at the beginning of preflight,
before QEMU or Ansible starts. The error does not echo rejected values because
an argument may itself contain sensitive data. The Go preflight owns validation;
the Bazel rule forwards the declaration without a second validation policy.

Before implementation, add real-runner prerequisite cases for short extra vars,
extra-variable files, connection, inventory, user, and private-key overrides.
Do not execute those cases against the vulnerable implementation: they could
redirect privileged tasks to the controller. Run them only after the guard is
present. The public smoke target uses the exact supported host-role options.

## Risks / Trade-offs

Previously accepted arbitrary flags now fail. No current consumer loses a
required option. Future argument support needs an explicit review of its
connection and inventory effects. Scenario variables stay in declared playbooks.

The negative cases stage an assertion-only probe with fact gathering and
privilege escalation disabled. If preflight regresses, these cases cannot
execute the deployment smoke tasks on the controller. This is declared test
input, not a mocked runner or VM lifecycle.

The first nine prerequisite cases passed. The accompanying allowed-options
smoke attempt failed in create after the guest booted but SSH did not become
ready. Serial observations showed a login prompt and a service-start failure;
the exact unit and cause are unverified. Cleanup succeeded. Evidence is retained
under `arguments-startup-failure` with stable defect ID
`MOLECULE-GUEST-SSH-STARTUP`; a fresh guest is used for a bounded retry.

The retry reproduced failed cloud-init and sshd startup and Btrfs checksum
errors. Retained input identity revealed that the cached image was
`a57d27af98d6017dcbd5505dbc3567d266d4f6cd802bf9f91ece8d0a077cb864`,
whereas the declared Fedora pin and previous successful tests used
`28680fe5b371a5a82ebf43a31926e086a168e59949d03969c5093e7071f90b7f`.
Host QEMU hashes were unchanged. This is a recurrence of the image-cache
mismatch recorded in the original runner change; the cause of the cached-byte
change is unverified. Restore through the owning Bazel fetch workflow and
verify the pinned bytes before counting new guest acceptance.

## Acceptance evidence

The owning `bazel_agent bazel fetch --force --repo=@org_fedora_cloud` restored
the declared SHA-256. On the restored bytes, the nine prerequisite cases and
public smoke scenario passed. Each rejection recorded zero guest PID, no
lifecycle phases, and successful cleanup. Smoke used the existing consumer's
`--skip-tags=ssh_sign_key` and `--diff` options, completed all six phases, and
retained the pinned image hash. The dynamic-name sentinel remained absent from
all artifacts. The negative fixture is assertion-only, so a future guard
regression cannot run deployment tasks locally.

Fresh test artifacts, source hashes, and runtime hashes are retained under
ignored `out/molecule-runner-split/arguments-green-evidence` and
`arguments-test-inputs-2.json`. The earlier two boot failures are preserved as
invalid-image observations, not successful acceptance. Final candidate gates
must confirm source/runtime input stability after archival and formatting.

## Review follow-up: special tag selectors

A subsequent review of `3af6e822b9e8` identified that syntactically valid special
skip selectors could suppress the untagged VM lifecycle and produce a successful
empty run. The pinned Ansible implementation and its
[special-tag documentation](https://docs.ansible.com/projects/ansible/latest/playbook_guide/playbooks_tags.html#special-tags)
identify `all`, `tagged`, `untagged`, `always`, and `never` as reserved selectors.
Preflight now rejects all five, including inside comma-separated lists. Ordinary
role tags remain supported. This completes the same argument-boundary change;
the preceding nine-case results describe the earlier candidate, not this fix.

Seven additional real-runner rejection cases were written before the guard
change. On 2026-09-26, all 16 prerequisite cases passed with no guest or lifecycle
execution and successful cleanup. The public smoke scenario retained the
ordinary host-role options and completed all six phases on the verified Fedora
pin. Runtime/source hashes and artifacts are under ignored
`out/molecule-runner-split/tags-green-evidence` and `tags-test-inputs.json`.
