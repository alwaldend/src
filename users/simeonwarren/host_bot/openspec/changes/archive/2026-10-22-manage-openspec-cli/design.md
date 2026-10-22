## Context

See [proposal.md](proposal.md). The playbook already provisions Node.js/npm
and installs Open Computer Use through `community.general.npm` as the managed user.

## Goals / Non-Goals

Reproduce the manually installed CLI through the existing npm installation
pattern. Deployment and repository CLI-wrapper changes are outside this task.

## Decisions

- Pin OpenSpec to 1.14.2, resolved from the primary npm registry on 2026-10-16,
  instead of retaining the moving `latest` tag. Its Node.js requirement is
  `>=20.19.0`; the observed host version is 22.23.1.
- Give OpenSpec its own version and npm-prefix inventory variables. Share the
  Node.js/npm prerequisite task with both npm tools, including their tags.
- Disable npm lifecycle scripts, matching the existing installation pattern.
  The published OpenSpec package has no install lifecycle scripts.

## Risks / Trade-offs

- npm resolves transitive packages at installation time → Pin the direct CLI
  version and retain the repository's established host npm mechanism.
- Offline packaging cannot prove host installation or idempotency → Inspect
  the packaged task and merged variables; leave live deployment to an explicit request.
