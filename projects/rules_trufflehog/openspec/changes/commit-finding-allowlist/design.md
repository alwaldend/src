## Context

See `proposal.md`. The current path filter hides whole files. A scan-only
inventory reports ten previously excluded findings. The user requires commit
and line matching and explicitly selected `projects/rules_trufflehog` as owner.

## Goals / Non-Goals

Provide a standalone module with consumer-supplied tools and exact suppression.
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
- Strict JSON entries match detector, path, full commit ID, and line. Reasons
  are required and duplicate coordinates are rejected. Stale entries are valid
  because a branch scan need not include every accepted historical commit.
- Stream original JSON records into annotated finding events and emit a final
  summary. Successful scanner codes are zero without findings or 183 with
  findings. Any other exit, missing metadata, malformed JSON, or contradiction
  is an error regardless of suppression counts.
- Both scan modes force `--no-verification` and `--no-update`. Scope is HEAD
  history with additional refs skipped, or all local refs when explicitly selected.

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
