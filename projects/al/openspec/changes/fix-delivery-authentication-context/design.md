## Context

See proposal.md for motivation. Direct read-only checks on 2026-11-04 reproduced HTTP 401 inside the sandbox and success outside it. `gh auth token` availability was checked without displaying its output: the sandbox could not retrieve the stored credential. Neither token environment variable was set. AL preserves the GitHub CLI credential environment, so this is an execution-context boundary rather than AL changing credential precedence.

## Goals / Non-Goals

Keep credentials in the GitHub CLI and preserve every publication guard. Do not modify the shared sandbox, provision credentials, automatically change Git transport, or make AL request execution escalation.

## Decisions

Proceed with a focused diagnostic and procedure correction. For a preflight diagnostic containing the GitHub CLI's `HTTP 401` status marker, append conditional recovery guidance and wrap the original error. Other failures keep their existing behavior. The marker selects advice only, never authorization or retry behavior.

Use `gh auth status --active --hostname HOST --json hosts` with an explicit projection of active/state/tokenSource for comparisons. JSON status exits zero even for invalid authentication; the procedure must inspect `state`. Never request `--show-token`, print raw credential files, or advise a new login based solely on a sandbox failure.

Alternatives rejected: changing transport cannot repair API authentication; copying tokens into the environment weakens credential isolation; changing host permissions affects concurrent users and needs a separate managed configuration migration. AL cannot itself grant sandbox permissions, so the agent uses the existing approval mechanism.

## Risks / Trade-offs

- HTTP 401 does not prove sandbox denial or token expiry: keep the advice conditional and compare contexts.
- CLI wording can change: an unrecognized marker retains the original safe failure; skill guidance remains available.
- Execution approval can fail: stop and report that boundary rather than bypassing it.

## Failure cases before implementation

- HTTP 401 is confused with SSH/HTTPS failure, causing an ineffective transport switch.
- A non-authentication failure receives misleading credential advice.
- Recovery guidance loses the original error category or exposes a dummy credential from stderr.
- A failed preflight continues to Git network operations, mutates HEAD, or retries automatically.
- JSON auth status returns zero with `state: error` and is mistaken for success.
- Outside-context failure or approval rejection triggers repeated login or escalation attempts.
