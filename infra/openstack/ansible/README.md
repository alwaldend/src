---
title: OpenStack Ansible
description: Host preparation and Kolla-Ansible orchestration
tags:
  - ansible
  - openstack
---

The playbook prepares the two bare-metal hosts, validates the storage and
network layout, renders a Kolla multinode inventory, and invokes the pinned
Kolla-Ansible release from a local virtual environment.

See the parent [OpenStack deployment documentation](../README.md) before
running a target.
