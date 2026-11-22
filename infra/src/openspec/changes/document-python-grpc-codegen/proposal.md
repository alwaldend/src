## Why

Python gRPC service bindings need a generator in addition to the message
bindings. The former gRPC 1.69.0 rules fail to load under Bazel 9.

## What Changes

- Select gRPC 1.78.0.bcr.1 for native Python stub generation and document the
  rules in repo-python.
- Preserve the Protobuf 33.4 pin aligned with the Android Java runtime.
  The tested gRPC 1.85.4 release raises Protobuf to 35.1.
- Remove the MLflow plugin's standalone Python code generator dependency.

## Capabilities

This documentation adds no runtime behavior.

## Impact

The guidance applies to Python gRPC targets in the root Bazel workspace.
