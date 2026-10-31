## Why

The requested host security review needs a reproducible Lynis audit with private,
retained evidence and no persistent scanner installation.

## What Changes

- Package the approved stable CISOfy Lynis source archive with an integrity pin.
- Add a reusable, explicitly invoked Ansible audit role with a bounded runtime,
  private controller reports, and cleanup of its remote audit directory.

## Capabilities

### New Capabilities

- `lynis-audit`: Run the packaged scanner on demand and retain private evidence.

### Modified Capabilities

None.

## Impact

Adds `roles/lynis_audit`. The dependency declaration is owned by
`third_party/com_github_cisofy_lynis`; a consuming playbook supplies the packaged
archive path and an absolute controller report directory. No live execution,
authentication, periodic scheduling, hardening, or automatic remediation is part
of implementation validation.
