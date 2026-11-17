---
title: Development host deployment
description: One deployment for both development hosts
---

```sh
bazel run //infra/dev/ansible
bazel run //infra/dev/ansible -- --limit secure1.dev.alwaldend.com
```

One invocation runs both host-specific plays using the `src_infra_dev` AppRole.
Workstation storage and services apply only to dev1; secure1 runs the Forgejo
worker and its build prerequisites.

Enable or reconcile only the secure worker and its build prerequisites with:

```sh
bazel run //infra/dev/ansible -- --limit secure1.dev.alwaldend.com --tags forgejo_runner
```

The deployment identity needs its exact-path read grant for the existing worker
record before initial deployment. The runner role preserves existing local
credentials. Proxy configuration additionally needs an exact read of
`secrets/alwaldend.com/vault1/approles/src_infra_threexui/http_proxies` by
`src_infra_dev`; the controller injects only the selected `t3code` password.
Apply the checked-in Vault policy through its owning workflow before deploying
the worker. Proxy credentials remain in the mode-0640 runner configuration;
the systemd transport units contain no credentials. Inspect the service with the
packaged SSH target:

```sh
bazel run //infra/dev:ssh.secure1 -- 'systemctl is-active forgejo-runner.service'
```

The workstation's managed known_hosts trusts secure1 certificates through the
canonical repository server CA and pins dev1's public host key. The dev1 key
was verified from its local host filesystem and strict SSH handshake during
the deployment; host-key checking remains enabled through hostname cutover.

Build `:ansible_bin` and run `:syntax` for offline validation.

Terraform publishes the cache device and capacity to the development AppRole's
`runner_cache` Vault record. The normal wrapper reads it into
`DEV_RUNNER_CACHE_METADATA`; secure1 group variables decode its `device` and
`size` fields. Apply VM provisioning before deployment so this record exists.
No output command or caller-supplied cache arguments are needed. The cache tasks
reject missing/invalid metadata, boot
positions, capacity mismatches, partitions, incompatible filesystems and mounts
used elsewhere before creating ext4 without forcing an overwrite. The filesystem
is mounted by UUID at `/var/cache/bazel`; the worker requires it before startup.
Ansible has no XO lookup executable, API request or authentication.

The action cache uses its private `disk_cache` directory with 700-GiB garbage
collection. Bazel output bases remain local. The guest receives no SMB credential
or network-filesystem client. Apply storage and the VM disk in the order in the
[provisioning guide](../tf_setup/README.md). Verify
live encryption, mount identity, write access and a repeated Bazel action only
after separately authorized deployment.

To install or reconcile only the pinned standalone fj binary on the current
dev1 host, run `bazel run //infra/dev/ansible:fj_install`. It uses the same
packaged playbook's `fj` tasks with a local connection and privilege escalation
disabled, independently of the Vault/SSH deployment wrapper. `-- --list-tasks`
and `-- --syntax-check` inspect this selection without installing it.
