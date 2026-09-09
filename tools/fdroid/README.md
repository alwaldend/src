---
title: F-Droid
description: Pinned fdroidserver environment for local F-Droid build testing
statuses:
  - active
languages:
  - python
tags:
  - android
  - fdroid
---

Pinned `fdroidserver` environment for local F-Droid builds. The wrapper
provisions an isolated Python virtualenv and delegates to `fdroid`. Use
`tools/gradle/gradle-wrapper` for the Android app's Gradle build; this tool
only drives F-Droid metadata and local builds.

## Usage

```sh
tools/fdroid/fdroid-wrapper <fdroid-subcommand>
```

Run from an F-Droid data directory containing `config.yml` and `metadata/`.
Caches go under the workspace's ignored `out/` tree.
