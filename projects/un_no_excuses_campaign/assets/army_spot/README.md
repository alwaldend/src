# Army spot — Blender remake

An editable, stylized recreation of the supplied 30-second
`NO-EXCUSES-Spot-11-In-the-Army-Canadian-English-mp4.mp4` cartoon.

- `army_spot.blend` contains the four characters, uniforms, desert setting,
  shot cameras, keyframed motion, and the packed original soundtrack.
- `army_spot.mp4` is the rendered 960×540 H.264 video with the complete
  original AAC soundtrack.
- `poster.png` is a representative frame, and `manifest.json` records file
  identities and the checks performed on the delivered files.

## Fidelity and editability

The remake follows the source's shot cuts and broad action, including the
salute and the two recruits' departure. Facial likeness, fabric construction,
gestures, lip synchronization, and the distant scenery remain approximate.
Independent visual reviews identified these differences; the result is not an
exact replica.

The original video is not used as a background, texture, or visual layer.
The scene's geometry, materials, cameras, and animation remain editable with
automatic scripts disabled. No runtime Python handler is required. The
soundtrack retains the source's rights and attribution requirements.

## Render the maintained scene

Use the repository-pinned Blender. The maintained `.blend` file is the source
for this revision. `//projects/un_no_excuses_campaign:army_spot_scene` generates the earlier
procedural starting point; it does not reproduce the subsequent sculpting,
animation, and staging changes.

The original attachment is needed only to extract the soundtrack for the MP4.
It is not checked in. Run from the repository root, replacing the first path
with the attachment's absolute path:

```sh
set -e
army_reference=/absolute/path/to/NO-EXCUSES-Spot-11-In-the-Army-Canadian-English-mp4.mp4
army_output="$PWD/out/army_spot/rebuild"
mkdir -p "$army_output/frames"
printf '%s  %s\n' \
  4aa8481103870d56131063ed62438f6999e3d0085a5826a647a03b6eb6e7917d \
  "$army_reference" | sha256sum -c -
ffmpeg -nostdin -n -i "$army_reference" -map 0:a:0 -c:a copy -vn \
  "$army_output/original_audio.m4a"
bazel run //third_party/org_blender_blender:blender -- \
  --background --factory-startup --disable-autoexec \
  "$PWD/projects/un_no_excuses_campaign/assets/army_spot/army_spot.blend" \
  --render-output "$army_output/frames/frame_" --render-format PNG \
  --frame-start 1 --frame-end 898 --render-anim
bazel run //projects/un_no_excuses_campaign:encode_frames -- \
  --frames "$army_output/frames" --output "$army_output/silent.mp4" \
  --count 898 --width 960 --height 540 --fps 30 --fps-base 1.001 \
  --hold-frames 4
ffmpeg -nostdin -n -i "$army_output/silent.mp4" \
  -i "$army_output/original_audio.m4a" \
  -map 0:v:0 -map 1:a:0 -c copy -movflags +faststart \
  "$army_output/army_spot.mp4"
```

The scene uses CPU Cycles, eight samples, and denoising. The source has 898
frames at 30000/1001 fps; the encoder holds the final frame for four more
frames to preserve the slightly longer soundtrack. The host H.264 decoder may
be a stub, so Blender's bundled decoder and encoder are used for verification.

## Controlling frames

Frames 1, 79, 120, 198, 262, 350, 429, 648, 688, 717, 778, and 842 mark
the source's hard cuts. Frame 301 controls the main close-up, and frame 569
controls the salute. Rear anatomy is inferred from the available footage.
