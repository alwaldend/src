---
title: Mikrotik
description: Mikrotik setup for dc1.alwaldend.com
languages:
  - rsc
tags:
  - mikrotik
---

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
