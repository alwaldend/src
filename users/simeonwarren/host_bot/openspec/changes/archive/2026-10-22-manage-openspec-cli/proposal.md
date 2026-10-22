## Why

OpenSpec was installed manually with `npm install -g @fission-ai/openspec@latest`.
Host Bot's checked-in Ansible must reproduce that tool installation.

## What Changes

- Install the scoped OpenSpec npm package at a pinned version under the managed user's local prefix.
- Reuse the existing Node.js/npm package provisioning and document the managed CLI.

## Capabilities

### New Capabilities

- `openspec-cli`: Reproducible user-local OpenSpec installation on Host Bot.

### Modified Capabilities

None.

## Impact

Host Bot's Ansible playbook, inventory variables, and README. The requested
package is already used by the repository; no repository dependency pin changes
or live deployment are included.
