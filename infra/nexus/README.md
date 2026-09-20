---
title: Nexus
description: Planned native Nexus Repository deployment on XCP-ng
---

This component currently contains a proposal and its OpenSpec validation
scaffolding. It contains no service implementation and establishes no deployed
state. See the [deployment proposal](openspec/changes/add-nexus-deployment/proposal.md)
for the agreed scope and [design](openspec/changes/add-nexus-deployment/design.md)
for architecture, dependencies, and remaining deployment inputs.

`infra/nexus` owns the deployment-specific configuration and documentation.
Its planned boundaries are `tf_setup` for the VM and DNS, `ansible` for native
service deployment, and `tf` for Nexus API configuration. Reusable service and
database behavior belongs to `projects/ansible_collection`; Vault identities
and XO grants remain with `infra/vault` and `infra/xcp_ng`. Traefik owns ingress.

The current Bazel targets package documentation and OpenSpec sources only.
The owner is registered in the repository's OpenSpec validation suite.
Implementation and live infrastructure operations require later user requests.
