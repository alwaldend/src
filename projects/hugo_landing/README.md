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
