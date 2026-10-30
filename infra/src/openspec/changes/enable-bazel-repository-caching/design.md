## Context

Bazel 9.3.1 defaults to a local repository contents cache. Shared configuration does not enable experimental remote repository caching; the resolved protobuf module exposes prefer_prebuilt_protoc with a false default.

## Goals / Non-Goals

Use existing cache backends and compiler artifacts. Do not provision services, change dependency versions, or impose protobuf settings on standalone nested modules.

## Decisions

Enable the repository-cache startup option in project.bazelrc, shared by root and nested modules. Enable the protobuf build setting only in root.bazelrc, where the module is available. Preserve default local cache location and final user overrides. Keep existing strip_prefix declarations because their explicit directory checks remain useful.

## Risks / Trade-offs

Experimental repository caching can expose repository materialization bugs; validate representative fetch and build operations. Cached repositories must be declared reproducible by their rules. Prebuilt protoc falls back to source on unsupported platforms. A startup-option change restarts the worktree Bazel server once.

## Migration Plan

Validate source and representative protobuf output, run repository quality and affected semantic lint, then publish. Rollback removes the two flags.

## Verification

The representative Python protobuf build passed and produced api_pb2.py with protobuf 6.33.4 metadata. Action graph inspection confirms GenProto and GenProtoDescriptorSet use the pinned prebuilt Linux x86-64 protoc; protobuf-internal ProtoCompile actions can still build the source compiler. OpenSpec strict validation, all 27 repository quality/environment tests, and semantic lint of tools/bazelrc and infra/src/openspec passed. Final candidate checks and publication state belong to the ignored delivery receipt.
