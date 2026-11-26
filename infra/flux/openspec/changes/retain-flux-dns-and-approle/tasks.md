## 1. Retire the Flux cluster source

- [x] 1.1 Remove the cluster packages and obsolete configuration, preserve the DNS targets and declarations, and update the main specification.
- [x] 1.2 Remove the former security inventory reference and use Harbor's own SOPS target in its README. Check both surviving consumers in the Bazel graph.
- [x] 1.3 Move the AppRole and policies into a Flux-owned state root. Verify both Vault roots offline and verify the indexed DNS policy relinquishment with a local fixture.
- [x] 1.4 Validate the DNS declarations, surviving build targets, and OpenSpec workspace offline.
- [x] 1.5 Move the remaining Flux-specific Vault objects into the project. Validate both roots and verify every old managed address has a non-destructive transfer declaration.
- [x] 1.6 Adopt all 53 live bindings, verify OIDC credentials and exclusive Flux state ownership, and remove the temporary transfer declarations.
