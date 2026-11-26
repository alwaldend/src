## Decisions

- Keep upstream artifact pins in their respective `third_party` packages.
  Each upstream package declares its own `:pull` and `:push` targets
  through a shared macro. The generic transfer command lives in `tools/artifact_mirror`.
  The third-party macro supplies the Forgejo owner and authentication.
- Store the community Flux chart and upstream Flux Operator chart as Helm
  archives, and stage image indexes in local OCI layouts.
- Publish charts to Forgejo's Helm registry and preserve multi-platform image
  indexes in its OCI registry.
- Verify artifact digests before publication. Use the third-party Forgejo login
  integration and an invocation-local ORAS credential file.
- Document the actor's upstream pulls, Forgejo pushes, and Vault authentication
  in the tool README with a Bazel-rendered diagram.
