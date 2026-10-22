## Why

PR #128 mixes the Anki CLI with shared Go dependency and agent skill changes.
The user requested a separate PR for the shared changes so each scope can be
reviewed independently.

## What Changes

- Extract the existing approved Go module pins, sums, and Bazel repository
  exposure without changing their versions or introducing dependencies.
- Extract the Cobra CLI and latest stable dependency guidance, together with
  their owner-local OpenSpec records and existing evaluation cases.
- Publish this prerequisite against master. The Anki CLI needs the upgraded
  SQLite collation API; its existing PR will be synchronized after this merges.

## Capabilities

No new behavioral specification: this change separates already implemented
dependency and workflow changes into an independently validated candidate.

## Impact

Shared go.mod/go.sum resolution, tools/go/include.MODULE.bazel, and the
tools/agents skills. No Anki source or collection content is included.

## Acceptance

The extracted files match PR #128's candidate byte for byte. Repository quality,
semantic lint, skill harness checks, and existing shared Go consumers pass.
Both PRs retain master as their base; no merge or stack is authorized.
