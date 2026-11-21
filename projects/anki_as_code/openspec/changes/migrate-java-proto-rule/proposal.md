## Why

Bazel reports a deprecation warning when building the Anki as code Java
protobuf target.

## What Changes

Load the Java protobuf rule from the pinned Protobuf dependency.

## Impact

The generated Java library retains its existing Bazel label and dependency
on the project's protobuf schema. This is part of the
[repository CI reporting follow-up](../../../../../infra/src/openspec/changes/improve-forgejo-ci-output/proposal.md).
