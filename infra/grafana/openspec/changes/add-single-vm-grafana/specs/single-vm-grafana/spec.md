## ADDED Requirements

### Requirement: Single Grafana guest

Terraform SHALL provision one persistent Fedora VM through Xen Orchestra.
Ansible SHALL install Grafana and enable it under systemd.

#### Scenario: Fresh installation

- **WHEN** operators apply the VM root and run the deployment playbook
- **THEN** Grafana runs on the guest and stores its SQLite database on the VM disk

### Requirement: Administrative access

Traefik SHALL serve Grafana at the DC1 DNS name over Vault-issued HTTPS and
redirect HTTP. Grafana SHALL listen on loopback. The firewall SHALL allow SSH
from any routed source.

The Grafana process SHALL be allowed to reach the Vault token and UserInfo
endpoints while SELinux is enforcing.

#### Scenario: Browser sign-in

- **WHEN** a Vault global administrator chooses the Vault sign-in option
- **THEN** Grafana uses the registered HTTPS callback and grants server administrator access

#### Scenario: User outside the OIDC assignment

- **WHEN** a Vault user outside the global administrators attempts OIDC sign-in
- **THEN** the application denies access

### Requirement: Collect proxy metrics

Traefik SHALL expose its Prometheus endpoint on the guest address. The guest
firewall SHALL allow the configured metrics port from the VictoriaMetrics collector.

#### Scenario: Scrape the Grafana proxy

- **WHEN** VictoriaMetrics scrapes the Grafana proxy from its collector address
- **THEN** Traefik returns Prometheus metrics

### Requirement: Ingress backend reachability

Grafana SHALL publish an ingress alias and a dedicated backend DNS name.
The guest SHALL admit shared ingress HTTPS from the ingress peer address.

#### Scenario: Route a public Grafana request

- **WHEN** shared ingress connects to the dedicated backend hostname
- **THEN** the guest serves HTTPS for the existing Grafana service hostname
