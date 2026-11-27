---
title: K3s
description: Setup k3s
tags:
  - ansible_role
---

Set `k3s_manifests_local_path`, `k3s_auth_config_local_path`, or
`k3s_admission_config_local_path` to an empty string to omit that file during
bootstrap.

Set `k3s_server_ca_vault_path` to a Vault KV v2 API path such as
`secrets/data/example/kubernetes_ca` to publish the public server CA after
K3s starts. The controller uses the injected Vault connection and credentials.
The role stores `server-ca.crt` as `certificate` and writes a new version
only when the certificate changes. An empty path disables publication.
