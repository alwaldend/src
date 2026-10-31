## Why

Small endpoints cannot reliably hold the ClamAV engine and signature database.
The official native client can stream files to a private controller without
installing packages or maintaining a custom file walker and scanner protocol.

## What Changes

- Add a bounded Ansible role using the official pinned `clamdscan --stream`.
- Unpack the official DEB data payload into private temporary endpoint storage.
- Forward only the controller Unix socket through a dedicated SSH master,
  established from a controller-only ControlPath before any endpoint command;
  capture caller arguments once and preserve unrelated connections.
- Retain native logs, summary counts, and private lifecycle receipts without
  claiming host cleanliness or complete coverage. Retain both boot-command
  results privately, including failed and unreachable errors.
- Keep the client and controller stream limits equal at 104857601 bytes and
  the engine file limit at 104857600 bytes. Enable `HeuristicAlerts` and
  `AlertExceedsMax` so oversized files produce native limit detections and
  traversal continues without a successful truncated-prefix verdict.
- Preserve native infected counts as detections, including limits and
  encrypted-content heuristics, rather than treating every finding as malware.
- Default the shared SSH wrapper to the approved three-rule `ansible` match
  through its existing devsec include, with an empty-list override. Normal
  owner deployments apply the policy; the scanner never changes server policy.

## Capabilities

### New Capabilities

- `clamav-scan`: Private native streaming with bounded execution and retained
  operational evidence.
- `ssh-hardening`: Shared forwarding defaults in normal SSH setup.

### Modified Capabilities

None.

## Impact

Owns `roles/clamav_scan`, `roles/ssh`, and this collection change. The caller
owns the daemon, signatures, configuration, fleet selection, and live
deployment. The initial fleet is the 11 supported hosts; older XCP-ng
compatibility is deferred.
The scanner adds no endpoint package installation, custom scanning helper,
quarantine, deletion, reboot, or public TCP listener.
