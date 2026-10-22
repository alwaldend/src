## Context

See proposal.md for motivation. `internal/cli.New` constructs the full Cobra
tree around an injected application interface. `docs_filegroup` packages
generated files, and the existing schema generator uses `run_binary`.

## Goals / Non-Goals

Generate a deterministic reference hierarchy from the actual command tree without
opening collections or linking the native server. Native Rust flags remain
owned by native help rather than being copied into the Go command definitions.

## Decisions

- Render each visible command's Cobra help inside a fenced text block beneath
  its command-path heading. This preserves aliases, arguments, inherited flags,
  defaults, and protobuf request help using existing dependencies. A separately
  maintained command catalog would drift; Cobra's doc subpackage would add
  otherwise unused Markdown/manpage dependencies.
- Traverse commands in Cobra's sorted order, initialize the normal help
  command/flags, and never execute the CLI tree. Mark argument-forwarding
  commands with instructions for obtaining their native help.
- Place the generator in `tools/generate_cli_docs`, the build output in
  `docs/cli/`, and retain `/docs/projects/anki_as_code/docs/cli/` as a short
  section index. Generate `_index.md` with the CLI description and child links,
  and `<command>.md` for each visible root command with all its nested help.
- Declare a single `run_binary` output directory so the build does not maintain
  a second root-command inventory. `docs_filegroup` packages that directory
  under its `cli` basename. Generated Markdown exists only in Bazel output.

## Risks / Trade-offs

- The API tree remains long on its own page; the short index and sibling pages
  keep other workflows independently navigable.
- Native forwarding has no Cobra flag metadata; describe this boundary and
  point readers to native `--help`.
- Failure cases include an oversized index, missing child links/pages, help
  leaking between root groups, lost nested commands/inherited flags/aliases,
  missing request fields, nondeterministic output, missing output arguments,
  and a file blocking the output directory. Update the executable E2E before
  implementation, then inspect the packaged site hierarchy and public deployment.
