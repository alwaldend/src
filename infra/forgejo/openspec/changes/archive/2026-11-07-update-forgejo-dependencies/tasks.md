## 1. Update existing pins

- [x] 1.1 Update server/source, runner, and CLI pins; verify publisher integrity,
      source extraction, and version output from the Bazel-built binaries.
- [x] 1.2 Update the provider archive and consuming constraint together; reconcile
      the Bzlmod lock and verify local provider packaging and Terraform checks.

## 2. Validate and deliver

- [x] 2.1 Build affected Ansible packages, validate the OpenSpec record, and pass
      formatting, repository quality, and affected-target semantic lint checks.
- [x] 2.2 Record verified output and validation evidence; verify the final diff
      contains only task-owned dependency and OpenSpec changes. Git and delivery
      receipts own subsequent commit, push, and pull-request state.
