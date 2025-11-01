---
title: Ansible
description: Ansible tree
---

## Setup Rancher

```sh
bazel run infra/ansible:bin.ansible_playbook -- -l rancher1 alwaldend.main.rancher1.yaml
```
