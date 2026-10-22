# SSH publication acceptance

The failure matrix was established before implementation. The user subsequently
requested rsync/rclone rather than a custom transfer protocol. The cases below
reflect that decision; immutable-file conflicts and a digest database are no
longer part of the contract.

| Case             | Observable acceptance                                                                                               |
| ---------------- | ------------------------------------------------------------------------------------------------------------------- |
| Host trust       | Unknown and mismatched host keys fail before publication.                                                           |
| Identity         | SSH uses an administrator; sudo runs rsync and activation as download; no publisher authorized keys are configured. |
| Ordinary release | Files and archives retain bytes and filenames, without extraction or an extra files directory.                      |
| Updates          | Identical repeats work; changed ordinary filenames follow rsync updates; unrelated retained files remain.           |
| Paths            | Reject traversal and control characters; spaces, Unicode, and shell punctuation remain literal filenames.           |
| Website          | Tar.gz and zip archives remain downloadable; HTML/assets are readable through the stable site URL.                  |
| Extraction       | Escaping/link/device/duplicate/malformed entries and missing index.html fail before selection.                      |
| Failures         | Interrupted or full-disk transfers preserve the selected site; completed ordinary files may remain.                 |
| Coordination     | Unique staging and a promotion lock preserve completed site versions during concurrent deployment.                  |
| Reselection      | Selecting a retained website changes only current and does not require local payload bytes.                         |
| Backends         | SSH does not need Soras or Ansible; OCI continues invoking Soras with the original tag/file.                        |
| Targets          | Both infra targets package their environment and host, and accept a runtime release directory.                      |

Verification records belong in evidence.md after the checks run. The earlier
custom-receiver implementation was removed before delivery.

## Independent publication steps (review acceptance, before implementation)

- Upload alone copies the manifest files without extracting or selecting a site.
- Extract alone reads the published archive even when local payloads are absent;
  it creates the versioned site without changing the selected release.
- Link alone selects an existing version without reading an archive, invoking
  rsync, or requiring local payload files; an absent or invalid site fails.
- Combined steps execute upload, extract, then link; duplicate, unknown, empty,
  or out-of-order step names fail before a remote change.
- An extraction or transfer failure preserves the current link; a repeated
  extraction preserves an existing completed site version.
- Quotes, spaces, Unicode, shell punctuation, and SSH options remain literal
  through POSIX shell quoting and rsync's transport parser.
- Errors identify the failed step, operation, and relevant path while retaining
  their underlying causes. Invalid archive checks keep their original causes.

Extend the isolated SSH/HTTP fixture for these cases before changing source.
Keep existing unknown-host, update, archive-safety, concurrent-transfer, disk
exhaustion, wrapper, and OCI coverage. The new repo-go skill's offline cases
cover contextual error wrapping, Bazel invocation, and dependency ownership.
