## ADDED Requirements

### Requirement: Operator web access

Flux SHALL expose its Operator UI through shared ingress using Vault OIDC
and Kubernetes read-only permissions.

#### Scenario: Access the Operator

- **WHEN** a user passes shared ingress client certificate verification and
  signs in through the Flux Vault OIDC provider
- **THEN** the UI impersonates the user with the Kubernetes view role
- **AND** the upstream verifies TLS for the public Operator hostname

#### Scenario: Access the Operator from DC1

- **WHEN** a user on the DC1 management or workload LAN opens the Operator hostname
- **THEN** DC1 DNS directs the browser to Traefik listening on TCP 443 on the Flux VM
- **AND** the user signs in through Vault OIDC

#### Scenario: Renew the backend certificate

- **WHEN** cert-manager reconciles the Operator ingress certificate
- **THEN** it uses a short-lived Kubernetes service account token to
  authenticate to the Vault signing role for the Operator hostname
