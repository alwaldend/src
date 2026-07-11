---
title: flux
description: Fluxcd deployment
tags:
  - flux
  - ansible
  - pve
  - k3s
---

## Links

- Site: https://fluxcd.io/

## Set up port forwarding

```sh
ssh -L 6443:127.0.0.1:6443 -N flux.alwaldend.com
```

## Run flux CLI

```sh
bazel run //infra/flux/cl:flux
```

## Show flux status

```sh
bazel run infra/flux/cl:flux -- get all -A
```

## Run bootstrap

```sh
bazel run //infra/flux/cl:flux.bootstrap
```

## Run flux operator

```sh
bazel run //infra/flux/cl:op
```

## Show oidc info

```sh
bazel run //infra/flux/cl:oidc
```

## Renew certificate

```sh
bazel run //infra/flux/cl:cmctl -- renew traefik-gateway-websecure-tls -n traefik
```

## Update secrets

- Generate secret id and the token:
  ```sh
  SECRET_ID=$(bazel run //infra/flux/cl:vault.secret_id | jq -r .data.secret_id)
  TOKEN=$(echo "${SECRET_ID}" | bazel run //infra/flux/cl:vault.ops_token | jq -r .auth.client_token)
  echo "Secret id: ${SECRET_ID}, Token: ${TOKEN}"
  ```
- Patch secret-id in [./cl/cert-manager/issuer-approle.yaml](./cl/cert-manager/issuer-approle.yaml)
- Patch token in [./cl/flux-system/flux-sops-secret.yaml](./cl/flux-system/flux-sops-secret.yaml)
- Encrypt:
  ```sh
  bazel run //infra/flux/cl:sops.encrypt infra/flux/cl/cert-manager/issuer-approle.yaml
  bazel run //infra/flux/cl:sops.encrypt infra/flux/cl/sops/flux-sops-secret.yaml
  ```
