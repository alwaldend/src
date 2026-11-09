---
title: Vault approle
description: Create a vault approle
tags:
  - vault
  - terraform
---

By default the member group grants AppRole bootstrap permissions. The role's
own KV namespace policy stays on the role token. Set
`allow_member_secret_access = true` for a personal role whose human members
should also access that namespace directly, including through Vault's UI.
This reuses the existing namespace policy and group; other roles retain the
default bootstrap-only member permissions.

Opted-in roles also receive list-only permissions on exact metadata ancestors
so the UI can navigate from the mount root. Those listings reveal folder names;
they grant no secret-value access outside the role's own namespace.
