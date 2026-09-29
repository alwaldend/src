## Why

The X Article uploader authenticates to the X API with OAuth credentials that
must never be committed. It has no component identity, so the repository's
injection path cannot supply them: no AppRole exists to authenticate the read,
no Vault path holds the credential reference, and no policy grants the
component access to that path alone.

## What Changes

- Declare a `src_projects_x_article_uploader` identity, composing the existing
  AppRole module rather than introducing a second authentication mechanism.
- Store the X API credential in the component's own Vault path, by reference
  only: no credential value enters checked-in source.
- Grant the identity read access to that path and nothing else, so the
  publication credential cannot read another component's secrets.
- Record the injected environment variable the publisher consumes, and the
  non-secret project wiring that selects the injection.

## Capabilities

### New Capabilities

- `x-article-uploader-identity`: component-scoped authentication and
  least-privilege access to the uploader's own X API credential.

### Modified Capabilities

None. Existing identity and policy requirements continue to apply.

## Impact

- `infra/vault/tf/**`: the new component identity and its policy, following the
  existing project-identity shape.
- The uploader's `al.lua`, its `al_config` target, and the plugin data and
  labels that select the injection belong to
  `projects/x_article_uploader`; this change owns the identity and policy they
  authenticate against.
- Nothing here authorizes a live credential write, role creation, or
  publication. Those remain separate authorized operations.
