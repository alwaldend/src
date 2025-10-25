---
title: Repo
description: Repository-specific information
---

## Setup

- Install packages:
  ```sh
  sudo dnf install \
    clang clang-format java-latest-openjdk-devel rust mesa-libGL-devel go \
    Xvfb python3 gmp-devel
  ```
- Install bazel: https://bazel.build/install/bazelisk
- Install nvm: https://github.com/nvm-sh/nvm
- Install node: `nvm install node`
- Install commandline tools to `${ANDROID_HOME}/cmdline-tools/latest`: https://developer.android.com/studio#command-tools
- Install NDK to `${ANDROID_NDK_HOME}` (`~/Android/Ndk`): https://github.com/android/ndk/wiki
- Install android tools: `sdkmanager "platforms;android-36" "build-tools;36.0.0"`
- Install qt: `bazel run //tools/qt:install`
- Install git hooks: `bazel run //tools/git_hooks:install`

## Development

- Build `compile_commands.json`: `bazel run :refresh_compile_commands`
- Run builds: `bazel build //...`
- Run tests: `bazel test //...`
- Replace a rule with another one:
  ```sh
  find \
    "(" -name "*.bazel" -o -name "*.bzl" -o -name ".bazelrc" -o -name "*.md" ")" \
    -type f \
    -exec sed -i 's|//bzl/rules/txt|//tools/txt|g' "{}" ";"`
  ```
- Remove a rule:
  ```sh
  find -name "BUILD.bazel" -type f -exec sed -i '/al_readme(/,/)/d' "{}" ";"
  ```
