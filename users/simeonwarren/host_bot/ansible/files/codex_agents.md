# Host Codex instructions

Test Codex model, provider, authentication, and shared-configuration
migrations in an isolated opt-in profile before changing a shared default.
Breaking global configuration is especially harmful on a host with concurrent
users or agents.

Checked-in Ansible is the desired state for this host. Mirror every
intentional host configuration change into
`users/simeonwarren/host_bot` before handoff.
