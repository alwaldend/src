---
title: Hugo Landing
description: Reusable Hugo landing site for project pages
statuses:
  - active
tags:
  - hugo
---

`hugo_landing` packages a reusable Hugo landing site. `al_hugo_landing`
combines the shared source with one project README and emits a Hugo-ready
source archive.
`al_hugo_landing_site` also builds the rendered site.

## Usage

```starlark
al_hugo_landing(
    name = "landing_site",
    docs = "//projects/goal:docs",
    project = "goal",
    title = "Goal",
    docs_url = "https://alwaldend.com/docs/projects/goal/",
    repository_url = "https://github.com/alwaldend/src",
)
```

```starlark
al_hugo_landing_site(
    name = "site",
    docs = "//projects/goal/docs:docs",
    project = "goal",
    title = "Goal",
    docs_url = "https://alwaldend.com/docs/projects/goal/",
    repository_url = "https://github.com/alwaldend/src",
)
```

## Site collection

Each project owns its rendered landing target at
`//projects/<project>/landing:site`. The shared macro remains in
`hugo_landing` so the visual shell and README-to-site transformation stay
reusable.

Agent workflow: [Add a project site](https://github.com/alwaldend/src/blob/master/projects/hugo_landing/skills/add-project-site/SKILL.md).

`projects/projects.bzl` lists the dedicated project sites; the
[project directory](../README.md) links to each site. The existing
`alwaldend.com` apex site covers that project. Nested rules modules use
`//projects/hugo_landing/landing:<project>` targets in the root workspace.

README-relative links lead to their source files on GitHub, and relative images
load from the matching source directory. Literal template examples are preserved;
Hugo shortcodes requiring the full documentation site are omitted.

## Publish project sites

Build and validate all sites from the repository root:

```sh
bazel_agent bazel test //projects:landings //infra/dns:config_test
```

The authenticated deployment command uses the operator's Git executable and
GitHub authentication. It intentionally accesses GitHub at runtime; builds do
not deploy. Supply an absolute, task-owned scratch directory:

```sh
bazel_agent bazel run //projects:deploy_landings -- --scratch /absolute/worktree/out/site-rollout/git
```

Each site is pushed to `alwaldend/<project-slug>-landing` on the `pages` branch
with `CNAME` and `.nojekyll`. Existing branch history is preserved, unchanged
sites produce no commit, and empty repositories can be bootstrapped. Use
`--dry-run` to stage and compare without pushing. The deployment tool validates
every input before publishing and stops on the first failed project.
Add `--project <project_id>` to publish only one registered site.

Repository creation and Pages settings belong to
[GitHub Terraform](../../infra/github/tf/README.md). DNS records belong to
`projects/<project>/dnsconfig.json` and are applied through
[DNSControl](../../infra/dns/README.md). Review each plan or preview and apply
only the authorized project-site changes. Verify the published commit, Pages
build, DNS, and HTTPS response before calling a rollout complete.
