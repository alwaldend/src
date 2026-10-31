## Context

See `proposal.md`. The existing Bazel tests own scan arguments and exclusions.
`repo-secrets` owns credential handling; this skill owns scan execution and triage.
The user explicitly requires original findings for local inspection.

## Goals / Non-Goals

Preserve enough evidence to inspect every finding and reproduce the scan.
This change does not rotate credentials or rewrite history. The companion
TruffleHog change excludes only the four findings the user confirmed as false
positives after the first updated scan.

## Decisions

- Keep the repository-wide procedure in `tools/agents/skills/repo-secret-scan`.
  Link the scanner README for scope ownership instead of duplicating exclusions.
- Retain raw output in restricted ignored workspace scratch and provide the path.
  Blanket redaction or automatic deletion before inspection obstructs triage;
  publishing raw credentials would create unnecessary additional exposure.
- Use the existing skill library, discovery generator, and offline evaluation
  rules. Live model evaluation cannot reproduce local Git and scanner execution.
- Prefer native same-line ignores for reviewed current synthetic fixtures;
  retain exact historical exclusions. The scanner README owns the protocol and
  explains that ignored lines are omitted before JSON output and can be inspected
  explicitly through `--no-ignore-tag` with verification disabled.

## Risks / Trade-offs

- Unverified findings can still be real credentials; inspect source context.
- Git history scans omit uncommitted and untracked content; report that limit.
- Retained raw output can contain credentials; keep it local, outside tracked
  source, with restrictive permissions, and remove it when inspection is complete.
