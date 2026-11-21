---
title: Grafana role
---

Run `alwaldend.main.grafana` with privilege escalation on Fedora. The role
installs the Fedora Grafana package, stores data in the RPM's SQLite
directory, listens on loopback, and checks database health after starting
the systemd service. Set `grafana_hostname` to the HTTPS hostname served by
the reverse proxy.

Supply the local recovery password, encryption key, and OIDC client secret
as files on the Ansible controller. Set the OIDC client ID, browser
authorization URL, token URL, UserInfo URL, and administrator group.
Members of that OIDC group receive Grafana server administrator access.
The role enables PKCE and keeps the local login available for recovery.

Set `grafana_feature_toggles` to enable deployment-specific Grafana features.
