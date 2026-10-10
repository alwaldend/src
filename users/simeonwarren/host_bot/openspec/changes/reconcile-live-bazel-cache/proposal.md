## Why

The live host-bot bazelrc omitted the already-declared repository contents cache, and its cache directory was absent. A new worktree therefore repeated pinned LLVM archive extraction.

## What Changes

Record the user-approved manual repair of the two live paths to the existing checked-in Ansible desired state. No configuration setting, dependency, storage allocation, or reusable role behavior changes.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None; this is a deployment reconciliation record without a specification delta.

## Impact

The live bazelrc matches `ansible/files/bazelrc`; the repository contents directory matches the owner, group, and mode declared in `ansible/group_vars/all.yaml`. Include this repair record in PR #136 with the extracted Anki models, as explicitly requested.
