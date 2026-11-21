# Src push-mirror evidence

Both live native HTTPS mirrors copied branches and tags from Forgejo to their
GitHub and GitLab destinations. GitLab returned a mirrored LFS object with
the same SHA-256 as its source pointer. The GitHub-only branches were copied
to Forgejo and both divergent histories were reconciled before activation.

The Vault-backed Terraform root bootstrapped the existing mirrors. GitLab
was adopted without a mirror API write. GitHub's HTTPS entry replaced its
failed SSH entry. The catalog now maps one explicit entry to each mirror.

The HTTPS credentials use repository push permissions. The temporary
GitHub App and SSH branch grants are absent from the desired configuration.
Removing those live grants and disposing of unused credentials require
their own approved operations.

## Ref filter verification

The Go tests covered filtered mirror adoption and rejection of a missing
trunk, a destination-only branch, or a differing branch tip. A local bare
destination seeded with a development branch was synchronized with Forgejo
17.0.1's three push refspecs and `git push --mirror`. The resulting inventory
contained trunk, a weekly release branch, and a tag. The existing development
branch was pruned and a second development branch was excluded.

On 2026-11-21, Terraform reconciled both src mirrors and completed their
state-address migration. The Forgejo API reported `master,releases/*` and
accepted a synchronization request. GitHub and GitLab each pruned 22
development branches. The source retained all 22 at their previous tips.
Both destinations matched every allowed source ref. The Terraform plan
after removing the completed address moves reported no changes.
