## Source validation

Supplied run: commit `8b91edb37de21007becef5a37d41ffef2371f102`,
2026-11-16. Local verification: 2026-11-17.

- The pinned runner accepts `.forgejo/workflows/ci.yaml` with `fetch-depth: 0`.
- `//projects/sri/releases:head` builds a populated release manifest with
  complete local history, exercising the failed release-generation boundary.
- Twelve tests pass across the CI action and executable, diagram source
  comparisons, runner registration/template acceptance, Terraform formatting,
  and Vault-login URL rejection.
- The owning diagram update target regenerated all stale projections; the
  WebP output was visually inspected.
- The packaged Ansible syntax check and deployment binary build pass.
- A local playbook renders the exact secure1 environment-merge task and actual
  runner and transport templates using dummy credentials. Reserved password
  characters survive URL encoding, runner timeout/cache defaults remain intact,
  both proxy-variable cases agree, and local/apex/nested-domain bypass succeeds.
  External and deceptive suffix hosts remain proxied.
- A local stream fixture executes the transport command from its source
  template: trusted TLS forwards the HTTP CONNECT bytes unchanged, while an
  untrusted certificate and wrong hostname fail before forwarding. The fixture
  found and corrected OpenSSL interactive command interpretation by adding
  `-nocommands`. Systemd unit verification passes.
- Read-only secure1 inspection confirms OpenSSL 3.5.9 supports `-nocommands`,
  the proxy hostname resolves, the worker is active and `zip` is absent.

Reproducible local receipts and scripts remain under ignored `out/forgejo-ci/`:
`tls-transport.json`, `verify_tls_transport.py`, `proxy-render.yaml`,
`proxy-selection.json`, and bounded command logs. They use dummy credentials.
Original scan findings are retained privately as `findings-baseline.raw.log`.

## Secret finding disposition

Complete reachable history reports fourteen findings before the new exclusions:
ten already suppressed and four unmatched. Each unmatched historical source
line was inspected at its full commit ID and checked against the original
finding. Three are synthetic URL-rejection credentials and one is a recorded
ComfyUI job UUID. Exact coordinates are in `tools/trufflehog/allowlist.json`.
The supplied shallow run attributed source to the merge commit; full checkout
restores actual introduction coordinates, including an already accepted NPM
detector collision in recorded invocation evidence.

## Final source gates

The repaired scan completes with fourteen total, fourteen suppressed and zero
unmatched findings; verification remains disabled. Thirty-three root checks
pass, covering repository formatting, Buildifier, all three owner OpenSpec
workspaces and the affected Vast.ai acceptance suite. AL delivery tests pass
from the owning nested workspace. Root semantic lint passes across the affected
configuration and source targets. The nested module has no `lint` configuration;
its affected Go source is also checked through the root module dependency.

## Live boundary

The user subsequently authorized the prepared policy and secure1 deployment.
The owner-local rollout receipts in
`infra/dev/openspec/changes/configure-runner-proxy/deployment.md` and
`infra/vault/openspec/changes/allow-runner-proxy-read/deployment.md` record
successful application, converged reruns, a checksum-verified Bazel fetch,
six successful upstream requests and direct local-forge access.
No complete new CI run was performed. The master-only workflow will not run
for the feature pull request; merging remains outside task authority.

## Session ergonomics

`ci-log-classification`: initial text matching repeated dependency fan-out
errors and overproduced output. The bounded unique-root-cause extraction in
ignored scratch resolved it; future diagnosis should classify downloader,
action and test summaries before reading individual lines.

`worktree-lfs-read-storage`: Git status attempted writes to the shared LFS
temporary directory under the filesystem sandbox. Read-only Git commands used
task-local LFS storage; authorized delivery uses managed storage. No shared
configuration was changed. These are task-local observations, not proposed
changes to canonical policies or tooling.
