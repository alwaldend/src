---
title: Cl1
description: cl1.residential.dc.alwaldend.com
---

## Links

- Argocd install (v3.3.2): https://github.com/argoproj/argo-cd/blob/master/manifests/install.yaml
- Cert-manager install: https://github.com/cert-manager/cert-manager/releases/download/v1.19.4/cert-manager.yaml
- Longhorn: https://github.com/longhorn/longhorn/releases/download/v1.11.0/longhorn.yaml

## Deployment

### Full deployment

```sh
bazel run //infra/alwaldend.com/dc/residential/cl1/ansible:deploy
```

### Only k3s

```sh
bazel run //infra/alwaldend.com/dc/residential/cl1/ansible:deploy -- --tags k3s
```

### Only kustomization

```sh
bazel run //infra/alwaldend.com/dc/residential/cl1/ansible:deploy -- --tags k3s_kustomize
```
