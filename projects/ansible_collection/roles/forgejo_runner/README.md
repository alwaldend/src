---
title: Forgejo runner
description: Install and register a Forgejo Actions runner
tags:
  - ansible_role
---

Installs the repository-pinned Linux x86-64 Forgejo runner, creates its service
account, registers it once, and manages `forgejo-runner.service`. Set
`forgejo_runner_dir` to choose the installation directory (default
`/opt/forgejo-runner`). It contains the root-owned `forgejo-runner` executable
and `config.yaml`, and runner-owned `home/` and `work/` directories. The unit
is installed in `/etc/systemd/system`. Registration and the daemon use the
same configuration and `work/.runner` file.

Supply `forgejo_runner_instance` and `forgejo_runner_token` for initial
registration. `forgejo_runner_name` defaults to the inventory hostname. The
deployment supplies the token through its existing Vault injection flow; the
role passes it through stdin with logging suppressed and a bounded registration
deadline. Neither the worker nor this role performs Vault or OIDC login.
Once `work/.runner` exists, deployment needs no registration token or instance
input. That file contains the worker credentials and is restricted to mode
`0600`. Changing labels or capacity updates the
daemon configuration; changing a registered runner's instance or identity
requires deliberate re-registration.

For an identity already created through Forgejo's API, supply
`forgejo_runner_credentials` with `id`, `uuid`, `token`, `instance` and `name`.
The instance and name must match the role's configured values. The role installs
this identity only when `work/.runner` is absent, with secret logging suppressed;
it does not register another worker or overwrite existing credentials.

`forgejo_runner_labels` and `forgejo_runner_capacity` control job selection and
concurrency. `forgejo_runner_config` supplies advanced settings using upstream
YAML keys. Its default mapping is shown in [defaults](defaults/main.yaml).
Replacing it replaces that mapping; unspecified upstream settings use the
runner's defaults. The role always sets `runner.file`, `runner.labels`,
`runner.capacity` and `host.workdir_parent` to match its own inputs.

For example, to change the job timeout and enable the cache:

```yaml
forgejo_runner_config:
  log:
    level: info
    job_level: info
  runner:
    timeout: 4h
    shutdown_timeout: 30s
    insecure: false
  cache:
    enabled: true
```

Set `forgejo_runner_install_packages: false` when the calling deployment supplies
all [runner packages](defaults/main.yaml), including on RPM-OSTree hosts.
Node.js is required for JavaScript actions such as checkout.
The `acl` package must be available for Ansible to execute registration as the
unprivileged runner account; the role includes it in its package defaults.
Complete any required reboot before starting this role. Container runtime
installation and configuration remain the caller's responsibility; Podman is
not installed for the default host executor. The default installation
and service paths are writable on OSTree hosts.

Host executor labels run jobs as the runner account without job isolation.
Restrict registration to trusted repositories and keep deployment credentials
out of the runner environment. The role does not grant sudo access. The Actions
cache service is disabled by default; enable it only after configuring its
network reachability. Configuration, binary, and service changes trigger a
handler to reload systemd and restart the runner.

Callers that install build tools into the runner's home can import `install`
first and `service` after their prerequisites. The normal role entry point runs
both phases. Installation defers restarts; the service phase marks prerequisites
ready and flushes pending handlers. If prerequisite setup fails, forced handlers
leave a fresh worker stopped.

Registration failures produce a fixed message. Check initial inputs,
connectivity, TLS trust and Ansible's ability to switch to the runner account;
keep secret-bearing command output suppressed.

Older path overrides are replaced by `forgejo_runner_dir`, preserving the
default layout. Migrate `forgejo_runner_timeout` and
`forgejo_runner_shutdown_timeout` to `runner.timeout` and
`runner.shutdown_timeout` under `forgejo_runner_config`, and
`forgejo_runner_cache_enabled` to `cache.enabled`. The service stop deadline
remains `forgejo_runner_service_stop_timeout`; keep it longer than the configured
runner shutdown timeout.

The [local acceptance fixture](../../test/forgejo_runner/README.md) exercises
registration and renders both templates without deploying a worker.
