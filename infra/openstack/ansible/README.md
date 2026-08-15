---
title: OpenStack Ansible
description: Host preparation for the two-node OpenStack deployment
---

This package contains only host validation and storage preparation.

The Bazel launcher invokes this playbook as one normal Ansible process. It then
invokes Kolla-Ansible directly. An Ansible task never starts another Ansible
process.

See the parent [OpenStack deployment documentation](../README.md).
