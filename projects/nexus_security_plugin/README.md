---
title: Nexus security plugin
description: Security plugin for Sonatype Nexus 3
statuses:
  - finished
languages:
  - java
tags:
  - plugin
  - sonatype_nexus
---

This plugin allows you to perform a check every time an artifact is requested from a repository

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/nexus_security_plugin

## Features

- A request handler that checks all requests to repositories
- Capability controlling the request handler
- A task that periodically updates the capability using a remote source
- [Sonatype Nexus 3](https://www.sonatype.com/products/sonatype-nexus-repository) plugin
- Java

## Arch

![diagram](./diagram.svg)
