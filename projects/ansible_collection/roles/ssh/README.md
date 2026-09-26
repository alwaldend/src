---
title: SSH
description: Setup ssh
tags:
  - ansible_role
---

Host certificate paths are sorted before rendering the hardened SSH
configuration to keep filesystem enumeration order from restarting sshd.

`ssh_client_ca_local_path` selects the controller-side client CA file; its
default remains the packaged `ssh_ca_clients.pub`.
