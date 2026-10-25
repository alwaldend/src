## Context

See proposal.md for scope. The baseline inventory is 250 open alerts from the
GitHub API, captured on 2026-10-25 in ignored `out/dependabot/alerts.jsonl`.
Native uv locks, pnpm locks, and Go manifests own dependency resolution.

## Goals / Non-Goals

Preserve affected tooling while removing vulnerable selections. Do not dismiss
alerts, alter live host configuration, or replace existing package managers.

## Decisions

- Prefer current stable releases verified through PyPI, npm, and the Go proxy.
  Use targeted security updates when a broad upgrade would invalidate a maintained
  consumer patch or change an unrelated contract.
- Raise Python minimums to 3.11 where secure Pillow/filelock releases require it.
  Bazel uses Python 3.14.4; retaining Python 3.10 would retain vulnerable packages.
- Refresh Promptfoo and rebase its telemetry/concurrency patch on the exact
  published distribution. Keep Mermaid's pinned rendering patch intact unless
  an upgrade proves necessary.
- Compare every captured advisory range with the final lock selections. Keep
  genuine unpatched exposure explicit rather than suppressing audit evidence.
- This record owns the coordinated repository dependency maintenance; affected
  manifests remain authoritative within their existing component boundaries.

## Risks / Trade-offs

- Major dependency updates can change APIs: run component E2E and offline
  validation targets plus representative CLI outputs.
- Some advisories have no patched version: inspect current registry versions
  and dependency paths before deciding whether removal is compatible.
- Toolchain and network failures can obstruct validation: preserve diagnostics
  separately from candidate failures.

## Decision Review

Verdict: proceed with secure releases and targeted overrides. Keeping existing
vulnerable selections fails the requested security outcome. Updating Mermaid
to a new major would also require reworking a rendering patch; a targeted KaTeX
override removes the vulnerable version with a smaller compatibility surface.
Verify a representative diagram before accepting that override.

## Execution State

Current candidate: dependency updates and Cordis YAML migration implemented.
The final version-range reconciliation resolves all 250 captured alerts; the
unpatched ecdsa implementation is excluded from unsigned DFU packaging.
The same-input DFU package comparison is identical with and without ecdsa.
Next action: rerun exact-candidate delivery checks and publish.

## Cordis YAML compatibility

The user explicitly requested retaining js-yaml 5 and fixing the integration.
The Include plugin owns the entry-list dialect; its v4 schema cannot be passed
to a v5 loader. Migrate its existing `!!js` tag through a focused package patch
and resolve the plugin and runtime to the same js-yaml release. Do not copy the
dialect into the runtime. Failure cases to verify before accepting: ordinary
entry-list parsing, dump/load preservation of `!!js` expression nodes, lazy
expression activation, persisted configuration reload, and malformed input
rejection. Existing runtime and stdio E2E tests cover ordinary configuration and malformed
input. A new runtime E2E covers expression execution, preservation during a
configuration write, scalar types, and restart; its repeatable result is
`yaml_roundtrip.json` in the Bazel test undeclared outputs.

## Session Review

- `DEPENDENCY-LOCK-RETENTION`: the first uv and pnpm updates retained some
  vulnerable transitive or peer versions. The baseline-to-candidate range
  reconciliation exposed them; manifest constraints and explicit peer ownership
  fixed the selections without hand-editing locks. Reuse the captured baseline
  and compare actual resolved versions on continuation.
- `YAML-SCHEMA-VERSION-MISMATCH`: the initial dependency upgrade mixed v4 and v5
  schema objects. Migrating the owning plugin and adding a behavior E2E resolved
  it; the final runtime test passed all ten cases.
- `NESTED-TMPDIR-PATH`: an initial nested-workspace test used a nonexistent
  relative temporary directory. An absolute task-local TMPDIR corrected it;
  all six Promptfoo tests then passed.
- `ANKI-CRATE-MISSING-MANIFEST`: the Anki-backed transport check could not build
  because the resolved percent-encoding-iri archive lacks Cargo.toml. The
  unchanged Rust dependency is outside this fix; the pure Go client builds and
  Al plugin transport tests pass with the new gRPC version.

The first discovery included an unnecessarily large worktree inventory and tool
catalog. Intent-to-add index entries used for formatting were normalized to
fully staged task-owned files after the delivery index-flag gate rejected them. Subsequent discovery was bounded to relevant metadata. No shared host
configuration, credentials, live infrastructure, or alert state was changed.

## Final Advisory Decision

New evidence reopens the initial decision to retain ecdsa. The pinned firmware
CMake commands invoke `adafruit-nrfutil dfu genpkg` without key-file arguments;
its signing module already handles a missing ecdsa import. Verdict: remove the
unused signing dependency through uv exclusions. This removes the vulnerable
implementation rather than asserting an upstream fix. Alternatives are to keep
the advisory or implement a new signing backend; neither is needed for the
repository's unsigned DFU outputs. Before accepting, rebuild the firmware and
compare DFU packages generated from identical firmware input with and without
the previously pinned ecdsa dependency. Signing through this helper is
intentionally unavailable. The first whole-firmware hash comparison differed:
SystemInfo.cpp embeds `__TIME__`, so separately compiled images are not suitable
controls. The same-input comparison verifies identical BIN, DAT, and manifest
member hashes, and the firmware build test passes without ecdsa.
