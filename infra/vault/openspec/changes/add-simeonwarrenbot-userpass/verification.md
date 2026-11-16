# Selected userpass verification

A disposable loopback Vault server exercised the exact candidate with Vault
provider 5.8.0 and random provider 3.9.1. It used in-memory Vault storage,
private temporary Terraform state, -dev-no-store-token, and an integrity check
on the shared token helper. No production endpoint was contacted.

Verified outcomes:

- Only the selected bot received a generated entry; the human was unselected.
- The initial and rotated passwords both had exactly 72 characters.
- The bot could not read, write, or delete its userpass data or metadata, or
  list its namespace. This follows from absent grants, without deny rules.
- The human could read the bot entry but not modify it, list the bot namespace,
  or read siblings. The infrastructure writer also could not read siblings.
- KV creation, metadata refresh, and rotation with computed metadata succeeded.
  Infrastructure metadata deletion remained denied.
- The matching stored password authenticated to the existing bot entity.
- Existing token policies and TTL/max-TTL settings survived initial management
  and subsequent password replacement.
- The final ordinary plan had no changes and exited zero.
- The shared token helper remained unchanged; the fixture server, state,
  credentials, and logs were removed.

Input hashes:

- userpass.tf: `b2c6f08c767f0b566fdc07b56c766d6ea04833aae0a50d89146fdedc95ba5f97`
- user_simeonwarrenbot.tf: `704df838b5d530c60a9d1643ef86f53e4f3153ac2464db6e91cf2bb5dcdb5e55`
- group_users.tf: `0583201526e6d3a9b5c420efc65d3236945655f5e5dbaae54850dae608c00342`
- user_simeonwarren.tf: `dcecf4379f404010c08a1afda17eec01efa34f4d50c06a29deb6fe203edc1d5b`

The fixture copied userpass.tf, user_simeonwarrenbot.tf, and group_users.tf
verbatim. It projected only the human entity and its new reader policy from
user_simeonwarren.tf, leaving unrelated existing human modules outside this
exercise. A pre-existing fixture bot account had non-default token lifetimes
and a known policy assignment to verify their preservation.

The task-local repeatable driver is
`out/simeonwarrenbot-auth/verify_userpass.py`, using the pinned Terraform and
Vault launchers beside it. Each local apply requires interactive review. The
sanitized result is `userpass-verification.json` beside the driver.

## Production deployment

The user explicitly requested "Fix comments and apply" on 2026-11-16, approving
the previously described scope. The implementation applied was commit
`4817189643df922384c1389794f74f5310cc4eab`. Two invocations of the owning
`//infra/vault/tf:tf.apply` target displayed normal text plans, reviewed before
answering yes. No saved plans, JSON plans, or state exports were used.

1. At 17:13 UTC, create userpass_writer and simeonwarrenbot_userpass_reader,
   and update the existing simeonwarren entity and infrastructure AppRole
   assignments: two additions, two in-place changes, no destroys.
2. At 17:14 UTC, use a fresh AppRole token to create the 72-character generator,
   store the selected userpass entry, then update the existing bot login:
   three managed resources added, no destroys or unrelated changes. The login
   resource's creation adopted the existing account through a password-only
   update; its original alias/entity remained intact.

The subsequent full, untargeted `//infra/vault/tf:tf.plan` with
`-detailed-exitcode` reported no changes and exited zero. It emitted an unrelated
Yandex provider version-service connectivity warning; resource refresh and
planning completed successfully.

Live verification kept credentials in memory and confirmed:

- The host token resolved to the existing human entity and could read the
  bot username/password, whose password length was 72.
- Human capabilities were exactly read for the entry's data/metadata and deny
  for a sibling or namespace listing; actual sibling reads and listing failed.
- The stored password authenticated to the original bot entity. Its data,
  metadata, sibling, and namespace capabilities were deny, and actual reads
  and listing returned HTTP 403.
- The temporary verification token was revoked and the shared host token
  helper remained byte-for-byte unchanged.

The task-local driver is `out/simeonwarrenbot-auth/verify_live_userpass.py`,
run through the existing injected host authentication and pinned Vault CLI.
It retains only the sanitized `live-userpass-verification.json` report. An
initial report attempt parsed the successful token-revocation text as JSON;
the driver was corrected and all checks passed on its next run.

Prior designs were never deployed. The infrastructure writer retains exact
metadata read/update without deletion, addressing the original review finding.
The outdated review thread remains unresolved because the delivery adapter
refused outdated-thread handling and automatic approval rejected the attempted
lower-level reply. Neither operation was retried.
