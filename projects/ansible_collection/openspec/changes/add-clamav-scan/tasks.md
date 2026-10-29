## 1. Native role

- [x] 1.1 Record the native decision and failure cases before adaptation; preserve
      superseded Python work in ignored scratch.
- [x] 1.2 Replace custom scanning with the extracted native client, bounded
      execution, native logs, and explicit runtime compatibility checks.
- [x] 1.3 Preserve dedicated SSH forwarding and authenticated recovery; remove
      caller resets and reconcile documentation/specifications.
- [x] 1.4 Validate final packaged role, task syntax, and official source contract.

## 2. Parent integration

- [x] 2.1 Configure the daemon, signatures, stream/engine limits, and supported
      fleet selection; keep XCP-ng compatibility deferred.
- [x] 2.2 Verify rendered native configuration with real clean, EICAR, and
      oversized-file E2E fixtures; retain native logs, receipts, configuration,
      and scanner identity.
- [ ] 2.3 Complete the parent-owned production scans
      and delivery tracked in `infra/security/openspec/changes/add-clamav-scans`;
      retain operational evidence. Focused review of the three source corrections
      is complete, but production scan completion is not established.
- [x] 2.4 Verify the revised native limit policy with equal 104857601-byte
      client/controller stream limits, controller `MaxFileSize 104857600`,
      `HeuristicAlerts yes`, and `AlertExceedsMax yes`. Retain native evidence
      of an oversized-file limit detection followed by later EICAR detection
      in the same traversal, with no successful prefix verdict.
- [x] 2.5 Verify controller-only ControlPath allocation and persistent master
      establishment before any endpoint command, single capture/restoration of
      caller arguments, explicit own-master mux operations, same-master
      collection, and private retention of both boot-command results including
      unreachable errors. Rerun applicable parent-owned validation.

Checked items record the preceding implementation and validation. In
particular, 2.2 used the former 104857600-byte server stream quota; its E2E
evidence is superseded for the revised policy. Task 2.4 passed locally against
the revised limits before the early-master change; it does not validate SSH.
Task 2.5 and the forwarded XOA canary passed against the final role; production
scan outcomes and delivery remain under parent ownership.

## 3. Shared SSH hardening

- [x] 3.1 Add the overridable shared `ansible` match default and distinct
      mapping to the existing devsec role.
- [x] 3.2 Extract the hardening entrypoint with required caller certificate
      paths, preserving normal setup and upstream handler behavior.
- [x] 3.3 Verify default, empty/custom overrides, required certificates,
      normal certificate derivation, and actual upstream template rendering
      offline; validate role/docs packaging and syntax.
- [x] 3.4 Reject missing or non-sequence certificate inputs before the native
      include, permit intentional empty lists, and verify the failure boundary
      offline. Move only this fixture's evidence into the original task root.
- [x] 3.5 Record acceptance of the parent-owned guarded native rollout on all
      nine approved hosts and the secure-runner reverse-Unix canary, including
      fresh SSH, rejected TCP forwards, preserved CA/boot, and owned cleanup.

The parent completed the approved nine-host rollout and bounded secure-runner
canary; see validation.md for receipt-backed acceptance and limitations. This
does not complete the full-fleet production scans or delivery in task 2.3,
which remains pending. The scanner never changes SSH server policy.

Items 3.2–3.4 describe the completed rollout implementation. The PR
simplification removes that helper and its one-time caller; normal SSH setup
now passes the shared default directly to the native hardening role.

## PR simplification

- [x] Simplify the SSH wrapper to its existing native hardening include and shared forwarding default; remove the completed rollout-only certificate helper.
- [x] Simplify native scan preparation and evidence collection without losing diagnostic results, private retention, transport ownership, or honest coverage reporting.
- [x] Validate the simplified role behavior against native tools and package the exact candidate; do not modify active scan runtimes.

## Selected-host performance controls

- [x] Add explicit background scheduling and native multiscan inputs with unchanged defaults and receipt fields; validate native session scanning and deploy the authorized single-host run with dedicated controller resources.

## Complete Ansible review

- [x] Discover and exclude actual block/character device nodes without opening them or excluding regular files under `/dev/shm`; retain the exclusions in the receipt.
- [x] Replace the asynchronous ClamAV exit marker with the finished native Ansible result, retaining synchronous password-transport workarounds.
- [x] Exercise device exclusions, native detection/error/timeout results, evidence retention, and failure cleanup; review all Ansible changes and callers before delivery.

## Native simplification

- [x] Use a short private controller runtime directory for SSH IPC, remove CWD/relative-path tricks, and compute the outcome once from registered results.
- [x] Verify successful and failed execution, private retention, owned cleanup, and the packaged candidate.
