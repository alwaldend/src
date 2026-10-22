## Why

The X Article uploader authenticates to the X API with OAuth credentials that
must never be committed. It has no component identity, so the repository's
injection path cannot supply them: no AppRole exists to authenticate the read,
no Vault path holds the credential reference, and no policy grants the
component access to that path alone.

## What Changes

- Declare a `src_projects_x_article_uploader` identity, composing the existing
  AppRole module rather than introducing a second authentication mechanism.
- Store the OAuth 1.0a credential — API key, API secret, access token, and
  access token secret — under the component's own AppRole subtree, by reference
  only: no credential value enters checked-in source.
- Read it through the shared module's own-subtree policy, so the identity cannot
  read another component's secrets.
- Record the injected environment variables the draft command consumes, and the
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
