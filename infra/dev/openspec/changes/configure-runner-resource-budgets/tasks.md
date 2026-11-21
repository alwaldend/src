## 1. Runner configuration

- [x] 1.1 Set host-relative job, CPU and memory budgets in the runner bazelrc; verify the pinned Bazel parses the exact file for build and test commands.
- [x] 1.2 Document the runner budgets and their scheduling-only scope; inspect the diff to confirm workstation defaults and runner capacity remain unchanged.

## 2. Verification

- [x] 2.1 Build the Ansible package, run its offline syntax check, validate the OpenSpec change and pass repository quality and affected semantic lint checks; record the results.
