---
title: OpenStack Ansible
description: Host preparation and Kolla-Ansible orchestration
tags:
  - ansible
  - openstack
---

The playbook validates both bare-metal hosts, prepares local storage, renders a
Kolla multinode inventory, and invokes the pinned Kolla-Ansible revision.

See the parent [OpenStack deployment documentation](../README.md) before you
run a target. The normal targets can change hosts and cloud services.
