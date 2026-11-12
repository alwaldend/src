## Why

Project lifecycle labels have no shared definitions, overlap, and drift between
READMEs and landing pages. Several READMEs omit them, and unfinished or historical
feature work has been mistaken for current active development.

## What Changes

- Define one lifecycle vocabulary: active, maintenance, completed, archived.
- Review all direct project owners and assign exactly one README status using
  the user-confirmed current feature-development focus.
- Generate landing status and the central status index from README metadata;
  reject invalid statuses and stale projections through inventory checks.
- Preserve release, deployment, and acceptance limitations independently.

## Capabilities

### New Capabilities

- `project-lifecycle`: Canonical definitions, project declarations, generated
  status projections, and current-project review evidence.

### Modified Capabilities

None.

## Impact

Shared project metadata, the repository inventory command and its E2E fixture,
project READMEs and landing front matter. Existing dependencies suffice. Old
taxonomy terms in_progress, finished, and abandoned are replaced; retain their
public URLs as redirects. No deployment or host configuration changes.
