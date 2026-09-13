---
title: GitHub
description: GitHub organization, repositories, access, and Pages configuration
---

This project manages the GitHub organization and repositories declared by the
[shared repository catalog](../repos/README.md), including existing forks,
project landing sites, organization membership, and developer access.

[Service Terraform](tf/README.md) adopts existing resources and preserves their
identities. It also manages repository default-branch rules and existing Pages
configuration. Site content is published through the owning project deployment
targets; this infrastructure package does not publish it.
