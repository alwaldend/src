---
title: GitHub Terraform
description: GitHub Pages repositories for project landing sites
---

This package creates the public repositories that host the project landing
pages. Each repository receives its built Hugo site from the reusable deploy
target `//projects:deploy_landings` in the root workspace. Terraform also manages
the custom domain and Pages source (`pages`, `/`).

For a new site, first plan and apply with `bootstrap_projects` containing only
the new projects whose `pages` branches do not yet exist. For example, pass
`-var='bootstrap_projects=["new_project"]'` to `//infra/github/tf:tf.plan` and
the reviewed `tf.apply`. This creates the repository without trying to enable
Pages on a nonexistent branch. Do not include an existing live site: doing so
would remove its Pages configuration.

Publish the built site with `//projects:deploy_landings`, then plan and apply
with the default empty `bootstrap_projects` to enable Pages. Keep all applies
limited to reviewed, authorized repository and Pages changes. DNS is managed
separately by `//infra/dns`; validate the public custom domain after rollout.
