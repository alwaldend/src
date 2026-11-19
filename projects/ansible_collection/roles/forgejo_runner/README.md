---
title: Forgejo runner
description: Install and register a Forgejo Actions runner
tags:
  - ansible_role
---

Installs the repository-pinned Linux x86-64 Forgejo runner, creates its service
account, registers it with Forgejo, and manages its systemd service. The binary
and configuration are root-owned. The home and work directories belong to the
runner account. Registration and the daemon use the same explicit configuration
and registration file, so repeated deployments reuse the existing runner.

Supply `forgejo_runner_instance` and `forgejo_runner_token` for initial
registration. `forgejo_runner_name` defaults to the inventory hostname. Inject
the token from Vault; the role passes it through standard input with logging
suppressed.
The registration file is mode `0600`. Changing labels or capacity updates the
daemon configuration; changing a registered runner's instance or identity
requires deliberate re-registration.

Set `forgejo_runner_install_packages: false` when the calling deployment supplies
all [runner packages](defaults/main.yaml), including on RPM-OSTree hosts.
Node.js is required for JavaScript actions such as checkout.
The `acl` package must be available for Ansible to execute registration as the
unprivileged runner account; both roles include it in their package defaults.
Complete any required reboot before starting this role. Container runtime
configuration remains the caller's responsibility. The default installation
and service paths are writable on OSTree hosts.

Host executor labels run jobs as the runner account without job isolation.
Restrict registration to trusted repositories and keep deployment credentials
out of the runner environment. The role does not grant sudo access. The Actions
cache service is disabled by default; enable it only after configuring its
network reachability. Configuration, binary, and service changes trigger a
handler to reload systemd and restart the runner.
