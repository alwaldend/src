"""Encode a rendered frame sequence using pinned Blender's bundled H.264 codec.

The output is silent. Supply the scene dimensions, timing, and optional
final-frame hold; mux a soundtrack separately when needed.
"""

import argparse
import pathlib
import sys

import bpy


def _main() -> None:
    """Load exactly the requested frame range and export a playable H.264 MP4."""
    parser: argparse.ArgumentParser = argparse.ArgumentParser(
        description=__doc__
    )
    parser.add_argument("--frames", type=pathlib.Path, required=True)
    parser.add_argument("--output", type=pathlib.Path, required=True)
    parser.add_argument("--count", type=int, required=True)
    parser.add_argument("--width", type=int, required=True)
    parser.add_argument("--height", type=int, required=True)
    parser.add_argument("--fps", type=int, default=30)
    parser.add_argument("--fps-base", type=float, default=1.0)
    parser.add_argument("--hold-frames", type=int, default=0)
    arguments: argparse.Namespace = parser.parse_args(
        sys.argv[sys.argv.index("--") + 1 :]
    )
    if (
        min(arguments.count, arguments.width, arguments.height, arguments.fps)
        <= 0
    ):
        parser.error("count, dimensions, and fps must be positive")
    if arguments.fps_base <= 0 or arguments.hold_frames < 0:
        parser.error("fps-base must be positive and hold-frames nonnegative")
    frames: pathlib.Path = arguments.frames.resolve()
    for frame in range(1, arguments.count + 1):
        if not (frames / f"frame_{frame:04}.png").is_file():
            raise FileNotFoundError(f"Missing animation frame {frame}")
    source_image: bpy.types.Image = bpy.data.images.load(
        str(frames / "frame_0001.png"), check_existing=True
    )
    source_width, source_height = source_image.size
    scene: bpy.types.Scene = bpy.context.scene
    scene.render.resolution_x = arguments.width
    scene.render.resolution_y = arguments.height
    scene.render.resolution_percentage = 100
    scene.render.fps = arguments.fps
    scene.render.fps_base = arguments.fps_base
    scene.frame_start = 1
    scene.frame_end = arguments.count + arguments.hold_frames
    scene.render.use_sequencer = True
    editor: bpy.types.SequenceEditor = scene.sequence_editor_create()
    strip: bpy.types.ImageStrip = editor.strips.new_image(
        "Rendered animation",
        str(frames / "frame_0001.png"),
        channel=1,
        frame_start=1,
    )
    strip.transform.scale_x = arguments.width / source_width
    strip.transform.scale_y = arguments.height / source_height
    for frame in range(2, arguments.count + 1):
        strip.elements.append(f"frame_{frame:04}.png")
    strip.frame_final_end = arguments.count + arguments.hold_frames + 1
    scene.view_settings.view_transform = "Standard"
    scene.render.image_settings.media_type = "VIDEO"
    scene.render.ffmpeg.format = "MPEG4"
    scene.render.ffmpeg.codec = "H264"
    scene.render.ffmpeg.constant_rate_factor = "HIGH"
    scene.render.ffmpeg.ffmpeg_preset = "GOOD"
    scene.render.ffmpeg.audio_codec = "NONE"
    scene.render.filepath = str(arguments.output.resolve())
    bpy.ops.render.render(animation=True)


if __name__ == "__main__":
    _main()
