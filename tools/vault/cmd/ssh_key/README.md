---
title: Vault SSH key bootstrap
description: Create or reuse an Ed25519 SSH key pair stored directly in Vault
---

`ssh_key ensure` generates a dedicated Ed25519 OpenSSH pair
in memory and stores `algorithm`, `private_key`, `public_key`, and `fingerprint`
in KV v2 using CAS=0. No private values enter Terraform attributes, output,
arguments, or local files. The pair is intentionally unencrypted for unattended
consumers; Vault and deployed file permissions protect it.

The command consumes the existing AL `VAULT_ADDR`, `VAULT_TOKEN`, and TLS client
configuration. `VAULT_SSH_KEY_MOUNT` and `VAULT_SSH_KEY_PATH`
identify its destination. HTTPS verification is mandatory and redirects are
rejected. The operation is bounded to one minute with retries disabled.

An existing record must contain a valid matching Ed25519 pair and fingerprint.
Missing data with existing metadata, including deleted or destroyed versions,
fails rather than generating a replacement. Permission and transport failures
also stop bootstrap. CAS prevents overwrites if another writer wins a race.
If persistence succeeds but its response is lost, the next invocation validates
and reuses that record. Restore Vault history rather than deleting metadata to
work around a failure.

The caller owns the selected Vault namespace, permissions, serialization and
deployment of the pair. Use an owning IaC backend lock to serialize execution;
concurrent standalone execution is unsupported. A creation-only Terraform
provisioner can invoke the command without storing private attributes. Normal
applies then preserve the resource, and explicit replacement reuses the cache.
Removing that resource retains the key; it does not detect later Vault drift
or rotate the identity.

The [Forgejo consumer acceptance fixture](../../../../infra/forgejo/test/signing_key/README.md)
executes this binary and its actual Terraform provisioner against isolated TLS
Vault services. The tool itself has no infrastructure source dependency.
