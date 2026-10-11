---
title: Mikrotik
description: Mikrotik setup for dc1.alwaldend.com
languages:
  - rsc
tags:
  - mikrotik
---

The [Terraform package](tf/README.md) prepares owner-local DNS management.
Router exports remain documentation inputs and are not applied by that package.

## DC1 architecture

This diagram describes the maintained architecture source, not live inventory or service health.
The [Mermaid source](docs/dc1.mmd) is authoritative; Bazel renders its light and dark SVGs.

[![DC1 architecture](docs/assets/dc1.svg)](docs/assets/dc1.svg)

## Deployment

- Open [Winbox](https://help.mikrotik.com/docs/spaces/ROS/pages/328129/WinBox)
- Open new terminal
- Run `/export`
- Copy output

## Links

- Website: https://mikrotik.com
- Docs: https://help.mikrotik.com/docs/spaces/ROS/pages/328155/Configuration+Management

## Exports

### Router1

{{< readfile file="router1.rsc" code="true" lang="rsc" >}}

### Router2

{{< readfile file="router2.rsc" code="true" lang="rsc" >}}
