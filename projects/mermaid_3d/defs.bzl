"""Hermetic spatial views from authoritative Mermaid definitions."""

load("//projects/mermaid_3d/internal:snapshot.bzl", _mermaid_3d_snapshot = "mermaid_3d_snapshot")
load("//tools/mermaid:theme.bzl", "MermaidThemeInfo")

mermaid_3d_snapshot = _mermaid_3d_snapshot

def _mermaid_3d_impl(ctx):
    theme = ctx.attr.theme[MermaidThemeInfo]
    args = ctx.actions.args()
    args.add("--input", ctx.file.src.path)
    args.add("--output", ctx.outputs.out.path)
    args.add("--model", ctx.outputs.model.path)
    args.add("--browser", ctx.executable._browser.path)
    args.add("--config", theme.config.path)
    args.add("--palette", theme.palette.path)
    args.add("--font", theme.font.path)
    ctx.actions.run(
        executable = ctx.executable._render,
        arguments = [args],
        inputs = [ctx.file.src, theme.config, theme.palette, theme.font, theme.font_license],
        outputs = [ctx.outputs.out, ctx.outputs.model],
        tools = [ctx.attr._render[DefaultInfo].files_to_run, ctx.attr._browser[DefaultInfo].files_to_run],
        env = {"BAZEL_BINDIR": ctx.bin_dir.path, "JS_BINARY__NO_CD_BINDIR": "1"},
        mnemonic = "Mermaid3D",
        progress_message = "Rendering spatial diagram %{label}",
    )
    return [DefaultInfo(files = depset([ctx.outputs.out]))]

mermaid_3d = rule(
    implementation = _mermaid_3d_impl,
    attrs = {
        "src": attr.label(allow_single_file = [".mmd"], mandatory = True),
        "theme": attr.label(providers = [MermaidThemeInfo], mandatory = True),
        "out": attr.output(mandatory = True),
        "model": attr.output(mandatory = True),
        "_render": attr.label(default = "//projects/mermaid_3d/cmd/render", executable = True, cfg = "exec"),
        "_browser": attr.label(default = "@com_google_chrome_headless_shell//:chrome_headless_shell_binary", executable = True, cfg = "exec"),
    },
)
