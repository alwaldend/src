"""Build actions for rendering Mermaid source diagrams."""

def _mermaid_svg_impl(ctx):
    source = ctx.file.src
    theme = ctx.file.theme
    output = ctx.outputs.out

    if not output.basename.endswith(".svg"):
        fail("out must name an .svg file, got: {}".format(output.basename))

    plain = ctx.attr.plain
    label_font = ctx.file._label_font

    args = ctx.actions.args()
    args.add("--input", source.path)
    args.add("--output", output.path)
    args.add("--browser", ctx.executable._browser.path)
    args.add("--background", "white")

    if plain:
        # A plain render asks Mermaid for its own defaults, so the repository
        # theme and the pinned fonts are deliberately not inputs.
        args.add("--plain")
    else:
        args.add("--config", theme.path)

        # The handwriting face travels inside the rendered document, so a
        # consumer that can load no webfont still shows the glyphs the
        # renderer measured.
        args.add("--font", label_font.path)

        # The render exposes exactly these directories to Chrome through
        # Fontconfig, so no family resolves from the host. Every file in an
        # exposed directory becomes reachable, so each holds one repository-owned
        # family.
        font_files = ctx.files._body_fonts + ctx.files._label_fonts
        font_directories = {}
        for font in font_files:
            font_directories[font.dirname] = None
        for directory in font_directories:
            args.add("--font-directory", directory)

    inputs = [source]
    if not plain:
        inputs.append(theme)
        inputs = inputs + ctx.files._body_fonts + ctx.files._label_fonts

    ctx.actions.run(
        arguments = [args],
        env = {
            "BAZEL_BINDIR": ctx.bin_dir.path,
            # Declared File.path values are relative to the action execroot.
            "JS_BINARY__NO_CD_BINDIR": "1",
        },
        executable = ctx.executable._render,
        inputs = inputs,
        mnemonic = "MermaidSvg",
        outputs = [output],
        progress_message = "Rendering Mermaid SVG %{label}",
        tools = [
            ctx.attr._browser[DefaultInfo].files_to_run,
            ctx.attr._render[DefaultInfo].files_to_run,
        ],
    )

    return [DefaultInfo(files = depset([output]))]

mermaid_svg = rule(
    implementation = _mermaid_svg_impl,
    doc = "Renders one Mermaid source file to a declared SVG output.",
    attrs = {
        "src": attr.label(
            allow_single_file = [".mmd"],
            mandatory = True,
            doc = "Mermaid source file.",
        ),
        "out": attr.output(
            mandatory = True,
            doc = "Rendered .svg output.",
        ),
        "plain": attr.bool(
            default = False,
            doc = "Render with Mermaid's own defaults: no repository theme, " +
                  "no paint-order post-processing, and no pinned fonts. The " +
                  "output is then neither themed nor hermetic, and is intended " +
                  "for comparing against the upstream appearance.",
        ),
        "theme": attr.label(
            allow_single_file = [".json"],
            default = "//tools/mermaid:theme.json",
            doc = "Mermaid configuration applied by default; a diagram's own " +
                  "Mermaid directives still override it.",
        ),
        "_label_font": attr.label(
            allow_single_file = [".ttf"],
            default = "//tools/mermaid/fonts:ArchitectsDaughter-Regular.ttf",
            doc = "Handwriting face embedded in the rendered document, so the " +
                  "output keeps its measured text metrics without a webfont.",
        ),
        "_body_fonts": attr.label(
            default = "@com_alwaldend_src_tools_drawio_fonts//:fonts",
            doc = "Pinned fonts the theme falls back to.",
        ),
        "_label_fonts": attr.label(
            default = "//tools/mermaid/fonts:fonts",
            doc = "Pinned handwriting font the theme labels diagrams with.",
        ),
        "_browser": attr.label(
            cfg = "exec",
            default = "@com_google_chrome_headless_shell//:chrome_headless_shell_binary",
            executable = True,
        ),
        "_render": attr.label(
            cfg = "exec",
            default = "//tools/mermaid/cmd/render",
            executable = True,
        ),
    },
)
