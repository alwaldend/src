## Context

See `proposal.md` for the requested outcome. Existing third-party command line
tools use `rules_binary_toolchain`, package-local locks, and generated root
module includes. The current checkout is isolated on `t3/add-ytdlp-repo`.

## Goals / Non-Goals

Provide a runnable, integrity-verified yt-dlp executable using existing rules.
This change does not install host software or add optional media tools.

## Decisions

- Use `third_party/com_github_yt_dlp_yt_dlp`, matching the upstream owner and
  artifact, with a repository-internal `yt-dlp` alias.
- Pin the latest stable release observed on 2026-10-23, `2026.10.21`, using
  its Linux x86_64 standalone executable and SHA-256 SRI. This matches the
  platform supported by existing CLI dependencies and needs no host Python
  or new Python dependency lock entries.
- Use the existing repository download action with `executable: true` and
  explicit Linux/x86_64 execution constraints. The upstream standalone binary
  contains its interpreter and bundled Python packages.
- Keep version and integrity exclusively in `binary_toolchain.json`; the
  README documents usage, platform, optional tools, and upstream licenses.
- Use `skip_specs: true`: this adds tooling under the existing reproducible
  Bazel development requirement without changing repository policy.

## Risks / Trade-offs

- Optional processing and site requirements vary. Document FFmpeg/ffprobe and
  JavaScript runtime requirements without adding unrequested dependencies.
- Upstream extractors age. Update the pinned lock through the repository
  workflow and verify checksums when upgrading.
- Startup alone cannot validate every website. Check the actual alias's
  version and extractor enumeration without relying on third-party media.

## Acceptance evidence

Observed on 2026-10-23 before delivery preparation:

- The publisher checksum manifest and release API agree on the executable's
  SHA-256. Bazel fetched it successfully with the recorded SRI.
- The package builds and its alias reports `2026.10.21`; extractor enumeration
  returns 1,755 entries, including `youtube`.
- Module regeneration adds only the new package include. Bzlmod dependency
  reconciliation leaves `MODULE.bazel.lock` unchanged.
- All 112 selected repository quality, Buildifier and OpenSpec tests pass;
  semantic lint of the package and strict change validation pass.
- Repository formatting produces no changes outside this task.

Raw logs and representative output are under ignored `out/add-yt-dlp/`.
The delivery receipt binds the final prepared candidate and its repeated
mandatory validation gates before publication.
