---
title: Go gRPC WASI build metadata
description: WASI dependency correction for the repository-pinned Go gRPC module
---

The root `go.mod` and `go.sum` pin [Go gRPC](https://github.com/grpc/grpc-go).
`tools/go/include.MODULE.bazel` applies `wasip1.patch` after Gazelle generates
the module's BUILD files. The pinned Gazelle platform list omits `wasip1`, so
the generated `internal/syscall` dependency selection loses `grpclog` on WASI.

The patch supplies `grpclog` for the default platform branch, matching the
existing non-Linux Go implementation without changing library source or versions.
Validate through the public `anki_as_code` WASI build and native CLI workflows.
