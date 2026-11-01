# Diagram validation

Observed on 2026-10-29 at 20:10 UTC in the supplied
`ops/host-security-updates-20261029` worktree. The authoritative source is
`infra/security/docs/security-scans.mmd`, SHA-256
`b31eeab4e09f37376788496cdd9ef854700d4212589e9c718d1bbbedce78b47e`.

## Checks

- `bazel_agent bazel build //infra/security:security_scans //infra/security:security_scans_preview //infra/security:security_scans_preview_dark //infra/security:docs`
  passed. The SVG target uses the maintained site theme, embedded font, and
  native renderer; WebP targets use the same source and natural proportions.
- A `bazel_agent bazel cquery` of `//infra/security:docs_assets.files` verified
  the actual package destinations:
  `content/docs/infra/security/assets/security-scans.svg` and
  `content/docs/infra/security/assets/security-scans.dark.svg`.
- Inspected the source and both rendered appearances. Browser measurements
  confirmed a `1040 × 513` SVG view box, two groups, six nodes, and seven
  paths. All nodes are contained in their intended group; title plates fit
  horizontally within their containers and remain inside the canvas. No
  path crosses a title plate. Titles straddle the top container border as
  specified by the shared theme.
- The bidirectional file-stream/result edge avoids the initial separate
  reverse edges routing behind nodes. Both `Run`, both `copy`, and both
  `Produces` labels remain readable. ClamAV's engine and fresh signatures
  are visibly on the controller; the ClamAV client is on the target. The final
  label follows the native-client architecture and contains no Python or
  controller-helper reference.
- Scoped buildifier check and `git diff --check` passed. Strict validation
  of `add-security-scan-diagram` passed. No full-repository formatter ran.

Private evidence and stable render copies are retained under
`out/host-security-updates-20261029/clamav/diagram/`. The two `geometry-*.html`
files embed the unmodified generated SVGs and repeat the browser measurements.
`sha256.txt` identifies the source and captured outputs. The parent owns final
inline preview publication. No live infrastructure operation or Git publication
was performed.

## Session ergonomics

`diagram-evidence-discovery`: an unnecessary full worktree inventory and
oversized combined reads caused truncated output. Prefer the current-checkout
query and bounded skill reads already used later in this session.

`diagram-output-capture`: generated files disappeared from the workspace's
`bazel-bin` path during concurrent work. The parent confirmed no clean command;
the cause is unverified. A bounded path query preceded restoring cached outputs
and immediately copying them into task scratch. Future capture should record
absolute output paths and retain preview copies immediately after the build.
No shared policy changes are proposed by this documentation task.
