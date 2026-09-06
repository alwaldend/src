---
title: XCP-ng
description: XCP-ng infrastructure
tags:
  - terraform
  - xcp-ng
---

Terraform in `tf` manages a Xen Orchestra resource set for every
entity in Vault's `approles` group, following `infra/pve/tf`. Resource sets
are the XO grouping and delegation mechanism; XCP-ng host pools represent
physical clusters and are not created per AppRole.

`resource_set_inventory` assigns a pool, template, storage repository and
network by name to each AppRole. Native provider lookups resolve their IDs
and reject ambiguous matches. Only Forgejo has an assignment by default;
other sets remain empty. Subjects are the exact synchronized OIDC users
selected by Vault issuer and immutable AppRole entity UUID. Bootstrap has
synchronized the 18 AppRole users; new AppRoles need their first OIDC login
before Terraform can bind their sets. `resource_set_cpu_limit` defaults to
32 CPUs per set. Creating a set does not allocate CPUs or storage.

The provider is pinned to `vatesfr/xenorchestra` 0.41.0. The single
`infra/xcp_ng/al.lua` configuration supplies administrator authentication for
XCP-ng infrastructure operations: its named `xcp_ng` Vault AppRole injects
`XOA_TOKEN` from key `xoa_token` in
`secrets/alwaldend.com/vault1/approles/src_infra_xcp_ng/xoa`.
Forgejo instead uses its own Vault identity through the
[XO login plugin](cmd/xo_login/README.md), which issues an invocation-scoped
XO session and revokes it on shutdown. It retains its own Terraform backend
and does not inherit the infrastructure administrator token. Token values
are not Terraform inputs or outputs.

The default endpoint is `wss://xoa.xcp-ng.alwaldend.com`. DNS maps the XO
appliance to `192.168.1.206` and the XCP-ng host to `192.168.1.213`.
TLS verification is enabled. `xoa_insecure` is an explicit temporary
bootstrap override for an appliance whose certificate has not been replaced.
The HTTP backend stores state under the AppRole's `tf_backend` Vault path.

## OIDC and certificates

`infra/vault/tf` owns the `src_infra_xcp_ng_provider` OIDC client, the
`src_infra_xcp_ng_users` and `src_infra_xcp_ng_admins` groups, and the
`src_infra_xcp_ng_pki_server` ACME role. Apply those scoped prerequisites
before the XO configuration. The callback is
`https://xoa.xcp-ng.alwaldend.com/signin/oidc/callback`.

XO must have the `auth-oidc` plugin installed and trust Vault's CA. The
Terraform provider has no plugin resource, so `terraform_data.oidc` invokes
the packaged Go JSON-RPC helper directly, without a shell. It configures
Vault discovery, the client, and the `user groups` scopes, enables autoload,
and verifies the plugin is loaded. The OIDC client secret remains sensitive
and is stored in the Vault-protected Terraform state. The provisioner signs
the XCP-ng AppRole in through Vault to create the synchronized OIDC groups.
Plugin reconciliation runs when configuration or the helper binary changes,
without native drift refresh. Use `-replace=terraform_data.oidc` to reapply
after manual plugin drift.

Native `xenorchestra_acl` resources manage the administrators group's pool
permissions. An import block adopts matching existing ACLs using their
discovered API IDs. The helper only reads these ACLs; it does not change them.
The users group grants login eligibility, with membership refreshed at login.
Per-AppRole resource-set membership and existing VM ownership use direct user
subjects. The Forgejo user receives administration permission on its own VM;
no shared pool-wide view grant is configured.

The helper remains for gaps in the pinned provider: OIDC plugin configuration
and bootstrap, discovery of external identity subjects and synchronized
groups, host registration, prepared VHD template import, and diagnostics.
Resources, ACLs and named infrastructure lookups use the native provider.

The resource-set management page is
[Self Service](https://xoa.xcp-ng.alwaldend.com/#/self). In the deployed XO
source revision `961b505cfb74cbf24aaeb8d61c5613f97e78de2c`, its
[menu entry](https://github.com/vatesfr/xen-orchestra/blob/961b505cfb74cbf24aaeb8d61c5613f97e78de2c/packages/xo-web/src/xo-app/menu/index.js#L298)
requires a global administrator. The
[page](https://github.com/vatesfr/xen-orchestra/blob/961b505cfb74cbf24aaeb8d61c5613f97e78de2c/packages/xo-web/src/xo-app/self/index.js#L697)
also requires `XOA_PLAN > 3`, corresponding to
[Premium or Community](https://github.com/vatesfr/xen-orchestra/blob/961b505cfb74cbf24aaeb8d61c5613f97e78de2c/packages/xo-web/src/common/xoa-plans.js).
The runtime edition has not been verified. Pool-level administrator ACLs do
not satisfy the page's global-administrator requirement.

The [certificate playbook](ansible/README.md) uses Certbot with Vault EAB,
installs XO's Vault CA trust, and configures renewal. It requires appliance
SSH access and inspected HTTPS certificate paths. Initial issuance and
renewal briefly interrupt XO management access, not running guest VMs.

## Commands

The XO server must be connected to an XCP-ng host before discovering
templates or creating VMs. `xo_inspect` prints selected operational metadata
without credentials or plugin configuration:

```sh
XOA_URL=wss://xoa.xcp-ng.alwaldend.com bazel_agent bazel run //infra/xcp_ng:xo_inspect
bazel_agent bazel run //infra/xcp_ng/tf:tf.plan
bazel_agent bazel run //infra/xcp_ng/tf:tf.apply
bazel_agent bazel test //infra/xcp_ng/tf:tf_tests.fmt_test //infra/xcp_ng/cmd/xo_config:xo_config_test
```

Formatting and helper tests run without live authentication. See
[Forgejo provisioning](../forgejo/tf_setup/README.md) for VM prerequisites.

To register `host1`, store its `xcpng_username` and `xcpng_password` fields
in `secrets/alwaldend.com/vault1/approles/src_infra_xcp_ng/host1`, then run:

```sh
XOA_URL=wss://xoa.xcp-ng.alwaldend.com XCPNG_HOST=192.168.1.213 bazel_agent bazel run //infra/xcp_ng:xo_register_host
```

The command uses the XOA token to register the host and supplies its separate
host credentials through AL. It waits for a connected pool and leaves an
already connected host unchanged. Duplicate registrations are rejected.
`XCPNG_INSECURE=true` explicitly allows the host's bootstrap certificate;
`XOA_INSECURE=true` separately controls verification of XO's certificate.

`//infra/xcp_ng:xo_import_template` imports the repository-pinned image from
`//third_party/org_fedora_cloud` through XO's disk-upload API. Supply
`XO_TEMPLATE_SHA256` from the image's pinned checksum, `XO_TEMPLATE_SR`,
`XO_BASE_TEMPLATE_ID` (a diskless Generic Linux BIOS template), and
`XO_TEMPLATE_NAME`, together with `XOA_URL`. The host must support QCOW2
imports. The command verifies the image before creating anything, creates an
unbooted template with only a boot disk at device 0, and verifies the final
disk layout. It rejects incomplete or conflicting prior imports for inspection.
This keeps cloud-init fresh for cloned VMs and prevents empty CD drives from
shifting the Forgejo data disks to unexpected device names.

Older XO releases, including 5.192.1, require VHD rather than QCOW2 uploads.
For those releases, convert the pinned image to dynamic VHD in task-private
scratch, verify the decoded contents against the original with `qemu-img
compare`, and record the converter version and converted file's SHA256.
Use `//infra/xcp_ng:xo_import_prepared_template` with `XO_TEMPLATE_IMAGE`
set to the absolute converted path, `XO_TEMPLATE_FORMAT=vhd`, and
`XO_TEMPLATE_SHA256` set to the converted checksum. The other inputs are
the same as for the pinned-image importer.

An interrupted import can resume on an inspected empty, halted task VM by
setting `XO_TEMPLATE_RESUME_ID`, `XO_TEMPLATE_PREVIOUS_FORMAT`, and
`XO_TEMPLATE_PREVIOUS_SHA256` to its exact prior identity and image marker.
Resume preserves that VM and refuses any attached disk or network interface.
