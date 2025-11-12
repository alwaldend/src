---
title: Bzlenv
description: Setup bazel environment
languages:
  - bzl
tags:
  - repo
  - bzl_rules
---

Rules to create a bazel environment which functions similar to venv

## Features

- Adds bazel-built tools to `${PATH}`
- Exports .env

## Usage

Create env and activate:
```sh
. "$(bazel run //tools/bzlenv)"
```

Activate existing env:
```sh
. .bzlenv/bin/activate
```

Deactivate env:
```sh
bzlenv_deactivate
```
