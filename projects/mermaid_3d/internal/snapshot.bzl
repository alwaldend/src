"""Offline raster snapshots of generated Mermaid 3D views."""

def _mermaid_3d_snapshot_impl(ctx):
    if ctx.outputs.out.extension not in ["png", "webp"]:
        fail("out must have a .png or .webp extension", attr = "out")
    for name, value in [("width", ctx.attr.width), ("height", ctx.attr.height)]:
        if value < 128 or value > 8192:
            fail("{} must be from 128 to 8192 pixels".format(name), attr = name)
    if ctx.attr.width * ctx.attr.height > 16777216:
        fail("Snapshot dimensions must not exceed 16777216 pixels total")

    args = ctx.actions.args()
    args.add("--input", ctx.file.src.path)
    args.add("--output", ctx.outputs.out.path)
    args.add("--browser", ctx.executable._browser.path)
    args.add("--color-scheme", ctx.attr.color_scheme)
    args.add("--camera", ctx.attr.camera)
    args.add("--width", ctx.attr.width)
    args.add("--height", ctx.attr.height)
    ctx.actions.run(
        executable = ctx.executable._snapshot,
        arguments = [args],
        inputs = [ctx.file.src],
        outputs = [ctx.outputs.out],
        tools = [ctx.attr._snapshot[DefaultInfo].files_to_run, ctx.attr._browser[DefaultInfo].files_to_run],
        env = {"BAZEL_BINDIR": ctx.bin_dir.path, "JS_BINARY__NO_CD_BINDIR": "1"},
        mnemonic = "Mermaid3DSnapshot",
        progress_message = "Capturing spatial diagram %{label}",
    )
    return [DefaultInfo(files = depset([ctx.outputs.out]))]

mermaid_3d_snapshot = rule(
    implementation = _mermaid_3d_snapshot_impl,
    doc = "Encode the existing viewer canvas offline as PNG or WebP, fitted to exact pixel dimensions.",
    attrs = {
        "src": attr.label(allow_single_file = [".html"], mandatory = True, doc = "Self-contained HTML from mermaid_3d."),
        "out": attr.output(mandatory = True, doc = "PNG (lossless) or WebP (browser quality 1)."),
        "color_scheme": attr.string(default = "light", values = ["light", "dark"]),
        "camera": attr.string(default = "isometric", values = ["isometric", "front"], doc = "The viewer's default spatial view or front view."),
        "width": attr.int(default = 1280, doc = "Output width in pixels, 128..8192; at most 16777216 pixels total."),
        "height": attr.int(default = 960, doc = "Output height in pixels, 128..8192; at most 16777216 pixels total."),
        "_snapshot": attr.label(default = "//projects/mermaid_3d/cmd/snapshot", executable = True, cfg = "exec"),
        "_browser": attr.label(default = "@com_google_chrome_headless_shell//:chrome_headless_shell_binary", executable = True, cfg = "exec"),
    },
)
