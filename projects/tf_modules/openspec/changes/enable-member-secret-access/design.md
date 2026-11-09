## Decision

Reuse the existing namespace policy and member group. An optional allow_member_secret_access flag adds the existing approle_secrets policy to group policies. Its default false preserves all component-role boundaries; only the personal user role enables it. This avoids duplicating ACL text or introducing a second identity group.

## Verification

The targeted live plan must update only the personal group's policy list in addition to the separately approved team-delivery grant. After authorized deployment, inspect the existing human login's mount visibility and own-secret capabilities without reading values. Verify that unrelated user namespaces remain denied. Offline formatting, owner validation, and semantic lint cover source structure.

## Metadata navigation

Direct access to the namespace did not establish UI navigation: the mount root could not be listed and the namespace without a trailing slash reported deny. For opted-in roles, derive exact metadata ancestor paths from secret_path and grant only list on each ancestor, with and without a trailing slash. Add list on the exact namespace path without a slash; the existing wildcard continues to govern the namespace and descendants. Ancestor folder names become visible, but secret values outside the namespace remain denied. Default-false roles receive no navigation rules.

Verify every parent from the mount root to imagegen, both path spellings in capability preflight, unrelated value denial, and unchanged token metadata version. Include a default-false component policy in the read-only plan as a negative control.
