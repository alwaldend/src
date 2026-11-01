---
title: yt-dlp
description: Pinned yt-dlp binary toolchain
---

The [binary lock](binary_toolchain.json) owns the upstream release and SHA-256
integrity for the Linux x86_64 standalone executable. It bundles Python and
upstream's Python dependencies and requires glibc 2.17 or newer.

Run it through the repository alias:

```sh
al bazel run //third_party/com_github_yt_dlp_yt_dlp:yt-dlp -- --version
al bazel run //third_party/com_github_yt_dlp_yt_dlp:yt-dlp -- --help
```

Pass download options and URLs after `--`. Use absolute output paths to select
a download destination. Update the checked-in lock to upgrade this dependency.

FFmpeg and ffprobe are needed for merging separate audio/video streams and
many post-processing operations. Full YouTube support also needs a supported
JavaScript runtime. These executables must be supplied by the caller; see
the upstream dependency documentation for supported runtimes and options.

## Upstream

- [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- [Release artifacts and checksums](https://github.com/yt-dlp/yt-dlp/releases)
- [Dependencies](https://github.com/yt-dlp/yt-dlp#dependencies)
- [License and bundled dependency licenses](https://github.com/yt-dlp/yt-dlp#license)
