---
title: OpenCanary role
---

Installs OpenCanary in a Python virtual environment on Fedora and runs it as
an unprivileged systemd service. Set `opencanary_services` to a mapping of
module names to port and protocol pairs. Enable the optional audited,
read-only Samba decoy with `opencanary_smb_enabled`.

Events go to the system journal. Only supply an address range the guest is
intended to monitor when permitting inbound access to decoy ports.
