## Why

Agent output settings currently depend on the AL Bazel wrapper selecting a
repository profile. Host Bot already manages a user Bazel configuration, which
can apply these preferences to ordinary Bazel commands in any workspace.

## What Changes

- Declare bounded, non-interactive output and failure-focused test reporting
  directly in the managed user's Bazel configuration.
- Document ordinary Bazel invocation and the portable output-settings block
  separately from Host Bot's machine-specific cache and resource settings.
- Preserve the existing Bazelisk and AL installation tasks; release stamping
  uses installed AL through the repository's release configuration.

## Capabilities

### New Capabilities

- `bazel-agent-output`: Apply agent output preferences through the managed
  user configuration without requiring a workspace-specific named profile.

### Modified Capabilities

None.

## Impact

Changes the checked-in Host Bot Bazel configuration and documentation. The
companion AL `unify-release-versioning` change removes the command wrappers
and updates repository setup instructions. No new dependency or live host
deployment is included.
