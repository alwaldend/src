## Purpose

Keep Forgejo authentication aligned with its canonical service URL while
permitting certificates for its independently routed local Git alias.

## ADDED Requirements

### Requirement: Canonical Forgejo callback and certificate aliases

The Forgejo OIDC client SHALL derive its allowed callback from the configured
Forgejo URL, defaulting to `https://forgejo.alwaldend.com`, followed by
`/user/oauth2/vault/callback`. Its PKI role SHALL continue to permit both
`forgejo.alwaldend.com` and `git.alwaldend.com`. Existing client, group and
AppRole identities SHALL be retained.

#### Scenario: Configure the canonical service

- **WHEN** Terraform uses the default Forgejo URL
- **THEN** the allowed callback is `https://forgejo.alwaldend.com/user/oauth2/vault/callback`
- **AND** the callback list contains no duplicate literal or old Git URL

#### Scenario: Override the service URL

- **WHEN** an operator supplies another Forgejo URL
- **THEN** the client callback follows that URL without retaining a fixed canonical callback
