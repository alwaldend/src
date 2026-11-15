## Context

The shared SSH role discovers keys before upstream hardening, which expects all standard key types. Cloud-init may supply only one trusted bootstrap identity.

## Decisions

Use the existing community.crypto.openssh_keypair module before discovery, with regenerate set to never. The shared ssh_host_keys_dir setting also supplies discovery and the upstream hardening role. Remove the secure1-specific shell command.

This preserves existing identities and handles every partial key set. Checking only the ECDSA filename misses other absent types; supplying additional private keys through Terraform would expand bootstrap secret handling unnecessarily.

## Risks / Trade-offs

Key generation must never rotate the seeded identity. Exercise empty, Ed25519-only and ECDSA-only key directories plus repeat convergence using the exact packaged task. Live host configuration is outside this source-validation operation.

Create RSA keys with the same ssh_host_rsa_key_size used by hardening so a
non-default size does not trigger key replacement after certificate signing.
The isolated 3072-bit case reproduced this defect before the correction.
