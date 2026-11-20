## 1. Infrastructure and service

- [x] 1.1 Define the XCP-ng guest and its Vault and Xen Orchestra integration.
- [x] 1.2 Package the Ansible deployment and document VMUI access.

## 2. Verification

- [x] 2.1 Validate the source and representative service behavior offline.

Terraform format and validation, Ansible syntax, and strict OpenSpec
validation passed locally. A checksum-verified binary ran with the service
flags and a disposable data directory. The health check, VMUI, and a
self-scraped metric query passed before and after a restart.

## 3. Ingress, DNS and scraping

- [x] 3.1 Define the DNS root, authenticated ingress route, and private HTTPS listener.
- [x] 3.2 Configure the VictoriaMetrics loopback scrape job.
- [x] 3.3 Validate the DNS, packaged playbooks, and listeners offline.

Offline checks covered both DNS views, the Terraform module schemas, and
the rendered Ansible configurations. The published binary served VMUI over
verified HTTPS and reported the configured loopback scrape target as up
with a disposable local certificate and data directory.
