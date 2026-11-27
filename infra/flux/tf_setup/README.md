---
title: Flux XCP-ng VM
---

This root provisions the Fedora VM and its disks on XCP-ng through Xen
Orchestra. It reads the host address from [Flux DNS](../dns/dnsconfig.json).

The `src_infra_flux` Vault AppRole signs in to Xen Orchestra through Vault
OIDC. AL injects the invocation-scoped XO token and revokes it on normal
shutdown. Flux can provision within its assigned resource set. The state
backend uses the same AppRole and the dedicated `tf_backend/xcp_ng_setup`
path.

Before the first provisioning run, synchronize the AppRole's XO identity if
needed, then reconcile the `src_infra_flux` resource set using the
[Xen Orchestra workflow](../../xcp_ng/tf/README.md).

```sh
XO_BOOTSTRAP_APPROLE=src_infra_flux bazel run //infra/xcp_ng/cmd/xo_login:bootstrap
bazel run //infra/xcp_ng/tf:tf.plan
bazel run //infra/xcp_ng/tf:tf.apply
```

Review the Flux VM plan before applying it.

```sh
bazel run //infra/flux/tf_setup:tf_setup.plan
bazel run //infra/flux/tf_setup:tf_setup.apply
```
