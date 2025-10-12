---
title: Android
description: Android tree
cascade:
  - categories:
      - kt
      - android
---

## Links

- Website: https://www.android.com/intl/en_us
- Rules: https://github.com/bazelbuild/rules_android
- Rules: https://github.com/bazelbuild/rules_kotlin

## Commands

- Pin:
  ```sh
  REPIN=1 bazel run @com_alwaldend_src_maven_android//:pin
  ```
- Check updates:
  ```sh
  bazel run @com_alwaldend_src_maven_android//:outdated
  ```
