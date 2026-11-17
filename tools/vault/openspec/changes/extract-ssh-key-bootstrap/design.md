## Context

The existing command has no Forgejo API calls or service-specific algorithm.
Only its names, private-key comment and packaging tie it to that consumer.

## Goals / Non-Goals

Extract the existing implementation without changing its stored schema, key
identity, Vault policy or deployment lifecycle. Do not add algorithms or rotation.

## Decisions

Use tools/vault/cmd/ssh_key, a repository-visible execution tool with the tools
transitive boundary. Inputs become VAULT_SSH_KEY_MOUNT and VAULT_SSH_KEY_PATH.
The command and key comment use neutral SSH terminology. No infrastructure
source becomes a dependency of the tool.

Move the existing combined CLI/Terraform acceptance fixture to
infra/forgejo/test/signing_key. It still executes the actual provisioner source
and new tool binary, retaining its established failure and state-exclusion cases.
Keeping this consumer-specific fixture in infra avoids reversing ownership.

## Risks / Trade-offs

The executable runfiles path and environment names must change together.
Verify the new tool path in the packaged setup target and execute the actual
Terraform declaration in the isolated fixture. Existing cached records remain
valid because validation does not depend on the private-key comment.
