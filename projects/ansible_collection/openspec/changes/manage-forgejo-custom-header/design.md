## Context

Forgejo reads its custom header template when it starts. The role already
stages startup inputs and retains an archive before replacing them.

## Decisions

- Accept the optional source as a literal Forgejo template under the caller's
  control. Resolve the destination under the work directory's custom tree.
- Determine header changes by content alongside the other startup inputs. Track
  the existing executable, configuration, and service unit separately so a
  missing optional header does not make an established installation partial.
- Stage the header before downtime and install it inside the activation
  boundary. Require the same archive coverage as for other startup inputs.
