---
title: Maven Install Gradle Converter
description: Converts rules_jvm_external maven_install.json locks to Gradle build inputs
statuses:
  - active
languages:
  - go
tags:
  - android
  - gradle
  - bazel
---

Converts `rules_jvm_external` `maven_install.json` lock files into Gradle
version catalogs and dependency verification metadata so Bazel-built Android
apps can maintain the parallel Gradle build required for F-Droid.

## Usage

Use the `al_gradle_lock` Starlark macro from
`main/bzl/gradle_lock.bzl` in your app's `gradle/BUILD.bazel`:

```starlark
load("//tools/maven_install_gradle_converter:defs.bzl", "al_gradle_lock")

al_gradle_lock(
    name = "update",
    lock_file = "//projects/<app>:maven_lock.json",
    tool = "//tools/maven_install_gradle_converter/cmd/gradle_lock_gen",
)
```

Run `bazel run :update_toml_write` to regenerate the checked-in version
catalog. The corresponding `*_write_test` target fails when it goes stale.

Gradle owns `verification-metadata.xml` because it must record the plugin
classpath in addition to the Bazel lock. Regenerate it with:

```sh
projects/<app>/../../tools/gradle/gradle-wrapper --project-cache-dir out/gradle-wrapper/project-cache --write-verification-metadata sha256
```
