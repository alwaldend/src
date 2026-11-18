---
title: Vault server certificate
---

Use `alwaldend.main.vault_server_certificate` to obtain a certificate from
Vault PKI for a managed host's existing private key. The role retains a trusted
certificate for the configured name while it has at least seven days remaining.
When renewal is needed, the role sends a signing request to Vault from the
Ansible controller, verifies the returned certificate, and installs it on the
host. The private key stays on the managed host.

Supply the Vault PKI signing path, certificate path, private key path, and
service group. The role uses the play's `date_time` fact, the host's CA trust,
and the Vault connection and token injected on the controller. Set
`vault_server_certificate_notify` to a handler when the service needs to reload
or restart after a certificate change.
