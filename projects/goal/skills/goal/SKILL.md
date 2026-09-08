---
name: goal
description: >-
  Disabled legacy goal workflow, retained only as a compatibility notice.
  Use the openspec skill for new specifications and resumable work.
---

# Goal (disabled)

This skill is disabled and removed from repository skill discovery. Do not
use it to initialize or pursue new goals. Use the `openspec` skill and the
affected owner's `<owner>/openspec/` workspace for specifications, change
tasks, and continuation. `infra/src/openspec/` is for repository evolution.
Select that owner when running the pinned command from the root Bazel
workspace: `OPENSPEC_PROJECT=<owner> bazel_agent bazel run //tools/openspec -- ...`.

The deprecated goal CLI retains its existing commands for legacy record
compatibility, inspection, validation, and interrupted-publication recovery.
Its warning goes to stderr; existing stdout formats remain unchanged.

Historical references below describe the old format. They do not establish
the current workflow or authorize additional work:

- [Record format](references/record-format-v1alpha1.md)
- [Lifecycle and evidence](references/lifecycle-and-evidence.md)
- [Sessions and concurrency](references/sessions-and-concurrency.md)
- [Graph organization](references/graph-organization.md)
- [Promotion and legacy import](references/promotion-and-migration.md)
