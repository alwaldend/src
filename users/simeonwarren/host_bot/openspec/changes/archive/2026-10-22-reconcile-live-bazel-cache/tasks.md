## 1. Reconcile live configuration

- [x] 1.1 Identify live drift against the existing Ansible-owned bazelrc and cache-directory declaration.
- [x] 1.2 Apply the user's exact manual-operation exception to those two paths without running Ansible.
- [x] 1.3 Verify byte equality, ownership, permissions, and acceptance of the restored Bazel configuration.
- [x] 1.4 Verify repository contents reuse across fresh task-owned output bases and record the evidence.

## 2. Deliver

- [x] 2.1 Prepare this bounded repair record for PR #136 and select combined exact-candidate checks; delivery receipts own validation and publication results.
