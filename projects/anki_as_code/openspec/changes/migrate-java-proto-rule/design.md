## Decision

Use `@protobuf//bazel:java_proto_library.bzl`, consistent with the other
Java protobuf targets in the repository. The Java wrapper continues to export
the generated library.
