---
title: VictoriaMetrics web credential setup
---

Run `bazel run //infra/victoria_metrics/tf_setup/web_auth:tf_setup.apply`
after configuring the component AppRole and before running Ansible. The
component AppRole writes the dashboard login and bcrypt verifier to `web_auth`
in its own Vault subtree. Ansible injects the verifier file for Traefik.

The credential state uses an admin-owned backend because it contains the
generated password. Reapply the guest playbook after rotating the credentials
to load the updated verifier.
