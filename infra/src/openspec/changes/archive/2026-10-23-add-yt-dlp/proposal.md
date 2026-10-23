## Why

The user requested yt-dlp, which is not currently available as a repository
dependency. Add a reproducibly pinned executable using the existing binary
toolchain mechanism.

## What Changes

- Add the upstream stable Linux x86_64 executable, verified against its
  publisher checksum, in `third_party/com_github_yt_dlp_yt_dlp`.
- Expose a repository-internal `yt-dlp` Bazel alias and document its use.
- Regenerate the root module and reconcile its dependency lock.

## Capabilities

### New Capabilities

None. This dependency integration follows the existing reproducible Bazel
development contract; it does not change repository specification requirements.

### Modified Capabilities

None. Set `skip_specs: true` for this tooling addition.

## Impact

Adds only the explicitly requested yt-dlp dependency and its generated module
integration. Acceptance requires the runnable alias to report the pinned
version, enumerate extractors, and pass repository validation. Optional media
processing tools and JavaScript runtimes remain caller-provided.
