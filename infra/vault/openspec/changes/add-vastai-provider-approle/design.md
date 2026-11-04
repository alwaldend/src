## Decisions

Use the existing vault_approle module, disable its optional Yandex policy,
and retain the existing operator bootstrap pattern. The project owns its
credential and backend paths under its own AppRole namespace. The bootstrap
root selects the project's existing AppRole identity for Vault writes. Its
standard namespace policy covers token persistence and soft deletion, so no
additional policy or cross-role secret access is required.

## Risks / Trade-offs

The identity definitions require a separately authorized Vault apply before
authentication works. Source checks do not deploy them or write credentials.
