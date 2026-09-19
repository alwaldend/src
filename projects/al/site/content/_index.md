---
title: Al
linkTitle: Al
description: Repository command runner and Bazel configuration rules
layout: landing
statuses:
  - active
languages:
  - bzl
  - go
tags:
  - bzl_rules
  - fp
  - proto
---

AL is a command runner that prepares credentials and environment
variables through plugins, runs your command, and cleans up afterward. Its Bazel
rules package commands together with the configuration and plugins they need.

## Features

- Plugin-based credential and environment preparation around any command
- Bazel rules that package a command with its configuration and plugins
- A source-keyed executable cache for running declared tools
