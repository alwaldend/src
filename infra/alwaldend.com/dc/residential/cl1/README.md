---
title: Cl1
description: cl1.residential.dc.alwaldend.com
---

## Links

- Argocd install source: https://github.com/argoproj/argo-cd/blob/76162a93105a5cf01a10a3faa8203621b487a0c3/manifests/ha/namespace-install.yaml
- Argocd kustomization example: https://github.com/argoproj/argoproj-deployments/blob/master/argocd/kustomization.yaml

## Deployment

- Setup hosts:
  ```sh
  bazel run //infra/alwaldend.com/dc/residential/cl1/ansible:deploy
  ```
- Deploy the chart:
  ```sh
  ssh -NL 127.0.0.1:6443:127.0.0.1:6443 ansible@mini1.residential.dc.alwaldend.com &
  bazel run //infra/alwaldend.com/dc/residential/cl1/helm:deploy
  ```
