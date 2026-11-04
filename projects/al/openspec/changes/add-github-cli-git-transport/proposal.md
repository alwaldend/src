## Why

SSH delivery can fail even when GitHub CLI authentication works. Add an explicit gh-authenticated HTTPS transport while retaining delivery guards and the SSH default.

## What Changes

- Add `--git-transport gh` to use the selected GitHub CLI as an isolated HTTPS credential helper.
- Bind transport selection to preparation receipts and preserve endpoint identity, lease, TLS and redirect guards.
- Document selection and authentication recovery in the delivery skill and extend its evaluation cases.

## Capabilities

### New Capabilities

- `github-cli-git-transport`: Explicit gh-authenticated GitHub HTTPS delivery.

### Modified Capabilities

None.

## Impact

AL delivery commands, isolated Git transport, preparation receipts, operator documentation and al-delivery skill. No dependencies, host configuration, or remote configuration changes. Separate PR from the Vast.ai integration.
