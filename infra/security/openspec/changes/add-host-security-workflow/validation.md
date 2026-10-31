# Validation evidence

The modern-host playbook, pinned scanner archive, collection roles, and reused
inventories passed packaging and Ansible syntax checks. A local download host
completed the audit/update canary with its boot identity unchanged. The later
fleet execution completed Lynis and package maintenance on ten modern hosts.
Private pre/post evidence records unchanged boots, native automatic updates,
installed versus running kernels, and service-state comparisons. No previously
active service became inactive in the post-update sample. Appliance results
are tracked in the sibling support-appliance-security change.

Two successful Ansible runs ended with a wrapper token-revocation timeout.
Their successful task recaps and post-update evidence establish package
execution separately from credential cleanup. This cleanup failure remains
an operational limitation; it was not treated as an update failure or ignored.

## Session ergonomics

- `security-mutable-runfiles`: an audit process lost its Bazel runfiles while
  other build work was occurring. The cause was not established. Task-local
  materialized runtime snapshots allowed affected audits to finish; the
  snapshots preserved the Python interpreter link required by the packaged
  virtual environment. This is an operational workaround, not a source
  deployment contract or a shared tooling change.
- `security-ssh-shell-output`: raw SSH boot-ID collection included shell
  integration escape sequences. Native command-module collection removed
  the ambiguity; unchanged boot identities were independently rechecked.
- `security-base-selection`: an initial assumption about stacking on PR 202
  was corrected to the user's explicit separate master-based PR. Task source
  was preserved during the base transition; PR 202 was not changed.
- `security-bounded-evidence`: oversized combined reads caused avoidable
  truncation. Subsequent reads extracted per-host diagnostic fields and kept
  full operational artifacts in private ignored scratch.

Exact candidate source gates and publication are recorded by the local
delivery receipts; this document does not substitute for those gates.

Appliance completion: XOA native OS updates and XCP-ng vendor packages were
installed with unchanged boot identities. XCP-ng required its documented
origin mirror and a readiness retry before restoring host placement. The
three other managed guests retained their boots and active service states.
Twelve configured candidates could not be inspected; none was classified as
inactive merely because access failed. Private evidence records each limit.

## PR simplification validation

On 2026-10-29, the revised maintenance playbook passed packaged Ansible syntax
and 13 offline behavioral cases with pinned Ansible 2.20.7/Python 3.14.4.
Fixture roles and a transport that rejects connections exercised full success,
audit/update failures, changed boot identity, role tags, and unreachable
outcomes. Ordinary failures remained failures after successful `always` checks;
audit failure skipped the updater, and changed boots failed even after a role
failure. Initial unreachable access was explicitly reported and excluded;
later transport failures remained unreachable without claiming boot validation.

The repeatable `boot-e2e/verify_boot.py` and its `receipt.json`, plus
`candidate-inspection.json` with matching packaged source hashes, are under
`out/host-security-updates-20261029/clamav/pr-simplification/`. These receipts
validate orchestration with fixtures, not live role or host behavior.

The SSH target now reuses the generated maintenance config through the existing
`al_binary_run` and a quoted shell `exec`, with the parent's explicit approval
to change this target's rule/load. Its AL config, injector dependency,
maintenance selector, and host-PATH SSH selection are unchanged. There was no
packaged SSH binary in the previous `tools/ssh` macro. The configuration path
and user arguments remain quoted, and a missing or empty config variable fails
before invoking SSH.

`security-ssh-config-expansion`: AL passes argv literally, OpenSSH does not
expand environment variables in `-F`, and `Include` is forbidden on its command
line. Those two offline checks failed before connection. Bazel also tokenizes
`args`, so the shell command must have an outer quote pair in the BUILD value.
`ssh-command.sh`, generated with Bazel's non-executing `--script_path`, retains
the shell command as one argument. The original direct-shell fixture receipt
`ssh-config-verification.json` verifies trust, identity, BatchMode, timeouts,
and a path containing spaces using only `ssh -G`.

The final `verify_ssh_wrapper.py` harness and
`ssh-wrapper-e2e/20261029T204619755885Z/receipt.json` verify exact source and
generated-command hashes, unchanged fixed AL arguments, native `ssh -G` policy
and `-oHostKeyAlias`, and exact forwarding through a PATH-selected fixture SSH
for spaces, empty arguments, and shell metacharacters. Both missing and empty
config variables fail before SSH invocation. The SSH package build, scoped
Buildifier check, and strict OpenSpec validation passed. No AL plugin or live
SSH operation was executed during these checks. All three delegated fixes are
ready for parent diff review and delivery.

Session friction: `security-bounded-evidence` recurred in two oversized initial
reads. Subsequent commands used bounded files and selected diagnostics.
`security-fixture-plugin-loader` required initializing the native Ansible
plugin loader before rendering a file lookup; no source defect was involved.
The parent's brief Bazel pause was honored while its runtime was captured;
no shared-cache repair or clean was attempted. No live operations, commits,
or publication were performed by this subtask.

## Complete Ansible follow-up review

Reviewed the maintenance, updater, SSH, Lynis, ClamAV, appliance and Host Bot
changes and their owner callers. Corrected kernel-device traversal, XOA failure
cleanup and missing SSH-directory provisioning; native async results replace
the asynchronous scanner marker. Fleet inventory visibility now uses one named
package group while retaining its exact consumer and XCP-ng's separate
appliance access. Consumer packaging passed.

All 41 changed non-OpenSpec YAML files parsed. Pinned syntax checks passed for
19 owner playbooks across all 13 inventory owners. Runtime fixtures and their
limits are summarized in the private
`out/host-security-updates-20261029/ansible-review/validation-summary.json`;
individual owner records retain the relevant acceptance. Existing package,
SSH and scanner policies were not deployed by this source review. Full fleet
scan completion remains separate from these implementation checks.
