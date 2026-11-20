---
title: Nas
date: 2026-08-23
description: TrueNAS deploy
tags:
  - nas
  - truenas
---

The [DNS Terraform package](dns/README.md) manages this owner's DNS records through
its dedicated Vault AppRole. The separate [TrueNAS root](truenas/README.md)
manages the existing pool, filesystem datasets, and Nextcloud lifecycle
metadata with its own state. Existing NAS resources are
adopted through that root; appliance installation remains outside these roots.
Namespaces and child datasets are managed independently by AppRole under
[approles/](approles/README.md), with separate state and plan/apply targets.

## Architecture

This diagram describes the maintained architecture source, not live inventory or service health.
The [Mermaid source](docs/architecture.mmd) is authoritative; Bazel renders its light and dark SVGs.

[![Truenas architecture](docs/assets/architecture.svg)](docs/assets/architecture.svg)

## Links

- Docs: https://www.truenas.com/
