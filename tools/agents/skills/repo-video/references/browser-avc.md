# H.264 through the existing browser

Use this when compatible MP4 output is required and the host FFmpeg H.264
encoder is unusable. Both the encoder and decoder can be stubs; a native
decoder error does not by itself invalidate a browser-produced file.

## Resolve the existing tools

From the repository root, build the existing Puppeteer links and generate a
launcher for the pinned browser. Use task-local temporary and profile storage.

```sh
bazel_agent bazel build //tools:node_modules/puppeteer
bazel_agent bazel run --script_path=out/video-export/chrome \
  //third_party/com_google_chrome_headless_shell:chrome_headless_shell
```

In a task-local CommonJS program using the available Node runtime:

```javascript
const path = require("node:path");
const puppeteer = require(
  path.resolve("bazel-bin/tools/node_modules/puppeteer"),
);
async function launchVideoBrowser() {
  return puppeteer.launch({
    executablePath: path.resolve("out/video-export/chrome"),
    headless: true,
    userDataDir: path.resolve("out/video-export/browser-profile"),
    args: [
      "--no-sandbox",
      "--disable-dev-shm-usage",
      "--disable-background-networking",
    ],
  });
}
```

Place asynchronous code in an async function when using CommonJS. The browser
launch convention matches existing repository browser tools. Do not copy a
foreign worktree's cache path or install another browser or npm package.

## Encode the intended frames

Load a task-local HTML file through `file://` so WebCodecs has a secure context.
Check `VideoEncoder.isConfigSupported`, then actually encode and flush one
frame before a long export. One tested portrait configuration is:

```javascript
const config = {
  codec: "avc1.640029",
  width: 1080,
  height: 1438,
  bitrate: 8000000,
  framerate: 30,
  hardwareAcceleration: "prefer-software",
  avc: { format: "annexb" },
};
```

Choose dimensions, H.264 level, and bitrate for the actual output; the example
is not a universal preset. For a still, draw it once into an `OffscreenCanvas`.
Submit `ceil(audioSeconds * fps)` frames with timestamps
`Math.round(index * 1000000 / fps)` in microseconds and a declared frame duration.
Request periodic keyframes so seeking works, and close each `VideoFrame`.

Apply backpressure to `encoder.encodeQueueSize` instead of accumulating every
frame in flight. Flush at the end and check asynchronous errors and frame counts.
When decoding an existing clip, flushing `VideoDecoder` between arbitrary
batches can require another keyframe; use queue backpressure instead.

Write encoded Annex B chunks to a task-local `.h264` file in output order.
Keep payloads out of progress logs. Close the browser in `finally`.

## Mux and validate MP4

For fixed-rate Annex B output, declare the same input rate during muxing. This
example encodes an incompatible input soundtrack to AAC; copy existing AAC
when it already meets the requested output contract.

```sh
timeout 180 ffmpeg -nostdin -hide_banner -loglevel warning -n \
  -fflags +genpts -r 30 -f h264 -i "$video_out/video.h264" \
  -i "$video_audio" -map 0:v:0 -map 1:a:0 \
  -c:v copy -c:a aac -b:a 256k -movflags +faststart \
  "$video_out/video.mp4"
```

Check the actual `avc1` video and `mp4a` audio entries and media timing. If native
probing cannot open the stub decoder, inspect the MP4 container and use the
working browser for decoding. Load the MP4 in an HTML video element, wait for
metadata, play it, and seek near the beginning, middle, and end. Record dimensions,
duration, readiness, and any media error.

For a verification image, a browser screenshot of the displayed video avoids
the tainted-canvas restriction that can prevent exporting a `file://` video
through `canvas.toDataURL`. Verify audio coverage separately; successful video
playback is not proof that the complete supplied soundtrack was retained.
