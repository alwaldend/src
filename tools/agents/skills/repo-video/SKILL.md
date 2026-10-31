---
name: repo-video
description: >-
  Create, edit, convert, mux, and verify video artifacts in this repository,
  including still images with soundtracks. Use repo-blender for Blender scenes
  and mermaid-diagrams for diagram authoring.
---

# Work with repository video

## Establish the media contract

Inspect the actual inputs and requested destination before choosing a command.
Preserve explicit container, codec, dimensions, timing, crop, and soundtrack
requirements. A valid file that the user's player or upload service rejects
does not fulfill the request.

Use `al-workspace` for task output and temporary storage. Keep generated
media, browser profiles, intermediate files, and reports under ignored
`out/<task>/`; preserve the source media. Artifact creation alone does not
require a source commit. Use the repository change and delivery procedures
when editing maintained source, rules, or assets.

For ordinary playable video, prefer H.264 video and AAC audio in an MP4 with
fast start. A `.mp4` extension alone says nothing about the codecs. Use WebM or
other formats when the user or consumer supports them; do not silently substitute
WebM for a compatibility requirement. Check the destination's current limits
when preparing an upload, rather than applying one platform's limits everywhere.

## Select available tools

Prefer an existing repository entry point and use `al-bazel`/`repo-bazel`
for Bazel commands. When no suitable entry point exists, use available local
FFmpeg and ffprobe. Inspect their version, encoder list, and selected encoder's
help. Before a long export with an unfamiliar encoder, actually encode a short
sample: a listed encoder can be a nonfunctional stub.

If host H.264 fails, the repository's pinned Chrome can supply a working
WebCodecs encoder without installing software. Read
[browser H.264](references/browser-avc.md) for tool resolution, encoding, and
playback verification. Check this path before declaring compatible MP4 output
unavailable. Do not install codecs or change persistent host configuration as
an implicit part of a media task.

Other owners supply specialized operations:

- `mermaid-diagrams` renders canonical diagrams; use its PNG or WebP output as
  video input instead of redrawing the diagram.
- `repo-blender` owns Blender scenes and their rendering.
- The [Video2X README](../../../../third_party/com_github_k4yt3x_video2x/README.md)
  owns pinned interpolation commands and GPU requirements. Its bundled FFmpeg
  libraries do not establish that a standalone FFmpeg executable is available.
- The [yt-dlp README](../../../../third_party/com_github_yt_dlp_yt_dlp/README.md)
  owns the pinned downloader when downloading is requested.

## Inspect and process explicit streams

Probe inputs with bounded `ffprobe -v error ... -of json` commands. Identify
video/audio streams, dimensions, frame rate, aspect ratio, duration, start
timestamps, rotation, and relevant color information. Missing or approximate
duration metadata is not a reason to guess the requested coverage.

Use explicit `-map` selections. Copy compatible streams when their contents and
timing need no change; encode only what the edit or destination requires. For
precise cuts, stream copying can be limited by keyframes: verify boundary frames
and re-encode where precision requires it. Keep diagrams and text readable;
do not crop, stretch, or lower resolution merely to satisfy a convenient preset.

For a still image with a complete soundtrack, read
[the still-image recipe](references/still-image.md). Container durations include
codec timing details, and `-shortest` or an arbitrary `-t` can lose audio packets.
Prefer finite video plus the complete audio input, and verify the final samples.

## Verify the resulting artifact

Check the produced file, not just the exit code or requested encoder settings:

- Confirm the actual container, codecs, required streams, dimensions, frame
  rate, duration, and fast-start placement where relevant.
- Decode or play the output with a working decoder. Inspect representative
  frames and seek near the beginning, middle, and end. If the host decoder is
  also a stub, use the pinned browser and record that verification method.
- For copied audio, compare decoded PCM hashes against the complete source;
  plausible container duration alone does not prove preservation. For lossy
  encoding or resampling, hashes normally differ: compare decoded sample
  coverage, timing, and representative audio instead of demanding equality.
- For a static video, compare decoded frames against the intended still with
  allowance for deliberate scaling and lossy encoding. Do not claim pixel
  identity after a lossy conversion.

Retain a small report with input/output paths, commands and tool versions,
actual media properties, and verification results. Keep media payloads out of
ordinary diagnostics. Return a usable absolute artifact path and identify its
format and material limitations. Creating a video does not authorize uploading
or publishing it, or committing user attachments and disposable intermediates.
