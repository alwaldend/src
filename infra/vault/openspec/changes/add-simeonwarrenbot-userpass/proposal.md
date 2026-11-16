## Why

Selected users need generated Vault passwords stored in per-user userpass
entries. Users must not receive automatic access to those entries;
simeonwarren needs one explicit read grant for the bot.

## What Changes

- Generate stable 72-character passwords for an explicit user list containing
  only simeonwarrenbot.
- Store username/password pairs at
  `secrets/alwaldend.com/users/<username>/userpass`.
- Grant infrastructure management of those selected entries and grant
  simeonwarren read access only to simeonwarrenbot/userpass.
- Omit user/group access grants and deny rules.
- Update each selected Vault login only after its password has been stored.

## Capabilities

### New Capabilities

- `bot-userpass`: Selected-user password generation, private per-user storage,
  and one human read grant.

### Modified Capabilities

None.

## Impact

Vault Terraform and this existing change record. No new dependency or identity
group. Earlier designs were never deployed and required no migration. The user
explicitly authorized the scoped live apply on 2026-11-16.
