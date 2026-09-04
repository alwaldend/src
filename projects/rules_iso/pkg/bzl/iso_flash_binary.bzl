"""Runnable wrapper that flashes an ISO image to a block device."""

def _impl(ctx):
    script = ctx.actions.declare_file("{}.sh".format(ctx.label.name))
    image = ctx.file.image
    runfiles = ctx.runfiles(files = [image])
    image_rlocation = image.short_path.removeprefix("../")
    ctx.actions.write(
        output = script,
        content = """#!/usr/bin/env sh
set -eu
RUNFILES_DIR="${{RUNFILES_DIR:-$0.runfiles}}"
image_runfiles="${{RUNFILES_DIR}}/{image_rlocation}"
device="${{1:-}}"
if [ -z "${{device}}" ]; then
    echo "Usage: bazel run {target} -- /dev/sdX" >&2
    exit 2
fi
if [ ! -b "${{device}}" ]; then
    echo "Not a block device: ${{device}}" >&2
    exit 1
fi
echo "Flashing ${{image_runfiles}} to ${{device}}. This will overwrite ${{device}}." >&2
dd if="${{image_runfiles}}" of="${{device}}" bs=4M conv=fsync status=progress
""".format(
            image_rlocation = image_rlocation,
            target = ctx.label,
        ),
        is_executable = True,
    )
    return [DefaultInfo(
        executable = script,
        files = depset([script]),
        runfiles = runfiles,
    )]

iso_flash_binary = rule(
    implementation = _impl,
    attrs = {
        "image": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
    },
    executable = True,
)
