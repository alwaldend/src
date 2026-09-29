## 1. Component identity and policy

- [ ] 1.1 Add the `src_projects_x_article_uploader` identity under `infra/vault/tf`, composing the existing AppRole module and joining the shared memberships the other project identities use, and register it in the stage's group lists; verify the Terraform root validates and the identity appears in the group membership
- [ ] 1.2 Add the least-privilege policy granting that identity read access to its own credential path only, following the shape of a sibling project identity; verify the policy grants the uploader's path and denies another component's credential path
- [ ] 1.3 Declare the credential path by reference with no value in checked-in source, and record the environment variable the publisher consumes; verify a repository search finds no credential value and the reference resolves in the configuration
- [ ] 1.4 Verify the changed Terraform packages build and their structural checks pass, without creating a role, writing a credential, or authenticating to Vault
