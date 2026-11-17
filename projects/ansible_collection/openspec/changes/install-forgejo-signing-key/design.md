## Context

The role already creates the service account and mounts the data disk before
validating and activating startup files. Signing files must exist at activation.

## Goals / Non-Goals

Install an optional supplied SSH pair, preserve the private identity, and retain
default backup coverage. Vault generation and rotation belong to the caller.

## Decisions

Use an opt-in pair dictionary, empty by default. Keep the path fixed below
forgejo_dir so the existing archive contract covers both files. Reject an
existing different private checksum before copying. Derive the public key with
ssh-keygen before accepting the supplied pair. Install mode 0600 private and
0644 public beneath a service-owned 0700 directory. No service restart is needed
when restoring the same missing pair; configuration changes use guarded activation.

## Risks / Trade-offs

Existing different keys require explicit operator reconciliation. The fixture
must test a genuine SSH-signed Forgejo API commit, repeats, public-file repair,
and private mismatch refusal before service downtime. Failure cases are also
recorded in the owning infrastructure change before implementation.
