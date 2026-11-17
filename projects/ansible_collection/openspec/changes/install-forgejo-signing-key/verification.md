# Verification

Observed 2026-11-17. The full eight-case
`//projects/ansible_collection/test/forgejo_upgrade:upgrade_test` passes with
real pinned Forgejo binaries. The signing case verifies an actual API-generated
commit with Git, private mode 0600, repeat public-file and permission repair,
unchanged service lifecycle, and rejection of a different installed private key.
It emits `ssh-commit-signing.json` alongside the existing recovery receipts.

Deployment packaging, semantic lint, repository quality, and this owner's
OpenSpec validation pass. The tests use a disposable local service-manager
fixture; no live host or credentials were used. Vault generation and its
separate acceptance evidence belong to the infra/forgejo change.
