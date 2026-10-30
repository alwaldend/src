## Why

The user requested the third-party dependencies from security PR #214 as a
separate pull request targeting master. Lynis needs an owner-local record so
its dependency can be reviewed independently of the consuming audit role.

## What Changes

- Extract the existing Lynis package from commit
  `5e4bf7ccb7e7948e12c12ae12b92d048bbba06e3`, preserving its version, URL,
  integrity, archive bytes, and Bazel interface.
- Retain the generated root module include and remove the README link to the
  consumer role that is outside this pull request.
- Verify the extracted archive, executable version, and repository gates.

## Capabilities

### New Capabilities

None. This extraction preserves the existing dependency contract.

### Modified Capabilities

None. `skip_specs: true` records that the extraction has no behavioral delta.

## Impact

Only `third_party/com_github_cisofy_lynis` and its generated root module include
are owned here. The companion ClamAV package retains its own change record.
Consumer roles, host configuration, and live audits remain outside this change.
