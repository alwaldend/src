## Context

The dev1 playbook manages the bot account's global Git configuration.

## Decisions

Declare the delivery key alongside the existing configuration. Use the
operator identity from public Git history. The Ansible deployment remains the
owner of the installed configuration.
