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
    go Xvfb python3 gmp-devel patch podman
  ```
- Install bazel: https://bazel.build/install/bazelisk
- Install nvm: https://github.com/nvm-sh/nvm
- Install node: `nvm install node`
- Install commandline tools to `${ANDROID_HOME}/cmdline-tools/latest`: https://developer.android.com/studio#command-tools
- Install NDK to `${ANDROID_NDK_HOME}` (`~/Android/Ndk`): https://github.com/android/ndk/wiki
- Install android tools: `sdkmanager "platforms;android-36" "build-tools;36.0.0"`
- Install qt: `bazel run //tools/qt:install`
- Install git hooks: `bazel run //tools/git_hooks:install`
- Setup .bzlenv: `. "$(bazel run //tools/bzlenv)"`

## .env file

- `RANCHER_ALWALDEND_COM_K3S_TOKEN`: k3s token for `//infra/charts/rancher.dc1.alwaldend.com`
- `DNSCONTROL_CLOUDFLAREAPI_ACCOUNT_ID`: Cloudflare account id for `//infra/dns`
- `DNSCONTROL_CLOUDFLAREAPI_API_TOKEN`: Cloudflare token for `//infra/dns`

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

## Export .env (bzlenv should export it automatically)

```sh
export $(cat .env | xargs)
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
