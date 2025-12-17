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

- Install [Gadgetbridge](https://gadgetbridge.org/)
- Download the firmware:
  ```sh
  oras pull docker.io/alwaldend/src:projects_infinitime_pinetime_mcuboot_app_dfu_1_15_0_zip_head
  ```
- Pair the watch in `Gadgetbridge`
- Install the firmware
