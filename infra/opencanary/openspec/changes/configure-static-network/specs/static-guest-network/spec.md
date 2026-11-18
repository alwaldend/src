## ADDED Requirements

### Requirement: Static networking and site-local DNS

OpenCanary SHALL use a single declared static IPv4 address on the regular Xen Orchestra network.

#### Scenario: Provision the guest and resolve its name

- **WHEN** the VM and its DNS root are applied
- **THEN** the guest IPv4 address matches the site-local A record
- **AND** Ansible can address the guest by its configured name

### Requirement: Authenticated TLS on the private web listener

The OpenCanary guest SHALL serve a Vault-issued HTTPS certificate for its
service name. The guest firewall SHALL allow the ingress peer to reach that
listener. Ansible SHALL replace a certificate that will expire within seven
days and retain the guest-generated private key.

#### Scenario: Configure or renew the certificate

- **WHEN** Ansible finds a missing, untrusted, or soon-to-expire certificate
- **THEN** it signs a request from the guest key and restarts OpenCanary using the issued certificate
