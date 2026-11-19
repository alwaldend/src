## Repair nested module resolution

- [x] Remove the unused versionless sibling dependencies from `rules_openspec`.
- [x] Ignore every nested module directory that owns a `MODULE.bazel`.
- [x] Regenerate the affected module locks through the pinned workflow.

## Validate

- [x] Build and test all twelve nested workspaces standalone.
- [x] Verify root target expansion no longer crosses a workspace boundary.
- [x] Verify root gates and `@rules_openspec` consumers still build.
