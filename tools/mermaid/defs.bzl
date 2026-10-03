"""Build actions for rendering Mermaid source diagrams."""

def _mermaid_svg_impl(ctx):
    source = ctx.file.src
    theme = ctx.file.theme
    output = ctx.outputs.out

    if not output.basename.endswith(".svg"):
        fail("out must name an .svg file, got: {}".format(output.basename))

    outputs = [output]
    if ctx.attr.native and (ctx.attr.plain or not ctx.file.palette or not ctx.outputs.dark_out):
        fail("native rendering requires a palette and dark_out, and cannot be plain")
    if not ctx.attr.native and (ctx.file.palette or ctx.outputs.dark_out):
        fail("palette and dark_out require native rendering")
    plain = ctx.attr.plain
    label_font = ctx.file.label_font or ctx.file._label_font

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
        if ctx.attr.native:
            args.add("--native")
            args.add("--palette", ctx.file.palette.path)
            args.add("--dark-output", ctx.outputs.dark_out.path)
            outputs.append(ctx.outputs.dark_out)

        # The selected label face travels inside the rendered document, so a
        # consumer that can load no webfont still shows the glyphs the
        # renderer measured.
        args.add("--font", label_font.path)

        # The render exposes exactly these directories to Chrome through
        # Fontconfig, so no family resolves from the host. Every file in an
        # exposed directory becomes reachable, so each holds one repository-owned
        # family.
        font_files = ctx.files._body_fonts + ctx.files._label_fonts + [label_font]
        font_directories = {}
        for font in font_files:
            font_directories[font.dirname] = None
        for directory in font_directories:
            args.add("--font-directory", directory)

    inputs = [source]
    if not plain:
        inputs.append(theme)
        if ctx.file.palette:
            inputs.append(ctx.file.palette)
        inputs = inputs + ctx.files._body_fonts + ctx.files._label_fonts + [label_font]

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
        outputs = outputs,
        progress_message = "Rendering Mermaid SVG %{label}",
        tools = [
            ctx.attr._browser[DefaultInfo].files_to_run,
            ctx.attr._render[DefaultInfo].files_to_run,
        ],
    )

    return [DefaultInfo(files = depset(outputs))]

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
        "native": attr.bool(
            default = False,
            doc = "Render light/dark variants using native Mermaid APIs without SVG post-processing.",
        ),
        "palette": attr.label(
            allow_single_file = [".css"],
            doc = "Stylesheet exposing --mermaid-* theme variables in light and dark media.",
        ),
        "dark_out": attr.output(
            doc = "Dark SVG output for native rendering.",
        ),
        "label_font": attr.label(
            allow_single_file = [".ttf"],
            doc = "Optional embedded font override matching the selected theme.",
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

def _mermaid_webp_impl(ctx):
    source = ctx.file.src
    theme = ctx.file.theme
    output = ctx.outputs.out

    if not output.basename.endswith(".webp"):
        fail("out must name a .webp file, got: {}".format(output.basename))
    if ctx.attr.native and not ctx.file.palette:
        fail("native rendering requires a palette")
    if not ctx.attr.native and ctx.file.palette:
        fail("palette requires native rendering")
    if not ctx.attr.native and ctx.attr.color_scheme != "light":
        fail("a dark color_scheme requires native rendering")
    if ctx.attr.scale <= 0:
        fail("scale must be a positive integer")
    if ctx.attr.aspect_ratio and (len(ctx.attr.aspect_ratio) != 2 or min(ctx.attr.aspect_ratio) <= 0):
        fail("aspect_ratio must contain two positive integers")

    label_font = ctx.file.label_font or ctx.file._label_font

    args = ctx.actions.args()
    args.add("--input", source.path)
    args.add("--output", output.path)
    args.add("--browser", ctx.executable._browser.path)
    args.add("--background", "white")
    args.add("--raster")
    args.add("--scale", str(ctx.attr.scale))
    if ctx.attr.aspect_ratio:
        for term in ctx.attr.aspect_ratio:
            args.add("--aspect-ratio", str(term))
    args.add("--config", theme.path)
    args.add("--font", label_font.path)
    if ctx.attr.native:
        args.add("--native")
        args.add("--palette", ctx.file.palette.path)
        args.add("--color-scheme", ctx.attr.color_scheme)

    # The render exposes exactly these directories to Chrome through
    # Fontconfig, so no family resolves from the host. Every file in an
    # exposed directory becomes reachable, so each holds one repository-owned
    # family.
    font_files = ctx.files._body_fonts + ctx.files._label_fonts + [label_font]
    font_directories = {}
    for font in font_files:
        font_directories[font.dirname] = None
    for directory in font_directories:
        args.add("--font-directory", directory)

    inputs = [source, theme] + ctx.files._body_fonts + ctx.files._label_fonts + [label_font]
    if ctx.file.palette:
        inputs.append(ctx.file.palette)

    ctx.actions.run(
        arguments = [args],
        env = {
            "BAZEL_BINDIR": ctx.bin_dir.path,
            # Declared File.path values are relative to the action execroot.
            "JS_BINARY__NO_CD_BINDIR": "1",
        },
        executable = ctx.executable._render,
        inputs = inputs,
        mnemonic = "MermaidWebp",
        outputs = [output],
        progress_message = "Rendering Mermaid WebP %{label}",
        tools = [
            ctx.attr._browser[DefaultInfo].files_to_run,
            ctx.attr._render[DefaultInfo].files_to_run,
        ],
    )

    return [DefaultInfo(files = depset([output]))]

mermaid_webp = rule(
    implementation = _mermaid_webp_impl,
    doc = "Renders one Mermaid source file to a WebP image a service " +
          "upload endpoint accepts, reusing the maintained SVG render " +
          "contract and encoding the raster through the pinned browser.",
    attrs = {
        "src": attr.label(
            allow_single_file = [".mmd"],
            mandatory = True,
            doc = "Mermaid source file.",
        ),
        "out": attr.output(
            mandatory = True,
            doc = "Rendered .webp output.",
        ),
        "scale": attr.int(
            default = 2,
            doc = "Raster scale factor. The image is authored at the " +
                  "diagram's natural size and scaled by this factor, so a " +
                  "node-sized diagram still carries crisp text.",
        ),
        "aspect_ratio": attr.int_list(
            doc = "Optional [width, height] ratio of positive integers. Adds " +
                  "theme-colored padding around the scaled diagram without " +
                  "cropping or stretching; empty keeps the natural canvas.",
        ),
        "native": attr.bool(
            default = False,
            doc = "Encode a native SVG using the consumer's palette without SVG post-processing.",
        ),
        "color_scheme": attr.string(
            default = "light",
            values = ["light", "dark"],
            doc = "Color scheme selected from the native palette; dark requires native rendering.",
        ),
        "palette": attr.label(
            allow_single_file = [".css"],
            doc = "Stylesheet exposing --mermaid-* theme variables; requires native rendering.",
        ),
        "theme": attr.label(
            allow_single_file = [".json"],
            default = "//tools/mermaid:theme.json",
            doc = "Mermaid configuration applied by default; a diagram's own " +
                  "Mermaid directives still override it.",
        ),
        "label_font": attr.label(
            allow_single_file = [".ttf"],
            doc = "Optional embedded font override matching the selected theme.",
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
