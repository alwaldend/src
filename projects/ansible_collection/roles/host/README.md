---
title: Host
description: Common host setup
tags:
  - ansible_role
---

Package upgrades precede configuration. The firewall is configured after SSH
package dependencies, and OS hardening runs last so their defaults cannot
overwrite the final host policy.
