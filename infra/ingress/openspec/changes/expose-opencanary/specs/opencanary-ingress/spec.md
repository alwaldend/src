## ADDED Requirements

### Requirement: Certificate-authenticated OpenCanary web ingress

The OpenCanary website SHALL inherit the ingress requirement for a verified
client certificate. Traefik SHALL use HTTPS to the web decoy and verify the
backend certificate against the configured server name.

#### Scenario: Reject a client without a certificate

- **WHEN** a client connects to the public HTTPS hostname without a trusted client certificate
- **THEN** ingress rejects the connection before forwarding it to OpenCanary

#### Scenario: Reject an invalid backend certificate

- **WHEN** the backend certificate is expired, untrusted, or issued for a different name
- **THEN** Traefik rejects the backend TLS connection
