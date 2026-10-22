## 1. Component identity and policy

- [x] 1.1 Add the `src_projects_x_article_uploader` identity under `infra/vault/tf`, composing the existing AppRole module and joining the shared memberships the other project identities use, and register it in the stage's group lists; verify the Terraform root validates and the identity appears in the group membership
- [x] 1.2 Keep the credential inside the identity's own AppRole subtree, read through the shared module's own-subtree policy, matching every other component; verify the policy grants the uploader's path and denies another component's credential path
- [x] 1.3 Declare the credential mount and path by reference with no value in checked-in source, and record the environment variables the draft command consumes; verify a repository search finds no credential value and both references resolve in the configuration
- [x] 1.3a Name the credential solely by reference in `al.lua` (mount and logical path, no value); verify a repository search finds no credential value and the reference resolves in the configuration
- [x] 1.4 Verify the changed Terraform packages build and their structural checks pass, without creating a role, writing a credential, or authenticating to Vault
