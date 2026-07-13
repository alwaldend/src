---
title: harbor
description: Harbor deployment
tags:
  - harbor
  - ansible
  - pve
---

## Links

- Site: https://goharbor.io/
- Releases: https://github.com/goharbor/harbor/releases

## CE Login

- Go to User Profile, copy CLI secret
- Run:
  ```sh
  podman login harbor.alwaldend.com
  ```

## Update secrets

- Generate approle secret id:
  ```sh
  SECRET_ID=$(bazel run //infra/harbor/cl:vault.secret_id | jq -r .data.secret_id)
  echo "Secret id: ${SECRET_ID}"
  ```
- Patch secret-id in [./cl/cert-manager/issuer-approle.yaml](./cl/cert-manager/issuer-approle.yaml)
- Encrypt:
  ```sh
  bazel run //infra/flux/cl:sops.encrypt infra/harbor/cl/cert-manager/issuer-approle.yaml
  ```
