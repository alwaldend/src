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

## Reprovision a landing-site certificate

GitHub issues the custom-domain certificate as part of its Pages build and
offers no direct reissue action. To force a new certificate for one live site,
plan and apply with only that project in `certificate_reprovision_projects`:

```sh
bazel_agent bazel run //infra/github/tf:tf.plan -- \
  -var='certificate_reprovision_projects=["affected_project"]'
```

Inspect that the plan removes only that project's Pages block, then apply the
reviewed changes. Restore the site with a second plan and apply that leaves both
`bootstrap_projects` and `certificate_reprovision_projects` empty. GitHub
schedules a new certificate once the custom domain returns.

Republish the site with `//projects:deploy_landings` only if the first apply
removed its `pages` branch content; normally the branch is unchanged. Verify the
pushed revision, Pages build state, DNS, and the certificate served for the
custom domain before closing the recovery.

Keep this variable empty outside an active recovery. Never name an existing
served site in `bootstrap_projects` to force a reissue.
