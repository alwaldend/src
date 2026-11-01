## Why

The security documentation needs a diagram that distinguishes target-side
Lynis auditing and ClamAV file streaming from controller-side malware analysis.
Readers must not infer that the full ClamAV engine runs on small target hosts.

## What Changes

- Add one authoritative Mermaid source under `infra/security/docs/`.
- Render site-themed light and dark SVG assets through the root package.
- Reference the diagram immediately after the security README frontmatter.
- Show Ansible execution and report-copy edges, scanner-produced reports,
  and SSH streaming between the target client and controller engine.

## Capabilities

No behavior changes. This documentation-only change declares
`skip_specs: true`; scanner behavior remains owned by its implementation.

## Impact

Only the security diagram source, root documentation build wiring, README
image reference, and this change record are affected. No live operations,
new dependencies, authentication changes, commits, or pushes are included.
