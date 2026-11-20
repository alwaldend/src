## ADDED Requirements

### Requirement: Single guest provisioning

The infrastructure configuration SHALL define one Fedora guest with a
persistent disk and a dedicated Xen Orchestra resource set.

#### Scenario: Provision an authorized deployment

- **WHEN** Terraform is applied with the address declared in the component DNS configuration and the resource set is ready
- **THEN** one VM is provisioned with that address on the guest network.

### Requirement: Metrics service and dashboard

The guest playbook SHALL install a checksum-verified VictoriaMetrics release,
start it under systemd, and store metrics in a persistent directory. The
service SHALL listen on loopback for local access and SSH forwarding.

#### Scenario: Explore metrics from the deployed service

- **WHEN** an operator forwards the guest listener over SSH after deployment
- **THEN** VMUI responds and the operator can query self-scraped metrics.

#### Scenario: Restart the service

- **WHEN** systemd restarts VictoriaMetrics
- **THEN** the service continues using the same stored metrics.

### Requirement: Authenticated ingress and scrape configuration

The deployment SHALL configure the loopback scrape job and expose the
dashboard through shared ingress requiring client certificates. The private
HTTPS listener SHALL obtain and renew its service-hostname certificate
through Vault ACME and accept traffic from ingress peers and the approved DC1
LANs. The upstream address and public ingress alias SHALL come from the
component's DNS root.

#### Scenario: Open the dashboard through shared ingress

- **WHEN** an operator with a trusted ingress client certificate opens VMUI
- **THEN** ingress verifies the guest HTTPS backend and forwards the request.

#### Scenario: Resolve the service and the ingress upstream

- **WHEN** clients and shared ingress resolve the dashboard names
- **THEN** the service hostname has a public CNAME to shared ingress and a DC1 A record for the guest
- **AND** the separate `dc1.` upstream hostname has a public A record for the guest's private address.

### Requirement: Shared web login

Terraform SHALL store a generated login in Vault for members of the
administrator group. The deployment SHALL configure Traefik with a
separately stored password verifier.

#### Scenario: Authenticate to the dashboard

- **WHEN** an operator who has passed ingress client authentication supplies the generated login
- **THEN** guest Traefik permits the dashboard request.

#### Scenario: Open the dashboard directly from the LAN

- **WHEN** a DC1 LAN operator opens the service hostname over HTTP using DC1 DNS
- **THEN** Traefik redirects to HTTPS, presents a certificate for the service hostname, and requires the shared login.
