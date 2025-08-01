---
title: Src
description: Source code
cascade:
  - type: docs
    _target:
      path: /**
---

## Links

- Docs: https://alwaldend.github.io/src/
- Dev release: https://github.com/alwaldend/src/releases/tag/dev
- Dev release artifacts: https://github.com/alwaldend/src/releases/download/dev/alwaldend-src.tar.gz

## Licence

AGPLv3, see [LICENSE](./LICENSE.txt)

## Setup

- Install packages: `sudo dnf install clang clang-format java-latest-openjdk-devel rust mesa-libGL-devel go Xvfb`
- Install bazel: https://bazel.build/install/bazelisk
- Install nvm: https://github.com/nvm-sh/nvm
- Install node: `nvm install node`
- Install commandline tools to `${ANDROID_HOME}/cmdline-tools/latest`: https://developer.android.com/studio#command-tools
- Install android tools: `sdkmanager "platforms;android-36" "build-tools;36.0.0"`
- Install qt: `bazel run //starlark/bazel/qt:install`
- Install git hooks: `bazel run //shell/git-hooks:install`

## Development

- Build `compile_commands.json`: `bazel run :refresh_compile_commands`
- Run builds: `bazel build //...`
- Run tests: `bazel test //...`
