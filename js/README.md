---
title: Js
description: Javascript tree
cascade:
  - categories:
      - js
---

## Links

- Website: https://www.ecma-international.org/publications-and-standards/standards/ecma-262/
- Rules: https://github.com/bazel-contrib/rules_nodejs
- Rules: https://github.com/aspect-build/rules_js

## Commands

- Install:
  ```sh
  bazel run -- @com_alwaldend_src_pnpm//:pnpm --dir "${PWD}" install --lockfile-only
  ```
