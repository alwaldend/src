# Src push-mirror evidence

Both live native HTTPS mirrors copied branches and tags from Forgejo to their
GitHub and GitLab destinations. GitLab returned a mirrored LFS object with
the same SHA-256 as its source pointer. The GitHub-only branches were copied
to Forgejo and both divergent histories were reconciled before activation.

The Vault-backed Terraform root bootstrapped the existing mirrors. GitLab
was adopted without a mirror API write. GitHub's HTTPS entry replaced its
failed SSH entry. The catalog now maps one explicit entry to each mirror.
The state-address moves preserve the two existing Terraform identities.

The HTTPS credentials use repository push permissions. The temporary
GitHub App and SSH branch grants are absent from the desired configuration.
Removing those live grants and disposing of unused credentials require
their own approved operations.
