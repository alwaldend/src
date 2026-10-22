## Why

PR #128 includes the Anki upstream schema dependency alongside the CLI. Extract the dependency into an independently reviewable prerequisite PR as requested.

## What Changes

- Publish the existing integrity-pinned Anki source archive and Bazel-generated Go schema targets separately.
- Regenerate the root module include through its owning tool.
- Keep both PRs against master; remove the matching dependency changes from #128 after this prerequisite merges.

## Capabilities

No behavior or requirements change; this is a publication split of an already validated dependency. `skip_specs: true` applies.

## Impact

`third_party/com_github_ankitects_anki`, root `MODULE.bazel`, and the shared repository change record. CLI implementation and collection data remain in their owning PRs.
