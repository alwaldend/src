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
    go Xvfb python3 gmp-devel patch podman git git-lfs fuse fuse-libs pcsc-lite pcsc-lite-devel \
    openssl gnutls-utils opensc openssl-pkcs11 libdnet
  ```
- Install bazelisk: https://bazel.build/install/bazelisk
- Symlink bazel:
  ```sh
  ln -s ~/.local/bin/bazelisk ~/.local/bin/bazel
  ```
- Install git hooks:
  ```sh
  bazel run //:write_git_hooks
  ```
- Install android tools:
  ```sh
  bazel run //tools/android:install
  ```
- Setup dotfiles (optional):
  ```sh
  bazel run //projects/dotfiles:install
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

## Vault login

```sh
bazel run //tools/vault -- login --method=userpass username="${USER}"
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

> This requires around 90G

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

## Adding a bazel module

```sh
mkdir projects/project_name
cd projects/project_name
ln -s ../../.bazelignore
ln -s ../../.bazeliskrc
ln -s ../../tools/bazelrc/main/bazelrc/bzl_project.bazelrc .bazelrc
```

## Ssh key using Yubikey's pgp doesn't work

Agent restart helps:
```sh
gpg-connect-agent updatestartuptty /bye
```

Check the key:
```sh
bazel run //tools/ykman -- info
```

## Remote timeout

If some repo rules are timing out (java ones, for example), you can add
`--remote_timeout 1000000`

## Update pnpm lock

```sh
bazel run //tools/pnpm -- --dir "${PWD}" install
```

Documentation: https://github.com/aspect-build/rules_js/blob/main/docs/pnpm.md#update_pnpm_lock

## Ssh config

Proxmox VMs have Kerberos enabled, which slows SSH down, you need to disable it
locally

`~/.ssh/config.d/config`:
```
GSSAPIAuthentication no
```
