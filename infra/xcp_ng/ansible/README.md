---
title: Xen Orchestra certificates
description: Vault ACME certificate issuance and renewal for the XOA appliance
---

This playbook installs Debian's Certbot package and obtains a certificate for
`xoa.xcp-ng.alwaldend.com` and `host1.xoa.xcp-ng.alwaldend.com` from the Vault
`src_infra_xcp_ng_pki_server` ACME role. The XCP-ng AppRole creates the
external account binding; its registration file is removed after issuance.

The target injects `xoa_ssh_password` from
`secrets/alwaldend.com/vault1/approles/src_infra_xcp_ng/xoa` for the appliance's
`xoa` SSH account and sudo. Ansible uses its `sshpass` password mechanism;
the controller must provide `sshpass` and a verified SSH host key. Passwords
are runtime environment inputs and are not stored in inventory.

Run the read-only inspection target before deploying a different appliance:

```sh
bazel_agent bazel run //infra/xcp_ng/ansible:ansible.inspect
bazel_agent bazel run //infra/xcp_ng/ansible:ansible
```

The checked-in settings match the inspected appliance: `xo-server` runs as
root and listens on HTTPS with `/etc/ssl/cert.pem` and `/etc/ssl/key.pem`.
Native `autoCert` is disabled. The deployment preserves these listener paths.
Apply the Vault PKI role before certificate issuance.

Certbot uses standalone HTTP-01 validation. Vault must resolve both names to
the appliance and reach TCP port 80. Issuance and renewal briefly stop
`xo-server` to free that port, interrupting management access but not guest VMs.
The existing TLS files are preserved with `.before-vault-acme` suffixes before
their paths become symlinks to Certbot's renewable files. The dedicated
`xo-acme-renew.timer` checks twice daily. Certbot is configured to renew the
seven-day Vault certificates with two days remaining; stop/start hooks run
when renewal is needed. Certificate keys and ACME account files remain root-only.

The playbook adds a service-specific Vault CA trust override using
`NODE_EXTRA_CA_CERTS` so XO can also validate Vault's OIDC endpoint. It verifies
the appliance HTTPS endpoint against that CA after deployment.

The deployment target is `//infra/xcp_ng/ansible:ansible`. Building
`//infra/xcp_ng/ansible:ansible_bin` packages the playbook without deploying it.

XO's [native ACME implementation](https://github.com/vatesfr/xen-orchestra/blob/master/%40xen-orchestra/mixins/SslCertificate.mjs)
does not pass EAB credentials to its ACME client. Certbot supports
[external account binding](https://eff-certbot.readthedocs.io/en/stable/using.html#configuration-file)
without changing Vault's mandatory-EAB policy.

## XCP-ng Ansible interpreter

The existing [inventory](inventory.yaml) selects
`/opt/ansible-python/3.14.4/bin/python3.14` for
`host1.xcp-ng.alwaldend.com` at host scope. This takes precedence over the
security consumer's `group_vars/all.yaml` setting of `/usr/bin/python3`.
Shared security workflows consume this owner inventory directly.

The user explicitly authorized a **manual bootstrap exception** for this
host. The parent implementation agent owns installation and validation of
the side-by-side runtime. Bootstrap must preserve vendor Python, `PATH`, and
system symlinks, and must not reboot the host. The runtime follows the
existing rules_python CPython 3.14.4 pin in
[`tools/py/include.MODULE.bazel`](../../../tools/py/include.MODULE.bazel);
there is no additional third-party pin. It is not maintained by `yum`.

The host-scoped inventory setting and this exception document the IaC
reconciliation. The parent reports successful installation and native module
verification. Evidence and pending shared-role audit validation are recorded in the
[owner change](../openspec/changes/export-security-inventory/tasks.md).
The parent must review the candidate before deployment and record the remaining
validation results there.

## XCP-ng host certificates

`//infra/xcp_ng/ansible:ansible.host` provisions a Vault ACME certificate for
`host1.xcp-ng.alwaldend.com` and `xcp-ng.alwaldend.com`. Inspect a different
host with `//infra/xcp_ng/ansible:ansible.inspect_host` before deployment.
The host must serve `/opt/xensource/www/.well-known/acme-challenge` through
HTTP or an HTTP-to-HTTPS redirect; the inspection verifies and cleans up a
public probe file. This host supports that path without changing XAPI settings.

XCP-ng's vendor Python 3.6 cannot run the current Ansible modules. The
certificate playbook explicitly uses `/usr/bin/python3` through Ansible's
script transport for a Python 3.6-compatible helper, without adding packages
to dom0. The helper downloads the immutable Lego release declared in
`third_party/com_github_go_acme_lego_bin/binary_toolchain.json`, verifies its
checksum, and installs it under root-only `/etc/xcp-ng-acme`.

The initial EAB is temporary. Only the ACME account remains for unattended
renewal; no Vault token is retained on the host. The client requests RSA4096
keys to match the Vault role. `xcp-ng-acme-renew.timer` runs twice daily and
uses XAPI's supported `xe host-server-certificate-install` command when the
issued certificate differs from the installed certificate. Installation is
retried independently of issuance. Renewal does not stop XAPI or running VMs.
The final deployment verifies HTTPS and enables the timer.

Set an absolute task-private `TMPDIR` for controller registration scratch.
Ansible removes the temporary registration script after execution. The
`bootstrap` tag updates the helper without ordering a certificate.

## OpenID Connect plugin

The XO playbook installs the official `xo-server-auth-oidc` 0.4.0 plugin with
pinned, checksum-verified upstream files and isolated dependencies. This
version supports Vault group synchronization while remaining compatible with
the installed XO core. Ansible verifies the plugin before restarting XO;
Terraform configures it and reconciles group ACLs through the API.

SSH connection reuse keeps deployment within the appliance's existing UFW
limit of six new SSH connections per 30 seconds. No firewall exception is
required.
