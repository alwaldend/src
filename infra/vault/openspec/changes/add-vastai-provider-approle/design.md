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

## Live verification

On 2026-11-04, a reviewed plan targeted only this project's AppRole module.
The apply created eight identity, policy, and role resources with zero updates
or destroys. The role attaches its standard namespace and named shared-read
policies and joins no general infrastructure group. The subsequent bootstrap
plan authenticated using the new project identity for its Vault destination.
