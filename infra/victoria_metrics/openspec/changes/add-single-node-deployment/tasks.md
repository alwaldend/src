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

## 4. Web credential setup

- [x] 4.1 Use the VictoriaMetrics AppRole to provision one web credential entry, grant `global_admins` read access, and verify authenticated VMUI access after Ansible loads the verifier.

## 5. VMUI dashboard

- [x] 5.1 Package a VMUI dashboard and verify that the pinned server serves its panels and executes their queries.

## 6. Architecture diagram

- [x] 6.1 Render the VM boundary within DC1 infrastructure, dashboard access, Traefik scrape ports, and a stack representing the local scrape targets.

## Grafana proxy coverage

- [x] Add the Grafana Traefik endpoint to the DC1 scrape job.
