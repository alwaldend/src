---
title: Flux workloads
---

Terraform installs the Flux Operator using the mirrored Forgejo chart and
image, then applies the FluxInstance after the operator is ready. The instance
uses the distribution embedded in the operator and pulls controller images
from Forgejo. The chart and images support anonymous reads.

The workload root also installs cert-manager with its CRDs and Traefik with
a non-root loopback HTTPS listener on the VM. Both pull their pinned charts
and images from Forgejo. Publish the [cert-manager](../../../third_party/io_cert_manager_cert_manager/README.md)
and [Traefik](../../../third_party/io_traefik_traefik/README.md) mirrors before
applying the root.

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

## Operator ingress

Open `https://operator.flux.alwaldend.com` and sign in with Vault OIDC.
On the DC1 management and workload LANs, DNS points directly to the Flux VM
and browsers need to trust the Vault server CA. Outside those LANs, use
[shared ingress](../../ingress/README.md) with a trusted client certificate.
The Flux user group receives the Kubernetes `view` role.

For the initial rollout, apply the [Flux AppRole root](../approle/README.md)
before running the [K3s playbook](../ansible/README.md). Then apply the workload
root. cert-manager uses a short-lived Flux service account token to renew the
backend certificate through Vault, which validates the token over the private
network. The workload plan needs the Vault OIDC client and the public API CA,
both supplied by AL. After the workload apply, activate the
[host HTTPS socket](../ansible/README.md#https-socket). When migrating the
existing listener, TCP 443 is briefly unavailable between those stages.

Before applying [shared ingress](../../ingress/README.md), apply the
[Flux DNS root](../dns/README.md) and verify that the `flux-operator` issuer and
`flux-operator-tls` certificate are ready in `flux-system`. Shared ingress
connects to `dc1.operator.flux.alwaldend.com` over WireGuard, with TLS
verification for the public hostname.
