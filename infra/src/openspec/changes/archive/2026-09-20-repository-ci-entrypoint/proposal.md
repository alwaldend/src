## Why

The secure workflow currently verifies runner infrastructure instead of building and testing the repository. Review requests one shared CI entry point and removal of the completed runner acceptance harness.

## What Changes

- Add `//tools/ci` to build and test every normal target in the root and declared nested Bazel workspaces, and invoke it through a local Node.js action from the secure workflow.
- Remove the runner smoke package, publication helper, playbook, and validation-branch configuration; retain historical acceptance evidence.
- Remove the redundant Node CA override after a read-only check on the deployed runner.
- Retain declarative Android prerequisites because Bazel resolves the NDK before it can run the installer.

## Capabilities

### New Capabilities

### Modified Capabilities

- `repository`: CI workflows use one shared build/test entry point.

## Impact

Workflow, tools, runner documentation and packaging, CI skill, and Vault's declared allowed refs. Removing the validation ref narrows desired Vault policy; this source change does not claim that policy has been deployed. The runner-owned specification records its corresponding contract update.
