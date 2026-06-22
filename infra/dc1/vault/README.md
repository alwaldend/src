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
bazel run //infra/dc1/vault/tf_setup:tf.apply # Create VMs (requires an active Vault host)
bazel run //infra/dc1/vault/ansible # Set up hosts (BM and VMs)
bazel run //infra/dc1/vault/tf:tf.apply # Configure vault
```

## Backup

```sh
bazel run //infra/dc1/vault:backup
```

## Unseal

With a working Vault:
```sh
bazel run //infra/dc1/vault:unseal
```

Without a working Vault:
```sh
bazel run //infra/dc1/vault:unseal_standalone
```

## Fix quorum

```sh
bazel run //infra/dc1/vault/ansible:fix_quorum
```

## Set up only VMs

```sh
bazel run //infra/dc1/vault/ansible:ansible.vm # Set up only VMs
```

## Set up only bare metal

```sh
bazel run //infra/dc1/vault/ansible:ansible.bm # Set up only bare metal
```

## Tf

Plan:
```sh
bazel run //infra/dc1/vault/tf:tf.plan
```

Apply:
```sh
bazel run //infra/dc1/vault/tf:tf.apply
```

Run terraform directly:
```sh
bazel run //infra/dc1/vault/tf:tf.direct -- -chdir="${PWD}" plan
```

## Replace VMs

```sh
bazel run //infra/dc1/vault/tf:tf.apply -- -replace 'module.vm_ha["host2"].proxmox_vm_qemu.vm' -replace 'module.vm_ha["host3"].proxmox_vm_qemu.vm
```

## Generate and import a user cilent certificate

```sh
username="username"
bazel run //infra/dc1/vault:gen_client_cert -- --user "${username}" --output_dir "${PWD}"
bazel run //tools/ykman -- piv certificates import 9A "${PWD}/${username}.pfx"
bazel run //tools/ykman -- piv keys import 9A "${PWD}/${username}.pfx"
```

## Generate a host cilent certificate

```sh
bazel run //infra/dc1/vault:gen_client_cert -- --host some-host --output_dir "${HOME}/.al/client_cert"
```

## Generate a user device cilent certificate

```sh
bazel run //infra/dc1/vault:gen_client_cert -- --hostsome-host --user username --output_dir "${HOME}/.al/client_cert"
```

## Unseal

- Prepare encrypted unseal token
- Run and input the encrypted token:
  ```sh
  bazel run //infra/dc1/vault:unseal
  ```

## Root token

- Prepare encrypted unseal token
- Run and input the encrypted token:
  ```sh
  bazel run //infra/dc1/vault:gen_root_token -- --pgp_key path_to_public_gpg_key_in_base64
  ```

## Generate an EAB for ACME

```sh
bazel run //infra/dc1/vault -- write -f pki/ica_servers/roles/ica_servers_dc1_pve1/acme/new-eab
```

## Sign a client ssh key

```sh
bazel run //:vault -- write ssh/clients/sign/admins ttl=30000000  public_key=@"${HOME}/.ssh/key"
```

## Revoke all tokens

```sh
bazel run //infra/dc1/vault -- token revoke -mode=path auth
```

## Vault certificates

Vault certificates (`tls_cert_file`, `tls_key_file`) should be updated manually


## Read OIDC client info

```sh
bazel run //infra/dc1/vault -- read identity/oidc/client/src_infra_dc1_forgejo1_provider
```
