---
title: OS
description: Common os setup
tags:
  - ansible_role
---

On Red Hat systems with UFW installed, the role applies the upstream
`ufw_ipt_sysctl` setting when sysctl and UFW-default management are enabled.
This keeps firewall reloads from overwriting the hardened kernel settings.
