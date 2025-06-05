---
title: Bazel
---

## Advantages

- Hermetic builds
- Dependency graph
- Local build and test cache
- Remote cache
- Remote execution
- Unified API for builds, tasks, and tests

## Disadvantages

- Additional layer of abstraction and complexity
- Weird edge cases and occasional bugs
- Lack of support, documentation, and information on the internet
- The majority of advantages are irrelevant for the majority of codebases

## .bazelrc has variables

- `%workspace%`: Workspace directory

https://bazel.build/run/bazelrc

## Local registry module is cached

If you changed a local registry module, but bazel still uses the old version,
restart bazel: `bazel shutdown`

https://github.com/bazelbuild/bazel/issues/20477#issuecomment-1851057077

## Output root ownership

Bazel checks ownership of the output root at startup, so it will fail with
`mkdir('/path/to/dir'): (error: 13): Permission denied` if the current user does
not own the directory.

Example:

```sh
$ bazel --client_debug
...
[FATAL 15:19:33.457 src/main/cpp/blaze_util_posix.cc:499] mkdir('/path/to/dir'): (error: 13): Permission denied
```

Traceback:
- https://github.com/bazelbuild/bazel/blob/d798ebde6c6394203a87b5f1a6b62ecfc3880991/src/main/cpp/blaze_util_posix.cc#L497
- https://github.com/bazelbuild/bazel/blob/d798ebde6c6394203a87b5f1a6b62ecfc3880991/src/main/cpp/util/file_posix.cc#L518
- https://github.com/bazelbuild/bazel/blob/d798ebde6c6394203a87b5f1a6b62ecfc3880991/src/main/cpp/util/file_posix.cc#L90
- https://github.com/bazelbuild/bazel/blob/d798ebde6c6394203a87b5f1a6b62ecfc3880991/src/main/cpp/util/file_posix.cc#L48
- https://github.com/bazelbuild/bazel/blob/d798ebde6c6394203a87b5f1a6b62ecfc3880991/src/main/cpp/util/file_posix.cc#L67
