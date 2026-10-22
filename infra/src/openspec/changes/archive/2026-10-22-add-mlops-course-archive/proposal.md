## Why

The user requested Yandex-Practicum/mlops-course in third_party. Upstream has no
tags, so the user authorized a commit archive instead.

## What Changes

- Add a checksum-pinned source archive and repository-internal source alias.
- Document upstream provenance and MIT licensing; regenerate the root module.
- Verify archive contents and the Bazel source target, then deliver the change.

## Capabilities

No specification requirements change. This dependency registration uses
`skip_specs: true`.

## Impact

Adds `third_party/com_github_yandex_practicum_mlops_course` and its generated
root module include. Course environments and lessons remain upstream content.
