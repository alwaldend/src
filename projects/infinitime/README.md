---
title: Infinitime
description: InfiniTime firmware fork with a text watchface and Pomodoro app
statuses:
  - maintenance
languages:
  - cpp
tags:
  - embedded
  - fork
---

This InfiniTime firmware fork adds a text watchface and a Pomodoro app. The
repository packages the pinned firmware source and its build dependencies.

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
during module resolution. The Python tooling manifest and native `uv.lock` live
in this project; npm and SDK archives retain their upstream pins. Python and
Node runtimes use the root workspace toolchains.
The standard pip extension reads the local `uv.lock` without fetching firmware.
The project re-exports npm with a separate extension identity so unrelated npm
targets do not fetch firmware to read its upstream lockfile.
`bazel mod deps` deliberately evaluates all extensions and can still fetch it.

`native_uv.patch` removes the upstream source's pip-tools lock target. This
project owns the effective Python dependencies in `pyproject.toml`; update
them with `//projects/infinitime:uv_lock.update` using the shared
[uv workflow](../../tools/py/README.md). Explicit lock generation uses sandboxed
`uv lock` and requires network access to PyPI. Firmware builds consume the
checked-in native lock without running resolution.

Validate the firmware with `bazel test //projects/infinitime:build_test`.

The Python environment excludes `ecdsa`, whose optional key-signing support in
`adafruit-nrfutil` is unused by this project's unsigned `dfu genpkg` commands.
This removes [CVE-2024-23342](https://github.com/advisories/GHSA-wj6h-64fc-37mp)
from the dependency graph. The helper does not support key generation or signed
DFU packages in this environment. Firmware image and DFU packaging continue
through the existing build workflow.

- Install [Gadgetbridge](https://gadgetbridge.org/)
- Download the firmware:
  ```sh
  oras pull docker.io/alwaldend/src:projects_infinitime_pinetime_mcuboot_app_dfu_1_15_0_zip_head
  ```
- Pair the watch in `Gadgetbridge`
- Install the firmware
