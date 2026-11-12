---
title: MCP Cordis
linkTitle: MCP Cordis
description: Workspace-local runtime packages behind a stable MCP server
layout: landing
date: 2026-08-30
# Generated from project README by //tools/repo_quality:inventory_update.
statuses:
  - maintenance
tags:
  - mcp
  - cordis
---

`mcp_cordis` is a standalone stdio MCP server that mounts workspace-local
runtime packages behind one stable server interface, so packages can change
without changing how a client connects.

## Features

- One stable stdio MCP server for workspace-local packages
- Reusable and disposable package definitions in the same layout
- Ordinary ESM files rather than a bespoke plugin format
