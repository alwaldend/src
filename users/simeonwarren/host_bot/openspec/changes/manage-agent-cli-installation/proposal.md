## Why

Host Bot manages Codex configuration but does not provision its executable,
and Grok CLI is absent. Rebuilding the host should install both requested
CLIs through the checked-in Ansible deployment.

## What Changes

- Install Codex and Grok for the managed user through their official shell
  installers, using the commands supplied by the user.
- Skip each installer when its command already exists and expose separate
  installation tags for a scoped deployment.
- Document the installation behavior and verify both commands on Host Bot.

## Capabilities

### New Capabilities

- `agent-cli-installation`: Provision the Codex and Grok executables for the
  managed Host Bot user without reinstalling existing commands.

### Modified Capabilities

None.

## Impact

Changes the Host Bot Ansible playbook and documentation. Grok is an explicitly
requested new host dependency. Both installer URLs intentionally resolve the
publisher's current release on first installation, as requested by the user;
this does not add a Bazel dependency or change Codex provider configuration.
The user authorized installation on Host Bot after recording the automation.
