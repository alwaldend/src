## Why

The repository needs a simple public home for release files and static
websites, using the same deployment on local XCP-ng and Yandex Cloud. SSH
publication should expose ordinary release files immediately after upload
and activate archived websites through a selected-release symlink.

## What Changes

- Add `infra/download`, with one VM in each environment, 100 GB of filesystem
  storage per VM, and independent deployment targets and Terraform state.
- Run the standard Traefik service in front of Nginx on each VM. Both receive
  traffic directly; `infra/ingress` is not part of the serving path.
- Serve `download.alwaldend.com` and `alwaldend.com` through split-horizon DNS:
  public records select Yandex Cloud and local records select XCP-ng.
- Expose public release files and JSON directory listings. Keep extracted
  website releases separately and serve the selected release through a
  `current` symlink.
- Integrate the linked Nginx role, component Vault identity, SSH release-tool
  deployment, and themed downloads browser in the existing website.
- Provide no backups, automatic replication, failover, or rollback command.
  Redeploying an existing website release changes its selected symlink.

## Capabilities

### New Capabilities

- `static-hosting`: Two directly reachable VM deployments, public release
  storage, selected website serving, split-horizon routing, and persistence.

### Modified Capabilities

None in this owner. Shared component changes retain their own specifications.

## Impact

This change coordinates the following owner-local plans:

| Owner              | Change                                                                                                                                       |
| ------------------ | -------------------------------------------------------------------------------------------------------------------------------------------- |
| Ansible collection | [Nginx role](../../../../../projects/ansible_collection/openspec/changes/add-nginx-role/proposal.md)                                         |
| Vault              | [Download AppRole](../../../../../infra/vault/openspec/changes/add-download-approle/proposal.md)                                             |
| Release tool       | [SSH deployment](../../../../../tools/release/openspec/changes/add-ssh-deployment/proposal.md)                                               |
| Main website       | [Downloads browser and SSH hosting](../../../../../projects/alwaldend.com/openspec/changes/add-download-browser-and-ssh-hosting/proposal.md) |

Infrastructure implementation will also update the existing Yandex folder and
XCP-ng resource-set assignments, and the apex DNS records in their current
`infra/dns` owner. New download records belong to `infra/download`; existing
apex records must not be duplicated. The website plan owns the production
deployment transition from GitHub Pages.

Acceptance requires the linked changes plus an end-to-end run through SSH
publication, Nginx listings, the website browser, site selection, and both DNS
views. This planning delivery does not provision hosts, change DNS, or publish
the website. Live operations require their own explicit environment scope.
