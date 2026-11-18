---
title: Ansible
description: Ansible setup
tags:
  - ansible
---

The Ansible package installs the [common Forgejo theme](../../../internal/assets/README.md)
through the managed custom header. Forgejo's Light, Dark, and Auto modes use
the website's neutral palette, font stack, and corner radii. Auto follows the
browser preference. The header also carries the mobile pull request diff styles.

Forgejo uses the shared Alwaldend monogram as its SVG favicon and navigation
logo. The `forgejo_branding` tag updates only those public assets.
