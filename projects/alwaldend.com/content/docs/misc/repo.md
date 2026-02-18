---
title: Repo
description: Repository-specific information
tags:
  - repo
---

## Setup

- Install packages:
  ```sh
  sudo dnf install \
    clang clang-format java-latest-openjdk-devel rust cargo mesa-libGL-devel \
    go Xvfb python3 gmp-devel patch podman git git-lfs fuse
  ```
- Install bazelisk: https://bazel.build/install/bazelisk
- Symlink bazel:
  ```sh
  ln -s ~/.local/bin/bazelisk ~/.local/bin/bazel
  ```
- Install git hooks:
  ```sh
  bazel run //tools/git_hooks:install
  ```
- Setup dotfiles (optional):
  ```sh
  bazel run //projects/dotfiles:install
  ```
- Install android tools:
  ```sh
  bazel run //tools/android:install
  ```
- Install nvm: https://github.com/nvm-sh/nvm
- Install node:
  ```sh
  nvm install node
  ```
- Install qt:
  ```sh
  bazel run //tools/qt:install
  ```
- Set up secrets:
  ```sh
  bazel run //tools/secrets:systemd_creds_edit
  ```

## Setup .bzlenv

```sh
. "$(bazel run //tools/bzlenv)"
```

## Compile commands

Build `compile_commands.json`:

```sh
bazel run :refresh_compile_commands
```

## Build everything

```sh
bazel build //...
```

## Test everything

```sh
bazel test //...
```

## Replace a rule with another one

```sh
find \
  "(" -name "*.bazel" -o -name "*.bzl" -o -name ".bazelrc" -o -name "*.md" ")" \
  -type f \
  -exec sed -i 's|//bzl/rules/txt|//tools/txt|g' "{}" ";"`
```

## Parse a tar manifest

```sh
bazel build //projects/alwaldend.com:site_source_archive && \
  cat bazel-bin/projects/alwaldend.com/site_source_archive.manifest | \
  jq -r --slurp ". | flatten | .[].dest" | \
  grep "^content"
```

## Remove a rule

```sh
find -name "BUILD.bazel" -type f -exec sed -i '/al_readme(/,/)/d' "{}" ";"
```

## Find all public targets (ignores `default_visibility`)

```sh
bazel query 'attr(visibility, "//visibility:public", //...)'
```
