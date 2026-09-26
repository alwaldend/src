---
title: Ansible collection
description: Ansible collection alwaldend.main with Bazel source packaging
statuses:
  - maintenance
tags:
  - ansible_collection
---

`alwaldend.main` is the repository's Ansible collection. Bazel packages its
Ansible source for reuse in repository automation.

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/ansible_collection
- Docs: https://docs.ansible.com/ansible/latest/collections_guide/index.html

## Features

- Ansible collection
- Bazel-managed Ansible source packaging.

## Galaxy

{{< readfile file="galaxy.yml" code="true" lang="yaml" >}}

## Disposable role tests

The initial Molecule scenarios exercise `host`, `traefik`, and `forgejo`
independently on fresh local QEMU guests. Use the individual role targets or the
explicit suite:

```sh
bazel_agent bazel test //projects/ansible_collection/roles/host:molecule_test
bazel_agent bazel test //projects/ansible_collection/roles/traefik:molecule_test
bazel_agent bazel test //projects/ansible_collection/roles/forgejo:molecule_test
bazel_agent bazel test //projects/ansible_collection:molecule_test
```

The [shared runner](../../tools/molecule/README.md) owns supported platforms,
tool and image inputs, acceleration, networking, execution policy, artifacts,
and cleanup recovery. Tests are opt-in and execute afresh when selected.
Guest package services require network access and remain mutable.

Scenarios live under `extensions/molecule`, outside production role packaging.
They test the packaged candidate, including the service binaries, with generated
temporary credentials. Fresh convergence, an unchanged application, configuration
updates, and persistence checks are described in the
[acceptance matrix](extensions/molecule/acceptance.md). Inspect Bazel's
undeclared test outputs for `result.json`, task outcomes, phase summaries, and
guest package versions. Private keys, passwords, and complete runtime trees
are not retained.

The host fixture signs SSH host certificates locally and excludes only the
external `ssh_sign_key` issuance step. Traefik disables Vault EAB/public ACME
and uses a local backend with generated TLS certificates. Forgejo uses SQLite
and a disposable local account and repository; OIDC and external databases are
excluded. These are role tests, not full deployment validation. They never use
a production inventory or change the developer host persistently.
