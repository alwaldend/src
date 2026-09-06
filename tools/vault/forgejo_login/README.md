---
title: Forgejo login
description: Get a short-lived Forgejo token using OIDC
languages:
  - go
tags:
  - forgejo
  - vault
---

The plugin registers token removal before issuance, using a unique token name.
Startup failure and normal shutdown attempt removal through the authenticated
Forgejo browser session, including when creation did not return a usable token.
Cleanup parses the complete applications settings page and verifies token absence
after deletion. It then logs out the invocation's browser session and verifies
that the retained original cookie no longer accesses authenticated settings.
Cleanup failures are reported.

The parser follows the [Forgejo 15.0.3 applications template](https://codeberg.org/forgejo/forgejo/src/tag/v15.0.3/templates/user/settings/applications.tmpl): all tokens appear on
one page, each in a `flex-item` row with a title and a matching delete button.
It requires the applications page marker and new-token link before accepting an
empty list. Login pages, invalid token IDs, duplicates, and responses larger than
4 MiB fail closed. A changed upstream template may require a parser update.
Network requests have a ten-second timeout and honor cancellation. Forced process
termination or an unavailable Forgejo server can prevent cleanup; Forgejo tokens
do not acquire an expiry merely because this plugin created them.
