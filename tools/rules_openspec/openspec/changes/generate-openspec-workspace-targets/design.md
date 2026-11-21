## Design

Treat `openspec/config.yaml` as the package-local discovery marker. Generate an `openspec_workspace` source target and a validation target in the same package, including when Gazelle creates its BUILD file. Keep existing source declarations and unrelated targets during Gazelle updates.

Use the pinned CLI to validate both current and archived artifacts. Keep archive validation local so coverage follows subsequent archives without editing BUILD files. Use the published rule module name for standalone consumers and an in-tree load label for the root generator. In standalone archives, remove
the workspace package directive that points to an enclosing monorepo group.

Remove the central validation targets and discovery code. Preserve the project inventory updater's existing project and site metadata projections.
