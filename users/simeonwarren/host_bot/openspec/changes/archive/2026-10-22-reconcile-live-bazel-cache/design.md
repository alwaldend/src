## Context

Read-only inspection found the missing live cache setting and directory. JVM thread evidence showed the fresh worktree unpacking the pinned LLVM archive. The checked-in host configuration already declared the intended setting and directory.

## Decision and authority

The user explicitly requested: "just modify the host manually and ensure its in sync with the role - no need to run ansible", and requested inclusion in the new PR. This authorizes the exact imperative exception for the bazelrc and repository contents cache directory; it does not authorize unrelated host changes.

Create the cache directory using the existing Ansible ownership and mode. Install the checked-in bazelrc at its existing managed destination with the role's ownership and mode. Do not invoke Ansible or repeat storage provisioning. Checked-in Ansible remains the sole owner of desired configuration.

## Verification

Compare the live file byte-for-byte with its checked-in owner, verify both paths' permissions and ownership, and confirm Bazel accepts the configuration. Keep detailed execution evidence under ignored task scratch. Verify materialized repository reuse across isolated task-owned Bazel fixtures.

## Scope

This record carries the manual-operation authorization and IaC reconciliation required by repository policy. Source configuration needed no amendment because it was already correct.
