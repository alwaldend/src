---
title: Flux workloads
---

Terraform installs the Flux Operator using the mirrored Forgejo chart and
image, then applies the FluxInstance after the operator is ready. The instance
uses the distribution embedded in the operator and pulls controller images
from Forgejo. The chart and images support anonymous reads. The operator UI
is disabled.

After [bootstrapping K3s](../ansible/README.md), run Terraform from any
host that can reach the VM's SSH port. AL obtains a Vault-signed SSH client
certificate, verifies the VM's SSH host certificate, and opens the API
forward on a private loopback port for the invocation. AL also fetches the
Kubernetes server CA and a Vault OIDC ID token for Helm. The K3s bootstrap
binding grants the Flux AppRole cluster administration through its Vault group.

```sh
bazel run //infra/flux/tf:tf.plan
bazel run //infra/flux/tf:tf.apply
```

Helm waits for the operator deployment. FluxInstance reconciliation continues
after the Helm release completes. Bump the local instance chart version
when changing its template so Terraform plans a Helm upgrade.
