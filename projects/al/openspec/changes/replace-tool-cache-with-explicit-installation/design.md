## Context

The only root AL tool declaration is `repo_delivery`. `bazel_agent` already
runs ordinary Bazel commands. The user explicitly replaced the initially
requested download approach with local builds and explicit installation.

## Goals / Non-Goals

Remove tool caching and provide an install target for standalone host execution.
Do not upload releases, add selection flags, install on the real host, or alter
shared Codex configuration.

## Decisions

- Remove the AL cached-tool interface and reserve its removed protobuf fields.
- Keep `//tools/repo_delivery` as the ordinary source-built default target.
- Add a reusable Go installer under the existing file_installer owner. It copies
  a standalone executable to a temporary sibling and atomically replaces the
  destination with mode 0755. Default destination is `~/.local/bin`; an explicit
  directory supports isolated validation.
- `//tools/repo_delivery:install` and `//projects/bazel_agent:install` wrap the
  same installer with their local binaries and command names. Build repo_delivery without cgo so it needs no runtime bundle.
- Update guidance to use the installed command or ordinary Bazel invocation.

## Risks / Trade-offs

Installed tools update only when the install target runs. No automatic source
refresh is promised. Installation supports standalone binaries without runfiles.

## Validation

Before implementing the installer, exercise installation and replacement in an
isolated directory, execute the copied binary, and reject missing sources and
invalid names while preserving a prior installation. Produce a JSON test receipt.
Run AL, runner and delivery tests, configured formatting and semantic lint.

## Acceptance evidence

Both tool-specific install targets completed in `out/tool-installation/bin`.
`acceptance.json` records mode 0755, executable hashes, direct delivery provider
execution, and the installed runner's exact argument forwarding. The installer
E2E test emits `installation.json` for repeated installation, executable startup,
and preservation after failed updates. No host installation or upload occurred.

Session observation `FORMAT-OUTPUT-BOUNDING`: whole-repository formatters emit
large file inventories. Subsequent formatter output was redirected to ignored
task logs rather than copied into session context. No shared workflow change
is proposed for this bounded operational adjustment.
