## 1. Select and declare cache policy

- [x] 1.1 Select the host cache policy from the fresh-output-base experiment and document its trade-offs.
- [x] 1.2 Update Ansible-owned settings and any required directory or permissions; build the packaged desired state and inspect its contents.

## 2. Reconcile this host

- [x] 2.1 Apply the authorized bounded host cache patch after checking concurrent state; verify desired/live bytes, ownership, and mode.
- [x] 2.2 Confirm the live setting is accepted and Go module reuse works from another fresh worktree.
- [x] 2.3 Validate OpenSpec, repository quality, and affected semantic lint for the final source candidate.
- [x] 2.4 Document workspace cleaning and shared-cache maintenance, checking the pinned Bazel help for cleanup scope and garbage-collection defaults.
