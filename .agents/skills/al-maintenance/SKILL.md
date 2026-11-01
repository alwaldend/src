---
name: al-maintenance
description: >-
  Refresh repository-generated metadata and dependency declarations through their
  owning update workflows. Use for stale inventories, module declarations, skill
  discovery, or explicit dependency/source maintenance; use al-checks for validation.
---

# Maintain repository source

Use al-workspace and al-openspec to establish the task workspace and maintained
change before edits. Read the owner's README and identify its generator and
maintained inputs. Existing generated projections are outputs, not new owners.

Use `al bazel` for repository targets. Common owners are:

- `//tools/bazel_module:update` for generated root module declarations;
- `//tools/repo_quality:inventory_update` for OpenSpec owners and the project index;
- `//.agents:write_skills` for skill discovery;
- the owning Gazelle target for generated BUILD files.

Change maintained inputs, run the owner workflow, and inspect the complete diff.
Run the same generator a second time to verify stability while inputs are unchanged.
Keep logs and disposable outputs under the task's ignored out directory.

For dependency changes, use repo-external-dependency and the owner's pinned
language tooling. Reusing existing first-party modules and repository pins is
allowed; genuinely new external software requires the user's specific approval.
Generate module locks through Bazel and language manifests through their owning tool.

Use repo-bazel for focused validation and al-checks when the changed inputs and
consumers require a complete workspace matrix. Finish with al-delivery. Updating
source does not authorize installation, deployment or live infrastructure changes.
