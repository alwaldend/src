# Diagnose GitHub authentication

GitHub API preflight and Git transport are separate. `--git-transport gh`
selects HTTPS for Git; both Git modes use the selected forge CLI for API calls.
HTTP 401 alone proves neither credential expiry nor a transport problem.

Apply `repo-secrets`. When a sandboxed preflight reports an authentication
failure, compare sanitized active-account status in the failing context and,
through the execution approval mechanism, outside that sandbox. Use the selected
`--forge-cli` executable and captured host; for the default CLI on github.com:

```sh
gh auth status --active --hostname github.com --json hosts \
  --jq '[.hosts[][] | {active, state, tokenSource}]'
```

Inspect `state`: JSON status can exit zero even when authentication failed.
`tokenSource: default` identifies the CLI's default credential lookup; it does
not prove the stored credential was accessible. Never request `--show-token`,
print raw credential files, or copy tokens into arguments or the environment
to work around the sandbox.

If the approved-context check succeeds while the sandboxed check fails, run
the same AL entry point in that working context through execution approval,
retaining the selected transport throughout the receipt-bound workflow. Do
not ask for a fresh login, change transport to fix the API failure, or change
global Git, CLI, or sandbox configuration. Approval rejection remains a stop:
report it, diagnose the stated reason, and do not route around it.

If authentication also fails in the approved context, diagnose that failure
before requesting `gh auth login` or a credential refresh. If the comparison
cannot be authorized, report credential access as unverified rather than
claiming the login expired. An API-only or transport-only success does not
prove the other channel works; verify both through normal delivery inspection.
