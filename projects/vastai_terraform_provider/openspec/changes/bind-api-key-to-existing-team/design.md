## Context

The stored personal issuer has key creation and context-selection permissions.
Directly requesting another account's full permissions returned HTTP 401. That
does not establish that the credential cannot create keys generally.

## Decisions

Extend the existing API-key resource with a replacement-only optional team ID.
Resolve the owner context before creation and persist the new key before using
the SDK's existing verified activation flow on that new key. No issuer context
is changed. An empty permission string omits the permissions field to request
the platform's default full-access key; an empty JSON object is not equivalent.
The team's owner context supplies its final team permissions. Reject scoped
permission input for this mode before API calls to avoid widening its scope.

This is preferred to modifying a shared credential, replacing the team, or
inventing additional permission overrides. Live activation remains an explicit
acceptance gate; the earlier isolated activation received HTTP 400, so a new
attempt must preserve the key for reviewed cleanup and report failure honestly.

## Failure cases before implementation

Terraform CLI fixtures cover successful binding, rejected activation, and a
missing owner context. They verify no team creation or credit mutation and
revocation after partial creation. No credential is printed or written to a
Vault destination after a failed binding.

## Provider authentication

Use the configured provider client when no resource issuer override is supplied.
Keep the optional `team_api_key` only for a management key generated during the
same bootstrap apply. Provider configuration requires known credentials, so
removing that bootstrap mechanism would require a separate workflow migration.
Removing redundant issuer metadata must preserve the generated key. Verify key
visibility through the provider before recording that removal and preserve the
existing key ID and value. A new explicit bootstrap issuer, key name,
permissions, or team binding retains replacement behavior. Existing-team callers and ordinary test resources omit
the resource credential and use their provider's Vault-injected identity.

E2E cases must prove provider selection with an unusable environment fallback,
all CRUD operations without a resource credential, and removal of a legacy
override without generating or revoking a key. A different explicit bootstrap
issuer must still plan replacement because it can represent a new team.
