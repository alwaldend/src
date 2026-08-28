---
title: Hugo
description: Hugo rules
languages:
  - bzl
tags:
  - bzl_rules
  - hugo
---

`al_hugo_binary` creates a runnable command for a Hugo site.
`al_hugo_run_binary` builds a site with the registered Hugo toolchain and takes
the site as a normal target dependency. Keeping the site separate from the
execution-platform tool prevents site sources from being rebuilt in the exec
configuration.
