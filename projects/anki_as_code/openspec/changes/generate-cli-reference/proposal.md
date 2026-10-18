## Why

The CLI has many nested commands, including descriptor-generated Anki RPCs.
A short reference index with a child page per root subcommand keeps their help
easy to navigate without maintaining or committing a second command inventory.

## What Changes

- Add a project-local tool under `tools/generate_cli_docs` that renders the
  existing Cobra tree into a short index and one page per root subcommand.
- Include each root subcommand's complete nested help on its own page without
  executing application commands.
- Generate the directory through Bazel and package the hierarchy with project docs.
- Keep the existing reference URL and generated Markdown out of Git, then
  redeploy the public site through its existing release workflow.

## Capabilities

### New Capabilities

None. This is documentation tooling; command behavior is unchanged.

### Modified Capabilities

None. This change declares `skip_specs: true`.

## Impact

Project-local generator, CLI target visibility, documentation build wiring,
and README. Reuses Cobra and standard-library packages without new dependencies.
