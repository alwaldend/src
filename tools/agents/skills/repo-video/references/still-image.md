# Still image with a complete soundtrack

Use a rendered PNG or WebP and a finite audio file. Inspect their dimensions,
duration, and stream types first. Choose frame rate and geometry for the user's
destination; 30 fps is a broadly usable example, not a requirement for every
static video. Ensure the chosen pixel format supports the dimensions without
silently cropping the image.

Set absolute `video_image`, `video_audio`, and `video_out` paths for the task.
Create task-local temporary storage through `al-workspace`. The following
example assumes a usable `libx264` encoder; if it is absent or fails, use
[the existing browser encoder](browser-avc.md).

```sh
video_seconds="$(timeout 30 ffprobe -v error -show_entries format=duration \
  -of default=noprint_wrappers=1:nokey=1 "$video_audio")"

timeout 300 ffmpeg -nostdin -hide_banner -loglevel warning -n \
  -f image2 -loop 1 -framerate 30 -i "$video_image" \
  -t "$video_seconds" -an -vf setsar=1 \
  -c:v libx264 -preset medium -tune stillimage -crf 18 -pix_fmt yuv420p \
  "$video_out/video-only.mp4"

timeout 180 ffmpeg -nostdin -hide_banner -loglevel warning -n \
  -i "$video_out/video-only.mp4" -i "$video_audio" \
  -map 0:v:0 -map 1:a:0 -c:v copy -c:a aac -b:a 256k \
  -movflags +faststart "$video_out/video.mp4"
```

Check that `video_seconds` is a usable numeric duration before encoding. If it
is unavailable, establish decoded coverage first. The first command bounds
only the silent video. The second reads the complete audio with no duration
cap or `-shortest`. Use audio stream copy instead of AAC conversion when the
source already meets the output requirements.

A final fraction of a video frame can extend past the audio. If exact endpoint
alignment is required, handle it deliberately after checking decoded coverage;
do not trim the audio to an approximate header duration.

## Verify copied audio

For a stream-copy export, decode both audio tracks to the same PCM format:

```sh
timeout 60 ffmpeg -nostdin -v error -n -i "$video_audio" \
  -map 0:a:0 -vn -c:a pcm_s16le -f hash -hash sha256 \
  "$video_out/source-audio.sha256"
timeout 60 ffmpeg -nostdin -v error -n -i "$video_out/video.mp4" \
  -map 0:a:0 -vn -c:a pcm_s16le -f hash -hash sha256 \
  "$video_out/output-audio.sha256"
cmp "$video_out/source-audio.sha256" "$video_out/output-audio.sha256"
```

A differing hash is a failed preservation check, not an expected container
change. Investigate missing packets, skip samples, and discard padding. Remux
the complete original audio into finite video without a trimming option, then
verify again. These equality checks do not apply to lossy AAC conversion;
check decoded sample coverage and expected codec padding instead.
