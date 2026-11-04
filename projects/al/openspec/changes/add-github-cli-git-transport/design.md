## Context

AL currently accepts only canonical SSH endpoints and performs network Git operations in private config-free bare directories. The forge API already uses gh. The user requested gh Git transport support and skill documentation in a separate PR.

## Goals / Non-Goals

**Goals:** Explicit gh-authenticated HTTPS selection for canonical GitHub SSH or HTTPS remotes, no changes to host or repository configuration, unchanged publication guards and default SSH behavior.

**Non-Goals:** Automatic fallback, importing global credential helpers, arbitrary HTTP endpoints, new dependencies, and SSH troubleshooting.

## Decisions

Proceed with `--git-transport gh`, retaining SSH as default. Importing global Git configuration would restore mutable URL rewriting and credential routing, so it is rejected. Fetch and push use canonical HTTPS endpoints derived only after validating the captured GitHub repository. Receipts continue binding original endpoint digests and additionally bind transport mode; legacy receipts imply SSH. Use the explicitly selected forge CLI as the sole host-scoped credential helper, shell-quoting its executable because Git runs credential helpers through a shell. Tokens flow directly between gh and Git and are never read by AL. Disable prompts and diagnostics that can expose credentials. Keep TLS verification, disabled redirects/proxies, isolated config and exact leases.

## Failure cases before implementation

- Default SSH or invalid mode must not silently enable HTTPS or fall back.
- gh mode must reject insecure, credential-bearing, escaped, noncanonical or non-GitHub endpoints before network operations.
- Ambient helpers, URL rewrites, TLS/proxy settings and diagnostic variables must not enter isolated transport.
- An executable path with spaces or shell metacharacters must remain one executable.
- Missing/failed authentication must fail without token disclosure or publication.
- A receipt prepared for one mode must be refused under another mode.
- Rebase transports must retain the explicitly selected mode and helper.
- Fetch, exact-lease push and verification must preserve existing remote state guards.

## Validation

Extend behavioral fixtures before implementation. Exercise real local Git snapshot/push operations under an audited gh transport and execute the real Git credential protocol with a dummy helper. Run the full delivery package tests, CLI help and live read-only GitHub inspection, skill packaging and offline evaluation validation, strict OpenSpec validation, root quality and semantic lint. Publish the separate PR using the candidate adapter.

## Execution evidence

The initial fixture build failed on the missing transport API before implementation. The complete delivery package and offline skill evaluation validation then passed. The discovery check correctly detected a stale generated skill projection; `//.agents:write_skills` refreshed it. A task-local candidate executable reports `git_transport: gh` and `delivery_transport_available: true`; its live read-only inspection successfully fetched master over isolated HTTPS while origin retained SSH configuration. Master advanced after the separate Vast.ai PR merged, so delivery preparation must replay this task's commit onto the captured current base.

Session review: use the existing strict JSON receipt shape whenever extending receipt fields, and regenerate skill discovery after editing canonical skill content. The local fixture exercises the real Git credential protocol with a quoted dummy executable and real local Git publication; live inspection additionally verifies GitHub HTTPS connectivity. Neither offline skill validation nor fixture routing alone proves hosted model behavior or real HTTPS publication.

The prepublication checks passed all 30 selected tests, including root quality, the full delivery fixture suite, offline skill evaluation validation and generated discovery freshness. Affected-target semantic lint and strict OpenSpec validation passed. SSH receipts retain their prior encoding; explicit gh receipts carry the new transport field. The exact prepared commit will be checked again before publication.

## Review follow-up: relative credential executable

Review thread PRRT_kwDOOgezIM6zkRZD identifies a working-directory mismatch: forge preflight executes relative CLI paths from the repository root, while Git executes the same helper string from an isolated transport directory. Extend the behavioral fixture before changing code. Resolve path-containing relative helper names against the selected forge directory; keep absolute paths and bare executable names unchanged. This retains the selected CLI identity without changing global configuration or helper quoting.

The follow-up fixture reproduced the reported failure after successful forge preflight: the relative helper could not be found from the transport directory, while the bare PATH name succeeded. After anchoring relative paths to the forge directory, both variants pass credential exchange and verified fixture publication, including executable names with spaces and shell metacharacters.
