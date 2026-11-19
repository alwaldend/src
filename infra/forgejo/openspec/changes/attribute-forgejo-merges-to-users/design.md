## Decision

Set `DEFAULT_TRUST_MODEL = collaborator` in the instance configuration.
Forgejo uses this policy where a repository inherits the default. The
collaborator policy allows Forgejo to retain the merging user's committer
identity while signing with the instance key. It also trusts signatures from
repository collaborators under Forgejo's verification rules.

The Ansible package renders and installs the configuration. After an authorized
deployment, a new pull request merge can be inspected for its Git committer
and verified using the public key from `/api/v1/signing-key.ssh`.

[Forgejo signing configuration](https://forgejo.org/docs/v17.0/admin/config-cheat-sheet/#repository---signing-repositorysigning)
