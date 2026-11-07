## Why

The user requested simplification of PR 273. Configuration inference expanded a
known local SQLite upgrade into a separate Forgejo configuration interpreter.

## What Changes

- Remove the custom INI planner and its filesystem-inference guarantees.
- Make `forgejo_backup_paths` the explicit archive contract and expose the
  configured SQLite file through `forgejo_database_path`.
- Retain native version checks, offline archives, protected root metadata,
  unchanged-file repair, health checks and safe recovery boundaries.
- Resolve host-side recovery/root aliases before checking storage separation.
- Keep the real-binary upgrade/restore and failure rehearsals; remove tests of
  the deleted interpreter.

## Capabilities

Modify `forgejo-controlled-upgrades` to describe declared local persistence.
The caller owns external roots and separately coordinated path migrations;
SQLite/HTTP and encryption-key binding guards remain native Ansible checks.

## Impact

The collection role, its packaging, specification and existing integration
fixture. No deployment, new dependency or live configuration change.
