## Why

The Anki CLI initially upgraded an existing SQLite driver to an old release
that supplied the required API. The user requested a default preference for
the latest release when upgrading dependencies.

## What Changes

The canonical dependency skill prefers the latest stable upstream release,
verifies freshness and compatibility, and documents exceptions. Its evaluation
cases cover choosing a current release rather than the oldest API-compatible
release.

## Capabilities

No behavioral specification delta. This is a procedure update.

## Impact

tools/agents/skills/repo-external-dependency. Existing reproducibility and
dependency approval requirements remain in force.
