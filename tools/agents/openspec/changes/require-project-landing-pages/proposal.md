## Why

The bridge extraction omitted its landing page because project-layout explicitly
allowed repository-only projects to omit one. New projects need landing pages
so project onboarding includes their visitor-facing entry point.

## What Changes

- Require every new first-party project under projects/ to include a landing
  page, including projects whose consumers are repository-local.
- Route implementation and validation to repo-hugo rather than duplicating
  the site procedure in project-layout.
- Add representative evaluation coverage for new repository-only projects;
  retain the distinction between generated build outputs and new projects.

## Capabilities

Procedure update only, with no runtime capability delta; skip_specs is true.

## Impact

The canonical project-layout skill and its existing offline evaluation suite.
The bridge project owns its landing content and registry entry in its existing
extraction change. No site deployment or new dependency is authorized.
