## Why

The stored LeetCode submissions belong to simeonwarren. Move them from the
shared data tree to `users/simeonwarren/leetcode`, retaining site inclusion.

## What Changes

- Move the submission tree with original submission bytes preserved.
- Update downloader output directories and Bazel labels.
- Keep generated submissions on the main site, as explicitly requested.
- Fix the delivery tool preparation receipt capacity so this move can be
  published through its guarded workflow.

## Capabilities

No new capability requirements. This relocation and delivery capacity fix
preserve existing behavior; `skip_specs: true` records that scope.

## Impact

Submission ownership, downloader paths, and the site dependency label change.
Preparation receipts gain a bounded capacity of 4 MiB; other typed records retain
256 KiB limits. No dependencies or live deployments are added. Acceptance requires
submission preservation, site packaging, and a validated pull request.

Owner-local records:

- [Downloader paths](../../../../../projects/leetcode_downloader/openspec/changes/relocate-simeonwarren-submission-output/proposal.md)
- [Site packaging](../../../../../projects/alwaldend.com/openspec/changes/consume-relocated-leetcode-submissions/proposal.md)
- [Delivery capacity](../../../../../tools/repo_delivery/openspec/changes/support-large-preparation-receipts/proposal.md)
