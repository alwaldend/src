---
title: Agents
linkTitle: Agents
description: Repository-wide reusable agent skills
layout: landing
statuses:
  - active
languages:
  - markdown
  - bzl
tags:
  - agent
  - skills
---

This project owns the reusable, cross-project agent skills used across
this repository. Each skill is a procedure with its own scope and ownership
rules, so an agent loads the guidance it needs for the task at hand instead of
one large document.

## Features

- One skill per procedure, with an explicit description and scope
- Discovery written into the agent skills directory from declared sources
- Product-specific skills stay with their product; only genuinely shared ones
  live here
