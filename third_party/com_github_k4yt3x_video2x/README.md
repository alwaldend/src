---
title: Video2X
description: Pinned Video2X frame interpolation executable
---

The [binary lock](binary_toolchain.json) owns the upstream release and SHA-256
integrity for the Linux x86_64 AppImage. It includes the Video2X CLI, FFmpeg
libraries, ncnn, and RIFE models 4.6, 4.25, and 4.26. The release asset has no
publisher-provided digest; the lock records the locally calculated download hash.

The host must provide an AVX2-capable CPU, a Vulkan loader, a compatible GPU
driver, and access to its render device. No desktop session is required.

Run from the repository root with task-local extraction and cache directories:

```sh
mkdir -p out/video2x/tmp out/video2x/cache
export TMPDIR="$PWD/out/video2x/tmp"
export XDG_CACHE_HOME="$PWD/out/video2x/cache"
al bazel run //third_party/com_github_k4yt3x_video2x:video2x -- \
  --appimage-extract-and-run --version
al bazel run //third_party/com_github_k4yt3x_video2x:video2x -- \
  --appimage-extract-and-run --list-devices
```

The AppImage option extracts and runs without FUSE or installation. Select the
device index reported by `--list-devices`. For example, this interpolates a
caller-provided video fivefold using RIFE 4.25 and writes a lossless FFV1 video:

```sh
al bazel run //third_party/com_github_k4yt3x_video2x:video2x -- \
  --appimage-extract-and-run \
  -i "$PWD/out/video2x/input.nut" -o "$PWD/out/video2x/interpolated.nut" \
  -p rife --rife-model rife-v4.25 -m 5 -t 100 -d 0 -a none \
  -c ffv1 --pix-fmt bgr0 --no-copy-streams --no-progress
```

Use absolute media paths. `-t 100`, also the default, disables scene-change
skipping; lower thresholds can repeat frames across large changes. With multiplier
`m`, this release emits `1 + m * (input_frames - 1)` frames. To interpolate a closed loop,
append its first frame to the input, then discard the final output endpoint.
RIFE processes three color channels and does not preserve input transparency.

For repeated runs, the upstream `--appimage-extract` option avoids repeated
extraction. Keep the extracted tree in task scratch. Set `APPDIR` to its absolute
`squashfs-root` path and `LD_LIBRARY_PATH` to its `usr/lib` directory when invoking
`usr/bin/video2x` directly.

## Upstream

- [Video2X source](https://github.com/k4yt3x/video2x)
- [Pinned release and bundled models](https://github.com/k4yt3x/video2x/releases/tag/6.4.0)
- [CLI options for the pinned release](https://github.com/k4yt3x/video2x/blob/6.4.0/tools/video2x/src/argparse.cpp)
- [AGPL-3.0 license](https://github.com/k4yt3x/video2x/blob/6.4.0/LICENSE)
- [RIFE model guidance](https://github.com/hzwer/Practical-RIFE#trained-model)
