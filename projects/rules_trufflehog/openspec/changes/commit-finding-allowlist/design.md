## Context

See `proposal.md`. The current path filter hides whole files. A scan-only
inventory reports ten previously excluded findings. The user requires commit
and line matching and explicitly selected `projects/rules_trufflehog` as owner.

## Goals / Non-Goals

Provide a standalone module with a module-owned scanner and exact suppression.
Consumers supply Git, repository location, and policy, with an optional scanner override.
This change does not alter Git history, add value hashes, or verify credentials.

## Decisions

- Use the existing Cobra, rules_go runfiles, and standard library; add no new
  external software. Reuse TruffleHog's generated `source_metadatapb.Git` type
  at the existing scanner pin, 3.99.7, instead of maintaining its field schema.
  The CLI uses `encoding/json` around an anonymous envelope, so retain a small
  adapter for `SourceMetadata.Data.Git`, preserve raw JSON, and continue validating
  the extracted location through the project-owned allowlist contract.
  Select validation support 1.3.3 in Bzlmod to match the Go graph; the parent's
  older transitive selection must not silently replace the required version.
- Create the nested module at the user-selected project path. The user clarified
  that first-party module registration and reuse of existing dependency pins
  require no separate dependency approval. Mirror that clarification into policy
  and its owning skills.
- A native Starlark test rule symlinks the Go executable, merges tool runfiles,
  and passes configuration through the test environment. The consumer supplies
  a root marker whose resolved location identifies its worktree.
- Keep the scanner's release lock and binary toolchain inside the module. The
  default scanner label resolves the module's registered execution-platform
  toolchain, and the public CLI target uses the same selection. Keep explicit
  scanner overrides for fixtures and custom builds; ordinary consumers omit it.
  Register every generated scanner toolchain with `:all`; archive execution
  constraints determine compatibility. The parent module relies on that
  registration and declares no platform-specific alias or registration.
  The `tools/trufflehog` CLI alias forwards to the module without another pin;
  remove the obsolete root third-party package.
- Strict JSON entries match detector, path, full commit ID, and line. Reasons
  are required and duplicate coordinates are rejected. Stale entries are valid
  because a branch scan need not include every accepted historical commit.
- Stream original JSON records into annotated finding events and emit a final
  summary. Successful scanner codes are zero without findings or 183 with
  findings. Any other exit, missing metadata, malformed JSON, or contradiction
  is an error regardless of suppression counts.
- Both scan modes force `--no-verification` and `--no-update`. Scope is HEAD
  history with additional refs skipped, or all local refs when explicitly selected.
- Keep one state-machine source in the module's `docs` package. The repository
  documentation consumer renders it through the existing shared Mermaid site
  pipeline and packages the light/dark SVGs once for both the README and landing.
  The independent rules module exports source content without depending on the
  parent site's renderer. A dark WebP and a screenshot of the current allowlist
  are review artifacts, not additional maintained diagram sources.
  Give TruffleHog its own nested state machine with explicit start and end
  markers. Keep finding preservation and allowlist matching outside that
  subprocess, within the wrapper. Use direct wording and distinguish the
  scanner's 0/183 exit from the wrapper's final status based on unmatched counts.
  Keep state labels to actions and exit statuses. Put scan scope, revision-local
  ignores, exact matching fields, and completion checks in a compact README key
  linked from the landing, with links to their owning documentation sections.
  Fork after launching the scanner once: one branch tracks its lifetime while
  the other reads and classifies stdout. Keep the finding loop entirely in the
  output branch. Join the branches only after stdout EOF and process exit;
  either may finish first, and buffered findings must still be processed.

## Risks / Trade-offs

- An approved tuple covers findings from that detector at the same immutable
  commit and line; moving a fixture requires a new reviewed tuple.
- Wrapper stdout retains original values; keep local raw reports restricted.
- Missing or renamed metadata must fail rather than silently weakening matching.

## Acceptance

Define the isolated subprocess E2E cases before implementation. Exercise exact
matching and each coordinate mismatch, scanner failures after accepted output,
truncated and missing-metadata output, invalid policy, timeouts, and the required
no-verification flag. Retain a repeatable JSON report in undeclared test outputs.
Validate independent module builds/tests, repository integration, and the real scan,
then publish through the existing task PR.
