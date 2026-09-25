---
title: Nginx
description: Native Nginx with validated consumer configuration
---

Installs Fedora's `nginx` package and manages its native systemd service.
The consumer supplies a complete `nginx_config_template`, including listeners,
HTTP settings, MIME includes, and document roots. Replacing the main config
also removes package-default virtual hosts from the active configuration.

Ansible renders and validates the entire candidate with `nginx -t -c` before
atomically replacing `nginx_config_path`. Invalid candidates leave the prior
configuration and running service intact. Valid changes notify a reload;
an unchanged run does not restart or reload Nginx. The role never manages
published content. Packages can be caller-managed with
`nginx_manage_packages: false`.

The download component owns its routes, storage permissions, SELinux labels,
and the Traefik TLS frontend. The [role validation record](../../openspec/changes/add-nginx-role/tasks.md)
links the consumer fixture and lifecycle results.
