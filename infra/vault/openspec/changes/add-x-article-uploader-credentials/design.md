## Context

The uploader needs a long-lived OAuth 1.0a credential at run time. Vault is the
only approved place for the value, and the read must be scoped so the uploader
cannot read another component's secrets.

## Decisions

### Keep the credential in the identity's own AppRole subtree

The credential lives at
`alwaldend.com/vault1/approles/src_projects_x_article_uploader/oauth1` on the
`secrets` mount, the same shape every other component uses: a value under the
component's own AppRole subtree, granted by the shared `vault_approle` module
without a bespoke policy. The uploader's `al.lua` names that mount and path
directly, exactly as every other component's injected reference does, so there
is no second declaration to keep in step.

This change previously placed the credential in a separate
`alwaldend.com/vault1/x_api/<name>` namespace, to keep the read-only grant
outside the role's own writable subtree: the shared module grants every role
`create`, `update`, `delete`, and `patch` on `approles/<name>/*`. That
placement was abandoned as unnecessary complexity. The subtree is the role's
own storage — its Terraform state and its secrets — and a component that can
write its own credential can always rotate it; the boundary that matters is
that this identity cannot read or write _another_ component's subtree, which
the shared module's per-role policies already enforce.

An earlier revision also carried the mount and path as separate checked-in files
read by both the identity and `al.lua`, with output preconditions comparing them
against the stage's mount and the role's own subtree. That indirection was
removed: it duplicated the reference the sibling components already write as a
literal, and the one literal is simpler to read and to change.

### Reference the credential by path only

The identity's policy grants the role's own subtree by construction, and
`al.lua` requests the credential from the same subtree by naming it directly. No
credential value is checked in; only the mount and the logical path are named.

Alternatives considered: a separate namespace outside the AppRole subtree (the
previous decision) and narrowing the shared module's base policy for this one
role. The separate namespace adds a convention the shared module already
covers; narrowing the shared policy changes a shared lifecycle relied on by
every component identity. Neither earns its cost.

## Risks / Trade-offs

[The identity can write as well as read its own credential] -> That is the
shared module's per-role contract for every component, and rotation is the
capability it exists to grant. The boundary this change must preserve is
between identities, not between one identity and its own storage, and the
shared policies grant nothing on another component's subtree.

[The new path needs an authorized operator write before first use] -> The value
is written by an authorized operator through an operator-authenticated target.
The component's own `vault.kv_put` authenticates as
`src_projects_x_article_uploader` and could write its own subtree, but the
credential is provisioned by an operator, for example
`bazel_agent bazel run //:vault -- kv put -mount secrets <path> <field>=-`
with the value read from stdin. No live write or apply is authorized by this
change.
