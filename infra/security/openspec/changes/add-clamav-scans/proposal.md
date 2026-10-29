## Why

The authorized host scan needs current ClamAV signatures without loading the
database on memory-constrained endpoints or coupling scans to maintenance.

## What Changes

- Package thin fleet and appliance imports of one shared three-play ClamAV
  lifecycle, retaining separate host defaults, connection variables, inventory
  exports, and authentication boundaries.
- Prepare an unprivileged temporary controller daemon, preserve one private
  freshclam cache, verify signatures, and retain lifecycle evidence using
  native Ansible modules, ClamAV tools, and a transient systemd user unit.
- Use one bounded daemon VERSION polling task to match the verified daily
  database and reject the client's local-version fallback. Keep signature
  metadata in `signatures.json` and the actual daemon version in lifecycle
  evidence; remove the redundant local engine-version loop and `engine.json`.
- Invoke the shared streaming scan role and stop the owned daemon after success,
  failure, or its bounded lifetime. Interrupted playbooks rely on the unit's
  maximum runtime. The appliance default is XOA; XCP-ng is unscanned because
  the native payload requires a newer glibc.
- Reuse the existing SSH hardening role's forwarding default. Retire the
  completed rollout's dedicated maintenance command and helper, preserving
  historical evidence. Normal owner deployment retains the upstream restart
  handler and does not provide the retired rollout's reload-only or rollback
  workflow.

## Capabilities

### New Capabilities

- `clamav-controller`: Private verified controller lifecycle and packaged scans.

### Modified Capabilities

None.

## Impact

Only `infra/security` source and its OpenSpec workspace. Consume the separately
owned pinned ClamAV payload and `alwaldend.main.clamav_scan` role. Existing
maintenance behavior is unchanged. Live operation and final delivery belong to
the parent task. Historical runtime evidence remains tied to its captured
candidate; simplified-package validation and production acceptance stay pending.
