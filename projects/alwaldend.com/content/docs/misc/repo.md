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
    clang clang-tools-extra java-latest-openjdk-devel rust cargo mesa-libGL-devel \
    go Xvfb python3 gmp-devel patch podman git git-lfs fuse fuse-libs pcsc-lite pcsc-lite-devel \
    openssl gnutls-utils opensc openssl-pkcs11 libdnet qemu-kvm qemu-img
  ```
- Install bazelisk: https://bazel.build/install/bazelisk
- Symlink bazel:
  ```sh
  ln -s ~/.local/bin/bazelisk ~/.local/bin/bazel
  ```
- Install AL for repository lifecycle commands and release stamping:
  ```sh
  bazel run //tools/al:install
  ```
  Keep `~/.local/bin` on `PATH` and rerun this command after updating AL's
  source. Release builds use `bazel build --config=release <target>` with
  the installed AL command supplying workspace status. See the
  [release workflow](../../../../al/README.md) for version selection and
  publication.
- Install git hooks:
  ```sh
  bazel run //:write_git_hooks
  ```
- Install android tools:
  ```sh
  ANDROID_HOME="${HOME}/Android/Sdk" \
    bazel run \
      --repo_env=ANDROID_HOME= \
      //tools/android:install
  ```
  The empty repository environment is required for the first installation so
  Bazel does not try to load the SDK before `sdkmanager` installs it.
  Afterwards, export `ANDROID_HOME` and `ANDROID_NDK_HOME`, or set both with
  `--repo_env` in the ignored `user.bazelrc`. The repository dotfiles export
  the expected paths.
- Setup dotfiles (optional):
  ```sh
  bazel run //projects/dotfiles:install
  ```
- Install nvm: https://github.com/nvm-sh/nvm
- Install node:
  ```sh
  nvm install node
  ```
- Qt is downloaded by Bazel; no host installation is required.

## Bazel output for agents

Use ordinary `bazel build`, `bazel test`, and `bazel run` commands. Bazel reads
`~/.bazelrc` for user defaults across workspaces. For bounded, non-interactive
output and failure-focused test reporting, add the **Agent output defaults**
block from the
[managed Host Bot configuration](https://github.com/alwaldend/src/blob/master/users/simeonwarren/host_bot/ansible/files/bazelrc)
to your own `~/.bazelrc`. The block uses standard Bazel options and needs no
named `agent` profile. Explicit command-line options override these defaults.

Copy only that output-settings block; the cache directories, SDK locations,
and resource limits elsewhere in the file belong to Host Bot. Ansible manages
the complete file on that host.

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
ln -s ../../tools/bazelrc/bzl_project.bazelrc .bazelrc
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

## Update certificates on Fedora

```sh
sudo cp data/ssl/alwaldend.com/pki_ca_servers.crt /etc/pki/ca-trust/source/anchors/
sudo update-ca-trust
```

## Run go mod tidy

```sh
find -name "go.mod" -exec sh -c 'cd "$(dirname "{}")" && bazel run @rules_go//go -- mod tidy' ";"
```
