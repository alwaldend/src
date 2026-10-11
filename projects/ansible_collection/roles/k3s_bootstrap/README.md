---
title: K3s bootstrap
description: Bootstrap k3s cluster
tags:
  - ansible_role
---

The role declares Python and uv as host packages, creates its Kubernetes virtual
environment with `uv venv`, and installs the Python client with `uv pip install`.
The environment is created only when its interpreter is missing; package installs
report changes only when uv installs or removes distributions. These operations
run during an explicitly authorized role deployment and require package-index
access. Repository validation builds and renders the role without deploying it.
