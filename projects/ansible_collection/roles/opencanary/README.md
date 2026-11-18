---
title: OpenCanary role
---

Installs OpenCanary in a Python virtual environment on Fedora and runs it as
an unprivileged systemd service. Set `opencanary_services` to a mapping of
module names to port and protocol pairs. Enable the optional audited,
read-only Samba decoy with `opencanary_smb_enabled`.

When HTTPS is enabled, set `opencanary_https_sign_path` to a Vault PKI signing
endpoint to maintain the HTTPS certificate with the
[Vault server certificate role](../vault_server_certificate/README.md).
Use the `opencanary_tls` tag to check or renew it on an existing host.

Events go to the system journal. Only supply an address range the guest is
intended to monitor when permitting inbound access to decoy ports.
