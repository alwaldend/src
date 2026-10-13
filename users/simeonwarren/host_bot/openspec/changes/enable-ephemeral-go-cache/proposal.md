## Why

The host already shares Bazel action and repository-contents caches, while
Gazelle's Go module download cache is scoped to each output base. The user
selected ephemeral module caching to avoid retaining another copy of module
downloads when completed repositories can be reused.

## What Changes

- Set `GO_REPOSITORY_EPHEMERAL_MODCACHE=1` in the managed user bazelrc.
- Apply that exact setting to this host's live bazelrc, as requested.
- Document the distinction between completed-repository reuse and the
  temporary Go download cache, including refetches on cache misses.
- Explain workspace cleaning, shared-cache retention, and the pinned Bazel
  garbage-collection defaults.

## Capabilities

No specification delta: this changes a developer-tool cache policy.

## Impact

The owning Ansible bazelrc, host README, and live user bazelrc are affected.
The companion repository change upgrades Gazelle to 0.53.0. Existing shared
cache directories and Codex permission configuration need no expansion.
