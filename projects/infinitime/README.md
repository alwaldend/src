---
title: Infinitime
description: Fork of InfiniTimeOrg/InfiniTime
statuses:
  - maintenance
languages:
  - cpp
tags:
  - embedded
  - fork
---

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/infinitime
- Fork code: https://github.com/alwaldend/com_github_infinitimeorg_infinitime
- Upstream code: https://github.com/InfiniTimeOrg/InfiniTime

## Features

- Extra Watchface: `Text`
- Extra app: `Pomodoro`
- C++, embedded

## Usage

The root workspace declares the pinned firmware with `use_repo_rule` in
`include.MODULE.bazel`, preserving recursive Git submodules and the CMake
toolchain patch. Firmware source is fetched when referenced, rather than
during module resolution. Its Python/npm lockfiles and SDK archives retain
the upstream pins; Python and Node runtimes use the root workspace toolchains.
The project re-exports pip and npm with separate extension identities so
unrelated Python/npm targets do not fetch firmware to read those lockfiles.
The isolated pip extension declares the Linux x86_64 host platform used by
the pinned ARM compiler archive.
`bazel mod deps` deliberately evaluates all extensions and can still fetch it.

Validate the firmware with `bazel test //projects/infinitime:build_test`.

- Install [Gadgetbridge](https://gadgetbridge.org/)
- Download the firmware:
  ```sh
  oras pull docker.io/alwaldend/src:projects_infinitime_pinetime_mcuboot_app_dfu_1_15_0_zip_head
  ```
- Pair the watch in `Gadgetbridge`
- Install the firmware
