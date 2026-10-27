## Why

The authorized site rollout fails because the repository documentation aggregate
references a missing Go gRPC documentation target. The user also requested an
audit of all docs targets and exclusion of generated releases and project
binaries from the website archive.

## What Changes

- Restore the missing documentation target at its existing package owner.
- Audit every root-workspace docs target and the website's actual packaged files.
- Exclude generated release documentation from the website graph; downloadable
  artifacts remain on the download service.
- Preserve ordinary documentation, images, fonts, CSS, and JavaScript.

## Capabilities

This repairs packaging and verifies the existing documentation boundary;
there is no new user-facing capability.

## Impact

Repository documentation BUILD declarations and the apex website consumer.
No external dependency is added. Release publication and host configuration
retain their separate owner-local changes.
