---
title: Rules Hugo
description: Bazel rules for Hugo sites
statuses:
  - active
languages:
  - bzl
tags:
  - bzl_rules
  - hugo
---

`rules_hugo` builds Hugo sites with a registered Hugo toolchain. It keeps the
site archive in the target configuration and the Hugo binary in the execution
configuration so site sources are not rebuilt for the execution platform.

`al_hugo_site` wraps a site source archive and its PostCSS tooling.
`al_hugo_run_binary` builds the site with the registered Hugo toolchain and
`al_hugo_binary` exposes a runnable Hugo command for a site. The optional
`al_hugo_worker` rule runs the build through a persistent worker.

## Module setup

```starlark
bazel_dep(name = "rules_hugo", version = "<VERSION>")
```

Register the Hugo toolchain in the root module:

```starlark
register_toolchains("@rules_hugo//main/bzl:all")
```
