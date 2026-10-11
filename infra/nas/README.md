---
title: Nas
description: TrueNAS deploy
tags:
  - nas
  - truenas
---

The [Terraform package](tf/README.md) prepares this owner's DNS records through
its dedicated Vault AppRole. NAS service provisioning remains outside that root.

## Architecture

This diagram describes the maintained architecture source, not live inventory or service health.
The [Mermaid source](docs/architecture.mmd) is authoritative; Bazel renders its light and dark SVGs.

[![Truenas architecture](docs/assets/architecture.svg)](docs/assets/architecture.svg)

## Links

- Docs: https://www.truenas.com/
