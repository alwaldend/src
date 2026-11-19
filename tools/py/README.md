---
title: Py
description: Python rules
languages:
  - py
  - bzl
tags:
  - bzl_rules
---

The repository Python dependency lock remains in root `requirements.txt`, next
to `pyproject.toml`. Update it with `//tools/py:requirements.update` and validate
it with `//tools/py:requirements.test`; the root labels remain compatibility
entry points. These upstream pip-tools workflows require network access to
resolve and validate package versions and run outside the action sandbox.
