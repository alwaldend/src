---
title: Forgejo runner
description: Forgejo Actions runner deployment
tags:
  - forgejo
  - ansible
  - pve
---

This project provisions one VM and configures it as a Forgejo Actions runner for
`https://forgejo.alwaldend.com`.

Before running the Ansible target, store a runner registration token at
`alwaldend.com/vault1/approles/src_infra_forgejo_runner/config` under the
`runner_token` key. The deployment registers the runner once; rotating its
registration requires removing `/opt/forgejo-runner/.runner` and rerunning the
playbook.
