## Context

Forgejo 17.0.1 paginates file diffs using `[ui] DIFF_PAGING_NUM`. Its diff
content is a flex item whose intrinsic width can push the file header and
actions outside the mobile viewport.

## Decisions

- Set the diff page size in the instance configuration and the branch cleanup
  default on the existing Terraform repository resources. The pinned provider
  sends the branch setting to Forgejo but does not read it back during refresh.
- Use the Forgejo custom header hook to apply a stylesheet only to Files Changed
  below 768 pixels. Let the content column shrink, wrap filenames, tabs, and
  controls, and keep wide code content scrollable inside its file panel.
- Use the Ansible role's optional startup input for the header, so its first
  installation and changes follow the same activation as the INI file.
