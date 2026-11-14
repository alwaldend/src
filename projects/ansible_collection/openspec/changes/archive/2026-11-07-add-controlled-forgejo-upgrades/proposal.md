## Why

The Forgejo role replaces its executable and reloads the service without a
consistent backup or upgrade acceptance checks. A database migration can make
restoring only the prior executable unsafe.

## What Changes

- Detect changed executable, configuration, and service inputs before activation;
  reject downgrades and unsupported backup layouts before stopping the service.
- Drain an active instance with its installed executable, stop it, and retain
  a protected offline archive before replacing any installed startup input.
- Start the candidate and verify HTTP health, running version, and doctor checks.
  Preserve recovery material and stop a rejected candidate without automatically
  starting an old executable against a possibly migrated database.
- Preserve fresh installation and unchanged-run behavior, and document restore
  steps and the supported local SQLite storage boundary.
- Exercise packaged Ansible tasks with local upgrade and failure fixtures.

## Capabilities

### New Capabilities

- `forgejo-controlled-upgrades`: offline upgrade backups and guarded activation.

### Modified Capabilities

None.

## Impact

`roles/forgejo`, collection filter support, and collection-owned integration
tests. Existing dependency pins remain unchanged except a test-only reference
to the previous Forgejo release for migration acceptance. No deployment,
production credentials, or persistent host configuration changes are authorized.
