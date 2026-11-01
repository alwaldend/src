## 1. Behavior and acceptance

Sections 1–3 record the initial implementation. The Go sync command and its
helpers are now removed. Infrastructure owns the replacement playbook and
adapted acceptance fixture through
[task 19](../../../../../infra/download/openspec/changes/add-static-hosting/tasks.md#19-infrastructure-owned-synchronization).

- [x] 1.1 Record transfer, trust, authoritative deletion, and failure boundaries before implementation in the design and specification.
- [x] 1.2 Add an isolated CLI-to-SSH/rsync fixture before implementation and demonstrate that the missing command fails it.

## 2. Command and integration

- [x] 2.1 Implement source capture, validation, rsync publication, private tree staging, atomic project-directory exchange, and dry run; pass the isolated fixture.
- [x] 2.2 Expose the authenticated download sync target under the download owner's change; build it and inspect its configured direction and credential references.
- [x] 2.3 Document the command, temporary-space requirement, authoritative deletion, failure behavior, and production-execution boundary.

## 3. Delivery

- [x] 3.1 Retain the end-to-end receipt, run formatting, package checks, semantic lint, and repository quality gates on the final candidate.
- [x] 3.2 Publish and verify the updated pull request with the implementation and evidence.

## 4. Local version collisions

- [x] 4.1 Add real SSH acceptance cases for duplicate uploads, existing extraction, concurrent same-version publication, and unchanged local trees after refusal before implementation.
- [x] 4.2 Publish new local releases atomically and reject existing versions, retaining authoritative public replacement through sync.
- [x] 4.3 Inject short-lived credentials into publication targets, document the collision contract, and validate the complete local-to-public path.

## 5. Extracted website retention

- [x] 5.1 Specify failure boundaries and extend the real SSH fixture before implementation.
- [x] 5.2 Remove obsolete extracted versions only after successful selection, retain archives, and verify authoritative sync carries the reduced tree.
- [x] 5.3 Document rollback from archives and validate and publish the updated candidate.
- [x] 5.4 Reproduce an inaccessible inherited login directory in the SSH fixture, fix remote working-directory handling, and retry only cleanup on the selected local release.
- [x] 5.5 Verify restrictive-umask publication and repair the selected local archive directory through the release tool, retaining both archive hashes.

## 6. Rebase and cleanup review

- [x] 6.1 Preserve trunk namespace support, third-party replacement uploads and explicit SSH-agent controls while replaying project-only collision refusal.
- [x] 6.2 Exercise mounted old-release refusal through the real SSH fixture before correcting cleanup; verify mounted private state and selection remain unchanged.
- [x] 6.3 Correct mount-boundary preflight and reconcile recovery/DNS documentation.

Exact-candidate validation and publication status belong to the guarded delivery
receipts under `out/rebase-review/`; all required gates run before publication.
