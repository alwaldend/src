---
title: Ansible
description: Ansible setup
tags:
  - ansible
---

The `ansible.time` target applies only the Vault role's time synchronization
tasks. It replaces the main NTP server entry, enables Chrony, restarts it when
configuration changes, and waits for synchronization. Its initial correction
policy can step a drifted clock after restart. It does not restart Vault.
Other Chrony settings and DHCP-provided sources are retained. Verification checks
Chrony's synchronized state and estimated remaining clock correction.

```sh
bazel run //infra/vault/ansible:ansible.time -- --limit host1.vault.alwaldend.com
```

The bare-metal inventory identity remains stable for its signed SSH host
certificate; `ansible_host` selects the hostname declared by Vault's DNS owner.
