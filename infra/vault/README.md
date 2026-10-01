---
title: Vault
description: Setup for vault.dc1.alwaldend.com
tags:
  - ansible
  - terraform
---

## Links

- Intermediate CA: https://developer.hashicorp.com/vault/tutorials/pki/pki-engine-external-ca
- ACME: https://developer.hashicorp.com/vault/docs/secrets/pki/acme

## Deployment

```sh
bazel run //infra/vault/tf_setup:tf.apply # Create VMs (requires an active Vault host)
bazel run //infra/vault/ansible # Set up hosts (BM and VMs)
bazel run //infra/vault/tf:tf.apply # Configure vault
```

## Backup

```sh
bazel run //infra/vault:backup
```

## Unseal

With a working Vault:

```sh
bazel run //infra/vault:unseal
```

Without a working Vault:

```sh
bazel run //infra/vault:unseal_standalone
```

## Fix quorum

```sh
bazel run //infra/vault/ansible:fix_quorum
```

## Set up only VMs

```sh
bazel run //infra/vault/ansible:ansible.vm # Set up only VMs
```

## Set up only bare metal

```sh
bazel run //infra/vault/ansible:ansible.bm # Set up only bare metal
```

## Tf

Plan:

```sh
bazel run //infra/vault/tf:tf.plan
```

Apply:

```sh
bazel run //infra/vault/tf:tf.apply
```

Run terraform directly:

```sh
bazel run //infra/vault/tf:tf.direct -- -chdir="${PWD}" plan
```

## Replace VMs

```sh
bazel run //infra/vault/tf:tf.apply -- -replace 'module.vm_ha["host2"].proxmox_vm_qemu.vm' -replace 'module.vm_ha["host3"].proxmox_vm_qemu.vm
```

## Generate and import a user cilent certificate

```sh
username="username"
bazel run //infra/vault:gen_client_cert -- --user "${username}" --output_dir "${PWD}"
bazel run //tools/ykman -- piv certificates import 9A "${PWD}/${username}.pfx"
bazel run //tools/ykman -- piv keys import 9A "${PWD}/${username}.pfx"
```

## Generate a host cilent certificate

```sh
bazel run //infra/vault:gen_client_cert -- --host some-host --output_dir "${HOME}/.al/client_cert"
```

## Generate a user device cilent certificate

```sh
bazel run //infra/vault:gen_client_cert -- --host some-host --user username --output_dir "${HOME}/.al/client_cert"
```

## Unseal

- Prepare encrypted unseal token
- Run and input the encrypted token:
  ```sh
  bazel run //infra/vault:unseal
  ```

## Root token

- Prepare encrypted unseal token
- Run and input the encrypted token:
  ```sh
  bazel run //infra/vault:gen_root_token -- --pgp_key path_to_public_gpg_key_in_base64
  ```

## Generate an EAB for ACME

```sh
bazel run //infra/vault -- write -f pki/ica_servers/roles/ica_servers_dc1_pve1/acme/new-eab
```

## Sign a client ssh key

```sh
bazel run //:vault -- write ssh/clients/sign/admins ttl=30000000  public_key=@"${HOME}/.ssh/key"
```

## Revoke all tokens

```sh
bazel run //infra/vault -- token revoke -mode=path auth
```

## Vault certificates

Renew the Vault server certificate through its Ansible role. The `vault_tls`
tag issues a certificate, installs its chain and private key, and reloads Vault
with SIGHUP. It stops on certificate errors and does not run the host setup,
package installation, or service restart tasks.

Run from the repository root. Keep SSH host verification enabled by deriving a
temporary `known_hosts` entry from the repository-owned public SSH CA:

```sh
vault_recovery_dir="${PWD}/out/vault-certificate-recovery"
mkdir -p "${vault_recovery_dir}"
printf '@cert-authority host1.vault.alwaldend.com,host1.vault.dc1.alwaldend.com %s\n' \
  "$(cat infra/vault/tf/output/ssh_ca_servers.crt)" \
  >"${vault_recovery_dir}/ssh_known_hosts"
```

For an expired server certificate, temporarily disable Vault certificate
verification for this recovery command only. The SSH options trust that CA
only for this invocation and retain the existing SSH certificate's host alias:

```sh
env VAULT_SKIP_VERIFY=true \
  ANSIBLE_SSH_COMMON_ARGS="-o StrictHostKeyChecking=yes \
-o HostKeyAlias=host1.vault.alwaldend.com \
-o UserKnownHostsFile=\"${vault_recovery_dir}/ssh_known_hosts\"" \
  bazel_agent bazel run //infra/vault/ansible:ansible -- \
  --limit host1.vault.alwaldend.com --tags vault_tls
```

Vault must be reachable and unsealed, and the host token must allow the
component's AppRole login. If that token has expired, connect the login YubiKey
and refresh it first:

```sh
bazel_agent bazel run //tools/vault/login -- simeonwarren --insecure
```

This login replaces the host token. Stop if it fails; a missing YubiKey cannot
be fixed by disabling TLS verification. Neither command changes the persistent
TLS verification configuration. Do not add the bypass to shell profiles,
`al.lua`, or Ansible defaults, and do not use verbose or diff output for the
certificate tasks. The temporary public `known_hosts` projection does not
modify the host's SSH configuration; do not replace it with an unverified
`ssh-keyscan` result or disable SSH host verification.

After renewal, verify Vault with certificate validation enabled:

```sh
env -u VAULT_SKIP_VERIFY bazel_agent bazel run //:vault.status
```

For renewal before expiry, omit `VAULT_SKIP_VERIFY=true` from the Ansible
command. The current inventory serves Vault from `host1`; review the inventory
before changing the host limit. Its inventory name remains
`host1.vault.alwaldend.com` to preserve the Vault node identity; `ansible_host`
connects through the DNS-owned name `host1.vault.dc1.alwaldend.com`. Certificate
alternate names include this connection hostname for the existing unseal
endpoint, as well as the shared Vault service names.

## Read OIDC client info

```sh
bazel run //infra/vault -- read identity/oidc/client/src_infra_dc1_forgejo1_provider
```

## Read all entity aliases

```sh
 bazel run //infra/vault -- list -format json identity/entity-alias/id | jq ".[]" | xargs "-I{}" bazel run //infra/vault -- read "identity/entity-alias/id/{}"
```
