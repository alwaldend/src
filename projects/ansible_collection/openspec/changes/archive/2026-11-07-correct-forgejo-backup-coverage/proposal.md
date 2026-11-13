## Why

Review of PR 273 found three gaps in the existing complete-backup guarantee:
quoted INI paths can be misresolved, legacy storage can be omitted, and custom
backup roots can exclude installed startup files.

## What Changes

- Normalize simple quoted INI values and reject unsupported syntax whose path
  interpretation cannot be established safely.
- Include effective legacy LFS and avatar paths, respecting modern replacement
  settings, in coverage and relocation validation.
- Require the installed executable, configuration and unit in custom backup
  coverage before downtime.
- Extend the real-binary integration rehearsal with external persistence,
  rejected incomplete coverage and complete restoration of those files.

## Capabilities

No new or modified requirements. This corrects implementation of the existing
`forgejo-controlled-upgrades` complete recovery and preflight requirements.

## Impact

The collection's Forgejo planner, preflight tasks, documentation and existing
integration fixture. No new dependency, deployment or live configuration change.
