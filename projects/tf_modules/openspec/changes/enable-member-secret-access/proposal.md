## Why

Human members of a personal AppRole need direct access to their own KV namespace for Vault UI use. Bootstrap membership currently grants only AppRole login capabilities.

## What Changes

Add an explicit opt-in flag that attaches the role's existing namespace policy to its existing member group. The default remains false; no new policy or group is created.

## Capabilities

- `member-secret-access`: Explicitly selected role members can access the role's own KV namespace.

## Impact

The personal user role opts in. All other role groups retain their existing policies. No mount visibility setting is changed.
