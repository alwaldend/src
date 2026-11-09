## Why

The existing Vast infrastructure identity must deliver its generated team token to the user's KV entry without authenticating as that user.

## What Changes

Grant the existing src_infra_vastai role lifecycle access to exactly the imagegen/vastai data entry and read access to its metadata. Add no AppRole or general user-namespace access.

Enable direct access for the human simeonwarren identity to its own KV namespace by reusing the personal role's existing namespace policy and member group.

## Capabilities

- `vastai-token-delivery`: The team writer publishes the generated token to one regular-user KV entry.
- `personal-kv-ui-access`: The human login can discover and access its own KV namespace.

## Impact

One policy, one existing AppRole attachment, and one existing personal-group policy attachment. The user explicitly authorized deploying both scoped permissions. Source validation does not deploy them.
