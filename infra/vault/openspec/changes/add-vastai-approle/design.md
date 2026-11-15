## Context

The existing GitHub AppRole declaration already provides the needed pattern.
See [proposal](proposal.md) for the requested identity.

## Goals / Non-Goals

Reuse the shared module's identity, namespace, and bootstrap policy. No new
policy implementation or live Vault operation is required.

## Decisions

Copy the GitHub declaration shape with `src_infra_vastai` as its name and
`simeonwarren` as operator. Disable Yandex policy attachment and omit the general
infrastructure group because neither capability is needed for this foundation.

## Risks / Trade-offs

Checked-in declarations do not establish deployed access. Deployment remains
a separately authorized apply through the existing Vault Terraform root.

## Authorized SSH reconciliation, 2026-11-15

The operator approves deployment of the existing bot identity and account SSH
registration, then separately approves the declared metadata-read policy and
AppRole attachment. The first identity apply finds no changes. The Vault policy
apply adds its exact-path read policy and attaches it to the account AppRole.
The stored public key has fingerprint BUxVSqTBim2ZUgOZs4oIBKJWs2PPSwOoxYG6/x+JuX8,
while this host public key and loaded agent match
MECEFy9M8xrFmdEYXb+NlwB6e7imUFhIdpf0Adb5w0w. Automatic approval review
rejects replacing the stored key under the earlier registration-only approval.
The operator then explicitly approves replacing only this identity's sshpubkey
with the current host public key and updating its Vast account registration.
Do not change private keys, host configuration, aliases or other identity fields.
Record targeted outcomes and do not apply unrelated account team deletions.

The reviewed public-key reconciliation updated only one Vault entity metadata
field, then updated the existing account key registration in place. Both applies
reported zero additions and zero deletions. The prior access-policy apply reported
one policy addition and one AppRole update. No private key or host configuration
changed.
