---
title: Download host
description: Btrfs storage, SSH publication, and native static hosting
---

`alwaldend.main.download_host` configures download hosts using the existing
`host`, `nginx`, and `traefik` roles. It checks that the selected content device
is a block device before configuring the host, creates Btrfs without
forcing replacement of an existing filesystem, and mounts by UUID. This role
is intended for the Fedora hosts used by `infra/download`.

The role creates only the publication roots: `projects/`, `sites/`, and private
`staging/`, owned by the SSH publisher. Publishers create individual projects,
sites, and releases. The mount root and service state remain root-owned;
privileged ownership changes never traverse per-site directories. Shared Nginx
and Traefik roles own their service units. Traefik uses its standard system-disk
location. A daily systemd timer runs incremental duperemove over projects
and sites with resource limits and private hash state. Uploading files,
extracting sites, and switching the current-release symlink are separate
publisher operations.

On SELinux-enabled hosts, the role persists `httpd_sys_content_t` labels for
static content and applies them with `restorecon`. This permits confined Nginx
to read the custom document root; ordinary Unix permissions alone do not grant
that access. See [Red Hat's Nginx configuration guide](https://docs.redhat.com/documentation/red_hat_enterprise_linux/10/html/deploying_web_servers_and_reverse_proxies/setting-up-and-configuring-nginx).

Filesystem creation, UUID-based mounting, directories, SELinux labels, and daily
maintenance live in `tasks/filesystem.yaml`. The pinned
[`community.general.btrfs_info`](https://docs.ansible.com/projects/ansible/13/collections/community/general/btrfs_info_module.html)
module discovers filesystem UUIDs; selection uses the canonical device path so
both direct device paths and by-id symlinks work. Attachment checks remain in
`tasks/main.yaml` before base-host setup.

## Caller contract

- Select target hosts through inventory or deployment wrapper limits and enable
  privilege escalation. The role does not restrict play size.
- Set `force_handlers: true` on the play so completed configuration changes run
  their notified handlers even if a later task fails.
- Supply `download_content_device` and `download_publisher_authorized_keys`.
  Obtain keys through the caller's injection mechanism; do not embed credentials
  in role defaults.
- Override `download_root` and `download_publisher` as needed. The content mount
  must be under root-owned, non-publisher-writable parent directories.
- Keep Traefik configuration and data on the system disk (the shared role
  defaults to `/opt/traefik`) so content-disk failure does not stop the proxy.
- Supply the shared roles' configuration, including `nginx_config_template`,
  Traefik static/dynamic templates, ACME settings, and base-host variables.
  The caller owns inventory, firewall policy, Vault references, and routing.

```yaml
- name: Configure the selected download VM
  hosts: download
  become: true
  force_handlers: true
  roles:
    - alwaldend.main.download_host
```

The collection packages the role automatically through its existing role
aggregation. Domain-specific serving templates stay with the deployment.

{{< readfile file="defaults/main.yaml" code="true" lang="yaml" >}}
