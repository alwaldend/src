---
title: Vault pki server
description: Create a pki config for a server
tags:
  - vault
  - terraform
---

`allowed_domains` defines permitted certificate names. `allow_subdomains`
defaults to `true` for existing consumers; set it to `false` to permit only
the explicitly listed names. EAB access is granted through
`eab_new_member_group_ids` for the role-scoped ACME directory.
