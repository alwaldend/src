## 1. Configure ordinary Bazel invocation

- [x] 1.1 Add the portable agent output defaults to the managed user bazelrc; inspect the diff to confirm cache, SDK, and resource settings are preserved.
- [x] 1.2 Update the host documentation and repository setup guidance; verify the instructions use ordinary Bazel commands and link to the canonical output-settings block.

## 2. Verify the checked-in configuration

- [x] 2.1 Validate the OpenSpec change and packaged Ansible input, then inspect a Bazel invocation that loads the candidate configuration without changing the live host.
- [x] 2.2 Verify the rendered setup documentation and record the source validation evidence.

Validation evidence: strict OpenSpec validation passed; the Ansible package and
site built with the candidate rc loaded explicitly. The rc announcement reports
the output defaults, and the rendered setup page contains the installation,
release-build, and user-configuration guidance. Logs and output hashes are under
`out/unify-release-versioning/api-naming/plain-bazel-*`. No host was deployed.
